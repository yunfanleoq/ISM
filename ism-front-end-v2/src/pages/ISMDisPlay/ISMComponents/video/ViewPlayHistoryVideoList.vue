<template>
  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"  x="0px" y="0px"  xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <foreignObject pointer-events="all" :width="detail.style.position.w" :height="detail.style.position.h">
        <div>
    <div class="h5videodiv" :style="{'width': width+'px','height': height+'px',}" >
      <LivePlayer :id="detail.identifier" :video-url="video_url"  :video-title="VideoTitle" :alt="VideoTitle+$t('configComponent.video.VideoOfflineTips')" :loop="true" :muted="false" aspect="fullscreen" fluent autoplay stretch></LivePlayer>
      <div  class="video-close" style="display: none;">关闭</div>
      <div  class="video-close channel-selector" @click="showSystemVideoModel()" >选择回放地址</div>
    </div>
          <system-history-video-model @onSelectVideo="onSelectVideo" :networkVideo="videoComponentData" ref="systemVideoModel"></system-history-video-model>
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
import {getVideoStatus } from "@/services/video";
import LivePlayer from '@liveqing/liveplayer'
import systemHistoryVideoModel from '@/components/systemHistoryVideoModel/systemHistoryVideoModel'
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
export default {
  mixins: [ISMChildAutoMixin],
  name: 'ism-view-play-history-video-list',
  inject: ['getNode'],
  i18n: require('../../../../i18n/language'),
  props: {

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
  components: {
    LivePlayer,
    systemHistoryVideoModel
  },
  data() {
    return {
      detail:null,
      IsToolBox:false,
      editMode:true,
      isStart:false,
      animateType:"blink",
      spinDirection:0,
      animateSpeed:0.5,
      animateSpinSpeed:0.5,
      blinkSpeed:0.5,
      startColor:"#74f808",
      stopColor:"#74f808",
      strokeColor:"#000000",
      fill:"#A1BFE2",
      strokeWidth:0.3,
      fillOpacity:1,
      strokeOpacity:1,
      videoComponentData:{},
      token:"",
      videoStyle:{
        height:300+"px",
        width:'auto'
      },
      VideoStatus:0,
      VideoTitle:"",
      videoType:1,
      video_url:"",
      temp_video_url:"",
      width:600,
      height:600,
      base:{
        "text": "configComponent.video.historyVideoList",
        "icon": "icon-shipinhuifang",
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
            "zIndex": -1,
            "transform": 0,
            "diy":[
              {
                "name":"configComponent.video.Token",
                "type":8,
                "value":{
                  "type":1,
                  "key":""
                },
                "key":"token",
              },
            ]
          }
        }
      }
    }
  },
  methods: {
    showSystemVideoModel(){
      this.$refs.systemVideoModel.showModal()
    },
    onSelectVideo(token){
      this.video_url = token.url
      this.VideoTitle = token.name
      this.visible =false
    },
    getSystemVideoList(number){
      let _t = this
      _t.videoPlayerList=[]
      let tableData={}

      for(let i=0;i<number;i++)
      {
        tableData.url = ""
        tableData.Name = ""
        tableData.status = 0
        _t.videoPlayerList.push(tableData)
        tableData={}
      }
    },
    initWebSock(){
      let _t=this

        _t.$EventBus.$on("RealAlarm", (data) => {
            if((data.DeviceUuid==_t.token)&&(data.AlarmID=="videoConnectStatusAlarm"))
            {
              _t.VideoStatus = data.Value
              if(_t.VideoStatus==0)
              {
                _t.video_url = _t.temp_video_url
              }
              else{
                _t.video_url=""
              }
            }
        });

    },
    getVideoStatus() {
      let _t = this
      let params = {
        uuid:this.token
      }
      getVideoStatus(params).then(function(res) {
          let data = res.data
        if(data.code!=0)
        {
          return
        }
        if(data.result.status==true)
        {
          _t.VideoStatus=0
          _t.video_url=_t.temp_video_url
        }
        else
        {
          _t.VideoStatus=1
          _t.video_url=""
        }
      });
    },
    initComponents(option){
      this.width = option.style.position.w
      this.height = option.style.position.h
      let i=0
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
    }
  },
  mounted() {
    let _t = this
    this.$nextTick(function(){
      this.initComponents(this.detail);
      if(true){
        let activeEvent = this.detail.identifier+"activeEvent"//动作数据
        let animateEvent = this.detail.identifier+"animateEvent"//动作数据

        _t.$EventBus.$on(activeEvent, (data) => {

        })
        _t.$EventBus.$on(animateEvent, (data) => {
          _t.isStart = data
        })
      }
    });
  },
  destroyed(){
    this.video_url=""
  },
  computed: {
    styleVar() {
      return {
        "height": this.detail.style.position.h+'px',
        "--foreColor": this.foreColor ,
        '--backColor':this.backColor,
        "--selectedColor": this.selectedColor ,
        '--hoverColor': this.hoverColor,
        '--selectedTextColor': this.selectedTextColor,
        '--TextFontSize': this.TextFontSize+'px',
        '--hoverTextColor': this.hoverTextColor,
        '--SearchColor': this.SearchColor,
        '--SearchBackColor': this.SearchBackColor,
        '--SearchBorderColor': this.SearchBorderColor,

        '--dateSelectColor': this.dateSelectColor,
        '--dateSelectBackColor': this.dateSelectBackColor,
        '--dateSelectBorderColor': this.dateSelectBorderColor,

        '--tableHeaderColor': this.tableHeaderColor,
        '--tableHeaderBackColor': this.tableHeaderBackColor,
        '--tableSplitColor': this.tableSplitColor,
        '--tableHoverColor':this.tableHoverColor
      };
    },
    animatedStyle(){
      return {
        "--blinkSpeed":this.blinkSpeed+'s',
        "--stopColor":this.stopColor,
        "--startColor":this.startColor,
        "--animateSpeed":this.animateSpeed+'s',
        "--animateSpinSpeed":this.animateSpinSpeed+'s'
      }
    },
    textAlign: function(){
      if(this.detail.style.textAlign == undefined) {
        return "center";
      } else {
        return this.detail.style.textAlign;
      }
    },
    lineHeight: function() {
      if(this.detail.style.lineHeight == undefined) {
        return this.detail.style.position.h;
      }
      return this.detail.style.lineHeight;
    },
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
    _t.$EventBus.$on('cell-editMode', (data) => {
      _t.editMode = data.edit
      _t.IsToolBox = data.toolbox
      _t.initComponents(_t.detail)
    })
    this.initComponents(this.detail);
  }

}
</script>

<style lang="less" scoped>

.video-close, .video-title {
  position: absolute;
  top: 5px;
  color: #fff;
  background-color: hsla(0,0%,50%,.5);
  padding: 2px 5px;
  border-radius: 2px;
  max-width: 120px;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}
.view-list .video-show  {
  border-left-color: transparent;
}
.video-close {
  font-size: 12px;
  cursor: pointer;
}
.h5video {
  position: relative;
  border: 1px solid black;
  background-color:#000000;
}

.h5video1 {
  width: 100%;
  height: 100%;
  border: 1px solid black;
  background-color:#000000;
}

.h5videodiv{
  border: 1px solid black;
  background-color:#000000;
  position:relative;
}
.s-full {
  position: absolute;
  left:35%;
  top:45%;
  color: white;
  z-index: 999;
  font-size: 20px;
  cursor: pointer;
}
.h5videodiv1{
  width: 100%;
  height: 100%;
  position:relative;
  background-color:#000000;
}

</style>
