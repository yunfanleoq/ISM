<template>
<div :style="animatedStyle" v-show="detail.style.visible==1||isStart ? true:false">
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
            opacity:DivOpacity,
            borderWidth: detail.style.borderWidth + 'px',
            borderStyle: detail.style.borderStyle,
            borderColor: detail.style.borderColor,
            transform: detail.style.transform==-1099? 'rotateY(180deg)':detail.style.transform==-1098?'rotateX(180deg)':'',
        }">
        <div :style="{
        fontSize: detail.style.fontSize + 'px',
        fontFamily: detail.style.fontFamily,
        color: detail.style.foreColor,
        'font-weight':detail.style.fontWeight,
        textAlign: textAlign,
        lineHeight: lineHeight + 'px',
    }">
          {{Variable}}{{ChartUnit}}
        </div>
  </div>
</div>
</template>

<script>
import {mapState} from "vuex";
import store from "@/store";
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
export default {
  mixins: [ISMChildAutoMixin],
    name: 'view-svg-variable',
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
        animateSpeed:0.5,
        animateSpinSpeed:0.5,
        spinDirection:0,
        blinkSpeed:0.5,
        isStart:false,
        Variable:"###",
        ChartUnit:"",
        ShowJinZhi:5,
        GetNodeObj:null,
        ShowFloatNumber:-1,
        base:{
          text: "configComponent.variable.Text",
          "icon": "icon-bianliangbiao",
          "isFontIcon": true,
          "info": {
            "type": "image",
            "action": [],
            "dataBind":[],
            "active": [
              {
                id:"ShowVariable",
                name:"configComponent.variable.ShowData",
                result:"",
                isExpression:false,
                condition:{
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
              }
            ],
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
                  id: "millcolorGrad",
                  name: "component.public.millcolorGrad",
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
                "w": 132,
                "h": 60
              },
              "visible":1,
              "backColor": "transparent",
              "foreColor": "#000000",
              "zIndex": -1,
              "transform": 0,
              fontWeight:400,
              textAlign: "center",
              fontSize: 74,
              fontFamily: "黑体",
              "diy":[
                {
                  "name":"component.public.fillOpacity",
                  "type":7,
                  "value":1,
                  "min":0,
                  "max":1,
                  "key":"fillOpacity",
                },
                {
                  "name":"configComponent.ChartPublic.ChartUnit",
                  "type":4,
                  "value":"",
                  "key":"ChartUnit",
                },
                {
                  "name":"configComponent.ChartPublic.ShowJinZhi",
                  "type":6,
                  "enumList":[
                    {option:'configComponent.ChartPublic.ShowJinZhiNo',value:5},
                    {option:'configComponent.ChartPublic.ShowJinZhi10',value:1},
                    {option:'configComponent.ChartPublic.ShowJinZhi16',value:2},
                    {option:'configComponent.ChartPublic.ShowJinZhi8',value:3},
                    {option:'configComponent.ChartPublic.ShowJinZhi2',value:4}
                  ],
                  "value":5,
                  "key":"ShowJinZhi",
                },
                {
                  "name":"configComponent.ChartPublic.ShowFloatNumber",
                  "type":6,
                  "enumList":[
                    {option:'configComponent.ChartPublic.ShowFloatOld',value:-1},
                    {option:'0',value:0},
                    {option:'1',value:1},
                    {option:'2',value:2},
                    {option:'3',value:3},
                    {option:'4',value:4},
                    {option:'5',value:5},
                    {option:'6',value:6},
                    {option:'7',value:7},
                    {option:'8',value:8},
                  ],
                  "value":-1,
                  "key":"ShowFloatNumber",
                },
              ]
            }
          }
        }
      }
    },
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
        this.DivOpacity = option.style.opacity
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
          else if(option.style.diy[i].key=="imageURL")
          {
            this.imageURL=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="ChartUnit")
          {
            this.ChartUnit=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="ShowJinZhi")
          {
            this.ShowJinZhi=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="ShowFloatNumber")
          {
            this.ShowFloatNumber=option.style.diy[i].value
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
                this.blinkSpeed = parseFloat(option.animate.animateElement[i].elementList[k].value)
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
          if(data.ID == "ShowVariable")
          {
            if(_t.ShowJinZhi!=5)
            {
              let tenNumber = parseInt(data.result)
              if(_t.ShowJinZhi==1)
              {
                _t.Variable = tenNumber.toString()
              }
              else if(_t.ShowJinZhi==2)
              {
                _t.Variable = tenNumber.toString(16).toUpperCase()
              }
              else if(_t.ShowJinZhi==3)
              {
                _t.Variable = tenNumber.toString(8)
              }
              else if(_t.ShowJinZhi==4)
              {
                _t.Variable = tenNumber.toString(2)
              }
            }
            else
            {
              _t.Variable = data.result
            }
            if(_t.ShowFloatNumber!=-1)
            {
              _t.Variable = parseFloat(_t.Variable).toFixed(_t.ShowFloatNumber)
            }
          }
        })
        _t.$EventBus.$on(animateEvent, (data) => {
            _t.isStart = data
            if (_t.animateType.includes("visible")) {
              if (data) {
                _t.detail.style.visible = true
              } else {
                _t.detail.style.visible = false
              }
            }
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
      this.initComponents(this.detail);
    }
}
</script>
<style scoped lang="less">
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
</style>
