<template>
   <div :style="nodeMouseStyleVar"
        @click.stop="onNodeClick($event,'click')"
        @dblclick.stop="onNodeClick($event,'dblclick')"
   >
      <div  class="ism-render"   :style="layerStyle" v-if="LayerContainerData.layer&&SelectPagerID!=''" ref="ismrender">
        <div :class="{'animated':true,[`${LayerContainerData.layer.animate}`]: true}">
          <div class="pager-graph-container">
            <div :ref="detail.identifier" ></div>
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
        </div>
      </div>
     <div v-else>
      <a-empty
          image="/static/images/original.png"
          :image-style="{
          height: '60px',
        }"
      >
        <span slot="description">{{$t("configComponent.pagerContainer.tips")}}</span>
      </a-empty>
     </div>
        <a-modal :visible="PopUpDialog"
                      v-drag-modal
                      :style="stylePopUpVar"
                      :centered="true"
                      :forceRender="true"
                      :destroyOnClose="true"
                      :width="(!isExternUrl?PopUpContainerConfigData.layer.width:linkInfoExternal.width)+'px'"
                      :title="typeof PopUpContainerConfigData.PageName!='undefined'&&PopUpContainerConfigData.PageName!=''?PopUpContainerConfigData.PageName:'对话框'"
                      :footer="null"
                      @cancel="ClosePopDialog"
                      :modal="false">
           <template #title></template>
          <div  @click="PopUpDialogClick">
            <div v-if="!isExternUrl">
              <div  class="ism-popup-render"  @contextmenu.prevent="componentRightDialogClick" :style="popUpStyle" v-if="PopUpContainerConfigData.layer" ref="ismrender1">
                <div class="pager-graph-container" :class="{'animated':true,[`${PopUpContainerConfigData.layer.animate}`]: true}">
                    <div :ref="detail.identifier+'ISMPopUpRunningContainer'" ></div>
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
       <div v-show="showKeyboard&&LayerContainerData.layer.virtuallyKey">
        <SimpleKeyboard ref="SimpleKeyboard" @onChange="onChangeKeyboard" @onKeyPress="onKeyPress"/>
      </div>
  </div>
</template>
<script>
import {register} from "@antv/x6-vue-shape";

let res_components = {}
import {formatDate} from "@/utils/common";
import SimpleKeyboard from '@/components/SimpleKeyboard/SimpleKeyboard.vue'
const loadingKey = 'updatable'
import { mapActions, mapState, mapMutations ,mapGetters} from 'vuex'
import {setData} from "@/services/device";
import {ComponentRestApi} from "@/services/RestApi";
import TemplateRender from 'vue-template-render'
import {ExecSysScript} from "@/services/system";
import {AUTH_TYPE, checkAuthorization, getAuthorization, setAuthorization} from "@/utils/request";
import Vue from 'vue'
import Contextmenu from "vue-contextmenujs"
Vue.use(Contextmenu);
import '@/utils/vmodalDrage'
import AKeepAlive from "@/components/cache/AKeepAlive";
import {getRealDataByUuid} from "@/services/device";
import ISMGroupNode from "@/pages/ISMDisPlay/ISMGroupNode.vue";
import {Graph} from "@antv/x6";
import {getDisplayModelLayerData,getDisplayModelPagerLayerData,getLayerDataStructByToken} from "@/services/displayModel";
import {Selection} from "@antv/x6";
res_components["AKeepAlive"] = AKeepAlive
res_components["SimpleKeyboard"] = SimpleKeyboard
export default {
    name: 'ism-view-pager-container',
    inject: ['getNode'],
    i18n: require('@/i18n/language'),
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
      },
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
            this.$message.loading({ content: 'Loading...', key: loadingKey, duration: 0 });
            this.CurrentRealUUIDList=[]
            this.CurrentModelUUIDList=[]
            this.selectDisplayPageDataStruct({page:page,callback:function (uuids,devices,isFound){
                _t.chargePage = false
                if (isFound === false) {
                  _t.$message.destroy(loadingKey)
                  _t.$message.error(_t.$t("readData.NotFindPage"))
                  return
                }
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
                        for(let k = 0,realDataLen =res.data.realData.length;k<realDataLen;k++)
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
                },300)
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

          if(this.LayerContainerData.layer.autoSize==1)
          {
            setTimeout(function (){
              _t.setScale()
            },1000)
          }
        },
        deep: true
      }
    },
    components: res_components,
    computed: {
    nodeMouseStyleVar(){
      return {
        '--NodeMouseStyle':(this.SelectCurrentNodeData!=null)&&(typeof this.SelectCurrentNodeData.action !='undefined')&&(this.SelectCurrentNodeData.action.length > 0)?'pointer':'',
      };
    },
    layerStyle:function () {
      let scale = 1
      let styles = [`transform:scale(${scale})`];
      if(this.LayerContainerData.layer.backColor) {
        styles.push(`background-color: ${this.LayerContainerData.layer.backColor}`);
      }
      if(this.LayerContainerData.layer.backgroundImage) {
        styles.push(`background-image: url("${this.LayerContainerData.layer.backgroundImage}")`);
        styles.push(`background-size:100% 100%`);
      }
      if(this.LayerContainerData.layer.width > 0) {
        styles.push(`width: ${this.detail.style.position.w}px`);
      }
      if(this.LayerContainerData.layer.height > 0) {
        styles.push(`height: ${this.detail.style.position.h}px`);
      }
      let style = styles.join(';');
      return style;
    },
    popUpStyle:function () {
      let scale = 1
      let styles = [`transform:scale(${scale})`];
      if(this.PopUpContainerConfigData.layer.backColor) {
        styles.push(`background-color: ${this.PopUpContainerConfigData.layer.backColor}`);
      }
      if(this.PopUpContainerConfigData.layer.backgroundImage) {
        styles.push(`background-image: url("${this.PopUpContainerConfigData.layer.backgroundImage}")`);
        styles.push(`background-size:100% 100%`);
      }
      if(this.PopUpContainerConfigData.layer.width > 0) {
        styles.push(`width: ${this.PopUpContainerConfigData.layer.width}px`);
      }
      if(this.PopUpContainerConfigData.layer.height > 0) {
        styles.push(`height: ${this.PopUpContainerConfigData.layer.height}px`);
      }
      let style = styles.join(';');
      return style;
    },
    stylePopUpVar() {
        return {
          '--popUpBackColor':this.PopUpContainerConfigData.layer.backColor,
        };
      },
  },
    data() {
        return {
          ParentVuex:null,
          CurrentPagerRealDataUuidList:[],
          CurrentPagerRealDeviceUuidList:[],
          CurrentPagerPopRealDataUuidList:[],
          CurrentPagerPopRealDeviceUuidList:[],
          SelectCurrentNodeData:null,
          PCPageList:[],
          PhonePageList:[],
          isMobile:false,
          ISMPagerRuningCavasContainer:null,
          ISMPopUpRunningContainer:null,
          destroyMainGraphRaf: null,
          pageRenderToken: 0,
          pageLoadingToken: 0,
          pageLoadingStartedAt: 0,
          pageLoadingEl: null,
          pendingClickPageLoadingToken: null,
          _isDestroyed: false,
          detail:null,
          IsToolBox:false,
          editMode:true,
          LayerContainerData: {
            name: '--',
            layer: {
              backColor: '',
              backgroundImage: '',
              widthHeightRatio: '',
              width: 300,
              height:600,
              autoSize:0
            },
            components: []
          }, //当前场景的组态数据
          PopUpContainerConfigData:{//弹窗界面
            name: '--',
            layer: {
              backColor: '',
              backgroundImage: '',
              widthHeightRatio: '',
              width: 300,
              height:600
            },
            components: []
          },
          settingDialog:false,
          IsConfirm:false,
          showKeyboard:false,
          changeIpt:'',//选择了哪个输入框
          setPasswordDialog:false,
          actionPasswordDialog:false,
          PopUpDialog:false,
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
          pendingTimers:[],
          eventHandlers:{},
          nodeHandlers:{},
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
          currentDisplayUUID:"",
          currentPageUUID:"",
          currentPopUpDisplayUUID:"",
          currentPopUpPageUUID:"",
          deviceUuid:"",
          SetPassword:"",
          AutoSetValue:"",
          setDataUuid:"",
          settingVisible:false,
          settingLoading:false,
          SetForm:null,
          ActionFormPassword:null,
          SetFormPassword:null,
          fullScreen:false,
          CurrentRealUUIDList:[],
          CurrentModelUUIDList:[],
          CurrentPopRealUUIDList:[],
          CurrentPopModelUUIDList:[],
          Url:"",
          width:600,
          height:600,
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
          SelectPagerID:"",
          LaunchType:0,
          base:{
            "text": "configComponent.pagerContainer.Text",
            "icon": "icon-rongqi",
            "isFontIcon": true,
            "info": {
              "type": "ism-view-pager-container",
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
                  "w": 500,
                  "h": 200
                },
                "visible":1,
                "zIndex": -1,
                "transform": 0,
                "diy":[
                  {
                    "name":"configComponent.pagerContainer.LaunchType",
                    "type":6,
                    "enumList":[
                      {option:'True',value:1},{option:'False',value:0}
                    ],
                    "value":0,
                    "key":"LaunchType",
                  },
                  {
                    "name":"displayConfig.Properties.ComponentPager",
                    "type":13,
                    "value":"",
                    "Appid":"",
                    "key":"SelectPagerID",
                  }
                ]
              }
            }
          }
        }
    },
    methods: {
      InitPagerRealData() {
        let _t = this
        const datauuid = this.CurrentPagerRealDataUuidList
        const devices = this.CurrentPagerRealDeviceUuidList
        _t.CurrentRealUUIDList=[]
        _t.CurrentModelUUIDList=[]
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
        },50)
      },
      InitPopUpPagerRealData() {
        let _t = this
        const datauuid = this.CurrentPagerPopRealDataUuidList
        const devices = this.CurrentPagerPopRealDeviceUuidList
        _t.CurrentPopRealUUIDList=[]
        _t.CurrentPopModelUUIDList=[]
        setTimeout(function (){
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
        },50)
      },
      onNodeClick(e,type){
        if(!this.editMode) {
          this.doComponentAction(this.SelectCurrentNodeData, type)
        }
      },
      getLayerDataStruct  (data) {
        let params={
          muid:data.uuid
        }
        let bangDingData=[]
        let bangDingDeviceSN=[]
        let isPopUp = data.isPopUp?data.isPopUp:false
        getDisplayModelLayerData(params).then(function (res){
          if(res.data.code==0)
          {
            let pageLayer = res.data.layer
            let is_find_home = 0
            if(pageLayer.length>0)
            {
              let pcPageData = []
              let phonePageData = []
              for(let i=0;i<pageLayer.length;i++)
              {
                let pageInfo = {
                  id: 9,
                  key: 0,
                  isEdit: false,
                  pageUuid: "",
                  pageModelUuid: "",
                  isNewItem: false,
                  title: '',
                  depth: 1,
                  pageType:1,
                  scopedSlots: { title: 'custom' },
                }
                pageInfo.id = pageLayer[i].ID
                pageInfo.key = i
                pageInfo.isComponents=false
                pageInfo.title = pageLayer[i].PageName
                pageLayer[i].AppName = res.data.Display.name
                pageInfo.IsHome = pageLayer[i].IsHome
                pageInfo.IsLogin = pageLayer[i].IsLogin
                pageInfo.AppName = pageLayer[i].AppName
                pageInfo.pageUuid = pageLayer[i].PageId
                pageInfo.pageType = pageLayer[i].PageType
                pageInfo.pageModelUuid = pageLayer[i].modelId
                pageInfo.children=[]

                pageLayer[i].name = pageLayer[i].PageName
                let tempConfigData = pageLayer[i]
                try{
                  tempConfigData.layer = JSON.parse(tempConfigData.layer)
                  if (tempConfigData.components=="")
                  {
                    tempConfigData.components=[]
                  }
                  else{
                    tempConfigData.components = JSON.parse(tempConfigData.components)
                  }

                }catch (e) {
                  tempConfigData.components.cells=[]
                  continue
                }
                for(let k=0;k<tempConfigData.components.cells.length;k++)
                {
                  let components = {
                    isComponents:true,
                    title:tempConfigData.components.cells[k].data.detail.name,
                    key:tempConfigData.components.cells[k].id,
                    cellid:tempConfigData.components.cells[k].data.detail.identifier
                  }
                  pageInfo.children.push(components)
                }

                if(pageLayer[i].IsHome==1&!isPopUp)
                {
                  if(data.pageType)
                  {
                    if(pageLayer[i].PageType==0)
                    {
                      is_find_home=1
                      for(let k=0;k<tempConfigData.components.cells.length;k++)
                      {
                        if(typeof tempConfigData.components.cells[k].data.detail.active!="undefined") {
                          for (let kv = 0; kv < tempConfigData.components.cells[k].data.detail.active.length; kv++) {
                            if(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID!="") {
                              bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                            }
                            if(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN!="") {
                              bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                            }
                          }
                        }

                        if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&tempConfigData.components.cells[k].data.detail.animate.condition.dataID!="")
                        {
                          if(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN!="") {
                            bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                          }
                          bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
                        }

                        if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&typeof tempConfigData.components.cells[k].data.detail.animate.move!="undefined")
                        {
                          if(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN!="") {
                            bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                          }
                          bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                          if(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN!="") {
                            bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                          }
                          bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
                        }
                        else{
                          if(typeof tempConfigData.components.cells[k].data.detail.animate!="undefined") {
                            tempConfigData.components.cells[k].data.detail.animate.move = {
                              x: {
                                deviceSN: "",
                                selectVideoType: 0,
                                isBandDevice: false,
                                bandType: 1,
                                dataID: "",
                                dataName: "",
                              },
                              y: {
                                deviceSN: "",
                                selectVideoType: 0,
                                isBandDevice: false,
                                bandType: 1,
                                dataID: "",
                                dataName: "",
                              },
                            }
                          }
                        }
                      }
                    }
                  }
                  else
                  {
                    if(pageLayer[i].PageType==1) {
                      is_find_home=1
                      for(let k=0;k<tempConfigData.components.cells.length;k++)
                      {
                        if(typeof tempConfigData.components.cells[k].data.detail.active!="undefined") {
                          for (let kv = 0; kv < tempConfigData.components.cells[k].data.detail.active.length; kv++) {
                            if(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID!="") {
                              bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                            }
                            if(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN!="") {
                              bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                            }
                          }
                        }

                        if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&tempConfigData.components.cells[k].data.detail.animate.condition.dataID!="")
                        {
                          if(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN!="") {
                            bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                          }
                          bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
                        }
                        if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&typeof tempConfigData.components.cells[k].data.detail.animate.move!="undefined")
                        {
                          if(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN!=="") {
                            bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                          }
                          bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                          if(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN!="") {
                            bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                          }
                          bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
                        }
                        else{
                          if(typeof tempConfigData.components.cells[k].data.detail.animate!="undefined") {
                            tempConfigData.components.cells[k].data.detail.animate.move = {
                              x: {
                                deviceSN: "",
                                selectVideoType: 0,
                                isBandDevice: false,
                                bandType: 1,
                                dataID: "",
                                dataName: "",
                              },
                              y: {
                                deviceSN: "",
                                selectVideoType: 0,
                                isBandDevice: false,
                                bandType: 1,
                                dataID: "",
                                dataName: "",
                              },
                            }
                          }
                        }
                      }
                    }
                  }


                }
                pageInfo.pageLayerData = tempConfigData
                if(pageLayer[i].PageType==1)
                {
                  pcPageData.push(pageInfo)
                }else{
                  phonePageData.push(pageInfo)
                }
              }
              if(is_find_home==0&!isPopUp)
              {
                let tempConfigData = pageLayer[0]
                for(let k=0;k<tempConfigData.components.cells.length;k++)
                {
                  if(typeof tempConfigData.components.cells[k].data.detail.active!="undefined") {
                    for (let kv = 0; kv < tempConfigData.components.cells[k].data.detail.active.length; kv++) {
                      if(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID!="") {
                        bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                      }
                      if(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN!="") {
                        bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                      }
                    }
                  }

                  if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&tempConfigData.components.cells[k].data.detail.animate.condition.dataID!="")
                  {
                    if(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN!="") {
                      bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
                  }
                  if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&typeof tempConfigData.components.cells[k].data.detail.animate.move!="undefined")
                  {
                    if(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN!="") {
                      bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                    if(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN!="") {
                      bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
                  }
                  else{
                    if(typeof tempConfigData.components.cells[k].data.detail.animate!="undefined") {
                      tempConfigData.components.cells[k].data.detail.animate.move = {
                        x: {
                          deviceSN: "",
                          selectVideoType: 0,
                          isBandDevice: false,
                          bandType: 1,
                          dataID: "",
                          dataName: "",
                        },
                        y: {
                          deviceSN: "",
                          selectVideoType: 0,
                          isBandDevice: false,
                          bandType: 1,
                          dataID: "",
                          dataName: "",
                        },
                      }
                    }
                  }
                }
              }
            }
            let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
            data.cb(0,res.data.Display.project_uuid,bangDingData,newbangDingDeviceSN)
          }
          else
          {
            data.cb(-1,"",bangDingData,bangDingDeviceSN)
          }

        })
      },
      selectPopUpPagerContainerDisplayPageDataStruct  (page) {
        let _t = this
        let PCPageInfo = this.PCPageList
        let PhonePageInfo = this.PhonePageList

        let pageid = page.page.pageUuid
        let bangDingData=[]
        let bangDingDeviceSN=[]

        for(let i=0,PCPageInfoLen=PCPageInfo.length;i<PCPageInfoLen;i++)
        {
          if(PCPageInfo[i].pageUuid==pageid)
          {
            let tempConfigData = PCPageInfo[i].pageLayerData
            tempConfigData.name = page.page.title
            for(let k=0,componentsLen=tempConfigData.components.cells.length;k<componentsLen;k++)
            {
              if(typeof tempConfigData.components.cells[k].data.detail.active!="undefined") {
                for (let kv = 0,componentsAlen=tempConfigData.components.cells[k].data.detail.active.length; kv < componentsAlen; kv++) {
                  if(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID!="") {
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                  }
                  if(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN!="") {
                    bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                  }
                }
              }

              if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&tempConfigData.components.cells[k].data.detail.animate.condition.dataID!="")
              {
                if(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN!="") {
                  bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                }
                bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
              }
              if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&typeof tempConfigData.components.cells[k].data.detail.animate.move!="undefined")
              {
                if(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN!="") {
                  bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                }
                bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                if(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN!="") {
                  bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                }
                bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
              }
              else{
                if(typeof tempConfigData.components.cells[k].data.detail.animate!="undefined") {
                  tempConfigData.components.cells[k].data.detail.animate.move = {
                    x: {
                      deviceSN: "",
                      selectVideoType: 0,
                      isBandDevice: false,
                      bandType: 1,
                      dataID: "",
                      dataName: "",
                    },
                    y: {
                      deviceSN: "",
                      selectVideoType: 0,
                      isBandDevice: false,
                      bandType: 1,
                      dataID: "",
                      dataName: "",
                    },
                  }
                }
              }
            }
            let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
            page.callback(0,tempConfigData,bangDingData,newbangDingDeviceSN)
            return
          }
        }

        for(let i=0,PhonePageInfoLen=PhonePageInfo.length;i<PhonePageInfoLen;i++)
        {
          if(PhonePageInfo[i].pageUuid==pageid) {
            let tempConfigData = PhonePageInfo[i].pageLayerData
            tempConfigData.name = page.page.title
            for(let k=0,componentsLen=tempConfigData.components.cells.length;k<componentsLen;k++)
            {
              if(typeof tempConfigData.components.cells[k].data.detail.active!="undefined") {
                for (let kv = 0,componentsAlen= tempConfigData.components.cells[k].data.detail.active.length; kv <componentsAlen; kv++) {
                  if(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID!="") {
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                  }
                  if(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN!="") {
                    bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                  }
                }
              }

              if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&tempConfigData.components.cells[k].data.detail.animate.condition.dataID!="")
              {
                if(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN!="") {
                  bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                }
                bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
              }
              if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&typeof tempConfigData.components.cells[k].data.detail.animate.move!="undefined")
              {
                if(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN!="") {
                  bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                }
                bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                if(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN!="") {
                  bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                }
                bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
              }
              else{
                if(typeof tempConfigData.components.cells[k].data.detail.animate!="undefined") {
                  tempConfigData.components.cells[k].data.detail.animate.move = {
                    x: {
                      deviceSN: "",
                      selectVideoType: 0,
                      isBandDevice: false,
                      bandType: 1,
                      dataID: "",
                      dataName: "",
                    },
                    y: {
                      deviceSN: "",
                      selectVideoType: 0,
                      isBandDevice: false,
                      bandType: 1,
                      dataID: "",
                      dataName: "",
                    },
                  }
                }
              }
            }
            let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
            page.callback(0,tempConfigData,bangDingData,newbangDingDeviceSN)
            return
          }
        }

        this.getLayerDataStruct({uuid:page.page.displayUUID,isPopUp:true,cb:function () {
            let PCPageInfo = _t.PCPageList
            let PhonePageInfo = _t.PhonePageList

            let pageid = page.page.pageUuid
            let bangDingData = []
            let bangDingDeviceSN = []
            for (let i = 0, PCPageInfoLen = PCPageInfo.length; i < PCPageInfoLen; i++) {
              if (PCPageInfo[i].pageUuid == pageid) {
                let tempConfigData = PCPageInfo[i].pageLayerData
                tempConfigData.name = page.page.title
                for (let k = 0, componentsLen = tempConfigData.components.cells.length; k < componentsLen; k++) {
                  if (typeof tempConfigData.components.cells[k].data.detail.active != "undefined") {
                    for (let kv = 0, componentsAlen = tempConfigData.components.cells[k].data.detail.active.length; kv < componentsAlen; kv++) {
                      if (tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID != "") {
                        bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                      }
                      if (tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN != "") {
                        bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                      }
                    }
                  }

                  if ((typeof tempConfigData.components.cells[k].data.detail.animate != "undefined") && tempConfigData.components.cells[k].data.detail.animate.condition.dataID != "") {
                    if (tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN != "") {
                      bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
                  }
                  if ((typeof tempConfigData.components.cells[k].data.detail.animate != "undefined") && typeof tempConfigData.components.cells[k].data.detail.animate.move != "undefined") {
                    if (tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN != "") {
                      bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                    if (tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN != "") {
                      bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
                  } else {
                    if (typeof tempConfigData.components.cells[k].data.detail.animate != "undefined") {
                      tempConfigData.components.cells[k].data.detail.animate.move = {
                        x: {
                          deviceSN: "",
                          selectVideoType: 0,
                          isBandDevice: false,
                          bandType: 1,
                          dataID: "",
                          dataName: "",
                        },
                        y: {
                          deviceSN: "",
                          selectVideoType: 0,
                          isBandDevice: false,
                          bandType: 1,
                          dataID: "",
                          dataName: "",
                        },
                      }
                    }
                  }
                }
                let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
                page.callback(0,tempConfigData, bangDingData, newbangDingDeviceSN)
                return
              }
            }

            for (let i = 0, PhonePageInfoLen = PhonePageInfo.length; i < PhonePageInfoLen; i++) {
              if (PhonePageInfo[i].pageUuid == pageid) {
                let tempConfigData = PhonePageInfo[i].pageLayerData
                tempConfigData.name = page.page.title
                for (let k = 0, componentsLen = tempConfigData.components.cells.length; k < componentsLen; k++) {
                  if (typeof tempConfigData.components.cells[k].data.detail.active != "undefined") {
                    for (let kv = 0, componentsAlen = tempConfigData.components.cells[k].data.detail.active.length; kv < componentsAlen; kv++) {
                      if (tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID != "") {
                        bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                      }
                    }
                  }

                  if ((typeof tempConfigData.components.cells[k].data.detail.animate != "undefined") && tempConfigData.components.cells[k].data.detail.animate.condition.dataID != "") {
                    if (tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN != "") {
                      bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
                  }
                  if ((typeof tempConfigData.components.cells[k].data.detail.animate != "undefined") && typeof tempConfigData.components.cells[k].data.detail.animate.move != "undefined") {
                    if (tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN != "") {
                      bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                    if (tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN != "") {
                      bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
                  } else {
                    if (typeof tempConfigData.components.cells[k].data.detail.animate != "undefined") {
                      tempConfigData.components.cells[k].data.detail.animate.move = {
                        x: {
                          deviceSN: "",
                          selectVideoType: 0,
                          isBandDevice: false,
                          bandType: 1,
                          dataID: "",
                          dataName: "",
                        },
                        y: {
                          deviceSN: "",
                          selectVideoType: 0,
                          isBandDevice: false,
                          bandType: 1,
                          dataID: "",
                          dataName: "",
                        },
                      }
                    }
                  }
                }
                let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
                page.callback(0, tempConfigData,bangDingData, newbangDingDeviceSN)
                return
              }
            }
            const tempConfigData ={ "name": "--", "layer": { "backColor": "", "backgroundImage": "", "widthHeightRatio": "", "width": 300, "height": 600 }, "components": [] }
            page.callback(-1,tempConfigData,bangDingData,bangDingDeviceSN)
          }});
      },
      selectDisplayPageContainerDataStruct (page)  {
        let _t = this
        let PCPageInfo = this.PCPageList
        let PhonePageInfo = this.PhonePageList

        let pageid = page.page.pageUuid
        let bangDingData=[]
        let bangDingDeviceSN=[]
        for(let i=0,PCPageInfoLen=PCPageInfo.length;i<PCPageInfoLen;i++)
        {
          if(PCPageInfo[i].pageUuid==pageid)
          {
            let tempConfigData = PCPageInfo[i].pageLayerData
            tempConfigData.name = page.page.title
            for(let k=0,componentsLen=tempConfigData.components.cells.length;k<componentsLen;k++)
            {
              if(typeof tempConfigData.components.cells[k].data.detail.active!="undefined") {
                for (let kv = 0,componentsALen=tempConfigData.components.cells[k].data.detail.active.length; kv < componentsALen; kv++) {
                  if(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID!="") {
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                  }
                  if(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN!="") {
                    bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                  }
                }
              }

              if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&tempConfigData.components.cells[k].data.detail.animate.condition.dataID!="")
              {
                if(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN!="") {
                  bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                }
                bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
              }
              if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&typeof tempConfigData.components.cells[k].data.detail.animate.move!="undefined")
              {
                if(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN!="") {
                  bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                }
                bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                if(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN!="") {
                  bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                }
                bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
              }
              else{
                if(typeof tempConfigData.components.cells[k].data.detail.animate!="undefined") {
                  tempConfigData.components.cells[k].data.detail.animate.move = {
                    x: {
                      deviceSN: "",
                      selectVideoType: 0,
                      isBandDevice: false,
                      bandType: 1,
                      dataID: "",
                      dataName: "",
                    },
                    y: {
                      deviceSN: "",
                      selectVideoType: 0,
                      isBandDevice: false,
                      bandType: 1,
                      dataID: "",
                      dataName: "",
                    },
                  }
                }
              }
            }
            let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
            page.callback(0,tempConfigData,bangDingData,newbangDingDeviceSN)
            return
          }
        }

        for(let i=0,PhonePageInfoLen=PhonePageInfo.length;i<PhonePageInfoLen;i++)
        {
          if(PhonePageInfo[i].pageUuid==pageid) {
            let tempConfigData = PhonePageInfo[i].pageLayerData
            tempConfigData.name = page.page.title
            this.currentPageUUID = tempConfigData.PageId
            for(let k=0,componentsLen=tempConfigData.components.cells.length;k<componentsLen;k++)
            {
              if(typeof tempConfigData.components.cells[k].data.detail.active!="undefined") {
                for (let kv = 0,componentsALen=tempConfigData.components.cells[k].data.detail.active.length; kv < componentsALen; kv++) {
                  if(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID!="") {
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                  }
                  if(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN!="") {
                    bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                  }
                }
              }

              if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&tempConfigData.components.cells[k].data.detail.animate.condition.dataID!="")
              {
                if(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN!="") {
                  bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                }
                bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
              }
              if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&typeof tempConfigData.components.cells[k].data.detail.animate.move!="undefined")
              {
                if(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN!="") {
                  bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                }
                bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                if(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN!="") {
                  bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                }
                bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
              }
              else{
                if(typeof tempConfigData.components.cells[k].data.detail.animate!="undefined") {
                  tempConfigData.components.cells[k].data.detail.animate.move = {
                    x: {
                      deviceSN: "",
                      selectVideoType: 0,
                      isBandDevice: false,
                      bandType: 1,
                      dataID: "",
                      dataName: "",
                    },
                    y: {
                      deviceSN: "",
                      selectVideoType: 0,
                      isBandDevice: false,
                      bandType: 1,
                      dataID: "",
                      dataName: "",
                    },
                  }
                }
              }
            }
            let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
            page.callback(0,tempConfigData,bangDingData,newbangDingDeviceSN)
            return
          }
        }

        this.getLayerDataStruct({uuid:page.page.displayUUID,cb:function (){
            let PCPageInfo = _t.PCPageList
            let PhonePageInfo = _t.PhonePageList

            let pageid = page.page.pageUuid
            let bangDingData=[]
            let bangDingDeviceSN=[]
            for(let i=0,PCPageInfoLen=PCPageInfo.length;i<PCPageInfoLen;i++)
            {
              if(PCPageInfo[i].pageUuid==pageid)
              {
                let tempConfigData = PCPageInfo[i].pageLayerData
                tempConfigData.name = page.page.title
                for(let k=0,componentsLen=tempConfigData.components.cells.length;k<componentsLen;k++)
                {
                  if(typeof tempConfigData.components.cells[k].data.detail.active!="undefined") {
                    for (let kv = 0,componentsALen=tempConfigData.components.cells[k].data.detail.active.length; kv < componentsALen; kv++) {
                      if(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID!="") {
                        bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                      }
                      if(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN!="") {
                        bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                      }
                    }
                  }

                  if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&tempConfigData.components.cells[k].data.detail.animate.condition.dataID!="")
                  {
                    if(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN!="") {
                      bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
                  }
                  if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&typeof tempConfigData.components.cells[k].data.detail.animate.move!="undefined")
                  {
                    if(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN!="") {
                      bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                    if(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN!="") {
                      bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
                  }
                  else{
                    if(typeof tempConfigData.components.cells[k].data.detail.animate!="undefined") {
                      tempConfigData.components.cells[k].data.detail.animate.move = {
                        x: {
                          deviceSN: "",
                          selectVideoType: 0,
                          isBandDevice: false,
                          bandType: 1,
                          dataID: "",
                          dataName: "",
                        },
                        y: {
                          deviceSN: "",
                          selectVideoType: 0,
                          isBandDevice: false,
                          bandType: 1,
                          dataID: "",
                          dataName: "",
                        },
                      }
                    }
                  }
                }
                let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
                page.callback(0,tempConfigData,bangDingData,newbangDingDeviceSN)
                return
              }
            }

            for(let i=0,PhonePageInfoLen=PhonePageInfo.length;i<PhonePageInfoLen;i++)
            {
              if(PhonePageInfo[i].pageUuid==pageid) {
                let tempConfigData = PhonePageInfo[i].pageLayerData
                tempConfigData.name = page.page.title
                for(let k=0,componentsLen=tempConfigData.components.cells.length;k<componentsLen;k++)
                {
                  if(typeof tempConfigData.components.cells[k].data.detail.active!="undefined") {
                    for (let kv = 0,componentsALen=tempConfigData.components.cells[k].data.detail.active.length; kv < componentsALen; kv++) {
                      if(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID!="") {
                        bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                      }
                      if(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN!="") {
                        bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                      }
                    }
                  }

                  if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&tempConfigData.components.cells[k].data.detail.animate.condition.dataID!="")
                  {
                    if(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN!="") {
                      bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
                  }
                  if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&typeof tempConfigData.components.cells[k].data.detail.animate.move!="undefined")
                  {
                    if(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN!="") {
                      bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                    if(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN!="") {
                      bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                    }
                    bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
                  }
                  else{
                    if(typeof tempConfigData.components.cells[k].data.detail.animate!="undefined") {
                      tempConfigData.components.cells[k].data.detail.animate.move = {
                        x: {
                          deviceSN: "",
                          selectVideoType: 0,
                          isBandDevice: false,
                          bandType: 1,
                          dataID: "",
                          dataName: "",
                        },
                        y: {
                          deviceSN: "",
                          selectVideoType: 0,
                          isBandDevice: false,
                          bandType: 1,
                          dataID: "",
                          dataName: "",
                        },
                      }
                    }
                  }
                }
                let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
                page.callback(0,tempConfigData,bangDingData,newbangDingDeviceSN)
                return
              }
            }
            page.callback(-1,null,bangDingData,bangDingDeviceSN)
          }});
      },
      getLayerPagerContainerDataStruct (data) {
        let params={
          pageid:data.pageid
        }
        let bangDingData=[]
        let bangDingDeviceSN=[]
        getDisplayModelPagerLayerData(params).then(function (res){
          if(res.data.code==0)
          {
            let pageLayer = res.data.layer
            let tempConfigData = pageLayer
            try{
              tempConfigData.layer = JSON.parse(tempConfigData.layer)
              if (tempConfigData.components=="")
              {
                tempConfigData.components=[]
              }
              else{
                tempConfigData.components = JSON.parse(tempConfigData.components)
              }
              for(let k=0;k<tempConfigData.components.cells.length;k++)
              {
                if(typeof tempConfigData.components.cells[k].data.detail.active!="undefined") {
                  for (let kv = 0; kv < tempConfigData.components.cells[k].data.detail.active.length; kv++) {
                    if(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID!="") {
                      bangDingData.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.dataID)
                    }
                    if(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN!="") {
                      bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.active[kv].condition.deviceSN)
                    }
                  }
                }

                if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&tempConfigData.components.cells[k].data.detail.animate.condition.dataID!="")
                {
                  if(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN!="") {
                    bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.condition.deviceSN)
                  }
                  bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.condition.dataID)
                }
                if((typeof tempConfigData.components.cells[k].data.detail.animate!="undefined")&&typeof tempConfigData.components.cells[k].data.detail.animate.move!="undefined")
                {
                  if(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN!="") {
                    bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.x.deviceSN)
                  }
                  bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.x.dataID)

                  if(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN!="") {
                    bangDingDeviceSN.push(tempConfigData.components.cells[k].data.detail.animate.move.y.deviceSN)
                  }
                  bangDingData.push(tempConfigData.components.cells[k].data.detail.animate.move.y.dataID)
                }
                else{
                  if(typeof tempConfigData.components.cells[k].data.detail.animate!="undefined") {
                    tempConfigData.components.cells[k].data.detail.animate.move = {
                      x: {
                        deviceSN: "",
                        selectVideoType: 0,
                        isBandDevice: false,
                        bandType: 1,
                        dataID: "",
                        dataName: "",
                      },
                      y: {
                        deviceSN: "",
                        selectVideoType: 0,
                        isBandDevice: false,
                        bandType: 1,
                        dataID: "",
                        dataName: "",
                      },
                    }
                  }
                }
              }
              let newbangDingDeviceSN = Array.from(new Set(bangDingDeviceSN));
              data.cb(0,tempConfigData,bangDingData,newbangDingDeviceSN)

            }catch (e) {
              console.log(e)
              data.cb(-3,null,bangDingData,bangDingDeviceSN)
            }
          }
          else
          {
            data.cb(-1,null,bangDingData,bangDingDeviceSN)
          }

        })
      },
      initComponents(option){
        this.width = option.style.position.w
        this.height = option.style.position.h
        let i=0
        for(i=0;i<option.style.diy.length;i++)
        {
          if(option.style.diy[i].key=="SelectPagerID")
          {
            let temp=option.style.diy[i].value
            if(temp!=this.SelectPagerID) {
              this.SelectPagerID=temp
              this.getPagerData(this.SelectPagerID)
            }
          }
          else if(option.style.diy[i].key=="LaunchType")
          {
            this.LaunchType = option.style.diy[i].value
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
                  this.ISMPagerRuningCavasContainer.togglePanning(_t.tempAutoPadding)
                  this.ISMPagerRuningCavasContainer.toggleMouseWheel(_t.tempAutoPadding)
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
      },
      destroyMainGraph() {
        window.removeEventListener('resize', this.setScale)
        if (this.ISMPagerRuningCavasContainer) {
          try {
            this.ISMPagerRuningCavasContainer.clearCells()
            this.ISMPagerRuningCavasContainer.off()
            this.ISMPagerRuningCavasContainer.dispose()
          } catch (e) {
            console.error('[ViewPagerContainer.destroyMainGraph] dispose failed, will force cleanup', e)
          } finally {
            if (this.ISMPagerRuningCavasContainer) {
              try { this.ISMPagerRuningCavasContainer.off() } catch(e) { console.warn(e) }
              try { this.ISMPagerRuningCavasContainer.dispose() } catch(e) { console.warn(e) }
              try {
                const refObj = this.detail ? this.detail.identifier : null
                if (refObj) {
                  const container = this.$refs[refObj]
                  if (container) {
                    this.destroyMainGraphRaf = requestAnimationFrame(() => {
                      container.innerHTML = ''
                      this.destroyMainGraphRaf = null
                    })
                  }
                }
              } catch(e) { console.warn(e) }
            }
            this.ISMPagerRuningCavasContainer = null
          }
        }
      },
      destroyPopUpGraph() {
        if (this.ISMPopUpRunningContainer) {
          try {
            this.ISMPopUpRunningContainer.clearCells()
            this.ISMPopUpRunningContainer.off()
            this.ISMPopUpRunningContainer.dispose()
          } catch (e) {
            console.error('[ViewPagerContainer.destroyPopUpGraph] dispose failed, will force cleanup', e)
          } finally {
            if (this.ISMPopUpRunningContainer) {
              try { this.ISMPopUpRunningContainer.off() } catch(e) { console.warn(e) }
              try { this.ISMPopUpRunningContainer.dispose() } catch(e) { console.warn(e) }
              try {
                const refObj = this.detail ? (this.detail.identifier + 'ISMPopUpRunningContainer') : null
                if (refObj) {
                  const container = this.$refs[refObj]
                  if (container) container.innerHTML = ''
                }
              } catch(e) { console.warn(e) }
            }
            this.ISMPopUpRunningContainer = null
          }
        }
        this.clearPendingTimers('popup')
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
          this.ISMPagerRuningCavasContainer.resize(this.detail.style.position.w, this.detail.style.position.h);
          const dd = this.detail.style.position.w / this.LayerContainerData.layer.width;
          const dh = this.detail.style.position.h / this.LayerContainerData.layer.height;
          this.ISMPagerRuningCavasContainer.scale( dd,  dh);
        }
        else {
          this.ISMPagerRuningCavasContainer.resize(this.detail.style.position.w, this.detail.style.position.h);
          this.ISMPagerRuningCavasContainer.scale( 1,  1);
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
        const w = window.innerWidth / this.LayerContainerData.layer.width;
        const h = window.innerHeight / this.LayerContainerData.layer.height;
        return [w, h];
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
          return action.action === 'DeviceView' && action.DeviceView && action.DeviceView.isContainer
        })
      },
      beginPageLoading(key = loadingKey) {
        this.$message.destroy(key)
        this.pageLoadingToken += 1
        this.pageLoadingStartedAt = Date.now()
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
        if (token && token === this.pageLoadingToken) {
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
      doComponentAction(component,active){
        let _this = this
        if(component==null)
        {
          return
        }
        if (active == 'click' && this.hasPageJumpAction(component, active)) {
          this.pendingClickPageLoadingToken = this.beginPageLoading(loadingKey)
        } else if (active == 'click') {
          this.pendingClickPageLoadingToken = null
        }
        if(active=='click')
        {
          clearTimeout(this.clickTimer);  //首先清除计时器
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
        this.$message.destroy(loadingKey)
        this.PopUpDialog  = false
        this.currentPopUpDisplayUUID = ""
        this.currentPopUpPageUUID = ""
        this.destroyPopUpGraph()
      },
      PopUpDialogClick(){

        if(this.IsAutoClose)
        {
          this.PopUpDialog  = false
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

        for (let k = 0,Datalen = realData.Data.length; k < Datalen; k++) {
          if((_t.CurrentRealUUIDList.indexOf(realData.Data[k].Uuid)==-1)&&(_t.CurrentModelUUIDList.indexOf(realData.Data[k].ModelDataUuid)==-1))
          {
            continue
          }
          for(let j = 0,componentsLen = _t.LayerContainerData.components.cells.length;j<componentsLen;j++)
          {
            if((typeof _t.LayerContainerData.components.cells[j].data.detail.animate!="undefined")&&(typeof _t.LayerContainerData.components.cells[j].data.detail.animate.condition!="undefined"))
            {
              let  condition = _t.LayerContainerData.components.cells[j].data.detail.animate.condition
              let  selectAnimate = _t.LayerContainerData.components.cells[j].data.detail.animate.selected
              let isExpression = _t.LayerContainerData.components.cells[j].data.detail.animate.isExpression
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
                        const cell = _t.ISMPagerRuningCavasContainer.getCellById(_t.LayerContainerData.components.cells[j].id)
                        if (isStart) {
                          cell.setVisible(true)
                          _t.LayerContainerData.components.cells[j].data.detail.style.visible = true
                        } else {
                          _t.LayerContainerData.components.cells[j].data.detail.style.visible = false
                          cell.setVisible(false)
                        }
                      }
                      this.$EventBus.$emit(_t.LayerContainerData.components.cells[j].data.detail.identifier + "animateEvent", isStart);
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
                      const cell = _t.ISMPagerRuningCavasContainer.getCellById(_t.LayerContainerData.components.cells[j].id)
                      if (isStart) {
                        cell.setVisible(true)
                        _t.LayerContainerData.components.cells[j].data.detail.style.visible = true
                      } else {
                        _t.LayerContainerData.components.cells[j].data.detail.style.visible = false
                        cell.setVisible(false)
                      }
                    }
                    this.$EventBus.$emit(_t.LayerContainerData.components.cells[j].data.detail.identifier + "animateEvent", isStart);
                  }
                }
              }
            }
            //位移动画
            if((typeof _t.LayerContainerData.components.cells[j].data.detail.animate!="undefined")&&(typeof _t.LayerContainerData.components.cells[j].data.detail.animate.move!="undefined")&&_t.LayerContainerData.components.cells[j].data.detail.animate.selected.includes("animateMove"))
            {
              let  conditionx = _t.LayerContainerData.components.cells[j].data.detail.animate.move.x
              let  conditiony = _t.LayerContainerData.components.cells[j].data.detail.animate.move.y
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
                    _t.LayerContainerData.components.cells[j].data.detail.style.position.x = RealValue
                    this.$EventBus.$emit(_t.LayerContainerData.components.cells[j].data.detail.identifier + "animateMove", animateMove);
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
                  _t.LayerContainerData.components.cells[j].data.detail.style.position.x = RealValue
                  this.$EventBus.$emit(_t.LayerContainerData.components.cells[j].data.detail.identifier + "animateMove", animateMove);
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
                    _t.LayerContainerData.components.cells[j].data.detail.style.position.y = RealValue
                    this.$EventBus.$emit(_t.LayerContainerData.components.cells[j].data.detail.identifier + "animateMove", animateMove);
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
                  _t.LayerContainerData.components.cells[j].data.detail.style.position.y = RealValue
                  this.$EventBus.$emit(_t.LayerContainerData.components.cells[j].data.detail.identifier + "animateMove", animateMove);
                }
              }
            }

            if((typeof _t.LayerContainerData.components.cells[j].data.detail.active!="undefined"))
            {
              let  active = _t.LayerContainerData.components.cells[j].data.detail.active
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
                      const cell = _t.ISMPagerRuningCavasContainer.getCellById(_t.LayerContainerData.components.cells[j].id)
                      const isEdge = cell.isEdge()
                      if(!isEdge) {
                        _t.$EventBus.$emit(_t.LayerContainerData.components.cells[j].data.detail.identifier + "activeEvent", activeData);
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
                    const cell = _t.ISMPagerRuningCavasContainer.getCellById(_t.LayerContainerData.components.cells[j].id)
                    if(cell) {
                      const isEdge = cell.isEdge()
                      if (!isEdge) {
                        _t.$EventBus.$emit(_t.LayerContainerData.components.cells[j].data.detail.identifier + "activeEvent", activeData);
                      } else {
                        let animation = ''
                        if (activeData.ID == "Forward" && activeData.result) {
                          animation = 'ant-line-forward 30s infinite linear'
                        }
                        if (activeData.ID == "Reverse" && activeData.result) {
                          animation = 'ant-line-inverse 30s infinite linear'
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

        if((!_t.chargePagePopUp)&&(_t.PopUpDialog))
        {
          for (let k = 0,realDataLen = realData.Data.length; k < realDataLen; k++) {
            if((_t.CurrentPopRealUUIDList.indexOf(realData.Data[k].Uuid)==-1)&&(_t.CurrentPopModelUUIDList.indexOf(realData.Data[k].ModelDataUuid)==-1))
            {
              continue
            }
            for(let j = 0,componentsLen = _t.PopUpContainerConfigData.components.cells.length;j<componentsLen;j++)
            {
              if((typeof _t.PopUpContainerConfigData.components.cells[j].data.detail.animate!="undefined")&&(typeof _t.PopUpContainerConfigData.components.cells[j].data.detail.animate.condition!="undefined"))
              {
                let  condition = _t.PopUpContainerConfigData.components.cells[j].data.detail.animate.condition
                let  selectAnimate = _t.PopUpContainerConfigData.components.cells[j].data.detail.animate.selected
                let isExpression = _t.PopUpContainerConfigData.components.cells[j].data.detail.animate.isExpression
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
                          const cell = _t.ISMPopUpRunningContainer.getCellById(_t.PopUpContainerConfigData.components.cells[j].id)
                          if (isStart) {
                            cell.setVisible(true)
                            _t.PopUpContainerConfigData.components.cells[j].data.detail.style.visible = true
                          } else {
                            _t.PopUpContainerConfigData.components.cells[j].data.detail.style.visible = false
                            cell.setVisible(false)
                          }
                        }
                        _t.$EventBus.$emit(_t.PopUpContainerConfigData.components.cells[j].data.detail.identifier + "animateEvent", isStart);
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
                        const cell = _t.ISMPopUpRunningContainer.getCellById(_t.PopUpContainerConfigData.components.cells[j].id)
                        if (isStart) {
                          cell.setVisible(true)
                          _t.PopUpContainerConfigData.components.cells[j].data.detail.style.visible = true
                        } else {
                          _t.PopUpContainerConfigData.components.cells[j].data.detail.style.visible = false
                          cell.setVisible(false)
                        }
                      }
                      _t.$EventBus.$emit(_t.PopUpContainerConfigData.components.cells[j].data.detail.identifier + "animateEvent", isStart);
                    }
                  }
                }
              }

              //位移动画
              if((typeof _t.PopUpContainerConfigData.components.cells[j].data.detail.animate!="undefined")&&(typeof _t.PopUpContainerConfigData.components.cells[j].data.detail.animate.move!="undefined")&&_t.PopUpContainerConfigData.components.cells[j].data.detail.animate.selected.includes("animateMove"))
              {
                let  conditionx = _t.PopUpContainerConfigData.components.cells[j].data.detail.animate.move.x
                let  conditiony = _t.PopUpContainerConfigData.components.cells[j].data.detail.animate.move.y
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
                      _t.PopUpContainerConfigData.components.cells[j].data.detail.style.position.x = RealValue
                      this.$EventBus.$emit(_t.PopUpContainerConfigData.components.cells[j].data.detail.identifier + "animateMove", animateMove);
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
                    _t.PopUpContainerConfigData.components.cells[j].data.detail.style.position.x = RealValue
                    this.$EventBus.$emit(_t.PopUpContainerConfigData.components.cells[j].data.detail.identifier + "animateMove", animateMove);
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
                      _t.PopUpContainerConfigData.components.cells[j].data.detail.style.position.y = RealValue
                      this.$EventBus.$emit(_t.PopUpContainerConfigData.components.cells[j].data.detail.identifier + "animateMove", animateMove);
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
                    _t.PopUpContainerConfigData.components.cells[j].data.detail.style.position.y = RealValue
                    this.$EventBus.$emit(_t.LayerContainerData.components.cells[j].data.detail.identifier + "animateMove", animateMove);
                  }
                }
              }

              if((typeof _t.PopUpContainerConfigData.components.cells[j].data.detail.active!="undefined"))
              {
                let  active = _t.PopUpContainerConfigData.components.cells[j].data.detail.active
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
                        const cell = _t.ISMPopUpRunningContainer.getCellById(_t.PopUpContainerConfigData.components.cells[j].id)
                        const isEdge = cell.isEdge()
                        if(!isEdge) {
                          _t.$EventBus.$emit(_t.PopUpContainerConfigData.components.cells[j].data.detail.identifier + "activeEvent", activeData);
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
                      const cell = _t.ISMPopUpRunningContainer.getCellById(_t.PopUpContainerConfigData.components.cells[j].id)
                      const isEdge = cell.isEdge()
                      if(!isEdge) {
                        _t.$EventBus.$emit(_t.PopUpContainerConfigData.components.cells[j].data.detail.identifier + "activeEvent", activeData);
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
      async showPage(linkInfo) {
        //  this.PopUpDialog = false
        if(typeof linkInfo.isPopUp!='undefined'&& linkInfo.isPopUp==true)
        {
          let _t = this
          this.IsAutoClose = linkInfo.autoClose
          if (linkInfo.linkType == "Inside") {
            if (this.currentPopUpDisplayUUID === linkInfo.Inside.displayUUID &&
                this.currentPopUpPageUUID === linkInfo.Inside.pageUUID) {
              this.cancelPendingPageLoading()
              return
            }
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
            _t.destroyPopUpGraph()
            this.CurrentPopRealUUIDList=[]
            this.CurrentPopModelUUIDList=[]
            _t.CurrentPagerPopRealDeviceUuidList=[]
            _t.PopUpContainerConfigData = []
            this.selectPopUpPagerContainerDisplayPageDataStruct({
              page: page, callback: function (res,pagerData,uuids,devices) {
                if (_t._isDestroyed) { _t.closePageLoading(loadingKey, loadingToken); return }
                if(res==0)
                {
                  _t.CurrentPagerPopRealDataUuidList = uuids
                  _t.CurrentPagerPopRealDeviceUuidList=devices
                  _t.PopUpContainerConfigData = pagerData
                  _t.chargePagePopUp = false
                  _t.PopUpDialog = true
                  _t.$nextTick(function () {
                    _t.trackTimeout(function (){
                      _t.destroyPopUpGraph()
                      _t.RunPagerPopUpCavasContainerInit(pagerData.layer.autoSize,pagerData.layer.Padding, loadingToken)
                    },100, 'popup')
                  })
                }
                else
                {
                  _t.currentPopUpDisplayUUID = ""
                  _t.currentPopUpPageUUID = ""
                  _t.closePageLoading(loadingKey, loadingToken)
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
                this.PopUpContainerConfigData.PageName = "外部网页"
              }
              else
              {
                this.PopUpContainerConfigData.PageName = linkInfo.title
              }
            }
            else {
              this.PopUpContainerConfigData.PageName = "外部网页"
            }
          }
        }
        else {
          if (linkInfo.linkType == "Inside") {
            if (this.currentDisplayUUID === linkInfo.Inside.displayUUID &&
                this.currentPageUUID === linkInfo.Inside.pageUUID) {
              this.cancelPendingPageLoading()
              return
            }
            this.currentDisplayUUID = linkInfo.Inside.displayUUID
            this.currentPageUUID = linkInfo.Inside.pageUUID
            let _t = this
            let page = {
              pageType: 1,
              displayUUID: linkInfo.Inside.displayUUID,
              pageUuid: linkInfo.Inside.pageUUID
            }

            const loadingToken = await this.consumePendingPageLoading()
            _t.destroyMainGraph()
            _t.clearPendingTimers('main')
            _t.CurrentRealUUIDList=[]
            _t.CurrentModelUUIDList=[]
            _t.CurrentPagerRealDataUuidList = []
            _t.CurrentPagerRealDeviceUuidList=[]
            _t.chargePage = true
            this.selectDisplayPageContainerDataStruct({
              page: page, callback: function (no,pageData,uuids,devices) {
                if (_t._isDestroyed) { _t.closePageLoading(loadingKey, loadingToken); return }
                if(no==0){
                  _t.chargePage = false
                  _t.currentDisplayUUID = page.displayUUID
                  _t.currentPageUUID = page.pageUuid
                  _t.LayerContainerData = pageData
                  _t.CurrentPagerRealDataUuidList = uuids
                  _t.CurrentPagerRealDeviceUuidList=devices
                  _t.RunPagerCavasContainerInit(pageData.layer.autoSize,pageData.layer.Padding, loadingToken)
                } else {
                  _t.chargePage = false
                  _t.currentDisplayUUID = ""
                  _t.currentPageUUID = ""
                  _t.closePageLoading(loadingKey, loadingToken)
                  _t.$message.error(_t.$t("readData.NotFindPage"))
                }
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
      },
      RunPagerCavasContainerInit(auto,panning,loadingToken = this.pageLoadingToken){
        let _t = this
        let pagerSize = auto==1?true:false
        let pagerPanning = panning==1?true:false
        this.pageRenderToken += 1
        if (this.destroyMainGraphRaf) {
          cancelAnimationFrame(this.destroyMainGraphRaf)
          this.destroyMainGraphRaf = null
        }
        if(this.ISMPagerRuningCavasContainer!=null)
        {
          try {
            this.ISMPagerRuningCavasContainer.dispose();
          } catch(e) {
            console.error('[ViewPagerContainer.RunPagerCavasContainerInit] dispose failed', e)
          }
          this.ISMPagerRuningCavasContainer = null
        }
        let refObj = this.detail.identifier
        let view = this.$refs[refObj]
        if (typeof view=="undefined")
        {
          return
        }
        this.ISMPagerRuningCavasContainer = new Graph({
          container: view,
          width: _t.detail.style.position.w,
          height: _t.detail.style.position.h,
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
            color: _t.LayerContainerData.layer.backColor,   // 背景底色（可选）
            image: _t.LayerContainerData.layer.backgroundImage,
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
        window.removeEventListener('resize', this.setScale);
        window.addEventListener('resize', this.setScale);
        const renderToken = this.pageRenderToken
        if(!_t.editMode) {
          //
          this.ISMPagerRuningCavasContainer.on('node:click', ({e, x, y, node, view}) => {
            e.stopPropagation()
            const component = node.prop().data.detail
            _t.SelectCurrentNodeData = component
            // e.preventDefault()
            //  _t.doComponentAction(component,'click')
          });
          this.ISMPagerRuningCavasContainer.on('node:dblclick', ({e, x, y, node, view}) => {
            e.stopPropagation()
            const component = node.prop().data.detail
            _t.SelectCurrentNodeData = component
            _t.doComponentAction(component, 'dblclick')
          });
          this.ISMPagerRuningCavasContainer.on('node:mousedown', ({e, x, y, node, view}) => {
            e.stopPropagation()
            const component = node.prop().data.detail
            _t.SelectCurrentNodeData = component
            _t.doComponentAction(component, 'mousedown')
          });
          this.ISMPagerRuningCavasContainer.on('node:mouseup', ({e, x, y, node, view}) => {
            e.stopPropagation()
            const component = node.prop().data.detail
            _t.SelectCurrentNodeData = component
            _t.doComponentAction(component, 'mouseup')
          });
          this.ISMPagerRuningCavasContainer.on('node:mouseenter', ({e, x, y, node, view}) => {
            e.stopPropagation()
            const component = node.prop().data.detail
            _t.SelectCurrentNodeData = component
            _t.doComponentAction(component, 'mouseenter')
          });
          this.ISMPagerRuningCavasContainer.on('node:mouseleave', ({e, x, y, node, view}) => {
            e.stopPropagation()
            const component = node.prop().data.detail
            _t.SelectCurrentNodeData = component
            _t.doComponentAction(component, 'mouseleave')
          });
          //页面空白处右键菜单
          this.ISMPagerRuningCavasContainer.on('blank:contextmenu', ({e, x, y, cell, view}) => {
            e.preventDefault()
            this.onContextLayerMenu(e)
          })
          this.ISMPagerRuningCavasContainer.on('cell:contextmenu', ({e, x, y, cell, view}) => {
            e.preventDefault()
            this.onContextLayerMenu(e)
          })
        }
        // ==============================
        // 监听页面是否渲染完成，完成后就发广播给节点，告诉节点是在编辑模式
        this.ISMPagerRuningCavasContainer.on('render:done', () => {
          if (renderToken !== this.pageRenderToken || this._isDestroyed) return
          this.closePageLoading(loadingKey, loadingToken)
          this.tempAutoSize=this.LayerContainerData.layer.autoSize
          this.setScale()
          this.$EventBus.$emit('cell-editMode',{
            edit:false,
            toolbox:false,
            source: 'ViewPagerContainer'
          })
          this.$EventBus.$emit('cell-vuex',this.ParentVuex)
          this.InitPagerRealData()
        });
        try{
          const components = JSON.parse(JSON.stringify(_t.LayerContainerData.components))
          if (components.cells && Array.isArray(components.cells)) {
            components.cells = components.cells.filter(cell => cell && cell.shape)
          }
          _t.ISMPagerRuningCavasContainer.fromJSON(components)
        }catch (e){
          _t.closePageLoading(loadingKey, loadingToken)
          _t.$message.error(_t.$t('Render.GetPageError'))
          console.log(e)
        }
      },
      RunPagerPopUpCavasContainerInit(auto,panning,loadingToken = this.pageLoadingToken){
        let _t = this
        const refObj = this.detail ? this.detail.identifier + "ISMPopUpRunningContainer" : null
        if (!this.PopUpDialog || !refObj || !this.$refs[refObj]) {
          this.$nextTick(function () {
            _t.trackTimeout(function () {
              if (_t.PopUpDialog) {
                _t.RunPagerPopUpCavasContainerInit(auto, panning, loadingToken)
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
            this.ISMPopUpRunningContainer.dispose();
          } catch(e) {
            console.error('[ViewPagerContainer.RunPagerPopUpCavasContainerInit] dispose failed', e)
          }
          this.ISMPopUpRunningContainer = null
        }
        let view = this.$refs[refObj]
        this.ISMPopUpRunningContainer = new Graph({
          container: view,
          width: _t.PopUpContainerConfigData.layer.width,
          height: _t.PopUpContainerConfigData.layer.height,
          panning: pagerPanning,
          mousewheel: pagerPanning,
          autoResize:pagerSize,
          grid: false,
          background: {
            color: _t.PopUpContainerConfigData.layer.backColor,   // 背景底色（可选）
            image: _t.PopUpContainerConfigData.layer.backgroundImage,
            size: '100% 100%',
            repeat: 'no-repeat',
            quality:1,
          },
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
        const popRenderToken = this.pageRenderToken
        this.ISMPopUpRunningContainer.on('scale', ( scale ) => {

        });
        //
        this.ISMPopUpRunningContainer.on('node:click', ( { e, x, y, node, view } ) => {
          const component = node.prop().data.detail
          _t.SelectCurrentNodeData = component
          _t.doComponentAction(component,'click')
        });
        this.ISMPopUpRunningContainer.on('node:dblclick', ( { e, x, y, node, view } ) => {
          const component = node.prop().data.detail
          _t.SelectCurrentNodeData = component
          _t.doComponentAction(component,'dblclick')
        });
        this.ISMPopUpRunningContainer.on('node:mousedown', ( { e, x, y, node, view } ) => {
          const component = node.prop().data.detail
          _t.SelectCurrentNodeData = component
          _t.doComponentAction(component,'mousedown')
        });
        this.ISMPopUpRunningContainer.on('node:mouseup', ( { e, x, y, node, view } ) => {
          const component = node.prop().data.detail
          _t.SelectCurrentNodeData = component
          _t.doComponentAction(component,'mouseup')
        });
        this.ISMPopUpRunningContainer.on('node:mouseenter', ( { e, x, y, node, view } ) => {
          const component = node.prop().data.detail
          _t.SelectCurrentNodeData = component
          _t.doComponentAction(component,'mouseenter')
        });
        this.ISMPopUpRunningContainer.on('node:mouseleave', ( { e, x, y, node, view } ) => {
          const component = node.prop().data.detail
          _t.SelectCurrentNodeData = component
          _t.doComponentAction(component,'mouseleave')
        });
        //页面空白处右键菜单
        this.ISMPopUpRunningContainer.on('blank:contextmenu', ({ node }) => {
          this.onContextLayerMenu()
        })
        // ==============================
        // 监听页面是否渲染完成，完成后就发广播给节点，告诉节点是在编辑模式
        this.ISMPopUpRunningContainer.on('render:done', () => {
          if (popRenderToken !== this.pageRenderToken || this._isDestroyed) return
          this.closePageLoading(loadingKey, loadingToken)
          this.InitPopUpPagerRealData()
          this.$EventBus.$emit('cell-editMode',{
            edit:false,
            toolbox:false,
            source: 'ViewPagerContainer'
          })
          this.$EventBus.$emit('cell-vuex',this.ParentVuex)
        });
        //==============================
        try{
          const components = JSON.parse(JSON.stringify(_t.PopUpContainerConfigData.components))
          if (components.cells && Array.isArray(components.cells)) {
            components.cells = components.cells.filter(cell => cell && cell.shape)
          }
          _t.ISMPopUpRunningContainer.fromJSON(components)
        }catch (e){
          _t.closePageLoading(loadingKey, loadingToken)
          _t.$message.error(_t.$t('Render.GetPageError'))
          console.log(e)
        }
      },
      getPagerData(pageid){
        let _t = this
        this.chargePage = true
        _t.CurrentPagerRealDataUuidList = []
        _t.CurrentPagerRealDeviceUuidList=[]
        this.$message.loading({content: 'Loading...', key: loadingKey, duration: 0});
        this.getLayerPagerContainerDataStruct({
          pageid: pageid, cb: function (errno, pagerData,datauuid,devices) {
            if (_t._isDestroyed) { _t.$message.destroy(loadingKey); return }
            _t.chargePage = false
            if(errno==0)
            {
              _t.LayerContainerData = pagerData
              _t.CurrentPagerRealDataUuidList = datauuid
              _t.CurrentPagerRealDeviceUuidList=devices
             _t.RunPagerCavasContainerInit(pagerData.layer.autoSize,pagerData.layer.Padding)
            }
            else
            {
              _t.$message.destroy(loadingKey)
            }
          }
        });
      }
    },
    mounted() {
      let _t = this
      let _this = this
      this.$nextTick(function(){
        _t.SetForm = _t.$form.createForm(this)
        _t.ActionFormPassword = _t.$form.createForm(this)
        _t.SetFormPassword = _t.$form.createForm(this)
        this.initComponents(this.detail);
        let activeEvent = this.detail.identifier+"activeEvent"//动作数据
        let animateEvent = this.detail.identifier+"animateEvent"//动作数据

        _t.eventHandlers.activeEvent = (data) => {


        }
        _t.$EventBus.$on(activeEvent, _t.eventHandlers.activeEvent)
        _t.eventHandlers.animateEvent = (data) => {
          _t.isStart = data
        }
        _t.$EventBus.$on(animateEvent, _t.eventHandlers.animateEvent)
        _t.eventHandlers.MenuConfigPage = (data) => {
          if(_t.LaunchType==1) {
            let wsData = data
            _t.getPagerData(wsData.PageID)
          }
        }
        _t.$EventBus.$on("MenuConfigPage", _t.eventHandlers.MenuConfigPage);
        _t.eventHandlers.onContainerSelectDevice = async (device)=>{
          if(_t.LaunchType!=1) {
            _t.cancelPendingPageLoading()
            return
          }

          if(device.type==1)
          {
            _t.showAppUUID = device.showUUID
            _t.showPageUUID = device.showPageUUID
            _t.SelectDeviceUuid = device.key

            let page={
              pageType:1,
              displayUUID:_t.showAppUUID,
              pageUuid:_t.showPageUUID
            }
            if((typeof device.isPopUp!='undefined')&&(device.isPopUp))
            {
              _t.isExternUrl = false
              _t.chargePagePopUp = true
              const loadingToken = await _t.consumePendingPageLoading()
              _t.destroyPopUpGraph()
              _t.clearPendingTimers('popup')
              _t.CurrentPopRealUUIDList=[]
              _t.CurrentPopModelUUIDList=[]
              _t.CurrentPagerPopRealDataUuidList=[]
              _t.CurrentPagerPopRealDeviceUuidList=[]
              _t.PopUpContainerConfigData = []
              _t.selectPopUpPagerContainerDisplayPageDataStruct({
                page: page, callback: function (res,pagerData,uuids,devices) {
                  if (_t._isDestroyed) { _t.closePageLoading(loadingKey, loadingToken); return }
                  if(res==0)
                  {
                    _t.CurrentPagerPopRealDataUuidList = uuids
                    _t.CurrentPagerPopRealDeviceUuidList=devices
                    _t.PopUpContainerConfigData = pagerData
                    _t.CurrentPagerPopRealDeviceUuidList.push(device.key)
                    _t.chargePagePopUp = false
                    _t.PopUpDialog = true
                    _t.$nextTick(function () {
                      _t.trackTimeout(function (){
                        _t.destroyPopUpGraph()
                        _t.RunPagerPopUpCavasContainerInit(pagerData.layer.autoSize,pagerData.layer.Padding, loadingToken)
                      },100, 'popup')
                    })
                  }
                  else
                  {
                    _t.currentPopUpDisplayUUID = ""
                    _t.currentPopUpPageUUID = ""
                    _t.chargePagePopUp = false
                    _t.closePageLoading(loadingKey, loadingToken)
                    _t.$message.error(_t.$t("readData.NotFindPage"))
                  }
                }
              })
            }
            else {
              _t.currentDisplayUUID = page.displayUUID
              _t.currentPageUUID = page.pageUuid
              const loadingToken = await _t.consumePendingPageLoading()
              _t.destroyMainGraph()
              _t.clearPendingTimers('main')
              _t.chargePage = true
              _t.CurrentRealUUIDList=[]
              _t.CurrentModelUUIDList=[]
              _t.CurrentPagerRealDataUuidList=[]
              _t.CurrentPagerRealDeviceUuidList=[]
              _t.selectDisplayPageContainerDataStruct({
                page: page, callback: function (no,pageData,uuids,devices) {
                  if (_t._isDestroyed) { _t.closePageLoading(loadingKey, loadingToken); return }
                  _t.chargePage = false
                  if(no==0){
                    _t.LayerContainerData = pageData
                    _t.CurrentPagerRealDataUuidList = uuids
                    _t.CurrentPagerRealDeviceUuidList=devices
                    _t.CurrentPagerRealDeviceUuidList.push(device.key)
                    _t.RunPagerCavasContainerInit(pageData.layer.autoSize,pageData.layer.Padding, loadingToken)
                  } else {
                    _t.currentDisplayUUID = ""
                    _t.currentPageUUID = ""
                    _t.closePageLoading(loadingKey, loadingToken)
                    _t.$message.error(_t.$t("readData.NotFindPage"))
                  }
                }
              })
            }
          }
        }
        _t.$EventBus.$on("onContainerSelectDevice", _t.eventHandlers.onContainerSelectDevice)
        _t.eventHandlers.readDataPush = (data) => {
          _t.DealWithUpdateData(data)
        }
        _t.$EventBus.$on("readDataPush", _t.eventHandlers.readDataPush);
        _t.eventHandlers.StaticData = (data) => {
          _t.DealWithUpdateData(data)
        }
        _t.$EventBus.$on("StaticData", _t.eventHandlers.StaticData);
        _t.eventHandlers.SystemData = (data) => {
          if(!_t.chargePage) {
            let realData = data
            let projectUuid = getAuthorization(AUTH_TYPE.AUTH1)
            for (let k = 0; k < realData.Data.length; k++) {
              for (let j = 0; j < _t.LayerContainerData.components.length; j++) {
                if ((typeof _t.LayerContainerData.components.cells[j].data.detail.animate != "undefined") && (typeof _t.LayerContainerData.components.cells[j].data.detail.animate.condition != "undefined")) {
                  let condition = _t.LayerContainerData.components.cells[j].data.detail.animate.condition
                  let isExpression = _t.LayerContainerData.components.cells[j].data.detail.animate.isExpression
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
                        _t.$EventBus.$emit(_t.LayerContainerData.components.cells[j].data.detail.identifier + "animateEvent", isStart);
                      }
                    }
                  }
                }
                if((typeof _t.LayerContainerData.components.cells[j].data.detail.animate!="undefined")&&(typeof _t.LayerContainerData.components.cells[j].data.detail.animate.move!="undefined")&&_t.LayerContainerData.components.cells[j].data.detail.animate.selected.includes("animateMove"))
                {
                  let  conditionx = _t.LayerContainerData.components.cells[j].data.detail.animate.move.x
                  let  conditiony = _t.LayerContainerData.components.cells[j].data.detail.animate.move.y
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
                        _t.LayerContainerData.components.cells[j].data.detail.style.position.x = RealValue
                        this.$EventBus.$emit(_t.LayerContainerData.components.cells[j].data.detail.identifier + "animateMove", animateMove);
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
                      _t.LayerContainerData.components.cells[j].data.detail.style.position.x = RealValue
                      this.$EventBus.$emit(_t.LayerContainerData.components.cells[j].data.detail.identifier + "animateMove", animateMove);
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
                        _t.LayerContainerData.components.cells[j].data.detail.style.position.y = RealValue
                        this.$EventBus.$emit(_t.LayerContainerData.components.cells[j].data.detail.identifier + "animateMove", animateMove);
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
                      _t.LayerContainerData.components.cells[j].data.detail.style.position.y = RealValue
                      this.$EventBus.$emit(_t.LayerContainerData.components.cells[j].data.detail.identifier + "animateMove", animateMove);
                    }
                  }
                }
                if ((typeof _t.LayerContainerData.components.cells[j].data.detail.active != "undefined")) {
                  let active = _t.LayerContainerData.components.cells[j].data.detail.active
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
                          _t.$EventBus.$emit(_t.LayerContainerData.components.cells[j].data.detail.identifier + "activeEvent", activeData);
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
            for (let k = 0; k < realData.Data.length; k++) {
              for (let j = 0; j < _t.PopUpContainerConfigData.components.length; j++) {
                if ((typeof _t.PopUpContainerConfigData.components[j].animate != "undefined") && (typeof _t.PopUpContainerConfigData.components[j].animate.condition != "undefined")) {
                  let condition = _t.PopUpContainerConfigData.components[j].animate.condition
                  let isExpression = _t.PopUpContainerConfigData.components[j].animate.isExpression
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
                        _t.$EventBus.$emit(_t.PopUpContainerConfigData.components[j].identifier + "animateEvent", isStart);
                      }
                    }
                  }
                }

                //位移动画
                if((typeof _t.PopUpContainerConfigData.components[j].animate!="undefined")&&(typeof _t.PopUpContainerConfigData.components[j].animate.move!="undefined")&&_t.PopUpContainerConfigData.components[j].animate.selected.includes("animateMove"))
                {
                  let  conditionx = _t.PopUpContainerConfigData.components[j].animate.move.x
                  let  conditiony = _t.PopUpContainerConfigData.components[j].animate.move.y
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
                        _t.LayerContainerData.components.cells[j].data.detail.style.position.x = RealValue
                        this.$EventBus.$emit(_t.PopUpContainerConfigData.components[j].identifier + "animateMove", animateMove);
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
                      _t.LayerContainerData.components.cells[j].data.detail.style.position.x = RealValue
                      this.$EventBus.$emit(_t.PopUpContainerConfigData.components[j].identifier + "animateMove", animateMove);
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
                        _t.PopUpContainerConfigData.components[j].style.position.y = RealValue
                        this.$EventBus.$emit(_t.PopUpContainerConfigData.components[j].identifier + "animateMove", animateMove);
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
                      _t.PopUpContainerConfigData.components[j].style.position.y = RealValue
                      this.$EventBus.$emit(_t.LayerContainerData.components.cells[j].data.detail.identifier + "animateMove", animateMove);
                    }
                  }
                }

                if ((typeof _t.PopUpContainerConfigData.components[j].active != "undefined")) {
                  let active = _t.PopUpContainerConfigData.components[j].active
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
                          _t.$EventBus.$emit(_t.PopUpContainerConfigData.components[j].identifier + "activeEvent", activeData);
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }
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
        }
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
          }
        _t.$EventBus.$on("ChargePage", _t.eventHandlers.ChargePage);
      });
    },
    created(){
      let _t = this
      this.GetNodeObj = this.getNode()
      this.nodeHandlers.changeData = ({ current }) => {
        if(current) {
          _t.detail = current.detail
        }
      }
      this.nodeHandlers.changeSize = ({ current }) => {
        _t.detail.style.position.w = current.width
        _t.detail.style.position.h = current.height
      }
      this.GetNodeObj.on('change:data', this.nodeHandlers.changeData)
      this.GetNodeObj.on('change:size', this.nodeHandlers.changeSize);
      this.detail = this.GetNodeObj.getData().detail
      this.editMode = this.GetNodeObj.getData().editMode
      this.showDeviceUuid = this.GetNodeObj.getData().showDeviceUuid
      this.IsToolBox = this.GetNodeObj.getData().IsToolBox
      this.eventHandlers.cellEditMode = (data) => {
        _t.editMode = data.edit
        _t.IsToolBox = data.toolbox
        // _t.initComponents(_t.detail)
      }
      _t.$EventBus.$on('cell-editMode', this.eventHandlers.cellEditMode)
      this.eventHandlers.cellVuex = (data) => {
        this.ParentVuex = data
        this.PCPageList = data.PStore.state.ISMDisPlayEditorTool.PCPageList
        this.PhonePageList = data.PStore.state.ISMDisPlayEditorTool.PhonePageList
      }
      this.$EventBus.$on('cell-vuex', this.eventHandlers.cellVuex)
    },
    beforeDestroy() {
      this._isDestroyed = true
      this.clearPendingTimers()
      // 清理 DOM 事件
      window.removeEventListener('resize', this.setScale)
      this.removeListener()
      // 清理 X6 node 事件监听
      if (this.GetNodeObj) {
        try { this.GetNodeObj.off('change:data', this.nodeHandlers.changeData) } catch(e) { console.warn(e) }
        try { this.GetNodeObj.off('change:size', this.nodeHandlers.changeSize) } catch(e) { console.warn(e) }
      }
      this.nodeHandlers = {}
      // 清理 EventBus 事件（精确 $off，只移除自身监听，不影响其他组件）
      try {
        const h = this.eventHandlers
        if (this.detail && this.detail.identifier) {
          this.$EventBus.$off(this.detail.identifier + 'activeEvent', h.activeEvent)
          this.$EventBus.$off(this.detail.identifier + 'animateEvent', h.animateEvent)
        }
        this.$EventBus.$off('cell-editMode', h.cellEditMode)
        this.$EventBus.$off('cell-vuex', h.cellVuex)
        if (h.RealAlarm) this.$EventBus.$off('RealAlarm', h.RealAlarm)
        if (h.MenuConfigPage) this.$EventBus.$off('MenuConfigPage', h.MenuConfigPage)
        if (h.onContainerSelectDevice) this.$EventBus.$off('onContainerSelectDevice', h.onContainerSelectDevice)
        if (h.readDataPush) this.$EventBus.$off('readDataPush', h.readDataPush)
        if (h.StaticData) this.$EventBus.$off('StaticData', h.StaticData)
        if (h.SystemData) this.$EventBus.$off('SystemData', h.SystemData)
        if (h.ChargePage) this.$EventBus.$off('ChargePage', h.ChargePage)
        this.eventHandlers = {}
      } catch(e) {
        console.error('[ViewPagerContainer.beforeDestroy] EventBus cleanup failed', e)
      }
      // 清理嵌套 X6 Graph（含异常保护）
      this.destroyPopUpGraph()
      this.destroyMainGraph()
      this.unmountBodyPageLoading()
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

</style>
<style lang="less" scoped>
::v-deep .x6-edge-tool-vertex {
  display: none !important;
}
::v-deep .ant-modal-root div[aria-hidden="true"] {
  display: none !important;
}
::v-deep .x6-node:hover {
  cursor: var(--NodeMouseStyle);
}
::v-deep .pager-graph-container {
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
  width : 5px;  /*高宽分别对应横竖滚动条的尺寸*/
  height: 9px;
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
