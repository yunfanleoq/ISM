import Vue from 'vue'
import BoundPointsPopup from '../BoundPointsPopup.vue'

let instance = null

/**
 * 全页只挂一张测点卡片到 document.body。
 * 挂在 X6 节点里时，position:fixed 会被画布变换吃掉，卡片会跟着鼠标跑，也收不到点击。
 */
export function getBoundPointsPopup(vm, create = true) {
  if (instance && !instance._isDestroyed) {
    return instance
  }
  if (!create || typeof document === 'undefined' || !vm) {
    return null
  }
  const Ctor = Vue.extend(BoundPointsPopup)
  instance = new Ctor({ parent: vm.$root || vm })
  instance.$mount()
  document.body.appendChild(instance.$el)
  return instance
}
