<template>
  <div class="left-panel">
    <!-- Tab 切换 -->
    <div class="panel-tabs">
      <div
        class="tab-item"
        :class="{active: activeTab === 'components'}"
        @click="activeTab = 'components'"
      >
        <i class="fas fa-th"></i> {{ $t('ISM3DEditor.componentLibrary') }}
      </div>
      <div
        class="tab-item"
        :class="{active: activeTab === '2d'}"
        @click="activeTab = '2d'"
      >
        <i class="fas fa-layer-group"></i> 2D
      </div>
      <div
        class="tab-item"
        :class="{active: activeTab === 'scene'}"
        @click="activeTab = 'scene'"
      >
        <i class="fas fa-sitemap"></i> {{ $t('ISM3DEditor.sceneTree') }}
      </div>
    </div>

    <!-- 组件库 Tab -->
    <div v-show="activeTab === 'components'" class="tab-content component-tab-layout">
      <div class="component-side-nav">
        <a-tooltip
          v-for="nav in componentNavs"
          :key="nav.key"
          placement="right"
        >
          <template #title>
            <span>{{ nav.label }}</span>
          </template>
          <button
            type="button"
            class="component-side-nav-item"
            :class="{ active: currentComponentNav === nav.key }"
            @click="selectComponentNav(nav.key)"
          >
            <i :class="nav.icon"></i>
          </button>
        </a-tooltip>
      </div>
      <div class="component-main-pane">
        <!-- 普通组件库 -->
        <template v-if="currentComponentNav !== 'external' && currentComponentNav !== 'local'">
        <div class="component-search">
          <a-input v-model="searchComp" placeholder="搜索 3D 组件..." size="small" :allow-clear="true">
            <a-icon slot="prefix" type="search" style="color:#555588" />
          </a-input>
        </div>
        <div v-if="filteredComponents.length===0" style="color:#555588;text-align:center;padding:20px;font-size:12px">
          无匹配组件
        </div>
        <a-collapse
          v-show="filteredComponents.length > 0"
          class="component-library-collapse"
          :bordered="false"
          :active-key="activeComponentGroupKeys"
          expand-icon-position="left"
          @change="onComponentGroupChange"
        >
          <template v-for="(cat, index) in filteredComponents">
            <a-collapse-panel :header="cat.name" :key="cat.name" :style="customStyle">
              <div class="ism-toolbox component-library-group">
                <div class="toolbox-group scrollbar">
                  <div
                    v-for="item in cat.items"
                    :key="item.type"
                    class="toolbox-item"
                    :class="'color-'+item.color"
                    draggable="true"
                    @dragstart="onDragStart($event, item)"
                    @click="$emit('add-object', item)"
                    :title="item.label + (item.description ? ' - ' + item.description : '') + ' - 单击或拖拽添加'"
                  >
                    <div class="toolbox-item-icon">
                      <i :class="item.icon" style="font-size:28px"></i>
                    </div>
                    <div class="toolbox-item-text">{{ item.label }}</div>
                  </div>
                </div>
              </div>
            </a-collapse-panel>
          </template>
        </a-collapse>
        </template>

        <!-- Sketchfab 外部资产 -->
        <SketchfabBrowser
          v-if="currentComponentNav === 'external'"
          @add-object="$emit('add-object', $event)"
        />
        <!-- 本地模型库 -->
        <LocalModelBrowser
          v-if="currentComponentNav === 'local'"
          @add-object="$emit('add-object', $event)"
        />
      </div>
    </div>

    <!-- 2D组件 Tab -->
    <div v-show="activeTab === '2d'" class="tab-content-2d-wrapper">
      <a-layout style="height: 100%;">
        <a-layout-sider width="45" style="background: #fff;border-right:1px solid #95B8E7;">
          <a-menu
              :default-selected-keys="['toolBox']"
              mode="inline"
              :style="{'background':'#ffffff'}"
              :inlineIndent=14
              @click="sideNavClick"
              :inline-collapsed="true"
          >
            <a-menu-item key="toolBox">
              <a-tooltip placement="right">
                <template #title>
                  <span>{{$t('displayConfig.ToolIcoTitle')}}</span>
                </template>
                <icon-font type="icon-zujianku" style="font-size: 20px"/>
              </a-tooltip>
            </a-menu-item>
            <a-menu-item key="pageCanvas">
              <a-tooltip placement="right">
                <template #title>
                  <span>{{$t('displayConfig.ToolIcoPageCanvasTitle')}}</span>
                </template>
                <icon-font type="icon-yemianyuansu" style="font-size: 20px"/>
              </a-tooltip>
            </a-menu-item>
            <a-menu-item key="diy">
              <a-tooltip placement="right">
                <template #title>
                  <span>{{$t('displayConfig.ToolDiyTitle')}}</span>
                </template>
                <icon-font type="icon-zu" style="font-size: 20px"/>
              </a-tooltip>
            </a-menu-item>
            <a-menu-item key="mes">
              <a-tooltip placement="right">
                <template #title>
                  <span>{{$t('displayConfig.ToolMesTitle')}}</span>
                </template>
                <icon-font type="icon-MES" style="font-size: 20px"/>
              </a-tooltip>
            </a-menu-item>
          </a-menu>
        </a-layout-sider>
        <a-layout-content class="toolbox-content-inner">
          <div class="component-search">
            <a-input v-model="search2D" placeholder="搜索 2D 组件..." size="small" :allow-clear="true">
              <a-icon slot="prefix" type="search" style="color:#555588" />
            </a-input>
          </div>
          <a-upload
              v-if="currentNav === 'diy'"
              style="margin-left: 20px; margin-bottom: 10px"
              name="file"
              :multiple="false"
              :action="localUpgradeUrl"
              :showUploadList="false"
              :beforeUpload="beforeUpload"
              @change="localUpgradeCharge"
          >
            <a-button> <a-icon type="upload" />  {{$t('displayConfig.uploadDiy')}}</a-button>
          </a-upload>
          <div v-if="currentGroups.length === 0" style="color:#999;text-align:center;padding:20px;font-size:12px">
            暂无组件
          </div>
          <a-collapse v-show="currentGroups.length > 0" accordion :bordered="false" default-active-key="1" expand-icon-position="left">
            <template v-for="(group, index1) in currentGroups">
              <a-collapse-panel :header="$t(group.title)" :key="index1" :style="customStyle">
                <div class="ism-toolbox" style="border-right:0px solid #95B8E7;height: 300px; overflow: hidden;overflow-y: auto;">
                  <div class="toolbox-group scrollbar">
                    <template v-for="(value, index) in validResultsList(group.items)">
                      <div
                          class="toolbox-item"
                          v-bind:key="index"
                          draggable="true"
                          @dragstart="onDragStart($event, value)"
                          @click="$emit('add-object', value)"
                          :title="group.isSequence ? $t(group.title)+(index+1) : $t(value.text)"
                      >
                        <template>
                          <div class="toolbox-item-icon">
                            <div v-if="value.icon && value.icon.indexOf('icon-')!=-1">
                              <icon-font :type="value.icon" :style="{ fontSize: '32px' }"/>
                            </div>
                            <div v-else-if="value.icon && value.icon.indexOf('svg-')!=-1">
                              <img alt="icon" :src="value.icon" style="width:32px;height: 32px"/>
                            </div>
                            <div v-else>
                              <img alt="icon" :src="value.icon" @click="clickImg($event)" style="width:32px;height: 32px"/>
                            </div>
                          </div>
                          <div class="toolbox-item-text" v-if="group.isSequence">{{$t(group.title)}}{{index+1}}</div>
                          <div class="toolbox-item-text" v-else>{{$t(value.text)}}</div>
                        </template>
                      </div>
                    </template>
                  </div>
                </div>
                <a-icon slot="extra" v-if="group.items.length>100" type="double-right" @click="ShowMoreItem(group)" />
              </a-collapse-panel>
            </template>
          </a-collapse>
        </a-layout-content>
      </a-layout>
    </div>

    <a-modal :visible="itemMore"
            :title="$t(ShowMoreItemArray.title)"
            :dialogStyle="{width:'300px',height:'500px','z-index':-1,left: '352px', top:'70px', position:'absolute'}"
            :iconCls="ShowMoreItemArray.icon"
             :footer="null"
             @cancel="itemMore=false"
             :destroyOnClose="true"
             :maskClosable="false"
             :maskStyle="{}"
            :mask="false">
      <div class="ism-toolbox " style="border-right:0px solid #95B8E7;" >
        <div class="toolbox-group scrollbar" v-if="(Array.isArray(ShowMoreItemArray.items))&&(ShowMoreItemArray.items.length>0)">
          <template v-for="(value,index) in MoreItemList(ShowMoreItemArray.items)" >
            <div
                class="toolbox-item"
                v-bind:key="index"
                draggable="true"
                @dragstart="onModalDragStart($event,value)"
                @click="$emit('add-object', value)"
            >
              <template >
                <div class="toolbox-item-icon">
                  <div v-if="value.icon && value.icon.indexOf('icon-')!=-1">
                    <icon-font :type="value.icon" :style="{ fontSize: '32px' }"/>
                  </div>
                  <div v-else>
                    <img alt="icon" :src="value.icon" style="width:32px;height: 32px"/>
                  </div>
                </div>
                <div class="toolbox-item-text" v-if="ShowMoreItemArray.isSequence"> {{$t(ShowMoreItemArray.title)}}{{index+21}}</div>
                <div class="toolbox-item-text" v-else-if="value.text"> {{$t(value.text)}}</div>
              </template>
            </div>
          </template>
        </div>
      </div>
    </a-modal>

    <div v-show="activeTab === 'scene'" class="tab-content">
      <div class="scene-object-list">
        <div
          v-for="(obj, index) in sceneObjects"
          :key="obj.id"
          class="scene-obj-item"
          :class="{selected: selectedId===obj.id, 'is-background': obj.isBackground}"
          @click="$emit('select', obj.id)"
        >
          <i :class="obj.isBackground ? 'fas fa-image' : (obj.icon || 'fas fa-cube')"></i>
          <span class="obj-name">{{ obj.name }}</span>
          <span v-if="obj.isBackground" class="obj-bg-tag">背景</span>
          <span class="obj-actions">
            <span class="obj-action-btn" @click.stop="$emit('move-up', obj.id)" :title="$t('ISM3DEditor.moveUp')"><i class="fas fa-arrow-up"></i></span>
            <span class="obj-action-btn" @click.stop="$emit('move-down', obj.id)" :title="$t('ISM3DEditor.moveDown')"><i class="fas fa-arrow-down"></i></span>
            <span class="obj-del" @click.stop="$emit('delete', obj.id)"><i class="fas fa-times"></i></span>
          </span>
        </div>
        <div v-if="sceneObjects.length===0" style="color:#444466;text-align:center;padding:14px;font-size:12px">
          {{ $t('ISM3DEditor.sceneEmpty') }}
        </div>
      </div>
    </div>

  </div>
</template>

<script>
import { COMPONENT_LIBRARY } from '../Objects/IndustrialObjects'
import store from '@/store'
import { TOOLBOX_GROUPS, getISMComponent } from '@/pages/ISMDisPlay/componentRegistry'
import Conduit from '@/pages/ISMDisPlay/ISMComponents/Conduit/config.json'
import PageNavigation from '@/pages/ISMDisPlay/ISMComponents/page/navigation.json'
import { DIYPICUPLOAD } from '@/services/api'
import { GetUserCustomPel } from '@/services/system'
import SketchfabBrowser from './SketchfabBrowser.vue'
import LocalModelBrowser from './LocalModelBrowser.vue'

import { getAssetList, getCategories, filterAssets } from '@/pages/ISM3DEditor/utils/SketchfabAssets'

export default {
  name: 'ToolboxPanel',
  components: {
    SketchfabBrowser,
    LocalModelBrowser,
  },
  props: {
    sceneObjects: { type: Array, default: () => [] },
    selectedId: { type: String, default: null },
    showGridSettings: { type: Boolean, default: false },
    gridSettings: { type: Object, default: () => ({}) },
  },
  data() {
    return {
      searchComp: '',
      search2D: '',
      activeComponentGroupKeys: COMPONENT_LIBRARY.length ? [COMPONENT_LIBRARY[0].name] : [],
      currentComponentNav: 'industrial',
      componentLibrary: COMPONENT_LIBRARY,
      activeTab: 'components', // 'components' 或 'scene'
      currentNav: 'toolBox',
      customStyle: 'background: #fff;border-radius: 4px;margin-bottom: 1px;border: 0;overflow: hidden',
      localUpgradeUrl: DIYPICUPLOAD,
      messageShowLoad: false,
      itemMore: false,
      ShowMoreItemArray: [],
      ShowMoreModel: false,
      externalSearch: '',
      externalCategory: 'all',
      externalAssets: [],
      externalCategories: [],
    }
  },
  created() {
    this.initToolboxData()
  },
  mounted() {
    this.initToolboxData()
  },
  i18n: require('@/i18n/language'),
  watch: {
    searchComp: function(val) {
      if (val) {
        this.activeComponentGroupKeys = this.filteredComponents.map(function(cat) { return cat.name })
      } else if (!this.activeComponentGroupKeys.length && this.componentLibrary.length) {
        this.activeComponentGroupKeys = [this.componentLibrary[0].name]
      }
    }
  },
  computed: {
    componentNavs() {
      return [
        { key: 'industrial', label: '工业设备', icon: 'fas fa-layer-group', groups: ['泵阀管件', '容器储罐', '换热分离', '输送存储', '仪器仪表', '管道支架', '其他设备', '工业流水线', '工业扩展'] },
        { key: 'basic', label: '基础组件', icon: 'fas fa-plus', groups: ['基础组件', '基础几何'] },
        { key: 'city', label: '城市设施', icon: 'fas fa-building', groups: ['智慧城市', '仓储物流', '交通设施', '室内设施', '消防安防', '绿化景观'] },
        { key: 'domain', label: '行业场景', icon: 'fas fa-seedling', groups: ['数字孪生', '能源电力', '医疗设施', '农业设施', '水利水务', '运动休闲', '实验室'] },
        { key: 'effect', label: '灯光特效', icon: 'fas fa-magic', groups: ['灯光', '视觉特效'] },
        { key: 'external', label: '外部资产', icon: 'fas fa-globe', groups: [] },
        { key: 'local', label: '本地模型', icon: 'fas fa-hdd', groups: [] },
      ]
    },
    currentComponentNavConfig() {
      var navs = this.componentNavs || []
      return navs.find(nav => nav.key === this.currentComponentNav) || navs[0]
    },
    navComponentLibrary() {
      var cfg = this.currentComponentNavConfig
      if (!cfg || !cfg.groups) return this.componentLibrary
      return this.componentLibrary.filter(function(group) {
        return cfg.groups.indexOf(group.name) !== -1
      })
    },
    currentGroups() {
      let groups = []
      const state = store.state.ISMDisPlayEditorTool

      switch (this.currentNav) {
        case 'toolBox':
          groups = state.toolBoxList || TOOLBOX_GROUPS || []
          break
        case 'pageCanvas':
          groups = state.PageCanVasList || []
          break
        case 'mes':
          groups = state.MesComponentsList || []
          break
        case 'diy':
          groups = state.DiyComponentsList || []
          break
      }

      return groups
    },
    allToolboxGroups() {
      const state = store.state.ISMDisPlayEditorTool
      const groups = []
      ;(state.toolBoxList || []).forEach(group => groups.push(group))
      ;(state.PageCanVasList || []).forEach(group => groups.push(group))
      ;(state.MesComponentsList || []).forEach(group => groups.push(group))
      ;(state.DiyComponentsList || []).forEach(group => groups.push(group))
      return groups
    },
    normalized2DComponents() {
      var list = []
      var self = this
      ;(this.allToolboxGroups || []).forEach(function(group, groupIndex) {
        if (!group || !group.items || group.items.length === 0) return
        var items = group.items.map(function(rawItem, itemIndex) {
          if (!rawItem) return null
          var info = rawItem.info || {}
          var label = rawItem.text || (group.isSequence ? (group.title || '组件') + (itemIndex + 1) : '未知')
          var icon = rawItem.icon || ''
          return {
            key: '2d_' + groupIndex + '_' + itemIndex + '_' + (info.type || 'unknown'),
            type: '2dComponent',
            typeName: info.type || '2DComponent',
            label: label,
            name: label,
            icon: icon,
            iconClass: icon && icon.indexOf('icon-') === 0 ? '' : (icon && icon.indexOf('svg-') === 0 ? '' : 'fas fa-vector-square'),
            color: 'cyan',
            is2DComponent: true,
            source2D: JSON.parse(JSON.stringify(info)),
            source2DMeta: {
              groupTitle: group.title || '',
              text: rawItem.text || '',
              icon: icon,
              isSequence: !!group.isSequence,
              index: itemIndex
            }
          }
        }).filter(Boolean)
        if (items.length > 0) {
          list.push({
            name: self.resolveGroupName(group),
            items: items
          })
        }
      })
      return list
    },
    filtered2DComponents() {
      if (!this.search2D) return this.normalized2DComponents
      const q = this.search2D.toLowerCase()
      return this.normalized2DComponents
        .map(cat => ({
          ...cat,
          items: cat.items.filter(i => {
            return (i.label || '').toLowerCase().includes(q) ||
              (i.typeName || '').toLowerCase().includes(q)
          })
        }))
        .filter(cat => cat.items.length > 0)
    },
    filteredComponents() {
      var source = this.searchComp ? this.componentLibrary : this.navComponentLibrary
      if (!this.searchComp) return source
      const q = this.searchComp.toLowerCase()
      return source
        .map(cat => ({
          ...cat,
          items: cat.items.filter(i => {
            var label = (i.label || '').toLowerCase()
            var typeName = (i.typeName || '').toLowerCase()
            return label.includes(q) || typeName.includes(q)
          })
        }))
        .filter(cat => cat.items.length > 0)
    }
  },
  methods: {
    selectComponentNav(key) {
      this.currentComponentNav = key
      this.searchComp = ''
      this.$nextTick(() => {
        this.activeComponentGroupKeys = this.filteredComponents.length ? [this.filteredComponents[0].name] : []
      })
    },
    onComponentGroupChange(keys) {
      this.activeComponentGroupKeys = Array.isArray(keys) ? keys : (keys ? [keys] : [])
    },
    sideNavClick(e) {
      this.currentNav = e.key
    },
    validResultsList(item) {
      return item.slice(0, 100)
    },
    MoreItemList(item) {
      return item.slice(100)
    },
    ShowMoreItem(items) {
      this.ShowMoreItemArray = items
      this.ShowMoreModel = true
      this.itemMore = true
    },
    beforeUpload() {
      this.$message.loading({ content: this.$t('dataModel.opcuaModel.ImportNodeIDLoading'), duration: 0 })
    },
    localUpgradeCharge(info) {
      let _t = this
      this.messageShowLoad = true
      if (info.file.status === 'done') {
        let result = info.file.response
        this.$message.destroy()
        if (result.Code == 0) {
          this.$message.success(`${info.file.name} ` + this.$t("SystemUpgrade.AuthUploadLoading"))
          _t.GetUserCustomPel()
        } else {
          this.$message.error(`${info.file.name} ` + this.$t("SystemUpgrade.DiyFileSaveError"))
        }
      } else if (info.file.status === 'error') {
        this.$message.destroy()
        this.$message.error(`${info.file.name} ` + this.$t("SystemUpgrade.DiyFileSaveError"))
      }
      this.messageShowLoad = false
    },
    GetUserCustomPel() {
      let _t = this
      GetUserCustomPel().then(function (res) {
        if (res.data.code == 0) {
          const ISMDiyComponentsList = []
          if ((res.data.list != null) && (res.data.list.length > 0)) {
            for (let i = 0; i < res.data.list.length; i++) {
              if ((res.data.list[i].FilePath != null) && (res.data.list[i].FilePath.length > 0)) {
                let toolBigScreenCustomPel = {
                  title: res.data.list[i].DirName,
                  icon: "icon-custom",
                  opened: false,
                  items: []
                }
                let imgList = res.data.list[i].FilePath.split(',')
                for (let j = 0; j < imgList.length; j++) {
                  let showPngVue = _t.createPngItem(imgList[j])
                  showPngVue.icon = imgList[j]
                  showPngVue.info.style.imageURL = imgList[j]
                  showPngVue.text = ''
                  toolBigScreenCustomPel.items.push(showPngVue)
                }
                ISMDiyComponentsList.push(toolBigScreenCustomPel)
              }
            }
          }
          store.commit('ISMDisPlayEditorTool/setDiyComponentsList', ISMDiyComponentsList)
        }
      })
    },
    initToolboxData() {
      const state = store.state.ISMDisPlayEditorTool
      if (!Array.isArray(state.toolBoxList) || state.toolBoxList.length === 0) {
        this.initToolBoxList(state)
      }
      if (!Array.isArray(state.PageCanVasList) || state.PageCanVasList.length === 0) {
        this.initPageCanvasList(state)
      }
      if (!Array.isArray(state.MesComponentsList) || state.MesComponentsList.length === 0) {
        this.initMesComponentsList(state)
      }
      if (!Array.isArray(state.DiyComponentsList)) {
        state.DiyComponentsList = []
      }
    },
    initToolBoxList(state) {
      const componentsStandard = require.context('@/pages/ISMDisPlay/ISMComponents/standard/', true, /\.vue$/)
      const componentsVideo = require.context('@/pages/ISMDisPlay/ISMComponents/video/', true, /\.vue$/)
      const componentsLogin = require.context('@/pages/ISMDisPlay/ISMComponents/login/', true, /\.vue$/)
      const componentsCanvas = require.context('@/pages/ISMDisPlay/ISMComponents/canvas/', true, /\.vue$/)
      const componentsCharts = require.context('@/pages/ISMDisPlay/ISMComponents/charts/', true, /\.vue$/)
      const componentsBigScreen = require.context('@/pages/ISMDisPlay/ISMComponents/bigScreen/', true, /\.vue$/)
      const componentsSvgArrows = require.context('@/pages/ISMDisPlay/ISMComponents/svg/arrows/', true, /\.vue$/)
      const componentsSvgElectric = require.context('@/pages/ISMDisPlay/ISMComponents/ComponentClassification/electric/', true, /\.vue$/)
      const componentsImages = require.context('@/pages/ISMDisPlay/ISMComponents/Images/', true, /\.vue$/)
      const deviceComponents = require.context('@/pages/ISMDisPlay/ISMComponents/device/', true, /\.vue$/)
      const mapComponents = require.context('@/pages/ISMDisPlay/ISMComponents/map/', true, /\.vue$/)
      const historyChartsComponents = require.context('@/pages/ISMDisPlay/ISMComponents/historyCharts/', true, /\.vue$/)
      const mesStandardComponents = require.context('@/pages/ISMDisPlay/ISMComponents/Mes/standard/', true, /\.vue$/)
      const componentPiping = require.context('../../../../../public/static/ISM/Conduit', true, /\.png$/)
      const componentElectricPngList = require.context('../../../../../public/static/ISM/systemImage/electric', true, /\.svg$/)
      const componentHVACList = require.context('../../../../../public/static/ISM/HVAC/', true, /\.png$/)
      const componentElectricMachineryList = require.context('../../../../../public/static/ISM/ElectricMachinery/', true, /\.png$/)
      const componentFanList = require.context('../../../../../public/static/ISM/Fan/', true, /\.png$/)
      const componentMercuryList = require.context('../../../../../public/static/ISM/Mercury/', true, /\.png$/)
      const componentBlenderList = require.context('../../../../../public/static/ISM/Blender/', true, /\.png$/)

      const toolStandardBoxList = { title: "displayConfig.ToolBox.Base.title", icon: "icon-standard-application", opened: false, items: [] }
      const toolVideoBoxList = { title: "displayConfig.ToolBox.Video.title", icon: "icon-standard-application", opened: false, items: [] }
      const toolLoginBoxList = { title: "displayConfig.ToolBox.login.title", icon: "icon-loginBox", opened: false, items: [] }
      const toolCanvasBoxList = { title: "displayConfig.ToolBox.Diagram.title", icon: "icon-standard-picture-empty", opened: false, items: [] }
      const toolMapBoxList = { title: "displayConfig.ToolBox.Map", icon: "icon-map", opened: false, items: [] }
      const toolHistoryChartsBoxList = { title: "displayConfig.ToolBox.HistoryCharts", icon: "icon-standard-chart-curve", opened: false, items: [] }
      const toolChartsBoxList = { title: "displayConfig.ToolBox.Charts.title", icon: "icon-standard-chart-bar", opened: false, items: [] }
      const toolSvgArrowsBoxList = { title: "displayConfig.ToolBox.Arrows", icon: "icon-standard-arrows", opened: false, isSequence: true, items: [] }
      const toolSvgElectricBoxList = { title: "displayConfig.ToolBox.Electric", icon: "icon-electric", opened: false, isSequence: true, vueCount: 0, items: [] }
      const toolBigScreenBoxContainerList = { title: "displayConfig.ToolBox.bigScreen.Container", icon: "icon-standard-big-screen", opened: false, items: [] }
      const toolDeviceContainerList = { title: "displayConfig.ToolBox.device.Container", icon: "icon-device", opened: false, items: [] }
      const toolImagesBoxList = { title: "configComponent.image.Text", icon: "icon-image", opened: false, items: [] }
      const toolMesStandardBoxList = { title: "displayConfig.ToolBox.MesStandard.title", icon: "icon-standard-chart-bar", opened: false, items: [] }
      const toolBoxHVACList = { title: "displayConfig.ToolBox.HVAC", icon: "icon-hvac", opened: false, items: [] }
      const toolBoxElectricMachineryList = { title: "displayConfig.ToolBox.ElectricMachinery", icon: "icon-electric-machinery", opened: false, items: [] }
      const toolBoxFanList = { title: "displayConfig.ToolBox.Fan", icon: "icon-fan", opened: false, items: [] }
      const toolBoxMercuryList = { title: "displayConfig.ToolBox.Mercury", icon: "icon-mercury", opened: false, items: [] }
      const toolBoxBlenderList = { title: "displayConfig.ToolBox.Blender", icon: "icon-blender", opened: false, items: [] }

      componentsStandard.keys().forEach(filePath => {
        const comp = componentsStandard(filePath)
        if (comp.default && comp.default.data) {
          const base = comp.default.data().base
          if (base && base.info) {
            base.info.type = comp.default.name
            toolStandardBoxList.items.push(base)
          }
        }
      })

      componentsVideo.keys().forEach(filePath => {
        const comp = componentsVideo(filePath)
        if (comp.default && comp.default.data) {
          const base = comp.default.data().base
          if (base && base.info) {
            base.info.type = comp.default.name
            toolVideoBoxList.items.push(base)
          }
        }
      })

      componentsLogin.keys().forEach(filePath => {
        const comp = componentsLogin(filePath)
        if (comp.default && comp.default.data) {
          const base = comp.default.data().base
          if (base && base.info) {
            base.info.type = comp.default.name
            toolLoginBoxList.items.push(base)
          }
        }
      })

      componentsCanvas.keys().forEach(filePath => {
        const comp = componentsCanvas(filePath)
        if (comp.default && comp.default.data) {
          const base = comp.default.data().base
          if (base && base.info) {
            base.info.type = comp.default.name
            toolCanvasBoxList.items.push(base)
          }
        }
      })

      componentsCharts.keys().forEach(filePath => {
        const comp = componentsCharts(filePath)
        if (comp.default && comp.default.data) {
          const base = comp.default.data().base
          if (base && base.info) {
            base.info.type = comp.default.name
            toolChartsBoxList.items.push(base)
          }
        }
      })

      componentsBigScreen.keys().forEach(filePath => {
        const comp = componentsBigScreen(filePath)
        if (comp.default && comp.default.data) {
          const base = comp.default.data().base
          if (base && base.info) {
            base.info.type = comp.default.name
            toolBigScreenBoxContainerList.items.push(base)
          }
        }
      })

      componentsSvgArrows.keys().forEach(filePath => {
        const comp = componentsSvgArrows(filePath)
        if (comp.default && comp.default.data) {
          const base = comp.default.data().base
          if (base && base.info) {
            base.info.type = comp.default.name
            toolSvgArrowsBoxList.items.push(base)
          }
        }
      })

      toolSvgElectricBoxList.vueCount = 0
      componentsSvgElectric.keys().forEach(filePath => {
        const comp = componentsSvgElectric(filePath)
        if (comp.default && comp.default.data) {
          const base = comp.default.data().base
          if (base && base.info) {
            base.info.type = comp.default.name
            toolSvgElectricBoxList.items.push(base)
            toolSvgElectricBoxList.vueCount++
          }
        }
      })

      deviceComponents.keys().forEach(filePath => {
        const comp = deviceComponents(filePath)
        if (comp.default && comp.default.data) {
          const base = comp.default.data().base
          if (base && base.info) {
            base.info.type = comp.default.name
            toolDeviceContainerList.items.push(base)
          }
        }
      })

      mapComponents.keys().forEach(filePath => {
        const comp = mapComponents(filePath)
        if (comp.default && comp.default.data) {
          const base = comp.default.data().base
          if (base && base.info) {
            base.info.type = comp.default.name
            toolMapBoxList.items.push(base)
          }
        }
      })

      historyChartsComponents.keys().forEach(filePath => {
        const comp = historyChartsComponents(filePath)
        if (comp.default && comp.default.data) {
          const base = comp.default.data().base
          if (base && base.info) {
            base.info.type = comp.default.name
            toolHistoryChartsBoxList.items.push(base)
          }
        }
      })

      componentsImages.keys().forEach(filePath => {
        const comp = componentsImages(filePath)
        if (comp.default && comp.default.data) {
          const base = comp.default.data().base
          if (base && base.info) {
            base.info.type = comp.default.name
            toolImagesBoxList.items.push(base)
          }
        }
      })

      mesStandardComponents.keys().forEach(filePath => {
        const comp = mesStandardComponents(filePath)
        if (comp.default && comp.default.data) {
          const base = comp.default.data().base
          if (base && base.info) {
            base.info.type = comp.default.name
            toolMesStandardBoxList.items.push(base)
          }
        }
      })

      toolSvgElectricBoxList.items.splice(toolSvgElectricBoxList.vueCount, toolSvgElectricBoxList.items.length - toolSvgElectricBoxList.vueCount)
      componentElectricPngList.keys().forEach(filePath => {
        const fileName = filePath.split('/').pop()
        const showPngVue = this.createPngItem('/static/ISM/systemImage/electric/' + fileName)
        showPngVue.text = ''
        showPngVue.icon = '/static/ISM/systemImage/electric/' + fileName
        showPngVue.info.style.imageURL = '/static/ISM/systemImage/electric/' + fileName
        toolSvgElectricBoxList.items.push(showPngVue)
      })

      const toolConduitBoxList = { title: Conduit.groupTitle, icon: Conduit.icon, opened: Conduit.opened, isSequence: Conduit.isSequence, items: [] }
      componentPiping.keys().forEach(filePath => {
        const fileName = filePath.split('/').pop()
        const showPngVue = this.createPngItem('/static/ISM/Conduit/' + fileName)
        toolConduitBoxList.items.push(showPngVue)
      })

      componentHVACList.keys().forEach(filePath => {
        const fileName = filePath.split('/').pop()
        const showPngVue = this.createPngItem('/static/ISM/HVAC/' + fileName)
        showPngVue.icon = '/static/ISM/HVAC/' + fileName
        showPngVue.info.style.imageURL = '/static/ISM/HVAC/' + fileName
        showPngVue.text = ''
        toolBoxHVACList.items.push(showPngVue)
      })

      componentElectricMachineryList.keys().forEach(filePath => {
        const fileName = filePath.split('/').pop()
        const showPngVue = this.createPngItem('/static/ISM/ElectricMachinery/' + fileName)
        showPngVue.icon = '/static/ISM/ElectricMachinery/' + fileName
        showPngVue.info.style.imageURL = '/static/ISM/ElectricMachinery/' + fileName
        showPngVue.text = ''
        toolBoxElectricMachineryList.items.push(showPngVue)
      })

      componentFanList.keys().forEach(filePath => {
        const fileName = filePath.split('/').pop()
        const showPngVue = this.createPngItem('/static/ISM/Fan/' + fileName)
        showPngVue.icon = '/static/ISM/Fan/' + fileName
        showPngVue.info.style.imageURL = '/static/ISM/Fan/' + fileName
        showPngVue.text = ''
        toolBoxFanList.items.push(showPngVue)
      })

      componentMercuryList.keys().forEach(filePath => {
        const fileName = filePath.split('/').pop()
        const showPngVue = this.createPngItem('/static/ISM/Mercury/' + fileName)
        showPngVue.icon = '/static/ISM/Mercury/' + fileName
        showPngVue.info.style.imageURL = '/static/ISM/Mercury/' + fileName
        showPngVue.text = ''
        toolBoxMercuryList.items.push(showPngVue)
      })

      componentBlenderList.keys().forEach(filePath => {
        const fileName = filePath.split('/').pop()
        const showPngVue = this.createPngItem('/static/ISM/Blender/' + fileName)
        showPngVue.icon = '/static/ISM/Blender/' + fileName
        showPngVue.info.style.imageURL = '/static/ISM/Blender/' + fileName
        showPngVue.text = ''
        toolBoxBlenderList.items.push(showPngVue)
      })

      const toolBoxList = []
      toolBoxList.push(toolStandardBoxList)
      toolBoxList.push(toolVideoBoxList)
      toolBoxList.push(toolLoginBoxList)
      toolBoxList.push(toolDeviceContainerList)
      toolBoxList.push(toolCanvasBoxList)
      toolBoxList.push(toolChartsBoxList)
      toolBoxList.push(toolHistoryChartsBoxList)
      toolBoxList.push(toolMapBoxList)
      toolBoxList.push(toolSvgArrowsBoxList)
      toolBoxList.push(toolSvgElectricBoxList)
      toolBoxList.push(toolConduitBoxList)
      toolBoxList.push(toolBoxHVACList)
      toolBoxList.push(toolBoxElectricMachineryList)
      toolBoxList.push(toolBoxFanList)
      toolBoxList.push(toolBoxMercuryList)
      toolBoxList.push(toolBoxBlenderList)

      state.toolBoxList = toolBoxList
    },
    createPngItem(iconUrl) {
      return {
        text: 'configComponent.image.Text',
        icon: iconUrl,
        isFontIcon: false,
        info: {
          type: 'ism-view-png-image',
          action: [],
          dataBind: [],
          animate: {
            selected: [],
            condition: { deviceSN: '', selectVideoType: 0, isBandDevice: false, bandType: 1, dataID: '', dataName: '', operator: '', OperatorValue: '', OperatorMaxValue: '' },
            isExpression: false,
            animateList: [
              { id: 'Forbidden', name: 'component.public.Forbidden' },
              { id: 'blink', name: 'component.public.animateBlink' },
              { id: 'Zoom', name: 'component.public.Zoom' },
              { id: 'animateSpin', name: 'component.public.animateSpin' }
            ],
            animateElement: [
              { id: 'blink', elementList: [{ name: 'component.public.animateSpeed', type: 7, value: 1, min: 0.1, key: 'blinkSpeed' }] },
              { id: 'millcolorGrad', elementList: [
                { name: 'component.public.startColor', type: 2, value: '#74f808', key: 'startColor' },
                { name: 'component.public.stopColor', type: 2, value: '#f30b0b', key: 'stopColor' },
                { name: 'component.public.animateSpeed', type: 7, value: 1, min: 0.1, key: 'animateSpeed' }
              ]},
              { id: 'animateSpin', elementList: [
                { name: 'component.public.animateSpinSpeed', type: 7, value: 1, min: 0.1, key: 'spinSpeed' },
                { name: 'configComponent.bigScreen.border.border89Direction', type: 6, value: 0, enumList: [{ value: 0, option: 'configComponent.bigScreen.border.border89DirectionForward' }, { value: 1, option: 'configComponent.bigScreen.border.border89DirectionNegative' }], min: 1, key: 'spinDirection' }
              ]}
            ]
          },
          style: {
            position: { x: 0, y: 0, w: 32, h: 32 },
            visible: 1,
            backColor: 'transparent',
            zIndex: -1,
            transform: 0,
            imageURL: iconUrl,
            diy: []
          }
        }
      }
    },
    initPageCanvasList(state) {
      const PageCanVasList = []

      const toolPageNavigationList = {
        title: PageNavigation.groupTitle,
        icon: PageNavigation.icon,
        opened: PageNavigation.opened,
        isSequence: PageNavigation.isSequence,
        items: []
      }
      try {
        const componentPageNavigationList = require.context('../../../../../public/static/ISM/page/navigation/', true, /\.png$/)
        componentPageNavigationList.keys().forEach(filePath => {
          const fileName = filePath.split('/').pop()
          const showPngVue = this.createImageItem(fileName, '/static/ISM/page/navigation/' + fileName)
          toolPageNavigationList.items.push(showPngVue)
        })
      } catch (e) {}
      if (toolPageNavigationList.items.length > 0) {
        PageCanVasList.push(toolPageNavigationList)
      }

      const toolBigScreenBoxContainerList = {
        title: 'displayConfig.ToolBox.bigScreen.Container',
        icon: 'icon-standard-big-screen',
        opened: false,
        items: []
      }
      try {
        const componentPageContainerList = require.context('../../../../../public/static/ISM/page/container/', true, /\.png$/)
        componentPageContainerList.keys().forEach(filePath => {
          const fileName = filePath.split('/').pop()
          const showPngVue = this.createImageItem(fileName, '/static/ISM/page/container/' + fileName, 300, 300)
          toolBigScreenBoxContainerList.items.push(showPngVue)
        })
      } catch (e) {}
      if (toolBigScreenBoxContainerList.items.length > 0) {
        PageCanVasList.push(toolBigScreenBoxContainerList)
      }

      const toolBigScreenBoxDecorateList = {
        title: 'displayConfig.ToolBox.bigScreen.Decorate',
        icon: 'icon-Decorate',
        opened: false,
        items: []
      }
      try {
        const componentPageDecorateList = require.context('../../../../../public/static/ISM/page/decorate/', true, /\.png$/)
        componentPageDecorateList.keys().forEach(filePath => {
          const fileName = filePath.split('/').pop()
          const showPngVue = this.createImageItem(fileName, '/static/ISM/page/decorate/' + fileName)
          toolBigScreenBoxDecorateList.items.push(showPngVue)
        })
      } catch (e) {}
      if (toolBigScreenBoxDecorateList.items.length > 0) {
        PageCanVasList.push(toolBigScreenBoxDecorateList)
      }

      const toolBigScreenBoxBackgroundBarList = {
        title: 'displayConfig.ToolBox.bigScreen.BackgroundBar',
        icon: 'icon-BackgroundBar',
        opened: false,
        items: []
      }
      try {
        const componentPageBackgroundBarList = require.context('../../../../../public/static/ISM/page/BackgroundBar/', true, /\.png$/)
        componentPageBackgroundBarList.keys().forEach(filePath => {
          const fileName = filePath.split('/').pop()
          const showPngVue = this.createImageItem(fileName, '/static/ISM/page/BackgroundBar/' + fileName)
          toolBigScreenBoxBackgroundBarList.items.push(showPngVue)
        })
      } catch (e) {}
      if (toolBigScreenBoxBackgroundBarList.items.length > 0) {
        PageCanVasList.push(toolBigScreenBoxBackgroundBarList)
      }

      const toolBigScreenBoxScreenedging1List = {
        title: 'displayConfig.ToolBox.bigScreen.Edging1',
        icon: 'icon-edding',
        opened: false,
        items: []
      }
      try {
        const componentPageScreenedging1List = require.context('../../../../../public/static/ISM/page/screenedging1/', true, /\.png$/)
        componentPageScreenedging1List.keys().forEach(filePath => {
          const fileName = filePath.split('/').pop()
          const showPngVue = this.createImageItem(fileName, '/static/ISM/page/screenedging1/' + fileName)
          toolBigScreenBoxScreenedging1List.items.push(showPngVue)
        })
      } catch (e) {}
      if (toolBigScreenBoxScreenedging1List.items.length > 0) {
        PageCanVasList.push(toolBigScreenBoxScreenedging1List)
      }

      const toolBigScreenBoxScreenedging2List = {
        title: 'displayConfig.ToolBox.bigScreen.Edging2',
        icon: 'icon-edding',
        opened: false,
        items: []
      }
      try {
        const componentPageScreenedging2List = require.context('../../../../../public/static/ISM/page/screenedging2/', true, /\.svg$/)
        componentPageScreenedging2List.keys().forEach(filePath => {
          const fileName = filePath.split('/').pop()
          const showPngVue = this.createImageItem(fileName, '/static/ISM/page/screenedging2/' + fileName)
          toolBigScreenBoxScreenedging2List.items.push(showPngVue)
        })
      } catch (e) {}
      if (toolBigScreenBoxScreenedging2List.items.length > 0) {
        PageCanVasList.push(toolBigScreenBoxScreenedging2List)
      }

      const toolBigScreenBoxScreenedging3List = {
        title: 'displayConfig.ToolBox.bigScreen.Edging3',
        icon: 'icon-edding',
        opened: false,
        items: []
      }
      try {
        const componentPageScreenedging3List = require.context('../../../../../public/static/ISM/page/screenedging3/', true, /\.png$/)
        componentPageScreenedging3List.keys().forEach(filePath => {
          const fileName = filePath.split('/').pop()
          const showPngVue = this.createImageItem(fileName, '/static/ISM/page/screenedging3/' + fileName)
          toolBigScreenBoxScreenedging3List.items.push(showPngVue)
        })
      } catch (e) {}
      if (toolBigScreenBoxScreenedging3List.items.length > 0) {
        PageCanVasList.push(toolBigScreenBoxScreenedging3List)
      }

      const toolBigScreenBox3DIconList = {
        title: 'displayConfig.ToolBox.bigScreen.3D_MODEL',
        icon: 'icon-3DMODEL',
        opened: false,
        items: []
      }
      try {
        const componentPage3DIconList = require.context('../../../../../public/static/ISM/page/3DIcon/', true, /\.svg$/)
        componentPage3DIconList.keys().forEach(filePath => {
          const fileName = filePath.split('/').pop()
          const showPngVue = this.createImageItem(fileName, '/static/ISM/page/3DIcon/' + fileName)
          toolBigScreenBox3DIconList.items.push(showPngVue)
        })
      } catch (e) {}
      if (toolBigScreenBox3DIconList.items.length > 0) {
        PageCanVasList.push(toolBigScreenBox3DIconList)
      }

      const toolBigScreenBoxTeIconList = {
        title: 'displayConfig.ToolBox.bigScreen.TeIcon',
        icon: 'icon-TeIcon',
        opened: false,
        items: []
      }
      try {
        const componentPageBigpicList = require.context('../../../../../public/static/ISM/page/bigpic/', true, /\.png$/)
        componentPageBigpicList.keys().forEach(filePath => {
          const fileName = filePath.split('/').pop()
          const showPngVue = this.createImageItem(fileName, '/static/ISM/page/bigpic/' + fileName)
          toolBigScreenBoxTeIconList.items.push(showPngVue)
        })
      } catch (e) {}
      if (toolBigScreenBoxTeIconList.items.length > 0) {
        PageCanVasList.push(toolBigScreenBoxTeIconList)
      }

      const toolBigScreenBoxllustrationList = {
        title: 'displayConfig.ToolBox.bigScreen.Illustration',
        icon: 'icon-Illustration',
        opened: false,
        items: []
      }
      try {
        const componentPageLlustrationList = require.context('../../../../../public/static/ISM/page/Illustration/', true, /\.png$/)
        componentPageLlustrationList.keys().forEach(filePath => {
          const fileName = filePath.split('/').pop()
          const showPngVue = this.createImageItem(fileName, '/static/ISM/page/Illustration/' + fileName)
          toolBigScreenBoxllustrationList.items.push(showPngVue)
        })
      } catch (e) {}
      if (toolBigScreenBoxllustrationList.items.length > 0) {
        PageCanVasList.push(toolBigScreenBoxllustrationList)
      }

      state.PageCanVasList = PageCanVasList
    },
    initMesComponentsList(state) {
      const MesComponentsList = []
      const MesStandardComponentsDir = require.context('@/pages/ISMDisPlay/ISMComponents/Mes/standard/', true, /\.vue$/)
      const MesStandardList = {
        title: 'displayConfig.ToolBox.MesStandard.title',
        icon: 'icon-standard-chart-bar',
        opened: false,
        items: []
      }
      MesStandardComponentsDir.keys().forEach(filePath => {
        const comp = MesStandardComponentsDir(filePath)
        if (comp.default && comp.default.data) {
          const base = comp.default.data().base
          if (base && base.info) {
            base.info.type = comp.default.name
            MesStandardList.items.push(base)
          }
        }
      })
      MesComponentsList.push(MesStandardList)
      state.MesComponentsList = MesComponentsList
    },
    normalizeGroups(groups) {
      var list = []
      var self = this
      ;(groups || []).forEach(function(group, groupIndex) {
        if (!group || !group.items || group.items.length === 0) return
        var items = group.items.map(function(rawItem, itemIndex) {
          if (!rawItem) return null
          var info = rawItem.info || {}
          var label = rawItem.text || (group.isSequence ? (group.title || '组件') + (itemIndex + 1) : '未知')
          var icon = rawItem.icon || ''
          return {
            key: '2d_' + groupIndex + '_' + itemIndex + '_' + (info.type || 'unknown'),
            type: '2dComponent',
            typeName: info.type || '2DComponent',
            label: label,
            name: label,
            icon: icon,
            iconClass: icon && icon.indexOf('icon-') === 0 ? '' : (icon && icon.indexOf('svg-') === 0 ? '' : 'fas fa-vector-square'),
            color: 'cyan',
            is2DComponent: true,
            source2D: JSON.parse(JSON.stringify(info)),
            source2DMeta: {
              groupTitle: group.title || '',
              text: rawItem.text || '',
              icon: icon,
              isSequence: !!group.isSequence,
              index: itemIndex
            }
          }
        }).filter(Boolean)
        if (items.length > 0) {
          list.push({
            name: self.resolveGroupName(group),
            items: items
          })
        }
      })
      return list
    },
    createImageItem(text, iconUrl, width = 32, height = 32) {
      return {
        text: "",
        icon: iconUrl,
        isFontIcon: false,
        info: {
          type: 'ism-view-png-image',
          action: [],
          dataBind: [],
          animate: {
            selected: [],
            condition: {
              deviceSN: '',
              selectVideoType: 0,
              isBandDevice: false,
              bandType: 1,
              dataID: '',
              dataName: '',
              operator: '',
              OperatorValue: '',
              OperatorMaxValue: ''
            },
            isExpression: false,
            animateList: [
              { id: 'Forbidden', name: 'component.public.Forbidden' },
              { id: 'blink', name: 'component.public.animateBlink' },
              { id: 'Zoom', name: 'component.public.Zoom' },
              { id: 'animateSpin', name: 'component.public.animateSpin' }
            ],
            animateElement: []
          },
          style: {
            position: { x: 0, y: 0, w: width, h: height },
            visible: 1,
            backColor: 'transparent',
            zIndex: -1,
            transform: 0,
            imageURL: iconUrl,
            diy: []
          }
        }
      }
    },
    resolveGroupName(group) {
      if (!group) return '2D'
      var title = group.title || '2D'
      return this.$te && this.$te(title) ? this.$t(title) : title
    },
    resolve2DLabel(item, group, index) {
      if (group && group.isSequence) {
        return this.resolveGroupName(group) + (index + 1)
      }
      var text = item && item.text ? item.text : ''
      if (text && this.$te && this.$te(text)) return this.$t(text)
      if (text) return text
      if (item && item.info && item.info.type) return item.info.type
      return '2D Component'
    },
    onDragStart(e, item) {
      const dragItem = item.info ? item.info : item
      this.$emit('drag-start', dragItem, e)
      if (e && e.dataTransfer) {
        e.dataTransfer.effectAllowed = 'copy'
        const itemType = item.type || (item.info && item.info.type) || 'ism-component'
        e.dataTransfer.setData('text/plain', itemType)
        e.dataTransfer.setData('application/x-ism3d-object', JSON.stringify(item))
      }
    },
    onModalDragStart(e, item) {
      const dragItem = item.info ? item.info : item
      this.$emit('drag-start', dragItem, e)
      if (e && e.dataTransfer) {
        e.dataTransfer.effectAllowed = 'copy'
        const itemType = item.type || (item.info && item.info.type) || 'ism-component'
        e.dataTransfer.setData('text/plain', itemType)
        e.dataTransfer.setData('application/x-ism3d-object', JSON.stringify(item))
      }
      this.itemMore = false
    }
  }
}
</script>

<style scoped>
.left-panel {
  width: 440px;
  min-width: 440px;
  background: #fafafa;
  border-right: 1px solid #e8e8e8;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  overflow: hidden;
}

/* Tab 切换样式 */
.panel-tabs {
  display: flex;
  border-bottom: 1px solid #e8e8e8;
  background: #fff;
  flex-shrink: 0;
  min-height: 40px;
}

.tab-item {
  flex: 1;
  text-align: center;
  padding: 9px 0;
  font-size: 12px;
  color: #666;
  cursor: pointer;
  border-bottom: 2px solid transparent;
  transition: all 0.15s;
  user-select: none;
}

.tab-item:hover {
  color: #13c2c2;
  background: #e6fffb;
}

.tab-item.active {
  color: #13c2c2;
  border-bottom-color: #13c2c2;
  font-weight: 600;
  background: #e6fffb;
}

.tab-item i {
  margin-right: 4px;
}

/* Tab 内容区 */
.tab-content {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

.component-tab-layout {
  flex-direction: row;
  overflow: hidden;
}

.component-side-nav {
  width: 46px;
  flex-shrink: 0;
  border-right: 1px solid #b5f5ec;
  background: #fff;
  padding-top: 6px;
}

.component-side-nav-item {
  position: relative;
  width: 45px;
  height: 44px;
  border: 0;
  border-left: 3px solid transparent;
  background: transparent;
  color: #13c2c2;
  cursor: pointer;
  outline: none;
}

.component-side-nav-item i {
  font-size: 18px;
}

.component-side-nav-item:hover,
.component-side-nav-item.active {
  background: #e6fffb;
  border-left-color: #13c2c2;
  color: #08979c;
}

.component-main-pane {
  min-width: 0;
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: #fff;
}

.panel-title {
  padding: 10px 14px;
  font-size: 12px;
  font-weight: 600;
  color: #666;
  text-transform: uppercase;
  letter-spacing: .8px;
  border-bottom: 1px solid #e8e8e8;
  background: #fff;
  flex-shrink: 0;
}

.component-search {
  padding: 8px 10px 6px;
  border-bottom: 1px solid #f0f0f0;
  flex-shrink: 0;
}

.component-library-collapse {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  background: #fff;
}

.component-library-collapse :deep(.ant-collapse-item) {
  border-bottom: 1px solid #f2f2f2;
}

.component-library-collapse :deep(.ant-collapse-header) {
  padding-top: 10px;
  padding-bottom: 10px;
  color: #222;
  font-size: 13px;
}

.component-library-group {
  max-height: none;
  overflow: visible;
}

.comp-img-icon {
  width: 32px;
  height: 32px;
  object-fit: contain;
  margin-bottom: 4px;
}

.toolbox-item.color-blue i { color: #13c2c2; }
.toolbox-item.color-green i { color: #52c41a; }
.toolbox-item.color-orange i { color: #fa8c16; }
.toolbox-item.color-purple i { color: #722ed1; }
.toolbox-item.color-pink i { color: #eb2f96; }
.toolbox-item.color-cyan i { color: #13c2c2; }
.toolbox-item.color-yellow i { color: #fadb14; }
.toolbox-item.color-red i { color: #ff4d4f; }
.toolbox-item.color-gray i { color: #8c8c8c; }
.toolbox-item.color-brown i { color: #a67c52; }

/* 2D组件Tab布局样式 */
.tab-content-2d-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.tab-content-2d-wrapper :deep(.ant-layout) {
  height: 100%;
}

.tab-content-2d-wrapper :deep(.ant-layout-sider) {
  flex-shrink: 0;
  overflow: hidden;
}

.tab-content-2d-wrapper :deep(.ant-layout-content) {
  overflow: hidden;
}

.toolbox-content-inner {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.toolbox-content-inner .component-search {
  padding: 6px;
  border-bottom: 1px solid #e8e8e8;
}

.toolbox-content-inner > .ant-collapse {
  flex: 1;
  overflow-y: auto;
}

/* 2D组件工具箱样式 - 与2D编辑器完全一致 */
.ism-toolbox {
  background-color: white;
  overflow-y: auto;
  overflow-x: hidden;
  max-height: 300px;
  .toolbox-group {
    display: flex;
    flex-wrap: wrap;
    padding: 5px;
    justify-content: flex-start;
    align-content: space-between;
    .toolbox-item {
      width: 70px;
      padding: 6px;
      color: #777;
      border: transparent solid 1px;
      &.base {
        width: 64px;
      }
      .toolbox-item-icon {
        text-align: center;
      }
      .toolbox-item-text {
        margin-top: 2px;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
        text-align: center;
        font-size: 12px;
      }
    }
    .toolbox-item:hover {
      border: #ccc solid 1px;
      background: #e6fffb;
      color: #13c2c2;
      border-radius: 6px;
      cursor: pointer;
    }
  }
}

/* 折叠面板样式覆盖 */
:deep(.ant-collapse-borderless > .ant-collapse-item > .ant-collapse-content > .ant-collapse-content-box) {
  padding-top: 0px;
}

:deep(.ant-collapse-content > .ant-collapse-content-box) {
  padding: 0px;
}

/* 滚动条样式 */
::-webkit-scrollbar {
  width : 5px;
  height: 9px;
}
::-webkit-scrollbar-thumb {
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
  box-shadow   : inset 0 0 5px rgba(0, 0, 0, 0.2);
  background   : #ededed;
  border-radius: 10px;
}

.scene-object-list {
  flex: 1;
  overflow-y: auto;
}

.scene-obj-item {
  display: flex;
  align-items: center;
  padding: 5px 14px;
  font-size: 12px;
  color: #666;
  cursor: pointer;
  gap: 7px;
  border-left: 2px solid transparent;
  transition: all .12s;
}

.scene-obj-item:hover {
  background: #f0f0f0;
}

.scene-obj-item.selected {
  background: #e6fffb;
  border-left-color: #13c2c2;
  color: #13c2c2;
}

.scene-obj-item i {
  font-size: 12px;
  color: #13c2c2;
}

.scene-obj-item .obj-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.scene-obj-item.is-background {
  opacity: 0.7;
}

.scene-obj-item.is-background .obj-name {
  color: #13c2c2;
  font-style: italic;
}

.scene-obj-item .obj-bg-tag {
  font-size: 10px;
  background: #13c2c2;
  color: #fff;
  padding: 1px 4px;
  border-radius: 3px;
  margin-left: 4px;
}

.scene-obj-item .obj-actions {
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity .12s;
}

.scene-obj-item:hover .obj-actions {
  opacity: 1;
}

.scene-obj-item .obj-action-btn {
  padding: 1px 4px;
  border-radius: 3px;
  color: #999;
  transition: all .12s;
}

.scene-obj-item .obj-action-btn:hover {
  color: #13c2c2;
  background: #e6fffb;
}

.scene-obj-item .obj-del {
  color: #ccc;
  font-size: 11px;
  padding: 1px 4px;
  border-radius: 3px;
  opacity: 0;
  transition: opacity .12s;
}

.scene-obj-item:hover .obj-del {
  opacity: 1;
  color: #ff4d4f;
}

.grid-settings-panel {
  padding: 8px 10px;
  border-top: 1px solid #e8e8e8;
  flex-shrink: 0;
}

.grid-settings-panel .prop-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.grid-settings-panel .prop-label {
  font-size: 11px;
  color: #666;
  min-width: 50px;
  flex-shrink: 0;
}
</style>
