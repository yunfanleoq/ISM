import {
  IEC104MODELNODEIDADD,IEC104MODELNODEIDDEL,IEC104MODELNODEIDEDIT,IEC104MODELNODEIDLIST
} from '@/services/api'
import {request, METHOD} from '@/utils/request'


export async function IEC104ModelNodeIDAdd(params) {
  return request(IEC104MODELNODEIDADD, METHOD.POST,params)
}
export async function IEC104ModelNodeIDDel(params) {
  return request(IEC104MODELNODEIDDEL, METHOD.POST,params)
}
export async function IEC104ModelNodeIDEdit(params) {
  return request(IEC104MODELNODEIDEDIT, METHOD.POST,params)
}

export async function IEC104ModelNodeIDList(params) {
  return request(IEC104MODELNODEIDLIST, METHOD.POST,params)
}

export default {
  IEC104ModelNodeIDAdd,
  IEC104ModelNodeIDDel,
  IEC104ModelNodeIDEdit,
  IEC104ModelNodeIDList
}
