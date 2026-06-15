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
import * as echarts from 'echarts';
import BaseView from '../../View';
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
export default {
  mixins: [ISMChildAutoMixin],
  name: 'ism-view-chart-gauge-4',
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
          "text": "configComponent.chartGauge.temperatureGauge",
          "icon": "data:image/png;base64,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",
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
                  "name":"configComponent.ChartPublic.ChartTitle",
                  "type":4,
                  "value":"Title",
                  "key":"ChartTitle",
                },
                {
                  "name":"configComponent.ChartPublic.ChartUnit",
                  "type":4,
                  "value":"°C",
                  "key":"ChartUnit",
                },
                {
                  "name":"configComponent.ChartPublic.ChartTitleFontSize",
                  "type":1,
                  "value":20,
                  "key":"ChartTitleFontSize",
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
        eventValue: '0.00',
        eventUnit: '',
        chartUnit:"",
        option: {
          tooltip: {
            formatter: '{a} <br/>{b} : {c}%'
          },
          series: [
            {
              type: 'gauge',
              radius: "100%",
              center: ['50%', '72%'],
              startAngle: 200,
              endAngle: -20,
              min: 0,
              max: 60,
              splitNumber: 12,
              itemStyle: {
                color: '#FFAB91'
              },
              progress: {
                show: true,
                width: 30
              },
              pointer: {
                show: false
              },
              axisLine: {
                lineStyle: {
                  width: 30
                }
              },
              axisTick: {
                distance: -45,
                splitNumber: 5,
                lineStyle: {
                  width: 2,
                  color: '#999'
                }
              },
              splitLine: {
                distance: -52,
                length: 14,
                lineStyle: {
                  width: 3,
                  color: '#999'
                }
              },
              axisLabel: {
                distance: -10,
                color: '#999',
                fontSize: 20,
                formatter: function (val) {
                  //解决刻度的值为浮点数问题
                  return Math.ceil(val);
                },
              },
              anchor: {
                show: false
              },
              title: {
                show: false
              },
              detail: {
                valueAnimation: true,
                width: '60%',
                lineHeight: 40,
                borderRadius: 8,
                offsetCenter: [0, '-15%'],
                fontSize: 60,
                fontWeight: 'bolder',
                formatter: null,
                rich: {
                  value: {
                    fontSize: 50,
                    fontWeight: 'bolder',
                    color: '#777'
                  },
                  unit: {
                    fontSize: 20,
                    color: '#999',
                    padding: [0, 0, -20, 10]
                  }
                },
                color: 'auto'
              },
              data: [
                {
                  value: 20
                }
              ]
            },
            {
              type: 'gauge',
              radius: "100%",
              center: ['50%', '72%'],
              startAngle: 200,
              endAngle: -20,
              min: 0,
              max: 60,
              itemStyle: {
                color: '#FD7347'
              },
              progress: {
                show: true,
                width: 8
              },
              pointer: {
                show: false
              },
              axisLine: {
                show: false
              },
              axisTick: {
                show: false
              },
              splitLine: {
                show: false
              },
              axisLabel: {
                show: false
              },
              detail: {
                show: false
              },
              data: [
                {
                  value: 20
                }
              ]
            }
          ]
        }
      }
  },
  methods: {
    formatter(value){
      return '{value|' + value.toFixed(0) + '}{unit|'+this.chartUnit+'}';
    },
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
      this.option.series[0].detail.formatter = this.formatter
      for( i=0;i<option.style.diy.length;i++)
      {
        if(option.style.diy[i].key=="ChartTitle")
        {
          this.option.series[0].data[0].name=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="ChartUnit")
        {
          this.chartUnit = option.style.diy[i].value
          this.option.tooltip.formatter = '{a} <br/>{b} : {c}'+this.chartUnit
        }
        else if(option.style.diy[i].key=="ChartTitleFontSize")
        {
          this.option.series[0].detail.fontSize = option.style.diy[i].value
          this.option.series[0].axisLabel.fontSize = option.style.diy[i].value
          this.option.series[0].axisLabel.distance = -(parseInt(option.style.diy[i].value)-10)
        }
        else if(option.style.diy[i].key=="ChartMinValue")
        {
          this.option.series[0].min = option.style.diy[i].value
          this.option.series[1].min = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="ChartMaxValue")
        {
          this.option.series[0].max = option.style.diy[i].value
          this.option.series[1].max = option.style.diy[i].value
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
      // _t.initComponents(_t.detail);
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
