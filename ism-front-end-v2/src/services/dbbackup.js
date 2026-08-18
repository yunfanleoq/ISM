import {
    DBBACKUP,GETTABLESLIST,GETBACKUPLIST,DBRESTORE,
    GETDBCONFIG,SETDBCONFIG,DBDOWN,DBDELETEBACKUP,
    HISHISTORYBACKUP,GETHISHISTORYBACKUPLIST,HISHISTORYBACKUPDOWN
} from '@/services/api'
import {request, METHOD} from '@/utils/request'

export async function DbBackup(params) {
    return request(DBBACKUP, METHOD.POST,params,{
        timeout:600000
    })
}

export async function GetTablesList(params) {
    return request(GETTABLESLIST, METHOD.POST,params,{
        timeout:600000
    })
}

export async function GetBackUpList(params) {
    return request(GETBACKUPLIST, METHOD.POST,params,{
        timeout:600000
    })
}
export async function DbRestore(params) {
    return request(DBRESTORE, METHOD.POST,params,{
        timeout:60*60*1000
    })
}
export async function DbDown(params) {
    return request(DBDOWN, METHOD.POST,params,{
        timeout:60*60*1000,
        responseType: 'blob'
    })
}
export async function DbDeleteBackup(params) {
    return request(DBDELETEBACKUP, METHOD.POST,params,{
        timeout:60000
    })
}
export async function HisDbBackup(params) {
    return request(HISHISTORYBACKUP, METHOD.POST,params,{
        timeout:600000
    })
}
export async function GetHisBackUpList(params) {
    return request(GETHISHISTORYBACKUPLIST, METHOD.POST,params,{
        timeout:600000
    })
}
export async function HisDbDown(params) {
    return request(HISHISTORYBACKUPDOWN, METHOD.POST,params,{
        timeout:60*60*1000,
        responseType: 'blob'
    })
}
export async function GetDbConfig(params) {
    return request(GETDBCONFIG, METHOD.POST,params,{
        timeout:600000
    })
}
export async function SetDbConfig(params) {
    return request(SETDBCONFIG, METHOD.POST,params,{
        timeout:600000
    })
}

export default {
    DbBackup,
    GetTablesList,
    GetBackUpList,
    DbRestore,
    SetDbConfig,
    GetDbConfig,
    DbDown,
    DbDeleteBackup,
    HisDbBackup,
    GetHisBackUpList,
    HisDbDown
}
