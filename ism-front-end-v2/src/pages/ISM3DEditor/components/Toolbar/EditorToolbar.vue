<template>
  <div class="editor-header">
    <button class="toolbar-btn" @click="$emit('save')"><i class="fas fa-save"></i> {{ $t('ISM3DEditor.save') }}</button>
    <button class="toolbar-btn" @click="$emit('import')"><i class="fas fa-file-import"></i> {{ $t('ISM3DEditor.import') }}</button>
    <button class="toolbar-btn" @click="$emit('export')"><i class="fas fa-file-export"></i> {{ $t('ISM3DEditor.export') }}</button>
    <button class="toolbar-btn" @click="$emit('load-model')" :title="$t('ISM3DEditor.loadModelTitle')"><i class="fas fa-cube"></i> {{ $t('ISM3DEditor.loadModel') }}</button>
    <div class="toolbar-divider"></div>
    <button class="toolbar-btn" :disabled="!canUndo" @click="$emit('undo')"><i class="fas fa-undo"></i> {{ $t('ISM3DEditor.undo') }}</button>
    <button class="toolbar-btn" :disabled="!canRedo" @click="$emit('redo')"><i class="fas fa-redo"></i> {{ $t('ISM3DEditor.redo') }}</button>
    <div class="toolbar-divider"></div>
    <button class="toolbar-btn" :disabled="!selected" @click="$emit('duplicate')" title="复制 (Ctrl+D)"><i class="fas fa-copy"></i> {{ $t('ISM3DEditor.duplicate') }}</button>
    <button class="toolbar-btn" :disabled="!selected" @click="$emit('focus')" :title="$t('ISM3DEditor.focusObject')"><i class="fas fa-crosshairs"></i> {{ $t('ISM3DEditor.focus') }}</button>
    <button class="toolbar-btn" :disabled="!selected" @click="$emit('delete')" title="删除 (Delete)"><i class="fas fa-trash"></i> {{ $t('ISM3DEditor.delete') }}</button>
    <div class="toolbar-divider"></div>
    <div class="mode-group">
      <button class="toolbar-btn" :class="{active:mode==='select'}" @click="$emit('set-mode','select')" title="选择 (W)"><i class="fas fa-mouse-pointer"></i> {{ $t('ISM3DEditor.select') }}</button>
      <button class="toolbar-btn" :class="{active:mode==='move'}" @click="$emit('set-mode','move')" title="移动 (E)"><i class="fas fa-arrows-alt"></i> {{ $t('ISM3DEditor.move') }}</button>
      <button class="toolbar-btn" :class="{active:mode==='rotate'}" @click="$emit('set-mode','rotate')" title="旋转 (R)"><i class="fas fa-sync-alt"></i> {{ $t('ISM3DEditor.rotate') }}</button>
      <button class="toolbar-btn" :class="{active:mode==='scale'}" @click="$emit('set-mode','scale')" :title="$t('ISM3DEditor.scaleTitle')"><i class="fas fa-expand-arrows-alt"></i> {{ $t('ISM3DEditor.scale') }}</button>
    </div>
    <div class="toolbar-divider"></div>
    <button class="toolbar-btn" :class="{active:showGridSettings}" @click="$emit('toggle-grid-settings')"><i class="fas fa-cog"></i> {{ $t('ISM3DEditor.sceneSettings') }}</button>
    <button class="toolbar-btn" :class="{active:showLightSettings}" @click="$emit('toggle-light-settings')" title="灯光设置"><i class="fas fa-lightbulb"></i> 灯光</button>
    <button class="toolbar-btn" :class="{active:showInspectPanel}" @click="$emit('toggle-inspect')" title="检查"><i class="fas fa-search"></i> 检查</button>
    <div class="toolbar-spacer"></div>
    <button class="toolbar-btn" @click="$emit('templates')"><i class="fas fa-palette"></i> {{ $t('ISM3DEditor.template') }}</button>
    <div class="toolbar-divider"></div>
    <button class="toolbar-btn" @click="$emit('preview')"><i class="fas fa-eye"></i> {{ $t('ISM3DEditor.preview') }}</button>
    <button class="toolbar-btn" @click="$emit('show-help')"><i class="fas fa-question-circle"></i></button>
  </div>
</template>

<script>
export default {
  name: 'EditorToolbar',
  i18n: require('@/i18n/language'),
  props: {
    mode: { type: String, default: 'select' },
    canUndo: { type: Boolean, default: false },
    canRedo: { type: Boolean, default: false },
    showGridSettings: { type: Boolean, default: false },
    showLightSettings: { type: Boolean, default: false },
    showInspectPanel: { type: Boolean, default: false },
    selected: { type: String, default: null },
  }
}
</script>

<style scoped>
.editor-header {
  height: 48px;
  background: #ffffff;
  display: flex;
  align-items: center;
  padding: 0 12px;
  gap: 8px;
  border-bottom: 1px solid #e8e8e8;
  z-index: 100;
  flex-shrink: 0;
}
.editor-logo {
  color: #13c2c2;
  font-size: 18px;
  font-weight: 700;
  margin-right: 12px;
  letter-spacing: 1px;
  white-space: nowrap;
}
.toolbar-divider {
  width: 1px;
  height: 24px;
  background: #e8e8e8;
  margin: 0 6px;
}
.toolbar-btn {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 5px 10px;
  border-radius: 4px;
  border: 1px solid transparent;
  cursor: pointer;
  font-size: 13px;
  color: #333;
  background: transparent;
  transition: all .15s;
  white-space: nowrap;
}
.toolbar-btn:hover { background: #e6fffb; border-color: #87e8de; color: #13c2c2; }
.toolbar-btn.active { background: #e6fffb; border-color: #13c2c2; color: #13c2c2; }
.toolbar-btn[disabled] { opacity: 0.4; cursor: not-allowed; }
.toolbar-btn[disabled]:hover { background: transparent; border-color: transparent; color: #333; }
.toolbar-btn i { font-size: 13px; }
.toolbar-spacer { flex: 1; }
.mode-group { display: flex; border: 1px solid #b5f5ec; border-radius: 4px; overflow: hidden; }
.mode-group .toolbar-btn { border: none; border-radius: 0; border-right: 1px solid #b5f5ec; }
.mode-group .toolbar-btn:last-child { border-right: none; }
</style>
