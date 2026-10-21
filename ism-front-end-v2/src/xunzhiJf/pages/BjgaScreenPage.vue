<template>
  <div class="bjga-host" :class="{ 'is-full': fullscreen }">
    <iframe class="bjga-frame" :src="src" title="公安局动环" />
  </div>
</template>

<script>
export default {
  name: 'BjgaScreenPage',
  props: {
    fullscreen: { type: Boolean, default: false },
  },
  computed: {
    src() {
      const hash = (this.$route.meta && this.$route.meta.bjgaHash) || this.$route.query.page || 'dashboard'
      return '/bjga/index.html#/' + String(hash).replace(/^#?\//, '')
    },
  },
}
</script>

<style scoped>
.bjga-host {
  height: calc(100vh - 112px);
  min-height: 640px;
  background: #050d24;
}
.bjga-host.is-full {
  height: 100vh;
  min-height: 0;
}
.bjga-frame {
  width: 100%;
  height: 100%;
  border: 0;
  background: #050d24;
}
</style>
