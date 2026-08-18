<template>
  <div xmlns="http://www.w3.org/1999/xhtml"
       v-show="(detail.style.visible==1||isStart) && !isRetiredOverviewHint && !isOverviewAlarmBackdrop && !isRetiredOverviewKpi"
       @mousedown="onTextMouseDown"
       @click="onTextClick"
       :title="pagerInfoTitle || (isEnergyOverviewMode ? energyOverviewCoverageTitle : (labelRole === 'deviceBreadcrumb' ? displayText : null))"
       :class="chromeClass"
       :style="boxStyle">
    <span v-if="labelRole === 'section'" class="vst-section-bar" aria-hidden="true"></span>
    <span v-if="labelRole === 'floorGroup' && !hasFloorIconPrefix" class="vst-floor-icon" aria-hidden="true">📋</span>
    <span v-if="onlineRatioParts" class="vst-text vst-online-ratio">
      <b>{{ onlineRatioParts[0] }}</b><i>/</i><span>{{ onlineRatioParts[1] }}</span>
    </span>
    <span v-else-if="deviceBreadcrumbParts.length" class="vst-text vst-device-breadcrumb">
      <span v-for="(part, index) in deviceBreadcrumbParts" :key="`${part}-${index}`" class="vst-breadcrumb-segment">
        <i>›</i><b :class="{ 'is-current': index === deviceBreadcrumbParts.length - 1 }">{{ part }}</b>
      </span>
    </span>
    <span v-else-if="isPagerInfoRole && pagerJumpEditing" class="vst-text vst-page-jump" @click.stop>
      第
      <input
        ref="pagerJumpInput"
        class="vst-page-jump-input"
        type="number"
        min="1"
        :max="pagerLiveTotalPages"
        v-model="pagerJumpInput"
        @keydown.enter.prevent="confirmPagerJump"
        @keydown.esc.prevent="cancelPagerJump"
        @blur="confirmPagerJump"
      />
      /{{ pagerLiveTotalPages }} 页 · 共 {{ pagerLiveTotalCount }} {{ pagerLiveUnit }}
    </span>
    <span v-else-if="isPagerInfoRole" class="vst-text vst-page-info-live">{{ livePagerInfoText }}</span>
    <span v-else class="vst-text">{{ displayText }}</span>
  </div>
</template>

<script>

import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'
import store from '@/store'
import { subscribeEnergyOverviewStats } from '@/services/energyOverview'

export default {
  mixins: [ISMChildAutoMixin],
    name: 'view-svg-text',
    inject: ['getNode'],
    data() {
      return {
        Text:"",
        DivOpacity:1,
        animateType:[],
        startColor:"#74f808",
        stopColor:"#74f808",
        animateSpeed:0.5,
        animateSpinSpeed:0.5,
        spinDirection:0,
        blinkSpeed:0.5,
        isStart:false,
        italic:false,
        imageURL:"",
        energyOverviewStats:null,
        energyOverviewError:false,
        pagerJumpEditing: false,
        pagerJumpInput: '1',
        detail: { style: { diy: [], position: { w: 100, h: 40 }, visible: 1, text: '' }, animate: { selected: [], animateElement: [] }, action: [] },
        IsToolBox:false,
        editMode:true,
        base:{
          text: "configComponent.label.Text",
          "icon": "icon-icon_svg_wenben",
          "isFontIcon": true,
          "info": {
            "type": "view-svg-text",
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
                "w": 100,
                "h": 40
              },
              "visible":1,
              "backColor": "transparent",
              "foreColor": "#000000",
              fontWeight:400,
              "zIndex": -1,
              "transform": 0,
              text: "标签",
              textAlign: "center",
              fontSize: 30,
              fontFamily: "Arial",
              letterSpacing:0,
              italic:0,
              "diy":[

              ]
            }
          }
        }
      }
    },
    computed: {
      energyOverviewRole() {
        const configuredRole = String((this.detail && this.detail.energyOverviewRole) || '')
        if (configuredRole) return configuredRole
        const legacyOverviewRoles = {
          'ov-stat-power-val': 'activePower',
          'ov-stat-energy-val': 'todayEnergy',
        }
        return legacyOverviewRoles[String((this.detail && this.detail.name) || '')] || ''
      },
      isEnergyOverviewMode() {
        return this.energyOverviewRole === 'activePower' || this.energyOverviewRole === 'todayEnergy'
      },
      energyOverviewCoverageTitle() {
        const stats = this.energyOverviewStats
        if (!stats || stats.dataStatus === 'ok') return null
        return `部分数据：有效 ${stats.validDevices || 0}/${stats.eligibleDevices || 0}，缺失 ${stats.missingDevices || 0}，歧义 ${stats.ambiguousDevices || 0}`
      },
      isRetiredOverviewHint() {
        const text = String((this.detail && this.detail.style && this.detail.style.text) || '').trim()
        return text === '● 与设备管理树顶级区域一致 · 10kV母线 → 馈线模块 · 点击下钻进线→母线→馈线'
          || text === '与设备管理树一致 · 点击区域下钻'
      },
      displayText() {
        if (this.isEnergyOverviewMode) return this.energyOverviewDisplayValue()
        const text = String((this.detail && this.detail.style && this.detail.style.text) || '')
        const name = String((this.detail && this.detail.name) || '')
        if (/stat-alarm-val$/.test(name) && /^\d+$/.test(text.trim())) {
          return text.trim().padStart(3, '0')
        }
        if (text === '拓扑概览 · 大屏主页面轮询') return '拓扑概览'
        return text.trim() === '›' ? '❯' : text
      },
      onlineRatioParts() {
        const match = String(this.displayText || '').trim().match(/^(\d+)\s*\/\s*(\d+)$/)
        return match ? [match[1], match[2]] : null
      },
      deviceBreadcrumbParts() {
        if (this.labelRole !== 'deviceBreadcrumb') return []
        return String(this.displayText || '')
          .replace(/^\s*›\s*/, '')
          .split(/\s*›\s*/)
          .map(part => part.trim())
          .filter(Boolean)
      },
      hasClickAction() {
        const actions = (this.detail && this.detail.action) || []
        return actions.some(a => a && a.type === 'click' && a.action === 'link')
      },
      textAlign: function(){
        if(this.detail.style.textAlign == undefined) {
          return "center";
        } else {
          return this.detail.style.textAlign;
        }
      },
      // flex 容器需要 justify-content 才能水平对齐文本；text-align 在 flex 单行文本里无效。
      // 仅当显式声明 textAlign 时才生效，未声明时保持 flex-start(左)，避免影响存量左对齐文本。
      justifyContent: function(){
        var a = this.detail.style.textAlign;
        if(a === 'center') return 'center';
        if(a === 'right')  return 'flex-end';
        return 'flex-start';
      },
      lineHeight: function() {
      if(this.detail.style.lineHeight == undefined) {
        return this.detail.style.position.h;
      }
      return this.detail.style.lineHeight;
    },
      /** 运行态导航层级：设备详情页禁用悬浮浮窗 */
      isDeviceDetailPage() {
        const nav = this.$store && this.$store.state.ISMDisPlayEditorTool
          ? this.$store.state.ISMDisPlayEditorTool.navContext
          : null
        return !!(nav && nav.kind === 'device')
      },
      isOverviewAlarmBackdrop() {
        const name = String((this.detail && this.detail.name) || '')
        return /stat-alarm-(?:glow|fill|val|bg|accent|icon|lab)$/.test(name)
      },
      /** 20260803：取消总功率/总能耗/在线设备顶栏数值卡（兼容旧 JSON 未重生） */
      isRetiredOverviewKpi() {
        const name = String((this.detail && this.detail.name) || '')
        if (/stat-(?:power|energy|online)-(?:glow|fill|val|bg|accent|icon|lab)$/.test(name)) {
          return true
        }
        if (this.isEnergyOverviewMode) {
          return true
        }
        const text = String((this.detail && this.detail.style && this.detail.style.text) || '').trim()
        return text === '总功率' || text === '今日用电量' || text === '在线设备'
      },
      /** 按文本/尺寸/链接自动识别大屏标签角色（不改模板 DB 也能图形化） */
      labelRole() {
        const style = this.detail.style || {}
        const diy = style.diy || []
        const roleItem = diy.find(d => d && d.key === 'labelRole')
        let role = roleItem && roleItem.value ? String(roleItem.value) : ''
        if (!role) {
          const h = style.position && style.position.h ? style.position.h : 40
          const t = String(style.text || '').trim()
          // 「返回上一级」含「上一级」，需单独匹配；勿只写「返回上级」
          if (/^←\s*/.test(t) || /返回总图|返回上一级|返回上级|返回首页/.test(t)
            || /^(?:‹|◀)\s*返回/.test(t)) role = 'navBack'
          else if (/馈线模块|分区标题|模块区/.test(t) || (/^馈线/.test(t) && h <= 22)) role = 'section'
          else if (/设备组$/.test(t.replace(/^📋\s*/, '')) || t === '设备组') role = 'floorGroup'
          else if (t === '›' || t === '→') role = 'breadcrumbArrow'
          else if (/^(?:‹|←|◀)?\s*(上一页|上页|前一页)/.test(t)) role = 'pagePrev'
          else if (/^(?:›|→|▶)?\s*(下一页|下页|后一页)/.test(t)) role = 'pageNext'
          else if (/第\s*\d+\s*\/\s*\d+\s*页/.test(t)) role = 'pageInfo'
          else if (this.hasClickAction && h >= 26 && /^←/.test(t)) role = 'navBack'
        }
        // 翻页角色按 Vuex 实时校正，避免仅改卡片页而 chrome 文案/可点状态卡住
        return this.resolveLivePagerRole(role)
      },
      isPagerInfoRole() {
        const r = this.labelRole
        return r === 'pageInfo' || r === 'detailPageInfo'
      },
      navContextLive() {
        const navStore = this.$store || store
        return navStore && navStore.state.ISMDisPlayEditorTool
          ? navStore.state.ISMDisPlayEditorTool.navContext
          : null
      },
      pagerLiveCurrent() {
        const nav = this.navContextLive
        if (!nav) return 1
        if (this.isDeviceListPagerNav(nav)) return (nav.pageIndex || 0) + 1
        if (this.isDatapointPagerNav(nav)) return (nav.datapointPageIndex || 0) + 1
        return 1
      },
      pagerLiveTotalPages() {
        const nav = this.navContextLive
        if (!nav) return 1
        if (this.isDeviceListPagerNav(nav)) {
          const size = Math.max(1, Number(nav.pageSize) || 49)
          const n = Number(nav.totalDevices) || 0
          // 以总数/页大小为准，避免陈旧 totalPages 把下一页提前锁死
          if (n > 0) return Math.max(1, Math.ceil(n / size))
          return Math.max(1, Number(nav.totalPages) || 1)
        }
        if (this.isDatapointPagerNav(nav)) {
          const size = Math.max(1, Number(nav.datapointPageSize) || 20)
          const n = Number(nav.totalDatapoints) || 0
          // 与 ViewRealTable.pagerTotalPages 一致：优先 ceil(total/pageSize)
          if (n > 0) return Math.max(1, Math.ceil(n / size))
          return Math.max(1, Number(nav.datapointTotalPages) || 1)
        }
        return 1
      },
      pagerLiveTotalCount() {
        const nav = this.navContextLive
        if (!nav) return 0
        if (this.isDeviceListPagerNav(nav)) return Number(nav.totalDevices) || 0
        if (this.isDatapointPagerNav(nav)) return Number(nav.totalDatapoints) || 0
        return 0
      },
      pagerLiveUnit() {
        const nav = this.navContextLive
        if (nav && this.isDatapointPagerNav(nav)) return '个测点'
        return '台'
      },
      livePagerInfoText() {
        const nav = this.navContextLive
        if (!nav) {
          return String((this.detail && this.detail.style && this.detail.style.text) || '')
        }
        if (this.isDeviceListPagerNav(nav)) {
          return `第 ${this.pagerLiveCurrent}/${this.pagerLiveTotalPages} 页 · 共 ${this.pagerLiveTotalCount} 台`
        }
        if (this.isDatapointPagerNav(nav)) {
          return `第 ${this.pagerLiveCurrent}/${this.pagerLiveTotalPages} 页 · 共 ${this.pagerLiveTotalCount} 个测点`
        }
        return String((this.detail && this.detail.style && this.detail.style.text) || '')
      },
      pagerInfoTitle() {
        if (!this.isPagerInfoRole || this.IsToolBox) return null
        return '点击输入页码跳转'
      },
      chromeClass() {
        const r = this.labelRole
        return r ? `vst-chrome vst-chrome--${r}` : ''
      },
      hasFloorIconPrefix() {
        const t = String((this.detail.style && this.detail.style.text) || '').trim()
        return /^📋/.test(t)
      },
      boxStyle() {
        const style = this.detail.style || {}
        const pos = style.position || {}
        const w = pos.w || 100
        const h = pos.h || 40
        const role = this.labelRole
        const base = {
          width: w + 'px',
          height: h + 'px',
          cursor: this.hasClickAction ? 'pointer' : 'default',
          overflow: 'hidden',
          display: 'flex',
          alignItems: 'center',
          justifyContent: this.justifyContent,
          lineHeight: 1.35,
          boxSizing: 'border-box',
          paddingTop: '2px',
          fontSize: (style.fontSize || 14) + 'px',
          fontFamily: style.fontFamily || 'Microsoft YaHei, PingFang SC, sans-serif',
          fontWeight: style.fontWeight || 400,
          color: style.foreColor || '#e2e8f0',
          textAlign: this.textAlign,
          whiteSpace: 'nowrap',
          textOverflow: 'ellipsis',
          position: 'relative',
        }
        if (role === 'globalOverview') {
          Object.assign(base, {
            color: '#8fb8cc',
            fontWeight: 600,
            letterSpacing: '0.3px',
            cursor: 'pointer',
          })
        } else if (role === 'navBack' || role === 'deviceListBack') {
          Object.assign(base, {
            background: 'linear-gradient(135deg, rgba(8, 42, 68, 0.92) 0%, rgba(12, 58, 88, 0.88) 100%)',
            border: '1px solid rgba(0, 229, 255, 0.45)',
            borderRadius: '4px',
            boxShadow: '0 0 10px rgba(0, 229, 255, 0.12), inset 0 1px 0 rgba(255,255,255,0.06)',
            paddingLeft: '10px',
            paddingRight: '12px',
            color: '#7ee8ff',
            fontWeight: 600,
            letterSpacing: '0.5px',
            cursor: 'pointer',
          })
        } else if (role === 'deviceBreadcrumb') {
          Object.assign(base, {
            // 标题左对齐后，面包屑从中后段起排，中间留空隙
            paddingLeft: '1020px',
            paddingRight: '290px',
            color: style.foreColor || '#8fb8cc',
            fontSize: Math.max(11, (style.fontSize || 12)) + 'px',
            fontWeight: 500,
            letterSpacing: '0.25px',
            justifyContent: 'flex-start',
            background: 'linear-gradient(90deg, transparent 52%, rgba(7, 48, 66, 0.22) 62%, transparent 90%)',
          })
        } else if (role === 'deviceInfoName') {
          Object.assign(base, {
            paddingLeft: '10px',
            color: '#dffaff',
            fontWeight: 600,
            letterSpacing: '0.35px',
            background: 'linear-gradient(90deg, rgba(8, 54, 75, 0.42), transparent)',
            borderLeft: '2px solid rgba(53, 225, 255, 0.7)',
          })
        } else if (role === 'section') {
          Object.assign(base, {
            background: 'linear-gradient(90deg, rgba(0, 229, 255, 0.18) 0%, rgba(0, 229, 255, 0.04) 72%, transparent 100%)',
            borderBottom: '1px solid rgba(0, 229, 255, 0.35)',
            paddingLeft: '8px',
            paddingRight: '6px',
            color: '#00e5ff',
            fontWeight: 600,
            fontSize: Math.max(12, (style.fontSize || 14)) + 'px',
            letterSpacing: '1px',
          })
        } else if (role === 'floorGroup') {
          Object.assign(base, {
            background: 'linear-gradient(90deg, rgba(14, 52, 78, 0.85) 0%, rgba(8, 30, 48, 0.6) 100%)',
            border: '1px solid rgba(0, 180, 220, 0.3)',
            borderLeft: '3px solid #00c8e8',
            borderRadius: '2px',
            paddingLeft: '6px',
            paddingRight: '8px',
            color: '#b8e4f8',
            fontWeight: 600,
          })
        } else if (role === 'breadcrumbArrow') {
          Object.assign(base, {
            color: '#38d9f5',
            fontSize: Math.max(16, (style.fontSize || 14)) + 'px',
            fontWeight: 700,
            justifyContent: 'center',
            textShadow: '0 0 8px rgba(0, 229, 255, 0.45)',
          })
        } else if (role === 'pagePrev' || role === 'pageNext' || role === 'pagePrevDisabled' || role === 'pageNextDisabled'
          || role === 'detailPagePrev' || role === 'detailPageNext' || role === 'detailPagePrevDisabled' || role === 'detailPageNextDisabled') {
          const disabled = role === 'pagePrevDisabled' || role === 'pageNextDisabled'
            || role === 'detailPagePrevDisabled' || role === 'detailPageNextDisabled'
          Object.assign(base, {
            background: disabled ? 'rgba(8, 30, 48, 0.4)' : 'rgba(8, 42, 68, 0.85)',
            border: `1px solid ${disabled ? 'rgba(60, 90, 120, 0.35)' : 'rgba(0, 229, 255, 0.4)'}`,
            borderRadius: '4px',
            paddingLeft: '10px',
            paddingRight: '10px',
            color: disabled ? '#5f7799' : '#7ee8ff',
            fontWeight: 600,
            cursor: disabled ? 'default' : 'pointer',
          })
        } else if (role === 'pageInfo' || role === 'detailPageInfo') {
          Object.assign(base, {
            color: '#9fb6d6',
            fontSize: Math.max(11, (style.fontSize || 12)) + 'px',
            justifyContent: 'center',
            cursor: this.IsToolBox ? 'default' : 'pointer',
          })
        }
        return base
      },
    },
    watch: {
      detail: {
        handler(newVal, oldVal) {
          this.initComponents(newVal);
          this.syncEnergyOverviewSubscription()
        },
        deep: true
      }
    },
    methods: {
      isDeviceListPagerNav(nav) {
        return !!(nav && nav.deviceListMode)
      },
      isDatapointPagerNav(nav) {
        return !!(nav && (nav.signalMode || nav.routeMode === 'signal'
          || nav.kind === 'device' || nav.allDatapoints || nav.datapoints))
      },
      resolveLivePagerRole(role) {
        const r = String(role || '')
        if (!r) return ''
        const nav = this.navContextLive
        if (!nav || this.IsToolBox) return r
        const cur = this.isDeviceListPagerNav(nav)
          ? (nav.pageIndex || 0)
          : (this.isDatapointPagerNav(nav) ? (nav.datapointPageIndex || 0) : 0)
        const totalPages = this.pagerLiveTotalPages
        if (r === 'pagePrev' || r === 'pagePrevDisabled') {
          return this.isDeviceListPagerNav(nav) && cur > 0 ? 'pagePrev' : 'pagePrevDisabled'
        }
        if (r === 'pageNext' || r === 'pageNextDisabled') {
          return this.isDeviceListPagerNav(nav) && cur < totalPages - 1 ? 'pageNext' : 'pageNextDisabled'
        }
        if (r === 'detailPagePrev' || r === 'detailPagePrevDisabled') {
          return this.isDatapointPagerNav(nav) && cur > 0 ? 'detailPagePrev' : 'detailPagePrevDisabled'
        }
        if (r === 'detailPageNext' || r === 'detailPageNextDisabled') {
          return this.isDatapointPagerNav(nav) && cur < totalPages - 1 ? 'detailPageNext' : 'detailPageNextDisabled'
        }
        if (r === 'pageInfo' && this.isDatapointPagerNav(nav) && !this.isDeviceListPagerNav(nav)) {
          return 'detailPageInfo'
        }
        return r
      },
      startPagerJump() {
        if (this.IsToolBox) return
        if (this.pagerLiveTotalPages <= 0) return
        this.pagerJumpEditing = true
        this.pagerJumpInput = String(this.pagerLiveCurrent)
        this.$nextTick(() => {
          const el = this.$refs.pagerJumpInput
          if (el && typeof el.focus === 'function') {
            el.focus()
            if (typeof el.select === 'function') el.select()
          }
        })
      },
      cancelPagerJump() {
        this.pagerJumpEditing = false
      },
      confirmPagerJump() {
        if (!this.pagerJumpEditing) return
        this.pagerJumpEditing = false
        let page = parseInt(String(this.pagerJumpInput || '').trim(), 10)
        if (!Number.isFinite(page)) return
        const total = Math.max(1, this.pagerLiveTotalPages)
        page = Math.max(1, Math.min(total, page))
        const nav = this.navContextLive
        if (!nav) return
        if (this.isDeviceListPagerNav(nav)) {
          const next = page - 1
          if (next === (nav.pageIndex || 0)) return
          this.$EventBus.$emit('NavPageChange', { pageIndex: next })
          return
        }
        if (this.isDatapointPagerNav(nav)) {
          const next = page - 1
          if (next === (nav.datapointPageIndex || 0)) return
          this.$EventBus.$emit('NavPageChange', { datapointPageIndex: next })
        }
      },
      energyOverviewDisplayValue() {
        const stats = this.energyOverviewStats
        if (!stats || stats.configured === false || this.energyOverviewError) return '--'
        const raw = stats.current && stats.current[this.energyOverviewRole]
        if (raw == null || raw === '') return '--'
        const value = raw && typeof raw === 'object' ? raw.value : raw
        const numeric = Number(value)
        const display = Number.isFinite(numeric)
          ? numeric.toFixed(2).replace(/\.?0+$/, '')
          : String(value)
        const defaultUnit = this.energyOverviewRole === 'activePower' ? 'kW' : 'kWh'
        const unit = (raw && typeof raw === 'object' && raw.unit)
          || (stats.units && stats.units[this.energyOverviewRole])
          || defaultUnit
        return unit ? `${display} ${unit}` : display
      },
      syncEnergyOverviewSubscription() {
        if (!this.isEnergyOverviewMode || this.IsToolBox) {
          if (this._energyOverviewUnsubscribe) {
            this._energyOverviewUnsubscribe()
            this._energyOverviewUnsubscribe = null
          }
          return
        }
        if (this._energyOverviewUnsubscribe) return
        this._energyOverviewUnsubscribe = subscribeEnergyOverviewStats((stats, error) => {
          this.energyOverviewStats = stats
          this.energyOverviewError = !!error
        })
      },
      // hover 预览：仅在运行态（非编辑/工具箱）且本 cell 绑定了设备时触发，
      // 不影响既有 onTextClick / GoPage 钻探逻辑。
      resolveDeviceBinding() {
        const d = this.detail || {}
        const actives = d.active || []
        for (let i = 0; i < actives.length; i++) {
          const c = actives[i] && actives[i].condition
          if (c && c.deviceSN) {
            return { uuid: c.deviceSN, name: c.DeviceName || (d.style && d.style.text) || '' }
          }
        }
        const ac = d.animate && d.animate.condition
        if (ac && ac.deviceSN) {
          return { uuid: ac.deviceSN, name: ac.DeviceName || (d.style && d.style.text) || '' }
        }
        return null
      },
      onTextMouseDown() {
        const role = this.labelRole
        if (!/^(?:(?:detailPage|page)(?:Prev|Next)(?:Disabled)?|deviceListBack|navBack|globalOverview)$/.test(role)) return
        this._pagerHandledByMouseDown = true
        this.onTextClick({ type: 'mousedown' })
      },
      onTextClick(event) {
        if (event && event.type === 'click' && this._pagerHandledByMouseDown) {
          this._pagerHandledByMouseDown = false
          return
        }
        const role = this.labelRole
        const isRuntimeControl = /^(?:(?:detailPage|page)(?:Prev|Next)(?:Disabled)?|pageInfo|detailPageInfo|deviceListBack|navBack|globalOverview)$/.test(role)
        // X6 运行态动态注入的分页文本偶尔仍携带 editMode=true。
        // 显式分页角色以 navContext 是否存在作为运行态门禁，不能在角色识别前被误拦截。
        if (this.IsToolBox || (this.editMode && !isRuntimeControl)) {
          return
        }
        if (role === 'breadcrumbArrow') return
        if (/Disabled$/.test(role)) return
        if (role === 'pageInfo' || role === 'detailPageInfo') {
          this.startPagerJump()
          return
        }
        // deviceListBack / navBack 同源：回设备列表（无列表上下文时由 RunTree 回首页）
        if (role === 'deviceListBack' || role === 'navBack') {
          this.$EventBus.$emit('ReturnToDeviceList')
          return
        }
        if (role === 'globalOverview') {
          const diy = (this.detail && this.detail.style && this.detail.style.diy) || []
          const nav = this.navContextLive
          const homePageUuid = String(
            (diy.find(d => d && d.key === 'homePageUuid') || {}).value
            || (nav && nav.homePageUuid)
            || '',
          )
          if (!homePageUuid) {
            // diy/nav 均无首页 uuid 时仍发事件，由 ISMRender/RunTree 用 modelId 兜底
            this.$EventBus.$emit('ReturnToGlobalOverview')
            return
          }
          this.$EventBus.$emit('GoPage', {
            IsPopUp: false,
            autoClose: false,
            linkType: 'Inside',
            ModelId: homePageUuid,
            PageUuid: homePageUuid,
            navContext: null,
          })
          return
        }
        if (role === 'pagePrev' || role === 'pageNext') {
          // X6 vue-shape 是独立 Vue 实例，未必注入 this.$store。
          // 使用共享 store 兜底，保证模板翻页控件与表格分页行为一致。
          const nav = this.navContextLive
          if (!nav || !nav.deviceListMode) return
          const cur = nav.pageIndex || 0
          const total = this.pagerLiveTotalPages
          const next = role === 'pagePrev' ? cur - 1 : cur + 1
          if (next < 0 || next >= total) return
          this.$EventBus.$emit('NavPageChange', { pageIndex: next })
          return
        }
        if (role === 'detailPagePrev' || role === 'detailPageNext') {
          const nav = this.navContextLive
          // 信号层测点翻页（与底部分页条同一事件，±1）
          if (nav && (nav.signalMode || nav.routeMode === 'signal' || nav.allDatapoints || nav.kind === 'device')) {
            const cur = nav.datapointPageIndex || 0
            const total = this.pagerLiveTotalPages
            const next = role === 'detailPagePrev' ? cur - 1 : cur + 1
            if (next < 0 || next >= total) return
            this.$EventBus.$emit('NavPageChange', { datapointPageIndex: next })
            return
          }
          if (!nav || !nav.detailPointMode) return
          const cur = nav.detailPageIndex || 0
          const total = nav.detailTotalPages || 1
          const next = role === 'detailPagePrev' ? cur - 1 : cur + 1
          if (next < 0 || next >= total) return
          this.$EventBus.$emit('NavPageChange', { detailPageIndex: next })
          return
        }
        const actions = (this.detail && this.detail.action) || []
        const clickAction = actions.find(
          a => a && a.type === 'click' && a.action === 'link' && a.link
        )
        if (!clickAction || !clickAction.link) {
          return
        }
        const link = clickAction.link
        this.$EventBus.$emit('GoPage', {
          IsPopUp: link.isPopUp,
          autoClose: link.autoClose,
          linkType: link.linkType,
          ModelId: link.Inside && link.Inside.displayUUID,
          PageUuid: link.Inside && link.Inside.pageUUID,
          width: link.width,
          height: link.height,
          External: link.External,
          title: link.title,
          OpenExternalType: link.OpenExternalType,
          // 层级模板：槽位重映射内嵌的子节点上下文
          navContext: link.navContext || null
        })
      },
      initComponents(option){
        if(this.IsToolBox)
        {
          return
        }
        if (!option || !option.style) return
        this.DivOpacity = option.style.opacity
        const diy = option.style.diy || []
        let i=0
        for( i=0;i<diy.length;i++)
        {
          if(diy[i].key=="strokeWidth")
          {
            this.strokeWidth=diy[i].value
          }
          else if(diy[i].key=="strokeFill")
          {
            this.fill=diy[i].value
          }
          else if(diy[i].key=="strokeColor")
          {
            this.strokeColor=diy[i].value
          }
          else if(diy[i].key=="fillOpacity")
          {
            this.fillOpacity=diy[i].value
          }
          else if(diy[i].key=="strokeOpacity")
          {
            this.strokeOpacity=diy[i].value
          }
          else if(diy[i].key=="imageURL")
          {
            this.imageURL=diy[i].value
          }
        }
        i=0
        if (this.isEnergyOverviewMode) {
          this.animateType = []
          this.isStart = false
          return
        }
        this.animateType = (option.animate && option.animate.selected) || []
        if(option.animate && option.animate.isExpression)
        {
          this.isStart = false
        }
        else
        {
          this.isStart = true
        }
      }
    },
    mounted() {
      let _t = this
      this.$nextTick(function(){
       this.initComponents(this.detail);
        this.syncEnergyOverviewSubscription()
        if (!this.detail || !this.detail.identifier) return
        let activeEvent = this.detail.identifier+"activeEvent"
        let animateEvent = this.detail.identifier+"animateEvent"

        _t.$EventBus.$on(activeEvent, (data) => {
          if (_t.isEnergyOverviewMode) return
          if((_t.editMode)&&(!this.IsToolBox)){
            return
          }
          if (data && data.result !== undefined) {
            let value = data.result
            if (typeof value === 'boolean') value = value ? '1' : '0'
            _t.detail.style.text = String(value)
          }
        })
        _t.$EventBus.$on(animateEvent, (data) => {
          if (_t.isEnergyOverviewMode) return
          if((_t.editMode)&&(!this.IsToolBox)){
            return
          }
          _t.isStart = data
        })

        // 信号层翻页：同步顶部页码文案与上/下一页可点状态
        _t._onNavDatapointPageUpdate = (nav) => {
          if (!nav || _t.IsToolBox) return
          const role = _t.labelRole
          if (!role || String(role).indexOf('detailPage') !== 0) return
          const cur = Number(nav.datapointPageIndex) || 0
          const size = Math.max(1, Number(nav.datapointPageSize) || 20)
          const n = Number(nav.totalDatapoints) || 0
          // 优先 ceil(total/pageSize)，避免陈旧 datapointTotalPages（如 63）锁死下一页
          const totalPages = n > 0
            ? Math.max(1, Math.ceil(n / size))
            : Math.max(1, Number(nav.datapointTotalPages) || 1)
          if (!_t.detail || !_t.detail.style) return
          if (!Array.isArray(_t.detail.style.diy)) _t.detail.style.diy = []
          const setRole = (r) => {
            const item = _t.detail.style.diy.find(d => d && d.key === 'labelRole')
            if (item) item.value = r
            else _t.detail.style.diy.push({ name: 'labelRole', type: 9, value: r, key: 'labelRole' })
          }
          if (role === 'detailPageInfo' || role === 'pageInfo') {
            _t.detail.style.text = `第 ${cur + 1}/${totalPages} 页 · 共 ${n} 个测点`
            setRole('detailPageInfo')
          } else if (role === 'detailPagePrev' || role === 'detailPagePrevDisabled') {
            setRole(cur > 0 ? 'detailPagePrev' : 'detailPagePrevDisabled')
          } else if (role === 'detailPageNext' || role === 'detailPageNextDisabled') {
            setRole(cur < totalPages - 1 ? 'detailPageNext' : 'detailPageNextDisabled')
          }
        }
        _t.$EventBus.$on('NavDatapointPageUpdate', _t._onNavDatapointPageUpdate)
      });
    },
    beforeDestroy() {
      if (this._energyOverviewUnsubscribe) {
        this._energyOverviewUnsubscribe()
        this._energyOverviewUnsubscribe = null
      }
      if (this._onNavDatapointPageUpdate && this.$EventBus) {
        this.$EventBus.$off('NavDatapointPageUpdate', this._onNavDatapointPageUpdate)
        this._onNavDatapointPageUpdate = null
      }
    },
    created(){
    let _t = this
    try {
      this.GetNodeObj = this.getNode()
      if (!this.GetNodeObj) {
        return
      }
      const nodeData = this.GetNodeObj.getData() || {}
      this.detail = nodeData.detail || this.detail
      this.editMode = nodeData.editMode || false
      this.showDeviceUuid = nodeData.showDeviceUuid || ''
      this.IsToolBox = nodeData.IsToolBox || false
      this.GetNodeObj.on('change:data', ({ current }) => {
        if(current && current.detail) {
          _t.detail = current.detail
        }
      })
      this.GetNodeObj.on('change:size', ({ current }) => {
        if (_t.detail && _t.detail.style && _t.detail.style.position) {
          _t.detail.style.position.w = current.width
          _t.detail.style.position.h = current.height
        }
      });
      _t.$EventBus.$on('cell-editMode', (data) => {
        _t.editMode = data.edit
        _t.IsToolBox = data.toolbox
      })
      this.initComponents(this.detail);
    } catch(e) {
      console.error('[ViewSvgText] created error:', e)
    }
  }
}
</script>
<style scoped>
.vst-text {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1 1 auto;
  min-width: 0;
  pointer-events: none;
}
.vst-page-info-live,
.vst-page-jump {
  pointer-events: auto;
}
.vst-page-jump {
  display: inline-flex;
  align-items: center;
  gap: 2px;
  overflow: visible;
}
.vst-page-jump-input {
  width: 42px;
  height: 20px;
  margin: 0 2px;
  padding: 0 4px;
  border: 1px solid rgba(0, 229, 255, 0.45);
  border-radius: 3px;
  background: rgba(8, 42, 68, 0.92);
  color: #7ee8ff;
  font-size: 12px;
  font-weight: 600;
  text-align: center;
  outline: none;
  -moz-appearance: textfield;
}
.vst-page-jump-input::-webkit-outer-spin-button,
.vst-page-jump-input::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}
.vst-section-bar {
  flex: 0 0 3px;
  align-self: stretch;
  margin-right: 6px;
  background: linear-gradient(180deg, #00e5ff, #0088bb);
  border-radius: 1px;
  box-shadow: 0 0 6px rgba(0, 229, 255, 0.5);
}
.vst-floor-icon {
  flex: 0 0 auto;
  margin-right: 4px;
  font-size: 0.9em;
  opacity: 0.9;
}
.vst-online-ratio b {
  color: #36e6a0;
  font-weight: 700;
  text-shadow: 0 0 9px rgba(54, 230, 160, 0.38);
}
.vst-online-ratio i {
  margin: 0 0.18em;
  color: rgba(159, 182, 214, 0.62);
  font-style: normal;
  font-weight: 400;
}
.vst-online-ratio > span {
  color: #e8f1ff;
  font-weight: 700;
}
.vst-device-breadcrumb {
  display: flex;
  align-items: center;
}
.vst-breadcrumb-segment {
  display: inline-flex;
  align-items: center;
  min-width: 0;
}
.vst-breadcrumb-segment i {
  flex: none;
  margin: 0 8px;
  color: #38d9f5;
  font-style: normal;
  font-weight: 700;
  text-shadow: 0 0 7px rgba(0, 229, 255, 0.36);
}
.vst-breadcrumb-segment b {
  overflow: hidden;
  color: #8fb8cc;
  font-weight: 500;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.vst-breadcrumb-segment b.is-current {
  color: #dffaff;
  font-weight: 650;
  text-shadow: 0 0 9px rgba(53, 225, 255, 0.28);
}
.vst-chrome--navBack:hover {
  border-color: rgba(0, 229, 255, 0.75);
  box-shadow: 0 0 14px rgba(0, 229, 255, 0.28), inset 0 1px 0 rgba(255,255,255,0.1);
  color: #aaf4ff;
}
.vst-chrome--navBack::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  width: 4px;
  height: 100%;
  background: linear-gradient(180deg, #00e5ff, transparent);
  border-radius: 4px 0 0 4px;
  opacity: 0.7;
}
</style>
<style>
.svg-el {
  transform-origin: center center;
}
</style>
