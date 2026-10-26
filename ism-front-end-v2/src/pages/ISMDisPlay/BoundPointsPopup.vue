<template>
  <div
    v-show="visible"
    class="bound-points-popup"
    :style="boxStyle"
    @mouseenter="onEnter"
    @mouseleave="onLeave"
    @mousedown.stop
    @click.stop
  >
    <div class="bpp-header">
      <span class="bpp-col bpp-col-idx">序号</span>
      <span class="bpp-col bpp-col-dev">设备</span>
      <span class="bpp-col bpp-col-name">测点</span>
      <span class="bpp-col bpp-col-val">数值</span>
      <button type="button" class="bpp-pin" @click.stop="onTitleAction">{{ pinned ? '点击关闭' : '点击钉住' }}</button>
      <button type="button" class="bpp-close" title="关闭" @click.stop="closeNow">×</button>
    </div>
    <div class="bpp-body" :style="bodyStyle">
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
const MAX_W = 720
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
      hideTimer: null,
      suppressTimer: null,
      ownerKey: '',
      suppressKey: '',
      boxMaxH: MAX_H
    }
  },
  computed: {
    boxStyle() {
      return {
        left: this.posX + 'px',
        top: this.posY + 'px',
        width: 'min(' + MAX_W + 'px, calc(100vw - 16px))',
        maxWidth: 'calc(100vw - 16px)',
        maxHeight: (this.boxMaxH || MAX_H) + 'px'
      }
    },
    bodyStyle() {
      const h = Math.max(120, (this.boxMaxH || MAX_H) - 34)
      return { maxHeight: h + 'px' }
    }
  },
  created() {
    this._popupId = 'bpp-' + Math.random().toString(36).slice(2)
    if (this.$EventBus) {
      this.$EventBus.$on('ism-bound-points-claim', this.onOtherClaim)
    }
  },
  mounted() {
    if (this.$el && this.$el.parentNode !== document.body) {
      document.body.appendChild(this.$el)
    }
    document.addEventListener('mousedown', this.onDocPointer, true)
    document.addEventListener('keydown', this.onKeydown)
  },
  beforeDestroy() {
    if (this.$EventBus) {
      this.$EventBus.$off('ism-bound-points-claim', this.onOtherClaim)
    }
    if (this.hideTimer) clearTimeout(this.hideTimer)
    if (this.suppressTimer) clearTimeout(this.suppressTimer)
    document.removeEventListener('mousedown', this.onDocPointer, true)
    document.removeEventListener('keydown', this.onKeydown)
  },
  methods: {
    onOtherClaim(id) {
      if (id !== this._popupId) this.hide(true)
    },
    show(payload, options) {
      if (!payload || !Array.isArray(payload.binds) || payload.binds.length === 0) {
        return
      }
      const ownerKey = payload.ownerKey || ''
      const force = !!(options && options.force)
      if (!force && this.suppressKey && ownerKey && ownerKey === this.suppressKey) {
        return
      }
      if (force) this.clearSuppress(ownerKey)
      if (this.$EventBus) this.$EventBus.$emit('ism-bound-points-claim', this._popupId)
      if (this.hideTimer) {
        clearTimeout(this.hideTimer)
        this.hideTimer = null
      }
      this.ownerKey = ownerKey || this.ownerKey
      this.computePosition(payload.clientX, payload.clientY)
      this.visible = true
      this.loadData(payload.binds)
    },
    closeNow() {
      if (this.hideTimer) {
        clearTimeout(this.hideTimer)
        this.hideTimer = null
      }
      if (this.ownerKey) this.suppressKey = this.ownerKey
      this.pinned = false
      this.visible = false
      this.armSuppressClear()
    },
    armSuppressClear() {
      if (!this.suppressKey) return
      if (this.suppressTimer) clearTimeout(this.suppressTimer)
      this.suppressTimer = setTimeout(() => {
        this.suppressTimer = null
        this.suppressKey = ''
      }, 320)
    },
    cancelSuppressClear() {
      if (this.suppressTimer) {
        clearTimeout(this.suppressTimer)
        this.suppressTimer = null
      }
    },
    clearSuppress(ownerKey) {
      if (this.suppressKey && ownerKey && ownerKey !== this.suppressKey) return
      this.suppressKey = ''
      if (this.suppressTimer) {
        clearTimeout(this.suppressTimer)
        this.suppressTimer = null
      }
    },
    togglePin(payload) {
      if (this.pinned && this.visible) {
        this.closeNow()
        return
      }
      this.pinned = true
      if (payload) this.show(payload, { force: true })
    },
    onTitleAction() {
      if (this.pinned) {
        this.closeNow()
        return
      }
      this.pinned = true
      this.clearSuppress(this.ownerKey)
    },
    onDocPointer(evt) {
      if (!this.pinned || !this.visible) return
      const el = this.$el
      if (el && evt.target && el.contains(evt.target)) return
      this.closeNow()
    },
    onKeydown(evt) {
      if (!this.visible) return
      if (evt.key === 'Escape' || evt.key === 'Esc') this.closeNow()
    },
    move(payload) {
      if (!this.visible || this.pinned || !payload) return
      this.computePosition(payload.clientX, payload.clientY)
    },
    hide(force) {
      if (force) {
        this.closeNow()
        return
      }
      if (this.pinned) return
      if (this.hideTimer) clearTimeout(this.hideTimer)
      this.hideTimer = setTimeout(() => {
        this.hideTimer = null
        if (this.pinned) return
        this.pinned = false
        this.visible = false
      }, HIDE_DELAY)
    },
    onEnter() {
      if (this.hideTimer) {
        clearTimeout(this.hideTimer)
        this.hideTimer = null
      }
    },
    onLeave() {
      if (this.pinned) return
      this.closeNow()
    },
    loadData(binds) {
      const uuids = binds.map((item) => item && (item.dataID || item.dataId)).filter(Boolean)
      const devices = []
      binds.forEach((item) => {
        const id = item && (item.deviceSN || item.duid)
        if (id && devices.indexOf(id) < 0) devices.push(id)
      })
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
      getRealDataByUuid({ uuid: uuids, devices: devices }).then((res) => {
        if (seq !== this.reqSeq) return
        const data = res && res.data ? res.data : {}
        const list = Array.isArray(data.realData) ? data.realData : []
        const byUuid = {}
        list.forEach((it) => {
          ;[it.uuid, it.Uuid, it.mduid, it.model_data_uuid].forEach((id) => {
            if (id) byUuid[id] = it
          })
        })
        this.rows = binds.map((item, index) => {
          const id = item && (item.dataID || item.dataId)
          const real = id ? byUuid[id] : null
          return {
            key: id || index,
            device: (real && (real.DeviceName || real.deviceName)) || (item && (item.DeviceName || item.deviceSN)) || '',
            name: (real && (real.name || real.Name)) || (item && item.dataName) || '',
            value: this.fmtValue(real && real.value),
            unit: (real && (real.unit || real.DataUnit || real.data_unit)) || (item && (item.unit || item.dataUnit || item.DataUnit)) || ''
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
      const margin = 8
      const offset = 16
      const estW = Math.min(MAX_W, Math.max(360, vw - margin * 2))
      const available = Math.max(160, vh - margin * 2)
      const boxH = Math.min(MAX_H, available)
      let x = (clientX || 0) + offset
      let y = (clientY || 0) + offset
      if (x + estW > vw - margin) x = (clientX || 0) - estW - offset
      if (x < margin) x = margin
      if (y + boxH > vh - margin) y = vh - margin - boxH
      if (y < margin) y = margin
      this.posX = x
      this.posY = y
      this.boxMaxH = Math.max(160, Math.min(MAX_H, vh - y - margin))
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
  flex: 0 0 34px;
}
.bpp-pin {
  flex: 0 0 auto;
  margin-left: 8px;
  border: 0;
  background: transparent;
  color: #fff;
  font-size: 12px;
  cursor: pointer;
  padding: 0;
}
.bpp-close {
  flex: 0 0 28px;
  width: 28px;
  height: 28px;
  margin-right: 4px;
  border: 0;
  background: transparent;
  color: #fff;
  font-size: 18px;
  line-height: 28px;
  cursor: pointer;
}
.bpp-body {
  flex: 1 1 auto;
  min-height: 0;
  overflow-y: auto;
  overflow-x: hidden;
  overscroll-behavior: contain;
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
.bpp-col-idx { width: 48px; flex: 0 0 48px; }
.bpp-col-dev { width: 90px; flex: 1 1 72px; min-width: 0; }
.bpp-col-name {
  flex: 1 1 auto;
  min-width: 180px;
  white-space: normal;
  word-break: break-all;
  overflow: visible;
  text-overflow: clip;
  line-height: 1.4;
  padding-top: 4px;
  padding-bottom: 4px;
}
.bpp-col-val { width: 128px; flex: 0 0 128px; text-align: right; }
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
