/**
 * 统一钻探 navContext schema
 *
 * 组织层（任意深度）：后端组织节点 → … → 设备
 * 信号层（2 层）：设备 → 测点（ism-view-real-table，无寄存器组）
 */

import { getModelDataPoints, getRealData } from '@/services/device'
import {
  DEFAULT_DEVICE_PAGE_SIZE,
  applyDeviceListPagination,
  paginateDevices,
} from './deviceListPager'
import { REAL_DATA_DEFAULT_PAGE_SIZE } from '@/utils/realDataBatch'
import { isDeviceNode, isPureDeviceContainer } from './drillDepth'
import { DASHBOARD_PERFORMANCE, normalizePageSize } from './dashboardPerformance'
import { splitNameByLastUnderscore } from './pointValueDisplay'

/** @typedef {'zone'|'room'|'device'|'datapoint'} NavLayer */

/** 信号层测点每页容量（10×8 高密度卡片；单设备 5000+ 测点必须服务端分页） */
export const DEFAULT_DATAPOINT_PAGE_SIZE = DASHBOARD_PERFORMANCE.datapointPageSize

/** @deprecated 兼容旧模板 */
export const MODBUS_TO_NAV_LAYER = {
  A1: 'zone',
  A2: 'room',
  A3: 'device',
  B2: 'datapoint',
  B3: 'datapoint',
  device: 'device',
  datapoint: 'datapoint',
}

export { isLeafDeviceContainer, isLeafGatewayContainer, isPureDeviceContainer, isPureGatewayContainer } from './drillDepth'

/** 树节点 → nav layer */
export function resolveModbusLayerFromNode(node) {
  if (!node) return null
  if (node.modbusLayer) return node.modbusLayer
  if (isDeviceNode(node)) return 'device'
  const layer = node.layer || node.kind
  if (layer === 'room') return 'room'
  if (layer === 'zone') return 'zone'
  if (layer === 'datapoint') return 'datapoint'
  return 'zone'
}

/** modbusLayer / kind → navContext.layer */
export function resolveNavLayer(modbusLayer, fallbackKind) {
  if (modbusLayer === 'device' || modbusLayer === 'A3') return 'device'
  if (modbusLayer === 'datapoint' || modbusLayer === 'B2' || modbusLayer === 'B3') return 'datapoint'
  if (modbusLayer && MODBUS_TO_NAV_LAYER[modbusLayer]) {
    return MODBUS_TO_NAV_LAYER[modbusLayer]
  }
  const k = fallbackKind || ''
  if (k === 'device' || k === 'gateway') return 'device'
  if (['zone', 'room', 'datapoint'].includes(k)) return k
  return k || 'zone'
}

/** 兼容旧模板 kind 字段 */
export function resolveLegacyKind(layer) {
  if (layer === 'device' || layer === 'datapoint') return 'device'
  return layer || 'zone'
}

/** 树子节点 → childNodes 列表项 */
export function treeChildToNodeItem(child) {
  if (!child) return null
  if (isDeviceNode(child)) {
    return {
      name: child.label || child.name || '',
      label: child.label || child.name || '',
      uuid: child.uuid || '',
      sid: child.sid != null ? child.sid : null,
      modelUuid: child.modelUuid || child.muid || '',
      muid: child.modelUuid || child.muid || '',
      status: child.status || 'off',
      kind: 'device',
      layer: 'device',
      treeDepth: child.treeDepth,
    }
  }
  return {
    name: child.label || child.name || '',
    label: child.label || child.name || '',
    uuid: child.uuid || '',
    sid: child.sid != null ? child.sid : null,
    kind: child.kind || 'zone',
    layer: child.layer || child.kind || 'zone',
    modbusLayer: child.modbusLayer || resolveModbusLayerFromNode(child),
    treeDepth: child.treeDepth,
  }
}

/** 从树节点收集直接 childNodes */
export function collectChildNodesFromTree(node) {
  return ((node && node.children) || [])
    .map(treeChildToNodeItem)
    .filter(Boolean)
}

/** 测点 → childNodes（信号层） */
export function datapointsToChildNodes(points) {
  return (points || []).map(p => ({
    name: p.name || p.Name || '',
    label: p.name || p.Name || '',
    uuid: p.uuid || p.Uuid || '',
    dataID: p.uuid || p.Uuid || '',
    unit: p.unit || p.Unit || p.dataUnit || '',
    value: p.value != null ? p.value : '',
    deviceName: p.deviceName || p.device_name || '',
    deviceUuid: p.deviceUuid || p.device_uuid || '',
    modelDataUuid: p.modelDataUuid || '',
    kind: 'datapoint',
    layer: 'datapoint',
  }))
}

/**
 * 构造基础 navContext（组织层：展示直接 children）
 */
export function createBaseNavContext(node, ancestors = []) {
  const modbusLayer = resolveModbusLayerFromNode(node)
  const layer = resolveNavLayer(modbusLayer, node.kind)
  const childNodes = collectChildNodesFromTree(node)
  const isDevList = isPureDeviceContainer(node)
  const base = {
    layer,
    modbusLayer,
    kind: resolveLegacyKind(layer),
    treeDepth: node.treeDepth,
    sid: node.sid != null ? node.sid : null,
    uuid: node.uuid || '',
    name: node.label || node.name || '',
    label: node.label || node.name || '',
    modelUuid: node.modelUuid || node.muid || '',
    muid: node.modelUuid || node.muid || '',
    childNodes,
    children: node.children || [],
    ancestors,
    deviceUuid: '',
    datapointPageIndex: 0,
    datapointPageSize: DEFAULT_DATAPOINT_PAGE_SIZE,
    pageIndex: 0,
    pageSize: DEFAULT_DEVICE_PAGE_SIZE,
    totalCount: childNodes.length,
  }
  if (isDevList) {
    base.childDevices = childNodes.map(c => ({
      name: c.name,
      uuid: c.uuid,
      code: c.uuid,
      modelUuid: c.modelUuid,
      status: c.status || 'off',
    }))
    base.deviceListMode = true
  } else {
    base.childDevices = []
  }
  return base
}

/** 设备列表分页（UPS 122 等） */
export function applyGatewayListPagination(nav, pageSizeOverride) {
  const fromDevices = (nav.allChildDevices || nav.childDevices || []).map(d => ({
    ...d,
    // 保留虚拟列头柜 kind；仅缺省时才回落为 device
    kind: d.kind || (d.virtualCabinet ? 'virtualCabinet' : 'device'),
    layer: d.layer || 'device',
  }))
  const all = nav.allChildNodes
    || nav.childNodes
    || fromDevices
  const size = pageSizeOverride || nav.pageSize || DEFAULT_DEVICE_PAGE_SIZE
  const p = paginateDevices(all, nav.pageIndex, size)
  return {
    ...nav,
    layer: nav.layer || 'zone',
    deviceListMode: true,
    // 虚拟柜列表模式标记不得被冲掉
    virtualCabinetListMode: !!(nav.virtualCabinetListMode
      || all.some(d => d && (d.kind === 'virtualCabinet' || d.virtualCabinet))),
    allChildNodes: all,
    childNodes: p.pageDevices,
    allChildDevices: all,
    childDevices: p.pageDevices,
    pageIndex: p.pageIndex,
    pageSize: p.pageSize,
    totalCount: p.totalDevices,
    totalDevices: p.totalDevices,
    totalPages: p.totalPages,
  }
}

/** 测点列表分页（信号层 ViewRealTable） */
export function applyDatapointPagination(nav, pageSizeOverride) {
  // 服务端已完成筛选与分页时，浏览器只保留当前页，不得再次切片或恢复全量缓存。
  if (nav && nav.serverPaged) {
    const size = normalizePageSize(pageSizeOverride || nav.datapointPageSize)
    const total = Math.max(0, Number(nav.totalDatapoints) || 0)
    const totalPages = Math.max(1, Math.ceil(total / size))
    const pageIndex = Math.min(Math.max(0, Number(nav.datapointPageIndex) || 0), totalPages - 1)
    return {
      ...nav,
      datapointPageIndex: pageIndex,
      datapointPageSize: size,
      datapointTotalPages: totalPages,
      detailPointMode: total > size,
    }
  }
  const all = nav.allDatapoints || nav.datapoints || nav.childNodes || []
  const size = pageSizeOverride || nav.datapointPageSize || DEFAULT_DATAPOINT_PAGE_SIZE
  const p = paginateDevices(all, nav.datapointPageIndex, size)
  return {
    ...nav,
    layer: 'device',
    kind: 'device',
    signalMode: true,
    allDatapoints: all,
    datapoints: p.pageDevices,
    childNodes: p.pageDevices,
    datapointPageIndex: p.pageIndex,
    datapointPageSize: p.pageSize,
    totalDatapoints: p.totalDevices,
    datapointTotalPages: p.totalPages,
    detailPointMode: p.totalDevices > size,
  }
}

/** 是否处于列表分页模式（组织层设备列表 / 虚拟列头柜列表） */
export function isNavListPaged(nav) {
  if (!nav) return false
  return !!(nav.deviceListMode || nav.virtualCabinetListMode)
}

/** 是否信号层测点分页 */
export function isSignalDatapointPaged(nav) {
  if (!nav) return false
  return !!(nav.signalMode && (nav.totalDatapoints || 0) > (nav.datapointPageSize || DEFAULT_DATAPOINT_PAGE_SIZE))
}

/** 页码文案（组织层设备列表） */
export function formatNavPageInfo(nav) {
  const cur = (nav.pageIndex || 0) + 1
  const total = nav.totalPages || 1
  const n = nav.totalCount != null ? nav.totalCount : (nav.totalDevices || 0)
  return `第 ${cur}/${total} 页 · 共 ${n} 台设备`
}

/** 测点页码文案 */
export function formatDatapointPageInfo(nav) {
  const cur = (nav.datapointPageIndex || 0) + 1
  const total = nav.datapointTotalPages || 1
  const n = nav.totalDatapoints || 0
  return `第 ${cur}/${total} 页 · 共 ${n} 个测点`
}

/** 测点全名 → 表格行显示名（去掉设备前缀） */
export function datapointRowLabel(fullName, deviceLabel = '') {
  const n = String(fullName || '').trim()
  if (!n) return ''
  const label = String(deviceLabel || '').trim()
  // 对齐数据仓库：最后一个 `_` 拆设备名/测点名
  const split = splitNameByLastUnderscore(n)
  if (label && split.deviceName === label && split.pointName) {
    return split.pointName
  }
  if (label) {
    const prefix = `${label}_`
    if (n.startsWith(prefix)) return n.slice(prefix.length)
    if (n.startsWith(label)) return n.slice(label.length).replace(/^[_\s]+/, '') || n
  }
  return split.pointName || n
}

/**
 * 按数据仓库设备名列过滤测点（last `_`）。
 * @param {object[]} points
 * @param {string} deviceLabel 虚拟设备名 / 逻辑设备名
 * @param {{ includeUnprefixed?: boolean }} [opts] includeUnprefixed=true 时纳入无设备名前缀点（UPS fallback 组）
 */
export function filterDatapointsForDevice(points, deviceLabel = '', opts = {}) {
  const list = Array.isArray(points) ? points : []
  const label = String(deviceLabel || '').trim()
  if (!label || !list.length) return list
  const includeUnprefixed = !!(opts && opts.includeUnprefixed)

  const nameOf = p => String(p.name || p.Name || p.label || '')

  // 1) 对齐数据仓库：split.deviceName === label；fallback 组另含无前缀点
  let filtered = list.filter(p => {
    const split = splitNameByLastUnderscore(nameOf(p))
    if (split.deviceName === label) return true
    if (includeUnprefixed && !split.deviceName) return true
    return false
  })
  if (filtered.length) return filtered

  // 2) 兼容旧 startsWith（单段前缀 / 共享模型短名）
  filtered = list.filter(p => {
    const n = nameOf(p)
    return n === label || n.startsWith(`${label}_`)
  })
  if (filtered.length) return filtered

  // 3) 去空白后再试（树节点名偶发带空格）
  const compact = label.replace(/\s+/g, '')
  if (compact && compact !== label) {
    filtered = list.filter(p => {
      const n = nameOf(p)
      const split = splitNameByLastUnderscore(n)
      if (split.deviceName === compact) return true
      return n === compact || n.startsWith(`${compact}_`)
    })
    if (filtered.length) return filtered
  }

  // 4) 仍无结果：不回退全量（共享模型 1800+ 会撑爆表格）
  return []
}

const _modelPointsCache = Object.create(null)
const _devicePointsCache = Object.create(null)

export function clearDeviceDatapointCache(muid) {
  if (muid) {
    delete _modelPointsCache[muid]
    Object.keys(_devicePointsCache).forEach(k => {
      if (k.startsWith(`${muid}|`)) delete _devicePointsCache[k]
    })
    return
  }
  Object.keys(_modelPointsCache).forEach(k => { delete _modelPointsCache[k] })
  Object.keys(_devicePointsCache).forEach(k => { delete _devicePointsCache[k] })
}

function extractModelPointsList(res) {
  const body = res && res.data
  if (!body) return []
  if (Array.isArray(body.list)) return body.list
  if (body.data && Array.isArray(body.data.list)) return body.data.list
  if (Array.isArray(body.data)) return body.data
  return []
}

function mapRealDataRows(rows) {
  return (Array.isArray(rows) ? rows : []).map(p => ({
    name: p.name || p.Name || '',
    label: p.name || p.Name || '',
    uuid: p.uuid || p.Uuid || '',
    unit: p.unit || p.Unit || p.DataUnit || p.data_unit || p.dataUnit || '',
    value: p.value != null ? p.value : (p.Value != null ? p.Value : ''),
    deviceName: p.DeviceName || p.device_name || p.deviceName || '',
    deviceUuid: p.duid || p.device_uuid || p.deviceUuid || '',
    modelDataUuid: p.mduid || p.model_data_uuid || p.ModelDataUuid || '',
  })).filter(p => {
    const n = String(p.name || '').trim()
    if (!n) return false
    // 系统内置点，不进测点表（如 device.DeviceStatus）
    if (/^device\./i.test(n)) return false
    if (/^system\./i.test(n)) return false
    return true
  })
}

/**
 * 拉取设备测点 —— 与「设备管理/数据仓库」同源：
 * - 名称、单位：device_real_data（库）
 * - 实时值：DeviceRealDataMapByUUID（内存，socket/采集写入；GetRealData 已覆盖）
 *
 * @param {string} muid 物模型 uuid
 * @param {string} deviceLabel 逻辑设备名（共享模型下按 name 前缀过滤）
 * @param {string} [deviceUuid] 设备 uuid（优先）
 */
export async function fetchDeviceDatapoints(muid, deviceLabel = '', deviceUuid = '') {
  const label = String(deviceLabel || '').trim()
  const uuid = String(deviceUuid || '').trim()
  const cacheKey = `${muid || ''}|${label}|${uuid}`
  if (Object.prototype.hasOwnProperty.call(_devicePointsCache, cacheKey)
    && Array.isArray(_devicePointsCache[cacheKey])
    && _devicePointsCache[cacheKey].length > 0) {
    return _devicePointsCache[cacheKey]
  }

  // 1) 主路径：分页拉取（禁止 fetchAll，大设备可达数千～数万点）
  // 有逻辑设备名时优先只按 namePrefix 拉（共享模型全名点），避免 OR device_uuid 混入短名/系统点
  if (uuid || label) {
    try {
      const pageSize = 200
      const collected = []
      let page = 1
      let total = Infinity
      const maxPages = 100 // 硬顶约 20000，避免大机柜列尾被截断
      while (page <= maxPages && collected.length < total) {
        const primary = label
          ? {
            muid: muid || undefined,
            namePrefix: label,
            deviceLabel: label,
            uuid: uuid || undefined,
            page,
            pageSize,
            IsRemoveGW: false,
          }
          : {
            uuid: uuid || undefined,
            muid: muid || undefined,
            page,
            pageSize,
            IsRemoveGW: false,
          }
        const res = await getRealData(primary)
        const body = res && res.data
        if (!body || body.code !== 0) break
        const rows = mapRealDataRows(body.realData)
        if (!rows.length) break
        collected.push(...rows)
        total = Math.max(0, Number(body.total) || collected.length)
        const hasMore = body.hasMore === true || collected.length < total
        if (!hasMore && rows.length < pageSize) break
        if (collected.length >= total) break
        page += 1
      }
      let list = collected
      // 前缀为空时再回退设备 uuid（同样分页）
      if (!list.length && uuid && label) {
        page = 1
        total = Infinity
        while (page <= maxPages && list.length < total) {
          const res2 = await getRealData({
            uuid,
            page,
            pageSize,
            IsRemoveGW: false,
          })
          const body2 = res2 && res2.data
          if (!body2 || body2.code !== 0) break
          const rows2 = mapRealDataRows(body2.realData)
          if (!rows2.length) break
          list = list.concat(rows2)
          total = Math.max(0, Number(body2.total) || list.length)
          if (!body2.hasMore && rows2.length < pageSize) break
          if (list.length >= total) break
          page += 1
        }
      }
      if (list.length) {
        _devicePointsCache[cacheKey] = list
        return list
      }
    } catch (e) {
      console.warn('[navContext] fetchDeviceDatapoints getRealData failed', e && e.message)
    }
  }

  // 2) 回退：物模型测点表（仅名称，无单位/实时值）
  if (!muid) return []
  try {
    let raw = _modelPointsCache[muid]
    if (!raw) {
      const res = await getModelDataPoints({ muid })
      raw = extractModelPointsList(res)
      if (!Array.isArray(raw)) raw = []
      if (raw.length) {
        _modelPointsCache[muid] = raw
      }
    }
    const list = raw.map(p => ({
      name: p.name || p.Name || '',
      label: p.name || p.Name || '',
      uuid: p.uuid || p.Uuid || '',
      unit: p.unit || p.Unit || p.dataUnit || p.data_unit || '',
      value: '',
    })).filter(p => p.name)
    let filtered = filterDatapointsForDevice(list, label)
    if (!label) filtered = list
    if (filtered.length) {
      _devicePointsCache[cacheKey] = filtered
    } else {
      console.warn('[navContext] fetchDeviceDatapoints empty after filter', {
        muid, label, uuid, modelPoints: list.length,
      })
    }
    return filtered
  } catch (e) {
    console.warn('[navContext] fetchDeviceDatapoints failed', muid, e && e.message)
    return []
  }
}

/**
 * 刷新当前页测点的内存实时值（名称/单位不变）
 * 与设备管理一致：按 uuid 列表走 GetRealData namePrefix 分页，或整表 fetchAll 后切片
 */
export async function refreshDeviceDatapointValues(points = []) {
  const list = Array.isArray(points) ? points : []
  if (!list.length) return list
  // 已有 value 字段时，用 GetRealDataByBindings 不合适；改为按测点 uuid 无法直接批量
  // 这里返回原列表，由调用方用 getRealData(fetchAll) 重拉；保留 hook 便于扩展
  return list
}

/** 由 nav 生成当前页测点表行配置（供页内翻页复用，避免整页 GoPage） */
export function buildSignalTablePageConfig(nav) {
  if (!nav) {
    return {
      rowDeviceNames: [],
      rowDeviceCodes: [],
      columnHeaders: ['实时值'],
      bindingMatrix: [],
      totalDatapoints: 0,
      datapointPageIndex: 0,
      datapointPageSize: DEFAULT_DATAPOINT_PAGE_SIZE,
      datapointTotalPages: 1,
      deviceLabel: '',
    }
  }
  const paged = applyDatapointPagination(nav)
  const deviceLabel = paged.name || paged.label || ''
  const points = paged.datapoints || paged.childNodes || []
  const bindingMatrix = points
    .map(p => {
      const pointName = p.name || p.label || ''
      if (!pointName) return null
      // 共享模型实点多挂在网关：优先 owner->name，否则仅测点全名（后端按 dataName 解析）
      const owner = p.deviceName || p.device_name || ''
      if (owner) return [`${owner}->${pointName}`]
      return [pointName]
    })
    .filter(Boolean)
  return {
    rowDeviceNames: points.map(p => datapointRowLabel(p.name || p.label, deviceLabel)),
    rowDeviceCodes: points.map(p => p.unit || ''),
    columnHeaders: ['实时值'],
    bindingMatrix,
    totalDatapoints: paged.totalDatapoints || 0,
    datapointPageIndex: paged.datapointPageIndex || 0,
    datapointPageSize: paged.datapointPageSize || DEFAULT_DATAPOINT_PAGE_SIZE,
    datapointTotalPages: paged.datapointTotalPages || 1,
    deviceLabel,
    pagedNav: paged,
  }
}

/**
 * 信号层：设备 → 测点 navContext
 * @param {object} deviceNode 设备树节点
 * @param {object[]} [datapoints] 已拉取的测点（可空，由 binding 异步补全）
 * @param {object[]} [ancestors]
 */
export function buildDeviceSignalContext(deviceNode, datapoints = [], ancestors = [], pageInfo = {}) {
  const points = datapointsToChildNodes(datapoints)
  const pageSize = normalizePageSize(pageInfo.pageSize || DEFAULT_DATAPOINT_PAGE_SIZE)
  const total = Math.max(points.length, Number(pageInfo.total) || 0)
  const vc = deviceNode.virtualCabinet || pageInfo.virtualCabinet || ''
  const parentLabel = deviceNode.parentDeviceLabel || pageInfo.parentDeviceLabel || ''
  const vcFallback = !!(
    deviceNode.virtualCabinetFallback
    || deviceNode.isFallbackGroup
    || pageInfo.virtualCabinetFallback
    || (vc && parentLabel && vc === parentLabel)
  )
  const base = {
    layer: 'device',
    modbusLayer: 'device',
    kind: 'device',
    routeMode: 'signal',
    signalMode: true,
    treeDepth: deviceNode.treeDepth,
    sid: deviceNode.sid != null ? deviceNode.sid : null,
    uuid: deviceNode.uuid || '',
    deviceUuid: deviceNode.uuid || '',
    gatewayUuid: deviceNode.uuid || '',
    name: deviceNode.label || deviceNode.name || '',
    label: deviceNode.label || deviceNode.name || '',
    code: deviceNode.code || deviceNode.uuid || '',
    status: deviceNode.status || 'off',
    modelUuid: deviceNode.modelUuid || deviceNode.muid || '',
    muid: deviceNode.modelUuid || deviceNode.muid || '',
    childNodes: points,
    datapoints: points,
    allDatapoints: pageInfo.serverPaged ? [] : points,
    children: [],
    ancestors,
    datapointPageIndex: Math.max(0, Number(pageInfo.page || 1) - 1),
    datapointPageSize: pageSize,
    totalDatapoints: total,
    serverPaged: !!pageInfo.serverPaged,
    datapointQuery: String(pageInfo.query || ''),
    datapointCategory: String(pageInfo.category || (vc && !vcFallback ? `${vc}_` : '') || ''),
    virtualCabinet: vc,
    virtualCabinetFallback: vcFallback,
    parentDeviceLabel: parentLabel,
    parentDeviceUuid: deviceNode.parentDeviceUuid || pageInfo.parentDeviceUuid || '',
    deviceListReturnContext: deviceNode.deviceListReturnContext || null,
    homePageUuid: deviceNode.homePageUuid || pageInfo.homePageUuid || '',
    childDevices: [],
  }
  return applyDatapointPagination(base)
}

/**
 * 仅请求当前设备测点的当前筛选页。
 * 不得使用 fetchAll：设备可能包含数千到数万测点。
 *
 * @param {string} [pointNamePrefix] 虚拟设备名（数据仓库 deviceName / last `_` 前段）。
 *   后端 uuid+namePrefix 是 OR（会拉回整台设备测点），虚拟柜必须走 uuid+category AND。
 * @param {boolean} [isFallbackGroup] UPS 等无前缀点归属真设备名时：uuid 拉取后按 last `_` 客户端过滤
 */
export async function fetchDeviceDatapointPage({
  muid = '',
  deviceLabel = '',
  deviceUuid = '',
  page = 1,
  pageSize = DEFAULT_DATAPOINT_PAGE_SIZE,
  query = '',
  category = '',
  pointNamePrefix = '',
  isFallbackGroup = false,
} = {}) {
  const prefix = String(pointNamePrefix || '').trim()
  const uuid = String(deviceUuid || '').trim()
  const label = String(deviceLabel || '').trim()
  const cat = String(category || '').trim()
  const pageNum = Math.max(1, Number(page) || 1)
  const size = normalizePageSize(pageSize)
  const q = String(query || '').trim() || undefined
  const fallback = !!isFallbackGroup

  // fallback 组含无 `_` 点位，category LIKE '设备名_%' 会漏点 → 按 uuid 拉页再客户端过滤
  if (fallback && prefix && uuid) {
    const collected = []
    let scanPage = 1
    let total = Infinity
    const maxPages = 50
    while (scanPage <= maxPages && collected.length < total) {
      const res = await getRealData({
        uuid,
        muid: muid || undefined,
        page: scanPage,
        pageSize: 100,
        query: q,
        IsRemoveGW: false,
      })
      const body = res && res.data
      if (!body || body.code !== 0) break
      const rows = mapRealDataRows(body.realData)
      if (!rows.length) break
      collected.push(...rows)
      total = Math.max(0, Number(body.total) || collected.length)
      if (!body.hasMore && rows.length < 100) break
      if (collected.length >= total) break
      scanPage += 1
    }
    let filtered = filterDatapointsForDevice(collected, prefix, { includeUnprefixed: true })
    if (q) {
      const qq = q.toLowerCase()
      filtered = filtered.filter(p => String(p.name || p.label || '').toLowerCase().includes(qq))
    }
    const start = (pageNum - 1) * size
    return {
      points: filtered.slice(start, start + size),
      total: filtered.length,
      page: pageNum,
      pageSize: size,
    }
  }

  // 虚拟设备：device_uuid = ? AND name LIKE '设备名_%'（用 category，避免 OR 全量）
  const payload = (prefix && uuid)
    ? {
      uuid,
      muid: muid || undefined,
      page: pageNum,
      pageSize: size,
      query: q,
      category: prefix.endsWith('_') ? prefix : `${prefix}_`,
      IsRemoveGW: false,
    }
    : {
      uuid: uuid || undefined,
      muid: muid || undefined,
      namePrefix: label || undefined,
      deviceLabel: label || undefined,
      page: pageNum,
      pageSize: size,
      query: q,
      category: cat || undefined,
      IsRemoveGW: false,
    }

  const res = await getRealData(payload)
  const body = res && res.data
  if (!body || body.code !== 0) {
    return { points: [], total: 0, page: 1, pageSize: size }
  }
  let points = mapRealDataRows(body.realData)
  // 客户端兜底：按 last `_` 设备名对齐，避免误 OR
  if (prefix) {
    points = filterDatapointsForDevice(points, prefix, { includeUnprefixed: false })
  }
  return {
    points,
    total: Math.max(0, Number(body.total) || 0),
    page: Math.max(1, Number(body.page) || 1),
    pageSize: normalizePageSize(body.pageSize || size),
  }
}

/** 从 navContext 推导测点分页请求参数（虚拟柜用 pointNamePrefix） */
export function resolveDatapointFetchParams(nav) {
  const n = nav || {}
  const vc = String(n.virtualCabinet || '').trim()
  const uuid = String(n.deviceUuid || n.uuid || '').trim()
  const muid = String(n.modelUuid || n.muid || '').trim()
  const parentLabel = String(n.parentDeviceLabel || '').trim()
  const isFallback = !!(
    n.virtualCabinetFallback
    || n.isFallbackGroup
    || (vc && parentLabel && vc === parentLabel)
  )
  if (vc && uuid) {
    return {
      muid,
      deviceUuid: uuid,
      pointNamePrefix: vc,
      deviceLabel: '',
      category: '',
      isFallbackGroup: isFallback,
    }
  }
  return {
    muid,
    deviceUuid: uuid,
    pointNamePrefix: '',
    deviceLabel: String(n.label || n.name || '').trim(),
    category: String(n.datapointCategory || '').trim(),
    isFallbackGroup: false,
  }
}

/** @deprecated 寄存器组主链路已禁用 */
export async function fetchRegisterGroups() {
  return []
}

/** @deprecated */
export async function fetchRegisterGroupPoints() {
  return []
}

/** @deprecated */
export function buildGatewayRegisterGroupContext(gatewayNode, _groups, ancestors = []) {
  return buildDeviceSignalContext(gatewayNode, [], ancestors)
}

/** @deprecated */
export function buildRegisterGroupDetailContext(parentNav, registerGroup, ancestors = []) {
  return buildDeviceSignalContext(
    { uuid: parentNav.deviceUuid || parentNav.uuid, label: registerGroup.name, modelUuid: parentNav.modelUuid },
    [],
    ancestors,
  )
}

/** @deprecated */
export function applyRegisterGroupPagination(nav, pageSizeOverride) {
  return applyDatapointPagination(nav, pageSizeOverride)
}

/** @deprecated */
export function registerGroupsToChildNodes(groups) {
  return datapointsToChildNodes(groups)
}

/** @deprecated A2 中间分组命名 */
export function isA2GroupName(name) {
  const n = String(name || '').trim()
  return /配电室_机房模块|_机房模块/.test(n)
}

/** @deprecated */
export function isA2GroupContainer(node) {
  if (!node) return false
  if (isA2GroupName(node.rawLabel || node.label || node.name)) return true
  const children = (node && node.children) || []
  if (!children.length) return false
  if (isPureDeviceContainer(node)) return false
  return children.some(c => c && (c.kind === 'zone' || c.kind === 'room'))
}
