package packages

import (
	"html/template"
	"reflect"

	"github.com/mattn/anko/env"
)

func init() {
	env.Packages["html/template"] = map[string]reflect.Value{
		"New": reflect.ValueOf(template.New),
	}

	env.PackageTypes["html/template"] = map[string]reflect.Type{}
}
