<template>
  <section v-show="visible" class="org-overview">
    <header class="org-header">
      <div class="org-heading">
        <span class="org-heading-icon">⚡</span>
        <div>
          <b>组织层级总览</b>
          <span>POWER DISTRIBUTION ORGANIZATION TOPOLOGY</span>
        </div>
      </div>
    </header>
    <div v-if="loading" class="org-empty">正在加载组织架构…</div>
    <div v-else-if="!roots.length" class="org-empty">暂无组织架构数据</div>
    <div v-else ref="orgTree" class="org-tree">
      <div ref="orgTreeStage" class="org-tree-stage" :style="treeStageStyle">
        <org-node
          v-for="node in roots"
          :key="stableNodeKey(node)"
          :node="node"
          :depth="1"
          @open-device-list="openDeviceList"
          @open-device="openDevice"
          @open-virtual-cabinet="openVirtualCabinet"
          @layout-dirty="onLayoutDirty"
        />
      </div>
    </div>
  </section>
</template>

<script>
import { getMonitorTree } from '@/services/device'

function nodeType(node) {
  const t = (node && node.value && node.value.type)
  if (t === 'virtualCabinet') return 'virtualCabinet'
  return Number(t)
}

/** 拓扑网格最多铺这么多真设备；超出部分点「全部」列表看剩余 */
const MAX_DEVICES_IN_ORG_TREE = 40
/** 子节点超过此数改用网格，避免横向一排把整树缩到看不清 */
const DENSE_CHILD_THRESHOLD = 6
/** 等比缩放下限：宁可滚动也不再缩小字号 */
const MIN_TREE_SCALE = 0.78

const OrgNode = {
  name: 'OrgNode',
  props: {
    node: { type: Object, required: true },
    depth: { type: Number, default: 1 },
  },
  data() {
    return {
      openedByPointer: false,
      // 默认展开前 2 层，更深默认收起
      expanded: this.depth <= 2,
    }
  },
  computed: {
    value() { return this.node.value || {} },
    kind() {
      const t = nodeType(this.node)
      if (t === 'virtualCabinet') return 'virtualCabinet'
      if (t === 1) return 'device'
      return 'organization'
    },
    name() {
      return this.node.text || this.value.name || this.value.virtualCabinet || '未命名'
    },
    orgChildren() {
      return (this.node.children || []).filter(c => nodeType(c) === 0)
    },
    /** 同时有子组织与直属 type=1（如「配电室」下 3A1/3A2 + 25 个配电室*） */
    hasMixedOrgAndDevices() {
      return this.orgChildren.length > 0 && this.directDeviceCount > 0
    },
    rawDeviceChildren() {
      if (this.kind !== 'organization') return []
      // 有子组织时仍保留直属 type=1，否则只剩 2 个机房模块
      return (this.node.children || []).filter(c => nodeType(c) === 1)
    },
    paintedSiblingDevices() {
      // 仅混合层把 type=1 当「同级房间/配电室」卡片画出；末级不在树内嵌设备网格
      if (!this.hasMixedOrgAndDevices) return []
      return this.rawDeviceChildren.slice(0, MAX_DEVICES_IN_ORG_TREE)
    },
    siblingDevicesTruncated() {
      return this.hasMixedOrgAndDevices
        && this.rawDeviceChildren.length > this.paintedSiblingDevices.length
    },
    children() {
      if (this.kind === 'device' || this.kind === 'virtualCabinet') return []
      // 总览树：只铺组织层 +（混合时）直属 type=1 卡片；绝不在树内展开设备网格
      return [...this.orgChildren, ...this.paintedSiblingDevices]
    },
    hasExpandableChildren() {
      return this.children.length > 0
    },
    isLeaf() {
      return this.kind === 'organization' && this.orgChildren.length === 0
    },
    directDeviceCount() {
      return (this.node.children || []).filter(c => nodeType(c) === 1).length
    },
    deviceCount() {
      if (this.kind === 'device' || this.kind === 'virtualCabinet') return 1
      const count = (n) => {
        const t = nodeType(n)
        if (t === 1) return 1
        if (t === 'virtualCabinet') return 0
        return (n.children || []).reduce((v, c) => v + count(c), 0)
      }
      return count(this.node)
    },
    nodeClass() {
      return {
        'is-root': this.depth === 1 && this.kind === 'organization',
        'is-branch': this.kind === 'organization' && this.depth > 1 && !this.isLeaf,
        'is-leaf': this.kind === 'organization' && this.isLeaf,
        'is-device': this.kind === 'device',
        'is-cabinet': this.kind === 'virtualCabinet',
        'is-collapsed': this.hasExpandableChildren && !this.expanded,
      }
    },
    icon() {
      if (this.kind === 'virtualCabinet') return '▣'
      if (this.kind === 'device') return '🔌'
      if (this.depth === 1) return '⌁'
      return this.isLeaf ? '▰' : '⬡'
    },
    clickable() {
      if (this.kind === 'virtualCabinet' || this.kind === 'device') return true
      return this.deviceCount > 0 || this.orgChildren.length > 0
    },
    cardTitle() {
      if (this.kind === 'virtualCabinet') return '查看列头柜测点'
      if (this.kind === 'device') return '查看设备测点'
      if (this.hasMixedOrgAndDevices) {
        return this.siblingDevicesTruncated
          ? `含 ${this.orgChildren.length} 子组织 + ${this.directDeviceCount} 直属（展示前 ${MAX_DEVICES_IN_ORG_TREE}）`
          : `含 ${this.orgChildren.length} 子组织 + ${this.directDeviceCount} 直属配电室/设备`
      }
      if (this.isLeaf && this.directDeviceCount > 0) {
        return `进入设备列表（${this.directDeviceCount} 台）`
      }
      if (this.deviceCount > 0) return '查看组织设备列表'
      return this.name
    },
    /** 子节点多时走网格，避免横向无限拉宽 */
    isDenseLayout() {
      return this.children.length >= DENSE_CHILD_THRESHOLD
    },
    childrenWrapClass() {
      return {
        'is-dense': this.isDenseLayout,
        'is-device-grid': this.isDenseLayout && this.kind === 'organization',
      }
    },
  },
  methods: {
    stableNodeKey(node) {
      const value = (node && node.value) || {}
      if (value.type === 'virtualCabinet') {
        return `vc-${value.parentDeviceUuid || ''}-${value.virtualCabinet || node.text}`
      }
      return String(value.uuid || node.key || (value.sid != null ? value.sid : (node.text || 'n')))
    },
    toggleExpand(event) {
      if (event) {
        event.preventDefault()
        event.stopPropagation()
      }
      if (!this.hasExpandableChildren) return
      this.expanded = !this.expanded
      this.$nextTick(() => this.$emit('layout-dirty'))
    },
    openCard(event) {
      if (event && event.type === 'click' && this.openedByPointer) {
        this.openedByPointer = false
        return
      }
      this.handleOpen()
    },
    openCardFromPointer() {
      if (!this.clickable) return
      this.openedByPointer = true
      this.handleOpen()
    },
    handleOpen() {
      if (this.kind === 'virtualCabinet') {
        this.$emit('open-virtual-cabinet', this.node)
        return
      }
      if (this.kind === 'device') {
        // 总览树内不展开列头柜，直接进测点/设备页
        this.$emit('open-device', this.node)
        return
      }
      // 末级组织：直接进设备列表，禁止树内嵌套设备网格（重叠/布局崩坏的根因）
      if (this.isLeaf && this.directDeviceCount > 0) {
        this.$emit('open-device-list', this.node)
        return
      }
      if (this.deviceCount > 0) this.$emit('open-device-list', this.node)
    },
  },
  template: `<div class="org-node" :class="nodeClass">
    <div
      class="org-card"
      :class="{ 'is-clickable': clickable }"
      :role="clickable ? 'button' : null"
      :tabindex="clickable ? 0 : null"
      :title="cardTitle"
      @mousedown.left="openCardFromPointer"
      @click="openCard"
      @keydown.enter="openCard"
      @keydown.space.prevent="openCard"
    >
      <span class="org-card-corner corner-tl"></span><span class="org-card-corner corner-br"></span>
      <button
        v-if="hasExpandableChildren"
        type="button"
        class="org-expand-btn"
        :title="expanded ? '收起' : '展开'"
        :aria-expanded="expanded ? 'true' : 'false'"
        @mousedown.stop.prevent
        @click.stop="toggleExpand"
      >{{ expanded ? '−' : '+' }}</button>
      <div class="org-card-icon"><span>{{ icon }}</span></div>
      <div class="org-card-main">
        <strong :title="name">{{ name }}</strong>
      </div>
      <div class="org-card-actions" v-if="kind === 'organization' && (hasMixedOrgAndDevices || (isLeaf && directDeviceCount > 0))">
        <button type="button" class="org-chip org-chip--accent" :title="siblingDevicesTruncated ? '查看全部' : '打开设备列表'" @mousedown.stop.prevent @click.stop="$emit('open-device-list', node)">{{ siblingDevicesTruncated ? '全部' : '列表' }}</button>
      </div>
      <div class="org-card-count" v-else-if="kind === 'organization'">
        <b>{{ deviceCount }}</b><span>设备</span>
      </div>
    </div>
    <div v-if="hasExpandableChildren && expanded" class="org-children" :class="childrenWrapClass">
      <div v-if="isDenseLayout" class="org-dense-stem" aria-hidden="true"></div>
      <org-node
        v-for="child in children"
        :key="stableNodeKey(child)"
        :node="child"
        :depth="depth + 1"
        @open-device-list="$emit('open-device-list', $event)"
        @open-device="$emit('open-device', $event)"
        @open-virtual-cabinet="$emit('open-virtual-cabinet', $event)"
        @layout-dirty="$emit('layout-dirty')"
      />
    </div>
  </div>`,
}

export default {
  name: 'ScadaOrgOverview',
  components: { OrgNode },
  props: { projectUuid: { type: String, required: true }, modelId: { type: String, required: true } },
  data: () => ({
    roots: [],
    loading: true,
    currentPage: '',
    goPageHandler: null,
    treeScale: 1,
    treeLayoutReady: false,
    treeResizeObserver: null,
  }),
  computed: {
    visible() { return !this.currentPage || this.currentPage === this.modelId },
    treeStageStyle() {
      return {
        transform: `scale(${this.treeScale})`,
        visibility: this.treeLayoutReady ? 'visible' : 'hidden',
      }
    },
  },
  methods: {
    stableNodeKey(node) {
      const value = (node && node.value) || {}
      return String(value.uuid || node.key || (value.sid != null ? value.sid : (node.text || 'organization')))
    },
    openDeviceList(node) {
      const value = (node && node.value) || {}
      const businessId = value.uuid || node.key || (value.sid != null ? String(value.sid) : '')
      if (!businessId) return
      this.$EventBus.$emit('OpenOrgDeviceList', { sid: value.sid, businessId: String(businessId) })
    },
    openDevice(node) {
      const value = (node && node.value) || {}
      this.$EventBus.$emit('OpenDeviceDatapoints', {
        name: node.text || value.name,
        label: node.text || value.name,
        uuid: value.uuid || node.key || '',
        deviceUuid: value.uuid || node.key || '',
        muid: value.muid || '',
        modelUuid: value.muid || '',
        sid: value.sid,
        type: 1,
        kind: 'device',
        status: value.Status === 1 ? 'on' : 'off',
      })
    },
    openVirtualCabinet(node) {
      const value = (node && node.value) || {}
      const prefix = value.virtualCabinet || node.text || ''
      const parentLabel = value.parentDeviceLabel || ''
      const isFallback = !!(
        value.virtualCabinetFallback
        || value.isFallbackGroup
        || (parentLabel && prefix === parentLabel)
      )
      this.$EventBus.$emit('OpenDeviceDatapoints', {
        name: prefix,
        label: prefix,
        kind: 'virtualCabinet',
        virtualCabinet: prefix,
        virtualCabinetFallback: isFallback,
        isFallbackGroup: isFallback,
        parentDeviceLabel: parentLabel,
        parentDeviceUuid: value.parentDeviceUuid || value.uuid || '',
        uuid: value.parentDeviceUuid || value.uuid || '',
        deviceUuid: value.parentDeviceUuid || value.uuid || '',
        muid: value.muid || '',
        modelUuid: value.muid || '',
        sid: value.sid,
        type: 1,
      })
    },
    onLayoutDirty() {
      this.$nextTick(() => this.updateTreeScale())
    },
    updateConnectorGeometry() {
      const stage = this.$refs.orgTreeStage
      if (!stage) return
      const scale = Math.max(0.0001, Number(this.treeScale) || 1)
      stage.querySelectorAll('.org-children').forEach(childrenEl => {
        // 网格布局不画横向总线，跳过几何计算
        if (childrenEl.classList.contains('is-dense')) return
        const directNodes = Array.from(childrenEl.children)
          .filter(el => el.classList && el.classList.contains('org-node'))
        const directCards = directNodes
          .map(node => Array.from(node.children).find(el => el.classList && el.classList.contains('org-card')))
          .filter(Boolean)
        if (!directCards.length) return
        const childrenRect = childrenEl.getBoundingClientRect()
        const centerOf = card => {
          const rect = card.getBoundingClientRect()
          return (rect.left + rect.width / 2 - childrenRect.left) / scale
        }
        const start = centerOf(directCards[0])
        const end = centerOf(directCards[directCards.length - 1])
        childrenEl.style.setProperty('--branch-start-x', `${start}px`)
        childrenEl.style.setProperty('--branch-span', `${Math.max(0, end - start)}px`)
      })
    },
    updateTreeScale() {
      const tree = this.$refs.orgTree
      const stage = this.$refs.orgTreeStage
      if (!tree || !stage) return
      this.updateConnectorGeometry()
      const availableWidth = Math.max(1, tree.clientWidth - 16)
      const availableHeight = Math.max(1, tree.clientHeight - 16)
      // transform 不影响 layout 尺寸；scrollWidth 即未缩放自然宽
      const naturalWidth = Math.max(1, stage.scrollWidth, stage.offsetWidth)
      const naturalHeight = Math.max(1, stage.scrollHeight, stage.offsetHeight)
      const fit = Math.min(1, availableWidth / naturalWidth, availableHeight / naturalHeight)
      const hasDense = !!stage.querySelector('.org-children.is-dense')
      // 密集网格：保字号可读，超出靠滚动；稀疏组织树仍可小幅缩放适配
      const floor = hasDense ? MIN_TREE_SCALE : 0.72
      this.treeScale = Math.max(floor, Math.floor(fit * 10000) / 10000)
      this.treeLayoutReady = true
    },
    setupTreeLayout() {
      this.updateTreeScale()
      if (this.treeResizeObserver) this.treeResizeObserver.disconnect()
      if (typeof ResizeObserver === 'undefined') return
      // 只观察视口容器尺寸；勿观察 stage（内容/伪元素变化会误触缩放 → 滚动抖动）
      this.treeResizeObserver = new ResizeObserver(() => {
        if (this._treeScaleRaf) cancelAnimationFrame(this._treeScaleRaf)
        this._treeScaleRaf = requestAnimationFrame(() => {
          this._treeScaleRaf = 0
          this.updateTreeScale()
        })
      })
      if (this.$refs.orgTree) this.treeResizeObserver.observe(this.$refs.orgTree)
    },
  },
  mounted() {
    this.currentPage = this.modelId
    this.goPageHandler = data => {
      if (data && !data.IsPopUp && data.PageUuid) this.currentPage = data.PageUuid
    }
    this.$EventBus.$on('GoPage', this.goPageHandler)
    getMonitorTree({}, { headers: { ProjectUuid: this.projectUuid } })
      .then(res => { this.roots = res && res.data && res.data.code === 0 ? (res.data.list || []) : [] })
      .catch(() => { this.roots = [] })
      .finally(() => {
        this.loading = false
        this.$nextTick(() => this.setupTreeLayout())
      })
  },
  beforeDestroy() {
    if (this.goPageHandler) this.$EventBus.$off('GoPage', this.goPageHandler)
    if (this.treeResizeObserver) this.treeResizeObserver.disconnect()
    if (this._treeScaleRaf) cancelAnimationFrame(this._treeScaleRaf)
  },
}
</script>

<style>
/* 无顶部 KPI：content_top≈72/1080 → top≈6.7vh；高度占满至底边 */
.org-overview{position:absolute;left:1.65vw;top:calc(6.7vh + 6px);width:64.75vw;height:calc(90.8vh - 6px);z-index:40;display:flex;flex-direction:column;overflow:hidden;padding:0 18px 10px;box-sizing:border-box;background:transparent;border:0;box-shadow:none;clip-path:polygon(0 0,100% 0,100% 100%,18px 100%,0 calc(100% - 18px));color:#d9f6ff}
.org-overview:before,.org-overview:after{display:none}
.org-header{position:relative;height:56px;flex:none;display:flex;align-items:center;justify-content:space-between;margin:3px 10px 0;padding:0 14px;border:0;background:linear-gradient(90deg,rgba(7,48,66,.2),transparent 42%,transparent 68%,rgba(7,48,66,.12));clip-path:polygon(9px 0,100% 0,100% calc(100% - 9px),calc(100% - 9px) 100%,0 100%,0 9px)}
.org-header:before{content:"";position:absolute;left:0;top:0;width:72px;height:15px;border-left:1px solid rgba(48,221,255,.55);border-top:1px solid rgba(48,221,255,.55);pointer-events:none}.org-header:after{content:"";position:absolute;left:0;right:0;bottom:0;height:1px;background:linear-gradient(90deg,rgba(49,229,255,.62),rgba(49,229,255,.08) 30%,transparent 52%,rgba(49,229,255,.08) 78%,rgba(49,229,255,.38));box-shadow:0 0 7px rgba(24,216,255,.18);pointer-events:none}
.org-heading{position:relative;display:flex;align-items:center;gap:11px;padding-left:3px}.org-heading-icon{display:grid;place-items:center;width:30px;height:30px;border:1px solid rgba(34,231,255,.45);background:radial-gradient(circle,rgba(0,207,243,.16),transparent 68%);color:#3beaff;text-shadow:0 0 10px #00dfff;transform:rotate(45deg);box-shadow:0 0 12px rgba(0,210,255,.08)}.org-heading-icon::first-letter{transform:rotate(-45deg)}.org-heading div{display:flex;flex-direction:column}.org-heading b{font-size:15px;letter-spacing:2px;color:#dffaff;text-shadow:0 0 8px rgba(54,221,255,.2)}.org-heading div span{font-size:8px;letter-spacing:1.35px;color:#4f87a3;margin-top:3px}
.org-summary{display:flex;align-items:center;gap:3px;font-size:10px;color:#7297ac}.org-summary>span{position:relative;padding:5px 10px;border:0;border-left:1px solid rgba(60,188,218,.26);background:linear-gradient(90deg,rgba(8,47,65,.38),transparent)}.org-summary>span:last-of-type{border-right:1px solid rgba(60,188,218,.16)}.org-summary i{font-style:normal;color:#3ee8ff;font-size:13px;font-weight:700}.org-summary em{margin-left:8px;padding:4px 7px;font-style:normal;color:#4de0ad;text-shadow:0 0 6px rgba(77,224,173,.5);background:rgba(21,92,77,.08);animation:orgStatusPulse 2.2s ease-in-out infinite}
.org-empty{display:grid;place-items:center;flex:1;color:#7897b6}.org-tree{position:relative;flex:1;display:flex;align-items:center;justify-content:center;overflow:hidden;padding:26px 14px 32px;background-image:linear-gradient(rgba(30,116,143,.055) 1px,transparent 1px),linear-gradient(90deg,rgba(30,116,143,.055) 1px,transparent 1px);background-size:32px 32px}.org-tree:after{content:"";position:absolute;inset:0;width:auto;pointer-events:none;z-index:0;background:linear-gradient(90deg,transparent,rgba(20,215,255,.045),transparent);background-size:36px 100%;background-repeat:no-repeat;background-position:-40px 0;animation:orgRadarSweep 8s linear infinite;contain:strict}.org-tree-stage{position:relative;z-index:1;flex:none;display:flex;align-items:center;justify-content:center;width:max-content;transform-origin:top center}.org-tree-stage>.org-node{min-width:max-content;justify-content:center}.org-node{--org-card-width:clamp(190px,14vw,250px);--branch-gap:clamp(64px,9vh,92px);--branch-half-gap:clamp(32px,4.5vh,46px);position:relative;display:flex;flex-direction:column;align-items:center}.org-card{position:relative;display:flex;align-items:center;width:var(--org-card-width);min-height:74px;padding:9px 11px;box-sizing:border-box;overflow:hidden;contain:paint;border:1px solid rgba(43,197,230,.42);background:linear-gradient(135deg,rgba(8,45,64,.96),rgba(5,25,39,.98));box-shadow:inset 0 0 18px rgba(0,185,230,.08),0 4px 12px rgba(0,0,0,.22)}
.org-card.is-clickable{cursor:pointer;transition:border-color .18s ease,box-shadow .18s ease,transform .18s ease}.org-card.is-clickable:hover,.org-card.is-clickable:focus-visible{outline:none;border-color:rgba(83,235,255,.82);box-shadow:inset 0 0 22px rgba(0,210,255,.13),0 0 14px rgba(0,210,255,.18);transform:translateY(-1px)}.org-card:before{content:"";position:absolute;top:0;bottom:0;width:55px;left:0;transform:skewX(-18deg) translateX(-70px);background:linear-gradient(90deg,transparent,rgba(91,232,255,.12),transparent);animation:orgCardScan 5.2s ease-in-out infinite;will-change:transform,opacity;pointer-events:none}.org-card:after{content:"";position:absolute;left:42px;right:9px;bottom:0;height:1px;background:linear-gradient(90deg,#19dfff,transparent)}.org-card-corner{position:absolute;width:8px;height:8px}.corner-tl{top:-1px;left:-1px;border-top:2px solid #31eaff;border-left:2px solid #31eaff}.corner-br{right:-1px;bottom:-1px;border-right:2px solid #31eaff;border-bottom:2px solid #31eaff}
.org-card-icon{position:relative;flex:none;display:grid;place-items:center;width:38px;height:38px;border:1px solid rgba(31,221,255,.45);clip-path:polygon(50% 0,100% 25%,100% 75%,50% 100%,0 75%,0 25%);background:rgba(0,186,225,.12);color:#39e9ff;text-shadow:0 0 9px #00dfff;font-size:18px}.org-card-icon:after{content:"";position:absolute;inset:4px;border:1px dashed rgba(83,235,255,.55);border-radius:50%;animation:orgRotor 5s linear infinite}.org-card-main{min-width:0;flex:1;margin-left:9px}.org-card-main strong{display:block;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;color:#dffaff;font-size:12px;margin:0}.org-card-count{flex:none;display:flex;flex-direction:column;align-items:flex-end;margin-left:7px}.org-card-count b{font-size:19px;line-height:21px;color:#64edbf;text-shadow:0 0 8px rgba(65,230,180,.35)}.org-card-count span{font-size:8px;color:#6d9a91}.org-expand-btn{appearance:none;flex:none;position:relative;z-index:2;width:22px;height:22px;margin-right:6px;padding:0;border:1px solid rgba(60,200,230,.55);border-radius:2px;background:rgba(8,42,58,.92);color:#7ad8f5;font-size:14px;line-height:1;cursor:pointer;pointer-events:auto}.org-expand-btn:hover{border-color:rgba(80,230,255,.9);color:#dffaff;background:rgba(12,60,80,.98)}
.org-node.is-root>.org-card{width:clamp(240px,18vw,290px);min-height:80px;border-color:rgba(52,233,255,.7);background:linear-gradient(135deg,rgba(5,79,103,.92),rgba(5,29,45,.98));box-shadow:inset 0 0 25px rgba(0,216,255,.13),0 0 14px rgba(0,210,255,.12)}.org-node.is-root>.org-card .org-card-icon{width:44px;height:44px}.org-node.is-leaf>.org-card{border-color:rgba(70,215,172,.35)}.org-node.is-leaf>.org-card .org-card-icon{color:#64edbf;border-color:rgba(70,215,172,.45);background:rgba(45,201,150,.09)}
.org-node.is-device>.org-card{border-color:rgba(100,180,255,.45);background:linear-gradient(135deg,rgba(12,40,72,.96),rgba(5,22,40,.98))}
.org-node.is-device>.org-card .org-card-icon{color:#7ec8ff;border-color:rgba(100,180,255,.5);background:rgba(40,120,200,.12)}
.org-node.is-cabinet{--org-card-width:clamp(120px,8vw,160px)}
.org-node.is-cabinet>.org-card{min-height:64px;border-color:rgba(70,215,172,.4);background:linear-gradient(135deg,rgba(12,55,48,.95),rgba(5,28,32,.98))}
.org-node.is-cabinet>.org-card .org-card-icon{width:30px;height:30px;font-size:14px;color:#64edbf;border-color:rgba(70,215,172,.5);background:rgba(45,201,150,.1)}
.org-node.is-cabinet>.org-card .org-card-main strong{font-size:11px}
.org-children{position:relative;display:flex;justify-content:space-around;align-items:flex-start;width:100%;gap:clamp(14px,1.8vw,30px);padding-top:var(--branch-gap)}.org-children:before{content:"";position:absolute;top:var(--branch-half-gap);left:var(--branch-start-x,calc(var(--org-card-width) / 2));width:var(--branch-span,calc(100% - var(--org-card-width)));height:2px;background:linear-gradient(90deg,rgba(0,213,255,.15),#17dfff 35%,#8ef5ff 50%,#17dfff 65%,rgba(0,213,255,.15));background-size:180px 100%;box-shadow:0 0 7px rgba(0,213,255,.32);animation:orgCurrentHorizontal 2.4s linear infinite}.org-children>.org-node:before{content:"";position:absolute;top:calc(0px - var(--branch-half-gap));left:50%;width:2px;height:var(--branch-half-gap);background:linear-gradient(180deg,#17dfff,#a2f8ff,#17dfff);background-size:100% 40px;box-shadow:0 0 6px rgba(0,213,255,.35);animation:orgCurrentVertical 1.8s linear infinite}.org-node>.org-children:after{content:"";position:absolute;top:0;left:50%;width:2px;height:var(--branch-half-gap);background:linear-gradient(180deg,#17dfff,#a2f8ff,#17dfff);background-size:100% 40px;box-shadow:0 0 6px rgba(0,213,255,.35);animation:orgCurrentVertical 1.8s linear infinite}
.org-legend{height:26px;flex:none;display:flex;align-items:center;gap:12px;padding:0 8px 0 28px;border-top:0;color:#557d94;font-size:9px;flex-wrap:wrap}.org-legend span{display:flex;align-items:center;gap:5px}.org-legend span:last-child{margin-left:auto}.org-legend i{display:inline-block;width:7px;height:7px;border:1px solid}.legend-root{color:#2be7ff;background:#12647a}.legend-branch{color:#2ca8ca;background:#0c4359}.legend-leaf{color:#58dcb0;background:#145541}.legend-device{color:#7ec8ff;background:#1a4a7a}.legend-cabinet{color:#64edbf;background:#0d4a3a}
.org-tree::-webkit-scrollbar{width:6px;height:6px}.org-tree::-webkit-scrollbar-track{background:rgba(5,25,39,.7)}.org-tree::-webkit-scrollbar-thumb{background:rgba(0,210,240,.35);border-radius:4px}
.org-overview{padding:0}
.org-header{margin:0;padding:0 14px}
.org-tree{min-height:0;padding:8px;overflow:auto;box-sizing:border-box;align-items:flex-start;justify-content:flex-start;isolation:isolate}
.org-tree-stage{box-sizing:border-box;padding:6px 0;min-width:100%;min-height:100%;align-items:flex-start;justify-content:center}
.org-node{--org-card-width:clamp(150px,10.5vw,220px)}
.org-children{gap:clamp(8px,.8vw,16px)}
.org-legend{padding:0 10px}
.org-card{min-height:72px}
.org-card-main strong{display:flex;align-items:center;min-height:30px;margin:0;overflow:visible;text-overflow:clip;white-space:normal;word-break:break-all;line-height:1.25;font-size:clamp(10px,.62vw,12px)}
.org-card-actions{flex:none;display:flex;flex-direction:column;gap:4px;margin-left:6px}
.org-chip{appearance:none;cursor:pointer;border:1px solid rgba(60,200,230,.45);background:rgba(8,42,58,.85);color:#9ad8ef;font-size:9px;letter-spacing:.5px;padding:3px 7px;line-height:1.2;border-radius:2px}
.org-chip:hover{border-color:rgba(80,230,255,.8);color:#dffaff;background:rgba(12,60,80,.95)}
.org-chip--accent{border-color:rgba(70,215,172,.5);color:#7aefc8}
.org-chip--accent:hover{border-color:rgba(100,235,190,.85);color:#b8ffe0}
/* 多子节点：固定宽度网格，不再横向无限拉宽导致整树缩成蚂蚁 */
.org-children.is-dense{display:grid;grid-template-columns:repeat(auto-fill,minmax(132px,1fr));gap:8px 8px;width:min(860px,78vw);max-width:860px;padding:28px 6px 6px;justify-content:stretch;align-items:stretch;box-sizing:border-box}
.org-children.is-dense:before,.org-children.is-dense:after{display:none}
.org-children.is-dense>.org-node:before{display:none}
.org-children.is-dense>.org-node{width:100%;min-width:0;align-items:stretch}
.org-children.is-dense>.org-node>.org-card{width:100%;min-height:58px;padding:7px 8px}
.org-children.is-dense>.org-node>.org-card .org-card-icon{width:28px;height:28px;font-size:13px}
.org-children.is-dense>.org-node>.org-card .org-card-main{margin-left:7px}
.org-children.is-dense>.org-node>.org-card .org-card-main strong{min-height:auto;font-size:11px;margin:0;display:-webkit-box;-webkit-line-clamp:2;-webkit-box-orient:vertical;overflow:hidden;word-break:break-word}
.org-children.is-dense>.org-node>.org-card .org-card-count{display:none}
.org-children.is-dense>.org-node>.org-card .org-expand-btn{width:18px;height:18px;font-size:12px;margin-right:4px}
.org-dense-stem{position:absolute;top:0;left:50%;width:2px;height:22px;transform:translateX(-50%);background:linear-gradient(180deg,#17dfff,#a2f8ff);box-shadow:0 0 6px rgba(0,213,255,.35);pointer-events:none}
/* 光晕/扫描只用 transform 或 background-position，禁止动画 left/top（会撑大 scrollWidth → 滚动条抖动） */
@keyframes orgRotor{to{transform:rotate(360deg)}}
@keyframes orgCardScan{0%,15%{transform:skewX(-18deg) translateX(-70px);opacity:0}35%,60%{opacity:1}82%,100%{transform:skewX(-18deg) translateX(280px);opacity:0}}
@keyframes orgCurrentHorizontal{to{background-position:180px 0}}
@keyframes orgCurrentVertical{to{background-position:0 40px}}
@keyframes orgRadarSweep{from{background-position:-40px 0}to{background-position:calc(100% + 40px) 0}}
@keyframes orgStatusPulse{0%,100%{opacity:.65}50%{opacity:1;text-shadow:0 0 10px rgba(77,224,173,.8)}}
@media (prefers-reduced-motion:reduce){.org-tree:after,.org-card:before,.org-card-icon:after,.org-children:before,.org-children>.org-node:before,.org-node>.org-children:after,.org-summary em{animation:none}}
</style>
