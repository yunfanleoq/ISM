<template>
  <component
    :is="navComp"
    v-if="navComp"
    :model-id="modelId"
    :project-uuid="projectUuid"
  />
</template>

<script>
export default {
  name: 'RunTreeNavLoader',
  props: {
    modelId: { type: String, required: true },
    projectUuid: { type: String, required: true },
  },
  data() {
    return { navComp: null }
  },
  created() {
    // 使用版本化 chunk 名，避开 writeToDisk 场景下同名旧 chunk 被浏览器复用。
    // 曾出现 import 返回 Vuex store（而不是导航 SFC）并被 Vue 当组件挂载，
    // 导致 mounted/fetchTree 均不执行，导航与画布同时空白。
    import(/* webpackChunkName: "ism-runtree-nav-v2" */ './ISMRunTreeNav.vue')
      .then(mod => {
        const comp = mod && mod.default
        if (comp && typeof comp.data === 'function' && comp.methods && comp.methods.fetchTree) {
          this.navComp = comp
          return
        }
        console.error('[RunTreeNavLoader] ISMRunTreeNav 导出无效，未挂载错误对象', {
          exportKeys: comp ? Object.keys(comp) : [],
        })
      })
      .catch(err => console.error('[RunTreeNavLoader] 加载失败:', err))
  },
}
</script>
