<template>
  <div class="ism-pageview">
    <SCADAMonitor
      v-if="scadaConfig"
      :title="scadaConfig.title"
      :subtitle="scadaConfig.subtitle"
      :project-uuid="scadaConfig.projectUuid"
    />
    <template v-else>
      <ISMRender :showUuid="$route.params.uid" showToken="" showDeviceUuid="" />
      <PreviewWatermark />
    </template>
  </div>
</template>

<script>
import ISMRender from './ISMRender'
import PreviewWatermark from '@/components/PreviewWatermark.vue'
import SCADAMonitor from '@/pages/SCADAMonitor/index'
import { getScadaConfigForDisplay } from '@/pages/SCADAMonitor/scadaDisplayRegistry'

export default {
  name: 'AppRun',
  components: {
    ISMRender,
    PreviewWatermark,
    SCADAMonitor,
  },
  computed: {
    scadaConfig() {
      return getScadaConfigForDisplay(this.$route.params.uid)
    },
  },
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
