import {
    DBBACKUP,HISDBBACKUP,GETHISBACKUPLIST,HISDBRESTORE,HISDBDELETEBACKUP,HISDBDOWNLOAD,HISDBBACKUPUPLOAD,
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
        // 后端硬上限 48 小时，结束后最多再拷 2 小时。页面要比这段更长，否则先弹失败、备份还在跑。
        timeout:51*60*60*1000
    })
}

export async function GetHisBackUpList(params) {
    return request(GETHISBACKUPLIST, METHOD.POST,params,{
        timeout:600000
    })
}

export async function HisDbRestore(params) {
    return request(HISDBRESTORE, METHOD.POST,params,{
        timeout:6*60*60*1000
    })
}

export async function HisDbDeleteBackup(params) {
    return request(HISDBDELETEBACKUP, METHOD.POST,params,{
        timeout:600000
    })
}

export async function HisDbDownload(params) {
    return request(HISDBDOWNLOAD, METHOD.POST, params, {
        timeout:6*60*60*1000,
        responseType: 'blob'
    })
}

export async function HisDbBackupUpload(formData) {
    return request(HISDBBACKUPUPLOAD, METHOD.POST, formData, {
        timeout:6*60*60*1000
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
    HisDbDownload,
    HisDbBackupUpload,
    GetTablesList,
    GetBackUpList,
    DbRestore,
    DbDeleteBackup,
    SetDbConfig,
    GetDbConfig,
    DbDown
}
