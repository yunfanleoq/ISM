import {
  HJ212MODELADD, HJ212MODELLIST, HJ212MODELDEL, HJ212MODELEDIT,
  HJ212MODELNODEIDADD,HJ212MODELNODEIDDEL,HJ212MODELNODEIDEDIT,HJ212MODELNODEIDLIST,HJ212TEMPLETEIDLIST
} from '@/services/api'
import {request, METHOD} from '@/utils/request'


/**
 * 模型添加
 */
export async function HJ212ModelAdd(params) {
  return request(HJ212MODELADD, METHOD.POST,params)
}

/**
 * 模型删除
 */
export async function HJ212ModelDelete(params) {
  return request(HJ212MODELDEL, METHOD.POST,params)
}

/**
 * modbus模型修改
 */
export async function HJ212ModelEdit(params) {
  return request(HJ212MODELEDIT, METHOD.POST,params)
}
/**
 * 模型列表
 */
export async function HJ212ModelList(params) {
  return request(HJ212MODELLIST, METHOD.POST,params)
}

export async function HJ212ModelNodeIDAdd(params) {
  return request(HJ212MODELNODEIDADD, METHOD.POST,params)
}
export async function HJ212ModelNodeIDDel(params) {
  return request(HJ212MODELNODEIDDEL, METHOD.POST,params)
}
export async function HJ212ModelNodeIDEdit(params) {
  return request(HJ212MODELNODEIDEDIT, METHOD.POST,params)
}

export async function HJ212ModelNodeIDList(params) {
  return request(HJ212MODELNODEIDLIST, METHOD.POST,params)
}
export async function HJ212TempleteIDList() {
  return request(HJ212TEMPLETEIDLIST, METHOD.POST)
}

export default {
  HJ212ModelAdd,
  HJ212ModelDelete,
  HJ212ModelEdit,
  HJ212ModelList,
  HJ212ModelNodeIDAdd,
  HJ212ModelNodeIDDel,
  HJ212ModelNodeIDEdit,
  HJ212ModelNodeIDList,
  HJ212TempleteIDList
}
