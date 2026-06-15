package iec104

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var closeLock sync.Mutex

func NewClient(option *ClientOption) *Client {
	return &Client{
		ClientOption: option,

		org:         ORG(0),
		coa:         COA(0x0001),
		status:      0,
		failedTimes: 0,
		ReadTimeOut: 15000,
		Retries:     5,
		waitGroup:   &sync.WaitGroup{},
		sendChan:    make(chan []byte, 1),
		recvChan:    make(chan *APDU, 2000),
		dataChan:    make(chan *APDU, 2000),
		cmdRspChan:  make(chan *cmdRsp, 2000),
	}
}

// Client in IEC 104 is also called as master or controlling station.
type Client struct {
	*ClientOption
	conn                 net.Conn // network channel with the iec104 substation/server
	waitGroup            *sync.WaitGroup
	cancel               context.CancelFunc
	sendChan             chan []byte // send data to server
	recvChan             chan *APDU  // receive apdu from server
	dataChan             chan *APDU  // make Client owner to handle data received from server by themselves
	cmdRspChan           chan *cmdRsp
	Retries              int
	org                  ORG    // originator address to identify controlling station when there are multiple controlling stations
	coa                  COA    // common address (or station address)
	ssn, rsn, recvICount uint16 // send sequence number, receive sequence number
	ifn                  uint16 // i-format frame number (for send S-frame data regularity)
	failedTimes          int
	ReadTimeOut          int

	status int32 // initial, connected, disconnected
}

func (c *Client) Connect() error {
	if err := c.dial(); err != nil {
		return err
	}

	// After the establishment of a TCP connection, send and receive sequence number should be set to zero.
	c.ssn, c.rsn, c.recvICount = 0, 0, 0

	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel
	go c.writingToSocket(ctx)
	go c.readingFromSocket(ctx)
	go c.handlingData(ctx)
	c.waitGroup.Add(3)
	c.onConnectHandler(c)
	c.handler.OnConnectHandler(c)
	return nil
}
func (c *Client) dial() (err error) {
	schema, address, timeout := c.server.Scheme, c.server.Host, c.connectTimeout
	var conn net.Conn
	switch schema {
	case "tcp":
		conn, err = net.DialTimeout("tcp", address, timeout)
	case "ssl", "tls", "tcps":
		conn, err = tls.DialWithDialer(&net.Dialer{Timeout: timeout}, "tcp", address, c.tc)
	default:
		return fmt.Errorf("unknown schema: %s", schema)
	}
	if err != nil {
		return err
	}

	c.status = 1
	c.conn = conn
	return
}

func (c *Client) writingToSocket(ctx context.Context) {
	_lg.Debugf("start goroutine for writing to socket")
	defer func() {
		_lg.Debugf("stop goroutine for writing to socket")
	}()

	for {
		select {
		case <-ctx.Done():
			c.waitGroup.Done()
			return
		case data := <-c.sendChan:
			if c.status != 1 { // 加这一句
				continue
			}
			var timeout time.Time
			if c.ReadTimeOut > 0 {
				timeout = time.Now().Add(time.Millisecond * 5000)
			}
			c.conn.SetWriteDeadline(timeout)
			if _, err := c.conn.Write(data); err != nil {
				c.Close()
				fmt.Println("write to socket:", err.Error())
			}
		}
	}
}
func (c *Client) readingFromSocket(ctx context.Context) {
	_lg.Debugf("start goroutine for reading from socket")
	defer func() {
		_lg.Debugf("stop goroutine for reading from socket")
	}()

	for {
		select {
		case <-ctx.Done():
			c.waitGroup.Done()
			return
		default:
			if c.status != 1 { // 加这一句
				time.Sleep(100 * time.Millisecond)
				continue
			}
			apdu, err := c.readFromSocket(ctx)
			if err != nil {
				fmt.Println("read from socket:", err)
				c.failedTimes++
				if c.failedTimes >= c.Retries {
					c.failedTimes = 0
					c.Close()
				}
				continue
			}
			c.failedTimes = 0
			switch apdu.frame.Type() {
			case FrameTypeU:
				uFrame, ok := apdu.frame.(*UFrame)
				if ok {
					switch uFrame.Cmd[0] {
					case UFrameFunctionStartDTA[0]:
						_lg.Debugf("receive u frame: StartDTA")
					case UFrameFunctionStartDTC[0]:
						_lg.Debugf("receive u frame: StartDTC")
						c.recvChan <- apdu
					case UFrameFunctionStopDTA[0]:
						_lg.Debugf("receive u frame: StopDTA")
					case UFrameFunctionStopDTC[0]:
						_lg.Debugf("receive u frame: StopDTC")
						c.recvChan <- apdu
					case UFrameFunctionTestFA[0]:
						_lg.Debugf("receive u frame: TestFA")
						c.sendUFrame(UFrameFunctionTestFC)
					case UFrameFunctionTestFC[0]:
						_lg.Debugf("receive u frame: TestFC")
						c.sendUFrame(UFrameFunctionTestFC)
					}
				}
			}
		}
	}
}
func (c *Client) readFromSocket(ctx context.Context) (*APDU, error) {

	apduLen, err := c.readApduHeader()
	if err != nil {
		return nil, err
	}

	apdu, err := c.readApduBody(apduLen)
	if err != nil {
		return nil, err
	}
	return apdu, nil
}

func (c *Client) handlingData(ctx context.Context) {
	_lg.Debugf("start goroutine for handling data received from server")
	defer func() {
		_lg.Debugf("stop goroutine for handling data received from server")
	}()

	for {
		select {
		case <-ctx.Done():
			c.waitGroup.Done()
			return
		case apdu := <-c.dataChan:
			if err := c.handleData(apdu); err != nil {
				_lg.Warnf("handle iFrame, got: %v", err)
			}
		}
	}
}
func (c *Client) handleData(apdu *APDU) error {
	defer func() {
		if err := recover(); err != nil {
			_lg.Errorf("client handler: %+v", err)
		}
	}()

	_lg.Debugf("handle iFrame: TypeID: %X, COT: %X", apdu.ASDU.typeID, apdu.ASDU.cot)

	switch apdu.typeID {
	case CIcNa1:
		return c.handler.GeneralInterrogationHandler(apdu)
	case CCiNa1:
		return c.handler.CounterInterrogationHandler(apdu)
	case CRdNa1:
		return c.handler.ReadCommandHandler(apdu)
	case CCsNa1:
		return c.handler.ClockSynchronizationHandler(apdu)
	case CTsNb1, CTsTa1:
		return c.handler.TestCommandHandler(apdu)
	case CRpNc1:
		return c.handler.ResetProcessCommandHandler(apdu)
	case CCdNa1:
		return c.handler.DelayAcquisitionCommandHandler(apdu)
	default:
		return c.handler.APDUHandler(apdu)
	}
}

// readApduHeader reads both startByte and apduLen, and returns apduLen
func (c *Client) readApduHeader() (uint8, error) { //
	buf := make([]byte, 2)
	var timeout time.Time
	if c.ReadTimeOut > 0 {
		timeout = time.Now().Add(time.Millisecond * time.Duration(c.ReadTimeOut))
	}

	c.conn.SetDeadline(timeout)
	n, err := c.conn.Read(buf)
	if err != nil {
		return 0, err
	}
	if n != 2 {
		return 0, errors.New("invalid data: empty")
	} else if buf[0] != startByte {
		return 0, fmt.Errorf("invalid data: unexpected start - % X, expected start - % X", buf[0], startByte)
	}
	return buf[1], nil
}
func (c *Client) readApduBody(apduLen uint8) (*APDU, error) {
	apduData := make([]byte, apduLen)
	var timeout time.Time
	if c.ReadTimeOut > 0 {
		timeout = time.Now().Add(time.Millisecond * time.Duration(c.ReadTimeOut))
	}

	c.conn.SetDeadline(timeout)
	n, err := c.conn.Read(apduData)
	if err != nil {
		return nil, err
	}

	for n < int(apduLen) {
		bufLen := int(apduLen) - n
		buf := make([]byte, bufLen)

		if c.ReadTimeOut > 0 {
			timeout = time.Now().Add(time.Millisecond * time.Duration(c.ReadTimeOut))
		}
		c.conn.SetDeadline(timeout)

		m, err := c.conn.Read(buf)
		if err != nil {
			return nil, err
		}
		apduData = append(apduData[:n], buf[:m]...)
		n = len(apduData)
	}
	_lg.Debugf("receive: [% X]", append([]byte{startByte, apduLen}, apduData...))

	apdu := new(APDU)
	if err := apdu.Parse(apduData); err != nil {
		return nil, err
	}

	switch apdu.frame.Type() {
	case FrameTypeI:
		if apdu.ASDU.cmdRsp != nil {
			select {
			case c.cmdRspChan <- apdu.ASDU.cmdRsp:
			default:
				// 通道满了就丢弃，不阻塞接收协程
			}
		}
		if apdu.ASDU.toBeHandled {
			c.dataChan <- apdu
		}
		var TCf2 = uint16(apdu.APCI.Cf2)
		var TCf1 = uint16(apdu.APCI.Cf1)
		c.rsn = ((uint16(TCf1|TCf2<<8))/2 + 1) * 2
		// c.recvICount++
		// if c.recvICount >= 7 {
		// 	c.recvICount = 0
		// 	c.SendTestFrame()
		// }
		if apdu.ASDU.sendSFrame {
			c.SendTestFrame()
		}
		//	c.incRsn()
	}

	return apdu, nil
}
func (c *Client) WaitExitAllGoroutine() {
	c.waitGroup.Wait()
}
func (c *Client) IsConnected() int32 {
	return c.status
}

func (c *Client) Close() {

	closeLock.Lock()
	defer closeLock.Unlock()

	c.status = 0
	if c.cancel != nil {
		c.cancel()
	}
	c.onDisconnectHandler(c)
	c.handler.OnDisConnectHandler(c)
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Client) SendGeneralInterrogation() {
	ios := []*InformationObject{
		{
			ioa: 0x000000,
			ies: []*InformationElement{
				{
					Format: []InformationElementType{QOI},
					Raw:    []byte{0x14},
				},
			},
		},
	}
	c.SendIFrame(&ASDU{
		typeID: CIcNa1,
		sq:     false,
		nObjs:  NOO(len(ios)),
		t:      false,
		cot:    CotAct,
		ios:    ios,
	})
}

func (c *Client) SendCounterInterrogation() {
	ios := []*InformationObject{
		{
			ioa: 0x000000,
			ies: []*InformationElement{
				{
					Format: []InformationElementType{QCC},
					Raw:    []byte{0x45},
				},
			},
		},
	}
	c.SendIFrame(&ASDU{
		typeID: CCiNa1,
		sq:     false,
		nObjs:  NOO(len(ios)),
		t:      false,
		cot:    CotAct,
		ios:    ios,
	})
}

func (c *Client) SendClockSynchronizationCmd(year byte, month byte, day byte, hour byte, min byte, mis uint16) {
	ie := &InformationElement{
		Format: []InformationElementType{SCO},
	}

	iBytes := serializeBigEndianUint16(uint16(mis))
	ie.Raw = append(ie.Raw, iBytes...)

	iBytesMin := []byte{byte(min)}
	ie.Raw = append(ie.Raw, iBytesMin...)

	iBytesHour := []byte{byte(hour)}
	ie.Raw = append(ie.Raw, iBytesHour...)

	iBytesDay := []byte{byte(day)}
	ie.Raw = append(ie.Raw, iBytesDay...)

	iBytesMonth := []byte{byte(month)}
	ie.Raw = append(ie.Raw, iBytesMonth...)

	iBytesYear := []byte{byte(year)}
	ie.Raw = append(ie.Raw, iBytesYear...)

	ios := []*InformationObject{
		{
			ioa: 0x000000,
			ies: []*InformationElement{ie},
		},
	}
	c.SendIFrame(&ASDU{
		typeID: CCsNa1,
		sq:     false,
		nObjs:  NOO(len(ios)),
		t:      false,
		cot:    CotAct,
		ios:    ios,
	})
}

func (c *Client) SendSingleCommand(address IOA, close bool) error {
	if atomic.LoadInt32(&c.status) != 1 {
		return errors.New("connection is closed, skip send command")
	}

	// select
	ie := &InformationElement{
		Format: []InformationElementType{SCO},
	}
	if close {
		ie.Raw = []byte{0x81}
	} else {
		ie.Raw = []byte{0x80}
	}
	ios := []*InformationObject{
		{
			ioa: address,
			ies: []*InformationElement{ie},
		},
	}
	c.SendIFrame(&ASDU{
		typeID: CScNa1,
		sq:     false,
		nObjs:  NOO(len(ios)),
		t:      false,
		cot:    CotAct,
		ios:    ios,
	})
	select {
	case rsp := <-c.cmdRspChan:
		if rsp.err != nil {
			return rsp.err
		}
	case <-time.After(time.Duration(c.ReadTimeOut) * time.Millisecond): // 错误修复
		return errors.New("timeout")
	}

	// execute
	ie = &InformationElement{
		Format: []InformationElementType{SCO},
	}
	if close {
		ie.Raw = []byte{0x01}
	} else {
		ie.Raw = []byte{0x00}
	}
	ios = []*InformationObject{
		{
			ioa: address,
			ies: []*InformationElement{ie},
		},
	}
	c.SendIFrame(&ASDU{
		typeID: CScNa1,
		sq:     false,
		nObjs:  NOO(len(ios)),
		t:      false,
		cot:    CotAct,
		ios:    ios,
	})
	select {
	case rsp := <-c.cmdRspChan:
		if rsp.err != nil {
			return rsp.err
		}
	case <-time.After(time.Duration(c.ReadTimeOut) * time.Millisecond): // 错误修复
		return errors.New("timeout")
	}
	return nil
}
func (c *Client) SendSetPointFloatCommand(address IOA, PointSetValue float32) error {
	// select
	if atomic.LoadInt32(&c.status) != 1 {
		return errors.New("connection is closed, skip send command")
	}

	ie := &InformationElement{
		Format: []InformationElementType{QOS},
	}

	iBytes := serializeBigEndianFloat32(PointSetValue)
	ie.Raw = append(ie.Raw, iBytes...)

	RawBytes := make([]byte, 1, 1)
	RawBytes[0] = 0x00
	ie.Raw = append(ie.Raw, RawBytes...)

	ios := []*InformationObject{
		{
			ioa: address,
			ies: []*InformationElement{ie},
		},
	}
	c.SendIFrame(&ASDU{
		typeID: CSeNc1,
		sq:     false,
		nObjs:  NOO(len(ios)),
		t:      false,
		cot:    CotReq,
		ios:    ios,
	})
	// select {
	// case rsp := <-c.cmdRspChan:
	// 	if rsp.err != nil {
	// 		return rsp.err
	// 	}
	// }
	// execute
	// ie = &InformationElement{
	// 	Format: []InformationElementType{QOS},
	// }

	// RawBytes1 := make([]byte, 1, 1)
	// RawBytes1[0] = 0x02
	// ie.Raw = append(ie.Raw, RawBytes1...)
	// iBytes1 := serializeBigEndianFloat32(PointSetValue)
	// ie.Raw = append(ie.Raw, iBytes1...)

	// ios = []*InformationObject{
	// 	{
	// 		ioa: address,
	// 		ies: []*InformationElement{ie},
	// 	},
	// }
	// c.SendIFrame(&ASDU{
	// 	typeID: CSeNc1,
	// 	sq:     false,
	// 	nObjs:  NOO(len(ios)),
	// 	t:      false,
	// 	cot:    CotActCon,
	// 	ios:    ios,
	// })
	// select {
	// case rsp := <-c.cmdRspChan:
	// 	if rsp.err != nil {
	// 		return rsp.err
	// 	}
	// }
	return nil
}
func (c *Client) SendSetPointInt16Command(address IOA, PointSetValue int16) error {
	// select
	if atomic.LoadInt32(&c.status) != 1 {
		return errors.New("connection is closed, skip send command")
	}

	ie := &InformationElement{
		Format: []InformationElementType{QOS},
	}

	iBytes := serializeLittleEndianInt16(PointSetValue)
	ie.Raw = append(ie.Raw, iBytes...)

	RawBytes := make([]byte, 1, 1)
	RawBytes[0] = 0x00
	ie.Raw = append(ie.Raw, RawBytes...)

	ios := []*InformationObject{
		{
			ioa: address,
			ies: []*InformationElement{ie},
		},
	}
	c.SendIFrame(&ASDU{
		typeID: CSeNb1,
		sq:     false,
		nObjs:  NOO(len(ios)),
		t:      false,
		cot:    CotAct,
		ios:    ios,
	})
	// select {
	// case rsp := <-c.cmdRspChan:
	// 	if rsp.err != nil {
	// 		return rsp.err
	// 	}
	// }
	// execute
	// ie = &InformationElement{
	// 	Format: []InformationElementType{QOS},
	// }

	// RawBytes1 := make([]byte, 1, 1)
	// RawBytes1[0] = 0x02
	// ie.Raw = append(ie.Raw, RawBytes1...)
	// iBytes1 := serializeBigEndianFloat32(PointSetValue)
	// ie.Raw = append(ie.Raw, iBytes1...)

	// ios = []*InformationObject{
	// 	{
	// 		ioa: address,
	// 		ies: []*InformationElement{ie},
	// 	},
	// }
	// c.SendIFrame(&ASDU{
	// 	typeID: CSeNc1,
	// 	sq:     false,
	// 	nObjs:  NOO(len(ios)),
	// 	t:      false,
	// 	cot:    CotActCon,
	// 	ios:    ios,
	// })
	// select {
	// case rsp := <-c.cmdRspChan:
	// 	if rsp.err != nil {
	// 		return rsp.err
	// 	}
	// }
	return nil
}

// 双点遥控（断路器/刀闸 标准双控：选择 → 执行 → 终止）
func (c *Client) SendDoubleCommand(address IOA, close bool) error {
	if atomic.LoadInt32(&c.status) != 1 {
		return errors.New("connection is closed, skip send command")
	}

	// ====================== 1. 遥控选择（SELECT）S/E=1 ======================
	selectCmd := byte(0x81)
	if close {
		selectCmd = 0x82 // 合闸: 10 → 0x82
	}

	ieSelect := &InformationElement{Format: []InformationElementType{DCO}, Raw: []byte{selectCmd}}
	iosSelect := []*InformationObject{{ioa: address, ies: []*InformationElement{ieSelect}}}

	c.SendIFrame(&ASDU{
		typeID: CDcNa1, // 双点遥控 TypeID=46
		sq:     false,
		nObjs:  1,
		t:      false,
		cot:    CotAct, // Cause=6 激活
		ios:    iosSelect,
	})

	// 等待返校
	select {
	case rsp := <-c.cmdRspChan:
		if rsp.err != nil {
			return fmt.Errorf("选择失败: %w", rsp.err)
		}
	case <-time.After(time.Duration(c.ReadTimeOut) * time.Millisecond):
		return errors.New("选择超时")
	}

	// ====================== 2. 遥控执行（EXECUTE）S/E=0 ======================
	execCmd := byte(0x01)
	if close {
		execCmd = 0x02
	}

	ieExec := &InformationElement{Format: []InformationElementType{DCO}, Raw: []byte{execCmd}}
	iosExec := []*InformationObject{{ioa: address, ies: []*InformationElement{ieExec}}}

	c.SendIFrame(&ASDU{
		typeID: CDcNa1,
		sq:     false,
		nObjs:  1,
		t:      false,
		cot:    CotAct, // Cause=6
		ios:    iosExec,
	})

	// 等待执行返校
	select {
	case <-c.cmdRspChan:
	case rsp := <-c.cmdRspChan:
		if rsp.err != nil {
			return fmt.Errorf("选择失败: %w", rsp.err)
		}
	case <-time.After(time.Duration(c.ReadTimeOut) * time.Millisecond):
		return errors.New("执行超时")
	}

	// ====================== 3. 激活终止 ACTTERM（必须发！释放锁）======================
	// 空命令体 + Cause=10 就是终止
	ieTerm := &InformationElement{Format: []InformationElementType{DCO}, Raw: []byte{0x00}}
	iosTerm := []*InformationObject{{ioa: address, ies: []*InformationElement{ieTerm}}}

	c.SendIFrame(&ASDU{
		typeID: CDcNa1,
		sq:     false,
		nObjs:  1,
		t:      false,
		cot:    CotActTerm, // <-- 关键：Cause=10 激活终止
		ios:    iosTerm,
	})

	// 清理通道残留帧
	select {
	case rsp := <-c.cmdRspChan:
		if rsp.err != nil {
			return fmt.Errorf("选择失败: %w", rsp.err)
		}
	case <-time.After(time.Duration(c.ReadTimeOut) * time.Millisecond):
		return errors.New("确认超时")
	}

	return nil
}

func (c *Client) SendIFrame(asdu *ASDU) {
	apci := &IFrame{
		SendSN: c.ssn,
		RecvSN: c.rsn / 2,
	}
	asdu.org = c.org
	asdu.coa = c.coa
	c.sendIFrame(apci, asdu)
}

func (c *Client) sendIFrame(apci *IFrame, asdu *ASDU) {
	c.incSsn()

	frame := c.buildFrame(append(apci.Data(), asdu.Data()...))
	_lg.Debugf("send i frame: [% X]", frame)
	c.sendChan <- frame
}

func (c *Client) SendTestFrame() {
	c.sendSFrame(&SFrame{
		RecvSN: c.rsn,
	})
}

func (c *Client) SendTestUFrame() {
	c.sendUFrame(UFrameFunctionTestFA)
}
func (c *Client) sendSFrame(x *SFrame) {
	frame := c.buildFrame(x.Data())
	_lg.Debugf("send s frame: [% X]", frame)
	c.sendChan <- frame
}

func (c *Client) sendUFrame(x UFrameFunction) {
	name := ""
	frame := c.buildFrame(x)
	switch x[0] {
	case UFrameFunctionStartDTA[0]:
		name = "StartDTA"
	case UFrameFunctionStartDTC[0]:
		name = "StartDTC"
	case UFrameFunctionStopDTA[0]:
		name = "StopDTA"
	case UFrameFunctionStopDTC[0]:
		name = "StopDTC"
	case UFrameFunctionTestFA[0]:
		name = "TestFA"
	case UFrameFunctionTestFC[0]:
		name = "TestFC"
	}
	_lg.Debugf("send u frame: %s - [% X]", name, frame)
	c.sendChan <- frame
}

func (c *Client) buildFrame(data []byte) []byte {
	frame := make([]byte, 0, 0)
	iBytes := serializeBigEndianUint16(uint16(len(data)))
	frame = append(frame, startByte)
	frame = append(frame, iBytes[1])
	frame = append(frame, data...)
	return frame
}

func (c *Client) incRsn() {
	c.rsn++
	if c.rsn == 1<<15 {
		c.rsn = 0
	}
}

func (c *Client) incSsn() {
	c.ssn++
	// 错误修复：使用 ssn 判断
	if c.ssn == 1<<15 {
		c.ssn = 0
	}
}
