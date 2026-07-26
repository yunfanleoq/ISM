/**
 * 测点实时值展示文案（不影响底层存数）。
 * 支路/开关类状态：1→合闸，0→分闸。
 * 通讯/在线类状态：Online→在线，Offline→离线。
 */

/** 名称命中则按合/分闸语义展示 */
const SWITCH_STATUS_RE = /(支路状态|开关状态|断路器状态|合分闸|分合闸|闸位状态|开关位)/

/** 名称命中则按在线/离线语义展示 */
const ONLINE_STATUS_RE = /(通讯状态|通信状态|在线状态|设备状态|连接状态|链路状态|DeviceStatus|OnlineStatus|CommStatus|ConnectionStatus)/i

const ONLINE_VALUE_MAP = {
  online: '在线',
  offline: '离线',
  on: '在线',
  off: '离线',
  connected: '在线',
  disconnected: '离线',
  running: '运行',
  stopped: '停止',
  stop: '停止',
  alarm: '告警',
  normal: '正常',
  fault: '故障',
  error: '故障',
  warning: '预警',
  idle: '空闲',
  'true': '在线',
  'false': '离线',
}

export function isSwitchStatusPoint(pointName) {
  return SWITCH_STATUS_RE.test(String(pointName || ''))
}

export function isOnlineStatusPoint(pointName) {
  return ONLINE_STATUS_RE.test(String(pointName || ''))
}

/** 系统内置点英文名 → 中文展示名（数据仓库等列表） */
const POINT_NAME_ZH_MAP = {
  'device.devicestatus': '设备状态',
  'devicestatus': '设备状态',
}

export function formatPointDisplayName(pointName) {
  const raw = String(pointName || '').trim()
  if (!raw) return raw
  const key = raw.toLowerCase()
  if (Object.prototype.hasOwnProperty.call(POINT_NAME_ZH_MAP, key)) {
    return POINT_NAME_ZH_MAP[key]
  }
  // 后缀匹配：xxx.DeviceStatus
  const leaf = key.split('.').pop()
  if (Object.prototype.hasOwnProperty.call(POINT_NAME_ZH_MAP, leaf)) {
    return POINT_NAME_ZH_MAP[leaf]
  }
  return raw
}

/**
 * 最后一个下划线前为设备名，之后为测点名。
 * 例：配电室2A1_T1_410_BC线电压 → 设备=配电室2A1_T1_410，测点=BC线电压
 * 无下划线时整串视为测点名，设备名为空。
 */
export function splitNameByLastUnderscore(fullName) {
  const raw = String(fullName || '').trim()
  if (!raw) {
    return { deviceName: '', pointName: '' }
  }
  const idx = raw.lastIndexOf('_')
  if (idx <= 0 || idx >= raw.length - 1) {
    return { deviceName: '', pointName: formatPointDisplayName(raw) }
  }
  return {
    deviceName: raw.slice(0, idx).trim(),
    pointName: formatPointDisplayName(raw.slice(idx + 1).trim()),
  }
}

function mapOnlineStatusValue(value) {
  if (value === undefined || value === null || value === '') return value
  if (value === '-' || value === '—') return value
  if (value === '在线' || value === '离线') return value

  const raw = typeof value === 'string' ? value.trim() : value
  const key = String(raw).toLowerCase()
  if (Object.prototype.hasOwnProperty.call(ONLINE_VALUE_MAP, key)) {
    return ONLINE_VALUE_MAP[key]
  }

  const n = Number(raw)
  if (raw === 1 || raw === '1' || n === 1 || raw === true) return '在线'
  if (raw === 0 || raw === '0' || n === 0 || raw === false) return '离线'
  return value
}

/**
 * @param {string} pointName 测点显示名（可含设备前缀）
 * @param {*} value 原始实时值
 * @returns {*} 展示用文案或原值
 */
export function formatPointDisplayValue(pointName, value) {
  if (value === undefined || value === null || value === '') return value
  if (value === '-' || value === '—') return value

  // 无点名时仍尝试翻译常见英文状态字面量（MQTT Status 等）
  if (!pointName) {
    const mapped = mapOnlineStatusValue(value)
    return mapped !== value ? mapped : value
  }

  if (isOnlineStatusPoint(pointName)) {
    return mapOnlineStatusValue(value)
  }

  // 即使点名未命中，常见英文状态字面量也统一中文化（现场客户可见）
  if (typeof value === 'string') {
    const key = value.trim().toLowerCase()
    if (Object.prototype.hasOwnProperty.call(ONLINE_VALUE_MAP, key)) {
      return ONLINE_VALUE_MAP[key]
    }
  }

  if (!isSwitchStatusPoint(pointName)) return value

  const raw = typeof value === 'string' ? value.trim() : value
  if (raw === '合闸' || raw === '分闸') return raw

  const n = Number(raw)
  if (raw === 1 || raw === '1' || n === 1) return '合闸'
  if (raw === 0 || raw === '0' || n === 0) return '分闸'
  if (raw === true) return '合闸'
  if (raw === false) return '分闸'
  return value
}

/**
 * 设备在线状态点（device.DeviceStatus / 通讯状态等）→ true/false；
 * 无有效值时返回 null，由调用方回退 monitortree Status。
 * 约定：1/Online/在线 → true；0/2/3/Offline/离线等 → false。
 */
export function isDeviceOnlineFromStatusValue(value) {
  if (value === undefined || value === null || value === '') return null
  if (value === '-' || value === '—') return null
  const mapped = formatPointDisplayValue('device.DeviceStatus', value)
  if (mapped === '在线') return true
  if (mapped === '离线') return false
  const n = Number(value)
  if (n === 1) return true
  if (Number.isFinite(n)) return false
  return null
}
