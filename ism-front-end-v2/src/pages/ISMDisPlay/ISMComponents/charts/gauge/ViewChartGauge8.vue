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
let chartData = 100;
export default {
  mixins: [ISMChildAutoMixin],
  name: 'ism-view-chart-gauge-8',
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
      tempOption:null,
        base:{
          "text": "configComponent.oil.title",
          "icon": "icon-yiwuyaodian-",
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
                "w": 220,
                "h": 380
              },
              "backColor": "transparent",
              foreColor:"#FFFF00",
              fontSize:50,
              "zIndex": 1,
              "transform": 0,
              "diy":[
                {
                  "name":"configComponent.ChartPublic.ChartUnit",
                  "type":4,
                  "value":"cm",
                  "key":"ChartUnit",
                },
                {
                  "name":"configComponent.ChartPublic.ChartMaxValue",
                  "type":7,
                  "value":200,
                  "key":"ChartMaxValue",
                },
              ]
            }
          }
        },
        strokeColor:"#000000",
        chartData:100,
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
        MaxValue:100,
        eventValue: '0.00',
        eventUnit: '',
        chartUnit:"",
        option: {
          tooltip: {
            trigger: 'none'
          },
          grid:{
            top:20,
            bottom:20
          },
          xAxis: {
            data: ["百分比"],
            axisTick: {
              show: false
            },
            axisLine: {
              show: false
            },
            axisLabel: {
              show: false,
              textStyle: {
                color: '#e54035'
              }
            }
          },
          yAxis: {
            splitLine: {
              show: false
            },
            axisTick: {
              show: false
            },
            axisLine: {
              show: false
            },
            axisLabel: {
              show: false
            }
          },
          series: [
            {
            name: '最上层立体圆',
            type: 'pictorialBar',
            symbolSize: [100, 45],
            symbolOffset: [0, -20],
            z: 12,
            itemStyle: {
              normal: {
                color: '#363F5E'
              }
            },
            data: [{
              value: 100,
              symbolPosition: 'end'
            }]
          },
            {
            name: '中间立体圆',
            type: 'pictorialBar',
            symbolSize: [100, 45],
            symbolOffset: [0, -20],
            z: 12,
            itemStyle: {
              normal: {
                color: '#FFFF00',
                opacity: 1
              }
            },
            data: [{
              value: chartData,
              symbolPosition: 'end'
            }]
          },
            {
              name: '最底部立体圆',
              type: 'pictorialBar',
              symbolSize: [100, 45],
              symbolOffset: [0, 20],
              z: 12,
              itemStyle: {
                normal: {
                  color: '#FFFF00',
                }
              },
              data: [100 - chartData]
          },
            {
            //底部立体柱
            stack: '1',
            type: 'bar',
            itemStyle: {
              normal: {
                color: '#FFFF00',
                opacity: .7
              }
            },
            label: {
              show: true,
              position: "inside",
              distance: 20,
              color: "#FFFE00",
              fontSize:50,
              formatter:'{c}'+'%'
            },
              clip:false,
              barWidth:100,
            data: [chartData]
          },
            {
            //上部立体柱
            stack: '1',
            type: 'bar',
            itemStyle: {
              normal: {
                color: '#36405E',
                opacity: .7
              }
            },
              clip:false,
              barWidth:100,
            data: [100 - chartData]
          }
          ]
        }
      }
  },
  beforeDestroy () {
    if (this.echartsView != null&&(typeof this.echartsView.dispose=="function")) {
      this.echartsView.dispose()
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
      for( i=0;i<option.style.diy.length;i++)
      {
        if(option.style.diy[i].key=="ChartUnit")
        {
          this.chartUnit = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="ChartMaxValue")
        {
          this.MaxValue = parseFloat(option.style.diy[i].value)
          this.option.series[0].data[0].value = parseFloat(option.style.diy[i].value)
          this.option.series[2].data[0] = parseFloat(option.style.diy[i].value)-chartData
          this.option.series[4].data[0] = parseFloat(option.style.diy[i].value)-chartData
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
      this.option.series[0].barWidth = this.detail.style.position.w
      this.option.series[1].barWidth = this.detail.style.position.w
      this.option.series[2].barWidth = this.detail.style.position.w
      this.option.series[0].symbolSize[0] = this.detail.style.position.w
      this.option.series[1].symbolSize[0] = this.detail.style.position.w
      this.option.series[2].symbolSize[0] = this.detail.style.position.w
      this.option.series[3].barWidth= this.detail.style.position.w
      this.option.series[4].barWidth= this.detail.style.position.w

      this.option.series[1].itemStyle.normal.color = this.detail.style.foreColor
      this.option.series[2].itemStyle.normal.color = this.detail.style.foreColor
      this.option.series[3].itemStyle.normal.color = this.detail.style.foreColor
      this.option.series[3].label.color = this.detail.style.foreColor
      this.option.series[3].label.fontSize = this.detail.style.fontSize
      this.option.series[3].label.formatter = '{c}'+this.chartUnit

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

        _t.initComponents(_t.detail);
        if(_t.tempOption!=null)
        {
          _t.echartsView.setOption(_t.tempOption, true);
          this.echartsView.resize();
        }
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
            if(data.result>_t.MaxValue)
            {
              data.result = _t.MaxValue
            }
            _t.option.series[3].data[0] = data.result
            _t.option.series[1].data[0].value = data.result
            _t.option.series[4].data[0] = _t.MaxValue-parseFloat(data.result)
            if(data.result/_t.MaxValue>=0.6)
            {
              _t.option.series[3].label.position = "inside"
            }
            else
            {
              _t.option.series[3].label.position = "top"
            }
            _t.tempOption = _t.option
            _t.echartsView.setOption(_t.option, true);
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
