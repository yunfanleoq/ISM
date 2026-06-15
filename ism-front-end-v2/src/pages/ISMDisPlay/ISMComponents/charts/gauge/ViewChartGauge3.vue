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
let chartUnit=""
export default {
  mixins: [ISMChildAutoMixin],
  name: 'ism-view-chart-gauge-3',
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
          "text": "configComponent.chartGauge.levelGauge",
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
                  "name":"configComponent.ChartPublic.splitNumber",
                  "type":1,
                  "value":10,
                  "key":"splitNumber",
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
        option: {
          tooltip: {
            formatter: '{a} <br/>{b} : {c}%'
          },
          series: [{
            name: 'Pressure',
            type: 'gauge',
            startAngle: 200,
            endAngle: -20,
            min: 0,
            max: 1,
            splitNumber: 8,
            axisLine: {
              roundCap: true,
              lineStyle: {
                width: 6,
                color: [
                  [0.25, '#FF6E76'],
                  [0.5, '#FDDD60'],
                  [0.75, '#58D9F9'],
                  [1, '#7CFFB2']
                ]
              }
            },
            progress: {
              show: true,
              roundCap: true,
              width: 18
            },
            itemStyle: {
              color: '#58D9F9',
              shadowColor: 'rgba(0,138,255,0.45)',
              shadowBlur: 10,
              shadowOffsetX: 2,
              shadowOffsetY: 2
            },
            pointer: {
              icon: 'path://M12.8,0.7l12,40.1H0.7L12.8,0.7z',
              length: '12%',
              width: 20,
              offsetCenter: [0, '-60%'],
              itemStyle: {
                color: 'auto'
              }
            },
            axisTick: {
              length: 12,
              lineStyle: {
                color: 'auto',
                width: 2
              }
            },
            splitLine: {
              length: 20,
              lineStyle: {
                color: 'auto',
                width: 5
              }
            },
            radius: "110%",
            center : ['50%', '72%'],    // 默认全局居中
            title: {
              show: true,
              offsetCenter: [0, '-20%'],
              textStyle: {
                fontSize: 10,
                color: '#ffffff'
              }
            },
            axisLabel: {
              color: '#464646',
              fontSize: 20,
              distance: -60,
              formatter: function (value) {
                if (value === 0.875) {
                  return 'A';
                } else if (value === 0.625) {
                  return 'B';
                } else if (value === 0.375) {
                  return 'C';
                } else if (value <= 0.125) {
                  return 'D';
                }
                return '';
              }
            },
            detail: {
              fontSize: 50,
              offsetCenter: [0, '0%'],
              valueAnimation: true,
              formatter: function (value) {
                return Math.round(value*100) + '分';
              },
              color: 'auto',
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
              }
            },
            data: [
              {
                value: 0,
                name: ''
              }
            ]
          }]
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
      for( i=0;i<option.style.diy.length;i++)
      {
        if(option.style.diy[i].key=="ChartTitle")
        {
          this.option.series[0].data[0].name=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="ChartUnit")
        {
          chartUnit = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="ChartTitleFontSize")
        {
          this.option.series[0].title.textStyle.fontSize = option.style.diy[i].value
          this.option.series[0].detail.fontSize = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="ChartTitleFontColor")
        {
          this.option.series[0].title.textStyle.color = option.style.diy[i].value
          this.option.series[0].detail.color = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="ChartWidth")
        {
          this.option.series[0].axisLine.lineStyle.width = parseInt(option.style.diy[i].value)
        }
        else if(option.style.diy[i].key=="LabelDis")
        {
          this.option.series[0].axisLabel.distance = parseInt(option.style.diy[i].value)
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
