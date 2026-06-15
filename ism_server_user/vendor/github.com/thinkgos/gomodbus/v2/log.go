package modbus

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"log"
)

// 内部调试实现.
type logger struct {
	provider LogProvider
	// has log output enabled,
	// 1: enable
	// 0: disable
	has uint32
}

// newLogger new logger with prefix.
func newLogger(prefix string) logger {
	//创建日志文件
	_, err := os.Stat("logs/modbus")

	if os.IsNotExist(err) {
		os.Mkdir("logs/modbus", os.ModePerm)
	}
	fileName := fmt.Sprintf("logs/modbus/modbus_log_%s_%s.log", prefix, time.Now().Format("2006-01-02"))
	logFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0744)
	if err != nil {
		panic(err)
	}
	return logger{
		provider: defaultLogger{log.New(logFile, prefix, log.LstdFlags)},
		has:      0,
	}
}

// LogMode set enable or disable log output when you has set logger.
func (sf *logger) LogMode(enable bool) {
	if enable {
		atomic.StoreUint32(&sf.has, 1)
	} else {
		atomic.StoreUint32(&sf.has, 0)
	}
}

// setLogProvider overwrite log provider.
func (sf *logger) SetLogProvider(p LogProvider) {
	if p != nil {
		sf.provider = p
	}
}

// Error Log ERROR level message.
func (sf logger) Error(prefix string, format string, v ...interface{}) {
	if atomic.LoadUint32(&sf.has) == 1 {
		sf.provider.Error(prefix, format, v...)
	}
}

// Debug Log DEBUG level message.
func (sf logger) Debug(prefix string, format string, v ...interface{}) {
	if atomic.LoadUint32(&sf.has) == 1 {
		sf.provider.Debug(prefix, format, v...)
	}
}

// default log.
type defaultLogger struct {
	*log.Logger
}

// check implement LogProvider interface.
var _ LogProvider = (*defaultLogger)(nil)

// Error Log ERROR level message.
func (sf defaultLogger) Error(prefix string, format string, v ...interface{}) {
	sf.Printf("[E]: "+format, v...)
}

// Debug Log DEBUG level message.
func (sf defaultLogger) Debug(prefix string, format string, v ...interface{}) {
	sf.Printf("[D]: "+format, v...)
}
