<template>

  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"   x="0px" y="0px"  viewBox="-8.92 1.199 41.632 29.575"   xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
    <path class="color" d="M-6.999,25.796
	c6.291-7.727,12.581-15.453,18.872-23.179c6.308,7.726,12.615,15.452,18.923,23.179c0.772,0.945,0.632,2.338-0.315,3.109
	c-0.947,0.77-2.34,0.627-3.109-0.32c-5.167-6.346-10.333-12.69-15.498-19.035c-5.16,6.337-10.319,12.674-15.479,19.01
	c-0.759,0.949-2.148,1.098-3.091,0.33S-7.775,26.731-6.999,25.796z"/>

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
    name: 'view-svg-arrows31',
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
       this.initComponents(this.detail);
        if((!_t.editMode)&&(!this.IsToolBox))
        {
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
