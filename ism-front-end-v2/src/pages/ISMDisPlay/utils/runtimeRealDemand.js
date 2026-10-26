/**
 * 旧组态运行态实时值名单。
 * 只保留当前画面里看得到的对象、以及鼠标正放着的对象。
 * 翻页、离开时由调用方换掉这份名单，上一批请求作废。
 */

function finite(value, fallback) {
  const n = Number(value)
  if (Number.isFinite(n)) return n
  const f = Number(fallback)
  return Number.isFinite(f) ? f : 0
}

export function cellBox(cell) {
  const pos = (cell && cell.position) || {}
  const size = (cell && cell.size) || {}
  const stylePos = cell && cell.data && cell.data.detail && cell.data.detail.style
    ? cell.data.detail.style.position
    : null
  return {
    x: finite(cell && cell.x, finite(pos.x, stylePos && stylePos.x)),
    y: finite(cell && cell.y, finite(pos.y, stylePos && stylePos.y)),
    w: finite(cell && cell.width, finite(size.width, stylePos && stylePos.w)),
    h: finite(cell && cell.height, finite(size.height, stylePos && stylePos.h)),
  }
}

export function isCellShown(cell) {
  if (!cell) return false
  if (cell.visible === false || cell.visible === 0) return false
  const style = cell.data && cell.data.detail && cell.data.detail.style
  if (!style || style.visible === undefined || style.visible === null) return true
  return style.visible !== 0 && style.visible !== false
}

export function boxIntersects(box, view, pad) {
  if (!box || !view) return true
  if (!(box.w > 0) || !(box.h > 0)) return true
  const p = pad || 0
  const right = box.x + box.w
  const bottom = box.y + box.h
  return right >= view.x - p
    && box.x <= view.x + view.w + p
    && bottom >= view.y - p
    && box.y <= view.y + view.h + p
}

export function readGraphViewport(graph) {
  if (!graph || typeof graph.scale !== 'function' || !graph.container) return null
  const scale = graph.scale() || {}
  const translated = typeof graph.translate === 'function' ? (graph.translate() || {}) : {}
  const sx = Number(scale.sx) || 1
  const sy = Number(scale.sy) || 1
  const w = graph.container.clientWidth
  const h = graph.container.clientHeight
  if (!w || !h) return null
  return {
    x: -(Number(translated.tx) || 0) / sx,
    y: -(Number(translated.ty) || 0) / sy,
    w: w / sx,
    h: h / sy,
  }
}

function pushPair(pairs, systemUuids, seen, dataID, deviceSN) {
  if (typeof dataID !== 'string' || !dataID) return
  const device = typeof deviceSN === 'string' ? deviceSN : ''
  const key = device + '\n' + dataID
  if (seen.has(key)) return
  seen.add(key)
  if (device) pairs.push({ uuid: dataID, device })
  else systemUuids.push(dataID)
}

function addCellPairs(cell, pairs, systemUuids, seen) {
  const detail = cell && cell.data && cell.data.detail
  if (!detail) return
  const active = detail.active
  if (active && active.length) {
    for (let i = 0; i < active.length; i++) {
      const condition = active[i] && active[i].condition
      if (!condition) continue
      pushPair(pairs, systemUuids, seen, condition.dataID, condition.deviceSN)
    }
  }
  const animate = detail && detail.animate
  if (!animate || typeof animate !== 'object') return
  if (animate.condition) {
    pushPair(pairs, systemUuids, seen, animate.condition.dataID, animate.condition.deviceSN)
  }
  const move = animate.move
  if (!move || typeof move !== 'object') return
  if (move.x) pushPair(pairs, systemUuids, seen, move.x.dataID, move.x.deviceSN)
  if (move.y) pushPair(pairs, systemUuids, seen, move.y.dataID, move.y.deviceSN)
}

function isHovered(cell, hoverId) {
  if (!hoverId || !cell) return false
  if (cell.id === hoverId) return true
  const identifier = cell.data && cell.data.detail && cell.data.detail.identifier
  return identifier === hoverId
}

/**
 * cells 必须是当前这一页（或当前分页容器里正在显示的那一页）。
 * viewport 有值时，画面外的对象不进名单；hoverId 对应的对象不论是否在画面内都进名单。
 */
export function selectRuntimeRealDemand(cells, { viewport = null, hoverId = '', pad = 48 } = {}) {
  const pairs = []
  const systemUuids = []
  const seen = new Set()
  const list = Array.isArray(cells) ? cells : []
  for (let i = 0; i < list.length; i++) {
    const cell = list[i]
    if (!cell) continue
    const hovered = isHovered(cell, hoverId)
    if (!hovered) {
      if (!isCellShown(cell)) continue
      if (viewport && !boxIntersects(cellBox(cell), viewport, pad)) continue
    }
    addCellPairs(cell, pairs, systemUuids, seen)
  }
  return { pairs, systemUuids }
}

export function demandSignature(demand) {
  const pairs = (demand && demand.pairs) || []
  const systemUuids = (demand && demand.systemUuids) || []
  const paired = []
  for (let i = 0; i < pairs.length; i++) {
    paired.push(pairs[i].device + '\n' + pairs[i].uuid)
  }
  paired.sort()
  const system = systemUuids.slice().sort()
  return paired.join('\t') + '\f' + system.join('\t')
}

/** 把一次 getRealDataByUuid 的结果收成按设备推送的 readDataPush 载荷。空值不进名单。 */
export function groupRealDataPush(rows) {
  const pushDataMap = {}
  const uuids = []
  const modelUuids = []
  const list = Array.isArray(rows) ? rows : []
  for (let k = 0; k < list.length; k++) {
    const row = list[k]
    if (!row || row.value === '') continue
    const duid = row.duid
    if (!pushDataMap[duid]) {
      pushDataMap[duid] = {
        DeviceUuid: duid,
        ProjectUuid: row.project_uuid,
        Cmd: 'RealData',
        Data: [],
      }
    }
    uuids.push(row.uuid)
    modelUuids.push(row.mduid)
    pushDataMap[duid].Data.push({
      Uuid: row.uuid,
      ModelDataUuid: row.mduid,
      Value: row.value,
    })
  }
  return { pushDataMap, uuids, modelUuids }
}
