<template>

  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"  x="0px" y="0px"  xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <foreignObject style="overflow:visible;" pointer-events="all" :width="detail.style.position.w" :height="detail.style.position.h">
        <div :style="boardContainerStyle">
         <div v-if="showEmptyState" :style="emptyStateStyle">{{ $t('configComponent.AlarmList.emptyText') }}</div>
         <dv-scroll-board :config="scrollConfig" :style="{width:detail.style.position.w+'px',height:detail.style.position.h+'px'}" />
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
import svgView from '../View';
import {GetCurrentAlarmList} from "@/services/alarm";
import {formatDate} from "@/utils/common";
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
  name: 'view-device-alarm-list',
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
      rowNum:5,
      scrollConfig:{
        header: ['列1', '列2', '列3'],
        data: [
          ['行1列1', '行1列2', '行1列3'],
          ['行2列1', '行2列2', '行2列3'],
          ['行3列1', '行3列2', '行3列3'],
          ['行4列1', '行4列2', '行4列3'],
          ['行5列1', '行5列2', '行5列3'],
          ['行6列1', '行6列2', '行6列3'],
          ['行7列1', '行7列2', '行7列3'],
          ['行8列1', '行8列2', '行8列3'],
          ['行9列1', '行9列2', '行9列3'],
          ['行10列1', '行10列2', '行10列3']
        ],
        index: true,
        columnWidth: [50],
        indexHeader:"",
        align: ['center'],
        carousel: 'single',
        waitTime:1000
      },
      base:{
        text: "configComponent.AlarmList.title",
        "icon": "icon-icon-test",
        "isFontIcon": true,
        "info": {
          "type": "image",
          "action": [],
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
              "h": 200
            },
            "visible":1,
            "backColor": "transparent",
            "zIndex": -1,
            "transform": 0,
            "diy":[
              {
                "name":"configComponent.AlarmList.rowNum",
                "type":1,
                "value":5,
                "min":5,
                "max":100,
                "key":"rowNum",
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
            ]
          }
        }
      }
    }
  },
  watch: {
    detail: {
      handler(newVal, oldVal) {
        if(this.editMode) {
          this.initComponents(newVal);
        }
      },
      deep: true
    }
  },
  computed: {
    showEmptyState() {
      return this.scrollConfig && Array.isArray(this.scrollConfig.data) && this.scrollConfig.data.length === 0
    },
    boardContainerStyle() {
      const hasData = this.scrollConfig && Array.isArray(this.scrollConfig.data) && this.scrollConfig.data.length > 0
      const emptyBg = this.detail && this.detail.style && this.detail.style.backColor
        ? this.detail.style.backColor
        : '#ffffff'
      return {
        position: 'relative',
        width: this.detail && this.detail.style ? this.detail.style.position.w + 'px' : '100%',
        height: this.detail && this.detail.style ? this.detail.style.position.h + 'px' : '100%',
        backgroundColor: hasData ? 'transparent' : emptyBg
      }
    },
    emptyStateStyle() {
      return {
        position: 'absolute',
        left: '0',
        top: '0',
        width: this.detail && this.detail.style ? this.detail.style.position.w + 'px' : '100%',
        height: this.detail && this.detail.style ? this.detail.style.position.h + 'px' : '100%',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        color: (this.detail && this.detail.style && this.detail.style.foreColor) ? this.detail.style.foreColor : '#666666',
        fontSize: '14px',
        zIndex: '1',
        pointerEvents: 'none'
      }
    }
  },
  methods: {
    initComponents(option){

      if(this.IsToolBox)
      {
        return
      }
      let i=0
      for( i=0;i<option.style.diy.length;i++)
      {
        if(option.style.diy[i].key=="rowNum")
        {
          this.rowNum=parseInt( option.style.diy[i].value)
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
      }
      this.updateTable([])
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
    updateTable(data){
      this.scrollConfig = {
        header:[
          this.$t('reporting.AlarmHistory.DeviceName'),
          this.$t('reporting.AlarmHistory.AlarmName'),
          this.$t('reporting.AlarmHistory.HappenTime'),
          this.$t('reporting.AlarmHistory.AlarmLevel'),
          this.$t('dataModel.editData.AlarmMessage')
        ],
        indexHeader:this.$t('configComponent.AlarmList.index'),
        index: true,
        columnWidth: [50],
        align: ['center'],
        carousel: 'single',
        rowNum:this.rowNum,
        headerBGC:this.headerBGC,
        oddRowBGC:this.oddRowBGC,
        evenRowBGC:this.evenRowBGC,
        waitTime:this.waitTime,
        data:data
      }

      clearInterval(this.AlarmTimer)
      let timer = this.waitTime * (data.length)
      if(timer<5000)
      {
        timer=5000
      }
      this.AlarmTimer = setInterval(this.QueryAlarmList, timer)
    },
    QueryAlarmList(){
      let _t = this


      _t.dataSource = []
      const params = {
        deviceList:this.SelectDevice,
        dataList:this.SelectAlarmData,
      }
      this.messageShowLoad=true
      _t.scrollConfig.data = []
      GetCurrentAlarmList(params).then(function (res){
        if(res.data.code==0)
        {
          let updateData = []
          for(let i = 0;i<res.data.list.length;i++)
          {
              let single = []
            single[0] = res.data.list[i].DeviceName
            single[1] = _t.$t(res.data.list[i].AlarmName)
            single[2] = formatDate(new Date(res.data.list[i].HappenTime),'yyyy-MM-dd hh:mm:ss')
            if(res.data.list[i].AlarmLevel==0)
            {
              single[3] ='<span style="color:#0099FF;">'+ _t.$t('dataModel.alarm.Tips')+'</span>'
            }
            else if(res.data.list[i].AlarmLevel==1)
            {
              single[3] = '<span style="color:#0099FF;">'+ _t.$t('dataModel.alarm.Minor')+'</span>'
            }
            else if(res.data.list[i].AlarmLevel==2)
            {
              single[3] = '<span style="color:yellow;">'+ _t.$t('dataModel.alarm.Importance')+'</span>'
            }
            else if(res.data.list[i].AlarmLevel==3)
            {
              single[3] = '<span style="color:orange;">'+ _t.$t('dataModel.alarm.Urgency')+'</span>'
            }
            else if(res.data.list[i].AlarmLevel==4)
            {
              single[3] = '<span style="color:red;">'+ _t.$t('dataModel.alarm.Deadly')+'</span>'
            }
            single[4] = '<span title='+_t.$t(res.data.list[i].AlarmMessage)+'>'+ _t.$t(res.data.list[i].AlarmMessage)+'</span>'
            updateData.push(single)
          }
          _t.updateTable (updateData)
        }
        _t.messageShowLoad=false
      }).catch(function(){
        _t.messageShowLoad=false
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
    }
  },
  beforeDestroy(){
    clearInterval(this.AlarmTimer)
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
    if(!_t.editMode)
    {
      _t.QueryAlarmList()
    }
    this.$nextTick(function(){
      this.initComponents(this.detail);
        let activeEvent = this.detail.identifier+"activeEvent"//动作数据
        let animateEvent = this.detail.identifier+"animateEvent"//动作数据

        _t.$EventBus.$on(activeEvent, _t._activeEventHandler = (data) => {


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
      _t.QueryAlarmList()
    })
  }
}
</script>
<style >
.svg-el {
  /*transform: rotate(45deg);*/
  transform-origin: center center;
}

</style>
