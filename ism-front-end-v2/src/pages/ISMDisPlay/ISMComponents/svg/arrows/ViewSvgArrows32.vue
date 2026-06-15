<template>

  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"   x="0px" y="0px"  viewBox="0 0 32.21 29.449"   xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
    <polygon class="color" points="1.417,14.871 12.249,1.417 19.664,1.417 8.5,14.871 19.664,28.323 12.249,28.323 1.417,14.871 "/>
      <polygon class="color" points="12.564,14.871 23.378,1.417 30.793,1.417 19.664,14.871 30.793,28.323 23.378,28.323 12.564,14.871 "/>


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
    name: 'view-svg-arrows32',
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
