<template>

  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"   x="0px" y="0px" viewBox="0 5.032 41.247 35.197"    xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el"  :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <polygon class="color" points="1.447,38.811 8.992,14.069 15.003,6.45 39.66,6.45 31.993,31.194 26.165,38.811 1.447,38.811 	"/>
      <polyline class="stroked" points="8.992,14.069 33.895,14.069 26.165,38.75 	"/>
      <line class="stroked" x1="33.895" y1="14.069" x2="39.66" y2="6.45"/>
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
    name: 'view-svg-3d-parallelogram1',
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
          "text": "displayConfig.ToolBox.Diagram.Parallelogram1",
          "icon": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACkAAAAjCAYAAAAJ+yOQAAAAAXNSR0IArs4c6QAABZBJREFUWEe1l3001GkUx78c6rRNelEnbSqSZrIYrKVJsqqlElvyshEqldN6KUpt1g5pU6dQ4dBaoYytUPJWtjcq2xR5SxhJEzuJXo4WKcTs+Y3VUck8M9s8//7uvd/P797nuc995CD9GgFgpPTuxJ7tcsSmHxjO0PzqF9qY8atH0cZ0SRtjOL+e7i6FFgFfuVnAny8tpNXKtb7RK9x81WQBSMU8mxiOooKcG08a6xdIA2kwd6FNgptvKPMLmpJMGO/cuIDqkkLh7fxMp462tlOSQqrM0Welum7Za/rljFkyAeTX3kVp4Z/oaGu9fzWLsxxAnUSQdB3jYyvWbl2vZWAiE8C2ly+QzYmCkbk10uMPXOSVcy0pIWLIKWqztlrZe7BNlzmOlwkhgAupcZg2cw7k5eWRfCTQp6mxPkoSyMVWqzfHOGzapSkrwKKCXLzubIfZsh9wJiH8YVbyEVsAFaSQpjPpuomWjh4ao2V0UBrqqtDT3QVrZ0+8fduDQwHrubxy7nwAfUSQarO1g82snIKUxivLKoko/esijM1twDQ2x4OqEiRHs0Me8SqDBgTF7kmGHqvQLzTBZOSo0TKB7OxoQ0ZSBJy9gkXx89Lim07GhNgB4JJC6ixf45Vh775Do/+MCT87aO3d22h5/AgLljpCKOxD+E630sria2YAOoggJ06Z7rHOL/SotuGCzw43EPByRhI0tAygTtdFY301jh8KjH5Qdcd7sOCw5WbosXI2B0ZajVOe3O9DWVNbWewmIfuntz3dSI5kw8UnBAqKI3D9/OnWYwf9HQBcJoVUM7Nyyl2/fb8WmaTkVnxeBep5FVi8wlXkHBPiWXE7P9sCwFNSSGtPdkyakflymY1j186fhoqqOui6RmgR8HEsbCentuKWy4e/+8nCMXTnJq3z3++mojpT8hQReqREB8N23TaMGj0GxdfOd8fu8XTo7e3NJIWcaGi27JJ38FG9T+tRJ136zSng16KyqABLHT1EEkkRu6rys1OsAfBJIb918w1NX2izRkwHl74t3bpyDrSxE0B1jpcvWvDbXp/s6jKuzVBJGTIVGlr6Ya5bft2mNltHbOGkxUyN2wdLuw0YO2GSKKNJEQEbn7cI4kkhRzKYrJv+YRwDBQVFsZD9DV6ysj9tagD3Sia+d/ERxU+N21eXezJ2JYAqUkhDuw07zlg7e00nIPzPRLJ8ltzIA+Tl8LWJJahrMZK98XpNGZe6ZYZcH6VAVWOOv4v37gMM5lxyRgktzx0/BBMLO0yaMg21d4vAiWIHND6o3kcMydCfd8U7OHYhTUny2VYoFEJObvjStz5vxtXMZKxy9xcx5aRE/50Wf2AVgGJSSPoSh02ZqzcH0iVMDrE5dUg6X7XD2NxaNDuG7VhTUlMmmh3fEEEqKSu7uG8LO6HHWkQsKqlh7qlY6LEWY+oMTZRzLyOLExVeX122fbg479WGocdK2/hThN3EyVMl1X7f/hPn6FX7S2QmR8LpRzb6+npx4nAgLz875SCABFLIqSYWtnmbdh3W/n+ElPfQbYlXzsWzZgFMl9jjahbnWe6pmJ+fPxH8Lk5vcCaXeAQcTpv3nS1NnBPZ94/TeelsIjR1vkFn+z/UNB53v7K4/04Us95B0nWMjrr67vVQVZfNmenueoOTMSFYtNINnCNBhTXlN+0BNIsDpL4PQI5lshbl+4Um6pM4kdoIIYTcfxL1NWV4WFOGR3X3Kgvz0t2Hazkfxh+ANHH2DEq3sHNXIQUgtaP6JtU/C3L+QENd1TPupQyv1687Ukn932VSnc7c4+KzO5B6a8hqxe33QxO/Npp/v/K99wuJHpVJeXU6s1Tb0HScgqLi538OUm+Bxw20p08auurulVBjVSsJ2GCbgXLL7B0zSIzSGnLKEQf9LxHm9lWlGA/dAAAAAElFTkSuQmCC",
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
