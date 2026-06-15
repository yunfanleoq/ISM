package cjt188

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"time"
)

var _ ClientProvider = (*RTUClientProvider)(nil)

type RTUClientProvider struct {
	serialPort
	logger
	PrefixHandler
	IsPrefix     bool
	BeforeEnable bool
	prefix       string
}

func (sf *RTUClientProvider) setPrefixHandler(handler PrefixHandler) {
	sf.PrefixHandler = handler
}

// SendAndRead 发送数据并读取返回值
func (sf *RTUClientProvider) SendAndRead(p *Protocol) (aduResponse []byte, err error) {
	bf := bytes.NewBuffer(make([]byte, 0))
	err = sf.EncodePrefix(bf, false)
	if err != nil {
		return nil, err
	}

	err = p.Encode(bf)
	if err != nil {
		return nil, err
	}
	return sf.SendRawFrameAndRead(bf.Bytes())
}
func (sf *RTUClientProvider) Send(p *Protocol) (err error) {
	bf := bytes.NewBuffer(make([]byte, 0))
	err = sf.EncodePrefix(bf, sf.BeforeEnable)
	if err != nil {
		return err
	}
	err = p.Encode(bf)
	if err != nil {
		return err
	}
	return sf.SendRawFrame(bf.Bytes())
}

// ReadRawFrame 读取返回数据
func (sf *RTUClientProvider) ReadRawFrame() (aduResponse []byte, err error) {
	readByte := make([]byte, 1024)
	readByteSize, _ := io.ReadFull(sf.port, readByte)
	if readByteSize > 0 {
		sf.Debug(sf.prefix, "rec <==[% x]", readByte[:readByteSize])
		var find_before_index int = 0
		for k, v := range readByte {
			if v == 0x68 {
				find_before_index = k
				break
			}
		}
		if readByteSize < 5 {
			return nil, fmt.Errorf("read length < 5")
		}
		if (readByteSize - find_before_index) < 10 {
			return nil, fmt.Errorf("read length error <10")
		}
		checkBuffer := readByte[find_before_index : readByteSize-2]
		var cs = 0
		for _, b := range checkBuffer {
			cs += int(b)
		}
		checkcs := readByte[readByteSize-2]
		if checkcs != byte(cs) {
			return nil, fmt.Errorf("check error %d,%d", checkcs, byte(cs))
		}
		packend := readByte[readByteSize-1]
		if packend != 0x16 {
			return nil, fmt.Errorf("package end char  error %x", packend)
		}

		return readByte[find_before_index:readByteSize], nil
	} else {
		return readByte, fmt.Errorf("time out")
	}
}
func (sf *RTUClientProvider) SendRawFrameAndRead(aduRequest []byte) (aduResponse []byte, err error) {
	sf.mu.Lock()
	defer sf.mu.Unlock()
	if err = sf.connect(); err != nil {
		return
	}
	err = sf.SendRawFrame(aduRequest)
	if err != nil {
		_ = sf.close()
		return
	}
	return sf.ReadRawFrame()
}
func (sf *RTUClientProvider) SendRawFrame(aduRequest []byte) (err error) {
	if err = sf.connect(); err != nil {
		return
	}
	// Send the request
	sf.Debug(sf.prefix, "sending ==> [% x]", aduRequest)
	//发送数据
	_, err = sf.port.Write(aduRequest)
	return err
}

// NewRTUClientProvider allocates and initializes a RTUClientProvider.
// it will use default /dev/ttyS0 19200 8 1 N and timeout 1000
func NewRTUClientProvider(prefix string, opts ...ClientProviderOption) *RTUClientProvider {
	p := &RTUClientProvider{
		prefix:        prefix,
		PrefixHandler: &DefaultPrefix{},
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// calculateDelay roughly calculates time needed for the next frame.
// See MODBUS over Serial Line - Specification and Implementation Guide (page 13).
func (sf *RTUClientProvider) calculateDelay(chars int) time.Duration {
	var characterDelay, frameDelay int // us

	if sf.BaudRate <= 0 || sf.BaudRate > 19200 {
		characterDelay = 750
		frameDelay = 1750
	} else {
		characterDelay = 15000000 / sf.BaudRate
		frameDelay = 35000000 / sf.BaudRate
	}
	return time.Duration(characterDelay*chars+frameDelay) * time.Microsecond
}
func (sf *RTUClientProvider) SetConnect(conn net.Conn) error {

	return nil
}

func (sf *RTUClientProvider) SetTCPTimeout(t time.Duration) {

}
func (sf *RTUClientProvider) SetBeforeCode(enable bool) error {
	sf.BeforeEnable = enable
	return nil
}
