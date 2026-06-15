/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:56:19
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

import (
	"ISMServer/utils/errmsg"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beevik/ntp"
)

type SystemParamsController struct {
	beego.Controller
}

// Windows 时区到 IANA 时区的映射表
var windowsToIANA = map[string]string{
	"UTC":                        "Etc/UTC",
	"China Standard Time":        "Asia/Shanghai",
	"Tokyo Standard Time":        "Asia/Tokyo",
	"India Standard Time":        "Asia/Kolkata",
	"Arabian Standard Time":      "Asia/Dubai",
	"Singapore Standard Time":    "Asia/Singapore",
	"SE Asia Standard Time":      "Asia/Bangkok",
	"GMT Standard Time":          "Europe/London",
	"Romance Standard Time":      "Europe/Paris",
	"W. Europe Standard Time":    "Europe/Berlin",
	"Russian Standard Time":      "Europe/Moscow",
	"Eastern Standard Time":      "America/New_York",
	"Central Standard Time":      "America/Chicago",
	"Mountain Standard Time":     "America/Denver",
	"Pacific Standard Time":      "America/Los_Angeles",
	"Venezuela Standard Time":    "America/Caracas",
	"Argentina Standard Time":    "America/Argentina/Buenos_Aires",
	"New Zealand Standard Time":  "Pacific/Auckland",
	"Fiji Standard Time":         "Pacific/Fiji",
	"AUS Eastern Standard Time":  "Australia/Sydney",
	"UTC-02":                     "Atlantic/South_Georgia",
	"Azores Standard Time":       "Atlantic/Azores",
	"Egypt Standard Time":        "Africa/Cairo",
	"South Africa Standard Time": "Africa/Johannesburg",
	"Korea Standard Time":        "Asia/Seoul",
	"Hong Kong Standard Time":    "Asia/Hong_Kong",
	"Philippine Standard Time":   "Asia/Manila",
	"Iran Standard Time":         "Asia/Tehran",
	"Arabic Standard Time":       "Asia/Baghdad",
	"Arab Standard Time":         "Asia/Riyadh",
	"Turkey Standard Time":       "Europe/Istanbul",
	"Pacific SA Standard Time":   "America/Santiago",
	"Hawaiian Standard Time":     "Pacific/Honolulu",
	"West Pacific Standard Time": "Pacific/Port_Moresby",
	"Tonga Standard Time":        "Pacific/Tongatapu",
}

// 获取当前系统时区
func getSystemTimeZone() (string, error) {
	var cmd *exec.Cmd

	// 根据操作系统选择命令
	switch runtime.GOOS {
	case "windows":
		// Windows 使用 tzutil 工具
		cmd = exec.Command("tzutil", "/g")
	case "linux":
		// Linux 使用 timedatectl 工具
		cmd = exec.Command("timedatectl")
	default:
		return "", fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}

	// 执行命令并获取输出
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("获取时区失败: %v", err)
	}

	// 解析输出
	output := strings.TrimSpace(out.String())
	if runtime.GOOS == "linux" {
		// 从 timedatectl 输出中提取时区
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "Time zone:") {
				parts := strings.Fields(line)
				if len(parts) >= 3 {
					return parts[2], nil
				}
			}
		}
		return "", fmt.Errorf("未找到时区信息")
	}
	ianaTimeZone := output
	// 如果是 Windows 系统，转换为 IANA 时区
	if runtime.GOOS == "windows" {
		ianaTimeZone, err = convertWindowsToIANA(output)
		if err != nil {
			return "", err
		}
	}
	// Windows 直接返回 tzutil 的输出
	return ianaTimeZone, nil
}

// 将 Windows 时区转换为 IANA 时区
func convertWindowsToIANA(windowsTimeZone string) (string, error) {
	ianaTimeZone, exists := windowsToIANA[windowsTimeZone]
	if !exists {
		return "", fmt.Errorf("无法找到对应的 IANA 时区: %s", windowsTimeZone)
	}
	return ianaTimeZone, nil
}

// 设置系统时区
func setSystemTimeZone(ianaTimeZone string) error {
	// IANA 时区与 Windows 时区的映射
	timeZoneMap := map[string]string{
		"Etc/UTC":                        "UTC",
		"Asia/Shanghai":                  "China Standard Time",
		"Asia/Tokyo":                     "Tokyo Standard Time",
		"Asia/Kolkata":                   "India Standard Time",
		"Asia/Dubai":                     "Arabian Standard Time",
		"Asia/Singapore":                 "Singapore Standard Time",
		"Asia/Bangkok":                   "SE Asia Standard Time",
		"Europe/London":                  "GMT Standard Time",
		"Europe/Paris":                   "Romance Standard Time",
		"Europe/Berlin":                  "W. Europe Standard Time",
		"Europe/Moscow":                  "Russian Standard Time",
		"America/New_York":               "Eastern Standard Time",
		"America/Chicago":                "Central Standard Time",
		"America/Denver":                 "Mountain Standard Time",
		"America/Los_Angeles":            "Pacific Standard Time",
		"America/Caracas":                "Venezuela Standard Time",
		"America/Argentina/Buenos_Aires": "Argentina Standard Time",
		"Pacific/Auckland":               "New Zealand Standard Time",
		"Pacific/Fiji":                   "Fiji Standard Time",
		"Australia/Sydney":               "AUS Eastern Standard Time",
		"Atlantic/South_Georgia":         "UTC-02",
		"Atlantic/Azores":                "Azores Standard Time",
		"Africa/Cairo":                   "Egypt Standard Time",
		"Africa/Johannesburg":            "South Africa Standard Time",
		"Asia/Seoul":                     "Korea Standard Time",
		"Asia/Hong_Kong":                 "Hong Kong Standard Time",
		"Asia/Manila":                    "Philippine Standard Time",
		"Asia/Jakarta":                   "SE Asia Standard Time",
		"Asia/Tehran":                    "Iran Standard Time",
		"Asia/Baghdad":                   "Arabic Standard Time",
		"Asia/Riyadh":                    "Arab Standard Time",
		"Europe/Istanbul":                "Turkey Standard Time",
		"Europe/Madrid":                  "Romance Standard Time",
		"Europe/Rome":                    "W. Europe Standard Time",
		"Europe/Amsterdam":               "W. Europe Standard Time",
		"Europe/Brussels":                "Romance Standard Time",
		"Europe/Zurich":                  "W. Europe Standard Time",
		"America/Toronto":                "Eastern Standard Time",
		"America/Vancouver":              "Pacific Standard Time",
		"America/Mexico_City":            "Central Standard Time (Mexico)",
		"America/Sao_Paulo":              "E. South America Standard Time",
		"America/Bogota":                 "SA Pacific Standard Time",
		"America/Lima":                   "SA Pacific Standard Time",
		"America/Santiago":               "Pacific SA Standard Time",
		"Pacific/Honolulu":               "Hawaiian Standard Time",
		"Pacific/Guam":                   "West Pacific Standard Time",
		"Pacific/Port_Moresby":           "West Pacific Standard Time",
		"Pacific/Tongatapu":              "Tonga Standard Time",
		"Antarctica/South_Pole":          "New Zealand Standard Time",
		"Antarctica/McMurdo":             "New Zealand Standard Time",
	}

	// 根据操作系统选择时区名称
	var timeZone string
	if runtime.GOOS == "windows" {
		timeZone = timeZoneMap[ianaTimeZone]
		if timeZone == "" {
			return fmt.Errorf("不支持的 IANA 时区: %s", ianaTimeZone)
		}
	} else {
		timeZone = ianaTimeZone
	}

	// 设置时区
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("tzutil", "/s", timeZone)
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("Windows 时区设置失败: %v", err)
		}
	case "linux":
		cmd := exec.Command("timedatectl", "set-timezone", timeZone)
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("Linux 时区设置失败: %v", err)
		}
	default:
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}
	return nil
}
func (c *SystemParamsController) TestNtpServer() {
	type ParamsData struct {
		NTPServer string `json:"NTPServer"`
		NTPPort   int    `json:"NTPPort"`
	}
	var pdata ParamsData
	var code int = -1
	data := c.Ctx.Input.RequestBody

	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到user对象中
		err := json.Unmarshal(data, &pdata)
		if err != nil {
			code = -3
		} else {
			_, err := ntp.QueryWithOptions(pdata.NTPServer, ntp.QueryOptions{
				Timeout: 30 * time.Second,
				Port:    pdata.NTPPort,
			})
			if err != nil {
				code = -4
			} else {
				code = 0
			}
		}
	} else {
		code = -2
	}
	result := map[string]interface{}{
		"code": code,
	}
	c.Data["json"] = result
	c.ServeJSON()
}
func hangSyncSystemTime(timeStr string) error {
	correctTime, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return fmt.Errorf("时间格式错误: %v", err)
	}
	switch runtime.GOOS {
	case "windows":
		// 设置日期（YYYY-MM-DD格式）
		dateCmd := exec.Command("cmd", "/C", "date", correctTime.Format("2006-01-02"))
		if err := dateCmd.Run(); err != nil {
			return fmt.Errorf("时间设置失败: %v (需要管理员权限)", err)
		}

		// 设置时间（HH:MM:SS格式）
		timeCmd := exec.Command("cmd", "/C", "time", correctTime.Format("15:04:05"))
		if err := timeCmd.Run(); err != nil {
			return fmt.Errorf("时间设置失败: %v (需要管理员权限)", err)
		}

		// 同步时区（北京时间）
		tzCmd := exec.Command("tzutil", "/s", "China Standard Time")
		tzCmd.Run() // 时区设置不阻断主流程
	case "linux":
		cmd := exec.Command("date", "-s", correctTime.Format("2006-01-02 15:04:05"))
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("时间设置失败: %v (需要管理员权限)", err)
		}
	default:
		return fmt.Errorf("unsupported OS")
	}
	return nil
}
func syncSystemTime(address string, port int) error {
	response, err := ntp.QueryWithOptions(address, ntp.QueryOptions{
		Timeout: 10 * time.Second,
		Port:    port,
	})

	if err != nil {
		return err
	}

	correctTime := time.Now().Add(response.ClockOffset)

	switch runtime.GOOS {
	case "windows":
		// 设置日期（YYYY-MM-DD格式）
		dateCmd := exec.Command("cmd", "/C", "date", correctTime.Format("2006-01-02"))
		if err := dateCmd.Run(); err != nil {
			return fmt.Errorf("时间设置失败: %v (需要管理员权限)", err)
		}

		// 设置时间（HH:MM:SS格式）
		timeCmd := exec.Command("cmd", "/C", "time", correctTime.Format("15:04:05"))
		if err := timeCmd.Run(); err != nil {
			return fmt.Errorf("时间设置失败: %v (需要管理员权限)", err)
		}

		// 同步时区（北京时间）
		tzCmd := exec.Command("tzutil", "/s", "China Standard Time")
		tzCmd.Run() // 时区设置不阻断主流程
	case "linux":
		cmd := exec.Command("date", "-s", correctTime.Format("2006-01-02 15:04:05"))
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("时间设置失败: %v (需要管理员权限)", err)
		}
		// 将系统时间写入硬件时钟
		cmd = exec.Command("hwclock", "--systohc")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("写入硬件时钟失败: %v (需要管理员权限)", err)
		}
	default:
		return fmt.Errorf("unsupported OS")
	}
	return nil
}
func (c *SystemParamsController) GetSystemWebData() {

	var presult = make(map[string]interface{})
	var err error

	presult["httpport"], err = config.Int("httpport")
	if err != nil {
		presult["httpport"] = 8081
	}

	presult["enablehttps"], err = config.Bool("enablehttps")
	if err != nil {
		presult["enablehttps"] = false
	}

	presult["httpsport"], err = config.Int("httpsport")
	if err != nil {
		presult["httpsport"] = 443
	}
	presult["wsport"], err = config.Int("wsport")
	if err != nil {
		presult["wsport"] = 443
	}

	presult["httpskeyfile"], err = config.String("httpskeyfile")
	if err != nil {
		presult["httpskeyfile"] = ""
	}
	presult["httpscertfile"], err = config.String("httpscertfile")
	if err != nil {
		presult["httpscertfile"] = ""
	}
	now := time.Now()
	formatted := now.Format("2006-01-02 15:04:05")
	presult["systemtime"] = formatted

	timeConf, timeerr := config.NewConfig("ini", "conf/systimeconfig.conf")
	getSystemTimeZone, err := getSystemTimeZone()
	if err != nil {
		getSystemTimeZone = "Asia/Shanghai"
	}
	if timeerr != nil {
		presult["TimeZone"] = getSystemTimeZone
		presult["CheckType"] = 0
		presult["NTPServer"] = "time.windows.com"
		presult["NTPPort"] = 123
		presult["NTPCheckTime"] = 1440
	} else {
		presult["TimeZone"] = getSystemTimeZone
		presult["CheckType"], err = timeConf.String("CheckType")
		if err != nil {
			presult["CheckType"] = 0
		}
		presult["NTPServer"], err = timeConf.String("NTPServer")
		if err != nil {
			presult["NTPServer"] = "time.windows.com"
		}
		presult["NTPPort"], err = timeConf.String("NTPPort")
		if err != nil {
			presult["NTPPort"] = 123
		}
		presult["NTPCheckTime"], err = timeConf.String("NTPCheckTime")
		if err != nil {
			presult["NTPCheckTime"] = 1440
		}
	}
	result := map[string]interface{}{
		"result": presult,
	}
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *SystemParamsController) SaveSystemWebData() {
	type ParamsData struct {
		HTTPPort    int    `json:"HTTPPort"`
		HTTPSEnable bool   `json:"HTTPSEnable"`
		HTTPsPort   int    `json:"HTTPsPort"`
		WSPort      int    `json:"WSPort"`
		HttpsCert   string `json:"HttpsCert"`
		HttpsKey    string `json:"HttpsKey"`
	}
	var pdata ParamsData
	var code int = -1
	data := c.Ctx.Input.RequestBody

	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到user对象中
		err := json.Unmarshal(data, &pdata)
		if err != nil {
			code = -3
		} else {
			if pdata.HTTPPort > 1 && pdata.HTTPPort <= 65534 {
				config.Set("httpport", fmt.Sprintf("%d", pdata.HTTPPort))
			}
			if pdata.HTTPsPort > 1 && pdata.HTTPsPort <= 65534 {
				config.Set("httpsport", fmt.Sprintf("%d", pdata.HTTPsPort))
			}
			if pdata.WSPort > 1 && pdata.WSPort <= 65534 {
				config.Set("wsport", fmt.Sprintf("%d", pdata.WSPort))
			}
			if !pdata.HTTPSEnable || pdata.HTTPSEnable {
				config.Set("enablehttps", fmt.Sprintf("%t", pdata.HTTPSEnable))
			}
			config.Set("httpscertfile", pdata.HttpsCert)
			config.Set("httpskeyfile", pdata.HttpsKey)
			config.SaveConfigFile("conf/app.conf")
			code = 0
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "修改了系统参数 "+fmt.Sprintf("%d", pdata.HTTPPort)+" "+fmt.Sprintf("%d", pdata.HTTPsPort), errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -2
	}
	result := map[string]interface{}{
		"code":   code,
		"result": "",
	}
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *SystemParamsController) GetSystemMqttData() {

	var presult = make(map[string]interface{})
	var err error
	var code int = -1
	iniconf, err := config.NewConfig("ini", "conf/mqtt.conf")
	if err != nil {
		code = -2
	} else {
		presult["isEnable"], err = iniconf.Bool("isEnable")
		if err != nil {
			presult["isEnable"] = false
		}
		presult["TLS"], err = iniconf.Bool("TLS")
		if err != nil {
			presult["TLS"] = false
		}
		presult["mqttCloudPlat"], err = iniconf.Int("mqttCloudPlat")
		if err != nil {
			presult["mqttCloudPlat"] = 1
		}
		presult["BrokerHost"], err = iniconf.String("MQTT::BrokerHost")
		if err != nil {
			presult["BrokerHost"] = "127.0.0.1"
		}
		presult["BrokerPort"], err = iniconf.Int("MQTT::BrokerPort")
		if err != nil {
			presult["BrokerPort"] = 1883
		}
		presult["UserName"], err = iniconf.String("MQTT::UserName")
		if err != nil {
			presult["UserName"] = ""
		}
		presult["PassWord"], err = iniconf.String("MQTT::PassWord")
		if err != nil {
			presult["PassWord"] = ""
		}
		presult["ClientID"], err = iniconf.String("MQTT::ClientID")
		if err != nil {
			presult["ClientID"] = ""
		}
		presult["SubscribeTopic"], err = iniconf.String("MQTT::SubscribeTopic")
		if err != nil {
			presult["SubscribeTopic"] = ""
		}
		presult["PublishTopic"], err = iniconf.String("MQTT::PublishTopic")
		if err != nil {
			presult["PublishTopic"] = ""
		}
		presult["CertPath"], err = iniconf.String("MQTT::certPath")
		if err != nil {
			presult["CertPath"] = ""
		}
	}

	result := map[string]interface{}{
		"code":   code,
		"result": presult,
	}
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *SystemParamsController) SaveSystemMqttData() {
	type ParamsData struct {
		IsEnable       bool   `json:"isEnable"`
		BrokerHost     string `json:"BrokerHost"`
		BrokerPort     int    `json:"BrokerPort"`
		UserName       string `json:"UserName"`
		PassWord       string `json:"PassWord"`
		ClientID       string `json:"ClientID"`
		TLS            bool   `json:"TLS"`
		CertPath       string `json:"certPath"`
		SubscribeTopic string `json:"SubscribeTopic"`
		PublishTopic   string `json:"PublishTopic"`
	}
	var pdata ParamsData
	var code int = -1
	data := c.Ctx.Input.RequestBody

	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到user对象中
		err := json.Unmarshal(data, &pdata)
		if err != nil {
			code = -3
		} else {
			iniconf, err := config.NewConfig("ini", "conf/mqtt.conf")
			if err != nil {
				code = -2
			}

			if !pdata.IsEnable || pdata.IsEnable {
				iniconf.Set("isEnable", fmt.Sprintf("%t", pdata.IsEnable))
			}
			iniconf.Set("mqtt::BrokerHost", pdata.BrokerHost)
			iniconf.Set("mqtt::BrokerPort", fmt.Sprintf("%d", pdata.BrokerPort))
			iniconf.Set("mqtt::UserName", pdata.UserName)
			iniconf.Set("mqtt::PassWord", pdata.PassWord)
			iniconf.Set("mqtt::ClientID", pdata.ClientID)
			if !pdata.TLS || pdata.TLS {
				iniconf.Set("mqtt::TLS", fmt.Sprintf("%t", pdata.TLS))
			}
			iniconf.Set("mqtt::certPath", pdata.CertPath)
			iniconf.Set("mqtt::SubscribeTopic", pdata.SubscribeTopic)
			iniconf.Set("mqtt::PublishTopic", pdata.PublishTopic)

			iniconf.SaveConfigFile("conf/mqtt.conf")
			code = 0
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "修改了MQTT客户端参数 "+pdata.BrokerHost+" "+fmt.Sprintf("%d", pdata.BrokerPort), errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -2
	}
	result := map[string]interface{}{
		"code":   code,
		"result": "",
	}
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *SystemParamsController) GetSystemModbusData() {

	var presult = make(map[string]interface{})
	var err error

	presult["modbusServerPort"], err = config.Int("modbusServerPort")
	if err != nil {
		presult["modbusServerPort"] = 3000
	}
	presult["iec104calldelaytime"], err = config.Int("iec104calldelaytime")
	if err != nil {
		presult["iec104calldelaytime"] = 3000
	}

	result := map[string]interface{}{
		"result": presult,
	}
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *SystemParamsController) SaveSystemModbusData() {
	type ParamsData struct {
		ModbusServerPort    int `json:"ModbusServerPort"`
		Iec104calldelaytime int `json:"iec104calldelaytime"`
	}
	var pdata ParamsData
	var code int = -1
	data := c.Ctx.Input.RequestBody

	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到user对象中
		err := json.Unmarshal(data, &pdata)
		if err != nil {
			code = -3
		} else {
			if pdata.ModbusServerPort > 1 && pdata.ModbusServerPort <= 65534 {
				config.Set("modbusServerPort", fmt.Sprintf("%d", pdata.ModbusServerPort))
			}
			config.Set("iec104calldelaytime", fmt.Sprintf("%d", pdata.Iec104calldelaytime))
			config.SaveConfigFile("conf/app.conf")
			code = 0
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "修改了Modbus 远程采集端口 "+fmt.Sprintf("%d", pdata.ModbusServerPort), errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -2
	}
	result := map[string]interface{}{
		"code":   code,
		"result": "",
	}
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *SystemParamsController) GetSystemHistoryConfig() {

	var presult = make(map[string]interface{})
	var err error

	HistoryConf, err := config.NewConfig("ini", "conf/historyData.conf")
	if err != nil {

		presult["HistoryRecordDbType"] = 1
	}
	presult["PartitionType"], err = HistoryConf.String("PartitionType")
	if err != nil {
		presult["PartitionType"] = "0"
	}
	presult["HistoryRecordDbType"], err = HistoryConf.String("HistoryRecordDbType")
	if err != nil || presult["HistoryRecordDbType"] == "" {
		presult["HistoryRecordDbType"] = "1"
	}
	presult["OnceWriteHistoryCounts"], err = HistoryConf.String("OnceWriteHistoryCounts")
	if err != nil || presult["OnceWriteHistoryCounts"] == "" {
		presult["OnceWriteHistoryCounts"] = "100"
	}

	presult["TDenginePort"], err = HistoryConf.String("TDengine::TDenginePort")
	if err != nil {
		presult["TDenginePort"] = "6041"
	}
	presult["TDengineHost"], err = HistoryConf.String("TDengine::TDengineHost")
	if err != nil {
		presult["TDengineHost"] = "127.0.0.1"
	}
	presult["TDengineUserName"], err = HistoryConf.String("TDengine::UserName")
	if err != nil {
		presult["TDengineUserName"] = ""
	}
	presult["TDenginePassWord"], err = HistoryConf.String("TDengine::PassWord")
	if err != nil {
		presult["TDenginePassWord"] = ""
	}

	presult["ChickHousePort"], err = HistoryConf.String("ChickHouse::ChickHousePort")
	if err != nil {
		presult["ChickHousePort"] = "9000"
	}
	presult["ChickHouseHost"], err = HistoryConf.String("ChickHouse::ChickHouseHost")
	if err != nil {
		presult["ChickHouseHost"] = "127.0.0.1"
	}
	presult["ChickHouseUserName"], err = HistoryConf.String("ChickHouse::UserName")
	if err != nil {
		presult["ChickHouseUserName"] = ""
	}
	presult["ChickHousePassWord"], err = HistoryConf.String("ChickHouse::PassWord")
	if err != nil {
		presult["ChickHousePassWord"] = ""
	}
	presult["ChickHouseDataBase"], err = HistoryConf.String("ChickHouse::DataBase")
	if err != nil {
		presult["ChickHouseDataBase"] = ""
	}
	presult["ChickHouseConnectTimeout"], err = HistoryConf.String("ChickHouse::ConnectTimeout")
	if err != nil {
		presult["ChickHouseConnectTimeout"] = "60s"
	}
	presult["ChickHouseReadTimeout"], err = HistoryConf.String("ChickHouse::ReadTimeout")
	if err != nil {
		presult["ChickHouseReadTimeout"] = "60s"
	}

	presult["InfluxdbUrl"], err = HistoryConf.String("Influxdb::Url")
	if err != nil {
		presult["InfluxdbUrl"] = ""
	}
	presult["InfluxdbToken"], err = HistoryConf.String("Influxdb::Token")
	if err != nil {
		presult["InfluxdbToken"] = ""
	}
	presult["InfluxdbOrg"], err = HistoryConf.String("Influxdb::Org")
	if err != nil {
		presult["InfluxdbOrg"] = ""
	}
	presult["InfluxdbBucket"], err = HistoryConf.String("Influxdb::Bucket")
	if err != nil {
		presult["InfluxdbBucket"] = ""
	}

	presult["PGHost"], err = HistoryConf.String("PG::Host")
	if err != nil {
		presult["PGHost"] = ""
	}
	presult["PGPort"], err = HistoryConf.String("PG::Port")
	if err != nil {
		presult["PGPort"] = "5432"
	}
	presult["PGUser"], err = HistoryConf.String("PG::User")
	if err != nil {
		presult["PGUser"] = ""
	}
	presult["PGPassWord"], err = HistoryConf.String("PG::PassWord")
	if err != nil {
		presult["PGPassWord"] = ""
	}
	presult["PGDbName"], err = HistoryConf.String("PG::DbName")
	if err != nil {
		presult["PGDbName"] = ""
	}

	result := map[string]interface{}{
		"result": presult,
	}
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *SystemParamsController) SaveSystemHistoryConfig() {
	type ParamsData struct {
		HistoryRecordDbType      string `json:"HistoryRecordDbType"`
		OnceWriteHistoryCounts   string `json:"OnceWriteHistoryCounts"`
		PartitionType            string `json:"PartitionType"`
		TDengineHost             string `json:"TDengineHost"`
		TDenginePassWord         string `json:"TDenginePassWord"`
		TDenginePort             string `json:"TDenginePort"`
		TDengineUserName         string `json:"TDengineUserName"`
		ChickHousePort           string `json:"ChickHousePort"`
		ChickHouseHost           string `json:"ChickHouseHost"`
		ChickHouseUserName       string `json:"ChickHouseUserName"`
		ChickHousePassWord       string `json:"ChickHousePassWord"`
		ChickHouseDataBase       string `json:"ChickHouseDataBase"`
		ChickHouseConnectTimeout string `json:"ChickHouseConnectTimeout"`
		ChickHouseReadTimeout    string `json:"ChickHouseReadTimeout"`
		InfluxdbUrl              string `json:"InfluxdbUrl"`
		InfluxdbToken            string `json:"InfluxdbToken"`
		InfluxdbOrg              string `json:"InfluxdbOrg"`
		InfluxdbBucket           string `json:"InfluxdbBucket"`
		PGHost                   string `json:"PGHost"`
		PGPort                   string `json:"PGPort"`
		PGDbName                 string `json:"PGDbName"`
		PGUser                   string `json:"PGUser"`
		PGPassWord               string `json:"PGPassWord"`
	}
	var pdata ParamsData
	var code int = -1
	data := c.Ctx.Input.RequestBody

	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到user对象中
		err := json.Unmarshal(data, &pdata)
		if err != nil {
			code = -3
		} else {
			HistoryConf, err := config.NewConfig("ini", "conf/historyData.conf")
			if err != nil {
				code = -3
			} else {
				HistoryConf.Set("HistoryRecordDbType", pdata.HistoryRecordDbType)
				HistoryConf.Set("PartitionType", pdata.PartitionType)
				HistoryConf.Set("OnceWriteHistoryCounts", pdata.OnceWriteHistoryCounts)

				HistoryConf.Set("TDengine::TDenginePort", pdata.TDenginePort)
				HistoryConf.Set("TDengine::TDengineHost", pdata.TDengineHost)
				HistoryConf.Set("TDengine::UserName", pdata.TDengineUserName)
				HistoryConf.Set("TDengine::PassWord", pdata.TDenginePassWord)

				HistoryConf.Set("ChickHouse::ChickHousePort", pdata.ChickHousePort)
				HistoryConf.Set("ChickHouse::ChickHouseHost", pdata.ChickHouseHost)
				HistoryConf.Set("ChickHouse::UserName", pdata.ChickHouseUserName)
				HistoryConf.Set("ChickHouse::PassWord", pdata.ChickHousePassWord)
				HistoryConf.Set("ChickHouse::DataBase", pdata.ChickHouseDataBase)
				HistoryConf.Set("ChickHouse::ConnectTimeout", pdata.ChickHouseConnectTimeout)
				HistoryConf.Set("ChickHouse::ReadTimeout", pdata.ChickHouseReadTimeout)

				HistoryConf.Set("Influxdb::Url", pdata.InfluxdbUrl)
				HistoryConf.Set("Influxdb::Token", pdata.InfluxdbToken)
				HistoryConf.Set("Influxdb::Org", pdata.InfluxdbOrg)
				HistoryConf.Set("Influxdb::Bucket", pdata.InfluxdbBucket)

				HistoryConf.Set("PG::Host", pdata.PGHost)
				HistoryConf.Set("PG::Port", pdata.PGPort)
				HistoryConf.Set("PG::User", pdata.PGUser)
				HistoryConf.Set("PG::PassWord", pdata.PGPassWord)
				HistoryConf.Set("PG::DbName", pdata.PGDbName)

				HistoryConf.SaveConfigFile("conf/historyData.conf")
				code = 0
			}

			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "修改了历史数据库 "+pdata.HistoryRecordDbType, errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -2
	}
	result := map[string]interface{}{
		"code":   code,
		"result": "",
	}
	c.Data["json"] = result
	c.ServeJSON()
	go func() {
		timer := time.NewTimer(3 * time.Second)
		<-timer.C
		logs.Info("配置更新,准备重启启动...")
		time.Sleep(100 * time.Millisecond)
		os.Exit(0)
	}()
}
func (c *SystemParamsController) SaveSystemTimeConfig() {
	type ParamsData struct {
		TimeZone     string `json:"TimeZone"`
		CheckType    int    `json:"CheckType"`
		NTPServer    string `json:"NTPServer"`
		NTPPort      int    `json:"NTPPort"`
		NTPCheckTime int    `json:"NTPCheckTime"`
		SetTime      string `json:"SetTime"`
	}
	var pdata ParamsData
	var code int = -1
	data := c.Ctx.Input.RequestBody

	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到user对象中
		err := json.Unmarshal(data, &pdata)
		if err != nil {
			code = -3
		} else {
			timeConf, err := config.NewConfig("ini", "conf/systimeconfig.conf")
			if err != nil {
				code = -4
			} else {
				timeConf.Set("TimeZone", pdata.TimeZone)
				timeConf.Set("CheckType", fmt.Sprintf("%d", pdata.CheckType))
				timeConf.Set("NTPServer", pdata.NTPServer)
				timeConf.Set("NTPPort", fmt.Sprintf("%d", pdata.NTPPort))
				timeConf.Set("NTPCheckTime", fmt.Sprintf("%d", pdata.NTPCheckTime))
				timeConf.SaveConfigFile("conf/systimeconfig.conf")
				setSystemTimeZone(pdata.TimeZone)
				if pdata.CheckType == 1 {
					err := syncSystemTime(pdata.NTPServer, pdata.NTPPort)
					if err != nil {
						code = -5
					} else {
						code = 0
						logs.Info("系统时间设置成功")
					}
				} else if pdata.CheckType == 0 {
					err := hangSyncSystemTime(pdata.SetTime)
					if err != nil {
						code = -5
					} else {
						code = 0
						logs.Info("系统时间设置成功")
					}
				}
			}
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "修改了系统时间 "+pdata.NTPServer, errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	} else {
		code = -2
	}
	result := map[string]interface{}{
		"code":   code,
		"result": "",
	}
	c.Data["json"] = result
	c.ServeJSON()
}
