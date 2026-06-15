<template>
  <div class="context-menu" v-if="visible" :style="{top: y+'px', left: x+'px'}" @mouseleave="$emit('close')">
    <div class="context-menu-item" @click="handleAction('duplicate')"><i class="fas fa-copy"></i> {{ $t('ISM3DEditor.ContextMenuCopy') }}</div>
    <div class="context-menu-item" @click="handleAction('focus')"><i class="fas fa-crosshairs"></i> {{ $t('ISM3DEditor.ContextMenuFocus') }}</div>
    <div class="context-menu-sep"></div>
    <div class="context-menu-item" @click="handleAction('reset-transform')"><i class="fas fa-undo"></i> {{ $t('ISM3DEditor.ContextMenuResetTransform') }}</div>
    <div class="context-menu-sep"></div>
    <div class="context-menu-item danger" @click="handleAction('delete')"><i class="fas fa-trash"></i> {{ $t('ISM3DEditor.ContextMenuDelete') }}</div>
  </div>
</template>

<script>
export default {
  name: 'ContextMenu',
  props: {
    visible: { type: Boolean, default: false },
    x: { type: Number, default: 0 },
    y: { type: Number, default: 0 },
  },
  methods: {
    handleAction(action) {
      this.$emit('action', action)
      this.$emit('close')
    }
  }
}
</script>

<style scoped>
.context-menu {
  position: fixed;
  background: #fff;
  border: 1px solid #e8e8e8;
  border-radius: 6px;
  padding: 4px 0;
  z-index: 9999;
  min-width: 140px;
  box-shadow: 0 8px 24px rgba(0,0,0,.15);
}
.context-menu-item {
  padding: 6px 14px;
  font-size: 12px;
  color: #333;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
}
.context-menu-item:hover { background: #e6fffb; color: #13c2c2; }
.context-menu-item.danger { color: #ff4d4f; }
.context-menu-item.danger:hover { background: #fff1f0; }
.context-menu-sep { height: 1px; background: #e8e8e8; margin: 4px 0; }
</style>
