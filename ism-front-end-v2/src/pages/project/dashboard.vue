<template>
  <div id="WaterMarkShow" class="analysis">
    <a-spin :tip="$t('dashBoard.loading')" size="large" :spinning="OptLoading">
      <!-- 1. 核心指标卡片区：渐变背景+悬浮动效 -->
      <a-row class="metrics-card-row" :gutter="[24, 24]">
        <a-col :sm="24" :md="12" :xl="4" v-for="(item, index) in metricsList" :key="index">
          <a-card
              :class="`metrics-card border-${item.type}`"
              @mouseenter="cardHover = index"
              @mouseleave="cardHover = -1"
          >
            <p class="metrics-label">{{ $t(item.i18nKey) }}</p>
            <div class="metrics-value-wrap">
              <span class="metrics-value">{{ SystemInfo[item.dataKey] }}</span>
              <span class="metrics-unit">{{ item.unit }}</span>
            </div>
            <!-- 悬浮时显示图标，增强交互反馈 -->
            <div class="metrics-icon" v-show="cardHover === index">
              <a-icon :type="item.icon" />
            </div>
          </a-card>
        </a-col>
      </a-row>

    <a-row class="resource-row" :gutter="[16, 16]">
      <a-col :xs="24" :xl="8">
        <a-card   class=" metrics-card-l border-info r" title="CPU" >
          <div class="resource-card-body">
            <div class="resource-chart cpu-chart-wrap">
              <div class="view-chart-gauge" ref="cpuCharts"></div>
            </div>
            <div class="resource-stats">
            <a-form class="resource-form resource-form-grid">
              <a-form-item :label="$t('dashBoard.CUPCount')">
                <span class="ant-form-text">
                  {{SystemInfo.CpuInfo.Number}}
                </span>
              </a-form-item>
              <a-form-item :label="$t('dashBoard.CUPKernelCount')">
                <span class="ant-form-text">
                  {{SystemInfo.CpuInfo.Cores}}
                </span>
              </a-form-item>
              <a-form-item :label="$t('dashBoard.CUPUsedCount')">
                <span class="ant-form-text">
                  {{SystemInfo.CpuInfo.UsedPercent[0]}} %
                </span>
              </a-form-item>
              <a-form-item :label="$t('dashBoard.LA5')">
                <span class="ant-form-text">
                  {{SystemInfo.CpuInfo.Load.load5}} %
                </span>
              </a-form-item>
              <a-form-item :label="$t('dashBoard.LA15')">
                <span class="ant-form-text">
                  {{SystemInfo.CpuInfo.Load.load15}} %
                </span>
              </a-form-item>
              <div class="resource-stat-placeholder"></div>
            </a-form>
            </div>
          </div>

        </a-card>
      </a-col>
      <a-col :xs="24" :xl="8">
        <a-card   class=" metrics-card-l border-danger r" title="硬盘">
          <div class="resource-card-body">
            <div class="resource-chart disk-chart-wrap">
              <div class="view-chart-gauge" ref="viewDisk"></div>
            </div>
            <div class="resource-stats">
              <a-form class="resource-form resource-form-grid">
                <a-form-item :label="$t('dashBoard.DiskCap')">
                <span class="ant-form-text">
                  {{SystemInfo.DiskInfo.total}} GB
                </span>
                </a-form-item>
                <a-form-item :label="$t('dashBoard.UsedDisk')">
                <span class="ant-form-text">
                  {{SystemInfo.DiskInfo.used}} GB
                </span>
                </a-form-item>
                <a-form-item :label="$t('dashBoard.MemRem')">
                <span class="ant-form-text">
                  {{SystemInfo.DiskInfo.free}} GB
                </span>
                </a-form-item>
                <a-form-item :label="$t('dashBoard.UsedDiskPro')">
                <span class="ant-form-text">
                  {{SystemInfo.DiskInfo.usedPercent}} %
                </span>
                </a-form-item>
                <div class="resource-stat-placeholder"></div>
                <div class="resource-stat-placeholder"></div>
              </a-form>
            </div>
          </div>
        </a-card>
      </a-col>
      <a-col :xs="24" :xl="8">
        <a-card   class="metrics-card-l border-success r" title="内存">
          <div class="resource-card-body">
            <div class="resource-chart mem-chart-wrap">
              <div class="view-chart-gauge" ref="memCharts"></div>
            </div>
            <div class="resource-stats">
              <a-form class="resource-form resource-form-grid">
                <a-form-item :label="$t('dashBoard.MemCap')">
                <span class="ant-form-text">
                  {{SystemInfo.MemInfo.total}} GB
                </span>
                </a-form-item>
                <a-form-item :label="$t('dashBoard.UsedMem')">
                <span class="ant-form-text">
                  {{SystemInfo.MemInfo.used}} GB
                </span>
                </a-form-item>
                <a-form-item :label="$t('dashBoard.MemRem')">
                  <span class="ant-form-text">
                    {{SystemInfo.MemInfo.free}} GB
                  </span>
                </a-form-item>

                <a-form-item :label="$t('dashBoard.ProgramUsed')">
                  <span class="ant-form-text">
                    {{SystemInfo.MemInfo.goUsed}} MB
                  </span>
                </a-form-item>

                <a-form-item :label="$t('dashBoard.UsedMemPro')">
                  <span class="ant-form-text">
                    {{SystemInfo.MemInfo.usedPercent}} %
                  </span>
                </a-form-item>
                <div class="resource-stat-placeholder"></div>
              </a-form>
            </div>
          </div>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16" style="margin-top: 10px">
      <a-card class="metrics-card-l border-success r"  title="运行环境信息">
        <a-row style="padding: 10px">
          <a-col  :span="8"  style="padding: 20px">
            <a-form :labelCol=" { span: 6 }"
                    :wrapperCol=" { span: 14 }">
              <a-form-item :label="$t('dashBoard.system')">
                <span class="ant-form-text">
                  {{SystemInfo.HostInfo.platform}}
                </span>
              </a-form-item>
              <a-form-item :label="$t('dashBoard.BootTime')">
                <span class="ant-form-text">
                   {{SystemInfo.HostInfo.bootTime}}
                </span>
              </a-form-item>
              <a-form-item :label="$t('dashBoard.RunTime')">
                <span class="ant-form-text">
                  {{SystemInfo.HostInfo.uptime}} 小时
                </span>
              </a-form-item>
              <a-form-item :label="$t('dashBoard.RunMem')">
                <span class="ant-form-text">
                   {{SystemInfo.MemInfo.used}} GB
                </span>
              </a-form-item>
            </a-form>
          </a-col>

          <a-col  :span="8"  style="padding: 20px">
            <a-form :labelCol=" { span: 6 }"
                    :wrapperCol=" { span: 14 }">
              <a-form-item :label="$t('dashBoard.systemType')">
                <span class="ant-form-text">
                  {{SystemInfo.HostInfo.kernelArch}}
                </span>
              </a-form-item>
              <a-form-item :label="$t('dashBoard.ProCount')">
                <span class="ant-form-text">
                  {{SystemInfo.HostInfo.procs}}
                </span>
              </a-form-item>
              <a-form-item :label="$t('dashBoard.MacID')">
                <span class="ant-form-text">
                  {{ProtectedID}}
                </span>
              </a-form-item>
              <a-form-item :label="$t('dashBoard.ProgramUsed')">
                <span class="ant-form-text">
                   {{SystemInfo.MemInfo.goUsed}} MB
                </span>
              </a-form-item>
            </a-form>
          </a-col>

          <a-col  :span="8"  style="padding: 20px">
            <a-form :labelCol=" { span: 6 }"
                    :wrapperCol=" { span: 14 }">
              <a-form-item :label="$t('dashBoard.Goroutine')">
                <span class="ant-form-text">
                  {{SystemInfo.Goroutine}}
                </span>
              </a-form-item>
              <a-form-item :label="$t('dashBoard.netRecvied')">
                <span class="ant-form-text">
                  {{SystemInfo.NetWork.Receive}} MB
                </span>
              </a-form-item>
              <a-form-item :label="$t('dashBoard.netSend')">
                <span class="ant-form-text">
                  {{SystemInfo.NetWork.Sent}} MB
                </span>
              </a-form-item>
              <a-form-item :label="$t('SystemUpgrade.Version')">
                <span class="ant-form-text">
                  {{systemVersion}}
                </span>
              </a-form-item>
            </a-form>
          </a-col>

        </a-row>
      </a-card>
    </a-row>
    </a-spin>
  </div>
</template>

<script>
import * as echarts from 'echarts';
import {GetSystemAnalysis, GetSystemParams} from "@/services/system";
import {mapState, mapMutations} from 'vuex'
import watermark from "@/utils/watermark.js";
export default {
  name: 'Analysis',
  i18n: require('../../i18n/language'),
  data () {
    return {
      loading: true,
      OptLoading:true,
      echartsView:null,
      echartsViewMem:null,
      echartsViewDisk:null,
      refreshTimer: null,
      cpuHistoryLimit: 20,
      cpuHistoryLabels: [],
      cpuHistoryValues: [],
      memHistoryLimit: 20,
      memHistoryLabels: [],
      memHistoryValues: [],
      optionCpu:{
        animation: false,
        grid: {
          top: 18,
          left: 4,
          right: 2,
          bottom: 14,
          containLabel: true
        },
        tooltip: {
          trigger: 'axis',
          formatter: function (params) {
            if (!params.length) return ''
            return params[0].axisValue + '<br/>CPU: ' + params[0].value + '%'
          }
        },
        xAxis: {
          type: 'category',
          boundaryGap: false,
          data: [],
          axisLine: {
            lineStyle: {
              color: '#d9e2f3'
            }
          },
          axisLabel: {
            color: '#7f8fa4',
            fontSize: 10
          },
        },
        yAxis: {
          type: 'value',
          min: 0,
          max: 100,
          splitNumber: 4,
          axisLabel: {
            color: '#7f8fa4',
            formatter: '{value}%'
          },
          splitLine: {
            lineStyle: {
              color: 'rgba(27, 153, 199, 0.12)'
            }
          }
        },
        series: [{
          name: 'CPU',
          type: 'line',
          smooth: true,
          symbol: 'none',
          data: [],
          lineStyle: {
            width: 3,
            color: '#1b99c7'
          },
          areaStyle: {
            color: 'rgba(27, 153, 199, 0.18)'
          }
        }],
      },
      optionDisk:{
        title: [
          {
            text: '00%',
            x: 'center',
            top: '42%',
            left:'50%',
            textStyle: {
              fontSize: 30,
              color: '#000000',
              fontFamily: 'DINAlternate-Bold, DINAlternate',
              fontWeight: '600',
            },
            textAlign: 'center',
          }],
        backgroundColor: '',
        polar: {
          radius: ['88%', '72%'],
          center: ['50%', '48%']
        },
        angleAxis: {
          max: 100,
          show: false,
        },
        radiusAxis: {
          type: 'category',
          show: true,
          axisLabel: {
            show: false,
          },
          axisLine: {
            show: false,
          },
          axisTick: {
            show: false,
          },
        },
        series: [
          {
            type: 'bar',
            roundCap: true,
            barWidth: 10,
            showBackground: true,
            backgroundStyle: {
              color: 'rgba(66, 66, 66, .2)',
            },
            data: [10],
            coordinateSystem: 'polar',

            itemStyle: {
              normal: {
                color: 'rgba(96, 244, 194, 1)',
              },
            },
          },
        ],
      },
      optionMem:{
        animation: false,
        grid: {
          top: 18,
          left: 4,
          right: 2,
          bottom: 14,
          containLabel: true
        },
        tooltip: {
          trigger: 'axis',
          formatter: function (params) {
            if (!params.length) return ''
            return params[0].axisValue + '<br/>MEM: ' + params[0].value + '%'
          }
        },
        xAxis: {
          type: 'category',
          boundaryGap: false,
          data: [],
          axisLine: {
            lineStyle: {
              color: '#d9e2f3'
            }
          },
          axisLabel: {
            color: '#7f8fa4',
            fontSize: 10
          },
        },
        yAxis: {
          type: 'value',
          min: 0,
          max: 100,
          splitNumber: 4,
          axisLabel: {
            color: '#7f8fa4',
            formatter: '{value}%'
          },
          splitLine: {
            lineStyle: {
              color: 'rgba(84, 199, 27, 0.12)'
            }
          }
        },
        series: [{
          name: 'MEM',
          type: 'line',
          smooth: true,
          symbol: 'none',
          data: [],
          lineStyle: {
            width: 3,
            color: '#54c71b'
          },
          areaStyle: {
            color: 'rgba(84, 199, 27, 0.18)'
          }
        }],
      },
      SystemInfo:{
        CpuInfo:{
          Number:0,
          Cores:0,
          UsedPercent:['00.00'],
          Load: {
            "load1": 0,
            "load5": 0,
            "load15": 0
          }
        },
        DiskInfo: {
          "total": 40,
          "free": 20,
          "used": 10,
          "usedPercent": 0,
        },
        MemInfo: {
          "available": 0,
          "free": 0,
          "goUsed": 0,
          "total": 0,
          "used": 0,
          "usedPercent": 0,
        },
        HostInfo: {
          "hostname": "DESKTOP-VDQD76M",
          "uptime": 0,
          "bootTime": 0,
          "procs": 0,
          "os": "windows",
          "platform": "Microsoft Windows 10 Pro",
          "platformFamily": "Standalone Workstation",
          "platformVersion": "10.0.19045.3693 Build 19045.3693",
          "kernelVersion": "10.0.19045.3693 Build 19045.3693",
          "kernelArch": "x86_64",
          "virtualizationSystem": "",
          "virtualizationRole": "",
          "hostId": "f60ccb2f-a5cd-4bd6-856a-3f0723fb27a7"
        },
        NetWork:{
          Receive:0,
          Sent:0
        }
      },
      IsLicense:true,
      cardHover: -1, // 卡片悬浮索引
      // 核心指标配置（统一管理，便于维护）
      metricsList: [
        { type: 'info', i18nKey: 'dashboard.APPCount', dataKey: 'AppCount', unit: '个', icon: 'appstore' },
        { type: 'info', i18nKey: 'dashboard.DeviceCount', dataKey: 'DeviceCount', unit: '个', icon: 'laptop' },
        { type: 'danger', i18nKey: 'dashboard.DeviceOfflineCount', dataKey: 'DeviceOffCount', unit: '个', icon: 'wifi-off' },
        { type: 'danger', i18nKey: 'dashboard.alarmCount', dataKey: 'AlarmCount', unit: '个', icon: 'exclamation-circle' },
        { type: 'success', i18nKey: 'dashboard.DeviceDataCount', dataKey: 'DataCount', unit: '个', icon: 'database' },
        { type: 'success', i18nKey: 'dashboard.VideoCount', dataKey: 'VideoCount', unit: '个', icon: 'play-circle' }
      ],
    }
  },
  computed: {
    ...mapState('setting', ['langList','isMobile','systemVersion','ProtectedID']),
  },
  watch: {
    $route () {
      this.$nextTick(() => {
        this.initMap()
        this.GetSystemAnalysis()
      })
    }
  },
  methods: {
    getCurrentTimeLabel() {
      const now = new Date()
      const pad = (value) => String(value).padStart(2, '0')
      return `${pad(now.getHours())}:${pad(now.getMinutes())}:${pad(now.getSeconds())}`
    },
    updateCpuChart(value) {
      const cpuValue = Number(parseFloat(value).toFixed(2))
      this.cpuHistoryLabels.push(this.getCurrentTimeLabel())
      this.cpuHistoryValues.push(cpuValue)
      if (this.cpuHistoryLabels.length > this.cpuHistoryLimit) {
        this.cpuHistoryLabels.shift()
        this.cpuHistoryValues.shift()
      }
      this.optionCpu.xAxis.data = [...this.cpuHistoryLabels]
      this.optionCpu.series[0].data = [...this.cpuHistoryValues]
      if (this.echartsView) {
        this.echartsView.setOption(this.optionCpu, true)
      }
    },
    updateMemChart(value) {
      const memValue = Number(parseFloat(value).toFixed(2))
      this.memHistoryLabels.push(this.getCurrentTimeLabel())
      this.memHistoryValues.push(memValue)
      if (this.memHistoryLabels.length > this.memHistoryLimit) {
        this.memHistoryLabels.shift()
        this.memHistoryValues.shift()
      }
      this.optionMem.xAxis.data = [...this.memHistoryLabels]
      this.optionMem.series[0].data = [...this.memHistoryValues]
      if (this.echartsViewMem) {
        this.echartsViewMem.setOption(this.optionMem, true)
      }
    },
    startAutoRefresh() {
      this.stopAutoRefresh()
      this.refreshTimer = setInterval(() => {
        this.GetSystemAnalysis()
      }, 3000)
    },
    stopAutoRefresh() {
      if (this.refreshTimer) {
        clearInterval(this.refreshTimer)
        this.refreshTimer = null
      }
    },
    resizeCharts() {
      this.$nextTick(() => {
        if (this.echartsView) {
          this.echartsView.resize()
        }
        if (this.echartsViewMem) {
          this.echartsViewMem.resize()
        }
        if (this.echartsViewDisk) {
          this.echartsViewDisk.resize()
        }
      })
    },
    GetSystemAnalysis(showLoading = false){
      let _t = this
      if (showLoading) {
        _t.OptLoading = true
      }
      GetSystemAnalysis().then(function (res){
        if (showLoading) {
          _t.OptLoading = false
        }
        // 合并 API 数据到默认值，防止 null 字段导致模板渲染崩溃
        const apiData = res.data.list || {}
        if (apiData.CpuInfo) {
          Object.assign(_t.SystemInfo.CpuInfo, apiData.CpuInfo)
          _t.SystemInfo.CpuInfo.UsedPercent[0] = parseFloat(_t.SystemInfo.CpuInfo.UsedPercent[0]).toFixed(2)
          _t.SystemInfo.CpuInfo.Load.load1 = parseFloat(_t.SystemInfo.CpuInfo.Load.load1).toFixed(2)
          _t.SystemInfo.CpuInfo.Load.load5 = parseFloat(_t.SystemInfo.CpuInfo.Load.load5).toFixed(2)
          _t.SystemInfo.CpuInfo.Load.load15 = parseFloat(_t.SystemInfo.CpuInfo.Load.load15).toFixed(2)
        }
        if (apiData.DiskInfo) {
          Object.assign(_t.SystemInfo.DiskInfo, apiData.DiskInfo)
          _t.SystemInfo.DiskInfo.total = (parseFloat(_t.SystemInfo.DiskInfo.total)/(1024*1024*1024)).toFixed(0)
          _t.SystemInfo.DiskInfo.free = (parseFloat(_t.SystemInfo.DiskInfo.free)/(1024*1024*1024)).toFixed(0)
          _t.SystemInfo.DiskInfo.used = (parseFloat(_t.SystemInfo.DiskInfo.used)/(1024*1024*1024)).toFixed(0)
          _t.SystemInfo.DiskInfo.usedPercent = (parseFloat(_t.SystemInfo.DiskInfo.usedPercent)).toFixed(0)
        }
        if (apiData.MemInfo) {
          Object.assign(_t.SystemInfo.MemInfo, apiData.MemInfo)
          _t.SystemInfo.MemInfo.total = (parseFloat(_t.SystemInfo.MemInfo.total)/(1024*1024*1024)).toFixed(0)
          _t.SystemInfo.MemInfo.used = (parseFloat(_t.SystemInfo.MemInfo.used)/(1024*1024*1024)).toFixed(0)
          _t.SystemInfo.MemInfo.free = (parseFloat(_t.SystemInfo.MemInfo.free)/(1024*1024*1024)).toFixed(0)
          _t.SystemInfo.MemInfo.goUsed = (parseFloat(_t.SystemInfo.MemInfo.goUsed)/(1024*1024)).toFixed(0)
          _t.SystemInfo.MemInfo.usedPercent = (parseFloat(_t.SystemInfo.MemInfo.usedPercent)).toFixed(0)
        }
        if (apiData.NetWork) {
          Object.assign(_t.SystemInfo.NetWork, apiData.NetWork)
          _t.SystemInfo.NetWork.Receive = (parseFloat(_t.SystemInfo.NetWork.Receive)/(1024*1024)).toFixed(0)
          _t.SystemInfo.NetWork.Sent = (parseFloat(_t.SystemInfo.NetWork.Sent)/(1024*1024)).toFixed(0)
        }
        if (apiData.HostInfo) {
          Object.assign(_t.SystemInfo.HostInfo, apiData.HostInfo)
          _t.SystemInfo.HostInfo.uptime = (parseInt(_t.SystemInfo.HostInfo.uptime)/(3600)).toFixed(0)
          _t.SystemInfo.HostInfo.bootTime = _t.SystemInfo.BootTime
        }
        if (apiData.Goroutine != null) _t.SystemInfo.Goroutine = apiData.Goroutine
        if (apiData.BootTime) _t.SystemInfo.BootTime = apiData.BootTime

        _t.updateCpuChart(_t.SystemInfo.CpuInfo.UsedPercent[0])
        _t.updateMemChart(_t.SystemInfo.MemInfo.usedPercent)

        _t.optionDisk.series[0].data[0] = _t.SystemInfo.DiskInfo.usedPercent
        _t.optionDisk.title[0].text = _t.SystemInfo.DiskInfo.usedPercent+ "%"
        if (_t.echartsViewDisk) {
          _t.echartsViewDisk.setOption(_t.optionDisk,true)
        }
        _t.resizeCharts()
      }).catch(function () {
        if (showLoading) {
          _t.OptLoading = false
        }
      })
    },
    initMap(){
      if (this.echartsView) {
        this.echartsView.dispose()
      }
      if (this.echartsViewMem) {
        this.echartsViewMem.dispose()
      }
      if (this.echartsViewDisk) {
        this.echartsViewDisk.dispose()
      }
      let view = this.$refs.cpuCharts
      let viewMem = this.$refs.memCharts
      let viewDisk = this.$refs.viewDisk

      if (view && view.clientWidth > 0 && view.clientHeight > 0) {
        this.echartsView = echarts.init(view, null)
        this.echartsView.setOption(this.optionCpu, true)
      }

      if (viewMem && viewMem.clientWidth > 0 && viewMem.clientHeight > 0) {
        this.echartsViewMem = echarts.init(viewMem, null)
        this.echartsViewMem.setOption(this.optionMem, true)
      }

      if (viewDisk && viewDisk.clientWidth > 0 && viewDisk.clientHeight > 0) {
        this.echartsViewDisk = echarts.init(viewDisk, null)
        this.echartsViewDisk.setOption(this.optionDisk, true)
      }

      this.resizeCharts()
      setTimeout(() => {
        this.resizeCharts()
      }, 100)
    },
    GetSystemParams(){
      let _t = this
      GetSystemParams().then(function (res) {
        _t.$store.state.setting.IsOEM = res.data.list.IsOEM
        _t.IsLicense = _t.$store.state.setting.IsOEM
        if(!_t.IsLicense)
        {
          watermark.set("www.ismctl.com","免费试用版")
        }
      }).catch(function(e){

      })
    }
  },
  mounted(){
    this.$nextTick(() => {
      this.initMap()
      this.GetSystemAnalysis(true)
      this.startAutoRefresh()
    })
  },
  created() {
    this.GetSystemParams()
  },
  deactivated() {
    this.stopAutoRefresh()
  },
  activated() {
    this.startAutoRefresh()
    if (!this.echartsView || !this.echartsViewMem || !this.echartsViewDisk) {
      this.$nextTick(() => this.initMap())
    }
  },
  beforeDestroy() {
    this.stopAutoRefresh()
    if (this.echartsView) {
      this.echartsView.dispose()
      this.echartsView = null
    }
    if (this.echartsViewMem) {
      this.echartsViewMem.dispose()
      this.echartsViewMem = null
    }
    if (this.echartsViewDisk) {
      this.echartsViewDisk.dispose()
      this.echartsViewDisk = null
    }
  },
}
</script>

<style lang="less" scoped>

:root {
  --panel-bg: rgba(255, 255, 255, 0.82);
  --panel-border: rgba(148, 163, 184, 0.16);
  --panel-shadow: 0 18px 40px rgba(15, 23, 42, 0.08);
  --panel-shadow-hover: 0 22px 48px rgba(15, 23, 42, 0.12);
  --text-strong: #132238;
  --text-muted: #63758c;
}

// 核心指标卡片样式
.metrics-card-row {
  margin-top: 0 !important;
}

.metrics-card {
  position: relative;
  border-radius: 18px;
  transition: transform 0.28s ease, box-shadow 0.28s ease, border-color 0.28s ease;
  overflow: hidden;
  border: 1px solid var(--panel-border);
  box-shadow: var(--panel-shadow);
  height: 100%;
  backdrop-filter: blur(14px);
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.92) 0%, rgba(248, 251, 255, 0.88) 100%);

  &:hover {
    transform: translateY(-4px);
    box-shadow: var(--panel-shadow-hover);
    border-color: rgba(59, 130, 246, 0.12);
  }

  & /deep/ .ant-card-body {
    padding: 14px 12px;
    text-align: center;
  }
}
.metrics-card-l {
  position: relative;
  border-radius: 18px;
  overflow: hidden;
  border: 1px solid var(--panel-border);
  box-shadow: var(--panel-shadow);
  height: 100%;
  backdrop-filter: blur(14px);
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.92) 0%, rgba(248, 251, 255, 0.88) 100%);
  transition: transform 0.28s ease, box-shadow 0.28s ease, border-color 0.28s ease;
  & /deep/ .ant-card-body {
    padding: 10px 8px 12px;
  }
  & /deep/ .ant-card-head {
    min-height: 44px;
    border-bottom: 1px solid rgba(148, 163, 184, 0.1);
    padding: 0 18px;
  }
  & /deep/ .ant-card-head-title {
    padding: 12px 0 10px;
    color: var(--text-strong);
    font-size: 18px;
    font-weight: 700;
    letter-spacing: 0.02em;
  }
  &:hover {
    transform: translateY(-4px);
    box-shadow: var(--panel-shadow-hover);
  }
}
.resource-row {
  margin-top: 18px;
}
.resource-card-body {
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: 12px;
  height: auto;
  min-height: 258px;
}
.resource-stats {
  display: block;
  min-width: 0;
  flex: 0 0 auto;
  height: auto;
  padding-left: 18px;
}
.resource-form {
  width: 100%;
}
.resource-form-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  column-gap: 18px;
  row-gap: 0;
  grid-auto-rows: minmax(28px, auto);
  width: 100%;
  align-content: start;
  height: auto;
}
.resource-stat-placeholder {
  visibility: hidden;
  height: 28px;
}
.resource-chart {
  display: flex;
  align-items: flex-end;
  justify-content: center;
  min-width: 0;
  flex: 0 0 190px;
  height: 190px;
  width: 100%;
  border-radius: 14px;
  background: linear-gradient(180deg, rgba(246, 250, 255, 0.78) 0%, rgba(252, 254, 255, 0.42) 100%);
}
.cpu-chart-wrap,
.mem-chart-wrap {
  align-items: flex-end;
  justify-content: center;
}
.disk-chart-wrap {
  align-items: center;
  justify-content: center;
}
.cpu-chart-wrap .view-chart-gauge,
.mem-chart-wrap .view-chart-gauge {
  max-width: none;
  width: 100%;
  height: 190px;
}
.disk-chart-wrap .view-chart-gauge {
  width: 100%;
  max-width: 240px;
  height: 190px;
}
.resource-form ::v-deep .ant-form-item {
  display: flex;
  align-items: center;
  margin-bottom: 0;
  min-height: 28px;
}
.resource-form ::v-deep .ant-form-item:last-child {
  margin-bottom: 0;
}
.resource-form ::v-deep .ant-form-item-label {
  text-align: left;
  line-height: 1.3;
  white-space: nowrap;
  flex: 0 0 auto;
  min-width: auto;
  padding-right: 4px;
}
.resource-form ::v-deep .ant-form-item-control-wrapper {
  flex: 1 1 auto;
  min-width: 0;
  white-space: nowrap;
  overflow: visible;
}
.resource-form ::v-deep .ant-form-text {
  white-space: nowrap;
}
.resource-form ::v-deep .ant-form-item-label > label {
  color: var(--text-muted);
  font-size: 12px;
  font-weight: 600;
  letter-spacing: 0.01em;
}
.resource-form ::v-deep .ant-form-item-control {
  line-height: 1.3;
}
.resource-form .ant-form-text {
  color: var(--text-strong);
  font-variant-numeric: tabular-nums;
  font-size: 14px;
  font-weight: 700;
}
.border-info {
  border-top: 4px solid #2b9fd8;
  background:
    radial-gradient(circle at top left, rgba(43, 159, 216, 0.09) 0%, rgba(255, 255, 255, 0) 40%),
    linear-gradient(180deg, rgba(255, 255, 255, 0.95) 0%, rgba(248, 251, 255, 0.9) 100%);
}

.border-danger {
  border-top: 4px solid #48cdb7;
  background:
    radial-gradient(circle at top left, rgba(72, 205, 183, 0.09) 0%, rgba(255, 255, 255, 0) 40%),
    linear-gradient(180deg, rgba(255, 255, 255, 0.95) 0%, rgba(246, 253, 252, 0.9) 100%);
}

.border-success {
  border-top: 4px solid #78bf4d;
  background:
    radial-gradient(circle at top left, rgba(120, 191, 77, 0.09) 0%, rgba(255, 255, 255, 0) 40%),
    linear-gradient(180deg, rgba(255, 255, 255, 0.95) 0%, rgba(249, 253, 246, 0.9) 100%);
}

.metrics-label {
  color: var(--text-muted);
  font-size: 13px;
  margin-bottom: 14px;
  font-weight: 700;
  letter-spacing: 0.04em;
  text-transform: uppercase;
}

.metrics-value-wrap {
  display: flex;
  justify-content: center;
  align-items: baseline;
}

.metrics-value {
  color: var(--text-strong);
  font-size: 32px;
  font-weight: 700;
  line-height: 1;
}

.metrics-unit {
  color: var(--text-muted);
  font-size: 13px;
  margin-left: 6px;
  font-weight: 600;
}

.metrics-icon {
  position: absolute;
  top: 14px;
  right: 14px;
  color: inherit;
  opacity: 0.5;
  transition: opacity 0.28s ease, transform 0.28s ease;
  &:hover {
    opacity: 1;
    transform: scale(1.08);
  }
}
#WaterMarkShow{
  position: relative;
  z-index: 99;
  overflow: hidden;
  padding: 18px;
  background:
    radial-gradient(circle at top left, rgba(18, 151, 214, 0.1) 0%, rgba(245, 248, 252, 0) 26%),
    radial-gradient(circle at top right, rgba(214, 28, 115, 0.08) 0%, rgba(245, 248, 252, 0) 28%),
    linear-gradient(180deg, #f5f8fc 0%, #eef3f8 100%);
  min-height: 100vh;
}
::v-deep .ant-form-item {
  margin-bottom: 0px;
}
::v-deep .active-ring-chart-container .active-ring-name{
  fontSize: 25px;
  color:#171515;
}
.extra-wrap{
  .extra-item{
    display: inline-block;
    margin-right: 24px;
    a:not(:first-child){
      margin-left: 24px;
    }
  }
}
.view-chart-gauge {
  width: 100%;
  max-width: none;
  height: 190px;
}
@media screen and (max-width: 992px){
  .resource-card-body {
    height: auto;
    gap: 10px;
  }
  .resource-stats {
    flex: 0 0 auto;
    height: auto;
  }
  .resource-form-grid {
    height: auto;
  }
  .resource-chart {
    flex-basis: 200px;
    height: 200px;
  }
  .resource-chart,
  .cpu-chart-wrap,
  .mem-chart-wrap,
  .disk-chart-wrap {
    justify-content: center;
  }
  .view-chart-gauge {
    max-width: none;
    height: 200px;
  }
  .extra-wrap .extra-item{
    display: none;
  }
}
@media screen and (max-width: 576px){
  .metrics-card-l {
    & /deep/ .ant-card-body {
      padding: 10px 12px;
    }
  }
  .resource-form ::v-deep .ant-form-item {
    margin-bottom: 8px;
  }
  .resource-form-grid {
    grid-template-columns: 1fr;
    grid-auto-rows: 26px;
    row-gap: 0;
  }
  .resource-chart {
    flex-basis: 180px;
    height: 180px;
  }
  .resource-form ::v-deep .ant-form-item-label,
  .resource-form ::v-deep .ant-form-item-control-wrapper {
    flex: 0 0 50%;
    max-width: 50%;
  }
  .view-chart-gauge {
    height: 180px;
  }
  .extra-wrap{
    display: none;
  }
}
.text-secondary {
  color: #6c757d!important;
  font-size: 16px;
  margin-top: 10px;
  text-align: center;
  font-family: "黑体";
}
.text-content {
  font-size: 29px;
  text-align: center;
  font-family: "黑体";
  margin-bottom: 5px;
}
::v-deep .ant-card{
  border-radius: 18px;
}
.border-success {
  border-top: 4px solid #78bf4d;
}
.border-danger {
  border-top: 4px solid  #48cdb7;
}
.border-info {
  border-top: 4px solid #2b9fd8;
}
::v-deep .ant-card-body {
  padding: 0px;
}
</style>