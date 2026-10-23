<template>
  <div class="xunzhi-monitor-overview">
    <a-alert
      type="success"
      show-icon
      message="场地监测 · 实时总览（S0 演示路径）"
      description="温湿度/漏水等详见环境平面大屏；本页为菜单接通占位 + 快捷跳转。"
    />
    <a-card title="监测入口" style="margin-top: 16px">
      <a-button type="primary" @click="$router.push('/xunzhi/jf/monitor/env')">环境监测</a-button>
      <a-button style="margin-left: 8px" @click="$router.push('/xunzhi/jf/viz/plane')">环境平面大屏</a-button>
      <a-button style="margin-left: 8px" @click="$router.push('/xunzhi/jf/alarm/active')">活动告警</a-button>
    </a-card>
  </div>
</template>

<script>
export default { name: 'XunzhiJfMonitorOverview' }
</script>

<style scoped>
.xunzhi-monitor-overview {
  padding: 16px;
}
</style>
