<template>
  <div class="ism-pageview">
    <ISMRender :showUuid="$route.params.uid" showToken="" showDeviceUuid="" />
    <DeviceHoverTooltip />
    <!-- 航信机房 NCC 大屏：可折叠层级导航树（覆盖在静态侧栏区域之上） -->
    <ISMRunTreeNav v-if="showRunTree" />
    <PreviewWatermark />
  </div>
</template>

<script>

import ISMRender from './ISMRender';
import PreviewWatermark from '@/components/PreviewWatermark.vue'
import DeviceHoverTooltip from './DeviceHoverTooltip.vue'
import ISMRunTreeNav from './ISMRunTreeNav.vue'

// 仅对航信机房 NCC 大屏启用折叠导航树
const NCC_DASHBOARD_UUID = '043135ad-44be-e5d8-89be-3e54883c23a8'

export default {
  name: 'AppRun',
  components: {
    ISMRender,
    PreviewWatermark,
    DeviceHoverTooltip,
    ISMRunTreeNav
  },
  computed: {
    showRunTree() {
      return this.$route.params.uid === NCC_DASHBOARD_UUID
    }
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


