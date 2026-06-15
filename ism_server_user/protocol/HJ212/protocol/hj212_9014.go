package HJ212Protocol

import "strings"

type HJ212_9014 struct {
}

// TODO 数据编码
func (entity *HJ212_9014) Encode() ([]byte, error) {
	var builder strings.Builder

	builder.WriteString("CP=&&&&")

	return []byte(builder.String()), nil
}

// 数据解码
func (entity *HJ212_9014) Decode(data string) error {

	return nil
}
