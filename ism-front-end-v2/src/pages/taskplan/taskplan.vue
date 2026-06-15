<template>
  <a-card>
    <a-space class="operator">
      <a-button @click="addVisible=true;isEdit=false" type="primary" icon="plus">{{$t('dataModel.newModel')}}</a-button>
      <a-button @click="refresh()"  type="default" icon="sync" :loading="refIconLoading">{{$t("dataModel.refModel")}}</a-button>
    </a-space>

    <a-spin style="padding: 1px;"  :spinning="messageShowLoad" tip="Loading...">
      <a-table rowKey="TaskName" :pagination="pagination" :columns="columns" :data-source="dataSource">
        <template v-for="(item, index) in columns" :slot="item.slotName">
          <span :key="index">{{ $t(item.slotName) }}</span>
        </template>
        <div slot="CronExpression" slot-scope="text,record" >
          <div v-if="record.TaskType == 1">
              {{text}}
          </div>
          <div v-else >
            {{text|TaskPlanFilters(that)}}
          </div>
        </div>
        <div slot="ExecPrvTime" slot-scope="text" >
          {{text|formatDate}}
        </div>
        <div slot="TaskStatus" slot-scope="status,record" >
          <a-switch :checked="status==0?false:true" @change="onChangeStatus($event,record)">
            <a-icon slot="checkedChildren" type="check" />
            <a-icon slot="unCheckedChildren" type="close" />
          </a-switch>
        </div>
        <div slot="TaskContent" slot-scope="text">
          <template
              v-for="(fragment, i) in supportTaskList"
          >
          <span
              v-if="fragment.value == text"
              :key="i"
              class="highlight"
          >
            {{ $t(fragment.name) }}
          </span>
          </template>
        </div>
        <div slot="action" slot-scope="text, record">
          <a @click="GoToEdit(record)" style="color: #13C2C2"><a-icon type="edit" />{{$t('dataModel.modelDetail')}}</a> |
          <a-popconfirm :title="$t('dataModel.deleteConfirm')" @confirm="deleteRecord(record.TaskUuid)">
            <a-icon slot="icon" type="question-circle-o" style="color: red" />
            <a-icon type="delete" theme="twoTone" two-tone-color="#eb2f96"/><a style="color: #eb2f96">{{$t('dataModel.delete')}}</a>
          </a-popconfirm>
        </div>
      </a-table>
    </a-spin>
    <a-drawer
          :title="isEdit?$t('TaskPlan.TaskEdit'):$t('TaskPlan.AddTaskTitle')"
          :width="720"
          :zIndex="900"
          :visible="addVisible"
          :body-style="{ paddingBottom: '80px' }"
          @close="onClose"
      >
      <a-form :form="PlanForm" layout="vertical" >
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item
                :label="$t('TaskPlan.TaskName')"
            >
              <a-input  autocomplete="autocomplete"

                        v-decorator="['TaskName', {rules: [{ required: true, message: $t('TaskPlan.TaskName'), whitespace: true}]}]"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item
                :label="$t('TaskPlan.TaskContent')"
            >
              <a-select  autocomplete="autocomplete" @change="ChargeTaskContent"

                        v-decorator="['TaskContent', {rules: [{ required: true,type:'number', message: $t('TaskPlan.TaskContent'), whitespace: true}]}]"
              >
                <a-select-option v-for="(task,index) in supportTaskList" :key="index" :value="task.value">
                  {{ $t(task.name)}}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item
                :label="$t('TaskPlan.TaskExecPlan')"
            >
          <a-row :gutter="16">
            <a-col :span="12">
              <a-select mode="multiple"
                    :allowClear="true"
                    v-model="TaskExecPlan"
          >
              <a-select-option v-for="(Item,index) in TaskExecPlanList" :key="index" :value="Item.value">
                {{ $t(Item.name) }}
              </a-select-option>
          </a-select>
            </a-col>
            <a-col :span="12" v-if="TaskExecPlan.includes(7)">
              <a-tooltip>
                <template slot="title">
                  {{$t('TaskPlan.CronExplain')}}
                </template>
                  <a-input  v-decorator="['TaskTime', {rules: [{ required: true, message: $t('TaskPlan.TaskExecPlan')}]}]" style="width: 100%" />
              </a-tooltip>
            </a-col>
            <a-col :span="12" v-else>
              <a-time-picker  v-decorator="['TaskTime', {rules: [{ required: true, message: $t('TaskPlan.TaskExecPlan')}]}]" style="width: 100%" />
            </a-col>
          </a-row>
        </a-form-item>

        <a-row :gutter="16" v-if="TaskContent==3||TaskContent==4">
          <a-col :span="12">
            <a-form-item
                :label="$t('TaskPlan.HistoryKeepDay')"
            >
              <a-input  autocomplete="autocomplete"

                        v-decorator="['HistoryKeepDay', {rules: [{ required: true, message: $t('TaskPlan.HistoryKeepDay'), whitespace: true}]}]"
              />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16" v-if="TaskContent==2">
          <a-button  type="link" icon="plus" @click="AddDeviceList">{{$t('TaskPlan.TaskAddSetData')}}</a-button>
          <a-col :span="24" v-for="(item,index) in SetDeviceList" :key="index">
            <a-col :span="12">
              <a-form-item
                  :label="$t('dataModel.CustomData.DataFrom')"
              >
                <a-input    v-model="item.dataName">
                  <a-tooltip placement="top" slot="addonAfter">
                    <template slot="title">
                      <span>{{$t('displayConfig.Properties.SelectValue')}}</span>
                    </template>
                    <icon-font @click="ShowDeviceDataModel();SelectDeviceIndex=index" type="icon-xuanzeshuju"  />
                  </a-tooltip>
                </a-input>
              </a-form-item>
            </a-col>

            <a-col :span="10">
              <a-form-item
                  :label="$t('TaskPlan.SetValue')"
              >
                <a-input  autocomplete="autocomplete" v-model="item.Value"
                />
              </a-form-item>
            </a-col>

            <a-col :span="2">
              <a-form-item
                  label="操作"
              >
              <icon-font @click="DelDeviceList(index)" type="icon-shanchu"  style="margin-top: 5px;font-size: 20px"/>
              </a-form-item>
            </a-col>
          </a-col>
        </a-row>
        <a-row :gutter="16" v-if="TaskContent==5">
          <a-col :span="12">
            <a-form-item
                :label="$t('TaskPlan.ReportTemplete')"
            >
            <a-select mode="multiple"
                      :allowClear="true"
                      v-model="reportTemplete"
            >
              <a-select-option v-for="(Item,index) in reportTempleteList" :key="index" :value="Item.Uuid">
                {{ $t(Item.Name) }}
              </a-select-option>
            </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16" v-if="TaskContent==6">
          <a-col :span="12">
            <a-form-item
                :label="$t('TaskPlan.scriptList')"
            >
              <a-select mode="multiple"
                        :allowClear="true"
                        v-model="scriptListSave"
              >
                <a-select-option v-for="(Item,index) in scriptList" :key="index" :value="Item.ScriptUuid">
                  {{ $t(Item.ScriptName) }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16" v-if="TaskContent==7">
          <a-col :span="12">
            <a-form-item
                :label="$t('TaskPlan.CheckVideoSize')"
            >
              <a-input  autocomplete="autocomplete"

                        v-decorator="['CheckVideoSize', {rules: [{ required: true, message: $t('TaskPlan.CheckVideoSize'), whitespace: true}]}]"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item
                :label="$t('TaskPlan.KeepVideoDays')"
            >
              <a-input  autocomplete="autocomplete"

                        v-decorator="['KeepVideoDays', {rules: [{ required: true, message: $t('TaskPlan.KeepVideoDays'), whitespace: true}]}]"
              />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16" v-if="TaskContent==9">
          <a-col :span="12">
            <a-form-item
                :label="$t('TaskPlan.SQLReportTemplate')"
            >
              <a-select mode="multiple"
                        :allowClear="true"
                        v-model="SQLReportTemplate"
              >
                <a-select-option v-for="(Item,index) in SqlTemplateSource" :key="index" :value="Item.Uuid">
                  {{ $t(Item.Name) }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="24">
          <a-form-item
            :label="$t('dataModel.modelDec')"
        >
          <a-textarea autocomplete="autocomplete"

                      v-decorator="['Description', {rules: [{ required: true, message: $t('dataModel.modelDec'), whitespace: true}]}]"
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
        <a-button  type="primary" :style="{ marginRight: '8px' }" v-if="!isEdit"  @click="AddTaskPlan()">
          {{$t('TaskPlan.TaskAdd')}}
        </a-button>
        <a-button  type="primary" :style="{ marginRight: '8px' }" v-if="isEdit"  @click="EditTaskPlan()">
          {{$t('TaskPlan.TaskEdit')}}
        </a-button>

        <a-button  @click="onClose">
          {{$t('device.CancelButton')}}
        </a-button>
      </div>
      </a-drawer>
    <device-data-model @onSelectDataModel="onSelectData" ref="deviceDataModel"></device-data-model>
  </a-card>
</template>

<script>
import deviceDataModel from "@/components/deviceDataModel/deviceDataModel";
import {modbusModelGroupAdd} from "@/services/modbusModel";
import moment from 'moment'
import {GetTaskPlanList,AddTaskPlan,EditTaskPlan,DelTaskPlan} from "@/services/taskplan";
import {formatDate} from "@/utils/common";
import {GetReportTempletes} from "@/services/reportTemplete";
import {GetScriptList} from "@/services/ismscripts";
import {GetSQLReportTempletes} from "@/services/SqlReportTemplete";
export default {
  name: 'TaskPlan',
  i18n: require('@/i18n/language'),
  data () {
    return {
      SqlTemplateSource:[],
      scriptListSave:[],
      pagination:{
        pageSize:15,
        showSizeChanger:true
      },

      SelectDeviceIndex:0,
      TaskContent:0,
      isEdit:false,
      TaskExecPlan:[],
      TaskTime:"",
      reportTemplete:[],
      SQLReportTemplate:[],
      TaskExecPlanList:[
        {
            value:1,
            name:"TaskPlan.TaskPlanMonday"
        },
        {
          value:2,
          name:"TaskPlan.TaskPlanTuesday"
        },
        {
          value:3,
          name:"TaskPlan.TaskPlanWednesday"
        },
        {
          value:4,
          name:"TaskPlan.TaskPlanThursday"
        },
        {
          value:5,
          name:"TaskPlan.TaskPlanFriday"
        },
        {
          value:6,
          name:"TaskPlan.TaskPlanSaturday"
        },
        {
          value:0,
          name:"TaskPlan.TaskPlanSunday"
        },
        {
          value:7,
          name:"TaskPlan.TaskPlanDiy"
        },
      ],
      DataRecordType:0,
      SetDeviceList:[
        {
          "deviceSN":"",
          "dataName":"",
          "DataSN":"",
          "Value":"",
        }
      ],
      supportTaskList:[
        {
          value:1,
          name:'TaskPlan.TaskBackDb'
        },
        {
          value:2,
          name:'TaskPlan.TaskSetData'
        },
        {
          value:3,
          name:'TaskPlan.TaskDelHistoryData'
        },
        {
          value:4,
          name:'TaskPlan.TaskDelHistoryAlarm'
        },
        {
          value:5,
          name:'TaskPlan.TimelyExport'
        },
        {
          value:6,
          name:'TaskPlan.TimelyScripts'
        },
        {
          value:7,
          name:'TaskPlan.TimelyCheckVideo'
        },
        {
          value:8,
          name:'TaskPlan.TimelySyncNTP'
        },
        {
          value:9,
          name:'TaskPlan.TimelyExportSQLReport'
        }
      ],
      alarmStatus:2,
      recordStatus:0,
      messageShowLoad:false,
      advanced: true,
      refIconLoading: false,
      columns: [
        {
          width: '10%',
          slotName: 'TaskPlan.TaskName',
          scopedSlots: { customRender: 'TaskName', title: 'TaskPlan.TaskName' },
          dataIndex: 'TaskName'
        },
        {
          width: '10%',
          slotName: 'TaskPlan.ExecTime',
          scopedSlots: { customRender: 'CronExpression', title: 'TaskPlan.ExecTime' },
          dataIndex: 'CronExpression',
        },
        {
          slotName: 'TaskPlan.TaskContent',
          width: '10%',
          scopedSlots: { customRender: 'TaskContent', title: 'TaskPlan.TaskContent' },
          dataIndex: 'TaskContent',
        },
        {
          width: '5%',
          slotName: 'TaskPlan.TaskStatus',
          scopedSlots: { customRender: 'TaskStatus', title: 'TaskPlan.TaskStatus' },
          dataIndex: 'Status',
        },
        {
          width: '10%',
          slotName: 'TaskPlan.ExecPrvTime',
          scopedSlots: { customRender: 'ExecPrvTime', title: 'TaskPlan.ExecPrvTime' },
          dataIndex: 'PrevTime',
        },
        {
          width: '5%',
          slotName: 'TaskPlan.ExecTimes',
          scopedSlots: { customRender: 'ExecTimes', title: 'TaskPlan.ExecTimes' },
          dataIndex: 'ExecuteTimes',
        },
        {
          width: '10%',
          slotName: 'dataModel.modelTableOpt',
          scopedSlots: { customRender: 'action',title: 'dataModel.modelTableOpt'}
        }
      ],
      DataTypeList:[
        {
          name:this.$t('dataModel.static.DataTypeBool'),
          value:2
        },
        {
          name:this.$t('dataModel.static.DataTypeBit'),
          value:3
        },
        {
          name:this.$t('dataModel.static.DataTypeInt'),
          value:1
        },
        {
          name:this.$t('dataModel.static.DataTypeDouble'),
          value:4
        }
      ],
      supportDeviceList:[],
      dataSource: [],
      reportTempleteList:[],
      scriptList:[],
      selectedRows: [],
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
    deviceDataModel
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
      let str = ""
     try {
       TaskPlan = TaskPlan.split(" ")

       let week = TaskPlan[5].split(",")
       for (let i = 0; i < week.length; i++) {
         for (let k = 0; k < that.TaskExecPlanList.length; k++) {
           if (parseInt(week[i]) == that.TaskExecPlanList[k].value) {
             str = str + that.$t(that.TaskExecPlanList[k].name) + ","
           }
         }
       }
       str = str + TaskPlan[2] + ":" + TaskPlan[1] + ":" + TaskPlan[0]
     }catch (e) {
       str  = value
     }
      return str
    }
  },
  mounted(){
    this.GetReportTempletes()
    this.GetSQLReportTemplates()
    this.getTaskList()
    this.GetScriptList()
  },
  activated(){

  },
  created(){

  },
  watch: {
    '$route' () {
      this.dataSource=[]
      this.getTaskList()
      this.GetScriptList()
    }
  },
  methods: {
    GetScriptList(){
      let _t = this
      this.scriptList = []
      GetScriptList().then(function (res){
        _t.refIconLoading = false
        if (res.data.code == 200) {
          _t.scriptList = res.data.list
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
    GetReportTempletes(){
      let _t = this
      _t.loading = true
      this.reportTempleteList=[]
      GetReportTempletes().then(function (res){
        _t.loading = false
        if (res.data.code == 0) {
          if(res.data.list==null)
          {
            _t.reportTempleteList=[]
          }
          else
          {
            _t.reportTempleteList = res.data.list
          }
        }
      }).finally(function (error) {
        _t.loading = false
      })
    },
    onClose() {
      this.addVisible = false;
    },
    DelDeviceList(index){
      this.SetDeviceList.splice(index,1)
    },
    AddDeviceList(){
      this.SetDeviceList.push({
        "deviceSN":"",
        "dataName":"",
        "DataSN":"",
        "Value":"",
      })
    },
    refresh(){
      this.refIconLoading=true
      this.getTaskList()
    },
    ChargeTaskContent(Value){
      this.TaskContent = Value
    },
    onChangeStatus(value,item){
      let _t = this
      item.Status = value?1:0
      const params = {
        uuid:item.TaskUuid,
        data: item
      }
      EditTaskPlan(params).then(function (res){
        if (res.data.code == 200) {
          _t.getTaskList()
          _t.$message.success(_t.$t('dataModel.static.EditSuccess'), 3)
        }
        else {
          _t.$message.error(_t.$t('dataModel.static.EditFailed'), 3)
        }
      })
    },
    GoToEdit(item){
      let _t = this
      this.isEdit = true
      this.addVisible = true
      this.taskStatus = item.Status
      _t.EditUUid = item.TaskUuid
      _t.TaskExecPlan=[]
      _t.SetDeviceList = []
      if(item.SetDeviceList=="")
      {
        _t.SetDeviceList=[]
      }
      else {
        try{
          _t.SetDeviceList = JSON.parse(item.SetDeviceList)
        }catch (e) {
          _t.SetDeviceList=[]
        }
      }
      let TaskTime=""
      let TaskType = item.TaskType
      if(TaskType == 1)
      {
        _t.TaskExecPlan.push(7)
        TaskTime = item.CronExpression
      }
      else if(TaskType == 2) {
        let CronExpression = item.CronExpression.split(' ')
        let TaskExecPlan = CronExpression[5].split(',')
        for (let i = 0; i < TaskExecPlan.length; i++) {
          _t.TaskExecPlan.push(parseInt(TaskExecPlan[i]))
        }
        TaskTime = moment(CronExpression[2]+":"+CronExpression[1]+":"+CronExpression[0],'HH:mm:ss')
      }
      _t.TaskContent = item.TaskContent
      if(_t.TaskContent==5)
      {
        try{
          _t.reportTemplete = JSON.parse(item.ReportTempleteList)
        }catch (e) {
          _t.reportTemplete=[]
        }
      }
      else if(_t.TaskContent==6)
      {
        try{
          _t.scriptListSave = JSON.parse(item.ScriptList)
        }catch (e) {
          _t.scriptListSave=[]
        }
      }
      else if(_t.TaskContent==9)
      {
        try{
          console.log(item.SQLReportTempleteList)
          _t.SQLReportTemplate = JSON.parse(item.SQLReportTempleteList)
        }catch (e) {
          _t.SQLReportTemplate=[]
        }
      }

      setTimeout(function (){
        _t.PlanForm.setFieldsValue(
            {
              TaskName:item.TaskName,
              TaskTime:TaskTime,
              TaskContent:item.TaskContent,
              Description:item.Description,
            })
        if(_t.TaskContent==3||_t.TaskContent==4)
        {
          _t.PlanForm.setFieldsValue(
              {
                HistoryKeepDay:item.KeepHistoryDay.toString(),
              })
        }
        if(_t.TaskContent==7)
        {
          _t.PlanForm.setFieldsValue(
              {
                CheckVideoSize:item.MaxDirSize.toString(),
                KeepVideoDays:item.MinFileAge.toString(),
              })
        }
      },200)
    },
    EditTaskPlan(){
      let _t = this
      this.PlanForm.validateFields((err) => {
        if (!err) {
          const params = {
            uuid:_t.EditUUid,
            data: {
              TaskName: _t.PlanForm.getFieldValue('TaskName'),
              TaskContent: _t.PlanForm.getFieldValue('TaskContent'),
              Description: _t.PlanForm.getFieldValue('Description'),
              Status:_t.taskStatus,
            }
          }
          if(_t.TaskExecPlan.includes(7))
          {
            params.data.CronExpression = _t.PlanForm.getFieldValue('TaskTime')
            params.data.TaskType=1
          }
          else {
            let TaskTime = moment(_t.PlanForm.getFieldValue('TaskTime')).format('HH:mm:ss');

            TaskTime = TaskTime.split(':')
            if (TaskTime.length != 3) {
              return
            }
            let CronExpression = TaskTime[2] + ' ' + TaskTime[1] + ' ' + TaskTime[0] + ' * *'
            if (_t.TaskExecPlan.length == 0) {
              CronExpression = CronExpression + ' ?'
            } else {
              CronExpression = CronExpression + ' ' + _t.TaskExecPlan.join(',')
            }
            params.data.CronExpression = CronExpression
            params.TaskType=2
          }

          if((_t.TaskContent==2)&&(_t.SetDeviceList.length==0))
          {
            _t.$message.error(_t.$t('TaskPlan.SetValueEmpty'), 3)
            return
          }
          else
          {
            params.data.SetDeviceList = JSON.stringify(_t.SetDeviceList)
          }
          if((_t.TaskContent==3)||(_t.TaskContent==4))
          {
            params.data.KeepHistoryDay = parseInt(_t.PlanForm.getFieldValue('HistoryKeepDay'))
            if(params.data.KeepHistoryDay<=0)
            {
              _t.$message.error(_t.$t('TaskPlan.KeepDayError'), 3)
              return
            }
          }
          if((_t.TaskContent==5)) {
            if (_t.reportTemplete.length <= 0) {
              _t.$message.error(_t.$t('TaskPlan.ReportTempleteEmpty'), 3)
              return
            }
            params.data.ReportTempleteList = JSON.stringify(_t.reportTemplete)
          }
          if((_t.TaskContent==6)) {
            if (_t.scriptListSave.length <= 0) {
              _t.$message.error(_t.$t('TaskPlan.ReportTempleteEmpty'), 3)
              return
            }
            params.data.ScriptList = JSON.stringify(_t.scriptListSave)
          }
          if((_t.TaskContent==7)) {
            params.data.MaxDirSize = parseInt(_t.PlanForm.getFieldValue('CheckVideoSize'))
            params.data.MinFileAge = parseInt(_t.PlanForm.getFieldValue('KeepVideoDays'))
          }
          if((_t.TaskContent==9)) {
            if (_t.SQLReportTemplate.length <= 0) {
              _t.$message.error(_t.$t('TaskPlan.ReportTempleteEmpty'), 3)
              return
            }
            params.data.SQLReportTempleteList = JSON.stringify(_t.SQLReportTemplate)
          }
          EditTaskPlan(params).then(function (res){
            if (res.data.code == 200) {
              _t.getTaskList()
              _t.addVisible = false;
              _t.isEdit = false
              _t.$message.success(_t.$t('dataModel.static.EditSuccess'), 3)
            }
            else {
              _t.$message.error(_t.$t('dataModel.static.EditFailed'), 3)
            }
          })
        }
      })
    },
    getTaskList(){
      let _t = this
      GetTaskPlanList().then(function (res){
        _t.refIconLoading = false
        if (res.data.code == 200) {
          _t.dataSource = res.data.list
        }
        else if (res.data.code == 2004)
        {
          _t.$message.error(_t.$t('displayModel.modelHavedBind'), 3)
        }
        else
        {
          _t.$message.error(_t.$t('displayModel.DelModelFailed'), 3)
        }
      })
    },
    AddTaskPlan(){
      let _t = this
      this.isEdit = false
      this.PlanForm.validateFields((err) => {
        if (!err) {
          const params = {
            TaskName:_t.PlanForm.getFieldValue('TaskName'),
            TaskContent:_t.PlanForm.getFieldValue('TaskContent'),
            Description:_t.PlanForm.getFieldValue('Description')
          }
          if(_t.TaskExecPlan.includes(7))
          {
            params.CronExpression = _t.PlanForm.getFieldValue('TaskTime')
            params.TaskType=1
          }
          else
          {
            let TaskTime = moment(_t.PlanForm.getFieldValue('TaskTime')).format('HH:mm:ss');

            TaskTime = TaskTime.split(':')
            if (TaskTime.length != 3) {
              return
            }
            let CronExpression = TaskTime[2] + ' ' + TaskTime[1] + ' ' + TaskTime[0] + ' * *'
            if (_t.TaskExecPlan.length == 0) {
              CronExpression = CronExpression + ' ?'
            } else {
              CronExpression = CronExpression + ' ' + _t.TaskExecPlan.join(',')
            }
            params.CronExpression = CronExpression
            params.TaskType=2
          }
          if((_t.TaskContent==2)&&(_t.SetDeviceList.length==0))
          {
            _t.$message.error(_t.$t('TaskPlan.SetValueEmpty'), 3)
              return
          }
          else
          {
            params.SetDeviceList = JSON.stringify(_t.SetDeviceList)
          }
          if((_t.TaskContent==3)||(_t.TaskContent==4))
          {
            params.KeepHistoryDay = parseInt(_t.PlanForm.getFieldValue('HistoryKeepDay'))
            if(params.KeepHistoryDay<=0)
            {
              _t.$message.error(_t.$t('TaskPlan.KeepDayError'), 3)
              return
            }
          }
          else if((_t.TaskContent==5))
          {
            if(_t.reportTemplete.length<=0)
            {
              _t.$message.error(_t.$t('TaskPlan.ReportTempleteEmpty'), 3)
              return
            }
            params.ReportTempleteList = JSON.stringify(_t.reportTemplete)
          }
          else if((_t.TaskContent==6))
          {
            if(_t.scriptListSave.length<=0)
            {
              _t.$message.error(_t.$t('TaskPlan.ScriptEmpty'), 3)
              return
            }
            params.ScriptList = JSON.stringify(_t.scriptListSave)
          }
          else if((_t.TaskContent==7))
          {
            params.MaxDirSize = parseInt(_t.PlanForm.getFieldValue('CheckVideoSize'))
            params.MinFileAge = parseInt(_t.PlanForm.getFieldValue('KeepVideoDays'))
          }
          else if((_t.TaskContent==9))
          {
            if(_t.SQLReportTemplate.length<=0)
            {
              _t.$message.error(_t.$t('TaskPlan.ReportTempleteEmpty'), 3)
              return
            }
            params.SQLReportTempleteList = JSON.stringify(_t.SQLReportTemplate)
          }
          AddTaskPlan(params).then(function (res){
              if (res.data.code == 2002) {
                _t.getTaskList()
                _t.addVisible = false;
                _t.$message.success(_t.$t('displayModel.AddModelSuccess'), 3)
              }  else if (res.data.code == 2001) {
                _t.$message.error(_t.$t('TaskPlan.AddExist'), 3)
              }
              else if (res.data.code == -3) {
                _t.$message.error(_t.$t('TaskPlan.AddExFailed'), 3)
              }else{
                _t.$message.error(_t.$t('TaskPlan.AddFailed'), 3)
              }
          })
        }
      })
    },
    ShowDeviceDataModel(index,type){
      this.$refs.deviceDataModel.showDataModal()
    },
    deleteRecord(uuid){
      let _t = this
      let params ={
        TaskUuid:uuid
      }
      DelTaskPlan(params).then(function (res){
        _t.refIconLoading = false
        if (res.data.code == 200) {
          _t.getTaskList()
          _t.$message.success(_t.$t('displayModel.DelModelSuccess'), 3)
        }
        else
        {
          _t.$message.error(_t.$t('displayModel.DelModelFailed'), 3)
        }
      })
    },
    onSelectData(selectData) {
      if(selectData.DeviceType!=-1)
      {
        this.$message.error(this.$t('TaskPlan.SelectError'), 3)
        return
      }
      this.SetDeviceList[this.SelectDeviceIndex].Value=""
      this.SetDeviceList[this.SelectDeviceIndex].deviceSN=selectData.DeviceSN
      this.SetDeviceList[this.SelectDeviceIndex].dataName=selectData.DeviceName+"\\"+selectData.name
      this.SetDeviceList[this.SelectDeviceIndex].DataSN=selectData.uuid
    },
    GetSQLReportTemplates(){
      let _t = this

      _t.loading = true
      this.SqlTemplateSource=[]
      GetSQLReportTempletes().then(function (res){
        _t.loading = false
        if (res.data.code == 0) {
          if(res.data.list==null)
          {
            _t.SqlTemplateSource=[]
          }
          else
          {
            _t.SqlTemplateSource = res.data.list
          }
        }
      }).finally(function (error) {
        _t.loading = false
      })
    },
  }
}
</script>

<style lang="less" scoped>
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
