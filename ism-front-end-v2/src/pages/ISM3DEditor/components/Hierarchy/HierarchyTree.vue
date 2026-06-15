/**
 * HierarchyTree - 层级树组件
 */
<template>
  <div class="hierarchy-tree">
    <div class="tree-header">
      <span class="tree-title">{{ $t('ISM3DEditor.hierarchy') }}</span>
      <div class="tree-actions">
        <a-tooltip :title="$t('ISM3DEditor.expandAll')">
          <a-button icon="down" size="small" type="text" @click="expandAll" />
        </a-tooltip>
        <a-tooltip :title="$t('ISM3DEditor.collapseAll')">
          <a-button icon="up" size="small" type="text" @click="collapseAll" />
        </a-tooltip>
      </div>
    </div>

    <div class="tree-content" @dragover.prevent @drop="handleDrop">
      <div v-if="objects.length === 0" class="empty-hint">
        {{ $t('ISM3DEditor.sceneEmpty') }}
      </div>

      <div v-else class="tree-list">
        <div
          v-for="obj in objects"
          :key="obj.id"
          class="tree-item"
          :class="{ selected: obj.id === selectedId }"
          @click="handleSelect(obj.id)"
          @contextmenu.prevent="showContextMenu($event, obj)"
          draggable="true"
          @dragstart="handleDragStart($event, obj)"
        >
          <span class="item-icon" v-html="getObjectIcon(obj.type)"></span>
          <span class="item-name">{{ obj.name || obj.type }}</span>
          <span class="item-type">{{ obj.type }}</span>

          <div class="item-actions">
            <a-tooltip :title="$t('ISM3DEditor.hideShow')">
              <a-icon
                :type="obj.visible ? 'eye' : 'eye-invisible'"
                @click.stop="toggleVisibility(obj)"
              />
            </a-tooltip>
            <a-tooltip :title="$t('ISM3DEditor.lockUnlock')">
              <a-icon
                :type="obj.locked ? 'lock' : 'unlock'"
                @click.stop="toggleLock(obj)"
              />
            </a-tooltip>
          </div>
        </div>
      </div>
    </div>

    <!-- 右键菜单 -->
    <div
      v-if="contextMenu.visible"
      class="context-menu"
      :style="{ top: contextMenu.y + 'px', left: contextMenu.x + 'px' }"
      @click="hideContextMenu"
    >
      <div class="menu-item" @click="handleContextAction('duplicate')">
        <a-icon type="copy" /> {{ $t('ISM3DEditor.duplicate') }}
      </div>
      <div class="menu-item" @click="handleContextAction('delete')">
        <a-icon type="delete" /> {{ $t('ISM3DEditor.delete') }}
      </div>
      <div class="menu-divider"></div>
      <div class="menu-item" @click="handleContextAction('focus')">
        <a-icon type="aim" /> {{ $t('ISM3DEditor.focusObject') }}
      </div>
      <div class="menu-item" @click="handleContextAction('duplicate')">
        <a-icon type="block" /> {{ $t('ISM3DEditor.hide') }}
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'HierarchyTree',
  i18n: require('@/i18n/language'),
  props: {
    objects: {
      type: Array,
      default: () => []
    },
    selectedId: {
      type: String,
      default: null
    }
  },

  data() {
    return {
      contextMenu: {
        visible: false,
        x: 0,
        y: 0,
        target: null
      }
    }
  },

  methods: {
    handleSelect(id) {
      this.$emit('select', id)
    },

    handleDragStart(event, obj) {
      event.dataTransfer.setData('application/json', JSON.stringify({
        type: 'move',
        objectId: obj.id
      }))
    },

    handleDrop(event) {
      const data = JSON.parse(event.dataTransfer.getData('application/json') || '{}')
      if (data.action === 'add') {
        // 从工具箱拖放
        this.$emit('add', { type: data.type })
      }
    },

    toggleVisibility(obj) {
      this.$emit('update', {
        id: obj.id,
        changes: { visible: !obj.visible }
      })
    },

    toggleLock(obj) {
      this.$emit('update', {
        id: obj.id,
        changes: { locked: !obj.locked }
      })
    },

    showContextMenu(event, obj) {
      this.contextMenu = {
        visible: true,
        x: event.clientX,
        y: event.clientY,
        target: obj
      }
    },

    hideContextMenu() {
      this.contextMenu.visible = false
    },

    handleContextAction(action) {
      const obj = this.contextMenu.target
      if (!obj) return

      switch (action) {
        case 'duplicate':
          this.$emit('duplicate', obj.id)
          break
        case 'delete':
          this.$emit('delete', obj.id)
          break
        case 'focus':
          this.$emit('focus', obj.id)
          break
      }

      this.hideContextMenu()
    },

    expandAll() {
      // 暂时不需要
    },

    collapseAll() {
      // 暂时不需要
    },

    getObjectIcon(type) {
      const icons = {
        box: '<svg viewBox="0 0 24 24"><path fill="currentColor" d="M21 16.5L12 21.5L3 16.5V7.5L12 2.5L21 7.5V16.5Z"/></svg>',
        sphere: '<svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="8" fill="none" stroke="currentColor" stroke-width="2"/></svg>',
        cylinder: '<svg viewBox="0 0 24 24"><ellipse cx="12" cy="6" rx="6" ry="2" fill="none" stroke="currentColor" stroke-width="2"/><line x1="6" y1="6" x2="6" y2="18" stroke="currentColor" stroke-width="2"/><line x1="18" y1="6" x2="18" y2="18" stroke="currentColor" stroke-width="2"/></svg>',
        tank: '<svg viewBox="0 0 24 24"><path d="M6 18V8C6 6 8 4 12 4C16 4 18 6 18 8V18" fill="none" stroke="currentColor" stroke-width="2"/></svg>',
        pump: '<svg viewBox="0 0 24 24"><rect x="4" y="8" width="10" height="8" rx="2" fill="none" stroke="currentColor" stroke-width="2"/></svg>',
        valve: '<svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="6" fill="none" stroke="currentColor" stroke-width="2"/></svg>',
        conveyor: '<svg viewBox="0 0 24 24"><rect x="2" y="10" width="20" height="4" fill="none" stroke="currentColor" stroke-width="2"/></svg>',
        motor: '<svg viewBox="0 0 24 24"><rect x="2" y="6" width="14" height="12" rx="2" fill="none" stroke="currentColor" stroke-width="2"/></svg>',
        dataPanel: '<svg viewBox="0 0 24 24"><rect x="2" y="4" width="20" height="16" rx="2" fill="none" stroke="currentColor" stroke-width="2"/></svg>',
        indicator: '<svg viewBox="0 0 24 24"><circle cx="12" cy="10" r="6" fill="none" stroke="currentColor" stroke-width="2"/></svg>',
        gauge: '<svg viewBox="0 0 24 24"><path d="M4 18A10 10 0 1 1 20 12" fill="none" stroke="currentColor" stroke-width="2"/></svg>',
        wall: '<svg viewBox="0 0 24 24"><rect x="2" y="4" width="20" height="16" fill="none" stroke="currentColor" stroke-width="2"/></svg>',
        floor: '<svg viewBox="0 0 24 24"><rect x="2" y="10" width="20" height="4" fill="none" stroke="currentColor" stroke-width="2"/></svg>',
        default: '<svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="6" fill="none" stroke="currentColor" stroke-width="2"/></svg>'
      }
      return icons[type] || icons.default
    }
  }
}
</script>

<style lang="less" scoped>
.hierarchy-tree {
  display: flex;
  flex-direction: column;
  height: 100%;

  .tree-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 12px;
    border-bottom: 1px solid #e8e8e8;

    .tree-title {
      font-weight: 500;
      color: #333;
    }

    .tree-actions {
      display: flex;
      gap: 4px;
    }
  }

  .tree-content {
    flex: 1;
    overflow-y: auto;
    padding: 4px 0;

    .empty-hint {
      padding: 20px;
      text-align: center;
      color: #666;
      font-size: 12px;
    }

    .tree-list {
      .tree-item {
        display: flex;
        align-items: center;
        padding: 6px 12px;
        cursor: pointer;
        transition: background 0.15s;
        user-select: none;

        &:hover {
          background: #f0f0f0;

          .item-actions {
            opacity: 1;
          }
        }

        &.selected {
          background: #e6fffb;

          .item-name,
          .item-type {
            color: #13c2c2;
          }

          .item-icon {
            color: #13c2c2;
          }
        }

        .item-icon {
          width: 16px;
          height: 16px;
          margin-right: 8px;
          color: #999;
          flex-shrink: 0;

          svg {
            width: 100%;
            height: 100%;
          }
        }

        .item-name {
          flex: 1;
          color: #333;
          font-size: 12px;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }

        .item-type {
          color: #999;
          font-size: 10px;
          margin-left: 8px;
        }

        .item-actions {
          display: flex;
          gap: 4px;
          opacity: 0;
          transition: opacity 0.15s;

          .anticon {
            font-size: 12px;
            color: #999;
            cursor: pointer;
            padding: 2px;

            &:hover {
              color: #13c2c2;
            }
          }
        }
      }
    }
  }

  .context-menu {
    position: fixed;
    background: #ffffff;
    border: 1px solid #d9d9d9;
    border-radius: 4px;
    padding: 4px 0;
    z-index: 1000;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    min-width: 120px;

    .menu-item {
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 8px 12px;
      color: #333;
      font-size: 12px;
      cursor: pointer;

      &:hover {
        background: #e6fffb;
        color: #13c2c2;
      }

      .anticon {
        font-size: 14px;
      }
    }

    .menu-divider {
      height: 1px;
      background: #e8e8e8;
      margin: 4px 0;
    }
  }

  :deep(.ant-btn-text) {
    color: #666;

    &:hover {
      color: #13c2c2;
      background: transparent;
    }
  }
}
</style>
