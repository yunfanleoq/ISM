<template>
  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"  x="0px" y="0px"  xml:space="preserve" :style="{'overflow': 'visible','width':(detail && detail.style && detail.style.position && detail.style.position.w) || '200px','height':(detail && detail.style && detail.style.position && detail.style.position.h) || '200px',}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <foreignObject style="overflow:visible;" pointer-events="all" :width="((detail && detail.style && detail.style.position && detail.style.position.w) || 200) - ((detail && detail.style && detail.style.borderWidth) || 0)*2" :height="((detail && detail.style && detail.style.position && detail.style.position.h) || 200) - ((detail && detail.style && detail.style.borderWidth) || 0)*2">
          <div
            class="view-chart-real-data"
            :class="{ 'overview-embedded-chart': isOverviewEmbeddedChart }"
            :ref="detail && detail.identifier ? detail.identifier : 'chart_default'"
            :style="{'overflow': 'visible','width':'100%','height':'100%'}"
          ></div>
      </foreignObject>
      <path
        v-if="isOverviewEmbeddedChart"
        class="overview-chart-frame"
        :class="`overview-chart-frame--${overviewChartKind}`"
        :d="overviewFramePath"
        fill="none"
        vector-effect="non-scaling-stroke"
        pointer-events="none"
      />
      <path
        v-if="isOverviewEmbeddedChart"
        class="overview-chart-frame-accent"
        :class="`overview-chart-frame-accent--${overviewChartKind}`"
        :d="overviewFrameAccentPath"
        fill="none"
        vector-effect="non-scaling-stroke"
        pointer-events="none"
      />
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
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
  name: 'ism-view-real-data-smooth-chart',
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
    }
  },
  computed: {
    isOverviewEmbeddedChart() {
      const name = String((this.detail && this.detail.name) || '')
      const pos = (this.detail && this.detail.style && this.detail.style.position) || {}
      return Number(pos.x) >= 1200 && /功率趋势|用电量趋势/.test(name)
    },
    overviewChartKind() {
      return /用电量|电度/.test(String((this.detail && this.detail.name) || '')) ? 'energy' : 'power'
    },
    overviewChartSize() {
      const pos = (this.detail && this.detail.style && this.detail.style.position) || {}
      return {
        width: Math.max(24, Number(pos.w) || 200),
        height: Math.max(24, Number(pos.h) || 200),
      }
    },
    overviewFramePath() {
      const w = this.overviewChartSize.width
      const h = this.overviewChartSize.height
      return `M 0.5 12 L 12 0.5 H ${w - 12} L ${w - 0.5} 12 V ${h - 12} L ${w - 12} ${h - 0.5} H 12 L 0.5 ${h - 12} Z`
    },
    overviewFrameAccentPath() {
      const w = this.overviewChartSize.width
      return `M 12 0.5 H ${Math.min(w - 12, 128)}`
    }
  },
  data() {
    return {
      detail:{},
      IsToolBox:false,
      editMode:true,
      strokeColor:"#000000",
      fill:"#A1BFE2",
      strokeWidth:0.3,
      ChartTimelyRefreshTimer:null,
      dimensionSampleTimer:null,
      fillOpacity:1,
      strokeOpacity:1,
      animateType:"blink",
      startColor:"#74f808",
      stopColor:"#74f808",
      animateSpeed:0.5,
      animateSpinSpeed:0.5,
      spinDirection:0,
      blinkSpeed:0.5,
      isStart:false,
      base:{
        "text": "configComponent.RealDataChart.PolygonalSmoothingTitle",
        "icon": "icon-zhexiantu-pinghuashuangxian",
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
                deviceName:"",
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
                deviceName:"",
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
                deviceName:"",
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
                deviceName:"",
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
                deviceName:"",
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
              "h": 300
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
                "value":"实时数据",
                "key":"ChartTitle",
              },
              {
                "name":"configComponent.ChartPublic.TimelyInitEcharts",
                "type":1,
                "value":60,
                "key":"TimelyInitEcharts",
              },
              {
                "name":"configComponent.ChartPublic.ChartTimelyRefresh",
                "type":1,
                "value":60,
                "key":"ChartTimelyRefresh",
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
                value:"HH:mm:ss",
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
            ]
          }
        }
      },
      EchartsXFormat:1,
      EchartsTheme:"dark",
      date: [],
      isFinish:0,
      yieldRate: [],
      yieldIndex: [],
      EchartsWidth:1,
      EchartsViewObj:null,
      timelySave:null,
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
          trigger: 'axis'
        },
        legend: {
          data: [],
          textStyle:{
            color:"",
            fontFamily:"",
            fontSize:"",
          }
        },
        grid: {
          left: '1%',
          right: '1%',
          bottom: '1%',
          containLabel: true
        },
        toolbox: { show: false },
        xAxis: {
          type: 'category',
          boundaryGap: false,
          axisLine: {
            show: true,
            lineStyle: {
              color: '#eeeeee'
            }
          },
          axisLabel: {
            show: true,
            rotate:40,
            color: '#fff',
          },
          data: []
        },
        yAxis: {
          type: 'value',
          min: 0,
          max: 100,
          splitLine: {
            show: true,
            lineStyle: { color: 'rgba(125, 160, 200, 0.18)', type: 'dashed' }
          },
          axisLabel: {
            show: true,
            color: '#fff',
          },
        },
        series: []
      },
      ShowChartVariable1IsCome:false,
      ShowChartVariable2IsCome:false,
      ShowChartVariable3IsCome:false,
      ShowChartVariable4IsCome:false,
      ShowChartVariable5IsCome:false,
      seriesMap:[]
    }
  },
  methods: {
    buildDimensionSample(chartTitle) {
      const timeAxis = ['00:00', '02:00', '04:00', '06:00', '08:00', '10:00',
        '12:00', '14:00', '16:00', '18:00', '20:00', '22:00']
      const sampleLine = (name, data, color) => ({
        name,
        type: 'line',
        smooth: true,
        symbol: 'none',
        isDimensionSample: true,
        lineStyle: { width: 2, type: 'dashed', color },
        areaStyle: { color: `${color}18` },
        data,
      })
      if (/用电量趋势|正有功电度|有功电度|用电量/.test(chartTitle)) {
        return {
          axis: timeAxis,
          unit: 'kWh',
          series: [
            sampleLine('累计用电量',
              [8200, 8260, 8310, 8350, 8420, 8510, 8620, 8730, 8840, 8930, 9020, 9100],
              '#4dabf7'),
          ],
        }
      }
      if (/功率趋势|总有功功率|总无功功率|总视在功率/.test(chartTitle)) {
        return {
          axis: timeAxis,
          unit: 'kW / kvar / kVA',
          series: [
            sampleLine('总有功功率',
              [380, 350, 325, 315, 345, 430, 555, 640, 615, 570, 505, 445], '#22d3ee'),
            sampleLine('总无功功率',
              [88, 82, 76, 74, 80, 98, 124, 142, 138, 128, 116, 102], '#a78bfa'),
            sampleLine('总视在功率',
              [410, 378, 351, 340, 372, 463, 598, 685, 660, 612, 546, 480], '#34d399'),
          ],
        }
      }
      return null
    },
    stopDimensionSampleAnimation() {
      if (this.dimensionSampleTimer) {
        clearInterval(this.dimensionSampleTimer)
        this.dimensionSampleTimer = null
      }
    },
    startDimensionSampleAnimation() {
      this.stopDimensionSampleAnimation()
      if (typeof window !== 'undefined'
        && window.matchMedia
        && window.matchMedia('(prefers-reduced-motion: reduce)').matches) return
      const samples = (this.option.series || []).filter(series => series.isDimensionSample)
      if (!samples.length) return
      this.dimensionSampleTimer = setInterval(() => {
        if (!this.echartsView) return
        samples.forEach(series => {
          const data = Array.isArray(series.data) ? series.data : []
          if (data.length < 2) return
          const min = Math.min(...data)
          const max = Math.max(...data)
          const range = Math.max(1, max - min)
          const last = Number(data[data.length - 1]) || 0
          const isEnergy = /用电量|电度/.test(series.name || '')
          const delta = isEnergy
            ? Math.max(range * 0.025, 0.1)
            : (Math.random() - 0.48) * range * 0.16
          data.shift()
          data.push(Number(Math.max(0, last + delta).toFixed(2)))
        })
        this.echartsView.setOption({ series: this.option.series }, false, true)
      }, 1800)
    },
    metricDisplayName(condition) {
      const raw = String((condition && condition.dataName) || '').trim()
      const deviceName = String((condition && condition.DeviceName) || '').trim()
      if (!raw) return '系统统计'
      if (deviceName && raw.indexOf(deviceName) === 0) {
        return raw.slice(deviceName.length).replace(/^[_\-\s]+/, '') || raw
      }
      const knownMetrics = ['总有功功率', '总无功功率', '总视在功率', '正有功电度', '有功电度', '用电量']
      return knownMetrics.find(metric => raw.endsWith(metric)) || raw
    },
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
    refreshChart(){
      if(this.echartsView==null)
      {
        return
      }
      this.date = []
      this.initComponents(this.detail)
    },
    initComponents(option){
      if(this.IsToolBox)
      {
        return
      }
      // 确保 option 存在
      if (!option || !option.style) {
        console.warn('ViewRealDataSmoothChart initComponents: option or option.style is undefined')
        return
      }
      let i=0

      let refObj = this.detail && this.detail.identifier ? this.detail.identifier : 'chart_' + Date.now()
      let view = this.$refs[refObj]
      // 确保 DOM 元素存在
      if (!view) {
        console.warn('ViewRealDataSmoothChart initComponents: cannot find DOM element with ref:', refObj)
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
      this.option.title.textStyle.color = option.style.foreColor
      this.option.title.textStyle.fontSize = option.style.fontSize
      this.option.title.textStyle.fontFamily  = option.style.fontFamily

      this.option.legend.textStyle.color = option.style.foreColor
      this.option.legend.textStyle.fontSize = option.style.fontSize
      this.option.legend.textStyle.fontFamily  = option.style.fontFamily

      this.option.xAxis.axisLabel.color = option.style.foreColor
      this.option.xAxis.axisLabel.fontSize = option.style.fontSize
      this.option.xAxis.axisLabel.fontFamily = option.style.fontFamily

      this.option.yAxis.axisLabel.color = option.style.foreColor
      this.option.yAxis.axisLabel.fontSize = option.style.fontSize
      this.option.yAxis.axisLabel.fontFamily = option.style.fontFamily

      const diy = option.style.diy || []
      for( i=0;i<diy.length;i++)
      {
        if(diy[i].key=="ChartTitle")
        {
          this.option.title.text=diy[i].value
        }
        else if(diy[i].key=="ChartTimelyRefresh")
        {
          this.ChartTimelyRefresh = parseInt(diy[i].value)
        }
        else if(diy[i].key=="YMax")
        {
          if(diy[i].value==0)
          {
            this.option.yAxis.max = 'dataMax'
          }
          else {
            this.option.yAxis.max = diy[i].value
          }
        }
        else if(diy[i].key=="TimelyInitEcharts")
        {
            this.TimelyInitEcharts = parseInt(diy[i].value)
        }
        else if(diy[i].key=="YMin")
        {
          if(diy[i].value==0)
          {
            this.option.yAxis.min = 'dataMin'
          }
          else {
            this.option.yAxis.min = diy[i].value
          }
        }
        else if(diy[i].key=="EchartsWidth")
        {
          this.EchartsWidth = parseInt(diy[i].value)
        }
        else if(diy[i].key=="EchartsXRotate")
        {
          this.option.xAxis.axisLabel.rotate = parseInt(diy[i].value)
        }
        else if(diy[i].key=="EchartsXFormat")
        {
          this.EchartsXFormat = diy[i].value
        }
        else if(diy[i].key=="EchartsXTheme")
        {
          this.EchartsTheme = diy[i].value
          if(this.echartsView) {
            this.echartsView.dispose()
            this.echartsView = null
          }
          this.echartsView = echarts.init(view, this.EchartsTheme);
        }
      }
      if (!this.echartsView) {
        this.echartsView = echarts.init(view, this.EchartsTheme);
      }
      i=0
      this.animateType = (option.animate && option.animate.selected) || []
      if(option.animate && option.animate.isExpression)
      {
        this.isStart = false
      }
      else
      {
        this.isStart = true
      }
      const animateElement = (option.animate && option.animate.animateElement) || []
      for( i=0;i<animateElement.length;i++)
      {
        if(animateElement[i].id=="millcolorGrad")
        {
          for(let k =0;k<animateElement[i].elementList.length;k++)
          {
            if(animateElement[i].elementList[k].key=="startColor")
            {
              this.startColor=animateElement[i].elementList[k].value
            }
            else if(animateElement[i].elementList[k].key=="stopColor")
            {
              this.stopColor=animateElement[i].elementList[k].value
            }
            else if(animateElement[i].elementList[k].key=="animateSpeed")
            {
              this.animateSpeed=animateElement[i].elementList[k].value
            }
          }
        }
        else if(animateElement[i].id=="blink")
        {
          for(let k =0;k<animateElement[i].elementList.length;k++) {
            if (animateElement[i].elementList[k].key == "blinkSpeed") {
              this.blinkSpeed = animateElement[i].elementList[k].value
            }
          }
        }
        else if(animateElement[i].id=="animateSpin")
        {
          for(let k =0;k<animateElement[i].elementList.length;k++) {
            if (animateElement[i].elementList[k].key == "spinSpeed") {
              this.animateSpinSpeed = animateElement[i].elementList[k].value
            }
            else if (animateElement[i].elementList[k].key == "spinDirection") {
              this.spinDirection = animateElement[i].elementList[k].value
            }
          }
        }
      }

      this.echartsView.resize()
      if(!this.editMode)
      {
        this.option.xAxis.data = [];
        for(i=0;i<this.option.series.length;i++)
        {
          this.option.series[i].data = []
          this.option.series[i].smooth= true
        }
        this.option.legend.data=[]
        this.option.series=[]
        for(let i =0;i<this.detail.active.length;i++)
        {
          if(this.detail.active[i].condition.dataName=="")
          {
            continue
          }
          const metricName = this.metricDisplayName(this.detail.active[i].condition)
          this.option.legend.data.push(metricName)
          let series= {
            // 图表只表达指标，不暴露某一具体设备名称，避免误读为设备详情趋势。
            name: metricName,
            type: 'line',
            smooth: true,
            dataID:this.detail.active[i].condition.deviceSN+this.detail.active[i].condition.dataID,
            data: [],
            lineStyle: {
              width: this.EchartsWidth,
            },
            symbolSize: this.EchartsWidth,
          }
          this.option.series.push(series)
          this.seriesMap[this.detail.active[i].condition.deviceSN+this.detail.active[i].condition.dataID]=[]
        }
        const boundMetricNames = this.detail.active
          .map(active => this.metricDisplayName(active.condition))
          .join(' ')
        const chartTitle = `${String(this.option.title.text || '')} ${boundMetricNames}`
        if (/功率|电度|用电量/.test(chartTitle)) {
          this.option.legend.show = true
        }
        const hasSamples = this.option.series.some(series =>
          Array.isArray(series.data) && series.data.length > 0
        )
        if (!hasSamples) {
          const dimensionSample = this.buildDimensionSample(chartTitle)
          if (dimensionSample) {
            this.option.xAxis.data = dimensionSample.axis
            this.option.xAxis.axisLabel.rotate = 0
            this.option.yAxis.name = dimensionSample.unit
            this.option.yAxis.min = 'dataMin'
            this.option.yAxis.max = 'dataMax'
            this.option.series.push(...dimensionSample.series)
            this.option.legend.data = dimensionSample.series.map(series => series.name)
          } else {
            this.option.xAxis.data = ['-5m', '-4m', '-3m', '-2m', '-1m', '现在']
          }
          this.option.graphic = []
        } else {
          this.option.graphic = []
        }
      }
      let _t = this
      // this.option.series[0].itemStyle.normal.color = this.progressColor
      setTimeout(function (){
        if (!_t.echartsView) {
          return
        }
        _t.echartsView.setOption(_t.option,true)
        _t.echartsView.resize()
        _t.startDimensionSampleAnimation()
      }, 100)
      if(!this.editMode){
        setTimeout(function (){
          _t.ReInitEcharts()
        }, 60000*this.TimelyInitEcharts)
      }
    },
    onResize() {
      if (this.echartsView) {
        this.echartsView.resize();
      }
    },
    ReInitEcharts(){
      clearInterval(this.ChartTimelyRefreshTimer)
      if (this.echartsView != null&&(typeof this.echartsView.dispose=="function")) {
        this.echartsView.clear();
        this.echartsView.dispose()
        this.echartsView = null
      }
      let defaultOption={
        title: {
          text: "",
          textStyle:{
            color:"",
            fontFamily:"",
            fontSize:"",
          }
        },
        tooltip: {
          trigger: 'axis'
        },
        legend: {
          data: [],
          textStyle:{
            color:"",
            fontFamily:"",
            fontSize:"",
          }
        },
        grid: {
          left: '1%',
          right: '1%',
          bottom: '1%',
          containLabel: true
        },
        toolbox: { show: false },
        xAxis: {
          type: 'category',
          boundaryGap: false,
          axisLine: {
            show: true,
            lineStyle: {
              color: '#eeeeee'
            }
          },
          axisLabel: {
            show: true,
            rotate:40,
            color: '#fff',
          },
          data: []
        },
        yAxis: {
          type: 'value',
          min: 0,
          max: 100,
          splitLine: {
            show: true,
            lineStyle: { color: 'rgba(125, 160, 200, 0.18)', type: 'dashed' }
          },
          axisLabel: {
            show: true,
            color: '#fff',
          },
        },
        series: []
      }
      this.option = defaultOption
      this.initComponents(this.detail)
    },
    updateView() {
      this.setOption(this.option);
    },
    // 获取当前时间
    getTime () {
      let ts = arguments[0] || 0;
      let t, h, i, s,ms,m,d,y;
      t = new Date();
      y = t.getFullYear()
      m = t.getMonth()+1
      d = t.getDate()
      h = t.getHours();
      i = t.getMinutes();
      s = t.getSeconds();
      ms = t.getMilliseconds();
      // 定义时间格式
      return y+"-"+(m < 10 ? '0' + m : m)+"-"+(d < 10 ? '0' + d : d)+" "+(h < 10 ? '0' + h : h) + ':' + (i < 10 ? '0' + i : i) + ':' + (s < 10 ? '0' + s : s)+"."+ms;
    },
    getSecondByDateSub(begin)
    {
      let beginSplit = begin.split(" ")
      let beginstr = ""
      let endDate = new Date();
      if(beginSplit.length==1)
      {
        beginstr = endDate.getFullYear()+"-"+(endDate.getMonth()+1)+"-"+endDate.getDate()+" "+begin
      }
      else
      {
        beginstr = begin
      }

      let beginDate = new Date(beginstr);
      let diff = endDate.getTime() - beginDate.getTime();
      let sec = diff / 1000;
      return sec;
    },
    // 添加实时数据
    addData : function(data) {
      let c_data = this.getTime()
      c_data = moment(c_data).format(this.EchartsXFormat);

      const hadDimensionSample = this.option.series.some(series => series.isDimensionSample)
      const realSeries = this.option.series.filter(series => !series.isDimensionSample)
      for(let i=0;i<realSeries.length;i++)
      {
        realSeries[i].data = this.seriesMap[realSeries[i].dataID] || []
      }
      const hasRealSamples = realSeries.some(series => series.data.length > 0)
      if (!hasRealSamples) {
        this.echartsView.setOption(this.option, true)
        return
      }
      this.stopDimensionSampleAnimation()
      if (hadDimensionSample) {
        this.option.series = realSeries
        this.option.legend.data = realSeries.map(series => series.name)
        this.option.graphic = []
        this.option.xAxis.data = []
        this.option.yAxis.name = ''
        this.option.yAxis.min = 'dataMin'
        this.option.yAxis.max = 'dataMax'
      }
      this.option.xAxis.data.push(c_data);

      for(let i=0;i<this.option.xAxis.data.length;i++)
      {
        if(this.getSecondByDateSub(this.option.xAxis.data[i] )>(this.ChartTimelyRefresh*60))
        {
          this.option.xAxis.data.splice(i,1)
          for(let k=0;k<this.option.series.length;k++)
          {
            this.option.series[k].data.splice(i,1)
          }
        }
      }
      // this.SetMemDataBySelf(JSON.stringify(this.option))
      // 重新将数组赋值给echarts选项
      this.echartsView.setOption(this.option);
    }
  },
  beforeDestroy () {
    clearInterval(this.ChartTimelyRefresh)
    this.stopDimensionSampleAnimation()
    if (this.echartsView != null&&(typeof this.echartsView.dispose=="function")) {
      this.echartsView.clear();
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
      // _t.initComponents(_t.detail);
    })
  },
  mounted() {
      let _t = this
    this.$nextTick(function(){
      this.initComponents(this.detail);
      let activeEvent = this.detail.identifier+"activeEvent"//动作数据
      let animateEvent = this.detail.identifier+"animateEvent"//动作数据
      _t.$EventBus.$on("DealWithRealDataFinish",(data) => {
        if(_t.isFinish==0)
        {
          return
        }

        if(_t.ShowChartVariable1IsCome == false)
        {
          for(let i =0;i<_t.detail.active.length;i++)
          {
            if(_t.detail.active[i].condition.dataName=="")
            {
              continue
            }
            if(_t.detail.active[i].id=="ShowChartVariable1")
            {
              const seriesKey = _t.detail.active[i].condition.deviceSN+_t.detail.active[i].condition.dataID
              const seriesData = _t.seriesMap[seriesKey]
              let seriesMapLength = seriesData ? seriesData.length : 0
              if(seriesMapLength>=1)
              {
                let seriesMapDataArray = seriesData
                let seriesMapData = seriesMapDataArray[seriesMapLength-1]
                seriesData.push(seriesMapData)
              }
            }
          }
        }
        else if(_t.ShowChartVariable2IsCome == false)
        {
          for(let i =0;i<_t.detail.active.length;i++)
          {
            if(_t.detail.active[i].condition.dataName=="")
            {
              continue
            }
            if(_t.detail.active[i].id=="ShowChartVariable2")
            {
              const seriesKey = _t.detail.active[i].condition.deviceSN+_t.detail.active[i].condition.dataID
              const seriesData = _t.seriesMap[seriesKey]
              let seriesMapLength = seriesData ? seriesData.length : 0
              if(seriesMapLength>=1)
              {
                let seriesMapDataArray = seriesData
                let seriesMapData = seriesMapDataArray[seriesMapLength-1]
                seriesData.push(seriesMapData)
              }
            }
          }
        }
        else if(_t.ShowChartVariable3IsCome == false)
        {
          for(let i =0;i<_t.detail.active.length;i++)
          {
            if(_t.detail.active[i].condition.dataName=="")
            {
              continue
            }
            if(_t.detail.active[i].id=="ShowChartVariable3")
            {
              const seriesKey = _t.detail.active[i].condition.deviceSN+_t.detail.active[i].condition.dataID
              const seriesData = _t.seriesMap[seriesKey]
              let seriesMapLength = seriesData ? seriesData.length : 0
              if(seriesMapLength>=1)
              {
                let seriesMapDataArray = seriesData
                let seriesMapData = seriesMapDataArray[seriesMapLength-1]
                seriesData.push(seriesMapData)
              }
            }
          }
        }
        else if(_t.ShowChartVariable4IsCome == false)
        {
          for(let i =0;i<_t.detail.active.length;i++)
          {
            if(_t.detail.active[i].condition.dataName=="")
            {
              continue
            }
            if(_t.detail.active[i].id=="ShowChartVariable4")
            {
              const seriesKey = _t.detail.active[i].condition.deviceSN+_t.detail.active[i].condition.dataID
              const seriesData = _t.seriesMap[seriesKey]
              let seriesMapLength = seriesData ? seriesData.length : 0
              if(seriesMapLength>=1)
              {
                let seriesMapDataArray = seriesData
                let seriesMapData = seriesMapDataArray[seriesMapLength-1]
                seriesData.push(seriesMapData)
              }
            }
          }
        }
        else if(_t.ShowChartVariable5IsCome == false)
        {
          for(let i =0;i<_t.detail.active.length;i++)
          {
            if(_t.detail.active[i].condition.dataName=="")
            {
              continue
            }
            if(_t.detail.active[i].id=="ShowChartVariable5")
            {
              const seriesKey = _t.detail.active[i].condition.deviceSN+_t.detail.active[i].condition.dataID
              const seriesData = _t.seriesMap[seriesKey]
              let seriesMapLength = seriesData ? seriesData.length : 0
              if(seriesMapLength>=1)
              {
                let seriesMapDataArray = seriesData
                let seriesMapData = seriesMapDataArray[seriesMapLength-1]
                seriesData.push(seriesMapData)
              }
            }
          }
        }
        _t.addData()
        _t.isFinish=0
        _t.ShowChartVariable5IsCome = false
        _t.ShowChartVariable4IsCome = false
        _t.ShowChartVariable3IsCome = false
        _t.ShowChartVariable2IsCome = false
        _t.ShowChartVariable1IsCome = false
      })
      _t.$EventBus.$on(activeEvent, (data) => {
        let valueObj = parseFloat(data.result)
        if(!isNaN(valueObj)) {
          let c_data = _t.getTime()
          c_data = moment(c_data).format(_t.EchartsXFormat);
          if(_t.option.xAxis.data.indexOf(c_data)!=-1)
          {
            return
          }
          if (data.ID == "ShowChartVariable1") {
            const seriesData = _t.seriesMap[data.DeviceSN+data.dataID]
            if (seriesData) {
              seriesData.push(parseFloat(data.result))
            }
            _t.ShowChartVariable1IsCome = true
          } else if (data.ID == "ShowChartVariable2") {
            const seriesData = _t.seriesMap[data.DeviceSN+data.dataID]
            if (seriesData) {
              seriesData.push(parseFloat(data.result))
            }
            _t.ShowChartVariable2IsCome = true
          } else if (data.ID == "ShowChartVariable3") {
            const seriesData = _t.seriesMap[data.DeviceSN+data.dataID]
            if (seriesData) {
              seriesData.push(parseFloat(data.result))
            }
            _t.ShowChartVariable3IsCome = true
          } else if (data.ID == "ShowChartVariable4") {
            const seriesData = _t.seriesMap[data.DeviceSN+data.dataID]
            if (seriesData) {
              seriesData.push(parseFloat(data.result))
            }
            _t.ShowChartVariable4IsCome = true
          } else if (data.ID == "ShowChartVariable5") {
            const seriesData = _t.seriesMap[data.DeviceSN+data.dataID]
            if (seriesData) {
              seriesData.push(parseFloat(data.result))
            }
            _t.ShowChartVariable5IsCome = true
          }
          _t.isFinish=1
        }
        else{
          console.log("isanna",data.result)
        }
      })
      _t.$EventBus.$on(animateEvent, (data) => {
        _t.isStart = data
      })
    });
  }
}
</script>

<style lang="less">
.view-chart-real-data {
  position: relative;
  height: 100%;
  width: 100%;
  box-sizing: border-box;
  text-align: center;
  display: flex;
  justify-content: center;
  align-items: center;
}
.view-chart-real-data.overview-embedded-chart {
  background:
    radial-gradient(circle at 16% 22%, rgba(0, 222, 255, 0.1), transparent 27%),
    radial-gradient(circle at 86% 74%, rgba(94, 92, 230, 0.08), transparent 30%),
    linear-gradient(rgba(42, 126, 151, 0.055) 1px, transparent 1px),
    linear-gradient(90deg, rgba(42, 126, 151, 0.055) 1px, transparent 1px),
    linear-gradient(145deg, rgba(8, 31, 48, 0.72), rgba(5, 17, 30, 0.46));
  background-size: auto, auto, 24px 24px, 24px 24px, auto;
  box-shadow: inset 0 0 28px rgba(0, 157, 198, 0.055);
  animation: overviewChartGridFlow 12s linear infinite;
}
.view-chart-real-data.overview-embedded-chart::before {
  content: "";
  position: absolute;
  z-index: 2;
  top: 0;
  left: 0;
  right: 0;
  height: 1px;
  pointer-events: none;
  background: linear-gradient(90deg, rgba(0, 229, 255, 0.48), rgba(0, 229, 255, 0.14) 36%, transparent 76%);
  box-shadow: 0 0 5px rgba(0, 229, 255, 0.12);
}
.overview-chart-frame {
  stroke-width: 1px !important;
  stroke-linejoin: miter !important;
  opacity: 0.42;
}
.overview-chart-frame--power { stroke: #16d9ee !important; }
.overview-chart-frame--energy { stroke: #6589ff !important; }
.overview-chart-frame-accent {
  stroke-width: 1.5px !important;
  opacity: 0.82;
  filter: drop-shadow(0 0 3px currentColor);
}
.overview-chart-frame-accent--power { stroke: #5cf1d2 !important; color: #5cf1d2; }
.overview-chart-frame-accent--energy { stroke: #9c7cff !important; color: #9c7cff; }
@keyframes overviewChartGridFlow {
  from { background-position: 0 0, 0 0, 0 0, 0 0, 0 0; }
  to { background-position: 0 0, 0 0, 24px 24px, 24px 24px, 0 0; }
}
@media (prefers-reduced-motion: reduce) {
  .view-chart-real-data.overview-embedded-chart { animation: none; }
}
</style>
