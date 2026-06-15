<template>
  <div :style="{'background':detail.style.backColor}">
    <OperationLogger :ShowPageCount="ShowPageSize"/>
  </div>
</template>

<script>
import {formatDate} from '@/utils/common';
import OperationLogger from "@/pages/journal/OperationLog/OperationLog";
import canvasView from '../View';
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
export default {
  mixins: [ISMChildAutoMixin],
  name: 'ComOperationLogger',
  i18n: require('@/i18n/language'),
  inject: ['getNode'],
  components:{
    OperationLogger
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
      ShowPageSize:10,
      loadExecl:null,
      isLoadExecl:false,
      exportName:"",
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
      messageShowLoad:false,
      refIconLoading: false,
      base:{
        text: "journal.title",
        "icon": "icon-caozuorizhi",
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
            ]
          }
        }
      }
    }
  },
  authorize: {
    // deleteRecord: 'delete'
  },
  mounted(){

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
    this.initComponents(this.detail)
  },
  watch: {
    '$route' () {
      this.getMonitorTree()
    },
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
        if(option.style.diy[i].key=="ShowCount")
        {
          this.ShowPageSize=parseInt(option.style.diy[i].value)
        }
        console.log("this.ShowPageSize",this.ShowPageSize)
      }

    },
  }
}
</script>

<style lang="less" scoped>

</style>
