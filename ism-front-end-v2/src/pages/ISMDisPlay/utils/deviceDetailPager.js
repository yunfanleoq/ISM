/**
 * 单设备详情页测点/参数行分页（懒加载：仅当前页绑点参与 getRealData 收集）
 */

import { detectPageNavRole } from './deviceListPager'

export const DEFAULT_DETAIL_PAGE_SIZE = 8

function getText(cell) {
  const style = cell && cell.data && cell.data.detail && cell.data.detail.style
  return style && typeof style.text === 'string' ? style.text : ''
}

function setText(cell, text) {
  const style = cell && cell.data && cell.data.detail && cell.data.detail.style
  if (style) style.text = text
}

function hasDataBinding(detail) {
  if (!detail) return false
  let found = false
  const walk = obj => {
    if (!obj || typeof obj !== 'object' || found) return
    if (Array.isArray(obj)) {
      obj.forEach(walk)
      return
    }
    if (obj.condition && typeof obj.condition === 'object') {
      const c = obj.condition
      if (c.dataName || c.dataID) found = true
    }
    Object.keys(obj).forEach(k => {
      if (k === 'condition') return
      const v = obj[k]
      if (v && typeof v === 'object') walk(v)
    })
  }
  walk(detail.active)
  walk(detail.animate)
  return found
}

/** 按 Y 坐标聚合带数据绑定的参数行 */
export function collectDetailPointRows(cells) {
  const rowMap = new Map()
  ;(cells || []).forEach((cell, idx) => {
    const detail = cell && cell.data && cell.data.detail
    if (!detail || !hasDataBinding(detail)) return
    const y = Math.round((cell.y || cell.position?.y || 0) / 8) * 8
    let row = rowMap.get(y)
    if (!row) {
      row = { y, cellIdxs: [] }
      rowMap.set(y, row)
    }
    row.cellIdxs.push(idx)
  })
  return [...rowMap.values()].sort((a, b) => a.y - b.y)
}

export function paginateDetailPoints(rows, pageIndex, pageSize) {
  const all = rows || []
  const size = Math.max(1, pageSize || DEFAULT_DETAIL_PAGE_SIZE)
  const total = all.length
  const totalPages = Math.max(1, Math.ceil(total / size))
  const idx = Math.max(0, Math.min(pageIndex || 0, totalPages - 1))
  const start = idx * size
  return {
    pageIndex: idx,
    pageSize: size,
    totalPoints: total,
    totalPages,
    pageRows: all.slice(start, start + size),
  }
}

export function isDeviceDetailPaged(nav) {
  return !!(nav && (nav.kind === 'device' || nav.kind === 'registerGroup') && nav.detailPointMode)
}

export function formatDetailPageInfo(nav) {
  const cur = (nav.detailPageIndex || 0) + 1
  const total = nav.detailTotalPages || 1
  const n = nav.detailTotalPoints || 0
  return `第 ${cur}/${total} 页 · 共 ${n} 个测点`
}

/**
 * 设备详情页：仅保留当前页测点行，其余数据绑定点位整组删除。
 * @returns {{ cells: object[], nav: object }}
 */
export function applyDeviceDetailPagination(cells, nav, pageSizeOverride) {
  if (!nav || (nav.kind !== 'device' && nav.kind !== 'registerGroup')) {
    return { cells, nav }
  }
  const rows = collectDetailPointRows(cells)
  if (!rows.length) {
    return { cells, nav: { ...nav, detailPointMode: false } }
  }
  const size = pageSizeOverride || nav.detailPageSize || DEFAULT_DETAIL_PAGE_SIZE
  const p = paginateDetailPoints(rows, nav.detailPageIndex, size)
  const keep = new Set()
  p.pageRows.forEach(r => r.cellIdxs.forEach(i => keep.add(i)))
  // 保留无绑定的静态格（标题/基本参数/图表框等）
  const filtered = cells.filter((cell, idx) => {
    const detail = cell && cell.data && cell.data.detail
    if (!detail) return true
    if (!hasDataBinding(detail)) return true
    return keep.has(idx)
  })
  const nextNav = {
    ...nav,
    detailPointMode: p.totalPoints > size,
    detailPageIndex: p.pageIndex,
    detailPageSize: p.pageSize,
    detailTotalPoints: p.totalPoints,
    detailTotalPages: p.totalPages,
  }
  return { cells: filtered, nav: nextNav }
}

/** 翻页控件（设备详情测点页） */
export function rewriteDetailPaginationCells(cells, nav) {
  if (!isDeviceDetailPaged(nav)) return
  const info = formatDetailPageInfo(nav)
  const cur = nav.detailPageIndex || 0
  const totalPages = nav.detailTotalPages || 1
  cells.forEach(cell => {
    const t = getText(cell)
    if (!t) return
    const role = detectPageNavRole(t)
    if (!role && !/\{\{\s*nav\.detailPageInfo\s*\}\}/.test(t)) return
    const detail = cell.data && cell.data.detail
    if (!detail) return
    if (!detail.style.diy) detail.style.diy = []
    const setRole = r => {
      const item = detail.style.diy.find(d => d && d.key === 'labelRole')
      if (item) item.value = r
      else detail.style.diy.push({ name: 'labelRole', type: 9, value: r, key: 'labelRole' })
    }
    if (role === 'pagePrev') {
      setRole(cur > 0 ? 'detailPagePrev' : 'detailPagePrevDisabled')
    } else if (role === 'pageNext') {
      setRole(cur < totalPages - 1 ? 'detailPageNext' : 'detailPageNextDisabled')
    } else if (role === 'pageInfo' || /\{\{\s*nav\.detailPageInfo\s*\}\}/.test(t)) {
      setRole('detailPageInfo')
      setText(cell, info)
    }
  })
}
