/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-09-07 14:45:55
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

import (
	"ISMServer/models"
	ismhj212 "ISMServer/protocol/HJ212"
	s7protocols "ISMServer/protocol/S7"
	bacnetprotocols "ISMServer/protocol/bacnet"
	cjt188protocols "ISMServer/protocol/cjt188"
	protocol_common "ISMServer/protocol/common"
	dlt645protocols "ISMServer/protocol/dlt645"
	iec104protocols "ISMServer/protocol/iec104"
	iec61850protocols "ISMServer/protocol/iec61850"
	modbusprotocols "ISMServer/protocol/modbus"
	mqttprotocols "ISMServer/protocol/mqtt"
	ismnode "ISMServer/protocol/netnode"
	opcuaprotocols "ISMServer/protocol/opcua"
	opcuapub "ISMServer/protocol/opcuapub"
	snmpprotocols "ISMServer/protocol/snmp"
	ismWebsocket "ISMServer/protocol/websocket"
	customDataTask "ISMServer/task/DealWithCustomData"
	ISMScriptFunc "ISMServer/task/ISMScript/func"
	staticDataTask "ISMServer/task/staticData"
	triggerAlarmTask "ISMServer/task/triggerAlarm"
	"ISMServer/utils/errmsg"
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/beego/beego/v2/core/config"
	beego "github.com/beego/beego/v2/server/web"
	"gorm.io/gorm"
)

type DeviceLibraryController struct {
	beego.Controller
}

type ICMP struct {
	Type        uint8
	Code        uint8
	CheckSum    uint16
	Identifier  uint16
	SequenceNum uint16
}

func getICMP(seq uint16) ICMP {
	icmp := ICMP{
		Type:        8,
		Code:        0,
		CheckSum:    0,
		Identifier:  0,
		SequenceNum: seq,
	}

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp)
	icmp.CheckSum = CheckSum(buffer.Bytes())
	buffer.Reset()

	return icmp
}

func sendICMPRequest(icmp ICMP, IPAddr string) error {
	conn, err := net.Dial("ip4:icmp", IPAddr)
	if err != nil {
		fmt.Printf("Fail to connect to remote host: %s\n", err)
		return err
	}
	defer conn.Close()

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp)

	if _, err := conn.Write(buffer.Bytes()); err != nil {
		return err
	}

	tStart := time.Now()

	conn.SetReadDeadline((time.Now().Add(time.Second * 2)))

	recv := make([]byte, 1024)
	receiveCnt, err := conn.Read(recv)

	if err != nil {
		return err
	}

	tEnd := time.Now()
	duration := tEnd.Sub(tStart).Nanoseconds() / 1e6

	fmt.Printf("%d bytes from %s: seq=%d time=%dms\n", receiveCnt, IPAddr, icmp.SequenceNum, duration)

	return err
}

func ping(host string) bool {
	var err error

	if err = sendICMPRequest(getICMP(uint16(1)), host); err != nil {
		return false
	}
	return true
}
func CheckSum(data []byte) uint16 {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += (sum >> 16)

	return uint16(^sum)
}

func (c *DeviceLibraryController) Ping() {

	type TreeList struct {
		IP string `json:"IP"`
	}
	var code int64 = 0
	var getParams TreeList
	data := c.Ctx.Input.RequestBody

	result := map[string]interface{}{
		"code": code,
	}

	err := json.Unmarshal(data, &getParams)
	if err != nil {
		code = -1
	} else {
		if !ping(getParams.IP) {
			result["code"] = -2
		} else {
			result["code"] = 0
		}
	}
	c.Data["json"] = result

	c.ServeJSON() //返回json格式

}
func (c *DeviceLibraryController) MonitorTree() {

	var getLists []*models.TreeList
	var code int64 = 0
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		getLists = models.GetAllDevices(0, ProjectUuid, false)
	} else {
		code = -1
		getLists = nil
	}
	if TreeList, OK := ismnode.ProjectDataCollects.Load(ProjectUuid); OK {
		TreeListTemp := TreeList.([]string)
		for _, item := range TreeListTemp {
			if TreeItem, IsExit := protocol_common.ISMNodeProjectDataCollects.Load(item); IsExit {
				PushTreeItem := TreeItem.(*models.TreeList)
				getLists = append(getLists, PushTreeItem)
			}
		}
	}
	result := map[string]interface{}{
		"code": code,
		"list": getLists,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *DeviceLibraryController) AddDeviceOrZone() {

	var getParams models.MonitorList
	var code int = 0

	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到对象中
		err := json.Unmarshal(data, &getParams)
		if err != nil {
			code = errmsg.NOTJSON
		} else {
			getParams.ProjectUuid = ProjectUuid
			code = models.AddDeviceOrZone(getParams)
			if code == 0 {
				go opcuapub.RefreshPublishedNodes()
			}
		}
	} else {
		code = -1
	}

	result := map[string]interface{}{
		"code": code,
	}
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了设备"+getParams.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	c.Data["json"] = result
	if getParams.DeviceType == 1 {
		go snmpprotocols.SnmpCloseChan()
	} else if getParams.DeviceType == 2 {
		go modbusprotocols.ModbusCloseChan()
	} else if getParams.DeviceType == 3 {
		go opcuaprotocols.OpcuaCloseChan()
	} else if getParams.DeviceType == 15 {
		go s7protocols.SimS7CloseChan()
	} else if getParams.DeviceType == 20 {
		go mqttprotocols.MqttCloseChan()
	} else if getParams.DeviceType == 30 {
		go dlt645protocols.DLT645CloseChan()
	} else if getParams.DeviceType == 40 {
		go iec104protocols.IEC104CloseChan()
	} else if getParams.DeviceType == 350 {
		go iec61850protocols.IEC61850CloseChan()
	} else if getParams.DeviceType == 470 {
		go ismhj212.HJ212CloseChan()
	} else if getParams.DeviceType == 490 {
		go cjt188protocols.Cjt188CloseChan()
	} else if getParams.DeviceType == 500 {
		go bacnetprotocols.BacnetCloseChan()
	}
	d := &ismnode.ISMNetNodeClientCtl{}
	go d.NetNodeUpdateDataCollect(ProjectUuid)
	go staticDataTask.PushStaticCloseChan()
	go customDataTask.CustomDataCloseChan()
	c.ServeJSON() //返回json格式
}
func (c *DeviceLibraryController) CopyDevice() {

	type CopyStu struct {
		CopyType       int    `json:"CopyType"`
		CopyUuid       string `json:"CopyUuid"`
		CopyDeviceType int    `json:"CopyDeviceType"`
		CopyCount      int    `json:"CopyCount"`
		Pid            int32  `json:"pid"`
	}
	var getParams models.MonitorList
	var GetCopyParsms CopyStu
	var code int = 0

	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到对象中
		err := json.Unmarshal(data, &GetCopyParsms)
		if err != nil {
			code = errmsg.NOTJSON
		} else {
			r1, DeviceData := models.GetDeviceData(GetCopyParsms.CopyUuid)
			if r1 == errmsg.SUCCSECODE {
				for i := 0; i < GetCopyParsms.CopyCount; i++ {
					getParams = DeviceData
					getParams.ID = 0
					getParams.Name = DeviceData.Name + "_copy_" + strconv.Itoa(i+1)
					code = models.AddDeviceOrZone(getParams)
				}
			} else {
				code = r1
			}
		}
	} else {
		code = -1
	}

	result := map[string]interface{}{
		"code": code,
	}
	c.Data["json"] = result
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "添加了设备"+getParams.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	c.ServeJSON() //返回json格式
}

func (c *DeviceLibraryController) DelDeviceOrZone() {

	var paramsJson map[string]interface{}

	var code int = 0
	var DeviceType int = 0

	data := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(data, &paramsJson)
	if err != nil {
		code = errmsg.NOTJSON
	} else {
		uuid := fmt.Sprintf("%s", paramsJson["uuid"])
		code, DeviceType = models.DelDeviceOrZone(uuid)
		if code == 0 {
			go opcuapub.RefreshPublishedNodes()
		}
	}

	result := map[string]interface{}{
		"code": code,
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	c.Data["json"] = result
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了设备", errmsg.JournalLevelInfo, c.Ctx.Input)
	go staticDataTask.PushStaticCloseChan()

	if DeviceType == 1 {
		go snmpprotocols.SnmpCloseChan()
	} else if DeviceType == 2 {
		go modbusprotocols.ModbusCloseChan()
	} else if DeviceType == 3 {
		go opcuaprotocols.OpcuaCloseChan()
	} else if DeviceType == 15 {
		go s7protocols.SimS7CloseChan()
	} else if DeviceType == 20 {
		go mqttprotocols.MqttCloseChan()
	} else if DeviceType == 30 {
		go dlt645protocols.DLT645CloseChan()
	} else if DeviceType == 40 {
		go iec104protocols.IEC104CloseChan()
	} else if DeviceType == 350 {
		go iec61850protocols.IEC61850CloseChan()
	} else if DeviceType == 470 {
		go ismhj212.HJ212CloseChan()
	} else if DeviceType == 490 {
		go cjt188protocols.Cjt188CloseChan()
	} else if DeviceType == 500 {
		go bacnetprotocols.BacnetCloseChan()
	}
	d := &ismnode.ISMNetNodeClientCtl{}
	go d.NetNodeUpdateDataCollect(ProjectUuid)
	go triggerAlarmTask.AlarmTriggerCloseChan()
	go customDataTask.CustomDataCloseChan()
	c.ServeJSON() //返回json格式
}
func (c *DeviceLibraryController) DelAllDevice() {

	type DelAllList struct {
		Uuid []string `json:"uuid"`
	}
	var paramsJson DelAllList

	var code int = 0
	var DeviceType int = 0

	data := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(data, &paramsJson)
	if err != nil {
		code = errmsg.NOTJSON
	} else {
		code, DeviceType = models.DelAllDevices(paramsJson.Uuid)
		if code == 0 {
			go opcuapub.RefreshPublishedNodes()
		}
	}

	result := map[string]interface{}{
		"code": code,
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	c.Data["json"] = result
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除了设备", errmsg.JournalLevelInfo, c.Ctx.Input)
	go staticDataTask.PushStaticCloseChan()

	if DeviceType == 1 {
		go snmpprotocols.SnmpCloseChan()
	} else if DeviceType == 2 {
		go modbusprotocols.ModbusCloseChan()
	} else if DeviceType == 3 {
		go opcuaprotocols.OpcuaCloseChan()
	} else if DeviceType == 15 {
		go s7protocols.SimS7CloseChan()
	} else if DeviceType == 20 {
		go mqttprotocols.MqttCloseChan()
	} else if DeviceType == 30 {
		go dlt645protocols.DLT645CloseChan()
	} else if DeviceType == 40 {
		go iec104protocols.IEC104CloseChan()
	} else if DeviceType == 350 {
		go iec61850protocols.IEC61850CloseChan()
	} else if DeviceType == 470 {
		go ismhj212.HJ212CloseChan()
	} else if DeviceType == 490 {
		go cjt188protocols.Cjt188CloseChan()
	} else if DeviceType == 500 {
		go bacnetprotocols.BacnetCloseChan()
	}
	d := &ismnode.ISMNetNodeClientCtl{}
	go d.NetNodeUpdateDataCollect(ProjectUuid)
	go triggerAlarmTask.AlarmTriggerCloseChan()
	go customDataTask.CustomDataCloseChan()
	c.ServeJSON() //返回json格式
}
func (c *DeviceLibraryController) EditDeviceOrZone() {

	type updateStu struct {
		Uuid       string             `json:"uuid"`
		UpdateData models.MonitorList `json:"editData"`
	}
	var getParams updateStu
	var code int = 0

	data := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(data, &getParams)
	if err != nil {
		code = errmsg.NOTJSON
	} else {

		code = models.UpdateDeviceOrZone(getParams.Uuid, getParams.UpdateData)
		if code == 0 {
			go opcuapub.RefreshPublishedNodes()
		}
	}

	result := map[string]interface{}{
		"code": code,
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了设备"+getParams.UpdateData.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	c.Data["json"] = result

	if getParams.UpdateData.DeviceType == 1 {
		go snmpprotocols.SnmpCloseChan()
	} else if getParams.UpdateData.DeviceType == 2 {
		go modbusprotocols.ModbusCloseChan()
	} else if getParams.UpdateData.DeviceType == 3 {
		go opcuaprotocols.OpcuaCloseChan()
	} else if getParams.UpdateData.DeviceType == 15 {
		go s7protocols.SimS7CloseChan()
	} else if getParams.UpdateData.DeviceType == 20 {
		go mqttprotocols.MqttCloseChan()
	} else if getParams.UpdateData.DeviceType == 30 {
		go dlt645protocols.DLT645CloseChan()
	} else if getParams.UpdateData.DeviceType == 40 {
		go iec104protocols.IEC104CloseChan()
	} else if getParams.UpdateData.DeviceType == 350 {
		go iec61850protocols.IEC61850CloseChan()
	} else if getParams.UpdateData.DeviceType == 470 {
		go ismhj212.HJ212CloseChan()
	} else if getParams.UpdateData.DeviceType == 490 {
		go cjt188protocols.Cjt188CloseChan()
	} else if getParams.UpdateData.DeviceType == 500 {
		go bacnetprotocols.BacnetCloseChan()
	}
	d := &ismnode.ISMNetNodeClientCtl{}
	go d.NetNodeUpdateDataCollect(ProjectUuid)
	go customDataTask.CustomDataCloseChan()
	c.ServeJSON() //返回json格式
}
func (c *DeviceLibraryController) SetDeviceStartOrStop() {

	type updateStu struct {
		Uuid       string             `json:"uuid"`
		UpdateData models.MonitorList `json:"editData"`
	}
	var getParams updateStu
	var code int = 0

	data := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(data, &getParams)
	if err != nil {
		code = errmsg.NOTJSON
	} else {

		code = models.SetDeviceEnable(getParams.Uuid, getParams.UpdateData)
		if code == 0 {
			go opcuapub.RefreshPublishedNodes()
		}
	}

	result := map[string]interface{}{
		"code": code,
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "编辑了设备"+getParams.UpdateData.Name, errmsg.JournalLevelInfo, c.Ctx.Input)
	c.Data["json"] = result

	if getParams.UpdateData.DeviceType == 1 {
		go snmpprotocols.SnmpCloseChan()
	} else if getParams.UpdateData.DeviceType == 2 {
		go modbusprotocols.ModbusCloseChan()
	} else if getParams.UpdateData.DeviceType == 3 {
		go opcuaprotocols.OpcuaCloseChan()
	} else if getParams.UpdateData.DeviceType == 15 {
		go s7protocols.SimS7CloseChan()
	} else if getParams.UpdateData.DeviceType == 20 {
		go mqttprotocols.MqttCloseChan()
	} else if getParams.UpdateData.DeviceType == 30 {
		go dlt645protocols.DLT645CloseChan()
	} else if getParams.UpdateData.DeviceType == 40 {
		go iec104protocols.IEC104CloseChan()
	} else if getParams.UpdateData.DeviceType == 350 {
		go iec61850protocols.IEC61850CloseChan()
	} else if getParams.UpdateData.DeviceType == 470 {
		go ismhj212.HJ212CloseChan()
	} else if getParams.UpdateData.DeviceType == 490 {
		go cjt188protocols.Cjt188CloseChan()
	} else if getParams.UpdateData.DeviceType == 500 {
		go bacnetprotocols.BacnetCloseChan()
	}
	d := &ismnode.ISMNetNodeClientCtl{}
	go d.NetNodeUpdateDataCollect(ProjectUuid)
	go customDataTask.CustomDataCloseChan()
	c.ServeJSON() //返回json格式
}
func (c *DeviceLibraryController) GetRealData() {

	var code int = 0
	var IsRemoveGW bool
	var getParams = make(map[string]interface{})
	var realData []models.DeviceRealData

	data := c.Ctx.Input.RequestBody
	SProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if SProjectUuid == "" {
		result := map[string]interface{}{
			"code":     -6,
			"realData": nil,
		}
		c.Data["json"] = result

		c.ServeJSON() //返回json格式
		return
	}
	//json数据封装到对象中
	err := json.Unmarshal(data, &getParams)
	if err != nil {
		code = errmsg.NOTJSON
	} else {
		uuid := fmt.Sprintf("%s", getParams["uuid"])
		IsRemoveGW = getParams["IsRemoveGW"].(bool)
		// ProjectUuid := getParams["ProjectUuid"].(string)
		if !IsRemoveGW {
			realData, code = models.GetRealData(uuid)
			dbtype, _ := config.Int("dbtype")
			if dbtype == 1 {
				for key, v := range realData {
					DeviceDataValue, isExist := protocol_common.DeviceRealDataMapByUUID.Load(v.Uuid)
					if isExist {
						realData[key].Value = DeviceDataValue.(string)
					}
				}
			}
		} else {
			if GetGwRealData, OK := ismnode.ProjectDeviceRealDataCollects.Load(uuid); OK {
				realData, OK = GetGwRealData.([]models.DeviceRealData)
				if !OK {
					code = -4
				}
			} else {
				code = -6
			}
		}
	}
	result := map[string]interface{}{
		"code":     code,
		"realData": realData,
	}

	c.Ctx.Output.Header("Content-Encoding", "gzip")
	gw := gzip.NewWriter(c.Ctx.ResponseWriter)
	defer gw.Close()

	err1 := json.NewEncoder(gw).Encode(result)
	if err1 != nil {
		fmt.Println("encode err: ", err1)
	}
}
func (c *DeviceLibraryController) GetRealDataByUuid() {

	var code int = 0
	type RecvUuid struct {
		Uuid        []string `json:"uuid"`
		Devices     []string `json:"devices"`
		DevicesName []string `json:"DevicesName"`
		DataName    []string `json:"DataName"`
	}
	var getParams RecvUuid

	var realData []models.DeviceRealData

	data := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(data, &getParams)
	if err != nil {
		code = errmsg.NOTJSON
	} else {
		realData, code = models.GetRealDataByUuid(getParams.Uuid, getParams.Devices)

		for key, v := range realData {
			DeviceDataValue, isExist := protocol_common.DeviceRealDataMapByUUID.Load(v.Uuid)
			if isExist {
				realData[key].Value = DeviceDataValue.(string)
			}
		}

		for _, v := range getParams.Uuid {
			if systemdata, IsExist := protocol_common.ISMSystemDataMapByUUID.Load(v); IsExist {
				realData = append(realData, systemdata.(models.DeviceRealData))
			}
			if tempdata, IsTrue := protocol_common.ISMNodeDeviceRealDataMapByUUID.Load(v); IsTrue {
				realData = append(realData, tempdata.(models.DeviceRealData))
			}
		}
	}
	result := map[string]interface{}{
		"code":     code,
		"realData": realData,
	}
	c.Ctx.Output.Header("Content-Encoding", "gzip")
	gw := gzip.NewWriter(c.Ctx.ResponseWriter)
	defer gw.Close()

	err1 := json.NewEncoder(gw).Encode(result)
	if err1 != nil {
		fmt.Println("encode err: ", err1)
	}
}

func (c *DeviceLibraryController) SetRealData() {

	var code int = 0

	var setParams = make(map[string]interface{})

	SProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if SProjectUuid == "" {
		result := map[string]interface{}{
			"code":     -16,
			"realData": nil,
		}
		c.Data["json"] = result

		c.ServeJSON() //返回json格式
		return
	}
	result := map[string]interface{}{
		"code": code,
	}
	data := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(data, &setParams)
	if err != nil {
		code = errmsg.NOTJSON
	} else {
		DeviceUuid := setParams["deviceUuid"].(string)
		DataUuid := setParams["dataUuid"].(string)
		SetValue := setParams["value"].(string)

		if RealDataList, IsTrue := ismnode.ProjectDeviceRealDataCollects.Load(DeviceUuid); IsTrue {
			RealData, OK := RealDataList.([]models.DeviceRealData)
			if OK {
				ProjectUuid := RealData[0].ProjectUuid
				d := &ismnode.ISMNetNodeServerHandlerCtl{}
				result = d.NetNodeSendApi(SProjectUuid, ProjectUuid, "SetRealData", setParams)
			} else {
				result["code"] = -7
			}

		} else {

			var staticData models.StaticData
			err1 := models.Db.Model(&models.StaticData{}).Where("uuid = ?", DataUuid).First(&staticData)

			var readData models.DeviceRealData
			err := models.Db.Model(&models.DeviceRealData{}).Where("model_data_uuid = ? and device_uuid = ?", DataUuid, DeviceUuid).First(&readData).Error
			if err != nil {
				code = errmsg.ERROR
			} else {
				if !errors.Is(err1.Error, gorm.ErrRecordNotFound) {
					var tempPushData protocol_common.PushRealDataWebData
					err := models.Db.Model(&models.DeviceRealData{}).Where("model_data_uuid = ? and device_uuid = ?", DataUuid, DeviceUuid).Update("value", SetValue).Error
					if err != nil {
						code = errmsg.ERROR
					}
					err2 := models.Db.Model(&models.StaticData{}).Where("uuid = ?", DataUuid).Update("data_default_value", SetValue).Error
					if err2 != nil {
						code = errmsg.ERROR
					}
					staticDataTask.PushStaticCloseChan()
					protocol_common.DeviceRealDataMapByUUID.Store(readData.Uuid, SetValue)
					protocol_common.DeviceRealDataMap.Store(readData.DeviceName+"->"+staticData.Name, SetValue)

					tempPushData.DeviceUuid = DeviceUuid
					tempPushData.ProjectUuid = readData.ProjectUuid

					tempPushData.Cmd = "RealData"

					var signleAlarm protocol_common.PushAlarm
					var signleHistoryData models.DevicesHistoryDataList

					signleAlarm.DeviceUuid = readData.DeviceUuid
					signleAlarm.ProjectUuid = readData.ProjectUuid
					signleAlarm.DataUuid = readData.Uuid
					signleAlarm.ModelDataUuid = readData.ModelDataUuid
					signleAlarm.AlarmLevel = readData.AlarmLevel

					signleHistoryData.DeviceUuid = readData.DeviceUuid
					signleHistoryData.ProjectUuid = readData.ProjectUuid
					signleHistoryData.DataUuid = readData.Uuid
					signleHistoryData.ModelDataUuid = readData.ModelDataUuid
					signleHistoryData.DataUnit = readData.DataUnit
					signleHistoryData.RecordInterval = readData.RecordInterval

					tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: readData.Uuid, ModelDataUuid: readData.ModelDataUuid, Value: SetValue})
					// protocol_common.GGatherDataQueue.QueuePush(tempPushData)
					go ismWebsocket.WSSend(tempPushData, tempPushData.ProjectUuid, 2)
					//设备主动告警信息
					if readData.IsAlarm == 1 && readData.AlarmShield == 0 {
						signleAlarm.Value = SetValue
						if signleAlarm.Value == "true" {
							signleAlarm.Value = "1"
						} else if signleAlarm.Value == "false" {
							signleAlarm.Value = "0"
						} else {
							value, err := strconv.ParseFloat(signleAlarm.Value, 32)
							if err == nil {
								if value >= 1 {
									signleAlarm.Value = "1"
								} else {
									signleAlarm.Value = "0"
								}
							} else {
								signleAlarm.Value = "0"
							}
						}
						signleAlarm.AlarmLevel = readData.AlarmLevel
						signleAlarm.AlarmClearMessage = readData.AlarmClearMessage
						signleAlarm.AlarmMessage = readData.AlarmMessage
						signleAlarm.DataName = readData.Name
						signleAlarm.DeviceName = readData.DeviceName
						signleAlarm.HappenTime = time.Now()
						protocol_common.GAlarmQueue.QueuePush(signleAlarm)
					} else if readData.IsRecord == 1 {
						//存储信息
						signleHistoryData.DataValue = SetValue
						signleHistoryData.DataName = readData.Name
						signleHistoryData.DeviceName = readData.DeviceName
						signleHistoryData.RecordTime = time.Now()
						signleHistoryData.RecordType = readData.RecordType
						signleHistoryData.RecordDataCharge = readData.RecordDataCharge
						protocol_common.GHistoryDataQueue.QueuePush(signleHistoryData)
					}

				} else {
					if readData.Auth == 1 {
						code = errmsg.ERROR
					} else {
						ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
						WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "设置了数据"+readData.DeviceName+"的"+readData.Name+"旧值:"+readData.Value+",新值:"+SetValue, errmsg.JournalLevelInfo, c.Ctx.Input)
						if readData.DeviceType == 1 { //SNMP
							snmpSetObj := &snmpprotocols.SnmpCtl{}
							code = snmpSetObj.SnmpSet(readData.Uuid, SetValue)
						} else if readData.DeviceType == 2 { //MODBUS
							modbusSetObj := &modbusprotocols.ModbusCtl{}
							code = modbusSetObj.ModbusSetData(readData.Uuid, SetValue)
						} else if readData.DeviceType == 3 { //OPCUA
							opcuaSetObj := &opcuaprotocols.OpcuaCtl{}
							code = opcuaSetObj.OPcuaDeviceSetData(readData.Uuid, SetValue)
						} else if readData.DeviceType == 15 { //西门子S7
							S7SetObj := &s7protocols.SimS7Ctl{}
							code = S7SetObj.SimS7DeviceSetData(readData.Uuid, SetValue)
						} else if readData.DeviceType == 20 { //Mqtt
							code = mqttprotocols.MqttSetPubData(readData.Uuid, SetValue)
						} else if readData.DeviceType == 40 { //IEC104设备
							IEC104SetObj := &iec104protocols.IEC1045Ctl{}
							code = IEC104SetObj.IEC104SetData(readData.Uuid, SetValue)
						} else if readData.DeviceType == 350 { //IEC61850设备
							IEC104SetObj := &iec61850protocols.IEC61850Ctl{}
							code = IEC104SetObj.IEC61850DeviceSetData(readData.Uuid, SetValue)
						} else if readData.DeviceType == 480 { //虚拟设备
							err := models.Db.Model(&models.DeviceRealData{}).Where("model_data_uuid = ? and device_uuid = ?", DataUuid, DeviceUuid).Update("value", SetValue).Error
							if err != nil {
								code = errmsg.ERROR
							} else {
								code = errmsg.SUCCSECODE
							}
						} else if readData.DeviceType == 500 { //BACnet
							BACnetSetObj := &bacnetprotocols.BacnetCtl{}
							code = BACnetSetObj.BACnetSetData(readData.Uuid, SetValue)
						}
						if code == 0 {
							protocol_common.DeviceRealDataMapByUUID.Store(readData.Uuid, SetValue)
							protocol_common.DeviceRealDataMap.Store(readData.DeviceName+"->"+readData.Name, SetValue)

							var tempPushData protocol_common.PushRealDataWebData
							tempPushData.DeviceUuid = DeviceUuid
							tempPushData.ProjectUuid = readData.ProjectUuid

							tempPushData.Cmd = "RealData"

							var signleAlarm protocol_common.PushAlarm
							var signleHistoryData models.DevicesHistoryDataList
							var pushTriggerAlarm protocol_common.TriggerRealData
							//触发器告警信息
							pushTriggerAlarm.DeviceUuid = readData.DeviceUuid
							pushTriggerAlarm.ProjectUuid = readData.ProjectUuid
							pushTriggerAlarm.DataUuid = readData.Uuid
							pushTriggerAlarm.DataName = readData.Name
							pushTriggerAlarm.DeviceName = readData.DeviceName
							pushTriggerAlarm.DataType = 1
							pushTriggerAlarm.AlarmShield = readData.AlarmShield
							pushTriggerAlarm.GatherTime = time.Now()

							pushTriggerAlarm.ModelDataUuid = readData.ModelDataUuid

							signleAlarm.DeviceUuid = readData.DeviceUuid
							signleAlarm.ProjectUuid = readData.ProjectUuid
							signleAlarm.DataUuid = readData.Uuid
							signleAlarm.ModelDataUuid = readData.ModelDataUuid
							signleAlarm.AlarmLevel = readData.AlarmLevel

							signleHistoryData.DeviceUuid = readData.DeviceUuid
							signleHistoryData.ProjectUuid = readData.ProjectUuid
							signleHistoryData.DataUuid = readData.Uuid
							signleHistoryData.ModelDataUuid = readData.ModelDataUuid
							signleHistoryData.DataUnit = readData.DataUnit
							signleHistoryData.RecordInterval = readData.RecordInterval

							pushTriggerAlarm.Value = SetValue
							tempPushData.Data = append(tempPushData.Data, protocol_common.UpdateStu{Uuid: readData.Uuid, ModelDataUuid: readData.ModelDataUuid, Value: SetValue})
							// protocol_common.GGatherDataQueue.QueuePush(tempPushData)
							go ismWebsocket.WSSend(tempPushData, tempPushData.ProjectUuid, 2)
							//设备主动告警信息
							if readData.IsAlarm == 1 && readData.AlarmShield == 0 {
								signleAlarm.Value = SetValue
								if signleAlarm.Value == "true" {
									signleAlarm.Value = "1"
								} else if signleAlarm.Value == "false" {
									signleAlarm.Value = "0"
								} else {
									value, err := strconv.ParseFloat(signleAlarm.Value, 32)
									if err == nil {
										if value >= 1 {
											signleAlarm.Value = "1"
										} else {
											signleAlarm.Value = "0"
										}
									} else {
										signleAlarm.Value = "0"
									}
								}
								signleAlarm.AlarmLevel = readData.AlarmLevel
								signleAlarm.AlarmClearMessage = readData.AlarmClearMessage
								signleAlarm.AlarmMessage = readData.AlarmMessage
								signleAlarm.DataName = readData.Name
								signleAlarm.DeviceName = readData.DeviceName
								signleAlarm.HappenTime = time.Now()
								protocol_common.GAlarmQueue.QueuePush(signleAlarm)
							} else if readData.IsRecord == 1 {
								//存储信息
								signleHistoryData.DataValue = SetValue
								signleHistoryData.DataName = readData.Name
								signleHistoryData.DeviceName = readData.DeviceName
								signleHistoryData.RecordTime = time.Now()
								signleHistoryData.RecordType = readData.RecordType
								signleHistoryData.RecordDataCharge = readData.RecordDataCharge
								protocol_common.GHistoryDataQueue.QueuePush(signleHistoryData)
							}
							//触发器队列
							_, isExist := protocol_common.DeviceAlarmTriggerMap.Load(pushTriggerAlarm.ModelDataUuid)
							if isExist {
								protocol_common.GTriggerDataQueue.Store(pushTriggerAlarm.DeviceUuid+"->"+pushTriggerAlarm.ModelDataUuid, pushTriggerAlarm)
							}
							//自定义数据队列
							_, isExistCustom := protocol_common.DeviceCustomDataMap.Load(pushTriggerAlarm.ModelDataUuid)
							if isExistCustom {
								protocol_common.GCustomDataQueue.Store(pushTriggerAlarm.DeviceUuid+"->"+pushTriggerAlarm.ModelDataUuid, pushTriggerAlarm)
							}
						}
					}

				}

			}
			result["code"] = code
		}
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
func (c *DeviceLibraryController) GetSupportDeviceType() {

	getSupputDevice := models.GetSupportDeviceList()

	result := map[string]interface{}{
		"list": getSupputDevice,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式

}
func (c *DeviceLibraryController) GetDeviceModelDataList() {

	var data interface{}

	var code int
	type RecvQueryParams struct {
		Type         int      `json:"getType"`
		SelectDevice []string `json:"SelectDevice"`
	}
	var params RecvQueryParams
	rawData := c.Ctx.Input.RequestBody

	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到user对象中
		err := json.Unmarshal(rawData, &params)
		if err != nil {
			code = -1
		} else {
			data, code = models.DeviceModelData(params.Type, params.SelectDevice, ProjectUuid)
		}
	} else {
		code = -1
		data = nil
	}

	result := map[string]interface{}{
		"code": code,
		"list": data,
	}

	c.Ctx.Output.Header("Content-Encoding", "gzip")
	gw := gzip.NewWriter(c.Ctx.ResponseWriter)
	defer gw.Close()

	err1 := json.NewEncoder(gw).Encode(&result)
	if err1 != nil {
		fmt.Println("encode err: ", err1)
	}
}

func (c *DeviceLibraryController) GetDataModelData() {
	var code int64 = 0
	var getModelByType = struct {
		DataModelType int `json:"type"`
	}{1}

	var getLists []models.SnmpDevicesDataModel

	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		//json数据封装到user对象中
		err := json.Unmarshal(data, &getModelByType)
		if err != nil {
			code = -1
		} else {
			getLists = models.DataModelGet(getModelByType.DataModelType, ProjectUuid)
		}
	} else {
		code = -1
		getLists = nil
	}
	result := map[string]interface{}{
		"code": code,
		"list": getLists,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
func (c *DeviceLibraryController) GetRealDataByName() {

	var code int = 0
	type QueryParams struct {
		DeviceList []string `json:"DeviceList"`
		DataList   []string `json:"DataList"`
	}
	var totalData []any

	var getParams QueryParams
	data := c.Ctx.Input.RequestBody
	SProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if SProjectUuid == "" {
		result := map[string]interface{}{
			"code":     -6,
			"realData": nil,
		}
		c.Data["json"] = result

		c.ServeJSON() //返回json格式
		return
	}
	//json数据封装到对象中
	err := json.Unmarshal(data, &getParams)
	if err != nil {
		code = errmsg.NOTJSON
	} else {

		// 3. 校验参数有效性（设备和数据列表长度需一致）
		if len(getParams.DeviceList) == 0 || len(getParams.DataList) == 0 {
			code = -7 // 参数为空错误码
		} else {

			// 4. 逐个获取数据（核心逻辑：单条调用，不批量）
			for i := 0; i < len(getParams.DeviceList); i++ {
				var singleData []any
				for k := 0; k < len(getParams.DataList); k++ {
					deviceName := getParams.DeviceList[i]
					dataName := getParams.DataList[k]

					// 跳过空值
					if deviceName == "" || dataName == "" {
						singleData = append(singleData, '-')
					}

					// 拼接成 "设备名称->数据名称" 格式
					deviceDataKey := fmt.Sprintf("%s->%s", deviceName, dataName)

					// 逐个调用GetDeviceData（每次仅传一个key）
					singleDataTemp := ISMScriptFunc.GetDeviceData(deviceDataKey)
					if singleDataTemp == nil {
						singleData = append(singleData, '-')
					} else {
						singleData = append(singleData, singleDataTemp)
					}
				}
				totalData = append(totalData, singleData)
			}
		}
	}
	result := map[string]interface{}{
		"code":     code,
		"realData": totalData,
	}

	c.Ctx.Output.Header("Content-Encoding", "gzip")
	gw := gzip.NewWriter(c.Ctx.ResponseWriter)
	defer gw.Close()

	err1 := json.NewEncoder(gw).Encode(result)
	if err1 != nil {
		fmt.Println("encode err: ", err1)
	}
}

// SyncDeviceRealData 补建缺失的 device_real_data（根据已有设备实例）
func (c *DeviceLibraryController) SyncDeviceRealData() {
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	var code int

	if ProjectUuid == "" {
		code = -1
	} else {
		created, skipped := models.SyncDeviceRealData(ProjectUuid)
		code = 0

		result := map[string]interface{}{
			"code":    code,
			"created": created,
			"skipped": skipped,
		}
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid,
			fmt.Sprintf("补建实时数据：创建%d条, 跳过%d条", created, skipped),
			errmsg.JournalLevelInfo, c.Ctx.Input)
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	result := map[string]interface{}{
		"code": code,
	}
	c.Data["json"] = result
	c.ServeJSON()
}

// BatchDisableAlarm 批量禁用项目下所有实时数据点的告警
func (c *DeviceLibraryController) BatchDisableAlarm() {
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	var code int
	if ProjectUuid == "" {
		code = -1
	} else {
		code = models.BatchDisableAlarm(ProjectUuid)
	}
	result := map[string]interface{}{"code": code}
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid,
		"批量禁用告警", errmsg.JournalLevelInfo, c.Ctx.Input)
	c.Data["json"] = result
	c.ServeJSON()
}

// BatchSetDeviceStatus 批量设置项目下设备状态
func (c *DeviceLibraryController) BatchSetDeviceStatus() {
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	var code int
	if ProjectUuid == "" {
		code = -1
		c.Data["json"] = map[string]interface{}{"code": code}
		c.ServeJSON()
		return
	}

	data := c.Ctx.Input.RequestBody
	var params struct {
		Status int `json:"status"`
	}
	if err := json.Unmarshal(data, &params); err != nil {
		c.Data["json"] = map[string]interface{}{"code": errmsg.NOTJSON}
		c.ServeJSON()
		return
	}

	code = models.BatchSetDeviceStatus(ProjectUuid, params.Status)
	result := map[string]interface{}{"code": code}
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid,
		fmt.Sprintf("批量设置设备状态: %d", params.Status),
		errmsg.JournalLevelInfo, c.Ctx.Input)
	c.Data["json"] = result
	c.ServeJSON()
}
