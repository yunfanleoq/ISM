<template>
  <viz-display-host :page-name="pageName" />
</template>

<script>
import VizDisplayHost from '@/xunzhiJf/pages/VizDisplayHost.vue'
import { XUNZHI_JF_NAV_25D_PAGE } from '@/xunzhiJf/resolveDisplayPage'

export default {
  name: 'XunzhiJfVizNav25dPage',
  components: { VizDisplayHost },
  data() {
    return { pageName: XUNZHI_JF_NAV_25D_PAGE }
  },
}
</script>
