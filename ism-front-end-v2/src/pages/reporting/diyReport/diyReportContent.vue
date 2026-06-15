<template>
  <a-card>
    <a-row>
      <a-col :span="8">
        <a-button icon="save"  style="margin-left:10px;margin-top:10px;margin-bottom: 10px" type="primary" @click="SaveReportTemplete">{{$t('diyReportTemplete.Save')}}</a-button>
        <a-button icon="reload"  style="margin-left:10px;margin-top:10px;margin-bottom: 10px" type="default" @click="LoadReportTemplete">{{$t('diyReportTemplete.Reload')}}</a-button>
        <a-button icon="backward"  style="margin-left:10px;margin-top:10px;margin-bottom: 10px" type="default" @click="back">{{$t('dataModel.modbusModel.Back')}}</a-button>
      </a-col>
      <a-col :span="12">
        <div style="margin-left:10px;margin-top:10px;margin-bottom: 10px;margin-right: 10px">
           <a-alert :message="$t('diyReportTemplete.TempleteTips')" type="info" show-icon />
        </div>
      </a-col>
    </a-row>
    <div class="hello">
      <div
          id="luckysheetContent"
          style="margin:0px;padding:0px;height:600px;top:5px;z-index: -1"
      ></div>

      <div v-show="isMaskShow" style="position: absolute;z-index: 1000000;left: 0px;top: 0px;bottom: 0px;right: 0px; background: rgba(255, 255, 255, 0.8); text-align: center;font-size: 40px;align-items:center;justify-content: center;display:flex;">Downloading</div>
      <device-history-data-model @onSelectDataModel="onSelectData" ref="deviceHistoryDataModel"></device-history-data-model>
    </div>
  </a-card>
</template>

<script>
import {getMonitorTree} from "@/services/device";
import {GetDeviceModelDataList} from "../../../services/device";
import {GetDataHistoryList, GetDiyDataHistoryList} from "../../../services/report";
import moment from 'moment';
import LuckyExcel from 'luckyexcel'
import DeviceHistoryDataModel from "@/components/deviceHistoryDataModel/deviceHistoryDataModel";
import {formatDate} from '@/utils/common';
import { exportExcel,getExcelData } from '@/utils/export'
import {SaveReportTemplete} from "@/services/reportTemplete";

export default {
  name: 'diyDataHistory',
  i18n: require('../../../i18n/language'),
  components:{
    DeviceHistoryDataModel
  },
  data () {
    return {
      isMaskShow: false,
      selectRow:0,
      selectCol:0,
      exportName:"",
      messageShowLoad:false,
      advanced: true,
      sheetOptions:{}
    }
  },
  authorize: {
    // deleteRecord: 'delete'
  },
  mounted(){

  },
  activated(){
    let _t = this
    this.sheetOptions = {
      container: 'luckysheetContent', // 设定DOM容器的id
      title: 'luckysheetContent', // 设定表格名称
      lang: 'zh', // 设定表格语言
      gridKey:_t.$route.params.uuid,
      plugins:['chart'],
      loading:{
        image:()=>{
          return `<svg viewBox="25 25 50 50" class="circular">
            <circle cx="50" cy="50" r="20" fill="none"></circle>
            </svg>`
        },
        imageClass:"loadingAnimation"
      },
      showinfobar:false,
      cellRightClickConfig: {
        customs: [
          {
            title: this.$t('diyReport.SelectData'),
            onClick: (clickEvent, event, params) => {
              _t.InsertData(clickEvent, event, params)
            },
          },
          {
            title: this.$t('diyReport.SelectDataMax'),
            onClick: (clickEvent, event, params) => {
              _t.InsertVar(clickEvent, event, params,4)
            },
          },
          {
            title: this.$t('diyReport.SelectDataMin'),
            onClick: (clickEvent, event, params) => {
              _t.InsertVar(clickEvent, event, params,5)
            },
          },
          {
            title: this.$t('diyReport.SelectDataDifference'),
            onClick: (clickEvent, event, params) => {
              _t.InsertVar(clickEvent, event, params,6)
            },
          },
          {
            title: this.$t('diyReport.SelectDataSum'),
            onClick: (clickEvent, event, params) => {
              _t.InsertVar(clickEvent, event, params,7)
            },
          },
          {
            title: this.$t('diyReport.SelectDataCount'),
            onClick: (clickEvent, event, params) => {
              _t.InsertVar(clickEvent, event, params,8)
            },
          },
          {
            title: this.$t('diyReport.SelectDataAverage'),
            onClick: (clickEvent, event, params) => {
              _t.InsertVar(clickEvent, event, params,9)
            },
          },

          {
            title: this.$t('diyReport.SelectDate'),
            onClick: (clickEvent, event, params) => {
              _t.InsertVar(clickEvent, event, params,1)
            },
          },
          {
            title: this.$t('diyReport.SelectDevice'),
            onClick: (clickEvent, event, params) => {
              _t.InsertVar(clickEvent, event, params,2)
            },
          },
          {
            title: this.$t('diyReport.SelectRange'),
            onClick: (clickEvent, event, params) => {
              _t.InsertVar(clickEvent, event, params,3)
            },
          },
        ],
        copy: true, // 复制
        copyAs: false, // 复制为
        paste: true, // 粘贴
        insertRow: false, // 插入行
        insertColumn: false, // 插入列
        deleteRow: false, // 删除选中行
        deleteColumn: false, // 删除选中列
        deleteCell: false, // 删除单元格
        hideRow: false, // 隐藏选中行和显示选中行
        hideColumn: false, // 隐藏选中列和显示选中列
        rowHeight: false, // 行高
        columnWidth: false, // 列宽
        clear: false, // 清除内容
        matrix: false, // 矩阵操作选区
        sort: false, // 排序选区
        filter: false, // 筛选选区
        chart: false, // 图表生成
        image: true, // 插入图片
        link: false, // 插入链接
        data: false, // 数据验证
        cellFormat: false // 设置单元格格式
      },
      showtoolbarConfig: {
        undoRedo: true, //撤销重做，注意撤消重做是两个按钮，由这一个配置决定显示还是隐藏
        paintFormat: true, //格式刷
        currencyFormat: false, //货币格式
        percentageFormat: false, //百分比格式
        numberDecrease: false, // '减少小数位数'
        numberIncrease: false, // '增加小数位数
        moreFormats: false, // '更多格式'
        font: true, // '字体'
        fontSize: true, // '字号大小'
        bold: true, // '粗体 (Ctrl+B)'
        italic: true, // '斜体 (Ctrl+I)'
        strikethrough: true, // '删除线 (Alt+Shift+5)'
        underline: true, // '下划线 (Alt+Shift+6)'
        textColor: true, // '文本颜色'
        fillColor: true, // '单元格颜色'
        border: true, // '边框'
        mergeCell: true, // '合并单元格'
        horizontalAlignMode: true, // '水平对齐方式'
        verticalAlignMode: true, // '垂直对齐方式'
        textWrapMode: true, // '换行方式'
        textRotateMode: true, // '文本旋转方式'
        image:true, // '插入图片'
        link:true, // '插入链接'
        chart: false, // '图表'（图标隐藏，但是如果配置了chart插件，右击仍然可以新建图表）
        postil:  false, //'批注'
        pivotTable: false,  //'数据透视表'
        function: false, // '公式'
        frozenMode: false, // '冻结方式'
        sortAndFilter: true, // '排序和筛选'
        conditionalFormat: true, // '条件格式'
        dataVerification: false, // '数据验证'
        splitColumn: true, // '分列'
        screenshot: false, // '截图'
        findAndReplace: true, // '查找替换'
        protection:false, // '工作表保护'
        print:false, // '打印'
        selectData:true
      },
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
    // In some cases, you need to use $nextTick
    this.$nextTick(() => {
      // luckysheet.create( this.sheetOptions);
      let myDate = new Date()
      let  value = '/static/reportTemplete/'+this.$route.params.uuid+'.xlsx?'+myDate.getMilliseconds()
      let name = this.$route.params.uuid
      let _t = this
      LuckyExcel.transformExcelToLuckyByUrl(value, name, function (exportJson, luckysheetfile) {
        if (exportJson.sheets == null || exportJson.sheets.length == 0) {
          return
        }

        _t.sheetOptions.data={}
        // luckysheet.destroy()
        _t.sheetOptions.data = exportJson.sheets
        _t.sheetOptions.title = exportJson.info.name,
          _t.sheetOptions.userInfo = exportJson.info.name.creator,
          luckysheet.create(_t.sheetOptions)
      })
    });
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
    '$route'() {
      luckysheet.destroy()
    }
  },
  methods: {
    onSelectData(selectData) {
      luckysheet.setCellValue(this.selectRow ,this.selectCol,"{{DataModel."+selectData.name+"}}")
    },
    InsertVar(clickEvent, event, params,type) {
      this.selectRow = params.rowIndex
      this.selectCol = params.columnIndex
      let ColValue = luckysheet.getCellValue(this.selectRow ,this.selectCol)
      if((ColValue!="")&&(ColValue!=null))
      {
        ColValue = ColValue.slice(0,ColValue.length-2);
      }
      else if((type!=1)&&(type!=2)&&(type!=3))
      {
        return
      }
      if(type==1)
      {
        luckysheet.setCellValue(this.selectRow ,this.selectCol,"{{DataModel.HistoryRecordDateTime}}")
      }
      else if(type==2)
      {
        luckysheet.setCellValue(this.selectRow ,this.selectCol,"{{DeviceName}}")
      }
      else if(type==3)
      {
        luckysheet.setCellValue(this.selectRow ,this.selectCol,"{{TimeRange}}")
      }
      else if(type==4)
      {
        luckysheet.setCellValue(this.selectRow ,this.selectCol,ColValue+"的时间段最大值}}")
      }
      else if(type==5)
      {
        luckysheet.setCellValue(this.selectRow ,this.selectCol,ColValue+"的时间段最小值}}")
      }
      else if(type==6)
      {
        luckysheet.setCellValue(this.selectRow ,this.selectCol,ColValue+"的时间段差值}}")
      }
      else if(type==7)
      {
        luckysheet.setCellValue(this.selectRow ,this.selectCol,ColValue+"的时间段和}}")
      }
      else if(type==8)
      {
        luckysheet.setCellValue(this.selectRow ,this.selectCol,ColValue+"的时间段数量}}")
      }
      else if(type==9)
      {
        luckysheet.setCellValue(this.selectRow ,this.selectCol,ColValue+"的时间段平均值}}")
      }
    },
    InsertData(clickEvent, event, params){
      this.$refs.deviceHistoryDataModel.showDataModal()
      this.selectRow = params.rowIndex
      this.selectCol = params.columnIndex
    },
    startDownload(){
      this.isLoadExecl = true
      this.loadExecl = this.$message.loading(this.$t("reporting.DataHistory.LoadingExecl"),0)
    },
    finishDownload(){
      this.$message.destroy(this.loadExecl)
      this.isLoadExecl = false
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
      const params ={
        SelectDevice:this.SelectDevice,
        getType:2
      }
      GetDeviceModelDataList(params).then(function (res){
        if(res.data.code==0)
        {
          _t.AlarmDataTree =res.data.list
        }
      })
    },
    onDateChange(date, dateString){
      this.SelectDateRange = dateString
    },
    LoadReportTemplete(){
      let myDate = new Date()
      let  value = '/static/reportTemplete/'+this.$route.params.uuid+'.xlsx?'+myDate.getMilliseconds()
      let name = '1111'
      let _t = this
      luckysheet.destroy()
      LuckyExcel.transformExcelToLuckyByUrl(value, name, function (exportJson, luckysheetfile) {
        if (exportJson.sheets == null || exportJson.sheets.length == 0) {
          return
        }

        _t.sheetOptions.data = exportJson.sheets
        _t.sheetOptions.title = exportJson.info.name,
        _t.sheetOptions.userInfo = exportJson.info.name.creator,
        luckysheet.create(_t.sheetOptions)
      })
    },
    back(){
      this.$router.push('/Reporting/DiyReportTemplete')
    },
    SaveReportTemplete(){
      let _t = this
      getExcelData(luckysheet.getluckysheetfile(),function (data){
        _t.dataSource = []
        const params = {
          Uuid:_t.$route.params.uuid,
          sheetData:data
        }
        _t.messageShowLoad=true
        SaveReportTemplete(params).then(function (res){
          if(res.data.code==0)
          {
            _t.$message.success(_t.$t('diyReportTemplete.SaveSuccess'), 3)
          }
          _t.messageShowLoad=false
        }).catch(function(){
          _t.messageShowLoad=false
          _t.$message.destroy()
          _t.$message.error(_t.$t('loginPage.serverError'), 3)
        })
      })
    }
  },
  destroyed () {
    luckysheet.destroy()
  },
}
</script>

<style lang="less" scoped>

::v-deep .ant-card-body {
  padding: 0px;
  zoom: 1;
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