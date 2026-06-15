<template>
  <a-card>
    <a-space class="operator">
      <a-button @click="addVisible=true;isEdit=false" type="primary" icon="plus">{{$t('dataModel.newModel')}}</a-button>
      <a-button @click="refresh()"  type="default" icon="sync" :loading="refIconLoading">{{$t("dataModel.refModel")}}</a-button>
    </a-space>

    <a-spin style="padding: 1px;"  :spinning="messageShowLoad" tip="Loading...">
      <a-table rowKey="InterfaceName" :pagination="pagination" :columns="columns" :data-source="dataSource">
        <template v-for="(item, index) in columns" :slot="item.slotName">
          <span :key="index">{{ $t(item.slotName) }}</span>
        </template>
        <div slot="InterfaceType" slot-scope="text" >
          <span v-if="text==1">{{ $t('ISMDataInterface.InterfaceUrl') }}</span>
          <span v-else-if="text==2">{{ $t('ISMDataInterface.InterfaceUrlPush') }}</span>
          <span v-else-if="text==3">{{ $t('ISMDataInterface.InterfaceMqttPush') }}</span>
          <span v-else-if="text==4">{{ $t('ISMDataInterface.InterfaceIEC104') }}</span>
          <span v-else-if="text==5">{{ $t('ISMDataInterface.InterfaceModbusTcp') }}</span>
          <span v-else-if="text==6">{{ $t('ISMDataInterface.InterfaceModbusRTU') }}</span>
        </div>
        <div slot="InterfaceStatus" slot-scope="InterfaceStatus,record" >
          <a-switch :checked="InterfaceStatus==0?false:true" @change="onChangeStatus($event,record)">
            <a-icon slot="checkedChildren" type="check" />
            <a-icon slot="unCheckedChildren" type="close" />
          </a-switch>
        </div>
        <div slot="action" slot-scope="text, record">
          <a @click="GoToEdit(record)" style="color: #13C2C2"><a-icon type="edit" />{{$t('dataModel.modelDetail')}}</a> |
          <a-popconfirm :title="$t('dataModel.deleteConfirm')" @confirm="deleteRecord(record.InterfaceUuid,record.InterfaceType)">
            <a-icon slot="icon" type="question-circle-o" style="color: red" />
            <a-icon type="delete" theme="twoTone" two-tone-color="#eb2f96"/><a style="color: #eb2f96">{{$t('dataModel.delete')}}</a>
          </a-popconfirm>
        </div>
      </a-table>
    </a-spin>
    <a-drawer
        :title="isEdit?$t('ISMDataInterface.editData'):$t('ISMDataInterface.addData')"
        :width="800"
        :visible="addVisible"
        :body-style="{ paddingBottom: '80px' }"
        @close="onClose"
    >
      <a-form :form="PlanForm" layout="vertical" >
        <a-row :gutter="16">
          <a-col :span="8">
            <a-form-item
                :label="$t('ISMDataInterface.InterfaceName')"
            >
              <a-input  autocomplete="autocomplete"

                        v-decorator="['InterfaceName', {rules: [{ required: true, message: $t('ISMDataInterface.InterfaceName'), whitespace: true}]}]"
              />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item
                :label="$t('ISMDataInterface.InterfaceType')"
            >
                <a-select @change="changeInterfaceType"
                        v-decorator="[
                  'InterfaceType',
                  {
                    rules: [{ required: true, message: $t('ISMDataInterface.InterfaceType') }],
                  },
                ]"
              >
                <a-select-option  v-for="(device,index) in supportTypeList" :key="index" :value=device.type>
                  {{ $t(device.name)  }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item
                :label="$t('ISMDataInterface.InterfaceDataTemplete')"
            >
              <a-select
                  mode="multiple"
                        v-decorator="[
                  'InterfaceDataTemplete',
                  {
                    rules: [{ required: true, message: $t('ISMDataInterface.InterfaceDataTemplete') }],
                  },
                ]"
              >
                <a-select-option  v-for="(device,index) in DataTempleteList" :key="index" :value=device.TempleteUuid>
                  {{ $t(device.TempleteName)  }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
<!--推送-->
        <a-row :gutter="16" v-if="InterfaceType==1">
          <a-col :span="12">
            <a-form-item
                :label="$t('ISMDataInterface.InterfaceAPIUrl')"
            >
              <a-input  autocomplete="autocomplete"

                        v-decorator="['InterfaceAPIUrl', {rules: [{ required: true, message: $t('ISMDataInterface.InterfaceAPIUrl'), whitespace: true}]}]"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item
                :label="$t('ISMDataInterface.InterfaceDataType')"
            >
              <a-select
                        v-decorator="[
                  'InterfaceDataType',
                  {
                    rules: [{ required: true, message: $t('ISMDataInterface.InterfaceDataType') }],
                  },
                ]"
              >
                <a-select-option  v-for="(device,index) in supportDataTypeList" :key="index" :value=device.type>
                  {{ $t(device.name)  }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16" v-if="InterfaceType==2">
          <a-col :span="8">
            <a-form-item
                :label="$t('ISMDataInterface.InterfaceAPIPushUrl')"
            >
              <a-input  autocomplete="autocomplete"

                        v-decorator="['InterfaceAPIPushUrl', {rules: [{ required: true, message: $t('ISMDataInterface.InterfaceAPIUrl'), whitespace: true}]}]"
              />
            </a-form-item>
          </a-col>
          <a-col :span="5">
            <a-form-item
                :label="$t('ISMDataInterface.InterfaceAPIPushMethod')"
            >
              <a-select 
                        v-decorator="[
                  'InterfaceAPIPushMethod',
                  {
                    rules: [{ required: true, message: $t('ISMDataInterface.InterfaceAPIPushMethod') }],
                  },
                ]"
              >
                <a-select-option  v-for="(device,index) in supportMethods" :key="index" :value=device.type>
                  {{ $t(device.name)  }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="5">
            <a-form-item
                :label="$t('ISMDataInterface.InterfaceDataType')"
            >
              <a-select
                        v-decorator="[
                  'InterfaceDataType',
                  {
                    rules: [{ required: true, message: $t('ISMDataInterface.InterfaceDataType') }],
                  },
                ]"
              >
                <a-select-option  v-for="(device,index) in supportDataTypeList" :key="index" :value=device.type>
                  {{ $t(device.name)  }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="5">
            <a-form-item
                :label="$t('ISMDataInterface.InterfaceAPIPushInterval')"
            >
              <a-input
                        v-decorator="[
                  'InterfaceAPIPushInterval',
                  {
                    rules: [{ required: true, message: $t('ISMDataInterface.InterfaceAPIPushInterval') }],
                  },
                ]"
              />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16" v-if="InterfaceType==2">
          <a-card :title="$t('ISMDataInterface.InterfaceAPIPushHeader')">
            <a slot="extra"  @click="addHeader">{{$t('ISMDataInterface.InterfaceAPIPushHeaderAdd')}}</a>
            <div v-for="(device,index) in HeaderList" :key="index">
              <a-col :span="7">
                <a-form-item
                    :label="$t('ISMDataInterface.InterfaceAPIPushHeaderName')"
                >
                  <a-input  autocomplete="autocomplete" v-model="device.HeaderName"
                  />
                </a-form-item>
              </a-col>
              <a-col :span="7">
                <a-form-item
                    :label="$t('ISMDataInterface.InterfaceAPIPushHeaderValue')"
                >
                  <a-input  autocomplete="autocomplete" v-model="device.HeaderValue"
                  />
                </a-form-item>
              </a-col>
              <a-col :span="7">
                <a-form-item
                    :label="$t('ISMDataInterface.InterfaceAPIPushHeaderDes')"
                >
                  <a-input  autocomplete="autocomplete" v-model="device.HeaderDes"
                  />
                </a-form-item>
              </a-col>
              <a-col :span="3">
                <a-form-item
                    :label="$t('ISMDataInterface.InterfaceAPIPushHeaderOpt')"
                >
                  <a-button @click="DelHeader(index)">{{$t('ISMDataInterface.InterfaceAPIPushHeaderDel')}}</a-button>
                </a-form-item>

              </a-col>
            </div>
          </a-card>
        </a-row>
        <a-row :gutter="16" style="margin-top: 10px" v-if="InterfaceType==2">
          <a-card :title="$t('ISMDataInterface.InterfaceAPIPushCookie')">
            <a slot="extra"  @click="addCookie">{{$t('ISMDataInterface.InterfaceAPIPushHeaderAdd')}}</a>
            <div v-for="(device,index) in CookieList" :key="index">
              <a-col :span="7">
                <a-form-item
                    :label="$t('ISMDataInterface.InterfaceAPIPushHeaderName')"
                >
                  <a-input  autocomplete="autocomplete" v-model="device.CookieName"
                  />
                </a-form-item>
              </a-col>
              <a-col :span="7">
                <a-form-item
                    :label="$t('ISMDataInterface.InterfaceAPIPushHeaderValue')"
                >
                  <a-input  autocomplete="autocomplete" v-model="device.CookieValue"
                  />
                </a-form-item>
              </a-col>
              <a-col :span="7">
                <a-form-item
                    :label="$t('ISMDataInterface.InterfaceAPIPushHeaderDes')"
                >
                  <a-input  autocomplete="autocomplete" v-model="device.CookieDes"
                  />
                </a-form-item>
              </a-col>
              <a-col :span="3">
                <a-form-item
                    :label="$t('ISMDataInterface.InterfaceAPIPushHeaderOpt')"
                >
                  <a-button @click="DelCookie(index)">{{$t('ISMDataInterface.InterfaceAPIPushHeaderDel')}}</a-button>
                </a-form-item>

              </a-col>
            </div>
          </a-card>
        </a-row>
<!--MQTT-->
        <a-row :gutter="16" v-if="InterfaceType==3">
          <a-col :span="8">
            <a-form-item
                :label="$t('ISMDataInterface.InterfaceMqttUrl')"
            >
              <a-input  autocomplete="autocomplete"

                        v-decorator="['InterfaceMqttUrl', {rules: [{ required: true, message: $t('ISMDataInterface.InterfaceMqttUrl'), whitespace: true}]}]"
              />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item
                :label="$t('ISMDataInterface.InterfaceMqttPort')"
            >
              <a-input  autocomplete="autocomplete"

                        v-decorator="['InterfaceMqttPort', {rules: [{ required: true, message: $t('ISMDataInterface.InterfaceMqttPort'), whitespace: true}]}]"
              />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item
                :label="$t('ISMDataInterface.InterfaceMqttID')"
            >
              <a-input  autocomplete="autocomplete"

                        v-decorator="['InterfaceMqttID', {rules: [{ required: true, message: $t('ISMDataInterface.InterfaceMqttID'), whitespace: true}]}]"
              />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item
                :label="$t('ISMDataInterface.InterfaceMqttUser')"
            >
              <a-input  autocomplete="autocomplete"

                        v-decorator="['InterfaceMqttUser', {rules: [{ required: false, message: $t('ISMDataInterface.InterfaceMqttUser'), whitespace: true}]}]"
              />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item
                :label="$t('ISMDataInterface.InterfaceMqttPw')"
            >
              <a-input  autocomplete="autocomplete"

                        v-decorator="['InterfaceMqttPw', {rules: [{ required: false, message: $t('ISMDataInterface.InterfaceMqttPw'), whitespace: true}]}]"
              />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item
                :label="$t('ISMDataInterface.InterfaceAPIPushInterval')"
            >
              <a-input
                  v-decorator="[
                  'InterfaceAPIPushInterval',
                  {
                    rules: [{ required: true, message: $t('ISMDataInterface.InterfaceAPIPushInterval') }],
                  },
                ]"
              />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16" v-if="InterfaceType==3">
          <a-col :span="24">
            <a-form-item>
              <span slot="label">
                {{$t('ISMDataInterface.InterfaceMqttSub')}}&nbsp;
                <a-tooltip :title="$t('ISMDataInterface.InterfaceMqttSubTips')">
                  <a-icon type="question-circle-o" />
                </a-tooltip>
              </span>
              <a-input  autocomplete="autocomplete"

                        v-decorator="['InterfaceMqttSub', {rules: [{ required: true, message: $t('ISMDataInterface.InterfaceMqttSub'), whitespace: true}]}]"
              />
            </a-form-item>
          </a-col>
        </a-row>
<!--IEC104-->
        <a-row :gutter="16" v-if="InterfaceType==4">
          <a-col :span="8">
            <a-form-item
                :label="$t('ISMDataInterface.InterfaceIEC104Port')"
            >
              <a-input  autocomplete="autocomplete"

                        v-decorator="['InterfaceIEC104Port', {rules: [{ required: true, message: $t('ISMDataInterface.InterfaceIEC104Port'), whitespace: true}]}]"
              />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item
                :label="$t('ISMDataInterface.InterfaceIEC104Addr')"
            >
              <a-input  autocomplete="autocomplete"

                        v-decorator="['InterfaceIEC104Addr', {rules: [{ required: true, message: $t('ISMDataInterface.InterfaceIEC104Addr'), whitespace: true}]}]"
              />
            </a-form-item>
          </a-col>
        </a-row>
<!--Modbus Tcp-->
        <a-row :gutter="16" v-if="InterfaceType==5">
            <a-col :span="8">
              <a-form-item
                  :label="$t('ISMDataInterface.InterfaceModbusTcpPort')"
              >
                <a-input  autocomplete="autocomplete"

                          v-decorator="['InterfaceModbusTcpPort', {rules: [{ required: true, message: $t('ISMDataInterface.InterfaceModbusTcpPort'), whitespace: true}]}]"
                />
              </a-form-item>
            </a-col>
            <a-col :span="8">
              <a-form-item
                  :label="$t('ISMDataInterface.InterfaceModbusTcpAddr')"
              >
                <a-input  autocomplete="autocomplete"

                          v-decorator="['InterfaceModbusTcpAddr', {rules: [{ required: true, message: $t('ISMDataInterface.InterfaceModbusTcpAddr'), whitespace: true}]}]"
                />
              </a-form-item>
            </a-col>
        </a-row>
<!--Modbus RTU-->
        <a-row :gutter="24" v-if="InterfaceType==6">
          <a-col :span="8">
            <a-form-item
                :label="$t('ISMDataInterface.InterfaceModbusTcpAddr')"
            >
              <a-input  autocomplete="autocomplete"

                        v-decorator="['InterfaceModbusTcpAddr', {rules: [{ required: true, message: $t('ISMDataInterface.InterfaceModbusTcpAddr'), whitespace: true}]}]"
              />
            </a-form-item>
          </a-col>

          <a-col :span="8">
            <a-form-item
                :label="$t('dataModel.modbusModel.SerialConnection')"

            >

              <a-select
                  v-decorator="['SerialPort', {rules: [{ required: true, message: $t('dataModel.modbusModel.SerialConnection'), whitespace: true}]}]"
              >
                <a-select-option v-for="options in COMList" :key="options" :value="options">
                  {{ options }}
                </a-select-option>

              </a-select>
            </a-form-item>
          </a-col>

          <a-col :span="8">
            <a-form-item
              :label="$t('dataModel.modbusModel.SerialBaud')"

          >

            <a-select
                v-decorator="['SerialPortBaud', {rules: [{ required: true, message: $t('dataModel.modbusModel.SerialBaud'), whitespace: true}]}]"
            >
              <a-select-option v-for="options in BaudList" :key="options" :value="options">
                {{ options }}
              </a-select-option>

            </a-select>
          </a-form-item>
          </a-col>

          <a-col :span="8">
            <a-form-item
              :label="$t('dataModel.modbusModel.SerialPortDataBit')"

          >

            <a-select
                v-decorator="['SerialPortDataBit', {rules: [{ required: true, message: $t('dataModel.modbusModel.SerialPortDataBit'), whitespace: true}]}]"
            >
              <a-select-option v-for="options in SerialDataBitList" :key="options" :value="options">
                {{ options }}
              </a-select-option>

            </a-select>
          </a-form-item>
          </a-col>

          <a-col :span="8">
              <a-form-item
            :label="$t('dataModel.modbusModel.SerialPortVerifyBit')"

        >

          <a-select
              v-decorator="['SerialPortVerifyBit', {rules: [{ required: true, message: $t('dataModel.modbusModel.SerialPortVerifyBit'), whitespace: true}]}]"
          >
            <a-select-option v-for="options in SerialVerifyList" :key="options" :value="options">
              {{ options }}
            </a-select-option>

          </a-select>
        </a-form-item>
          </a-col>

          <a-col :span="8">
            <a-form-item
        :label="$t('dataModel.modbusModel.SerialPortStopBit')"

    >

      <a-select
          v-decorator="['SerialPortStopBit', {rules: [{ required: true, message: $t('dataModel.modbusModel.SerialPortStopBit'), whitespace: true}]}]"
      >
        <a-select-option v-for="options in SerialStopBitList" :key="options" :value="options">
          {{ options }}
        </a-select-option>

      </a-select>
    </a-form-item>
          </a-col>

          <a-col :span="8">
            <a-form-item :label="$t('dataModel.modbusModel.SerialPortFlow')">

              <a-select
                  v-decorator="['SerialPortFlow', {rules: [{ required: true, message: $t('dataModel.modbusModel.SerialPortFlow'), whitespace: true}]}]"
              >
                <a-select-option v-for="options in SerialFlowList" :key="options" :value="options">
                  {{ options }}
                </a-select-option>

              </a-select>
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
import {COMListGet} from "@/services/modbusModel";
import moment from 'moment'
import {formatDate} from "@/utils/common";
import {
  AddInterfaceData,
  DelInterfaceData,
  EditInterfaceData,
  GetInterfaceDataList,
  EditInterfaceStatus
} from "../../services/datainterface";
import {GetTempleteDataList} from "../../services/datatemplete";
export default {
  name: 'DataInterface',
  i18n: require('@/i18n/language'),
  data() {
    return {
      COMList:[],
      InterfaceType:1,
      InterfaceStatus:1,
      DataTempleteList:[],
      BaudList:["1200","2400","4800","9600","19200","38400","115200"],
      SerialDataBitList:["5","6","7","8"],
      SerialStopBitList:["1","2","1.5"],
      SerialFlowList:["None","XOn/XOff","RTS/CTS"],
      SerialVerifyList:["None","Even","Odd"],
      InterfaceApiServer:{
        url:"",
        type:""
      },
      supportTypeList:[
        {
          type:1,
          name:"ISMDataInterface.InterfaceUrl"
        },
        {
          type:2,
          name:"ISMDataInterface.InterfaceUrlPush"
        },
        {
          type:3,
          name:"ISMDataInterface.InterfaceMqttPush"
        },
        {
          type:4,
          name:"ISMDataInterface.InterfaceIEC104"
        },
        {
          type:5,
          name:"ISMDataInterface.InterfaceModbusTcp"
        },
        {
          type:6,
          name:"ISMDataInterface.InterfaceModbusRTU"
        },
      ],
      supportDataTypeList:[
        {
          type:"application/json",
          name:"ISMDataInterface.InterfaceDataTypeJSON"
        },
        {
          type:"application/xml",
          name:"ISMDataInterface.InterfaceDataTypeXML"
        },
        {
          type:"text/plain",
          name:"ISMDataInterface.InterfaceDataTypeText"
        },
        {
          type:"text/html",
          name:"ISMDataInterface.InterfaceDataTypeHtml"
        },
      ],
      supportMethods:[
        {
          type:1,
          name:"GET"
        },
        {
          type:2,
          name:"POST"
        },
      ],
      pagination: {
        pageSize: 15,
        showSizeChanger: true
      },
      HeaderList:[
        {
          HeaderName:"",
          HeaderValue:"",
          HeaderDes:""
        }
      ],
      CookieList:[
        {
          CookieName:"",
          CookieValue:"",
          CookieDes:""
        }
      ],
      ScriptType: 0,
      isCharge: true,
      CodeContent: "",
      isEdit: false,
      messageShowLoad: false,
      advanced: true,
      refIconLoading: false,
      columns: [
        {
          width: '20%',
          slotName: 'ISMDataInterface.InterfaceName',
          scopedSlots: {customRender: 'InterfaceName', title: 'ISMDataInterface.InterfaceName'},
          dataIndex: 'InterfaceName'
        },
        {
          slotName: 'ISMDataInterface.InterfaceType',
          width: '30%',
          scopedSlots: {customRender: 'InterfaceType', title: 'ISMDataInterface.InterfaceType'},
          dataIndex: 'InterfaceType',
        },
        {
          width: '10%',
          slotName: 'ISMDataInterface.InterfaceStatus',
          scopedSlots: {customRender: 'InterfaceStatus', title: 'ISMDataInterface.InterfaceStatus'},
          dataIndex: 'InterfaceStatus'
        },
        {
          width: '10%',
          slotName: 'dataModel.modelTableOpt',
          scopedSlots: {customRender: 'action', title: 'dataModel.modelTableOpt'}
        }
      ],
      dataSource: [],
      addVisible: false,
      error: '',
      editUuid: "",
      editVisible: false,
      PlanForm: this.$form.createForm(this),
      editForm: this.$form.createForm(this),
      textAreValue: "",
      that: this,
      value: 1
    }
  },
  components: {

  },
  authorize: {
    // deleteRecord: 'delete'
  },
  filters: {
    formatDate(time) {
      let date = new Date(time)
      return formatDate(date, 'yyyy-MM-dd hh:mm:ss')
    }
  },
  mounted() {
  },
  activated() {
    this.GetTempleteList()
    this.GetInterfaceList()
  },
  created() {

  },
  methods: {
    getCommList(){
      let _t = this
      _t.messageShowLoad = true
      this.$nextTick(function(){
        _t.PlanForm.setFieldsValue(
            {
              SerialPort:_t.COMList[0],
              SerialPortBaud:"9600",
              SerialPortDataBit:"8",
              SerialPortVerifyBit:"None",
              SerialPortStopBit:"1",
              SerialPortFlow:"None",
            })
      });

      COMListGet().then(function (res){
        _t.COMList = res.data

        _t.PlanForm.setFieldsValue(
            {
              SerialPort:_t.COMList[0],
              SerialPortBaud:"9600",
              SerialPortDataBit:"8",
              SerialPortVerifyBit:"None",
              SerialPortStopBit:"1",
              SerialPortFlow:"None",
            })
        _t.messageShowLoad = false
      })
    },
    onChangeStatus(value,item){
      let _t = this
      item.InterfaceStatus = value?1:0
      const params = {
        uuid:item.InterfaceUuid,
        data: item
      }
      EditInterfaceData(params).then(function (res){
        if (res.data.code == 200) {
          _t.GetInterfaceList()
          _t.$message.success(_t.$t('dataModel.static.EditSuccess'), 3)
        }
        else {
          _t.$message.error(_t.$t('dataModel.static.EditFailed'), 3)
        }
      })
    },
    addHeader(){
      let header = {
        HeaderName:"",
        HeaderValue:"",
        HeaderDes:""
      }
      this.HeaderList.push(header)
    },
    DelHeader(index){
      this.HeaderList.splice(index,1)
    },
    addCookie(){
      let header = {
        CookieName:"",
        CookieValue:"",
        CookieDes:""
      }
      this.CookieList.push(header)
    },
    DelCookie(index){
      this.CookieList.splice(index,1)
    },
    GetTempleteList() {
      let _t = this
      this.DataTempleteList = []
      GetTempleteDataList().then(function (res) {
        if (res.data.code == 200) {
          _t.DataTempleteList = res.data.list
        } else if (res.data.code == 2001) {
          _t.$message.error(_t.$t('displayModel.ModelExist'), 3)
        } else if (res.data.code == 2003) {
          _t.$message.error(_t.$t('displayModel.AddModelFailed'), 3)
        }
      })
    },
    changeInterfaceType(value){
      this.InterfaceType = value
      if(this.InterfaceType==6)
      {
        this.getCommList()
      }
    },
    onClose() {
      this.addVisible = false;
    },
    GetInterfaceList() {
      let _t = this
      this.dataSource = []
      GetInterfaceDataList().then(function (res) {
        _t.refIconLoading = false
        if (res.data.code == 200) {
          _t.dataSource = res.data.list
          _t.addVisible = false;
        } else if (res.data.code == 2001) {
          _t.$message.error(_t.$t('displayModel.ModelExist'), 3)
        } else if (res.data.code == 2003) {
          _t.$message.error(_t.$t('displayModel.AddModelFailed'), 3)
        }
      })
    },
    AddScript() {
      let _t = this
      this.PlanForm.validateFields((err) => {
        if (!err) {
          let InterfaceDataUuidArray = _t.PlanForm.getFieldValue('InterfaceDataTemplete')
          let params = {
            InterfaceName: _t.PlanForm.getFieldValue('InterfaceName'),
            InterfaceType: _t.InterfaceType,
            InterfaceStatus: 1,
            InterfaceDataUuid:InterfaceDataUuidArray.join(",")
          }
          if(_t.InterfaceType==1)
          {
            params.InterfaceContent = JSON.stringify(
              {
                url:_t.PlanForm.getFieldValue('InterfaceAPIUrl'),
                type:_t.PlanForm.getFieldValue('InterfaceDataType')
              }
            )
          }else if(_t.InterfaceType==2)
          {
            params.InterfaceContent = JSON.stringify(
                {
                  url:_t.PlanForm.getFieldValue('InterfaceAPIPushUrl'),
                  type:_t.PlanForm.getFieldValue('InterfaceDataType'),
                  method:parseInt(_t.PlanForm.getFieldValue('InterfaceAPIPushMethod')),
                  interval:parseInt(_t.PlanForm.getFieldValue('InterfaceAPIPushInterval')),
                  header:_t.HeaderList,
                  cookie:_t.CookieList,
                }
            )
          }else if(_t.InterfaceType==3)
          {
            params.InterfaceContent = JSON.stringify(
                {
                  host:_t.PlanForm.getFieldValue('InterfaceMqttUrl'),
                  port:_t.PlanForm.getFieldValue('InterfaceMqttPort'),
                  id:_t.PlanForm.getFieldValue('InterfaceMqttID'),
                  user:_t.PlanForm.getFieldValue('InterfaceMqttUser'),
                  password:_t.PlanForm.getFieldValue('InterfaceMqttPw'),
                  interval:parseInt(_t.PlanForm.getFieldValue('InterfaceAPIPushInterval')),
                  subject:_t.PlanForm.getFieldValue('InterfaceMqttSub'),
                }
            )
          }else if(_t.InterfaceType==4)
          {
            params.InterfaceContent = JSON.stringify(
                {
                  addr:parseInt(_t.PlanForm.getFieldValue('InterfaceIEC104Addr')),
                  port:parseInt(_t.PlanForm.getFieldValue('InterfaceIEC104Port')),
                }
            )
          }else if(_t.InterfaceType==5)
          {
            params.InterfaceContent = JSON.stringify(
                {
                  addr:parseInt(_t.PlanForm.getFieldValue('InterfaceModbusTcpAddr')),
                  port:parseInt(_t.PlanForm.getFieldValue('InterfaceModbusTcpPort')),
                }
            )
          }
          else if(_t.InterfaceType==6)
          {
            params.InterfaceContent = JSON.stringify(
                {
                  addr:parseInt(_t.PlanForm.getFieldValue('InterfaceModbusTcpAddr')),
                  SerialPort:_t.PlanForm.getFieldValue('SerialPort'),
                  SerialPortBaud:parseInt(_t.PlanForm.getFieldValue('SerialPortBaud')),
                  SerialPortDataBit:parseInt(_t.PlanForm.getFieldValue('SerialPortDataBit')),
                  SerialPortVerifyBit:_t.PlanForm.getFieldValue('SerialPortVerifyBit'),
                  SerialPortStopBit:_t.PlanForm.getFieldValue('SerialPortStopBit'),
                  SerialPortFlow:_t.PlanForm.getFieldValue('SerialPortFlow'),
                }
            )
          }else{
            return
          }
          AddInterfaceData(params).then(function (res) {
            if (res.data.code == 2002) {
              _t.GetInterfaceList()
              _t.addVisible = false;
              _t.$message.success(_t.$t('ISMDataTemplete.AddSuccess'), 3)
            } else {
              _t.$message.error(_t.$t('ISMDataTemplete.AddFailed'), 3)
            }
          })
        }
      })
    },
    EditScript() {
      let _t = this
      this.PlanForm.validateFields((err) => {
        if (!err) {
          let InterfaceDataUuidArray = _t.PlanForm.getFieldValue('InterfaceDataTemplete')
          let params = {
            uuid: _t.EditUUid,
            data: {
              InterfaceName: _t.PlanForm.getFieldValue('InterfaceName'),
              InterfaceType: _t.InterfaceType,
              InterfaceStatus: _t.InterfaceStatus,
              InterfaceDataUuid: InterfaceDataUuidArray.join(",")
            }
          }
          if(_t.InterfaceType==1)
          {
            params.data.InterfaceContent = JSON.stringify(
                {
                  url:_t.PlanForm.getFieldValue('InterfaceAPIUrl'),
                  type:_t.PlanForm.getFieldValue('InterfaceDataType'),
                }
            )
          }else if(_t.InterfaceType==2)
          {
            params.data.InterfaceContent = JSON.stringify(
                {
                  url:_t.PlanForm.getFieldValue('InterfaceAPIPushUrl'),
                  type:_t.PlanForm.getFieldValue('InterfaceDataType'),
                  method:parseInt(_t.PlanForm.getFieldValue('InterfaceAPIPushMethod')),
                  interval:parseInt(_t.PlanForm.getFieldValue('InterfaceAPIPushInterval')),
                  header:_t.HeaderList,
                  cookie:_t.CookieList,
                }
            )
          }
          else if(_t.InterfaceType==3)
          {
            params.data.InterfaceContent = JSON.stringify(
                {
                  host:_t.PlanForm.getFieldValue('InterfaceMqttUrl'),
                  port:_t.PlanForm.getFieldValue('InterfaceMqttPort'),
                  id:_t.PlanForm.getFieldValue('InterfaceMqttID'),
                  user:_t.PlanForm.getFieldValue('InterfaceMqttUser'),
                  password:_t.PlanForm.getFieldValue('InterfaceMqttPw'),
                  interval:parseInt(_t.PlanForm.getFieldValue('InterfaceAPIPushInterval')),
                  subject:_t.PlanForm.getFieldValue('InterfaceMqttSub'),
                }
            )
          }else if(_t.InterfaceType==4)
          {
            params.data.InterfaceContent = JSON.stringify(
                {
                  addr:parseInt(_t.PlanForm.getFieldValue('InterfaceIEC104Addr')),
                  port:parseInt(_t.PlanForm.getFieldValue('InterfaceIEC104Port')),
                }
            )
          }else if(_t.InterfaceType==5)
          {
            params.data.InterfaceContent = JSON.stringify(
                {
                  addr:parseInt(_t.PlanForm.getFieldValue('InterfaceModbusTcpAddr')),
                  port:parseInt(_t.PlanForm.getFieldValue('InterfaceModbusTcpPort')),
                }
            )
          }
          else if(_t.InterfaceType==6)
          {
            _t.getCommList()
            params.data.InterfaceContent = JSON.stringify(
                {
                  addr:parseInt(_t.PlanForm.getFieldValue('InterfaceModbusTcpAddr')),
                  SerialPort:_t.PlanForm.getFieldValue('SerialPort'),
                  SerialPortBaud:parseInt(_t.PlanForm.getFieldValue('SerialPortBaud')),
                  SerialPortDataBit:parseInt(_t.PlanForm.getFieldValue('SerialPortDataBit')),
                  SerialPortVerifyBit:_t.PlanForm.getFieldValue('SerialPortVerifyBit'),
                  SerialPortStopBit:_t.PlanForm.getFieldValue('SerialPortStopBit'),
                  SerialPortFlow:_t.PlanForm.getFieldValue('SerialPortFlow'),
                }
            )
          }
          EditInterfaceData(params).then(function (res) {
            if (res.data.code == 200) {
              _t.GetInterfaceList()
              _t.addVisible = false;
              _t.$message.success(_t.$t('ISMDataTemplete.EditSuccess'), 3)
            } else {
              _t.$message.error(_t.$t('ISMDataTemplete.EditFailed'), 3)
            }
          })
        }
      })
    },
    GoToEdit(item) {
      let _t = this
      _t.isCharge = false
      this.isEdit = true
      this.addVisible = true
      _t.InterfaceType = item.InterfaceType
      _t.EditUUid = item.InterfaceUuid
      _t.InterfaceStatus = item.InterfaceStatus
      setTimeout(function () {
        _t.isCharge = true
        let InterfaceDataUuidArray=item.InterfaceDataUuid.split(",");
        _t.PlanForm.setFieldsValue(
            {
              InterfaceName: item.InterfaceName,
              InterfaceType: _t.InterfaceType,
              InterfaceDataTemplete: InterfaceDataUuidArray,
            })
        if(_t.InterfaceType==1){
          try {
            let APIUrl = JSON.parse(item.InterfaceContent)
            _t.PlanForm.setFieldsValue(
                {
                  InterfaceAPIUrl: APIUrl.url,
                  InterfaceDataType:APIUrl.type,
                })
          }catch (e){
            _t.PlanForm.setFieldsValue(
                {
                  InterfaceAPIUrl: "",
                  InterfaceDataType:"",
                })
          }
        }else if(_t.InterfaceType==2){
          try {
            let APIUrl = JSON.parse(item.InterfaceContent)
            _t.PlanForm.setFieldsValue(
                {
                  InterfaceAPIPushUrl: APIUrl.url,
                  InterfaceDataType:APIUrl.type,
                  InterfaceAPIPushMethod:APIUrl.method,
                  InterfaceAPIPushInterval:APIUrl.interval
                })
            _t.HeaderList = APIUrl.header
            _t.CookieList=APIUrl.cookie
          }catch (e){
            _t.HeaderList=[]
            _t.CookieList=[]
            _t.PlanForm.setFieldsValue(
                {
                  InterfaceAPIUrl: "",
                  InterfaceDataType:"",
                  InterfaceAPIPushMethod:"",
                  InterfaceAPIPushInterval:0
                })
          }
        }else if(_t.InterfaceType==3){
          try {
            let APIUrl = JSON.parse(item.InterfaceContent)
            _t.PlanForm.setFieldsValue(
                {
                  InterfaceMqttUrl: APIUrl.host,
                  InterfaceMqttPort:APIUrl.port,
                  InterfaceMqttID:APIUrl.id,
                  InterfaceMqttUser:APIUrl.user,
                  InterfaceMqttPw:APIUrl.password,
                  InterfaceAPIPushInterval:APIUrl.interval,
                  InterfaceMqttSub:APIUrl.subject,
                })
          }catch (e){
            _t.HeaderList=[]
            _t.CookieList=[]
            _t.PlanForm.setFieldsValue(
                {
                  InterfaceMqttUrl: "",
                  InterfaceMqttPort:"",
                  InterfaceMqttID:"",
                  InterfaceMqttUser:"",
                  InterfaceMqttPw:"",
                  InterfaceAPIPushInterval:1000,
                  InterfaceMqttSub:"",
                })
          }
        }else if(_t.InterfaceType==4){
          try {
            let APIUrl = JSON.parse(item.InterfaceContent)
            _t.PlanForm.setFieldsValue(
                {
                  InterfaceIEC104Addr: APIUrl.addr.toString(),
                  InterfaceIEC104Port:APIUrl.port.toString()
                })
          }catch (e){

          }
        }else if(_t.InterfaceType==5){
          try {
            let APIUrl = JSON.parse(item.InterfaceContent)
            _t.PlanForm.setFieldsValue(
                {
                  InterfaceModbusTcpAddr: APIUrl.addr.toString(),
                  InterfaceModbusTcpPort:APIUrl.port.toString()
                })
          }catch (e){

          }
        }else if(_t.InterfaceType==6){
          try {
            let APIUrl = JSON.parse(item.InterfaceContent)
            _t.PlanForm.setFieldsValue(
                {
                  InterfaceModbusTcpAddr: APIUrl.addr.toString(),
                  SerialPort:APIUrl.SerialPort,
                  SerialPortBaud:APIUrl.SerialPortBaud.toString(),
                  SerialPortDataBit:APIUrl.SerialPortDataBit.toString(),
                  SerialPortVerifyBit:APIUrl.SerialPortVerifyBit,
                  SerialPortStopBit:APIUrl.SerialPortStopBit,
                  SerialPortFlow:APIUrl.SerialPortFlow,
                })
          }catch (e){

          }
        }
      }, 200)
    },
    refresh() {
      this.refIconLoading = true
      this.GetInterfaceList()
    },
    deleteRecord(uuid,type) {
      let _t = this
      const params = {
        InterfaceUuid: uuid,
        InterfaceType:type
      }
      DelInterfaceData(params).then(function (res) {
        if (res.data.code == 200) {
          _t.GetInterfaceList()
          _t.addVisible = false;
          _t.$message.success(_t.$t('ISMDataTemplete.DelSuccess'), 3)
        } else {
          _t.$message.error(_t.$t('ISMDataTemplete.DelFailed'), 3)
        }
      })
    },
  }
}
</script>

<style lang="less" scoped>
::v-deep .search {
  margin-bottom: 54px;
}

::v-deep .ant-form-item {
  margin-bottom: 1px;
}

::v-deep .ant-row .ant-form-item {
  margin-bottom: 1px;
}

::v-deep .fold {
  width: calc(100% - 216px);
  display: inline-block
}

::v-deep .operator {
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

::v-deep .ant-table-thead > tr > th {
  padding: 10px 10px;
  overflow-wrap: break-word;
}

::v-deep .ant-table-tbody > tr > td {
  padding: 7px 5px;
  overflow-wrap: break-word;
}

::v-deep .ant-table-thead > tr > th {
  color: #909399;
  font-weight: 500;
  text-align: left;
  //background: #f8f8f8;
  //border-bottom: 1px solid #e8e8e8;
  transition: background .3s ease;
}
</style>
