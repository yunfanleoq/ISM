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
       <a-progress :strokeLinecap="strokeLinecap"  :stroke-color="{
        from: strokeColor,
        to: strokeEndColor,
      }" :showInfo="showInfo" :type="strokeType" :percent="percent" :strokeWidth="strokeWidth">
         <template #format="percent">
            <span :style="{
        fontSize: detail.style.fontSize + 'px',
        fontFamily: detail.style.fontFamily,
        'font-weight':detail.style.fontWeight,
        color: detail.style.foreColor,
        'font-style':detail.style.italic?'oblique':'normal'
    }">{{ percent }}{{ChartUnit}}</span>
         </template>
        </a-progress>
    </div>
  </div>

</template>

<script>
import {mapState} from "vuex";
import store from "@/store";
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
    name: 'view-progress-bars',
    i18n: require('../../../../i18n/language'),
    inject: ['getNode'],
    data() {
      return {
        detail:null,
        IsToolBox:false,
        DivOpacity:1,
        Text:"",
        foreColor:"#000000",
        backColor:"#ffffff",
        strokeColor:"#0becc7",
        strokeEndColor:"#87d068",
        strokeWidth:10,
        strokeType:"line",
        strokeLinecap:"round",
        showInfo:true,
        fill:"#A1BFE2",
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
        StatusList:[],
        sliderMax:100,
        sliderMin:0,
        sliderValue:10,
        HavedSliderColor:"",
        NoHavedSliderColor:"",
        SliderColor:"",
        SliderHeight:4,
        HoverHavedSliderColor:"",
        HoverNoHavedSliderColor:"",
        HoverSliderColor:"",
        Direction:1,
        ChartUnit:"%",
        percent:60,
        StatusValue:"loading",
        metreFullScale:100,
        metreMinScale:0,
        base:{
          text: "configComponent.ProgressBar.title",
          "icon": "icon-a-012",
          "isFontIcon": true,
          "info": {
            "type": "image",
            "action": [],
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
            "style": {
              "position": {
                "x": 0,
                "y": 0,
                "w": 300,
                "h": 30
              },
              "visible":1,
              "backColor": "transparent",
              "zIndex": -1,
              "foreColor": "#000000",
              fontWeight:400,
              "transform": 0,
              fontSize: 14,
              fontFamily: "Arial",
              "diy":[
                {
                  "name":"configComponent.ProgressBar.Type",
                  type:6,
                  value:0,
                  enumList:[
                    {
                      value:0,
                      option:"configComponent.ProgressBar.line"
                    },
                    {
                      value:1,
                      option:"configComponent.ProgressBar.circle"
                    }
                  ],
                  "key":"Type",
                },
                {
                  "name":"configComponent.ProgressBar.showInfo",
                  type:6,
                  value:0,
                  enumList:[
                    {
                      value:0,
                      option:"configComponent.ProgressBar.showInfoTrue"
                    },
                    {
                      value:1,
                      option:"configComponent.ProgressBar.showInfoFalse"
                    }
                  ],
                  "key":"ShowInfo",
                },
                {
                  "name":"configComponent.ProgressBar.strokeLinecap",
                  type:6,
                  value:0,
                  enumList:[
                    {
                      value:0,
                      option:"configComponent.ProgressBar.strokeLinecapround"
                    },
                    {
                      value:1,
                      option:"configComponent.ProgressBar.strokeLinecapsquare"
                    }
                  ],
                  "key":"strokeLinecap",
                },
                {
                  "name":"configComponent.ProgressBar.strokeWidth",
                  "type": 1,
                  "value": 10,
                  "key":"strokeWidth",
                },
                {
                  "name":"configComponent.ProgressBar.strokeColor",
                  "type": 2,
                  "value": "#0becc7",
                  "key":"strokeColor",
                },
                {
                  "name":"configComponent.ProgressBar.strokeEndColor",
                  "type": 2,
                  "value": "#87d068",
                  "key":"strokeEndColor",
                },
                {
                  "name":"configComponent.ChartPublic.ChartUnit",
                  "type":4,
                  "value":"%",
                  "key":"ChartUnit",
                },
                {
                  "name":"configComponent.ProgressBar.metreMinScale",
                  "type": 1,
                  "value": 0,
                  "key":"metreMinScale",
                },
                {
                  "name":"configComponent.ProgressBar.metreFullScale",
                  "type": 1,
                  "value": 100,
                  "key":"metreFullScale",
                },

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
        this.DivOpacity = option.style.opacity
        this.backColor=this.detail.style.backColor
        this.foreColor=this.detail.style.foreColor
        let i=0
        for( i=0;i<option.style.diy.length;i++)
        {
          if(option.style.diy[i].key=="Type")
          {
            if(option.style.diy[i].value==0)
            {
              this.strokeType="line"
            }
            else if(option.style.diy[i].value==1)
            {
              this.strokeType="circle"
            }
          }
          else if(option.style.diy[i].key=="ShowInfo")
          {
            if(option.style.diy[i].value==0)
            {
              this.showInfo=true
            }
            else
            {
              this.showInfo=false
            }
          }
          else if(option.style.diy[i].key=="strokeLinecap")
          {
            if(option.style.diy[i].value==0)
            {
              this.strokeLinecap="round"
            }
            else
            {
              this.strokeLinecap="square"
            }
          }
          else if(option.style.diy[i].key=="strokeWidth")
          {
            this.strokeWidth=parseInt(option.style.diy[i].value)
          }
          else if(option.style.diy[i].key=="strokeColor")
          {
            this.strokeColor=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="strokeEndColor")
          {
            this.strokeEndColor=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="ChartUnit")
          {
            this.ChartUnit=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="metreFullScale")
          {
            this.metreFullScale=parseFloat(option.style.diy[i].value)
          }
          else if(option.style.diy[i].key=="metreMinScale")
          {
            this.metreMinScale=parseFloat(option.style.diy[i].value)
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
            if(data.ID == "ShowVariable")
            {
              _t.percent = (parseFloat(data.result)-_t.metreMinScale)/(_t.metreFullScale-_t.metreMinScale)*100
            }
          })
          _t.$EventBus.$on(animateEvent, (data) => {
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
    this.initComponents(this.detail);
  }
}
</script>
<style lang="less">
::v-deep .ant-progress-show-info .ant-progress-outer {
  margin-right: calc(-2em - 15px);
  padding-right: calc(2em + 8px);
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
