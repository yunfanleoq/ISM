package packages

import (
	"encoding/base64"
	"reflect"

	"github.com/mattn/anko/env"
)

func init() {
	env.Packages["encoding/base64"] = map[string]reflect.Value{
		"DecodeString":   reflect.ValueOf(base64.StdEncoding.DecodeString),
		"EncodeToString": reflect.ValueOf(base64.StdEncoding.EncodeToString),
	}

	env.PackageTypes["encoding/base64"] = map[string]reflect.Type{}
}
