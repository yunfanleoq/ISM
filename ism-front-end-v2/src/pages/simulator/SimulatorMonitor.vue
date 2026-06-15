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
            title="A20 电力仪表"
            :value="a20Count"
            prefix="⚡"
            :value-style="{ color: '#52c41a' }"
          />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card" :bordered="false">
          <a-statistic
            title="A40 电力仪表"
            :value="a40Count"
            prefix="🔋"
            :value-style="{ color: '#722ed1' }"
          />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card" :bordered="false">
          <a-statistic
            title="UPS 电源"
            :value="upsCount"
            prefix="💡"
            :value-style="{ color: '#fa8c16' }"
          />
        </a-card>
      </a-col>
    </a-row>

    <!-- 过滤和选择 -->
    <a-row :gutter="16" class="filter-row">
      <a-col :span="8">
        <span style="margin-right: 8px">设备类型:</span>
        <a-radio-group v-model="filterType" button-style="solid" @change="onFilterChange">
          <a-radio-button value="all">全部 ({{ totalSlaves }})</a-radio-button>
          <a-radio-button value="A20">A20 ({{ a20Count }})</a-radio-button>
          <a-radio-button value="A40">A40 ({{ a40Count }})</a-radio-button>
          <a-radio-button value="UPS">UPS ({{ upsCount }})</a-radio-button>
        </a-radio-group>
      </a-col>
      <a-col :span="8">
        <span style="margin-right: 8px">选择设备:</span>
        <a-select
          v-model="selectedSlave"
          style="width: 200px"
          placeholder="选择从站 ID"
          show-search
          :filter-option="filterSlaveOption"
          @change="onSlaveChange"
        >
          <a-select-option
            v-for="item in filteredSlaves"
            :key="item.id"
            :value="item.id"
          >
            从站 {{ item.id }} — {{ item.type }}
          </a-select-option>
        </a-select>
      </a-col>
      <a-col :span="8">
        <span style="margin-right: 8px">数据对比模式:</span>
        <a-switch
          v-model="compareMode"
          checked-children="多设备"
          un-checked-children="单设备"
        />
      </a-col>
    </a-row>

    <!-- 加载状态 -->
    <a-spin :spinning="loading" tip="正在获取模拟器数据...">
      <!-- 单个设备详情 -->
      <template v-if="!compareMode && currentSlave">
        <a-divider orientation="left">
          <a-icon type="hdd" />
          从站 {{ currentSlave.id }} — {{ currentSlave.type }} 实时数据
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
        >
          <template #addr="text">
            <a-tag color="blue">{{ text }}</a-tag>
          </template>
          <template #raw="text">
            <span class="raw-value">{{ text }}</span>
          </template>
          <a-table-column-group title="Holding Registers — 实时数据">
            <a-table-column title="地址" data-index="addr" :width="100" align="center">
              <template #default="{ text }">
                <a-tag color="blue">{{ text }}</a-tag>
              </template>
            </a-table-column>
            <a-table-column title="参数名称" data-index="name" :width="200" />
            <a-table-column title="原始值" data-index="raw" :width="120" align="right">
              <template #default="{ text }">
                <span class="raw-value">{{ text }}</span>
              </template>
            </a-table-column>
            <a-table-column title="单位" data-index="unit" :width="80" align="center">
              <template #default="{ text }">
                <span v-if="text" class="unit-badge">{{ text }}</span>
                <span v-else>-</span>
              </template>
            </a-table-column>
            <a-table-column title="换算" data-index="scale" :width="100" align="center">
              <template #default="{ text }">
                <a-tag v-if="text" color="geekblue">{{ text }}</a-tag>
                <span v-else>-</span>
              </template>
            </a-table-column>
          </a-table-column-group>
        </a-table>
      </template>

      <!-- 多设备对比模式 -->
      <template v-if="compareMode">
        <a-divider orientation="left">
          <a-icon type="table" />
          多设备对比视图（关键参数摘要）
        </a-divider>

        <a-table
          :columns="compareColumns"
          :data-source="compareData"
          :pagination="{ pageSize: 20, showSizeChanger: true, showTotal: t => `共 ${t} 项` }"
          :row-key="r => r.name"
          size="middle"
          bordered
          :scroll="{ x: 1200 }"
        >
          <template #name="text">
            <strong>{{ text }}</strong>
          </template>
        </a-table>
      </template>

      <!-- 无数据提示 -->
      <a-empty v-if="!compareMode && !currentSlave && !loading" description="请选择一个从站设备查看详细数据" />
    </a-spin>
  </div>
</template>

<script>
export default {
  name: 'SimulatorMonitor',
  data() {
    return {
      connected: false,
      loading: false,
      autoRefresh: true,
      refreshInterval: 2,
      updateCount: 0,
      filterType: 'all',
      selectedSlave: null,
      compareMode: false,

      // API 数据
      totalSlaves: 0,
      a20Count: 0,
      a40Count: 0,
      upsCount: 0,
      allSlaves: [],
      allSlavesData: {},

      // 定时器
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
      return this.allSlavesData[this.selectedSlave] || null
    },
    hrColumns() {
      return [
        { title: '地址', dataIndex: 'addr', key: 'addr', width: 80, align: 'center',
          customRender: (text) => ({ children: text != null ? `HR ${text}` : '', attrs: {} }) },
        { title: '参数名称', dataIndex: 'name', key: 'name', width: 200 },
        { title: '原始值', dataIndex: 'raw', key: 'raw', width: 120, align: 'right',
          customRender: (text) => ({ children: text != null ? text.toLocaleString() : '-', attrs: {} }) },
        { title: '单位', dataIndex: 'unit', key: 'unit', width: 80, align: 'center' },
        { title: '换算', dataIndex: 'scale', key: 'scale', width: 100, align: 'center' },
      ]
    },
    compareColumns() {
      return [
        { title: '参数名称', dataIndex: 'name', key: 'name', fixed: 'left', width: 180 },
        ...this.getCompareDeviceColumns(),
      ]
    },
    compareData() {
      // 提取所有设备共有的参数名，展示原始值
      const paramMap = {}
      for (const sid of Object.keys(this.allSlavesData)) {
        const slave = this.allSlavesData[sid]
        if (!slave) continue
        for (const hr of slave.hr) {
          if (!paramMap[hr.name]) {
            paramMap[hr.name] = { name: hr.name }
          }
          paramMap[hr.name][`slave_${slave.id}`] = hr.raw
        }
      }
      return Object.values(paramMap)
    },
  },
  methods: {
    filterSlaveOption(input, option) {
      const text = option.componentOptions.children[0].text || ''
      return text.toLowerCase().includes(input.toLowerCase())
    },
    onFilterChange() {
      this.selectedSlave = null
    },
    onSlaveChange(value) {
      if (value) {
        this.fetchSlaveDetail(value)
      }
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
        const resp = await fetch('http://127.0.0.1:5040/api/summary')
        if (!resp.ok) throw new Error('API error')
        const summary = await resp.json()
        this.connected = true
        this.totalSlaves = summary.total_slaves
        this.a20Count = summary.a20_range[1] - summary.a20_range[0] + 1
        this.a40Count = summary.a40_range[1] - summary.a40_range[0] + 1
        this.upsCount = summary.ups_range[1] - summary.ups_range[0] + 1
        this.allSlaves = summary.slaves

        // 预加载所有设备数据
        await this.fetchAllSlavesPreview()
        this.updateCount++

        // 如果当前选中了设备，静默更新详情
        if (this.selectedSlave) {
          await this.fetchSlaveDetail(this.selectedSlave)
        }
      } catch (e) {
        if (!silent) console.error('模拟器连接失败:', e.message)
        this.connected = false
      } finally {
        if (!silent) this.loading = false
      }
    },

    async fetchAllSlavesPreview() {
      try {
        const resp = await fetch('http://127.0.0.1:5040/api/slaves')
        if (!resp.ok) return
        const data = await resp.json()
        const slaveList = []
        for (const item of data) {
          const transformed = this.transformSlaveData(item)
          this.$set(this.allSlavesData, item.slave, transformed)
          slaveList.push({ id: item.slave, type: item.device_type })
        }
        this.allSlaves = slaveList
      } catch (e) {
        console.error('批量获取失败:', e.message)
      }
    },

    async fetchSlaveDetail(sid) {
      try {
        const resp = await fetch(`http://127.0.0.1:5040/api/slave/${sid}`)
        if (!resp.ok) return
        const data = await resp.json()
        this.$set(this.allSlavesData, sid, this.transformSlaveData(data))
      } catch (e) {
        console.error('获取从站详情失败:', e.message)
      }
    },

    getRegisterUnit(name) {
      const unitMap = {
        '电压': 'V',
        '电流': 'A',
        '频率': 'Hz',
        '有功功率': 'kW',
        '无功功率': 'kvar',
        '视在功率': 'kVA',
        '功率因数': '',
        '有功电度': 'kWh',
        '谐波畸变率': '%',
      }
      for (const [key, unit] of Object.entries(unitMap)) {
        if (name.includes(key)) return unit
      }
      return ''
    },

    transformSlaveData(rawData) {
      const result = {
        id: rawData.slave,
        type: rawData.device_type,
        hr: [],
        di: [],
      }
      // transform holding_registers
      if (rawData.holding_registers) {
        for (const [addr, reg] of Object.entries(rawData.holding_registers)) {
          result.hr.push({
            addr: parseInt(addr),
            name: reg.name,
            raw: reg.raw != null ? reg.raw : reg.value,
            value: reg.value,
            unit: this.getRegisterUnit(reg.name),
            scale: '',
          })
        }
        result.hr.sort((a, b) => a.addr - b.addr)
      }
      // transform discrete_inputs
      if (rawData.discrete_inputs) {
        if (Array.isArray(rawData.discrete_inputs)) {
          result.di = rawData.discrete_inputs.map((di, idx) => ({
            addr: di.addr != null ? di.addr : idx,
            name: di.name || `DI ${idx}`,
            value: di.value,
            label: di.value === 1 ? '正常' : '异常',
          }))
        } else {
          for (const [addr, val] of Object.entries(rawData.discrete_inputs)) {
            const v = typeof val === 'object' ? val.value : val
            result.di.push({
              addr: parseInt(addr),
              name: typeof val === 'object' ? (val.name || `DI ${addr}`) : `DI ${addr}`,
              value: v,
              label: v === 1 ? '正常' : '异常',
            })
          }
        }
        result.di.sort((a, b) => a.addr - b.addr)
      }
      return result
    },

    getCompareDeviceColumns() {
      const sids = Object.keys(this.allSlavesData)
        .filter(sid => {
          const slave = this.allSlavesData[sid]
          if (!slave) return false
          if (this.filterType !== 'all' && slave.type !== this.filterType) return false
          return true
        })
        .sort((a, b) => Number(a) - Number(b))
        .slice(0, 10) // 最多显示10个设备，避免表格过宽

      return sids.map(sid => {
        const slave = this.allSlavesData[sid]
        const title = slave ? `Slave ${sid} (${slave.type})` : `Slave ${sid}`
        return {
          title,
          dataIndex: `slave_${sid}`,
          key: `slave_${sid}`,
          width: 140,
          align: 'right',
        }
      })
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
