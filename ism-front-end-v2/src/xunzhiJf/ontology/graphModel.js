import { BRIDGE_LINKS, ORG_LINKS } from './orgTree'

export const LAYER_COLOR = {
  physical: '#3ec6e0',
  business: '#f0a202',
  system: '#3ddc97',
  bridge: '#b388ff',
  gap: '#ff6b7a',
}

export const LAYER_LABEL = {
  physical: '物理世界',
  business: '业务与工作流',
  system: '系统世界',
  bridge: '桥',
  gap: '缺口',
}

const STRUCTURAL = [
  'Module', 'Page', 'Workflow', 'Phase', 'Agent', 'Skill', 'Tool', 'ProtocolDriver',
  'Room', 'Device', 'UPS', 'ColumnHeader', 'PrecisionAC', 'Rack', 'Alarm', 'AlarmPolicy',
  'WorkOrder', 'ITAsset', 'CapacityAccount', 'PueSnapshot', 'Gap', 'User', 'Role',
]

export function layerColor(layer, gap) {
  if (gap) return LAYER_COLOR.gap
  return LAYER_COLOR[layer] || '#8aa0b8'
}

function relNodeId(relation) {
  return 'rel:' + relation
}

const BRIDGE_ANCHOR = {
  physical: ['Room', 'Device', 'UPS', 'PrecisionAC', 'Point'],
  business: ['Alarm', 'WorkOrder', 'CapacityAccount', 'PueSnapshot'],
  system: ['Module', 'Page', 'Workflow', 'Agent'],
}

function pickAnchors(entities, layer) {
  const want = BRIDGE_ANCHOR[layer] || []
  const got = []
  want.forEach(type => {
    const hit = entities.find(e => e.type === type)
    if (hit) got.push(hit)
  })
  if (got.length) return got
  return entities.filter(e => e.layer === layer).slice(0, 4)
}

export function schemaToGraph(schema, layer) {
  const entities = ((schema && schema.entity_types) || []).filter(e => !layer || layer === 'all' || e.layer === layer)
  const relations = ((schema && schema.relation_types) || []).filter(r => !layer || layer === 'all' || r.layer === layer)
  const present = {}
  entities.forEach(e => { present[e.type] = e })
  const nodes = entities.map(e => ({
    id: e.type,
    label: e.label || e.type,
    color: e.color || layerColor(e.layer),
    layer: e.layer,
    kind: 'entity',
  }))
  const edges = []
  relations.forEach(rel => {
    const fromStar = !rel.from || rel.from.indexOf('*') >= 0
    const toStar = !rel.to || rel.to.indexOf('*') >= 0
    if (!fromStar && !toStar) {
      rel.from.forEach(src => {
        rel.to.forEach(dst => {
          if (present[src] && present[dst]) {
            edges.push({ from: src, to: dst, label: rel.label || rel.relation })
          }
        })
      })
      return
    }
    const layer = rel.layer || 'bridge'
    const id = relNodeId(rel.relation)
    nodes.push({
      id,
      label: rel.label || rel.relation,
      color: rel.color || layerColor(layer),
      layer,
      kind: 'relation',
    })
    if (layer === 'bridge') {
      ;['physical', 'business', 'system'].forEach(side => {
        pickAnchors(entities, side).forEach(e => {
          edges.push({ from: id, to: e.type, label: rel.label || rel.relation, layer: 'bridge' })
        })
      })
      return
    }
    const pool = entities.filter(e => e.layer === layer)
    pool.forEach(e => edges.push({ from: id, to: e.type, label: rel.label || rel.relation, layer }))
  })
  return { nodes, edges }
}

export function instancesToGraph(rawNodes, rawEdges, opts) {
  const options = opts || {}
  const showPoints = !!options.showPoints
  let list = (rawNodes || []).filter(n => showPoints || n.node_type !== 'Point')
  const rank = type => {
    const i = STRUCTURAL.indexOf(type)
    return i < 0 ? 100 : i
  }
  list = list.slice().sort((a, b) => rank(a.node_type) - rank(b.node_type))
  const cap = options.cap || 160
  const drawn = list.slice(0, cap)
  const keep = {}
  drawn.forEach(n => { keep[n.node_key] = true })
  const nodes = drawn.map(n => ({
    id: n.node_key,
    label: n.label || n.node_type || n.node_key,
    color: layerColor(n.layer, n.gap),
    layer: n.gap ? 'gap' : n.layer,
    kind: 'instance',
    raw: n,
  }))
  const edges = []
  ;(rawEdges || []).forEach(e => {
    if (keep[e.from] && keep[e.to] && edges.length < 420) {
      edges.push({ from: e.from, to: e.to, label: e.relation })
    }
  })
  return {
    nodes,
    edges,
    hidden: Math.max(0, (rawNodes || []).length - drawn.length),
    total: (rawNodes || []).length,
  }
}

export function overviewNodes(counts) {
  return ['physical', 'business', 'system'].map(layer => ({
    id: 'layer:' + layer,
    label: LAYER_LABEL[layer] + (counts && counts[layer] != null ? ' · ' + counts[layer] : ''),
    layer,
    kind: 'layer',
    color: LAYER_COLOR[layer],
  }))
}

function typeLabel(schema, type) {
  const ent = ((schema && schema.entity_types) || []).find(e => e.type === type)
  return (ent && ent.label) || type
}

function focusRing(center, others, edges) {
  const nodes = [Object.assign({}, center, { focus: true })].concat(others.map(n => Object.assign({}, n, { focus: false })))
  return { nodes, edges, hidden: 0, total: nodes.length }
}

export function drillSchema(schema, stack) {
  const full = schemaToGraph(schema, 'all')
  if (!stack.length) {
    const counts = {}
    full.nodes.filter(n => n.kind === 'entity').forEach(n => {
      counts[n.layer] = (counts[n.layer] || 0) + 1
    })
    return { nodes: overviewNodes(counts), edges: [], hidden: 0, total: 3 }
  }
  const last = stack[stack.length - 1]
  if (last.kind === 'layer') {
    const org = buildOrgSchema(schema)
    const nodes = org.nodes.filter(n => n.layer === last.layer)
    const ids = {}
    nodes.forEach(n => { ids[n.id] = true })
    const edges = org.edges.filter(e => ids[e.from] && ids[e.to])
    return { nodes, edges, layout: 'org', hidden: 0, total: nodes.length }
  }
  const org = buildOrgSchema(schema)
  const center = org.nodes.find(n => n.id === last.id) || full.nodes.find(n => n.id === last.id) || {
    id: last.id, label: last.label, layer: last.layer, kind: last.kind || 'entity', color: layerColor(last.layer),
  }
  const edges = org.edges.filter(e => e.from === center.id || e.to === center.id)
  if (!edges.length) {
    const fallback = full.edges.filter(e => e.from === center.id || e.to === center.id)
    const ids = {}
    fallback.forEach(e => { ids[e.from] = true; ids[e.to] = true })
    delete ids[center.id]
    return focusRing(center, full.nodes.filter(n => ids[n.id]), fallback)
  }
  const ids = {}
  edges.forEach(e => { ids[e.from] = true; ids[e.to] = true })
  const nodes = org.nodes.filter(n => ids[n.id])
  if (!nodes.some(n => n.id === center.id)) nodes.unshift(center)
  return { nodes, edges, layout: 'org', hidden: 0, total: nodes.length }
}

export function drillInstances(rawNodes, rawEdges, stack, showPoints, schema) {
  const all = rawNodes || []
  const usable = all.filter(n => showPoints || n.node_type !== 'Point')
  const byKey = {}
  all.forEach(n => { byKey[n.node_key] = n })
  if (!stack.length) {
    const counts = { physical: 0, business: 0, system: 0 }
    usable.forEach(n => {
      if (counts[n.layer] != null && !n.gap) counts[n.layer] += 1
    })
    return { nodes: overviewNodes(counts), edges: [], hidden: 0, total: 3 }
  }
  const last = stack[stack.length - 1]
  if (last.kind === 'layer') {
    const inLayer = last.layer === 'gap'
      ? all.filter(n => n.gap || n.node_type === 'Gap')
      : usable.filter(n => n.layer === last.layer && !n.gap)
    const grouped = {}
    inLayer.forEach(n => {
      const t = n.node_type || 'Node'
      if (!grouped[t]) grouped[t] = 0
      grouped[t] += 1
    })
    const nodes = Object.keys(grouped).map(t => ({
      id: 'type:' + last.layer + ':' + t,
      label: typeLabel(schema, t) + ' · ' + grouped[t],
      layer: last.layer,
      kind: 'type',
      nodeType: t,
      color: layerColor(last.layer),
    }))
    const rootId = 'layer:' + last.layer
    nodes.unshift({
      id: rootId,
      label: LAYER_LABEL[last.layer] || last.layer,
      layer: last.layer,
      kind: 'layer',
      core: true,
      color: layerColor(last.layer),
    })
    const present = {}
    nodes.forEach(n => { present[n.id] = true })
    const edges = []
    const linked = {}
    ORG_LINKS.filter(link => link.layer === last.layer).forEach(link => {
      const from = link.parent.indexOf('layer:') === 0 ? link.parent : ('type:' + last.layer + ':' + link.parent)
      const to = 'type:' + last.layer + ':' + link.child
      if (!present[from] || !present[to]) return
      edges.push({ from, to, label: link.rel, layer: last.layer, org: true })
      linked[to] = true
    })
    nodes.forEach(n => {
      if (n.id === rootId || linked[n.id]) return
      edges.push({ from: rootId, to: n.id, label: '归入', layer: last.layer, org: true })
    })
    return { nodes, edges, layout: 'org', hidden: 0, total: nodes.length }
  }
  if (last.kind === 'type') {
    const list = usable.filter(n => n.layer === last.layer && n.node_type === last.nodeType)
    const drawn = list.slice(0, 36)
    const keep = {}
    const nodes = drawn.map(n => ({
      id: n.node_key,
      label: n.label || n.node_type,
      color: layerColor(n.layer, n.gap),
      layer: n.layer,
      kind: 'instance',
      nodeType: n.node_type,
      raw: n,
    }))
    nodes.forEach(n => { keep[n.id] = true })
    const edges = []
    ;(rawEdges || []).forEach(e => {
      if (keep[e.from] && keep[e.to]) edges.push({ from: e.from, to: e.to, label: e.relation })
    })
    return { nodes, edges, hidden: Math.max(0, list.length - drawn.length), total: list.length }
  }
  const centerRaw = byKey[last.id]
  const center = {
    id: last.id,
    label: (centerRaw && centerRaw.label) || last.label || last.id,
    color: layerColor((centerRaw && centerRaw.layer) || last.layer, centerRaw && centerRaw.gap),
    layer: (centerRaw && centerRaw.layer) || last.layer,
    kind: 'instance',
    nodeType: centerRaw && centerRaw.node_type,
    raw: centerRaw,
    focus: true,
  }
  const edges = []
  const others = []
  let pointCount = 0
  const seen = {}
  ;(rawEdges || []).forEach(e => {
    if (e.from !== last.id && e.to !== last.id) return
    const otherId = e.from === last.id ? e.to : e.from
    const other = byKey[otherId]
    if (!other || seen[otherId]) return
    seen[otherId] = true
    if (!showPoints && other.node_type === 'Point') {
      pointCount += 1
      return
    }
    edges.push({ from: e.from, to: e.to, label: e.relation })
    others.push({
      id: other.node_key,
      label: other.label || other.node_type,
      color: layerColor(other.layer, other.gap),
      layer: other.gap ? 'gap' : other.layer,
      kind: 'instance',
      nodeType: other.node_type,
      raw: other,
    })
  })
  if (pointCount) {
    const pid = 'points:' + last.id
    others.push({
      id: pid,
      label: '测点 · ' + pointCount,
      color: LAYER_COLOR.physical,
      layer: 'physical',
      kind: 'type',
      nodeType: 'Point',
    })
  }
  return focusRing(center, others.slice(0, 28), edges)
}

export function applyLegend(graph, showLayer) {
  if (!graph || !showLayer) return graph
  const nodes = (graph.nodes || []).filter(n => showLayer[n.layer] !== false)
  const ids = {}
  nodes.forEach(n => { ids[n.id] = true })
  const edges = (graph.edges || []).filter(e => {
    if (!ids[e.from] || !ids[e.to]) return false
    if (e.layer && showLayer[e.layer] === false) return false
    return true
  })
  return Object.assign({}, graph, { nodes, edges })
}

function entityNode(ent) {
  return {
    id: ent.type,
    label: ent.label || ent.type,
    color: ent.color || layerColor(ent.layer),
    layer: ent.layer,
    kind: 'entity',
  }
}

export function buildOrgSchema(schema) {
  const entities = (schema && schema.entity_types) || []
  const byType = {}
  const nodes = []
  entities.forEach(ent => {
    const node = entityNode(ent)
    byType[ent.type] = node
    nodes.push(node)
  })
  const seen = {}
  nodes.forEach(n => { seen[n.id] = true })
  ;['physical', 'business', 'system'].forEach(layer => {
    if (!nodes.some(n => n.layer === layer)) return
    const id = 'layer:' + layer
    if (seen[id]) return
    nodes.push({
      id,
      label: LAYER_LABEL[layer],
      layer,
      kind: 'layer',
      color: LAYER_COLOR[layer],
      core: true,
    })
    seen[id] = true
  })
  const edges = []
  const linked = {}
  ORG_LINKS.forEach(link => {
    if (!seen[link.parent] || !seen[link.child]) return
    edges.push({
      from: link.parent,
      to: link.child,
      label: link.rel,
      layer: link.layer,
      org: true,
    })
    linked[link.child] = true
  })
  nodes.forEach(n => {
    if (n.kind !== 'entity' || linked[n.id]) return
    const core = 'layer:' + n.layer
    if (!seen[core]) return
    edges.push({ from: core, to: n.id, label: '归入', layer: n.layer, org: true })
    linked[n.id] = true
  })
  BRIDGE_LINKS.forEach(link => {
    if (!seen[link.from] || !seen[link.to]) return
    edges.push({ from: link.from, to: link.to, label: link.rel, layer: 'bridge' })
  })
  return { nodes, edges, layout: 'org', hidden: 0, total: nodes.length }
}

export function wholeSchema(schema, showLayer) {
  return applyLegend(buildOrgSchema(schema), showLayer)
}

export function wholeInstances(rawNodes, rawEdges, showPoints, showLayer) {
  const rank = type => {
    const i = STRUCTURAL.indexOf(type)
    return i < 0 ? 100 : i
  }
  let list = (rawNodes || []).filter(n => showPoints || n.node_type !== 'Point')
  list = list.filter(n => {
    const layer = n.gap || n.node_type === 'Gap' ? 'gap' : (n.layer || 'system')
    return !showLayer || showLayer[layer] !== false
  })
  list = list.slice().sort((a, b) => rank(a.node_type) - rank(b.node_type)).slice(0, 96)
  const keep = {}
  const nodes = list.map(n => {
    const layer = n.gap || n.node_type === 'Gap' ? 'gap' : (n.layer || 'system')
    keep[n.node_key] = true
    return {
      id: n.node_key,
      label: n.label || n.node_type || n.node_key,
      color: layerColor(layer === 'gap' ? n.layer : layer, layer === 'gap'),
      layer,
      kind: 'instance',
      nodeType: n.node_type,
      raw: n,
    }
  })
  const layerOf = {}
  nodes.forEach(n => { layerOf[n.id] = n.layer })
  const edges = []
  ;(rawEdges || []).forEach(e => {
    if (keep[e.from] && keep[e.to] && edges.length < 320) {
      const fromLayer = layerOf[e.from]
      const toLayer = layerOf[e.to]
      edges.push({
        from: e.from,
        to: e.to,
        label: e.relation,
        layer: fromLayer && fromLayer === toLayer ? fromLayer : 'bridge',
      })
    }
  })
  return {
    nodes,
    edges,
    layout: 'layers',
    hidden: Math.max(0, (rawNodes || []).length - nodes.length),
    total: (rawNodes || []).length,
  }
}
