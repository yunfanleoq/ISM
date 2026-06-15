<template>
  <div class="ism-pageview">
    <ISMRender :showUuid="$route.params.uid" showToken="" showDeviceUuid="" />
    <PreviewWatermark />
  </div>
</template>

<script>

import ISMRender from './ISMRender';
import PreviewWatermark from '@/components/PreviewWatermark.vue'

export default {
  name: 'AppRun',
  components: {
    ISMRender,
    PreviewWatermark
  },
  methods: {

  },
  created() {

  },
  mounted() {

  }
}
</script>

<style lang="less">
.ism-pageview {
  height: 100vh;
  width: 100%;
  position: relative;
  overflow: hidden;
}
</style>


