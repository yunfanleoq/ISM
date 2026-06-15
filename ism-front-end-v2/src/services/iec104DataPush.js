import {
  IEC104DATAPUSHADD,IEC104DATAPUSHDEL,IEC104DATAPUSHEDIT,IEC104DATAPUSHLIST
} from '@/services/api'
import {request, METHOD} from '@/utils/request'


export async function IEC104DataPushAdd(params) {
  return request(IEC104DATAPUSHADD, METHOD.POST,params)
}
export async function IEC104DataPushDel(params) {
  return request(IEC104DATAPUSHDEL, METHOD.POST,params)
}
export async function IEC104DataPushEdit(params) {
  return request(IEC104DATAPUSHEDIT, METHOD.POST,params)
}

export async function IEC104DataPushList(params) {
  return request(IEC104DATAPUSHLIST, METHOD.POST,params)
}

export default {
  IEC104DataPushAdd,
  IEC104DataPushDel,
  IEC104DataPushEdit,
  IEC104DataPushList
}
