<template>

  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"   x="0px" y="0px" viewBox="-1.387 0 29.678 29.449"  xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <circle class="color" cx="13.452" cy="14.87" r="13.422"/>
      <polygon class="stroked" points="3.02,11.49 3.02,18.25 16.113,18.25 16.113,22.706 24.455,14.897 16.113,7.034 16.113,11.49 3.02,11.49 "/>
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
    name: 'view-svg-arrows4',
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
          "icon": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAAAXNSR0IArs4c6QAABYtJREFUWEfFl3lQU3cQxzcCARK8gIgcJgQCgccRQQTUioJWIlLUqqg4limVioU6lV4WREY5tB1LOyoODjYdq3jjgVQDVailqCgSwvEgmBCICCIiHiSQINB5TxIFEuVldPr7K/N+u9/9/Hb3bX6PBP/zIhGMb2Fr77TSwNAo6EWfik4ijTPB/AcHB3oNjciy/hd9Rfeb7p4BgM6x6o4JgMl2D1cpVfHmNGur2R9+bOnsMdPMcqrdsBiPHrRAQ/Xt7ut/nX30uKOtnWxMzpCKak69DeSNAHQWx22gv5eHeM+hrYj6hmlCMXubHr7fq+iGXN4eKVpR2jHOwCRKJhbW6nLUCeCIeMeRSKQtMYl7HWjW08YUeKRRR9s9yErb3Dg4OPiLBK3Yr01EK4CL16z0yZNpq2OS9jvoFXmEU1ZKXGNXV8fJesGNhJF6owCcEZ84i6k2W95VcHVADKKzvTWjobY883WIYQB0hONmTDLM27b/HOGTV968Cq1NdyFkTYzOpKXGLW9UDr4Ik6GvemIYgB2TXfZVGs9Xn5pfPnkQZOJamBGwGHzmLtYKgfXEr4lRt1qkIj+1gQaAhXDCHRCf3etik5n61L0o7yhYWtmCBBWANd0R/Bcs1SqTk7lD2oiWbxWjQvwV1QBgp0/KPO9rYkrVJz5gALb2TsBkc4B/OhvMadbwQfDKUVq9PXJIiV2myYIagMbxD7wVv+uw/cBAv14AxRdzwI7JBranH2Aa/FPZQDGbAPNDI0bpZfwQ2SS8WewLAB04gJ0j8mXYutjdfoEfUQrOHII6wXUwpY4nBCLvfgpL1mzCAdSrMJcH4wwMYOGyyGFaZcUXFXk5mVtbJOg+HIDuiJz/+sc/lk6ymAJ5R/aCT0AI2DBYhAB0GRddOAJ9KiUEr9qgMXnS+RB+/v6TCzIJuuxlBpjOt9J4V2Ziv981AKZ57dIJ6H7aBdzwaDAwMMRBEqMW3m6RNvjiANMcXatSDxV4vC8ATLe08CxIRUII/3wrkI1NYduG4Op7kjpPdQmqUw7x3d8nADYjju5Lhsj4XWDLcIKkDdwamQT1GCqBS0Uar9DrfQFI6gRw5dxh4K6KBoaT21AJFglapPXeLwHs2VeTs/KCsNS86x4QVZVBaUEuBIdH4yfHlkrZAztiwopamkQLcAAGC0mLiE1OcJk+C3L2JQPZhAITzWmE3gJF9zMICFmNDyD1qikvgYp/C/Dmm2LD0Dyvr7wBxzJ3pDeL0UQcgEqlcvyCll+IjE9ntDaLoUfxnFBwzPhOCR84/kGaOVB54yqgglI87a9DYbaHMxKay4rOLZXL5ULNKKazkNqUbD5COPKQg3oUY4OovOQySOuFePDxkyxGSSZFc1GZGMWbQQPg6OadGhgaETWXG/4qhwRo1ABdHQ+gTSbB065tmpbwT7UV5x/jSWortg0DwHvB2b1u58FLLgTiakwxgM72+2BCoeInNyIba5XZvjGkvrmhxlW9Oew+4Dpj9iYrG4dvP41PJ/yXfCIrHVRKBazfnAIkkvar5u8ZCdL21sY9dXeuH9AKgD1ksjkn54WunRsYGqFXKXRlrzj/WNu1/OMlUpFw9es2WlHpLOQad1U0e86iFVb6lGOkT2lhbjv/dLZIJkbnjdzTeS1nOLnne/rNd1/52XevXmA9aM789lNzVdnfNc13a0K1ub/xw8TV0z+1r79vzdovtk9lIV6ErkpiVCA/fmDnAyMjoxN1lTfxjicMMOQwxd7Z4yDZ2HR6UNj6iT4B3Mm6Ohz73y//h99VlHfkqUrZU9nUUL0RAB6+KXFj+jYcEqCykOmbFXLFErIx2dKUOt7UbII5vtX97DH0yJ/3qJSqRxQq5U8xWrkXAORjqRgRAG39M2no4RPsI3ksAUfa/AfdfRc/VYb6AQAAAABJRU5ErkJggg==",
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
