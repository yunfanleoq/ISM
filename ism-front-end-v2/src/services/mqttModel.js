import {
  MQTTMODELADD, MQTTMODELLIST, MQTTMODELDEL, MQTTMODELEDIT,
  MQTTMODELNODEIDADD,MQTTMODELNODEIDDEL,MQTTMODELNODEIDEDIT,MQTTMODELNODEIDLIST
} from '@/services/api'
import {request, METHOD} from '@/utils/request'


/**
 * 模型添加
 */
export async function MqttModelAdd(params) {
  return request(MQTTMODELADD, METHOD.POST,params)
}

/**
 * 模型删除
 */
export async function MqttModelDelete(params) {
  return request(MQTTMODELDEL, METHOD.POST,params)
}

/**
 * modbus模型修改
 */
export async function MqttModelEdit(params) {
  return request(MQTTMODELEDIT, METHOD.POST,params)
}
/**
 * 模型列表
 */
export async function MqttModelList(params) {
  return request(MQTTMODELLIST, METHOD.POST,params)
}

export async function MqttModelNodeIDAdd(params) {
  return request(MQTTMODELNODEIDADD, METHOD.POST,params)
}
export async function MqttModelNodeIDDel(params) {
  return request(MQTTMODELNODEIDDEL, METHOD.POST,params)
}
export async function MqttModelNodeIDEdit(params) {
  return request(MQTTMODELNODEIDEDIT, METHOD.POST,params)
}

export async function MqttModelNodeIDList(params) {
  return request(MQTTMODELNODEIDLIST, METHOD.POST,params)
}

export default {
  MqttModelAdd,
  MqttModelDelete,
  MqttModelEdit,
  MqttModelList,
  MqttModelNodeIDAdd,
  MqttModelNodeIDDel,
  MqttModelNodeIDEdit,
  MqttModelNodeIDList
}
