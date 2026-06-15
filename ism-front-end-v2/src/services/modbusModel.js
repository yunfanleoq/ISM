import {
  COMLIST, MODBUSMODELADD, MODBUSMODELLIST, MODBUSMODELDEL, MODBUSMODELEDIT,MODBUSMODELGROUPDEL,MODBUSMODELGROUPADD,MODBUSMODELGROUPLIST,
  MODBUSMODELREGISTERADDRESSLIST,MODBUSMODELREGISTERADDRESSDEL,MODBUSMODELREGISTERADDRESSEDIT,MODBUSMODELGROUPEDIT,MODBUSMODELREGISTERADDRESSADD
} from '@/services/api'
import {request, METHOD} from '@/utils/request'

/**
 * 串口列表获取
 */
export async function COMListGet() {
  return request(COMLIST, METHOD.POST)
}

/**
 * 模型添加
 */
export async function ModbusModelAdd(params) {
  return request(MODBUSMODELADD, METHOD.POST,params)
}

/**
 * 模型列表
 */
export async function DeviceModellist(params) {
  return request(MODBUSMODELLIST, METHOD.POST,params)
}

/**
 * 模型删除
 */
export async function modbusModelDelete(params) {
  return request(MODBUSMODELDEL, METHOD.POST,params)
}

/**
 * modbus模型修改
 */
export async function modbusModelEdit(params) {
  return request(MODBUSMODELEDIT, METHOD.POST,params)
}

/**
 * modbus寄存器组添加
 */
export async function modbusModelGroupAdd(params) {
  return request(MODBUSMODELGROUPADD, METHOD.POST,params)
}
/**
 * modbus寄存器组编辑
 */
export async function modbusModelGroupEdit(params) {
  return request(MODBUSMODELGROUPEDIT, METHOD.POST,params)
}
/**
 * modbus寄存器组添加
 */
export async function modbusModelGroupDel(params) {
  return request(MODBUSMODELGROUPDEL, METHOD.POST,params)
}

/**
 * modbus寄存器组列表
 */
export async function modbusModelGroupList(params) {
  return request(MODBUSMODELGROUPLIST, METHOD.POST,params)
}

/**
 * modbus寄存器组里的寄存器列表
 */
export async function modbusModelRegisterList(params) {
  return request(MODBUSMODELREGISTERADDRESSLIST, METHOD.POST,params)
}

/**
 * modbus寄存器组里的寄存器编辑
 */
export async function modbusModelRegisterEdit(params) {
  return request(MODBUSMODELREGISTERADDRESSEDIT, METHOD.POST,params)
}
/**
 * modbus寄存器组里的寄存器添加
 */
export async function modbusModelRegisterAdd(params) {
  return request(MODBUSMODELREGISTERADDRESSADD, METHOD.POST,params)
}

/**
 * modbus寄存器组里的寄存器删除
 */
export async function modbusModelRegisterDel(params) {
  return request(MODBUSMODELREGISTERADDRESSDEL, METHOD.POST,params)
}


export default {
  COMListGet,
  modbusModelGroupEdit,
  modbusModelRegisterAdd,
  ModbusModelAdd,
  DeviceModellist,
  modbusModelDelete,
  modbusModelEdit,
  modbusModelGroupAdd,
  modbusModelGroupDel,
  modbusModelGroupList,
  modbusModelRegisterList
}
