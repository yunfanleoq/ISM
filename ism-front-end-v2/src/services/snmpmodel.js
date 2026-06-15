import {
  SNMPMODELADD,SNMPMODELLIST,SNMPMODELDELETE, SNMPMODELEDIT,SNMPMODELSINGLE,
  SAVEMIB,GETMIBS,DELETEMIBS,MODELDATAEDIT,GETHistoryMIBS
} from '@/services/api'
import {request, METHOD} from '@/utils/request'

/**
 * snmp模型添加
 */
export async function snmpModelAdd(params) {
  return request(SNMPMODELADD, METHOD.POST,params)
}

/**
 * snmp单个模型获取
 */
export async function getSnmpModelDetail(params) {
  return request(SNMPMODELSINGLE, METHOD.POST,params)
}

/**
 * snmp模型修改
 */
export async function snmpModelEdit(params) {
  return request(SNMPMODELEDIT, METHOD.POST,params)
}

/**
 * snmp Mib保存
 */
export async function snmpModelMibSave(params) {
  return request(SAVEMIB, METHOD.POST,params,{
    timeout:600000000
  })
}

/**
 * snmp模型列表
 */
export async function snmpModelList(params) {
  return request(SNMPMODELLIST, METHOD.POST,params)
}

/**
 * snmp模型删除
 */
export async function snmpModelDelete(params) {
  return request(SNMPMODELDELETE, METHOD.POST,params)
}

/**
 * snmp mib获取
 */
export async function snmpModelGetMibs(params) {
  return request(GETMIBS, METHOD.POST,params)
}
/**
 * snmp mib 删除
 */
export async function snmpModelDeleteMibs(params) {
  return request(DELETEMIBS, METHOD.POST,params)
}

/**
 * 数据编辑
 */
export async function modelDataEdit(params) {
  return request(MODELDATAEDIT, METHOD.POST,params)
}

/**
 * 通过设备模型获取数据
 */
export async function getDatasByUuid(params) {
  return request(GETMIBS, METHOD.POST,params)
}

export async function getHistoryDatasByUuid(params) {
  return request(GETHistoryMIBS, METHOD.POST,params)
}
export default {
  snmpModelAdd,
  snmpModelList,
  snmpModelDelete,
  snmpModelEdit,
  getSnmpModelDetail,
  snmpModelMibSave,
  snmpModelGetMibs,
  snmpModelDeleteMibs,
  modelDataEdit,
  getDatasByUuid,
  getHistoryDatasByUuid
}
