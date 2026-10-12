/** noRunTree=1 时从 page JSON 移除组态内嵌「机房监视树」格（非 CSS 隐藏） */
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
  if (/左侧运行态树与 monitortree 同源/.test(text)) return true
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
 * 从 cells 数组 **删除** 监视树相关格（fromJSON 前调用），并左移加宽实时表。
 * @returns {{ removed: number }}
 */
export function removeEmbeddedMonitorTreeCells(cells) {
  if (!Array.isArray(cells) || !cells.length) {
    return { removed: 0 }
  }
  const treeLeft = 24
  let removed = 0
  let shiftWidth = 0

  for (let i = cells.length - 1; i >= 0; i--) {
    const cell = cells[i]
    const det = cell && cell.data && cell.data.detail
    if (!isEmbeddedMonitorTreeCell(det)) continue
    if (typeof cell.x === 'number' && cell.x < 400) {
      shiftWidth = Math.max(shiftWidth, (cell.width || 320) + Math.max(0, cell.x - treeLeft))
    }
    cells.splice(i, 1)
    removed++
  }

  if (shiftWidth > 0) {
    const expand = (cell) => {
      if (!cell || typeof cell.x !== 'number' || cell.x < 200) return
      const delta = cell.x - treeLeft
      if (delta <= 0) return
      cell.x = treeLeft
      cell.width = (cell.width || 1528) + delta
      syncCellGeometry(cell)
    }
    cells.forEach((cell) => {
      if (cell.shape === 'ism-view-real-table') expand(cell)
      const det = cell.data && cell.data.detail
      if (det && det.name === 'plane-sec-table') expand(cell)
    })
  }

  return { removed }
}

/** @deprecated 使用 removeEmbeddedMonitorTreeCells */
export function patchLayerCellsForNoRunTree(cells) {
  const { removed } = removeEmbeddedMonitorTreeCells(cells)
  return removed > 0
}

export function layerNameIsEnvPlane(layerData) {
  return !!(layerData && layerData.name === XUNZHI_JF_ENV_PLANE_PAGE)
}
