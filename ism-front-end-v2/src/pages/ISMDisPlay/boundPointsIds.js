export function asText(value) {
  if (Array.isArray(value)) {
    return value.length ? asText(value[0]) : ''
  }
  if (value === null || value === undefined) return ''
  return String(value)
}

export function bindId(item) {
  return asText(item && (item.dataID || item.dataId))
}

export function deviceId(item) {
  return asText(item && (item.deviceSN || item.deviceUuid || item.duid))
}
