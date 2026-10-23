<template>
  <viz-display-host :page-name="pageName" />
</template>

<script>
import VizDisplayHost from '@/xunzhiJf/pages/VizDisplayHost.vue'
import { XUNZHI_JF_ENV_PLANE_PAGE } from '@/xunzhiJf/resolveDisplayPage'

export default {
  name: 'XunzhiJfVizPlanePage',
  components: { VizDisplayHost },
  data() {
    return { pageName: XUNZHI_JF_ENV_PLANE_PAGE }
  },
}
</script>
