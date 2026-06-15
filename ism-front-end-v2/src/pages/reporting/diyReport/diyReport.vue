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
                    show-search
                    tree-node-filter-prop="title"
                    v-model="SelectDevice"
                    style="width: 100%"
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
                  :label="$t('diyReport.TimeGe')"
                  :labelCol="{span: 5}"
                  :wrapperCol="{span: 18, offset: 1}"
              >
                <a-select  style="width: 100%"  v-model="SelectAlarmData">
                    <a-select-option v-for="(alarmItem,itemIndex) in TimeGeList" :key="itemIndex" :value=alarmItem.value>
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
            <a-col :md="4" :sm="24" v-if="SelectDateType=='Diy'">
              <span style="float: right; margin-top: 3px;">
                <a-button :disabled="messageShowLoad" type="primary" @click="QueryHistoryDataList">{{$t('reporting.AlarmHistory.Query')}}</a-button>
              </span>
            </a-col>
            <a-col v-else :md="2" :sm="24" >
              <span style="float: right; margin-top: 3px;">
                <a-button :disabled="messageShowLoad" type="primary" @click="QueryHistoryDataList">{{$t('reporting.AlarmHistory.Query')}}</a-button>
              </span>
            </a-col>

            <a-col  :md="2" :sm="24" >
              <span style="margin-left:10px;float: left; margin-top: 3px;">
                <a-button  type="default"  icon="download" @click="loadTemplete">{{$t('diyReportTemplete.Load')}}</a-button>
              </span>
            </a-col>
<!--            <a-col  :md="2" :sm="24" >-->
<!--              <span style="margin-left:40px;float: left; margin-top: 3px;">-->
<!--                <a-button  type="dashed"  icon="upload" @click="exportExcelBtn">{{$t('diyReportTemplete.Export')}}</a-button>-->
<!--              </span>-->
<!--            </a-col>-->
            <a-col  :md="10" :sm="24" >
              <span style="float: right; margin-top: 3px;">
                <a-alert :message="TempleteName" type="info" show-icon v-if="TempleteName!=''"/>
              </span>
            </a-col>
          </a-row>
        </div>
      </a-form>
    </div>
    <div class="hello">
      <div
          id="luckysheetQuery"
          style="margin:0px;padding:0px;height:600px;width:100%;left: 0px;top: 100px;bottom:0px;"
      ></div>

      <report-templete-model @OnSelectTemplete="onSelectTemplete" ref="reportTempleteModel"></report-templete-model>
    </div>
  </a-card>


</template>

<script>
import {getMonitorTree} from "@/services/device";
import { GetDiyDataHistoryList} from "../../../services/report";
import moment from 'moment';
import LuckyExcel from 'luckyexcel'
import {formatDate} from '@/utils/common';
import { exportExcel,getExcelData } from '@/utils/export'
import reportTempleteModel from "@/components/reportTempleteModel/reportTempleteModel";
import 'moment/locale/zh-cn';
import  'moment/locale/en-ie';
import  'moment/locale/zh-tw';

export default {
  name: 'diyDataHistory',
  i18n: require('../../../i18n/language'),
  components:{
    reportTempleteModel
  },
  data () {
    return {
      moment,
      pagination:{
        pageSize:15,
        showSizeChanger:true
      },
      TempleteName:"",
      isMaskShow: false,
      selectRow:0,
      selectCol:0,
      loadExecl:null,
      isLoadExecl:false,
      exportName:"",
      json_fields_cn: {
        "数据名称": "DataName",    //常规字段
        "设备": "DeviceName", //支持嵌套属性
        "数据值": "DataValue",
        "数据单位": "DataUnit",
        "记录时间": {
          field: "RecordTime",
          //自定义回调函数
          callback: value => {
            let date = new Date(value)
            return formatDate(date,'yyyy-MM-dd hh:mm:ss')
          }
        }
      },
      json_fields_en: {
        "DataName": "DataName",    //常规字段
        "DeviceName": "DeviceName", //支持嵌套属性
        "DataValue": "DataValue",
        "DataUnit": "DataUnit",
        "RecordTime": {
          field: "RecordTime",
          //自定义回调函数
          callback: value => {
            let date = new Date(value)
            return formatDate(date,'yyyy-MM-dd hh:mm:ss')
          }
        }
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
      SelectAlarmData:7,
      deviceTreeData:[],
      AlarmDataTree:[],
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
      form: this.$form.createForm(this),
      messageShowLoad:false,
      advanced: true,
      refIconLoading: false,
      templeteUuid:"",
      columns: [
        {
          slotName: 'reporting.DataHistory.DeviceName',
          width: '10%',
          scopedSlots: { customRender: 'serial', title: 'reporting.DataHistory.DeviceName' },
          dataIndex: 'DeviceName',
        },
        {
          width: '10%',
          slotName: 'reporting.DataHistory.DataName',
          scopedSlots: { customRender: 'serial', title: 'reporting.DataHistory.DataName' },
          dataIndex: 'DataName',
        },
        {
          width: '10%',
          slotName: 'reporting.DataHistory.DataValue',
          scopedSlots: { customRender: 'DataValue', title: 'reporting.DataHistory.DataValue' },
          dataIndex: 'DataValue',
        },
        {
          width: '10%',
          slotName: 'reporting.DataHistory.DataUnit',
          scopedSlots: { customRender: 'DataUnit', title: 'reporting.DataHistory.DataUnit' },
          dataIndex: 'DataUnit',
        },
        {
          width: '10%',
          slotName: 'reporting.DataHistory.RecordTime',
          scopedSlots: { customRender: 'RecordTime', title: 'reporting.DataHistory.RecordTime' },
          dataIndex: 'RecordTime',
        }
      ],
      dataSource: [],
      conditionExpress:"",
      selectedRows: [],
      sheetOptions:{},
    }
  },
  authorize: {
    // deleteRecord: 'delete'
  },
  mounted(){
    this.sheetOptions = {
      container: 'luckysheetQuery', // 设定DOM容器的id
      lang: 'zh', // 设定表格语言
      plugins:['chart'],
      allowEdit: false,
      showinfobar:false,
      showtoolbar:false
    }
    if(this.$i18n.locale=="CN")
    {
      this.sheetOptions.lang='zh'
      this.json_fields = this.json_fields_cn
    }
    else
    {
      this.sheetOptions.lang='en'
      this.json_fields = this.json_fields_en
    }
    this.exportName = this.$t('reporting.DataHistory.exportName')+"."+formatDate( new Date(),'yyyy-MM-dd hh:mm:ss')+".xlsx"

    let _t = this
    luckysheet.destroy()
    // In some cases, you need to use $nextTick
    this.$nextTick(() => {
      this.luckysheetOpt =  luckysheet.create(this.sheetOptions);
    });
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
    let _t = this
    luckysheet.destroy()
    // In some cases, you need to use $nextTick
    this.$nextTick(() => {
      this.luckysheetOpt =  luckysheet.create(this.sheetOptions);
    });
  },
  watch: {
    '$route'() {
      luckysheet.destroy()
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
    loadTemplete(){
      this.$refs.reportTempleteModel.showModal()
    },
    onSelectTemplete(selectData) {
      this.TempleteName = this.$t('diyReportTemplete.LoadTempleteName')+selectData.Name
      this.templeteUuid = selectData.Uuid
      this.LoadReportTemplete(selectData.Uuid)
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
    exportExcelBtn() {
      // console.log(luckysheet.getluckysheetfile())
      exportExcel(luckysheet.getluckysheetfile(), this.exportName)
    },
    LoadReportTemplete(uuid){
      let myDate = new Date()
      let  value = '/static/reportTemplete/'+uuid+'.xlsx?'+myDate.getMilliseconds()
      let name = uuid
      let _t = this
      LuckyExcel.transformExcelToLuckyByUrl(value, name, function (exportJson, luckysheetfile) {
        if (exportJson.sheets == null || exportJson.sheets.length == 0) {
          return
        }
        luckysheet.destroy()
        _t.sheetOptions.data = exportJson.sheets
        _t.sheetOptions.title = exportJson.info.name,
            _t.sheetOptions.userInfo = exportJson.info.name.creator,
            luckysheet.create(_t.sheetOptions)
      })
    },
    LoadRemoteReport(path){
      let myDate = new Date()
      let  value = '/'+path+'?'+myDate.getMilliseconds()
      // let  value = '/static/HistoryData/dd.xlsx'
      let name = formatDate(new Date(),'yyyy-MM-dd hh:mm:ss')
      let _t = this
      LuckyExcel.transformExcelToLuckyByUrl(value, name, function (exportJson, luckysheetfile) {
          if (exportJson.sheets == null || exportJson.sheets.length == 0) {
            return
          }
          luckysheet.destroy()
          _t.sheetOptions.data = exportJson.sheets
          _t.sheetOptions.title = exportJson.info.name
          _t.sheetOptions.userInfo = exportJson.info.name.creator
          luckysheet.create(_t.sheetOptions)
      })
    },
    QueryHistoryDataList(){
      let _t = this
      if(_t.templeteUuid=="")
      {
        _t.$message.error(_t.$t("diyReportTemplete.TempleteEmpty"))
        return
      }
      getExcelData(luckysheet.getluckysheetfile(),function (data){
        _t.dataSource = []
        const params = {
          deviceList:_t.SelectDevice,
          dataList:_t.SelectAlarmData,
          dateRange:_t.SelectDateRange,
          dateType:_t.SelectDateType,
          timeIn:_t.SelectAlarmData,
          tUuid:_t.templeteUuid,
        }
        if(params.dateRange==""||(params.dateRange[0]==""))
        {
          _t.$message.error(_t.$t("reporting.DataHistory.SelectDateError"))
          return
        }
        if(params.deviceList=="")
        {
          _t.$message.error(_t.$t("diyReport.SelectDeviceTips"))
          return
        }
        if(params.timeIn==0)
        {
          _t.$message.error(_t.$t("diyReport.SelectInTips"))
          return
        }
        _t.$message.loading({ content: _t.$t('diyReport.QueryLoading'),duration: 0 });
        _t.messageShowLoad=true
        GetDiyDataHistoryList(params).then(function (res){
          if(res.data.code==0)
          {
            _t.$message.destroy()
            _t.LoadRemoteReport(res.data.path)
            const elink = document.createElement('a')
            elink.href = res.data.path
            elink.setAttribute('download', _t.exportName)
            elink.style.display = 'none'
            document.body.appendChild(elink)
            setTimeout(() => {
              elink.click()
              document.body.removeChild(elink)
            }, 66)
          }
          _t.messageShowLoad=false
        }).catch(function(){
          _t.messageShowLoad=false
          _t.$message.destroy()
          _t.$message.error(_t.$t('loginPage.serverError'), 3)
        })
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
