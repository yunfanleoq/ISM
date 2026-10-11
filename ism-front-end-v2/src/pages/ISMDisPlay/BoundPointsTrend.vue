<template>
  <div ref="host" class="bpp-chart"></div>
</template>

<script>
import * as echarts from 'echarts'

const COLORS = ['#1f8f4a', '#1a6fb5', '#c47b12', '#6b4c9a']

export default {
  name: 'bound-points-trend',
  props: {
    times: { type: Array, default: () => [] },
    series: { type: Array, default: () => [] }
  },
  watch: {
    times() { this.render() },
    series: { deep: true, handler() { this.render() } }
  },
  mounted() {
    this.render()
  },
  beforeDestroy() {
    if (this.chart) {
      this.chart.dispose()
      this.chart = null
    }
  },
  methods: {
    render() {
      const host = this.$refs.host
      if (!host) return
      if (!this.chart) this.chart = echarts.init(host)
      const series = (this.series || []).map((item, index) => ({
        name: item.name,
        type: 'line',
        smooth: true,
        showSymbol: false,
        data: item.data,
        lineStyle: { width: 1.5, color: COLORS[index % COLORS.length] },
        itemStyle: { color: COLORS[index % COLORS.length] }
      }))
      this.chart.setOption({
        color: COLORS,
        grid: { left: 40, right: 12, top: 12, bottom: 36 },
        tooltip: { trigger: 'axis' },
        legend: {
          type: 'scroll',
          bottom: 0,
          textStyle: { fontSize: 11, color: '#1a1a1a' }
        },
        xAxis: {
          type: 'category',
          data: this.times,
          axisLabel: { fontSize: 10, color: '#3d6b4a' },
          axisLine: { lineStyle: { color: '#b7d4be' } }
        },
        yAxis: {
          type: 'value',
          scale: true,
          axisLabel: { fontSize: 10, color: '#3d6b4a' },
          splitLine: { lineStyle: { color: '#d5e6d8' } }
        },
        series: series
      }, true)
      this.chart.resize()
    }
  }
}
</script>

<style scoped>
.bpp-chart {
  height: 168px;
  width: 100%;
  background: #fff;
}
</style>
