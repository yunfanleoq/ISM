import {
    ALARMTRIGGERADD,ALARMTRIGGERDEL,ALARMTRIGGEREDIT,ALARMTRIGGERLIST,CURRENTALARMLIST,UPDATECURRENTALARM,SHIELDALARMLIST
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
 */
export async function GetCurrentAlarmList(params) {
    return request(CURRENTALARMLIST, METHOD.POST,params)
}

/**
 * 操作告警
 */
export async function UpdateCurrentAlarm(params) {
    return request(UPDATECURRENTALARM, METHOD.POST,params)
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
    UpdateCurrentAlarm
}