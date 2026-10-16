<template>
  <div class="xunzhi-jf-skeleton">
    <a-alert
      type="warning"
      show-icon
      message="演示骨架页 · 非正式验收"
      :description="alertDesc"
    />
    <a-card :title="title" style="margin-top: 16px">
      <p><strong>菜单码</strong>：{{ menuCode }}</p>
      <p><strong>子模块</strong>：{{ moduleId }}</p>
      <p><strong>权限</strong>：{{ perm }}</p>
      <p v-if="routePath"><strong>路由</strong>：{{ routePath }}</p>
      <p class="hint">功能未齐时仍保留侧栏入口；正式页按契约切片替换本占位。</p>
    </a-card>
  </div>
</template>

<script>
export default {
  name: 'XunzhiJfSkeletonPage',
  computed: {
    metaProps() {
      return (this.$route.meta && this.$route.meta.skeletonProps) || {}
    },
    title() {
      return this.metaProps.title || '占位'
    },
    menuCode() {
      return this.metaProps.menuCode || ''
    },
    moduleId() {
      return this.metaProps.moduleId || ''
    },
    perm() {
      return this.metaProps.perm || ''
    },
    routePath() {
      return this.metaProps.routePath || ''
    },
    sliceHint() {
      return this.metaProps.sliceHint || ''
    },
    alertDesc() {
      const parts = ['建设中 / Mock / 依赖未齐']
      if (this.sliceHint) parts.push(this.sliceHint)
      return parts.join(' · ')
    },
  },
}
</script>

<style scoped>
.xunzhi-jf-skeleton {
  padding: 16px;
}
.hint {
  color: rgba(0, 0, 0, 0.45);
  margin-top: 12px;
}
</style>
