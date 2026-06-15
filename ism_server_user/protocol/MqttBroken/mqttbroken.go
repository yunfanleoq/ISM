// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 mochi-mqtt, mochi-co
// SPDX-FileContributor: mochi-co

package ismmqttbroken

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"time"

	appconf "github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/config"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/packets"
)

func MqttBrokenServer() {

	IsStart, errPort := appconf.Bool("EnableMqttBreoken")
	if errPort != nil {
		IsStart = false
	}
	if !IsStart {
		logs.Info("禁止启动 MQTT Broken,如果启用请在conf/app.conf里面添加EnableMqttBreoken=true")
		return
	}

	logs.Info("启动MQTT Broken.......")

	level := new(slog.LevelVar)

	configBytes, err1 := os.ReadFile("conf/mqtt_broken_config.json")
	if err1 != nil {
		logs.Error("MQTT Broker启动失败", err1)
	}

	options, err1 := config.FromBytes(configBytes)
	if err1 != nil {
		logs.Error("MQTT Broker启动失败", err1)
	}

	options.InlineClient = true
	server := mqtt.New(options)

	server.Log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
	level.Set(slog.LevelError)

	_ = server.AddHook(new(auth.AllowHook), nil)

	// Add custom hook (ExampleHook) to the server
	err := server.AddHook(new(ISMHook), &ISMHookOptions{
		Server: server,
	})
	if err != nil {
		logs.Error("MQTT Broker启动失败", err)
	}

	// err = server.AddHook(new(badger.Hook), &badger.Options{
	// 	Path: "data/mqttbroken",
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// Start the server

	serverErr := server.Serve()
	if serverErr != nil {
		logs.Error("MQTT Broker启动失败", serverErr)
	}
}

// Options contains configuration settings for the hook.
type ISMHookOptions struct {
	Server *mqtt.Server
}

type ISMHook struct {
	mqtt.HookBase
	config *ISMHookOptions
}

func (h *ISMHook) ID() string {
	return "events-example"
}

func (h *ISMHook) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnect,
		mqtt.OnDisconnect,
	}, []byte{b})
}

func (h *ISMHook) Init(config any) error {
	h.Log.Info("initialised")
	if _, ok := config.(*ISMHookOptions); !ok && config != nil {
		return mqtt.ErrInvalidConfigType
	}

	h.config = config.(*ISMHookOptions)
	if h.config.Server == nil {
		return mqtt.ErrInvalidConfigType
	}
	return nil
}

func (h *ISMHook) OnConnect(cl *mqtt.Client, pk packets.Packet) error {
	type DeviceStatus struct {
		ClientID    string
		Status      string
		HappendTime time.Time
	}
	var pushdata DeviceStatus

	pushdata.ClientID = cl.ID
	pushdata.HappendTime = time.Now()
	pushdata.Status = "Online"

	data, _ := json.Marshal(pushdata)

	err := h.config.Server.Publish("Broker/clients/status", data, false, 0)
	if err != nil {
		h.Log.Error("hook.publish", "error", err)
	}
	return nil
}

func (h *ISMHook) OnDisconnect(cl *mqtt.Client, err error, expire bool) {

	type DeviceStatus struct {
		ClientID    string
		Status      string
		HappendTime time.Time
	}
	var pushdata DeviceStatus

	pushdata.ClientID = cl.ID
	pushdata.HappendTime = time.Now()
	pushdata.Status = "Offline"

	data, _ := json.Marshal(pushdata)

	err1 := h.config.Server.Publish("Broker/clients/status", data, false, 0)
	if err1 != nil {
		h.Log.Error("hook.publish", "error", err)
	}
}
