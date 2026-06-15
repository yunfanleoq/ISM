<template>
  <a-card :body-style="{padding: '24px 32px'}" :bordered="false">
    <a-form @submit="onSubmit" :form="form">
      <a-alert type="error" :closable="true" v-show="error" :message="error" showIcon style="margin-bottom: 24px;" />
      <a-form-item
          :label="$t('dataModel.static.DataName')"
          :labelCol="{span: 7}"
          :wrapperCol="{span: 10}"
      >
        <a-input  autocomplete="autocomplete"

                  v-decorator="['name', {rules: [{ required: true, validator: isValidateTxtNonSpec, message: $t('device.deviceNameVal'), whitespace: true}]}]"
        />
      </a-form-item>
      <a-form-item
          :label="$t('dataModel.static.DataType')"
          :labelCol="{span: 7}"
          :wrapperCol="{span: 10}"
      >
        <a-select  style="" autocomplete="autocomplete"

                   v-decorator="['DataType', {initialValue:'1',rules: [{ required: true, message: $t('dataModel.static.DataType'), whitespace: true}]}]"
        >
          <a-select-option v-for="options in DataTypeList" :key="options.value" :value="options.value.toString()">
            {{ $t(options.name) }}
          </a-select-option>
        </a-select>

      </a-form-item>
      <a-form-item
          :labelCol="{span: 7}"
          :wrapperCol="{span: 10}"
      >
         <span slot="label">
          {{$t('dataModel.static.DataDeviceType')}}&nbsp;
          <a-tooltip :title="$t('dataModel.static.DataDevicePublicTips')">
            <a-icon type="question-circle-o" />
          </a-tooltip>
      </span>
        <a-select  style="" autocomplete="autocomplete"

                   v-decorator="['DataDeviceType', {initialValue:'158',rules: [{ required: true, message: $t('dataModel.static.DataDeviceType'), whitespace: true}]}]"
        >
          <a-select-option  v-for="(device,index) in supportDeviceList" :key="index" :value=device.type.toString()>
            {{ $t(device.name) }}
          </a-select-option>
        </a-select>

      </a-form-item>
      <a-form-item
          :label="$t('dataModel.static.DataDefaultValue')"
          :labelCol="{span: 7}"
          :wrapperCol="{span: 10}"
      >
        <a-input autocomplete="autocomplete"

                 v-decorator="['DataDefaultValue', {rules: [{ required: true, message: $t('dataModel.static.DataDefaultValue'), whitespace: true,initialValue:162}]}]"
        />
      </a-form-item>
      <a-form-item
          :label="$t('dataModel.static.DataUnit')"
          :labelCol="{span: 7}"
          :wrapperCol="{span: 10}"
      >
        <a-input autocomplete="autocomplete"

                 v-decorator="['DataUnit', {rules: [{ required: true, message: $t('dataModel.static.DataUnit'), whitespace: true}]}]"
        />
      </a-form-item>
      <a-form-item
          :label="$t('dataModel.static.DataDec')"
          :labelCol="{span: 7}"
          :wrapperCol="{span: 10}"
      >
        <Mtextarea   v-model="textAreValue"
                     rows="4"
                     :showWordLimit="true"
                     :maxLength="100"
                     :autoSize="false"
                     v-decorator="['description', { rules: [{ required: true, message: $t('dataModel.static.DataDec') }] }]"
        />
      </a-form-item>
      <a-form-item style="margin-top: 24px" :wrapperCol="{span: 10, offset: 7}">
        <a-button type="primary" htmlType="submit">{{$t('dataModel.add')}}</a-button>
        <a-button style="margin-left: 8px" @click="onBlackCLK()">{{$t('dataModel.back')}}</a-button>
      </a-form-item>
    </a-form>
  </a-card>
</template>

<script>
import {snmpModelAdd} from "../../../services/snmpmodel";
import {getSupportDeviceList} from "../../../services/device";
import ProjectLayout from "../../../layouts/ProjectLayout";
import Mtextarea from '@/components/textarea/index'
import {StaticModelAdd} from "../../../services/staticmodel";
export default {
  name: 'StaticModelDetail',
  i18n: require('../../../i18n/language'),
  data () {
    return {
      error: '',
      form: this.$form.createForm(this),
      version:1,
      textAreValue:"",
      supportDeviceList:[],
      securityLevel:1,
      DataTypeList:[
        {
          name:this.$t('dataModel.static.DataTypeInt'),
          value:1
        },
        {
          name:this.$t('dataModel.static.DataTypeString'),
          value:2
        },
        {
          name:this.$t('dataModel.static.DataTypeDouble'),
          value:3
        },
        {
          name:this.$t('dataModel.static.DataTypeJson'),
          value:4
        }
      ],
      value: 1
    }
  },
  components: {Mtextarea},
  computed: {
    desc() {
      return this.$t('pageDesc')
    }
  },
  mounted(){
    this.getSupportDevice()
  },
  methods: {
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
    getSupportDevice(){
      let _t = this
      getSupportDeviceList().then(function (res){
        let publicDevice = {
          name: "dataModel.static.DataDevicePublic",
          type: '158'
        }
        for(let i=0;i<res.data.list.length;i++)
        {
          if(res.data.list[i].type!=7&&res.data.list[i].type!=6)
          {
            _t.supportDeviceList.push(res.data.list[i])
          }
        }
        _t.supportDeviceList.push(publicDevice)
      })
    },
    onSubmit (e) {
      e.preventDefault()
      this.form.validateFields((err) => {
        if (!err) {
          this.logging = true
          const params = {
            Name:this.form.getFieldValue('name'),
            DataDeviceType:parseInt(this.form.getFieldValue('DataDeviceType')),
            DataType:parseInt(this.form.getFieldValue('DataType')),
            DataDefaultValue:this.form.getFieldValue('DataDefaultValue'),
            DataUnit:this.form.getFieldValue('DataUnit'),
            DataDescription:this.form.getFieldValue('description'),
          };
          StaticModelAdd(params).then(this.addResponse)
        }
      })
    },
    onBlackCLK(){
      this.$router.push('/DeviceModel/StaticData')
    },
    addResponse(res) {
      this.logging = false
      if (res.data.code == 0) {
        this.$message.success(this.$t('dataModel.modelAddSuccess'), 3)
        this.$router.push('/DeviceModel/StaticData')
      }
      else if (res.data.code == 3001)
      {
        this.$message.error(this.$t('dataModel.modelNameRepeat'), 3)
      }
      else {
        this.$message.error(this.$t('dataModel.modelAddFailed'), 3)
      }
    },
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
