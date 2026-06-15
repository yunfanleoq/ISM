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
                                opacity:detail.style.opacity,
                                borderWidth: detail.style.borderWidth + 'px',
                                borderStyle: detail.style.borderStyle,
                                borderColor: detail.style.borderColor,
                                transform: detail.style.transform==-1099? 'rotateY(180deg)':detail.style.transform==-1098?'rotateX(180deg)':'',
                            }">
        <a-dropdown :placement="placement">
            <a  :style="{
                     fontSize: detail.style.fontSize + 'px',
                      fontFamily: detail.style.fontFamily,
                      'font-weight':detail.style.fontWeight,
                      color: detail.style.foreColor,
                      'background-color':detail.style.backColor,
                      'text-align':detail.style.textAlign,
                      display:'block',
                      'font-style':detail.style.italic?'oblique':'normal',
                      lineHeight: lineHeight + 'px',
                  }" @click="e => e.preventDefault()">
             {{detail.style.text}}
            </a>
                  <div  slot="overlay" :style="styleVar">
                    <a-menu  :style="{
                      color: detail.style.foreColor,
                      'background-color':detail.style.backColor,
                  }"  >
                          <a-menu-item :style="{
                      width: detail.style.position.w+'px',
                     fontSize: detail.style.fontSize + 'px',
                      fontFamily: detail.style.fontFamily,
                      'text-align':detail.style.textAlign,
                      'font-weight':detail.style.fontWeight,
                      'font-style':detail.style.italic?'oblique':'normal',
                      lineHeight: lineHeight + 'px',
                      height: lineHeight + 'px',
                  }" v-for="(item,index) in MenuList" :key="index" @click="JumpPage(item)">
                            {{item.MenuName}}
                          </a-menu-item>

                      </a-menu>
                  </div>
          </a-dropdown>

  </div>
  </div>

</template>

<script>
import {setData} from "@/services/device";
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
    name: 'view-drop-down',
    i18n: require('../../../../i18n/language'),
    data() {
      return {
        detail:null,
        IsToolBox:false,
        editMode:true,
        Text:"",
        foreColor:"#000000",
        backColor:"#ffffff",
        hoverForeColor:"#000000a6",
        hoverBackColor:"#DAEAF6",
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
        MenuList:[],
        placement:"bottomCenter",
        ClickType:0,
        base:{
          text: "configComponent.DropDown.title",
          "icon": "icon-icon__xuanzexialacaidan",
          "isFontIcon": true,
          "info": {
            "type": "image",
            "action": [],
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
            "active": [
              {
                id:"ControlStatus",
                name:"configComponent.status.DropMenu",
                result:0,
                isStatus:false,
                isMenu:true,
                isSwitch:false,
                isImageStatus:false,
                isTextStatus:false,
                isExpression:false,
                condition:{
                  PageList:[

                  ]
                },
              }
            ],
            "style": {
              "position": {
                "x": 0,
                "y": 0,
                "w": 130,
                "h": 40
              },
              "visible":1,
              "backColor": "#ffffff",
              "foreColor": "#000000",
              "zIndex": -1,
              "transform": 0,
              fontSize: 30,
              fontFamily: "Arial",
              italic:0,
              text: "下拉菜单",
              textAlign: "center",
              fontWeight:400,
              "diy":[
                {
                  "name":"configComponent.Menu.ClickType",
                  "type":6,
                  "enumList":[
                    {option:'configComponent.Menu.ClickTypeEmit',value:1},{option:'configComponent.Menu.ClickTypeJump',value:0}
                  ],
                  "value":0,
                  "key":"ClickType",
                },
                {
                  "name":"configComponent.ComBox.hoverForce",
                  "type": 2,
                  "value": "#000000",
                  "key":"hoverForce",
                },
                {
                  "name":"configComponent.ComBox.hoverBack",
                  "type": 2,
                  "value": "#6ACBFF",
                  "key":"hoverBack",
                },
                {
                  "name":"configComponent.ComBox.placement",
                  "type":6,
                  "enumList":[
                    {option:'configComponent.ComBox.placementDown',value:"bottomCenter"},{option:'configComponent.ComBox.placementUp',value:"topCenter"}
                  ],
                  "value":"bottomCenter",
                  "key":"placement",
                },
              ]
            }
          }
        }
      }
    },
    inject: ['getNode'],
    computed: {
      animatedStyle(){
        return {
          "--tableBorderColor": this.tableBorderColor ,
          "--blinkSpeed":this.blinkSpeed+'s',
          "--stopColor":this.stopColor,
          "--startColor":this.startColor,
          "--animateSpeed":this.animateSpeed+'s',
          "--animateSpinSpeed":this.animateSpinSpeed+'s',
          '--fontWeight':this.detail.style.fontWeight,
          '--fontSize':this.detail.style.fontSize+ 'px',
          '--fontFamily':this.detail.style.fontFamily,
          '--fontStyle':this.detail.style.italic?'oblique':'normal',
          "--height": this.detail.style.position.h+'px',
          "--width": this.detail.style.position.w+'px',
          "--foreColor": this.foreColor ,
          '--backColor':this.backColor,
        }
      },
      styleVar() {
        return {
          "--height": this.detail.style.position.h+'px',
          "--width": this.detail.style.position.w+'px',
          "--foreColor": this.foreColor ,
          '--backColor':this.backColor,
          "--hoverForeColor": this.hoverForeColor ,
          '--hoverBackColor':this.hoverBackColor,
          '--fontWeight':this.detail.style.fontWeight,
          '--fontSize':this.detail.style.fontSize+ 'px',
          '--fontFamily':this.detail.style.fontFamily,
          '--fontStyle':this.detail.style.italic?'oblique':'normal',
        };
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
      JumpPage(item){
        if(this.ClickType==1)
        {
          this.$EventBus.$emit("MenuConfigPage", item);
        }
        else
        {
          this.$EventBus.$emit("ChargePage", item);
        }
      },
      initComponents(option){
        if(this.IsToolBox)
        {
          return
        }
        this.backColor=this.detail.style.backColor
        this.foreColor=this.detail.style.foreColor
        this.MenuList = option.active[0].condition.PageList
        let i=0
        for( i=0;i<option.style.diy.length;i++)
        {
          if(option.style.diy[i].key=="hoverBack")
          {
            this.hoverBackColor=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="hoverForce")
          {
            this.hoverForeColor=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="placement")
          {
            this.placement=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="ClickType")
          {
            this.ClickType=option.style.diy[i].value
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
        if((_t.editMode)&&(!this.IsToolBox)){
          return
        }
        if(data.ID == "ControlStatus")
        {
          let temp  =  parseFloat(data.result)
          if(!isNaN(temp))
          {
            _t.StatusValue=temp
          }
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
<style lang="less" scoped>
.svg-el {
  /*transform: rotate(45deg);*/
  transform-origin: center center;
}
::v-deep li{
  width: var(--width);
  padding: 1px;
}
::v-deep li a{
  color: var(--foreColor);
}
::v-deep li:hover{
  background-color:  var(--hoverBackColor);
  color: var(--hoverForeColor);
}
::v-deep li:hover a{
  //background-color:  var(--hoverBackColor);
  color: var(--hoverForeColor);
}
::v-deep .ant-menu .ant-menu-item-selected {
  background-color:  var(--hoverBackColor);
  color: var(--hoverForeColor);
}
::v-deep .ant-menu{
  background-color:  var(--hoverBackColor);
}
</style>
