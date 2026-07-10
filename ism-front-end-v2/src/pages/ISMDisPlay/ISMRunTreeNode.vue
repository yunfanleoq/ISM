<template>
  <li class="rt-node">
    <div
      class="rt-row"
      :class="[`kind-${node.kind}`, { selected: selectedId === node.id }]"
      @click="onRowClick"
    >
      <!-- 展开/收起箭头 -->
      <span class="rt-arrow" @click.stop="onArrowClick">
        <template v-if="hasChildren">{{ isExpanded ? '▾' : '▸' }}</template>
        <template v-else> </template>
      </span>
      <span class="rt-icon">{{ node.icon }}</span>
      <span class="rt-label" :title="node.label">{{ node.label }}</span>
      <span v-if="node.kind === 'device'" class="rt-dot" :class="node.status"></span>
      <span v-else-if="node.count != null" class="rt-count">{{ node.count }}台</span>
    </div>

    <ul v-if="hasChildren" v-show="isExpanded" class="rt-children">
      <ism-runtree-node
        v-for="child in node.children"
        :key="child.id"
        :node="child"
        :expanded-map="expandedMap"
        :selected-id="selectedId"
        @toggle="$emit('toggle', $event)"
        @select="$emit('select', $event)"
      />
    </ul>
  </li>
</template>

<script>
// 递归树节点：异步自引用，避免同步循环依赖导致模块 export 错乱
const RunTreeNodeSelf = () => import(/* webpackChunkName: "ism-runtree-node" */ './ISMRunTreeNode.vue')

export default {
  name: 'IsmRuntreeNode',
  components: { 'ism-runtree-node': RunTreeNodeSelf },
  props: {
    node: { type: Object, required: true },
    expandedMap: { type: Object, required: true },
    selectedId: { type: String, default: '' },
  },
  computed: {
    hasChildren() {
      return Array.isArray(this.node.children) && this.node.children.length > 0
    },
    isExpanded() {
      return !!this.expandedMap[this.node.id]
    },
    /** 仅真实设备叶节点可钻探；组织容器（含 UPS/配电室 等目录）不可 */
    isNavigableDevice() {
      const n = this.node
      if (!n) return false
      if (n.type === 1) return true
      if (n.kind === 'device' || n.kind === 'gateway') return true
      if (n.layer === 'gateway') return true
      return false
    },
  },
  methods: {
    onArrowClick() {
      if (this.hasChildren) this.$emit('toggle', this.node.id)
    },
    onRowClick() {
      // 非叶子：只展开/收起目录；叶子设备：才进入详情页
      if (this.hasChildren) {
        this.$emit('toggle', this.node.id)
        return
      }
      if (this.isNavigableDevice) {
        this.$emit('select', this.node)
      }
    },
  },
}
</script>

<style scoped>
.rt-node { list-style: none; }

/* 节点行：科技感卡片底（深色半透明 + 细边框 + 圆角），行间留白即为分隔 */
.rt-row {
  position: relative;
  display: flex;
  align-items: center;
  gap: 6px;
  min-height: 25px;
  padding: 2px 8px 2px 5px;
  margin-bottom: 3px;
  border-radius: 6px;
  cursor: pointer;
  color: #cfe0f5;
  font-size: 12px;
  line-height: 1.25;
  background: rgba(13, 26, 46, 0.42);
  border: 1px solid rgba(30, 58, 95, 0.35);
  transition: background 0.18s, border-color 0.18s, box-shadow 0.18s, color 0.18s;
  user-select: none;
  overflow: visible;
}
/* 左侧高亮条：默认隐藏，hover/选中点亮（克制的发光点缀） */
.rt-row::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 2px;
  height: 0;
  background: linear-gradient(180deg, transparent, #00e5ff, transparent);
  border-radius: 2px;
  transition: height 0.2s ease;
}
.rt-row:hover {
  background: rgba(0, 229, 255, 0.07);
  border-color: rgba(0, 229, 255, 0.3);
  color: #e8f1ff;
}
.rt-row:hover::before { height: 58%; }
.rt-row.selected {
  background: rgba(0, 229, 255, 0.13);
  border-color: rgba(0, 229, 255, 0.5);
  color: #00e5ff;
  box-shadow: 0 0 9px rgba(0, 229, 255, 0.16), inset 0 0 7px rgba(0, 229, 255, 0.07);
}
.rt-row.selected::before { height: 76%; box-shadow: 0 0 6px #00e5ff; }

.rt-arrow {
  width: 12px;
  flex-shrink: 0;
  text-align: center;
  color: #5f7799;
  font-size: 9px;
}
.rt-icon { flex-shrink: 0; font-size: 12px; line-height: 1; }
.rt-label {
  flex: 0 1 auto;
  min-width: max-content;
  white-space: nowrap;
}
.rt-count {
  flex-shrink: 0;
  font-size: 9.5px;
  color: #6b86ab;
  padding: 0 5px;
  line-height: 14px;
  border: 1px solid rgba(30, 58, 95, 0.6);
  border-radius: 8px;
  background: rgba(13, 26, 46, 0.5);
}
.rt-dot { flex-shrink: 0; width: 7px; height: 7px; border-radius: 50%; }
/* 在线点：缓慢呼吸（3s 周期，低频不刺眼） */
.rt-dot.on {
  background: #10e0a0;
  box-shadow: 0 0 6px #10e0a0;
  animation: rtDotPulse 3s ease-in-out infinite;
}
.rt-dot.off { background: #5f7799; }
@keyframes rtDotPulse {
  0%, 100% { opacity: 0.6; }
  50% { opacity: 1; }
}

/* 各层级字号/配色：层次清晰，逐级递减 */
.rt-row.kind-root {
  font-size: 13px;
  font-weight: 700;
  color: #00e5ff;
  background: rgba(0, 229, 255, 0.06);
  border-color: rgba(0, 229, 255, 0.28);
}
.rt-row.kind-root .rt-label { text-shadow: 0 0 6px rgba(0, 229, 255, 0.4); }
.rt-row.kind-cabinet { font-size: 12.5px; }
.rt-row.kind-cabinet .rt-label { font-weight: 600; }
.rt-row.kind-zone { font-size: 12px; }
.rt-row.kind-group { font-size: 11.5px; color: #b9cce6; }
.rt-row.kind-device { font-size: 11px; }

/* 子级缩进 + 竖向连接导引线（兼具结构分隔） */
.rt-children {
  list-style: none;
  margin: 0 0 0 7px;
  padding-left: 8px;
  border-left: 1px solid rgba(30, 58, 95, 0.5);
}
</style>
