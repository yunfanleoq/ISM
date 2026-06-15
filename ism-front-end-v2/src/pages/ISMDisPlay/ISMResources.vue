<template>
  <div style="border-right:0px solid #95B8E7; overflow: hidden;
            overflow-y: auto;">
    <a-collapse v-model="activeKey"  :bordered="false" :expand-icon-position="expandIconPosition">
      <a-collapse-panel key="1" :header="$t('ISMResources.PCHeader')" :style="customStyle">
        <a-tooltip placement="top" slot="extra">
          <template slot="title">
            <span>{{$t('ISMResources.NewPage')}}</span>
          </template>
          <icon-font  style="font-size: 15px" type="icon-xinjian3" @click="addModelPage($event,1)" />
        </a-tooltip>

        <a-tree :multiple="false" :selectedKeys="PCPageCheckKey" :treeData="PCPageList" :blockNode="true" @select="chargePage">
              <template slot="custom" slot-scope="item">
                <div v-if="item.isEdit&&selectTreeKey==item.key">
                  <a-input v-model="item.title" style="width: 180px;height: 24px" />
                  <span class="tree-save_icon " @click="SureEdit(item)">
                    <a-icon type="check-circle" />
                  </span>
                  <span class="tree-cancle_icon " @click="cancelEdit(item)">
                  <a-icon type="close-circle" />
                  </span>
                </div>
                <div v-else>
                    <span class="node-title" :title="item.title">
                        <template slot="title">
                          <span>{{$t('ISMResources.PageTipsIsHome')}}</span>
                        </template>
                      <a-icon type="home" v-if="item.IsHome==1"/>
                      <span class="empty-home-ico" v-else></span>
                      {{item.title}}
                    </span>
                    <span class="" style="margin-left: 0px" v-if="selectTreeKey==item.key">
                      <a-tooltip placement="top" v-if="item.IsLogin!=1">
                        <template slot="title">
                          <span>{{$t('ISMResources.PageTipsHome')}}</span>
                        </template>
                       <span class="icon-wrap" @click="setHomePage(item)">
                            <a-icon type="home" />
                       </span>
                      </a-tooltip>
                      <a-tooltip placement="top">
                        <template slot="title">
                          <span>{{$t('ISMResources.PageTipsEdit')}}</span>
                        </template>
                        <span class="icon-wrap" @click="editPage(item)">
                          <a-icon type="form" />
                        </span>
                       </a-tooltip>
                      <a-tooltip placement="top">
                        <template slot="title">
                          <span>{{$t('ISMResources.PageTipsCopy')}}</span>
                        </template>
                        <span class="icon-wrap" @click="copyPage(item)">
                          <a-icon type="copy" />
                        </span>
                       </a-tooltip>
                      <a-tooltip placement="top">
                          <template slot="title">
                            <span>{{$t('ISMResources.PageTipsDel')}}</span>
                          </template>
                          <span class="icon-wrap" @click="delPage(item)">
                            <a-icon type="delete" />
                          </span>
                       </a-tooltip>
                    </span>
                  </div>
            </template>
        </a-tree>
      </a-collapse-panel>
      <a-collapse-panel key="2" :header="$t('ISMResources.PhoneHeader')" :style="customStyle">
        <a-tooltip placement="top" slot="extra">
          <template slot="title">
            <span>{{$t('ISMResources.NewPage')}}</span>
          </template>
          <icon-font  style="font-size: 15px" type="icon-xinjian3" @click="addModelPage($event,0)" />
        </a-tooltip>
        <a-tree  :multiple="false" :selectedKeys="PhonePageCheckKey" :treeData="PhonePageList" :blockNode="true" @select="chargePage">
          <template slot="custom" slot-scope="item">
            <div v-if="item.isEdit&&selectTreeKey==item.key">
              <a-input v-model="item.title" style="width: 180px;height: 24px" />
              <span class="tree-save_icon " @click="SureEdit(item)">
                    <a-icon type="check-circle" />
                  </span>
              <span class="tree-cancle_icon " @click="cancelEdit(item)">
                  <a-icon type="close-circle" />
                  </span>
            </div>
            <div v-else >
              <span class="node-title" :title="item.title">
                 <a-icon type="home" v-if="item.IsHome==1"/>
                 <span class="empty-home-ico" v-else></span>
                {{item.title}}
              </span>
              <span class="" style="margin-left: 0px" v-if="selectTreeKey==item.key">
                <a-tooltip placement="top">
                  <template slot="title">
                    <span>{{$t('ISMResources.PageTipsHome')}}</span>
                  </template>
                 <span class="icon-wrap" @click="setHomePage(item)">
                      <a-icon type="home" />
                 </span>
                </a-tooltip>
                 <a-tooltip placement="top">
                  <template slot="title">
                    <span>{{$t('ISMResources.PageTipsEdit')}}</span>
                  </template>
                  <span class="icon-wrap" @click="editPage(item)">
                    <a-icon type="form" />
                  </span>
                 </a-tooltip>
                 <a-tooltip placement="top">
                  <template slot="title">
                    <span>{{$t('ISMResources.PageTipsDel')}}</span>
                  </template>
                  <span class="icon-wrap" @click="delPage(item)">
                    <a-icon type="delete" />
                  </span>
                 </a-tooltip>
                </span>
            </div>
          </template>
        </a-tree>
      </a-collapse-panel>
    </a-collapse>
    <a-modal v-drag-modal :destroyOnClose="true" v-model="AddPageVisible" :title="$t('ISMResources.NewPage')">
      <template slot="footer">
        <a-button type="primary"  @click="addModelPageRequest">
          {{ $t('displayModel.ModelSure') }}
        </a-button>
        <a-button  @click="AddPageVisible=false">
          {{ $t('displayModel.ModelCancel') }}
        </a-button>
      </template>
      <a-form :form="AddModelPageForm" :label-col="{ span: 5 }" :wrapper-col="{ span: 18 }">
        <a-form-item :label="$t('ISMResources.PageName')">
          <a-input
              v-decorator="['PageName', { rules: [{ required: true, message: $t('ISMResources.PageName') }] }]"
          />
        </a-form-item>
        <a-form-item :label="$t('ISMResources.PageType')">
          <a-select v-decorator="['PageType', { initialValue:0,rules: [{ required: true, message: $t('ISMResources.PageType') }] }]">
            <a-select-option  v-for="(type,index) in AddPageType" :key="index" :value=type.value>
              {{ $t(type.name) }}
            </a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item :label="$t('ISMResources.PageSize')">
          <a-select @change="handleSelectPageSizeChange"   autocomplete="autocomplete" v-if="pageType==1"
                    v-decorator="['PageSize', {rules: [{ required: true, message: $t('ISMResources.PageSize'), whitespace: true}]}]"
          >
            <a-select-option value="1">1920*1080 {{$t('ISMResources.PageSizeFirst')}}</a-select-option>
            <a-select-option value="2">1060*520</a-select-option>
            <a-select-option value="3">1440*900</a-select-option>
            <a-select-option value="4">1366*768</a-select-option>
          </a-select>
          <a-select @change="handleSelectPageSizeChange"   autocomplete="autocomplete" v-if="pageType==0"
                    v-decorator="['PageSize', {rules: [{ required: true, message: $t('ISMResources.PageSize'), whitespace: true}]}]"
          >
            <a-select-opt-group>
              <span slot="label" style="font-size: 14px"><a-icon type="apple" /> iPhone</span>
              <a-select-option value="5">iPhone 5/SE</a-select-option>
              <a-select-option value="6">iPhone XR</a-select-option>
              <a-select-option value="7">iPhone 12 Pro</a-select-option>
              <a-select-option value="8">iPad</a-select-option>
              <a-select-option value="9">iPad Pro</a-select-option>
              <a-select-option value="10">iPad Air</a-select-option>
              <a-select-option value="11">iPad Mini</a-select-option>
              <a-select-option value="12">iPhone 6/7/8</a-select-option>
              <a-select-option value="13">iPhone 6/7/8 Plus</a-select-option>
              <a-select-option value="14">iPhone X</a-select-option>
            </a-select-opt-group>
            <a-select-opt-group>
              <span slot="label" style="font-size: 14px"><a-icon type="android" /> Android</span>
                <a-select-option value="15">Samsung Galaxy S8+</a-select-option>
                <a-select-option value="16">Samsung Galaxy S20 Ultra</a-select-option>
                <a-select-option value="17">Samsung Galaxy A15/71</a-select-option>
            </a-select-opt-group>

            <a-select-opt-group>
              <span slot="label" style="font-size: 14px"><a-icon type="tablet" /> {{$t('ISMResources.PageSizeDiy')}}</span>
              <a-select-option value="18">{{$t('ISMResources.PageSizeDiy')}}</a-select-option>
            </a-select-opt-group>
          </a-select>
        </a-form-item>

      </a-form>
    </a-modal>
  </div>
</template>
<script>
import store from "../../store";
import { mapActions, mapGetters, mapState, mapMutations } from 'vuex'

import {
  displayModelDelete,
  DisplayModelPageAdd,
  DisplayModelPageDel,
  DisplayModelPageEdit, DisplayModelPageSetHome,DisplayModelPageCopy
} from "../../services/displayModel";
import {updateLayerDataStruct} from "../../store/ISM/actions";
import {normalizeISMScene} from "@/pages/ISMDisPlay/utils/ismSceneNormalizer";
export default {
  name: 'ISMResources',
  data() {
    return {
      AddPageType:[
        {name:"ISMResources.PageShowType",value:0},
        {name:"ISMResources.PageLoginType",value:1},
      ],
      activeKey: ['1'],
      selectTreeKey:0,
      AddPageVisible:false,
      AddModelPageForm:this.$form.createForm(this),
      customStyle: 'background: #fff;border-radius: 4px;margin-bottom: 1px;border: 0;overflow: hidden',
      codeArray: [],
      pageSize:"",
      pageType:1,
      checkPageInfo:null,
      PCPageCheckKey:[],
      PhonePageCheckKey:[],
      logging:false,
      expandIconPosition: 'left',
    };
  },
  i18n: require('../../i18n/language'),
  computed: {
    ...mapState({
      configObject: state => store.state.ISMDisPlayEditorTool.selectedComponent,
      selectPageUuid: state => store.state.ISMDisPlayEditorTool.selectPageUuid,
      configData: state => store.state.ISMDisPlayEditorTool.LayerData,
      PCPageList: state => store.state.ISMDisPlayEditorTool.PCPageList,
      PhonePageList: state => store.state.ISMDisPlayEditorTool.PhonePageList,
      ISMCavasContainer:state => store.state.ISMDisPlayEditorTool.ISMCavasContainer,
    }),
  },
  watch: {
    configData: {
      handler(newVal) {
        this.updateLayerDataStruct(newVal)
      }
    }
  },
  components: {

  },
  methods: {
    ...mapActions('ISMDisPlayEditorTool',[
      'selectLayerDataStruct',
      'saveLayerDataStruct',
      'setGroupList',
      'getLayerDataStruct',
      'selectParentLayerDataStruct',
      'updateLayerDataStruct',
        'updateAllLayerDataStruct'
    ]),
    ...mapMutations('ISMDisPlayEditorTool',[
      'ClearRedo',
      'ClearUndo',
      'ClearisHistoryOp'
    ]),
    handleSelectPageSizeChange(value){
        this.pageSize = value
    },
    doSaveLayerData(uuid){
      let _t = this
      let params = {
        uuid:uuid,
        pageid:this.selectPageUuid,
        LayerData:null,
      }
      const LayerData = JSON.parse(JSON.stringify(this.configData))
      const normalizedLayerData = normalizeISMScene({
        layer: LayerData.layer,
        components: this.ISMCavasContainer.toJSON()
      })
      LayerData.layer = {
        ...LayerData.layer,
        ...normalizedLayerData.layer
      }
      LayerData.components = normalizedLayerData.components
      params.LayerData = LayerData
      this.saveLayerDataStruct(params).then(function (res){
        if(res.data.code == 200)
        {
          let uid = _t.$route.params.uid
          _t.updateAllLayerDataStruct({pageType:_t.isMobile,uuid:uid,cb:function (){}});
        }
        else
        {
          _t.$message.error(_t.$t('displayModel.SaveDataFailed'))
        }
      })
    },
    chargePage(key,info){
      if(info.selectedNodes.length>0)
      {
        let page = info.selectedNodes[0].data.props.dataRef
        if(page.isComponents)
        {
          this.selectParentLayerDataStruct(key[0])
          this.$nextTick(function(){
            this.$EventBus.$emit("ClickComponentsEvent", key[0]);
          });
        }
        else
        {
          if(key.length>0)
          {
            this.selectTreeKey = key[0]
          }
          if (page.pageType == 1) {
            this.PCPageCheckKey = key;
            this.PhonePageCheckKey = []
          } else {
            this.PhonePageCheckKey = key;
            this.PCPageCheckKey = []
          }
          this.checkPageInfo = page
          //保存上个页面
          this.doSaveLayerData(this.$route.params.uid)
          this.selectLayerDataStruct(page)
          document.title = page.AppName + ' | ' + page.title
        }
        // this.setGroupList()
      }
    },
    addModelPage(e,type) {
      e.stopPropagation()
      this.AddPageVisible=true
      this.pageType = type
    },
    addModelPageRequest(e) {
      let _t = this
      this.AddModelPageForm.validateFields((err) => {
        if (!err) {
          this.logging = true
          const params = {
            modelUuid:this.$route.params.uid,
            name:this.AddModelPageForm.getFieldValue('PageName'),
            isLogin:this.AddModelPageForm.getFieldValue('PageType'),
            size:this.pageSize,
            pageType:parseInt(this.pageType)
          };
          DisplayModelPageAdd(params).then(function (res){
            if (res.data.code == 4002) {
              _t.$message.success(_t.$t('ISMResources.PageAddSuccess'), 3)
              _t.getLayerDataStruct({uuid:_t.$route.params.uid,cb:function (){
                  _t.spinning = false
                }});
              _t.AddPageVisible=false
            }else if (res.data.code == 4005) {
              _t.$message.error(_t.$t('ISMResources.PageCountOut'), 3)
            }else if (res.data.code == 4006) {
              _t.$message.error(_t.$t('ISMResources.LoginPageNoAuth'), 3)
            }
            else{
              _t.$message.error(_t.$t('ISMResources.PageAddFailed'), 3)
            }
          }).catch(function (){
            _t.logging = false
            _t.$message.error(_t.$t('loginPage.serverError'), 3)
          })
        }
      })
    },
    cancelEdit(item){
      if(item.pageType==1) {
        for(let i=0;i<this.PCPageList.length;i++){
            this.PCPageList[i].isEdit=false
        }
      }else{
        for(let i=0;i<this.PhonePageList.length;i++){
            this.PhonePageList[i].isEdit=false
        }
      }
    },
    editPage(item) {
      if(item.pageType==1) {

        for(let i=0;i<this.PCPageList.length;i++){
            if(item.id==this.PCPageList[i].id)
            {
              this.PCPageList[i].isEdit=true
              break
            }
        }
      }else{
        for(let i=0;i<this.PhonePageList.length;i++){
          if(item.id==this.PhonePageList[i].id)
          {
            this.PhonePageList[i].isEdit=true
            break
          }
        }
      }
    },
    setHomePage(item){
      let _t = this
      this.logging = true
      const params = {
        muid:item.pageModelUuid,
        pageid:item.pageUuid,
        pageType:item.pageType,
      };
      DisplayModelPageSetHome(params).then(function (res){
        if (res.data.code == 4002) {
          _t.$message.success(_t.$t('ISMResources.PageUpdateSuccess'), 3)
          if(item.pageType==1) {

            for(let i=0;i<_t.PCPageList.length;i++){
              if(item.id==_t.PCPageList[i].id)
              {
                _t.PCPageList[i].IsHome=1
              }
              else
              {
                _t.PCPageList[i].IsHome=0
              }
            }
          }else{
            for(let i=0;i<_t.PhonePageList.length;i++){
              if(item.id==_t.PhonePageList[i].id)
              {
                _t.PhonePageList[i].IsHome=1
              }
              else{
                _t.PhonePageList[i].IsHome=0
              }
            }
          }
        }else{
          _t.$message.error(_t.$t('ISMResources.PageUpdateFailed'), 3)
        }
      }).catch(function (){
        _t.logging = false
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
    },
    copyPage(item){
      let _t = this
      this.logging = true
      const params = {
        muid:item.pageModelUuid,
        pageid:item.pageUuid,
        pageType:item.pageType,
      };
      DisplayModelPageCopy(params).then(function (res){
        if (res.data.code == 200) {
          _t.$message.success(_t.$t('ISMResources.PageUpdateSuccess'), 3)
          _t.getLayerDataStruct({uuid:_t.$route.params.uid,cb:function (){
              _t.spinning = false
            }});
        }else if (res.data.code == 4005) {
          _t.$message.error(_t.$t('ISMResources.PageCountOut'), 3)
        }else{
          _t.$message.error(_t.$t('ISMResources.PageUpdateFailed'), 3)
        }
      }).catch(function (){
        _t.logging = false
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
    },
    SureEdit(item){
      let _t = this
      this.logging = true
      const params = {
        modelUuid:item.pageModelUuid,
        pageUuid:item.pageUuid,
        name:item.title,
      };
      DisplayModelPageEdit(params).then(function (res){
        if (res.data.code == 4002) {
          _t.$message.success(_t.$t('ISMResources.PageUpdateSuccess'), 3)
          if(item.pageType==1) {

            for(let i=0;i<_t.PCPageList.length;i++){
              if(item.id==_t.PCPageList[i].id)
              {
                _t.PCPageList[i].title=item.title
                _t.PCPageList[i].isEdit=false
                break
              }
            }
          }else{
            for(let i=0;i<_t.PhonePageList.length;i++){
              if(item.id==_t.PhonePageList[i].id)
              {
                _t.PhonePageList[i].title=item.title
                _t.PhonePageList[i].isEdit=false
                break
              }
            }
          }
          _t.selectLayerDataStruct(item)
        }else{
          _t.$message.error(_t.$t('ISMResources.PageUpdateFailed'), 3)
        }
      }).catch(function (){
        _t.logging = false
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
    },
    delPage(item) {
      let _t = this
      this.$confirm({
        content: _t.$t('ISMResources.PageDelSure'),
        okText: _t.$t('displayModel.ConfirmOk'),
        onOk() {
          const params={
            modelUuid:_t.$route.params.uid,
            pageId:item.pageUuid
          }
          DisplayModelPageDel(params).then(function (res){
            _t.logging = false
            if (res.data.code == 4002) {
              _t.$message.success(_t.$t('displayModel.DelModelSuccess'), 3)
              if(item.pageType==1) {

                for(let i=0;i<_t.PCPageList.length;i++){
                  if(item.id==_t.PCPageList[i].id)
                  {
                    _t.PCPageList.splice(i,1)
                    break
                  }
                }
              }else{
                for(let i=0;i<_t.PhonePageList.length;i++){
                  if(item.id==_t.PhonePageList[i].id)
                  {
                    _t.PhonePageList.splice(i,1)
                    break
                  }
                }
              }
            }
            else
            {
              _t.$message.error(_t.$t('displayModel.DelModelFailed'), 3)
            }
          })
        },
        cancelText: _t.$t('displayModel.ConfirmCancel'),
        onCancel() {

        },
      });
    }
  },
};
</script>
<style lang="less" scoped>
::v-deep .ant-tree-switcher.ant-tree-switcher_open {
  .icon-plus {
    //background-image: url("./minus.png") ; // 展开节点时的icon
    background-size: 24px;
    width: 24px;
    height: 24px;
  }
}
::v-deep .ant-tree-switcher.ant-tree-switcher_close {
  .icon-plus {
    //background-image: url("./add.png") ; // 收起节点时的icon
    background-size: 24px;
    width: 24px;
    height: 24px;
  }
}
.empty-home-ico {
  display: inline-block;
  width: 14px;
  height: 14px;
  margin: 0;
  line-height: 14px;
  text-align: center;
  vertical-align: top;
  border: 0 none;
  outline: none;
  cursor: pointer;
}
.icon-wrap {
  margin: 0 6px;
  vertical-align: 0.45em;
}
.node-title {
  width: 180px;
  overflow:hidden;
  text-overflow:ellipsis;
  display: inline-block;
  /*禁止换行显示*/
  white-space:nowrap;
}
.tree-save_icon{
  margin-left: 10px;
}
.tree-cancle_icon {
  margin: 0 10px;
}
</style>
