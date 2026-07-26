/**
 * 层级模板页：按 NavContext 解析相对绑点（deviceSN/dataID）与动态表格行。
 * 约定：
 * - 设备模板页（nav.kind==='device'）：deviceSN 为空的相对绑点灌当前设备 uuid
 * - 容器模板页（root/zone/room/cabinet/floor）：由 slotRemap 槽位重映射逐槽注入，
 *   纯空占位 condition 不再整页灌同一 uuid（否则悬浮提示全页同列表）
 * - dataName 有值且 dataID 空/占位 → 按物模型测点名解析 dataID
 * - ViewRealTable diy rowSource=navChildren → 组织层子设备列表
 * - ViewRealTable diy rowSource=navDatapoints → 信号层测点列表（单设备）
 * - 文本占位符：{{nav.name}} / {{nav.deviceCount}} / {{nav.onlineCount}} /
 *   {{nav.offlineCount}} / {{nav.abnormalCount}} / {{nav.childCount}}
 */

import { GetDataModelData } from '@/services/system'
import { getModelDataPoints } from '@/services/device'
import { remapContainerCells, CONTAINER_KINDS } from './slotRemap'
import { sanitizeGraphComponents } from './graphCellSanitizer'
import {
  applyDeviceListPagination,
  isDeviceListNav,
  isFloorGroupArtifact,
  formatPageInfo,
  rewritePaginationCells,
  detectPageNavRole,
  deviceListPageSizeForNav,
} from './deviceListPager'
import {
  applyDatapointPagination,
  fetchDeviceDatapointPage,
  buildDeviceSignalContext,
  buildSignalTablePageConfig,
  formatDatapointPageInfo,
  isSignalDatapointPaged,
  DEFAULT_DATAPOINT_PAGE_SIZE,
  resolveDatapointFetchParams,
} from './navContext'

const dpCache = Object.create(null)

/** 空 layer 模板页（bootstrap 脚本）的默认画布 — 与 build_ncc_dashboard 一致 */
export const DEFAULT_TEMPLATE_LAYER = {
  width: 1920,
  height: 1080,
  autoSize: 1,
  Padding: 0,
  gridSize: 10,
  backColor: '#0a0e17',
  backgroundImage: '',
}

const TPL_MAIN_X = 16
const TPL_MAIN_W = 1920 - 32
const TPL_TITLE_Y = 72
const TPL_SUBTITLE_Y = 112
const TPL_TABLE_Y = 152
const TPL_TABLE_H = 1080 - TPL_TABLE_Y - 32
const SIGNAL_PAGER_Y = 72
const SIGNAL_PAGER_LAYOUT = Object.freeze({
  prev: { x: 1480, w: 88 },
  info: { x: 1572, w: 176 },
  next: { x: 1752, w: 88 },
})
const DEVICE_LIST_PAGER_Y = TPL_SUBTITLE_Y
const DEVICE_LIST_PAGER_LAYOUT = SIGNAL_PAGER_LAYOUT

export function ensureTemplatePageLayer(pageData) {
  if (!pageData) return pageData
  const raw = pageData.layer
  const base = (raw && typeof raw === 'object' && !Array.isArray(raw)) ? raw : {}
  const layer = { ...DEFAULT_TEMPLATE_LAYER, ...base }
  if (base.background && !layer.backColor) {
    layer.backColor = base.background
  }
  if (!layer.width) layer.width = DEFAULT_TEMPLATE_LAYER.width
  if (!layer.height) layer.height = DEFAULT_TEMPLATE_LAYER.height
  if (layer.autoSize == null) layer.autoSize = 1
  pageData.layer = layer
  return pageData
}

function templateBaseAnimate() {
  return {
    selected: [],
    animateElement: [],
    animateList: [],
    isExpression: false,
    condition: {
      deviceSN: '', selectVideoType: 0, isBandDevice: false, bandType: 1,
      dataID: '', dataName: '', operator: '', OperatorValue: '', OperatorMaxValue: '',
    },
  }
}

function makeTemplateTextCell(id, x, y, w, h, text, opts = {}) {
  const z = opts.z != null ? opts.z : 10
  const foreColor = opts.foreColor || '#cfe0f5'
  const fontSize = opts.fontSize || 14
  return {
    shape: 'view-svg-text',
    id,
    x, y,
    width: w,
    height: h,
    zIndex: z,
    visible: true,
    position: { x, y },
    size: { width: w, height: h },
    data: {
      editMode: false,
      IsToolBox: false,
      detail: {
        type: 'view-svg-text',
        identifier: id,
        name: opts.name || id,
        style: {
          position: { x, y, w, h },
          visible: 1,
          text: String(text),
          foreColor,
          fontSize,
          fontFamily: 'Microsoft YaHei',
          fontWeight: 400,
          backColor: 'transparent',
          zIndex: z,
          transform: 0,
          diy: opts.diy || [],
        },
        animate: templateBaseAnimate(),
        active: [],
        action: [],
        dataBind: [],
      },
    },
  }
}

/** 测点页翻页控件按角色布局，不能继承旧模板的设备组坐标。 */
function layoutSignalPagerChrome(cells) {
  (cells || []).forEach(cell => {
    const detail = cell && cell.data && cell.data.detail
    const style = detail && detail.style
    const diy = (style && style.diy) || []
    const role = String((diy.find(d => d && d.key === 'labelRole') || {}).value || '')
    const slot = /detailPagePrev/.test(role)
      ? SIGNAL_PAGER_LAYOUT.prev
      : (/detailPageNext/.test(role)
        ? SIGNAL_PAGER_LAYOUT.next
        : (role === 'detailPageInfo' ? SIGNAL_PAGER_LAYOUT.info : null))
    if (!slot || !style) return
    const h = Number(style.position && style.position.h) || 28
    cell.x = slot.x
    cell.y = SIGNAL_PAGER_Y
    cell.width = slot.w
    cell.height = h
    cell.position = { ...(cell.position || {}), x: slot.x, y: SIGNAL_PAGER_Y }
    cell.size = { ...(cell.size || {}), width: slot.w, height: h }
    style.position = { ...(style.position || {}), x: slot.x, y: SIGNAL_PAGER_Y, w: slot.w, h }
  })
}

/** 设备卡片页翻页控件统一靠右紧凑排列，与底部分页使用同一页码口径。 */
function layoutDeviceListPagerChrome(cells) {
  (cells || []).forEach(cell => {
    const detail = cell && cell.data && cell.data.detail
    const style = detail && detail.style
    const diy = (style && style.diy) || []
    const role = String((diy.find(d => d && d.key === 'labelRole') || {}).value || '')
    const slot = /^pagePrev/.test(role)
      ? DEVICE_LIST_PAGER_LAYOUT.prev
      : (/^pageNext/.test(role)
        ? DEVICE_LIST_PAGER_LAYOUT.next
        : (role === 'pageInfo' ? DEVICE_LIST_PAGER_LAYOUT.info : null))
    if (!slot || !style) return
    const h = Number(style.position && style.position.h) || 28
    cell.x = slot.x
    cell.y = DEVICE_LIST_PAGER_Y
    cell.width = slot.w
    cell.height = h
    cell.position = { ...(cell.position || {}), x: slot.x, y: DEVICE_LIST_PAGER_Y }
    cell.size = { ...(cell.size || {}), width: slot.w, height: h }
    style.position = { ...(style.position || {}), x: slot.x, y: DEVICE_LIST_PAGER_Y, w: slot.w, h }
  })
}

function makeTemplateTableCell(id, x, y, w, h, diy, opts = {}) {
  const z = opts.z != null ? opts.z : 8
  return {
    shape: 'ism-view-real-table',
    id,
    x, y,
    width: w,
    height: h,
    zIndex: z,
    visible: true,
    position: { x, y },
    size: { width: w, height: h },
    data: {
      editMode: false,
      IsToolBox: false,
      detail: {
        type: 'ism-view-real-table',
        identifier: id,
        name: opts.name || '实时数据表',
        style: {
          position: { x, y, w, h },
          visible: 1,
          foreColor: '#cfe0f5',
          fontSize: 14,
          fontFamily: 'Microsoft YaHei',
          fontWeight: 400,
          backColor: 'transparent',
          zIndex: z,
          transform: 0,
          diy,
        },
        animate: templateBaseAnimate(),
        active: [],
        action: [],
        dataBind: [],
      },
    },
  }
}

/** 将任意历史实时表 cell 收敛为运行态唯一卡片容器，避免模板差异泄漏到组织/设备页。 */
function normalizeRuntimeTableCell(cell, id, rowSource, pageSize, name) {
  const table = cell || makeTemplateTableCell(
    id, TPL_MAIN_X, TPL_TABLE_Y, TPL_MAIN_W, TPL_TABLE_H, [], { name },
  )
  const detail = table.data && table.data.detail
  if (!detail || !detail.style) return table
  table.x = TPL_MAIN_X
  table.y = TPL_TABLE_Y
  table.width = TPL_MAIN_W
  table.height = TPL_TABLE_H
  table.position = { ...(table.position || {}), x: TPL_MAIN_X, y: TPL_TABLE_Y }
  table.size = { ...(table.size || {}), width: TPL_MAIN_W, height: TPL_TABLE_H }
  detail.name = name
  detail.style.position = {
    ...(detail.style.position || {}),
    x: TPL_MAIN_X,
    y: TPL_TABLE_Y,
    w: TPL_MAIN_W,
    h: TPL_TABLE_H,
  }
  if (!Array.isArray(detail.style.diy)) detail.style.diy = []
  const setDiy = (key, value) => {
    const item = detail.style.diy.find(d => d && d.key === key)
    if (item) item.value = value
    else detail.style.diy.push({ name: key, type: 9, value, key })
  }
  setDiy('rowSource', rowSource)
  setDiy('ShowCount', pageSize)
  setDiy('themeName', 'dark')
  return table
}

export function clearNavDpCache() {
  Object.keys(dpCache).forEach(k => { delete dpCache[k] })
}

/**
 * @param {string} modelUuid 物模型 muid
 * @returns {Promise<Record<string, string>>} dataName -> data uuid
 */
export async function loadDpMapByModel(modelUuid) {
  if (!modelUuid) return {}
  if (dpCache[modelUuid]) return dpCache[modelUuid]
  try {
    const res = await getModelDataPoints({ muid: modelUuid })
    const list = (res && res.data && res.data.list) || []
    const map = Object.create(null)
    if (Array.isArray(list)) {
      for (const item of list) {
        const name = item.name || item.Name || item.dataName
        const uuid = item.uuid || item.Uuid || item.dataID
        if (name && uuid) map[String(name)] = String(uuid)
      }
    }
    // 兜底：旧接口
    if (!Object.keys(map).length) {
      try {
        const res2 = await GetDataModelData({ type: 2 })
        const list2 = (res2 && res2.data && res2.data.list) || []
        for (const item of list2) {
          if ((item.muid || item.Muid) !== modelUuid) continue
          const name = item.name || item.Name
          const uuid = item.uuid || item.Uuid
          if (name && uuid) map[String(name)] = String(uuid)
        }
      } catch (e2) { /* ignore */ }
    }
    dpCache[modelUuid] = map
    return map
  } catch (e) {
    console.warn('[navContextBinding] loadDpMap failed', modelUuid, e && e.message)
    dpCache[modelUuid] = {}
    return {}
  }
}

function walkConditions(obj, visit) {
  if (!obj || typeof obj !== 'object') return
  if (Array.isArray(obj)) {
    obj.forEach(v => walkConditions(v, visit))
    return
  }
  if (obj.condition && typeof obj.condition === 'object') {
    visit(obj.condition)
  }
  Object.keys(obj).forEach(k => {
    if (k === 'condition') return
    const v = obj[k]
    if (v && typeof v === 'object') walkConditions(v, visit)
  })
}

/**
 * 是否相对设备绑点。
 * 设备模板页：空 deviceSN 也视为相对（整页只关于一台设备）。
 * 容器模板页：仅显式标记（__NAV__/{{device}}）才灌，纯空占位交给槽位重映射。
 */
function isRelativeDeviceBinding(cond, navKind) {
  if (!cond || typeof cond !== 'object') return false
  if (cond.isBandDevice === true) return false
  const sn = cond.deviceSN
  if (sn === '__NAV__' || sn === '{{device}}') return true
  if (navKind === 'device' || navKind === 'datapoint') return true
  return false
}

function collectDeviceChildren(node) {
  const out = []
  if (!node) return out
  const walk = n => {
    if (!n) return
    if (n.kind === 'gateway' || n.kind === 'device') {
      out.push({
        name: n.label || n.name || '',
        uuid: n.uuid || n.deviceUuid || '',
        code: n.code || n.uuid || '',
        modelUuid: n.modelUuid || n.muid || '',
        status: n.status || 'off',
        kind: 'gateway',
        layer: 'gateway',
      })
      return
    }
    (n.children || []).forEach(walk)
  }
  ;(node.children || []).forEach(walk)
  return out
}

/** 页级统计（占位符解析用） */
function navStats(navContext, childDevices) {
  const all = navContext.allChildDevices
    || (isDeviceListNav(navContext) ? (navContext.childDevices || navContext.childNodes || childDevices) : null)
  const devices = all || childDevices || navContext.childNodes || []
  const online = devices.filter(d => d && d.status === 'on').length
  const containers = (navContext.children || []).filter(c =>
    c && c.kind && c.kind !== 'device' && c.kind !== 'gateway',
  )
  const deviceCount = navContext.totalCount != null
    ? navContext.totalCount
    : (navContext.totalDevices != null
      ? navContext.totalDevices
      : devices.length)
  return {
    deviceCount,
    deviceTotal: deviceCount,
    onlineCount: online,
    offlineCount: deviceCount - online,
    abnormalCount: deviceCount - online,
    childCount: containers.length || (navContext.childNodes || []).length,
    pageInfo: navContext.totalDevices != null || navContext.totalCount != null
      ? formatPageInfo(navContext) : '',
    detailPageInfo: isSignalDatapointPaged(navContext) ? formatDatapointPageInfo(navContext) : '',
  }
}

function applyTextPlaceholders(style, navContext, stats) {
  if (!style || typeof style.text !== 'string' || style.text.indexOf('{{') < 0) return
  const name = navContext.name || navContext.label || ''
  style.text = style.text
    .replace(/\{\{\s*nav\.name\s*\}\}/g, name)
    .replace(/\{\{\s*nav\.label\s*\}\}/g, navContext.label || name)
    .replace(/\{\{\s*nav\.deviceCount\s*\}\}/g, String(stats.deviceCount))
    .replace(/\{\{\s*nav\.deviceTotal\s*\}\}/g, String(stats.deviceTotal))
    .replace(/\{\{\s*nav\.onlineCount\s*\}\}/g, String(stats.onlineCount))
    .replace(/\{\{\s*nav\.offlineCount\s*\}\}/g, String(stats.offlineCount))
    .replace(/\{\{\s*nav\.abnormalCount\s*\}\}/g, String(stats.abnormalCount))
    .replace(/\{\{\s*nav\.childCount\s*\}\}/g, String(stats.childCount))
    .replace(/\{\{\s*nav\.pageInfo\s*\}\}/g, stats.pageInfo || '')
    .replace(/\{\{\s*nav\.detailPageInfo\s*\}\}/g, stats.detailPageInfo || '')
}

/** 静态模板文案改写（无占位符的硬编码样本文本） */
function rewriteStaticNavTexts(style, navContext, stats) {
  if (!style || typeof style.text !== 'string') return
  let t = style.text
  const name = navContext.name || navContext.label || ''
  if (isDeviceListNav(navContext)) {
    if (isFloorGroupArtifact(t)) {
      style.text = ''
      return
    }
  }
  if (t === 'RootZone' && navContext.kind !== 'root' && navContext.kind !== 'home') {
    style.text = name
    return
  }
  if (/^RootZone\s*[›|>]\s*/.test(t)) {
    style.text = t.replace(/^RootZone/, name)
    return
  }
  if (/^\d+台设备$/.test(t)) {
    style.text = `${stats.deviceCount}台设备`
    return
  }
  const subM = t.match(/^(\d+)台设备(\s*[·|].*)$/)
  if (subM) {
    style.text = `${stats.deviceCount}台设备${subM[2]}`
    return
  }
  const sumM = t.match(/^共\s*(\d+)\s*台设备/)
  if (sumM) {
    style.text = t.replace(sumM[1], String(stats.deviceCount))
  }
}

/**
 * 顶部 x≈833 的“全局总览”是固定系统入口，不属于随 nav.name 改写的面包屑。
 * 统一写入独立角色和当前项目首页目标，绕开旧 page_id/静态模板重映射。
 */
function normalizeGlobalOverviewCell(cell, homePageUuid) {
  if (!cell || !homePageUuid) return
  const detail = cell.data && cell.data.detail
  const style = detail && detail.style
  if (!detail || !style || typeof style.text !== 'string') return
  const diy = Array.isArray(style.diy) ? style.diy : (style.diy = [])
  const roleItem = diy.find(d => d && d.key === 'labelRole')
  const role = String((roleItem || {}).value || '')
  const position = style.position || {}
  // X6 top-level 坐标是运行态真实位置；部分旧模板的 style.position 被整组复用，不能用于区分相邻面包屑。
  const x = Number(cell.x != null ? cell.x : position.x)
  const y = Number(cell.y != null ? cell.y : position.y)
  const width = Number(cell.width != null ? cell.width : position.w)
  const isFixedHeaderSlot = y >= 0 && y <= 56 && x >= 720 && x < 840 && width >= 60
  if (role !== 'globalOverview' && !isFixedHeaderSlot) return

  style.text = '📊 全局总览'
  if (roleItem) roleItem.value = 'globalOverview'
  else diy.push({ name: 'labelRole', type: 9, value: 'globalOverview', key: 'labelRole' })
  const homeItem = diy.find(d => d && d.key === 'homePageUuid')
  if (homeItem) homeItem.value = homePageUuid
  else diy.push({ name: 'homePageUuid', type: 9, value: homePageUuid, key: 'homePageUuid' })

  if (!Array.isArray(detail.action)) detail.action = []
  let action = detail.action.find(a => a && a.type === 'click' && a.action === 'link')
  if (!action) {
    action = { type: 'click', action: 'link', link: {} }
    detail.action.push(action)
  }
  if (!action.link) action.link = {}
  action.link.linkType = 'Inside'
  action.link.isPopUp = false
  action.link.autoClose = false
  action.link.navContext = null
  action.link.Inside = {
    ...(action.link.Inside || {}),
    displayUUID: homePageUuid,
    pageUUID: homePageUuid,
  }
}

function compactBreadcrumbSegment(value, maxLength = 24) {
  const text = String(value || '').trim()
  if (text.length <= maxLength) return text
  const side = Math.max(6, Math.floor((maxLength - 1) / 2))
  return `${text.slice(0, side)}…${text.slice(-side)}`
}

/** 点位页顶部路径：固定首页入口之后仅展示真实组织层级、父设备与当前虚拟柜/设备。 */
function buildDeviceBreadcrumb(nav) {
  const organizations = []
  const add = (value, kind = '') => {
    const text = String(value || '').trim()
    const normalizedKind = String(kind || '').trim().toLowerCase()
    if (normalizedKind === 'root'
      || normalizedKind === 'home'
      || normalizedKind === 'devicelist'
      || normalizedKind === 'device-list'
      || /^(?:RootZone|root|home|📊\s*全局总览|全局总览|设备列表)$/i.test(text)) return
    if (text && text !== nav.name && organizations[organizations.length - 1] !== text) {
      organizations.push(text)
    }
  }
  ;(nav.ancestors || []).forEach(item => add(item && (item.label || item.name), item && item.kind))
  add(
    nav.deviceListReturnContext && (nav.deviceListReturnContext.label || nav.deviceListReturnContext.name),
    nav.deviceListReturnContext && nav.deviceListReturnContext.kind,
  )
  // 虚拟列头柜：补上父设备名（若 ancestors 未带）
  if (nav.virtualCabinet && nav.parentDeviceLabel) {
    add(nav.parentDeviceLabel, 'device')
  }

  let visibleOrganizations = organizations
  if (visibleOrganizations.length > 3) {
    visibleOrganizations = [
      visibleOrganizations[0],
      '…',
      visibleOrganizations[visibleOrganizations.length - 1],
    ]
  }
  visibleOrganizations = visibleOrganizations.map(name => compactBreadcrumbSegment(name, 18))
  const deviceName = compactBreadcrumbSegment(nav.name || nav.label || '未命名设备', 36)
  return `› ${[...visibleOrganizations, deviceName].join(' › ')}`
}

function normalizeDeviceBreadcrumbCell(cells, nav) {
  const candidate = (cells || []).find(cell => {
    const detail = cell && cell.data && cell.data.detail
    const style = detail && detail.style
    const position = (style && style.position) || {}
    const text = String((style && style.text) || '').trim()
    return cell && cell.shape === 'view-svg-text'
      && !text
      && Number(cell.x != null ? cell.x : position.x) === 0
      && Number(cell.y != null ? cell.y : position.y) === 0
      && Number(cell.width != null ? cell.width : position.w) >= 1200
      && Number(cell.height != null ? cell.height : position.h) <= 64
  })
  if (!candidate) return
  const detail = candidate.data.detail
  const style = detail.style
  style.text = buildDeviceBreadcrumb(nav)
  style.textAlign = 'left'
  style.foreColor = '#8fb8cc'
  style.fontSize = 12
  style.fontWeight = 500
  if (!Array.isArray(style.diy)) style.diy = []
  const role = style.diy.find(d => d && d.key === 'labelRole')
  if (role) role.value = 'deviceBreadcrumb'
  else style.diy.push({ name: 'labelRole', type: 9, value: 'deviceBreadcrumb', key: 'labelRole' })
  detail.name = '当前设备面包屑'
}

/**
 * 顶栏系统标题左对齐，并略向左收；与面包屑之间留出中间空隙。
 * 匹配 build_ncc_dashboard 的 header-logo / header-title / header-subtitle 槽位。
 */
function normalizeHeaderTitleLayout(cells) {
  const moveCellX = (cell, style, position, y, w, targetX) => {
    const x = Number(cell.x != null ? cell.x : position.x)
    if (Math.abs(x - targetX) <= 4) return
    cell.x = targetX
    if (cell.position) cell.position.x = targetX
    style.position = {
      ...position,
      x: targetX,
      y,
      w,
      h: Number(cell.height != null ? cell.height : position.h) || position.h,
    }
  }
  ;(cells || []).forEach(cell => {
    if (!cell || cell.shape !== 'view-svg-text') return
    const detail = cell.data && cell.data.detail
    const style = detail && detail.style
    if (!style) return
    const position = style.position || {}
    const x = Number(cell.x != null ? cell.x : position.x)
    const y = Number(cell.y != null ? cell.y : position.y)
    const w = Number(cell.width != null ? cell.width : position.w)
    const text = String(style.text || '')
    const name = String((detail && detail.name) || '')
    const id = String(cell.id || '')
    const isLogo = y >= 0 && y <= 50
      && x >= 100 && x <= 280
      && w <= 48
      && (/⚡|header-logo/i.test(text) || /header-logo/i.test(name) || /header-logo/i.test(id))
    if (isLogo) {
      moveCellX(cell, style, position, y, w, 16)
      return
    }
    const isTitle = y >= 0 && y <= 28
      && x >= 40 && x <= 400
      && w >= 280
      && (/监控系统/.test(text) || /header-title/i.test(name) || /header-title/i.test(id))
    const isSubtitle = y > 28 && y <= 50
      && x >= 40 && x <= 400
      && w >= 280
      && (/POWER MONITORING|DATA CENTER/i.test(text)
        || /header-subtitle/i.test(name)
        || /header-subtitle/i.test(id))
    if (!isTitle && !isSubtitle) return
    style.textAlign = 'left'
    // 标题靠左（logo 旁），与面包屑（约 x≥1020）之间留中间空隙
    moveCellX(cell, style, position, y, w, isTitle ? 56 : 58)
  })
}

/**
 * 设备详情模板由任一设备样本生成时，标题和“设备名称”字段会带有该样本名称。
 * 信号层复用模板时，仅替换与标题样本完全相同的文本，避免误改测点名称或业务文案。
 */
function rewriteStaticDeviceName(style, templateDeviceName, currentDeviceName) {
  if (!style || typeof style.text !== 'string' || !templateDeviceName || !currentDeviceName) return
  if (style.text === templateDeviceName) {
    style.text = currentDeviceName
  } else if (style.text === `🔧 ${templateDeviceName}`) {
    style.text = `🔧 ${currentDeviceName}`
  }
}

function findTemplateDeviceName(cells) {
  for (const cell of cells || []) {
    const style = cell && cell.data && cell.data.detail && cell.data.detail.style
    const text = style && style.text
    const match = typeof text === 'string' && text.match(/^🔧\s+(.+?)\s*$/)
    if (match && match[1]) return match[1]
  }
  return ''
}

function isViewRealTableCell(cell) {
  if (!cell) return false
  if (cell.shape === 'ism-view-real-table') return true
  const type = cell.data && cell.data.detail && cell.data.detail.type
  return type === 'ism-view-real-table'
}

/** 清理模板在运行态页眉左侧写死的样本设备名，保留统一动态控件。 */
function isLegacyRuntimeHeaderCell(cell, maxX) {
  const style = ((((cell || {}).data || {}).detail || {}).style || {})
  const text = String(style.text || '').trim()
  const position = style.position || {}
  const diy = style.diy || []
  const labelRole = String((diy.find(d => d && d.key === 'labelRole') || {}).value || '')
  return !!text
    && Number(position.y) > 56
    && Number(position.y) < 145
    && Number(position.x) < maxX
    && !/全局总览/.test(text)
    && !/^(deviceListBack|deviceInfo)/.test(labelRole)
}

/**
 * 纯设备列表（childrenList）：无 room 模板槽位时注入标题 + ViewRealTable(navChildren) + 翻页控件
 */
export function ensureDeviceListPageLayout(cells, nav) {
  if (!isDeviceListNav(nav)) return cells || []
  // x<180 的旧样本设备名与 x=210 开始的统一组织标题属于不同区域。
  let list = (cells || []).filter(cell => {
    const style = ((((cell || {}).data || {}).detail || {}).style || {})
    const diy = style.diy || []
    const role = String((diy.find(d => d && d.key === 'labelRole') || {}).value || '')
    // 模板缓存可能刚被点位页注入动态 chrome；返回列表时必须清除点位专属控件。
    if (/^(deviceListBack|deviceInfo|deviceBreadcrumb|detailPage)/.test(role)) return false
    return !isLegacyRuntimeHeaderCell(cell, 180)
  })
  const pageSize = deviceListPageSizeForNav(nav)
  const uid = nav.sid != null ? nav.sid : (nav.uuid || 'list')

  const table = normalizeRuntimeTableCell(
    list.find(isViewRealTableCell),
    `devlist-vrt-${uid}`,
    'navChildren',
    pageSize,
    '设备卡片列表',
  )
  list = list.filter(cell => !isViewRealTableCell(cell))
  list.push(table)

  const hasTitle = list.some(c => {
    const detail = ((c || {}).data || {}).detail || {}
    const text = (detail.style && detail.style.text) || ''
    return /\{\{\s*nav\.name\s*\}\}/.test(String(text))
  })
  if (!hasTitle) {
    list.unshift(makeTemplateTextCell(
      `devlist-title-${uid}`, TPL_MAIN_X, TPL_TITLE_Y, TPL_MAIN_W, 36,
      '{{nav.name}}', { foreColor: '#00e5ff', fontSize: 20, name: '设备列表标题' },
    ))
  }

  const hasPageInfo = list.some(c => {
    const detail = ((c || {}).data || {}).detail || {}
    const text = (detail.style && detail.style.text) || ''
    return detectPageNavRole(text) === 'pageInfo' || /\{\{\s*nav\.pageInfo\s*\}\}/.test(text)
  })
  if (!hasPageInfo) {
    const mkNavText = (id, x, text, role) => makeTemplateTextCell(
      id, x, TPL_SUBTITLE_Y, role === 'pageInfo' ? 280 : 88, 28, text,
      {
        foreColor: '#b9cce6',
        fontSize: 14,
        diy: role ? [{ name: 'labelRole', type: 9, value: role, key: 'labelRole' }] : [],
      },
    )
    list.push(
      mkNavText(`devlist-prev-${uid}`, DEVICE_LIST_PAGER_LAYOUT.prev.x, '‹ 上一页', 'pagePrev'),
      mkNavText(`devlist-info-${uid}`, DEVICE_LIST_PAGER_LAYOUT.info.x, '{{nav.pageInfo}}', 'pageInfo'),
      mkNavText(`devlist-next-${uid}`, DEVICE_LIST_PAGER_LAYOUT.next.x, '下一页 ›', 'pageNext'),
    )
  }

  return list
}

/**
 * 设备信号层模板页：标题 + ViewRealTable(navDatapoints) + 页码占位
 * 与 navChildren（多设备各一行）不同：navDatapoints = 单设备、测点各占一行
 */
export function ensureDeviceDatapointPageLayout(cells, nav) {
  if (!nav || !nav.signalMode) return cells || []
  // 旧设备组详情页遗留了“← Y11设备组”和静态“🔧 设备名”页眉。
  // 当前链路由左侧组织树承载上下文，设备页只展示按需加载的测点表，
  // 因此不能保留这些会误导层级关系的历史 chrome。
  let list = (cells || []).filter(cell => {
    const style = ((((cell || {}).data || {}).detail || {}).style || {})
    const text = String(style.text || '').trim()
    const position = style.position || {}
    const diy = style.diy || []
    const labelRole = String((diy.find(d => d && d.key === 'labelRole') || {}).value || '')
    // 客户反馈：编号/模型/状态 UUID 串对运维无意义，测点页不再展示。
    if (labelRole === 'deviceInfoMeta') return false
    // 旧大屏在顶栏写死了“区域 → 机房 → 设备组 → 设备”的多段面包屑及跳转。
    // 新架构仅保留全局总览入口，其余上下文由实时组织树提供。
    const isLegacyHeaderBreadcrumb = Number(position.y) <= 56
      && Number(position.x) >= 740
      && Number(position.x) < 1640
      && !/全局总览/.test(text)
    // 画布 y=56~145、x<720 是统一设备信息栏。历史模板在这里写死样本设备名，
    // 无论具体字符串为何都应清理；全局入口、有效返回按钮和本次动态角色必须保留。
    return !(isLegacyHeaderBreadcrumb
      || isLegacyRuntimeHeaderCell(cell, 720)
      || /设备组/.test(text)
      || /^🔧\s*(?!\{\{).+/.test(text)
      || /^●\s*(离线|运行|运行中)$/.test(text))
  })
  normalizeHeaderTitleLayout(list)
  normalizeDeviceBreadcrumbCell(list, nav)
  // 历史模板以裸 › 表示面包屑层级，视觉像小方块且曾误触发分页。
  // 统一替换为明确的方向箭头并附上独立语义角色。
  list.forEach(cell => {
    const style = cell && cell.data && cell.data.detail && cell.data.detail.style
    if (!style || String(style.text || '').trim() !== '›') return
    style.text = '➜'
    if (!Array.isArray(style.diy)) style.diy = []
    const role = style.diy.find(d => d && d.key === 'labelRole')
    if (role) role.value = 'breadcrumbArrow'
    else style.diy.push({ name: 'labelRole', type: 9, value: 'breadcrumbArrow', key: 'labelRole' })
  })
  const pageSize = nav.datapointPageSize || DEFAULT_DATAPOINT_PAGE_SIZE
  const uid = nav.deviceUuid || nav.uuid || nav.sid || 'device'
  const roleOf = cell => {
    const diy = (((((cell || {}).data || {}).detail || {}).style || {}).diy) || []
    return String((diy.find(d => d && d.key === 'labelRole') || {}).value || '')
  }
  const cellByRole = role => list.find(cell => roleOf(cell) === role)
  if (!list.some(cell => roleOf(cell) === 'deviceListBack')) {
    list.push(makeTemplateTextCell(
      `devpoints-back-${uid}`, 16, SIGNAL_PAGER_Y, 120, 30, '‹ 返回上一级',
      {
        foreColor: '#7ee8ff',
        fontSize: 13,
        name: '返回设备列表',
        diy: [{ name: 'labelRole', type: 9, value: 'deviceListBack', key: 'labelRole' }],
      },
    ))
  }
  const deviceNameText = `设备：${nav.name || nav.label || '未命名设备'}`
  const deviceNameCell = cellByRole('deviceInfoName')
  if (deviceNameCell) {
    deviceNameCell.data.detail.style.text = deviceNameText
  } else {
    list.push(makeTemplateTextCell(
      `devpoints-device-${uid}`, 152, SIGNAL_PAGER_Y, 500, 30,
      deviceNameText,
      {
        foreColor: '#dffaff',
        fontSize: 15,
        name: '当前设备名称',
        diy: [{ name: 'labelRole', type: 9, value: 'deviceInfoName', key: 'labelRole' }],
      },
    ))
  }

  const table = normalizeRuntimeTableCell(
    list.find(isViewRealTableCell),
    `devpoints-vrt-${uid}`,
    'navDatapoints',
    pageSize,
    '测点实时数据卡片',
  )
  list = list.filter(cell => !isViewRealTableCell(cell))
  list.push(table)

  // 顶部翻页：上一页 / 页码 / 下一页（与底部分页条同步，均走 NavPageChange）
  const hasTopPagerBtns = list.some(c => {
    const diy = (((((c || {}).data || {}).detail || {}).style || {}).diy) || []
    return diy.some(d => d && d.key === 'labelRole'
      && /detailPagePrev|detailPageNext/.test(String(d.value || '')))
  })
  if (!hasTopPagerBtns) {
    const mkNavText = (id, x, text, role, w) => makeTemplateTextCell(
      id, x, TPL_SUBTITLE_Y, w || (role === 'detailPageInfo' ? 280 : 88), 28, text,
      {
        foreColor: '#b9cce6',
        fontSize: 14,
        diy: role ? [{ name: 'labelRole', type: 9, value: role, key: 'labelRole' }] : [],
      },
    )
    // 若已有仅页码文案，仍补上左右按钮
    const hasInfoOnly = list.some(c => {
      const diy = (((((c || {}).data || {}).detail || {}).style || {}).diy) || []
      return diy.some(d => d && d.key === 'labelRole' && String(d.value) === 'detailPageInfo')
    })
    if (!hasInfoOnly) {
      list.push(mkNavText(`devpoints-info-${uid}`, TPL_MAIN_X + 96, '{{nav.detailPageInfo}}', 'detailPageInfo', 300))
    }
    list.push(
      mkNavText(`devpoints-prev-${uid}`, TPL_MAIN_X, '‹ 上一页', 'detailPagePrev', 88),
      mkNavText(`devpoints-next-${uid}`, TPL_MAIN_X + 410, '下一页 ›', 'detailPageNext', 88),
    )
  }

  return list
}

/** @deprecated 使用 ensureDeviceDatapointPageLayout */
export function ensureSignalViewRealTable(cells, nav) {
  return ensureDeviceDatapointPageLayout(cells, nav)
}

/**
 * 信号层：将 ViewRealTable 配置为当前页测点行（rowSource=navDatapoints）
 * 行 = 测点；绑点格式 deviceName->dataName；仅当前页参与 QueryRealData
 */
export function remapDeviceRealtimeTable(cells, nav) {
  if (!nav || !nav.signalMode) return cells
  const cfg = buildSignalTablePageConfig(nav)
  const paged = cfg.pagedNav || applyDatapointPagination(nav)
  const names = cfg.rowDeviceNames
  const units = cfg.rowDeviceCodes
  const bindings = (cfg.bindingMatrix || []).map(row => (row && row[0]) || '').filter(Boolean).join(';')

  return (cells || []).map(cell => {
    if (!isViewRealTableCell(cell) || !cell.data || !cell.data.detail) return cell
    const next = JSON.parse(JSON.stringify(cell))
    const detail = next.data.detail
    if (!Array.isArray(detail.style.diy)) detail.style.diy = []
    const diy = detail.style.diy
    const setDiy = (key, value) => {
      const item = diy.find(d => d && d.key === key)
      if (item) item.value = value
      else diy.push({ name: key, type: 9, value, key })
    }
    setDiy('rowSource', 'navDatapoints')
    setDiy('columnHeaders', (cfg.columnHeaders || ['实时值']).join(','))
    setDiy('rowDeviceNames', names.join('\n'))
    setDiy('rowDeviceCodes', units.join('\n'))
    setDiy('rowBindings', bindings)
    // 写入分页元数据：X6 vue-shape 内 this.$store 常不可用，表格需从 diy 读总数
    setDiy('navTotalDatapoints', String(cfg.totalDatapoints || 0))
    setDiy('navDatapointPageIndex', String(cfg.datapointPageIndex || 0))
    setDiy('navDatapointPageSize', String(cfg.datapointPageSize || DEFAULT_DATAPOINT_PAGE_SIZE))
    setDiy('navDatapointTotalPages', String(cfg.datapointTotalPages || 1))
    const showCount = diy.find(d => d && d.key === 'ShowCount')
    const pageSize = cfg.datapointPageSize || DEFAULT_DATAPOINT_PAGE_SIZE
    if (showCount) {
      showCount.value = pageSize
    } else {
      diy.push({ name: 'ShowCount', type: 1, value: pageSize, key: 'ShowCount' })
    }
    return next
  })
}

function rewriteSignalPaginationCells(cells, nav) {
  if (!nav || !nav.signalMode) return
  const info = formatDatapointPageInfo(nav)
  const cur = nav.datapointPageIndex || 0
  const totalPages = nav.datapointTotalPages || 1
  cells.forEach(cell => {
    const style = cell && cell.data && cell.data.detail && cell.data.detail.style
    const t = style && typeof style.text === 'string' ? style.text : ''
    if (!t) return
    const role = detectPageNavRole(t)
    if (!role && !/\{\{\s*nav\.detailPageInfo\s*\}\}/.test(t)) return
    const detail = cell.data.detail
    if (!detail.style.diy) detail.style.diy = []
    const setRole = r => {
      const item = detail.style.diy.find(d => d && d.key === 'labelRole')
      if (item) item.value = r
      else detail.style.diy.push({ name: 'labelRole', type: 9, value: r, key: 'labelRole' })
    }
    if (role === 'pagePrev') {
      setRole(cur > 0 ? 'detailPagePrev' : 'detailPagePrevDisabled')
    } else if (role === 'pageNext') {
      setRole(cur < totalPages - 1 ? 'detailPageNext' : 'detailPageNextDisabled')
    } else if (role === 'pageInfo' || /\{\{\s*nav\.detailPageInfo\s*\}\}/.test(t)) {
      setRole('detailPageInfo')
      style.text = info
    }
  })
}

/**
 * 解析单页 cells：容器页槽位重映射 + 设备页上下文注入 + 占位符 + 动态表格
 * @param {object} components 页 components 对象（含 cells）
 * @param {object} navContext
 * @param {object} [opts]
 * @param {Record<string,string>} [opts.dpMap] 当前节点物模型测点图
 * @param {Record<string,Record<string,string>>} [opts.dpMaps] muid -> 测点图（floor 行用）
 * @param {object[]} [opts.childDevices]
 * @param {object} [opts.treeIndex] 导航树索引（槽位重映射钥匙）
 * @param {object} [opts.templateMap] 层级模板页映射
 */
export function resolveComponentsWithNavContext(components, navContext, opts = {}) {
  if (!components || !Array.isArray(components.cells) || !navContext) {
    return components
  }
  let nav = navContext
  const deviceUuid = nav.deviceUuid || nav.gatewayUuid || nav.uuid || ''
  const deviceName = nav.name || nav.label || ''
  const navKind = nav.kind || ''
  const navLayer = nav.layer || navKind
  const dpMap = opts.dpMap || {}
  let childDevices = opts.childDevices || nav.childDevices || []

  if (nav.signalMode) {
    nav = applyDatapointPagination(nav)
  }

  // 深拷贝整页，避免污染缓存模板页
  let cells
  try {
    cells = JSON.parse(JSON.stringify(components.cells))
  } catch (e) {
    return components
  }

  // 信号层：单设备测点表（navDatapoints：一行一个测点，底部分页）
  if (nav.signalMode) {
    cells = ensureDeviceDatapointPageLayout(cells, nav)
    cells = remapDeviceRealtimeTable(cells, nav)
  }

  const stats = navStats(nav, childDevices)

  // 纯设备列表：注入表格 + 翻页（无 room 模板槽位时）
  if (isDeviceListNav(nav)) {
    cells = ensureDeviceListPageLayout(cells, nav)
  }

  // 顶栏系统标题：左对齐 + 与面包屑中间留空隙（首页/列表/测点页共用）
  normalizeHeaderTitleLayout(cells)

  // 容器模板页：槽位重映射（click link / 悬浮绑点 / 槽位文本）
  if (CONTAINER_KINDS.indexOf(navKind) >= 0) {
    try {
      cells = remapContainerCells(cells, nav, {
        index: opts.treeIndex || { byOldPageId: Object.create(null), bySid: Object.create(null), parentBySid: Object.create(null) },
        templateMap: opts.templateMap || {},
        dpMaps: opts.dpMaps || {},
      })
    } catch (e) {
      console.warn('[navContextBinding] slotRemap failed:', e && e.message)
    }
  }

  // 设备列表翻页控件（room 模板无槽位时仍生效）
  if (isDeviceListNav(nav)) {
    rewritePaginationCells(cells, nav)
    layoutDeviceListPagerChrome(cells)
  }
  if (nav.signalMode) {
    rewriteSignalPaginationCells(cells, nav)
    layoutSignalPagerChrome(cells)
  }

  const homePageUuid = nav.homePageUuid || opts.homePageUuid || ''
  cells.forEach(cell => normalizeGlobalOverviewCell(cell, homePageUuid))

  const templateDeviceName = nav.signalMode ? findTemplateDeviceName(cells) : ''

  cells.forEach(cell => {
    if (!cell || !cell.data || !cell.data.detail) return
    const nextDetail = cell.data.detail

    walkConditions(nextDetail, cond => {
      if (isRelativeDeviceBinding(cond, navKind) && deviceUuid) {
        cond.deviceSN = deviceUuid
        if (deviceName) cond.DeviceName = deviceName
      }
      const dataName = cond.dataName
      if (dataName && dpMap[dataName] && (!cond.dataID || cond.dataID === '' || cond.dataID === '__NAV__')) {
        cond.dataID = dpMap[dataName]
      }
    })

    // ViewRealTable: diy rowSource=navChildren
    if (Array.isArray(nextDetail.style && nextDetail.style.diy)) {
      const diy = nextDetail.style.diy
      const rowSourceItem = diy.find(d => d && d.key === 'rowSource')
      if (rowSourceItem && String(rowSourceItem.value) === 'navChildren' && childDevices.length) {
        const names = childDevices.map(d => d.name || '')
        const codes = childDevices.map(d => d.code || d.uuid || '')
        const colHeadersItem = diy.find(d => d && d.key === 'columnHeaders')
        const headers = colHeadersItem && colHeadersItem.value
          ? String(colHeadersItem.value).split(/[\n,|]/).map(s => s.trim()).filter(Boolean)
          : ['状态']
        const bindings = childDevices.map(d =>
          headers.map(h => `${d.name || ''}->${h}`).join(',')
        ).join(';')
        const setDiy = (key, value) => {
          const item = diy.find(d => d && d.key === key)
          if (item) item.value = value
          else diy.push({ name: key, type: 9, value, key })
        }
        setDiy('rowDeviceNames', names.join('\n'))
        setDiy('rowDeviceCodes', codes.join('\n'))
        setDiy('rowBindings', bindings)
      }
    }

    // 文本占位符 {{nav.*}} + 硬编码样本改写
    applyTextPlaceholders(nextDetail.style, nav, stats)
    rewriteStaticNavTexts(nextDetail.style, nav, stats)
    rewriteStaticDeviceName(nextDetail.style, templateDeviceName, deviceName)
  })

  return sanitizeGraphComponents({ ...components, cells }, { tag: 'navContextBinding' })
}

/**
 * 异步：加载测点图并解析 components
 * @param {object} components
 * @param {object} navContext
 * @param {object} [opts] {treeIndex, templateMap}
 */
export async function resolvePageComponentsAsync(components, navContext, opts = {}) {
  if (!navContext) return components
  let nav = navContext
  if (nav.signalMode) {
    const muid = nav.modelUuid || nav.muid || ''
    const fetchParams = resolveDatapointFetchParams(nav)
    if ((muid || fetchParams.deviceUuid) && !nav.serverPaged) {
      const pointPage = await fetchDeviceDatapointPage({
        ...fetchParams,
        page: (Number(nav.datapointPageIndex) || 0) + 1,
        pageSize: nav.datapointPageSize || DEFAULT_DATAPOINT_PAGE_SIZE,
        query: nav.datapointQuery || '',
      })
      const vc = String(nav.virtualCabinet || '').trim()
      const parentLabel = String(nav.parentDeviceLabel || '').trim()
      const vcFallback = !!(
        nav.virtualCabinetFallback
        || nav.isFallbackGroup
        || (vc && parentLabel && vc === parentLabel)
      )
      nav = buildDeviceSignalContext(
        {
          label: nav.label, name: nav.name, uuid: nav.uuid || nav.deviceUuid,
          sid: nav.sid, modelUuid: muid, muid,
          virtualCabinet: vc,
          virtualCabinetFallback: vcFallback,
          parentDeviceLabel: parentLabel,
          parentDeviceUuid: nav.parentDeviceUuid || '',
        },
        pointPage.points,
        nav.ancestors || [],
        {
          ...pointPage,
          serverPaged: true,
          query: nav.datapointQuery,
          category: (vc && !vcFallback)
            ? `${vc.replace(/_+$/, '')}_`
            : nav.datapointCategory,
          virtualCabinet: vc,
          virtualCabinetFallback: vcFallback,
          parentDeviceLabel: parentLabel,
          parentDeviceUuid: nav.parentDeviceUuid || '',
        },
      )
      nav.datapointPageIndex = Math.max(0, Number(pointPage.page || 1) - 1)
    } else {
      nav = applyDatapointPagination(nav)
    }
  } else if (isDeviceListNav(nav)) {
    nav = applyDeviceListPagination(nav)
  }
  const modelUuid = nav.modelUuid || nav.muid || ''
  let dpMap = {}
  if (modelUuid) {
    dpMap = await loadDpMapByModel(modelUuid)
  }
  const childDevices = nav.childDevices || nav.childNodes || collectDeviceChildren(nav)

  // 设备列表页：仅当前页设备加载测点图（懒加载，避免 122 台并发）
  const dpMaps = {}
  const needRowDp = nav.deviceListMode || nav.kind === 'floor' || isDeviceListNav(nav)
  if (needRowDp && childDevices.length) {
    const muids = [...new Set(childDevices.map(d => d.modelUuid).filter(Boolean))]
    const maps = await Promise.all(muids.map(m => loadDpMapByModel(m)))
    muids.forEach((m, i) => { dpMaps[m] = maps[i] })
  }

  const resolved = resolveComponentsWithNavContext(components, nav, {
    dpMap,
    dpMaps,
    childDevices,
    treeIndex: opts.treeIndex || null,
    templateMap: opts.templateMap || null,
  })
  return resolved
}

export { collectDeviceChildren }
