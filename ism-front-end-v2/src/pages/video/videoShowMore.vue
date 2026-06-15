<template>
  <div ref="videoDiv" :id="id">
    <a-layout >
      <a-card v-if="!fullScreen">
        <a-radio-group default-value="horizontal" @change="chargeVideoShow">
          <a-radio-button value="single">
            {{ $t('VideoManager.single') }}
          </a-radio-button>
          <a-radio-button value="four">
            {{ $t('VideoManager.four') }}
          </a-radio-button>
          <a-radio-button value="nine">
            {{ $t('VideoManager.nine') }}
          </a-radio-button>
          <a-radio-button value="sixth">
            {{ $t('VideoManager.sixth') }}
          </a-radio-button>
          <a-radio-button value="full" @click="toggleScreen">
            <a-icon class="action" :type="fullScreen ? 'fullscreen-exit' : 'fullscreen'" />
            <span v-if="fullScreen">{{ $t('VideoManager.ExitFull') }}</span>
            <span v-if="!fullScreen">{{ $t('VideoManager.full') }}</span>
          </a-radio-button>
        </a-radio-group>
      </a-card>
      <a-layout-content style="margin: 0 0px">
          <a-card style="padding: 0px;" id="viewCard">
            <a-list :grid="gridShow" :data-source="videoPlayerList">
              <a-list-item slot="renderItem" slot-scope="item, index">
                <a-card :style="videoStyle">
                  <LivePlayer custom-buttons="对讲" :video-url="item.url"  :video-title="item.Name" :alt="VideoTitle+$t('configComponent.video.VideoOfflineTips')" aspect="fullscreen" fluent autoplay live stretch>
                  </LivePlayer>
                  <div  class="video-close" style="display: none;">关闭</div>
                  <div  class="video-close channel-selector" @click="showSystemVideoModel(index)" >选择通道</div>
                </a-card>
              </a-list-item>
            </a-list>
          </a-card>
      </a-layout-content>
      <system-video-model @onSelectVideo="onSelectVideo" :networkVideo="videoComponentData" ref="systemVideoModel"></system-video-model>
    </a-layout>
  </div>
</template>

<script>
import LivePlayer from '@liveqing/liveplayer'
import systemVideoModel from '@/components/systemVideoModel/systemVideoModel'
import {mapState} from "vuex";
export default {
  name: 'VideoManager',
  i18n: require('../../i18n/language'),
  components: {
    LivePlayer,
    systemVideoModel
  },
  provide() {
    return {
      videoDiv: this
    }
  },
  data () {
    return {
      id: `${new Date().getTime()}-${Math.floor(Math.random() * 10)}`,
      fullScreen: false,
      windowWidth: document.documentElement.clientWidth,  //实时屏幕宽度
      windowHeight: document.documentElement.clientHeight,   //实时屏幕高度
      videoComponentData:{},
      gridShow:{
        gutter: 0,
        column: 1
      },
      VideoTitle:"",
      selectIndex:0,
      videoStyle:{
        height:300+"px",
        width:'auto'
      },
      videoType:1,
      video_url:"",
      temp_video_url:"",
      messageShowLoad:false,
      width:600,
      height:600,
      videoPlayerList:[],
      base:{
        "text": "configComponent.video.Text",
        "icon": "icon-shipinjiankong",
        "isFontIcon": true,
        "info": {
          "type": "image",
          "action": [],
          "dataBind":[],
          "style": {
            "position": {
              "x": 0,
              "y": 0,
              "w": 200,
              "h": 200
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
  authorize: {

  },
  computed: {
    ...mapState('setting', ['videoServer']),
  },
  mounted(){
    this.getSystemVideoList(1)
    let that = this;
    // <!--把window.onresize事件挂在到mounted函数上-->
    window.onresize = () => {
      return (() => {
        window.fullHeight = document.documentElement.clientHeight;
        window.fullWidth = document.documentElement.clientWidth;
        that.windowHeight = window.fullHeight;  // 高
        that.windowWidth = window.fullWidth; // 宽
        that.videoStyle.height = that.windowHeight+'px'
      })()
    };
  },
  activated(){

  },
  created(){
    this.addListener()
    this.videoStyle.height = document.documentElement.clientHeight+'px'
  },
  beforeDestroy() {
    this.removeListener()
  },
  watch: {
    windowHeight (val) {
      let that = this;
      this.videoStyle.height = (that.windowHeight/that.gridShow.column)+'px'
    },
    windowWidth (val) {
      let that = this;
    }
  },
  methods: {
    toggleScreen() {
      if (this.fullScreen) {
        this.outFullScreen()
      } else {
        this.inFullScreen()
      }
    },
    inFullScreen() {
      const el = this.$refs.videoDiv
      el.classList.add('beauty-scroll')
      if (el.requestFullscreen) {
        el.requestFullscreen()
        return true
      } else if (el.webkitRequestFullScreen) {
        el.webkitRequestFullScreen()
        return true
      } else if (el.mozRequestFullScreen) {
        el.mozRequestFullScreen()
        return true
      } else if (el.msRequestFullscreen) {
        el.msRequestFullscreen()
        return true
      }
      this.$message.warn('对不起，您的浏览器不支持全屏模式')

      return false
    },
    outFullScreen() {
      if (document.exitFullscreen) {
        document.exitFullscreen()
      } else if (document.webkitCancelFullScreen) {
        document.webkitCancelFullScreen();
      } else if (document.mozCancelFullScreen) {
        document.mozCancelFullScreen()
      } else if (document.msExitFullscreen) {
        document.msExitFullscreen()
      }
      this.$refs.videoDiv.classList.remove('beauty-scroll')
    },
    addListener() {
      document.addEventListener('fullscreenchange', this.fullScreenListener)
      document.addEventListener('webkitfullscreenchange', this.fullScreenListener)
      document.addEventListener('mozfullscreenchange', this.fullScreenListener)
      document.addEventListener('msfullscreenchange', this.fullScreenListener)
    },
    removeListener() {
      document.removeEventListener('fullscreenchange', this.fullScreenListener)
      document.removeEventListener('webkitfullscreenchange', this.fullScreenListener)
      document.removeEventListener('mozfullscreenchange', this.fullScreenListener)
      document.removeEventListener('msfullscreenchange', this.fullScreenListener)
    },
    fullScreenListener(e) {
      if (e.target.id === this.id) {
        this.fullScreen = !this.fullScreen
      }
    },
    showSystemVideoModel(index){
      this.$refs.systemVideoModel.showModal()
      this.selectIndex = index
    },
    chargeVideoShow(e){
      if(e.target.value=="single")
      {
        this.gridShow.column=1
        this.videoStyle.height = (this.windowHeight/1)+'px'
        this.getSystemVideoList(1)
      }
      else  if(e.target.value=="four")
      {
        this.gridShow.column=2
        this.videoStyle.height = (this.windowHeight/2)+'px'
        this.getSystemVideoList(4)
      }
      else  if(e.target.value=="nine")
      {
        this.gridShow.column=3
        this.getSystemVideoList(9)
        this.videoStyle.height = (this.windowHeight/3)+'px'
      }
      else  if(e.target.value=="sixth")
      {
        this.gridShow.column=4
        this.videoStyle.height = (this.windowHeight/4)+'px'
        this.getSystemVideoList(16)
      }
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
    onSelectVideo(token){
      console.log(token)
      if(token.type==1)
      {
        this.videoPlayerList[this.selectIndex].url = "webrtc://"+this.videoServer+"/webrtcstream/"+token.key
      }
      else if(token.type==3)
      {
        this.videoPlayerList[this.selectIndex].url = token.url
      }
      else if(token.type==4)
      {
        this.videoPlayerList[this.selectIndex].url = "http://"+location.host+"/"+token.url
      }
      else
      {
        this.videoPlayerList[this.selectIndex].url = token.key
      }
      this.videoPlayerList[this.selectIndex].Name = token.name
      this.visible =false
    },
  }
}
</script>

<style lang="less" scoped>
::v-deep .ant-card-body {
  padding: 0px;
  zoom: 1;
}
::v-deep  .ant-card-bordered {
   border: 1px solid #f0f0f0;
}
::v-deep .ant-list-grid .ant-col > .ant-list-item {
  margin-bottom: 0px;
}
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
  //right: 5px;
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
