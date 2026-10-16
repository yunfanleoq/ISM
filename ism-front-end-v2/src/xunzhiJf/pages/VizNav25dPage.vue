<template>
  <div class="xunzhi-nav25d-page">
    <a-alert
      type="info"
      show-icon
      banner
      message="试验切片 · 非正式验收"
      description="本页为 Three.js 2.5D 观感演示（假机柜阵列），非 BIM/现场模型。正式监视请使用平面组态兜底。"
    />
    <div class="xunzhi-nav25d-toolbar">
      <a-radio-group v-model="viewMode" button-style="solid" size="small">
        <a-radio-button value="25d">2.5D 试验</a-radio-button>
        <a-radio-button value="plane">平面组态兜底</a-radio-button>
      </a-radio-group>
      <template v-if="viewMode === '25d'">
        <span class="xunzhi-nav25d-label">机房/楼层</span>
        <a-select
          v-model="activeFloor"
          size="small"
          style="width: 200px"
          :options="floorOptions"
        />
      </template>
      <a-button v-if="viewMode === 'plane'" size="small" type="link" @click="viewMode = '25d'">
        返回 2.5D 试验
      </a-button>
      <a-button size="small" @click="$router.push('/xunzhi/jf/viz/plane')">打开环境平面大屏</a-button>
    </div>
    <div class="xunzhi-nav25d-body">
      <xunzhi-jf-viz25d v-if="viewMode === '25d'" :key="vizKey" :floor-id="activeFloor" />
      <viz-display-host v-else :page-name="pageName" />
    </div>
  </div>
</template>

<script>
import VizDisplayHost from '@/xunzhiJf/pages/VizDisplayHost.vue'
import XunzhiJfViz25d from '@/xunzhiJf/components/XunzhiJfViz25d.vue'
import { XUNZHI_JF_NAV_25D_PAGE } from '@/xunzhiJf/resolveDisplayPage'

const FLOOR_OPTIONS = [
  { value: 'hall-a', label: 'A 区主机房' },
  { value: 'hall-b', label: 'B 区配电间' },
  { value: 'hall-c', label: 'C 区电池间' },
]

export default {
  name: 'XunzhiJfVizNav25dPage',
  components: { VizDisplayHost, XunzhiJfViz25d },
  data() {
    return {
      viewMode: '25d',
      activeFloor: 'hall-a',
      floorOptions: FLOOR_OPTIONS,
      pageName: XUNZHI_JF_NAV_25D_PAGE,
      vizKey: 0,
    }
  },
  watch: {
    viewMode(mode) {
      if (mode === '25d') this.vizKey += 1
    },
  },
}
</script>

<style scoped>
.xunzhi-nav25d-page {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0;
  gap: 8px;
}
.xunzhi-nav25d-toolbar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 12px;
  padding: 0 4px;
}
.xunzhi-nav25d-label {
  font-size: 13px;
  color: rgba(0, 0, 0, 0.65);
}
.xunzhi-nav25d-body {
  flex: 1;
  min-height: 420px;
  position: relative;
}
</style>
