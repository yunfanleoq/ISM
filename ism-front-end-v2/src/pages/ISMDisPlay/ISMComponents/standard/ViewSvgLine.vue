<template>

  <svg  xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"  :width="detail.style.position.w" :height="detail.style.position.h" x="0" y="0"    xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <defs>
        <marker
            :id="detail.identifier+'start'"
            :fill="fill"
            markerUnits="strokeWidth"
            :stroke-width="strokeWidth"
            viewBox="0 0 10 10"
            refX="7"
            refY="3.5"
            markerWidth="10"
            markerHeight="10"
            orient="auto"
        >
          <circle cy="3.5" cx="7" r="2" v-if="ConnectType==1"/>
          <polygon points="0 0, 10 3.5, 0 7" v-if="ConnectType==0"></polygon>
        </marker>
        <marker
          :id="detail.identifier+'end'"
          :fill="fill"
          markerUnits="strokeWidth"
          :stroke-width="strokeWidth"
          viewBox="0 0 10 10"
          refX="7"
          refY="3.5"
          markerWidth="10"
          markerHeight="10"
          orient="auto"
      >
          <circle cy="3.5" cx="7" r="2" v-if="ConnectType==1"/>
          <polygon points="10 0, 0 3.5, 10 7" v-if="ConnectType==0"></polygon>
        </marker>
      </defs>
    <g class="svg-el" :class="{'animated':true,[`${detail.style.animate}`]: true}"  :style="{'width':detail.style.position.w,'height':detail.style.position.h,'text-align':'center','opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <line :marker-start="ConnectSelect==1||ConnectSelect==2?'url(#' + detail.identifier + 'end)':''" :stroke="strokeColor" x1="0" :y1="detail.style.position.h/2" :x2="detail.style.position.w" :y2="detail.style.position.h/2" :stroke-width="strokeWidth" :marker-end="ConnectSelect==0||ConnectSelect==2?'url(#' + detail.identifier + 'start)':''"/>
<!--      闪烁-->
        <animate v-if="(isStart)&&animateType.includes('blink')&&!IsToolBox" attributeName="opacity"
                 values="0.1;1;0.1" :dur="blinkSpeed+'s'"
                 repeatCount="indefinite"/>

      <animate v-if="(statusStart)&&!IsToolBox" attributeName="opacity"
               values="0.1;1;0.1" :dur="blinkStatusSpeed+'s'"
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
import {mapState} from "vuex";
import store from "@/store";
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
    name: 'view-svg-line',
    inject: ['getNode'],
    data() {
      return {
        detail:null,
        IsToolBox:false,
        editMode:true,
        Text:"",
        DivOpacity:1,
        strokeColor:"#000000",
        fill:"#A1BFE2",
        strokeWidth:0.3,
        fillOpacity:1,
        strokeOpacity:1,
        animateType:"blink",
        startColor:"#74f808",
        stopColor:"#74f808",
        ConnectSelect:1,
        animateSpeed:0.5,
        animateSpinSpeed:0.5,
        spinDirection:0,
        blinkSpeed:0.5,
        blinkStatusSpeed:0.5,
        isStart:false,
        statusStart:false,
        ConnectType:0,
        base:{
          "text": "displayConfig.ToolBox.Diagram.line",
          "icon": "icon-one-line-arrow",
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
            "active": [
              {
                id:"ControlStatus",
                name:"configComponent.status.ControlStatus",
                result:0,
                isStatus:true,
                isSwitch:false,
                isImageStatus:false,
                isTextStatus:false,
                isLineStatus:true,
                isExpression:false,
                condition:{
                  deviceSN:"",
                  isBandDevice:false,
                  bandType:1,
                  dataID: "",
                  dataName: "",
                  IsManual:false,
                  StatusList:[
                    {
                      "StatusOpt":"==",
                      "TextColor":"#000000",
                      "Blink":'0',
                      "BlinkSpeed":1,
                      "value2":1,
                      "value":0
                    }
                  ]
                },
              }
            ],
            "style": {
              "position": {
                "x": 0,
                "y": 0,
                "w": 300,
                "h": 50
              },
              "visible":1,
              "backColor": "transparent",
              "zIndex": -1,
              "transform": 0,
              "diy":[
                {
                  "name":"component.public.strokeWidth",
                  "type":7,
                  "value":1,
                  "min":0,
                  "key":"strokeWidth",
                },
                {
                  "name":"component.line.ConnectColor",
                  "type":2,
                  "value":"#f7032c",
                  "key":"strokeFill",
                },
                {
                  "name":"component.public.strokeColor",
                  "type":2,
                  "value":"#000000",
                  "key":"strokeColor",
                },
                {
                  "name":"component.public.fillOpacity",
                  "type":7,
                  "value":1,
                  "min":0,
                  "max":1,
                  "key":"fillOpacity",
                },
                {
                  name:"component.line.ConnectSelect",
                  type:6,
                  value:3,
                  enumList:[
                    {
                      value:0,
                      option:"component.line.ConnectStart"
                    },
                    {
                      value:1,
                      option:"component.line.ConnectEnd"
                    },
                    {
                      value:2,
                      option:"component.line.ConnectBoth"
                    },
                    {
                      value:3,
                      option:"component.line.ConnectNothing"
                    },
                  ],
                  min:1,
                  key:"ConnectSelect",
                }
                ,
                {
                  name:"component.line.ConnectType",
                  type:6,
                  value:0,
                  enumList:[
                    {
                      value:0,
                      option:"component.line.ConnectTypeArrow"
                    },
                    {
                      value:1,
                      option:"component.line.ConnectTypeCircle"
                    }
                  ],
                  min:1,
                  key:"ConnectType",
                }
              ]
            }
          }
        }
      }
    },
    computed: {
    ...mapState({
      ISMCavasContainer:state => store.state.ISMDisPlayEditorTool.ISMCavasContainer,
    }),
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
        for(let i=0;i<option.active[0].condition.StatusList.length;i++)
        {
          option.active[0].condition.StatusList[i].statusStart=0
        }
        this.StatusList = option.active[0].condition.StatusList
        if(option.active[0].condition.StatusList.length)
        {
          this.StatusValue = parseFloat(option.active[0].condition.StatusList[0].value)
        }
        let i=0
        for( i=0;i<option.style.diy.length;i++)
        {
          if(option.style.diy[i].key=="strokeWidth")
          {
            this.strokeWidth=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="strokeFill")
          {
            this.fill=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="strokeColor")
          {
            this.strokeColor=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="fillOpacity")
          {
            this.fillOpacity=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="strokeOpacity")
          {
            this.strokeOpacity=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="ConnectSelect")
          {
            this.ConnectSelect=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="ConnectType")
          {
            this.ConnectType=option.style.diy[i].value
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
      }
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
           if(data.ID == "ControlStatus")
           {
             let TempStatusValue = parseFloat(data.result)
             _t.StatusValue = TempStatusValue
             let StatusList=_t.detail.active[0].condition.StatusList
             let StatusListLen=_t.detail.active[0].condition.StatusList.length
             for(let i=0;i<StatusListLen;i++)
             {
               switch(StatusList[i].StatusOpt)
               {
                 case "==":{
                   if(_t.StatusValue==StatusList[i].value)
                   {
                     _t.strokeColor = StatusList[i].TextColor
                     if(StatusList[i].Blink=="1")
                     {
                       _t.detail.active[0].condition.StatusList[i].statusStart=true
                     }
                     else
                     {
                       _t.detail.active[0].condition.StatusList[i].statusStart=false
                     }
                   }
                   else
                   {
                     _t.detail.active[0].condition.StatusList[i].statusStart=false
                   }
                   break
                 }
                 case ">":{
                   if(_t.StatusValue>StatusList[i].value)
                   {
                     _t.strokeColor = StatusList[i].TextColor
                     if(StatusList[i].Blink=="1")
                     {
                       _t.detail.active[0].condition.StatusList[i].statusStart=true
                     }
                     else
                     {
                       _t.detail.active[0].condition.StatusList[i].statusStart=false
                     }
                   }
                   else
                   {
                     _t.detail.active[0].condition.StatusList[i].statusStart=false
                   }
                   break
                 }
                 case ">=":{
                   if(_t.StatusValue>=StatusList[i].value)
                   {
                     _t.strokeColor = StatusList[i].TextColor
                     if(StatusList[i].Blink=="1")
                     {
                       _t.detail.active[0].condition.StatusList[i].statusStart=true
                     }
                     else
                     {
                       _t.detail.active[0].condition.StatusList[i].statusStart=false
                     }
                   }
                   else
                   {
                     _t.detail.active[0].condition.StatusList[i].statusStart=false
                   }
                   break
                 }
                 case "<":{
                   if(_t.StatusValue<StatusList[i].value)
                   {
                     _t.strokeColor = StatusList[i].TextColor
                     if(StatusList[i].Blink=="1")
                     {
                       _t.detail.active[0].condition.StatusList[i].statusStart=true
                     }
                     else
                     {
                       _t.detail.active[0].condition.StatusList[i].statusStart=false
                     }
                   }
                   else
                   {
                     _t.detail.active[0].condition.StatusList[i].statusStart=false
                   }
                   break
                 }
                 case "<=":{
                   if(_t.StatusValue<=StatusList[i].value)
                   {
                     _t.strokeColor = StatusList[i].TextColor
                     if(StatusList[i].Blink=="1")
                     {
                       _t.detail.active[0].condition.StatusList[i].statusStart=true
                     }
                     else
                     {
                       _t.detail.active[0].condition.StatusList[i].statusStart=false
                     }
                   }
                   else
                   {
                     _t.detail.active[0].condition.StatusList[i].statusStart=false
                   }
                   break
                 }
                 case "<>":{
                   if((_t.StatusValue>=StatusList[i].value)&&(_t.StatusValue<=StatusList[i].value2))
                   {
                     _t.strokeColor = StatusList[i].TextColor
                     if(StatusList[i].Blink=="1")
                     {
                       _t.detail.active[0].condition.StatusList[i].statusStart=true
                     }
                     else
                     {
                       _t.detail.active[0].condition.StatusList[i].statusStart=false
                     }
                   }
                   else
                   {
                     _t.detail.active[0].condition.StatusList[i].statusStart=false
                   }
                   break
                 }
                 case "<!>":{
                   if((_t.StatusValue<=StatusList[i].value)||(_t.StatusValue>=StatusList[i].value2))
                   {
                     _t.strokeColor = StatusList[i].TextColor
                     if(StatusList[i].Blink=="1")
                     {
                       _t.detail.active[0].condition.StatusList[i].statusStart=true
                     }
                     else
                     {
                       _t.detail.active[0].condition.StatusList[i].statusStart=false
                     }
                   }
                   else
                   {
                     _t.detail.active[0].condition.StatusList[i].statusStart=false
                   }
                   break
                 }
                 case "!=":{
                   if(_t.StatusValue!=StatusList[i].value)
                   {
                     _t.strokeColor = StatusList[i].TextColor
                     if(StatusList[i].Blink=="1")
                     {
                       _t.detail.active[0].condition.StatusList[i].statusStart=true
                     }
                     else
                     {
                       _t.detail.active[0].condition.StatusList[i].statusStart=false
                     }
                   }
                   else
                   {
                     _t.detail.active[0].condition.StatusList[i].statusStart=false
                   }
                   break
                 }
               }
             }
             let is_find = 0
             for(let i=0;i<StatusListLen;i++)
             {
                if(_t.detail.active[0].condition.StatusList[i].statusStart)
                {
                  is_find=1;
                  _t.blinkStatusSpeed =_t.detail.active[0].condition.StatusList[i].BlinkSpeed
                  break
                }
             }
             if(is_find==1)
             {
               _t.statusStart=true

             }
             else{
               _t.statusStart=false
             }
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
    this.initComponents(this.detail);
  }
}
</script>
<style >
.svg-el {
  /*transform: rotate(45deg);*/
  transform-origin: center center;
}
</style>
