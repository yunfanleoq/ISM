<template>
  <a-card  >
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
            :label="$t('dataModel.modbusModel.connection')"
            :labelCol="{span: 7}"
            :wrapperCol="{span: 10}"
        >
          <a-select @change="handleSelectChange"   autocomplete="autocomplete"
                    v-decorator="['modbusConnectType', {rules: [{ required: true, message: $t('dataModel.modbusModel.connection'), whitespace: true}]}]"
          >
            <a-select-option value="Serial">{{$t('dataModel.modbusModel.SerialConnection')}}</a-select-option>
            <a-select-option value="TCPClient">{{$t('dataModel.modbusModel.TCPClientConnection')}}</a-select-option>
            <a-select-option value="TCPServer">{{$t('dataModel.modbusModel.TCPServerConnection')}}</a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item
            :label="$t('dataModel.TimeOut')"
            :labelCol="{span: 7}"
            :wrapperCol="{span: 10}">
          <a-input  v-decorator="['timeout', {rules: [{ required: true, message: $t('dataModel.TimeOut'), whitespace: true,initialValue:'1000'}]}]"/>
        </a-form-item>
        <a-form-item
            :labelCol="{span: 7}"
            :wrapperCol="{span: 10}">
           <span slot="label">
            {{$t('dataModel.DataFormat')}}&nbsp;
            <a-tooltip :title="$t('dataModel.DataFormatTips')">
              <a-icon type="question-circle-o" />
            </a-tooltip>
          </span>
          <a-select  autocomplete="autocomplete"
                     v-decorator="['DataFormat', {rules: [{ required: true, message: $t('dataModel.DataFormat'), whitespace: true}]}]"
          >
            <a-select-option value="BigEndian">{{$t('dataModel.DataFormatBigEndian')}}</a-select-option>
            <a-select-option value="LittleEndian">{{$t('dataModel.DataFormatLittleEndian')}}</a-select-option>
          </a-select>
        </a-form-item>
        <div v-if="modbusConnectType=='Serial'">
          <a-form-item
              :label="$t('dataModel.modbusModel.SerialConnection')"
              :labelCol="{span: 7}"
              :wrapperCol="{span: 10}"
          >

            <a-select
                v-decorator="['SerialPort', {rules: [{ required: true, message: $t('dataModel.modbusModel.SerialConnection'), whitespace: true}]}]"
            >
              <a-select-option v-for="options in COMList" :key="options" :value="options">
                {{ options }}
              </a-select-option>

            </a-select>
          </a-form-item>

          <a-form-item
              :label="$t('dataModel.modbusModel.SerialBaud')"
              :labelCol="{span: 7}"
              :wrapperCol="{span: 10}"
          >

            <a-select
                v-decorator="['SerialPortBaud', {rules: [{ required: true, message: $t('dataModel.modbusModel.SerialBaud'), whitespace: true}]}]"
            >
              <a-select-option v-for="options in BaudList" :key="options" :value="options">
                {{ options }}
              </a-select-option>

            </a-select>
          </a-form-item>

          <a-form-item
              :label="$t('dataModel.modbusModel.SerialPortDataBit')"
              :labelCol="{span: 7}"
              :wrapperCol="{span: 10}"
          >

            <a-select
                v-decorator="['SerialPortDataBit', {rules: [{ required: true, message: $t('dataModel.modbusModel.SerialPortDataBit'), whitespace: true}]}]"
            >
              <a-select-option v-for="options in SerialDataBitList" :key="options" :value="options">
                {{ options }}
              </a-select-option>

            </a-select>
          </a-form-item>

          <a-form-item
              :label="$t('dataModel.modbusModel.SerialPortVerifyBit')"
              :labelCol="{span: 7}"
              :wrapperCol="{span: 10}"
          >

            <a-select
                v-decorator="['SerialPortVerifyBit', {rules: [{ required: true, message: $t('dataModel.modbusModel.SerialPortVerifyBit'), whitespace: true}]}]"
            >
              <a-select-option v-for="options in SerialVerifyList" :key="options" :value="options">
                {{ options }}
              </a-select-option>

            </a-select>
          </a-form-item>

          <a-form-item
              :label="$t('dataModel.modbusModel.SerialPortStopBit')"
              :labelCol="{span: 7}"
              :wrapperCol="{span: 10}"
          >

            <a-select
                v-decorator="['SerialPortStopBit', {rules: [{ required: true, message: $t('dataModel.modbusModel.SerialPortStopBit'), whitespace: true}]}]"
            >
              <a-select-option v-for="options in SerialStopBitList" :key="options" :value="options">
                {{ options }}
              </a-select-option>

            </a-select>
          </a-form-item>

          <a-form-item
              :label="$t('dataModel.modbusModel.SerialPortFlow')"
              :labelCol="{span: 7}"
              :wrapperCol="{span: 10}"
          >

            <a-select
                v-decorator="['SerialPortFlow', {rules: [{ required: true, message: $t('dataModel.modbusModel.SerialPortFlow'), whitespace: true}]}]"
            >
              <a-select-option v-for="options in SerialFlowList" :key="options" :value="options">
                {{ options }}
              </a-select-option>

            </a-select>
          </a-form-item>
        </div>

        <a-form-item style="margin-top: 24px" :wrapperCol="{span: 10, offset: 7}">
          <a-button type="primary" htmlType="submit">{{$t('dataModel.add')}}</a-button>
          <a-button style="margin-left: 8px" @click="onBlackCLK()">{{$t('dataModel.back')}}</a-button>
        </a-form-item>
      </a-form>
    </a-spin>
  </a-card>
</template>

<script>
import {
  getSnmpModelDetail,
} from "@/services/snmpmodel";

import {
  modbusModelEdit,
} from "@/services/modbusModel";

import { uuid } from 'vue-uuid';
import {COMListGet} from "@/services/modbusModel";
import {displayModelList, getDisplayModelLayerData} from "@/services/displayModel";

const dataSource= []
export default {
  name: 'ModbusModelImport',
  i18n: require('../../../i18n/language'),
  data () {
    return {
      error: '',
      COMList:[],
      configurationModel:[],
      displayPageList:[],
      BaudList:["1200","2400","4800","9600","19200","38400","115200"],
      SerialDataBitList:["5","6","7","8"],
      SerialStopBitList:["1","2","1.5"],
      SerialFlowList:["None","XOn/XOff","RTS/CTS"],
      ModbusSerialMode:["RTU","ASCII"],
      ModbusTCPMode:["RTU","ASCII","TCP/IP"],
      SerialVerifyList:["None","Even","Odd"],
      modbusConnectType:"Serial",
      messageShowLoad:false,
      form: this.$form.createForm(this),
    }
  },
  activated() {
      this.getCommList()
  },
  mounted() {
    this.getConfigurationModel()
    this.getCommList()

  },
  computed: {
    desc() {
      return this.$t('pageDesc')
    }
  },
  methods: {
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
    getCommList(){
      let _t = this
      COMListGet().then(function (res){
        _t.COMList = res.data
        _t.getSingleModelDetail()
      })
    },
    getSingleModelDetail(){
      let _t = this
      const params={
        uuid:this.$route.params.uid
      }
      this.messageShowLoad = true
      getSnmpModelDetail(params).then(function (res){
        _t.GetDisplayPage(res.data.data.configUid)
        _t.messageShowLoad = false
        _t.modbusConnectType = res.data.data.CJT188ConnectType
        setTimeout(function (){
          _t.form.setFieldsValue(
              {
                name:res.data.data.name,
                dec:res.data.data.dec,
                timeout:res.data.data.CJT188Timeout.toString(),
                DataFormat:res.data.data.CJT188DataFormat,
                configurationModel:res.data.data.configUid,
                configurationPageUUID:res.data.data.PageUUID,
                modbusConnectType:res.data.data.CJT188ConnectType,
              })

          if(_t.modbusConnectType=="Serial")
          {
            _t.form.setFieldsValue(
                {
                  ModbusType:res.data.data.CJT188ConnectMode,
                  SerialPort:res.data.data.CJT188ConnectCOMName,
                  SerialPortBaud:res.data.data.CJT188SerialBaud.toString(),
                  SerialPortDataBit:res.data.data.CJT188SerialBits.toString(),
                  SerialPortVerifyBit:res.data.data.CJT188SerialParity,
                  SerialPortStopBit:res.data.data.CJT188SerialStopBits,
                  SerialPortFlow:res.data.data.CJT188SerialFlow,
                })
          }
          else
          {
            _t.form.setFieldsValue(
                {
                  IpAddress:res.data.data.CJT188TCPClientIpaddress,
                  Port:res.data.data.port.toString(),
                  ModbusType:res.data.data.CJT188ConnectMode,
                }
            )
          }
        },300)
      })
    },
    onSubmit (e) {
      e.preventDefault()
      let _t = this
      this.form.validateFields((err) => {
        if (!err) {
          this.messageShowLoad = true
          const params = {
            uuid:this.$route.params.uid,
            data: {
              name:this.form.getFieldValue('name'),
              dec:this.form.getFieldValue('dec'),
              type:490,
              gatherNumber:0,
              port:parseInt(this.form.getFieldValue('Port')),
              CJT188Timeout:parseInt(this.form.getFieldValue('timeout')),
              CJT188DataFormat:this.form.getFieldValue('DataFormat'),
              configUid:this.form.getFieldValue('configurationModel'),
              PageUUID:this.form.getFieldValue('configurationPageUUID'),
              CJT188TCPClientIpaddress:this.form.getFieldValue('IpAddress'),
              CJT188ConnectType:this.form.getFieldValue('modbusConnectType'),
              CJT188ConnectMode:this.form.getFieldValue('ModbusType'),
              CJT188ConnectCOMName:this.form.getFieldValue('SerialPort'),
              CJT188SerialBaud:parseInt(this.form.getFieldValue('SerialPortBaud')),
              CJT188SerialBits:parseInt(this.form.getFieldValue('SerialPortDataBit')),
              CJT188SerialParity:this.form.getFieldValue('SerialPortVerifyBit'),
              CJT188SerialStopBits:this.form.getFieldValue('SerialPortStopBit'),
              CJT188SerialFlow:this.form.getFieldValue('SerialPortFlow'),
            }
          };
          modbusModelEdit(params).then(function (res){
            _t.messageShowLoad = false
            if (res.data.code == 200) {
              _t.$message.success(_t.$t('dataModel.editSuccess'), 3)
            }
            else
            {
              _t.$message.error(_t.$t('dataModel.editFailed'), 3)
            }
          }).catch(function (){
            _t.messageShowLoad = false
            _t.$message.error(_t.$t('loginPage.serverError'), 3)
          })
        }
      })
    },
    onBlackCLK(){
      this.$router.push('/DeviceModel/CJT188Model')
    },
    handleSelectChange(value) {
      this.modbusConnectType = value
    },
  }
}
</script>


<style lang="less" scoped>
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

</style>