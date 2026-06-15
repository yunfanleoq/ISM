<template>
  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"  x="0px" y="0px"  xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <foreignObject class="history-trend-foreign-object" pointer-events="all" :width="detail.style.position.w" :height="detail.style.position.h">
        <div class="history-trend-layout">
          <div :class="['history-trend-theme-root', formThemeClass]" :key="`theme-root-${EchartsTheme}`" ref="themeRoot">
            <a-form layout="horizontal" :class="formThemeClass" :style="formThemeStyle" :key="`form-${EchartsTheme}`">
              <div>
                <a-row >
                  <a-col :md="8" :sm="24" >
                    <a-form-item
                        :label="$t('reporting.AlarmHistory.SelectDate')"
                        :labelCol="{span: 8}"
                        :wrapperCol="{span: 16}"
                        style="margin-bottom: 0;"
                    >
                      <div :class="['date-picker-row', datePickerThemeClass]">
                        <a-date-picker
                          :key="`date-picker-${EchartsTheme}`"
                          :defaultValue="moment()"
                          style="width: 100%"
                          @change="onDateChange"
                          size="default"
                          :placeholder="$t('reporting.AlarmHistory.DateDayType')"
                          :class="datePickerThemeClass"
                          :dropdownClassName="datePickerDropdownClass"
                          :getCalendarContainer="getDatePickerPopupContainer"
                        />
                      </div>
                    </a-form-item>
                  </a-col>
                </a-row>
              </div>
            </a-form>
          </div>
          <div class="view-chart-real-data" :ref="detail.identifier"></div>
        </div>
      </foreignObject>
      <!--      闪烁-->
      <animate v-if="isStart&&animateType.includes('blink')&&!IsToolBox" attributeName="opacity"
               values="0.1;1;0.1" :dur="blinkSpeed+'s'"
               repeatCount="indefinite"/>
      <!--渐变-->
      <animate v-if="isStart&&animateType.includes('millcolorGrad')&&!IsToolBox" attributeName="fill"
               :values="startColor+';'+stopColor+';'+startColor" :dur="animateSpeed+'s'"
               repeatCount="indefinite"/>
      <!--缩放      -->
      <animateTransform v-if="isStart&&animateType.includes('Zoom')&&!IsToolBox" attributeName="transform"   begin="0s" dur="0.6s" type="scale" values="0.9;1;0.9" repeatCount="indefinite"/>
      <!--      顺时针旋转-->
      <animateTransform v-if="isStart&&animateType.includes('animateSpin')&&!IsToolBox&&spinDirection==0" attributeType="XML" attributeName="transform" :dur="animateSpinSpeed+'s'" type="rotate" from="0 0 0" to="360 0 0" repeatCount="indefinite" />
      <!--      逆时针旋转-->
      <animateTransform v-if="isStart&&animateType.includes('animateSpin')&&!IsToolBox&&spinDirection==1" attributeType="XML" attributeName="transform" :dur="animateSpinSpeed+'s'" type="rotate" from="360 0 0" to="0 0 0" repeatCount="indefinite" />
  </g>
</svg>

</template>

<script>
import * as echarts from 'echarts'
import BaseView from '../View';
import moment from 'moment/moment'
import chalk from '@/theme/echarts/chalk'
import essos from '@/theme/echarts/essos'
import dark from '@/theme/echarts/dark'
import infographic from '@/theme/echarts/infographic'
import macarons from '@/theme/echarts/macarons'
import roma from '@/theme/echarts/roma'
import shine from '@/theme/echarts/shine'
import vintage from '@/theme/echarts/vintage'
import purplePassion from '@/theme/echarts/purple-passion'
import walden from '@/theme/echarts/walden'
import westeros from '@/theme/echarts/westeros'
import wonderland from '@/theme/echarts/wonderland'
import {GetChartDataHistoryTrendByDate} from "@/services/report";
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
export default {
  mixins: [ISMChildAutoMixin],
  name: 'ism-view-history-trend-chart-by-date',
  i18n: require('@/i18n/language'),
  inject: ['getNode'],
  props: {

  },
  watch: {
    detail: {
      handler(newVal, oldVal) {
        if(this.editMode){
          this.initComponents(newVal)
        }
        this.onResize()
      },
      deep: true
    },
    // 监听主题变化，强制更新视图
    EchartsTheme: {
      handler(newTheme, oldTheme) {
        if (newTheme !== oldTheme) {
          // 强制更新视图，确保日期选择器样式同步更新
          this.$nextTick(() => {
            this.onResize()
            this.syncThemeStyles()
          })
        }
      }
    }
  },
  data() {
    return {
      moment,
      detail:{},
      IsToolBox:false,
      editMode:true,
      strokeColor:"#000000",
      fill:"#A1BFE2",
      strokeWidth:0.3,
      ChartTimelyRefreshTimer:null,
      ChartTimelyRefresh:5*60*1000,
      fillOpacity:1,
      TimeInterval:60,
      strokeOpacity:1,
      animateType:"blink",
      startColor:"#74f808",
      stopColor:"#74f808",
      animateSpeed:0.5,
      animateSpinSpeed:0.5,
      ValueType:5,
      spinDirection:0,
      blinkSpeed:0.5,
      isStart:false,
      IsShowDate:0,
      // 深色主题列表
      darkThemes: ['dark', 'chalk', 'purplePassion', 'westeros', 'wonderland'],
      // 浅色主题列表
      lightThemes: ['essos', 'infographic', 'macarons', 'roma', 'shine', 'vintage', 'walden'],
      base:{
        "text": "configComponent.RealDataChart.HistoryTrendDateTitle",
        "icon": "icon-qushitu",
        "isFontIcon": true,
        "info": {
          "type": "real-data-chart",
          "action": [],
          "active": [
            {
              id:"ShowChartVariable1",
              name:"configComponent.variable.ShowData",
              result:"",
              isExpression:false,
              condition:{
                deviceSN:"",
                selectVideoType:0,
                isBandDevice:false,
                bandType:1,
                dataID: "",
                dataName: "",
                operator:"",
                OperatorValue:"",
                OperatorMaxValue:"",
              },
            },
            {
              id:"ShowChartVariable2",
              name:"configComponent.variable.ShowData",
              result:"",
              isExpression:false,
              condition:{
                deviceSN:"",
                selectVideoType:0,
                isBandDevice:false,
                bandType:1,
                dataID: "",
                dataName: "",
                operator:"",
                OperatorValue:"",
                OperatorMaxValue:"",
              },
            },
            {
              id:"ShowChartVariable3",
              name:"configComponent.variable.ShowData",
              result:"",
              isExpression:false,
              condition:{
                deviceSN:"",
                selectVideoType:0,
                isBandDevice:false,
                bandType:1,
                dataID: "",
                dataName: "",
                operator:"",
                OperatorValue:"",
                OperatorMaxValue:"",
              },
            },
            {
              id:"ShowChartVariable4",
              name:"configComponent.variable.ShowData",
              result:"",
              isExpression:false,
              condition:{
                deviceSN:"",
                selectVideoType:0,
                isBandDevice:false,
                bandType:1,
                dataID: "",
                dataName: "",
                operator:"",
                OperatorValue:"",
                OperatorMaxValue:"",
              },
            },
            {
              id:"ShowChartVariable5",
              name:"configComponent.variable.ShowData",
              result:"",
              isExpression:false,
              condition:{
                deviceSN:"",
                selectVideoType:0,
                isBandDevice:false,
                bandType:1,
                dataID: "",
                dataName: "",
                operator:"",
                OperatorValue:"",
                OperatorMaxValue:"",
              },
            }
          ],
          "animate": {
            "selected": [],
            "condition":{
              deviceSN:"",
              selectVideoType:0,
              isBandDevice:false,
              bandType:1,
              dataID: "",
              dataName: "",
              operator:"",
              OperatorValue:"",
              OperatorMaxValue:"",
            },
            "isExpression": false,
            "animateList": [
              {
                id: "Forbidden",
                name: "component.public.Forbidden",
              },
              {
                id: "blink",
                name: "component.public.animateBlink",
              },
              {
                id: "Zoom",
                name: "component.public.Zoom",
              },
              {
                id: "animateSpin",
                name: "component.public.animateSpin",
              },
            ],
            "animateElement": [
              {
                id: "blink",
                elementList:[
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"blinkSpeed",
                  },
                ]
              },
              {
                id: "millcolorGrad",
                elementList:[
                  {
                    "name": "component.public.startColor",
                    "type": 2,
                    "value": "#74f808",
                    "key": "startColor",
                  },
                  {
                    "name": "component.public.stopColor",
                    "type": 2,
                    "value": "#f30b0b",
                    "key": "stopColor",
                  },
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"animateSpeed",
                  },
                ]
              },
              {
                id: "animateSpin",
                elementList:[
                  {
                    "name":"component.public.animateSpinSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"spinSpeed",
                  },
                  {
                    name:"configComponent.bigScreen.border.border89Direction",
                    type:6,
                    value:0,
                    enumList:[
                      {
                        value:0,
                        option:"configComponent.bigScreen.border.border89DirectionForward"
                      },
                      {
                        value:1,
                        option:"configComponent.bigScreen.border.border89DirectionNegative"
                      }
                    ],
                    min:1,
                    key:"spinDirection",
                  }
                ]
              },
            ],
          },
          "style": {
            "position": {
              "x": 0,
              "y": 0,
              "w": 700,
              "h": 430
            },
            "backColor": "transparent",
            "foreColor": "#ffffff",
            "fontSize": 14,
            fontFamily: "Arial",
            "zIndex": 1,
            "transform": 0,
            "diy":[
              {
                "name":"configComponent.ChartPublic.ChartTitle",
                "type":4,
                "value":"历史趋势",
                "key":"ChartTitle",
              },
              {
                "name":"configComponent.ChartPublic.TimelyInitEcharts",
                "type":1,
                "value":60,
                "key":"TimelyInitEcharts",
              },
              {
                "name":"configComponent.ChartPublic.TimeInterval",
                "type":1,
                "value":60,
                "key":"TimeInterval",
              },
              {
                "name":"configComponent.ChartPublic.TimelyRefresh",
                "type":1,
                "value":5,
                "key":"TimelyRefresh",
              },
              {
                "name":"configComponent.ChartPublic.ChartRange",
                type:6,
                value:5,
                enumList:[
                  {
                    value:0,
                    option:"configComponent.ChartPublic.PeriodList.oneHour"
                  },
                  {
                    value:1,
                    option:"configComponent.ChartPublic.PeriodList.FiveHour"
                  },
                  {
                    value:2,
                    option:"configComponent.ChartPublic.PeriodList.FivethHour"
                  },
                  {
                    value:3,
                    option:"configComponent.ChartPublic.PeriodList.oneDay"
                  },
                  {
                    value:4,
                    option:"configComponent.ChartPublic.PeriodList.threeDay"
                  },
                  {
                    value:5,
                    option:"configComponent.ChartPublic.PeriodList.sevenDay"
                  },
                  {
                    value:6,
                    option:"configComponent.ChartPublic.PeriodList.fivethDay"
                  },
                  {
                    value:7,
                    option:"configComponent.ChartPublic.PeriodList.onemonth"
                  },
                ],
                "key":"ChartRange",
              },
              {
                "name":"configComponent.ChartPublic.YMax",
                "type":1,
                "value":0,
                "key":"YMax",
              },
              {
                "name":"configComponent.ChartPublic.YMin",
                "type":1,
                "value":0,
                "key":"YMin",
              },
              {
                "name":"configComponent.ChartPublic.IsShowDate",
                type:6,
                value:0,
                enumList:[
                  {
                    value:0,
                    option:"False"
                  },
                  {
                    value:1,
                    option:"True"
                  },
                ],
                "key":"IsShowDate",
              },
              {
                "name":"configComponent.ChartPublic.EchartsWidth",
                "type":1,
                "value":2,
                "key":"EchartsWidth",
              },
              {
                "name":"configComponent.ChartPublic.EchartsXRotate",
                "type":1,
                "value":0,
                "key":"EchartsXRotate",
              },
              {
                "name":"configComponent.ChartPublic.EchartsXFormat",
                type:4,
                value:"YYYY-MM-DD HH:mm:ss",
                "key":"EchartsXFormat",
              },
              {
                "name":"configComponent.ChartPublic.EchartsXTheme",
                type:6,
                value:"dark",
                enumList:[
                  {
                    value:"chalk",
                    option:"chalk"
                  },
                  {
                    value:"essos",
                    option:"essos"
                  },
                  {
                    value:"dark",
                    option:"dark"
                  },
                  {
                    value:"infographic",
                    option:"infographic"
                  },
                  {
                    value:"macarons",
                    option:"macarons"
                  },
                  {
                    value:"roma",
                    option:"roma"
                  },
                  {
                    value:"shine",
                    option:"shine"
                  },
                  {
                    value:"vintage",
                    option:"vintage"
                  },
                  {
                    value:"purplePassion",
                    option:"purplePassion"
                  },
                  {
                    value:"walden",
                    option:"walden"
                  },
                  {
                    value:"westeros",
                    option:"westeros"
                  },
                  {
                    value:"wonderland",
                    option:"wonderland"
                  }
                ],
                "key":"EchartsXTheme",
              },
              {
                "name":"configComponent.ChartPublic.IsShowZoom",
                type:6,
                value:1,
                enumList:[
                  {
                    value:0,
                    option:"False"
                  },
                  {
                    value:1,
                    option:"True"
                  },
                ],
                "key":"IsShowZoom",
              },
            ]
          }
        }
      },
      EchartsXFormat:1,
      EchartsTheme:"dark",
      SelectDateRange: moment().format("YYYY-MM-DD"),
      date: [],
      yieldRate: [],
      EchartsWidth:1,
      yieldIndex: [],
      timelySave:null,
      getList:[],
      option : {
        title: {
          text: "",
          textStyle:{
            color:"",
            fontFamily:"",
            fontSize:"",
          }
        },
        tooltip: {
          trigger: 'axis',
          confine: true
        },
        legend: {
          data: ['Email', 'Union Ads', 'Video Ads', 'Direct', 'Search Engine'],
          textStyle:{
            color:"",
            fontFamily:"",
            fontSize:"",
          }
        },
        grid: {
          left: '1%',
          right: '1%',
          bottom: '10%',
          containLabel: true
        },
        toolbox: {
          feature: {
            dataZoom: {
              yAxisIndex: 'none'
            },
            restore: {},
            saveAsImage: {}
          }
        },
        xAxis: {
          type: 'time',
          min: moment().startOf('day').valueOf(),
          max: moment().endOf('day').valueOf(),
          interval: 4 * 60 * 60 * 1000,
          boundaryGap: false,
          axisLine: {
            show: false, //是否显示坐标刻度
            lineStyle: {
              color: '#eeeeee'
            }
          },
          axisLabel: {
            show: true,
            rotate:70,
            margin: 15,
            formatter: function (value) {
              return moment(value).format('HH:mm')
            } , // 仅显示时分
            textStyle: {
              color: '#fff'
            }
          },
          data: []
        },
        yAxis: {
          type: 'value',
          min:100,
          max:1000,
          boundaryGap: [0, '100%'],
          axisLabel: {
            show: true,
            textStyle: {
              color: '#fff'
            }
          },
        },
        dataZoom: [
          {
            type: 'inside',
            xAxisIndex: 0,
            start: 0,
            end: 100,
            zoomLock: false
          },
          {
            type: 'slider',
            xAxisIndex: 0,
          }
        ],
        series: [
          {
            name: 'Email',
            type: 'line',
            stack: 'Total',
            data: [['2025/04/29 13:38:30',120], ['2025/04/29 13:38:40',132], ['2025/04/29 13:38:50',101], ['2025/04/29 13:38:51',134], ['2025/04/29 13:38:52',90], ['2025/04/29 13:38:53',230], ['2025/04/29 13:38:56',210]]
          },
          {
            name: 'Union Ads',
            type: 'line',
            stack: 'Total',
            data: [['2025/04/29 13:38:30',120], ['2025/04/29 13:38:40',132], ['2025/04/29 13:38:50',101], ['2025/04/29 13:38:51',134], ['2025/04/29 13:38:52',90], ['2025/04/29 13:38:53',230], ['2025/04/29 13:38:56',210]]
          },
          {
            name: 'Video Ads',
            type: 'line',
            stack: 'Total',
            data: [['2025/04/29 13:38:30',120], ['2025/04/29 13:38:40',132], ['2025/04/29 13:38:50',101], ['2025/04/29 13:38:51',134], ['2025/04/29 13:38:52',90], ['2025/04/29 13:38:53',230], ['2025/04/29 13:38:56',210]]
          },
          {
            name: 'Direct',
            type: 'line',
            stack: 'Total',
            data: [['2025/04/29 13:38:30',120], ['2025/04/29 13:38:40',132], ['2025/04/29 13:38:50',101], ['2025/04/29 13:38:51',134], ['2025/04/29 13:38:52',90], ['2025/04/29 13:38:53',230], ['2025/04/29 13:38:56',210]]
          },
          {
            name: 'Search Engine',
            type: 'line',
            stack: 'Total',
            data: [['2025/04/29 13:38:30',120], ['2025/04/29 13:38:40',132], ['2025/04/29 13:38:50',101], ['2025/04/29 13:38:51',134], ['2025/04/29 13:38:52',90], ['2025/04/29 13:38:53',230], ['2025/04/29 13:38:56',210]]
          }
        ]
      },
      TimelyInitEcharts:60,
      seriesMap:[]
    }
  },
    computed: {
      // 判断当前主题是否为深色主题
      isDarkTheme() {
        return this.darkThemes.includes(this.EchartsTheme)
      },
      // 表单主题样式类
      formThemeClass() {
        return this.isDarkTheme ? 'form-dark-theme' : 'form-light-theme'
      },
      // 日期选择器主题样式类
      datePickerThemeClass() {
        return this.isDarkTheme ? 'date-picker-dark' : 'date-picker-light'
      },
      datePickerDropdownClass() {
        return this.isDarkTheme ? 'date-picker-dropdown-dark' : 'date-picker-dropdown-light'
      },
      chartThemeBackground() {
        const backgroundMap = {
          chalk: 'rgba(41,52,65,1)',
          dark: '#100C2A',
          essos: 'rgba(242,234,191,0.15)',
          infographic: 'rgba(0,0,0,0)',
          macarons: 'rgba(0,0,0,0)',
          purplePassion: 'rgba(91,92,110,1)',
          roma: 'rgba(0,0,0,0)',
          shine: 'rgba(0,0,0,0)',
          vintage: '#fef8ef',
          walden: 'rgba(252,252,252,0)',
          westeros: 'rgba(0,0,0,0)',
          wonderland: 'rgba(255,255,255,0)'
        }
        return backgroundMap[this.EchartsTheme] || 'rgba(0,0,0,0)'
      },
      // 表单主题内联样式
      formThemeStyle() {
        return {
          backgroundColor: this.chartThemeBackground,
          padding: '4px 8px',
          borderRadius: '0'
        }
      }
    },
    methods: {
      waitChartContainerReady(view, callback, retryCount = 0) {
        if (!view) {
          return
        }
        if (view.clientWidth > 0 && view.clientHeight > 0) {
          callback()
          return
        }
        if (retryCount >= 10) {
          callback()
          return
        }
        const raf = typeof requestAnimationFrame === 'function' ? requestAnimationFrame : (fn) => setTimeout(fn, 16)
        raf(() => this.waitChartContainerReady(view, callback, retryCount + 1))
      },
      getDatePickerPopupContainer(triggerNode) {
        if (typeof document !== 'undefined' && document.body) {
          return document.body
        }
        return triggerNode && triggerNode.parentNode ? triggerNode.parentNode : triggerNode
      },
      applyFixedTimeAxis() {
        const selectedDate = this.SelectDateRange || moment().format("YYYY-MM-DD")
        const dayStart = moment(selectedDate + ' 00:00:00', 'YYYY-MM-DD HH:mm:ss')
        const dayEnd = moment(selectedDate + ' 23:59:59', 'YYYY-MM-DD HH:mm:ss')

        this.option.xAxis.type = 'time'
        this.option.xAxis.min = dayStart.valueOf()
        this.option.xAxis.max = dayEnd.valueOf()
        this.option.xAxis.interval = 4 * 60 * 60 * 1000
        this.option.xAxis.axisLabel.formatter = function (value) {
          return moment(value).format('HH:mm')
        }
      },
      getTooltipPosition(point, params, dom, rect, size) {
        const contentSize = size && size.contentSize ? size.contentSize : [0, 0]
        const viewSize = size && size.viewSize ? size.viewSize : [0, 0]
        const hasSliderZoom = Array.isArray(this.option.dataZoom) && this.option.dataZoom.some(item => item.type === 'slider')
        const bottomReserve = hasSliderZoom ? 48 : 12
        const gap = 12

        let x = point[0] + gap
        let y = point[1] - contentSize[1] - gap

        const maxX = Math.max(viewSize[0] - contentSize[0] - gap, gap)
        const maxY = Math.max(viewSize[1] - contentSize[1] - bottomReserve, gap)

        if (x > maxX) {
          x = point[0] - contentSize[0] - gap
        }
        if (x < gap) {
          x = gap
        }

        if (y < gap) {
          y = point[1] + gap
        }
        if (y > maxY) {
          y = maxY
        }
        if (y < gap) {
          y = gap
        }

        return [x, y]
      },
      onDateChange(date, dateString){
        this.SelectDateRange = dateString
        this.applyFixedTimeAxis()
        this.timelyTimeGetHistory()
      },
    refreshChart(){
      if(this.echartsView==null)
      {
        return
      }
      this.date = []
      this.initComponents(this.detail)
    },
    syncThemeStyles() {
      const root = this.$refs.themeRoot
      if (!root) {
        return
      }

      const isDark = this.isDarkTheme
      const labelColor = isDark ? 'rgba(255, 255, 255, 0.85)' : 'rgba(0, 0, 0, 0.85)'
      const inputBg = isDark ? 'rgba(40, 40, 60, 0.8)' : '#fff'
      const inputBorder = isDark ? 'rgba(255, 255, 255, 0.2)' : '#d9d9d9'
      const iconColor = isDark ? 'rgba(255, 255, 255, 0.65)' : 'rgba(0, 0, 0, 0.45)'
      const panelBg = isDark ? 'rgba(20, 20, 30, 0.98)' : '#fff'

      root.querySelectorAll('.ant-form-item-label label, .ant-form-item-label .ant-form-item-required').forEach((node) => {
        node.style.color = labelColor
      })

      root.querySelectorAll('.ant-calendar-picker-input, .ant-input').forEach((node) => {
        node.style.backgroundColor = inputBg
        node.style.borderColor = inputBorder
        node.style.color = labelColor
      })

      root.querySelectorAll('.ant-calendar-picker-icon, .ant-calendar-picker-clear').forEach((node) => {
        node.style.color = iconColor
      })

      root.querySelectorAll('.ant-calendar, .ant-calendar-picker-container, .ant-calendar-date-panel, .ant-calendar-month-panel, .ant-calendar-year-panel, .ant-calendar-decade-panel').forEach((node) => {
        node.style.backgroundColor = panelBg
        node.style.color = labelColor
      })
    },
    initComponents(option){
      if(this.IsToolBox)
      {
        return
      }
      let _t = this
      let i=0
      let refObj = this.detail.identifier
      let view = this.$refs[refObj]
      if (!view) {
        return
      }
      if (view.clientWidth === 0 || view.clientHeight === 0) {
        this.waitChartContainerReady(view, () => this.initComponents(option))
        return
      }
      if(!this.editMode) {
        if (this.echartsView != null && (typeof this.echartsView.dispose == "function")) {
          this.echartsView.dispose()
          this.echartsView = null
        }
      }
      this.fillOpacity=option.style.opacity
      this.SelectDateRange = this.SelectDateRange || moment().format("YYYY-MM-DD")
      this.applyFixedTimeAxis()
      this.option.tooltip.position = this.getTooltipPosition
      this.option.title.textStyle.color = option.style.foreColor
      this.option.title.textStyle.fontSize = option.style.fontSize
      this.option.title.textStyle.fontFamily  = option.style.fontFamily

      this.option.legend.textStyle.color = option.style.foreColor
      this.option.legend.textStyle.fontSize = option.style.fontSize
      this.option.legend.textStyle.fontFamily  = option.style.fontFamily

      this.option.xAxis.axisLabel.interval=30
      this.option.xAxis.axisLabel.textStyle.color = option.style.foreColor
      this.option.xAxis.axisLabel.textStyle.fontSize = option.style.fontSize
      this.option.xAxis.axisLabel.textStyle.fontFamily  = option.style.fontFamily

      this.option.yAxis.axisLabel.textStyle.color = option.style.foreColor
      this.option.yAxis.axisLabel.textStyle.fontSize = option.style.fontSize
      this.option.yAxis.axisLabel.textStyle.fontFamily  = option.style.fontFamily

      for( i=0;i<option.style.diy.length;i++)
      {
        if(option.style.diy[i].key=="ChartTitle")
        {
          this.option.title.text=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="IsShowZoom")
        {
          if(option.style.diy[i].value!=1)
          {
            this.option.dataZoom=[]
            this.option.grid.bottom='1%'
          }
          else
          {
            this.option.dataZoom=[
              {
                type: 'inside',
                xAxisIndex: 0,
                start: 0,
                end: 100,
                zoomLock: false
              },
              {
                type: 'slider',
                xAxisIndex: 0,
              }
            ]
            this.option.grid.bottom='10%'
          }
        }
        else if(option.style.diy[i].key=="TimelyRefresh")
        {
          this.ChartTimelyRefresh = parseInt(option.style.diy[i].value)*60*1000
        }
        else if(option.style.diy[i].key=="YMax")
        {
          if(option.style.diy[i].value==0)
          {
            this.option.yAxis.max = 'dataMax'
          }
          else {
            this.option.yAxis.max = option.style.diy[i].value
          }
        }
        else if(option.style.diy[i].key=="TimelyInitEcharts")
        {
            this.TimelyInitEcharts = parseInt(option.style.diy[i].value)
        }
        else if(option.style.diy[i].key=="YMin")
        {
          if(option.style.diy[i].value==0)
          {
            this.option.yAxis.min = 'dataMin'
          }
          else {
            this.option.yAxis.min = option.style.diy[i].value
          }
        }
        else if(option.style.diy[i].key=="ChartRange")
        {
          this.ChartRange = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="TimeInterval")
        {
          this.TimeInterval = parseInt(option.style.diy[i].value)
        }
        else if(option.style.diy[i].key=="IsShowDate")
        {
          this.IsShowDate = parseInt(option.style.diy[i].value)
        }
        else if(option.style.diy[i].key=="EchartsWidth")
        {
          this.EchartsWidth = parseInt(option.style.diy[i].value)
        }
        else if(option.style.diy[i].key=="EchartsXRotate")
        {
          this.option.xAxis.axisLabel.rotate = parseInt(option.style.diy[i].value)
        }
        else if(option.style.diy[i].key=="EchartsXFormat")
        {
          this.EchartsXFormat = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="EchartsXTheme")
        {
          this.EchartsTheme = option.style.diy[i].value
          if(this.echartsView) {
            this.echartsView.dispose()
            this.echartsView = null
          }
          this.echartsView = echarts.init(view, this.EchartsTheme);
          // 主题变化时强制更新视图，确保日期选择器样式同步
          this.$nextTick(() => {
            this.onResize()
            this.syncThemeStyles()
          })
        }
      }
      if (!this.echartsView) {
        this.echartsView = echarts.init(view, this.EchartsTheme);
      }
      i=0
      this.animateType = option.animate.selected
      if(option.animate.isExpression)
      {
        this.isStart = false
      }
      else
      {
        this.isStart = true
      }
      for( i=0;i<option.animate.animateElement.length;i++)
      {
        if(option.animate.animateElement[i].id=="millcolorGrad")
        {
          for(let k =0;k<option.animate.animateElement[i].elementList.length;k++)
          {
            if(option.animate.animateElement[i].elementList[k].key=="startColor")
            {
              this.startColor=option.animate.animateElement[i].elementList[k].value
            }
            else if(option.animate.animateElement[i].elementList[k].key=="stopColor")
            {
              this.stopColor=option.animate.animateElement[i].elementList[k].value
            }
            else if(option.animate.animateElement[i].elementList[k].key=="animateSpeed")
            {
              this.animateSpeed=option.animate.animateElement[i].elementList[k].value
            }
          }
        }
        else if(option.animate.animateElement[i].id=="blink")
        {
          for(let k =0;k<option.animate.animateElement[i].elementList.length;k++) {
            if (option.animate.animateElement[i].elementList[k].key == "blinkSpeed") {
              this.blinkSpeed = option.animate.animateElement[i].elementList[k].value
            }
          }
        }
        else if(option.animate.animateElement[i].id=="animateSpin")
        {
          for(let k =0;k<option.animate.animateElement[i].elementList.length;k++) {
            if (option.animate.animateElement[i].elementList[k].key == "spinSpeed") {
              this.animateSpinSpeed = option.animate.animateElement[i].elementList[k].value
            }
            else if (option.animate.animateElement[i].elementList[k].key == "spinDirection") {
              this.spinDirection = option.animate.animateElement[i].elementList[k].value
            }
          }
        }
      }

      this.echartsView.resize()
      if(!this.editMode)
      {
        for(i=0;i<this.option.series.length;i++)
        {
          this.option.series[i].data = []
          this.option.series[i].smooth= true
          this.seriesMap[this.option.series[i].dataID] = []
        }
        this.option.xAxis.data = [];
        this.option.legend.data=[]
        this.option.series=[]
        this.getList =[]
        for(let i =0;i<this.detail.active.length;i++)
        {
          if(this.detail.active[i].condition.dataName=="")
          {
            continue
          }
          let listObj = {}
          this.option.legend.data.push(typeof this.detail.active[i].condition.DeviceName!="undefined" ?this.detail.active[i].condition.DeviceName+"-"+this.detail.active[i].condition.dataName:this.detail.active[i].condition.dataName)
          let series= {
            name: typeof this.detail.active[i].condition.DeviceName!="undefined" ?this.detail.active[i].condition.DeviceName+"-"+this.detail.active[i].condition.dataName:this.detail.active[i].condition.dataName,
            type: 'line',
            dataID:this.detail.active[i].condition.deviceSN+this.detail.active[i].condition.dataID,
            data: [],
            itemStyle: {
              normal: {
                lineStyle: {
                  width:this.EchartsWidth
                }
              }
            },
            symbolSize: this.EchartsWidth,
          }
          listObj.DeviceUuid = this.detail.active[i].condition.deviceSN
          listObj.ModelDataUuid = this.detail.active[i].condition.dataID
          this.getList.push(listObj)
          this.option.series.push(series)
          this.seriesMap[this.detail.active[i].condition.deviceSN+this.detail.active[i].condition.dataID]=[]
        }
        let _t = this
        clearInterval(this.ChartTimelyRefreshTimer)
        this.ChartTimelyRefreshTimer = setInterval(function (){
          if(!_t.editMode){
           _t.timelyTimeGetHistory()
          }
        },this.ChartTimelyRefresh)
        this.timelyTimeGetHistory()
      }
      else
      {
        this.echartsView.setOption(this.option);
      }
      this.$nextTick(() => {
        this.syncThemeStyles()
      })
      if(!this.editMode){
        setTimeout(function (){
          _t.ReInitEcharts()
        }, 60000*this.TimelyInitEcharts)
      }
    },
    ReInitEcharts(){
      clearInterval(this.ChartTimelyRefreshTimer)
      if (this.echartsView != null&&(typeof this.echartsView.dispose=="function")) {
        this.echartsView.clear();
        this.echartsView.dispose()
        this.echartsView = null
      }
      let defaultOption= {
            title: {
              text: "",
              textStyle:{
                color:"",
                fontFamily:"",
                fontSize:"",
              }
            },
            tooltip: {
              trigger: 'axis',
              confine: true
            },
            legend: {
              data: ['Email', 'Union Ads', 'Video Ads', 'Direct', 'Search Engine'],
              textStyle:{
                color:"",
                fontFamily:"",
                fontSize:"",
              }
            },
            grid: {
              left: '1%',
              right: '1%',
              bottom: '12%',
              containLabel: true
            },
            toolbox: {
              feature: {
                saveAsImage: {}
              }
            },
            xAxis: {
              type: 'time',
              min: moment().startOf('day').valueOf(),
              max: moment().endOf('day').valueOf(),
              interval: 4 * 60 * 60 * 1000,
              boundaryGap: false,
              axisLine: {
                show: false, //是否显示坐标刻度
                lineStyle: {
                  color: '#eeeeee'
                }
              },
              axisLabel: {
                show: true,
                rotate:40,
                formatter: function (value) {
                  return moment(value).format('HH:mm')
                },
                textStyle: {
                  color: '#fff'
                }
              },
              data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
            },
            yAxis: {
              type: 'value',
              min:100,
              max:1000,
              axisLabel: {
                show: true,
                textStyle: {
                  color: '#fff'
                }
              },
            },
            series: [
              {
                name: 'Email',
                type: 'line',
                stack: 'Total',
                data: [120, 132, 101, 134, 90, 230, 210]
              },
              {
                name: 'Union Ads',
                type: 'line',
                stack: 'Total',
                data: [220, 182, 191, 234, 290, 330, 310]
              },
              {
                name: 'Video Ads',
                type: 'line',
                stack: 'Total',
                data: [150, 232, 201, 154, 190, 330, 410]
              },
              {
                name: 'Direct',
                type: 'line',
                stack: 'Total',
                data: [320, 332, 301, 334, 390, 330, 320]
              },
              {
                name: 'Search Engine',
                type: 'line',
                stack: 'Total',
                data: [820, 932, 901, 934, 1290, 1330, 1320]
              }
            ]
      }
      this.option = defaultOption
      this.option.tooltip.position = this.getTooltipPosition
      this.applyFixedTimeAxis()
      this.initComponents(this.detail)
    },
    timelyTimeGetHistory(){
      let _t = this
      if (!this.echartsView) {
        return
      }
      const params ={
        HistoryTime:this.SelectDateRange,
        List: this.getList
      }
      if(this.getList.length==0)
      {
        _t.echartsView.setOption(_t.option,true)
        return
      }
      GetChartDataHistoryTrendByDate(params).then(function (res){
        if(res.data.code==0&&res.data.data!=null)
        {
          _t.option.title.text = "历史趋势"
          for(let i=0;i<_t.option.series.length;i++)
          {
            _t.option.series[i].data =[]
            _t.seriesMap[_t.option.series[i].dataID] = []
          }
          _t.option.xAxis.data = []
          for(let i=0;i<res.data.data.length;i++)
          {
            let singleArray={
              value:[]
            }
            const item = res.data.data[i]
            let DateTime = moment(item.RecordTime).format(_t.EchartsXFormat);
            singleArray.value[0]=DateTime
            singleArray.value[1]=item.DataValue ? item.DataValue : 0
            const seriesKey = item.DeviceUuid + item.ModelDataUuid
            if (_t.seriesMap[seriesKey]) {
              _t.seriesMap[seriesKey].push(singleArray)
            }
          }
          _t.option.xAxis.data = Array.from(new Set(_t.option.xAxis.data))
          for(let i=0;i<_t.option.series.length;i++)
          {
            _t.option.series[i].data = _t.seriesMap[_t.option.series[i].dataID];
          }
        }
        else
        {
          _t.option.title.text = "暂无数据"
          _t.option.xAxis.data = []
          for(let i=0;i<_t.option.series.length;i++)
          {
            _t.option.series[i].data =[]
            _t.seriesMap[_t.option.series[i].dataID] = []
          }
        }
        _t.echartsView.setOption(_t.option,true)
      })
    },
    onResize() {
      if (this.echartsView) {
        this.echartsView.resize();
      }
    },
    updateView() {
      this.setOption(this.option);
    },
    // 获取当前时间
  },
  beforeDestroy () {
    clearInterval(this.ChartTimelyRefreshTimer)
    if (this.echartsView != null&&(typeof this.echartsView.dispose=="function")) {
      this.echartsView.dispose()
    }
  },
  created(){
    let _t = this
    const node = this.getNode()
    node.on('change:data', ({ current }) => {
      if(current) {
        _t.detail = current.detail
      }
    })
    node.on('change:size', ({ current }) => {
      _t.detail.style.position.w = current.width
      _t.detail.style.position.h = current.height
      _t.$nextTick(() => {
        _t.onResize()
      })
    });
    this.detail = node.getData().detail
    this.editMode = node.getData().editMode
    this.showDeviceUuid = node.getData().showDeviceUuid
    this.IsToolBox = node.getData().IsToolBox
    _t.$EventBus.$on('cell-editMode', (data) => {
      _t.editMode = data.edit
      _t.IsToolBox = data.toolbox
    })
  },
  mounted() {
    let _t = this
    this.$nextTick(function(){
      this.initComponents(this.detail);
      this.syncThemeStyles()
      let activeEvent = this.detail.identifier+"activeEvent"//动作数据
      let animateEvent = this.detail.identifier+"animateEvent"//动作数据

      _t.$EventBus.$on(activeEvent, (data) => {

      })
      _t.$EventBus.$on(animateEvent, (data) => {
        _t.isStart = data
      })
    });
  }
}
</script>

<style lang="less">
.history-trend-foreign-object {
  overflow: hidden;
}

.history-trend-layout {
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
  overflow: hidden;
}

.view-chart-real-data {
  margin-top: 0;
  height: auto;
  width: 100%;
  flex: 1 1 auto;
  min-height: 0;
  text-align: center;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
}

.history-trend-theme-root {
  width: 100%;
  flex: 0 0 auto;
}

.history-trend-theme-root :deep(.ant-form) {
  margin-bottom: 0 !important;
}

.history-trend-theme-root :deep(.ant-form-item) {
  margin-bottom: 0 !important;
}

.history-trend-theme-root :deep(.ant-form-item-label) {
  line-height: 24px !important;
}

.history-trend-theme-root :deep(.ant-form-item-control) {
  line-height: 24px !important;
}

.history-trend-theme-root :deep(.ant-calendar-picker),
.history-trend-theme-root :deep(.ant-calendar-picker-input),
.history-trend-theme-root :deep(.ant-input) {
  height: 24px !important;
  line-height: 24px !important;
  border-radius: 0 !important;
}

.history-trend-theme-root :deep(.ant-calendar-picker-input) {
  padding-top: 0 !important;
  padding-bottom: 0 !important;
}

.date-picker-row {
  width: 100%;
}


// 深色主题表单样式
.form-dark-theme {
  :deep(.ant-form-item-label > label) {
    color: rgba(255, 255, 255, 0.85) !important;
  }
  :deep(.ant-form-item-label) {
    .ant-form-item-required {
      color: rgba(255, 255, 255, 0.85) !important;
    }
    label {
      color: rgba(255, 255, 255, 0.85) !important;
    }
  }
  :deep(.ant-form-item-explain) {
    color: #ff4d4f;
  }
}

// 浅色主题表单样式
.form-light-theme {
  :deep(.ant-form-item-label > label) {
    color: rgba(0, 0, 0, 0.85) !important;
  }
  :deep(.ant-form-item-label) {
    .ant-form-item-required {
      color: rgba(0, 0, 0, 0.85) !important;
    }
    label {
      color: rgba(0, 0, 0, 0.85) !important;
    }
  }
}

// 深色主题日期选择器
.date-picker-dark {
  :deep(.ant-calendar-picker-input) {
    background-color: rgba(40, 40, 60, 0.8) !important;
    border-color: rgba(255, 255, 255, 0.2) !important;
    color: rgba(255, 255, 255, 0.85) !important;
  }
  :deep(.ant-calendar-picker-icon),
  :deep(.ant-calendar-picker-clear) {
    color: rgba(255, 255, 255, 0.65) !important;
  }
  :deep(.ant-calendar-picker-input:hover) {
    border-color: rgba(255, 255, 255, 0.4) !important;
  }
  :deep(.ant-calendar) {
    background-color: rgba(20, 20, 30, 0.95) !important;
    border-color: rgba(255, 255, 255, 0.15) !important;
  }
  :deep(.ant-calendar-header) {
    border-bottom-color: rgba(255, 255, 255, 0.15) !important;
  }
  :deep(.ant-calendar-title) {
    color: rgba(255, 255, 255, 0.85) !important;
  }
  :deep(.ant-calendar-month-panel-title) {
    color: rgba(255, 255, 255, 0.85) !important;
  }
  :deep(.ant-calendar-year-panel-title) {
    color: rgba(255, 255, 255, 0.85) !important;
  }
  :deep(.ant-calendar-day) {
    color: rgba(255, 255, 255, 0.85) !important;
  }
  :deep(.ant-calendar-day:hover) {
    background-color: rgba(255, 255, 255, 0.1) !important;
  }
  :deep(.ant-calendar-selected-day) {
    background-color: #1890ff !important;
    border-color: #1890ff !important;
  }
  :deep(.ant-calendar-today) {
    border-color: rgba(255, 255, 255, 0.3) !important;
  }
  :deep(.ant-calendar-month) {
    color: rgba(255, 255, 255, 0.85) !important;
  }
  :deep(.ant-calendar-month:hover) {
    background-color: rgba(255, 255, 255, 0.1) !important;
  }
  :deep(.ant-calendar-year) {
    color: rgba(255, 255, 255, 0.85) !important;
  }
  :deep(.ant-calendar-year:hover) {
    background-color: rgba(255, 255, 255, 0.1) !important;
  }
}

// 浅色主题日期选择器
.date-picker-light {
  :deep(.ant-calendar-picker-input) {
    background-color: #fff !important;
    border-color: #d9d9d9 !important;
    color: rgba(0, 0, 0, 0.85) !important;
  }
  :deep(.ant-calendar-picker-icon),
  :deep(.ant-calendar-picker-clear) {
    color: rgba(0, 0, 0, 0.45) !important;
  }
}

.date-picker-dropdown-dark {
  z-index: 2100 !important;

  .ant-calendar {
    background-color: rgba(20, 20, 30, 0.98) !important;
    border-color: rgba(255, 255, 255, 0.15) !important;
    color: rgba(255, 255, 255, 0.85) !important;
  }
  .ant-calendar-input,
  .ant-calendar-time-picker-input {
    background-color: rgba(40, 40, 60, 0.8) !important;
    border-color: rgba(255, 255, 255, 0.2) !important;
    color: rgba(255, 255, 255, 0.85) !important;
  }
  .ant-calendar-header,
  .ant-calendar-month-panel-header,
  .ant-calendar-year-panel-header,
  .ant-calendar-decade-panel-header,
  .ant-calendar-footer {
    border-color: rgba(255, 255, 255, 0.15) !important;
  }
  .ant-calendar-next-century-btn,
  .ant-calendar-next-year-btn,
  .ant-calendar-prev-century-btn,
  .ant-calendar-prev-year-btn,
  .ant-calendar-month-select,
  .ant-calendar-year-select,
  .ant-calendar-decade-panel-century-select,
  .ant-calendar-decade-panel-year-select,
  .ant-calendar-month-panel-year-select,
  .ant-calendar-year-panel-decade-select,
  .ant-calendar-today-btn,
  .ant-calendar-time-picker-btn,
  .ant-calendar-ok-btn,
  .ant-calendar-title,
  .ant-calendar-month-panel-title,
  .ant-calendar-year-panel-title,
  .ant-calendar-decade-panel-title,
  .ant-calendar-date,
  .ant-calendar-month-panel-month,
  .ant-calendar-year-panel-year,
  .ant-calendar-decade-panel-decade {
    color: rgba(255, 255, 255, 0.85) !important;
  }
  .ant-calendar-date:hover,
  .ant-calendar-month-panel-month:hover,
  .ant-calendar-year-panel-year:hover,
  .ant-calendar-decade-panel-decade:hover {
    background-color: rgba(255, 255, 255, 0.1) !important;
  }
  .ant-calendar-selected-day .ant-calendar-date,
  .ant-calendar-month-panel-selected-cell .ant-calendar-month-panel-month,
  .ant-calendar-year-panel-selected-cell .ant-calendar-year-panel-year,
  .ant-calendar-decade-panel-selected-cell .ant-calendar-decade-panel-decade {
    background: #1890ff !important;
    color: #fff !important;
  }
  .ant-calendar-today .ant-calendar-date {
    border-color: rgba(255, 255, 255, 0.45) !important;
  }
}

.date-picker-dropdown-light {
  z-index: 2100 !important;

  .ant-calendar {
    background-color: #fff !important;
    color: rgba(0, 0, 0, 0.85) !important;
  }
}

.form-dark-theme .ant-form-item-label > label,
.form-dark-theme .ant-form-item-label label,
.form-dark-theme .ant-form-item-label .ant-form-item-required {
  color: rgba(255, 255, 255, 0.85) !important;
}

.form-dark-theme .ant-form-item-explain {
  color: #ff4d4f !important;
}

.form-light-theme .ant-form-item-label > label,
.form-light-theme .ant-form-item-label label,
.form-light-theme .ant-form-item-label .ant-form-item-required {
  color: rgba(0, 0, 0, 0.85) !important;
}

.date-picker-dark .ant-input,
.date-picker-dark .ant-calendar-picker-input {
  background-color: rgba(40, 40, 60, 0.8) !important;
  border-color: rgba(255, 255, 255, 0.2) !important;
  color: rgba(255, 255, 255, 0.85) !important;
}

.date-picker-dark .ant-input:hover,
.date-picker-dark .ant-calendar-picker-input:hover,
.date-picker-dark:hover .ant-input {
  border-color: rgba(255, 255, 255, 0.4) !important;
}

.date-picker-dark .ant-calendar-picker-icon,
.date-picker-dark .ant-calendar-picker-clear {
  color: rgba(255, 255, 255, 0.65) !important;
}

.date-picker-light .ant-input,
.date-picker-light .ant-calendar-picker-input {
  background-color: #fff !important;
  border-color: #d9d9d9 !important;
  color: rgba(0, 0, 0, 0.85) !important;
}

.date-picker-light .ant-calendar-picker-icon,
.date-picker-light .ant-calendar-picker-clear {
  color: rgba(0, 0, 0, 0.45) !important;
}
</style>
