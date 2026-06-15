<template>

  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"   x="0px" y="0px" viewBox="0 0 43.357 29.449"  xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <polygon points="2.607,6.323 4.201,14.87 2.607,23.467 1.417,23.517 3.011,14.87 1.417,6.247"/>
      <polygon points="7.741,6.55 9.336,14.87 7.741,23.24 6.174,23.291 7.741,14.87 6.174,6.5 "/>
      <polygon points="13.281,6.85 14.874,14.87 13.281,22.939 10.726,23.064 12.295,14.87 10.726,6.7 "/>
      <polygon points="19.024,7.127 20.592,14.87 19.024,22.636 15.255,22.838 16.847,14.87 15.255,6.926 "/>
      <path d="M28.307,2.301c2.715,5.821,7.613,10.337,13.634,12.569c-6.021,2.231-10.919,6.747-13.634,12.567
        c-0.048,0.104-0.121,0.113-0.22,0.119c-0.099,0.008-0.207,0.016-0.185-0.119c0.295-1.768,0.591-3.535,0.886-5.303
        c-2.857,0.146-5.715,0.291-8.573,0.439c0.522-2.57,1.043-5.137,1.566-7.704c-0.523-2.563-1.046-5.128-1.569-7.691
        c2.859,0.143,5.719,0.284,8.576,0.428c-0.295-1.769-0.591-3.537-0.886-5.306c-0.022-0.133,0.086-0.125,0.185-0.118
        C28.186,2.188,28.259,2.199,28.307,2.301z"/>
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
    name: 'view-svg-arrows1',
    inject: ['getNode'],
    data() {
      return {
        detail:null,
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
