package cjt188

import (
	"bytes"
	"fmt"
	"sync"
)

var (
	// check implements Client interface.
	_ Client = (*client)(nil)
)

type client struct {
	ClientProvider
	mu sync.Mutex
}

func (c *client) Read(address Address, itemCode int32) (*ReadData, bool, error) {
	// c.ClientProvider.SendCurrentAddress(address.strValue)
	// c.ClientProvider.SendCurrentCode(address.strValue)
	resp, err := c.ClientProvider.SendAndRead(ReadRequest(address, itemCode))
	if err != nil {
		return nil, false, err
	}
	decode, err := Decode(bytes.NewBuffer(resp))
	if err != nil {
		return nil, false, err
	}
	recvCode := Bytes2Int(decode.Data.(*ReadData).dataType)
	if address.strValue != decode.Address.strValue {
		return nil, false, fmt.Errorf("address :%s,recvaddress:%s", address.strValue, decode.Address.strValue)
	}
	if uint16(itemCode) != recvCode {
		return nil, false, fmt.Errorf("sendCode :%d,recvCode:%d", itemCode, recvCode)
	}
	return decode.Data.(*ReadData), decode.Control.IsState(HasNext), err
}

// Broadcast 设备广播
func (c *client) Broadcast(p InformationElement, control Control) error {
	var err error
	bf := bytes.NewBuffer(make([]byte, 0))
	err = p.Encode(bf)
	if err != nil {
		return err
	}
	return c.Send(NewProtocol(NewAddress(BroadcastAddress, LittleEndian), p, &control))
}

func (c *client) ReadWithBlock(address Address, data ReadRequestData) (*Protocol, error) {
	resp, err := c.ClientProvider.SendAndRead(ReadRequestWithBlock(address, data))
	if err != nil {
		return nil, err
	}
	return Decode(bytes.NewBuffer(resp))
}

// Option custom option
type Option func(c *client)

func NewClient(p ClientProvider, opts ...Option) Client {
	c := &client{ClientProvider: p}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
