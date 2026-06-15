<template>
  <a-card>
    <a-space class="operator">
      <a-button @click="UserAdd" type="primary" icon="user-add">{{$t('dataModel.newModel')}}</a-button>
      <a-button @click="refresh()"  type="default" icon="sync" :loading="refIconLoading">{{$t("dataModel.refModel")}}</a-button>
    </a-space>
    <a-spin style="padding: 1px;"  :spinning="messageShowLoad" tip="Loading...">
      <a-table rowKey="ID"  :pagination="pagination" :columns="columns" :data-source="UserList">
        <template v-for="(item, index) in columns" :slot="item.slotName">
          <span :key="index">{{ $t(item.slotName) }}</span>
        </template>
        <div slot="role" slot-scope="text">
          <span v-if="text=='Operator'">{{ $t("account.settings.UserList.RoleOperator") }}</span>
          <span v-if="text=='User'">{{ $t("account.settings.UserList.RoleUser") }}</span>
        </div>
        <div slot="action" slot-scope="text, record">
          <a-popconfirm :title="$t('dataModel.deleteConfirm')" @confirm="UserDel(record)">
            <a-icon slot="icon" type="question-circle-o" style="color: red" />
            <a-icon type="delete" theme="twoTone" two-tone-color="#eb2f96"/><a style="color: #eb2f96">{{$t('dataModel.delete')}}</a>
          </a-popconfirm>
        </div>
      </a-table>
    </a-spin>
  </a-card>
</template>

<script>
import {SystemUserList,SystemUserDel} from "@/services/user";

export default {
  name: 'UserAdd',
  i18n: require('../../../i18n/language'),
  data () {
    return {
      pagination:{
        pageSize:15,
        showSizeChanger:true
      },
      // cropper
      UserList:[],
      refIconLoading:false,
      messageShowLoad:false,
      columns: [
        {
          width: '7%',
          slotName: 'dataModel.modelTableIndex',
          scopedSlots: { customRender: 'serial', title: 'dataModel.modelTableIndex' },
          dataIndex: 'ID'
        },
        {
          width: '15%',
          slotName: 'account.settings.basic.UserName',
          scopedSlots: { customRender: 'serial', title: 'account.settings.basic.UserName' },
          dataIndex: 'Username',
        },
        {
          width: '15%',
          slotName: 'account.settings.basic.name',
          scopedSlots: { customRender: 'serial', title: 'account.settings.basic.name' },
          dataIndex: 'name',
        },
        {
          slotName: 'account.settings.basic.position',
          width: '12%',
          scopedSlots: { customRender: 'serial', title: 'account.settings.basic.position' },
          dataIndex: 'job',
        },
        {
          width: '12%',
          slotName: 'account.settings.basic.phone',
          scopedSlots: { customRender: 'serial', title: 'account.settings.basic.phone' },
          dataIndex: 'phone',
        },
        {
          width: '15%',
          slotName: 'account.settings.basic.email',
          scopedSlots: { customRender: 'serial', title: 'account.settings.basic.email' },
          dataIndex: 'email',
        },
        {
          width: '10%',
          slotName: 'account.settings.basic.Role',
          scopedSlots: { customRender: 'role', title: 'account.settings.basic.Role' },
          dataIndex: 'role',
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
      this.SystemUserList()
    }
  },
  created(){
    this.SystemUserList()
  },
  methods: {
    refresh(){
      this.refIconLoading = true
      this.SystemUserList()
    },
    SystemUserList(){
      let _t = this
      _t.UserList = []
      this.messageShowLoad = true
      SystemUserList().then(function (res){
        if(res.data.code==0)
        {
          _t.UserList = res.data.List
        }
      }).finally(function (error) {
        _t.messageShowLoad = false
        _t.refIconLoading=false
      })

    },
    UserAdd(){
      this.$router.push('/Setting/UserAdd')
    },
    UserDel(item){
      let _t = this
     const params={
        id:item.ID,
        Name:item.name
     }
      _t.messageShowLoad = true
      SystemUserDel(params).then(function (res){
        if(res.data.code==200)
        {
          _t.SystemUserList()
        }
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
