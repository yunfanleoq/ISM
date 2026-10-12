<template>
  <div>
    <a-alert
      type="info"
      show-icon
      style="margin: 12px 12px 0"
      message="告警中心 · 活动告警（S1 演示）"
      description="与平台「实时告警」同一数据源 GetCurrentAlarmList；非正式验收。"
    />
    <current-alarm />
  </div>
</template>

<script>
import CurrentAlarm from '@/pages/alarm/currentAlarm/currentAlarm.vue'

export default {
  name: 'XunzhiJfAlarmActiveHost',
  components: { CurrentAlarm },
}
</script>
