package cjt188

import (
	"net"
	"time"

	"github.com/goburrow/serial"
)

type ClientProvider interface {
	// Connect try to connect the remote server
	Connect() error
	// IsConnected returns a bool signifying whether
	// the client is connected or not.
	IsConnected() bool
	SetBeforeCode(enable bool) error
	LogMode(enable bool)
	SetConnect(conn net.Conn) error
	// Close disconnect the remote server
	Close() error
	setSerialConfig(config serial.Config)
	setPrefixHandler(handler PrefixHandler)
	// setTCPTimeout set tcp connect & read timeout
	SetTCPTimeout(t time.Duration)
	SetLogProvider(p LogProvider)
	SendAndRead(*Protocol) (aduResponse []byte, err error)
	SendRawFrameAndRead(aduRequest []byte) (aduResponse []byte, err error)
	SendRawFrame(aduRequest []byte) (err error)
	ReadRawFrame() (aduResponse []byte, err error)
	Send(*Protocol) (err error)
}

// LogProvider  log message levels only Debug and Error
type LogProvider interface {
	Error(prefix string, format string, v ...interface{})
	Debug(prefix string, format string, v ...interface{})
}
