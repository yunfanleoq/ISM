<template>
  <a-layout style="background: #fff;height: 100%;">
    <a-layout-header style="position: fixed; z-index: 1000; width: 100%;background: #fff;height: 50px">
      <ISMHeader ref="ISMHeader"  />
    </a-layout-header>
    <a-layout style="margin-top: 50px;height: 100%;">
      <a-layout-sider v-show="leftCollapsed"  width="430" style="background: #fff;position: relative;">
        <a-card :title="$t('displayConfig.WestTitle')" style="background: #fff;height: 100%;overflow-y: auto;">
          <a-layout  style="background: #fff;height: 100%;">
            <a-layout-sider width="45" style="background: #fff;border-right:1px solid #95B8E7;">
              <a-menu
                  :default-selected-keys="['ISMResources']"
                  mode="inline"
                  :style="{'background':'#ffffff'}"
                  :inlineIndent=14
                  @click="sideNavClick"
                  :inline-collapsed="true"
              >
                <a-menu-item key="ISMResources">
                  <a-tooltip placement="right">
                    <template #title>
                      <span>{{$t('displayConfig.PageIcoTitle')}}</span>
                    </template>
                    <icon-font type="icon-iconset0335"  style="font-size: 20px"/>
                  </a-tooltip>

                </a-menu-item>
                <a-menu-item key="ISMToolBox">
                  <a-tooltip placement="right">
                    <template #title>
                      <span>{{$t('displayConfig.ToolIcoTitle')}}</span>
                    </template>
                    <icon-font type="icon-zujianku" style="font-size: 20px"/>
                  </a-tooltip>
                </a-menu-item>
                <a-menu-item key="ISMPageCanvas">
                  <a-tooltip placement="right">
                    <template #title>
                      <span>{{$t('displayConfig.ToolIcoPageCanvasTitle')}}</span>
                    </template>
                    <icon-font type="icon-yemianyuansu"  style="font-size: 20px"/>
                  </a-tooltip>
                </a-menu-item>
                <a-menu-item key="ISMDiyToolBox">
                  <a-tooltip placement="right">
                    <template #title>
                      <span>{{$t('displayConfig.ToolDiyTitle')}}</span>
                    </template>
                    <icon-font type="icon-zu" style="font-size: 20px"/>
                  </a-tooltip>
                </a-menu-item>
                <a-menu-item key="ISMMESToolBox">
                  <a-tooltip placement="right">
                    <template #title>
                      <span>{{$t('displayConfig.ToolMesTitle')}}</span>
                    </template>
                    <icon-font type="icon-MES" style="font-size: 20px"/>
                  </a-tooltip>
                </a-menu-item>
              </a-menu>
            </a-layout-sider>
            <a-layout-content>
              <ISMPageCanvas @drag-start="handleDragStart" ref="ISMPageCanvas"  v-show="selectMenuKey=='ISMPageCanvas'"/>
              <ISMResources ref="ISMResources"  v-show="selectMenuKey=='ISMResources'"/>
              <ISMToolBox @drag-start="handleDragStart" ref="ISMToolBox"  v-show="selectMenuKey=='ISMToolBox'"/>
              <ISMDiyToolBox  @drag-start="handleDragStart" ref="ISMDiyToolBox"  v-show="selectMenuKey=='ISMDiyToolBox'"/>
              <ISMMesToolBox  @drag-start="handleDragStart" ref="ISMMESToolBox"  v-show="selectMenuKey=='ISMMESToolBox'"/>
            </a-layout-content>
          </a-layout>
        </a-card>
        <div class="side-toggle-btn side-toggle-left" @click="setLeftCollapsed(false)" :title="$t('displayConfig.collapse')">
          <a-icon type="double-left" />
        </div>
      </a-layout-sider>
      <a-layout-sider v-show="!leftCollapsed"  width="25" style="background: #fff;border-right:1px solid #95B8E7;">
        <div class="side-toggle-btn side-toggle-left-expand" @click="setLeftCollapsed(true)" :title="$t('displayConfig.expand')">
          <a-icon type="double-right" />
          <span class="side-toggle-label">{{$t('displayConfig.toolbox')}}</span>
        </div>
      </a-layout-sider>


      <a-layout-content class="ISMCanvas" style="background: #fff;height: 100%;">
          <ISMCanvas ref="ISMCanvas"/>
      </a-layout-content>

      <a-layout-sider v-show="rightCollapsed"  width="430" style="background: #fff;position: relative;">
        <a-card :title="$t('displayConfig.EastTitle')" style="background: #fff;height: 100%;overflow-y: auto;">
          <ISMProperties ref="ISMProperties" />
        </a-card>
        <div class="side-toggle-btn side-toggle-right" @click="setRightCollapsed(false)" :title="$t('displayConfig.collapse')">
          <a-icon type="double-right" />
        </div>
      </a-layout-sider>
      <a-layout-sider v-show="!rightCollapsed"  width="25" style="background: #fff;border-left:1px solid #95B8E7;">
        <div class="side-toggle-btn side-toggle-right-expand" @click="setRightCollapsed(true)" :title="$t('displayConfig.expand')">
          <a-icon type="double-left" />
          <span class="side-toggle-label">{{$t('displayConfig.properties')}}</span>
        </div>
      </a-layout-sider>
    </a-layout>
  </a-layout>
</template>
<script>
import ISMHeader from './ISMHeader';
import ISMCanvas from './ISMCanvas';
import ISMToolBox from './ISMToolBox';
import ISMPageCanvas from './ISMPageCanvas';
import ISMResources from './ISMResources';
import ISMMesToolBox from './ISMMesBox';
import ISMDiyToolBox from './ISMDiyBox';
import ISMProperties from './ISMProperties';
import ISMLogger from './ISMLogger';
import {mapMutations, mapState} from "vuex";
import store from "@/store";
export default {
  name: 'ISMDisPlayEditor',
  components: {
    ISMCanvas,
    ISMToolBox,
    ISMPageCanvas,
    ISMProperties,
    ISMHeader,
    ISMResources,
    ISMMesToolBox,
    ISMDiyToolBox
    // ISMLogger
  },
  i18n: require('../../i18n/language'),
  data() {
    return {
      leftCollapsed:true,
      rightCollapsed:false,
      selectMenuKey:"ISMResources"
    }
  },
  methods: {
    ...mapMutations('ISMDisPlayEditorTool', [
      'setLoggerList',
    ]),
    handleDragStart(nodeType,e) {
      this.ISMCavasDND.startDrag(nodeType, e)
    },
    sideNavClick(e){
      this.selectMenuKey = e.key
    },
    setLeftCollapsed(what){
      this.leftCollapsed = what
    },
    setRightCollapsed(what){
      this.rightCollapsed = what
    }
  },
  computed: {
    ...mapState('setting', ['langList','isMobile','lang',]),
    ...mapState({
      ISMCavasDND: state => store.state.ISMDisPlayEditorTool.ISMCavasDND,
    }),
    IsLicense () {
      return this.$store.state.setting.IsOEM
    },
  },
  mounted(){

  },
  created(){
    this.setLoggerList({level:"info",content:"页面初始化完成"})
  }
}

</script>

<style scoped>
::v-deep .ant-layout-header {
  padding: 0 0px;
}
::v-deep  .ant-card-extra {
  padding: 10px 0;
}
::v-deep .ant-card-body {
  padding: 2px;
  height: 100%;
  background: #fff;
}
::v-deep  .ant-card-head-title {
  padding: 10px 0;
}
::v-deep  .ant-card-head {
  padding: 0 5px;
}
.layout-panel-center{
  width: 100%;
  height: 100%;
}
.side-toggle-btn {
  position: absolute;
  top: 0;
  bottom: 0;
  width: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  background: linear-gradient(90deg, transparent, rgba(24, 144, 255, 0.05));
  transition: all 0.3s ease;
  z-index: 10;
}
.side-toggle-btn:hover {
  width: 20px;
  background: linear-gradient(90deg, transparent, rgba(24, 144, 255, 0.15));
}
.side-toggle-btn svg {
  font-size: 14px;
  color: #999;
  transition: all 0.3s ease;
}
.side-toggle-btn:hover svg {
  color: #1890ff;
}
.side-toggle-left {
  right: 0;
  border-radius: 4px 0 0 4px;
}
.side-toggle-right {
  left: 0;
  border-radius: 0 4px 4px 0;
}
.side-toggle-left-expand, .side-toggle-right-expand {
  position: absolute;
  top: 0;
  bottom: 0;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
}
.side-toggle-left-expand:hover, .side-toggle-right-expand:hover {
  background: rgba(24, 144, 255, 0.08);
}
.side-toggle-left-expand svg, .side-toggle-right-expand svg {
  font-size: 16px;
  color: #1890ff;
}
.side-toggle-label {
  font-size: 10px;
  color: #666;
  writing-mode: vertical-rl;
  text-orientation: mixed;
  margin-top: 8px;
  letter-spacing: 2px;
}
</style>