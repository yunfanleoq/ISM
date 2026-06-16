<template>
  <!-- 折叠态：仅一个竖向把手；展开态：完整树 -->
  <div class="ism-runtree" :class="{ collapsed }">
    <!-- 收起后的竖向把手 -->
    <div v-if="collapsed" class="rt-handle" @click="collapsed = false" title="展开导航树">
      <span class="rt-handle-icon">›</span>
      <span class="rt-handle-text">导航</span>
    </div>

    <!-- 展开的树面板 -->
    <div v-else class="rt-panel">
      <div class="rt-head">
        <span class="rt-title">📍 设备导航</span>
        <span class="rt-collapse-btn" @click="collapsed = true" title="收起导航树">‹</span>
      </div>

      <div class="rt-body">
        <div v-if="loading" class="rt-empty">⏳ 加载设备树中…</div>
        <div v-else-if="roots.length === 0" class="rt-empty">⚠️ 暂无设备数据</div>

        <ul v-else class="rt-tree">
          <ism-runtree-node
            v-for="node in roots"
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
      </div>
    </div>
  </div>
</template>

<script>
import { getMonitorTree } from '@/services/device'
import sha1 from 'crypto-js/sha1'
import Hex from 'crypto-js/enc-hex'
import ISMRunTreeNode from './ISMRunTreeNode.vue'

// uuid5(DNS, name).hex —— 与 build_ncc_dashboard.py 的 page_id 生成完全一致
const DNS_NS_HEX = '6ba7b8109dad11d180b400c04fd430c8'
function utf8Bytes(s) {
  return Array.from(unescape(encodeURIComponent(s))).map(c => c.charCodeAt(0))
}
function bytesToHex(b) {
  return b.map(x => x.toString(16).padStart(2, '0')).join('')
}
function hexToBytes(h) {
  const a = []
  for (let i = 0; i < h.length; i += 2) a.push(parseInt(h.substr(i, 2), 16))
  return a
}
function uuid5Hex(name) {
  const allHex = DNS_NS_HEX + bytesToHex(utf8Bytes(name))
  const hashHex = sha1(Hex.parse(allHex)).toString(Hex)
  const b = hexToBytes(hashHex).slice(0, 16)
  b[6] = (b[6] & 0x0f) | 0x50
  b[8] = (b[8] & 0x3f) | 0x80
  return bytesToHex(b)
}
const pageIdBuilding = sid => uuid5Hex(`ncc-dash-bldg-${sid}`)
const pageIdFloor = (bldgSid, key) => uuid5Hex(`ncc-dash-floor-${bldgSid}-${key}`)
const pageIdDevice = sid => uuid5Hex(`ncc-dash-dev-${sid}`)

function floorKey(name) {
  const parts = String(name || '').split('_')
  return parts.length >= 3 ? parts[2] : 'default'
}

export default {
  name: 'ISMRunTreeNav',
  components: { 'ism-runtree-node': ISMRunTreeNode },
  props: {
    // 该树仅服务航信机房 NCC 大屏
    projectUuid: { type: String, default: '31bc90be-ebc4-dd61-ba9d-ce6e075e40e2' },
    modelId: { type: String, default: '043135ad-44be-e5d8-89be-3e54883c23a8' },
  },
  data() {
    return {
      loading: true,
      roots: [],
      expandedMap: {},   // id -> bool
      selectedId: '',
      collapsed: false,
      totalDevices: 0,
    }
  },
  mounted() {
    this.fetchTree()
  },
  methods: {
    async fetchTree() {
      this.loading = true
      try {
        const res = await getMonitorTree({ headers: { ProjectUuid: this.projectUuid } })
        if (res.data && res.data.code === 0 && Array.isArray(res.data.list)) {
          this.roots = this.buildForest(res.data.list)
        }
      } catch (e) {
        console.warn('[RunTree] 获取设备树失败:', e && e.message)
      }
      this.loading = false
      // 默认展开：机房根 + 配电室（机柜默认收起，保持紧凑）
      this.$nextTick(() => this.expandTopLevels())
    },

    /** 递归把 monitortree 节点转成带 page_id 的展示树 */
    buildForest(nodes) {
      const out = []
      for (const node of nodes || []) {
        const t = this.transform(node)
        if (t) out.push(t)
      }
      return out
    },

    transform(node) {
      const v = node.value || {}
      const sid = v.Sid
      const type = v.type
      const name = node.text || v.Name || '未命名'
      const children = node.children || []

      // 设备节点 (type=1)
      if (type === 1) {
        return {
          id: node.key || `dev-${sid}`,
          label: name,
          icon: '🔌',
          kind: 'device',
          status: v.Status === 1 ? 'on' : 'off',
          pageId: pageIdDevice(sid),
          children: [],
        }
      }

      // 区域节点 (type=0)
      const type1 = children.filter(c => (c.value || {}).type === 1)
      const type0 = children.filter(c => (c.value || {}).type === 0)

      // 机柜：直接挂载设备 → 按设备名 parts[2] 分「设备组」
      if (type1.length > 0) {
        const groups = {}
        for (const c of type1) {
          const cname = c.text || (c.value || {}).Name || ''
          const key = floorKey(cname)
          if (!groups[key]) groups[key] = []
          groups[key].push(c)
        }
        const groupChildren = Object.keys(groups).sort().map(key => ({
          id: `floor-${sid}-${key}`,
          label: `${key}设备组`,
          icon: '📋',
          kind: 'group',
          count: groups[key].length,
          pageId: pageIdFloor(sid, key),
          children: groups[key].map(c => this.transform(c)),
        }))
        return {
          id: node.key || `bldg-${sid}`,
          label: name,
          icon: '🗄',
          kind: 'cabinet',
          count: type1.length,
          pageId: pageIdBuilding(sid),
          children: groupChildren,
        }
      }

      // 容器区域：机房根(Sid=1) 或 楼层/配电室 —— 继续向下递归
      const isRoot = sid === 1
      const childNodes = type0.map(c => this.transform(c)).filter(Boolean)
      return {
        id: node.key || `zone-${sid}`,
        label: name,
        icon: isRoot ? '🏭' : '🏬',
        kind: isRoot ? 'root' : 'zone',
        // 机房根钻探回总览页；中间配电室无独立 page，仅作容器
        pageId: isRoot ? this.modelId : null,
        children: childNodes,
      }
    },

    expandTopLevels() {
      const map = {}
      const walk = (nodes, depth) => {
        for (const n of nodes) {
          if (depth <= 1 && n.children && n.children.length) map[n.id] = true
          if (n.children) walk(n.children, depth + 1)
        }
      }
      walk(this.roots, 0)
      this.expandedMap = map
      // 统计设备总数
      let total = 0
      const count = nodes => nodes.forEach(n => {
        if (n.kind === 'device') total++
        if (n.children) count(n.children)
      })
      count(this.roots)
      this.totalDevices = total
    },

    onToggle(id) {
      this.$set(this.expandedMap, id, !this.expandedMap[id])
    },

    onSelect(node) {
      if (!node.pageId) {
        // 无独立页面的容器节点：点击=展开/收起
        if (node.children && node.children.length) this.onToggle(node.id)
        return
      }
      this.selectedId = node.id
      // 复用现有钻探链路：ISMRender 在 mounted 时监听 GoPage
      this.$EventBus.$emit('GoPage', {
        ModelId: this.modelId,
        PageUuid: node.pageId,
        IsPopUp: false,
        AutoClose: false,
        linkType: 'Inside',
      })
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
  width: 11.98vw;          /* 230 / 1920 */
  height: 94.81vh;         /* (1080-56) / 1080 */
  z-index: 50;
  font-family: 'Microsoft YaHei', 'PingFang SC', sans-serif;
  pointer-events: none;    /* 容器透传，子元素各自开启 */
}
.ism-runtree.collapsed { width: 34px; }

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
.rt-handle-icon { font-size: 18px; font-weight: 700; }
.rt-handle-text { writing-mode: vertical-rl; font-size: 12px; letter-spacing: 2px; color: #9fb6d6; }

/* 树面板 */
.rt-panel {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #0b1322;
  border-right: 1px solid #1e3a5f;
  box-shadow: 2px 0 16px rgba(0, 0, 0, 0.35);
}
.rt-head {
  height: 38px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 10px;
  border-bottom: 1px solid #16314f;
}
.rt-title { font-size: 13px; font-weight: 600; color: #9fb6d6; letter-spacing: 1px; }
.rt-collapse-btn {
  cursor: pointer; color: #00e5ff; font-size: 18px; font-weight: 700;
  width: 22px; height: 22px; line-height: 20px; text-align: center;
  border: 1px solid #1e3a5f; border-radius: 4px;
}
.rt-collapse-btn:hover { background: rgba(0, 229, 255, 0.12); }

.rt-body { flex: 1; overflow-y: auto; overflow-x: hidden; padding: 6px 4px; }
.rt-tree { list-style: none; margin: 0; padding: 0; }
.rt-empty { text-align: center; color: #5f7799; padding: 24px 10px; font-size: 12px; }

.rt-foot {
  height: 26px; flex-shrink: 0;
  display: flex; align-items: center; justify-content: center;
  border-top: 1px solid #16314f;
}
.rt-foot-total { font-size: 11px; color: #5f7799; }

/* 滚动条 */
.rt-body::-webkit-scrollbar { width: 6px; }
.rt-body::-webkit-scrollbar-track { background: transparent; }
.rt-body::-webkit-scrollbar-thumb { background: #1e3a5f; border-radius: 3px; }
.rt-body::-webkit-scrollbar-thumb:hover { background: #2c5a8f; }
</style>
