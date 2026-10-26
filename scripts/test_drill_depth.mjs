#!/usr/bin/env node
/**
 * 自测 detectDrillDepth — 组织层（treeDepth 1~4）+ 信号层（device→datapoints）
 */
const ORG_MAX_DEPTH = 4
const SIGNAL_LAYERS = 2

function analyzeNodeChildren(node) {
  const children = (node && node.children) || []
  const deviceChildren = children.filter(c => c && isDeviceNode(c))
  const containerChildren = children.filter(c => c && !isDeviceNode(c))
  return { deviceChildren, containerChildren, total: children.length }
}

function isDeviceNode(node) {
  if (!node) return false
  if (node.type === 1) return true
  if (node.kind === 'gateway' || node.kind === 'device') return true
  return isLeafDeviceContainer(node)
}

function isLeafDeviceContainer(node) {
  if (!node || node.kind === 'device' || node.kind === 'root') return false
  return analyzeNodeChildren(node).total === 0
}

function isPureDeviceContainer(node) {
  const children = (node && node.children) || []
  return children.length > 0 && children.every(c => c && isDeviceNode(c))
}

function resolveTreeDepth(node, index) {
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

function resolveRouteMode(node) {
  if (!node) return 'org'
  if (isDeviceNode(node)) return 'signal'
  const { deviceChildren, containerChildren } = analyzeNodeChildren(node)
  if (deviceChildren.length > 0 && containerChildren.length === 0) {
    return deviceChildren.length > 1 ? 'childrenList' : 'org'
  }
  return 'org'
}

function detectDrillDepth(node, index) {
  const treeDepth = resolveTreeDepth(node, index)
  const orgDepth = Math.min(ORG_MAX_DEPTH, Math.max(1, treeDepth))
  const routeMode = resolveRouteMode(node)
  let totalLayers = orgDepth
  if (routeMode === 'signal') {
    totalLayers = orgDepth + SIGNAL_LAYERS - 1
  } else {
    totalLayers = Math.min(ORG_MAX_DEPTH + SIGNAL_LAYERS, orgDepth + SIGNAL_LAYERS)
  }
  return { treeDepth, orgDepth, routeMode, totalLayers }
}

function makeIndex(nodes) {
  const bySid = Object.create(null)
  const parentBySid = Object.create(null)
  const walk = (n, p) => {
    if (n.sid != null) { bySid[n.sid] = n; parentBySid[n.sid] = p }
    ;(n.children || []).forEach(c => walk(c, n))
  }
  nodes.forEach(r => walk(r, null))
  return { bySid, parentBySid }
}

const root = {
  kind: 'root', sid: 1, treeDepth: 1, label: 'RootZone', children: [
    { kind: 'device', sid: 1227940994, treeDepth: 2, label: '数据机房报警解析', children: [] },
    {
      kind: 'zone', sid: 1713132134, treeDepth: 2, label: 'UPS报警解析',
      children: Array.from({ length: 122 }, (_, i) => ({
        kind: 'device', sid: 1000 + i, treeDepth: 3, label: `配电室1A1_U${i}`,
      })),
    },
    {
      kind: 'zone', sid: 730335501, treeDepth: 2, label: '配电室',
      children: [
        {
          kind: 'room', sid: 1131570899, treeDepth: 3, label: '配电室_机房模块3A1',
          children: [
            { kind: 'device', sid: 573695476, treeDepth: 4, label: '机房模块3A1_1' },
          ],
        },
      ],
    },
  ],
}
const index = makeIndex([root])

const cases = [
  { name: 'RootZone', node: root, expectDepth: 1, expectRoute: 'org' },
  { name: '数据机房报警解析', node: root.children[0], expectDepth: 2, expectRoute: 'signal' },
  { name: 'UPS报警解析', node: root.children[1], expectDepth: 2, expectRoute: 'childrenList' },
  { name: '配电室', node: root.children[2], expectDepth: 2, expectRoute: 'org' },
  { name: '配电室_机房模块3A1', node: root.children[2].children[0], expectDepth: 3, expectRoute: 'org' },
  { name: '机房模块3A1_1', node: root.children[2].children[0].children[0], expectDepth: 4, expectRoute: 'signal' },
]

let pass = 0, fail = 0
for (const c of cases) {
  const d = detectDrillDepth(c.node, index)
  const ok = d.treeDepth === c.expectDepth && d.routeMode === c.expectRoute
  if (ok) pass++; else fail++
  console.log(`${ok ? '✅' : '❌'} ${c.name}: treeDepth=${d.treeDepth}(期望${c.expectDepth}) route=${d.routeMode}(期望${c.expectRoute}) total=${d.totalLayers}`)
}
console.log(`\n叶容器→设备: ${isDeviceNode(root.children[0]) ? '✅' : '❌'}`)
console.log(`纯设备容器: ${isPureDeviceContainer(root.children[1]) ? '✅' : '❌'}`)
console.log(`结果: ${pass} 通过, ${fail} 失败`)
process.exit(fail > 0 ? 1 : 0)
