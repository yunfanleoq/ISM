<template>
  <a-card size="small" title="生产数据" style="margin-top: 8px">
    <p v-if="loading" class="muted">正在读取生产数据…</p>
    <template v-else>
      <p v-if="live.note" class="muted">{{ live.note }}</p>
      <p v-if="live.title && (live.fields || []).length" class="live-title">{{ live.title }}</p>
      <div v-for="(field, i) in live.fields || []" :key="'f' + i" class="live-row">
        <span>{{ field.label }}</span>
        <b>{{ field.value }}</b>
      </div>
      <div v-for="(row, i) in live.rows || []" :key="'r' + i" class="live-row">
        <b>{{ rowText(row) }}</b>
      </div>
      <p v-if="!(live.fields || []).length && !(live.rows || []).length && !live.note" class="muted">还没有生产数据。</p>
    </template>
    <div class="live-log">日志</div>
    <p v-if="!(live.logs || []).length" class="muted">没有相关日志。</p>
    <div v-for="(line, i) in live.logs || []" :key="'l' + i" class="live-line">
      <em>{{ line.time }}</em> {{ line.kind }} {{ line.text }}
    </div>
  </a-card>
</template>

<script>
export default {
  name: 'NodeLive',
  props: {
    live: { type: Object, default: () => ({ fields: [], rows: [], logs: [], note: '' }) },
    loading: { type: Boolean, default: false },
  },
  methods: {
    rowText(row) {
      return (row || []).map(cell => cell.value).filter(Boolean).join(' · ')
    },
  },
}
</script>

<style scoped>
.muted { color: #829ab1; margin: 0 0 6px; }
.live-title { margin: 0 0 6px; font-weight: 600; color: #102a43; }
.live-row { display: flex; justify-content: space-between; gap: 8px; font-size: 12px; padding: 3px 0; border-bottom: 1px solid #f0f4f8; }
.live-row span { color: #627d98; flex: none; }
.live-row b { font-weight: 600; text-align: right; word-break: break-all; }
.live-log { margin: 10px 0 4px; font-size: 12px; font-weight: 600; color: #334e68; }
.live-line { font-size: 12px; color: #334e68; padding: 3px 0; word-break: break-all; }
.live-line em { color: #829ab1; font-style: normal; margin-right: 4px; }
</style>
