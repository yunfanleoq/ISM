<template>
  <div xmlns="http://www.w3.org/1999/xhtml"
       v-show="detail.style.visible==1||isStart"
       @click="onTextClick"
       :class="chromeClass"
       :style="boxStyle">
    <span v-if="labelRole === 'section'" class="vst-section-bar" aria-hidden="true"></span>
    <span v-if="labelRole === 'floorGroup' && !hasFloorIconPrefix" class="vst-floor-icon" aria-hidden="true">📋</span>
    <span class="vst-text">{{ detail.style.text }}</span>
  </div>
</template>

<script>

import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin'

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
      /** 按文本/尺寸/链接自动识别大屏标签角色（不改模板 DB 也能图形化） */
      labelRole() {
        const style = this.detail.style || {}
        const diy = style.diy || []
        const roleItem = diy.find(d => d && d.key === 'labelRole')
        if (roleItem && roleItem.value) return String(roleItem.value)
        const h = style.position && style.position.h ? style.position.h : 40
        const t = String(style.text || '').trim()
        if (/^←\s*/.test(t) || /返回总图|返回上级|返回首页/.test(t)) return 'navBack'
        if (/馈线模块|分区标题|模块区/.test(t) || (/^馈线/.test(t) && h <= 22)) return 'section'
        if (/设备组$/.test(t.replace(/^📋\s*/, '')) || t === '设备组') return 'floorGroup'
        if (/^(‹|←|◀|上一页|上页)/.test(t)) return 'pagePrev'
        if (/^(›|→|▶|下一页|下页)/.test(t)) return 'pageNext'
        if (/第\s*\d+\s*\/\s*\d+\s*页/.test(t)) return 'pageInfo'
        if (this.hasClickAction && h >= 26 && /^←/.test(t)) return 'navBack'
        return ''
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
        if (role === 'navBack') {
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
          })
        }
        return base
      },
    },
    watch: {
      detail: {
        handler(newVal, oldVal) {
          this.initComponents(newVal);
        },
        deep: true
      }
    },
    methods: {
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
      onTextClick() {
        if (this.IsToolBox || this.editMode) {
          return
        }
        const role = this.labelRole
        if (role === 'pagePrev' || role === 'pageNext') {
          const nav = this.$store && this.$store.state.ISMDisPlayEditorTool
            ? this.$store.state.ISMDisPlayEditorTool.navContext
            : null
          if (!nav || !nav.deviceListMode) return
          const cur = nav.pageIndex || 0
          const total = nav.totalPages || 1
          const next = role === 'pagePrev' ? cur - 1 : cur + 1
          if (next < 0 || next >= total) return
          this.$EventBus.$emit('NavPageChange', { pageIndex: next })
          return
        }
        if (role === 'detailPagePrev' || role === 'detailPageNext') {
          const nav = this.$store && this.$store.state.ISMDisPlayEditorTool
            ? this.$store.state.ISMDisPlayEditorTool.navContext
            : null
          // 信号层测点翻页（与底部分页条同一事件，±1）
          if (nav && (nav.signalMode || nav.routeMode === 'signal' || nav.allDatapoints || nav.kind === 'device')) {
            const cur = nav.datapointPageIndex || 0
            const total = Math.max(1, Number(nav.datapointTotalPages) || 1)
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
        if (!this.detail || !this.detail.identifier) return
        let activeEvent = this.detail.identifier+"activeEvent"
        let animateEvent = this.detail.identifier+"animateEvent"

        _t.$EventBus.$on(activeEvent, (data) => {
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
          if((_t.editMode)&&(!this.IsToolBox)){
            return
          }
          _t.isStart = data
        })

        // 信号层翻页：同步顶部页码文案与上/下一页可点状态
        _t._onNavDatapointPageUpdate = (nav) => {
          if (!nav || _t.editMode || _t.IsToolBox) return
          const role = _t.labelRole
          if (!role || String(role).indexOf('detailPage') !== 0) return
          const cur = nav.datapointPageIndex || 0
          const totalPages = Math.max(1, Number(nav.datapointTotalPages)
            || Math.ceil((Number(nav.totalDatapoints) || 0) / Math.max(1, Number(nav.datapointPageSize) || 20))
            || 1)
          const n = Number(nav.totalDatapoints) || 0
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
