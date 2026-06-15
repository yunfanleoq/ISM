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
  name: 'ism-view-chart-category-1',
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
        "text": "configComponent.category1.title",
        "icon": "data:image/webp;base64,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",
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
              "w": 342,
              "h": 300
            },
            "backColor": "transparent",
            "foreColor": "#000000",
            fontFamily: "Arial",
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
                "name":"configComponent.ChartPublic.ChartAxisTickColor",
                "type":2,
                "value":"#5B7AD8",
                "key":"ChartAxisTickColor",
              },
              {
                "name":"configComponent.ChartPublic.EchartsWidth",
                "type":1,
                "value":20,
                "key":"EchartsWidth",
              },
              {
                "name":"configComponent.ChartPublic.IsShowLabel",
                "type":6,
                "value":0,
                "enumList":[
                  { "value":1, "option":"Yes" },
                  { "value":0, "option":"No" }
                ],
                "key":"IsShowLabel",
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
        title: {
          text: 'Waterfall Chart',
          textStyle:{

          }
        },

        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'shadow'
          },
          formatter: function (params) {
            let tar = params[0];
            return tar.name +' : ' + tar.value;
          }
        },
        grid: {
          left: '4%',
          right: '4%',
          top:'4%',
          bottom: '3%',
          containLabel: true
        },
        xAxis: {
          type: 'category',
          axisLabel:{
            textStyle:{},
          },
          data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
        },
        yAxis: {
          axisLabel:{
            textStyle:{},
          },
          type: 'value'
        },
        series: [
          {
            data: [120, 200, 150, 80, 70, 110, 130],
            type: 'bar',
            barWidth:20,
            textStyle:{

            },
            itemStyle: {
              color: '#a90000'
            },
            label: {
              show: true, // 必选：显示数值
              position: 'top', // 数值位置：top（柱子顶部）/inside（柱子内部）/bottom（底部）
              align: 'center', // 对齐方式
              verticalAlign: 'middle', // 垂直对齐
              fontSize: 12, // 字体大小
              color: '#000', // 字体颜色
              // 可选：自定义数值格式（比如保留2位小数、加单位）
              formatter: function(params) {
                // params.value 是当前柱子的数值
                return params.value; // 基础显示：仅数值
                // 进阶：保留2位小数 + 单位 → return params.value.toFixed(2) + ' 单位';
              }
            }
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
      this.option.title.textStyle.color = option.style.foreColor
      this.option.title.textStyle.fontFamily  = option.style.fontFamily

      this.option.xAxis.axisLabel.textStyle.color = option.style.foreColor
      this.option.xAxis.axisLabel.textStyle.fontSize = option.style.fontSize
      this.option.xAxis.axisLabel.textStyle.fontFamily  = option.style.fontFamily

      this.option.yAxis.axisLabel.textStyle.color = option.style.foreColor
      this.option.yAxis.axisLabel.textStyle.fontSize = option.style.fontSize
      this.option.yAxis.axisLabel.textStyle.fontFamily  = option.style.fontFamily

      this.option.series[0].label.color = option.style.foreColor
      this.option.series[0].label.fontSize = option.style.fontSize
      this.option.series[0].label.fontFamily  = option.style.fontFamily

      for( i=0;i<option.style.diy.length;i++)
      {
        if(option.style.diy[i].key=="ChartTitle")
        {
          this.option.title.text=option.style.diy[i].value
          if(this.option.title.text!="")
          {
            this.option.grid.top='15%'
          }
          else
          {
            this.option.grid.top='4%'
          }
        }
        else if(option.style.diy[i].key=="ChartUnit")
        {
          // this.option.series[0].detail.formatter='{value}'+option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="IsShowLabel")
        {
          this.option.series[0].label.show = option.style.diy[i].value==1?true:false
        }
        else if(option.style.diy[i].key=="ChartTitleFontSize")
        {
          this.option.title.textStyle.fontSize = option.style.diy[i].value
          this.option.xAxis.axisLabel.textStyle.fontSize = option.style.diy[i].value
          this.option.yAxis.axisLabel.textStyle.fontSize = option.style.diy[i].value
          this.option.series[0].label.fontSize = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="ChartTitleFontColor")
        {
          this.option.title.textStyle.color = option.style.diy[i].value
          this.option.series[0].label.color = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="ChartAxisTickColor")
        {
          this.option.series[0].itemStyle.color = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="EchartsWidth")
        {
          this.option.series[0].barWidth = parseInt(option.style.diy[i].value)
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
      if(!this.editMode)
      {
        this.option.xAxis.data=[]
        this.option.series[0].data=[]
        for(let i =0;i<this.detail.active.length;i++)
        {
          if(this.detail.active[i].condition.dataName=="")
          {
            continue
          }
          this.option.xAxis.data.push(this.detail.active[i].condition.dataName)
          let series= {
            value:100,
            name: this.detail.active[i].condition.DeviceName+"-"+this.detail.active[i].condition.dataName,
            dataID:this.detail.active[i].condition.deviceSN+this.detail.active[i].condition.dataID,
          }
          this.option.series[0].data.push(series)

        }
      }
      this.echartsView.setOption(this.option,true);
    },
    UpdateChartData : function(data) {
      for(let i=0;i<this.option.series[0].data.length;i++)
      {
        if(this.option.series[0].data[i].dataID == (data.DeviceSN+data.dataID))
        {
          this.option.series[0].data[i].value = data.result
        }
      }
      // 重新将数组赋值给echarts选项
      let _t = this
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
      if(current && current.detail) {
        _t.detail = current.detail
        // 当数据变化时重新初始化
        if (_t.$refs && !_t.IsToolBox) {
          _t.$nextTick(() => {
            _t.initComponents(_t.detail)
          })
        }
      }
    })
    node.on('change:size', ({ current }) => {
      if (_t.detail && _t.detail.style && _t.detail.style.position) {
        _t.detail.style.position.w = current.width
        _t.detail.style.position.h = current.height
      }
    });
    const nodeData = node.getData()
    this.detail = nodeData.detail || {
      type: 'chart-gauge',
      active: [],
      animate: { selected: [], condition: {}, isExpression: false },
      style: {
        position: { x: 0, y: 0, w: 200, h: 200 },
        foreColor: '#333333',
        fontFamily: 'Arial',
        fontSize: 14,
        borderWidth: 0,
        diy: []
      }
    }
    this.editMode = nodeData.editMode !== undefined ? nodeData.editMode : true
    this.showDeviceUuid = nodeData.showDeviceUuid || ''
    this.IsToolBox = nodeData.IsToolBox !== undefined ? nodeData.IsToolBox : false
    _t.$EventBus.$on('cell-editMode', (data) => {
      _t.editMode = data.edit
      _t.IsToolBox = data.toolbox
    })
  },
  mounted() {
    let _t = this
    this.$nextTick(function(){
      // 确保 detail 有 identifier（3D场景需要）
      if (!_t.detail.identifier) {
        _t.detail.identifier = 'chart_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9)
      }
      _t.initComponents(_t.detail);
      let activeEvent = _t.detail.identifier + 'activeEvent'
      let animateEvent = _t.detail.identifier + 'animateEvent'

      _t.$EventBus.$on(activeEvent, (data) => {
        _t.UpdateChartData(data)
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