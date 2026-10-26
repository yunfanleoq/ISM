import {ADDREPORTTEMPLETE,DelREPORTTEMPLETE,GETREPORTTEMPLETE,EDITREPORTTEMPLETE,SAVEREPORTTEMPLETE,GETREPORTTEMPLETEFILE,GETDIYREPORTFILE,HANDEXPORT} from '@/services/api'
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
export async function GetReportTempleteFile(uuid) {
    return request(GETREPORTTEMPLETEFILE, METHOD.POST, {Uuid: uuid},{
        timeout:600000,
        responseType: 'arraybuffer'
    })
}
export async function GetDiyReportFile(path) {
    return request(GETDIYREPORTFILE, METHOD.POST, {path: path},{
        timeout:600000,
        responseType: 'arraybuffer'
    })
}
export function isPackedFile(buf) {
    if (!buf || buf.byteLength < 4) {
        return false
    }
    const u8 = new Uint8Array(buf)
    return u8[0] === 0x50 && u8[1] === 0x4b
}
export function packedFileError(buf) {
    try {
        const text = new TextDecoder('utf-8').decode(buf)
        const obj = JSON.parse(text)
        if (obj && obj.msg) {
            return obj.msg
        }
    } catch (e) {}
    return ''
}
export function downloadArrayBuffer(buf, filename) {
    const blob = new Blob([buf], {type: 'application/octet-stream'})
    const url = URL.createObjectURL(blob)
    const elink = document.createElement('a')
    elink.href = url
    elink.setAttribute('download', filename || 'report')
    elink.style.display = 'none'
    document.body.appendChild(elink)
    setTimeout(() => {
        elink.click()
        document.body.removeChild(elink)
        URL.revokeObjectURL(url)
    }, 66)
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
    GetReportTempleteFile,
    GetDiyReportFile,
    HandExportReportTemplete
}
