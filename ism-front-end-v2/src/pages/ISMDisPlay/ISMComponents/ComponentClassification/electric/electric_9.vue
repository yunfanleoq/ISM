<template>

  <svg  xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"   x="0px" y="0px" viewBox="0 0 32 32"   :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">

    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <rect x="1" y="10" width="6" height="12" :stroke="strokeColor" :fill="handcartIn ? closeColor : 'transparent'" :stroke-width="strokeWidth"/>
      <circle cy="16" cx="10" :r="ConnectDiameter" :stroke="strokeColor" :fill="fill" :stroke-width="strokeWidth"/>
      <circle cy="16" cx="30" :r="ConnectDiameter" :stroke="strokeColor" :fill="fill" :stroke-width="strokeWidth"/>
      <line :stroke="strokeColor" :stroke-width="strokeWidth" y1="16" x1="7" y2="16" x2="10"/>
      <line v-if="!isClose" :stroke="strokeColor" :stroke-width="strokeWidth" y1="16" x1="10" y2="6" x2="28"/>
      <line v-else :stroke="closeColor" :stroke-width="strokeWidth" y1="16" x1="10" y2="16" x2="30"/>
      <animate v-if="isStart&&animateType.includes('blink')&&!IsToolBox" attributeName="opacity"
               values="0.1;1;0.1" :dur="blinkSpeed+'s'"
               repeatCount="indefinite"/>
      <animate v-if="isStart&&animateType.includes('millcolorGrad')&&!IsToolBox" attributeName="fill"
               :values="startColor+';'+stopColor+';'+startColor" :dur="animateSpeed+'s'"
               repeatCount="indefinite"/>
      <animateTransform v-if="isStart&&animateType.includes('Zoom')&&!IsToolBox" attributeName="transform"   begin="0s" dur="0.6s" type="scale" values="0.9;1;0.9" repeatCount="indefinite"/>
      <animateTransform v-if="isStart&&animateType.includes('animateSpin')&&!IsToolBox&&spinDirection==0" attributeType="XML" attributeName="transform" :dur="animateSpinSpeed+'s'" type="rotate" from="0 0 0" to="360 0 0" repeatCount="indefinite" />
      <animateTransform v-if="isStart&&animateType.includes('animateSpin')&&!IsToolBox&&spinDirection==1" attributeType="XML" attributeName="transform" :dur="animateSpinSpeed+'s'" type="rotate" from="360 0 0" to="0 0 0" repeatCount="indefinite" />
    </g>
  </svg>

</template>

<script>
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

function asActiveBool(result) {
  return result === true || result === 1 || result === '1' || result === 'true'
}

export default {
  mixins: [ISMChildAutoMixin],
  name: 'view-svg-electric9',
  inject: ['getNode'],
  data() {
    return {
      detail:{},
      IsToolBox:false,
      editMode:true,
      isClose:false,
      handcartIn:false,
      strokeColor:"#00FFFF",
      closeColor:"#d81e06",
      fill:"#00FFFF",
      ConnectDiameter:1,
      strokeWidth:1,
      fillOpacity:1,
      strokeOpacity:1,
      animateType:[],
      startColor:"#74f808",
      stopColor:"#74f808",
      animateSpeed:0.5,
      animateSpinSpeed:0.5,
      spinDirection:0,
      blinkSpeed:0.5,
      isStart:false,
      base:{
        "text": "component.Electric.BreakerSwitch",
        "icon": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAAAXNSR0IArs4c6QAAAOlJREFUWEftlDESwiAQRf9ex0avY6MHMIfwEGasTeMxbC218SBeAF0BIwaBITAp3MykyZD8/x4bCBNfNHE+pMCfGVBqCaLj59zVN6DUDMAKwBrA9Rm+BdHZlqhXgGl16ALA4XUT3b7/urIFhrQc6iivU6CnnQPoftH6zpx8Axm0ZQqMoM0voGl5oHiaeZKje5t6xwaFjag1MYsOBlaJr6UpH0XsFlErTagH+zMgjuAfQ1aT1bjK9CAqB2jNfiuC2sK9BZQNdw2G2xBNdz4h/MPovi3k1ZIATEgBsSAGBADDwqYSvMLer7OAAAAAElFTkSuQmCC",
        "isFontIcon": true,
        "info": {
          "type": "image",
          "action": [],
          "active": [
            {
              id:"HandcartPosition",
              name:"component.Electric.HandcartPosition",
              result:"",
              isExpression:true,
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
            },
            {
              id:"BreakerStatus",
              name:"component.Electric.BreakerStatus",
              result:"",
              isExpression:true,
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
              { id: "Forbidden", name: "component.public.Forbidden" },
              { id: "blink", name: "component.public.animateBlink" },
              { id: "millcolorGrad", name: "component.public.millcolorGrad" },
              { id: "Zoom", name: "component.public.Zoom" },
              { id: "animateSpin", name: "component.public.animateSpin" }
            ],
            "animateElement": [
              {
                id: "blink",
                elementList:[
                  { name:"component.public.animateSpeed", type:7, value:1, min:0.1, key:"blinkSpeed" }
                ]
              },
              {
                id: "millcolorGrad",
                elementList:[
                  { name: "component.public.startColor", type: 2, value: "#74f808", key: "startColor" },
                  { name: "component.public.stopColor", type: 2, value: "#f30b0b", key: "stopColor" },
                  { name:"component.public.animateSpeed", type:7, value:1, min:0.1, key:"animateSpeed" }
                ]
              },
              {
                id: "animateSpin",
                elementList:[
                  { name:"component.public.animateSpinSpeed", type:7, value:1, min:0.1, key:"spinSpeed" },
                  {
                    name:"configComponent.bigScreen.border.border89Direction",
                    type:6,
                    value:0,
                    enumList:[
                      { value:0, option:"configComponent.bigScreen.border.border89DirectionForward" },
                      { value:1, option:"configComponent.bigScreen.border.border89DirectionNegative" }
                    ],
                    min:1,
                    key:"spinDirection",
                  }
                ]
              }
            ]
          },
          "style": {
            "position": { "x": 0, "y": 0, "w": 32, "h": 32 },
            "visible":1,
            "backColor": "transparent",
            "zIndex": -1,
            "transform": 0,
            "diy":[
              { name:"component.Electric.ElectronicDeviceWidth", type:7, value:1, min:0.1, key:"strokeWidth" },
              { name:"component.Electric.ElectronicDeviceColor", type:2, value:"#00FFFF", key:"strokeColor" },
              { name:"component.Electric.CloseColor", type:2, value:"#d81e06", key:"closeColor" },
              { name:"component.Electric.ConnectDiameter", type:1, value:"1", min:1, key:"ConnectDiameter" },
              { name:"component.Electric.ConnectColor", type:2, value:"#00FFFF", key:"strokeFill" },
              { name:"component.public.fillOpacity", type:7, value:1, min:0.1, max:1, key:"fillOpacity" }
            ]
          }
        }
      }
    }
  },
  watch: {
    detail: {
      handler(newVal) {
        if(this.editMode) {
          this.initComponents(newVal);
        }
      },
      deep: true
    }
  },
  methods: {
    initComponents(option){
      if(this.IsToolBox || !option)
      {
        return
      }
      const diy = (option.style && option.style.diy) || []
      for(let i=0;i<diy.length;i++)
      {
        if(diy[i].key=="strokeWidth")
        {
          this.strokeWidth=diy[i].value
        }
        else if(diy[i].key=="strokeFill")
        {
          this.fill=diy[i].value
        }
        else if(diy[i].key=="strokeColor")
        {
          this.strokeColor=diy[i].value
        }
        else if(diy[i].key=="closeColor")
        {
          this.closeColor=diy[i].value || "#d81e06"
        }
        else if(diy[i].key=="fillOpacity")
        {
          this.fillOpacity=diy[i].value
        }
        else if(diy[i].key=="strokeOpacity")
        {
          this.strokeOpacity=diy[i].value
        }
        else if(diy[i].key=="ConnectDiameter")
        {
          this.ConnectDiameter=diy[i].value
        }
      }
      this.animateType = (option.animate && option.animate.selected) || []
      if(option.animate && option.animate.isExpression)
      {
        this.isStart = false
      }
      else
      {
        this.isStart = true
      }
      const animateElement = (option.animate && option.animate.animateElement) || []
      for(let i=0;i<animateElement.length;i++)
      {
        if(animateElement[i].id=="millcolorGrad")
        {
          for(let k =0;k<animateElement[i].elementList.length;k++)
          {
            if(animateElement[i].elementList[k].key=="startColor")
            {
              this.startColor=animateElement[i].elementList[k].value
            }
            else if(animateElement[i].elementList[k].key=="stopColor")
            {
              this.stopColor=animateElement[i].elementList[k].value
            }
            else if(animateElement[i].elementList[k].key=="animateSpeed")
            {
              this.animateSpeed=animateElement[i].elementList[k].value
            }
          }
        }
        else if(animateElement[i].id=="blink")
        {
          for(let k =0;k<animateElement[i].elementList.length;k++) {
            if (animateElement[i].elementList[k].key == "blinkSpeed") {
              this.blinkSpeed = animateElement[i].elementList[k].value
            }
          }
        }
        else if(animateElement[i].id=="animateSpin")
        {
          for(let k =0;k<animateElement[i].elementList.length;k++) {
            if (animateElement[i].elementList[k].key == "spinSpeed") {
              this.animateSpinSpeed = animateElement[i].elementList[k].value
            }
            else if (animateElement[i].elementList[k].key == "spinDirection") {
              this.spinDirection = animateElement[i].elementList[k].value
            }
          }
        }
      }
    }
  },
  mounted() {
    let _t = this
    this.$nextTick(function(){
      _t.initComponents(_t.detail);
        let activeEvent = this.detail.identifier+"activeEvent"
        let animateEvent = this.detail.identifier+"animateEvent"

        _t.$EventBus.$on(activeEvent, (data) => {
          if(data.ID == "GroupStrokeColor")
          {
            _t.strokeColor = data.result
            return
          }
          if(data.ID == "BreakerStatus" || data.ID == "ElectronicClose")
          {
            _t.isClose = asActiveBool(data.result)
          }
          else if(data.ID == "HandcartPosition")
          {
            _t.handcartIn = asActiveBool(data.result)
          }
        })
        _t.$EventBus.$on(animateEvent, (data) => {
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
      _t.initComponents(_t.detail);
    })
  }
}
</script>
<style >
.svg-el {
  transform-origin: center center;
}
</style>
