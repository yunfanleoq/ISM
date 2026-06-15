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
            opacity:DivOpacity,
            borderWidth: detail.style.borderWidth + 'px',
            borderStyle: detail.style.borderStyle,
            borderColor: detail.style.borderColor,
            transform: detail.style.transform==-1099? 'rotateY(180deg)':detail.style.transform==-1098?'rotateX(180deg)':'',
        }">
   <a-input :readOnly="editMode" type="text" @input="onInput" @focus="onInputFocus()"  v-model="TextValue" :style="{'padding-left':leftDis+'px','border':0,'font-size': detail.style.fontSize+'px','border-radius':detail.style.BorderEdges+'px','outline':'none','color':detail.style.foreColor,'background-color':detail.style.backColor,'height': detail.style.position.h+'px','width':detail.style.position.w+'px'}"/>
  <SimpleKeyboard v-show="showKeyboardpassword&&virtuallyKey" ref="SimpleKeyboardloginPassword" @onChange="onChangeKeyboard" @onKeyPress="onKeyPress" keyboardClass="simple-keyboard-user"/>
</div>
  </div>

</template>

<script>
import SimpleKeyboard from '@/components/SimpleKeyboard/SimpleKeyboard.vue'
import {setData} from "@/services/device";
import {mapState} from "vuex";
import store from "@/store";
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
export default {
  mixins: [ISMChildAutoMixin],
    name: 'view-svg-input-text',
    inject: ['getNode'],
    i18n: require('@/i18n/language'),
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
        leftDis:10,
        showKeyboardpassword:false,
        virtuallyKey:false,
        isStart:false,
        italic:false,
        imageURL:"",
        TextValue:"Input",
        ChartUnit:"",
        ShowJinZhi:5,
        InputTimer:null,
        IsInput:false,
        base:{
          text: "configComponent.input.Text",
          "icon": "icon-tubiao-shuru_wenbenshurukuang",
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
            "style": {
              "position": {
                "x": 0,
                "y": 0,
                "w": 100,
                "h": 40
              },
              "visible":1,
              "backColor": "#ffffff",
              "foreColor": "#000000",
              "zIndex": -1,
              "transform": 0,
              fontSize: 30,
              letterSpacing:0,
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
                  "name":"component.public.leftDis",
                  "type":1,
                  "value":10,
                  "key":"leftDis",
                },
                {
                  "name":"displayConfig.Properties.virtuallyKey",
                  "type":6,
                  "enumList":[
                    {option:'True',value:1},{option:'False',value:0}
                  ],
                  "value":0,
                  "key":"virtuallyKey",
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
    components: {
      SimpleKeyboard
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
      onKeyPress(button){
        if (button === '{enter}' || button === '{close}') {
          this.closekeyboard()
        }
      },
      onInput(res) {
        let _t = this
        this.IsInput = true
        if(this.InputTimer!=null)
        {
          clearTimeout(this.InputTimer)
        }
        this.InputTimer = setTimeout(function (){
          _t.IsInput = false
          _t.InputTimer=null
          _t.ManualSetData()
        },1500)
      },
      onInputFocus(res) {

        this.showKeyboardpassword = true
        this.changeIpt = res
        // 父组件调用子组件的方法
        this.$refs.SimpleKeyboardloginPassword.onKeyPress('{clear}')
      },
      ManualSetData(){
        let _t = this

        if(this.TextValue==null||this.TextValue=="")
        {
          return
        }
        this.settingLoading=true
        let params = {
          deviceUuid:this.detail.active[0].condition.deviceSN,
          dataUuid:this.detail.active[0].condition.dataID,
          value:this.TextValue.toString(),
        };
        this.$message.destroy()
        setData(params).then(function (res){
          if(res.data.code!=0)
          {
            _t.$message.error(_t.$t("readData.SetFailed"))
          }
          else
          {
            _t.$message.success(_t.$t("readData.SetSuccess"))
          }
        }).catch(function (error) {
          _t.settingLoading = false
          _t.$message.error(_t.$t("readData.SetFailed"))
        }).finally(function (error) {
          _t.settingLoading = false
        })
      },
      // 给输入框赋值
      onChangeKeyboard(input) {
        this.TextValue = input
      },
      // 点击关闭隐藏键盘
      closekeyboard() {
        this.showKeyboardpassword = false
      },
      initComponents(option){
        if(this.IsToolBox)
        {
          return
        }
        let i=0
        for( i=0;i<option.style.diy.length;i++)
        {
          if(option.style.diy[i].key=="ChartUnit")
          {
            this.ChartUnit=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="ShowJinZhi")
          {
            this.ShowJinZhi=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="leftDis")
          {
            this.leftDis=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="virtuallyKey")
          {
            this.virtuallyKey=option.style.diy[i].value
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

        _t.$EventBus.$on(activeEvent, _t._activeEventHandler = (data) => {
          if((_t.editMode)&&(!this.IsToolBox)){
            return
          }
          if(_t.IsInput)
          {
            return
          }
          if(data.ID == "ShowVariable")
          {
            if(_t.ShowJinZhi!=5)
            {
              let tenNumber = parseInt(data.result)
              if(_t.ShowJinZhi==1)
              {
                _t.TextValue = tenNumber.toString()+_t.ChartUnit
              }
              else if(_t.ShowJinZhi==2)
              {
                _t.TextValue = tenNumber.toString(16).toUpperCase()+_t.ChartUnit
              }
              else if(_t.ShowJinZhi==3)
              {
                _t.TextValue = tenNumber.toString(8)+_t.ChartUnit
              }
              else if(_t.ShowJinZhi==4)
              {
                _t.TextValue = tenNumber.toString(2)+_t.ChartUnit
              }
            }
            else
            {
              _t.TextValue = data.result+_t.ChartUnit
            }
          }
        })
        _t.$EventBus.$on(animateEvent, _t._animateEventHandler = (data) => {
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
    _t.$EventBus.$on('cell-editMode', _t._cellEditModeHandler = (data) => {
      _t.editMode = data.edit
      _t.IsToolBox = data.toolbox
    })
    this.initComponents(this.detail);
  },
  beforeDestroy() {
    // 清理 EventBus 监听
    if (this._activeEventHandler) {
      this.$EventBus.$off(this.detail.identifier + 'activeEvent', this._activeEventHandler)
    }
    if (this._animateEventHandler) {
      this.$EventBus.$off(this.detail.identifier + 'animateEvent', this._animateEventHandler)
    }
    if (this._cellEditModeHandler) {
      this.$EventBus.$off('cell-editMode', this._cellEditModeHandler)
    }
  }
}
</script>
<style lang="less" scoped>
.svg-el {
  /*transform: rotate(45deg);*/
  transform-origin: center center;
}
// 键盘样式
   .simple-keyboard-user {
     bottom: 0;
     left: 0%;
     width: 600px;
     color: #000;
     z-index: 999999999;
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
