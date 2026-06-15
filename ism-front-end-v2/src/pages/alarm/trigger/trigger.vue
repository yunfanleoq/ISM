<template>
  <a-card>
    <a-spin style="z-index:9999;margin-right: -300px" size="large" :spinning="ShowRegisterLoading" tip="Loading...">
      <a-drawer
        :title="$t('alarm.trigger.addTriggerTitle')"
        :width="720"
        :visible="visible"
        :body-style="{ paddingBottom: '80px' }"
        :after-visible-change="afterVisibleChange"
        @close="onClose"
    >
      <a-form :form="form" layout="vertical" >
        <a-row :gutter="16">
          <a-col :span="24">
            <a-form-item style="margin-bottom:5px" :label="$t('alarm.trigger.TriggerAddName')">
              <a-input  autocomplete="autocomplete"
                        v-decorator="['TriggerAddName', {rules: [{ required: true, message: $t('alarm.trigger.TriggerAddName'), whitespace: true}]}]"
              />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="8">
            <a-form-item style="margin-bottom:5px" :label="$t('alarm.trigger.DeviceType')">
              <a-select  @change="changeDeviceType" :disabled="isEdit"
                        v-decorator="[
                  'DeviceType',
                  {
                    rules: [{ required: true, message: $t('alarm.trigger.DeviceType') }],
                  },
                ]"
              >
                <a-select-option  v-for="(device,index) in supportDeviceList" :key="index" :value=device.type>
                  {{ $t(device.name)  }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item style="margin-bottom:5px" :label="$t('alarm.trigger.DeviceModel')">
              <a-select  @change="changeDeviceModel" allowClear :disabled="isEdit"
                        v-decorator="[
                  'DeviceModel',
                  {
                    rules: [{ required: true, message: $t('alarm.trigger.DeviceModel') }],
                  },
                ]"
              >
                <a-select-option v-for="(model,index) in supportDeviceModelList" :key="index" :value="model.uuid">
                  {{ model.modelName }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item style="margin-bottom:5px" :label="$t('alarm.trigger.DeviceModelData')">
              <a-select
                  show-search
                  option-filter-prop="children"
                  :filter-option="filterOption"
                  v-decorator="[
                  'DeviceModelData',
                  {
                    rules: [{ required: true, message: $t('alarm.trigger.DeviceModel') }],
                  },
                ]"
              >
                <a-select-option v-for="(model,index) in dataModelList" :key="index" :value="model.uuid">
                  {{ model.name }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="16">
          <a-col :span="8">
            <a-form-item style="margin-bottom:5px" :label="$t('alarm.trigger.TriggerCondition')">
              <a-select
                  show-search
                  @change="changeCondition"
                  option-filter-prop="children"
                  v-decorator="[
                  'TriggerCondition',
                  {
                    rules: [{ required: true, message: $t('alarm.trigger.TriggerCondition') }],
                  },
                ]"
              >
                <a-select-option value=">">
                  {{$t('alarm.trigger.TriggerConditionGreater')}}
                </a-select-option>
                <a-select-option value=">=">
                  {{$t('alarm.trigger.TriggerConditionGreaterEq')}}
                </a-select-option>
                <a-select-option value="<">
                  {{$t('alarm.trigger.TriggerConditionLess')}}
                </a-select-option>
                <a-select-option value="<=">
                  {{$t('alarm.trigger.TriggerConditionLessEq')}}
                </a-select-option>
                <a-select-option value="=">
                  {{$t('alarm.trigger.TriggerConditionEqual')}}
                </a-select-option>
                <a-select-option value="!=">
                  {{$t('alarm.trigger.TriggerConditionNot')}}
                </a-select-option>
                <a-select-option value="&&">
                  {{$t('alarm.trigger.TriggerConditionRange')}}
                </a-select-option>
                <a-select-option value="||">
                  {{$t('alarm.trigger.TriggerConditionOutRange')}}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item style="margin-bottom:5px" :label="$t('alarm.trigger.TriggerConditionXValue')">
              <a-input  autocomplete="autocomplete"
                        v-decorator="['XValue', {rules: [{ required: true, message: $t('alarm.trigger.TriggerConditionXValue'), whitespace: true}]}]"
              />
            </a-form-item>
          </a-col>
          <a-col :span="8" v-if="(conditionExpress=='&&')||(conditionExpress=='||')">
            <a-form-item style="margin-bottom:5px" :label="$t('alarm.trigger.TriggerConditionYValue')">
              <a-input  autocomplete="autocomplete"
                        v-decorator="['YValue', {rules: [{ required: true, message: $t('alarm.trigger.TriggerConditionYValue'), whitespace: true}]}]"
              />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="24">
            <a-form-item style="margin-bottom:5px">
              <span slot="label">
                {{$t('alarm.trigger.TriggerKeepTime')}}
                <a-tooltip :title="$t('alarm.trigger.TriggerKeepTimeTips')">
                  <a-icon type="question-circle-o" />
                </a-tooltip>
              </span>
              <a-input-number :min="500" style="width: 100%"
                  v-decorator="[
                  'TriggerKeepTime',
                  {
                    rules: [{ required: true, message: $t('alarm.trigger.TriggerKeepTime') }],
                  },
                ]"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item style="margin-bottom:20px" :label="$t('alarm.trigger.TriggerAlarmSwitch')">
              <a-switch @change="onTriggerAlarmSwitchChange" v-decorator="['TriggerAlarmSwitch', {rules: [{required: true, message: $t('alarm.trigger.TriggerAlarmSwitch')}],initialValue: false,valuePropName: 'checked'}] ">
                <a-icon slot="checkedChildren" type="check" />
                <a-icon slot="unCheckedChildren" type="close" />
              </a-switch>
            </a-form-item>
          </a-col>
        </a-row>
        <transition name="plus-icon">
          <a-row :gutter="16" v-if="TriggerAlarmSwitchStatus">
            <section class="code-box">
              <section class="code-box-meta markdown">
                <h4>{{$t('alarm.trigger.TriggerAlarmSwitch')}}</h4>
                <div class="code-box-actions">
                  <a-col :span="24" >
                    <a-form-item :label="$t('dataModel.AlarmLevel')">
                      <a-select   autocomplete="autocomplete"  v-decorator="[
                  'AlarmLevel',
                  {
                    rules: [{ required: true, message: $t('dataModel.AlarmLevel') }],
                  },
                ]">
                        <a-select-option value='0'>{{$t('dataModel.alarm.Tips')}}</a-select-option>
                        <a-select-option value='1'>{{$t('dataModel.alarm.Minor')}}</a-select-option>
                        <a-select-option value='2'>{{$t('dataModel.alarm.Importance')}}</a-select-option>
                        <a-select-option value='3'>{{$t('dataModel.alarm.Urgency')}}</a-select-option>
                        <a-select-option value='4'>{{$t('dataModel.alarm.Deadly')}}</a-select-option>
                      </a-select>
                    </a-form-item>
                  </a-col>
                  <a-col :span="12" >
                    <a-form-item style="margin-bottom: 0px" :label="$t('alarm.trigger.TriggerAlarmShowText')">
                      <a-textarea
                          v-decorator="[
                  'TriggerAlarmShowText',
                  {
                    rules: [{ required: true, message: $t('alarm.trigger.TriggerAlarmShowText') }],
                  },
                ]"
                          :rows="4"
                      />
                    </a-form-item>
                  </a-col>
                  <a-col :span="12">
                    <a-form-item style="margin-bottom: 0px" :label="$t('alarm.trigger.TriggerAlarmHideText')">
                      <a-textarea
                          v-decorator="[
                  'TriggerAlarmHideText',
                  {
                    rules: [{ required: false, message: $t('alarm.trigger.TriggerAlarmHideText') }],
                  },
                ]"
                          :rows="4"
                      />
                    </a-form-item>
                  </a-col>
                </div>
              </section>
            </section>
          </a-row>
        </transition>

        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item style="margin-bottom:20px" :label="$t('alarm.trigger.TriggerLinkageSwitch')">
              <a-switch @change="onTriggerLinkageSwitch" v-decorator="['TriggerLinkageSwitch', {rules: [{required: true, message: $t('alarm.trigger.TriggerLinkageSwitch')}],initialValue: false,valuePropName: 'checked'}] ">
                <a-icon slot="checkedChildren" type="check" />
                <a-icon slot="unCheckedChildren" type="close" />
              </a-switch>
            </a-form-item>
          </a-col>
        </a-row>
        <transition name="plus-icon">
          <a-row :gutter="16" v-if="TriggerLinkageSwitchStatus">
          <section  class="code-box">
            <section class="code-box-meta markdown">
              <h4>{{$t('alarm.trigger.TriggerLinkageSwitch')}}</h4>
              <div class="code-box-actions">
                <a-col :span="24">
                  <a-form-item style="margin-bottom: 0px" :label="$t('alarm.trigger.DeviceLinkModelData')">
                    <a-select
                        show-search
                        option-filter-prop="children"
                        :filter-option="filterOption"
                        v-decorator="[
                  'DeviceLinkModelData',
                  {
                    rules: [{ required: true, message: $t('alarm.trigger.DeviceLinkModelData') }],
                  },
                ]"
                    >
                      <a-select-option v-for="(model,index) in dataModelList" :key="index" :value="model.uuid">
                        {{ model.name }}
                      </a-select-option>
                    </a-select>
                  </a-form-item>
                </a-col>
                <a-col :span="24">
                  <a-form-item style="margin-bottom: 0px" :label="$t('alarm.trigger.TriggerLinkageAlarmValue')">
                    <a-input
                               v-decorator="[
                  'TriggerLinkageAlarmValue',
                  {
                    rules: [{ required: true, message: $t('alarm.trigger.TriggerLinkageAlarmValue') }],
                  },
                ]"
                    >
                    </a-input>
                  </a-form-item>
                </a-col>
                <a-col :span="24">
                  <a-form-item style="margin-bottom: 0px" :label="$t('alarm.trigger.TriggerLinkageAlarmClearValue')">
                    <a-input
                        v-decorator="[
                  'TriggerLinkageAlarmClearValue',
                  {
                    rules: [{ required: false, message: $t('alarm.trigger.TriggerLinkageAlarmClearValue') }],
                  },
                ]"
                    >
                    </a-input>
                  </a-form-item>
                </a-col>
              </div>
            </section>
          </section>
        </a-row>
        </transition>
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
        <a-button  type="primary" v-if="!isEdit" :style="{ marginRight: '8px' }" @click="AddTriggerBtn">
          {{$t('alarm.trigger.AddTriggerBtn')}}
        </a-button>
        <a-button  type="primary" v-else :style="{ marginRight: '8px' }" @click="EditTriggerBtn">
          {{$t('alarm.trigger.EditTriggerBtn')}}
        </a-button>

        <a-button  @click="onClose">
          {{$t('device.CancelButton')}}
        </a-button>
      </div>
    </a-drawer>
    </a-spin>

      <a-space class="operator">
        <a-button @click="ShowAddTrigger('add')" type="primary" icon="plus">{{$t('dataModel.newModel')}}</a-button>
        <a-button @click="GetTriggerList"  type="default" icon="sync" :loading="messageShowLoad">{{$t("dataModel.refModel")}}</a-button>
      </a-space>
      <div>
      <a-spin style="padding: 1px;"  :spinning="messageShowLoad" tip="Loading...">
        <a-table :pagination="pagination" :columns="columns" :data-source="dataSource" rowKey="TriggerName" >
        <template v-for="(item, index) in columns" :slot="item.slotName">
          <span :key="index">{{ $t(item.slotName) }}</span>
        </template>
        <div slot="TriggerType" slot-scope="TriggerType">
          <span v-if="TriggerType==1">
              {{$t("alarm.trigger.TriggerLinkType")}}
          </span>
          <span v-else-if="TriggerType==2">
              {{$t("alarm.trigger.TriggerAlarmType")}}
          </span>
          <span v-else-if="TriggerType==3">
              {{$t("alarm.trigger.TriggerLinkAndAlarmType")}}
          </span>
        </div>

          <div slot="UpdatedAt" slot-scope="UpdatedAt">
            <span>
              {{UpdatedAt|formatDate}}
            </span>
          </div>
        <div slot="action" slot-scope="text, record">
          <a type="link"   @click="gotoEditTrigger(record)" style="cursor: pointer;color: #13C2C2"><a-icon type="edit" /><span style="margin-left: 2px;">{{$t('alarm.trigger.Detail')}}</span></a> |
          <a-popconfirm :title="$t('dataModel.deleteConfirm')" @confirm="deleteRecord(record)">
            <a-icon slot="icon" type="question-circle-o" style="color: red" />
            <a-icon type="delete" theme="twoTone" two-tone-color="#eb2f96"/><a style="color: #eb2f96">{{$t('dataModel.delete')}}</a>
          </a-popconfirm>
        </div>
      </a-table>
      </a-spin>
      </div>
  </a-card>
</template>

<script>
import {getSupportDeviceList} from "../../../services/device";
import {getDatasByUuid,snmpModelList} from "../../../services/snmpmodel";
import {AlarmTriggerAdd,GetAlarmTriggerList,AlarmTriggerDel,AlarmTriggerEdit} from "../../../services/alarm";
import {formatDate} from '@/utils/common';
export default {
  name: 'trigger',
  i18n: require('../../../i18n/language'),
  data () {
    return {
      pagination:{
        pageSize:15,
        showSizeChanger:true
      },
      editTriggerID:0,
      dataModelList:[],
      dataLinkModelList:[],
      ShowRegisterLoading:false,
      supportDeviceList:[],
      supportLinkDeviceList:[],
      supportDeviceModelList:[],
      supportLinkDeviceModelList:[],
      visible:false,
      isEdit:false,
      DeviceType:0,
      form: this.$form.createForm(this),
      messageShowLoad:false,
      advanced: true,
      TriggerAlarmSwitchStatus:false,
      TriggerLinkageSwitchStatus:false,
      refIconLoading: false,
      columns: [
        {
          width: '20%',
          slotName: 'alarm.trigger.TriggerAddName',
          scopedSlots: { customRender: 'TriggerName', title: 'alarm.trigger.TriggerAddName' },
          dataIndex: 'TriggerName',
        },
        {
          slotName: 'alarm.trigger.TriggerType',
          width: '10%',
          scopedSlots: { customRender: 'TriggerType',title: 'alarm.trigger.TriggerType'},
          dataIndex: 'TriggerType',
        },
        {
          slotName: 'alarm.trigger.TriggerCondition',
          width: '10%',
          scopedSlots: { customRender: 'TriggerCondition',title: 'alarm.trigger.TriggerCondition'},
          dataIndex: 'TriggerCondition',
        },
        {
          slotName: 'alarm.trigger.TriggerUpdatedAt',
          width: '10%',
          scopedSlots: { customRender: 'UpdatedAt',title: 'alarm.trigger.TriggerUpdatedAt'},
          dataIndex: 'UpdatedAt',
        },
        {
          slotName: 'dataModel.modelTableOpt',
          width: '10%',
          scopedSlots: { customRender: 'action',title: 'dataModel.modelTableOpt'}
        }
      ],
      dataSource: [],
      conditionExpress:"",
      selectedRows: []
    }
  },
  authorize: {
    // deleteRecord: 'delete'
  },
  mounted(){
    this.getSupportDevice()
  },
  activated(){

  },
  filters: {
    formatDate(time) {
      let date = new Date(time)
      return formatDate(date,'yyyy-MM-dd hh:mm:ss')
    },
  },
  created(){
    this.dataSource=[]
    this.getModelList(0)
    this.GetTriggerList()
  },
  watch: {
    '$route' () {
      this.getSupportDevice()
    }
  },
  methods: {
    onTriggerAlarmSwitchChange(checked) {
      this.TriggerAlarmSwitchStatus = checked
    },
    onTriggerLinkageSwitch(checked) {
      this.TriggerLinkageSwitchStatus = checked
    },
    getDeviceDataModel(type,uuid){
      let params={
        muid:uuid,
        type:this.DeviceType
      }
      if(type==1)
      {
        this.dataModelList = []
      }
      else{
        this.dataLinkModelList = []
      }

      let _t = this
      this.messageShowLoad = true
      getDatasByUuid(params).then(function (res){
        _t.messageShowLoad = false
        if(res.data.mibs.length>0)
        {
          let temp={}
          for(let i=0;i<res.data.mibs.length;i++)
          {
            temp.id = i.toString()
            temp.name = res.data.mibs[i].name
            temp.uuid = res.data.mibs[i].uuid
            if(type==1)
            {
              _t.dataModelList.push(temp)
            }
            else{
              _t.dataLinkModelList.push(temp)
            }

            temp={}
          }
        }
      })
    },
    changeDeviceType(value){
      this.DeviceType = value
      this.dataModelList=[]
      this.supportDeviceModelList=[]
      if(this.isEdit!=true)
      {
        this.form.setFieldsValue(
            {
              DeviceModel:"",
              DeviceModelData:"",
            })
      }

      this.getModelList(1)
    },
    changeDeviceModel(value){
      this.dataModelList=[]
      if(this.isEdit!=true) {
        this.form.setFieldsValue(
            {
              DeviceModelData: "",
            })
      }
      this.getDeviceDataModel(1,value)
    },
    changeCondition(value){
        this.conditionExpress=value
    },
    getSupportDevice(){
      let _t = this
      getSupportDeviceList().then(function (res){
        _t.supportDeviceList =res.data.list
        _t.supportLinkDeviceList =res.data.list
      })
    },
    getModelList(type){
      let _t = this
      if(type==1)
      {
        this.supportDeviceModelList=[]
      }
      else if(type==2){
        this.supportLinkDeviceModelList=[]
      }
      else{
        this.supportDeviceModelList=[]
        this.supportLinkDeviceModelList=[]
      }

      const  params= {
        type:this.DeviceType
      }
      snmpModelList(params).then(function (res){
        let tempData={}
        if(res.data.list!=null)
        {
          for(let i=0;i<res.data.list.length;i++)
          {
            tempData.uuid = res.data.list[i].uuid
            tempData.modelName = res.data.list[i].name
            if(type==1)
            {
              _t.supportDeviceModelList.push(tempData)
            }
            else  if(type==2){
              _t.supportLinkDeviceModelList.push(tempData)
            }
            else{
              _t.supportDeviceModelList.push(tempData)
              _t.supportLinkDeviceModelList.push(tempData)
            }
            tempData={}
          }
        }
      })
    },
    filterOption(input, option) {
      return (
          option.componentOptions.children[0].text.toLowerCase().indexOf(input.toLowerCase()) >= 0
      );
    },
    deleteRecord(recond) {
      let params={
        ID:parseInt(recond.ID),
        name:recond.TriggerName
      }
      let _t = this
      AlarmTriggerDel(params).then(function (res) {
        if(res.data.code==0)
        {
          _t.dataSource = _t.dataSource.filter(item => item.ID !== recond.ID)
          _t.selectedRows = _t.selectedRows.filter(item => item.ID !== recond.ID)
        }
        else if(res.data.code==2004)
        {
            _t.$message.error(_t.$t("dataModel.modelBand"))
        }
      }).catch(function(){
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
    },
    ShowAddTrigger(){
      this.form.setFieldsValue(
          {
            TriggerAddName:"",
            DeviceType:"",
            DeviceModel:"",
            DeviceModelData:"",
            TriggerCondition:"",
            XValue:"",
            TriggerAlarmShowText:"",
            TriggerAlarmHideText:"",
            TriggerAlarmSwitch:false,
            TriggerLinkageSwitch:false,
          })
      this.TriggerAlarmSwitchStatus=false
      this.TriggerLinkageSwitchStatus=false
      this.visible = true;
      this.isEdit = false;
    },
    AddTriggerBtn(){
      if((this.TriggerAlarmSwitchStatus==false)&&(this.TriggerLinkageSwitchStatus==false))
      {
        this.$message.error(this.$t('alarm.trigger.TriggerParamsError'), 3)
        return
      }
      this.form.validateFields((err) => {
        if (!err) {
          let params = {
            TriggerName:this.form.getFieldValue('TriggerAddName'),
            TriggerDeviceModelUuid:this.form.getFieldValue('DeviceModel'),
            TriggerModelDataUuid:this.form.getFieldValue('DeviceModelData'),
            TriggerCondition:this.form.getFieldValue('TriggerCondition'),
            TriggerDeviceType:parseInt(this.form.getFieldValue('DeviceType')),
            TriggerXValue:this.form.getFieldValue('XValue'),
            TriggerKeepTime:parseInt(this.form.getFieldValue('TriggerKeepTime')),
          }
          if ((params.TriggerCondition=='&&')||(params.TriggerCondition=='||')){
            params.TriggerYValue=this.form.getFieldValue('YValue')
          }
          if(this.TriggerAlarmSwitchStatus)
          {
            params.TriggerAlarmShowText=this.form.getFieldValue('TriggerAlarmShowText')
            params.TriggerAlarmHideText=this.form.getFieldValue('TriggerAlarmHideText')
            params.TriggerAlarmLevel=parseInt(this.form.getFieldValue('AlarmLevel'))
          }
          if(this.TriggerLinkageSwitchStatus)
          {
            params.TriggerLinkDeviceType=parseInt(this.form.getFieldValue('DeviceType')),
            params.TriggerLinkdeviceModelUuid=this.form.getFieldValue('DeviceModel')
            params.TriggerLinkModelDataUuid=this.form.getFieldValue('DeviceLinkModelData')
            params.TriggerLinkageAlarmValue=this.form.getFieldValue('TriggerLinkageAlarmValue')
            params.TriggerLinkageAlarmClearValue=this.form.getFieldValue('TriggerLinkageAlarmClearValue')
          }
          if((this.TriggerAlarmSwitchStatus==true)&&(this.TriggerLinkageSwitchStatus==true))
          {
            params.TriggerType = 3
          }
          else if(this.TriggerLinkageSwitchStatus){
            params.TriggerType = 1
          }
          else if(this.TriggerAlarmSwitchStatus){
            params.TriggerType = 2
          }
          else{
            this.$message.error(this.$t('alarm.trigger.TriggerParamsError'), 3)
            return
          }

          let _t = this
          AlarmTriggerAdd(params).then(function (res) {
            if(res.data.code==0)
            {
              _t.GetTriggerList()
              _t.visible = false
              _t.$message.success(_t.$t('alarm.trigger.AddSuccess'), 3)
            }
            else if(res.data.code==3001){
              _t.$message.error(_t.$t('alarm.trigger.NameExist'), 3)
            }
            else if(res.data.code==3004){
              _t.$message.error(_t.$t('alarm.trigger.BangDing'), 3)
            }
            else {
              _t.$message.error(_t.$t('alarm.trigger.AddFailed'), 3)
            }
          }).catch(function(){
            _t.$message.error(_t.$t('loginPage.serverError'), 3)
          })
        }
      })
    },
    EditTriggerBtn(){
      if((this.TriggerAlarmSwitchStatus==false)&&(this.TriggerLinkageSwitchStatus==false))
      {
        this.$message.error(this.$t('alarm.trigger.TriggerParamsError'), 3)
        return
      }
      this.form.validateFields((err) => {
        if (!err) {
          let params = {
            ID:parseInt(this.editTriggerID),
            TriggerName:this.form.getFieldValue('TriggerAddName'),
            TriggerDeviceModelUuid:this.form.getFieldValue('DeviceModel'),
            TriggerModelDataUuid:this.form.getFieldValue('DeviceModelData'),
            TriggerCondition:this.form.getFieldValue('TriggerCondition'),
            TriggerDeviceType:parseInt(this.form.getFieldValue('DeviceType')),
            TriggerXValue:this.form.getFieldValue('XValue'),
            TriggerKeepTime:parseInt(this.form.getFieldValue('TriggerKeepTime')),
          }
          if ((params.TriggerCondition=='&&')||(params.TriggerCondition=='||')){
            params.TriggerYValue=this.form.getFieldValue('YValue')
          }
          if(this.TriggerAlarmSwitchStatus)
          {
            params.TriggerAlarmShowText=this.form.getFieldValue('TriggerAlarmShowText')
            params.TriggerAlarmHideText=this.form.getFieldValue('TriggerAlarmHideText')
            params.TriggerAlarmLevel=parseInt(this.form.getFieldValue('AlarmLevel'))
          }
          if(this.TriggerLinkageSwitchStatus)
          {
            params.TriggerLinkDeviceType=parseInt(this.form.getFieldValue('DeviceType'))
            params.TriggerLinkdeviceModelUuid=this.form.getFieldValue('DeviceModel')
            params.TriggerLinkModelDataUuid=this.form.getFieldValue('DeviceLinkModelData')
            params.TriggerLinkageAlarmValue=this.form.getFieldValue('TriggerLinkageAlarmValue')
            params.TriggerLinkageAlarmClearValue=this.form.getFieldValue('TriggerLinkageAlarmClearValue')
          }
          if((this.TriggerAlarmSwitchStatus==true)&&(this.TriggerLinkageSwitchStatus==true))
          {
            params.TriggerType = 3
          }
          else if(this.TriggerLinkageSwitchStatus){
            params.TriggerType = 1
          }
          else if(this.TriggerAlarmSwitchStatus){
            params.TriggerType = 2
          }
          let _t = this
          AlarmTriggerEdit(params).then(function (res) {
            if(res.data.code==0)
            {
              _t.GetTriggerList()
              _t.visible = false
              _t.$message.success(_t.$t('alarm.trigger.EditSuccess'), 3)
            }
            else if(res.data.code==3001){
              _t.$message.error(_t.$t('alarm.trigger.NameExist'), 3)
            }
            else if(res.data.code==3002){
              _t.$message.error(_t.$t('alarm.trigger.AddFailed'), 3)
            }
            else if(res.data.code==3004){
              _t.$message.error(_t.$t('alarm.trigger.BangDing'), 3)
            }
            else {
              _t.$message.error(_t.$t('alarm.trigger.AddFailed'), 3)
            }
          }).catch(function(){
            _t.$message.error(_t.$t('loginPage.serverError'), 3)
          })
        }
      })
    },
    GetTriggerList(){
      let _t = this
      this.dataSource=[]
      this.messageShowLoad=true
      GetAlarmTriggerList().then(function (res) {
        _t.messageShowLoad=false
        if(res.data.code==0)
        {
          _t.dataSource = res.data.list
        }
      }).catch(function(){
        _t.messageShowLoad=false
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
    },
    gotoEditTrigger(item){
      let _t = this
      this.isEdit = true
      this.editTriggerID = item.ID
      this.ShowRegisterLoading = true
      this.changeDeviceType(item.TriggerDeviceType)
      this.changeDeviceModel(item.TriggerDeviceModelUuid)

      if(item.TriggerType==3){
        this.TriggerAlarmSwitchStatus=true
        this.TriggerLinkageSwitchStatus=true
      }else if(item.TriggerType==1){
        this.TriggerAlarmSwitchStatus=false
        this.TriggerLinkageSwitchStatus=true
      }else if(item.TriggerType==2){
        this.TriggerLinkageSwitchStatus=false
        this.TriggerAlarmSwitchStatus=true
      }
      this.conditionExpress = item.TriggerCondition
      setTimeout(function (){
        _t.form.setFieldsValue(
            {
              TriggerAddName:item.TriggerName,
              DeviceType:item.TriggerDeviceType,
              DeviceModel:item.TriggerDeviceModelUuid,
              DeviceModelData:item.TriggerModelDataUuid,
              TriggerCondition:item.TriggerCondition,
              XValue:item.TriggerXValue,
              TriggerKeepTime:item.TriggerKeepTime.toString()
            })
          if ((item.TriggerCondition=='&&')||(item.TriggerCondition=='||')){
            _t.form.setFieldsValue(
                {
                  YValue:item.TriggerYValue,
                })
          }

          if(item.TriggerType==3){
            _t.form.setFieldsValue(
                {
                  TriggerAlarmSwitch:true,
                  TriggerLinkageSwitch:true,
                  TriggerAlarmShowText:item.TriggerAlarmShowText,
                  TriggerAlarmHideText:item.TriggerAlarmHideText,
                  AlarmLevel:item.TriggerAlarmLevel.toString(),
                  DeviceLinkModelData:item.TriggerLinkModelDataUuid,
                  TriggerLinkageAlarmValue:item.TriggerLinkageAlarmValue,
                  TriggerLinkageAlarmClearValue:item.TriggerLinkageAlarmClearValue
                })
          }else if(item.TriggerType==1){
            _t.form.setFieldsValue(
                {
                  TriggerAlarmSwitch:false,
                  TriggerLinkageSwitch:true,
                  DeviceLinkModelData:item.TriggerLinkModelDataUuid,
                  TriggerLinkageAlarmValue:item.TriggerLinkageAlarmValue,
                  TriggerLinkageAlarmClearValue:item.TriggerLinkageAlarmClearValue
                })
          }else if(item.TriggerType==2){
            _t.form.setFieldsValue(
                {
                  TriggerAlarmSwitch:true,
                  TriggerLinkageSwitch:false,
                  TriggerAlarmShowText:item.TriggerAlarmShowText,
                  TriggerAlarmHideText:item.TriggerAlarmHideText,
                  AlarmLevel:item.TriggerAlarmLevel.toString(),
                })
          }
        _t.ShowRegisterLoading = false
      },1000)
      this.visible = true;
    },
    afterVisibleChange(val) {
      if(!val)
      {
        this.edit = false
      }
    },
    onClose() {
      this.visible = false;
    },
  }
}
</script>

<style lang="less">

.plus-icon-enter-active{
  transition: opacity .5s;
}
.plus-icon-enter{
  opacity: 0;
}
.plus-icon-leave-active{
  transition: opacity .5s;
}
.plus-icon-leave-to{
  opacity: 0;
}
.plus-icon-enter-to{
  opacity: 1;
}

.code-box-actions {
  padding-top: 12px;
  text-align: center;
  opacity: .7;
  transition: opacity .3s;
}
.code-box-meta .demo-description>h4, .code-box-meta>h4 {
  position: absolute;
  top: -14px;
  padding: 1px 8px;
  margin-left: 16px;
  color: #777;
  border-radius: 2px 2px 0 0;
  background: #fff;
  font-size: 14px;
  width: auto;
}
.code-box {
  border: 1px solid #ebedf0;
  border-radius: 2px;
  display: inline-block;
  width: 100%;
  position: relative;
  padding: 20px;
  transition: all .2s;
}

.search{
  margin-bottom: 54px;
}
.fold{
  width: calc(100% - 216px);
  display: inline-block
}
.operator{
  margin-bottom: 18px;
}
@media screen and (max-width: 900px) {
  .fold {
    width: 100%;
  }
}

.ant-table-thead > tr > th
{
  padding: 10px 10px;
  overflow-wrap: break-word;
}
.ant-table-tbody > tr > td {
  padding: 7px 5px;
  overflow-wrap: break-word;
}

.ant-table-thead>tr>th {
  color: #909399;
  font-weight: 500;
  text-align: left;
  //background: #f8f8f8;
  //border-bottom: 1px solid #e8e8e8;
  transition: background .3s ease;
}
</style>
