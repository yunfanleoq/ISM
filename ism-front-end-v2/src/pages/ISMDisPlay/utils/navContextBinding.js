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
  fetchDeviceDatapoints,
  buildDeviceSignalContext,
  buildSignalTablePageConfig,
  formatDatapointPageInfo,
  isSignalDatapointPaged,
  DEFAULT_DATAPOINT_PAGE_SIZE,
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
    ;(n.children || []).forEach(walk)
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

function isViewRealTableCell(cell) {
  if (!cell) return false
  if (cell.shape === 'ism-view-real-table') return true
  const type = cell.data && cell.data.detail && cell.data.detail.type
  return type === 'ism-view-real-table'
}

/**
 * 纯设备列表（childrenList）：无 room 模板槽位时注入标题 + ViewRealTable(navChildren) + 翻页控件
 */
export function ensureDeviceListPageLayout(cells, nav) {
  if (!isDeviceListNav(nav)) return cells || []
  const list = [...(cells || [])]
  const pageSize = deviceListPageSizeForNav(nav)
  const uid = nav.sid != null ? nav.sid : (nav.uuid || 'list')

  if (!list.some(isViewRealTableCell)) {
    const id = `devlist-vrt-${uid}`
    list.push(makeTemplateTableCell(
      id, TPL_MAIN_X, TPL_TABLE_Y, TPL_MAIN_W, TPL_TABLE_H,
      [
        { name: 'rowSource', type: 9, value: 'navChildren', key: 'rowSource' },
        { name: 'columnHeaders', type: 9, value: '在线状态', key: 'columnHeaders' },
        { name: 'ShowCount', type: 1, value: pageSize, key: 'ShowCount' },
        { name: 'tableHeaderColor', type: 2, value: '#f8fbff', key: 'tableHeaderColor' },
        { name: 'tableHeaderBackColor', type: 2, value: '#1d3557', key: 'tableHeaderBackColor' },
        { name: 'tableSplitColor', type: 2, value: '#263449', key: 'tableSplitColor' },
        { name: 'tableHoverColor', type: 2, value: '#1e3a5f', key: 'tableHoverColor' },
        { name: 'themeName', type: 6, value: 'dark', key: 'themeName' },
      ],
    ))
  }

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
      mkNavText(`devlist-prev-${uid}`, TPL_MAIN_X, '‹ 上一页', 'pagePrev'),
      mkNavText(`devlist-info-${uid}`, TPL_MAIN_X + 96, '{{nav.pageInfo}}', 'pageInfo'),
      mkNavText(`devlist-next-${uid}`, TPL_MAIN_X + 396, '下一页 ›', 'pageNext'),
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
  const list = [...(cells || [])]
  const pageSize = nav.datapointPageSize || DEFAULT_DATAPOINT_PAGE_SIZE
  const uid = nav.deviceUuid || nav.uuid || nav.sid || 'device'

  if (!list.some(isViewRealTableCell)) {
    const id = `devpoints-vrt-${uid}`
    list.push(makeTemplateTableCell(
      id, TPL_MAIN_X, TPL_TABLE_Y, TPL_MAIN_W, TPL_TABLE_H,
      [
        { name: 'rowSource', type: 9, value: 'navDatapoints', key: 'rowSource' },
        { name: 'columnHeaders', type: 9, value: '实时值', key: 'columnHeaders' },
        { name: 'ShowCount', type: 1, value: pageSize, key: 'ShowCount' },
        { name: 'tableHeaderColor', type: 2, value: '#f8fbff', key: 'tableHeaderColor' },
        { name: 'tableHeaderBackColor', type: 2, value: '#1d3557', key: 'tableHeaderBackColor' },
        { name: 'tableSplitColor', type: 2, value: '#263449', key: 'tableSplitColor' },
        { name: 'tableHoverColor', type: 2, value: '#1e3a5f', key: 'tableHoverColor' },
        { name: 'themeName', type: 6, value: 'dark', key: 'themeName' },
      ],
      { name: '设备测点表' },
    ))
  }

  const hasTitle = list.some(c => {
    const text = ((((c || {}).data || {}).detail || {}).style || {}).text || ''
    return /\{\{\s*nav\.name\s*\}\}/.test(String(text))
  })
  if (!hasTitle) {
    list.unshift(makeTemplateTextCell(
      `devpoints-title-${uid}`, TPL_MAIN_X, TPL_TITLE_Y, TPL_MAIN_W, 36,
      '{{nav.name}}', { foreColor: '#00e5ff', fontSize: 20, name: '设备标题' },
    ))
  }

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
  if (!isSignalDatapointPaged(nav)) return
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
  }
  if (isSignalDatapointPaged(nav)) {
    rewriteSignalPaginationCells(cells, nav)
  }

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
    const devLabel = nav.label || nav.name || ''
    if (muid && !(nav.allDatapoints && nav.allDatapoints.length)) {
      let points = await fetchDeviceDatapoints(muid, devLabel, nav.uuid || nav.deviceUuid || '')
      if (!points.length && devLabel) {
        const short = String(devLabel).split('_').filter(Boolean).pop()
        if (short && short !== devLabel) {
          points = await fetchDeviceDatapoints(muid, short, nav.uuid || nav.deviceUuid || '')
        }
      }
      nav = buildDeviceSignalContext(
        {
          label: nav.label, name: nav.name, uuid: nav.uuid || nav.deviceUuid,
          sid: nav.sid, modelUuid: muid, muid,
        },
        points,
        nav.ancestors || [],
      )
      nav.datapointPageIndex = navContext.datapointPageIndex || 0
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
