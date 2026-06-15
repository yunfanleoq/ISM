package packages

import (
	"reflect"

	"github.com/mattn/anko/env"
	"github.com/xuri/excelize/v2"
)

func init() {
	env.Packages["excelize"] = map[string]reflect.Value{
		"OpenFile": reflect.ValueOf(excelize.OpenFile),
		"NewFile":  reflect.ValueOf(excelize.NewFile),
	}
}
