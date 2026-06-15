<template>
  <div>
<!--    寄存器组的寄存器地址配置-->
    <a-card  v-if="!registerGroupListTable" style="min-height: 400px">
      <a-space class="operator">

        <a-button type="primary" @click="registerGroupListTable=true"> <a-icon type="backward" />
          {{$t("dataModel.modbusModel.Back")}}</a-button>

        <a-button type="danger" :disabled="selectDataTableUuid.length==0" @click="deleteRecord"> <a-icon type="delete"  />
          {{$t('dataModel.delete')}}</a-button>

      </a-space>
      <a-spin style="padding: 1px;"  :spinning="messageShowLoad" tip="Loading...">
        <a-table :loading="false" :pagination="pagination" row-key="uuid" :row-selection="rowSelection" :columns="columns" :data-source="dataSource" class="ant-table-tbody">
          <template v-for="(item, index) in columns" :slot="item.slotName">
            <span :key="index">{{ $t(item.slotName) }}</span>
          </template>
          <template slot="oidName" slot-scope="text, record">
            <a-input :value="text" @change="onCellChange(record, $event)"/>
          </template>
          <template slot="action" slot-scope="text, record">
            <div class="editable-row-operations">
            <span >
              <a  @click="() => edit(record)">
                <a-icon type="edit"/> {{$t('dataModel.edit')}}</a>
            </span>
            </div>
          </template>
        </a-table>
      </a-spin>
      <div>
      <a-spin style="z-index:9999;margin-top: -600px;margin-right: -300px" size="large" :spinning="ShowRegisterLoading" tip="Loading...">
        <a-drawer
          :title="$t('dataModel.EditDataModel')"
          :width="720"
          :visible="editVisible"
          :body-style="{ paddingBottom: '80px' }"
          @close="onClose"
      >
        <a-form :form="EditForm" layout="vertical" >
          <a-row :gutter="16">
            <a-col :span="12">
              <a-form-item :label="$t('dataModel.editData.dataName')">
                <a-input
                    v-decorator="[
                  'name',
                  {
                    rules: [{ required: true, message: $t('dataModel.editData.dataName') }],
                  },
                ]"
                />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item :label="$t('dataModel.editData.dataType')">
                <a-select   autocomplete="autocomplete"  @change="changeDataType"  v-decorator="[
                  'dataType',
                  {
                    rules: [{ required: true, message: $t('dataModel.editData.dataType') }],
                  },
                ]">
                  <a-select-option value="Short">Signed</a-select-option>
                  <a-select-option value="Unsigned short">Unsigned</a-select-option>
                  <a-select-option value="Long">Long</a-select-option>
                  <a-select-option value="Float">Float</a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
          </a-row>
          <a-row :gutter="16">
            <a-col :span="12" v-if="dataType=='Long'||dataType=='Float'">
              <a-form-item :label="$t('dataModel.editData.DataOrder')">
                <a-select   autocomplete="autocomplete"  v-decorator="[
                  'ByteOrder',
                  {
                    rules: [{ required: false, message: $t('dataModel.editData.DataOrder') }],
                  },
                ]">
                  <a-select-option value="ABCD">ABCD</a-select-option>
                  <a-select-option value="CDAB">CDAB</a-select-option>
                  <a-select-option value="BADC">BADC</a-select-option>
                  <a-select-option value="DCBA">DCBA</a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :span="12" >

              <a-form-item :label="$t('dataModel.modbusModel.ConversionExpression')">
                <a-tooltip placement="top">
                  <template slot="title">
                    <span>{{$t('dataModel.modbusModel.ConversionExpressionTips')}}</span>
                  </template>
                  <a-input   autocomplete="autocomplete"   v-decorator="[
                  'ConversionExpression',
                  {
                    rules: [{ required: false, message: $t('dataModel.modbusModel.ConversionExpression') }],
                  },
                ]">
                  </a-input>
                </a-tooltip>
              </a-form-item>
            </a-col>

            <a-col :span="12" >
              <a-form-item :label="$t('dataModel.editData.dataAuth')">
                <a-select   autocomplete="autocomplete"  v-decorator="[
                  'dataAuth',
                  {
                    rules: [{ required: true, message: $t('dataModel.editData.dataAuth') }],
                  },
                ]">
                  <a-select-option value="ReadOnly">ReadOnly</a-select-option>
                  <a-select-option value="ReadWrite">ReadWrite</a-select-option>
                </a-select>
              </a-form-item>
            </a-col>

            <a-col :span="12" >
              <a-form-item :label="$t('dataModel.editData.dataUnit')">
                <a-input   autocomplete="autocomplete"   v-decorator="[
                  'dataUnit',
                  {
                    rules: [{ required: false, message: $t('dataModel.editData.dataUnit') }],
                  },
                ]">
                </a-input>
              </a-form-item>
            </a-col>

            <a-col :span="12" >

              <a-form-item :label="$t('dataModel.editData.dataAlarm')">
                <a-tooltip placement="top">
                  <template slot="title">
                    <span>{{$t('dataModel.editData.dataAlarmTips')}}</span>
                  </template>
                  <a-select  @change="alarmCharge"   autocomplete="autocomplete"  v-decorator="[
                  'dataAlarm',
                  {
                    rules: [{ required: true, message: $t('dataModel.editData.dataAlarm') }],
                  },
                ]">
                    <a-select-option value="1">{{$t('dataModel.editData.dataAlarmYes')}}</a-select-option>
                    <a-select-option value="0">{{$t('dataModel.editData.dataAlarmNo')}}</a-select-option>
                  </a-select>
                </a-tooltip>
              </a-form-item>

            </a-col>

            <!--            告警等级-->
            <div v-if="alarmStatus">
              <a-col :span="12" >
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
              <a-col :span="12">
                <a-form-item :label="$t('dataModel.editData.AlarmMessage')">
                  <a-input   autocomplete="autocomplete"   v-decorator="[
                  'AlarmMessage',
                  {
                    rules: [{ required: true, message: $t('dataModel.editData.AlarmMessage') }],
                  },
                ]">
                  </a-input>
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item :label="$t('dataModel.editData.AlarmClearMessage')">
                  <a-input   autocomplete="autocomplete"   v-decorator="[
                  'AlarmClearMessage',
                  {
                    rules: [{ required: false, message: $t('dataModel.editData.AlarmClearMessage') }],
                  },
                ]">
                  </a-input>
                </a-form-item>
              </a-col>
            </div>


            <!--存储            -->
            <div v-else>
              <a-col :span="12" >
                <a-form-item :label="$t('dataModel.editData.dataRecord')">
                  <a-select   @change="recordCharge" autocomplete="autocomplete"  v-decorator="[
                  'dataRecord',
                  {
                    rules: [{ required: true, message: $t('dataModel.editData.dataRecord') }],
                  },
                ]">
                    <a-select-option value='1' selectd>{{$t('dataModel.editData.dataAlarmYes')}}</a-select-option>
                    <a-select-option value='0'>{{$t('dataModel.editData.dataAlarmNo')}}</a-select-option>
                  </a-select>
                </a-form-item>
              </a-col>
              <a-col :span="12" v-if="recordStatus" >
                <a-form-item :label="$t('dataModel.dataRecordType')">
                  <a-select   autocomplete="autocomplete"  @change="chargeDataRecordType" v-decorator="[
                  'dataRecordType',
                  {
                    rules: [{ required: true, message: $t('dataModel.dataRecordType') }],
                    initialValue: '0'
                  },
                ]">
                    <a-select-option value=1>{{$t('dataModel.dataRecordTimely')}}</a-select-option>
                    <a-select-option value=0>{{$t('dataModel.dataRecordCharge')}}</a-select-option>
                    <a-select-option value=2>{{$t('dataModel.dataRecordNow')}}</a-select-option>
                    <a-select-option value=3>{{$t('dataModel.dataRecordChangeRate')}}</a-select-option>
                  </a-select>
                </a-form-item>
              </a-col>
              <a-col :span="12" v-if="(recordStatus)&&((DataRecordType==0)||(DataRecordType==3))">
                <a-form-item :label="DataRecordType==0?$t('dataModel.dataRecordChargeValue'):$t('dataModel.dataRecordChangeRateValue')">
                  <a-input   :min="1"  style="width: 100%" autocomplete="autocomplete"  v-decorator="[
                  'dataRecordChargeValue',
                  {
                    rules: [{ required: true, message: $t('dataModel.dataRecordChargeValue') }],
                  },
                ]">
                  </a-input>
                </a-form-item>
              </a-col>
              <a-col :span="12" v-if="(recordStatus)&&(DataRecordType==1)">
                <a-form-item :label="$t('dataModel.editData.dataRecordTime')">
                  <a-input-number   :min="1"  style="width: 100%" autocomplete="autocomplete"  v-decorator="[
                  'dataRecordTime',
                  {
                    rules: [{ required: true, message: $t('dataModel.editData.dataRecordTime') }],
                  },
                ]">
                  </a-input-number>
                </a-form-item>
              </a-col>
            </div>
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
          <a-button  type="primary" :style="{ marginRight: '8px' }" @click="save()">
            {{$t('device.EditButton')}}
          </a-button>

          <a-button  @click="onClose">
            {{$t('device.CancelButton')}}
          </a-button>
        </div>
      </a-drawer>
      </a-spin>
      </div>
    </a-card>

<!--    寄存器组表格-->
    <a-card  v-if="registerGroupListTable" style="min-height: 400px">
      <a-space class="operator">

        <a-button type="primary" @click="RegisterVisible=true"> <a-icon type="plus" />
          {{$t("dataModel.modbusModel.RegisterGroup")}}</a-button>

        <a-button type="default" @click="onBlackCLK()"> <a-icon type="backward" />
          {{$t("dataModel.modbusModel.Back")}}</a-button>

      </a-space>
      <a-spin style="padding: 1px;"  :spinning="messageShowLoad" tip="Loading...">
        <a-table  :pagination="pagination" row-key="uuid" :columns="registerGroupColumns" :data-source="registerGroupDataSource" class="ant-table-tbody">
          <template v-for="(item, index) in registerGroupColumns" :slot="item.slotName">
            <span :key="index">{{ $t(item.slotName) }}</span>
          </template>
          <template slot="function" slot-scope="text">
            <span v-if="text==1"> 0{{text}} {{$t('dataModel.modbusModel.FunctionReadCoils')}}</span>
            <span v-else-if="text==2"> 0{{text}} {{$t('dataModel.modbusModel.FunctionReadDisCrete')}}</span>
            <span v-else-if="text==3"> 0{{text}} {{$t('dataModel.modbusModel.FunctionReadHoldingRegisters')}}</span>
            <span v-else-if="text==4"> 0{{text}} {{$t('dataModel.modbusModel.FunctionReadInputRegisters')}}</span>
          </template>
          <template slot="action" slot-scope="text, record">
            <div class="editable-row-operations">
              <span >
                <a  @click="() => registerConfig(record.uuid)">
                  <a-icon type="caret-right" /> {{$t('dataModel.modbusModel.RegisterGroupConfig')}}</a>
              </span>
              <a-divider type="vertical" />
              <a-popconfirm :title="$t('dataModel.deleteConfirm')" @confirm="deleteRegisterGroupRecord(record.uuid)">
                <a-icon slot="icon" type="question-circle-o" style="color: red" />
                <a-icon type="delete" theme="twoTone" two-tone-color="#eb2f96"/><a style="color: #eb2f96">{{$t('dataModel.delete')}}</a>
              </a-popconfirm>
            </div>
          </template>
        </a-table>
      </a-spin>
    </a-card>
<!--    添加寄存器组-->
    <a-modal  size="large" :bodyStyle="{ width: '800px',height: '210px'}"  v-model="RegisterVisible" :title="$t('dataModel.modbusModel.RegisterGroup')" on-ok="handleOk">
      <template slot="footer">

        <div >
          <a-button key="submit" type="primary" @click="addRegisterGroup">
            {{ $t('component.deviceDataModel.submit')}}
          </a-button>
          <a-button key="back" @click="RegisterVisible=false">
            {{$t('component.deviceDataModel.cancel')}}
          </a-button>
        </div>

      </template>
        <a-form :form="RegisterForm" layout="vertical" @submit="addRegisterGroup">
          <a-alert type="error" :closable="true" v-show="error" :message="error" showIcon style="margin-bottom: 24px;" />

          <a-form-item
              :label="$t('dataModel.modbusModel.RegisterGroupName')"
              :labelCol="{span: 4}"
              :wrapperCol="{span: 10}"
          >
            <a-input autocomplete="autocomplete"
                     v-decorator="['RegisterGroupName', {rules: [{ required: true, message: $t('dataModel.modbusModel.RegisterGroupName'), whitespace: true}]}]"
            />
          </a-form-item>

          <a-form-item :label="$t('dataModel.modbusModel.RegisterFunction')"
                       :labelCol="{span: 4}"
                       :wrapperCol="{span: 10}"
          >
            <a-select
                v-decorator="['RegisterFunction', {rules: [{ required: true, message: $t('dataModel.modbusModel.RegisterFunction'), whitespace: true}]}]"
            >
              <a-select-option value="1">
                01 {{$t('dataModel.modbusModel.FunctionReadCoils')}}
              </a-select-option>
              <a-select-option value="2">
                02 {{$t('dataModel.modbusModel.FunctionReadDisCrete')}}
              </a-select-option>
              <a-select-option value="3">
                03 {{$t('dataModel.modbusModel.FunctionReadHoldingRegisters')}}
              </a-select-option>
              <a-select-option value="4">
                04 {{$t('dataModel.modbusModel.FunctionReadInputRegisters')}}
              </a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item
              :label="$t('dataModel.modbusModel.RegisterStartAddress')"
              :labelCol="{span: 4}"
              :wrapperCol="{span: 10}"
          >
            <a-input autocomplete="autocomplete"
                     v-decorator="['RegisterStartAddress', {rules: [{ required: true, message: $t('dataModel.modbusModel.RegisterStartAddress'), whitespace: true}]}]"
            />
          </a-form-item>
          <a-form-item
              :label="$t('dataModel.modbusModel.RegisterCount')"
              :labelCol="{span: 4}"
              :wrapperCol="{span: 10}"
          >
            <a-input autocomplete="autocomplete"
                     v-decorator="['RegisterCount', {rules: [{ required: true, message: $t('dataModel.modbusModel.RegisterCount'), whitespace: true}]}]"
            />
          </a-form-item>
        </a-form>
    </a-modal>
  </div>
</template>

<script>
import {
  getSnmpModelDetail, modelDataEdit, snmpModelDeleteMibs, snmpModelEdit,
} from "../../../services/snmpmodel";

import {
  modbusModelEdit,
  modbusModelGroupAdd,
  modbusModelGroupDel,
  modbusModelGroupList, modbusModelRegisterDel,
  modbusModelRegisterEdit,
  modbusModelRegisterList
} from "../../../services/modbusModel";

import { uuid } from 'vue-uuid';
import {COMListGet} from "@/services/modbusModel";

const dataSource= []

export default {
  name: 'ModbusModelRegister',
  i18n: require('../../../i18n/language'),
  components: {
  },
  data () {
    return {
      pagination:{
        pageSize:15,
        showSizeChanger:true
      },
      DataRecordType:0,
      registerGroupListTable:true,
      error: '',
      ShowRegisterLoading:false,
      alarmStatus:0,
      recordStatus:0,
      EditForm:this.$form.createForm(this),
      RegisterForm: this.$form.createForm(this),
      selectDataTableUuid:[],
      rowSelection:{
        onSelect:this.onDataTableSelect,
        onSelectAll:this.onDataTableSelectAll
      },
      RegisterVisible:false,
      editVisible:false,
      messageShowLoad:false,
      RegisterMessageShowLoad:false,
      columns: [
        {
          slotName: this.$t("dataModel.oidName"),
          scopedSlots: {  customRender: 'oidName' ,title:this.$t("dataModel.oidName")},
          width: '20%',
          dataIndex: 'name'
        },
        {
          slotName: this.$t("dataModel.modbusModel.RegisterAddress"),
          scopedSlots: {  customRender: 'RegisterAddress' ,title:this.$t("dataModel.modbusModel.RegisterAddress")},
          width: '20%',
          align:"center",
          dataIndex: 'registerAddress'
        },
        {
          slotName:this.$t("dataModel.oidType"),
          scopedSlots: {  customRender: 'oidType' ,title:this.$t("dataModel.oidType") },
          width: '15%',
          align:"center",
          dataIndex: 'type',
        },
        {
          slotName:this.$t("dataModel.oidAuth"),
          scopedSlots: {  customRender: 'oidAuth'  ,title:this.$t("dataModel.oidAuth")},
          width: '15%',
          align:"center",
          dataIndex: 'auth',
        },
        {
          slotName:this.$t("dataModel.DataUnit"),
          scopedSlots: {  customRender: 'dataUnit'  ,title:this.$t("dataModel.DataUnit")},
          width: '10%',
          align:"center",
          dataIndex: 'unit'
        },
        {
          title: this.$t('dataModel.modelTableOpt'),
          width: '15%',
          scopedSlots: { customRender: 'action' }
        }
      ],
      registerGroupColumns: [
        {
          slotName: this.$t("dataModel.modbusModel.RegisterGroupName"),
          scopedSlots: {  customRender: 'name' ,title:this.$t("dataModel.modbusModel.RegisterGroupName")},
          width: '15%',
          dataIndex: 'name'
        },
        {
          slotName:this.$t("dataModel.modbusModel.RegisterFunction"),
          scopedSlots: {  customRender: 'function'  ,title:this.$t("dataModel.modbusModel.RegisterFunction")},
          width: '15%',
          align:"left",
          dataIndex: 'function'
        },
        {
          slotName:this.$t("dataModel.modbusModel.RegisterStartAddress"),
          scopedSlots: {  customRender: 'registerStart' ,title:this.$t("dataModel.modbusModel.RegisterStartAddress") },
          width: '10%',
          align:"center",
          dataIndex: 'registerStart',
        },
        {
          slotName:this.$t("dataModel.modbusModel.RegisterCount"),
          scopedSlots: {  customRender: 'registerCount'  ,title:this.$t("dataModel.modbusModel.RegisterCount")},
          width: '10%',
          align:"center",
          dataIndex: 'registerCount',
        },
        {
          title: this.$t('dataModel.modelTableOpt'),
          width: '20%',
          scopedSlots: { customRender: 'action' }
        }
      ],
      dataSource,
      dataType:"Short",
      registerGroupDataSource:[],
      selectedRows: [],
    }
  },
  created(){
    this.registerGroupList()
    this.registerGroupListTable=true
  },
  activated() {

  },
  mounted() {

  },
  computed: {

  },
  methods: {
    onCellChange(rc, e) {
      let params = {
        uuid:rc.uuid,
        name:e.target.value,
        auth:rc.auth,
        type:rc.type,
        ByteOrder:rc.ByteOrder,
        unit:rc.unit,
        conversionExpression:rc.conversionExpression,
        alarm:rc.alarm,
        alarmLevel:rc.alarmLevel,
        AlarmMessage:rc.AlarmMessage,
        AlarmClearMessage:rc.AlarmClearMessage,
        record:rc.record,
        RecordType:rc.RecordType,
        recordInterval:rc.recordInterval,
        RecordDataCharge:rc.RecordDataCharge,
      }
      let _t = this
      modbusModelRegisterEdit(params).then(function (res){
        if(res.data.code==200)
        {
          const newData = [..._t.dataSource];
          const target = newData.filter(item => _t.editingKey === item.uuid)[0];
          if (target) {
            target.name = e.target.value
            target.auth=rc.auth
            target.type=rc.type
            target.ByteOrder=rc.ByteOrder
            target.unit=rc.unit
            target.conversionExpression=rc.conversionExpression
            target.alarm=rc.alarm
            target.alarmLevel=rc.alarmLevel
            target.AlarmMessage = rc.AlarmMessage
            target.AlarmClearMessage =rc.AlarmClearMessage
            target.record=rc.record
            target.RecordType=rc.RecordType
            target.recordInterval=rc.recordInterval
            target.RecordDataCharge=rc.RecordDataCharge
            _t.dataSource = newData;
          }
          // _t.$message.success(_t.$t("dataModel.saveSuccess"));
          _t.editVisible = false;
        }
        else
        {
          _t.$message.error(_t.$t("dataModel.saveFailed"));
        }
      })
    },
    changeDataType(value){
      this.dataType = value
    },
    onClose() {
      this.editVisible = false;
    },
    findDataTableSelectIndex (key) {
      for(let i=0;i<this.selectDataTableUuid.length;i++)
      {
        if(this.selectDataTableUuid[i]==key)
        {
          return i
        }
      }
      return -1
    },
    onDataTableSelect (record, selected, selectedRows)  {
      if(selected)
      {
        this.selectDataTableUuid.push(record.uuid)
      }
      else
      {
        const index=this.findDataTableSelectIndex(record.uuid)
        if(index!=-1)
        {
          this.selectDataTableUuid.splice(index,1)
        }
      }
    },
    onDataTableSelectAll(selected, selectedRows, changeRows) {
      if(selected)
      {
        for(let i=0;i<selectedRows.length;i++)
        {
          const index=this.findDataTableSelectIndex(selectedRows[i].uuid)
          if(index==-1)
          {
            this.selectDataTableUuid.push(selectedRows[i].uuid)
          }
        }

      }
      else
      {
        for(let i=0;i<changeRows.length;i++)
        {
          const index=this.findDataTableSelectIndex(changeRows[i].uuid)
          if(index!=-1)
          {
            this.selectDataTableUuid.splice(index,1)
          }
        }
      }
    },
    edit(item) {
      let _t = this
      this.editingKey = item.uuid
      this.alarmStatus = item.alarm
      this.recordStatus = item.record
      this.ShowRegisterLoading = true
      this.DataRecordType = item.RecordType
      if(item.recordInterval==0)
      {
        item.recordInterval=1
      }
      this.dataType = item.type
      setTimeout(function (){
        _t.EditForm.setFieldsValue(
            {
              name:item.name,
              dataType:item.type,
              dataAuth:item.auth,
              dataUnit:item.unit,
              ConversionExpression:item.conversionExpression,
            })
            if(item.type=="Long"||item.type=="Float")
            {
              _t.EditForm.setFieldsValue(
                  {
                    ByteOrder:item.ByteOrder,
                  })
            }
            if (item.alarm==1){
              _t.EditForm.setFieldsValue(
                  {
                    dataAlarm:item.alarm.toString(),
                    AlarmLevel:item.alarmLevel.toString(),
                    dataRecord:item.record.toString(),
                    AlarmMessage :item.AlarmMessage,
                    AlarmClearMessage : item.AlarmClearMessage,
                  })
            }
            else  if (item.record==1)
            {
              _t.EditForm.setFieldsValue(
                  {
                    dataAlarm:item.alarm.toString(),
                    dataRecord:item.record.toString(),
                    dataRecordType:item.RecordType.toString(),
                    dataRecordChargeValue:item.RecordDataCharge.toString(),
                    dataRecordTime:item.recordInterval.toString(),
                  })
            }
            else
            {
              _t.EditForm.setFieldsValue(
                  {
                    dataRecord:item.record.toString(),
                    dataAlarm:item.alarm.toString(),
                  })
            }
        _t.ShowRegisterLoading = false
      },1000)
      this.editVisible = true;
    },
    save() {
      this.EditForm.validateFields((err) => {
        if (!err) {
          let params = {
            uuid:this.editingKey,
            name:this.EditForm.getFieldValue('name'),
            auth:this.EditForm.getFieldValue('dataAuth'),
            type:this.EditForm.getFieldValue('dataType'),
            ByteOrder:this.EditForm.getFieldValue('ByteOrder'),
            unit:this.EditForm.getFieldValue('dataUnit'),
            conversionExpression:this.EditForm.getFieldValue('ConversionExpression'),
            alarm:parseInt(this.EditForm.getFieldValue('dataAlarm')),
            alarmLevel:parseInt(this.EditForm.getFieldValue('AlarmLevel')),
            AlarmMessage:this.EditForm.getFieldValue('AlarmMessage'),
            AlarmClearMessage:this.EditForm.getFieldValue('AlarmClearMessage'),
            record:parseInt(this.EditForm.getFieldValue('dataRecord')),
            RecordType:parseInt(this.EditForm.getFieldValue('dataRecordType')),
            recordInterval:parseInt(this.EditForm.getFieldValue('dataRecordTime')),
            RecordDataCharge:this.EditForm.getFieldValue('dataRecordChargeValue')?this.EditForm.getFieldValue('dataRecordChargeValue').toString():"",
          }
          let _t = this
          modbusModelRegisterEdit(params).then(function (res){
            if(res.data.code==200)
            {
              const newData = [..._t.dataSource];
              const target = newData.filter(item => _t.editingKey === item.uuid)[0];
              if (target) {
                target.name = _t.EditForm.getFieldValue('name')
                target.auth=_t.EditForm.getFieldValue('dataAuth')
                target.type=_t.EditForm.getFieldValue('dataType')
                target.ByteOrder=_t.EditForm.getFieldValue('ByteOrder')
                target.unit=_t.EditForm.getFieldValue('dataUnit')
                target.conversionExpression=_t.EditForm.getFieldValue('ConversionExpression'),
                target.alarm=parseInt(_t.EditForm.getFieldValue('dataAlarm'))
                target.alarmLevel=parseInt(_t.EditForm.getFieldValue('AlarmLevel')),
                target.AlarmMessage = _t.EditForm.getFieldValue('AlarmMessage')
                target.AlarmClearMessage = _t.EditForm.getFieldValue('AlarmClearMessage')
                target.record=parseInt(_t.EditForm.getFieldValue('dataRecord'))
                target.RecordType=parseInt(_t.EditForm.getFieldValue('dataRecordType'))
                target.recordInterval=parseInt(_t.EditForm.getFieldValue('dataRecordTime'))
                target.RecordDataCharge=_t.EditForm.getFieldValue('dataRecordChargeValue')?_t.EditForm.getFieldValue('dataRecordChargeValue').toString():""
                _t.dataSource = newData;
              }
              _t.$message.success(_t.$t("dataModel.saveSuccess"));
              _t.editVisible = false;
            }
            else
            {
              _t.$message.error(_t.$t("dataModel.saveFailed"));
            }
          })
        }
      })
    },
    addRegisterGroup(){
      this.RegisterForm.validateFields((err) => {
        if (!err) {

          const params = {
            muid:this.$route.params.uid,
            name:this.RegisterForm.getFieldValue('RegisterGroupName'),
            function:parseInt(this.RegisterForm.getFieldValue('RegisterFunction')),
            registerStart:parseInt(this.RegisterForm.getFieldValue('RegisterStartAddress')),
            registerCount:parseInt(this.RegisterForm.getFieldValue('RegisterCount')),
          }
          let numReg = /^[0-9]*$/
          let numRe = new RegExp(numReg)
          if((!numRe.test(params.registerCount))||(!numRe.test(params.registerStart)))
          {
            this.$message.error(this.$t("dataModel.modbusModel.MustBeNumber"));
            return;
          }
          if((params.registerCount>125)||(params.registerCount<=0)||(!numRe.test(params.registerCount))||(!numRe.test(params.registerStart)))
          {
            this.$message.error(this.$t("dataModel.modbusModel.RegisterQuantityError"));
            return;
          }
          this.RegisterVisible = false;
          let _t = this
          this.messageShowLoad = true
          modbusModelGroupAdd(params).then(function (res){
            _t.messageShowLoad = false
            if(res.data.code==2002)
            {
              _t.registerGroupList()
              _t.$message.success(_t.$t("dataModel.saveSuccess"));
            }
            else if(res.data.code==2001)
            {
              _t.$message.error(_t.$t("dataModel.modbusModel.RegisterExist"));
            }
            else
            {
              _t.$message.error(_t.$t("dataModel.saveFailed"));
            }
          }).catch(function (){
            _t.messageShowLoad = false
            _t.$message.error(_t.$t('loginPage.serverError'), 3)
          })
        }
      })
    },
    registerGroupList(){
      let _t = this
      this.messageShowLoad = true
      _t.registerGroupDataSource=[]
      let params = {
        muid:this.$route.params.uid
      }
      modbusModelGroupList(params).then(function (res){
        _t.messageShowLoad = false
        _t.registerGroupDataSource = res.data.list==null?[]:res.data.list
      })
    },
    deleteRegisterGroupRecord(uuid){
      let _t = this
      let parsms = {
        uuid:uuid
      }
      this.messageShowLoad = true
      modbusModelGroupDel(parsms).then(function (res){
        if(res.data.code==0)
        {
          _t.registerGroupDataSource = _t.registerGroupDataSource.filter(item => item.uuid !== uuid)
          _t.selectedRows = _t.selectedRows.filter(item => item.uuid !== uuid)
          _t.$message.success(_t.$t("dataModel.modbusModel.RegisterDeLSuccess"));
        }
        else if(res.data.code==2002)
        {
          _t.$message.success(_t.$t("dataModel.modbusModel.RegisterDeLFailed"));
        }
        _t.messageShowLoad = false
      })
    },
    registerConfig(uuid){
      this.registerGroupListTable = false
      this.registerAddressList(uuid)
    },
    registerAddressList(uuid){
      this.messageShowLoad = true
      let _t = this
      _t.dataSource=[]
      let params = {
        uuid:uuid
      }
      modbusModelRegisterList(params).then(function (res){
        _t.dataSource = res.data.list==null?[]:res.data.list
        _t.messageShowLoad = false
      })
    },
    deleteRecord() {
      let _t = this

      this.$confirm({
        content: _t.$t('dataModel.deleteConfirm'),
        okText: _t.$t('displayModel.ConfirmOk'),
        onOk() {
          if(_t.selectDataTableUuid.length==0)
          {
            return
          }
          let params={
            uuid:_t.selectDataTableUuid
          }
          modbusModelRegisterDel(params).then(function (res) {
            if(res.data.code==0)
            {
              _t.$message.success(_t.$t("dataModel.deleteSuccess"));
              for(let i=0;i<_t.selectDataTableUuid.length;i++)
              {
                _t.dataSource = _t.dataSource.filter(item => item.uuid !== _t.selectDataTableUuid[i])
                _t.selectedRows = _t.selectedRows.filter(item => item.uuid !== _t.selectDataTableUuid[i])
              }
              _t.selectDataTableUuid=[]
            }
            else {
              _t.$message.error(_t.$t("dataModel.deleteFailed"));
            }
          })
        },
        cancelText: _t.$t('displayModel.ConfirmCancel'),
        onCancel() {

        },
      });


    },
    onBlackCLK(){
      this.$router.push('/DeviceModel/ModbusModel')
    },
    alarmCharge(value){
      this.alarmStatus=parseInt(value)
      let _t = this
      this.$nextTick(function(){
        _t.EditForm.setFieldsValue(
            {
              AlarmLevel:"0",
            })
      });
    },
    recordCharge(value){
      this.recordStatus=parseInt(value)
    },
    chargeDataRecordType(value){
      this.DataRecordType = parseInt(value)
    }
  }
}
</script>


<style lang="less" >
.ant-table-tbody > tr > td {
  padding: 1px 1px;
  overflow-wrap: break-word;
}
 .ant-form-item {
   margin-bottom: 5px;

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
