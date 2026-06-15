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
  </a-card>
  </a-spin>
</template>

<script>
import {OpcuaModelAdd, OpcuaModelDelete,OpcuaModelEdit,OpcuaModellist} from "@/services/opcuaModel";
import {displayModelList, getDisplayModelLayerData} from "@/services/displayModel";

export default {
  name: 'OpcuaModelAdd',
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
      version:1,
      value: 1
    }
  },
  mounted(){
    this.getConfigurationModel()
    this.$nextTick(function (){
      this.form.setFieldsValue(
          {
            gatherNumber:"30",
            opcuaConnectType:"1",
            SecurityPolicy:"1",
            SecurityModes:"1",
            Authentication:"1"
          })
    })

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
            configUid:this.form.getFieldValue('configurationModel'),
            PageUUID:this.form.getFieldValue('configurationPageUUID'),
            type:3,
            gatherNumber:parseInt(this.form.getFieldValue('gatherNumber')),
            OPCUAConnectType:parseInt(this.form.getFieldValue('opcuaConnectType')),
            OPCUASecurityPolicies:parseInt(this.form.getFieldValue('SecurityPolicy')),
            OPCUASecurityModes:parseInt(this.form.getFieldValue('SecurityModes')),
            OPCUAAuthModes:parseInt(this.form.getFieldValue('Authentication')),
          };
          if(params.OPCUAAuthModes==2)
          {
            params.OPCUAConnectUserName = this.form.getFieldValue('AuthenticationUser')
            params.OPCUAConnectPassword = this.form.getFieldValue('AuthenticationPassword')
          }
          else if(params.OPCUAAuthModes==3)
          {
            params.OPCUACertificatePath = this.form.getFieldValue('CertificatePath')
            params.OPCUAPrivateKeyPath = this.form.getFieldValue('CertificatePrivateKey')
          }
          OpcuaModelAdd(params).then(function (res){
            _t.messageShowLoad = false
            if (res.data.code == 2002) {
              _t.$message.success(_t.$t('dataModel.modelAddSuccess'), 3)
              _t.$router.push('/DeviceModel/OPCUAModel')
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

<style scoped>
.ant-form-item {
  margin-bottom: 5px;

}
</style>