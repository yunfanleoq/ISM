package packages

import (
	"reflect"

	"github.com/gorilla/websocket"
	"github.com/mattn/anko/env"
)

func init() {

	env.Packages["websocket"] = map[string]reflect.Value{}

	env.PackageTypes["websocket"] = map[string]reflect.Type{
		"Upgrader": reflect.TypeOf(websocket.Upgrader{}),
	}
}
