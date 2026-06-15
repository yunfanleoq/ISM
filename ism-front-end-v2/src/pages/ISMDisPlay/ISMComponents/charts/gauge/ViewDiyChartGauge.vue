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
import 'echarts-liquidfill'
import BaseView from '../../View';
import {ComponentRestApi} from "@/services/RestApi";
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
export default {
  mixins: [ISMChildAutoMixin],
  name: 'ism-view-chart-gauge-diy',
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
      ChartOption:"",
      detail:{},
      IsToolBox:false,
      editMode:true,
      base:{
        "text": "configComponent.chartGauge.diyManGauge",
        "icon": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAAAXNSR0IArs4c6QAAAmpJREFUWEftVktym0AQfS2qHF9BO7JTDuFEOomlAzhhq5VgpS0VHwByEhH5EtqFHVdwpQo61cMfwzAgubzJ7KSa6ff69etuCB986IPx8Z/AbAW+nI8nTrG7bPbxNWWcRUDAAawFmFN8vobEZAKr83FLQNDIOuYUm7kkJhFYnY42WfjTI/lsEsYECnCR3h6o+SwSxgSaddeYLppqTCMCq5fjgRiuidsJCLMUnqknRgmsXo5rYoj0xocJ7uVh75k8GCeQG09cr9pu7NxUAY7uAvDiF21eo6L9DhoTRkzwLg/7iE/3NiwOwBzRt79aJQYV4POnE1hlHQMcSiDphIWFAwPb1hwg7CrgRfYIUMMv7OpI9BJQGSy42+8xMtopNXJfBEwIy1rz77tDG7hBkWhHX1/DvvINKzAYkF1kqixqB2iBRb2CtOMktu8v3+yNNwScH0kARuz/XHq5El1JVR6x1BdEUqK+wRSDyJOsBRgpAhBsLLDz/WXUVKJFwHlKTqDC7YQQBE9Y8/l+C2adAcuYlV8UcIZHoDU/4i6JioDzlGxBrSWTZwqEI2oU93Kjyg/neyJkhwaXkNiU5VAECrZ9S6bOrHhUGLSxE2qXjwDXyjMi/3m5kT9qBXLJtMumpYaFdensIgHd2zoRhuc/L6uOaHugv27d7qnq2DDY2JSsStkN1tuGRoHFpNwaSH1tPghcXtbugqKmMvWGvgF0q8EV8xrsDv2VgXa6GthIgSaK4yRrZKpN+9VgRLDUoJn0lTy6jjskhoZL1ddjkhuZcCxIZVJLTcrWaB17exMCU0F09yeV4JbAk034HuAS8x9Vrjgw6gGMwgAAAABJRU5ErkJggg==",
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
            "zIndex": 1,
            "transform": 0,
            "diy":[
              {
                "name":"component.public.ChartOption",
                "type":11,
                "value":"",
                "rows":30,
                "key":"ChartOption",
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
      echartTimer: null,
      eventValue: '0.00',
      eventUnit: '',
      option:{
        backgroundColor: '#0F224C',
        series: [
          {
            type: 'liquidFill',
            shape: 'container',
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
            waveAnimation: false,
            amplitude:0,
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
        if(option.style.diy[i].key=="ChartOption")
        {
          this.ChartOption=option.style.diy[i].value
          if(this.ChartOption!="") {
            this.echartsView = echarts.init(view, null);
            let func = new Function(
                'echarts',
                'ComponentRestApi',
                'echartTimer',
                option.style.diy[i].value
            ).bind(this)
            func(echarts, ComponentRestApi, this.echartTimer)
          }
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
    clearInterval(this.echartTimer)
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
      }
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
        if(data.ID == "ShowData")
        {
          if( _t.option.series[0].data[0]!=data.result)
          {
            let value = parseInt(data.result)/this.progressMax
            _t.option.series[0].label.normal.formatter = data.result+this.chartUnit
            _t.option.series[0].data[0] = value
            _t.option.series[0].data[1] = value
            _t.option.series[0].data[2] = value
            _t.option.series[0].data[3] = value
            setTimeout(this.echartsView.setOption(_t.option,true), 500)
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
