import { createISMNode } from './ismSceneNormalizer'

function createAiNode(shape, name, x, y, width, height, style = {}) {
  return createISMNode({
    shape,
    name,
    x,
    y,
    width,
    height,
    zIndex: typeof style.zIndex === 'number' ? style.zIndex : -1,
    style: {
      text: style.text || name,
      ...style
    }
  })
}

function createAiPanel(name, x, y, width, height, title) {
  return [
    createAiNode('dv-border-box1', name, x, y, width, height, { zIndex: 1 }),
    createAiNode('view-svg-text', title, x + 16, y + 12, width - 32, 32, {
      text: title,
      foreColor: '#8ed8ff',
      fontSize: 20,
      fontWeight: 600,
      zIndex: 2
    })
  ]
}

export function buildLocalAiSceneGraph({ prompt, width, height, append }) {
  const offsetX = append ? 40 : 0
  const offsetY = append ? 40 : 0
  const w = Math.max(width, 960)
  const h = Math.max(height, 540)
  const title = prompt.indexOf('配电') !== -1 ? '配电室智能监控' : prompt.indexOf('空调') !== -1 ? '空调机房智能监控' : '工业设备智能监控'
  const cells = []
  cells.push(createAiNode('view-svg-text', 'AI场景标题', offsetX + w * 0.28, offsetY + 24, w * 0.44, 56, {
    text: title,
    foreColor: '#e9f7ff',
    fontSize: 34,
    fontWeight: 700,
    zIndex: 10
  }))
  cells.push(...createAiPanel('左侧设备状态面板', offsetX + 40, offsetY + 110, w * 0.25, h * 0.34, '设备状态'))
  cells.push(...createAiPanel('右侧实时数据面板', offsetX + w * 0.72, offsetY + 110, w * 0.24, h * 0.34, '实时数据'))
  cells.push(...createAiPanel('底部趋势分析面板', offsetX + 40, offsetY + h * 0.62, w * 0.56, h * 0.28, '趋势分析'))
  cells.push(...createAiPanel('底部报警面板', offsetX + w * 0.63, offsetY + h * 0.62, w * 0.33, h * 0.28, '报警事件'))

  const centerX = offsetX + w * 0.35
  const centerY = offsetY + h * 0.18
  const cardWidth = w * 0.13
  const cardHeight = 84
  const deviceNames = prompt.indexOf('配电') !== -1
    ? ['主变压器', '高压柜', '低压柜', '馈线柜', '母联柜', '电容柜']
    : ['冷水机组', '循环水泵', '冷却塔', '新风机组', '空压机', '末端阀门']
  deviceNames.forEach((deviceName, index) => {
    const col = index % 3
    const row = Math.floor(index / 3)
    const x = centerX + col * (cardWidth + 28)
    const y = centerY + row * (cardHeight + 42)
    cells.push(createAiNode('view-svg-text', deviceName, x, y, cardWidth, cardHeight, {
      text: `${deviceName}\n运行正常`,
      backColor: 'rgba(8, 40, 70, 0.78)',
      foreColor: index === 2 ? '#ffd666' : '#73f5c4',
      fontSize: 18,
      fontWeight: 600,
      borderWidth: 1,
      BorderEdges: 8,
      borderColor: index === 2 ? '#d6a20f' : '#1fbf9b',
      zIndex: 5
    }))
  })

  cells.push(createAiNode('view-device-real-data-table', 'AI实时数据表', offsetX + w * 0.735, offsetY + 170, w * 0.205, h * 0.21, {
    foreColor: '#d7ecff',
    fontSize: 14,
    zIndex: 6
  }))
  cells.push(createAiNode('view-svg-text', 'AI趋势占位', offsetX + 80, offsetY + h * 0.7, w * 0.48, h * 0.12, {
    text: '趋势图区域\n可替换为实时/历史曲线组件',
    backColor: 'rgba(7, 33, 56, 0.72)',
    foreColor: '#8ed8ff',
    fontSize: 22,
    BorderEdges: 6,
    borderWidth: 1,
    borderColor: '#1f8fff',
    zIndex: 6
  }))
  cells.push(createAiNode('view-svg-text', 'AI报警列表占位', offsetX + w * 0.66, offsetY + h * 0.7, w * 0.27, h * 0.12, {
    text: '暂无高优先级报警\n设备巡检正常',
    backColor: 'rgba(62, 30, 18, 0.68)',
    foreColor: '#ffd666',
    fontSize: 20,
    BorderEdges: 6,
    borderWidth: 1,
    borderColor: '#d6a20f',
    zIndex: 6
  }))
  return {
    layer: {
      width: w,
      height: h,
      backColor: '#061321',
      backgroundImage: ''
    },
    components: {
      cells
    }
  }
}
