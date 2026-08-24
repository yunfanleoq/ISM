<template>
  <!--
    运行态浮层：定时拉取真实告警表 GetCurrentAlarmList。
    RC08bate-20260724：顶部 KPI 四卡已从大屏删除，不再覆盖「在线设备/活跃告警」数值；
    仅保留右下角活跃告警列表面板（可配条数、可滚动；「历史查询」开全屏抽屉）。
    仅在总览页(page === modelId)显示；钻探到子页时自动隐藏。
  -->
  <div class="scada-alarm-root">
    <div v-show="visible" class="scada-alarm" :style="panelStyle">
    <div class="sa-inner">
      <div class="sa-head">
        <span class="sa-title">活跃告警</span>
        <span class="sa-head-right">
          <button
            class="sa-history-btn"
            type="button"
            title="打开告警历史查询"
            aria-label="打开告警历史查询"
            @click="openHistoryDrawer"
          >历史查询</button>
          <span class="sa-delay-label">条数</span>
          <input
            v-model.number="rowLimit"
            class="sa-delay-input sa-row-limit-input"
            type="number"
            min="5"
            max="200"
            step="1"
            title="面板最多展示条数，超出可滚动查看；修改后自动保存"
            @change="saveRowLimit"
          >
          <span class="sa-delay-label">延迟</span>
          <input
            v-model.number="startupAlarmDelayMinutes"
            class="sa-delay-input"
            type="number"
            min="0"
            max="1440"
            step="1"
            title="服务启动后 N 分钟内不启用告警判定（仅建基线），结束后新边沿才告警；修改后自动保存"
            :disabled="savingStartupDelay"
            @change="saveStartupAlarmDelay"
          >
          <span class="sa-delay-label">分</span>
          <span class="sa-badge" :class="{ 'sa-badge-ok': activeCount === 0 }">
            ● {{ activeCount }} 条
          </span>
          <button
            class="sa-clear-alarm"
            type="button"
            title="一键清除当前项目告警"
            aria-label="一键清除当前项目告警"
            :disabled="clearing"
            @click="clearAllAlarms"
          >🧹</button>
        </span>
      </div>

      <div v-if="loading && !alarms.length" class="sa-empty sa-loading">⏳ 加载告警中…</div>

      <div v-else-if="!alarms.length" class="sa-empty sa-ok">
        <div class="sa-ok-icon">✓ 当前无活跃告警</div>
        <div class="sa-ok-sub">全园区设备运行正常</div>
      </div>

      <div v-else class="sa-list">
        <div
          v-for="(a, i) in shownAlarms"
          :key="a.ID || (a.DeviceUuid + a.DataUuid + i)"
          class="sa-row"
          :class="{ 'sa-row-busy': clearingRow }"
          title="点击清除该条告警"
          @click="clearOneAlarm(a)"
        >
          <span class="sa-dot" :style="{ background: levelColor(a.AlarmLevel) }"></span>
          <span class="sa-dev" :title="a.DeviceName">{{ a.DeviceName }}</span>
          <span class="sa-name" :title="alarmText(a)">{{ alarmText(a) }}</span>
          <span class="sa-time">{{ shortTime(a.HappenTime) }}</span>
        </div>
        <div v-if="moreCount > 0" class="sa-more">
          已显示 {{ rowLimit }} / {{ activeCount }} · 调大「条数」可看更多
        </div>
      </div>
    </div>
    </div>

    <ScadaAlarmHistoryDrawer
      :visible.sync="historyDrawerVisible"
      :project-uuid="projectUuid"
    />
  </div>
</template>

<script>
import { GetCurrentAlarmList, ClearAllCurrentAlarm, UpdateCurrentAlarm } from '@/services/alarm'
import { GetAlarmNoticeByType, UpdateAlarmNoticeByType } from '@/services/alarmNotice'
import ScadaAlarmHistoryDrawer from './ScadaAlarmHistoryDrawer.vue'

const POLL_INTERVAL = 15000 // 15s 轮询
const ROW_LIMIT_KEY = 'ism.scadaAlarmRowLimit'
const DEFAULT_ROW_LIMIT = 20
const MIN_ROW_LIMIT = 5
const MAX_ROW_LIMIT = 200

function readStoredRowLimit() {
  try {
    const n = Number(localStorage.getItem(ROW_LIMIT_KEY))
    if (Number.isInteger(n) && n >= MIN_ROW_LIMIT && n <= MAX_ROW_LIMIT) return n
  } catch (e) { /* ignore */ }
  return DEFAULT_ROW_LIMIT
}

export default {
  name: 'ScadaAlarmPanel',
  components: { ScadaAlarmHistoryDrawer },
  props: {
    projectUuid: { type: String, required: true },
    modelId: { type: String, required: true },
  },
  data() {
    return {
      loading: true,
      alarms: [],
      currentPage: this.modelId, // 初始即总览页
      pollTimer: null,
      reqSeq: 0, // 注意：不能用下划线前缀，Vue 不会代理 data 中以 _ 开头的属性
      clearing: false,
      clearingRow: false,
      startupAlarmDelayMinutes: 10,
      savingStartupDelay: false,
      rowLimit: readStoredRowLimit(),
      historyDrawerVisible: false,
    }
  },
  computed: {
    // 仅总览页（page === modelId）显示告警面板；钻探子页时隐藏
    visible() {
      return this.currentPage === this.modelId
    },
    activeCount() {
      return this.alarms.length
    },
    shownAlarms() {
      const limit = Math.max(MIN_ROW_LIMIT, Math.min(MAX_ROW_LIMIT, Number(this.rowLimit) || DEFAULT_ROW_LIMIT))
      return this.alarms.slice(0, limit)
    },
    moreCount() {
      const limit = Math.max(MIN_ROW_LIMIT, Math.min(MAX_ROW_LIMIT, Number(this.rowLimit) || DEFAULT_ROW_LIMIT))
      return Math.max(0, this.alarms.length - limit)
    },
    // 画布 1920×1080：右下角告警区（与 build_ncc_dashboard 单功率趋势 + alarm_h=540 对齐）
    // 无顶部 KPI 时 alarm_y≈516；历史查询走标题栏按钮 → 全屏抽屉
    panelStyle() {
      return {
        left: `${(1312 / 1920) * 100}vw`,
        top: `${(516 / 1080) * 100}vh`,
        width: `${(584 / 1920) * 100}vw`,
        height: `${(540 / 1080) * 100}vh`,
      }
    },
  },
  mounted() {
    this.$EventBus && this.$EventBus.$on('GoPage', this.onGoPage)
    this.fetchAlarms()
    this.fetchStartupAlarmDelay()
    this.pollTimer = setInterval(() => {
      this.fetchAlarms()
    }, POLL_INTERVAL)
  },
  beforeDestroy() {
    if (this.pollTimer) clearInterval(this.pollTimer)
    this.$EventBus && this.onGoPage && this.$EventBus.$off('GoPage', this.onGoPage)
  },
  methods: {
    openHistoryDrawer() {
      this.historyDrawerVisible = true
    },
    // 弹窗不改变底层页面；普通切换才更新当前页
    onGoPage(data) {
      if (!data || data.IsPopUp) return
      if (data.PageUuid) this.currentPage = data.PageUuid
    },
    async fetchAlarms() {
      // 子页隐藏时不必拉取，省请求
      if (!this.visible) return
      const seq = ++this.reqSeq
      try {
        const res = await GetCurrentAlarmList(
          { deviceList: [], dataList: [] },
          { headers: { ProjectUuid: this.projectUuid } }
        )
        if (seq !== this.reqSeq) return // 防竞态：仅采纳最新一次
        const list = res && res.data && res.data.code === 0 && Array.isArray(res.data.list)
          ? res.data.list
          : []
        this.alarms = list
      } catch (e) {
        if (seq === this.reqSeq) {
          console.warn('[ScadaAlarmPanel] 获取活跃告警失败:', e && e.message)
        }
      } finally {
        if (seq === this.reqSeq) this.loading = false
      }
    },
    saveRowLimit() {
      let n = Number(this.rowLimit)
      if (!Number.isInteger(n) || n < MIN_ROW_LIMIT || n > MAX_ROW_LIMIT) {
        this.$message.error(`展示条数必须是 ${MIN_ROW_LIMIT}～${MAX_ROW_LIMIT} 的整数`)
        this.rowLimit = readStoredRowLimit()
        return
      }
      this.rowLimit = n
      try {
        localStorage.setItem(ROW_LIMIT_KEY, String(n))
      } catch (e) { /* ignore */ }
    },
    async clearOneAlarm(item) {
      if (!item || this.clearingRow || this.clearing) return
      const deviceUuid = item.DeviceUuid || item.deviceUuid || item.duid
      const dataUuid = item.DataUuid || item.dataUuid || item.uuid
      if (!deviceUuid || !dataUuid) {
        this.$message.error('该告警缺少设备或测点标识，无法清除')
        return
      }
      this.clearingRow = true
      try {
        const res = await UpdateCurrentAlarm({
          type: 1,
          update: {
            duid: deviceUuid,
            uuid: dataUuid,
          }
        })
        if (res && res.data && res.data.code === 0) {
          this.alarms = this.alarms.filter(a =>
            !((a.DeviceUuid || a.deviceUuid) === deviceUuid && (a.DataUuid || a.dataUuid) === dataUuid)
          )
          this.$message.success('该条告警已清除')
          await this.fetchAlarms()
        } else {
          this.$message.error('清除该条告警失败')
        }
      } catch (e) {
        this.$message.error('清除该条告警失败')
      } finally {
        this.clearingRow = false
      }
    },
    async clearAllAlarms() {
      if (this.clearing) return
      if (!window.confirm(
        '确定清除当前项目的全部实时告警吗？\n\n' +
        '说明：仅消除实时状态，告警记录会保留并可在「历史查询」中查看，不会物理删除。'
      )) return
      this.clearing = true
      try {
        const res = await ClearAllCurrentAlarm(
          {skipOfflineResync:true},
          {headers: {ProjectUuid: this.projectUuid}}
        )
        if (res && res.data && res.data.code === 0) {
          this.alarms = []
          this.$message.success('实时状态已清除，记录已进入历史（' + (res.data.count || 0) + '）')
          await this.fetchAlarms()
        } else {
          this.$message.error('清除实时告警失败')
        }
      } catch (e) {
        this.$message.error('清除实时告警失败')
      } finally {
        this.clearing = false
      }
    },
    async fetchStartupAlarmDelay() {
      try {
        const res = await GetAlarmNoticeByType(
          {type: 'StartupAlarmDelay'},
          {headers: {ProjectUuid: this.projectUuid}}
        )
        const params = res && res.data && res.data.code === 0 && res.data.list
          ? JSON.parse(res.data.list.AlarmNoticeParams)
          : null
        if (params && Number.isInteger(Number(params.DelayMinutes))) {
          this.startupAlarmDelayMinutes = Number(params.DelayMinutes)
        }
      } catch (e) {
        console.warn('[ScadaAlarmPanel] 获取启动告警延迟失败:', e && e.message)
      }
    },
    async saveStartupAlarmDelay() {
      const minutes = Number(this.startupAlarmDelayMinutes)
      if (!Number.isInteger(minutes) || minutes < 0 || minutes > 1440) {
        this.$message.error('延迟时间必须是 0 到 1440 的整数分钟')
        await this.fetchStartupAlarmDelay()
        return
      }
      this.savingStartupDelay = true
      try {
        const res = await UpdateAlarmNoticeByType(
          {type: 'StartupAlarmDelay', params: JSON.stringify({DelayMinutes: minutes})},
          {headers: {ProjectUuid: this.projectUuid}}
        )
        if (!res || !res.data || res.data.code !== 0) {
          this.$message.error('告警延迟保存失败')
          await this.fetchStartupAlarmDelay()
        }
      } catch (e) {
        this.$message.error('告警延迟保存失败')
        await this.fetchStartupAlarmDelay()
      } finally {
        this.savingStartupDelay = false
      }
    },
    alarmText(a) {
      const msg = a.AlarmMessage || a.AlarmName || '告警'
      return this.$te(msg) ? this.$t(msg) : msg
    },
    levelColor(level) {
      switch (Number(level)) {
        case 4: return '#ff3b30' // 致命
        case 3: return '#ff6b35' // 紧急
        case 2: return '#ffd60a' // 重要
        default: return '#3b82f6' // 次要/提示
      }
    },
    shortTime(t) {
      if (!t) return ''
      const d = new Date(t)
      if (isNaN(d.getTime())) return String(t).slice(11, 19)
      const p = n => String(n).padStart(2, '0')
      return `${p(d.getHours())}:${p(d.getMinutes())}:${p(d.getSeconds())}`
    },
  },
}
</script>

<style scoped>
/* 透明全屏根容器：承载定位浮层，自身不挡交互 */
.scada-alarm-root {
  position: absolute;
  inset: 0;
  pointer-events: none;
  z-index: 60;
  font-family: 'Microsoft YaHei', 'PingFang SC', sans-serif;
}

.scada-alarm {
  position: absolute;
  pointer-events: auto;
  font-family: 'Microsoft YaHei', 'PingFang SC', sans-serif;
  box-sizing: border-box;
  padding: 0 8px 8px;
}
/* 告警区嵌入右侧统一外框，仅用分隔光轨区分内容，避免完整边框套叠。 */
.sa-inner {
  position: relative;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  background: linear-gradient(180deg, rgba(11, 28, 43, 0.96), rgba(9, 19, 32, 0.96));
  border: 0;
  border-radius: 0;
  padding: 9px 10px 6px;
  box-sizing: border-box;
  box-shadow: inset 0 12px 22px rgba(0, 91, 126, 0.035);
  clip-path: polygon(0 0, 100% 0, 100% 100%, 14px 100%, 0 calc(100% - 14px));
}
.sa-inner::before {
  content: "";
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 1px;
  pointer-events: none;
  background: linear-gradient(90deg, rgba(255, 107, 53, 0.72), rgba(0, 229, 255, 0.25) 28%, transparent 72%);
  box-shadow: 0 0 7px rgba(255, 107, 53, 0.16);
}
.sa-inner::after {
  content: "";
  position: absolute;
  left: -1px;
  bottom: 14px;
  width: 20px;
  height: 1px;
  pointer-events: none;
  transform: rotate(45deg);
  transform-origin: left center;
  background: linear-gradient(90deg, rgba(0, 229, 255, 0.55), rgba(0, 229, 255, 0.12));
  box-shadow: 0 0 5px rgba(0, 229, 255, 0.2);
}
.sa-head {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 6px;
}
.sa-title {
  font-size: 0.72vw;
  font-weight: 600;
  color: #ff6b35;
  letter-spacing: 1px;
}
.sa-badge {
  font-size: 0.6vw;
  color: #ff6b35;
  font-weight: 600;
}
.sa-badge-ok { color: #10e0a0; }
.sa-head-right {
  display: flex;
  align-items: center;
  gap: 4px;
}
.sa-delay-label {
  color: #5f7799;
  font-size: 0.55vw;
}
.sa-delay-input {
  width: 22px;
  height: 16px;
  padding: 0 2px;
  border: 1px solid rgba(0, 229, 255, 0.3);
  border-radius: 2px;
  background: rgba(8, 30, 48, 0.9);
  color: #9fefff;
  font-size: 0.56vw;
  line-height: 14px;
  text-align: center;
  outline: none;
  appearance: textfield;
}
.sa-row-limit-input {
  width: 28px;
}
.sa-delay-input::-webkit-inner-spin-button,
.sa-delay-input::-webkit-outer-spin-button {
  margin: 0;
  -webkit-appearance: none;
}
.sa-delay-input:focus {
  border-color: rgba(0, 229, 255, 0.8);
}
.sa-clear-alarm {
  width: 20px;
  height: 18px;
  padding: 0;
  border: 1px solid rgba(0, 229, 255, 0.45);
  border-radius: 3px;
  background: rgba(0, 229, 255, 0.08);
  color: #00e5ff;
  cursor: pointer;
  font-size: 12px;
  font-weight: 500;
  line-height: 16px;
}
.sa-clear-alarm:hover:not(:disabled) {
  background: rgba(0, 229, 255, 0.2);
}
.sa-clear-alarm:disabled {
  border-color: rgba(95, 119, 153, 0.3);
  color: #5f7799;
  cursor: not-allowed;
}
.sa-history-btn {
  height: 18px;
  padding: 0 6px;
  margin-right: 2px;
  border: 1px solid rgba(255, 107, 53, 0.55);
  border-radius: 3px;
  background: rgba(255, 107, 53, 0.12);
  color: #ff9a6b;
  cursor: pointer;
  font-size: 0.55vw;
  font-weight: 600;
  line-height: 16px;
  letter-spacing: 0.5px;
  white-space: nowrap;
}
.sa-history-btn:hover {
  background: rgba(255, 107, 53, 0.28);
  color: #ffc4a8;
}

.sa-empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
}
.sa-loading { color: #5f7799; font-size: 0.62vw; }
.sa-ok-icon { color: #10e0a0; font-size: 0.7vw; font-weight: 600; }
.sa-ok-sub { color: #9fb6d6; font-size: 0.58vw; margin-top: 4px; }

.sa-list {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  overflow-x: hidden;
}
.sa-row {
  display: flex;
  align-items: center;
  gap: 6px;
  height: 1.85vh;
  min-height: 22px;
  padding: 0 4px;
  border-radius: 3px;
  background: rgba(16, 29, 51, 0.7);
  margin-bottom: 3px;
  cursor: pointer;
}
.sa-row:hover {
  background: rgba(0, 229, 255, 0.12);
}
.sa-row-busy {
  pointer-events: none;
  opacity: 0.6;
}
.sa-dot {
  flex-shrink: 0;
  width: 7px;
  height: 7px;
  border-radius: 50%;
  box-shadow: 0 0 5px currentColor;
}
.sa-dev {
  flex-shrink: 0;
  max-width: 38%;
  color: #e8f1ff;
  font-size: 0.58vw;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.sa-name {
  flex: 1;
  color: #9fb6d6;
  font-size: 0.58vw;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.sa-time {
  flex-shrink: 0;
  color: #5f7799;
  font-size: 0.54vw;
  font-variant-numeric: tabular-nums;
}
.sa-more {
  color: #5f7799;
  font-size: 0.54vw;
  text-align: center;
  padding: 4px 0 2px;
  flex-shrink: 0;
}

.sa-list::-webkit-scrollbar { width: 5px; }
.sa-list::-webkit-scrollbar-track { background: transparent; }
.sa-list::-webkit-scrollbar-thumb { background: #1e3a5f; border-radius: 3px; }
.sa-list::-webkit-scrollbar-thumb:hover { background: #2c5a8f; }
</style>
