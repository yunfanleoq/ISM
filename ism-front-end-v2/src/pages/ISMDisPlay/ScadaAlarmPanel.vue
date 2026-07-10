<template>
  <!--
    运行态「活跃告警」浮层：定时拉取真实告警表 GetCurrentAlarmList（后端 = clear_time
    未消除的活跃告警），替换原本只读 monitor_list.status 的静态快照。覆盖两处：
      1) 右下角告警列表面板（画布 x=1312 y=872 w=584 h=184）
      2) 顶部 KPI「活跃告警」卡数值（画布卡片中心 x≈1674 y=104）
    仅在总览页(page === modelId)显示；钻探到子页时自动隐藏。
  -->
  <div class="scada-alarm-root">
    <!-- 顶部 KPI「活跃告警」实时计数（径向渐隐底，避免硬边遮卡片渐变） -->
    <div v-show="visible" class="sa-kpi" :style="kpiStyle">
      <span class="sa-kpi-num" :style="{ color: kpiColor }">{{ activeCount }}</span>
    </div>

    <div v-show="visible" class="scada-alarm" :style="panelStyle">
    <div class="sa-inner">
      <div class="sa-head">
        <span class="sa-title">活跃告警</span>
        <span class="sa-badge" :class="{ 'sa-badge-ok': activeCount === 0 }">
          ● {{ activeCount }} 条
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
        >
          <span class="sa-dot" :style="{ background: levelColor(a.AlarmLevel) }"></span>
          <span class="sa-dev" :title="a.DeviceName">{{ a.DeviceName }}</span>
          <span class="sa-name" :title="alarmText(a)">{{ alarmText(a) }}</span>
          <span class="sa-time">{{ shortTime(a.HappenTime) }}</span>
        </div>
        <div v-if="moreCount > 0" class="sa-more">
          +{{ moreCount }} 更多告警 · 左侧导航树查看
        </div>
      </div>
    </div>
    </div>
  </div>
</template>

<script>
import { GetCurrentAlarmList } from '@/services/alarm'

const POLL_INTERVAL = 15000 // 15s 轮询
const ROW_LIMIT = 5         // 面板高度有限，最多展示 5 行，其余折叠为“更多”

export default {
  name: 'ScadaAlarmPanel',
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
      return this.alarms.slice(0, ROW_LIMIT)
    },
    moreCount() {
      return Math.max(0, this.alarms.length - ROW_LIMIT)
    },
    // 画布 1920×1080，autoSize=1 铺满视口 → 用 vw/vh 与 cells 坐标对齐
    panelStyle() {
      return {
        left: `${(1312 / 1920) * 100}vw`,
        top: `${(872 / 1080) * 100}vh`,
        width: `${(584 / 1920) * 100}vw`,
        height: `${(184 / 1080) * 100}vh`,
      }
    },
    // 顶部 KPI「活跃告警」卡数值居中于卡片中心(画布 x≈1674 y=104)，盖一块径向渐隐底
    kpiStyle() {
      return {
        left: `${(1574 / 1920) * 100}vw`,
        top: `${(102 / 1080) * 100}vh`,
        width: `${(200 / 1920) * 100}vw`,
        height: `${(40 / 1080) * 100}vh`,
      }
    },
    kpiColor() {
      return this.activeCount > 0 ? '#ff6b35' : '#e8f1ff'
    },
  },
  mounted() {
    this.$EventBus && this.$EventBus.$on('GoPage', this.onGoPage)
    this.fetchAlarms()
    this.pollTimer = setInterval(this.fetchAlarms, POLL_INTERVAL)
  },
  beforeDestroy() {
    if (this.pollTimer) clearInterval(this.pollTimer)
    this.$EventBus && this.onGoPage && this.$EventBus.$off('GoPage', this.onGoPage)
  },
  methods: {
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
/* 透明全屏根容器：承载两块定位浮层，自身不挡交互 */
.scada-alarm-root {
  position: absolute;
  inset: 0;
  pointer-events: none;
  z-index: 60;
  font-family: 'Microsoft YaHei', 'PingFang SC', sans-serif;
}

/* 顶部 KPI「活跃告警」实时计数：径向渐隐底融入卡片，无硬边 */
.sa-kpi {
  position: absolute;
  display: flex;
  align-items: center;
  justify-content: center;
  pointer-events: none;
  background: radial-gradient(ellipse at center, #0d1c2d 58%, rgba(13, 28, 45, 0) 100%);
}
.sa-kpi-num {
  font-size: 1.6vw;
  font-weight: 700;
  line-height: 1;
  letter-spacing: 0.5px;
}

.scada-alarm {
  position: absolute;
  pointer-events: auto;
  font-family: 'Microsoft YaHei', 'PingFang SC', sans-serif;
  box-sizing: border-box;
  padding: 4px;
}
/* 内层略微内缩，露出底层 box13 的霓虹边角，背景不透明以遮住静态“无活跃告警”文字 */
.sa-inner {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  background: linear-gradient(180deg, rgba(13, 23, 38, 0.96), rgba(10, 14, 23, 0.96));
  border: 1px solid rgba(30, 58, 95, 0.6);
  border-radius: 4px;
  padding: 8px 10px;
  box-sizing: border-box;
  box-shadow: inset 0 0 20px rgba(0, 50, 90, 0.12);
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
  overflow-y: auto;
  overflow-x: hidden;
}
.sa-row {
  display: flex;
  align-items: center;
  gap: 6px;
  height: 1.85vh;
  padding: 0 4px;
  border-radius: 3px;
  background: rgba(16, 29, 51, 0.7);
  margin-bottom: 3px;
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
  padding-top: 2px;
}

.sa-list::-webkit-scrollbar { width: 5px; }
.sa-list::-webkit-scrollbar-track { background: transparent; }
.sa-list::-webkit-scrollbar-thumb { background: #1e3a5f; border-radius: 3px; }
.sa-list::-webkit-scrollbar-thumb:hover { background: #2c5a8f; }
</style>
