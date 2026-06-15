<template>
<div>
  <input type="text" id="appLoginUsername"  v-model="appLoginUsername" @focus="onInputFocus('SetAutoPassword')" autocomplete="off" placeholder="用户名" :style="{'padding-left':leftDis+'px','border':0,'font-size': detail.style.fontSize+'px','border-radius':detail.style.BorderEdges+'px','outline':'none','color':detail.style.foreColor,'background-color':detail.style.backColor,'height': detail.style.position.h+'px','width':detail.style.position.w+'px'}"/>

    <SimpleKeyboardUser v-show="showKeyboarduser&&virtuallyKey" ref="SimpleKeyboardloginUser" @onChange="onChangeKeyboard" @onKeyPress="onKeyPress" keyboardClass="simple-keyboard-passowrd"/>
</div>
</template>

<script>
import svgView from '../View';
import SimpleKeyboardUser from '@/components/SimpleKeyboard/SimpleKeyboard.vue'
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
export default {
  mixins: [ISMChildAutoMixin],
    name: 'view-svg-login-user-name',
    inject: ['getNode'],
    data() {
      return {
        detail:null,
        IsToolBox:false,
        editMode:true,
        strokeColor:"#000000",
        fill:"#A1BFE2",
        strokeWidth:0.3,
        leftDis:10,
        appLoginUsername:"",
        showKeyboarduser:false,
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
        italic:false,
        imageURL:"",
        virtuallyKey:false,
        base:{
          "text": "displayConfig.ToolBox.login.loginUser",
          "icon": "icon-yonghuming",
          "isFontIcon": true,
          "info": {
            "type": "image",
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
                "w": 250,
                "h": 40
              },
              "visible":1,
              "backColor": "#ffffff",
              "foreColor": "#000000",
              fontWeight:400,
              "zIndex": 3,
              BorderEdges:5,
              "transform": 0,
              fontSize: 20,
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
              ]
            }
          }
        }
      }
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
    watch: {
      detail: {
        handler(newVal, oldVal) {
          this.initComponents(newVal);
        },
        deep: true
      }
    },
    components: {
      SimpleKeyboardUser
    },
    methods: {
      onKeyPress(button){
        if (button === '{enter}' || button === '{close}') {
          this.closekeyboard()
        }
      },
// inpuit获取焦点显示虚拟键盘
      onInputFocus(res) {
        this.showKeyboarduser = true
        this.changeIpt = res
        // 父组件调用子组件的方法
        this.$refs.SimpleKeyboardloginUser.onKeyPress('{clear}')
      },
      // 给输入框赋值
      onChangeKeyboard(input) {
        this.appLoginUsername = input
      },
      // 点击关闭隐藏键盘
      closekeyboard() {
        this.showKeyboarduser = false
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
          else if(option.style.diy[i].key=="leftDis")
          {
            this.leftDis=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="virtuallyKey")
          {
            this.virtuallyKey=option.style.diy[i].value
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
    }
}
</script>
<style lang="less" scoped>
.svg-el {
  /*transform: rotate(45deg);*/
  transform-origin: center center;
}
// 键盘样式
.simple-keyboard-passowrd {
  bottom: 0;
  left: 0%;
  width: 600px;
  color: #000;
  z-index: 6000;
}
</style>
