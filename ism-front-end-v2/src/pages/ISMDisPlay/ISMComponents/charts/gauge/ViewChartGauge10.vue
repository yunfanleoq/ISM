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
  name: 'ism-view-chart-gauge-10',
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
          "text": "configComponent.AnimationRing.title",
          "icon": "icon-huanxingtu",
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
                "w": 269,
                "h": 270
              },
              "backColor": "#000000",
              foreColor:"#00C0FF",
              fontSize:30,
              "zIndex": 1,
              "transform": 0,
              "diy":[
                {
                  "name":"configComponent.ChartPublic.ChartUnit",
                  "type":4,
                  "value":"kg/m³",
                  "key":"ChartUnit",
                },
                {
                  "name":"configComponent.AnimationRing.RingColor",
                  "type":2,
                  "value":"#0CD3DB",
                  "key":"RingColor",
                },
                {
                  "name":"configComponent.ChartPublic.ChartWidth",
                  "type":7,
                  "value":20,
                  "key":"ChartWidth",
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
        MaxValue:100,
        eventValue: '0.00',
        eventUnit: '',
        chartUnit:"",
        ChartWidth:10,
        angle:0,
        RingColor:"#0CD3DB",
        value:55.0,
      }
  },
  beforeDestroy () {
    if (this._drawFrameId) {
      cancelAnimationFrame(this._drawFrameId)
      this._drawFrameId = null
    }
    if (this.echartsView != null&&(typeof this.echartsView.dispose=="function")) {
      this.echartsView.dispose()
    }
  },
  methods: {
    getCirclePoint(x0, y0, r, angle) {
      let x1 = x0 + r * Math.cos(angle * Math.PI / 180)
      let y1 = y0 + r * Math.sin(angle * Math.PI / 180)
      return {
        x: x1,
        y: y1
      }
    },
    optionFunc(){
      let _t = this
      let option={
            backgroundColor:_t.detail.style.backColor,
            title: {
              text: '{a|'+ _t.value +'}{c|'+_t.chartUnit+'}',
              x: 'center',
              y: 'center',
              textStyle: {
                rich:{
                  a: {
                    fontSize: _t.detail.style.fontSize,
                    color: _t.detail.style.foreColor,
                  },

                  c: {
                    fontSize: _t.detail.style.fontSize-10,
                    color: _t.detail.style.foreColor,
                    padding: [5,0]
                  }
                }
              }
            },
            legend: {
              type: "plain",
              orient: "vertical",
              right: 0,
              top: "10%",
              align: "auto",
              data: [],
              textStyle: {
                color: "white",
                fontSize: 16,
                padding: [10, 1, 10, 0]
              },
              selectedMode:false
            },
            series: [
            {
              name: "ring1",
              type: 'custom',
              radius: ['65%', '95%'],
              coordinateSystem: "none",
              renderItem: function(params, api) {
                  return {
                  type: 'arc',
                  shape: {
                    cx: api.getWidth() / 2,
                    cy: api.getHeight() / 2,
                    r: Math.min(api.getWidth(), api.getHeight()) / 2 * 0.85,
                    startAngle: (0+_t.angle) * Math.PI / 180,
                    endAngle: (90+_t.angle) * Math.PI / 180
                  },
                  style: {
                    stroke: _t.RingColor,
                    fill: "transparent",
                    lineWidth: _t.ChartWidth
                  },
                  silent: true
                };
              },
              data: [0]
            },
            {
              name: "ring2",
              type: 'custom',
              coordinateSystem: "none",
              radius: "110%",
              center : ['50%', '58%'],    // 默认全局居中
              renderItem: function(params, api) {
              return {
            type: 'arc',
            shape: {
              cx: api.getWidth() / 2,
              cy: api.getHeight() / 2,
              r: Math.min(api.getWidth(), api.getHeight()) / 2 * 0.85,
              startAngle: (180+_t.angle) * Math.PI / 180,
              endAngle: (270+_t.angle) * Math.PI / 180
            },
            style: {
              stroke:  _t.RingColor,
              fill: "transparent",
              lineWidth: _t.ChartWidth
            },
            silent: true
          };
              },
              data: [0]
            },
            {
              name: "ring3",
              type: 'custom',
              radius: "110%",
              center : ['50%', '58%'],    // 默认全局居中
              coordinateSystem: "none",
              renderItem: function(params, api) {
                return {
                  type: 'arc',
                  shape: {
                    cx: api.getWidth() / 2,
                    cy: api.getHeight() / 2,
                    r: Math.min(api.getWidth(), api.getHeight()) / 2 * 1,
                    startAngle: (270+-_t.angle) * Math.PI / 180,
                    endAngle: (40+-_t.angle) * Math.PI / 180
                  },
                  style: {
                    stroke:  _t.RingColor,
                    fill: "transparent",
                    lineWidth: _t.ChartWidth
                  },
                  silent: true
                };
              },
              data: [0]
            },
        ]
    }
      return option
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
      for( i=0;i<option.style.diy.length;i++)
      {
        if(option.style.diy[i].key=="ChartUnit")
        {
          this.chartUnit = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="RingColor")
        {
          this.RingColor= option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="ChartWidth")
        {
          this.ChartWidth= option.style.diy[i].value
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

      // this.option.series[0].symbolSize[0] = this.detail.style.position.w
      // this.option.series[1].symbolSize[0] = this.detail.style.position.w
      // this.option.series[2].symbolSize[0] = this.detail.style.position.w
      // this.option.series[3].symbolSize[0] = this.detail.style.position.w
      // this.option.series[4].symbolSize[0] = this.detail.style.position.w
      // this.option.series[5].barWidth= this.detail.style.position.w
      // this.option.series[6].barWidth= this.detail.style.position.w
      // this.option.series[5].label.color = this.detail.style.foreColor
      // this.option.series[5].label.fontSize = this.detail.style.fontSize
      // this.option.series[5].label.distance = parseInt(this.detail.style.fontSize)+20
      // this.option.series[5].label.formatter = '{c}'+chartUnit
      this.echartsView.resize()
      this.echartsView.setOption(this.optionFunc(),true);

    },
    draw(){
      let _t = this
      _t.angle = _t.angle+3
      this.echartsView.setOption(this.optionFunc(), true)
      this._drawFrameId = requestAnimationFrame(this.draw);
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
    this._drawFrameId = requestAnimationFrame(this.draw);
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
            _t.value = data.result
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
