<template>

  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"   x="0px" y="0px"  viewBox="0 -0.004 61.111 29.453"    xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <path class="color" d="M11.739,26.906
	c0,0.473,0,0.944,0,1.417c-3.319,0-6.64,0-9.958,0c0-7.591,0-15.182,0-22.773c3.318,0,6.639,0,9.958,0c0,0.47,0,0.937,0,1.405
	c1.466,0.042,2.964,0.214,4.371,0.018c1.867-0.259,3.574-1.164,5.251-2.137c1.674-0.974,3.318-2.016,5.128-2.684
	c1.812-0.668,3.787-0.964,5.636-0.545c1.85,0.417,3.57,1.55,5.193,2.52c6.439,0,12.878,0,19.317,0
	c1.692-0.184,3.171,1.144,3.171,2.848s-1.479,3.031-3.171,2.847c-4.06,0-8.12,0-12.18,0c1.515,0.164,2.635,1.482,2.553,3
	c-0.081,1.522-1.337,2.71-2.86,2.71c-1.04,0-2.078,0-3.117,0c1.572-0.169,2.984,0.969,3.155,2.54
	c0.169,1.574-0.968,2.986-2.54,3.155c-0.952,0-1.901,0-2.854,0c1.583,0,2.864,1.283,2.864,2.865s-1.281,2.864-2.864,2.864
	c-3.237,0.378-5.031,0.664-6.835,0.84c-1.803,0.178-3.616,0.248-5.426,0.209c-1.812-0.038-3.62-0.186-5.412-0.429
	c-1.793-0.242-3.569-0.581-5.37-0.678C14.421,26.828,13.077,26.887,11.739,26.906z"/>
      <line class="stroked" x1="11.739" y1="6.955" x2="11.739" y2="26.906"/>
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
    name: 'view-svg-arrows36',
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
          "icon": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAAAXNSR0IArs4c6QAAA91JREFUWEfFl3tM01cUx79F6IOBymyA0VBB2YBKQLTGKXEx0al/GN3cgrpIdESJMT4xBqbglhmWEmdJZzbNeKZqo0ESUTaVjKBAEbEWVOoPtK2VR2lpR0VkPFq75ZJCqpZiWUvvXzf3d37nfH7n3N950ODlRZui/QgA8fwV67YODPTrKLlUAuDeVHRNBYAbk7j88Pqte/ZHxy+BsbsDjTXX69oeNJx9LJdecBXCVQB2TMKy75PTMvbO5y0at2W1vsalsznlHc/aLipktZfegiA2/nUANnrmEkDsoiRh8q6MQ/NiFr6jz2TQ4/YfEqmuSz1k/5DJChweGuxn2J+FzY0y3qv+s6FdTZUQACaARJtAHIAW274NQO/Yi1Fx/NRlqzYWrv5iu6tedigvr7upF2Xv2kyL5iWmfrrm68JQ7nx0aZ6AE/EJutuVqL56Ib1DTeUB8IuK46ckff7l4ZXrv+H5+MxwC4DJqIPw6I79tNiEpSmbd2eJI2MSxhUrFfdx/vTxW0EfhhjM5mHu2uS0pXH8z0CjuRQxp6DmkWGUCL8rdgjQa+iGvksDJuuDUY/QGSRK7l19JiMByJkQoK/XgMjoePdatdPWIqvRnTyyLcWjAGdO7O0KDGJrWf4BZot5xMfXj24lDDNns4dkdTdaWpvqcz0KIJdW6ksLBHlajbIYABuA0c4JBpIfPApAjOUd/fZy852qVAD9juLpcQDJbycab5bmbwCg9wpAQW56Ze2NyykAeqYdYPCfVxBl7xRR8voMAMPTBlAh+RUWs9n00mSsryoX5wKoneh/9sgd0HWqUXQyQ9728O4PAK45Sya00PDI7I3bDvy4fM2mcTmSCf9vIqKa6iFI37ISwG2nAI5qgTsAiBd+FxxMUyma870CQCpr0anMncoWWaFXAB7erdadytye7OwCjnZE7g6BrOY6QjgRKC0QlD1oqN4HoHtaPXD13C+m5jtVZSqqSQBANVk5dbsHiEFb9tsx2dd7JARE6V9XSqhzouNfAaC84oFb1ySKYmHmFrsGd0IOj4RALMqqqLoiJhdQM2UPKB/LSecCBpOJOcEczAwi/YTz9dpigebpI5QV/SxUyGqPAXhjRnjvYkS64otncipn0H1brVb4hHzE/XgeL3FxcGg4O5gzF/4Bs8DyDwCZiDrVrdBrNdA+V0l7ezoHtM9V7Sqq6ScAzyYDHr2EZOCIil2YFRIWMdjf97df4Kw55hcmA71cLDoEoMKmhAwDPADcBfyk1QwGi0OGKl8/+ohGqVD2dGqUANQAGgFY3sfwmMzYZBRmOwgH0GHbv7CfjFxR6orsfxUAHHRqRxReAAAAAElFTkSuQmCC",
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
