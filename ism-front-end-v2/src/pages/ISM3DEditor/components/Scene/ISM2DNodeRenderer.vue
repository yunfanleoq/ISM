<template>
  <div class="ism-2d-node-renderer">
    <div
      class="ism-render-wrapper"
      :class="{ 'ism-render-wrapper-clickable': detail && detail.action && detail.action.length > 0 }"
      :style="wrapperStyle"
    >
      <component
        v-if="resolvedComponent"
        :is="resolvedComponent"
        :detail="detail"
        :editMode="editMode"
        :selected="selected"
        :showDeviceUuid="showDeviceUuid"
        :IsToolBox="false"
      />
      <div v-else class="ism-2d-node-fallback">{{ detail && detail.type ? detail.type : '2D Component' }}</div>
    </div>
  </div>
</template>

<script>
import { getISMComponent } from '@/pages/ISMDisPlay/componentRegistry'
import { mapActions, mapMutations, mapState } from 'vuex'

export default {
  name: 'ISM2DNodeRenderer',
  props: {
    objectData: {
      type: Object,
      required: true
    },
    editMode: {
      type: Boolean,
      default: false
    },
    selected: {
      type: Boolean,
      default: false
    },
    showDeviceUuid: {
      type: String,
      default: ''
    }
  },
  provide() {
    return {
      getNode: () => this.nodeApi
    }
  },
  computed: {
    detail() {
      let source2D = this.objectData && this.objectData.source2D ? this.objectData.source2D : {}
      // 确保 identifier 存在（3D场景中可能缺失）
      if (!source2D.identifier) {
        source2D.identifier = 'scene_2d_' + (this.objectData && this.objectData.id ? this.objectData.id : Date.now())
      }
      // 确保 active 数组存在
      if (!Array.isArray(source2D.active)) {
        source2D.active = []
      }
      // 确保 animate 对象存在
      if (!source2D.animate || typeof source2D.animate !== 'object') {
        source2D.animate = {
          selected: [],
          condition: {},
          isExpression: false,
          animateList: [],
          animateElement: []
        }
      }
      // 确保 style 对象存在
      if (!source2D.style || typeof source2D.style !== 'object') {
        source2D.style = {
          position: { x: 0, y: 0, w: 200, h: 200 },
          foreColor: '#333333',
          fontFamily: 'Arial',
          fontSize: 14,
          borderWidth: 0,
          diy: []
        }
      }
      return source2D
    },
    wrapperStyle() {
      const style = this.detail && this.detail.style ? this.detail.style : {}
      const position = style.position || {}
      const zIndex = parseInt(style.zIndex, 10)
      return {
        width: (position.w || 0) + 'px',
        height: (position.h || 0) + 'px',
        backgroundColor: style.backColor || 'transparent',
        zIndex: isNaN(zIndex) ? 1 : zIndex,
        borderWidth: (style.borderWidth || 0) + 'px',
        borderStyle: style.borderStyle || 'solid',
        borderColor: style.borderColor || '#000000'
      }
    },
    resolvedComponent() {
      const type = this.detail.type
      console.log('Resolving component for type:', type)
      const component = getISMComponent(type)
      if (!component) {
        console.warn('Component not found for type:', type)
      }
      return component
    },
    nodeData() {
      var position = this.detail && this.detail.style && this.detail.style.position ? this.detail.style.position : {}
      return {
        detail: this.detail,
        editMode: this.editMode,
        showDeviceUuid: this.showDeviceUuid || '',
        IsToolBox: false,
        width: position.w || 0,
        height: position.h || 0
      }
    }
  },
  watch: {
    detail: {
      deep: true,
      handler(newVal) {
        if (!this.nodeApi || !newVal) return
        this.nodeApi.emit('change:data', { current: this.nodeData })
        this.nodeApi.emit('change:size', { current: this.nodeData })
      }
    },
    editMode() {
      if (!this.nodeApi) return
      this.$EventBus.$emit('cell-editMode', { edit: this.editMode, toolbox: false })
      this.nodeApi.emit('change:data', { current: this.nodeData })
    },
    showDeviceUuid() {
      if (!this.nodeApi) return
      this.nodeApi.emit('change:data', { current: this.nodeData })
    }
  },
  created() {
    const listeners = {}
    this.nodeApi = {
      getData: () => this.nodeData,
      on: (event, handler) => {
        if (!listeners[event]) listeners[event] = []
        listeners[event].push(handler)
      },
      off: (event, handler) => {
        if (!listeners[event]) return
        if (!handler) {
          listeners[event] = []
          return
        }
        listeners[event] = listeners[event].filter((fn) => fn !== handler)
      },
      emit: (event, payload) => {
        (listeners[event] || []).forEach((fn) => fn(payload))
      }
    }
    this.$EventBus.$emit('cell-editMode', { edit: this.editMode, toolbox: false })
    this.$EventBus.$emit('cell-vuex', {
      PMapState: mapState,
      PMapActions: mapActions,
      PMapMutations: mapMutations,
      PStore: this.$store,
      PRouter: this.$router,
      loginUuid: this.$route && this.$route.params ? this.$route.params.uid : ''
    })
  }
}
</script>

<style scoped>
.ism-2d-node-renderer {
  width: 100%;
  height: 100%;
  position: relative;
}

.ism-render-wrapper {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

.ism-render-wrapper-clickable {
  cursor: pointer;
}

.ism-2d-node-fallback {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px dashed #13c2c2;
  background: rgba(255, 255, 255, 0.92);
  color: #13c2c2;
  font-size: 12px;
  text-align: center;
  overflow: hidden;
}
</style>
