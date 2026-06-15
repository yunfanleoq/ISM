<template>
<div>
  <div class="run-login-graph-container" :class="{'animated':true,[`${configData.layer.animate}`]: true}" >
    <div ref="ISMRunningLoginContainer" ></div>
  </div>
</div>
</template>

<script>
import {formatDate} from "@/utils/common";

const loadingKey = 'updatable'
import ISMBase from './ISMBase';
import store from "@/store";
import { mapActions, mapState, mapMutations ,mapGetters} from 'vuex'
import Contextmenu from "vue-contextmenujs"
import Vue from 'vue'
Vue.use(Contextmenu);
import {ExecSysScript, GetDisplayLoginPage} from "@/services/system";
import {register} from "@antv/x6-vue-shape";
import ISMGroupNode from "@/pages/ISMDisPlay/ISMGroupNode.vue";
import {Graph} from "@antv/x6";
import {AUTH_TYPE, checkAuthorization, setAuthorization} from "@/utils/request";
import {getRealDataByUuid, setData} from "@/services/device";
import TemplateRender from "vue-template-render";
import {ComponentRestApi} from "@/services/RestApi";
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
export default {
  mixins: [ISMChildAutoMixin],
  name: 'ISMRenderLogin',
  extends: ISMBase,
  i18n: require('../../i18n/language'),
  components: {

  },
  props: {
    showUuid: {
      type: String,
      required: true
    },
    showDeviceUuid: {
      type: String,
      required: true
    }
  },
  watch: {
    showUuid: {
      handler(newVal, oldVal) {
        if(newVal!="")
        {
          let _t =this
          let page={
            pageType:1,
            displayUUID:newVal,
            pageUuid:this.showPageUUID
          }
          _t.chargePage = true
          this.$message.loading({ content: 'Loading...',loadingKey,duration: 0 });
          this.selectDisplayPageDataStruct({page:page,callback:function (uuids,devices,isFound){
            _t.$message.destroy();
            _t.chargePage = false
            if (isFound === false) {
              _t.$message.error(_t.$t("readData.NotFindPage"))
            }
          }})
        }
      },
      deep: true
    },
    fullScreen:{
      handler(newVal, oldVal) {
        let _t = this
        if(newVal==false)
        {
          const popContainer = document.getElementById("popContainer")
          popContainer.style.overflow="visible"
        }

        if(this.configData.layer.autoSize==1)
        {
          setTimeout(function (){
            _t.setScale()
          },1000)
        }
      },
      deep: true
    }
  },
  computed: {
    ...mapState('setting', ['langList','isMobile','lang',]),
    ...mapGetters('account', ['user']),
    ...mapState({
      selectedValue:state => store.state.ISMDisPlayEditorTool.selectedValue,
      configData: state => store.state.ISMDisPlayEditorTool.LayerData,
      PopUpConfigData: state => store.state.ISMDisPlayEditorTool.PopUpConfigData,
    }),
    layerStyle:function () {
      let scale = 1
      let styles = [`transform:scale(${scale})`];
      if(this.configData.layer.backColor) {
        styles.push(`background-color: ${this.configData.layer.backColor}`);
      }
      if(this.configData.layer.backgroundImage) {
        styles.push(`background-image: url("${this.configData.layer.backgroundImage}")`);
        styles.push(`background-size:100% 100%`);
      }
      if(this.configData.layer.width > 0) {
        styles.push(`width: ${this.configData.layer.width}px`);
      }
      if(this.configData.layer.height > 0) {
        styles.push(`height: ${this.configData.layer.height}px`);
      }
      let style = styles.join(';');
      return style;
    }
  },
  data() {
    return {
      tempAutoSize:0,
      tempAutoPadding:0,
      SelectCurrentNodeData:null,
      ISMPopUpRunningContainer:null,
      ISMRuningLoginCavasContainer:null,
      chargePage:false,
      Zoom:100,
      clickTimer:null,
      linkInfoExternal:{
        url:"",
        width:1024,
        height:768
      },
      chargePagePopUp:false,
      isExternUrl:false,
      isPopUpOpen:false,
      SelectDeviceUuid:"",
      showPageUUID:"",
      showAppUUID:"",
      deviceUuid:"",
      SetPassword:"",
      AutoSetValue:"",
      setDataUuid:"",
      settingVisible:false,
      settingLoading:false,
      SetForm:this.$form.createForm(this),
      SetFormPassword:this.$form.createForm(this),
      fullScreen:false
    }
  },
  created(){
    let _t = this
    this.addListener()
    if(this.showUuid!="")
    {
      this.chargePage = true
    }
    this.getDisplayLoginPage()
  },
  methods: {
    ...mapMutations('ISMDisPlayEditorTool',[
      'setlayerZoom',
    ]),
    ...mapMutations('account', ['setUser', 'setPermissions', 'setRoles']),
    ...mapActions('ISMDisPlayEditorTool',[
      'getLoginLayerDataStruct',
    ]),
    RunLoginCavasContainerInit(auto,panning){
      let _t = this
      let pagerSize = auto==1?true:false
      let pagerPanning = panning==1?true:false
      if(this.ISMRuningLoginCavasContainer!=null)
      {
        this.ISMRuningLoginCavasContainer.dispose();
        this.ISMRuningLoginCavasContainer = null
      }
      register({
        shape: "view-ism-group-node",
        width: 100,
        height: 100,
        component: ISMGroupNode,
        data: {
          locked:false,
          UpdateNodeFlag:true,
          editMode: true,
          showDeviceUuid:"",
          IsToolBox:false,
          detail:ISMGroupNode
        }
      })
      this.ISMRuningLoginCavasContainer = new Graph({
        container: this.$refs.ISMRunningLoginContainer,
        width: _t.configData.layer.width,
        height: _t.configData.layer.height,
        panning: pagerPanning,
        mousewheel: pagerPanning,
        autoResize:pagerSize,
        grid: false,
        connecting: {
          router: 'manhattan',
          connector: {
            name: 'rounded',
            args: {
              radius: 8,
            },
          },
          anchor: 'center',
          connectionPoint: 'anchor',
          allowBlank: false,
          snap:false
        },
        background: {
          color: _t.configData.layer.backColor,   // 背景底色（可选）
          image: _t.configData.layer.backgroundImage,
          size: 'cover',
          repeat: 'no-repeat',
          quality:1,
        },
        interacting: (cellView) => {
          return {
            nodeMovable:false,
            magnetConnectable: false,  // 禁止节点移动
            edgeMovable: false,  // 禁止边移动
            edgeLabelMovable: false,  // 禁止箭头移动
            arrowheadMovable: false,  // 禁止顶点移动
            vertexMovable: false,  // 禁止添加顶点
            vertexAddable: false,  // 禁止删除顶点
            vertexDeletable: false,  // 禁止删除顶点
          };
        },
      })
      this.ISMRuningLoginCavasContainer.off('node:click')
      this.ISMRuningLoginCavasContainer.off('node:dblclick')
      this.ISMRuningLoginCavasContainer.off('node:mousedown')
      this.ISMRuningLoginCavasContainer.off('node:mouseup')
      this.ISMRuningLoginCavasContainer.off('node:mouseenter')
      this.ISMRuningLoginCavasContainer.off('node:mouseleave')
      this.ISMRuningLoginCavasContainer.off('render:done')
      this.ISMRuningLoginCavasContainer.off('blank:contextmenu')
      this.ISMRuningLoginCavasContainer.off('cell:contextmenu')
      this.ISMRuningLoginCavasContainer.off('scale')
      window.removeEventListener('resize', this.setScale);
      window.addEventListener('resize', this.setScale);
      this.ISMRuningLoginCavasContainer.on('scale', ( scale ) => {

      });
      //
      this.ISMRuningLoginCavasContainer.on('node:click', ( { e, x, y, node, view } ) => {
        const component = node.prop().data.detail
        _t.SelectCurrentNodeData = component
        _t.doComponentAction(component,'click')
      });
      this.ISMRuningLoginCavasContainer.on('node:dblclick', ( { e, x, y, node, view } ) => {
        const component = node.prop().data.detail
        _t.SelectCurrentNodeData = component
        _t.doComponentAction(component,'dblclick')
      });
      this.ISMRuningLoginCavasContainer.on('node:mousedown', ( { e, x, y, node, view } ) => {
        const component = node.prop().data.detail
        _t.SelectCurrentNodeData = component
        _t.doComponentAction(component,'mousedown')
      });
      this.ISMRuningLoginCavasContainer.on('node:mouseup', ( { e, x, y, node, view } ) => {
        const component = node.prop().data.detail
        _t.SelectCurrentNodeData = component
        _t.doComponentAction(component,'mouseup')
      });
      this.ISMRuningLoginCavasContainer.on('node:mouseenter', ( { e, x, y, node, view } ) => {
        const component = node.prop().data.detail
        _t.SelectCurrentNodeData = component
        _t.doComponentAction(component,'mouseenter')
      });
      this.ISMRuningLoginCavasContainer.on('node:mouseleave', ( { e, x, y, node, view } ) => {
        const component = node.prop().data.detail
        _t.SelectCurrentNodeData = component
        _t.doComponentAction(component,'mouseleave')
      });
      //页面空白处右键菜单
      this.ISMRuningLoginCavasContainer.on('blank:contextmenu', ({ e, x, y, cell, view }) => {
        this.onContextLayerMenu(e)
      })
      this.ISMRuningLoginCavasContainer.on('cell:contextmenu', ({ e, x, y, cell, view }) => {
        this.onContextLayerMenu(e)
      })
      //==============================
      //监听页面是否渲染完成，完成后就发广播给节点，告诉节点是在编辑模式
      this.ISMRuningLoginCavasContainer.on('render:done', () => {
        this.tempAutoSize=this.configData.layer.autoSize
        this.setScale()
        this.$EventBus.$emit('cell-editMode',{
          edit:false,
          toolbox:false
        })
        this.$EventBus.$emit('cell-vuex',{
          PMapState:mapState,
          PMapActions:mapActions,
          PMapMutations:mapMutations,
          PStore:this.$store,
          loginUuid:this.showUuid,
          PRouter:this.$router
        })
      });
      try{
        const components = JSON.parse(JSON.stringify(_t.configData.components))
        if (components.cells && Array.isArray(components.cells)) {
          components.cells = components.cells.filter(cell => cell && cell.shape)
        }
        _t.ISMRuningLoginCavasContainer.fromJSON(components)
      }catch (e){
        _t.$message.error(_t.$t('Render.GetPageError'))
        console.log(e)
      }
    },
    formatDateTime(time){
      let date = new Date(time)
      return formatDate(date,'yyyy-MM-dd hh:mm:ss')
    },
    loadPager(){
      let _t = this
      if(this.showUuid!="")
      {
        if(this.showToken!="")
        {
          this.chargePage = true
          this.$message.loading({content: 'Loading...', loadingKey, duration: 0});
          this.CurrentRealUUIDList=[]
          this.CurrentModelUUIDList=[]
          this.getLayerDataStructByTokenData({
            pageType: this.isMobile, token:this.showToken,uuid: this.showUuid, cb: function (errno, project_uuid,expireAt,token, datauuid,devices) {
              document.title = _t.configData.AppName
              _t.$message.destroy()
              _t.chargePage = false
              if (errno == 0) {
                setAuthorization({token: project_uuid}, AUTH_TYPE.AUTH1)
                setAuthorization({token: token, expireAt: expireAt})
                setAuthorization({token: _t.showToken}, AUTH_TYPE.AUTH3)
              }
              else  if (errno == -5)
              {
                document.title = _t.$t('displayConfig.Properties.NoAuth')
                _t.$message.warn(_t.$t('displayConfig.Properties.NoAuth'))
                return
              }
              else  if (errno == -4)
              {
                document.title = _t.$t('displayConfig.Properties.NoToken')
                _t.$message.warn(_t.$t('displayConfig.Properties.NoToken'))
                return
              }
              else if (errno != 0)
              {
                return
              }
              _t.RunCavasContainerInit(_t.configData.layer.autoSize,_t.configData.layer.Padding)
              setTimeout(function (){
                let newuuids=[]
                for(let i=0;i<datauuid.length;i++)
                {
                  if((typeof datauuid[i]!="undefined")&&(datauuid[i].length>0))
                  {
                    newuuids.push(datauuid[i])
                  }
                }
                let newdevices=[]
                for(let i=0;i<devices.length;i++)
                {
                  if((typeof devices[i]!="undefined")&&(devices[i].length>0))
                  {
                    newdevices.push(devices[i])
                  }
                }
                getRealDataByUuid({uuid: newuuids,devices:newdevices}).then(function (res) {
                  _t.settingLoading = false
                  if (res.data.code == 0) {
                    for (let k = 0; k < res.data.realData.length; k++) {
                      if(res.data.realData[k].value=="")
                      {
                        continue
                      }
                      let pushData = {
                        DeviceUuid: res.data.realData[k].duid,
                        ProjectUuid: res.data.realData[k].project_uuid,
                        Cmd: "RealData",
                        Data: []
                      }
                      _t.CurrentRealUUIDList.push(res.data.realData[k].uuid)
                      _t.CurrentModelUUIDList.push(res.data.realData[k].mduid)
                      let DataObj = {
                        Uuid: res.data.realData[k].uuid,
                        ModelDataUuid: res.data.realData[k].mduid,
                        Value: res.data.realData[k].value
                      }
                      pushData.Data.push(DataObj)
                      _t.$EventBus.$emit("readDataPush", pushData);
                    }
                  }
                })
              },300)
            }
          });
        }
        else {
          this.$message.loading({content: 'Loading...', loadingKey, duration: 0});
          this.CurrentRealUUIDList=[]
          this.CurrentModelUUIDList=[]
          this.getLayerDataStruct({
            pageType: this.isMobile, uuid: this.showUuid, cb: function (errno, project_uuid, datauuid,devices) {
              document.title = _t.configData.AppName
              if (errno == 0) {
                if (!checkAuthorization(AUTH_TYPE.AUTH1)) {
                  setAuthorization({token: project_uuid}, AUTH_TYPE.AUTH1)
                }
              } else {
                _t.$router.push('/login')
                _t.$message.destroy()
                return
              }
              _t.$message.destroy()
              _t.RunCavasContainerInit(_t.configData.layer.autoSize,_t.configData.layer.Padding)
              setTimeout(function (){
                let newuuids=[]
                for(let i=0;i<datauuid.length;i++)
                {
                  if((typeof datauuid[i]!="undefined")&&(datauuid[i].length>0))
                  {
                    newuuids.push(datauuid[i])
                  }
                }
                let newdevices=[]
                for(let i=0;i<devices.length;i++)
                {
                  if((typeof devices[i]!="undefined")&&(devices[i].length>0))
                  {
                    newdevices.push(devices[i])
                  }
                }
                getRealDataByUuid({uuid: newuuids,devices:newdevices}).then(function (res) {
                  _t.settingLoading = false
                  if (res.data.code == 0) {
                    for (let k = 0; k < res.data.realData.length; k++) {
                      if(res.data.realData[k].value=="")
                      {
                        continue
                      }
                      let pushData = {
                        DeviceUuid: res.data.realData[k].duid,
                        ProjectUuid: res.data.realData[k].project_uuid,
                        Cmd: "RealData",
                        Data: []
                      }
                      _t.CurrentRealUUIDList.push(res.data.realData[k].uuid)
                      _t.CurrentModelUUIDList.push(res.data.realData[k].mduid)
                      let DataObj = {
                        Uuid: res.data.realData[k].uuid,
                        ModelDataUuid: res.data.realData[k].mduid,
                        Value: res.data.realData[k].value
                      }
                      pushData.Data.push(DataObj)
                      _t.$EventBus.$emit("readDataPush", pushData);
                    }
                  }
                })
              },300)
            }
          });
        }
      }
    },
    onKeyPress(button){
      if (button === '{enter}' || button === '{close}') {
        this.closekeyboard()
      }
    },
    // inpuit获取焦点显示虚拟键盘
    onInputFocus(res) {
      this.showKeyboard = true
      this.changeIpt = res
      // 父组件调用子组件的方法
      this.$refs.SimpleKeyboard.onKeyPress('{clear}')
    },
    // 给输入框赋值
    onChangeKeyboard(input) {
      if (this.changeIpt == 'SetPassword') {
        this.SetPasswordFormValue = input
      } else if (this.changeIpt == 'SetValue') {
        this.SetValueFormValue = input
      }else if (this.changeIpt == 'SetAutoPassword') {
        this.SetAutoPasswordFormValue = input
      } else if (this.changeIpt == 'SetActionPassword') {
        this.ActionPasswordValue = input
      }

    },
    // 点击关闭隐藏键盘
    closekeyboard() {
      this.showKeyboard = false
    },
    onContextLayerMenu(event) {
      let _t = this
      if(!this.isCommentRightClick) {
        this.$contextmenu({
          items: [
            {
              label: _t.$t('Render.FullScreen'),
              icon: "el-icon-full-screen",
              onClick: () => {
                _t.inFullScreen()
              }
            },
            {
              label: _t.$t('Render.ExitFullScreen'),
              icon: "el-icon-tuichuquanping",
              onClick: () => {
                _t.outFullScreen()
              }
            },
            {
              label: _t.$t('Render.Reload'),
              icon: "el-icon-lunxun",
              onClick: () => {
                _t.loadPager()
              }
            },
            {
              label: _t.tempAutoSize?_t.$t('Render.ScreenUndo'):_t.$t('Render.AutoSize'),
              icon: "el-icon-agora_zishiyingshipinfenbianshuai",
              onClick: () => {
                _t.tempAutoSize = !_t.tempAutoSize
                _t.setScale()
              }
            },
            {
              label: _t.tempAutoPadding?_t.$t('Render.NodeDragDisable'):_t.$t('Render.NodeDragEnable'),
              icon: _t.tempAutoPadding?"el-icon-jinzhi":"el-icon-tuodong",
              onClick: () => {
                _t.tempAutoPadding=!_t.tempAutoPadding
                this.ISMRuningLoginCavasContainer.togglePanning(_t.tempAutoPadding)
                this.ISMRuningLoginCavasContainer.toggleMouseWheel(_t.tempAutoPadding)
              }
            },
          ],
          event, // 鼠标事件信息
          divided:true,
          customClass: "custom-class", // 自定义菜单 class
          zIndex: 10000, // 菜单样式 z-index
          minWidth: 230 // 主菜单最小宽度
        });
        return false;
      }
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
    inFullScreen() {
      const el = document.body
      const popContainer = document.getElementById("popContainer")
      popContainer.style.overflow="auto"
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
      const el = document.body
      const popContainer = document.getElementById("popContainer")
      popContainer.style.overflow="visible"
      if (el.exitFullscreen) {
        el.exitFullscreen()
      } else if (document.webkitCancelFullScreen) {
        document.webkitCancelFullScreen();
      } else if (document.mozCancelFullScreen) {
        document.mozCancelFullScreen()
      } else if (document.msExitFullscreen) {
        document.msExitFullscreen()
      }
      this.$refs.ismrender.classList.remove('beauty-scroll')
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
      // const el = document.body
      this.fullScreen = document.body.fullScreen || document.mozFullScreen || document.webkitIsFullScreen;
      if(!this.fullScreen){
        const popContainer = document.getElementById("popContainer")
        popContainer.style.overflow="visible"
      }
    },
    componentRightClick($event){
      //this.$refs.componentMenu.showContextMenu($event.pageX,$event.pageY)
    },
    componentRightDialogClick($event){
      //this.$refs.componentMenuDialog.showContextMenu($event.pageX,$event.pageY)
    },
    setScreenUndoDialog() {
      const appRef = this.$refs.ismrender1
      let arr = this.getScale();
      arr[0]=1
      arr[1]=1
      appRef.style.transform =
          "scale(" + arr[0] + "," + arr[1] + ")";
    },
    // 设置比例
    setScale() {
      if(this.tempAutoSize==1) {
        this.ISMRuningLoginCavasContainer.resize(window.innerWidth, window.innerHeight);
        const dd = window.innerWidth / this.configData.layer.width;
        const dh = window.innerHeight / this.configData.layer.height;
        this.ISMRuningLoginCavasContainer.scale( dd,  dh);
      }
      else {
        this.ISMRuningLoginCavasContainer.resize(this.configData.layer.width, this.configData.layer.height);
        this.ISMRuningLoginCavasContainer.scale( 1,  1);
      }

    },
    ActionSoundSpeech : function() {
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
    setScaleDialog() {
      const appRef = this.$refs.ismrender1
      let arr = this.getScale();
      appRef.style.transform =
          "scale(" + arr[0] + "," + arr[1] + ")";
    },
    getScale() {
      const w = window.innerWidth / this.configData.layer.width;
      const h = window.innerHeight / this.configData.layer.height;
      return [w, h];
    },
    doComponentAction(component,active){
      let _this = this
      if(active=='click')
      {
        clearTimeout(this.clickTimer);  //首先清除计时器
        this.clickTimer = setTimeout(() => {
          if(typeof component.action!='undefined'&&component.action.length)
          {
            for(let i = 0,actionLen = component.action.length; i <actionLen ; i++) {
              let action = component.action[i];
              if(action.type == active) {
                if(typeof action.actionVoice!='undefined'&&action.actionVoice.length) {
                  this.ActionSoundSpeech().start({
                    container: "#ActionSppekContent",
                    Lang: this.$i18n.locale,
                    rate: 1
                  }, action.actionVoice);
                }
                this.handleComponentAction(component,action);
              }
            }
          }
          else if(typeof component.active!='undefined'&&component.active.length) {
            if (active == "click") {
              for (let i = 0,activeLen = component.active.length; i < activeLen; i++) {
                let componentActive = component.active[i];
                if (componentActive.isSwitch == true) {
                  if(typeof componentActive.condition.actionVoice!='undefined'&&componentActive.condition.actionVoice.length) {
                    this.ActionSoundSpeech().start({
                      container: "#ActionSppekContent",
                      Lang: this.$i18n.locale,
                      rate: 1
                    }, componentActive.condition.actionVoice);
                  }
                  _this.$confirm({
                    title: _this.$t('displayConfig.Properties.SecondConfirm'),
                    content: _this.$t('displayConfig.Properties.SecondConfirmContent'),
                    onOk() {
                      _this.handleComponentSwitchActive(componentActive);
                    },
                    onCancel() {

                    },
                    class: 'test',
                  });

                }
              }
            }
          }
        }, 300);
        return
      }
      else if(active=='dblclick')
      {
        clearTimeout(this.clickTimer);  //首先清除计时器
      }
      if(typeof component.action!='undefined'&&component.action.length)
      {
        for(let i = 0,actionLen = component.action.length; i < actionLen; i++) {
          let action = component.action[i];
          if(action.type == active) {
            if(typeof action.actionVoice!='undefined'&&action.actionVoice.length) {
              this.ActionSoundSpeech().start({
                container: "#ActionSppekContent",
                Lang: this.$i18n.locale,
                rate: 1
              }, action.actionVoice);
            }
            this.handleComponentAction(component,action);
          }
        }
      }
      else if(typeof component.active!='undefined'&&component.active.length) {
        if (active == "click") {
          for (let i = 0; i < component.active.length; i++) {
            let componentActive = component.active[i];
            if (componentActive.isSwitch == true) {
              if(typeof componentActive.condition
                  .actionVoice!='undefined'&&componentActive.condition
                  .actionVoice.length) {
                this.ActionSoundSpeech().start({
                  container: "#ActionSppekContent",
                  Lang: this.$i18n.locale,
                  rate: 1
                }, componentActive.condition.actionVoice);
              }
              _this.$confirm({
                title: _this.$t('displayConfig.Properties.SecondConfirm'),
                content: _this.$t('displayConfig.Properties.SecondConfirmContent'),
                onOk() {
                  _this.handleComponentSwitchActive(componentActive);
                },
                onCancel() {

                },
                class: 'test',
              });
            }
          }
        }
      }
    },
    handleComponentSwitchActive(active){
      let _this = this;
      try {
        if(this.showToken=="") {
          if ((typeof this.user.Role != "undefined") && (this.user.Role != "Admin")) {
            if ((typeof active.condition.actionAuth != "undefined") && (active.condition.actionAuth.length > 0)) {
              let auth = active.condition.actionAuth.find(item => item === this.user.Role)
              if (auth != this.user.Role) {
                this.$message.error(this.$t('displayModel.NoAuthentication'), 3)
                return
              }
            }
          }
        }
      } catch (e) {
        this.$message.error(this.$t('displayModel.NoAuthentication'), 3)
        return
      }
      let setValue = active.condition
      if(setValue.isBandDevice)
      {
        this.deviceUuid =setValue.deviceSN
      }
      else
      {
        this.deviceUuid =this.SelectDeviceUuid
      }
      this.setDataUuid = setValue.dataID
      if(setValue.IsManual)
      {
        this.SetPasswordFormValue=""
        this.settingDialog = true
        this.settingVisible = true
      }
      else
      {
        let findFlag = 0
        if(typeof setValue.AutoSet !="undefined"&&setValue.AutoSet.length<=0)
        {
          this.$message.error(this.$t("component.public.ParamsConfigError"))
          return
        }
        for(let i =0,AutoSetLen = setValue.AutoSet.length;i<AutoSetLen;i++)
        {
          if(active.result==setValue.AutoSet[i].value)
          {
            if(i==0)
            {
              setValue.AutoSetValue = setValue.AutoSet[1].value
            }
            else
            {
              setValue.AutoSetValue = setValue.AutoSet[0].value
            }
            findFlag=1;
            this.AutoSetValue = setValue.AutoSetValue
            if((typeof setValue.SetPassword=="undefined")||(setValue.SetPassword==""))
            {
              this.AutoSetData(setValue.AutoSetValue)
            }
            else
            {
              this.SetPassword = setValue.SetPassword
              this.SetAutoPasswordFormValue=""
              this.setPasswordDialog = true
            }

            break
          }
        }
        if(!findFlag)
        {
          setValue.AutoSetValue = setValue.AutoSet[0].value
          this.AutoSetValue = setValue.AutoSetValue
          if((typeof setValue.SetPassword=="undefined")||(setValue.SetPassword==""))
          {
            this.AutoSetData(setValue.AutoSetValue)
          }
          else
          {
            this.SetPassword = setValue.SetPassword
            this.SetAutoPasswordFormValue=""
            this.setPasswordDialog = true
          }
        }
      }
    },
    ClosePopDialog(){
      this.PopUpDialog  = false
      if(this.ISMPopUpRunningContainer!=null)
      {
        this.ISMPopUpRunningContainer.dispose();
        this.ISMPopUpRunningContainer = null
      }
    },
    PopUpDialogClick(){

      if(this.IsAutoClose)
      {
        this.PopUpDialog  = false
        if(this.ISMPopUpRunningContainer!=null)
        {
          this.ISMPopUpRunningContainer.dispose();
          this.ISMPopUpRunningContainer = null
        }
      }
    },
    doLoopSetValue(counter,setvalue,delay){
      let _this = this

      let setValue = setvalue[counter]
      if((typeof setValue=="undefined"))
      {
        return
      }
      if((typeof setValue.isBandDevice!="undefined")&&(setValue.isBandDevice))
      {
        this.deviceUuid =setValue.deviceSN
      }
      else
      {
        this.deviceUuid =this.SelectDeviceUuid
      }
      this.setDataUuid = setValue.dataID
      this.SetPassword = setValue.SetPassword?setValue.SetPassword:""
      if(setValue.IsManual)
      {
        this.SetPasswordFormValue=""
        this.settingDialog = true
        this.settingVisible = true
      }
      else
      {
        this.AutoSetValue = setValue.AutoSetValue
        if((typeof setValue.SetPassword=="undefined")||(setValue.SetPassword==""))
        {
          _this.AutoSetData(setValue.AutoSetValue)
        }
        else
        {
          this.SetAutoPasswordFormValue=""
          this.setPasswordDialog = true
        }
      }


      if (counter < setvalue.length) {
        setTimeout(function() {
          _this.doLoopSetValue(counter + 1,setvalue,delay);
        }, delay);
      }
    },
    handleComponentAction(component,action){

      if(typeof action.action=="undefined")
      {
        return
      }
      let _this = this;
      try {
        if(this.showToken=="") {
          if ((typeof this.user.Role != "undefined") && (this.user.Role != "Admin")) {
            if ((typeof action.actionAuth != "undefined") && (action.actionAuth.length > 0)) {
              let auth = action.actionAuth.find(item => item === this.user.Role)
              if (auth != this.user.Role) {
                this.$message.error(this.$t('displayModel.NoAuthentication'), 3)
                return
              }
            }
          }
        }
      } catch (e) {
        this.$message.error(this.$t('displayModel.NoAuthentication'), 3)
        return
      }

      if((typeof action.ActionPassword=="undefined")||(action.ActionPassword==""))
      {
        if(action.action == 'visible'){
          if(action.showItems.length > 0) {
            action.showItems.forEach(identifier => {
              _this.showComponent(identifier,true);
            });
          }
          if(action.hideItems.length > 0) {
            action.hideItems.forEach(identifier => {
              _this.showComponent(identifier,false);
            });
          }
        }
        else if(action.action == 'RestApi') {
          let jsonT = `<div>${action.RestApi.Params}</div>`
          let params = TemplateRender.render(jsonT,{data: {msg: 'hello'}}).outerHTML
          let divReg = new RegExp("<div>","g")
          let div2Reg = new RegExp("</div>","g")
          params = params.replace(divReg,"")
          params = params.replace(div2Reg,"")
          if(typeof action.actionConfirm !== 'undefined' && action.actionConfirm)
          {
            _this.$confirm({
              title: _this.$t('displayConfig.Properties.SecondConfirm'),
              content: _this.$t('displayConfig.Properties.SecondConfirmContent'),
              onOk() {
                ComponentRestApi(action.RestApi.Type, action.RestApi.Url, action.RestApi.Params).then(function (res) {

                }).catch(function () {

                })
              },
              onCancel() {

              },
              class: 'test',
            });

          }
          else
          {
            ComponentRestApi(action.RestApi.Type, action.RestApi.Url, action.RestApi.Params).then(function (res) {

            }).catch(function () {

            })
          }
        }
        else if(action.action == 'link') {
          if(typeof action.actionConfirm !== 'undefined' && action.actionConfirm) {
            _this.$confirm({
              title: _this.$t('displayConfig.Properties.SecondConfirm'),
              content: _this.$t('displayConfig.Properties.SecondConfirmContent'),
              onOk() {
                _this.showPage(action.link);
              },
              onCancel() {

              },
              class: 'test',
            });
          }
          else
          {
            _this.showPage(action.link);
          }
        }
        else if(action.action == 'SetValue') {
          if(typeof action.SetDelay == 'undefined') {
            action.SetDelay=1000
          }else if(action.SetDelay == ""){
            action.SetDelay=1000
          }
          if(typeof action.actionConfirm !== 'undefined' && action.actionConfirm) {
            _this.$confirm({
              title: _this.$t('displayConfig.Properties.SecondConfirm'),
              content: _this.$t('displayConfig.Properties.SecondConfirmContent'),
              onOk() {
                _this.doLoopSetValue(0,action.setValue,action.SetDelay)
              },
              onCancel() {

              },
              class: 'test',
            });
          }
          else
          {
            _this.doLoopSetValue(0,action.setValue,action.SetDelay)
          }
        }
        else if(action.action == 'DeviceView'){
          const onSelectData = {
            key:action.DeviceView.key,
            showUUID:action.DeviceView.showUUID,
            showPageUUID:action.DeviceView.showPageUUID,
            type:action.DeviceView.type,
            isPopUp:action.DeviceView.isPopUp?action.DeviceView.isPopUp:null,
            selectKey:action.DeviceView.key,
            isContainer:action.DeviceView.isContainer,
          }
          this.PopUpDialog = false
          if(typeof action.actionConfirm !== 'undefined' && action.actionConfirm) {
            _this.$confirm({
              title: _this.$t('displayConfig.Properties.SecondConfirm'),
              content: _this.$t('displayConfig.Properties.SecondConfirmContent'),
              onOk() {
                if(onSelectData.isContainer!="undefined"&&onSelectData.isContainer)
                {
                  _this.$EventBus.$emit("onContainerSelectDevice", onSelectData);
                }
                else
                {
                  _this.$EventBus.$emit("onSelectDevice", onSelectData);
                }
              },
              onCancel() {

              },
              class: 'test',
            });
          }
          else
          {
            if(onSelectData.isContainer!="undefined"&&onSelectData.isContainer)
            {
              _this.$EventBus.$emit("onContainerSelectDevice", onSelectData);
            }
            else
            {
              _this.$EventBus.$emit("onSelectDevice", onSelectData);
            }
          }

        }
        else if(action.action == 'SysScript') {
          if(typeof action.actionConfirm !== 'undefined' && action.actionConfirm) {
            _this.$confirm({
              title: _this.$t('displayConfig.Properties.SecondConfirm'),
              content: _this.$t('displayConfig.Properties.SecondConfirmContent'),
              onOk() {
                _this.doSysScript(action.ScriptList);
              },
              onCancel() {

              },
              class: 'test',
            });
          }
          else
          {
            _this.doSysScript(action.ScriptList);
          }

        }
        else if(action.action == 'Animation') {
          if(action.animationStatus=="start")
          {
            this.$EventBus.$emit(component.identifier+"animateEvent", 1);
          }
          else
          {
            this.$EventBus.$emit(component.identifier+"animateEvent", 0);
          }
        }
      }
      else
      {
        this.ActionComponent = component
        this.ActionEvent = action
        this.ActionPasswordSet = action.ActionPassword
        this.ActionPasswordValue=""
        this.actionPasswordDialog = true
      }
    },
    doSysScript(Script){
      let _t =this
      let params = {
        Script:Script
      };
      ExecSysScript(params).then(function (res){
        _t.settingLoading=false
        if(res.data.code==0)
        {
          _t.$message.success(_t.$t(res.data.msg))
          _t.settingVisible = false
        }
        else
        {
          _t.$message.error(_t.$t(res.data.msg))
        }
      }).catch(function (error) {
        _t.settingLoading = false
        _t.$message.error(_t.$t("readData.SetFailed"))
      }).finally(function (error) {
        _t.settingLoading = false
      })
    },
    ManualSetData(){
      let _t = this
      this.showKeyboard=false
      if((this.SetPassword!=""))
      {
        if(this.SetPasswordFormValue!=this.SetPassword)
        {
          _t.$message.error(_t.$t("readData.SetPasswordError"))
          return
        }
      }

      let params = {
        deviceUuid:_t.deviceUuid,
        dataUuid:_t.setDataUuid,
        value:_t.SetValueFormValue,
      };
      if(_t.IsConfirm)
      {
        _t.$confirm({
          title: _t.$t('displayConfig.Properties.SecondConfirm'),
          content: _t.$t('displayConfig.Properties.SecondConfirmContent'),
          onOk() {
            _t.settingLoading=true
            setData(params).then(function (res) {
              _t.settingLoading = false
              if (res.data.code == 0) {
                _t.$message.success(_t.$t("readData.SetSuccess"))
                _t.settingVisible = false
                _t.settingDialog = false
              } else {
                _t.$message.error(_t.$t("readData.SetFailed"))
              }
            }).catch(function (error) {
              _t.settingLoading = false
              _t.$message.error(_t.$t("readData.SetFailed"))
            }).finally(function (error) {
              _t.settingLoading = false
            })
          },
          onCancel() {

          },
          class: 'test',
        });
      }
      else {
        _t.settingLoading=true
        setData(params).then(function (res) {
          _t.settingLoading = false
          if (res.data.code == 0) {
            _t.$message.success(_t.$t("readData.SetSuccess"))
            _t.settingVisible = false
            _t.settingDialog = false
          } else {
            _t.$message.error(_t.$t("readData.SetFailed"))
          }
        }).catch(function (error) {
          _t.settingLoading = false
          _t.$message.error(_t.$t("readData.SetFailed"))
        }).finally(function (error) {
          _t.settingLoading = false
        })
      }
    },
    PasswordSetAction(){
      let _this = this
      if(this.ActionPasswordSet==this.ActionPasswordValue)
      {
        let action = this.ActionEvent
        let component = this.ActionComponent
        this.actionPasswordDialog = false
        this.ActionEvent=null
        this.ActionComponent=null
        if(action.action == 'visible'){
          if(action.showItems.length > 0) {
            action.showItems.forEach(identifier => {
              _this.showComponent(identifier,true);
            });
          }
          if(action.hideItems.length > 0) {
            action.hideItems.forEach(identifier => {
              _this.showComponent(identifier,false);
            });
          }
        }
        else if(action.action == 'RestApi') {
          let jsonT = `<div>${action.RestApi.Params}</div>`
          let params = TemplateRender.render(jsonT,{data: {msg: 'hello'}}).outerHTML
          let divReg = new RegExp("<div>","g")
          let div2Reg = new RegExp("</div>","g")
          params = params.replace(divReg,"")
          params = params.replace(div2Reg,"")
          if(typeof action.actionConfirm !== 'undefined' && action.actionConfirm)
          {
            _this.$confirm({
              title: _this.$t('displayConfig.Properties.SecondConfirm'),
              content: _this.$t('displayConfig.Properties.SecondConfirmContent'),
              onOk() {
                ComponentRestApi(action.RestApi.Type, action.RestApi.Url, action.RestApi.Params).then(function (res) {

                }).catch(function () {

                })
              },
              onCancel() {

              },
              class: 'test',
            });

          }
          else
          {
            ComponentRestApi(action.RestApi.Type, action.RestApi.Url, action.RestApi.Params).then(function (res) {

            }).catch(function () {

            })
          }
        }
        else if(action.action == 'link') {
          if(typeof action.actionConfirm !== 'undefined' && action.actionConfirm) {
            _this.$confirm({
              title: _this.$t('displayConfig.Properties.SecondConfirm'),
              content: _this.$t('displayConfig.Properties.SecondConfirmContent'),
              onOk() {
                _this.showPage(action.link);
              },
              onCancel() {

              },
              class: 'test',
            });
          }
          else
          {
            _this.showPage(action.link);
          }
        }
        else if(action.action == 'SetValue') {
          for(let i = 0,setValueLen = action.setValue.length;i<setValueLen;i++)
          {
            let setValue = action.setValue[i]
            if(setValue.isBandDevice)
            {
              this.deviceUuid =setValue.deviceSN
            }
            else
            {
              this.deviceUuid =this.SelectDeviceUuid
            }
            if(typeof action.actionConfirm !== 'undefined' && action.actionConfirm)
            {
              this.IsConfirm = true
            }
            else
            {
              this.IsConfirm = false
            }
            this.setDataUuid = setValue.dataID
            this.SetPassword = setValue.SetPassword?setValue.SetPassword:""
            if(setValue.IsManual)
            {
              this.SetPasswordFormValue=""
              this.settingDialog = true
              this.settingVisible = true
            }
            else
            {
              this.AutoSetValue = setValue.AutoSetValue
              if((typeof setValue.SetPassword=="undefined")||(setValue.SetPassword==""))
              {
                if(typeof action.actionConfirm !== 'undefined' && action.actionConfirm) {
                  _this.$confirm({
                    title: _this.$t('displayConfig.Properties.SecondConfirm'),
                    content: _this.$t('displayConfig.Properties.SecondConfirmContent'),
                    onOk() {
                      _this.AutoSetData(setValue.AutoSetValue)
                    },
                    onCancel() {

                    },
                    class: 'test',
                  });
                }
                else
                {
                  _this.AutoSetData(setValue.AutoSetValue)
                }

              }
              else
              {
                this.SetAutoPasswordFormValue=""
                this.setPasswordDialog = true
              }
            }
          }
        }
        else if(action.action == 'DeviceView'){
          const onSelectData = {
            key:action.DeviceView.key,
            showUUID:action.DeviceView.showUUID,
            showPageUUID:action.DeviceView.showPageUUID,
            type:action.DeviceView.type,
            isPopUp:action.DeviceView.isPopUp?action.DeviceView.isPopUp:null,
            selectKey:action.DeviceView.key,
            isContainer:action.DeviceView.isContainer,
          }
          this.PopUpDialog = false
          if(typeof action.actionConfirm !== 'undefined' && action.actionConfirm) {
            _this.$confirm({
              title: _this.$t('displayConfig.Properties.SecondConfirm'),
              content: _this.$t('displayConfig.Properties.SecondConfirmContent'),
              onOk() {
                if(onSelectData.isContainer!="undefined"&&onSelectData.isContainer)
                {
                  _this.$EventBus.$emit("onContainerSelectDevice", onSelectData);
                }
                else
                {
                  _this.$EventBus.$emit("onSelectDevice", onSelectData);
                }
              },
              onCancel() {

              },
              class: 'test',
            });
          }
          else
          {
            if(onSelectData.isContainer!="undefined"&&onSelectData.isContainer)
            {
              _this.$EventBus.$emit("onContainerSelectDevice", onSelectData);
            }
            else
            {
              _this.$EventBus.$emit("onSelectDevice", onSelectData);
            }
          }

        }
        else if(action.action == 'SysScript') {
          if(typeof action.actionConfirm !== 'undefined' && action.actionConfirm) {
            _this.$confirm({
              title: _this.$t('displayConfig.Properties.SecondConfirm'),
              content: _this.$t('displayConfig.Properties.SecondConfirmContent'),
              onOk() {
                _this.doSysScript(action.ScriptList);
              },
              onCancel() {

              },
              class: 'test',
            });
          }
          else
          {
            _this.doSysScript(action.ScriptList);
          }

        }
        else if(action.action == 'Animation') {
          if(action.animationStatus=="start")
          {
            this.$EventBus.$emit(component.identifier+"animateEvent", 1);
          }
          else
          {
            this.$EventBus.$emit(component.identifier+"animateEvent", 0);
          }
        }
      }
      else
      {
        this.$message.error(this.$t('displayModel.NoAuthentication'), 3)
      }
    },
    PasswordSetData(){
      let _t = this
      this.showKeyboard=false
      if(this.SetAutoPasswordFormValue==this.SetPassword)
      {
        let params = {
          deviceUuid:_t.deviceUuid,
          dataUuid:_t.setDataUuid,
          value:this.AutoSetValue.toString(),
        };
        if(_t.IsConfirm)
        {
          _t.$confirm({
            title: _t.$t('displayConfig.Properties.SecondConfirm'),
            content: _t.$t('displayConfig.Properties.SecondConfirmContent'),
            onOk() {
              _t.settingLoading=true
              setData(params).then(function (res) {
                _t.settingLoading = false
                if (res.data.code == 0) {
                  _t.$message.success(_t.$t("readData.SetSuccess"))
                  _t.settingVisible = false
                  _t.setPasswordDialog = false
                } else {
                  _t.$message.error(_t.$t("readData.SetFailed"))
                }
              }).catch(function (error) {
                _t.settingLoading = false
                _t.$message.error(_t.$t("readData.SetFailed"))
              }).finally(function (error) {
                _t.settingLoading = false
              })
            },
            onCancel() {

            },
            class: 'test',
          });
        }
        else {
          _t.settingLoading=true
          setData(params).then(function (res) {
            _t.settingLoading = false
            if (res.data.code == 0) {
              _t.$message.success(_t.$t("readData.SetSuccess"))
              _t.settingVisible = false
              _t.setPasswordDialog = false
            } else {
              _t.$message.error(_t.$t("readData.SetFailed"))
            }
          }).catch(function (error) {
            _t.settingLoading = false
            _t.$message.error(_t.$t("readData.SetFailed"))
          }).finally(function (error) {
            _t.settingLoading = false
          })
        }
      }
      else
      {
        _t.$message.error(_t.$t("readData.SetPasswordError"))
      }
    },
    AutoSetData(value){
      let _t = this
      this.settingLoading=true
      let params = {
        deviceUuid:_t.deviceUuid,
        dataUuid:_t.setDataUuid,
        value:value.toString()
      };
      _t.$message.loading({ content: 'Loading...' })
      setData(params).then(function (res){
        _t.settingLoading=false
        if(res.data.code==0)
        {
          _t.$message.success(_t.$t("readData.SetSuccess"))

          _t.settingVisible = false
        }
        else
        {
          _t.$message.error(_t.$t("readData.SetFailed"))
        }

      }).catch(function (error) {
        _t.settingLoading = false
        _t.$message.error(_t.$t("readData.SetFailed"))
      }).finally(function (error) {
        _t.settingLoading = false
        setTimeout(function (){
          _t.$message.destroy()
        },500)
      })
    },
    DealWithUpdateData(realData){
      let _t = this

      for (let k = 0,Datalen = realData.Data.length; k < Datalen; k++) {
        if((_t.CurrentRealUUIDList.indexOf(realData.Data[k].Uuid)==-1)&&(_t.CurrentModelUUIDList.indexOf(realData.Data[k].ModelDataUuid)==-1))
        {
          continue
        }
        for(let j = 0,componentsLen = _t.configData.components.cells.length;j<componentsLen;j++)
        {
          if((typeof _t.configData.components.cells[j].data.detail.animate!="undefined")&&(typeof _t.configData.components.cells[j].data.detail.animate.condition!="undefined"))
          {
            let  condition = _t.configData.components.cells[j].data.detail.animate.condition
            let  selectAnimate = _t.configData.components.cells[j].data.detail.animate.selected
            let isExpression = _t.configData.components.cells[j].data.detail.animate.isExpression
            //动画
            if(condition.isBandDevice)
            {
              let isStart = false
              if(realData.DeviceUuid==condition.deviceSN) {

                if ((condition.dataID == realData.Data[k].ModelDataUuid)||(condition.dataID == realData.Data[k].Uuid)) {
                  const RealValue = parseFloat(realData.Data[k].Value)
                  if (isExpression) {
                    const OperatorValue = parseFloat(condition.OperatorValue)
                    switch (condition.operator) {
                      case "==": {
                        if (RealValue == OperatorValue) {
                          isStart = true
                        } else {
                          isStart = false
                        }
                        break
                      }
                      case ">": {
                        if (RealValue > OperatorValue) {
                          isStart = true
                        } else {
                          isStart = false
                        }
                        break
                      }
                      case ">=": {
                        if (RealValue >= OperatorValue) {
                          isStart = true
                        } else {
                          isStart = false
                        }
                        break
                      }
                      case "<": {
                        if (RealValue < OperatorValue) {
                          isStart = true
                        } else {
                          isStart = false
                        }
                        break
                      }
                      case "<=": {
                        if (RealValue <= OperatorValue) {
                          isStart = true
                        } else {
                          isStart = false
                        }
                        break
                      }
                      case "!=": {
                        if (RealValue != OperatorValue) {
                          isStart = true
                        } else {
                          isStart = false
                        }
                        break
                      }
                      case "<>": {
                        const OperatorMaxValue = parseFloat(condition.OperatorMaxValue)
                        if ((RealValue >= OperatorValue) && (RealValue <= OperatorMaxValue)) {
                          isStart = true
                        } else {
                          isStart = false
                        }
                        break
                      }
                      case "<!>": {
                        const OperatorMaxValue = parseFloat(condition.OperatorMaxValue)
                        if ((RealValue <= OperatorValue) || (RealValue >= OperatorMaxValue)) {
                          isStart = true
                        } else {
                          isStart = false
                        }
                        break
                      }
                    }
                    if (selectAnimate.includes("visible")) {
                      if (isStart) {
                        _t.configData.components.cells[j].data.detail.style.visible = true
                      } else {
                        _t.configData.components.cells[j].data.detail.style.visible = false
                      }
                    }
                    this.$EventBus.$emit(_t.configData.components.cells[j].data.detail.identifier + "animateEvent", isStart);
                  }
                }
              }
            }
            else if(realData.DeviceUuid==_t.SelectDeviceUuid)
            {
              let isStart = false
              if ((condition.dataID == realData.Data[k].ModelDataUuid)||(condition.dataID == realData.Data[k].Uuid)) {
                const RealValue = parseFloat(realData.Data[k].Value)
                if (isExpression) {
                  const OperatorValue = parseFloat(condition.OperatorValue)
                  switch (condition.operator) {
                    case "==": {
                      if (RealValue == OperatorValue) {
                        isStart = true
                      } else {
                        isStart = false
                      }
                      break
                    }
                    case ">": {
                      if (RealValue > OperatorValue) {
                        isStart = true
                      } else {
                        isStart = false
                      }
                      break
                    }
                    case ">=": {
                      if (RealValue >= OperatorValue) {
                        isStart = true
                      } else {
                        isStart = false
                      }
                      break
                    }
                    case "<": {
                      if (RealValue < OperatorValue) {
                        isStart = true
                      } else {
                        isStart = false
                      }
                      break
                    }
                    case "<=": {
                      if (RealValue <= OperatorValue) {
                        isStart = true
                      } else {
                        isStart = false
                      }
                      break
                    }
                    case "!=": {
                      if (RealValue != OperatorValue) {
                        isStart = true
                      } else {
                        isStart = false
                      }
                      break
                    }
                    case "<>": {
                      const OperatorMaxValue = parseFloat(condition.OperatorMaxValue)
                      if ((RealValue >= OperatorValue) && (RealValue <= OperatorMaxValue)) {
                        isStart = true
                      } else {
                        isStart = false
                      }
                      break
                    }
                    case "<!>": {
                      const OperatorMaxValue = parseFloat(condition.OperatorMaxValue)
                      if ((RealValue <= OperatorValue) || (RealValue >= OperatorMaxValue)) {
                        isStart = true
                      } else {
                        isStart = false
                      }
                      break
                    }
                  }
                  if (selectAnimate.includes("visible")) {
                    if (isStart) {
                      _t.configData.components.cells[j].data.detail.style.visible = true
                    } else {
                      _t.configData.components.cells[j].data.detail.style.visible = false
                    }
                  }
                  this.$EventBus.$emit(_t.configData.components.cells[j].data.detail.identifier + "animateEvent", isStart);
                }
              }
            }
          }
          //位移动画
          if((typeof _t.configData.components.cells[j].data.detail.animate!="undefined")&&(typeof _t.configData.components.cells[j].data.detail.animate.move!="undefined")&&_t.configData.components.cells[j].data.detail.animate.selected.includes("animateMove"))
          {
            let  conditionx = _t.configData.components.cells[j].data.detail.animate.move.x
            let  conditiony = _t.configData.components.cells[j].data.detail.animate.move.y
            //动画
            if(conditionx.isBandDevice)
            {
              if(realData.DeviceUuid==conditionx.deviceSN) {

                if ((conditionx.dataID == realData.Data[k].ModelDataUuid)||(conditionx.dataID == realData.Data[k].Uuid)) {

                  const RealValue = parseFloat(realData.Data[k].Value)
                  let animateMove = {
                    ID: "animateMoveX",
                    DeviceSN: conditionx.deviceSN,
                    dataID: conditionx.dataID,
                    result: RealValue
                  }
                  _t.configData.components.cells[j].data.detail.style.position.x = RealValue
                  this.$EventBus.$emit(_t.configData.components.cells[j].data.detail.identifier + "animateMove", animateMove);
                }
              }
            }
            else if(realData.DeviceUuid==_t.SelectDeviceUuid)
            {
              if ((conditionx.dataID == realData.Data[k].ModelDataUuid)||(conditionx.dataID == realData.Data[k].Uuid)) {
                const RealValue = parseFloat(realData.Data[k].Value)
                let animateMove = {
                  ID: "animateMoveX",
                  DeviceSN: conditionx.deviceSN,
                  dataID: conditionx.dataID,
                  result: RealValue
                }
                _t.configData.components.cells[j].data.detail.style.position.x = RealValue
                this.$EventBus.$emit(_t.configData.components.cells[j].data.detail.identifier + "animateMove", animateMove);
              }
            }

            //动画
            if(conditiony.isBandDevice)
            {
              if(realData.DeviceUuid==conditiony.deviceSN) {

                if ((conditiony.dataID == realData.Data[k].ModelDataUuid)||(conditiony.dataID == realData.Data[k].Uuid)) {

                  const RealValue = parseFloat(realData.Data[k].Value)
                  let animateMove = {
                    ID: "animateMoveY",
                    DeviceSN: conditiony.deviceSN,
                    dataID: conditiony.dataID,
                    result: RealValue
                  }
                  _t.configData.components.cells[j].data.detail.style.position.y = RealValue
                  this.$EventBus.$emit(_t.configData.components.cells[j].data.detail.identifier + "animateMove", animateMove);
                }
              }
            }
            else if(realData.DeviceUuid==_t.SelectDeviceUuid)
            {
              if ((conditiony.dataID == realData.Data[k].ModelDataUuid)||(conditiony.dataID == realData.Data[k].Uuid)) {
                const RealValue = parseFloat(realData.Data[k].Value)
                let animateMove = {
                  ID: "animateMoveY",
                  DeviceSN: conditiony.deviceSN,
                  dataID: conditiony.dataID,
                  result: RealValue
                }
                _t.configData.components.cells[j].data.detail.style.position.y = RealValue
                this.$EventBus.$emit(_t.configData.components.cells[j].data.detail.identifier + "animateMove", animateMove);
              }
            }
          }

          if((typeof _t.configData.components.cells[j].data.detail.active!="undefined"))
          {
            let  active = _t.configData.components.cells[j].data.detail.active
            //动作
            let tempResult = ""
            for (let activeIndex = 0; activeIndex < active.length; activeIndex++) {
              if(active[activeIndex].condition.isBandDevice)
              {
                if(realData.DeviceUuid==active[activeIndex].condition.deviceSN) {
                  if ((active[activeIndex].condition.dataID == realData.Data[k].ModelDataUuid)||(active[activeIndex].condition.dataID == realData.Data[k].Uuid)) {
                    if ((typeof active[activeIndex].isExpression != 'undefined') && (active[activeIndex].isExpression)) {
                      const RealValue = parseFloat(realData.Data[k].Value)
                      const OperatorValue = parseFloat(active[activeIndex].condition.OperatorValue)
                      switch (active[activeIndex].condition.operator) {
                        case "==": {
                          if (RealValue == OperatorValue) {
                            tempResult = true
                          } else {
                            tempResult = false
                          }
                          break
                        }
                        case ">": {
                          if (RealValue > OperatorValue) {
                            tempResult = true
                          } else {
                            tempResult = false
                          }
                          break
                        }
                        case ">=": {
                          if (RealValue >= OperatorValue) {
                            tempResult = true
                          } else {
                            tempResult = false
                          }
                          break
                        }
                        case "<": {
                          if (RealValue < OperatorValue) {
                            tempResult = true
                          } else {
                            tempResult = false
                          }
                          break
                        }
                        case "<=": {
                          if (RealValue <= OperatorValue) {
                            tempResult = true
                          } else {
                            tempResult = false
                          }
                          break
                        }
                        case "!=": {
                          if (RealValue != OperatorValue) {
                            tempResult = true
                          } else {
                            tempResult = false
                          }
                          break
                        }
                        case "<>": {
                          const OperatorMaxValue = parseFloat(active[activeIndex].condition.OperatorMaxValue)
                          if ((RealValue >= OperatorValue) && (RealValue <= OperatorMaxValue)) {
                            tempResult = true
                          } else {
                            tempResult = false
                          }
                          break
                        }
                        case "<!>": {
                          const OperatorMaxValue = parseFloat(active[activeIndex].condition.OperatorMaxValue)
                          if ((RealValue <= OperatorValue) || (RealValue >= OperatorMaxValue)) {
                            tempResult = true
                          } else {
                            tempResult = false
                          }
                          break
                        }
                      }
                    } else {
                      tempResult = realData.Data[k].Value
                    }

                    active[activeIndex].result = tempResult
                    let activeData = {
                      ID: active[activeIndex].id,
                      DeviceSN: active[activeIndex].condition.deviceSN,
                      dataID: active[activeIndex].condition.dataID,
                      index: activeIndex,
                      result: tempResult
                    }
                    const cell = _t.ISMRuningLoginCavasContainer.getCellById(_t.configData.components.cells[j].id)
                    const isEdge = cell.isEdge()
                    if(!isEdge) {
                      _t.$EventBus.$emit(_t.configData.components.cells[j].data.detail.identifier + "activeEvent", activeData);
                    }else{
                      let animation = ''
                      if(activeData.ID == "Forward"&&activeData.result)
                      {
                        animation = 'ant-line-forward 30s infinite linear'
                      }
                      if(activeData.ID == "Reverse"&&activeData.result)
                      {
                        animation = 'ant-line-inverse 30s infinite linear'
                      }
                      cell.setAttrs({
                        line:{
                          style:{
                            animation:animation,
                          }
                        }
                      });
                    }
                  }
                }
              }
              else if(realData.DeviceUuid==_t.SelectDeviceUuid)
              {
                if ((active[activeIndex].condition.dataID == realData.Data[k].ModelDataUuid)||(active[activeIndex].condition.dataID == realData.Data[k].Uuid)) {
                  if((typeof active[activeIndex].isExpression!='undefined')&&(active[activeIndex].isExpression))
                  {
                    const RealValue = parseFloat(realData.Data[k].Value)
                    const OperatorValue = parseFloat(active[activeIndex].condition.OperatorValue)
                    switch (active[activeIndex].condition.operator){
                      case "==":{
                        if(RealValue==OperatorValue)
                        {
                          tempResult = true
                        }
                        else{
                          tempResult = false
                        }
                        break
                      }
                      case ">":{
                        if(RealValue>OperatorValue)
                        {
                          tempResult = true
                        }
                        else{
                          tempResult = false
                        }
                        break
                      }
                      case ">=":{
                        if(RealValue>=OperatorValue)
                        {
                          tempResult = true
                        }
                        else{
                          tempResult = false
                        }
                        break
                      }
                      case "<":{
                        if(RealValue<OperatorValue)
                        {
                          tempResult = true
                        }
                        else{
                          tempResult = false
                        }
                        break
                      }
                      case "<=":{
                        if(RealValue<=OperatorValue)
                        {
                          tempResult = true
                        }
                        else{
                          tempResult = false
                        }
                        break
                      }
                      case "!=":{
                        if(RealValue!=OperatorValue)
                        {
                          tempResult = true
                        }
                        else{
                          tempResult = false
                        }
                        break
                      }
                      case "<>":{
                        const OperatorMaxValue = parseFloat(active[activeIndex].condition.OperatorMaxValue)
                        if((RealValue>=OperatorValue)&&(RealValue<=OperatorMaxValue))
                        {
                          tempResult = true
                        }
                        else{
                          tempResult = false
                        }
                        break
                      }
                      case "<!>":{
                        const OperatorMaxValue = parseFloat(active[activeIndex].condition.OperatorMaxValue)
                        if((RealValue<=OperatorValue)||(RealValue>=OperatorMaxValue))
                        {
                          tempResult = true
                        }
                        else{
                          tempResult = false
                        }
                        break
                      }
                    }
                  }
                  else
                  {
                    tempResult = realData.Data[k].Value
                  }
                  active[activeIndex].result = tempResult
                  let activeData = {
                    ID:active[activeIndex].id,
                    DeviceSN:active[activeIndex].condition.deviceSN,
                    dataID:active[activeIndex].condition.dataID,
                    index:activeIndex,
                    result:tempResult
                  }
                  const cell = _t.ISMRuningLoginCavasContainer.getCellById(_t.configData.components.cells[j].id)
                  const isEdge = cell.isEdge()
                  if(!isEdge) {
                    _t.$EventBus.$emit(_t.configData.components.cells[j].data.detail.identifier + "activeEvent", activeData);
                  }else{
                    let animation = ''
                    if(activeData.ID == "Forward"&&activeData.result)
                    {
                      animation = 'ant-line-forward 30s infinite linear'
                    }
                    if(activeData.ID == "Reverse"&&activeData.result)
                    {
                      animation = 'ant-line-inverse 30s infinite linear'
                    }
                    cell.setAttrs({
                      line:{
                        style:{
                          animation:animation,
                        }
                      }
                    });
                  }
                }
              }
            }
          }
        }
      }

      if((!_t.chargePagePopUp)&&(_t.PopUpDialog))
      {
        for (let k = 0,realDataLen = realData.Data.length; k < realDataLen; k++) {
          if((_t.CurrentPopRealUUIDList.indexOf(realData.Data[k].Uuid)==-1)&&(_t.CurrentPopModelUUIDList.indexOf(realData.Data[k].ModelDataUuid)==-1))
          {
            continue
          }
          for(let j = 0,componentsLen = _t.PopUpConfigData.components.cells.length;j<componentsLen;j++)
          {
            if((typeof _t.PopUpConfigData.components.cells[j].data.detail.animate!="undefined")&&(typeof _t.PopUpConfigData.components.cells[j].data.detail.animate.condition!="undefined"))
            {
              let  condition = _t.PopUpConfigData.components.cells[j].data.detail.animate.condition
              let  selectAnimate = _t.PopUpConfigData.components.cells[j].data.detail.animate.selected
              let isExpression = _t.PopUpConfigData.components.cells[j].data.detail.animate.isExpression
              //动画
              if(condition.isBandDevice)
              {
                let isStart = false
                if(realData.DeviceUuid==condition.deviceSN) {
                  if ((condition.dataID == realData.Data[k].ModelDataUuid)||(condition.dataID == realData.Data[k].Uuid)) {
                    const RealValue = parseFloat(realData.Data[k].Value)
                    if (isExpression) {
                      const OperatorValue = parseFloat(condition.OperatorValue)
                      switch (condition.operator) {
                        case "==": {
                          if (RealValue == OperatorValue) {
                            isStart = true
                          } else {
                            isStart = false
                          }
                          break
                        }
                        case ">": {
                          if (RealValue > OperatorValue) {
                            isStart = true
                          } else {
                            isStart = false
                          }
                          break
                        }
                        case ">=": {
                          if (RealValue >= OperatorValue) {
                            isStart = true
                          } else {
                            isStart = false
                          }
                          break
                        }
                        case "<": {
                          if (RealValue < OperatorValue) {
                            isStart = true
                          } else {
                            isStart = false
                          }
                          break
                        }
                        case "<=": {
                          if (RealValue <= OperatorValue) {
                            isStart = true
                          } else {
                            isStart = false
                          }
                          break
                        }
                        case "!=": {
                          if (RealValue != OperatorValue) {
                            isStart = true
                          } else {
                            isStart = false
                          }
                          break
                        }
                        case "<>": {
                          const OperatorMaxValue = parseFloat(condition.OperatorMaxValue)
                          if ((RealValue >= OperatorValue) && (RealValue <= OperatorMaxValue)) {
                            isStart = true
                          } else {
                            isStart = false
                          }
                          break
                        }
                        case "<!>": {
                          const OperatorMaxValue = parseFloat(condition.OperatorMaxValue)
                          if ((RealValue <= OperatorValue) || (RealValue >= OperatorMaxValue)) {
                            isStart = true
                          } else {
                            isStart = false
                          }
                          break
                        }
                      }
                      if (selectAnimate.includes("visible")) {
                        if (isStart) {
                          _t.PopUpConfigData.components.cells[j].data.detail.style.visible = true
                        } else {
                          _t.PopUpConfigData.components.cells[j].data.detail.style.visible = false
                        }
                      }
                      _t.$EventBus.$emit(_t.PopUpConfigData.components.cells[j].data.detail.identifier + "animateEvent", isStart);
                    }
                  }
                }
              }
              else if (realData.DeviceUuid==_t.SelectDeviceUuid)
              {
                let isStart = false
                if ((condition.dataID == realData.Data[k].ModelDataUuid)||(condition.dataID == realData.Data[k].Uuid)) {
                  const RealValue = parseFloat(realData.Data[k].Value)
                  if (isExpression) {
                    const OperatorValue = parseFloat(condition.OperatorValue)
                    switch (condition.operator) {
                      case "==": {
                        if (RealValue == OperatorValue) {
                          isStart = true
                        } else {
                          isStart = false
                        }
                        break
                      }
                      case ">": {
                        if (RealValue > OperatorValue) {
                          isStart = true
                        } else {
                          isStart = false
                        }
                        break
                      }
                      case ">=": {
                        if (RealValue >= OperatorValue) {
                          isStart = true
                        } else {
                          isStart = false
                        }
                        break
                      }
                      case "<": {
                        if (RealValue < OperatorValue) {
                          isStart = true
                        } else {
                          isStart = false
                        }
                        break
                      }
                      case "<=": {
                        if (RealValue <= OperatorValue) {
                          isStart = true
                        } else {
                          isStart = false
                        }
                        break
                      }
                      case "!=": {
                        if (RealValue != OperatorValue) {
                          isStart = true
                        } else {
                          isStart = false
                        }
                        break
                      }
                      case "<>": {
                        const OperatorMaxValue = parseFloat(condition.OperatorMaxValue)
                        if ((RealValue >= OperatorValue) && (RealValue <= OperatorMaxValue)) {
                          isStart = true
                        } else {
                          isStart = false
                        }
                        break
                      }
                      case "<!>": {
                        const OperatorMaxValue = parseFloat(condition.OperatorMaxValue)
                        if ((RealValue <= OperatorValue) || (RealValue >= OperatorMaxValue)) {
                          isStart = true
                        } else {
                          isStart = false
                        }
                        break
                      }
                    }
                    if (selectAnimate.includes("visible")) {
                      if (isStart) {
                        _t.PopUpConfigData.components.cells[j].data.detail.style.visible = true
                      } else {
                        _t.PopUpConfigData.components.cells[j].data.detail.style.visible = false
                      }
                    }
                    _t.$EventBus.$emit(_t.PopUpConfigData.components.cells[j].data.detail.identifier + "animateEvent", isStart);
                  }
                }
              }
            }

            //位移动画
            if((typeof _t.PopUpConfigData.components.cells[j].data.detail.animate!="undefined")&&(typeof _t.PopUpConfigData.components.cells[j].data.detail.animate.move!="undefined")&&_t.PopUpConfigData.components.cells[j].data.detail.animate.selected.includes("animateMove"))
            {
              let  conditionx = _t.PopUpConfigData.components.cells[j].data.detail.animate.move.x
              let  conditiony = _t.PopUpConfigData.components.cells[j].data.detail.animate.move.y
              //动画
              if(conditionx.isBandDevice)
              {
                if(realData.DeviceUuid==conditionx.deviceSN) {

                  if ((conditionx.dataID == realData.Data[k].ModelDataUuid)||(conditionx.dataID == realData.Data[k].Uuid)) {

                    const RealValue = parseFloat(realData.Data[k].Value)
                    let animateMove = {
                      ID: "animateMoveX",
                      DeviceSN: conditionx.deviceSN,
                      dataID: conditionx.dataID,
                      result: RealValue
                    }
                    _t.configData.components.cells[j].data.detail.style.position.x = RealValue
                    this.$EventBus.$emit(_t.PopUpConfigData.components.cells[j].data.detail.identifier + "animateMove", animateMove);
                  }
                }
              }
              else if(realData.DeviceUuid==_t.SelectDeviceUuid)
              {
                if ((conditionx.dataID == realData.Data[k].ModelDataUuid)||(conditionx.dataID == realData.Data[k].Uuid)) {
                  const RealValue = parseFloat(realData.Data[k].Value)
                  let animateMove = {
                    ID: "animateMoveX",
                    DeviceSN: conditionx.deviceSN,
                    dataID: conditionx.dataID,
                    result: RealValue
                  }
                  _t.configData.components.cells[j].data.detail.style.position.x = RealValue
                  this.$EventBus.$emit(_t.PopUpConfigData.components.cells[j].data.detail.identifier + "animateMove", animateMove);
                }
              }

              //动画
              if(conditiony.isBandDevice)
              {
                if(realData.DeviceUuid==conditiony.deviceSN) {

                  if ((conditiony.dataID == realData.Data[k].ModelDataUuid)||(conditiony.dataID == realData.Data[k].Uuid)) {

                    const RealValue = parseFloat(realData.Data[k].Value)
                    let animateMove = {
                      ID: "animateMoveY",
                      DeviceSN: conditiony.deviceSN,
                      dataID: conditiony.dataID,
                      result: RealValue
                    }
                    _t.PopUpConfigData.components.cells[j].data.detail.style.position.y = RealValue
                    this.$EventBus.$emit(_t.PopUpConfigData.components.cells[j].data.detail.identifier + "animateMove", animateMove);
                  }
                }
              }
              else if(realData.DeviceUuid==_t.SelectDeviceUuid)
              {
                if ((conditiony.dataID == realData.Data[k].ModelDataUuid)||(conditiony.dataID == realData.Data[k].Uuid)) {
                  const RealValue = parseFloat(realData.Data[k].Value)
                  let animateMove = {
                    ID: "animateMoveY",
                    DeviceSN: conditiony.deviceSN,
                    dataID: conditiony.dataID,
                    result: RealValue
                  }
                  _t.PopUpConfigData.components.cells[j].data.detail.style.position.y = RealValue
                  this.$EventBus.$emit(_t.configData.components.cells[j].data.detail.identifier + "animateMove", animateMove);
                }
              }
            }

            if((typeof _t.PopUpConfigData.components.cells[j].data.detail.active!="undefined"))
            {
              let  active = _t.PopUpConfigData.components.cells[j].data.detail.active
              let tempResult = ""
              for (let activeIndex = 0; activeIndex < active.length; activeIndex++) {
                if(active[activeIndex].condition.isBandDevice)
                {
                  if(realData.DeviceUuid==active[activeIndex].condition.deviceSN) {
                    if ((active[activeIndex].condition.dataID == realData.Data[k].ModelDataUuid)||(active[activeIndex].condition.dataID == realData.Data[k].Uuid)) {
                      if ((typeof active[activeIndex].isExpression != 'undefined') && (active[activeIndex].isExpression)) {
                        const RealValue = parseFloat(realData.Data[k].Value)
                        const OperatorValue = parseFloat(active[activeIndex].condition.OperatorValue)
                        switch (active[activeIndex].condition.operator) {
                          case "==": {
                            if (RealValue == OperatorValue) {
                              tempResult = true
                            } else {
                              tempResult = false
                            }
                            break
                          }
                          case ">": {
                            if (RealValue > OperatorValue) {
                              tempResult = true
                            } else {
                              tempResult = false
                            }
                            break
                          }
                          case ">=": {
                            if (RealValue >= OperatorValue) {
                              tempResult = true
                            } else {
                              tempResult = false
                            }
                            break
                          }
                          case "<": {
                            if (RealValue < OperatorValue) {
                              tempResult = true
                            } else {
                              tempResult = false
                            }
                            break
                          }
                          case "<=": {
                            if (RealValue <= OperatorValue) {
                              tempResult = true
                            } else {
                              tempResult = false
                            }
                            break
                          }
                          case "!=": {
                            if (RealValue != OperatorValue) {
                              tempResult = true
                            } else {
                              tempResult = false
                            }
                            break
                          }
                          case "<>": {
                            const OperatorMaxValue = parseFloat(active[activeIndex].condition.OperatorMaxValue)
                            if ((RealValue >= OperatorValue) && (RealValue <= OperatorMaxValue)) {
                              tempResult = true
                            } else {
                              tempResult = false
                            }
                            break
                          }
                          case "<!>": {
                            const OperatorMaxValue = parseFloat(active[activeIndex].condition.OperatorMaxValue)
                            if ((RealValue <= OperatorValue) || (RealValue >= OperatorMaxValue)) {
                              tempResult = true
                            } else {
                              tempResult = false
                            }
                            break
                          }
                        }
                      } else {
                        tempResult = realData.Data[k].Value
                      }
                      active[activeIndex].result = tempResult
                      let activeData = {
                        ID: active[activeIndex].id,
                        DeviceSN: active[activeIndex].condition.deviceSN,
                        dataID: active[activeIndex].condition.dataID,
                        index: activeIndex,
                        result: tempResult
                      }
                      const cell = _t.ISMRuningLoginCavasContainer.getCellById(_t.configData.components.cells[j].id)
                      const isEdge = cell.isEdge()
                      if(!isEdge) {
                        _t.$EventBus.$emit(_t.PopUpConfigData.components.cells[j].data.detail.identifier + "activeEvent", activeData);
                      }else{
                        let animation = ''
                        if(activeData.ID == "Forward"&&activeData.result)
                        {
                          animation = 'ant-line-forward 30s infinite linear'
                        }
                        if(activeData.ID == "Reverse"&&activeData.result)
                        {
                          animation = 'ant-line-inverse 30s infinite linear'
                        }
                        cell.setAttrs({
                          line:{
                            style:{
                              animation:animation,
                            }
                          }
                        });
                      }
                    }
                  }
                }
                else if (realData.DeviceUuid==_t.SelectDeviceUuid)
                {
                  if ((active[activeIndex].condition.dataID == realData.Data[k].ModelDataUuid)||(active[activeIndex].condition.dataID == realData.Data[k].Uuid)) {
                    if((typeof active[activeIndex].isExpression!='undefined')&&(active[activeIndex].isExpression))
                    {
                      const RealValue = parseFloat(realData.Data[k].Value)
                      const OperatorValue = parseFloat(active[activeIndex].condition.OperatorValue)
                      switch (active[activeIndex].condition.operator){
                        case "==":{
                          if(RealValue==OperatorValue)
                          {
                            tempResult = true
                          }
                          else{
                            tempResult = false
                          }
                          break
                        }
                        case ">":{
                          if(RealValue>OperatorValue)
                          {
                            tempResult = true
                          }
                          else{
                            tempResult = false
                          }
                          break
                        }
                        case ">=":{
                          if(RealValue>=OperatorValue)
                          {
                            tempResult = true
                          }
                          else{
                            tempResult = false
                          }
                          break
                        }
                        case "<":{
                          if(RealValue<OperatorValue)
                          {
                            tempResult = true
                          }
                          else{
                            tempResult = false
                          }
                          break
                        }
                        case "<=":{
                          if(RealValue<=OperatorValue)
                          {
                            tempResult = true
                          }
                          else{
                            tempResult = false
                          }
                          break
                        }
                        case "!=":{
                          if(RealValue!=OperatorValue)
                          {
                            tempResult = true
                          }
                          else{
                            tempResult = false
                          }
                          break
                        }
                        case "<>":{
                          const OperatorMaxValue = parseFloat(active[activeIndex].condition.OperatorMaxValue)
                          if((RealValue>=OperatorValue)&&(RealValue<=OperatorMaxValue))
                          {
                            tempResult = true
                          }
                          else{
                            tempResult = false
                          }
                          break
                        }
                        case "<!>":{
                          const OperatorMaxValue = parseFloat(active[activeIndex].condition.OperatorMaxValue)
                          if((RealValue<=OperatorValue)||(RealValue>=OperatorMaxValue))
                          {
                            tempResult = true
                          }
                          else{
                            tempResult = false
                          }
                          break
                        }
                      }
                    }
                    else
                    {
                      tempResult = realData.Data[k].Value
                    }
                    active[activeIndex].result = tempResult
                    let activeData = {
                      ID:active[activeIndex].id,
                      DeviceSN:active[activeIndex].condition.deviceSN,
                      dataID:active[activeIndex].condition.dataID,
                      index:activeIndex,
                      result:tempResult
                    }
                    const cell = _t.ISMRuningLoginCavasContainer.getCellById(_t.configData.components.cells[j].id)
                    const isEdge = cell.isEdge()
                    if(!isEdge) {
                      _t.$EventBus.$emit(_t.PopUpConfigData.components.cells[j].data.detail.identifier + "activeEvent", activeData);
                    }else{
                      let animation = ''
                      if(activeData.ID == "Forward"&&activeData.result)
                      {
                        animation = 'ant-line-forward 30s infinite linear'
                      }
                      if(activeData.ID == "Reverse"&&activeData.result)
                      {
                        animation = 'ant-line-inverse 30s infinite linear'
                      }
                      cell.setAttrs({
                        line:{
                          style:{
                            animation:animation,
                          }
                        }
                      });
                    }
                  }
                }
              }
            }
          }
        }
      }
      _t.$EventBus.$emit("DealWithRealDataFinish", true);

    },
    showComponent(identifier,visible) {
      let spirits = this.$refs['spirit'];
      for(let i = 0,spiritsLen = spirits.length; i < spiritsLen; i++){
        let spirit = spirits[i];
        if(spirit.detail.groupID == identifier) {
          spirit.detail.style.visible = visible;
        }
      }
    },
    showPage(linkInfo) {
      //  this.PopUpDialog = false
      if(linkInfo && linkInfo.linkType == "Inside" && linkInfo.Inside && (linkInfo.Inside.displayType === 2 || linkInfo.Inside.displayType === '2')) {
        this.$router.push({
          path: `/DisPlay3DRunApp/${linkInfo.Inside.displayUUID}`,
          query: linkInfo.Inside.pageUUID ? { pageId: linkInfo.Inside.pageUUID } : {}
        })
        return
      }
      if(typeof linkInfo.isPopUp!='undefined'&& linkInfo.isPopUp==true)
      {
        let _t = this
        this.IsAutoClose = linkInfo.autoClose
        if (linkInfo.linkType == "Inside") {
          let page = {
            pageType: 1,
            displayUUID: linkInfo.Inside.displayUUID,
            pageUuid: linkInfo.Inside.pageUUID
          }
          _t.isExternUrl = false
          _t.chargePagePopUp = true
          this.CurrentPopRealUUIDList=[]
          this.CurrentPopModelUUIDList=[]
          // this.$message.loading({content: 'Loading...', loadingKey, duration: 0});
          this.selectPopUpDisplayPageDataStruct({
            page: page, callback: function (res,uuids,devices) {
              if(res==0)
              {
                _t.$message.destroy();
                _t.chargePagePopUp = false
                _t.PopUpDialog = true
                _t.$message.loading({content: 'Loading...', loadingKey, duration: 0});
                setTimeout(function (){
                  _t.$message.destroy()
                  _t.RunPopUpCavasContainerInit(_t.PopUpConfigData.layer.autoSize,_t.PopUpConfigData.layer.Padding)
                },500)
                // _t.PopUpDialog.center()
                setTimeout(function (){
                  if(uuids.length>0)
                  {
                    let newuuids=[]
                    for(let i=0;i<uuids.length;i++)
                    {
                      if((typeof uuids[i]!="undefined")&&(uuids[i].length>0))
                      {
                        newuuids.push(uuids[i])
                      }
                    }
                    let newdevices=[]
                    for(let i=0;i<devices.length;i++)
                    {
                      if((typeof devices[i]!="undefined")&&(devices[i].length>0))
                      {
                        newdevices.push(devices[i])
                      }
                    }
                    getRealDataByUuid({uuid:newuuids,devices:newdevices}).then(function (res){
                      _t.settingLoading=false
                      if(res.data.code==0)
                      {
                        for(let k = 0,realDataLen = res.data.realData.length;k<realDataLen;k++)
                        {
                          let pushData = {
                            DeviceUuid:res.data.realData[k].duid,
                            ProjectUuid:res.data.realData[k].project_uuid,
                            Cmd:"RealData",
                            Data:[]
                          }
                          _t.CurrentPopRealUUIDList.push(res.data.realData[k].uuid)
                          _t.CurrentPopModelUUIDList.push(res.data.realData[k].mduid)
                          let DataObj = {
                            Uuid: res.data.realData[k].uuid,
                            ModelDataUuid: res.data.realData[k].mduid,
                            Value: res.data.realData[k].value
                          }
                          pushData.Data.push(DataObj)
                          _t.$EventBus.$emit("readDataPush", pushData);
                        }
                      }
                    })
                  }
                },300)
              }
              else
              {
                _t.$message.error(_t.$t("readData.NotFindPage"))
                setTimeout(function (){
                  _t.$message.destroy();
                },1000)
              }
            }
          })
        }
        else
        {
          _t.chargePagePopUp = false
          _t.isExternUrl = true
          _t.linkInfoExternal.width = linkInfo.width
          _t.linkInfoExternal.height = linkInfo.height
          _t.linkInfoExternal.url = linkInfo.External
          // _t.$refs.PopUpDialog.vcenter()
          _t.PopUpDialog = true
          if(typeof linkInfo.title!='undefined')
          {
            if(linkInfo.title=="")
            {
              this.PopUpConfigData.PageName = "外部网页"
            }
            else
            {
              this.PopUpConfigData.PageName = linkInfo.title
            }
          }
          else {
            this.PopUpConfigData.PageName = "外部网页"
          }
        }
      }
      else {
        if (linkInfo.linkType == "Inside") {
          let _t = this
          let page = {
            pageType: 1,
            displayUUID: linkInfo.Inside.displayUUID,
            pageUuid: linkInfo.Inside.pageUUID
          }

          _t.CurrentRealUUIDList=[]
          _t.CurrentModelUUIDList=[]
          // this.$message.loading({content: 'Loading...', loadingKey, duration: 0});
          this.selectDisplayPageDataStruct({
            page: page, callback: function (uuids,devices,isFound) {
              _t.$message.destroy();
              if (isFound === false) {
                _t.$message.error(_t.$t("readData.NotFindPage"))
                return
              }
              _t.RunCavasContainerInit(_t.configData.layer.autoSize,_t.configData.layer.Padding)
              setTimeout(function (){
                if(uuids.length>0)
                {
                  let newuuids=[]
                  for(let i=0;i<uuids.length;i++)
                  {
                    if((typeof uuids[i]!="undefined")&&(uuids[i].length>0))
                    {
                      newuuids.push(uuids[i])
                    }
                  }
                  let newdevices=[]
                  for(let i=0;i<devices.length;i++)
                  {
                    if((typeof devices[i]!="undefined")&&(devices[i].length>0))
                    {
                      newdevices.push(devices[i])
                    }
                  }
                  getRealDataByUuid({uuid:newuuids,devices:newdevices}).then(function (res){
                    _t.settingLoading=false
                    if(res.data.code==0)
                    {
                      for(let k = 0,realDataLen=res.data.realData.length;k<realDataLen;k++)
                      {
                        let pushData = {
                          DeviceUuid:res.data.realData[k].duid,
                          ProjectUuid:res.data.realData[k].project_uuid,
                          Cmd:"RealData",
                          Data:[]
                        }
                        _t.CurrentRealUUIDList.push(res.data.realData[k].uuid)
                        _t.CurrentModelUUIDList.push(res.data.realData[k].mduid)
                        let DataObj = {
                          Uuid: res.data.realData[k].uuid,
                          ModelDataUuid: res.data.realData[k].mduid,
                          Value: res.data.realData[k].value
                        }
                        pushData.Data.push(DataObj)
                        _t.$EventBus.$emit("readDataPush", pushData);
                      }
                    }
                  })
                }
              },500)
            }
          })
        }
        else if (linkInfo.linkType == "External") {
          if(typeof linkInfo.OpenExternalType=='undefined') {
            window.open(linkInfo.External, '_blank')
          }else if( linkInfo.OpenExternalType=='self'){
            window.open(linkInfo.External, '_self')
          }else if(linkInfo.OpenExternalType=='new'){
            window.open(linkInfo.External, '_blank')
          }
        }
      }
    },
    getDisplayLoginPage(){
      let _t = this
      this.getLoginLayerDataStruct({pageType:this.isMobile?0:1,uuid:this.showUuid,cb:function (errno,project_uuid){
        if(errno!=0)
        {
          _t.$router.push('/login')
        }
        else
        {
          _t.RunLoginCavasContainerInit(_t.configData.layer.autoSize,_t.configData.layer.Padding)
          document.title=_t.configData.AppName
        }

          _t.$message.destroy()
          _t.chargePage = false
        }});
    },
  },
  mounted() {
    let _t = this
    _t.$EventBus.$on("AppLoginSuccess", (data) => {
      _t.setUser(data.user)
      _t.setRoles(data.roles)
    })
  },
  beforeDestroy() {
    speechSynthesis.cancel()
  }
}
</script>

<style lang="less">
::v-deep .x6-edge-tool-vertex {
  display: none !important;
}
::v-deep .run-login-graph-container {
  position: absolute;
  height: 100%;
  width: 100%;
  flex:1;
  overflow: auto;
}
.mymenu {
  padding: 2px 0;
  z-index: 100000;
}
.mymenu .menu-active {
  border-radius: 0;
  border-color: transparent;
}

.ism-render {
  position: absolute;
  transform-origin: 0 0 0; /*指定缩放的基本点*/
  -moz-transform-origin: 0 0 0;
  -ms-transform-origin: 0 0 0;
  -webkit-transform-origin: 0 0 0;
  -o-transform-origin: 0 0 0;
  .ism-render-wrapper {
    position: absolute;
  }

  .ism-render-wrapper-clickable {
    cursor: pointer;
  }
}
.ism-popup-render {
  transform-origin: 0 0 0; /*指定缩放的基本点*/
  -moz-transform-origin: 0 0 0;
  -ms-transform-origin: 0 0 0;
  -webkit-transform-origin: 0 0 0;
  -o-transform-origin: 0 0 0;
  .ism-render-wrapper {
    position: absolute;
  }

  .ism-render-wrapper-clickable {
    cursor: pointer;
  }
}

::-webkit-scrollbar {
  /*滚动条整体样式*/
  width : 5px;  /*高宽分别对应横竖滚动条的尺寸*/
  height: 9px;
}
::-webkit-scrollbar-thumb {
  /*滚动条里面小方块*/
  border-radius   : 10px;
  background-color: skyblue;
  background-image: -webkit-linear-gradient(
      45deg,
      rgba(255, 255, 255, 0.2) 25%,
      transparent 25%,
      transparent 50%,
      rgba(255, 255, 255, 0.2) 50%,
      rgba(255, 255, 255, 0.2) 75%,
      transparent 75%,
      transparent
  );
}
::-webkit-scrollbar-track {
  /*滚动条里面轨道*/
  box-shadow   : inset 0 0 5px rgba(0, 0, 0, 0.2);
  background   : #ededed;
  border-radius: 1px;
}
</style>
