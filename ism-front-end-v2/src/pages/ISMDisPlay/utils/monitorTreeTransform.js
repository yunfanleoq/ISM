/**
 * 设备树 → 导航树 transform（与设备管理 DeviceTree / ISMEditorDeviceTree 一致）
 * 规则：type=0 容器 / type=1 设备，递归展开，不插入 floor/设备组 中间层。
 * treeDepth：从 RootZone 起算，1~4（组织层上限）。
 */

/** 为 transform 节点附加 treeDepth（parentDepth 默认 0 = RootZone 之上） */
export function attachTreeDepth(node, parentDepth = 0) {
  if (!node) return node
  const depth = parentDepth + 1
  node.treeDepth = depth
  ;(node.children || []).forEach(c => attachTreeDepth(c, depth))
  return node
}

export function isRoomZoneName(name) {
  const n = String(name || '').trim()
  return n.endsWith('配电室') || n.endsWith('模块') || /机房|房间/.test(n)
}

import { isA2GroupName } from './navContext'

/**
 * 容器 kind 判定（三至五层钻探：A1=zone / A2=room / A3=gateway）
 * @param {{ type:number, sid:number, pid:number, name:string, childKinds:string[] }} meta
 */
export function resolveMonitorNodeKind(meta) {
  const { type, sid, pid, name, childKinds } = meta
  // A3 网络转机（monitor_list type=1）
  if (type === 1) return 'gateway'
  if (sid === 1) return 'root'
  // A1 分类容器（RootZone 直属）
  if (pid === 1) return 'zone'
  // A2 中间分组（配电室_机房模块3A1 等）
  if (isA2GroupName(name)) return 'room'
  // 纯设备列表叶容器（UPS 122 台，无中间分组）
  if (childKinds.length && childKinds.every(k => k === 'device' || k === 'gateway')) {
    return 'zone'
  }
  if (isRoomZoneName(name)) return 'room'
  return 'zone'
}

export function iconForMonitorKind(kind, pid) {
  if (kind === 'root') return '🏭'
  if (kind === 'gateway' || kind === 'device') return '🔌'
  if (kind === 'cabinet') return '🗄'
  if (kind === 'zone' && pid === 1) return '🏬'
  if (kind === 'room') return '📁'
  return '📁'
}

/** 统计子树内设备/网关节点总数（含所有层级，非仅直属） */
export function countDevicesInSubtree(children) {
  let total = 0
  for (const c of children || []) {
    if (c.kind === 'device' || c.kind === 'gateway') total++
    total += countDevicesInSubtree(c.children)
  }
  return total
}

/** 去掉重复的 RootZone（sid=1 为权威根） */
export function normalizeRootNodes(nodes) {
  const list = nodes || []
  const rootZones = list.filter(n => {
    const v = n.value || {}
    return (n.text || v.Name) === 'RootZone' && v.type === 0
  })
  if (rootZones.length <= 1) return list
  return list.filter(n => {
    const v = n.value || {}
    const name = n.text || v.Name || ''
    if (name === 'RootZone' && v.sid !== 1) return false
    return true
  })
}
