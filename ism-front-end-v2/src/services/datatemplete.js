import {
    ADDTEMPLETEDATA,
    DELTEMPLETEDATA,
    EDITTEMPLETEDATA,
    TEMPLETEDATALIST,
    CHECKSCRIPT
} from '@/services/api'
import {request, METHOD} from '@/utils/request'

export async function AddTempleteData(params) {
    return request(ADDTEMPLETEDATA, METHOD.POST,params)
}
export async function DelTempleteData(params) {
    return request(DELTEMPLETEDATA, METHOD.POST,params)
}
export async function EditTempleteData(params) {
    return request(EDITTEMPLETEDATA, METHOD.POST,params)
}
export async function GetTempleteDataList() {
    return request(TEMPLETEDATALIST, METHOD.POST)
}
export async function CheckTempleteData(params) {
    return request(CHECKSCRIPT, METHOD.POST,params)
}
export default {
    AddTempleteData,
    DelTempleteData,
    EditTempleteData,
    GetTempleteDataList,
    CheckTempleteData
}
