export const RECORD_TYPE_HEADER = '存储类型(变化存储、定时存储、即时存储、变化百分比、整点存储)'

export function recordTypeText(value) {
  switch (Number(value)) {
    case 1: return '定时存储'
    case 2: return '即时存储'
    case 3: return '变化百分比'
    case 4: return '整点存储'
    default: return '变化存储'
  }
}
