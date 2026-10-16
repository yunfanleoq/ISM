<template>
  <div class="xunzhi-home">
    <a-alert type="info" show-icon message="循智·机房环境工作台（演示骨架）" description="KPI 卡片占位；非正式验收。" />
    <a-row :gutter="16" style="margin-top: 16px">
      <a-col :span="6"><a-card title="活动告警">—</a-card></a-col>
      <a-col :span="6"><a-card title="测点在线">—</a-card></a-col>
      <a-col :span="6"><a-card title="PUE">—</a-card></a-col>
      <a-col :span="6"><a-card title="机柜/U位">—</a-card></a-col>
    </a-row>
    <a-card title="演示快捷入口" style="margin-top: 16px">
      <a-button type="primary" @click="$router.push('/xunzhi/jf/monitor/overview')">实时总览</a-button>
      <a-button style="margin-left: 8px" @click="$router.push('/xunzhi/jf/viz/plane')">环境平面大屏</a-button>
      <a-button style="margin-left: 8px" @click="$router.push('/xunzhi/jf/alarm/active')">活动告警</a-button>
    </a-card>
  </div>
</template>

<script>
export default { name: 'XunzhiJfHomeWorkbench' }
</script>

<style scoped>
.xunzhi-home {
  padding: 16px;
}
</style>
