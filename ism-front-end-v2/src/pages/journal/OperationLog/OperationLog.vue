<template>
  <a-card>
    <div >
      <a-form layout="horizontal">
        <div :class="advanced ? null: 'fold'">
          <a-row >
            <a-col :md="8" :sm="24" >
              <a-form-item
                  :label="$t('reporting.AlarmHistory.DateType')"
                  :labelCol="{span: 5}"
                  :wrapperCol="{span: 18, offset: 1}"
              >
                <a-radio-group v-model="SelectDateType" @change="chargeDateType">
                  <a-radio-button value="Day">
                    {{$t('reporting.AlarmHistory.DateDayType')}}
                  </a-radio-button>
                  <a-radio-button value="Weekly">
                    {{$t('reporting.AlarmHistory.DateWeeklyType')}}
                  </a-radio-button>
                  <a-radio-button value="Month">
                    {{$t('reporting.AlarmHistory.DateMonthType')}}
                  </a-radio-button>
                  <a-radio-button value="Diy">
                    {{$t('reporting.AlarmHistory.DateDiyType')}}
                  </a-radio-button>
                </a-radio-group>
              </a-form-item>
            </a-col>
            <a-col :md="8" :sm="24" >
              <a-form-item
                  :label="$t('reporting.AlarmHistory.SelectDate')"
                  :labelCol="{span: 5}"
                  :wrapperCol="{span: 18, offset: 1}"
              >
                <a-date-picker :defaultValue="moment()" style="width: 100%" @change="onDateChange"  size="default" :placeholder="$t('reporting.AlarmHistory.DateDayType')" v-if="SelectDateType=='Day'"/>
                <a-month-picker :defaultValue="moment()" style="width: 100%" @change="onDateChange" size="default" :placeholder="$t('reporting.AlarmHistory.DateMonthType')" v-if="SelectDateType=='Month'"/>
                <a-week-picker  :defaultValue="moment()"   style="width: 100%" @change="onWeeklyDateChange" size="default" :placeholder="$t('reporting.AlarmHistory.DateWeeklyType')" v-if="SelectDateType=='Weekly'"/>
                <a-range-picker :default-value="[moment().add(-1, 'day'),moment()]"  :showTime="true" @change="onDateChange" size="default" v-if="SelectDateType=='Diy'"/>
              </a-form-item>
            </a-col>
            <a-col :md="4" :sm="24" v-if="SelectDateType=='Diy'">
              <span style="float: right; margin-top: 3px;">
                <a-button :disabled="messageShowLoad" type="primary" @click="QueryJournal">{{$t('reporting.AlarmHistory.Query')}}</a-button>
              </span>
            </a-col>
            <a-col v-else :md="2" :sm="24" >
              <span style="float: right; margin-top: 3px;">
                <a-button :disabled="messageShowLoad" type="primary" @click="QueryJournal">{{$t('reporting.AlarmHistory.Query')}}</a-button>
              </span>
            </a-col>
            <a-col  :md="2" :sm="24" >
              <span style="float: left; margin-top: 3px;">
                <a-button :disabled="isLoadExecl" type="default" style="margin-left: 5px" @click="handleExport">{{$t('reporting.AlarmHistory.Export')}}</a-button>
              </span>
            </a-col>
          </a-row>
        </div>
      </a-form>
    </div>

    <a-spin style="padding: 1px;"  :spinning="messageShowLoad" tip="Loading...">
      <a-table  :pagination="pagination" :columns="columns" :data-source="dataSource" rowKey="Time">
      <template v-for="(item, index) in columns" :slot="item.slotName">
        <span :key="index">{{ $t(item.slotName) }}</span>
      </template>
        <span slot="Content" slot-scope="Content">
         {{Content|ContentSplit(that)}}
        </span>
        <span slot="Time" slot-scope="Time">
         {{Time|formatDate}}
        </span>
        <span slot="JournalLevel" slot-scope="JournalLevel">
         <span style="color: #44aaff" v-if="JournalLevel==1001">
           {{$t('journal.JournalLevelInfo')}}
         </span>
          <span style="color: #ed1329" v-else-if="JournalLevel==1003">
            {{$t('journal.JournalLevelError')}}
          </span>
          <span style="color: yellow" v-else-if="JournalLevel==1002">
            {{$t('journal.JournalLevelWarning')}}
          </span>
        </span>

    </a-table>
    </a-spin>

  </a-card>


</template>

<script>
import {GetOperationLog} from "../../../services/journal";
import moment from 'moment';

import {formatDate} from '@/utils/common';
import {exportExcel} from '@/services/excelExport'
import 'moment/locale/zh-cn';
import  'moment/locale/en-ie';
import  'moment/locale/zh-tw';
export default {
  name: 'OperationLogger',
  i18n: require('../../../i18n/language'),
  components:{
  },
  props: {
    ShowPageCount: {
      type: Number,
      default: 5
    }
  },
  data () {
    return {
      moment,
      pagination:{
        pageSize:15,
        showSizeChanger:true
      },
      loadExecl:null,
      isLoadExecl:false,
      exportName:"",
      json_fields_cn: {
        "内容": {
          field: "Content",
          //自定义回调函数
          callback: value => {
            let ContentSplit = value.split("&")
            let VarContent = ""
            for (let i=0;i<ContentSplit.length;i++)
            {
              if(i%2==0){
                VarContent = VarContent+this.$t(ContentSplit[i])
              }
              else{
                VarContent = VarContent+ContentSplit[i]
              }
            }

            return VarContent
          }
        },    //常规字段
        "时间": {
          field: "Time",
          //自定义回调函数
          callback: value => {
            let date = new Date(value)
            return formatDate(date,'yyyy-MM-dd hh:mm:ss')
          }
        },
        "等级": {
          field: "JournalLevel",
          //自定义回调函数
          callback: value => {
            switch (value){
              case '1001':{
                return this.$t('journal.JournalLevelInfo')
              }
              case '1002':{
                return this.$t('journal.JournalLevelError')
              }
              case '1003':{
                return this.$t('journal.JournalLevelWarning')
              }
            }
          }
        },
        "操作者": {
          field: "Operator",
          //自定义回调函数
          callback: value => {
            return this.$t(value)
          }
        },
        "请求信息": {
          field: "ClientInfo",
          //自定义回调函数
          callback: value => {
            return this.$t(value)
          }
        },
      },
      json_fields_en: {
        "Content": {
          field: "Content",
          //自定义回调函数
          callback: value => {
            let ContentSplit = value.split("&")
            let VarContent = ""
            for (let i=0;i<ContentSplit.length;i++)
            {
              if(i%2==0){
                VarContent = VarContent+this.$t(ContentSplit[i])
              }
              else{
                VarContent = VarContent+ContentSplit[i]
              }
            }

            return VarContent
          }
        },    //常规字段
        "Time": {
          field: "Time",
          //自定义回调函数
          callback: value => {
            let date = new Date(value)
            return formatDate(date,'yyyy-MM-dd hh:mm:ss')
          }
        },
        "Level": {
          field: "JournalLevel",
          //自定义回调函数
          callback: value => {
            switch (value){
              case '1001':{
                return this.$t('journal.JournalLevelInfo')
              }
              case '1002':{
                return this.$t('journal.JournalLevelError')
              }
              case '1003':{
                return this.$t('journal.JournalLevelWarning')
              }
            }
          }
        },
        "Operator": {
          field: "Operator",
          //自定义回调函数
          callback: value => {
            return this.$t(value)
          }
        },
        "RequestInfo": {
          field: "ClientInfo",
          //自定义回调函数
          callback: value => {
            return this.$t(value)
          }
        },
      },
      json_fields:{},
      json_meta: [
        [
          {
            " key ": " charset ",
            " value ": " utf- 8 "
          }
        ]
      ],
      SelectDateType: 'Day',
      SelectDevice:[],
      SelectDateRange:moment().format("YYYY-MM-DD"),
      SelectAlarmData:[],
      messageShowLoad:false,
      advanced: true,
      that:this,
      refIconLoading: false,
      columns: [
        {
          slotName: 'journal.OperationLog.content',
          width: '20%',
          scopedSlots: { customRender: 'Content', title: 'journal.OperationLog.content' },
          dataIndex: 'Content',
        },
        {
          width: '10%',
          slotName: 'journal.OperationLog.time',
          scopedSlots: { customRender: 'Time', title: 'journal.OperationLog.time' },
          dataIndex: 'Time',
        },
        {
          width: '10%',
          slotName: 'journal.OperationLog.JournalLevel',
          scopedSlots: { customRender: 'JournalLevel', title: 'journal.OperationLog.JournalLevel' },
          dataIndex: 'JournalLevel',
        },
        {
          width: '10%',
          slotName: 'journal.OperationLog.operator',
          scopedSlots: { customRender: 'operator', title: 'journal.OperationLog.operator' },
          dataIndex: 'Operator',
        },
        {
          width: '15%',
          slotName: 'journal.OperationLog.ClientInfo',
          scopedSlots: { customRender: 'ClientInfo', title: 'journal.OperationLog.ClientInfo' },
          dataIndex: 'ClientInfo',
        },
      ],
      dataSource: [],
      conditionExpress:"",
      selectedRows: []
    }
  },
  authorize: {
    // deleteRecord: 'delete'
  },
  mounted(){
    if(this.$i18n.locale=="CN")
    {
      this.json_fields = this.json_fields_cn
    }
    else
    {
      this.json_fields = this.json_fields_en
    }
    // if(this.$i18n.locale=="CN")
    // {
    //   this.lag = timeCN
    // }
    // else if(this.$i18n.locale=="HK")
    // {
    //   this.lag = timeTW
    // }
    // else
    // {
    //   this.lag = timeEn
    // }

    this.pagination.pageSize = this.ShowPageCount
    this.exportName = this.$t('journal.title')+"."+formatDate( new Date(),'yyyy-MM-dd hh:mm:ss')+".xls"
  },
  activated(){

  },
  filters: {
    formatDate(time) {
      let date = new Date(time)
      return formatDate(date,'yyyy-MM-dd hh:mm:ss')
    },
    ContentSplit(Content,that){
      let ContentSplit = Content.split("&")
      let VarContent = ""
      for (let i=0;i<ContentSplit.length;i++)
      {
        if(i%2==0){
          VarContent = VarContent+that.$t(ContentSplit[i])
        }
        else{
          VarContent = VarContent+ContentSplit[i]
        }
      }
      return VarContent
    },
  },
  created(){

  },
  watch: {

  },
  methods: {
    chargeDateType(e){
      let type = e.target.value
      if(type=="Day")
      {
        this.SelectDateRange = moment().format("YYYY-MM-DD")
      }
      else  if(type=="Weekly")
      {
        const startDate = moment().day(1).format('YYYY-MM-DD'); // 周一日期
        const endDate = moment().day(7).format('YYYY-MM-DD'); // 周日日期
        this.SelectDateRange = [startDate,endDate]
      }
      else  if(type=="Month")
      {
        this.SelectDateRange = moment().format("YYYY-MM")
      }
      else{
        this.SelectDateRange =[moment().add(-1, 'day'),moment()]
      }
    },
    async handleExport(){
      this.isLoadExecl = true
      this.loadExecl = this.$message.loading(this.$t("reporting.DataHistory.LoadingExecl"),0)
      try {
        await exportExcel(this.dataSource, this.json_fields, this.exportName.replace('.xls', ''))
      } finally {
        this.$message.destroy(this.loadExecl)
        this.isLoadExecl = false
      }
    },
    filterOption(input, option) {
      return (
          option.componentOptions.children[0].text.toLowerCase().indexOf(input.toLowerCase()) >= 0
      );
    },
    onDateChange(date, dateString){
      this.SelectDateRange = dateString
    },
    onWeeklyDateChange(date, dateString){
      const startDate = moment(date).day(1).format('YYYY-MM-DD'); // 周一日期
      const endDate = moment(date).day(7).format('YYYY-MM-DD'); // 周日日期
      this.SelectDateRange = [startDate,endDate]
    },
    QueryJournal(){
      let _t = this

      _t.dataSource = []
      const params = {
        dateType:this.SelectDateType,
        dateRange:this.SelectDateRange,
      }
      if(params.dateRange==""||(params.dateRange[0]==""))
      {
        this.$message.error(this.$t("reporting.AlarmHistory.SelectDateError"))
        return
      }
      this.messageShowLoad=true
      GetOperationLog(params).then(function (res){
        if(res.data.code==0)
        {
          _t.dataSource =res.data.list
        }
        _t.messageShowLoad=false
      })
    }
  }
}
</script>

<style lang="less">

.ant-select-dropdown-menu-item-group-title {
  height: 32px;
  padding: 0 12px;
  color: rgba(0, 0, 0, 0.45);
  font-size: 12px;
  background-color: #f8f8f8;
  line-height: 32px;
}

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