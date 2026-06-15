<template>
  <a-card>
    <a-button @click="addVisible=true;textAreValue='';SelectDevice=[];SelectAlarmData=7;SelectPeriod=3" type="primary">  {{$t('configComponent.video.AddVideo')}} </a-button>
    <a-divider type="vertical" />
    <a-button @click="GetReportTempletes()" :disabled="messageShowLoad" type="default" >  {{$t('configComponent.video.VideoRefresh')}} </a-button>
    <div style="margin-top: 5px">
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
          <a-button type="link" :disabled="IsExport"  @click="goToEdit(record)" style="margin-right:2px;padding:0;cursor: pointer;color: #13C2C2"><a-icon type="edit" /><span style="margin-left: 2px;">{{$t('monitor.nodeEdit')}}</span></a-button> |
          <a-button type="link"  :disabled="IsExport" @click="HandExport(record)" style="margin-right:2px;padding:0;cursor: pointer;color: #13C2C2"><a-icon type="export" /><span style="margin-left: 2px;">{{$t('diyReportTemplete.HandExport')}}</span></a-button> |
          <router-link :disabled="IsExport" :to="`/Reporting/ReportTempleteContent/${record.Uuid}`"  style="color: darkorange"><a-icon type="file-excel" />{{$t('diyReportTemplete.TempleteContent')}}</router-link> |
          <a-popconfirm
              :disabled="IsExport"
              :title="$t('configComponent.video.DelVideoConfirm')"
              :ok-text="$t('component.systemImageModel.delImageYes')"
              :cancel-text="$t('component.systemImageModel.delImageNo')"
              @confirm="DelTemplateAction(record.Uuid)"
          >
            <a-icon slot="icon" type="question-circle-o" style="color: red" />
            <a-icon type="delete" theme="twoTone" two-tone-color="#eb2f96"/><a style="color: #eb2f96">{{$t('configComponent.video.DelVideo')}}</a>
          </a-popconfirm>
        </div>
      </a-table>
    </div>
    <a-modal v-model="addVisible" :title="$t('configComponent.video.AddVideo')" @ok="AddTemplateAction">
      <a-form :form="AddForm" :label-col="{ span: 4 }" :wrapper-col="{ span: 19 }">
        <a-form-item
            :label="$t('diyReportTemplete.Name')"
        >
          <a-input
              v-decorator="[
                'Name',
                { rules: [{ required: true}] },
              ]"
          />
        </a-form-item>
        <a-form-item
            :label="$t('reporting.AlarmHistory.DeviceList')"
        >
          <a-tree-select
              show-search
              tree-node-filter-prop="title"
              v-model="SelectDevice"
              style="width: 100%"
              tree-checkable
              allow-clear
              :dropdown-style="{ maxHeight: '400px', overflow: 'auto' }"
              :tree-data="deviceTreeData"
              :replace-fields="{ value: 'key',title:'text'}"
              :placeholder="$t('reporting.AlarmHistory.DeviceList')"
              tree-default-expand-all
          >
          </a-tree-select>
        </a-form-item>
        <a-form-item
            :label="$t('diyReport.TimeGe')"
        >
          <a-select  style="width: 100%"  v-model="SelectAlarmData">
            <a-select-option v-for="(alarmItem,itemIndex) in TimeGeList" :key="itemIndex" :value=alarmItem.value>
              {{ $t(alarmItem.name) }}
            </a-select-option>
          </a-select>

        </a-form-item>
        <a-form-item
            :label="$t('diyReportTemplete.Period')"
        >
          <a-select  style="width: 100%"  v-model="SelectPeriod">
            <a-select-option v-for="(alarmItem,itemIndex) in Period" :key="itemIndex" :value=alarmItem.value>
              {{ $t(alarmItem.name) }}
            </a-select-option>
          </a-select>

        </a-form-item>
        <a-form-item
            :label="$t('diyReportTemplete.Describe')"
        >
          <Mtextarea   v-model="textAreValue"
                       rows="4"
                       :showWordLimit="true"
                       :maxLength="100"
                       :autoSize="false"
                       v-decorator="['Describe', { rules: [{ required: true, message: $t('diyReportTemplete.Describe') }] }]"
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal v-model="editVisible"  :title="$t('configComponent.video.AddVideo')" @ok="EditTempleteAction">
      <a-form :form="EditForm" :label-col="{ span: 4 }" :wrapper-col="{ span: 19 }">
        <a-form-item
            :label="$t('diyReportTemplete.Name')"
        >
          <a-input
              v-decorator="[
                'Name',
                { rules: [{ required: true}] },
              ]"
          />
        </a-form-item>
        <a-form-item
            :label="$t('reporting.AlarmHistory.DeviceList')"
        >
          <a-tree-select
              show-search
              tree-node-filter-prop="title"
              v-model="SelectDevice"
              style="width: 100%"
              tree-checkable
              allow-clear
              :dropdown-style="{ maxHeight: '400px', overflow: 'auto' }"
              :tree-data="deviceTreeData"
              :replace-fields="{ value: 'key',title:'text'}"
              :placeholder="$t('reporting.AlarmHistory.DeviceList')"
              tree-default-expand-all
          >
          </a-tree-select>
        </a-form-item>
        <a-form-item
            :label="$t('diyReport.TimeGe')"
        >
          <a-select  style="width: 100%"  v-model="SelectAlarmData">
            <a-select-option v-for="(alarmItem,itemIndex) in TimeGeList" :key="itemIndex" :value=alarmItem.value>
              {{ $t(alarmItem.name) }}
            </a-select-option>
          </a-select>

        </a-form-item>
        <a-form-item
            :label="$t('diyReportTemplete.Period')"
        >
          <a-select  style="width: 100%"  v-model="SelectPeriod">
            <a-select-option v-for="(alarmItem,itemIndex) in Period" :key="itemIndex" :value=alarmItem.value>
              {{ $t(alarmItem.name) }}
            </a-select-option>
          </a-select>

        </a-form-item>
        <a-form-item
            :label="$t('diyReportTemplete.Describe')"
        >
          <Mtextarea   v-model="textAreValue"
                       rows="4"
                       :showWordLimit="true"
                       :maxLength="100"
                       :autoSize="false"
                       v-decorator="['Describe', { rules: [{ required: true, message: $t('diyReportTemplete.Describe') }] }]"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>


</template>

<script>
import {formatDate} from '@/utils/common';
import Mtextarea from "@/components/textarea";
import {HandExportReportTemplete,GetReportTempletes,AddReportTemplete,DelReportTemplete,EditReportTemplete} from "@/services/reportTemplete";
import {getMonitorTree} from "@/services/device";
export default {
  name: 'diyDataHistoryTemplete',
  i18n: require('../../../i18n/language'),
  components:{
    Mtextarea
  },
  data () {
    return {
      messageShowLoad:false,
      pagination:{
        pageSize:15,
        showSizeChanger:true
      },
      textAreValue:"",
      editVisible:false,
      addVisible:false,
      SelectAlarmData:7,
      SelectPeriod:3,
      IsExport:false,
      TimeGeList:[
        {
          "value":1,
          "name":"diyReport.TimeGeOneMin"
        },
        {
          "value":2,
          "name":"diyReport.TimeGeFiveMin"
        },
        {
          "value":3,
          "name":"diyReport.TimeGeTenMin"
        },
        {
          "value":4,
          "name":"diyReport.TimeGeTenFiveMin"
        },
        {
          "value":5,
          "name":"diyReport.TimeGeTenThreeMin"
        },
        {
          "value":6,
          "name":"diyReport.TimeGeTenOneHour"
        },
        {
          "value":7,
          "name":"diyReport.TimeGeTenOneDay"
        },
      ],
      Period:[
        {
          "value":1,
          "name":"diyReportTemplete.PeriodList.oneDay"
        },
        {
          "value":2,
          "name":"diyReportTemplete.PeriodList.threeDay"
        },
        {
          "value":3,
          "name":"diyReportTemplete.PeriodList.sevenDay"
        },
        {
          "value":4,
          "name":"diyReportTemplete.PeriodList.fivethDay"
        },
        {
          "value":5,
          "name":"diyReportTemplete.PeriodList.onemonth"
        }
      ],
      AddForm:this.$form.createForm(this),
      EditForm:this.$form.createForm(this),
      columns: [
        {
          slotName: 'diyReportTemplete.Name',
          width: '15%',
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
          width: '20%',
          scopedSlots: { customRender: 'Opt',title: 'configComponent.video.TableVideoOpt'}
        }
      ],
      dataSource: [],
      SelectDevice:[],
      deviceTreeData:[],
    }
  },
  authorize: {
    // deleteRecord: 'delete'
  },
  mounted(){
    this.getMonitorTree()
    this.GetReportTempletes()
  },
  activated(){

  },
  filters: {
    formatDate(time) {
      let date = new Date(time)
      return formatDate(date,'yyyy-MM-dd hh:mm:ss')
    },
  },
  created(){

  },
  watch: {

  },
  methods: {
    getMonitorTree(){
      let _t = this
      this.deviceTreeData=[]
      getMonitorTree().then(function (res){
        if(res.data.code==0)
        {
          _t.deviceTreeData =res.data.list
        }
      })
    },
    HandExport(rc){
      let _t = this
      let params = {
        uuid:rc.Uuid
      };
      _t.$message.loading({ content: _t.$t('diyReportTemplete.ExportLoading'),duration: 0 });
      _t.IsExport=true
      HandExportReportTemplete(params).then(function (res){
        if(res.data.code==0)
        {
          _t.$message.destroy()
          const elink = document.createElement('a')
          elink.href = res.data.path
          elink.setAttribute('download', rc.Name+".zip")
          elink.style.display = 'none'
          document.body.appendChild(elink)
          setTimeout(() => {
            elink.click()
            document.body.removeChild(elink)
          }, 66)
        }else{
          _t.IsExport=false
          _t.$message.destroy()
          _t.$message.error(_t.$t('diyReportTemplete.HandExportFailed'))
        }
      }).finally(function (error) {
        _t.IsExport = false
      })
    },
    AddTemplateAction(){
      let _t = this
      this.AddForm.validateFields((err) => {
        if (!err) {
          let params = {
            Name:this.AddForm.getFieldValue('Name'),
            DeviceUuids:JSON.stringify(this.SelectDevice),
            TimeGe:this.SelectAlarmData,
            Period:this.SelectPeriod,
            Describe:this.textAreValue
          };
          AddReportTemplete(params).then(function (res){
            if(res.data.code==0)
            {
              _t.GetReportTempletes()
              _t.addVisible = false;
            }else{
              _t.$message.success(_t.$t('alarm.trigger.AddFailed'))
            }
          })
        }
      })
    },
    EditTempleteAction(){
      let _t = this
      this.EditForm.validateFields((err) => {
        if (!err) {
          let params = {
            uuid:_t.editUuid,
            data: {
              Name:this.EditForm.getFieldValue('Name'),
              DeviceUuids:JSON.stringify(this.SelectDevice),
              TimeGe:this.SelectAlarmData,
              Period:this.SelectPeriod,
              Describe:this.textAreValue
            }
          };
          EditReportTemplete(params).then(function (res){
            if(res.data.code==0)
            {
              _t.GetReportTempletes()
              _t.$message.success(_t.$t('alarm.trigger.EditSuccess'))
              _t.editVisible = false;
            }else{
              _t.$message.success(_t.$t('VideoManager.EditFailed'))
            }
          })
        }
      })
    },
    goToEdit(item){
      this.editUuid = item.Uuid
      this.textAreValue = item.Describe
      if(item.TimeGe==0)
      {
        this.SelectAlarmData=7
      }
      else {
        this.SelectAlarmData = item.TimeGe
      }
      if(item.Period==0){
        this.SelectPeriod=3
      }
      else {
        this.SelectPeriod = item.Period
      }
      try{
        this.SelectDevice = JSON.parse(item.DeviceUuids)
      }catch (e) {
        this.SelectDevice = []
      }

      let _t = this
      this.editVisible = true
      this.$message.loading(this.$t("monitor.loading"), 0.5)
      setTimeout(function (){
        _t.EditForm.setFieldsValue(
            {
              Name:item.Name,
              Describe:item.Describe,
            })
      },200)
    },
    DelTemplateAction(uuid){
      let _t = this

      let params = {
        Uuid:uuid
      };
      DelReportTemplete(params).then(function (res){
        if(res.data.code==0)
        {
          _t.GetReportTempletes()
          _t.$message.success(_t.$t('component.systemImageModel.delImageSuccess'))
          _t.addVisible = false;
        }
        else
        {
          _t.$message.error(_t.$t('component.systemImageModel.delImageFailed'))
        }
      })
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
  }
}
</script>

<style lang="less">

.plus-icon-enter-active{
  transition: opacity .5s;
}
.plus-icon-enter{
  opacity: 0;
}
.plus-icon-leave-active{
  transition: opacity .5s;
}
.plus-icon-leave-to{
  opacity: 0;
}
.plus-icon-enter-to{
  opacity: 1;
}

.code-box-actions {
  padding-top: 12px;
  text-align: center;
  opacity: .7;
  transition: opacity .3s;
}
.code-box-meta .demo-description>h4, .code-box-meta>h4 {
  position: absolute;
  top: -14px;
  padding: 1px 8px;
  margin-left: 16px;
  color: #777;
  border-radius: 2px 2px 0 0;
  background: #fff;
  font-size: 14px;
  width: auto;
}
.code-box {
  border: 1px solid #ebedf0;
  border-radius: 2px;
  display: inline-block;
  width: 100%;
  position: relative;
  margin: 0 0 16px;
  transition: all .2s;
}

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
