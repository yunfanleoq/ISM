import {GETHISTORYALARMLIST,GETHISTORYDATALIST,
    GETDIYHISTORYDATALIST,GETCHARTHISTORYDATALIST,GETCHARTHISTORYTRENDDATALIST,GETHISTORYDATAREPORT,GETHISTORYHOUR,
    GETHISTORYDIFFERENCEDAY,GETHISTORYDIFFERENCEWEEKLY,GETHISTORYDIFFERENCEMONTH,GETHISTORYDIFFERENCEYEAR,GETCHARTHISTORYTRENDDATABYDATE,
    GETHISTORYDIYDIFFERENCEREPORT
}
from '@/services/api'
import {request, METHOD} from '@/utils/request'

//历史告警
export async function GetAlarmHistoryList(params) {
  return request(GETHISTORYALARMLIST, METHOD.POST, params,{
        timeout:600000000
    })
}

//历史数据
export async function GetDataHistoryList(params) {
  return request(GETHISTORYDATALIST, METHOD.POST, params,{
        timeout:600000000
    })
}
//历史报表
export async function GetDataHistoryReport(params) {
    return request(GETHISTORYDATAREPORT, METHOD.POST, params,{
        timeout:600000000
    })
}
//自定义数据
export async function GetDiyDataHistoryList(params) {
  return request(GETDIYHISTORYDATALIST, METHOD.POST, params,{
        timeout:600000000
    })
}

export async function GetChartDataHistoryList(params) {
    return request(GETCHARTHISTORYDATALIST, METHOD.POST, params,{
        timeout:600000000
    })
}
export async function GetChartDataHistoryTrendList(params) {
    return request(GETCHARTHISTORYTRENDDATALIST, METHOD.POST, params,{
        timeout:600000000
    })
}
export async function GetHistoryExeclFile(Url) {
    return request(Url, METHOD.GET, {

    },{
        timeout:600000000
    })
}
export async function GetHistoryHour(params) {
    return request(GETHISTORYHOUR, METHOD.POST, params,{
        timeout:600000000
    })
}
export async function GetHistoryDayDifference(params) {
    return request(GETHISTORYDIFFERENCEDAY, METHOD.POST, params,{
        timeout:600000000
    })
}
export async function GetHistoryWeeklyDifference(params) {
    return request(GETHISTORYDIFFERENCEWEEKLY, METHOD.POST, params,{
        timeout:600000000
    })
}
export async function GetHistoryMonthDifference(params) {
    return request(GETHISTORYDIFFERENCEMONTH, METHOD.POST, params,{
        timeout:600000000
    })
}
export async function GetHistoryYearDifference(params) {
    return request(GETHISTORYDIFFERENCEYEAR, METHOD.POST, params,{
        timeout:600000000
    })
}
export async function GetHistoryDiyDifferenceReport(params) {
    return request(GETHISTORYDIYDIFFERENCEREPORT, METHOD.POST, params,{
        timeout:600000000
    })
}
export async function GetChartDataHistoryTrendByDate(params) {
    return request(GETCHARTHISTORYTRENDDATABYDATE, METHOD.POST, params,{
        timeout:600000000
    })
}
export default {
  GetAlarmHistoryList,
  GetDataHistoryList,
  GetDiyDataHistoryList,
  GetChartDataHistoryList,
  GetChartDataHistoryTrendList,
    GetHistoryHour,
    GetHistoryDayDifference,
    GetHistoryWeeklyDifference,
    GetHistoryMonthDifference,
    GetHistoryYearDifference,
    GetHistoryDiyDifferenceReport,
    GetChartDataHistoryTrendByDate
}