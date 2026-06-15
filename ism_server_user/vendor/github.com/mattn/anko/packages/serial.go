package packages

import (
	"reflect"

	"github.com/goburrow/serial"
	"github.com/mattn/anko/env"
)

func init() {
	env.Packages["serial"] = map[string]reflect.Value{
		"Open": reflect.ValueOf(serial.Open),
	}
	env.PackageTypes["serial"] = map[string]reflect.Type{
		"Config": reflect.TypeOf(serial.Config{}),
	}
}
