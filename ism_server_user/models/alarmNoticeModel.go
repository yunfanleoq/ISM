/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-06-16 10:11:34
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package models

import (
	"ISMServer/utils/errmsg"
	"errors"

	"encoding/json"

	createUUID "github.com/go-basic/uuid"
	"gorm.io/gorm"
)

// 设备告警表
type AlarmNotice struct {
	gorm.Model

	Uuid              string `gorm:"index;type:varchar(250);not null" json:"Uuid" validate:"required" label:"UUID"`
	AlarmNoticeType   string `gorm:"index;type:varchar(250);not null" json:"AlarmNoticeType" validate:"required,min=4,max=250" label:"通知类型"`
	AlarmNoticeParams string `gorm:"type:text;not null" json:"AlarmNoticeParams" validate:"required" label:"通知类型的参数"`
	ProjectUuid       string `gorm:"index;type:varchar(250);not null" json:"ProjectUuid" validate:"required,min=4,max=250" label:"项目uuid"`
}
type AlarmNoticeMail struct {
	IsEnable         bool   `json:"IsEnable" label:"是否使能"`
	MailServerIP     string `json:"MailServerIP" label:"邮件中心IP"`
	MailServerPort   int    `json:"MailServerPort" label:"邮件中心端口"`
	MailSendUser     string `json:"MailSendUser" label:"邮件发送账号"`
	MailSendPassword string `json:"MailSendPassword" label:"邮件发送密码"`
	MailSendAddress  string `json:"MailSendAddress" label:"邮件发送地址"`
	TLS              bool   `json:"TLS" label:"邮件TLS"`
	MailTo           string `json:"MailTo" label:"邮件发送地址"`
	MailSendUserName string `json:"MailSendUserName" label:"邮件发件人"`
	MailSendSubject  string `json:"MailSendSubject" label:"邮件主题"`
}
type AlarmNoticeWeChat struct {
	IsEnable     bool   `json:"IsEnable" label:"是否使能"`
	EnterpriseID string `json:"EnterpriseID" label:"企业ID"`
	AgentId      int    `json:"AgentId" label:"应用ID"`
	Secret       string `json:"Secret" label:"应用秘钥"`
}
type AlarmNoticeSpeeker struct {
	IsEnable  bool   `json:"IsEnable" label:"是否使能"`
	IpAddress string `json:"IpAddress" label:"IP 地址"`
	Volume    int    `json:"Volume" label:"音量"`
}
type AlarmNoticeAliyunSms struct {
	AccessKeyId     string `json:"AccessKeyId" label:"AccessKey ID"`
	AccessKeySecret string `json:"AccessKeySecret" label:"AccessKey Secret"`
	SignName        string `json:"SignName" label:"短信签名名称"`
	TemplateCode    string `json:"TemplateCode" label:"短信模板CODE"`
}
type AlarmNoticeIhuyiVoice struct {
	AccessKeyId     string `json:"AccessKeyId" label:"AccessKey ID"`
	AccessKeySecret string `json:"AccessKeySecret" label:"AccessKey Secret"`
	TemplateCode    string `json:"TemplateCode" label:"语音模板CODE"`
}
type AlarmNoticeSms struct {
	IsEnable          bool                 `json:"IsEnable" label:"是否使能"`
	SmsType           int                  `json:"SmsType" label:"短信平台"`
	AliYunSms         AlarmNoticeAliyunSms `json:"AliYunSms" label:"短信平台"`
	EveryDaySendCount int                  `json:"EveryDaySendCount" label:"每天最多发送多少"`
	SendAlarmLevel    []string             `json:"SendAlarmLevel" label:"发送的告警等级"`
	PhoneNumbers      string               `json:"PhoneNumbers" label:"手机号码"`
}
type AlarmNoticeVoice struct {
	IsEnable          bool                  `json:"IsEnable" label:"是否使能"`
	VoiceType         int                   `json:"VoiceType" label:"语音平台"`
	IhuyiVoice        AlarmNoticeIhuyiVoice `json:"IhuyiVoice" label:"短信平台"`
	EveryDaySendCount int                   `json:"EveryDaySendCount" label:"每天最多发送多少"`
	SendAlarmLevel    []string              `json:"SendAlarmLevel" label:"发送的告警等级"`
	PhoneNumbers      string                `json:"PhoneNumbers" label:"手机号码"`
}

/*
*
告警通知获取
*/
func AlarmNoticeGetAll(ProjectUuid string, AlarmNoticeType string) (int, AlarmNotice) {

	var getAlarmNotice AlarmNotice

	code := Db.Model(&AlarmNotice{}).Where("alarm_notice_type = ? and project_uuid = ?", AlarmNoticeType, ProjectUuid).First(&getAlarmNotice)
	if errors.Is(code.Error, gorm.ErrRecordNotFound) {
		getAlarmNotice.ProjectUuid = ProjectUuid
		getAlarmNotice.Uuid = createUUID.New()
		getAlarmNotice.AlarmNoticeType = AlarmNoticeType
		if AlarmNoticeType == "Mail" {
			var mail AlarmNoticeMail
			mail.IsEnable = true
			mail.MailServerIP = ""
			mail.MailServerPort = 25
			mail.MailSendUser = ""
			mail.MailSendPassword = ""
			mail.MailSendAddress = ""
			mail.TLS = true
			jsonByte, _ := json.Marshal(mail)
			getAlarmNotice.AlarmNoticeParams = string(jsonByte)
		} else if AlarmNoticeType == "Phone" {
			var phone AlarmNoticeSms
			phone.IsEnable = true
			phone.SmsType = 1 //阿里云

			phone.PhoneNumbers = ""
			phone.EveryDaySendCount = 10

			phone.AliYunSms.AccessKeyId = ""
			phone.AliYunSms.AccessKeySecret = ""
			phone.AliYunSms.SignName = ""
			phone.AliYunSms.TemplateCode = ""

			jsonByte, _ := json.Marshal(phone)
			getAlarmNotice.AlarmNoticeParams = string(jsonByte)
		} else if AlarmNoticeType == "dingTalk" {
			getAlarmNotice.AlarmNoticeParams = "[{\"IsEnable\":true,\"Webhook\":\"\",\"Secret\":\"\"}]"
		} else if AlarmNoticeType == "weChat" {
			var mail AlarmNoticeWeChat
			mail.IsEnable = true
			mail.EnterpriseID = ""
			mail.AgentId = 0
			mail.Secret = ""
			jsonByte, _ := json.Marshal(mail)
			getAlarmNotice.AlarmNoticeParams = string(jsonByte)
		} else if AlarmNoticeType == "IhiyiVoice" {
			var voice AlarmNoticeVoice
			voice.IsEnable = true
			voice.VoiceType = 1
			voice.IhuyiVoice.AccessKeyId = ""
			voice.IhuyiVoice.AccessKeySecret = ""
			voice.IhuyiVoice.TemplateCode = ""
			voice.EveryDaySendCount = 10
			voice.PhoneNumbers = ""
			jsonByte, _ := json.Marshal(voice)
			getAlarmNotice.AlarmNoticeParams = string(jsonByte)
		}
		Db.Model(&AlarmNotice{}).Create(&getAlarmNotice)
	}
	return errmsg.SUCCSECODE, getAlarmNotice
}

/*
*
告警通知更新
*/
func AlarmNoticeUpdate(AlarmNoticeType string, ProjectUuid string, params string) int {
	code := Db.Model(&AlarmNotice{}).Where("alarm_notice_type = ? and project_uuid = ?", AlarmNoticeType, ProjectUuid).Update("alarm_notice_params", params).Error
	if code != nil {
		return errmsg.ERROR_DATABASE
	}
	return errmsg.SUCCSECODE
}
