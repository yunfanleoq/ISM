<template>
  <div class="scada-container">
    <!-- 顶部标题栏 + DataV 装饰 -->
    <div class="scada-header">
      <dv-decoration-3 class="header-deco" />
      <div class="header-left">
        <span class="header-logo">⚡</span>
        <span class="header-title">{{ pageTitle }}</span>
        <span class="header-subtitle">{{ pageSubtitle }}</span>
      </div>
      <div class="header-center">
        <div class="breadcrumb">
          <span class="bread-item" :class="{active: drillLevel === 0}" @click="goToLevel('overview')">📊 全局总览</span>
          <template v-if="drillLevel >= 1">
            <span class="bread-arrow">›</span>
            <span class="bread-item active" @click="goToLevel('overview')">{{ selectedBuilding?.name }}</span>
          </template>
          <template v-if="drillLevel >= 2">
            <span class="bread-arrow">›</span>
            <span class="bread-item active">{{ selectedFloor?.name }}</span>
          </template>
          <template v-if="drillLevel >= 3">
            <span class="bread-arrow">›</span>
            <span class="bread-item active">{{ selectedDevice?.name }}</span>
          </template>
        </div>
      </div>
      <div class="header-right">
        <span class="header-clock">{{ currentTime }}</span>
        <a class="header-history-btn" @click="goToHistoryData">📋 历史数据</a>
        <span class="header-status-dot" :class="wsConnected ? 'online' : 'offline'"></span>
        <span>{{ wsConnected ? '系统运行中' : '数据连接中…' }}</span>
      </div>
      <dv-decoration-3 class="header-deco" :reverse="true" />
    </div>

    <!-- 主体内容 -->
    <div class="scada-body">
      <!-- 左侧侧边栏 - DataV border-box-8 包裹 -->
      <div class="scada-sidebar">
        <dv-border-box-8 :dur="3" class="sidebar-wrapper">
          <div class="sidebar-inner">
            <div class="sidebar-panel">
              <div class="panel-title">📍 建筑导航</div>
              <div class="building-tree">
                <div v-for="bldg in buildings" :key="bldg.id" class="tree-bldg">
                  <div class="tree-bldg-header"
                       :class="{selected: selectedBuilding?.id === bldg.id}"
                       @click="selectBuilding(bldg)">
                    <span class="tree-icon">🏢</span>
                    <span class="tree-name">{{ bldg.name }}</span>
                    <span class="tree-count">{{ countBuildingDevices(bldg) }}台</span>
                  </div>
                  <div v-show="selectedBuilding?.id === bldg.id || drillLevel <= 0" class="tree-floors">
                    <div v-for="flr in bldg.floors" :key="flr.id" class="tree-floor"
                         :class="{selected: selectedFloor?.id === flr.id}"
                         @click="selectFloor(bldg, flr)">
                      <span class="tree-icon">{{ flr.name.includes('UPS') ? '🔋' : '📋' }}</span>
                      <span class="tree-name">{{ flr.name }}</span>
                      <span class="tree-badge" :class="getFloorStatus(flr)">{{ getFloorAlarmCount(flr) }}</span>
                    </div>
                  </div>
                </div>
                <div v-if="buildings.length === 0 && !loading" class="tree-empty">
                  ⚠️ 暂无设备数据，请先在设备管理中配置项目设备
                </div>
                <div v-if="loading" class="tree-empty">⏳ 加载设备树中…</div>
              </div>
            </div>

            <!-- 告警快速列表 - DataV border-box-1 包裹 -->
            <div class="sidebar-panel alarm-panel">
              <dv-border-box-1>
                <div class="alarm-inner">
                  <div class="panel-title">🚨 实时告警 <span class="alarm-count">{{ activeAlarms.length }}</span></div>
                  <div class="alarm-list">
                    <div v-for="(alarm, idx) in activeAlarms.slice(0, 8)" :key="idx" class="alarm-item" :class="'level-'+alarm.level">
                      <span class="alarm-time">{{ alarm.time }}</span>
                      <span class="alarm-device">{{ alarm.device }}</span>
                      <span class="alarm-msg">{{ alarm.msg }}</span>
                    </div>
                    <div v-if="activeAlarms.length === 0" class="alarm-empty">✅ 当前无告警</div>
                  </div>
                </div>
              </dv-border-box-1>
            </div>
          </div>
        </dv-border-box-8>
      </div>

      <!-- 中间主区域 -->
      <div class="scada-main">
        <!-- 概览层 -->
        <div v-if="drillLevel === 0" class="level-overview">
          <!-- 第一行：统计卡片 + DataV decoration-5 背景 -->
          <div class="stats-row">
            <div class="stat-card">
              <dv-decoration-5 class="stat-card-bg" />
              <div class="stat-icon" style="background:linear-gradient(135deg,#667eea,#764ba2)">⚡</div>
              <div class="stat-info">
                <div class="stat-value">{{ stats.totalPower }} <small>kW</small></div>
                <div class="stat-label">总功率</div>
              </div>
            </div>
            <div class="stat-card">
              <dv-decoration-5 class="stat-card-bg" />
              <div class="stat-icon" style="background:linear-gradient(135deg,#f093fb,#f5576c)">📊</div>
              <div class="stat-info">
                <div class="stat-value">{{ stats.todayEnergy }} <small>kWh</small></div>
                <div class="stat-label">今日用电量</div>
              </div>
            </div>
            <div class="stat-card">
              <dv-decoration-5 class="stat-card-bg" />
              <div class="stat-icon" style="background:linear-gradient(135deg,#4facfe,#00f2fe)">🖥</div>
              <div class="stat-info">
                <div class="stat-value">{{ stats.onlineDevices }}<small>/{{ stats.totalDevices }}</small></div>
                <div class="stat-label">在线设备</div>
              </div>
            </div>
            <div class="stat-card">
              <dv-decoration-5 class="stat-card-bg" />
              <div class="stat-icon" style="background:linear-gradient(135deg,#fa709a,#fee140)">🔔</div>
              <div class="stat-info">
                <div class="stat-value">{{ activeAlarms.length }}</div>
                <div class="stat-label">活跃告警</div>
              </div>
            </div>
          </div>

          <!-- 第二行：左实时状态 + 右趋势图 -->
          <div class="content-row">
            <!-- 电力单线图概览 -->
            <div class="panel panel-line-diagram">
              <dv-border-box-1>
                <div class="diagram-inner">
                  <div class="panel-title">⚡ 系统拓扑概览</div>
                  <div class="line-diagram">
                    <div v-for="bldg in buildings" :key="bldg.id" class="diagram-bldg" @click="selectBuilding(bldg)">
                      <div class="diagram-bldg-name">🏢 {{ bldg.name }}</div>
                      <div class="diagram-floors">
                        <div v-for="flr in bldg.floors" :key="flr.id" class="diagram-floor" @click.stop="selectFloor(bldg, flr)">
                          <div class="diagram-floor-name">{{ flr.name }}</div>
                          <div class="diagram-device-dots">
                            <span v-for="dev in flr.devices" :key="dev.uid"
                                  class="device-dot" :class="dev.status"
                                  :title="dev.name"></span>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </dv-border-box-1>
            </div>

            <!-- 右侧趋势图 -->
            <div class="panel panel-chart">
              <dv-border-box-1>
                <div class="chart-inner">
                  <div class="panel-title">📈 功率趋势 (24h)</div>
                  <div ref="trendChart" class="chart-container"></div>
                </div>
              </dv-border-box-1>
            </div>
          </div>

          <!-- 第三行：设备状态网格 -->
          <div class="content-row">
            <div class="panel panel-full">
              <dv-border-box-1>
                <div class="grid-inner">
                  <div class="panel-title">🏭 设备运行状态总览</div>
                  <div class="device-grid">
                    <div v-for="bldg in buildings" :key="bldg.id" class="device-grid-section" v-show="bldg.floors.some(f => f.devices.length > 0)">
                      <div class="grid-section-title">🏢 {{ bldg.name }}</div>
                      <div class="grid-section-devices">
                        <template v-for="flr in bldg.floors">
                          <div v-for="dev in flr.devices" :key="dev.uid"
                               class="device-card" :class="dev.status"
                               @click="selectDevice(bldg, flr, dev)">
                            <div class="device-card-icon">
                              <ElectricIcon :type="dev.type" :status="dev.status" :size="36" />
                            </div>
                            <div class="device-card-name">{{ dev.name }}</div>
                            <div class="device-card-status" :class="dev.status">
                              {{ statusLabel(dev.status) }}
                            </div>
                          </div>
                        </template>
                      </div>
                    </div>
                  </div>
                </div>
              </dv-border-box-1>
            </div>
          </div>
        </div>

        <!-- 建筑层级 -->
        <div v-else-if="drillLevel === 1" class="level-building">
          <div class="level-header">
            <span class="back-btn" @click="goToLevel('overview')">← 返回总览</span>
            <h2>🏢 {{ selectedBuilding?.name }}</h2>
            <span class="level-subtitle">{{ countBuildingDevices(selectedBuilding) }}台设备 · {{ countBuildingAlarms(selectedBuilding) }}条告警</span>
          </div>
          <div class="floor-cards">
            <div v-for="flr in selectedBuilding?.floors || []" :key="flr.id"
                 class="floor-card" @click="selectFloor(selectedBuilding, flr)">
              <dv-border-box-1>
                <div class="floor-card-inner">
                  <div class="floor-card-header">
                    <span class="floor-card-name">{{ flr.name.includes('UPS') ? '🔋' : '📋' }} {{ flr.name }}</span>
                    <span class="floor-card-count">{{ flr.devices.length }}台设备</span>
                  </div>
                  <div class="floor-card-stats">
                    <span class="floor-stat running">🟢 {{ countDevicesByStatus(flr.devices, 'running') }}运行</span>
                    <span class="floor-stat alarm">🔴 {{ countDevicesByStatus(flr.devices, 'alarm') }}告警</span>
                    <span class="floor-stat stopped">⏸ {{ countDevicesByStatus(flr.devices, 'stopped') }}停止</span>
                  </div>
                </div>
              </dv-border-box-1>
            </div>
          </div>
        </div>

        <!-- 楼层层级 -->
        <div v-else-if="drillLevel === 2" class="level-floor">
          <div class="level-header">
            <span class="back-btn" @click="selectBuilding(selectedBuilding)">← {{ selectedBuilding?.name }}</span>
            <h2>📋 {{ selectedFloor?.name }}</h2>
            <span class="level-subtitle">{{ selectedFloor?.devices?.length || 0 }}台设备</span>
          </div>
          <div class="device-table">
            <dv-border-box-1>
              <div class="table-inner">
                <table>
                  <thead>
                    <tr>
                      <th>设备名称</th>
                      <th>图标</th>
                      <th>协议</th>
                      <th>状态</th>
                      <th>实时功率</th>
                      <th>电流A/B/C</th>
                      <th>温度</th>
                      <th>操作</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="dev in selectedFloor?.devices || []" :key="dev.uid" :class="'row-'+dev.status">
                      <td class="dev-name">{{ dev.name }}</td>
                      <td><ElectricIcon :type="dev.type" :status="dev.status" :size="28" /></td>
                      <td><span class="protocol-tag">{{ dev.protocol?.toUpperCase() }}</span></td>
                      <td><span class="status-tag" :class="dev.status">{{ statusLabel(dev.status) }}</span></td>
                      <td>{{ getRealValue(dev.uid, 'power') }} kW</td>
                      <td>{{ getRealValue(dev.uid, 'Ia') }}/{{ getRealValue(dev.uid, 'Ib') }}/{{ getRealValue(dev.uid, 'Ic') }} A</td>
                      <td>{{ getRealValue(dev.uid, 'temp') }}°C</td>
                      <td><button class="detail-btn" @click="selectDevice(selectedBuilding, selectedFloor, dev)">详情</button></td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </dv-border-box-1>
          </div>
        </div>

        <!-- 设备详情层级 -->
        <div v-else-if="drillLevel === 3" class="level-device">
          <div class="level-header">
            <span class="back-btn" @click="selectFloor(selectedBuilding, selectedFloor)">← {{ selectedFloor?.name }}</span>
            <h2>{{ selectedDevice?.type === 'ups' ? '🔋' : '🔧' }} {{ selectedDevice?.name }}</h2>
            <span class="status-tag" :class="selectedDevice?.status">{{ statusLabel(selectedDevice?.status) }}</span>
          </div>
          <div class="device-detail-grid">
            <!-- 设备 SVG 图形预览 -->
            <div class="detail-panel detail-svg">
              <div class="panel-title">🔌 设备拓扑</div>
              <div class="svg-preview">
                <ElectricIcon :type="selectedDevice?.type" :status="selectedDevice?.status" :size="120" />
              </div>
            </div>
            <div class="detail-panel">
              <div class="panel-title">📋 基本参数</div>
              <table class="param-table">
                <tr><td>设备名称</td><td>{{ selectedDevice?.name }}</td></tr>
                <tr><td>设备类型</td><td>{{ selectedDevice?.type }}</td></tr>
                <tr><td>通信协议</td><td>{{ selectedDevice?.protocol?.toUpperCase() }}</td></tr>
                <tr><td>运行状态</td><td><span class="status-tag" :class="selectedDevice?.status">{{ statusLabel(selectedDevice?.status) }}</span></td></tr>
                <tr><td>Modbus地址</td><td>{{ selectedDevice?.modbusAddr || '-' }}</td></tr>
                <tr><td>打包时间(ms)</td><td>{{ selectedDevice?.packTime || '-' }}</td></tr>
                <tr><td>所属建筑</td><td>{{ selectedBuilding?.name }}</td></tr>
                <tr><td>所属楼层</td><td>{{ selectedFloor?.name }}</td></tr>
              </table>
            </div>
            <div class="detail-panel">
              <div class="panel-title">📊 实时参数</div>
              <table class="param-table">
                <template v-if="selectedDevice?.type === 'ups'">
                  <tr><td>输出A相电压</td><td>{{ getRealValue(selectedDevice?.uid, 'Ua') }} V</td></tr>
                  <tr><td>输出B相电压</td><td>{{ getRealValue(selectedDevice?.uid, 'Ub') }} V</td></tr>
                  <tr><td>输出C相电压</td><td>{{ getRealValue(selectedDevice?.uid, 'Uc') }} V</td></tr>
                  <tr><td>输出A相电流</td><td>{{ getRealValue(selectedDevice?.uid, 'Ia') }} A</td></tr>
                  <tr><td>输出总有功功率</td><td>{{ getRealValue(selectedDevice?.uid, 'power') }} kW</td></tr>
                  <tr><td>输出功率因数</td><td>{{ formatPf(selectedDevice?.uid) }}</td></tr>
                </template>
                <template v-else>
                  <tr><td>A相电压</td><td>{{ getRealValue(selectedDevice?.uid, 'Ua') }} V</td></tr>
                  <tr><td>B相电压</td><td>{{ getRealValue(selectedDevice?.uid, 'Ub') }} V</td></tr>
                  <tr><td>C相电压</td><td>{{ getRealValue(selectedDevice?.uid, 'Uc') }} V</td></tr>
                  <tr><td>A相电流</td><td>{{ getRealValue(selectedDevice?.uid, 'Ia') }} A</td></tr>
                  <tr><td>有功功率</td><td>{{ getRealValue(selectedDevice?.uid, 'power') }} kW</td></tr>
                  <tr><td>功率因数</td><td>{{ formatPf(selectedDevice?.uid) }}</td></tr>
                </template>
              </table>
            </div>
            <div class="detail-panel detail-chart">
              <div class="panel-title">📈 24小时功率曲线</div>
              <div ref="deviceChart" class="chart-container"></div>
            </div>
            <div class="detail-panel">
              <div class="panel-title">🔔 设备告警</div>
              <div class="alarm-history" v-if="deviceAlarms.length > 0">
                <div v-for="(alarm, idx) in deviceAlarms" :key="idx" class="alarm-history-item">
                  <span class="time">{{ alarm.time }}</span>
                  <span class="level" :class="'level-'+alarm.level">{{ levelLabel(alarm.level) }}</span>
                  <span>{{ alarm.msg }}</span>
                </div>
              </div>
              <div v-else class="alarm-empty">✅ 该设备无告警记录</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import * as echarts from 'echarts'
import { getMonitorTree, GetDeviceModelDataList, getRealDataByUuid } from '@/services/device'
import { GetChartDataHistoryTrendList } from '@/services/report'
import {
  DEVICE_TYPE_COLORS,
  DEVICE_TYPE_COMPONENT,
  STATUS_LABELS,
  DATA_KEY_MAP,
  floorKeyForDevice,
  inferDeviceType,
} from './config'
import {
  buildPointCatalog,
  findModelDataUuid,
  ingestRealDataResponse,
  ingestWsPush,
  readCachedValue,
  realDataKey,
} from './scadaRealData'
import ElectricIcon from './ElectricIcon.vue'

export default {
  name: 'SCADAMonitor',
  components: { ElectricIcon },
  props: {
    title: {
      type: String,
      default: '配电房智能监控系统',
    },
    subtitle: {
      type: String,
      default: 'SCADA Power Monitor',
    },
    projectUuid: {
      type: String,
      default: '31bc90be-ebc4-dd61-ba9d-ce6e075e40e2',
    },
  },
  data() {
    return {
      buildings: [],
      drillLevel: 0,
      selectedBuilding: null,
      selectedFloor: null,
      selectedDevice: null,
      currentTime: '',
      loading: false,
      wsConnected: false,
      stats: {
        totalPower: '0',
        todayEnergy: '0',
        onlineDevices: 0,
        totalDevices: 0,
      },
      activeAlarms: [],
      realTimeData: {},
      pendingRealData: {},
      pointCatalog: {},
      subscribeModelUuids: [],
      subscribeDeviceUuids: [],
      refreshIntervalMs: 1000,
      dataPollTimer: null,
      flushTimer: null,
      trendChart: null,
      deviceChart: null,
      powerHistory: new Array(24).fill(null),
      devicePowerHistory: new Array(24).fill(null),
      lastRecordedHour: -1,
      historyRecordingTimer: null,
    }
  },
  computed: {
    pageTitle() {
      return this.title
    },
    pageSubtitle() {
      return this.subtitle
    },
    deviceAlarms() {
      if (!this.selectedDevice) return []
      return this.activeAlarms.filter(a => a.deviceUid === this.selectedDevice?.uid)
    },
  },
  mounted() {
    this.updateClock()
    this.clockTimer = setInterval(this.updateClock, 1000)
    this.fetchDeviceTree()
    this.setupEventBus()
    this.startDataTimers()
    this.$nextTick(() => {
      this.initTrendChart()
    })
  },
  beforeDestroy() {
    clearInterval(this.clockTimer)
    clearInterval(this.historyRecordingTimer)
    this.stopDataTimers()
    this.trendChart?.dispose()
    this.deviceChart?.dispose()
    if (this._readDataHandler) {
      this.$EventBus.$off('readDataPush', this._readDataHandler)
    }
    if (this._realAlarmHandler) {
      this.$EventBus.$off('RealAlarm', this._realAlarmHandler)
    }
  },
  methods: {
    // ========== 数据获取 ==========
    async fetchDeviceTree() {
      this.loading = true
      try {
        const res = await getMonitorTree({ headers: { ProjectUuid: this.projectUuid } })
        if (res.data && res.data.code === 0 && res.data.list && res.data.list.length > 0) {
          this.buildings = this.transformTreeToBuildings(res.data.list)
        }
      } catch (e) {
        console.warn('[SCADA] 获取设备树失败:', e.message)
      }
      this.loading = false
      if (this.buildings.length === 0) {
        console.warn('[SCADA] 当前项目无设备数据，请先在设备管理中配置')
      } else {
        this.recalcStats()
        await this.fetchDeviceDataPoints()
      }
      this.$nextTick(() => {
        if (!this.trendChart) this.initTrendChart()
      })
    },

    /** 拉取数据点目录并订阅实时数据 */
    async fetchDeviceDataPoints() {
      const allDeviceUids = []
      for (const bldg of this.buildings) {
        for (const flr of bldg.floors) {
          for (const dev of flr.devices) {
            if (dev.uid) allDeviceUids.push(dev.uid)
          }
        }
      }
      if (allDeviceUids.length === 0) {
        console.warn('[SCADA] fetchDeviceDataPoints: 没有设备 UID，跳过')
        return
      }

      try {
        const res = await GetDeviceModelDataList({ getType: 0, SelectDevice: allDeviceUids })
        if (res.data?.code === 0 && res.data.list) {
          const { catalog, modelDataUuids, deviceUuids } = buildPointCatalog(this.buildings, res.data.list)
          this.pointCatalog = catalog
          this.subscribeModelUuids = modelDataUuids
          this.subscribeDeviceUuids = deviceUuids
          console.log('[SCADA] 点位目录:', Object.keys(catalog).length, '订阅设备:', deviceUuids.length)
          await this.pullRealDataSnapshot()
          this.fetchTrendHistoryData()
        } else {
          console.warn('[SCADA] GetDeviceModelDataList 返回异常 code:', res.data?.code)
        }
      } catch (e) {
        console.warn('[SCADA] 获取数据点映射失败:', e.message)
      }
    },

    /** 调用 getRealDataByUuid 拉取当前快照（同时触发后端 WS 推送注册） */
    async pullRealDataSnapshot() {
      if (!this.subscribeModelUuids.length || !this.subscribeDeviceUuids.length) return
      try {
        const res = await getRealDataByUuid({
          uuid: this.subscribeModelUuids,
          devices: this.subscribeDeviceUuids,
        })
        if (res.data?.code === 0 && Array.isArray(res.data.realData)) {
          ingestRealDataResponse(this.pendingRealData, res.data.realData)
          this.flushRealDataToView()
          console.log('[SCADA] 实时快照:', res.data.realData.filter(r => r.value !== '').length, '/', res.data.realData.length)
        }
      } catch (e) {
        console.warn('[SCADA] pullRealDataSnapshot 失败:', e.message)
      }
    },

    startDataTimers() {
      this.stopDataTimers()
      // UI 合并 pending → realTimeData，默认 1s 刷新一次（避免 WS 高频推送导致界面抖动）
      this.flushTimer = setInterval(() => this.flushRealDataToView(), this.refreshIntervalMs)
      // 轮询快照作为 WS 兜底，间隔略长以减轻后端压力
      this.dataPollTimer = setInterval(() => this.pullRealDataSnapshot(), this.refreshIntervalMs * 5)
    },

    stopDataTimers() {
      if (this.dataPollTimer) {
        clearInterval(this.dataPollTimer)
        this.dataPollTimer = null
      }
      if (this.flushTimer) {
        clearInterval(this.flushTimer)
        this.flushTimer = null
      }
    },

    /** 将 pending 缓冲合并到 realTimeData 并刷新 UI（默认 1s 一次） */
    flushRealDataToView() {
      const keys = Object.keys(this.pendingRealData)
      if (!keys.length) return
      for (const key of keys) {
        this.$set(this.realTimeData, key, this.pendingRealData[key])
      }
      this.recalcStats()
      this.recordPowerHistory()
      this.refreshTrendChart()
      if (this.drillLevel === 3) {
        this.refreshDeviceChart()
      }
    },

    transformTreeToBuildings(treeList) {
      const buildings = []
      // RootZone (Sid=1) children → buildings (type=0)
      function findBuildings(nodes) {
        for (const node of nodes) {
          if (node.value && node.value.Sid === 1) {
            return findBuildings(node.children || [])
          }
          if (node.value && node.value.type === 0) {
            const bldg = buildBuilding(node)
            if (bldg && bldg.floors.length > 0) {
              buildings.push(bldg)
            }
          }
          if (node.children && node.children.length > 0) {
            findBuildings(node.children)
          }
        }
      }
      function buildBuilding(node) {
        const floors = []
        const directDevices = []
        if (node.children) {
          for (const child of node.children) {
            if (child.value && child.value.type === 0) {
              const flr = buildFloor(child)
              if (flr && flr.devices.length > 0) {
                floors.push(flr)
              }
            } else if (child.value && child.value.type === 1) {
              directDevices.push(deviceFromNode(child))
            }
          }
        }
        if (directDevices.length > 0) {
          const grouped = {}
          for (const dev of directDevices) {
            const key = floorKeyForDevice(dev.name, dev.muid)
            if (!grouped[key]) {
              grouped[key] = { id: `${node.key}_${key}`, name: key, devices: [] }
            }
            grouped[key].devices.push(dev)
          }
          floors.push(...Object.values(grouped).sort((a, b) => a.name.localeCompare(b.name, 'zh-CN')))
        }
        return {
          id: node.key,
          name: node.text || node.value?.Name || '未命名',
          floors,
        }
      }
      function buildFloor(node) {
        const devices = []
        function collectDevices(children) {
          if (!children) return
          for (const child of children) {
            if (child.value && child.value.type === 1) {
              devices.push(deviceFromNode(child))
            }
            if (child.children) {
              collectDevices(child.children)
            }
          }
        }
        collectDevices(node.children || [])
        return {
          id: node.key,
          name: node.text || node.value?.Name || '未命名',
          devices,
        }
      }
      function deviceFromNode(node) {
        const v = node.value || {}
        let extraData = {}
        if (v.extra) {
          try {
            extraData = typeof v.extra === 'string' ? JSON.parse(v.extra) : v.extra
          } catch (e) {
            extraData = {}
          }
        }
        const modbusAddr = extraData?.Modbus?.address || extraData?.modbus?.address || v.slaveId || ''
        const packTime = extraData?.Modbus?.packTime ?? extraData?.modbus?.packTime ?? ''
        const deviceName = node.text || v.Name || '未知设备'
        const muid = v.muid || ''
        return {
          uid: node.key,
          name: deviceName,
          type: inferDeviceType(deviceName, muid, v.deviceType),
          protocol: v.protocol || (deviceName.startsWith('UPS') ? 'modbus' : 'modbus'),
          status: v.Status === 1 ? 'running' : v.Status === 0 ? 'stopped' : 'offline',
          muid,
          configUid: v.configUid || '',
          modbusAddr: modbusAddr,
          packTime: packTime,
          rawValue: v,
        }
      }
      findBuildings(treeList || [])
      return buildings
    },

    /** 从 realTimeData 聚合真实统计值 */
    recalcStats() {
      let totalDevices = 0
      let onlineDevices = 0
      for (const bldg of this.buildings) {
        for (const flr of bldg.floors) {
          totalDevices += flr.devices.length
          onlineDevices += flr.devices.filter(d => d.status === 'running').length
        }
      }
      this.stats.totalDevices = totalDevices
      this.stats.onlineDevices = onlineDevices

      // 从 WebSocket / 轮询快照聚合总功率和今日用电量
      let totalPower = 0
      let todayEnergy = 0
      const powerCandidates = DATA_KEY_MAP.power || []
      const energyCandidates = (DATA_KEY_MAP.energy || []).concat(DATA_KEY_MAP.todayEnergy || [])

      for (const entry of Object.values(this.pointCatalog)) {
        const data = this.realTimeData[realDataKey(entry.deviceUid, entry.modelDataUuid)]
        if (!data || data.Value === undefined || data.Value === null || data.Value === '') continue
        const v = Number(data.Value)
        if (Number.isNaN(v)) continue

        if (powerCandidates.some(name => entry.dataName === name || entry.dataName.includes(name))) {
          totalPower += v
        }
        if (energyCandidates.some(name => entry.dataName === name || entry.dataName.includes(name))) {
          todayEnergy += v
        }
      }

      this.stats.totalPower = totalPower > 0 ? String(Math.round(totalPower * 10) / 10) : '0'
      this.stats.todayEnergy = todayEnergy > 0 ? String(Math.round(todayEnergy * 10) / 10) : '0'
    },

    // ========== EventBus 实时数据 ==========
    setupEventBus() {
      this._readDataHandler = (wsData) => {
        this.wsConnected = true
        ingestWsPush(this.pendingRealData, wsData)
      }
      this._realAlarmHandler = (wsData) => {
        this.wsConnected = true
        if (wsData.Data && Array.isArray(wsData.Data)) {
          for (const alarm of wsData.Data) {
            const exists = this.activeAlarms.find(a => a.rawId === alarm.ID)
            if (!exists) {
              this.activeAlarms.unshift({
                rawId: alarm.ID,
                time: alarm.HappenTime || new Date().toLocaleTimeString('zh-CN', { hour12: false }),
                device: alarm.DeviceName || alarm.DeviceUuid || '',
                deviceUid: alarm.DeviceUuid || '',
                msg: alarm.AlarmMsg || alarm.DataName || '',
                level: alarm.AlarmLevel || 0,
              })
            }
          }
          if (this.activeAlarms.length > 100) {
            this.activeAlarms = this.activeAlarms.slice(0, 50)
          }
        }
      }
      this.$EventBus.$on('readDataPush', this._readDataHandler)
      this.$EventBus.$on('RealAlarm', this._realAlarmHandler)
    },

    // ========== 实时数据取值 ==========
    getRealValue(deviceUid, key) {
      if (!deviceUid) return '--'
      const v = readCachedValue(this.realTimeData, deviceUid, key, this.pointCatalog)
      if (v === null) return '--'
      if (key === 'pf') {
        const pf = v > 1 ? v / 100 : v
        return pf.toFixed(2)
      }
      return Number.isInteger(v) ? String(v) : (Math.round(v * 100) / 100).toString()
    },

    findDataPointUuid(deviceUid, key) {
      return findModelDataUuid(this.pointCatalog, deviceUid, key)
    },

    // ========== 时钟 ==========
    updateClock() {
      const now = new Date()
      this.currentTime = now.toLocaleString('zh-CN', { hour12: false })
    },

    // ========== 图表 ==========
    initTrendChart() {
      if (this.trendChart) return
      const el = this.$refs.trendChart
      if (!el) return
      this.trendChart = echarts.init(el)
      const hours = Array.from({ length: 24 }, (_, i) => i + ':00')
      const data = [...this.powerHistory]
      this.trendChart.setOption({
        grid: { top: 10, right: 20, bottom: 30, left: 50 },
        xAxis: { type: 'category', data: hours, axisLabel: { color: '#94a3b8', fontSize: 10, interval: 3 } },
        yAxis: { type: 'value', name: 'kW', axisLabel: { color: '#94a3b8' }, splitLine: { lineStyle: { color: '#1e293b' } } },
        series: [{
          data, type: 'line', smooth: true, connectNulls: false,
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(52,152,219,0.3)' },
              { offset: 1, color: 'rgba(52,152,219,0.02)' }
            ])
          },
          lineStyle: { color: '#3498db', width: 2 },
          itemStyle: { color: '#3498db' },
          symbol: 'none',
        }]
      })
    },
    /** 刷新趋势图数据 */
    refreshTrendChart() {
      if (!this.trendChart) return
      this.trendChart.setOption({ series: [{ data: [...this.powerHistory] }] })
    },
    initDeviceChart() {
      this.$nextTick(() => {
        if (this.deviceChart) this.deviceChart.dispose()
        const el = this.$refs.deviceChart
        if (!el) return
        this.deviceChart = echarts.init(el)
        const hours = Array.from({ length: 24 }, (_, i) => i + ':00')
        const data = [...this.devicePowerHistory]
        this.deviceChart.setOption({
          grid: { top: 10, right: 20, bottom: 30, left: 45 },
          xAxis: { type: 'category', data: hours, axisLabel: { color: '#94a3b8', fontSize: 9, interval: 4 } },
          yAxis: { type: 'value', name: 'kW', axisLabel: { color: '#94a3b8' }, splitLine: { lineStyle: { color: '#1e293b' } } },
          series: [{
            data, type: 'bar',
            itemStyle: {
              color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                { offset: 0, color: '#3498db' },
                { offset: 1, color: '#2c3e50' }
              ])
            }
          }],
          tooltip: { trigger: 'axis' },
        })
      })
    },
    /** 刷新设备曲线图 */
    refreshDeviceChart() {
      if (!this.deviceChart) return
      this.deviceChart.setOption({ series: [{ data: [...this.devicePowerHistory] }] })
    },

    /** 每小时记录一次当前总功率到 powerHistory */
    recordPowerHistory() {
      const now = new Date()
      const currentHour = now.getHours()
      if (currentHour === this.lastRecordedHour) return
      this.lastRecordedHour = currentHour

      const powerCandidates = DATA_KEY_MAP.power || []
      let totalPower = 0
      for (const entry of Object.values(this.pointCatalog)) {
        if (!powerCandidates.some(name => entry.dataName === name || entry.dataName.includes(name))) continue
        const data = this.realTimeData[realDataKey(entry.deviceUid, entry.modelDataUuid)]
        if (!data || data.Value === undefined || data.Value === null || data.Value === '') continue
        const v = Number(data.Value)
        if (Number.isNaN(v)) continue
        totalPower += v
      }

      const newHistory = [...this.powerHistory]
      newHistory[currentHour] = Math.round(totalPower)
      this.powerHistory = newHistory

      if (this.selectedDevice?.uid) {
        const devicePower = readCachedValue(this.realTimeData, this.selectedDevice.uid, 'power', this.pointCatalog)
        if (devicePower !== null) {
          const deviceHistory = [...this.devicePowerHistory]
          deviceHistory[currentHour] = Math.round(devicePower)
          this.devicePowerHistory = deviceHistory
        }
      }

      this.refreshTrendChart()
      this.refreshDeviceChart()
    },

    /** 格式化功率因数显示（处理 0-100 与 0-1 两种量程） */
    formatPf(deviceUid) {
      return this.getRealValue(deviceUid, 'pf')
    },

    /** 从后端历史数据库拉取最近 24h 的功率趋势数据 */
    async fetchTrendHistoryData() {
      const powerCandidates = DATA_KEY_MAP.power || []
      const deviceList = []
      for (const entry of Object.values(this.pointCatalog)) {
        if (powerCandidates.some(name => entry.dataName === name || entry.dataName.includes(name))) {
          deviceList.push({
            DeviceUuid: entry.deviceUid,
            ModelDataUuid: entry.modelDataUuid,
          })
        }
      }
      if (deviceList.length === 0) return

      try {
        const res = await GetChartDataHistoryTrendList({
          TimeIn: 1,
          List: deviceList,
          HistoryTime: 24 * 60, // 最近 24 小时（1440 分钟）
        })
        if (res.data && res.data.code === 0 && res.data.data) {
          const rawList = res.data.data
          this.buildPowerHistoryFromRaw(rawList)
        }
      } catch (e) {
        console.warn('[SCADA] 获取历史趋势数据失败:', e.message)
      }
      // 启动定时刷新（每 5 分钟从 DB 同步最新历史数据）
      if (this.historyRecordingTimer) clearInterval(this.historyRecordingTimer)
      this.historyRecordingTimer = setInterval(() => {
        this.fetchTrendHistoryData()
      }, 5 * 60 * 1000)
    },

    /** 将原始历史数据按小时桶聚合到 powerHistory */
    buildPowerHistoryFromRaw(rawList) {
      // 按小时分桶：hourBucket[0..23] = { sum, count }
      const hourBuckets = new Array(24).fill(null).map(() => ({ sum: 0, count: 0 }))
      const now = new Date()
      const currentHour = now.getHours()

      for (const item of rawList) {
        if (!item.RecordTime || !item.DataValue) continue
        const t = new Date(item.RecordTime)
        const hour = t.getHours()
        const v = Number(item.DataValue)
        if (isNaN(v)) continue
        // 只统计当前小时及之前 23 小时内的数据
        const hourDiff = currentHour - hour
        const targetIdx = hourDiff < 0 ? hourDiff + 24 : hourDiff
        const idx = 23 - Math.min(targetIdx, 23)
        if (idx >= 0 && idx < 24) {
          hourBuckets[idx].sum += v
          hourBuckets[idx].count += 1
        }
      }

      const newHistory = [...this.powerHistory]
      for (let i = 0; i < 24; i++) {
        if (hourBuckets[i].count > 0) {
          newHistory[i] = Math.round(hourBuckets[i].sum / hourBuckets[i].count)
        }
      }
      this.powerHistory = newHistory
      this.refreshTrendChart()
      this.refreshDeviceChart()
    },

    // ========== 层级导航 ==========
    goToHistoryData() {
      this.$router.push('/Reporting/DataHistoryReport')
    },
    goToLevel(level) {
      if (level === 'overview') {
        this.drillLevel = 0
        this.selectedBuilding = null
        this.selectedFloor = null
        this.selectedDevice = null
        this.deviceChart?.dispose()
        this.deviceChart = null
      }
    },
    selectBuilding(bldg) {
      this.selectedBuilding = bldg
      this.selectedFloor = null
      this.selectedDevice = null
      this.drillLevel = 1
    },
    selectFloor(bldg, flr) {
      this.selectedBuilding = bldg
      this.selectedFloor = flr
      this.selectedDevice = null
      this.drillLevel = 2
    },
    selectDevice(bldg, flr, dev) {
      this.selectedBuilding = bldg
      this.selectedFloor = flr
      this.selectedDevice = dev
      this.drillLevel = 3
      this.devicePowerHistory = new Array(24).fill(null)
      this.fetchDeviceTrendHistory(dev.uid)
      this.initDeviceChart()
    },

    /** 拉取单台设备 24h 功率历史 */
    async fetchDeviceTrendHistory(deviceUid) {
      const modelDataUuid = findModelDataUuid(this.pointCatalog, deviceUid, 'power')
      if (!modelDataUuid) return
      try {
        const res = await GetChartDataHistoryTrendList({
          TimeIn: 1,
          List: [{ DeviceUuid: deviceUid, ModelDataUuid: modelDataUuid }],
          HistoryTime: 24 * 60,
        })
        if (res.data?.code === 0 && res.data.data) {
          this.buildDevicePowerHistoryFromRaw(res.data.data)
        }
      } catch (e) {
        console.warn('[SCADA] 获取设备历史趋势失败:', e.message)
      }
    },

    buildDevicePowerHistoryFromRaw(rawList) {
      const hourBuckets = new Array(24).fill(null).map(() => ({ sum: 0, count: 0 }))
      const now = new Date()
      const currentHour = now.getHours()
      for (const item of rawList) {
        if (!item.RecordTime || item.DataValue === undefined || item.DataValue === null || item.DataValue === '') continue
        const t = new Date(item.RecordTime)
        const hour = t.getHours()
        const v = Number(item.DataValue)
        if (Number.isNaN(v)) continue
        const hourDiff = currentHour - hour
        const targetIdx = hourDiff < 0 ? hourDiff + 24 : hourDiff
        const idx = 23 - Math.min(targetIdx, 23)
        if (idx >= 0 && idx < 24) {
          hourBuckets[idx].sum += v
          hourBuckets[idx].count += 1
        }
      }
      const newHistory = [...this.devicePowerHistory]
      for (let i = 0; i < 24; i++) {
        if (hourBuckets[i].count > 0) {
          newHistory[i] = Math.round(hourBuckets[i].sum / hourBuckets[i].count)
        }
      }
      this.devicePowerHistory = newHistory
      this.refreshDeviceChart()
    },

    // ========== 统计函数 ==========
    countBuildingDevices(bldg) {
      if (!bldg) return 0
      return bldg.floors.reduce((sum, f) => sum + f.devices.length, 0)
    },
    countBuildingAlarms(bldg) {
      if (!bldg) return 0
      let count = 0
      for (const flr of bldg.floors) {
        count += flr.devices.filter(d => d.status === 'alarm').length
      }
      return count
    },
    getFloorStatus(flr) {
      const hasAlarm = flr.devices.some(d => d.status === 'alarm')
      return hasAlarm ? 'alarm' : 'normal'
    },
    getFloorAlarmCount(flr) {
      if (!flr) return 0
      return flr.devices.filter(d => d.status === 'alarm').length
    },
    countDevicesByStatus(devices, status) {
      return devices.filter(d => d.status === status).length
    },
    statusLabel(status) {
      return STATUS_LABELS[status] || status
    },
    levelLabel(level) {
      const map = { 0: '提示', 1: '次要', 2: '重要', 3: '紧急' }
      return map[level] || '未知'
    },
  },
}

</script>

<style scoped>
.scada-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #0a0e17;
  color: #e2e8f0;
  font-family: 'Microsoft YaHei', 'PingFang SC', sans-serif;
  overflow: hidden;
}

/* ====== Header ====== */
.scada-header {
  height: 56px;
  background: linear-gradient(180deg, #141c2b 0%, #0d1220 100%);
  border-bottom: 2px solid #1e3a5f;
  display: flex;
  align-items: center;
  padding: 0 20px;
  flex-shrink: 0;
  z-index: 10;
  position: relative;
}
.header-deco {
  position: absolute;
  top: 0;
  width: 200px;
  height: 56px;
  pointer-events: none;
}
.header-deco:first-child { left: 0; }
.header-deco:last-child { right: 0; }
.header-left { display: flex; align-items: center; gap: 10px; margin-left: 210px; }
.header-logo { font-size: 28px; }
.header-title { font-size: 20px; font-weight: 700; background: linear-gradient(90deg, #60a5fa, #34d399); -webkit-background-clip: text; -webkit-text-fill-color: transparent; }
.header-subtitle { font-size: 12px; color: #64748b; margin-left: 8px; }
.header-center { flex: 1; display: flex; justify-content: center; }
.header-right { display: flex; align-items: center; gap: 8px; font-size: 13px; color: #94a3b8; }
.header-clock { font-family: 'Courier New', monospace; color: #60a5fa; }
.header-history-btn {
  cursor: pointer; padding: 3px 12px; border-radius: 4px;
  background: rgba(96,165,250,.12); color: #60a5fa;
  border: 1px solid rgba(96,165,250,.25);
  font-size: 12px; transition: all .2s; text-decoration: none;
}
.header-history-btn:hover { background: rgba(96,165,250,.25); border-color: #60a5fa; color: #93c5fd; }
.header-status-dot { width: 8px; height: 8px; border-radius: 50%; }
.header-status-dot.online { background: #22c55e; box-shadow: 0 0 8px #22c55e; }
.header-status-dot.offline { background: #f59e0b; box-shadow: 0 0 8px #f59e0b; animation: pulse 1.5s infinite; }
@keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.3; } }

/* Breadcrumb */
.breadcrumb { display: flex; align-items: center; gap: 4px; font-size: 14px; }
.bread-item { cursor: pointer; color: #64748b; padding: 4px 8px; border-radius: 4px; transition: all .2s; }
.bread-item:hover { color: #60a5fa; background: rgba(96,165,250,.1); }
.bread-item.active { color: #e2e8f0; font-weight: 600; }
.bread-arrow { color: #334155; margin: 0 2px; }

/* ====== Body ====== */
.scada-body { flex: 1; display: flex; overflow: hidden; }

/* ====== Sidebar ====== */
.scada-sidebar {
  width: 280px;
  background: #0d1220;
  border-right: 1px solid #1e293b;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}
.sidebar-wrapper { flex: 1; display: flex; flex-direction: column; }
.sidebar-inner { flex: 1; overflow-y: auto; padding: 4px; }
.sidebar-panel { padding: 8px 4px; border-bottom: 1px solid #1e293b; margin-bottom: 4px; }
.sidebar-panel.alarm-panel { flex: 1; overflow-y: auto; border-bottom: none; }
.alarm-inner { padding: 8px; }
.panel-title {
  font-size: 13px; font-weight: 600; color: #94a3b8;
  text-transform: uppercase; letter-spacing: 1px; margin-bottom: 10px;
  display: flex; align-items: center; gap: 6px;
}

/* Building Tree */
.tree-empty { text-align: center; color: #475569; padding: 20px 10px; font-size: 12px; }
.tree-bldg { margin-bottom: 2px; }
.tree-bldg-header {
  display: flex; align-items: center; gap: 6px; padding: 8px 10px;
  border-radius: 6px; cursor: pointer; transition: all .2s; font-size: 13px;
}
.tree-bldg-header:hover, .tree-bldg-header.selected { background: rgba(96,165,250,.12); color: #60a5fa; }
.tree-icon { font-size: 14px; }
.tree-name { flex: 1; }
.tree-count { font-size: 11px; color: #475569; }
.tree-floors { margin-left: 10px; border-left: 1px solid #1e293b; padding-left: 10px; }
.tree-floor {
  display: flex; align-items: center; gap: 6px; padding: 6px 8px;
  border-radius: 4px; cursor: pointer; margin: 2px 0; font-size: 12px; transition: all .2s;
}
.tree-floor:hover, .tree-floor.selected { background: rgba(96,165,250,.08); color: #93c5fd; }
.tree-badge { font-size: 10px; padding: 1px 6px; border-radius: 10px; background: #1e293b; }
.tree-badge.alarm { background: rgba(239,68,68,.2); color: #ef4444; }

/* Alarm List */
.alarm-count { background: #ef4444; color: #fff; font-size: 11px; padding: 1px 8px; border-radius: 10px; font-weight: 700; }
.alarm-list { display: flex; flex-direction: column; gap: 4px; max-height: 300px; overflow-y: auto; }
.alarm-item {
  display: flex; flex-direction: column; gap: 2px; padding: 6px 8px;
  border-radius: 6px; background: #0f172a; border-left: 3px solid #334155; font-size: 11px;
}
.alarm-item.level-3 { border-left-color: #dc2626; background: rgba(220,38,38,.06); }
.alarm-item.level-2 { border-left-color: #ef4444; background: rgba(239,68,68,.06); }
.alarm-item.level-1 { border-left-color: #f59e0b; background: rgba(245,158,11,.06); }
.alarm-item.level-0 { border-left-color: #3b82f6; background: rgba(59,130,246,.06); }
.alarm-time { color: #64748b; font-size: 10px; }
.alarm-device { color: #94a3b8; font-weight: 600; }
.alarm-msg { color: #cbd5e1; }
.alarm-empty { text-align: center; color: #475569; padding: 12px; font-size: 12px; }

/* ====== Main Content ====== */
.scada-main { flex: 1; overflow-y: auto; padding: 16px; background: #0a0e17; }

/* ====== Stats Row ====== */
.stats-row { display: grid; grid-template-columns: repeat(4, 1fr); gap: 16px; margin-bottom: 16px; }
.stat-card {
  background: linear-gradient(135deg, #141c2b, #0d1220);
  border: 1px solid #1e293b; border-radius: 12px; padding: 20px;
  display: flex; align-items: center; gap: 16px; transition: all .3s; cursor: pointer;
  position: relative; overflow: hidden;
}
.stat-card-bg { position: absolute; top: -10px; right: -10px; width: 80px; height: 80px; opacity: .15; }
.stat-card:hover { border-color: #334155; transform: translateY(-2px); box-shadow: 0 8px 25px rgba(0,0,0,.3); }
.stat-icon {
  width: 52px; height: 52px; border-radius: 14px;
  display: flex; align-items: center; justify-content: center; font-size: 24px; flex-shrink: 0;
  z-index: 1;
}
.stat-info { z-index: 1; }
.stat-value { font-size: 28px; font-weight: 700; color: #f1f5f9; line-height: 1; }
.stat-value small { font-size: 14px; font-weight: 400; color: #64748b; }
.stat-label { font-size: 13px; color: #64748b; margin-top: 4px; }

/* ====== Content Grid ====== */
.content-row { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; margin-bottom: 16px; }
.panel {
  background: linear-gradient(135deg, #141c2b, #0d1220);
  border: 1px solid #1e293b; border-radius: 12px;
}
.panel-full { grid-column: 1 / -1; }
.diagram-inner, .chart-inner, .grid-inner { padding: 12px; }
.panel-title {
  font-size: 13px; font-weight: 600; color: #94a3b8;
  text-transform: uppercase; letter-spacing: 1px; margin-bottom: 12px;
}
.chart-container { width: 100%; height: 280px; }

/* ====== Line Diagram ====== */
.line-diagram { display: flex; flex-direction: column; gap: 16px; }
.diagram-bldg {
  background: #0f172a; border-radius: 10px; padding: 12px; cursor: pointer;
  transition: all .2s; border: 1px solid transparent;
}
.diagram-bldg:hover { border-color: #334155; background: #131c2e; }
.diagram-bldg-name { font-size: 14px; font-weight: 600; color: #e2e8f0; margin-bottom: 8px; }
.diagram-floor { padding: 6px 10px; border-radius: 6px; margin: 4px 0; background: #0a0e17; cursor: pointer; }
.diagram-floor:hover { background: #111827; }
.diagram-floor-name { font-size: 12px; color: #94a3b8; margin-bottom: 4px; }
.diagram-device-dots { display: flex; gap: 4px; flex-wrap: wrap; }
.device-dot { width: 8px; height: 8px; border-radius: 50%; background: #22c55e; }
.device-dot.alarm { background: #ef4444; box-shadow: 0 0 4px #ef4444; }
.device-dot.stopped { background: #64748b; }

/* ====== Device Grid ====== */
.device-grid { display: flex; flex-direction: column; gap: 16px; }
.grid-section-title { font-size: 13px; color: #94a3b8; font-weight: 600; margin-bottom: 8px; }
.grid-section-devices { display: grid; grid-template-columns: repeat(auto-fill, minmax(140px, 1fr)); gap: 10px; }
.device-card {
  background: #0f172a; border: 1px solid #1e293b; border-radius: 8px;
  padding: 12px 8px; text-align: center; cursor: pointer; transition: all .2s;
}
.device-card:hover { border-color: #3b82f6; transform: translateY(-2px); box-shadow: 0 4px 12px rgba(59,130,246,.15); }
.device-card.alarm { border-color: rgba(239,68,68,.4); background: rgba(239,68,68,.05); }
.device-card.alarm:hover { border-color: #ef4444; box-shadow: 0 4px 12px rgba(239,68,68,.2); }
.device-card-icon { margin-bottom: 6px; display: flex; justify-content: center; }
.device-card-name { font-size: 11px; color: #cbd5e1; margin-bottom: 4px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.device-card-status { font-size: 10px; padding: 2px 8px; border-radius: 10px; display: inline-block; }
.device-card-status.running { background: rgba(34,197,94,.15); color: #22c55e; }
.device-card-status.alarm { background: rgba(239,68,68,.15); color: #ef4444; }
.device-card-status.stopped { background: rgba(100,116,139,.15); color: #94a3b8; }

/* ====== Level Headers ====== */
.level-header { margin-bottom: 20px; display: flex; align-items: center; gap: 12px; flex-wrap: wrap; }
.level-header h2 { font-size: 22px; color: #f1f5f9; margin: 0; }
.level-subtitle { font-size: 13px; color: #64748b; }
.back-btn {
  cursor: pointer; padding: 6px 14px; border-radius: 6px;
  background: rgba(96,165,250,.1); color: #60a5fa; font-size: 13px;
  transition: all .2s; border: 1px solid transparent;
}
.back-btn:hover { background: rgba(96,165,250,.2); border-color: #3b82f6; }

/* ====== Floor Cards ====== */
.floor-cards { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 16px; }
.floor-card {
  background: #0f172a; border: 1px solid #1e293b; border-radius: 12px;
  cursor: pointer; transition: all .3s;
}
.floor-card:hover { border-color: #3b82f6; transform: translateY(-3px); box-shadow: 0 12px 30px rgba(0,0,0,.3); }
.floor-card-inner { padding: 16px; }
.floor-card-header { display: flex; justify-content: space-between; margin-bottom: 12px; }
.floor-card-name { font-size: 16px; font-weight: 600; color: #e2e8f0; }
.floor-card-count { font-size: 12px; color: #64748b; }
.floor-card-stats { display: flex; gap: 16px; font-size: 12px; }
.floor-stat { padding: 2px 8px; border-radius: 4px; }
.floor-stat.running { background: rgba(34,197,94,.1); color: #22c55e; }
.floor-stat.alarm { background: rgba(239,68,68,.1); color: #ef4444; }
.floor-stat.stopped { background: rgba(100,116,139,.1); color: #94a3b8; }

/* ====== Device Table ====== */
.device-table { background: #0f172a; border-radius: 12px; border: 1px solid #1e293b; overflow: hidden; margin-bottom: 16px; }
.table-inner { padding: 4px; }
.device-table table { width: 100%; border-collapse: collapse; font-size: 13px; }
.device-table th {
  background: #141c2b; color: #64748b; font-weight: 600; padding: 10px 12px;
  text-align: left; border-bottom: 1px solid #1e293b; font-size: 11px;
  text-transform: uppercase; letter-spacing: .5px;
}
.device-table td { padding: 10px 12px; border-bottom: 1px solid #1a2332; color: #cbd5e1; }
.device-table tr:hover { background: rgba(59,130,246,.04); }
.device-table .row-alarm { background: rgba(239,68,68,.04); }
.dev-name { font-weight: 600; color: #e2e8f0; }
.protocol-tag {
  font-size: 10px; padding: 2px 8px; border-radius: 4px;
  background: rgba(59,130,246,.15); color: #60a5fa; font-weight: 600;
}
.status-tag { font-size: 11px; padding: 3px 10px; border-radius: 10px; font-weight: 600; }
.status-tag.running { background: rgba(34,197,94,.15); color: #22c55e; }
.status-tag.alarm { background: rgba(239,68,68,.15); color: #ef4444; }
.status-tag.stopped { background: rgba(100,116,139,.15); color: #94a3b8; }
.status-tag.offline { background: rgba(100,116,139,.15); color: #475569; }
.detail-btn {
  cursor: pointer; padding: 4px 14px; border-radius: 6px;
  background: rgba(96,165,250,.12); color: #60a5fa;
  border: 1px solid rgba(96,165,250,.3); font-size: 12px; transition: all .2s;
}
.detail-btn:hover { background: rgba(96,165,250,.25); }

/* ====== Device Detail Grid ====== */
.device-detail-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.detail-panel { background: #0f172a; border: 1px solid #1e293b; border-radius: 12px; padding: 18px; }
.detail-panel.detail-chart { grid-column: 1 / -1; }
.detail-panel.detail-svg { grid-column: 1 / -1; }
.svg-preview {
  display: flex; justify-content: center; align-items: center;
  min-height: 130px; background: #0a0e17; border-radius: 8px;
}
.detail-chart .chart-container { height: 250px; }
.param-table { width: 100%; border-collapse: collapse; font-size: 13px; }
.param-table td { padding: 8px 0; border-bottom: 1px solid #1a2332; color: #cbd5e1; }
.param-table td:first-child { color: #64748b; width: 40%; }
.param-table td:last-child { font-weight: 600; color: #e2e8f0; }

/* Alarm History */
.alarm-history { display: flex; flex-direction: column; gap: 8px; }
.alarm-history-item {
  display: flex; gap: 12px; align-items: center; padding: 8px 12px;
  border-radius: 6px; background: #141c2b; font-size: 12px; color: #cbd5e1;
}
.alarm-history-item .time { color: #64748b; font-size: 11px; min-width: 100px; }
.alarm-history-item .level {
  font-size: 10px; padding: 1px 8px; border-radius: 10px; font-weight: 600;
  min-width: 50px; text-align: center;
}
.level.level-3 { background: rgba(220,38,38,.15); color: #dc2626; }
.level.level-2 { background: rgba(239,68,68,.15); color: #ef4444; }
.level.level-1 { background: rgba(245,158,11,.15); color: #f59e0b; }
.level.level-0 { background: rgba(59,130,246,.15); color: #3b82f6; }

/* Scrollbar */
.scada-main::-webkit-scrollbar, .scada-sidebar::-webkit-scrollbar, .sidebar-inner::-webkit-scrollbar { width: 6px; }
.scada-main::-webkit-scrollbar-track, .scada-sidebar::-webkit-scrollbar-track, .sidebar-inner::-webkit-scrollbar-track { background: transparent; }
.scada-main::-webkit-scrollbar-thumb, .scada-sidebar::-webkit-scrollbar-thumb, .sidebar-inner::-webkit-scrollbar-thumb { background: #1e293b; border-radius: 3px; }
.scada-main::-webkit-scrollbar-thumb:hover, .scada-sidebar::-webkit-scrollbar-thumb:hover { background: #334155; }
</style>
