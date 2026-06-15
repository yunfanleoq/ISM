import {ADDSQLREPORTTEMPLETE,DelSQLREPORTTEMPLETE,GETSQLREPORTTEMPLETE,EDITSQLREPORTTEMPLETE,EXPORTSQLREPORTTEMPLETE,VIEWSQLREPORTTEMPLETE} from '@/services/api'
import {request, METHOD} from '@/utils/request'

//历史告警
export async function GetSQLReportTempletes(params) {
  return request(GETSQLREPORTTEMPLETE, METHOD.POST, params,{
        timeout:600000000
    })
}

//历史数据
export async function AddSQLReportTemplete(params) {
  return request(ADDSQLREPORTTEMPLETE, METHOD.POST, params,{
        timeout:600000000
    })
}
//自定义数据
export async function DelSQLReportTemplete(params) {
  return request(DelSQLREPORTTEMPLETE, METHOD.POST, params,{
        timeout:600000000
    })
}
export async function EditSQLReportTemplete(params) {
    return request(EDITSQLREPORTTEMPLETE, METHOD.POST, params,{
        timeout:600000000
    })
}
export async function ExportReportTemplate(params) {
    return request(EXPORTSQLREPORTTEMPLETE, METHOD.POST, params,{
        timeout:600000000
    })
}
export async function ViewReportTemplate(params) {
    return request(VIEWSQLREPORTTEMPLETE, METHOD.POST, params,{
        timeout:600000000
    })
}
export default {
    GetSQLReportTempletes,
    AddSQLReportTemplete,
    DelSQLReportTemplete,
    EditSQLReportTemplete,
    ViewReportTemplate
}
