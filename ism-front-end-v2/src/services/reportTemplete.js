import {ADDREPORTTEMPLETE,DelREPORTTEMPLETE,GETREPORTTEMPLETE,EDITREPORTTEMPLETE,SAVEREPORTTEMPLETE,HANDEXPORT} from '@/services/api'
import {request, METHOD} from '@/utils/request'

//历史告警
export async function GetReportTempletes(params) {
  return request(GETREPORTTEMPLETE, METHOD.POST, params,{
        timeout:600000000
    })
}

//历史数据
export async function AddReportTemplete(params) {
  return request(ADDREPORTTEMPLETE, METHOD.POST, params,{
        timeout:600000000
    })
}
//自定义数据
export async function DelReportTemplete(params) {
  return request(DelREPORTTEMPLETE, METHOD.POST, params,{
        timeout:600000000
    })
}
export async function EditReportTemplete(params) {
    return request(EDITREPORTTEMPLETE, METHOD.POST, params,{
        timeout:600000000
    })
}
export async function SaveReportTemplete(params) {
    return request(SAVEREPORTTEMPLETE, METHOD.POST, params,{
        timeout:600000000
    })
}

//手动导出
export async function HandExportReportTemplete(params) {
    return request(HANDEXPORT, METHOD.POST, params,{
        timeout:600000000
    })
}
export default {
    GetReportTempletes,
    AddReportTemplete,
    DelReportTemplete,
    EditReportTemplete,
    SaveReportTemplete,
    HandExportReportTemplete
}
