<template>

  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"   x="0px" y="0px" viewBox="-17.4 13.044 41.848 18.323"  xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
     <path class="color" d="M10.129,24.628
	c0-4.825-4.092-8.156-8.201-9.438c-4.11,1.282-8.2,4.612-8.2,9.438c-3.237,0-6.474,0-9.711,0c0-3.208,1.903-5.89,4.406-7.614
	c2.516-1.734,5.641-2.552,8.65-2.552c3.235,0,6.472,0,9.709,0c3.008,0,6.135,0.818,8.651,2.552c2.5,1.725,4.404,4.406,4.404,7.614
	c1.064,0,2.128,0,3.192,0c-2.682,1.774-5.366,3.549-8.048,5.322c-2.684-1.773-5.365-3.548-8.049-5.322
	C7.998,24.628,9.063,24.628,10.129,24.628z"/>
      <path class="stroked" d="M1.928,15.19 c-1.599-0.498-3.184-0.729-4.854-0.729"/>

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
    name: 'view-svg-arrows18',
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
          "icon": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAAAXNSR0IArs4c6QAAA8FJREFUWEdjZBhgwDjA9jMQ44AETi6e7u/fvohwcvMy/Pn9i+H3r59Y3c3GzvGCgZHxwa8f3yMZGBgeEOM5fA5Q4OEXXK9jbKeiY2rHo2lgCTbvzYvHDDcunmDYu3ExQ/20TQwiEjJwe25cOM7w+sUThlUz276wsLLNeff6+URCDsHlAAVuXv7zeU2zBDSgFiP7Znpz7pe7188dYWJi1ijrWaaA7AiYutWzOz4c271+wbvXzwvxhQRWB4hLK9xPKulUwGY5yOCT+zefeP38sScDA4OCuLTCfkocgeEAUUnZ7SllPR7YLH/z4glDR1H4DqjlMI8p2HmFr08u7TbA5tPqJNcXOhb2oTuWzzqCTR7dAQo+UVnnQ1MrBLApnlSb+uLskZ2S6HJCopL9Ff0rQ8SlFRAJAklRbYrHkUd3r9kSdICqjsn5msnrsPoE5PvWvOAJuOIUn15QtP35/8cXWygghwBe37fkBl24feWMIZ4EVR9f2Jrj5Bcrgk0NrlCAOwAUjP2rThZg0wzy/eT6zAUPbl1MxJei7bzCz+NKC0sn1z/YtW6+I3q2hDsAlPK7lhxSwGYBKAi3LJsWyMDAcACfA+SUtQ43z9lhg00NqIxYOKmm9Nn92z3I8jAHKAQmFN4PiMeeZVdMb9mwfdUskAMIgYSuJYeaSUmMYAfgC36QfE9p7IbLZw4S4wCG+MLW17jSQXdpzJ0rZw65IkcD2AH6Fs4dRe3zy3HFf3GkFSjuFxDyPkgeXzrYt2nxm4X91aYYDjB18N6fUz/dAZsFV04ffNFdFguqCIiqXHRN7NeXdC8OIDYdgEOACFeLEuN7kBoFNf35jTM3J+BSX5/ui5KbwA5wC0q8H53biDUHQLOPIrEOAEVnXH5zObYKCmQGVgdEZtVe9whN1cBmyYZFE4+sn9+LtRjFpl5KUbWkuG1hNy4HoCdoUAiAS0BdU3us5f+aud2ESkCYO0AhCML2lf0rG7BVZiCFc7tL5x7atnIJVNMBcBTw8AuexxXEfEIiS9ELD2xqQVn5z+9fKUJiUn9yG2cK4AoBUKF25cxhhge3LoM8zEhMk4zY6GcAlaa42gbIhoBKxZltBeCKjaoOAKUxn6isAlzVOcgRoHqlLs3zw9fPHwVBfGo7gAHUoKnoW+mBOxd4f3hw6zK8XqG6A/A105CDHhYltHAAOLujRwV60NPaARhRUZ+OGvQ0dwByVID6ErBUj56laBUFMHvqLZz8Sy6fPvAHlurp7QAGHl6B1V8+f5iKqzVF6xAgWIgNuAMAmZS5MO9upIcAAAAASUVORK5CYII=",
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
