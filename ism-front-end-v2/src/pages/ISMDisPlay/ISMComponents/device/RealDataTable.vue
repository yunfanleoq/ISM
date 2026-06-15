<template>

  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"  x="0px" y="0px"  xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <foreignObject :style="styleVar" pointer-events="all" :width="detail.style.position.w" :height="detail.style.position.h">
         <dv-scroll-board v-if="editMode||detail.style.visible" :ref="detail.identifier"  :config="scrollConfig" :style="{
          width:detail.style.position.w+'px',height:detail.style.position.h+'px'}" />
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
import svgView from '../View';
import {GetCurrentAlarmList} from "@/services/alarm";
import {formatDate} from "@/utils/common";
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
  name: 'view-device-real-data-table',
  inject: ['getNode'],
  i18n: require('../../../../i18n/language'),
  data() {
    return {
      detail:null,
      IsToolBox:false,
      editMode:true,
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
      imageURL:"",
      AlarmTimer:null,
      IsShowIndex:1,
      RealDataTableHeaderHeight:50,
      updateData:[],
      updateTimer:null,
      scrollConfig:{
        header: [],
        data: [],
        index: true,
        columnWidth: [50],
        indexHeader:"",
        align: ['center'],
        carousel: 'single',
        waitTime:1000
      },
      base:{
        text: "configComponent.RealDataTable.title",
        "icon": "icon-shishishuju_huaban",
        "isFontIcon": true,
        "info": {
          "type": "image",
          "action": [],
          "active": [],
          "dataBind":[],
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
              "w": 400,
              "h": 300
            },
            "visible":1,
            "backColor": "transparent",
            foreColor:"#ffffff",
            fontFamily:"宋体",
            fontSize:12,
            fontWeight:400,
            "zIndex": -1,
            "transform": 0,
            "diy":[
              {
                "name":"configComponent.RealDataTable.RealDataTableHeaderHeight",
                "type":1,
                "value":50,
                "min":1,
                "max":100,
                "key":"RealDataTableHeaderHeight",
              },
              {
                "name":"configComponent.AlarmList.headerBGC",
                "type": 2,
                "value": "#00BAFF",
                "key": "headerBGC",
              },
              {
                "name":"configComponent.AlarmList.oddRowBGC",
                "type": 2,
                "value": "#0a2732",
                "key": "oddRowBGC",
              },
              {
                "name":"configComponent.AlarmList.evenRowBGC",
                "type": 2,
                "value": "#003b51",
                "key": "evenRowBGC",
              },
              {
                "name":"configComponent.AlarmList.waitTime",
                "type":7,
                "value":1000,
                "min":100,
                "max":10000,
                "key":"waitTime",
              },
              {
                "name":"configComponent.AlarmList.rowNum",
                "type":1,
                "value":5,
                "min":5,
                "max":100,
                "key":"rowNum",
              },
              {
                name:"configComponent.AlarmList.index",
                type:6,
                value:0,
                enumList:[
                  {
                    value:0,
                    option:"configComponent.AlarmList.indexShow"
                  },
                  {
                    value:1,
                    option:"configComponent.AlarmList.indexHide"
                  }
                ],
                min:1,
                key:"IsShowIndex",
              }
            ]
          }
        }
      }
    }
  },
  computed: {
    styleVar() {
      return {
        "--foreColor": this.detail.style.foreColor ,
        '--fontFamily': this.detail.style.fontFamily,
        '--fontWeight': this.detail.style.fontWeight,
        "--fontSize": this.detail.style.fontSize+'px' ,
      };
    }
  },
  watch: {
    detail: {
      handler(newVal, oldVal) {
        this.detail.style.visible = newVal.style.visible
        if(this.editMode) {
          this.initComponents(newVal);
        }
      },
      deep: true
    }
  },
  methods: {
    updateTableValue(PushData){
      this.updateData =  this.scrollConfig.data
      for(let i=0;i<this.detail.active.length;i++)
      {
        if(this.detail.active[i].id==PushData.ID)
        {
          this.updateData[i][2]='<span style="color:'+this.detail.style.foreColor+';font-family:'+this.detail.style.fontFamily+';font-size:'+this.detail.style.fontSize+'px">'+PushData.result+(this.detail.active[i].condition.dataUnit!='undefined'?this.detail.active[i].condition.dataUnit:"")+'</span>'
          break;
        }
      }
    },
    updateDataTimer(){
      let refObj = this.detail.identifier
      let view = this.$refs[refObj]
      this.scrollConfig.data = this.updateData
      view.updateRows(this.updateData)
    },
    initComponents(option){
      if(this.IsToolBox)
      {
        return
      }
      let i=0
      for( i=0;i<option.style.diy.length;i++)
      {
        if(option.style.diy[i].key=="RealDataTableHeaderHeight")
        {
          this.RealDataTableHeaderHeight=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="headerBGC")
        {
          this.headerBGC=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="oddRowBGC")
        {
          this.oddRowBGC=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="evenRowBGC")
        {
          this.evenRowBGC=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="waitTime")
        {
          this.waitTime=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="rowNum")
        {
          this.rowNum=parseInt( option.style.diy[i].value)
        }
        else if(option.style.diy[i].key=="IsShowIndex")
        {
          this.IsShowIndex=option.style.diy[i].value
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
      i=0
      let data = []
      for (i = 0; i < option.active.length; i++) {
        let single = []
        single[0] = '<span style="color:'+option.style.foreColor+';font-family:'+option.style.fontFamily+';font-size:'+option.style.fontSize+'px">'+option.active[i].condition.DeviceName+'</span>'
        single[1] = '<span style="color:'+option.style.foreColor+';font-family:'+option.style.fontFamily+';font-size:'+option.style.fontSize+'px">'+option.active[i].condition.dataName+'</span>'
        single[2] = '<span style="color:'+option.style.foreColor+';font-family:'+option.style.fontFamily+';font-size:'+option.style.fontSize+'px">'+"-"+'</span>'
        data.push(single)
      }
     this.updateTable(data)
      if((!this.editMode)&&(!this.IsToolBox)) {
        clearInterval(this.updateTimer)
        this.updateTimer = setInterval(this.updateDataTimer, this.waitTime * option.active.length)
      }
    },
    updateTable(data){
      this.scrollConfig = {
        header:[
          this.$t('configComponent.RealDataTable.RealDataDeviceName'),
          this.$t('configComponent.RealDataTable.RealDataName'),
          this.$t('configComponent.RealDataTable.RealDataValue'),
        ],
        indexHeader:this.$t('configComponent.AlarmList.index'),
        columnWidth: this.detail.style.position.w/2,
        align: ['center'],
        carousel: 'single',
        headerHeight:this.RealDataTableHeaderHeight,
        rowNum:this.rowNum,
        headerBGC:this.headerBGC,
        oddRowBGC:this.oddRowBGC,
        evenRowBGC:this.evenRowBGC,
        waitTime:this.waitTime,
        data:data
      }
      if(this.IsShowIndex==0)
      {
        this.scrollConfig.index = true
      }
      else
      {
        this.scrollConfig.index = false
      }
      this.$forceUpdate()
    },
  },
  beforeDestroy(){
    clearInterval(this.updateTimer)
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
  },
  mounted() {
    let _t = this
    this.$nextTick(function(){
      this.initComponents(this.detail);
        let activeEvent = this.detail.identifier+"activeEvent"//动作数据
        let animateEvent = this.detail.identifier+"animateEvent"//动作数据

        _t.$EventBus.$on(activeEvent, _t._activeEventHandler = (data) => {
         _t.updateTableValue(data)
        })
        _t.$EventBus.$on(animateEvent, _t._animateEventHandler = (data) => {
          _t.isStart = data
        })

    });
  },
  created(){
    let _t = this
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
  }
}
</script>
<style scoped lang="less">
.svg-el {
  /*transform: rotate(45deg);*/
  transform-origin: center center;
}
::v-deep .dv-scroll-board .header .header-item {
  font-size: var(--fontSize);
  font-family: var(--fontFamily);
  font-weight:var(--fontWeight);
  color: var(--foreColor);
}
</style>