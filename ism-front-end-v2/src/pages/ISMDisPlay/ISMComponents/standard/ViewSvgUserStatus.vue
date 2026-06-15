<template>

  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"  x="0px" y="0px"  xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el" :class="{'animated':true,[`${detail.style.animate}`]: true}" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <foreignObject style="overflow:visible;" pointer-events="all" :width="detail.style.position.w" :height="detail.style.position.h">
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
        <a-dropdown>
            <div class="header-avatar" style="cursor: pointer">
              <a-avatar class="avatar" size="small" shape="circle" v-if="user!=null" :src="user.avatar"/>
              <span v-if="user!=null" class="name" :style="{
        fontSize: detail.style.fontSize + 'px',
        fontFamily: detail.style.fontFamily,
        color: detail.style.foreColor,
        'font-weight':detail.style.fontWeight,
        textAlign: textAlign,
        lineHeight: lineHeight + 'px',
    }">{{user.name}}</span>
              <span v-else class="name" :style="{
        fontSize: detail.style.fontSize + 'px',
        fontFamily: detail.style.fontFamily,
        color: detail.style.foreColor,
        'font-weight':detail.style.fontWeight,
        textAlign: textAlign,
        lineHeight: lineHeight + 'px',
    }">无登录用户</span>
            </div>
          <a-menu v-if="user!=null" :class="['avatar-menu']" slot="overlay" :style="{
        fontSize: detail.style.fontSize + 'px',
        fontFamily: detail.style.fontFamily,
        color: detail.style.foreColor,
        'background-color': detail.style.backColor,
        lineHeight: lineHeight + 'px',
    }">
              <template >
                <a-menu-item >
                  <a target="_blank" @click="gotoSetting" :style="{
        fontSize: detail.style.fontSize + 'px',
        fontFamily: detail.style.fontFamily,
        color: detail.style.foreColor,
    }">
                    <a-icon type="user" />
                    <span style="margin-left: 5px">{{$t('header.Profile')}}</span>
                  </a>
                </a-menu-item>
                <a-menu-divider />
              </template>
              <a-menu-item @click="logout" :style="{
        fontSize: detail.style.fontSize + 'px',
        fontFamily: detail.style.fontFamily,
        color: detail.style.foreColor,
    }">
                  <a-icon style="margin-right: 8px;" type="poweroff" />
                  <span>{{$t('header.Logout')}}</span>
              </a-menu-item>
            </a-menu>
        </a-dropdown>
          </div>
        </div>
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
import {mapGetters} from 'vuex'
import Cookie from 'js-cookie'
import router from '@/router'
import {logout} from "@/services/user";
import account from "@/store/modules/account";
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
export default {
  mixins: [ISMChildAutoMixin],
    name: 'view-svg-user-status',
    inject: ['getNode'],
    components: {},
    i18n: require('../../../../i18n/language'),
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
    data() {
      return {
        user:{},
        PRouter:{},
        detail:null,
        IsToolBox:false,
        editMode:true,
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
        base:{
          text: "configComponent.UserStatus.Text",
          "icon": "icon-techreport-",
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
                "w": 127,
                "h": 34
              },
              "visible":1,
              "backColor": "transparent",
              "foreColor": "#000000",
              "zIndex": -1,
              "transform": 0,
              fontWeight:400,
              textAlign: "center",
              fontSize: 15,
              fontFamily: "Arial",
              "diy":[
                {
                  "name":"component.public.fillOpacity",
                  "type":7,
                  "value":1,
                  "min":0,
                  "max":1,
                  "key":"fillOpacity",
                },
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
    methods: {
      gotoSetting() {
        if(this.editMode)
        {
          return
        }
        this.PRouter.push("/Setting/Account")
      },
      logout() {
        if(this.editMode)
        {
          return
        }
        logout()
        let LoginFrom = localStorage.getItem("LoginFrom")
        if(LoginFrom!=null)
        {
          this.PRouter.push(LoginFrom)
          location.reload()
        }
        else
        {
          this.PRouter.push('/login')
          Cookie.remove('ProjectUuid')
        }

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
          else if(option.style.diy[i].key=="imageURL")
          {
            this.imageURL=option.style.diy[i].value
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
      _t.$EventBus.$on('cell-vuex', (data) => {
        _t.PRouter=data.PRouter
        _t.initComponents(_t.detail)
      })
      const userGet = localStorage.getItem(process.env.VUE_APP_USER_KEY)
      _t.user = JSON.parse(userGet)
  }
}
</script>
<style >
.svg-el {
  /*transform: rotate(45deg);*/
  transform-origin: center center;
}

</style>
