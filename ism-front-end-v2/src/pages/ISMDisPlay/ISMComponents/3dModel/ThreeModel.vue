<template>
  <div :style="containerStyle" class="three-model-container">
    <div v-if="isEditing" class="three-model-placeholder">
      <a-empty :description="$t('3D模型组件')" />
      <div class="model-info" v-if="componentData && resolvedModelUrl">
        <p>{{ $t('模型') }}: {{ getFileName(resolvedModelUrl) }}</p>
        <p>{{ $t('尺寸') }}: {{ componentData.width }}x{{ componentData.height }}</p>
      </div>
    </div>
    <vue3dLoader
      v-if="!isEditing && showRender"
      ref="threeLoader"
      :filePath="resolvedModelUrl"
      :width="componentData.width || 300"
      :height="componentData.height || 200"
      :position="componentData.position || { x: 0, y: 0, z: 0 }"
      :rotation="componentData.rotation || { x: 0, y: 0, z: 0 }"
      :scale="componentData.scale || { x: 1, y: 1, z: 1 }"
      :cameraPosition="componentData.cameraPosition || { x: 0, y: 0, z: 5 }"
      :backgroundColor="componentData.backgroundColor || 0x000000"
      :backgroundAlpha="componentData.backgroundAlpha !== undefined ? componentData.backgroundAlpha : 0"
      :lights="componentData.lights || defaultLights"
      :autoPlay="componentData.autoPlay !== undefined ? componentData.autoPlay : true"
      :showFps="false"
      :enableDraco="componentData.enableDraco || false"
      @load="onModelLoad"
      @error="onModelError"
    />
  </div>
</template>

<script>
import vue3dLoader from '@/components/3dLoader/vue3dLoader.vue'

export default {
  name: 'ThreeModel',
  components: {
    vue3dLoader
  },
  props: {
    nodeData: {
      type: Object,
      default: () => ({})
    },
    isEditing: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      showRender: true,
      _watchers: [],
      defaultLights: [
        {
          type: 'AmbientLight',
          color: 0x404040
        },
        {
          type: 'DirectionalLight',
          position: { x: 1, y: 1, z: 1 },
          color: 0xffffff,
          intensity: 1
        }
      ]
    }
  },
  computed: {
    componentData() {
      return (this.nodeData && this.nodeData.data && this.nodeData.data.detail) || {}
    },
    resolvedModelUrl() {
      return this.componentData.modelPath || this.componentData.modelUrl || ''
    },
    containerStyle() {
      const detail = this.componentData || {}
      return {
        width: (detail.width || 300) + 'px',
        height: (detail.height || 200) + 'px',
        overflow: 'hidden'
      }
    }
  },
  watch: {
    resolvedModelUrl() {
      this.reloadModel()
    },
    'componentData.width'() {
      this.reloadModel()
    },
    'componentData.height'() {
      this.reloadModel()
    }
  },
  methods: {
    getFileName(url) {
      if (!url) return ''
      const parts = url.split('/')
      return parts[parts.length - 1]
    },
    onModelLoad(object) {
      this.$emit('load', object)
      // 如果配置了自动旋转
      if (this.componentData.autoRotate) {
        this.startAutoRotate()
      }
      // 数据绑定：应用绑定的变换
      this.applyDataBinding()
    },
    onModelError(error) {
      console.error('3D Model Load Error:', error)
      this.$message && this.$message.error(this.$t('模型加载失败'))
    },
    reloadModel() {
      this.showRender = false
      this.$nextTick(() => {
        this.showRender = true
      })
    },
    startAutoRotate() {
      const loader = this.$refs.threeLoader
      if (!loader) return
      const object = loader.object
      if (!object) return

      const animate = () => {
        if (!this.componentData.autoRotate) return
        object.rotation.y += 0.01
        this._rotateFrameId = requestAnimationFrame(animate)
      }
      animate()
    },
    applyDataBinding() {
      if (!this.componentData.bindData || !this.componentData.bindData.enabled) return
      const bindData = this.componentData.bindData
      const loader = this.$refs.threeLoader
      if (!loader || !loader.object) return

      // 先清理旧的 watcher（防止重复调用时泄漏）
      this._cleanupWatchers()

      // 绑定旋转
      if (bindData.rotationBind) {
        // 从数据源获取值并更新旋转
        this._watchers.push(this.$watch(() => this.getBoundValue(bindData.rotationBind), (val) => {
          if (loader.object) {
            loader.object.rotation.y = (val || 0) * Math.PI / 180
          }
        }))
      }

      // 绑定显隐
      if (bindData.visibleBind) {
        this._watchers.push(this.$watch(() => this.getBoundValue(bindData.visibleBind), (val) => {
          if (loader.object) {
            loader.object.visible = !!val
          }
        }))
      }
    },
    _cleanupWatchers() {
      if (!this._watchers) return
      this._watchers.forEach(unwatch => { try { unwatch() } catch(e) {} })
      this._watchers = []
    },
    getBoundValue(bindKey) {
      // 从 Vuex store 或数据源获取绑定值
      // 格式：device001.temperature
      if (!bindKey) return null
      const store = this.$store && this.$store.state.ISMDisPlayEditorTool
      if (store && store.realTimeData) {
        return store.realTimeData[bindKey]
      }
      return null
    },
    // 更新组件数据（供属性面板调用）
    updateComponentData(key, value) {
      if (!this.nodeData || !this.nodeData.data) return
      if (!this.nodeData.data.detail) {
        this.$set(this.nodeData.data, 'detail', {})
      }
      this.$set(this.nodeData.data.detail, key, value)
    }
  },
  beforeDestroy() {
    this.showRender = false
    // 清理动态 $watch
    this._cleanupWatchers()
    // 取消动画帧
    if (this._rotateFrameId) {
      cancelAnimationFrame(this._rotateFrameId)
      this._rotateFrameId = null
    }
  }
}
</script>

<style scoped>
.three-model-container {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}
.three-model-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.05);
  border: 1px dashed #ccc;
}
.model-info {
  margin-top: 8px;
  font-size: 12px;
  color: #666;
}
.model-info p {
  margin: 2px 0;
}
</style>
