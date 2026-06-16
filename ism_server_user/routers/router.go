/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:26
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-08-21 09:42:42
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package routers

import (
	"ISMServer/controllers"
	"ISMServer/middleware"
	"ISMServer/models"
	DataInterface "ISMServer/protocol/DataInterface"
	"ISMServer/utils/errmsg"
	"errors"
	"net/http"
	"strings"

	"github.com/beego/beego/v2/core/config"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"gorm.io/gorm"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.UserController{}, "*:Login")
	beego.Router("/routes", &controllers.UserController{}, "*:Routes")
	//snmp
	beego.Router("/snmpmodeladd", &controllers.SnmpDeviceModelController{}, "*:ModelAdd")
	beego.Router("/snmpmodellist", &controllers.SnmpDeviceModelController{}, "*:ModelList")
	beego.Router("/snmpmodeldel", &controllers.SnmpDeviceModelController{}, "*:ModelDel")
	beego.Router("/snmpmodelget", &controllers.SnmpDeviceModelController{}, "*:ModelGet")
	beego.Router("/snmpmodeledit", &controllers.SnmpDeviceModelController{}, "*:ModelEdit")
	beego.Router("/snmpmodelimport", &controllers.SnmpDeviceModelController{}, "*:ModelImport")
	beego.Router("/snmpmodelsavemib", &controllers.SnmpDeviceModelController{}, "*:ModelSaveMib")
	beego.Router("/snmpmodelgetmibs", &controllers.SnmpDeviceModelController{}, "*:ModelGetMibs")
	beego.Router("/snmpmodeldeletemibs", &controllers.SnmpDeviceModelController{}, "*:ModelDeleteMibs")
	beego.Router("/modelDataEdit", &controllers.SnmpDeviceModelController{}, "*:ModelDataEdit")
	beego.Router("/snmpmodelimportxml", &controllers.SnmpDeviceModelController{}, "*:ModelImportXml")

	beego.Router("/snmpmodelgethistorymibs", &controllers.SnmpDeviceModelController{}, "*:ModelGetHistoryMibs")
	beego.Router("/supportDevice", &controllers.DeviceLibraryController{}, "*:GetSupportDeviceType")

	beego.Router("/monitortree", &controllers.DeviceLibraryController{}, "*:MonitorTree")
	beego.Router("/monitorAdd", &controllers.DeviceLibraryController{}, "*:AddDeviceOrZone")
	beego.Router("/monitorCopy", &controllers.DeviceLibraryController{}, "*:CopyDevice")
	beego.Router("/monitorEdit", &controllers.DeviceLibraryController{}, "*:EditDeviceOrZone")
	beego.Router("/monitorDel", &controllers.DeviceLibraryController{}, "*:DelDeviceOrZone")
	beego.Router("/syncDeviceRealData", &controllers.DeviceLibraryController{}, "*:SyncDeviceRealData")
	beego.Router("/DeviceRealDataDisableAlarm", &controllers.DeviceLibraryController{}, "*:BatchDisableAlarm")
	beego.Router("/MonitorBatchSetStatus", &controllers.DeviceLibraryController{}, "*:BatchSetDeviceStatus")
	beego.Router("/monitorAllDel", &controllers.DeviceLibraryController{}, "*:DelAllDevice")
	beego.Router("/ping", &controllers.DeviceLibraryController{}, "*:Ping")
	beego.Router("/getRealData", &controllers.DeviceLibraryController{}, "*:GetRealData")
	beego.Router("/getRealDataByUuid", &controllers.DeviceLibraryController{}, "*:GetRealDataByUuid")
	beego.Router("/setData", &controllers.DeviceLibraryController{}, "*:SetRealData")
	beego.Router("/SetDeviceStartOrStop", &controllers.DeviceLibraryController{}, "*:SetDeviceStartOrStop")

	beego.Router("/getDataModelData", &controllers.DeviceLibraryController{}, "*:GetDataModelData")

	//显示模型
	beego.Router("/displayModelAdd", &controllers.DisplayModelController{}, "*:ModelAdd")
	beego.Router("/displayModelList", &controllers.DisplayModelController{}, "*:ModelList")
	beego.Router("/displayModelDel", &controllers.DisplayModelController{}, "*:ModelDel")
	beego.Router("/displayModelGet", &controllers.DisplayModelController{}, "*:ModelGet")
	beego.Router("/displayModelEdit", &controllers.DisplayModelController{}, "*:ModelEdit")
	beego.Router("/getDisplayModelLayerData", &controllers.DisplayModelController{}, "*:ModelLayerGet")
	beego.Router("/getDisplayModelLayerDataByToken", &controllers.DisplayModelController{}, "*:ModelLayerGetByToken")
	beego.Router("/saveDisplayModelLayerData", &controllers.DisplayModelController{}, "*:ModelLayerSave")
	beego.Router("/DisplayModelPageAdd", &controllers.DisplayModelController{}, "*:ModelPageAdd")
	beego.Router("/DisplayModelPageDel", &controllers.DisplayModelController{}, "*:ModelPageDel")
	beego.Router("/DisplayModelPageEdit", &controllers.DisplayModelController{}, "*:ModelPageEdit")
	beego.Router("/DisplayModelPageSetHome", &controllers.DisplayModelController{}, "*:ModelPageSetHome")
	beego.Router("/GetUserDisplayList", &controllers.DisplayModelController{}, "*:GetUserModelList")
	beego.Router("/displayImageUpload/:suuid", &controllers.DisplayModelController{}, "*:DisplayImageUpload")
	beego.Router("/DisplayModelPageCopy", &controllers.DisplayModelController{}, "*:ModelLayerCopy")
	beego.Router("/DisplayTempleteList", &controllers.DisplayModelController{}, "*:DisplayTempleteList")
	beego.Router("/DisplayTempleteGet", &controllers.DisplayModelController{}, "*:DisplayTempleteGet")
	beego.Router("/DisplayLoginPage", &controllers.DisplayModelController{}, "*:ModelGetLoginPage")
	beego.Router("/ModelAddUser", &controllers.DisplayModelController{}, "*:ModelAddUser")
	beego.Router("/ModelDelUser", &controllers.DisplayModelController{}, "*:ModelDelUser")
	beego.Router("/GetModelUsers", &controllers.DisplayModelController{}, "*:GetModelUsers")
	beego.Router("/getDisplayModelPagerLayerData", &controllers.DisplayModelController{}, "*:ModelLayerPagerGet")
	//图片管理
	beego.Router("/systemImageList", &controllers.SystemImageController{}, "*:ImageList")
	beego.Router("/systemImageUpload", &controllers.SystemImageController{}, "*:ImageUpload")
	beego.Router("/systemImageDel", &controllers.SystemImageController{}, "*:ImageDel")

	//设备模型
	beego.Router("/GetDeviceModelDataList", &controllers.DeviceLibraryController{}, "*:GetDeviceModelDataList")

	//modbus
	beego.Router("/Comlist", &controllers.ModbusDeviceModelController{}, "*:Comlist")
	beego.Router("/modbusModelAdd", &controllers.ModbusDeviceModelController{}, "*:ModelAdd")
	beego.Router("/modbusModelList", &controllers.ModbusDeviceModelController{}, "*:ModelList")
	beego.Router("/modbusModelDel", &controllers.ModbusDeviceModelController{}, "*:ModelDel")
	beego.Router("/modbusModelEdit", &controllers.ModbusDeviceModelController{}, "*:ModelEdit")
	beego.Router("/modbusModelRegisterGroupAdd", &controllers.ModbusDeviceModelController{}, "*:ModelRegisterGroupAdd")
	beego.Router("/modbusModelRegisterGroupEdit", &controllers.ModbusDeviceModelController{}, "*:ModelRegisterGroupEdit")
	beego.Router("/modbusModelRegisterGroupDel", &controllers.ModbusDeviceModelController{}, "*:ModelRegisterGroupDel")
	beego.Router("/modbusModelRegisterGroupList", &controllers.ModbusDeviceModelController{}, "*:ModelRegisterGroupList")
	beego.Router("/modbusModelRegisterList", &controllers.ModbusDeviceModelController{}, "*:ModelRegisterList")
	beego.Router("/modbusModelRegisterEdit", &controllers.ModbusDeviceModelController{}, "*:ModelRegisterEdit")
	beego.Router("/modbusModelRegisterAdd", &controllers.ModbusDeviceModelController{}, "*:ModelRegisterAdd")
	beego.Router("/modbusModelRegisterDel", &controllers.ModbusDeviceModelController{}, "*:ModelRegisterDel")

	//DLT645
	beego.Router("/dlt645ModelAdd", &controllers.DLT645DeviceModelController{}, "*:ModelAdd")
	beego.Router("/dlt64ModelList", &controllers.DLT645DeviceModelController{}, "*:ModelList")
	beego.Router("/dlt64ModelDel", &controllers.DLT645DeviceModelController{}, "*:ModelDel")
	beego.Router("/dlt64ModelEdit", &controllers.DLT645DeviceModelController{}, "*:ModelEdit")
	beego.Router("/dlt64ModelNodeIDAdd", &controllers.DLT645DeviceModelController{}, "*:ModelDataAdd")
	beego.Router("/dlt64ModelNodeIDDel", &controllers.DLT645DeviceModelController{}, "*:ModelDataDel")
	beego.Router("/dlt64ModelNodeIDEdit", &controllers.DLT645DeviceModelController{}, "*:ModelDataEdit")
	beego.Router("/dlt64ModelNodeIDList", &controllers.DLT645DeviceModelController{}, "*:ModelDataList")

	//CJT188
	beego.Router("/cjt188ModelAdd", &controllers.CJT188DeviceModelController{}, "*:ModelAdd")
	beego.Router("/cjt188ModelList", &controllers.CJT188DeviceModelController{}, "*:ModelList")
	beego.Router("/cjt188ModelDel", &controllers.CJT188DeviceModelController{}, "*:ModelDel")
	beego.Router("/cjt188ModelEdit", &controllers.CJT188DeviceModelController{}, "*:ModelEdit")
	beego.Router("/cjt188ModelNodeIDAdd", &controllers.CJT188DeviceModelController{}, "*:ModelDataAdd")
	beego.Router("/cjt188ModelNodeIDDel", &controllers.CJT188DeviceModelController{}, "*:ModelDataDel")
	beego.Router("/cjt188ModelNodeIDEdit", &controllers.CJT188DeviceModelController{}, "*:ModelDataEdit")
	beego.Router("/cjt188ModelNodeIDList", &controllers.CJT188DeviceModelController{}, "*:ModelDataList")

	//IEC104
	beego.Router("/iec104ModelAdd", &controllers.IEC104ModelController{}, "*:ModelAdd")
	beego.Router("/iec104ModelList", &controllers.IEC104ModelController{}, "*:ModelList")
	beego.Router("/iec104ModelDel", &controllers.IEC104ModelController{}, "*:ModelDel")
	beego.Router("/iec104ModelEdit", &controllers.IEC104ModelController{}, "*:ModelEdit")
	beego.Router("/iec104ModelNodeIDAdd", &controllers.IEC104ModelController{}, "*:ModelDataAdd")
	beego.Router("/iec104ModelNodeIDDel", &controllers.IEC104ModelController{}, "*:ModelDataDel")
	beego.Router("/iec104ModelNodeIDEdit", &controllers.IEC104ModelController{}, "*:ModelDataEdit")
	beego.Router("/iec104ModelNodeIDList", &controllers.IEC104ModelController{}, "*:ModelDataList")

	//OPCUA
	beego.Router("/opcuaModelAdd", &controllers.OPCUAController{}, "*:ModelAdd")
	beego.Router("/opcuaModelList", &controllers.OPCUAController{}, "*:ModelList")
	beego.Router("/opcuaModelDel", &controllers.OPCUAController{}, "*:ModelDel")
	beego.Router("/opcuaModelEdit", &controllers.OPCUAController{}, "*:ModelEdit")

	beego.Router("/opcuaNodeidImport/:muid", &controllers.OPCUAController{}, "*:NodeIDImport")

	beego.Router("/OpcuaModelNodeIDAdd", &controllers.OPCUAController{}, "*:ModelDataAdd")
	beego.Router("/OpcuaModelNodeIDDel", &controllers.OPCUAController{}, "*:ModelDataDel")
	beego.Router("/OpcuaModelNodeIDEdit", &controllers.OPCUAController{}, "*:ModelDataEdit")
	beego.Router("/OpcuaModelNodeIDList", &controllers.OPCUAController{}, "*:ModelDataList")

// 	//IEC61850
// 	beego.Router("/IEC61850ModelAdd", &controllers.IEC61850Controller{}, "*:ModelAdd")
// 	beego.Router("/IEC61850ModelList", &controllers.IEC61850Controller{}, "*:ModelList")
// 	beego.Router("/IEC61850ModelDel", &controllers.IEC61850Controller{}, "*:ModelDel")
// 	beego.Router("/IEC61850ModelEdit", &controllers.IEC61850Controller{}, "*:ModelEdit")
// 	beego.Router("/IEC61850NodeidImport/:muid", &controllers.IEC61850Controller{}, "*:NodeIDImport")
// 	beego.Router("/IEC61850ModelNodeIDAdd", &controllers.IEC61850Controller{}, "*:ModelDataAdd")
// 	beego.Router("/IEC61850ModelNodeIDDel", &controllers.IEC61850Controller{}, "*:ModelDataDel")
// 	beego.Router("/IEC61850ModelNodeIDEdit", &controllers.IEC61850Controller{}, "*:ModelDataEdit")
// 	beego.Router("/IEC61850ModelNodeIDList", &controllers.IEC61850Controller{}, "*:ModelDataList")

	//MQTT
	beego.Router("/mqttModelAdd", &controllers.MqttController{}, "*:ModelAdd")
	beego.Router("/mqttModelList", &controllers.MqttController{}, "*:ModelList")
	beego.Router("/mqttModelDel", &controllers.MqttController{}, "*:ModelDel")
	beego.Router("/mqttModelEdit", &controllers.MqttController{}, "*:ModelEdit")

	beego.Router("/mqttModelNodeIDAdd", &controllers.MqttController{}, "*:ModelDataAdd")
	beego.Router("/mqttModelNodeIDDel", &controllers.MqttController{}, "*:ModelDataDel")
	beego.Router("/mqttModelNodeIDEdit", &controllers.MqttController{}, "*:ModelDataEdit")
	beego.Router("/mqttModelNodeIDList", &controllers.MqttController{}, "*:ModelDataList")

	//HJ212
	beego.Router("/hj212ModelAdd", &controllers.HJ212Controller{}, "*:ModelAdd")
	beego.Router("/hj212ModelList", &controllers.HJ212Controller{}, "*:ModelList")
	beego.Router("/hj212ModelDel", &controllers.HJ212Controller{}, "*:ModelDel")
	beego.Router("/hj212ModelEdit", &controllers.HJ212Controller{}, "*:ModelEdit")

	beego.Router("/hj212ModelNodeIDAdd", &controllers.HJ212Controller{}, "*:ModelDataAdd")
	beego.Router("/hj212ModelNodeIDDel", &controllers.HJ212Controller{}, "*:ModelDataDel")
	beego.Router("/hj212ModelNodeIDEdit", &controllers.HJ212Controller{}, "*:ModelDataEdit")
	beego.Router("/hj212ModelNodeIDList", &controllers.HJ212Controller{}, "*:ModelDataList")
	beego.Router("/hj212ModelDataTemplete", &controllers.HJ212Controller{}, "*:ModelDataTemplete")

	//视频
	beego.Router("/getVideoStatus", &controllers.VideoController{}, "*:GetVideoStatus")
	beego.Router("/addVideo", &controllers.VideoController{}, "*:AddRTSPStream")
	beego.Router("/getVideoList", &controllers.VideoController{}, "*:GetAllRTSPStreamList")
	beego.Router("/delVideo", &controllers.VideoController{}, "*:DelRTSPStream")
	beego.Router("/editVideo", &controllers.VideoController{}, "*:EditRTSPStream")
	beego.Router("/setVideoStopOrStart", &controllers.VideoController{}, "*:SetVideoStopOrStart")
	beego.Router("/livestream", &controllers.VideoController{}, "*:Livestream")
	beego.Router("/webrtcstream/:suuid", &controllers.VideoController{}, "*:WEBRTCLivestream")
	beego.Router("/ws/livestream/?:suuid", &controllers.VideoController{}, "*:WSLivestream")
	beego.Router("/codec/:uuid", &controllers.VideoController{}, "*:Codec")
	beego.Router("/GetMonibucaVideoList", &controllers.VideoController{}, "*:GetMonibucaVideoList")
	beego.Router("/GetMonibucaHistoryVideoList", &controllers.VideoController{}, "*:MonibucaHistoryVideoList")

	beego.Router("/SnapToJpg", &controllers.VideoController{}, "*:SnapRTSPStreamToJpg")
	beego.Router("/GetSnapTree", &controllers.VideoController{}, "*:GetSnapTree")
	beego.Router("/GetSnapVideoTree", &controllers.VideoController{}, "*:GetSnapVideoTree")
	beego.Router("/PlayVideo/?:path/", &controllers.VideoController{}, "*:PlayVideo")
	beego.Router("/GetRecordVideoTree", &controllers.VideoController{}, "*:GetRecordVideoTree")

	beego.Router("/PtzControl", &controllers.VideoController{}, "*:PtzControl")

	//授权页面
	beego.Router("/checkLicense", &controllers.AuthLicenseController{}, "*:GetLicense")
	beego.Router("/saveLicense", &controllers.AuthLicenseController{}, "*:SaveLicense")

	//告警报告
	beego.Router("/GetAlarmHistoryList", &controllers.ReportController{}, "*:GetAlarmHistoryList")
	beego.Router("/GetDataHistoryList", &controllers.ReportController{}, "*:GetDataHistoryList")
	beego.Router("/GetDiyDataHistoryList", &controllers.ReportController{}, "*:GetDiyDataHistoryList")
	//实时告警
	beego.Router("/GetCurrentAlarmList", &controllers.AlarmController{}, "*:GetCurrentAlarmList")
	beego.Router("/AlarmOpt", &controllers.AlarmController{}, "*:AlarmOpt")
	beego.Router("/AlarmShieldList", &controllers.AlarmController{}, "*:GetCurrentShieldAlarmList")
	//告警联动
	beego.Router("/AlarmTriggerAdd", &controllers.AlarmController{}, "*:AlarmTriggerAdd")
	beego.Router("/AlarmTriggerDel", &controllers.AlarmController{}, "*:AlarmTriggerDel")
	beego.Router("/AlarmTriggerEdit", &controllers.AlarmController{}, "*:AlarmTriggerEdit")
	beego.Router("/GetAlarmTriggerList", &controllers.AlarmController{}, "*:GetAlarmTriggerList")

	//用户
	beego.Router("/uploadUserAvatar", &controllers.UserController{}, "*:UploadAvatar")
	beego.Router("/GetUserInfo", &controllers.UserController{}, "*:GetUserInfo")
	beego.Router("/SetUserInfo", &controllers.UserController{}, "*:SetUserInfo")
	beego.Router("/SetUserPassword", &controllers.UserController{}, "*:SetUserPassword")
	beego.Router("/SystemDisplayUserList", &controllers.UserController{}, "*:SystemDisplayUserList")

	beego.Router("/SystemUserAdd", &controllers.UserController{}, "*:SystemUserAdd")
	beego.Router("/SystemUserDel", &controllers.UserController{}, "*:SystemUserDel")
	beego.Router("/SystemUserList", &controllers.UserController{}, "*:SystemUserList")
	//系统
	beego.Router("/SystemRolesList", &controllers.ISMSystem{}, "*:RolesList")
	beego.Router("/GetSystemFonts", &controllers.ISMSystem{}, "*:GetSystemFonts")
	beego.Router("/SetDebug", &controllers.ISMSystem{}, "*:SetDebug")
	beego.Router("/GetSystemAnalysis", &controllers.ISMSystem{}, "*:GetSystemAnalysis")
	beego.Router("/GetSystemNetworkInfo", &controllers.ISMSystem{}, "*:GetSystemNetworkInfo")
	beego.Router("/SaveSystemNetworkInfo", &controllers.ISMSystem{}, "*:SaveSystemNetworkInfo")
	beego.Router("/RebootSystem", &controllers.ISMSystem{}, "*:RebootSystem")
	beego.Router("/RebootISMSystem", &controllers.ISMSystem{}, "*:RebootISMSystem")

	//项目
	beego.Router("/ProjectAdd", &controllers.ProjectController{}, "*:ProjectAdd")
	beego.Router("/ProjectDel", &controllers.ProjectController{}, "*:ProjectDel")
	beego.Router("/ProjectEdit", &controllers.ProjectController{}, "*:ProjectEdit")
	beego.Router("/ProjectFixCreator", &controllers.ProjectController{}, "*:FixProjectCreator")
	beego.Router("/ProjectList", &controllers.ProjectController{}, "*:ProjectList")
	beego.Router("/ExportProject", &controllers.ProjectController{}, "*:ExportProject")
	beego.Router("/ImportProject", &controllers.ProjectController{}, "*:ImportProject")

	//静态数据模型
	beego.Router("/AddStaticData", &controllers.StaticDataController{}, "*:AddStaticData")
	beego.Router("/EditStaticData", &controllers.StaticDataController{}, "*:EditStaticData")
	beego.Router("/DelStaticData", &controllers.StaticDataController{}, "*:DelStaticData")
	beego.Router("/GetStaticDataList", &controllers.StaticDataController{}, "*:StaticDataList")

	//自定义数据模型
	beego.Router("/AddCustomData", &controllers.CustomDataController{}, "*:AddCustomData")
	beego.Router("/EditCustomData", &controllers.CustomDataController{}, "*:EditCustomData")
	beego.Router("/DelCustomData", &controllers.CustomDataController{}, "*:DelCustomData")
	beego.Router("/GetCustomDataList", &controllers.CustomDataController{}, "*:GetCustomDataList")

	//日志
	beego.Router("/GetJournal", &controllers.ISMJournal{}, "*:GetJournalList")

	//告警通知
	beego.Router("/GetAlarmNotice", &controllers.AlarmNoticeController{}, "*:GetAlarmNoticeParams")
	beego.Router("/UpdateAlarmNotice", &controllers.AlarmNoticeController{}, "*:UpdateAlarmNoticeParams")
	beego.Router("/TestSend", &controllers.AlarmNoticeController{}, "*:TestSend")

	//数据库操作
	beego.Router("/GetTablesList", &controllers.DbOptController{}, "*:GetTablesList")
	beego.Router("/DbBackUp", &controllers.DbOptController{}, "*:DbBackUp")
	beego.Router("/GetBackUpList", &controllers.DbOptController{}, "*:GetBackUpList")
	beego.Router("/DbRestore", &controllers.DbOptController{}, "*:DbRestore")
	beego.Router("/GetDbConfig", &controllers.DbOptController{}, "*:GetDbConfig")
	beego.Router("/SetDbConfig", &controllers.DbOptController{}, "*:SetDbConfig")
	beego.Router("/DbDown", &controllers.DbOptController{}, "*:DbDown")

	//系统数据
	beego.Router("/GetSystemData", &controllers.ISMSystem{}, "*:GetSystemData")
	beego.Router("/GetCustomPel", &controllers.ISMSystem{}, "*:GetCustomPel")
	beego.Router("/GetSystemDeviceInfo", &controllers.ISMSystem{}, "*:GetSystemDeviceInfo")
	beego.Router("/GetSystemParams", &controllers.ISMSystem{}, "*:GetSystemParams")
	beego.Router("/GetSystemParamsList", &controllers.ISMSystem{}, "*:GetSystemParamsList")
	beego.Router("/GetPhysicalIDCheck", &controllers.ISMSystem{}, "*:GetPhysicalIDCheck")
	beego.Router("/WitePhysicalID", &controllers.ISMSystem{}, "*:WitePhysicalID")
	beego.Router("/GetAuthLicenseInfo", &controllers.ISMSystem{}, "*:GetAuthLicenseInfo")

	//RESTFul设备采集
	beego.Router("/api/v1/PushDeviceData", &controllers.DeviceRestFulController{}, "*:PushDeviceData")

	beego.Router("/AddRESTFulData", &controllers.DeviceRestFulController{}, "*:AddRESTFulData")
	beego.Router("/EditRESTFulData", &controllers.DeviceRestFulController{}, "*:EditRESTFulData")
	beego.Router("/DelRESTFulData", &controllers.DeviceRestFulController{}, "*:DelRESTFulData")
	beego.Router("/GetRESTFulDataList", &controllers.DeviceRestFulController{}, "*:RESTFulDataList")

	beego.Router("/AddRESTFulModel", &controllers.DeviceRestFulController{}, "*:AddRESTFulModel")
	beego.Router("/EditRESTFulModel", &controllers.DeviceRestFulController{}, "*:EditRESTFulModel")
	beego.Router("/DelRESTFulModel", &controllers.DeviceRestFulController{}, "*:DelRESTFulModel")
	beego.Router("/GetRESTFulModelList", &controllers.DeviceRestFulController{}, "*:RESTFulModelList")
	//虚拟设备
	beego.Router("/AddVirtualDeviceData", &controllers.VirtualDeviceController{}, "*:AddVirtualDeviceData")
	beego.Router("/EditVirtualDeviceData", &controllers.VirtualDeviceController{}, "*:EditVirtualDeviceData")
	beego.Router("/DelVirtualDeviceData", &controllers.VirtualDeviceController{}, "*:DelVirtualDeviceData")
	beego.Router("/GetVirtualDeviceDataList", &controllers.VirtualDeviceController{}, "*:VirtualDeviceDataList")

	beego.Router("/AddVirtualDeviceModel", &controllers.VirtualDeviceController{}, "*:AddVirtualDeviceModel")
	beego.Router("/EditVirtualDeviceModel", &controllers.VirtualDeviceController{}, "*:EditVirtualDeviceModel")
	beego.Router("/DelVirtualDeviceModel", &controllers.VirtualDeviceController{}, "*:DelVirtualDeviceModel")
	beego.Router("/GetVirtualDeviceModelList", &controllers.VirtualDeviceController{}, "*:VirtualDeviceModelList")

	//S7
	beego.Router("/AddS7Data", &controllers.SimS7Controller{}, "*:AddS7Data")
	beego.Router("/EditS7Data", &controllers.SimS7Controller{}, "*:EditS7Data")
	beego.Router("/DelS7Data", &controllers.SimS7Controller{}, "*:DelS7Data")
	beego.Router("/GetS7DataList", &controllers.SimS7Controller{}, "*:S7DataList")

	beego.Router("/AddS7Model", &controllers.SimS7Controller{}, "*:AddS7Model")
	beego.Router("/EditS7Model", &controllers.SimS7Controller{}, "*:EditS7Model")
	beego.Router("/DelS7Model", &controllers.SimS7Controller{}, "*:DelS7Model")
	beego.Router("/GetS7ModelList", &controllers.SimS7Controller{}, "*:S7ModelList")

	//Token
	beego.Router("/GetAccessTokenList", &controllers.UserController{}, "*:GetAccessTokenList")
	beego.Router("/CreateAccessToken", &controllers.UserController{}, "*:CreateAccessToken")
	beego.Router("/DelAccessToken", &controllers.UserController{}, "*:DelAccessToken")
	//升级
	beego.Router("/OnlineCheckUpgrade", &controllers.ISMSystem{}, "*:OnlineCheckUpgrade")
	beego.Router("/LocalUpgrade", &controllers.ISMSystem{}, "*:LocalUpgrade")
	beego.Router("/BeginUpgrade", &controllers.ISMSystem{}, "*:BeginUpgrade")

	beego.Router("/AuthUpload", &controllers.ISMSystem{}, "*:AuthUpload")
	beego.Router("/BackupUpload", &controllers.ISMSystem{}, "*:BackupUpload")

	beego.Router("/DiyUpload", &controllers.ISMSystem{}, "*:DiyUpload")

	//任务计划
	beego.Router("/AddTaskPlan", &controllers.TaskPlanController{}, "*:AddTaskPlan")
	beego.Router("/EditTaskPlan", &controllers.TaskPlanController{}, "*:EditTaskPlan")
	beego.Router("/DelTaskPlan", &controllers.TaskPlanController{}, "*:DelTaskPlan")
	beego.Router("/GetTaskPlanList", &controllers.TaskPlanController{}, "*:GetTaskPlanList")

	//报表模板
	beego.Router("/AddReportTemplete", &controllers.ReportTempleteController{}, "*:AddReportTemplete")
	beego.Router("/DelReportTemplete", &controllers.ReportTempleteController{}, "*:DelReportTemplete")
	beego.Router("/GetReportTemplete", &controllers.ReportTempleteController{}, "*:GetReportTemplete")
	beego.Router("/EditReportTemplete", &controllers.ReportTempleteController{}, "*:EditReportTemplete")
	beego.Router("/SaveReportTemplete", &controllers.ReportTempleteController{}, "*:SaveReportTemplete")
	beego.Router("/HandExport", &controllers.ReportTempleteController{}, "*:HandExport")

	//
	beego.Router("/UpdateDataModel/:muid", &controllers.ISMSystem{}, "*:UpdateDataModel")

	//系统脚本
	beego.Router("/AddScript", &controllers.ISMScriptController{}, "*:AddScript")
	beego.Router("/DelScript", &controllers.ISMScriptController{}, "*:DelScript")
	beego.Router("/EditScript", &controllers.ISMScriptController{}, "*:EditScript")
	beego.Router("/GetScriptList", &controllers.ISMScriptController{}, "*:GetScriptList")
	beego.Router("/CheckScript", &controllers.ISMScriptController{}, "*:CheckScript")
	beego.Router("/ExecSysScript", &controllers.ISMSystem{}, "*:ExecSysScript")
	beego.Router("/DisableScript", &controllers.ISMScriptController{}, "*:DisableScript")

	//内存数据
	beego.Router("/GetMemData", &controllers.MemDbController{}, "*:GetMemData")
	beego.Router("/SetMemData", &controllers.MemDbController{}, "*:SetMemData")

	//ISM组网
	beego.Router("/AddOutConnect", &controllers.ISMNetworkController{}, "*:AddOutConnect")
	beego.Router("/EditOutConnect", &controllers.ISMNetworkController{}, "*:EditOutConnect")
	beego.Router("/DelOutConnect", &controllers.ISMNetworkController{}, "*:DelOutConnect")
	beego.Router("/OptOutConnect", &controllers.ISMNetworkController{}, "*:OptOutConnect")
	beego.Router("/GetConnectOut", &controllers.ISMNetworkController{}, "*:GetConnectOut")
	beego.Router("/GetInConnectList", &controllers.ISMNetworkController{}, "*:GetInConnectList")

	//
	beego.Router("/GetNodeConfig", &controllers.ISMNetworkController{}, "*:GetNodeConfig")
	beego.Router("/SetNodeConfig", &controllers.ISMNetworkController{}, "*:SetNodeConfig")
	//
	beego.Router("/GetChartDataHistoryList", &controllers.ReportController{}, "*:GetChartDataHistoryList")
	beego.Router("/GetTrendChartData", &controllers.ReportController{}, "*:GetTrendChartData")
	beego.Router("/GetTrendChartDataByDate", &controllers.ReportController{}, "*:GetTrendChartDataByDate")
	beego.Router("/GetDataHistoryReport", &controllers.ReportController{}, "*:GetDataHistoryReport")
	beego.Router("/GetHistoryHour", &controllers.ReportController{}, "*:GetHistoryHour")
	beego.Router("/GetHistoryDayDifferenceReport", &controllers.ReportController{}, "*:GetHistoryDayDifferenceReport")
	beego.Router("/GetHistoryWeeklyDifferenceReport", &controllers.ReportController{}, "*:GetHistoryWeeklyDifferenceReport")
	beego.Router("/GetHistoryMonthDifferenceReport", &controllers.ReportController{}, "*:GetHistoryMonthDifferenceReport")
	beego.Router("/GetHistoryYearDifferenceReport", &controllers.ReportController{}, "*:GetHistoryYearDifferenceReport")

	//系统参数接口
	beego.Router("/GetSystemWebData", &controllers.SystemParamsController{}, "*:GetSystemWebData")
	beego.Router("/SaveSystemWebData", &controllers.SystemParamsController{}, "*:SaveSystemWebData")
	beego.Router("/GetSystemMqttData", &controllers.SystemParamsController{}, "*:GetSystemMqttData")
	beego.Router("/SaveSystemMqttData", &controllers.SystemParamsController{}, "*:SaveSystemMqttData")
	beego.Router("/GetSystemModbusData", &controllers.SystemParamsController{}, "*:GetSystemModbusData")
	beego.Router("/SaveSystemModbusData", &controllers.SystemParamsController{}, "*:SaveSystemModbusData")
	beego.Router("/GetSystemHistoryConfig", &controllers.SystemParamsController{}, "*:GetSystemHistoryConfig")
	beego.Router("/SaveSystemHistoryConfig", &controllers.SystemParamsController{}, "*:SaveSystemHistoryConfig")
	beego.Router("/SaveSystemTimeConfig", &controllers.SystemParamsController{}, "*:SaveSystemTimeConfig")
	beego.Router("/TestNtpServer", &controllers.SystemParamsController{}, "*:TestNtpServer")
	//系统数据模版
	beego.Router("/GetSystemDataTemplete", &controllers.SystemDataTempleteController{}, "*:GetSystemDataTemplete")
	beego.Router("/AddSystemDataTemplete", &controllers.SystemDataTempleteController{}, "*:AddSystemDataTemplete")
	beego.Router("/EditSystemDataTemplete", &controllers.SystemDataTempleteController{}, "*:EditSystemDataTemplete")
	beego.Router("/DelSystemDataTemplete", &controllers.SystemDataTempleteController{}, "*:DelSystemDataTemplete")
	//系统数据接口列表
	beego.Router("/GetSystemDataInterfaceList", &controllers.SystemDataInterfaceController{}, "*:GetSystemDataInterfaceList")
	beego.Router("/AddSystemDataInterface", &controllers.SystemDataInterfaceController{}, "*:AddSystemDataInterface")
	beego.Router("/EditSystemDataInterface", &controllers.SystemDataInterfaceController{}, "*:EditSystemDataInterface")
	beego.Router("/DelSystemDataInterface", &controllers.SystemDataInterfaceController{}, "*:DelSystemDataInterface")
	beego.Router("/SystemExecSQLQuery", &controllers.SystemDataInterfaceController{}, "*:SystemExecSQLQuery")
	//IEC104 数据推送接口
	beego.Router("/IEC104PushDataAdd", &controllers.IEC104DataPushController{}, "*:ModelDataAdd")
	beego.Router("/IEC104PushDataDel", &controllers.IEC104DataPushController{}, "*:ModelDataDel")
	beego.Router("/IEC104PushDataEdit", &controllers.IEC104DataPushController{}, "*:ModelDataEdit")
	beego.Router("/IEC104PushDataList", &controllers.IEC104DataPushController{}, "*:ModelDataList")
	beego.Router("/UpdateIEC104DataModel/:muid", &controllers.IEC104DataPushController{}, "*:UpdateIEC104DataModel")

	//IEC104 数据推送接口
	beego.Router("/ModbusTcpPushDataAdd", &controllers.ModbusTcpDataPushController{}, "*:ModelDataAdd")
	beego.Router("/ModbusTcpPushDataDel", &controllers.ModbusTcpDataPushController{}, "*:ModelDataDel")
	beego.Router("/ModbusTcpPushDataEdit", &controllers.ModbusTcpDataPushController{}, "*:ModelDataEdit")
	beego.Router("/ModbusTcpPushDataList", &controllers.ModbusTcpDataPushController{}, "*:ModelDataList")
	beego.Router("/UpdateModbusTcpPushDataModel/:muid", &controllers.ModbusTcpDataPushController{}, "*:UpdateModbusTcpDataModel")

	// ===== AI 项目生成器 OpenAPI =====
// 	beego.Router("/api/v1/open/status", &controllers.OpenApiController{}, "*:Status")
// 	beego.Router("/api/v1/open/projects", &controllers.OpenApiController{}, "*:GetProjects")
// 	beego.Router("/api/v1/open/project/import", &controllers.OpenApiController{}, "*:ImportProject")
// 	beego.Router("/api/v1/open/project/import-file", &controllers.OpenApiController{}, "*:ImportProjectFile")
// 	beego.Router("/api/v1/open/excel/preview", &controllers.OpenApiController{}, "*:ExcelPreview")
// 	// ===== end AI OpenAPI =====

	// 验证用户是否已经登录
	beego.InsertFilter("/*", beego.BeforeExec, FilterUser)

	CorsEnable, err := config.Bool("CorsEnable")
	if err == nil && CorsEnable {
		beego.InsertFilter("/*", beego.BeforeRouter, corsFunc)
		beego.InsertFilter("/*", beego.BeforeStatic, corsFunc)
	}
}

var success = []byte("SUPPORT OPTIONS")

var corsFunc = func(ctx *context.Context) {
	origin := ctx.Input.Header("Origin")
	ctx.Output.Header("Access-Control-Allow-Methods", "OPTIONS,DELETE,POST,GET,PUT,PATCH")
	ctx.Output.Header("Access-Control-Max-Age", "3600")
	ctx.Output.Header("Access-Control-Allow-Headers", "X-Custom-Header,accept,Content-Type,Access-Token,Origin,Authorization,Access-Control-Allow-Headers,projectuuid ")
	ctx.Output.Header("Access-Control-Allow-Credentials", "true")
	ctx.Output.Header("Access-Control-Allow-Origin", origin)
	if ctx.Input.Method() == http.MethodOptions {
		// options请求，返回200
		ctx.Output.SetStatus(http.StatusOK)
		_ = ctx.Output.Body(success)
	}
}

func findCustomAPI(url string) bool {
	for _, v := range DataInterface.GetUrlInterfaceData {
		if v.InterfaceUrl == url {
			return true
		}
	}
	return false
}

var FilterUser = func(ctx *context.Context) {
	ShareAppToken := ctx.Request.Header.Get("ShareAppToken")
	if ShareAppToken != "" {
		var getToken models.UserApiAccessToken
		err := models.Db.Model(&models.UserApiAccessToken{}).Where("access_token = ?", ShareAppToken).First(&getToken).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.Redirect(302, "/")
		}
	} else {
		token := ctx.Request.Header.Get("Authorization")
		if token == "" {
			curl := findCustomAPI(ctx.Request.RequestURI)
			if !strings.Contains(ctx.Request.RequestURI, "/GetRecordVideoTree") &&
				!strings.Contains(ctx.Request.RequestURI, "/PlayVideo") &&
				!strings.Contains(ctx.Request.RequestURI, "/GetSystemNetworkInfo") &&
				!strings.Contains(ctx.Request.RequestURI, "/GetSnapTree") &&
				!strings.Contains(ctx.Request.RequestURI, "/SnapToJpg") &&
				!strings.Contains(ctx.Request.RequestURI, "/WitePhysicalID") &&
				!strings.Contains(ctx.Request.RequestURI, "/GetPhysicalIDCheck") &&
				!strings.Contains(ctx.Request.RequestURI, "/UpdateDataModel") &&
				!strings.Contains(ctx.Request.RequestURI, "/UpdateIEC104DataModel") &&
				!strings.Contains(ctx.Request.RequestURI, "/UpdateModbusTcpPushDataModel") &&
				!strings.Contains(ctx.Request.RequestURI, "/setData") &&
				!strings.Contains(ctx.Request.RequestURI, "/GetSystemParams") &&
				!strings.Contains(ctx.Request.RequestURI, "/GetSystemDeviceInfo") &&
				!strings.Contains(ctx.Request.RequestURI, "/getDisplayModelLayerData") &&
				!strings.Contains(ctx.Request.RequestURI, "/getDisplayModelLayerDataByToken") &&
				!strings.Contains(ctx.Request.RequestURI, "/getRealDataByUuid") &&
				!strings.Contains(ctx.Request.RequestURI, "/GetCustomPel") &&
				!strings.Contains(ctx.Request.RequestURI, "/DisplayLoginPage") &&
				!strings.Contains(ctx.Request.RequestURI, "/LocalUpgrade") &&
				!strings.Contains(ctx.Request.RequestURI, "/AuthUpload") &&
				!strings.Contains(ctx.Request.RequestURI, "/BackupUpload") &&
				!strings.Contains(ctx.Request.RequestURI, "/DiyUpload") &&
				!strings.Contains(ctx.Request.RequestURI, "/api/v1/PushDeviceData") &&
				!strings.Contains(ctx.Request.RequestURI, "/api/v1/open") &&
				!strings.Contains(ctx.Request.RequestURI, "/IEC61850NodeidImport") &&
				!strings.Contains(ctx.Request.RequestURI, "/opcuaNodeidImport") &&
				!strings.Contains(ctx.Request.RequestURI, "/displayImageUpload") &&
				!strings.Contains(ctx.Request.RequestURI, "/GetAuthLicenseInfo") &&
				ctx.Request.RequestURI != "/login" &&
				ctx.Request.RequestURI != "/snmpmodelimportxml" &&
				ctx.Request.RequestURI != "/snmpmodelimport" &&
				ctx.Request.RequestURI != "/systemImageUpload" &&
				!strings.Contains(ctx.Request.RequestURI, "/webrtcstream/") &&
				!strings.Contains(ctx.Request.RequestURI, "ticket") &&
				!strings.Contains(ctx.Request.RequestURI, "/codec/") &&
				ctx.Request.RequestURI != "/" &&
				ctx.Request.RequestURI != "/saveLicense" &&
				ctx.Request.RequestURI != "/checkLicense" &&
				!curl {
				ctx.Redirect(302, "/")
			}
		} else {
			result, _, _, _, _ := middleware.JwtToken(token)
			if !strings.Contains(ctx.Request.RequestURI, "/GetRecordVideoTree") &&
				!strings.Contains(ctx.Request.RequestURI, "/PlayVideo") &&
				!strings.Contains(ctx.Request.RequestURI, "/GetSystemNetworkInfo") &&
				!strings.Contains(ctx.Request.RequestURI, "/GetSnapTree") &&
				!strings.Contains(ctx.Request.RequestURI, "/SnapToJpg") &&
				!strings.Contains(ctx.Request.RequestURI, "/WitePhysicalID") &&
				!strings.Contains(ctx.Request.RequestURI, "/GetPhysicalIDCheck") &&
				!strings.Contains(ctx.Request.RequestURI, "/UpdateDataModel") &&
				!strings.Contains(ctx.Request.RequestURI, "/UpdateIEC104DataModel") &&
				!strings.Contains(ctx.Request.RequestURI, "/UpdateModbusTcpPushDataModel") &&
				!strings.Contains(ctx.Request.RequestURI, "/setData") &&
				!strings.Contains(ctx.Request.RequestURI, "/GetSystemParams") &&
				!strings.Contains(ctx.Request.RequestURI, "/GetSystemDeviceInfo") &&
				!strings.Contains(ctx.Request.RequestURI, "/getDisplayModelLayerData") &&
				!strings.Contains(ctx.Request.RequestURI, "/getDisplayModelLayerDataByToken") &&
				!strings.Contains(ctx.Request.RequestURI, "/getRealDataByUuid") &&
				!strings.Contains(ctx.Request.RequestURI, "/GetCustomPel") &&
				!strings.Contains(ctx.Request.RequestURI, "/DisplayLoginPage") &&
				!strings.Contains(ctx.Request.RequestURI, "/LocalUpgrade") &&
				!strings.Contains(ctx.Request.RequestURI, "/AuthUpload") &&
				!strings.Contains(ctx.Request.RequestURI, "/BackupUpload") &&
				!strings.Contains(ctx.Request.RequestURI, "/DiyUpload") &&
				!strings.Contains(ctx.Request.RequestURI, "/api/v1/PushDeviceData") &&
				!strings.Contains(ctx.Request.RequestURI, "/api/v1/open") &&
				!strings.Contains(ctx.Request.RequestURI, "/IEC61850NodeidImport") &&
				!strings.Contains(ctx.Request.RequestURI, "/opcuaNodeidImport") &&
				!strings.Contains(ctx.Request.RequestURI, "/displayImageUpload") &&
				!strings.Contains(ctx.Request.RequestURI, "/GetAuthLicenseInfo") &&
				result != errmsg.SUCCSE &&
				ctx.Request.RequestURI != "/snmpmodelimportxml" &&
				ctx.Request.RequestURI != "/snmpmodelimport" &&
				ctx.Request.RequestURI != "/systemImageUpload" &&
				!strings.Contains(ctx.Request.RequestURI, "/webrtcstream/") &&
				!strings.Contains(ctx.Request.RequestURI, "/codec/") &&
				ctx.Request.RequestURI != "/login" &&
				ctx.Request.RequestURI != "/" &&
				ctx.Request.RequestURI != "/saveLicense" &&
				ctx.Request.RequestURI != "/checkLicense" {
				ctx.Redirect(302, "/")
			}
		}
	}
}
