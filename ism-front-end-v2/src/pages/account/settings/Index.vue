<template>
  <div class="page-header-index-wide">
    <a-card :bordered="false" :bodyStyle="{ padding: '16px 0', height: '100%' }" :style="{ height: '100%' }">
      <div class="account-settings-info-main" :class="{ 'mobile': isMobile }">
        <div class="account-settings-info-left">
          <a-menu
              :class="['avatar-menu']"
            :mode="isMobile ? 'horizontal' : 'inline'"
            :style="{ border: '0', width: isMobile ? '560px' : 'auto'}"
            :selectedKeys="selectedKeys"
            type="inner"
          >
            <a-menu-item :class="openKeys=='basic'?'ant-menu-item ant-menu-item-selected':''" key="/account/settings/basic" @click="onOpenChange('basic')">
                {{ $t('account.settings.menuMap.basic') }}
            </a-menu-item>
            <a-menu-item :class="openKeys=='security'?'ant-menu-item ant-menu-item-selected':''" @click="onOpenChange('security')">
                {{ $t('account.settings.menuMap.security') }}
            </a-menu-item>
          </a-menu>
        </div>
        <div class="account-settings-info-right">
          <basic-setting v-if="openKeys=='basic'"></basic-setting>
          <security-setting v-if="openKeys=='security'"></security-setting>
        </div>
      </div>
    </a-card>
  </div>



</template>

<script>

import BasicSetting from './BasicSetting'
import SecuritySetting from './Security'
export default {
  i18n: require('../../../i18n/language'),
  data () {
    return {
      // horizontal  inline
      mode: 'inline',
      page: {},
      openKeys: "basic",
      selectedKeys: [],
      isMobile:false,
      // cropper
      preview: {},
      option: {
        img: '/avatar2.jpg',
        info: true,
        size: 1,
        outputType: 'jpeg',
        canScale: false,
        autoCrop: true,
        // 只有自动截图开启 宽度高度才生效
        autoCropWidth: 180,
        autoCropHeight: 180,
        fixedBox: true,
        // 开启宽度和高度比例
        fixed: true,
        fixedNumber: [1, 1]
      },

      pageTitle: ''
    }
  },
  components: {
    BasicSetting,
    SecuritySetting
  },
  mounted () {
    this.updateMenu()
  },
  methods: {
    onOpenChange (openKeys) {
     this.openKeys = openKeys
    },
    updateMenu () {
      const routes = this.$route.matched.concat()
      this.selectedKeys = [ routes.pop().path ]
    }
  },
  watch: {
    '$route' (val) {
      this.updateMenu()
    }
  }
}
</script>

<style lang="less" scoped>
// 主容器样式补充
.account-settings-info-main {
  display: flex;
  width: 100%;
  height: 100%;
  overflow: hidden;

  // 移动设备布局调整
  &.mobile {
    flex-direction: column;
    padding: 0 16px;
  }
}
// 左侧菜单容器样式
.account-settings-info-left {
  flex: 0 0 220px;
  padding: 0 16px;
  border-right: 1px solid #f0f2f5;
  box-sizing: border-box;

  // 移动设备菜单容器样式
  &.mobile {
    flex: none;
    border-right: none;
    border-bottom: 1px solid #f0f2f5;
    padding: 0 0 16px;
    margin-bottom: 16px;
    width: 100% !important;
  }
}
// 右侧内容区域样式
.account-settings-info-right {
  flex: 1;
  padding: 0 24px;
  overflow-y: auto;
  box-sizing: border-box;
  height: 100%;

  // 移动设备内容区域样式
  &.mobile {
    padding: 0;
    height: calc(100% - 60px);
  }
}

// 原有footer样式保留
.footer {
  padding: 20px 20px 15px; /* 大幅降低整体高度 */
  text-align: center;
  background-color: #f9fafb; /* 浅灰底色，比纯白更显质感 */
  border-top: 1px solid #f0f2f5; /* 细边框分隔 */

  .links {
    margin-bottom: 12px; /* 减少链接区与版权区间距 */

    a {
      color: #4e5969; /* 低饱和深灰，柔和不刺眼 */
      font-size: 14px;
      margin: 0 15px; /* 适当缩减链接间距 */
      text-decoration: none;
      transition: all 0.2s ease;

      &:hover {
        color: #2c83d8; /* 低饱和蓝色，hover反馈温和 */
      }

      .anticon {
        margin-right: 5px;
        vertical-align: middle;
        font-size: 15px;
      }
    }
  }

  .copyright {
    color: #86909c; /* 中灰，区分层级 */
    font-size: 13px;
    line-height: 1.5; /* 紧凑行高 */
    margin: 5px 0;

    a {
      color: #6b7785; /* 比正文稍深，突出可点击性 */
      text-decoration: none;
      transition: color 0.2s ease;

      &:hover {
        color: #2c83d8;
      }
    }

    .anticon {
      margin: 0 4px;
      vertical-align: middle;
      color: #86909c;
    }
  }

  .copyright:last-child {
    margin-top: 8px;
    font-size: 12px;
    color: #aeb4bb; /* 浅灰，弱化最底部信息 */
  }
}
</style>
