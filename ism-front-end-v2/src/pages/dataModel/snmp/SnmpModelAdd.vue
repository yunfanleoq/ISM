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
          :label="$t('dataModel.Port')"
          :labelCol="{span: 7}"
          :wrapperCol="{span: 10}"
      >
        <a-input autocomplete="autocomplete"

                    v-decorator="['port', {rules: [{ required: true, message: $t('dataModel.Port'), whitespace: true,initialValue:162}]}]"
        />
      </a-form-item>
<!--      <a-form-item-->
<!--          :label="$t('dataModel.TimeOut')"-->
<!--          :labelCol="{span: 7}"-->
<!--          :wrapperCol="{span: 10}"-->
<!--      >-->
<!--        <a-input autocomplete="autocomplete"-->
<!--                 -->
<!--                 v-decorator="['timeout', {rules: [{ required: true, message: $t('dataModel.TimeOut'), whitespace: true,initialValue:'1000'}]}]"-->
<!--        />-->
<!--      </a-form-item>-->
<!--      <a-form-item-->
<!--          :label="$t('dataModel.FailedTimes')"-->
<!--          :labelCol="{span: 7}"-->
<!--          :wrapperCol="{span: 10}"-->
<!--      >-->
<!--        <a-input autocomplete="autocomplete"-->
<!--                 -->
<!--                 v-decorator="['failedTimes', {rules: [{ required: true, message: $t('dataModel.FailedTimes'), whitespace: true,initialValue:'5'}]}]"-->
<!--        />-->
<!--      </a-form-item>-->
<!--      <a-form-item-->
<!--          :label="$t('dataModel.Interval')"-->
<!--          :labelCol="{span: 7}"-->
<!--          :wrapperCol="{span: 10}"-->
<!--      >-->
<!--        <a-input  style=""-->
<!--                  -->
<!--                  v-decorator="['interval', {rules: [{ required: true, message: $t('dataModel.Interval'), whitespace: true,initialValue:'1000'}]}]"-->
<!--        >-->
<!--        </a-input>-->
<!--      </a-form-item>-->
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
        <a-select @change="handleSelectChange" style="" autocomplete="autocomplete"

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
            :label="$t('dataModel.snmpSecurityLevel')"
            :labelCol="{span: 7}"
            :wrapperCol="{span: 10}"
        >
          <a-select style="" autocomplete="autocomplete" @change="securityLevelSelectChange"

                    v-decorator="['snmpSecurityLevel', {rules: [{ required: true, message: $t('dataModel.snmpSecurityLevel'), whitespace: true}]}]">
            <a-select-option value=1>{{$t('dataModel.snmpLevelNoAuthNoPrivacy')}}</a-select-option>
            <a-select-option value=2>{{$t('dataModel.snmpLevelAuthNoPrivacy')}}</a-select-option>
            <a-select-option value=3>{{$t('dataModel.snmpLevelAuthPrivacy')}}</a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item v-if="securityLevel==2||securityLevel==3"
            :label="$t('dataModel.snmpUserName')"
            :labelCol="{span: 7}"
            :wrapperCol="{span: 10}"
        >
          <a-input autocomplete="autocomplete"

                   v-decorator="['snmpUserName', {rules: [{ required: true, message: $t('dataModel.snmpUserName'), whitespace: true}]}]"/>
        </a-form-item>
        <a-form-item v-if="securityLevel==2||securityLevel==3"
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

        <a-form-item v-if="securityLevel==2||securityLevel==3"
            :label="$t('dataModel.snmpUserPassword')"
            :labelCol="{span: 7}"
            :wrapperCol="{span: 10}"
        >
          <a-input autocomplete="autocomplete"

                   v-decorator="['snmpUserPassword', {rules: [{ required: true, message: $t('dataModel.snmpUserPassword'), whitespace: true}]}]"/>
        </a-form-item>
        <a-form-item v-if="securityLevel==3"
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

        <a-form-item v-if="securityLevel==3"
            :label="$t('dataModel.snmpPrivacyPassword')"
            :labelCol="{span: 7}"
            :wrapperCol="{span: 10}"
        >
          <a-input autocomplete="autocomplete"

                   v-decorator="['snmpPrivacyPassword', {rules: [{ required: true, message: $t('dataModel.snmpPrivacyPassword'), whitespace: true}]}]"/>
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
import {snmpModelAdd} from "../../../services/snmpmodel";
import {displayModelList, getDisplayModelLayerData} from "@/services/displayModel";

export default {
  name: 'BasicForm',
  i18n: require('../../../i18n/language'),
  data () {
    return {
      error: '',
      messageShowLoad:false,
      configurationModel:[],
      displayPageList:[],
      form: this.$form.createForm(this),
      version:1,
      securityLevel:1,
      value: 1
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
      this.form.validateFields((err) => {
        if (!err) {
          this.messageShowLoad = true
          const params = {
            name:this.form.getFieldValue('name'),
            dec:this.form.getFieldValue('dec'),
            configUid:this.form.getFieldValue('configurationModel'),
            PageUUID:this.form.getFieldValue('configurationPageUUID'),
            type:1,
            port:parseInt(this.form.getFieldValue('port')),
            version:parseInt(this.form.getFieldValue('version')),
            writecomm:this.form.getFieldValue('writecomm'),
            readcomm:this.form.getFieldValue('readcomm'),
            gatherNumber:parseInt(this.form.getFieldValue('gatherNumber')),
            snmpUserName:this.form.getFieldValue('snmpUserName'),
            snmpSecurityLevel:parseInt(this.form.getFieldValue('snmpSecurityLevel')),
            snmpAuthAlgorithm:parseInt(this.form.getFieldValue('snmpAuthAlgorithm')),
            snmpUserPassword:this.form.getFieldValue('snmpUserPassword'),
            snmpPrivacyAlgorithm:parseInt(this.form.getFieldValue('snmpPrivacyAlgorithm')),
            snmpPrivacyPassword:this.form.getFieldValue('snmpPrivacyPassword'),
          };
          snmpModelAdd(params).then(this.afterLogin)
        }
      })
    },
    onBlackCLK(){
      this.$router.push('/DeviceModel/SnmpModel')
    },
    afterLogin(res) {
      this.messageShowLoad = false
      if (res.data.code == 2002) {
        this.$message.success(this.$t('dataModel.modelAddSuccess'), 3)
        this.$router.push('/DeviceModel/SnmpModel')
      }
      else if (res.data.code == 2001)
      {
        this.$message.error(this.$t('dataModel.modelNameRepeat'), 3)
      }
      else {
        this.$message.error(this.$t('dataModel.modelAddFailed'), 3)
      }
    },
    handleSelectChange(value) {
      this.version=parseInt(value)
    },
    securityLevelSelectChange(value) {
      this.securityLevel=parseInt(value)
    }
  }
}
</script>

<style >
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
  /*background: #f8f8f8;*/
  /*border-bottom: 1px solid #e8e8e8;*/
  transition: background .3s ease;
}
.ant-form-item {
  margin-bottom: 5px;

}
</style>