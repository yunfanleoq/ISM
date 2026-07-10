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
    import(/* webpackChunkName: "ism-runtree-nav" */ './ISMRunTreeNav.vue')
      .then(mod => {
        const comp = mod && mod.default
        if (comp) {
          this.navComp = comp
          return
        }
        console.error('[RunTreeNavLoader] ISMRunTreeNav export default 为空')
      })
      .catch(err => console.error('[RunTreeNavLoader] 加载失败:', err))
  },
}
</script>
