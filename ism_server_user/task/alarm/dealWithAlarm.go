/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:36
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-08-04 15:53:50
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package alarmTask

import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"mime"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"gopkg.in/gomail.v2"
)

var AlarmConfigRWMutex sync.Mutex
var AlarmConfig map[string][]models.AlarmNotice

func UpdateAlarmNoticeConfig(ProjectUuid string) int {

	AlarmConfigRWMutex.Lock()
	AlarmConfig = make(map[string][]models.AlarmNotice)
	var findAlarmNotice []models.AlarmNotice
	code := models.Db.Model(&models.AlarmNotice{}).Select("alarm_notice_type", "alarm_notice_params", "project_uuid").Where("project_uuid = ?", ProjectUuid).Find(&findAlarmNotice).Error
	if code != nil {
		AlarmConfigRWMutex.Unlock()
		return -2
	}
	AlarmConfig[ProjectUuid] = findAlarmNotice
	AlarmConfigRWMutex.Unlock()
	return 0
}
func convertToGBK(s string) ([]byte, error) {
	// 将字符串转换为rune切片处理Unicode字符
	r := []rune(s)
	// 初始化GBK编码器
	encoder := simplifiedchinese.GBK.NewEncoder()
	// 转换字符串到GBK
	gbkBytes, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(string(r)), encoder))
	if err != nil {
		return nil, err
	}
	return gbkBytes, nil
}
func GetAlarmNoticeConfig() int {

	AlarmConfigRWMutex.Lock()
	AlarmConfig = make(map[string][]models.AlarmNotice)
	var findAlarmNotice []models.AlarmNotice
	code := models.Db.Model(&models.AlarmNotice{}).Select("alarm_notice_type", "alarm_notice_params", "project_uuid").Where("id > 0").Find(&findAlarmNotice).Error
	if code != nil {
		AlarmConfigRWMutex.Unlock()
		return -2
	}
	for _, item := range findAlarmNotice {
		AlarmConfig[item.ProjectUuid] = append(AlarmConfig[item.ProjectUuid], item)
	}
	AlarmConfigRWMutex.Unlock()
	return 0
}

var DeviceAlarmTemp = make(map[string]protocol_common.PushAlarm, protocol_common.AlarmCacheCount)

func SendDingDingMessage(config string, sendAlarm protocol_common.PushAlarm) int {

	type DingTalkItem struct {
		IsEnable bool   `jons:"IsEnable"`
		Webhook  string `jons:"Webhook"`
		Secret   string `jons:"Secret"`
	}
	var DingTalkConfig []DingTalkItem
	err := json.Unmarshal([]byte(config), &DingTalkConfig)
	if err != nil {
		return -1
	}
	SendHtmlContent, res := DingTalkTempleteString(sendAlarm)
	if res != 0 {
		return -1
	}

	for _, talk := range DingTalkConfig {
		if talk.IsEnable && talk.Webhook != "" && talk.Secret != "" {
			d := DingTalk{
				Webhook:  talk.Webhook,
				Secret:   talk.Secret,
				EnableAt: true,
				AtAll:    true,
			}
			d.SendMessage(SendHtmlContent)
		}
	}
	return 0
}
func SendWeChatMessage(config string, sendAlarm protocol_common.PushAlarm) int {

	var weTalkConfig models.AlarmNoticeWeChat
	err := json.Unmarshal([]byte(config), &weTalkConfig)
	if err != nil {
		return -1
	}
	SendHtmlContent, res := DingTalkTempleteString(sendAlarm)
	if res != 0 {
		return -1
	}

	if weTalkConfig.IsEnable {
		d := WeChat{
			Corpid:     weTalkConfig.EnterpriseID,
			Corpsecret: weTalkConfig.Secret,
			Agentid:    weTalkConfig.AgentId,
		}
		go d.SendMessage(SendHtmlContent)
	}

	return 0
}
func EmailTempleteString(sendAlarm protocol_common.PushAlarm) (string, int) {

	tmpl, err := template.ParseFiles("data/tpl/alarm.tpl")

	if err != nil {
		return "", -1
	}
	var buf bytes.Buffer
	err3 := tmpl.Execute(&buf, sendAlarm)
	if err3 != nil {
		fmt.Println(err3)
	}
	return buf.String(), 0
}
func DingTalkTempleteString(sendAlarm protocol_common.PushAlarm) (string, int) {

	tmpl, err := template.ParseFiles("data/tpl/DingTalk.tpl")

	if err != nil {
		return "", -1
	}
	var buf bytes.Buffer
	err3 := tmpl.Execute(&buf, sendAlarm)
	if err3 != nil {
		fmt.Println(err3)
	}
	return buf.String(), 0
}
func SpeekerTempleteString(sendAlarm protocol_common.PushAlarm) (string, int) {

	tmpl, err := template.ParseFiles("data/tpl/Speeker.tpl")

	if err != nil {
		return "", -1
	}
	var buf bytes.Buffer
	err3 := tmpl.Execute(&buf, sendAlarm)
	if err3 != nil {
		fmt.Println(err3)
	}
	return buf.String(), 0
}
func MessageTempleteString(sendAlarm protocol_common.PushAlarm) (string, int) {

	tmpl, err := template.ParseFiles("data/tpl/message.tpl")

	if err != nil {
		return "", -1
	}
	var buf bytes.Buffer
	err3 := tmpl.Execute(&buf, sendAlarm)
	if err3 != nil {
		fmt.Println(err3)
	}
	return buf.String(), 0
}
func SendMail(config string, sendAlarm protocol_common.PushAlarm) int {
	var EmailConfig models.AlarmNoticeMail
	err := json.Unmarshal([]byte(config), &EmailConfig)
	if err != nil {
		return -1
	}
	if !EmailConfig.IsEnable {
		return -2
	}
	SendHtmlContent, res := EmailTempleteString(sendAlarm)
	if res != 0 {
		return -1
	}

	port := EmailConfig.MailServerPort

	MailTo := strings.Split(EmailConfig.MailTo, ";")
	if len(MailTo) == 0 {
		return -1
	}

	for _, mail := range MailTo {
		m := gomail.NewMessage()
		m.SetHeader("From", mime.QEncoding.Encode("UTF-8", EmailConfig.MailSendUserName)+"<"+EmailConfig.MailSendUser+">")

		m.SetHeader("To", mail)
		m.SetHeader("Subject", EmailConfig.MailSendSubject)
		m.SetBody("text/html", SendHtmlContent)
		d := gomail.NewDialer(EmailConfig.MailServerIP, port, EmailConfig.MailSendUser, EmailConfig.MailSendPassword)
		d.TLSConfig = &tls.Config{InsecureSkipVerify: EmailConfig.TLS}
		err1 := d.DialAndSend(m)
		if err1 != nil {
			logs.Error("To:", mail, "##", "Send Email Failed! Err:", err1)
		} else {
			logs.Error("To:", mail, "##", "Send Email Successfully!")
		}
	}

	return 0
}
func SendSms(config string, sendAlarm protocol_common.PushAlarm) int {
	var PhoneConfig models.AlarmNoticeSms
	var isTrue int = 0
	err := json.Unmarshal([]byte(config), &PhoneConfig)
	if err != nil {
		return -1
	}
	if !PhoneConfig.IsEnable {
		return -2
	}
	for _, v := range PhoneConfig.SendAlarmLevel {
		if v == fmt.Sprint(sendAlarm.AlarmLevel) {
			isTrue = 1
			break
		}
	}

	if isTrue == 1 {
		//阿里云
		if PhoneConfig.SmsType == 1 {
			d := AliyunSms{
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
			go d.SendAliyunSms(sendAlarm)
		}
	}

	return 0
}
func SendVoice(config string, sendAlarm protocol_common.PushAlarm) int {
	var PhoneConfig models.AlarmNoticeVoice
	var isTrue int = 0
	err := json.Unmarshal([]byte(config), &PhoneConfig)
	if err != nil {
		return -1
	}
	if !PhoneConfig.IsEnable {
		return -2
	}
	for _, v := range PhoneConfig.SendAlarmLevel {
		if v == fmt.Sprint(sendAlarm.AlarmLevel) {
			isTrue = 1
			break
		}
	}

	if isTrue == 1 {
		//互亿
		if PhoneConfig.VoiceType == 1 {
			d := IhuyiVoice{
				AlarmNoticeVoice: models.AlarmNoticeVoice{
					PhoneNumbers: PhoneConfig.PhoneNumbers,
					IhuyiVoice: models.AlarmNoticeIhuyiVoice{
						AccessKeyId:     PhoneConfig.IhuyiVoice.AccessKeyId,
						AccessKeySecret: PhoneConfig.IhuyiVoice.AccessKeySecret,
						TemplateCode:    PhoneConfig.IhuyiVoice.TemplateCode,
					},
				},
			}
			go d.SendIhuyiVoice(sendAlarm)
		}
	}

	return 0
}
func SendSpeekerCX830E(sendAlarm protocol_common.PushAlarm) int {
	var SpeekerConfig models.AlarmNoticeSpeeker
	SpeekerConf, err := config.NewConfig("ini", "conf/Speeker.conf")
	if err != nil {
		return -1
	}
	SpeekerConfig.IsEnable, err = SpeekerConf.Bool("isEnable")
	if err != nil || !SpeekerConfig.IsEnable {
		return -3
	}
	SpeekerConfig.IpAddress, err = SpeekerConf.String("CX830E::IpAddress")
	if SpeekerConfig.IpAddress == "" || err != nil {
		return -5
	}
	SpeekerConfig.Volume, err = SpeekerConf.Int("CX830E::Volume")
	if err != nil {
		return -5
	}
	conn, err := net.Dial("tcp", SpeekerConfig.IpAddress)
	if err != nil {
		return -4
	}
	defer conn.Close()
	AlarmClearMessage, err1 := convertToGBK(sendAlarm.AlarmClearMessage)
	AlarmMessage, err2 := convertToGBK(sendAlarm.AlarmMessage)
	if err1 != nil || err2 != nil {
		return -3
	}

	sendAlarm.AlarmClearMessage = string(AlarmClearMessage)
	sendAlarm.AlarmMessage = string(AlarmMessage)
	SpeekerContent, res := SpeekerTempleteString(sendAlarm)
	if res != 0 {
		return -1
	}
	SpeekerContent = "#[v" + fmt.Sprintf("%d", SpeekerConfig.Volume) + "]" + SpeekerContent
	conn.Write([]byte(SpeekerContent))
	return 0
}
func SendAlarmNotice(alarm protocol_common.PushAlarm) int {
	AlarmConfigRWMutex.Lock()
	AlarmConfigItem, ok := AlarmConfig[alarm.ProjectUuid]
	SendSpeekerCX830E(alarm)
	if ok {
		for _, config := range AlarmConfigItem {
			if config.AlarmNoticeType == "Mail" {
				SendMail(config.AlarmNoticeParams, alarm)
			} else if config.AlarmNoticeType == "Phone" {
				SendSms(config.AlarmNoticeParams, alarm)
			} else if config.AlarmNoticeType == "dingTalk" {
				SendDingDingMessage(config.AlarmNoticeParams, alarm)
			} else if config.AlarmNoticeType == "weChat" {
				SendWeChatMessage(config.AlarmNoticeParams, alarm)
			} else if config.AlarmNoticeType == "IhiyiVoice" {
				SendVoice(config.AlarmNoticeParams, alarm)
			}
		}
	}
	AlarmConfigRWMutex.Unlock()
	return 0
}

func DealWithAlarm() {
	GetAlarmNoticeConfig()
	for {
		data, code := protocol_common.GAlarmQueue.QueuePull()
		if data == nil {
			time.Sleep(time.Millisecond * 1000)
			continue
		}
		if code != -1 {
			var build strings.Builder
			var updateAlarm models.DevicesAlarmList
			alarm := data.(protocol_common.PushAlarm)
			build.WriteString(alarm.DeviceUuid)
			build.WriteString(alarm.DataUuid)
			key := build.String()
			alarmTemp, isExist := DeviceAlarmTemp[key]

			updateAlarm.AlarmName = alarm.DataName
			updateAlarm.DeviceUuid = alarm.DeviceUuid
			updateAlarm.ProjectUuid = alarm.ProjectUuid
			updateAlarm.DeviceName = alarm.DeviceName
			updateAlarm.DataUuid = alarm.DataUuid
			updateAlarm.ModelDataUuid = alarm.ModelDataUuid
			updateAlarm.HappenTime = alarm.HappenTime
			updateAlarm.AlarmLevel = alarm.AlarmLevel

			updateAlarm.KeepTime = 0
			alarm.Cmd = "RealAlarm"

			var AlarmMessage bytes.Buffer
			t1 := template.New("AlarmMessage")
			tmpl, _ := t1.Parse(alarm.AlarmMessage)
			if tmpl != nil {
				err3 := tmpl.Execute(&AlarmMessage, alarm)
				if err3 != nil {
					updateAlarm.AlarmMessage = alarm.AlarmMessage
				} else {
					updateAlarm.AlarmMessage = AlarmMessage.String()
				}
			} else {
				updateAlarm.AlarmMessage = alarm.AlarmMessage
			}

			var AlarmClearMessage bytes.Buffer
			t2 := template.New("AlarmClearMessage")
			tmpl2, _ := t2.Parse(alarm.AlarmClearMessage)
			if tmpl2 != nil {
				err4 := tmpl2.Execute(&AlarmClearMessage, alarm)
				if err4 != nil {
					updateAlarm.AlarmClearMessage = alarm.AlarmClearMessage
				} else {
					updateAlarm.AlarmClearMessage = AlarmClearMessage.String()
				}
			} else {
				updateAlarm.AlarmClearMessage = alarm.AlarmClearMessage
			}
			alarm.AlarmClearMessage = updateAlarm.AlarmClearMessage
			alarm.AlarmMessage = updateAlarm.AlarmMessage

			if !isExist {
				oldValue, isexit := protocol_common.DeviceRealDataMapByUUID.Load(alarm.DeviceUuid + alarm.DataUuid)
				if alarm.Value == "1" {
					ClearTime, _ := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")
					updateAlarm.ClearTime = ClearTime
					models.Db.Model(&models.DevicesAlarmList{}).Create(&updateAlarm)
					alarm.ID = updateAlarm.ID
					if isexit {
						if oldValue != alarm.Value {
							protocol_common.PushGAlarmQueue.QueuePush(alarm)
							if updateAlarm.DataUuid == "sys.suid.device.status" {
								models.Db.Model(&models.MonitorList{}).Where("uuid = ?", alarm.DeviceUuid).Update("status", 0)
							}
							go SendAlarmNotice(alarm)
						}
					} else {
						protocol_common.PushGAlarmQueue.QueuePush(alarm)
						if updateAlarm.DataUuid == "sys.suid.device.status" {
							models.Db.Model(&models.MonitorList{}).Where("uuid = ?", alarm.DeviceUuid).Update("status", 0)
						}
						go SendAlarmNotice(alarm)
					}
				} else {
					if isexit {
						if oldValue != alarm.Value {
							if updateAlarm.DataUuid == "sys.suid.device.status" {
								protocol_common.PushGAlarmQueue.QueuePush(alarm)
								go SendAlarmNotice(alarm)
								models.Db.Model(&models.MonitorList{}).Where("uuid = ?", alarm.DeviceUuid).Update("status", 1)
							}
						}
					} else {
						if updateAlarm.DataUuid == "sys.suid.device.status" {
							protocol_common.PushGAlarmQueue.QueuePush(alarm)
							go SendAlarmNotice(alarm)
							models.Db.Model(&models.MonitorList{}).Where("uuid = ?", alarm.DeviceUuid).Update("status", 1)
						}
					}
				}
				DeviceAlarmTemp[key] = alarm
			} else {
				if alarmTemp.Value != alarm.Value {
					var status int = 0
					if alarm.Value == "1" {
						ClearTime, _ := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")
						updateAlarm.ClearTime = ClearTime
						models.Db.Model(&models.DevicesAlarmList{}).Create(&updateAlarm)
						alarm.ID = updateAlarm.ID
						status = 0
						protocol_common.PushGAlarmQueue.QueuePush(alarm)
						if alarm.AlarmMessage != "" {
							go SendAlarmNotice(alarm)
						}
					} else {
						updateAlarm.ClearTime = alarm.HappenTime
						updateAlarm.KeepTime = (float64)((alarm.HappenTime.UnixMilli() - alarmTemp.HappenTime.UnixMilli()) / 1000.0)
						status = 1
						models.Db.Model(&models.DevicesAlarmList{}).Where("ID = ? AND device_uuid = ? AND data_uuid = ?", alarmTemp.ID, alarm.DeviceUuid, alarm.DataUuid).Updates(models.DevicesAlarmList{ClearTime: updateAlarm.ClearTime, KeepTime: updateAlarm.KeepTime})
						protocol_common.PushGAlarmQueue.QueuePush(alarm)
						if alarm.AlarmClearMessage != "" {
							go SendAlarmNotice(alarm)
						}
					}
					if updateAlarm.DataUuid == "sys.suid.device.status" {
						models.Db.Model(&models.MonitorList{}).Where("uuid = ?", alarm.DeviceUuid).Update("status", status)
					}

					DeviceAlarmTemp[key] = alarm
				}
			}
		}
		time.Sleep(time.Millisecond * 100)
	}
}
