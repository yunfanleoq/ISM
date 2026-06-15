/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:36
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-08-04 16:09:59
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package alarmTask

import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

type AliyunSms struct {
	models.AlarmNoticeSms
}

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func (t *AliyunSms) createClient(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}
func (t *AliyunSms) SendAliyunSms(sendAlarm protocol_common.PushAlarm) (int, string) {
	client, _err := t.createClient(tea.String(t.AliYunSms.AccessKeyId), tea.String(t.AliYunSms.AccessKeySecret))
	if _err != nil {
		return -1, ""
	}
	if sendAlarm.DataName == "device.DeviceStatus" {
		sendAlarm.DataName = "设备状态"
	}
	if sendAlarm.Value == "1" {
		if sendAlarm.AlarmMessage == "device.DeviceStatusOffline" {
			sendAlarm.AlarmMessage = "设备离线"
		} else if sendAlarm.AlarmMessage == "VideoManager.VideoOffline" {
			sendAlarm.AlarmMessage = "视频设备离线"
		}
	} else {
		if sendAlarm.AlarmClearMessage == "device.DeviceStatusOnline" {
			sendAlarm.AlarmMessage = "设备在线"
		} else if sendAlarm.AlarmClearMessage == "VideoManager.VideoOnline" {
			sendAlarm.AlarmMessage = "视频设备在线"
		} else {
			sendAlarm.AlarmMessage = sendAlarm.AlarmClearMessage
		}
	}
	smsmap := map[string]interface{}{
		"device_name":   sendAlarm.DeviceName,
		"data_name":     sendAlarm.DataName,
		"alarm_message": sendAlarm.AlarmMessage,
		"happen_time":   sendAlarm.HappenTime.Format("2006-01-02 15:04:05"),
	}
	smsjson, _ := json.Marshal(smsmap)
	SendTo := strings.Split(t.PhoneNumbers, ";")
	if len(SendTo) == 0 {
		return -2, ""
	}
	for _, phone := range SendTo {
		sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
			SignName:      tea.String(t.AliYunSms.SignName),
			TemplateCode:  tea.String(t.AliYunSms.TemplateCode),
			PhoneNumbers:  tea.String(phone),
			TemplateParam: tea.String(string(smsjson)),
		}
		// 复制代码运行请自行打印 API 的返回值
		res, _err := client.SendSms(sendSmsRequest)
		if _err != nil {
			continue
		}
		ResCode := res.Body.Code
		if *ResCode == "OK" {
			continue
		} else {
			fmt.Println(string(smsjson))
			fmt.Println(*res.Body.Message)
			return -3, *res.Body.Message
		}
	}

	return 0, ""
}
func (t *AliyunSms) SendAliyunSmsTest(sms string) (int, string) {
	client, _err := t.createClient(tea.String(t.AliYunSms.AccessKeyId), tea.String(t.AliYunSms.AccessKeySecret))
	if _err != nil {
		return -1, ""
	}
	timeGet := time.Now()
	smsmap := map[string]interface{}{
		"device_name":   "测试设备",
		"data_name":     "测试告警名称",
		"alarm_message": "告警消息",
		"happen_time":   timeGet.Format("2006-01-02 15:04:05"),
	}
	smsjson, _ := json.Marshal(smsmap)
	SendTo := strings.Split(t.PhoneNumbers, ";")
	if len(SendTo) == 0 {
		return -2, ""
	}
	for _, phone := range SendTo {
		sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
			SignName:      tea.String(t.AliYunSms.SignName),
			TemplateCode:  tea.String(t.AliYunSms.TemplateCode),
			PhoneNumbers:  tea.String(phone),
			TemplateParam: tea.String(string(smsjson)),
		}
		// 复制代码运行请自行打印 API 的返回值
		res, _err := client.SendSms(sendSmsRequest)
		if _err != nil {
			return -2, ""
		}
		ResCode := res.Body.Code
		if *ResCode == "OK" {
			continue
		} else {
			return -3, *res.Body.Message
		}
	}

	return 0, ""
}
