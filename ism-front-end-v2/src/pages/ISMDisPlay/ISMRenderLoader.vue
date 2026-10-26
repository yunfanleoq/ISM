<template>
  <component
    :is="renderComp"
    v-if="renderComp"
    :showUuid="showUuid"
    showToken=""
    showDeviceUuid=""
  />
</template>

<script>
/**
 * 薄加载器：ISMRender 含 ISMBase + 200+ 子组件，必须与 pageView 完全隔离在独立 async chunk，
 * 否则 webpack 循环依赖会导致 pageView export default 为 undefined → setting 'render' 崩溃。
 */
export default {
  name: 'ISMRenderLoader',
  props: {
    showUuid: { type: String, required: true },
  },
  data() {
    return { renderComp: null }
  },
  created() {
    import(/* webpackChunkName: "ism-render" */ './ISMRender.vue')
      .then(mod => {
        const comp = mod && mod.default
        if (comp) {
          this.renderComp = comp
          return
        }
        console.error('[ISMRenderLoader] ISMRender export default 为空')
      })
      .catch(err => console.error('[ISMRenderLoader] 加载失败:', err))
  },
}
</script>
