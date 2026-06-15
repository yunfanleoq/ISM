package HJ212Protocol

import (
	"fmt"
	"strings"

	hjerrors "ISMServer/protocol/HJ212/errors"
)

type Message struct {
	Header Header
	Body   Entity
}

// Encode 消息编码
func (message *Message) Encode() ([]byte, error) {
	// 消息头编码
	header, err := message.Header.Encode()
	if err != nil {
		return nil, err
	}

	// 消息体编码
	body, err := message.Body.Encode()
	if err != nil {
		return nil, err
	}

	// CRC校验
	crc := CRCCheckout(string(header) + string(body))
	crcStr := fmt.Sprintf("%04X", crc)

	// 初始化builder
	builder := strings.Builder{}

	// 写入报文头
	packetHeader := fmt.Sprintf("##%04d", len(header)+len(body))
	builder.WriteString(packetHeader)

	// 写入报文体
	builder.WriteString(string(header))
	builder.WriteString(string(body))

	// 写入CRC校验码
	builder.WriteString(crcStr)

	// 写入报文结束标识符
	builder.WriteString("\r\n")

	return []byte(builder.String()), nil
}

// Decode 消息解码
func (message *Message) Decode(data []byte) error {
	dataStr := string(data)

	// 包头解码
	headerIndex := strings.Index(dataStr, "CP")
	if err := message.Header.Decode(dataStr[6:headerIndex]); err != nil {
		return err
	}

	// 正文解码
	bodyStr := dataStr[headerIndex:]
	entity, err := message.decode(message.Header.CN, bodyStr)
	if err != nil {

	}
	message.Body = entity

	return nil
}

func (message *Message) decode(typ uint16, data string) (Entity, error) {
	// 获取消息body对应的实体解码类
	entity, ok := entityMapper[typ]
	if !ok {
		return nil, hjerrors.ErrEntityNotFound
	}

	if err := entity.Decode(data); err != nil {
		return nil, err
	}

	return entity, nil
}
