<template>

  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"  x="0px" y="0px"  xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el" :class="{'animated':true,[`${detail.style.animate}`]: true}" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <foreignObject style="overflow:visible;" pointer-events="all" :width="detail.style.position.w" :height="detail.style.position.h">
        <div :style="styleVar">
            <a-table :columns="columns" :data-source="TableData" :pagination="pagination" :loading="loading">
              <template #dateCell="{text}">
                {{ isDate(text) ? formatDate(text) : '无效日期' }}
              </template>
            </a-table>
        </div>

      </foreignObject>
<!--      闪烁-->
        <animate v-if="isStart&&animateType.includes('blink')&&!IsToolBox" attributeName="opacity"
                 values="0.1;1;0.1" :dur="blinkSpeed+'s'"
                 repeatCount="indefinite"/>
<!--渐变-->
        <animate v-if="isStart&&animateType.includes('millcolorGrad')&&!IsToolBox" attributeName="fill"
                 :values="startColor+';'+stopColor+';'+startColor" :dur="animateSpeed+'s'"
                 repeatCount="indefinite"/>
<!--缩放      -->
      <animateTransform v-if="isStart&&animateType.includes('Zoom')&&!IsToolBox" attributeName="transform"   begin="0s" dur="0.6s" type="scale" values="0.9;1;0.9" repeatCount="indefinite"/>
<!--      顺时针旋转-->
      <animateTransform v-if="isStart&&animateType.includes('animateSpin')&&!IsToolBox&&spinDirection==0" attributeType="XML" attributeName="transform" :dur="animateSpinSpeed+'s'" type="rotate" from="0 0 0" to="360 0 0" repeatCount="indefinite" />
<!--      逆时针旋转-->
      <animateTransform v-if="isStart&&animateType.includes('animateSpin')&&!IsToolBox&&spinDirection==1" attributeType="XML" attributeName="transform" :dur="animateSpinSpeed+'s'" type="rotate" from="360 0 0" to="0 0 0" repeatCount="indefinite" />
  </g>
</svg>

</template>

<script>
import {ExecSysScript, ISMExecSqlQuery} from "@/services/system";
import {formatDate} from "@/utils/common";
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
export default {
  mixins: [ISMChildAutoMixin],
    name: 'view-svg-sql-query',
    inject: ['getNode'],
    computed: {
      animatedStyle(){
        return {
          "--blinkSpeed":this.blinkSpeed+'s',
          "--stopColor":this.stopColor,
          "--startColor":this.startColor,
          "--animateSpeed":this.animateSpeed+'s',
          "--animateSpinSpeed":this.animateSpinSpeed+'s'
        }
      },
      textAlign: function(){
        if(this.detail.style.textAlign == undefined) {
          return "center";
        } else {
          return this.detail.style.textAlign;
        }
      },
      lineHeight: function() {
        if(this.detail.style.lineHeight == undefined) {
          return this.detail.style.position.h;
        }
        return this.detail.style.lineHeight;
      },
      styleVar() {
        return {
          "height": this.detail.style.position.h + 'px',
          "--foreColor": this.detail.style.foreColor,
          '--backColor': this.detail.style.backColor,
          "--selectedColor": this.selectedColor ,
          '--hoverColor': this.hoverColor,
          '--selectedTextColor': this.selectedTextColor,
          '--TextFontSize': this.TextFontSize+'px',
          '--hoverTextColor': this.hoverTextColor,
          '--SearchColor': this.SearchColor,
          '--SearchBackColor': this.SearchBackColor,
          '--SearchBorderColor': this.SearchBorderColor,

          '--tableHeaderColor': this.tableHeaderColor,
          '--tableHeaderBackColor': this.tableHeaderBackColor,
          '--tableSplitColor': this.tableSplitColor,
          '--tableHoverColor': this.tableHoverColor,

          '--dateSelectColor': this.dateSelectColor,
          '--dateSelectBackColor': this.dateSelectBackColor,
          '--dateSelectBorderColor': this.dateSelectBorderColor,

          '--fontWeight': this.detail.style.fontWeight,
          '--textAlign': this.detail.style.textAlign,
          '--fontSize': this.detail.style.fontSize + 'px',
          '--letterSpacing': this.detail.style.letterSpacing+'px',
          '--font-style': this.detail.style.italic?'oblique':'normal',
          '--fontFamily': this.detail.style.fontFamily
        };
      },
    },
    filters: {
      formatDate(time) {
        let date = new Date(time)
        return formatDate(date,'yyyy-MM-dd hh:mm:ss')
      },
    },
    data() {
      return {
        detail:null,
        IsToolBox:false,
        editMode:true,
        Text:"",
        DivOpacity:1,
        pagination: {
          pageSize: 15,
          showSizeChanger: false,
          hideOnSinglePage:true,
          showLessItems:true,
          simple:true
        },
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
        loading: false,
        strokeColor:"#000000",
        fill:"#A1BFE2",
        strokeWidth:0.3,
        fillOpacity:1,
        strokeOpacity:1,
        columns:[],
        animateType:"blink",
        startColor:"#74f808",
        stopColor:"#74f808",
        animateSpeed:0.5,
        animateSpinSpeed:0.5,
        spinDirection:0,
        blinkSpeed:0.5,
        isStart:false,
        italic:false,
        TableError:0,
        SqlQuery:"",
        DataFrom:1,
        TimelyExecSql:null,
        SqlQueryTimely:30,
        SystemScript:"",
        TableData:[
          {
            key: '1',
            name: 'John Brown',
            age: 32,
            address: 'New York No. 1 Lake Park',
            tags: ['nice', 'developer'],
          },
          {
            key: '2',
            name: 'Jim Green',
            age: 42,
            address: 'London No. 1 Lake Park',
            tags: ['loser'],
          },
        ],
        base:{
          text: "configComponent.sqlTable.Text",
          "icon": "icon-biaoge",
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
                "w": 619,
                "h": 460
              },
              "visible":1,
              "zIndex": -1,
              "transform": 0,
              "backColor": "transparent",
              "foreColor": "#000000",
              fontWeight:400,
              textAlign: "left",
              fontSize: 14,
              fontFamily: "Arial",
              italic:0,
              "diy":[
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
                {
                  "name":"configComponent.sqlTable.TableHeader",
                  "type":14,
                  "value":"[{\"title\":\"年龄\",\"dataIndex\":\"age\"\},{\"title\":\"名称\",\"dataIndex\":\"name\"},{\"title\":\"地址\",\"dataIndex\":\"address\"}]",
                  "rows":30,
                  "key":"TableHeader",
                },
                {
                  name:"configComponent.sqlTable.DataFrom",
                  type:6,
                  value:1,
                  enumList:[
                    {
                      value:0,
                      option:"configComponent.sqlTable.DataFromScript"
                    },
                    {
                      value:1,
                      option:"configComponent.sqlTable.DataFromSql"
                    }
                  ],
                  min:1,
                  key:"DataFrom",
                },
                {
                  "name":"configComponent.sqlTable.SqlQuery",
                  "type":30,
                  "value":"",
                  "rows":30,
                  "key":"SqlQuery",
                },
                {
                  "name":"configComponent.sqlTable.SystemScript",
                  "type":31,
                  "value":"",
                  "key":"SystemScript",
                },
                {
                  "name":"configComponent.sqlTable.SqlQueryTimely",
                  "type":1,
                  "value":30,
                  "key":"SqlQueryTimely",
                },
                {
                  "name":"configComponent.sqlTable.PageSize",
                  "type":1,
                  "value":10,
                  "key":"PageSize",
                }
              ]
            }
          }
        }
      }
    },
    watch: {
      detail: {
        handler(newVal, oldVal) {
          this.initComponents(newVal);
        },
        deep: true
      }
    },
    methods: {
      initComponents(option){
        let _t = this
        if(this.IsToolBox)
        {
          return
        }
        let i=0
        clearInterval(this.TimelyExecSql)
        for( i=0;i<option.style.diy.length;i++)
        {
          if(option.style.diy[i].key=="TableHeader")
          {
              try{
                this.columns=JSON.parse(option.style.diy[i].value)
              }catch (e) {
                this.TableError = -1
              }
          }
          else if(option.style.diy[i].key=="SqlQuery")
          {
            this.SqlQuery=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="SqlQueryTimely")
          {
            this.SqlQueryTimely=option.style.diy[i].value*1000
          }
          else if(option.style.diy[i].key=="PageSize")
          {
            this.pagination.pageSize=parseInt( option.style.diy[i].value)
          }else if(option.style.diy[i].key=="hoverColor")
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
          else if(option.style.diy[i].key=="DataFrom")
          {
            this.DataFrom=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="SystemScript")
          {
            this.SystemScript=option.style.diy[i].value
          }
        }
        i=0
        this.animateType = option.animate.selected
        if(option.animate.isExpression)
        {
          this.isStart = false
        }
        else
        {
          this.isStart = true
        }
        for( i=0;i<option.animate.animateElement.length;i++)
        {
          if(option.animate.animateElement[i].id=="millcolorGrad")
          {
            for(let k =0;k<option.animate.animateElement[i].elementList.length;k++)
            {
              if(option.animate.animateElement[i].elementList[k].key=="startColor")
              {
                this.startColor=option.animate.animateElement[i].elementList[k].value
              }
              else if(option.animate.animateElement[i].elementList[k].key=="stopColor")
              {
                this.stopColor=option.animate.animateElement[i].elementList[k].value
              }
              else if(option.animate.animateElement[i].elementList[k].key=="animateSpeed")
              {
                this.animateSpeed=option.animate.animateElement[i].elementList[k].value
              }
            }
          }
          else if(option.animate.animateElement[i].id=="blink")
          {
            for(let k =0;k<option.animate.animateElement[i].elementList.length;k++) {
              if (option.animate.animateElement[i].elementList[k].key == "blinkSpeed") {
                this.blinkSpeed = option.animate.animateElement[i].elementList[k].value
              }
            }
          }
          else if(option.animate.animateElement[i].id=="animateSpin")
          {
            for(let k =0;k<option.animate.animateElement[i].elementList.length;k++) {
              if (option.animate.animateElement[i].elementList[k].key == "spinSpeed") {
                this.animateSpinSpeed = option.animate.animateElement[i].elementList[k].value
              }
              else if (option.animate.animateElement[i].elementList[k].key == "spinDirection") {
                this.spinDirection = option.animate.animateElement[i].elementList[k].value
              }
            }
          }
        }

        if(this.SqlQueryTimely<=5000)
        {
          this.SqlQueryTimely=5000
        }
        if(_t.DataFrom==1)
        {
          _t.ExecSqlQuery()
        }
        else
        {
          _t.ExecScriptQuery()
        }
          this.TimelyExecSql = setInterval(function () {
            if(_t.DataFrom==1)
            {
              _t.ExecSqlQuery()
            }
            else
            {
              _t.ExecScriptQuery()
            }
          },this.SqlQueryTimely)

      },
      ExecSqlQuery(){
        let _t = this
        if(this.SqlQuery.length==0)
        {
          return
        }
        let params = {
          sql: this.SqlQuery
        }
        this.loading = true;
        ISMExecSqlQuery(params).then(function (res) {
          let result = res.data
          _t.TableData=[]
          if(result.code==0)
          {
            _t.TableData = result.data
          }
          _t.loading = false;
        }).catch(function(e){
          _t.loading = false;
        })
      },
      ExecScriptQuery(){
        let _t =this
        if(this.SystemScript.length==0)
        {
          return
        }
        let params = {
          Script:[]
        };
        params.Script[0]=this.SystemScript
        this.loading = true;
        ExecSysScript(params).then(function (res){
          _t.TableData=[]
          let result = res.data
          if(result.code==0)
          {
            _t.TableData = result.result
          }
        }).catch(function (error) {
          console.log(error)
          _t.loading = false;
        }).finally(function (error) {
          _t.loading = false;
        })
      }
    },
    beforeDestroy() {
      clearInterval(this.TimelyExecSql)
    },
    mounted() {
      let _t = this
      this.$nextTick(function(){
       this.initComponents(this.detail);
        let activeEvent = this.detail.identifier+"activeEvent"//动作数据
        let animateEvent = this.detail.identifier+"animateEvent"//动作数据

        _t.$EventBus.$on(activeEvent, (data) => {


        })
        _t.$EventBus.$on(animateEvent, (data) => {
          if((_t.editMode)&&(!this.IsToolBox)){
            return
          }
          _t.isStart = data
        })

      });
    },
    created(){
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
    })
  }
}
</script>
<style lang="less" scoped>
.svg-el {
  /*transform: rotate(45deg);*/
  transform-origin: center center;
}


::v-deep  .ant-table-thead > tr > th
{
  padding: 10px 10px;
  overflow-wrap: break-word;
}
::v-deep  .ant-table-tbody > tr > td {
  padding: 7px 5px;
  overflow-wrap: break-word;
  text-align: var(--textAlign);
  font-size:var(--fontSize);
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
  font-weight: var(--fontWeight);
  text-align: var(--textAlign);
  font-style:var(--font-style);
  font-size:var(--fontSize);
  font-family:var(--fontFamily);
  background: var(--tableHeaderBackColor);
  border-bottom: 1px solid var(--tableSplitColor);
  transition: background .3s ease;
}
::v-deep .ant-pagination-prev {
  background-color: var(--backColor);
  border: 0px solid var(--foreColor);
}

::v-deep .ant-pagination-next {
  background-color: var(--backColor);
  border: 0px solid var(--foreColor);
}
::v-deep  .ant-pagination-next:hover  .anticon{
  color: var(--tableHoverColor); /* 悬停时的边框色 */
}
::v-deep  .ant-pagination-prev:hover  .anticon{
  color: var(--tableHoverColor); /* 悬停时的边框色 */
}
::v-deep .ant-pagination-next:focus .anticon,
.ant-pagination-next-active .anticon {
  color: var(--foreColor);
}
::v-deep .ant-pagination-prev:focus .anticon,
.ant-pagination-prev-active .anticon {
  color: var(--foreColor);
}
::v-deep  .ant-pagination-item-link {
  background-color: var(--backColor);
  border: 1px solid var(--backColor);
  color: var(--foreColor);
  font-style:var(--font-style);
  font-size:var(--fontSize);
}

::v-deep  .ant-pagination-simple-pager{
  color: var(--foreColor);
  font-style:var(--font-style);
  font-size:var(--fontSize);
}
::v-deep  .ant-pagination-simple-pager>input{
  color: var(--foreColor);
  background-color: var(--backColor);
  font-style:var(--font-style);
  font-size:var(--fontSize);
  border: 1px solid var(--foreColor);
}


</style>
