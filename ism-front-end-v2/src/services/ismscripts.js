import {
    ADDSCRIPT,
    DELSCRIPT,
    EDITSCRIPT,
    SCRIPTLIST,
    CHECKSCRIPT
} from '@/services/api'
import {request, METHOD} from '@/utils/request'

export async function AddScript(params) {
    return request(ADDSCRIPT, METHOD.POST,params)
}
export async function DelScript(params) {
    return request(DELSCRIPT, METHOD.POST,params)
}
export async function EditScript(params) {
    return request(EDITSCRIPT, METHOD.POST,params)
}
export async function GetScriptList() {
    return request(SCRIPTLIST, METHOD.POST)
}
export async function CheckScript(params) {
    return request(CHECKSCRIPT, METHOD.POST,params)
}
export default {
    AddScript,
    DelScript,
    EditScript,
    GetScriptList,
    CheckScript
}
