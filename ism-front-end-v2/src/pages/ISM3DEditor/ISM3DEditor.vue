/**
 * ISM3DEditor - 3D工业场景编辑器
 * 主容器组件
 * 使用组件本地状态管理（不依赖 Vuex）
 */
<template>
  <div class="ism-3d-editor" @keydown="handleKeyDown" tabindex="0" ref="editorRoot">
    <EditorToolbar
      :mode="editMode"
      :can-undo="history.past.length > 0"
      :can-redo="history.future.length > 0"
      :show-grid-settings="showGridSettings"
      :show-light-settings="showLightSettings"
      :show-inspect-panel="showInspectPanel"
      :selected="selectedObjectId"
      @save="handleSave"
      @import="handleImport"
      @export="handleExport"
      @load-model="openModelFileSelector"
      @undo="handleUndo"
      @redo="handleRedo"
      @set-mode="handleSetMode"
      @toggle-grid-settings="toggleGridSettings"
      @toggle-light-settings="toggleLightSettings"
      @toggle-inspect="toggleInspectPanel"
      @preview="handlePreview"
      @show-help="showHelp"
      @templates="showTemplates"
      @duplicate="handleDuplicateObject(selectedObjectId)"
      @focus="handleFocusSelected"
      @delete="handleDeleteSelected"
    />
    <input
      ref="modelFileInput"
      class="editor-hidden-file-input"
      type="file"
      accept=".glb,.gltf"
      @change="handleModelSelect"
    />

    <div class="editor-page-bar">
      <div class="editor-page-select">
        <span class="editor-page-label">{{ $t('ISM3DEditor.page') }}</span>
        <a-select
          :value="currentPageId"
          size="small"
          style="width: 220px"
          :loading="pageLoading || pageSwitching"
          :disabled="pageLoading || pageSwitching"
          @change="handleSwitchPage"
        >
          <a-select-option v-for="page in scenePages" :key="page.id" :value="page.id">
            <span>{{ page.title }}</span>
            <span v-if="page.isHome" class="editor-page-home">{{ $t('ISM3DEditor.homePage') }}</span>
          </a-select-option>
        </a-select>
      </div>
      <div class="editor-page-actions">
        <a-button size="small" icon="plus" :disabled="pageLoading || pageSwitching" @click="openPageDialog('add')">{{ $t('ISM3DEditor.add') }}</a-button>
        <a-button size="small" icon="edit" :disabled="!currentPageId || pageLoading || pageSwitching" @click="openPageDialog('edit')">{{ $t('ISM3DEditor.rename') }}</a-button>
        <a-button size="small" icon="copy" :disabled="!currentPageId || pageLoading || pageSwitching" @click="copyCurrentPage">{{ $t('ISM3DEditor.copy') }}</a-button>
        <a-button size="small" icon="home" :disabled="!currentPageId || currentPageIsHome || pageLoading || pageSwitching" @click="setCurrentPageHome">{{ $t('ISM3DEditor.setAsHome') }}</a-button>
        <a-button size="small" icon="delete" :disabled="scenePages.length <= 1 || pageLoading || pageSwitching" @click="deleteCurrentPage">{{ $t('ISM3DEditor.delete') }}</a-button>
      </div>
      <div v-if="pageLoading || pageSwitching" class="editor-page-loading-inline">
        <a-spin size="small" />
        <span>{{ $t('ISM3DEditor.pageLoading') }}</span>
      </div>
    </div>
    <a-modal
      :visible="pageDialogVisible"
      :title="pageDialogMode === 'add' ? $t('ISM3DEditor.addPage') : $t('ISM3DEditor.renamePage')"
      :confirm-loading="pageDialogLoading"
      @ok="submitPageDialog"
      @cancel="pageDialogVisible = false"
    >
      <a-input v-model="pageForm.name" :placeholder="$t('ISM3DEditor.inputPageName')" @pressEnter="submitPageDialog" />
    </a-modal>

    <div class="editor-main">
      <div v-if="pageLoading || pageSwitching" class="editor-page-loading-mask">
        <a-spin />
        <span>{{ $t('ISM3DEditor.pageLoading') }}</span>
      </div>
      <div class="editor-left-panels" :class="{ collapsed: leftPanelCollapsed }"><button class="left-panel-toggle" @click="leftPanelCollapsed = !leftPanelCollapsed" :title="leftPanelCollapsed ? $t('ISM3DEditor.expandToolbox') : $t('ISM3DEditor.collapseToolbox')"><i :class="leftPanelCollapsed ? 'fas fa-chevron-right' : 'fas fa-chevron-left'"></i></button><ToolboxPanel
        class="editor-toolbox"
        :sceneObjects="sceneObjects"
        :selectedId="selectedObjectId"
        :showGridSettings="showGridSettings"
        :gridSettings="gridSettings"
        @add-object="handleAddObject"
        @drag-start="onDragStart"
        @select="handleSelectObject"
        @delete="handleDeleteObject"
        @grid-update="updateGrid"
      /></div>

      <!-- 中间画布 -->
      <div class="editor-canvas-container">
        <SceneCanvas
          ref="sceneCanvas"
          :selected-id="selectedObjectId"
          :selected-obj="selectedObject"
          :objects="sceneObjects"
          :mode="editMode"
          :grid-size="gridSize"
          :grid-color="gridColor"
          @object-selected="handleObjectSelected"
          @object-transformed="handleObjectTransformed"
          @overlay-edit-start="handleOverlayEditStart"
          @overlay-updated="handleOverlayUpdated"
          @focus-selected="handleFocusSelected"
          @frame-all="handleFrameAll"
          @delete-selected="handleDeleteSelected"
          @drop-object="handleDropObject"
          @initialized="onSceneInitialized"
          @create-flow-pipe="handleCreateFlowPipe"
          @drawing-pipe-start="handleDrawingPipeStart"
          @drawing-pipe-cancel="handleDrawingPipeCancel"
          @move-up="handleMoveUp"
          @move-down="handleMoveDown"
          @toggle-lock="handleToggleLock"
          @load-model="openModelFileSelector"
          @gltf-animations-loaded="handleGLTFAnimationsLoaded"
        />
      </div>

      <!-- 右侧面板 -->
      <div class="editor-right-panels" :class="{ collapsed: rightPanelCollapsed }"><button class="right-panel-toggle" @click="rightPanelCollapsed = !rightPanelCollapsed" :title="rightPanelCollapsed ? $t('ISM3DEditor.expandProperties') : $t('ISM3DEditor.collapseProperties')"><i :class="rightPanelCollapsed ? 'fas fa-chevron-left' : 'fas fa-chevron-right'"></i></button>
        <!-- 层级树 (已隐藏)
        <HierarchyTree
          class="editor-hierarchy"
          :objects="sceneObjects"
          :selected-id="selectedObjectId"
          @select="handleSelectObject"
          @delete="handleDeleteObject"
          @duplicate="handleDuplicateObject"
          @focus="handleFocusObject"
          @update="handleTreeUpdate"
        />
        -->

        <ISM2DPropertiesBridge
          v-if="selectedObject && selectedObject.type === '2dComponent'"
          class="editor-properties"
          :current-object="selectedObject"
          :scene-objects="sceneObjects"
          @prop-change="handlePropChange"
          @scene-sync="handleSceneSync"
        />
        <PropertiesPanel
          v-else
          class="editor-properties"
          :current-object="selectedObject"
          :scene-objects="sceneObjects"
          @prop-change="handlePropChange"
          @transform-change="handleTransformChange"
          @text-change="handleTextChange"
          @material-change="handlePropChange"
          @anim-change="handlePropChange"
          @duplicate="handleDuplicateObject"
          @focus="handleFocusObject"
          @delete="handleDeleteSelected"
        />
        <AdvancedWorkbench
          v-if="false"
          class="editor-advanced-workbench"
          :scene-objects="sceneObjects"
          :selected-id="selectedObjectId"
          :scene-extras="sceneExtras"
          @select-object="handleSelectObject"
          @update-object="handleWorkbenchUpdateObject"
          @add-asset="handleWorkbenchAddAsset"
          @add-light="handleWorkbenchAddLight"
          @update-extras="handleSceneExtrasUpdate"
          @save-camera="handleSaveCameraView"
          @update-camera="handleUpdateCameraView"
          @apply-camera="handleApplyCameraView"
          @move-object="handleWorkbenchMoveObject"
          @group-selected="handleGroupSelected"
          @batch-layer="handleWorkbenchBatchLayer"
          @duplicate-object="handleDuplicateObject"
          @delete-object="handleDeleteObject"
        />
      </div>

      <!-- 场景设置面板 -->
      <div v-if="showGridSettings" class="light-settings-dropdown" @click.self="showGridSettings = false">
        <div class="scene-settings-card" :style="dropdownPos.scene">
          <div class="light-settings-header">
            <span class="light-settings-title"><i class="fas fa-cog" style="margin-right:6px;color:#13c2c2;"></i>{{ $t('ISM3DEditor.sceneSettings') }}</span>
            <i class="far fa-question-circle" style="color:#999;cursor:pointer;font-size:13px;" title="调整场景环境效果"></i>
          </div>
          <div class="light-settings-body scene-settings-body">
            <a-tabs class="scene-settings-tabs" default-active-key="basic" size="small" :animated="false">
              <a-tab-pane key="basic" tab="基础">
            <div class="light-section">
              <div class="light-section-title"><i class="fas fa-cube" style="margin-right:6px;color:#666;"></i>{{ $t('ISM3DEditor.modelOptimize') }}</div>
              <a-switch size="small" :checked="!!gridSettings.modelOptimize" @change="onModelOptimizeChange" />
            </div>
            <div class="light-section">
              <div class="light-section-title"><i class="fas fa-image" style="margin-right:6px;color:#666;"></i>{{ $t('ISM3DEditor.backgroundType') }}</div>
              <a-radio-group :value="gridSettings.backgroundMode || 'solid'" size="small" @change="e => onBackgroundModeChange(e.target.value)">
                <a-radio-button value="solid">{{ $t('ISM3DEditor.solidColor') }}</a-radio-button>
                <a-radio-button value="gradient">{{ $t('ISM3DEditor.gradient') }}</a-radio-button>
                <a-radio-button value="image">{{ $t('ISM3DEditor.image') }}</a-radio-button>
              </a-radio-group>
              <div v-if="gridSettings.backgroundMode === 'gradient'" class="light-row" style="margin-top:6px;">
                <span class="light-label">渐变色</span>
                <div class="light-color-wrap">
                  <input type="color" class="light-color-input" :value="gridSettings.backgroundColor2 || '#dbefff'" @input="e => patchGridSettings({ backgroundColor2: e.target.value })" />
                </div>
              </div>
            </div>
            <div class="light-section">
              <div class="light-section-title"><i class="fas fa-globe" style="margin-right:6px;color:#666;"></i>{{ $t('ISM3DEditor.environmentPreset') }}</div>
              <a-radio-group :value="gridSettings.environmentPreset || 'sky'" size="small" @change="e => onEnvironmentPresetChange(e.target.value)">
                <a-radio-button value="sky">{{ $t('ISM3DEditor.sky') }}</a-radio-button>
                <a-radio-button value="sunset">{{ $t('ISM3DEditor.sunset') }}</a-radio-button>
                <a-radio-button value="ocean">{{ $t('ISM3DEditor.ocean') }}</a-radio-button>
                <a-radio-button value="forest">{{ $t('ISM3DEditor.forest') }}</a-radio-button>
                <a-radio-button value="twilight">{{ $t('ISM3DEditor.twilight') }}</a-radio-button>
                <a-radio-button value="night">{{ $t('ISM3DEditor.night') }}</a-radio-button>
              </a-radio-group>
            </div>
              </a-tab-pane>
              <a-tab-pane key="environment" tab="天空/地面">
            <div class="light-section">
              <div class="light-section-title"><i class="fas fa-layer-group" style="margin-right:6px;color:#666;"></i>场景选择</div>
              <a-select :value="gridSettings.sceneEnvironmentPreset || 'clearSky'" size="small" style="width:100%;" @change="onSceneEnvironmentPresetChange">
                <a-select-option v-for="preset in sceneEnvironmentPresets" :key="preset.key" :value="preset.key">{{ preset.label }}</a-select-option>
              </a-select>
              <div class="scene-preset-desc">{{ currentSceneEnvironmentPreset.description }}</div>
            </div>
            <div class="light-section">
              <div class="light-section-title"><i class="fas fa-cloud" style="margin-right:6px;color:#666;"></i>天空盒</div>
              <div class="light-row">
                <span class="light-label">启用</span>
                <a-switch size="small" :checked="gridSettings.skyboxEnabled !== false" @change="onSkyboxEnabledChange" />
              </div>
              <div class="light-row">
                <span class="light-label">样式</span>
                <a-select :value="gridSettings.skyboxPreset || 'horizon'" size="small" style="flex:1;" @change="onSkyboxPresetChange">
                  <a-select-option value="horizon">地平线</a-select-option>
                  <a-select-option value="engineeringSky">工程天空</a-select-option>
                  <a-select-option value="blueSky">晴空</a-select-option>
                  <a-select-option value="sunsetGlow">晚霞</a-select-option>
                  <a-select-option value="deepNight">夜空</a-select-option>
                  <a-select-option value="customImage">自定义贴图</a-select-option>
                  <a-select-option value="customHdri">HDRI</a-select-option>
                </a-select>
              </div>
              <div v-if="gridSettings.skyboxPreset === 'customImage'" class="light-row">
                <span class="light-label">贴图</span>
                <a-input :value="gridSettings.skyboxImage || ''" size="small" style="flex:1;" placeholder="天空盒全景图 URL" @change="e => patchGridSettings({ skyboxImage: e.target.value })" />
                <a-button size="small" @click="openSystemImageSelector('skyboxImage')">{{ $t('ISM3DEditor.selectImage') }}</a-button>
              </div>
              <div v-if="gridSettings.skyboxPreset === 'customHdri'" class="light-row">
                <span class="light-label">HDRI</span>
                <a-input :value="gridSettings.skyboxHdri || ''" size="small" style="flex:1;" placeholder=".hdr URL" @change="e => patchGridSettings({ skyboxHdri: e.target.value })" />
              </div>
            </div>
            <div class="light-section">
              <div class="light-section-title"><i class="fas fa-mountain" style="margin-right:6px;color:#666;"></i>地面</div>
              <div class="light-row">
                <span class="light-label">启用</span>
                <a-switch size="small" :checked="gridSettings.groundEnabled !== false" @change="onGroundEnabledChange" />
              </div>
              <div class="light-row">
                <span class="light-label">风格</span>
                <a-select :value="gridSettings.groundStyle || 'slate'" size="small" style="flex:1;" @change="onGroundStyleChange">
                  <a-select-option value="slate">岩板</a-select-option>
                  <a-select-option value="sand">沙地</a-select-option>
                  <a-select-option value="concrete">混凝土</a-select-option>
                  <a-select-option value="oceanDark">深海</a-select-option>
                  <a-select-option value="tunnelSection">隧道剖面</a-select-option>
                  <a-select-option value="geologySection">地层剖面</a-select-option>
                </a-select>
              </div>
              <template v-if="gridSettings.groundStyle === 'tunnelSection' || gridSettings.groundStyle === 'geologySection'">
                <div class="light-row">
                  <span class="light-label">长度</span>
                  <a-input-number :value="gridSettings.sectionLength || 220" :min="20" :max="1000" :step="1" style="flex:1;" size="small" @change="e => patchGridSettings({ sectionLength: e || 220 })" />
                </div>
                <div class="light-row">
                  <span class="light-label">宽度</span>
                  <a-input-number :value="gridSettings.sectionWidth || 36" :min="6" :max="200" :step="1" style="flex:1;" size="small" @change="e => patchGridSettings({ sectionWidth: e || 36 })" />
                </div>
                <div class="light-row">
                  <span class="light-label">高度</span>
                  <a-input-number :value="gridSettings.sectionHeight || 18" :min="4" :max="120" :step="1" style="flex:1;" size="small" @change="e => patchGridSettings({ sectionHeight: e || 18 })" />
                </div>
              </template>
              <div v-if="gridSettings.groundStyle === 'tunnelSection'" class="light-row">
                <span class="light-label">半径</span>
                <a-input-number :value="gridSettings.tunnelRadius || 18" :min="2" :max="120" :step="1" style="flex:1;" size="small" @change="e => patchGridSettings({ tunnelRadius: e || 18 })" />
              </div>
            </div>
              </a-tab-pane>
              <a-tab-pane key="lighting" tab="光照">
            <div class="light-section">
              <div class="light-section-title"><i class="fas fa-sun" style="margin-right:6px;color:#666;"></i>{{ $t('ISM3DEditor.sceneLighting') }}</div>
              <a-select :value="gridSettings.lightingPreset || 'day'" size="small" @change="onLightingPresetChange" style="width:100%">
                <a-select-option value="day">{{ $t('ISM3DEditor.daylight') }}</a-select-option>
                <a-select-option value="evening">{{ $t('ISM3DEditor.evening') }}</a-select-option>
                <a-select-option value="night">{{ $t('ISM3DEditor.night') }}</a-select-option>
                <a-select-option value="studio">{{ $t('ISM3DEditor.studio') }}</a-select-option>
                <a-select-option value="industrial">{{ $t('ISM3DEditor.industrial') }}</a-select-option>
              </a-select>
            </div>
            <div class="light-section">
              <div class="light-section-title"><i class="fas fa-layer-group" style="margin-right:6px;color:#666;"></i>{{ $t('ISM3DEditor.floorReflection') }}</div>
              <a-radio-group :value="gridSettings.floorReflection || 'none'" size="small" @change="e => onFloorReflectionChange(e.target.value)">
                <a-radio-button value="none">{{ $t('ISM3DEditor.none') }}</a-radio-button>
                <a-radio-button value="matte">{{ $t('ISM3DEditor.matte') }}</a-radio-button>
                <a-radio-button value="low">{{ $t('ISM3DEditor.low') }}</a-radio-button>
                <a-radio-button value="medium">{{ $t('ISM3DEditor.medium') }}</a-radio-button>
                <a-radio-button value="strong">{{ $t('ISM3DEditor.strong') }}</a-radio-button>
              </a-radio-group>
            </div>
            <div class="light-section">
              <div class="light-section-title"><i class="fas fa-expand" style="margin-right:6px;color:#666;"></i>{{ $t('ISM3DEditor.enhanceDepth') }}</div>
              <a-switch size="small" :checked="!!gridSettings.enhanceDepth" @change="onEnhanceDepthChange" />
            </div>
              </a-tab-pane>
              <a-tab-pane key="background" tab="背景">
            <div class="light-section">
              <div class="light-section-title"><i class="fas fa-fill-drip" style="margin-right:6px;color:#666;"></i>{{ $t('ISM3DEditor.backgroundColor') }}</div>
              <div class="light-row">
                <div class="light-color-wrap">
                  <input type="color" class="light-color-input" :value="gridSettings.backgroundColor || '#ffffff'" @input="e => patchGridSettings({ backgroundColor: e.target.value })" />
                </div>
                <span style="color:#999;font-size:11px;">{{ gridSettings.backgroundColor || '#ffffff' }}</span>
              </div>
            </div>
            <div v-if="gridSettings.backgroundMode === 'image'" class="light-section">
              <div class="light-section-title"><i class="fas fa-image" style="margin-right:6px;color:#666;"></i>{{ $t('ISM3DEditor.backgroundImage') }}</div>
              <div style="display:flex;gap:6px;">
                <a-input :value="gridSettings.backgroundImage || ''" :placeholder="$t('ISM3DEditor.imageUrl')" :size="'small'" style="flex:1;" @change="e => patchGridSettings({ backgroundImage: e.target.value, backgroundMode: 'image' })" />
                <a-button :size="'small'" @click="openSystemImageSelector()">{{ $t('ISM3DEditor.selectImage') }}</a-button>
              </div>
            </div>
            <div class="light-section">
              <div class="light-section-title"><i class="fas fa-cloud-upload-alt" style="margin-right:6px;color:#666;"></i>{{ $t('ISM3DEditor.backgroundModel') }}</div>
              <div class="scene-bg-model-drop" :class="{dragover: bgModelDragOver}" @dragover.prevent="bgModelDragOver = true" @dragleave="bgModelDragOver = false" @drop.prevent="handleBgModelDrop" @click="openBgModelFileSelector">
                <i class="fas fa-image"></i>
                <div class="scene-bg-model-title">{{ $t('ISM3DEditor.loadModelTitle') }}</div>
                <div class="scene-bg-model-hint">{{ $t('ISM3DEditor.dragOrClickModel') }}</div>
              </div>
              <input ref="bgModelFileInput" class="editor-hidden-file-input" type="file" accept=".glb,.gltf" @change="handleBgModelSelect" />
            </div>
              </a-tab-pane>
              <a-tab-pane key="grid" tab="网格">
            <div class="light-section">
              <div class="light-section-title"><i class="fas fa-border-all" style="margin-right:6px;color:#666;"></i>网格设置</div>
              <div class="light-row">
                <span class="light-label">大小</span>
                <a-input-number :value="gridSize || 10" :min="1" :max="100" @change="onGridSizeChange" style="flex:1;" :size="'small'" />
              </div>
              <div class="light-row">
                <span class="light-label">分割</span>
                <a-input-number :value="gridSettings.divisions || 20" :min="5" :max="100" @change="onGridDivisionsChange" style="flex:1;" :size="'small'" />
              </div>
              <div class="light-row">
                <span class="light-label">颜色</span>
                <div class="light-color-wrap">
                  <input type="color" class="light-color-input" :value="gridColor || '#444444'" @input="onGridColorChange" />
                </div>
              </div>
              <div class="light-row">
                <span class="light-label">中线</span>
                <div class="light-color-wrap">
                  <input type="color" class="light-color-input" :value="gridSettings.colorCenterLine || '#111111'" @input="onCenterLineColorChange" />
                </div>
              </div>
              <div class="light-row">
                <span class="light-label">显示</span>
                <a-switch :checked="showGrid" @change="toggleGrid" />
              </div>
            </div>
              </a-tab-pane>
            </a-tabs>
          </div>
        </div>
      </div>

      <!-- 灯光设置面板 -->
      <div v-if="showLightSettings" class="light-settings-dropdown" @click.self="showLightSettings = false">
        <div class="light-settings-card" :style="dropdownPos.light">
          <div class="light-settings-header">
            <span class="light-settings-title"><i class="fas fa-lightbulb" style="margin-right:6px;color:#13c2c2;"></i>灯光设置</span>
            <i class="far fa-question-circle" style="color:#999;cursor:pointer;font-size:13px;" title="调整场景光照效果"></i>
          </div>
          <div class="light-settings-body">
            <!-- HDRI 环境 -->
            <div class="light-section">
              <div class="light-section-title"><i class="fas fa-globe" style="margin-right:6px;color:#666;"></i>HDRI 环境</div>
              <div class="light-row">
                <span class="light-label">强度</span>
                <div class="light-slider-wrap">
                  <input type="range" class="light-slider" min="0" max="5" step="0.1"
                    :value="(gridSettings.lightSettings && gridSettings.lightSettings.envIntensity) || 1.0"
                    @input="e => onLightSettingChange('envIntensity', parseFloat(e.target.value))" />
                </div>
                <span class="light-value">{{ ((gridSettings.lightSettings && gridSettings.lightSettings.envIntensity) || 1.0).toFixed(1) }}</span>
              </div>
            </div>
            <!-- 环境光 -->
            <div class="light-section">
              <div class="light-section-title"><i class="far fa-lightbulb" style="margin-right:6px;color:#666;"></i>环境光</div>
              <div class="light-row">
                <span class="light-label">颜色</span>
                <div class="light-color-wrap">
                  <input type="color" class="light-color-input"
                    :value="(gridSettings.lightSettings && gridSettings.lightSettings.ambientColor) || '#ffffff'"
                    @input="e => onLightSettingChange('ambientColor', e.target.value)" />
                </div>
              </div>
              <div class="light-row">
                <span class="light-label">强度</span>
                <div class="light-slider-wrap">
                  <input type="range" class="light-slider" min="0" max="3" step="0.1"
                    :value="(gridSettings.lightSettings && gridSettings.lightSettings.ambientIntensity) || 1.0"
                    @input="e => onLightSettingChange('ambientIntensity', parseFloat(e.target.value))" />
                </div>
                <span class="light-value">{{ ((gridSettings.lightSettings && gridSettings.lightSettings.ambientIntensity) || 1.0).toFixed(1) }}</span>
              </div>
            </div>
            <!-- 方向光 -->
            <div class="light-section">
              <div class="light-section-title"><i class="fas fa-sun" style="margin-right:6px;color:#666;"></i>方向光</div>
              <div class="light-row">
                <span class="light-label">颜色</span>
                <div class="light-color-wrap">
                  <input type="color" class="light-color-input"
                    :value="(gridSettings.lightSettings && gridSettings.lightSettings.directionalColor) || '#ffffff'"
                    @input="e => onLightSettingChange('directionalColor', e.target.value)" />
                </div>
              </div>
              <div class="light-row">
                <span class="light-label">强度</span>
                <div class="light-slider-wrap">
                  <input type="range" class="light-slider" min="0" max="3" step="0.1"
                    :value="(gridSettings.lightSettings && gridSettings.lightSettings.directionalIntensity) || 0.9"
                    @input="e => onLightSettingChange('directionalIntensity', parseFloat(e.target.value))" />
                </div>
                <span class="light-value">{{ ((gridSettings.lightSettings && gridSettings.lightSettings.directionalIntensity) || 0.9).toFixed(1) }}</span>
              </div>
            </div>
            <div class="light-hint">
              <i class="fas fa-info-circle" style="margin-right:4px;color:#13c2c2;"></i>
              拖拽场景中的太阳图标调整方向光位置
            </div>
          </div>
        </div>
      </div>

      <div v-if="showInspectPanel" class="light-settings-dropdown" @click.self="showInspectPanel = false">
        <div class="inspect-settings-card" :style="dropdownPos.inspect">
          <div class="light-settings-header">
            <span class="light-settings-title"><i class="fas fa-search" style="margin-right:6px;color:#13c2c2;"></i>检查</span>
          </div>
          <div class="inspect-panel-body">
            <div class="metric-grid inspect-metric-grid">
              <div><b>{{ sceneObjects.length }}</b><span>对象</span></div>
              <div><b>{{ externalModelCount }}</b><span>外部模型</span></div>
              <div><b>{{ animatedCount }}</b><span>动画对象</span></div>
              <div><b>{{ hiddenCount }}</b><span>隐藏对象</span></div>
              <div><b>{{ lockedCount }}</b><span>锁定对象</span></div>
              <div><b>{{ invalidTimelineCount }}</b><span>失效关键帧</span></div>
              <div><b>{{ invalidEventCount }}</b><span>失效事件</span></div>
              <div><b>{{ sceneExtras.resourceLibrary ? sceneExtras.resourceLibrary.length : 0 }}</b><span>资源</span></div>
            </div>
            <div class="inspect-actions">
              <a-button size="small" @click="handleWorkbenchBatchLayer({ action: 'showAll', ids: sceneObjects.map(item => item.id) })" :disabled="!hiddenCount">显示全部</a-button>
              <a-button size="small" @click="handleWorkbenchBatchLayer({ action: 'unlockAll', ids: sceneObjects.map(item => item.id) })" :disabled="!lockedCount">解锁全部</a-button>
              <a-button size="small" @click="repairTimelineTracks" :disabled="!invalidTimelineCount">清理失效关键帧</a-button>
              <a-button size="small" @click="repairInvalidEvents" :disabled="!invalidEventCount">清理失效事件</a-button>
            </div>
            <a-alert v-for="item in performanceWarnings" :key="item" type="warning" :message="item" show-icon style="margin-bottom:6px" />
          </div>
        </div>
      </div>

      <system-image-model
        ref="systemImageModel"
        :networkImageUrl="''"
        @onSelectImage="onSelectBackgroundImage"
      ></system-image-model>

      <!-- 帮助对话框 -->
      <a-modal v-model="showHelpDialog" :title="$t('ISM3DEditor.helpTitle')" :width="600" :footer="null">
        <div style="padding: 16px 0; max-height: 400px; overflow-y: auto;">
          <h3>{{ $t('ISM3DEditor.keyboardShortcuts') }}</h3>
          <ul style="line-height: 2;">
            <li><b>Delete/Backspace</b> - {{ $t('ISM3DEditor.deleteObject') }}</li>
            <li><b>Ctrl+Z</b> - {{ $t('ISM3DEditor.undoAction') }}</li>
            <li><b>Ctrl+Shift+Z / Ctrl+Y</b> - {{ $t('ISM3DEditor.redoAction') }}</li>
            <li><b>Ctrl+D</b> - {{ $t('ISM3DEditor.copyObject') }}</li>
            <li><b>Ctrl+S</b> - {{ $t('ISM3DEditor.saveAction') }}</li>
            <li><b>W</b> - {{ $t('ISM3DEditor.selectMode') }}</li>
            <li><b>E</b> - {{ $t('ISM3DEditor.moveMode') }}</li>
            <li><b>R</b> - {{ $t('ISM3DEditor.rotateMode') }}</li>
            <li><b>S</b> - {{ $t('ISM3DEditor.keyboardScaleMode') }}</li>
            <li><b>G</b> - {{ $t('ISM3DEditor.keyboardFocus') }}</li>
            <li><b>F</b> - {{ $t('ISM3DEditor.keyboardFrameAll') }}</li>
            <li><b>Escape</b> - {{ $t('ISM3DEditor.keyboardCancel') }}</li>
            <li><b>{{ $t('ISM3DEditor.directionKeys') }}</b> - {{ $t('ISM3DEditor.keyboardArrow') }}</li>
          </ul>
          <h3>{{ $t('ISM3DEditor.mouseActions') }}</h3>
          <ul style="line-height: 2;">
            <li><b>{{ $t('ISM3DEditor.leftClickAction') }}</b> - {{ $t('ISM3DEditor.mouseLeft') }}</li>
            <li><b>{{ $t('ISM3DEditor.rightDragPan') }}</b> - {{ $t('ISM3DEditor.mouseRight') }}</li>
            <li><b>{{ $t('ISM3DEditor.middleDragRotate') }}</b> - {{ $t('ISM3DEditor.mouseMiddle') }}</li>
            <li><b>{{ $t('ISM3DEditor.wheelZoom') }}</b> - {{ $t('ISM3DEditor.mouseWheel') }}</li>
          </ul>
        </div>
      </a-modal>

      <!-- 模板选择面板（高科技风格） -->
      <SceneTemplatePicker
        :visible="showTemplatesModal"
        @close="showTemplatesModal = false"
        @select="handleLoadTemplate"
      />
    </div>
  </div>
</template>

<script>
import { message } from 'ant-design-vue'
import { METHOD, request } from '@/utils/request'
import { SYSTEMIMAGEUPLOAD } from '@/services/api'
import {
  getDisplayModelLayerData,
  getDisplayModelPagerLayerData,
  setDisplayModelLayerData,
  DisplayModelPageAdd,
  DisplayModelPageDel,
  DisplayModelPageEdit,
  DisplayModelPageSetHome,
  DisplayModelPageCopy
} from '@/services/displayModel'
import { COMP_COLORS } from './components/Objects/IndustrialObjects'
import EditorToolbar from './components/Toolbar/EditorToolbar.vue'
import ToolboxPanel from './components/Toolbox/ToolboxPanel.vue'
import SceneCanvas from './components/Scene/SceneCanvas.vue'
import HierarchyTree from './components/Hierarchy/HierarchyTree.vue'
import PropertiesPanel from './components/Properties/PropertiesPanel.vue'
import ISM2DPropertiesBridge from './components/Properties/ISM2DPropertiesBridge.vue'
import AdvancedWorkbench from './components/Workbench/AdvancedWorkbench.vue'
import systemImageModel from '@/components/systemImageModel/systemImageModel.vue'
import SceneTemplatePicker from './components/SceneTemplatePicker.vue'
import { createEmptyScenePayload, parseScenePayloadOrEmpty } from './utils/scenePayload'
import { SCENE_TEMPLATES, getTemplateById } from './utils/SceneTemplates'
import { ensureGLTFAnimationGroups, syncLegacyGLTFAnimationFields } from './utils/GLTFAnimationGroups'
import { createDefaultPositionBindings, ensurePositionBindings } from './utils/positionBindings'
import { defaultGridSettings, mergeSceneSettings, getEnvironmentPresetColors } from './utils/sceneSettings'
import { sm4EncryptBase64, sm4DecryptBase64 } from '@/utils/smUtils'
import '@/pages/ISMDisPlay/componentRegistry'
import store from "@/store";
import { mapActions, mapGetters, mapState, mapMutations } from 'vuex'
const generateId = () => 'obj_' + Math.random().toString(36).substr(2, 9)

const SCENE_ENVIRONMENT_PRESETS = [
  {
    key: 'clearSky',
    label: '默认天空',
    description: '通用浅色天空和基础地面，适合普通模型编辑。',
    settings: {
      backgroundMode: 'gradient',
      environmentPreset: 'sky',
      skyboxEnabled: true,
      skyboxPreset: 'horizon',
      skyboxImage: '',
      skyboxHdri: '',
      groundEnabled: true,
      groundStyle: 'slate',
      lightingPreset: 'day',
      floorReflection: 'none',
      enhanceDepth: false,
      backgroundImage: '',
      backgroundColor: '#e8f5ff',
      backgroundColor2: '#b7ddff'
    }
  },
  {
    key: 'tunnel',
    label: '隧道工程',
    description: '蓝灰天空、暗色工程地面和半圆隧道剖面。',
    settings: {
      backgroundMode: 'gradient',
      environmentPreset: 'ocean',
      skyboxEnabled: true,
      skyboxPreset: 'engineeringSky',
      skyboxImage: '',
      skyboxHdri: '',
      groundEnabled: true,
      groundStyle: 'tunnelSection',
      sectionLength: 220,
      sectionWidth: 36,
      sectionHeight: 18,
      tunnelRadius: 18,
      lightingPreset: 'industrial',
      floorReflection: 'none',
      enhanceDepth: true,
      backgroundImage: '',
      backgroundColor: '#d8eef5',
      backgroundColor2: '#7f9ba8',
      colorGrid: '#244760',
      colorCenterLine: '#1c5c86'
    }
  },
  {
    key: 'geology',
    label: '地层剖面',
    description: '偏暖天空和地层断面，用于地下空间、管廊和地质场景。',
    settings: {
      backgroundMode: 'gradient',
      environmentPreset: 'sunset',
      skyboxEnabled: true,
      skyboxPreset: 'sunsetGlow',
      skyboxImage: '',
      skyboxHdri: '',
      groundEnabled: true,
      groundStyle: 'geologySection',
      sectionLength: 220,
      sectionWidth: 36,
      sectionHeight: 18,
      lightingPreset: 'evening',
      floorReflection: 'none',
      enhanceDepth: true,
      backgroundImage: '',
      backgroundColor: '#fff0d9',
      backgroundColor2: '#d59a62'
    }
  },
  {
    key: 'ocean',
    label: '海面空间',
    description: '高亮蓝天和深色地面，适合港口、海洋和开阔空间。',
    settings: {
      backgroundMode: 'gradient',
      environmentPreset: 'ocean',
      skyboxEnabled: true,
      skyboxPreset: 'blueSky',
      skyboxImage: '',
      skyboxHdri: '',
      groundEnabled: true,
      groundStyle: 'oceanDark',
      lightingPreset: 'day',
      floorReflection: 'low',
      enhanceDepth: false,
      backgroundImage: '',
      backgroundColor: '#e5fbff',
      backgroundColor2: '#8fd8ef'
    }
  },
  {
    key: 'cityDusk',
    label: '城市黄昏',
    description: '夕阳天空、混凝土地面和柔和反射，适合城市/园区场景。',
    settings: {
      backgroundMode: 'gradient',
      environmentPreset: 'sunset',
      skyboxEnabled: true,
      skyboxPreset: 'sunsetGlow',
      skyboxImage: '',
      skyboxHdri: '',
      groundEnabled: true,
      groundStyle: 'concrete',
      lightingPreset: 'evening',
      floorReflection: 'low',
      enhanceDepth: true,
      backgroundImage: '',
      backgroundColor: '#fff0d9',
      backgroundColor2: '#ffb36f'
    }
  },
  {
    key: 'nightOps',
    label: '夜间运维',
    description: '深色夜景和冷色补光，适合监控、运维和数字孪生大屏。',
    settings: {
      backgroundMode: 'gradient',
      environmentPreset: 'night',
      skyboxEnabled: true,
      skyboxPreset: 'deepNight',
      skyboxImage: '',
      skyboxHdri: '',
      groundEnabled: true,
      groundStyle: 'slate',
      lightingPreset: 'night',
      floorReflection: 'matte',
      enhanceDepth: true,
      backgroundImage: '',
      backgroundColor: '#101827',
      backgroundColor2: '#273452'
    }
  }
]

function createDefaultAnimateConfig() {
  return {
    selected: [],
    condition: {
      deviceSN: "",
      selectVideoType: 0,
      isBandDevice: false,
      bandType: 1,
      dataID: "",
      dataName: "",
      operator: "",
      OperatorValue: "",
      OperatorMaxValue: "",
    },
    isExpression: false,
    animateList: [
      { id: "Forbidden", name: "禁用" },
      { id: "autoRotate", name: "自旋" },
      { id: "floatAnim", name: "浮动" },
      { id: "blink", name: "闪烁" },
      { id: "visible", name: "显隐" },
    ],
    animateElement: [
      {
        id: "autoRotate",
        elementList: [
          { name: "速度", type: 7, value: 1, min: -5, max: 5, key: "rotateSpeed" },
          { name: "轴向", type: 6, value: 'y', options: ['x', 'y', 'z'], key: "rotateAxis" },
        ]
      },
      {
        id: "floatAnim",
        elementList: [
          { name: "幅度", type: 7, value: 0.15, min: 0.01, max: 1, key: "floatRange" },
          { name: "速度", type: 7, value: 2, min: 0.5, max: 5, key: "floatSpeed" },
        ]
      },
      {
        id: "blink",
        elementList: [
          { name: "速度", type: 7, value: 6, min: 1, max: 10, key: "blinkSpeed" },
          { name: "最低透明度", type: 7, value: 0.2, min: 0, max: 0.8, key: "blinkMin" },
        ]
      }
    ],
  }
}

function createDefaultInteractionConfig() {
  return {
    click: { enabled: false, type: 'none', payload: '', target: '_blank', width: 960, height: 640, title: '', routePath: '', displayMode: 'window', autoClose: false, showTargets: '', hideTargets: '', showTargetsList: [], hideTargetsList: [], deviceKey: '', showUUID: '', showPageUUID: '', isPopUp: false, isContainer: false, deviceType: '', actionConfirm: false, actionVoice: '', actionAuth: [] },
    dblclick: { enabled: false, type: 'none', payload: '', target: '_blank', width: 960, height: 640, title: '', routePath: '', displayMode: 'window', autoClose: false, showTargets: '', hideTargets: '', showTargetsList: [], hideTargetsList: [], deviceKey: '', showUUID: '', showPageUUID: '', isPopUp: false, isContainer: false, deviceType: '', actionConfirm: false, actionVoice: '', actionAuth: [] },
    mouseenter: { enabled: false, type: 'none', payload: '', target: '_blank', width: 960, height: 640, title: '', routePath: '', displayMode: 'window', autoClose: false, showTargets: '', hideTargets: '', showTargetsList: [], hideTargetsList: [], deviceKey: '', showUUID: '', showPageUUID: '', isPopUp: false, isContainer: false, deviceType: '', actionConfirm: false, actionVoice: '', actionAuth: [] },
    mouseleave: { enabled: false, type: 'none', payload: '', target: '_blank', width: 960, height: 640, title: '', routePath: '', displayMode: 'window', autoClose: false, showTargets: '', hideTargets: '', showTargetsList: [], hideTargetsList: [], deviceKey: '', showUUID: '', showPageUUID: '', isPopUp: false, isContainer: false, deviceType: '', actionConfirm: false, actionVoice: '', actionAuth: [] },
  }
}

function cloneSceneObjects(objects) {
  return JSON.parse(JSON.stringify(objects || []))
}

function cloneGLTFAnimations(animations) {
  return Array.isArray(animations) ? JSON.parse(JSON.stringify(animations)) : []
}

function parseImportedSceneFile(content) {
  const raw = String(content || '').trim()
  if (!raw) throw new Error('empty scene file')
  let decrypted = ''
  try {
    decrypted = sm4DecryptBase64(raw)
  } catch (e) {
    decrypted = ''
  }
  const jsonText = decrypted || raw
  return JSON.parse(jsonText)
}

function createDefaultSceneExtras() {
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

function buildActionFromInteraction(type, interaction) {
  const baseAction = {
    type: type,
    actionAuth: interaction.actionAuth || [],
    actionVoice: interaction.actionVoice || '',
    actionConfirm: !!interaction.actionConfirm,
  }

  if (interaction.type === 'link') {
    return {
      ...baseAction,
      action: 'link',
      link: {
        linkType: 'External',
        External: interaction.payload || '',
        OpenExternalType: interaction.target === '_self' ? 'self' : 'new',
        isPopUp: false,
        autoClose: false,
        width: interaction.width || 960,
        height: interaction.height || 640,
        title: interaction.title || '',
      }
    }
  }

  if (interaction.type === 'route') {
    return {
      ...baseAction,
      action: 'link',
      link: {
        linkType: 'Inside',
        External: '',
        OpenExternalType: 'new',
        isPopUp: interaction.displayMode === 'popup',
        autoClose: !!interaction.autoClose,
        width: interaction.width || 960,
        height: interaction.height || 640,
        title: interaction.title || '',
        routePath: interaction.routePath || '',
      }
    }
  }

  if (interaction.type === 'deviceView') {
    return {
      ...baseAction,
      action: 'DeviceView',
      DeviceView: {
        key: interaction.deviceKey || '',
        showUUID: interaction.showUUID || '',
        showPageUUID: interaction.showPageUUID || '',
        type: interaction.deviceType || '',
        isPopUp: !!interaction.isPopUp,
        isContainer: !!interaction.isContainer,
        routePath: interaction.routePath || '',
      }
    }
  }

  if (interaction.type === 'visible') {
    return {
      ...baseAction,
      action: 'visible',
      showItems: interaction.showTargetsList || [],
      hideItems: interaction.hideTargetsList || [],
    }
  }

  if (interaction.type === 'popupText') {
    return {
      ...baseAction,
      action: 'popupText',
      popupText: interaction.payload || '',
    }
  }

  return {
    ...baseAction,
    action: interaction.type || 'none',
  }
}

function evaluateCondition(operator, realValue, compareValue, compareMaxValue) {
  switch (operator) {
    case '==': return realValue == compareValue
    case '!=': return realValue != compareValue
    case '>': return realValue > compareValue
    case '>=': return realValue >= compareValue
    case '<': return realValue < compareValue
    case '<=': return realValue <= compareValue
    case '<>': return !isNaN(compareMaxValue) && realValue >= compareValue && realValue <= compareMaxValue
    case '<!>': return !isNaN(compareMaxValue) && (realValue < compareValue || realValue > compareMaxValue)
    default: return false
  }
}

export default {
  name: 'ISM3DEditor',
  i18n: require('@/i18n/language'),
  components: {
    EditorToolbar,
    ToolboxPanel,
    SceneCanvas,
    HierarchyTree,
    PropertiesPanel,
    ISM2DPropertiesBridge,
    AdvancedWorkbench,
    systemImageModel,
    SceneTemplatePicker,
  },

  data() {
    return {
      SelectPagerID:"",
      scenePages: [],
      currentPageId: "",
      pageLoading: false,
      pageSwitching: false,
      pageDialogVisible: false,
      pageDialogLoading: false,
      pageDialogMode: 'add',
      pageForm: {
        name: ''
      },
      // 场景对象列表
      sceneObjects: [],
      isCharge: false,
      leftPanelCollapsed: false,
      _savedSceneSettings: null,
      // 选中对象 ID
      selectedObjectId: null,
      rightPanelCollapsed: false,
      // 编辑模式：select / move / rotate / scale
      editMode: 'select',
      // 相机模式
      cameraMode: 'orbit',
      objectCounter: 0,
      // 网格显示
      showGrid: true,
      // 网格大小
      gridSize: 10,
      // 网格颜色
      gridColor: '#444444',
      // 网格设置显示
      showGridSettings: false,
      // 灯光设置显示
      showLightSettings: false,
      showInspectPanel: false,
      // 下拉面板位置
      dropdownPos: { scene: { left: '0px' }, light: { left: '0px' }, inspect: { left: '0px' } },
      bgModelDragOver: false,
      // 线框模式
      wireframeMode: false,
      showHelpDialog: false,
      showTemplatesModal: false,
      SCENE_TEMPLATES,
      sceneEnvironmentPresets: SCENE_ENVIRONMENT_PRESETS,
      _dragItem: null,
      // 撤销/重做历史
      history: {
        past: [],
        future: [],
        maxHistory: 50
      },
      // 网格设置
      gridSettings: defaultGridSettings(),
      sceneExtras: createDefaultSceneExtras(),
      // 页面关闭时间记录
      _beforeUnload_time: 0
    }
  },

  computed: {
    selectedObject() {
      if (!this.selectedObjectId) return null
      return this.sceneObjects.find(obj => obj.id === this.selectedObjectId) || null
    },
    currentSceneEnvironmentPreset() {
      const current = this.gridSettings.sceneEnvironmentPreset || 'clearSky'
      return this.sceneEnvironmentPresets.find(preset => preset.key === current) || this.sceneEnvironmentPresets[0]
    },
    currentPage() {
      return this.scenePages.find(page => page.id === this.currentPageId) || null
    },
    currentPageIsHome() {
      return !!(this.currentPage && this.currentPage.isHome)
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

  mounted() {
    let _t =this
    // 聚焦编辑器以接收键盘事件
    this.$refs.editorRoot && this.$refs.editorRoot.focus()
    // 尝试从服务端加载场景数据
    _t.loadSceneFromServer()
    this.$EventBus.$on("readDataPush", this.DealWithUpdateData)
    this.$EventBus.$on("StaticData", this.DealWithUpdateData)
    this.$EventBus.$on("SystemData", this.DealWithUpdateData)
    window.addEventListener('keydown', this.handleGlobalKeyDown, true)
    // 添加页面关闭提示
    window.addEventListener('beforeunload', this.beforeunloadHandler)
    window.addEventListener('unload', this.unloadHandler)
  },

  beforeDestroy() {
    this.selectedObjectId = null
    this.sceneObjects = []
    // 清理 EventBus 监听
    this.$EventBus.$off("readDataPush", this.DealWithUpdateData)
    this.$EventBus.$off("StaticData", this.DealWithUpdateData)
    this.$EventBus.$off("SystemData", this.DealWithUpdateData)
    window.removeEventListener('keydown', this.handleGlobalKeyDown, true)
    // 移除页面关闭提示监听
    window.removeEventListener('beforeunload', this.beforeunloadHandler)
    window.removeEventListener('unload', this.unloadHandler)
  },

  methods: {
    beforeunloadHandler (e) {
      if (!this.isCharge) {
        return
      }
      this._beforeUnload_time = new Date().getTime()
      e = e || window.event
      if (e) {
        e.returnValue = '关闭提示'
      }
      return '关闭提示'
    },
    unloadHandler () {
      if (!this.isCharge) {
        return
      }
      this._gap_time = new Date().getTime() - this._beforeUnload_time
      if (this._gap_time <= 5) {
        console.log("页面关闭")
      }
    },
    ...mapActions('ISMDisPlayEditorTool',[
      'saveLayerDataStruct',
      'setGroupList',
      'setLayerData',
      'SyncLayerData',
      'getLayerDataStruct',
    ]),
    // ========================
    //  对象 CRUD
    // ========================

    /** 添加对象 */
    addObject({ type, options = {} }) {
      const id = generateId()
      var colorKey = options.colorKey || options.defaultColor || options.color || 'blue'
      var mappedColor = COMP_COLORS[colorKey] || colorKey
      if (options.color && options.color.charAt(0) === '#') {
        mappedColor = options.color
      }
      var isTextObject = type === 'text3d' || type === 'textPlain3d' || type === 'dataText'
      var objectColor = isTextObject ? (options.color && options.color.charAt(0) === '#' ? options.color : '#ffffff') : (options.color || mappedColor)
      var label = options.label || type
      var obj = {
        id: id,
        type: type,
        name: options.name || label + '_' + (++this.objectCounter),
        typeName: options.typeName || type,
        x: options.x || 0,
        y: options.y || 0,
        z: options.z || 0,
        rx: options.rx || 0,
        ry: options.ry || 0,
        rz: options.rz || 0,
        sx: options.sx || 1,
        sy: options.sy || 1,
        sz: options.sz || 1,
        color: objectColor,
        opacity: options.opacity !== undefined ? options.opacity : 1,
        metalness: options.metalness !== undefined ? options.metalness : 0.3,
        roughness: options.roughness !== undefined ? options.roughness : 0.7,
        wireframe: options.wireframe || false,
        showShadow: !!options.showShadow,
        visible: options.visible !== false,
        locked: options.locked || false,
        icon: options.icon || '',
        category: options.category || '',
        groupName: options.groupName || '',
        intensity: options.intensity !== undefined ? options.intensity : undefined,
        distance: options.distance !== undefined ? options.distance : undefined,
        angle: options.angle !== undefined ? options.angle : undefined,
        penumbra: options.penumbra !== undefined ? options.penumbra : undefined,
        emissive: options.emissive || '#000000',
        textureData: options.textureData || '',
        mediaUrl: options.mediaUrl || options.imageUrl || options.videoUrl || '',
        imageUrl: options.imageUrl || options.mediaUrl || '',
        videoUrl: options.videoUrl || options.mediaUrl || '',
        webUrl: options.webUrl || '',
        mediaWidth: options.mediaWidth !== undefined ? options.mediaWidth : (type === 'video3d' ? 1.6 : 1.4),
        mediaAspect: options.mediaAspect !== undefined ? options.mediaAspect : (type === 'video3d' ? 16 / 9 : 4 / 3),
        mediaAutoplay: options.mediaAutoplay !== undefined ? options.mediaAutoplay : true,
        mediaLoop: options.mediaLoop !== undefined ? options.mediaLoop : true,
        mediaMuted: options.mediaMuted !== undefined ? options.mediaMuted : true,
        uiX: options.uiX !== undefined ? options.uiX : 40,
        uiY: options.uiY !== undefined ? options.uiY : 40,
        uiWidth: options.uiWidth !== undefined ? options.uiWidth : (type === 'webEmbed' ? 420 : 180),
        uiHeight: options.uiHeight !== undefined ? options.uiHeight : (type === 'webEmbed' ? 260 : 100),
        uiRotation: options.uiRotation !== undefined ? options.uiRotation : 0,
        textContent: options.textContent || '',
        fontSize: options.fontSize || 16,
        textBgColor: options.textBgColor || (type === 'text3d' || type === 'textPlain3d' || type === 'dataText' ? '#00000000' : '#000000'),
        labelOpacity: options.labelOpacity !== undefined ? options.labelOpacity : 1,
        textBgOpacity: options.textBgOpacity !== undefined ? options.textBgOpacity : (type === 'text3d' || type === 'textPlain3d' || type === 'dataText' ? 0 : 0.2),
        labelFaceCamera: options.labelFaceCamera !== undefined ? options.labelFaceCamera : true,
        labelFixedSize: options.labelFixedSize !== undefined ? options.labelFixedSize : true,
        labelFontFamily: options.labelFontFamily || '系统默认',
        labelRenderMode: options.labelRenderMode || 'component',
        labelShowBorder: options.labelShowBorder !== undefined ? !!options.labelShowBorder : false,
        labelBorderWidth: options.labelBorderWidth !== undefined ? options.labelBorderWidth : 1,
        labelBorderColor: options.labelBorderColor || '#ffffff',
        autoRotate: options.autoRotate || false,
        rotateSpeed: options.rotateSpeed !== undefined ? options.rotateSpeed : 1,
        rotateAxis: options.rotateAxis || 'y',
        floatAnim: options.floatAnim || false,
        floatRange: options.floatRange !== undefined ? options.floatRange : 0.15,
        floatSpeed: options.floatSpeed !== undefined ? options.floatSpeed : 2,
        blink: options.blink || false,
        blinkSpeed: options.blinkSpeed !== undefined ? options.blinkSpeed : 6,
        blinkMin: options.blinkMin !== undefined ? options.blinkMin : 0.2,
        gltfAnimationPlaying: !!options.gltfAnimationPlaying,
        gltfAnimationName: options.gltfAnimationName || '',
        gltfAnimationNames: Array.isArray(options.gltfAnimationNames) ? options.gltfAnimationNames.slice() : (options.gltfAnimationName ? [options.gltfAnimationName] : []),
        gltfAnimationSpeed: options.gltfAnimationSpeed !== undefined ? options.gltfAnimationSpeed : 1,
        gltfAnimationLoop: options.gltfAnimationLoop !== false,
        gltfAnimations: cloneGLTFAnimations(options.gltfAnimations),
        gltfAnimationGroups: ensureGLTFAnimationGroups(options, { createDefault: true }),
        gltfAnimationConditionEnabled: !!options.gltfAnimationConditionEnabled,
        gltfAnimationCondition: options.gltfAnimationCondition ? JSON.parse(JSON.stringify(options.gltfAnimationCondition)) : {
          isBandDevice: false,
          deviceSN: '',
          DeviceName: '',
          dataID: '',
          dataName: '',
          operator: '',
          OperatorValue: '',
          OperatorMaxValue: ''
        },
        animate: options.animate ? JSON.parse(JSON.stringify(options.animate)) : createDefaultAnimateConfig(),
        interactions: options.interactions ? JSON.parse(JSON.stringify(options.interactions)) : createDefaultInteractionConfig(),
        wsKey: options.wsKey || '',
        bindProp: options.bindProp || '',
        bindTransform: options.bindTransform || 'direct',
        bindScale: options.bindScale !== undefined ? options.bindScale : 1,
        bindOffset: options.bindOffset !== undefined ? options.bindOffset : 0,
        positionBindings: options.positionBindings ? JSON.parse(JSON.stringify(options.positionBindings)) : createDefaultPositionBindings(),
        realTimeValue: options.realTimeValue !== undefined ? options.realTimeValue : '',
        dataFormat: options.dataFormat || '{value}',
        is2DComponent: !!options.is2DComponent,
        label2D: options.label2D || options.label || options.name || '',
        source2D: options.source2D ? JSON.parse(JSON.stringify(options.source2D)) : null,
        source2DMeta: options.source2DMeta ? JSON.parse(JSON.stringify(options.source2DMeta)) : null,
        action: options.action ? JSON.parse(JSON.stringify(options.action)) : [],
        isExternalModel: !!options.isExternalModel,
        isBackground: !!options.isBackground,
        fitSize: options.fitSize !== undefined ? options.fitSize : undefined,
        modelPath: options.modelPath || '',
        materialOverridden: !!options.materialOverridden
      }
      if (obj.type === '2dComponent' && obj.source2D) {
        obj.source2D.style = obj.source2D.style || {}
        obj.source2D.style.position = obj.source2D.style.position || {}

        if (options.base && options.base.info && options.base.info.style && options.base.info.style.position) {
          const basePos = options.base.info.style.position
          if (!obj.source2D.style.position.w && basePos.w) {
            obj.source2D.style.position.w = basePos.w
          }
          if (!obj.source2D.style.position.h && basePos.h) {
            obj.source2D.style.position.h = basePos.h
          }
        }

        this.sync2DWrapperFields(obj)
      }
      obj.action = Array.isArray(obj.action) ? obj.action : this.getObjectActions(obj)
      this.commitSceneChange(function(list) {
        list.push(obj)
      })
      this.selectedObjectId = id
      return obj
    },

    /** 删除对象 */
    removeObject(objectId) {
      const idx = this.sceneObjects.findIndex(o => o.id === objectId)
      if (idx !== -1) {
        let removed = null
        this.commitSceneChange(function(list) {
          removed = list[idx]
          list.splice(idx, 1)
        })
        if (this.selectedObjectId === objectId) {
          this.selectedObjectId = null
        }
        // 如果删除的是背景模型，且场景中已无其他背景模型，恢复网格显示
        if (removed && removed.isBackground) {
          var stillHasBg = this.sceneObjects.some(function(o) { return o.isBackground })
          if (!stillHasBg) {
            this.showGrid = true
            if (this.$refs.sceneCanvas && this.$refs.sceneCanvas.toggleGrid) {
              this.$refs.sceneCanvas.toggleGrid(true)
            }
          }
        }
      }
    },

    updateObject(objectId, changes) {
      const obj = this.sceneObjects.find(o => o.id === objectId)
      if (obj && !obj.locked) {
        Object.assign(obj, changes)
        this.sync2DStyleFromWrapperFields(obj, changes)
      }
    },

    handleGLTFAnimationsLoaded(payload) {
      if (!payload || !payload.objectId) return
      const animations = cloneGLTFAnimations(payload.animations)
      const obj = this.sceneObjects.find(o => o.id === payload.objectId)
      if (!obj) return
      this.$set(obj, 'gltfAnimations', animations)
      const firstAnimationKey = animations.length ? (typeof animations[0] === 'string' ? animations[0] : animations[0].key) : ''
      if (!Array.isArray(obj.gltfAnimationNames)) {
        this.$set(obj, 'gltfAnimationNames', obj.gltfAnimationName ? [obj.gltfAnimationName] : [])
      }
      if (!obj.gltfAnimationName && firstAnimationKey) {
        this.$set(obj, 'gltfAnimationName', firstAnimationKey)
      }
      if (!obj.gltfAnimationNames.length && firstAnimationKey) {
        this.$set(obj, 'gltfAnimationNames', [firstAnimationKey])
      }
      if (obj.gltfAnimationSpeed === undefined) this.$set(obj, 'gltfAnimationSpeed', 1)
      if (obj.gltfAnimationLoop === undefined) this.$set(obj, 'gltfAnimationLoop', true)
      if (obj.gltfAnimationPlaying === undefined) this.$set(obj, 'gltfAnimationPlaying', false)
      const groups = ensureGLTFAnimationGroups(obj, { defaultAnimationKey: firstAnimationKey })
      this.$set(obj, 'gltfAnimationGroups', groups)
      syncLegacyGLTFAnimationFields(obj, groups)
      this.sceneObjects = [...this.sceneObjects]
    },

    applyGLTFAnimationCondition(obj, realData, currentData) {
      if (!obj) return
      const groups = ensureGLTFAnimationGroups(obj, { createDefault: false })
      if (!groups.length) return
      let changed = false
      const nextGroups = groups.map(group => {
        if (!group.conditionEnabled || !group.condition || !group.condition.dataID) return group
        const condition = group.condition
        if (condition.isBandDevice && realData.DeviceUuid !== condition.deviceSN) return group
        if (condition.dataID !== currentData.ModelDataUuid && condition.dataID !== currentData.Uuid) return group
        let isStart = true
        if (condition.operator) {
          const realValue = parseFloat(currentData.Value)
          const operatorValue = parseFloat(condition.OperatorValue)
          const operatorMaxValue = parseFloat(condition.OperatorMaxValue)
          if (isNaN(realValue) || isNaN(operatorValue)) return group
          isStart = evaluateCondition(condition.operator, realValue, operatorValue, operatorMaxValue)
        }
        if (group.playing === isStart) return group
        changed = true
        return { ...group, playing: isStart }
      })
      if (!changed) return
      const legacy = {}
      syncLegacyGLTFAnimationFields(legacy, nextGroups)
      this.updateObject(obj.id, {
        gltfAnimationGroups: nextGroups,
        ...legacy
      })
    },
    applyPositionBindings(obj, realData, currentData) {
      if (!obj || !obj.positionBindings) return
      const changes = {}
      let changed = false
      ;['x', 'y', 'z'].forEach((axis) => {
        const binding = obj.positionBindings[axis]
        if (!binding || !binding.dataID) return
        if (binding.isBandDevice && realData.DeviceUuid !== binding.deviceSN) return
        if (binding.dataID !== currentData.ModelDataUuid && binding.dataID !== currentData.Uuid) return
        const rawValue = parseFloat(currentData.Value)
        if (isNaN(rawValue)) return
        let nextValue = rawValue
        if (binding.transform === 'scale') {
          nextValue = rawValue * (binding.scale !== undefined ? Number(binding.scale) : 1)
        } else if (binding.transform === 'offset') {
          nextValue = rawValue + (binding.offset !== undefined ? Number(binding.offset) : 0)
        }
        changes[axis] = nextValue
        changed = true
      })
      if (changed) {
        this.updateObject(obj.id, changes)
      }
    },
    applyGenericDataBinding(obj, realData, currentData) {
      if (!obj || !obj.bindProp || !obj.dataID) return
      if (obj.isBandDevice && realData.DeviceUuid !== obj.deviceSN) return
      if (obj.dataID !== currentData.ModelDataUuid && obj.dataID !== currentData.Uuid) return
      const raw = currentData.Value
      const rawNumber = parseFloat(raw)
      let value = raw
      if (!isNaN(rawNumber)) {
        value = rawNumber
        if (obj.bindTransform === 'scale') {
          value = rawNumber * (obj.bindScale !== undefined ? Number(obj.bindScale) : 1)
        } else if (obj.bindTransform === 'offset') {
          value = rawNumber + (obj.bindOffset !== undefined ? Number(obj.bindOffset) : 0)
        }
      }
      const changes = {}
      if (obj.bindProp === 'visible') {
        changes.visible = !(value === false || value === 'false' || value === '0' || value === 0)
      } else if (obj.bindProp === 'textContent') {
        changes.textContent = String(raw)
      } else if (obj.bindProp === 'color') {
        changes.color = typeof raw === 'string' && raw.charAt(0) === '#' ? raw : obj.color
      } else if (obj.bindProp === 'opacity') {
        if (isNaN(rawNumber)) return
        changes.opacity = Math.max(0, Math.min(1, Number(value)))
      } else if (obj.bindProp === 'scale') {
        if (isNaN(rawNumber)) return
        changes.sx = Number(value)
        changes.sy = Number(value)
        changes.sz = Number(value)
      } else {
        return
      }
      this.updateObject(obj.id, changes)
    },

    handleWorkbenchUpdateObject(payload) {
      if (!payload || !payload.id) return
      const id = payload.id
      const changes = payload.changes || {}
      this.commitSceneChange(function(list) {
        const obj = list.find(o => o.id === id)
        if (obj) {
          Object.assign(obj, changes)
          this.sync2DStyleFromWrapperFields(obj, changes)
        }
      })
      if (this.$refs.sceneCanvas && this.$refs.sceneCanvas.syncMeshes) {
        this.$refs.sceneCanvas.syncMeshes(this.sceneObjects)
      }
    },

    handleWorkbenchAddAsset(asset) {
      if (!asset) return
      if (asset.source === 'catalog' && asset.fallbackType) {
        const obj = this.addObject({
          type: asset.fallbackType,
          options: {
            label: asset.name,
            name: asset.name,
            typeName: asset.typeName || asset.fallbackType,
            icon: asset.icon,
            color: asset.color || 'blue',
            category: asset.category || 'model',
            modelPath: asset.modelPath || '',
            sourceAssetId: asset.id,
            sourceAssetLicense: asset.license || '',
            sourceAssetUrl: asset.sourceUrl || '',
            isCatalogAsset: true,
            sx: asset.sx || 1,
            sy: asset.sy || 1,
            sz: asset.sz || 1
          }
        })
        if (asset.modelPath) {
          this.$nextTick(() => {
            if (this.$refs.sceneCanvas && this.$refs.sceneCanvas.loadGLTFModelUrl) {
              this.$refs.sceneCanvas.loadGLTFModelUrl(asset.modelPath, obj.id, (loadedId) => {
                this.markCatalogAssetAsGLTF(loadedId, asset)
              }, { silent: true })
            }
          })
        }
        return
      }
      if (asset.type === 'gltf' && asset.modelPath) {
        const obj = this.addObject({
          type: 'gltf',
          options: {
            label: asset.name,
            name: asset.name,
            typeName: asset.typeName || 'GLTF Model',
            icon: asset.icon,
            category: asset.category || 'model',
            modelPath: asset.modelPath,
            isExternalModel: true,
            materialOverridden: false,
            sx: asset.sx || 1,
            sy: asset.sy || 1,
            sz: asset.sz || 1
          }
        })
        this.$nextTick(() => {
          if (this.$refs.sceneCanvas && this.$refs.sceneCanvas.loadGLTFModelUrl) {
            this.$refs.sceneCanvas.loadGLTFModelUrl(asset.modelPath, obj.id)
          }
        })
        return
      }
      if (asset.snapshot) {
        const item = JSON.parse(JSON.stringify(asset.snapshot))
        item.name = item.name || asset.name
        item.x = (item.x || 0) + 1
        item.z = (item.z || 0) + 1
        this.addObject({ type: item.type, options: item })
        return
      }
      this.addObject({
        type: asset.type,
        options: {
          label: asset.name,
          typeName: asset.typeName || asset.type,
          icon: asset.icon,
          color: asset.color || 'blue'
        }
      })
    },

    markCatalogAssetAsGLTF(objectId, asset) {
      if (!objectId || !asset) return
      this.commitSceneChange(function(list) {
        const obj = list.find(o => o.id === objectId)
        if (!obj) return
        obj.type = 'gltf'
        obj.typeName = asset.typeName || obj.typeName || 'GLTF Model'
        obj.modelPath = asset.modelPath || obj.modelPath || ''
        obj.isExternalModel = true
        obj.isCatalogAsset = true
        obj.fallbackType = asset.fallbackType || obj.fallbackType || ''
        obj.sourceAssetId = asset.id || obj.sourceAssetId || ''
        obj.sourceAssetLicense = asset.license || obj.sourceAssetLicense || ''
        obj.sourceAssetUrl = asset.sourceUrl || obj.sourceAssetUrl || ''
        obj.materialOverridden = false
      })
    },

    handleWorkbenchAddLight(payload) {
      if (!payload) return
      const options = payload.options || {}
      this.addObject({
        type: payload.type,
        options: {
          label: options.label || payload.type,
          color: options.color || '#fff1b8',
          y: 3,
          z: 2,
          intensity: 1,
          distance: payload.type === 'point_light' || payload.type === 'spot_light' ? 10 : undefined,
          angle: payload.type === 'spot_light' ? Math.PI / 6 : undefined,
          penumbra: payload.type === 'spot_light' ? 0.25 : undefined
        }
      })
    },

    handleSceneExtrasUpdate(nextExtras) {
      this.sceneExtras = Object.assign(createDefaultSceneExtras(), nextExtras || {})
      if (this.$refs.sceneCanvas && this.$refs.sceneCanvas.applySceneExtras) {
        this.$refs.sceneCanvas.applySceneExtras(this.sceneExtras)
      }
      this.isCharge = true
    },

    handleSaveCameraView() {
      if (!this.$refs.sceneCanvas || !this.$refs.sceneCanvas.__3d) return
      const three = this.$refs.sceneCanvas.__3d
      const nextExtras = JSON.parse(JSON.stringify(this.sceneExtras || createDefaultSceneExtras()))
      nextExtras.cameraViews = nextExtras.cameraViews || []
      const index = nextExtras.cameraViews.length + 1
      nextExtras.cameraViews.push({
        id: 'cam_' + Date.now(),
        name: '视角 ' + index,
        position: {
          x: three.camera.position.x,
          y: three.camera.position.y,
          z: three.camera.position.z
        },
        target: {
          x: three.orbit.target.x,
          y: three.orbit.target.y,
          z: three.orbit.target.z
        },
        fov: three.camera.fov
      })
      this.handleSceneExtrasUpdate(nextExtras)
      message.success('已保存当前相机视角')
    },

    handleApplyCameraView(camera) {
      if (!camera || !this.$refs.sceneCanvas || !this.$refs.sceneCanvas.restoreCamera) return
      this.$refs.sceneCanvas.restoreCamera({
        cameraPosition: camera.position,
        cameraTarget: camera.target,
        cameraFov: camera.fov
      })
    },

    handleUpdateCameraView(cameraId) {
      if (!cameraId || !this.$refs.sceneCanvas || !this.$refs.sceneCanvas.__3d) return
      const three = this.$refs.sceneCanvas.__3d
      const nextExtras = JSON.parse(JSON.stringify(this.sceneExtras || createDefaultSceneExtras()))
      nextExtras.cameraViews = (nextExtras.cameraViews || []).map(item => {
        if (item.id !== cameraId) return item
        return Object.assign({}, item, {
          position: {
            x: three.camera.position.x,
            y: three.camera.position.y,
            z: three.camera.position.z
          },
          target: {
            x: three.orbit.target.x,
            y: three.orbit.target.y,
            z: three.orbit.target.z
          },
          fov: three.camera.fov,
          updatedAt: Date.now()
        })
      })
      this.handleSceneExtrasUpdate(nextExtras)
      message.success('已更新相机视角')
    },

    handleWorkbenchMoveObject(payload) {
      if (!payload || !payload.id) return
      const index = this.sceneObjects.findIndex(item => item.id === payload.id)
      if (index === -1) return
      const nextIndex = index + payload.direction
      if (nextIndex < 0 || nextIndex >= this.sceneObjects.length) return
      this.commitSceneChange(function(list) {
        const moved = list.splice(index, 1)[0]
        list.splice(nextIndex, 0, moved)
        this.normalize2DLayerZIndex(list)
      })
    },

    handleWorkbenchBatchLayer(payload) {
      if (!payload || !payload.action || !Array.isArray(payload.ids) || !payload.ids.length) return
      const idSet = new Set(payload.ids)
      if (payload.action === 'showAll' || payload.action === 'hideAll' || payload.action === 'lockAll' || payload.action === 'unlockAll' || payload.action === 'setGroup') {
        const vm = this
        this.commitSceneChange(function(list) {
          list.forEach(function(obj) {
            if (!obj || !idSet.has(obj.id)) return
            if (payload.action === 'showAll') obj.visible = true
            if (payload.action === 'hideAll') obj.visible = false
            if (payload.action === 'lockAll') obj.locked = true
            if (payload.action === 'unlockAll') obj.locked = false
            if (payload.action === 'setGroup') obj.groupName = payload.groupName || ''
            if (payload.action === 'showAll') vm.sync2DStyleFromWrapperFields(obj, { visible: true })
            if (payload.action === 'hideAll') vm.sync2DStyleFromWrapperFields(obj, { visible: false })
          })
        })
        if (this.$refs.sceneCanvas && this.$refs.sceneCanvas.syncMeshes) {
          this.$refs.sceneCanvas.syncMeshes(this.sceneObjects)
        }
      }
    },

    handleGroupSelected(payload) {
      const targetId = (payload && payload.id) || this.selectedObjectId
      if (!targetId) return
      this.handleWorkbenchUpdateObject({
        id: targetId,
        changes: {
          groupName: (payload && payload.groupName) || '默认分组'
        }
      })
    },

    /** 复制对象 */
    duplicateObject(objectId) {
      const src = this.sceneObjects.find(o => o.id === objectId)
      if (!src) return
      const options = {
        type: src.type,
        typeName: src.typeName,
        name: src.name + '_copy',
        x: src.x + 2,
        y: src.y,
        z: src.z,
        rx: src.rx,
        ry: src.ry,
        rz: src.rz,
        sx: src.sx,
        sy: src.sy,
        sz: src.sz,
        color: src.color,
        icon: src.icon,
        category: src.category,
        groupName: src.groupName,
        intensity: src.intensity,
        distance: src.distance,
        angle: src.angle,
        penumbra: src.penumbra,
        opacity: src.opacity,
        metalness: src.metalness,
        roughness: src.roughness,
        wireframe: src.wireframe,
        showShadow: !!src.showShadow,
        textContent: src.textContent,
        fontSize: src.fontSize,
        textBgColor: src.textBgColor,
        labelOpacity: src.labelOpacity,
        textBgOpacity: src.textBgOpacity,
        labelFaceCamera: src.labelFaceCamera,
        labelFixedSize: src.labelFixedSize,
        labelFontFamily: src.labelFontFamily,
        labelRenderMode: src.labelRenderMode,
        labelShowBorder: src.labelShowBorder,
        labelBorderWidth: src.labelBorderWidth,
        labelBorderColor: src.labelBorderColor,
        mediaUrl: src.mediaUrl,
        imageUrl: src.imageUrl,
        videoUrl: src.videoUrl,
        webUrl: src.webUrl,
        mediaWidth: src.mediaWidth,
        mediaAspect: src.mediaAspect,
        mediaAutoplay: src.mediaAutoplay,
        mediaLoop: src.mediaLoop,
        mediaMuted: src.mediaMuted,
        uiX: src.uiX,
        uiY: src.uiY,
        uiWidth: src.uiWidth,
        uiHeight: src.uiHeight,
        uiRotation: src.uiRotation,
        autoRotate: src.autoRotate,
        rotateSpeed: src.rotateSpeed,
        rotateAxis: src.rotateAxis,
        floatAnim: src.floatAnim,
        floatRange: src.floatRange,
        floatSpeed: src.floatSpeed,
        blink: src.blink,
        blinkSpeed: src.blinkSpeed,
        blinkMin: src.blinkMin,
        interactions: src.interactions,
        wsKey: src.wsKey,
        bindProp: src.bindProp,
        bindTransform: src.bindTransform,
        bindScale: src.bindScale,
        bindOffset: src.bindOffset,
        positionBindings: src.positionBindings ? JSON.parse(JSON.stringify(src.positionBindings)) : undefined,
        realTimeValue: src.realTimeValue,
        dataFormat: src.dataFormat,
      }

      if (src.type === '2dComponent') {
        options.source2D = src.source2D
        options.label2D = src.label2D
        options.is2DComponent = true
      }

      this.addObject(options)
    },

    /** 选中对象 */
    selectObject(objectId) {
      this.selectedObjectId = objectId
      // 保存选中时的快照，供属性编辑撤销使用
      this._propSnapshot = cloneSceneObjects(this.sceneObjects)
      this._propSnapshotTaken = false
    },

    // ========================
    // ========================

    /** 推入历史快照 */
    pushHistory() {
      const snapshot = cloneSceneObjects(this.sceneObjects)
      this.history.past.push(snapshot)
      if (this.history.past.length > this.history.maxHistory) {
        this.history.past.shift()
      }
      // 新操作后清空 future
      this.history.future = []
    },

    commitSceneChange(mutator) {
      this.pushHistory()
      const nextObjects = cloneSceneObjects(this.sceneObjects)
      mutator.call(this, nextObjects)
      this.sceneObjects = nextObjects
      this.isCharge = true
    },

    normalize2DLayerZIndex(list) {
      var zIndex = 1
      list.forEach(function(obj) {
        if (!obj || obj.type !== '2dComponent' || !obj.source2D) return
        if (!obj.source2D.style) obj.source2D.style = {}
        obj.source2D.style.zIndex = zIndex++
      })
    },

    getObjectActions(obj) {
      if (!obj || !obj.interactions) return []
      return Object.keys(obj.interactions)
        .filter(eventType => {
          const interaction = obj.interactions[eventType]
          return interaction && interaction.enabled && interaction.type && interaction.type !== 'none'
        })
        .map(eventType => buildActionFromInteraction(eventType, obj.interactions[eventType]))
    },

    syncObjectActionCompat(obj) {
      if (!obj) return
      this.$set(obj, 'action', this.getObjectActions(obj))
    },

    syncAllObjectActions() {
      this.sceneObjects.forEach(obj => {
        this.syncObjectActionCompat(obj)
      })
    },

    flush2DStoreSelectionToScene() {
      const selectedComponent = store && store.state && store.state.ISMDisPlayEditorTool
        ? store.state.ISMDisPlayEditorTool.selectedComponent
        : null
      if (!selectedComponent || !selectedComponent.identifier) return
      const obj = this.sceneObjects.find(item =>
        item && item.type === '2dComponent' && item.source2D && item.source2D.identifier === selectedComponent.identifier
      )
      if (!obj) return
      this.$set(obj, 'source2D', JSON.parse(JSON.stringify(selectedComponent)))
      this.sync2DWrapperFields(obj)
    },

    /** 撤销 */
    undo() {
      if (this.history.past.length === 0) return
      // 保存当前状态到 future
      this.history.future.push(cloneSceneObjects(this.sceneObjects))
      const prev = this.history.past.pop()
      this.sceneObjects = cloneSceneObjects(prev)
      if (this.selectedObjectId && !this.sceneObjects.find(o => o.id === this.selectedObjectId)) {
        this.selectedObjectId = null
      }
    },

    /** 重做 */
    redo() {
      if (this.history.future.length === 0) return
      this.history.past.push(cloneSceneObjects(this.sceneObjects))
      const next = this.history.future.pop()
      this.sceneObjects = cloneSceneObjects(next)
    },

    // ========================
    //  事件处理
    // ========================

    /** 从工具箱添加对象 */
    handleAddObject(item) {
      if (!item) return

      if (item.is2DComponent || (item.info && item.info.type)) {
        const source2D = item.source2D || (item.info ? JSON.parse(JSON.stringify(item.info)) : {})
        source2D.style = source2D.style || {}
        source2D.style.position = source2D.style.position || {}

        if (item.base && item.base.info && item.base.info.style && item.base.info.style.position) {
          const basePos = item.base.info.style.position
          if (!source2D.style.position.w && basePos.w) {
            source2D.style.position.w = basePos.w
          }
          if (!source2D.style.position.h && basePos.h) {
            source2D.style.position.h = basePos.h
          }
        }

        const label = item.text || item.label || item.name || item.typeName || (item.info && item.info.type) || '2D Component'

        this.addObject({
          type: '2dComponent',
          options: {
            ...item,
            source2D: source2D,
            label2D: label,
            is2DComponent: true,
            typeName: item.info && item.info.type ? item.info.type : '2DComponent'
          }
        })
        return
      }
      if (item.type === 'gltf' && item.modelPath) {
        const obj = this.addObject({
          type: 'gltf',
          options: {
            ...item,
            label: item.name || item.label || 'GLTF Model',
            modelPath: item.modelPath,
            isExternalModel: true,
            materialOverridden: false,
            fitSize: item.fitSize || 2
          }
        })
        this.$nextTick(() => {
          if (this.$refs.sceneCanvas && this.$refs.sceneCanvas.loadGLTFModelUrl) {
            this.$refs.sceneCanvas.loadGLTFModelUrl(item.modelPath, obj.id)
          }
        })
        return
      }
      if (item.type === 'gltf-buffer' && item.buffer) {
        const obj = this.addObject({
          type: 'gltf',
          options: {
            ...item,
            label: item.name || '本地模型',
            isExternalModel: true,
            materialOverridden: false,
            fitSize: item.fitSize || 2,
          }
        })
        this.$nextTick(() => {
          if (this.$refs.sceneCanvas && this.$refs.sceneCanvas.loadGLTFModelBuffer) {
            this.$refs.sceneCanvas.loadGLTFModelBuffer(item.buffer, obj.id)
          }
        })
        return
      }
      this.addObject({ type: item.type, options: item })
    },

    onDragStart(item, e) {
      this._dragItem = item
      if (e && e.dataTransfer) {
        e.dataTransfer.effectAllowed = 'copy'
        e.dataTransfer.setData('text/plain', item.type || '')
        e.dataTransfer.setData('application/x-ism3d-object', JSON.stringify(item))
      }
    },

    /** 画布区域放置对象（使用鼠标在 3D 空间的实际位置） */
    handleDropObject({ type, item: payloadItem, x, y, z, screenX, screenY }) {
      const item = payloadItem || this._dragItem
      if (!item) return
      this._dragItem = null
      const dropX = Number.isFinite(Number(x)) ? Number(x) : 0
      const dropY = Number.isFinite(Number(y)) ? Number(y) : 0
      const dropZ = Number.isFinite(Number(z)) ? Number(z) : 0
      if (item.snapshot) {
        const snapshot = JSON.parse(JSON.stringify(item.snapshot))
        snapshot.x = dropX
        snapshot.y = dropY
        snapshot.z = dropZ
        snapshot.locked = false
        this.addObject({ type: snapshot.type, options: snapshot })
        return
      }
      if (item.type === 'gltf' && item.modelPath) {
        const obj = this.addObject({
          type: 'gltf',
          options: {
            ...item,
            label: item.name || item.label || 'GLTF Model',
            x: dropX,
            y: dropY,
            z: dropZ,
            modelPath: item.modelPath,
            isExternalModel: true,
            materialOverridden: false,
            fitSize: item.fitSize || 2
          }
        })
        this.$nextTick(() => {
          if (this.$refs.sceneCanvas && this.$refs.sceneCanvas.loadGLTFModelUrl) {
            this.$refs.sceneCanvas.loadGLTFModelUrl(item.modelPath, obj.id)
          }
        })
        return
      }
      if (item.type === 'gltf-buffer' && item.buffer) {
        const obj = this.addObject({
          type: 'gltf',
          options: {
            ...item,
            label: item.name || '本地模型',
            x: dropX,
            y: dropY,
            z: dropZ,
            isExternalModel: true,
            materialOverridden: false,
            fitSize: item.fitSize || 2,
          }
        })
        this.$nextTick(() => {
          if (this.$refs.sceneCanvas && this.$refs.sceneCanvas.loadGLTFModelBuffer) {
            this.$refs.sceneCanvas.loadGLTFModelBuffer(item.buffer, obj.id)
          }
        })
        return
      }
      if (item.is2DComponent || (item.info && item.info.type)) {
        const source2D = item.source2D || (item.info ? JSON.parse(JSON.stringify(item.info)) : {})

        source2D.style = source2D.style || {}
        source2D.style.position = source2D.style.position || {}

        source2D.style.position.x = Math.max(0, Math.round(screenX || 0))
        source2D.style.position.y = Math.max(0, Math.round(screenY || 0))

        if (item.base && item.base.info && item.base.info.style && item.base.info.style.position) {
          const basePos = item.base.info.style.position
          if (!source2D.style.position.w && basePos.w) {
            source2D.style.position.w = basePos.w
          }
          if (!source2D.style.position.h && basePos.h) {
            source2D.style.position.h = basePos.h
          }
        }

        const label = item.text || item.label || item.name || item.typeName || (item.info && item.info.type) || '2D Component'
        this.addObject({
          type: '2dComponent',
          options: {
            ...item,
            source2D: source2D,
            label2D: label,
            is2DComponent: true,
            typeName: item.info && item.info.type ? item.info.type : '2DComponent',
            x: 0,
            y: 0,
            z: 0
          }
        })
        return
      }
      // 使用鼠标放置的实际坐标，z 默认放在地面上
      const isUiBasic = item.type === 'uiLabel' || item.type === 'uiImage' || item.type === 'webEmbed'
      this.addObject({
        type: type,
        options: {
          ...item,
          label2D: item.label || item.name || item.typeName,
          x: isUiBasic ? 0 : dropX,
          y: isUiBasic ? 0 : dropY,
          z: isUiBasic ? 0 : dropZ,
          uiX: isUiBasic ? Math.max(0, Math.round(screenX || 40)) : item.uiX,
          uiY: isUiBasic ? Math.max(0, Math.round(screenY || 40)) : item.uiY
        }
      })
    },

    /** 从画布选中对象 */
    handleObjectSelected(objectId) {
      this.selectObject(objectId)
    },

    handleObjectTransformed({ objectId, transform }) {
      if (!transform || Object.keys(transform).length === 0) {
        return
      }
      this.pushHistory()
      this.updateObject(objectId, transform)
    },

    handleMoveUp(objectId) {
      var idx = this.sceneObjects.findIndex(function(o) { return o.id === objectId })
      if (idx > 0) {
        this.pushHistory()
        var arr = this.sceneObjects.slice()
        var item = arr.splice(idx, 1)[0]
        arr.splice(idx - 1, 0, item)
        this.normalize2DLayerZIndex(arr)
        this.sceneObjects = arr
      }
    },

    handleMoveDown(objectId) {
      var idx = this.sceneObjects.findIndex(function(o) { return o.id === objectId })
      if (idx < this.sceneObjects.length - 1) {
        this.pushHistory()
        var arr = this.sceneObjects.slice()
        var item = arr.splice(idx, 1)[0]
        arr.splice(idx + 1, 0, item)
        this.normalize2DLayerZIndex(arr)
        this.sceneObjects = arr
      }
    },

    handleToggleLock(objectId) {
      var obj = this.sceneObjects.find(function(o) { return o.id === objectId })
      if (obj) {
        this.pushHistory()
        obj.locked = !obj.locked
        this.sceneObjects = this.sceneObjects.slice()
      }
    },

    /** 从层级树选中 */
    handleOverlayEditStart() {
      this.pushHistory()
    },

    handleOverlayUpdated({ objectId }) {
      const obj = this.sceneObjects.find(item => item.id === objectId)
      if (!obj) return
      this.sceneObjects = [...this.sceneObjects]
      this.selectedObjectId = objectId
    },

    handleSelectObject(objectId) {
      this.selectObject(objectId)
    },

    /** 从层级树删除 */
    handleDeleteObject(objectId) {
      this.removeObject(objectId)
    },

    /** 删除选中对象 */
    handleDeleteSelected() {
      if (this.selectedObjectId) {
        this.removeObject(this.selectedObjectId)
      }
    },

    /** 聚焦选中对象 */
    handleFocusSelected() {
      var canvas = this.$refs.sceneCanvas
      if (canvas && canvas.focusSelected) {
        canvas.focusSelected()
      } else if (canvas && canvas.focusOn && this.selectedObject) {
        canvas.focusOn(this.selectedObject)
      }
    },

    /** 显示全景 */
    handleFrameAll() {
      var canvas = this.$refs.sceneCanvas
      if (canvas && canvas.frameAll && canvas.__3d) {
        canvas.frameAll(this.sceneObjects, canvas.__3d.meshMap)
      }
    },

    /** 层级树更新对象 */
    handleTreeUpdate({ id, changes }) {
      this.updateObject(id, changes)
    },

    /** 聚焦对象 */
    handleFocusObject(id) {
      this.selectObject(id)
    },

    /** 复制 */
    handleDuplicateObject(id) {
      this.duplicateObject(id)
    },

    handlePropChange() {
      var sel = this.selectedObject
      console.log('[Editor] handlePropChange selId=' + (sel && sel.id) + ' selType=' + (sel && sel.type) + ' selColor=' + (sel && sel.color) + ' selFontSize=' + (sel && sel.fontSize) + ' selFontFamily=' + (sel && sel.labelFontFamily))
      if (this._propSnapshot && !this._propSnapshotTaken) {
        this.history.past.push(this._propSnapshot)
        if (this.history.past.length > this.history.maxHistory) {
          this.history.past.shift()
        }
        this.history.future = []
        this._propSnapshotTaken = true
      }

      this.sceneObjects = [...this.sceneObjects]

      if (this.$refs.sceneCanvas && this.$refs.sceneCanvas.syncMeshes) {
        this.$refs.sceneCanvas.syncMeshes(this.sceneObjects)
      }

      // 实时更新场景中的文字对象
      if (this.selectedObject && this.$refs.sceneCanvas && this.$refs.sceneCanvas.updateTextObject) {
        this.$refs.sceneCanvas.updateTextObject(this.selectedObject)
      }
    },

    handleSceneSync() {
      this.sceneObjects.forEach(obj => this.sync2DWrapperFields(obj))
      this.sceneObjects = [...this.sceneObjects]
    },
    sync2DWrapperFields(obj) {
      if (!obj || obj.type !== '2dComponent' || !obj.source2D) return
      obj.action = Array.isArray(obj.source2D.action) ? JSON.parse(JSON.stringify(obj.source2D.action)) : []
      if (obj.source2D.style && obj.source2D.style.visible !== undefined) {
        obj.visible = obj.source2D.style.visible !== 0 && obj.source2D.style.visible !== false
      }
      if (obj.source2D.animate) {
        obj.animate = JSON.parse(JSON.stringify(obj.source2D.animate))
      }
      if (obj.source2D.name) {
        obj.name = obj.source2D.name
      }
      if (obj.source2D.text && !obj.label2D) {
        obj.label2D = obj.source2D.text
      }
    },

    sync2DStyleFromWrapperFields(obj, changes) {
      if (!obj || obj.type !== '2dComponent' || !obj.source2D || !changes) return
      if (!obj.source2D.style) obj.source2D.style = {}
      if (changes.visible !== undefined) {
        obj.source2D.style.visible = changes.visible ? 1 : 0
      }
    },

    handleTransformChange() {
      this.pushHistory()
    },

    handleTextChange() {
      this.pushHistory()
    },

    /** 创建流动管道 */
    handleCreateFlowPipe(pipeData) {
      var id = 'flowPipe_' + Date.now()
      var newPipe = {
        id: id,
        type: 'flowPipe',
        name: '流动管道',
        points: pipeData.points,
        color: pipeData.color || '#f4ead8',
        radius: pipeData.radius || 0.1,
        flowSpeed: pipeData.flowSpeed || 1.0,
        highlightColor: pipeData.highlightColor || '#ff6a00',
        flowDashLength: pipeData.flowDashLength || 3,
        x: 0,
        y: 0,
        z: 0,
        sx: 1,
        sy: 1,
        sz: 1,
        rx: 0,
        ry: 0,
        rz: 0,
        visible: true,
        locked: false
      }

      // 先记录快照，再添加新管道
      this.pushHistory()
      this.sceneObjects.push(newPipe)
      this.selectObject(id)
    },

    /** 开始绘制管道 */
    handleDrawingPipeStart() {
      this.$message.info('开始绘制管道：左键点击添加点，双击结束')
    },

    /** 取消绘制管道 */
    handleDrawingPipeCancel() {
      this.$message.info('已取消绘制管道')
    },

    handleModelDrop(e) {
      e.preventDefault()
      var files = e.dataTransfer.files
      if (files && files.length > 0) {
        this.loadModelFile(files[0])
      }
    },

    /** 选择并加载模型文件 */
    handleModelSelect(e) {
      var file = e.target.files[0]
      if (file) {
        this.loadModelFile(file)
        // 清空 input，允许重复选择同一个文件
        e.target.value = ''
      }
    },

    openModelFileSelector() {
      if (this.$refs.modelFileInput) {
        this.$refs.modelFileInput.click()
      }
    },

    openBgModelFileSelector() {
      if (this.$refs.bgModelFileInput) {
        this.$refs.bgModelFileInput.click()
      }
    },

    /** 加载模型文件 */
    async uploadModelFileToServer(file) {
      var formData = new FormData()
      formData.append('file', file)
      const res = await request(SYSTEMIMAGEUPLOAD, METHOD.POST, formData)
      const result = res && res.data ? res.data : {}
      console.log('模型上传响应:', result)
      if (result.Code !== 200 && result.Code !== 2002) {
        throw new Error(result.Message || '模型上传失败')
      }
      var path = result.Path || result.path || ''
      console.log('模型上传成功，路径:', path)
      return path
    },

    async loadModelFile(file, isBackground) {
      if (!file) return
      var fileName = file.name.toLowerCase()
      if (!fileName.endsWith('.gltf') && !fileName.endsWith('.glb')) {
        message.error('不支持的文件格式，请选择 .gltf 或 .glb 文件')
        return
      }

      var id = 'model_' + Math.random().toString(36).substr(2, 9)
      var obj = {
        id: id,
        type: 'gltf',
        name: file.name.replace(/\.[^/.]+$/, ''),
        typeName: 'GLTFModel',
        x: 0, y: 0, z: 0,
        rx: 0, ry: 0, rz: 0,
        sx: 1, sy: 1, sz: 1,
        // 不设置默认颜色，保留模型原始材质
        opacity: isBackground ? 0.5 : undefined,
        visible: true,
        locked: !!isBackground,
        isExternalModel: true,
        isBackground: !!isBackground,
        modelPath: '',
        showShadow: false,
        // 标记用户是否主动修改过材质，未修改时重载保留原始材质/贴图
        materialOverridden: false
      }

      try {
        obj.modelPath = await this.uploadModelFileToServer(file)
        console.log('模型上传成功，路径:', obj.modelPath)
      } catch (err) {
        message.error(err.message || '模型上传失败')
        console.error('模型上传失败:', err)
        return
      }

      var self = this
      this.pushHistory()
      this.sceneObjects.push(obj)
      if (!isBackground) {
        this.selectedObjectId = id
      }

      var reader = new FileReader()
      reader.onload = function(e) {
        var arrayBuffer = e.target.result
        if (self.$refs.sceneCanvas && self.$refs.sceneCanvas.loadGLTFModelBuffer) {
          self.$refs.sceneCanvas.loadGLTFModelBuffer(arrayBuffer, id, function(loadedId, mesh) {
            if (isBackground && mesh) {
              mesh.traverse(function(child) {
                if (child.isMesh && child.material) {
                  child.material.transparent = true
                  child.material.opacity = 0.5
                  child.material.needsUpdate = true
                }
              })
              self.showGrid = false
              if (self.$refs.sceneCanvas && self.$refs.sceneCanvas.toggleGrid) {
                self.$refs.sceneCanvas.toggleGrid(false)
              }
            }
            message.success('模型加载成功：' + file.name)
            self.pushHistory()
          })
        } else {
          message.error('场景画布未初始化')
        }
      }
      reader.readAsArrayBuffer(file)
    },

    handleBgModelDrop(e) {
      e.preventDefault()
      this.bgModelDragOver = false
      var files = e.dataTransfer.files
      if (files && files.length > 0) {
        this.loadModelFile(files[0], true)
      }
    },

    /** 选择并加载背景模型文件 */
    handleBgModelSelect(e) {
      var file = e.target.files[0]
      this.bgModelDragOver = false
      if (file) {
        this.loadModelFile(file, true)
        // 清空 input，允许重复选择同一个文件
        e.target.value = ''
      }
    },

    /** 更新网格和场景附加设置 */
    updateGrid() {
      if (this.$refs.sceneCanvas && this.$refs.sceneCanvas.updateGrid) {
        this.$refs.sceneCanvas.updateGrid(this.gridSettings)
      }
      if (this.$refs.sceneCanvas && this.$refs.sceneCanvas.applySceneExtras) {
        this.$refs.sceneCanvas.applySceneExtras(this.sceneExtras, this.gridSettings)
      }
    },

    normalizeGridSettings(settings) {
      const normalized = Object.assign({}, settings || {})
      const colorFields = [
        ['backgroundColor', '#ffffff'],
        ['backgroundColor2', '#dbefff'],
        ['colorGrid', '#444444'],
        ['colorCenterLine', '#111111']
      ]
      colorFields.forEach(([key, fallback]) => {
        if (normalized[key] !== undefined) {
          normalized[key] = this.normColorValue(normalized[key], fallback)
        }
      })
      if (normalized.lightSettings) {
        const lightSettings = Object.assign({}, normalized.lightSettings)
        if (lightSettings.ambientColor !== undefined) {
          lightSettings.ambientColor = this.normColorValue(lightSettings.ambientColor, '#ffffff')
        }
        if (lightSettings.directionalColor !== undefined) {
          lightSettings.directionalColor = this.normColorValue(lightSettings.directionalColor, '#ffffff')
        }
        normalized.lightSettings = lightSettings
      }
      return normalized
    },

    patchGridSettings(patch) {
      this.gridSettings = this.normalizeGridSettings(Object.assign({}, this.gridSettings, patch || {}))
      if (patch && patch.size !== undefined) {
        this.gridSize = patch.size || 10
      }
      if (patch && patch.colorGrid !== undefined) {
        this.gridColor = this.gridSettings.colorGrid || '#444444'
      }
      this.updateGrid()
    },

    applySceneSettings(sceneSettings) {
      if (!sceneSettings) return
      this.gridSettings = this.normalizeGridSettings(mergeSceneSettings(this.gridSettings, sceneSettings))
      this.gridSize = this.gridSettings.size || 10
      this.gridColor = this.gridSettings.colorGrid || '#444444'
      if (sceneSettings.showGrid !== undefined) this.showGrid = sceneSettings.showGrid
    },

    buildSceneSettingsPayload() {
      const defaults = defaultGridSettings()
      return {
        backgroundColor: this.gridSettings.backgroundColor || defaults.backgroundColor,
        backgroundColor2: this.gridSettings.backgroundColor2 || defaults.backgroundColor2,
        backgroundImage: this.gridSettings.backgroundImage || '',
        backgroundMode: this.gridSettings.backgroundMode || defaults.backgroundMode,
        sceneEnvironmentPreset: this.gridSettings.sceneEnvironmentPreset || defaults.sceneEnvironmentPreset,
        environmentPreset: this.gridSettings.environmentPreset || defaults.environmentPreset,
        skyboxEnabled: this.gridSettings.skyboxEnabled !== false,
        skyboxPreset: this.gridSettings.skyboxPreset || defaults.skyboxPreset,
        skyboxImage: this.gridSettings.skyboxImage || '',
        skyboxHdri: this.gridSettings.skyboxHdri || '',
        groundEnabled: this.gridSettings.groundEnabled !== false,
        groundStyle: this.gridSettings.groundStyle || defaults.groundStyle,
        sectionLength: this.gridSettings.sectionLength || defaults.sectionLength,
        sectionWidth: this.gridSettings.sectionWidth || defaults.sectionWidth,
        sectionHeight: this.gridSettings.sectionHeight || defaults.sectionHeight,
        tunnelRadius: this.gridSettings.tunnelRadius || defaults.tunnelRadius,
        lightingPreset: this.gridSettings.lightingPreset || defaults.lightingPreset,
        floorReflection: this.gridSettings.floorReflection || defaults.floorReflection,
        enhanceDepth: !!this.gridSettings.enhanceDepth,
        modelOptimize: !!this.gridSettings.modelOptimize,
        gridSize: this.gridSettings.size || defaults.size,
        gridDivisions: this.gridSettings.divisions || defaults.divisions,
        gridColorCenterLine: this.gridSettings.colorCenterLine || defaults.colorCenterLine,
        gridColorGrid: this.gridSettings.colorGrid || defaults.colorGrid,
        showGrid: this.showGrid,
        lightSettings: this.gridSettings.lightSettings || defaults.lightSettings,
      }
    },

    emit2DComponentRuntimeEvents(obj, realData, currentData, realValue) {
      if (!obj || obj.type !== '2dComponent' || !obj.source2D) return
      const identifier = obj.source2D.identifier
      if (!identifier) return

      const activeList = Array.isArray(obj.source2D.active) ? obj.source2D.active : []
      activeList.forEach((activeItem, activeIndex) => {
        const activeCondition = activeItem && activeItem.condition ? activeItem.condition : null
        if (!activeCondition || !activeCondition.dataID) return
        const deviceMatched = activeCondition.isBandDevice ? realData.DeviceUuid === activeCondition.deviceSN : true
        const dataMatched = activeCondition.dataID === currentData.ModelDataUuid || activeCondition.dataID === currentData.Uuid
        if (!deviceMatched || !dataMatched) return

        let activeResult = currentData.Value
        if (activeItem.isExpression) {
          const operatorValue = parseFloat(activeCondition.OperatorValue)
          const operatorMaxValue = parseFloat(activeCondition.OperatorMaxValue)
          if (!isNaN(realValue) && !isNaN(operatorValue)) {
            activeResult = evaluateCondition(activeCondition.operator, realValue, operatorValue, operatorMaxValue)
          }
        }

        this.$EventBus.$emit(identifier + 'activeEvent', {
          ID: activeItem.id,
          DeviceSN: activeCondition.deviceSN,
          dataID: activeCondition.dataID,
          index: activeIndex,
          result: activeResult
        })
      })

      const moveConfig = obj.source2D.animate && obj.source2D.animate.move ? obj.source2D.animate.move : null
      const moveSelected = obj.source2D.animate && Array.isArray(obj.source2D.animate.selected) ? obj.source2D.animate.selected : []
      if (moveConfig && moveSelected.includes('animateMove') && !isNaN(realValue)) {
        ['x', 'y'].forEach((axis) => {
          const moveCondition = moveConfig[axis]
          if (!moveCondition || !moveCondition.dataID) return
          const deviceMatched = moveCondition.isBandDevice ? realData.DeviceUuid === moveCondition.deviceSN : true
          const dataMatched = moveCondition.dataID === currentData.ModelDataUuid || moveCondition.dataID === currentData.Uuid
          if (!deviceMatched || !dataMatched) return
          if (obj.source2D.style && obj.source2D.style.position) {
            obj.source2D.style.position[axis] = realValue
          }
          this.$EventBus.$emit(identifier + 'animateMove', {
            ID: axis === 'x' ? 'animateMoveX' : 'animateMoveY',
            DeviceSN: moveCondition.deviceSN,
            dataID: moveCondition.dataID,
            result: realValue
          })
        })
      }
    },

    // ========================
    //  数据绑定（动画控制）
    // ========================
    DealWithUpdateData(realData) {
      if (!realData || !realData.Data || !Array.isArray(realData.Data)) {
        return;
      }

      for (let k = 0, Datalen = realData.Data.length; k < Datalen; k++) {
        for (let j = 0, objLen = this.sceneObjects.length; j < objLen; j++) {
          const obj = this.sceneObjects[j];
          const currentData = realData.Data[k];

          this.applyGLTFAnimationCondition(obj, realData, currentData)

          if (!obj.animate || !obj.animate.condition) {
            continue;
          }

          const condition = obj.animate.condition;
          const selectAnimate = obj.animate.selected;
          const identifier = obj.source2D && obj.source2D.identifier ? obj.source2D.identifier : '';
          const realValue = parseFloat(currentData.Value);

          if (obj.type === 'dataText') {
            const wsMatched = obj.wsKey && (obj.wsKey === currentData.ModelDataUuid || obj.wsKey === currentData.Uuid || obj.wsKey === currentData.Name);
            const dataMatched = condition && condition.dataID && (condition.dataID === currentData.ModelDataUuid || condition.dataID === currentData.Uuid);
            if (wsMatched || dataMatched) {
              this.updateObject(obj.id, { realTimeValue: currentData.Value })
            }
          }
          this.applyPositionBindings(obj, realData, currentData)
          this.applyGenericDataBinding(obj, realData, currentData)
          this.emit2DComponentRuntimeEvents(obj, realData, currentData, realValue)

          // 绑定到指定设备
          if (condition.isBandDevice) {
            let isStart = false;

            if (realData.DeviceUuid == condition.deviceSN) {
              if ((condition.dataID == currentData.ModelDataUuid) ||
                  (condition.dataID == currentData.Uuid)) {

                const RealValue = parseFloat(currentData.Value);

                if (obj.animate.isExpression || condition.operator) {
                  const OperatorValue = parseFloat(condition.OperatorValue);

                  switch (condition.operator) {
                    case "==": isStart = (RealValue == OperatorValue); break;
                    case ">": isStart = (RealValue > OperatorValue); break;
                    case ">=": isStart = (RealValue >= OperatorValue); break;
                    case "<": isStart = (RealValue < OperatorValue); break;
                    case "<=": isStart = (RealValue <= OperatorValue); break;
                    case "!=": isStart = (RealValue != OperatorValue); break;
                    case "<>": {
                      const OperatorMaxValue = parseFloat(condition.OperatorMaxValue);
                      isStart = (RealValue >= OperatorValue) && (RealValue <= OperatorMaxValue);
                      break;
                    }
                    case "<!>": {
                      const OperatorMaxValue = parseFloat(condition.OperatorMaxValue);
                      isStart = (RealValue <= OperatorValue) || (RealValue >= OperatorMaxValue);
                      break;
                    }
                  }
                }

                if (!obj.animate.isExpression && !condition.operator) {
                  isStart = true;
                }
                if (selectAnimate && selectAnimate.length > 0) {
                  const changes = {};

                  selectAnimate.forEach(animId => {
                    switch (animId) {
                      case "autoRotate": changes.autoRotate = isStart; break;
                      case "floatAnim": changes.floatAnim = isStart; break;
                      case "blink": changes.blink = isStart; break;
                      case "visible": changes.visible = isStart; break;
                      case "flowAnim": changes.flowAnim = isStart; break;
                    }
                  });

                  // 更新对象属性（触发响应式）
                  this.updateObject(obj.id, changes);
                  if (obj.type === '2dComponent' && identifier) {
                    this.$EventBus.$emit(identifier + 'animateEvent', isStart);
                  }
                }
              }
            }
          }
          else if (realData.DeviceUuid == this.SelectDeviceUuid) {
            let isStart = false;

            if ((condition.dataID == currentData.ModelDataUuid) ||
                (condition.dataID == currentData.Uuid)) {

              const RealValue = parseFloat(currentData.Value);

              if (obj.animate.isExpression || condition.operator) {
                const OperatorValue = parseFloat(condition.OperatorValue);

                switch (condition.operator) {
                  case "==": isStart = (RealValue == OperatorValue); break;
                  case ">": isStart = (RealValue > OperatorValue); break;
                  case ">=": isStart = (RealValue >= OperatorValue); break;
                  case "<": isStart = (RealValue < OperatorValue); break;
                  case "<=": isStart = (RealValue <= OperatorValue); break;
                  case "!=": isStart = (RealValue != OperatorValue); break;
                  case "<>": {
                    const OperatorMaxValue = parseFloat(condition.OperatorMaxValue);
                    isStart = (RealValue >= OperatorValue) && (RealValue <= OperatorMaxValue);
                    break;
                  }
                  case "<!>": {
                    const OperatorMaxValue = parseFloat(condition.OperatorMaxValue);
                    isStart = (RealValue <= OperatorValue) || (RealValue >= OperatorMaxValue);
                    break;
                  }
                }
              }

              if (!obj.animate.isExpression && !condition.operator) {
                isStart = true;
              }
              if (selectAnimate && selectAnimate.length > 0) {
                const changes = {};

                selectAnimate.forEach(animId => {
                  switch (animId) {
                    case "autoRotate": changes.autoRotate = isStart; break;
                    case "floatAnim": changes.floatAnim = isStart; break;
                    case "blink": changes.blink = isStart; break;
                    case "visible": changes.visible = isStart; break;
                    case "flowAnim": changes.flowAnim = isStart; break;
                  }
                });

                // 更新对象属性（触发响应式）
                this.updateObject(obj.id, changes);
                if (obj.type === '2dComponent' && identifier) {
                  this.$EventBus.$emit(identifier + 'animateEvent', isStart);
                }
              }
            }
          }
        }
      }
    },

    // ========================
    // ========================

    /** 相机模式切换 */
    handleCameraChange(mode) {
      this.cameraMode = mode
    },

    /** 閫傚簲绐楀彛 */
    handleZoomFit() {
    },

    /** 保存场景 */
    handleSave() {
      this.saveSceneToServer()
    },

    /** 导入场景 JSON */
    handleImport() {
      const input = document.createElement('input')
      input.type = 'file'
      input.accept = '.json'
      input.onchange = (e) => {
        const file = e.target.files[0]
        if (!file) return
        const reader = new FileReader()
        reader.onload = (event) => {
          try {
            const data = parseImportedSceneFile(event.target.result)
            // 校验文件格式
            if (!data.objects || !Array.isArray(data.objects)) {
              message.error(this.$t('ISM3DEditor.invalidSceneFormat'))
              return
            }
            const valid = data.objects.every(obj => obj.id && obj.type)
            if (!valid) {
              message.error(this.$t('ISM3DEditor.invalidSceneObjectFields'))
              return
            }
            // 保存历史快照
            if (this.sceneObjects.length > 0) {
              this.pushHistory()
            }
            this.normalizeSceneObjects(data.objects)
            this.sceneObjects = data.objects
            this.selectedObjectId = null
            if (data.sceneSettings) {
              this.applySceneSettings(data.sceneSettings)
            }
            this.sceneExtras = Object.assign(createDefaultSceneExtras(), data.sceneExtras || {})
            message.success('已导入场景（' + data.objects.length + ' 个对象）')
            // 重建 GLTF 模型并应用背景
            this.$nextTick(() => {
              if (this.$refs.sceneCanvas && this.$refs.sceneCanvas.__3d) {
                this.restoreGLTFModels(data.objects)
                this.updateGrid()
                if (this.$refs.sceneCanvas.applySceneExtras) {
                  this.$refs.sceneCanvas.applySceneExtras(this.sceneExtras)
                }
                // 同步网格显隐状态到 SceneCanvas
                if (this.$refs.sceneCanvas.toggleGrid) {
                  this.$refs.sceneCanvas.toggleGrid(this.showGrid)
                }
              }
            })
          } catch (err) {
            message.error('场景文件解析失败：' + err.message)
          }
        }
        reader.readAsText(file)
      }
      input.click()
    },

    /** 导出场景 JSON */
    handleExport() {
      if (this.sceneObjects.length === 0) {
        message.warning('场景为空，无法导出')
        return
      }
      const data = {
        version: '1.0',
        generator: 'ISM3DEditor',
        exportTime: new Date().toISOString(),
        objectCount: this.sceneObjects.length,
        sceneSettings: this.buildSceneSettingsPayload(),
        sceneExtras: JSON.parse(JSON.stringify(this.sceneExtras || createDefaultSceneExtras())),
        objects: this.sceneObjects.map(obj => ({
          id: obj.id,
          type: obj.type,
          name: obj.name || '',
          typeName: obj.typeName || obj.type,
          icon: obj.icon || '',
          groupName: obj.groupName || '',
          intensity: obj.intensity !== undefined ? obj.intensity : undefined,
          distance: obj.distance !== undefined ? obj.distance : undefined,
          angle: obj.angle !== undefined ? obj.angle : undefined,
          penumbra: obj.penumbra !== undefined ? obj.penumbra : undefined,
          category: obj.category || '',
          x: obj.x !== undefined ? obj.x : 0,
          y: obj.y !== undefined ? obj.y : 0,
          z: obj.z !== undefined ? obj.z : 0,
          rx: obj.rx !== undefined ? obj.rx : 0,
          ry: obj.ry !== undefined ? obj.ry : 0,
          rz: obj.rz !== undefined ? obj.rz : 0,
          sx: obj.sx !== undefined ? obj.sx : 1,
          sy: obj.sy !== undefined ? obj.sy : 1,
          sz: obj.sz !== undefined ? obj.sz : 1,
          // GLTF 模型保留自身材质，普通对象使用编辑器材质字段
          color: (obj.type === 'text3d' || obj.type === 'textPlain3d' || obj.type === 'dataText' || obj.type === 'uiLabel') ? (obj.color || '#ffffff') : (obj.type === 'gltf' ? obj.color : (obj.color || '#4a90d9')),
          opacity: obj.opacity !== undefined ? obj.opacity : 1,
          metalness: obj.type === 'gltf' ? obj.metalness : (obj.metalness !== undefined ? obj.metalness : 0.3),
          roughness: obj.type === 'gltf' ? obj.roughness : (obj.roughness !== undefined ? obj.roughness : 0.7),
          wireframe: obj.type === 'gltf' ? obj.wireframe : (obj.wireframe || false),
          showShadow: !!obj.showShadow,
          visible: obj.visible !== false,
          locked: obj.locked || false,
          textContent: obj.textContent || '',
          fontSize: obj.fontSize || 16,
          textBgColor: obj.textBgColor || ((obj.type === 'text3d' || obj.type === 'textPlain3d' || obj.type === 'dataText' || obj.type === 'uiLabel') ? '#00000000' : '#000000'),
          labelOpacity: obj.labelOpacity !== undefined ? obj.labelOpacity : 1,
          textBgOpacity: obj.textBgOpacity !== undefined ? obj.textBgOpacity : ((obj.type === 'text3d' || obj.type === 'textPlain3d' || obj.type === 'dataText' || obj.type === 'uiLabel') ? 0 : 0.2),
          labelFaceCamera: obj.labelFaceCamera !== undefined ? obj.labelFaceCamera : true,
          labelFixedSize: obj.labelFixedSize !== undefined ? obj.labelFixedSize : true,
          labelFontFamily: obj.labelFontFamily || '系统默认',
          labelRenderMode: obj.labelRenderMode || 'component',
          labelShowBorder: obj.labelShowBorder !== undefined ? !!obj.labelShowBorder : false,
          labelBorderWidth: obj.labelBorderWidth !== undefined ? obj.labelBorderWidth : 1,
          labelBorderColor: obj.labelBorderColor || '#ffffff',
          autoRotate: obj.autoRotate || false,
          rotateSpeed: obj.rotateSpeed !== undefined ? obj.rotateSpeed : 1,
          rotateAxis: obj.rotateAxis || 'y',
          floatAnim: obj.floatAnim || false,
          floatRange: obj.floatRange !== undefined ? obj.floatRange : 0.15,
          floatSpeed: obj.floatSpeed !== undefined ? obj.floatSpeed : 2,
          blink: obj.blink || false,
          blinkSpeed: obj.blinkSpeed !== undefined ? obj.blinkSpeed : 6,
          blinkMin: obj.blinkMin !== undefined ? obj.blinkMin : 0.2,
          interactions: obj.interactions ? JSON.parse(JSON.stringify(obj.interactions)) : createDefaultInteractionConfig(),
          action: Array.isArray(obj.action) ? JSON.parse(JSON.stringify(obj.action)) : this.getObjectActions(obj),
          wsKey: obj.wsKey || '',
          bindProp: obj.bindProp || '',
          bindTransform: obj.bindTransform || 'direct',
          bindScale: obj.bindScale !== undefined ? obj.bindScale : 1,
          bindOffset: obj.bindOffset !== undefined ? obj.bindOffset : 0,
          positionBindings: obj.positionBindings ? JSON.parse(JSON.stringify(obj.positionBindings)) : undefined,
          realTimeValue: obj.realTimeValue !== undefined ? obj.realTimeValue : '',
          dataFormat: obj.dataFormat || '{value}',
          is2DComponent: !!obj.is2DComponent,
          label2D: obj.label2D || '',
          source2D: obj.source2D ? JSON.parse(JSON.stringify(obj.source2D)) : null,
          source2DMeta: obj.source2DMeta ? JSON.parse(JSON.stringify(obj.source2DMeta)) : null,
          animate: obj.animate ? JSON.parse(JSON.stringify(obj.animate)) : null,
          isExternalModel: obj.isExternalModel || false,
          isBackground: obj.isBackground || false,
          fitSize: obj.fitSize,
          modelPath: obj.modelPath || '',
          materialOverridden: !!obj.materialOverridden,
          emissive: obj.emissive || '#000000',
          textureData: obj.textureData || '',
          mediaUrl: obj.mediaUrl || '',
          imageUrl: obj.imageUrl || '',
          videoUrl: obj.videoUrl || '',
          webUrl: obj.webUrl || '',
          mediaWidth: obj.mediaWidth !== undefined ? obj.mediaWidth : undefined,
          mediaAspect: obj.mediaAspect !== undefined ? obj.mediaAspect : undefined,
          mediaAutoplay: obj.mediaAutoplay !== undefined ? obj.mediaAutoplay : true,
          mediaLoop: obj.mediaLoop !== undefined ? obj.mediaLoop : true,
          mediaMuted: obj.mediaMuted !== undefined ? obj.mediaMuted : true,
          uiX: obj.uiX !== undefined ? obj.uiX : undefined,
          uiY: obj.uiY !== undefined ? obj.uiY : undefined,
          uiWidth: obj.uiWidth !== undefined ? obj.uiWidth : undefined,
          uiHeight: obj.uiHeight !== undefined ? obj.uiHeight : undefined,
          uiRotation: obj.uiRotation !== undefined ? obj.uiRotation : 0,
          points: obj.type === 'flowPipe' ? (obj.points || []) : undefined,
          radius: obj.type === 'flowPipe' ? (obj.radius !== undefined ? obj.radius : 0.1) : undefined,
          flowSpeed: obj.type === 'flowPipe' ? (obj.flowSpeed !== undefined ? obj.flowSpeed : 1.0) : undefined,
          flowDirection: obj.type === 'flowPipe' ? (obj.flowDirection || 'forward') : undefined,
          flowAnim: obj.type === 'flowPipe' ? (obj.flowAnim !== undefined ? obj.flowAnim : true) : undefined,
          highlightColor: obj.type === 'flowPipe' ? (obj.highlightColor || '#ff6a00') : undefined,
          flowDashLength: obj.type === 'flowPipe' ? (obj.flowDashLength !== undefined ? obj.flowDashLength : 3) : undefined,
        }))
      }
      const encryptedData = sm4EncryptBase64(JSON.stringify(data))
      const blob = new Blob([encryptedData], { type: 'application/json' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      const timestamp = new Date().toISOString().replace(/[:.]/g, '-').substring(0, 19)
      a.download = 'scene-' + timestamp + '.json'
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      URL.revokeObjectURL(url)
      message.success('场景已导出（' + this.sceneObjects.length + ' 个对象）')
    },

    /** 撤销 */
    handleUndo() {
      this.undo()
    },

    /** 重做 */
    handleRedo() {
      this.redo()
    },

    /** 设置编辑模式 */
    handleSetMode(mode) {
      this.editMode = mode
    },

    /** 切换网格 */
    toggleGrid() {
      this.showGrid = !this.showGrid
      if (this.$refs.sceneCanvas && this.$refs.sceneCanvas.toggleGrid) {
        this.$refs.sceneCanvas.toggleGrid(this.showGrid)
      }
    },

    /** 切换线框模式 */
    toggleWireframe() {
      this.wireframeMode = !this.wireframeMode
    },

    /** 棰勮 */
    async handlePreview() {
      const uid = this.$route && this.$route.params && this.$route.params.uid
      if (uid) {
        this.flush2DStoreSelectionToScene()
        const saveSuccess = await this.saveSceneToServer({ silentSuccess: true })
        if (!saveSuccess) {
          return
        }
        const route = this.$router.resolve({
          path: '/DisPlay3DRunApp/' + uid,
          query: {
            pageId: this.currentPageId || this.SelectPagerID
          }
        })
        window.open(route.href, '_blank')
      } else {
        message.warning(this.$t('ISM3DEditor.missingAppIdPreview'))
      }
    },

    /** 切换网格设置面板 */
    toggleGridSettings() {
      this.showGridSettings = !this.showGridSettings
      this.showLightSettings = false
      this.showInspectPanel = false
      if (this.showGridSettings) {
        this.$nextTick(() => {
          var btns = Array.from(document.querySelectorAll('.editor-header .toolbar-btn'))
          var idx = btns.findIndex(function(b) { return b.textContent.indexOf('场景设置') !== -1 || b.innerText.indexOf('场景设置') !== -1 })
          if (idx >= 0) {
            this.dropdownPos.scene = { left: btns[idx].getBoundingClientRect().left + 'px' }
          }
        })
      }
    },

    /** 切换灯光设置面板 */
    toggleLightSettings() {
      this.showLightSettings = !this.showLightSettings
      this.showGridSettings = false
      this.showInspectPanel = false
      if (this.showLightSettings) {
        this.$nextTick(() => {
          var btns = Array.from(document.querySelectorAll('.editor-header .toolbar-btn'))
          var idx = btns.findIndex(function(b) { return b.textContent.indexOf('灯光') !== -1 || b.innerText.indexOf('灯光') !== -1 })
          if (idx >= 0) {
            this.dropdownPos.light = { left: btns[idx].getBoundingClientRect().left + 'px' }
          }
        })
      }
    },

    /** 灯光设置变更 */
    onLightSettingChange(key, value) {
      var ls = Object.assign({}, this.gridSettings.lightSettings || {})
      ls[key] = value
      this.patchGridSettings({ lightSettings: ls })
      if (this.$refs.sceneCanvas && this.$refs.sceneCanvas.applyLightSettings) {
        this.$refs.sceneCanvas.applyLightSettings(ls)
      }
    },

    toggleInspectPanel() {
      this.showInspectPanel = !this.showInspectPanel
      this.showGridSettings = false
      this.showLightSettings = false
      if (this.showInspectPanel) {
        this.$nextTick(() => {
          var btns = Array.from(document.querySelectorAll('.editor-header .toolbar-btn'))
          var idx = btns.findIndex(function(b) { return b.textContent.indexOf('检查') !== -1 || b.innerText.indexOf('检查') !== -1 })
          if (idx >= 0) {
            this.dropdownPos.inspect = { left: btns[idx].getBoundingClientRect().left + 'px' }
          }
        })
      }
    },

    repairTimelineTracks() {
      const ids = new Set(this.sceneObjects.map(item => item.id))
      const nextExtras = JSON.parse(JSON.stringify(this.sceneExtras || {}))
      nextExtras.timelineTracks = (nextExtras.timelineTracks || []).filter(item => item && item.objectId && ids.has(item.objectId))
      this.handleSceneExtrasUpdate(nextExtras)
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
          this.handleWorkbenchUpdateObject({
            id: obj.id,
            patch: {
              interactions: {
                ...(obj.interactions || {}),
                events: nextEvents
              }
            }
          })
        }
      })
    },

    showHelp() {
      this.showHelpDialog = true
    },

    showTemplates() {
      this.showTemplatesModal = true
    },

    /** 加载场景模板 */
    handleLoadTemplate(tpl) {
      if (!tpl || !tpl.objects) return
      const self = this
      // 如果当前场景非空，保存当前状态到历史
      if (self.sceneObjects.length > 0) {
        self.pushHistory()
      }
      const newObjects = tpl.objects.map(obj => {
        const newId = generateId()
        return { ...JSON.parse(JSON.stringify(obj)), id: newId }
      })
      self.normalizeSceneObjects(newObjects)
      self.sceneObjects = newObjects
      self.selectedObjectId = null
      self.showTemplatesModal = false
      self.showGrid = tpl.sceneSettings && tpl.sceneSettings.showGrid !== undefined ? !!tpl.sceneSettings.showGrid : true
      if (tpl.sceneSettings) {
        self.applySceneSettings(tpl.sceneSettings)
        self._savedSceneSettings = tpl.sceneSettings
      }
      self.pushHistory()
      message.success('已加载模板：' + tpl.name + '，' + tpl.objects.length + ' 个对象')
      self.$nextTick(() => {
        if (self.$refs.sceneCanvas && self.$refs.sceneCanvas.__3d) {
          self.updateGrid()
          if (self.$refs.sceneCanvas.toggleGrid) {
            self.$refs.sceneCanvas.toggleGrid(self.showGrid)
          }
          if (tpl.sceneSettings && self.$refs.sceneCanvas.restoreCamera) {
            self.$refs.sceneCanvas.restoreCamera(tpl.sceneSettings)
          }
        }
      })
    },

    closeHelp() {
      this.showHelpDialog = false
    },

    /** 网格大小变化 */
    onGridSizeChange(value) {
      this.patchGridSettings({ size: value || 10 })
    },

    /** 网格颜色变化 */
    onGridColorChange(e) {
      this.patchGridSettings({ colorGrid: e.target.value || '#444444' })
    },

    /** 网格分段变化 */
    onGridDivisionsChange(value) {
      this.patchGridSettings({ divisions: value || 20 })
    },

    onCenterLineColorChange(e) {
      this.patchGridSettings({ colorCenterLine: e.target.value || '#111111' })
    },

    onBackgroundModeChange(value) {
      const mode = value || 'solid'
      const patch = { backgroundMode: mode }
      if (mode !== 'image') {
        Object.assign(patch, getEnvironmentPresetColors(this.gridSettings.environmentPreset || 'sky'))
      }
      this.patchGridSettings(patch)
    },

    onEnvironmentPresetChange(value) {
      const preset = value || 'sky'
      const patch = Object.assign({ environmentPreset: preset }, getEnvironmentPresetColors(preset))
      if (this.gridSettings.backgroundMode === 'image') {
        patch.backgroundMode = 'solid'
        patch.backgroundImage = ''
      }
      this.patchGridSettings(patch)
    },

    onSceneEnvironmentPresetChange(value) {
      const preset = this.sceneEnvironmentPresets.find(item => item.key === value) || this.sceneEnvironmentPresets[0]
      if (!preset) return
      this.patchGridSettings(Object.assign({
        sceneEnvironmentPreset: preset.key
      }, preset.settings))
    },

    onSkyboxEnabledChange(value) {
      this.patchGridSettings({ skyboxEnabled: !!value })
    },

    onSkyboxPresetChange(value) {
      this.patchGridSettings({ skyboxPreset: value || 'horizon' })
    },

    onGroundEnabledChange(value) {
      this.patchGridSettings({ groundEnabled: !!value })
    },

    onGroundStyleChange(value) {
      this.patchGridSettings({ groundStyle: value || 'slate' })
    },

    onLightingPresetChange(value) {
      this.patchGridSettings({ lightingPreset: value || 'day' })
    },

    onFloorReflectionChange(value) {
      this.patchGridSettings({ floorReflection: value || 'none' })
    },

    onEnhanceDepthChange(value) {
      this.patchGridSettings({ enhanceDepth: !!value })
    },

    onModelOptimizeChange(value) {
      this.patchGridSettings({ modelOptimize: !!value })
    },

    openSystemImageSelector(target) {
      this._imageSelectTarget = target || 'background'
      if (this.$refs.systemImageModel) {
        this.$refs.systemImageModel.showModal(0)
      }
    },

    /** 选择背景图片 */
    onSelectBackgroundImage(url) {
      if (this._imageSelectTarget === 'skyboxImage') {
        this.patchGridSettings({ skyboxImage: url, skyboxPreset: 'customImage', skyboxEnabled: true })
        return
      }
      this.patchGridSettings({ backgroundImage: url, backgroundMode: 'image' })
    },

    handleResetCamera() {
      if (this.$refs.sceneCanvas && this.$refs.sceneCanvas.resetCamera) {
        this.$refs.sceneCanvas.resetCamera()
      }
    },

    // ========================
    // ========================

    handleGlobalKeyDown(e) {
      var key = e.key || ''
      var keyLower = key.toLowerCase()
      if ((e.ctrlKey || e.metaKey) && keyLower === 's') {
        e.preventDefault()
        this.handleSave()
      }
    },

    handleKeyDown(e) {
      if (e.defaultPrevented) return
      var key = e.key
      var keyLower = key.toLowerCase()
      var isInput = e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA' || e.target.tagName === 'SELECT'

      // Delete / Backspace - 删除选中对象
      if ((key === 'Delete' || key === 'Backspace') && this.selectedObjectId) {
        if (isInput) return
        e.preventDefault()
        this.handleDeleteSelected()
        return
      }

      // Ctrl+Z - 撤销
      if ((e.ctrlKey || e.metaKey) && keyLower === 'z' && !e.shiftKey) {
        e.preventDefault()
        this.undo()
        return
      }

      // Ctrl+Shift+Z / Ctrl+Y - 重做
      if ((e.ctrlKey || e.metaKey) && ((keyLower === 'z' && e.shiftKey) || keyLower === 'y')) {
        e.preventDefault()
        this.redo()
        return
      }

      // Ctrl+D - 复制
      if ((e.ctrlKey || e.metaKey) && keyLower === 'd') {
        e.preventDefault()
        if (this.selectedObjectId) {
          this.duplicateObject(this.selectedObjectId)
        }
        return
      }

      // Ctrl+S - 保存
      if ((e.ctrlKey || e.metaKey) && keyLower === 's') {
        e.preventDefault()
        this.handleSave()
        return
      }

      // W - 选择模式
      if (keyLower === 'w' && !e.ctrlKey && !e.metaKey) {
        if (isInput) return
        this.editMode = 'select'
        return
      }

      // E - 移动模式
      if (keyLower === 'e' && !e.ctrlKey && !e.metaKey) {
        if (isInput) return
        this.editMode = 'move'
        return
      }

      // R - 旋转模式
      if (keyLower === 'r' && !e.ctrlKey && !e.metaKey) {
        if (isInput) return
        this.editMode = 'rotate'
        return
      }

      // S - 缩放模式
      if (keyLower === 's' && !e.ctrlKey && !e.metaKey) {
        if (isInput) return
        this.editMode = 'scale'
        return
      }

      // G - 聚焦选中对象
      if (keyLower === 'g' && !e.ctrlKey && !e.metaKey) {
        if (isInput) return
        if (this.selectedObjectId) {
          this.handleFocusSelected()
        }
        return
      }

      // F - 全景视图
      if (keyLower === 'f' && !e.ctrlKey && !e.metaKey) {
        if (isInput) return
        this.handleFrameAll()
        return
      }

      // Escape - 取消选择和取消绘制
      if (key === 'Escape') {
        this.selectedObjectId = null
        // 通知 SceneCanvas 取消正在绘制的管道
        if (this.$refs.sceneCanvas && this.$refs.sceneCanvas.cancelDrawingPipe) {
          this.$refs.sceneCanvas.cancelDrawingPipe()
        }
        return
      }

      if (this.selectedObjectId && !e.ctrlKey && !e.metaKey) {
        if (isInput) return
        const step = e.shiftKey ? 0.1 : 0.01
        const offsetValue = function(value, delta) {
          const base = Number(value)
          return Number(((isNaN(base) ? 0 : base) + delta).toFixed(4))
        }
        let changes = null

        switch (key) {
          case 'ArrowLeft':
            changes = { x: offsetValue(this.selectedObject.x, -step) }
            break
          case 'ArrowRight':
            changes = { x: offsetValue(this.selectedObject.x, step) }
            break
          case 'ArrowUp':
            changes = { z: offsetValue(this.selectedObject.z, -step) }
            break
          case 'ArrowDown':
            changes = { z: offsetValue(this.selectedObject.z, step) }
            break
        }

        if (changes) {
          e.preventDefault()
          // 先记快照再修改，确保 Ctrl+Z 一次回到移动前
          this.pushHistory()
          this.updateObject(this.selectedObjectId, changes)
        }
      }
    },

    /** 统一规范化场景对象中的颜色字段，确保 <input type="color"> 永不报错 */
    normalizeSceneObjects(objects) {
      var self = this
      if (!Array.isArray(objects)) return
      objects.forEach(function(obj) {
        if (!obj) return
        if (obj.type === 'text3d' || obj.type === 'dataText' || obj.type === 'textPlain3d' || obj.type === 'uiLabel') {
          if (typeof obj.color !== 'string' || obj.color.charAt(0) !== '#') obj.color = '#ffffff'
        } else if (obj.type === 'gltf' && !obj.materialOverridden) {
        } else {
          obj.color = self.normColorValue(obj.color)
        }
        if (obj.type === 'text3d' || obj.type === 'dataText' || obj.type === 'textPlain3d' || obj.type === 'uiLabel') {
          obj.textBgColor = self.normColorValue(obj.textBgColor, '#00000000')
          if (obj.fontSize === undefined) obj.fontSize = 16
          if (obj.labelOpacity === undefined) obj.labelOpacity = 1
          if (obj.textBgOpacity === undefined) obj.textBgOpacity = obj.textBgColor === '#00000000' ? 0 : 0.2
          if (obj.labelFaceCamera === undefined) obj.labelFaceCamera = true
          if (obj.labelFixedSize === undefined) obj.labelFixedSize = true
          if (obj.labelFontFamily === undefined) obj.labelFontFamily = '系统默认'
          if (obj.labelRenderMode === undefined) obj.labelRenderMode = 'component'
          if (obj.labelShowBorder === undefined) obj.labelShowBorder = false
          if (obj.labelBorderWidth === undefined) obj.labelBorderWidth = 1
          if (obj.labelBorderColor === undefined) obj.labelBorderColor = '#ffffff'
        }
        if (obj.type === 'image3d' || obj.type === 'uiImage' || obj.type === 'video3d' || obj.type === 'webEmbed') {
          if (obj.mediaUrl === undefined) obj.mediaUrl = obj.imageUrl || obj.videoUrl || ''
          if (obj.imageUrl === undefined) obj.imageUrl = obj.type === 'image3d' || obj.type === 'uiImage' ? obj.mediaUrl : ''
          if (obj.videoUrl === undefined) obj.videoUrl = obj.type === 'video3d' ? obj.mediaUrl : ''
          if (obj.webUrl === undefined) obj.webUrl = ''
          if (obj.mediaWidth === undefined) obj.mediaWidth = obj.type === 'video3d' ? 1.6 : 1.4
          if (obj.mediaAspect === undefined) obj.mediaAspect = obj.type === 'video3d' ? 16 / 9 : 4 / 3
          if (obj.mediaAutoplay === undefined) obj.mediaAutoplay = true
          if (obj.mediaLoop === undefined) obj.mediaLoop = true
          if (obj.mediaMuted === undefined) obj.mediaMuted = true
          if (obj.uiX === undefined) obj.uiX = 40
          if (obj.uiY === undefined) obj.uiY = 40
          if (obj.uiWidth === undefined) obj.uiWidth = obj.type === 'webEmbed' ? 420 : 180
          if (obj.uiHeight === undefined) obj.uiHeight = obj.type === 'webEmbed' ? 260 : 100
          if (obj.uiRotation === undefined) obj.uiRotation = 0
        }
        // 确保 emissive 存在（GLTF 模型未修改材质时不覆盖）
        if (obj.type === 'gltf' && !obj.materialOverridden) {
        } else if (!obj.emissive) {
          obj.emissive = '#000000'
        }
        if (obj.showShadow === undefined) obj.showShadow = false
        if (obj.autoRotate === undefined) obj.autoRotate = false
        if (obj.rotateSpeed === undefined) obj.rotateSpeed = 1
        if (obj.rotateAxis === undefined) obj.rotateAxis = 'y'
        if (obj.floatAnim === undefined) obj.floatAnim = false
        if (obj.floatRange === undefined) obj.floatRange = 0.15
        if (obj.floatSpeed === undefined) obj.floatSpeed = 2
        if (obj.blink === undefined) obj.blink = false
        if (obj.blinkSpeed === undefined) obj.blinkSpeed = 6
        if (obj.blinkMin === undefined) obj.blinkMin = 0.2
        if (obj.gltfAnimationPlaying === undefined) obj.gltfAnimationPlaying = false
        if (obj.gltfAnimationName === undefined) obj.gltfAnimationName = ''
        if (!Array.isArray(obj.gltfAnimationNames)) obj.gltfAnimationNames = obj.gltfAnimationName ? [obj.gltfAnimationName] : []
        if (obj.gltfAnimationSpeed === undefined) obj.gltfAnimationSpeed = 1
        if (obj.gltfAnimationLoop === undefined) obj.gltfAnimationLoop = true
        if (!Array.isArray(obj.gltfAnimations)) obj.gltfAnimations = []
        obj.gltfAnimationGroups = ensureGLTFAnimationGroups(obj)
        syncLegacyGLTFAnimationFields(obj, obj.gltfAnimationGroups)
        if (obj.gltfAnimationConditionEnabled === undefined) obj.gltfAnimationConditionEnabled = false
        if (!obj.gltfAnimationCondition) {
          obj.gltfAnimationCondition = {
            isBandDevice: false,
            deviceSN: '',
            DeviceName: '',
            dataID: '',
            dataName: '',
            operator: '',
            OperatorValue: '',
            OperatorMaxValue: ''
          }
        }
        if (obj.wsKey === undefined) obj.wsKey = ''
        if (obj.bindProp === undefined) obj.bindProp = ''
        if (obj.bindTransform === undefined) obj.bindTransform = 'direct'
        if (obj.bindScale === undefined) obj.bindScale = 1
        if (obj.bindOffset === undefined) obj.bindOffset = 0
        ensurePositionBindings(obj)
        if (obj.realTimeValue === undefined) obj.realTimeValue = ''
        // 贴图数据
        if (obj.textureData === undefined) obj.textureData = ''
        if (obj.type === 'flowPipe') {
          if (obj.flowDirection === undefined) obj.flowDirection = 'forward'
          if (obj.flowAnim === undefined) obj.flowAnim = true
          if (obj.flowDashLength === undefined) obj.flowDashLength = 3
        }
        if (obj.dataFormat === undefined) obj.dataFormat = '{value}'
        if (obj.type === '2dComponent') {
          obj.is2DComponent = true
          if (obj.label2D === undefined) obj.label2D = obj.name || obj.typeName || '2D Component'
          if (obj.source2D === undefined) obj.source2D = null
          if (obj.source2DMeta === undefined) obj.source2DMeta = null
          self.sync2DWrapperFields(obj)
        }
        if (!obj.interactions) {
          obj.interactions = createDefaultInteractionConfig()
        }
        if ((!obj.interactions || Object.keys(obj.interactions).length === 0) && Array.isArray(obj.action)) {
          obj.interactions = createDefaultInteractionConfig()
          obj.action.forEach(actionItem => {
            if (!actionItem || !actionItem.type) return
            const interaction = obj.interactions[actionItem.type]
            if (!interaction) return
            interaction.enabled = true
            interaction.actionAuth = actionItem.actionAuth || []
            interaction.actionVoice = actionItem.actionVoice || ''
            interaction.actionConfirm = !!actionItem.actionConfirm
            if (actionItem.action === 'link' && actionItem.link) {
              if (actionItem.link.linkType === 'Inside') {
                interaction.type = 'route'
                interaction.routePath = actionItem.link.routePath || ''
                interaction.displayMode = actionItem.link.isPopUp ? 'popup' : 'window'
                interaction.autoClose = !!actionItem.link.autoClose
              } else {
                interaction.type = 'link'
                interaction.payload = actionItem.link.External || ''
                interaction.target = actionItem.link.OpenExternalType === 'self' ? '_self' : '_blank'
              }
              interaction.width = actionItem.link.width || interaction.width
              interaction.height = actionItem.link.height || interaction.height
              interaction.title = actionItem.link.title || ''
            } else if (actionItem.action === 'DeviceView' && actionItem.DeviceView) {
              interaction.type = 'deviceView'
              interaction.deviceKey = actionItem.DeviceView.key || ''
              interaction.showUUID = actionItem.DeviceView.showUUID || ''
              interaction.showPageUUID = actionItem.DeviceView.showPageUUID || ''
              interaction.deviceType = actionItem.DeviceView.type || ''
              interaction.isPopUp = !!actionItem.DeviceView.isPopUp
              interaction.isContainer = !!actionItem.DeviceView.isContainer
              interaction.routePath = actionItem.DeviceView.routePath || ''
            } else if (actionItem.action === 'visible') {
              interaction.type = 'visible'
              interaction.showTargetsList = actionItem.showItems || []
              interaction.hideTargetsList = actionItem.hideItems || []
              interaction.showTargets = interaction.showTargetsList.join(',')
              interaction.hideTargets = interaction.hideTargetsList.join(',')
            } else if (actionItem.action === 'popupText') {
              interaction.type = 'popupText'
              interaction.payload = actionItem.popupText || ''
            }
          })
        }
        if (!obj.animate) {
          obj.animate = {
            selected: [],
            condition: {
              deviceSN: "",
              selectVideoType: 0,
              isBandDevice: false,
              bandType: 1,
              dataID: "",
              dataName: "",
              operator: "",
              OperatorValue: "",
              OperatorMaxValue: "",
            },
            isExpression: false,
            animateList: [
              { id: "Forbidden", name: "禁用" },
              { id: "autoRotate", name: "自旋" },
              { id: "floatAnim", name: "浮动" },
              { id: "blink", name: "闪烁" },
              { id: "visible", name: "显隐" },
            ],
            animateElement: [
              {
                id: "autoRotate",
                elementList: [
                  { name: "速度", type: 7, value: 1, min: -5, max: 5, key: "rotateSpeed" },
                  { name: "轴向", type: 6, value: 'y', options: ['x','y','z'], key: "rotateAxis" },
                ]
              },
              {
                id: "floatAnim",
                elementList: [
                  { name: "幅度", type: 7, value: 0.15, min: 0, max: 1, key: "floatRange" },
                  { name: "速度", type: 7, value: 2, min: 0.1, max: 10, key: "floatSpeed" },
                ]
              },
              {
                id: "blink",
                elementList: [
                  { name: "速度", type: 7, value: 6, min: 1, max: 20, key: "blinkSpeed" },
                  { name: "最低透明度", type: 7, value: 0.2, min: 0, max: 1, key: "blinkMin" },
                ]
              },
              {
                id: "visible",
                elementList: [
                  { name: "显隐条件", type: 8, key: "condition" },
                ]
              },
            ]
          }
        }
        // 补全 gltf 相关字段
        if (obj.type === 'gltf') {
          if (obj.isExternalModel === undefined) obj.isExternalModel = true
          if (obj.isBackground === undefined) obj.isBackground = false
          if (obj.modelPath === undefined) obj.modelPath = ''
        }
        if (!Array.isArray(obj.action)) {
          obj.action = self.getObjectActions(obj)
        }
      })
    },
    restoreGLTFModels(objects) {
      var self = this
      var hasBackgroundModel = false
      objects.forEach(function(obj) {
        if (obj.type !== 'gltf') return
        if (obj.isBackground) hasBackgroundModel = true
        if (obj.modelPath) {
          console.log('从路径恢复 GLTF 模型:', obj.modelPath, 'ID:', obj.id)
          if (self.$refs.sceneCanvas && self.$refs.sceneCanvas.loadGLTFModelUrl) {
            self.$refs.sceneCanvas.loadGLTFModelUrl(obj.modelPath, obj.id, function(loadedId, mesh) {
              console.log('GLTF 模型加载成功:', loadedId)
              if (obj.isBackground && mesh) {
                mesh.traverse(function(child) {
                  if (child.isMesh && child.material) {
                    child.material.transparent = true
                    child.material.opacity = obj.opacity !== undefined ? obj.opacity : 0.5
                    child.material.needsUpdate = true
                  }
                })
              }
            })
          } else {
            console.error('场景画布未初始化，无法加载 GLTF 模型')
          }
          return
        } else {
          console.warn('GLTF 模型缺少 modelPath:', obj.id)
        }
      })
      if (hasBackgroundModel) {
        self.showGrid = false
        if (self.$refs.sceneCanvas && self.$refs.sceneCanvas.toggleGrid) {
          self.$refs.sceneCanvas.toggleGrid(false)
        }
      }
    },

    normColorValue(val, fallback) {
      fallback = fallback || '#000000'
      if (!val || typeof val !== 'string') return fallback
      var s = val.trim()
      if (COMP_COLORS[s]) return COMP_COLORS[s]
      // 已带 # 前缀
      if (s.charAt(0) === '#') {
        if (/^#[0-9a-fA-F]{3}$/.test(s)) {
          return '#' + s[1] + s[1] + s[2] + s[2] + s[3] + s[3]
        }
        if (/^#[0-9a-fA-F]{6}$/.test(s)) return s.toLowerCase()
        if (/^#[0-9a-fA-F]{8}$/.test(s)) return s.slice(0, 7).toLowerCase()
      }
      return fallback
    },

    // ========================
    // ========================

    /** SceneCanvas Three.js 初始化完成回调：若场景数据已加载，立即同步到画布 */
    onSceneInitialized() {
      this.$nextTick(() => {
        if (this.sceneObjects && this.sceneObjects.length > 0) {
          this._syncSceneToCanvas({ objects: this.sceneObjects, sceneSettings: this._savedSceneSettings })
        }
      })
    },

    _syncSceneToCanvas(dataJson) {
      var self = this
      var maxWait = 3000
      var interval = 50
      var waited = 0

      function trySync() {
        if (!self.$refs.sceneCanvas || !self.$refs.sceneCanvas.__3d) {
          if (waited >= maxWait) {
            console.warn('_syncSceneToCanvas: SceneCanvas 未就绪，超时放弃')
            return
          }
          waited += interval
          setTimeout(trySync, interval)
          return
        }

        var objects = dataJson.objects
        var sceneSettings = dataJson.sceneSettings

        self.$refs.sceneCanvas.syncMeshes(objects)
        self.restoreGLTFModels(objects)
        self.updateGrid()
        if (self.$refs.sceneCanvas.applySceneExtras) {
          self.$refs.sceneCanvas.applySceneExtras(self.sceneExtras)
        }
        if (self.$refs.sceneCanvas.toggleGrid) {
          self.$refs.sceneCanvas.toggleGrid(self.showGrid)
        }
        if (sceneSettings && self.$refs.sceneCanvas.restoreCamera) {
          self.$refs.sceneCanvas.restoreCamera(sceneSettings)
        }
      }

      trySync()
    },

    /** 从服务端加载场景 */
    normalizeScenePage(layer, index) {
      const pageId = layer && (layer.PageId || layer.pageUuid || layer.PageUuid || layer.id)
      const title = layer && (layer.PageName || layer.name || layer.title)
      return {
        id: pageId ? String(pageId) : String(index),
        pageUuid: pageId ? String(pageId) : String(index),
        title: title || (this.$t('ISM3DEditor.pageDefaultPrefix') + (index + 1)),
        pageType: layer && layer.PageType !== undefined ? layer.PageType : 1,
        pageModelUuid: (layer && (layer.modelId || layer.pageModelUuid || layer.modelUuid)) || (this.$route && this.$route.params && this.$route.params.uid),
        isHome: !!(layer && (layer.IsHome === 1 || layer.IsHome === '1' || layer.IsHome === true || layer.isHome === true)),
        raw: layer
      }
    },

    normalizeScenePages(layers) {
      const list = Array.isArray(layers) ? layers : []
      const pages = list.map((layer, index) => this.normalizeScenePage(layer, index)).filter(page => page.id)
      if (pages.length && !pages.some(page => page.isHome)) {
        pages[0].isHome = true
      }
      return pages
    },

    getSelectedSceneLayer(layers, targetPageId) {
      const list = Array.isArray(layers) ? layers : []
      if (!list.length) return null
      const query = (this.$route && this.$route.query) || {}
      const queryPageId = query.pageId || query.pageid || query.page || query.pageUuid
      const wanted = targetPageId || queryPageId || this.currentPageId || this.SelectPagerID
      if (wanted) {
        const matched = list.find(layer => String(layer.PageId || layer.pageUuid || layer.PageUuid || layer.id) === String(wanted))
        if (matched) return matched
      }
      return list.find(layer => layer.IsHome === 1 || layer.IsHome === '1' || layer.IsHome === true || layer.isHome === true) || list[0]
    },

    applyScenePageLayer(pageLayer) {
      let dataJson = parseScenePayloadOrEmpty(pageLayer && pageLayer.components)
      if (!dataJson.objects || !Array.isArray(dataJson.objects)) {
        message.error(this.$t('ISM3DEditor.invalidSceneFormat'))
        return false
      }
      const valid = dataJson.objects.every(obj => obj.id && obj.type)
      if (!valid) {
        message.error(this.$t('ISM3DEditor.invalidSceneObjectFields'))
        return false
      }
      if (this.sceneObjects.length > 0) {
        this.pushHistory()
      }
      this.normalizeSceneObjects(dataJson.objects)
      this.sceneObjects = dataJson.objects
      this._savedSceneSettings = dataJson.sceneSettings || null
      this.selectedObjectId = null
      if (dataJson.sceneSettings) {
        this.applySceneSettings(dataJson.sceneSettings)
      } else {
        this.applySceneSettings(defaultGridSettings())
      }
      this.sceneExtras = Object.assign(createDefaultSceneExtras(), dataJson.sceneExtras || {})
      this.$nextTick(() => { this._syncSceneToCanvas(dataJson) })
      return true
    },

    async loadSceneFromServer(targetPageId) {
      const uid = this.$route && this.$route.params && this.$route.params.uid
      if (!uid) return

      this.pageLoading = true
      try {
        let params = {
          muid:uid,
        }
        const res = await getDisplayModelLayerData(params)
        if(res.data.code==0)
        {
          const layers = Array.isArray(res.data.layer) ? res.data.layer : []
          this.scenePages = this.normalizeScenePages(layers)
          let pageLayer = this.getSelectedSceneLayer(layers, targetPageId)
          if (!pageLayer) return null
          const selectedPageId = String(pageLayer.PageId || pageLayer.pageUuid || pageLayer.PageUuid || pageLayer.id)
          try {
            const applied = this.applyScenePageLayer(pageLayer)
            if (!applied) return null
            this.SelectPagerID = selectedPageId
            this.currentPageId = selectedPageId
            this.isCharge = false
            return this.currentPageId
          } catch (parseErr) {
            // 场景数据解析失败，使用空场景
            console.log(parseErr)
          }
        }
      } catch (err) {
        // 接口不存在或加载失败，使用空场景
        console.log(err)
      } finally {
        this.pageLoading = false
      }
      return null
    },

    /** 保存场景到服务端 */
    async saveCurrentSceneBeforePageChange() {
      if (!this.currentPageId && !this.SelectPagerID) return true
      return this.saveSceneToServer({ silentSuccess: true })
    },

    async handleSwitchPage(pageId) {
      if (!pageId || String(pageId) === String(this.currentPageId)) return
      if (this.pageSwitching || this.pageLoading) return
      const oldPageId = this.currentPageId
      const oldSelectPagerID = this.SelectPagerID
      this.pageSwitching = true
      try {
        const saved = await this.saveCurrentSceneBeforePageChange()
        if (!saved) {
          this.currentPageId = oldPageId
          this.SelectPagerID = oldSelectPagerID
          return
        }
        const loadedPageId = await this.loadSceneFromServer(pageId)
        if (!loadedPageId) {
          this.currentPageId = oldPageId
          this.SelectPagerID = oldSelectPagerID
        }
      } finally {
        this.pageSwitching = false
      }
    },

    openPageDialog(mode) {
      this.pageDialogMode = mode
      this.pageForm.name = mode === 'edit' && this.currentPage ? this.currentPage.title : ''
      this.pageDialogVisible = true
    },

    isPageApiSuccess(res, codes) {
      const successCodes = codes || [0, 200, 4002]
      return !!(res && res.data && successCodes.indexOf(res.data.code) !== -1)
    },

    getPageIdFromPageResponse(res) {
      const data = res && res.data
      if (!data || typeof data !== 'object') return ''
      const candidates = [
        data.PageId,
        data.pageId,
        data.pageid,
        data.pageUuid,
        data.PageUuid,
        data.uuid,
        data.id,
        data.data && data.data.PageId,
        data.data && data.data.pageId,
        data.data && data.data.pageUuid,
        data.data && data.data.PageUuid,
        data.data && data.data.uuid,
        data.data && data.data.id,
        data.page && data.page.PageId,
        data.page && data.page.pageId,
        data.page && data.page.pageUuid,
        data.page && data.page.PageUuid,
        data.page && data.page.uuid,
        data.page && data.page.id,
      ]
      const found = candidates.find(value => value !== undefined && value !== null && value !== '')
      return found ? String(found) : ''
    },

    findAddedPageId(previousPageIds, name) {
      const oldIds = new Set((previousPageIds || []).map(id => String(id)))
      const added = this.scenePages.filter(page => page && page.id && !oldIds.has(String(page.id)))
      if (!added.length) return ''
      const sameName = added.find(page => page.title === name)
      return String((sameName || added[added.length - 1]).id)
    },

    async initializeBlankScenePage(pageId) {
      if (!pageId) return false
      this.currentPageId = String(pageId)
      this.SelectPagerID = String(pageId)
      this.sceneObjects = []
      this.selectedObjectId = null
      this.applySceneSettings(defaultGridSettings())
      this.sceneExtras = createDefaultSceneExtras()
      this._savedSceneSettings = defaultGridSettings()
      this.$nextTick(() => {
        this._syncSceneToCanvas(createEmptyScenePayload())
      })
      return this.saveSceneToServer({ silentSuccess: true })
    },

    async submitPageDialog() {
      const uid = this.$route && this.$route.params && this.$route.params.uid
      const name = (this.pageForm.name || '').trim()
      if (!uid || !name) {
        message.warning(this.$t('ISM3DEditor.inputPageName'))
        return
      }
      this.pageDialogLoading = true
      try {
        let res
        const previousPageIds = this.scenePages.map(page => page.id)
        if (this.pageDialogMode === 'add') {
          const saved = await this.saveCurrentSceneBeforePageChange()
          if (!saved) return
          res = await DisplayModelPageAdd({
            modelUuid: uid,
            name: name,
            isLogin: 0,
            size: '1',
            pageType: 1
          })
        } else {
          res = await DisplayModelPageEdit({
            modelUuid: uid,
            pageUuid: this.currentPageId,
            name: name
          })
        }
        if (this.isPageApiSuccess(res)) {
          this.pageDialogVisible = false
          message.success(this.pageDialogMode === 'add' ? this.$t('ISM3DEditor.addPageSuccess') : this.$t('ISM3DEditor.renamePageSuccess'))
          if (this.pageDialogMode === 'add') {
            const responsePageId = this.getPageIdFromPageResponse(res)
            await this.loadSceneFromServer(this.currentPageId)
            const responsePageExists = responsePageId && this.scenePages.some(page => String(page.id) === String(responsePageId))
            const newPageId = responsePageExists ? responsePageId : this.findAddedPageId(previousPageIds, name)
            if (newPageId) {
              await this.initializeBlankScenePage(newPageId)
              await this.loadSceneFromServer(newPageId)
            }
          } else {
            await this.loadSceneFromServer(this.currentPageId)
          }
        } else {
          message.error((res && res.data && res.data.message) || this.$t('ISM3DEditor.pageOperationFailed'))
        }
      } catch (err) {
        console.log(err)
        message.error(this.$t('ISM3DEditor.pageOperationFailed'))
      } finally {
        this.pageDialogLoading = false
      }
    },

    async copyCurrentPage() {
      const uid = this.$route && this.$route.params && this.$route.params.uid
      if (!uid || !this.currentPageId) return
      const saved = await this.saveSceneToServer({ silentSuccess: true })
      if (!saved) return
      try {
        const res = await DisplayModelPageCopy({
          muid: uid,
          pageid: this.currentPageId,
          pageType: this.currentPage ? this.currentPage.pageType : 1
        })
        if (this.isPageApiSuccess(res, [0, 200])) {
          message.success(this.$t('ISM3DEditor.copyPageSuccess'))
          await this.loadSceneFromServer(this.currentPageId)
        } else {
          message.error((res && res.data && res.data.message) || this.$t('ISM3DEditor.copyPageFailed'))
        }
      } catch (err) {
        console.log(err)
        message.error(this.$t('ISM3DEditor.copyPageFailed'))
      }
    },

    async setCurrentPageHome() {
      const uid = this.$route && this.$route.params && this.$route.params.uid
      if (!uid || !this.currentPageId) return
      try {
        const res = await DisplayModelPageSetHome({
          muid: uid,
          pageid: this.currentPageId,
          pageType: this.currentPage ? this.currentPage.pageType : 1
        })
        if (this.isPageApiSuccess(res)) {
          message.success(this.$t('ISM3DEditor.setHomeSuccess'))
          await this.loadSceneFromServer(this.currentPageId)
        } else {
          message.error((res && res.data && res.data.message) || this.$t('ISM3DEditor.setHomeFailed'))
        }
      } catch (err) {
        console.log(err)
        message.error(this.$t('ISM3DEditor.setHomeFailed'))
      }
    },

    async deleteCurrentPage() {
      const uid = this.$route && this.$route.params && this.$route.params.uid
      if (!uid || !this.currentPageId || this.scenePages.length <= 1) return
      if (!window.confirm(this.$t('ISM3DEditor.confirmDeleteCurrentPage'))) return
      try {
        const res = await DisplayModelPageDel({
          modelUuid: uid,
          pageId: this.currentPageId
        })
        if (this.isPageApiSuccess(res)) {
          message.success(this.$t('ISM3DEditor.deletePageSuccess'))
          const nextPage = this.scenePages.find(page => page.id !== this.currentPageId)
          await this.loadSceneFromServer(nextPage && nextPage.id)
        } else {
          message.error((res && res.data && res.data.message) || this.$t('ISM3DEditor.deletePageFailed'))
        }
      } catch (err) {
        console.log(err)
        message.error(this.$t('ISM3DEditor.deletePageFailed'))
      }
    },

    async saveSceneToServer(options = {}) {
      const silentSuccess = !!options.silentSuccess
      const uid = this.$route && this.$route.params && this.$route.params.uid
      if (!uid) {
        message.warning(this.$t('ISM3DEditor.missingAppIdSave'))
        return false
      }
      this.flush2DStoreSelectionToScene()
      // 保存前记录 GLTF 对象属性
      console.log('=== 保存场景前 ===')
      this.sceneObjects.filter(o => o.type === 'gltf').forEach(gltfObj => {
        console.log('保存前 GLTF 对象:', gltfObj.id, 'color:', gltfObj.color, 'opacity:', gltfObj.opacity,
          'metalness:', gltfObj.metalness, 'roughness:', gltfObj.roughness, 'wireframe:', gltfObj.wireframe)
      })
      let _t=this
      try {
        this.SelectPagerID = this.currentPageId || this.SelectPagerID
        let params = {
          uuid:uid,
          pageid:this.SelectPagerID ,
          LayerData:{},
        }
        const payload = {
          version: '1.0',
          generator: 'ISM3DEditor',
          exportTime: new Date().toISOString(),
          objectCount: this.sceneObjects.length,
          sceneSettings: Object.assign({}, this.buildSceneSettingsPayload(), {
            cameraPosition: this.$refs.sceneCanvas && this.$refs.sceneCanvas.__3d ? {
              x: this.$refs.sceneCanvas.__3d.camera.position.x,
              y: this.$refs.sceneCanvas.__3d.camera.position.y,
              z: this.$refs.sceneCanvas.__3d.camera.position.z,
            } : { x: 6, y: 5, z: 8 },
            cameraTarget: this.$refs.sceneCanvas && this.$refs.sceneCanvas.__3d ? {
              x: this.$refs.sceneCanvas.__3d.orbit.target.x,
              y: this.$refs.sceneCanvas.__3d.orbit.target.y,
              z: this.$refs.sceneCanvas.__3d.orbit.target.z,
            } : { x: 0, y: 0, z: 0 },
            cameraFov: this.$refs.sceneCanvas && this.$refs.sceneCanvas.__3d ? this.$refs.sceneCanvas.__3d.camera.fov : 45,
          }),
          sceneExtras: JSON.parse(JSON.stringify(this.sceneExtras || createDefaultSceneExtras())),
          objects: this.sceneObjects.map(obj => ({
            id: obj.id,
            type: obj.type,
            name: obj.name || '',
            typeName: obj.typeName || obj.type,
            icon: obj.icon || '',
            groupName: obj.groupName || '',
            intensity: obj.intensity !== undefined ? obj.intensity : undefined,
            distance: obj.distance !== undefined ? obj.distance : undefined,
            angle: obj.angle !== undefined ? obj.angle : undefined,
            penumbra: obj.penumbra !== undefined ? obj.penumbra : undefined,
            category: obj.category || '',
            x: obj.x !== undefined ? obj.x : 0,
            y: obj.y !== undefined ? obj.y : 0,
            z: obj.z !== undefined ? obj.z : 0,
            rx: obj.rx !== undefined ? obj.rx : 0,
            ry: obj.ry !== undefined ? obj.ry : 0,
            rz: obj.rz !== undefined ? obj.rz : 0,
            sx: obj.sx !== undefined ? obj.sx : 1,
            sy: obj.sy !== undefined ? obj.sy : 1,
            sz: obj.sz !== undefined ? obj.sz : 1,
            color: (obj.type === 'text3d' || obj.type === 'textPlain3d' || obj.type === 'dataText' || obj.type === 'uiLabel') ? (obj.color || '#ffffff') : (obj.type === 'gltf' ? (obj.materialOverridden ? obj.color : undefined) : (obj.color || '#4a90d9')),
            opacity: obj.type === 'gltf' ? (obj.materialOverridden ? obj.opacity : undefined) : (obj.opacity !== undefined ? obj.opacity : 1),
            metalness: obj.type === 'gltf' ? (obj.materialOverridden ? obj.metalness : undefined) : (obj.metalness !== undefined ? obj.metalness : 0.3),
            roughness: obj.type === 'gltf' ? (obj.materialOverridden ? obj.roughness : undefined) : (obj.roughness !== undefined ? obj.roughness : 0.7),
            wireframe: obj.type === 'gltf' ? (obj.materialOverridden ? obj.wireframe : undefined) : (obj.wireframe || false),
            showShadow: !!obj.showShadow,
            emissive: obj.type === 'gltf' ? (obj.materialOverridden ? obj.emissive : undefined) : (obj.emissive || '#000000'),
            textureData: obj.type === 'gltf' ? (obj.materialOverridden ? obj.textureData : undefined) : (obj.textureData || ''),
            mediaUrl: obj.mediaUrl || '',
            imageUrl: obj.imageUrl || '',
            videoUrl: obj.videoUrl || '',
            webUrl: obj.webUrl || '',
            mediaWidth: obj.mediaWidth !== undefined ? obj.mediaWidth : undefined,
            mediaAspect: obj.mediaAspect !== undefined ? obj.mediaAspect : undefined,
            mediaAutoplay: obj.mediaAutoplay !== undefined ? obj.mediaAutoplay : true,
            mediaLoop: obj.mediaLoop !== undefined ? obj.mediaLoop : true,
            mediaMuted: obj.mediaMuted !== undefined ? obj.mediaMuted : true,
            uiX: obj.uiX !== undefined ? obj.uiX : undefined,
            uiY: obj.uiY !== undefined ? obj.uiY : undefined,
            uiWidth: obj.uiWidth !== undefined ? obj.uiWidth : undefined,
            uiHeight: obj.uiHeight !== undefined ? obj.uiHeight : undefined,
            uiRotation: obj.uiRotation !== undefined ? obj.uiRotation : 0,
            visible: obj.visible !== false,
            locked: obj.locked || false,
            textContent: obj.textContent || '',
            fontSize: obj.fontSize || 16,
            textBgColor: obj.textBgColor || ((obj.type === 'text3d' || obj.type === 'textPlain3d' || obj.type === 'dataText' || obj.type === 'uiLabel') ? '#00000000' : '#000000'),
            labelOpacity: obj.labelOpacity !== undefined ? obj.labelOpacity : 1,
            textBgOpacity: obj.textBgOpacity !== undefined ? obj.textBgOpacity : ((obj.type === 'text3d' || obj.type === 'textPlain3d' || obj.type === 'dataText' || obj.type === 'uiLabel') ? 0 : 0.2),
            labelFaceCamera: obj.labelFaceCamera !== undefined ? obj.labelFaceCamera : true,
            labelFixedSize: obj.labelFixedSize !== undefined ? obj.labelFixedSize : true,
            labelFontFamily: obj.labelFontFamily || '系统默认',
            labelRenderMode: obj.labelRenderMode || 'component',
            labelShowBorder: obj.labelShowBorder !== undefined ? !!obj.labelShowBorder : false,
            labelBorderWidth: obj.labelBorderWidth !== undefined ? obj.labelBorderWidth : 1,
            labelBorderColor: obj.labelBorderColor || '#ffffff',
            autoRotate: obj.autoRotate || false,
            rotateSpeed: obj.rotateSpeed !== undefined ? obj.rotateSpeed : 1,
            rotateAxis: obj.rotateAxis || 'y',
            floatAnim: obj.floatAnim || false,
            floatRange: obj.floatRange !== undefined ? obj.floatRange : 0.15,
            floatSpeed: obj.floatSpeed !== undefined ? obj.floatSpeed : 2,
            blink: obj.blink || false,
            blinkSpeed: obj.blinkSpeed !== undefined ? obj.blinkSpeed : 6,
            blinkMin: obj.blinkMin !== undefined ? obj.blinkMin : 0.2,
            gltfAnimationPlaying: !!obj.gltfAnimationPlaying,
            gltfAnimationName: obj.gltfAnimationName || '',
            gltfAnimationNames: Array.isArray(obj.gltfAnimationNames) ? obj.gltfAnimationNames.slice() : (obj.gltfAnimationName ? [obj.gltfAnimationName] : []),
            gltfAnimationSpeed: obj.gltfAnimationSpeed !== undefined ? obj.gltfAnimationSpeed : 1,
            gltfAnimationLoop: obj.gltfAnimationLoop !== false,
            gltfAnimations: cloneGLTFAnimations(obj.gltfAnimations),
            gltfAnimationGroups: ensureGLTFAnimationGroups(obj),
            gltfAnimationConditionEnabled: !!obj.gltfAnimationConditionEnabled,
            gltfAnimationCondition: obj.gltfAnimationCondition ? JSON.parse(JSON.stringify(obj.gltfAnimationCondition)) : null,
            interactions: obj.interactions ? JSON.parse(JSON.stringify(obj.interactions)) : createDefaultInteractionConfig(),
            action: Array.isArray(obj.action) ? JSON.parse(JSON.stringify(obj.action)) : this.getObjectActions(obj),
            wsKey: obj.wsKey || '',
            bindProp: obj.bindProp || '',
            bindTransform: obj.bindTransform || 'direct',
            bindScale: obj.bindScale !== undefined ? obj.bindScale : 1,
            bindOffset: obj.bindOffset !== undefined ? obj.bindOffset : 0,
            realTimeValue: obj.realTimeValue !== undefined ? obj.realTimeValue : '',
            dataFormat: obj.dataFormat || '{value}',
            is2DComponent: !!obj.is2DComponent,
            label2D: obj.label2D || '',
            source2D: obj.source2D ? JSON.parse(JSON.stringify(obj.source2D)) : null,
            source2DMeta: obj.source2DMeta ? JSON.parse(JSON.stringify(obj.source2DMeta)) : null,
            animate: obj.animate ? JSON.parse(JSON.stringify(obj.animate)) : null,
            isExternalModel: obj.isExternalModel || false,
            isBackground: obj.isBackground || false,
            fitSize: obj.fitSize,
            modelPath: obj.modelPath || '',
            materialOverridden: !!obj.materialOverridden,
            points: obj.type === 'flowPipe' ? (obj.points || []) : undefined,
            radius: obj.type === 'flowPipe' ? (obj.radius !== undefined ? obj.radius : 0.1) : undefined,
            flowSpeed: obj.type === 'flowPipe' ? (obj.flowSpeed !== undefined ? obj.flowSpeed : 1.0) : undefined,
            flowDirection: obj.type === 'flowPipe' ? (obj.flowDirection || 'forward') : undefined,
            flowAnim: obj.type === 'flowPipe' ? (obj.flowAnim !== undefined ? obj.flowAnim : true) : undefined,
            highlightColor: obj.type === 'flowPipe' ? (obj.highlightColor || '#ff6a00') : undefined,
            flowDashLength: obj.type === 'flowPipe' ? (obj.flowDashLength !== undefined ? obj.flowDashLength : 3) : undefined,
          }))
        }
        console.log('保存场景数据，包含 GLTF 模型数量:', payload.objects.filter(function(o) { return o.type === 'gltf' && o.modelPath }).length)
        payload.objects.filter(function(o) { return o.type === 'gltf' }).forEach(function(gltfObj) {
          console.log('GLTF 模型保存信息:', gltfObj.id, '路径:', gltfObj.modelPath, '名称:', gltfObj.name,
            'color:', gltfObj.color, 'opacity:', gltfObj.opacity, 'metalness:', gltfObj.metalness, 'roughness:', gltfObj.roughness)
        })
        params.LayerData.components = JSON.stringify(payload)
        params.LayerData.layer={}
        const res = await this.saveLayerDataStruct(params)
        if(res.data.code == 200)
        {
          console.log('=== 保存场景后 ===')
          this.sceneObjects.filter(o => o.type === 'gltf').forEach(gltfObj => {
            console.log('保存后 GLTF 对象:', gltfObj.id, 'color:', gltfObj.color, 'opacity:', gltfObj.opacity,
              'metalness:', gltfObj.metalness, 'roughness:', gltfObj.roughness, 'wireframe:', gltfObj.wireframe)
          })
          // 保存成功，清除未保存修改标记
          this.isCharge = false
          if (!silentSuccess) {
            _t.$message.success(_t.$t('displayModel.SaveDataSuccess'))
          }
          return true
        }
        _t.$message.error(_t.$t('displayModel.SaveDataFailed'))
        return false
      } catch (err) {
        console.error('保存场景失败:', err)
        message.error('保存失败：' + (err.message || '未知错误'))
        return false
      }
    }
  }
}
</script>

<style lang="less">
.editor-page-bar {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  gap: 12px;
  height: 38px;
  padding: 0 12px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.08);
  background: #f7f9fb;
  flex: none;
}

.editor-page-select,
.editor-page-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.editor-page-actions {
  flex-wrap: wrap;
  justify-content: flex-start;
}

.editor-page-loading-inline {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #13a8a8;
  font-size: 12px;
  white-space: nowrap;
}

.editor-page-loading-mask {
  position: absolute;
  left: 50%;
  top: 50%;
  z-index: 40;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  color: #1f2933;
  font-size: 13px;
  background: rgba(255, 255, 255, 0.92);
  border: 1px solid rgba(19, 194, 194, 0.28);
  border-radius: 6px;
  box-shadow: 0 10px 30px rgba(15, 35, 45, 0.18);
  transform: translate(-50%, -50%);
  pointer-events: none;
}

.editor-page-label {
  color: #4b5563;
  font-size: 12px;
  white-space: nowrap;
}

.editor-page-home {
  margin-left: 6px;
  color: #13c2c2;
  font-size: 12px;
}

.scene-settings-modal {
  .ant-modal {
    top: 22px;
  }

  .ant-modal-content {
    max-height: calc(100vh - 44px);
    overflow: hidden;
    border-radius: 10px;
    background: rgba(247, 250, 250, 0.96);
    box-shadow: 0 20px 50px rgba(34, 55, 66, 0.22);
    backdrop-filter: blur(8px);
  }

  .ant-modal-header {
    padding: 20px 18px 14px;
    border-bottom: 0;
    background: transparent;
  }

  .ant-modal-title {
    font-size: 15px;
    font-weight: 700;
    color: #243241;
  }

  .ant-modal-close-x {
    width: 48px;
    height: 48px;
    line-height: 48px;
    color: #7a8a96;
  }

  .ant-modal-body {
    max-height: calc(100vh - 104px);
    padding: 0;
    background: transparent;
  }

  .scene-settings-panel {
    max-height: calc(100vh - 104px);
    overflow-y: auto;
    padding: 0 18px 22px;

    &::-webkit-scrollbar {
      width: 6px;
    }

    &::-webkit-scrollbar-thumb {
      border-radius: 6px;
      background: rgba(118, 148, 158, 0.45);
    }
  }

  .scene-settings-form {
    padding-bottom: 10px;

    .ant-form-item {
      display: flex;
      align-items: flex-start;
      margin-bottom: 14px !important;
      padding: 0;
    }

    .ant-form-item-label {
      width: 104px;
      padding: 5px 12px 0 0;
      text-align: left;
      line-height: 20px;
      flex: none;
      float: none;

      label {
        color: #26323c;
        font-size: 12px;
        font-weight: 600;
      }
    }

    .ant-form-item-control-wrapper {
      flex: 1;
      min-width: 0;
      width: auto;
      float: none;
    }

    .ant-form-item-control,
    .ant-form-item-children {
      line-height: 24px;
    }

    .ant-radio-group {
      display: grid;
      grid-template-columns: repeat(3, minmax(0, 1fr));
      gap: 8px;
      width: 100%;
    }

    .ant-radio-button-wrapper {
      height: 28px;
      min-width: 0;
      padding: 0 8px;
      border: 0;
      border-radius: 14px;
      background: rgba(255, 255, 255, 0.86);
      color: #5d6b77;
      line-height: 28px;
      text-align: center;
      box-shadow: 0 1px 8px rgba(38, 65, 75, 0.04);
      transition: all 0.18s ease;

      &::before {
        display: none;
      }
    }

    .ant-radio-button-wrapper-checked {
      background: #13c2c2;
      color: #fff;
      font-weight: 600;
      box-shadow: 0 8px 16px rgba(19, 194, 194, 0.22);
    }

    .ant-select,
    .ant-input,
    .ant-input-number {
      width: 100%;
    }

    .ant-select-selection,
    .ant-input,
    .ant-input-number {
      min-height: 34px;
      border-color: transparent;
      border-radius: 10px;
      background: rgba(255, 255, 255, 0.92);
      box-shadow: inset 0 0 0 1px rgba(214, 224, 229, 0.8);
    }

    .ant-select-selection__rendered {
      line-height: 34px;
    }

    .ant-switch {
      background-color: rgba(206, 216, 222, 0.9);
    }

    .ant-switch-checked {
      background-color: #13c2c2;
    }

    input[type='color'] {
      width: 48px !important;
      height: 28px !important;
      padding: 3px;
      border: 0;
      border-radius: 7px;
      background: rgba(255, 255, 255, 0.9);
      box-shadow: inset 0 0 0 1px rgba(205, 215, 220, 0.9);
      cursor: pointer;
    }
  }
}

.editor-hidden-file-input {
  display: none;
}

.scene-bg-model-drop {
  min-height: 92px;
  border: 1px dashed #13c2c2;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.72);
  color: #13c2c2;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  cursor: pointer;
  transition: background 0.18s ease, border-color 0.18s ease, box-shadow 0.18s ease;

  &:hover,
  &.dragover {
    background: #e6fffb;
    border-color: #13c2c2;
    box-shadow: 0 0 0 3px rgba(19, 194, 194, 0.12);
  }

  i {
    font-size: 26px;
    margin-bottom: 4px;
  }

  .scene-bg-model-title {
    font-size: 13px;
    font-weight: 500;
    line-height: 1.4;
  }

  .scene-bg-model-hint {
    color: #7b8a9a;
    font-size: 11px;
    line-height: 1.4;
  }
}

.ism-3d-editor {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100vh;
  overflow: hidden;
  background: #ffffff;
  outline: none;

    .editor-main {
      position: relative;
      display: flex;
      flex: 1;
      overflow: hidden;
      min-height: 0;

    .editor-left-panels {
      position: relative;
      width: 440px;
      min-width: 440px;
      flex-shrink: 0;
      transition: width 0.18s ease, min-width 0.18s ease;

      &.collapsed {
        width: 0;
        min-width: 0;
        overflow: visible;

        .editor-toolbox {
          display: none;
        }
      }

      .left-panel-toggle {
        position: absolute;
        right: -22px;
        top: 12px;
        z-index: 3;
        width: 22px;
        height: 48px;
        border: 1px solid #b5f5ec;
        background: #ffffff;
        color: #13c2c2;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 0;
        border-radius: 0 4px 4px 0;
        box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
      }
    }

    .editor-toolbox {
      width: 440px;
      height: 100%;
      flex-shrink: 0;
      background: #ffffff;
      border-right: 1px solid #b5f5ec;
    }

    .editor-canvas-container {
      flex: 1;
      display: flex;
      flex-direction: column;
      overflow: hidden;
    }



    .editor-right-panels {
      position: relative;
      width: 420px;
      flex-shrink: 0;
      display: flex;
      flex-direction: column;
      background: #ffffff;
      border-left: 1px solid #b5f5ec;
      transition: width 0.18s ease;

      &.collapsed {
        width: 0;
        min-width: 0;
        overflow: visible;
        border-left: 0;

        .editor-properties,
        .editor-hierarchy {
          display: none;
        }
      }

      .right-panel-toggle {
        position: absolute;
        left: -22px;
        top: 12px;
        z-index: 3;
        width: 22px;
        height: 48px;
        border: 1px solid #b5f5ec;
        border-right: 1px solid #b5f5ec;
        background: #ffffff;
        color: #13c2c2;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 0;
        border-radius: 4px 0 0 4px;
        box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
      }

        .editor-hierarchy {
          height: 220px;
          min-height: 150px;
          border-bottom: 1px solid #b5f5ec;
          overflow: auto;
        }

      .editor-properties {
        flex: 1;
        overflow: auto;
        padding-left: 0;
      }
    }
  }
}

/* 灯光设置下拉面板 */
.light-settings-dropdown {
  position: fixed;
  top: 48px;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 99;
}

.light-settings-card {
  position: absolute;
  top: 6px;
  width: 280px;
  background: #ffffff;
  border-radius: 10px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  overflow: hidden;
  font-size: 13px;
  color: #333;
}

.scene-settings-card {
  position: absolute;
  top: 6px;
  width: 340px;
  max-height: calc(100vh - 80px);
  overflow-y: auto;
  background: #ffffff;
  border-radius: 10px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  overflow: hidden;
  font-size: 13px;
  color: #333;
}

.inspect-settings-card {
  position: absolute;
  top: 6px;
  width: 420px;
  max-height: calc(100vh - 80px);
  overflow-y: auto;
  background: #ffffff;
  border-radius: 10px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  font-size: 13px;
  color: #333;
}

.inspect-panel-body {
  padding: 12px;
}

.inspect-metric-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 8px;
  margin-bottom: 10px;
}

.inspect-metric-grid div {
  padding: 10px 6px;
  border: 1px solid #edf0f5;
  border-radius: 6px;
  text-align: center;
  background: #fafbfc;
}

.inspect-metric-grid b {
  display: block;
  font-size: 18px;
  color: #13c2c2;
}

.inspect-metric-grid span {
  font-size: 12px;
  color: #666;
}

.inspect-actions {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px;
  margin-bottom: 10px;
}

.light-settings-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid #f0f0f0;
  font-weight: 600;
  font-size: 14px;
}

.light-settings-body {
  padding: 12px 16px;
}

.scene-settings-body {
  padding: 8px 12px 12px;
}

.scene-settings-tabs /deep/ .ant-tabs-bar {
  margin: 0 0 12px;
}

.scene-settings-tabs /deep/ .ant-tabs-nav .ant-tabs-tab {
  margin-right: 8px;
  padding: 6px 2px 8px;
  font-size: 12px;
}

.scene-settings-tabs /deep/ .ant-tabs-content {
  overflow: visible;
}

.scene-settings-tabs /deep/ .ant-tabs-tabpane {
  max-height: calc(100vh - 150px);
  overflow-y: auto;
  overflow-x: hidden;
  padding-right: 2px;
}

.scene-preset-desc {
  margin-top: 6px;
  color: #7a8a96;
  font-size: 11px;
  line-height: 1.45;
}

.light-section {
  margin-bottom: 14px;
}

.light-section:last-child {
  margin-bottom: 0;
}

.light-section-title {
  font-weight: 600;
  font-size: 13px;
  color: #333;
  margin-bottom: 8px;
  display: flex;
  align-items: center;
}

.light-row {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
  gap: 8px;
}

.light-row:last-child {
  margin-bottom: 0;
}

.light-label {
  width: 36px;
  font-size: 12px;
  color: #666;
  flex-shrink: 0;
}

.light-slider-wrap {
  flex: 1;
  display: flex;
  align-items: center;
}

.light-slider {
  width: 100%;
  height: 4px;
  border-radius: 2px;
  background: #e8e8e8;
  outline: none;
  -webkit-appearance: none;
  appearance: none;
}

.light-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background: #13c2c2;
  cursor: pointer;
  border: 2px solid #fff;
  box-shadow: 0 1px 4px rgba(0,0,0,0.2);
}

.light-slider::-moz-range-thumb {
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background: #13c2c2;
  cursor: pointer;
  border: 2px solid #fff;
  box-shadow: 0 1px 4px rgba(0,0,0,0.2);
}

.light-value {
  width: 32px;
  text-align: right;
  font-size: 12px;
  color: #666;
  flex-shrink: 0;
}

.light-color-wrap {
  flex: 1;
  display: flex;
  justify-content: flex-end;
}

.light-color-input {
  width: 28px;
  height: 20px;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  padding: 0;
  cursor: pointer;
}

.light-hint {
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px solid #f0f0f0;
  font-size: 11px;
  color: #999;
  line-height: 1.5;
}
</style>
