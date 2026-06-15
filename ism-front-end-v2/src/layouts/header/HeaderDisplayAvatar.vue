<template>
  <a-dropdown>
    <div class="header-avatar" style="cursor: pointer">
      <a-avatar class="avatar" size="small" shape="circle" :src="user.avatar"/>
      <span class="name">{{user.name}}</span>
    </div>
    <a-menu :class="['avatar-menu']" slot="overlay">
      <template v-if="is_project">
        <a-menu-item >
          <router-link to="/Setting/Account">
            <a-icon type="user" />
            <span style="margin-left: 5px">{{$t('header.Profile')}}</span>
          </router-link>

        </a-menu-item>
        <a-menu-divider />
      </template>
      <a-menu-item @click="logout">
        <a-icon style="margin-right: 8px;" type="poweroff" />
        <span>{{$t('header.Logout')}}</span>
      </a-menu-item>
    </a-menu>
  </a-dropdown>
</template>

<script>
import {mapGetters} from 'vuex'
import {logout} from '@/services/user'
import { Icon } from 'ant-design-vue';
import Cookie from 'js-cookie'

export default {
  name: 'HeaderDisplayAvatar',
  i18n: require('../../i18n/language'),
  computed: {
    ...mapGetters('account', ['user']),
  },
  created(){
  },
  data() {
    return {
      is_project: false
    }
  },
  methods: {
    logout() {
      logout()
      this.$router.push('/login')
      Cookie.remove('ProjectUuid')
    },
    ClearProjectInfo() {

    }
  }
}
</script>

<style lang="less">
  .header-avatar{
    display: inline-flex;
    .avatar, .name{
      align-self: center;
    }
    .avatar{
      margin-right: 8px;
    }
    .name{
      font-weight: 500;
    }
  }
  .avatar-menu{
    width: 150px;
  }

</style>
