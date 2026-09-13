import { uuid } from 'vue-uuid'

export const GROUP_SHAPE = 'view-ism-group-node'

export function withGraphBatch(graph, name, fn) {
  if (graph && typeof graph.batchUpdate === 'function') {
    return graph.batchUpdate(name, fn)
  }
  if (graph && typeof graph.startBatch === 'function') {
    graph.startBatch(name)
    try {
      return fn()
    } finally {
      graph.stopBatch(name)
    }
  }
  return fn()
}

export function isGroupNode(cell) {
  return !!(cell && cell.isNode && cell.isNode() && cell.prop && cell.prop().shape === GROUP_SHAPE)
}

/** 选区提升到顶层组，子节点去重，避免全选时组和内容被当成平级。 */
export function normalizeTopLevelNodes(cells) {
  const seen = {}
  const result = []
  ;(cells || []).forEach((cell) => {
    if (!cell || !cell.isNode || !cell.isNode()) return
    let top = cell
    let parent = typeof cell.getParent === 'function' ? cell.getParent() : null
    while (parent && parent.isNode && parent.isNode()) {
      top = parent
      parent = typeof parent.getParent === 'function' ? parent.getParent() : null
    }
    const id = top.id
    if (!id || seen[id]) return
    seen[id] = true
    result.push(top)
  })
  return result
}

export function selectedTopLevelNodes(graph) {
  if (!graph || typeof graph.getSelectedCells !== 'function') return []
  return normalizeTopLevelNodes(graph.getSelectedCells())
}

export function getNodeAngle(node) {
  const raw = node && node.prop ? node.prop().angle : 0
  const n = parseInt(raw, 10)
  return Number.isFinite(n) ? n : 0
}

function collectSelfAndDescendants(node) {
  const list = [node]
  const children = node && typeof node.getChildren === 'function' ? node.getChildren() : null
  ;(children || []).forEach((child) => {
    if (child && child.isNode && child.isNode()) {
      collectSelfAndDescendants(child).forEach((item) => {
        list.push(item)
      })
    }
  })
  return list
}

function leafNodes(top) {
  return collectSelfAndDescendants(top).filter((node) => {
    const children = typeof node.getChildren === 'function' ? node.getChildren() : null
    return !children || !children.length
  })
}

function fitGroupBox(group) {
  const children = (typeof group.getChildren === 'function' ? group.getChildren() : null) || []
  const nodes = children.filter((cell) => cell && cell.isNode && cell.isNode())
  if (!nodes.length) return
  const bboxes = nodes.map((node) => node.getBBox())
  const left = Math.min.apply(null, bboxes.map((b) => b.x))
  const top = Math.min.apply(null, bboxes.map((b) => b.y))
  const right = Math.max.apply(null, bboxes.map((b) => b.x + b.width))
  const bottom = Math.max.apply(null, bboxes.map((b) => b.y + b.height))
  group.setPosition(left, top)
  if (typeof group.resize === 'function') {
    group.resize(right - left, bottom - top)
  }
  const data = typeof group.getData === 'function' ? group.getData() : null
  if (data && data.detail && data.detail.style && data.detail.style.position) {
    const next = Object.assign({}, data)
    next.detail = Object.assign({}, data.detail)
    next.detail.style = Object.assign({}, data.detail.style)
    next.detail.style.position = Object.assign({}, data.detail.style.position, {
      w: right - left,
      h: bottom - top,
    })
    group.setData(next, { overwrite: true })
  }
}

function fitGroupTree(top) {
  collectSelfAndDescendants(top)
    .filter(isGroupNode)
    .reverse()
    .forEach(fitGroupBox)
}

function unionCenter(nodes) {
  let minX = Infinity
  let minY = Infinity
  let maxX = -Infinity
  let maxY = -Infinity
  nodes.forEach((node) => {
    const b = node.getBBox()
    minX = Math.min(minX, b.x)
    minY = Math.min(minY, b.y)
    maxX = Math.max(maxX, b.x + b.width)
    maxY = Math.max(maxY, b.y + b.height)
  })
  return { cx: (minX + maxX) / 2, cy: (minY + maxY) / 2 }
}

function toggleStyleTransform(node, magic) {
  const data = (node && typeof node.getData === 'function' ? node.getData() : null) || {}
  if (!data.detail) return
  const style = data.detail.style || {}
  const next = Object.assign({}, data)
  next.detail = Object.assign({}, data.detail)
  next.detail.style = Object.assign({}, style)
  next.detail.style.transform = style.transform == magic ? 0 : magic
  next.UpdateNodeFlag = !data.UpdateNodeFlag
  node.setData(next, { overwrite: true })
}

export function rotateSelectedNodes(graph, delta) {
  const tops = selectedTopLevelNodes(graph)
  if (!tops.length) return 0
  const rad = (delta * Math.PI) / 180
  const cos = Math.cos(rad)
  const sin = Math.sin(rad)
  withGraphBatch(graph, 'rotate', () => {
    tops.forEach((top) => {
      const all = leafNodes(top)
      const { cx, cy } = unionCenter(all.length ? all : [top])
      all.forEach((node) => {
        const bbox = node.getBBox()
        const px = bbox.x + bbox.width / 2
        const py = bbox.y + bbox.height / 2
        const nx = cx + (px - cx) * cos - (py - cy) * sin
        const ny = cy + (px - cx) * sin + (py - cy) * cos
        const pos = node.getPosition()
        node.setPosition(pos.x + (nx - px), pos.y + (ny - py))
        let angle = getNodeAngle(node) + delta
        if (angle >= 360 || angle <= -360) {
          angle = angle % 360
        }
        node.rotate(angle, { absolute: true })
      })
      fitGroupTree(top)
    })
  })
  return tops.length
}

/** axis: 'x' 左右镜像（旧 -1099/rotateY），'y' 上下镜像（旧 -1098/rotateX） */
export function flipSelectedNodes(graph, axis) {
  const tops = selectedTopLevelNodes(graph)
  if (!tops.length) return 0
  const magic = axis === 'x' ? -1099 : -1098
  withGraphBatch(graph, 'flip', () => {
    tops.forEach((top) => {
      const all = leafNodes(top)
      const { cx, cy } = unionCenter(all.length ? all : [top])
      all.forEach((node) => {
        const bbox = node.getBBox()
        const px = bbox.x + bbox.width / 2
        const py = bbox.y + bbox.height / 2
        const pos = node.getPosition()
        if (axis === 'x') {
          node.setPosition(pos.x + 2 * (cx - px), pos.y)
        } else {
          node.setPosition(pos.x, pos.y + 2 * (cy - py))
        }
        toggleStyleTransform(node, magic)
      })
      fitGroupTree(top)
    })
  })
  return tops.length
}

export function alignSelectedNodes(graph, type) {
  const nodes = selectedTopLevelNodes(graph)
  if (nodes.length < 2) return 0
  withGraphBatch(graph, 'align-' + type, () => {
    if (type === 'l') {
      const sorted = nodes.slice().sort((a, b) => a.getBBox().x - b.getBBox().x)
      const baseX = sorted[0].getBBox().x
      sorted.slice(1).forEach((node) => {
        const pos = node.getPosition()
        const dx = baseX - node.getBBox().x
        node.setPosition(pos.x + dx, pos.y)
      })
      return
    }
    if (type === 'r') {
      const sorted = nodes.slice().sort((a, b) => {
        const ba = a.getBBox()
        const bb = b.getBBox()
        return (bb.x + bb.width) - (ba.x + ba.width)
      })
      const base = sorted[0].getBBox()
      const baseRight = base.x + base.width
      sorted.slice(1).forEach((node) => {
        const pos = node.getPosition()
        const bbox = node.getBBox()
        const dx = baseRight - (bbox.x + bbox.width)
        node.setPosition(pos.x + dx, pos.y)
      })
      return
    }
    if (type === 't') {
      const sorted = nodes.slice().sort((a, b) => a.getBBox().y - b.getBBox().y)
      const baseY = sorted[0].getBBox().y
      sorted.slice(1).forEach((node) => {
        const pos = node.getPosition()
        const dy = baseY - node.getBBox().y
        node.setPosition(pos.x, pos.y + dy)
      })
      return
    }
    if (type === 'b') {
      const sorted = nodes.slice().sort((a, b) => {
        const ba = a.getBBox()
        const bb = b.getBBox()
        return (bb.y + bb.height) - (ba.y + ba.height)
      })
      const base = sorted[0].getBBox()
      const baseBottom = base.y + base.height
      sorted.slice(1).forEach((node) => {
        const pos = node.getPosition()
        const bbox = node.getBBox()
        const dy = baseBottom - (bbox.y + bbox.height)
        node.setPosition(pos.x, pos.y + dy)
      })
    }
  })
  return nodes.length
}

export function arrangeSelectedNodes(graph, direction) {
  const nodes = selectedTopLevelNodes(graph)
  if (nodes.length < 2) return 0
  withGraphBatch(graph, 'arrange-' + direction, () => {
    if (direction === 'Vertical') {
      const sorted = nodes.slice().sort((a, b) => a.getBBox().y - b.getBBox().y)
      const baseNode = sorted[0]
      const maxY = sorted[sorted.length - 1].getBBox().y
      const minY = baseNode.getBBox().y
      const totalHeight = sorted.slice(1).reduce((sum, node) => sum + node.getSize().height, 0)
      const spacing = (maxY - minY - totalHeight) / (sorted.length - 1)
      let currentY = minY + baseNode.getSize().height + spacing
      sorted.slice(1).forEach((node) => {
        const pos = node.getPosition()
        node.setPosition(pos.x, currentY)
        currentY += node.getSize().height + spacing
      })
      return
    }
    const sorted = nodes.slice().sort((a, b) => a.getBBox().x - b.getBBox().x)
    const baseNode = sorted[0]
    const maxX = sorted[sorted.length - 1].getBBox().x
    const minX = baseNode.getBBox().x
    const totalWidth = sorted.slice(1).reduce((sum, node) => sum + node.getSize().width, 0)
    const spacing = (maxX - minX - totalWidth) / (sorted.length - 1)
    let currentX = minX + baseNode.getSize().width + spacing
    sorted.slice(1).forEach((node) => {
      node.setPosition(currentX, baseNode.getPosition().y)
      currentX += node.getSize().width + spacing
    })
  })
  return nodes.length
}

function buildDefaultGroupNodeData(width, height) {
  return {
    locked: false,
    UpdateNodeFlag: true,
    editMode: true,
    showDeviceUuid: '',
    IsToolBox: false,
    detail: {
      identifier: uuid.v1(),
      name: '节点组',
      type: 'image',
      isCanvas: true,
      action: [],
      dataBind: [],
      active: [
        {
          id: 'Forward',
          name: 'component.ViewCanvasMoveLineArrow.Forward',
          result: '',
          isExpression: true,
          condition: {
            deviceSN: '',
            selectVideoType: 0,
            isBandDevice: false,
            bandType: 1,
            dataID: '',
            dataName: '',
            operator: '',
            OperatorValue: '',
            OperatorMaxValue: '',
          },
        },
        {
          id: 'Reverse',
          name: 'component.ViewCanvasMoveLineArrow.Reverse',
          result: '',
          isExpression: true,
          condition: {
            deviceSN: '',
            selectVideoType: 0,
            isBandDevice: false,
            bandType: 1,
            dataID: '',
            dataName: '',
            operator: '',
            OperatorValue: '',
            OperatorMaxValue: '',
          },
        },
        {
          id: 'ControlStatus',
          name: 'configComponent.status.ControlStatus',
          result: 0,
          isStatus: true,
          isSwitch: false,
          isImageStatus: false,
          isTextStatus: false,
          isLineStatus: true,
          isExpression: false,
          condition: {
            deviceSN: '',
            isBandDevice: false,
            bandType: 1,
            dataID: '',
            dataName: '',
            IsManual: false,
            StatusList: [
              {
                StatusOpt: '==',
                TextColor: '#d81e06',
                Blink: '0',
                BlinkSpeed: 1,
                value2: 1,
                value: 1,
              },
            ],
          },
        },
      ],
      animate: {
        selected: [],
        condition: {
          deviceSN: '',
          selectVideoType: 0,
          isBandDevice: false,
          bandType: 1,
          dataID: '',
          dataName: '',
          operator: '',
          OperatorValue: '',
          OperatorMaxValue: '',
        },
        isExpression: false,
        animateList: [],
        animateElement: [
          {
            id: 'blink',
            elementList: [
              {
                name: 'component.public.animateSpeed',
                type: 7,
                value: 1,
                min: 0.1,
                key: 'blinkSpeed',
              },
            ],
          },
          {
            id: 'millcolorGrad',
            elementList: [
              {
                name: 'component.public.startColor',
                type: 2,
                value: '#74f808',
                key: 'startColor',
              },
              {
                name: 'component.public.stopColor',
                type: 2,
                value: '#f30b0b',
                key: 'stopColor',
              },
              {
                name: 'component.public.animateSpeed',
                type: 7,
                value: 1,
                min: 0.1,
                key: 'animateSpeed',
              },
            ],
          },
          {
            id: 'animateSpin',
            elementList: [
              {
                name: 'component.public.animateSpinSpeed',
                type: 7,
                value: 1,
                min: 0.1,
                key: 'spinSpeed',
              },
              {
                name: 'configComponent.bigScreen.border.border89Direction',
                type: 6,
                value: 0,
                enumList: [
                  {
                    value: 0,
                    option: 'configComponent.bigScreen.border.border89DirectionForward',
                  },
                  {
                    value: 1,
                    option: 'configComponent.bigScreen.border.border89DirectionNegative',
                  },
                ],
                min: 1,
                key: 'spinDirection',
              },
            ],
          },
        ],
      },
      style: {
        position: { x: 0, y: 0, w: width, h: height },
        points: [],
        visible: 1,
        zIndex: -1000,
        transform: 0,
        backColor: '',
        foreColor: '',
        borderWidth: 2,
        BorderEdges: 0,
        opacity: 1,
        borderStyle: 'solid',
        borderColor: '#13c2c2',
        diy: [],
      },
    },
  }
}

export function createGroupFromCells(graph, cells) {
  const members = normalizeTopLevelNodes(cells)
  if (!graph || !members.length) return null
  if (members.length === 1 && isGroupNode(members[0])) return members[0]
  return withGraphBatch(graph, 'create-group', () => {
    const bboxes = members.map((cell) => cell.getBBox())
    const left = Math.min.apply(null, bboxes.map((b) => b.x))
    const top = Math.min.apply(null, bboxes.map((b) => b.y))
    const right = Math.max.apply(null, bboxes.map((b) => b.x + b.width))
    const bottom = Math.max.apply(null, bboxes.map((b) => b.y + b.height))
    const width = right - left
    const height = bottom - top
    const parent = graph.addNode({
      shape: GROUP_SHAPE,
      x: left,
      y: top,
      width,
      height,
      zIndex: -1000,
      data: buildDefaultGroupNodeData(width, height),
    })
    members.forEach((node) => {
      parent.addChild(node)
    })
    return parent
  })
}
