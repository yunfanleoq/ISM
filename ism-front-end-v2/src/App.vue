<template>
  <a-config-provider :locale="locale" :get-popup-container="popContainer">
    <router-view/>
  </a-config-provider>
</template>

<script>
import {enquireScreen} from './utils/util'
import {mapState, mapMutations} from 'vuex'
import themeUtil from '@/utils/themeUtil';
import {getI18nKey} from '@/utils/routerUtil'

import {formatDate} from '@/utils/common';
import {checkAuthorization,getAuthorization} from '@/utils/request'
import {AUTH_TYPE} from "./utils/request";
import {GetSystemAuthInfo,GetSystemCodeCheck,GetSystemParams} from "@/services/system";


export default {
  name: 'App',
  data() {
    return {
      wsUrl: "",
      websocket:null,
      tempProject:"",
      wsProtocol:"ws",
      connectWebsocket:null,
      connectSSE:null,
      sseReconnectTimer: null,
      sseReconnectAttempt: 0,
      locale: {},
      eventSource: null, // SSE 连接实例
      messages: [],      // 存储接收的消息
      isConnected: false // 连接状态
    }
  },
  filters: {
    formatDate(time) {
      let date = new Date(time)
      return formatDate(date,'yyyy-MM-dd hh:mm:ss')
    },
  },
  created () {
    this.CheckPhysicalCode()
    this.GetSystemAuthInfo()
    this.GetSystemParams()
   // this.connectWs()
    this.setHtmlTitle()

    enquireScreen(isMobile => this.setDevice(isMobile))
    this.checkBrowser()
  },
  mounted() {
    this.setWeekModeTheme(this.weekMode)
    window.addEventListener('online', this.handleSSEOnline)
    document.addEventListener('visibilitychange', this.handleSSEVisibilityChange)
  },
  beforeDestroy() {
    window.removeEventListener('online', this.handleSSEOnline)
    document.removeEventListener('visibilitychange', this.handleSSEVisibilityChange)
    this.DisconnectSSE()
  },
  watch: {
    weekMode(val) {
      this.setWeekModeTheme(val)
    },
    lang(val) {
      this.setLanguage(val)
      this.setHtmlTitle()
    },
    $route() {
      this.setHtmlTitle()
    },
    'theme.mode': function(val) {
      let closeMessage = this.$message.loading(`您选择了主题模式 ${val}, 正在切换...`)
      themeUtil.changeThemeColor(this.theme.color, val).then(closeMessage)
    },
    'theme.color': function(val) {
      let closeMessage = this.$message.loading(`您选择了主题色 ${val}, 正在切换...`)
      themeUtil.changeThemeColor(val, this.theme.mode).then(closeMessage)
    },
    'layout': function() {
      window.dispatchEvent(new Event('resize'))
    }
  },
  computed: {
    ...mapState('setting', ['layout', 'isMobile','theme', 'weekMode', 'lang'])
  },
  methods: {
    ...mapMutations('setting', ['setDevice','setLang']),
    setSystemLanguage(lan){
      let getlang = localStorage.getItem("lang")
      if((getlang=="null")||(getlang==null)||(getlang==""))
      {
        getlang=lan
      }
      this.setLanguage(getlang)
      this.setLang(getlang)
    },
    setWeekModeTheme(weekMode) {
      if (weekMode) {
        document.body.classList.add('week-mode')
      } else {
        document.body.classList.remove('week-mode')
      }
    },
    formatDateTime(time){
      let date = new Date(time)
      return formatDate(date,'yyyy-MM-dd hh:mm:ss')
    },
    checkBrowser(){
      let userAgent = navigator.userAgent; //取得浏览器的userAgent字符串
      if (userAgent.indexOf("Chrome") > -1){//判断是否Chrome浏览器
        return "Chrome";
      }else{
        if(!this.isMobile){
          this.$message.error("请使用Microsoft Edge或者谷歌浏览器！！！", 30)
        }
      }
    },
    connectWs(){
      let _t = this
      if('https:' == document.location.protocol)
      {
        this.connectSSE = setInterval(function () {
          _t.ConnectSSE()
        }, 1000)
      }
      else
      {
        this.connectWebsocket = setInterval(function () {
          _t.webSocketLink()
        }, 1000)
      }
    },
    webSocketLink(){
      if((!checkAuthorization())&&(!checkAuthorization(AUTH_TYPE.AUTH1)))
      {
        return
      }
      clearInterval(this.connectWebsocket)
      let port=this.$store.state.setting.WSPort
      let Address=this.$store.state.setting.WSAddress
      if(Address =="local" || Address =="")
      {
        Address = location.hostname
      }
      this.wsProtocol = 'https:' == document.location.protocol ? 'wss': 'ws';
      if(document.location.protocol=="https:")
      {
        port=443
      }
      this.wsUrl= this.wsProtocol+'://'+Address+':'+port+'/ws?token='+getAuthorization()+"&project="+getAuthorization(AUTH_TYPE.AUTH1)
      let _this = this;
      let projectCheck = {
        timeout: 5000,//5秒
        timeoutObj: null,
        reset: function(){
          clearInterval(this.timeoutObj);
          return this;
        },
        start: function(){
          this.timeoutObj = setInterval(function(){
            if(checkAuthorization(AUTH_TYPE.AUTH1))
            {
              let projectUuid =getAuthorization(AUTH_TYPE.AUTH1)
               if(_this.tempProject!=projectUuid)
               {
                 _this.tempProject = projectUuid
                 wsInit.connect()
               }else if(_this.websocket==null)
               {
                 wsInit.connect()
               }
               else if((_this.websocket.readyState!=1)){
                 wsInit.connect()
               }
            }
            else{
              if (_this.websocket) {
                _this.websocket.onclose = null
                _this.websocket.onmessage = null
                _this.websocket.onerror = null
                console.log("没有授权")
                _this.websocket.close()
                _this.websocket = null
              }
            }
          }, this.timeout)
        }
      };
      let wsInit={
        timerWait:null,
        waitMsg:false,
        reConnectTimer:null,
        connect:function (){
          if(!checkAuthorization(AUTH_TYPE.AUTH1))
          {
            projectCheck.reset().start();
            return
          }
          try{
            let port=_this.$store.state.setting.WSPort
            let Address=_this.$store.state.setting.WSAddress
            if(Address =="local" || Address =="")
            {
              Address = location.hostname
            }
            if(document.location.protocol=="https:")
            {
              port=443
            }
            _this.tempProject = getAuthorization(AUTH_TYPE.AUTH1)
            _this.wsUrl= _this.wsProtocol+'://'+Address+':'+port+'/ws?token='+getAuthorization()+"&project="+getAuthorization(AUTH_TYPE.AUTH1)
            if (_this.websocket) {
              _this.websocket.onclose = null
              _this.websocket.onmessage = null
              _this.websocket.onerror = null
              console.log("close")
              _this.websocket.close()
              _this.websocket = null
            }
            _this.websocket = new WebSocket(_this.wsUrl);
            _this.websocket.onopen = function(){
              console.log("websock连接成功");
              clearInterval(wsInit.reConnectTimer);

              _this.websocket.onclose = wsInit.onClose
              _this.websocket.onmessage =wsInit.onMessage
              _this.websocket.onerror = wsInit.onError
              projectCheck.reset().start();
            };
          }catch (e){
            projectCheck.reset().start();
            _this.websocket = null
          }
        },
        onClose:function (e){
          console.log("websock关闭");
          _this.websocket.onclose = null
          _this.websocket.onmessage = null
          _this.websocket.onerror = null
          _this.websocket.close()
          _this.websocket = null
        },
        onError:function (){
          console.log("websock错误");
          _this.websocket.onclose = null
          _this.websocket.onmessage = null
          _this.websocket.onerror = null
          _this.websocket.close()
          _this.websocket = null
        },
        onMessage:function (event) {
          this.waitMsg = true
          clearTimeout(this.timerWait)
          try{
            let wsData = JSON.parse(event.data)
            if(wsData.Cmd=="RealData")
            {
              _this.$EventBus.$emit("readDataPush", wsData);
            }
            else if(wsData.Cmd=="RealAlarm")
            {
              _this.$EventBus.$emit("RealAlarm", wsData);
            }
            else if(wsData.Cmd=="SystemData")
            {
              _this.$EventBus.$emit("SystemData", wsData);
            }
            else if(wsData.Cmd=="StaticData")
            {
              _this.$EventBus.$emit("StaticData", wsData);
            }
            else if(wsData.Cmd=="PlayVoice")
            {
              _this.$EventBus.$emit("PlayVoice", wsData);
            }
            else if(wsData.Cmd=="GoPage")
            {
              _this.$EventBus.$emit("GoPage", wsData);
            }
          }catch (e) {
            console.log(e)
          }
        }
      }
      wsInit.connect()
    },
    setLanguage(lang) {
      this.$i18n.locale = lang
      switch (lang) {
        case 'CN':
          this.locale = require('ant-design-vue/es/locale-provider/zh_CN').default
          break
        case 'HK':
          this.locale = require('ant-design-vue/es/locale-provider/zh_TW').default
          break
        case 'US':
        default:
          this.locale = require('ant-design-vue/es/locale-provider/en_US').default
          break
      }
    },
    setHtmlTitle() {
      const route = this.$route
      if (!route || !route.path) {
        return
      }
      const matched = route.matched
      const last = matched && matched.length ? matched[matched.length - 1] : null
      const key = route.path === '/'
        ? 'home.name'
        : getI18nKey((last && last.path) || route.path)
      document.title = this.$store.state.setting.SystemAuthAPPName
        ? this.$store.state.setting.SystemAuthAPPName + ' | ' + this.$t(key)
        : 'Loading' + ' | ' + this.$t(key)
    },
    popContainer() {
      return document.getElementById("popContainer")
    },
    CheckPhysicalCode(){
      let _t = this
      GetSystemCodeCheck().then(function (res) {
        _t.$store.state.setting.ProtectedID = res.data.id
        if(res.data.code==-1)
        {
          _t.$message.error("授权信息丢失，请联系管理员获取.", 6)
          _t.$router.push('/auth')
        }
        else if(res.data.code==-3)
        {
          _t.$message.error("您的软件授权已过期，请联系管理员续费或更新授权!", 6)
          _t.$router.push('/auth')
        }
      }).catch(function(e){
        _t.$message.error("授权信息丢失了，请联系管理员重新获取.", 6)
        _t.$router.push('/auth')
      })
    },
    GetSystemAuthInfo(){
      let _t = this
      _t.$store.state.setting.skeletonLoading = true
      GetSystemAuthInfo().then(function (res) {
        let auth = res.data || {}
        // 循安现场品牌：后端 applyXunanBranding 为准；前端再兜底去掉「零界/zerobound」
        const xunanLogo = '/static/branding/logo-xunan-hexagon.png'
        const looksLegacy = /零界|zerobound|临界|ai\s*$/i.test(
          String(auth.systemName || '') + ' ' + String(auth.SystemAPPName || '') + ' ' + String(auth.systemCompany || '')
        )
        _t.$store.state.setting.systemName = looksLegacy || !auth.systemName
          ? '循安电力监控平台'
          : auth.systemName
        _t.$store.state.setting.systemUrl = auth.systemUrl
        _t.$store.state.setting.SystemDynamicUrl = auth.SystemDynamicUrl
        _t.$store.state.setting.SystemLogo = (looksLegacy || !auth.SystemLogo || String(auth.SystemLogo).indexOf('data:image') === 0)
          ? xunanLogo
          : auth.SystemLogo
        _t.$store.state.setting.systemLoginBg = auth.systemLoginBg
        _t.$store.state.setting.SystemAuthAPPName = looksLegacy || !auth.SystemAPPName
          ? '循安'
          : auth.SystemAPPName
        _t.$store.state.setting.systemCompany = looksLegacy || !auth.systemCompany
          ? '北京循安科技有限公司'
          : auth.systemCompany
        _t.$store.state.setting.systemBg = auth.systemBg
        _t.$store.state.setting.skeletonLoading = false
        _t.setHtmlTitle()
      }).catch(function(e){
        _t.$message.error("授权信息丢失了，请联系管理员重新获取.", 6)
      })
    },
    GetSystemParams(){
      let _t = this
      GetSystemParams().then(function (res) {
        _t.$store.state.setting.WSPort = res.data.list.WSPort
        _t.$store.state.setting.WSAddress = res.data.list.WSAddress
        _t.$store.state.setting.IsLicense = res.data.list.IsLicense
        _t.$store.state.setting.IsOEM = res.data.list.IsOEM
        let lang = res.data.list.DefaultLang
        _t.connectWs()
        _t.setSystemLanguage(lang)
      }).catch(function(e){
        _t.connectWs()
      })
    },
    // 建立 SSE 连接
    ConnectSSE() {
      if((!checkAuthorization())&&(!checkAuthorization(AUTH_TYPE.AUTH1)))
      {
        return
      }
      clearInterval(this.connectSSE)
      if (!navigator.onLine || document.visibilityState === 'hidden') {
        return
      }
      if (this.eventSource &&
        (this.eventSource.readyState === EventSource.CONNECTING ||
          this.eventSource.readyState === EventSource.OPEN)) {
        return
      }
      clearTimeout(this.sseReconnectTimer)
      this.sseReconnectTimer = null

      // SSE 与当前页面同源，生产环境由反向代理转发到 Beego。
      const source = new EventSource(
        `${window.location.origin}/SSEPushData?token=${encodeURIComponent(getAuthorization())}`
      )
      this.eventSource = source

      // 2. 监听连接成功
      source.onopen = () => {
        console.log("SSE 连接已打开");
        this.isConnected = true;
        this.sseReconnectAttempt = 0;
      };
      let _this = this;
      // 3. 监听后端推送的消息
      source.onmessage = (event) => {
        try {
          let wsData = JSON.parse(event.data)
          if(wsData.Cmd=="RealData")
          {
            _this.$EventBus.$emit("readDataPush", wsData);
          }
          else if(wsData.Cmd=="RealAlarm")
          {
            _this.$EventBus.$emit("RealAlarm", wsData);
          }
          else if(wsData.Cmd=="SystemData")
          {
            _this.$EventBus.$emit("SystemData", wsData);
          }
          else if(wsData.Cmd=="StaticData")
          {
            _this.$EventBus.$emit("StaticData", wsData);
          }
          else if(wsData.Cmd=="PlayVoice")
          {
            _this.$EventBus.$emit("PlayVoice", wsData);
          }
          else if(wsData.Cmd=="GoPage")
          {
            _this.$EventBus.$emit("GoPage", wsData);
          }
        } catch (e) {
          console.error("解析 SSE 数据失败：", e);
        }
      };

      // 4. 监听错误（如连接断开）
      source.onerror = (error) => {
        if (_this.eventSource !== source) return
        console.error("SSE 错误：", error);
        _this.isConnected = false;
        // 关闭 EventSource 自带重连，统一交给带退避的单一重连定时器。
        source.close()
        _this.eventSource = null
        _this.scheduleSSEReconnect()
      };
    },
    scheduleSSEReconnect() {
      if (this.sseReconnectTimer || !navigator.onLine ||
        document.visibilityState === 'hidden') {
        return
      }
      const delays = [1000, 2000, 5000, 10000, 30000]
      const delay = delays[Math.min(this.sseReconnectAttempt, delays.length - 1)]
      this.sseReconnectAttempt += 1
      this.sseReconnectTimer = setTimeout(() => {
        this.sseReconnectTimer = null
        this.ConnectSSE()
      }, delay)
    },
    handleSSEOnline() {
      this.sseReconnectAttempt = 0
      this.ConnectSSE()
    },
    handleSSEVisibilityChange() {
      if (document.visibilityState === 'visible' && !this.isConnected) {
        this.ConnectSSE()
      }
    },
    DisconnectSSE() {
      clearInterval(this.connectSSE)
      clearTimeout(this.sseReconnectTimer)
      this.sseReconnectTimer = null
      if (this.eventSource) {
        this.eventSource.close();
        this.eventSource = null;
        this.isConnected = false;
      }
    }
  }
}
</script>

<style lang="less">
  .ant-notification-notice-message{
    color: #e10a40;
    font-size: 24px;
    font-weight: bold;
  }
</style>
