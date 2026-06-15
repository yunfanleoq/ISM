<template>
      <div class="boxer" :style="styleVar">
        <div class="header">
          <div class="inner-header flex">
            <div class="box1">
              <div class="middle1">
                <div class="middle_left">
                  <a-skeleton :loading="skeletonLoading" active>
                    <img :src="systemLoginBg"  style="width: 670px;height: 600px;">
                  </a-skeleton>
                </div>
                <div class="middle_right">
                  <div class="mr_box">
                    <div>
                      <div>
                        <a-skeleton :loading="skeletonLoading" active>
                        <img :src="systemLogo" alt="logo" style="width: 128px;height: 128px;">
                        </a-skeleton>
                      </div>
                      <h2>{{$t("loginPage.logonBtn")}}  </h2>
                    </div>
                    <a-form @submit="onSubmit" :form="form">
                      <a-alert type="error" :closable="true" v-show="error" :message="error" showIcon style="margin-bottom: 24px;" />
                      <a-form-item>
                        <a-input
                            class="form-control"
                            autocomplete="autocomplete"
                            size="large"
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
                            v-decorator="['password', {rules: [{ required: true, message: $t('loginPage.logonPasswordTips'), whitespace: true}]}]"
                        >
                          <a-icon slot="prefix" type="lock" size="large" style="position: absolute;top: -5px;right:-15px;color: #13c2c2;font-size: 20px;"/>
                        </a-input>
                      </a-form-item>
                      <a-form-item style="margin-bottom: 5px;">
                        <a-checkbox style="float: left"
                                    v-decorator="['autologin']"
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
      </div>
</template>

<script>
import CommonLayout from '@/layouts/CommonLayout'
import {login, getRoutesConfig} from '@/services/user'
import {AUTH_TYPE, checkAuthorization, setAuthorization} from '@/utils/request'
import {loadRoutes} from '@/utils/routerUtil'
import {mapState, mapMutations} from 'vuex'
import md5 from 'js-md5';
import {GetSystemParams} from "@/services/system";
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
    systemName () {
      return this.$store.state.setting.systemName
    },

    systemLogo () {
      return this.$store.state.setting.SystemLogo
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
    if(this.isMobile)
    {
      this.$router.push('/loginPhone')
    }
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
        loadRoutes(user.Menu)
        setAuthorization({token: loginRes.data.token, expireAt: loginRes.data.expireAt})
        if(roles[0].id=="User")
        {
          setAuthorization({token: user.ProjectUUID},AUTH_TYPE.AUTH1)
          this.$router.push('/UserDisplayList/'+user.Uuid)
        }
        else if(roles[0].id=="Operator")
        {
          setAuthorization({token: user.ProjectUUID},AUTH_TYPE.AUTH1)
          this.$router.push('/dashboard')
        }
        else
        {
          this.$router.push('/project')
        }

        this.$message.success(this.$t('loginPage.logonSuccess'), 3)
      } else {
        this.$message.error(this.$t('loginPage.logonFailed'), 3)
      }
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
            _t.$router.push('/project')
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
  //background: linear-gradient(120deg, rgba(84,58,183,1) 0%, rgba(0,172,193,1) 100%);
  color:white;
}

.inner-header {
  height:100%;
  width:100%;
  margin: 0;
  padding: 0;
}

.flex { /*Flexbox for containers*/
  display: flex;
  justify-content: center;
  align-items: center;
  text-align: center;
}

.waves {
  position:relative;
  width: 100%;
  height:15vh;
  margin-bottom:-7px; /*Fix for safari gap*/
  min-height:100px;
  max-height:150px;
}

.content {
  position:relative;
  height:20vh;
  text-align:center;
  background-color: white;
}

/* Animation */

.parallax > use {
  animation: move-forever 25s cubic-bezier(.55,.5,.45,.5)     infinite;
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
/*Shrinking for mobile*/
@media (max-width: 768px) {
  .waves {
    height:40px;
    min-height:40px;
  }
  .content {
    height:30vh;
  }
  h1 {
    font-size:24px;
  }
}
  /*APP安卓下载*/
  .logo{
    position: absolute;
    top:10px;
    left: 10px;
  }
  .logo img{
    width: 70px;
  }
  .boxer {
    background-image:var(--bgImage);
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
    box-shadow: 0px 0px 13px 5px rgba(14,25,80,.2);
    border-radius: 10px;
  }
  .middle1 {
    overflow: hidden;
  }
  .middle_left {
    float: left;
    width: 500px;
  }
  .middle_left img {

  }
  .middle_right {
    width: 500px;
    height: 600px;
    float: right;
    background: #ffffff;
  }
  .mr_box {
    width: 318px;
    margin: 0 auto;
    margin-top: 50px;
  }
  .mr_box form{
    text-align: center;
  }
  .mr_box h2 {
    letter-spacing: 2px;
    margin-bottom: 50px;
    display: block;
    text-align: center;
    color: #13c2c2;
    font-size: 30px;
  }
  input::-webkit-input-placeholder {
    color: #c1c1c1;
  }
  .form-group {
    margin-bottom: 44px;
    position: relative;
  }
  .mar_b {
    margin-bottom:20px;
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
    background: #13c2c2;
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
  /* 记住密码 */
  .f_pass{
    height: 40px;
    line-height: 40px;
    padding: 0 8px;
    margin: 6px 0;
  }
  .font-s {
    position: relative;
  }
  .custom-control {
    position: relative;
    display: block;
    min-height: 1.5rem;
    padding-left: 1.5rem;
  }
  .small, small {
    font-size: 80%;
    font-weight: 400;
    text-align: right;
  }
  input[type=checkbox], input[type=radio] {
    -webkit-box-sizing: border-box;
    box-sizing: border-box;
    padding: 0;
  }
  .custom-control-input {
    position: absolute;
    z-index: -1;
    opacity: 0;
  }
  .custom-control-label{
    font-size: 15px;
  }
  .custom-control-input:not(:disabled):active~.custom-control-label::before {
    color: #fff;
    background-color: #e5ebfa;
    border-color: #e5ebfa;
  }
  .custom-control-input:focus:not(:checked)~.custom-control-label::before {
    border-color: #bac8f3;
  }
  .custom-control-input:focus~.custom-control-label::before {
    -webkit-box-shadow: 0 0 0 0.2rem rgba(78,115,223,.25);
    box-shadow: 0 0 0 0.2rem rgba(78,115,223,.25);
  }
  .custom-control-input:checked~.custom-control-label::before {
    color: #fff;
    border-color: #059df6;
    background-color: #059df6;
  }
  .custom-checkbox .custom-control-label::before {
    border-radius:4px;
  }
  .custom-control-label::before, .custom-file-label, .custom-select {
    -webkit-transition: background-color .15s ease-in-out,border-color .15s ease-in-out,-webkit-box-shadow .15s ease-in-out;
    transition: background-color .15s ease-in-out,border-color .15s ease-in-out,-webkit-box-shadow .15s ease-in-out;
    transition: background-color .15s ease-in-out,border-color .15s ease-in-out,box-shadow .15s ease-in-out;
    transition: background-color .15s ease-in-out,border-color .15s ease-in-out,box-shadow .15s ease-in-out,-webkit-box-shadow .15s ease-in-out;
  }
  .custom-control-label::before {
    position: absolute;
    top:1px;
    right: 64px;
    display: block;
    width:16px;
    height: 16px;
    pointer-events: none;
    content: "";
    background-color: #fff;
    border: #b7b9cc solid 1px;
  }
  .custom-checkbox .custom-control-input:checked~.custom-control-label::after {
    background:url('../../assets/images/duigou.png') no-repeat;
  }
  .custom-control-label::after {
    position: absolute;
    top:5px;
    right: 62px;
    display: block;
    width: 16px;
    height: 16px;
    content: "";
    background: no-repeat 50%/50% 50%;
  }
  /* 密码错误 */
  .f_pass_n{

  }
  .f_pass i{
    margin:11px  6px 0 0;
  }
  .f_pass p{
    color: #dc0a0a;
    text-align: left;
  }
  @media screen and (min-width: 960px) and (max-width:1200px){
    .box1 {
      margin-top: 144px;
      width: 960px;;
    }
    .mr_box{
      margin-top: 40px;
    }
    .mr_box h2 {
      margin-bottom: 26px;
    }
    .middle_left {
      width: 560px;
    }
    .middle_left img {
      width: 100%;
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
    .login {
      margin-top: 20px;
    }
    .form-group {
      margin-bottom: 36px;
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
    .mr_box h2 {
      margin-bottom: 34px;
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
    .login {
      margin-top: 18px;
    }
    .form-group {
      margin-bottom: 28px;
    }
  }
</style>
