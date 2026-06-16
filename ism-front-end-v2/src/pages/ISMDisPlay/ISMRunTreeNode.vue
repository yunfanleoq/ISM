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
export default {
  name: 'ism-runtree-node',
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
  },
  methods: {
    onArrowClick() {
      if (this.hasChildren) this.$emit('toggle', this.node.id)
    },
    onRowClick() {
      // 行点击：可钻探的节点触发钻探；纯容器节点则展开/收起
      if (this.node.pageId) {
        this.$emit('select', this.node)
        if (this.hasChildren && !this.isExpanded) this.$emit('toggle', this.node.id)
      } else if (this.hasChildren) {
        this.$emit('toggle', this.node.id)
      }
    },
  },
}
</script>

<style scoped>
.rt-node { list-style: none; }
.rt-row {
  display: flex;
  align-items: center;
  gap: 5px;
  height: 28px;
  padding: 0 6px;
  border-radius: 5px;
  cursor: pointer;
  color: #cfe0f5;
  font-size: 12.5px;
  transition: background 0.15s, color 0.15s;
  user-select: none;
}
.rt-row:hover { background: rgba(0, 229, 255, 0.08); color: #e8f1ff; }
.rt-row.selected {
  background: rgba(0, 229, 255, 0.16);
  color: #00e5ff;
  box-shadow: inset 2px 0 0 #00e5ff;
}

.rt-arrow {
  width: 14px;
  flex-shrink: 0;
  text-align: center;
  color: #5f7799;
  font-size: 10px;
}
.rt-icon { flex-shrink: 0; font-size: 13px; }
.rt-label {
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.rt-count { flex-shrink: 0; font-size: 10px; color: #5f7799; }
.rt-dot { flex-shrink: 0; width: 7px; height: 7px; border-radius: 50%; }
.rt-dot.on { background: #10e0a0; box-shadow: 0 0 5px #10e0a0; }
.rt-dot.off { background: #5f7799; }

/* 各层级字号/颜色微调，层次更清晰 */
.kind-root > .rt-icon, .rt-row.kind-root .rt-label { font-weight: 700; }
.rt-row.kind-root { color: #00e5ff; }
.rt-row.kind-cabinet .rt-label { font-weight: 600; }

/* 子级缩进 */
.rt-children { list-style: none; margin: 0; padding-left: 14px; }
</style>
