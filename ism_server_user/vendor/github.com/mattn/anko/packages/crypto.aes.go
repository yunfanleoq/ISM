package packages

import (
	"crypto/aes"
	"reflect"

	"github.com/mattn/anko/env"
)

func init() {
	env.Packages["crypto/aes"] = map[string]reflect.Value{
		"NewCipher": reflect.ValueOf(aes.NewCipher),
	}

	env.PackageTypes["crypto/aes"] = map[string]reflect.Type{}
}
