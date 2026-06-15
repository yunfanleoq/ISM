package packages

import (
	"crypto/sha256"
	"reflect"

	"github.com/mattn/anko/env"
)

func init() {
	env.Packages["crypto/sha256"] = map[string]reflect.Value{
		"New":    reflect.ValueOf(sha256.New),
		"New224": reflect.ValueOf(sha256.New224),
		"Sum224": reflect.ValueOf(sha256.Sum224),
		"Sum256": reflect.ValueOf(sha256.Sum256),
	}

	env.PackageTypes["crypto/sha256"] = map[string]reflect.Type{}
}
