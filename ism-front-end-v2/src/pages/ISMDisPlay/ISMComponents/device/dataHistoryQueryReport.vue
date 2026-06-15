<template>
  <div :style="{'background':detail.style.backColor}">
    <div :style="styleVar">
    <DataHistoryQueryReport :tablePagination="pagination"/>
    </div>
  </div>
</template>

<script>
import DataHistoryQueryReport from "@/pages/reporting/dataHistoryReport/dataHistoryQueryReport";
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
  name: 'DataHistoryReportComponents',
  inject: ['getNode'],
  i18n: require('@/i18n/language'),
  components:{
    DataHistoryQueryReport
  },
  computed: {
    styleVar() {
      return {
        "height": this.detail.style.position.h+'px',
        "--foreColor": this.foreColor ,
        '--backColor':this.backColor,
        "--selectedColor": this.selectedColor ,
        '--hoverColor': this.hoverColor,
        '--selectedTextColor': this.selectedTextColor,
        '--TextFontSize': this.TextFontSize+'px',
        '--hoverTextColor': this.hoverTextColor,
        '--SearchColor': this.SearchColor,
        '--SearchBackColor': this.SearchBackColor,
        '--SearchBorderColor': this.SearchBorderColor,

        '--dateSelectColor': this.dateSelectColor,
        '--dateSelectBackColor': this.dateSelectBackColor,
        '--dateSelectBorderColor': this.dateSelectBorderColor,

        '--tableHeaderColor': this.tableHeaderColor,
        '--tableHeaderBackColor': this.tableHeaderBackColor,
        '--tableSplitColor': this.tableSplitColor,
        '--tableHoverColor':this.tableHoverColor
      };
    },
    CardStyle:function () {
      let styles = [];
      if(this.detail.style.backColor) {
        styles.push(`background-color: ${this.detail.style.backColor}`);
      }
      if(this.detail.style.foreColor) {
        styles.push(`color: ${this.detail.style.foreColor}`);
      }
      let style = styles.join(';');
      return style;
    },
  },
  data () {
    return {
      detail:null,
      IsToolBox:false,
      editMode:true,
      dataZ:[],
      valueData: '',
      treePageSize: 100,
      scrollPage: 1,
      frontDataZ:[],
      pagination:{
        pageSize:10,
        showSizeChanger:false,
        hideOnSinglePage:true,
        showLessItems:true,
        simple:true
      },
      foreColor:"#000000",
      backColor:"#ffffff",
      tableSplitColor:"#000",
      tableHoverColor:"#fff",
      tableHeaderColor:"",
      tableHeaderBackColor:"",
      SearchColor:"#000000",
      SearchBorderColor:"#108ec4",
      SearchBackColor:"#ffffff",
      dateSelectColor:"",
      dateSelectBackColor:"",
      dateSelectBorderColor:"",
      TextFontSize:14,
      loadExecl:null,
      isLoadExecl:false,
      exportName:"",

      base:{
        text: "configComponent.DataHistoryListReport.title",
        "icon": "icon-a-5_lishibaobiao",
        "isFontIcon": true,
        "info": {
          "type": "image",
          "action": [],
          "dataBind":[],
          "animate": {
            "selected": [],
            "condition":{
              deviceSN:"",
              selectVideoType:0,
              isBandDevice:false,
              bandType:1,
              dataID: "",
              dataName: "",
              operator:"",
              OperatorValue:"",
              OperatorMaxValue:"",
            },
            "isExpression": false,
            "animateList": [
              {
                id: "Forbidden",
                name: "component.public.Forbidden",
              },
              {
                id: "blink",
                name: "component.public.animateBlink",
              },
              {
                id: "Zoom",
                name: "component.public.Zoom",
              },
              {
                id: "animateSpin",
                name: "component.public.animateSpin",
              },
            ],
            "animateElement": [
              {
                id: "blink",
                elementList:[
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"blinkSpeed",
                  },
                ]
              },
              {
                id: "millcolorGrad",
                elementList:[
                  {
                    "name": "component.public.startColor",
                    "type": 2,
                    "value": "#74f808",
                    "key": "startColor",
                  },
                  {
                    "name": "component.public.stopColor",
                    "type": 2,
                    "value": "#f30b0b",
                    "key": "stopColor",
                  },
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"animateSpeed",
                  },
                ]
              },
              {
                id: "animateSpin",
                elementList:[
                  {
                    "name":"component.public.animateSpinSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"spinSpeed",
                  },
                  {
                    name:"configComponent.bigScreen.border.border89Direction",
                    type:6,
                    value:0,
                    enumList:[
                      {
                        value:0,
                        option:"configComponent.bigScreen.border.border89DirectionForward"
                      },
                      {
                        value:1,
                        option:"configComponent.bigScreen.border.border89DirectionNegative"
                      }
                    ],
                    min:1,
                    key:"spinDirection",
                  }
                ]
              },
            ],
          },
          "style": {
            "position": {
              "x": 0,
              "y": 0,
              "w": 600,
              "h": 400
            },
            "visible":1,
            "backColor": "#ffffff",
            "foreColor": "#000000",
            "zIndex": -1,
            "transform": 0,
            "diy":[
              {
                "name":"configComponent.DeviceTree.ShowCount",
                "type":1,
                "value":5,
                "min":1,
                "max":100,
                "key":"ShowCount",
              },
              {
                "name":"configComponent.DeviceTree.SearchColor",
                "type":2,
                "value":"#000000",
                "key":"SearchColor",
              },
              {
                "name":"configComponent.DeviceTree.SearchBackColor",
                "type":2,
                "value":"#ffffff",
                "key":"SearchBackColor",
              },
              {
                "name":"configComponent.DeviceTree.SearchBorderColor",
                "type":2,
                "value":"#cbc6c6",
                "key":"SearchBorderColor",
              },
              {
                "name":"configComponent.DataHistoryList.dateSelectColor",
                "type":2,
                "value":"#36cfc9",
                "key":"dateSelectColor",
              },
              {
                "name":"configComponent.DataHistoryList.dateSelectBackColor",
                "type":2,
                "value":"#ffffff",
                "key":"dateSelectBackColor",
              },
              {
                "name":"configComponent.DataHistoryList.dateSelectBorderColor",
                "type":2,
                "value":"#ebedf0",
                "key":"dateSelectBorderColor",
              },

              {
                "name":"configComponent.DataHistoryList.tableHeaderColor",
                "type":2,
                "value":"#000000",
                "key":"tableHeaderColor",
              },
              {
                "name":"configComponent.DataHistoryList.tableHeaderBackColor",
                "type":2,
                "value":"#fafafa",
                "key":"tableHeaderBackColor",
              },
              {
                "name":"configComponent.DataHistoryList.tableSplitColor",
                "type":2,
                "value":"#ebedf0",
                "key":"tableSplitColor",
              },
              {
                "name":"configComponent.DataHistoryList.tableHoverColor",
                "type":2,
                "value":"#ffffff",
                "key":"tableHoverColor",
              },
            ]
          }
        }
      }
    }
  },
  authorize: {

  },
  mounted(){

  },
  activated(){

  },
  created(){


    this.$nextTick(function(){
      this.initComponents(this.detail);
    });
    let _t = this
    this.GetNodeObj = this.getNode()
    this.GetNodeObj.on('change:data', ({ current }) => {
      if(current) {
        _t.detail = current.detail
      }
    })
    this.GetNodeObj.on('change:size', ({ current }) => {
      _t.detail.style.position.w = current.width
      _t.detail.style.position.h = current.height
    });
    this.detail = this.GetNodeObj.getData().detail
    this.editMode = this.GetNodeObj.getData().editMode
    this.showDeviceUuid = this.GetNodeObj.getData().showDeviceUuid
    this.IsToolBox = this.GetNodeObj.getData().IsToolBox
    _t.$EventBus.$on('cell-editMode', (data) => {
      _t.editMode = data.edit
      _t.IsToolBox = data.toolbox
      _t.initComponents(_t.detail)
    })
  },
  watch: {
    detail: {
      handler(newVal, oldVal) {
        if(this.editMode) {
          this.initComponents(newVal);
        }
      },
      deep: true
    }
  },
  methods: {
    initComponents(option){
      if(this.IsToolBox)
      {
        return
      }
      let i=0
      this.backColor=this.detail.style.backColor
      this.foreColor=this.detail.style.foreColor
      for( i=0;i<option.style.diy.length;i++)
      {
        if(option.style.diy[i].key=="hoverColor")
        {
          this.hoverColor=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="selectedColor")
        {
          this.selectedColor=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="TextFontSize")
        {
          this.TextFontSize=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="selectedTextColor")
        {
          this.selectedTextColor=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="SearchColor")
        {
          this.SearchColor=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="SearchBackColor")
        {
          this.SearchBackColor=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="SearchBorderColor")
        {
          this.SearchBorderColor=option.style.diy[i].value
        }

        else if(option.style.diy[i].key=="dateSelectColor")
        {
          this.dateSelectColor=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="dateSelectBackColor")
        {
          this.dateSelectBackColor=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="dateSelectBorderColor")
        {
          this.dateSelectBorderColor=option.style.diy[i].value
        }

        else if(option.style.diy[i].key=="tableHeaderColor")
        {
          this.tableHeaderColor=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="tableHeaderBackColor")
        {
          this.tableHeaderBackColor=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="tableSplitColor")
        {
          this.tableSplitColor=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="tableHoverColor")
        {
          this.tableHoverColor=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="ShowCount")
        {
          this.pagination.pageSize=parseInt(option.style.diy[i].value)
        }

      }
    },
  }
}
</script>
<style lang="less" scoped>
::v-deep .ant-form-item-label > label {
  color: var(--foreColor)
}
::v-deep .ant-card {
  background-color: var(--backColor);
}
::v-deep  .ant-table-placeholder {
  background:var(--backColor);
  border-top: 1px solid var(--backColor);
  border-bottom: 1px solid var(--backColor);
  border-radius: 0 0 4px 4px;
}
::v-deep  .ant-btn-primary {
  color: var(--foreColor);
  background-color: var(--SearchBackColor);
  border-color: var(--SearchBorderColor);
}
::v-deep .ant-input .ant-radio-button-wrapper {
  background-color: var(--SearchBackColor);
  border: 1px solid var(--SearchBorderColor);
}
::v-deep  .ant-select-selection {
  background-color: var(--SearchBackColor)
}
::v-deep  .ant-select-selection {
  border: 1px solid var(--SearchBorderColor);
}
::v-deep  .ant-radio-button-wrapper{
  background: var(--SearchBackColor);
  border: 1px solid var(--SearchBorderColor);
}
::v-deep  .ant-input {
  background-color: var(--SearchBackColor);
  border: 1px solid var(--SearchBorderColor);
}
::v-deep .ant-radio-button-wrapper-checked {
  color: var(--dateSelectColor);
  background: var(--dateSelectBackColor);
  border-color:  var(--dateSelectBorderColor);
  -webkit-box-shadow: -1px 0 0 0 var(--dateSelectBorderColor);
  box-shadow: -1px 0 0 0 var(--dateSelectBorderColor);
}
::v-deep  .plus-icon-enter-active{
  transition: opacity .5s;
}
::v-deep .plus-icon-enter{
  opacity: 0;
}
::v-deep .plus-icon-leave-active{
  transition: opacity .5s;
}
::v-deep .plus-icon-leave-to{
  opacity: 0;
}
::v-deep .plus-icon-enter-to{
  opacity: 1;
}

::v-deep .code-box-actions {
  padding-top: 12px;
  text-align: center;
  opacity: .7;
  transition: opacity .3s;
}
::v-deep .code-box-meta .demo-description>h4, .code-box-meta>h4 {
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
::v-deep .code-box {
  border: 1px solid #ebedf0;
  border-radius: 2px;
  display: inline-block;
  width: 100%;
  position: relative;
  margin: 0 0 16px;
  transition: all .2s;
}

::v-deep .search{
  margin-bottom: 54px;
}
::v-deep .fold{
  width: calc(100% - 216px);
  display: inline-block
}
::v-deep .operator{
  margin-bottom: 18px;
}
@media screen and (max-width: 900px) {
  .fold {
    width: 100%;
  }
}

::v-deep  .ant-table-thead > tr > th
{
  padding: 10px 10px;
  overflow-wrap: break-word;
}
::v-deep  .ant-table-tbody > tr > td {
  padding: 7px 5px;
  overflow-wrap: break-word;
  color: var(--foreColor);
}

::v-deep .ant-table-tbody > tr:hover:not(.ant-table-expanded-row):not(.ant-table-row-selected) > td {
  background: var(--tableHoverColor);
  cursor: default;
}

::v-deep .ant-table-tbody > tr > td {
  border-bottom: 1px solid var(--tableSplitColor);
}

::v-deep  .ant-table-thead>tr>th {
  color: var(--tableHeaderColor);
  font-weight: 500;
  text-align: left;
  background: var(--tableHeaderBackColor);
  //border-bottom: 1px solid #e8e8e8;
  transition: background .3s ease;
}
::v-deep .ant-pagination-item{
  background-color: var(--SearchBackColor);
  border: 1px solid var(--SearchBorderColor);
}
::v-deep .ant-pagination-item-active {
  font-weight: 500;
  background:  var(--dateSelectBackColor);
  border-color: var(--dateSelectBorderColor);
}
::v-deep .ant-pagination-prev {
  background-color: var(--SearchBackColor);
  border: 1px solid var(--SearchBorderColor);
}

::v-deep .ant-pagination-next {
  background-color: var(--SearchBackColor);
  border: 1px solid var(--SearchBorderColor);
}

::v-deep  .ant-pagination-item-link {
  background-color: var(--SearchBackColor);
  border: 1px solid var(--SearchBorderColor);
}
</style>