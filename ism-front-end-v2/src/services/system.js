import {
  SYSTEMROLESLIST,
  SYSTEMROLEADD,
  SYSTEMROLEEDIT,
  SYSTEMROLEDEL,
  SYSTEMROLEPERMISSIONS,
  SYSTEMROLEUPDATEPERMISSIONS,
  SYSTEMFONTS,
  SETDEBUG,
  GETSYSTEMDATAMODEL,
  LOCALUPGRADE,
  ONLINECHECKUPGRADE,
  BEGINUPGRADE,
  GETDATAMODELDATA,
  GETMEMDATA,
  SETMEMDATA,
  DISPLAYTEMPLETELIST,
  DISPLAYTEMPLETEGET,
  DISPLAYLOGINPAGE,
  EXECSYSTEMSCRIPT,
  SYSTEMAUTHINFO,
  DISABLESYSTEMSCRIPT,
  GETCUSTOMPEL,
  GETDEVICEINFO,
  SYSTEMPARAMS,
  GETSYSTEMANALYSIS,
  SYSTECODE,
  GETSYSTEMNETWORK,
  SYSTECODEWRITE,
  SAVENETWORKINFO,
  REBOOTSYSTEM,
  GETSYSTEMMQTTDATA,
  SAVESYSTEMMQTTDATA,
  GETSYSTEMMODBUSDATA,
  SAVESYSTEMMODBUSDATA,
  REBOOTISMSYSTEM,
  GETSYSTEMPARAMS,
  SAVESYSTEMPARAMS,
  GETSYSTEMHISTORYCONFIG,
  SAVESYSTEMHISTORYCONFIG,
  SAVESYSEXECSQLQUERY,
  SAVESYSTIME,
  TESTSYSNTPCONFIG,
  GETSYSTEMOPCUADATA,
  SAVESSYSTEMOPCUADATA,
} from '@/services/api'
import {request, METHOD} from '@/utils/request'


/**
 * 角色列表
 */
export async function systemRolesList() {
  return request(SYSTEMROLESLIST, METHOD.POST)
}
export async function systemRoleAdd(params) {
  return request(SYSTEMROLEADD, METHOD.POST, params)
}
export async function systemRoleEdit(params) {
  return request(SYSTEMROLEEDIT, METHOD.POST, params)
}
export async function systemRoleDel(params) {
  return request(SYSTEMROLEDEL, METHOD.POST, params)
}
export async function systemRolePermissions(params) {
  return request(SYSTEMROLEPERMISSIONS, METHOD.POST, params)
}
export async function systemRoleUpdatePermissions(params) {
  return request(SYSTEMROLEUPDATEPERMISSIONS, METHOD.POST, params)
}
/**
 * 字体列表
 */
export async function systemFontsList() {
  return request(SYSTEMFONTS, METHOD.POST)
}

export async function SetDebug(params) {
  return request(SETDEBUG, METHOD.POST,params)
}
export async function GetSystemData() {
  return request(GETSYSTEMDATAMODEL, METHOD.POST)
}

export async function GetDataModelData(params) {
  return request(GETDATAMODELDATA, METHOD.POST,params)
}

export async function OnlineCheckUpgrade() {
  return request(ONLINECHECKUPGRADE, METHOD.POST,{
    timeout:12000000
  })
}
export async function LocalUpgrade() {
  return request(LOCALUPGRADE, METHOD.POST)
}
export async function BeginUpgrade() {
  return request(BEGINUPGRADE, METHOD.POST)
}

export async function GetMemData(params) {
  return request(GETMEMDATA, METHOD.POST,params)
}
export async function SetMemData(params) {
  return request(SETMEMDATA, METHOD.POST,params)
}
export async function GetSystemAuthInfo(params) {
  return request(SYSTEMAUTHINFO, METHOD.GET,params,{
    timeout:12000000
  })
}
export async function GetSystemPageTemplete(params) {
  return request(DISPLAYTEMPLETELIST, METHOD.POST)
}
export async function GetSystemPageTempleteContent(params) {
  return request(DISPLAYTEMPLETEGET, METHOD.POST,params)
}
export async function GetDisplayLoginPage(params) {
  return request(DISPLAYLOGINPAGE, METHOD.POST,params)
}
export async function ExecSysScript(params) {
  return request(EXECSYSTEMSCRIPT, METHOD.POST,params)
}
export async function DisableSysScript(params) {
  return request(DISABLESYSTEMSCRIPT, METHOD.POST,params)
}
export async function GetUserCustomPel() {
  return request(GETCUSTOMPEL, METHOD.POST)
}
export async function GetDeviceInfo() {
  return request(GETDEVICEINFO, METHOD.POST)
}
export async function GetSystemAnalysis() {
  return request(GETSYSTEMANALYSIS, METHOD.POST)
}
export async function GetSystemParams() {
  return request(SYSTEMPARAMS, METHOD.POST)
}
export async function GetSystemCodeCheck(params) {
  return request(SYSTECODE, METHOD.POST,params)
}
export async function WriteSystemCode(params) {
  return request(SYSTECODEWRITE, METHOD.POST,params)
}
export async function GetSystemNetwork(params) {
  return request(GETSYSTEMNETWORK, METHOD.POST,params)
}
export async function RebootSystem(params) {
  return request(REBOOTSYSTEM, METHOD.POST,params)
}
export async function SaveSystemNetworkInfo(params) {
  return request(SAVENETWORKINFO, METHOD.POST,params)
}
export async function GetSystemWebParams(params) {
  return request(GETSYSTEMPARAMS, METHOD.POST,params)
}
export async function SaveSystemWebParams(params) {
  return request(SAVESYSTEMPARAMS, METHOD.POST,params)
}
export async function RebootISMSystem(params) {
  return request(REBOOTISMSYSTEM, METHOD.POST,params)
}
export async function GetSystemMqttData(params) {
  return request(GETSYSTEMMQTTDATA, METHOD.POST,params)
}
export async function SaveSystemMqttData(params) {
  return request(SAVESYSTEMMQTTDATA, METHOD.POST,params)
}
export async function GetSystemModbusData(params) {
  return request(GETSYSTEMMODBUSDATA, METHOD.POST,params)
}
export async function SaveSystemModbusData(params) {
  return request(SAVESYSTEMMODBUSDATA, METHOD.POST,params)
}
export async function GetSystemHistoryConfig(params) {
  return request(GETSYSTEMHISTORYCONFIG, METHOD.POST,params)
}
export async function SaveSystemHistoryConfig(params) {
  return request(SAVESYSTEMHISTORYCONFIG, METHOD.POST,params)
}
export async function ISMExecSqlQuery(params) {
  return request(SAVESYSEXECSQLQUERY, METHOD.POST,params)
}
export async function SaveSystemTimeConfig(params) {
  return request(SAVESYSTIME, METHOD.POST,params)
}
export async function TestNTPConfig(params) {
  return request(TESTSYSNTPCONFIG, METHOD.POST,params)
}
export async function GetSystemOpcuaData(params) {
  return request(GETSYSTEMOPCUADATA, METHOD.POST,params)
}
export async function SaveSystemOpcuaData(params) {
  return request(SAVESSYSTEMOPCUADATA, METHOD.POST,params)
}
export default {
  TestNTPConfig,
  SaveSystemTimeConfig,
  ISMExecSqlQuery,
  GetSystemHistoryConfig,
  SaveSystemHistoryConfig,
  GetSystemModbusData,
  SaveSystemModbusData,
  GetSystemMqttData,
  SaveSystemMqttData,
  systemRolesList,
  systemRoleAdd,
  systemRoleEdit,
  systemRoleDel,
  systemRolePermissions,
  systemRoleUpdatePermissions,
  systemFontsList,
  SetDebug,
  GetSystemData,
  OnlineCheckUpgrade,
  LocalUpgrade,
  BeginUpgrade,
  GetDataModelData,
  GetMemData,
  SetMemData,
  GetSystemAuthInfo,
  GetSystemPageTemplete,
  GetSystemPageTempleteContent,
  GetDisplayLoginPage,
  ExecSysScript,
  GetUserCustomPel,
  DisableSysScript,
  GetDeviceInfo,
  GetSystemNetwork,
  GetSystemParams,
  GetSystemAnalysis,
  GetSystemCodeCheck,
  WriteSystemCode,
  RebootSystem,
  RebootISMSystem,
  SaveSystemNetworkInfo,
  GetSystemWebParams,
  SaveSystemWebParams,
  GetSystemOpcuaData,
  SaveSystemOpcuaData,
}