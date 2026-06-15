<template>
  <a-card>
    <a-space class="operator">
      <a-button @click="TokenAdd" type="primary" icon="safety">{{$t('dataModel.newModel')}}</a-button>
      <a-button @click="refresh()"  type="default" icon="sync" :loading="refIconLoading">{{$t("dataModel.refModel")}}</a-button>
    </a-space>
    <a-spin style="padding: 1px;"  :spinning="messageShowLoad" tip="Loading...">
      <a-table rowKey="ID"  :pagination="pagination" :columns="columns" :data-source="TokenList">
        <template v-for="(item, index) in columns" :slot="item.slotName">
          <span :key="index">{{ $t(item.slotName) }}</span>
        </template>
        <div slot="CreatedAt" slot-scope="text">
          {{text|formatDate}}
        </div>
        <div slot="action" slot-scope="text, record">
          <a-popconfirm :title="$t('dataModel.deleteConfirm')" @confirm="TokenDel(record)">
            <a-icon slot="icon" type="question-circle-o" style="color: red" />
            <a-icon type="delete" theme="twoTone" two-tone-color="#eb2f96"/><a style="color: #eb2f96">{{$t('dataModel.delete')}}</a>
          </a-popconfirm>
        </div>
      </a-table>
    </a-spin>
  </a-card>
</template>

<script>
import {GetAccessTokenList,CreateAccessToken,DelAccessToken} from "@/services/user";
import {formatDate} from "@/utils/common";

export default {
  name: 'ApiToken',
  i18n: require('@/i18n/language'),
  data () {
    return {
      pagination:{
        pageSize:15,
        showSizeChanger:true
      },
      // cropper
      TokenList:[],
      refIconLoading:false,
      messageShowLoad:false,
      columns: [
        {
          width: '8%',
          slotName: 'dataModel.modelTableIndex',
          scopedSlots: { customRender: 'serial', title: 'dataModel.modelTableIndex' },
          dataIndex: 'ID'
        },
        {
          width: '65%',
          slotName: 'account.settings.token.token',
          scopedSlots: { customRender: 'AccessToken', title: 'account.settings.token.token' },
          dataIndex: 'AccessToken',
        },
        {
          slotName: 'account.settings.token.CreateDate',
          width: '15%',
          scopedSlots: { customRender: 'CreatedAt', title: 'account.settings.token.CreateDate' },
          dataIndex: 'CreatedAt',
        },
        {
          slotName: 'dataModel.modelTableOpt',
          scopedSlots: { customRender: 'action',title: 'dataModel.modelTableOpt'}
        }
      ],
    }
  },
  watch: {
    '$route'() {
      this.SystemTokenList()
    }
  },
  filters: {
    formatDate(time) {
      let date = new Date(time)
      return formatDate(date,'yyyy-MM-dd hh:mm:ss')
    },
  },
  created(){
    this.SystemTokenList()
  },
  methods: {
    refresh(){
      this.refIconLoading = true
      this.SystemTokenList()
    },
    SystemTokenList(){
      let _t = this
      _t.TokenList = []
      this.messageShowLoad = true
      GetAccessTokenList().then(function (res){
        if(res.data.code==0)
        {
          _t.TokenList = res.data.list
        }
      }).finally(function (error) {
        _t.messageShowLoad = false
        _t.refIconLoading=false
      })

    },
    TokenAdd(){
      let _t = this
      const params={

      }
      _t.messageShowLoad = true
      CreateAccessToken(params).then(function (res){
          _t.SystemTokenList()
      }).finally(function (error) {
        _t.messageShowLoad = false
        _t.refIconLoading=false
      })
    },
    TokenDel(item){
      let _t = this
     const params={
       accesstoken:item.AccessToken,
     }
      _t.messageShowLoad = true
      DelAccessToken(params).then(function (res){
          _t.SystemTokenList()
      }).finally(function (error) {
        _t.messageShowLoad = false
        _t.refIconLoading=false
      })
    }
  }
}
</script>

<style lang="less" scoped>
.operator{
  margin-bottom: 18px;
}
::v-deep .ant-table-thead > tr > th
{
  padding: 10px 10px;
  overflow-wrap: break-word;
}
::v-deep .ant-table-tbody > tr > td {
  padding: 7px 5px;
  overflow-wrap: break-word;
}

::v-deep .ant-table-thead>tr>th {
  color: #909399;
  font-weight: 500;
  text-align: left;
  //background: #f8f8f8;
  //border-bottom: 1px solid #e8e8e8;
  transition: background .3s ease;
}
</style>
