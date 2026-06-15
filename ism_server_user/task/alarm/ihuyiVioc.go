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
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const strUrl = "https://api.ihuyi.com/vm/Submit.json"

type IhuyiVoice struct {
	models.AlarmNoticeVoice
}

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func (t *IhuyiVoice) createClient(accessKeyId string, accessKeySecret string, templateid string) *url.Values {
	v := url.Values{}
	v.Set("account", accessKeyId)      //APIID（用户中心【云语音】-【语音通知】-【产品总览】查看）
	v.Set("password", accessKeySecret) //1、APIKEY（用户中心【云语音】-【语音通知】-【产品总览】查看）2、动态密码（生成动态密码方式请看该文档末尾的说明）
	v.Set("templateid", templateid)    //语音模板ID（使用模板变量方式发送时必填）调试阶段可使用系统默认模
	return &v
}
func (t *IhuyiVoice) SendIhuyiVoice(sendAlarm protocol_common.PushAlarm) (int, string) {
	fmt.Println(t.IhuyiVoice)
	v := t.createClient(t.IhuyiVoice.AccessKeyId, t.IhuyiVoice.AccessKeySecret, t.IhuyiVoice.TemplateCode)

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
	SendTo := strings.Split(t.PhoneNumbers, ";")
	if len(SendTo) == 0 {
		return -2, ""
	}
	v.Set("time", fmt.Sprintf("%10d", time.Now().Unix()))
	v.Set("content", sendAlarm.DeviceName+"|"+sendAlarm.DataName+"|"+sendAlarm.AlarmMessage+"|"+sendAlarm.HappenTime.Format("2006-01-02 15:04:05"))
	for _, phone := range SendTo {
		v.Set("mobile", phone)
		body := strings.NewReader(v.Encode()) //把form数据编码
		client := &http.Client{}
		req, _ := http.NewRequest("POST", strUrl, body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req) //发送
		if err != nil {
			fmt.Println(err)
		}
		res, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(res))
		resp.Body.Close() //一定要关闭resp.Body
	}

	return 0, ""
}
func (t *IhuyiVoice) SendIhuyiTest() (int, string) {
	v := t.createClient(t.IhuyiVoice.AccessKeyId, t.IhuyiVoice.AccessKeySecret, t.IhuyiVoice.TemplateCode)
	timeGet := time.Now()
	SendTo := strings.Split(t.PhoneNumbers, ";")
	if len(SendTo) == 0 {
		return -2, ""
	}
	v.Set("time", fmt.Sprintf("%10d", time.Now().Unix()))
	v.Set("content", "测试设备|测试告警名称|告警消息|"+timeGet.Format("2006-01-02 15:04:05"))
	for _, phone := range SendTo {
		v.Set("mobile", phone)

		body := strings.NewReader(v.Encode()) //把form数据编码
		client := &http.Client{}
		req, _ := http.NewRequest("POST", strUrl, body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req) //发送
		if err != nil {
			fmt.Println(err)
		}
		resp.Body.Close() //一定要关闭resp.Body
	}

	return 0, ""
}
