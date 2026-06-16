<template>
  <div
    v-show="visible"
    class="device-hover-tooltip"
    :style="boxStyle"
  >
    <div class="dht-header">
      <span class="dht-col dht-col-idx">序号</span>
      <span class="dht-col dht-col-dev">设备名称</span>
      <span class="dht-col dht-col-name">数据名称</span>
      <span class="dht-col dht-col-val">数值</span>
    </div>
    <div class="dht-body">
      <div v-if="loading" class="dht-tip">数据加载中…</div>
      <div v-else-if="error" class="dht-tip dht-tip-err">{{ error }}</div>
      <div v-else-if="rows.length === 0" class="dht-tip">暂无测点数据</div>
      <div
        v-for="(row, i) in rows"
        v-else
        :key="i"
        class="dht-row"
        :class="{ 'dht-row-odd': i % 2 === 1 }"
      >
        <span class="dht-col dht-col-idx"><i class="dht-badge">{{ i + 1 }}</i></span>
        <span class="dht-col dht-col-dev" :title="deviceName">{{ deviceName }}</span>
        <span class="dht-col dht-col-name" :title="row.name">{{ row.name }}</span>
        <span class="dht-col dht-col-val">{{ row.value }}<em v-if="row.unit" class="dht-unit">{{ row.unit }}</em></span>
      </div>
    </div>
  </div>
</template>

<script>
import { getRealData } from '@/services/device'

const CACHE_TTL = 8000        // ms — reuse a device's point list within this window
const HIDE_DELAY = 180        // ms — debounce hide to avoid flicker between cells
const MAX_W = 560
const MAX_H = 420

export default {
  name: 'device-hover-tooltip',
  data() {
    return {
      visible: false,
      loading: false,
      error: '',
      deviceUuid: '',
      deviceName: '',
      rows: [],
      posX: 0,
      posY: 0,
      reqSeq: 0,
      hideTimer: null,
      cache: {}                // { [uuid]: { ts, rows } }
    }
  },
  computed: {
    boxStyle() {
      return {
        left: this.posX + 'px',
        top: this.posY + 'px',
        maxWidth: MAX_W + 'px',
        maxHeight: MAX_H + 'px'
      }
    }
  },
  created() {
    this.$EventBus.$on('device-hover-show', this.onShow)
    this.$EventBus.$on('device-hover-move', this.onMove)
    this.$EventBus.$on('device-hover-hide', this.onHide)
  },
  beforeDestroy() {
    this.$EventBus.$off('device-hover-show', this.onShow)
    this.$EventBus.$off('device-hover-move', this.onMove)
    this.$EventBus.$off('device-hover-hide', this.onHide)
    if (this.hideTimer) clearTimeout(this.hideTimer)
  },
  methods: {
    onShow(payload) {
      if (!payload || !payload.deviceUuid) return
      if (this.hideTimer) { clearTimeout(this.hideTimer); this.hideTimer = null }
      this.deviceName = payload.deviceName || ''
      this.computePosition(payload.clientX, payload.clientY)
      this.visible = true
      // Same device already shown → keep current rows, just reposition.
      if (this.deviceUuid === payload.deviceUuid && this.rows.length > 0 && !this.error) {
        return
      }
      this.deviceUuid = payload.deviceUuid
      this.loadData(payload.deviceUuid)
    },
    onMove(payload) {
      if (!this.visible || !payload) return
      this.computePosition(payload.clientX, payload.clientY)
    },
    onHide() {
      if (this.hideTimer) clearTimeout(this.hideTimer)
      this.hideTimer = setTimeout(() => {
        this.visible = false
        this.hideTimer = null
      }, HIDE_DELAY)
    },
    loadData(uuid) {
      const cached = this.cache[uuid]
      if (cached && (Date.now() - cached.ts) < CACHE_TTL) {
        this.rows = cached.rows
        this.loading = false
        this.error = ''
        return
      }
      const seq = ++this.reqSeq
      this.loading = true
      this.error = ''
      this.rows = []
      getRealData({ uuid: uuid, IsRemoveGW: false }).then((res) => {
        if (seq !== this.reqSeq) return        // a newer hover superseded this one
        const data = res && res.data ? res.data : {}
        if (data.code !== 0 && data.code !== undefined && data.code !== null && data.code !== 200) {
          // some builds use 0 as success; tolerate missing code
        }
        const list = Array.isArray(data.realData) ? data.realData : []
        const rows = list.map((it) => ({
          name: it.name || it.Name || '',
          value: this.fmtValue(it.value),
          unit: it.unit || it.DataUnit || ''
        }))
        this.rows = rows
        this.loading = false
        this.cache[uuid] = { ts: Date.now(), rows }
      }).catch((e) => {
        if (seq !== this.reqSeq) return
        this.loading = false
        this.error = '测点数据获取失败'
        // eslint-disable-next-line no-console
        console.error('[DeviceHoverTooltip] getRealData error:', e)
      })
    },
    fmtValue(v) {
      if (v === null || v === undefined || v === '') return '—'
      return String(v)
    },
    computePosition(clientX, clientY) {
      const vw = window.innerWidth || document.documentElement.clientWidth
      const vh = window.innerHeight || document.documentElement.clientHeight
      const offset = 16
      // Estimate box size (clamped) for overflow guarding.
      const estW = Math.min(MAX_W, 520)
      const estH = Math.min(MAX_H, 60 + this.rows.length * 30 + 40)
      let x = (clientX || 0) + offset
      let y = (clientY || 0) + offset
      if (x + estW > vw - 8) x = (clientX || 0) - estW - offset
      if (x < 8) x = 8
      if (y + estH > vh - 8) y = vh - estH - 8
      if (y < 8) y = 8
      this.posX = x
      this.posY = y
    }
  }
}
</script>

<style scoped>
.device-hover-tooltip {
  position: fixed;
  z-index: 99999;
  pointer-events: none;            /* never steal hover from the canvas cell */
  display: flex;
  flex-direction: column;
  min-width: 360px;
  border: 1px solid rgba(0, 229, 255, 0.55);
  border-radius: 4px;
  box-shadow: 0 8px 28px rgba(0, 0, 0, 0.55), 0 0 0 1px rgba(0, 229, 255, 0.15) inset;
  background: rgba(7, 24, 38, 0.92);
  backdrop-filter: blur(2px);
  font-family: 'Microsoft YaHei', 'PingFang SC', sans-serif;
  overflow: hidden;
}
.dht-header {
  display: flex;
  align-items: center;
  height: 34px;
  background: linear-gradient(90deg, #16a6d8 0%, #1bb6e6 100%);
  color: #ffffff;
  font-size: 13px;
  font-weight: 600;
  flex: 0 0 auto;
}
.dht-body {
  overflow-y: auto;
  overflow-x: hidden;
}
.dht-row {
  display: flex;
  align-items: center;
  min-height: 30px;
  font-size: 12px;
  color: #dff1fb;
  background: rgba(9, 46, 66, 0.85);
  border-top: 1px solid rgba(0, 229, 255, 0.08);
}
.dht-row-odd {
  background: rgba(6, 30, 46, 0.85);
}
.dht-col {
  padding: 4px 10px;
  box-sizing: border-box;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.dht-col-idx { flex: 0 0 56px;  text-align: center; }
.dht-col-dev { flex: 0 0 130px; }
.dht-col-name { flex: 1 1 180px; }
.dht-col-val { flex: 0 0 110px; color: #00e5ff; font-weight: 600; }
.dht-unit {
  margin-left: 3px;
  font-style: normal;
  font-size: 11px;
  color: #7fb8d4;
  font-weight: 400;
}
.dht-badge {
  display: inline-block;
  min-width: 18px;
  height: 18px;
  line-height: 18px;
  padding: 0 4px;
  font-style: normal;
  font-size: 11px;
  color: #ffffff;
  background: #1391c9;
  border-radius: 3px;
}
.dht-tip {
  padding: 14px;
  font-size: 12px;
  color: #9fc4dc;
  text-align: center;
}
.dht-tip-err { color: #ff8a73; }
.dht-body::-webkit-scrollbar { width: 6px; }
.dht-body::-webkit-scrollbar-thumb {
  background: rgba(0, 229, 255, 0.35);
  border-radius: 3px;
}
</style>
