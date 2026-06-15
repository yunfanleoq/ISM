<template>

  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"   x="0px" y="0px" viewBox="0 0 51.033 43.984"    xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el"  :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
     <polygon class="color" points="13.567,1.418 37.575,1.418 49.616,21.942 37.575,42.567 13.456,42.567 1.416,21.942 13.567,1.418 	"/>
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
    name: 'view-svg-3d-hexagon3',
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
          "text": "displayConfig.ToolBox.Diagram.Hexagon3",
          "icon": "data:image/png;base64,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",
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
