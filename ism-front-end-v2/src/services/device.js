import {
  DEVICESNMPADD,MONITORTREE,MONITORADD,SUPPORTDEVICELIST,
  GETDEVICEMODELDATALIST,PINGICMP,MONITORCOPY,
  MONITORDEL,MONITOREDIT,MONITORREALDATA,SETDATA,MONITORREALDATABYUUID,SETDEVICESTARTORSTOP,MONITORDELALL,
  GETREALDATATOTABLE
} from '@/services/api'
import {request, METHOD} from '@/utils/request'

/**
 * snmp设备添加
 */
export async function deviceOrZoneAdd(params) {
  return request(DEVICESNMPADD, METHOD.POST,params,{
    timeout:12000000
  })
}


export async function getMonitorTree(config) {
  return request(MONITORTREE, METHOD.POST, {}, config)
}

export async function addMonitor(params) {
  return request(MONITORADD, METHOD.POST,params,{
    timeout:12000000
  })
}

export async function TestConnect(params) {
  return request(PINGICMP, METHOD.POST,params,{
    timeout:12000000
  })
}

export async function editMonitor(params) {
  return request(MONITOREDIT, METHOD.POST,params,{
    timeout:12000000
  })
}

export async function delMonitor(params) {
  return request(MONITORDEL, METHOD.POST,params,{
    timeout:12000000
  })
}

export async function getRealData(params) {
  return request(MONITORREALDATA, METHOD.POST,params)
}

export async function getRealDataByUuid(params) {
  return request(MONITORREALDATABYUUID, METHOD.POST,params)
}

export async function getSupportDeviceList(params) {
  return request(SUPPORTDEVICELIST, METHOD.POST,params)
}


export async function setData(params) {
  return request(SETDATA, METHOD.POST,params,{
    timeout:10000
  })
}

export async function GetDeviceModelDataList(params) {
  return request(GETDEVICEMODELDATALIST, METHOD.POST,params)
}


export async function SetDeviceStartOrStop(params) {
  return request(SETDEVICESTARTORSTOP, METHOD.POST,params)
}

export async function CopyDevices(params) {
  return request(MONITORCOPY, METHOD.POST,params,{
    timeout:12000000
  })
}
export async function delAllMonitor(params) {
  return request(MONITORDELALL, METHOD.POST,params,{
    timeout:12000000
  })
}
export async function GetRealDataToTable(params) {
  return request(GETREALDATATOTABLE, METHOD.POST,params,{
    timeout:12000000
  })
}
export default {
  deviceOrZoneAdd,
  getMonitorTree,
  addMonitor,
  getSupportDeviceList,
  delMonitor,
  editMonitor,
  getRealData,
  setData,
  GetDeviceModelDataList,
  getRealDataByUuid,
  SetDeviceStartOrStop,
  TestConnect,
  CopyDevices,
  delAllMonitor,
  GetRealDataToTable
}
