/**
 * 设备内虚拟设备（原虚拟列头柜）：
 * 按数据仓库同款规则拆分——测点名最后一个 `_` 前为设备名（splitNameByLastUnderscore）。
 * 无设备名前缀的点位归入组织树底层真设备名（fallbackDeviceName）。
 * 不写 monitor_list；首页 / 侧栏 / 钻探共用懒加载缓存。
 */

import { fetchDeviceDatapoints } from './navContext'
import { splitNameByLastUnderscore } from './pointValueDisplay'

const _cabinetCache = Object.create(null)
const _inflight = Object.create(null)

/** 缓存版本：规则从 first `_` 改为 last `_` 后自动失效旧缓存 */
const CACHE_VER = 'dw-last-v1'

/**
 * 与数据仓库「设备名」列一致：最后一个 `_` 前的段。
 * @param {string} name
 */
export function extractPointPrefix(name) {
  return splitNameByLastUnderscore(name).deviceName || ''
}

/**
 * 列头 A→J，再列尾 A→J，其余按拼音/字面量
 * @param {string} a
 * @param {string} b
 */
export function compareCabinetPrefix(a, b) {
  const rank = (p) => {
    const m = String(p).match(/^([A-Za-z])列(头|尾)$/)
    if (!m) return { series: 2, letter: 99, raw: p }
    const letter = m[1].toUpperCase().charCodeAt(0) - 65
    const series = m[2] === '头' ? 0 : 1
    return { series, letter, raw: p }
  }
  const ra = rank(a)
  const rb = rank(b)
  if (ra.series !== rb.series) return ra.series - rb.series
  if (ra.letter !== rb.letter) return ra.letter - rb.letter
  return String(ra.raw).localeCompare(String(rb.raw), 'zh')
}

/**
 * @param {object[]} points
 * @param {string} [fallbackDeviceName] 无设备名前缀时归属的真设备名
 * @returns {{ prefix: string, count: number, points: object[], isFallbackGroup: boolean }[]}
 */
export function groupDatapointsByPrefix(points, fallbackDeviceName = '') {
  const fallback = String(fallbackDeviceName || '').trim()
  const map = Object.create(null)
  ;(points || []).forEach((p) => {
    const name = p && (p.name || p.Name || p.label)
    let prefix = extractPointPrefix(name)
    let fromFallback = false
    if (!prefix) {
      if (!fallback) return
      prefix = fallback
      fromFallback = true
    }
    if (!map[prefix]) {
      map[prefix] = {
        prefix,
        count: 0,
        points: [],
        isFallbackGroup: false,
      }
    }
    if (fromFallback || prefix === fallback) {
      map[prefix].isFallbackGroup = true
    }
    map[prefix].count += 1
    map[prefix].points.push(p)
  })
  return Object.keys(map)
    .sort(compareCabinetPrefix)
    .map(k => map[k])
}

/** @param {{ prefix: string, count: number }[]} groups */
export function shouldUseVirtualCabinetLayer(groups) {
  return Array.isArray(groups) && groups.length >= 2
}

export function virtualCabinetCacheKey(device) {
  if (!device) return ''
  const uuid = String(device.uuid || device.deviceUuid || '').trim()
  const muid = String(device.modelUuid || device.muid || '').trim()
  const sid = device.sid != null ? String(device.sid) : ''
  const fallback = String(
    device.fallbackDeviceName || device.label || device.name || device.rawLabel || '',
  ).trim()
  return `${CACHE_VER}|${uuid}|${muid}|${sid}|${fallback}`
}

export function clearVirtualCabinetCache(key) {
  if (key) {
    delete _cabinetCache[key]
    delete _inflight[key]
    return
  }
  Object.keys(_cabinetCache).forEach(k => { delete _cabinetCache[k] })
  Object.keys(_inflight).forEach(k => { delete _inflight[k] })
}

/**
 * 伪树/列表子节点（kind=virtualCabinet）
 * @param {object} deviceNode 真设备
 * @param {{ prefix: string, count: number, isFallbackGroup?: boolean }[]} groups
 */
export function buildVirtualCabinetChildNodes(deviceNode, groups) {
  const parentUuid = deviceNode.uuid || deviceNode.deviceUuid || ''
  const parentLabel = deviceNode.label || deviceNode.name || deviceNode.rawLabel || ''
  const muid = deviceNode.modelUuid || deviceNode.muid || ''
  const sid = deviceNode.sid
  return (groups || []).map((g, index) => {
    const isFallback = !!(g.isFallbackGroup || (parentLabel && g.prefix === parentLabel))
    return {
      id: `vc-${parentUuid || sid || parentLabel}-${g.prefix}-${index}`,
      label: g.prefix,
      name: g.prefix,
      icon: '▣',
      kind: 'virtualCabinet',
      layer: 'device',
      type: 1,
      virtualCabinet: g.prefix,
      virtualCabinetFallback: isFallback,
      isFallbackGroup: isFallback,
      parentDeviceLabel: parentLabel,
      parentDeviceUuid: parentUuid,
      sid,
      uuid: parentUuid,
      deviceUuid: parentUuid,
      modelUuid: muid,
      muid,
      status: deviceNode.status || 'off',
      count: g.count,
      pointCount: g.count,
      code: g.prefix,
      children: [],
    }
  })
}

/**
 * 懒加载：拉测点名 → 分组 → 缓存（只缓存分组摘要，不长期持有全量 points）
 * @param {object} device
 * @param {string} [fallbackDeviceName] 默认取 device.label
 * @returns {Promise<{ groups: object[], cabinets: object[], enabled: boolean }>}
 */
export async function ensureVirtualCabinetsForDevice(device, fallbackDeviceName) {
  const fallback = String(
    fallbackDeviceName
    || (device && (device.fallbackDeviceName || device.label || device.name || device.rawLabel))
    || '',
  ).trim()
  const deviceWithFallback = device
    ? { ...device, fallbackDeviceName: fallback, label: device.label || fallback }
    : device
  const key = virtualCabinetCacheKey(deviceWithFallback)
  if (!key || key === `${CACHE_VER}||||`) {
    return { groups: [], cabinets: [], enabled: false }
  }
  if (_cabinetCache[key]) {
    return _cabinetCache[key]
  }
  if (_inflight[key]) {
    return _inflight[key]
  }

  const job = (async () => {
    const muid = device.modelUuid || device.muid || ''
    const uuid = device.uuid || device.deviceUuid || ''
    // 不用设备显示名做 namePrefix：点名设备段来自 last `_`，与组织树名可能不同
    const points = await fetchDeviceDatapoints(muid, '', uuid)
    const groups = groupDatapointsByPrefix(points, fallback).map(g => ({
      prefix: g.prefix,
      count: g.count,
      isFallbackGroup: !!g.isFallbackGroup,
    }))
    const enabled = shouldUseVirtualCabinetLayer(groups)
    const cabinets = enabled
      ? buildVirtualCabinetChildNodes({ ...device, label: device.label || fallback }, groups)
      : []
    const result = { groups, cabinets, enabled }
    _cabinetCache[key] = result
    return result
  })()

  _inflight[key] = job
  try {
    return await job
  } finally {
    delete _inflight[key]
  }
}

/**
 * 真设备 → 虚拟柜列表 navContext（复用设备列表模板）
 */
export function buildVirtualCabinetListContext(deviceNode, cabinets, ancestors = [], extras = {}) {
  const all = Array.isArray(cabinets) ? cabinets : []
  const label = deviceNode.label || deviceNode.name || ''
  const childDevices = all.map(c => ({
    name: c.name || c.label,
    label: c.label || c.name,
    uuid: c.uuid || c.deviceUuid,
    deviceUuid: c.deviceUuid || c.uuid,
    code: c.code || c.virtualCabinet || c.name,
    modelUuid: c.modelUuid || c.muid,
    muid: c.muid || c.modelUuid,
    sid: c.sid,
    kind: 'virtualCabinet',
    virtualCabinet: c.virtualCabinet || c.name,
    virtualCabinetFallback: !!(c.virtualCabinetFallback || c.isFallbackGroup),
    isFallbackGroup: !!(c.virtualCabinetFallback || c.isFallbackGroup),
    parentDeviceLabel: label,
    parentDeviceUuid: deviceNode.uuid || '',
    status: c.status || 'off',
    pointCount: c.pointCount || c.count,
  }))
  return {
    layer: 'device',
    modbusLayer: 'device',
    kind: 'device',
    routeMode: 'childrenList',
    deviceListMode: true,
    virtualCabinetListMode: true,
    signalMode: false,
    treeDepth: deviceNode.treeDepth,
    sid: deviceNode.sid != null ? deviceNode.sid : null,
    uuid: deviceNode.uuid || '',
    deviceUuid: deviceNode.uuid || '',
    gatewayUuid: deviceNode.uuid || '',
    name: label,
    label,
    code: deviceNode.code || deviceNode.uuid || '',
    modelUuid: deviceNode.modelUuid || deviceNode.muid || '',
    muid: deviceNode.modelUuid || deviceNode.muid || '',
    status: deviceNode.status || 'off',
    childNodes: all,
    allChildNodes: all,
    childDevices,
    allChildDevices: childDevices,
    children: [],
    ancestors: ancestors || [],
    pageIndex: 0,
    pageSize: extras.pageSize || 49,
    totalCount: all.length,
    totalDevices: all.length,
    homePageUuid: extras.homePageUuid || deviceNode.homePageUuid || '',
    deviceListReturnContext: extras.deviceListReturnContext || null,
  }
}

export function isVirtualCabinetNode(node) {
  if (!node) return false
  if (node.kind === 'virtualCabinet') return true
  if (node.virtualCabinet) return true
  return false
}

/** 列头/列尾前缀形态，如 A列头、E列尾 */
export function looksLikeCabinetPrefix(name) {
  return /^[A-Za-z]列[头尾]$/.test(String(name || '').trim())
}

/**
 * 设备列表卡点击时：在虚拟柜列表页把「E列头」还原为虚拟柜节点。
 * @param {object} device 卡片 source
 * @param {object|null} listNav 当前列表 navContext
 */
export function normalizeVirtualCabinetClick(device, listNav) {
  if (!device) return null
  const nav = listNav || null
  const prefix = String(
    device.virtualCabinet
    || ((nav && nav.virtualCabinetListMode && looksLikeCabinetPrefix(device.name || device.label))
      ? (device.name || device.label)
      : '')
    || (isVirtualCabinetNode(device) ? (device.name || device.label) : '')
    || '',
  ).trim()
  if (!prefix && !isVirtualCabinetNode(device)) return null

  const cabinetPrefix = prefix || String(device.virtualCabinet || device.name || device.label || '').trim()
  if (!cabinetPrefix) return null

  const parentUuid = String(
    device.parentDeviceUuid
    || (nav && nav.virtualCabinetListMode && (nav.deviceUuid || nav.uuid))
    || device.deviceUuid
    || device.uuid
    || '',
  ).trim()
  const parentLabel = String(
    device.parentDeviceLabel
    || (nav && nav.virtualCabinetListMode && (nav.name || nav.label))
    || '',
  ).trim()
  const muid = String(
    device.modelUuid || device.muid
    || (nav && (nav.modelUuid || nav.muid))
    || '',
  ).trim()
  const isFallback = !!(
    device.virtualCabinetFallback
    || device.isFallbackGroup
    || (parentLabel && cabinetPrefix === parentLabel)
  )

  return {
    id: device.id || `vc-${parentUuid}-${cabinetPrefix}`,
    label: cabinetPrefix,
    name: cabinetPrefix,
    kind: 'virtualCabinet',
    layer: 'device',
    type: 1,
    virtualCabinet: cabinetPrefix,
    virtualCabinetFallback: isFallback,
    isFallbackGroup: isFallback,
    parentDeviceLabel: parentLabel,
    parentDeviceUuid: parentUuid,
    uuid: parentUuid,
    deviceUuid: parentUuid,
    modelUuid: muid,
    muid,
    sid: device.sid != null ? device.sid : (nav && nav.sid),
    status: device.status || 'off',
    pointCount: device.pointCount || device.count,
  }
}

/**
 * 测点是否属于某虚拟设备名（对齐数据仓库 last `_`）。
 * fallback 组：另含无设备名前缀的点位。
 */
export function pointBelongsToVirtualDevice(fullName, deviceName, isFallbackGroup = false) {
  const label = String(deviceName || '').trim()
  if (!label) return false
  const split = splitNameByLastUnderscore(fullName)
  if (split.deviceName === label) return true
  if (isFallbackGroup && !split.deviceName) return true
  return false
}

/** 展示用测点名：有设备名前缀则去掉，否则整串（再经 format） */
export function displayPointNameForVirtualDevice(fullName, deviceName) {
  const split = splitNameByLastUnderscore(fullName)
  const label = String(deviceName || '').trim()
  if (label && split.deviceName === label) {
    return split.pointName || fullName
  }
  return split.pointName || String(fullName || '').trim()
}
