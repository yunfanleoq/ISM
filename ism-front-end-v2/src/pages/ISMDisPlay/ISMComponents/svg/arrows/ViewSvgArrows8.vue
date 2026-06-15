<template>

  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"   x="0px" y="0px" viewBox="0 3.842 43.357 22.055"  xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <polygon class="color" points="41.941,14.87 30.069,5.26 32.367,9.999 1.345,5.26 3.637,14.87 1.345,24.48 32.367,19.608 30.069,24.48 41.941,14.87 	"/>
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
    name: 'view-svg-arrows8',
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
          "icon": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACsAAAAWCAYAAABZuWWzAAAAAXNSR0IArs4c6QAABJJJREFUSEvFlwtMW2UUxw9lEwqbUXkIERHlNcZwbGyixI35GoMtJjqDxGncBrgtG1EzcJuyLUhClLC5IJiIiRrcIBuK4x2mIAzCSzJglLWlpbTc29d99P0CSmtu15KmA3q7DLzJzW3ynfM/v/O9zqkX3H/8AMBg/71WHwYAWDwJ5gUAIeGRm+uY/ht0FotFOGvQc2cE7AkAkDu9nmjSsd3j57/xjEGvPQoAUjoOlA0FC7EvJrfmnClLX1gwg0ZJACET42oF7i0WTiIkJp5neHuLwQLTE3d6elyS0NIN5GLnHRAU+ndgSPjj3PHBbAAYpaNjgw0NjzqXe/ZSSWTctgd8HAmolQRoFDgQcpQgpIgGkyFGFYGhFrAK5wy6SUTAcV2NFZc4OmHn4bjElMsi3jh7bKCzCABuuQO2wQLAG9kFpbW7M7KC3Dk4jxsNOlsCVCIkJsHUCswbEXBRrYpgmufmORZYmGbf6et1WQ21XcMnOj6p6Uj+t2/9292mHL7dWowI2N+tFN8B+9zed490HcorivAEdjlbq9UKagUOGqU9EZmYwKQiNS5DTUpcIrYsWEVGk5Zn0mnCXk3L/OhgdsGTXU01RHNtZRUuRS4DALmUtgMWEnbsGjhV9GOyr9+GR8G7rIbJqLedC2o1ZIgADwwJC4pLfMVmP9L3l7yr+Rp/tL/zEwC45yqyCLtpa0r5oVPn88Kj4lcV1p24iDcBdVUlg+PDPWcBoMvZfhGW2vBpB7N/2Zma4U5v1cdVpBw6blYrx4Y6Logm71U4Ai7CAkDM/qzjbZnHvnxh1WloBKBuoX8ar5Lt9T9XYKjwEgBonWEZm7YmDxWU1SStW7eehtzamAx1t2Ad9b+inLuD7zvDQkRsQldheX3q+sd81oaERhTW8G2orSzuR4XcXGfYqNff/vD6x5+XbKehsSYm3c21RE973R881nAhABDOsPtOXvyh7aU9B9YEZKUgOo0S2q5XTTXXVH4NANUPHLC47SlXso4XfhoRveV/hRWLeNBYXc4a6GzIB4D2pa+uLUntp7+p3sv037iqsFoVaatuKgUOEhGffCo41Lxj176nqaDskT5oqamk7tgTVI1Yrig8uzs9szH7i7LER0FqNs+DWoHZoNQKqotDSVwq0qDTXNzH12/WoNOMcsYGqJ4h5oOTF4rT3suB3vbf8c6G3xqm2CMXAUCyUrl97Wh+6Y3U/VmBnsAa9VpQKTBb+cRlKK4mqbaSgypx+azBoGPNGfWTkhk+VTZlTu+cPYZv5OZtTSe++v7Noe5WRe+tG6USIb8MABaWY3AcsJzzFX/+FBWftKSdbYaopkSBU1AEIRFpZgQcOYPBsBp0mrs81nCfCxBOI+n0dw6fvqYm5ZrOpqt5ANDkzud+8534csuxc1cyrBaLbfmo/UTIEBKd5iK4FDEzvLykFrAK7O2eY5aofxIP23xDWERMvw/T32+KPZILAEPuQKlxCjb4medjbz4REMww6XXjRqOWLxHyWfaZooAoODMdMQ9sDvgy/T8zGfU5ACCk6+fYBgHL9ZB0hR7CjgkARk/8/gPOQTyEqKH4cwAAAABJRU5ErkJggg==",
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
