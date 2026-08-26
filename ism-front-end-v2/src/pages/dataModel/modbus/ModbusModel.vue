<template>
  <a-card>
      <a-space class="operator">
        <a-button @click="modbusAdd()" type="primary" icon="plus">{{$t('dataModel.newModel')}}</a-button>
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
      <a-table rowKey="modelName" :pagination="pagination" :columns="columns" :data-source="dataSource">
      <template v-for="(item, index) in columns" :slot="item.slotName">
        <span :key="index">{{ $t(item.slotName) }}</span>
      </template>

      <div slot="action" slot-scope="text, record">
        <router-link :to="`/DeviceModel/ModbusDetail/${record.key}`" style="color: #13C2C2"><a-icon type="edit" />{{$t('dataModel.modelDetail')}}</router-link> |
        <router-link :to="`/DeviceModel/ModbusRegister/${record.key}`"  style="color: darkorange"><a-icon type="import" />{{$t('dataModel.modbusModel.ModbusRegister')}}</router-link> |
        <a-popconfirm :title="$t('dataModel.deleteConfirm')" @confirm="deleteRecord(record.key)">
          <a-icon slot="icon" type="question-circle-o" style="color: red" />
          <a-icon type="delete" theme="twoTone" two-tone-color="#eb2f96"/><a style="color: #eb2f96">{{$t('dataModel.delete')}}</a>
        </a-popconfirm>
    </div>
    </a-table>
    </a-spin>
  </a-card>
</template>

<script>
import {
  DeviceModellist,
  modbusModelDelete,
  modbusModelGroupList,
  modbusModelRegisterList,
} from "../../../services/modbusModel";
import { LOCALUPGATEALLMODBUSDATAMODEL } from "@/services/api";
import { AUTH_TYPE, getAuthorization } from "@/utils/request";
import { exportExcelWithStyle } from "@/services/excelExport.js"
import axios from 'axios'

// 全量导入点位量大（OceanBase 尤其慢），需远长于全局 30s
const IMPORT_ALL_TIMEOUT_MS = 2 * 60 * 60 * 1000

const alarmText = (value) => {
  switch (value) {
    case 0: return "否"
    case 1: return "是"
    default: return "否"
  }
}
const alarmLevelText = (value) => {
  switch (value) {
    case 0: return "提示"
    case 1: return "次要"
    case 2: return "重要"
    case 3: return "紧急"
    case 4: return "致命"
    default: return "提示"
  }
}
const recordText = (value) => {
  switch (value) {
    case 0: return "否"
    case 1: return "是"
    default: return "否"
  }
}
const recordTypeText = (value) => {
  switch (Number(value)) {
    case 0: return "变化存储"
    case 1: return "定时存储"
    case 2: return "即时存储"
    case 3: return "变化百分比"
    case 4: return "整点存储"
    default: return "变化存储"
  }
}

export default {
  name: 'ModbusModelList',
  i18n: require('../../../i18n/language'),
  data () {
    return {
      pagination:{
        pageSize:15,
        showSizeChanger:true
      },
      messageShowLoad:false,
      advanced: true,
      refIconLoading: false,
      exportLoading: false,
      importLoading: false,
      importAllUrl: LOCALUPGATEALLMODBUSDATAMODEL,
      exportFields: {
        "模型名称": "modelName",
        "寄存器组名称": "registerGroupName",
        "数据名称": "name",
        "寄存器地址": "registerAddress",
        "权限(ReadOnly,ReadWrite)": "auth",
        "类型": "type",
        "字节序": "ByteOrder",
        "单位": "unit",
        "转换关系": "conversionExpression",
        "是否告警(是,否)": {
          field: "alarm",
          callback: alarmText,
        },
        "告警等级(提示、次要、重要、紧急、致命)": {
          field: "alarmLevel",
          callback: alarmLevelText,
        },
        "告警消息": "AlarmMessage",
        "告警消除消息": "AlarmClearMessage",
        "报警触发值(0,1)": {
          field: "alarmOnValue",
          callback: value => (value === 0 || value === '0') ? '0' : '1',
        },
        "是否存储(是,否)": {
          field: "record",
          callback: recordText,
        },
        "存储类型(变化存储、定时存储、即时存储、变化百分比、整点存储)": {
          field: "RecordType",
          callback: recordTypeText,
        },
        "定时时间": "recordInterval",
        "变化值": "RecordDataCharge",
        "保留小数": "FloatAccuracy",
        "模型类型(勿修改)": "modeltype",
        "组ID(勿修改)": "registerGroupUuid",
        "数据ID(勿修改)": "uuid",
        "模型ID(勿修改)": "muid",
      },
      columns: [
        {
          width: '10%',
          slotName: 'dataModel.modelTableIndex',
          scopedSlots: { customRender: 'serial', title: 'dataModel.modelTableIndex' },
          dataIndex: 'no'
        },
        {
          width: '20%',
          slotName: 'dataModel.modelName',
          scopedSlots: { customRender: 'serial', title: 'dataModel.modelName' },
          dataIndex: 'modelName',
        },
        {
          slotName: 'dataModel.modelDec',
          width: '30%',
          scopedSlots: { customRender: 'serial', title: 'dataModel.modelDec' },
          dataIndex: 'modelDec',
        },
        {
          width: '10%',
          slotName: 'dataModel.modbusModel.connection',
          scopedSlots: { customRender: 'serial', title: 'dataModel.modbusModel.connection' },
          dataIndex: 'modbusConnectType',
        },
        {
          width: '10%',
          slotName: 'dataModel.modbusModel.ModbusType',
          scopedSlots: { customRender: 'serial', title: 'dataModel.modbusModel.ModbusType' },
          dataIndex: 'modbusConnectMode',
        },
        {
          slotName: 'dataModel.modelTableOpt',
          scopedSlots: { customRender: 'action',title: 'dataModel.modelTableOpt'}
        }
      ],
      dataSource: [],
      selectedRows: []
    }
  },
  authorize: {
    // deleteRecord: 'delete'
  },
  mounted(){

  },
  activated(){

  },
  created(){
    this.dataSource=[]
    this.getModelList()
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
  watch: {
    '$route' () {
     this.dataSource=[]

     this.getModelList()
    }
  },
  methods: {
    refresh(){
      this.refIconLoading=true
      this.getModelList()
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
          detail = `${detail}（请求超时，请拆分 Excel 后重试，或直连后端导入）`
        } else if (err.response && err.response.status) {
          detail = `${detail}（HTTP ${err.response.status}）`
        } else if (err.message) {
          detail = `${detail}（${err.message}）`
        }
        this.$message.error(detail)
      }
    },
    mapExportRow(item, model, group) {
      const source = {
        ...item,
        modelName: model.modelName,
        registerGroupName: group.name,
        registerGroupUuid: group.uuid,
        muid: model.key,
        modeltype: item.modeltype != null ? item.modeltype : 2,
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
      if (this.exportLoading) {
        return
      }
      if (!this.dataSource.length) {
        this.$message.warning(this.$t('dataModel.exportAllEmpty'))
        return
      }
      this.exportLoading = true
      try {
        const allRows = []
        for (const model of this.dataSource) {
          const groupRes = await modbusModelGroupList({ muid: model.key })
          const groups = groupRes.data.list || []
          for (const group of groups) {
            const regRes = await modbusModelRegisterList({ uuid: group.uuid })
            const regs = regRes.data.list || []
            for (const item of regs) {
              allRows.push(this.mapExportRow(item, model, group))
            }
          }
        }
        if (!allRows.length) {
          this.$message.warning(this.$t('dataModel.exportAllNoPoints'))
          return
        }
        const stamp = new Date().toISOString().slice(0, 19).replace(/[-:T]/g, '')
        await exportExcelWithStyle(allRows, this.exportFields, `Modbus全量点位_${stamp}`, '', false)
        this.$message.success(`${this.$t('dataModel.exportAllSuccess')}${allRows.length}`)
      } catch (e) {
        this.$message.error(this.$t('dataModel.exportAllFailed'))
      } finally {
        this.exportLoading = false
      }
    },
    getModelList(){
      this.dataSource=[]
      let _t = this
      const  params= {
        type:2
      }
      this.messageShowLoad=true
      DeviceModellist(params).then(function (res){
        let tableData={}
        _t.refIconLoading=false
        _t.messageShowLoad=false
        if(res.data.list!=null)
        {
          for(let i=0;i<res.data.list.length;i++)
          {
            tableData.key = res.data.list[i].uuid
            tableData.no = res.data.list[i].ID
            tableData.modelName = res.data.list[i].name
            tableData.modelDec = res.data.list[i].dec
            tableData.modbusConnectType = res.data.list[i].modbusConnectType
            tableData.modbusConnectMode = res.data.list[i].modbusConnectMode
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
      modbusModelDelete(params).then(function (res) {
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
    modbusAdd(){
      this.$router.push('/DeviceModel/ModbusAdd')
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
