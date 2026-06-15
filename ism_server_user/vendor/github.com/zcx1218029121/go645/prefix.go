package go645

import (
	"bytes"
	"io"
)

var _ PrefixHandler = (*DefaultPrefix)(nil)

type PrefixHandler interface {
	EncodePrefix(buffer *bytes.Buffer, enable bool) error
	DecodePrefix(reader io.Reader) ([]byte, error)
}

type DefaultPrefix struct {
}

func (d DefaultPrefix) EncodePrefix(buffer *bytes.Buffer, enable bool) error {
	if enable {
		var beforeCode = make([]byte, 4)
		beforeCode[0] = 0xFE
		beforeCode[1] = 0xFE
		beforeCode[2] = 0xFE
		beforeCode[3] = 0xFE
		buffer.Write(beforeCode)
	}
	return nil
}

func (d DefaultPrefix) DecodePrefix(reader io.Reader) ([]byte, error) {
	return nil, nil
}
