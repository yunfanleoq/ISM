<template>

  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"   x="0px" y="0px" viewBox="0 0 33.451 41.712"  xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
     <path class="color" d="M9.551,1.417
	c-2.213,2.323-5.92,2.323-8.133,0c0,3.676,0,7.352,0,11.026c0,2.814,0.17,5.526,1.038,8.261c0.867,2.735,2.146,5.1,3.936,7.338
	c1.807,2.251,3.785,3.998,6.423,5.254c2.668,1.266,5.307,1.474,8.156,1.474c0,1.843,0,3.68,0,5.524
	c3.689-3.688,7.377-7.374,11.063-11.063c-3.687-3.688-7.374-7.374-11.063-11.063c0,1.845,0,3.682,0,5.526
	c-1.528,0-2.978-0.103-4.438-0.661c-1.456-0.556-2.65-1.364-3.761-2.449c-1.117-1.087-1.947-2.257-2.528-3.707
	c-0.584-1.461-0.692-2.905-0.692-4.434C9.551,8.769,9.551,5.093,9.551,1.417z"/>
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
    name: 'view-svg-arrows3',
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
          "icon": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACEAAAApCAYAAAC7t0ACAAAAAXNSR0IArs4c6QAABI5JREFUWEdjZBgEgJGBgYGXjYNzh4GFS/upA5u3YHGTKTsnZ9/P799dGRgYftDCzWBH6JjYn5OUVWR58vDWRiYmppNXzxy5xcDAwG/m4JMkICRqff/mxY+3r56zoKkjzBx9j2fXTdW+e/08w6tnDxm+fPzwlpdfUFhcWoFBUUOfoac0+uTlM4cd6OIIXEE9YhwhbmznuT+vcaYmrpDoKok+efXsYQ8GBoYPNEuYsDQx4qODdzQkGBgYwIXVYAgJSTU9s2PVE9co4EqYzTmBD+9cPWtAy9zBp2/htL2ofYEVLkdMqU+/dvrQdhMGBobvtMqi3Lpm9ntKOheD6gasYHZH8cUjO1eDyokXtHIEo6K67r6GGVtBdQNWsGXp1Fur53TGMjAwnKKVIxjU9cy3FbTO9eTi4cNqx42LJxgWTqwue3b/djfNHKGqYzIjpbwnXUJGCasdP398Y5hQlXT62vljngwMDG+JcIgEKVEHyqIgEFXSuahX18wBpBkrOH9sN8OqOR3EhEa1jIKq75MHt/0ZGBheEuFgcDkBAlax+c0bXQLiRfBpWj+///2tK6dyr507uhSPusDQtMp1184c2nz13NEqBgaGK4QcAnOEiIWT//7M2sk6+DR8+vCWYX5PxeEf37/MxOOQiPyWOcvfvnrGcP7ITpBDUgmFCMwRDMqahrNTyntTpORV8Dr84/s3DNuWT7t/7ezRSY/uXZ+JpewILOlctE7XzIHh2vljDOvm9e68feV0Cb4QgTuCgYHBOa6gdYWzfyzeKAG58P//fwynDmxl2L1u3q5///+uuXv1wi4GBoaHUNfDHQHinz60jWH/xsV4QwTZEezKWoZbClrnufAJCBOKRrD8h7evGO5cPcNw8+LJB08f3n7FwMj4/MXje9yJRe0uoJCAAVCIbFkyGadDkB0B0uOfVNI1x947gmBooLvyz+9fDF8+fWD48+c3g4i4NIYn8DkE3RFsBpZO+yKz6qxxlRlEBREORRCHTNl89dwRlMSK7giQdktLF//mjOrJzpRYiEsvthDB5giQ/oSU8r4uW48QUXo4BJcjmJU0DWcGJRR565rZ4yxFKXEgKEQ2LOjbcfPSqQRcjgCZLymtpNEXl9tkr2FgIUmJhdj0fvvyiWFOZ/GLs0d2BuBzBEivvKySRkdISpmjgaWLOLUcAnLAhoUT7u9cMyeHgYFhGyFHgOwVU9YxabJ2CQy2cQ8WYefgosgt6A4AGUaMI2CWBuuaORS7Bycp4qtt8bkQmwNIdQRIvYKkrHK0iKSMj5m9t7KCup4oqOdOTOi8e/WMYf3C/seHtq3MAEUBsmNJCQlkfaAS1VBD38KPlZXVlJmN/beImDQHr4AQ/+cP79ncQ5PlxaTk4eqRQqCMgYFhDXpokesIdHNAlY0AAwMDJ6giLOlcNAFWd+CKAmqEBL6oh9eixDiAnDRBTM5oCk4qqXXwiWLYsmwaPBvi00it6EC2IzCtom/dzStnHh3csiwPmHY2EnI5LRwRo2tm13v51KEiBgYGfG1RuNto4Qh9BgYGa2AhN41QCMDkAVPSC5Z8jvb9AAAAAElFTkSuQmCC",
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
