<template>
  <div
    v-show="visible"
    class="bound-points-popup"
    :style="boxStyle"
    @mouseenter="onEnter"
    @mouseleave="onLeave"
    @mousedown.stop
  >
    <div class="bpp-title">
      <span class="bpp-title-text">{{ cardTitle }}</span>
      <span class="bpp-actions">
        <button type="button" class="bpp-title-hint" @click.stop="onTitleAction">{{ pinned ? '点击关闭' : '点击钉住' }}</button>
        <button type="button" class="bpp-close" title="关闭" @click.stop="closeNow">×</button>
      </span>
    </div>
    <div class="bpp-body">
      <div v-if="loading" class="bpp-tip">数据加载中…</div>
      <div v-else-if="error" class="bpp-tip bpp-tip-err">{{ error }}</div>
      <div v-else-if="rows.length === 0" class="bpp-tip">暂无绑定测点</div>
      <template v-else>
        <div v-for="group in groups" :key="group.name" class="bpp-section">
          <div class="bpp-section-title">{{ group.name }}</div>
          <div
            v-for="(row, i) in group.rows"
            :key="row.key || i"
            class="bpp-row"
            :class="{ 'bpp-row-odd': i % 2 === 1 }"
          >
            <span class="bpp-name" :title="row.name">{{ row.name }}</span>
            <span class="bpp-val">{{ row.value }}<em v-if="row.unit" class="bpp-unit">{{ row.unit }}</em></span>
          </div>
        </div>
        <div class="bpp-section">
          <div class="bpp-section-title">历史曲线 · 近24小时</div>
          <div v-if="historyLoading" class="bpp-tip">曲线加载中…</div>
          <div v-else-if="historySeries.length === 0" class="bpp-tip">暂无历史曲线</div>
          <bound-points-trend v-else :times="historyTimes" :series="historySeries" />
        </div>
        <div class="bpp-section">
          <div class="bpp-section-title">设备档案</div>
          <div v-if="archiveLoading" class="bpp-tip">档案加载中…</div>
          <div v-else-if="archives.length === 0" class="bpp-tip">暂无设备档案</div>
          <template v-else>
            <div v-for="item in archives" :key="item.uuid" class="bpp-archive">
              <div class="bpp-row"><span class="bpp-name">设备名称</span><span class="bpp-val">{{ item.name || '—' }}</span></div>
              <div class="bpp-row bpp-row-odd"><span class="bpp-name">所属区域</span><span class="bpp-val">{{ item.zone || '—' }}</span></div>
              <div class="bpp-row"><span class="bpp-name">数据模型</span><span class="bpp-val">{{ item.modelName || '—' }}</span></div>
              <div class="bpp-row bpp-row-odd"><span class="bpp-name">协议</span><span class="bpp-val">{{ typeLabel(item.deviceType) }}</span></div>
              <div class="bpp-row"><span class="bpp-name">通信状态</span><span class="bpp-val">{{ statusLabel(item.status) }}</span></div>
              <div class="bpp-row bpp-row-odd"><span class="bpp-name">采集周期</span><span class="bpp-val">{{ item.interval ? item.interval + ' 秒' : '—' }}</span></div>
              <div class="bpp-row"><span class="bpp-name">描述</span><span class="bpp-val">{{ item.description || '—' }}</span></div>
            </div>
          </template>
        </div>
      </template>
    </div>
  </div>
</template>

<script>
import { getRealDataByUuid, getDeviceArchive } from '@/services/device'
import { GetChartDataHistoryList } from '@/services/report'
import BoundPointsTrend from './BoundPointsTrend.vue'
import { bindId, deviceId } from './boundPointsIds'

const HIDE_DELAY = 420
const MAX_W = 460
const MAX_H = 760
const HISTORY_MINUTES = 24 * 60
const HISTORY_BUCKET = 60
const MAX_SERIES = 3

export default {
  name: 'bound-points-popup',
  components: { BoundPointsTrend },
  data() {
    return {
      visible: false,
      pinned: false,
      loading: false,
      error: '',
      rows: [],
      archives: [],
      archiveLoading: false,
      historyTimes: [],
      historySeries: [],
      historyLoading: false,
      posX: 0,
      posY: 0,
      reqSeq: 0,
      hideTimer: null,
      suppressTimer: null,
      currentKey: '',
      ownerKey: '',
      suppressKey: ''
    }
  },
  computed: {
    cardTitle() {
      const names = []
      this.rows.forEach((row) => {
        if (row.device && names.indexOf(row.device) === -1) names.push(row.device)
      })
      if (names.length === 0) return '测点数据'
      if (names.length === 1) return names[0]
      return names[0] + ' 等' + names.length + '台'
    },
    groups() {
      const order = []
      const bucket = {}
      this.rows.forEach((row) => {
        const name = row.section || '实时数据'
        if (!bucket[name]) {
          bucket[name] = []
          order.push(name)
        }
        bucket[name].push(row)
      })
      return order.map((name) => ({ name: name, rows: bucket[name] }))
    },
    boxStyle() {
      return {
        left: this.posX + 'px',
        top: this.posY + 'px',
        maxWidth: MAX_W + 'px',
        maxHeight: MAX_H + 'px'
      }
    }
  },
  mounted() {
    document.addEventListener('mousedown', this.onDocPointer, true)
    document.addEventListener('keydown', this.onKeydown)
  },
  beforeDestroy() {
    if (this.hideTimer) clearTimeout(this.hideTimer)
    if (this.suppressTimer) clearTimeout(this.suppressTimer)
    document.removeEventListener('mousedown', this.onDocPointer, true)
    document.removeEventListener('keydown', this.onKeydown)
  },
  methods: {
    bindKey(binds) {
      return (binds || []).map((item) => bindId(item) + '@' + deviceId(item)).join('|')
    },
    show(payload, options) {
      if (!payload || !Array.isArray(payload.binds) || payload.binds.length === 0) {
        return
      }
      const key = this.bindKey(payload.binds)
      const ownerKey = payload.ownerKey || key
      const force = !!(options && options.force)
      if (!force && this.suppressKey && ownerKey === this.suppressKey) {
        return
      }
      if (force) this.clearSuppress(ownerKey)
      if (this.pinned && this.visible && key !== this.currentKey && !force) {
        return
      }
      if (this.hideTimer) {
        clearTimeout(this.hideTimer)
        this.hideTimer = null
      }
      this.ownerKey = ownerKey
      const same = this.visible && key && key === this.currentKey
      if (!same) {
        this.computePosition(payload.clientX, payload.clientY)
        this.currentKey = key
        this.loadData(payload.binds)
      }
      this.visible = true
    },
    closeNow() {
      if (this.hideTimer) {
        clearTimeout(this.hideTimer)
        this.hideTimer = null
      }
      if (this.ownerKey) this.suppressKey = this.ownerKey
      this.armSuppressFallback()
      this.pinned = false
      this.visible = false
      this.currentKey = ''
    },
    armSuppressFallback() {
      if (this.suppressTimer) clearTimeout(this.suppressTimer)
      this.suppressTimer = setTimeout(() => {
        this.suppressTimer = null
        this.suppressKey = ''
      }, 500)
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
      const key = this.bindKey(payload && payload.binds)
      if (this.pinned && this.visible && key === this.currentKey) {
        this.closeNow()
        return
      }
      this.pinned = true
      this.show(payload, { force: true })
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
    move() {
      // 卡片停在出现位置。跟着鼠标走时，点不到「钉住」和「关闭」。
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
        this.currentKey = ''
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
      const list = Array.isArray(binds) ? binds : []
      const uuids = []
      const devices = []
      list.forEach((item) => {
        const id = bindId(item)
        const dev = deviceId(item)
        if (id && uuids.indexOf(id) === -1) uuids.push(id)
        if (dev && devices.indexOf(dev) === -1) devices.push(dev)
      })
      const seq = ++this.reqSeq
      this.loading = true
      this.error = ''
      this.archives = []
      this.historyTimes = []
      this.historySeries = []
      this.rows = list.map((item, index) => ({
        key: bindId(item) || index,
        device: (item && (item.DeviceName || deviceId(item))) || '',
        deviceUuid: deviceId(item),
        modelDataUuid: '',
        name: (item && item.dataName) || '',
        section: (item && item.section) || '',
        value: '—',
        unit: (item && (item.unit || item.dataUnit)) || '',
        record: 0
      }))
      if (uuids.length === 0) {
        this.loading = false
        this.loadExtras(seq)
        return
      }
      const params = { uuid: uuids }
      if (devices.length) params.devices = devices
      getRealDataByUuid(params).then((res) => {
        if (seq !== this.reqSeq) return
        const data = res && res.data ? res.data : {}
        const realList = Array.isArray(data.realData) ? data.realData : []
        const byUuid = {}
        const byModel = {}
        realList.forEach((it) => {
          const id = it.uuid || it.Uuid
          const modelId = it.mduid || it.ModelDataUuid
          const dev = it.duid || it.DeviceUuid || ''
          if (id) byUuid[id] = it
          if (modelId) {
            byModel[dev + '|' + modelId] = it
            if (!byModel[modelId]) byModel[modelId] = it
          }
        })
        this.rows = list.map((item, index) => {
          const id = bindId(item)
          const dev = deviceId(item)
          const real = (id && byUuid[id]) || (id && dev && byModel[dev + '|' + id]) || (id && byModel[id]) || null
          return {
            key: id || index,
            device: (real && (real.DeviceName || real.deviceName)) || (item && (item.DeviceName || dev)) || '',
            deviceUuid: (real && (real.duid || real.DeviceUuid)) || dev,
            modelDataUuid: (real && (real.mduid || real.ModelDataUuid)) || '',
            name: (real && (real.name || real.Name)) || (item && item.dataName) || '',
            section: (item && item.section) || '',
            value: this.fmtValue(real && (real.value != null && real.value !== '' ? real.value : real.Value)),
            unit: (real && (real.unit || real.DataUnit)) || (item && (item.unit || item.dataUnit)) || '',
            record: real && real.record ? real.record : 0
          }
        })
        this.loading = false
        this.loadExtras(seq)
      }).catch(() => {
        if (seq !== this.reqSeq) return
        this.loading = false
        this.error = '测点数据获取失败'
      })
    },
    loadExtras(seq) {
      const deviceIds = []
      this.rows.forEach((row) => {
        if (row.deviceUuid && deviceIds.indexOf(row.deviceUuid) === -1) deviceIds.push(row.deviceUuid)
      })
      this.loadArchive(seq, deviceIds)
      this.loadHistory(seq)
    },
    loadArchive(seq, deviceIds) {
      if (deviceIds.length === 0) {
        this.archives = []
        this.archiveLoading = false
        return
      }
      this.archiveLoading = true
      getDeviceArchive({ uuid: deviceIds }).then((res) => {
        if (seq !== this.reqSeq) return
        const data = res && res.data ? res.data : {}
        this.archives = Array.isArray(data.list) ? data.list : []
        this.archiveLoading = false
      }).catch(() => {
        if (seq !== this.reqSeq) return
        this.archives = []
        this.archiveLoading = false
      })
    },
    loadHistory(seq) {
      const picked = []
      const numeric = this.rows.filter((row) => row.deviceUuid && row.modelDataUuid && this.isNumeric(row.value))
      const recorded = numeric.filter((row) => row.record === 1)
      const source = recorded.length ? recorded : numeric
      source.forEach((row) => {
        if (picked.length >= MAX_SERIES) return
        picked.push(row)
      })
      if (picked.length === 0) {
        this.historyTimes = []
        this.historySeries = []
        this.historyLoading = false
        return
      }
      this.historyLoading = true
      GetChartDataHistoryList({
        TimeIn: HISTORY_BUCKET,
        HistoryTime: HISTORY_MINUTES,
        List: picked.map((row) => ({
          DeviceUuid: row.deviceUuid,
          ModelDataUuid: row.modelDataUuid
        }))
      }).then((res) => {
        if (seq !== this.reqSeq) return
        const body = res && res.data ? res.data : {}
        const list = body.code === 0 && Array.isArray(body.data) ? body.data : []
        this.applyHistory(picked, list)
        this.historyLoading = false
      }).catch(() => {
        if (seq !== this.reqSeq) return
        this.historyTimes = []
        this.historySeries = []
        this.historyLoading = false
      })
    },
    applyHistory(picked, list) {
      const times = []
      const buckets = picked.map(() => [])
      list.forEach((point) => {
        const stamp = point && point.HistoryRecordDateTime
        if (!stamp) return
        times.push(this.fmtTime(stamp))
        const items = Array.isArray(point.dataList) ? point.dataList : []
        picked.forEach((row, index) => {
          const hit = items.find((item) => item && item.DeviceUuid === row.deviceUuid && item.ModelDataUuid === row.modelDataUuid)
          const count = hit && hit.count != null ? Number(hit.count) : 1
          const raw = hit ? (hit.average != null && hit.average !== '-' ? hit.average : hit.Value) : null
          const num = hit && count > 0 && this.isNumeric(raw) ? Number(raw) : null
          buckets[index].push(num)
        })
      })
      const series = []
      picked.forEach((row, index) => {
        if (buckets[index].some((v) => v != null)) {
          series.push({ name: row.name || row.device, data: buckets[index] })
        }
      })
      this.historyTimes = series.length ? times : []
      this.historySeries = series
    },
    fmtTime(stamp) {
      const text = String(stamp)
      return text.length >= 16 ? text.slice(5, 16) : text
    },
    isNumeric(v) {
      if (v === null || v === undefined || v === '' || v === '—' || v === '-') return false
      return !isNaN(Number(v))
    },
    statusLabel(status) {
      if (status === 1) return '在线'
      if (status === 2) return '未激活'
      if (status === 3) return '已停用'
      return '离线'
    },
    typeLabel(name) {
      if (!name) return '—'
      const text = this.$t(name)
      return text && text !== name ? text : name
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
  min-width: 320px;
  width: 420px;
  border: 1px solid #1f8f4a;
  border-radius: 2px;
  box-shadow: 0 10px 28px rgba(0, 0, 0, 0.35);
  background: #f4f7f4;
  font-family: 'Microsoft YaHei', 'PingFang SC', sans-serif;
  overflow: hidden;
  color: #1a1a1a;
}
.bpp-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 32px;
  padding: 0 10px;
  background: #1f8f4a;
  color: #fff;
  flex: 0 0 auto;
}
.bpp-title-text {
  font-size: 14px;
  font-weight: 600;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.bpp-actions {
  display: flex;
  align-items: center;
  flex: 0 0 auto;
  margin-left: 8px;
}
.bpp-title-hint {
  flex: 0 0 auto;
  border: 0;
  background: transparent;
  color: #fff;
  font-size: 12px;
  opacity: 0.95;
  cursor: pointer;
  padding: 0;
}
.bpp-close {
  margin-left: 8px;
  border: 0;
  background: transparent;
  color: #fff;
  font-size: 18px;
  line-height: 1;
  cursor: pointer;
  padding: 0 2px;
}
.bpp-body {
  overflow-y: auto;
  overflow-x: hidden;
  max-height: 700px;
}
.bpp-section-title {
  height: 26px;
  line-height: 26px;
  padding: 0 10px;
  background: #2ea85a;
  color: #fff;
  font-size: 13px;
  font-weight: 600;
}
.bpp-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-height: 28px;
  padding: 0 10px;
  font-size: 13px;
  background: #fff;
}
.bpp-row-odd {
  background: #eef6ee;
}
.bpp-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-right: 12px;
}
.bpp-val {
  flex: 0 0 auto;
  font-variant-numeric: tabular-nums;
  color: #0b3d1e;
}
.bpp-unit {
  margin-left: 4px;
  color: #3d6b4a;
  font-style: normal;
}
.bpp-tip {
  padding: 12px;
  color: #3d6b4a;
}
.bpp-tip-err {
  color: #c43838;
}
</style>
