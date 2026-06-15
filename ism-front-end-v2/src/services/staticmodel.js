import {
  GETSTATICDATALIST,DELSTATICDATA,EDITSTATICDATA,ADDSTATICDATA
} from '@/services/api'
import {request, METHOD} from '@/utils/request'

/**
 * 模型添加
 */
export async function StaticModelAdd(params) {
  return request(ADDSTATICDATA, METHOD.POST,params)
}

/**
 * snmp单个模型获取
 */
export async function StaticModelEdit(params) {
  return request(EDITSTATICDATA, METHOD.POST,params)
}

/**
 * snmp模型修改
 */
export async function StaticModelDel(params) {
  return request(DELSTATICDATA, METHOD.POST,params)
}

/**
 * snmp Mib保存
 */
export async function GetStaticModelList(params) {
  return request(GETSTATICDATALIST, METHOD.POST,params)
}

export default {
  GetStaticModelList,
  StaticModelDel,
  StaticModelEdit,
  StaticModelAdd
}
