/**
 * 纯设备容器列表分页（导航树与设备管理一致，无 floor/设备组 中间层）
 */

// 1888×896 运行态内容区扣除外壳/边框/内边距后，可稳定容纳 7×7 张 240×110 卡片。
// 底部分页已移除，49 是设备上下文、顶部页码和卡片切片的唯一容量口径。
export const DEFAULT_DEVICE_PAGE_SIZE = 49

/** 设备列表只允许使用后端声明的 deviceList 模板，不猜名称、不硬编码记录 ID。 */
export function resolveDeviceListTemplateId(templateMap) {
  const map = templateMap || {}
  return map.deviceList || ''
}

/** childrenList 模式推荐每页条数（表格懒加载，与 ViewRealTable 分页一致） */
export function deviceListPageSizeForNav(nav) {
  const n = nav && nav.pageSize
  if (n && n > 0) return n
  return DEFAULT_DEVICE_PAGE_SIZE
}

/** 样本 floor/cabinet 页里的「设备组」面包屑/标题（纯设备列表应剔除） */
export function isFloorGroupArtifact(text) {
  const t = String(text || '').trim()
  if (!t) return false
  const core = t.replace(/^📋\s*/, '').replace(/^←\s*/, '')
  if (/default设备组/.test(core)) return true
  if (/^[A-Z]\d+设备组$/.test(core)) return true
  if (/设备组$/.test(core) && core !== '设备组') return true
  return false
}

import { isPureGatewayContainer, isNavListPaged, formatNavPageInfo } from './navContext'

/** 子节点是否全是 A3 转机（叶容器，如 UPS报警解析 122台） */
export function isPureDeviceContainer(node) {
  return isPureGatewayContainer(node)
}

/** navContext 是否处于「列表分页」模式（A3 转机 / B2 寄存器组 / 虚拟列头柜） */
export function isDeviceListNav(nav) {
  if (!nav) return false
  if (nav.virtualCabinetListMode) return true
  if (isNavListPaged(nav)) return true
  const children = nav.children || []
  const hasContainers = children.some(c =>
    c && c.kind && c.kind !== 'device' && c.kind !== 'gateway' && c.kind !== 'virtualCabinet',
  )
  if (hasContainers) return false
  const total = nav.totalDevices != null
    ? nav.totalDevices
    : (nav.allChildDevices || nav.childDevices || nav.allChildNodes || nav.childNodes || []).length
  return total > 0 && !hasContainers
}

export function paginateDevices(devices, pageIndex, pageSize) {
  const all = devices || []
  const size = Math.max(1, pageSize || DEFAULT_DEVICE_PAGE_SIZE)
  const total = all.length
  const totalPages = Math.max(1, Math.ceil(total / size))
  const idx = Math.max(0, Math.min(pageIndex || 0, totalPages - 1))
  const start = idx * size
  return {
    pageIndex: idx,
    pageSize: size,
    totalDevices: total,
    totalPages,
    pageDevices: all.slice(start, start + size),
  }
}

/**
 * 将 navContext 的 childDevices 切片为当前页，保留 allChildDevices 全量列表。
 * @param {object} nav
 * @param {number} [pageSizeOverride] 模板槽位数（动态）
 */
export function applyDeviceListPagination(nav, pageSizeOverride) {
  if (!nav) return nav
  const all = nav.allChildDevices || nav.childDevices || []
  const size = pageSizeOverride || nav.pageSize || DEFAULT_DEVICE_PAGE_SIZE
  const p = paginateDevices(all, nav.pageIndex, size)
  return {
    ...nav,
    allChildDevices: all,
    childDevices: p.pageDevices,
    pageIndex: p.pageIndex,
    pageSize: p.pageSize,
    totalDevices: p.totalDevices,
    totalPages: p.totalPages,
    deviceListMode: p.totalDevices > 0,
    virtualCabinetListMode: !!(nav.virtualCabinetListMode
      || all.some(d => d && (d.kind === 'virtualCabinet' || d.virtualCabinet))),
  }
}

/** 页码展示文案 */
export function formatPageInfo(nav) {
  const cur = (nav.pageIndex || 0) + 1
  const total = nav.totalPages || 1
  const devTotal = nav.totalDevices || 0
  return `第 ${cur}/${total} 页 · 共 ${devTotal} 台`
}

function getCellText(cell) {
  const style = cell && cell.data && cell.data.detail && cell.data.detail.style
  return style && typeof style.text === 'string' ? style.text : ''
}

function setCellText(cell, text) {
  const style = cell && cell.data && cell.data.detail && cell.data.detail.style
  if (style) style.text = text
}

/** 翻页控件文本与角色（上一页/下一页/页码）— 设备列表容器页通用 */
export function rewritePaginationCells(cells, nav) {
  if (!isDeviceListNav(nav)) return
  const info = formatPageInfo(nav)
  const cur = nav.pageIndex || 0
  const totalPages = nav.totalPages || 1
  cells.forEach(cell => {
    const t = getCellText(cell)
    if (!t) return
    const role = detectPageNavRole(t)
    if (!role && !/\{\{\s*nav\.pageInfo\s*\}\}/.test(t)) return
    const detail = cell.data && cell.data.detail
    if (!detail) return
    if (!detail.style.diy) detail.style.diy = []
    const setRole = r => {
      const item = detail.style.diy.find(d => d && d.key === 'labelRole')
      if (item) item.value = r
      else detail.style.diy.push({ name: 'labelRole', type: 9, value: r, key: 'labelRole' })
    }
    if (role === 'pagePrev') {
      setRole(cur > 0 ? 'pagePrev' : 'pagePrevDisabled')
      if (cur <= 0) setCellText(cell, '‹ 上一页')
    } else if (role === 'pageNext') {
      setRole(cur < totalPages - 1 ? 'pageNext' : 'pageNextDisabled')
      if (cur >= totalPages - 1) setCellText(cell, '下一页 ›')
    } else if (role === 'pageInfo' || /\{\{\s*nav\.pageInfo\s*\}\}/.test(t)) {
      setRole('pageInfo')
      setCellText(cell, info)
      const act = detail.action
      if (Array.isArray(act)) {
        act.forEach(a => {
          if (a && a.link && a.link.Inside) a.link.Inside.pageUUID = ''
        })
      }
    }
  })
}

/** 从模板槽位目标推算每页容量 */
export function inferPageSizeFromTargets(targets) {
  if (!targets || !targets.length) return DEFAULT_DEVICE_PAGE_SIZE
  const deviceRows = targets.filter(t => t.entry && t.entry.oldKind === 'device').length
  if (deviceRows > 0) return deviceRows
  const floorSlots = targets.filter(t => t.entry && t.entry.oldKind === 'floor').length
  if (floorSlots > 0) return floorSlots
  const buildingSlots = targets.filter(t => t.entry && t.entry.oldKind === 'building').length
  if (buildingSlots > 0) return buildingSlots
  return DEFAULT_DEVICE_PAGE_SIZE
}

// 仅带有明确页码语义的文案才是分页按钮。
// 裸 ‹/›/→ 常用于面包屑层级分隔，不能被误判成翻页控件。
const PAGE_PREV_RE = /^(?:(?:‹|←|◀)\s*)?(上一页|上页|前一页)\b/
const PAGE_NEXT_RE = /^(?:(?:›|→|▶)\s*)?(下一页|下页|后一页)\b/

export function detectPageNavRole(text) {
  const t = String(text || '').trim()
  if (!t) return ''
  if (PAGE_PREV_RE.test(t)) return 'pagePrev'
  if (PAGE_NEXT_RE.test(t)) return 'pageNext'
  if (/第\s*\d+\s*\/\s*\d+\s*页/.test(t) || /\{\{\s*nav\.pageInfo\s*\}\}/.test(t)) return 'pageInfo'
  return ''
}

/** 为 navContext 构造翻页后的副本（不改变 pageUuid） */
export function navContextForPage(nav, pageIndex) {
  if (!nav) return null
  return applyDeviceListPagination({
    ...nav,
    pageIndex: Math.max(0, pageIndex),
  })
}
