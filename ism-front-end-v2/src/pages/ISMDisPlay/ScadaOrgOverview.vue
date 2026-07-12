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
      <div class="org-summary">
        <span><i>{{ organizationCount }}</i> 组织节点</span>
        <span><i>{{ totalDeviceCount }}</i> 设备总数</span>
        <span><i>{{ maxDepth }}</i> 组织层级</span>
        <em>● 实时同步</em>
      </div>
    </header>
    <div v-if="loading" class="org-empty">正在加载组织架构…</div>
    <div v-else-if="!roots.length" class="org-empty">暂无组织架构数据</div>
    <div v-else ref="orgTree" class="org-tree">
      <div ref="orgTreeStage" class="org-tree-stage" :style="treeStageStyle">
        <org-node
          v-for="node in roots"
          :key="node.key || node.id"
          :node="node"
          :depth="1"
          @open-device-list="openDeviceList"
        />
      </div>
    </div>
    <footer class="org-legend">
      <span><i class="legend-root"></i>系统根节点</span>
      <span><i class="legend-branch"></i>组织区域</span>
      <span><i class="legend-leaf"></i>末级组织</span>
      <span>点击含设备的组织卡片，可查看统一设备列表与点位数据</span>
    </footer>
  </section>
</template>

<script>
import { getMonitorTree } from '@/services/device'

const OrgNode = {
  name: 'OrgNode',
  props: {
    node: { type: Object, required: true },
    depth: { type: Number, default: 1 },
  },
  data: () => ({ openedByPointer: false }),
  computed: {
    // 组织总览只展开组织节点；设备作为叶层统计，不把数百台设备平铺成清单。
    children() { return (this.node.children || []).filter(child => (child.value || {}).type === 0) },
    value() { return this.node.value || {} },
    name() { return this.node.text || this.value.name || '未命名组织' },
    isLeaf() { return this.children.length === 0 },
    directDeviceCount() {
      return (this.node.children || []).filter(child => (child.value || {}).type === 1).length
    },
    deviceCount() {
      const count = n => (n.value || {}).type === 1 ? 1 : (n.children || []).reduce((v, c) => v + count(c), 0)
      return count(this.node)
    },
    nodeClass() {
      return {
        'is-root': this.depth === 1,
        'is-branch': this.depth > 1 && !this.isLeaf,
        'is-leaf': this.isLeaf,
      }
    },
    icon() {
      if (this.depth === 1) return '⌁'
      return this.isLeaf ? '▰' : '⬡'
    },
  },
  methods: {
    openDeviceList(event) {
      if (event && event.type === 'click' && this.openedByPointer) {
        this.openedByPointer = false
        return
      }
      if (this.deviceCount > 0) this.$emit('open-device-list', this.node)
    },
    openDeviceListFromPointer() {
      if (this.deviceCount <= 0) return
      this.openedByPointer = true
      this.$emit('open-device-list', this.node)
    },
  },
  template: `<div class="org-node" :class="nodeClass">
    <div
      class="org-card"
      :class="{ 'is-clickable': deviceCount > 0 }"
      :role="deviceCount > 0 ? 'button' : null"
      :tabindex="deviceCount > 0 ? 0 : null"
      :title="deviceCount > 0 ? '查看组织设备' : name"
      @mousedown.left="openDeviceListFromPointer"
      @click="openDeviceList"
      @keydown.enter="openDeviceList"
      @keydown.space.prevent="openDeviceList"
    >
      <span class="org-card-corner corner-tl"></span><span class="org-card-corner corner-br"></span>
      <div class="org-card-icon"><span>{{ icon }}</span></div>
      <div class="org-card-main">
        <small>LEVEL {{ depth }}</small>
        <strong :title="name">{{ name }}</strong>
        <div class="org-card-meta">
          <span>{{ children.length }} 个子组织</span>
          <span v-if="directDeviceCount">{{ directDeviceCount }} 台直属设备</span>
        </div>
      </div>
      <div class="org-card-count"><b>{{ deviceCount }}</b><span>设备</span></div>
    </div>
    <div v-if="children.length" class="org-children">
      <org-node
        v-for="child in children"
        :key="child.key || child.id"
        :node="child"
        :depth="depth + 1"
        @open-device-list="$emit('open-device-list', $event)"
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
    organizationCount() {
      const count = node => 1 + (node.children || [])
        .filter(child => (child.value || {}).type === 0)
        .reduce((sum, child) => sum + count(child), 0)
      return this.roots.reduce((sum, root) => sum + count(root), 0)
    },
    totalDeviceCount() {
      const count = node => (node.value || {}).type === 1
        ? 1
        : (node.children || []).reduce((sum, child) => sum + count(child), 0)
      return this.roots.reduce((sum, root) => sum + count(root), 0)
    },
    maxDepth() {
      const depth = node => {
        const children = (node.children || []).filter(child => (child.value || {}).type === 0)
        return children.length ? 1 + Math.max(...children.map(depth)) : 1
      }
      return this.roots.length ? Math.max(...this.roots.map(depth)) : 0
    },
  },
  methods: {
    openDeviceList(node) {
      const value = (node && node.value) || {}
      if (value.sid == null) return
      this.$EventBus.$emit('OpenOrgDeviceList', { sid: value.sid })
    },
    updateTreeScale() {
      const tree = this.$refs.orgTree
      const stage = this.$refs.orgTreeStage
      if (!tree || !stage) return
      const availableWidth = Math.max(1, tree.clientWidth - 16)
      const availableHeight = Math.max(1, tree.clientHeight - 16)
      const naturalWidth = Math.max(1, stage.scrollWidth, stage.offsetWidth)
      const naturalHeight = Math.max(1, stage.scrollHeight, stage.offsetHeight)
      const next = Math.min(1, availableWidth / naturalWidth, availableHeight / naturalHeight)
      this.treeScale = Math.floor(next * 10000) / 10000
      this.treeLayoutReady = true
    },
    setupTreeLayout() {
      this.updateTreeScale()
      if (this.treeResizeObserver) this.treeResizeObserver.disconnect()
      if (typeof ResizeObserver === 'undefined') return
      this.treeResizeObserver = new ResizeObserver(() => this.updateTreeScale())
      if (this.$refs.orgTree) this.treeResizeObserver.observe(this.$refs.orgTree)
      if (this.$refs.orgTreeStage) this.treeResizeObserver.observe(this.$refs.orgTreeStage)
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
  },
}
</script>

<style>
.org-overview{position:absolute;left:1.65vw;top:calc(17.2vh + 6px);width:64.75vw;height:calc(79.3vh - 6px);z-index:40;display:flex;flex-direction:column;overflow:hidden;padding:0 18px 10px;box-sizing:border-box;background:transparent;border:0;box-shadow:none;clip-path:polygon(0 0,100% 0,100% 100%,18px 100%,0 calc(100% - 18px));color:#d9f6ff}
.org-overview:before,.org-overview:after{display:none}
.org-header{position:relative;height:56px;flex:none;display:flex;align-items:center;justify-content:space-between;margin:3px 10px 0;padding:0 14px;border:0;background:linear-gradient(90deg,rgba(7,48,66,.2),transparent 42%,transparent 68%,rgba(7,48,66,.12));clip-path:polygon(9px 0,100% 0,100% calc(100% - 9px),calc(100% - 9px) 100%,0 100%,0 9px)}
.org-header:before{content:"";position:absolute;left:0;top:0;width:72px;height:15px;border-left:1px solid rgba(48,221,255,.55);border-top:1px solid rgba(48,221,255,.55);pointer-events:none}.org-header:after{content:"";position:absolute;left:0;right:0;bottom:0;height:1px;background:linear-gradient(90deg,rgba(49,229,255,.62),rgba(49,229,255,.08) 30%,transparent 52%,rgba(49,229,255,.08) 78%,rgba(49,229,255,.38));box-shadow:0 0 7px rgba(24,216,255,.18);pointer-events:none}
.org-heading{position:relative;display:flex;align-items:center;gap:11px;padding-left:3px}.org-heading-icon{display:grid;place-items:center;width:30px;height:30px;border:1px solid rgba(34,231,255,.45);background:radial-gradient(circle,rgba(0,207,243,.16),transparent 68%);color:#3beaff;text-shadow:0 0 10px #00dfff;transform:rotate(45deg);box-shadow:0 0 12px rgba(0,210,255,.08)}.org-heading-icon::first-letter{transform:rotate(-45deg)}.org-heading div{display:flex;flex-direction:column}.org-heading b{font-size:15px;letter-spacing:2px;color:#dffaff;text-shadow:0 0 8px rgba(54,221,255,.2)}.org-heading div span{font-size:8px;letter-spacing:1.35px;color:#4f87a3;margin-top:3px}
.org-summary{display:flex;align-items:center;gap:3px;font-size:10px;color:#7297ac}.org-summary>span{position:relative;padding:5px 10px;border:0;border-left:1px solid rgba(60,188,218,.26);background:linear-gradient(90deg,rgba(8,47,65,.38),transparent)}.org-summary>span:last-of-type{border-right:1px solid rgba(60,188,218,.16)}.org-summary i{font-style:normal;color:#3ee8ff;font-size:13px;font-weight:700}.org-summary em{margin-left:8px;padding:4px 7px;font-style:normal;color:#4de0ad;text-shadow:0 0 6px rgba(77,224,173,.5);background:rgba(21,92,77,.08);animation:orgStatusPulse 2.2s ease-in-out infinite}
.org-empty{display:grid;place-items:center;flex:1;color:#7897b6}.org-tree{position:relative;flex:1;display:flex;align-items:center;justify-content:center;overflow:hidden;padding:26px 14px 32px;background-image:linear-gradient(rgba(30,116,143,.055) 1px,transparent 1px),linear-gradient(90deg,rgba(30,116,143,.055) 1px,transparent 1px);background-size:32px 32px}.org-tree:after{content:"";position:absolute;top:0;bottom:0;width:36px;pointer-events:none;background:linear-gradient(90deg,transparent,rgba(20,215,255,.045),transparent);animation:orgRadarSweep 8s linear infinite}.org-tree-stage{flex:none;display:flex;align-items:center;justify-content:center;width:max-content;transform-origin:center center}.org-tree-stage>.org-node{min-width:max-content;justify-content:center}.org-node{--org-card-width:clamp(190px,14vw,250px);--branch-gap:clamp(64px,9vh,92px);--branch-half-gap:clamp(32px,4.5vh,46px);position:relative;display:flex;flex-direction:column;align-items:center}.org-card{position:relative;display:flex;align-items:center;width:var(--org-card-width);min-height:74px;padding:9px 11px;box-sizing:border-box;overflow:hidden;border:1px solid rgba(43,197,230,.42);background:linear-gradient(135deg,rgba(8,45,64,.96),rgba(5,25,39,.98));box-shadow:inset 0 0 18px rgba(0,185,230,.08),0 4px 12px rgba(0,0,0,.22)}
.org-card.is-clickable{cursor:pointer;transition:border-color .18s ease,box-shadow .18s ease,transform .18s ease}.org-card.is-clickable:hover,.org-card.is-clickable:focus-visible{outline:none;border-color:rgba(83,235,255,.82);box-shadow:inset 0 0 22px rgba(0,210,255,.13),0 0 14px rgba(0,210,255,.18);transform:translateY(-1px)}.org-card:before{content:"";position:absolute;top:0;bottom:0;width:55px;left:-70px;transform:skewX(-18deg);background:linear-gradient(90deg,transparent,rgba(91,232,255,.12),transparent);animation:orgCardScan 5.2s ease-in-out infinite}.org-card:after{content:"";position:absolute;left:42px;right:9px;bottom:0;height:1px;background:linear-gradient(90deg,#19dfff,transparent)}.org-card-corner{position:absolute;width:8px;height:8px}.corner-tl{top:-1px;left:-1px;border-top:2px solid #31eaff;border-left:2px solid #31eaff}.corner-br{right:-1px;bottom:-1px;border-right:2px solid #31eaff;border-bottom:2px solid #31eaff}
.org-card-icon{position:relative;flex:none;display:grid;place-items:center;width:38px;height:38px;border:1px solid rgba(31,221,255,.45);clip-path:polygon(50% 0,100% 25%,100% 75%,50% 100%,0 75%,0 25%);background:rgba(0,186,225,.12);color:#39e9ff;text-shadow:0 0 9px #00dfff;font-size:18px}.org-card-icon:after{content:"";position:absolute;inset:4px;border:1px dashed rgba(83,235,255,.55);border-radius:50%;animation:orgRotor 5s linear infinite}.org-card-main{min-width:0;flex:1;margin-left:9px}.org-card-main small{display:block;color:#467c98;font-size:8px;letter-spacing:1px}.org-card-main strong{display:block;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;color:#dffaff;font-size:12px;margin:4px 0}.org-card-meta{display:flex;gap:5px;color:#688fa6;font-size:8px}.org-card-count{flex:none;display:flex;flex-direction:column;align-items:flex-end;margin-left:7px}.org-card-count b{font-size:19px;line-height:21px;color:#64edbf;text-shadow:0 0 8px rgba(65,230,180,.35)}.org-card-count span{font-size:8px;color:#6d9a91}
.org-node.is-root>.org-card{width:clamp(240px,18vw,290px);min-height:80px;border-color:rgba(52,233,255,.7);background:linear-gradient(135deg,rgba(5,79,103,.92),rgba(5,29,45,.98));box-shadow:inset 0 0 25px rgba(0,216,255,.13),0 0 14px rgba(0,210,255,.12)}.org-node.is-root>.org-card .org-card-icon{width:44px;height:44px}.org-node.is-leaf>.org-card{border-color:rgba(70,215,172,.35)}.org-node.is-leaf>.org-card .org-card-icon{color:#64edbf;border-color:rgba(70,215,172,.45);background:rgba(45,201,150,.09)}
.org-children{position:relative;display:flex;justify-content:space-around;align-items:flex-start;width:100%;gap:clamp(14px,1.8vw,30px);padding-top:var(--branch-gap)}.org-children:before{content:"";position:absolute;top:var(--branch-half-gap);left:calc(var(--org-card-width) / 2);right:calc(var(--org-card-width) / 2);height:2px;background:linear-gradient(90deg,rgba(0,213,255,.15),#17dfff 35%,#8ef5ff 50%,#17dfff 65%,rgba(0,213,255,.15));background-size:180px 100%;box-shadow:0 0 7px rgba(0,213,255,.32);animation:orgCurrentHorizontal 2.4s linear infinite}.org-children>.org-node:before{content:"";position:absolute;top:calc(0px - var(--branch-half-gap));left:50%;width:2px;height:var(--branch-half-gap);background:linear-gradient(180deg,#17dfff,#a2f8ff,#17dfff);background-size:100% 40px;box-shadow:0 0 6px rgba(0,213,255,.35);animation:orgCurrentVertical 1.8s linear infinite}.org-node>.org-children:after{content:"";position:absolute;top:0;left:50%;width:2px;height:var(--branch-half-gap);transform:translateY(-100%);background:linear-gradient(180deg,#17dfff,#a2f8ff,#17dfff);background-size:100% 40px;box-shadow:0 0 6px rgba(0,213,255,.35);animation:orgCurrentVertical 1.8s linear infinite}
.org-legend{height:26px;flex:none;display:flex;align-items:center;gap:15px;padding:0 8px 0 28px;border-top:0;color:#557d94;font-size:9px}.org-legend span{display:flex;align-items:center;gap:5px}.org-legend span:last-child{margin-left:auto}.org-legend i{display:inline-block;width:7px;height:7px;border:1px solid}.legend-root{color:#2be7ff;background:#12647a}.legend-branch{color:#2ca8ca;background:#0c4359}.legend-leaf{color:#58dcb0;background:#145541}
.org-tree::-webkit-scrollbar{width:6px;height:6px}.org-tree::-webkit-scrollbar-track{background:rgba(5,25,39,.7)}.org-tree::-webkit-scrollbar-thumb{background:rgba(0,210,240,.35);border-radius:4px}
/* 组织名称属于关键导航信息：允许完整换行，不以省略号隐藏 3A1/3A2 等后缀。 */
.org-overview{padding:0}
.org-header{margin:0;padding:0 14px}
.org-tree{min-height:0;padding:8px;overflow:hidden;box-sizing:border-box}
.org-tree-stage{box-sizing:border-box;padding:6px 0}
.org-node{--org-card-width:clamp(150px,10.5vw,220px)}
.org-children{gap:clamp(8px,.8vw,16px)}
.org-legend{padding:0 10px}
.org-card{min-height:88px}
.org-card-main strong{display:flex;align-items:center;min-height:30px;margin:3px 0;overflow:visible;text-overflow:clip;white-space:normal;word-break:break-all;line-height:1.25;font-size:clamp(10px,.62vw,12px)}
.org-card-meta{flex-wrap:wrap;line-height:1.2}
@keyframes orgRotor{to{transform:rotate(360deg)}}@keyframes orgCardScan{0%,15%{left:-70px;opacity:0}35%,60%{opacity:1}82%,100%{left:115%;opacity:0}}@keyframes orgCurrentHorizontal{to{background-position:180px 0}}@keyframes orgCurrentVertical{to{background-position:0 40px}}@keyframes orgRadarSweep{from{left:-40px}to{left:100%}}@keyframes orgStatusPulse{0%,100%{opacity:.65}50%{opacity:1;text-shadow:0 0 10px rgba(77,224,173,.8)}}
@media (prefers-reduced-motion:reduce){.org-tree:after,.org-card:before,.org-card-icon:after,.org-children:before,.org-children>.org-node:before,.org-node>.org-children:after,.org-summary em{animation:none}}
</style>
