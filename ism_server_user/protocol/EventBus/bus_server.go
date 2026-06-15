/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 09:00:01
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package EventBusServer

import (
	"time"

	evbus "github.com/asaskevich/EventBus"
)

func EventBusStart() {
	serverBus := evbus.NewServer(":2020", "/_server_bus_b", evbus.New())
	serverBus.Start()

	for {
		serverBus.EventBus().Publish("topic", 10)
		time.Sleep(time.Second * 1)
	}

}
