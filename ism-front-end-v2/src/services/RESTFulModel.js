import {
  RESTFulMODELADD, RESTFulMODELLIST, RESTFulMODELDEL, RESTFulMODELEDIT,
  RESTFulMODELDATAADD,RESTFulMODELDATADEL,RESTFulMODELDATAEDIT,RESTFulMODELDATALIST
} from '@/services/api'
import {request, METHOD} from '@/utils/request'


/**
 * 模型添加
 */
export async function RESTFulModelAdd(params) {
  return request(RESTFulMODELADD, METHOD.POST,params)
}

/**
 * 模型删除
 */
export async function RESTFulModelDelete(params) {
  return request(RESTFulMODELDEL, METHOD.POST,params)
}

/**
 * modbus模型修改
 */
export async function RESTFulModelEdit(params) {
  return request(RESTFulMODELEDIT, METHOD.POST,params)
}
/**
 * 模型列表
 */
export async function RESTFulModellist(params) {
  return request(RESTFulMODELLIST, METHOD.POST,params)
}

export async function RESTFulModelDataAdd(params) {
  return request(RESTFulMODELDATAADD, METHOD.POST,params)
}
export async function RESTFulModelDataDel(params) {
  return request(RESTFulMODELDATADEL, METHOD.POST,params)
}
export async function RESTFulModelDataEdit(params) {
  return request(RESTFulMODELDATAEDIT, METHOD.POST,params)
}

export async function RESTFulModelDataList(params) {
  return request(RESTFulMODELDATALIST, METHOD.POST,params)
}

export default {
  RESTFulModelAdd,
  RESTFulModelDelete,
  RESTFulModelEdit,
  RESTFulModellist,
  RESTFulModelDataAdd,
  RESTFulModelDataDel,
  RESTFulModelDataEdit,
  RESTFulModelDataList
}
