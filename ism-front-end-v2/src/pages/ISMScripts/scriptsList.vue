<template>
  <a-card>
    <a-space class="operator">
      <a-button @click="addVisible=true;isEdit=false" type="primary" icon="plus">{{$t('dataModel.newModel')}}</a-button>
      <a-button @click="refresh()"  type="default" icon="sync" :loading="refIconLoading">{{$t("dataModel.refModel")}}</a-button>
    </a-space>

    <a-spin style="padding: 1px;"  :spinning="messageShowLoad" tip="Loading...">
      <a-table rowKey="ScriptName" :pagination="pagination" :columns="columns" :data-source="dataSource">
        <template v-for="(item, index) in columns" :slot="item.slotName">
          <span :key="index">{{ $t(item.slotName) }}</span>
        </template>
        <div slot="ScriptType" slot-scope="text" >
          <span v-if="text==0">{{ $t('ISMScripts.ScriptAuto') }}</span>
          <span v-else-if="text==1">{{ $t('ISMScripts.ScriptHandle') }}</span>
        </div>
        <div slot="action" slot-scope="text, record">
          <a @click="doSysScript(record.ScriptUuid)" style="color: #e89924"><icon-font type="icon-gongjulan-zhihang" />{{$t('ISMScripts.ExecScript')}}</a> |
          <a v-if="record.IsDisable==0" @click="DoDisableSysScript(record.ScriptUuid,record.IsDisable)" style="color: #d81e06"><icon-font type="icon-jinzhitishi" />
            <span >{{$t('ISMScripts.DisableScript')}}</span>
          </a>
          <a v-if="record.IsDisable==1" @click="DoDisableSysScript(record.ScriptUuid,record.IsDisable)" style="color: #1296db"><icon-font type="icon-cangpeitubiao_qiyong" />
            <span>{{$t('ISMScripts.EnableScript')}}</span>
          </a> |
          <a @click="GoToEdit(record)" style="color: #13C2C2"><a-icon type="edit" />{{$t('dataModel.modelDetail')}}</a> |
          <a-popconfirm :title="$t('dataModel.deleteConfirm')" @confirm="deleteRecord(record.ScriptUuid)">
            <a-icon slot="icon" type="question-circle-o" style="color: red" />
            <a-icon type="delete" theme="twoTone" two-tone-color="#eb2f96"/><a style="color: #eb2f96">{{$t('dataModel.delete')}}</a>
          </a-popconfirm>
        </div>
      </a-table>
    </a-spin>
    <a-drawer
          :title="isEdit?$t('ISMScripts.editScript'):$t('ISMScripts.addScript')"
          :width="800"
          :visible="addVisible"
          :body-style="{ paddingBottom: '80px' }"
          @close="onClose"
      >
      <a-form :form="PlanForm" layout="vertical" >
        <a-row :gutter="16">
          <a-col :span="8">
            <a-form-item
                :label="$t('ISMScripts.ScriptName')"
            >
              <a-input  autocomplete="autocomplete"

                        v-decorator="['ScriptName', {rules: [{ required: true, message: $t('ISMScripts.ScriptName'), whitespace: true}]}]"
              />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item
                :label="$t('ISMScripts.ScriptType')"
            >
              <a-select  autocomplete="autocomplete" @change="chargeScriptType"

                         v-decorator="['ScriptType', {initialValue:0,rules: [{ type:'number',required: true, message: $t('ISMScripts.ScriptType'), whitespace: true}]}]"
              >
                <a-select-option :key="0" :value=0>{{ $t('ISMScripts.ScriptAuto') }}</a-select-option>
                <a-select-option :key="1" :value=1>{{ $t('ISMScripts.ScriptHandle') }}</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="8" v-if="ScriptType==0">
            <a-form-item
                :label="$t('ISMScripts.ScriptDelay')"
            >
              <a-input  autocomplete="autocomplete"

                        v-decorator="['ScriptDelay', {rules: [{ required: false, message: $t('ISMScripts.ScriptDelay'), whitespace: true}]}]"
              />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="24">
            <a-form-item>
                <span slot="label">
                {{$t('ISMScripts.ScriptContent')}}&nbsp;
                <a-tooltip title="全屏显示">
                  <a-icon type="fullscreen" @click="toggleEditorFullscreen"/>
                </a-tooltip>
              </span>
              <code-editor v-if="isCharge"
                  :value="CodeContent"
                  language="javascript"
                  @input="changeTextarea"
              >
              </code-editor>
              <a-button v-if="isFullscreen" class="fullscreen-btn" @click="toggleEditorFullscreen">
               退出全屏
              </a-button>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="24">
          <a-form-item
            :label="$t('ISMScripts.ScriptDes')"
        >
          <a-textarea autocomplete="autocomplete"

                      v-decorator="['ScriptDes', {rules: [{ required: true, message: $t('ISMScripts.ScriptDes'), whitespace: true}]}]"
          />
        </a-form-item>
          </a-col>
        </a-row>
      </a-form>
      <div
          :style="{
          position: 'absolute',
          right: 0,
          bottom: 0,
          width: '100%',
          borderTop: '1px solid #e9e9e9',
          padding: '10px 16px',
          background: '#fff',
          textAlign: 'right',
          zIndex: 1,
        }"
      >
        <a-button type="danger" :style="{ marginRight: '8px' }" @click="CheckScript">{{$t('ISMScripts.ScriptCheck')}}</a-button>
        <a-button  type="primary" :style="{ marginRight: '8px' }" v-if="!isEdit"  @click="AddScript()">
          {{$t('TaskPlan.TaskAdd')}}
        </a-button>
        <a-button  type="primary" :style="{ marginRight: '8px' }" v-if="isEdit"  @click="EditScript()">
          {{$t('TaskPlan.TaskEdit')}}
        </a-button>

        <a-button  @click="onClose">
          {{$t('device.CancelButton')}}
        </a-button>
      </div>
      </a-drawer>
  </a-card>
</template>

<script>
import codeEditor from '@/components/CodeEditor/index'
import moment from 'moment'
import {ExecSysScript,DisableSysScript} from "@/services/system";
import {formatDate} from "@/utils/common";
import {AddScript, CheckScript, DelScript, EditScript, GetScriptList} from "@/services/ismscripts";
export default {
  name: 'ScriptsList',
  i18n: require('@/i18n/language'),
  data () {
    return {
      isFullscreen:false,
      pagination:{
        pageSize:15,
        showSizeChanger:true
      },
      ScriptType:0,
      isCharge:true,
      CodeContent: "",
      isEdit:false,
      messageShowLoad:false,
      advanced: true,
      refIconLoading: false,
      columns: [
        {
          width: '10%',
          slotName: 'ISMScripts.ScriptName',
          scopedSlots: { customRender: 'ScriptName', title: 'ISMScripts.ScriptName' },
          dataIndex: 'ScriptName'
        },
        {
          slotName: 'ISMScripts.ScriptType',
          width: '10%',
          scopedSlots: { customRender: 'ScriptType', title: 'ISMScripts.ScriptType' },
          dataIndex: 'ScriptType',
        },
        {
          slotName: 'ISMScripts.ScriptDes',
          width: '30%',
          scopedSlots: { customRender: 'Description', title: 'ISMScripts.ScriptDes' },
          dataIndex: 'Description',
        },
        {
          width: '15%',
          slotName: 'dataModel.modelTableOpt',
          scopedSlots: { customRender: 'action',title: 'dataModel.modelTableOpt'}
        }
      ],
      dataSource: [],
      addVisible:false,
      error: '',
      editUuid:"",
      editVisible:false,
      PlanForm: this.$form.createForm(this),
      editForm: this.$form.createForm(this),
      textAreValue:"",
      that:this,
      value: 1
    }
  },
  components: {
    codeEditor,
  },
  authorize: {
    // deleteRecord: 'delete'
  },
  filters: {
    formatDate(time) {
      let date = new Date(time)
      return formatDate(date,'yyyy-MM-dd hh:mm:ss')
    },
    TaskPlanFilters(value,that){
      let TaskPlan = value
      TaskPlan = TaskPlan.split(" ")
      let str = ""

      let week = TaskPlan[5].split(",")
      for (let i=0;i<week.length;i++)
      {
        for(let k=0;k<that.TaskExecPlanList.length;k++)
        {
          if(parseInt(week[i])==that.TaskExecPlanList[k].value){
            str=str+that.$t(that.TaskExecPlanList[k].name)+","
          }
        }
      }
      str = str+TaskPlan[2]+":"+TaskPlan[1]+":"+TaskPlan[0]
      return str
    }
  },
  mounted(){
  },
  activated(){

  },
  created(){
    this.GetScriptList()
  },
  watch: {
    '$route' () {
      this.dataSource=[]
      this.GetScriptList()
    }
  },
  methods: {
    toggleEditorFullscreen () {
      this.isFullscreen = !this.isFullscreen
      const editorEl = document.querySelector('.CodeMirror')
      if(this.isFullscreen) {
        editorEl.classList.add('editor-fullscreen')
      }
      else{
        editorEl.classList.remove('editor-fullscreen')
      }
    },
    chargeScriptType(val){
      this.ScriptType=val
    },
    changeTextarea(val) {
      this.CodeContent = val
    },
    onClose() {
      this.addVisible = false;
    },
    doSysScript(Scriptuuid){
      let Script = []
      let _t =this
      Script.push(Scriptuuid)
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
    DoDisableSysScript(Scriptuuid,disable){
      let _t =this
      if(disable)
      {
        disable=0
      }
      else {
        disable = 1
      }
      let params = {
        uuid:Scriptuuid,
        disable:typeof disable=="undefined"?1:disable
      };
      DisableSysScript(params).then(function (res){
        _t.settingLoading=false
        if(res.data.code==200)
        {
          _t.$message.success("成功")
          _t.settingVisible = false
          _t.GetScriptList()
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
    GetScriptList(){
      let _t = this
      this.dataSource = []
      GetScriptList().then(function (res){
        _t.refIconLoading = false
        if (res.data.code == 200) {
          _t.dataSource = res.data.list
          _t.addVisible = false;
        }
        else if (res.data.code == 2001) {
          _t.$message.error(_t.$t('displayModel.ModelExist'), 3)
        }
        else if (res.data.code == 2003) {
          _t.$message.error(_t.$t('displayModel.AddModelFailed'), 3)
        }
      })
    },
    AddScript(){
      let  _t = this
      this.PlanForm.validateFields((err) => {
        if (!err) {
          const params = {
            ScriptName:_t.PlanForm.getFieldValue('ScriptName'),
            Delay:parseInt(_t.PlanForm.getFieldValue('ScriptDelay')),
            ScriptContent:_t.CodeContent,
            ScriptType:parseInt(_t.PlanForm.getFieldValue('ScriptType')),
            Description:_t.PlanForm.getFieldValue('ScriptDes')
          }
          AddScript(params).then(function (res){
            if (res.data.code == 2002) {
              _t.GetScriptList()
              _t.addVisible = false;
              _t.$message.success(_t.$t('ISMScripts.ScriptAddSuccess'), 3)
            }
            else {
              _t.$message.error(_t.$t('ISMScripts.ScriptAddFailed'), 3)
            }
          })
        }
      })
    },
    EditScript(){
      let  _t = this
      this.PlanForm.validateFields((err) => {
        if (!err) {
          const params = {
            Uuid:_t.EditUUid,
            data:{
              ScriptName:_t.PlanForm.getFieldValue('ScriptName'),
              Delay:parseInt(_t.PlanForm.getFieldValue('ScriptDelay')),
              ScriptContent:_t.CodeContent,
              ScriptType:parseInt(_t.PlanForm.getFieldValue('ScriptType')),
              Description:_t.PlanForm.getFieldValue('ScriptDes')
            }
          }
          EditScript(params).then(function (res){
            if (res.data.code == 200) {
              _t.GetScriptList()
              _t.addVisible = false;
              _t.$message.success(_t.$t('ISMScripts.ScriptEditSuccess'), 3)
            }
            else {
              _t.$message.error(_t.$t('ISMScripts.ScriptEditFailed'), 3)
            }
          })
        }
      })
    },
    GoToEdit(item){
      let _t = this
      _t.isCharge = false
      this.isEdit = true
      this.ScriptType = item.ScriptType
      this.addVisible = true
      _t.EditUUid = item.ScriptUuid
      _t.CodeContent = item.ScriptContent
      setTimeout(function (){
        _t.isCharge = true
        _t.PlanForm.setFieldsValue(
            {
              ScriptName:item.ScriptName,
              ScriptDelay:item.Delay.toString(),
              ScriptType:item.ScriptType,
              ScriptDes:item.Description,
            })
      },200)
    },
    refresh(){
      this.refIconLoading=true
      this.GetScriptList()
    },
    ChargeTaskContent(Value){
      this.TaskContent = Value
    },
    deleteRecord(uuid){
      let _t = this
      const params = {
        ScriptUuid:uuid
      }
      DelScript(params).then(function (res){
        if (res.data.code == 200) {
          _t.GetScriptList()
          _t.addVisible = false;
          _t.$message.success(_t.$t('ISMScripts.ScriptDelSuccess'), 3)
        }
        else {
          _t.$message.error(_t.$t('ISMScripts.ScriptDelFailed'), 3)
        }
      })
    },
    CheckScript(){
      let _t = this
      const params = {
        Script:_t.CodeContent
      }
      CheckScript(params).then(function (res){
        if (res.data.code <0) {
          _t.$message.error(res.data.msg, 3)
        }
        else {
          _t.$message.success(_t.$t('ISMScripts.ScriptCheckResult')+res.data.msg, 3)
        }
      })
    }
  }
}
</script>

<style lang="less" scoped>
::v-deep .fullscreen-btn{
  position: fixed;
  bottom: 20px;
  right: 20px;
  z-index: 9999;
}
::v-deep .editor-fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  height: auto;
  z-index: 9000;
}
::v-deep .search{
  margin-bottom: 54px;
}
::v-deep .ant-form-item {
  margin-bottom: 1px;
}
::v-deep .ant-row .ant-form-item {
  margin-bottom: 1px;
}

::v-deep .fold{
  width: calc(100% - 216px);
  display: inline-block
}
::v-deep .operator{
  margin-bottom: 18px;
}
@media screen and (max-width: 900px) {
  .fold {
    width: 100%;
  }
}
::v-deep .ant-form-item {
  margin-bottom: 10px;
}
::v-deep .ant-table-thead > tr > th
{
  padding: 10px 10px;
  overflow-wrap: break-word;
}
::v-deep .ant-table-tbody > tr > td {
  padding: 7px 5px;
  overflow-wrap: break-word;
}

::v-deep .ant-table-thead>tr>th {
  color: #909399;
  font-weight: 500;
  text-align: left;
  //background: #f8f8f8;
  //border-bottom: 1px solid #e8e8e8;
  transition: background .3s ease;
}
</style>
