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
    this.getSingleModelDetail()
  },
  mounted() {
    this.getConfigurationModel()
    this.getSingleModelDetail()

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
    getSingleModelDetail(){
      let _t = this
      const params={
        uuid:this.$route.params.uid
      }
      this.messageShowLoad = true
      getSnmpModelDetail(params).then(function (res){
        _t.GetDisplayPage(res.data.data.configUid)
        _t.messageShowLoad = false
        setTimeout(function (){
          _t.form.setFieldsValue(
              {
                name:res.data.data.name,
                dec:res.data.data.dec,
                timeout:res.data.data.DLT645Timeout.toString(),
                configurationModel:res.data.data.configUid,
                configurationPageUUID:res.data.data.PageUUID,
              })
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
              type:40,
              gatherNumber:0,
              configUid:this.form.getFieldValue('configurationModel'),
              PageUUID:this.form.getFieldValue('configurationPageUUID'),
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
      this.$router.push('/DeviceModel/IEC104Model')
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