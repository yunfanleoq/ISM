import {
  MODBUSTCPDATAPUSHADD,MODBUSTCPDATAPUSHDEL,MODBUSTCPDATAPUSHEDIT,MODBUSTCPDATAPUSHLIST
} from '@/services/api'
import {request, METHOD} from '@/utils/request'


export async function ModbusTcpDataPushAdd(params) {
  return request(MODBUSTCPDATAPUSHADD, METHOD.POST,params)
}
export async function ModbusTcpDataPushDel(params) {
  return request(MODBUSTCPDATAPUSHDEL, METHOD.POST,params)
}
export async function ModbusTcpDataPushEdit(params) {
  return request(MODBUSTCPDATAPUSHEDIT, METHOD.POST,params)
}

export async function ModbusTcpDataPushList(params) {
  return request(MODBUSTCPDATAPUSHLIST, METHOD.POST,params)
}

export default {
  ModbusTcpDataPushAdd,
  ModbusTcpDataPushDel,
  ModbusTcpDataPushEdit,
  ModbusTcpDataPushList
}
