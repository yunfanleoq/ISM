/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-09-05 15:34:10
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package DataInterface

import (
	"ISMServer/models"
	ISMScript "ISMServer/task/ISMScript/func"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var GMqttInterfaceChan chan bool
var mqttPushWg sync.WaitGroup

type MqttPushCtl struct {
	waitGroup  *sync.WaitGroup
	PushData   MqttPushInterfaceData
	MqttClient MQTT.Client
	ConnStatus bool
}

func (c *MqttPushCtl) InitMqttPush() {

}

// 重连后重新订阅
func (cr *MqttPushCtl) onMqttConnect(c MQTT.Client) {
	cr.ConnStatus = true
	logs.Info("MQTT 接口%s,连接成功", cr.PushData.InterfaceName)
}

// 连接丢失处理
func (cr *MqttPushCtl) onConnectionLost(c MQTT.Client, err error) {
	cr.ConnStatus = false
	logs.Error("MQTT 接口%s,连接断开,原因:%s", cr.PushData.InterfaceName, err)
}

func (cr *MqttPushCtl) pubhander(client MQTT.Client, msg MQTT.Message) {

}

func (c *MqttPushCtl) MqttClientConnect() {

	var brokerHost string = ""

	for {
		//检测协程是否主动退出
		select {
		case <-GMqttInterfaceChan:
			return
		default:
			time.Sleep(1 * time.Millisecond) // 降低空转频率
		}
		BrokerHost := c.PushData.InterfaceHost
		BrokerPort := c.PushData.InterfacePort
		UserName := c.PushData.InterfaceUser
		PassWord := c.PushData.InterfacePassword
		ClientID := c.PushData.InterfaceCid

		opts := MQTT.NewClientOptions()

		brokerHost = "mqtt://" + BrokerHost + ":" + BrokerPort

		opts = opts.AddBroker(brokerHost)

		opts.SetClientID(ClientID)
		if UserName != "" && PassWord != "" {
			opts.SetUsername(UserName)
			opts.SetPassword(PassWord)
		}
		opts.SetCleanSession(false)
		opts.SetAutoReconnect(true)
		opts.SetOnConnectHandler(c.onMqttConnect)
		opts.SetConnectionLostHandler(c.onConnectionLost)
		opts.SetKeepAlive(60 * 2 * time.Second)
		opts.SetDefaultPublishHandler(c.pubhander)
		c.MqttClient = MQTT.NewClient(opts)
		if token := c.MqttClient.Connect(); token.Wait() && token.Error() == nil {
			logs.Info("Pusb Mqtt Connect %s Mqtt Broker Sucess\n", brokerHost)
			break
		} else {
			logs.Error("Pusb Mqtt Connect %s Mqtt Broker Error %v\n", brokerHost, token.Error())
		}
		time.Sleep(time.Duration(c.PushData.InterfaceDataInterval) * time.Millisecond)
	}

}
func (c *MqttPushCtl) MqttPushPthread() {

	c.MqttClientConnect()
	for {
		//检测协程是否主动退出
		select {
		case <-GMqttInterfaceChan:
			c.MqttClient.Disconnect(0)
			c.waitGroup.Done()
			return
		default:
			time.Sleep(1 * time.Millisecond) // 降低空转频率
		}
		if !c.MqttClient.IsConnected() {
			time.Sleep(time.Duration(c.PushData.InterfaceDataInterval) * time.Millisecond)
			continue
		}
		for _, data := range c.PushData.InterfaceDataContent {

			// 定义正则表达式，匹配 {{}} 中间的字符串
			re := regexp.MustCompile(`\{\{(.*?)\}\}`)

			// 使用函数动态替换匹配的内容
			replacedText := re.ReplaceAllStringFunc(data, func(match string) string {
				// 去掉 {{ 和 }}，只保留中间的内容
				inner := strings.TrimPrefix(match, "{{")
				inner = strings.TrimSuffix(inner, "}}")

				value := ISMScript.GetDeviceRealData(inner)
				if value == nil {
					value = 0
				}
				_, ok := value.(string)
				if ok {
					return value.(string)
				} else {
					tempValue, ok := value.(float64)
					if ok {
						truncated := math.Trunc(tempValue)
						diff := tempValue - truncated
						istrue := math.Abs(diff) > 1e-9
						if !istrue {
							return fmt.Sprintf("%d", int64(tempValue))
						} else {
							return fmt.Sprintf("%.5f", float64(tempValue))
						}
					} else if tempValue1, ok1 := value.(int); ok1 {
						return fmt.Sprintf("%d", int64(tempValue1))
					} else {
						return fmt.Sprintf("%d", int64(tempValue1))
					}
				}
			})
			for _, tipc := range c.PushData.InterfaceSubject {
				c.MqttClient.Publish(tipc, 2, true, replacedText)
			}
		}
		time.Sleep(time.Duration(c.PushData.InterfaceDataInterval) * time.Millisecond)
	}
}

type MqttPushInterfaceData struct {
	InterfaceName         string
	InterfaceHost         string
	InterfacePort         string
	InterfaceCid          string
	InterfaceUser         string
	InterfacePassword     string
	InterfaceSubject      []string
	InterfaceDataContent  []string
	InterfaceDataInterval int
}

func isMqttPushChanClose() bool {
	select {
	case _, received := <-GMqttInterfaceChan:
		return !received
	default:
	}
	return false
}
func MqttPushInterfaceCloseChan() {
	isOpen := isMqttPushChanClose()
	if !isOpen && GMqttInterfaceChan != nil {
		close(GMqttInterfaceChan)
	}
}
func MqttPushInterfaceStart() {

	var is_starting = 0
	type mqttPushData struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Id       string `json:"id"`
		User     string `json:"user"`
		Password string `json:"password"`
		Interval int    `json:"interval"`
		Subject  string `json:"subject"`
	}
	for {
		if is_starting == 1 {
			mqttPushWg.Wait()
		}
		MqttPushInterfaceCloseChan()
		GMqttInterfaceChan = make(chan bool)
		var getData []models.SystemDataInterface
		err := models.Db.Model(&models.SystemDataInterface{}).Where("interface_type = 3").Select("*").Find(&getData).Error
		if err != nil {
			time.Sleep(10 * time.Second)
			continue
		}
		for _, v := range getData {
			if v.InterfaceStatus == 0 {
				continue
			}
			var SUrlInterfaceData MqttPushInterfaceData
			var edata mqttPushData
			SUrlInterfaceData.InterfaceName = v.InterfaceName
			err := json.Unmarshal([]byte(v.InterfaceContent), &edata)
			if err != nil {
				continue
			}
			SUrlInterfaceData.InterfaceHost = edata.Host
			SUrlInterfaceData.InterfacePort = edata.Port
			SUrlInterfaceData.InterfaceCid = edata.Id
			SUrlInterfaceData.InterfaceUser = edata.User
			SUrlInterfaceData.InterfacePassword = edata.Password

			ds := strings.Split(edata.Subject, ",")
			SUrlInterfaceData.InterfaceSubject = append(SUrlInterfaceData.InterfaceSubject, ds...)
			if len(SUrlInterfaceData.InterfaceSubject) == 0 {
				continue
			}
			SUrlInterfaceData.InterfaceDataInterval = edata.Interval
			d := strings.Split(v.InterfaceDataUuid, ",")
			for _, uuid := range d {
				var getDataContent models.SystemDataTemplete
				err := models.Db.Model(&models.SystemDataTemplete{}).Where("templete_uuid = ?", uuid).Select("*").Find(&getDataContent).Error
				if err == nil {
					SUrlInterfaceData.InterfaceDataContent = append(SUrlInterfaceData.InterfaceDataContent, getDataContent.TempleteContent)
				}
			}
			if len(SUrlInterfaceData.InterfaceDataContent) == 0 {
				continue
			}
			durl := &MqttPushCtl{waitGroup: &mqttPushWg, PushData: SUrlInterfaceData}
			go durl.MqttPushPthread()
			mqttPushWg.Add(1)
			is_starting = 1
		}
		if is_starting == 0 {
			time.Sleep(10 * time.Second)
			continue
		}
	}

}
