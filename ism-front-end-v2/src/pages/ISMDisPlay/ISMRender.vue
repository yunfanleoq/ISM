<template>
  <div class="ism-render-root" :style="nodeMouseStyleVar" ref="ismrender">
    <div class="run-graph-container" :class="{'animated':true,[`${configData.layer.animate}`]: true}">
      <div ref="ISMRunningContainer" ></div>
    </div>
    <div v-show="pageLoading" class="ism-page-loading">
      <div class="ism-page-loading-panel">
        <div class="ism-page-loading-scan"></div>
        <div class="ism-page-loading-content">
          <div class="ism-page-loading-orbit">
            <div class="ism-page-loading-spinner"></div>
            <div class="ism-page-loading-core"></div>
          </div>
          <div class="ism-page-loading-copy">
            <div class="ism-page-loading-title">页面加载中</div>
            <div class="ism-page-loading-subtitle">正在渲染组件，请稍候</div>
          </div>
        </div>
        <div class="ism-page-loading-meter">
          <span></span><span></span><span></span><span></span>
        </div>
      </div>
    </div>

    <a-modal :maskClosable="false" :visible="settingDialog"
             :footer="null"
             @cancel="settingDialog=false"
             v-drag-modal
             :width="isMobile?250:300"
             :destroyOnClose="true"
             :title="$t('monitor.Set')"
    >
      <div  style="padding: 5px;">
        <a-form :form="SetForm">
          <a-form-item :label="$t('displayConfig.Properties.SetPassword')" style="margin-top: 1px;margin-bottom: 2px;"
                       :labelCol="{span: isMobile?8:6}"
                       :wrapperCol="{span: isMobile?16:18}"
                       v-if="SetPassword!=''"
          >
            <a-input @focus="onInputFocus('SetPassword')" :style="{'text-security':'disc', '-webkit-text-security':'disc'}" v-model="SetPasswordFormValue"/>
          </a-form-item>
          <a-form-item :label="$t('monitor.SetValue')" style="margin-top: 2px;margin-bottom: 2px;"
                       :labelCol="{span:isMobile?8:6}"
                       :wrapperCol="{span: isMobile?16:18}"
          >
            <a-textarea auto-size @focus="onInputFocus('SetValue')" v-model="SetValueFormValue"/>
          </a-form-item>
        </a-form>
      </div>
      <a-divider />
      <div class="dialog-button">
        <a-button key="submit" type="primary" :loading="settingLoading" @click="ManualSetData">
          {{ $t('component.deviceDataModel.submit')}}
        </a-button>
        <a-button style="margin-left: 10px" key="back" @click="settingDialog=false;showKeyboard=false">
          {{$t('component.deviceDataModel.cancel')}}
        </a-button>
      </div>
    </a-modal>

    <a-modal :maskClosable="false" :visible="setPasswordDialog"
             :footer="null"
             @cancel="setPasswordDialog=false"
             v-drag-modal
             :width="isMobile?250:300"
             :destroyOnClose="true"
             :title="$t('displayConfig.Properties.SetPasswordTips')"
             :draggable="true"
             :resizable="true"
             :modal="false">
      <div class="f-full" style="padding: 5px;">
        <a-form :form="SetFormPassword">
          <a-form-item :label="$t('displayConfig.Properties.SetPassword')" style="margin-top: 15px"
                       :labelCol="{span: isMobile?8:6}"
                       :wrapperCol="{span: isMobile?16:18}"
          >
            <a-input type="text" autocomplete=off :style="{'text-security':'disc', '-webkit-text-security':'disc'}"
                     @focus="onInputFocus('SetAutoPassword')" v-model="SetAutoPasswordFormValue"
            />
          </a-form-item>
        </a-form>
      </div>
      <a-divider />
      <div class="dialog-button">
        <a-button key="submit" type="primary" :loading="settingLoading" @click="PasswordSetData">
          {{ $t('component.deviceDataModel.submit')}}
        </a-button>
        <a-button style="margin-left: 10px;" key="back" @click="setPasswordDialog=false;showKeyboard=false">
          {{$t('component.deviceDataModel.cancel')}}
        </a-button>
      </div>
    </a-modal>

    <a-modal :maskClosable="false" :visible="actionPasswordDialog"
             :footer="null"
             @cancel="actionPasswordDialog=false"
             v-drag-modal
             :width="isMobile?250:300"
             :destroyOnClose="true"
             :title="$t('displayConfig.Properties.SetPasswordTips')"
             :draggable="true"
             :resizable="true"
             :modal="false">
      <div class="f-full" style="padding: 5px;">
        <a-form :form="ActionFormPassword">
          <a-form-item :label="$t('displayConfig.Properties.SetPassword')" style="margin-top: 15px"
                       :labelCol="{span: isMobile?8:6}"
                       :wrapperCol="{span: isMobile?16:18}"
          >
            <a-input type="text" autocomplete=off :style="{'text-security':'disc', '-webkit-text-security':'disc'}"
                     @focus="onInputFocus('SetActionPassword')" v-model="ActionPasswordValue"
            />
          </a-form-item>
        </a-form>
      </div>
      <a-divider />
      <div class="dialog-button">
        <a-button key="submit" type="primary" :loading="settingLoading" @click="PasswordSetAction">
          {{ $t('component.deviceDataModel.submit')}}
        </a-button>
        <a-button style="margin-left: 10px;" key="back" @click="actionPasswordDialog=false;showKeyboard=false">
          {{$t('component.deviceDataModel.cancel')}}
        </a-button>
      </div>
    </a-modal>
    <a-modal :visible="PopUpDialog"
             v-drag-modal
             :centered="true"
             :style="stylePopUpVar"
             :forceRender="true"
             :destroyOnClose="true"
             :width="(!isExternUrl?PopUpConfigData.layer.width:linkInfoExternal.width)+'px'"
             :title="typeof PopUpConfigData.PageName!='undefined'&&PopUpConfigData.PageName!=''?PopUpConfigData.PageName:'对话框'"
             :footer="null"
             @cancel="ClosePopDialog"
             :modal="false">
      <div  @click="PopUpDialogClick">
        <div v-if="!isExternUrl">
          <div  class="ism-popup-render"  @contextmenu.prevent="componentRightDialogClick" :style="popUpStyle" v-if="PopUpConfigData.layer" ref="ismrender1">
            <div class="run-graph-container" :class="{'animated':true,[`${PopUpConfigData.layer.animate}`]: true}">
              <div ref="ISMPopUpRunningContainer" ></div>
            </div>
            <!--          <Menu ref="componentMenuDialog"  menuCls="mymenu" @itemClick="componentMenuClickDialog($event)">-->
            <!--            <MenuItem value="autosize" iconCls="icon-autosize" :text="$t('Render.AutoSize')"></MenuItem>-->
            <!--            <MenuItem value="ScreenUndo" iconCls="icon-screenundo" :text="$t('Render.ScreenUndo')"></MenuItem>-->
            <!--          </Menu>-->
          </div>
        </div>
        <div v-else >
          <iframe scrolling="auto"  :src="linkInfoExternal.url" :width="linkInfoExternal.width+'px'" :height="linkInfoExternal.height+'px'" frameborder="0"></iframe>
        </div>
      </div>
    </a-modal>
    <div v-show="showKeyboard&&configData.layer.virtuallyKey">
      <SimpleKeyboard ref="SimpleKeyboard" @onChange="onChangeKeyboard" @onKeyPress="onKeyPress"/>
    </div>
    <LockScreen/>
  </div>
</template>

<script>
import {formatDate} from "@/utils/common";
import SimpleKeyboard from '@/components/SimpleKeyboard/SimpleKeyboard.vue'
const loadingKey = 'updatable'
import ISMBase from './ISMBase';
import store from "@/store";
import { mapActions, mapState, mapMutations ,mapGetters} from 'vuex'
import {setData} from "@/services/device";
import {ComponentRestApi} from "@/services/RestApi";
import TemplateRender from 'vue-template-render'
import {ExecSysScript} from "@/services/system";
import {AUTH_TYPE, checkAuthorization, getAuthorization, setAuthorization} from "@/utils/request";
import {getRealDataByUuid} from "@/services/device";
import Vue from 'vue'
import Contextmenu from "vue-contextmenujs"
Vue.use(Contextmenu);
import '@/utils/vmodalDrage'
import AKeepAlive from "@/components/cache/AKeepAlive";
import {Graph, Shape} from "@antv/x6";
import {Keyboard} from "@antv/x6";
import {register} from "@antv/x6-vue-shape";
import ISMGroupNode from "@/pages/ISMDisPlay/ISMGroupNode.vue";
import {snapdom} from "@/utils/snapdom.mjs";
import LockScreen from "@/components/lockScreen/LockScreen";
import Hammer from 'hammerjs';
let isISMGroupNodeRegistered = false;
export default {
  name: 'ISMRender',
  extends: ISMBase,
  i18n: require('../../i18n/language'),
  components: {
    SimpleKeyboard,
    LockScreen,
  },
  props: {
    showUuid: {
      type: String,
      required: true
    },
    showToken: {
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
          console.log('[watch:showUuid] ENTRY newVal=', newVal, 'oldVal=', oldVal, 'currentDisplayUUID=', this.currentDisplayUUID, 'currentPageUUID=', this.currentPageUUID)
          let _t =this
          // 立即作废旧请求 token，并销毁旧 Graph
          const requestToken = ++this.mainPageRequestToken
          console.log('[watch:showUuid] requestToken=', requestToken, 'mainPageRequestToken=', this.mainPageRequestToken)
          _t.destroyMainGraph()
          _t.clearPendingTimers('main')
          let page={
            pageType:1,
            displayUUID:newVal,
            pageUuid:this.currentPageUUID
          }
          _t.chargePage = true
          this.$message.loading({ content: 'Loading...', key: 'updatable', duration: 0 });
          this.CurrentRealUUIDList=[]
          this.CurrentModelUUIDList=[]
          this.CurrentPagerRealDataUuidList=[]
          this.CurrentPagerRealDeviceUuidList=[]
          console.log('[watch:showUuid] calling selectDisplayPageDataStruct, page=', JSON.parse(JSON.stringify(page)))
          this.selectDisplayPageDataStruct({page:page,callback:function (uuids,devices,isFound){
              console.log('[watch:showUuid] callback fired, isFound=', isFound, 'uuids.length=', uuids?.length, 'devices.length=', devices?.length)
              if (requestToken !== _t.mainPageRequestToken || _t._isDestroyed) {
                return
              }
              _t.$message.destroy('updatable');
              _t.chargePage = false
              if (isFound === false) {
                _t.currentDisplayUUID = ""
                _t.currentPageUUID = ""
                _t.$message.error(_t.$t("readData.NotFindPage"))
                return
              }
              _t.CurrentPagerRealDataUuidList = uuids
              _t.CurrentPagerRealDeviceUuidList=devices
              _t.RunCavasContainerInit(_t.configData.layer.autoSize,_t.configData.layer.Padding)
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
    nodeMouseStyleVar(){
      return {
        '--NodeMouseStyle':(this.SelectCurrentNodeData!=null)&&(typeof this.SelectCurrentNodeData.action !='undefined')&&(this.SelectCurrentNodeData.action.length > 0)?'pointer':'',
      };
    },
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
    },
    popUpStyle:function () {
      let scale = 1
      let styles = [`transform:scale(${scale})`];
      if(this.PopUpConfigData.layer.backColor) {
        styles.push(`background-color: ${this.PopUpConfigData.layer.backColor}`);
      }
      if(this.PopUpConfigData.layer.backgroundImage) {
        styles.push(`background-image: url("${this.PopUpConfigData.layer.backgroundImage}")`);
        styles.push(`background-size:100% 100%`);
      }
      if(this.PopUpConfigData.layer.width > 0) {
        styles.push(`width: ${this.PopUpConfigData.layer.width}px`);
      }
      if(this.PopUpConfigData.layer.height > 0) {
        styles.push(`height: ${this.PopUpConfigData.layer.height}px`);
      }
      let style = styles.join(';');
      return style;
    },
    stylePopUpVar() {
      return {
        '--popUpBackColor':this.PopUpConfigData.layer.backColor,
      };
    },
  },
  data() {
    return {
      hammerManager:null,
      tempAutoSizePager:false,
      CurrentPagerRealDataUuidList:[],
      CurrentPagerRealDeviceUuidList:[],
      CurrentPagerPopRealDataUuidList:[],
      CurrentPagerPopRealDeviceUuidList:[],
      tempAutoSize:0,
      tempAutoPadding:0,
      SelectCurrentNodeData:null,
      ISMPopUpRunningContainer:null,
      ISMRuningCavasContainer:null,
      settingDialog:false,
      IsConfirm:false,
      showKeyboard:false,
      changeIpt:'',//选择了哪个输入框
      setPasswordDialog:false,
      actionPasswordDialog:false,
      PopUpDialog:false,
      pageLoading:false,
      pageLoadingEl:null,
      chargePage:false,
      IsAutoClose:false,
      Zoom:100,
      SetPasswordFormValue:"",
      SetAutoPasswordFormValue:"",
      ActionPasswordValue:"",
      ActionPasswordSet:"",
      ActionComponent:null,
      ActionEvent:null,
      SetValueFormValue:"",
      clickTimer:null,
      pendingClickPageLoadingToken:null,
      zoomDebounceTimer:null,
      pendingTimers:[],
      pageLoadingToken:0,
      pageLoadingStartedAt:0,
      lastGoPageKey:"",
      lastGoPageAt:0,
      destroyMainGraphRaf:null,
      mainPageRequestToken:0,
      popUpPageRequestToken:0,
      pageRenderToken:0,
      popUpRenderToken:0,
      // EventBus handler 引用，用于精确 $off（Vue 2 不代理 _ 开头的 data 属性，所以不用下划线）
      eventHandlers: {},
      linkInfoExternal:{
        url:"",
        width:1024,
        height:768
      },
      chargePagePopUp:false,
      isExternUrl:false,
      isPopUpOpen:false,
      SelectDeviceUuid:"",
      currentDisplayUUID: "",
      currentPageUUID: "",
      currentPopUpDisplayUUID: "",
      currentPopUpPageUUID: "",
      deviceUuid:"",
      SetPassword:"",
      AutoSetValue:"",
      setDataUuid:"",
      settingVisible:false,
      settingLoading:false,
      SetForm:this.$form.createForm(this),
      ActionFormPassword:this.$form.createForm(this),
      SetFormPassword:this.$form.createForm(this),
      fullScreen:false,
      CurrentRealUUIDList:[],
      CurrentRealDataNameList:[],
      CurrentModelUUIDList:[],
      CurrentPopRealUUIDList:[],
      CurrentPopModelUUIDList:[]
    }
  },
  created(){
    this.addListener()
  },
  methods: {
    ...mapMutations('ISMDisPlayEditorTool',[
      'setlayerZoom',
    ]),
    ...mapActions('ISMDisPlayEditorTool',[
      'getLayerDataStruct',
      "lockScreen",
      'getLayerDataStructByTokenData',
      'selectDisplayPageDataStruct',
      'selectPopUpDisplayPageDataStruct'
    ]),
    async exportPng(){

      const el = this.$refs.ISMRunningContainer;
      const result = await snapdom(el, { scale: 2 });

      const img = await result.toPng();
      document.body.appendChild(img);

      await result.download({ format: 'png', filename: this.configData.name });
    },
    InitPagerRealData() {
      console.log('[InitPagerRealData] ENTRY, pageRenderToken=', this.pageRenderToken, 'CurrentPagerRealDataUuidList.length=', this.CurrentPagerRealDataUuidList?.length, 'CurrentPagerRealDeviceUuidList.length=', this.CurrentPagerRealDeviceUuidList?.length)
      const renderToken = this.pageRenderToken
      let _t = this
      const datauuid = this.CurrentPagerRealDataUuidList
      const devices = this.CurrentPagerRealDeviceUuidList
      _t.CurrentRealUUIDList=[]
      _t.CurrentModelUUIDList=[]
      this.trackTimeout(function (){
        console.log('[InitPagerRealData] trackTimeout callback, renderToken=', renderToken, 'pageRenderToken=', _t.pageRenderToken, '_isDestroyed=', _t._isDestroyed)
        if (_t._isDestroyed || renderToken !== _t.pageRenderToken) {
          console.warn('[InitPagerRealData] ABORTED: token mismatch or destroyed')
          return
        }
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
        console.log('[InitPagerRealData] calling getRealDataByUuid, uuids=', newuuids.length, 'devices=', newdevices.length)
        getRealDataByUuid({uuid: newuuids,devices:newdevices}).then(function (res) {
          console.log('[InitPagerRealData] getRealDataByUuid response, code=', res.data?.code, 'realData.length=', res.data?.realData?.length)
          if (_t._isDestroyed || renderToken !== _t.pageRenderToken) {
            console.warn('[InitPagerRealData] response ABORTED: token mismatch or destroyed')
            return
          }
          _t.settingLoading = false
          if (res.data.code == 0) {
            // 按 DeviceUuid 分组，批量 emit，减少事件数量和 handler 调用次数
            let pushDataMap = {}
            for (let k = 0; k < res.data.realData.length; k++) {
              if(res.data.realData[k].value=="")
              {
                continue
              }
              let duid = res.data.realData[k].duid
              if (!pushDataMap[duid]) {
                pushDataMap[duid] = {
                  DeviceUuid: duid,
                  ProjectUuid: res.data.realData[k].project_uuid,
                  Cmd: "RealData",
                  Data: []
                }
              }
              _t.CurrentRealUUIDList.push(res.data.realData[k].uuid)
              _t.CurrentModelUUIDList.push(res.data.realData[k].mduid)
              let DataObj = {
                Uuid: res.data.realData[k].uuid,
                ModelDataUuid: res.data.realData[k].mduid,
                Value: res.data.realData[k].value
              }
              pushDataMap[duid].Data.push(DataObj)
            }
            for (let duid in pushDataMap) {
              _t.$EventBus.$emit("readDataPush", pushDataMap[duid])
            }
          }
        })
      },0, 'main')
    },
    InitPopUpPagerRealData() {
      const renderToken = this.popUpRenderToken
      let _t = this
      const datauuid = this.CurrentPagerPopRealDataUuidList
      const devices = this.CurrentPagerPopRealDeviceUuidList
      _t.CurrentPopRealUUIDList=[]
      _t.CurrentPopModelUUIDList=[]
      this.trackTimeout(function (){
        if (_t._isDestroyed || renderToken !== _t.popUpRenderToken) {
          return
        }
        if(datauuid.length>0)
        {
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
          getRealDataByUuid({uuid:newuuids,devices:newdevices}).then(function (res){
            if (_t._isDestroyed || renderToken !== _t.popUpRenderToken) {
              return
            }
            _t.settingLoading=false
            if(res.data.code==0)
            {
              // 按 DeviceUuid 分组，批量 emit，减少事件数量和 handler 调用次数
              let pushDataMap = {}
              for(let k = 0,realDataLen = res.data.realData.length;k<realDataLen;k++)
              {
                let duid = res.data.realData[k].duid
                if (!pushDataMap[duid]) {
                  pushDataMap[duid] = {
                    DeviceUuid: duid,
                    ProjectUuid: res.data.realData[k].project_uuid,
                    Cmd: "RealData",
                    Data: []
                  }
                }
                _t.CurrentPopRealUUIDList.push(res.data.realData[k].uuid)
                _t.CurrentPopModelUUIDList.push(res.data.realData[k].mduid)
                let DataObj = {
                  Uuid: res.data.realData[k].uuid,
                  ModelDataUuid: res.data.realData[k].mduid,
                  Value: res.data.realData[k].value
                }
                pushDataMap[duid].Data.push(DataObj)
              }
              for (let duid in pushDataMap) {
                _t.$EventBus.$emit("readDataPush", pushDataMap[duid])
              }
            }
          })
        }
      },0, 'popup')
    },
    getInitialPageId() {
      return this.$route && this.$route.query ? this.$route.query.pageId : ""
    },
    loadPager(){
      let _t = this
      if(this.showUuid!="")
      {
        if(this.showToken!="")
        {
          const requestToken = ++this.mainPageRequestToken
          // 立即销毁旧 Graph，阻断旧实例继续消费数据/定时器
          _t.destroyMainGraph()
          _t.clearPendingTimers('main')
          // 标记正在切换，阻止数据推送期间无效的 DealWithUpdateData
          _t.chargePage = true
          this.$message.loading({content: 'Loading...', key: 'updatable', duration: 0});
          this.CurrentRealUUIDList=[]
          this.CurrentModelUUIDList=[]
          _t.CurrentPagerRealDataUuidList=[]
          _t.CurrentPagerRealDeviceUuidList=[]
          this.getLayerDataStructByTokenData({
            pageType: this.isMobile, token:this.showToken,uuid: this.showUuid, cb: function (errno, project_uuid,expireAt,token, datauuid,devices) {
              if (requestToken !== _t.mainPageRequestToken || _t._isDestroyed) {
                return
              }
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
              const initialPageId = _t.getInitialPageId()
              if(initialPageId) {
                _t.showPage({
                  linkType: "Inside",
                  Inside: {
                    displayUUID: _t.showUuid,
                    pageUUID: initialPageId,
                    displayType: 1
                  }
                })
                return
              }
              _t.CurrentPagerRealDataUuidList = datauuid
              _t.CurrentPagerRealDeviceUuidList = devices
              _t.RunCavasContainerInit(_t.configData.layer.autoSize,_t.configData.layer.Padding)
            }
          });
        }
        else {
          const requestToken = ++this.mainPageRequestToken
          // 立即销毁旧 Graph，阻断旧实例继续消费数据/定时器
          _t.destroyMainGraph()
          _t.clearPendingTimers('main')
          // 标记正在切换，阻止数据推送期间无效的 DealWithUpdateData
          _t.chargePage = true
          this.$message.loading({content: 'Loading...', key: 'updatable', duration: 0});
          this.CurrentRealUUIDList=[]
          this.CurrentModelUUIDList=[]
          this.CurrentPagerRealDataUuidList=[]
          this.CurrentPagerRealDeviceUuidList=[]
          this.getLayerDataStruct({
            pageType: this.isMobile, uuid: this.showUuid, cb: function (errno, project_uuid, datauuid,devices) {
              if (requestToken !== _t.mainPageRequestToken || _t._isDestroyed) {
                return
              }
              document.title = _t.configData.AppName
              if (errno == 0) {
                if (!checkAuthorization(AUTH_TYPE.AUTH1)) {
                  setAuthorization({token: project_uuid}, AUTH_TYPE.AUTH1)
                }
              } else {
                _t.$router.push('/login')
                _t.$message.destroy()
                _t.chargePage = false
                return
              }
              _t.$message.destroy()
              _t.chargePage = false
              const initialPageId = _t.getInitialPageId()
              if(initialPageId) {
                _t.showPage({
                  linkType: "Inside",
                  Inside: {
                    displayUUID: _t.showUuid,
                    pageUUID: initialPageId,
                    displayType: 1
                  }
                })
                return
              }
              _t.CurrentPagerRealDataUuidList = datauuid
              _t.CurrentPagerRealDeviceUuidList = devices
              _t.RunCavasContainerInit(_t.configData.layer.autoSize,_t.configData.layer.Padding)
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
      if (this.$refs.SimpleKeyboard) {
        this.$refs.SimpleKeyboard.onKeyPress('{clear}')
      }
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
              label: _t.$t('Render.LockScreen'),
              icon: "el-icon-suoping",
              onClick: () => {
                _t.lockScreen()
              }
            },
            {
              label: _t.tempAutoSize==1?_t.$t('Render.ScreenUndo'):_t.$t('Render.AutoSize'),
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
                this.ISMRuningCavasContainer.togglePanning(_t.tempAutoPadding)
                this.ISMRuningCavasContainer.toggleMouseWheel(_t.tempAutoPadding)
              }
            },
            {
              label: _t.$t('Render.ExportPng'),
              icon: "el-icon-daochupng",
              onClick: () => {
                _t.exportPng()
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
    formatDateTime(time){
      let date = new Date(time)
      return formatDate(date,'yyyy-MM-dd hh:mm:ss')
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
    trackTimeout(handler, delay, scope) {
      const timer = setTimeout(() => {
        this.pendingTimers = this.pendingTimers.filter(item => item.timer !== timer)
        handler()
      }, delay)
      this.pendingTimers.push({ timer, scope: scope || 'general' })
      return timer
    },
    clearPendingTimers(scope) {
      const keep = []
      this.pendingTimers.forEach(item => {
        if (!scope || item.scope === scope || item.scope === 'general') {
          clearTimeout(item.timer)
        } else {
          keep.push(item)
        }
      })
      this.pendingTimers = keep
      if ((!scope || scope === 'main') && this.clickTimer) {
        clearTimeout(this.clickTimer)
        this.clickTimer = null
      }
      if ((!scope || scope === 'main') && this.zoomDebounceTimer) {
        clearTimeout(this.zoomDebounceTimer)
        this.zoomDebounceTimer = null
      }
    },
    buildGoPageKey(wsData) {
      return JSON.stringify({
        ModelId: wsData.ModelId || "",
        PageUuid: wsData.PageUuid || "",
        IsPopUp: !!wsData.IsPopUp,
        AutoClose: !!wsData.AutoClose,
        linkType: typeof wsData.linkType !== "undefined" ? wsData.linkType : "Inside",
        width: wsData.width || "",
        height: wsData.height || "",
        External: wsData.External || "",
        OpenExternalType: wsData.OpenExternalType || ""
      })
    },
    shouldHandleGoPage(wsData) {
      const now = Date.now()
      const key = this.buildGoPageKey(wsData)
      if (this.lastGoPageKey === key && now - this.lastGoPageAt < 5000) {
        return false
      }
      this.lastGoPageKey = key
      this.lastGoPageAt = now
      return true
    },
    destroyHammerManager() {
      if (this.hammerManager) {
        this.hammerManager.destroy();
        this.hammerManager = null;
      }
      if (this.zoomDebounceTimer) {
        clearTimeout(this.zoomDebounceTimer)
        this.zoomDebounceTimer = null
      }
    },
    destroyMainGraph() {
      this.destroyHammerManager()
      window.removeEventListener('resize', this.setScale)
      this.pageRenderToken += 1

      // 用 try/finally 保证即使 clearCells/dispose 抛异常，
      // Graph 引用也一定被清空，避免下次进入复用路径叠加 DOM
      try {
        if (this.ISMRuningCavasContainer) {
          // 先异步=false 时同步清理 cell views；若抛异常进入 catch
          this.ISMRuningCavasContainer.clearCells();
          this.ISMRuningCavasContainer.off();
          this.ISMRuningCavasContainer.dispose();
        }
      } catch (e) {
        console.error('[destroyMainGraph] dispose failed, will force cleanup', e)
      } finally {
        // 无论 dispose 成功与否，强制清理
        if (this.ISMRuningCavasContainer) {
          try { this.ISMRuningCavasContainer.off() } catch(e) { console.warn(e) }
          try { this.ISMRuningCavasContainer.dispose() } catch(e) { console.warn(e) }
          // 延迟清空容器 DOM，避免同步操作阻塞主线程导致页面切换卡顿
          try {
            const container = this.$refs.ISMRunningContainer
            if (container) {
              this.destroyMainGraphRaf = requestAnimationFrame(() => {
                console.log('[destroyMainGraph] RAF callback executing, clearing innerHTML, childCount was=', container.children?.length)
                container.innerHTML = ''
                this.destroyMainGraphRaf = null
              })
            }
            } catch(e) { console.warn(e) }
        }
        this.ISMRuningCavasContainer = null
      }

      // Graph 销毁后，全量清理子组件可能残留的 EventBus 监听器
      // （子组件 beforeDestroy 可能缺失或不完整）
      const identifiers = this.collectAllIdentifiers(false)
      const popIdentifiers = new Set(this.PopUpDialog ? this.collectAllIdentifiers(true) : [])
      for (const id of identifiers) {
        if (!popIdentifiers.has(id)) {
          this.$EventBus.$off(id + "activeEvent")
          this.$EventBus.$off(id + "animateEvent")
        }
      }

      // 立即重新注册 ISMRender 自身的 EventBus 监听器
      this.rebindISMHandlers()
    },
    destroyPopUpGraph() {
      this.popUpRenderToken += 1

      try {
        if (this.ISMPopUpRunningContainer) {
          this.ISMPopUpRunningContainer.clearCells();
          this.ISMPopUpRunningContainer.off();
          this.ISMPopUpRunningContainer.dispose();
        }
      } catch (e) {
        console.error('[destroyPopUpGraph] dispose failed, will force cleanup', e)
      } finally {
        if (this.ISMPopUpRunningContainer) {
          try { this.ISMPopUpRunningContainer.off() } catch(e) { console.warn(e) }
          try { this.ISMPopUpRunningContainer.dispose() } catch(e) { console.warn(e) }
          // 延迟清空容器 DOM，避免同步操作阻塞主线程
          try {
            const container = this.$refs.ISMPopUpRunningContainer
            if (container) {
              requestAnimationFrame(() => { container.innerHTML = '' })
            }
          } catch(e) { console.warn(e) }
        }
        this.ISMPopUpRunningContainer = null
      }
      this.clearPendingTimers('popup')

      // 清理弹窗子组件可能残留的 EventBus 监听器
      const popIdentifiers = this.collectAllIdentifiers(true)
      const mainIdentifiers = new Set(this.ISMRuningCavasContainer ? this.collectAllIdentifiers(false) : [])
      for (const id of popIdentifiers) {
        if (!mainIdentifiers.has(id)) {
          this.$EventBus.$off(id + "activeEvent")
          this.$EventBus.$off(id + "animateEvent")
        }
      }

      // 重新注册 ISMRender 自身 handler
      this.rebindISMHandlers()
    },
    collectAllIdentifiers(popUpOnly) {
      const configs = popUpOnly
        ? [this.PopUpConfigData]
        : [this.configData]
      const ids = []
      for (const config of configs) {
        if (config && config.components && config.components.cells) {
          for (const cell of config.components.cells) {
            const id = cell.data && cell.data.detail && cell.data.detail.identifier
            if (id) ids.push(id)
          }
        }
      }
      return ids
    },
    rebindISMHandlers() {
      const h = this.eventHandlers
      // 先 $off 再 $on，防止 destroyMainGraph/destroyPopUpGraph 后重复叠加 handler
      if (h.readDataPush) {
        this.$EventBus.$off("readDataPush", h.readDataPush)
        this.$EventBus.$on("readDataPush", h.readDataPush)
      }
      if (h.onSelectDevice) {
        this.$EventBus.$off("onSelectDevice", h.onSelectDevice)
        this.$EventBus.$on("onSelectDevice", h.onSelectDevice)
      }
      if (h.SystemData) {
        this.$EventBus.$off("SystemData", h.SystemData)
        this.$EventBus.$on("SystemData", h.SystemData)
      }
      if (h.StaticData) {
        this.$EventBus.$off("StaticData", h.StaticData)
        this.$EventBus.$on("StaticData", h.StaticData)
      }
      if (h.RealAlarm) {
        this.$EventBus.$off("RealAlarm", h.RealAlarm)
        this.$EventBus.$on("RealAlarm", h.RealAlarm)
      }
      if (h.ChargePage) {
        this.$EventBus.$off("ChargePage", h.ChargePage)
        this.$EventBus.$on("ChargePage", h.ChargePage)
      }
      if (h.PlayVoice) {
        this.$EventBus.$off("PlayVoice", h.PlayVoice)
        this.$EventBus.$on("PlayVoice", h.PlayVoice)
      }
      if (h.GoPage) {
        this.$EventBus.$off("GoPage", h.GoPage)
        this.$EventBus.$on("GoPage", h.GoPage)
      }
    },
    getGraphNodeDetail(node) {
      if (!node) return null
      let data = null
      try {
        if (typeof node.getData === 'function') {
          data = node.getData()
        }
      } catch (e) {
        console.warn(e)
      }
      if (!data) {
        try {
          const props = typeof node.prop === 'function' ? node.prop() : null
          data = props && props.data ? props.data : null
        } catch (e) {
          console.warn(e)
        }
      }
      if (!data) return null
      return data.detail || data
    },
    handleGraphNodeAction(node, active) {
      try {
        const component = this.getGraphNodeDetail(node)
        if (!component) return
        this.SelectCurrentNodeData = component
        this.doComponentAction(component, active)
      } catch (e) {
        console.warn(e)
      }
    },
    registerISMGroupNode() {
      if (isISMGroupNodeRegistered) {
        return
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
      isISMGroupNodeRegistered = true
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
        this.ISMRuningCavasContainer.resize(window.innerWidth, window.innerHeight);
        const dd = window.innerWidth / this.configData.layer.width;
        const dh = window.innerHeight / this.configData.layer.height;
        this.ISMRuningCavasContainer.scale( dd,  dh);
      }
      else {
        this.ISMRuningCavasContainer.resize(this.configData.layer.width, this.configData.layer.height);
        this.ISMRuningCavasContainer.scale( 1,  1);
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
    hasPageJumpAction(component, active) {
      if (!component || !Array.isArray(component.action)) {
        return false
      }
      return component.action.some(action => {
        if (!action || action.type !== active || action.actionConfirm) {
          return false
        }
        if (action.action === 'link') {
          return true
        }
        return action.action === 'DeviceView' && !(action.DeviceView && action.DeviceView.isContainer)
      })
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
    /**
 * 处理组件的交互动作，根据触发类型（单击/双击）执行对应的动作或激活状态
 * @param {Object} component - 目标组件对象，包含 action（动作列表）或 active（激活状态列表）属性
 * @param {string} active - 触发类型，支持 'click'（单击）和 'dblclick'（双击）
 * @description
 * 单击时通过 300ms 延迟区分单击与双击事件：
 *   - 若组件存在 action 列表，则匹配类型后播放语音并调用 handleComponentAction 执行动作；
 *   - 若组件存在 active 列表，则查找 isSwitch 为 true 的激活项，播放语音后根据是否配置二次确认弹窗，
 *     决定直接调用 handleComponentSwitchActive 或经用户确认后调用。
 * 双击时清除单击延迟计时器，使单击动作不再触发。
 */
  doComponentAction(component,active){
      let _this = this
      if(active=='click')
      {
        clearTimeout(this.clickTimer);  //首先清除计时器
        if (this.hasPageJumpAction(component, active)) {
          this.pendingClickPageLoadingToken = this.beginPageLoading(loadingKey)
        } else {
          this.pendingClickPageLoadingToken = null
        }
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
                  if(typeof componentActive.condition.ConfirmationDialog!='undefined'&&componentActive.condition.ConfirmationDialog==1) {
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
                  else
                  {
                    _this.handleComponentSwitchActive(componentActive);
                  }
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
      if (active == 'dblclick' && this.pendingClickPageLoadingToken) {
        this.closePageLoading(loadingKey, this.pendingClickPageLoadingToken)
        this.pendingClickPageLoadingToken = null
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
      this.$message.destroy()
      this.PopUpDialog  = false
      this.popUpPageRequestToken += 1
      this.currentPopUpDisplayUUID = ""
      this.currentPopUpPageUUID = ""
      this.destroyPopUpGraph()
    },
    PopUpDialogClick(){

      this.$message.destroy()
      if(this.IsAutoClose)
      {
        this.PopUpDialog  = false
        this.popUpPageRequestToken += 1
        this.currentPopUpDisplayUUID = ""
        this.currentPopUpPageUUID = ""
        this.destroyPopUpGraph()
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
        this.trackTimeout(function() {
          _this.doLoopSetValue(counter + 1,setvalue,delay);
        }, delay);
      }
    },
    handleComponentAction(component,action){

      if(typeof action.action=="undefined")
      {
        this.cancelPendingPageLoading()
        return
      }
      let _this = this;
      try {
        if(this.showToken=="") {
          if ((typeof this.user.Role != "undefined") && (this.user.Role != "Admin")) {
            if ((typeof action.actionAuth != "undefined") && (action.actionAuth.length > 0)) {
              let auth = action.actionAuth.find(item => item === this.user.Role)
              if (auth != this.user.Role) {
                this.cancelPendingPageLoading()
                this.$message.error(this.$t('displayModel.NoAuthentication'), 3)
                return
              }
            }
          }
        }
      } catch (e) {
        this.cancelPendingPageLoading()
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
          this.ClosePopDialog()
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
        this.cancelPendingPageLoading()
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
          this.ClosePopDialog()
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
      if(typeof realData=="undefined" || typeof realData.Data=="undefined" || realData.Data==null || realData.Data.length==0)
      {
        return
      }
      for (let k = 0,Datalen = realData.Data.length; k < Datalen; k++) {
        if((_t.CurrentRealUUIDList.indexOf(realData.Data[k].Uuid)==-1)&&(_t.CurrentModelUUIDList.indexOf(realData.Data[k].ModelDataUuid)==-1))
        {
          continue
        }
        if (typeof _t.configData.components === "undefined" || !_t.configData.components.cells) {
          return;
        }
        const cells = _t.configData.components.cells;
        if (!Array.isArray(cells) || cells.length === 0) {
          return;
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
                      const cell = _t.ISMRuningCavasContainer.getCellById(_t.configData.components.cells[j].id)
                      if (isStart) {
                        cell.setVisible(true)
                        _t.configData.components.cells[j].data.detail.style.visible = true
                      } else {
                        _t.configData.components.cells[j].data.detail.style.visible = false
                        cell.setVisible(false)
                      }
                    }
                    _t.$EventBus.$emit(_t.configData.components.cells[j].data.detail.identifier + "animateEvent", isStart);
                  }
                }
              }
            }
            else if(realData.DeviceUuid==_t.SelectDeviceUuid || (!_t.SelectDeviceUuid && realData.DeviceUuid==condition.deviceSN))
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
                    const cell = _t.ISMRuningCavasContainer.getCellById(_t.configData.components.cells[j].id)
                    if (isStart) {
                      cell.setVisible(true)
                      _t.configData.components.cells[j].data.detail.style.visible = true
                    } else {
                      _t.configData.components.cells[j].data.detail.style.visible = false
                      cell.setVisible(false)
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
                  const cell = _t.ISMRuningCavasContainer.getCellById(_t.configData.components.cells[j].id)
                  const Parent  = cell.getParent()
                  if(Parent)
                  {
                    Parent.getChildren().forEach(child => {
                      const pos = cell.position()
                      child.translate(RealValue, pos.y);
                    })
                  }
                  const pos = cell.position()
                  cell.position(RealValue, pos.y);
                }
              }
            }
            else if(realData.DeviceUuid==_t.SelectDeviceUuid || (!_t.SelectDeviceUuid && realData.DeviceUuid==conditionx.deviceSN))
            {
                  if ((conditionx.dataID == realData.Data[k].ModelDataUuid)||(conditionx.dataID == realData.Data[k].Uuid)) {
                    const RealValue = parseFloat(realData.Data[k].Value)
                    const cell = _t.ISMPopUpRunningContainer.getCellById(_t.PopUpConfigData.components.cells[j].id)
                    const Parent  = cell.getParent()
                if(Parent)
                {
                  Parent.getChildren().forEach(child => {
                    const pos = cell.position()
                    child.translate(RealValue, pos.y);
                  })
                }
                const pos = cell.position()
                cell.position(RealValue, pos.y);
              }
            }

            //动画
            if(conditiony.isBandDevice)
            {
              if(realData.DeviceUuid==conditiony.deviceSN) {

                  if ((conditiony.dataID == realData.Data[k].ModelDataUuid)||(conditiony.dataID == realData.Data[k].Uuid)) {

                    const RealValue = parseFloat(realData.Data[k].Value)
                    const cell = _t.ISMPopUpRunningContainer.getCellById(_t.PopUpConfigData.components.cells[j].id)
                    const Parent  = cell.getParent()
                  if(Parent)
                  {
                    Parent.getChildren().forEach(child => {
                      const pos = cell.position()
                      child.translate(pos.x,RealValue);
                    })
                  }
                  const pos = cell.position()
                  cell.position(pos.x,RealValue);
                }
              }
            }
            else if(realData.DeviceUuid==_t.SelectDeviceUuid || (!_t.SelectDeviceUuid && realData.DeviceUuid==conditiony.deviceSN))
            {
              if ((conditiony.dataID == realData.Data[k].ModelDataUuid)||(conditiony.dataID == realData.Data[k].Uuid)) {
                const RealValue = parseFloat(realData.Data[k].Value)
                const cell = _t.ISMRuningCavasContainer.getCellById(_t.configData.components.cells[j].id)
                const Parent  = cell.getParent()
                if(Parent)
                {
                  Parent.getChildren().forEach(child => {
                    const pos = cell.position()
                    child.translate(pos.x,RealValue);
                  })
                }
                const pos = cell.position()
                cell.position(pos.x,RealValue);
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
                    const cell = _t.ISMRuningCavasContainer.getCellById(_t.configData.components.cells[j].id)
                    if (!cell) continue
                    const isEdge = cell.isEdge()
                    if(!isEdge) {
                      _t.$EventBus.$emit(_t.configData.components.cells[j].data.detail.identifier + "activeEvent", activeData);
                    }else{
                      let animation = ''
                      if(activeData.ID == "Forward")
                      {
                        _t.configData.components.cells[j].data.detail.active[activeIndex].ForwardResult = activeData.result
                      }
                      else if(activeData.ID == "Reverse")
                      {
                        _t.configData.components.cells[j].data.detail.active[activeIndex].ReverseResult = activeData.result
                      }
                      const ForwardResult = _t.configData.components.cells[j].data.detail.active[0]?.ForwardResult
                      const ReverseResult = _t.configData.components.cells[j].data.detail.active[1]?.ReverseResult
                      if(ForwardResult)
                      {
                        animation = 'ant-line-forward 30s infinite linear'
                      }
                      else if(ReverseResult)
                      {
                        animation = 'ant-line-inverse 30s infinite linear'
                      }
                      else
                      {
                        animation = 'none'
                      }
                      cell.setAttrs({
                        line: {
                          style: {
                            animation: animation,
                          }
                        }
                      });
                    }
                  }
                }
              }
              else if(realData.DeviceUuid==_t.SelectDeviceUuid || (!_t.SelectDeviceUuid && realData.DeviceUuid==active[activeIndex].condition.deviceSN))
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
                  const cell = _t.ISMRuningCavasContainer.getCellById(_t.configData.components.cells[j].id)
                  const isEdge = cell.isEdge()
                  if(!isEdge) {
                    _t.$EventBus.$emit(_t.configData.components.cells[j].data.detail.identifier + "activeEvent", activeData);
                  }else{
                    let animation = ''
                    if(activeData.ID == "Forward")
                    {
                      _t.configData.components.cells[j].data.detail.active[activeIndex].ForwardResult = activeData.result
                    }
                    else if(activeData.ID == "Reverse")
                    {
                      _t.configData.components.cells[j].data.detail.active[activeIndex].ReverseResult = activeData.result
                    }
                    const ForwardResult = _t.configData.components.cells[j].data.detail.active[0].ForwardResult
                    const ReverseResult = _t.configData.components.cells[j].data.detail.active[1].ReverseResult
                    if(ForwardResult)
                    {
                      animation = 'ant-line-forward 30s infinite linear'
                    }
                    else if(ReverseResult)
                    {
                      animation = 'ant-line-inverse 30s infinite linear'
                    }
                    else
                    {
                      animation = 'none'
                    }
                    cell.setAttrs({
                      line: {
                        style: {
                          animation: animation,
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
          if (typeof _t.PopUpConfigData.components === "undefined" || !_t.PopUpConfigData.components.cells) {
            return;
          }
          const cells = _t.PopUpConfigData.components.cells;
          if (!Array.isArray(cells) || cells.length === 0) {
            return;
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
                        const cell = _t.ISMPopUpRunningContainer.getCellById(_t.PopUpConfigData.components.cells[j].id)
                        if (isStart) {
                          cell.setVisible(true)
                          _t.PopUpConfigData.components.cells[j].data.detail.style.visible = true
                        } else {
                          _t.PopUpConfigData.components.cells[j].data.detail.style.visible = false
                          cell.setVisible(false)
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
                      const cell = _t.ISMPopUpRunningContainer.getCellById(_t.PopUpConfigData.components.cells[j].id)
                      if (isStart) {
                        cell.setVisible(true)
                        _t.PopUpConfigData.components.cells[j].data.detail.style.visible = true
                      } else {
                        _t.PopUpConfigData.components.cells[j].data.detail.style.visible = false
                        cell.setVisible(false)
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
                    const cell = _t.ISMPopUpRunningContainer.getCellById(_t.PopUpConfigData.components.cells[j].id)
                    const Parent  = cell.getParent()
                    if(Parent)
                    {
                      Parent.getChildren().forEach(child => {
                        const pos = cell.position()
                        child.translate(RealValue, pos.y);
                      })
                    }
                    const pos = cell.position()
                    cell.position(RealValue, pos.y);
                  }
                }
              }
              else if(realData.DeviceUuid==_t.SelectDeviceUuid)
              {
                if ((conditionx.dataID == realData.Data[k].ModelDataUuid)||(conditionx.dataID == realData.Data[k].Uuid)) {
                  const RealValue = parseFloat(realData.Data[k].Value)
                  const cell = _t.ISMPopUpRunningContainer.getCellById(_t.PopUpConfigData.components.cells[j].id)
                  const Parent  = cell.getParent()
                  if(Parent)
                  {
                    Parent.getChildren().forEach(child => {
                      const pos = cell.position()
                      child.translate(RealValue, pos.y);
                    })
                  }
                  const pos = cell.position()
                  cell.position(RealValue, pos.y);
                }
              }

              //动画
              if(conditiony.isBandDevice)
              {
                if(realData.DeviceUuid==conditiony.deviceSN) {

                  if ((conditiony.dataID == realData.Data[k].ModelDataUuid)||(conditiony.dataID == realData.Data[k].Uuid)) {

                    const RealValue = parseFloat(realData.Data[k].Value)
                    const cell = _t.ISMPopUpRunningContainer.getCellById(_t.PopUpConfigData.components.cells[j].id)
                    const Parent  = cell.getParent()
                    if(Parent)
                    {
                      Parent.getChildren().forEach(child => {
                        const pos = cell.position()
                        child.translate(pos.x,RealValue);
                      })
                    }
                    const pos = cell.position()
                    cell.position(pos.x,RealValue);
                  }
                }
              }
              else if(realData.DeviceUuid==_t.SelectDeviceUuid)
              {
                if ((conditiony.dataID == realData.Data[k].ModelDataUuid)||(conditiony.dataID == realData.Data[k].Uuid)) {
                  const RealValue = parseFloat(realData.Data[k].Value)
                  const cell = _t.ISMPopUpRunningContainer.getCellById(_t.PopUpConfigData.components.cells[j].id)
                  const Parent  = cell.getParent()
                  if(Parent)
                  {
                    Parent.getChildren().forEach(child => {
                      const pos = cell.position()
                      child.translate(pos.x,RealValue);
                    })
                  }
                  const pos = cell.position()
                  cell.position(pos.x,RealValue);
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
                      const cell = _t.ISMPopUpRunningContainer.getCellById(_t.PopUpConfigData.components.cells[j].id)
                      const isEdge = cell.isEdge()
                      if(!isEdge) {
                        _t.$EventBus.$emit(_t.PopUpConfigData.components.cells[j].data.detail.identifier + "activeEvent", activeData);
                      }else{
                        let animation = ''
                        if(activeData.ID == "Forward")
                        {
                          _t.PopUpConfigData.components.cells[j].data.detail.active[activeIndex].ForwardResult = activeData.result
                        }
                        else if(activeData.ID == "Reverse")
                        {
                          _t.PopUpConfigData.components.cells[j].data.detail.active[activeIndex].ReverseResult = activeData.result
                        }
                        const ForwardResult = _t.PopUpConfigData.components.cells[j].data.detail.active[0]?.ForwardResult
                        const ReverseResult = _t.PopUpConfigData.components.cells[j].data.detail.active[1]?.ReverseResult
                        if(ForwardResult)
                        {
                          animation = 'ant-line-forward 30s infinite linear'
                        }
                        else if(ReverseResult)
                        {
                          animation = 'ant-line-inverse 30s infinite linear'
                        }
                        else
                        {
                          animation = 'none'
                        }
                        cell.setAttrs({
                          line: {
                            style: {
                              animation: animation,
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
                    const cell = _t.ISMPopUpRunningContainer.getCellById(_t.PopUpConfigData.components.cells[j].id)
                    if (cell) {
                      const isEdge = cell.isEdge()
                      if (!isEdge) {
                        _t.$EventBus.$emit(_t.PopUpConfigData.components.cells[j].data.detail.identifier + "activeEvent", activeData);
                      } else {
                        let animation = ''
                        if(activeData.ID == "Forward")
                        {
                          _t.PopUpConfigData.components.cells[j].data.detail.active[activeIndex].ForwardResult = activeData.result
                        }
                        else if(activeData.ID == "Reverse")
                        {
                          _t.PopUpConfigData.components.cells[j].data.detail.active[activeIndex].ReverseResult = activeData.result
                        }
                        const ForwardResult = _t.PopUpConfigData.components.cells[j].data.detail.active[0]?.ForwardResult
                        const ReverseResult = _t.PopUpConfigData.components.cells[j].data.detail.active[1]?.ReverseResult
                        if(ForwardResult)
                        {
                          animation = 'ant-line-forward 30s infinite linear'
                        }
                        else if(ReverseResult)
                        {
                          animation = 'ant-line-inverse 30s infinite linear'
                        }
                        else
                        {
                          animation = 'none'
                        }
                        cell.setAttrs({
                          line: {
                            style: {
                              animation: animation,
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
      }
      _t.$EventBus.$emit("DealWithRealDataFinish", true);

    },
    showComponent(identifier,visible) {
      const cell = this.ISMRuningCavasContainer.getCellById(identifier)
      cell.setVisible(visible)
      let getChildren = cell.getChildren()
      if(getChildren!=null)
      {
        getChildren.forEach(child => {
          child.setVisible(visible)
          let NodeInfoData = child.getData()
          NodeInfoData.UpdateNodeFlag = new Date().getTime()
          NodeInfoData.detail.style.visible = visible
          child.setData(NodeInfoData, {overwrite: true})
        })
      }

      let nodeData = cell.getData()
      nodeData.UpdateNodeFlag = new Date().getTime()
      nodeData.detail.style.visible = visible
      cell.setData(nodeData, {overwrite: true})
    },
    beginPageLoading(key = loadingKey) {
      this.$message.destroy(key)
      this.pageLoadingToken += 1
      this.pageLoadingStartedAt = Date.now()
      this.pageLoading = true
      this.mountBodyPageLoading()
      return this.pageLoadingToken
    },
    openPageLoading(key = loadingKey) {
      const token = this.beginPageLoading(key)
      return new Promise(resolve => {
        this.$nextTick(() => {
          const raf = window.requestAnimationFrame || function (cb) { return setTimeout(cb, 0) }
          raf(() => resolve(token))
        })
      })
    },
    consumePendingPageLoading() {
      const token = this.pendingClickPageLoadingToken
      this.pendingClickPageLoadingToken = null
      if (token && token === this.pageLoadingToken && this.pageLoading) {
        return Promise.resolve(token)
      }
      return this.openPageLoading(loadingKey)
    },
    cancelPendingPageLoading() {
      if (this.pendingClickPageLoadingToken) {
        this.closePageLoading(loadingKey, this.pendingClickPageLoadingToken)
        this.pendingClickPageLoadingToken = null
      }
    },
    closePageLoading(key = loadingKey, token = this.pageLoadingToken) {
      if (token !== this.pageLoadingToken) {
        return
      }
      const elapsed = Date.now() - this.pageLoadingStartedAt
      const close = () => {
        if (token !== this.pageLoadingToken) {
          return
        }
        this.pageLoading = false
        this.unmountBodyPageLoading()
        this.$message.destroy(key)
      }
      if (elapsed < 400) {
        setTimeout(close, 400 - elapsed)
      } else {
        close()
      }
    },
    mountBodyPageLoading() {
      if (typeof document === 'undefined') {
        return
      }
      if (!this.pageLoadingEl) {
        if (!document.getElementById('ism-page-loading-style')) {
          const style = document.createElement('style')
          style.id = 'ism-page-loading-style'
          style.textContent = `
            .ism-body-page-loading{position:fixed;left:0;top:0;right:0;bottom:0;z-index:2147483647;display:flex;align-items:center;justify-content:center;background:radial-gradient(circle at 50% 42%,rgba(14,116,144,.26),transparent 32%),linear-gradient(180deg,rgba(2,6,23,.52),rgba(2,6,23,.7));pointer-events:none}
            .ism-body-page-loading .ism-page-loading-panel{position:relative;width:310px;padding:28px 30px 24px;border:1px solid rgba(103,232,249,.56);background:linear-gradient(180deg,rgba(9,26,46,.94),rgba(4,13,27,.9));box-shadow:0 24px 70px rgba(0,0,0,.48),0 0 36px rgba(34,211,238,.24),inset 0 1px 0 rgba(255,255,255,.14);backdrop-filter:blur(14px);overflow:hidden}
            .ism-body-page-loading .ism-page-loading-panel:before{content:"";position:absolute;inset:10px;border:1px solid rgba(125,211,252,.14);pointer-events:none}
            .ism-body-page-loading .ism-page-loading-scan{position:absolute;left:0;top:0;width:100%;height:2px;background:linear-gradient(90deg,transparent,rgba(34,211,238,.98),rgba(147,197,253,.95),transparent);animation:ism-page-loading-scan 1.8s ease-in-out infinite}
            .ism-body-page-loading .ism-page-loading-content{position:relative;display:flex;align-items:center;gap:18px}
            .ism-body-page-loading .ism-page-loading-orbit{position:relative;width:54px;height:54px;flex:0 0 54px;display:flex;align-items:center;justify-content:center}
            .ism-body-page-loading .ism-page-loading-spinner{position:absolute;inset:0;border:3px solid rgba(125,211,252,.18);border-top-color:#22d3ee;border-right-color:#60a5fa;border-radius:50%;box-shadow:0 0 24px rgba(34,211,238,.38);animation:ism-page-loading-spin .9s linear infinite}
            .ism-body-page-loading .ism-page-loading-core{width:18px;height:18px;border-radius:50%;background:radial-gradient(circle,#e0faff 0,#22d3ee 40%,rgba(34,211,238,.1) 70%);box-shadow:0 0 20px rgba(34,211,238,.75);animation:ism-page-loading-breathe 1.3s ease-in-out infinite}
            .ism-body-page-loading .ism-page-loading-copy{min-width:0}
            .ism-body-page-loading .ism-page-loading-title{color:#e0faff;font-size:18px;font-weight:600;line-height:26px;letter-spacing:0}
            .ism-body-page-loading .ism-page-loading-subtitle{margin-top:4px;color:rgba(224,250,255,.72);font-size:13px;line-height:19px;letter-spacing:0}
            .ism-body-page-loading .ism-page-loading-meter{position:relative;display:flex;gap:6px;margin-top:24px}
            .ism-body-page-loading .ism-page-loading-meter span{height:3px;flex:1;background:#22d3ee;box-shadow:0 0 12px rgba(34,211,238,.45);animation:ism-page-loading-pulse 1.2s ease-in-out infinite}
            .ism-body-page-loading .ism-page-loading-meter span:nth-child(2){background:#38bdf8;animation-delay:.15s}
            .ism-body-page-loading .ism-page-loading-meter span:nth-child(3){background:#60a5fa;animation-delay:.3s}
            .ism-body-page-loading .ism-page-loading-meter span:nth-child(4){background:#93c5fd;animation-delay:.45s}
            @keyframes ism-page-loading-spin{to{transform:rotate(360deg)}}
            @keyframes ism-page-loading-scan{0%{transform:translateX(-120%);opacity:0}20%,80%{opacity:1}100%{transform:translateX(120%);opacity:0}}
            @keyframes ism-page-loading-pulse{0%,100%{opacity:.45;transform:scaleX(.72)}50%{opacity:1;transform:scaleX(1)}}
            @keyframes ism-page-loading-breathe{0%,100%{transform:scale(.82);opacity:.72}50%{transform:scale(1);opacity:1}}
          `
          document.head.appendChild(style)
        }
        const el = document.createElement('div')
        el.setAttribute('data-ism-page-loading', 'true')
        el.className = 'ism-body-page-loading'
        el.innerHTML = '<div class="ism-page-loading-panel"><div class="ism-page-loading-scan"></div><div class="ism-page-loading-content"><div class="ism-page-loading-orbit"><div class="ism-page-loading-spinner"></div><div class="ism-page-loading-core"></div></div><div class="ism-page-loading-copy"><div class="ism-page-loading-title">页面加载中</div><div class="ism-page-loading-subtitle">正在渲染组件，请稍候</div></div></div><div class="ism-page-loading-meter"><span></span><span></span><span></span><span></span></div></div>'
        Object.assign(el.style, {
          pointerEvents: 'none'
        })
        this.pageLoadingEl = el
      }
      if (!document.body.contains(this.pageLoadingEl)) {
        document.body.appendChild(this.pageLoadingEl)
      }
    },
    unmountBodyPageLoading() {
      if (this.pageLoadingEl && this.pageLoadingEl.parentNode) {
        this.pageLoadingEl.parentNode.removeChild(this.pageLoadingEl)
      }
    },
    async showPage(linkInfo) {
      const _debugTag = '[showPage]'
      console.log(_debugTag, '===== ENTRY =====', JSON.parse(JSON.stringify(linkInfo)))
      //  this.PopUpDialog = false
      try {
        if(linkInfo && linkInfo.linkType == "Inside" && linkInfo.Inside && (linkInfo.Inside.displayType === 2 || linkInfo.Inside.displayType === '2')) {
          this.cancelPendingPageLoading()
          this.$router.push({
            path: `/DisPlay3DRunApp/${linkInfo.Inside.displayUUID}`,
            query: linkInfo.Inside.pageUUID ? { pageId: linkInfo.Inside.pageUUID } : {}
          })
          return
        }
        if(typeof linkInfo.isPopUp!='undefined'&& linkInfo.isPopUp==true)
        {
          let _t = this
          const requestToken = ++this.popUpPageRequestToken
          this.popUpRenderToken += 1
          this.IsAutoClose = linkInfo.autoClose
          if (linkInfo.linkType == "Inside") {
            if (this.PopUpDialog &&
                this.ISMPopUpRunningContainer &&
                this.currentPopUpDisplayUUID === linkInfo.Inside.displayUUID &&
                this.currentPopUpPageUUID === linkInfo.Inside.pageUUID) {
              return
            }
            // 先标记目标页，防止异步加载期间重复触发 GoPage
            this.currentPopUpDisplayUUID = linkInfo.Inside.displayUUID
            this.currentPopUpPageUUID = linkInfo.Inside.pageUUID
            let page = {
              pageType: 1,
              displayUUID: linkInfo.Inside.displayUUID,
              pageUuid: linkInfo.Inside.pageUUID
            }
            _t.isExternUrl = false
            _t.chargePagePopUp = true
            const loadingToken = await this.consumePendingPageLoading()
            if (requestToken !== _t.popUpPageRequestToken || _t._isDestroyed) {
              _t.closePageLoading(loadingKey, loadingToken)
              return
            }
            _t.destroyPopUpGraph()
            this.CurrentPopRealUUIDList=[]
            this.CurrentPopModelUUIDList=[]
            this.CurrentPagerPopRealDataUuidList=[]
            this.CurrentPagerPopRealDeviceUuidList=[]
            this.selectPopUpDisplayPageDataStruct({
              page: page, callback: function (res,uuids,devices) {
                console.log('[showPage-PopUp] callback fired, res=', res, 'isFound/effective=', res===0)
                if (requestToken !== _t.popUpPageRequestToken || _t._isDestroyed) {
                  _t.closePageLoading(loadingKey, loadingToken)
                  console.warn('[showPage-PopUp] callback ABORTED: token mismatch or destroyed')
                  return
                }
                if(res==0)
                {
                  console.log('[showPage-PopUp] page loaded, opening dialog')
                  _t.chargePagePopUp = false
                  _t.PopUpDialog = true
                  _t.CurrentPagerPopRealDataUuidList=uuids
                  _t.CurrentPagerPopRealDeviceUuidList=devices
                  _t.$nextTick(function () {
                    _t.trackTimeout(function (){
                      console.log('[showPage-PopUp] delayed init fired')
                      if (requestToken !== _t.popUpPageRequestToken || _t._isDestroyed) {
                        _t.closePageLoading(loadingKey, loadingToken)
                        console.warn('[showPage-PopUp] delayed init ABORTED')
                        return
                      }
                      _t.destroyPopUpGraph()
                      console.log('[showPage-PopUp] calling RunPopUpCavasContainerInit')
                        _t.RunPopUpCavasContainerInit(_t.PopUpConfigData.layer.autoSize,_t.PopUpConfigData.layer.Padding, loadingToken)
                    },100, 'popup')
                  })

                  // _t.PopUpDialog.center()

                }
                else
                {
                  _t.currentPopUpDisplayUUID = ""
                  _t.currentPopUpPageUUID = ""
                  _t.chargePagePopUp = false
                  _t.closePageLoading(loadingKey, loadingToken);
                  _t.$message.error(_t.$t("readData.NotFindPage"))
                }
              }
            })
          }
          else
          {
            this.cancelPendingPageLoading()
            _t.chargePagePopUp = false
            _t.isExternUrl = true
            _t.currentPopUpDisplayUUID = ""
            _t.currentPopUpPageUUID = ""
            _t.destroyPopUpGraph()
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
            // 已在当前页则跳过，防止 GoPage 高频推送导致反复销毁重建
            if (this.currentDisplayUUID === linkInfo.Inside.displayUUID &&
                this.currentPageUUID === linkInfo.Inside.pageUUID) {
              this.cancelPendingPageLoading()
              console.log('[showPage] skip: already on target page', linkInfo.Inside)
              return
            }
            // 先标记目标页，防止异步加载期间重复触发 GoPage
            this.currentDisplayUUID = linkInfo.Inside.displayUUID
            this.currentPageUUID = linkInfo.Inside.pageUUID
            let _t = this
            const requestToken = ++this.mainPageRequestToken
            this.pageRenderToken += 1
            _t.chargePage = true
            const loadingToken = await this.consumePendingPageLoading()
            if (requestToken !== _t.mainPageRequestToken || _t._isDestroyed) {
              _t.closePageLoading(loadingKey, loadingToken)
              return
            }
            // 立即销毁旧 Graph，阻断旧实例继续消费数据/定时器
            _t.destroyMainGraph()
            _t.clearPendingTimers('main')
            // 标记正在切换，阻止数据推送期间无效的 DealWithUpdateData
            let page = {
              pageType: 1,
              displayUUID: linkInfo.Inside.displayUUID,
              pageUuid: linkInfo.Inside.pageUUID
            }

            _t.CurrentRealUUIDList=[]
            _t.CurrentModelUUIDList=[]
            this.CurrentPagerRealDataUuidList=[]
            this.CurrentPagerRealDeviceUuidList=[]

            this.selectDisplayPageDataStruct({
              page: page, callback: function (uuids,devices,isFound) {
                if (requestToken !== _t.mainPageRequestToken || _t._isDestroyed) {
                  _t.closePageLoading(loadingKey, loadingToken)
                  console.warn("[showPage] callback ABORTED: token mismatch or destroyed", "req=", requestToken, "cur=", _t.mainPageRequestToken, "destroyed=", _t._isDestroyed)
                  return
                }
                console.log('[showPage-Main] callback fired, isFound=', isFound, 'uuids=', uuids?.length, 'devices=', devices?.length)
                _t.chargePage = false
                if (isFound === false) {
                  _t.closePageLoading(loadingKey, loadingToken)
                  console.error('[showPage-Main] page not found, displayUUID=', page.displayUUID)
                  _t.currentDisplayUUID = ""
                  _t.currentPageUUID = ""
                  _t.$message.error(_t.$t("readData.NotFindPage"))
                  return
                }
                // 设置 currentDisplayUUID 会触发 watch，但这是原有逻辑，保留
                console.log('[showPage-Main] setting currentDisplayUUID=', linkInfo.Inside.displayUUID, 'currentPageUUID=', linkInfo.Inside.pageUUID)
                _t.currentDisplayUUID = linkInfo.Inside.displayUUID
                _t.currentPageUUID = linkInfo.Inside.pageUUID
                _t.CurrentPagerRealDataUuidList = uuids
                _t.CurrentPagerRealDeviceUuidList = devices
                console.log('[showPage-Main] calling RunCavasContainerInit')
                _t.RunCavasContainerInit(_t.configData.layer.autoSize,_t.configData.layer.Padding, loadingToken)
              }
            })
          }
          else if (linkInfo.linkType == "External") {
            this.cancelPendingPageLoading()
            if(typeof linkInfo.OpenExternalType=='undefined') {
              window.open(linkInfo.External, '_blank')
            }else if( linkInfo.OpenExternalType=='self'){
              window.open(linkInfo.External, '_self')
            }else if(linkInfo.OpenExternalType=='new'){
              window.open(linkInfo.External, '_blank')
            }
          }
        }
      }catch (e) {
        this.chargePage = false
        this.chargePagePopUp = false
        this.closePageLoading(loadingKey)
        console.error(e)
      }finally {
        console.log("[showPage] finally - nothing to clean here")
      }
    },
    initPinchZoom() {
      try {
        const CONFIG = {
          MIN_ZOOM: 0.1,
          MAX_ZOOM: 3,
          DEBOUNCE_TIME: 30,
          PINCH_THRESHOLD: 0.1,
        };

        const container = this.$refs.ISMRunningContainer; // DOM 容器
        const graph = this.ISMRuningCavasContainer; // X6 Graph 实例
        if (!container || !graph) {
          console.error('容器或 X6 实例未找到');
          return;
        }

        // 手动维护缩放比例（初始值 1，不受 X6 API 影响）
        this.destroyHammerManager()
        this.currentAbsoluteZoom = this.currentAbsoluteZoom || 1;
        let lastScale = 1;
        const containerRect = container.getBoundingClientRect();

        // 坐标转换工具函数
        const convertToContainerCoord = (clientX, clientY) => ({
          x: clientX - containerRect.left,
          y: clientY - containerRect.top,
        });

        // 初始化 Hammer
        const mc = new Hammer.Manager(container, {
          touchAction: 'manipulation',
          recognizers: [[Hammer.Pinch, { pointers: 2, threshold: CONFIG.PINCH_THRESHOLD }]]
        });

        // 捏合事件
        mc.on('pinch', (ev) => {
          const relativeScale = ev.scale / lastScale;
          lastScale = ev.scale;

          // 计算新的绝对缩放比例（手动维护，不依赖 X6 API）
          let newZoom = this.currentAbsoluteZoom * relativeScale;
          newZoom = Math.max(CONFIG.MIN_ZOOM, Math.min(newZoom, CONFIG.MAX_ZOOM));

          // 缩放中心点
          const center = convertToContainerCoord(ev.center.x, ev.center.y);

          // 防抖执行缩放
          clearTimeout(this.zoomDebounceTimer);
          this.zoomDebounceTimer = setTimeout(() => {
            // X6 缩放：v1.x 和 v2.x 都支持的调用方式
            graph.zoom(newZoom, {
              x: center.x,
              y: center.y,
              absolute: true, // 关键：绝对缩放，基于新比例
            });

            // 更新手动维护的缩放比例
            this.currentAbsoluteZoom = newZoom;
          }, CONFIG.DEBOUNCE_TIME);
        });

        // 捏合结束重置
        mc.on('pinchend pinchcancel', () => {
          lastScale = 1;
          clearTimeout(this.zoomDebounceTimer);
          this.zoomDebounceTimer = null
        });

        this.hammerManager = mc;
      } catch (error) {
        console.error('捏合缩放初始化异常：', error);
      }
    },
    RunCavasContainerInit(auto,panning,loadingToken = this.pageLoadingToken){
      let _t = this
      let pagerSize = auto==1?true:false
      let pagerPanning = panning==1?true:false
      this.rebuildRef=true
      // 取消上一次 destroyMainGraph 残留的 RAF 清理任务，避免它在新 Graph 渲染完成后清空 DOM
      if (this.destroyMainGraphRaf) {
        cancelAnimationFrame(this.destroyMainGraphRaf)
        this.destroyMainGraphRaf = null
        console.log('[RunCavasContainerInit] cancelled pending destroyMainGraph RAF')
      }
      // 无论是否复用 Graph，都先自增 token，作废上一轮所有异步回调（render:done / trackTimeout）
      this.pageRenderToken += 1
      if(this.tempAutoSizePager!=pagerSize)
      {
        this.tempAutoSizePager = pagerSize
        if(this.ISMRuningCavasContainer!=null)
        {
          try {
            this.destroyHammerManager()
            this.ISMRuningCavasContainer.clearCells();
            // 移除事件
            this.ISMRuningCavasContainer.off();
            this.ISMRuningCavasContainer.dispose()
          } catch(e) {
            console.error('[RunCavasContainerInit] dispose failed', e)
          }
          this.ISMRuningCavasContainer=null
        }
        this.rebuildRef=false
      }

      let width = _t.configData.layer.width
      let height = _t.configData.layer.height
      this.registerISMGroupNode()
      if(pagerSize)
      {
        width = window.innerWidth
        height = window.innerHeight
      }
      if(this.ISMRuningCavasContainer!=null)
      {
        console.log("[RunCavasContainerInit] reuse existing Graph instance")
        this.ISMRuningCavasContainer.resize(width,height)
        this.ISMRuningCavasContainer.drawBackground ( {
          color: _t.configData.layer.backColor,   // 背景底色（可选）
          image: _t.configData.layer.backgroundImage,
          size: '100% 100%',
          repeat: 'no-repeat',
          quality:1,
        });

        if(pagerPanning)
        {
          this.ISMRuningCavasContainer.enablePanning()
          this.ISMRuningCavasContainer.enableMouseWheel()
        }
        else
        {
          this.ISMRuningCavasContainer.disablePanning()
          this.ISMRuningCavasContainer.disableMouseWheel()
        }
      }
      else{
        console.log("[RunCavasContainerInit] create new Graph instance")
        this.ISMRuningCavasContainer = new Graph({
          container: this.$refs.ISMRunningContainer,
          width: _t.configData.layer.width,
          height: _t.configData.layer.height,
          panning: pagerPanning,
          mousewheel: pagerPanning,
          autoResize:false,
          grid: false,
          virtual:false,
          async:false,
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
            size: '100% 100%',
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
        this.initPinchZoom()
        this.ISMRuningCavasContainer.use(  new Keyboard({
                enabled: true,
                global: true,
              }),
        )
        if(pagerPanning)
        {
          this.ISMRuningCavasContainer.enablePanning()
          this.ISMRuningCavasContainer.enableMouseWheel()
        }
        else
        {
          this.ISMRuningCavasContainer.disablePanning()
          this.ISMRuningCavasContainer.disableMouseWheel()
        }
      }

      this.tempAutoSize=this.configData.layer.autoSize
      if(pagerSize)
      {
        this.setScale()
      }
      else
      {

        this.ISMRuningCavasContainer.resize(width,height)
        this.ISMRuningCavasContainer.scale( 1,  1);
      }
      console.log("[RunCavasContainerInit] Graph initialized with size:", width, height, "autoSize:", pagerSize, "panning:", pagerPanning)
      this.ISMRuningCavasContainer.clearCells()
      this.ISMRuningCavasContainer.off('node:click')
      this.ISMRuningCavasContainer.off('node:dblclick')
      this.ISMRuningCavasContainer.off('node:mousedown')
      this.ISMRuningCavasContainer.off('node:mouseup')
      this.ISMRuningCavasContainer.off('node:mouseenter')
      this.ISMRuningCavasContainer.off('node:mouseleave')
      this.ISMRuningCavasContainer.off('render:done')
      this.ISMRuningCavasContainer.off('blank:contextmenu')
      this.ISMRuningCavasContainer.off('cell:contextmenu')
      this.ISMRuningCavasContainer.off('scale')

      window.removeEventListener('resize', this.setScale);
      window.addEventListener('resize', this.setScale);
      // token 已在函数入口处自增，此处只捕获当前值供回调使用
      const renderToken = this.pageRenderToken
      this.ISMRuningCavasContainer.on('scale', ( scale ) => {

      });
      //
      this.ISMRuningCavasContainer.on('node:click', ( { e, x, y, node, view } ) => {
        _t.handleGraphNodeAction(node, 'click')
      });
      this.ISMRuningCavasContainer.on('node:dblclick', ( { e, x, y, node, view } ) => {
        _t.handleGraphNodeAction(node, 'dblclick')
      });
      this.ISMRuningCavasContainer.on('node:mousedown', ( { e, x, y, node, view } ) => {
        _t.handleGraphNodeAction(node, 'mousedown')
      });
      this.ISMRuningCavasContainer.on('node:mouseup', ( { e, x, y, node, view } ) => {
        _t.handleGraphNodeAction(node, 'mouseup')
      });
      this.ISMRuningCavasContainer.on('node:mouseenter', ( { e, x, y, node, view } ) => {
        _t.handleGraphNodeAction(node, 'mouseenter')
      });
      this.ISMRuningCavasContainer.on('node:mouseleave', ( { e, x, y, node, view } ) => {
        _t.handleGraphNodeAction(node, 'mouseleave')
      });
      //页面空白处右键菜单
      this.ISMRuningCavasContainer.on('blank:contextmenu', ({ e, x, y, cell, view }) => {
        this.onContextLayerMenu(e)
      })
      this.ISMRuningCavasContainer.on('cell:contextmenu', ({ e, x, y, cell, view }) => {
        this.onContextLayerMenu(e)
      })

      this.ISMRuningCavasContainer.bindKey('ctrl+l', () => {
        this.lockScreen()
        return false
      })
      //==============================
      //监听页面是否渲染完成，完成后就发广播给节点，告诉节点是在编辑模式
      this.ISMRuningCavasContainer.on('render:done', () => {
        console.log("[RunCavasContainerInit] render done")  
        const domContainer = this.$refs.ISMRunningContainer
        console.log('[RunCavasContainerInit] DOM check at render:done: container exists=', !!domContainer, 'childCount=', domContainer?.children?.length, 'svg exists=', !!domContainer?.querySelector('svg'))
        if (renderToken !== this.pageRenderToken) {
          console.warn('[RunCavasContainerInit] render:done ABORTED: token mismatch, renderToken=', renderToken, 'pageRenderToken=', this.pageRenderToken)
          return
        }
         console.log("[RunCavasContainerInit] render done33333333")  
        this.tempAutoSize=this.configData.layer.autoSize
        // this.setScale()
        this.$EventBus.$emit('cell-editMode',{
          edit:false,
          toolbox:false
        })
        this.$EventBus.$emit('cell-vuex',{
          PMapState:mapState,
          PMapActions:mapActions,
          PMapMutations:mapMutations,
          PStore:this.$store,
          PRouter:this.$router
        })
        this.closePageLoading(loadingKey, loadingToken)
        this.InitPagerRealData()
      });
      try{
        console.log("[RunCavasContainerInit] loading components from JSON, cells count=", _t.configData.components?.cells?.length)
        const hasObserver = !!_t.configData.components?.__ob__
        console.log("[RunCavasContainerInit] components has __ob__ (Vue Observer)=", hasObserver)
        const components = JSON.parse(JSON.stringify(_t.configData.components))
        if (components.cells && Array.isArray(components.cells)) {
          const originalCount = components.cells.length
          components.cells = components.cells.filter(cell => cell && cell.shape)
          const filteredCount = components.cells.length
          if (filteredCount < originalCount) {
            console.warn(`[RunCavasContainerInit] filtered ${originalCount - filteredCount} cells that are missing the shape property`)
          }
        }
        _t.ISMRuningCavasContainer.fromJSON(components)
        console.log("[RunCavasContainerInit] fromJSON completed, graph cells count=", _t.ISMRuningCavasContainer.getCells()?.length)
        this.$nextTick(() => {
          const raf = window.requestAnimationFrame || function (cb) { return setTimeout(cb, 0) }
          raf(() => this.closePageLoading(loadingKey, loadingToken))
        })
      }catch (e){
        _t.closePageLoading(loadingKey, loadingToken)
        _t.$message.error(_t.$t('Render.GetPageError'))
        console.error('[RunCavasContainerInit] fromJSON ERROR:', e)
      }
    },
    RunPopUpCavasContainerInit(auto,panning,loadingToken = this.pageLoadingToken){
      let _t = this
      if (!this.PopUpDialog || !this.$refs.ISMPopUpRunningContainer) {
        this.$nextTick(function () {
          _t.trackTimeout(function () {
            if (_t.PopUpDialog && !_t._isDestroyed) {
              _t.RunPopUpCavasContainerInit(auto, panning, loadingToken)
            }
          }, 50, 'popup')
        })
        return
      }
      let pagerSize = auto==1?true:false
      let pagerPanning = panning==1?true:false
      if(this.ISMPopUpRunningContainer!=null)
      {
        try {
          this.ISMPopUpRunningContainer.clearCells();
          this.ISMPopUpRunningContainer.off();
          this.ISMPopUpRunningContainer.dispose();
        } catch(e) {
          console.error('[RunPopUpCavasContainerInit] dispose failed', e)
        }
        this.ISMPopUpRunningContainer = null
      }
      this.registerISMGroupNode()
      const renderToken = ++this.popUpRenderToken
      this.ISMPopUpRunningContainer = new Graph({
        container: this.$refs.ISMPopUpRunningContainer,
        width: _t.PopUpConfigData.layer.width,
        height: _t.PopUpConfigData.layer.height,
        panning: pagerPanning,
        mousewheel: pagerPanning,
        autoResize:pagerSize,
        grid: false,
        virtual:false,
        async:false,
        background: {
          color: _t.PopUpConfigData.layer.backColor,   // 背景底色（可选）
          image: _t.PopUpConfigData.layer.backgroundImage,
          size: '100% 100%',
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
      this.ISMPopUpRunningContainer.on('scale', ( scale ) => {

      });
      //
      this.ISMPopUpRunningContainer.on('node:click', ( { e, x, y, node, view } ) => {
        _t.handleGraphNodeAction(node, 'click')
      });
      this.ISMPopUpRunningContainer.on('node:dblclick', ( { e, x, y, node, view } ) => {
        _t.handleGraphNodeAction(node, 'dblclick')
      });
      this.ISMPopUpRunningContainer.on('node:mousedown', ( { e, x, y, node, view } ) => {
        _t.handleGraphNodeAction(node, 'mousedown')
      });
      this.ISMPopUpRunningContainer.on('node:mouseup', ( { e, x, y, node, view } ) => {
        _t.handleGraphNodeAction(node, 'mouseup')
      });
      this.ISMPopUpRunningContainer.on('node:mouseenter', ( { e, x, y, node, view } ) => {
        _t.handleGraphNodeAction(node, 'mouseenter')
      });
      this.ISMPopUpRunningContainer.on('node:mouseleave', ( { e, x, y, node, view } ) => {
        _t.handleGraphNodeAction(node, 'mouseleave')
      });
      //页面空白处右键菜单
      this.ISMPopUpRunningContainer.on('blank:contextmenu', ({ node }) => {
        this.onContextLayerMenu()
      })
      this.ISMPopUpRunningContainer.on('render:done', () => {
        if (renderToken !== this.popUpRenderToken || this._isDestroyed) {
          this.closePageLoading(loadingKey, loadingToken)
          return
        }
        this.closePageLoading(loadingKey, loadingToken)
        this.$EventBus.$emit('cell-editMode',{
          edit:false,
          toolbox:false
        })
        this.$EventBus.$emit('cell-vuex',{
          PMapState:mapState,
          PMapActions:mapActions,
          PMapMutations:mapMutations,
          PStore:this.$store,
          PRouter:this.$router
        })
        this.InitPopUpPagerRealData()
      });
      //==============================
      try{
        const components = JSON.parse(JSON.stringify(_t.PopUpConfigData.components))
        if (components.cells && Array.isArray(components.cells)) {
          components.cells = components.cells.filter(cell => cell && cell.shape)
        }
        _t.ISMPopUpRunningContainer.fromJSON(components)
      }catch (e){
        _t.closePageLoading(loadingKey, loadingToken)
        _t.$message.error(_t.$t('Render.GetPageError'))
        console.log(e)
      }
    }
  },
  mounted() {
    this.setPasswordDialog = false
    this.settingDialog = false
    this.PopUpDialog = false
    let _t = this
    let _this = this
    this.$nextTick(function(){
      _t.loadPager()
      _t.eventHandlers.readDataPush && _t.$EventBus.$off("readDataPush", _t.eventHandlers.readDataPush)
      _t.eventHandlers.onSelectDevice && _t.$EventBus.$off("onSelectDevice", _t.eventHandlers.onSelectDevice)
      _t.eventHandlers.SystemData && _t.$EventBus.$off("SystemData", _t.eventHandlers.SystemData)
      _t.eventHandlers.StaticData && _t.$EventBus.$off("StaticData", _t.eventHandlers.StaticData)
      _t.eventHandlers.RealAlarm && _t.$EventBus.$off("RealAlarm", _t.eventHandlers.RealAlarm)
      _t.eventHandlers.ChargePage && _t.$EventBus.$off("ChargePage", _t.eventHandlers.ChargePage)

      _t.eventHandlers.PlayVoice && _t.$EventBus.$off("PlayVoice", _t.eventHandlers.PlayVoice)
      _t.eventHandlers.PlayVoice = (data) => {
        let wsData = data
        _this.alarmSoundSpeech().start({container:"#sppekContent",Lang:_this.$i18n.locale,rate:1},wsData.VoiceString);
      }
      _t.$EventBus.$on("PlayVoice", _t.eventHandlers.PlayVoice);

      _t.eventHandlers.GoPage && _t.$EventBus.$off("GoPage", _t.eventHandlers.GoPage)
      _t.eventHandlers.GoPage = (data) => {
        if(!_t.shouldHandleGoPage(data)) {
          console.warn("[GoPage] BLOCKED by shouldHandleGoPage (800ms debounce)")
          return
        }
        // 非弹窗切换时，若主页面正在加载中则直接忽略，防止高频点击叠加请求
        if(_t.chargePage && !(data.IsPopUp)) {
          console.warn("[GoPage] BLOCKED by chargePage guard", "chargePage=", _t.chargePage, "IsPopUp=", data.IsPopUp)
          return
        }
        let JumpWindow = localStorage.getItem("JumpWindow")
        let JumpWindowEnable = false
        if((JumpWindow=="null")||(JumpWindow==null)||(JumpWindow==""))
        {
          JumpWindowEnable=true
        }
        else
        {
          try{
            JumpWindow = JSON.parse(JumpWindow)
            JumpWindowEnable = JumpWindow.enable
          }catch (e)
          {
            JumpWindowEnable=true
          }

        }
        if(JumpWindowEnable) {
          let wsData = data
          let linkInfo = {
            isPopUp: wsData.IsPopUp,
            autoClose: wsData.AutoClose,
            linkType: typeof wsData.linkType != "undefined" ? wsData.linkType : "Inside",
            Inside: {
              displayUUID: wsData.ModelId,
              pageUUID: wsData.PageUuid,
            },
            width: wsData.width,
            height: wsData.height,
            External: wsData.External,
            title: wsData.External,
            OpenExternalType: wsData.OpenExternalType
          }
          _t.showPage(linkInfo)
        } else {
          console.warn("[GoPage] BLOCKED by JumpWindowEnable=false")
        }
      }
      _t.$EventBus.$on("GoPage", _t.eventHandlers.GoPage);

      _t.eventHandlers.onSelectDevice && _t.$EventBus.$off("onSelectDevice", _t.eventHandlers.onSelectDevice)
      _t.eventHandlers.onSelectDevice = async (device)=>{
        if(device.type==1)
        {
          _t.SelectDeviceUuid = device.key

          let page={
            pageType:1,
            displayUUID:device.showUUID,
            pageUuid:device.showPageUUID
          }
          if((typeof device.isPopUp!='undefined')&&(device.isPopUp))
          {
            _t.chargePagePopUp = true
            const loadingToken = await _t.consumePendingPageLoading()
            // _t.$message.loading({content: 'Loading...', key: 'updatable', duration: 0});
            _t.selectPopUpDisplayPageDataStruct({
              page: page, callback: function (res,uuids,devices) {
                _t.chargePagePopUp = false
                if(res==0)
                {
                  _t.PopUpDialog = true
                  _t.CurrentPagerPopRealDataUuidList=uuids
                  _t.CurrentPagerPopRealDeviceUuidList.push(_t.SelectDeviceUuid)
                  _t.$nextTick(function () {
                    _t.trackTimeout(function (){
                      _t.destroyPopUpGraph()
                      _t.RunPopUpCavasContainerInit(_t.PopUpConfigData.layer.autoSize,_t.PopUpConfigData.layer.Padding, loadingToken)
                    },100, 'popup')
                  })
                }
                else{
                  _t.closePageLoading(loadingKey, loadingToken)
                  _t.$message.error(_t.$t('Render.GetPageError'))
                }
              }
            })
          }
          else {
            _t.currentDisplayUUID = device.showUUID
            _t.currentPageUUID = device.showPageUUID
            const requestToken = ++this.mainPageRequestToken
            _t.pageRenderToken += 1
            _t.chargePage = true
            const loadingToken = await _t.consumePendingPageLoading()
            // 立即销毁旧 Graph，阻断旧实例继续消费数据/定时器
            _t.destroyMainGraph()
            _t.clearPendingTimers('main')
            // 标记正在切换，阻止数据推送期间无效的 DealWithUpdateData
            // _t.$message.loading({content: 'Loading...', key: 'updatable', duration: 0});
            _t.CurrentRealUUIDList=[]
            _t.CurrentModelUUIDList=[]
            _t.CurrentPagerRealDataUuidList=[]
            _t.CurrentPagerRealDeviceUuidList=[]
            _t.selectDisplayPageDataStruct({
              page: page, callback: function (uuids,devices,isFound) {
                if (requestToken !== _t.mainPageRequestToken || _t._isDestroyed) {
                  _t.closePageLoading(loadingKey, loadingToken)
                  return
                }
                _t.chargePage = false
                if (isFound === false) {
                  _t.closePageLoading(loadingKey, loadingToken)
                  _t.currentDisplayUUID = ""
                  _t.currentPageUUID = ""
                  _t.$message.error(_t.$t("readData.NotFindPage"))
                  return
                }
                _t.currentDisplayUUID = page.displayUUID
                _t.currentPageUUID = page.pageUuid
                _t.CurrentPagerRealDataUuidList = uuids
                _t.CurrentPagerRealDeviceUuidList[0]=_t.SelectDeviceUuid
                _t.RunCavasContainerInit(_t.configData.layer.autoSize,_t.configData.layer.Padding, loadingToken)
              }
            })
          }
        }
      }
      _t.$EventBus.$on("onSelectDevice", _t.eventHandlers.onSelectDevice);
      _t.eventHandlers.readDataPush = (data) => {
        if(!_t.chargePage) {
          _t.DealWithUpdateData(data)
        }
      }
      _t.$EventBus.$on("readDataPush", _t.eventHandlers.readDataPush);
      _t.eventHandlers.StaticData = (data) => {
        if(!_t.chargePage) {
          _t.DealWithUpdateData(data)
        }
      }
      _t.$EventBus.$on("StaticData", _t.eventHandlers.StaticData);
      _t.eventHandlers.SystemData = (data) => {

        if(!_t.chargePage) {
          let realData = data
          let projectUuid = getAuthorization(AUTH_TYPE.AUTH1)
          if (typeof _t.configData.components === "undefined" || !_t.configData.components.cells) {
            return;
          }
          const cells = _t.configData.components.cells;
          if (!Array.isArray(cells) || cells.length === 0) {
            return;
          }
          for (let k = 0; k < realData.Data.length; k++) {
            for (let j = 0; j < _t.configData.components.cells.length; j++) {
              if ((typeof _t.configData.components.cells[j].data.detail.animate != "undefined") && (typeof _t.configData.components.cells[j].data.detail.animate.condition != "undefined")) {
                let condition = _t.configData.components.cells[j].data.detail.animate.condition
                let isExpression = _t.configData.components.cells[j].data.detail.animate.isExpression
                //动画
                if ((typeof condition != "undefined") && (realData.ProjectUuid == projectUuid)) {
                  let isStart = false
                  if (condition.dataID == realData.Data[k].Uuid) {
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
                      }
                      _t.$EventBus.$emit(_t.configData.components.cells[j].data.detail.identifier + "animateEvent", isStart);
                    }
                  }
                }
              }
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
              if ((typeof _t.configData.components.cells[j].data.detail.active != "undefined")) {
                let active = _t.configData.components.cells[j].data.detail.active
                //动作
                if (typeof active != "undefined") {
                  let tempResult = ""
                  for (let activeIndex = 0; activeIndex < active.length; activeIndex++) {
                    if (((realData.ProjectUuid == projectUuid))) {
                      if (active[activeIndex].condition.dataID == realData.Data[k].Uuid) {
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
                          }
                        } else {
                          tempResult = realData.Data[k].Value
                        }
                        active[activeIndex].result = tempResult
                        let activeData = {
                          ID: active[activeIndex].id,
                          dataID: active[activeIndex].condition.dataID,
                          index: activeIndex,
                          result: tempResult
                        }
                        _t.$EventBus.$emit(_t.configData.components.cells[j].data.detail.identifier + "activeEvent", activeData);
                      }
                    }
                  }
                }
              }
            }
          }
        }

        if(!_t.chargePagePopUp) {
          let realData = data
          let projectUuid = getAuthorization(AUTH_TYPE.AUTH1)
          if (typeof _t.PopUpConfigData.components === "undefined" || !_t.PopUpConfigData.components.cells) {
            return;
          }
          const cells = _t.PopUpConfigData.components.cells;
          if (!Array.isArray(cells) || cells.length === 0) {
            return;
          }
          for (let k = 0; k < realData.Data.length; k++) {
            for (let j = 0; j < _t.PopUpConfigData.components.cells.length; j++) {
              if ((typeof _t.PopUpConfigData.components.cells[j].data.detail.animate != "undefined") && (typeof _t.PopUpConfigData.components.cells[j].data.detail.animate.condition != "undefined")) {
                let condition = _t.PopUpConfigData.components.cells[j].data.detail.animate.condition
                let isExpression = _t.PopUpConfigData.components.cells[j].data.detail.animate.isExpression
                //动画
                if ((typeof condition != "undefined") && (realData.ProjectUuid == projectUuid)) {
                  let isStart = false
                  if (condition.dataID == realData.Data[k].Uuid) {
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

              if ((typeof _t.PopUpConfigData.components.cells[j].data.detail.active != "undefined")) {
                let active = _t.PopUpConfigData.components.cells[j].data.detail.active
                //动作
                if (typeof active != "undefined") {
                  let tempResult = ""
                  for (let activeIndex = 0; activeIndex < active.length; activeIndex++) {
                    if (((realData.ProjectUuid == projectUuid))) {
                      if (active[activeIndex].condition.dataID == realData.Data[k].Uuid) {
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
                          }
                        } else {
                          tempResult = realData.Data[k].Value
                        }
                        active[activeIndex].result = tempResult
                        let activeData = {
                          ID: active[activeIndex].id,
                          dataID: active[activeIndex].condition.dataID,
                          index: activeIndex,
                          result: tempResult
                        }
                        _t.$EventBus.$emit(_t.PopUpConfigData.components.cells[j].data.detail.identifier + "activeEvent", activeData);
                      }
                    }
                  }
                }
              }
            }
          }
        }
      };
      _t.$EventBus.$on("SystemData", _t.eventHandlers.SystemData);
      _t.eventHandlers.RealAlarm = (data) => {
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
      };
      _t.$EventBus.$on("RealAlarm", _t.eventHandlers.RealAlarm);
      _t.eventHandlers.ChargePage = (data) => {
        let linkInfo = {
          isPopUp:data.IsPopUp,
          linkType:"Inside",
          Inside:{}
        }
        if(data.PageID!="") {
          if(data.DisPlayID!=undefined&&data.DisPlayID!="")
          {
            linkInfo.Inside.displayUUID = data.DisPlayID
          }
          else {
            linkInfo.Inside.displayUUID = this.$route.params.uid
          }
          linkInfo.Inside.pageUUID = data.PageID
          _t.showPage(linkInfo)
        }
      };
      _t.$EventBus.$on("ChargePage", _t.eventHandlers.ChargePage);


    });
  },
  beforeDestroy() {
    speechSynthesis.cancel()
    this.$notification.destroy()
    this.clearPendingTimers()
    this.mainPageRequestToken += 1
    this.popUpPageRequestToken += 1

    // === 清理子组件 EventBus 泄漏（全量 $off） ===
    // cell-editMode / cell-vuex：160+ 子组件 $on 但从未 $off

    // 动态事件名：identifier+"activeEvent" / identifier+"animateEvent"
    // 遍历 Graph Cell 数据，收集所有 identifier 并批量 $off

    // === 清理自身 EventBus 监听（精确 $off） ===
    const h = this.eventHandlers
    this.$EventBus.$off("readDataPush", h.readDataPush)
    this.$EventBus.$off("onSelectDevice", h.onSelectDevice)
    this.$EventBus.$off("SystemData", h.SystemData)
    this.$EventBus.$off("StaticData", h.StaticData)
    this.$EventBus.$off("RealAlarm", h.RealAlarm)
    this.$EventBus.$off("ChargePage", h.ChargePage)
    this.$EventBus.$off("PlayVoice", h.PlayVoice)
    this.$EventBus.$off("GoPage", h.GoPage)
    this.$EventBus.$off("cell-editMode")
    this.$EventBus.$off("cell-vuex")
    this.$EventBus.$off("MenuConfigPage")
    this.$EventBus.$off("onContainerSelectDevice")
    this.$EventBus.$off("DealWithRealDataFinish")
    this.eventHandlers = {}

    // 清理 X6 Graph（已有内部 try/catch/finally 保护）
    // destroyMainGraph 和 destroyPopUpGraph 内部已完整处理：
    // clearCells() → off() → dispose() → container.innerHTML='' → 置 null
    this.destroyMainGraph()
    this.destroyPopUpGraph()
    this.unmountBodyPageLoading()

    // 清理 DOM 事件监听器
    window.removeEventListener('resize', this.setScale);
    this.removeListener();
  }
}
</script>
<style>
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

@keyframes ant-line-forward {
  to { stroke-dashoffset: -1000; }
}
@keyframes ant-line-inverse {
  to { stroke-dashoffset: 1000; }
}

</style>

<style lang="less" scoped>
::v-deep .x6-edge-tool-vertex {
  display: none !important;
}
/* 设置可点击节点的鼠标样式 */
::v-deep .x6-node:hover {
  cursor: var(--NodeMouseStyle);
}
.ism-render-root {
  width: 100%;
  height: 100%;
}
.ism-page-loading {
  position: fixed;
  inset: 0;
  z-index: 2147483647;
  display: flex;
  align-items: center;
  justify-content: center;
  background:
    radial-gradient(circle at 50% 42%, rgba(14, 116, 144, 0.26), transparent 32%),
    linear-gradient(180deg, rgba(2, 6, 23, 0.52), rgba(2, 6, 23, 0.7));
  pointer-events: none;
}
.ism-page-loading-panel {
  position: relative;
  width: 310px;
  padding: 28px 30px 24px;
  border: 1px solid rgba(103, 232, 249, 0.56);
  background: linear-gradient(180deg, rgba(9, 26, 46, 0.94), rgba(4, 13, 27, 0.9));
  box-shadow:
    0 24px 70px rgba(0, 0, 0, 0.48),
    0 0 36px rgba(34, 211, 238, 0.24),
    inset 0 1px 0 rgba(255, 255, 255, 0.14);
  backdrop-filter: blur(14px);
  overflow: hidden;
}
.ism-page-loading-panel::before {
  content: "";
  position: absolute;
  inset: 10px;
  border: 1px solid rgba(125, 211, 252, 0.14);
  pointer-events: none;
}
.ism-page-loading-scan {
  position: absolute;
  left: 0;
  top: 0;
  width: 100%;
  height: 2px;
  background: linear-gradient(90deg, transparent, rgba(34, 211, 238, 0.98), rgba(147, 197, 253, 0.95), transparent);
  animation: ism-page-loading-scan 1.8s ease-in-out infinite;
}
.ism-page-loading-content {
  position: relative;
  display: flex;
  align-items: center;
  gap: 18px;
}
.ism-page-loading-orbit {
  position: relative;
  width: 54px;
  height: 54px;
  flex: 0 0 54px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.ism-page-loading-spinner {
  position: absolute;
  inset: 0;
  border: 3px solid rgba(125, 211, 252, 0.18);
  border-top-color: #22d3ee;
  border-right-color: #60a5fa;
  border-radius: 50%;
  box-shadow: 0 0 24px rgba(34, 211, 238, 0.38);
  animation: ism-page-loading-spin 0.9s linear infinite;
}
.ism-page-loading-core {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: radial-gradient(circle, #e0faff 0, #22d3ee 40%, rgba(34, 211, 238, 0.1) 70%);
  box-shadow: 0 0 20px rgba(34, 211, 238, 0.75);
  animation: ism-page-loading-breathe 1.3s ease-in-out infinite;
}
.ism-page-loading-copy {
  min-width: 0;
}
.ism-page-loading-title {
  color: #e0faff;
  font-size: 18px;
  font-weight: 600;
  line-height: 26px;
  letter-spacing: 0;
}
.ism-page-loading-subtitle {
  margin-top: 4px;
  color: rgba(224, 250, 255, 0.72);
  font-size: 13px;
  line-height: 19px;
  letter-spacing: 0;
}
.ism-page-loading-meter {
  position: relative;
  display: flex;
  gap: 6px;
  margin-top: 24px;
}
.ism-page-loading-meter span {
  height: 3px;
  flex: 1;
  background: #22d3ee;
  box-shadow: 0 0 12px rgba(34, 211, 238, 0.45);
  animation: ism-page-loading-pulse 1.2s ease-in-out infinite;
}
.ism-page-loading-meter span:nth-child(2) {
  background: #38bdf8;
  animation-delay: 0.15s;
}
.ism-page-loading-meter span:nth-child(3) {
  background: #60a5fa;
  animation-delay: 0.3s;
}
.ism-page-loading-meter span:nth-child(4) {
  background: #93c5fd;
  animation-delay: 0.45s;
}
@keyframes ism-page-loading-spin {
  to {
    transform: rotate(360deg);
  }
}
@keyframes ism-page-loading-scan {
  0% {
    transform: translateX(-120%);
    opacity: 0;
  }
  20%,
  80% {
    opacity: 1;
  }
  100% {
    transform: translateX(120%);
    opacity: 0;
  }
}
@keyframes ism-page-loading-pulse {
  0%,
  100% {
    opacity: 0.45;
    transform: scaleX(0.72);
  }
  50% {
    opacity: 1;
    transform: scaleX(1);
  }
}
@keyframes ism-page-loading-breathe {
  0%,
  100% {
    transform: scale(0.82);
    opacity: 0.72;
  }
  50% {
    transform: scale(1);
    opacity: 1;
  }
}
::v-deep .run-graph-container {
  position: absolute;
  height: 100%;
  width: 100%;
  flex:1;
  overflow: auto;
}
::v-deep  .ant-modal-header {
  padding: 10px 16px;
  background-color: var(--popUpBackColor);
  border-bottom: 0px solid #f0f0f0;
}
::v-deep  .ant-modal-close-x {
  width: 45px;
  height: 45px;
  line-height: 45px;
}
// 键盘样式
.simple-keyboard {
  position: absolute;
  bottom: 0;
  left: 5%;
  width: 90%;
  color: #000;
  z-index: 999999999;
}
.dialog-button{
  padding: 5px;
  text-align: right;
}
::v-deep  .ant-modal-body {
  padding: 0px;
}
::v-deep .ant-modal-content {
  position: relative;
  background-color: #fff;
  background-clip: padding-box;
  border: 0;
  border-radius: 2px;
  pointer-events: auto;
}

::v-deep .ant-divider-horizontal {
  margin: 5px 0;
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
  width : 4px;  /*高宽分别对应横竖滚动条的尺寸*/
  height: 4px;
}
body:-ms-fullscreen {
  overflow: auto!important;
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
