<template>
  <div class="simulator-monitor">
    <!-- 顶部标题栏 -->
    <div class="monitor-header">
      <h1 class="monitor-title">
        <a-icon type="api" style="margin-right: 8px" />
        Modbus 模拟器实时监控
      </h1>
      <div class="header-actions">
        <a-tag :color="connected ? 'green' : 'red'">
          {{ connected ? '已连接' : '未连接' }}
        </a-tag>
        <a-badge :status="connected ? 'processing' : 'default'" :text="`${updateCount} 次更新`" />
        <a-switch
          v-model="autoRefresh"
          checked-children="自动刷新"
          un-checked-children="手动"
          @change="onAutoRefreshToggle"
        />
        <a-button
          type="primary"
          icon="reload"
          :loading="loading"
          @click="fetchData"
          style="margin-left: 8px"
        >
          刷新
        </a-button>
        <a-input-number
          v-model="refreshInterval"
          :min="1"
          :max="60"
          :step="1"
          addon-after="秒"
          @change="onIntervalChange"
          style="width: 110px; margin-left: 8px"
        />
      </div>
    </div>

    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-card class="stat-card" :bordered="false">
          <a-statistic
            title="模拟设备总数"
            :value="totalSlaves"
            prefix="📡"
            :value-style="{ color: '#1890ff' }"
          />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card" :bordered="false">
          <a-statistic
            title="网关 IP 数"
            :value="ipCount"
            prefix="🌐"
            :value-style="{ color: '#52c41a' }"
          />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card" :bordered="false">
          <a-statistic
            title="数据点总数"
            :value="totalPoints"
            prefix="📊"
            :value-style="{ color: '#722ed1' }"
          />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card" :bordered="false">
          <a-statistic
            title="设备型号数"
            :value="deviceTypes.length"
            prefix="🏷️"
            :value-style="{ color: '#fa8c16' }"
          />
        </a-card>
      </a-col>
    </a-row>

    <!-- 设备型号分布 -->
    <a-card
      v-if="deviceTypes.length"
      class="type-dist-card"
      size="small"
      :bordered="false"
      style="margin-bottom: 20px"
    >
      <template #title>
        <a-icon type="pie-chart" style="margin-right: 6px" />设备型号分布
      </template>
      <a-tag
        v-for="dt in deviceTypes"
        :key="dt.type"
        color="blue"
        style="margin: 4px 6px 4px 0"
      >
        {{ dt.type }}：{{ dt.count }} 台
      </a-tag>
    </a-card>

    <!-- 过滤和选择 -->
    <a-row :gutter="16" class="filter-row">
      <a-col :span="8">
        <span style="margin-right: 8px">设备类型:</span>
        <a-radio-group v-model="filterType" button-style="solid" @change="onFilterChange">
          <a-radio-button value="all">全部 ({{ totalSlaves }})</a-radio-button>
          <a-radio-button v-for="dt in deviceTypes" :key="dt.type" :value="dt.type">
            {{ dt.type }} ({{ dt.count }})
          </a-radio-button>
        </a-radio-group>
      </a-col>
      <a-col :span="8">
        <span style="margin-right: 8px">选择设备:</span>
        <a-select
          v-model="selectedSlave"
          style="width: 320px"
          placeholder="选择设备（IP / 从站）"
          show-search
          :filter-option="filterSlaveOption"
          @change="onSlaveChange"
        >
          <a-select-option
            v-for="item in filteredSlaves"
            :key="item.id"
            :value="item.id"
          >
            {{ item.ip }} / 从站{{ item.slave_id }} — {{ item.type }}
          </a-select-option>
        </a-select>
      </a-col>
      <a-col :span="8">
        <span style="margin-right: 8px">数据对比模式:</span>
        <a-switch
          v-model="compareMode"
          checked-children="多设备"
          un-checked-children="单设备"
          @change="onCompareToggle"
        />
        <span style="margin-left: 8px; color: #999; font-size: 12px">
          按需对比前 {{ compareLimit }} 台
        </span>
      </a-col>
    </a-row>

    <!-- 加载状态 -->
    <a-spin :spinning="loading" tip="正在获取模拟器数据...">
      <!-- 单个设备详情 -->
      <template v-if="!compareMode && currentSlave">
        <a-divider orientation="left">
          <a-icon type="hdd" />
          {{ currentSlave.ip }} / 从站{{ currentSlave.slave_id }} — {{ currentSlave.type }} 实时数据
        </a-divider>

        <!-- 离散输入状态 -->
        <a-row :gutter="16" style="margin-bottom: 16px">
          <a-col :span="8" v-for="di in currentSlave.di" :key="di.addr">
            <a-card size="small" :bordered="true">
              <div style="display: flex; justify-content: space-between; align-items: center">
                <span>{{ di.name }} (DI {{ di.addr }})</span>
                <a-tag :color="di.value === 1 ? 'green' : 'red'">
                  {{ di.label }}
                </a-tag>
              </div>
            </a-card>
          </a-col>
        </a-row>

        <!-- Holding Registers 表格 -->
        <a-table
          :columns="hrColumns"
          :data-source="currentSlave.hr"
          :pagination="false"
          :row-key="r => r.addr"
          size="middle"
          bordered
        />
      </template>

      <!-- 多设备对比模式 -->
      <template v-if="compareMode">
        <a-divider orientation="left">
          <a-icon type="table" />
          多设备对比视图（关键参数摘要 · 共 {{ compareIds.length }} 台）
        </a-divider>

        <a-table
          :columns="compareColumns"
          :data-source="compareData"
          :pagination="{ pageSize: 20, showSizeChanger: true, showTotal: t => `共 ${t} 项` }"
          :row-key="r => r.name"
          size="middle"
          bordered
          :scroll="{ x: 1200 }"
        />
      </template>

      <!-- 无数据提示 -->
      <a-empty v-if="!compareMode && !currentSlave && !loading" description="请选择一个从站设备查看详细数据" />
    </a-spin>
  </div>
</template>

<script>
const API_BASE = 'http://127.0.0.1:5040'

export default {
  name: 'SimulatorMonitor',
  data() {
    return {
      connected: false,
      loading: false,
      autoRefresh: true,
      refreshInterval: 10,
      updateCount: 0,
      filterType: 'all',
      selectedSlave: null,
      compareMode: false,
      compareLimit: 15,

      // 摘要统计（新版 5040 API）
      totalSlaves: 0,
      ipCount: 0,
      totalPoints: 0,
      deviceTypes: [],   // [{ type, count }]
      byIp: [],          // [{ ip, slaves }]

      // 设备清单（轻量，仅 {id, slave_id, ip, type}），详情按需加载
      allSlaves: [],
      slavesData: {},    // id -> 已转换的设备详情（仅按需拉取）
      compareIds: [],    // 当前对比的设备 id 列表

      timer: null,
    }
  },
  computed: {
    filteredSlaves() {
      if (this.filterType === 'all') return this.allSlaves
      return this.allSlaves.filter(s => s.type === this.filterType)
    },
    currentSlave() {
      if (!this.selectedSlave) return null
      return this.slavesData[this.selectedSlave] || null
    },
    hrColumns() {
      return [
        { title: '地址', dataIndex: 'addr', key: 'addr', width: 90, align: 'center',
          customRender: (text) => ({ children: text != null ? `HR ${text}` : '', attrs: {} }) },
        { title: '参数名称', dataIndex: 'name', key: 'name', width: 220 },
        { title: '原始值', dataIndex: 'raw', key: 'raw', width: 120, align: 'right',
          customRender: (text) => ({ children: text != null ? Number(text).toLocaleString() : '-', attrs: {} }) },
        { title: '换算值', dataIndex: 'value', key: 'value', width: 120, align: 'right',
          customRender: (text) => ({ children: text != null ? text : '-', attrs: {} }) },
        { title: '单位', dataIndex: 'unit', key: 'unit', width: 80, align: 'center' },
      ]
    },
    compareColumns() {
      const cols = [
        { title: '参数名称', dataIndex: 'name', key: 'name', fixed: 'left', width: 200 },
      ]
      // id 含 "#" 与 "." 不能直接做 dataIndex（antd 会按路径解析），统一映射为 col0/col1...
      this.compareIds.forEach((id, idx) => {
        const slave = this.slavesData[id]
        const title = slave ? `${slave.ip}/从站${slave.slave_id}` : id
        cols.push({
          title,
          dataIndex: `col${idx}`,
          key: `col${idx}`,
          width: 150,
          align: 'right',
        })
      })
      return cols
    },
    compareData() {
      // 以参数名为行，对比各设备的原始值
      const paramMap = {}
      this.compareIds.forEach((id, idx) => {
        const slave = this.slavesData[id]
        if (!slave) return
        for (const hr of slave.hr) {
          if (!paramMap[hr.name]) paramMap[hr.name] = { name: hr.name }
          paramMap[hr.name][`col${idx}`] = hr.raw
        }
      })
      return Object.values(paramMap)
    },
  },
  methods: {
    filterSlaveOption(input, option) {
      const text = (option.componentOptions.children[0].text || '').trim()
      return text.toLowerCase().includes(input.toLowerCase())
    },
    onFilterChange() {
      this.selectedSlave = null
      if (this.compareMode) this.loadCompare()
    },
    onSlaveChange(value) {
      if (value) this.fetchSlaveDetail(value)
    },
    onCompareToggle(checked) {
      if (checked) this.loadCompare()
    },
    onAutoRefreshToggle(checked) {
      if (checked) {
        this.startAutoRefresh()
      } else {
        this.stopAutoRefresh()
      }
    },
    onIntervalChange() {
      if (this.autoRefresh) {
        this.stopAutoRefresh()
        this.startAutoRefresh()
      }
    },

    async fetchData(silent = false) {
      if (!silent) this.loading = true
      try {
        const resp = await fetch(`${API_BASE}/api/summary`)
        if (!resp.ok) throw new Error('API error')
        const summary = await resp.json()
        this.connected = true
        this.totalSlaves = summary.total_slaves || 0
        this.ipCount = summary.ip_count || 0
        this.totalPoints = summary.total_points || 0
        this.deviceTypes = summary.device_types || []
        this.byIp = summary.by_ip || []
        // 轻量清单：{id, slave_id, ip, type}，不预加载任何寄存器明细
        this.allSlaves = summary.slaves || []
        this.updateCount++

        // 详情仅对当前可见对象按需刷新
        if (!this.compareMode && this.selectedSlave) {
          await this.fetchSlaveDetail(this.selectedSlave)
        } else if (this.compareMode && this.compareIds.length) {
          await this.refreshCompare()
        }
      } catch (e) {
        if (!silent) console.error('模拟器连接失败:', e.message)
        this.connected = false
      } finally {
        if (!silent) this.loading = false
      }
    },

    async fetchSlaveDetail(id) {
      try {
        // id 形如 "172.31.4.12#5"，"#" 必须编码，否则会被当成 URL fragment 丢失
        const resp = await fetch(`${API_BASE}/api/slave/${encodeURIComponent(id)}`)
        if (!resp.ok) return
        const data = await resp.json()
        this.$set(this.slavesData, id, this.transformSlaveData(data))
      } catch (e) {
        console.error('获取从站详情失败:', e.message)
      }
    },

    async loadCompare() {
      // 进入对比模式 / 切换过滤类型：仅按需拉取前 N 台，避免全量预加载
      this.compareIds = this.filteredSlaves.slice(0, this.compareLimit).map(s => s.id)
      await this.refreshCompare()
    },

    async refreshCompare() {
      await Promise.all(this.compareIds.map(id => this.fetchSlaveDetail(id)))
    },

    getRegisterUnit(name) {
      const unitMap = {
        '线电压': 'V',
        '相电压': 'V',
        '电压': 'V',
        '电流': 'A',
        '频率': 'Hz',
        '有功功率': 'kW',
        '无功功率': 'kvar',
        '视在功率': 'kVA',
        '功率因数': '',
        '有功电度': 'kWh',
        '有功电能': 'kWh',
        '谐波畸变率': '%',
        '畸变率': '%',
        '温度': '℃',
        '湿度': '%',
      }
      for (const [key, unit] of Object.entries(unitMap)) {
        if (name.includes(key)) return unit
      }
      return ''
    },

    transformSlaveData(rawData) {
      const result = {
        id: rawData.slave,            // "ip#slave"
        slave_id: rawData.slave_id,
        ip: rawData.ip,
        type: rawData.device_type,
        hr: [],
        di: [],
      }
      // holding_registers: { addr: { name, raw, value } }
      if (rawData.holding_registers) {
        for (const [addr, reg] of Object.entries(rawData.holding_registers)) {
          result.hr.push({
            addr: parseInt(addr),
            name: reg.name,
            raw: reg.raw != null ? reg.raw : reg.value,
            value: reg.value,
            unit: this.getRegisterUnit(reg.name || ''),
          })
        }
        result.hr.sort((a, b) => a.addr - b.addr)
      }
      // discrete_inputs: { addr: { name, value } }
      if (rawData.discrete_inputs) {
        for (const [addr, val] of Object.entries(rawData.discrete_inputs)) {
          const v = typeof val === 'object' ? val.value : val
          result.di.push({
            addr: parseInt(addr),
            name: typeof val === 'object' ? (val.name || `DI ${addr}`) : `DI ${addr}`,
            value: v,
            label: v === 1 ? '正常' : '异常',
          })
        }
        result.di.sort((a, b) => a.addr - b.addr)
      }
      return result
    },

    startAutoRefresh() {
      this.stopAutoRefresh()
      this.timer = setInterval(() => {
        if (this.autoRefresh) {
          this.fetchData(true)  // 静默刷新，不显示 loading 遮罩
        }
      }, this.refreshInterval * 1000)
    },

    stopAutoRefresh() {
      if (this.timer) {
        clearInterval(this.timer)
        this.timer = null
      }
    },
  },

  mounted() {
    this.fetchData()
    if (this.autoRefresh) {
      this.startAutoRefresh()
    }
  },

  beforeDestroy() {
    this.stopAutoRefresh()
  },
}
</script>

<style lang="less" scoped>
.simulator-monitor {
  min-height: 100vh;
  background: #f0f2f5;
  padding: 24px;
}

.monitor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 16px 24px;
  background: #fff;
  border-radius: 4px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}

.monitor-title {
  margin: 0;
  font-size: 22px;
  font-weight: 600;
  color: #1a1a1a;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  text-align: center;
  border-radius: 4px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
}

.type-dist-card {
  border-radius: 4px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
}

.filter-row {
  margin-bottom: 20px;
  padding: 12px 16px;
  background: #fff;
  border-radius: 4px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
  display: flex;
  align-items: center;
}

.raw-value {
  font-family: 'SF Mono', 'Menlo', 'Consolas', monospace;
  font-weight: 600;
  color: #1890ff;
  font-size: 14px;
}

.unit-badge {
  display: inline-block;
  background: #e6f7ff;
  color: #1890ff;
  padding: 2px 8px;
  border-radius: 3px;
  font-size: 12px;
  font-weight: 500;
}
</style>
