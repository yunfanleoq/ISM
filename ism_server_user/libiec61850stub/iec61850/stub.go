// Package iec61850 is a compile-time stub for macOS / platforms without
// prebuilt libiec61850. API matches github.com/jifanchn/go-libiec61850/iec61850
// so ISM protocol code can build locally.
//
// On Kylin/Linux x86_64 production builds, remove the go.mod replace directive
// so the real CGO vendor library (linux64) is linked.
package iec61850

import (
	"errors"
	"time"
)

type FunctionalConstraint int
type MMSType int
type IedClientError int

type Option func(client *IedClient)

type IedClient struct {
	withoutTimestamps bool
}

type GoMmsValue struct {
	Type  MMSType
	Value interface{}
}

const (
	IEC61850_FC_ST FunctionalConstraint = iota
	IEC61850_FC_MX
	IEC61850_FC_SP
	IEC61850_FC_SV
	IEC61850_FC_CF
	IEC61850_FC_DC
	IEC61850_FC_SG
	IEC61850_FC_SE
	IEC61850_FC_SR
	IEC61850_FC_OR
	IEC61850_FC_BL
	IEC61850_FC_EX
	IEC61850_FC_CO
	IEC61850_FC_US
	IEC61850_FC_MS
	IEC61850_FC_RP
	IEC61850_FC_BR
	IEC61850_FC_LG
	IEC61850_FC_GO

	IEC61850_FC_ALL  FunctionalConstraint = 99
	IEC61850_FC_NONE FunctionalConstraint = -1
)

const (
	IED_ERROR_OK IedClientError = 0
)

var errStubUnavailable = errors.New("iec61850 stub: native libiec61850 is not available on this platform")

func NewIedClient(options ...Option) *IedClient {
	client := &IedClient{}
	for _, op := range options {
		if op != nil {
			op(client)
		}
	}
	return client
}

func ConnectTimeout(timeout time.Duration) Option {
	return func(c *IedClient) {}
}

func RequestTimeout(timeout time.Duration) Option {
	return func(c *IedClient) {}
}

func WithoutTimestamps(flag bool) Option {
	return func(c *IedClient) {
		c.withoutTimestamps = flag
	}
}

func (client *IedClient) Connect(hostname string, tcpPort int) error {
	return errStubUnavailable
}

func (client *IedClient) Close() {}

func (client *IedClient) ReadBoolean(objectRef string, constraint FunctionalConstraint) (bool, error) {
	return false, errStubUnavailable
}

func (client *IedClient) ReadVisibleString(objectRef string, constraint FunctionalConstraint) (string, error) {
	return "", errStubUnavailable
}

func (client *IedClient) ReadInt32(objectRef string, constraint FunctionalConstraint) (int32, error) {
	return 0, errStubUnavailable
}

func (client *IedClient) ReadUnsigned32(objectRef string, constraint FunctionalConstraint) (uint32, error) {
	return 0, errStubUnavailable
}

func (client *IedClient) ReadInt64(objectRef string, constraint FunctionalConstraint) (int64, error) {
	return 0, errStubUnavailable
}

func (client *IedClient) ReadFloat(objectRef string, constraint FunctionalConstraint) (float64, error) {
	return 0, errStubUnavailable
}

func (client *IedClient) WirteBoolean(objectRef string, constraint FunctionalConstraint, setvalue bool) IedClientError {
	return 1
}

func (client *IedClient) WirteVisibleString(objectRef string, constraint FunctionalConstraint, setvalue string) IedClientError {
	return 1
}

func (client *IedClient) WirteInt32(objectRef string, constraint FunctionalConstraint, setvalue int32) IedClientError {
	return 1
}

func (client *IedClient) WirteUnsigned32(objectRef string, constraint FunctionalConstraint, setvalue uint32) IedClientError {
	return 1
}

func (client *IedClient) WirteFloat(objectRef string, constraint FunctionalConstraint, setvalue float32) IedClientError {
	return 1
}
