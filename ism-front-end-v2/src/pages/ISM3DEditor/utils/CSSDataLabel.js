/**
 * CSSDataLabel — CSS2DRenderer 驱动的动态数据标签
 *
 * 功能：
 * - 基于 DIV + CSS2DRenderer 的 HTML 标签
 * - 支持动态数据绑定（实时刷新数值）
 * - 气泡样式、进度条、报警闪烁等预设样式
 *
 * 用法：
 *   const label = new CSSDataLabel()
 *   const el = label.create({
 *     text: '温度: 36°C',
 *     style: 'bubble',
 *     dataBind: { key: 'temp', value: 36, unit: '°C' }
 *   })
 *   scene.add(el)
 *
 *   // 动态更新
 *   label.update(el, { value: 42 })
 */
import * as THREE from 'three'
import { CSS2DObject } from 'three/examples/jsm/renderers/CSS2DRenderer.js'

// ===== 预设样式 =====
const PRESET_STYLES = {
  /** 气泡标签 */
  bubble: {
    container: 'padding:4px 12px;background:rgba(0,20,50,0.85);border:1px solid #00aaff;border-radius:16px;color:#e0f0ff;font-size:13px;font-family:"Microsoft YaHei",sans-serif;white-space:nowrap;pointer-events:none;backdrop-filter:blur(4px);box-shadow:0 0 8px rgba(0,170,255,0.3)',
    value: 'color:#00ffff;font-weight:bold;font-size:16px'
  },
  /** 数据面板 */
  panel: {
    container: 'padding:8px 14px;background:rgba(0,10,30,0.9);border:1px solid #2266aa;border-radius:6px;color:#aaccff;font-size:12px;font-family:"Microsoft YaHei",sans-serif;line-height:1.6;pointer-events:none;min-width:100px',
    title: 'color:#88bbff;font-size:11px;margin-bottom:2px',
    value: 'color:#00ff88;font-weight:bold;font-size:18px'
  },
  /** 告警标签 */
  alert: {
    container: 'padding:4px 12px;background:rgba(50,0,0,0.9);border:1px solid #ff4444;border-radius:4px;color:#ffaaaa;font-size:13px;font-family:"Microsoft YaHei",sans-serif;white-space:nowrap;pointer-events:none;animation:cssLabelBlink 1s infinite',
    value: 'color:#ff4444;font-weight:bold'
  },
  /** 迷你仪表 */
  gauge: {
    container: 'padding:6px 10px;background:rgba(0,20,40,0.85);border:1px solid #2266aa;border-radius:4px;color:#aaccee;font-size:12px;font-family:"Microsoft YaHei",sans-serif;pointer-events:none;min-width:80px;text-align:center',
    value: 'color:#ffcc00;font-weight:bold;font-size:20px',
    bar: 'width:100%%;height:3px;background:rgba(0,100,200,0.3);border-radius:2px;margin-top:3px;overflow:hidden'
  }
}

// 注入闪烁动画
let _styleInjected = false
function _injectStyles() {
  if (_styleInjected) return
  const style = document.createElement('style')
  style.textContent = `
    @keyframes cssLabelBlink {
      0%, 100% { opacity: 1; }
      50% { opacity: 0.4; }
    }
  `
  document.head.appendChild(style)
  _styleInjected = true
}

export class CSSDataLabel {
  constructor() {
    _injectStyles()
    this._labels = [] // 跟踪所有创建的标签
  }

  /**
   * 创建数据标签
   * @param {Object} options
   * @param {string} options.text - 显示文本（支持 {{value}} 占位符）
   * @param {string} [options.title] - 标题
   * @param {string} [options.style='bubble'] - 预设样式名
   * @param {Object} [options.dataBind] - 数据绑定 { key, value, unit, max }
   * @param {THREE.Vector3} [options.position] - 世界坐标位置
   * @returns {CSS2DObject}
   */
  create(options = {}) {
    const style = PRESET_STYLES[options.style] || PRESET_STYLES.bubble
    const data = options.dataBind || {}

    const div = document.createElement('div')

    // 构建 innerHTML
    let html = ''
    if (options.title && (options.style === 'panel')) {
      html += `<div style="${style.title || ''}">${options.title}</div>`
    }
    html += `<span style="${style.value || ''}">${data.value !== undefined ? data.value : ''}</span>`
    if (data.unit) {
      html += `<span style="font-size:12px;color:#aaccff"> ${data.unit}</span>`
    }

    // 仪表进度条
    if (options.style === 'gauge' && data.max) {
      const pct = Math.min(100, Math.max(0, (data.value / data.max) * 100))
      html += `<div style="${style.bar}"><div style="width:${pct}%;height:100%;background:linear-gradient(90deg,#00ff88,#ffcc00);border-radius:2px"></div></div>`
    }

    div.innerHTML = html
    div.setAttribute('style', style.container)

    // data 属性（用于后续更新）
    div.dataset.labelStyle = options.style || 'bubble'
    div.dataset.labelTitle = options.title || ''
    div.dataset.labelUnit = data.unit || ''
    div.dataset.labelMax = data.max || ''

    const obj = new CSS2DObject(div)
    if (options.position) {
      obj.position.copy(options.position)
    }
    obj.userData.isCSSLabel = true
    obj.userData.labelConfig = {
      style: options.style,
      title: options.title,
      dataBind: data
    }

    this._labels.push(obj)
    return obj
  }

  /**
   * 更新标签数据显示
   * @param {CSS2DObject} labelObj
   * @param {Object} data - { value, unit, max }
   */
  update(labelObj, data = {}) {
    if (!labelObj || !labelObj.userData.isCSSLabel) return
    const div = labelObj.element
    if (!div) return

    const style = div.dataset.labelStyle || 'bubble'
    const title = div.dataset.labelTitle || ''
    const unit = data.unit !== undefined ? data.unit : div.dataset.labelUnit
    const max = data.max !== undefined ? data.max : div.dataset.labelMax

    const preset = PRESET_STYLES[style] || PRESET_STYLES.bubble

    let html = ''
    if (title && style === 'panel') {
      html += `<div style="${preset.title || ''}">${title}</div>`
    }
    html += `<span style="${preset.value || ''}">${data.value !== undefined ? data.value : ''}</span>`
    if (unit) {
      html += `<span style="font-size:12px;color:#aaccff"> ${unit}</span>`
    }

    if (style === 'gauge' && max) {
      const pct = Math.min(100, Math.max(0, (data.value / max) * 100))
      html += `<div style="${preset.bar}"><div style="width:${pct}%;height:100%;background:linear-gradient(90deg,#00ff88,#ffcc00);border-radius:2px"></div></div>`
    }

    div.innerHTML = html

    // 更新 dataset
    if (data.unit !== undefined) div.dataset.labelUnit = data.unit
    if (data.max !== undefined) div.dataset.labelMax = data.max
  }

  /**
   * 批量更新
   * @param {Array<{label: CSS2DObject, data: Object}>} updates
   */
  updateBatch(updates) {
    for (const u of updates) {
      this.update(u.label, u.data)
    }
  }

  /**
   * 销毁标签
   */
  destroy(labelObj) {
    if (!labelObj) return
    const idx = this._labels.indexOf(labelObj)
    if (idx > -1) this._labels.splice(idx, 1)
    if (labelObj.element && labelObj.element.parentNode) {
      labelObj.element.parentNode.removeChild(labelObj.element)
    }
  }

  /**
   * 销毁所有标签
   */
  destroyAll() {
    for (const label of [...this._labels]) {
      this.destroy(label)
    }
    this._labels = []
  }
}

export default CSSDataLabel
