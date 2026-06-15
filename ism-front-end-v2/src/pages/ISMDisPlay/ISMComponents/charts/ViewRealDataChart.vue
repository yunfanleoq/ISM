<template>
  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"  x="0px" y="0px"  xml:space="preserve" :style="{'overflow': 'visible','width':(detail && detail.style && detail.style.position && detail.style.position.w) || '200px','height':(detail && detail.style && detail.style.position && detail.style.position.h) || '200px',}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <foreignObject style="overflow:visible;" pointer-events="all" :width="((detail && detail.style && detail.style.position && detail.style.position.w) || 200) - ((detail && detail.style && detail.style.borderWidth) || 0)*2" :height="((detail && detail.style && detail.style.position && detail.style.position.h) || 200) - ((detail && detail.style && detail.style.borderWidth) || 0)*2">
          <div
            class="view-chart-real-data"
            :ref="detail && detail.identifier ? detail.identifier : 'chart_default'"
            :style="{'overflow': 'visible','width':'100%','height':'100%'}"
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
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
  name: 'ism-view-real-data-chart',
  i18n: require('@/i18n/language'),
  inject: ['getNode'],
  props: {

  },
  watch: {
    detail: {
      handler(newVal, oldVal) {
        if(this.editMode){
          this.initComponents(newVal)
          this.onResize()
        }
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
        "text": "configComponent.RealDataChart.PolygonalTitle",
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
          bottom: '1%',
          containLabel: true
        },
        toolbox: {
          feature: {
            saveAsImage: {}
          }
        },
        xAxis: {
          type: 'category',
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
            data: [120, 132, 101, 134, 120, 230, 210]
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
            data: [820, 932, 901, 934, 345, 456, 785]
          }
        ]
      },
      TimelyInitEcharts:60,
      ShowChartVariable1IsCome:false,
      ShowChartVariable2IsCome:false,
      ShowChartVariable3IsCome:false,
      ShowChartVariable4IsCome:false,
      ShowChartVariable5IsCome:false,
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
        console.warn('ViewRealDataChart initComponents: option or option.style is undefined')
        return
      }
      let i=0

      let refObj = this.detail && this.detail.identifier ? this.detail.identifier : 'chart_' + Date.now()
      let view = this.$refs[refObj]
      // 确保 DOM 元素存在
      if (!view) {
        console.warn('ViewRealDataChart initComponents: cannot find DOM element with ref:', refObj)
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
        else if(option.style.diy[i].key=="ChartTimelyRefresh")
        {
          this.ChartTimelyRefresh = parseInt(option.style.diy[i].value)
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
        else if(option.style.diy[i].key=="EchartsWidth")
        {
          this.EchartsWidth = parseInt(option.style.diy[i].value)
        }
        else if(option.style.diy[i].key=="EchartsXRotate")
        {
          this.option.xAxis.axisLabel.rotate = parseInt(option.style.diy[i].value)
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

      this.echartsView.resize()
      if(!this.editMode)
      {
        this.option.xAxis.data = [];
        for(i=0;i<this.option.series.length;i++)
        {
          this.option.series[i].data = []
          this.option.series[i].smooth= false
        }
        this.option.legend.data=[]
        this.option.series=[]
        for(let i =0;i<this.detail.active.length;i++)
        {
          if(this.detail.active[i].condition.dataName=="")
          {
            continue
          }
          this.option.legend.data.push(this.detail.active[i].condition.DeviceName+"-"+this.detail.active[i].condition.dataName)
          let series= {
            name: this.detail.active[i].condition.DeviceName+"-"+this.detail.active[i].condition.dataName,
            type: 'line',
            large: true,
            largeThreshold: 10000,
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
          this.option.series.push(series)
          this.seriesMap[this.detail.active[i].condition.deviceSN+this.detail.active[i].condition.dataID]=[]
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
      }, 100)
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
              bottom: '1%',
              containLabel: true
        },
        toolbox: {
          feature: {
            saveAsImage: {}
          }
        },
        xAxis: {
          type: 'category',
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
            data: [120, 132, 101, 134, 120, 230, 210]
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
            data: [820, 932, 901, 934, 345, 456, 785]
          }
        ]
      }
      this.option = defaultOption
      this.initComponents(this.detail)
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

      for(let i=0;i<this.option.series.length;i++)
      {
        this.option.series[i].data = this.seriesMap[this.option.series[i].dataID] || []
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
      // 重新将数组赋值给echarts选项
      this.echartsView.setOption(this.option,true);
    }
  },
  beforeDestroy () {
    clearInterval(this.ChartTimelyRefreshTimer)
    if (this.echartsView != null&&(typeof this.echartsView.dispose=="function")) {
      this.echartsView.clear();
      this.echartsView.dispose()
      this.echartsView = null
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
              break
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
              break
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
              break
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
              break
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
              break
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
  height: 100%;
  width: 100%;
  text-align: center;
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>
