import {
  DNP3MODELADD, DNP3MODELLIST, DNP3MODELDEL, DNP3MODELEDIT,
  DNP3MODELNODEIDADD, DNP3MODELNODEIDDEL, DNP3MODELNODEIDEDIT, DNP3MODELNODEIDLIST
} from '@/services/api'
import {request, METHOD} from '@/utils/request'

export async function DNP3ModelAdd(params) {
  return request(DNP3MODELADD, METHOD.POST, params)
}

export async function DNP3ModelList(params) {
  return request(DNP3MODELLIST, METHOD.POST, params)
}

export async function DNP3ModelDel(params) {
  return request(DNP3MODELDEL, METHOD.POST, params)
}

export async function DNP3ModelEdit(params) {
  return request(DNP3MODELEDIT, METHOD.POST, params)
}

export async function DNP3ModelNodeIDAdd(params) {
  return request(DNP3MODELNODEIDADD, METHOD.POST, params)
}

export async function DNP3ModelNodeIDDel(params) {
  return request(DNP3MODELNODEIDDEL, METHOD.POST, params)
}

export async function DNP3ModelNodeIDEdit(params) {
  return request(DNP3MODELNODEIDEDIT, METHOD.POST, params)
}

export async function DNP3ModelNodeIDList(params) {
  return request(DNP3MODELNODEIDLIST, METHOD.POST, params)
}

export default {
  DNP3ModelAdd,
  DNP3ModelList,
  DNP3ModelDel,
  DNP3ModelEdit,
  DNP3ModelNodeIDAdd,
  DNP3ModelNodeIDDel,
  DNP3ModelNodeIDEdit,
  DNP3ModelNodeIDList
}