<template>
  <div>
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
        <a-table  :pagination="pagination" rowKey="uuid" :columns="registerGroupColumns" :data-source="registerGroupDataSource" class="ant-table-tbody">
          <template v-for="(item, index) in registerGroupColumns" :slot="item.slotName">
            <span :key="index">{{ $t(item.slotName) }}</span>
          </template>
          <template slot="PointType" slot-scope="text">
            <span v-if="text=='BinaryInput'"> {{$t('dataModel.DNP3Model.PointTypeBinaryInput')}}</span>
            <span v-else-if="text=='BinaryOutput'"> {{$t('dataModel.DNP3Model.PointTypeBinaryOutput')}}</span>
            <span v-else-if="text=='CounterInput'"> {{$t('dataModel.DNP3Model.PointTypeCounterInput')}}</span>
            <span v-else-if="text=='AnalogInput'"> {{$t('dataModel.DNP3Model.PointTypeAnalogInput')}}</span>
            <span v-else-if="text=='AnalogOutput'"> {{$t('dataModel.DNP3Model.PointTypeAnalogOutput')}}</span>
            <span v-else-if="text=='DoubleBitBinaryInput'"> {{$t('dataModel.DNP3Model.PointTypeDoubleBitBinaryInput')}}</span>
            <span v-else>{{ text }}</span>
          </template>
          <template slot="Class" slot-scope="text">
            <span v-if="text==0"> {{$t('dataModel.DNP3Model.ClassNone')}}</span>
            <span v-else-if="text==1"> {{$t('dataModel.DNP3Model.Class1')}}</span>
            <span v-else-if="text==2"> {{$t('dataModel.DNP3Model.Class2')}}</span>
            <span v-else-if="text==3"> {{$t('dataModel.DNP3Model.Class3')}}</span>
          </template>
          <template slot="NodeIDDataType" slot-scope="text">
            <span v-if="text=='1'"> Boolean</span>
            <span v-else-if="text=='2'"> Byte</span>
            <span v-else-if="text=='3'"> Short</span>
            <span v-else-if="text=='4'"> UInt16</span>
            <span v-else-if="text=='5'"> Int16</span>
            <span v-else-if="text=='6'"> UInt32</span>
            <span v-else-if="text=='7'"> Int32</span>
            <span v-else-if="text=='8'"> UInt64</span>
            <span v-else-if="text=='9'"> Int64</span>
            <span v-else-if="text=='10'"> Float</span>
            <span v-else-if="text=='11'"> Double</span>
            <span v-else-if="text=='12'"> String</span>
            <span v-else>{{ text }}</span>
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
                  :label="$t('dataModel.DNP3Model.PointType')"
              >
                <a-select class="DataType" autocomplete="autocomplete"
                          v-decorator="['PointType', {initialValue:'BinaryInput',rules: [{ required: true, message: $t('dataModel.DNP3Model.PointType'), whitespace: true}]}]"
                >
                  <a-select-option value="BinaryInput">{{$t('dataModel.DNP3Model.PointTypeBinaryInput')}}</a-select-option>
                  <a-select-option value="BinaryOutput">{{$t('dataModel.DNP3Model.PointTypeBinaryOutput')}}</a-select-option>
                  <a-select-option value="CounterInput">{{$t('dataModel.DNP3Model.PointTypeCounterInput')}}</a-select-option>
                  <a-select-option value="AnalogInput">{{$t('dataModel.DNP3Model.PointTypeAnalogInput')}}</a-select-option>
                  <a-select-option value="AnalogOutput">{{$t('dataModel.DNP3Model.PointTypeAnalogOutput')}}</a-select-option>
                  <a-select-option value="DoubleBitBinaryInput">{{$t('dataModel.DNP3Model.PointTypeDoubleBitBinaryInput')}}</a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item
                  :label="$t('dataModel.DNP3Model.Group')"
              >
                <a-input-number style="width: 100%" autocomplete="autocomplete"
                         v-decorator="['Group', {initialValue:1,rules: [{ required: true, message: $t('dataModel.DNP3Model.Group')}]}]"
                />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item
                  :label="$t('dataModel.DNP3Model.Variation')"
              >
                <a-input-number style="width: 100%" autocomplete="autocomplete"
                         v-decorator="['Variation', {initialValue:1,rules: [{ required: true, message: $t('dataModel.DNP3Model.Variation')}]}]"
                />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item
                  :label="$t('dataModel.DNP3Model.Index')"
              >
                <a-input-number style="width: 100%" :min="0" autocomplete="autocomplete"
                         v-decorator="['Index', {initialValue:0,rules: [{ required: true, message: $t('dataModel.DNP3Model.Index')}]}]"
                />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item
                  :label="$t('dataModel.DNP3Model.Class')"
              >
                <a-select class="DataType" autocomplete="autocomplete"
                          v-decorator="['Class', {initialValue:0,rules: [{ required: true, message: $t('dataModel.DNP3Model.Class')}]}]"
                >
                  <a-select-option :value="0">{{$t('dataModel.DNP3Model.ClassNone')}}</a-select-option>
                  <a-select-option :value="1">{{$t('dataModel.DNP3Model.Class1')}}</a-select-option>
                  <a-select-option :value="2">{{$t('dataModel.DNP3Model.Class2')}}</a-select-option>
                  <a-select-option :value="3">{{$t('dataModel.DNP3Model.Class3')}}</a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item
                  :label="$t('dataModel.opcuaModel.NodeIDDataType')"
              >
                <a-select class="DataType" autocomplete="autocomplete"
                          v-decorator="['NodeIDDataType', {initialValue:'10',rules: [{ required: true, message: $t('dataModel.opcuaModel.NodeIDDataType'), whitespace: true}]}]"
                >
                  <a-select-option value="1">Boolean</a-select-option>
                  <a-select-option value="2">Byte</a-select-option>
                  <a-select-option value="3">Short</a-select-option>
                  <a-select-option value="4">UInt16</a-select-option>
                  <a-select-option value="5">Int16</a-select-option>
                  <a-select-option value="6">UInt32</a-select-option>
                  <a-select-option value="7">Int32</a-select-option>
                  <a-select-option value="8">UInt64</a-select-option>
                  <a-select-option value="9">Int64</a-select-option>
                  <a-select-option value="10">Float</a-select-option>
                  <a-select-option value="11">Double</a-select-option>
                  <a-select-option value="12">String</a-select-option>
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
  DNP3ModelNodeIDAdd,
  DNP3ModelNodeIDDel,
  DNP3ModelNodeIDEdit,
  DNP3ModelNodeIDList
} from "@/services/dnp3Model";
import {ImportNodeID, LOCALUPGATEDATAMODEL} from "@/services/api";
import { exportExcelWithStyle } from "@/services/excelExport.js"

export default {
  name: 'DNP3ModelNodeID',
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
      error: '',
      ShowRegisterLoading:false,
      alarmStatus:0,
      recordStatus:0,
      RegisterForm: this.$form.createForm(this),
      RegisterVisible:false,
      messageShowLoad:false,
      localUpgradeUrl:LOCALUPGATEDATAMODEL,
      registerGroupColumns: [
        {
          width: '5%',
          slotName: 'dataModel.modelTableIndex',
          scopedSlots: { customRender: 'serial', title: 'dataModel.modelTableIndex' },
          customRender: (text, record, index) => `${index + 1}`,
        },
        {
          width: '12%',
          slotName: 'dataModel.opcuaModel.NodeIDName',
          scopedSlots: { customRender: 'name', title: 'dataModel.opcuaModel.NodeIDName' },
          dataIndex: 'name',
        },
        {
          width: '10%',
          slotName: 'dataModel.DNP3Model.PointType',
          scopedSlots: { customRender: 'PointType', title: 'dataModel.DNP3Model.PointType' },
          dataIndex: 'PointType',
        },
        {
          width: '8%',
          slotName: 'dataModel.DNP3Model.Group',
          scopedSlots: { customRender: 'Group', title: 'dataModel.DNP3Model.Group' },
          dataIndex: 'Group',
        },
        {
          width: '8%',
          slotName: 'dataModel.DNP3Model.Variation',
          scopedSlots: { customRender: 'Variation', title: 'dataModel.DNP3Model.Variation' },
          dataIndex: 'Variation',
        },
        {
          width: '8%',
          slotName: 'dataModel.DNP3Model.Index',
          scopedSlots: { customRender: 'Index', title: 'dataModel.DNP3Model.Index' },
          dataIndex: 'Index',
        },
        {
          width: '8%',
          slotName: 'dataModel.DNP3Model.Class',
          scopedSlots: { customRender: 'Class', title: 'dataModel.DNP3Model.Class' },
          dataIndex: 'Class',
        },
        {
          width: '8%',
          slotName: 'dataModel.opcuaModel.NodeIDDataType',
          scopedSlots: { customRender: 'NodeIDDataType', title: 'dataModel.opcuaModel.NodeIDDataType' },
          dataIndex: 'type',
        },
        {
          width: '8%',
          slotName: 'dataModel.opcuaModel.NodeIDAccessLevel',
          scopedSlots: { customRender: 'auth', title: 'dataModel.opcuaModel.NodeIDAccessLevel' },
          dataIndex: 'auth',
        },
        {
          width: '8%',
          slotName: 'dataModel.editData.dataUnit',
          scopedSlots: { customRender: 'unit', title: 'dataModel.editData.dataUnit' },
          dataIndex: 'unit',
        },
        {
          slotName: 'dataModel.modelTableOpt',
          scopedSlots: { customRender: 'action',title: 'dataModel.modelTableOpt'}
        }
      ],
      registerGroupDataSource: [],
    }
  },
  created(){
    this.NodeIdList()
  },
  methods: {
    isSpec(str){
      let specialRegex = /[!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?]+/;
      return specialRegex.test(str);
    },
    isValidateTxtNonSpec (rule, value, callback) {
      if (value != null && value !== '') {
        let numStr = value.charAt(0);
        if ((this.isSpec(value)) || (value.indexOf(' ') !== -1)||(!isNaN(parseFloat(numStr)) && isFinite(numStr))) {
          callback(new Error(this.$t('device.deviceNameVal')))
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
      this.$message.loading({ content: this.$t('dataModel.opcuaModel.ImportNodeIDLoading'), duration: 0 });
    },
    localUpgradeCharge(info){
      this.registerGroupDataSource=[]
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
      this.messageShowLoad = false
    },
    handleExport(){
      let _t = this
      let exportData = []
      for(let i=0;i<this.registerGroupDataSource.length;i++){
        let item = this.registerGroupDataSource[i]
        let pointTypeLabel = item.PointType
        let classLabel = item.Class
        let typeLabel = item.type
        let authLabel = item.auth
        exportData.push({
          no: i+1,
          name: item.name,
          PointType: pointTypeLabel,
          Group: item.Group,
          Variation: item.Variation,
          Index: item.Index,
          Class: classLabel,
          type: typeLabel,
          auth: authLabel,
          unit: item.unit,
          conversionExpression: item.conversionExpression,
        })
      }
      let columns = [
        { header: this.$t('dataModel.modelTableIndex'), key: 'no', width: 8 },
        { header: this.$t('dataModel.opcuaModel.NodeIDName'), key: 'name', width: 20 },
        { header: this.$t('dataModel.DNP3Model.PointType'), key: 'PointType', width: 18 },
        { header: this.$t('dataModel.DNP3Model.Group'), key: 'Group', width: 10 },
        { header: this.$t('dataModel.DNP3Model.Variation'), key: 'Variation', width: 12 },
        { header: this.$t('dataModel.DNP3Model.Index'), key: 'Index', width: 10 },
        { header: this.$t('dataModel.DNP3Model.Class'), key: 'Class', width: 10 },
        { header: this.$t('dataModel.opcuaModel.NodeIDDataType'), key: 'type', width: 12 },
        { header: this.$t('dataModel.opcuaModel.NodeIDAccessLevel'), key: 'auth', width: 14 },
        { header: this.$t('dataModel.editData.dataUnit'), key: 'unit', width: 10 },
        { header: this.$t('dataModel.modbusModel.ConversionExpression'), key: 'conversionExpression', width: 18 },
      ]
      exportExcelWithStyle(exportData, columns, 'DNP3DataModel')
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
      if(item.recordInterval==0)
      {
        item.recordInterval=1
      }
      setTimeout(function (){
        _t.RegisterForm.setFieldsValue(
            {
              NodeIDName:item.name,
              PointType:item.PointType,
              Group:item.Group,
              Variation:item.Variation,
              Index:item.Index,
              Class:item.Class,
              NodeIDDataType:item.type,
              NodeIDAccessLevel:item.auth,
              dataUnit:item.unit,
              NodeIDDec:item.Description,
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
    },
    save() {
      this.RegisterForm.validateFields((err) => {
        if (!err) {
          let params = {
            uuid:this.editingKey,
            muid:this.$route.params.uid,
            data: {
              name: this.RegisterForm.getFieldValue('NodeIDName'),
              PointType: this.RegisterForm.getFieldValue('PointType'),
              Group: this.RegisterForm.getFieldValue('Group'),
              Variation: this.RegisterForm.getFieldValue('Variation'),
              Index: this.RegisterForm.getFieldValue('Index'),
              Class: this.RegisterForm.getFieldValue('Class'),
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
          DNP3ModelNodeIDEdit(params).then(function (res){
            if(res.data.code==2002)
            {
              _t.$message.success(_t.$t("dataModel.saveSuccess"));
              _t.RegisterVisible = false;
              _t.NodeIdList()
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
      DNP3ModelNodeIDList(params).then(function (res){
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
            modeltype:510,
            name:this.RegisterForm.getFieldValue('NodeIDName'),
            PointType:this.RegisterForm.getFieldValue('PointType'),
            Group:this.RegisterForm.getFieldValue('Group'),
            Variation:this.RegisterForm.getFieldValue('Variation'),
            Index:this.RegisterForm.getFieldValue('Index'),
            Class:this.RegisterForm.getFieldValue('Class'),
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
          DNP3ModelNodeIDAdd(params).then(function (res){
            _t.messageShowLoad = false
            if(res.data.code==2002)
            {
              _t.NodeIdList()
              _t.$message.success(_t.$t("dataModel.saveSuccess"));
            }
            else if(res.data.code==2001)
            {
              _t.$message.error(_t.$t("dataModel.DNP3Model.IsErrorExist"));
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
      DNP3ModelNodeIDDel(params).then(function (res) {
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
      this.$router.push('/DeviceModel/DNP3Model')
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
  transition: background .3s ease;
}

.DataType::-webkit-scrollbar {
  width:4px;
  height:4px;
}

.DataType::-webkit-scrollbar-thumb {
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

.DataType::-webkit-scrollbar-track {
  box-shadow   : inset 0 0 5px rgba(0, 0, 0, 0.2);
  background   : #ededed;
  border-radius: 10px;
}
</style>