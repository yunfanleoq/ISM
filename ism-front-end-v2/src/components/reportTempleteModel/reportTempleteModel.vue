<template>
  <a-card>
    <a-modal v-drag-modal :width="'800px'" :footer="null" v-model="addVisible" :title="$t('diyReportTemplete.TempleteList')">
      <a-spin tip="Loading..." :spinning="loading">
        <a-table rowKey="Name" :pagination="pagination" :columns="columns" :data-source="dataSource" >
        <template v-for="(item, index) in columns" :slot="item.slotName">
          <span :key="index">{{ $t(item.slotName) }}</span>
        </template>
        <span slot="Status" slot-scope="Status">
                  <div v-if="Status==1" style="color: #74f808">
                    {{$t('configComponent.video.VideoOnline')}}
                  </div>
                 <div v-else-if="Status==0" style="color: #ea1111">
                    {{$t('configComponent.video.VideoOffline')}}
                 </div>
            </span>
        <div slot="Opt" slot-scope="text, record">
          <a type="link"   @click="SelectTemplete(record)" style="cursor: pointer;color: #13C2C2"><a-icon type="download" /><span style="margin-left: 2px;">{{$t('diyReportTemplete.Load')}}</span></a>
        </div>
      </a-table>
      </a-spin>
    </a-modal>
  </a-card>
</template>
<script>
import {GetReportTempletes} from "@/services/reportTemplete";
export default {
  name: 'deviceDataModel',
  i18n: require('../../i18n/language'),
  data() {
    return {
      pagination:{
        pageSize:15,
        showSizeChanger:true
      },
      loading:false,
      addVisible:false,
      columns: [
        {
          slotName: 'diyReportTemplete.Name',
          width: '20%',
          scopedSlots: { customRender: 'Name', title: 'diyReportTemplete.Name' },
          dataIndex: 'Name',
        },
        {
          width: '20%',
          slotName: 'diyReportTemplete.Describe',
          scopedSlots: { customRender: 'Describe', title: 'diyReportTemplete.Describe' },
          dataIndex: 'Describe',
        },
        {
          slotName: 'configComponent.video.TableVideoOpt',
          width: '10%',
          scopedSlots: { customRender: 'Opt',title: 'configComponent.video.TableVideoOpt'}
        }
      ],
      dataSource: [],
    };
  },
  watch: {
    '$route'() {

    }
  },
  mounted(){
    this.GetReportTempletes()
  },
  create(){
    this.GetReportTempletes()
  },
  methods: {
    showModal() {
      this.addVisible = true
      this.GetReportTempletes()
    },
    handleCancel() {
      this.addVisible = false;
    },
    SelectTemplete(record){
      this.$emit("OnSelectTemplete", record);
      this.addVisible = false;
    },
    GetReportTempletes(){
      let _t = this

      _t.loading = true
      this.dataSource=[]
      GetReportTempletes().then(function (res){
        _t.loading = false
        if (res.data.code == 0) {
          if(res.data.list==null)
          {
            _t.dataSource=[]
          }
          else
          {
            _t.dataSource = res.data.list
          }

        }
      }).finally(function (error) {
        _t.loading = false
      })
    }
  },
}
</script>

<style scoped lang="less">
::v-deep .ant-modal-body {
  padding: 5px;
  }
::v-deep .ant-select-dropdown {
  z-index: 900000;
}
  #DataGrid .datagrid-field-td{
    display: none;
  }
  #DataGrid .datagrid-header {
     border-style: unset;
  }
</style>
