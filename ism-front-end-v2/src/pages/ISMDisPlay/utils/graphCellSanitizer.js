/**
 * X6 fromJSON 入图前清理 components/cells：
 * - 统一为 { cells: [...] }（兼容 legacy 裸数组）
 * - 从 shape / detail.type 推断并补全 shape
 * - 过滤缺少合法 shape 的 cell（含嵌套 children）
 * - JSON 深拷贝，剥离 Vue Observer，避免 fromJSON 读到残缺字段
 */

function isValidShape(shape) {
  return typeof shape === 'string' && shape.trim().length > 0
}

/** 历史/注入 cells 中未注册的 shape → X6 已注册 shape */
const SHAPE_ALIASES = {
  'ism-view-svg-text': 'view-svg-text',
  'ViewSvgText': 'view-svg-text',
  'ism-view-real-table': 'ism-view-real-table',
}

function normalizeCellShape(shape) {
  if (!shape || typeof shape !== 'string') return shape
  const trimmed = shape.trim()
  return SHAPE_ALIASES[trimmed] || trimmed
}

function inferCellShape(cell) {
  if (!cell || typeof cell !== 'object') return null
  const candidates = [
    cell.shape,
    cell.data && cell.data.shape,
    cell.data && cell.data.detail && cell.data.detail.type,
    cell.detail && cell.detail.type,
  ]
  for (let i = 0; i < candidates.length; i += 1) {
    const normalized = normalizeCellShape(candidates[i])
    if (isValidShape(normalized)) return normalized
  }
  if (cell.source != null && cell.target != null) return 'edge'
  return null
}

/**
 * @param {object} cell
 * @param {string} path 日志路径
 * @param {string[]} dropped
 * @returns {object|null}
 */
function sanitizeCellTree(cell, path, dropped) {
  if (!cell || typeof cell !== 'object') {
    dropped.push(`${path}: not an object`)
    return null
  }
  const shape = inferCellShape(cell)
  if (!isValidShape(shape)) {
    dropped.push(`${path}: missing shape (${String(cell.shape)})`)
    return null
  }
  const next = { ...cell, shape }
  if (Array.isArray(next.children) && next.children.length) {
    const children = []
    next.children.forEach((child, i) => {
      const cleaned = sanitizeCellTree(child, `${path}.children[${i}]`, dropped)
      if (cleaned) children.push(cleaned)
    })
    next.children = children
  }
  return next
}

function deepCloneCells(cells) {
  try {
    return JSON.parse(JSON.stringify(cells))
  } catch (e) {
    return cells
  }
}

/**
 * @param {object|object[]|null|undefined} components
 * @param {{ tag?: string }} [opts]
 * @returns {{ cells: object[] }}
 */
export function sanitizeGraphComponents(components, opts = {}) {
  const tag = opts.tag || 'sanitizeGraphComponents'
  let cells = []
  if (Array.isArray(components)) {
    cells = components
  } else if (components && Array.isArray(components.cells)) {
    cells = components.cells
  } else if (components && components.cells == null) {
    return { cells: [] }
  } else {
    return { cells: [] }
  }

  const dropped = []
  const cleaned = []
  cells.forEach((cell, i) => {
    if (cell == null) {
      dropped.push(`[${i}]: null cell`)
      return
    }
    const item = sanitizeCellTree(cell, `[${i}]`, dropped)
    if (item) cleaned.push(item)
  })

  if (dropped.length) {
    console.warn(
      `[${tag}] dropped ${dropped.length}/${cells.length} invalid cells:`,
      dropped.slice(0, 8),
      dropped.length > 8 ? `…(+${dropped.length - 8} more)` : '',
    )
  }
  return { cells: deepCloneCells(cleaned) }
}
