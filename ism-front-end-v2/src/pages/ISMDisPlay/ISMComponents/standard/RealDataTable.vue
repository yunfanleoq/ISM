<template>
  <div v-show="detail.style.visible==1 ||isStart? true:false" :style="animatedStyle">
  <div :class="{
          'animated':true,[`${detail.style.animate}`]: true,
          'color-animation':isStart&&animateType.includes('millcolorGrad')&&!IsToolBox,
          'blink-animation':isStart&&animateType.includes('blink')&&!IsToolBox,
          'scale-animation':isStart&&animateType.includes('Zoom')&&!IsToolBox,
          'rotate-animation':isStart&&animateType.includes('animateSpin')&&!IsToolBox&&spinDirection==0,
          'rotate-anti-animation':isStart&&animateType.includes('animateSpin')&&!IsToolBox&&spinDirection==1
        }"
       :style="{
                                width: detail.style.position.w + 'px',
                                height: detail.style.position.h + 'px',
                                'background-color': detail.style.backColor,
                                'border-radius':detail.style.BorderEdges+'px',
                                opacity:detail.style.opacity,
                                borderWidth: detail.style.borderWidth + 'px',
                                borderStyle: detail.style.borderStyle,
                                borderColor: detail.style.borderColor,
                                transform: detail.style.transform==-1099? 'rotateY(180deg)':detail.style.transform==-1098?'rotateX(180deg)':'',
                            }">
          <div class="usr-table-wrapper" :style="{'height':detail.style.position.h+'px','width':detail.style.position.w+'px'}">
            <div class="usr-table-scroll">
              <table border="1" class="usr_table_style">
                <!-- 表格主体 		 -->
                <tr v-for="tr in tableRowIndexes" :key="tr">
                  <template v-for="td in tableCols">
                    <td :key="td"></td>
                  </template>
                </tr>

              </table>
            </div>
          </div>

  </div>
  </div>
</template>

<script>

import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
function getCell(dom) {
  if (dom.tagName === 'BODY') return null
  return (dom.tagName === 'TH' || dom.tagName === 'TD') ? dom : getCell(dom.parentNode)
}
export default {
  mixins: [ISMChildAutoMixin],
  name: 'view-data-table',
  inject: ['getNode'],
  i18n: require('@/i18n/language'),
  data() {
    return {
      detail:null,
      IsToolBox:false,
      dataComplete: false, // 数据加载完成标志
      reflush: true, // 强制刷新表格
      isShowNoData: false,
      DivOpacity:1,
      strokeColor:"#000000",
      fill:"#A1BFE2",
      strokeWidth:0.3,
      fillOpacity:1,
      strokeOpacity:1,
      animateType:"blink",
      startColor:"#74f808",
      stopColor:"#74f808",
      animateSpeed:0.5,
      animateSpinSpeed:0.5,
      spinDirection:0,
      blinkSpeed:0.5,
      isStart:false,
      imageURL:"",
      AlarmTimer:null,
      IsShowIndex:1,
      RealDataTableHeaderHeight:50,
      tableRows:5,
      tableCols:5,
      tableBorderColor:"#b3b3b3",
      editMode:true,
      base:{
        text: "configComponent.table.title",
        "icon": "icon-charubiaoge",
        "isFontIcon": true,
        "info": {
          "type": "image",
          "action": [],
          "active": [],
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
              "w": 210,
              "h": 140
            },
            "visible":1,
            "backColor": "transparent",
            "zIndex": -1,
            "transform": 0,
            "diy":[
              {
                "name":"configComponent.table.tableBorderColor",
                "type": 2,
                "value": "#b3b3b3",
                "key":"tableBorderColor",
              },
              {
                "name":"configComponent.table.rows",
                "type":1,
                "value":5,
                "min":1,
                "max":100,
                "key":"rows",
              },
              {
                "name":"configComponent.DeviceTree.ShowCount",
                "type":1,
                "value":5,
                "min":1,
                "max":100,
                "key":"ShowCount",
              },
              {
                "name":"configComponent.table.cols",
                  "type":1,
                  "value":5,
                  "min":1,
                  "max":100,
                  "key":"cols",
              }
            ]
          }
        }
      }
    }
  },
  components: {

  },
  computed: {
    tableRowIndexes(){
      return Array.from({ length: this.tableRows }, (item, index) => index + 1)
    },
    animatedStyle(){
      return {
        "--tableBorderColor": this.tableBorderColor ,
        "--blinkSpeed":this.blinkSpeed+'s',
        "--stopColor":this.stopColor,
        "--startColor":this.startColor,
        "--animateSpeed":this.animateSpeed+'s',
        "--animateSpinSpeed":this.animateSpinSpeed+'s',
        "--SearchColor": "#000000",
        "--SearchBackColor": "#ffffff",
        "--SearchBorderColor": "#cbc6c6",
        "--toolbarAccentSoft": "rgba(47, 111, 237, 0.12)",
        "--toolbarAccent": "#2f6fed"
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
    }
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
      for( i=0;i<option.style.diy.length;i++)
      {
        if(option.style.diy[i].key=="tableBorderColor")
        {
          this.tableBorderColor=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="rows")
        {
          this.tableRows=parseInt(option.style.diy[i].value)
        }
        else if(option.style.diy[i].key=="cols")
        {
          this.tableCols=parseInt(option.style.diy[i].value)
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
    },
  },
  beforeDestroy(){

  },
  mounted() {
    let _t = this
    this.$nextTick(function(){
      this.initComponents(this.detail);
      let activeEvent = this.detail.identifier+"activeEvent"//动作数据
      let animateEvent = this.detail.identifier+"animateEvent"//动作数据

      _t.$EventBus.$on(activeEvent, (data) => {
        if((_t.editMode)&&(!this.IsToolBox)){
          return
        }
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
    const node = this.getNode()
    node.on('change:data', ({ current }) => {
      if(current) {
        _t.detail = current.detail
      }
    })
    node.on('change:size', ({ current }) => {
      _t.detail.style.position.w = current.width
      _t.detail.style.position.h = current.height
    });
    this.detail = node.getData().detail
    this.editMode = node.getData().editMode
    this.showDeviceUuid = node.getData().showDeviceUuid
    this.IsToolBox = node.getData().IsToolBox
    _t.$EventBus.$on('cell-editMode', (data) => {
      _t.editMode = data.edit
      _t.IsToolBox = data.toolbox
    })
    this.initComponents(this.detail);
  }
}
</script>
<style lang="less">
.svg-el {
  /*transform: rotate(45deg);*/
  transform-origin: center center;
}

/* 使用animation关键帧 */
.color-animation {
  animation: colorChange var(--animateSpeed) linear infinite;
}

@keyframes colorChange {
  0% { background-color: var(--startColor); }
  100% { background-color:  var(--stopColor); }
}
/* 使用animation关键帧 */
.blink-animation {
  animation: blink var(--blinkSpeed) linear infinite;
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0; }
}
/*缩放*/
.scale-animation {
  animation: pulse 0.6s infinite;
}

@keyframes pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.5); }
}
/*顺时针旋转*/
.rotate-animation {
  animation: clockwiseRotate var(--animateSpinSpeed) linear infinite;
  transform-origin: center;
}

@keyframes clockwiseRotate {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
/*逆时针旋转*/
.rotate-anti-animation {
  animation: counterClockwiseRotate var(--animateSpinSpeed) linear infinite;
  transform-origin: center;
}

@keyframes counterClockwiseRotate {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(-360deg); }
}
.usr-table-wrapper {
  display: flex;
  flex-direction: column;
  min-height: 0;
  overflow: hidden;
}

.usr-table-scroll {
  flex: 1 1 auto;
  min-height: 0;
  overflow: hidden;
}

.usr_table_style {
  width: 100%;
  height: calc(100% - 2px);
  text-align: center;
  border-collapse: collapse;
  border: var(--tableBorderColor);
}

.table-pagination-bar {
  flex: 0 0 auto;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 0;
  height: 0;
  padding: 0;
  margin: 0 auto;
  width: fit-content;
  border: none;
  border-radius: 0;
  background: transparent;
  backdrop-filter: none;
  box-shadow: none;
  opacity: 0;
  transform: translateY(4px);
  pointer-events: none;
  overflow: hidden;
  transition: opacity 0.18s ease, transform 0.18s ease;
}

.table-pagination-bar.visible {
  min-height: 18px;
  height: 24px;
  padding: 0 8px;
  margin: 0 auto 2px;
  width: fit-content;
  border-width: 1px;
  opacity: 1;
  transform: translateY(0);
  pointer-events: auto;
}

::v-deep .table-pagination-bar .ant-pagination-simple {
  font-size: 11px !important;
}

::v-deep .table-pagination-bar .ant-pagination-prev,
::v-deep .table-pagination-bar .ant-pagination-next,
::v-deep .table-pagination-bar .ant-pagination-simple-pager,
::v-deep .table-pagination-bar .ant-pagination-item-link {
  min-width: 18px !important;
  height: 18px !important;
  line-height: 16px !important;
}

::v-deep .table-pagination-bar .ant-pagination-simple-pager input {
  height: 16px !important;
  min-width: 36px !important;
  padding: 0 4px !important;
}

::v-deep .table-pagination-bar .ant-pagination-item,
::v-deep .table-pagination-bar .ant-pagination-prev,
::v-deep .table-pagination-bar .ant-pagination-next,
::v-deep .table-pagination-bar .ant-pagination-item-link {
  color: var(--SearchColor) !important;
  background: var(--SearchBackColor) !important;
  border-color: var(--SearchBorderColor) !important;
}

::v-deep .table-pagination-bar .ant-pagination-item-active {
  background: var(--toolbarAccentSoft) !important;
  border-color: var(--toolbarAccent) !important;
}
</style>
