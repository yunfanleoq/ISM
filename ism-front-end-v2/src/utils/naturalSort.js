/**
 * 名称升序：中文/英文按 locale，数字按数值（2 < 10）
 */
export function compareNaturalName(a, b) {
  const sa = a == null ? '' : String(a)
  const sb = b == null ? '' : String(b)
  return sa.localeCompare(sb, 'zh-CN', { numeric: true, sensitivity: 'base' })
}

/**
 * 递归按名称升序排序设备树节点
 * @param {Array} nodes
 * @returns {Array}
 */
export function sortMonitorTreeByName(nodes) {
  if (!Array.isArray(nodes) || nodes.length === 0) {
    return nodes || []
  }
  const sorted = nodes.slice().sort((x, y) => {
    const nx = (x && (x.text || x.title || x.name)) || ''
    const ny = (y && (y.text || y.title || y.name)) || ''
    return compareNaturalName(nx, ny)
  })
  sorted.forEach(node => {
    if (node && Array.isArray(node.children) && node.children.length > 0) {
      node.children = sortMonitorTreeByName(node.children)
    }
  })
  return sorted
}
