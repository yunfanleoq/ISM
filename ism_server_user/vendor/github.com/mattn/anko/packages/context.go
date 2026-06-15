package packages

import (
	"context"
	"reflect"

	"github.com/mattn/anko/env"
)

func init() {
	env.Packages["context"] = map[string]reflect.Value{
		"Background": reflect.ValueOf(context.Background),
	}
	var Context context.Context
	env.PackageTypes["context"] = map[string]reflect.Type{
		"Context": reflect.TypeOf(&Context).Elem(),
	}
}
