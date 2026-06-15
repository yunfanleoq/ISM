<template>
  <a-card>
    <div >
      <a-form layout="horizontal">
        <div :class="advanced ? null: 'fold'">
          <a-row >
            <a-col :md="8" :sm="24" >
              <a-form-item
                  :label="$t('reporting.AlarmHistory.DeviceList')"
                  :labelCol="{span: 5}"
                  :wrapperCol="{span: 18, offset: 1}"
              >
                <a-tree-select
                    tree-node-filter-prop="title"
                    show-search
                    v-model="SelectDevice"
                    @change="SelectTreeDevice"
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
            </a-col>
            <a-col :md="8" :sm="24" >
              <a-form-item
                  :label="$t('reporting.AlarmHistory.DataList')"
                  :labelCol="{span: 5}"
                  :wrapperCol="{span: 18, offset: 1}"
              >
                <a-select  @dropdownVisibleChange="GetDeviceModelDataList"
                           @popupScroll="handlePopupScroll"
                           @search="handleSearch" allowClear show-search optionFilterProp="children" mode="multiple"  style="width: 100%" :token-separators="[',']" v-model="SelectAlarmData">

                  <a-select-option v-for="(alarmItem,itemIndex) in frontDataZ" :key="itemIndex" :value=alarmItem.uuid>
                    {{ $t(alarmItem.name) }}
                  </a-select-option>
                </a-select>

              </a-form-item>
            </a-col>
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
          </a-row>
          <a-row >
            <a-col :md="8" :sm="24" >
              <a-form-item
                  :label="$t('reporting.AlarmHistory.SelectDate')"
                  :labelCol="{span: 5}"
                  :wrapperCol="{span: 18, offset: 1}"
              >
                <a-date-picker :defaultValue="moment()" style="width: 100%" @change="onDateChange"  size="default" :placeholder="$t('reporting.AlarmHistory.DateDayType')" v-if="SelectDateType=='Day'"/>
                <a-month-picker :defaultValue="moment()" style="width: 100%" @change="onDateChange" size="default" :placeholder="$t('reporting.AlarmHistory.DateMonthType')" v-if="SelectDateType=='Month'"/>
                <a-week-picker :defaultValue="moment()" style="width: 100%" @change="onWeeklyDateChange" size="default" :placeholder="$t('reporting.AlarmHistory.DateWeeklyType')" v-if="SelectDateType=='Weekly'"/>
                <a-range-picker :default-value="[moment().add(-1, 'day'),moment()]" :showTime="true" @change="onDateChange" size="default" v-if="SelectDateType=='Diy'"/>
              </a-form-item>
            </a-col>
            <a-col :md="3" :sm="24" v-if="SelectDateType=='Diy'">
              <span style="float: right; margin-top: 3px;">
                <a-button :disabled="messageShowLoad" type="primary" @click="QueryAlarmList">{{$t('reporting.AlarmHistory.Query')}}</a-button>
              </span>
            </a-col>
            <a-col v-else :md="2" :sm="24" >
              <span style="float: right; margin-top: 3px;">
                <a-button :disabled="messageShowLoad" type="primary" @click="QueryAlarmList">{{$t('reporting.AlarmHistory.Query')}}</a-button>
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
      <a-table :pagination="pagination" :columns="columns" :data-source="dataSource" rowKey="ID">
      <template v-for="(item, index) in columns" :slot="item.slotName">
        <span :key="index">{{ $t(item.slotName) }}</span>
      </template>
        <span slot="AlarmName" slot-scope="AlarmName">
         {{$t(AlarmName)}}
        </span>
        <span slot="HappenTime" slot-scope="HappenTime">
         {{HappenTime|formatDate}}
        </span>
        <span slot="ClearTime" slot-scope="ClearTime">
         {{ClearTime|formatDate}}
        </span>
        <span slot="AlarmMessage" slot-scope="AlarmMessage">
         {{$t(AlarmMessage)}}
        </span>

        <span slot="AlarmLevel" slot-scope="AlarmLevel">
          <a-tag style="background-color:#0099FF ;" v-if="AlarmLevel==0">
           {{ $t('dataModel.alarm.Tips')}}
         </a-tag>
          <a-tag style="background-color:#0099FF ;" v-if="AlarmLevel==1">
           {{ $t('dataModel.alarm.Minor')}}
         </a-tag>
          <a-tag style="background-color:yellow ;" v-if="AlarmLevel==2">
           {{ $t('dataModel.alarm.Importance')}}
         </a-tag>
          <a-tag style="background-color:orange ;" v-if="AlarmLevel==3">
           {{ $t('dataModel.alarm.Urgency')}}
         </a-tag>
          <a-tag style="background-color:red ;" v-else-if="AlarmLevel==4">
           {{ $t('dataModel.alarm.Deadly')}}
         </a-tag>
        </span>
    </a-table>
    </a-spin>

  </a-card>


</template>

<script>
import {getMonitorTree} from "@/services/device";
import {GetDeviceModelDataList} from "../../../services/device";
import {GetAlarmHistoryList} from "../../../services/report";
import moment from 'moment';

import {formatDate} from '@/utils/common';
import {exportExcel} from '@/services/excelExport'
import 'moment/locale/zh-cn';
import  'moment/locale/en-ie';
import  'moment/locale/zh-tw';

export default {
  name: 'AlarmHistory',
  i18n: require('../../../i18n/language'),
  components:{
  },
  data () {
    return {
      moment,
      dataZ:[],
      valueData: '',
      treePageSize: 100,
      scrollPage: 1,
      frontDataZ:[],
      pagination:{
        pageSize:15,
        showSizeChanger:true
      },
      loadExecl:null,
      isLoadExecl:false,
      exportName:"",
      json_fields_cn: {
        "名称": {
          field: "AlarmName",
          //自定义回调函数
          callback: value => {
            return this.$t(value)
          }
        },    //常规字段
        "设备": "DeviceName", //支持嵌套属性
        "告警时间": {
          field: "HappenTime",
          //自定义回调函数
          callback: value => {
            let date = new Date(value)
            return formatDate(date,'yyyy-MM-dd hh:mm:ss')
          }
        },
        "告警信息": {
          field: "AlarmMessage",
          //自定义回调函数
          callback: value => {
            return this.$t(value)
          }
        },
        "清除时间": {
          field: "ClearTime",
          //自定义回调函数
          callback: value => {
            let date = new Date(value)
            return formatDate(date,'yyyy-MM-dd hh:mm:ss')
          }
        },
        "告警等级": {
          field: "alarmLevel",
          //自定义回调函数
          callback: value => {
            switch (value){
              case 0:{
                return this.$t('dataModel.alarm.Tips')
              }
              case 1:{
                return this.$t('dataModel.alarm.Minor')
              }
              case 2:{
                return this.$t('dataModel.alarm.Importance')
              }
              case 3:{
                return this.$t('dataModel.alarm.Urgency')
              }
              case 4:{
                return this.$t('dataModel.alarm.Deadly')
              }
            }
          }
        },
        "保持时间": "KeepTime",
      },
      json_fields_en: {
        "AlarmName": {
          field: "AlarmName",
          //自定义回调函数
          callback: value => {
            return this.$t(value)
          }
        },    //常规字段
        "DeviceName": "DeviceName", //支持嵌套属性
        "AlarmMessage": {
          field: "AlarmMessage",
          //自定义回调函数
          callback: value => {
            return this.$t(value)
          }
        },
        "HappenTime": {
          field: "HappenTime",
          //自定义回调函数
          callback: value => {
            let date = new Date(value)
            return formatDate(date,'yyyy-MM-dd hh:mm:ss')
          }
        },
        "ClearTime": {
          field: "ClearTime",
          //自定义回调函数
          callback: value => {
            let date = new Date(value)
            return formatDate(date,'yyyy-MM-dd hh:mm:ss')
          }
        },
        "AlarmLevel": {
          field: "alarmLevel",
          //自定义回调函数
          callback: value => {
            switch (value){
              case 0:{
                return this.$t('dataModel.alarm.Tips')
              }
              case 1:{
                return this.$t('dataModel.alarm.Minor')
              }
              case 2:{
                return this.$t('dataModel.alarm.Importance')
              }
              case 3:{
                return this.$t('dataModel.alarm.Urgency')
              }
              case 4:{
                return this.$t('dataModel.alarm.Deadly')
              }
            }
          }
        },
        "KeepTime": "KeepTime",
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
      deviceTreeData:[],
      AlarmDataTree:[],
      form: this.$form.createForm(this),
      messageShowLoad:false,
      advanced: true,
      refIconLoading: false,
      columns: [
        {
          slotName: 'reporting.AlarmHistory.DeviceName',
          width: '10%',
          scopedSlots: { customRender: 'DeviceName', title: 'reporting.AlarmHistory.DeviceName' },
          dataIndex: 'DeviceName',
        },
        {
          width: '10%',
          slotName: 'reporting.AlarmHistory.AlarmName',
          scopedSlots: { customRender: 'AlarmName', title: 'reporting.AlarmHistory.AlarmName' },
          dataIndex: 'AlarmName',
        },
        {
          width: '10%',
          slotName: 'dataModel.editData.AlarmMessage',
          scopedSlots: { customRender: 'AlarmMessage', title: 'dataModel.editData.AlarmMessage' },
          dataIndex: 'AlarmMessage',
        },
        {
          width: '10%',
          slotName: 'reporting.AlarmHistory.HappenTime',
          scopedSlots: { customRender: 'HappenTime', title: 'reporting.AlarmHistory.HappenTime' },
          dataIndex: 'HappenTime',
        },
        {
          width: '10%',
          slotName: 'reporting.AlarmHistory.ClearTime',
          scopedSlots: { customRender: 'ClearTime', title: 'reporting.AlarmHistory.ClearTime' },
          dataIndex: 'ClearTime',
        },
        {
          width: '7%',
          slotName: 'reporting.AlarmHistory.AlarmLevel',
          scopedSlots: { customRender: 'AlarmLevel', title: 'reporting.AlarmHistory.AlarmLevel' },
          dataIndex: 'AlarmLevel',
        },
        {
          width: '7%',
          slotName: 'reporting.AlarmHistory.KeepTime',
          scopedSlots: { customRender: 'KeepTime', title: 'reporting.AlarmHistory.KeepTime' },
          dataIndex: 'KeepTime',
        }
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
    this.exportName = this.$t('reporting.AlarmHistory.exportName')+"."+formatDate( new Date(),'yyyy-MM-dd hh:mm:ss')+".xls"
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
    this.getMonitorTree()
    this.GetDeviceModelDataList()
  },
  watch: {
    '$route' () {
      this.getMonitorTree()
    }
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
    handleSearch (val) {
      this.valueData = val
      if (!val) {
        this.GetDeviceModelDataList()
      } else {
        this.frontDataZ = []
        this.scrollPage = 1
        this.dataZ.forEach(item => {
          if (item.name.indexOf(val) >= 0) {
            this.frontDataZ.push(item)
          }
        })
        this.allDataZ = this.frontDataZ
        this.frontDataZ = this.frontDataZ.slice(0, 100)
      }
    },
//下拉框下滑事件
    handlePopupScroll (e) {
      const { target } = e
      const scrollHeight = target.scrollHeight - target.scrollTop
      const clientHeight = target.clientHeight
      // 下拉框不下拉的时候
      if (scrollHeight === 0 && clientHeight === 0) {
        this.scrollPage = 1
        console.log(this.scrollPage)
      } else {
        // 当下拉框滚动条到达底部的时候
        if (scrollHeight < clientHeight + 5) {
          this.scrollPage = this.scrollPage + 1
          const scrollPage = this.scrollPage// 获取当前页
          const treePageSize = this.treePageSize * (scrollPage || 1)// 新增数据量
          const newData = [] // 存储新增数据
          let max = '' // max 为能展示的数据的最大条数
          if (this.dataZ.length > treePageSize) {
            // 如果总数据的条数大于需要展示的数据
            max = treePageSize
          } else {
            // 否则
            max = this.dataZ.length
          }
          // 判断是否有搜索
          if (this.valueData) {
            this.allDataZ.forEach((item, index) => {
              if (index < max) { // 当data数组的下标小于max时
                newData.push(item)
              }
            })
          } else {
            this.dataZ.forEach((item, index) => {
              if (index < max) { // 当data数组的下标小于max时
                newData.push(item)
              }
            })
          }

          this.frontDataZ = newData // 将新增的数据赋值到要显示的数组中
        }
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
    SelectTreeDevice(value,node,extera){
      this.GetDeviceModelDataList()
    },
    GetDeviceModelDataList(){
      let _t = this
      this.AlarmDataTree=[]
      _t.dataZ=[]
      _t.frontDataZ=[]
      const params ={
          SelectDevice:this.SelectDevice,
          getType:1
      }
      GetDeviceModelDataList(params).then(function (res){
        if(res.data.code==0)
        {
          for(let i=0;i<res.data.list.length;i++)
          {
            if(typeof res.data.list[i].DataList!="undefined"&&res.data.list[i].DataList!=null) {
              for (let j = 0; j < res.data.list[i].DataList.length; j++) {
                _t.dataZ.push(res.data.list[i].DataList[j])
                _t.AlarmDataTree.push(res.data.list[i].DataList[j])
              }
            }
          }
          _t.frontDataZ = _t.dataZ.slice(0, 100)
        }
      })
    },
    onDateChange(date, dateString){
      this.SelectDateRange = dateString
    },
    onWeeklyDateChange(date, dateString){
      const startDate = moment(date).day(1).format('YYYY-MM-DD'); // 周一日期
      const endDate = moment(date).day(7).format('YYYY-MM-DD'); // 周日日期
      this.SelectDateRange = [startDate,endDate]
    },
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
    QueryAlarmList(){
      let _t = this


      _t.dataSource = []
      const params = {
        deviceList:this.SelectDevice,
        dataList:this.SelectAlarmData,
        dateType:this.SelectDateType,
        dateRange:this.SelectDateRange,
      }
      if(params.dateRange==""||(params.dateRange[0]==""))
      {
        this.$message.error(this.$t("reporting.AlarmHistory.SelectDateError"))
        return
      }
      this.messageShowLoad=true
      GetAlarmHistoryList(params).then(function (res){
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