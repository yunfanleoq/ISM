<template>

  <div :style="animatedStyle" v-show="detail.style.visible==1||isStart ? true:false">
    <div :class="{
          'animated':true,[`${detail.style.animate}`]: true,
          'color-animation':isStart&&animateType.includes('millcolorGrad')&&!IsToolBox,
          'blink-animation':isStart&&animateType.includes('blink')&&!IsToolBox,
          'scale-animation':isStart&&animateType.includes('Zoom')&&!IsToolBox,
          'rotate-animation':isStart&&animateType.includes('animateSpin')&&!IsToolBox&&spinDirection==0,
          'rotate-anti-animation':isStart&&animateType.includes('animateSpin')&&!IsToolBox&&spinDirection==1,

        }"
         :style="{
            'animation':'ant-line-forward 30s infinite linear',
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
    </div>
  </div>

</template>

<script>

import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
    name: 'view-ism-group-node',
    inject: ['getNode'],
    data() {
      return {
        detail:null,
        IsToolBox:false,
        editMode:true,
        Text:"",
        DivOpacity:1,
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
        Variable:"###",
        ChartUnit:"",
        ShowJinZhi:5,
        GetNodeObj:null,
        base:{
          text: "configComponent.variable.Text",
          "icon": "icon-bianliangbiao",
          "isFontIcon": true,
          "info": {
            "type": "image",
            "action": [],
            "dataBind":[],
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
                "w": 132,
                "h": 60
              },
              "visible":1,
              "backColor": "transparent",
              "foreColor": "#000000",
              "zIndex": -1,
              "transform": 0,
              fontWeight:400,
              textAlign: "center",
              fontSize: 74,
              fontFamily: "黑体",
              "diy":[
                {
                  "name":"component.public.fillOpacity",
                  "type":7,
                  "value":1,
                  "min":0,
                  "max":1,
                  "key":"fillOpacity",
                },
                {
                  "name":"configComponent.ChartPublic.ChartUnit",
                  "type":4,
                  "value":"",
                  "key":"ChartUnit",
                },
                {
                  "name":"configComponent.ChartPublic.ShowJinZhi",
                  "type":6,
                  "enumList":[
                    {option:'configComponent.ChartPublic.ShowJinZhiNo',value:5},
                    {option:'configComponent.ChartPublic.ShowJinZhi10',value:1},
                    {option:'configComponent.ChartPublic.ShowJinZhi16',value:2},
                    {option:'configComponent.ChartPublic.ShowJinZhi8',value:3},
                    {option:'configComponent.ChartPublic.ShowJinZhi2',value:4}
                  ],
                  "value":5,
                  "key":"ShowJinZhi",
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
      ensureGroupControlStatus(option){
        if (!option) {
          return
        }
        if (!Array.isArray(option.active)) {
          option.active = []
        }
        const hasStatus = option.active.some(item => item && item.id === 'ControlStatus')
        if (hasStatus) {
          return
        }
        option.active.push({
          id:"ControlStatus",
          name:"configComponent.status.ControlStatus",
          result:0,
          isStatus:true,
          isSwitch:false,
          isImageStatus:false,
          isTextStatus:false,
          isLineStatus:true,
          isExpression:false,
          condition:{
            deviceSN:"",
            isBandDevice:false,
            bandType:1,
            dataID: "",
            dataName: "",
            IsManual:false,
            StatusList:[
              {
                "StatusOpt":"==",
                "TextColor":"#d81e06",
                "Blink":'0',
                "BlinkSpeed":1,
                "value2":1,
                "value":1
              }
            ]
          },
        })
      },
      applyGroupStatusToChildren(color){
        const parent = this.GetNodeObj
        if (!parent || typeof parent.getChildren !== 'function' || !color) {
          return
        }
        parent.getChildren().forEach(child => {
          const data = child.getData() || {}
          const detail = data.detail
          if (!detail || !detail.style) {
            return
          }
          if (Array.isArray(detail.style.diy)) {
            detail.style.diy.forEach(item => {
              if (item && (item.key === 'strokeColor' || item.key === 'strokeFill')) {
                item.value = color
              }
            })
          }
          child.setData(Object.assign({}, data, { detail: detail, UpdateNodeFlag: Date.now() }))
          if (detail.identifier) {
            this.$EventBus.$emit(detail.identifier + 'activeEvent', { ID: 'GroupStrokeColor', result: color })
          }
        })
      },
      initComponents(option){
        if(this.IsToolBox)
        {
          return
        }
        this.DivOpacity = option.style.opacity
        let i=0
        const diy = (option.style && option.style.diy) || []
        for( i=0;i<diy.length;i++)
        {
          if(diy[i].key=="strokeWidth")
          {
            this.strokeWidth=diy[i].value
          }
          else if(diy[i].key=="strokeFill")
          {
            this.fill=diy[i].value
          }
          else if(diy[i].key=="strokeColor")
          {
            this.strokeColor=diy[i].value
          }
          else if(diy[i].key=="fillOpacity")
          {
            this.fillOpacity=diy[i].value
          }
          else if(diy[i].key=="strokeOpacity")
          {
            this.strokeOpacity=diy[i].value
          }
          else if(diy[i].key=="imageURL")
          {
            this.imageURL=diy[i].value
          }
          else if(diy[i].key=="ChartUnit")
          {
            this.ChartUnit=diy[i].value
          }
          else if(diy[i].key=="ShowJinZhi")
          {
            this.ShowJinZhi=diy[i].value
          }
        }
        i=0
        this.animateType = (option.animate && option.animate.selected) || []
        this.ensureGroupControlStatus(option)
        if(option.animate && option.animate.isExpression)
        {
          this.isStart = false
        }
        else
        {
          this.isStart = true
        }
        const animateElement = (option.animate && option.animate.animateElement) || []
        for( i=0;i<animateElement.length;i++)
        {
          if(animateElement[i].id=="millcolorGrad")
          {
            for(let k =0;k<animateElement[i].elementList.length;k++)
            {
              if(animateElement[i].elementList[k].key=="startColor")
              {
                this.startColor=animateElement[i].elementList[k].value
              }
              else if(animateElement[i].elementList[k].key=="stopColor")
              {
                this.stopColor=animateElement[i].elementList[k].value
              }
              else if(animateElement[i].elementList[k].key=="animateSpeed")
              {
                this.animateSpeed=animateElement[i].elementList[k].value
              }
            }
          }
          else if(animateElement[i].id=="blink")
          {
            for(let k =0;k<animateElement[i].elementList.length;k++) {
              if (animateElement[i].elementList[k].key == "blinkSpeed") {
                this.blinkSpeed = parseFloat(animateElement[i].elementList[k].value)
              }
            }
          }
          else if(animateElement[i].id=="animateSpin")
          {
            for(let k =0;k<animateElement[i].elementList.length;k++) {
              if (animateElement[i].elementList[k].key == "spinSpeed") {
                this.animateSpinSpeed = animateElement[i].elementList[k].value
              }
              else if (animateElement[i].elementList[k].key == "spinDirection") {
                this.spinDirection = animateElement[i].elementList[k].value
              }
            }
          }
        }
      }
    },
    mounted() {
      let _t = this
      // 保存闭包引用，用于 beforeDestroy 精确 $off
      this._handlers = {}
      this.$nextTick(function(){
        this.initComponents(this.detail);
        let activeEvent = this.detail.identifier+"activeEvent"//动作数据
        let animateEvent = this.detail.identifier+"animateEvent"//动作数据

        _t._handlers.activeEvent = (data) => {
          if((_t.editMode)&&(!this.IsToolBox)){
            return
          }
          if(data.ID == "ShowVariable")
          {
            if(_t.ShowJinZhi!=5)
            {
              let tenNumber = parseInt(data.result)
              if(_t.ShowJinZhi==1)
              {
                _t.Variable = tenNumber.toString()
              }
              else if(_t.ShowJinZhi==2)
              {
                _t.Variable = tenNumber.toString(16).toUpperCase()
              }
              else if(_t.ShowJinZhi==3)
              {
                _t.Variable = tenNumber.toString(8)
              }
              else if(_t.ShowJinZhi==4)
              {
                _t.Variable = tenNumber.toString(2)
              }
            }
            else
            {
              _t.Variable = data.result
            }
          }
          else if(data.ID == "ControlStatus")
          {
            const statusValue = parseFloat(data.result)
            const activeList = (_t.detail && _t.detail.active) || []
            const statusActive = activeList.find(item => item && item.id === 'ControlStatus') || activeList[0]
            const statusList = (statusActive && statusActive.condition && statusActive.condition.StatusList) || []
            for(let i=0;i<statusList.length;i++)
            {
              const rule = statusList[i]
              let matched = false
              switch(rule.StatusOpt)
              {
                case "==": matched = statusValue == rule.value; break
                case ">": matched = statusValue > rule.value; break
                case ">=": matched = statusValue >= rule.value; break
                case "<": matched = statusValue < rule.value; break
                case "<=": matched = statusValue <= rule.value; break
                case "!=": matched = statusValue != rule.value; break
                default: matched = statusValue == rule.value
              }
              if(matched)
              {
                _t.applyGroupStatusToChildren(rule.TextColor)
                break
              }
            }
          }
        }
        _t.$EventBus.$on(activeEvent, _t._handlers.activeEvent)
        _t._handlers.animateEvent = (data) => {
          _t.isStart = data
          if (_t.animateType.includes("visible")) {
            if (data) {
              _t.detail.style.visible = true
            } else {
              _t.detail.style.visible = false
            }
          }
        }
        _t.$EventBus.$on(animateEvent, _t._handlers.animateEvent)
      });
    },
    beforeDestroy() {
      // 清理 X6 Node 事件监听
      if (this.GetNodeObj) {
        this.GetNodeObj.off('change:data')
        this.GetNodeObj.off('change:size')
        this.GetNodeObj = null
      }
      // 清理 EventBus 监听
      const h = this._handlers || {}
      if (this.detail && this.detail.identifier) {
        this.$EventBus.$off(this.detail.identifier + "activeEvent", h.activeEvent)
        this.$EventBus.$off(this.detail.identifier + "animateEvent", h.animateEvent)
      }
      this.$EventBus.$off('cell-editMode', h.cellEditMode)
      this._handlers = null
    },
    created(){
    let _t = this
    this._handlers = this._handlers || {}
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
    this._handlers.cellEditMode = (data) => {
      _t.editMode = data.edit
      _t.IsToolBox = data.toolbox
      _t.initComponents(_t.detail)
    }
    _t.$EventBus.$on('cell-editMode', _t._handlers.cellEditMode)
    this.initComponents(this.detail);
  }
}
</script>
<style>
@keyframes ant-line-forward {
  to { stroke-dashoffset: -1000; }
}
@keyframes ant-line-inverse {
  to { stroke-dashoffset: 1000; }
}
</style>
<style scoped lang="less">
.svg-el {
  /*transform: rotate(45deg);*/
  transform-origin: center center;
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

