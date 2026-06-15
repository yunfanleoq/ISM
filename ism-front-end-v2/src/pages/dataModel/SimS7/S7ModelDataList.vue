<template>
  <div>
    <!--    寄存器组表格-->
    <a-card  style="min-height: 400px">
      <a-space class="operator">

        <a-button type="primary" @click="RegisterVisible=true;isEdit=false"> <a-icon type="plus" />
          {{$t("dataModel.SimS7Model.AddData")}}</a-button>

        <a-button type="default" @click="onBlackCLK()"> <a-icon type="backward" />
          {{$t("dataModel.opcuaModel.Back")}}</a-button>


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

      </a-space>
      <a-spin style="padding: 1px;"  :spinning="messageShowLoad" tip="Loading...">
        <a-table  :pagination="pagination" rowKey="name" :columns="registerGroupColumns" :data-source="registerGroupDataSource" class="ant-table-tbody">
          <template v-for="(item, index) in registerGroupColumns" :slot="item.slotName">
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
          <template slot="name" slot-scope="text, record,index, column">
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
          <template slot="DataFrom" slot-scope="text">
            <span v-if="text==0"> {{$t('dataModel.SimS7Model.DB')}}</span>
            <span v-if="text==1"> {{$t('dataModel.SimS7Model.InputArea')}}</span>
            <span v-if="text==2"> {{$t('dataModel.SimS7Model.OutputArea')}}</span>
            <span v-if="text==3"> {{$t('dataModel.SimS7Model.MemoryArea')}}</span>
          </template>
          <template slot="NodeIDDataType" slot-scope="text">
            <span v-if="text==0"> Bool</span>
            <span v-else-if="text==1"> Byte</span>
            <span v-else-if="text==2"> SINT </span>
            <span v-else-if="text==3"> INT</span>
            <span v-else-if="text==4"> DINT</span>
            <span v-else-if="text==5"> USINT</span>
            <span v-else-if="text==6"> UINT</span>
            <span v-else-if="text==7"> UDINT</span>
            <span v-else-if="text==8"> LINT</span>
            <span v-else-if="text==9"> REAL</span>
            <span v-else-if="text==10"> LREAL</span>
            <span v-else-if="text==11">TIME</span>
            <span v-else-if="text==12">LTIME</span>
            <span v-else-if="text==13">S5TIME</span>
            <span v-else-if="text==14">DATE</span>
            <span v-else-if="text==15">DATE_AND_TIME</span>
            <span v-else-if="text==16">WSTRING</span>
            <span v-else-if="text==17">String</span>
            <span v-else-if="text==18">ULINT</span>
            <span v-else-if="text==19">WORD</span>
            <span v-else-if="text==20">DWORD</span>
            <span v-else-if="text==21">LWORD</span>
            <span v-else-if="text==22">CHAR</span>
            <span v-else-if="text==23">WCHAR</span>
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
                         v-decorator="['NodeIDName', {rules: [{ required: true,validator: isValidateTxtNonSpec, message: $t('device.deviceNameVal'), whitespace: true}]}]"
                />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item
                  :label="$t('dataModel.SimS7Model.DataFrom')"
              >
                <a-select autocomplete="autocomplete" @change="DataFromCharge"
                         v-decorator="['DataFromType', { initialValue:0,rules: [{ required: true,type: 'number', message: $t('dataModel.SimS7Model.DataFrom'), whitespace: true}]}]"
                >
                  <a-select-option v-for="DataFromList in [{name:'dataModel.SimS7Model.DB',value:0},{name:'dataModel.SimS7Model.InputArea',value:1},{name:'dataModel.SimS7Model.OutputArea',value:2},{name:'dataModel.SimS7Model.MemoryArea',value:3}]" :key="DataFromList.value" :value="DataFromList.value">
                    {{ $t(DataFromList.name) }}
                  </a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :span="12" v-if="DataFrom==0">
              <a-form-item
                  :label="$t('dataModel.SimS7Model.DBIndex')"
              >
                <a-input  autocomplete="autocomplete"
                          v-decorator="['DBIndex', {rules: [{ required: true, message: $t('dataModel.SimS7Model.DBIndex'), whitespace: true}]}]"
                >
                </a-input>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item
                  :label="DataFrom==0?$t('dataModel.SimS7Model.DBOffset'):$t('dataModel.SimS7Model.DBAddress')"
              >
                <a-input   autocomplete="autocomplete"
                           v-decorator="['DBOffset', {rules: [{ required: true, message: DataFrom==0?$t('dataModel.SimS7Model.DBOffset'):$t('dataModel.SimS7Model.DBAddress'), whitespace: true}]}]"
                >

                </a-input>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item
                  :label="$t('dataModel.opcuaModel.NodeIDDataType')"
              >
                <a-select class="DataType" autocomplete="autocomplete" @change="DataTypeCharge"
                          v-decorator="['NodeIDDataType', {rules: [{ required: true, message: $t('dataModel.opcuaModel.NodeIDDataType'), whitespace: true}]}]"
                >

<!--                  <a-select-option value="0" v-if="DataFrom==0||DataFrom==1||DataFrom==2||DataFrom==3">Bool</a-select-option>-->
<!--                  <a-select-option value="1" v-if="DataFrom==0||DataFrom==1||DataFrom==2||DataFrom==3">Byte</a-select-option>-->
<!--                  <a-select-option value="2" v-if="DataFrom==0||DataFrom==3">SINT</a-select-option>-->
<!--                  <a-select-option value="3" v-if="DataFrom==0||DataFrom==3">INT</a-select-option>-->
<!--                  <a-select-option value="4" v-if="DataFrom==0||DataFrom==3">DINT</a-select-option>-->
<!--                  <a-select-option value="5" v-if="DataFrom==0||DataFrom==3">USINT</a-select-option>-->
<!--                  <a-select-option value="6" v-if="DataFrom==0||DataFrom==3">UINT</a-select-option>-->
<!--                  <a-select-option value="7" v-if="DataFrom==0||DataFrom==3">UDINT</a-select-option>-->
<!--                  <a-select-option value="8" v-if="DataFrom==0||DataFrom==3">LINT</a-select-option>-->
<!--                  <a-select-option value="18" v-if="DataFrom==0||DataFrom==3">ULINT</a-select-option>-->
<!--                  <a-select-option value="19" v-if="DataFrom==0||DataFrom==3">WORD</a-select-option>-->
<!--                  <a-select-option value="20" v-if="DataFrom==0||DataFrom==3">DWORD</a-select-option>-->
<!--                  <a-select-option value="21" v-if="DataFrom==0||DataFrom==3">LWORD</a-select-option>-->
<!--                  <a-select-option value="9" v-if="DataFrom==0||DataFrom==3">REAL</a-select-option>-->
<!--                  <a-select-option value="10" v-if="DataFrom==0||DataFrom==3">LREAL</a-select-option>-->

<!--                  <a-select-option value="16" v-if="DataFrom==0">WSTRING</a-select-option>-->
<!--                  <a-select-option value="17" v-if="DataFrom==0">STRING</a-select-option>-->

<!--                  <a-select-option value="11" v-if="DataFrom==0||DataFrom==3">TIME</a-select-option>-->
<!--                  <a-select-option value="12" v-if="DataFrom==0||DataFrom==3">LTIME</a-select-option>-->
<!--                  <a-select-option value="13" v-if="DataFrom==0||DataFrom==3">S5TIME</a-select-option>-->
<!--                  <a-select-option value="14" v-if="DataFrom==0">DATE</a-select-option>-->
<!--                  <a-select-option value="15" v-if="DataFrom==0">DATE_AND_TIME</a-select-option>-->
<!--                  <a-select-option value="22" v-if="DataFrom==0||DataFrom==3">CHAR</a-select-option>-->
<!--                  <a-select-option value="23" v-if="DataFrom==0||DataFrom==3">WCHAR</a-select-option>-->

                  <a-select-option value="0" >Bool</a-select-option>
                  <a-select-option value="1" >Byte</a-select-option>
                  <a-select-option value="2" >SINT</a-select-option>
                  <a-select-option value="3" >INT</a-select-option>
                  <a-select-option value="4">DINT</a-select-option>
                  <a-select-option value="5" >USINT</a-select-option>
                  <a-select-option value="6" >UINT</a-select-option>
                  <a-select-option value="7" >UDINT</a-select-option>
                  <a-select-option value="8" >LINT</a-select-option>
                  <a-select-option value="18" >ULINT</a-select-option>
                  <a-select-option value="19" >WORD</a-select-option>
                  <a-select-option value="20" >DWORD</a-select-option>
                  <a-select-option value="21" >LWORD</a-select-option>
                  <a-select-option value="9" >REAL</a-select-option>
                  <a-select-option value="10" >LREAL</a-select-option>

                  <a-select-option value="16" >WSTRING</a-select-option>
                  <a-select-option value="17" >STRING</a-select-option>

                  <a-select-option value="11" >TIME</a-select-option>
                  <a-select-option value="12" >LTIME</a-select-option>
                  <a-select-option value="13" >S5TIME</a-select-option>
                  <a-select-option value="14" >DATE</a-select-option>
                  <a-select-option value="15" >DATE_AND_TIME</a-select-option>
                  <a-select-option value="22" >CHAR</a-select-option>
                  <a-select-option value="23">WCHAR</a-select-option>

                </a-select>
              </a-form-item>
            </a-col>
            <a-col :span="12" v-if="S7DataType=='16'||S7DataType=='17'">
              <a-form-item :label="$t('dataModel.SimS7Model.StringLength')">
                <a-tooltip placement="top">
                  <template slot="title">
                    <span>{{$t('dataModel.SimS7Model.StringLength')}}</span>
                  </template>
                  <a-input   autocomplete="autocomplete"   v-decorator="[
                        'StringLength',
                        {
                          initialValue:'254',
                          rules: [{ required: true, message: $t('dataModel.SimS7Model.StringLength') }],
                        },
                      ]">
                  </a-input>
                </a-tooltip>
              </a-form-item>
            </a-col>
            <a-col :span="12" v-if="S7DataType=='20'||S7DataType=='21'||S7DataType=='1'||S7DataType=='19'">
              <a-form-item :label="$t('dataModel.SimS7Model.IsHaveUnsigned')">
                <a-tooltip placement="top">
                  <template slot="title">
                    <span>{{$t('dataModel.SimS7Model.IsHaveUnsigned')}}</span>
                  </template>
                  <a-select   autocomplete="autocomplete"   v-decorator="[
                        'IsHaveUnsigned',
                        {
                          initialValue:'0',
                          rules: [{ required: true,message: $t('dataModel.SimS7Model.IsHaveUnsigned') }],
                        },
                      ]">
                    <a-select-option value='0' >{{$t('dataModel.SimS7Model.IsSigned')}}</a-select-option>
                    <a-select-option value='1' >{{$t('dataModel.SimS7Model.IsUnsigned')}}</a-select-option>
                  </a-select>
                </a-tooltip>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item
                  :label="$t('dataModel.opcuaModel.NodeIDAccessLevel')"
              >
                <a-select  class="DataType" autocomplete="autocomplete"
                           v-decorator="['NodeIDAccessLevel', {rules: [{ required: true, message: $t('dataModel.opcuaModel.NodeIDAccessLevel'), whitespace: true}]}]"
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
  SimS7ModelDataAdd,
  SimS7ModelDataDel,
  SimS7ModelDataEdit,
  SimS7ModelDataList
} from "@/services/SimS7";
import {LOCALUPGATEDATAMODEL} from "@/services/api";
import { exportExcelWithStyle } from "@/services/excelExport.js"
const dataSource= []
const loadingKey = 'updatable'
export default {
  name: 'SimS7DataList',
  i18n: require('../../../i18n/language'),
  data () {
    return {
      isEdit:false,
      pagination:{
        pageSize:15,
        showSizeChanger:true
      },
      localUpgradeUrl:LOCALUPGATEDATAMODEL+"/"+this.$route.params.uid,
      json_fields_cn: {
        "数据名称": "name",    //常规字段
        "数据来源": {
          field: "DataFromType",
              //自定义回调函数
              callback: value => {
            switch (value){
              case 0:{
                return "DB"
              }
              case 1:{
                return "I区"
              }
              case 2:{
                return "Q区"
              }
              case 3:{
                return "M区"
              }
            }
          }
        },
        "数据块序号": "DBIndex",
        "数据偏移": "DBOffset",
        "字符串长度": "StringMaxLength",
        "是否有符号": {
          field: "IsHaveUnsigned",
          //自定义回调函数
          callback: value => {
            switch (value){
              case 0:{
                return "有符号"
              }
              case 1:{
                return "无符号"
              }
            }
          }
        },
        "权限(ReadOnly,ReadWrite)": "auth",
        "类型": {
            field: "type",
                //自定义回调函数
                callback: value => {
            switch (value){
              case "0":{
                return "Bool"
              }
              case "1":{
                return "Byte"
              }
              case "2":{
                return "SINT"
              }
              case "3":{
                return "INT"
              }
              case "4":{
                return "DINT"
              }
              case "5":{
                return "USINT"
              }
              case "6":{
                return "UINT"
              }
              case "7":{
                return "UDINT"
              }
              case "8":{
                return "LINT"
              }
              case "9":{
                return "REAL"
              }
              case "10":{
                return "LREAL"
              }
              case "11":{
                return "TIME"
              }
              case "12":{
                return "LTIME"
              }
              case "13":{
                return "S5TIME"
              }
              case "14":{
                return "DATE"
              }
              case "15":{
                return "DATE_AND_TIME"
              }
              case "16":{
                return "WSTRING"
              }
              case "17":{
                return "String"
              }
              case "18":{
                return "ULINT"
              }
              case "19":{
                return "WORD"
              }
              case "20":{
                return "DWORD"
              }
              case "21":{
                return "LWORD"
              }
              case "22":{
                return "CHAR"
              }
              case "23":{
                return "WCHAR"
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
        "数据ID(勿修改，留空即可新增)": "uuid",
      },
      json_meta: [
        [
          {
            " key ": " charset ",
            " value ": " utf- 8 "
          }
        ]
      ],
      exportName:"西门子S7数据模板",
      DataFrom:0,
      DataRecordType:0,
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
      searchText: '',
      searchInput: null,
      searchedColumn: '',
      S7DataType:0,
      RegisterVisible:false,
      editVisible:false,
      messageShowLoad:false,
      RegisterMessageShowLoad:false,
      registerGroupColumns: [
        {
          slotName: this.$t("dataModel.opcuaModel.NodeIDName"),
          scopedSlots: {  filterDropdown: 'filterDropdown', filterIcon: 'filterIcon', customRender: 'name' ,title:this.$t("dataModel.opcuaModel.NodeIDName")},
          width: '15%',
          dataIndex: 'name',
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
          slotName:this.$t("dataModel.SimS7Model.DataFrom"),
          scopedSlots: {  customRender: 'DataFrom'  ,title:this.$t("dataModel.SimS7Model.DataFrom")},
          width: '15%',
          align:"left",
          dataIndex: 'DataFromType'
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
    beforeUpload(){
      this.$message.loading({ content: this.$t('dataModel.opcuaModel.ImportNodeIDLoading'),loadingKey,duration: 0 });
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
    localUpgradeCharge(info){
      this.dataSource=[]
      this.messageShowLoad = true
      if (info.file.status === 'done') {
        let result = info.file.response
        this.$message.destroy();
        if(result.Code==0) {
          this.$message.success(`${info.file.name} `+this.$t("dataModel.importSuccess"));
          this.NodeIdList(this.registerUUID)
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
    handleSearch(selectedKeys, confirm, dataIndex) {
      confirm();
      this.searchText = selectedKeys[0];
      this.searchedColumn = dataIndex;
    },
    handleReset(clearFilters) {
      clearFilters();
      this.searchText = '';
    },
    DataFromCharge(e){
      this.DataFrom = e
    },
    DataTypeCharge(e){
      this.S7DataType = e
    },
    chargeDataRecordType(value){
      this.DataRecordType = parseInt(value)
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
      this.DataFrom = item.DataFromType

      if(item.recordInterval==0)
      {
        item.recordInterval=1
      }
      this.S7DataType = item.type
      setTimeout(function (){
        _t.RegisterForm.setFieldsValue(
            {
              NodeIDName:item.name,
              DataFromType:item.DataFromType,
              NodeIDDataType:item.type,
              NodeIDAccessLevel:item.auth,
              dataUnit:item.unit,
              DBOffset:item.DBOffset,
              StringLength:item.StringMaxLength,
              IsHaveUnsigned:item.IsHaveUnsigned?item.IsHaveUnsigned.toString():'0',
              ConversionExpression:item.conversionExpression,
              NodeIDDec:item.Description,
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
        if(_t.DataFrom==0){
          _t.RegisterForm.setFieldsValue(
              {
                DBIndex:item.DBIndex.toString()
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
              DataFromType:this.RegisterForm.getFieldValue('DataFromType'),
              auth: this.RegisterForm.getFieldValue('NodeIDAccessLevel'),
              type: this.RegisterForm.getFieldValue('NodeIDDataType'),
              unit: this.RegisterForm.getFieldValue('dataUnit'),
              DBOffset :this.RegisterForm.getFieldValue('DBOffset'),
              StringMaxLength:parseInt(this.RegisterForm.getFieldValue('StringLength')),
              IsHaveUnsigned:parseInt(this.RegisterForm.getFieldValue('IsHaveUnsigned')),
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
          if(this.DataFrom==0)
          {
            params.data.DBIndex = parseInt(this.RegisterForm.getFieldValue('DBIndex'))
          }
          let _t = this
          SimS7ModelDataEdit(params).then(function (res){
            if(res.data.code==2002)
            {
              const newData = [..._t.registerGroupDataSource];
              const target = newData.filter(item => _t.editingKey === item.uuid)[0];
              if (target) {
                _t.DataFrom = parseInt(_t.RegisterForm.getFieldValue('DataFromType'))
                target.name = _t.RegisterForm.getFieldValue('NodeIDName')
                target.NodeIDPath = _t.RegisterForm.getFieldValue('NodeIDPath')
                target.DataFromType=parseInt(_t.RegisterForm.getFieldValue('DataFromType'))
                target.DBIndex=parseInt(_t.RegisterForm.getFieldValue('DBIndex'))
                target.DBOffset=_t.RegisterForm.getFieldValue('DBOffset')
                target.StringMaxLength=parseInt(_t.RegisterForm.getFieldValue('StringLength'))
                target.IsHaveUnsigned=parseInt(_t.RegisterForm.getFieldValue('IsHaveUnsigned'))
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
      SimS7ModelDataList(params).then(function (res){
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
            modeltype:15,
            name:this.RegisterForm.getFieldValue('NodeIDName'),
            DataFromType:this.RegisterForm.getFieldValue('DataFromType'),
            DBIndex:parseInt(this.RegisterForm.getFieldValue('DBIndex')),
            DBOffset:this.RegisterForm.getFieldValue('DBOffset'),
            StringMaxLength:parseInt(this.RegisterForm.getFieldValue('StringLength')),
            IsHaveUnsigned:parseInt(this.RegisterForm.getFieldValue('IsHaveUnsigned')),
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
          SimS7ModelDataAdd(params).then(function (res){
            _t.messageShowLoad = false
            if(res.data.code==2002)
            {
              _t.NodeIdList()
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
    deleteRecord(uuid,muid) {
      let _t = this
      let params={
        uuid:uuid,
        muid:muid
      }
      SimS7ModelDataDel(params).then(function (res) {
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
      this.$router.push('/DeviceModel/SiemensS7Model')
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