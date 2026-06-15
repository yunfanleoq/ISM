<template>

  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"   x="0px" y="0px" viewBox="0 -0.277 42.563 42.703"    xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el"  :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <polygon class="color" points="1.417,30.909 1.417,21.005 21.279,1.14 31.243,1.14 41.145,11.042 41.145,21.005 21.279,40.871 11.379,40.871 1.417,30.909 	"/>

      <polyline class="stroked" points="11.379,11.042 21.279,11.042 31.243,20.946 31.243,30.909 	"/>
      <line class="stroked" x1="21.163" y1="11.042" x2="31.065" y2="1.14"/>
      <line class="stroked" x1="31.065" y1="20.946" x2="41.027" y2="11.042"/>
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
    name: 'view-svg-3d-octagon1',
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
          "text": "displayConfig.ToolBox.Diagram.Octagon1",
          "icon": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAAAXNSR0IArs4c6QAABCZJREFUWEe9l2tsU2UYgN8mBwYM1B+ToezSgc6wJo6uiZHuHMPU0Qo7W8QVpK3G/SAzdDMOiBpGINEfJAZNlBWzSywWd6sda0xwvSxYaDsjrE7GZpGk3WAz2ZAZDLqtF1bzHTizl3Par2XZ96vpefs9z3s53zkVwKOtXAB4MsEWawDAmQgheAR+LiVTnCosfqE8a/3Ghdh9rg//TExNjP0z/cfNuvEbV7v5OOkK5FJyRVOFSlO5IWdT3N6Djj5Yk/kYFJWUwoDtnK+/V3fE6+GWSEcgj5QrmmiVho6Fh8NhcDvNi3DW7IHEN41ez1BXrG2qAinDWeBP/b0+m+lMo3c0WiIVgeTwtY9DkVjKO1YD/b1jNpPuqG/01w42CFcgj5RVayvV9RXZOQVRAKbsqOfrnkgIX2xHf++4zaRv9I26GQkcgQJKVn2qQqXZFQ9fgF+cVmz4/+0wjdvO6Y57PUP6ZAL5b7575EIoGNy0MmN1XGnv3Z2BjNWZQKvrUrqbgwE/aD8+cHfIZatKJJBPyRVaWl2/K3ujkBdw7YodQsEQiKWvYkkguNthhr//uj3bcfoTik8AwU/T6vqdieAs8drlHyEUup9UIuCfZ27TrOwcCPjn7n16WLWdS0D4MPN4eJh/aoavXIT7wQCIpeWclQj455jMs57Kg2dFEhh1OzgFhJR8j5ZW12FlHksaGbyEMoOSUlnUJf/8LANf/3Q+PCOSMNe4BBh45Vv1O1FgohUOAwh4mjcy6AD//BxIyB3MFlxwLgEhKVN8VfX2e/JkcJxJG3E7wD83CyIJCW6HBTbkCGFzUUnUTyMrUEiW7/6y6p0GWapwgUAA6CDiWkMuGwxftgMpr4bNW8RxIfbznf/qTn74Eipkyf6PPrtEyhSZONnhxiAx05nP4fWaQ3E/mRz7Hfq6mrucVmMDI/DByXa7SEKtw90cJw4JWI1tIFPsjwqPhAPA1BIKRN+jXAIR8PcBYJp9FixLBRh4d3On02JsYOHLJjDpuw7nDc2dA5aeKPiyCIgkFPxgaO5wWXoORmbODgYzAwdP6OzFL76ypEO4EA7D9/ov4M7UrXaH2Yjgt7mGFwlQr+2tbd1dc+i5lRmrcAYcK2bC60HnwFVDywl0JHLC2RYAQRDk3tqjTdtpZfFSSEz4PNDX3drusvJnHtkC5jNBEKV7ahubymjlVq6XD6y0AeABvOVbl7UHnUC8mccJPPxCqtQc16YrgeBmQ8tZp4WB/4kjzfVMkyo1x7RltCqlStzy/gaW71pTgi/OAIfptn0HjmlfrlSJcdrBwA1teqfVeBg3c74WRLowEmW0UpyxCv3H5F4Ibja06V1pwBNVgKVt2/FGzdcFhVvXEitWxD13Z6YniYmbN6wuM5P5HZyex8Ykey1H8c+jF5sEm8+kC0d7/ge9LvHk95aDfgAAAABJRU5ErkJggg==",
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
          this.initComponents(newVal);
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
