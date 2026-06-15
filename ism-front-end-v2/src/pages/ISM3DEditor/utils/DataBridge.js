/**
 * DataBridge - 数据桥接器
 * 负责与工业数据源连接，处理数据绑定和更新
 */
import * as THREE from 'three'

/**
 * 数据转换函数
 */
const TransformFunctions = {
  // 直接映射
  direct: (value) => value,

  // 缩放
  scale: (value, params) => value * (params.scale || 1),

  // 偏移
  offset: (value, params) => value + (params.offset || 0),

  // 范围映射
  range: (value, params) => {
    const { inMin = 0, inMax = 100, outMin = 0, outMax = 1 } = params
    const clamped = Math.min(Math.max(value, inMin), inMax)
    return ((clamped - inMin) * (outMax - outMin)) / (inMax - inMin) + outMin
  },

  // 颜色映射
  color: (value, params) => {
    const { inMin = 0, inMax = 100, colorLow = '#00ff00', colorHigh = '#ff0000' } = params
    const t = (Math.min(Math.max(value, inMin), inMax) - inMin) / (inMax - inMin)
    const c1 = new THREE.Color(colorLow)
    const c2 = new THREE.Color(colorHigh)
    c1.lerp(c2, t)
    return '#' + c1.getHexString()
  }
}

/**
 * 属性更新函数
 */
const PropertyUpdaters = {
  'position.x': (object, value) => { object.position.x = value },
  'position.y': (object, value) => { object.position.y = value },
  'position.z': (object, value) => { object.position.z = value },
  'rotation.x': (object, value) => { object.rotation.x = THREE.MathUtils.degToRad(value) },
  'rotation.y': (object, value) => { object.rotation.y = THREE.MathUtils.degToRad(value) },
  'rotation.z': (object, value) => { object.rotation.z = THREE.MathUtils.degToRad(value) },
  'scale.x': (object, value) => { object.scale.x = value },
  'scale.y': (object, value) => { object.scale.y = value },
  'scale.z': (object, value) => { object.scale.z = value },
  'appearance.color': (object, value) => {
    if (object.material && object.material.color) {
      object.material.color.set(value)
    }
  },
  'appearance.opacity': (object, value) => {
    if (object.material) {
      object.material.opacity = value
      object.material.transparent = value < 1
      object.material.needsUpdate = true
    }
  },
  'visible': (object, value) => { object.visible = value }
}

/**
 * 数据桥接器主类
 */
export class DataBridge {
  constructor() {
    this.dataSources = new Map()
    this.bindings = new Map()
    this.objects = new Map()
    this.connections = new Map()
    this.enabled = true
  }

  /**
   * 添加数据源
   */
  addDataSource(id, config) {
    let source

    switch (config.type) {
      case 'websocket':
        source = new WebSocketDataSource(id, config)
        break
      case 'mqtt':
        source = new MQTTDataSource(id, config)
        break
      case 'http':
        source = new HTTPDataSource(id, config)
        break
      case 'static':
        source = new StaticDataSource(id, config)
        break
      default:
        console.warn(`Unknown data source type: ${config.type}`)
        return null
    }

    this.dataSources.set(id, source)
    return source
  }

  /**
   * 移除数据源
   */
  removeDataSource(id) {
    const source = this.dataSources.get(id)
    if (source) {
      source.disconnect()
      this.dataSources.delete(id)
    }
  }

  /**
   * 添加绑定
   */
  addBinding(objectId, binding) {
    if (!this.bindings.has(objectId)) {
      this.bindings.set(objectId, [])
    }
    this.bindings.get(objectId).push(binding)
  }

  /**
   * 移除绑定
   */
  removeBinding(objectId, bindingId) {
    const bindings = this.bindings.get(objectId)
    if (bindings) {
      const index = bindings.findIndex(b => b.id === bindingId)
      if (index !== -1) {
        bindings.splice(index, 1)
      }
    }
  }

  /**
   * 注册对象
   */
  registerObject(objectId, threeObject) {
    this.objects.set(objectId, threeObject)
  }

  /**
   * 注销对象
   */
  unregisterObject(objectId) {
    this.objects.delete(objectId)
  }

  /**
   * 启用/禁用
   */
  setEnabled(enabled) {
    this.enabled = enabled
  }

  /**
   * 连接所有数据源
   */
  connectAll() {
    this.dataSources.forEach(source => {
      source.connect()
      source.on('data', (data) => this.handleData(source.id, data))
    })
  }

  /**
   * 断开所有连接
   */
  disconnectAll() {
    this.dataSources.forEach(source => {
      source.disconnect()
    })
  }

  /**
   * 处理数据
   */
  handleData(sourceId, data) {
    if (!this.enabled) return

    this.bindings.forEach((bindings, objectId) => {
      const object = this.objects.get(objectId)
      if (!object) return

      bindings.forEach(binding => {
        if (binding.dataSource !== sourceId) return
        if (!binding.enabled) return

        // 获取数据值
        let value = data
        if (binding.dataPath) {
          value = this.getNestedValue(data, binding.dataPath)
        }

        if (value === undefined) return

        // 转换值
        const transformFn = TransformFunctions[binding.transform?.type] || TransformFunctions.direct
        const transformedValue = transformFn(value, binding.transform)

        // 更新对象
        const updater = PropertyUpdaters[binding.property]
        if (updater) {
          updater(object, transformedValue)
        }
      })
    })
  }

  /**
   * 获取嵌套属性值
   */
  getNestedValue(obj, path) {
    return path.split('.').reduce((current, key) => current?.[key], obj)
  }

  /**
   * 更新绑定
   */
  updateBinding(objectId, bindingId, changes) {
    const bindings = this.bindings.get(objectId)
    if (!bindings) return

    const binding = bindings.find(b => b.id === bindingId)
    if (binding) {
      Object.assign(binding, changes)
    }
  }

  /**
   * 清空所有
   */
  clear() {
    this.disconnectAll()
    this.dataSources.clear()
    this.bindings.clear()
    this.objects.clear()
  }
}

/**
 * 数据源基类
 */
class DataSource {
  constructor(id, config) {
    this.id = id
    this.config = config
    this.connected = false
    this.listeners = []
  }

  connect() {
    this.connected = true
  }

  disconnect() {
    this.connected = false
  }

  emit(data) {
    this.listeners.forEach(cb => cb(data))
  }

  on(callback) {
    this.listeners.push(callback)
  }

  off(callback) {
    this.listeners = this.listeners.filter(cb => cb !== callback)
  }
}

/**
 * WebSocket 数据源
 */
class WebSocketDataSource extends DataSource {
  constructor(id, config) {
    super(id, config)
    this.socket = null
    this.reconnectAttempts = 0
    this.maxReconnectAttempts = 5
  }

  connect() {
    super.connect()

    try {
      this.socket = new WebSocket(this.config.url)

      this.socket.onopen = () => {
        this.reconnectAttempts = 0
        console.log(`WebSocket connected: ${this.id}`)
      }

      this.socket.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          this.emit(data)
        } catch (e) {
          this.emit(event.data)
        }
      }

      this.socket.onclose = () => {
        this.connected = false
        this.attemptReconnect()
      }

      this.socket.onerror = (error) => {
        console.error(`WebSocket error: ${this.id}`, error)
      }
    } catch (error) {
      console.error(`Failed to connect WebSocket: ${this.id}`, error)
    }
  }

  disconnect() {
    super.disconnect()
    if (this.socket) {
      this.socket.close()
      this.socket = null
    }
  }

  attemptReconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++
      console.log(`Reconnecting WebSocket: ${this.id} (attempt ${this.reconnectAttempts})`)
      setTimeout(() => this.connect(), 2000 * this.reconnectAttempts)
    }
  }
}

/**
 * HTTP 数据源
 */
class HTTPDataSource extends DataSource {
  constructor(id, config) {
    super(id, config)
    this.interval = null
  }

  connect() {
    super.connect()

    const fetchData = async () => {
      try {
        const response = await fetch(this.config.url)
        const data = await response.json()
        this.emit(data)
      } catch (error) {
        console.error(`HTTP fetch error: ${this.id}`, error)
      }
    }

    fetchData()

    if (this.config.interval && this.config.interval > 0) {
      this.interval = setInterval(fetchData, this.config.interval)
    }
  }

  disconnect() {
    super.disconnect()
    if (this.interval) {
      clearInterval(this.interval)
      this.interval = null
    }
  }
}

/**
 * 静态数据源
 */
class StaticDataSource extends DataSource {
  constructor(id, config) {
    super(id, config)
    this.data = config.data || {}
  }

  connect() {
    super.connect()
    this.emit(this.data)
  }
}

/**
 * MQTT 数据源（简化实现）
 */
class MQTTDataSource extends DataSource {
  constructor(id, config) {
    super(id, config)
    this.client = null
  }

  connect() {
    super.connect()
    // MQTT 需要额外的库支持，这里只是占位实现
    console.warn('MQTT support requires additional library (mqtt.js)')
  }

  disconnect() {
    super.disconnect()
    if (this.client) {
      this.client.end()
      this.client = null
    }
  }
}

export default DataBridge