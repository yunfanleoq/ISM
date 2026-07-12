<template>
  <div class="ism-pageview">
    <ISMRenderLoader :show-uuid="dashboardUuid" />
    <ScadaOrgOverview
      v-if="showRunTree"
      :model-id="dashboardUuid"
      :project-uuid="projectUuid"
    />
    <RunTreeNavLoader
      v-if="showRunTree"
      :model-id="dashboardUuid"
      :project-uuid="projectUuid"
    />
    <ScadaAlarmPanel
      v-if="showRunTree"
      :model-id="dashboardUuid"
      :project-uuid="projectUuid"
    />
    <BackToAdminButton :project-uuid="projectUuid" />
    <PreviewWatermark />
  </div>
</template>

<script>
import PreviewWatermark from '@/components/PreviewWatermark.vue'
import BackToAdminButton from '@/components/BackToAdminButton.vue'

const ISMRenderLoader = () => import(/* webpackChunkName: "ism-render-loader" */ './ISMRenderLoader.vue')
const RunTreeNavLoader = () => import(/* webpackChunkName: "ism-runtree-loader" */ './RunTreeNavLoader.vue')
const ScadaAlarmPanel = () => import(/* webpackChunkName: "scada-alarm-panel" */ './ScadaAlarmPanel.vue')
const ScadaOrgOverview = () => import(/* webpackChunkName: "scada-org-overview" */ './ScadaOrgOverview.vue')

import {
  applyHomeProjectAuth,
  HOME_PROJECT_UUID,
  resolveHomeProjectUuid,
  shouldShowRunTreeNav,
} from '@/config/homeDashboard.js'
import { AUTH_TYPE, getAuthorization } from '@/utils/request'

export default {
  name: 'AppRunShell',
  components: {
    ISMRenderLoader,
    ScadaOrgOverview,
    PreviewWatermark,
    RunTreeNavLoader,
    ScadaAlarmPanel,
    BackToAdminButton,
  },
  props: {
    uid: { type: String, required: true },
  },
  created() {
    if (shouldShowRunTreeNav(this.uid, this.$store)) {
      applyHomeProjectAuth(this.$store)
    }
  },
  computed: {
    dashboardUuid() {
      return this.uid
    },
    showRunTree() {
      return shouldShowRunTreeNav(this.dashboardUuid, this.$store)
    },
    projectUuid() {
      const uid = this.dashboardUuid
      if (shouldShowRunTreeNav(uid, this.$store)) {
        return resolveHomeProjectUuid(uid, this.$store) || HOME_PROJECT_UUID
      }
      const resolved = resolveHomeProjectUuid(uid, this.$store)
      if (resolved) return resolved
      return getAuthorization(AUTH_TYPE.AUTH1) || ''
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

.ism-pageview::before {
  content: "";
  position: absolute;
  inset: 0;
  z-index: 45;
  pointer-events: none;
  mix-blend-mode: screen;
  background:
    radial-gradient(38% 42% at 22% 28%, rgba(34, 211, 238, 0.14), transparent 60%),
    radial-gradient(34% 40% at 80% 30%, rgba(59, 130, 246, 0.12), transparent 62%),
    radial-gradient(48% 52% at 60% 90%, rgba(20, 170, 210, 0.10), transparent 64%);
  background-size: 200% 200%, 200% 200%, 200% 200%;
  animation: ismHaloDrift 18s ease-in-out infinite alternate;
}
@keyframes ismHaloDrift {
  0%   { background-position:   0% 0%, 100%  0%, 50% 100%; }
  50%  { background-position:  28% 38%, 70% 28%, 42% 72%; }
  100% { background-position:  10% 18%, 92% 10%, 62% 92%; }
}

.ism-pageview::after {
  content: "";
  position: absolute;
  inset: 0;
  z-index: 46;
  pointer-events: none;
  mix-blend-mode: screen;
  background: linear-gradient(115deg, transparent 44%, rgba(130, 215, 255, 0.10) 50%, transparent 56%);
  background-size: 250% 250%;
  animation: ismScanSweep 14s linear infinite;
}
@keyframes ismScanSweep {
  0%   { background-position: -60% -60%; opacity: 0; }
  14%  { opacity: 0.3; }
  84%  { opacity: 0.3; }
  100% { background-position: 160% 160%; opacity: 0; }
}
</style>
