import {
  IEC61850MODELADD, IEC61850MODELLIST, IEC61850MODELDEL, IEC61850MODELEDIT,
  IEC61850MODELNODEIDADD,IEC61850MODELNODEIDDEL,IEC61850MODELNODEIDEDIT,IEC61850MODELNODEIDLIST
} from '@/services/api'
import {request, METHOD} from '@/utils/request'


/**
 * 模型添加
 */
export async function IEC61850ModelAdd(params) {
  return request(IEC61850MODELADD, METHOD.POST,params)
}

/**
 * 模型删除
 */
export async function IEC61850ModelDelete(params) {
  return request(IEC61850MODELDEL, METHOD.POST,params)
}

/**
 * modbus模型修改
 */
export async function IEC61850ModelEdit(params) {
  return request(IEC61850MODELEDIT, METHOD.POST,params)
}
/**
 * 模型列表
 */
export async function IEC61850Modellist(params) {
  return request(IEC61850MODELLIST, METHOD.POST,params)
}

export async function IEC61850ModelNodeIDAdd(params) {
  return request(IEC61850MODELNODEIDADD, METHOD.POST,params)
}
export async function IEC61850ModelNodeIDDel(params) {
  return request(IEC61850MODELNODEIDDEL, METHOD.POST,params)
}
export async function IEC61850ModelNodeIDEdit(params) {
  return request(IEC61850MODELNODEIDEDIT, METHOD.POST,params)
}

export async function IEC61850ModelNodeIDList(params) {
  return request(IEC61850MODELNODEIDLIST, METHOD.POST,params)
}

export default {
  IEC61850ModelAdd,
  IEC61850ModelDelete,
  IEC61850ModelEdit,
  IEC61850Modellist,
  IEC61850ModelNodeIDAdd,
  IEC61850ModelNodeIDDel,
  IEC61850ModelNodeIDEdit,
  IEC61850ModelNodeIDList
}
