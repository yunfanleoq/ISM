<template>
  <a-card>
    <a-space class="operator">
      <a-button @click="addModel" type="primary" icon="plus">{{$t('dataModel.RESTFulData.AddModel')}}</a-button>
      <a-button @click="refresh()"  type="default" icon="sync" :loading="refIconLoading">{{$t("dataModel.refModel")}}</a-button>
      <a-button @click="handleExportAll()" type="default" icon="export" :loading="exportLoading">{{$t("dataModel.exportAllPoints")}}</a-button>
      <a-upload
          name="file"
          accept=".xlsx"
          :multiple="false"
          :customRequest="customImportAllRequest"
          :showUploadList="false"
          :disabled="importLoading || exportLoading"
          @change="handleImportAllChange"
      >
        <a-button type="default" icon="import" :loading="importLoading">{{$t("dataModel.importAllPoints")}}</a-button>
      </a-upload>
    </a-space>

    <a-spin style="padding: 1px;"  :spinning="messageShowLoad || exportLoading || importLoading" :tip="spinTip">
      <a-table rowKey="key" :pagination="pagination" :columns="columns" :data-source="dataSource">
        <template v-for="(item, index) in columns" :slot="item.slotName">
          <span :key="index">{{ $t(item.slotName) }}</span>
        </template>
        <template slot="no" slot-scope="text, record,index, column">
          {{index+1}}
        </template>
        <div slot="OPCUAConnectType" slot-scope="value">
          <span v-if="value==1"> {{$t('dataModel.opcuaModel.connectionTcp')}}</span>
          <span v-else-if="value==2"> {{$t('dataModel.opcuaModel.connectionHttps')}}</span>
        </div>
        <div slot="OPCUAAuthModes" slot-scope="value">
          <span v-if="value==1"> {{$t('dataModel.opcuaModel.AuthenticationAnonymous')}}</span>
          <span v-else-if="value==2"> {{$t('dataModel.opcuaModel.AuthenticationUserPassword')}}</span>
          <span v-else-if="value==3"> {{$t('dataModel.opcuaModel.AuthenticationCertificate')}}</span>
        </div>
        <div slot="action" slot-scope="text, record">
          <router-link to="" @click.native="GoToEdit(record)" style="color: #13C2C2"><a-icon type="edit" />{{$t('dataModel.modelDetail')}}</router-link> |
          <router-link :to="`/DeviceModel/VirtualDeviceData/${record.key}`"  style="color: darkorange"><a-icon type="unordered-list" />{{$t('dataModel.opcuaModel.NodeIDConfig')}}</router-link> |
          <a-popconfirm :title="$t('dataModel.deleteConfirm')" @confirm="deleteRecord(record.key)">
            <a-icon slot="icon" type="question-circle-o" style="color: red" />
            <a-icon type="delete" theme="twoTone" two-tone-color="#eb2f96"/><a style="color: #eb2f96">{{$t('dataModel.delete')}}</a>
          </a-popconfirm>
        </div>
      </a-table>
    </a-spin>

    <a-modal v-model="AddVisible" :title="IsEdit?$t('dataModel.RESTFulData.EditModel'):$t('dataModel.RESTFulData.AddModel')" @ok="onAddSubmit">
      <a-form :form="RESTFulForm" :label-col="{ span: 6 }" :wrapper-col="{ span: 15 }">
        <a-alert type="error" :closable="true" v-show="error" :message="error" showIcon style="margin-bottom: 24px;" />
        <a-form-item
            :label="$t('dataModel.RESTFulData.ModelName')"
        >
          <a-input  autocomplete="autocomplete"

                    v-decorator="['name', {rules: [{ required: true, message: $t('dataModel.static.DataName'), whitespace: true}]}]"
          />
        </a-form-item>
        <a-form-item
            :label="$t('device.deviceConfigurationModelName')"
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
            :label="$t('dataModel.static.DataDec')"
        >
          <Mtextarea   v-model="textAreValue"
                       rows="4"
                       :showWordLimit="true"
                       :maxLength="100"
                       :autoSize="false"
                       v-decorator="['description', { rules: [{ required: true, message: $t('dataModel.static.DataDec') }] }]"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>

<script>
import Mtextarea from "@/components/textarea";
import {displayModelList, getDisplayModelLayerData} from "@/services/displayModel";
import {
  VirtualDeviceModelAdd,
  VirtualDeviceModelDelete,
  VirtualDeviceModelEdit,
  VirtualDeviceModellist
} from "@/services/VirtualDeviceModel";
import { LOCALUPGATEALLVIRTUALDEVICEDATAMODEL, EXPORTALLVIRTUALDEVICEDATAMODEL } from "@/services/api";
import { AUTH_TYPE, getAuthorization } from "@/utils/request";
import axios from 'axios'

const IMPORT_ALL_TIMEOUT_MS = 2 * 60 * 60 * 1000

const alarmText = (value) => {
  switch (Number(value)) {
    case 1: return "是"
    default: return "否"
  }
}
const alarmLevelText = (value) => {
  switch (Number(value)) {
    case 1: return "次要"
    case 2: return "重要"
    case 3: return "紧急"
    case 4: return "致命"
    default: return "提示"
  }
}
const recordText = (value) => {
  switch (Number(value)) {
    case 1: return "是"
    default: return "否"
  }
}
const recordTypeText = (value) => {
  switch (Number(value)) {
    case 1: return "定时存储"
    case 2: return "即时存储"
    default: return "变化存储"
  }
}

export default {
  name: 'VirtualDeviceModel',
  i18n: require('../../../i18n/language'),
  data () {
    return {
      pagination:{
        pageSize:15,
        showSizeChanger:true
      },
      RESTFulForm:this.$form.createForm(this),
      AddVisible:false,
      messageShowLoad:false,
      advanced: true,
      refIconLoading: false,
      exportLoading: false,
      importLoading: false,
      importAllUrl: LOCALUPGATEALLVIRTUALDEVICEDATAMODEL,
      exportFields: {
        "模型名称": "modelName",
        "数据名称": "name",
        "权限(ReadOnly,ReadWrite)": "auth",
        "数据类型": "type",
        "单位": "unit",
        "转换关系": "conversionExpression",
        "是否告警(是,否)": { field: "alarm", callback: alarmText },
        "告警等级(提示、次要、重要、紧急、致命)": { field: "alarmLevel", callback: alarmLevelText },
        "告警消息": "AlarmMessage",
        "告警消除消息": "AlarmClearMessage",
        "报警触发值(0,1)": {
          field: "alarmOnValue",
          callback: value => (value === 0 || value === '0') ? '0' : '1',
        },
        "是否存储(是,否)": { field: "record", callback: recordText },
        "存储类型(变化存储、定时存储、即时存储)": { field: "RecordType", callback: recordTypeText },
        "定时时间": "recordInterval",
        "变化值": "RecordDataCharge",
        "描述": "Description",
        "模型类型(勿修改)": "modeltype",
        "数据ID(勿修改)": "uuid",
        "模型ID(勿修改)": "muid",
      },
      columns: [
        {
          width: '10%',
          slotName: 'dataModel.modelTableIndex',
          scopedSlots: { customRender: 'no', title: 'dataModel.modelTableIndex' },
          dataIndex: 'no'
        },
        {
          width: '20%',
          slotName: 'dataModel.modelName',
          scopedSlots: { customRender: 'modelName', title: 'dataModel.modelName' },
          dataIndex: 'modelName',
        },
        {
          slotName: 'dataModel.modelDec',
          width: '40%',
          scopedSlots: { customRender: 'modelDec', title: 'dataModel.modelDec' },
          dataIndex: 'modelDec',
        },
        {
          slotName: 'dataModel.modelTableOpt',
          scopedSlots: { customRender: 'action',title: 'dataModel.modelTableOpt'}
        }
      ],
      dataSource: [],
      configurationModel:[],
      displayPageList:[],
      IsEdit:false,
      EditUuid:"",
      selectedRows: [],
      textAreValue:"",
      error: '',
    }
  },
  authorize: {
    // deleteRecord: 'delete'
  },
  components: {
    Mtextarea,
  },
  computed: {
    spinTip() {
      if (this.importLoading) {
        return this.$t('dataModel.importAllLoading')
      }
      if (this.exportLoading) {
        return this.$t('dataModel.exportAllLoading')
      }
      return 'Loading...'
    },
    importHeaders() {
      const headers = {}
      const token = getAuthorization(AUTH_TYPE.BEARER)
      const projectUuid = getAuthorization(AUTH_TYPE.AUTH1)
      if (token) {
        headers.Authorization = token
      }
      if (projectUuid) {
        headers.ProjectUuid = projectUuid
      }
      return headers
    },
  },
  mounted(){

  },
  activated(){

  },
  created(){
    this.dataSource=[]
    this.getConfigurationModel()
    this.getModelList()
  },
  watch: {
    '$route' () {
      this.dataSource=[]

      this.getModelList()
    }
  },
  methods: {
    mapExportRow(item, model) {
      const source = {
        ...item,
        modelName: model.modelName,
        muid: model.key,
        modeltype: item.modeltype != null ? item.modeltype : 480,
      }
      const row = {}
      for (const key in this.exportFields) {
        const fieldConfig = this.exportFields[key]
        if (typeof fieldConfig === 'string') {
          row[key] = source[fieldConfig]
        } else if (typeof fieldConfig === 'object' && fieldConfig.field) {
          const rawValue = source[fieldConfig.field]
          row[key] = fieldConfig.callback ? fieldConfig.callback(rawValue) : rawValue
        }
      }
      return row
    },
    async handleExportAll() {
      if (this.exportLoading) return
      this.exportLoading = true
      try {
        const res = await axios.post(EXPORTALLVIRTUALDEVICEDATAMODEL, {}, {
          headers: {
            ...this.importHeaders,
          },
          responseType: 'blob',
          timeout: IMPORT_ALL_TIMEOUT_MS,
        })
        const data = res.data
        if (data instanceof Blob && data.type && data.type.indexOf('application/json') !== -1) {
          const text = await data.text()
          let msg = this.$t('dataModel.exportAllFailed')
          try {
            const json = JSON.parse(text)
            if (json.msg) msg = json.msg
          } catch (e) { /* ignore */ }
          this.$message.error(msg)
          return
        }
        const blob = new Blob([data], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
        const url = window.URL.createObjectURL(blob)
        const a = document.createElement('a')
        const disp = (res.headers && (res.headers['content-disposition'] || res.headers['Content-Disposition'])) || ''
        let fileName = `虚拟设备全量点位_${new Date().toISOString().slice(0, 19).replace(/[-:T]/g, '')}.xlsx`
        const m = /filename\*?=(?:UTF-8'')?["']?([^"';]+)/i.exec(disp)
        if (m && m[1]) {
          try { fileName = decodeURIComponent(m[1]) } catch (e) { fileName = m[1] }
        }
        a.href = url
        a.download = fileName
        document.body.appendChild(a)
        a.click()
        document.body.removeChild(a)
        window.URL.revokeObjectURL(url)
        this.$message.success(this.$t('dataModel.exportAllSuccess'))
      } catch (e) {
        this.$message.error(this.$t('dataModel.exportAllFailed'))
      } finally {
        this.exportLoading = false
      }
    },
    customImportAllRequest({ file, onSuccess, onError, onProgress }) {
      const formData = new FormData()
      formData.append('file', file)
      axios.post(this.importAllUrl, formData, {
        headers: {
          ...this.importHeaders,
        },
        timeout: IMPORT_ALL_TIMEOUT_MS,
        onUploadProgress: (e) => {
          if (e.total > 0 && onProgress) {
            onProgress({ percent: Math.round((e.loaded / e.total) * 100) }, file)
          }
        },
      }).then((res) => {
        onSuccess(res.data, file)
      }).catch((err) => {
        onError(err)
      })
    },
    handleImportAllChange(info) {
      if (info.file.status === 'uploading') {
        this.importLoading = true
        return
      }
      if (info.file.status === 'done') {
        this.importLoading = false
        const result = info.file.response || {}
        if (result.Code === 0) {
          const detail = result.message || `${this.$t('dataModel.importAllSuccess')}`
          this.$message.success(detail)
          this.getModelList()
        } else if (result.Code === -2 || result.Code === -6) {
          this.$message.error(result.message || this.$t('dataModel.importAllFormatError'))
        } else {
          this.$message.error(result.message || this.$t('dataModel.importAllFailed'))
        }
        return
      }
      if (info.file.status === 'error') {
        this.importLoading = false
        const err = info.file.error || {}
        let detail = this.$t('dataModel.importAllFailed')
        if (err.code === 'ECONNABORTED' || /timeout/i.test(err.message || '')) {
          detail = `${detail}（请求超时，请拆分 Excel 后重试）`
        } else if (err.response && err.response.status) {
          detail = `${detail}（HTTP ${err.response.status}）`
        } else if (err.message) {
          detail = `${detail}（${err.message}）`
        }
        this.$message.error(detail)
      }
    },
    addModel(){
      this.AddVisible=true
      this.IsEdit=false
      this.EditUuid=""
      this.textAreValue=""
      this.RESTFulForm.setFieldsValue(
          {
            name:"",
            description:"",
            configurationModel:"",
            configurationPageUUID:"",
          })
    },
    GoToEdit(item){
      let _t = this
      _t.AddVisible=true
      if(item.configUid!="")
      {
        this.GetDisplayPage(item.configUid)
      }

      setTimeout(function (){
        _t.IsEdit = true
        _t.EditUuid=item.key
        _t.textAreValue = item.modelDec
        _t.RESTFulForm.setFieldsValue(
            {
              name:item.modelName,
              description:item.modelDec,
              configurationModel:item.configUid,
              configurationPageUUID:item.PageUUID,
            })

      },500)
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
    onAddSubmit(e){
      e.preventDefault()
      let _t = this
      if (this.IsEdit)
      {
        this.RESTFulForm.validateFields((err) => {
          if (!err) {
            this.messageShowLoad = true
            const params = {
              uuid:_t.EditUuid,
              data: {
                name:this.RESTFulForm.getFieldValue('name'),
                dec:this.RESTFulForm.getFieldValue('description'),
                configUid:this.RESTFulForm.getFieldValue('configurationModel'),
                PageUUID:this.RESTFulForm.getFieldValue('configurationPageUUID'),
                type:480,
              }
            };
            VirtualDeviceModelEdit(params).then(function (res){
              _t.messageShowLoad = false
              if (res.data.code == 200) {
                _t.$message.success(_t.$t('dataModel.editSuccess'), 3)
              }
              else
              {
                _t.$message.error(_t.$t('dataModel.editFailed'), 3)
              }
              _t.getModelList()
              _t.AddVisible=false
            }).catch(function (){
              _t.messageShowLoad = false
              _t.$message.error(_t.$t('loginPage.serverError'), 3)
            })
          }
        })
      }
      else
      {
        this.RESTFulForm.validateFields((err) => {
          if (!err) {
            this.logging = true
            const params = {
              name:this.RESTFulForm.getFieldValue('name'),
              dec:this.RESTFulForm.getFieldValue('description'),
              configUid:this.RESTFulForm.getFieldValue('configurationModel'),
              PageUUID:this.RESTFulForm.getFieldValue('configurationPageUUID'),
              type:480,
            };
            VirtualDeviceModelAdd(params).then(function (res){
              _t.messageShowLoad = false
              if (res.data.code == 2002) {
                _t.$message.success(_t.$t('dataModel.modelAddSuccess'), 3)
              }
              _t.AddVisible=false
              _t.getModelList()
            }).catch(function (){
              _t.messageShowLoad = false
              _t.$message.error(_t.$t('loginPage.serverError'), 3)
            })
          }
        })
      }
    },
    refresh(){
      this.refIconLoading=true
      this.getModelList()
    },
    getModelList(){
      this.dataSource=[]
      let _t = this
      const  params= {
        type:480
      }
      this.messageShowLoad=true
      VirtualDeviceModellist(params).then(function (res){
        let tableData={}
        _t.refIconLoading=false
        _t.messageShowLoad=false
        if(res.data.list!=null)
        {
          for(let i=0;i<res.data.list.length;i++)
          {
            tableData.key = res.data.list[i].uuid
            tableData.no = i+1
            tableData.modelName = res.data.list[i].name
            tableData.modelDec = res.data.list[i].dec
            tableData.configUid =  res.data.list[i].configUid
            tableData.PageUUID =  res.data.list[i].PageUUID
            _t.dataSource.push(tableData)
            tableData={}
          }
        }

      }).catch(function (){
        _t.messageShowLoad = false
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
    },
    deleteRecord(key) {
      let params={
        uuid:key
      }
      let _t = this
      VirtualDeviceModelDelete(params).then(function (res) {
        if(res.data.code==200)
        {
          _t.dataSource = _t.dataSource.filter(item => item.key !== key)
          _t.selectedRows = _t.selectedRows.filter(item => item.key !== key)
        }
        else if(res.data.code==2004)
        {
          _t.$message.error(_t.$t("dataModel.modelBand"))
        }
      })
    },
  }
}
</script>

<style lang="less">
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
</style>