/** noRunTree=1 时隐藏组态页左侧「机房监视树」占位格并放宽实时表布局 */
import { XUNZHI_JF_ENV_PLANE_PAGE } from '@/xunzhiJf/resolveDisplayPage'

export const EMBEDDED_MONITOR_TREE_CELL_NAMES = new Set([
  'plane-sec-tree',
  'plane-monitor-tree',
])

export function isEmbeddedMonitorTreeCell(detail) {
  if (!detail) return false
  const name = String(detail.name || '')
  if (EMBEDDED_MONITOR_TREE_CELL_NAMES.has(name)) return true
  const text = String((detail.style && detail.style.text) || '')
  if (/机房[监监]视树/.test(text)) return true
  if (/^①\s*机房[监监]视树/.test(text.trim())) return true
  return false
}

function syncCellGeometry(cell) {
  const det = cell.data && cell.data.detail
  if (!det || !det.style) return
  const pos = det.style.position || {}
  pos.x = cell.x
  pos.y = cell.y
  pos.w = cell.width
  pos.h = cell.height
  det.style.position = pos
  if (cell.position) {
    cell.position.x = cell.x
    cell.position.y = cell.y
  }
  if (cell.size) {
    cell.size.width = cell.width
    cell.size.height = cell.height
  }
}

/**
 * 就地修改 layer cells（fromJSON 前调用）。仅影响寻智平面 noRunTree 场景。
 */
export function patchLayerCellsForNoRunTree(cells) {
  if (!Array.isArray(cells) || !cells.length) return false
  let changed = false
  let hiddenTreeWidth = 0
  let treeLeft = 24

  for (const cell of cells) {
    const det = cell.data && cell.data.detail
    if (!isEmbeddedMonitorTreeCell(det)) continue
    if (typeof cell.x === 'number' && cell.x < 400) {
      hiddenTreeWidth = Math.max(hiddenTreeWidth, (cell.width || 320) + Math.max(0, cell.x - treeLeft))
    }
    if (det.style) det.style.visible = 0
    cell.visible = false
    changed = true
  }

  const realTable = cells.find(c => c && c.shape === 'ism-view-real-table')
  if (realTable && hiddenTreeWidth > 0 && typeof realTable.x === 'number' && realTable.x >= 200) {
    const shift = realTable.x - treeLeft
    realTable.x = treeLeft
    realTable.width = (realTable.width || 1528) + shift
    syncCellGeometry(realTable)
    changed = true
  }

  const secTable = cells.find(c => {
    const det = c && c.data && c.data.detail
    return det && det.name === 'plane-sec-table'
  })
  if (secTable && hiddenTreeWidth > 0 && typeof secTable.x === 'number' && secTable.x >= 200) {
    const shift = secTable.x - treeLeft
    secTable.x = treeLeft
    secTable.width = (secTable.width || 1528) + shift
    syncCellGeometry(secTable)
    changed = true
  }

  return changed
}

export function layerNameIsEnvPlane(layerData) {
  return !!(layerData && layerData.name === XUNZHI_JF_ENV_PLANE_PAGE)
}
