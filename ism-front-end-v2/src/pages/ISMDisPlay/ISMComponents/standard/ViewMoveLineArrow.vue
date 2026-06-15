<template>
  <div
               class="view-line-arrow"
               :style="styleVar"
               @mousemove="onMousemove($event)" @mouseup="onMouseUp($event)">
              <canvas  :ref="detail.identifier" :width="detail.style.position.w" :height="detail.style.position.h"  :style="{
                                width: detail.style.position.w + 'px',
                                height: detail.style.position.h + 'px',
                            }">
      Your browser does not support the HTML5 canvas tag.
    </canvas>

             <div v-if="editMode && selected">
        <div v-for="(pass,index) in points"  :key="index">
            <div  @contextmenu.stop="PointRightClick($event,index)">
              <div   v-if="!pass.isArrow"  class="passby" @mousedown.stop="arrowPassDown(pass,$event,index)" :style="{
                              left: strokeWidth<10?pass.x-(strokeWidth)/2  + 'px':pass.x-(strokeWidth)/2-5  + 'px',
                              top: strokeWidth<10?pass.y-(strokeWidth)/2  + 'px':pass.y-(strokeWidth)/2 -5 + 'px',
                          }">

              </div>

              <div v-else    class="passbx" @mousedown.stop="arrowPassDown(pass,$event,index)" :style="{
                                left: strokeWidth<10?pass.x-(strokeWidth)/2  + 'px':pass.x-(strokeWidth)/2-5  + 'px',
                                top: strokeWidth<10?pass.y-(strokeWidth)/2 + 'px':pass.y-(strokeWidth)/2-5 + 'px',
                            }">

              </div>
            </div>
        </div>
    </div>
          </div>
</template>

<script>
import canvasView from '../View';
import store from "../../../../store";
import { mapActions, mapState, mapMutations } from 'vuex'
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
export default {
  mixins: [ISMChildAutoMixin],
    name: 'ViewCanvasMoveLineArrow',
    inject: ['getNode'],
    i18n: require('@/i18n/language'),
    data() {
        return {
          spinDirection:0,
          timer:null,
          DeletePointIndex:-1,
          identifier:"",
          width:0,
          strokeColor:"#000000",
          fill:"#A1BFE2",
          strokeWidth:0.3,
          fillOpacity:1,
          strokeOpacity:1,
          startColor:"#74f808",
          stopColor:"#74f808",
          animateSpeed:0.5,
          animateSpinSpeed:0.5,
          blinkSpeed:0.5,
          RightClick:false,
          movePoints:[],
          tempMoveHeight:0,
          tempMoveWidth:0,
          height:0,
          lineDashOffset:0,
          MoveBrokenLineInterval:"",
          flag: false,
          passItem: {},
          currentPoint:{},
          moveIndex :0,
          ProPoint:{},
          isStart:false,
          points: [], //控制点（包含起始和终点）
          FACTOR_H: 5, //箭头 水平高度倍数
          FACTOR_V: 4, //箭头 垂直长度倍数
          animateType:"blink",
          base:{
            "text": "displayConfig.ToolBox.Diagram.MoveBrokenLine",
            "icon": "icon-chushuiliuliangicon",
            "isFontIcon": true,
            "info": {
              "type": "image",
              isCanvas:true,
              "action": [],
              "dataBind":[],
              "active": [
                {
                  id:"Forward",
                  name:"component.ViewCanvasMoveLineArrow.Forward",
                  result:"",
                  isExpression:true,
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
                },
                {
                  id:"Reverse",
                  name:"component.ViewCanvasMoveLineArrow.Reverse",
                  result:"",
                  isExpression:true,
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
                },
              ],
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
                  "w": 600,
                  "h": 50
                },
                "points": [],
                "visible":1,
                "zIndex": -1,
                "transform": 0,
                "backColor": "transparent",
                foreColor:"#F1A94B",
                "diy":[
                  {
                    "name":"displayConfig.ToolBox.Diagram.MoveBrokenLineBackColor",
                    "type":2,
                    "value":"#6D7B92",
                    "key":"MoveBrokenLineBackColor",
                  },
                  {
                    "name":"component.public.strokeWidth",
                    "type":1,
                    "value":20,
                    "min":1,
                    "key":"strokeWidth",
                  },
                  {
                    name:"displayConfig.ToolBox.Diagram.MoveBrokenLineConditionEnable",
                    type:6,
                    value:0,
                    enumList:[
                      {
                        value:0,
                        option:"component.public.Forbidden"
                      },
                      {
                        value:1,
                        option:"component.public.Enable"
                      }
                    ],
                    min:1,
                    key:"MoveBrokenLineConditionEnable",
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
              }
            }
          }
        }
    },
    destroyed () {
      // 组件销毁，关闭定时执行
      cancelAnimationFrame(this.timer)
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
    computed: {
      ...mapState({
        selectedValue:state => store.state.ISMDisPlayEditorTool.selectedValue,
        ISMCavasContainer:state => store.state.ISMDisPlayEditorTool.ISMCavasContainer,
      }),
      styleVar() {
        let strokeWidth = 10
        if(this.strokeWidth<10)
        {
          strokeWidth = 10
        }
        else
        {
          strokeWidth = parseInt(this.strokeWidth)+10
        }
        return {
          "--arrowWidth": (parseInt(strokeWidth))+"px" ,
        };
      },
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
    methods: {
      PointRightClick($event,index){
        $event.preventDefault()
        this.flag = false
        this.DeletePointIndex =index
        let _t = this
        this.$contextmenu({
          items: [
            {
              label: _t.$t('displayConfig.Canvas.menu.Delete'),
              icon: "el-icon-shanchu",
              onClick: () => {
                if(this.points.length>3)
                {
                  let ArrowCount = 0
                  for(let i=0;i<this.points.length;i++)
                  {
                    if(this.points[i].isArrow)
                    {
                      ArrowCount++
                    }
                  }
                  if((!this.points[this.DeletePointIndex].isArrow))
                  {
                    if((this.points[this.DeletePointIndex-1].isArrow)&&(ArrowCount==1))
                    {
                      this.points.splice(this.DeletePointIndex,1)
                    }
                    else
                    {
                      this.points.splice(this.DeletePointIndex-1,2)
                    }
                  }
                  else
                  {
                    if(ArrowCount>1)
                    {
                      this.points.splice(this.DeletePointIndex,1)
                    }
                  }

                  this.reDraw()
                  this.DeletePointIndex=-1
                }
              }
            }
          ],
          event, // 鼠标事件信息
          divided:true,
          customClass: "custom-class", // 自定义菜单 class
          zIndex: 10000, // 菜单样式 z-index
          minWidth: 230 // 主菜单最小宽度
        });
      },
      drawArrow(ctx, x2, y2, lineWidth, color) { // (x1, y1)是线段起点  (x2, y2)是线段终点
          ctx.beginPath(); // 坐标原点 => (x2, y2)
          ctx.moveTo(x2, y2);
          ctx.lineTo(x2 - lineWidth * this.FACTOR_H, y2 - lineWidth * this.FACTOR_V);
          ctx.lineTo(x2 - lineWidth * this.FACTOR_H, y2 + lineWidth * this.FACTOR_V);
          ctx.closePath();
          ctx.fillStyle = color; //设置线的颜色状态
          ctx.fill();
      },
      drawDashLine(ctx) {
          ctx.save()
          ctx.beginPath();
          ctx.lineWidth = this.strokeWidth; //设置线宽状态
          ctx.strokeStyle = this.strokeColor; //设置线的颜色状态
          ctx.lineCap = "round"
          ctx.setLineDash([this.strokeWidth,parseInt(this.strokeWidth)+10]);//实线 去掉此句

          if(this.MoveBrokenLineConditionEnable)
          {
            if(this.Forward==true)
            {
              for (let index = 0; index < this.points.length; index++) {
                const begin = this.points[index],
                    end = this.points[index + 1];
                ctx.moveTo(begin.x, begin.y);
                ctx.lineTo(end.x, end.y);
                if (index == this.points.length - 2)
                  break;
              }
              ctx.lineDashOffset = -(this.lineDashOffset+parseInt(this.strokeWidth))
            }
            else if(this.Reverse==true)
            {
              for (let index = 0; index < this.points.length; index++) {
                const begin = this.points[index],
                    end = this.points[index + 1];
                ctx.moveTo(begin.x, begin.y);
                ctx.lineTo(end.x, end.y);
                if (index == this.points.length - 2)
                  break;
              }
              ctx.lineDashOffset = this.lineDashOffset+parseInt(this.strokeWidth)
            }
            else
            {
              ctx.lineDashOffset = 0
            }
          }
          else
          {
            for (let index = 0; index < this.points.length; index++) {
              const begin = this.points[index],
                  end = this.points[index + 1];
              ctx.moveTo(begin.x, begin.y);
              ctx.lineTo(end.x, end.y);
              if (index == this.points.length - 2)
                break;
            }
            if(this.spinDirection==0)
            {
              ctx.lineDashOffset = -(this.lineDashOffset+parseInt(this.strokeWidth))
            }
            else
            {
              ctx.lineDashOffset = this.lineDashOffset+parseInt(this.strokeWidth)
            }
          }

          ctx.stroke(); //进行绘制
          ctx.closePath();
        ctx.restore()
      },
      drawLine(ctx) {
        ctx.save()
        ctx.beginPath();
        ctx.lineWidth = parseInt(this.strokeWidth)+5//设置线宽状态
        ctx.strokeStyle = this.backColor; //设置线的颜色状态
        ctx.lineCap = "round"
        for (let index = 0; index < this.points.length; index++) {
          const begin = this.points[index],
              end = this.points[index + 1];
          ctx.moveTo(begin.x, begin.y);
          ctx.lineTo(end.x, end.y);
          if (index == this.points.length - 2)
            break;
        }
        ctx.stroke(); //进行绘制
        ctx.closePath();
        ctx.restore()
      },
      getRad(degree){
        return degree/180*Math.PI;
      },
      reDraw() {
        this.lineDashOffset++
        let w = this.detail.style.position.w;
        let h = this.detail.style.position.h;
        let refObj = this.detail.identifier
        let el = this.$refs[refObj];
        let ctx = el.getContext("2d");

        ctx.clearRect(0, 0, w, h);
        this.drawLine(ctx);
        this.drawDashLine(ctx);
        if(this.lineDashOffset>1000)
        {
          this.lineDashOffset =8
        }
      },
      onResize() {
          this.reDraw();
      },
      intervalTimer () {
        this.reDraw()
        this.timer = requestAnimationFrame(this.intervalTimer)
      },
      arrowPassDown(pass, event, index) {
          if(event.button!=0)
          {
            this.flag = false;
            return
          }
          this.flag = true;
          pass.startX = event.pageX;
          pass.startY = event.pageY;
          pass.temp = {};
          pass.temp.x = pass.x;
          pass.temp.y = pass.y;
          this.passItem = pass;
          this.moveIndex = index
          if (pass.isArrow)
          {
            this.ProPoint = this.points[index-1]
            if(typeof this.ProPoint=='undefined')
            {
              let x2 = pass.x-50
              let y2 = pass.y
              let point2 = {
                x:x2,
                y:y2,
                isArrow:false,
              }
              this.ProPoint = point2
            }
            this.currentPoint = this.points[index]
          }

      },
      onMousemove(event) {
        if (!this.flag) {
          return
        }
        let zoom = this.selectedValue/100
        let dx = (event.pageX - this.passItem.startX)/zoom,
            dy = (event.pageY - this.passItem.startY)/zoom;

        event.cancelBubble = true;
        this.passItem.x = this.passItem.temp.x + dx;
        this.passItem.y = this.passItem.temp.y + dy;
        this.$nextTick(function(){
          this.reDraw();
        })
      },
      onMouseUp(event) {
        if(event.button!=0)
        {
          this.flag = false;
          return
        }
          if (!this.flag)
          {
            return
          }
        let zoom = this.selectedValue/100
        let moveDx = (event.pageX - this.passItem.startX)/zoom,
            moveDy = (event.pageY - this.passItem.startY)/zoom;
          this.flag = false;
          let proPoint = this.ProPoint
          let pass = this.currentPoint
          let nextPoint = this.points[this.moveIndex+1]
          if(!pass.isArrow)
          {
            return
          }
          if(typeof nextPoint=='undefined')
          {
            let x2 = pass.x+50
            let y2 = pass.y
            let point2 = {
              x:x2,
              y:y2,
              isArrow:false,
            }
            nextPoint = point2
          }
          let dx = Math.abs(event.pageX - this.passItem.startX)
          if(dx<2)
          {
            return
          }
          let x1 = proPoint.x+(pass.x-proPoint.x)/2
          let y1 = proPoint.y+(pass.y-proPoint.y)/2
          if(x1<10)
          {
            return
          }
          let point1 = {
            x:x1,
            y:y1,
            isArrow:true,
          }
          let x2 = pass.x+(nextPoint.x-pass.x)/2
          let y2 = pass.y+(nextPoint.y-pass.y)/2
          let point2 = {
            x:x2,
            y:y2,
            isArrow:true,
          }

          this.points[this.moveIndex].isArrow = false
          this.points.splice(this.moveIndex, 0, point1)
          this.points.splice(this.moveIndex + 2, 0, point2)

        },
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
          else if(option.style.diy[i].key=="MoveBrokenLineBackColor")
          {
            this.backColor=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="MoveBrokenLineInterval")
          {
            this.MoveBrokenLineInterval=option.style.diy[i].value
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
          else if(option.style.diy[i].key=="spinDirection")
          {
            this.spinDirection=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="MoveBrokenLineConditionEnable")
          {
            this.MoveBrokenLineConditionEnable=option.style.diy[i].value
          }
        }
        this.strokeColor = option.style.foreColor
        this.animateType = option.animate.selected
        if(this.detail.style.points.length==0)
        {
          let points=[
            {
              "x": 50,
              "y": this.strokeWidth,
              "isArrow":false
            },
            {
              "x": 350,
              "y": this.strokeWidth,
              "isArrow":true
            },
            {
              "x": 580,
              "y": this.strokeWidth,
              "isArrow":false
            },
          ]
          this.detail.style.points = points;
        }
        this.points = this.detail.style.points;
        this.onResize();
      },
    },
    mounted() {
      let _t = this
      this.$nextTick(function(){
       this.intervalTimer()
        this.initComponents(this.detail);
          let activeEvent = this.detail.identifier+"activeEvent"//动作数据
          let animateEvent = this.detail.identifier+"animateEvent"//动作数据

          _t.$EventBus.$on(activeEvent, (data) => {
            if(data.ID == "Forward")
            {
              _t.Forward = data.result
            }
            else if(data.ID == "Reverse")
            {
              _t.Reverse = data.result
            }
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
    this.initComponents(this.detail);
  }
}
</script>

<style >
.mymenu {
  padding: 2px 0;
}
.mymenu .menu-active {
  border-radius: 0;
  border-color: transparent;
}
.passby {
  position: absolute;
  height: var(--arrowWidth);
  width: var(--arrowWidth);
  background: #0cf;
  border-radius: 50%;
  margin-left: -1px;
  left: 50%;
  cursor: crosshair;
}
.view-line-arrow {
    height: 100%;
    width: 100%;
    position: relative;
    overflow: visible ;
}
.passbx {
  position: absolute;
  height: var(--arrowWidth);
  width: var(--arrowWidth);
  cursor: move;
  background: #c20629;
  border-radius: 50%;
  top: -20px;
  margin-left: -1px;
  left: 50%;
  cursor: crosshair;
}
</style>
