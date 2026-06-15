<template>

  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"   x="0px" y="0px" viewBox="0 -0.688 32.828 43.113"    xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el"  :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
        <polygon class="color" points="
      1.284,20.83 6.294,7.634 13.2,0.73 20.105,0.73 26.399,7.023 31.411,20.83 31.411,30.791 21.387,40.871 11.305,40.871
      1.284,30.791 1.284,20.83 	"/>
        <polyline class="stroked" points="1.284,20.159 11.305,10.811 21.387,10.811 31.411,20.771 	"/>
        <line class="stroked" x1="11.305" y1="10.811" x2="13.2" y2="0.73"/>
        <line class="stroked" x1="21.387" y1="10.811" x2="20.105" y2="0.73"/>
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
    name: 'view-svg-3d-octagon2',
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
          "text": "displayConfig.ToolBox.Diagram.Octagon2",
          "icon": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAAAXNSR0IArs4c6QAABPpJREFUWEftl3tMk1cYxp8iiHjhplhQYEA1sCGIXFo6LhaBggyqbdggAmrVxSXKNOBCmMni2EUdxJgFp8lkopMp0+HCxRYECgWKhUonKgwVnEADExzIQOXa5fuEpmChZaJ/ef5pvp73ed7fOd+5vB8Fs2/rAYxMIzMCUDZD/0syymzyGxubhVFtaSfXegdQbWnOY1O1hZcz5o2NjWY3yiXJABS6eOsKYOLs5ZdA9w/jd3W02/D4iZinr/+Sv1SUB6VSidryArFMLEgEINMGoQuAk4d/6LdsLp+7lLoSd25UghW+RaNvg1wCKJWwsqVBlJdVX1NWkNrRev/8TBDaAFgbOHFH2JE7GFY2NDTKJRgdHcUaTz+Nnu0PmtDe8ie8AzdBqRyDVJTfXV2Uc/4PaelhAI80iaYDWOG0jrnVy3/jxz7sSAejhYtJbYXgV6x28YSltYNGgL7ex5CW5iKYx1f1NzfUQXjpdGFNWf5uAA+nCjUCUG1o0f6hkefDt+yZNyEg3u3vZ48jImYP9A3mawQgYq5kHsPmbfuhp6eSQlYhxNWLp+KbG+rSdQEwZWzg5O787Ds/wwULVfHdne24VVuOgIiYGddV8ZVMeK3/ACbmFqq4sbFRnElLrhYLLm6e+ipemgGqjV103N6vzrrQ16uG+fzZAIp++wmKB3dh9Q5tRoDO1hascnZHEHf7pDhi/WSlp+xoa2k4o94xFcCUwQrP3ZmUphq94q+7kFeXwMHRFY6uDG27CsRob0pL0f+kF+6+bBibLiU1Q4PPkZF6oPJ6SS4HQM+E0SQA6gq76Nh9KZmudJYhEXBTKkJnWwvcvDeAam2vNbl6ALH4iFfm7hsCW9p7ZFe9VDR45vjnsf90Ki5rAjCls8JzdyWl+Q309ZKjNjO3gCszEPr6BrNKPhHc/bcCdZWFsLS2hysjAM+e9uP00QSRTCyMADBAxKlmwHKl/Ucxn355boHRIsMGeRXcmEGwW73mfyVWF42MDENedQ3Pn/6LdT5sNNVLn/zwdXzUyOBgoTqAqfv77KvuvmwmRU8PbsxALDY2e+Xk6gb3bsvIg8xxLYNY0AKZWEDsiCFyBkxMlh1gBHGOeviG6Dm5Mec0sbpZV0cr6iqL0NfT9ST/wkliMYoJgKXBPH5jEHebhZU1DUooQQGF/J3UJj0SDy9iKBRtp7m6ixLDQ0OQiYWQFOc01UtFHoTaftPW/UU8fsKq1zZ0DcYnUvbcrhHl+ZIAUZ8cFIRF7XZ8kwA/f/9FbfGVzOC3AG9nYHwX7Cvi8RPf6C7IPJZ8Q5SXFUgC0FkRWZy4vUwbh3ffyEZ4USVlCGvK8qInThF7Tx92CosT+6HL+E34OkiIiqmmLL+7qijnl5vXS74hihP1Y2yBoxsznu4XutsnJJJmtGjJnDL0Pn4EUX7WrdpyYZriQeM5jfXA+J+sAE7MkZDIXWQlPBft4b07yM9Kr6gpLyC+FWrVPac7yB09/EMPs3k7uE5rvV+JQS65NlBecCFbLik+BKBtqtlMN4mJi5fvATqLs50ZxLU2mE8WSTq3/r4eiAWXmmTighPNDfJTAIY1ibVeZUvMlkX4BPMOBXO3uy+ztNYJgKgjBdk/SiqE2UkAKmcSaQUYF7t6+G9MtVvl7Ga+fOV0X8ZkaGfb/fmK1vsldRWFBwE0ayPWFYDwWQ7AUpvheD+RmKz5tLX/ANb66LM/x2+fAAAAAElFTkSuQmCC",
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
