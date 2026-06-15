import {
  OPCUAMODELADD, OPCUAMODELLIST, OPCUAMODELDEL, OPCUAMODELEDIT,
  OPCUAMODELNODEIDADD,OPCUAMODELNODEIDDEL,OPCUAMODELNODEIDEDIT,OPCUAMODELNODEIDLIST
} from '@/services/api'
import {request, METHOD} from '@/utils/request'


/**
 * 模型添加
 */
export async function OpcuaModelAdd(params) {
  return request(OPCUAMODELADD, METHOD.POST,params)
}

/**
 * 模型删除
 */
export async function OpcuaModelDelete(params) {
  return request(OPCUAMODELDEL, METHOD.POST,params)
}

/**
 * modbus模型修改
 */
export async function OpcuaModelEdit(params) {
  return request(OPCUAMODELEDIT, METHOD.POST,params)
}
/**
 * 模型列表
 */
export async function OpcuaModellist(params) {
  return request(OPCUAMODELLIST, METHOD.POST,params)
}

export async function OpcuaModelNodeIDAdd(params) {
  return request(OPCUAMODELNODEIDADD, METHOD.POST,params)
}
export async function OpcuaModelNodeIDDel(params) {
  return request(OPCUAMODELNODEIDDEL, METHOD.POST,params)
}
export async function OpcuaModelNodeIDEdit(params) {
  return request(OPCUAMODELNODEIDEDIT, METHOD.POST,params)
}

export async function OpcuaModelNodeIDList(params) {
  return request(OPCUAMODELNODEIDLIST, METHOD.POST,params)
}

export default {
  OpcuaModelAdd,
  OpcuaModelDelete,
  OpcuaModelEdit,
  OpcuaModellist,
  OpcuaModelNodeIDAdd,
  OpcuaModelNodeIDDel,
  OpcuaModelNodeIDEdit,
  OpcuaModelNodeIDList
}
