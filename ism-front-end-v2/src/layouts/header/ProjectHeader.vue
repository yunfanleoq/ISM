<template>
  <a-layout-header :class="[headerTheme, 'admin-header']">
    <div :class="['admin-header-wide', layout, pageWidth]" style="position: relative;">
      <router-link  to="/Project" :class="['logo', isMobile ? null : 'pc', headerTheme]">
        <img width="32" :src="systemLogo" />
        <h1 >{{systemName}}</h1>
      </router-link>
      <a-divider v-if="isMobile" type="vertical" />
      <div v-if="layout !== 'side' && !isMobile" class="admin-header-menu" :style="`width: ${menuWidth};`">
        <i-menu class="head-menu" :theme="headerTheme" mode="horizontal" :options="menuData" @select="onSelect"/>
      </div>
      <div :class="['admin-header-right', headerTheme]">
          <header-avatar class="header-item"/>
          <a-dropdown class="lang header-item">
            <div>
              <a-icon type="global"/> {{langAlias}}
            </div>
            <a-menu @click="val => setLang(val.key)" :selected-keys="[lang]" slot="overlay">
              <a-menu-item v-for=" lang in langList" :key="lang.key">   {{lang.name}}</a-menu-item>
            </a-menu>
          </a-dropdown>
      </div>
    </div>
  </a-layout-header>
</template>

<script>
import HeaderAvatar from './HeaderAvatar'
import IMenu from '@/components/menu/menu'
import {mapState, mapMutations} from 'vuex'

import {formatDate} from '@/utils/common';
export default {
  name: 'ProjectHeader',
  components: {IMenu, HeaderAvatar},
  props: ['collapsed', 'menuData'],
  i18n: require('../../i18n/language'),
  data() {
    return {
      searchActive: false
    }
  },
  created(){
    let _this = this
    this.$EventBus.$off("RealAlarm")
    this.$EventBus.$on("RealAlarm", (data) => {
      let wsData = data
      let messageTitle = ''
      let messageContent = ''
      let AlarmLevel = ''
      let iconColor = '#0099FF'
      let SpeechAlarmContent = ""
      let HappenTime = _this.formatDateTime(wsData.HappenTime)
      if(wsData.Value=="1")
      {
        messageTitle=_this.$t('monitor.Notification.MessageAlarmTitle')
        messageContent = this.$t(wsData.AlarmMessage)
        SpeechAlarmContent = _this.$t('alarm.Speech.tips')
      }
      else{
        messageTitle=_this.$t('monitor.Notification.MessageClearAlarmTitle')
        messageContent = this.$t(wsData.AlarmClearMessage)
      }
      //此为摄像头的特殊数据ID需要单独处理
      if(wsData.DataUuid == "videoConnectStatusAlarm")
      {
        wsData.AlarmLevel = 3
        wsData.DataName=_this.$t('monitor.Notification.public.videoDataName')
        if (wsData.Value == "1")
        {
          messageContent = _this.$t('monitor.Notification.public.videoOffline')
        }
        else{
          messageContent =_this.$t('monitor.Notification.public.videoOnline')
        }
      }
      if(wsData.AlarmLevel==0)
      {
        iconColor = '#0099FF'
        AlarmLevel = _this.$t("dataModel.alarm.Tips")
      }
      else if(wsData.AlarmLevel==1)
      {
        iconColor = '#0066FF'
        AlarmLevel = _this.$t("dataModel.alarm.Minor")
      }
      else if(wsData.AlarmLevel==2)
      {
        iconColor = 'yellow'
        AlarmLevel =_this.$t("dataModel.alarm.Importance")
      }
      else if(wsData.AlarmLevel==3)
      {
        iconColor = 'orange'
        AlarmLevel = _this.$t("dataModel.alarm.Urgency")
      }
      else if(wsData.AlarmLevel==4)
      {
        iconColor = 'red'
        AlarmLevel = _this.$t("dataModel.alarm.Deadly")
      }
      SpeechAlarmContent = SpeechAlarmContent+_this.$t('monitor.Notification.DeviceTitle')
      SpeechAlarmContent = SpeechAlarmContent+_this.$t(wsData.DeviceName)
      SpeechAlarmContent = SpeechAlarmContent+_this.$t('monitor.Notification.DataTitle')
      SpeechAlarmContent = SpeechAlarmContent+_this.$t(wsData.DataName)
      SpeechAlarmContent = SpeechAlarmContent+_this.$t('dataModel.AlarmLevel')
      SpeechAlarmContent = SpeechAlarmContent+AlarmLevel
      SpeechAlarmContent = SpeechAlarmContent+_this.$t('monitor.Notification.HappenTime')
      SpeechAlarmContent = SpeechAlarmContent+HappenTime
      SpeechAlarmContent = SpeechAlarmContent+_this.$t('monitor.Notification.Message')
      SpeechAlarmContent = SpeechAlarmContent+messageContent
      let AlarmWindow = localStorage.getItem("AlarmWindow")

      if((AlarmWindow=="null")||(AlarmWindow==null)||(AlarmWindow==""))
      {
        AlarmWindow={}
        AlarmWindow.enable = true
        AlarmWindow.isClose = true
      }
      else
      {
        AlarmWindow = JSON.parse(AlarmWindow)
      }
      if(AlarmWindow.enable)
      {
        let isTrue = false
        if(typeof AlarmWindow.Level!="undefined") {
          isTrue = AlarmWindow.Level.find(item => item == wsData.AlarmLevel);
          if(AlarmWindow.Level.length==0)
          {
            isTrue=true
          }
        }else{
          isTrue=true
        }
        if(isTrue)
        {
          _this.$notification.warning({
            message: <span >{messageTitle}</span>,
            description:<div >
              <span style="font-size: 14px;font-weight: bold;">{_this.$t('monitor.Notification.DeviceTitle')}:</span><span>{_this.$t(wsData.DeviceName)}</span><br/>
              <span style="font-size: 14px;font-weight: bold;">{_this.$t('monitor.Notification.DataTitle')}:</span><span>{_this.$t(wsData.DataName)}</span><br/>
              <span style="font-size: 14px;font-weight: bold;">{_this.$t('dataModel.AlarmLevel')}:</span><span>{AlarmLevel}</span><br/>
              <span style="font-size: 14px;font-weight: bold;">{_this.$t('monitor.Notification.HappenTime')}:</span><span>{HappenTime}</span><br/>
              <span style="font-size: 14px;font-weight: bold;">{_this.$t('monitor.Notification.Message')}:</span><span>{messageContent}</span><br/>
            </div>,
            class: 'fancy-notification',
            placement:'bottomRight',
            duration:AlarmWindow.isClose?5:null,
            icon: <a-icon type="alert" theme="filled" style={{     'position': 'absolute',
              'margin-left': '4px',
              'font-size': '24px',
              'line-height': '24px','color':'red'}}/>,
            style: {
              padding: '10px 5px',
            },
          });
          _this.alarmSoundSpeech().start({container:"#sppekContent",Lang:_this.$i18n.locale,rate:1},SpeechAlarmContent);
        }
      }
    });

    this.$EventBus.$off("PlayVoice")
    this.$EventBus.$on("PlayVoice", (data) => {
      let wsData = data
      _this.alarmSoundSpeech().start({container:"#sppekContent",Lang:_this.$i18n.locale,rate:1},wsData.VoiceString);
    });
  },
  computed: {
    ...mapState('setting', ['langList','theme', 'isMobile', 'layout', 'systemName', 'lang', 'pageWidth']),
    headerTheme () {
      if (this.layout == 'side' && this.theme.mode == 'dark' && !this.isMobile) {
        return 'light'
      }
      return this.theme.mode
    },
    systemLogo () {
      return this.$store.state.setting.SystemLogo
    },
    langAlias() {
      let lang = this.langList.find(item => item.key == this.lang)
      return lang.alias
    },
    menuWidth() {
      const {layout, searchActive} = this
      const headWidth = layout === 'head' ? '100% - 188px' : '100%'
      const extraWidth = searchActive ? '600px' : '400px'
      return `calc(${headWidth} - ${extraWidth})`
    }
  },
  methods: {
    formatDateTime(time){
      let date = new Date(time)
      return formatDate(date,'yyyy-MM-dd hh:mm:ss')
    },
    toggleCollapse () {
      this.$emit('toggleCollapse')
    },
    alarmSoundSpeech : function() {
      let speechInstance = new SpeechSynthesisUtterance();
      return {
        start:function (opitions,content) {
          let getSpeech = localStorage.getItem("Speech")
          if((getSpeech=="null")||(getSpeech==null)||(getSpeech==""))
          {
            getSpeech={}
            getSpeech.enable = true
            getSpeech.speed = 1
          }
          else
          {
            getSpeech = JSON.parse(getSpeech)
          }

          if(getSpeech.enable)
          {
            let lang=opitions.Lang===undefined||""?"zh-CN":opitions.Lang;
            if(content!='') {
              speechSynthesis.cancel()
              speechInstance.text = content;
              speechInstance.lang = "zh-CN";
              speechInstance.pitch = 2;
              speechInstance.rate = getSpeech.speed;
              speechSynthesis.speak(speechInstance);
            }
          }
        }
      }
    },
    onSelect (obj) {
        this.$emit('menuSelect', obj)
      },
      ...mapMutations('setting', ['setLang'])
  },
  beforeDestroy() {
    speechSynthesis.cancel()
    this.$notification.destroy()
  }
}
</script>
<style lang="less" >
// 核心色彩体系（不变）
@primary: #13c2c2;
@primary-light: #e6fffa;
@primary-hover: #b5f5ec;
@primary-selected-bg: #0ea6a6;
@text-main: #1e293b;
@text-sub: #64748b;
@border: #e2e8f0;
@white: #ffffff;
// 美化后的 logo 样式
.logo {
  height: 50px;
  display: flex;
  align-items: center;
  padding: 0 20px;
  border-radius: 8px;
  position: relative;
  overflow: hidden;
  transition: all 0.4s cubic-bezier(0.25, 0.8, 0.25, 1);

  // logo 图片美化
  img {
    height: 36px;
    width: auto;
    object-fit: contain;
    border-radius: 4px;
    transition: all 0.4s ease;
  }

  &:hover img {
    transform: scale(1.12);
  }

  // 系统名称文字美化
  h1 {
    margin: 0 0 0 12px;
    font-size: 16px;
    color:#0d0d0d !important;
    font-weight: 700;
    font-family: '黑体', 'Microsoft YaHei', sans-serif;
    white-space: nowrap;
    transition: all 0.3s ease;
    line-height: 1;
  }

  &:hover h1 {
    transform: scale(1.05);
  }
}
/* 保持类名统一为 fancy-notification */
.fancy-notification {
  border: none !important;
  padding: 0 !important;
  margin: 0 !important;
  overflow: visible !important;
  position: relative;
}

/* 外层容器样式 */
.fancy-notification.ant-notification-notice {
  border: none !important;
  box-shadow: none !important;
}

/* 核心内容容器：科技感设计 */
.fancy-notification .ant-notification-notice-content {
  position: relative;
  background: #0f1736 !important; /* 深太空蓝背景 */
  border-radius: 8px !important;
  padding: 20px !important;
  padding-right: 45px !important; /* 关闭按钮空间 */
  border: 1px solid rgba(0, 240, 255, 0.3) !important;

  /* 霓虹发光边框 */
  box-shadow:
      0 0 5px rgba(0, 240, 255, 0.5),
      0 0 10px rgba(0, 240, 255, 0.3),
      inset 0 0 10px rgba(0, 240, 255, 0.2) !important;

  /* 网格纹理背景 */
  background-image:
      linear-gradient(rgba(15, 23, 54, 0.8) 1px, transparent 1px),
      linear-gradient(90deg, rgba(15, 23, 54, 0.8) 1px, transparent 1px) !important;
  background-size: 20px 20px !important;
  transition: all 0.3s ease;
}

/* 悬浮脉冲光效 */
.fancy-notification .ant-notification-notice-content::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  border-radius: 8px;
  border: 1px solid rgba(0, 240, 255, 0);
  box-shadow: 0 0 20px rgba(0, 240, 255, 0);
  transition: all 0.5s ease;
}
.fancy-notification:hover .ant-notification-notice-content::before {
  border-color: rgba(0, 240, 255, 0.8);
  box-shadow: 0 0 30px rgba(0, 240, 255, 0.6);
  animation: pulse 2s infinite;
}

/* 标题：故障艺术效果 */
.fancy-notification .ant-notification-notice-message {
  color: #00f0ff !important;
  font-size: 18px !important;
  font-weight: 600 !important;
  margin-bottom: 12px !important;
  display: flex;
  align-items: center;
  gap: 10px;
  text-shadow: 0 0 20px rgba(0, 240, 255, 0.7);
  position: relative;
  animation: glitch 3s infinite;
}
.fancy-notification .ant-notification-notice-message::before {
  content: attr(data-text);
  position: absolute;
  left: -1px;
  color: #ff00e4;
  opacity: 0.3;
}

/* 描述文字：终端风格 */
.fancy-notification .ant-notification-notice-description {
  color: #84f2ff !important;
  font-size: 16px !important;
  line-height: 1.6 !important;
  font-family: 'Consolas', 'Monaco', monospace !important;
  text-shadow: 0 0 2px rgba(132, 242, 255, 0.5);
}

/* 关闭按钮：科技感销毁按钮 */
.fancy-notification .ant-notification-notice-close {
  position: absolute !important;
  top: 18px !important;
  right: 18px !important;
  z-index: 10 !important;
  width: 26px !important;
  height: 26px !important;
  color: #ff4d4f !important;
  background: rgba(255, 77, 79, 0.1) !important;
  border-radius: 50% !important;
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  transition: all 0.3s;
  border: 1px solid rgba(255, 77, 79, 0.3) !important;
}
.fancy-notification .ant-notification-notice-close:hover {
  color: #fff !important;
  background: #ff4d4f !important;
  transform: scale(1.1) rotate(90deg);
  box-shadow: 0 0 15px rgba(255, 77, 79, 0.8) !important;
}

/* 入场动画 */
.fancy-notification {
  animation: techIn 0.6s cubic-bezier(0.16, 1, 0.3, 1) forwards;
  opacity: 0;
  transform: translateX(100%) scale(0.8);
}

/* 退场动画 */
.fancy-notification.ant-notification-notice-close-ping {
  animation: techOut 0.5s ease forwards !important;
}

/* 动画关键帧 */
@keyframes techIn {
  0% {
    opacity: 0;
    transform: translateX(100%) scale(0.8);
    box-shadow: 0 0 0 rgba(0, 240, 255, 0);
  }
  70% {
    opacity: 1;
    transform: translateX(-10px) scale(1.05);
    box-shadow: 0 0 20px rgba(0, 240, 255, 0.6);
  }
  100% {
    opacity: 1;
    transform: translateX(0) scale(1);
    box-shadow: 0 0 10px rgba(0, 240, 255, 0.4);
  }
}

@keyframes techOut {
  0% {
    opacity: 1;
    transform: translateX(0) scale(1);
  }
  100% {
    opacity: 0;
    transform: translateX(50px) scale(0.5);
  }
}

@keyframes pulse {
  0%, 100% { box-shadow: 0 0 20px rgba(0, 240, 255, 0.6); }
  50% { box-shadow: 0 0 40px rgba(0, 240, 255, 0.8); }
}

@keyframes glitch {
  0%, 95%, 100% { transform: translate(0); }
  96% { transform: translate(-1px, 1px); }
  97% { transform: translate(1px, -1px); }
  98% { transform: translate(-1px, -1px); }
  99% { transform: translate(1px, 1px); }
}

</style>
<style lang="less" scoped>
@import "index";
</style>