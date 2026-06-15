<template>
  <a-spin style="padding: 1px;"  :spinning="messageShowLoad" tip="Loading...">
  <a-card :body-style="{padding: '24px 32px'}" :bordered="false">
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
  </a-card>
  </a-spin>
</template>

<script>
import {COMListGet, ModbusModelAdd} from "@/services/modbusModel";
import {displayModelList, getDisplayModelLayerData} from "@/services/displayModel";

export default {
  name: 'ModbusModelAdd',
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
      version:1,
      value: 1
    }
  },
  mounted(){
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
    onSubmit (e) {
      e.preventDefault()
      let _t = this
      this.form.validateFields((err) => {
        if (!err) {
          this.messageShowLoad = true
          const params = {
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
          };
          ModbusModelAdd(params).then(function (res){
            _t.messageShowLoad = false
            if (res.data.code == 2002) {
              _t.$message.success(_t.$t('dataModel.modelAddSuccess'), 3)
              _t.$router.push('/DeviceModel/CJT188Model')
            }
          }).catch(function (){
            _t.messageShowLoad = false
            _t.$message.error(_t.$t('loginPage.serverError'), 3)
          })
        }
      })
    },
    onBlackCLK(){
      this.$router.push('/DeviceModel/DLT645Model')
    },
    getCommList(){
      let _t = this
      _t.messageShowLoad = true
      this.$nextTick(function(){
        _t.form.setFieldsValue(
            {
              timeout:"1000",
              modbusConnectType:_t.modbusConnectType,
              ModbusType:"RTU",
              DataFormat:"BigEndian",
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

        _t.form.setFieldsValue(
            {
              modbusConnectType:_t.modbusConnectType,
              ModbusType:"RTU",
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
    handleSelectChange(value){
      this.modbusConnectType = value
    },
  }
}
</script>

<style scoped>
.ant-form-item {
  margin-bottom: 5px;

}
</style>