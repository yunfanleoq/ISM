<template>

  <svg  xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"   x="0px" y="0px" viewBox="0 0 64 64"   :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">

    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <circle cx="32" cy="3"  :r="ConnectDiameter" stroke-width="0.0"/>
      <line :fill="fill" :stroke="strokeColor" :stroke-width="strokeWidth" x1="32" y1="10" x2="32" y2="3"/>
      <rect :fill="fill" :stroke="strokeColor" :stroke-width="strokeWidth" x="43" y="18" width="6" height="18"/>
      <line :fill="fill" :stroke="strokeColor" :stroke-width="strokeWidth" x1="43" y1="20" x2="48" y2="20"/>
      <line :fill="fill" :stroke="strokeColor" :stroke-width="strokeWidth" x1="43" y1="22" x2="48" y2="22"/>
      <line :fill="fill" :stroke="strokeColor" :stroke-width="strokeWidth" x1="43" y1="32" x2="49" y2="32"/>
      <line :fill="fill" :stroke="strokeColor" :stroke-width="strokeWidth" x1="43" y1="34" x2="48" y2="34"/>
      <line :fill="fill" :stroke="strokeColor" :stroke-width="strokeWidth" x1="18" y1="22" x2="13" y2="29"/>
      <line :fill="fill" :stroke="strokeColor" :stroke-width="strokeWidth" x1="32" y1="41" x2="32" y2="47"/>
      <line :fill="fill" :stroke="strokeColor" :stroke-width="strokeWidth" x1="28" y1="47" x2="36" y2="47"/>
      <line :fill="fill" :stroke="strokeColor" :stroke-width="strokeWidth" x1="29" y1="49" x2="35" y2="49"/>
      <line :fill="fill" :stroke="strokeColor" :stroke-width="strokeWidth" x1="31" y1="51" x2="33" y2="51"/>
      <line :fill="fill" :stroke="strokeColor" :stroke-width="strokeWidth" x1="18" y1="10" x2="46" y2="10"/>
      <line :fill="fill" :stroke="strokeColor" :stroke-width="strokeWidth" x1="46" y1="10" x2="46" y2="18"/>
      <line :fill="fill" :stroke="strokeColor" :stroke-width="strokeWidth" x1="18" y1="10" x2="18" y2="22"/>
      <line :fill="fill" :stroke="strokeColor" :stroke-width="strokeWidth" x1="18" y1="30" x2="18" y2="41"/>
      <line :fill="fill" :stroke="strokeColor" :stroke-width="strokeWidth" x1="18" y1="41" x2="46" y2="41"/>
      <line :fill="fill" :stroke="strokeColor" :stroke-width="strokeWidth" x1="46" y1="41" x2="46" y2="36"/>

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
import svgView from '../../View';
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
  name: 'view-svg-electric7',
  inject: ['getNode'],
  data() {
    return {
      detail:{},
      IsToolBox:false,
      editMode:true,
      strokeColor:"#00FFFF",
      fill:"#00FFFF",
      ConnectDiameter:1,
      strokeWidth:1,
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
      base:{
        "text": "configComponent.image.Text",
        "icon": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAEAAAABACAYAAACqaXHeAAAAAXNSR0IArs4c6QAAAipJREFUeF7tmj1OxDAQhd9cgRJBCTVIlBT0aA8AHIEDULKU9CBOwHIARE9BiQQlPyWIkisMOE7Y7LKBaOLYjj2utsjYnm+ex2OvCZk3ytx/hAXA/FIEgGg9VCBCA3gvAazkCuC4BHCiAAIRCL0EVAG6BOwuMPAcwDwGsAbAbmvt20756W17k+JLs22+gsiM26m5yQHMlwCWAUgcMQ5IwH2A6KCT90Z8XTso7Jn9JjOH4ykAVYALAg4l2Wo6DsfTJdCK+H8fOYzIf0O5TrrDUkAFGpitHzoUUsMBYJ1vKnzG0mpySABMsbXfsEQUQGgFTAA8SyfRMvFFqgDmJQBPAC5AZEviPpo9b0S4BJhN9N9AdNSH3z99RpkEmQ+/pb8Hou1enZ8euCLaBZg3ADwA2ATRY44A7gBcgei8d+ejUwDzKYBVEDUlJfdMoskBzCMAZ6X0P9172tBjNACsHEcguq5l6N3iN9FNb0CiAjDvpY/ToAL4sxCaSC9IUzkMZQ9AT4PSg9iQlkDmFyLTatD8yvBKrL71Otx2h7MEFECNQIQKyP7f4czfB0hPQA6lLJ2CmyQoHV0BeH5YsSBQqgCpekV2zObiZKtmu+iR1H2vFytzE/ergN8AqkfS9UdSCQOw9fy8Cuox8ep8cYsnknIXo+wBVPBmQXiPfDUN/wpQAF3Wj3vbcApw74uoRwUgwpaQkSogoWCKXFEFiLAlZKQKSCiYIldUASJsCRl9ARc1HlDz91x5AAAAAElFTkSuQmCC",
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
              "w": 32,
              "h": 32
            },
            "visible":1,
            "backColor": "transparent",
            "zIndex": -1,
            "transform": 0,
            "diy":[
              {
                "name":"component.Electric.ElectronicDeviceWidth",
                "type":7,
                "value":1,
                "min":0.1,
                "key":"strokeWidth",
              },
              {
                "name":"component.Electric.ElectronicDeviceColor",
                "type":2,
                "value":"#00FFFF",
                "key":"strokeColor",
              },
              {
                "name":"component.Electric.ConnectDiameter",
                "type":1,
                "value":"1",
                "min":1,
                "key":"ConnectDiameter",
              },
              {
                "name":"component.Electric.ConnectColor",
                "type":2,
                "value":"#00FFFF",
                "key":"strokeFill",
              },
              {
                "name":"component.public.fillOpacity",
                "type":7,
                "value":1,
                "min":0.1,
                "max":1,
                "key":"fillOpacity",
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
        else if(option.style.diy[i].key=="ConnectDiameter")
        {
          this.ConnectDiameter=option.style.diy[i].value
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
	  _t.initComponents(_t.detail)
    })
  }
}
</script>
<style >
.svg-el {
  /*transform: rotate(45deg);*/
  transform-origin: center center;
}
</style>
