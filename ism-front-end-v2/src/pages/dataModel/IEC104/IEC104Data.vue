<template>
  <div>
    <!--    寄存器组表格-->
    <a-card  style="min-height: 400px">
      <a-space class="operator">

        <a-button type="primary" @click="RegisterVisible=true;isEdit=false"> <a-icon type="plus" />
          {{$t("dataModel.DLT645Model.AddData")}}</a-button>



        <a-button type="link" @click="handleExport"> <a-icon type="export" />{{$t('dataModel.export')}}</a-button>

        <a-upload
            name="file"
            :multiple="false"
            :action=localUpgradeUrl
            :showUploadList="false"
            :beforeUpload="beforeUpload"
            @change="localUpgradeCharge"
        >
          <a-button type="link"> <a-icon type="import" />
            {{$t('dataModel.import')}}
          </a-button>
        </a-upload>

        <a-button type="link" @click="onBlackCLK()"> <a-icon type="backward" />
          {{$t("dataModel.opcuaModel.Back")}}</a-button>
      </a-space>
      <a-spin style="padding: 1px;"  :spinning="messageShowLoad" tip="Loading...">
        <a-table  :pagination="pagination" rowKey="name" :columns="registerGroupColumns" :data-source="registerGroupDataSource" class="ant-table-tbody">
          <template v-for="(item, index) in registerGroupColumns" :slot="item.slotName">
            <span :key="index">{{ $t(item.slotName) }}</span>
          </template>
          <template slot="DataCategory" slot-scope="text">
            <span v-if="text==1"> {{$t('dataModel.IEC104Model.DataCategoryYaoXin')}}</span>
            <span v-else-if="text==2"> {{$t('dataModel.IEC104Model.DataCategoryYaoCe')}}</span>
            <span v-else-if="text==3"> {{$t('dataModel.IEC104Model.DataCategoryMaiChong')}}</span>
            <span v-else-if="text==4"> {{$t('dataModel.IEC104Model.DataCategoryYaoKong')}}</span>
            <span v-else-if="text==5"> {{$t('dataModel.IEC104Model.DataCategoryYaoTiao')}}</span>
          </template>
          <template slot="NodeIDDataType" slot-scope="text">
            <span v-if="text==8"> Int</span>
            <span v-else-if="text==10"> Float</span>
          </template>
          <template slot="auth" slot-scope="text">
            <span v-if="text=='ReadOnly'"> ReadOnly</span>
            <span v-else-if="text=='ReadWrite'"> ReadWrite</span>
            <span v-else-if="text=='WriteOnly'"> WriteOnly</span>
          </template>
          <template slot="action" slot-scope="text, record">
            <div class="editable-row-operations">
              <span >
                <a  @click="() => edit(record)">
                  <a-icon type="edit" /> {{$t('dataModel.opcuaModel.NodeIDEdit')}}</a>
              </span>
              <a-popconfirm :title="$t('dataModel.deleteConfirm')" @confirm="deleteRecord(record.uuid,record.muid)">
                <a-icon slot="icon" type="question-circle-o" style="color: red" />
                <a-icon type="delete" theme="twoTone" two-tone-color="#eb2f96"/><a style="color: #eb2f96">{{$t('dataModel.delete')}}</a>
              </a-popconfirm>
            </div>
          </template>
        </a-table>
      </a-spin>
    </a-card>
    <!--    添加节点-->

    <a-drawer
        :title="isEdit?$t('dataModel.opcuaModel.EditNodeID'):$t('dataModel.opcuaModel.AddNodeID')"
        :width="720"
        :visible="RegisterVisible"
        :body-style="{ paddingBottom: '80px' }"
        @close="onClose"
    >
      <a-spin style="padding: 1px;"  :spinning="ShowRegisterLoading" tip="Loading...">
        <a-form :form="RegisterForm" layout="vertical" @submit="AddNodeId">
          <a-alert type="error" :closable="true" v-show="error" :message="error" showIcon style="margin-bottom: 24px;" />
          <a-row :gutter="16">
            <a-col :span="12">
              <a-form-item
                  :label="$t('dataModel.opcuaModel.NodeIDName')"
              >
                <a-input autocomplete="autocomplete"
                         v-decorator="['NodeIDName', {rules: [{ required: true, validator: isValidateTxtNonSpec, message: $t('device.deviceNameVal'), whitespace: true}]}]"
                />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item
                  :label="$t('dataModel.IEC104Model.DataFlag')"
              >
                <a-input autocomplete="autocomplete"
                         v-decorator="['DataFlag', {rules: [{ required: true, message: $t('dataModel.IEC104Model.DataFlag'), whitespace: true}]}]"
                />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item
                  :label="$t('dataModel.IEC104Model.DataCategory')"
              >
                <a-select @change="DataCategoryAlarmCharge" class="DataType" autocomplete="autocomplete"
                          v-decorator="['DataCategory', {initialValue:'1',rules: [{ required: true, message: $t('dataModel.IEC104Model.DataCategory'), whitespace: true}]}]"
                >

                  <a-select-option value="1">{{$t('dataModel.IEC104Model.DataCategoryYaoXin')}}</a-select-option>
                  <a-select-option value="2">{{$t('dataModel.IEC104Model.DataCategoryYaoCe')}}</a-select-option>
                  <a-select-option value="3">{{$t('dataModel.IEC104Model.DataCategoryMaiChong')}}</a-select-option>
                  <a-select-option value="4">{{$t('dataModel.IEC104Model.DataCategoryYaoKong')}}</a-select-option>
                  <a-select-option value="5">{{$t('dataModel.IEC104Model.DataCategoryYaoTiao')}}</a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col v-if="SelectDataCategory==4" :span="12">
              <a-form-item
                  :label="$t('dataModel.IEC104Model.DataCategoryYaoKongType')"
              >
                <a-select class="DataType" autocomplete="autocomplete"
                          v-decorator="['DataCategoryYaoKongType', {initialValue:'1',rules: [{ required: true, message: $t('dataModel.IEC104Model.DataCategoryYaoKongType'), whitespace: true}]}]"
                >

                  <a-select-option value="1">{{$t('dataModel.IEC104Model.DataCategoryYaoKongSingle')}}</a-select-option>
                  <a-select-option value="2">{{$t('dataModel.IEC104Model.DataCategoryYaoKongDouble')}}</a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col v-if="SelectDataCategory==5||SelectDataCategory==2" :span="12">
              <a-form-item
                  :label="$t('dataModel.IEC104Model.DataCategoryYaoTiaoType')"
              >
                <a-select  @change="DataCategoryYaoTiaoGuiYiCharge" class="DataType" autocomplete="autocomplete"
                          v-decorator="['DataCategoryYaoTiaoType', {initialValue:'3',rules: [{ required: true, message: $t('dataModel.IEC104Model.DataCategoryYaoTiaoType'), whitespace: true}]}]"
                >

                  <a-select-option value="1">{{$t('dataModel.IEC104Model.DataCategoryYaoTiaoGuiYi')}}</a-select-option>
                  <a-select-option value="2">{{$t('dataModel.IEC104Model.DataCategoryYaoTiaoBZ')}}</a-select-option>
                  <a-select-option value="3">{{$t('dataModel.IEC104Model.DataCategoryYaoTiaoFloat')}}</a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :span="12" v-if="SelectDataCategory==2&&DataCategoryYaoTiaoType==1">
              <a-form-item :label="$t('dataModel.IEC104Model.DataCategoryYaoTiaoGuiYiED')">
                <a-tooltip placement="top">
                  <template slot="title">
                    <span>{{$t('dataModel.IEC104Model.DataCategoryYaoTiaoGuiYiED')}}</span>
                  </template>
                  <a-input   autocomplete="autocomplete"   v-decorator="[
                        'DataCategoryYaoTiaoGuiYiED',
                        {
                          initialValue:'32768',rules: [{ required: true, message: $t('dataModel.IEC104Model.DataCategoryYaoTiaoGuiYiED') }],
                        },
                      ]">
                  </a-input>
                </a-tooltip>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item
                  :label="$t('dataModel.opcuaModel.NodeIDDataType')"
              >
                <a-select class="DataType" autocomplete="autocomplete"
                          v-decorator="['NodeIDDataType', {initialValue:'10',rules: [{ required: true, message: $t('dataModel.opcuaModel.NodeIDDataType'), whitespace: true}]}]"
                >

                  <a-select-option value="8">Int</a-select-option>
                  <a-select-option value="10">Float</a-select-option>
                </a-select>
              </a-form-item>
            </a-col>


            <a-col :span="12">
              <a-form-item
                  :label="$t('dataModel.opcuaModel.NodeIDAccessLevel')"
              >
                <a-select  class="DataType" autocomplete="autocomplete"
                           v-decorator="['NodeIDAccessLevel', {initialValue:'ReadWrite',rules: [{ required: true, message: $t('dataModel.opcuaModel.NodeIDAccessLevel'), whitespace: true}]}]"
                >
                  <a-select-option value="ReadOnly">ReadOnly</a-select-option>
                  <a-select-option value="ReadWrite">ReadWrite</a-select-option>
                  <a-select-option value="WriteOnly">WriteOnly</a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :span="12">
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
            <a-col :span="12">
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
            <a-col :span="12">
              <a-form-item :label="$t('dataModel.editData.dataAlarm')">
                <a-tooltip placement="top">
                  <template slot="title">
                    <span>{{$t('dataModel.editData.dataAlarmTips')}}</span>
                  </template>
                  <a-select  @change="alarmCharge"   autocomplete="autocomplete"  v-decorator="[
                          'dataAlarm',
                          {
                           initialValue:'0', rules: [{ required: true, message: $t('dataModel.editData.dataAlarm') }],
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
              <a-col :span="12">
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

              <a-col :span="12">
                <a-form-item :label="$t('dataModel.editData.dataRecord')">
                  <a-select   @change="recordCharge" autocomplete="autocomplete"  v-decorator="[
                      'dataRecord',
                      {
                       initialValue:'0', rules: [{ required: true, message: $t('dataModel.editData.dataRecord') }],
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
          <a-row :gutter="16">
            <a-col :span="24">
              <a-form-item
                  :label="$t('dataModel.opcuaModel.NodeIDDec')"
              >
                <a-textarea autocomplete="autocomplete"
                            v-decorator="['NodeIDDec', {rules: [{ required: false, message: $t('dataModel.opcuaModel.NodeIDDec'), whitespace: true}]}]"
                />
              </a-form-item>
            </a-col>
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

        <a-button v-if="isEdit" key="submit" type="primary" :style="{ marginRight: '8px' }" @click="save">
          {{ $t('component.deviceDataModel.Edit')}}
        </a-button>
        <a-button v-else key="submit" type="primary" :style="{ marginRight: '8px' }" @click="AddNodeId">
          {{ $t('component.deviceDataModel.submit')}}
        </a-button>
        <a-button key="back" @click="RegisterVisible=false">
          {{$t('component.deviceDataModel.cancel')}}
        </a-button>

      </div>
    </a-drawer>
  </div>
</template>

<script>
import {
  IEC104ModelNodeIDAdd,
  IEC104ModelNodeIDDel,
  IEC104ModelNodeIDEdit,
  IEC104ModelNodeIDList
} from "@/services/iec104Model";
import {ImportNodeID, LOCALUPGATEDATAMODEL} from "@/services/api";
import { exportExcelWithStyle } from "@/services/excelExport.js"
const dataSource= []
const loadingKey = 'updatable'
export default {
  name: 'IEC104ModelNodeID',
  i18n: require('../../../i18n/language'),
  data () {
    return {
      isEdit:false,
      pagination:{
        pageSize:15,
        showSizeChanger:true
      },
      DataRecordType:0,
      importUrl:ImportNodeID+"/"+this.$route.params.uid,
      registerGroupListTable:true,
      error: '',
      NodeIDAccessLevel:[],
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
      SelectDataCategory:1,
      DataCategoryYaoTiaoType:3,
      RegisterMessageShowLoad:false,
      json_fields_cn: {
        "数据名称": "name",    //常规字段
        "数据类别": {
          field: "DataCategory",
          //自定义回调函数
          callback: value => {
          switch (value){
            case 1:{
              return "遥信"
            }
            case 2:{
              return "遥测"
            }
            case 3:{
              return "脉冲(电量)"
            }
            case 4:{
              return "遥控"
            }
            case 5:{
              return "遥调(设点)"
            }
          }
        }
        },
        "信息点地址": "DataPoint",
        "权限(ReadOnly,ReadWrite)": "auth",
        "类型": {
              field: "type",
              //自定义回调函数
              callback: value => {
                switch (value){
                  case '8':{
                    return "Int"
                  }
                  case '10':{
                    return "Float"
                  }
                }
            }
        },
        "单位": "unit",
        "转换关系": "conversionExpression",
        "是否告警(是,否)": {
          field: "alarm",
          //自定义回调函数
          callback: value => {
            switch (value){
              case 0:{
                return "否"
              }
              case 1:{
                return "是"
              }
            }
          }
        },
        "告警等级(提示、次要、重要、紧急、致命)": {
          field: "alarmLevel",
          //自定义回调函数
          callback: value => {
            switch (value){
              case 0:{
                return "提示"
              }
              case 1:{
                return "次要"
              }
              case 2:{
                return "重要"
              }
              case 3:{
                return "紧急"
              }
              case 4:{
                return "致命"
              }
            }
          }
        },
        "告警消息": "AlarmMessage",
        "告警消除消息": "AlarmClearMessage",
        "是否存储(是,否)": {
          field: "record",
          //自定义回调函数
          callback: value => {
            switch (value){
              case 0:{
                return "否"
              }
              case 1:{
                return "是"
              }
            }
          }
        },
        "存储类型(变化存储、定时存储、即时存储)": {
          field: "RecordType",
          //自定义回调函数
          callback: value => {
            switch (value){
              case 0:{
                return "变化存储"
              }
              case 1:{
                return "定时存储"
              }
              case 2:{
                return "即时存储"
              }
            }
          }
        },
        "定时时间": "recordInterval",
        "变化值": "RecordDataCharge",
        "保留小数": "FloatAccuracy",
        "模型类型(勿修改)": "modeltype",
        "数据ID(勿修改)": "uuid"
      },
      json_meta: [
        [
          {
            " key ": " charset ",
            " value ": " utf- 8 "
          }
        ]
      ],
      exportName:"IEC104数据模板",
      localUpgradeUrl:LOCALUPGATEDATAMODEL+"/"+this.$route.params.uid,
      registerGroupColumns: [
        {
          slotName: this.$t("dataModel.opcuaModel.NodeIDName"),
          scopedSlots: {  customRender: 'name' ,title:this.$t("dataModel.opcuaModel.NodeIDName")},
          width: '15%',
          dataIndex: 'name'
        },
        {
          slotName:this.$t("dataModel.IEC104Model.DataCategory"),
          scopedSlots: {  customRender: 'DataCategory'  ,title:this.$t("dataModel.IEC104Model.DataCategory")},
          width: '15%',
          align:"left",
          dataIndex: 'DataCategory'
        },
        {
          slotName:this.$t("dataModel.IEC104Model.DataFlag"),
          scopedSlots: {  customRender: 'DataPoint'  ,title:this.$t("dataModel.IEC104Model.DataFlag")},
          width: '15%',
          align:"left",
          dataIndex: 'DataPoint'
        },
        {
          slotName:this.$t("dataModel.opcuaModel.NodeIDDataType"),
          scopedSlots: {  customRender: 'NodeIDDataType' ,title:this.$t("dataModel.opcuaModel.NodeIDDataType") },
          width: '15%',
          align:"center",
          dataIndex: 'type',
        },
        {
          slotName:this.$t("dataModel.opcuaModel.NodeIDAccessLevel"),
          scopedSlots: {  customRender: 'NodeIDAccessLevel'  ,title:this.$t("dataModel.opcuaModel.NodeIDAccessLevel")},
          width: '15%',
          align:"center",
          dataIndex: 'auth',
        },
        {
          title: this.$t('dataModel.modelTableOpt'),
          width: '15%',
          scopedSlots: { customRender: 'action' }
        }
      ],
      dataSource,
      registerGroupDataSource:[],
      selectedRows: [],
    }
  },
  created(){
    this.NodeIdList()
    this.registerGroupListTable=true
  },
  activated() {

  },
  mounted() {

  },
  computed: {

  },
  methods: {
    async handleExport() {
      const data = this.registerGroupDataSource.map(item => {
        const row = {};
        for (const key in this.json_fields_cn) {
          const fieldConfig = this.json_fields_cn[key];
          if (typeof fieldConfig === 'string') {
            row[key] = item[fieldConfig];
          } else if (typeof fieldConfig === 'object' && fieldConfig.field) {
            const rawValue = item[fieldConfig.field];
            if (fieldConfig.callback) {
              row[key] = fieldConfig.callback(rawValue);
            } else {
              row[key] = rawValue;
            }
          }
        }
        return row;
      });
      await exportExcelWithStyle(data, this.json_fields_cn, this.exportName, '', false);
    },
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
    chargeDataRecordType(value){
      this.DataRecordType = parseInt(value)
    },
    beforeUpload(){
      this.$message.loading({ content: this.$t('dataModel.opcuaModel.ImportNodeIDLoading'),loadingKey,duration: 0 });
    },
    localUpgradeCharge(info){
      this.dataSource=[]
      this.messageShowLoad = true
      if (info.file.status === 'done') {
        let result = info.file.response
        this.$message.destroy();
        if(result.Code==0) {
          this.$message.success(`${info.file.name} `+this.$t("dataModel.importSuccess"));
          this.NodeIdList()
        }
        else if(result.Code==-2)
        {
          this.$message.error(`${info.file.name} `+this.$t("dataModel.FormatError"));
        }
        else
        {
          this.$message.error(`${info.file.name} `+this.$t("SystemUpgrade.UpgradeFileSaveError"));
        }
      }
      else if (info.file.status === 'uploading') {
        //this.$message.success(`${info.file.name} `+this.$t("SystemUpgrade.BeginUpgradeUploading"));
      }

      this.messageShowLoad = false
    },
    nodeidImport(info) {
      this.dataSource=[]
      let _t = this
      this.messageShowLoad = true
      if (info.file.status === 'done') {
        let result = info.file.response
        if(result.Code==0)
        {
          this.NodeIdList()

          this.$message.success(this.$t('dataModel.opcuaModel.ImportNodeIDSuccess'))
          setTimeout(function (){
            _t.$message.destroy();
          },1000)

        }
        else
        {
          this.$message.error(this.$t('dataModel.opcuaModel.ImportNodeIDFailed'))
          setTimeout(function (){
            _t.$message.destroy();
          },1000)
        }

      }

      this.messageShowLoad = false
    },
    onClose() {
      this.RegisterVisible = false;
    },
    edit(item) {
      let _t = this
      this.isEdit = true
      this.RegisterVisible=true
      this.editingKey = item.uuid
      this.alarmStatus = item.alarm
      this.recordStatus = item.record
      this.ShowRegisterLoading = true
      this.DataRecordType = item.RecordType
      this.SelectDataCategory = item.DataCategory
      this.DataCategoryYaoTiaoType = item.DataCategoryYaoTiaoType
      if(item.recordInterval==0)
      {
        item.recordInterval=1
      }
      setTimeout(function (){
        _t.RegisterForm.setFieldsValue(
            {
              NodeIDName:item.name,
              DataFlag:item.DataPoint.toString(),
              DataCategory:item.DataCategory.toString(),
              NodeIDDataType:item.type,
              NodeIDAccessLevel:item.auth,
              dataUnit:item.unit,
              NodeIDDec:item.Description,
              DataCategoryYaoKongType:item.DataCategoryYaoKongType?item.DataCategoryYaoKongType.toString():'1',
              DataCategoryYaoTiaoType:item.DataCategoryYaoTiaoType?item.DataCategoryYaoTiaoType.toString():'2',
              DataCategoryYaoTiaoGuiYiED:item.DataCategoryYaoTiaoGuiYiED?item.DataCategoryYaoTiaoGuiYiED.toString():'32768',
              ConversionExpression:item.conversionExpression,
            })
        if (item.alarm==1){
          _t.RegisterForm.setFieldsValue(
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
          _t.RegisterForm.setFieldsValue(
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
          _t.RegisterForm.setFieldsValue(
              {
                dataRecord:item.record.toString(),
                dataAlarm:item.alarm.toString(),
              })
        }
        _t.ShowRegisterLoading = false
      },500)
      this.editVisible = true;
    },
    save() {
      this.RegisterForm.validateFields((err) => {
        if (!err) {
          let params = {
            uuid:this.editingKey,
            muid:this.$route.params.uid,
            data: {
              name: this.RegisterForm.getFieldValue('NodeIDName'),
              DataPoint:parseInt(this.RegisterForm.getFieldValue('DataFlag')),
              DataCategory:parseInt(this.RegisterForm.getFieldValue('DataCategory')),
              DataCategoryYaoKongType:parseInt(this.RegisterForm.getFieldValue('DataCategoryYaoKongType'))?parseInt(this.RegisterForm.getFieldValue('DataCategoryYaoKongType')):0,
              DataCategoryYaoTiaoType:parseInt(this.RegisterForm.getFieldValue('DataCategoryYaoTiaoType'))?parseInt(this.RegisterForm.getFieldValue('DataCategoryYaoTiaoType')):0,
              DataCategoryYaoTiaoGuiYiED:parseInt(this.RegisterForm.getFieldValue('DataCategoryYaoTiaoGuiYiED'))?parseInt(this.RegisterForm.getFieldValue('DataCategoryYaoTiaoGuiYiED')):0,
              auth: this.RegisterForm.getFieldValue('NodeIDAccessLevel'),
              type: this.RegisterForm.getFieldValue('NodeIDDataType'),
              unit: this.RegisterForm.getFieldValue('dataUnit'),
              conversionExpression: this.RegisterForm.getFieldValue('ConversionExpression'),
              alarm: parseInt(this.RegisterForm.getFieldValue('dataAlarm')),
              record: parseInt(this.RegisterForm.getFieldValue('dataRecord')),
              Description: this.RegisterForm.getFieldValue('NodeIDDec'),
            }
          }
          if (params.data.alarm==1)
          {
            params.data.alarmLevel= parseInt(this.RegisterForm.getFieldValue('AlarmLevel'))
            params.data.AlarmMessage= this.RegisterForm.getFieldValue('AlarmMessage')
            params.data.AlarmClearMessage=this.RegisterForm.getFieldValue('AlarmClearMessage')
          }
          if (params.data.record==1)
          {
            params.data.recordInterval=  parseInt(this.RegisterForm.getFieldValue('dataRecordTime'))
            params.data.RecordType=parseInt(this.RegisterForm.getFieldValue('dataRecordType'))
            params.data.recordInterval=parseInt(this.RegisterForm.getFieldValue('dataRecordTime'))
            params.data.RecordDataCharge=this.RegisterForm.getFieldValue('dataRecordChargeValue')?this.RegisterForm.getFieldValue('dataRecordChargeValue').toString():""
          }
          let _t = this
          IEC104ModelNodeIDEdit(params).then(function (res){
            if(res.data.code==2002)
            {
              const newData = [..._t.registerGroupDataSource];
              const target = newData.filter(item => _t.editingKey === item.uuid)[0];
              if (target) {
                target.name = _t.RegisterForm.getFieldValue('NodeIDName')
                target.DataPoint = _t.RegisterForm.getFieldValue('DataFlag')
                target.DataCategory = _t.RegisterForm.getFieldValue('DataCategory')
                target.DataCategoryYaoKongType = _t.RegisterForm.getFieldValue('DataCategoryYaoKongType')
                target.DataCategoryYaoTiaoType = _t.RegisterForm.getFieldValue('DataCategoryYaoTiaoType')
                target.DataCategoryYaoTiaoGuiYiED=_t.RegisterForm.getFieldValue('DataCategoryYaoTiaoGuiYiED')
                target.auth=_t.RegisterForm.getFieldValue('NodeIDAccessLevel')
                target.type=_t.RegisterForm.getFieldValue('NodeIDDataType')
                target.unit=_t.RegisterForm.getFieldValue('dataUnit')
                target.conversionExpression=_t.RegisterForm.getFieldValue('ConversionExpression')
                target.alarm=parseInt(_t.RegisterForm.getFieldValue('dataAlarm'))
                target.alarmLevel=parseInt(_t.RegisterForm.getFieldValue('AlarmLevel'))
                target.AlarmMessage = _t.RegisterForm.getFieldValue('AlarmMessage')
                target.AlarmClearMessage = _t.RegisterForm.getFieldValue('AlarmClearMessage')
                target.record=parseInt(_t.RegisterForm.getFieldValue('dataRecord'))
                target.RecordType=parseInt(_t.RegisterForm.getFieldValue('dataRecordType'))
                target.recordInterval=parseInt(_t.RegisterForm.getFieldValue('dataRecordTime'))
                target.RecordDataCharge=_t.RegisterForm.getFieldValue('dataRecordChargeValue')?_t.RegisterForm.getFieldValue('dataRecordChargeValue').toString():""
                _t.registerGroupDataSource = newData;
              }
              _t.$message.success(_t.$t("dataModel.saveSuccess"));
              _t.RegisterVisible = false;
            }
            else
            {
              _t.$message.error(_t.$t("dataModel.saveFailed"));
            }
          })
        }
      })
    },
    NodeIdList(){
      this.messageShowLoad = true
      const params = {
        muid:this.$route.params.uid,
      }
      this.RegisterVisible = false;
      let _t = this
      _t.registerGroupDataSource = []
      IEC104ModelNodeIDList(params).then(function (res){
        _t.messageShowLoad = false
        if(res.data.code==0)
        {
          _t.registerGroupDataSource = res.data.list
        }
      }).catch(function (){
        _t.messageShowLoad = false
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
    },
    AddNodeId(){
      this.RegisterForm.validateFields((err) => {
        if (!err) {
          this.messageShowLoad = true
          const params = {
            muid:this.$route.params.uid,
            modeltype:40,
            name:this.RegisterForm.getFieldValue('NodeIDName'),
            DataPoint:parseInt(this.RegisterForm.getFieldValue('DataFlag')),
            DataCategory:parseInt(this.RegisterForm.getFieldValue('DataCategory')),
            DataCategoryYaoKongType:parseInt(this.RegisterForm.getFieldValue('DataCategoryYaoKongType'))?parseInt(this.RegisterForm.getFieldValue('DataCategoryYaoKongType')):0,
            DataCategoryYaoTiaoType:parseInt(this.RegisterForm.getFieldValue('DataCategoryYaoTiaoType'))?parseInt(this.RegisterForm.getFieldValue('DataCategoryYaoTiaoType')):0,
            DataCategoryYaoTiaoGuiYiED:parseInt(this.RegisterForm.getFieldValue('DataCategoryYaoTiaoGuiYiED'))?parseInt(this.RegisterForm.getFieldValue('DataCategoryYaoTiaoGuiYiED')):0,
            auth:this.RegisterForm.getFieldValue('NodeIDAccessLevel'),
            type:this.RegisterForm.getFieldValue('NodeIDDataType'),
            unit:this.RegisterForm.getFieldValue('dataUnit'),
            conversionExpression:this.RegisterForm.getFieldValue('ConversionExpression'),
            alarm:parseInt(this.RegisterForm.getFieldValue('dataAlarm')),
            alarmLevel:parseInt(this.RegisterForm.getFieldValue('AlarmLevel')),
            AlarmMessage:this.RegisterForm.getFieldValue('AlarmMessage'),
            AlarmClearMessage:this.RegisterForm.getFieldValue('AlarmClearMessage'),
            record:parseInt(this.RegisterForm.getFieldValue('dataRecord')),
            RecordType:this.RegisterForm.getFieldValue('dataRecordType')?parseInt(this.RegisterForm.getFieldValue('dataRecordType')):0,
            recordInterval:this.RegisterForm.getFieldValue('dataRecordTime')?parseInt(this.RegisterForm.getFieldValue('dataRecordTime')):0,
            RecordDataCharge:this.RegisterForm.getFieldValue('dataRecordChargeValue')?this.RegisterForm.getFieldValue('dataRecordChargeValue').toString():"",

            Description:this.RegisterForm.getFieldValue('NodeIDDec'),
          }
          this.RegisterVisible = false;
          let _t = this
          IEC104ModelNodeIDAdd(params).then(function (res){
            _t.messageShowLoad = false
            if(res.data.code==2002)
            {
              _t.NodeIdList()
              _t.$message.success(_t.$t("dataModel.saveSuccess"));
            }
            else if(res.data.code==2001)
            {
              _t.$message.error(_t.$t("dataModel.IEC104Model.IsErrorExist"));
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
    deleteRecord(uuid,muid) {
      let _t = this
      let params={
        uuid:uuid,
        muid:muid
      }
      IEC104ModelNodeIDDel(params).then(function (res) {
        if(res.data.code==200)
        {
          _t.$message.success(_t.$t("dataModel.deleteSuccess"));
          _t.NodeIdList()
        }
        else {
          _t.$message.error(_t.$t("dataModel.deleteFailed"));
        }
      })
    },
    onBlackCLK(){
      this.$router.push('/DeviceModel/IEC104Model')
    },
    alarmCharge(value){
      this.alarmStatus=parseInt(value)
      let _t = this
      this.$nextTick(function(){
        _t.RegisterForm.setFieldsValue(
            {
              AlarmLevel:"0",
            })
      });
    },
    DataCategoryAlarmCharge(value){
      this.SelectDataCategory=parseInt(value)
    },
    DataCategoryYaoTiaoGuiYiCharge(value){
      this.DataCategoryYaoTiaoType=parseInt(value)
    },
    recordCharge(value){
      this.recordStatus=parseInt(value)
    },
  }
}
</script>


<style lang="less" >
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

.DataType::-webkit-scrollbar {/*滚动条整体样式*/
  width:4px;/*高宽分别对应横竖滚动条的尺寸*/
  height:4px;
}

.DataType::-webkit-scrollbar-thumb {/*滚动条里面小方块*/
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

.DataType::-webkit-scrollbar-track {/*滚动条里面轨道*/
  /*滚动条里面轨道*/
  box-shadow   : inset 0 0 5px rgba(0, 0, 0, 0.2);
  background   : #ededed;
  border-radius: 10px;
}

</style>