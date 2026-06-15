import {
    DBBACKUP,GETTABLESLIST,GETBACKUPLIST,DBRESTORE,
    GETDBCONFIG,SETDBCONFIG,DBDOWN
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
        timeout:60*60*1000
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
    DbDown
}
