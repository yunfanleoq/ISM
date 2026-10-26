export  function formatDate (date, fmt) {
    if (/(y+)/.test(fmt)) {
        fmt = fmt.replace(RegExp.$1, (date.getFullYear() + '').substr(4 - RegExp.$1.length))
    }
    let o = {
        'M+': date.getMonth() + 1,
        'd+': date.getDate(),
        'h+': date.getHours(),
        'm+': date.getMinutes(),
        's+': date.getSeconds()
    }
    for (let k in o) {
        if (new RegExp(`(${k})`).test(fmt)) {
            let str = o[k] + ''
            fmt = fmt.replace(RegExp.$1, RegExp.$1.length === 1 ? str : padLeftZero(str))
        }
    }
    return fmt
}
/** 将后端返回的时间（墙钟字符串或 ISO）解析为 Date，避免无时区字符串被当成 UTC */
export function parseLocalDateTime(time) {
    if (time == null || time === '') return null
    if (time instanceof Date) return time
    if (typeof time === 'number') return new Date(time)
    const s = String(time).trim()
    // 已有时区标记：交给原生解析
    if (/[zZ]$|[+-]\d{2}:?\d{2}$/.test(s)) {
        return new Date(s)
    }
    // "2006-01-02 15:04:05" / "2006-01-02T15:04:05" → 按本地墙钟
    const m = s.match(/^(\d{4})-(\d{2})-(\d{2})[ T](\d{2}):(\d{2}):(\d{2})/)
    if (m) {
        return new Date(
            Number(m[1]), Number(m[2]) - 1, Number(m[3]),
            Number(m[4]), Number(m[5]), Number(m[6])
        )
    }
    return new Date(s)
}
function padLeftZero (str) {
    return ('00' + str).substr(str.length)
}
export default {
    formatDate,
    parseLocalDateTime
}