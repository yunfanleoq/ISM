<template>
  <a-card v-if="defaultTabKey==1" >
    <a-spin style="padding: 1px;"  :spinning="messageShowLoad" tip="Loading...">
      <a-form @submit="onSubmit" :form="form">
      <a-alert type="error" :closable="true" v-show="error" :message="error" showIcon style="margin-bottom: 24px;" />
      <a-form-item
          :label="$t('dataModel.modelName')"
          :labelCol="{span: 7}"
          :wrapperCol="{span: 10}"
      >
        <a-input  autocomplete="autocomplete"

                  v-decorator="['name', {rules: [{ required: true, message: $t('dataModel.modelName'), whitespace: true}]}]"
        />
      </a-form-item>
      <a-form-item
          :label="$t('dataModel.modelDec')"
          :labelCol="{span: 7}"
          :wrapperCol="{span: 10}"
      >
        <a-textarea autocomplete="autocomplete"

                    v-decorator="['dec', {rules: [{ required: true, message: $t('dataModel.modelDec'), whitespace: true}]}]"
        />
      </a-form-item>
      <a-form-item
          :label="$t('device.deviceConfigurationModelName')"
          :labelCol="{span: 7}"
          :wrapperCol="{span: 10}"
      >
        <a-select  @select="GetDisplayPage"
                   v-decorator="[
                'configurationModel',
                {
                  rules: [{ required: false, message: $t('device.deviceConfigurationModelName') }],
                },
              ]"
        >
          <a-select-option v-for="(model,index) in configurationModel" :key="index" :value="model.uuid">
            {{ model.name }}
          </a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item
            :label="$t('device.deviceConfigurationPageName')"
            :labelCol="{span: 7}"
            :wrapperCol="{span: 10}"
        >
          <a-select
              v-decorator="[
                  'configurationPageUUID',
                  {
                    rules: [{ required: false, message: $t('device.deviceConfigurationPageName') }],
                  },
                ]"
          >
            <a-select-option v-for="(page,index) in displayPageList" :key="index" :value="page.value">
              {{ page.label }}
            </a-select-option>
          </a-select>
        </a-form-item>
      <a-form-item
          :label="$t('dataModel.Port')"
          :labelCol="{span: 7}"
          :wrapperCol="{span: 10}"
      >
        <a-input autocomplete="autocomplete"

                 v-decorator="['port', {rules: [{ required: true, message: $t('dataModel.Port'), whitespace: true,initialValue:162}]}]"
        />
      </a-form-item>
      <a-form-item
            :label="$t('dataModel.gatherNumber')"
            :labelCol="{span: 7}"
            :wrapperCol="{span: 10}"
        >
          <a-input autocomplete="autocomplete"

                   v-decorator="['gatherNumber', {rules: [{ required: true, message: $t('dataModel.gatherNumber'), whitespace: true}]}]"
          />
        </a-form-item>
      <a-form-item
          :label="$t('dataModel.snmpVersion')"
          :labelCol="{span: 7}"
          :wrapperCol="{span: 10}"
      >
        <a-select @change="handleSelectChange"  autocomplete="autocomplete"

                  v-decorator="['version', {rules: [{ required: true, message: $t('dataModel.snmpVersion'), whitespace: true}]}]"
        >
          <a-select-option value=1>V1</a-select-option>
          <a-select-option value=2>V2</a-select-option>
          <a-select-option value=3>V3</a-select-option>
        </a-select>
      </a-form-item>
      <div v-if="version!==3">
        <a-form-item
            :label="$t('dataModel.snmpWriteComm')"
            :labelCol="{span: 7}"
            :wrapperCol="{span: 10}"
        >
          <a-input autocomplete="autocomplete"

                   v-decorator="['writecomm', {rules: [{ required: true, message: $t('dataModel.snmpWriteComm'), whitespace: true}]}]"/>
        </a-form-item>
        <a-form-item
            :label="$t('dataModel.snmpReadComm')"
            :labelCol="{span: 7}"
            :wrapperCol="{span: 10}"
        >
          <a-input autocomplete="autocomplete"

                   v-decorator="['readcomm', {rules: [{ required: true, message: $t('dataModel.snmpReadComm'), whitespace: true}]}]"/>
        </a-form-item>
      </div>

      <div v-if="version===3">
        <a-form-item
            :label="$t('dataModel.snmpUserName')"
            :labelCol="{span: 7}"
            :wrapperCol="{span: 10}"
        >
          <a-input-password autocomplete="autocomplete"

                   v-decorator="['snmpUserName', {rules: [{ required: true, message: $t('dataModel.snmpUserName'), whitespace: true}]}]"/>
        </a-form-item>
        <a-form-item
            :label="$t('dataModel.snmpSecurityLevel')"
            :labelCol="{span: 7}"
            :wrapperCol="{span: 10}"
        >
          <a-select style="" autocomplete="autocomplete"

                    v-decorator="['snmpSecurityLevel', {rules: [{ required: true, message: $t('dataModel.snmpSecurityLevel'), whitespace: true}]}]">
            <a-select-option value=1>{{$t('dataModel.snmpLevelNoAuthNoPrivacy')}}</a-select-option>
            <a-select-option value=2>{{$t('dataModel.snmpLevelAuthNoPrivacy')}}</a-select-option>
            <a-select-option value=3>{{$t('dataModel.snmpLevelAuthPrivacy')}}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item
            :label="$t('dataModel.snmpAuthAlgorithm')"
            :labelCol="{span: 7}"
            :wrapperCol="{span: 10}"
        >
          <a-select  style="" autocomplete="autocomplete"

                     v-decorator="['snmpAuthAlgorithm', {rules: [{ required: true, message: $t('dataModel.snmpAuthAlgorithm'), whitespace: true}]}]">
            <a-select-option value=1>MD5</a-select-option>
            <a-select-option value=2>SHA</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item
            :label="$t('dataModel.snmpUserPassword')"
            :labelCol="{span: 7}"
            :wrapperCol="{span: 10}"
        >
          <a-input-password autocomplete="autocomplete"

                   v-decorator="['snmpUserPassword', {rules: [{ required: true, message: $t('dataModel.snmpUserPassword'), whitespace: true}]}]"/>
        </a-form-item>
        <a-form-item
            :label="$t('dataModel.snmpPrivacyAlgorithm')"
            :labelCol="{span: 7}"
            :wrapperCol="{span: 10}"
        >
          <a-select style="" autocomplete="autocomplete"

                    v-decorator="['snmpPrivacyAlgorithm', {rules: [{ required: true, message: $t('dataModel.snmpPrivacyAlgorithm'), whitespace: true}]}]">
            <a-select-option value=1>DES</a-select-option>
            <a-select-option value=2>AES</a-select-option>
            <a-select-option value=3>AES192</a-select-option>
            <a-select-option value=4>AES256</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item
            :label="$t('dataModel.snmpPrivacyPassword')"
            :labelCol="{span: 7}"
            :wrapperCol="{span: 10}"
        >
          <a-input-password autocomplete="autocomplete"

                   v-decorator="['snmpPrivacyPassword', {rules: [{ required: true, message: $t('dataModel.snmpPrivacyPassword'), whitespace: true}]}]"/>
        </a-form-item>
      </div>
      <a-form-item style="margin-top: 24px" :wrapperCol="{span: 10, offset: 7}">
        <a-button type="primary" htmlType="submit">{{$t('dataModel.edit')}}</a-button>
        <a-button style="margin-left: 8px" @click="onBlackCLK()">{{$t('dataModel.back')}}</a-button>
      </a-form-item>
    </a-form>
    </a-spin>
  </a-card>
  <a-card v-else-if="defaultTabKey==2"  >
    <a-space class="operator">
      <a-upload
          name="file"
          :multiple="false"
          :action=importUrl
          :showUploadList="false"
          :beforeUpload="beforeUpload"
          @change="mibImport"
      >
        <a-button type="primary"> <a-icon type="upload" />
          {{$t("dataModel.importMib")}}</a-button>
      </a-upload>

      <a-upload
          name="file"
          :multiple="false"
          :action=importXML
          :showUploadList="false"
          :beforeUpload="beforeUpload"
          @change="mibImportxml"
      >
        <a-button type="primary"> <a-icon type="upload" />
          {{$t("dataModel.importXML")}}</a-button>
      </a-upload>

      <a-button @click="saveMibData()"> <a-icon type="save" />
        {{$t("dataModel.saveMib")}}</a-button>

      <a-button type="danger" :disabled="selectDataTableUuid.length==0" @click="deleteRecord"> <a-icon type="delete"  />
        {{$t('dataModel.delete')}}</a-button>

      <a-button   @click="backToList"> <a-icon type="backward" />
        {{$t("dataModel.modbusModel.Back")}}</a-button>

      <a-button  type="link">
        <a href="http://www.ismctl.com/ism/templete.zip">{{$t("dataModel.SnmpTempleteFile")}}</a>
      </a-button>
    </a-space>
    <a-spin style="padding: 1px;"  :spinning="messageShowLoad" tip="Loading...">
      <a-table rowKey='oidPath' :pagination="pagination" :row-selection="rowSelection" :columns="columns" :data-source="dataSource" class="ant-table-tbody">
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
        <template slot="customRender" slot-scope="text, record, index, column">
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
          >{{ fragment }}</mark
          >
          <template v-else>{{ fragment }}</template>
        </template>
      </span>
          <template v-else>
            {{ text }}
          </template>
        </template>

      <template v-for="(item, index) in columns" :slot="item.slotName">
        <span :key="index">{{ $t(item.slotName) }}</span>
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
      <a-drawer
          :title="$t('dataModel.EditDataModel')"
          :width="720"
          :visible="editVisible"
          :body-style="{ paddingBottom: '80px' }"
          @close="onClose"
      >
        <a-spin style="padding: 1px;"  :spinning="messageShowLoad" tip="Loading...">
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
              <a-form-item :label="$t('dataModel.editData.OID')">
                <a-input
                    v-decorator="[
                  'oid',
                  {
                    rules: [{ required: true, message: $t('dataModel.editData.OID') }],
                  },
                ]"
                />
              </a-form-item>
            </a-col>
          </a-row>
          <a-row :gutter="16">

            <a-col :span="12">
              <a-form-item :label="$t('dataModel.editData.dataType')">
                <a-select   autocomplete="autocomplete"   v-decorator="[
                  'dataType',
                  {
                    rules: [{ required: true, message: $t('dataModel.editData.dataType') }],
                  },
                ]">
                  <a-select-option value="Integer">Integer</a-select-option>
                  <a-select-option value="OctetString">OctetString</a-select-option>
                  <a-select-option value="OID">OID</a-select-option>
                  <a-select-option value="Gauge">Gauge</a-select-option>
                  <a-select-option value="Counter32">Counter32</a-select-option>
                  <a-select-option value="IpAddress">IpAddress</a-select-option>
                  <a-select-option value="TimeTicks">TimeTicks</a-select-option>
                  <a-select-option value="Counter64">Counter64</a-select-option>
                  <a-select-option value="UnsignedInteger">UnsignedInteger</a-select-option>
                </a-select>
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
                  <a-select-option value="ReadCreate">ReadCreate</a-select-option>
                  <a-select-option value="AccessibleForNotify">AccessibleForNotify</a-select-option>
                  <a-select-option value="NotAccessible">NotAccessible</a-select-option>
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
                  <a-select   @change="alarmCharge"  autocomplete="autocomplete"  v-decorator="[
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
            <div v-else>
              <a-col :span="12" >
                <a-form-item :label="$t('dataModel.editData.dataRecord')">
                  <a-select  @change="recordCharge" autocomplete="autocomplete"  v-decorator="[
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
        </a-spin>
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
  </a-card>
</template>

<script>
import {
  snmpModelEdit,
  getSnmpModelDetail,
  snmpModelMibSave,
  snmpModelGetMibs, snmpModelDeleteMibs,modelDataEdit
} from "@/services/snmpmodel";
import {IMPORTMIB,IMPORTMIBXML} from '@/services/api'
import { uuid } from 'vue-uuid';
import {displayModelList, getDisplayModelLayerData} from "@/services/displayModel";

const dataSource= []
export default {
  name: 'SnmpModelImport',
  i18n: require('../../../i18n/language'),
  data () {
    this.cacheData = dataSource.map(item => ({ ...item }));
    return {
      pagination:{
        pageSize:15,
        showSizeChanger:true
      },
      DataRecordType:0,
      configurationModel:[],
      displayPageList:[],
      searchText: '',
      searchInput: null,
      searchedColumn: '',
      alarmStatus:0,
      recordStatus:0,
      selectDataTableUuid:[],
      rowSelection:{
        onSelect:this.onDataTableSelect,
        onSelectAll:this.onDataTableSelectAll
      },
      editVisible:false,
      error: '',
      editingKey: '',
      importUrl:IMPORTMIB,
      importXML:IMPORTMIBXML,
      defaultTabKey:1,
      messageShowLoad:false,
      EditForm:this.$form.createForm(this),
      form: this.$form.createForm(this),
      version:1,
      columns: [
        {
          slotName: this.$t("dataModel.oidName"),
          scopedSlots: { filterDropdown: 'filterDropdown', filterIcon: 'filterIcon', customRender: 'oidName' ,title:this.$t("dataModel.oidName")},
          width: '20%',
          dataIndex: 'oidName',
          onFilter: (value, record) =>
              record.oidName
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
          slotName: this.$t("dataModel.oidPath"),
          scopedSlots: {  customRender: 'oidPath' ,title:this.$t("dataModel.oidPath")},
          width: '25%',
          dataIndex: 'oidPath',
        },
        {
          slotName:this.$t("dataModel.oidType"),
          scopedSlots: {  customRender: 'oidType' ,title:this.$t("dataModel.oidType") },
          width: '10%',
          dataIndex: 'oidType',
        },
        {
          slotName:this.$t("dataModel.oidAuth"),
          scopedSlots: {  customRender: 'oidAuth'  ,title:this.$t("dataModel.oidAuth")},
          width: '10%',
          dataIndex: 'oidAuth',
        },
        {
          slotName:this.$t("dataModel.DataUnit"),
          scopedSlots: {  customRender: 'dataUnit'  ,title:this.$t("dataModel.DataUnit")},
          width: '10%',
          dataIndex: 'dataUnit'
        },
        {
          title: this.$t('dataModel.modelTableOpt'),
          width: '15%',
          scopedSlots: { customRender: 'action' }
        }
      ],
      dataSource,
      selectedRows: [],
      value: 1
    }
  },
  activated() {
    this.defaultTabKey = this.$route.params.tab
    if(this.defaultTabKey==1)
    {
      this.getSingleModelDetail()
    }
    else
    {
      this.getAllMibs()
    }
  },
  mounted() {
    this.getConfigurationModel()
  },
  computed: {
    desc() {
      return this.$t('pageDesc')
    }
  },
  methods: {
    beforeUpload(){
      this.$message.loading({ content: this.$t('dataModel.opcuaModel.ImportNodeIDLoading'),duration: 0 });
    },
    getConfigurationModel(){
      this.configurationModel=[]
      let _t = this
      const  params= {
        DisplayType: 1
      }
      displayModelList(params).then(function (res){
        let tableData={}
        if(res.data.list!=null)
        {
          for(let i=0;i<res.data.list.length;i++)
          {
            tableData.name = res.data.list[i].name
            tableData.description = res.data.list[i].description
            tableData.uuid = res.data.list[i].displayUid
            _t.configurationModel.push(tableData)
            tableData={}
          }
        }

      })
    },
    GetDisplayPage(uuid){
      let params={
        muid:uuid
      }
      let _t = this
      getDisplayModelLayerData(params).then(function (res){
        _t.displayPageList = []
        if(res.data.code==0)
        {
          let pageLayer = res.data.layer
          if(pageLayer.length>0)
          {
            for(let i=0;i<pageLayer.length;i++)
            {
              let pageInfo = {}
              pageInfo.label = pageLayer[i].PageName
              pageInfo.value = pageLayer[i].PageId
              pageInfo.pageType = pageLayer[i].PageType
              pageInfo.pageModelUuid = pageLayer[i].modelId
              _t.displayPageList.push(pageInfo)
            }
          }
        }
      })
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
    onClose() {
      this.editVisible = false;
      this.messageShowLoad = false
    },
    getAllMibs(){
      let params={
        muid:this.$route.params.uid,
        type:1
      }
      this.dataSource = []
      let _t = this
      this.messageShowLoad = true
      snmpModelGetMibs(params).then(function (res){
        _t.messageShowLoad = false
        if(res.data.mibs.length>0)
        {
          let temp={}
          for(let i=0;i<res.data.mibs.length;i++)
          {
            temp.key = i.toString()
            temp.oidName = res.data.mibs[i].name
            temp.oidPath = res.data.mibs[i].oid
            temp.oidType = res.data.mibs[i].type
            temp.oidAuth = res.data.mibs[i].auth
            temp.uuid = res.data.mibs[i].uuid
            temp.dataUnit = res.data.mibs[i].unit
            temp.alarm=parseInt(res.data.mibs[i].alarm)
            temp.alarmLevel = res.data.mibs[i].alarmLevel
            temp.AlarmMessage =res.data.mibs[i].AlarmMessage
            temp.AlarmClearMessage = res.data.mibs[i].AlarmClearMessage
            temp.conversionExpression = res.data.mibs[i].conversionExpression
            temp.record=parseInt(res.data.mibs[i].record)
            temp.RecordType=parseInt(res.data.mibs[i].RecordType)
            temp.RecordDataCharge=res.data.mibs[i].RecordDataCharge.toString()
            temp.recordInterval=parseInt(res.data.mibs[i].recordInterval.toString())
            _t.dataSource.push(temp)
            temp={}
          }
          _t.cacheData = _t.dataSource.map(item => ({ ...item }));
        }
      })
    },
    handleChange(value, key, column) {
      const newData = [...this.dataSource];
      const target = newData.filter(item => key === item.key)[0];
      if (target) {
        target[column] = value;
        this.dataSource = newData;
      }
    },
    selectChange(value,key, column){
      const newData = [...this.dataSource];
      const target = newData.filter(item => key === item.key)[0];
      if (target) {
        target[column] = value;
        this.dataSource = newData;
      }
    },
    edit(item) {
      let _t = this
      this.messageShowLoad = true
      this.editingKey = item.uuid
      this.alarmStatus = parseInt(item.alarm)
      this.recordStatus = parseInt(item.record)
      this.ShowRegisterLoading = true
      this.DataRecordType = item.RecordType ?item.RecordType:0
      if(item.recordInterval==0)
      {
        item.recordInterval=1
      }
      setTimeout(function (){
        _t.EditForm.setFieldsValue(
            {
              name:item.oidName,
              oid:item.oidPath,
              dataType:item.oidType,
              dataAuth:item.oidAuth,
              dataUnit:item.dataUnit,
              ConversionExpression:item.conversionExpression,
            })

        if (item.alarm==1){
          _t.EditForm.setFieldsValue(
              {
                dataAlarm:item.alarm.toString(),
                AlarmLevel:(typeof item.alarmLevel!="undefined")?item.alarmLevel.toString():"",
                dataRecord:(typeof item.record!="undefined")?item.record.toString():"0",
                AlarmMessage :(typeof item.AlarmMessage!="undefined")?item.AlarmMessage:"",
                AlarmClearMessage : (typeof item.AlarmClearMessage!="undefined")?item.AlarmClearMessage:"",
              })
        }
        else  if (item.record==1)
        {
          _t.EditForm.setFieldsValue(
              {
                dataAlarm:(typeof item.alarm!="undefined")?item.alarm.toString():"0",
                dataRecord:item.record.toString(),
                dataRecordType:(typeof item.RecordType!="undefined")?item.RecordType.toString():"",
                dataRecordChargeValue:(typeof item.RecordDataCharge!="undefined")?item.RecordDataCharge.toString():"",
                dataRecordTime:(typeof item.recordInterval!="undefined")?item.recordInterval.toString():"",
              })
        }
        else
        {
          _t.EditForm.setFieldsValue(
              {
                dataRecord:item.record?item.record.toString():"0",
                dataAlarm:item.alarm?item.alarm.toString():"0",
              })
        }
        _t.messageShowLoad = false
        _t.ShowRegisterLoading = false
      },500)

      this.editVisible = true;
    },
    save() {
      this.EditForm.validateFields((err) => {
        if (!err) {
          this.messageShowLoad = true
          let params = {
            muid:this.$route.params.uid,
            uuid:this.editingKey,
            editData:{
              name:this.EditForm.getFieldValue('name'),
              oid:this.EditForm.getFieldValue('oid'),
              auth:this.EditForm.getFieldValue('dataAuth'),
              type:this.EditForm.getFieldValue('dataType'),
              unit:this.EditForm.getFieldValue('dataUnit'),
              alarm:parseInt(this.EditForm.getFieldValue('dataAlarm')),
              alarmLevel:parseInt(this.EditForm.getFieldValue('AlarmLevel')),
              AlarmMessage:this.EditForm.getFieldValue('AlarmMessage'),
              AlarmClearMessage:this.EditForm.getFieldValue('AlarmClearMessage'),
              conversionExpression:this.EditForm.getFieldValue('ConversionExpression'),
              record:this.EditForm.getFieldValue('dataRecord')?parseInt(this.EditForm.getFieldValue('dataRecord')):0,
              RecordType:this.EditForm.getFieldValue('dataRecordType')?parseInt(this.EditForm.getFieldValue('dataRecordType')):0,
              recordInterval:this.EditForm.getFieldValue('dataRecordTime')?parseInt(this.EditForm.getFieldValue('dataRecordTime')):0,
              RecordDataCharge:this.EditForm.getFieldValue('dataRecordChargeValue')?this.EditForm.getFieldValue('dataRecordChargeValue').toString():"",
            }
          }
          let _t = this
          modelDataEdit(params).then(function (res){
            if(res.data.code==0)
            {
              const newData = [..._t.dataSource];
              const target = newData.filter(item => _t.editingKey === item.uuid)[0];
              if (target) {
                    target.oidName = _t.EditForm.getFieldValue('name')
                    target.oidPath=_t.EditForm.getFieldValue('oid')
                    target.oidAuth=_t.EditForm.getFieldValue('dataAuth')
                    target.oidType=_t.EditForm.getFieldValue('dataType')
                    target.dataUnit=_t.EditForm.getFieldValue('dataUnit')
                    target.alarm=parseInt(_t.EditForm.getFieldValue('dataAlarm'))
                    target.conversionExpression=_t.EditForm.getFieldValue('ConversionExpression')
                    target.alarmLevel=parseInt(_t.EditForm.getFieldValue('AlarmLevel'))
                    target.AlarmMessage = _t.EditForm.getFieldValue('AlarmMessage')
                    target.AlarmClearMessage = _t.EditForm.getFieldValue('AlarmClearMessage')
                    target.record=parseInt(_t.EditForm.getFieldValue('dataRecord'))
                    target.recordInterval=parseInt(_t.EditForm.getFieldValue('dataRecordTime'))
                    target.RecordType=parseInt(_t.EditForm.getFieldValue('dataRecordType'))
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
            _t.messageShowLoad = false
          }).finally(function (){
              _t.messageShowLoad = false
          })
        }
      })
    },
    saveMibData(){
      let saveParams=[]
      this.messageShowLoad = true
      for(let i=0;i<this.dataSource.length;i++)
      {
        let temp = {}
        temp.muid = this.$route.params.uid
        temp.name = this.dataSource[i].oidName
        temp.oid = this.dataSource[i].oidPath
        if( this.dataSource[i].uuid==""|| typeof this.dataSource[i].uuid == "undefined")
        {
          temp.uuid = uuid.v1()
        }
        else
        {
          temp.uuid = this.dataSource[i].uuid
        }
        temp.auth = this.dataSource[i].oidAuth
        temp.type = this.dataSource[i].oidType

        temp.conversionExpression = (typeof this.dataSource[i].conversionExpression!="undefined")? this.dataSource[i].conversionExpression:""

        temp.alarmLevel = (typeof this.dataSource[i].alarmLevel!="undefined")?this.dataSource[i].alarmLevel:0
        temp.AlarmClearMessage = (typeof this.dataSource[i].AlarmClearMessage!="undefined")?this.dataSource[i].AlarmClearMessage:""
        temp.AlarmMessage =  (typeof this.dataSource[i].AlarmMessage!="undefined")?this.dataSource[i].AlarmMessage:""
        temp.alarm=(typeof this.dataSource[i].alarm!="undefined")?parseInt(this.dataSource[i].alarm):0

        temp.record=(typeof this.dataSource[i].record!="undefined")?parseInt(this.dataSource[i].record):0
        temp.recordInterval=(typeof this.dataSource[i].recordInterval!="undefined")?parseInt(this.dataSource[i].recordInterval):0
        temp.RecordType = (typeof this.dataSource[i].RecordType!="undefined")?this.dataSource[i].RecordType:0
        temp.RecordDataCharge = (typeof this.dataSource[i].RecordDataCharge!="undefined")?this.dataSource[i].RecordDataCharge:""
        saveParams.push(temp)
        temp={}
      }
      let _t = this
      snmpModelMibSave(saveParams).then(function (res){
        if(res.data.code==0)
        {
          _t.$message.success(_t.$t("dataModel.saveSuccess"));
          _t.getAllMibs()
        }
        else
        {
          _t.$message.error(_t.$t("dataModel.saveFailed"));
        }
        _t.messageShowLoad = false
      }).catch(function (){
        _t.messageShowLoad = false
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
    },
    mibImport(info) {
      this.dataSource=[]
      this.messageShowLoad = true
      if (info.file.status === 'done') {
        let result = info.file.response
        this.$message.destroy();
        if(result.Code==0) {
          let tempData = {}
          if (result.Nodes != null) {
            this.$message.success(`${info.file.name} `+this.$t("dataModel.parseSuccess"));
            for (let i = 0; i < result.Nodes.length; i++) {
              if ((result.Nodes[i].Access != 'Unknown')&&(result.Nodes[i].Oid!=null)) {
                tempData.key = i.toString()
                tempData.oidName = result.Nodes[i].Name
                tempData.oidPath = result.Nodes[i].Oid.join(".")+'.0'
                if(result.Nodes[i].Type==null)
                {
                  tempData.oidType = result.Nodes[i].Decl
                }
                else {
                  tempData.oidType = result.Nodes[i].Type.BaseType
                }

                tempData.oidAuth = result.Nodes[i].Access
                tempData.uuid = ""
                tempData.dataAlarm=0
                tempData.dataRecord=0
                tempData.dataRecordTime=0
                this.dataSource.push(tempData)
                tempData = {}
              }
            }
            this.cacheData = this.dataSource.map(item => ({ ...item }));
          }
          else
          {
            this.$message.error(`${info.file.name} `+this.$t("dataModel.parseFailed"));
          }
        }
        else
        {
          this.$message.error(`${info.file.name} `+this.$t("dataModel.parseFailed"));
        }
      }
      this.messageShowLoad = false
    },
    mibImportxml(info) {
      this.dataSource=[]
      this.messageShowLoad = true
      if (info.file.status === 'done') {
        let result = info.file.response
        this.$message.destroy();
        if(result.Code==0) {
          let tempData = {}
          try {
            if (result.Nodes != null && result.Nodes.Instances != null) {
              this.$message.success(`${info.file.name} ` + this.$t("dataModel.parseSuccess"));
              for (let i = 0; i < result.Nodes.Instances.Instance.length; i++) {
                tempData.key = i.toString()
                tempData.oidName = result.Nodes.Instances.Instance[i].Name
                tempData.oidPath = result.Nodes.Instances.Instance[i].Oid
                if(tempData.oidPath.substr(0, 1) == "."){
                  tempData.oidPath = tempData.oidPath.slice(1);
                }
                tempData.oidType = result.Nodes.Instances.Instance[i].ValueType
                tempData.oidAuth = "ReadWrite"
                tempData.uuid = ""
                tempData.dataAlarm = 0
                tempData.dataRecord = 0
                tempData.dataRecordTime = 0
                this.dataSource.push(tempData)
                tempData = {}
              }

              this.cacheData = this.dataSource.map(item => ({...item}));
            } else {
              this.$message.error(`${info.file.name} ` + this.$t("dataModel.parseFailed"));
            }
          }catch (e) {
            this.$message.error(`${info.file.name} ` + this.$t("dataModel.parseFailed"));
          }
        }
        else
        {
          this.$message.error(`${info.file.name} `+this.$t("dataModel.parseFailed"));
        }
      }
      this.messageShowLoad = false
    },
    getSingleModelDetail(){
      let _t = this
      const params={
        uuid:this.$route.params.uid
      }
      this.messageShowLoad = true
      getSnmpModelDetail(params).then(function (res){
        _t.version = res.data.data.version
        _t.GetDisplayPage(res.data.data.configUid)
        _t.$nextTick(function (){
          _t.messageShowLoad = false
          _t.form.setFieldsValue(
              {
                name:res.data.data.name,
                dec:res.data.data.dec,
                configurationModel:res.data.data.configUid,
                configurationPageUUID:res.data.data.PageUUID,
                port:res.data.data.port.toString(),
                version:res.data.data.version.toString(),
                gatherNumber:res.data.data.gatherNumber.toString(),
              })

          if(_t.version!=3)
          {
            _t.form.setFieldsValue(
                {
                  writecomm:res.data.data.writecomm,
                  readcomm:res.data.data.readcomm,
                }
            )
          }
          else
          {
            _t.form.setFieldsValue(
                {
                  snmpUserName:res.data.data.snmpUserName,
                  snmpSecurityLevel:res.data.data.snmpSecurityLevel.toString(),
                  snmpAuthAlgorithm:res.data.data.snmpAuthAlgorithm.toString(),
                  snmpUserPassword:res.data.data.snmpUserPassword,
                  snmpPrivacyAlgorithm:res.data.data.snmpPrivacyAlgorithm.toString(),
                  snmpPrivacyPassword:res.data.data.snmpPrivacyPassword,
                }
            )
          }
        },500)
      })
    },
    onSubmit (e) {
      e.preventDefault()
      this.form.validateFields((err) => {
        if (!err) {
          this.logging = true
          this.messageShowLoad = true
          const params = {
            uuid:this.$route.params.uid,
            data: {
              name: this.form.getFieldValue('name'),
              dec: this.form.getFieldValue('dec'),
              configUid:this.form.getFieldValue('configurationModel'),
              PageUUID:this.form.getFieldValue('configurationPageUUID'),
              port:parseInt(this.form.getFieldValue('port')),
              version: parseInt(this.form.getFieldValue('version')),
              writecomm: this.form.getFieldValue('writecomm'),
              readcomm: this.form.getFieldValue('readcomm'),
              gatherNumber:parseInt(this.form.getFieldValue('gatherNumber')),
              snmpUserName: this.form.getFieldValue('snmpUserName'),
              snmpSecurityLevel: parseInt(this.form.getFieldValue('snmpSecurityLevel')),
              snmpAuthAlgorithm: parseInt(this.form.getFieldValue('snmpAuthAlgorithm')),
              snmpUserPassword: this.form.getFieldValue('snmpUserPassword'),
              snmpPrivacyAlgorithm: parseInt(this.form.getFieldValue('snmpPrivacyAlgorithm')),
              snmpPrivacyPassword: this.form.getFieldValue('snmpPrivacyPassword'),
            }
          };
          snmpModelEdit(params).then(this.afterLogin)
        }
      })
    },
    afterLogin(res) {
      this.logging = false
      if (res.data.code == 200) {
        this.$message.success(this.$t('dataModel.editSuccess'), 3)
      }
      else
      {
        this.$message.error(this.$t('dataModel.editFailed'), 3)
      }
      this.messageShowLoad = false
    },
    onBlackCLK(){
      this.$router.push('/DeviceModel/SnmpModel')
    },
    handleSelectChange(value) {
      this.version=parseInt(value)
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
          _t.messageShowLoad = true
          let params={
            muid:_t.$route.params.uid,
            uuid:_t.selectDataTableUuid
          }
          snmpModelDeleteMibs(params).then(function (res) {
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
            _t.messageShowLoad = false
          })
        },
        cancelText: _t.$t('displayModel.ConfirmCancel'),
        onCancel() {

        },
      });


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
    backToList(){
      this.$router.push("/DeviceModel/SnmpModel")
    },
    chargeDataRecordType(value){
      this.DataRecordType = parseInt(value)
    }
  }
}
</script>


<style lang="less">
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
.editable-row-operations a {
  margin-right: 8px;
}
.ant-table-tbody > tr > td {
  padding: 7px 7px;
  overflow-wrap: break-word;
}
.ant-form-item {
  margin-bottom: 5px;

}
</style>