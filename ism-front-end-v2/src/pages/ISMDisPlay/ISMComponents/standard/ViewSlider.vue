<template>

  <div :style="animatedStyle" v-show="detail.style.visible==1 ||isStart? true:false">
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
            'background-color': detail.style.backColor,
            'border-radius':detail.style.BorderEdges+'px',
            opacity:DivOpacity,
            borderWidth: detail.style.borderWidth + 'px',
            borderStyle: detail.style.borderStyle,
            borderColor: detail.style.borderColor,
            transform: detail.style.transform==-1099? 'rotateY(180deg)':detail.style.transform==-1098?'rotateX(180deg)':'',
        }">
        <a-slider v-if="Direction==2" v-model="sliderValue" @afterChange="handleChange"  vertical :min="sliderMin" :max="sliderMax" />
        <a-slider v-else v-model="sliderValue" @afterChange="handleChange"  :min="sliderMin" :max="sliderMax" />
    </div>
  </div>

</template>

<script>
import {setData} from "@/services/device";
import {mapState} from "vuex";
import store from "@/store";
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
    name: 'view-slider',
    i18n: require('../../../../i18n/language'),
    inject: ['getNode'],
    data() {
      return {
        isCharging:false,
        detail:null,
        IsToolBox:false,
        editMode:true,
        Text:"",
        DivOpacity:1,
        foreColor:"#000000",
        backColor:"#ffffff",
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
        StatusValue:"loading",
        base:{
          text: "configComponent.Slider.title",
          "icon": "icon-huadongshuru",
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
                "w": 200,
                "h": 50
              },
              "visible":1,
              "backColor": "transparent",
              "zIndex": -1,
              "diy":[
                {
                  "name":"configComponent.Slider.Max",
                  "type": 1,
                  "value": 100,
                  "key":"Max",
                },
                {
                  "name":"configComponent.Slider.Min",
                  "type": 1,
                  "value": 0,
                  "key":"Min",
                },
                {
                  "name":"configComponent.Slider.Height",
                  "type": 1,
                  "value": 4,
                  "key":"Height",
                },
                {
                  "name":"configComponent.Slider.HavedSliderColor",
                  "type": 2,
                  "value": "#91d5ff",
                  "key":"HavedSliderColor",
                },
                {
                  "name":"configComponent.Slider.NoHavedSliderColor",
                  "type": 2,
                  "value": "#f5f5f5",
                  "key":"NoHavedSliderColor",
                },
                {
                  "name":"configComponent.Slider.SliderColor",
                  "type": 2,
                  "value": "#91d5ff",
                  "key":"SliderColor",
                },
                {
                  "name":"configComponent.Slider.HoverHavedSliderColor",
                  "type": 2,
                  "value": "#5fa9d7",
                  "key":"HoverHavedSliderColor",
                },
                {
                  "name":"configComponent.Slider.HoverNoHavedSliderColor",
                  "type": 2,
                  "value": "#f5f5f5",
                  "key":"HoverNoHavedSliderColor",
                },
                {
                  "name":"configComponent.Slider.Direction",
                  "type":6,
                  "enumList":[
                    {option:'configComponent.Slider.Directionh',value:1},{option:'configComponent.Slider.Directionv',value:2}
                  ],
                  "value":1,
                  "key":"Direction",
                },
                {
                  "name":"configComponent.Slider.HoverSliderColor",
                  "type": 2,
                  "value": "#5fa9d7",
                  "key":"HoverSliderColor",
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
          "--animateSpinSpeed":this.animateSpinSpeed+'s',
          "--height": this.SliderHeight+'px',
          "--width": this.detail.style.position.w+'px',
          '--HavedSliderColor':this.HavedSliderColor,
          '--NoHavedSliderColor':this.NoHavedSliderColor,
          '--SliderColor':this.SliderColor,
          '--HoverHavedSliderColor':this.HoverHavedSliderColor,
          '--HoverNoHavedSliderColor':this.HoverNoHavedSliderColor,
          '--HoverSliderColor':this.HoverSliderColor,
          "--heightHandle": (this.SliderHeight+10)+'px',
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
      ManualSetData(value){
        let _t = this

        this.settingLoading=true
        let params = {
          deviceUuid:this.detail.active[0].condition.deviceSN,
          dataUuid:this.detail.active[0].condition.dataID,
          value:value.toString(),
        };

        setData(params).then(function (res){
          if(res.data.code!=0)
          {
            _t.$message.error(_t.$t("readData.SetFailed"))
          }
        }).catch(function (error) {
          _t.settingLoading = false
          _t.$message.error(_t.$t("readData.SetFailed"))
        }).finally(function (error) {
          _t.settingLoading = false
        })
      },
      handleChange(value) {
        if((!this.editMode))
        {
          this.isCharging=true
          this.ManualSetData(value)
          setTimeout(function (){
            this.isCharging=false
          },1000)
        }
      },
      initComponents(option){
        if(this.IsToolBox)
        {
          return
        }
        this.backColor=this.detail.style.backColor
        this.foreColor=this.detail.style.foreColor
        let i=0
        for( i=0;i<option.style.diy.length;i++)
        {
          if(option.style.diy[i].key=="Max")
          {
            this.sliderMax=parseInt(option.style.diy[i].value)
          }
          else if(option.style.diy[i].key=="Min")
          {
            this.sliderMin=parseInt(option.style.diy[i].value)
          }
          else if(option.style.diy[i].key=="HavedSliderColor")
          {
            this.HavedSliderColor=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="NoHavedSliderColor")
          {
            this.NoHavedSliderColor=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="SliderColor")
          {
            this.SliderColor=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="HoverHavedSliderColor")
          {
            this.HoverHavedSliderColor=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="HoverNoHavedSliderColor")
          {
            this.HoverNoHavedSliderColor=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="HoverSliderColor")
          {
            this.HoverSliderColor=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="Height")
          {
            this.SliderHeight=parseInt(option.style.diy[i].value)
          }
          else if(option.style.diy[i].key=="Direction")
          {
            this.Direction=parseInt(option.style.diy[i].value)
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
              if(!_t.isCharging) {
                _t.sliderValue = parseFloat(data.result)
              }
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
.ant-slider-rail {
  height: var(--height);
  background-color: var(--NoHavedSliderColor);
}
.ant-slider-handle .ant-tooltip-open {
  border-color: #de2559;
}
.ant-slider-track {
  height: var(--height);
  background-color:var(--HavedSliderColor);
}
 .ant-slider-handle {
   height: var(--heightHandle);
   width: var(--heightHandle);
   border: solid 2px var(--SliderColor);
 }

.ant-slider:hover .ant-slider-track {
  background-color: var(--HoverHavedSliderColor);
}
.ant-slider:hover .ant-slider-rail {
  background-color: var(--HoverNoHavedSliderColor);
}
.ant-slider:hover .ant-slider-handle {
  border: solid 2px var(--HoverSliderColor);
}
</style>
