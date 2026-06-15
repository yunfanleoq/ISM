/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:36
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-06-15 10:33:48
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package ISMNetwork

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/beego/beego/v2/core/logs"
)

type SystemNetworkInterfaceIP struct {
	IP   string
	Mask string
}
type SystemNetworkInterface struct {
	Name       string
	MacAddress string
	GateWay    string
	IPv4       []SystemNetworkInterfaceIP
	IPv6       []SystemNetworkInterfaceIP
}

func InitNetwork() {

	if (runtime.GOARCH != "arm64") && (runtime.GOARCH != "arm") {
		return
	} else {
		var systemNetwork []SystemNetworkInterface
		filename := "conf/NetWorker.json"
		netcontent, err := os.ReadFile(filename)
		if err == nil {
			logs.Info("开始初始化网络...")
			err = json.Unmarshal(netcontent, &systemNetwork)
			if err == nil {
				for _, v := range systemNetwork {
					for num, item := range v.IPv4 {

						var interName string
						if num != 0 {
							interName = v.Name + ":" + fmt.Sprintf("%d", (num-1))
						} else {
							interName = v.Name
						}

						logs.Info("设置%s的IP地址:%s,Mask:%s\n", interName, item.IP, item.Mask)
						// 构建ifconfig命令
						cmd := exec.Command("ifconfig", interName, item.IP, "netmask", item.Mask)
						var out bytes.Buffer
						cmd.Stdout = &out
						cmd.Stderr = &out
						err := cmd.Run()
						if err != nil {
							logs.Error(fmt.Sprintf("设置IP地址 error: %s, output: %s", err, out.String()))
						}
					}

					if v.MacAddress != "" {
						var out bytes.Buffer
						logs.Info("设置%s的MAC地址:%s\n", v.Name, v.MacAddress)
						cmd := exec.Command("ifconfig", v.Name, "down")
						err := cmd.Run()
						if err != nil {
							logs.Error(fmt.Sprintf("停止网口 error: %s, output: %s", err, out.String()))
						}

						// 构建ifconfig命令
						cmd = exec.Command("ifconfig", v.Name, "hw", "ether", v.MacAddress)

						cmd.Stdout = &out
						cmd.Stderr = &out
						err = cmd.Run()
						if err != nil {
							logs.Error(fmt.Sprintf("设置MAC地址错误 error: %s, output: %s", err, out.String()))
						}
						cmd = exec.Command("ifconfig", v.Name, "up")
						cmd.Stdout = &out
						cmd.Stderr = &out
						err = cmd.Run()
						if err != nil {
							logs.Error(fmt.Sprintf("启动网卡错误 error: %s, output: %s", err, out.String()))
						}
					}

					if v.GateWay != "" {
						logs.Info("设置%s的网关地址:%s\n", v.Name, v.GateWay)
						// 构建ifconfig命令
						cmd := exec.Command("ip", "route", "add", "default", "via", v.GateWay, "dev", v.Name)
						var out bytes.Buffer
						cmd.Stdout = &out
						cmd.Stderr = &out
						err := cmd.Run()
						if err != nil {
							logs.Error(fmt.Sprintf("设置网关错误 error: %s, output: %s", err, out.String()))
						}
					}
				}
			}
		}
	}
}
