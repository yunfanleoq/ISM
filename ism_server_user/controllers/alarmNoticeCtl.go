/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-06-16 09:39:21
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	alarmTask "ISMServer/task/alarm"
	"ISMServer/utils/errmsg"
	"crypto/tls"
	"encoding/json"
	"mime"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
	"gopkg.in/gomail.v2"
)

type AlarmNoticeController struct {
	beego.Controller
}

func (c *AlarmNoticeController) GetAlarmNoticeParams() {

	var list models.AlarmNotice

	var getParams map[string]interface{}
	var code = 0

	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {

		data := c.Ctx.Input.RequestBody
		//json数据封装到对象中
		err := json.Unmarshal(data, &getParams)
		if err != nil {
			code = errmsg.NOTJSON
		} else {
			if getParams["type"] != nil {
				code, list = models.AlarmNoticeGetAll(ProjectUuid, getParams["type"].(string))
			} else {
				code = -3
			}

		}

	} else {
		code = -1
	}
	result := map[string]interface{}{
		"code": code,
		"list": list,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
func (c *AlarmNoticeController) UpdateAlarmNoticeParams() {

	var list models.AlarmNotice

	var getParams map[string]interface{}
	var code = 0

	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {

		data := c.Ctx.Input.RequestBody
		//json数据封装到对象中
		err := json.Unmarshal(data, &getParams)
		if err != nil {
			code = errmsg.NOTJSON
		} else {
			code = models.AlarmNoticeUpdate(getParams["type"].(string), ProjectUuid, getParams["params"].(string))
			alarmTask.UpdateAlarmNoticeConfig(ProjectUuid)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "更新告警参数"+getParams["params"].(string), errmsg.JournalLevelInfo, c.Ctx.Input)
		}

	} else {
		code = -1
	}
	result := map[string]interface{}{
		"code": code,
		"list": list,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
func SendMailTest(config string) int {
	var EmailConfig models.AlarmNoticeMail
	err := json.Unmarshal([]byte(config), &EmailConfig)
	if err != nil {
		return -1
	}
	port := EmailConfig.MailServerPort
	m := gomail.NewMessage()
	m.SetHeader("From", mime.QEncoding.Encode("UTF-8", EmailConfig.MailSendUserName)+"<"+EmailConfig.MailSendUser+">")
	MailTo := strings.Split(EmailConfig.MailTo, ";")
	if len(MailTo) == 0 {
		return -1
	}

	for _, mail := range MailTo {
		m.SetHeader("To", mail)
		m.SetHeader("Subject", EmailConfig.MailSendSubject)
		m.SetBody("text/html", "Test Email")
		d := gomail.NewDialer(EmailConfig.MailServerIP, port, EmailConfig.MailSendUser, EmailConfig.MailSendPassword)

		d.TLSConfig = &tls.Config{InsecureSkipVerify: EmailConfig.TLS}
		err1 := d.DialAndSend(m)
		if err1 != nil {
			return -3
		}
	}

	return 0
}
func (c *AlarmNoticeController) TestSend() {

	var getParams map[string]interface{}
	var code = 0
	var msg string = ""
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {

		data := c.Ctx.Input.RequestBody
		//json数据封装到对象中
		err := json.Unmarshal(data, &getParams)
		if err != nil {
			code = errmsg.NOTJSON
		} else {
			if getParams["type"].(string) == "Mail" {
				code = SendMailTest(getParams["params"].(string))
			} else if getParams["type"].(string) == "dingTalk" {
				code, msg = SendDingDingMessageTest(getParams["params"].(string), int(getParams["index"].(float64)))
			} else if getParams["type"].(string) == "weChat" {
				code, msg = SendWeChatMessage(getParams["params"].(string))
			} else if getParams["type"].(string) == "Phone" {
				code, msg = SendSmsTest(getParams["params"].(string))
			} else if getParams["type"].(string) == "IhiyiVoice" {
				code, msg = SendVoiceTest(getParams["params"].(string))
			} else {
				code = -3
			}

		}

	} else {
		code = -1
	}
	result := map[string]interface{}{
		"code": code,
		"msg":  msg,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
func SendDingDingMessageTest(config string, index int) (int, string) {

	type DingTalkItem struct {
		IsEnable bool   `jons:"IsEnable"`
		Webhook  string `jons:"Webhook"`
		Secret   string `jons:"Secret"`
	}
	var DingTalkConfig []DingTalkItem
	err := json.Unmarshal([]byte(config), &DingTalkConfig)
	if err != nil {
		return -1, ""
	}
	talk := DingTalkConfig[index]
	if talk.IsEnable && talk.Webhook != "" && talk.Secret != "" {
		d := alarmTask.DingTalk{
			Webhook:  talk.Webhook,
			Secret:   talk.Secret,
			EnableAt: true,
			AtAll:    true,
		}
		err, result := d.SendMessage("Test Message")
		if err != 0 {
			return -2, result
		}
	} else {
		return -3, ""
	}

	return 0, ""
}
func SendWeChatMessage(config string) (int, string) {
	var sendAlarm protocol_common.PushAlarm
	var weTalkConfig models.AlarmNoticeWeChat
	err := json.Unmarshal([]byte(config), &weTalkConfig)
	if err != nil {
		return -1, ""
	}
	SendHtmlContent, res := alarmTask.DingTalkTempleteString(sendAlarm)
	if res != 0 {
		return -1, ""
	}

	if weTalkConfig.IsEnable {
		d := alarmTask.WeChat{
			Corpid:     weTalkConfig.EnterpriseID,
			Corpsecret: weTalkConfig.Secret,
			Agentid:    weTalkConfig.AgentId,
		}
		return d.SendMessage(SendHtmlContent)
	}

	return -2, ""
}

func SendSmsTest(config string) (int, string) {
	var PhoneConfig models.AlarmNoticeSms
	err := json.Unmarshal([]byte(config), &PhoneConfig)
	if err != nil {
		return -1, ""
	}

	//阿里云
	if PhoneConfig.SmsType == 1 {
		d := alarmTask.AliyunSms{
			AlarmNoticeSms: models.AlarmNoticeSms{
				PhoneNumbers: PhoneConfig.PhoneNumbers,
				AliYunSms: models.AlarmNoticeAliyunSms{
					AccessKeyId:     PhoneConfig.AliYunSms.AccessKeyId,
					AccessKeySecret: PhoneConfig.AliYunSms.AccessKeySecret,
					SignName:        PhoneConfig.AliYunSms.SignName,
					TemplateCode:    PhoneConfig.AliYunSms.TemplateCode,
				},
			},
		}
		return d.SendAliyunSmsTest("22345")
	}

	return -2, ""
}
func SendVoiceTest(config string) (int, string) {
	var PhoneConfig models.AlarmNoticeVoice
	err := json.Unmarshal([]byte(config), &PhoneConfig)
	if err != nil {
		return -1, ""
	}

	//阿里云
	if PhoneConfig.VoiceType == 1 {
		d := alarmTask.IhuyiVoice{
			AlarmNoticeVoice: models.AlarmNoticeVoice{
				PhoneNumbers: PhoneConfig.PhoneNumbers,
				IhuyiVoice: models.AlarmNoticeIhuyiVoice{
					AccessKeyId:     PhoneConfig.IhuyiVoice.AccessKeyId,
					AccessKeySecret: PhoneConfig.IhuyiVoice.AccessKeySecret,
					TemplateCode:    PhoneConfig.IhuyiVoice.TemplateCode,
				},
			},
		}
		return d.SendIhuyiTest()
	}

	return -2, ""
}
