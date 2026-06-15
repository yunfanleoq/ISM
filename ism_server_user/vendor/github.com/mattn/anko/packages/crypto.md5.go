package packages

import (
	"crypto/md5"
	"reflect"

	"github.com/mattn/anko/env"
)

func init() {
	env.Packages["crypto/md5"] = map[string]reflect.Value{
		"New": reflect.ValueOf(md5.New),
		"Sum": reflect.ValueOf(md5.Sum),
	}

	env.PackageTypes["crypto/md5"] = map[string]reflect.Type{}
}
