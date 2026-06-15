import {
  ADDCUSTOMDATA,EDITCUSTOMDATA,DELCUSTOMDATA,GETCUSTOMDATA
} from '@/services/api'
import {request, METHOD} from '@/utils/request'

/**
 * 模型添加
 */
export async function CustomDataAdd(params) {
  return request(ADDCUSTOMDATA, METHOD.POST,params)
}

/**
 * snmp单个模型获取
 */
export async function CustomDataEdit(params) {
  return request(EDITCUSTOMDATA, METHOD.POST,params)
}

/**
 * snmp模型修改
 */
export async function CustomDataDel(params) {
  return request(DELCUSTOMDATA, METHOD.POST,params)
}

/**
 * snmp Mib保存
 */
export async function CustomDataList(params) {
  return request(GETCUSTOMDATA, METHOD.POST,params)
}

export default {
  CustomDataEdit,
  CustomDataAdd,
  CustomDataDel,
  CustomDataList
}
