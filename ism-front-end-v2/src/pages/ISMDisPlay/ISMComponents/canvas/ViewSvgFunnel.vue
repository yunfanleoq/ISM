<template>

  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"   x="0px" y="0px" viewBox="0 0 43.986 43.984"    xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el"  :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
       <path class="color" d="M26.178,41.639
        c-0.583,1.245-7.79,1.245-8.372,0C12.423,30.142,7.041,18.645,1.658,7.147C0.796,5.306,4.069,4.015,5.265,3.59
        c1.79-0.636,3.616-1.021,5.48-1.328c3.738-0.616,7.465-0.844,11.246-0.844c3.779,0,7.505,0.229,11.243,0.844
        c1.864,0.307,3.688,0.693,5.478,1.327c1.197,0.425,4.475,1.716,3.613,3.558C36.943,18.645,31.56,30.142,26.178,41.639z"/>
            <path class="stroked" d="M42.326,7.147
        c-1.122,2.396-7.84,3.504-9.696,3.787c-3.538,0.539-7.066,0.747-10.638,0.747c-3.57,0-7.096-0.208-10.631-0.746
        c-1.857-0.283-8.582-1.392-9.703-3.788"/>
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
    name: 'view-svg-3d-funnel',
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
          "text": "displayConfig.ToolBox.Diagram.Funnel",
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
