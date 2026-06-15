import {
  BACNETMODELADD, BACNETMODELLIST, BACNETMODELDEL, BACNETMODELEDIT,
  BACNETMODELNODEIDADD,BACNETMODELNODEIDDEL,BACNETMODELNODEIDEDIT,BACNETMODELNODEIDLIST
} from '@/services/api'
import {request, METHOD} from '@/utils/request'


/**
 * 模型添加
 */
export async function BACnetModelAdd(params) {
  return request(BACNETMODELADD, METHOD.POST,params)
}

/**
 * 模型删除
 */
export async function BACnetModelDelete(params) {
  return request(BACNETMODELDEL, METHOD.POST,params)
}

/**
 * modbus模型修改
 */
export async function BACnetModelEdit(params) {
  return request(BACNETMODELEDIT, METHOD.POST,params)
}
/**
 * 模型列表
 */
export async function BACnetModelList(params) {
  return request(BACNETMODELLIST, METHOD.POST,params)
}

export async function BACnetModelNodeIDAdd(params) {
  return request(BACNETMODELNODEIDADD, METHOD.POST,params)
}
export async function BACnetModelNodeIDDel(params) {
  return request(BACNETMODELNODEIDDEL, METHOD.POST,params)
}
export async function BACnetModelNodeIDEdit(params) {
  return request(BACNETMODELNODEIDEDIT, METHOD.POST,params)
}

export async function BACnetModelNodeIDList(params) {
  return request(BACNETMODELNODEIDLIST, METHOD.POST,params)
}

export default {
  BACnetModelAdd,
  BACnetModelDelete,
  BACnetModelEdit,
  BACnetModelList,
  BACnetModelNodeIDAdd,
  BACnetModelNodeIDDel,
  BACnetModelNodeIDEdit,
  BACnetModelNodeIDList
}
