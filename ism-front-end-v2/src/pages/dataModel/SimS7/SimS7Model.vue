<template>
  <a-card>
      <a-space class="operator">
        <a-button @click="S7ModelAdd()" type="primary" icon="plus">{{$t('dataModel.newModel')}}</a-button>
        <a-button @click="refresh()"  type="default" icon="sync" :loading="refIconLoading">{{$t("dataModel.refModel")}}</a-button>
      </a-space>

    <a-spin style="padding: 1px;"  :spinning="messageShowLoad" tip="Loading...">
      <a-table rowKey="modelName" :pagination="pagination" :columns="columns" :data-source="dataSource">
      <template v-for="(item, index) in columns" :slot="item.slotName">
        <span :key="index">{{ $t(item.slotName) }}</span>
      </template>
        <template slot="serial" slot-scope="text, record,index, column">
          {{index+1}}
        </template>
      <div slot="action" slot-scope="text, record">
        <router-link :to="`/DeviceModel/SimS7Detail/${record.key}`" style="color: #13C2C2"><a-icon type="edit" />{{$t('dataModel.modelDetail')}}</router-link> |
        <router-link :to="`/DeviceModel/S7DataList/${record.key}`"  style="color: darkorange"><a-icon type="import" />{{$t('dataModel.opcuaModel.NodeIDConfig')}}</router-link> |
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
import {DeviceModellist,modbusModelDelete} from "../../../services/modbusModel";
import {SimS7Modellist,SimS7ModelDelete} from "@/services/SimS7";
export default {
  name: 'SimS7ModelList',
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
      columns: [
        {
          width: '10%',
          slotName: 'dataModel.modelTableIndex',
          scopedSlots: { customRender: 'serial', title: 'dataModel.modelTableIndex' },
          dataIndex: 'no'
        },
        {
          width: '30%',
          slotName: 'dataModel.modelName',
          scopedSlots: { customRender: 'modelName', title: 'dataModel.modelName' },
          dataIndex: 'modelName',
        },
        {
          slotName: 'dataModel.modelDec',
          width: '35%',
          scopedSlots: { customRender: 'modelDec', title: 'dataModel.modelDec' },
          dataIndex: 'modelDec',
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
    getModelList(){
      this.dataSource=[]
      let _t = this
      const  params= {
        type:15
      }
      this.messageShowLoad=true
      SimS7Modellist(params).then(function (res){
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
      SimS7ModelDelete(params).then(function (res) {
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
    S7ModelAdd(){
      this.$router.push('/DeviceModel/SimS7Add')
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
