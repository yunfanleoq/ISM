<template>

  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"  x="0px" y="0px"  xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el" :class="{'animated':true,[`${detail.style.animate}`]: true}" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <foreignObject style="overflow:visible;" pointer-events="all" :width="detail.style.position.w" :height="detail.style.position.h">
        <div :style="animatedStyle" v-show="detail.style.visible==1 ||isStart? true:false">
          <div :class="{
                'animated':true,[`${detail.style.animate}`]: true,
                'color-animation':isStart&&animateType.includes('millcolorGrad')&&!IsToolBox,
                'blink-animation':isStart&&animateType.includes('blink')&&!IsToolBox,
                'scale-animation':isStart&&animateType.includes('Zoom')&&!IsToolBox,
                'rotate-animation':isStart&&animateType.includes('animateSpin')&&!IsToolBox&&spinDirection==0,
                'rotate-anti-animation':isStart&&animateType.includes('animateSpin')&&!IsToolBox&&spinDirection==1
              }"
               :style="{
                                      width: detail.style.position.w + 'px',
                                      height: detail.style.position.h + 'px',
                                      'background-color': detail.style.backColor,
                                      'border-radius':detail.style.BorderEdges+'px',
                                      opacity:detail.style.opacity,
                                      borderWidth: detail.style.borderWidth + 'px',
                                      borderStyle: detail.style.borderStyle,
                                      borderColor: detail.style.borderColor,
                                      transform: detail.style.transform==-1099? 'rotateY(180deg)':detail.style.transform==-1098?'rotateX(180deg)':'',
                                  }">
              <div class="view-url" >
                <iframe scrolling="no" allowtransparency="true" :style="{'z-index':-10000,'position':'relative'}" :src="Url" :width="width" :height="height" frameborder="0"></iframe>
              </div>
          </div>
        </div>
      </foreignObject>
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

import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
export default {
  mixins: [ISMChildAutoMixin],
    name: 'ism-view-weather',
    inject: ['getNode'],
    props: {

    },
    computed: {
      animatedStyle(){
        return {
          "--blinkSpeed":this.blinkSpeed+'s',
          "--stopColor":this.stopColor,
          "--startColor":this.startColor,
          "--animateSpeed":this.animateSpeed+'s',
          "--animateSpinSpeed":this.animateSpinSpeed+'s'
        }
      },
      textAlign: function(){
        if(this.detail.style.textAlign == undefined) {
          return "center";
        } else {
          return this.detail.style.textAlign;
        }
      },
      lineHeight: function() {
      if(this.detail.style.lineHeight == undefined) {
        return this.detail.style.position.h;
      }
      return this.detail.style.lineHeight;
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
    components: {},
    data() {
        return {
          detail:null,
          IsToolBox:false,
          editMode:true,
          Url:"",
          width:600,
          city:"阜阳",
          weatherStyle:"ya",
          height:600,
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
            "text": "configComponent.weather.Text",
            "icon": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABIAAAASCAYAAABWzo5XAAAAAXNSR0IArs4c6QAAAgtJREFUOE+t1D1oE3EYx/Hvc5dI0hZiVFBEUUQt1IKtFgkOUtBFUUEQVLBzvQiiS1scbEWQFsW42NSlEVxEIYMOuohFpKVoa6mbg6Fi8AWsVVJyebtHTjQmNgkpeOvzez489/zvf8J/eqQeR6O+TrHsMTebj/pPCtpsWvbl0t66oELUfw406DbnR/z3UeIeK31v2VAu6usXkYMoW1FMDF6jjKAcV9Fhr2WPLZkollAfsFkhcOx5643GxXcPDCFSaQVi6CWz277i1sqg0TltEcUdfcefxuCP2a+Hx0OrK+9SH3ss+1AZFEtoJ/Cs2vK7njQsKYnwQZUNgg78miiW0LbfyMpq0KZPcfbNnC4ri5Awz6S3FCcaTegJgbJT+BfMFaDlTS9N2SQbU5MEsklQHnnC6aNF6M6c9qQyDGULIAIrDGjw/qW+2fBlsZxunY9zKtHVrY7e9lhptw36p3Ui4xAqjQZ8sK4RPqdgIVNl1RAenJS7cpaU9L3SMHCrUnSVD+bt2t++bRK82S4L0jelD1GOVIp7Dcg5taGCQ/u1PTLjQkmU9RXj7otrbSifZ+f1kMxK75RGRDlfz+WtkHmbydIW2Stp6XmpzYbwAlizXMwwOHB1lzwtHv/Fad3vOAwD2+vGhIHB3VL8lRTv2oVx9ZsG2zwmHQofVcmWoo7QJMpaw2ACeD/UId9L6z8BwQG7Lj8o4mQAAAAASUVORK5CYII=",
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
                    id: "blink",
                    name: "component.public.animateBlink",
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
                  "w": 270,
                  "h": 137
                },
                "visible":1,
                "zIndex": -1,
                "backColor": "transparent",
                "foreColor": "#000000",
                fontSize: 12,
                "transform": 0,
                "diy":[
                  {
                    "name":"configComponent.weather.style",
                    "type":6,
                    "enumList":[
                      {option:'configComponent.weather.style1',value:1},
                      {option:'configComponent.weather.style2',value:2},
                      {option:'configComponent.weather.style3',value:3},
                      {option:'configComponent.weather.style4',value:4},
                      {option:'configComponent.weather.style5',value:5},
                      {option:'configComponent.weather.style6',value:6},
                    ],
                    "value":1,
                    "key":"style",
                  },
                  {
                    "name":"configComponent.weather.City",
                    "type":4,
                    "value":"阜阳",
                    "key":"City",
                  },
                ]
              }
            }
          }
        }
    },
    methods: {
        initComponents(option){
          this.width = option.style.position.w
          this.height = option.style.position.h
          let i=0
          for(i=0;i<option.style.diy.length;i++)
          {
            if(option.style.diy[i].key=="style")
            {
              if(option.style.diy[i].value==1)
              {
                this.weatherStyle = "ya"
              }
              else if(option.style.diy[i].value==2)
              {
                this.weatherStyle = "tz"
              }
              else if(option.style.diy[i].value==3)
              {
                this.weatherStyle = "tx"
              }
              else if(option.style.diy[i].value==4)
              {
                this.weatherStyle = "tg"
              }
              else if(option.style.diy[i].value==5)
              {
                this.weatherStyle = "tk"
              }
              else if(option.style.diy[i].value==6)
              {
                this.weatherStyle = "tt"
              }
            }
            else if(option.style.diy[i].key=="City")
            {
              if(option.style.diy[i].value!="") {
                this.city = option.style.diy[i].value.replace("市", "")
                this.city = this.city.replace("县", "")
                this.city = this.city.replace("区", "")
              }
            }
          }
          let  foreColor = option.style.foreColor.replace("#", "")
          this.Url = "https://widget.tianqiapi.com/?style="+this.weatherStyle+"&skin=pitaya&city="+this.city+"&fontsize="+option.style.fontSize+"&color="+foreColor
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
    this.GetNodeObj = this.getNode()
    this.GetNodeObj.on('change:data', ({ current }) => {
      if(current) {
        _t.detail = current.detail
      }
    })
    this.GetNodeObj.on('change:size', ({ current }) => {
      _t.detail.style.position.w = current.width
      _t.detail.style.position.h = current.height
    });
    this.detail = this.GetNodeObj.getData().detail
    this.editMode = this.GetNodeObj.getData().editMode
    this.showDeviceUuid = this.GetNodeObj.getData().showDeviceUuid
    this.IsToolBox = this.GetNodeObj.getData().IsToolBox
    _t.$EventBus.$on('cell-editMode', (data) => {
      _t.editMode = data.edit
      _t.IsToolBox = data.toolbox
	  _t.initComponents(_t.detail)
    })
    this.initComponents(this.detail);
  }
}
</script>

<style lang="less">
.view-url {
    position:absolute;
    z-index:-1;
}
</style>
