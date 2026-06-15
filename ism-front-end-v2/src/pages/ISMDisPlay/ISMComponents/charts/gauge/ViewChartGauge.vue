<template>
  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"  x="0px" y="0px"  xml:space="preserve" :style="{'overflow': 'visible','width':(detail && detail.style && detail.style.position && detail.style.position.w) || '200px','height':(detail && detail.style && detail.style.position && detail.style.position.h) || '200px'}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <foreignObject style="overflow:visible;" pointer-events="all" :width="((detail && detail.style && detail.style.position && detail.style.position.w) || 200) - ((detail && detail.style && detail.style.borderWidth) || 0)*2" :height="((detail && detail.style && detail.style.position && detail.style.position.h) || 200) - ((detail && detail.style && detail.style.borderWidth) || 0)*2">
        <div class="view-chart-gauge" :ref="detail && detail.identifier ? detail.identifier : 'chart_default'" :style="{'overflow': 'visible','width':(detail && detail.style && detail.style.position && detail.style.position.w) || '200px','height':(detail && detail.style && detail.style.position && detail.style.position.h) || '200px',}">
            Click to bind data.
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
import BaseView from '../../View';
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
  name: 'ism-view-chart-gauge-0',
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
      ShowDataResult:0,
        base:{
          "text": "configComponent.chartGauge.diyGauge",
          "icon": "icon-gauge-dashboard-",
          "isFontIcon": true,
          "info": {
            "type": "chart-gauge",
            "action": [],
            "active": [
              {
                id:"ShowData",
                name:"configComponent.ChartPublic.ShowData",
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
            "dataBind": [],
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
                "w": 342,
                "h": 300
              },
              "backColor": "transparent",
              "zIndex": 1,
              "transform": 0,
              "diy":[
                {
                  "name":"configComponent.ChartPublic.splitNumber",
                  "type":1,
                  "value":10,
                  "key":"splitNumber",
                },
                {
                  "name":"configComponent.ChartPublic.ChartTitle",
                  "type":4,
                  "value":"Title",
                  "key":"ChartTitle",
                },
                {
                  "name":"configComponent.ChartPublic.ChartUnit",
                  "type":4,
                  "value":"%",
                  "key":"ChartUnit",
                },
                {
                  "name":"configComponent.ChartPublic.ChartTitleFontSize",
                  "type":1,
                  "value":20,
                  "key":"ChartTitleFontSize",
                },
                {
                  "name":"configComponent.ChartPublic.ChartTitleFontColor",
                  "type":2,
                  "value":"#000000",
                  "key":"ChartTitleFontColor",
                },
                {
                  "name":"configComponent.ChartPublic.ChartMinValue",
                  "type":7,
                  "value":0,
                  "key":"ChartMinValue",
                },
                {
                  "name":"configComponent.ChartPublic.ChartMaxValue",
                  "type":7,
                  "value":100,
                  "key":"ChartMaxValue",
                },
                {
                  "name":"configComponent.ChartPublic.ChartAxisTickColor",
                  "type":2,
                  "value":"#000000",
                  "key":"ChartAxisTickColor",
                },
                {
                  "name":"configComponent.ChartPublic.ChartWidth",
                  "type":7,
                  "value":20,
                  "key":"ChartWidth",
                },
                {
                  "name":"configComponent.ChartPublic.LabelDis",
                  "type":7,
                  "value":25,
                  "key":"LabelDis",
                },
                {
                  "name":"configComponent.ChartPublic.ChartSplitLineWidth",
                  "type":7,
                  "value":5,
                  "key":"ChartSplitLineWidth",
                },
                {
                  "name":"configComponent.ChartPublic.ChartSplitLineHeight",
                  "type":7,
                  "value":5,
                  "key":"ChartSplitLineHeight",
                },
                {
                  "name":"configComponent.ChartPublic.Area1Range",
                  "type":4,
                  "value":"0~30",
                  "key":"Area1Range",
                },
                {
                  "name":"configComponent.ChartPublic.Area1Color",
                  "type":2,
                  "value":"#4dabf7",
                  "key":"Area1Color",
                },
                {
                  "name":"configComponent.ChartPublic.Area2Range",
                  "type":4,
                  "value":"30~60",
                  "key":"Area2Range",
                },
                {
                  "name":"configComponent.ChartPublic.Area2Color",
                  "type":2,
                  "value":"#69db7c",
                  "key":"Area2Color",
                },
                {
                  "name":"configComponent.ChartPublic.Area3Color",
                  "type":2,
                  "value":"#ff6b6b",
                  "key":"Area3Color",
                },
              ]
            }
          }
        },
        strokeColor:"#000000",
        fill:"#A1BFE2",
        strokeWidth:0.3,
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
        echart: null,
        eventValue: '0.00',
        eventUnit: '',
        option: {
          tooltip: {
                formatter: "{a} <br/>{b} : {c}%"
            },
          series: [
              {
                type: 'gauge',
                min:10,
                max:100,
                animation:true,
                splitNumber: 9,
                radius: "110%",
                center : ['50%', '58%'],    // 默认全局居中
                title: {
                  show: true,
                  textStyle: {
                    fontSize: 10,
                    color: '#ffffff'
                  }
                },
                itemStyle: {
                  color: '#0EA80E'
                },
                axisTick: {
                  distance: 0,
                  lineStyle: {
                    width: 2,
                    color: 'auto',
                  },
                },
                axisLine: {
                  lineStyle: {
                    width: 20,
                    color: [
                      [0.2, "#4dabf7"],
                      [0.65, "#69db7c"],
                      [1, "#ff6b6b"]
                    ]
                  }
                },
                pointer: {
                  itemStyle: {
                    color: 'auto'
                  }
                },
                splitLine: {
                  distance: -10,
                  show:true,
                  lineStyle: {
                    width: 2,
                    color: 'auto'
                  }
                },
                axisLabel: {
                  splitNumber: 20,
                  show: true,
                  distance: 25,
                  color: "#ffffff",
                  fontSize:39
                },
                tooltip: {
                  show:false,
                },
                detail: {
                  formatter: '{value}%',
                  color: "#df1f1f",
                  fontSize: 12,
                },
                data: [
                    {
                    value: 0,
                    name: ''
                  }
                ]
             }
          ]
        }
      }
  },
  methods: {
    initComponents(option){
      if(this.IsToolBox)
      {
        return
      }
      // 确保 option 存在
      if (!option || !option.style) {
        console.warn('Chart initComponents: option or option.style is undefined')
        return
      }
      let refObj = this.detail && this.detail.identifier ? this.detail.identifier : 'chart_' + Date.now()
      let view = this.$refs[refObj]
      // 确保 DOM 元素存在
      if (!view) {
        console.warn('Chart initComponents: cannot find DOM element with ref:', refObj)
        return
      }
      let i=0
      let minrang=0
      for( i=0;i<option.style.diy.length;i++)
      {
        if(option.style.diy[i].key=="ChartTitle")
        {
          this.option.series[0].data[0].name=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="ChartUnit")
        {
          this.option.series[0].detail.formatter='{value}'+option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="ChartTitleFontSize")
        {
          this.option.series[0].title.textStyle.fontSize = option.style.diy[i].value
          this.option.series[0].detail.fontSize = option.style.diy[i].value
          this.option.series[0].axisLabel.fontSize = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="ChartTitleFontColor")
        {
          this.option.series[0].title.textStyle.color = option.style.diy[i].value
          this.option.series[0].detail.color = option.style.diy[i].value
          this.option.series[0].axisLabel.color = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="ChartMinValue")
        {
          this.option.series[0].min = parseInt(option.style.diy[i].value)
        }
        else if(option.style.diy[i].key=="ChartMaxValue")
        {
          this.option.series[0].max = parseInt(option.style.diy[i].value)
        }
        else if(option.style.diy[i].key=="ChartAxisTickColor")
        {
          // this.option.series[0].axisTick.lineStyle.color = option.style.diy[i].value
          // this.option.series[0].splitLine.lineStyle.color = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="ChartSplitLineColor")
        {
          // this.option.series[0].splitLine.lineStyle.color = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="ChartWidth")
        {
          this.option.series[0].axisLine.lineStyle.width = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="LabelDis")
        {
          this.option.series[0].axisLabel.distance = parseInt(option.style.diy[i].value)
        }
        else if(option.style.diy[i].key=="ChartSplitLineWidth")
        {
          this.option.series[0].splitLine.lineStyle.width = option.style.diy[i].value
          this.option.series[0].axisTick.lineStyle.width = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="ChartSplitLineHeight")
        {
          this.option.series[0].splitLine.length = parseInt(option.style.diy[i].value)+20
          this.option.series[0].axisTick.length = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="Area1Range")
        {
          let range = option.style.diy[i].value.split("~")
          minrang = range[0]
          this.option.series[0].axisLine.lineStyle.color[0][0] = (parseInt(range[1])-parseInt(range[0]))/(this.option.series[0].max-this.option.series[0].min)
        }
        else if(option.style.diy[i].key=="Area2Range")
        {
          let range2 = option.style.diy[i].value.split("~")
          this.option.series[0].axisLine.lineStyle.color[1][0] = (parseInt(range2[1])-parseInt(minrang))/(this.option.series[0].max-this.option.series[0].min)
          this.option.series[0].axisLine.lineStyle.color[2][0] = 1
        }
        else if(option.style.diy[i].key=="Area1Color")
        {
          this.option.series[0].axisLine.lineStyle.color[0][1] = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="Area2Color")
        {
          this.option.series[0].axisLine.lineStyle.color[1][1] = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="Area3Color")
        {
          this.option.series[0].axisLine.lineStyle.color[2][1] = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="splitNumber")
        {
          this.option.series[0].splitNumber = parseInt(option.style.diy[i].value)
        }
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

      if (!this.echartsView) {
       this.echartsView = echarts.init(view, null);
      }
      else
      {
       this.echartsView.resize()
      }
      let _t = this
      // this.option.series[0].itemStyle.normal.color = this.progressColor
      setTimeout(function (){
        _t.echartsView.setOption(_t.option,true)
        _t.echartsView.resize()
      }, 100)
    },
    onResize() {
        if (this.echartsView) {
           this.echartsView.resize();
        }
    },
    updateView() {
          this.setOption(this.option);
      },
  },
  beforeDestroy () {
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
    });
    node.on('change:visible', ({ current }) => {
      if(current)
      {
        _t.option.series[0].data[0].value =_t.ShowDataResult
        _t.initComponents(_t.detail);
      }
    });
    this.detail = node.getData().detail
    this.editMode = node.getData().editMode
    this.showDeviceUuid = node.getData().showDeviceUuid
    this.IsToolBox = node.getData().IsToolBox
    _t.$EventBus.$on('cell-editMode', (data) => {
      _t.editMode = data.edit
      _t.IsToolBox = data.toolbox
	  // _t.initComponents(_t.detail)
    })
  },
  mounted() {
      let _t = this
      this.$nextTick(function(){
        this.initComponents(this.detail);
        let activeEvent = this.detail.identifier+"activeEvent"//动作数据
        let animateEvent = this.detail.identifier+"animateEvent"//动作数据

        _t.$EventBus.$on(activeEvent, (data) => {
          if(data.ID == "ShowData")
          {
            this.ShowDataResult = data.result
            this.option.series[0].data[0].value = data.result
            this.echartsView.setOption(this.option, true);
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
.view-chart-gauge {
    height: 100%;
    width: 100%;
}
</style>
