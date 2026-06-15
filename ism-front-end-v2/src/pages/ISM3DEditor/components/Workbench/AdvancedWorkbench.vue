<template>
  <div class="advanced-workbench">
    <a-tabs v-model="activeTab" size="small" :animated="false">
      <a-tab-pane v-if="false" key="assets" :tab="$t('ISM3DEditor.resources')">
        <div class="workbench-toolbar">
          <a-input-search v-model="assetKeyword" size="small" :placeholder="$t('ISM3DEditor.searchModelObject')" />
          <a-button size="small" icon="plus" @click="createAssetFromSelection" :disabled="!selectedObject">{{ $t('ISM3DEditor.favoriteSelected') }}</a-button>
        </div>
        <div class="inline-form">
          <span>{{ $t('ISM3DEditor.modelUrl') }}</span>
          <a-input v-model="modelUrl" size="small" :placeholder="$t('ISM3DEditor.glbGltfUrl')" />
          <a-button size="small" @click="addExternalModelAsset" :disabled="!modelUrl">{{ $t('ISM3DEditor.add') }}</a-button>
        </div>
        <div class="asset-filters">
          <a-tag :color="assetCategory === 'all' ? 'green' : ''" @click="assetCategory = 'all'">{{ $t('ISM3DEditor.all') }}</a-tag>
          <a-tag :color="assetCategory === 'model' ? 'green' : ''" @click="assetCategory = 'model'">{{ $t('ISM3DEditor.model') }}</a-tag>
          <a-tag :color="assetCategory === 'object' ? 'green' : ''" @click="assetCategory = 'object'">{{ $t('ISM3DEditor.object') }}</a-tag>
          <a-tag :color="assetCategory === 'effect' ? 'green' : ''" @click="assetCategory = 'effect'">{{ $t('ISM3DEditor.effect') }}</a-tag>
          <a-tag :color="assetCategory === 'favorite' ? 'green' : ''" @click="assetCategory = 'favorite'">{{ $t('ISM3DEditor.favorite') }}</a-tag>
        </div>
        <div class="asset-filters asset-source-filters">
          <a-tag :color="assetSource === 'all' ? 'blue' : ''" @click="assetSource = 'all'">{{ $t('ISM3DEditor.allSources') }}</a-tag>
          <a-tag :color="assetSource === 'builtin' ? 'blue' : ''" @click="assetSource = 'builtin'">{{ $t('ISM3DEditor.builtIn') }}</a-tag>
          <a-tag :color="assetSource === 'catalog' ? 'blue' : ''" @click="assetSource = 'catalog'">{{ $t('ISM3DEditor.catalog') }}</a-tag>
          <a-tag :color="assetSource === 'custom' ? 'blue' : ''" @click="assetSource = 'custom'">{{ $t('ISM3DEditor.favorite') }}</a-tag>
          <a-tag :color="assetSource === 'external' ? 'blue' : ''" @click="assetSource = 'external'">{{ $t('ISM3DEditor.externalModel') }}</a-tag>
        </div>
        <div class="asset-grid">
          <div
            v-for="asset in filteredAssets"
            :key="asset.id"
            class="asset-card"
            draggable="true"
            @dragstart="handleAssetDragStart($event, asset)"
            @click="addAsset(asset)"
          >
            <div class="asset-thumb">
              <img v-if="asset.thumbnail" :src="asset.thumbnail" />
              <i v-else :class="asset.icon || 'fas fa-cube'"></i>
            </div>
            <div class="asset-name" :title="asset.name">{{ asset.name }}</div>
            <span v-if="asset.license" class="asset-license">{{ asset.license }}</span>
            <span class="asset-source">{{ assetSourceLabel(asset) }}</span>
            <a-icon
              v-if="isCustomAsset(asset)"
              class="asset-remove"
              type="close"
              @click.stop="removeAsset(asset.id)"
            />
          </div>
        </div>
      </a-tab-pane>

      <a-tab-pane v-if="false" key="layers" :tab="$t('ISM3DEditor.layers')">
        <div class="workbench-toolbar">
          <a-input-search v-model="layerKeyword" size="small" :placeholder="$t('ISM3DEditor.searchObject')" />
          <a-button size="small" icon="folder-add" @click="assignGroup" :disabled="!selectedObject">{{ $t('ISM3DEditor.group') }}</a-button>
        </div>
        <div class="inline-form">
          <span>{{ $t('ISM3DEditor.group') }}</span>
          <a-select v-model="groupFilter" size="small">
            <a-select-option value="">{{ $t('ISM3DEditor.all') }}</a-select-option>
            <a-select-option v-for="group in groupOptions" :key="group" :value="group">{{ group }}</a-select-option>
          </a-select>
          <a-input v-model="groupNameDraft" size="small" :placeholder="$t('ISM3DEditor.groupName')" />
        </div>
        <div class="layer-group-actions">
          <a-button size="small" @click="assignGroupToFiltered" :disabled="!filteredObjectIds.length">{{ $t('ISM3DEditor.batchGroup') }}</a-button>
          <a-button size="small" @click="clearGroupForFiltered" :disabled="!filteredObjectIds.length">{{ $t('ISM3DEditor.clearGroup') }}</a-button>
          <a-button size="small" @click="groupFilter = ''" :disabled="!groupFilter">{{ $t('ISM3DEditor.allGroups') }}</a-button>
        </div>
        <div class="layer-actions">
          <a-button size="small" @click="$emit('batch-layer', { action: 'showAll', ids: filteredObjectIds })">{{ $t('ISM3DEditor.showAll') }}</a-button>
          <a-button size="small" @click="$emit('batch-layer', { action: 'hideAll', ids: filteredObjectIds })">{{ $t('ISM3DEditor.hideAll') }}</a-button>
          <a-button size="small" @click="$emit('batch-layer', { action: 'lockAll', ids: filteredObjectIds })">{{ $t('ISM3DEditor.lockAll') }}</a-button>
          <a-button size="small" @click="$emit('batch-layer', { action: 'unlockAll', ids: filteredObjectIds })">{{ $t('ISM3DEditor.unlock') }}</a-button>
          <a-button size="small" @click="$emit('duplicate-object', selectedId)" :disabled="!selectedObject">{{ $t('ISM3DEditor.duplicate') }}</a-button>
          <a-button size="small" type="danger" @click="$emit('delete-object', selectedId)" :disabled="!selectedObject">{{ $t('ISM3DEditor.delete') }}</a-button>
        </div>
        <div class="layer-list">
          <div
            v-for="obj in filteredObjects"
            :key="obj.id"
            class="layer-row"
            :class="{ selected: obj.id === selectedId }"
            @click="$emit('select-object', obj.id)"
          >
            <i :class="obj.icon || iconForObject(obj)"></i>
            <span class="layer-name" :title="obj.name">{{ obj.name }}</span>
            <span v-if="obj.groupName" class="layer-group" :title="obj.groupName">{{ obj.groupName }}</span>
            <a-tooltip :title="$t('ISM3DEditor.hideShow')">
              <a-icon :type="obj.visible === false ? 'eye-invisible' : 'eye'" @click.stop="patchObject(obj.id, { visible: obj.visible === false })" />
            </a-tooltip>
            <a-tooltip :title="$t('ISM3DEditor.lockUnlock')">
              <a-icon :type="obj.locked ? 'lock' : 'unlock'" @click.stop="patchObject(obj.id, { locked: !obj.locked })" />
            </a-tooltip>
            <a-tooltip :title="$t('ISM3DEditor.moveUp')">
              <a-icon type="arrow-up" @click.stop="$emit('move-object', { id: obj.id, direction: -1 })" />
            </a-tooltip>
            <a-tooltip :title="$t('ISM3DEditor.moveDown')">
              <a-icon type="arrow-down" @click.stop="$emit('move-object', { id: obj.id, direction: 1 })" />
            </a-tooltip>
          </div>
        </div>
      </a-tab-pane>

      <a-tab-pane v-if="false" key="materials" :tab="$t('ISM3DEditor.material')">
        <div class="preset-grid">
          <div
            v-for="preset in materialPresets"
            :key="preset.id"
            class="preset-card"
            :class="{ disabled: !selectedObject }"
            @click="applyMaterial(preset)"
          >
            <span class="material-swatch" :style="{ background: preset.color }"></span>
            <div>
              <div class="preset-title">{{ preset.name }}</div>
              <div class="preset-desc">M {{ preset.metalness }} / R {{ preset.roughness }}</div>
            </div>
          </div>
        </div>
        <a-divider style="margin:8px 0" />
        <div class="inline-form">
          <span>{{ $t('ISM3DEditor.textureUrl') }}</span>
          <a-input v-model="textureUrl" size="small" :placeholder="$t('ISM3DEditor.textureUrlPlaceholder')" />
          <a-button size="small" @click="applyTexture" :disabled="!selectedObject">{{ $t('ISM3DEditor.apply') }}</a-button>
        </div>
      </a-tab-pane>

      <a-tab-pane v-if="false" key="lighting" :tab="$t('ISM3DEditor.lighting')">
        <div class="button-grid">
          <a-button size="small" @click="addLight('point_light')">{{ $t('ISM3DEditor.pointLight') }}</a-button>
          <a-button size="small" @click="addLight('spot_light')">{{ $t('ISM3DEditor.spotLight') }}</a-button>
          <a-button size="small" @click="addLight('ambient_light')">{{ $t('ISM3DEditor.ambientLight') }}</a-button>
          <a-button size="small" @click="addLight('directional_light')">{{ $t('ISM3DEditor.directionalLight') }}</a-button>
        </div>
        <template v-if="selectedLightObject">
          <a-divider style="margin:8px 0" />
          <div class="inline-form">
            <span>{{ $t('ISM3DEditor.color') }}</span>
            <input class="workbench-color-input" type="color" :value="selectedLightObject.color || '#ffffff'" @input="patchObject(selectedLightObject.id, { color: $event.target.value })" />
            <span class="form-value">{{ selectedLightObject.color || '#ffffff' }}</span>
          </div>
          <div class="inline-form">
            <span>{{ $t('ISM3DEditor.intensity') }}</span>
            <a-slider :value="selectedLightObject.intensity === undefined ? 1 : selectedLightObject.intensity" :min="0" :max="10" :step="0.1" @change="patchObject(selectedLightObject.id, { intensity: $event })" />
            <a-input-number :value="selectedLightObject.intensity === undefined ? 1 : selectedLightObject.intensity" size="small" :min="0" :max="10" :step="0.1" @change="patchObject(selectedLightObject.id, { intensity: $event })" />
          </div>
          <div v-if="selectedLightObject.type === 'point_light' || selectedLightObject.type === 'spot_light'" class="inline-form">
            <span>{{ $t('ISM3DEditor.distance') }}</span>
            <a-slider :value="selectedLightObject.distance === undefined ? 10 : selectedLightObject.distance" :min="0" :max="80" :step="1" @change="patchObject(selectedLightObject.id, { distance: $event })" />
            <a-input-number :value="selectedLightObject.distance === undefined ? 10 : selectedLightObject.distance" size="small" :min="0" :max="80" :step="1" @change="patchObject(selectedLightObject.id, { distance: $event })" />
          </div>
          <div v-if="selectedLightObject.type === 'spot_light'" class="inline-form">
            <span>{{ $t('ISM3DEditor.angle') }}</span>
            <a-slider :value="selectedLightAngle" :min="5" :max="90" :step="1" @change="updateLightAngle" />
            <a-input-number :value="selectedLightAngle" size="small" :min="5" :max="90" :step="1" @change="updateLightAngle" />
          </div>
          <div v-if="selectedLightObject.type === 'spot_light'" class="inline-form">
            <span>{{ $t('ISM3DEditor.penumbra') }}</span>
            <a-slider :value="selectedLightObject.penumbra === undefined ? 0.25 : selectedLightObject.penumbra" :min="0" :max="1" :step="0.05" @change="patchObject(selectedLightObject.id, { penumbra: $event })" />
            <a-input-number :value="selectedLightObject.penumbra === undefined ? 0.25 : selectedLightObject.penumbra" size="small" :min="0" :max="1" :step="0.05" @change="patchObject(selectedLightObject.id, { penumbra: $event })" />
          </div>
        </template>
        <a-divider style="margin:8px 0" />
        <div class="inline-form">
          <span>环境</span>
          <a-select :value="sceneExtras.environment" size="small" @change="updateSceneExtra('environment', $event)">
            <a-select-option value="studio">Studio</a-select-option>
            <a-select-option value="industrial">Industrial</a-select-option>
            <a-select-option value="night">Night</a-select-option>
            <a-select-option value="outdoor">Outdoor</a-select-option>
          </a-select>
        </div>
        <div class="inline-form">
          <span>雾效</span>
          <a-switch size="small" :checked="!!sceneExtras.fogEnabled" @change="updateSceneExtra('fogEnabled', $event)" />
          <input type="color" :value="sceneExtras.fogColor || '#d8e4f0'" @input="updateSceneExtra('fogColor', $event.target.value)" />
        </div>
      </a-tab-pane>

      <a-tab-pane v-if="false" key="camera" tab="相机">
        <div class="workbench-toolbar">
          <a-button size="small" icon="camera" @click="saveCameraView">保存视角</a-button>
          <a-button size="small" icon="play-circle" @click="setDefaultCamera" :disabled="!selectedCameraId">设为默认</a-button>
          <a-button size="small" icon="sync" @click="$emit('update-camera', selectedCameraId)" :disabled="!selectedCameraId">更新</a-button>
        </div>
        <div class="camera-list">
          <div
            v-for="camera in sceneExtras.cameraViews"
            :key="camera.id"
            class="camera-row"
            :class="{ selected: camera.id === selectedCameraId }"
            @click="selectedCameraId = camera.id"
          >
            <i class="fas fa-video"></i>
            <a-input
              class="camera-name-input"
              :value="camera.name"
              size="small"
              @click.stop
              @change="renameCamera(camera.id, $event.target.value)"
            />
            <a-tag v-if="sceneExtras.defaultCameraId === camera.id" color="blue">默认</a-tag>
            <a-icon type="select" @click.stop="$emit('apply-camera', camera)" />
            <a-icon type="sync" @click.stop="$emit('update-camera', camera.id)" />
            <a-icon type="delete" @click.stop="removeCamera(camera.id)" />
          </div>
        </div>
      </a-tab-pane>

      <a-tab-pane v-if="false" key="events" tab="事件">
        <div v-if="!selectedObject" class="empty-tip">选择对象后配置交互事件</div>
        <template v-else>
          <div class="inline-form">
            <span>触发</span>
            <a-select v-model="eventDraft.trigger" size="small">
              <a-select-option value="click">点击</a-select-option>
              <a-select-option value="hover">悬停</a-select-option>
              <a-select-option value="dblclick">双击</a-select-option>
              <a-select-option value="data">数据条件</a-select-option>
            </a-select>
          </div>
          <div class="inline-form">
            <span>动作</span>
            <a-select v-model="eventDraft.action" size="small">
              <a-select-option value="focus">聚焦对象</a-select-option>
              <a-select-option value="camera">切换相机</a-select-option>
              <a-select-option value="toggleVisible">切换显隐</a-select-option>
              <a-select-option value="openPanel">打开弹窗</a-select-option>
              <a-select-option value="callApi">调用接口</a-select-option>
            </a-select>
          </div>
          <div class="inline-form">
            <span>参数</span>
            <a-select v-if="eventDraft.action === 'camera'" v-model="eventDraft.value" size="small" placeholder="选择相机">
              <a-select-option v-for="camera in sceneExtras.cameraViews" :key="camera.id" :value="camera.id">{{ camera.name }}</a-select-option>
            </a-select>
            <a-select v-else-if="eventDraft.action === 'toggleVisible'" v-model="eventDraft.value" size="small" placeholder="选择对象">
              <a-select-option v-for="obj in sceneObjects" :key="obj.id" :value="obj.id">{{ obj.name }}</a-select-option>
            </a-select>
            <a-input v-else v-model="eventDraft.value" size="small" :placeholder="eventValuePlaceholder" />
            <a-select v-if="eventDraft.action === 'callApi'" v-model="eventDraft.method" size="small" class="event-method-select">
              <a-select-option value="get">GET</a-select-option>
              <a-select-option value="post">POST</a-select-option>
              <a-select-option value="put">PUT</a-select-option>
              <a-select-option value="delete">DELETE</a-select-option>
            </a-select>
            <a-button size="small" @click="saveEvent">{{ editingEventId ? $t('ISM3DEditor.save') : $t('ISM3DEditor.add') }}</a-button>
            <a-button v-if="editingEventId" size="small" @click="cancelEventEdit">{{ $t('ISM3DEditor.cancel') }}</a-button>
          </div>
          <div class="event-list">
            <div v-for="evt in currentEvents" :key="evt.id" class="event-row" :class="{ disabled: evt.enabled === false }">
              <span>{{ evt.trigger }} -> {{ evt.action }}<em v-if="evt.method && evt.action === 'callApi'"> / {{ evt.method.toUpperCase() }}</em><em v-if="evt.value"> / {{ evt.value }}</em></span>
              <a-icon :type="evt.enabled === false ? 'play-circle' : 'pause-circle'" @click="toggleEventEnabled(evt)" />
              <a-icon type="copy" @click="duplicateEvent(evt)" />
              <a-icon type="edit" @click="editEvent(evt)" />
              <a-icon type="delete" @click="removeEvent(evt.id)" />
            </div>
          </div>
        </template>
      </a-tab-pane>

      <a-tab-pane v-if="false" key="timeline" tab="时间轴">
        <div class="workbench-toolbar">
          <a-button size="small" icon="plus" @click="addKeyframe" :disabled="!selectedObject">添加关键帧</a-button>
          <a-switch size="small" :checked="!!sceneExtras.timelineEnabled" @change="updateSceneExtra('timelineEnabled', $event)" />
        </div>
        <div class="inline-form">
          <span>速度</span>
          <a-input-number :value="sceneExtras.timelineSpeed || 1" size="small" :min="0.1" :max="8" :step="0.1" @change="updateSceneExtra('timelineSpeed', $event || 1)" />
          <a-button size="small" @click="clearTimelineForSelected" :disabled="!selectedObject">清当前</a-button>
        </div>
        <div class="inline-form">
          <span>时长</span>
          <a-input-number :value="sceneExtras.timelineDuration || 0" size="small" :min="0" :step="0.5" @change="updateSceneExtra('timelineDuration', $event || 0)" />
          <a-switch size="small" :checked="sceneExtras.timelineLoop !== false" checked-children="循环" un-checked-children="单次" @change="updateSceneExtra('timelineLoop', $event)" />
        </div>
        <div class="inline-form">
          <span>属性</span>
          <a-select v-model="timelineDraft.property" size="small">
            <a-select-option value="transform">Transform</a-select-option>
            <a-select-option value="position">Position</a-select-option>
            <a-select-option value="rotation">Rotation</a-select-option>
            <a-select-option value="scale">Scale</a-select-option>
          </a-select>
          <a-input-number v-model="timelineDraft.time" size="small" :min="0" :step="0.5" />
        </div>
        <div class="workbench-toolbar">
          <a-button size="small" type="danger" @click="clearAllTimeline" :disabled="!(sceneExtras.timelineTracks && sceneExtras.timelineTracks.length)">清空时间轴</a-button>
        </div>
        <div class="timeline-list">
          <div v-for="track in sortedTimelineTracks" :key="track.id" class="timeline-row" :class="{ disabled: track.enabled === false }">
            <span>{{ track.objectName }}</span>
            <span>{{ track.property }}</span>
            <a-input-number :value="track.time" size="small" :min="0" :step="0.5" @change="updateTimelineTrack(track.id, { time: $event || 0 })" />
            <a-icon :type="track.enabled === false ? 'pause-circle' : 'play-circle'" @click="updateTimelineTrack(track.id, { enabled: track.enabled === false })" />
            <a-icon type="reload" @click="captureTimelineTrack(track)" />
            <a-icon type="delete" @click="removeTimelineTrack(track.id)" />
          </div>
        </div>
      </a-tab-pane>

      <a-tab-pane v-if="false" key="inspect" tab="检查">
        <div class="metric-grid">
          <div><b>{{ sceneObjects.length }}</b><span>{{ $t('ISM3DEditor.object') }}</span></div>
          <div><b>{{ externalModelCount }}</b><span>{{ $t('ISM3DEditor.externalModel') }}</span></div>
          <div><b>{{ animatedCount }}</b><span>动画对象</span></div>
          <div><b>{{ hiddenCount }}</b><span>隐藏对象</span></div>
          <div><b>{{ lockedCount }}</b><span>锁定对象</span></div>
          <div><b>{{ invalidTimelineCount }}</b><span>失效关键帧</span></div>
          <div><b>{{ invalidEventCount }}</b><span>失效事件</span></div>
          <div><b>{{ sceneExtras.resourceLibrary ? sceneExtras.resourceLibrary.length : 0 }}</b><span>资源</span></div>
        </div>
        <div class="inspect-actions">
          <a-button size="small" @click="$emit('batch-layer', { action: 'showAll', ids: sceneObjects.map(item => item.id) })" :disabled="!hiddenCount">显示全部</a-button>
          <a-button size="small" @click="$emit('batch-layer', { action: 'unlockAll', ids: sceneObjects.map(item => item.id) })" :disabled="!lockedCount">解锁全部</a-button>
          <a-button size="small" @click="repairTimelineTracks" :disabled="!invalidTimelineCount">清理失效关键帧</a-button>
          <a-button size="small" @click="repairInvalidEvents" :disabled="!invalidEventCount">清理失效事件</a-button>
        </div>
        <a-alert v-for="item in performanceWarnings" :key="item" type="warning" :message="item" show-icon style="margin-bottom:6px" />
      </a-tab-pane>
    </a-tabs>
  </div>
</template>

<script>
import { COMPONENT_LIBRARY } from '../Objects/IndustrialObjects'
import { INDUSTRIAL_ASSETS } from '../../assets/industrialAssets'

function clone(data) {
  return JSON.parse(JSON.stringify(data || {}))
}

function defaultExtras() {
  return {
    environment: 'studio',
    fogEnabled: false,
    fogColor: '#d8e4f0',
    defaultCameraId: '',
    cameraViews: [],
    eventFlows: [],
    timelineEnabled: false,
    timelineSpeed: 1,
    timelineLoop: true,
    timelineDuration: 0,
    timelineTracks: [],
    resourceLibrary: []
  }
}

export default {
  name: 'AdvancedWorkbench',
  i18n: require('@/i18n/language'),
  props: {
    sceneObjects: { type: Array, default: () => [] },
    selectedId: { type: String, default: null },
    sceneExtras: { type: Object, default: () => defaultExtras() }
  },
  data() {
    return {
      activeTab: 'camera',
      assetKeyword: '',
      assetCategory: 'all',
      assetSource: 'all',
      layerKeyword: '',
      groupFilter: '',
      groupNameDraft: '',
      textureUrl: '',
      modelUrl: '',
      editingEventId: '',
      selectedCameraId: '',
      eventDraft: {
        trigger: 'click',
        action: 'focus',
        value: '',
        method: 'get'
      },
      timelineDraft: {
        property: 'transform',
        time: 1
      },
      materialPresets: [
        { id: 'steel', name: '拉丝金属', color: '#9aa4ad', metalness: 0.85, roughness: 0.28, opacity: 1, emissive: '#000000' },
        { id: 'glass', name: '工业玻璃', color: '#8bd3ff', metalness: 0, roughness: 0.05, opacity: 0.38, emissive: '#000000' },
        { id: 'warning', name: '警示橙', color: '#ff8a1c', metalness: 0.15, roughness: 0.42, opacity: 1, emissive: '#2a1000' },
        { id: 'rubber', name: '深色橡胶', color: '#1f2328', metalness: 0, roughness: 0.86, opacity: 1, emissive: '#000000' },
        { id: 'emissive', name: '发光蓝', color: '#2f8cff', metalness: 0.05, roughness: 0.25, opacity: 1, emissive: '#1f6fff' }
      ]
    }
  },
  computed: {
    selectedObject() {
      return this.sceneObjects.find(item => item.id === this.selectedId) || null
    },
    selectedLightObject() {
      if (!this.selectedObject || !this.isLightObject(this.selectedObject)) return null
      return this.selectedObject
    },
    selectedLightAngle() {
      if (!this.selectedLightObject || this.selectedLightObject.angle === undefined) return 30
      return Math.round((Number(this.selectedLightObject.angle) || 0) * 180 / Math.PI)
    },
    baseAssets() {
      const libraryAssets = []
      COMPONENT_LIBRARY.forEach(group => {
        (group.items || []).forEach(item => {
          libraryAssets.push({
            id: 'builtin_' + item.type,
            name: item.label,
            type: item.type,
            category: item.type === 'gltf' ? 'model' : (group.name === '视觉特效' ? 'effect' : 'object'),
            icon: item.icon,
            color: item.color,
            typeName: item.typeName,
            source: 'builtin',
            groupName: group.name || '',
            license: item.type === 'gltf' ? '导入' : '免费'
          })
        })
      })
      const customAssets = (this.sceneExtras.resourceLibrary || []).map(item => ({
        ...item,
        category: item.category || 'favorite',
        source: item.isExternalModel || item.modelPath ? 'external' : 'custom',
        license: item.license || '私有'
      }))
      const catalogAssets = (INDUSTRIAL_ASSETS || []).map(item => ({
        ...item,
        source: 'catalog',
        category: item.category || 'model',
        icon: item.icon || 'fas fa-cube',
        license: item.license || 'Project Asset'
      }))
      return customAssets.concat(catalogAssets).concat(libraryAssets)
    },
    filteredAssets() {
      const keyword = this.assetKeyword.trim().toLowerCase()
      return this.baseAssets.filter(item => {
        const matchCategory = this.assetCategory === 'all' || item.category === this.assetCategory
        const matchSource = this.assetSource === 'all' || item.source === this.assetSource
        const matchKeyword = !keyword ||
          (item.name || '').toLowerCase().indexOf(keyword) !== -1 ||
          (item.type || '').toLowerCase().indexOf(keyword) !== -1 ||
          (item.typeName || '').toLowerCase().indexOf(keyword) !== -1 ||
          (item.groupName || '').toLowerCase().indexOf(keyword) !== -1
        return matchCategory && matchSource && matchKeyword
      }).slice(0, 80)
    },
    filteredObjects() {
      const keyword = this.layerKeyword.trim().toLowerCase()
      return this.sceneObjects.filter(item => {
        const matchKeyword = !keyword ||
          (item.name || '').toLowerCase().indexOf(keyword) !== -1 ||
          (item.type || '').toLowerCase().indexOf(keyword) !== -1 ||
          (item.groupName || '').toLowerCase().indexOf(keyword) !== -1
        const matchGroup = !this.groupFilter || item.groupName === this.groupFilter
        return matchKeyword && matchGroup
      })
    },
    groupOptions() {
      const groups = []
      this.sceneObjects.forEach(item => {
        const name = item && item.groupName
        if (name && groups.indexOf(name) === -1) groups.push(name)
      })
      return groups.sort()
    },
    filteredObjectIds() {
      return this.filteredObjects.map(item => item.id)
    },
    sortedTimelineTracks() {
      return (this.sceneExtras.timelineTracks || []).slice().sort((a, b) => {
        const objectName = String(a.objectName || '').localeCompare(String(b.objectName || ''))
        if (objectName !== 0) return objectName
        const propertyName = String(a.property || '').localeCompare(String(b.property || ''))
        if (propertyName !== 0) return propertyName
        return (Number(a.time) || 0) - (Number(b.time) || 0)
      })
    },
    currentEvents() {
      if (!this.selectedObject) return []
      return (this.selectedObject.interactions && this.selectedObject.interactions.events) || []
    },
    eventValuePlaceholder() {
      if (this.eventDraft.action === 'openPanel') return '输入弹窗 URL'
      if (this.eventDraft.action === 'callApi') return '输入接口 URL'
      if (this.eventDraft.trigger === 'data') return '数据点 UUID/名称，可留空'
      return '参数'
    },
    externalModelCount() {
      return this.sceneObjects.filter(item => item.type === 'gltf' || item.isExternalModel).length
    },
    animatedCount() {
      return this.sceneObjects.filter(item => item.autoRotate || item.floatAnim || item.blink || item.flowAnim).length
    },
    hiddenCount() {
      return this.sceneObjects.filter(item => item.visible === false).length
    },
    lockedCount() {
      return this.sceneObjects.filter(item => item.locked).length
    },
    invalidTimelineCount() {
      const ids = new Set(this.sceneObjects.map(item => item.id))
      return (this.sceneExtras.timelineTracks || []).filter(item => item && item.objectId && !ids.has(item.objectId)).length
    },
    invalidEventCount() {
      const ids = new Set(this.sceneObjects.map(item => item.id))
      const cameraIds = new Set((this.sceneExtras.cameraViews || []).map(item => item.id))
      let count = 0
      this.sceneObjects.forEach(obj => {
        const events = (obj.interactions && obj.interactions.events) || []
        events.forEach(evt => {
          if (!evt) return
          if (evt.action === 'toggleVisible' && evt.value && !ids.has(evt.value)) count++
          if (evt.action === 'camera' && evt.value && !cameraIds.has(evt.value)) count++
        })
      })
      return count
    },
    performanceWarnings() {
      const warnings = []
      if (this.sceneObjects.length > 250) warnings.push('对象数量较多，建议按区域拆分或合并静态模型。')
      if (this.externalModelCount > 30) warnings.push('外部模型数量较多，建议使用压缩 GLB 并开启资源复用。')
      if (this.animatedCount > 80) warnings.push('动画对象较多，运行端可能出现帧率波动。')
      if (this.invalidTimelineCount > 0) warnings.push('存在指向已删除对象的时间轴关键帧，可一键清理。')
      if (this.invalidEventCount > 0) warnings.push('存在指向已删除对象或相机的交互事件，可一键清理。')
      const embeddedModels = this.sceneObjects.filter(item => item.gltfBase64 || item.modelBase64).length
      if (embeddedModels > 0) warnings.push('存在内嵌模型数据，建议改为资源 URL 保存。')
      if (warnings.length === 0) warnings.push('当前场景复杂度正常。')
      return warnings
    }
  },
  methods: {
    iconForObject(obj) {
      if (!obj) return 'fas fa-cube'
      if (obj.type === 'gltf') return 'fas fa-file-import'
      if (obj.type === '2dComponent') return 'fas fa-layer-group'
      if (obj.type && obj.type.indexOf('light') !== -1) return 'fas fa-lightbulb'
      return 'fas fa-cube'
    },
    isLightObject(obj) {
      return !!(obj && obj.type && obj.type.indexOf('light') !== -1)
    },
    updateLightAngle(value) {
      if (!this.selectedLightObject) return
      const degree = Number(value)
      this.patchObject(this.selectedLightObject.id, {
        angle: (isNaN(degree) ? 30 : degree) * Math.PI / 180
      })
    },
    assignGroup() {
      if (!this.selectedObject) return
      const name = (this.groupNameDraft || this.groupFilter || '默认分组').trim()
      if (!name) return
      this.groupNameDraft = name
      this.groupFilter = name
      this.$emit('group-selected', { id: this.selectedObject.id, groupName: name })
    },
    assignGroupToFiltered() {
      const name = (this.groupNameDraft || this.groupFilter || '默认分组').trim()
      if (!name || !this.filteredObjectIds.length) return
      this.groupNameDraft = name
      this.groupFilter = name
      this.$emit('batch-layer', { action: 'setGroup', ids: this.filteredObjectIds, groupName: name })
    },
    clearGroupForFiltered() {
      if (!this.filteredObjectIds.length) return
      this.$emit('batch-layer', { action: 'setGroup', ids: this.filteredObjectIds, groupName: '' })
      this.groupFilter = ''
    },
    patchObject(id, changes) {
      this.$emit('update-object', { id, changes })
    },
    addAsset(asset) {
      this.$emit('add-asset', clone(asset))
    },
    assetSourceLabel(asset) {
      if (!asset) return ''
      if (asset.source === 'builtin') return '内置'
      if (asset.source === 'catalog') return '目录'
      if (asset.source === 'external') return 'URL'
      return '收藏'
    },
    isCustomAsset(asset) {
      return !!(asset && (asset.source === 'custom' || asset.source === 'external'))
    },
    removeAsset(id) {
      if (!id) return
      const nextExtras = clone(this.sceneExtras)
      nextExtras.resourceLibrary = (nextExtras.resourceLibrary || []).filter(item => item.id !== id)
      this.$emit('update-extras', nextExtras)
    },
    addExternalModelAsset() {
      const url = (this.modelUrl || '').trim()
      if (!url) return
      const asset = {
        id: 'url_model_' + encodeURIComponent(url).replace(/%/g, '').slice(0, 80),
        name: url.split('/').pop() || 'External Model',
        type: 'gltf',
        typeName: 'GLTF Model',
        category: 'model',
        source: 'external',
        icon: 'fas fa-file-import',
        modelPath: url,
        isExternalModel: true,
        license: 'URL'
      }
      const nextExtras = clone(this.sceneExtras)
      nextExtras.resourceLibrary = nextExtras.resourceLibrary || []
      nextExtras.resourceLibrary = nextExtras.resourceLibrary.filter(item => item.id !== asset.id && item.modelPath !== asset.modelPath)
      nextExtras.resourceLibrary.unshift(asset)
      this.$emit('update-extras', nextExtras)
      this.$emit('add-asset', asset)
      this.modelUrl = ''
    },
    handleAssetDragStart(e, asset) {
      if (!e || !e.dataTransfer) return
      e.dataTransfer.effectAllowed = 'copy'
      e.dataTransfer.setData('application/x-ism3d-asset', JSON.stringify(asset))
    },
    createAssetFromSelection() {
      if (!this.selectedObject) return
      const asset = {
        id: 'asset_' + this.selectedObject.id,
        name: this.selectedObject.name,
        type: this.selectedObject.type,
        typeName: this.selectedObject.typeName,
        category: 'favorite',
        source: 'custom',
        icon: this.iconForObject(this.selectedObject),
        snapshot: clone(this.selectedObject)
      }
      const nextExtras = clone(this.sceneExtras)
      nextExtras.resourceLibrary = nextExtras.resourceLibrary || []
      nextExtras.resourceLibrary = nextExtras.resourceLibrary.filter(item => item.id !== asset.id)
      nextExtras.resourceLibrary.unshift(asset)
      this.$emit('update-extras', nextExtras)
    },
    applyMaterial(preset) {
      if (!this.selectedObject) return
      this.patchObject(this.selectedObject.id, {
        color: preset.color,
        opacity: preset.opacity,
        metalness: preset.metalness,
        roughness: preset.roughness,
        emissive: preset.emissive,
        materialOverridden: true
      })
    },
    applyTexture() {
      if (!this.selectedObject || !this.textureUrl) return
      this.patchObject(this.selectedObject.id, {
        textureData: this.textureUrl,
        materialOverridden: true
      })
    },
    addLight(type) {
      const map = {
        point_light: { label: '点光源', color: '#ffd666' },
        spot_light: { label: '聚光灯', color: '#fff1b8' },
        ambient_light: { label: '环境光', color: '#d6e4ff' },
        directional_light: { label: '平行光', color: '#f0f5ff' }
      }
      this.$emit('add-light', { type, options: map[type] || {} })
    },
    updateSceneExtra(key, value) {
      const nextExtras = clone(this.sceneExtras)
      nextExtras[key] = value
      this.$emit('update-extras', nextExtras)
    },
    saveCameraView() {
      this.$emit('save-camera')
    },
    setDefaultCamera() {
      const nextExtras = clone(this.sceneExtras)
      nextExtras.defaultCameraId = this.selectedCameraId
      this.$emit('update-extras', nextExtras)
    },
    removeCamera(id) {
      const nextExtras = clone(this.sceneExtras)
      nextExtras.cameraViews = (nextExtras.cameraViews || []).filter(item => item.id !== id)
      if (nextExtras.defaultCameraId === id) nextExtras.defaultCameraId = ''
      if (this.selectedCameraId === id) this.selectedCameraId = ''
      this.$emit('update-extras', nextExtras)
    },
    renameCamera(id, name) {
      const nextName = (name || '').trim()
      if (!id || !nextName) return
      const nextExtras = clone(this.sceneExtras)
      nextExtras.cameraViews = (nextExtras.cameraViews || []).map(item => {
        if (item.id !== id) return item
        return { ...item, name: nextName }
      })
      this.$emit('update-extras', nextExtras)
    },
    saveEvent() {
      if (!this.selectedObject) return
      const events = this.currentEvents.slice()
      const existing = events.find(item => item.id === this.editingEventId)
      const nextEvent = {
        id: this.editingEventId || ('evt_' + Date.now()),
        trigger: this.eventDraft.trigger,
        action: this.eventDraft.action,
        value: this.eventDraft.value,
        method: this.eventDraft.method || 'get',
        enabled: existing && existing.enabled === false ? false : true
      }
      const index = events.findIndex(item => item.id === nextEvent.id)
      if (index >= 0) {
        events.splice(index, 1, nextEvent)
      } else {
        events.push(nextEvent)
      }
      this.patchObject(this.selectedObject.id, {
        interactions: {
          ...(this.selectedObject.interactions || {}),
          events
        }
      })
      this.cancelEventEdit()
    },
    editEvent(evt) {
      if (!evt) return
      this.editingEventId = evt.id
      this.eventDraft = {
        trigger: evt.trigger || 'click',
        action: evt.action || 'focus',
        value: evt.value || '',
        method: evt.method || 'get'
      }
    },
    cancelEventEdit() {
      this.editingEventId = ''
      this.eventDraft = {
        trigger: 'click',
        action: 'focus',
        value: '',
        method: 'get'
      }
    },
    addEvent() {
      this.saveEvent()
      this.eventDraft.value = ''
    },
    removeEvent(id) {
      if (!this.selectedObject) return
      const events = this.currentEvents.filter(item => item.id !== id)
      this.patchObject(this.selectedObject.id, {
        interactions: {
          ...(this.selectedObject.interactions || {}),
          events
        }
      })
      if (this.editingEventId === id) this.cancelEventEdit()
    },
    patchEvents(events) {
      if (!this.selectedObject) return
      this.patchObject(this.selectedObject.id, {
        interactions: {
          ...(this.selectedObject.interactions || {}),
          events
        }
      })
    },
    toggleEventEnabled(evt) {
      if (!evt) return
      const events = this.currentEvents.map(item => {
        if (item.id !== evt.id) return item
        return {
          ...item,
          enabled: item.enabled === false
        }
      })
      this.patchEvents(events)
    },
    duplicateEvent(evt) {
      if (!evt) return
      const events = this.currentEvents.slice()
      events.push({
        ...evt,
        id: 'evt_' + Date.now(),
        enabled: evt.enabled !== false
      })
      this.patchEvents(events)
    },
    addKeyframe() {
      if (!this.selectedObject) return
      const nextExtras = clone(this.sceneExtras)
      nextExtras.timelineTracks = nextExtras.timelineTracks || []
      const property = this.timelineDraft.property || 'transform'
      const time = Number(this.timelineDraft.time) || 0
      nextExtras.timelineTracks.push({
        id: 'kf_' + Date.now(),
        objectId: this.selectedObject.id,
        objectName: this.selectedObject.name,
        property: property,
        time: time,
        enabled: true,
        value: this.getTimelineValue(property)
      })
      nextExtras.timelineEnabled = true
      nextExtras.timelineDuration = Math.max(Number(nextExtras.timelineDuration) || 0, time)
      this.timelineDraft.time = time + 1
      this.$emit('update-extras', nextExtras)
    },
    getTimelineValue(property) {
      const obj = this.selectedObject || {}
      if (property === 'position') return { x: obj.x || 0, y: obj.y || 0, z: obj.z || 0 }
      if (property === 'rotation') return { x: obj.rx || 0, y: obj.ry || 0, z: obj.rz || 0 }
      if (property === 'scale') return { x: obj.sx || 1, y: obj.sy || 1, z: obj.sz || 1 }
      return {
        x: obj.x || 0,
        y: obj.y || 0,
        z: obj.z || 0,
        rx: obj.rx || 0,
        ry: obj.ry || 0,
        rz: obj.rz || 0,
        sx: obj.sx || 1,
        sy: obj.sy || 1,
        sz: obj.sz || 1
      }
    },
    removeTimelineTrack(id) {
      const nextExtras = clone(this.sceneExtras)
      nextExtras.timelineTracks = (nextExtras.timelineTracks || []).filter(item => item.id !== id)
      this.$emit('update-extras', nextExtras)
    },
    updateTimelineTrack(id, changes) {
      const nextExtras = clone(this.sceneExtras)
      nextExtras.timelineTracks = (nextExtras.timelineTracks || []).map(item => {
        if (item.id !== id) return item
        return { ...item, ...changes }
      })
      this.$emit('update-extras', nextExtras)
    },
    captureTimelineTrack(track) {
      if (!track || !this.selectedObject || track.objectId !== this.selectedObject.id) return
      this.updateTimelineTrack(track.id, {
        objectName: this.selectedObject.name,
        value: this.getTimelineValue(track.property || 'transform')
      })
    },
    clearTimelineForSelected() {
      if (!this.selectedObject) return
      const nextExtras = clone(this.sceneExtras)
      nextExtras.timelineTracks = (nextExtras.timelineTracks || []).filter(item => item.objectId !== this.selectedObject.id)
      this.$emit('update-extras', nextExtras)
    },
    clearAllTimeline() {
      const nextExtras = clone(this.sceneExtras)
      nextExtras.timelineTracks = []
      this.$emit('update-extras', nextExtras)
    },
    repairTimelineTracks() {
      const ids = new Set(this.sceneObjects.map(item => item.id))
      const nextExtras = clone(this.sceneExtras)
      nextExtras.timelineTracks = (nextExtras.timelineTracks || []).filter(item => item && item.objectId && ids.has(item.objectId))
      this.$emit('update-extras', nextExtras)
    },
    repairInvalidEvents() {
      const ids = new Set(this.sceneObjects.map(item => item.id))
      const cameraIds = new Set((this.sceneExtras.cameraViews || []).map(item => item.id))
      this.sceneObjects.forEach(obj => {
        const events = (obj.interactions && obj.interactions.events) || []
        const nextEvents = events.filter(evt => {
          if (!evt) return false
          if (evt.action === 'toggleVisible' && evt.value && !ids.has(evt.value)) return false
          if (evt.action === 'camera' && evt.value && !cameraIds.has(evt.value)) return false
          return true
        })
        if (nextEvents.length !== events.length) {
          this.patchObject(obj.id, {
            interactions: {
              ...(obj.interactions || {}),
              events: nextEvents
            }
          })
        }
      })
    }
  }
}
</script>

<style scoped>
.advanced-workbench {
  height: 360px;
  min-height: 260px;
  border-bottom: 1px solid #e8e8e8;
  background: #fff;
  overflow: hidden;
}
.advanced-workbench ::v-deep(.ant-tabs) {
  height: 100%;
}
.advanced-workbench ::v-deep(.ant-tabs-content) {
  height: calc(100% - 38px);
  overflow: auto;
  padding: 0 8px 8px;
}
.workbench-toolbar {
  display: flex;
  gap: 6px;
  align-items: center;
  margin-bottom: 8px;
}
.layer-actions {
  display: grid;
  grid-template-columns: repeat(6, minmax(0, 1fr));
  gap: 4px;
  margin-bottom: 8px;
}
.layer-group-actions {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 4px;
  margin-bottom: 8px;
}
.layer-actions .ant-btn {
  padding: 0 4px;
  font-size: 11px;
}
.layer-group-actions .ant-btn {
  padding: 0 4px;
  font-size: 11px;
}
.asset-filters {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
  margin-bottom: 8px;
}
.asset-filters .ant-tag {
  cursor: pointer;
}
.asset-grid,
.preset-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8px;
}
.asset-card {
  position: relative;
  min-width: 0;
  padding: 6px;
  border: 1px solid #edf0f5;
  border-radius: 6px;
  cursor: pointer;
  background: #fafbfc;
}
.asset-card:hover,
.preset-card:hover {
  border-color: #13c2c2;
}
.asset-thumb {
  height: 52px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fff;
  border-radius: 5px;
  overflow: hidden;
}
.asset-thumb img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}
.asset-thumb i {
  font-size: 24px;
  color: #13c2c2;
}
.asset-name {
  margin-top: 4px;
  font-size: 12px;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.asset-license {
  position: absolute;
  right: 22px;
  top: 4px;
  font-size: 10px;
  color: #fff;
  background: #13c2c2;
  border-radius: 3px;
  padding: 0 4px;
}
.asset-source {
  position: absolute;
  left: 6px;
  top: 4px;
  max-width: 52px;
  padding: 0 4px;
  border-radius: 3px;
  color: #595959;
  background: #f0f2f5;
  font-size: 10px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.asset-remove {
  position: absolute;
  right: 5px;
  top: 5px;
  width: 14px;
  height: 14px;
  line-height: 14px;
  border-radius: 50%;
  color: #8c8c8c;
  background: #fff;
  font-size: 10px;
  text-align: center;
  box-shadow: 0 0 0 1px #e8e8e8;
}
.asset-remove:hover {
  color: #ff4d4f;
}
.layer-list,
.camera-list,
.timeline-list,
.event-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.layer-row,
.camera-row,
.timeline-row,
.event-row {
  display: flex;
  align-items: center;
  gap: 8px;
  min-height: 28px;
  padding: 4px 6px;
  border: 1px solid #edf0f5;
  border-radius: 4px;
  background: #fff;
  font-size: 12px;
}
.layer-row.selected,
.camera-row.selected {
  border-color: #13c2c2;
  background: #e6fffb;
}
.camera-name-input {
  flex: 1;
  min-width: 0;
}
.camera-name-input ::v-deep(.ant-input) {
  height: 22px;
  padding: 0 4px;
  font-size: 12px;
}
.layer-name {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.layer-group {
  max-width: 72px;
  flex-shrink: 0;
  padding: 0 5px;
  border-radius: 3px;
  color: #08979c;
  background: #e6fffb;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.event-row span {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.event-row em {
  color: #8c8c8c;
  font-style: normal;
}
.event-row.disabled {
  opacity: 0.5;
}
.event-method-select {
  width: 92px;
  flex: 0 0 92px;
}
.timeline-row.disabled {
  opacity: 0.5;
}
.timeline-row .ant-input-number {
  width: 64px;
}
.preset-card {
  display: flex;
  gap: 8px;
  align-items: center;
  min-width: 0;
  padding: 8px;
  border: 1px solid #edf0f5;
  border-radius: 6px;
  cursor: pointer;
}
.preset-card.disabled {
  opacity: 0.45;
  cursor: not-allowed;
}
.material-swatch {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  border: 1px solid #d9d9d9;
  flex-shrink: 0;
}
.preset-title {
  font-size: 12px;
  color: #333;
}
.preset-desc {
  font-size: 10px;
  color: #888;
}
.inline-form {
  display: grid;
  grid-template-columns: 58px 1fr auto;
  gap: 6px;
  align-items: center;
  margin-bottom: 8px;
  font-size: 12px;
}
.form-value {
  min-width: 68px;
  padding: 2px 8px;
  border-radius: 10px;
  background: #f5f7fa;
  color: #7a8793;
  font-size: 11px;
  line-height: 18px;
  text-align: center;
}
.workbench-color-input,
.inline-form input[type='color'] {
  width: 34px;
  height: 26px;
  padding: 2px;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  background: #fff;
  cursor: pointer;
}
.button-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
}
.metric-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
  margin-bottom: 8px;
}
.metric-grid div {
  padding: 8px;
  border: 1px solid #edf0f5;
  border-radius: 6px;
  text-align: center;
  background: #fafbfc;
}
.metric-grid b {
  display: block;
  font-size: 18px;
  color: #13c2c2;
}
.metric-grid span {
  font-size: 11px;
  color: #666;
}
.inspect-actions {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 6px;
  margin-bottom: 8px;
}
.inspect-actions .ant-btn {
  padding: 0 6px;
  font-size: 11px;
}
.empty-tip {
  padding: 32px 8px;
  color: #999;
  text-align: center;
}
</style>
