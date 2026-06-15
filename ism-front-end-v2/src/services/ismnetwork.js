import {
    SETNODECONFIG,
    GETNODECONFIG,
    ADDOUTCONNECT,
    OPTOUTCONNECT,
    EDITOUTCONNECT,
    DELOUTCONNECT,
    GETOUTCONNECT,
    GETINCONNECT
} from '@/services/api'
import {request, METHOD} from '@/utils/request'

/**
 * 获取操作日志
 */
export async function GetNodeConfig(params) {
    return request(GETNODECONFIG, METHOD.POST,params)
}
/**
 * 获取系统日志
 */
export async function SetNodeConfig(params) {
    return request(SETNODECONFIG, METHOD.POST,params)
}

export async function AddOutConnect(params) {
    return request(ADDOUTCONNECT, METHOD.POST,params)
}

export async function EditOutConnect(params) {
    return request(EDITOUTCONNECT, METHOD.POST,params)
}
export async function DelOutConnect(params) {
    return request(DELOUTCONNECT, METHOD.POST,params)
}
export async function OptOutConnect(params) {
    return request(OPTOUTCONNECT, METHOD.POST,params)
}
export async function GetConnectOut(params) {
    return request(GETOUTCONNECT, METHOD.POST,params)
}
export async function GetConnectIn(params) {
    return request(GETINCONNECT, METHOD.POST,params)
}

export default {
    GetNodeConfig,
    SetNodeConfig,
    AddOutConnect,
    EditOutConnect,
    DelOutConnect,
    OptOutConnect,
    GetConnectOut,
    GetConnectIn
}
