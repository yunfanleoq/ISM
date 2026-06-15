<template>
  <ISMProperties />
</template>

<script>
import store from '@/store'
import ISMProperties from '@/pages/ISMDisPlay/ISMProperties.vue'
import { TOOLBOX_GROUPS, getISMComponent } from '@/pages/ISMDisPlay/componentRegistry'
import { setGroupList } from '@/store/ISM/actions'

function deepClone(value) {
  return JSON.parse(JSON.stringify(value))
}

function ensureLayerData() {
  const state = store.state.ISMDisPlayEditorTool
  if (!state.LayerData) {
    state.LayerData = { name: '--', layer: { backColor: '', backgroundImage: '', width: 1600, height: 900, autoSize: 0, Padding: 1 }, components: [] }
  }
  if (!Array.isArray(state.LayerData.components)) {
    state.LayerData.components = []
  }
  return state.LayerData
}

function ensureArrayState(state, key) {
  if (!Array.isArray(state[key])) {
    state[key] = []
  }
}

function ensureToolboxLists() {
  const state = store.state.ISMDisPlayEditorTool
  if (!Array.isArray(state.toolBoxList) || !state.toolBoxList.length) {
    state.toolBoxList = TOOLBOX_GROUPS
  }
  ensureArrayState(state, 'PageCanVasList')
  ensureArrayState(state, 'MesComponentsList')
  ensureArrayState(state, 'DiyComponentsList')
  ensureArrayState(state, 'PCPageList')
  ensureArrayState(state, 'PhonePageList')
  ensureArrayState(state, 'GroupList')
  ensureArrayState(state, 'undoStack')
  ensureArrayState(state, 'redoStack')
  ensureArrayState(state, 'selectedComponents')
  ensureArrayState(state, 'copySrcItems')
  ensureArrayState(state, 'loggerList')
  if (!state.selectedComponentMap || typeof state.selectedComponentMap !== 'object') {
    state.selectedComponentMap = {}
  }
  if (state.selectedValue === undefined || state.selectedValue === null) {
    state.selectedValue = 100
  }
  if (typeof state.selectPageUuid !== 'string') {
    state.selectPageUuid = ''
  }
  if (typeof state.selectPageContainerUuid !== 'string') {
    state.selectPageContainerUuid = ''
  }
  if (!state.PopUpConfigData) {
    state.PopUpConfigData = { name: '--', layer: { backColor: '', backgroundImage: '', width: 300, height: 600, Padding: 1 }, components: [] }
  }
  if (!state.PopUpContainerConfigData) {
    state.PopUpContainerConfigData = { name: '--', layer: { backColor: '', backgroundImage: '', width: 300, height: 600, Padding: 1 }, components: [] }
  }
}

function buildLayerComponentsFromSceneObjects(sceneObjects) {
  if (!Array.isArray(sceneObjects)) return []
  return sceneObjects
    .filter(item => item && item.type === '2dComponent' && item.source2D)
    .map(item => {
      const source2D = item.source2D
      ensureSource2DDefaults(source2D)
      return source2D
    })
}

function getLayerComponentByIdentifier(identifier) {
  const layerData = ensureLayerData()
  return layerData.components.find(item => item.identifier === identifier) || null
}

function getGroupedComponents(source2D) {
  const layerData = ensureLayerData()
  if (!source2D || !source2D.groupID) return []
  return layerData.components.filter(item => item.groupID === source2D.groupID)
}

function createCanvasCell(source2D, vm) {
  ensureSource2DDefaults(source2D)
  return {
    id: source2D.identifier || '',
    data: {
      detail: source2D,
      props: {
        dataRef: source2D.dataRef || null
      }
    },
    getData() {
      return this.data
    },
    setData(nodeData) {
      if (nodeData && nodeData.detail) {
        Object.keys(source2D).forEach(key => {
          if (!(key in nodeData.detail)) {
            vm.$delete(source2D, key)
          }
        })
        Object.keys(nodeData.detail).forEach(key => {
          vm.$set(source2D, key, nodeData.detail[key])
        })
      }
    },
    setVisible(visible) {
      source2D.style.visible = visible ? 1 : 0
    },
    isNode() {
      return true
    },
    isEdge() {
      return false
    },
    prop() {
      return {
        data: this.data
      }
    }
  }
}

function createVirtualCanvasContainer(vm) {
  const listeners = {}
  return {
    batchUpdate(fn) {
      if (typeof fn === 'function') fn()
      vm.$emit('prop-change')
      vm.$emit('scene-sync')
    },
    getCells() {
      return buildLayerComponentsFromSceneObjects(vm.sceneObjects).map(item => createCanvasCell(item, vm))
    },
    getSelectedCells() {
      const current = vm.currentObject && vm.currentObject.source2D ? vm.currentObject.source2D : null
      if (!current) return []
      return getGroupedComponents(current).map(item => createCanvasCell(item, vm))
    },
    getCellById(id) {
      const cell = getLayerComponentByIdentifier(id)
      return cell ? createCanvasCell(cell, vm) : null
    },
    select() {
      return
    },
    on(event, handler) {
      if (!listeners[event]) listeners[event] = []
      listeners[event].push(handler)
    },
    off(event, handler) {
      if (!listeners[event]) return
      if (!handler) {
        listeners[event] = []
        return
      }
      listeners[event] = listeners[event].filter(fn => fn !== handler)
    }
  }
}

function mergeMissing(target, defaults) {
  if (!defaults || typeof defaults !== 'object') return target
  if (Array.isArray(defaults)) {
    if (!Array.isArray(target)) {
      return deepClone(defaults)
    }
    if (target.length === 0 && defaults.length > 0) {
      return deepClone(defaults)
    }
    return target
  }
  if (!target || typeof target !== 'object') {
    return deepClone(defaults)
  }
  Object.keys(defaults).forEach(key => {
    const currentValue = target[key]
    const defaultValue = defaults[key]
    if (currentValue === undefined || currentValue === null) {
      target[key] = deepClone(defaultValue)
    } else if (typeof defaultValue === 'object' && defaultValue !== null && !Array.isArray(defaultValue)) {
      target[key] = mergeMissing(currentValue, defaultValue)
    } else if (Array.isArray(defaultValue) && !Array.isArray(currentValue)) {
      target[key] = deepClone(defaultValue)
    } else if (Array.isArray(defaultValue) && currentValue.length === 0 && defaultValue.length > 0) {
      target[key] = deepClone(defaultValue)
    }
  })
  return target
}

function hydrateSource2DFromBase(source2D) {
  if (!source2D || !source2D.type) return
  const component = getISMComponent(source2D.type)
  if (!component || typeof component.data !== 'function') return
  try {
    const base = component.data().base
    const baseInfo = base && base.info ? deepClone(base.info) : null
    if (!baseInfo) return
    mergeMissing(source2D, baseInfo)
  } catch (err) {
    // ignore per-component base hydration failures
  }
}

function ensureSource2DDefaults(source2D) {
  hydrateSource2DFromBase(source2D)
  if (!source2D.style) {
    source2D.style = {}
  }
  if (!source2D.style.position) {
    source2D.style.position = {}
  }
  if (!Array.isArray(source2D.style.diy)) {
    source2D.style.diy = []
  }
  if (!Array.isArray(source2D.action)) {
    source2D.action = []
  }
  if (!Array.isArray(source2D.active)) {
    source2D.active = []
  }
  if (!Array.isArray(source2D.dataBind)) {
    source2D.dataBind = []
  }
  source2D.dataBind.forEach(item => {
    if (item.DeviceName === undefined) item.DeviceName = ''
    if (item.isBandDevice === undefined) item.isBandDevice = false
    if (item.deviceSN === undefined) item.deviceSN = ''
    if (item.dataName === undefined) item.dataName = ''
    if (item.dataID === undefined) item.dataID = ''
  })
  if (!source2D.animate) {
    source2D.animate = {}
  }
  if (!Array.isArray(source2D.animate.selected)) {
    source2D.animate.selected = []
  }
  if (!Array.isArray(source2D.animate.animateList)) {
    source2D.animate.animateList = []
  }
  if (!Array.isArray(source2D.animate.animateElement)) {
    source2D.animate.animateElement = []
  }
  if (source2D.animate.isExpression === undefined) {
    source2D.animate.isExpression = false
  }
  if (!source2D.animate.condition) {
    source2D.animate.condition = {
      dataID: '',
      dataName: '',
      deviceSN: '',
      bandType: 1,
      isBandDevice: false,
      operator: '==',
      OperatorValue: '',
      OperatorMaxValue: '',
      isExpression: false,
      expression: '',
      dataUnit: '',
      DeviceName: ''
    }
  }
  if (!source2D.animate.move) {
    source2D.animate.move = {
      x: { deviceSN: '', selectVideoType: 0, isBandDevice: false, bandType: 1, dataID: '', dataName: '', DeviceName: '' },
      y: { deviceSN: '', selectVideoType: 0, isBandDevice: false, bandType: 1, dataID: '', dataName: '', DeviceName: '' }
    }
  }
  if (source2D.style.visible === undefined) {
    source2D.style.visible = 1
  }
  if (source2D.style.opacity === undefined) {
    source2D.style.opacity = 1
  }
  if (source2D.style.animate === undefined) {
    source2D.style.animate = ''
  }
  if (source2D.style.zIndex === undefined) {
    source2D.style.zIndex = -1
  }
  if (source2D.style.transform === undefined) {
    source2D.style.transform = 0
  }
  if (source2D.style.position.x === undefined) {
    source2D.style.position.x = 0
  }
  if (source2D.style.position.y === undefined) {
    source2D.style.position.y = 0
  }
  if (source2D.style.position.w === undefined) {
    source2D.style.position.w = 160
  }
  if (source2D.style.position.h === undefined) {
    source2D.style.position.h = 80
  }
  if (!source2D.style.BorderEdges && source2D.style.BorderEdges !== 0) {
    source2D.style.BorderEdges = 0
  }
  if (source2D.style.borderWidth === undefined) {
    source2D.style.borderWidth = 0
  }
  if (source2D.style.borderStyle === undefined) {
    source2D.style.borderStyle = 'solid'
  }
  if (source2D.style.borderColor === undefined) {
    source2D.style.borderColor = 'transparent'
  }
  if (source2D.style.backColor === undefined) {
    source2D.style.backColor = 'transparent'
  }
  if (source2D.style.foreColor === undefined) {
    source2D.style.foreColor = '#333333'
  }
  if (source2D.style.text === undefined) {
    source2D.style.text = ''
  }
  if (source2D.style.textAlign === undefined) {
    source2D.style.textAlign = 'center'
  }
  if (source2D.style.fontFamily === undefined) {
    source2D.style.fontFamily = 'Arial'
  }
  if (source2D.style.fontWeight === undefined) {
    source2D.style.fontWeight = 400
  }
  if (source2D.style.fontSize === undefined) {
    source2D.style.fontSize = 14
  }
  if (source2D.style.letterSpacing === undefined) {
    source2D.style.letterSpacing = 0
  }
  if (source2D.style.italic === undefined) {
    source2D.style.italic = 0
  }
  if (source2D.style.radius === undefined) {
    source2D.style.radius = 0
  }
  if (source2D.style.lineWidth === undefined) {
    source2D.style.lineWidth = 1
  }

  source2D.action.forEach(action => {
    if (!Array.isArray(action.actionAuth)) {
      action.actionAuth = []
    }
    if (action.ActionPassword === undefined) {
      action.ActionPassword = ''
    }
    if (action.ActionVoice === undefined) {
      action.ActionVoice = ''
    }
    if (action.actionVoice === undefined) {
      action.actionVoice = action.ActionVoice || ''
    }
    if (action.SecondConfirm === undefined) {
      action.SecondConfirm = false
    }
    if (action.actionConfirm === undefined) {
      action.actionConfirm = !!action.SecondConfirm
    }
    if (action.SetDelay === undefined) {
      action.SetDelay = '1000'
    }
    if (!action.link) {
      action.link = {
        linkType: 'Inside',
        Inside: { displayUUID: '', pageUUID: '' },
        isPopUp: false,
        autoClose: false,
        External: '',
        OpenExternalType: 'new',
        width: '',
        height: '',
        title: ''
      }
    }
    if (!action.DeviceView) {
      action.DeviceView = {
        key: '',
        showUUID: '',
        showPageUUID: '',
        type: '',
        selectKey: '',
        isPopUp: false,
        isContainer: false
      }
    }
    if (!action.RestApi) {
      action.RestApi = {
        Name: '',
        IsSystem: '1',
        Type: 'Post',
        Url: '',
        Params: '{}'
      }
    }
    if (!Array.isArray(action.setValue)) {
      action.setValue = []
    }
    action.setValue.forEach(setValue => {
      if (setValue.deviceSN === undefined) setValue.deviceSN = ''
      if (setValue.DeviceName === undefined) setValue.DeviceName = ''
      if (setValue.IsManual === undefined) setValue.IsManual = false
      if (setValue.AutoSetValue === undefined) setValue.AutoSetValue = ''
      if (setValue.SetPassword === undefined) setValue.SetPassword = ''
      if (setValue.isBandDevice === undefined) setValue.isBandDevice = false
      if (setValue.dataID === undefined) setValue.dataID = ''
      if (setValue.dataName === undefined) setValue.dataName = ''
    })
    if (!Array.isArray(action.ScriptList)) {
      action.ScriptList = []
    }
    if (action.animationStatus === undefined) {
      action.animationStatus = ''
    }
  })

  source2D.active.forEach(active => {
    if (active.isExpression === undefined) {
      active.isExpression = false
    }
    if (!active.condition) {
      active.condition = {}
    }
    if (active.condition.deviceSN === undefined) active.condition.deviceSN = ''
    if (active.condition.selectVideoType === undefined) active.condition.selectVideoType = 0
    if (active.condition.isBandDevice === undefined) active.condition.isBandDevice = false
    if (active.condition.bandType === undefined) active.condition.bandType = 1
    if (active.condition.dataID === undefined) active.condition.dataID = ''
    if (active.condition.dataUnit === undefined) active.condition.dataUnit = ''
    if (active.condition.dataName === undefined) active.condition.dataName = ''
    if (active.condition.DeviceName === undefined) active.condition.DeviceName = ''
    if (active.condition.operator === undefined) active.condition.operator = ''
    if (active.condition.OperatorValue === undefined) active.condition.OperatorValue = ''
    if (active.condition.OperatorMaxValue === undefined) active.condition.OperatorMaxValue = ''
    if (!Array.isArray(active.condition.actionAuth)) active.condition.actionAuth = []
    if (active.condition.ActionVoice === undefined) active.condition.ActionVoice = ''
    if (active.condition.SetPassword === undefined) active.condition.SetPassword = ''
    if (!Array.isArray(active.condition.StatusList)) active.condition.StatusList = []
    if (!Array.isArray(active.condition.PageList)) active.condition.PageList = []
    active.condition.StatusList.forEach(status => {
      if (status.Text === undefined) status.Text = 'Text'
      if (status.Image === undefined) status.Image = ''
      if (status.value === undefined) status.value = 0
      if (status.value2 === undefined) status.value2 = ''
      if (status.StatusOpt === undefined) status.StatusOpt = '=='
      if (status.TextColor === undefined) status.TextColor = '#333333'
      if (status.BlinkSpeed === undefined) status.BlinkSpeed = ''
    })
    active.condition.PageList.forEach(page => {
      if (page.MenuName === undefined) page.MenuName = ''
      if (page.IsPopUp === undefined) page.IsPopUp = false
      if (page.DisPlayID === undefined) page.DisPlayID = ''
      if (page.PageID === undefined) page.PageID = ''
    })
  })
}

function buildNodeAttrs(source2D) {
  return {
    body: {
      fill: source2D.style.backColor || 'transparent',
      stroke: source2D.style.borderColor || 'transparent',
      strokeWidth: source2D.style.borderWidth || 0,
      opacity: source2D.style.opacity === undefined ? 1 : source2D.style.opacity
    },
    text: {
      text: source2D.style.text || '',
      fill: source2D.style.foreColor || '#333333',
      fontSize: source2D.style.fontSize || 14,
      fontWeight: source2D.style.fontWeight || 'normal'
    },
    line: {},
    wrap: {}
  }
}

function isEdgeLikeComponent(source2D) {
  const type = ((source2D && source2D.type) || '').toLowerCase()
  if (!type) return false
  return (
    type.includes('line') ||
    type.includes('arrow') ||
    type.includes('dashed') ||
    type.includes('conduit')
  )
}

function ensureEdgeDiyDefaults(source2D) {
  if (!Array.isArray(source2D.style.diy)) {
    source2D.style.diy = []
  }
  const ensureDiy = (key, value) => {
    const found = source2D.style.diy.find(item => item.key === key)
    if (!found) {
      source2D.style.diy.push({ key, value })
    } else if (found.value === undefined || found.value === null || found.value === '') {
      found.value = value
    }
  }
  ensureDiy('strokeWidth', source2D.style.lineWidth || 4)
  ensureDiy('MoveBrokenLineInterval', 3)
  ensureDiy('strokeLength', 10)
  ensureDiy('strokeSpace', 5)
  ensureDiy('MoveBrokenLineConditionEnable', 0)
  ensureDiy('strokeBgWidth', 25)
  ensureDiy('strokeLinejoin', 1)
  ensureDiy('strokeLineType', 0)
  ensureDiy('strokeLineMarkerStyle', 0)
  ensureDiy('strokeLineMarker', 0)
  ensureDiy('strokeLineMarkerColor', source2D.style.foreColor || '#06ad8e')
  ensureDiy('strokeLineMarkerWidth', 10)
  ensureDiy('strokeLineMarkerHeight', 10)
}

function upsertDiyValue(source2D, key, value) {
  if (!Array.isArray(source2D.style.diy)) {
    source2D.style.diy = []
  }
  const found = source2D.style.diy.find(item => item.key === key)
  if (found) {
    found.value = value
  } else {
    source2D.style.diy.push({ key, value })
  }
}

function createVirtualNode(vm) {
  const target = vm.currentObject
  return {
    prop() {
      ensureSource2DDefaults(target.source2D)
      if (isEdgeLikeComponent(target.source2D)) {
        ensureEdgeDiyDefaults(target.source2D)
      }
      const style = target.source2D && target.source2D.style ? target.source2D.style : {}
      const pos = style.position || {}
      const propData = {
        angle: style.transform || 0,
        visible: style.visible === undefined ? true : style.visible === 1,
        position: { x: pos.x || 0, y: pos.y || 0 },
        size: { width: pos.w || 160, height: pos.h || 80 },
        attrs: buildNodeAttrs(target.source2D),
        data: {
          detail: target.source2D,
          props: {
            dataRef: target.source2D.dataRef || null
          }
        }
      }
      if (isEdgeLikeComponent(target.source2D)) {
        propData.connectType = 'connected-edge'
      }
      return propData
    },
    getData() {
      return {
        detail: target.source2D,
        props: {
          dataRef: target.source2D.dataRef || null
        }
      }
    },
    setData(nodeData) {
      if (nodeData && nodeData.detail) {
        ensureSource2DDefaults(nodeData.detail)
        const nextDetail = deepClone(nodeData.detail)
        Object.keys(target.source2D || {}).forEach(key => {
          if (!(key in nextDetail)) {
            vm.$delete(target.source2D, key)
          }
        })
        Object.keys(nextDetail).forEach(key => {
          vm.$set(target.source2D, key, nextDetail[key])
        })
        vm.syncBridgeState()
        vm.$emit('prop-change')
        vm.$emit('scene-sync')
      }
    },
    setZIndex(z) {
      vm.$set(target.source2D.style, 'zIndex', z)
      vm.syncBridgeState()
    },
    rotate(angle) {
      vm.$set(target.source2D.style, 'transform', angle)
      vm.syncBridgeState()
    },
    setVisible(visible) {
      vm.$set(target.source2D.style, 'visible', visible ? 1 : 0)
      vm.syncBridgeState()
    },
    position(x, y) {
      vm.$set(target.source2D.style.position, 'x', x)
      vm.$set(target.source2D.style.position, 'y', y)
      vm.syncBridgeState()
    },
    resize(width, height) {
      vm.$set(target.source2D.style.position, 'w', width)
      vm.$set(target.source2D.style.position, 'h', height)
      vm.syncBridgeState()
    },
    isEdge() {
      return isEdgeLikeComponent(target.source2D)
    },
    setAttrs(attrs) {
      const source2D = target.source2D
      ensureSource2DDefaults(source2D)
      if (isEdgeLikeComponent(source2D)) {
        ensureEdgeDiyDefaults(source2D)
      }
      if (attrs && attrs.body) {
        if (attrs.body.fill !== undefined) vm.$set(source2D.style, 'backColor', attrs.body.fill)
        if (attrs.body.stroke !== undefined) vm.$set(source2D.style, 'borderColor', attrs.body.stroke)
        if (attrs.body.strokeWidth !== undefined) vm.$set(source2D.style, 'borderWidth', attrs.body.strokeWidth)
        if (attrs.body.opacity !== undefined) vm.$set(source2D.style, 'opacity', attrs.body.opacity)
      }
      if (attrs && attrs.text) {
        if (attrs.text.text !== undefined) vm.$set(source2D.style, 'text', attrs.text.text)
        if (attrs.text.fill !== undefined) vm.$set(source2D.style, 'foreColor', attrs.text.fill)
        if (attrs.text.fontSize !== undefined) vm.$set(source2D.style, 'fontSize', attrs.text.fontSize)
        if (attrs.text.fontWeight !== undefined) vm.$set(source2D.style, 'fontWeight', attrs.text.fontWeight)
      }
      if (attrs && attrs.line) {
        if (attrs.line.strokeOpacity !== undefined) vm.$set(source2D.style, 'opacity', attrs.line.strokeOpacity)
        const lineStyle = attrs.line.style || {}
        if (lineStyle.stroke !== undefined) vm.$set(source2D.style, 'foreColor', lineStyle.stroke)
        if (lineStyle.strokeWidth !== undefined) {
          vm.$set(source2D.style, 'lineWidth', lineStyle.strokeWidth)
          upsertDiyValue(source2D, 'strokeWidth', lineStyle.strokeWidth)
        }
        if (lineStyle.strokeDasharray !== undefined) {
          const dashText = `${lineStyle.strokeDasharray}`.trim()
          const dashParts = dashText.split(/\s+/)
          upsertDiyValue(source2D, 'strokeLength', dashParts[0] || '')
          upsertDiyValue(source2D, 'strokeSpace', dashParts[1] || '')
          upsertDiyValue(source2D, 'strokeLineType', dashText ? 0 : 1)
        }
        if (lineStyle.strokeLinejoin !== undefined) {
          const joinMap = { miter: 0, round: 1, bevel: 2 }
          upsertDiyValue(source2D, 'strokeLinejoin', joinMap[lineStyle.strokeLinejoin] !== undefined ? joinMap[lineStyle.strokeLinejoin] : 1)
        }
        if (lineStyle.animation !== undefined) {
          const animationText = `${lineStyle.animation}`
          if (!animationText) {
            upsertDiyValue(source2D, 'MoveBrokenLineConditionEnable', 1)
          } else {
            upsertDiyValue(source2D, 'MoveBrokenLineConditionEnable', 0)
            const match = animationText.match(/(\d+(?:\.\d+)?)s/)
            if (match) {
              upsertDiyValue(source2D, 'MoveBrokenLineInterval', Number(match[1]))
            }
          }
        }
        if (attrs.line.sourceMarker || attrs.line.targetMarker) {
          const marker = attrs.line.sourceMarker || attrs.line.targetMarker
          const markerMap = {
            classic: 0,
            diamond: 1,
            cross: 2,
            circle: 3,
            circlePlus: 4,
            ellipse: 5
          }
          upsertDiyValue(source2D, 'strokeLineMarker', marker && marker.name ? 1 : 0)
          if (marker && marker.name !== undefined) {
            upsertDiyValue(source2D, 'strokeLineMarkerStyle', markerMap[marker.name] !== undefined ? markerMap[marker.name] : 0)
          }
          if (marker && marker.fill !== undefined) {
            upsertDiyValue(source2D, 'strokeLineMarkerColor', marker.fill)
          }
          if (marker && marker.width !== undefined) {
            upsertDiyValue(source2D, 'strokeLineMarkerWidth', marker.width)
          }
          if (marker && marker.height !== undefined) {
            upsertDiyValue(source2D, 'strokeLineMarkerHeight', marker.height)
          }
        }
      }
      if (attrs && attrs.wrap) {
        if (attrs.wrap.stroke !== undefined) vm.$set(source2D.style, 'backColor', attrs.wrap.stroke)
        if (attrs.wrap.strokeWidth !== undefined) upsertDiyValue(source2D, 'strokeBgWidth', attrs.wrap.strokeWidth)
        if (attrs.wrap.strokeLinejoin !== undefined) {
          const joinMap = { miter: 0, round: 1, bevel: 2 }
          upsertDiyValue(source2D, 'strokeLinejoin', joinMap[attrs.wrap.strokeLinejoin] !== undefined ? joinMap[attrs.wrap.strokeLinejoin] : 1)
        }
      }
      vm.syncBridgeState()
      vm.$emit('prop-change')
      vm.$emit('scene-sync')
    },
    setProp() {},
    getChildren() {
      const source2D = target.source2D
      return getGroupedComponents(source2D)
        .filter(item => item.identifier !== source2D.identifier)
        .map(item => ({
          setVisible(visible) {
            ensureSource2DDefaults(item)
            item.style.visible = visible ? 1 : 0
          },
          getData() {
            return {
              detail: item,
              props: {
                dataRef: item.dataRef || null
              }
            }
          },
          setData(nodeData) {
            if (nodeData && nodeData.detail) {
              Object.assign(item, deepClone(nodeData.detail))
            }
          }
        }))
    }
  }
}

export default {
  name: 'ISM2DPropertiesBridge',
  components: {
    ISMProperties
  },
  props: {
    currentObject: { type: Object, required: true },
    sceneObjects: { type: Array, default: () => [] }
  },
  data() {
    return {
      syncingSelectedComponent: false
    }
  },
  computed: {
    selectedComponent() {
      return store.state.ISMDisPlayEditorTool.selectedComponent
    }
  },
  watch: {
    currentObject: {
      immediate: true,
      handler() {
        this.syncBridgeState()
      }
    },
    sceneObjects: {
      handler() {
        this.syncBridgeState()
      }
    },
    selectedComponent: {
      deep: true,
      handler(newVal) {
        this.syncSelectedComponentToScene(newVal)
      }
    }
  },
  methods: {
    buildSelectedComponentMap(source2D) {
      if (!source2D) return {}
      const items = getGroupedComponents(source2D)
      if (!items.length) {
        return { [source2D.identifier]: source2D }
      }
      return items.reduce((acc, item) => {
        acc[item.identifier] = item
        return acc
      }, {})
    },
    syncBridgeState() {
      if (!this.currentObject || !this.currentObject.source2D) return
      const state = store.state.ISMDisPlayEditorTool
      const layerData = ensureLayerData()
      ensureToolboxLists()
      const source2D = this.currentObject.source2D
      ensureSource2DDefaults(source2D)
      const layerComponents = buildLayerComponentsFromSceneObjects(this.sceneObjects)
      layerData.components = layerComponents
      const virtualNode = createVirtualNode(this)
      const nodeProps = virtualNode.prop()
      nodeProps.data = {
        detail: source2D,
        props: {
          dataRef: source2D.dataRef || null
        }
      }

      const selectedComponentMap = this.buildSelectedComponentMap(source2D)
      const selectedComponents = Object.keys(selectedComponentMap)

      state.selectedIsLayer = false
      state.selectedComponent = source2D
      state.selectedComponents = selectedComponents
      state.selectedComponentMap = selectedComponentMap
      state.selectedNode = virtualNode
      state.selectedNodePops = nodeProps
      state.ISMCavasContainer = createVirtualCanvasContainer(this)
      setGroupList({ state })
    },
    syncSelectedComponentToScene(source2D) {
      if (this.syncingSelectedComponent) return
      if (!this.currentObject || !this.currentObject.source2D || !source2D) return
      if (source2D.identifier && this.currentObject.source2D.identifier && source2D.identifier !== this.currentObject.source2D.identifier) return
      this.syncingSelectedComponent = true
      ensureSource2DDefaults(source2D)
      Object.keys(this.currentObject.source2D).forEach(key => {
        if (!(key in source2D)) this.$delete(this.currentObject.source2D, key)
      })
      Object.keys(source2D).forEach(key => {
        this.$set(this.currentObject.source2D, key, deepClone(source2D[key]))
      })
      this.$emit('scene-sync')
      this.$nextTick(() => {
        this.syncingSelectedComponent = false
      })
    }
  }
}
</script>
