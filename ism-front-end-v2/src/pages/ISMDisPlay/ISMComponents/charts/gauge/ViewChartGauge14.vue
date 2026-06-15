<template>
  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"  x="0px" y="0px"  xml:space="preserve" :style="{'overflow': 'visible','width':(detail && detail.style && detail.style.position && detail.style.position.w) || '200px','height':(detail && detail.style && detail.style.position && detail.style.position.h) || '200px'}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <foreignObject style="overflow:visible;" pointer-events="all" :width="((detail && detail.style && detail.style.position && detail.style.position.w) || 200) - ((detail && detail.style && detail.style.borderWidth) || 0)*2" :height="((detail && detail.style && detail.style.position && detail.style.position.h) || 200) - ((detail && detail.style && detail.style.borderWidth) || 0)*2">
        <div class="view-chart-gauge" :ref="detail.identifier" :style="{'overflow': 'visible'}">
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
import 'echarts-liquidfill'
import BaseView from '../../View';
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
export default {
  mixins: [ISMChildAutoMixin],
  name: 'ism-view-chart-gauge-14',
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
      base:{
        "text": "configComponent.chartGauge.CircularFlow",
        "icon": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAARQAAAERCAYAAAC6vVjKAAAACXBIWXMAAA7EAAAOxAGVKw4bAAAE7WlUWHRYTUw6Y29tLmFkb2JlLnhtcAAAAAAAPD94cGFja2V0IGJlZ2luPSLvu78iIGlkPSJXNU0wTXBDZWhpSHpyZVN6TlRjemtjOWQiPz4gPHg6eG1wbWV0YSB4bWxuczp4PSJhZG9iZTpuczptZXRhLyIgeDp4bXB0az0iQWRvYmUgWE1QIENvcmUgNS42LWMxNDggNzkuMTY0MDM2LCAyMDE5LzA4LzEzLTAxOjA2OjU3ICAgICAgICAiPiA8cmRmOlJERiB4bWxuczpyZGY9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkvMDIvMjItcmRmLXN5bnRheC1ucyMiPiA8cmRmOkRlc2NyaXB0aW9uIHJkZjphYm91dD0iIiB4bWxuczp4bXA9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC8iIHhtbG5zOmRjPSJodHRwOi8vcHVybC5vcmcvZGMvZWxlbWVudHMvMS4xLyIgeG1sbnM6cGhvdG9zaG9wPSJodHRwOi8vbnMuYWRvYmUuY29tL3Bob3Rvc2hvcC8xLjAvIiB4bWxuczp4bXBNTT0iaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wL21tLyIgeG1sbnM6c3RFdnQ9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC9zVHlwZS9SZXNvdXJjZUV2ZW50IyIgeG1wOkNyZWF0b3JUb29sPSJBZG9iZSBQaG90b3Nob3AgMjEuMCAoV2luZG93cykiIHhtcDpDcmVhdGVEYXRlPSIyMDIzLTAyLTEzVDA5OjU4OjA4KzA4OjAwIiB4bXA6TW9kaWZ5RGF0ZT0iMjAyMy0wMi0xM1QxMDowMTo0NiswODowMCIgeG1wOk1ldGFkYXRhRGF0ZT0iMjAyMy0wMi0xM1QxMDowMTo0NiswODowMCIgZGM6Zm9ybWF0PSJpbWFnZS9wbmciIHBob3Rvc2hvcDpDb2xvck1vZGU9IjMiIHhtcE1NOkluc3RhbmNlSUQ9InhtcC5paWQ6MWUwNzEwYzYtM2ViZS04YjQzLTkyODItYzYwNjI3MzBkMWNiIiB4bXBNTTpEb2N1bWVudElEPSJ4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",
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
              "w": 150,
              "h": 150
            },
            "backColor": "transparent",
            foreColor:"#ffffff",
            fontSize:40,
            fontFamily:"数字字体-2",
            "zIndex": 1,
            "transform": 0,
            "diy":[
              {
                "name":"configComponent.ChartPublic.ChartUnit",
                "type":4,
                "value":"%",
                "key":"ChartUnit",
              },
              {
                "name":"configComponent.CircularFlow.progressMax",
                "type":1,
                "value":100,
                "min":1,
                "key":"progressMax",
              },
              {
                "name":"configComponent.CircularFlow.HighColor",
                "type":2,
                "value":"#446bf5",
                "key":"HighColor",
              },
              {
                "name":"configComponent.CircularFlow.LowColor",
                "type":2,
                "value":"#2ca3e2",
                "key":"progressColor",
              },
              {
                "name":"configComponent.CircularFlow.surplusColor",
                "type":2,
                "value":"#0F224C",
                "key":"surplusColor",
              }
            ]
          }
        }
      },
      progressMax:100,
      strokeColor:"#000000",
      fill:"#A1BFE2",
      strokeWidth:0.3,
      fillOpacity:1,
      tempValue:"default",
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
      echart: null,
      eventValue: '0.00',
      eventUnit: '',
      option:{
        backgroundColor: '#0F224C',
        series: [
          {
            type: 'liquidFill',
            radius: '100%',
            center: ['50%', '50%'],
            color: [
              {
                type: 'linear',
                x: 0,
                y: 0,
                x2: 0,
                y2: 1,
                colorStops: [
                  {
                    offset: 0,
                    color: '#446bf5',
                  },
                  {
                    offset: 1,
                    color: '#2ca3e2',
                  }
                ],
                globalCoord: false,
              },
            ],
            data: [0.45,0.45,0.45,0.45], // data个数代表波浪数
            backgroundStyle: {
              borderWidth: 1,
              color: 'RGBA(51, 66, 127, 0.7)',
            },
            label: {
              normal: {
                textStyle: {
                  fontSize: 50,
                  color: '#fff',
                },
              },
            },
            outline: {
              show: false,
              borderDistance: 0,
              itemStyle: {
                borderWidth: 2,
                borderColor: '#112165',
              },
            },
          },
        ],
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
        if(option.style.diy[i].key=="ChartUnit")
        {
          this.chartUnit = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="progressColor")
        {
          this.progressColor = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="HighColor")
        {
          this.option.series[0].color[0].colorStops[0].color = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="surplusColor")
        {
          this.surplusColor = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="ChartWidth")
        {
          this.option.series[0].barWidth = parseInt(option.style.diy[i].value)
        }
        else if(option.style.diy[i].key=="progressMax")
        {
          this.progressMax = parseFloat(option.style.diy[i].value)
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

      this.option.backgroundColor = this.detail.style.backColor

      this.option.series[0].label.normal.formatter = '45'+this.chartUnit
      this.option.series[0].label.normal.textStyle.color = this.detail.style.foreColor
      this.option.series[0].label.normal.textStyle.fontSize = this.detail.style.fontSize
      this.option.series[0].label.normal.textStyle.fontFamily = this.detail.style.fontFamily
      this.option.series[0].color[0].colorStops[1].color = this.progressColor

      this.option.series[0].backgroundStyle.color = this.surplusColor
      // this.option.series[0].itemStyle.normal.color = this.progressColor
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
          if( _t.tempValue !=data.result)
          {
            let value = parseFloat(data.result)/this.progressMax
            console.log(value)
            _t.option.series[0].label.normal.formatter = data.result+this.chartUnit
            _t.option.series[0].data[0] = value
            _t.option.series[0].data[1] = value
            _t.option.series[0].data[2] = value
            _t.option.series[0].data[3] = value
            setTimeout(this.echartsView.setOption(_t.option,true), 500)
            _t.tempValue = data.result
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
