package packages

import (
	"encoding/hex"
	"reflect"

	"github.com/mattn/anko/env"
)

func init() {
	env.Packages["encoding/hex"] = map[string]reflect.Value{
		"Decode":         reflect.ValueOf(hex.Decode),
		"DecodeString":   reflect.ValueOf(hex.DecodeString),
		"DecodedLen":     reflect.ValueOf(hex.DecodedLen),
		"Dump":           reflect.ValueOf(hex.Dump),
		"Dumper":         reflect.ValueOf(hex.Dumper),
		"Encode":         reflect.ValueOf(hex.Encode),
		"EncodeToString": reflect.ValueOf(hex.EncodeToString),
		"EncodedLen":     reflect.ValueOf(hex.EncodedLen),
		"NewDecoder":     reflect.ValueOf(hex.NewDecoder),
		"NewEncoder":     reflect.ValueOf(hex.NewEncoder),
	}

	env.PackageTypes["encoding/hex"] = map[string]reflect.Type{}
}
