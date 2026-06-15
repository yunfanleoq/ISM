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
            :label="$t('dataModel.gatherNumber')"
            :labelCol="{span: 7}"
            :wrapperCol="{span: 10}"
        >
          <a-input autocomplete="autocomplete"

                   v-decorator="['gatherNumber', {rules: [{ required: true, message: $t('dataModel.gatherNumber'), whitespace: true}]}]"
          />
        </a-form-item>
        <a-form-item
            :label="$t('dataModel.opcuaModel.connection')"
            :labelCol="{span: 7}"
            :wrapperCol="{span: 10}"
        >
          <a-select @change="handleConnectionSelectChange"   autocomplete="autocomplete"
                    v-decorator="['opcuaConnectType', {rules: [{ required: true, message: $t('dataModel.opcuaModel.connection'), whitespace: true}]}]"
          >
            <a-select-option value="1">{{$t('dataModel.opcuaModel.connectionTcp')}}</a-select-option>
            <a-select-option value="2">{{$t('dataModel.opcuaModel.connectionHttps')}}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item
            :label="$t('dataModel.opcuaModel.SecurityPolicy')"
            :labelCol="{span: 7}"
            :wrapperCol="{span: 10}">

          <a-select @change="handleSecurityPolicyChange"   autocomplete="autocomplete"
                    v-decorator="['SecurityPolicy', {rules: [{ required: true, message: $t('dataModel.opcuaModel.SecurityPolicy'), whitespace: true}]}]"
          >
            <a-select-option value="1">{{$t('dataModel.opcuaModel.SecurityPolicyNone')}}</a-select-option>
            <a-select-option value="2">{{$t('dataModel.opcuaModel.SecurityPolicyBasic128Rsa15')}}</a-select-option>
            <a-select-option value="3">{{$t('dataModel.opcuaModel.SecurityPolicyBasic256')}}</a-select-option>
            <a-select-option value="4">{{$t('dataModel.opcuaModel.SecurityPolicyBasic256Sha256')}}</a-select-option>
            <a-select-option value="5">{{$t('dataModel.opcuaModel.SecurityPolicyAes128Sha256Rsa0aep')}}</a-select-option>
            <a-select-option value="6">{{$t('dataModel.opcuaModel.SecurityPolicyAes256Sha256RsaPss')}}</a-select-option>
          </a-select>

        </a-form-item>
        <a-form-item
            :label="$t('dataModel.opcuaModel.SecurityModes')"
            :labelCol="{span: 7}"
            :wrapperCol="{span: 10}">

          <a-select @change="handleSecurityModesChange"   autocomplete="autocomplete"
                    v-decorator="['SecurityModes', {rules: [{ required: true, message: $t('dataModel.opcuaModel.SecurityModes'), whitespace: true}]}]"
          >
            <a-select-option value="1">{{$t('dataModel.opcuaModel.SecurityModesNone')}}</a-select-option>
            <a-select-option value="2">{{$t('dataModel.opcuaModel.SecurityModesSign')}}</a-select-option>
            <a-select-option value="3">{{$t('dataModel.opcuaModel.SecurityModesSignAndEncrypt')}}</a-select-option>
          </a-select>

        </a-form-item>
        <a-form-item
            :label="$t('dataModel.opcuaModel.Authentication')"
            :labelCol="{span: 7}"
            :wrapperCol="{span: 10}">

          <a-select @change="handleAuthenticationChange"   autocomplete="autocomplete"
                    v-decorator="['Authentication', {rules: [{ required: true, message: $t('dataModel.opcuaModel.Authentication'), whitespace: true}]}]"
          >
            <a-select-option value="1">{{$t('dataModel.opcuaModel.AuthenticationAnonymous')}}</a-select-option>
            <a-select-option value="2">{{$t('dataModel.opcuaModel.AuthenticationUserPassword')}}</a-select-option>
            <a-select-option value="3">{{$t('dataModel.opcuaModel.AuthenticationCertificate')}}</a-select-option>
          </a-select>

        </a-form-item>
        <div v-if="Authentication==2">
          <a-form-item
              :label="$t('dataModel.opcuaModel.AuthenticationUser')"
              :labelCol="{span: 7}"
              :wrapperCol="{span: 10}">
            <a-input  v-decorator="['AuthenticationUser', {rules: [{ required: true, message: $t('dataModel.opcuaModel.AuthenticationUser'), whitespace: true}]}]"/>
          </a-form-item>
          <a-form-item
              :label="$t('dataModel.opcuaModel.AuthenticationPassword')"
              :labelCol="{span: 7}"
              :wrapperCol="{span: 10}">
            <a-input  v-decorator="['AuthenticationPassword', {rules: [{ required: true, message: $t('dataModel.opcuaModel.AuthenticationPassword'), whitespace: true}]}]"/>
          </a-form-item>
        </div>
        <div v-else-if="Authentication==3">
          <a-form-item
              :label="$t('dataModel.opcuaModel.CertificatePath')"
              :labelCol="{span: 7}"
              :wrapperCol="{span: 10}">
            <a-input v-model="CertificatePath" v-decorator="['CertificatePath', {rules: [{ required: true, message: $t('dataModel.opcuaModel.CertificatePath'), whitespace: true}]}]"/>
          </a-form-item>
          <a-form-item
              :label="$t('dataModel.opcuaModel.CertificatePrivateKey')"
              :labelCol="{span: 7}"
              :wrapperCol="{span: 10}">
            <a-input  v-model="CertificatePrivateKey" v-decorator="['CertificatePrivateKey', {rules: [{ required: true, message: $t('dataModel.opcuaModel.CertificatePrivateKey'), whitespace: true}]}]"/>
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
  OpcuaModelEdit,
} from "@/services/opcuaModel";

import {displayModelList, getDisplayModelLayerData} from "@/services/displayModel";

const dataSource= []
export default {
  name: 'OPCUAModelDetail',
  i18n: require('../../../i18n/language'),
  data () {
    return {
      error: '',
      COMList:[],
      CertificatePrivateKey:"",
      CertificatePath:"",
      OPCUAConnectType:1,
      SecurityModes:1,
      SecurityPolicy:1,
      Authentication:1,
      configurationModel:[],
      displayPageList:[],
      messageShowLoad:false,
      form: this.$form.createForm(this),
    }
  },
  activated() {

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
        _t.modbusConnectType = res.data.data.modbusConnectType
        _t.form.setFieldsValue(
            {
              name:res.data.data.name,
              dec:res.data.data.dec,
              configurationModel:res.data.data.configUid,
              configurationPageUUID:res.data.data.PageUUID,
              gatherNumber:res.data.data.gatherNumber.toString(),
              opcuaConnectType:res.data.data.OPCUAConnectType.toString(),
              SecurityPolicy:res.data.data.OPCUASecurityPolicies.toString(),
              SecurityModes:res.data.data.OPCUASecurityModes.toString(),
              Authentication:res.data.data.OPCUAAuthModes.toString(),
            })
        _t.Authentication = res.data.data.OPCUAAuthModes
        setTimeout(function (){
          if(res.data.data.OPCUAAuthModes==2)
          {
            _t.form.setFieldsValue(
                {
                  AuthenticationUser:res.data.data.OPCUAConnectUserName,
                  AuthenticationPassword:res.data.data.OPCUAConnectPassword,
                })
          }
          else if(res.data.data.OPCUAAuthModes==3)
          {
            _t.form.setFieldsValue(
                {
                  CertificatePath:res.data.data.OPCUACertificatePath,
                  CertificatePrivateKey:res.data.data.OPCUAPrivateKeyPath,
                })
          }
        },500)

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
              configUid:this.form.getFieldValue('configurationModel'),
              PageUUID:this.form.getFieldValue('configurationPageUUID'),
              type:3,
              gatherNumber:parseInt(this.form.getFieldValue('gatherNumber')),
              OPCUAConnectType:parseInt(this.form.getFieldValue('opcuaConnectType')),
              OPCUASecurityPolicies:parseInt(this.form.getFieldValue('SecurityPolicy')),
              OPCUASecurityModes:parseInt(this.form.getFieldValue('SecurityModes')),
              OPCUAAuthModes:parseInt(this.form.getFieldValue('Authentication')),
            }
          };

          if(params.data.OPCUAAuthModes==2)
          {
            params.data.OPCUAConnectUserName = this.form.getFieldValue('AuthenticationUser')
            params.data.OPCUAConnectPassword = this.form.getFieldValue('AuthenticationPassword')
          }
          else if(params.data.OPCUAAuthModes==3)
          {
            params.data.OPCUACertificatePath = this.form.getFieldValue('CertificatePath')
            params.data.OPCUAPrivateKeyPath = this.form.getFieldValue('CertificatePrivateKey')
          }
          OpcuaModelEdit(params).then(function (res){
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
      this.$router.push('/DeviceModel/OPCUAModel')
    },
    handleSelectChange(value) {
      this.modbusConnectType = value
    },
    handleConnectionSelectChange(value){
      this.OPCUAConnectType = value
    },
    handleSecurityModesChange(v){
      this.SecurityModes = v
    },
    handleSecurityPolicyChange(v){
      this.SecurityPolicy = v
    },
    handleAuthenticationChange(v){
      this.Authentication = v
    }
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