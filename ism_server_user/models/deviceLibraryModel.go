/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-08-24 10:29:08
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package models

import (
	protocol_common "ISMServer/protocol/common"
	"ISMServer/utils/errmsg"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/go-basic/uuid"
	"gorm.io/gorm"
)

// 数据模型结构体
type DevicesModel struct {
	gorm.Model

	Name             string `gorm:"index;type:varchar(250);not null" json:"name" validate:"required,min=4,max=250" label:"模型名称"`
	Described        string `gorm:"type:varchar(250);not null" json:"dec" validate:"required,min=1,max=250" label:"描述"`
	Uuid             string `gorm:"index;type:varchar(250);not null" json:"uuid" validate:"required" label:"记录的UUID"`
	Type             int    `gorm:"type:int;DEFAULT:1;" json:"type" validate:"required" label:"协议类型"`
	GatherNumber     int    `gorm:"type:int;DEFAULT:30;not null" json:"gatherNumber" validate:"required" label:"一次采集数量"`
	ProjectUuid      string `gorm:"index;type:varchar(250);not null" json:"project_uuid" validate:"required" label:"项目的UUID"`
	ConfigurationUid string `gorm:"index;type:varchar(250);" json:"configUid" validate:"required,min=1,max=250" label:"展示模型uuid"`
	PageUUID         string `gorm:"index;type:varchar(250);" json:"PageUUID" validate:"required,min=1,max=250" label:"展示模型页面uuid"`
	//SNMP协议
	Version              int    `gorm:"type:int;DEFAULT:2;" json:"version" validate:"required" label:"协议版本"`
	Port                 int    `gorm:"type:int;DEFAULT:161;" json:"port" validate:"required" label:"端口"`
	Writecomm            string `gorm:"type:varchar(250);" json:"writecomm" validate:"required,min=2,max=250" label:"写公共体"`
	Readcomm             string `gorm:"type:varchar(250);" json:"readcomm" validate:"required,min=2,max=250" label:"读公共体"`
	SnmpUserName         string `gorm:"type:varchar(250);" json:"snmpUserName" validate:"required,min=2,max=250" label:"V3 用户名"`
	SnmpSecurityLevel    int    `gorm:"type:int;;" json:"snmpSecurityLevel" validate:"required" label:"v3 鉴权方式"`
	SnmpAuthAlgorithm    int    `gorm:"type:int;;" json:"snmpAuthAlgorithm" validate:"required" label:"v3 鉴权加密算法"`
	SnmpUserPassword     string `gorm:"type:varchar(250);" json:"snmpUserPassword" validate:"required,min=2,max=250" label:"V3 密码"`
	SnmpPrivacyAlgorithm int    `gorm:"type:int;;" json:"snmpPrivacyAlgorithm" validate:"required" label:"v3 隐私加密算法"`
	SnmpPrivacyPassword  string `gorm:"type:varchar(250);" json:"snmpPrivacyPassword" validate:"required,min=2,max=250"  label:"v3 隐私密码"`

	//modbus
	ModbusConnectType        string `gorm:"type:varchar(250);not null;DEFAULT:''" json:"modbusConnectType" validate:"required" label:"modbus连接方式"`
	ModbusConnectMode        string `gorm:"type:varchar(250);not null;DEFAULT:''" json:"modbusConnectMode" validate:"required" label:"modbus模式"`
	ModbusConnectCOMName     string `gorm:"type:varchar(250);" json:"modbusCom" validate:"required" label:"连接的串口名称"`
	ModbusSerialBaud         int    `gorm:"type:int;DEFAULT:9600;" json:"serialBaud" validate:"required" label:"modbus串口波特率"`
	ModbusSerialBits         int    `gorm:"type:int;" json:"serialBits" validate:"required" label:"modbus串口数据位"`
	ModbusSerialParity       string `gorm:"type:varchar(250);" json:"serialParity" validate:"required" label:"modbus串口校验位"`
	ModbusSerialStopBits     string `gorm:"type:varchar(250);" json:"serialStopBits" validate:"required" label:"modbus串口停止位"`
	ModbusSerialFlow         string `gorm:"type:varchar(250);" json:"serialFlow" validate:"required" label:"modbus串口流控"`
	Timeout                  int    `gorm:"type:int;" json:"timeout" validate:"required" label:"连接超时"`
	DataFormat               string `gorm:"type:varchar(250);" json:"DataFormat" validate:"required" label:"数据格式"`
	ModbusTCPClientIpaddress string `gorm:"type:varchar(250);" json:"modbusClientIpaddress" validate:"required" label:"modbus客户端IP"`

	//OPCUA
	OPCUAConnectType      int `gorm:"type:int;not null;DEFAULT:1" json:"OPCUAConnectType" validate:"required" label:"OPCUA连接方式"`
	OPCUASecurityPolicies int `gorm:"type:int;DEFAULT:1" json:"OPCUASecurityPolicies" validate:"required" label:"OPCUA 加密模型"`
	OPCUASecurityModes    int `gorm:"type:int;not null;DEFAULT:1" json:"OPCUASecurityModes" validate:"required" label:"OPCUA安全模式"`
	OPCUAAuthModes        int `gorm:"type:int;not null;DEFAULT:1" json:"OPCUAAuthModes" validate:"required" label:"OPCUA授权模式"`

	OPCUATLSPolicies     int    `gorm:"type:int;DEFAULT:0" json:"OPCUATLSPolicies" validate:"required" label:"OPCUA TLS版本"`
	OPCUAConnectUserName string `gorm:"type:varchar(250);" json:"OPCUAConnectUserName" validate:"required" label:"OPCUA连接用户名"`
	OPCUAConnectPassword string `gorm:"type:varchar(250);" json:"OPCUAConnectPassword" validate:"required" label:"OPCUA连接密码"`
	OPCUACertificatePath string `gorm:"type:varchar(250);" json:"OPCUACertificatePath" validate:"required" label:"OPCUA连接证书路径"`
	OPCUAPrivateKeyPath  string `gorm:"type:varchar(250);" json:"OPCUAPrivateKeyPath" validate:"required" label:"OPCUA连接私阴路径"`

	//DLT645
	DLT645ConnectType        string `gorm:"type:varchar(250);not null;DEFAULT:''" json:"DLT645ConnectType" validate:"required" label:"DLT645连接方式"`
	DLT645ConnectMode        string `gorm:"type:varchar(250);not null;DEFAULT:''" json:"DLT645ConnectMode" validate:"required" label:"DLT645模式"`
	DLT645ConnectCOMName     string `gorm:"type:varchar(250);" json:"DLT645ConnectCOMName" validate:"required" label:"连接的串口名称"`
	DLT645SerialBaud         int    `gorm:"type:int;DEFAULT:9600;" json:"DLT645SerialBaud" validate:"required" label:"DLT645串口波特率"`
	DLT645SerialBits         int    `gorm:"type:int;" json:"DLT645SerialBits" validate:"required" label:"DLT645串口数据位"`
	DLT645SerialParity       string `gorm:"type:varchar(250);" json:"DLT645SerialParity" validate:"required" label:"DLT645串口校验位"`
	DLT645SerialStopBits     string `gorm:"type:varchar(250);" json:"DLT645SerialStopBits" validate:"required" label:"DLT645串口停止位"`
	DLT645SerialFlow         string `gorm:"type:varchar(250);" json:"DLT645SerialFlow" validate:"required" label:"DLT645串口流控"`
	DLT645Timeout            int    `gorm:"type:int;" json:"DLT645Timeout" validate:"required" label:"连接超时"`
	DLT645DataFormat         string `gorm:"type:varchar(250);" json:"DLT645DataFormat" validate:"required" label:"数据格式"`
	DLT645TCPClientIpaddress string `gorm:"type:varchar(250);" json:"DLT645TCPClientIpaddress" validate:"required" label:"DLT645客户端IP"`

	//Mqtt
	MqttSetDataFormat string `gorm:"type:text;" json:"MqttSetDataFormat" validate:"required" label:"MQtt下发的数据格式"`

	//HJ212
	HJ212ConnectType string `gorm:"type:text;" json:"HJ212ConnectType" validate:"required" label:"连接类型"`

	//CJT188
	CJT188ConnectType        string `gorm:"type:varchar(250);not null;DEFAULT:''" json:"CJT188ConnectType" validate:"required" label:"DLT645连接方式"`
	CJT188ConnectMode        string `gorm:"type:varchar(250);not null;DEFAULT:''" json:"CJT188ConnectMode" validate:"required" label:"DLT645模式"`
	CJT188ConnectCOMName     string `gorm:"type:varchar(250);" json:"CJT188ConnectCOMName" validate:"required" label:"连接的串口名称"`
	CJT188SerialBaud         int    `gorm:"type:int;DEFAULT:9600;" json:"CJT188SerialBaud" validate:"required" label:"DLT645串口波特率"`
	CJT188SerialBits         int    `gorm:"type:int;" json:"CJT188SerialBits" validate:"required" label:"DLT645串口数据位"`
	CJT188SerialParity       string `gorm:"type:varchar(250);" json:"CJT188SerialParity" validate:"required" label:"DLT645串口校验位"`
	CJT188SerialStopBits     string `gorm:"type:varchar(250);" json:"CJT188SerialStopBits" validate:"required" label:"DLT645串口停止位"`
	CJT188SerialFlow         string `gorm:"type:varchar(250);" json:"CJT188SerialFlow" validate:"required" label:"DLT645串口流控"`
	CJT188Timeout            int    `gorm:"type:int;" json:"CJT188Timeout" validate:"required" label:"连接超时"`
	CJT188DataFormat         string `gorm:"type:varchar(250);" json:"CJT188DataFormat" validate:"required" label:"数据格式"`
	CJT188TCPClientIpaddress string `gorm:"type:varchar(250);" json:"CJT188TCPClientIpaddress" validate:"required" label:"DLT645客户端IP"`
}

// 设备告警表
type DevicesAlarmList struct {
	gorm.Model

	AlarmName         string    `gorm:"type:varchar(250);not null" json:"AlarmName" validate:"required,min=4,max=250" label:"告警名称"`
	DeviceUuid        string    `gorm:"type:varchar(250);not null" json:"DeviceUuid" label:"设备UUID"`
	ProjectUuid       string    `gorm:"type:varchar(250);not null" json:"project_uuid" validate:"required" label:"项目的UUID"`
	DeviceName        string    `gorm:"type:varchar(250);not null" json:"DeviceName" validate:"required" label:"设备名称"`
	DataUuid          string    `gorm:"type:varchar(250);not null" json:"DataUuid" validate:"required" label:"数据UUID"`
	ModelDataUuid     string    `gorm:"type:varchar(250);not null" json:"ModelDataUuid" validate:"required" label:"模型数据UUID"`
	HappenTime        time.Time `gorm:"type:datetime;not null" json:"HappenTime" validate:"required" label:"发生时间"`
	ClearTime         time.Time `gorm:"type:datetime;" json:"ClearTime" validate:"required" label:"清除时间"`
	KeepTime          float64   `gorm:"type:double;" json:"KeepTime" validate:"required" label:"告警保持时间"`
	AlarmMessage      string    `gorm:"type:text;" json:"AlarmMessage" validate:"required" label:"告警显示"`
	AlarmClearMessage string    `gorm:"type:text;" json:"AlarmClearMessage" validate:"required" label:"告警消除显示"`
	AlarmLevel        int       `gorm:"type:int;" json:"AlarmLevel" validate:"required" label:"告警等级"`
}

// 设备历史数据
type DevicesHistoryDataList struct {
	gorm.Model

	DataName         string    `gorm:"index;type:varchar(250);not null" json:"DataName" validate:"required,min=4,max=250" label:"数据名称"`
	DeviceUuid       string    `gorm:"index;type:varchar(250);not null" json:"DeviceUuid" label:"设备UUID"`
	ProjectUuid      string    `gorm:"index;type:varchar(250);not null" json:"project_uuid" validate:"required" label:"项目的UUID"`
	DeviceName       string    `gorm:"type:varchar(250);not null" json:"DeviceName" validate:"required" label:"设备名称"`
	DataUuid         string    `gorm:"index;type:varchar(250);not null" json:"DataUuid" validate:"required" label:"数据UUID"`
	ModelDataUuid    string    `gorm:"index;type:varchar(250);not null" json:"ModelDataUuid" validate:"required" label:"模型数据UUID"`
	RecordTime       time.Time `gorm:"index;type:datetime;not null" json:"RecordTime" validate:"required" label:"记录时间"`
	DataUnit         string    `gorm:"type:varchar(250);" json:"DataUnit" validate:"required" label:"数据单位"`
	DataValue        string    `gorm:"type:varchar(250);not null" json:"DataValue" validate:"required" label:"记录数据"`
	RecordInterval   int       `gorm:"type:int;" json:"RecordInterval" validate:"required" label:"记录间隔"`
	RecordType       int       `gorm:"index;type:int;" json:"RecordType" validate:"required" label:"存储方式"`
	RecordDataCharge string    `gorm:"type:varchar(250);" json:"RecordDataCharge" validate:"required" label:""`
	RecordDataTimely string    `gorm:"type:varchar(250);" json:"RecordDataTimely" validate:"required" label:""`
}

// 设备历史数据
type DevicesCHHistoryData struct {
	DataName      string    `gorm:"index;type:varchar(250);not null" json:"DataName" validate:"required,min=4,max=250" label:"数据名称"`
	DeviceUuid    string    `gorm:"index;type:varchar(250);not null" json:"DeviceUuid" label:"设备UUID"`
	ProjectUuid   string    `gorm:"index;type:varchar(250);not null" json:"project_uuid" validate:"required" label:"项目的UUID"`
	DeviceName    string    `gorm:"index;type:varchar(250);not null" json:"DeviceName" validate:"required" label:"设备名称"`
	DataUuid      string    `gorm:"index;type:varchar(250);not null" json:"DataUuid" validate:"required" label:"数据UUID"`
	ModelDataUuid string    `gorm:"index;type:varchar(250);not null" json:"ModelDataUuid" validate:"required" label:"模型数据UUID"`
	RecordTime    time.Time `gorm:"index;type:datetime;not null" json:"RecordTime" validate:"required" label:"记录时间"`
	DataUnit      string    `gorm:"type:varchar(250);" json:"DataUnit" validate:"required" label:"数据单位"`
	DataValue     string    `gorm:"type:varchar(250);not null" json:"DataValue" validate:"required" label:"记录数据"`
}

// 设备历史数据
type DevicesPgHistoryData struct {
	DataName      string    `gorm:"index;type:varchar(250);not null" json:"DataName" validate:"required,min=4,max=250" label:"数据名称"`
	DeviceUuid    string    `gorm:"index;type:varchar(250);not null" json:"DeviceUuid" label:"设备UUID"`
	ProjectUuid   string    `gorm:"index;type:varchar(250);not null" json:"project_uuid" validate:"required" label:"项目的UUID"`
	DeviceName    string    `gorm:"index;type:varchar(250);not null" json:"DeviceName" validate:"required" label:"设备名称"`
	DataUuid      string    `gorm:"index;type:varchar(250);not null" json:"DataUuid" validate:"required" label:"数据UUID"`
	ModelDataUuid string    `gorm:"index;type:varchar(250);not null" json:"ModelDataUuid" validate:"required" label:"模型数据UUID"`
	RecordTime    time.Time `gorm:"index;type:TIMESTAMPTZ;not null" json:"RecordTime" validate:"required" label:"记录时间"`
	DataUnit      string    `gorm:"type:varchar(250);" json:"DataUnit" validate:"required" label:"数据单位"`
	DataValue     string    `gorm:"type:varchar(250);not null" json:"DataValue" validate:"required" label:"记录数据"`
}

// 监控树表
type MonitorList struct {
	gorm.Model

	Sid                 int32  `gorm:"index;type:int;not null" json:"sid" validate:"required" label:"节点的ID"`
	Pid                 int32  `gorm:"index;type:int;not null" json:"pid" validate:"required" label:"父节点的ID"`
	Name                string `gorm:"index;type:varchar(250);not null;" json:"name" validate:"required,min=4,max=250" label:"名称"`
	Type                int    `gorm:"type:int;not null" json:"type" validate:"required" label:"节点类型"`
	Timeout             int    `gorm:"type:int;DEFAULT:5;not null" json:"timeout" validate:"required" label:"超时"`
	IsEnable            int    `gorm:"index;type:int;DEFAULT:1;not null" json:"IsEnable" validate:"required" label:"是否启停"`
	ProjectUuid         string `gorm:"index;type:varchar(250);not null" json:"project_uuid" validate:"required" label:"项目的UUID"`
	Interval            int    `gorm:"type:int;DEFAULT:5;not null" json:"interval" validate:"required" label:"采集间隔"`
	FailedTimes         int    `gorm:"type:int;DEFAULT:5;not null" json:"failedTimes" validate:"required" label:"通信失败次数"`
	Described           string `gorm:"type:varchar(250);" json:"description" validate:"required,min=1,max=250" label:"描述"`
	OfflineClear        int    `gorm:"type:int;DEFAULT:0;not null" json:"offlineClear" validate:"required" label:"离线是否恢复默认值"`
	OfflineDefaultValue string `gorm:"type:varchar(250);" json:"offlineDefaultValue" validate:"required,min=1,max=250" label:"离线默认值"`
	DeviceType          int    `gorm:"type:int;" json:"deviceType" validate:"required" label:"节点设备类型"`
	Muid                string `gorm:"index;type:varchar(250);" json:"muid" validate:"required,min=1,max=250" label:"数据模型uuid"`
	ConfigurationUid    string `gorm:"type:varchar(250);" json:"configUid" validate:"required,min=1,max=250" label:"展示模型uuid"`
	PageUUID            string `gorm:"index;type:varchar(250);" json:"PageUUID" validate:"required,min=1,max=250" label:"展示模型页面uuid"`
	Uuid                string `gorm:"index;type:varchar(250);" json:"uuid" validate:"required,min=1,max=250" label:"设备uuid"`
	ExtraData           string `gorm:"type:text;" json:"extra" validate:"required,min=1,max=250" label:"额外数据可以存储其他数据类型"`
	Status              int    `gorm:"type:int;DEFAULT:0;not null" json:"Status" validate:"required" label:"在线状态"`
	Longitude           string `gorm:"type:varchar(250)" json:"longitude" validate:"required" label:"经度"`
	Latitude            string `gorm:"type:varchar(250)" json:"latitude" validate:"required" label:"纬度"`
	ProjectUUID         string `gorm:"-" json:"ProjectUUID"`
	IsRemoteGw          bool   `gorm:"-" json:"IsRemoteGw"`
}

type Menu struct {
	Title string      `json:"title"`
	Key   string      `json:"key"`
	Value MonitorList `json:"value"`
}

type TreeList struct {
	Text        string                 `json:"text"`
	Key         string                 `json:"key"`
	IconCls     string                 `json:"iconCls"`
	Id          int                    `json:"id"`
	Value       MonitorList            `json:"value"`
	ScopedSlots map[string]interface{} `json:"scopedSlots"`
	Children    []*TreeList            `json:"children"`
}
type GetTreeAndConfig struct {
	MonitorList
	ProjectUuid      string `gorm:"type:varchar(250);not null" json:"project_uuid" validate:"required" label:"项目的UUID"`
	ConfigurationUid string `gorm:"type:varchar(250);" json:"configUid" validate:"required,min=1,max=250" label:"展示模型uuid"`
}

type ConfigStu struct {
	ConfigurationUid string `json:"configUid"  label:"展示模型uuid"`
	PageUUID         string `json:"PageUUID" label:"展示模型页面uuid"`
}

// 获取系统支持的设备类型列表
func GetSupportDeviceList() []DevicesSupportList {

	var getDevicesSupputList []DevicesSupportList

	Db.Model(&DevicesSupportList{}).Select("*").Where("type!=6 and type!=7").Find(&getDevicesSupputList)

	return getDevicesSupputList
}
func monitTree(menu []*GetTreeAndConfig, pid int32, IsRemoteGw bool, ProjectUUID string) []*TreeList {
	treeList := []*TreeList{}
	for _, v := range menu {
		if v.Pid == pid {
			ScopedSlots := make(map[string]interface{}, 1)
			ScopedSlots["title"] = v.Name

			v.MonitorList.ConfigurationUid = v.ConfigurationUid
			v.MonitorList.PageUUID = v.PageUUID
			v.MonitorList.ProjectUUID = ProjectUUID
			v.MonitorList.IsRemoteGw = IsRemoteGw
			node := &TreeList{
				Id:          int(v.ID),
				Key:         v.Uuid,
				IconCls:     "icon-remove-3",
				Text:        v.Name,
				ScopedSlots: ScopedSlots,
				Value:       v.MonitorList,
			}
			node.Children = append(node.Children, monitTree(menu, v.Sid, IsRemoteGw, ProjectUUID)...)
			treeList = append(treeList, node)
		}
	}
	return treeList
}

/*
*
递归获取树形菜单
*/
func GetAllDevices(pid int32, ProjectUuid string, IsRemoteGw bool) []*TreeList {
	var ZoneMenu []*GetTreeAndConfig
	var DeviceMenu []*GetTreeAndConfig
	Db.Raw("SELECT * FROM  monitor_list WHERE project_uuid=? and type=0", ProjectUuid).Scan(&ZoneMenu)
	Db.Raw("SELECT monitor_list.*,devices_model.configuration_uid,devices_model.page_uuid FROM devices_model, monitor_list WHERE monitor_list.muid = devices_model.uuid and monitor_list.project_uuid=? and monitor_list.type=1", ProjectUuid).Scan(&DeviceMenu)
	DeviceMenu = append(DeviceMenu, ZoneMenu...)
	treeList := monitTree(DeviceMenu, pid, IsRemoteGw, ProjectUuid)
	return treeList
}

/*
*
添加设备或者区域
*/
func AddDeviceOrZone(params MonitorList) int {

	var existList MonitorList

	if params.Type == 1 {
		existError := Db.Model(&MonitorList{}).Where("name = ?", params.Name).First(&existList)

		if !errors.Is(existError.Error, gorm.ErrRecordNotFound) {
			//添加的资源和设备已经存在
			return errmsg.ERROR_DEVICE_EXIST
		}
	}
	if params.Pid != 0 {
		existError := Db.Model(&MonitorList{}).Where("sid = ?", params.Pid).First(&existList)
		if errors.Is(existError.Error, gorm.ErrRecordNotFound) {
			//添加的资源和设备已经存在
			return errmsg.ERROR
		}
	}
	if params.DeviceType == 480 {
		params.Status = 1
	} else {
		params.Status = 2
	}
	if params.Type == 1 {
		if params.DeviceType != 5 {
			addressError := Db.Model(&MonitorList{}).Where("muid = ? AND extra_data = ?", params.Muid, params.ExtraData).First(&existList)

			if !errors.Is(addressError.Error, gorm.ErrRecordNotFound) {
				//添加的资源和设备已经存在
				return errmsg.ERROR_DEVICE_ADDRESS_EXIST
			}
		}
	}

	rand.Seed(time.Now().UnixNano())
	params.Sid = rand.Int31()
	params.Uuid = uuid.New()

	if params.Type == 1 {
		//设备静态数据
		var getDevicePublicStaticData []StaticData
		var getDeviceStaticData []StaticData
		var writeDeviceRealData []DeviceRealData
		var getAlarmTrigger []AlarmTrigger
		var getCustomData []CustomData
		staticPublicError := Db.Model(&StaticData{}).Where("data_device_type = ? and project_uuid = ?", 158, params.ProjectUuid).Find(&getDevicePublicStaticData).Error
		if staticPublicError != nil {
			return errmsg.ERROR
		}
		staticDeviceError := Db.Model(&StaticData{}).Where("project_uuid = ? and data_device_type = ?", params.ProjectUuid, params.DeviceType).Find(&getDeviceStaticData).Error
		if staticDeviceError != nil {
			return errmsg.ERROR
		}
		TriggerDeviceError := Db.Model(&AlarmTrigger{}).Where("project_uuid = ? and trigger_device_type = ? and trigger_device_model_uuid = ?", params.ProjectUuid, params.DeviceType, params.Muid).Find(&getAlarmTrigger).Error
		if TriggerDeviceError != nil {
			return errmsg.ERROR
		}
		CustomDeviceError := Db.Model(&CustomData{}).Where("project_uuid = ? and device_type = ? and select_data_model_u_uid = ?", params.ProjectUuid, params.DeviceType, params.Muid).Find(&getCustomData).Error
		if CustomDeviceError != nil {
			return errmsg.ERROR
		}
		//snmp设备类型
		if params.DeviceType == 1 {
			var getDeviceData []SnmpDevicesDataModel
			err := Db.Model(&SnmpDevicesDataModel{}).Where("muid = ?", params.Muid).Find(&getDeviceData).Error
			if err != nil {
				return errmsg.ERROR_DATABASE
			}
			for key, _ := range getDeviceData {
				var tempDeviceRealData DeviceRealData
				if getDeviceData[key].Auth == "ReadOnly" {
					tempDeviceRealData.Auth = 1
				} else if getDeviceData[key].Auth == "ReadWrite" {
					tempDeviceRealData.Auth = 2
				} else {
					tempDeviceRealData.Auth = 3
				}
				tempDeviceRealData.DataUnit = getDeviceData[key].DataUnit
				tempDeviceRealData.ProjectUuid = params.ProjectUuid
				tempDeviceRealData.DeviceName = params.Name
				tempDeviceRealData.Name = getDeviceData[key].Name
				tempDeviceRealData.Oid = getDeviceData[key].Oid
				tempDeviceRealData.Uuid = uuid.New()
				tempDeviceRealData.ModelDataUuid = getDeviceData[key].Uuid

				isInteger := strings.Contains(getDeviceData[key].OidType, "Integer")

				if isInteger {
					tempDeviceRealData.Type = 1
				} else if getDeviceData[key].OidType == "OctetString" {
					tempDeviceRealData.Type = 2
				} else {
					tempDeviceRealData.Type = 3
				}

				tempDeviceRealData.Value = ""
				tempDeviceRealData.Muid = params.Muid
				tempDeviceRealData.DeviceUuid = params.Uuid
				tempDeviceRealData.DeviceType = 1

				tempDeviceRealData.IsRecord = getDeviceData[key].IsRecord
				tempDeviceRealData.RecordInterval = getDeviceData[key].RecordInterval
				tempDeviceRealData.RecordDataCharge = getDeviceData[key].RecordDataCharge
				tempDeviceRealData.RecordType = getDeviceData[key].RecordType

				tempDeviceRealData.IsAlarm = getDeviceData[key].IsAlarm
				tempDeviceRealData.AlarmLevel = getDeviceData[key].AlarmLevel
				tempDeviceRealData.AlarmMessage = getDeviceData[key].AlarmMessage
				tempDeviceRealData.AlarmClearMessage = getDeviceData[key].AlarmClearMessage

				writeDeviceRealData = append(writeDeviceRealData, tempDeviceRealData)
			}
		} else if params.DeviceType == 2 { //modbus设备类型

			var getDeviceData []ModbusDevicesDataModel
			err := Db.Model(&ModbusDevicesDataModel{}).Where("muid = ?", params.Muid).Find(&getDeviceData).Error
			if err != nil {
				return errmsg.ERROR_DATABASE
			}
			for key, _ := range getDeviceData {
				var tempDeviceRealData DeviceRealData
				if getDeviceData[key].Auth == "ReadOnly" {
					tempDeviceRealData.Auth = 1
				} else if getDeviceData[key].Auth == "ReadWrite" {
					tempDeviceRealData.Auth = 2
				} else {
					tempDeviceRealData.Auth = 3
				}
				tempDeviceRealData.DataUnit = getDeviceData[key].DataUnit
				tempDeviceRealData.ProjectUuid = params.ProjectUuid
				tempDeviceRealData.DeviceName = params.Name
				tempDeviceRealData.Name = getDeviceData[key].Name
				tempDeviceRealData.Uuid = uuid.New()
				tempDeviceRealData.ModelDataUuid = getDeviceData[key].Uuid
				tempDeviceRealData.Type = 1
				tempDeviceRealData.Value = ""
				tempDeviceRealData.Muid = params.Muid
				tempDeviceRealData.DeviceUuid = params.Uuid
				tempDeviceRealData.DeviceType = 2

				tempDeviceRealData.IsRecord = getDeviceData[key].IsRecord
				tempDeviceRealData.RecordInterval = getDeviceData[key].RecordInterval
				tempDeviceRealData.RecordDataCharge = getDeviceData[key].RecordDataCharge
				tempDeviceRealData.RecordType = getDeviceData[key].RecordType

				tempDeviceRealData.IsAlarm = getDeviceData[key].IsAlarm
				tempDeviceRealData.AlarmLevel = getDeviceData[key].AlarmLevel
				tempDeviceRealData.AlarmMessage = getDeviceData[key].AlarmMessage
				tempDeviceRealData.AlarmClearMessage = getDeviceData[key].AlarmClearMessage

				writeDeviceRealData = append(writeDeviceRealData, tempDeviceRealData)
			}

		} else if params.DeviceType == 3 { //OPCUA设备类型
			var getDeviceData []OpcuaDevicesDataModel
			err := Db.Model(&OpcuaDevicesDataModel{}).Where("muid = ?", params.Muid).Find(&getDeviceData).Error
			if err != nil {
				return errmsg.ERROR_DATABASE
			}
			for key, _ := range getDeviceData {
				var tempDeviceRealData DeviceRealData
				if getDeviceData[key].Auth == "ReadOnly" {
					tempDeviceRealData.Auth = 1
				} else if getDeviceData[key].Auth == "ReadWrite" {
					tempDeviceRealData.Auth = 2
				} else {
					tempDeviceRealData.Auth = 3
				}
				tempDeviceRealData.DataUnit = getDeviceData[key].DataUnit
				tempDeviceRealData.ProjectUuid = params.ProjectUuid
				tempDeviceRealData.DeviceName = params.Name
				tempDeviceRealData.Name = getDeviceData[key].Name
				tempDeviceRealData.Uuid = uuid.New()
				tempDeviceRealData.ModelDataUuid = getDeviceData[key].Uuid
				tempDeviceRealData.Type = 1
				tempDeviceRealData.Value = ""
				tempDeviceRealData.Muid = params.Muid
				tempDeviceRealData.DeviceUuid = params.Uuid
				tempDeviceRealData.DeviceType = 3
				tempDeviceRealData.IsRecord = getDeviceData[key].IsRecord
				tempDeviceRealData.RecordInterval = getDeviceData[key].RecordInterval
				tempDeviceRealData.RecordDataCharge = getDeviceData[key].RecordDataCharge
				tempDeviceRealData.RecordType = getDeviceData[key].RecordType

				tempDeviceRealData.IsAlarm = getDeviceData[key].IsAlarm
				tempDeviceRealData.AlarmLevel = getDeviceData[key].AlarmLevel
				tempDeviceRealData.AlarmMessage = getDeviceData[key].AlarmMessage
				tempDeviceRealData.AlarmClearMessage = getDeviceData[key].AlarmClearMessage
				writeDeviceRealData = append(writeDeviceRealData, tempDeviceRealData)
			}
		} else if params.DeviceType == 5 { //RESTFul设备类型
			var getDeviceData []RESTFulDataModel
			params.Status = 1
			err := Db.Model(&RESTFulDataModel{}).Where("muid = ?", params.Muid).Find(&getDeviceData).Error
			if err != nil {
				return errmsg.ERROR_DATABASE
			}
			for key, _ := range getDeviceData {
				var tempDeviceRealData DeviceRealData
				if getDeviceData[key].Auth == "ReadOnly" {
					tempDeviceRealData.Auth = 1
				} else if getDeviceData[key].Auth == "ReadWrite" {
					tempDeviceRealData.Auth = 2
				} else {
					tempDeviceRealData.Auth = 3
				}
				tempDeviceRealData.DataUnit = getDeviceData[key].DataUnit
				tempDeviceRealData.ProjectUuid = params.ProjectUuid
				tempDeviceRealData.DeviceName = params.Name
				tempDeviceRealData.Name = getDeviceData[key].Name
				tempDeviceRealData.Uuid = uuid.New()
				tempDeviceRealData.ModelDataUuid = getDeviceData[key].Uuid
				tempDeviceRealData.Type = 1
				tempDeviceRealData.Value = ""
				tempDeviceRealData.Muid = params.Muid
				tempDeviceRealData.DeviceUuid = params.Uuid
				tempDeviceRealData.DeviceType = 5
				tempDeviceRealData.IsRecord = getDeviceData[key].IsRecord
				tempDeviceRealData.RecordInterval = getDeviceData[key].RecordInterval
				tempDeviceRealData.RecordDataCharge = getDeviceData[key].RecordDataCharge
				tempDeviceRealData.RecordType = getDeviceData[key].RecordType

				tempDeviceRealData.IsAlarm = getDeviceData[key].IsAlarm
				tempDeviceRealData.AlarmLevel = getDeviceData[key].AlarmLevel
				tempDeviceRealData.AlarmMessage = getDeviceData[key].AlarmMessage
				tempDeviceRealData.AlarmClearMessage = getDeviceData[key].AlarmClearMessage
				writeDeviceRealData = append(writeDeviceRealData, tempDeviceRealData)
			}
		} else if params.DeviceType == 15 { //西门子S7设备模型
			var getDeviceData []SimS7DataModel
			err := Db.Model(&SimS7DataModel{}).Where("muid = ?", params.Muid).Find(&getDeviceData).Error
			if err != nil {
				return errmsg.ERROR_DATABASE
			}
			for key, _ := range getDeviceData {
				var tempDeviceRealData DeviceRealData
				if getDeviceData[key].Auth == "ReadOnly" {
					tempDeviceRealData.Auth = 1
				} else if getDeviceData[key].Auth == "ReadWrite" {
					tempDeviceRealData.Auth = 2
				} else {
					tempDeviceRealData.Auth = 3
				}
				tempDeviceRealData.DataUnit = getDeviceData[key].DataUnit
				tempDeviceRealData.ProjectUuid = params.ProjectUuid
				tempDeviceRealData.DeviceName = params.Name
				tempDeviceRealData.Name = getDeviceData[key].Name
				tempDeviceRealData.Uuid = uuid.New()
				tempDeviceRealData.ModelDataUuid = getDeviceData[key].Uuid
				tempDeviceRealData.Type = 1
				tempDeviceRealData.Value = ""
				tempDeviceRealData.Muid = params.Muid
				tempDeviceRealData.DeviceUuid = params.Uuid
				tempDeviceRealData.DeviceType = params.DeviceType
				tempDeviceRealData.IsRecord = getDeviceData[key].IsRecord
				tempDeviceRealData.RecordInterval = getDeviceData[key].RecordInterval
				tempDeviceRealData.RecordDataCharge = getDeviceData[key].RecordDataCharge
				tempDeviceRealData.RecordType = getDeviceData[key].RecordType

				tempDeviceRealData.IsAlarm = getDeviceData[key].IsAlarm
				tempDeviceRealData.AlarmLevel = getDeviceData[key].AlarmLevel
				tempDeviceRealData.AlarmMessage = getDeviceData[key].AlarmMessage
				tempDeviceRealData.AlarmClearMessage = getDeviceData[key].AlarmClearMessage
				writeDeviceRealData = append(writeDeviceRealData, tempDeviceRealData)
			}
		} else if params.DeviceType == 20 { //Mqtt设备类型
			var getDeviceData []MqttDevicesDataModel
			err := Db.Model(&MqttDevicesDataModel{}).Where("muid = ?", params.Muid).Find(&getDeviceData).Error
			if err != nil {
				return errmsg.ERROR_DATABASE
			}
			for key, _ := range getDeviceData {
				var tempDeviceRealData DeviceRealData
				if getDeviceData[key].Auth == "ReadOnly" {
					tempDeviceRealData.Auth = 1
				} else if getDeviceData[key].Auth == "ReadWrite" {
					tempDeviceRealData.Auth = 2
				} else {
					tempDeviceRealData.Auth = 3
				}
				tempDeviceRealData.DataUnit = getDeviceData[key].DataUnit
				tempDeviceRealData.ProjectUuid = params.ProjectUuid
				tempDeviceRealData.DeviceName = params.Name
				tempDeviceRealData.Name = getDeviceData[key].Name
				tempDeviceRealData.Uuid = uuid.New()
				tempDeviceRealData.ModelDataUuid = getDeviceData[key].Uuid
				tempDeviceRealData.Type = 1
				tempDeviceRealData.Value = ""
				tempDeviceRealData.Muid = params.Muid
				tempDeviceRealData.DeviceUuid = params.Uuid
				tempDeviceRealData.DeviceType = params.DeviceType
				tempDeviceRealData.IsRecord = getDeviceData[key].IsRecord
				tempDeviceRealData.RecordInterval = getDeviceData[key].RecordInterval
				tempDeviceRealData.RecordDataCharge = getDeviceData[key].RecordDataCharge
				tempDeviceRealData.RecordType = getDeviceData[key].RecordType

				tempDeviceRealData.IsAlarm = getDeviceData[key].IsAlarm
				tempDeviceRealData.AlarmLevel = getDeviceData[key].AlarmLevel
				tempDeviceRealData.AlarmMessage = getDeviceData[key].AlarmMessage
				tempDeviceRealData.AlarmClearMessage = getDeviceData[key].AlarmClearMessage
				writeDeviceRealData = append(writeDeviceRealData, tempDeviceRealData)
			}
		} else if params.DeviceType == 30 { //DLT645设备类型
			var getDeviceData []Dlt645DevicesDataModel
			err := Db.Model(&Dlt645DevicesDataModel{}).Where("muid = ?", params.Muid).Find(&getDeviceData).Error
			if err != nil {
				return errmsg.ERROR_DATABASE
			}
			for key, _ := range getDeviceData {
				var tempDeviceRealData DeviceRealData
				if getDeviceData[key].Auth == "ReadOnly" {
					tempDeviceRealData.Auth = 1
				} else if getDeviceData[key].Auth == "ReadWrite" {
					tempDeviceRealData.Auth = 2
				} else {
					tempDeviceRealData.Auth = 3
				}
				tempDeviceRealData.DataUnit = getDeviceData[key].DataUnit
				tempDeviceRealData.ProjectUuid = params.ProjectUuid
				tempDeviceRealData.DeviceName = params.Name
				tempDeviceRealData.Name = getDeviceData[key].Name
				tempDeviceRealData.Uuid = uuid.New()
				tempDeviceRealData.ModelDataUuid = getDeviceData[key].Uuid
				tempDeviceRealData.Type = 1
				tempDeviceRealData.Value = ""
				tempDeviceRealData.Muid = params.Muid
				tempDeviceRealData.DeviceUuid = params.Uuid
				tempDeviceRealData.DeviceType = params.DeviceType
				tempDeviceRealData.IsRecord = getDeviceData[key].IsRecord
				tempDeviceRealData.RecordInterval = getDeviceData[key].RecordInterval
				tempDeviceRealData.RecordDataCharge = getDeviceData[key].RecordDataCharge
				tempDeviceRealData.RecordType = getDeviceData[key].RecordType

				tempDeviceRealData.IsAlarm = getDeviceData[key].IsAlarm
				tempDeviceRealData.AlarmLevel = getDeviceData[key].AlarmLevel
				tempDeviceRealData.AlarmMessage = getDeviceData[key].AlarmMessage
				tempDeviceRealData.AlarmClearMessage = getDeviceData[key].AlarmClearMessage
				writeDeviceRealData = append(writeDeviceRealData, tempDeviceRealData)
			}
		} else if params.DeviceType == 40 { //IEC104设备类型
			var getDeviceData []IEC104DevicesDataModel
			err := Db.Model(&IEC104DevicesDataModel{}).Where("muid = ?", params.Muid).Find(&getDeviceData).Error
			if err != nil {
				return errmsg.ERROR_DATABASE
			}
			for key, _ := range getDeviceData {
				var tempDeviceRealData DeviceRealData
				if getDeviceData[key].Auth == "ReadOnly" {
					tempDeviceRealData.Auth = 1
				} else if getDeviceData[key].Auth == "ReadWrite" {
					tempDeviceRealData.Auth = 2
				} else {
					tempDeviceRealData.Auth = 3
				}
				tempDeviceRealData.DataUnit = getDeviceData[key].DataUnit
				tempDeviceRealData.ProjectUuid = params.ProjectUuid
				tempDeviceRealData.DeviceName = params.Name
				tempDeviceRealData.Name = getDeviceData[key].Name
				tempDeviceRealData.Uuid = uuid.New()
				tempDeviceRealData.ModelDataUuid = getDeviceData[key].Uuid
				tempDeviceRealData.Type = 1
				tempDeviceRealData.Value = ""
				tempDeviceRealData.Muid = params.Muid
				tempDeviceRealData.DeviceUuid = params.Uuid
				tempDeviceRealData.DeviceType = params.DeviceType
				tempDeviceRealData.IsRecord = getDeviceData[key].IsRecord
				tempDeviceRealData.RecordInterval = getDeviceData[key].RecordInterval
				tempDeviceRealData.RecordDataCharge = getDeviceData[key].RecordDataCharge
				tempDeviceRealData.RecordType = getDeviceData[key].RecordType

				tempDeviceRealData.IsAlarm = getDeviceData[key].IsAlarm
				tempDeviceRealData.AlarmLevel = getDeviceData[key].AlarmLevel
				tempDeviceRealData.AlarmMessage = getDeviceData[key].AlarmMessage
				tempDeviceRealData.AlarmClearMessage = getDeviceData[key].AlarmClearMessage
				writeDeviceRealData = append(writeDeviceRealData, tempDeviceRealData)
			}
		} else if params.DeviceType == 350 { //IEC61850设备类型
			var getDeviceData []IEC61850DevicesDataModel
			err := Db.Model(&IEC61850DevicesDataModel{}).Where("muid = ?", params.Muid).Find(&getDeviceData).Error
			if err != nil {
				return errmsg.ERROR_DATABASE
			}
			for key, _ := range getDeviceData {
				var tempDeviceRealData DeviceRealData
				if getDeviceData[key].Auth == "ReadOnly" {
					tempDeviceRealData.Auth = 1
				} else if getDeviceData[key].Auth == "ReadWrite" {
					tempDeviceRealData.Auth = 2
				} else {
					tempDeviceRealData.Auth = 3
				}
				tempDeviceRealData.DataUnit = getDeviceData[key].DataUnit
				tempDeviceRealData.ProjectUuid = params.ProjectUuid
				tempDeviceRealData.DeviceName = params.Name
				tempDeviceRealData.Name = getDeviceData[key].Name
				tempDeviceRealData.Uuid = uuid.New()
				tempDeviceRealData.ModelDataUuid = getDeviceData[key].Uuid
				tempDeviceRealData.Type = 1
				tempDeviceRealData.Value = ""
				tempDeviceRealData.Muid = params.Muid
				tempDeviceRealData.DeviceUuid = params.Uuid
				tempDeviceRealData.DeviceType = params.DeviceType
				tempDeviceRealData.IsRecord = getDeviceData[key].IsRecord
				tempDeviceRealData.RecordInterval = getDeviceData[key].RecordInterval
				tempDeviceRealData.RecordDataCharge = getDeviceData[key].RecordDataCharge
				tempDeviceRealData.RecordType = getDeviceData[key].RecordType

				tempDeviceRealData.IsAlarm = getDeviceData[key].IsAlarm
				tempDeviceRealData.AlarmLevel = getDeviceData[key].AlarmLevel
				tempDeviceRealData.AlarmMessage = getDeviceData[key].AlarmMessage
				tempDeviceRealData.AlarmClearMessage = getDeviceData[key].AlarmClearMessage
				writeDeviceRealData = append(writeDeviceRealData, tempDeviceRealData)
			}
		} else if params.DeviceType == 470 { //HJ212设备类型
			var getDeviceData []HJ212DevicesDataModel
			err := Db.Model(&HJ212DevicesDataModel{}).Where("muid = ?", params.Muid).Find(&getDeviceData).Error
			if err != nil {
				return errmsg.ERROR_DATABASE
			}
			for key, _ := range getDeviceData {
				var tempDeviceRealData DeviceRealData
				if getDeviceData[key].Auth == "ReadOnly" {
					tempDeviceRealData.Auth = 1
				} else if getDeviceData[key].Auth == "ReadWrite" {
					tempDeviceRealData.Auth = 2
				} else {
					tempDeviceRealData.Auth = 3
				}
				tempDeviceRealData.DataUnit = getDeviceData[key].DataUnit
				tempDeviceRealData.ProjectUuid = params.ProjectUuid
				tempDeviceRealData.DeviceName = params.Name
				tempDeviceRealData.Name = getDeviceData[key].Name
				tempDeviceRealData.Uuid = uuid.New()
				tempDeviceRealData.ModelDataUuid = getDeviceData[key].Uuid
				tempDeviceRealData.Type = 1
				tempDeviceRealData.Value = ""
				tempDeviceRealData.Muid = params.Muid
				tempDeviceRealData.DeviceUuid = params.Uuid
				tempDeviceRealData.DeviceType = params.DeviceType
				tempDeviceRealData.IsRecord = getDeviceData[key].IsRecord
				tempDeviceRealData.RecordInterval = getDeviceData[key].RecordInterval
				tempDeviceRealData.RecordDataCharge = getDeviceData[key].RecordDataCharge
				tempDeviceRealData.RecordType = getDeviceData[key].RecordType

				tempDeviceRealData.IsAlarm = getDeviceData[key].IsAlarm
				tempDeviceRealData.AlarmLevel = getDeviceData[key].AlarmLevel
				tempDeviceRealData.AlarmMessage = getDeviceData[key].AlarmMessage
				tempDeviceRealData.AlarmClearMessage = getDeviceData[key].AlarmClearMessage
				writeDeviceRealData = append(writeDeviceRealData, tempDeviceRealData)
			}
		} else if params.DeviceType == 480 { //虚拟设备
			var getDeviceData []VirtualDeviceDataModel
			err := Db.Model(&VirtualDeviceDataModel{}).Where("muid = ?", params.Muid).Find(&getDeviceData).Error
			if err != nil {
				return errmsg.ERROR_DATABASE
			}
			for key, _ := range getDeviceData {
				var tempDeviceRealData DeviceRealData
				if getDeviceData[key].Auth == "ReadOnly" {
					tempDeviceRealData.Auth = 1
				} else if getDeviceData[key].Auth == "ReadWrite" {
					tempDeviceRealData.Auth = 2
				} else {
					tempDeviceRealData.Auth = 3
				}
				tempDeviceRealData.DataUnit = getDeviceData[key].DataUnit
				tempDeviceRealData.ProjectUuid = params.ProjectUuid
				tempDeviceRealData.DeviceName = params.Name
				tempDeviceRealData.Name = getDeviceData[key].Name
				tempDeviceRealData.Uuid = uuid.New()
				tempDeviceRealData.ModelDataUuid = getDeviceData[key].Uuid
				tempDeviceRealData.Type = 1
				tempDeviceRealData.Value = ""
				tempDeviceRealData.Muid = params.Muid
				tempDeviceRealData.DeviceUuid = params.Uuid
				tempDeviceRealData.DeviceType = params.DeviceType
				tempDeviceRealData.IsRecord = getDeviceData[key].IsRecord
				tempDeviceRealData.RecordInterval = getDeviceData[key].RecordInterval
				tempDeviceRealData.RecordDataCharge = getDeviceData[key].RecordDataCharge
				tempDeviceRealData.RecordType = getDeviceData[key].RecordType

				tempDeviceRealData.IsAlarm = getDeviceData[key].IsAlarm
				tempDeviceRealData.AlarmLevel = getDeviceData[key].AlarmLevel
				tempDeviceRealData.AlarmMessage = getDeviceData[key].AlarmMessage
				tempDeviceRealData.AlarmClearMessage = getDeviceData[key].AlarmClearMessage
				writeDeviceRealData = append(writeDeviceRealData, tempDeviceRealData)
			}
		} else if params.DeviceType == 490 { //CJT188设备类型
			var getDeviceData []CJT188DevicesDataModel
			err := Db.Model(&CJT188DevicesDataModel{}).Where("muid = ?", params.Muid).Find(&getDeviceData).Error
			if err != nil {
				return errmsg.ERROR_DATABASE
			}
			for key, _ := range getDeviceData {
				var tempDeviceRealData DeviceRealData
				if getDeviceData[key].Auth == "ReadOnly" {
					tempDeviceRealData.Auth = 1
				} else if getDeviceData[key].Auth == "ReadWrite" {
					tempDeviceRealData.Auth = 2
				} else {
					tempDeviceRealData.Auth = 3
				}
				tempDeviceRealData.DataUnit = getDeviceData[key].DataUnit
				tempDeviceRealData.ProjectUuid = params.ProjectUuid
				tempDeviceRealData.DeviceName = params.Name
				tempDeviceRealData.Name = getDeviceData[key].Name
				tempDeviceRealData.Uuid = uuid.New()
				tempDeviceRealData.ModelDataUuid = getDeviceData[key].Uuid
				tempDeviceRealData.Type = 1
				tempDeviceRealData.Value = ""
				tempDeviceRealData.Muid = params.Muid
				tempDeviceRealData.DeviceUuid = params.Uuid
				tempDeviceRealData.DeviceType = params.DeviceType
				tempDeviceRealData.IsRecord = getDeviceData[key].IsRecord
				tempDeviceRealData.RecordInterval = getDeviceData[key].RecordInterval
				tempDeviceRealData.RecordDataCharge = getDeviceData[key].RecordDataCharge
				tempDeviceRealData.RecordType = getDeviceData[key].RecordType

				tempDeviceRealData.IsAlarm = getDeviceData[key].IsAlarm
				tempDeviceRealData.AlarmLevel = getDeviceData[key].AlarmLevel
				tempDeviceRealData.AlarmMessage = getDeviceData[key].AlarmMessage
				tempDeviceRealData.AlarmClearMessage = getDeviceData[key].AlarmClearMessage
				writeDeviceRealData = append(writeDeviceRealData, tempDeviceRealData)
			}
		} else if params.DeviceType == 500 { //BACnet设备类型
			var getDeviceData []BacnetDevicesDataModel
			err := Db.Model(&BacnetDevicesDataModel{}).Where("muid = ?", params.Muid).Find(&getDeviceData).Error
			if err != nil {
				return errmsg.ERROR_DATABASE
			}
			for key, _ := range getDeviceData {
				var tempDeviceRealData DeviceRealData
				if getDeviceData[key].Auth == "ReadOnly" {
					tempDeviceRealData.Auth = 1
				} else if getDeviceData[key].Auth == "ReadWrite" {
					tempDeviceRealData.Auth = 2
				} else {
					tempDeviceRealData.Auth = 3
				}
				tempDeviceRealData.DataUnit = getDeviceData[key].DataUnit
				tempDeviceRealData.ProjectUuid = params.ProjectUuid
				tempDeviceRealData.DeviceName = params.Name
				tempDeviceRealData.Name = getDeviceData[key].Name
				tempDeviceRealData.Uuid = uuid.New()
				tempDeviceRealData.ModelDataUuid = getDeviceData[key].Uuid
				tempDeviceRealData.Type = 1
				tempDeviceRealData.Value = ""
				tempDeviceRealData.Muid = params.Muid
				tempDeviceRealData.DeviceUuid = params.Uuid
				tempDeviceRealData.DeviceType = params.DeviceType
				tempDeviceRealData.IsRecord = getDeviceData[key].IsRecord
				tempDeviceRealData.RecordInterval = getDeviceData[key].RecordInterval
				tempDeviceRealData.RecordDataCharge = getDeviceData[key].RecordDataCharge
				tempDeviceRealData.RecordType = getDeviceData[key].RecordType

				tempDeviceRealData.IsAlarm = getDeviceData[key].IsAlarm
				tempDeviceRealData.AlarmLevel = getDeviceData[key].AlarmLevel
				tempDeviceRealData.AlarmMessage = getDeviceData[key].AlarmMessage
				tempDeviceRealData.AlarmClearMessage = getDeviceData[key].AlarmClearMessage
				writeDeviceRealData = append(writeDeviceRealData, tempDeviceRealData)
			}
		}

		//公共数据
		for _, PublicStaticData := range getDevicePublicStaticData {
			var tempDeviceRealData DeviceRealData
			tempDeviceRealData.Auth = 2
			tempDeviceRealData.ProjectUuid = params.ProjectUuid
			tempDeviceRealData.DeviceName = params.Name
			tempDeviceRealData.Name = PublicStaticData.Name
			tempDeviceRealData.Uuid = uuid.New()
			tempDeviceRealData.ModelDataUuid = PublicStaticData.Uuid

			tempDeviceRealData.Type = PublicStaticData.DataType

			tempDeviceRealData.Value = PublicStaticData.DataDefaultValue
			tempDeviceRealData.Muid = params.Muid
			tempDeviceRealData.DeviceUuid = params.Uuid
			tempDeviceRealData.DeviceType = params.DeviceType
			writeDeviceRealData = append(writeDeviceRealData, tempDeviceRealData)
		}
		//静态数据
		for _, DeviceStaticData := range getDeviceStaticData {
			var tempDeviceRealData DeviceRealData
			tempDeviceRealData.Auth = 2
			tempDeviceRealData.ProjectUuid = params.ProjectUuid
			tempDeviceRealData.DeviceName = params.Name
			tempDeviceRealData.Name = DeviceStaticData.Name
			tempDeviceRealData.Uuid = uuid.New()
			tempDeviceRealData.ModelDataUuid = DeviceStaticData.Uuid

			tempDeviceRealData.Type = DeviceStaticData.DataType

			tempDeviceRealData.Value = DeviceStaticData.DataDefaultValue
			tempDeviceRealData.Muid = params.Muid
			tempDeviceRealData.DeviceUuid = params.Uuid
			tempDeviceRealData.DeviceType = params.DeviceType

			tempDeviceRealData.IsAlarm = DeviceStaticData.DataAlarm
			tempDeviceRealData.IsRecord = DeviceStaticData.IsRecord
			tempDeviceRealData.RecordInterval = DeviceStaticData.RecordInterval
			tempDeviceRealData.RecordDataCharge = DeviceStaticData.RecordDataCharge
			tempDeviceRealData.RecordType = DeviceStaticData.RecordType
			tempDeviceRealData.AlarmLevel = DeviceStaticData.AlarmLevel
			tempDeviceRealData.AlarmClearMessage = DeviceStaticData.AlarmMessage
			tempDeviceRealData.AlarmMessage = DeviceStaticData.AlarmClearMessage
			writeDeviceRealData = append(writeDeviceRealData, tempDeviceRealData)
		}
		//触发器
		for _, AlarmTriggerItem := range getAlarmTrigger {

			var insertRealData DeviceRealData
			insertRealData.DeviceName = params.Name
			insertRealData.ProjectUuid = params.ProjectUuid
			insertRealData.Name = AlarmTriggerItem.TriggerName
			insertRealData.Uuid = uuid.New()
			insertRealData.ModelDataUuid = AlarmTriggerItem.Uuid
			insertRealData.Type = 1
			insertRealData.Value = ""
			insertRealData.Muid = AlarmTriggerItem.TriggerDeviceModelUuid
			insertRealData.DeviceUuid = params.Uuid
			insertRealData.DeviceType = params.DeviceType
			insertRealData.IsAlarm = 1
			insertRealData.IsRecord = 0
			insertRealData.RecordInterval = 0
			insertRealData.AlarmLevel = AlarmTriggerItem.TriggerAlarmLevel
			insertRealData.AlarmClearMessage = AlarmTriggerItem.TriggerAlarmHideText
			insertRealData.AlarmMessage = AlarmTriggerItem.TriggerAlarmShowText
			writeDeviceRealData = append(writeDeviceRealData, insertRealData)
		}
		//自定义数据
		for _, CustomDataItem := range getCustomData {

			var insertRealData DeviceRealData
			insertRealData.DeviceName = params.Name
			insertRealData.ProjectUuid = params.ProjectUuid
			insertRealData.Name = CustomDataItem.Name
			insertRealData.Uuid = uuid.New()
			insertRealData.ModelDataUuid = CustomDataItem.Uuid
			insertRealData.Type = CustomDataItem.DataType
			insertRealData.Value = CustomDataItem.DataDefaultValue
			insertRealData.Muid = CustomDataItem.SelectDataModelUUid
			insertRealData.DeviceUuid = params.Uuid
			insertRealData.DeviceType = 4 //自定义数据类型
			insertRealData.IsAlarm = CustomDataItem.DataAlarm
			insertRealData.IsRecord = CustomDataItem.IsRecord
			insertRealData.RecordInterval = CustomDataItem.RecordInterval
			insertRealData.RecordDataCharge = CustomDataItem.RecordDataCharge
			insertRealData.RecordType = CustomDataItem.RecordType
			insertRealData.AlarmLevel = CustomDataItem.AlarmLevel
			insertRealData.AlarmClearMessage = CustomDataItem.AlarmMessage
			insertRealData.AlarmMessage = CustomDataItem.AlarmClearMessage
			writeDeviceRealData = append(writeDeviceRealData, insertRealData)
		}
		//==========================================系统自带数据
		//设备断线告警的ID
		var tempDeviceRealData DeviceRealData
		tempDeviceRealData.Auth = 2
		tempDeviceRealData.ProjectUuid = params.ProjectUuid
		tempDeviceRealData.DeviceName = params.Name
		tempDeviceRealData.Name = "device.DeviceStatus"
		tempDeviceRealData.Uuid = "sys.suid.device.status"
		tempDeviceRealData.ModelDataUuid = "sys.suid.device.status"

		tempDeviceRealData.Type = 1

		tempDeviceRealData.Value = "0"
		tempDeviceRealData.Muid = params.Muid
		tempDeviceRealData.DeviceUuid = params.Uuid
		tempDeviceRealData.DeviceType = params.DeviceType

		tempDeviceRealData.IsAlarm = 1
		tempDeviceRealData.IsRecord = 0
		tempDeviceRealData.RecordInterval = 0
		tempDeviceRealData.AlarmLevel = 4
		tempDeviceRealData.AlarmClearMessage = "device.DeviceStatusOnline"
		tempDeviceRealData.AlarmMessage = "device.DeviceStatusOffline"
		writeDeviceRealData = append(writeDeviceRealData, tempDeviceRealData)
		//==================================================================
		writeResult := Db.Model(&DeviceRealData{}).CreateInBatches(&writeDeviceRealData, 20)

		if writeResult.Error != nil {
			return errmsg.ERROR
		}
	}
	result := Db.Model(&MonitorList{}).Create(&params)

	if result.Error != nil {
		return errmsg.ERROR_DATABASE
	}
	return errmsg.SUCCSECODE
}

/*
*
获取单个设备信息
*/
func GetDeviceData(device_uuid string) (int, MonitorList) {

	var existList MonitorList
	existError := Db.Model(&MonitorList{}).Where("uuid = ?", device_uuid).First(&existList)

	if existError.Error != nil {
		return errmsg.ERROR_DATABASE, existList
	}
	return errmsg.SUCCSECODE, existList
}

/*
*
删除设备或者区域
*/
func DelAllDevices(uuid []string) (int, int) {

	var delSnmpModels MonitorList
	var delReadData DeviceRealData
	// var delAlarmData DevicesAlarmList
	// var delHistoryData DevicesHistoryDataList

	err := Db.Model(&MonitorList{}).Unscoped().Where("uuid in ?", uuid).Delete(&delSnmpModels).Error
	if err != nil {
		return errmsg.ERROR, delSnmpModels.DeviceType
	}
	err1 := Db.Model(&DeviceRealData{}).Unscoped().Where("device_uuid in ?", uuid).Delete(&delReadData).Error
	if err1 != nil {
		return errmsg.ERROR, delSnmpModels.DeviceType
	}

	// err2 := Db.Model(&DevicesAlarmList{}).Unscoped().Where("device_uuid = ?", uuid).Delete(&delAlarmData).Error
	// if err2 != nil {
	// 	return errmsg.ERROR, delSnmpModels.DeviceType
	// }

	// err3 := Db.Model(&DevicesHistoryDataList{}).Unscoped().Where("device_uuid = ?", uuid).Delete(&delHistoryData).Error
	// if err3 != nil {
	// 	return errmsg.ERROR, delSnmpModels.DeviceType
	// }

	return errmsg.SUCCSECODE, delSnmpModels.DeviceType
}

/*
*
删除设备或者区域
*/
func DelDeviceOrZone(uuid string) (int, int) {

	var delSnmpModels MonitorList
	var delReadData DeviceRealData
	var delAlarmData DevicesAlarmList
	// var delHistoryData DevicesHistoryDataList

	err6 := Db.Model(&MonitorList{}).Where("uuid = ?", uuid).First(&delSnmpModels).Error
	if err6 != nil {
		return errmsg.ERROR, -1
	}
	if delSnmpModels.Pid == 0 {
		return errmsg.ROOTZONE_NOT_DEL, delSnmpModels.DeviceType
	}

	err := Db.Model(&MonitorList{}).Unscoped().Where("uuid = ?", uuid).Delete(&delSnmpModels).Error
	if err != nil {
		return errmsg.ERROR, delSnmpModels.DeviceType
	}
	err1 := Db.Model(&DeviceRealData{}).Unscoped().Where("device_uuid = ?", uuid).Delete(&delReadData).Error
	if err1 != nil {
		return errmsg.ERROR, delSnmpModels.DeviceType
	}

	err2 := Db.Model(&DevicesAlarmList{}).Unscoped().Where("device_uuid = ? and clear_time<?", uuid, "2007-01-02 15:04:05").Delete(&delAlarmData).Error
	if err2 != nil {
		return errmsg.ERROR, delSnmpModels.DeviceType
	}

	// err3 := Db.Model(&DevicesHistoryDataList{}).Unscoped().Where("device_uuid = ?", uuid).Delete(&delHistoryData).Error
	// if err3 != nil {
	// 	return errmsg.ERROR, delSnmpModels.DeviceType
	// }

	return errmsg.SUCCSECODE, delSnmpModels.DeviceType
}

/*
*
更新设备或者区域
*/
func UpdateDeviceOrZone(uuid string, updateData MonitorList) int {

	var updateAlarmData DevicesAlarmList
	var updateRealData DeviceRealData
	// var updateHistoryData DevicesHistoryDataList

	err := Db.Model(&MonitorList{}).Where("uuid = ?", uuid).Updates(&updateData).Error
	if err != nil {
		return errmsg.ERROR
	}
	updateRealData.DeviceName = updateData.Name
	err4 := Db.Model(&DeviceRealData{}).Where("device_uuid = ?", uuid).Updates(&updateRealData).Error
	if err4 != nil {
		return errmsg.ERROR
	}

	updateAlarmData.DeviceName = updateData.Name
	err2 := Db.Model(&DevicesAlarmList{}).Where("device_uuid = ?", uuid).Updates(&updateAlarmData).Error
	if err2 != nil {
		return errmsg.ERROR
	}
	// updateHistoryData.DeviceName = updateData.Name
	// err3 := Db.Model(&DevicesHistoryDataList{}).Where("device_uuid = ?", uuid).Updates(&updateHistoryData).Error
	// if err3 != nil {
	// 	return errmsg.ERROR
	// }
	return errmsg.SUCCSECODE
}

/*
*
更新设备启停
*/
func SetDeviceEnable(uuid string, updateData MonitorList) int {

	if updateData.IsEnable == 1 {
		updateData.Status = 0
	} else {
		updateData.Status = 3
	}
	protocol_common.DeviceRealDataMapByUUID.Store("sys.suid.device.status", fmt.Sprintf("%d", updateData.Status))
	protocol_common.DeviceRealDataMap.Store(updateData.Name+"->"+"设备状态", fmt.Sprintf("%d", updateData.Status))
	err := Db.Model(&MonitorList{}).Select("is_enable", "status").Where("uuid = ?", uuid).Updates(&updateData).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSECODE

}

/*
*
获取设备实时数据
*/
func GetRealData(uuid string) ([]DeviceRealData, int) {

	var readDataList []DeviceRealData
	err := Db.Model(&DeviceRealData{}).Where("device_uuid = ?", uuid).Find(&readDataList).Error
	if err != nil {
		return nil, errmsg.ERROR
	}
	return readDataList, errmsg.SUCCSECODE
}

/*
*
获取设备实时数据
*/
func GetRealDataByUuid(uuid, devices []string) ([]DeviceRealData, int) {

	var readDataList []DeviceRealData
	err := Db.Model(&DeviceRealData{}).Select("model_data_uuid,device_uuid,uuid,project_uuid,value").Where("model_data_uuid in ? and device_uuid in ?", uuid, devices).Find(&readDataList).Error
	if err != nil {
		return nil, errmsg.ERROR
	}
	return readDataList, errmsg.SUCCSECODE
}

// 设备模型列表
func DeviceModelData(getType int, devicelist []string, ProjectUuid string) (interface{}, int) {

	var getDeviceModels []DevicesModel
	var getDeviceList []MonitorList
	var result []interface{}

	err := Db.Model(&MonitorList{}).Where("uuid in ?", devicelist).Find(&getDeviceList).Error
	if err == nil && len(getDeviceList) > 0 {
		var DeviceType []int
		var muid []string
		for _, device := range getDeviceList {
			DeviceType = append(DeviceType, device.DeviceType)
			muid = append(muid, device.Muid)
		}
		err1 := Db.Model(&DevicesModel{}).Where("type in ? and project_uuid = ? and uuid in ?", DeviceType, ProjectUuid, muid).Find(&getDeviceModels).Error
		if err1 != nil {
			return result, errmsg.MODEL_HAVED_BAND
		}
	} else {
		err1 := Db.Model(&DevicesModel{}).Where("type != 0 and project_uuid = ?", ProjectUuid).Find(&getDeviceModels).Error

		if err1 != nil {
			return result, errmsg.MODEL_HAVED_BAND
		}
	}
	for _, device := range getDeviceModels {

		single_device := make(map[string]interface{})
		single_device["lable"] = device.Name

		if device.Type == 1 {
			var getMibs []SnmpDevicesDataModel
			if getType == 1 {
				Db.Model(&SnmpDevicesDataModel{}).Where("muid = ? AND is_alarm = 1", device.Uuid).Select("*").Find(&getMibs)
			} else if getType == 2 {
				Db.Model(&SnmpDevicesDataModel{}).Where("muid = ? AND is_record = 1", device.Uuid).Select("*").Find(&getMibs)
			} else {
				Db.Model(&SnmpDevicesDataModel{}).Where("muid = ?", device.Uuid).Select("*").Find(&getMibs)
			}
			single_device["DataList"] = getMibs
		} else if device.Type == 2 {
			var getRegisterAddressList []ModbusDevicesDataModel

			if getType == 1 {
				Db.Model(&ModbusDevicesDataModel{}).Select("*").Where("muid = ? AND is_alarm = 1", device.Uuid).Find(&getRegisterAddressList)
			} else if getType == 2 {
				Db.Model(&ModbusDevicesDataModel{}).Select("*").Where("muid = ? AND is_record = 1", device.Uuid).Find(&getRegisterAddressList)
			} else {
				Db.Model(&ModbusDevicesDataModel{}).Select("*").Where("muid = ?", device.Uuid).Find(&getRegisterAddressList)
			}
			single_device["DataList"] = getRegisterAddressList
		} else if device.Type == 3 {
			var getRegisterAddressList []OpcuaDevicesDataModel

			if getType == 1 {
				Db.Model(&OpcuaDevicesDataModel{}).Select("*").Where("muid = ? AND is_alarm = 1", device.Uuid).Find(&getRegisterAddressList)
			} else if getType == 2 {
				Db.Model(&OpcuaDevicesDataModel{}).Select("*").Where("muid = ? AND is_record = 1", device.Uuid).Find(&getRegisterAddressList)
			} else {
				Db.Model(&OpcuaDevicesDataModel{}).Select("*").Where("muid = ?", device.Uuid).Find(&getRegisterAddressList)
			}
			single_device["DataList"] = getRegisterAddressList
		} else if device.Type == 15 {
			var getRegisterAddressList []SimS7DataModel

			if getType == 1 {
				Db.Model(&SimS7DataModel{}).Select("*").Where("muid = ? AND is_alarm = 1", device.Uuid).Find(&getRegisterAddressList)
			} else if getType == 2 {
				Db.Model(&SimS7DataModel{}).Select("*").Where("muid = ? AND is_record = 1", device.Uuid).Find(&getRegisterAddressList)
			} else {
				Db.Model(&SimS7DataModel{}).Select("*").Where("muid = ?", device.Uuid).Find(&getRegisterAddressList)
			}
			single_device["DataList"] = getRegisterAddressList
		} else if device.Type == 5 {
			var getRegisterAddressList []RESTFulDataModel

			if getType == 1 {
				Db.Model(&RESTFulDataModel{}).Select("*").Where("muid = ? AND is_alarm = 1", device.Uuid).Find(&getRegisterAddressList)
			} else if getType == 2 {
				Db.Model(&RESTFulDataModel{}).Select("*").Where("muid = ? AND is_record = 1", device.Uuid).Find(&getRegisterAddressList)
			} else {
				Db.Model(&RESTFulDataModel{}).Select("*").Where("muid = ?", device.Uuid).Find(&getRegisterAddressList)
			}
			single_device["DataList"] = getRegisterAddressList
		} else if device.Type == 20 {
			var getRegisterAddressList []MqttDevicesDataModel

			if getType == 1 {
				Db.Model(&MqttDevicesDataModel{}).Select("*").Where("muid = ? AND is_alarm = 1", device.Uuid).Find(&getRegisterAddressList)
			} else if getType == 2 {
				Db.Model(&MqttDevicesDataModel{}).Select("*").Where("muid = ? AND is_record = 1", device.Uuid).Find(&getRegisterAddressList)
			} else {
				Db.Model(&MqttDevicesDataModel{}).Select("*").Where("muid = ?", device.Uuid).Find(&getRegisterAddressList)
			}
			single_device["DataList"] = getRegisterAddressList
		} else if device.Type == 30 {
			var getRegisterAddressList []Dlt645DevicesDataModel

			if getType == 1 {
				Db.Model(&Dlt645DevicesDataModel{}).Select("*").Where("muid = ? AND is_alarm = 1", device.Uuid).Find(&getRegisterAddressList)
			} else if getType == 2 {
				Db.Model(&Dlt645DevicesDataModel{}).Select("*").Where("muid = ? AND is_record = 1", device.Uuid).Find(&getRegisterAddressList)
			} else {
				Db.Model(&Dlt645DevicesDataModel{}).Select("*").Where("muid = ?", device.Uuid).Find(&getRegisterAddressList)
			}
			single_device["DataList"] = getRegisterAddressList
		} else if device.Type == 40 {
			var getRegisterAddressList []IEC104DevicesDataModel

			if getType == 1 {
				Db.Model(&IEC104DevicesDataModel{}).Select("*").Where("muid = ? AND is_alarm = 1", device.Uuid).Find(&getRegisterAddressList)
			} else if getType == 2 {
				Db.Model(&IEC104DevicesDataModel{}).Select("*").Where("muid = ? AND is_record = 1", device.Uuid).Find(&getRegisterAddressList)
			} else {
				Db.Model(&IEC104DevicesDataModel{}).Select("*").Where("muid = ?", device.Uuid).Find(&getRegisterAddressList)
			}
			single_device["DataList"] = getRegisterAddressList
		} else if device.Type == 350 {
			var getRegisterAddressList []IEC61850DevicesDataModel

			if getType == 1 {
				Db.Model(&IEC61850DevicesDataModel{}).Select("*").Where("muid = ? AND is_alarm = 1", device.Uuid).Find(&getRegisterAddressList)
			} else if getType == 2 {
				Db.Model(&IEC61850DevicesDataModel{}).Select("*").Where("muid = ? AND is_record = 1", device.Uuid).Find(&getRegisterAddressList)
			} else {
				Db.Model(&IEC61850DevicesDataModel{}).Select("*").Where("muid = ?", device.Uuid).Find(&getRegisterAddressList)
			}
			single_device["DataList"] = getRegisterAddressList
		} else if device.Type == 470 {
			var getRegisterAddressList []HJ212DevicesDataModel

			if getType == 1 {
				Db.Model(&HJ212DevicesDataModel{}).Select("*").Where("muid = ? AND is_alarm = 1", device.Uuid).Find(&getRegisterAddressList)
			} else if getType == 2 {
				Db.Model(&HJ212DevicesDataModel{}).Select("*").Where("muid = ? AND is_record = 1", device.Uuid).Find(&getRegisterAddressList)
			} else {
				Db.Model(&HJ212DevicesDataModel{}).Select("*").Where("muid = ?", device.Uuid).Find(&getRegisterAddressList)
			}
			single_device["DataList"] = getRegisterAddressList
		} else if device.Type == 480 {
			var getRegisterAddressList []VirtualDeviceDataModel

			if getType == 1 {
				Db.Model(&VirtualDeviceDataModel{}).Select("*").Where("muid = ? AND is_alarm = 1", device.Uuid).Find(&getRegisterAddressList)
			} else if getType == 2 {
				Db.Model(&VirtualDeviceDataModel{}).Select("*").Where("muid = ? AND is_record = 1", device.Uuid).Find(&getRegisterAddressList)
			} else {
				Db.Model(&VirtualDeviceDataModel{}).Select("*").Where("muid = ?", device.Uuid).Find(&getRegisterAddressList)
			}
			single_device["DataList"] = getRegisterAddressList
		} else if device.Type == 490 {
			var getRegisterAddressList []CJT188DevicesDataModel

			if getType == 1 {
				Db.Model(&CJT188DevicesDataModel{}).Select("*").Where("muid = ? AND is_alarm = 1", device.Uuid).Find(&getRegisterAddressList)
			} else if getType == 2 {
				Db.Model(&CJT188DevicesDataModel{}).Select("*").Where("muid = ? AND is_record = 1", device.Uuid).Find(&getRegisterAddressList)
			} else {
				Db.Model(&CJT188DevicesDataModel{}).Select("*").Where("muid = ?", device.Uuid).Find(&getRegisterAddressList)
			}
			single_device["DataList"] = getRegisterAddressList
		} else if device.Type == 500 {
			var getRegisterAddressList []BacnetDevicesDataModel

			if getType == 1 {
				Db.Model(&BacnetDevicesDataModel{}).Select("*").Where("muid = ? AND is_alarm = 1", device.Uuid).Find(&getRegisterAddressList)
			} else if getType == 2 {
				Db.Model(&BacnetDevicesDataModel{}).Select("*").Where("muid = ? AND is_record = 1", device.Uuid).Find(&getRegisterAddressList)
			} else {
				Db.Model(&BacnetDevicesDataModel{}).Select("*").Where("muid = ?", device.Uuid).Find(&getRegisterAddressList)
			}
			single_device["DataList"] = getRegisterAddressList
		}

		result = append(result, single_device)
	}

	//自定义数据
	var getCustomDataList []CustomData
	single_CustomData_device := make(map[string]interface{})
	single_CustomData_device["lable"] = "alarm.current.CustomModel"
	if getType == 1 {
		Db.Model(&CustomData{}).Select("*").Where("data_alarm = 1 and project_uuid = ? ", ProjectUuid).Find(&getCustomDataList)
	} else if getType == 2 {
		Db.Model(&CustomData{}).Select("*").Where("is_record = 1 and project_uuid = ? ", ProjectUuid).Find(&getCustomDataList)
	} else {
		Db.Model(&CustomData{}).Select("*").Where("ID > 0 and project_uuid = ? ", ProjectUuid).Find(&getCustomDataList)
	}
	single_CustomData_device["DataList"] = getCustomDataList
	result = append(result, single_CustomData_device)

	//自定义数据
	var getStaticDataList []StaticData
	single_StaticData_device := make(map[string]interface{})
	single_StaticData_device["lable"] = "alarm.current.SystemVarModel"
	if getType == 1 {
		Db.Model(&StaticData{}).Select("*").Where("data_alarm = 1 and project_uuid = ? ", ProjectUuid).Find(&getStaticDataList)
	} else if getType == 2 {
		Db.Model(&StaticData{}).Select("*").Where("is_record = 1 and project_uuid = ? ", ProjectUuid).Find(&getStaticDataList)
	} else {
		Db.Model(&StaticData{}).Select("*").Where("ID > 0 and project_uuid = ? ", ProjectUuid).Find(&getStaticDataList)
	}
	single_StaticData_device["DataList"] = getStaticDataList
	result = append(result, single_StaticData_device)

	if getType == 1 {
		single_system_alarm := make(map[string]interface{})
		single_system_alarm["lable"] = "alarm.current.SystemModel"

		var systemAlarmList []SnmpDevicesDataModel

		var sysAlarm SnmpDevicesDataModel
		sysAlarm.Name = "device.DeviceStatus"
		sysAlarm.Uuid = "sys.suid.device.status"
		systemAlarmList = append(systemAlarmList, sysAlarm)

		single_system_alarm["DataList"] = systemAlarmList
		result = append(result, single_system_alarm)
	}
	if getType == 1 {
		single_trigger := make(map[string]interface{})
		single_trigger["lable"] = "alarm.trigger.TriggerTitle"
		var getMibs []AlarmTrigger
		var pushMibs []SnmpDevicesDataModel

		Db.Model(&AlarmTrigger{}).Where("id != 0 and trigger_type = 2 and project_uuid = ?", ProjectUuid).Select("*").Find(&getMibs)
		if len(getMibs) > 0 {
			for _, trigger := range getMibs {
				var temp SnmpDevicesDataModel
				temp.Name = trigger.TriggerName
				temp.Uuid = trigger.Uuid
				pushMibs = append(pushMibs, temp)
			}
		}
		single_trigger["DataList"] = pushMibs
		result = append(result, single_trigger)
	}

	return result, errmsg.SUCCSECODE
}

// mib获取
func DataModelGet(modelType int, project_uuid string) []SnmpDevicesDataModel {

	var getMibs []SnmpDevicesDataModel
	if modelType == 6 { //静态数据
		Db.Model(&StaticData{}).Where("ID > 0 and project_uuid=?", project_uuid).Select("*").Find(&getMibs)
	} else if modelType == 7 { //自定义数据
		Db.Model(&CustomData{}).Where("ID > 0 and project_uuid=?", project_uuid).Select("*").Find(&getMibs)
	}
	return getMibs
}

// SyncDeviceRealData 根据项目的 monitor_list 设备实例，补建缺失的 device_real_data
// 用于设备已创建但实时数据点未自动生成时的修复场景（常见于批量导入）
func SyncDeviceRealData(projectUuid string) (int, int) {
	var devices []MonitorList
	err := Db.Model(&MonitorList{}).Where("project_uuid = ? AND type >= 1 AND muid != '' AND deleted_at IS NULL", projectUuid).Find(&devices).Error
	if err != nil {
		return 0, 0
	}

	created := 0
	skipped := 0

	for _, dev := range devices {
		// 查询该设备型号的数据模型
		var dataModels []ModbusDevicesDataModel
		err := Db.Model(&ModbusDevicesDataModel{}).Where("muid = ? AND deleted_at IS NULL", dev.Muid).Find(&dataModels).Error
		if err != nil {
			continue
		}

		for _, dm := range dataModels {
			// 检查是否已存在
			var exist DeviceRealData
			existErr := Db.Model(&DeviceRealData{}).Where("device_uuid = ? AND model_data_uuid = ?", dev.Uuid, dm.Uuid).First(&exist).Error
			if !errors.Is(existErr, gorm.ErrRecordNotFound) {
				skipped++
				continue
			}

			auth := 3
			if dm.Auth == "ReadOnly" {
				auth = 1
			} else if dm.Auth == "ReadWrite" {
				auth = 2
			}

			rd := DeviceRealData{
				Uuid:                 uuid.New(),
				ProjectUuid:          projectUuid,
				DeviceName:           dev.Name,
				DeviceUuid:           dev.Uuid,
				Name:                 dm.Name,
				ModelDataUuid:        dm.Uuid,
				Muid:                 dev.Muid,
				DataUnit:             dm.DataUnit,
				Auth:                 auth,
				Type:                 1,
				DeviceType:           2,
				Value:                "",
				IsRecord:             dm.IsRecord,
				RecordInterval:       dm.RecordInterval,
				RecordDataCharge:     dm.RecordDataCharge,
				RecordType:           dm.RecordType,
				IsAlarm:              dm.IsAlarm,
				AlarmLevel:           dm.AlarmLevel,
				AlarmMessage:         dm.AlarmMessage,
				AlarmClearMessage:    dm.AlarmClearMessage,
				ConversionExpression: dm.ConversionExpression,
			}
			err = Db.Create(&rd).Error
			if err != nil {
				continue
			}
			created++
		}
	}
	return created, skipped
}

// BatchDisableAlarm 批量禁用项目下所有实时数据点的告警
func BatchDisableAlarm(projectUuid string) int {
	err := Db.Model(&DeviceRealData{}).Where("project_uuid = ? AND deleted_at IS NULL", projectUuid).
		Updates(map[string]interface{}{"is_alarm": 0, "alarm_shield": 1}).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSECODE
}

// BatchSetDeviceStatus 按 project_uuid 批量设置设备状态
func BatchSetDeviceStatus(projectUuid string, status int) int {
	err := Db.Model(&MonitorList{}).Where("project_uuid = ? AND deleted_at IS NULL AND type >= 1", projectUuid).
		Update("status", status).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSECODE
}
