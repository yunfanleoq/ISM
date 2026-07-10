import {
    ALARMTRIGGERADD,ALARMTRIGGERDEL,ALARMTRIGGEREDIT,ALARMTRIGGERLIST,CURRENTALARMLIST,UPDATECURRENTALARM,CLEARALLCURRENTALARM,SHIELDALARMLIST,ALARMEVENTFEED,ALARMTRIGGEREXPORT,ALARMTRIGGERIMPORT
} from '@/services/api'
import {request, METHOD} from '@/utils/request'

/**
 * 触发器添加
 */
export async function AlarmTriggerAdd(params) {
    return request(ALARMTRIGGERADD, METHOD.POST,params)
}

/**
 * 触发器编辑
 */
export async function AlarmTriggerEdit(params) {
    return request(ALARMTRIGGEREDIT, METHOD.POST,params)
}

/**
 * 触发器删除
 */
export async function AlarmTriggerDel(params) {
    return request(ALARMTRIGGERDEL, METHOD.POST,params)
}
/**
 * 触发器获取
 */
export async function GetAlarmTriggerList() {
    return request(ALARMTRIGGERLIST, METHOD.POST)
}

/**
 * 实时告警
 * @param {Object} params 请求体（deviceList / dataList，可空）
 * @param {Object} [config] 透传 axios 配置（如 { headers: { ProjectUuid } }，大屏按路由指定项目）
 */
export async function GetCurrentAlarmList(params, config) {
    return request(CURRENTALARMLIST, METHOD.POST,params, config)
}

/**
 * 操作告警
 */
export async function UpdateCurrentAlarm(params) {
    return request(UPDATECURRENTALARM, METHOD.POST,params)
}

/**
 * 一键清除实时告警（可按当前筛选条件批量清除）
 */
export async function ClearAllCurrentAlarm(params) {
    return request(CLEARALLCURRENTALARM, METHOD.POST, params || {})
}
export async function GetAlarmEventFeed(params, config) {
    return request(ALARMEVENTFEED, METHOD.POST, params || {}, config)
}

export async function AlarmTriggerExport() {
    return request(ALARMTRIGGEREXPORT, METHOD.POST)
}

export async function AlarmTriggerImport(params) {
    return request(ALARMTRIGGERIMPORT, METHOD.POST, params)
}

/**
 * 屏蔽告警
 */
export async function GetCurrentShieldAlarmList(params) {
    return request(SHIELDALARMLIST, METHOD.POST,params)
}
export default {
    AlarmTriggerAdd,
    AlarmTriggerEdit,
    AlarmTriggerDel,
    GetAlarmTriggerList,
    GetCurrentAlarmList,
    GetCurrentShieldAlarmList,
    UpdateCurrentAlarm,
    ClearAllCurrentAlarm
}