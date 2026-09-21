<template>
  <div
    v-show="visible"
    class="bound-points-popup"
    :style="boxStyle"
    @mouseenter="onEnter"
    @mouseleave="onLeave"
  >
    <div class="bpp-header">
      <span class="bpp-col bpp-col-idx">序号</span>
      <span class="bpp-col bpp-col-dev">设备</span>
      <span class="bpp-col bpp-col-name">测点</span>
      <span class="bpp-col bpp-col-val">数值</span>
    </div>
    <div class="bpp-body">
      <div v-if="loading" class="bpp-tip">数据加载中…</div>
      <div v-else-if="error" class="bpp-tip bpp-tip-err">{{ error }}</div>
      <div v-else-if="rows.length === 0" class="bpp-tip">暂无绑定测点</div>
      <div
        v-for="(row, i) in rows"
        v-else
        :key="row.key || i"
        class="bpp-row"
        :class="{ 'bpp-row-odd': i % 2 === 1 }"
      >
        <span class="bpp-col bpp-col-idx"><i class="bpp-badge">{{ i + 1 }}</i></span>
        <span class="bpp-col bpp-col-dev" :title="row.device">{{ row.device }}</span>
        <span class="bpp-col bpp-col-name" :title="row.name">{{ row.name }}</span>
        <span class="bpp-col bpp-col-val">{{ row.value }}<em v-if="row.unit" class="bpp-unit">{{ row.unit }}</em></span>
      </div>
    </div>
  </div>
</template>

<script>
import { getRealDataByUuid } from '@/services/device'

const HIDE_DELAY = 180
const MAX_W = 560
const MAX_H = 420

export default {
  name: 'bound-points-popup',
  data() {
    return {
      visible: false,
      pinned: false,
      loading: false,
      error: '',
      rows: [],
      posX: 0,
      posY: 0,
      reqSeq: 0,
      hideTimer: null
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
  beforeDestroy() {
    if (this.hideTimer) clearTimeout(this.hideTimer)
  },
  methods: {
    show(payload) {
      if (!payload || !Array.isArray(payload.binds) || payload.binds.length === 0) {
        return
      }
      if (this.hideTimer) {
        clearTimeout(this.hideTimer)
        this.hideTimer = null
      }
      this.computePosition(payload.clientX, payload.clientY)
      this.visible = true
      this.loadData(payload.binds)
    },
    togglePin(payload) {
      if (this.pinned && this.visible) {
        this.pinned = false
        this.hide(true)
        return
      }
      this.pinned = true
      this.show(payload)
    },
    move(payload) {
      if (!this.visible || this.pinned || !payload) return
      this.computePosition(payload.clientX, payload.clientY)
    },
    hide(force) {
      if (this.pinned && !force) return
      if (this.hideTimer) clearTimeout(this.hideTimer)
      this.hideTimer = setTimeout(() => {
        this.visible = false
        this.pinned = false
        this.hideTimer = null
      }, force ? 0 : HIDE_DELAY)
    },
    onEnter() {
      if (this.hideTimer) {
        clearTimeout(this.hideTimer)
        this.hideTimer = null
      }
    },
    onLeave() {
      this.hide(false)
    },
    loadData(binds) {
      const uuids = binds.map((item) => item && (item.dataID || item.dataId)).filter(Boolean)
      const seq = ++this.reqSeq
      this.loading = true
      this.error = ''
      this.rows = binds.map((item, index) => ({
        key: (item && item.dataID) || index,
        device: (item && (item.DeviceName || item.deviceSN)) || '',
        name: (item && item.dataName) || '',
        value: '—',
        unit: ''
      }))
      if (uuids.length === 0) {
        this.loading = false
        return
      }
      getRealDataByUuid({ uuid: uuids }).then((res) => {
        if (seq !== this.reqSeq) return
        const data = res && res.data ? res.data : {}
        const list = Array.isArray(data.realData) ? data.realData : []
        const byUuid = {}
        list.forEach((it) => {
          const id = it.uuid || it.Uuid
          if (id) byUuid[id] = it
        })
        this.rows = binds.map((item, index) => {
          const id = item && (item.dataID || item.dataId)
          const real = id ? byUuid[id] : null
          return {
            key: id || index,
            device: (real && (real.DeviceName || real.deviceName)) || (item && (item.DeviceName || item.deviceSN)) || '',
            name: (real && (real.name || real.Name)) || (item && item.dataName) || '',
            value: this.fmtValue(real && real.value),
            unit: (real && (real.unit || real.DataUnit)) || ''
          }
        })
        this.loading = false
      }).catch(() => {
        if (seq !== this.reqSeq) return
        this.loading = false
        this.error = '测点数据获取失败'
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
.bound-points-popup {
  position: fixed;
  z-index: 99999;
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
.bpp-header {
  display: flex;
  align-items: center;
  height: 34px;
  background: linear-gradient(90deg, #16a6d8 0%, #1bb6e6 100%);
  color: #ffffff;
  font-size: 13px;
  font-weight: 600;
  flex: 0 0 auto;
}
.bpp-body {
  overflow-y: auto;
  overflow-x: hidden;
}
.bpp-row {
  display: flex;
  align-items: center;
  min-height: 30px;
  color: #d7f6ff;
  font-size: 12px;
}
.bpp-row-odd {
  background: rgba(255, 255, 255, 0.04);
}
.bpp-col {
  padding: 0 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.bpp-col-idx { width: 56px; flex: 0 0 56px; }
.bpp-col-dev { width: 120px; flex: 1 1 120px; }
.bpp-col-name { width: 140px; flex: 1 1 140px; }
.bpp-col-val { width: 110px; flex: 0 0 110px; text-align: right; }
.bpp-badge {
  display: inline-block;
  min-width: 18px;
  padding: 0 4px;
  border-radius: 9px;
  background: #16a6d8;
  color: #fff;
  text-align: center;
  font-style: normal;
}
.bpp-unit {
  margin-left: 4px;
  color: #8ad7ee;
  font-style: normal;
}
.bpp-tip {
  padding: 12px;
  color: #9ad8ea;
}
.bpp-tip-err {
  color: #ff8a8a;
}
</style>
