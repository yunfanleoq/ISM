<template>

  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"   x="0px" y="0px"  viewBox="0 0 18.945 29.449"    xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <path class="color" d="M10.2,25.702
	c0,0.235,0.016,0.477,0.048,0.725c0.079,0.631,0.255,1.309,0.69,1.599c0.437,0.289,1.134,0.192,1.736,0.004
	c0.605-0.188,1.117-0.467,1.633-0.763c0.515-0.296,1.034-0.608,1.559-0.933c0.525-0.323,1.055-0.657,1.661-0.657
	c-0.134,0-0.2-0.261-0.227-0.342c-0.114-0.358-0.166-0.725-0.215-1.094c-0.068-0.512-0.111-1.024-0.146-1.538
	c-0.068-1.03-0.095-2.06-0.095-3.091c0-1.032,0.026-2.062,0.095-3.091c0.033-0.514,0.078-1.026,0.146-1.537
	c0.049-0.37,0.1-0.736,0.215-1.096c0.026-0.081,0.093-0.342,0.227-0.342c-0.606,0-1.136,0.334-1.661,0.657
	c-0.524,0.323-1.044,0.637-1.559,0.933c-0.516,0.296-1.027,0.573-1.633,0.762c-0.306,0.096-0.638,0.168-0.945,0.179
	c0-2.562,0-5.125,0-7.687c0-0.356-0.067-0.703-0.23-0.994c-0.283-0.499-0.854-0.831-1.436-0.89c-0.583-0.058-1.18,0.16-1.756,0.373
	C7.732,7.094,7.177,7.302,6.597,7.4c-0.581,0.097-1.184,0.082-1.79,0.082c0-2.021,0-4.043,0-6.064
	c-1.13,4.043-2.26,8.086-3.39,12.13c1.13,4.043,2.26,8.086,3.39,12.13c0-2.022,0-4.044,0-6.064c0.606,0,1.209,0.015,1.79-0.083
	c0.579-0.097,1.134-0.306,1.71-0.52c0.576-0.213,1.173-0.43,1.756-0.373c0.046,0.006,0.092,0.012,0.137,0.02
	C10.2,21.005,10.2,23.353,10.2,25.702z"/>
      <path class="stroked" d="M11.73,16.077
	c-0.299,0.009-0.575-0.039-0.791-0.182c-0.436-0.29-0.611-0.968-0.69-1.599c-0.08-0.63-0.063-1.215,0.071-1.788
	c0.135-0.573,0.386-1.135,0.646-1.68c0.262-0.546,0.533-1.074,0.672-1.671c0.058-0.25,0.092-0.511,0.092-0.767"/>
      <line class="stroked" x1="10.2" y1="13.572" x2="10.2" y2="18.656"/>
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
    name: 'view-svg-arrows38',
    inject: ['getNode'],
    data() {
      return {
        detail:null,IsToolBox:false,editMode:true,strokeColor:"#000000",
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
        base:{
          "text": "configComponent.image.Text",
          "icon": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAAAXNSR0IArs4c6QAAA4hJREFUWEdjZBhgwEjAfgcGBgYFIF5AK3fidYCQqOT1Xz++X/ny+UPoQDgggYGBYT4Pr8CaAXEAD7/geRlFDYMnd68PiANA8X4/pbyXYcW05gFxwPzAhMIEYXEZQg4AORSEQeAAOekEVyL837v8GMP1C8cZFk2sOfHrx/dKkOFSiqomP758kRYWl3b49P6NgL65IwMbB6fAmxdPWS6fPvDn6+ePExgYGBpJcQg2ByRYOPlPzqydzPPmxROGS6f2vwEZyMrGISIqIcOgYWCJ1XyQ2vm9FXeunDmUSkpoYDiAk4vndUHrXBFcFuHzXbqX5p8f37/2MjAwVBAbCugOSFBQ0+1vnLlVgFgDYOpAIVCf7v3/y6f3SaQUXCgOAGW93IYZBuT4vjjSisEvOvfvvN7yFHIdACp29y/c/4hUzzNsWNgP1mNk5fanNs0TlAaILrqRQ2B+Snlvgq0HaaUuKOjbC8MYQLnm0e2rFDngPzm+P7prLcOiibUMM7deI8YBCkKikvmg0Hr3+nkhiIaHALnxD3LA/N4Khjk7b+NyANhSVjb2AFN7LwFWNnaB3Wvnw0tX5CggKwdcPn2QYWJNCkPHwv0M3z5/hEXBAZilrGzsHB5haRKwqJ3enPvlxL6NubB0guwABU4untOklgEgB2xaPIlB29iG4caF43/evX7xBFRCGtt6KGDLTbmBhl8+fXiry8DA8AAlCqBFbYmyukE3qBIiFoAcsH3lTIaynmVEacHrAFDFws3Lf37apstEF0SkOiDeUQ7F4xhFMQ+vwOqIrNoQ9OwIym4gsH/zkg8g+ualUw9AFdLfP3++/vjxVXXqhotsxAQBQQeAQoFfWOz6pDVnOG5cOM4woTr5DQ+/4BluXqEXIAse3Lq4EK2ycVBQ011PbPFNjAMYYFny9YsnhNoDIDdR3wEwQ10CEwSIaBHRxAGgwmO7trGNxvmjuwk1yUhyAKFcAE9HUoqqJT++fGn/9eP7BgKtYpIcUJ/u/eHBrcuBsHSEr18AbpgS0SwnyQGrZ3d82LJsGrzphrdjAsqSoCChZgiAqu7d6xdc+PLxvSFGSYglH4NCAVxk4gEEQwBUhoDKj/2blzJAG66grIxZFBNTkGBzJCjBpldN0EAu+2GWnj647cPLpw9ADRSsrWVCnVNi3aQgp6y1+N+/vyqc3LwvQCUk1FK4T3EZRC0HEOtQDHUA9MS1MF1yvEUAAAAASUVORK5CYII=",
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
                  "name":"component.public.strokeWidth",
                  "type":7,
                  "value":0.6,
                  "min":0,
                  "key":"strokeWidth",
                },
                {
                  "name":"component.public.strokeFill",
                  "type":2,
                  "value":"#A1BFE2",
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
                  "name":"component.public.strokeOpacity",
                  "type":7,
                  "value":1,
                  "min":0,
                  "max":1,
                  "key":"strokeOpacity",
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
       if(true){
        let activeEvent = this.detail.identifier+"activeEvent"//动作数据
        let animateEvent = this.detail.identifier+"animateEvent"//动作数据

        _t.$EventBus.$on(activeEvent, (data) => {


        })
        _t.$EventBus.$on(animateEvent, (data) => {
          _t.isStart = data
        })
      }
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
