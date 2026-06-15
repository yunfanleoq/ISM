<template>

  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"   x="0px" y="0px"  viewBox="0 0 43.358 29.449"    xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <polygon points="1.417,14.871 12.249,1.417 19.664,1.417 8.5,14.871 19.664,28.323 12.249,28.323 1.417,14.871 	"/>
      <polygon points="12.564,14.871 23.378,1.417 30.793,1.417 19.664,14.871 30.793,28.323 23.378,28.323 12.564,14.871 	"/>
      <polygon points="23.693,14.871 34.542,1.417 41.941,1.417 30.793,14.871 41.941,28.323 34.542,28.323 23.693,14.871 	"/>
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
    name: 'view-svg-arrows37',
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
