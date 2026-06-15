/**
 * ISM X6 Cell 通用清理 mixin
 *
 * 解决问题：X6 VueShape 创建的 Vue 组件实例在 cell.remove → view.unmount → $destroy 时
 * 会触发 beforeDestroy，但子组件注册的 EventBus 监听器（cell-editMode / activeEvent / animateEvent）
 * 没有被清理，导致闭包持有 Vue 实例引用，无法被 GC。
 *
 * 使用方法：在子组件中混入即可
 *   import { ismCellCleanup } from '@/mixins/ismCellCleanup'
 *   export default { mixins: [ismCellCleanup], ... }
 */

function getIdentifier() {
  try {
    return this.detail && this.detail.identifier
  } catch (e) {
    return null
  }
}

export const ismCellCleanup = {
  data() {
    return {
      _busHandlers: null,
      _nodeObj: null,
    }
  },
  beforeDestroy() {
    // 清理 X6 Node 事件
    if (this._nodeObj) {
      this._nodeObj.off('change:data')
      this._nodeObj.off('change:size')
      this._nodeObj = null
    }

    // 清理 EventBus 监听
    const h = this._busHandlers || {}
    const id = getIdentifier.call(this)
    if (id) {
      this.$EventBus && this.$EventBus.$off(id + 'activeEvent', h.activeEvent)
      this.$EventBus && this.$EventBus.$off(id + 'animateEvent', h.animateEvent)
    }
    this.$EventBus && this.$EventBus.$off('cell-editMode', h.cellEditMode)
    this._busHandlers = null
  },
  methods: {
    /**
     * 在 created 中替代直接 this.getNode() + this.$EventBus.$on
     * 自动保存引用以便 beforeDestroy 清理
     */
    initISMCellListeners(options) {
      const node = this.getNode()
      this._nodeObj = node
      if (!this._busHandlers) {
        this._busHandlers = {}
      }
      const h = this._busHandlers

      // X6 Node 事件（默认注册 change:data 和 change:size）
      node.on('change:data', ({ current }) => {
        if (current) {
          this.detail = current.detail
        }
      })
      node.on('change:size', ({ current }) => {
        if (this.detail && this.detail.style) {
          this.detail.style.position.w = current.width
          this.detail.style.position.h = current.height
        }
      })

      // EventBus cell-editMode 监听
      h.cellEditMode = options.onCellEditMode || ((data) => {
        this.editMode = data.edit
        this.IsToolBox = data.toolbox
        if (this.initComponents) {
          this.initComponents(this.detail)
        }
      })
      this.$EventBus.$on('cell-editMode', h.cellEditMode)

      // 初始化 detail
      this.detail = node.getData().detail
      this.editMode = node.getData().editMode
      this.showDeviceUuid = node.getData().showDeviceUuid
      this.IsToolBox = node.getData().IsToolBox
    },

    /**
     * 在 mounted 中调用，注册 activeEvent / animateEvent
     */
    initISMCellActiveListeners(onActiveEvent, onAnimateEvent) {
      const id = getIdentifier.call(this)
      if (!id || !this._busHandlers) return
      const h = this._busHandlers

      const activeEvent = id + 'activeEvent'
      const animateEvent = id + 'animateEvent'
      if (h.activeEvent) {
        this.$EventBus.$off(activeEvent, h.activeEvent)
      }
      if (h.animateEvent) {
        this.$EventBus.$off(animateEvent, h.animateEvent)
      }

      if (onActiveEvent) {
        h.activeEvent = onActiveEvent
        this.$EventBus.$on(activeEvent, h.activeEvent)
      }
      if (onAnimateEvent) {
        h.animateEvent = onAnimateEvent
        this.$EventBus.$on(animateEvent, h.animateEvent)
      }
    },
  },
}
