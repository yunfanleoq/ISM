<template>
  <div v-show="detail.style.visible==1 ||isStart? true:false">
    <div :class="{
          'animated':true,[`${detail.style.animate}`]: true
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
  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"  x="0px" y="0px"  xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el" :class="{'animated':true,[`${detail.style.animate}`]: true}" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <template v-for="(item,index) in StatusList" >
          <template v-if="(item.StatusOpt == '==')||(typeof item.StatusOpt == 'undefined')">
            <image preserveAspectRatio="none meet" :key="index" v-if="StatusValue==item.value" :width="detail.style.position.w" :height="detail.style.position.h" :href="item.Image"></image>
          </template>
          <template v-if="item.StatusOpt == '>'" >
            <image preserveAspectRatio="none meet" :key="index" v-if="StatusValue>item.value" :width="detail.style.position.w" :height="detail.style.position.h" :href="item.Image"></image>
          </template>
          <template v-if="item.StatusOpt == '>='">
                <image preserveAspectRatio="none meet" :key="index" v-if="StatusValue>=item.value" :width="detail.style.position.w" :height="detail.style.position.h" :href="item.Image"></image>
          </template>
          <template v-if="item.StatusOpt == '<'">
               <image preserveAspectRatio="none meet" :key="index" v-if="StatusValue<item.value" :width="detail.style.position.w" :height="detail.style.position.h" :href="item.Image"></image>
          </template>
          <template v-if="item.StatusOpt == '<='">
                <image preserveAspectRatio="none meet" :key="index" v-if="StatusValue<=item.value" :width="detail.style.position.w" :height="detail.style.position.h" :href="item.Image"></image>
          </template>
          <template v-if="item.StatusOpt == '<>'" >
              <image preserveAspectRatio="none meet" :key="index" v-if="(StatusValue>=item.value)&&(StatusValue<=item.value2)" :width="detail.style.position.w" :height="detail.style.position.h" :href="item.Image"></image>
          </template>
          <template v-if="item.StatusOpt == '<!>'" >
              <image preserveAspectRatio="none meet" :key="index" v-if="(StatusValue<=item.value)||(StatusValue>=item.value2)" :width="detail.style.position.w" :height="detail.style.position.h" :href="item.Image"></image>
          </template>
          <template v-if="item.StatusOpt == '!='" >
                <image preserveAspectRatio="none meet" :key="index" v-if="(StatusValue!=item.value)" :width="detail.style.position.w" :height="detail.style.position.h" :href="item.Image"></image>
          </template>
      </template>

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
    </div>
  </div>
</template>

<script>
import {mapState} from "vuex";
import store from "@/store";
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

export default {
  mixins: [ISMChildAutoMixin],
    name: 'view-svg-image-status',
    i18n: require('../../../../i18n/language'),
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
        StatusList:[],
        StatusValue:0,
        base:{
          text: "configComponent.status.ImageTitle",
          "icon": "icon-icon_status",
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
                name:"configComponent.status.ControlStatus",
                result:0,
                isStatus:true,
                isSwitch:false,
                isImageStatus:true,
                isTextStatus:false,
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
                      "value2":1,
                      "Image":"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAAAXNSR0IArs4c6QAAAWZJREFUWEftlzFOAzEQRV9OAQUUCCooETXkBmmg5xS0iDvQ0wINNwBq6qQCJRJCAnEI0Ece5KxMNh57WQlhKYW16/nPf8beyYCex6BnfVIA68BmR2BPwHMcuwlwCxx0JG5h74ChTWIA7foxPNBLXQzb3BYgN+ZSoIdyYI6wMoU5LAe+Nhk78A/Q5sANsBdS8gCMHOlxp+DjB7Hce8QFcAUcAqfAfQDZD/Nr4CjDCRfAK7AC7ACTILYNjIE3YPXPA3hSsAFME864UqA4yxahauUE2AUugbMobYrjBtDitmMocbmloVtOx1p1oiK12ikCWFRrsbgEdTosdTFEJwApcYNtQpwHZ6p9CxaJpyDeawIsI56C0H1S7ECOeBNC8yIAj7hBFBdhiXjxPVAqXgRQQ9wNUEvcBVBTvBUg1ZZbG12rTbd4a8BLsyuOCTN6jOxXL4BjW/Xbf81mgH7fI7epzN5u24LeAT4B+PijITzWRRcAAAAASUVORK5CYII=",
                      "value":0
                    }
                  ]
                },
              }
            ],
            "style": {
              "position": {
                "x": 0,
                "y": 0,
                "w": 50,
                "h": 50
              },
              "visible":1,
              "zIndex": -1,
              "transform": 0,
              "diy":[

              ]
            }
          }
        }
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
    methods: {
      initComponents(option){
        if(this.IsToolBox)
        {
          return
        }
        this.fillOpacity = option.style.opacity
        this.StatusList = option.active[0].condition.StatusList
        if(option.active[0].condition.StatusList.length)
        {
          this.StatusValue = parseFloat(option.active[0].condition.StatusList[0].value)
        }
        let i=0
        for( i=0;i<option.style.diy.length;i++)
        {
          if(option.style.diy[i].key=="On")
          {
            this.SwitchOn=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="Off")
          {
            this.SwitchOff=option.style.diy[i].value
          }
          else if(option.style.diy[i].key=="SwitchStyle")
          {
            this.switchStyleIndex=option.style.diy[i].value
          }
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
            let TempStatusValue = parseFloat(data.result)
            _t.StatusValue = TempStatusValue
            // for(let i=0;i<_t.detail.active[0].condition.StatusList.length;i++)
            // {
            //   if(TempStatusValue==_t.detail.active[0].condition.StatusList[i].value)
            //   {
            //     _t.StatusValue = TempStatusValue
            //     break
            //   }
            // }
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
    _t.$EventBus.$on('cell-editMode', (data) => {
      _t.editMode = data.edit
      _t.IsToolBox = data.toolbox
	  // _t.initComponents(_t.detail)
    })
  }
}
</script>
<style >
.svg-el {
  /*transform: rotate(45deg);*/
  transform-origin: center center;
}
::-webkit-scrollbar {
  /*滚动条整体样式*/
  width : 5px;  /*高宽分别对应横竖滚动条的尺寸*/
  height: 9px;
}
::-webkit-scrollbar-thumb {
  /*滚动条里面小方块*/
  border-radius   : 10px;
  background-color: skyblue;
  background-image: -webkit-linear-gradient(
      45deg,
      rgba(255, 255, 255, 0.2) 25%,
      transparent 25%,
      transparent 50%,
      rgba(255, 255, 255, 0.2) 50%,
      rgba(255, 255, 255, 0.2) 75%,
      transparent 75%,
      transparent
  );
}
::-webkit-scrollbar-track {
  /*滚动条里面轨道*/
  box-shadow   : inset 0 0 5px rgba(0, 0, 0, 0.2);
  background   : #ededed;
  border-radius: 10px;
}
</style>
