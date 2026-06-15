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
	TCPServerDefaultTimeout = 1 * time.Second
	// TCPDefaultAutoReconnect TCP Default auto reconnect count
	TCPServerDefaultAutoReconnect = 1
	tcpServerAduMaxSize           = 1024
)

// TCPClientProvider implements ClientProvider interface.
type Dlt645TCPServerProvider struct {
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
var _ ClientProvider = (*Dlt645TCPServerProvider)(nil)

// request pool, all TCP client use this pool

// NewTCPClientProvider allocates a new TCPClientProvider.
func NewDlt645TCPServerProvider(prefix string, modbusType byte, address string, opts ...ClientProviderOption) *Dlt645TCPServerProvider {
	p := &Dlt645TCPServerProvider{
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
func (sf *Dlt645TCPServerProvider) setPrefixHandler(handler PrefixHandler) {
	sf.PrefixHandler = handler
}

// SendAndRead 发送数据并读取返回值
func (sf *Dlt645TCPServerProvider) SendAndRead(p *Protocol) (aduResponse []byte, err error) {
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
func (sf *Dlt645TCPServerProvider) Send(p *Protocol) (err error) {
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
func (sf *Dlt645TCPServerProvider) ReadRawFrame() (aduResponse []byte, err error) {

	var timeout time.Time
	if sf.timeout > 0 {
		timeout = time.Now().Add(sf.timeout)
	} else {
		timeout = time.Now().Add(1000)
	}
	if err = sf.conn.SetDeadline(timeout); err != nil {
		return nil, err
	}
	head := make([]byte, 200)
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
func (sf *Dlt645TCPServerProvider) SendRawFrameAndRead(aduRequest []byte) (aduResponse []byte, err error) {
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
func (sf *Dlt645TCPServerProvider) SendRawFrame(aduRequest []byte) (err error) {
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
func (sf *Dlt645TCPServerProvider) Connect() error {
	sf.mu.Lock()
	err := sf.connect()
	sf.mu.Unlock()
	return err
}

// Caller must hold the mutex before calling this method.
func (sf *Dlt645TCPServerProvider) connect() error {
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
func (sf *Dlt645TCPServerProvider) IsConnected() bool {
	sf.mu.Lock()
	b := sf.isConnected()
	sf.mu.Unlock()
	return b
}

// Caller must hold the mutex before calling this method.
func (sf *Dlt645TCPServerProvider) isConnected() bool {
	return sf.conn != nil
}

// SetAutoReconnect set auto reconnect  retry count
func (sf *Dlt645TCPServerProvider) SetAutoReconnect(cnt byte) {
	sf.mu.Lock()
	sf.autoReconnect = cnt
	if sf.autoReconnect > 6 {
		sf.autoReconnect = 6
	}
	sf.mu.Unlock()
}

// Close closes current connection.
func (sf *Dlt645TCPServerProvider) Close() (err error) {
	sf.mu.Lock()
	if sf.conn != nil {
		err = sf.conn.Close()
		sf.conn = nil
	}
	sf.mu.Unlock()
	return
}

func (sf *Dlt645TCPServerProvider) SetSerialConfig(serial.Config) {}

func (sf *Dlt645TCPServerProvider) SetTCPTimeout(t time.Duration) {
	sf.timeout = t
}

// flush flushes pending data in the connection,
// returns io.EOF if connection is closed.
func (sf *Dlt645TCPServerProvider) flush(b []byte) (err error) {
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
func (sf *Dlt645TCPServerProvider) SetConnect(conn net.Conn) error {
	sf.conn = conn
	return nil
}
func (sf *Dlt645TCPServerProvider) SetBeforeCode(enable bool) error {
	sf.BeforeEnable = enable
	return nil
}
