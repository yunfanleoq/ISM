package cjt188

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/goburrow/serial"
)

const (
	// TCPDefaultTimeout TCP Default timeout
	TCPDefaultTimeout = 1 * time.Second
	// TCPDefaultAutoReconnect TCP Default auto reconnect count
	TCPDefaultAutoReconnect = 1
	tcpAduMaxSize           = 1024
)

// TCPClientProvider implements ClientProvider interface.
type CJT188TCPClientProvider struct {
	logger
	address string
	serialPort
	protocolType byte
	mu           sync.Mutex
	// TCP connection
	conn net.Conn
	// Connect & Read timeout
	timeout time.Duration
	// if > 0, when disconnect,it will try to reconnect the remote
	// but if we active close self,it will not to reconnect
	// if == 0 auto reconnect not active
	autoReconnect byte
	// For synchronization between messages of server & client
	transactionID uint32
	// request
	prefix       string
	BeforeEnable bool
	PrefixHandler
}

// check TCPClientProvider implements the interface ClientProvider underlying method
var _ ClientProvider = (*CJT188TCPClientProvider)(nil)

// request pool, all TCP client use this pool

// NewTCPClientProvider allocates a new TCPClientProvider.
func NewCjt188TCPClientProvider(prefix string, modbusType byte, address string, opts ...ClientProviderOption) *CJT188TCPClientProvider {
	p := &CJT188TCPClientProvider{
		address:       address,
		protocolType:  modbusType,
		timeout:       TCPDefaultTimeout,
		autoReconnect: TCPDefaultAutoReconnect,
		prefix:        prefix,
		PrefixHandler: &DefaultPrefix{},
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}
func (sf *CJT188TCPClientProvider) setPrefixHandler(handler PrefixHandler) {
	sf.PrefixHandler = handler
}

// SendAndRead 发送数据并读取返回值
func (sf *CJT188TCPClientProvider) SendAndRead(p *Protocol) (aduResponse []byte, err error) {
	bf := bytes.NewBuffer(make([]byte, 0))
	err = sf.EncodePrefix(bf, sf.BeforeEnable)
	if err != nil {
		return nil, err
	}

	err = p.Encode(bf)
	if err != nil {
		return nil, err
	}
	// sf.Debugf("sending ==> [% x]", bf.Bytes())
	return sf.SendRawFrameAndRead(bf.Bytes())
}
func (sf *CJT188TCPClientProvider) Send(p *Protocol) (err error) {
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
func (sf *CJT188TCPClientProvider) ReadRawFrame() (aduResponse []byte, err error) {

	var timeout time.Time
	if sf.timeout > 0 {
		timeout = time.Now().Add(sf.timeout)
	}
	if err = sf.conn.SetDeadline(timeout); err != nil {
		return nil, err
	}
	head := make([]byte, 1024)
	size, err := sf.conn.Read(head[:])
	if err != nil {
		return nil, err
	}
	var find_before_index int = 0
	for k, v := range head {
		if v == 0x68 {
			find_before_index = k
			break
		}
	}
	if size < 5 {
		return nil, fmt.Errorf("read length error")
	}
	if (size - find_before_index) < 10 {
		return nil, fmt.Errorf("read length error <10")
	}
	checkBuffer := head[find_before_index : size-2]
	var cs = 0
	for _, b := range checkBuffer {
		cs += int(b)
	}
	checkcs := head[size-2]
	if checkcs != byte(cs) {
		return nil, fmt.Errorf("check error %d,%d", checkcs, byte(cs))
	}
	packend := head[size-1]
	if packend != 0x16 {
		return nil, fmt.Errorf("package end char  error %x", packend)
	}
	sf.Debug(sf.prefix, "recv head ==> [% x]\n", head[find_before_index:size])

	return head[find_before_index:size], nil
}
func (sf *CJT188TCPClientProvider) SendRawFrameAndRead(aduRequest []byte) (aduResponse []byte, err error) {
	sf.mu.Lock()
	defer sf.mu.Unlock()
	// if err = sf.connect(); err != nil {
	// 	return
	// }
	err = sf.SendRawFrame(aduRequest)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	return sf.ReadRawFrame()
}
func (sf *CJT188TCPClientProvider) SendRawFrame(aduRequest []byte) (err error) {
	// if err = sf.connect(); err != nil {
	// 	return
	// }
	var timeout time.Time
	if sf.timeout > 0 {
		timeout = time.Now().Add(sf.timeout)
	}
	if err = sf.conn.SetDeadline(timeout); err != nil {
		return err
	}
	// Send the request
	sf.Debug(sf.prefix, "sending ==> [% x]", aduRequest)
	//发送数据
	_, err = sf.conn.Write(aduRequest)
	return err
}

// Connect establishes a new connection to the address in Address.
// Connect and Close are exported so that multiple requests can be done with one session
func (sf *CJT188TCPClientProvider) Connect() error {
	sf.mu.Lock()
	err := sf.connect()
	sf.mu.Unlock()
	return err
}

// Caller must hold the mutex before calling this method.
func (sf *CJT188TCPClientProvider) connect() error {
	dialer := &net.Dialer{Timeout: sf.timeout}
	conn, err := dialer.Dial("tcp", sf.address)
	if err != nil {
		return err
	}
	sf.conn = conn
	return nil
}

// IsConnected returns a bool signifying whether
// the client is connected or not.
func (sf *CJT188TCPClientProvider) IsConnected() bool {
	sf.mu.Lock()
	b := sf.isConnected()
	sf.mu.Unlock()
	return b
}

// Caller must hold the mutex before calling this method.
func (sf *CJT188TCPClientProvider) isConnected() bool {
	return sf.conn != nil
}

// SetAutoReconnect set auto reconnect  retry count
func (sf *CJT188TCPClientProvider) SetAutoReconnect(cnt byte) {
	sf.mu.Lock()
	sf.autoReconnect = cnt
	if sf.autoReconnect > 6 {
		sf.autoReconnect = 6
	}
	sf.mu.Unlock()
}

// Close closes current connection.
func (sf *CJT188TCPClientProvider) Close() (err error) {
	sf.mu.Lock()
	if sf.conn != nil {
		err = sf.conn.Close()
		sf.conn = nil
	}
	sf.mu.Unlock()
	return
}

func (sf *CJT188TCPClientProvider) SetSerialConfig(serial.Config) {}

func (sf *CJT188TCPClientProvider) SetTCPTimeout(t time.Duration) {
	sf.timeout = t
}

// flush flushes pending data in the connection,
// returns io.EOF if connection is closed.
func (sf *CJT188TCPClientProvider) flush(b []byte) (err error) {
	if err = sf.conn.SetReadDeadline(time.Now()); err != nil {
		return
	}
	// timeout setting will be reset when reading
	if _, err = sf.conn.Read(b); err != nil {
		// Ignore timeout error
		if netError, ok := err.(net.Error); ok && netError.Timeout() {
			err = nil
		}
	}
	return
}
func (sf *CJT188TCPClientProvider) SetConnect(conn net.Conn) error {

	return nil
}
func (sf *CJT188TCPClientProvider) SetBeforeCode(enable bool) error {
	sf.BeforeEnable = enable
	return nil
}
