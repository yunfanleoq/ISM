<template>

  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"  x="0px" y="0px"  xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el" :class="{'animated':true,[`${detail.style.animate}`]: true}" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <foreignObject style="overflow:visible;" pointer-events="all" :width="detail.style.position.w" :height="detail.style.position.h">
        <a-empty v-if="!isShowDocx"></a-empty>
        <VueOfficeDocx v-else :src="docx" :style="{width:detail.style.position.w+'px',height:detail.style.position.h+'px'}" @rendered="rendered" />
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
import VueOfficeDocx from '@vue-office/docx'
import '@vue-office/docx/lib/index.css'

import BaseView from '../View';
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
export default {
  mixins: [ISMChildAutoMixin],
    name: 'ism-view-docx',
    inject: ['getNode'],
    props: {

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
    components: {
      VueOfficeDocx
    },
    computed: {
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
    data() {
        return {
          detail:null,
          IsToolBox:false,
          editMode:true,
          Url:"",
          width:600,
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
          isShowDocx: true,
          docx: 'https://501351981.github.io/vue-office/examples/dist/static/test-files/test.docx',
          editorOptions: {
            toolbar: true,
                width: '800px',
                height: '600px'
          },
          base:{
            "text": "configComponent.viewDoc.title",
            "icon": "icon-word",
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
                  "w": 200,
                  "h": 200
                },
                "visible":1,
                "zIndex": -1,
                "transform": 0,
                "diy":[
                  {
                    "name":"configComponent.viewPdf.doc",
                    "type":12,
                    "value":"",
                    "key":"docUrl",
                  },
                ]
              }
            }
          }
        }
    },
    methods: {
      initComponents(option){
          let that = this
          that.isShowDocx=false

          this.width = option.style.position.w
          this.height = option.style.position.h
          let i=0
          for(i=0;i<option.style.diy.length;i++)
          {
            if(option.style.diy[i].key=="docUrl")
            {
              if(option.style.diy[i].value!="") {
                this.docx = "http://" + location.host + "/" + option.style.diy[i].value
                // this.docx="http://127.0.0.1:8081/"+option.style.diy[i].value
                console.log(this.docx)
                setTimeout(function () {
                  that.isShowDocx = true
                }, 500)
              }
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

        },
      rendered(){
        console.log("渲染完成")
      }
    },
    beforeDestroy() {

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
//.vue-office-pdf-wrapper{
//  height:auto;
//  width:auto;
//}
//.vue-office-pdf-wrapper canvas{
//  height:100%;
//  width:100%
//}
</style>
