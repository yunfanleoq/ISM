/**
 * ISM 子组件内存泄漏自动修复 Mixin（零侵入版）
 *
 * 解决问题：ISMComponents 下 ~200 个子组件的 EventBus $on / X6 Node 事件未清理
 *
 * 原理：
 *   1. created() 中用 Object.defineProperty 将 this.$EventBus 替换为追踪代理
 *      所有 this.$EventBus.$on / _t.$EventBus.$on 调用自动被记录
 *   2. beforeDestroy() 中用记录的精确 handler 引用进行 $off，不影响其他组件
 *   3. X6 Node 事件（change:data / change:size）用广播 off（安全，每个组件独立节点）
 *
 * 使用方法（组件只需加一行）：
 *   import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
 *   export default { mixins: [ISMChildAutoMixin], ... }
 *
 * 注意：组件中已有的 this.$EventBus.$on / _t.$EventBus.$on 代码不需要任何修改
 */

export default {
  created() {
    // 仅对有 $EventBus 的 ISM 子组件生效
    if (!this.$EventBus) return

    const self = this
    const realBus = this.$EventBus
    const tracked = []

    // 保存追踪列表，beforeDestroy 时使用
    this._ismAutoTrackedHandlers = tracked

    // 创建代理对象，拦截 $on / $once 调用以自动追踪 handler 引用
    const proxy = {
      $on(event, handler) {
        tracked.push({ event, handler })
        return realBus.$on(event, handler)
      },
      $off(event, handler) {
        return realBus.$off(event, handler)
      },
      $once(event, handler) {
        tracked.push({ event, handler })
        return realBus.$once(event, handler)
      },
      $emit(...args) {
        return realBus.$emit(...args)
      },
    }

    // 用代理替换当前组件实例的 $EventBus（不影响其他组件，因为每个 Vue 实例独立）
    Object.defineProperty(this, '$EventBus', {
      value: proxy,
      writable: false,
      configurable: true,
    })
  },

  beforeDestroy() {
    // 1. 精确清理 EventBus（只移除当前组件的监听，不影响其他组件）
    const tracked = this._ismAutoTrackedHandlers
    if (tracked && tracked.length > 0) {
      // beforeDestroy 中 $EventBus 可能已被其他逻辑修改，用 originalBus
      const bus = Object.getOwnPropertyDescriptor(this, '$EventBus')
      const realBus = bus && bus.value && bus.value.$off ? bus.value : (this.$EventBus || null)
      for (let i = 0; i < tracked.length; i++) {
        try {
          if (realBus) realBus.$off(tracked[i].event, tracked[i].handler)
        } catch (e) {
          // ignore
        }
      }
      this._ismAutoTrackedHandlers = null
    }

    // 2. X6 Node 事件清理（广播 off 在节点上安全 — 每个组件有独立的 GetNodeObj）
    if (this.GetNodeObj) {
      try { this.GetNodeObj.off('change:data') } catch (e) {}
      try { this.GetNodeObj.off('change:size') } catch (e) {}
    }
  },
}
