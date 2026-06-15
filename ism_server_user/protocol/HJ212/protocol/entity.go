package HJ212Protocol

// Entity Body实体
type Entity interface {
	Encode() ([]byte, error)
	Decode(data string) error
}
