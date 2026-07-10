<template>
  <!-- 折叠态：仅一个竖向把手；展开态：完整树 -->
  <div
    class="ism-runtree"
    :class="{ collapsed }"
    :style="panelStyle"
  >
    <!-- 收起后的竖向把手 -->
    <div v-if="collapsed" class="rt-handle" @click="collapsed = false" title="展开导航树">
      <span class="rt-handle-icon">›</span>
      <span class="rt-handle-text">导航</span>
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
import { resolveXunanSubstationPageId, XUNAN_MODEL_ID } from '@/config/xunanDashboardPages'
import {
  pageIdBuilding,
  pageIdDevice,
  pageIdRoom,
  pageIdZone,
  buildNavTreeIndex,
  buildNavContextForNode,
  resolveTemplatePageIdForKind,
} from './utils/navTreeIndex'
import {
  normalizeRootNodes,
  resolveMonitorNodeKind,
  iconForMonitorKind,
  isRoomZoneName,
  countDevicesInSubtree,
} from './utils/monitorTreeTransform'
import { resolveDeviceListTemplateId, deviceListPageSizeForNav } from './utils/deviceListPager'
import { resolveDeviceSignalTemplateId } from './utils/deviceSignalTemplate'
import {
  buildDeviceSignalContext,
  fetchDeviceDatapoints,
  applyGatewayListPagination,
} from './utils/navContext'
import { detectDrillDepth, resolveOnSelectPageUuid, isDeviceNode } from './utils/drillDepth'
import store from '@/store'

// 异步加载子组件，避免与父 SFC 同 chunk 时 export default 错乱
const RunTreeNode = () => import(/* webpackChunkName: "ism-runtree-node" */ './ISMRunTreeNode.vue')

/** 提取「X模块 / X配电室」的编码 X，合并后统一展示为 X变电所 */
function zoneDisplayKey(name) {
  const n = String(name || '').trim()
  const m = n.match(/^(.+?)(模块|配电室)/)
  return m ? m[1] : null
}

function substationKeyFromName(name) {
  const n = String(name || '').trim()
  if (!n) return null
  if (n.toUpperCase().startsWith('ECC')) return 'ECC'
  const head = n.split(/[-_及]/)[0].replace(/配电室/g, '').replace(/模块/g, '').trim()
  let m = head.match(/(\d+[AB]\d+)/i)
  if (m) return m[1].toUpperCase()
  const compact = n.replace(/[^0-9A-Za-z]/g, '')
  m = compact.match(/(\d+[AB]\d+)/i)
  return m ? m[1].toUpperCase() : null
}

/** 界面展示：库内「4A1配电室」→「4A1变电所」 */
function displaySubstationName(name) {
  const key = substationKeyFromName(name)
  if (key) return `${key}变电所`
  const n = String(name || '').trim()
  if (!n) return '变电所'
  if (n.toUpperCase().startsWith('ECC')) return 'ECC变电所'
  const m2 = n.match(/^(\d+[AB])/i)
  if (m2) return `${m2[1].toUpperCase()}变电所`
  const cleaned = n.replace(/(配电室|模块).*/, '').trim()
  return cleaned ? `${cleaned}变电所` : '变电所'
}

/** 是否机柜节点：仅有 type=1 子节点（遗留 buildSubstationForest 用） */
function isCabinetRawNode(node) {
  const children = node.children || []
  const type1 = children.filter(c => ((c.value || {}).type) === 1)
  const type0 = children.filter(c => ((c.value || {}).type) === 0)
  return type1.length > 0 && type0.length === 0
}

/** 变电所 zone → uuid5 配电室页（与 build_ncc_dashboard.py page_id_room 一致） */
function resolveZonePageId(modelId, zoneKey, roomSid) {
  if (!zoneKey || zoneKey === '_other') return null
  if (roomSid) return pageIdRoom(roomSid)
  const mapped = resolveXunanSubstationPageId(modelId, zoneKey)
  return mapped || null
}

/** 机柜归属变电所编码：优先机柜名(1A1_U11柜→1A1)，否则向上找 2A1模块/ECC配电室 */
function substationKeyForCabinet(cabNode, ancestorNames) {
  const cabName = cabNode.text || (cabNode.value || {}).Name || ''
  let key = substationKeyFromName(cabName)
  if (key) return key
  for (let i = ancestorNames.length - 1; i >= 0; i--) {
    const an = ancestorNames[i]
    if (String(an).toUpperCase().startsWith('ECC')) return 'ECC'
    key = substationKeyFromName(an) || zoneDisplayKey(an)
    if (key) return String(key).toUpperCase()
  }
  return '_other'
}

/** ECC 置顶，其余按编码/名称字母序（1A < 1B < 2A …） */
function isEccTreeNode(node) {
  const key = String(node.zoneKey || '').toUpperCase()
  if (key === 'ECC') return true
  return String(node.label || '').toUpperCase().startsWith('ECC')
}

function sortKeyForTreeNode(node) {
  return String(node.zoneKey || node.label || '')
}

function sortZoneSiblings(nodes) {
  return [...(nodes || [])].sort((a, b) => {
    const aEcc = isEccTreeNode(a)
    const bEcc = isEccTreeNode(b)
    if (aEcc !== bEcc) return aEcc ? -1 : 1
    return sortKeyForTreeNode(a).localeCompare(sortKeyForTreeNode(b), 'en', {
      numeric: true,
      sensitivity: 'base',
    })
  })
}

export default {
  name: 'ISMRunTreeNav',
  components: { 'ism-runtree-node': RunTreeNode },
  props: {
    projectUuid: { type: String, default: '3ec5821f-b512-2adb-3e1c-473720d0a93e' },
    modelId: { type: String, default: 'b8b4c094-faa9-a22a-1d0d-037539b27a6c' },
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
    this.$nextTick(() => {
      this.measureWidth()
    })
  },
  beforeDestroy() {
    if (this._rtTreeRO) this._rtTreeRO.disconnect()
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
      return resolveTemplatePageIdForKind(
        this.templateMap, node.kind, node.modelUuid || node.muid || '')
    },

    buildNavContext(node) {
      return buildNavContextForNode(node, this.navIndex || null)
    },

    async fetchTree() {
      this.loading = true
      try {
        const res = await getMonitorTree({ headers: { ProjectUuid: this.projectUuid } })
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
      // 默认展开：仅第一层（RootZone），配电室/机柜/设备组默认收起
      this.$nextTick(() => {
        this.expandTopLevels()
        this.measureWidth()
      })
    },

    /** 从 RootZone 子树收集机柜，按 1A1/2A2/ECC 等变电所编码分组（跳过 1A/2A 区域层） */
    buildSubstationForest(rootNode) {
      const map = new Map()

      const walk = (node, ancestors, lastRoomSid, lastRoomName) => {
        const name = node.text || (node.value || {}).Name || ''
        const isRoom = isRoomZoneName(name) || !!zoneDisplayKey(name)
        const roomSid = isRoom ? (node.value || {}).sid : lastRoomSid
        const roomName = isRoom ? name : lastRoomName
        const nextAncestors = [...ancestors, name]
        const children = node.children || []
        const type1 = children.filter(c => ((c.value || {}).type) === 1)
        const type0 = children.filter(c => ((c.value || {}).type) === 0)

        // 混合容器（如「配电室」）：type=1 子项按变电所编码分组，type=0 子项继续下钻
        if (type1.length > 0 && type0.length > 0) {
          for (const c of type1) {
            const cname = c.text || (c.value || {}).Name || ''
            let key = substationKeyFromName(cname) || substationKeyForCabinet(c, ancestors)
            if (!map.has(key)) map.set(key, { cabs: [], roomSid: null, roomName: '' })
            const entry = map.get(key)
            entry.cabs.push(c)
            if (roomSid && isRoomZoneName(roomName)) {
              entry.roomSid = roomSid
              entry.roomName = roomName
            } else if (roomSid && !entry.roomSid) {
              entry.roomSid = roomSid
              entry.roomName = roomName
            }
          }
          for (const c of type0) {
            walk(c, nextAncestors, roomSid, roomName)
          }
          return
        }

        if (isCabinetRawNode(node)) {
          const key = substationKeyForCabinet(node, ancestors)
          if (!map.has(key)) map.set(key, { cabs: [], roomSid: null, roomName: '' })
          const entry = map.get(key)
          entry.cabs.push(node)
          if (roomSid && isRoomZoneName(roomName)) {
            entry.roomSid = roomSid
            entry.roomName = roomName
          } else if (roomSid && !entry.roomSid) {
            entry.roomSid = roomSid
            entry.roomName = roomName
          }
          return
        }

        for (const c of type0) {
          walk(c, nextAncestors, roomSid, roomName)
        }
      }

      for (const child of rootNode.children || []) {
        if ((child.value || {}).type === 0) {
          walk(child, ['RootZone'], null, '')
        }
      }

      const subs = []
      for (const [key, entry] of map.entries()) {
        const cabinets = entry.cabs
          .map(n => this.transform(n))
          .filter(Boolean)
          .sort((a, b) => String(a.label).localeCompare(String(b.label), 'en', { numeric: true }))
        const label = key === 'ECC'
          ? 'ECC变电所'
          : (key === '_other' ? '其他变电所' : `${key}变电所`)
        subs.push({
          id: `sub-${key}`,
          label,
          icon: '🏛',
          kind: 'zone',
          zoneKey: key === '_other' ? null : key,
          pageId: resolveZonePageId(this.modelId, key === '_other' ? null : key, entry.roomSid),
          children: cabinets,
        })
      }
      return this.sortZoneTree(subs)
    },

    /** 合并同级「X模块」与「X配电室」（非根子树兜底） */
    mergeAliasZones(nodes) {
      if (!nodes || !nodes.length) return []
      const zoneMap = new Map()
      const rest = []
      for (const n of nodes) {
        const key = n.zoneKey
        if (key && (n.kind === 'zone' || n.kind === 'root')) {
          if (zoneMap.has(key)) {
            const exist = zoneMap.get(key)
            exist.children = this.mergeAliasZones([
              ...(exist.children || []),
              ...(n.children || []),
            ])
            if (isRoomZoneName(n.rawLabel) && n.pageId) {
              exist.pageId = n.pageId
              exist.id = n.id
            } else if (!exist.pageId && n.pageId) {
              exist.pageId = n.pageId
            }
            continue
          }
          zoneMap.set(key, {
            ...n,
            label: `${key}变电所`,
            pageId: isRoomZoneName(n.rawLabel) ? n.pageId : null,
            children: this.mergeAliasZones(n.children || []),
          })
          continue
        }
        rest.push({
          ...n,
          children: n.children && n.children.length
            ? this.mergeAliasZones(n.children)
            : n.children,
        })
      }
      return this.sortZoneTree([...zoneMap.values(), ...rest])
    },

    /** 递归排序各层子节点：ECC 第一，其余字母序 */
    sortZoneTree(nodes) {
      if (!nodes || !nodes.length) return []
      const sorted = sortZoneSiblings(nodes)
      return sorted.map(n => ({
        ...n,
        children: n.children && n.children.length
          ? this.sortZoneTree(n.children)
          : n.children,
      }))
    },

    /** 递归把 monitortree 节点转成带 page_id 的展示树（层级与设备管理 DeviceTree 一致） */
    buildForest(nodes, parentDepth = 0) {
      const normalized = normalizeRootNodes(nodes)
      const out = []
      for (const node of normalized) {
        const v = node.value || {}
        const name = node.text || v.Name || ''
        const treeDepth = parentDepth + 1
        if (name === 'RootZone' && v.sid === 1) {
          const children = (node.children || [])
            .map(c => this.transform(c, treeDepth))
            .filter(Boolean)
          out.push({
            id: node.key || 'zone-1',
            label: name,
            icon: '🏭',
            kind: 'root',
            treeDepth,
            pageId: this.modelId,
            children: this.sortTreeChildren(children),
          })
          continue
        }
        const t = this.transform(node, parentDepth)
        if (t) out.push(t)
      }
      return this.sortZoneTree(out)
    },

    sortTreeChildren(nodes) {
      return [...(nodes || [])].sort((a, b) =>
        String(a.label).localeCompare(String(b.label), 'en', { numeric: true, sensitivity: 'base' }),
      )
    },

    resolveContainerPageId(sid, pid) {
      if (pid === 1) return pageIdZone(sid)
      return pageIdRoom(sid)
    },

    transform(node, parentDepth = 0) {
      const v = node.value || {}
      const sid = v.sid
      const pid = v.pid
      const type = v.type
      const name = node.text || v.Name || '未命名'
      const rawChildren = node.children || []
      const treeDepth = parentDepth + 1

      // 设备 (type=1) = 组织层叶节点 / 信号层入口
      if (type === 1) {
        return {
          id: node.key || `dev-${sid}`,
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
          pageId: pageIdDevice(sid),
          children: [],
        }
      }

      const childNodes = []
      for (const c of rawChildren) {
        const t = this.transform(c, treeDepth)
        if (t) childNodes.push(t)
      }
      const sortedChildren = this.sortTreeChildren(childNodes)

      // 叶容器即设备（无子节点，如「数据机房报警解析」）
      if (sortedChildren.length === 0) {
        return {
          id: node.key || `dev-${sid}`,
          rawLabel: name,
          label: name,
          icon: '🔌',
          kind: 'device',
          layer: 'device',
          type: 0,
          treeDepth,
          sid,
          uuid: v.uuid || node.key || '',
          modelUuid: v.muid || '',
          muid: v.muid || '',
          status: v.Status === 1 ? 'on' : 'off',
          pageId: pageIdDevice(sid) || this.resolveContainerPageId(sid, pid),
          children: [],
        }
      }

      const childKinds = sortedChildren.map(c => c.kind)
      const kind = resolveMonitorNodeKind({ type, sid, pid, name, childKinds })
      const deviceCount = countDevicesInSubtree(sortedChildren)
      let pageId = ''
      if (kind === 'root') {
        pageId = this.modelId
      } else if (kind === 'cabinet') {
        pageId = pageIdBuilding(sid)
      } else {
        pageId = this.resolveContainerPageId(sid, pid)
      }
      return {
        id: node.key || `${kind}-${sid}`,
        rawLabel: name,
        label: name,
        icon: iconForMonitorKind(kind, pid),
        kind,
        layer: kind === 'room' ? 'room' : (kind === 'zone' ? 'zone' : kind),
        treeDepth,
        sid,
        uuid: v.uuid || '',
        count: deviceCount || undefined,
        pageId,
        children: sortedChildren,
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

    onToggle(id) {
      this.$set(this.expandedMap, id, !this.expandedMap[id])
      this.measureWidth()
    },

    async onSelect(node) {
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

      if (depth.routeMode === 'signal' || isDeviceNode(node)) {
        const muid = node.modelUuid || node.muid || ''
        const devLabel = node.label || node.rawLabel || node.name || ''
        const points = muid
          ? await fetchDeviceDatapoints(muid, devLabel, node.uuid || node.deviceUuid || '')
          : []
        const ancestors = (navContext.ancestors || [])
        navContext = buildDeviceSignalContext(node, points, ancestors)
        if (!points.length) {
          console.warn('[ISMRunTreeNav] device datapoints empty', {
            muid, devLabel, nodeId: node.id, label: node.label, rawLabel: node.rawLabel,
          })
        }
        const pageList = (store.state.ISMDisPlayEditorTool && store.state.ISMDisPlayEditorTool.PCPageList) || []
        pageUuid = resolveDeviceSignalTemplateId(this.templateMap, pageList, muid)
          || resolveOnSelectPageUuid(node, this.templateMap, this.navIndex)
          || this.resolveTemplatePageId({ kind: 'device', modelUuid: muid })
      } else if (depth.routeMode === 'childrenList') {
        const listPageSize = deviceListPageSizeForNav(navContext)
        navContext = applyGatewayListPagination({
          ...navContext,
          pageSize: listPageSize,
        })
        const pageList = (store.state.ISMDisPlayEditorTool && store.state.ISMDisPlayEditorTool.PCPageList) || []
        pageUuid = resolveDeviceListTemplateId(this.templateMap, pageList)
          || resolveOnSelectPageUuid(node, this.templateMap, this.navIndex)
      } else {
        pageUuid = resolveOnSelectPageUuid(node, this.templateMap, this.navIndex)
          || this.resolveTemplatePageId(node)
        const pageList = (store.state.ISMDisPlayEditorTool && store.state.ISMDisPlayEditorTool.PCPageList) || []
        pageUuid = pageUuid || resolveDeviceListTemplateId(this.templateMap, pageList)
      }

      if (!pageUuid) {
        if (node.children && node.children.length) this.onToggle(node.id)
        const isDevice = depth.routeMode === 'signal' || isDeviceNode(node)
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
.ism-runtree.collapsed { width: 34px; min-width: 34px; max-width: 34px; }

.rt-handle,
.rt-panel { pointer-events: auto; }

/* 收起把手 */
.rt-handle {
  width: 30px;
  height: 120px;
  margin-top: 12px;
  background: linear-gradient(180deg, #0e1a2e, #0b1322);
  border: 1px solid #1e3a5f;
  border-left: none;
  border-radius: 0 8px 8px 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 6px;
  cursor: pointer;
  color: #00e5ff;
  box-shadow: 0 0 12px rgba(0, 229, 255, 0.15);
}
.rt-handle:hover { background: linear-gradient(180deg, #11243d, #0d1a30); }
/* 把手缓慢呼吸光（3s 周期，克制） */
.rt-handle { animation: rtHandleBreath 3s ease-in-out infinite; }
@keyframes rtHandleBreath {
  0%, 100% { box-shadow: 0 0 8px rgba(0, 229, 255, 0.12); }
  50% { box-shadow: 0 0 16px rgba(0, 229, 255, 0.26); }
}
.rt-handle-icon { font-size: 18px; font-weight: 700; }
.rt-handle-text { writing-mode: vertical-rl; font-size: 12px; letter-spacing: 2px; color: #9fb6d6; }

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
