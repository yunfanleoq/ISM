<template>

  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"   x="0px" y="0px"  viewBox="0.309 1.2 23.174 41.93"  xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
     <path class="color" d="M14.352,39.216
	c0.02,1.362-1.043,2.48-2.413,2.497c-1.362,0.014-2.471-1.083-2.471-2.442c0-0.02,0-0.038,0-0.055c0-9.213,0-18.426,0-27.639
	c-1.317,1.829-2.636,3.66-3.953,5.489c-0.666,0.949-1.957,1.191-2.91,0.519c-0.944-0.668-1.164-1.978-0.486-2.917
	c0.01-0.012,0.019-0.023,0.028-0.035c2.667-3.693,5.336-7.384,8.005-11.075c0.838-1.234,2.646-1.258,3.516-0.045
	c0.009,0.012,0.019,0.023,0.028,0.038c2.652,3.696,5.305,7.389,7.957,11.084c0.685,0.935,0.501,2.24-0.443,2.929
	c-0.935,0.685-2.251,0.472-2.922-0.472c-0.009-0.014-0.019-0.026-0.028-0.038c-1.303-1.826-2.605-3.653-3.908-5.477
	C14.352,20.789,14.352,30.003,14.352,39.216z"/>
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
    name: 'view-svg-arrows21',
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
          "icon": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABcAAAApCAYAAADeWrJ9AAAAAXNSR0IArs4c6QAAA6VJREFUSEtjZCAN8DAwMHwhVgsjkQqVVXVMSji5ePQ/vHt199Gda5MZGBhOEdJLlOFKmoZz3EOSky2c/Bju37jEsH/rsm0HtyxLZ2BgeILPAoKGy6lq5XmEpE20dguCm3P9/DGGud2l7a+fP66ixHBzt6DEFZHZdQpMTMwo5uxZv+Dt0Z1ra+7dvDgDlwX4XC6ia2q/OSav0UJCRglD/88f3xhmtRUePHN4ez4DA8NFbBbgNFxMRqEhLrepXtfMAafPH9+7zrBmbveRC8f2eGJLRbgMLw1OKkn1i81TJZQizh7ewbBmTlfls0d3OtDVYjNcU8/MYW5cQaulqKQsIbPB8gsnVj+8c+XMhEd3rk9A1oBuOLucis7s5NLOWAU1XbwG//71k4GVjR2s5sPblwwzWvPPXD9/LISBgeEhTCOK4VKKqqWBcUVdZg7ecIO3rZjB8OHdq49RWXX8MMH9W5Yx3Lp46nN69QRemNjNS6cYZrUXdL958aQMm+F2XhHp88PTq+FJ49SBLQwrZrTudQ1MNPQMTxOCabpy5hDDtObcVfH5LQ7mTr5iMPEtS6e8O3VoR+3DW5emgcRgLpfWMLDckFbZZyIsJg336vJpLVdO7NvYG5Pb1OMalCCMbHh3aUy+roldVGxBi7m4tAJY6s/vXwzTm7MPnTm8s5CBgeEc2HB+IbGZ6dUT0rSNbODBsXvt/LdLptT3MDAw7InJbdqBxXBQ+D5z9I1ekVDULgfT+PThbYZpTdnHnty74QYynE/XzH5nSediC5iC88f2MKyb37Pk0Z1rsQwMDCZ4DF8rJadSbmrvXR6UVCwI079kcv2Z3evmx4BdLq2ovrKofUGYiLg0w+0rZxnWL+g9dfXskQAGBobnhAxnYGDgMLX33m5o7epg7RrEAEpFE2tTj18+dSAAFuY2iuq6zWo65twXTux98fLp/TkMDAyboC7B63KoGm9pJY0YLUMr5Qc3L76/feXsVJB+5KQowsDAIAoKRwYGho9IiZwYw0HKQYlekYGB4RPUDHhqwZdhiDUcwwyC5TkRYY7TYaOGYw2aAQ8Wn4D4okWBCQXwHHhszzqGma0FGQwMDDPxJbMBd/loOscaPaPBMhoskBAYzaFkF7mmrsGJU2NyGk1hJpzcv+nVtKacVKQWAtmGs8sqaSyIzWuxV9c3k3zx5B7DnM7iXbevnAUVufcpLXJB+pWBXZNADX2LiFfPHh5+9/r5CgYGhpP4DAbJAQD4YvaBwNcAAwAAAABJRU5ErkJggg==",
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
