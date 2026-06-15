import {
  DLT645MODELNODEIDADD,DLT645MODELNODEIDDEL,DLT645MODELNODEIDEDIT,DLT645MODELNODEIDLIST
} from '@/services/api'
import {request, METHOD} from '@/utils/request'


export async function DLT645ModelNodeIDAdd(params) {
  return request(DLT645MODELNODEIDADD, METHOD.POST,params)
}
export async function DLT645ModelNodeIDDel(params) {
  return request(DLT645MODELNODEIDDEL, METHOD.POST,params)
}
export async function DLT645ModelNodeIDEdit(params) {
  return request(DLT645MODELNODEIDEDIT, METHOD.POST,params)
}

export async function DLT645ModelNodeIDList(params) {
  return request(DLT645MODELNODEIDLIST, METHOD.POST,params)
}

export default {
  DLT645ModelNodeIDAdd,
  DLT645ModelNodeIDDel,
  DLT645ModelNodeIDEdit,
  DLT645ModelNodeIDList
}
