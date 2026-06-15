<template>
  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"  x="0px" y="0px"  xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <foreignObject style="overflow:visible;" pointer-events="all" :width="detail.style.position.w" :height="detail.style.position.h">
        <div class="view-chart-gauge" :ref="detail.identifier" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
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
import BaseView from '../View';
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
  name: 'ism-view-alarm-status',
  inject: ['getNode'],
  i18n: require('../../../../i18n/language'),
  props: {

  },
  watch: {
    detail: {
      handler(newVal, oldVal) {
        this.initComponents(newVal)
      },
      deep: true
    }
  },
  data() {
    return {
        detail:null,
        IsToolBox:false,
        editMode:true,
        base:{
          "text": "configComponent.AlarmStatistics.title",
          "icon": "icon-anquangaojingtongji",
          "isFontIcon": true,
          "info": {
            "type": "chart-gauge",
            "action": [],
            "active": [
              {
                id:"AlarmCount",
                name:"SystemData.AlarmCount",
                result:"",
                isExpression:false,
                condition:{
                  deviceSN:"",
                  selectVideoType:0,
                  isBandDevice:false,
                  bandType:1,
                  dataID: "ism.SystemData.AlarmCount",
                  dataName: "告警总数",
                  operator:"",
                  OperatorValue:"",
                  OperatorMaxValue:"",
                },
              },
              {
                id:"TipsAlarmCount",
                name:"SystemData.TipsAlarmCount",
                result:"",
                isExpression:false,
                condition:{
                  deviceSN:"",
                  selectVideoType:0,
                  isBandDevice:false,
                  bandType:1,
                  dataID: "ism.SystemData.TipsAlarmCount",
                  dataName: "提示告警总数",
                  operator:"",
                  OperatorValue:"",
                  OperatorMaxValue:"",
                },
              },
              {
                id:"MinorAlarmCount",
                name:"SystemData.MinorAlarmCount",
                result:"",
                isExpression:false,
                condition:{
                  deviceSN:"",
                  selectVideoType:0,
                  isBandDevice:false,
                  bandType:1,
                  dataID: "ism.SystemData.MinorAlarmCount",
                  dataName: "次要告警总数",
                  operator:"",
                  OperatorValue:"",
                  OperatorMaxValue:"",
                },
              },
              {
                id:"ImportanceAlarmCount",
                name:"SystemData.ImportanceAlarmCount",
                result:"",
                isExpression:false,
                condition:{
                  deviceSN:"",
                  selectVideoType:0,
                  isBandDevice:false,
                  bandType:1,
                  dataID: "ism.SystemData.ImportanceAlarmCount",
                  dataName: "重要告警总数",
                  operator:"",
                  OperatorValue:"",
                  OperatorMaxValue:"",
                },
              },
              {
                id:"UrgencyAlarmCount",
                name:"SystemData.UrgencyAlarmCount",
                result:"",
                isExpression:false,
                condition:{
                  deviceSN:"",
                  selectVideoType:0,
                  isBandDevice:false,
                  bandType:1,
                  dataID: "ism.SystemData.UrgencyAlarmCount",
                  dataName: "紧急告警总数",
                  operator:"",
                  OperatorValue:"",
                  OperatorMaxValue:"",
                }
              },
              {
                id:"DeadlyAlarmCount",
                name:"SystemData.DeadlyAlarmCount",
                result:"",
                isExpression:false,
                condition:{
                  deviceSN:"",
                  selectVideoType:0,
                  isBandDevice:false,
                  bandType:1,
                  dataID: "ism.SystemData.DeadlyAlarmCount",
                  dataName: "致命告警总数",
                  operator:"",
                  OperatorValue:"",
                  OperatorMaxValue:"",
                },
              },
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
                "w": 300,
                "h": 400,
              },
              "backColor": "#000000",
              foreColor:"#00C0FF",
              fontSize:12,
              fontFamily:"数字字体-2",
              "zIndex": -3,
              "diy":[
                {
                  "name":"configComponent.AlarmStatistics.TipsTextColor",
                  "type":2,
                  "value":"#a5f1f1",
                  "key":"TipsTextColor",
                },
                {
                  "name":"configComponent.AlarmStatistics.MinorTextColor",
                  "type":2,
                  "value":'#378ffb',
                  "key":"MinorTextColor",
                },
                {
                  "name":"configComponent.AlarmStatistics.ImportanceColor",
                  "type":2,
                  "value":'#FFFF00',
                  "key":"ImportanceColor",
                },
                {
                  "name":"configComponent.AlarmStatistics.UrgencyColor",
                  "type":2,
                  "value":'#FFA500',
                  "key":"UrgencyColor",
                },
                {
                  "name":"configComponent.AlarmStatistics.DeadlyColor",
                  "type":2,
                  "value":'#FF0000',
                  "key":"DeadlyColor",
                },
                {
                  "name":"configComponent.AlarmStatistics.ISShowLegend",
                  "type":6,
                  "enumList":[
                    {option:'True',value:1},{option:'False',value:0}
                  ],
                  "value":0,
                  "key":"ISShowLegend",
                },
                {
                  "name":"configComponent.ChartPublic.EchartsOutside",
                  "type":1,
                  "value":40,
                  "key":"EchartsOutside",
                },
                {
                  "name":"configComponent.ChartPublic.EchartsInside",
                  "type":1,
                  "value":60,
                  "key":"EchartsInside",
                }
              ]
            }
          }
        },
        strokeColor:"#000000",
        fill:"#A1BFE2",
        strokeWidth:0.3,
        EchartsOutside:65,
        EchartsInside:65,
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
        TipsTextColor:"#57b7f7",
        MinorTextColor:"#0583d7",
        ImportanceColor:"yellow",
        UrgencyColor:"orange",
        DeadlyColor:"red",
        AlarmCount:100,
        TipsCount:40,
        MinorCount:30,
        ImportanceCount:10,
        UrgencyCount:10,
        ISShowLegend:false,
        DeadlyCount:10,
        trafficWay:[],
        trafficWayColor:[]
      }
  },
  methods: {
    optionFunc(){
      let _t = this
      this.trafficWay = [{
        name: this.$t('dataModel.alarm.Tips'),
        value: _t.TipsCount*1000
      },{
        name: this.$t('dataModel.alarm.Minor'),
        value: _t.MinorCount*1000
      },{
        name: this.$t('dataModel.alarm.Importance'),
        value: _t.ImportanceCount*1000
      },{
        name: this.$t('dataModel.alarm.Urgency'),
        value: _t.UrgencyCount*1000
      },{
        name: this.$t('dataModel.alarm.Deadly'),
        value: _t.DeadlyCount*1000
      }];

      let data = [];
      this.trafficWayColor=[_t.TipsTextColor,_t.MinorTextColor,_t.ImportanceColor,_t.UrgencyColor,_t.DeadlyColor,]
      for (let i = 0; i < this.trafficWay.length; i++) {
        data.push(
            {
              value: _t.trafficWay[i].value,
              name: _t.trafficWay[i].name,
              itemStyle: {
                normal: {
                  borderWidth: 5,
                  shadowBlur: 20,
                  borderColor:_t.trafficWayColor[i],
                  shadowColor: _t.trafficWayColor[i]
                }
              }
            },
            {
              value: 2,
              name: '',
              itemStyle: {
                normal: {
                  label: {
                    show: false
                  },
                  labelLine: {
                    show: false
                  },
                  color: 'rgba(0, 0, 0, 0)',
                  borderColor: 'rgba(0, 0, 0, 0)',
                  borderWidth: 0
                }
              }
            }
        );
      }
      let option = {
        backgroundColor: _t.detail.style.backColor,
        color : _t.trafficWayColor,
        tooltip: {
          show: false
        },
        legend: {
          show:_t.ISShowLegend,
          icon: "circle",
          orient: 'horizontal',
          data:[_t.$t('dataModel.alarm.Tips'),_t.$t('dataModel.alarm.Minor'),_t.$t('dataModel.alarm.Importance'),_t.$t('dataModel.alarm.Urgency'),_t.$t('dataModel.alarm.Deadly')],
          bottom: '2%',
          align: 'right',
          textStyle: {
            fontSize: _t.detail.style.fontSize,
            color:  _t.detail.style.foreColor,
            fontFamily: _t.detail.style.fontFamily,
          },
          itemGap: 20
        },
        toolbox: {
          show: false
        },
        series: [
          {
            name: '',
            type: 'pie',
            roundCap: true,
            radius: [_t.EchartsInside+'%', _t.EchartsOutside+'%'],
            // center : ['70%', '50%'],    // 默认全局居中
            hoverAnimation: true,
            itemStyle: {
              normal: {
                label: {
                  show: true,
                  position: 'outside',
                  color:  _t.detail.style.foreColor,
                  fontSize: _t.detail.style.fontSize,
                  fontFamily: _t.detail.style.fontFamily,
                  formatter: function(params) {
                    let percent = 0;
                    let total = _t.AlarmCount*1000;
                    percent = ((params.value / total) * 100).toFixed(0);
                    if(params.name !== '') {
                      return params.name + '\n' + '\n' + _t.$t('SystemData.Count')+'：' + params.value/1000 + '\n' + '\n' + _t.$t('configComponent.DeviceStatus.Percent')+'：' + percent + '%';
                    }else {
                      return '';
                    }
                  },
                },
                labelLine: {
                  length:20,
                  length2:10,
                  show: true,
                  color: _t.detail.style.foreColor,
                }
              }
            },
            data: data
          }
        ]
      }
      return option
    },
    initComponents(option){
      if(this.IsToolBox)
      {
        return
      }
      let refObj = this.detail.identifier
      let view = this.$refs[refObj]
      let i=0
      for( i=0;i<option.style.diy.length;i++)
      {
        if(option.style.diy[i].key=="TipsTextColor")
        {
          this.TipsTextColor = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="MinorTextColor")
        {
          this.MinorTextColor = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="ImportanceColor")
        {
          this.ImportanceColor = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="UrgencyColor")
        {
          this.UrgencyColor = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="DeadlyColor")
        {
          this.DeadlyColor = option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="ISShowLegend")
        {
           if(option.style.diy[i].value==1)
          {
            this.ISShowLegend = true
          }
           else
           {
             this.ISShowLegend =false
           }

        }
        else if(option.style.diy[i].key=="EchartsInside")
        {
          this.EchartsInside = parseInt(option.style.diy[i].value)
        }
        else if(option.style.diy[i].key=="EchartsOutside")
        {
          this.EchartsOutside = parseInt(option.style.diy[i].value)
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


      this.$nextTick(function (){
        if (!this.echartsView) {
          this.echartsView = echarts.init(view, null);
        }
        else
        {
          this.echartsView.resize()
        }
        let echartOptions = this.optionFunc()
        this.echartsView.setOption(echartOptions,true)
        this.$forceUpdate()
      })
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
    this.$nextTick(function(){
        let activeEvent = this.detail.identifier+"activeEvent"//动作数据
        let animateEvent = this.detail.identifier+"animateEvent"//动作数据

        _t.$EventBus.$on(activeEvent, _t._activeEventHandler = (data) => {
          let charge = 0
          console.log(data)
          if(data.ID == "AlarmCount")
          {
            if(_t.AlarmCount!=data.result)
            {
              charge = 1
              _t.AlarmCount = data.result
            }
          }
          else if(data.ID == "TipsAlarmCount")
          {
            if(_t.TipsCount!=data.result) {
              charge = 1
              _t.TipsCount = data.result
            }
          }
          else if(data.ID == "MinorAlarmCount")
          {
            if(_t.MinorCount!=data.result) {
              charge = 1
              _t.MinorCount = data.result
            }
          }
          else if(data.ID == "ImportanceAlarmCount")
          {
            if(_t.ImportanceCount!=data.result) {
              charge = 1
              _t.ImportanceCount = data.result
            }
          }
          else if(data.ID == "UrgencyAlarmCount")
          {
            if(_t.UrgencyCount!=data.result) {
              charge = 1
              _t.UrgencyCount = data.result
            }
          }
          else if(data.ID == "DeadlyAlarmCount")
          {
            if(_t.DeadlyCount!=data.result) {
              charge = 1
              _t.DeadlyCount = data.result
            }
          }
          if(charge)
          {
            let echartOptions = this.optionFunc()
            setTimeout(this.echartsView.setOption(echartOptions,true), 500)
            this.$forceUpdate()
          }
        })
        _t.$EventBus.$on(animateEvent, _t._animateEventHandler = (data) => {
          _t.isStart = data
        })
    });
    this.GetNodeObj = this.getNode()
    this.GetNodeObj.on('change:data', ({ current }) => {
      if(current) {
        _t.detail = current.detail
      }
    })
    this.GetNodeObj.on('change:size', ({ current }) => {
      _t.detail.style.position.w = current.width
      _t.detail.style.position.h = current.height
    });
    this.detail = this.GetNodeObj.getData().detail
    this.editMode = this.GetNodeObj.getData().editMode
    this.showDeviceUuid = this.GetNodeObj.getData().showDeviceUuid
    this.IsToolBox = this.GetNodeObj.getData().IsToolBox
    _t.$EventBus.$on('cell-editMode', _t._cellEditModeHandler = (data) => {
      _t.editMode = data.edit
      _t.IsToolBox = data.toolbox
	  _t.initComponents(_t.detail)
    })
  },
  mounted() {
    this.$nextTick(function(){
        this.initComponents(this.detail);
      });
    },
  beforeDestroy() {
    // 清理 ECharts 实例
    if (this.echartsView) {
      this.echartsView.dispose()
      this.echartsView = null
    }
    // 清理 EventBus 监听
    if (this._activeEventHandler) {
      this.$EventBus.$off(this.detail.identifier + 'activeEvent', this._activeEventHandler)
    }
    if (this._animateEventHandler) {
      this.$EventBus.$off(this.detail.identifier + 'animateEvent', this._animateEventHandler)
    }
    if (this._cellEditModeHandler) {
      this.$EventBus.$off('cell-editMode', this._cellEditModeHandler)
    }
  }
}
</script>

<style lang="less">
.view-chart-gauge {
    height: 100%;
    width: 100%;
}
</style>
