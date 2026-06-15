package packages

import (
	"reflect"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mattn/anko/env"
)

func init() {
	env.Packages["mqtt"] = map[string]reflect.Value{
		"NewClientOptions": reflect.ValueOf(mqtt.NewClientOptions),
		// "AddBroker":        reflect.ValueOf(mqtt.AddBroker),
		// "SetClientID":              reflect.ValueOf(mqtt.SetClientID),
		// "SetUsername":              reflect.ValueOf(mqtt.SetUsername),
		// "SetPassword":              reflect.ValueOf(mqtt.SetPassword),
		// "SetKeepAlive":             reflect.ValueOf(mqtt.SetKeepAlive),
		// "SetDefaultPublishHandler": reflect.ValueOf(mqtt.SetDefaultPublishHandler),
		"NewClient": reflect.ValueOf(mqtt.NewClient),
		// "Connect":                  reflect.ValueOf(mqtt.Connect),
		// "Unsubscribe":              reflect.ValueOf(mqtt.Unsubscribe),
		// "Subscribe":                reflect.ValueOf(mqtt.Subscribe),
	}
	var message mqtt.Message
	var client mqtt.Client
	var messageHandler mqtt.MessageHandler
	var connectionLostHandler mqtt.ConnectionLostHandler
	var onConnectHandler mqtt.OnConnectHandler
	env.PackageTypes["mqtt"] = map[string]reflect.Type{
		"MessageHandler":        reflect.TypeOf(&messageHandler).Elem(),
		"Client":                reflect.TypeOf(&client).Elem(),
		"Message":               reflect.TypeOf(&message).Elem(),
		"ConnectionLostHandler": reflect.TypeOf(&connectionLostHandler).Elem(),
		"OnConnectHandler":      reflect.TypeOf(&onConnectHandler).Elem(),
	}
}
