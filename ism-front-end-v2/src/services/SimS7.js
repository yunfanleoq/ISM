import {
    SimS7MODELADD, SimS7MODELLIST, SimS7MODELDEL, SimS7MODELEDIT,
    SimS7MODELDATAADD,SimS7MODELDATADEL,SimS7MODELDATAEDIT,SimS7MODELDATALIST
} from '@/services/api'
import {request, METHOD} from '@/utils/request'


/**
 * 模型添加
 */
export async function SimS7ModelAdd(params) {
    return request(SimS7MODELADD, METHOD.POST,params)
}

/**
 * 模型删除
 */
export async function SimS7ModelDelete(params) {
    return request(SimS7MODELDEL, METHOD.POST,params)
}

/**
 * modbus模型修改
 */
export async function SimS7ModelEdit(params) {
    return request(SimS7MODELEDIT, METHOD.POST,params)
}
/**
 * 模型列表
 */
export async function SimS7Modellist(params) {
    return request(SimS7MODELLIST, METHOD.POST,params)
}

export async function SimS7ModelDataAdd(params) {
    return request(SimS7MODELDATAADD, METHOD.POST,params)
}
export async function SimS7ModelDataDel(params) {
    return request(SimS7MODELDATADEL, METHOD.POST,params)
}
export async function SimS7ModelDataEdit(params) {
    return request(SimS7MODELDATAEDIT, METHOD.POST,params)
}

export async function SimS7ModelDataList(params) {
    return request(SimS7MODELDATALIST, METHOD.POST,params)
}

export default {
    SimS7ModelAdd,
    SimS7ModelDelete,
    SimS7ModelEdit,
    SimS7Modellist,
    SimS7ModelDataAdd,
    SimS7ModelDataDel,
    SimS7ModelDataEdit,
    SimS7ModelDataList
}
