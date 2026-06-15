<template>

  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"  x="0px" y="0px"  xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el" :class="{'animated':true,[`${detail.style.animate}`]: true}" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">

      <image preserveAspectRatio="none meet" :style="{cursor:editMode?'':'pointer'}" v-if="isOn" :width="detail.style.position.w" :height="detail.style.position.h" :href="SwitchOnImg"></image>
      <image preserveAspectRatio="none meet" :style="{cursor:editMode?'':'pointer'}"  v-else :width="detail.style.position.w" :height="detail.style.position.h" :href="SwitchOffImg"></image>
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
import svgView from '../View';
import {setData} from "@/services/device";
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
    name: 'view-svg-switch',
    i18n: require('../../../../i18n/language'),
    inject: ['getNode'],
    data() {
      return {
        detail:null,
        IsToolBox:false,
        editMode:true,
        DivOpacity:1,
        Text:"",
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
        isOn:0,
        SwitchOn:0,
        SwitchOff:1,
        SwitchOffImg:"",
        tempImgIndex:0,
        SwitchOnImg:"",
        switchStyleIndex:0,
        switchStyle:[
          {
            "on":"/static/images/switch/4.svg",
            "off":"/static/images/switch/3.svg",
          },
          {
            "on":"/static/images/switch/1.png",
            "off":"/static/images/switch/2.png",
          },
          {
            "on":"/static/images/switch/5.png",
            "off":"/static/images/switch/6.png",
          },
          {
            "on":"/static/images/switch/7.png",
            "off":"/static/images/switch/8.png",
          },
          {
            "on":"/static/images/switch/10.svg",
            "off":"/static/images/switch/19.svg",
          },
          {
            "on":"/static/images/switch/12.svg",
            "off":"/static/images/switch/11.svg",
          },
          {
            "on":"/static/images/switch/14.svg",
            "off":"/static/images/switch/13.svg",
          },
          {
            "on":"/static/images/switch/15.png",
            "off":"/static/images/switch/16.png",
          },
          {
            "on":"/static/images/switch/17.svg",
            "off":"/static/images/switch/18.svg",
          },
          {
            "on":"/static/images/switch/9.svg",
            "off":"/static/images/switch/20.svg",
          },
          {
            "on":"/static/images/switch/21.png",
            "off":"/static/images/switch/22.png",
          },
          {
            "on":"/static/images/switch/1_1.png",
            "off":"/static/images/switch/1_0.png",
          },
          {
            "on":"/static/images/switch/3_1.png",
            "off":"/static/images/switch/3_0.png",
          },
          {
            "on":"/static/images/switch/9_1.png",
            "off":"/static/images/switch/9_0.png",
          },
          {
            "on":"/static/images/switch/4_1.png",
            "off":"/static/images/switch/4_0.png",
          },
          {
            "on":"/static/images/switch/58_1.svg",
            "off":"/static/images/switch/58_0.svg",
          },
          {
            "on":"/static/images/switch/23.png",
            "off":"/static/images/switch/24.png",
          }
        ],
        base:{
          text: "configComponent.switch.Text",
          "icon": "icon-STSkaiguan",
          "isFontIcon": true,
          "info": {
            "type": "image",
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
                id:"ControlSwitch",
                name:"component.switch.ControlSwitch",
                result:0,
                isSwitch:true,
                isExpression:false,
                condition:{
                  deviceSN:"",
                  isBandDevice:false,
                  bandType:1,
                  dataID: "",
                  dataName: "",
                  SetPassword:"",
                  IsManual:false,
                  actionAuth:[],
                  AutoSet:[
                    {
                      "name":"configComponent.switch.Off",
                      "value":0
                    },
                    {
                      "name":"configComponent.switch.On",
                      "value":1
                    },
                  ]
                },
              }
            ],
            "style": {
              "position": {
                "x": 0,
                "y": 0,
                "w": 90,
                "h": 170
              },
              "visible":1,
              "zIndex": -1,
              "transform": 0,
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
                  "name":"component.switch.OnImg",
                  "type":5,
                  "value":"/static/images/switch/4.svg",
                  "key":"OnImg",
                },
                {
                  "name":"component.switch.OffImg",
                  "type":5,
                  "value":"/static/images/switch/3.svg",
                  "key":"OffImg",
                },
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
          if(option.style.diy[i].key=="On")
          {
            this.SwitchOn=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="Off")
          {
            this.SwitchOff=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="SwitchStyle")
          {
            this.switchStyleIndex=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="OnImg")
          {
            this.SwitchOnImg = option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="OffImg")
          {
            this.SwitchOffImg = option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="strokeWidth")
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
            if(data.ID == "ControlSwitch")
            {
                if(data.result==_t.detail.active[0].condition.AutoSet[0].value)
                {
                  _t.isOn = 0
                }
                else if(data.result==_t.detail.active[0].condition.AutoSet[1].value)
                {
                  _t.isOn = 1
                }
                else
                {
                  _t.isOn = 0
                }
              _t.detail.active[0].result = data.result
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
<style >
.svg-el {
  /*transform: rotate(45deg);*/
  transform-origin: center center;
}

</style>
