<template>
      <div class="boxer xunan-login" :style="styleVar">
        <div class="header">
          <div class="inner-header flex">
            <div class="box1">
              <div class="middle1">
                <div class="middle_left">
                  <div class="brand-panel">
                    <a-skeleton :loading="skeletonLoading" active>
                      <img :src="systemLogo" alt="循安" class="brand-logo" />
                    </a-skeleton>
                    <h1 class="brand-title">循安科技电力监控平台</h1>
                    <p class="brand-subtitle">POWER MONITORING &amp; INTELLIGENT MANAGEMENT</p>
                    <ul class="brand-features">
                      <li>实时数据采集与监控</li>
                      <li>智能告警与故障预警</li>
                      <li>能耗分析与报表统计</li>
                    </ul>
                  </div>
                </div>
                <div class="middle_right">
                  <div class="mr_box">
                    <div class="login-head">
                      <a-skeleton :loading="skeletonLoading" active>
                        <img :src="systemLogo" alt="logo" class="login-card-logo" />
                      </a-skeleton>
                      <h2>登录</h2>
                      <p class="login-welcome">欢迎登录，请输入您的账号信息</p>
                    </div>
                    <a-form @submit="onSubmit" :form="form">
                      <a-alert type="error" :closable="true" v-show="error" :message="error" showIcon style="margin-bottom: 24px;" />
                      <a-form-item>
                        <a-input
                            class="form-control"
                            autocomplete="autocomplete"
                            size="large"
                            :placeholder="$t('loginPage.logonUserTips')"
                            v-decorator="['name', {rules: [{ required: true, message:  $t('loginPage.logonUserTips'), whitespace: true}]}]"
                        >
                          <a-icon slot="prefix" type="user" style="position: absolute;top: -5px;right:-15px;color: #13c2c2;font-size: 20px;"/>
                        </a-input>
                      </a-form-item>
                      <a-form-item style="margin-bottom: 5px;">
                        <a-input
                            class="form-control"
                            size="large"
                            autocomplete="autocomplete"
                            type="password"
                            :placeholder="$t('loginPage.logonPasswordTips')"
                            v-decorator="['password', {rules: [{ required: true, message: $t('loginPage.logonPasswordTips'), whitespace: true}]}]"
                        >
                          <a-icon slot="prefix" type="lock" size="large" style="position: absolute;top: -5px;right:-15px;color: #13c2c2;font-size: 20px;"/>
                        </a-input>
                      </a-form-item>
                      <a-form-item style="margin-bottom: 5px;">
                        <a-checkbox style="float: left"
                                    v-decorator="['autologin', { valuePropName: 'checked' }]"
                        >{{$t('loginPage.AutoLogin')}}</a-checkbox>
                      </a-form-item>
                      <a-form-item>
                        <a-button :loading="logging" class="login" style="width: 100%;margin-top: 2px" size="large" htmlType="submit" type="primary">{{$t('loginPage.logonBtn')}}</a-button>
                      </a-form-item>
                      <div>
                        <a-breadcrumb>
                          <a-breadcrumb-item   v-for=" lang in langList" :key="lang.key" ><a @click="setLang(lang.key)">{{lang.name}}</a></a-breadcrumb-item>
                        </a-breadcrumb>
                      </div>
                    </a-form>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <!--Waves Container-->
          <div>
            <svg class="waves" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 24 150 28" preserveAspectRatio="none" shape-rendering="auto">
              <defs>
                <path id="gentle-wave" d="M-160 44c30 0 58-18 88-18s 58 18 88 18 58-18 88-18 58 18 88 18 v44h-352z" />
              </defs>
              <g class="parallax">
                <use xlink:href="#gentle-wave" x="48" y="0" fill="rgba(255,255,255,0.7" />
                <use xlink:href="#gentle-wave" x="48" y="3" fill="rgba(255,255,255,0.5)" />
                <use xlink:href="#gentle-wave" x="48" y="5" fill="rgba(255,255,255,0.3)" />
                <use xlink:href="#gentle-wave" x="48" y="7" fill="#fff" />
              </g>
            </svg>
          </div>
          <!--Waves end-->
        </div>
        <!-- 整页底部居中，对齐登录参考图；勿放在登录卡片内 -->
        <div class="xunan-footer">
          © 北京循安科技有限公司 - 循安科技电力监控平台
        </div>
      </div>
</template>

<script>
import CommonLayout from '@/layouts/CommonLayout'
import {login, getRoutesConfig} from '@/services/user'
import {AUTH_TYPE, checkAuthorization, setAuthorization} from '@/utils/request'
import {loadRoutes} from '@/utils/routerUtil'
import {mapState, mapMutations, mapGetters} from 'vuex'
import md5 from 'js-md5';
import {GetSystemParams} from "@/services/system";
import {ProjectList} from "@/services/project";
import {applyHomeProjectAuth} from '@/config/homeDashboard'
export default {
  name: 'Login',
  components: {},
  i18n: require('../../i18n/language'),
  data () {
    return {
      logging: false,
      error: '',
      form: this.$form.createForm(this)
    }
  },
  computed: {
    ...mapState('setting', ['langList','isMobile','lang','skeletonLoading']),
    ...mapGetters('setting', ['homeDashboardPath']),
    systemName () {
      return this.$store.state.setting.systemName
    },

    systemLogo () {
      const logo = this.$store.state.setting.SystemLogo
      if (!logo || String(logo).indexOf('data:image') === 0) {
        return '/static/branding/logo-xunan-hexagon.png'
      }
      return logo
    },
    systemLoginBg () {
      return this.$store.state.setting.systemLoginBg
    },
    styleVar() {
      return {
        '--bgImage':`url('${this.$store.state.setting.systemBg}')`,
      };
    },
  },
  mounted() {
  },
  created(){
    this.GetSystemCas()
    let autologin = localStorage.getItem("autologin")
    let User =  localStorage.getItem("User")
    let Password =  localStorage.getItem("Password")
    if((autologin)&&(User)&&(Password))
    {
      this.autoLogin(User,Password)
    }
  },
  methods: {
    ...mapMutations('setting', ['setLang']),
    ...mapMutations('account', ['setUser', 'setPermissions', 'setRoles','setRoutesConfig']),
    autoLogin(user,password){
      let _t = this
      const Username = user
      const passwordMd5 = password
      login(Username, passwordMd5).then(this.afterLogin).catch(function(){
        _t.logging = false
        _t.$message.error(_t.$t('loginPage.serverError'), 3)
      })
    },
    onSubmit (e) {
      e.preventDefault()
      this.form.validateFields((err) => {
        if (!err) {
          this.logging = true
          let _t = this
          const Username = this.form.getFieldValue('name')
          const password = md5(this.form.getFieldValue('password'))
          const autologin = this.form.getFieldValue('autologin')
          if(autologin)
          {
            localStorage.setItem("autologin",autologin)
            localStorage.setItem("User",Username)
            localStorage.setItem("Password",password)
          }
          login(Username, password).then(this.afterLogin).catch(function(){
            _t.logging = false
            _t.$message.error(_t.$t('loginPage.serverError'), 3)
          })
        }
      })
    },
    afterLogin(res) {
      this.logging = false
      const loginRes = res.data
      if (loginRes.code == 1000) {
        const {user, roles,} = loginRes.data
        this.setUser(user)
        this.setRoles(roles)
        this.setRoutesConfig(user.Menu)
        localStorage.setItem("LoginFrom",'/login')
        setAuthorization({token: loginRes.data.token, expireAt: loginRes.data.expireAt})
        this.$message.success(this.$t('loginPage.logonSuccess'), 3)
        const _t = this
        const proceedAfterLogin = () => {
          loadRoutes(user.Menu)
          _t.enterAfterAuth(roles, user)
        }
        this.$store.dispatch('setting/fetchSystemHomeDashboard')
          .then(proceedAfterLogin)
          .catch(function () {
            proceedAfterLogin()
          })
      } else {
        this.$message.error(this.$t('loginPage.logonFailed'), 3)
      }
    },
    // 根据当前用户可访问项目数量决定登录落地页：
    //  - User/Operator：已绑定单一项目，写入 ProjectUuid 后直接进大屏
    //  - Admin/其它：仅 1 个项目则自动选中并进大屏；多个项目则进项目选择页
    enterAfterAuth(roles, user) {
      const _t = this
      const roleId = (roles && roles[0]) ? roles[0].id : ''
      const goHome = () => {
        applyHomeProjectAuth(_t.$store)
        _t.$router.push(_t.homeDashboardPath)
      }
      if (roleId == "User" || roleId == "Operator") {
        if (user && user.ProjectUUID) {
          setAuthorization({token: user.ProjectUUID}, AUTH_TYPE.AUTH1)
        }
        goHome()
        return
      }
      ProjectList().then(function (res) {
        const list = (res && res.data && res.data.code == 0 && res.data.list) ? res.data.list : []
        if (list.length == 1 && list[0].ProjectInfo) {
          setAuthorization({token: list[0].ProjectInfo.uuid}, AUTH_TYPE.AUTH1)
          goHome()
        } else {
          _t.$router.push('/project')
        }
      }).catch(function () {
        applyHomeProjectAuth(_t.$store)
        goHome()
      })
    },
    GetSystemCas(){
      let _t = this
      GetSystemParams().then(function (res) {
        if(checkAuthorization()&&res.data.list.UserData!=null)
        {
          const {user, roles,Menu} = res.data.list.UserData
          _t.setUser(user)
          _t.setRoles(roles)
          _t.setRoutesConfig(user.Menu)
          loadRoutes(user.Menu)
          if(roles[0].id=="Admin")
          {
            _t.$store.dispatch('setting/fetchSystemHomeDashboard').then(function () {
              loadRoutes(user.Menu)
              _t.enterAfterAuth(roles, user)
            })
          }
        }
      }).catch(function(e){
        console.log(e)
      })
    }
  }
}
</script>

<style lang="less" scoped>
.header {
  position:relative;
  text-align:center;
  height: 90%;
  background-size: 100% 100%;
  color:white;
}

.inner-header {
  height:100%;
  width:100%;
  margin: 0;
  padding: 0;
}

.flex {
  display: flex;
  justify-content: center;
  align-items: center;
  text-align: center;
}

.waves {
  position:relative;
  width: 100%;
  height:15vh;
  margin-bottom:-7px;
  min-height:100px;
  max-height:150px;
}

.parallax > use {
  animation: move-forever 25s cubic-bezier(.55,.5,.45,.5) infinite;
}
.parallax > use:nth-child(1) {
  animation-delay: -2s;
  animation-duration: 7s;
}
.parallax > use:nth-child(2) {
  animation-delay: -3s;
  animation-duration: 10s;
}
.parallax > use:nth-child(3) {
  animation-delay: -4s;
  animation-duration: 13s;
}
.parallax > use:nth-child(4) {
  animation-delay: -5s;
  animation-duration: 20s;
}
@keyframes move-forever {
  0% {
    transform: translate3d(-90px,0,0);
  }
  100% {
    transform: translate3d(85px,0,0);
  }
}
@media (max-width: 768px) {
  .waves {
    height:40px;
    min-height:40px;
  }
  h1 {
    font-size:24px;
  }
}
  .boxer {
    background-image: linear-gradient(135deg, #0a1a3a 0%, #0d2b55 45%, #0a3d5c 100%);
    height: 100%;
    position: relative;
    overflow: hidden;
    background-repeat: no-repeat;
    background-size: cover;
  }
  .box1 {
    width: 1156px;
    margin:6% auto 0;
    overflow: hidden;
    box-shadow: 0px 0px 24px 8px rgba(0, 40, 80, .35);
    border-radius: 12px;
    background: rgba(8, 28, 58, 0.55);
    border: 1px solid rgba(19, 194, 194, 0.25);
  }
  .middle1 {
    overflow: hidden;
  }
  .middle_left {
    float: left;
    width: 560px;
    height: 600px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: radial-gradient(ellipse at 30% 20%, rgba(19,194,194,0.18), transparent 55%),
      linear-gradient(160deg, rgba(10,40,80,0.9), rgba(8,24,50,0.95));
  }
  .brand-panel {
    padding: 40px 48px;
    text-align: left;
    color: #e8f7f7;
  }
  .brand-logo {
    width: 88px;
    height: 88px;
    object-fit: contain;
    margin-bottom: 20px;
  }
  .brand-title {
    margin: 0 0 10px;
    font-size: 32px;
    font-weight: 700;
    letter-spacing: 2px;
    color: #ffffff;
    line-height: 1.3;
  }
  .brand-subtitle {
    margin: 0 0 28px;
    font-size: 12px;
    letter-spacing: 1.5px;
    color: rgba(160, 220, 230, 0.85);
  }
  .brand-features {
    list-style: none;
    padding: 0;
    margin: 0;
  }
  .brand-features li {
    margin: 12px 0;
    padding: 12px 16px;
    border-radius: 8px;
    background: rgba(12, 40, 72, 0.65);
    border: 1px solid rgba(19, 194, 194, 0.35);
    color: #d7f3f3;
    font-size: 15px;
  }
  .middle_right {
    width: 500px;
    height: 600px;
    float: right;
    background: rgba(255, 255, 255, 0.96);
  }
  .mr_box {
    width: 318px;
    margin: 0 auto;
    margin-top: 40px;
  }
  .login-head {
    text-align: center;
    margin-bottom: 8px;
  }
  .login-card-logo {
    width: 72px;
    height: 72px;
    object-fit: contain;
  }
  .login-welcome {
    margin: -8px 0 24px;
    color: #6b7c8a;
    font-size: 13px;
  }
  .mr_box form{
    text-align: center;
  }
  .mr_box h2 {
    letter-spacing: 2px;
    margin: 12px 0 8px;
    display: block;
    text-align: center;
    color: #13c2c2;
    font-size: 28px;
  }
  .form-control {
    height: 30px;
    width: 320px;
  }
  .login {
    height: 50px;
    line-height: 50px;
    width: 82%;
    border: none;
    border-radius: 60px;
    background: linear-gradient(90deg, #13c2c2, #08979c);
    color: #fff;
    font-size: 18px;
    letter-spacing: 2px;
    cursor: pointer;
    transition: .3s all linear;
    margin-top: 40px;
  }
  .login:hover {
    background: #035757;
    transition: .3s all linear;
  }
  .xunan-footer {
    position: absolute;
    left: 0;
    right: 0;
    bottom: 28px;
    z-index: 5;
    text-align: center;
    color: rgba(210, 230, 240, 0.88);
    font-size: 13px;
    letter-spacing: 0.5px;
    line-height: 1.4;
    pointer-events: none;
    white-space: nowrap;
  }
  @media screen and (min-width: 960px) and (max-width:1200px){
    .box1 {
      margin-top: 144px;
      width: 960px;;
    }
    .mr_box{
      margin-top: 40px;
    }
    .middle_left {
      width: 560px;
      height: 470px;
    }
    .middle1 {
      overflow: hidden;
      height: 470px;
    }
    .middle_right {
      width: 400px;
      overflow: hidden;
      height: 470px;
    }
    .brand-title {
      font-size: 26px;
    }
    .login {
      margin-top: 20px;
    }
  }
  @media screen and (max-width: 960px) {
    .box1 {
      margin-top: 80px;
      width: 418px;
    }
    .mr_box{
      margin-top: 40px;
    }
    .middle_left {
      display: none;
      float: none;
    }
    .middle_right {
      float: none;
      margin: 0 auto;
      overflow: hidden;
      height: 470px;
      width: 418px;
    }
    .xunan-footer {
      bottom: 16px;
      white-space: normal;
      padding: 0 16px;
    }
    .login {
      margin-top: 18px;
    }
  }
</style>
