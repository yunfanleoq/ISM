/**
 * 任意组织层级 + 信号层（设备 → 测点）。
 * 深度与路由均来自 monitorTree 实际结构。
 */

import { resolveDeviceListTemplateId } from './deviceListPager'
import { resolveDeviceSignalTemplateId } from './deviceSignalTemplate'
import { isRoomZoneName } from './monitorTreeTransform'

/** 信号层：真设备 →（可选虚拟列头柜）→ 测点 */
export const SIGNAL_LAYERS = 3

/**
 * 分析树节点子项构成（transform 后）
 * @returns {{ deviceChildren: object[], containerChildren: object[], total: number }}
 */
export function analyzeNodeChildren(node) {
  const children = (node && node.children) || []
  const deviceChildren = children.filter(c => c && isDeviceNode(c))
  const containerChildren = children.filter(c => c && !isDeviceNode(c))
  return { deviceChildren, containerChildren, total: children.length }
}

/** 是否设备节点（type=1 / gateway / device） */
export function isDeviceNode(node) {
  if (!node) return false
  if (node.type === 1) return true
  if (node.kind === 'gateway' || node.kind === 'device') return true
  if (node.layer === 'gateway') return true
  return isLeafDeviceContainer(node)
}

/** 叶容器：type=0 无子节点，自身即设备（如「数据机房报警解析」） */
export function isLeafDeviceContainer(node) {
  if (!node || node.kind === 'device' || node.kind === 'root') return false
  // 组织/配电室/机房即使暂无子节点，也不得当设备直进测点页（避免「一点配电室全是点位」）
  if (node.kind === 'organization' || node.kind === 'zone' || node.kind === 'room') return false
  const label = node.label || node.name || node.rawLabel || ''
  if (isRoomZoneName(label)) return false
  const { total } = analyzeNodeChildren(node)
  return total === 0
}

/** @deprecated 兼容旧引用 */
export function isLeafGatewayContainer(node) {
  return isLeafDeviceContainer(node)
}

/** @deprecated 兼容旧引用 */
export function isGatewayNode(node) {
  return isDeviceNode(node)
}

/** 子节点是否全是设备（如 UPS 122 台） */
export function isPureDeviceContainer(node) {
  const children = (node && node.children) || []
  return children.length > 0 && children.every(c => c && isDeviceNode(c))
}

/** @deprecated 兼容旧引用 */
export function isPureGatewayContainer(node) {
  return isPureDeviceContainer(node)
}

/**
 * 从节点读取 treeDepth（transform 时写入）；无则按祖先链推算
 */
export function resolveTreeDepth(node, index) {
  if (!node) return 1
  if (node.treeDepth != null) return node.treeDepth
  if (!index || node.sid == null) return 1
  let depth = 1
  let p = index.parentBySid[node.sid]
  while (p) {
    depth++
    p = p.sid != null ? index.parentBySid[p.sid] : null
  }
  return depth
}

/** 组织层深度（由实际祖先链决定，不设上限） */
export function resolveOrgDepth(node, index) {
  return Math.max(1, resolveTreeDepth(node, index))
}

/**
 * 决定 onSelect 路由模式
 * @returns {'org'|'childrenList'|'signal'}
 * - org：非设备节点，画布展示直接 children（容器卡片）
 * - childrenList：非设备节点且子项全是设备（分页列表；含仅 1 台，保证 配电室→设备→点位）
 * - signal：设备节点，进入测点实时表
 */
export function resolveRouteMode(node) {
  if (!node) return 'org'
  if (isDeviceNode(node)) return 'signal'
  const { deviceChildren, containerChildren } = analyzeNodeChildren(node)
  // 直属全是设备 → 设备列表（即使只有 1 台也不跳过列表层直接铺点位）
  if (deviceChildren.length > 0 && containerChildren.length === 0) {
    return 'childrenList'
  }
  return 'org'
}

/**
 * 钻探深度摘要
 * @returns {{
 *   treeDepth: number,
 *   orgDepth: number,
 *   routeMode: string,
 *   totalLayers: number,
 *   signalLayers: number,
 *   isDevice: boolean,
 *   pureDeviceList: boolean,
 *   hasSubContainers: boolean,
 *   modbusLayer: string,
 *   remainingLayers: number,
 * }}
 */
export function detectDrillDepth(node, index) {
  const { deviceChildren, containerChildren } = analyzeNodeChildren(node)
  const treeDepth = resolveTreeDepth(node, index)
  const orgDepth = resolveOrgDepth(node, index)
  const routeMode = resolveRouteMode(node)
  const isDevice = isDeviceNode(node)
  const pureDeviceList = !isDevice && isPureDeviceContainer(node) && deviceChildren.length > 0
  const hasSubContainers = containerChildren.length > 0

  let totalLayers = orgDepth
  if (routeMode === 'signal') {
    totalLayers = orgDepth + SIGNAL_LAYERS - 1
  } else {
    totalLayers = orgDepth + SIGNAL_LAYERS
  }

  const remainingLayers = routeMode === 'signal'
    ? SIGNAL_LAYERS
    : SIGNAL_LAYERS

  const modbusLayer = isDevice ? 'device' : (orgDepth <= 1 ? 'zone' : (orgDepth <= 2 ? 'zone' : 'room'))

  return {
    treeDepth,
    orgDepth,
    routeMode,
    totalLayers,
    signalLayers: SIGNAL_LAYERS,
    remainingLayers,
    isDevice,
    pureDeviceList,
    hasSubContainers,
    modbusLayer,
  }
}

/**
 * onSelect 目标模板页
 */
export function resolveOnSelectPageUuid(node, templateMap, index) {
  const depth = detectDrillDepth(node, index)
  return depth.routeMode === 'signal'
    ? resolveDeviceSignalTemplateId(templateMap)
    : resolveDeviceListTemplateId(templateMap)
}

/**
 * 统一 navContext.layer 字段
 */
export function layerKindFromDepth(depth, node) {
  if (!depth) return 'zone'
  if (depth.routeMode === 'signal') return 'device'
  if (depth.routeMode === 'childrenList') {
    return node && node.kind === 'room' ? 'room' : 'zone'
  }
  if (node && node.kind === 'room') return 'room'
  if (node && node.kind === 'zone') return 'zone'
  if (node && node.kind === 'cabinet') return 'cabinet'
  return 'zone'
}
