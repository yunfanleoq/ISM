<template>
  <!-- 模板完全不变，和你最初的代码一致 -->
  <a-layout-sider
    :theme="sideTheme"
    :class="['side-menu', 'beauty-scroll', isMobile ? null : 'shadow']"
    width="256px"
    :collapsible="collapsible"
    v-model="collapsed"
    :trigger="null"
  >
    <div :class="['logo', theme]">
      <router-link to="/Project">
        <img :src="systemLogo">
        <h1>{{systemName}}</h1>
      </router-link>
    </div>
    <i-menu
      :theme="theme"
      :collapsed="collapsed"
      :options="menuData"
      @select="onSelect"
      class="menu"
    />
  </a-layout-sider>
</template>

<script>
// JS 逻辑完全不变，和你最初的代码一致
import IMenu from './menu'
import {mapState} from 'vuex'
export default {
  name: 'SideMenu',
  i18n: require('../../i18n/language'),
  components: {IMenu},
  props: {
    collapsible: {
      type: Boolean,
      required: false,
      default: false
    },
    collapsed: {
      type: Boolean,
      required: false,
      default: false
    },
    menuData: {
      type: Array,
      required: true
    },
    theme: {
      type: String,
      required: false,
      default: 'dark'
    }
  },
  computed: {
    sideTheme() {
      return this.theme == 'light' ? this.theme : 'dark'
    },
    systemLogo () {
      return this.$store.state.setting.SystemLogo
    },
    ...mapState('setting', ['isMobile', 'systemName'])
  },
  methods: {
    onSelect (obj) {
      this.$emit('menuSelect', obj)
    }
  }
}
</script>

<style lang="less" scoped>
@import "index";

// 核心色彩体系（不变）
@primary: #13c2c2;
@primary-light: #e6fffa;
@primary-hover: #b5f5ec;
@primary-selected-bg: #0ea6a6;
@text-main: #1e293b;
@text-sub: #64748b;
@border: #e2e8f0;
@white: #ffffff;

// 全局动画定义（不变）
@keyframes hoverLift {
  from { transform: translateY(0); box-shadow: none; }
  to { transform: translateY(-2px); box-shadow: 0 4px 12px rgba(19, 194, 194, 0.15); }
}
@keyframes iconShift {
  from { transform: translateX(0); }
  to { transform: translateX(3px); }
}
@keyframes arrowRotate {
  from { transform: rotate(0); }
  to { transform: rotate(180deg); }
}

// 1. 侧边栏容器基础样式（不变）
::v-deep .side-menu {
  width: 256px !important;
  min-width: 256px !important;
  max-width: 256px !important;
  background-color: @white;
  border-right: 1px solid @border;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow-x: hidden !important;
  overflow-y: auto;

  &.ant-layout-sider-collapsed {
    .logo h1 {
      opacity: 0;
      width: 0;
      transition: all 0.2s ease;
    }
    .ant-menu-item,
    .ant-menu-submenu-title {
      padding: 0 28px !important;
    }
  }
}

// 2. Logo 区域（不变）
.side-menu .logo {
  height: 50px;
  display: flex;
  align-items: center;
  padding: 0 5px;
  background-color: @white;
  border-bottom: 1px solid @border;
  transition: all 0.3s ease;
  overflow: hidden;
  border-radius: 0px;
  &:hover {
    background-color: @primary-light;
  }

  img {
    height: 30px;
    width: auto;
    object-fit: contain;
    transition: transform 0.3s ease;
  }

  &:hover img {
    transform: scale(1.08);
  }

  h1 {
    margin: 0 0 0 10px;
    font-size: 18px;
    color:#0d0d0d !important;
    font-weight: 700;
    text-align: center;
    font-family: '黑体', 'Microsoft YaHei', sans-serif;
    white-space: nowrap;
    transition: all 0.3s ease;
    line-height: 1;
  }
  &:hover h1 {
    transform: scale(1.08);
  }
}

// 3. 菜单容器（不变）
::v-deep .menu {
  padding: 1px 0;
  background-color: @white;
  overflow-x: hidden !important;
  overflow-y: auto;
  max-height: calc(100vh - 60px);

  // 滚动条美化（不变）
  &::-webkit-scrollbar {
    width: 1px;
    height: 0;
  }
  &::-webkit-scrollbar-thumb {
    border-radius: 3px;
    background: @primary;
    opacity: 0.3;
  }
  &::-webkit-scrollbar-track {
    background: @primary-light;
  }
  /* 新增：强制刷新选中样式，解决首次加载不生效问题 */
  &.ant-menu {
    &-light {
      &.ant-menu {
        &-vertical,
        &-vertical-left,
        &-inline {
          > .ant-menu-item-selected,
          > .ant-menu-submenu > .ant-menu-submenu-title-selected {
            // 强制覆盖选中样式
            background: linear-gradient(90deg, @primary-selected-bg 0%, @primary 100%) !important;
            color: @white !important;
          }
        }
      }
    }
  }
}

// 4. 关键修复：给样式加“父容器类”，权重超越组件库默认样式
// 原代码是 .ant-menu-light .ant-menu-item，现在加 .menu 父类，变成 .menu .ant-menu-light .ant-menu-item
::v-deep .menu .ant-menu-light {
  // 4.1 菜单项基础样式（现在权重足够，首次加载就生效）
  .ant-menu-item,
  .ant-menu-submenu-title {
    height: 44px;
    margin: 6px 10px;
    padding: 0 18px;
    font-size: 14px;
    line-height: 44px;
    border-radius: 8px;
    color: @text-main;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  // 4.2 菜单项图标/箭头基础样式（首次加载生效）
  .ant-menu-item .anticon,
  .ant-menu-submenu-title .anticon {
    font-size: 18px;
    color: @text-sub;
    margin-right: 14px;
    transition: all 0.3s ease;
  }
  .ant-menu-submenu-arrow {
    color: @text-sub;
    font-size: 16px;
    transition: all 0.3s ease;
  }

  // 4.3 鼠标悬停效果（不变，权重依然足够）
  .ant-menu-item:hover,
  .ant-menu-submenu-title:hover {
    background: linear-gradient(90deg, @primary-hover 0%, @primary-light 100%);
    color: @primary !important;

    .anticon {
      color: @primary !important;
    }

    .ant-menu-submenu-arrow {
      color: @primary !important;
    }
  }

  // 4.4 选中效果（不变，权重依然足够）
  .ant-menu-item-selected,
  .ant-menu-submenu-title-selected {
    &,
    &:active,
    &:focus {
      background: linear-gradient(90deg, @primary-selected-bg 0%, @primary 100%) !important;
      color: @white !important;
      font-weight: 600 !important;
      box-shadow: 0 4px 16px rgba(14, 166, 166, 0.25) !important;
      border-left: 3px solid @white !important;
    }

    .anticon,
    & > span,
    & > span > span,
    & > a,
    & > i,
    .ant-menu-submenu-arrow {
      color: @white !important;
      margin-right: 14px;
    }
  }
}

// 5. 子菜单展开动画（不变）
::v-deep .ant-menu-sub {
  animation: subMenuExpand 0.2s ease forwards;
  background-color: @white !important;
  border-left: 2px solid @primary;
  border-radius: 0 8px 8px 0;
  padding-left: 10px !important;
  width: 100% !important;
}

@keyframes subMenuExpand {
  from { opacity: 0; transform: translateX(-5px); }
  to { opacity: 1; transform: translateX(0); }
}

// 6. 移动端适配（不变）
::v-deep .is-mobile .side-menu {
  position: fixed;
  z-index: 1000;
  height: 100vh;
  width: 256px !important;
  max-width: 256px !important;
  animation: mobileSlideIn 0.3s ease forwards;
  box-shadow: 2px 0 20px rgba(19, 194, 194, 0.15);
  overflow-x: hidden !important;
}

// 4. 菜单项基础样式（不变）
::v-deep .ant-menu-item .anticon,
.ant-menu-submenu-title .anticon {
  font-size: 14px;
  color: @text-sub;
  margin-right: 10px;
  transition: all 0.3s ease;
}

::v-deep .ant-menu-vertical .ant-menu-item,
.ant-menu-vertical-left .ant-menu-item,
.ant-menu-vertical-right .ant-menu-item,
.ant-menu-inline .ant-menu-item,
.ant-menu-vertical .ant-menu-submenu-title,
.ant-menu-vertical-left .ant-menu-submenu-title,
.ant-menu-submenu .ant-menu-vertical-right .ant-menu-submenu-title,
.ant-menu-inline .ant-menu-submenu-title {
  height: 44px;
  padding: 0 18px;
  font-size: 14px;
  line-height: 44px;
  border-radius: 0px;
  color: @text-main;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

// 5. 子菜单箭头（不变）
::v-deep .ant-menu-submenu-arrow {
  color: @text-sub;
  font-size: 16px;
  transition: all 0.3s ease;
}

// 6. 鼠标悬停效果（不变）
::v-deep .ant-menu-light .ant-menu-item:hover,
::v-deep .ant-menu-light .ant-menu-submenu-title:hover {
  background: linear-gradient(90deg, @primary-hover 0%, @primary-light 100%);
  color: @primary !important;

  .anticon {
    color: @primary !important;
  }

  .ant-menu-submenu-arrow {
    color: @primary !important;
  }
}

// 7. 选中效果（强化优先级，解决首次加载问题）
::v-deep .ant-menu-light .ant-menu-item-selected,
::v-deep .ant-menu-light .ant-menu-submenu-title-selected {
  /* 核心优化：通过多层父类限定提升优先级，确保首次加载生效 */
  &,
  &:active,
  &:focus {
    background: linear-gradient(90deg, @primary-selected-bg 0%, @primary 100%) !important;
    color: @white !important;
    font-weight: 600 !important;
    box-shadow: 0 4px 16px rgba(14, 166, 166, 0.25) !important;
    border-left: 0px solid @white !important;
  }

  // 强制覆盖所有子元素颜色
  .anticon,
  & > span,
  & > span > span,
  & > a,
  & > i,
  .ant-menu-submenu-arrow {
    color: @white !important;
    margin-right: 14px;
  }
}

@keyframes mobileSlideIn {
  from { transform: translateX(-100%); }
  to { transform: translateX(0); }
}

// =================================================================
// == 仅新增：鼠标悬停左侧边框（不改变任何原有布局） ==
// =================================================================
::v-deep .ant-menu-light .ant-menu-item:hover,
::v-deep .ant-menu-light .ant-menu-submenu-title:hover,
::v-deep .ant-menu-sub .ant-menu-item:hover {
  box-shadow: -3px 0 0 0 @primary inset !important; /* 内部阴影实现边框，不占布局空间 */
}
</style>