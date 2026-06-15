import {
  CJT188MODELNODEIDADD,CJT188MODELNODEIDDEL,CJT188MODELNODEIDEDIT,CJT188MODELNODEIDLIST
} from '@/services/api'
import {request, METHOD} from '@/utils/request'


export async function CJT188ModelNodeIDAdd(params) {
  return request(CJT188MODELNODEIDADD, METHOD.POST,params)
}
export async function CJT188ModelNodeIDDel(params) {
  return request(CJT188MODELNODEIDDEL, METHOD.POST,params)
}
export async function CJT188ModelNodeIDEdit(params) {
  return request(CJT188MODELNODEIDEDIT, METHOD.POST,params)
}

export async function CJT188ModelNodeIDList(params) {
  return request(CJT188MODELNODEIDLIST, METHOD.POST,params)
}

export default {
  CJT188ModelNodeIDAdd,
  CJT188ModelNodeIDDel,
  CJT188ModelNodeIDEdit,
  CJT188ModelNodeIDList
}
