<template>
  <common-layout>
    <div class="top">
      <div class="header">
        <img alt="logo" class="logo" :src="systemLogo" />
      </div>
    </div>
    <div class="login">
      <a-form @submit="onSubmit" :form="form">
        <a-tabs size="large" :tabBarStyle="{textAlign: 'center'}" style="padding: 0 2px;">
          <a-tab-pane tab="登录" key="1">
            <a-alert type="error" :closable="true" v-show="error" :message="error" showIcon style="margin-bottom: 24px;" />
            <a-form-item>
              <a-input
                  class="form-control"
                  autocomplete="autocomplete"
                  size="large"
                  v-decorator="['name', {rules: [{ required: true, message:  $t('loginPage.logonUserTips'), whitespace: true}]}]"
              >
                <a-icon slot="prefix" type="user" style="position: absolute;right:-15px;color: #13c2c2;font-size: 20px;"/>
              </a-input>
            </a-form-item>
            <a-form-item>
              <a-input
                  class="form-control"
                  size="large"
                  autocomplete="autocomplete"
                  type="password"
                  v-decorator="['password', {rules: [{ required: true, message: $t('loginPage.logonPasswordTips'), whitespace: true}]}]"
              >
                <a-icon slot="prefix" type="lock" size="large" style="position: absolute;right:-15px;color: #13c2c2;font-size: 20px;"/>
              </a-input>
            </a-form-item>
          </a-tab-pane>
        </a-tabs>
        <a-form-item>
          <a-button :loading="logging" class="login" style="width: 100%;margin-top: 2px" size="large" htmlType="submit" type="primary">{{$t('loginPage.logonBtn')}}</a-button>
        </a-form-item>
      </a-form>
      <div style="margin-top: 12px; text-align: center;">
        <a @click="$router.push('/login')">使用桌面版登录</a>
      </div>
      <div>
        <a-breadcrumb>
          <a-breadcrumb-item   v-for=" lang in langList" :key="lang.key" ><a @click="setLang(lang.key)">{{lang.name}}</a></a-breadcrumb-item>
        </a-breadcrumb>
      </div>
    </div>
  </common-layout>
</template>

<script>
import CommonLayout from '@/layouts/CommonLayout'
import {login} from '@/services/user'
import {AUTH_TYPE, setAuthorization} from '@/utils/request'
import {loadRoutes} from '@/utils/routerUtil'
import {mapMutations, mapGetters} from 'vuex'
import {mapState} from 'vuex'
import md5 from 'js-md5'
import {ProjectList} from '@/services/project'
import {applyHomeProjectAuth} from '@/config/homeDashboard'
export default {
  name: 'LoginPhone',
  i18n: require('@/i18n/language'),
  components: {CommonLayout},
  data () {
    return {
      logging: false,
      error: '',
      form: this.$form.createForm(this)
    }
  },
  created(){
    if(!this.isMobile)
    {
      this.$router.push('/login')
    }
    else
    {
     let phoneUser =  localStorage.getItem("phoneUser")
     let phonePassword =  localStorage.getItem("phonePassword")
      if((phoneUser)&&(phonePassword))
      {
        this.autoLogin(phoneUser,phonePassword)
      }
    }
  },
  computed: {
    ...mapState('setting', ['langList','isMobile','lang',]),
    ...mapGetters('setting', ['homeDashboardPath']),
    systemName () {
      return this.$store.state.setting.systemName
    },
    systemLogo () {
      return this.$store.state.setting.SystemLogo
    },
  },
  methods: {
    ...mapMutations('setting', ['setLang']),
    ...mapMutations('account', ['setUser', 'setPermissions', 'setRoles', 'setRoutesConfig']),
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
          localStorage.setItem("phoneUser",Username)
          localStorage.setItem("phonePassword",password)
          login(Username, password).then(this.afterLogin).catch(function(){
            _t.logging = false
            localStorage.removeItem("phoneUser")
            localStorage.removeItem("phonePassword")
            _t.$message.error(_t.$t('loginPage.serverError'), 3)
          })
        }
      })
    },
    afterLogin(res) {
      this.logging = false
      const loginRes = res.data
      if (loginRes.code == 1000) {
        const {user, roles} = loginRes.data
        this.setUser(user)
        this.setRoles(roles)
        this.setRoutesConfig(user.Menu)
        setAuthorization({token: loginRes.data.token, expireAt: loginRes.data.expireAt})
        localStorage.setItem('LoginFrom', '/loginPhone')
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
        localStorage.removeItem('phoneUser')
        localStorage.removeItem('phonePassword')
        this.$message.error(this.$t('loginPage.logonFailed'), 3)
      }
    },
    // 与 Login.vue 保持一致：Admin/单项目自动进首页大屏，多项目进项目选择页
    enterAfterAuth(roles, user) {
      const _t = this
      const roleId = (roles && roles[0]) ? roles[0].id : ''
      const goHome = () => {
        applyHomeProjectAuth(_t.$store)
        _t.$router.push(_t.homeDashboardPath)
      }
      if (roleId === 'User' || roleId === 'Operator') {
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
        _t.$router.push('/project')
      })
    }
  }
}
</script>

<style lang="less" scoped>
  .common-layout{
    .top {
      text-align: center;
      .header {
        height: 44px;
        line-height: 44px;
        a {
          text-decoration: none;
        }
        .logo {
          height: 44px;
          vertical-align: top;
          margin-right: 16px;
        }
        .title {
          font-size: 33px;
          color: @title-color;
          font-family: 'Myriad Pro', 'Helvetica Neue', Arial, Helvetica, sans-serif;
          font-weight: 600;
          position: relative;
          top: 2px;
        }
      }
      .desc {
        font-size: 14px;
        color: @text-color-second;
        margin-top: 12px;
        margin-bottom: 40px;
      }
    }
    .login{
      width: 368px;
      margin: 0 auto;
      @media screen and (max-width: 576px) {
        width: 95%;
      }
      @media screen and (max-width: 320px) {
        .captcha-button{
          font-size: 14px;
        }
      }
      .icon {
        font-size: 24px;
        color: @text-color-second;
        margin-left: 16px;
        vertical-align: middle;
        cursor: pointer;
        transition: color 0.3s;

        &:hover {
          color: @primary-color;
        }
      }
    }
  }
</style>
