<template>
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
            opacity:DivOpacity,
            borderWidth: detail.style.borderWidth + 'px',
            borderStyle: detail.style.borderStyle,
            borderColor: detail.style.borderColor,
            transform: detail.style.transform==-1099? 'rotateY(180deg)':detail.style.transform==-1098?'rotateX(180deg)':'',
        }">
        <div >
            <carousel v-if="isUpdate" :center="true" :autoplaySpeed="500" :autoplay="true" :nav="false" :loop="true" :items="1" :style="{'height':detail.style.position.h+'px','width':detail.style.position.w+'px','overflow': 'hidden'}">
              <img :style="{'height':(detail.style.position.h-30)+'px','width':(detail.style.position.w-30)+'px','overflow': 'hidden'}" v-for="item in ImgList" :src="item" :key="item"/>
            </carousel>
          </div>
    </div>
  </div>

</template>

<script>
import carousel from 'vue-owl-carousel'
import {mapState} from "vuex";
import store from "@/store";
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
export default {
  mixins: [ISMChildAutoMixin],
  name: 'view-slide',
  inject: ['getNode'],
  i18n: require('../../../../i18n/language'),
  data() {
    return {
      detail:null,
      IsToolBox:false,
      editMode:true,
      DivOpacity:1,
      isUpdate:false,
      dataComplete: false, // 数据加载完成标志
      reflush: true, // 强制刷新表格
      isShowNoData: false,
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
      imageURL:"",
      AlarmTimer:null,
      speeds:10,
      IsShowIndex:1,
      RealDataTableHeaderHeight:50,
      tableRows:5,
      direction:"left",
      tableCols:5,
      Variable:"###",
      ChartUnit:"",
      ImgList:[],
      tableBorderColor:"#b3b3b3",
      base:{
        text: "configComponent.slide.title",
        "icon": "icon-zujian-huandengpian",
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
              "w": 210,
              "h": 140
            },
            "visible":1,
            "backColor": "transparent",
            "zIndex": -1,
            "transform": 0,
            "diy":[
              {
                "name":"displayConfig.Properties.ImageUrl",
                "type":5,
                "value":"data:image/png;base64,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",
                "key":"imageURL",
              },
              {
                "name":"displayConfig.Properties.ImageUrl",
                "type":5,
                "value":"data:image/png;base64,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",
                "key":"imageURL1",
              },
              {
                "name":"displayConfig.Properties.ImageUrl",
                "type":5,
                "value":"data:image/png;base64,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",
                "key":"imageURL2",
              },
              {
                "name":"displayConfig.Properties.ImageUrl",
                "type":5,
                "value":"",
                "key":"imageURL3",
              },
              {
                "name":"displayConfig.Properties.ImageUrl",
                "type":5,
                "value":"",
                "key":"imageURL4",
              },
              {
                "name":"displayConfig.Properties.ImageUrl",
                "type":5,
                "value":"",
                "key":"imageURL5",
              },
              {
                "name":"displayConfig.Properties.ImageUrl",
                "type":5,
                "value":"",
                "key":"imageURL6",
              },
              {
                "name":"displayConfig.Properties.ImageUrl",
                "type":5,
                "value":"",
                "key":"imageURL7",
              },
              {
                "name":"displayConfig.Properties.ImageUrl",
                "type":5,
                "value":"",
                "key":"imageURL8",
              },
              {
                "name":"displayConfig.Properties.ImageUrl",
                "type":5,
                "value":"",
                "key":"imageURL9",
              }
            ]
          }
        }
      }
    }
  },
  components: {
    carousel
  },
  computed: {
    styleVar() {
      return {
        "--tableBorderColor": this.tableBorderColor ,
      };
    },
    ...mapState({
      ISMCavasContainer:state => store.state.ISMDisPlayEditorTool.ISMCavasContainer,
    }),
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
  methods: {
    initComponents(option){
      let that = this
      if(this.IsToolBox)
      {
        return
      }
      this.ImgList=[]
      this.isUpdate=false
      let i=0
      for( i=0;i<option.style.diy.length;i++)
      {
        if(option.style.diy[i].key=="imageURL")
        {
          if(option.style.diy[i].value!="")
          {
            this.ImgList.push(option.style.diy[i].value)
          }
        }
        else if(option.style.diy[i].key=="imageURL1")
        {
          if(option.style.diy[i].value!="")
          {
            this.ImgList.push(option.style.diy[i].value)
          }
        }
        else if(option.style.diy[i].key=="imageURL2")
        {
          if(option.style.diy[i].value!="")
          {
            this.ImgList.push(option.style.diy[i].value)
          }
        }
        else if(option.style.diy[i].key=="imageURL3")
        {
          if(option.style.diy[i].value!="")
          {
            this.ImgList.push(option.style.diy[i].value)
          }
        }
        else if(option.style.diy[i].key=="imageURL4")
        {
          if(option.style.diy[i].value!="")
          {
            this.ImgList.push(option.style.diy[i].value)
          }
        }
        else if(option.style.diy[i].key=="imageURL5")
        {
          if(option.style.diy[i].value!="")
          {
            this.ImgList.push(option.style.diy[i].value)
          }
        }
        else if(option.style.diy[i].key=="imageURL6")
        {
          if(option.style.diy[i].value!="")
          {
            this.ImgList.push(option.style.diy[i].value)
          }
        }
        else if(option.style.diy[i].key=="imageURL7")
        {
          if(option.style.diy[i].value!="")
          {
            this.ImgList.push(option.style.diy[i].value)
          }
        }
        else if(option.style.diy[i].key=="imageURL8")
        {
          if(option.style.diy[i].value!="")
          {
            this.ImgList.push(option.style.diy[i].value)
          }
        }
        else if(option.style.diy[i].key=="imageURL9")
        {
          if(option.style.diy[i].value!="")
          {
            this.ImgList.push(option.style.diy[i].value)
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
      setTimeout(function (){
        that.isUpdate=true
      },1000)

    },
  },
  beforeDestroy(){

  },
  mounted() {
    let _t = this
    this.$nextTick(function(){
      this.initComponents(this.detail);
      let activeEvent = this.detail.identifier+"activeEvent"//动作数据
      let animateEvent = this.detail.identifier+"animateEvent"//动作数据

      _t.$EventBus.$on(activeEvent, (data) => {
        if((_t.editMode)&&(!this.IsToolBox)){
          return
        }
        if(data.ID == "ShowVariable")
        {
          _t.Variable = data.result
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
    })
  }
}
</script>
<style scoped lang="less">
.svg-el {
  /*transform: rotate(45deg);*/
  transform-origin: center center;
}
.marquee-container{
  overflow: hidden;
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
</style>
