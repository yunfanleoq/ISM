<template>

  <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"   x="0px" y="0px"   xml:space="preserve" :style="{'overflow': 'visible','width':detail.style.position.w,'height':detail.style.position.h,}">
    <g class="svg-el" :style="{'opacity':fillOpacity,'stroke-opacity':strokeOpacity,'stroke':strokeColor,'stroke-width':strokeWidth,'stroke-linecap':'round','stroke-linejoin':'round','fill':fill}">
      <image style="cursor: pointer" preserveAspectRatio="none meet" @mouseover="mouseenter" @mouseout="mouseleave" @click="loginRequest" :width="detail.style.position.w" :height="detail.style.position.h" :href="imageURL"></image>
      <foreignObject style="overflow:visible;" pointer-events="all" :width="detail.style.position.w" :height="detail.style.position.h">
        <div
            @mouseover="mouseenter"
            @mouseout="mouseleave"
            @click="loginRequest"
            :style="{
              cursor: 'pointer',
        fontSize: detail.style.fontSize + 'px',
        fontFamily: detail.style.fontFamily,
        'font-weight':detail.style.fontWeight,
        color: detail.style.foreColor,
        textAlign: textAlign,
        'font-style':detail.style.italic?'oblique':'normal',
        lineHeight: lineHeight + 'px',
    }">
          {{detail.style.text}}
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
import svgView from '../View';
import {AUTH_TYPE, setAuthorization} from '@/utils/request'
import md5 from 'js-md5';
import {login, getRoutesConfig} from '@/services/user'
import {mapState, mapMutations} from 'vuex'
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
export default {
  mixins: [ISMChildAutoMixin],
  name: 'view-svg-login-button',
  inject: ['getNode'],
  i18n: require('@/i18n/language'),
  data() {
    return {
      loginUuid:"",
      PStore:{},
      PMapMutations:{},
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
      mouseleaveUrl:"",
      mouseenterUrl:"",
      textX:0,
      textY:0,
      base:{
        "text": "displayConfig.ToolBox.login.loginBtn",
        "icon": "icon-dengluzhanghao",
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
              "w": 157,
              "h": 69
            },
            "visible":1,
            "zIndex": -1,
            "transform": 0,
            text: "登录",
            "foreColor": "#2BF3D1",
            fontSize: 30,
            fontFamily: "Arial",
            "diy":[
              {
                "name":"configComponent.Button.textY",
                "type":7,
                "value":45,
                "key":"textY",
              },
              {
                "name":"configComponent.Button.textX",
                "type":7,
                "value":50,
                "key":"textX",
              },
              {
                "name":"component.public.fillOpacity",
                "type":7,
                "value":1,
                "key":"fillOpacity",
              },
              {
                "name":"configComponent.Button.BackImg",
                "type":5,
                "value":"/static/ISM/systemImage/button/72-1.png",
                "key":"imageURL",
              },
              {
                "name":"configComponent.Button.hoverForce",
                "type":5,
                "value":"/static/ISM/systemImage/button/72-0.png",
                "key":"imageEnter",
              }
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
    },
    IsLicense () {
      return this.PStore.state.setting.IsLicense
    },
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
          this.mouseleaveUrl = this.imageURL
        }
        else if(option.style.diy[i].key=="imageEnter")
        {
          this.mouseenterUrl=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="textY")
        {
          this.textY=option.style.diy[i].value
        }
        else if(option.style.diy[i].key=="textX")
        {
          this.textX=option.style.diy[i].value
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
    },
    mouseenter(e){
      this.imageURL =this.mouseenterUrl
    },
    mouseleave(e){
      this.imageURL =this.mouseleaveUrl
    },
    afterLogin(res) {
      this.logging = false
      const loginRes = res.data
      if (loginRes.code == 1000) {
        const {user, roles} = loginRes.data

        setAuthorization({token: user.ProjectUUID},AUTH_TYPE.AUTH1)
        this.$EventBus.$emit('AppLoginSuccess', {user, roles})
        localStorage.setItem("LoginFrom",'/AppRun/'+this.loginUuid)
        setAuthorization({token: loginRes.data.token, expireAt: loginRes.data.expireAt})
        this.PRouter.push('/AppRun/'+this.loginUuid).catch(err => (
            console.log(err)
        ))
        this.$message.success(this.$t('loginPage.logonSuccess'), 3)
      } else if (loginRes.code == -17) {
        this.$message.error(this.$t('loginPage.logonFailedNoAuth'), 3)
      } else {
        this.$message.error(this.$t('loginPage.logonFailed'), 3)
      }
    },
    loginRequest(){
      let _t = this
      if(_t.editMode){
        return
      }
      if(!_t.IsLicense){
        document.title = _t.$t('displayConfig.Properties.NoAuth')
        _t.$message.warn(_t.$t('displayConfig.Properties.NoAuth'))
        return
      }
      let appLoginUsername=document.getElementById("appLoginUsername");
      let appLoginPassword=document.getElementById("appLoginPassword");
      if ((appLoginUsername==null)||(appLoginPassword==null))
      {
        this.$message.error("用户名、密码不能为空！！")
        return
      }
      appLoginPassword = appLoginPassword.value;
      appLoginUsername = appLoginUsername.value;
      if(appLoginPassword==""||appLoginUsername=="")
      {
        this.$message.error("用户名、密码不能为空！！")
        return
      }
      login(appLoginUsername, md5(appLoginPassword),this.loginUuid).then(this.afterLogin).catch(function(e){
        _t.logging = false
        console.log(e)
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
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
    _t.$EventBus.$on('cell-vuex', (data) => {
      _t.PRouter=data.PRouter
     _t.loginUuid = data.loginUuid
      _t.PStore=data.PStore
      _t.PMapMutations = data.PMapMutations('account', ['setUser', 'setPermissions', 'setRoles'])
    })
  }
}
</script>
<style >
.svg-el {
  /*transform: rotate(45deg);*/
  transform-origin: center center;
}

</style>
