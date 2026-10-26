<template>
  <!-- 折叠态：仅一个竖向把手；展开态：完整树 -->
  <div
    class="ism-runtree"
    :class="{ collapsed }"
    :style="panelStyle"
  >
    <!-- 收起后的竖向把手 -->
    <div
      v-if="collapsed"
      class="rt-handle"
      role="button"
      tabindex="0"
      aria-label="展开导航树"
      title="展开导航树"
      @click="collapsed = false"
      @keydown.enter="collapsed = false"
      @keydown.space.prevent="collapsed = false"
    >
      <span class="rt-handle-icon">›</span>
    </div>

    <!-- 展开的树面板 -->
    <div v-else class="rt-panel">
      <div class="rt-head">
        <span class="rt-title">📍 设备导航</span>
        <div class="rt-head-actions">
          <span class="rt-act-btn" @click="collapseOneLayer" title="逐层折叠（收起当前最深层）">−</span>
          <span class="rt-act-btn" @click="collapseAllLayers" title="全部折叠（仅保留第一层）">⊟</span>
          <span class="rt-collapse-btn" @click="collapsed = true" title="收起导航面板">‹</span>
        </div>
      </div>

      <div ref="rtBody" class="rt-body" @scroll="syncHScroll">
        <div v-if="loading" class="rt-empty">⏳ 加载设备树中…</div>
        <div v-else-if="!(roots && roots.length)" class="rt-empty">⚠️ 暂无设备数据</div>

        <ul v-else ref="rtTree" class="rt-tree">
          <ism-runtree-node
            v-for="node in (roots || [])"
            :key="node.id"
            :node="node"
            :expanded-map="expandedMap"
            :selected-id="selectedId"
            @toggle="onToggle"
            @select="onSelect"
          />
        </ul>
      </div>

      <div class="rt-foot">
        <span class="rt-foot-total">共 {{ totalDevices }} 台设备</span>
        <span v-if="hScrollMax > 0" class="rt-foot-hint">↔ 可横向滚动查看全称</span>
      </div>
      <div v-if="hScrollMax > 0" class="rt-hscroll-bar">
        <span class="rt-hscroll-icon">◀</span>
        <input
          class="rt-hscroll"
          type="range"
          min="0"
          :max="hScrollMax"
          :value="hScrollLeft"
          @input="onHScrollInput"
        />
        <span class="rt-hscroll-icon">▶</span>
      </div>
    </div>
  </div>
</template>

<script>
import { getMonitorTree } from '@/services/device'
import { displayModelTemplateMap } from '@/services/displayModel'
import {
  pageIdDevice,
  buildNavTreeIndex,
  buildNavContextForNode,
  resolveTemplatePageIdForKind,
} from './utils/navTreeIndex'
import { normalizeRootNodes } from './utils/monitorTreeTransform'
import { resolveDeviceListTemplateId, deviceListPageSizeForNav } from './utils/deviceListPager'
import { resolveDeviceSignalTemplateId } from './utils/deviceSignalTemplate'
import {
  buildDeviceSignalContext,
  fetchDeviceDatapointPage,
  applyGatewayListPagination,
} from './utils/navContext'
import {
  ensureVirtualCabinetsForDevice,
  buildVirtualCabinetListContext,
  isVirtualCabinetNode,
  normalizeVirtualCabinetClick,
} from './utils/virtualCabinet'
import { DASHBOARD_PERFORMANCE } from './utils/dashboardPerformance'
import { detectDrillDepth, resolveOnSelectPageUuid, isDeviceNode } from './utils/drillDepth'
import store from '@/store'

// 异步加载子组件，避免与父 SFC 同 chunk 时 export default 错乱
const RunTreeNode = () => import(/* webpackChunkName: "ism-runtree-node" */ './ISMRunTreeNode.vue')

export default {
  name: 'ISMRunTreeNav',
  components: { 'ism-runtree-node': RunTreeNode },
  props: {
    projectUuid: { type: String, required: true },
    modelId: { type: String, required: true },
  },
  data() {
    return {
      loading: true,
      roots: [],
      expandedMap: {},   // id -> bool
      selectedId: '',
      collapsed: false,
      totalDevices: 0,
      hScrollLeft: 0,
      hScrollMax: 0,
      contentWidth: 168,
      templateMap: null,
      navIndex: null,
    }
  },
  mounted() {
    if (this.projectUuid) {
      try {
        sessionStorage.setItem('ProjectUuid', this.projectUuid)
      } catch (e) { /* ignore */ }
      import('js-cookie').then(({ default: Cookie }) => {
        Cookie.set('ProjectUuid', this.projectUuid)
      })
    }
    this.fetchTemplateMap()
    this.fetchTree()
    this.$EventBus.$on('OpenOrgDeviceList', this.onOrgDeviceList)
    this.$EventBus.$on('OpenDeviceDatapoints', this.onDeviceDatapoints)
    this.$EventBus.$on('ReturnToDeviceList', this.onReturnToDeviceList)
    this.$EventBus.$on('ReturnToGlobalOverview', this.onReturnToGlobalOverview)
    this.$nextTick(() => {
      this.measureWidth()
    })
  },
  beforeDestroy() {
    if (this._rtTreeRO) this._rtTreeRO.disconnect()
    this.$EventBus.$off('OpenOrgDeviceList', this.onOrgDeviceList)
    this.$EventBus.$off('OpenDeviceDatapoints', this.onDeviceDatapoints)
    this.$EventBus.$off('ReturnToDeviceList', this.onReturnToDeviceList)
    this.$EventBus.$off('ReturnToGlobalOverview', this.onReturnToGlobalOverview)
  },
  watch: {
    projectUuid(val) {
      if (val && !(this.roots && this.roots.length)) this.fetchTree()
    },
    collapsed(val) {
      if (!val) this.measureWidth()
    },
    loading(val) {
      if (!val) this.measureWidth()
    },
    expandedMap: {
      deep: true,
      handler() {
        this.measureWidth()
      },
    },
  },
  computed: {
    isXunanDash() {
      // 循安大屏已改用 build_ncc_dashboard 动态 page_id，需启用 uuid5 钻探
      return false
    },
    panelStyle() {
      if (this.collapsed) return {}
      const w = this.contentWidth
      return { width: `${w}px`, minWidth: `${w}px` }
    },
  },
  methods: {
    async fetchTemplateMap() {
      if (!this.modelId) return
      try {
        const res = await displayModelTemplateMap({ muid: this.modelId })
        if (res && res.data && res.data.code === 0 && res.data.map) {
          this.templateMap = res.data.map
          try {
            store.commit('ISMDisPlayEditorTool/setNavTemplateMap', res.data.map)
          } catch (e) { /* ignore */ }
        }
      } catch (e) {
        console.warn('[RunTree] templateMap failed', e && e.message)
      }
    },

    resolveTemplatePageId(node) {
      if (!this.templateMap) return ''
      return resolveTemplatePageIdForKind(this.templateMap, node.kind, '')
    },

    buildNavContext(node) {
      return buildNavContextForNode(node, this.navIndex || null)
    },

    async fetchTree() {
      this.loading = true
      try {
        // ProjectUuid 是后端用于隔离组织树的 HTTP 请求头，不是接口请求体字段。
        // 传入 params 会使后端无法识别项目范围，从而返回空树。
        const res = await getMonitorTree({}, { headers: { ProjectUuid: this.projectUuid } })
        if (res.data && res.data.code === 0 && Array.isArray(res.data.list)) {
          const normalized = normalizeRootNodes(res.data.list)
          this.roots = this.buildForest(normalized)
          // 建索引：旧 page_id 反查 + sid 定位（槽位重映射与 GoPage 兜底依赖）
          this.navIndex = buildNavTreeIndex(this.roots)
          try {
            store.commit('ISMDisPlayEditorTool/setNavTreeIndex', this.navIndex)
          } catch (e) { /* ignore */ }
          if (!this.roots.length) {
            const rawCount = res.data.list.length
            if (rawCount === 0) {
              console.warn('[RunTree] monitortree 返回空列表, project=', this.projectUuid)
            } else {
              console.warn('[RunTree] monitortree 有数据但构建为空, project=', this.projectUuid, 'raw=', rawCount)
            }
          }
        } else {
          console.warn('[RunTree] monitortree 异常:', res.data && res.data.code, 'project=', this.projectUuid)
        }
      } catch (e) {
        console.warn('[RunTree] 获取设备树失败:', e && e.message)
      }
      this.loading = false
      // 默认只展开后端返回的顶层节点；展开状态始终按业务标识保存。
      this.$nextTick(() => {
        this.expandTopLevels()
        this.measureWidth()
      })
    },

    /** 按后端原始父子关系递归构树，不推断组织名称、层数或节点顺序。 */
    buildForest(nodes, parentDepth = 0) {
      return normalizeRootNodes(nodes).map((node, index) =>
        this.transform(node, parentDepth, `root-${index}`),
      ).filter(Boolean)
    },

    transform(node, parentDepth = 0, fallbackId = '') {
      const v = node.value || {}
      const sid = v.sid
      const type = Number(v.type)
      const name = node.text || v.Name || '未命名'
      const treeDepth = parentDepth + 1
      const businessId = String(v.uuid || node.key || (sid != null ? sid : fallbackId))

      if (type === 1) {
        return {
          id: `device-${businessId}`,
          label: name,
          icon: '🔌',
          kind: 'device',
          layer: 'device',
          type: 1,
          treeDepth,
          sid,
          uuid: v.uuid || node.key || '',
          modelUuid: v.muid || '',
          muid: v.muid || '',
          status: v.Status === 1 ? 'on' : 'off',
          pageId: sid != null ? pageIdDevice(sid) : '',
          children: [],
        }
      }

      const children = (node.children || []).map((child, index) =>
        this.transform(child, treeDepth, `${businessId}-${index}`),
      ).filter(Boolean)
      const countDevices = items => items.reduce((sum, child) =>
        sum + (child.type === 1 ? 1 : countDevices(child.children || [])), 0)
      return {
        id: `organization-${businessId}`,
        rawLabel: name,
        label: name,
        icon: '▦',
        kind: parentDepth === 0 ? 'root' : 'organization',
        layer: 'organization',
        type: 0,
        treeDepth,
        sid,
        uuid: v.uuid || node.key || '',
        count: countDevices(children) || undefined,
        pageId: parentDepth === 0 ? this.modelId : '',
        children,
      }
    },

    expandTopLevels() {
      this.expandedMap = this.buildDefaultExpandedMap()
      // 统计设备总数
      let total = 0
      const count = nodes => nodes.forEach(n => {
        if (n.kind === 'device' || n.kind === 'gateway') total++
        if (n.children) count(n.children)
      })
      count(this.roots || [])
      this.totalDevices = total
    },

    /** 默认仅展开第一层（根节点） */
    buildDefaultExpandedMap() {
      const map = {}
      const walk = (nodes, depth) => {
        for (const n of nodes) {
          if (depth === 0 && n.children && n.children.length) map[n.id] = true
          if (n.children) walk(n.children, depth + 1)
        }
      }
      walk(this.roots || [], 0)
      return map
    },

    /** 全部折叠：回到默认，只保留第一层展开 */
    collapseAllLayers() {
      this.expandedMap = this.buildDefaultExpandedMap()
      this.measureWidth()
    },

    /** 逐层折叠：每次收起当前已展开的最深一层 */
    collapseOneLayer() {
      const expandedIds = Object.keys(this.expandedMap).filter(id => this.expandedMap[id])
      if (!expandedIds.length) return

      const depthOf = {}
      const walk = (nodes, depth) => {
        for (const n of nodes) {
          depthOf[n.id] = depth
          if (n.children && n.children.length) walk(n.children, depth + 1)
        }
      }
      walk(this.roots || [], 0)

      let maxDepth = -1
      for (const id of expandedIds) {
        const d = depthOf[id]
        if (d != null && d > maxDepth) maxDepth = d
      }
      if (maxDepth <= 0) return

      const map = { ...this.expandedMap }
      for (const id of expandedIds) {
        if (depthOf[id] === maxDepth) delete map[id]
      }
      this.expandedMap = map
      this.measureWidth()
    },

    async onToggle(id) {
      const willExpand = !this.expandedMap[id]
      if (willExpand) {
        const node = this.findNodeById(id)
        if (node && isDeviceNode(node) && !isVirtualCabinetNode(node)
          && !(node.children && node.children.length)
          && !node.virtualCabinetsChecked) {
          try {
            const fallbackName = node.label || node.rawLabel || node.name || ''
            const { enabled, cabinets } = await ensureVirtualCabinetsForDevice(node, fallbackName)
            this.$set(node, 'virtualCabinetsChecked', true)
            if (enabled && cabinets.length) {
              this.$set(node, 'children', cabinets)
              this.$set(node, 'hasVirtualCabinets', true)
            }
          } catch (e) {
            console.warn('[RunTree] virtual cabinets load failed', e && e.message)
            this.$set(node, 'virtualCabinetsChecked', true)
          }
        }
      }
      this.$set(this.expandedMap, id, willExpand)
      this.measureWidth()
    },

    findNodeById(id) {
      let found = null
      const walk = (nodes) => {
        for (const n of nodes || []) {
          if (n.id === id) {
            found = n
            return
          }
          if (n.children && n.children.length) walk(n.children)
          if (found) return
        }
      }
      walk(this.roots)
      return found
    },

    onOrgDeviceList(payload) {
      const sid = payload && payload.sid
      const businessId = payload && payload.businessId
      const node = sid != null && this.navIndex && this.navIndex.bySid
        ? this.navIndex.bySid[sid]
        : (businessId && this.navIndex && this.navIndex.byBusinessId
          ? this.navIndex.byBusinessId[String(businessId)]
          : null)
      if (!node) {
        this.$message && this.$message.warning('设备导航数据尚未加载完成，请稍后重试')
        return
      }
      // 画面层级：配电室 → 直属设备列表 → 单设备点位。
      // 有直属设备时绝不扁平化子孙点位；仅当直属全是组织容器时才下钻收集设备。
      const directChildren = node.children || []
      const directDevices = directChildren.filter(c => c && isDeviceNode(c))
      const organizationDevices = []
      if (directDevices.length) {
        organizationDevices.push(...directDevices)
      } else {
        const collectDevices = current => {
          ;(current.children || []).forEach(child => {
            if (!child) return
            if (isDeviceNode(child)) organizationDevices.push(child)
            else collectDevices(child)
          })
        }
        collectDevices(node)
      }
      if (!organizationDevices.length) {
        this.$message && this.$message.warning('该组织暂无设备')
        return
      }
      this.onSelect({ ...node, children: organizationDevices }, { forceDeviceList: true })
    },

    onDeviceDatapoints(device) {
      if (!device) return
      const currentNav = store.state.ISMDisPlayEditorTool
        ? store.state.ISMDisPlayEditorTool.navContext
        : null
      const returnContext = device.deviceListReturnContext
        || (currentNav && (currentNav.deviceListMode || currentNav.virtualCabinetListMode)
          ? { ...currentNav }
          : null)

      if (isVirtualCabinetNode(device) || device.virtualCabinet
        || (currentNav && currentNav.virtualCabinetListMode)) {
        const node = normalizeVirtualCabinetClick(device, currentNav)
          || {
            id: device.id || `vc-${device.uuid}-${device.name}`,
            label: device.virtualCabinet || device.name || device.label || '',
            name: device.virtualCabinet || device.name || device.label || '',
            kind: 'virtualCabinet',
            layer: 'device',
            type: 1,
            virtualCabinet: device.virtualCabinet || device.name || device.label || '',
            parentDeviceLabel: device.parentDeviceLabel || (currentNav && currentNav.name) || '',
            parentDeviceUuid: device.parentDeviceUuid
              || (currentNav && (currentNav.deviceUuid || currentNav.uuid))
              || device.deviceUuid || device.uuid || '',
            uuid: device.parentDeviceUuid
              || (currentNav && (currentNav.deviceUuid || currentNav.uuid))
              || device.deviceUuid || device.uuid || '',
            deviceUuid: device.parentDeviceUuid
              || (currentNav && (currentNav.deviceUuid || currentNav.uuid))
              || device.deviceUuid || device.uuid || '',
            modelUuid: device.modelUuid || device.muid || (currentNav && (currentNav.modelUuid || currentNav.muid)) || '',
            muid: device.muid || device.modelUuid || (currentNav && (currentNav.muid || currentNav.modelUuid)) || '',
            sid: device.sid != null ? device.sid : (currentNav && currentNav.sid),
            status: device.status || 'off',
          }
        if (!node.parentDeviceUuid && !node.uuid) {
          console.warn('[RunTree] virtual cabinet click missing parent uuid', device)
        }
        this.onSelect(node, { returnContext })
        return
      }

      const sid = device.sid
      const indexedNode = sid != null && this.navIndex && this.navIndex.bySid
        ? this.navIndex.bySid[sid]
        : null
      const node = indexedNode || {
        ...device,
        id: device.id || `dev-${sid || device.uuid || device.name}`,
        label: device.label || device.name || '',
        kind: 'device',
        layer: 'device',
        type: 1,
        modelUuid: device.modelUuid || device.muid || '',
      }
      this.onSelect(node, { returnContext })
    },

    onReturnToDeviceList() {
      const currentNav = store.state.ISMDisPlayEditorTool
        ? store.state.ISMDisPlayEditorTool.navContext
        : null
      const returnContext = currentNav && currentNav.deviceListReturnContext
      if (returnContext && (returnContext.deviceListMode || returnContext.virtualCabinetListMode)) {
        const pageUuid = resolveDeviceListTemplateId(this.templateMap)
        if (!pageUuid) {
          this.onReturnToGlobalOverview()
          return
        }
        try {
          store.commit('ISMDisPlayEditorTool/setNavContext', returnContext)
        } catch (e) { /* ignore */ }
        this.$EventBus.$emit('GoPage', {
          ModelId: this.modelId,
          PageUuid: pageUuid,
          IsPopUp: false,
          AutoClose: false,
          linkType: 'Inside',
          navContext: returnContext,
        })
        this.collapsed = true
        return
      }
      // 无列表返回栈（树直进测点 / 上下文丢失）：回全局总览，避免按钮「没反应」
      this.onReturnToGlobalOverview()
    },

    onReturnToGlobalOverview() {
      const homePageUuid = this.modelId
      if (!homePageUuid) return
      try {
        store.commit('ISMDisPlayEditorTool/setNavContext', null)
      } catch (e) { /* ignore */ }
      this.selectedId = null
      this.$EventBus.$emit('GoPage', {
        ModelId: homePageUuid,
        PageUuid: homePageUuid,
        IsPopUp: false,
        AutoClose: false,
        linkType: 'Inside',
        navContext: null,
      })
      this.collapsed = true
    },

    async onSelect(node, options = {}) {
      const forceDeviceList = !!options.forceDeviceList
      // 组织节点只承担展开/收缩职责，禁止跳入旧的 zone/room/building 页面。
      // 组织总览卡片可显式进入直属设备列表，设备叶节点进入按需加载的测点详情。
      if (!isDeviceNode(node) && !isVirtualCabinetNode(node) && !forceDeviceList) {
        if (node.children && node.children.length) this.onToggle(node.id)
        return
      }
      const depth = detectDrillDepth(node, this.navIndex)
      let pageUuid = ''
      let navContext = this.buildNavContext(node)
      navContext = {
        ...navContext,
        routeMode: depth.routeMode,
        drillDepth: depth.totalLayers,
        orgDepth: depth.orgDepth,
        treeDepth: depth.treeDepth,
        remainingLayers: depth.remainingLayers,
      }

      if (forceDeviceList) {
        navContext = applyGatewayListPagination({
          ...navContext,
          routeMode: 'childrenList',
          pageSize: deviceListPageSizeForNav(navContext),
        })
        pageUuid = resolveDeviceListTemplateId(this.templateMap)
      } else if (isVirtualCabinetNode(node)) {
        const muid = node.modelUuid || node.muid || ''
        const prefix = node.virtualCabinet || node.label || node.name || ''
        const parentUuid = node.parentDeviceUuid || node.deviceUuid || node.uuid || ''
        const parentLabel = node.parentDeviceLabel || ''
        const isFallback = !!(
          node.virtualCabinetFallback
          || node.isFallbackGroup
          || (parentLabel && prefix === parentLabel)
        )
        const pointPage = parentUuid || muid
          ? await fetchDeviceDatapointPage({
            muid,
            deviceUuid: parentUuid,
            pointNamePrefix: prefix,
            isFallbackGroup: isFallback,
            page: 1,
            pageSize: DASHBOARD_PERFORMANCE.datapointPageSize,
          })
          : { points: [], total: 0, page: 1, pageSize: DASHBOARD_PERFORMANCE.datapointPageSize }
        const ancestors = [
          ...(navContext.ancestors || []),
          ...(parentLabel ? [{ label: parentLabel, name: parentLabel, kind: 'device' }] : []),
        ]
        navContext = buildDeviceSignalContext({
          ...node,
          label: prefix,
          name: prefix,
          uuid: parentUuid,
          deviceUuid: parentUuid,
          modelUuid: muid,
          muid,
          virtualCabinet: prefix,
          virtualCabinetFallback: isFallback,
          parentDeviceLabel: parentLabel,
          parentDeviceUuid: parentUuid,
        }, pointPage.points, ancestors, {
          ...pointPage,
          serverPaged: true,
          virtualCabinet: prefix,
          virtualCabinetFallback: isFallback,
          parentDeviceLabel: parentLabel,
          parentDeviceUuid: parentUuid,
          category: isFallback ? '' : `${prefix}_`,
        })
        navContext = {
          ...navContext,
          virtualCabinet: prefix,
          virtualCabinetFallback: isFallback,
          parentDeviceLabel: parentLabel,
          parentDeviceUuid: parentUuid,
        }
        if (options.returnContext && (options.returnContext.deviceListMode || options.returnContext.virtualCabinetListMode)) {
          navContext = {
            ...navContext,
            deviceListReturnContext: options.returnContext,
          }
        }
        pageUuid = resolveDeviceSignalTemplateId(this.templateMap)
          || resolveOnSelectPageUuid(node, this.templateMap, this.navIndex)
      } else if (depth.routeMode === 'signal' || isDeviceNode(node)) {
        const muid = node.modelUuid || node.muid || ''
        const devLabel = node.label || node.rawLabel || node.name || ''
        // 多设备名（数据仓库 last `_`）时先虚拟设备列表，不直进全量测点
        let cabinetInfo = null
        try {
          cabinetInfo = await ensureVirtualCabinetsForDevice(
            node,
            node.label || node.rawLabel || node.name || '',
          )
        } catch (e) {
          console.warn('[RunTree] virtual cabinet probe failed', e && e.message)
        }
        if (cabinetInfo && cabinetInfo.enabled && cabinetInfo.cabinets.length) {
          if (!(node.children && node.children.length)) {
            this.$set(node, 'children', cabinetInfo.cabinets)
            this.$set(node, 'hasVirtualCabinets', true)
            this.$set(node, 'virtualCabinetsChecked', true)
          }
          const ancestors = (navContext.ancestors || [])
          navContext = applyGatewayListPagination(
            buildVirtualCabinetListContext(node, cabinetInfo.cabinets, ancestors, {
              pageSize: deviceListPageSizeForNav(navContext),
              homePageUuid: this.modelId,
              deviceListReturnContext: options.returnContext || null,
            }),
            deviceListPageSizeForNav(navContext),
          )
          pageUuid = resolveDeviceListTemplateId(this.templateMap)
        } else {
          const pointPage = muid || node.uuid || node.deviceUuid
            ? await fetchDeviceDatapointPage({
              muid,
              deviceLabel: '',
              deviceUuid: node.uuid || node.deviceUuid || '',
              page: 1,
              pageSize: DASHBOARD_PERFORMANCE.datapointPageSize,
            })
            : { points: [], total: 0, page: 1, pageSize: DASHBOARD_PERFORMANCE.datapointPageSize }
          const ancestors = (navContext.ancestors || [])
          navContext = buildDeviceSignalContext(node, pointPage.points, ancestors, {
            ...pointPage,
            serverPaged: true,
          })
          {
            const curNav = store.state.ISMDisPlayEditorTool
              ? store.state.ISMDisPlayEditorTool.navContext
              : null
            const listReturn = options.returnContext
              || (curNav && (curNav.deviceListMode || curNav.virtualCabinetListMode) ? { ...curNav } : null)
            if (listReturn && (listReturn.deviceListMode || listReturn.virtualCabinetListMode)) {
              navContext = {
                ...navContext,
                deviceListReturnContext: listReturn,
              }
            }
          }
          if (!pointPage.points.length && !pointPage.total) {
            console.warn('[ISMRunTreeNav] device datapoints empty', {
              muid, devLabel, nodeId: node.id, label: node.label, rawLabel: node.rawLabel,
            })
          }
          pageUuid = resolveDeviceSignalTemplateId(this.templateMap)
            || resolveOnSelectPageUuid(node, this.templateMap, this.navIndex)
        }
      } else if (depth.routeMode === 'childrenList') {
        const listPageSize = deviceListPageSizeForNav(navContext)
        navContext = applyGatewayListPagination({
          ...navContext,
          pageSize: listPageSize,
        })
        pageUuid = resolveDeviceListTemplateId(this.templateMap)
          || resolveOnSelectPageUuid(node, this.templateMap, this.navIndex)
      } else {
        pageUuid = resolveOnSelectPageUuid(node, this.templateMap, this.navIndex)
          || resolveDeviceListTemplateId(this.templateMap)
      }

      // 固定系统入口必须回当前展示模型首页，而不是模板映射中的旧 home 快照。
      navContext = {
        ...navContext,
        homePageUuid: this.modelId,
      }

      if (!pageUuid) {
        if (node.children && node.children.length) this.onToggle(node.id)
        const isDevice = depth.routeMode === 'signal' || isDeviceNode(node) || isVirtualCabinetNode(node)
        const hint = isDevice
          ? '未配置设备测点模板页，请先执行 scripts/bootstrap_device_signal_template_mysql.sh'
          : '未配置设备列表模板页，请先执行 scripts/bootstrap_device_list_template_mysql.sh'
        this.$message && this.$message.warning(hint)
        return
      }
      this.selectedId = node.id
      try {
        store.commit('ISMDisPlayEditorTool/setNavContext', navContext)
      } catch (e) { /* ignore */ }
      this.$EventBus.$emit('GoPage', {
        ModelId: this.modelId,
        PageUuid: pageUuid,
        IsPopUp: false,
        AutoClose: false,
        linkType: 'Inside',
        navContext,
      })
      this.collapsed = true
    },

    /** 按当前可见树内容实测宽度，避免侧栏过宽 */
    measureWidth() {
      if (this.collapsed) return
      this.$nextTick(() => {
        const tree = this.$refs.rtTree
        if (tree && !this._rtTreeRO) {
          this._rtTreeRO = new ResizeObserver(() => this.measureWidth())
          this._rtTreeRO.observe(tree)
        }
        const head = this.$el && this.$el.querySelector('.rt-head')
        const foot = this.$el && this.$el.querySelector('.rt-foot')
        let w = 0
        if (tree) w = Math.max(w, tree.scrollWidth)
        if (head) w = Math.max(w, head.scrollWidth)
        if (foot) w = Math.max(w, foot.scrollWidth)
        const bodyPad = 12
        const minW = 168
        const maxW = Math.min(Math.floor(window.innerWidth * 0.4), 480)
        const next = Math.min(maxW, Math.max(minW, w + bodyPad))
        // 幂等：宽度无实质变化时直接返回，避免 ResizeObserver 自反馈循环导致宽度逐帧增长
        if (Math.abs(next - this.contentWidth) < 2) return
        this.contentWidth = next
        this.syncHScroll()
      })
    },

    syncHScroll() {
      const el = this.$refs.rtBody
      if (!el) return
      this.hScrollLeft = el.scrollLeft
      this.hScrollMax = Math.max(0, el.scrollWidth - el.clientWidth)
    },

    onHScrollInput(e) {
      const el = this.$refs.rtBody
      if (!el) return
      const v = Number(e.target.value)
      el.scrollLeft = v
      this.hScrollLeft = v
    },
  },
}
</script>

<style scoped>
/* 覆盖在 AppRun 画布左侧栏区域之上（autoSize=1 → 画布铺满视口，按比例对齐） */
.ism-runtree {
  position: absolute;
  left: 0;
  top: 5.18vh;            /* 56 / 1080 */
  max-width: 40vw;
  height: 94.81vh;         /* (1080-56) / 1080 */
  z-index: 10050;
  font-family: 'Microsoft YaHei', 'PingFang SC', sans-serif;
  pointer-events: none;    /* 容器透传，子元素各自开启 */
  transition: width 0.2s ease, min-width 0.2s ease;
}
.ism-runtree.collapsed { width: 16px; min-width: 16px; max-width: 16px; }

.rt-handle,
.rt-panel { pointer-events: auto; }

/* 收起把手 */
.rt-handle {
  width: 14px;
  height: 34px;
  margin-top: 8px;
  background: rgba(8, 28, 43, 0.34);
  border: 1px solid rgba(53, 199, 226, 0.22);
  border-left: none;
  border-radius: 0 5px 5px 0;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: rgba(91, 231, 250, 0.68);
  opacity: 0.42;
  backdrop-filter: blur(2px);
  box-shadow: 0 0 5px rgba(0, 229, 255, 0.08);
  transition: width 0.18s ease, opacity 0.18s ease, background 0.18s ease;
}
.rt-handle:hover,
.rt-handle:focus-visible {
  width: 16px;
  opacity: 0.92;
  outline: none;
  background: rgba(10, 48, 67, 0.66);
  box-shadow: 0 0 8px rgba(0, 229, 255, 0.18);
}
.rt-handle-icon { font-size: 14px; line-height: 1; font-weight: 700; }

/* 树面板：深色渐变底 + 内发光 + 科技角标 */
.rt-panel {
  position: relative;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  background: linear-gradient(180deg, rgba(11, 19, 34, 0.97), rgba(8, 14, 26, 0.97));
  border-right: 1px solid rgba(30, 58, 95, 0.85);
  box-shadow: 2px 0 18px rgba(0, 0, 0, 0.4), inset 0 0 30px rgba(0, 50, 90, 0.12);
}
/* 四角科技角标（左上 + 右下两枚 L 形括号，点到为止） */
.rt-panel::before,
.rt-panel::after {
  content: '';
  position: absolute;
  width: 13px;
  height: 13px;
  border: 1.5px solid rgba(0, 229, 255, 0.55);
  pointer-events: none;
  z-index: 3;
}
.rt-panel::before { top: 4px; left: 4px; border-right: none; border-bottom: none; }
.rt-panel::after { bottom: 4px; right: 4px; border-left: none; border-top: none; }

.rt-head {
  position: relative;
  height: 40px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 6px;
  padding: 0 8px 0 12px;
  background: linear-gradient(90deg, rgba(0, 229, 255, 0.06), transparent);
}
.rt-head-actions {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
}
.rt-act-btn {
  cursor: pointer;
  color: #7da0c8;
  font-size: 13px;
  font-weight: 700;
  width: 20px;
  height: 20px;
  line-height: 18px;
  text-align: center;
  border: 1px solid rgba(30, 58, 95, 0.8);
  border-radius: 4px;
  transition: background 0.15s, color 0.15s, box-shadow 0.15s;
}
.rt-act-btn:hover {
  color: #00e5ff;
  background: rgba(0, 229, 255, 0.1);
  box-shadow: 0 0 6px rgba(0, 229, 255, 0.2);
}
/* 标题下渐变分隔线（青→透明） */
.rt-head::after {
  content: '';
  position: absolute;
  left: 0;
  bottom: 0;
  width: 100%;
  height: 1px;
  /* 底色渐变线 + 一段缓慢横向流动的高光（4.5s 周期） */
  background:
    linear-gradient(90deg, rgba(0, 229, 255, 0.5), rgba(0, 229, 255, 0.08) 60%, transparent),
    linear-gradient(90deg, transparent 0%, rgba(150, 240, 255, 0.95) 50%, transparent 100%);
  background-size: 100% 100%, 42% 100%;
  background-repeat: no-repeat, no-repeat;
  background-position: 0 0, -45% 0;
  animation: rtUnderlineFlow 4.5s linear infinite;
}
@keyframes rtUnderlineFlow {
  0%   { background-position: 0 0, -45% 0; }
  100% { background-position: 0 0, 145% 0; }
}
/* 标题：轻微 glow + 缓慢呼吸（3.2s，避免高频闪烁） */
.rt-title {
  font-size: 13px;
  font-weight: 600;
  color: #bfe9ff;
  letter-spacing: 1px;
  animation: rtTitleBreath 3.2s ease-in-out infinite;
}
@keyframes rtTitleBreath {
  0%, 100% { text-shadow: 0 0 6px rgba(0, 229, 255, 0.3); }
  50% { text-shadow: 0 0 12px rgba(0, 229, 255, 0.6); }
}
.rt-collapse-btn {
  cursor: pointer; color: #00e5ff; font-size: 18px; font-weight: 700;
  width: 22px; height: 22px; line-height: 20px; text-align: center;
  border: 1px solid rgba(30, 58, 95, 0.8); border-radius: 4px;
  transition: background 0.15s, box-shadow 0.15s;
}
.rt-collapse-btn:hover { background: rgba(0, 229, 255, 0.12); box-shadow: 0 0 8px rgba(0, 229, 255, 0.25); }

/* 树体：极淡科技网格底（低透明度，不干扰阅读） */
.rt-body {
  position: relative;
  flex: 1;
  overflow-y: auto;
  overflow-x: auto;
  padding: 8px 6px;
  background-image:
    linear-gradient(rgba(30, 58, 95, 0.05) 1px, transparent 1px),
    linear-gradient(90deg, rgba(30, 58, 95, 0.05) 1px, transparent 1px);
  background-size: 22px 22px;
}
.rt-tree {
  list-style: none;
  margin: 0;
  padding: 0;
  width: max-content;
}
.rt-empty { text-align: center; color: #5f7799; padding: 24px 10px; font-size: 12px; }

.rt-foot {
  position: relative;
  height: 28px; flex-shrink: 0;
  display: flex; align-items: center; justify-content: center;
  gap: 8px;
  background: linear-gradient(90deg, transparent, rgba(0, 229, 255, 0.04), transparent);
}
/* 底部渐变分隔线 */
.rt-foot::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  width: 100%;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(0, 229, 255, 0.35), transparent);
}
.rt-foot-total { font-size: 11px; color: #7da0c8; letter-spacing: 0.5px; }
.rt-foot-hint { font-size: 10px; color: #5f7799; white-space: nowrap; }

.rt-hscroll-bar {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px 8px;
  pointer-events: auto;
}
.rt-hscroll-icon {
  font-size: 9px;
  color: #5f7799;
  flex-shrink: 0;
  user-select: none;
}
.rt-hscroll {
  flex: 1;
  height: 4px;
  margin: 0;
  accent-color: #00e5ff;
  cursor: ew-resize;
}

/* 滚动条（纵向 + 横向） */
.rt-body::-webkit-scrollbar { width: 6px; height: 6px; }
.rt-body::-webkit-scrollbar-track { background: transparent; }
.rt-body::-webkit-scrollbar-thumb { background: #1e3a5f; border-radius: 3px; }
.rt-body::-webkit-scrollbar-thumb:hover { background: #2c5a8f; }
</style>
