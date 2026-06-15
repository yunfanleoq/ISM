<template>
  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"  x="0px" y="0px"  xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <foreignObject style="overflow:visible;" pointer-events="all" :width="detail.style.position.w" :height="detail.style.position.h">
          <div
            class="view-chart-real-data"
            :ref="detail.identifier"
            :style="{'width': detail.style.position.w + 'px', 'height': detail.style.position.h + 'px'}"
          ></div>
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
import {GetChartDataHistoryList} from "@/services/report";
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
export default {
  mixins: [ISMChildAutoMixin],
  name: 'ism-view-history-data-chart',
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
  data() {
    return {
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
      base:{
        "text": "configComponent.RealDataChart.PolygonalHistoryTitle",
        "icon": "icon-duijizhexian",
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
                "value":"历史数据",
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
                "name":"configComponent.ChartPublic.ValueType",
                type:6,
                value:5,
                enumList:[
                  {
                    value:0,
                    option:"configComponent.ChartPublic.ValueList.raw"
                  },
                  {
                    value:1,
                    option:"configComponent.ChartPublic.ValueList.max"
                  },
                  {
                    value:2,
                    option:"configComponent.ChartPublic.ValueList.min"
                  },
                  {
                    value:3,
                    option:"configComponent.ChartPublic.ValueList.diff"
                  },
                  {
                    value:4,
                    option:"configComponent.ChartPublic.ValueList.sum"
                  },
                  {
                    value:5,
                    option:"configComponent.ChartPublic.ValueList.average"
                  },
                ],
                "key":"ValueType",
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
          trigger: 'axis'
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
            dataZoom: {
              yAxisIndex: 'none'
            },
            restore: {},
            saveAsImage: {}
          }
        },
        xAxis: {
          type: 'time',
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
            formatter: '{yy}-{MM}-{dd}\n{HH}:{mm}:{ss}' , // 仅显示时分
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
        view.style.width = this.detail.style.position.w + 'px'
        view.style.height = this.detail.style.position.h + 'px'
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
      let _t = this
      let i=0
      let refObj = this.detail.identifier
      let view = this.$refs[refObj]
      // 防止 DOM 未就绪（foreignObject 内嵌 SVG 场景）
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
        // 非编辑模式：每次都重新 init echarts
        this.echartsView = echarts.init(view, this.EchartsTheme)
      }
      this.fillOpacity=option.style.opacity
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
            this.option.grid.bottom='12%'
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
        else if(option.style.diy[i].key=="ValueType")
        {
          this.ValueType = parseInt(option.style.diy[i].value)
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

      if (this.echartsView) {
        this.echartsView.resize()
      }
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
              trigger: 'axis'
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
      this.initComponents(this.detail)
    },
    timelyTimeGetHistory(){
      let _t = this
      if (!this.echartsView) {
        return
      }
      let ChartRangeTime = 0
      if(this.ChartRange==0)
      {
        ChartRangeTime = 60
      }
      else if(this.ChartRange==1){
        ChartRangeTime = 60*5
      }
      else if(this.ChartRange==2){
        ChartRangeTime = 60*15
      }
      else if(this.ChartRange==3){
        ChartRangeTime = 60*24
      }
      else if(this.ChartRange==4){
        ChartRangeTime = 60*24*3
      }
      else if(this.ChartRange==5){
        ChartRangeTime = 60*24*7
      }
      else if(this.ChartRange==6){
        ChartRangeTime = 60*24*15
      }
      else if(this.ChartRange==7){
        ChartRangeTime = 60*24*30
      }
      const params ={
        TimeIn:this.TimeInterval,
        HistoryTime:ChartRangeTime,
        List: this.getList
      }
      if(this.getList.length==0)
      {
        _t.echartsView.setOption(_t.option,true)
        return
      }
      GetChartDataHistoryList(params).then(function (res){
        if(res.data.code==0&&res.data.data!=null)
        {
          for(let i=0;i<_t.option.series.length;i++)
          {
            _t.option.series[i].data =[]
            _t.seriesMap[_t.option.series[i].dataID] = []
          }
          _t.option.xAxis.data = []
          for(let i=0;i<res.data.data.length;i++)
          {
            let DateTime = moment(res.data.data[i].HistoryRecordDateTime).format(_t.EchartsXFormat);
            // _t.option.xAxis.data.push(DateTime)
            let temp_average=""
            let temp_diff=""
            let temp_max=""
            let temp_min=""
            let temp_sum=""
            for(let k=0;k<res.data.data[i].dataList.length;k++)
            {
              let item = res.data.data[i].dataList[k]
              if (!Object.keys(_t.seriesMap).includes(item.DeviceUuid + item.ModelDataUuid))
              {
                continue
              }
              let singleArray={
                value:[]
              }

              if(item.average=="-") {
                item.average=temp_average
              }else{
                temp_average=item.average
              }

              if(item.diff=="-") {
                item.diff=temp_diff
              }else{
                temp_diff=item.diff
              }

               if(item.max=="-") {
                item.max=temp_max
                }else{
                 temp_max=item.max
               }

               if(item.min=="-") {
                item.min=temp_min
                }else{
                 temp_min=item.min
               }

               if(item.sum=="-") {
                item.sum=temp_sum
              }else{
                 temp_sum=item.sum
               }
              singleArray.value[0]=DateTime
              if(_t.ValueType==0)
              {
                singleArray.value[1]=item.Value ? item.Value : 0
              }
              else if(_t.ValueType==1)
              {
                singleArray.value[1]=item.max ? item.max : 0
              }
              else if(_t.ValueType==2)
              {
                singleArray.value[1]=item.min ? item.min : 0
              }
              else if(_t.ValueType==3)
              {
                singleArray.value[1]=item.diff ? item.diff : 0
              }
              else if(_t.ValueType==4)
              {
                singleArray.value[1]=item.sum ? item.sum : 0
              }
              else if(_t.ValueType==5)
              {
                singleArray.value[1]=item.average ? item.average : 0
              }
              const seriesKey = item.DeviceUuid + item.ModelDataUuid
              if (_t.seriesMap[seriesKey]) {
                _t.seriesMap[seriesKey].push(singleArray)
              }
            }
          }
          for(let i=0;i<_t.option.series.length;i++)
          {
            _t.option.series[i].data = _t.seriesMap[_t.option.series[i].dataID] || [];
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
        const view = _t.$refs[_t.detail.identifier]
        if (view) {
          view.style.width = current.width + 'px'
          view.style.height = current.height + 'px'
        }
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
.view-chart-real-data {
  height: 100%;
  width: 100%;
  text-align: center;
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>
