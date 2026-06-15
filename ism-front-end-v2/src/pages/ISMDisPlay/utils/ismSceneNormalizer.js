import { uuid } from 'vue-uuid'

const DEFAULT_CONDITION = {
  deviceSN: '',
  selectVideoType: 0,
  isBandDevice: false,
  bandType: 1,
  dataID: '',
  dataName: '',
  operator: '',
  OperatorValue: '',
  OperatorMaxValue: ''
}

function toNumber(value, fallback) {
  const number = Number(value)
  return Number.isFinite(number) ? number : fallback
}

function clone(value) {
  return JSON.parse(JSON.stringify(value || {}))
}

function readCellPosition(cell, sourcePosition) {
  const graphPosition = cell && cell.position ? cell.position : {}
  return {
    x: toNumber(graphPosition.x, toNumber(cell && cell.x, toNumber(sourcePosition && sourcePosition.x, 0))),
    y: toNumber(graphPosition.y, toNumber(cell && cell.y, toNumber(sourcePosition && sourcePosition.y, 0)))
  }
}

function readCellSize(cell, sourcePosition) {
  const graphSize = cell && cell.size ? cell.size : {}
  return {
    width: toNumber(graphSize.width, toNumber(cell && cell.width, toNumber(sourcePosition && sourcePosition.w, 100))),
    height: toNumber(graphSize.height, toNumber(cell && cell.height, toNumber(sourcePosition && sourcePosition.h, 40)))
  }
}

export function createISMNode({
  shape,
  name,
  x = 0,
  y = 0,
  width = 100,
  height = 40,
  zIndex = -1,
  style = {},
  detail = {},
  data = {},
  attrs,
  ports
}) {
  const safeShape = shape || detail.type || 'view-svg-text'
  const safeX = toNumber(x, 0)
  const safeY = toNumber(y, 0)
  const safeWidth = toNumber(width, 100)
  const safeHeight = toNumber(height, 40)
  const safeZIndex = toNumber(zIndex, toNumber(style.zIndex, -1))
  const detailStyle = {
    position: {
      x: safeX,
      y: safeY,
      w: safeWidth,
      h: safeHeight
    },
    visible: 1,
    backColor: 'transparent',
    foreColor: '#d7ecff',
    fontWeight: 400,
    zIndex: safeZIndex,
    transform: 0,
    text: name || safeShape,
    textAlign: 'center',
    fontSize: 18,
    fontFamily: 'Arial',
    letterSpacing: 0,
    italic: 0,
    borderWidth: 0,
    BorderEdges: 0,
    opacity: 1,
    borderStyle: 'solid',
    borderColor: '#1f8fff',
    diy: [],
    ...clone(style),
    position: {
      ...(style && style.position ? style.position : {}),
      x: safeX,
      y: safeY,
      w: safeWidth,
      h: safeHeight
    },
    zIndex: safeZIndex
  }
  const nodeDetail = {
    identifier: uuid.v1(),
    name: name || safeShape,
    type: safeShape,
    action: [],
    dataBind: [],
    active: [],
    animate: {
      selected: [],
      condition: { ...DEFAULT_CONDITION },
      isExpression: false,
      animateList: [
        { id: 'blink', name: 'component.public.animateBlink' },
        { id: 'Zoom', name: 'component.public.Zoom' },
        { id: 'animateSpin', name: 'component.public.animateSpin' }
      ],
      animateElement: []
    },
    ...clone(detail),
    type: safeShape,
    style: detailStyle
  }
  return {
    id: uuid.v1(),
    shape: safeShape,
    x: safeX,
    y: safeY,
    width: safeWidth,
    height: safeHeight,
    position: { x: safeX, y: safeY },
    size: { width: safeWidth, height: safeHeight },
    visible: true,
    zIndex: safeZIndex,
    ...(attrs ? { attrs } : {}),
    ...(ports ? { ports } : {}),
    data: {
      locked: false,
      UpdateNodeFlag: true,
      editMode: true,
      showDeviceUuid: '',
      IsToolBox: false,
      ...clone(data),
      detail: nodeDetail
    }
  }
}

export function normalizeISMCell(cell, options = {}) {
  if (!cell || typeof cell !== 'object' || !cell.shape) {
    return null
  }
  const allowedShapes = options.allowedShapes || []
  const isEdge = cell.shape === 'edge'
  if (!isEdge && allowedShapes.length > 0 && !allowedShapes.includes(cell.shape)) {
    return null
  }
  const sourceData = clone(cell.data)
  const sourceDetail = clone(sourceData.detail)
  const sourceStyle = clone(sourceDetail.style)
  const sourcePosition = sourceStyle.position || {}
  const cellPosition = readCellPosition(cell, sourcePosition)
  const cellSize = readCellSize(cell, sourcePosition)
  const x = cellPosition.x
  const y = cellPosition.y
  const width = cellSize.width
  const height = cellSize.height
  const zIndex = toNumber(cell.zIndex, toNumber(sourceStyle.zIndex, options.fallbackZIndex || -1))

  if (isEdge) {
    return {
      ...cell,
      zIndex,
      data: {
        locked: false,
        UpdateNodeFlag: true,
        editMode: true,
        showDeviceUuid: '',
        IsToolBox: false,
        ...sourceData
      }
    }
  }

  return {
    ...cell,
    id: cell.id || uuid.v1(),
    shape: cell.shape,
    x,
    y,
    width,
    height,
    position: { x, y },
    size: { width, height },
    visible: cell.visible !== false,
    zIndex,
    data: {
      locked: false,
      UpdateNodeFlag: true,
      editMode: true,
      showDeviceUuid: '',
      IsToolBox: false,
      ...sourceData,
      detail: {
        identifier: sourceDetail.identifier || uuid.v1(),
        name: sourceDetail.name || cell.name || cell.shape,
        type: cell.shape,
        action: Array.isArray(sourceDetail.action) ? sourceDetail.action : [],
        dataBind: Array.isArray(sourceDetail.dataBind) ? sourceDetail.dataBind : [],
        active: Array.isArray(sourceDetail.active) ? sourceDetail.active : [],
        ...sourceDetail,
        animate: {
          selected: sourceDetail.animate && sourceDetail.animate.selected || [],
          condition: (sourceDetail.animate && sourceDetail.animate.condition) || { ...DEFAULT_CONDITION },
          isExpression: sourceDetail.animate ? !!sourceDetail.animate.isExpression : false,
          animateList: (sourceDetail.animate && sourceDetail.animate.animateList) || [],
          animateElement: (sourceDetail.animate && sourceDetail.animate.animateElement) || []
        },
        type: cell.shape,
        style: {
          visible: 1,
          backColor: 'transparent',
          foreColor: '#d7ecff',
          fontWeight: 400,
          transform: 0,
          textAlign: 'center',
          fontSize: 18,
          fontFamily: 'Arial',
          letterSpacing: 0,
          italic: 0,
          borderWidth: 0,
          BorderEdges: 0,
          opacity: 1,
          borderStyle: 'solid',
          borderColor: '#1f8fff',
          diy: [],
          ...sourceStyle,
          position: { ...sourcePosition, x, y, w: width, h: height },
          zIndex
        }
      }
    }
  }
}

export function normalizeISMScene(scene, options = {}) {
  const source = scene || {}
  const sourceComponents = source.components || source.graph || {}
  const cells = Array.isArray(sourceComponents.cells) ? sourceComponents.cells : []
  console.log('[normalizeISMScene] input cells:', cells.length, 'shapes:', cells.map(c => c?.shape || 'NOSHAPE').slice(0,5))
  const normalizedCells = cells
    .map(cell => normalizeISMCell(cell, options))
    .filter(Boolean)
  console.log('[normalizeISMScene] output cells:', normalizedCells.length)
  const layer = source.layer || {}
  return {
    layer: {
      ...clone(layer),
      width: toNumber(layer.width, options.width || 1920),
      height: toNumber(layer.height, options.height || 1080),
      backColor: layer.backColor !== undefined ? layer.backColor : (options.backColor || '#061321'),
      backgroundImage: layer.backgroundImage !== undefined ? layer.backgroundImage : ''
    },
    components: {
      ...sourceComponents,
      cells: normalizedCells
    }
  }
}
