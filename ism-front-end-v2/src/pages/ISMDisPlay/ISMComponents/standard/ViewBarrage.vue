<template>
  <div :style="animatedStyle" v-show="detail.style.visible==1||isStart ? true:false">
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
      <div >
        <vue-marquee v-if="isUpdate" :showProgress="false" :duration="speeds" :direction="direction" :style="{'height':detail.style.position.h+'px','width':detail.style.position.w+'px','overflow': 'hidden'}" >
          <vue-marquee-slide v-for="i in 10" :key="i" >
                <span :style="{
                        fontSize: detail.style.fontSize + 'px',
                        fontFamily: detail.style.fontFamily,
                        'font-weight':detail.style.fontWeight,
                        color: detail.style.foreColor,
                        textAlign: textAlign,
                        'font-style':detail.style.italic?'oblique':'normal',
                        lineHeight: lineHeight + 'px',
                        padding:'0 15px',
                        whiteSpace: 'nowrap'
                    }"
                > {{detail.style.text}}{{Variable}}{{ChartUnit}} </span>
          </vue-marquee-slide>
        </vue-marquee>
      </div>

    </div>
  </div>
</template>

<script>
import { Marquee, Slide } from 'vue-marquee-component-fix';
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
export default {
  mixins: [ISMChildAutoMixin],
  name: 'view-barrage',
  i18n: require('../../../../i18n/language'),
  inject: ['getNode'],
  data() {
    return {
      detail:null,
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
      Variable:"",
      ChartUnit:"",
      tableBorderColor:"#b3b3b3",
      base:{
        text: "configComponent.barrage.title",
        "icon": "icon-m-gundongwenzi",
        "isFontIcon": true,
        "info": {
          "type": "image",
          "action": [],
          "active": [
            {
              id:"ShowVariable",
              name:"configComponent.variable.ShowData",
              result:"",
              isExpression:false,
              condition:{
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
            }
          ],
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
            "foreColor": "#000000",
            fontWeight:400,
            "zIndex": -1,
            "transform": 0,
            text: "Test",
            textAlign: "center",
            fontSize: 30,
            fontFamily: "Arial",
            italic:0,
            "diy":[
              {
                "name":"configComponent.ChartPublic.ChartUnit",
                "type":4,
                "value":"",
                "key":"ChartUnit",
              },
              {
                "name":"configComponent.barrage.speeds",
                "type":1,
                "value":56000,
                "min":1,
                "key":"speeds",
              },
              {
                name:"configComponent.barrage.direction",
                type:6,
                value:0,
                enumList:[
                  {
                    value:0,
                    option:"configComponent.barrage.directionLeft"
                  },
                  {
                    value:1,
                    option:"configComponent.barrage.directionRight"
                  },
                  {
                    value:2,
                    option:"configComponent.barrage.directionUp"
                  },
                  {
                    value:3,
                    option:"configComponent.barrage.directionDown"
                  }
                ],
                min:1,
                key:"direction",
              }
            ]
          }
        }
      }
    }
  },
  components: {
    [Marquee.name]: Marquee,
    [Slide.name]: Slide
  },
  computed: {
    animatedStyle(){
      return {
        "--tableBorderColor": this.tableBorderColor ,
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
    end() {
      console.log(this.$refs.danmaku);
    },
    initComponents(option){
      let that = this
      if(this.IsToolBox)
      {
        return
      }
      this.isUpdate=false
      let i=0
      for( i=0;i<option.style.diy.length;i++)
      {
        if(option.style.diy[i].key=="direction")
        {
          if(option.style.diy[i].value==0)
          {
            this.direction="left"
          }
          else if(option.style.diy[i].value==1)
          {
            this.direction="right"
          }
          else if(option.style.diy[i].value==2)
          {
            this.direction="top"
          }
          else if(option.style.diy[i].value==3)
          {
            this.direction="bottom"
          }
        }
        else if(option.style.diy[i].key=="speeds")
        {
          this.speeds=parseInt(option.style.diy[i].value)
        }
        else if(option.style.diy[i].key=="ChartUnit")
        {
          this.ChartUnit=option.style.diy[i].value
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
    let _t = this
    // 清理EventBus事件监听
    let activeEvent = this.detail.identifier+"activeEvent"
    let animateEvent = this.detail.identifier+"animateEvent"
    
    // 清理定时器
    if(this.AlarmTimer){
      clearInterval(this.AlarmTimer)
      this.AlarmTimer = null
    }
    
    // 清理节点事件监听
    const node = this.getNode()
    if(node){
      node.off('change:data')
      node.off('change:size')
    }
    
  },
  mounted() {
    let _t = this
    this.$nextTick(function(){
      this.initComponents(this.detail);
      let activeEvent = this.detail.identifier+"activeEvent"//动作数据
      let animateEvent = this.detail.identifier+"animateEvent"//动作数据

      _t.$EventBus.$on(activeEvent, (data) => {
        if(data.ID == "ShowVariable")
        {
          _t.Variable = data.result
        }
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
    })
    this.initComponents(this.detail);
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