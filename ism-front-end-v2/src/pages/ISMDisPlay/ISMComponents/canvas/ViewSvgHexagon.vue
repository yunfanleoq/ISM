<template>

  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"   x="0px" y="0px" viewBox="0 2.42 44.907 40.005"    xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el"  :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <polygon class="color" points="1.296,26.252 10.242,11.331 17.719,3.837 34.611,3.837 43.49,18.758 34.611,33.679 27.134,41.109 10.242,41.109 1.296,26.252 	"/>
      <polyline class="stroked" points="10.242,11.331 27.134,11.331 34.611,3.837 	"/>
      <polyline class="stroked" points="27.134,11.331 36.011,26.252 43.49,18.758 	"/>
      <line class="stroked" x1="27.134" y1="41.174" x2="36.011" y2="26.252"/>
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
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
    name: 'view-svg-3d-hexagon1',
  inject: ['getNode'],
    data() {
      return {
        detail:{},
        IsToolBox:false,
        editMode:true,
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
        base:{
          "text": "displayConfig.ToolBox.Diagram.Hexagon1",
          "icon": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACwAAAAoCAYAAACFFRgXAAAAAXNSR0IArs4c6QAABoxJREFUWEfFmXdUU1ccx7901yMFcVtAREHCTFSUoMUJqKgoUVHkKERrq6JVhp7jUcaxQhGt1YKDhKEVUaZYg+xRkSlbJKwyxEXFihQnSM+LDQZI8l4w0Pcn+d7f73O/7/7u/d2HAuTzfAWgWw6h2sliKJAJKPyurjdjztkvhymOo6AVK1FWGaMCdCunXDlvCqBCWpwPBVan0U05NmwXC20DY5l4icTEK7nfWIvkmBBoGUxH+rXw0OrSPMfBAlbXoTM5LLYrdVgFBaD7/cohYFNiQzFv2XqoT9HDtYsBTZEc380AUiRBD9RhNR06k0sKK7RRTHYh7Fyr9ZiopSdQPGm5j0Af56zKkmwrAM/EQQ8EWE2XYcpZ6eBsOdVwppS3RzgpPvyDpjrBMhCFFQYq+CP+VXSQ394HTXUn5QGspstgclY6uJDASp6HADY2BHOXvndWVN3V1QmOz56ynNS4tQCq+kaSxWE1HQYz0MbBZXEvZyUb2Y+agE2JDYGZBFjhgNqKQlw+482tvl3w7UCB6TS66eVpcyy0FZVUJL1pqbvE09ZHqC4rgPWm3dDQ0ifdUa6c/6UhNuRnYsfIEBVTclhD29CTRjfxoJsuIk0kSdDR3oZbmTzYOXlBUWkEaZzHD5vB8XXJ4JfkLAPQIRxABZix0Hpj2MbdP9JIs5AIirNT0N7WCrMltpRC5aZdfR7JPeL8+EHTWcrAOnRm5AYnj9Xqk3UpJSETxYQcw8y5VlDV1CGTovPNa5z1/qEkP4O3BkAtMYDM4WX2u7y45qscx5JGpyiorypDeUEmVtjv7DdCQUEB3SIHCyGoLi9AZKD3qerbhTvIgMfSGEzeDo/T0wWFJscnLe48xnw9CfozvukPLTiye289McFH6+J+O0kU4A2JDmto63uw2G6ehrPmyxH1Xai/Hz/EtbAAsDa7YdhwotET84gc4y33GxHsty+5siR7hSRgxgJr+7BNu70/uNAkzbbwZiJedvyD2RYsqYYQS4RYKkU3kzounDjgJBaYRmdG2Dl5rJFXoYkj6u5+i+igozA1t8GEiVNI3yK/NBeZvPC74oCt7Hd6BZnbyK/QJNHUVRaDX5wLK7ttUoEJ2MSooOdFWYn2fYHfFZr7qemKyiNJZy0PAdFeTtDQhi6D6N37PwRsTmocGvgloQ01FY69gDW0DdxZbFevwSg0SZNrbbmHhAgOWGxXfDFseC8ZAVuWm4am+uo75Xlp1sReLArMWLDC/sKmPd7yOSFksD8/k4euN6/BXLSqZ1RVWR7K8tIxYvT4F9cjOHuEp10P8FAUmqQ5dHV2IjrID2ZWthinqgkh7LTZloji+KbfKcleLuwnhMBW9k6eXHMWe8AXSRkMFSutuX0LRFupSaOjNC8d5qsckMm71BAb2rtjI4DHausbJ+w6FEgfqkKTNLkAz+346JOPYfvdfsF16XKgN6e6rGCrqF7g8FTDWZH2Oz1XExfB/+shGp0wf0/BOp6sywDXx7k0O/UK0db1unUIlwR9/vINYQ7OPkNecEKD8jN46Op6A+bClSjI5L2K5B51e9Rc92tfA3uKbqKW3gEW2+2QkcmCITeZaNaTooMEvUXHs6dSb86i29poHTozfrt7wAylEaOGFJq4QatOmgoawxS/h/k3RXGPsAGkiu2J+vxx6QYnd64Fa8v4oSImdoaq8nxYrduGP/klCD9zWOrXn369hI7RrHC7HZ7rhB83BhP87dsuwf47x3INxqlpgvuTS3lWUtQ6AHck5RXX/BjNW7Y+zNHFd9C3jMIbCXj54jlmW9ggP53XdfXcib13G/k/SzNJbHupNkV3/2q222E6c+GgGfzkrweIv3RG0EM011chNuR4QkXRDQcAj2QGBjBKUIAH/Y2VVEYPCnRq3HnBMfzpZ58jOvhYEr8kmzggGsmSSbuELl637UDQkrVbJ5AFkfX3en4pKgqzoGVgjJjQY4n84hwCtolKHKm3Zh1Dk7D1Ow7aaWgbUIlFWRMdfBRjVSfhRkJEwn+wd6kOJrvmG5gttb242c2P/NsSxYzF2cmorShCXWXR9cp3zjZTHCqQkQFDXZO2z2aLqw+DaU6qJUv8tLUF547vR3vbk+Sa27eIw0EmWErAAFRoRibx3x/0n6U8cgwZk9Tfw08fQnluRv69phobAPcGEoyqa5b0mfMDR01Q7/koJ2uylvuNwxvr7rxqa22ZS/xrQ9bxQj1VYEKvNtAkIuOIJqX4Q+L8C30FbPI28Qz/AAAAAElFTkSuQmCC",
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
      this.initComponents(this.detail)
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
      _t.initComponents(_t.detail);
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
