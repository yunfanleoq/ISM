<template>

  <svg  xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none meet"   x="0px" y="0px" viewBox="0 0 48 48"   :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">

    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <circle cx="33" cy="4"  :r="ConnectDiameter" stroke-width="0.0"/>
      <rect fill="none" :stroke="strokeColor" :stroke-width="strokeWidth" x="9" y="13" width="8" height="19"/>
      <line fill="none" :stroke="strokeColor" :stroke-width="strokeWidth" x1="9" y1="15" x2="17" y2="15"/>
      <line fill="none" :stroke="strokeColor" :stroke-width="strokeWidth" x1="9" y1="17" x2="17" y2="17"/>
      <line fill="none" :stroke="strokeColor" :stroke-width="strokeWidth" x1="9" y1="28" x2="17" y2="28"/>
      <line fill="none" :stroke="strokeColor" :stroke-width="strokeWidth" x1="9" y1="30" x2="17" y2="30"/>
      <line fill="none" :stroke="strokeColor" :stroke-width="strokeWidth" x1="13" y1="32" x2="13" y2="36"/>
      <line fill="none" :stroke="strokeColor" :stroke-width="strokeWidth" x1="9" y1="36" x2="17" y2="36"/>
      <line fill="none" :stroke="strokeColor" :stroke-width="strokeWidth" x1="16" y1="38" x2="10" y2="38"/>
      <line fill="none" :stroke="strokeColor" :stroke-width="strokeWidth" x1="12" y1="40" x2="14" y2="40"/>
      <circle fill="none" :stroke="strokeColor" :stroke-width="strokeWidth" cx="33" cy="41" r="6"/>
      <circle fill="none" :stroke="strokeColor" :stroke-width="strokeWidth" cx="33" cy="31" r="6"/>
      <circle fill="none" :stroke="strokeColor" :stroke-width="strokeWidth" cx="25" cy="36" r="6"/>
      <line fill="none" :stroke="strokeColor" :stroke-width="strokeWidth" x1="33" y1="25" x2="33" y2="4"/>
      <polyline fill="none" :stroke="strokeColor" :stroke-width="strokeWidth" points="13,13 13,10 33,10 "/>
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
    name: 'view-svg-electric2',
  inject: ['getNode'],
    data() {
      return {
        detail:{},
        IsToolBox:false,
        editMode:true,
        strokeColor:"#00FFFF",
        fill:"#00FFFF",
        ConnectDiameter:1,
        strokeWidth:1,
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
          "icon": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADAAAAAwCAYAAABXAvmHAAAAAXNSR0IArs4c6QAAAoBJREFUaEPtmC1sFUEQx38jqUNgCIo0kCAIaUDgaHAlKHAkoNF1IABBHRoNCQ5U01pwCBoEggQSBAnBIHBFTrtvb7l7d+92993tfbzmNrm8l+x9zG9m/rszK6z4kBW3nwmgdQRVv8/eIXKhybuGj4Dqrwzg3KoCPMkAnk0ATTzQ+hnVKQLLO9F5bfknFz1xA/iASE8asMbbj6YZgwCYdbuRxyrMvWhgPmWMx8zII9AGpnOAcMq0SgF6AHgDnPXkvAH4jci9RpKYAEJuU50i4PVRDynk3+pbGtCHiCeAKYU8Hgh3ZKEcD82HV7mOy+mQgaH5Ewugegq4BTzMGF8Ce4j8CzEX54dJIdVHgLneA4eZQWvAJrCDyE4sRP8Aqq+B08A2IvZIxQ1Vc7TyAviLyP0YiBiAdKWE9fx1RG4HltZd4GNMJPoDsDn/B9ioeL5MYyPxGTgT0kR3ALlIL2b2mbS5hMhWTGqgaqLwCpG3vvu7AZgX6ZfMgDvA+eP/T2NSA9XnM4GLmN/aEQNQbuLdGaYT4HxHVidS1ceAOT40V1ikyQAMe2xP7BOp6l3gwUzANj38Ik2WQlWBLd76QyItztt31os0qYjjAXIP12VsMUKq+8DXWTrZ8e3/ThwToeyhGA2Y7f5qwabqsQocAFeANURMrtcPq5FrwDrwA3iX3XwZuAn8PC4xPqXcyMoAZRGb748YwIq4DFH08AEiexRFOqoUigfw77SjEPF8JKzni2P0y2gIwEZrcbU56EbmXVYWTI6qlFjWeHf/aIq5pgDVDXCgcjoVgNXHAA1NSgCfyO1cBy1laoA8Eq6pd/2CKSU6aOq7ALAQ7ljFdWx5MbfEN48A4aW6QHId0doAAAAASUVORK5CYII=",
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
                  "name":"component.Electric.ElectronicDeviceWidth",
                  "type":7,
                  "value":1,
                  "min":0.1,
                  "key":"strokeWidth",
                },
                {
                  "name":"component.Electric.ElectronicDeviceColor",
                  "type":2,
                  "value":"#00FFFF",
                  "key":"strokeColor",
                },
                {
                  "name":"component.Electric.ConnectDiameter",
                  "type":1,
                  "value":"1",
                  "min":1,
                  "key":"ConnectDiameter",
                },
                {
                  "name":"component.Electric.ConnectColor",
                  "type":2,
                  "value":"#00FFFF",
                  "key":"strokeFill",
                },
                {
                  "name":"component.public.fillOpacity",
                  "type":7,
                  "value":1,
                  "min":0.1,
                  "max":1,
                  "key":"fillOpacity",
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
          else if(option.style.diy[i].key=="ConnectDiameter")
          {
            this.ConnectDiameter=option.style.diy[i].value
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
