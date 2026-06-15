<template>
  <div :style="animatedStyle" v-show="detail.style.visible==1||isStart ? true:false">
    <div :class="{
          'animated':true,[`${detail.style.animate}`]: true,
          'color-animation':isStart&&animateType.includes('millcolorGrad')&&!IsToolBox,
          'blink-animation':isStart&&animateType.includes('blink')&&!IsToolBox,
          'scale-animation':isStart&&animateType.includes('Zoom')&&!IsToolBox,
          'rotate-animation':isStart&&animateType.includes('animateSpin')&&!IsToolBox&&spinDirection==0,
          'rotate-anti-animation':isStart&&animateType.includes('animateSpin')&&!IsToolBox&&spinDirection==1
        }" :style="{
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
            <div class="example">
                     <a-spin tip="Loading..." :spinning="isCharge">
                              <vue3dLoader v-if="!isCharge"
                                  :backgroundAlpha="0"
                                  :width="detail.style.position.w"
                                  :height="detail.style.position.h"
                                  :filePath="filePath"
                                  :mtlPath="mtlPath"
                                  :parallelLoad="true"
                                  :rotation="rotation"
                                  :cameraPosition="positon"
                                  @process="onLoadProcess"
                              ></vue3dLoader>
                       <div class="process" v-if="loadProcess!=100">
                         loadding: {{ loadProcess + "%" }}
                        </div>
                     </a-spin>

            </div>
    </div>
  </div>
</template>

<script>
import svgView from '../View';
import  vue3dLoader  from "@/components/3dLoader/vue3dLoader";
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
export default {
  mixins: [ISMChildAutoMixin],
    name: 'view-3d-model',
    inject: ['getNode'],
    data() {
      return {
        DivOpacity:1,
        rotation: {
          x: 0,
          y:0,
          z: 0,
        },
        positon: {
          x: 0,
          y:0,
          z: 0,
        },
        ShowLabel: null,
        strokeColor:"#000000",
        fill:"#A1BFE2",
        strokeWidth:0.3,
        fillOpacity:1,
        filePath:"",
        mtlPath:"",
        strokeOpacity:1,
        animateType:"blink",
        startColor:"#74f808",
        stopColor:"#74f808",
        animateSpeed:0.5,
        animateSpinSpeed:0.5,
        spinDirection:0,
        blinkSpeed:0.5,
        isStart:false,
        imageURL:"",
        isCharge:false,
        loadProcess:0,
        rotationAnimationFrame:null,
        detail:null,
        IsToolBox:false,
        editMode:true,
        base:{
          text: "configComponent.View3DModel.Text",
          "icon": "icon-a-ziyuan17",
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
                "w": 100,
                "h": 100
              },
              "visible":1,
              "backColor": "transparent",
              "foreColor": "#000000",
              "zIndex": -1,
              "transform": 0,
              text: "Test",
              fontWeight:400,
              textAlign: "center",
              fontSize: 30,
              fontFamily: "Arial",
              "diy":[
                {
                  "name":"configComponent.View3DModel.ModelFile",
                  "type":10,
                  "value":"/models/gltf/DamagedHelmet.gltf",
                  "key":"ModelFile",
                },
                {
                  "name":"configComponent.View3DModel.ModelTml",
                  "type":10,
                  "value":"",
                  "key":"ModelTml",
                },
                {
                  "name":"configComponent.View3DModel.Rotation",
                  "type":6,
                  "enumList":[
                    {option:'True',value:1},{option:'False',value:0}
                  ],
                  "value":0,
                  "key":"Rotation",
                },
                {
                  "name":"configComponent.View3DModel.RotationPostion",
                  "type":6,
                  "enumList":[
                    {option:'X轴',value:1},{option:'Y轴',value:0},{option:'Z轴',value:2}
                  ],
                  "value":0,
                  "key":"RotationPostion",
                },
                {
                  "name":"configComponent.View3DModel.CameraPositionX",
                  "type":7,
                  "value":0,
                  "key":"CameraPositionX",
                },
                {
                  "name":"configComponent.View3DModel.CameraPositionY",
                  "type":7,
                  "value":0,
                  "key":"CameraPositionY",
                },
                {
                  "name":"configComponent.View3DModel.CameraPositionZ",
                  "type":7,
                  "value":0,
                  "key":"CameraPositionZ",
                },
              ]
            }
          }
        }
      }
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
    components: {vue3dLoader  },
    watch: {
      detail: {
        handler(newVal, oldVal) {
          let _this = this
          if(this.editMode) {
            this.initComponents(newVal);
          }
        },
        deep: true
      },
      '$route' () {
        cancelAnimationFrame(this.rotationAnimationFrame)
      }
    },
    methods: {
      rotate() {
        this.rotationAnimationFrame = requestAnimationFrame(this.rotate);
        if(this.RotationPostion==0)
        {
          this.rotation.y += 0.01;
        }
        else if(this.RotationPostion==1)
        {
          this.rotation.x += 0.01;
        }
        else if(this.RotationPostion==2)
        {
          this.rotation.z += 0.01;
        }
      },
      initComponents(option){
        this.isCharge=true;
        let _t = this
        if(this.IsToolBox)
        {
          return
        }
        let i=0
        for( i=0;i<option.style.diy.length;i++)
        {
          if(option.style.diy[i].key=="ModelFile")
          {
              this.filePath = option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="ModelTml")
          {
            this.mtlPath = option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="Rotation")
          {
            this.Rotation = option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="RotationPostion")
          {
            this.RotationPostion = option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="CameraPositionX")
          {
            this.positon.x = option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="CameraPositionY")
          {
            this.positon.y = option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="CameraPositionZ")
          {
            this.positon.z = option.style.diy[i].value
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
        if(this.Rotation)
        {
          cancelAnimationFrame(this.rotationAnimationFrame)
          this.rotate()
        }
        else
        {
          cancelAnimationFrame(this.rotationAnimationFrame)
        }
        setTimeout(function (){
          _t.isCharge=false;
        },100)
      },
      onLoadProcess(event,index){
        this.loadProcess = Math.floor((event.loaded / event.total) * 100);
      },
    },
    mounted() {
      let _t = this
      this.$nextTick(function(){
       this.ShowLabel = {}
       this.ShowLabel.value = [
          // text label
          {
            text: "I'm Text Label",
            textStyle: {
              fontFamily: "Arial",
              fontSize: 12,
              fontWeight: 600,
              lineHeight: 1,
              color: "#ffffff",
              borderWidth: 8,
              borderRadius: 0,
              borderColor: "#000000",
              backgroundColor: "rgba(0,0,0,1)",
            },
            sid: 3, // 自定义标识，可有可无
          },
        ];
       this.initComponents(this.detail);

        let activeEvent = this.detail.identifier+"activeEvent"//动作数据
        let animateEvent = this.detail.identifier+"animateEvent"//动作数据

        _t.$EventBus.$on(activeEvent, (data) => {
          if((_t.editMode)&&(!this.IsToolBox)){
            return
          }

        })
        _t.$EventBus.$on(animateEvent, (data) => {
          if((_t.editMode)&&(!this.IsToolBox)){
            return
          }
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
      })
      this.initComponents(this.detail);
    },
    beforeDestroy() {
      cancelAnimationFrame(this.rotationAnimationFrame)
      this.rotationAnimationFrame = null
    }
}
</script>
<style lang="less">
.process {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translateX(-50%);
  padding: 3px 8px;
  background-color: aquamarine;
  color: red;
}

/* 使用animation关键帧 */
.color-animation {
  animation: colorChange var(--animateSpeed) linear infinite;
}

@keyframes colorChange {
  0% { background-color: var(--startColor); }
  100% { background-color:  var(--stopColor); }
}
/* 使用animation关键帧 */
.blink-animation {
  animation: blink var(--blinkSpeed) linear infinite;
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0; }
}
/*缩放*/
.scale-animation {
  animation: pulse 0.6s infinite;
}

@keyframes pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.5); }
}
/*顺时针旋转*/
.rotate-animation {
  animation: clockwiseRotate var(--animateSpinSpeed) linear infinite;
  transform-origin: center;
}

@keyframes clockwiseRotate {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
/*逆时针旋转*/
.rotate-anti-animation {
  animation: counterClockwiseRotate var(--animateSpinSpeed) linear infinite;
  transform-origin: center;
}

@keyframes counterClockwiseRotate {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(-360deg); }
}
.svg-el {
  /*transform: rotate(45deg);*/
  transform-origin: center center;
}
.example {
  text-align: center;
}
</style>
