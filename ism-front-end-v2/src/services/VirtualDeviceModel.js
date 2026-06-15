import {
  VIRTUALDEVICEMODELADD, VIRTUALDEVICEMODELLIST, VIRTUALDEVICEMODELDEL, VIRTUALDEVICEMODELEDIT,
  VIRTUALDEVICEMODELDATAADD,VIRTUALDEVICEMODELDATADEL,VIRTUALDEVICEMODELDATAEDIT,VIRTUALDEVICEMODELDATALIST
} from '@/services/api'
import {request, METHOD} from '@/utils/request'


/**
 * 模型添加
 */
export async function VirtualDeviceModelAdd(params) {
  return request(VIRTUALDEVICEMODELADD, METHOD.POST,params)
}

/**
 * 模型删除
 */
export async function VirtualDeviceModelDelete(params) {
  return request(VIRTUALDEVICEMODELDEL, METHOD.POST,params)
}

/**
 * modbus模型修改
 */
export async function VirtualDeviceModelEdit(params) {
  return request(VIRTUALDEVICEMODELEDIT, METHOD.POST,params)
}
/**
 * 模型列表
 */
export async function VirtualDeviceModellist(params) {
  return request(VIRTUALDEVICEMODELLIST, METHOD.POST,params)
}

export async function VirtualDeviceModelDataAdd(params) {
  return request(VIRTUALDEVICEMODELDATAADD, METHOD.POST,params)
}
export async function VirtualDeviceModelDataDel(params) {
  return request(VIRTUALDEVICEMODELDATADEL, METHOD.POST,params)
}
export async function VirtualDeviceModelDataEdit(params) {
  return request(VIRTUALDEVICEMODELDATAEDIT, METHOD.POST,params)
}

export async function VirtualDeviceModelDataList(params) {
  return request(VIRTUALDEVICEMODELDATALIST, METHOD.POST,params)
}

export default {
  VirtualDeviceModelAdd,
  VirtualDeviceModelDelete,
  VirtualDeviceModelEdit,
  VirtualDeviceModellist,
  VirtualDeviceModelDataAdd,
  VirtualDeviceModelDataDel,
  VirtualDeviceModelDataEdit,
  VirtualDeviceModelDataList
}
