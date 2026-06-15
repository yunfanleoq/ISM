package packages

import (
	"reflect"

	"github.com/beego/beego/v2/core/logs"
	"github.com/mattn/anko/env"
)

func init() {
	env.Packages["logs"] = map[string]reflect.Value{
		"Error":   reflect.ValueOf(logs.Error),
		"Warning": reflect.ValueOf(logs.Warning),
		"Info":    reflect.ValueOf(logs.Info),
		"Trace":   reflect.ValueOf(logs.Trace),
		"Notice":  reflect.ValueOf(logs.Notice),
	}

	env.PackageTypes["logs"] = map[string]reflect.Type{}
}
