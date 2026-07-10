/**
 * 导航树索引：
 * - uuid5Hex 旧页 page_id 生成（与 build_ncc_dashboard.py 完全一致）
 * - 旧 page_id → {oldKind, sid, key} 反查表（模板 cells 里的 click link 用它做「槽位钥匙」）
 * - navContext 构造（含祖先链 ancestors，用于面包屑重写）
 */

import {
  DEFAULT_DEVICE_PAGE_SIZE,
  resolveDeviceListTemplateId,
} from './deviceListPager'
import {
  createBaseNavContext,
  applyGatewayListPagination,
  buildDeviceSignalContext,
  isPureDeviceContainer,
} from './navContext'
import { detectDrillDepth, isLeafDeviceContainer, isDeviceNode } from './drillDepth'
import sha1 from 'crypto-js/sha1'
import Hex from 'crypto-js/enc-hex'

// uuid5(DNS, name).hex
const DNS_NS_HEX = '6ba7b8109dad11d180b400c04fd430c8'

function utf8Bytes(s) {
  return Array.from(unescape(encodeURIComponent(s))).map(c => c.charCodeAt(0))
}
function bytesToHex(b) {
  return b.map(x => x.toString(16).padStart(2, '0')).join('')
}
function hexToBytes(h) {
  const a = []
  for (let i = 0; i < h.length; i += 2) a.push(parseInt(h.substr(i, 2), 16))
  return a
}
export function uuid5Hex(name) {
  const allHex = DNS_NS_HEX + bytesToHex(utf8Bytes(name))
  const hashHex = sha1(Hex.parse(allHex)).toString(Hex)
  const b = hexToBytes(hashHex).slice(0, 16)
  b[6] = (b[6] & 0x0f) | 0x50
  b[8] = (b[8] & 0x3f) | 0x80
  return bytesToHex(b)
}

export const pageIdBuilding = sid => uuid5Hex(`ncc-dash-bldg-${sid}`)
export const pageIdFloor = (bldgSid, key) => uuid5Hex(`ncc-dash-floor-${bldgSid}-${key}`)
export const pageIdDevice = sid => uuid5Hex(`ncc-dash-dev-${sid}`)
export const pageIdRoom = sid => uuid5Hex(`ncc-dash-room-${sid}`)
export const pageIdZone = sid => uuid5Hex(`ncc-dash-zone-${sid}`)

/** 设备名 → 设备组 key（与 build_ncc_dashboard.py 分组规则一致） */
export function floorKey(name) {
  const parts = String(name || '').split('_')
  return parts.length >= 3 ? parts[2] : 'default'
}

export function sortByLabel(list, field = 'label') {
  return [...(list || [])].sort((a, b) =>
    String(a[field] || '').localeCompare(String(b[field] || ''), 'en', {
      numeric: true,
      sensitivity: 'base',
    }),
  )
}

/**
 * 遍历导航树（ISMRunTreeNav.buildForest 产物）构建索引。
 * @returns {{
 *   bySid: Object<number, object>,
 *   parentBySid: Object<number, object|null>,
 *   byOldPageId: Object<string, {oldKind:string, sid:number, key?:string}>,
 *   roots: object[],
 * }}
 */
export function buildNavTreeIndex(roots) {
  const bySid = Object.create(null)
  const parentBySid = Object.create(null)
  const byOldPageId = Object.create(null)

  const put = (pid, entry) => {
    if (!byOldPageId[pid]) byOldPageId[pid] = entry
  }

  const walk = (node, parent) => {
    if (!node) return
    const sid = node.sid
    if (sid != null) {
      bySid[sid] = node
      parentBySid[sid] = parent || null
      if (node.kind === 'device' || node.kind === 'gateway') {
        put(pageIdDevice(sid), { oldKind: 'device', sid })
      } else {
        // 旧脚本对同一容器 sid 可能生成 zone/room/building 三种页，全部注册
        put(pageIdZone(sid), { oldKind: 'zone', sid })
        put(pageIdRoom(sid), { oldKind: 'room', sid })
        put(pageIdBuilding(sid), { oldKind: 'building', sid })
        // 设备组中间层（仅模板页钻探用；导航树不再生成 kind=floor 节点）
        const floorKids = (node.children || []).filter(c => c && c.kind === 'floor')
        floorKids.forEach(f => {
          const key = f.floorKey || floorKey(f.label)
          const bldgSid = f.cabinetSid != null ? f.cabinetSid : sid
          put(pageIdFloor(bldgSid, key), { oldKind: 'floor', sid: bldgSid, key })
        })
        // 直属设备的分组页 floor-{sid}-{key}（兼容旧平铺结构）
        const direct = (node.children || []).filter(c => c && (c.kind === 'device' || c.kind === 'gateway'))
        const keys = new Set(direct.map(d => floorKey(d.label || d.name)))
        keys.forEach(key => {
          put(pageIdFloor(sid, key), { oldKind: 'floor', sid, key })
        })
      }
    }
    ;(node.children || []).forEach(c => walk(c, node))
  }
  ;(roots || []).forEach(r => walk(r, null))

  return { bySid, parentBySid, byOldPageId, roots: roots || [] }
}

/** 深收集节点子孙 A3 转机（带在线状态） */
export function collectNodeDevices(node) {
  const out = []
  const walk = n => {
    if (!n) return
    if (n.kind === 'gateway' || n.kind === 'device') {
      out.push({
        name: n.label || n.name || '',
        uuid: n.uuid || '',
        code: n.uuid || '',
        modelUuid: n.modelUuid || n.muid || '',
        status: n.status || 'off',
        kind: 'device',
        layer: 'device',
      })
      return
    }
    ;(n.children || []).forEach(walk)
  }
  if (node && (node.kind === 'gateway' || node.kind === 'device' || isDeviceNode(node))) {
    walk(node)
  } else {
    ;((node && node.children) || []).forEach(walk)
  }
  return out
}

/** ancestors 链的轻量元素（不含 children 树，避免序列化膨胀） */
function liteAncestor(nav) {
  const out = {
    kind: nav.kind,
    sid: nav.sid != null ? nav.sid : null,
    uuid: nav.uuid || '',
    label: nav.label || nav.name || '',
    name: nav.name || nav.label || '',
  }
  // 虚拟层级（直属设备组/设备组，无 sid）需要自带设备列表才能被面包屑重建
  if (out.sid == null && Array.isArray(nav.childDevices)) {
    out.childDevices = nav.childDevices
  }
  return out
}

/**
 * 为树节点构造 navContext（含祖先链）。
 * @param {object} node 树节点
 * @param {object} index buildNavTreeIndex 产物
 * @param {object[]} [ancestorsOverride] 显式祖先链（虚拟层级下钻时用）
 */
export function buildNavContextForNode(node, index, ancestorsOverride) {
  let ancestors = ancestorsOverride
  if (!ancestors) {
    ancestors = []
    if (index && node) {
      if (node.sid != null) {
        const chain = []
        let p = index.parentBySid[node.sid]
        while (p) {
          chain.unshift(p)
          p = p.sid != null ? index.parentBySid[p.sid] : null
        }
        ancestors = chain.map(n => liteAncestor({
          kind: n.kind,
          layer: n.layer || n.kind,
          modbusLayer: n.modbusLayer,
          sid: n.sid,
          uuid: n.uuid,
          label: n.label,
          name: n.label,
          treeDepth: n.treeDepth,
        }))
      }
    }
  }
  const base = createBaseNavContext(node, ancestors)
  const depth = detectDrillDepth(node, index)
  const enriched = {
    ...base,
    routeMode: depth.routeMode,
    drillDepth: depth.totalLayers,
    orgDepth: depth.orgDepth,
    treeDepth: depth.treeDepth,
    remainingLayers: depth.remainingLayers,
    modbusLayer: depth.modbusLayer || base.modbusLayer,
  }
  const pureList = !isDeviceNode(node) && isPureDeviceContainer(node)
  if (pureList) {
    return applyGatewayListPagination({
      ...enriched,
      allChildNodes: enriched.childNodes,
      allChildDevices: enriched.childDevices,
      pageIndex: 0,
      pageSize: DEFAULT_DEVICE_PAGE_SIZE,
    })
  }
  if (isDeviceNode(node) || isLeafDeviceContainer(node)) {
    return buildDeviceSignalContext(node, [], ancestors)
  }
  return enriched
}

/** 子 navContext（槽位下钻）：ancestors = 父链 + 父自身 */
export function childAncestors(parentNav) {
  return [...(parentNav.ancestors || []), liteAncestor(parentNav)]
}

/** 按节点 kind/layer/物模型解析模板页 id */
export function resolveTemplatePageIdForKind(templateMap, kind, modelUuid) {
  const map = templateMap || {}
  const roomFallback = map.room || map.zone || ''
  if (kind === 'root' || kind === 'home') return map.home || roomFallback
  if (kind === 'zone') return map.zone || roomFallback
  if (kind === 'room') return map.room || map.zone || ''
  // 设备（含原 gateway）→ 测点实时表模板
  if (kind === 'gateway' || kind === 'device') {
    const byModel = map.deviceByModel || {}
    if (modelUuid && byModel[modelUuid]) return byModel[modelUuid]
    return map.deviceDefault || map.cabinet || roomFallback
  }
  if (kind === 'cabinet') return map.cabinet || map.room || map.zone || ''
  // 禁用 floor / registerGroup 主链路
  if (kind === 'floor' || kind === 'registerGroup') return map.room || map.zone || roomFallback
  if (kind === 'device') {
    const byModel = map.deviceByModel || {}
    if (modelUuid && byModel[modelUuid]) return byModel[modelUuid]
    return map.deviceDefault || roomFallback
  }
  return roomFallback
}

/**
 * GoPage 兜底：把旧 page_id 转成「模板页 + navContext」。
 * 仅用于链接未经槽位重映射的场景（如首页初载后的卡片点击）。
 * @returns {{pageUuid: string, navContext: object}|null}
 */
export function resolveOldPageTarget(oldPageId, index, templateMap) {
  if (!oldPageId || !index || !index.byOldPageId) return null
  const entry = index.byOldPageId[oldPageId]
  if (!entry) return null
  const node = index.bySid[entry.sid]
  if (!node) return null

  if (entry.oldKind === 'floor') {
    // 旧 floor 设备组 → room 列表模板 + A3 转机分页（禁用 floor 主链路）
    const nav = applyGatewayListPagination(buildNavContextForNode(node, index))
    const pageUuid = resolveDeviceListTemplateId(templateMap)
      || resolveTemplatePageIdForKind(templateMap, node.kind, nav.modelUuid)
    return pageUuid ? { pageUuid, navContext: nav } : null
  }

  const nav = buildNavContextForNode(node, index)
  const pageUuid = resolveTemplatePageIdForKind(templateMap, node.kind, nav.modelUuid)
  return pageUuid ? { pageUuid, navContext: nav } : null
}
