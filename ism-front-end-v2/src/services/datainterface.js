import {
    ADDINTERFACEDATA,
    DELINTERFACEDATA,
    EDITINTERFACEDATA,
    INTERFACEDATALIST,
    EDITINTERFACESTATUS
} from '@/services/api'
import {request, METHOD} from '@/utils/request'

export async function AddInterfaceData(params) {
    return request(ADDINTERFACEDATA, METHOD.POST,params)
}
export async function DelInterfaceData(params) {
    return request(DELINTERFACEDATA, METHOD.POST,params)
}
export async function EditInterfaceData(params) {
    return request(EDITINTERFACEDATA, METHOD.POST,params)
}
export async function GetInterfaceDataList() {
    return request(INTERFACEDATALIST, METHOD.POST)
}
export async function EditInterfaceStatus() {
    return request(EDITINTERFACESTATUS, METHOD.POST)
}
export default {
    AddInterfaceData,
    DelInterfaceData,
    EditInterfaceData,
    GetInterfaceDataList,
    EditInterfaceStatus
}
