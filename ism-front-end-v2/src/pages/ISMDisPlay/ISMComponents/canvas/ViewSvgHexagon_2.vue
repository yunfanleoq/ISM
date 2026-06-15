<template>

  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"   x="0px" y="0px" viewBox="0 2.42 44.907 40.005"    xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el"  :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <polygon class="color" points="1.416,26.193 8.396,8.909 13.582,0 24.289,0 29.408,8.909 35.992,26.193 27.146,41.15 10.259,41.15 1.416,26.193 "/>
      <polyline class="stroked" points="1.416,26.193 10.259,11.235 13.582,0 "/>
      <polyline class="stroked" points="35.992,26.193 27.146,11.235 24.289,0 "/>
      <line class="stroked" x1="10.259" y1="11.235" x2="27.146" y2="11.235"/>
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
    name: 'view-svg-3d-hexagon2',
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
          "text": "displayConfig.ToolBox.Diagram.Hexagon2",
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
