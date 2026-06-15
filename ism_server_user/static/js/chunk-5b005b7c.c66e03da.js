(window.webpackJsonp=window.webpackJsonp||[]).push([["chunk-5b005b7c"],{
/***/"774d":
/***/function(e,t,i){"use strict";
/* harmony import */i("7a6d");
/* harmony import */},
/***/"7a6d":
/***/function(e,t,i){
// extracted by mini-css-extract-plugin
/***/},
/***/af4d:
/***/function(e,t,i){"use strict";
// ESM COMPAT FLAG
i.r(t);
// CONCATENATED MODULE: ./node_modules/cache-loader/dist/cjs.js?{"cacheDirectory":"node_modules/.cache/vue-loader","cacheIdentifier":"50a92b3c-vue-loader-template"}!./node_modules/cache-loader/dist/cjs.js??ref--13-0!./node_modules/thread-loader/dist/cjs.js!./node_modules/babel-loader/lib!./node_modules/vue-loader/lib/loaders/templateLoader.js??ref--6!./node_modules/cache-loader/dist/cjs.js??ref--1-0!./node_modules/vue-loader/lib??vue-loader-options!./src/pages/reporting/diyReport/diyReportContent.vue?vue&type=template&id=1ab90901&scoped=true
i("14d9"),i("fb6a"),i("b0c0");var o=i("600d"),s=(i("c2ec"),i("c1df"),i("f7e6")),a=i.n(s),n=i("fded"),l=i("cf45"),c=(i("f2d9"),i("0340")),r=i("f64d"),h={name:"diyDataHistory",i18n:i("89fe"),components:{
// DownloadExcel,
DeviceHistoryDataModel:n.a},data:function(){return{isMaskShow:!1,selectRow:0,selectCol:0,exportName:"",messageShowLoad:!1,advanced:!0,sheetOptions:{}}},authorize:{},mounted:function(){},activated:function(){var e=this,t=this;this.sheetOptions={container:"luckysheetContent",
// 设定DOM容器的id
title:"luckysheetContent",
// 设定表格名称
lang:"zh",
// 设定表格语言
gridKey:t.$route.params.uuid,plugins:["chart"],loading:{image:function(){return'<svg viewBox="25 25 50 50" class="circular">\n            <circle cx="50" cy="50" r="20" fill="none"></circle>\n            </svg>'},imageClass:"loadingAnimation"},showinfobar:!1,cellRightClickConfig:{customs:[{title:this.$t("diyReport.SelectData"),onClick:function(e,i,o){t.InsertData(e,i,o)}},{title:this.$t("diyReport.SelectDataMax"),onClick:function(e,i,o){t.InsertVar(e,i,o,4)}},{title:this.$t("diyReport.SelectDataMin"),onClick:function(e,i,o){t.InsertVar(e,i,o,5)}},{title:this.$t("diyReport.SelectDataDifference"),onClick:function(e,i,o){t.InsertVar(e,i,o,6)}},{title:this.$t("diyReport.SelectDataSum"),onClick:function(e,i,o){t.InsertVar(e,i,o,7)}},{title:this.$t("diyReport.SelectDataCount"),onClick:function(e,i,o){t.InsertVar(e,i,o,8)}},{title:this.$t("diyReport.SelectDataAverage"),onClick:function(e,i,o){t.InsertVar(e,i,o,9)}},{title:this.$t("diyReport.SelectDate"),onClick:function(e,i,o){t.InsertVar(e,i,o,1)}},{title:this.$t("diyReport.SelectDevice"),onClick:function(e,i,o){t.InsertVar(e,i,o,2)}},{title:this.$t("diyReport.SelectRange"),onClick:function(e,i,o){t.InsertVar(e,i,o,3)}}],copy:!0,
// 复制
copyAs:!1,
// 复制为
paste:!0,
// 粘贴
insertRow:!1,
// 插入行
insertColumn:!1,
// 插入列
deleteRow:!1,
// 删除选中行
deleteColumn:!1,
// 删除选中列
deleteCell:!1,
// 删除单元格
hideRow:!1,
// 隐藏选中行和显示选中行
hideColumn:!1,
// 隐藏选中列和显示选中列
rowHeight:!1,
// 行高
columnWidth:!1,
// 列宽
clear:!1,
// 清除内容
matrix:!1,
// 矩阵操作选区
sort:!1,
// 排序选区
filter:!1,
// 筛选选区
chart:!1,
// 图表生成
image:!0,
// 插入图片
link:!1,
// 插入链接
data:!1,
// 数据验证
cellFormat:!1},showtoolbarConfig:{undoRedo:!0,
//撤销重做，注意撤消重做是两个按钮，由这一个配置决定显示还是隐藏
paintFormat:!0,
//格式刷
currencyFormat:!1,
//货币格式
percentageFormat:!1,
//百分比格式
numberDecrease:!1,
// '减少小数位数'
numberIncrease:!1,
// '增加小数位数
moreFormats:!1,
// '更多格式'
font:!0,
// '字体'
fontSize:!0,
// '字号大小'
bold:!0,
// '粗体 (Ctrl+B)'
italic:!0,
// '斜体 (Ctrl+I)'
strikethrough:!0,
// '删除线 (Alt+Shift+5)'
underline:!0,
// '下划线 (Alt+Shift+6)'
textColor:!0,
// '文本颜色'
fillColor:!0,
// '单元格颜色'
border:!0,
// '边框'
mergeCell:!0,
// '合并单元格'
horizontalAlignMode:!0,
// '水平对齐方式'
verticalAlignMode:!0,
// '垂直对齐方式'
textWrapMode:!0,
// '换行方式'
textRotateMode:!0,
// '文本旋转方式'
image:!0,
// '插入图片'
link:!0,
// '插入链接'
chart:!1,
// '图表'（图标隐藏，但是如果配置了chart插件，右击仍然可以新建图表）
postil:!1,
//'批注'
pivotTable:!1,
//'数据透视表'
function:!1,
// '公式'
frozenMode:!1,
// '冻结方式'
sortAndFilter:!0,
// '排序和筛选'
conditionalFormat:!0,
// '条件格式'
dataVerification:!1,
// '数据验证'
splitColumn:!0,
// '分列'
screenshot:!1,
// '截图'
findAndReplace:!0,
// '查找替换'
protection:!1,
// '工作表保护'
print:!1,
// '打印'
selectData:!0}},"CN"==this.$i18n.locale?(this.sheetOptions.lang="zh",this.json_fields=this.json_fields_cn):(this.sheetOptions.lang="en",this.json_fields=this.json_fields_en),this.exportName=this.$t("reporting.DataHistory.exportName")+"."+Object(l.a)(new Date,"yyyy-MM-dd hh:mm:ss")+".xlsx",
// In some cases, you need to use $nextTick
this.$nextTick((function(){
// luckysheet.create( this.sheetOptions);
var t=new Date,i="/static/reportTemplete/"+e.$route.params.uuid+".xlsx?"+t.getMilliseconds(),o=e.$route.params.uuid,s=e;a.a.transformExcelToLuckyByUrl(i,o,(function(e,t){null!=e.sheets&&0!=e.sheets.length&&(s.sheetOptions.data={},
// luckysheet.destroy()
s.sheetOptions.data=e.sheets,s.sheetOptions.title=e.info.name,s.sheetOptions.userInfo=e.info.name.creator,luckysheet.create(s.sheetOptions))}))}))},filters:{formatDate:function(e){var t=new Date(e);return Object(l.a)(t,"yyyy-MM-dd hh:mm:ss")}},created:function(){},watch:{$route:function(){luckysheet.destroy()}},methods:{onSelectData:function(e){luckysheet.setCellValue(this.selectRow,this.selectCol,"{{DataModel."+e.name+"}}")},InsertVar:function(e,t,i,o){this.selectRow=i.rowIndex,this.selectCol=i.columnIndex;var s=luckysheet.getCellValue(this.selectRow,this.selectCol);if(""!=s&&null!=s)s=s.slice(0,s.length-2);else if(1!=o&&2!=o&&3!=o)return;1==o?luckysheet.setCellValue(this.selectRow,this.selectCol,"{{DataModel.HistoryRecordDateTime}}"):2==o?luckysheet.setCellValue(this.selectRow,this.selectCol,"{{DeviceName}}"):3==o?luckysheet.setCellValue(this.selectRow,this.selectCol,"{{TimeRange}}"):4==o?luckysheet.setCellValue(this.selectRow,this.selectCol,s+"的时间段最大值}}"):5==o?luckysheet.setCellValue(this.selectRow,this.selectCol,s+"的时间段最小值}}"):6==o?luckysheet.setCellValue(this.selectRow,this.selectCol,s+"的时间段差值}}"):7==o?luckysheet.setCellValue(this.selectRow,this.selectCol,s+"的时间段和}}"):8==o?luckysheet.setCellValue(this.selectRow,this.selectCol,s+"的时间段数量}}"):9==o&&luckysheet.setCellValue(this.selectRow,this.selectCol,s+"的时间段平均值}}")},InsertData:function(e,t,i){this.$refs.deviceHistoryDataModel.showDataModal(),this.selectRow=i.rowIndex,this.selectCol=i.columnIndex},startDownload:function(){this.isLoadExecl=!0,this.loadExecl=this.$message.loading(this.$t("reporting.DataHistory.LoadingExecl"),0)},finishDownload:function(){this.$message.destroy(this.loadExecl),this.isLoadExecl=!1},filterOption:function(e,t){return t.componentOptions.children[0].text.toLowerCase().indexOf(e.toLowerCase())>=0},SelectTreeDevice:function(e,t,i){this.GetDeviceModelDataList()},GetDeviceModelDataList:function(){var e=this;this.AlarmDataTree=[];var t={SelectDevice:this.SelectDevice,getType:2};Object(o.b)(t).then((function(t){0==t.data.code&&(e.AlarmDataTree=t.data.list)}))},onDateChange:function(e,t){this.SelectDateRange=t},LoadReportTemplete:function(){var e=new Date,t="/static/reportTemplete/"+this.$route.params.uuid+".xlsx?"+e.getMilliseconds(),i=this;luckysheet.destroy(),a.a.transformExcelToLuckyByUrl(t,"1111",(function(e,t){null!=e.sheets&&0!=e.sheets.length&&(i.sheetOptions.data=e.sheets,i.sheetOptions.title=e.info.name,i.sheetOptions.userInfo=e.info.name.creator,luckysheet.create(i.sheetOptions))}))},back:function(){this.$router.push("/Reporting/DiyReportTemplete")},SaveReportTemplete:function(){var e=this;Object(c.b)(luckysheet.getluckysheetfile(),(function(t){e.dataSource=[];var i={Uuid:e.$route.params.uuid,sheetData:t};e.messageShowLoad=!0,Object(r.f)(i).then((function(t){0==t.data.code&&e.$message.success(e.$t("diyReportTemplete.SaveSuccess"),3),e.messageShowLoad=!1})).catch((function(){e.messageShowLoad=!1,e.$message.destroy(),e.$message.error(e.$t("loginPage.serverError"),3)}))}))}},destroyed:function(){luckysheet.destroy()}},d=(i("774d"),i("2877")),u=Object(d.a)(h,(function(){var e=this,t=e._self._c;return t("a-card",[t("a-row",[t("a-col",{attrs:{span:8}},[t("a-button",{staticStyle:{"margin-left":"10px","margin-top":"10px","margin-bottom":"10px"},attrs:{icon:"save",type:"primary"},on:{click:e.SaveReportTemplete}},[e._v(e._s(e.$t("diyReportTemplete.Save")))]),t("a-button",{staticStyle:{"margin-left":"10px","margin-top":"10px","margin-bottom":"10px"},attrs:{icon:"reload",type:"default"},on:{click:e.LoadReportTemplete}},[e._v(e._s(e.$t("diyReportTemplete.Reload")))]),t("a-button",{staticStyle:{"margin-left":"10px","margin-top":"10px","margin-bottom":"10px"},attrs:{icon:"backward",type:"default"},on:{click:e.back}},[e._v(e._s(e.$t("dataModel.modbusModel.Back")))])],1),t("a-col",{attrs:{span:12}},[t("div",{staticStyle:{"margin-left":"10px","margin-top":"10px","margin-bottom":"10px","margin-right":"10px"}},[t("a-alert",{attrs:{message:e.$t("diyReportTemplete.TempleteTips"),type:"info","show-icon":""}})],1)])],1),t("div",{staticClass:"hello"},[t("div",{staticStyle:{margin:"0px",padding:"0px",height:"600px",top:"5px","z-index":"-1"},attrs:{id:"luckysheetContent"}}),t("div",{directives:[{name:"show",rawName:"v-show",value:e.isMaskShow,expression:"isMaskShow"}],staticStyle:{position:"absolute","z-index":"1000000",left:"0px",top:"0px",bottom:"0px",right:"0px",background:"rgba(255, 255, 255, 0.8)","text-align":"center","font-size":"40px","align-items":"center","justify-content":"center",display:"flex"}},[e._v("Downloading")]),t("device-history-data-model",{ref:"deviceHistoryDataModel",on:{onSelectDataModel:e.onSelectData}})],1)],1)}),[],!1,null,"1ab90901",null)
/* harmony default export */;t.default=u.exports}}]);