<template>
  <a-card>
      <a-space class="operator">
        <a-button @click="addVisible=true" type="primary" icon="plus">{{$t('dataModel.newModel')}}</a-button>
        <a-button @click="refresh()"  type="default" icon="sync" :loading="refIconLoading">{{$t("dataModel.refModel")}}</a-button>
      </a-space>

    <a-spin style="padding: 1px;"  :spinning="messageShowLoad" tip="Loading...">
      <a-table rowKey="no" :pagination="pagination" :columns="columns" :data-source="dataSource">
      <template v-for="(item, index) in columns" :slot="item.slotName">
        <span :key="index">{{ $t(item.slotName) }}</span>
      </template>
        <div
            slot="filterDropdown"
            slot-scope="{ setSelectedKeys, selectedKeys, confirm, clearFilters, column }"
            style="padding: 8px"
        >
          <a-input
              v-ant-ref="c => (searchInput = c)"
              :placeholder="`Search ${column.dataIndex}`"
              :value="selectedKeys[0]"
              style="width: 188px; margin-bottom: 8px; display: block;"
              @change="e => setSelectedKeys(e.target.value ? [e.target.value] : [])"
              @pressEnter="() => handleSearch(selectedKeys, confirm, column.dataIndex)"
          />
          <a-button
              type="primary"
              icon="search"
              size="small"
              style="width: 90px; margin-right: 8px"
              @click="() => handleSearch(selectedKeys, confirm, column.dataIndex)"
          >

            {{$t('readData.Search')}}
          </a-button>
          <a-button size="small" style="width: 90px" @click="() => handleReset(clearFilters)">

            {{$t('readData.Reset')}}
          </a-button>
        </div>
        <a-icon
            slot="filterIcon"
            slot-scope="filtered"
            type="search"
            :style="{ color: filtered ? '#108ee9' : undefined }"
        />
        <template slot="Name" slot-scope="text, record,index, column">
            <span v-if="searchText && searchedColumn === column.dataIndex">
                      <template
                          v-for="(fragment, i) in text
                          .toString()
                          .split(new RegExp(`(?<=${searchText})|(?=${searchText})`, 'i'))"
                      >
                       <mark
                           v-if="fragment.toLowerCase() === searchText.toLowerCase()"
                           :key="i"
                           class="highlight"
                       >{{ $t(fragment) }}</mark>
                        <template v-else>{{ $t(fragment) }}</template>
                      </template>
                    </span>
          <template v-else>
            {{ $t(text) }}
          </template>
        </template>
      <div slot="DataType" slot-scope="text" >
          <template
              v-for="(fragment, i) in DataTypeList"
          >
            <span
                v-if="fragment.value === text"
                :key="i"
                class="highlight"
            >
              {{ fragment.name }}
            </span>
          </template>
      </div>
      <div slot="DataDeviceType" slot-scope="text">
        <template
            v-for="(fragment, i) in supportDeviceList"
        >
          <span
              v-if="fragment.type == text"
              :key="i"
              class="highlight"
          >
            {{ $t(fragment.name) }}
          </span>
        </template>
      </div>
      <div slot="action" slot-scope="text, record">
        <a @click="GoToEdit(record)" style="color: #13C2C2"><a-icon type="edit" />{{$t('dataModel.modelDetail')}}</a> |
        <a-popconfirm :title="$t('dataModel.deleteConfirm')" @confirm="deleteRecord(record.key)">
          <a-icon slot="icon" type="question-circle-o" style="color: red" />
          <a-icon type="delete" theme="twoTone" two-tone-color="#eb2f96"/><a style="color: #eb2f96">{{$t('dataModel.delete')}}</a>
        </a-popconfirm>
    </div>
    </a-table>
    </a-spin>
    <a-modal v-drag-modal v-model="addVisible" :title="$t('configComponent.video.AddVideo')" @ok="addCustomData">
      <a-form :form="form" :label-col="{ span: 7 }" :wrapper-col="{ span: 15 }">
        <a-alert type="error" :closable="true" v-show="error" :message="error" showIcon style="margin-bottom: 24px;" />
        <a-form-item
            :label="$t('dataModel.static.DataName')"
        >
          <a-input  autocomplete="autocomplete"

                    v-decorator="['name', {rules: [{ required: true, validator: isValidateTxtNonSpec,message: $t('device.deviceNameVal'), whitespace: true}]}]"
          />
        </a-form-item>
        <a-form-item
            :label="$t('dataModel.static.DataType')"
        >
          <a-select  style="" autocomplete="autocomplete"

                     v-decorator="['DataType', {initialValue:'1',rules: [{ required: true, message: $t('dataModel.static.DataType'), whitespace: true}]}]"
          >
            <a-select-option v-for="options in DataTypeList" :key="options.value" :value="options.value.toString()">
              {{ $t(options.name) }}
            </a-select-option>
          </a-select>

        </a-form-item>
        <a-form-item
            :label="$t('dataModel.CustomData.DataFrom')"
        >
          <a-input    v-decorator="['DataFrom', {rules: [{ required: true, message: $t('dataModel.static.DataType'), whitespace: true}]}]">
            <a-tooltip placement="top" slot="addonAfter">
              <template slot="title">
                <span>{{$t('displayConfig.Properties.SelectValue')}}</span>
              </template>
              <icon-font @click="ShowDeviceDataModel()" type="icon-xuanzeshuju"  />
            </a-tooltip>
          </a-input>
        </a-form-item>
        <a-form-item :label="$t('dataModel.modbusModel.ConversionExpression')">
          <a-tooltip placement="top">
            <template slot="title">
              <span>{{$t('dataModel.CustomData.ExpressTips')}}</span>
            </template>
            <a-input   autocomplete="autocomplete"   v-decorator="[
                  'ConversionExpression',
                  {
                    rules: [{ required: true, message: $t('dataModel.modbusModel.ConversionExpression') }],
                  },
                ]">
            </a-input>
          </a-tooltip>
        </a-form-item>
        <a-form-item
            :label="$t('dataModel.static.DataDefaultValue')"
        >
          <a-input autocomplete="autocomplete"

                   v-decorator="['DataDefaultValue', {rules: [{ required: true, message: $t('dataModel.static.DataDefaultValue'), whitespace: true,initialValue:162}]}]"
          />
        </a-form-item>

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

        <div v-if="alarmStatus==1">
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
          <a-form-item :label="$t('dataModel.editData.AlarmMessage')">
            <a-input   autocomplete="autocomplete"   v-decorator="[
                  'AlarmMessage',
                  {
                    rules: [{ required: true, message: $t('dataModel.editData.AlarmMessage') }],
                  },
                ]">
            </a-input>
          </a-form-item>
          <a-form-item :label="$t('dataModel.editData.AlarmClearMessage')">
            <a-input   autocomplete="autocomplete"   v-decorator="[
                  'AlarmClearMessage',
                  {
                    rules: [{ required: false, message: $t('dataModel.editData.AlarmClearMessage') }],
                  },
                ]">
            </a-input>
          </a-form-item>
        </div>
        <!--存储            -->
        <div v-if="alarmStatus==0">
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
          <div  v-if="recordStatus" >
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
          </div>
          <div :span="12" v-if="(recordStatus)&&((DataRecordType==0)||(DataRecordType==3))">
            <a-form-item :label="$t('dataModel.dataRecordChargeValue')">
              <a-input   :min="1"  style="width: 100%" autocomplete="autocomplete"  v-decorator="[
                  'dataRecordChargeValue',
                  {
                    rules: [{ required: true, message: $t('dataModel.dataRecordChargeValue') }],
                  },
                ]">
              </a-input>
            </a-form-item>
          </div>
          <div  v-if="(recordStatus)&&(DataRecordType==1)">
            <a-form-item :label="$t('dataModel.editData.dataRecordTime')">
              <a-input-number   :min="1"  style="width: 100%" autocomplete="autocomplete"  v-decorator="[
                  'dataRecordTime',
                  {
                    rules: [{ required: true, message: $t('dataModel.editData.dataRecordTime') }],
                  },
                ]">
              </a-input-number>
            </a-form-item>
          </div>
        </div>
        <a-form-item
            :label="$t('dataModel.static.DataUnit')"
        >
          <a-input autocomplete="autocomplete"

                   v-decorator="['DataUnit', {rules: [{ required: false, message: $t('dataModel.static.DataUnit'), whitespace: true}]}]"
          />
        </a-form-item>
        <a-form-item
            :label="$t('dataModel.static.DataDec')"
        >
          <Mtextarea   v-model="textAreValue"
                       rows="2"
                       :showWordLimit="true"
                       :maxLength="100"
                       :autoSize="false"
                       v-decorator="['description', { rules: [{ required: true, message: $t('dataModel.static.DataDec') }] }]"
          />
        </a-form-item>
      </a-form>
    </a-modal>
    <a-modal v-model="editVisible" :title="$t('dataModel.static.EditTitle')" @ok="onEditSubmit">
      <a-form :form="editForm" :label-col="{ span: 7 }" :wrapper-col="{ span: 15 }">
        <a-alert type="error" :closable="true" v-show="error" :message="error" showIcon style="margin-bottom: 24px;" />
        <a-form-item
            :label="$t('dataModel.static.DataName')"
        >
          <a-input  autocomplete="autocomplete"

                    v-decorator="['name', {rules: [{ required: true, validator: isValidateTxtNonSpec,message: $t('device.deviceNameVal'), whitespace: true}]}]"
          />
        </a-form-item>
        <a-form-item
            :label="$t('dataModel.static.DataType')"
        >
          <a-select  style="" autocomplete="autocomplete"

                     v-decorator="['DataType', {initialValue:'1',rules: [{ required: true, message: $t('dataModel.static.DataType'), whitespace: true}]}]"
          >
            <a-select-option v-for="options in DataTypeList" :key="options.value" :value="options.value.toString()">
              {{ $t(options.name) }}
            </a-select-option>
          </a-select>

        </a-form-item>
        <a-form-item
            :label="$t('dataModel.CustomData.DataFrom')"
        >
          <a-input   disabled v-decorator="['DataFrom', {rules: [{ required: true, message: $t('dataModel.static.DataType'), whitespace: true}]}]">
            <a-tooltip placement="top" slot="addonAfter">
              <template slot="title">
                <span>{{$t('displayConfig.Properties.SelectValue')}}</span>
              </template>
              <icon-font @click="ShowDeviceDataModel()" type="icon-xuanzeshuju"  />
            </a-tooltip>
          </a-input>
        </a-form-item>
        <a-form-item :label="$t('dataModel.modbusModel.ConversionExpression')">
          <a-tooltip placement="top">
            <template slot="title">
              <span>{{$t('dataModel.CustomData.ExpressTips')}}</span>
            </template>
            <a-input   autocomplete="autocomplete"   v-decorator="[
                  'ConversionExpression',
                  {
                    rules: [{ required: true, message: $t('dataModel.modbusModel.ConversionExpression') }],
                  },
                ]">
            </a-input>
          </a-tooltip>
        </a-form-item>
        <a-form-item
            :label="$t('dataModel.static.DataDefaultValue')"
        >
          <a-input autocomplete="autocomplete"

                   v-decorator="['DataDefaultValue', {rules: [{ required: true, message: $t('dataModel.static.DataDefaultValue'), whitespace: true,initialValue:162}]}]"
          />
        </a-form-item>

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

        <div v-if="alarmStatus==1">
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
          <a-form-item :label="$t('dataModel.editData.AlarmMessage')">
            <a-input   autocomplete="autocomplete"   v-decorator="[
                  'AlarmMessage',
                  {
                    rules: [{ required: true, message: $t('dataModel.editData.AlarmMessage') }],
                  },
                ]">
            </a-input>
          </a-form-item>
          <a-form-item :label="$t('dataModel.editData.AlarmClearMessage')">
            <a-input   autocomplete="autocomplete"   v-decorator="[
                  'AlarmClearMessage',
                  {
                    rules: [{ required: false, message: $t('dataModel.editData.AlarmClearMessage') }],
                  },
                ]">
            </a-input>
          </a-form-item>
        </div>
        <!--存储            -->
        <div v-if="alarmStatus==0">
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
          <div  v-if="recordStatus" >
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
          </div>
          <div v-if="(recordStatus)&&((DataRecordType==0)||(DataRecordType==3))">
            <a-form-item :label="DataRecordType==0?$t('dataModel.dataRecordChargeValue'):$t('dataModel.dataRecordChangeRateValue')">
              <a-input   :min="1"  style="width: 100%" autocomplete="autocomplete"  v-decorator="[
                  'dataRecordChargeValue',
                  {
                    rules: [{ required: true, message: $t('dataModel.dataRecordChargeValue') }],
                  },
                ]">
              </a-input>
            </a-form-item>
          </div>
          <div  v-if="(recordStatus)&&(DataRecordType==1)">
            <a-form-item :label="$t('dataModel.editData.dataRecordTime')">
              <a-input-number   :min="1"  style="width: 100%" autocomplete="autocomplete"  v-decorator="[
                  'dataRecordTime',
                  {
                    rules: [{ required: true, message: $t('dataModel.editData.dataRecordTime') }],
                  },
                ]">
              </a-input-number>
            </a-form-item>
          </div>
        </div>
        <a-form-item
            :label="$t('dataModel.static.DataUnit')"
        >
          <a-input autocomplete="autocomplete"

                   v-decorator="['DataUnit', {rules: [{ required: false, message: $t('dataModel.static.DataUnit'), whitespace: true}]}]"
          />
        </a-form-item>
        <a-form-item
            :label="$t('dataModel.static.DataDec')"
        >
          <Mtextarea   v-model="textAreValue"
                       rows="4"
                       :showWordLimit="true"
                       :maxLength="100"
                       :autoSize="false"
                       v-decorator="['description', { rules: [{ required: true, message: $t('dataModel.static.DataDec') }] }]"
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <device-data-model @onSelectDataModel="onSelectData" ref="deviceDataModel"></device-data-model>
  </a-card>
</template>

<script>
import {CustomDataAdd,CustomDataEdit, CustomDataDel,CustomDataList} from "@/services/customDataModel";
import {getSupportDeviceList} from "@/services/device";
import Mtextarea from '@/components/textarea/index'
import deviceDataModel from "@/components/deviceDataModel/deviceDataModel";

export default {
  name: 'CustomDataModel',
  i18n: require('../../../i18n/language'),
  data () {
    return {
      pagination:{
        pageSize:15,
        showSizeChanger:true
      },
      DataRecordType:0,
      DataFromInfo:{
        DeviceSN: "",
        IsDevice: false,
        name: "通断电",
        uuid: "",
      },
      searchText: '',
      searchInput: null,
      searchedColumn: '',
      alarmStatus:2,
      recordStatus:0,
      messageShowLoad:false,
      advanced: true,
      refIconLoading: false,
      columns: [
        {
          width: '10%',
          slotName: 'dataModel.modelTableIndex',
          scopedSlots: { customRender: 'serial', title: 'dataModel.modelTableIndex' },
          dataIndex: 'no'
        },
        {
          width: '20%',
          slotName: 'dataModel.static.DataName',
          scopedSlots: {  filterDropdown: 'filterDropdown', filterIcon: 'filterIcon', customRender: 'Name', title: 'dataModel.static.DataName' },
          dataIndex: 'Name',
          onFilter: (value, record) =>
              record.name
                  .toString()
                  .toLowerCase()
                  .includes(value.toLowerCase()),
          onFilterDropdownVisibleChange: visible => {
            if (visible) {
              setTimeout(() => {
                this.searchInput.focus();
              }, 0);
            }
          },
        },
        {
          slotName: 'dataModel.static.DataDeviceType',
          width: '10%',
          scopedSlots: { customRender: 'DataDeviceType', title: 'dataModel.static.DataDeviceType' },
          dataIndex: 'DataDeviceType',
        },
        {
          width: '10%',
          slotName: 'dataModel.static.DataType',
          scopedSlots: { customRender: 'DataType', title: 'dataModel.static.DataType' },
          dataIndex: 'DataType',
        },
        {
          width: '8%',
          slotName: 'dataModel.static.DataUnit',
          scopedSlots: { customRender: 'DataUnit', title: 'dataModel.static.DataUnit' },
          dataIndex: 'DataUnit',
        },
        {
          width: '15%',
          slotName: 'dataModel.static.DataDefaultValue',
          scopedSlots: { customRender: 'DataDefaultValue', title: 'dataModel.static.DataDefaultValue' },
          dataIndex: 'DataDefaultValue',
        },
        {
          width: '15%',
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
      selectedRows: [],
      addVisible:false,
      error: '',
      editUuid:"",
      editVisible:false,
      form: this.$form.createForm(this),
      editForm: this.$form.createForm(this),
      textAreValue:"",
      value: 1
    }
  },
  components: {
    Mtextarea,
    deviceDataModel,
  },
  authorize: {
    // deleteRecord: 'delete'
  },
  mounted(){
    let _t = this
  },
  activated(){

  },
  created(){
    this.getSupportDevice()
  },
  watch: {
    '$route' () {
     this.dataSource=[]
     this.getModelList()
    }
  },
  methods: {
    isSpec(s) {
      let pattern = /[~!@#$%^&*<>|'-]/gi
      return pattern.test(s)
    },
    isValidateTxtNonSpec (rule, value, callback) {
      if (value != null && value !== '') {
        let numStr = value.charAt(0);
        if ((this.isSpec(value)) || (value.indexOf(' ') !== -1)||(!isNaN(parseFloat(numStr)) && isFinite(numStr))) {
          callback(new Error('不能包含特殊字符或空格'))
        } else {
          callback()
        }
      } else {
        callback()
      }
    },
    handleSearch(selectedKeys, confirm, dataIndex) {
      confirm();
      this.searchText = selectedKeys[0];
      this.searchedColumn = dataIndex;
    },
    handleReset(clearFilters) {
      clearFilters();
      this.searchText = '';
    },
    chargeDataRecordType(value){
      this.DataRecordType = parseInt(value)
    },
    alarmCharge(value){
      this.alarmStatus=parseInt(value)
    },
    recordCharge(value){
      this.recordStatus=parseInt(value)
    },
    getSupportDevice(){
      let _t = this
      getSupportDeviceList().then(function (res){
        let publicDevice = {
          name: "dataModel.static.DataDevicePublic",
          type: '158'
        }
        _t.supportDeviceList =res.data.list
        _t.supportDeviceList.push(publicDevice)

        _t.dataSource=[]
        _t.getModelList()
      })
    },
    refresh(){
      this.refIconLoading=true
      this.getModelList()
    },
    getModelList(){
      this.dataSource=[]
      let _t = this
      const  params= {
        DisplayType: 1
      }
      this.messageShowLoad=true
      CustomDataList(params).then(function (res){
        let tableData={}
        _t.refIconLoading=false
        _t.messageShowLoad=false
        if(res.data.list!=null)
        {
          for(let i=0;i<res.data.list.length;i++)
          {
            tableData.key = res.data.list[i].uuid
            tableData.no = res.data.list[i].ID
            tableData.Name = res.data.list[i].name
            tableData.DataType = res.data.list[i].DataType
            tableData.DataDeviceType = res.data.list[i].DeviceType.toString()
            tableData.DataUnit = res.data.list[i].unit
            tableData.DataDefaultValue = res.data.list[i].DataDefaultValue
            tableData.DataDescription = res.data.list[i].DataDescription
            tableData.tag = res.data.list[i]
            _t.dataSource.push(tableData)
            tableData={}
          }
        }

    })
    },
    deleteRecord(key) {
      let params={
        uuid:key
      }
      let _t = this
      CustomDataDel(params).then(function (res) {
        if(res.data.code==0)
        {
          _t.dataSource = _t.dataSource.filter(item => item.key !== key)
          _t.selectedRows = _t.selectedRows.filter(item => item.key !== key)
        }
        else
        {
            _t.$message.error(_t.$t("dataModel.modelBand"))
        }
      })
    },
    addCustomData (e) {
      e.preventDefault()
      this.form.validateFields((err) => {
        if (!err) {
          this.logging = true
          const params = {
            Name:this.form.getFieldValue('name'),
            DataType:parseInt(this.form.getFieldValue('DataType')),
            DataDefaultValue:this.form.getFieldValue('DataDefaultValue'),
            unit:this.form.getFieldValue('DataUnit'),
            DataFromName:this.form.getFieldValue('DataFrom'),
            conversionExpression:this.form.getFieldValue('ConversionExpression'),
            DataDescription:this.form.getFieldValue('description'),
            DeviceSN:this.DataFromInfo.DeviceSN,
            IsDevice:this.DataFromInfo.IsDevice?1:0,
            DataUuid:this.DataFromInfo.DataUuid,
            DeviceType : this.DataFromInfo.DeviceType,
            SelectDataModelUUid : this.DataFromInfo.selectDataModelUUid
          };
          if (this.alarmStatus==1){
            params.dataAlarm = parseInt(this.form.getFieldValue('dataAlarm'))
            params.AlarmLevel=parseInt(this.form.getFieldValue('AlarmLevel'))
            params.dataRecord=0
            params.AlarmMessage=this.form.getFieldValue('AlarmMessage')
            params.AlarmClearMessage =this.form.getFieldValue('AlarmClearMessage')
          }
          else  if (this.recordStatus==1)
          {
            params.dataAlarm = 0
            params.DataRecord = parseInt(this.form.getFieldValue('dataRecord'))
            params.recordInterval=this.form.getFieldValue('dataRecordTime')?parseInt(this.form.getFieldValue('dataRecordTime')):0
            params.RecordType=this.form.getFieldValue('dataRecordType')?parseInt(this.form.getFieldValue('dataRecordType')):0
            params.RecordDataCharge=this.form.getFieldValue('dataRecordChargeValue')?this.form.getFieldValue('dataRecordChargeValue'):"0"
          }else{
            params.dataRecord=0
            params.dataAlarm = 0
          }
          CustomDataAdd(params).then(this.addResponse)
        }
      })
    },
    onEditSubmit (e) {
      e.preventDefault()
      let _t = this
      this.editForm.validateFields((err) => {
        if (!err) {
          this.logging = true
          const params = {
            Name:this.editForm.getFieldValue('name'),
            DataType:parseInt(this.editForm.getFieldValue('DataType')),
            DataDefaultValue:this.editForm.getFieldValue('DataDefaultValue'),
            unit:this.editForm.getFieldValue('DataUnit'),
            DataFromName:this.editForm.getFieldValue('DataFrom'),
            conversionExpression:this.editForm.getFieldValue('ConversionExpression'),
            DataDescription:this.editForm.getFieldValue('description'),
          };
          if (this.alarmStatus==1){
            params.dataAlarm = parseInt(this.editForm.getFieldValue('dataAlarm'))
            params.AlarmLevel=parseInt(this.editForm.getFieldValue('AlarmLevel'))
            params.dataRecord=0
            params.AlarmMessage=this.editForm.getFieldValue('AlarmMessage')
            params.AlarmClearMessage =this.editForm.getFieldValue('AlarmClearMessage')
          }
          else  if (this.recordStatus==1)
          {
            params.dataAlarm = 0
            params.dataRecord = parseInt(this.editForm.getFieldValue('dataRecord'))
            params.recordInterval=this.editForm.getFieldValue('dataRecordTime')?parseInt(this.editForm.getFieldValue('dataRecordTime')):0
            params.RecordType=this.editForm.getFieldValue('dataRecordType')?parseInt(this.editForm.getFieldValue('dataRecordType')):0
            params.RecordDataCharge=this.editForm.getFieldValue('dataRecordChargeValue')?this.editForm.getFieldValue('dataRecordChargeValue').toString():""
          }
          else
          {
            params.dataRecord=0
            params.dataAlarm = 0
          }
          const PassParams = {
            uuid:_t.editUuid,
            Data: params
          };
          CustomDataEdit(PassParams).then(function (res){
            _t.logging = false
            if (res.data.code == 0) {
              _t.$message.success(_t.$t('dataModel.static.EditSuccess'), 2)
              _t.getModelList()
              _t.editVisible = false
            }
            else if (res.data.code == 3001)
            {
              _t.$message.error(_t.$t('dataModel.modelNameRepeat'), 2)
            }
            else {
              _t.$message.error(_t.$t('dataModel.static.EditFailed'), 2)
            }
          })
        }
      })
    },
    addResponse(res) {
      this.logging = false
      if (res.data.code == 0) {
        this.$message.success(this.$t('dataModel.modelAddSuccess'), 3)
        this.getModelList()
        this.addVisible = false
      }
      else if (res.data.code == 3001)
      {
        this.$message.error(this.$t('dataModel.modelNameRepeat'), 3)
      }
      else {
        this.$message.error(this.$t('dataModel.modelAddFailed'), 3)
      }
    },
    handleMenuClick (e) {
      if (e.key === 'delete') {
        this.remove()
      }
    },
    GoToEdit(item){
      this.editVisible = true
      this.editUuid = item.key
      let _t = this
      this.$message.loading(this.$t("monitor.loading"), 0.2)
      this.alarmStatus = item.tag.dataAlarm
      this.recordStatus = item.tag.dataRecord
      this.DataRecordType = item.tag.RecordType?item.tag.RecordType:0
      setTimeout(function (){
        _t.textAreValue = item.DataDescription
        _t.editForm.setFieldsValue(
            {
              name:item.Name,
              DataFrom:item.tag.DataFromName,
              DataType:item.DataType.toString(),
              DataDefaultValue:item.DataDefaultValue,
              DataUnit:item.DataUnit,
              ConversionExpression:item.tag.conversionExpression,
              dataAlarm:item.tag.dataAlarm.toString(),
              AlarmLevel:item.tag.AlarmLevel.toString(),
              AlarmMessage:item.tag.AlarmMessage,
              AlarmClearMessage:item.tag.AlarmClearMessage,
              dataRecord:item.tag.dataRecord.toString(),
              dataRecordTime:item.tag.recordInterval,
              description:item.tag.DataDescription,
            })

            if (item.tag.dataAlarm==1){
              _t.editForm.setFieldsValue(
                  {
                    dataAlarm:item.tag.dataAlarm.toString(),
                    AlarmLevel:item.tag.AlarmLevel.toString(),
                    dataRecord:item.tag.dataRecord.toString(),
                    AlarmMessage :item.tag.AlarmMessage,
                    AlarmClearMessage : item.tag.AlarmClearMessage,
                  })
            }
            else  if (item.tag.dataRecord==1)
            {
              _t.editForm.setFieldsValue(
                  {
                    dataAlarm:item.tag.dataAlarm.toString(),
                    dataRecord:item.tag.dataRecord.toString(),
                    dataRecordType:item.tag.RecordType.toString(),
                    dataRecordChargeValue:item.tag.RecordDataCharge.toString(),
                    dataRecordTime:item.tag.recordInterval.toString(),
                  })
            }
            else
            {
              _t.editForm.setFieldsValue(
                  {
                    dataRecord:item.tag.dataRecord.toString(),
                    dataAlarm:item.tag.dataAlarm.toString(),
                  })
            }
      },500)
    },
    onSelectData(selectData) {
      this.DataFromInfo.DeviceSN=selectData.DeviceSN
      this.DataFromInfo.IsDevice=selectData.IsDevice
      this.DataFromInfo.DataUuid=selectData.uuid
      this.DataFromInfo.DeviceType = selectData.DeviceType
      this.DataFromInfo.selectDataModelUUid = selectData.selectDataModelUUid

      if(selectData.DeviceType==-1)
      {
        this.$message.error(this.$t('dataModel.CustomData.DataErrorTips'), 3)
        return
      }
      this.form.setFieldsValue(
          {
            DataFrom:selectData.name,
          })
    },
    ShowDeviceDataModel(index,type){
      this.$refs.deviceDataModel.showDataModal()
    },
  }
}
</script>

<style lang="less" scoped>
::v-deep .search{
  margin-bottom: 54px;
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
