package packages

import (
	"encoding/binary"
	"reflect"

	"github.com/mattn/anko/env"
)

func init() {
	env.Packages["encoding/binary"] = map[string]reflect.Value{
		"Read":   reflect.ValueOf(binary.Read),
		"Write": reflect.ValueOf(binary.Write),
		"AppendUvarint": reflect.ValueOf(binary.AppendUvarint),
		"AppendVarint": reflect.ValueOf(binary.AppendVarint),
		"PutUvarint": reflect.ValueOf(binary.PutUvarint),
		"PutVarint": reflect.ValueOf(binary.PutVarint),
		"ReadUvarint": reflect.ValueOf(binary.ReadUvarint),
		"ReadVarint": reflect.ValueOf(binary.ReadVarint),
		"Size": reflect.ValueOf(binary.Size),
		"Uvarint": reflect.ValueOf(binary.Uvarint),
		"Varint": reflect.ValueOf(binary.Varint),
		"BigEndian": reflect.ValueOf(binary.BigEndian),
		"LittleEndian": reflect.ValueOf(binary.LittleEndian),
	}
	
	env.PackageTypes["encoding/binary"] = map[string]reflect.Type{
		
	}
}