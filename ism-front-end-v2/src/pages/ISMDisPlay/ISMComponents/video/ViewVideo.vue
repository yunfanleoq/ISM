<template>
        <div :style="animatedStyle" v-show="detail.style.visible==1 ||isStart? true:false">
          <div :class="{
                'animated':true,[`${detail.style.animate}`]: true,
                'color-animation':isStart&&animateType.includes('millcolorGrad')&&!IsToolBox,
                'blink-animation':isStart&&animateType.includes('blink')&&!IsToolBox,
                'scale-animation':isStart&&animateType.includes('Zoom')&&!IsToolBox,
                'rotate-animation':isStart&&animateType.includes('animateSpin')&&!IsToolBox&&spinDirection==0,
                'rotate-anti-animation':isStart&&animateType.includes('animateSpin')&&!IsToolBox&&spinDirection==1
              }"
               :style="{
                                      width: detail.style.position.w + 'px',
                                      height: detail.style.position.h + 'px',
                                      'background-color': detail.style.backColor,
                                      'border-radius':detail.style.BorderEdges+'px',
                                      opacity:detail.style.opacity,
                                      borderWidth: detail.style.borderWidth + 'px',
                                      borderStyle: detail.style.borderStyle,
                                      borderColor: detail.style.borderColor,
                                      transform: detail.style.transform==-1099? 'rotateY(180deg)':detail.style.transform==-1098?'rotateX(180deg)':'',
                                  }">
    <div v-if="videoType==1" >
      <div class="h5videodiv" :style="{'width': width+'px','height': height+'px',}" >
        <LivePlayer :id="detail.identifier" :video-url="video_url"  :video-title="isShowTitle?VideoTitle:''" :alt="VideoTitle+$t('configComponent.video.VideoOfflineTips')" aspect="fullscreen" fluent autoplay live stretch>
          <div @click="isShowControl=!isShowControl" style="position:absolute;left:5px;top:5px;color:#d21717;cursor: pointer"><icon-font type="icon-PTZ" style="font-size: 32px;color: #c81010"></icon-font></div>
          <div v-if="isShowControl" class="ptz">
            <div class="outer-ring">
          <!-- 上 -->
          <a-icon type="caret-up" class="caret-up" @click="ptzController(0)" />
              <!-- 下 -->
              <a-icon type="caret-down" class="caret-down" @click="ptzController(1)" />
              <!-- 左 -->
              <a-icon type="caret-left" class="caret-left" @click="ptzController(2)" />
              <!-- 右 -->
              <a-icon type="caret-right" class="caret-right" @click="ptzController(3)" />
              <div class="inner-ring">
            <a-icon type="plus" class="plus" @click="ptzController(8)"/>
                <div class="line"></div>
                <a-icon type="minus" class="minus" @click="ptzController(9)"/>
          </div>
              <!--         左上 -->
              <!--         <a-icon type="caret-left" class="ob-caret-left" @click="ptzController(4)" />-->
              <!--         右上 -->
              <!--         <a-icon type="caret-up" class="ob-caret-up" @click="ptzController(6)"  />-->
              <!--         右下 -->
              <!--         <a-icon type="caret-right" class="ob-caret-right" @click="ptzController(7)" />-->
              <!--         左下 -->
              <!--         <a-icon type="caret-down" class="ob-caret-down" @click="ptzController(5)" />-->
        </div>
          </div>
        </LivePlayer>
      </div>
    </div>
          <div v-else-if="videoType==0" >
      <div class="h5videodiv" :style="{'width': width+'px','height': height+'px',}" >
        <LivePlayer :id="detail.identifier" :video-url="video_url"  :video-title="VideoTitle" :alt="VideoTitle+$t('configComponent.video.VideoOfflineTips')" aspect="fullscreen" fluent autoplay loop stretch>
        </LivePlayer>
      </div>
    </div>
          <div v-else-if="videoType==4||videoType==3" >
      <div class="h5videodiv" :style="{'width': width+'px','height': height+'px',}" >
        <LivePlayer :id="detail.identifier" :video-url="video_url"  :video-title="VideoTitle" :alt="VideoTitle+$t('configComponent.video.VideoOfflineTips')" aspect="fullscreen" fluent autoplay loop stretch>
        </LivePlayer>
      </div>
    </div>
          <div v-else>
      <div class="root" :style="{'width': width+'px','height': height+'px',}" >
      <div  :ref="detail.identifier"></div>
      </div>
    </div>
  </div>
        </div>



</template>

<script>
import {getVideoStatus,PtzControl } from "@/services/video";
import LivePlayer from '@liveqing/liveplayer'
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
export default {
  mixins: [ISMChildAutoMixin],
  name: 'ism-view-base-video',
  inject: ['getNode'],
  i18n: require('../../../../i18n/language'),
  props: {

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
    LivePlayer
  },
  data() {
    return {
      PMapState:null,
      videoServer:"",
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
      jessibuca: null,
      version: '',
      wasm: false,
      vc: "ff",
      playing: false,
      quieting: true,
      loaded: false, // mute
      showOperateBtns: true,
      showBandwidth: true,
      err: "",
      speed: 0,
      performance: "",
      volume: 1,
      rotate: 0,
      useWCS: false,
      useMSE: true,
      useOffscreen: false,
      recording: false,
      recordType: 'webm',
      scale: 0,
      isShowControl:false,
      token:"",
      VideoStatus:0,
      VideoTitle:"",
      videoType:1,
      isControling:false,
      video_url:"",
      temp_video_url:"",
      width:600,
      height:600,
      Ptzx:0,
      Ptzy:0,
      isRecord:0,
      Ptzz:0,
      isShowTitle:0,
      base:{
        "text": "configComponent.video.Text",
        "icon": "icon-shipinjiankong",
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
                "name":"configComponent.video.IsShowTitle",
                type:6,
                value:0,
                enumList:[
                  {
                    value:0,
                    option:"configComponent.video.ShowTitleNo"
                  },
                  {
                    value:1,
                    option:"configComponent.video.ShowTitleYes"
                  }
                ],
                "key":"isShowTitle",
              },
              {
                "name":"configComponent.video.Token",
                "type":8,
                "value":{
                  "type":1,
                  "key":""
                },
                "key":"token",
              }
            ]
          }
        }
      }
    }
  },
  methods: {
    createVideo(options) {
      options = options || {};
      let _t = this
      this.jessibuca = new window.Jessibuca(
          Object.assign(
              {
                container: _t.$refs[_t.detail.identifier],
                videoBuffer: Number(6), // 缓存时长
                isResize: false,
                useWCS: this.useWCS,
                useMSE: this.useMSE,
                text: "",
                loadingText: "疯狂加载中...",
                hasAudio:true,
                debug: false,
                supportDblclickFullscreen: true,
                showBandwidth: this.showBandwidth, // 显示网速
                operateBtns: {
                  fullscreen: this.showOperateBtns,
                  screenshot: this.showOperateBtns,
                  play: this.showOperateBtns,
                  audio: this.showOperateBtns,
                },
                vod: this.vod,
                forceNoOffscreen: !this.useOffscreen,
                isNotMute: true,
                timeout: 300
              },
              options
          )
      );
    },
    // 云台控制
    ptzController(conno) {
      let _t = this
      if(this.isControling)
      {
        return
      }

      if(conno==0)
      {
        this.Ptzy = this.Ptzy-0.1;
        if(this.Ptzy<=-1)
        {
          this.Ptzy=-1
        }
      }
      else if(conno==1)
      {
        this.Ptzy = this.Ptzy+0.1;
        if(this.Ptzy>=1)
        {
          this.Ptzy=1
        }
      }
      else if(conno==2)
      {
        this.Ptzx = this.Ptzx+0.1;
		if(this.Ptzx>=1)
        {
          this.Ptzx=1
        }
      }
      else if(conno==3)
      {
        this.Ptzx = this.Ptzx-0.1;
		if(this.Ptzx<=-1)
        {
          this.Ptzx=-1
        }
      }
      else if(conno==8)
      {
        this.Ptzz = this.Ptzz-0.1;
		if(this.Ptzx<=-1)
        {
          this.Ptzx=-1
        }
      }
      else if(conno==9)
      {
        this.Ptzz = this.Ptzz+0.1;
		if(this.Ptzx>=1)
        {
          this.Ptzx=1
        }
      }
      let data = {
        Ptzx:this.Ptzx,
        Ptzy:this.Ptzy,
        Ptzz:this.Ptzz,
        uuid:this.token
      }
      this.isControling = true
      PtzControl(data).then(res => {
        _t.isControling = false
      })
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
      if (this.jessibuca) {
        this.jessibuca.destroy();
      }
      this.width = option.style.position.w
      this.height = option.style.position.h
      for(let i=0;i<option.style.diy.length;i++)
      {
        if(option.style.diy[i].key=="token")
        {
          this.videoType = option.style.diy[i].value.type
          if(this.videoType==1)
          {
            this.token=option.style.diy[i].value.key
            this.video_url = "webrtc://"+this.videoServer+"/webrtcstream/"+this.token
          }
          else if(this.videoType==3)
          {
            this.video_url = option.style.diy[i].value.url
          }
          else if(this.videoType==4)
          {
            this.video_url = "http://"+location.host+"/"+option.style.diy[i].value.url
            console.log(this.video_url)
          }
          else
          {
            this.video_url = option.style.diy[i].value.key
          }
          this.VideoTitle = option.style.diy[i].value.name
          this.temp_video_url = this.video_url
        }
        else if(option.style.diy[i].key=="isShowTitle")
        {
          this.isShowTitle = option.style.diy[i].value
        }
      }
      if(this.videoType==3)
      {
        let _t = this
        setTimeout(function (){
          _t.createVideo()
          setTimeout(function (){
            _t.jessibuca.play(_t.video_url);
          },1000)
        },500)


      }
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
      if(this.videoType==1&&this.token!="") {
        this.getVideoStatus()
      }
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
    if (this.jessibuca) {
      this.jessibuca.destroy();
    }
    this.video_url=""
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
    _t.$EventBus.$on('cell-vuex', (data) => {
      _t.videoServer = data.PStore.state.setting.videoServer
      _t.initComponents(_t.detail)
    })
  }
}
</script>

<style lang="less">
.root {
  background: rgba(13, 14, 27, 0.7);
  display: flex;
  place-content: center;
}
.ptz{
  .outer-ring{
    z-index :9999;
    position: absolute;
    width: 150px;
    height: 150px;
    bottom: 10%;
    left: 2px;
    background-color: #fff;
    border-radius: 50%;
    box-shadow: inset 0 0 25px #e8e8e8, 0 1px 0 #c3c3c3, 0 2px 0 #c9c9c9, 0 2px 3px #333;
    i{
      font-size: 20px;
      cursor: pointer;
      &:hover{
        color: #00B176;
      }
    }
    .caret-up{
      position: absolute;
      top: 6px;
      color: #00ccff;
      left: 75px;
      margin-left: -10px;
    }
    .caret-down{
      position: absolute;
      bottom: 6px;
      left: 75px;
      color: #00ccff;
      margin-left: -10px;
    }
    .caret-left{
      position: absolute;
      left: 6px;
      top: 75px;
      color: #00ccff;
      margin-top: -10px;
    }
    .caret-right{
      position: absolute;
      right: 6px;
      top: 75px;
      color: #00ccff;
      margin-top: -10px;
    }
    .ob-caret-left{
      transform: rotate(45deg);
      position: absolute;
      top: 24px;
      color: #00ccff;
      left: 24px;
    }
    .ob-caret-up{
      transform: rotate(45deg);
      position: absolute;
      top: 24px;
      color: #00ccff;
      right: 24px;
    }
    .ob-caret-right{
      transform: rotate(45deg);
      position: absolute;
      bottom: 24px;
      color: #00ccff;
      right: 24px;
    }
    .ob-caret-down{
      transform: rotate(45deg);
      position: absolute;
      bottom: 24px;
      color: #00ccff;
      left: 24px;
    }
    .inner-ring{
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%,-50%);
      width: 80px;
      height: 80px;
      border-radius: 50%;
      background-color: #fff;
      border: 1px solid #ddd;
      .plus{
        position: absolute;
        top: 10px;
        color: #00ccff;
        left: 40px;
        margin-left: -10px;
      }
      .line{
        height: 1px;
        width: 100%;
        color: #00ccff;
        background-color: #ddd;
        position: absolute;
        top: 39px;
        left: 0;
      }
      .minus{
        position: absolute;
        bottom: 10px;
        color: #00ccff;
        left: 40px;
        margin-left: -10px;
      }
    }
  }
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
