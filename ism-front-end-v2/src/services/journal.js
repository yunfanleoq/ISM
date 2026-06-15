import {
    JOURNALGET
} from '@/services/api'
import {request, METHOD} from '@/utils/request'

/**
 * 获取操作日志
 */
export async function GetOperationLog(params) {
    return request(JOURNALGET, METHOD.POST,params)
}
/**
 * 获取系统日志
 */
export async function GetSystemRunLog(params) {
    return request(JOURNALGET, METHOD.POST,params)
}


export default {
    GetSystemRunLog,
    GetOperationLog,
}