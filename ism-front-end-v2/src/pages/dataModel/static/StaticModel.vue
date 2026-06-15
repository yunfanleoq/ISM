<template>
  <a-card>
      <a-space class="operator">
        <a-button @click="addVisible=true" type="primary" icon="plus">{{$t('dataModel.newModel')}}</a-button>
        <a-button @click="refresh()"  type="default" icon="sync" :loading="refIconLoading">{{$t("dataModel.refModel")}}</a-button>
      </a-space>

    <a-spin style="padding: 1px;"  :spinning="messageShowLoad" tip="Loading...">
      <a-table rowKey="Name" :pagination="pagination" :columns="columns" :data-source="dataSource">
      <template v-for="(item, index) in columns" :slot="item.slotName">
        <span :key="index">{{ $t(item.slotName) }}</span>
      </template>
        <template slot="serial" slot-scope="text, record,index, column">
           {{index+1}}
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
    <a-modal v-model="addVisible" :title="$t('configComponent.video.AddVideo')" @ok="onSubmit">
      <a-form :form="form" :label-col="{ span: 7 }" :wrapper-col="{ span: 15 }">
        <a-alert type="error" :closable="true" v-show="error" :message="error" showIcon style="margin-bottom: 24px;" />
        <a-form-item
            :label="$t('dataModel.static.DataName')"
        >
          <a-input  autocomplete="autocomplete"

                    v-decorator="['name', {rules: [{ required: true, message: $t('dataModel.static.DataName'), whitespace: true}]}]"
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
        >
         <span slot="label">
          {{$t('dataModel.static.DataDeviceType')}}&nbsp;
          <a-tooltip :title="$t('dataModel.static.DataDevicePublicTips')">
            <a-icon type="question-circle-o" />
          </a-tooltip>
      </span>
          <a-select  style="" autocomplete="autocomplete"

                     v-decorator="['DataDeviceType', {initialValue:'158',rules: [{ required: true, message: $t('dataModel.static.DataDeviceType'), whitespace: true}]}]"
          >
            <a-select-option  v-for="(device,index) in supportDeviceList" :key="index" :value=device.type.toString()>
              {{ $t(device.name) }}
            </a-select-option>
          </a-select>

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
            :label="$t('dataModel.DataAddCount')"
        >
          <a-input autocomplete="autocomplete"

                   v-decorator="['DataAddCount', {initialValue:'1',rules: [{ required: true, message: $t('dataModel.DataAddCount'), whitespace: true}]}]"
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

                    v-decorator="['name', {rules: [{ required: true, message: $t('dataModel.static.DataName'), whitespace: true}]}]"
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
        >
         <span slot="label">
          {{$t('dataModel.static.DataDeviceType')}}&nbsp;
          <a-tooltip :title="$t('dataModel.static.DataDevicePublicTips')">
            <a-icon type="question-circle-o" />
          </a-tooltip>
      </span>
          <a-select  style="" autocomplete="autocomplete" disabled

                     v-decorator="['DataDeviceType', {initialValue:'158',rules: [{ required: true, message: $t('dataModel.static.DataDeviceType'), whitespace: true}]}]"
          >
            <a-select-option  v-for="(device,index) in supportDeviceList" :key="index" :value=device.type.toString()>
              {{ $t(device.name) }}
            </a-select-option>
          </a-select>

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
  </a-card>
</template>

<script>
import {GetStaticModelList,StaticModelEdit, StaticModelDel,StaticModelAdd} from "../../../services/staticmodel";
import {getSupportDeviceList} from "../../../services/device";
import Mtextarea from '@/components/textarea/index'

export default {
  name: 'StaticModel',
  i18n: require('../../../i18n/language'),
  data () {
    return {
      pagination:{
        pageSize:15,
        showSizeChanger:true
      },
      searchText: '',
      searchInput: null,
      searchedColumn: '',
      DataRecordType:0,
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
          scopedSlots: {filterDropdown: 'filterDropdown', filterIcon: 'filterIcon', customRender: 'Name', title: 'dataModel.static.DataName' },
          dataIndex: 'Name',
          onFilter: (value, record) =>
              record.Name
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
          name:this.$t('dataModel.static.DataTypeInt'),
          value:1
        },
        {
          name:this.$t('dataModel.static.DataTypeString'),
          value:2
        },
        {
          name:this.$t('dataModel.static.DataTypeDouble'),
          value:3
        },
        {
          name:this.$t('dataModel.static.DataTypeJson'),
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
  components: {Mtextarea},
  authorize: {
    // deleteRecord: 'delete'
  },
  mounted(){
    let _t = this
    this.getSupportDevice()
  },
  activated(){

  },
  created(){

  },
  watch: {
    '$route' () {
     this.dataSource=[]
     this.getModelList()
    }
  },
  methods: {
    handleSearch(selectedKeys, confirm, dataIndex) {
      confirm();
      this.searchText = selectedKeys[0];
      this.searchedColumn = dataIndex;
    },
    handleReset(clearFilters) {
      clearFilters();
      this.searchText = '';
    },
    recordCharge(value){
      this.recordStatus=parseInt(value)
    },
    chargeDataRecordType(value){
      this.DataRecordType = parseInt(value)
    },
    alarmCharge(value){
      this.alarmStatus=parseInt(value)
    },
    getSupportDevice(){
      let _t = this
      getSupportDeviceList().then(function (res){
        let publicDevice = {
          name: "dataModel.static.DataDevicePublic",
          type: '158'
        }
        for(let i=0;i<res.data.list.length;i++)
        {
          if(res.data.list[i].type!=7&&res.data.list[i].type!=6)
          {
            _t.supportDeviceList.push(res.data.list[i])
          }
        }
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
      GetStaticModelList(params).then(function (res){
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
            tableData.DataDeviceType = res.data.list[i].DataDeviceType.toString()
            tableData.DataUnit = res.data.list[i].DataUnit
            tableData.DataDefaultValue = res.data.list[i].DataDefaultValue
            tableData.DataDescription = res.data.list[i].DataDescription

            tableData.dataAlarm = res.data.list[i].dataAlarm
            tableData.dataRecord = res.data.list[i].dataRecord
            tableData.RecordType = res.data.list[i].RecordType
            tableData.AlarmLevel = res.data.list[i].AlarmLevel
            tableData.AlarmMessage = res.data.list[i].AlarmMessage
            tableData.AlarmClearMessage = res.data.list[i].AlarmClearMessage
            tableData.recordInterval = res.data.list[i].recordInterval
            tableData.RecordDataCharge = res.data.list[i].RecordDataCharge

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
      StaticModelDel(params).then(function (res) {
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
    onSubmit (e) {
      e.preventDefault()
      this.form.validateFields((err) => {
        if (!err) {
          this.logging = true
          const params = {
            name:this.form.getFieldValue('name'),
            DataDeviceType:parseInt(this.form.getFieldValue('DataDeviceType')),
            DataType:parseInt(this.form.getFieldValue('DataType')),
            DataDefaultValue:this.form.getFieldValue('DataDefaultValue'),
            DataUnit:this.form.getFieldValue('DataUnit'),
            Count:parseInt(this.form.getFieldValue('DataAddCount')),
            DataDescription:this.form.getFieldValue('description'),
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
            params.RecordInterval=this.form.getFieldValue('dataRecordTime')?parseInt(this.form.getFieldValue('dataRecordTime')):0
            params.RecordType=this.form.getFieldValue('dataRecordType')?parseInt(this.form.getFieldValue('dataRecordType')):0
            params.RecordDataCharge=this.form.getFieldValue('dataRecordChargeValue')?this.form.getFieldValue('dataRecordChargeValue').toString():"0"
          }else{
            params.dataRecord=0
            params.dataAlarm = 0
          }
          StaticModelAdd(params).then(this.addResponse)
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
              name: this.editForm.getFieldValue('name'),
              DataDeviceType: parseInt(this.editForm.getFieldValue('DataDeviceType')),
              DataType: parseInt(this.editForm.getFieldValue('DataType')),
              DataDefaultValue: this.editForm.getFieldValue('DataDefaultValue'),
              DataUnit: this.editForm.getFieldValue('DataUnit'),
              DataDescription: this.editForm.getFieldValue('description'),
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
            params.RecordInterval=this.editForm.getFieldValue('dataRecordTime')?parseInt(this.editForm.getFieldValue('dataRecordTime')):0
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
          StaticModelEdit(PassParams).then(function (res){
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
      this.alarmStatus = item.dataAlarm
      this.recordStatus = item.dataRecord
      this.DataRecordType = item.RecordType?item.RecordType:0

      this.$message.loading(this.$t("monitor.loading"), 0.2)
      setTimeout(function (){
        _t.textAreValue = item.DataDescription
        _t.editForm.setFieldsValue(
            {
              name:item.Name,
              DataDeviceType:item.DataDeviceType.toString(),
              DataType:item.DataType.toString(),
              DataDefaultValue:item.DataDefaultValue,
              DataUnit:item.DataUnit,
              description:item.DataDescription,
              dataAlarm:item.dataAlarm.toString(),
              AlarmLevel:item.AlarmLevel.toString(),
              AlarmMessage:item.AlarmMessage,
              AlarmClearMessage:item.AlarmClearMessage,
              dataRecord:item.dataRecord.toString(),
              dataRecordTime:item.recordInterval,
            })
        if (item.dataAlarm==1){
          _t.editForm.setFieldsValue(
              {
                dataAlarm:item.dataAlarm.toString(),
                AlarmLevel:item.AlarmLevel.toString(),
                dataRecord:item.dataRecord.toString(),
                AlarmMessage :item.AlarmMessage,
                AlarmClearMessage : item.AlarmClearMessage,
              })
        }
        else  if (item.dataRecord==1)
        {
          _t.editForm.setFieldsValue(
              {
                dataAlarm:item.dataAlarm.toString(),
                dataRecord:item.dataRecord.toString(),
                dataRecordType:item.RecordType.toString(),
                dataRecordChargeValue:item.RecordDataCharge.toString(),
                dataRecordTime:item.recordInterval.toString(),
              })
        }
        else
        {
          _t.editForm.setFieldsValue(
              {
                dataRecord:item.dataRecord.toString(),
                dataAlarm:item.dataAlarm.toString(),
              })
        }
      },500)
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
