<template>
  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"  x="0px" y="0px"  xml:space="preserve" :style="{'overflow': 'visible','width':(detail && detail.style && detail.style.position && detail.style.position.w) || '200px','height':(detail && detail.style && detail.style.position && detail.style.position.h) || '200px'}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <foreignObject style="overflow:visible;" pointer-events="all" :width="((detail && detail.style && detail.style.position && detail.style.position.w) || 200) - ((detail && detail.style && detail.style.borderWidth) || 0)*2" :height="((detail && detail.style && detail.style.position && detail.style.position.h) || 200) - ((detail && detail.style && detail.style.borderWidth) || 0)*2">
        <div class="view-chart-gauge" :ref="detail && detail.identifier ? detail.identifier : 'chart_default'" :style="{'overflow': 'visible','width':'100%','height':'100%',}">
            <canvas data-type="linear-gauge" :ref="detail && detail.identifier ? detail.identifier : 'chart_default'" :width="((detail && detail.style && detail.style.position && detail.style.position.w) || 200) - ((detail && detail.style && detail.style.borderWidth) || 0)*2" :height="((detail && detail.style && detail.style.position && detail.style.position.h) || 200) - ((detail && detail.style && detail.style.borderWidth) || 0)*2"

            ></canvas>
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
import BaseView from '../../View';
import { LinearGauge } from 'canvas-gauges'
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
  name: 'ism-view-measuring',
  inject: ['getNode'],
  props: {

  },
  components: {

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
        base:{
          "text": "configComponent.measuring.title",
          "icon": "icon-liangchi",
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
                "w": 170,
                "h": 360
              },
              "backColor": "transparent",
              "foreColor": "#000000",
              fontFamily: "Arial",
              fontSize: 18,
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
                  "value":"%",
                  "key":"ChartUnit",
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
                  name:"configComponent.TemperatureGauge.TicksPosition",
                  type:6,
                  value:0,
                  enumList:[
                    {
                      value:0,
                      option:"configComponent.TemperatureGauge.TicksPositionLeft"
                    },
                    {
                      value:1,
                      option:"configComponent.TemperatureGauge.TicksPositionRight"
                    }
                  ],
                  min:1,
                  key:"TicksPosition",
                },
                {
                  "name":"configComponent.TemperatureGauge.colorMajorTicks",
                  "type":2,
                  "value":"#000000",
                  "key":"colorMajorTicks",
                },
                {
                  "name":"configComponent.TemperatureGauge.colorMinorTicks",
                  "type":2,
                  "value":"#ADA9A9",
                  "key":"colorMinorTicks",
                },
                {
                  "name":"configComponent.TemperatureGauge.colorBarProgress",
                  "type":2,
                  "value":"#327ac0",
                  "key":"colorBarProgress",
                },
                {
                  "name":"configComponent.TemperatureGauge.colorBar",
                  "type":2,
                  "value":"#f5f5f5",
                  "key":"colorBar",
                },
                {
                  "name":"configComponent.TemperatureGauge.LowValue",
                  "type":7,
                  "value":-10,
                  "key":"LowValue",
                },
                {
                  "name":"configComponent.TemperatureGauge.LowColor",
                  "type":2,
                  "value":"yellow",
                  "key":"LowColor",
                },
                {
                  "name":"configComponent.TemperatureGauge.HighValue",
                  "type":7,
                  "value":65,
                  "key":"HighValue",
                },
                {
                  "name":"configComponent.TemperatureGauge.HighColor",
                  "type":2,
                  "value":"red",
                  "key":"HighColor",
                },
                {
                  "name":"configComponent.TemperatureGauge.barWidth",
                  "type":7,
                  "value":10,
                  "key":"barWidth",
                },
                {
                  "name":"configComponent.TemperatureGauge.majorTicksStep",
                  "type":1,
                  "value":10,
                  "min":1,
                  "key":"majorTicksStep",
                },
                {
                  "name":"configComponent.TemperatureGauge.minorTicksStep",
                  "type":1,
                  "value":5,
                  "min":1,
                  "key":"minorTicksStep",
                },
                {
                  "name":"configComponent.TemperatureGauge.isShowValue",
                  "type":6,
                  "enumList":[
                    {option:'True',value:1},{option:'False',value:0}
                  ],
                  "value":1,
                  "key":"isShowValue",
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
        isShowValue:1,
        eventValue: '0.00',
        eventUnit: '',
        option: {
          renderTo: null,
          width: 0,
          height:0,
          units: "",

          minValue: 0,
          startAngle: 90,
          ticksAngle: 180,
          valueBox: false,
          valueBoxStroke:0,
          colorValueBoxShadow:false,
          colorValueBoxBackground:false,
          valueDec:2,
          valueInt:2,
          maxValue: 220,
          majorTicksStep:10,
          majorTicks: [0,20,30,50,60,80,90,100],
          minorTicks: 10,
          strokeTicks: true,
          highlights: [
            {
              "from": 60,
              "to": 100,
              "color": "rgba(200, 50, 50, .75)"
            },
            {
              "from": 10,
              "to": 20,
              "color": "#000000"
            },
          ],
          colorPlate: "",
          borderShadowWidth: 0,
          borders: false,
          needleType: "arrow",
          needleWidth: 5,
          needleCircleSize: 7,
          needleCircleOuter: true,
          needleCircleInner: false,
          animationDuration: 1500,
          animationRule: "linear",
          barBeginCircle:false,
          colorMajorTicks: "#Bfe66a",
          colorMinorTicks: "#ffe66a",
          colorTitle: "#eee",
          colorUnits: "#ccc",
          colorNumbers: "#eee",
          colorBarProgress: "#327ac0",
          colorBar: "#f5f5f5",
          animateOnInit: true,
          animatedValue: true,
          barWidth: 10,
          fontNumbers: "Arial",
          fontTitle: "Arial",
          fontUnits: "Arial",
          fontValue: "Arial",
          tickSide: "left",
          numberSide: "left",
          needleSide: "left",
          fontNumbersSize: 20,
          fontTitleSize: 24,
          fontUnitsSize: 22,
          fontValueSize: 26,
          value: 0
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
      this.option.renderTo = view
      this.option.width=this.detail.style.position.w
      this.option.height=this.detail.style.position.h
      this.option.colorPlate=this.detail.style.backColor
      this.option.colorTitle = this.detail.style.foreColor
      this.option.colorUnits = this.detail.style.foreColor
      this.option.colorNumbers = this.detail.style.foreColor
      this.option.strokeTicks=true
      this.option.fontNumbers= this.detail.style.fontFamily
      this.option.fontTitle= this.detail.style.fontFamily
      this.option.fontUnits= this.detail.style.fontFamily
      this.option.fontValue= this.detail.style.fontFamily

      this.option.fontNumbersSize= this.detail.style.fontSize
      this.option.fontTitleSize= this.detail.style.fontSize
      this.option.fontUnitsSize= this.detail.style.fontSize
      this.option.fontValueSize= this.detail.style.fontSize
      for( i=0;i<option.style.diy.length;i++)
      {
        if(option.style.diy[i].key=="ChartTitle")
        {
          this.option.title=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="ChartUnit")
        {
          this.eventUnit = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="ChartMinValue")
        {
          this.option.minValue = parseFloat(option.style.diy[i].value)
        }
        else if(option.style.diy[i].key=="ChartMaxValue")
        {
          this.option.maxValue = parseFloat(option.style.diy[i].value)
        }
        else if(option.style.diy[i].key=="colorMajorTicks")
        {
          this.option.colorMajorTicks = parseInt(option.style.diy[i].value)
        }
        else if(option.style.diy[i].key=="colorMinorTicks")
        {
          this.option.colorMinorTicks = parseInt(option.style.diy[i].value)
        }
        else if(option.style.diy[i].key=="colorBarProgress")
        {
          this.option.colorBarProgress = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="colorBar")
        {
          this.option.colorBar = option.style.diy[i].value
        }

        else if(option.style.diy[i].key=="LowValue")
        {
          this.option.highlights[0].from =  parseFloat(this.option.minValue)
          this.option.highlights[0].to =  parseFloat(option.style.diy[i].value)
        }
        else if(option.style.diy[i].key=="LowColor")
        {
          this.option.highlights[0].color = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="HighValue")
        {
          this.option.highlights[1].from = parseFloat(option.style.diy[i].value)
          this.option.highlights[1].to = parseFloat(this.option.maxValue)
        }
        else if(option.style.diy[i].key=="HighColor")
        {
          this.option.highlights[1].color = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="TicksPosition")
        {
          let position = ""
          if(option.style.diy[i].value==0)
          {
            position = "left"
          }
          else if(option.style.diy[i].value==1)
          {
            position = "right"
          }
          this.option.tickSide = position
          this.option.numberSide = position
          this.option.needleSide = position
        }
        else if(option.style.diy[i].key=="barWidth")
        {
          this.option.barWidth = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="majorTicksStep")
        {
          this.majorTicksStep = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="minorTicksStep")
        {
          this.option.minorTicks = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="isShowValue")
        {
          this.isShowValue = option.style.diy[i].value
        }
      }
      if(this.isShowValue)
      {
        this.option.units = this.option.value+this.eventUnit
      }
      else
      {
        this.option.units = ""
      }
      this.option.majorTicks=[]
      for(let k=this.option.minValue;k<=this.option.maxValue;k=k+parseInt(this.majorTicksStep))
      {
        this.option.majorTicks.push(k)
      }
      if(this.option.majorTicks[this.option.majorTicks.length-1]!=this.option.maxValue)
      {
        this.option.majorTicks.push(this.option.maxValue)
      }

      this.echartsView = new LinearGauge(this.option).draw();
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
    },
    onResize() {

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
            this.echartsView.update({ value:parseFloat(data.result)})
            if(this.isShowValue)
            {
              this.echartsView.update ({ units:parseFloat(data.result)+this.eventUnit})
            }
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
