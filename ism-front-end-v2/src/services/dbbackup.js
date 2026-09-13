import {
    DBBACKUP,HISDBBACKUP,GETHISBACKUPLIST,HISDBRESTORE,HISDBDELETEBACKUP,HISDBBACKUPUPLOAD,
    GETTABLESLIST,GETBACKUPLIST,DBRESTORE,DBDELETEBACKUP,
    GETDBCONFIG,SETDBCONFIG,DBDOWN
} from '@/services/api'
import {request, METHOD} from '@/utils/request'

export async function DbBackup(params) {
    return request(DBBACKUP, METHOD.POST,params,{
        timeout:600000
    })
}

export async function HisDbBackup(params) {
    return request(HISDBBACKUP, METHOD.POST,params,{
        timeout:60*60*1000
    })
}

export async function GetHisBackUpList(params) {
    return request(GETHISBACKUPLIST, METHOD.POST,params,{
        timeout:600000
    })
}

export async function HisDbRestore(params) {
    return request(HISDBRESTORE, METHOD.POST,params,{
        timeout:60*60*1000
    })
}

export async function HisDbDeleteBackup(params) {
    return request(HISDBDELETEBACKUP, METHOD.POST,params,{
        timeout:600000
    })
}

export async function HisDbBackupUpload(formData) {
    return request(HISDBBACKUPUPLOAD, METHOD.POST, formData, {
        timeout:60*60*1000
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
export async function DbDeleteBackup(params) {
    return request(DBDELETEBACKUP, METHOD.POST,params,{
        timeout:600000
    })
}
export async function DbDown(params) {
    return request(DBDOWN, METHOD.POST,params,{
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
    HisDbBackup,
    GetHisBackUpList,
    HisDbRestore,
    HisDbDeleteBackup,
    HisDbBackupUpload,
    GetTablesList,
    GetBackUpList,
    DbRestore,
    DbDeleteBackup,
    SetDbConfig,
    GetDbConfig,
    DbDown
}
