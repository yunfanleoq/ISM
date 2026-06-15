/**
 * DataBindingPanel - data binding panel
 */
<template>
  <div class="data-binding-panel">
    <div class="panel-header">
      <span class="panel-title">{{ $t('ISM3DEditor.dataBinding') }}</span>
      <a-button
        v-if="object"
        type="primary"
        size="small"
        icon="plus"
        @click="showAddBinding = true"
      />
    </div>

    <div v-if="!object" class="empty-state">
      <a-empty :description="$t('ISM3DEditor.pleaseSelectObject')" />
    </div>

    <div v-else class="binding-content">
      <div class="binding-section">
        <div class="section-header">
          <span>{{ $t('ISM3DEditor.dataSource') }}</span>
          <a-button type="link" size="small" icon="setting" @click="showDataSourceModal = true" />
        </div>
        <div class="section-body">
          <div v-if="dataSources.length === 0" class="empty-hint">No data sources configured</div>
          <div v-else class="source-list">
            <div v-for="source in dataSources" :key="source.id" class="source-item">
              <span class="source-name">{{ source.name }}</span>
              <span class="source-type">{{ source.type }}</span>
            </div>
          </div>
        </div>
      </div>

      <div class="binding-section">
        <div class="section-header">
          <span>{{ $t('ISM3DEditor.propertyBinding') }}</span>
        </div>
        <div class="section-body">
          <div v-if="objectBindings.length === 0" class="empty-hint">No bindings</div>
          <div v-else class="binding-list">
            <div v-for="binding in objectBindings" :key="binding.id" class="binding-item">
              <div class="binding-info">
                <span class="binding-property">{{ binding.property }}</span>
                <span class="binding-arrow">-></span>
                <span class="binding-source">{{ binding.dataSource }}</span>
              </div>
              <div class="binding-actions">
                <a-switch
                  :checked="binding.enabled"
                  size="small"
                  @change="toggleBinding(binding)"
                />
                <a-icon type="delete" class="delete-btn" @click="removeBinding(binding.id)" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <a-modal
      v-model="showAddBinding"
      :title="$t('ISM3DEditor.addBinding')"
      :ok-text="$t('ISM3DEditor.add')"
      :cancel-text="$t('ISM3DEditor.cancel')"
      @ok="handleAddBinding"
    >
      <a-form :label-col="{ span: 6 }" :wrapper-col="{ span: 18 }">
        <a-form-item :label="$t('ISM3DEditor.objectProperty')">
          <a-select v-model="newBinding.property" :placeholder="$t('ISM3DEditor.selectProperty')">
            <a-select-option value="position.x">Position X</a-select-option>
            <a-select-option value="position.y">Position Y</a-select-option>
            <a-select-option value="position.z">Position Z</a-select-option>
            <a-select-option value="rotation.x">Rotation X</a-select-option>
            <a-select-option value="rotation.y">Rotation Y</a-select-option>
            <a-select-option value="rotation.z">Rotation Z</a-select-option>
            <a-select-option value="scale.x">Scale X</a-select-option>
            <a-select-option value="scale.y">Scale Y</a-select-option>
            <a-select-option value="scale.z">Scale Z</a-select-option>
            <a-select-option value="appearance.color">Color</a-select-option>
            <a-select-option value="appearance.opacity">Opacity</a-select-option>
            <a-select-option value="visible">Visible</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item :label="$t('ISM3DEditor.dataSource')">
          <a-select v-model="newBinding.dataSource" :placeholder="$t('ISM3DEditor.selectDataSource')">
            <a-select-option v-for="source in dataSources" :key="source.id" :value="source.id">
              {{ source.name }}
            </a-select-option>
            <a-select-option value="custom">{{ $t('ISM3DEditor.customData') }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item :label="$t('ISM3DEditor.dataPath')">
          <a-input v-model="newBinding.dataPath" placeholder="temperature, value" />
        </a-form-item>
        <a-form-item :label="$t('ISM3DEditor.transformType')">
          <a-select v-model="newBinding.transformType">
            <a-select-option value="direct">{{ $t('ISM3DEditor.directMapping') }}</a-select-option>
            <a-select-option value="scale">{{ $t('ISM3DEditor.scaleTransform') }}</a-select-option>
            <a-select-option value="offset">{{ $t('ISM3DEditor.offsetTransform') }}</a-select-option>
            <a-select-option value="range">{{ $t('ISM3DEditor.rangeMapping') }}</a-select-option>
            <a-select-option value="expression">{{ $t('ISM3DEditor.expressionTransform') }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item v-if="newBinding.transformType === 'scale'" :label="$t('ISM3DEditor.scaleFactor')">
          <a-input-number v-model="newBinding.scale" :min="0" :step="0.1" style="width: 100%" />
        </a-form-item>
        <a-form-item v-if="newBinding.transformType === 'offset'" :label="$t('ISM3DEditor.offsetValue')">
          <a-input-number v-model="newBinding.offset" :step="1" style="width: 100%" />
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal
      v-model="showDataSourceModal"
      :title="$t('ISM3DEditor.configureDataSource')"
      :footer="null"
      width="600px"
    >
      <div class="data-source-config">
        <a-list :data-source="dataSources" bordered>
          <template slot="renderItem" slot-scope="item">
            <a-list-item>
              <a-list-item-meta :title="item.name" :description="item.type" />
              <template slot="actions">
                <a-icon type="edit" />
                <a-icon type="delete" @click="removeDataSource(item.id)" />
              </template>
            </a-list-item>
          </template>
        </a-list>

        <a-divider>{{ $t('ISM3DEditor.addDataSource') }}</a-divider>

        <a-form :label-col="{ span: 6 }" :wrapper-col="{ span: 18 }">
          <a-form-item :label="$t('ISM3DEditor.name')">
            <a-input v-model="newSource.name" placeholder="Data source name" />
          </a-form-item>
          <a-form-item :label="$t('ISM3DEditor.type')">
            <a-select v-model="newSource.type" :placeholder="$t('ISM3DEditor.selectType')">
              <a-select-option value="websocket">{{ $t('ISM3DEditor.websocket') }}</a-select-option>
              <a-select-option value="mqtt">{{ $t('ISM3DEditor.mqtt') }}</a-select-option>
              <a-select-option value="http">{{ $t('ISM3DEditor.httpApi') }}</a-select-option>
              <a-select-option value="websocket-api">{{ $t('ISM3DEditor.websocketApi') }}</a-select-option>
              <a-select-option value="static">{{ $t('ISM3DEditor.staticData') }}</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item v-if="newSource.type !== 'static'" :label="$t('ISM3DEditor.address')">
            <a-input v-model="newSource.url" placeholder="ws:// or http://" />
          </a-form-item>
          <a-form-item :label="$t('ISM3DEditor.refreshInterval')">
            <a-input-number v-model="newSource.interval" :min="100" :step="100" />
            <span style="margin-left: 8px">ms</span>
          </a-form-item>
          <a-form-item :wrapper-col="{ span: 18, offset: 6 }">
            <a-button type="primary" @click="addDataSource">{{ $t('ISM3DEditor.add') }}</a-button>
          </a-form-item>
        </a-form>
      </div>
    </a-modal>
  </div>
</template>

<script>
import { message } from 'ant-design-vue'

export default {
  name: 'DataBindingPanel',
  i18n: require('@/i18n/language'),
  props: {
    object: {
      type: Object,
      default: null
    },
    bindings: {
      type: Array,
      default: () => []
    }
  },

  data() {
    return {
      showAddBinding: false,
      showDataSourceModal: false,
      dataSources: [],
      newBinding: {
        property: '',
        dataSource: '',
        dataPath: '',
        transformType: 'direct',
        scale: 1,
        offset: 0
      },
      newSource: {
        name: '',
        type: 'websocket',
        url: '',
        interval: 1000
      }
    }
  },

  computed: {
    objectBindings() {
      if (!this.object) return []
      return this.bindings.filter(b => b.objectId === this.object.id)
    }
  },

  methods: {
    handleAddBinding() {
      if (!this.newBinding.property || !this.newBinding.dataSource) {
        message.warning(this.$t('ISM3DEditor.fillCompleteInfo'))
        return
      }

      this.$emit('bind', {
        objectId: this.object.id,
        binding: {
          property: this.newBinding.property,
          dataSource: this.newBinding.dataSource,
          dataPath: this.newBinding.dataPath,
          transform: {
            type: this.newBinding.transformType,
            scale: this.newBinding.scale,
            offset: this.newBinding.offset
          }
        }
      })

      this.showAddBinding = false
      this.resetNewBinding()
    },

    resetNewBinding() {
      this.newBinding = {
        property: '',
        dataSource: '',
        dataPath: '',
        transformType: 'direct',
        scale: 1,
        offset: 0
      }
    },

    removeBinding(bindingId) {
      this.$emit('unbind', {
        objectId: this.object.id,
        bindingId
      })
    },

    toggleBinding(binding) {
      this.$emit('update-binding', {
        bindingId: binding.id,
        changes: { enabled: !binding.enabled }
      })
    },

    addDataSource() {
      if (!this.newSource.name) {
        message.warning(this.$t('ISM3DEditor.enterDataSourceName'))
        return
      }

      this.dataSources.push({
        id: 'source_' + Date.now(),
        ...this.newSource
      })

      message.success(this.$t('ISM3DEditor.dataSourceAdded'))
      this.resetNewSource()
    },

    removeDataSource(sourceId) {
      const index = this.dataSources.findIndex(s => s.id === sourceId)
      if (index !== -1) {
        this.dataSources.splice(index, 1)
      }
    },

    resetNewSource() {
      this.newSource = {
        name: '',
        type: 'websocket',
        url: '',
        interval: 1000
      }
    }
  }
}
</script>

<style lang="less" scoped>
.data-binding-panel {
  display: flex;
  flex-direction: column;
  height: 100%;

  .panel-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 12px;
    border-bottom: 1px solid #3c3c3c;

    .panel-title {
      font-weight: 500;
      color: #ccc;
    }
  }

  .empty-state {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .binding-content {
    flex: 1;
    overflow-y: auto;
  }

  .binding-section {
    border-bottom: 1px solid #3c3c3c;

    .section-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 8px 12px;
      background: #2d2d2d;
      color: #ccc;
      font-size: 12px;
      font-weight: 500;
    }

    .section-body {
      padding: 8px 12px;

      .empty-hint {
        color: #666;
        font-size: 12px;
        text-align: center;
        padding: 12px;
      }
    }
  }

  .source-list {
    .source-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 6px 8px;
      background: #3c3c3c;
      border-radius: 4px;
      margin-bottom: 4px;

      .source-name {
        color: #ccc;
        font-size: 12px;
      }

      .source-type {
        color: #666;
        font-size: 10px;
        text-transform: uppercase;
      }
    }
  }

  .binding-list {
    .binding-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 8px;
      background: #3c3c3c;
      border-radius: 4px;
      margin-bottom: 6px;

      .binding-info {
        display: flex;
        align-items: center;
        gap: 6px;

        .binding-property {
          color: #4a90d9;
          font-size: 11px;
          font-family: monospace;
        }

        .binding-arrow {
          color: #666;
        }

        .binding-source {
          color: #2ecc71;
          font-size: 11px;
          font-family: monospace;
        }
      }

      .binding-actions {
        display: flex;
        align-items: center;
        gap: 8px;

        .delete-btn {
          color: #e74c3c;
          cursor: pointer;

          &:hover {
            color: #ff6b6b;
          }
        }
      }
    }
  }

  .data-source-config {
    .ant-list-item {
      background: #2d2d2d;
    }
  }

  :deep(.ant-form-item-label > label) {
    color: #888;
  }

  :deep(.ant-input),
  :deep(.ant-select-selection) {
    background: #3c3c3c;
    border-color: #4a4a4a;
    color: #ccc;
  }

  :deep(.ant-btn-primary) {
    background: #4a90d9;
    border-color: #4a90d9;

    &:hover {
      background: #5a9ee0;
    }
  }
}
</style>
