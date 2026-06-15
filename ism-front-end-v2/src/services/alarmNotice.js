import {
    GETALARMNOTICEPAMAMS,UPDATEALARMNOTICEPAMAMS,TESTESEND
} from '@/services/api'
import {request, METHOD} from '@/utils/request'

/**
 * 获取告警通知
 */
export async function GetAlarmNoticeByType(params) {
    return request(GETALARMNOTICEPAMAMS, METHOD.POST,params)
}

/**
 * 更新告警通知
 */
export async function UpdateAlarmNoticeByType(params) {
    return request(UPDATEALARMNOTICEPAMAMS, METHOD.POST,params)
}

/**
 * 测试Email
 */
export async function TestSend(params) {
    return request(TESTESEND, METHOD.POST,params)
}
export default {
    GetAlarmNoticeByType,
    UpdateAlarmNoticeByType,
    TestSend
}