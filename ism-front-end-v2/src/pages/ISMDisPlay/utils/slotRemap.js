/**
 * 容器模板页「槽位重映射」：
 * 模板页 cells 中 click link 仍指向样本时代的旧 page_id（已软删除）。
 * 以旧 page_id 作为「槽位钥匙」（byOldPageId 反查得 {oldKind, sid, key}），
 * 将 cells 分组为槽位，与当前 navContext 的子节点/设备组对齐后逐槽重写：
 *   - click link → 对应层级模板页 + 子 navContext（内嵌在 link.navContext）
 *   - 悬浮/数据绑点 deviceSN → 对应子设备 uuid
 *   - 槽位静态文本（组名/台数/状态统计）按当前子节点改写
 * 多余槽位整组删除；面包屑链接重定向到当前祖先。
 */

import {
  buildNavContextForNode,
  childAncestors,
  collectNodeDevices,
  resolveTemplatePageIdForKind,
  sortByLabel,
} from './navTreeIndex'
import {
  applyDeviceListPagination,
  inferPageSizeFromTargets,
  rewritePaginationCells,
  isDeviceListNav,
  isFloorGroupArtifact,
  resolveDeviceListTemplateId,
} from './deviceListPager'
import {
  applyGatewayListPagination,
  buildDeviceSignalContext,
  treeChildToNodeItem,
} from './navContext'
import { isDeviceNode } from './drillDepth'
import { sanitizeGraphComponents } from './graphCellSanitizer'

// ---------- 基础工具 ----------

function getClickLink(cell) {
  const detail = cell && cell.data && cell.data.detail
  if (!detail || !Array.isArray(detail.action)) return null
  const act = detail.action.find(
    a => a && a.type === 'click' && a.action === 'link' && a.link && a.link.Inside,
  )
  return act ? act.link : null
}

function setCellLink(cell, pageUuid, navContext) {
  const link = getClickLink(cell)
  if (!link) return
  link.Inside.pageUUID = pageUuid
  link.navContext = navContext
}

function getText(cell) {
  const style = cell && cell.data && cell.data.detail && cell.data.detail.style
  return style && typeof style.text === 'string' ? style.text : ''
}

function setText(cell, text) {
  const style = cell && cell.data && cell.data.detail && cell.data.detail.style
  if (style) style.text = text
}

/** 悬浮绑定：写入 animate.condition（ViewSvgText.resolveDeviceBinding 读取的位置） */
function setHoverDevice(cell, uuid, name) {
  const detail = cell && cell.data && cell.data.detail
  if (!detail || !uuid) return
  if (!detail.animate) detail.animate = {}
  if (!detail.animate.condition) {
    detail.animate.condition = { deviceSN: '', isBandDevice: false, dataID: '', dataName: '' }
  }
  detail.animate.condition.deviceSN = uuid
  detail.animate.condition.DeviceName = name || ''
}

function walkConds(obj, visit) {
  if (!obj || typeof obj !== 'object') return
  if (Array.isArray(obj)) {
    obj.forEach(v => walkConds(v, visit))
    return
  }
  if (obj.condition && typeof obj.condition === 'object') visit(obj.condition)
  Object.keys(obj).forEach(k => {
    if (k === 'condition') return
    const v = obj[k]
    if (v && typeof v === 'object') walkConds(v, visit)
  })
}

function deviceStats(devices) {
  const total = (devices || []).length
  const online = (devices || []).filter(d => d.status === 'on').length
  return { total, online, offline: total - online }
}

/** 样本 floor/cabinet 页里的「设备组」面包屑/标题（纯设备列表应剔除） */
function stripFloorGroupArtifactCells(cells, drop) {
  cells.forEach((cell, i) => {
    if (drop.has(i)) return
    const t = getText(cell)
    if (isFloorGroupArtifact(t)) drop.add(i)
  })
}

/** 设备/转机行内统计：单卡 1 条，禁止父级总数写入子卡 */
function rewriteDeviceRowTexts(cells, idxs, device) {
  const on = device && device.status === 'on'
  const online = on ? 1 : 0
  const offline = on ? 0 : 1
  const isRg = device && (device.kind === 'registerGroup' || device.layer === 'registerGroup')
  idxs.forEach(i => {
    const t = getText(cells[i])
    if (!t) return
    let m
    if (/^\d+台设备$/.test(t)) {
      setText(cells[i], isRg ? '1组' : (on ? '运行' : '离线'))
    } else if ((m = t.match(/^🟢 \d+(运行|在线)$/))) {
      setText(cells[i], `🟢 ${online}${m[1]}`)
    } else if (/^🔴 \d+告警$/.test(t)) {
      setText(cells[i], '🔴 0告警')
    } else if ((m = t.match(/^⏸ \d+停止$/))) {
      setText(cells[i], `⏸ ${offline}停止`)
    } else if (/^\d+$/.test(t)) {
      const next = idxs.indexOf(i) >= 0 && idxs.indexOf(i) + 1 < idxs.length
        ? getText(cells[idxs[idxs.indexOf(i) + 1]])
        : ''
      if (next === '设备') setText(cells[i], '1')
      else if (next === '在线' || next === '运行') setText(cells[i], String(online))
      else if (next === '异常' || next === '停止') setText(cells[i], String(offline))
    } else if (t === '离线' || t === '运行' || t === '在线') {
      setText(cells[i], on ? '运行' : '离线')
    } else if (isFloorGroupArtifact(t)) {
      setText(cells[i], '')
    }
  })
}

/** 纯设备列表页：改写返回面包屑为父容器名 */
function rewriteDeviceListBackLink(cells, idxs, nav) {
  const parent = (nav.ancestors || []).slice(-1)[0]
  const label = parent ? (parent.label || parent.name) : (nav.label || nav.name || '')
  if (!label) return
  idxs.forEach(i => {
    const t = getText(cells[i])
    if (!t) return
    if (/^←\s*/.test(t)) setText(cells[i], `← ${label}`)
    else if (t.length > 1 && !/^[›|#·]+$/.test(t) && !/台设备/.test(t)) {
      setText(cells[i], label)
    }
  })
}

// ---------- 槽位文本改写 ----------

const STAT_LABELS = new Set(['设备', '在线', '异常', '运行', '告警', '停止'])

/**
 * 槽位卡片文本按当前子节点改写。
 * @param {object[]} slotCells 槽位内 cells（按模板顺序）
 * @param {{title:string, devices:object[]}} info
 */
function rewriteSlotTexts(slotCells, info) {
  const { total, online, offline } = deviceStats(info.devices)
  for (let i = 0; i < slotCells.length; i++) {
    const t = getText(slotCells[i])
    if (!t) continue
    let m
    if (/^\d+台设备$/.test(t)) {
      setText(slotCells[i], `${total}台设备`)
    } else if ((m = t.match(/^🟢 \d+(运行|在线)$/))) {
      setText(slotCells[i], `🟢 ${online}${m[1]}`)
    } else if (/^🔴 \d+告警$/.test(t)) {
      // 无告警统计来源，维持 0
      setText(slotCells[i], '🔴 0告警')
    } else if ((m = t.match(/^⏸ \d+停止$/))) {
      setText(slotCells[i], `⏸ ${offline}停止`)
    } else if (/^\d+$/.test(t)) {
      // 纯数字：语义由槽位内紧随其后的标签 cell 决定（设备/在线/异常）
      const next = i + 1 < slotCells.length ? getText(slotCells[i + 1]) : ''
      if (next === '设备') setText(slotCells[i], String(total))
      else if (next === '在线' || next === '运行') setText(slotCells[i], String(online))
      else if (next === '异常' || next === '停止') setText(slotCells[i], String(offline))
    } else if (STAT_LABELS.has(t) || /^[›←—\-·|#]+$/.test(t)) {
      // 标签/符号不动
    } else if (/设备组$/.test(t.replace(/^📋\s*/, '')) || /设备列表$/.test(t)) {
      const prefix = t.startsWith('📋') ? '📋 ' : ''
      const title = /设备列表$/.test(info.title) ? info.title : `${info.title} · 设备列表`
      setText(slotCells[i], `${prefix}${title}`)
    } else if (t.length > 1 && !/^[🟢🔴⏸📋📊🏢🏛]/.test(t)) {
      // 其余非符号文本视为标题（样本子节点名）
      setText(slotCells[i], info.title)
    }
  }
}

/** 槽位整体重写：链接 + 悬浮 + 文本 */
function rewriteSlot(slotCells, pageUuid, navContext, title) {
  const devices = navContext.childDevices || []
  slotCells.forEach(cell => {
    if (pageUuid) setCellLink(cell, pageUuid, navContext)
    if (devices.length) setHoverDevice(cell, devices[0].uuid, devices[0].name)
  })
  rewriteSlotTexts(slotCells, { title, devices })
}

// ---------- 面包屑 ----------

/** 按旧页 kind 从当前祖先链中挑对应层级 */
function pickAncestor(nav, oldKind) {
  const list = nav.ancestors || []
  const match = kinds => {
    for (let i = list.length - 1; i >= 0; i--) {
      if (kinds.includes(list[i].kind)) return list[i]
    }
    return null
  }
  if (oldKind === 'building') return match(['cabinet']) || list[list.length - 1] || null
  if (oldKind === 'room') return match(['room', 'zone']) || list[list.length - 1] || null
  if (oldKind === 'zone') return match(['zone']) || null
  return list[list.length - 1] || null
}

/** 由祖先链元素重建 navContext（真实节点走树，虚拟层级用自带 childDevices） */
function ancestorNavContext(anc, nav, index) {
  const list = nav.ancestors || []
  const pos = list.indexOf(anc)
  const prefix = pos >= 0 ? list.slice(0, pos) : []
  if (anc.sid != null && index.bySid[anc.sid]) {
    return buildNavContextForNode(index.bySid[anc.sid], index, prefix)
  }
  return {
    sid: null,
    uuid: anc.uuid || '',
    name: anc.name || anc.label || '',
    label: anc.label || anc.name || '',
    kind: anc.kind,
    modelUuid: '',
    muid: '',
    childDevices: anc.childDevices || [],
    children: [],
    ancestors: prefix,
  }
}

function rewriteBreadcrumbSlot(slotCells, target, nav, index, templateMap) {
  const anc = pickAncestor(nav, target.entry.oldKind)
  if (!anc) return
  const ancNav = ancestorNavContext(anc, nav, index)
  const kind = ancNav.kind === 'root' ? 'home' : ancNav.kind
  const pageUuid = resolveTemplatePageIdForKind(templateMap, kind, ancNav.modelUuid)
  if (!pageUuid) return
  slotCells.forEach(cell => {
    setCellLink(cell, pageUuid, ancNav)
    const t = getText(cell)
    if (!t) return
    const m = t.match(/^←\s*(.*)$/)
    if (m) setText(cell, `← ${ancNav.label}`)
    else if (t.length > 1 && !/^[›|#·]+$/.test(t)) setText(cell, ancNav.label)
  })
}

// ---------- 各层级重映射 ----------

/** 收集槽位：旧 page_id → {entry, cellIdxs} */
function collectTargets(cells, index) {
  const map = new Map()
  cells.forEach((cell, i) => {
    const link = getClickLink(cell)
    if (!link) return
    const pid = link.Inside.pageUUID
    const entry = index.byOldPageId[pid]
    if (!entry) return
    let t = map.get(pid)
    if (!t) {
      t = { pageId: pid, entry, cellIdxs: [] }
      map.set(pid, t)
    }
    t.cellIdxs.push(i)
  })
  return [...map.values()].sort((a, b) => a.cellIdxs[0] - b.cellIdxs[0])
}

/** 模板页内未知 page_id：按 Y 坐标 + 文案推断槽位（floor 模板样本 link 与 uuid5 钥匙不一致） */
function augmentIndexFromCells(index, cells) {
  if (!index || !Array.isArray(cells)) return index
  const byOldPageId = { ...index.byOldPageId }
  const unknown = new Map()
  cells.forEach((cell, i) => {
    const link = getClickLink(cell)
    if (!link) return
    const pid = link.Inside.pageUUID
    if (!pid || byOldPageId[pid]) return
    let g = unknown.get(pid)
    if (!g) {
      g = { cellIdxs: [], minY: Infinity }
      unknown.set(pid, g)
    }
    g.cellIdxs.push(i)
    const y = cell.y || (cell.position && cell.position.y) || 0
    if (y < g.minY) g.minY = y
  })
  if (!unknown.size) return index

  const sorted = [...unknown.entries()].sort((a, b) => a[1].minY - b[1].minY)
  sorted.forEach(([pid, g]) => {
    const texts = g.cellIdxs.map(i => getText(cells[i]))
    const joined = texts.join('|')
    let oldKind = 'device'
    if (/^←|←\s/.test(joined) || texts.some(t => /^←/.test(t))) {
      oldKind = 'building'
    } else if (/详情\s*›/.test(joined) || texts.some(t => t === '详情 ›')) {
      oldKind = 'device'
    } else if (/全局总览/.test(joined)) {
      oldKind = 'home'
    } else if (/设备组/.test(joined)) {
      oldKind = 'floor'
    } else if (texts.some(t => /配电室|机房|模块/.test(t) && !/详情/.test(t))) {
      oldKind = 'zone'
    }
    byOldPageId[pid] = { oldKind, sid: null }
  })
  return { ...index, byOldPageId }
}

const cellsOf = (cells, idxs) => idxs.map(i => cells[i])

/** 当前 nav 的直接 children（仅一层，不展开子树） */
function getDirectChildNodes(nav) {
  const fromCtx = (nav && nav.childNodes) || []
  if (fromCtx.length) return sortByLabel(fromCtx, 'label')
  return sortByLabel(
    ((nav && nav.children) || []).map(treeChildToNodeItem).filter(Boolean),
    'label',
  )
}

function childNodeToDeviceItem(c) {
  return {
    name: c.name || c.label || '',
    uuid: c.uuid || '',
    code: c.uuid || '',
    modelUuid: c.modelUuid || c.muid || '',
    status: c.status || 'off',
    kind: 'device',
    layer: 'device',
    sid: c.sid,
  }
}

/**
 * 组织层：按 navContext 直接 children 填槽位（容器卡片 + 设备卡片混合）
 * 每个槽位 click link 携带子节点 navContext，逐层下钻
 */
function remapOrgContainerPage(cells, targets, nav, ctx) {
  const { index, templateMap } = ctx
  const directChildren = getDirectChildNodes(nav)
  const drop = new Set()

  // 子项全是设备且 >1 → 分页设备列表（UPS 122）
  if (directChildren.length > 1 && directChildren.every(c => isDeviceNode(c))) {
    const devices = directChildren.map(childNodeToDeviceItem)
    const listNav = applyGatewayListPagination({
      ...nav,
      allChildNodes: sortByLabel(directChildren, 'label'),
      childNodes: sortByLabel(directChildren, 'label'),
      allChildDevices: sortByLabel(devices, 'name'),
      childDevices: sortByLabel(devices, 'name'),
      deviceListMode: true,
    })
    return remapDeviceListPage(cells, targets, listNav, ctx)
  }

  const cardTargets = []
  const backTargets = []
  targets.forEach(target => {
    if (target.entry.oldKind === 'building' || target.entry.oldKind === 'floor') {
      cardTargets.push(target)
    } else if (target.entry.oldKind === 'zone' || target.entry.oldKind === 'room' || target.entry.oldKind === 'cabinet') {
      backTargets.push(target)
      rewriteBreadcrumbSlot(cellsOf(cells, target.cellIdxs), target, nav, index, templateMap)
    }
  })

  cardTargets.sort((a, b) => (a.cellIdxs[0] || 0) - (b.cellIdxs[0] || 0))

  cardTargets.forEach((target, r) => {
    const child = directChildren[r]
    const idxs = target.cellIdxs || []
    if (!child) {
      idxs.forEach(i => drop.add(i))
      return
    }
    const slotCells = cellsOf(cells, idxs)
    const treeNode = (child.sid != null && index.bySid[child.sid]) ? index.bySid[child.sid] : child

    if (isDeviceNode(treeNode) || child.kind === 'device') {
      const devName = child.label || child.name || ''
      const childNav = buildDeviceSignalContext(
        {
          label: devName,
          name: devName,
          uuid: child.uuid,
          sid: child.sid,
          modelUuid: child.modelUuid || child.muid,
          muid: child.modelUuid || child.muid,
        },
        [],
        childAncestors(nav),
      )
      const pageUuid = resolveTemplatePageIdForKind(templateMap, 'device', child.modelUuid || child.muid)
      slotCells.forEach(cell => {
        if (pageUuid) setCellLink(cell, pageUuid, childNav)
        setHoverDevice(cell, child.uuid, devName)
        const t = getText(cell)
        if (t && t !== '详情 ›' && !/^\d+台设备$/.test(t) && !/^[🟢🔴⏸]/.test(t) && !isFloorGroupArtifact(t)) {
          setText(cell, devName)
        }
      })
      rewriteSlotTexts(slotCells, {
        title: devName,
        devices: [{ uuid: child.uuid, status: child.status || 'off' }],
      })
    } else {
      const childNav = buildNavContextForNode(treeNode, index, childAncestors(nav))
      const pageUuid = resolveTemplatePageIdForKind(templateMap, childNav.kind, childNav.modelUuid)
      rewriteSlot(slotCells, pageUuid, childNav, childNav.label || childNav.name)
    }
  })

  return drop
}

/** root/home 页：链接 sid 即真实目标，逐链接直连 */
function remapRootPage(cells, targets, nav, ctx) {
  const { index, templateMap } = ctx
  const drop = new Set()
  targets.forEach(target => {
    const { entry } = target
    const node = index.bySid[entry.sid]
    if (!node) return
    // 旧 floor 槽位 → room 列表 + A3 转机分页（禁用 floor 主链路）
    if (entry.oldKind === 'floor') {
      const childNav = applyGatewayListPagination(buildNavContextForNode(node, index))
      const pageUuid = resolveDeviceListTemplateId(templateMap)
        || resolveTemplatePageIdForKind(templateMap, node.kind, childNav.modelUuid)
      if (!pageUuid) return
      rewriteSlot(cellsOf(cells, target.cellIdxs), pageUuid, childNav, `${childNav.label} · 转机列表`)
      return
    }
    const childNav = buildNavContextForNode(node, index)
    const pageUuid = resolveTemplatePageIdForKind(templateMap, node.kind, childNav.modelUuid)
    if (!pageUuid) return
    cellsOf(cells, target.cellIdxs).forEach(cell => {
      setCellLink(cell, pageUuid, childNav)
      const devs = childNav.kind === 'device'
        ? [{ uuid: childNav.uuid, name: childNav.name }]
        : childNav.childDevices
      if (devs && devs.length) setHoverDevice(cell, devs[0].uuid, devs[0].name)
    })
  })
  return drop
}

/** zone/room 页：直接 children 填槽（容器 + 设备混合） */
function remapZonePage(cells, targets, nav, ctx) {
  return remapOrgContainerPage(cells, targets, nav, ctx)
}

/** cabinet 页：子容器或设备列表 */
function remapCabinetPage(cells, targets, nav, ctx) {
  const childContainers = (nav.children || []).filter(c =>
    c && c.kind && !isDeviceNode(c),
  )
  if (childContainers.length) {
    const remappedTargets = targets.map(t => {
      if (t.entry.oldKind === 'floor' || t.entry.oldKind === 'building') {
        return { ...t, entry: { ...t.entry, oldKind: 'building' } }
      }
      return t
    })
    return remapOrgContainerPage(cells, remappedTargets, nav, ctx)
  }
  return remapDeviceListPage(cells, targets, nav, ctx)
}

/** zone/room/cabinet/floor：设备行槽位 → 当前页设备（按 y 坐标聚合整行 cells） */
function remapDeviceListPage(cells, targets, nav, ctx) {
  const { index, templateMap, dpMaps } = ctx
  const drop = new Set()
  const listMode = isDeviceListNav(nav)
  const pageSize = inferPageSizeFromTargets(targets)
  const pagedNav = nav.deviceListMode
    ? applyGatewayListPagination(nav, pageSize)
    : applyDeviceListPagination(nav, pageSize)
  const devices = sortByLabel(
    pagedNav.childNodes || pagedNav.childDevices || [],
    'name',
  )

  const rowTargets = []
  const cardTargets = []
  const backTargets = []
  targets.forEach(target => {
    if (target.entry.oldKind === 'device') rowTargets.push(target)
  })
  const hasDeviceRows = rowTargets.length > 0
  targets.forEach(target => {
    if (target.entry.oldKind === 'device') return
    if (target.entry.oldKind === 'floor' || target.entry.oldKind === 'building') {
      if (listMode && hasDeviceRows) {
        target.cellIdxs.forEach(i => drop.add(i))
      } else {
        cardTargets.push(target)
      }
    } else if (listMode && (target.entry.oldKind === 'zone' || target.entry.oldKind === 'room' || target.entry.oldKind === 'cabinet')) {
      backTargets.push(target)
    } else {
      rewriteBreadcrumbSlot(cellsOf(cells, target.cellIdxs), target, nav, index, templateMap)
    }
  })
  // 无 device 行槽位时，复用 floor/building 卡片槽位逐台展示（zone 模板兜底）
  if (!hasDeviceRows && cardTargets.length) {
    cardTargets.forEach(t => { rowTargets.push(t) })
  }
  if (!rowTargets.length) {
    if (listMode) stripFloorGroupArtifactCells(cells, drop)
    return drop
  }

  rowTargets.forEach((target, r) => {
    const device = devices[r]
    const idxs = target.cellIdxs || []
    if (!device) {
      idxs.forEach(i => drop.add(i))
      return
    }
    const isDev = isDeviceNode(device) || device.kind === 'device' || device.kind === 'gateway'
    const childNav = isDev
      ? buildDeviceSignalContext(
        { label: device.name, name: device.name, uuid: device.uuid, sid: device.sid,
          modelUuid: device.modelUuid, muid: device.modelUuid },
        [],
        childAncestors(pagedNav),
      )
      : {
        sid: null, uuid: device.uuid, kind: 'device',
        modelUuid: device.modelUuid || '', muid: device.modelUuid || '',
        name: device.name, label: device.name,
        childDevices: [], children: [],
        ancestors: childAncestors(pagedNav),
      }
    const pageUuid = resolveTemplatePageIdForKind(templateMap, 'device', device.modelUuid)
    const dpMap = (dpMaps && dpMaps[device.modelUuid]) || {}
    idxs.forEach(i => {
      const cell = cells[i]
      const detail = cell.data && cell.data.detail
      if (!detail) return
      setHoverDevice(cell, device.uuid, device.name)
      const link = getClickLink(cell)
      if (link && pageUuid) {
        setCellLink(cell, pageUuid, childNav)
        const t = getText(cell)
        if (t && t !== '详情 ›' && !/^\d+台设备$/.test(t) && !/^[🟢🔴⏸]/.test(t) && !isFloorGroupArtifact(t)) {
          setText(cell, device.name)
        }
      }
      walkConds(detail, cond => {
        if (typeof cond.deviceSN !== 'undefined') {
          cond.deviceSN = device.uuid
          cond.DeviceName = device.name
        }
        if (cond.dataName && dpMap[cond.dataName]) {
          cond.dataID = dpMap[cond.dataName]
        }
      })
    })
    if (listMode) {
      rewriteDeviceRowTexts(cells, idxs, device)
    } else if (!hasDeviceRows) {
      rewriteSlotTexts(cellsOf(cells, idxs), { title: device.name, devices: [device] })
    }
  })

  if (listMode) {
    backTargets.forEach(t => rewriteDeviceListBackLink(cells, t.cellIdxs, nav))
    stripFloorGroupArtifactCells(cells, drop)
  }

  rewritePaginationCells(cells, pagedNav)
  return drop
}

// ---------- 主入口 ----------

/**
 * 容器模板页槽位重映射（就地修改 cells，返回过滤后的新数组）。
 * @param {object[]} cells 已深拷贝的 cells
 * @param {object} nav navContext
 * @param {object} ctx {index, templateMap, dpMaps}
 * @returns {object[]} 重写后的 cells
 */
export function remapContainerCells(cells, nav, ctx) {
  if (!Array.isArray(cells) || !nav) {
    return cells
  }
  const baseIndex = (ctx && ctx.index && ctx.index.byOldPageId)
    ? ctx.index
    : { byOldPageId: Object.create(null), bySid: Object.create(null), parentBySid: Object.create(null) }
  const index = augmentIndexFromCells(baseIndex, cells)
  const targets = collectTargets(cells, index)
  if (!targets.length) return cells

  let drop
  const kind = nav.kind
  const layer = nav.layer || kind
  const remapCtx = { ...(ctx || {}), index }
  if (kind === 'root' || kind === 'home') {
    drop = remapRootPage(cells, targets, nav, remapCtx)
  } else if (kind === 'zone' || kind === 'room' || layer === 'zone' || layer === 'room') {
    drop = remapZonePage(cells, targets, nav, remapCtx)
  } else if (kind === 'cabinet' || kind === 'gateway' || layer === 'gateway' || nav.signalMode) {
    drop = remapCabinetPage(cells, targets, nav, remapCtx)
  } else if (kind === 'floor') {
    drop = remapDeviceListPage(cells, targets, nav, remapCtx)
  } else {
    return cells
  }
  const filtered = drop && drop.size
    ? cells.filter((_, i) => !drop.has(i))
    : cells
  return sanitizeGraphComponents(filtered, { tag: 'slotRemap' }).cells
}

export { remapOrgContainerPage }
export const CONTAINER_KINDS = ['root', 'home', 'zone', 'room', 'cabinet', 'gateway', 'floor']
