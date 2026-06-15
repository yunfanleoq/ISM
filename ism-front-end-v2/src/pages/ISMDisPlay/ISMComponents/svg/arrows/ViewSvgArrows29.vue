<template>

  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"   x="0px" y="0px"  viewBox="0 3.031 43.791 24.749"  xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
    <path class="color" d="M19.468,24.523
	c-2.511-0.07-4.837,2.017-7.308,1.84c-1.236-0.088-2.507-0.742-2.699-1.749c-0.192-1.006,0.698-2.364,0.82-3.711
	c0.122-1.347-0.523-2.684-1.321-2.535c-0.798,0.146-1.751,1.776-2.933,2.355C4.844,21.3,3.432,20.825,2.59,19.96
	c-0.842-0.867-1.113-2.121-1.168-3.558c-0.111-2.87,0.641-6.454,2.568-6.416c0.964,0.018,2.222,0.941,3.531,1.246
	s2.668-0.01,2.855-0.959c0.186-0.95-0.801-2.537-0.823-3.703C9.53,5.404,10.47,4.658,11.553,4.488
	c2.169-0.34,4.909,1.625,7.43,1.478c2.527,1.503,5.519,2.343,8.464,2.237c2.344,2.086,5.432,3.489,8.596,3.528
	c1.57,2.131,4.55,1.748,6.331,3.531c-1.966,1.665-4.714,2.142-6.464,4.14c-3.162,0.089-6.288,1.197-8.928,2.93
	C24.307,21.879,21.579,22.922,19.468,24.523z"/>
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
    name: 'view-svg-arrows29',
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
          "icon": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACsAAAAYCAYAAABjswTDAAAAAXNSR0IArs4c6QAABJ9JREFUWEfNl3tMU2cYxp8iNwGBwUAxgAgbonPgdMJwixmOJTLmDFkymRhGli3iMA5lGRlhEUeWQUaYyRaJ0ShqNjUKMhHkJhRYW0pLC21pgdIWWigt9Ma1hba4HUIX2BRaBtMvOX+d53ueX3ve837vIcG25QSABMBg27bVURPB1ixPX/+gdJ+N/nGE2GQ2qbrbW+4CYADgWWOwGhprYU9k5F27EB4VM5c5Oa6DpJujEPJYTmIBq8tommmdHNeJpUI+HQBxrcmyCnZTQHBb7qWHux2d1v8LggBXKwcxNTGGro5WrWJAxKLVleXM/+vTq0ltFWxQ6KtXP/0qP2XLyzuXzZaJBeCzKOrOtj/U+qnJuhe8fZX0hnIy8UAAtC1rsITAKlgvX7/CQ0fTTh84nGx1lslkhEImhnZkCOOjGkyOjarHxzR2bFrtWamQ/7PVRguEC2G9ALwJQAdACqB/gW5jbEIK/eMT2VvsHRxXkjO3Z3bWDHLFLVUXi8qlk8vPAGi3xcwC6xEUGnE+NiE5XqdSrBuQ9MjVwzKKkMcqAjAQEhbx4eagbdmfZRYE2GL+NG2fkAdWc5WWx2wuEQnYtwF0AhhaztsCG3vkeNYv7yWmbiM2mM1GkMtvok/IHV/v4qZ6cZO/e1TMB94eXj7L+dl0f7CvB72dLFU3hy7XaYb1ABSdzOZiAGVPMrLApuYUlRdtDYuwKWw1xTPTBujUSrTU39fQHpXly/uE1wn4hRlzsNt3RT9IyymK3+BBlO2zX62NFcMMcsWIWjlIEwnY9wAQ3WSKgLXf8do+ambhrb3PHnMxgXJAgm4OfaS1sbKX20o+TcBuiHjjQNWZH4r3PW+wBM+4ToPqu5d1HEb9twSsZ3jk25UZ+dejnydYon65jKYRSk1Jv4BNvQjgMgHrELYrmvrNT7dffxIs0dxnDHq4uLn/b7+F0ViJxoqbLVxG468AbgAYJcLnXjAfv4DvT+VeygoM2bEISNLNQW3pVY5aOWhMOnl2T+BLr6wZ8OPHs9ColKBU3dGWXCkgDow780f035mW1hUZcyipKOlkzm4HR2JkBQTtNNSUXmliNVd/AcD1o88zb8QfTQtdLVqDfhIykQByaS8U/b26QZlISQJJ1U6rywPwYKk+i5CwiONunt7Je946GKpSyHSCdppYyGMW/HUE1wIIezch5eGxU98F/VdYlWIAAjYVlNpS5jp7RzKPQSZak8SWE8zC4A8gGMAMAA7R2+ZvvJ+YmlUcdyTV2xZYs9mEUfUwpg16mE1GSEV8oulTO1rqzwFosvWLw6qpC0BuRt61bMvwvRywuKsdAhZNw2ltkBgMeo6jk7P+8azZTshjPgLwOwDjch5LlsESm7eGhu+98PWPvx201PPTtMPyflBq7mkHxPwqZnNV4fxUZVoJ2EphnWIPf3L+2JfnUkkku0UexBvcw2VCxGdrJD0c+YROo+ezKekAqKsFuNDHqjIIDN6evjNyf+L+uMQoZxdXTIxpIe0VjNAb7guHB6UVQ7LeagCi+Vl4LTjnPK2CnU9/x8cvIMbV3Wuzs7PrRFcHlQBkA5CvGd0/jP8EEYT5OuC3EOwAAAAASUVORK5CYII=",
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
