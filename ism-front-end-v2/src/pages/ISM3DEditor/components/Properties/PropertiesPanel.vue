<template>
  <div class="right-panel">
    <div class="no-select-tip" v-if="!currentObject">
      <i class="fas fa-mouse-pointer"></i>
      点击场景中的对象
      <br>
      查看并编辑属性
    </div>

    <template v-else>
      <a-tabs v-model="activeKey" size="small" class="prop-tabs">
        <a-tab-pane key="style" :tab="$t('ISM3DEditor.style')" :forceRender="true">
          <div class="tab-section">
            <div class="tab-section-title">{{ $t('ISM3DEditor.basicInfo') }}</div>
            <a-form layout="horizontal" :label-col="{ span: 6 }" :wrapper-col="{ span: 17 }" size="small">
              <a-form-item :label="$t('ISM3DEditor.name')">
                <a-input :value="currentObject.name" @change="onPropChange('name', $event)" />
              </a-form-item>
              <a-form-item v-if="currentObject.type === 'gltf'" :label="$t('ISM3DEditor.visible')">
                <a-input :value="currentObject.modelPath || ''" disabled />
              </a-form-item>
              <a-form-item :label="$t('ISM3DEditor.visible')">
                <a-switch :checked="currentObject.visible !== false" @change="onPropSwitch('visible', $event)" size="small" />
              </a-form-item>
            </a-form>
          </div>

          <div class="tab-section" v-if="isTextLabelObject(currentObject)">
            <div class="tab-section-title">{{ $t('ISM3DEditor.labelSettings') }}</div>
            <a-form layout="horizontal" :label-col="{ span: 6 }" :wrapper-col="{ span: 17 }" size="small" class="label-form">
              <a-form-item v-if="currentObject.type !== 'dataText'" :label="$t('ISM3DEditor.labelContent')">
                <a-textarea
                  :value="currentObject.textContent || ''"
                  @change="onPropChange('textContent', $event)"
                  :placeholder="$t('ISM3DEditor.labelContentPlaceholder')"
                  :auto-size="{ minRows: 2, maxRows: 6 }"
                />
              </a-form-item>
              <a-form-item label="文字透明度">
                <a-slider :value="n(currentObject.labelOpacity, 1)" :min="0" :max="1" :step="0.01" @change="onValChange('labelOpacity', $event)" />
              </a-form-item>
              <a-form-item label="字体">
                <a-select :value="currentObject.labelFontFamily || '系统默认'" size="small" @change="onValChange('labelFontFamily', $event)">
                  <a-select-option value="系统默认">系统默认</a-select-option>
                  <a-select-option value="Microsoft YaHei">微软雅黑</a-select-option>
                  <a-select-option value="Arial">Arial</a-select-option>
                  <a-select-option value="黑体">黑体</a-select-option>
                  <a-select-option value="楷体">楷体</a-select-option>
                  <a-select-option value="隶书">隶书</a-select-option>
                  <a-select-option value="宋体">宋体</a-select-option>
                  <a-select-option value="数字字体-1">数字字体-1</a-select-option>
                  <a-select-option value="数字字体-2">数字字体-2</a-select-option>
                  <a-select-option value="数字字体-3">数字字体-3</a-select-option>
                  <a-select-option value="数字字体-4">数字字体-4</a-select-option>
                  <a-select-option value="数字字体-5">数字字体-5</a-select-option>
                  <a-select-option value="数字字体-6">数字字体-6</a-select-option>
                  <a-select-option value="数字字体-7">数字字体-7</a-select-option>
                  <a-select-option value="数字字体-8">数字字体-8</a-select-option>
                  <a-select-option value="数字字体-9">数字字体-9</a-select-option>
                  <a-select-option value="数字字体-10">数字字体-10</a-select-option>
                  <a-select-option value="数字字体-11">数字字体-11</a-select-option>
                  <a-select-option value="数字字体-12">数字字体-12</a-select-option>
                </a-select>
              </a-form-item>
              <a-form-item :label="$t('ISM3DEditor.fontSize')">
                <a-input-number :value="n(currentObject.fontSize, 16)" :min="12" :max="120" :step="1" @change="onValChange('fontSize', $event)" style="width:100%" />
              </a-form-item>
              <a-form-item :label="$t('ISM3DEditor.foregroundColor')">
                <div class="color-row">
                  <input type="color" :value="normColor(currentObject.color)" @input="onColorInput('color', $event)" class="cip" />
                  <div class="color-input-wrap">
                    <a-input :value="currentObject.color" @change="onPropInput('color', $event)" />
                    <button type="button" @click="clearColorValue('color')" class="clear-btn" title="清除"><i class="fas fa-times"></i></button>
                  </div>
                </div>
              </a-form-item>
              <a-form-item :label="$t('ISM3DEditor.backgroundColor')">
                <div class="color-row">
                  <input type="color" :value="normBgColor(currentObject.textBgColor)" @input="onBgColorInput($event)" class="cip" />
                  <div class="color-input-wrap">
                    <a-input :value="currentObject.textBgColor" @change="onBgColorTextInput($event)" placeholder="#00000000" />
                    <button type="button" @click="clearBgColor" class="clear-btn" :title="$t('ISM3DEditor.clearBackgroundColor')"><i class="fas fa-times"></i></button>
                  </div>
                </div>
              </a-form-item>
              <a-form-item :label="$t('ISM3DEditor.opacity')">
                <a-slider :value="n(currentObject.textBgOpacity, 0)" :min="0" :max="1" :step="0.01" @change="onValChange('textBgOpacity', $event)" />
              </a-form-item>
              <a-form-item label="面向相机">
                <a-switch :checked="currentObject.labelFaceCamera !== false" @change="onPropSwitch('labelFaceCamera', $event)" size="small" />
              </a-form-item>
              <a-form-item label="固定大小">
                <a-switch :checked="currentObject.labelFixedSize !== false" @change="onPropSwitch('labelFixedSize', $event)" size="small" />
              </a-form-item>
              <a-form-item label="显示边框">
                <a-switch :checked="!!currentObject.labelShowBorder" @change="onPropSwitch('labelShowBorder', $event)" size="small" />
              </a-form-item>
              <a-form-item label="渲染模式">
                <a-select :value="currentObject.labelRenderMode || 'component'" size="small" @change="onValChange('labelRenderMode', $event)">
                  <a-select-option value="component">组件</a-select-option>
                  <a-select-option value="sprite">精灵</a-select-option>
                </a-select>
              </a-form-item>
              <template v-if="!!currentObject.labelShowBorder">
                <a-form-item label="边框宽度">
                  <a-input-number :value="n(currentObject.labelBorderWidth, 1)" :min="1" :max="8" :step="1" @change="onValChange('labelBorderWidth', $event)" style="width:100%" />
                </a-form-item>
                <a-form-item label="边框颜色">
                  <div class="color-row">
                    <input type="color" :value="normColor(currentObject.labelBorderColor || '#ffffff')" @input="onColorInput('labelBorderColor', $event)" class="cip" />
                    <div class="color-input-wrap">
                      <a-input :value="currentObject.labelBorderColor || '#ffffff'" @change="onPropInput('labelBorderColor', $event)" />
                      <button type="button" @click="clearColorValue('labelBorderColor', '#ffffff')" class="clear-btn" title="Clear"><i class="fas fa-times"></i></button>
                    </div>
                  </div>
                </a-form-item>
              </template>
            </a-form>
          </div>

          <div class="tab-section" v-if="currentObject.type === 'uiLabel'">
            <div class="tab-section-title">屏幕位置</div>
            <a-form layout="horizontal" :label-col="{ span: 6 }" :wrapper-col="{ span: 17 }" size="small">
              <a-form-item label="屏幕 X">
                <a-input-number :value="n(currentObject.uiX, 40)" :min="0" :step="1" @change="onValChange('uiX', $event)" style="width:100%" />
              </a-form-item>
              <a-form-item label="屏幕 Y">
                <a-input-number :value="n(currentObject.uiY, 40)" :min="0" :step="1" @change="onValChange('uiY', $event)" style="width:100%" />
              </a-form-item>
              <a-form-item label="宽度">
                <a-input-number :value="n(currentObject.uiWidth, 180)" :min="24" :step="1" @change="onValChange('uiWidth', $event)" style="width:100%" />
              </a-form-item>
              <a-form-item label="高度">
                <a-input-number :value="n(currentObject.uiHeight, 48)" :min="24" :step="1" @change="onValChange('uiHeight', $event)" style="width:100%" />
              </a-form-item>
            </a-form>
          </div>

          <div class="tab-section" v-if="isBasicMediaObject(currentObject)">
            <div class="tab-section-title">媒体设置</div>
            <a-form layout="horizontal" :label-col="{ span: 6 }" :wrapper-col="{ span: 17 }" size="small">
              <a-form-item v-if="currentObject.type === 'image3d' || currentObject.type === 'uiImage'" label="图片">
                <div class="texture-row">
                  <div class="texture-preview" @click="openTexturePicker" :title="currentObject.mediaUrl || currentObject.imageUrl || currentObject.textureData ? '更换图片' : '选择图片'">
                    <img v-if="currentObject.mediaUrl || currentObject.imageUrl || currentObject.textureData" :src="currentObject.mediaUrl || currentObject.imageUrl || currentObject.textureData" class="texture-thumb" />
                    <i v-else class="fas fa-image" style="font-size:24px;color:#ccc"></i>
                  </div>
                  <div class="texture-actions">
                    <a-button size="small" @click="openTexturePicker" icon="picture">
                      {{ currentObject.mediaUrl || currentObject.imageUrl || currentObject.textureData ? '更换图片' : '选择图片' }}
                    </a-button>
                    <a-button v-if="currentObject.mediaUrl || currentObject.imageUrl || currentObject.textureData" size="small" type="danger" @click="clearMediaImage" icon="delete">清除</a-button>
                  </div>
                </div>
              </a-form-item>
              <a-form-item v-if="currentObject.type === 'video3d'" label="视频地址">
                <a-input :value="currentObject.mediaUrl || currentObject.videoUrl || ''" @change="onMediaUrlChange($event)" placeholder="mp4/webm 地址" />
              </a-form-item>
              <a-form-item v-if="currentObject.type === 'webEmbed'" label="网页地址">
                <a-input :value="currentObject.webUrl || currentObject.mediaUrl || ''" @change="onWebUrlChange($event)" placeholder="https://..." />
              </a-form-item>
              <template v-if="currentObject.type === 'image3d' || currentObject.type === 'video3d'">
                <a-form-item label="平面宽度">
                  <a-input-number :value="n(currentObject.mediaWidth, currentObject.type === 'video3d' ? 1.6 : 1.4)" :min="0.1" :max="20" :step="0.1" @change="onValChange('mediaWidth', $event)" style="width:100%" />
                </a-form-item>
                <a-form-item label="宽高比">
                  <a-input-number :value="n(currentObject.mediaAspect, currentObject.type === 'video3d' ? 1.7778 : 1.3333)" :min="0.1" :max="8" :step="0.01" @change="onValChange('mediaAspect', $event)" style="width:100%" />
                </a-form-item>
              </template>
              <template v-if="currentObject.type === 'uiImage' || currentObject.type === 'webEmbed'">
                <a-form-item label="屏幕 X">
                  <a-input-number :value="n(currentObject.uiX, 40)" :min="0" :step="1" @change="onValChange('uiX', $event)" style="width:100%" />
                </a-form-item>
                <a-form-item label="屏幕 Y">
                  <a-input-number :value="n(currentObject.uiY, 40)" :min="0" :step="1" @change="onValChange('uiY', $event)" style="width:100%" />
                </a-form-item>
                <a-form-item label="宽度">
                  <a-input-number :value="n(currentObject.uiWidth, currentObject.type === 'webEmbed' ? 420 : 180)" :min="24" :step="1" @change="onValChange('uiWidth', $event)" style="width:100%" />
                </a-form-item>
                <a-form-item label="高度">
                  <a-input-number :value="n(currentObject.uiHeight, currentObject.type === 'webEmbed' ? 260 : 100)" :min="24" :step="1" @change="onValChange('uiHeight', $event)" style="width:100%" />
                </a-form-item>
              </template>
              <template v-if="currentObject.type === 'video3d'">
                <a-form-item label="自动播放">
                  <a-switch :checked="currentObject.mediaAutoplay !== false" @change="onPropSwitch('mediaAutoplay', $event)" size="small" />
                </a-form-item>
                <a-form-item label="循环">
                  <a-switch :checked="currentObject.mediaLoop !== false" @change="onPropSwitch('mediaLoop', $event)" size="small" />
                </a-form-item>
                <a-form-item label="静音">
                  <a-switch :checked="currentObject.mediaMuted !== false" @change="onPropSwitch('mediaMuted', $event)" size="small" />
                </a-form-item>
              </template>
            </a-form>
          </div>

          <div class="tab-section" v-if="currentObject.type === 'flowPipe'">
            <div class="tab-section-title">{{ $t('ISM3DEditor.flowPipeSettings') }}</div>
            <a-form layout="horizontal" :label-col="{ span: 6 }" :wrapper-col="{ span: 17 }" size="small">
              <a-form-item :label="$t('ISM3DEditor.pipeColor')">
                <div class="color-row">
                  <input type="color" :value="normColor(currentObject.color)" @input="onColorInput('color', $event)" class="cip" />
                  <div class="color-input-wrap">
                    <a-input :value="currentObject.color" @change="onPropInput('color', $event)" />
                    <button type="button" @click="clearColorValue('color')" class="clear-btn" title="清除"><i class="fas fa-times"></i></button>
                  </div>
                </div>
              </a-form-item>
              <a-form-item :label="$t('ISM3DEditor.highlightColor')">
                <div class="color-row">
                  <input type="color" :value="normColor(currentObject.highlightColor || '#ff6a00')" @input="onColorInput('highlightColor', $event)" class="cip" />
                  <div class="color-input-wrap">
                    <a-input :value="currentObject.highlightColor || '#ff6a00'" @change="onPropInput('highlightColor', $event)" />
                    <button type="button" @click="clearColorValue('highlightColor', '#ff6a00')" class="clear-btn" title="清除"><i class="fas fa-times"></i></button>
                  </div>
                </div>
              </a-form-item>
              <a-form-item :label="$t('ISM3DEditor.pipeRadius')">
                <a-input-number :value="n(currentObject.radius, 0.1)" :min="0.01" :max="2" :step="0.01" @change="onValChange('radius', $event)" style="width:100%" />
              </a-form-item>
              <a-form-item :label="$t('ISM3DEditor.flowSpeed')">
                <a-slider :value="n(currentObject.flowSpeed, 1.0)" :min="0" :max="5" :step="0.1" @change="onValChange('flowSpeed', $event)" />
              </a-form-item>
              <a-form-item :label="$t('ISM3DEditor.flowDashLength')">
                <a-input-number :value="n(currentObject.flowDashLength, 3)" :min="0.05" :max="3" :step="0.05" @change="onValChange('flowDashLength', $event)" style="width:100%" />
              </a-form-item>
              <a-form-item :label="$t('ISM3DEditor.flowDirection')">
                <a-select :value="currentObject.flowDirection || 'forward'" @change="onValChange('flowDirection', $event)" style="width:100%">
                  <a-select-option value="forward">{{ $t('ISM3DEditor.forward') }}</a-select-option>
                  <a-select-option value="backward">{{ $t('ISM3DEditor.backward') }}</a-select-option>
                </a-select>
              </a-form-item>
            </a-form>
          </div>

          <div class="tab-section">
            <div class="tab-section-title">{{ $t('ISM3DEditor.transform') }}</div>
            <a-form layout="horizontal" :label-col="{ span: 6 }" :wrapper-col="{ span: 17 }" size="small">
              <a-form-item :label="$t('ISM3DEditor.positionX')"><a-input-number :value="n(currentObject.x, 0)" @change="onValChange('x', $event)" :step="0.1" style="width:100%" /></a-form-item>
              <a-form-item :label="$t('ISM3DEditor.positionY')"><a-input-number :value="n(currentObject.y, 0)" @change="onValChange('y', $event)" :step="0.1" style="width:100%" /></a-form-item>
              <a-form-item :label="$t('ISM3DEditor.positionZ')"><a-input-number :value="n(currentObject.z, 0)" @change="onValChange('z', $event)" :step="0.1" style="width:100%" /></a-form-item>
              <a-form-item :label="$t('ISM3DEditor.rotationX')"><a-input-number :value="n(currentObject.rx, 0)" @change="onValChange('rx', $event)" :step="1" style="width:100%" /></a-form-item>
              <a-form-item :label="$t('ISM3DEditor.rotationY')"><a-input-number :value="n(currentObject.ry, 0)" @change="onValChange('ry', $event)" :step="1" style="width:100%" /></a-form-item>
              <a-form-item :label="$t('ISM3DEditor.rotationZ')"><a-input-number :value="n(currentObject.rz, 0)" @change="onValChange('rz', $event)" :step="1" style="width:100%" /></a-form-item>
              <a-form-item :label="$t('ISM3DEditor.scaleX')"><a-input-number :value="n(currentObject.sx, 1)" @change="onValChange('sx', $event)" :step="0.1" :min="0.01" style="width:100%" /></a-form-item>
              <a-form-item :label="$t('ISM3DEditor.scaleY')"><a-input-number :value="n(currentObject.sy, 1)" @change="onValChange('sy', $event)" :step="0.1" :min="0.01" style="width:100%" /></a-form-item>
              <a-form-item :label="$t('ISM3DEditor.scaleZ')"><a-input-number :value="n(currentObject.sz, 1)" @change="onValChange('sz', $event)" :step="0.1" :min="0.01" style="width:100%" /></a-form-item>
            </a-form>
          </div>
        </a-tab-pane>

        <a-tab-pane key="material" :tab="$t('ISM3DEditor.material')" :forceRender="true">
          <div class="tab-section material-section">
            <a-form layout="horizontal" :label-col="{ span: 5 }" :wrapper-col="{ span: 18 }" size="small" class="material-form">
              <a-form-item :label="$t('ISM3DEditor.color')">
                <div class="color-row">
                  <input type="color" :value="normColor(currentObject.color)" @input="onColorInput('color', $event)" class="cip" />
                  <div class="color-input-wrap">
                    <a-input :value="currentObject.color" @change="onPropInput('color', $event)" />
                    <button type="button" @click="clearColorValue('color')" class="clear-btn" title="清除"><i class="fas fa-times"></i></button>
                  </div>
                </div>
              </a-form-item>
              <a-form-item :label="$t('ISM3DEditor.opacity')">
                <a-slider :value="n(currentObject.opacity, 1)" :min="0" :max="1" :step="0.01" @change="onValChange('opacity', $event)" />
              </a-form-item>
              <a-form-item :label="$t('ISM3DEditor.metalness')">
                <a-slider :value="n(currentObject.metalness, 0.3)" :min="0" :max="1" :step="0.01" @change="onValChange('metalness', $event)" />
              </a-form-item>
              <a-form-item :label="$t('ISM3DEditor.roughness')">
                <a-slider :value="n(currentObject.roughness, 0.7)" :min="0" :max="1" :step="0.01" @change="onValChange('roughness', $event)" />
              </a-form-item>
              <a-form-item :label="$t('ISM3DEditor.emissive')">
                <div class="color-row">
                  <input type="color" :value="normColor(currentObject.emissive)" @input="onColorInput('emissive', $event)" class="cip" />
                  <div class="color-input-wrap">
                    <a-input :value="currentObject.emissive" @change="onPropInput('emissive', $event)" />
                    <button type="button" @click="clearColorValue('emissive', '#000000')" class="clear-btn" title="清除"><i class="fas fa-times"></i></button>
                  </div>
                </div>
              </a-form-item>
              <a-form-item :label="$t('ISM3DEditor.wireframe')">
                <a-switch :checked="!!currentObject.wireframe" @change="onPropSwitch('wireframe', $event)" size="small" />
              </a-form-item>
              <a-form-item :label="$t('ISM3DEditor.showShadow')">
                <a-switch :checked="!!currentObject.showShadow" @change="onPropSwitch('showShadow', $event)" size="small" />
              </a-form-item>
              <a-divider style="margin:4px 0 6px" />
              <a-form-item :label="$t('ISM3DEditor.texture')">
                <div class="texture-row">
                  <div class="texture-preview" @click="openTexturePicker" :title="currentObject.textureData ? $t('ISM3DEditor.changeTexture') : $t('ISM3DEditor.selectTexture')">
                    <img v-if="currentObject.textureData" :src="currentObject.textureData" class="texture-thumb" />
                    <i v-else class="fas fa-image" style="font-size:24px;color:#ccc"></i>
                  </div>
                  <div class="texture-actions">
                    <a-button size="small" @click="openTexturePicker" icon="picture">
                      {{ currentObject.textureData ? $t('ISM3DEditor.change') : $t('ISM3DEditor.selectImage') }}
                    </a-button>
                    <a-button v-if="currentObject.textureData" size="small" type="danger" @click="clearTexture" icon="delete">{{ $t('ISM3DEditor.clear') }}</a-button>
                  </div>
                </div>
              </a-form-item>
            </a-form>
          </div>
        </a-tab-pane>

        <a-tab-pane key="animate" :tab="$t('ISM3DEditor.animation')" :forceRender="true">
          <div class="tab-section animate-section">
            <div class="tab-section-title">{{ $t('ISM3DEditor.animation') }}</div>
            <a-form layout="horizontal" :label-col="{ span: 6 }" :wrapper-col="{ span: 17 }" size="small" class="animate-form">
              <a-form-item :label="$t('ISM3DEditor.selfRotate')">
                <a-switch :checked="!!currentObject.autoRotate" @change="onPropSwitch('autoRotate', $event)" />
              </a-form-item>
              <div v-if="!!currentObject.autoRotate">
                <a-form-item :label="$t('ISM3DEditor.speed')">
                  <a-slider :value="n(currentObject.rotateSpeed, 1)" :min="-5" :max="5" :step="0.1" @change="onAnimVal('rotateSpeed', $event)" />
                </a-form-item>
                <a-form-item :label="$t('ISM3DEditor.axis')">
                  <a-select :value="currentObject.rotateAxis || 'y'" @change="onAnimVal('rotateAxis', $event)" style="width:100%">
                    <a-select-option value="y">{{ $t('ISM3DEditor.yAxis') }}</a-select-option>
                    <a-select-option value="x">{{ $t('ISM3DEditor.xAxis') }}</a-select-option>
                    <a-select-option value="z">{{ $t('ISM3DEditor.zAxis') }}</a-select-option>
                  </a-select>
                </a-form-item>
              </div>
              <a-divider style="margin:4px 0 6px" />
              <a-form-item :label="$t('ISM3DEditor.float')">
                <a-switch :checked="!!currentObject.floatAnim" @change="onPropSwitch('floatAnim', $event)" />
              </a-form-item>
              <div v-if="!!currentObject.floatAnim">
                <a-form-item :label="$t('ISM3DEditor.range')">
                  <a-slider :value="n(currentObject.floatRange, 0.15)" :min="0.01" :max="1" :step="0.01" @change="onAnimVal('floatRange', $event)" />
                </a-form-item>
                <a-form-item :label="$t('ISM3DEditor.speed')">
                  <a-slider :value="n(currentObject.floatSpeed, 2)" :min="0.5" :max="5" :step="0.1" @change="onAnimVal('floatSpeed', $event)" />
                </a-form-item>
              </div>
              <a-divider style="margin:4px 0 6px" />
              <a-form-item :label="$t('ISM3DEditor.blink')">
                <a-switch :checked="!!currentObject.blink" @change="onPropSwitch('blink', $event)" />
              </a-form-item>
              <div v-if="!!currentObject.blink">
                <a-form-item :label="$t('ISM3DEditor.speed')">
                  <a-slider :value="n(currentObject.blinkSpeed, 6)" :min="1" :max="10" :step="0.5" @change="onAnimVal('blinkSpeed', $event)" />
                </a-form-item>
                <a-form-item :label="$t('ISM3DEditor.minimum')">
                  <a-slider :value="n(currentObject.blinkMin, 0.2)" :min="0" :max="0.8" :step="0.05" @change="onAnimVal('blinkMin', $event)" />
                </a-form-item>
              </div>
              <a-divider style="margin:4px 0 6px" v-if="currentObject.type === 'flowPipe'" />
              <a-form-item :label="$t('ISM3DEditor.flow')" v-if="currentObject.type === 'flowPipe'">
                <a-switch :checked="!!currentObject.flowAnim" @change="onPropSwitch('flowAnim', $event)" />
              </a-form-item>
            </a-form>
          </div>

          <div class="tab-section animate-section">
            <div class="tab-section-title">{{ $t('ISM3DEditor.animationBinding') }}</div>
            <a-form layout="horizontal" :label-col="{ span: 6 }" :wrapper-col="{ span: 17 }" size="small" class="animate-form animate-bind-form">
              <a-form-item :label="$t('ISM3DEditor.enableBinding')">
                <a-switch :checked="!!(currentObject.animate && currentObject.animate.condition.isBandDevice)" @change="onAnimateBindChange('isBandDevice', $event)" />
              </a-form-item>
              <template v-if="currentObject.animate && currentObject.animate.condition.isBandDevice">
                <a-form-item :label="$t('ISM3DEditor.bindAnimation')">
                  <a-checkbox-group class="bind-animation-group" :value="currentObject.animate.selected || []" @change="onAnimateSelectedChange">
                    <a-checkbox value="autoRotate">{{ $t('ISM3DEditor.selfRotate') }}</a-checkbox>
                    <a-checkbox value="floatAnim">{{ $t('ISM3DEditor.float') }}</a-checkbox>
                    <a-checkbox value="blink">{{ $t('ISM3DEditor.blink') }}</a-checkbox>
                    <a-checkbox value="visible">{{ $t('ISM3DEditor.visible') }}</a-checkbox>
                    <a-checkbox value="flowAnim" v-if="currentObject.type === 'flowPipe'">{{ $t('ISM3DEditor.flow') }}</a-checkbox>
                  </a-checkbox-group>
                </a-form-item>
                <a-form-item :label="$t('ISM3DEditor.deviceName')">
                  <a-input :value="currentObject.animate.condition.DeviceName || ''" @change="onAnimateBindChange('deviceSN', $event.target.value)" />
                </a-form-item>
                <a-form-item :label="$t('ISM3DEditor.dataName')">
                  <a-input-group compact class="data-select-group">
                    <a-input :value="currentObject.animate.condition.dataName || ''" disabled />
                    <a-button type="primary" @click="openDataSelector('animate')">
                      <i class="fas fa-database"></i> {{ currentObject.animate.condition.dataName ? $t('ISM3DEditor.reselect') : $t('ISM3DEditor.selectData') }}
                    </a-button>
                  </a-input-group>
                </a-form-item>
                <a-form-item :label="$t('ISM3DEditor.operator')">
                  <a-select :value="currentObject.animate.condition.operator || ''" @change="onAnimateBindChange('operator', $event)">
                    <a-select-option value="==">==</a-select-option>
                    <a-select-option value=">">></a-select-option>
                    <a-select-option value=">=">>=</a-select-option>
                    <a-select-option value="<">&lt;</a-select-option>
                    <a-select-option value="<=">&lt;=</a-select-option>
                    <a-select-option value="!=">!=</a-select-option>
                    <a-select-option value="<>">&lt;&gt;</a-select-option>
                    <a-select-option value="<!>">&lt;!&gt;</a-select-option>
                  </a-select>
                </a-form-item>
                <a-form-item :label="$t('ISM3DEditor.comparisonValue')">
                  <a-input-number :value="n(currentObject.animate.condition.OperatorValue, 0)" :step="0.1" style="width:100%" @change="onAnimateBindChange('OperatorValue', $event)" />
                </a-form-item>
              </template>
            </a-form>
          </div>
        </a-tab-pane>

        <a-tab-pane key="modelAnimation" :tab="$t('ISM3DEditor.modelAnimation')" v-if="currentObject.type === 'gltf' || currentObject.isExternalModel" :forceRender="true">
          <div class="tab-section">
            <div class="tab-section-title model-animation-title">
              <span>{{ $t('ISM3DEditor.modelAnimation') }}</span>
              <a-button size="small" type="primary" ghost @click="addGLTFAnimationGroup">
                <i class="fas fa-plus"></i> {{ $t('ISM3DEditor.addGroup') }}              </a-button>
            </div>
            <div v-if="gltfAnimationOptions.length === 0" class="model-animation-empty">{{ $t('ISM3DEditor.noBuiltInAnimation') }}</div>
            <div v-for="group in gltfAnimationGroups" :key="group.id" class="model-animation-group">
              <div class="model-animation-group-header">
                <a-input :value="group.name" size="small" class="model-animation-group-name" @change="setGLTFAnimationGroupProp(group, 'name', $event.target.value)" />
                <a-switch :checked="!!group.playing" size="small" @change="setGLTFAnimationGroupProp(group, 'playing', $event)" />
                <a-button size="small" type="link" :disabled="gltfAnimationGroups.length <= 1" @click="removeGLTFAnimationGroup(group.id)">
                  {{ $t('ISM3DEditor.delete') }}
                </a-button>
              </div>
              <a-form layout="horizontal" :label-col="{ span: 6 }" :wrapper-col="{ span: 17 }" size="small">
                <a-form-item label="动画片段">
                  <a-select mode="multiple" :value="group.animationNames || []" :disabled="gltfAnimationOptions.length === 0" @change="onGLTFAnimationGroupNamesChange(group, $event)" style="width:100%">
                    <a-select-option v-for="item in gltfAnimationOptions" :key="item.key" :value="item.key">{{ item.label }}</a-select-option>
                  </a-select>
                </a-form-item>
                <a-form-item label="播放速度">
                  <a-slider :value="n(group.speed, 1)" :min="0" :max="5" :step="0.1" @change="setGLTFAnimationGroupProp(group, 'speed', $event)" />
                </a-form-item>
                <a-form-item label="循环播放">
                  <a-switch :checked="group.loop !== false" @change="setGLTFAnimationGroupProp(group, 'loop', $event)" />
                </a-form-item>
                <a-form-item label="条件控制">
                  <a-switch :checked="!!group.conditionEnabled" @change="setGLTFAnimationGroupProp(group, 'conditionEnabled', $event)" />
                </a-form-item>
                <template v-if="group.conditionEnabled">
                  <a-form-item label="动画数据">
                    <a-input-group compact class="data-select-group">
                      <a-input :value="(group.condition && group.condition.dataName) || ''" disabled />
                      <a-button type="primary" @click="openGLTFAnimationGroupDataSelector(group)">
                        <i class="fas fa-database"></i> {{ group.condition && group.condition.dataName ? '重新选择' : '选择数据' }}
                      </a-button>
                    </a-input-group>
                  </a-form-item>
                  <a-form-item label="绑定设备">
                    <a-switch :checked="!!(group.condition && group.condition.isBandDevice)" @change="setGLTFAnimationGroupCondition(group, 'isBandDevice', $event)" />
                  </a-form-item>
                  <a-form-item :label="$t('ISM3DEditor.deviceName')" v-if="group.condition && group.condition.isBandDevice">
                    <a-input :value="group.condition.DeviceName || group.condition.deviceSN || ''" @change="setGLTFAnimationGroupCondition(group, 'deviceSN', $event.target.value)" />
                  </a-form-item>
                  <a-form-item :label="$t('ISM3DEditor.operator')">
                    <a-select :value="(group.condition && group.condition.operator) || ''" allowClear @change="setGLTFAnimationGroupCondition(group, 'operator', $event || '')">
                      <a-select-option value="==">==</a-select-option>
                      <a-select-option value=">">></a-select-option>
                      <a-select-option value=">=">>=</a-select-option>
                      <a-select-option value="<">&lt;</a-select-option>
                      <a-select-option value="<=">&lt;=</a-select-option>
                      <a-select-option value="!=">!=</a-select-option>
                      <a-select-option value="<>">&lt;&gt;</a-select-option>
                      <a-select-option value="<!>">&lt;!&gt;</a-select-option>
                    </a-select>
                  </a-form-item>
                  <a-form-item :label="$t('ISM3DEditor.comparisonValue')" v-if="group.condition && group.condition.operator">
                    <a-input-number :value="n(group.condition.OperatorValue, 0)" :step="0.1" style="width:100%" @change="setGLTFAnimationGroupCondition(group, 'OperatorValue', $event)" />
                  </a-form-item>
                  <a-form-item label="最大值" v-if="group.condition && (group.condition.operator === '<>' || group.condition.operator === '<!>')">
                    <a-input-number :value="n(group.condition.OperatorMaxValue, 0)" :step="0.1" style="width:100%" @change="setGLTFAnimationGroupCondition(group, 'OperatorMaxValue', $event)" />
                  </a-form-item>
                </template>
              </a-form>
            </div>
          </div>
        </a-tab-pane>

        <a-tab-pane key="events" :tab="$t('displayConfig.Properties.TabHeaterBehavior')" :forceRender="true">
          <div class="tab-section event-section">
            <div v-for="(event, index) in currentObject.action" :key="index" class="action-card">
              <div class="action-card-head">
                <span>{{ $t('displayConfig.Properties.TabHeaterBehavior') }}-{{ index + 1 }}</span>
                <a-icon type="delete" @click="delBindAction(index)" theme="twoTone" two-tone-color="#eb2f96" class="action-delete" />
              </div>
              <a-form layout="horizontal" :label-col="{ span: 7 }" :wrapper-col="{ span: 16 }" size="small">
                <a-form-item :label="$t('displayConfig.Properties.ComponentEvent')">
                  <a-select v-model="event.type" @change="emitActionChange">
                    <a-select-option v-for="option in actionEventOptions" :key="option.value" :value="option.value">{{ $t(option.label) }}</a-select-option>
                  </a-select>
                </a-form-item>
                <a-form-item :label="$t('displayConfig.Properties.ComponentActionAuth')">
                  <a-select mode="multiple" v-model="event.actionAuth" @change="emitActionChange">
                    <a-select-option v-for="option in roleOptions" :key="option.RoleId" :value="option.RoleId">{{ option.Name }}</a-select-option>
                  </a-select>
                </a-form-item>
                <a-form-item :label="$t('displayConfig.Properties.ComponentActionPassword')">
                  <a-input type="text" v-model="event.ActionPassword" :style="{'text-security':'disc', '-webkit-text-security':'disc'}" @change="emitActionChange" />
                </a-form-item>
                <a-form-item :label="$t('displayConfig.Properties.ComponentActionVoice')">
                  <a-input v-model="event.actionVoice" @change="emitActionChange" />
                </a-form-item>
                <a-form-item :label="$t('displayConfig.Properties.ComponentAction')">
                  <a-select v-model="event.action" @change="changeAction($event, index)">
                    <a-select-option v-for="option in actionTypeOptions" :key="option.value" :value="option.value">{{ $t(option.label) }}</a-select-option>
                  </a-select>
                </a-form-item>
                <a-form-item :label="$t('displayConfig.Properties.SecondConfirm')">
                  <a-switch v-model="event.actionConfirm" @change="emitActionChange">
                    <a-icon slot="checkedChildren" type="check" />
                    <a-icon slot="unCheckedChildren" type="close" />
                  </a-switch>
                </a-form-item>
                <div v-if="event.action === 'visible'">
                  <a-form-item :label="$t('displayConfig.Properties.ComponentClickShow')">
                    <a-select allowClear mode="multiple" v-model="event.showItems" @change="emitActionChange">
                      <a-select-option v-for="obj in sceneObjects" :key="'show_' + obj.id" :value="obj.id">{{ obj.name || obj.typeName || obj.id }}</a-select-option>
                    </a-select>
                  </a-form-item>
                  <a-form-item :label="$t('displayConfig.Properties.ComponentClickHide')">
                    <a-select allowClear mode="multiple" v-model="event.hideItems" @change="emitActionChange">
                      <a-select-option v-for="obj in sceneObjects" :key="'hide_' + obj.id" :value="obj.id">{{ obj.name || obj.typeName || obj.id }}</a-select-option>
                    </a-select>
                  </a-form-item>
                </div>
                <div v-if="event.action === 'link'">
                  <a-form-item :label="$t('displayConfig.Properties.linkType')">
                    <a-select v-model="event.link.linkType" @change="emitActionChange" allowClear>
                      <a-select-option value="Inside">{{ $t('displayConfig.Properties.linkInside') }}</a-select-option>
                      <a-select-option value="External">{{ $t('displayConfig.Properties.linkExternal') }}</a-select-option>
                    </a-select>
                  </a-form-item>
                  <a-form-item v-if="event.link.linkType === 'Inside'" :label="$t('displayConfig.Properties.linkIAppUUID')">
                    <a-select v-model="event.link.Inside.displayUUID" allowClear @change="changeLinkInsideDisplay(event.link)">
                      <a-select-option v-for="option in configurationModel" :key="option.uuid" :value="option.uuid">
                        {{ option.name }}{{ option.displayType == 2 ? '（3D场景）' : '' }}
                      </a-select-option>
                    </a-select>
                  </a-form-item>
                  <a-form-item v-if="event.link.linkType === 'Inside'" :label="$t('displayConfig.Properties.linkIAppPageUUID')">
                    <a-select v-model="event.link.Inside.pageUUID" allowClear @change="emitActionChange">
                      <a-select-option v-for="option in generateTargetPage(event.link.Inside.displayUUID)" :key="option.value" :value="option.value">
                        {{ option.label }}
                      </a-select-option>
                    </a-select>
                  </a-form-item>
                  <a-form-item v-if="event.link.linkType !== 'Inside'" label="网页来源">
                    <a-select :value="event.link.ExternalSource || 'Url'" @change="value => changeLinkExternalSource(event.link, value)">
                      <a-select-option value="Url">输入网址</a-select-option>
                      <a-select-option value="Page">选择页面</a-select-option>
                    </a-select>
                  </a-form-item>
                  <template v-if="event.link.linkType !== 'Inside' && event.link.ExternalSource === 'Page'">
                    <a-form-item label="应用">
                      <a-select v-model="event.link.ExternalPage.displayUUID" allowClear @change="changeLinkExternalDisplay(event.link)">
                        <a-select-option v-for="option in get2DConfigurationModel()" :key="option.uuid" :value="option.uuid">
                          {{ option.name }}
                        </a-select-option>
                      </a-select>
                    </a-form-item>
                    <a-form-item label="页面">
                      <a-select v-model="event.link.ExternalPage.pageUUID" allowClear @change="changeLinkExternalPage(event.link)">
                        <a-select-option v-for="option in generateTargetPage(event.link.ExternalPage.displayUUID)" :key="option.value" :value="option.value">
                          {{ option.label }}
                        </a-select-option>
                      </a-select>
                    </a-form-item>
                  </template>
                  <a-form-item v-if="event.link.linkType !== 'Inside'" :label="$t('displayConfig.Properties.linkExternalUrl')">
                    <a-input v-model="event.link.External" :disabled="event.link.ExternalSource === 'Page'" @change="emitActionChange" />
                  </a-form-item>
                  <a-form-item v-if="event.link.linkType !== 'Inside'" :label="$t('displayConfig.Properties.OpenLinkExternalType')">
                    <a-select v-model="event.link.OpenExternalType" @change="emitActionChange">
                      <a-select-option value="self">{{ $t('displayConfig.Properties.OpenLinkExternalSelf') }}</a-select-option>
                      <a-select-option value="new">{{ $t('displayConfig.Properties.OpenLinkExternalNew') }}</a-select-option>
                    </a-select>
                  </a-form-item>
                  <a-form-item :label="$t('displayConfig.Properties.isLinkPopUp')">
                    <a-checkbox v-model="event.link.isPopUp" @change="emitActionChange">{{ $t('displayConfig.Properties.isLinkPopUp') }}</a-checkbox>
                  </a-form-item>
                  <a-form-item :label="$t('displayConfig.Properties.autoClose')" v-if="event.link.isPopUp">
                    <a-checkbox v-model="event.link.autoClose" @change="emitActionChange">{{ $t('displayConfig.Properties.autoClose') }}</a-checkbox>
                  </a-form-item>
                  <a-form-item :label="$t('displayConfig.Properties.linkExternalWidth')" v-if="event.link.isPopUp">
                    <a-input v-model="event.link.width" @change="emitActionChange" />
                  </a-form-item>
                  <a-form-item :label="$t('displayConfig.Properties.linkExternalHeight')" v-if="event.link.isPopUp">
                    <a-input v-model="event.link.height" @change="emitActionChange" />
                  </a-form-item>
                  <a-form-item :label="$t('displayConfig.Properties.linkExternalTitle')" v-if="event.link.isPopUp">
                    <a-input v-model="event.link.title" @change="emitActionChange" />
                  </a-form-item>
                </div>
                <div v-if="event.action === 'SetValue'">
                  <a-form-item :label="$t('displayConfig.Properties.SetDelay')">
                    <a-input v-model="event.SetDelay" @change="emitActionChange" />
                  </a-form-item>
                  <div class="set-value-actions">
                    <a-button size="small" @click="addBindSetValue(index)">{{ $t('displayConfig.Properties.AddSetValue') }}</a-button>
                  </div>
                  <div v-for="(setValueItem, setValueIndex) in event.setValue" :key="setValueIndex" class="set-value-card">
                    <div class="set-value-head">
                      <span>{{ $t('displayConfig.Properties.SetValueInfo') }}-{{ setValueIndex + 1 }}</span>
                      <a-icon type="close-circle" @click="delBindSetValue(index, setValueIndex)" theme="twoTone" two-tone-color="#eb2f96" class="action-delete" />
                    </div>
                    <a-form-item :label="$t('displayConfig.Properties.IsManualLabel')">
                      <a-checkbox v-model="setValueItem.IsManual" @change="emitActionChange">{{ $t('displayConfig.Properties.IsManual') }}</a-checkbox>
                    </a-form-item>
                    <a-form-item :label="$t('displayConfig.Properties.ComponentBandData')">
                      <a-input v-model="setValueItem.dataName" @change="emitActionChange">
                        <a-tooltip placement="top" slot="addonAfter">
                          <template slot="title">
                            <span>{{ $t('displayConfig.Properties.SelectValue') }}</span>
                          </template>
                          <i class="fas fa-database" @click="openActionSetValueSelector(index, setValueIndex)" />
                        </a-tooltip>
                      </a-input>
                    </a-form-item>
                    <a-form-item :label="$t('displayConfig.Properties.AutoSetValue')" v-if="!setValueItem.IsManual">
                      <a-input v-model="setValueItem.AutoSetValue" @change="emitActionChange" />
                    </a-form-item>
                    <a-form-item :label="$t('displayConfig.Properties.SetPassword')">
                      <a-input v-model="setValueItem.SetPassword" @change="emitActionChange" />
                    </a-form-item>
                  </div>
                </div>
                <div v-if="event.action === 'RestApi'">
                  <a-form-item :label="$t('displayConfig.Properties.RestApiName')">
                    <a-input v-model="event.RestApi.Name" @change="emitActionChange" />
                  </a-form-item>
                  <a-form-item :label="$t('displayConfig.Properties.RestApiFrom')">
                    <a-select v-model="event.RestApi.IsSystem" @change="emitActionChange">
                      <a-select-option value="1">{{ $t('displayConfig.Properties.ExternRestApi') }}</a-select-option>
                      <a-select-option value="2">{{ $t('displayConfig.Properties.SystemRestApi') }}</a-select-option>
                    </a-select>
                  </a-form-item>
                  <a-form-item :label="$t('displayConfig.Properties.RestApiType')">
                    <a-select v-model="event.RestApi.Type" @change="emitActionChange">
                      <a-select-option value="Post">Post</a-select-option>
                      <a-select-option value="Get">Get</a-select-option>
                    </a-select>
                  </a-form-item>
                  <a-form-item :label="$t('displayConfig.Properties.RestApiUrl')">
                    <a-input v-model="event.RestApi.Url" @change="emitActionChange" />
                  </a-form-item>
                  <a-form-item :label="$t('displayConfig.Properties.RestApiParam')">
                    <a-textarea :rows="6" v-model="event.RestApi.Params" @change="emitActionChange" />
                  </a-form-item>
                </div>
                <div v-if="event.action === 'SysScript'">
                  <a-form-item :label="$t('displayConfig.Properties.action.SysScript')">
                    <a-select mode="multiple" v-model="event.ScriptList" @change="emitActionChange" />
                  </a-form-item>
                </div>
                <div v-if="event.action === 'DeviceView'">
                  <a-form-item :label="$t('displayConfig.Properties.SelectDevice')">
                    <a-input v-model="event.DeviceView.key" @change="emitActionChange" />
                  </a-form-item>
                  <a-form-item :label="$t('displayConfig.Properties.isLinkPopUp')">
                    <a-checkbox v-model="event.DeviceView.isPopUp" @change="emitActionChange">{{ $t('displayConfig.Properties.isLinkPopUp') }}</a-checkbox>
                  </a-form-item>
                  <a-form-item :label="$t('displayConfig.Properties.isContainer')">
                    <a-checkbox v-model="event.DeviceView.isContainer" @change="emitActionChange">{{ $t('displayConfig.Properties.isContainer') }}</a-checkbox>
                  </a-form-item>
                </div>
                <div v-if="event.action === 'Animation'">
                  <a-form-item :label="$t('displayConfig.Properties.action.animationStatus')">
                    <a-select v-model="event.animationStatus" @change="emitActionChange" allowClear>
                      <a-select-option value="start">{{ $t('displayConfig.Properties.action.animationStart') }}</a-select-option>
                      <a-select-option value="stop">{{ $t('displayConfig.Properties.action.animationStop') }}</a-select-option>
                    </a-select>
                  </a-form-item>
                </div>
              </a-form>
            </div>
            <div class="event-add">
              <a-button block type="primary" @click="addBindAction">{{ $t('displayConfig.Properties.ComponentAddBandActionBtn') }}</a-button>
            </div>
          </div>
        </a-tab-pane>

        <a-tab-pane key="data" tab="数据" :forceRender="true">
          <div class="tab-section data-bind-section" v-if="showGenericDataBinding">
            <div class="tab-section-title">{{ $t('ISM3DEditor.dataBinding') }}</div>
            <a-form layout="horizontal" :label-col="{ span: 5 }" :wrapper-col="{ span: 18 }" size="small" class="data-bind-form">
              <a-form-item :label="$t('ISM3DEditor.dataName')">
                <a-input-group compact class="data-select-group">
                  <a-input :value="currentObject.dataName || ''" disabled />
                  <a-button type="primary" @click="openDataSelector('binding')">
                    <i class="fas fa-database"></i> {{ currentObject.dataName ? '重新选择' : '选择数据' }}
                  </a-button>
                </a-input-group>
              </a-form-item>
              <a-form-item label="绑定属性">
                <a-select :value="currentObject.bindProp || ''" @change="onBindPropChange" allowClear>
                  <a-select-option v-for="option in bindPropOptions" :key="option.value" :value="option.value">{{ option.label }}</a-select-option>
                </a-select>
              </a-form-item>
              <a-form-item v-if="showBindTransformControls" :label="$t('ISM3DEditor.transformType')">
                <a-select :value="currentObject.bindTransform || 'direct'" @change="onBindTransformChange">
                  <a-select-option value="direct">{{ $t('ISM3DEditor.direct') }}</a-select-option>
                  <a-select-option value="scale">{{ $t('ISM3DEditor.scale') }}</a-select-option>
                  <a-select-option value="offset">{{ $t('ISM3DEditor.offsetTransform') }}</a-select-option>
                </a-select>
              </a-form-item>
              <a-form-item v-if="currentObject.bindProp && currentObject.bindTransform==='scale'" label="系数">
                <a-input-number :value="n(currentObject.bindScale, 1)" :min="0" :step="0.1" style="width:100%" @change="onBindScaleChange" />
              </a-form-item>
              <a-form-item v-if="showBindTransformControls && currentObject.bindTransform==='offset'" :label="$t('ISM3DEditor.offsetTransform')">
                <a-input-number :value="n(currentObject.bindOffset, 0)" :step="1" style="width:100%" @change="onBindOffsetChange" />
              </a-form-item>
            </a-form>
          </div>

          <div class="tab-section data-bind-section" v-if="showPositionDataBinding">
            <div class="tab-section-title">三坐标绑定</div>
            <a-form layout="horizontal" :label-col="{ span: 5 }" :wrapper-col="{ span: 18 }" size="small" class="data-bind-form">
              <template v-for="axis in positionBindingAxes">
                <a-form-item :key="axis.key + '_data'" :label="axis.label">
                  <a-input-group compact class="data-select-group">
                    <a-input :value="(currentObject.positionBindings && currentObject.positionBindings[axis.key] && currentObject.positionBindings[axis.key].dataName) || ''" disabled />
                    <a-button type="primary" @click="openPositionBindingSelector(axis.key)">
                      <i class="fas fa-database"></i> {{ currentObject.positionBindings && currentObject.positionBindings[axis.key] && currentObject.positionBindings[axis.key].dataName ? '重新选择' : '选择数据' }}
                    </a-button>
                  </a-input-group>
                </a-form-item>
                <a-form-item :key="axis.key + '_transform'" label="变换">
                  <a-select :value="getPositionBinding(axis.key).transform || 'direct'" @change="onPositionBindingChange(axis.key, 'transform', $event)">
                    <a-select-option value="direct">{{ $t('ISM3DEditor.direct') }}</a-select-option>
                    <a-select-option value="scale">{{ $t('ISM3DEditor.scale') }}</a-select-option>
                    <a-select-option value="offset">{{ $t('ISM3DEditor.offsetTransform') }}</a-select-option>
                  </a-select>
                </a-form-item>
                <a-form-item v-if="getPositionBinding(axis.key).transform === 'scale'" :key="axis.key + '_scale'" label="系数">
                  <a-input-number :value="n(getPositionBinding(axis.key).scale, 1)" :step="0.1" style="width:100%" @change="onPositionBindingChange(axis.key, 'scale', $event || 1)" />
                </a-form-item>
                <a-form-item v-if="getPositionBinding(axis.key).transform === 'offset'" :key="axis.key + '_offset'" :label="$t('ISM3DEditor.offsetTransform')">
                  <a-input-number :value="n(getPositionBinding(axis.key).offset, 0)" :step="1" style="width:100%" @change="onPositionBindingChange(axis.key, 'offset', $event || 0)" />
                </a-form-item>
              </template>
            </a-form>
          </div>

          <div class="tab-section data-bind-section" v-if="showRealtimeDataDisplay">
            <div class="tab-section-title">{{ $t('ISM3DEditor.realTimeDataDisplay') }}</div>
            <a-form layout="horizontal" :label-col="{ span: 5 }" :wrapper-col="{ span: 18 }" size="small" class="data-bind-form">
              <a-form-item :label="$t('ISM3DEditor.dataName')">
                <a-input-group compact class="data-select-group">
                  <a-input :value="currentObject.dataName || ''" disabled />
                  <a-button type="primary" @click="openDataSelector('dataText')">
                    <i class="fas fa-database"></i> {{ currentObject.dataName ? '重新选择' : '选择数据' }}
                  </a-button>
                </a-input-group>
              </a-form-item>
              <a-form-item label="显示模板">
                <a-input :value="currentObject.dataFormat || '{value}'" @change="onPropChange('dataFormat', $event)" placeholder="如：温度: {value}" />
              </a-form-item>
              <a-form-item label="当前值">
                <a-input :value="currentObject.realTimeValue || ''" disabled />
              </a-form-item>
            </a-form>
          </div>
        </a-tab-pane>
        <a-tab-pane v-if="false" key="eventsLegacy" tab="事件" :forceRender="true">
          <div class="tab-section event-section" v-for="eventItem in interactionEventOptions" :key="eventItem.key">
            <div class="event-section-head">
              <span class="tab-section-title">{{ eventItem.label }}</span>
              <a-switch size="small" :checked="!!getInteraction(eventItem.key).enabled" @change="onInteractionChange(eventItem.key, 'enabled', $event)" />
            </div>
            <a-form v-if="!!getInteraction(eventItem.key).enabled" layout="horizontal" :label-col="{ span: 6 }" :wrapper-col="{ span: 17 }" size="small">
              <a-form-item label="动作">
                <a-select :value="getInteraction(eventItem.key).type || 'none'" @change="onInteractionTypeChange(eventItem.key, $event)">
                  <a-select-option value="none">无</a-select-option>
                  <a-select-option value="link">外部链接</a-select-option>
                  <a-select-option value="route">内部页面</a-select-option>
                  <a-select-option value="deviceView">设备视图</a-select-option>
                  <a-select-option value="visible">显隐组件</a-select-option>
                  <a-select-option value="popupText">弹窗文本</a-select-option>
                </a-select>
              </a-form-item>

              <template v-if="getInteraction(eventItem.key).type === 'link'">
                <a-form-item label="链接">
                  <a-input :value="getInteraction(eventItem.key).payload || ''" placeholder="https://..." @change="onInteractionInput(eventItem.key, 'payload', $event)" />
                </a-form-item>
                <a-form-item label="打开方式">
                  <a-select :value="getInteraction(eventItem.key).target || '_blank'" @change="onInteractionChange(eventItem.key, 'target', $event)">
                    <a-select-option value="_blank">新窗口</a-select-option>
                    <a-select-option value="_self">当前窗口</a-select-option>
                  </a-select>
                </a-form-item>
              </template>

              <template v-if="getInteraction(eventItem.key).type === 'route'">
                <a-form-item label="页面路径">
                  <a-input :value="getInteraction(eventItem.key).routePath || ''" placeholder="/AppRun/..." @change="onInteractionInput(eventItem.key, 'routePath', $event)" />
                </a-form-item>
                <a-form-item label="显示方式">
                  <a-select :value="getInteraction(eventItem.key).displayMode || 'window'" @change="onInteractionChange(eventItem.key, 'displayMode', $event)">
                    <a-select-option value="window">新窗口</a-select-option>
                    <a-select-option value="popup">弹窗</a-select-option>
                  </a-select>
                </a-form-item>
              </template>

              <template v-if="getInteraction(eventItem.key).type === 'deviceView'">
                <a-form-item label="设备 Key">
                  <a-input :value="getInteraction(eventItem.key).deviceKey || ''" @change="onInteractionInput(eventItem.key, 'deviceKey', $event)" />
                </a-form-item>
                <a-form-item label="页面 UUID">
                  <a-input :value="getInteraction(eventItem.key).showPageUUID || ''" @change="onInteractionInput(eventItem.key, 'showPageUUID', $event)" />
                </a-form-item>
                <a-form-item label="弹窗">
                  <a-switch size="small" :checked="!!getInteraction(eventItem.key).isPopUp" @change="onInteractionChange(eventItem.key, 'isPopUp', $event)" />
                </a-form-item>
              </template>

              <template v-if="getInteraction(eventItem.key).type === 'visible'">
                <a-form-item label="显示">
                  <a-select mode="multiple" :value="getInteraction(eventItem.key).showTargetsList || []" @change="onInteractionTargetsChange(eventItem.key, 'showTargetsList', 'showTargets', $event)">
                    <a-select-option v-for="obj in sceneObjects" :key="'show_' + obj.id" :value="obj.id">{{ obj.name || obj.typeName || obj.id }}</a-select-option>
                  </a-select>
                </a-form-item>
                <a-form-item label="隐藏">
                  <a-select mode="multiple" :value="getInteraction(eventItem.key).hideTargetsList || []" @change="onInteractionTargetsChange(eventItem.key, 'hideTargetsList', 'hideTargets', $event)">
                    <a-select-option v-for="obj in sceneObjects" :key="'hide_' + obj.id" :value="obj.id">{{ obj.name || obj.typeName || obj.id }}</a-select-option>
                  </a-select>
                </a-form-item>
              </template>

              <template v-if="getInteraction(eventItem.key).type === 'popupText'">
                <a-form-item label="文本">
                  <a-textarea :value="getInteraction(eventItem.key).payload || ''" :auto-size="{ minRows: 2, maxRows: 5 }" @change="onInteractionInput(eventItem.key, 'payload', $event)" />
                </a-form-item>
              </template>

              <template v-if="getInteraction(eventItem.key).type === 'route' || getInteraction(eventItem.key).type === 'link'">
                <a-form-item label="标题">
                  <a-input :value="getInteraction(eventItem.key).title || ''" @change="onInteractionInput(eventItem.key, 'title', $event)" />
                </a-form-item>
                <a-form-item label="宽高">
                  <a-input-group compact>
                    <a-input-number :value="n(getInteraction(eventItem.key).width, 960)" :min="320" style="width:50%" @change="onInteractionChange(eventItem.key, 'width', $event)" />
                    <a-input-number :value="n(getInteraction(eventItem.key).height, 640)" :min="240" style="width:50%" @change="onInteractionChange(eventItem.key, 'height', $event)" />
                  </a-input-group>
                </a-form-item>
              </template>

              <a-form-item label="二次确认">
                <a-switch size="small" :checked="!!getInteraction(eventItem.key).actionConfirm" @change="onInteractionChange(eventItem.key, 'actionConfirm', $event)" />
              </a-form-item>
              <a-form-item label="语音提示">
                <a-input :value="getInteraction(eventItem.key).actionVoice || ''" @change="onInteractionInput(eventItem.key, 'actionVoice', $event)" />
              </a-form-item>
            </a-form>
          </div>
          <div class="tab-section data-bind-section data-bind-empty" v-if="!hasDataBindingPanels">
            当前组件暂无可配置的数据绑定
          </div>
        </a-tab-pane>
      </a-tabs>

    </template>

    <device-data-model ref="deviceDataModel" @onSelectDataModel="onSelectData"></device-data-model>
    <system-image-model ref="systemImageModel" :networkImageUrl="''" @onSelectImage="onTextureSelected"></system-image-model>
  </div>
</template>

<script>
import deviceDataModel from "@/components/deviceDataModel/deviceDataModel.vue";
import systemImageModel from "@/components/systemImageModel/systemImageModel.vue";
import { displayModelList, getDisplayModelLayerData } from "@/services/displayModel";
import {
  createGLTFAnimationCondition,
  ensureGLTFAnimationGroups,
  normalizeGLTFAnimationGroup,
  syncLegacyGLTFAnimationFields
} from '../../utils/GLTFAnimationGroups'
import { createDefaultPositionBinding, ensurePositionBindings, POSITION_BINDING_AXES } from '../../utils/positionBindings'

export default {
  name: 'PropertiesPanel',
  components: {
    deviceDataModel,
    systemImageModel
  },
  props: {
    currentObject: { type: Object, default: null },
    sceneObjects: { type: Array, default: () => [] }
  },
  data() {
    return {
      activeKey: 'style',
      dataSelectMode: 'animate',
      editingGLTFAnimationGroupId: '',
      editingPositionBindingAxis: '',
      editingActionIndex: -1,
      editingSetValueIndex: -1,
      configurationModel: [],
      displayPageList: new Map()
    }
  },
  i18n: require('@/i18n/language'),
  mounted() {
    this.getConfigurationModel()
  },
  computed: {
    gltfAnimationOptions() {
      if (!this.currentObject || !Array.isArray(this.currentObject.gltfAnimations)) return []
      const counts = {}
      this.currentObject.gltfAnimations.forEach((item, index) => {
        const name = typeof item === 'string' ? item : (item.name || item.label || ('Animation ' + (index + 1)))
        counts[name] = (counts[name] || 0) + 1
      })
      return this.currentObject.gltfAnimations.map((item, index) => {
        if (typeof item === 'string') {
          return {
            key: item,
            name: item,
            label: counts[item] > 1 ? item + ' #' + (index + 1) : item
          }
        }
        const name = item.name || item.label || item.key || ('Animation ' + (index + 1))
        return {
          key: item.key || name,
          name,
          label: item.label || (counts[name] > 1 ? name + ' #' + (index + 1) : name)
        }
      })
    },
    firstGLTFAnimation() {
      return this.gltfAnimationOptions.length ? this.gltfAnimationOptions[0].key : ''
    },
    gltfAnimationGroups() {
      if (!this.currentObject) return []
      return Array.isArray(this.currentObject.gltfAnimationGroups) ? this.currentObject.gltfAnimationGroups : []
    },
    selectedGLTFAnimationNames() {
      if (!this.currentObject) return []
      if (Array.isArray(this.currentObject.gltfAnimationNames)) return this.currentObject.gltfAnimationNames
      if (this.currentObject.gltfAnimationName) return [this.currentObject.gltfAnimationName]
      return this.firstGLTFAnimation ? [this.firstGLTFAnimation] : []
    },
    gltfAnimationCondition() {
      return this.currentObject && this.currentObject.gltfAnimationCondition ? this.currentObject.gltfAnimationCondition : {}
    },
    actionEventOptions() {
      return [
        { label: 'displayConfig.Properties.action.click', value: 'click' },
        { label: 'displayConfig.Properties.action.MouseDown', value: 'mousedown' },
        { label: 'displayConfig.Properties.action.MouseUp', value: 'mouseup' },
        { label: 'displayConfig.Properties.action.dbClick', value: 'dblclick' },
        { label: 'displayConfig.Properties.action.mouseenter', value: 'mouseenter' },
        { label: 'displayConfig.Properties.action.mouseleave', value: 'mouseleave' }
      ]
    },
    actionTypeOptions() {
      return [
        { label: 'displayConfig.Properties.action.openLink', value: 'link' },
        { label: 'displayConfig.Properties.action.SetValue', value: 'SetValue' },
        { label: 'displayConfig.Properties.action.Visible', value: 'visible' },
        { label: 'displayConfig.Properties.action.SysScript', value: 'SysScript' },
        { label: 'displayConfig.Properties.action.DeviceView', value: 'DeviceView' },
        { label: 'displayConfig.Properties.action.RestApi', value: 'RestApi' },
        { label: 'displayConfig.Properties.action.animation', value: 'Animation' }
      ]
    },
    roleOptions() {
      return [
        { RoleId: 'Operator', Name: this.$t('account.settings.UserList.RoleOperator') },
        { RoleId: 'User', Name: this.$t('account.settings.UserList.RoleUser') }
      ]
    },
    bindPropOptions() {
      if (!this.currentObject) return []
      const options = [
        { label: this.$t('ISM3DEditor.visible'), value: 'visible' }
      ]
      if (this.isTextLabelObject(this.currentObject)) {
        options.push({ label: this.$t('ISM3DEditor.color'), value: 'color' })
        options.push({ label: this.$t('ISM3DEditor.labelContent'), value: 'textContent' })
        return options
      }
      if (this.currentObject.type === 'gltf' || this.currentObject.isExternalModel) {
        return options
      }
      options.push({ label: this.$t('ISM3DEditor.color'), value: 'color' })
      options.push({ label: this.$t('ISM3DEditor.opacity'), value: 'opacity' })
      if (!this.isBasicMediaObject(this.currentObject)) {
        options.push({ label: this.$t('ISM3DEditor.scale'), value: 'scale' })
      }
      return options
    },
    showGenericDataBinding() {
      return !!(this.currentObject && this.bindPropOptions.length > 0 && this.currentObject.type !== 'dataText' && this.currentObject.type !== '2dComponent')
    },
    showPositionDataBinding() {
      if (!this.currentObject) return false
      return this.currentObject.type !== '2dComponent' && !this.isBasicMediaObject(this.currentObject)
    },
    showRealtimeDataDisplay() {
      return !!(this.currentObject && this.currentObject.type === 'dataText')
    },
    hasDataBindingPanels() {
      return this.showGenericDataBinding || this.showPositionDataBinding || this.showRealtimeDataDisplay
    },
    showBindTransformControls() {
      return this.currentObject && (this.currentObject.bindProp === 'opacity' || this.currentObject.bindProp === 'scale')
    },
    positionBindingAxes() {
      return POSITION_BINDING_AXES.map((axis) => ({
        key: axis,
        label: this.$t('ISM3DEditor.position' + axis.toUpperCase())
      }))
    },
    interactionEventOptions() {
      return [
        { key: 'click', label: '单击' },
        { key: 'dblclick', label: '双击' },
        { key: 'mouseenter', label: '鼠标移入' },
        { key: 'mouseleave', label: '鼠标移出' }
      ]
    }
  },
  watch: {
    currentObject: {
      handler(newVal) {
        if (!newVal) return
        if (!newVal.animate) {
          this.$set(newVal, 'animate', {
            selected: [],
            condition: {
              isBandDevice: false,
              deviceSN: '',
              DeviceName: '',
              dataID: '',
              dataName: '',
              operator: '',
              OperatorValue: '',
              OperatorMaxValue: ''
            },
            isExpression: false,
            expression: '',
            animateList: [],
            animateElement: []
          })
        }
        if (newVal.bindProp === undefined) {
          this.$set(newVal, 'bindProp', '')
          this.$set(newVal, 'bindTransform', 'direct')
          this.$set(newVal, 'bindScale', 1)
          this.$set(newVal, 'bindOffset', 0)
          this.$set(newVal, 'wsKey', '')
        }
        if (!this.showGenericDataBinding || (newVal.bindProp && !this.bindPropOptions.some(option => option.value === newVal.bindProp))) {
          this.$set(newVal, 'bindProp', '')
          this.$set(newVal, 'bindTransform', 'direct')
          this.$set(newVal, 'bindScale', 1)
          this.$set(newVal, 'bindOffset', 0)
        } else if (newVal.bindProp !== 'opacity' && newVal.bindProp !== 'scale' && newVal.bindTransform !== 'direct') {
          this.$set(newVal, 'bindTransform', 'direct')
          this.$set(newVal, 'bindScale', 1)
          this.$set(newVal, 'bindOffset', 0)
        }
        ensurePositionBindings(newVal)
        if (newVal.realTimeValue === undefined) {
          this.$set(newVal, 'realTimeValue', '')
        }
        if (newVal.dataFormat === undefined) {
          this.$set(newVal, 'dataFormat', '{value}')
        }
        if (newVal.dataID === undefined) {
          this.$set(newVal, 'dataID', '')
        }
        if (newVal.dataName === undefined) {
          this.$set(newVal, 'dataName', '')
        }
        if (newVal.deviceSN === undefined) {
          this.$set(newVal, 'deviceSN', '')
        }
        if (newVal.DeviceName === undefined) {
          this.$set(newVal, 'DeviceName', '')
        }
        if (newVal.isBandDevice === undefined) {
          this.$set(newVal, 'isBandDevice', false)
        }
        this.ensureActionList(newVal)
        if (newVal.type === 'gltf' || newVal.isExternalModel) {
          if (!Array.isArray(newVal.gltfAnimations)) this.$set(newVal, 'gltfAnimations', [])
          if (newVal.gltfAnimationName === undefined) this.$set(newVal, 'gltfAnimationName', '')
          if (!Array.isArray(newVal.gltfAnimationNames)) {
            const names = newVal.gltfAnimationName ? [newVal.gltfAnimationName] : []
            this.$set(newVal, 'gltfAnimationNames', names)
          }
          if (newVal.gltfAnimationPlaying === undefined) this.$set(newVal, 'gltfAnimationPlaying', false)
          if (newVal.gltfAnimationSpeed === undefined) this.$set(newVal, 'gltfAnimationSpeed', 1)
          if (newVal.gltfAnimationLoop === undefined) this.$set(newVal, 'gltfAnimationLoop', true)
          if (newVal.gltfAnimationConditionEnabled === undefined) this.$set(newVal, 'gltfAnimationConditionEnabled', false)
          if (!newVal.gltfAnimationCondition) {
            this.$set(newVal, 'gltfAnimationCondition', {
              isBandDevice: false,
              deviceSN: '',
              DeviceName: '',
              dataID: '',
              dataName: '',
              operator: '',
              OperatorValue: '',
              OperatorMaxValue: ''
            })
          }
          const groups = ensureGLTFAnimationGroups(newVal, { defaultAnimationKey: this.firstGLTFAnimation })
          this.$set(newVal, 'gltfAnimationGroups', groups)
          syncLegacyGLTFAnimationFields(newVal, groups)
        }
      },
      immediate: true
    }
  },
  methods: {
    n(val, fallback) {
      const num = parseFloat(val)
      return isNaN(num) ? fallback : num
    },
    isTextLabelObject(obj) {
      return obj && (obj.type === 'text3d' || obj.type === 'dataText' || obj.type === 'textPlain3d' || obj.type === 'uiLabel')
    },
    isBasicMediaObject(obj) {
      return obj && (obj.type === 'image3d' || obj.type === 'uiImage' || obj.type === 'video3d' || obj.type === 'webEmbed')
    },
    createDefaultLink() {
      return {
        linkType: 'External',
        Inside: {
          displayUUID: '',
          pageUUID: ''
        },
        External: '',
        ExternalSource: 'Url',
        ExternalPage: {
          displayUUID: '',
          pageUUID: ''
        },
        OpenExternalType: 'new',
        isPopUp: false,
        autoClose: false,
        width: 960,
        height: 640,
        title: ''
      }
    },
    createDefaultSetValue() {
      return {
        IsManual: false,
        isBandDevice: false,
        deviceSN: '',
        DeviceName: '',
        dataID: '',
        dataName: '',
        AutoSetValue: '',
        SetPassword: ''
      }
    },
    createDefaultAction() {
      return {
        type: 'click',
        action: 'visible',
        actionAuth: [],
        ActionPassword: '',
        actionVoice: '',
        actionConfirm: false,
        showItems: [],
        hideItems: [],
        link: this.createDefaultLink(),
        setValue: [],
        SetDelay: 1000,
        RestApi: {
          Name: '',
          IsSystem: '1',
          Type: 'Post',
          Url: '',
          Params: '{}'
        },
        DeviceView: {
          key: '',
          showUUID: '',
          showPageUUID: '',
          type: '',
          isPopUp: false,
          isContainer: false,
          routePath: ''
        },
        ScriptList: [],
        animationStatus: 'start'
      }
    },
    normalizeActionItem(action) {
      const next = Object.assign(this.createDefaultAction(), action || {})
      next.actionAuth = Array.isArray(next.actionAuth) ? next.actionAuth : []
      next.showItems = Array.isArray(next.showItems) ? next.showItems : []
      next.hideItems = Array.isArray(next.hideItems) ? next.hideItems : []
      next.link = Object.assign(this.createDefaultLink(), next.link || {})
      next.link.Inside = Object.assign({ displayUUID: '', pageUUID: '' }, next.link.Inside || {})
      next.link.ExternalPage = Object.assign({ displayUUID: '', pageUUID: '' }, next.link.ExternalPage || {})
      next.setValue = Array.isArray(next.setValue) ? next.setValue.map(item => Object.assign(this.createDefaultSetValue(), item || {})) : []
      next.RestApi = Object.assign(this.createDefaultAction().RestApi, next.RestApi || {})
      next.DeviceView = Object.assign(this.createDefaultAction().DeviceView, next.DeviceView || {})
      next.ScriptList = Array.isArray(next.ScriptList) ? next.ScriptList : []
      return next
    },
    ensureActionList(obj) {
      if (!obj) return []
      const source = Array.isArray(obj.action) ? obj.action : []
      const normalized = source.map(item => this.normalizeActionItem(item))
      this.$set(obj, 'action', normalized)
      return obj.action
    },
    emitActionChange() {
      if (this.currentObject) this.ensureActionList(this.currentObject)
      this.$emit('prop-change')
    },
    addBindAction() {
      if (!this.currentObject) return
      const list = this.ensureActionList(this.currentObject)
      list.push(this.createDefaultAction())
      this.emitActionChange()
    },
    delBindAction(index) {
      if (!this.currentObject || !Array.isArray(this.currentObject.action)) return
      this.currentObject.action.splice(index, 1)
      this.emitActionChange()
    },
    changeAction(value, index) {
      if (!this.currentObject || !Array.isArray(this.currentObject.action)) return
      const event = this.currentObject.action[index]
      if (!event) return
      this.$set(event, 'action', value)
      const normalized = this.normalizeActionItem(event)
      this.$set(this.currentObject.action, index, normalized)
      this.emitActionChange()
    },
    addBindSetValue(index) {
      if (!this.currentObject || !Array.isArray(this.currentObject.action)) return
      const event = this.currentObject.action[index]
      if (!event) return
      if (!Array.isArray(event.setValue)) this.$set(event, 'setValue', [])
      event.setValue.push(this.createDefaultSetValue())
      this.emitActionChange()
    },
    delBindSetValue(index, setValueIndex) {
      if (!this.currentObject || !Array.isArray(this.currentObject.action)) return
      const event = this.currentObject.action[index]
      if (!event || !Array.isArray(event.setValue)) return
      event.setValue.splice(setValueIndex, 1)
      this.emitActionChange()
    },
    openActionSetValueSelector(index, setValueIndex) {
      this.editingActionIndex = index
      this.editingSetValueIndex = setValueIndex
      this.openDataSelector('actionSetValue')
    },
    getConfigurationModel() {
      const vm = this
      this.configurationModel = []
      ;[1, 2].forEach(function(displayType) {
        displayModelList({ DisplayType: displayType }).then(function(res) {
          const list = res && res.data ? res.data.list : null
          if (!Array.isArray(list)) return
          list.forEach(function(item) {
            const tableData = {
              name: item.name,
              description: item.description,
              uuid: item.displayUid,
              displayType: displayType
            }
            vm.configurationModel.push(tableData)
            vm.getDisplayPage(item.displayUid)
          })
        })
      })
    },
    getDisplayPage(uuid) {
      if (!uuid) return
      const vm = this
      getDisplayModelLayerData({ muid: uuid }).then(function(res) {
        const data = res && res.data ? res.data : {}
        const pageLayer = Array.isArray(data.layer) ? data.layer : (data.pageLayer || [])
        const displayArray = []
        if (Array.isArray(pageLayer)) {
          pageLayer.forEach(function(page) {
            displayArray.push({
              label: page.PageName || page.title || page.name || '',
              value: page.PageId || page.id || page.uuid || '',
              pageType: page.PageType,
              pageModelUuid: page.modelId
            })
          })
        }
        vm.displayPageList.set(uuid, displayArray)
        vm.displayPageList = new Map(vm.displayPageList)
      })
    },
    generateTargetPage(uuid) {
      return this.displayPageList.get(uuid) || []
    },
    get2DConfigurationModel() {
      return this.configurationModel.filter(item => item.displayType !== 2)
    },
    ensureLinkExternalPage(link) {
      if (!link.ExternalPage) {
        this.$set(link, 'ExternalPage', { displayUUID: '', pageUUID: '' })
      }
      return link.ExternalPage
    },
    buildExternalPageUrl(displayUUID, pageUUID) {
      if (!displayUUID || !pageUUID) return ''
      const basePath = window.location.origin + window.location.pathname
      return basePath + '#/AppRun/' + displayUUID + '?pageId=' + pageUUID
    },
    changeLinkExternalSource(link, value) {
      this.$set(link, 'ExternalSource', value)
      this.ensureLinkExternalPage(link)
      if (value === 'Page') {
        link.External = this.buildExternalPageUrl(link.ExternalPage.displayUUID, link.ExternalPage.pageUUID)
      }
      this.emitActionChange()
    },
    changeLinkExternalDisplay(link) {
      const externalPage = this.ensureLinkExternalPage(link)
      externalPage.pageUUID = ''
      link.External = ''
      this.emitActionChange()
    },
    changeLinkExternalPage(link) {
      const externalPage = this.ensureLinkExternalPage(link)
      link.External = this.buildExternalPageUrl(externalPage.displayUUID, externalPage.pageUUID)
      this.emitActionChange()
    },
    getLinkInsideDisplayType(uuid) {
      const target = this.configurationModel.find(item => item.uuid === uuid)
      return target ? target.displayType : 1
    },
    changeLinkInsideDisplay(link) {
      if (!link || !link.Inside) {
        this.emitActionChange()
        return
      }
      const displayType = this.getLinkInsideDisplayType(link.Inside.displayUUID)
      this.$set(link.Inside, 'displayType', displayType)
      this.$set(link.Inside, 'pageUUID', '')
      this.emitActionChange()
    },
    createDefaultInteraction() {
      return {
        enabled: false,
        type: 'none',
        payload: '',
        target: '_blank',
        width: 960,
        height: 640,
        title: '',
        routePath: '',
        displayMode: 'window',
        autoClose: false,
        showTargets: '',
        hideTargets: '',
        showTargetsList: [],
        hideTargetsList: [],
        deviceKey: '',
        showUUID: '',
        showPageUUID: '',
        isPopUp: false,
        isContainer: false,
        deviceType: '',
        actionConfirm: false,
        actionVoice: '',
        actionAuth: []
      }
    },
    ensureInteractions(obj) {
      if (!obj) return {}
      if (!obj.interactions || Array.isArray(obj.interactions)) {
        this.$set(obj, 'interactions', {})
      }
      this.interactionEventOptions.forEach((eventItem) => {
        if (!obj.interactions[eventItem.key]) {
          this.$set(obj.interactions, eventItem.key, this.createDefaultInteraction())
        } else {
          obj.interactions[eventItem.key] = Object.assign(this.createDefaultInteraction(), obj.interactions[eventItem.key])
        }
      })
      return obj.interactions
    },
    getInteraction(eventKey) {
      if (!this.currentObject) return this.createDefaultInteraction()
      const interactions = this.currentObject.interactions || {}
      return interactions[eventKey] || this.createDefaultInteraction()
    },
    emitInteractionChange() {
      this.$emit('prop-change')
    },
    onInteractionChange(eventKey, prop, value) {
      const interactions = this.ensureInteractions(this.currentObject)
      const interaction = interactions[eventKey]
      this.$set(interaction, prop, value)
      if (prop === 'enabled' && value && (!interaction.type || interaction.type === 'none')) {
        this.$set(interaction, 'type', 'popupText')
      }
      if (prop === 'enabled' && !value) {
        this.$set(interaction, 'type', 'none')
      }
      this.emitInteractionChange()
    },
    onInteractionInput(eventKey, prop, e) {
      this.onInteractionChange(eventKey, prop, e && e.target ? e.target.value : e)
    },
    onInteractionTypeChange(eventKey, value) {
      const interactions = this.ensureInteractions(this.currentObject)
      const interaction = interactions[eventKey]
      this.$set(interaction, 'type', value || 'none')
      this.$set(interaction, 'enabled', !!value && value !== 'none')
      this.emitInteractionChange()
    },
    onInteractionTargetsChange(eventKey, listProp, textProp, values) {
      const list = Array.isArray(values) ? values : []
      const interactions = this.ensureInteractions(this.currentObject)
      const interaction = interactions[eventKey]
      this.$set(interaction, listProp, list)
      this.$set(interaction, textProp, list.join(','))
      this.emitInteractionChange()
    },
    normColor(val) {
      if (typeof val === 'number') return '#' + val.toString(16).padStart(6, '0')
      if (!val || typeof val !== 'string') return '#000000'
      const s = val.trim()
      const map = { blue:'#4da6ff',green:'#4dffa6',orange:'#ffaa4d',purple:'#c17bff',pink:'#ff7bb5',cyan:'#4dfff0',yellow:'#fffb4d',red:'#ff4d4f',gray:'#8c8c8c',brown:'#a67c52' }
      if (map[s]) return map[s]
      if (s[0] === '#') {
        if (/^#[0-9a-fA-F]{3}$/.test(s)) return '#' + s[1] + s[1] + s[2] + s[2] + s[3] + s[3]
        if (/^#[0-9a-fA-F]{6}$/.test(s)) return s.toLowerCase()
        if (/^#[0-9a-fA-F]{8}$/.test(s)) return s.slice(0, 7).toLowerCase()
      }
      return '#000000'
    },
    normBgColor(val) {
      if (typeof val === 'number') return '#' + val.toString(16).padStart(8, '0')
      if (!val || typeof val !== 'string') return '#000000'
      var s = val.trim()
      if (s === '#00000000' || s === 'transparent') return '#000000'
      if (s[0] === '#') {
        if (/^#[0-9a-fA-F]{3}$/.test(s)) return '#' + s[1] + s[1] + s[2] + s[2] + s[3] + s[3]
        if (/^#[0-9a-fA-F]{6}$/.test(s)) return s.toLowerCase()
        if (/^#[0-9a-fA-F]{8}$/.test(s)) return s.slice(0, 7).toLowerCase()
      }
      return '#000000'
    },
    onBgColorInput(e) {
      var hex = e.target.value
      var alpha = this.currentObject.textBgColor && this.currentObject.textBgColor.length === 9
        ? this.currentObject.textBgColor.slice(7, 9) : 'ff'
      if (alpha === '00') alpha = 'ff'
      this.setProp('textBgColor', hex + alpha)
      if (!this.n(this.currentObject.textBgOpacity, 0)) {
        this.setProp('textBgOpacity', 0.48)
      }
      this.$emit('prop-change')
    },
    onBgColorTextInput(e) {
      const value = e && e.target ? e.target.value : e
      this.setProp('textBgColor', value)
      if (value && value !== '#00000000' && value !== 'transparent' && !this.n(this.currentObject.textBgOpacity, 0)) {
        this.setProp('textBgOpacity', 0.48)
      }
      this.$emit('prop-change')
    },
    clearBgColor() {
      this.setProp('textBgColor', '#00000000')
      this.setProp('textBgOpacity', 0)
      this.$emit('prop-change')
    },
    clearColorValue(prop, fallback) {
      this.setProp(prop, fallback || '#ffffff')
      this.$emit('material-change')
    },
    setProp(prop, val) {
      this.$set(this.currentObject, prop, val)
      // 如果修改的是 GLTF 模型的材质属性，标记用户已主动修改，重载时保留自定义材质
      var materialProps = ['color', 'opacity', 'metalness', 'roughness', 'wireframe', 'emissive', 'textureData']
      if (materialProps.indexOf(prop) !== -1 && this.currentObject) {
        this.$set(this.currentObject, 'materialOverridden', true)
      }
      this.$emit('prop-change')
    },
    // ===== 贴图功能（使用项目图片组件） =====
    openTexturePicker() {
      if (this.$refs.systemImageModel) {
        this.$refs.systemImageModel.showModal(0)
      }
    },
    onTextureSelected(url) {
      if (!url) return
      if (this.currentObject && (this.currentObject.type === 'image3d' || this.currentObject.type === 'uiImage')) {
        this.$set(this.currentObject, 'mediaUrl', url)
        this.$set(this.currentObject, 'imageUrl', url)
        this.$set(this.currentObject, 'textureData', url)
        this.$set(this.currentObject, 'materialOverridden', true)
        this.$emit('prop-change')
        return
      }
      // Text sprites do not support user texture maps.
      if (this.currentObject && (this.currentObject.type === 'text3d' || this.currentObject.type === 'dataText')) {
        if (this.$message) this.$message.warning('文字对象不支持贴图，请选中模型对象后再设置')
        return
      }
      this.$set(this.currentObject, 'textureData', url)
      this.$set(this.currentObject, 'materialOverridden', true)
      this.$emit('prop-change')
    },
    clearMediaImage() {
      if (!this.currentObject) return
      this.$set(this.currentObject, 'mediaUrl', '')
      this.$set(this.currentObject, 'imageUrl', '')
      this.$set(this.currentObject, 'textureData', '')
      this.$set(this.currentObject, 'materialOverridden', true)
      this.$emit('prop-change')
    },
    onMediaUrlChange(e) {
      var url = e && e.target ? e.target.value : e
      this.setProp('mediaUrl', url)
      if (this.currentObject.type === 'image3d' || this.currentObject.type === 'uiImage') {
        this.setProp('imageUrl', url)
        this.setProp('textureData', url)
      } else if (this.currentObject.type === 'video3d') {
        this.setProp('videoUrl', url)
      }
      this.$emit('prop-change')
    },
    onWebUrlChange(e) {
      var url = e && e.target ? e.target.value : e
      this.setProp('webUrl', url)
      this.setProp('mediaUrl', url)
      this.$emit('prop-change')
    },
    clearTexture() {
      this.$set(this.currentObject, 'textureData', '')
      this.$set(this.currentObject, 'materialOverridden', true)
      this.$emit('prop-change')
    },
    onPropChange(prop, e) {
      this.setProp(prop, e.target.value)
    },
    onPropInput(prop, e) {
      this.setProp(prop, e.target ? e.target.value : e)
    },
    onPropSwitch(prop, checked) {
      this.setProp(prop, checked)
      if (prop === 'wireframe' || prop === 'showShadow' || prop === 'color' || prop === 'emissive') {
        this.$emit('material-change')
      }
    },
    onValChange(prop, val) {
      this.setProp(prop, val)
    },
    onColorInput(prop, e) {
      var val = e.target.value
      console.log('[PropertiesPanel] onColorInput prop=' + prop + ' val=' + val + ' obj.type=' + (this.currentObject && this.currentObject.type) + ' obj.id=' + (this.currentObject && this.currentObject.id))
      this.setProp(prop, val)
      this.$emit('material-change')
    },
    onAnimVal(prop, val) {
      this.setProp(prop, val)
      this.$emit('anim-change')
    },
    onBindPropChange(value) {
      this.$set(this.currentObject, 'bindProp', value)
      if (!value || (value !== 'opacity' && value !== 'scale')) {
        this.$set(this.currentObject, 'bindTransform', 'direct')
        this.$set(this.currentObject, 'bindScale', 1)
        this.$set(this.currentObject, 'bindOffset', 0)
      }
      this.$emit('prop-change')
    },
    onBindTransformChange(value) {
      this.$set(this.currentObject, 'bindTransform', value)
      this.$emit('prop-change')
    },
    onBindScaleChange(value) {
      this.$set(this.currentObject, 'bindScale', value || 1)
      this.$emit('prop-change')
    },
    onBindOffsetChange(value) {
      this.$set(this.currentObject, 'bindOffset', value || 0)
      this.$emit('prop-change')
    },
    getPositionBinding(axis) {
      if (!this.currentObject) return {}
      ensurePositionBindings(this.currentObject)
      if (!this.currentObject.positionBindings[axis]) {
        this.$set(this.currentObject.positionBindings, axis, createDefaultPositionBinding())
      }
      return this.currentObject.positionBindings[axis]
    },
    onPositionBindingChange(axis, prop, value) {
      const binding = this.getPositionBinding(axis)
      this.$set(binding, prop, value)
      this.$emit('prop-change')
    },
    openPositionBindingSelector(axis) {
      this.editingPositionBindingAxis = axis
      this.openDataSelector('positionBinding')
    },
    onAnimateBindChange(prop, value, isAnimate = false) {
      if (!this.currentObject.animate) return
      if (isAnimate) {
        this.$set(this.currentObject.animate, prop, value)
      } else {
        this.$set(this.currentObject.animate.condition, prop, value)
      }
      this.$emit('prop-change')
    },
    onAnimateSelectedChange(checkedValues) {
      if (!this.currentObject.animate) return
      this.$set(this.currentObject.animate, 'selected', checkedValues)
      this.$emit('prop-change')
    },
    onGLTFAnimationNamesChange(values) {
      const names = Array.isArray(values) ? values : []
      this.$set(this.currentObject, 'gltfAnimationNames', names)
      this.$set(this.currentObject, 'gltfAnimationName', names[0] || '')
      this.$emit('prop-change')
    },
    onGLTFAnimationConditionChange(prop, value) {
      if (!this.currentObject) return
      if (!this.currentObject.gltfAnimationCondition) {
        this.$set(this.currentObject, 'gltfAnimationCondition', {})
      }
      this.$set(this.currentObject.gltfAnimationCondition, prop, value)
      this.$emit('prop-change')
    },
    commitGLTFAnimationGroups(groups) {
      const nextGroups = ensureGLTFAnimationGroups({
        ...this.currentObject,
        gltfAnimationGroups: groups
      }, { defaultAnimationKey: this.firstGLTFAnimation })
      this.$set(this.currentObject, 'gltfAnimationGroups', nextGroups)
      syncLegacyGLTFAnimationFields(this.currentObject, nextGroups)
      this.$emit('prop-change')
    },
    addGLTFAnimationGroup() {
      if (!this.currentObject) return
      const groups = ensureGLTFAnimationGroups(this.currentObject, { defaultAnimationKey: this.firstGLTFAnimation })
      const group = normalizeGLTFAnimationGroup({
        id: 'gltf_anim_group_' + Date.now() + '_' + groups.length,
        name: '动画组' + (groups.length + 1),
        playing: false,
        animationNames: this.firstGLTFAnimation ? [this.firstGLTFAnimation] : [],
        speed: 1,
        loop: true,
        conditionEnabled: false,
        condition: createGLTFAnimationCondition()
      }, groups.length)
      groups.push(group)
      this.commitGLTFAnimationGroups(groups)
    },
    removeGLTFAnimationGroup(groupId) {
      if (!this.currentObject) return
      const groups = ensureGLTFAnimationGroups(this.currentObject, { defaultAnimationKey: this.firstGLTFAnimation })
        .filter(group => group.id !== groupId)
      this.commitGLTFAnimationGroups(groups.length ? groups : [normalizeGLTFAnimationGroup({}, 0)])
    },
    setGLTFAnimationGroupProp(group, prop, value) {
      if (!this.currentObject || !group) return
      const groups = ensureGLTFAnimationGroups(this.currentObject, { defaultAnimationKey: this.firstGLTFAnimation })
      const target = groups.find(item => item.id === group.id)
      if (!target) return
      target[prop] = value
      this.commitGLTFAnimationGroups(groups)
    },
    setGLTFAnimationGroupCondition(group, prop, value) {
      if (!this.currentObject || !group) return
      const groups = ensureGLTFAnimationGroups(this.currentObject, { defaultAnimationKey: this.firstGLTFAnimation })
      const target = groups.find(item => item.id === group.id)
      if (!target) return
      target.condition = createGLTFAnimationCondition(target.condition || {})
      target.condition[prop] = value
      if (prop === 'deviceSN') target.condition.DeviceName = value
      this.commitGLTFAnimationGroups(groups)
    },
    onGLTFAnimationGroupNamesChange(group, values) {
      this.setGLTFAnimationGroupProp(group, 'animationNames', Array.isArray(values) ? values : [])
    },
    openGLTFAnimationGroupDataSelector(group) {
      this.editingGLTFAnimationGroupId = group && group.id ? group.id : ''
      this.openDataSelector('gltfAnimationGroup')
    },
    openDataSelector(mode) {
      this.dataSelectMode = mode
      if (this.$refs.deviceDataModel) {
        this.$refs.deviceDataModel.showDataModal()
      }
    },
    showDataModel() {
      this.openDataSelector('animate')
    },
    onSelectData(selectionDataInfo) {
      if (this.dataSelectMode === 'positionBinding') {
        const axis = this.editingPositionBindingAxis || 'x'
        const binding = this.getPositionBinding(axis)
        this.$set(binding, 'dataID', selectionDataInfo.uuid || '')
        this.$set(binding, 'dataName', selectionDataInfo.name || '')
        this.$set(binding, 'dataUnit', selectionDataInfo.unit || '')
        this.$set(binding, 'isBandDevice', !!selectionDataInfo.IsDevice)
        this.$set(binding, 'wsKey', selectionDataInfo.uuid || '')
        if (selectionDataInfo.IsDevice) {
          this.$set(binding, 'deviceSN', selectionDataInfo.DeviceSN || '')
          this.$set(binding, 'DeviceName', selectionDataInfo.DeviceName || '')
        } else {
          this.$set(binding, 'deviceSN', '')
          this.$set(binding, 'DeviceName', '')
        }
        this.$emit('prop-change')
        return
      }

      if (this.dataSelectMode === 'actionSetValue') {
        const action = this.currentObject && this.currentObject.action
          ? this.currentObject.action[this.editingActionIndex]
          : null
        const setValueItem = action && Array.isArray(action.setValue)
          ? action.setValue[this.editingSetValueIndex]
          : null
        if (!setValueItem) return
        this.$set(setValueItem, 'dataID', selectionDataInfo.uuid || '')
        this.$set(setValueItem, 'dataName', selectionDataInfo.name || '')
        this.$set(setValueItem, 'isBandDevice', !!selectionDataInfo.IsDevice)
        if (selectionDataInfo.IsDevice) {
          this.$set(setValueItem, 'deviceSN', selectionDataInfo.DeviceSN || '')
          this.$set(setValueItem, 'DeviceName', selectionDataInfo.DeviceName || '')
        } else {
          this.$set(setValueItem, 'deviceSN', '')
          this.$set(setValueItem, 'DeviceName', '')
        }
        this.$emit('prop-change')
        return
      }

      if (this.dataSelectMode === 'binding') {
        this.$set(this.currentObject, 'dataID', selectionDataInfo.uuid || '')
        this.$set(this.currentObject, 'dataName', selectionDataInfo.name || '')
        this.$set(this.currentObject, 'dataUnit', selectionDataInfo.unit || '')
        this.$set(this.currentObject, 'isBandDevice', !!selectionDataInfo.IsDevice)
        this.$set(this.currentObject, 'wsKey', selectionDataInfo.uuid || '')
        if (selectionDataInfo.IsDevice) {
          this.$set(this.currentObject, 'deviceSN', selectionDataInfo.DeviceSN || '')
          this.$set(this.currentObject, 'DeviceName', selectionDataInfo.DeviceName || '')
        } else {
          this.$set(this.currentObject, 'deviceSN', '')
          this.$set(this.currentObject, 'DeviceName', '')
        }
        this.$emit('prop-change')
        return
      }

      if (this.dataSelectMode === 'dataText') {
        this.$set(this.currentObject, 'dataID', selectionDataInfo.uuid || '')
        this.$set(this.currentObject, 'dataName', selectionDataInfo.name || '')
        this.$set(this.currentObject, 'dataUnit', selectionDataInfo.unit || '')
        this.$set(this.currentObject, 'isBandDevice', !!selectionDataInfo.IsDevice)
        this.$set(this.currentObject, 'wsKey', selectionDataInfo.uuid || '')
        if (selectionDataInfo.IsDevice) {
          this.$set(this.currentObject, 'deviceSN', selectionDataInfo.DeviceSN || '')
          this.$set(this.currentObject, 'DeviceName', selectionDataInfo.DeviceName || '')
        } else {
          this.$set(this.currentObject, 'deviceSN', '')
          this.$set(this.currentObject, 'DeviceName', '')
        }
        this.$emit('prop-change')
        return
      }

      if (this.dataSelectMode === 'gltfAnimation') {
        if (!this.currentObject.gltfAnimationCondition) this.$set(this.currentObject, 'gltfAnimationCondition', {})
        const target = this.currentObject.gltfAnimationCondition
        this.$set(target, 'dataID', selectionDataInfo.uuid || '')
        this.$set(target, 'dataName', selectionDataInfo.name || '')
        this.$set(target, 'isBandDevice', !!selectionDataInfo.IsDevice)
        if (selectionDataInfo.IsDevice) {
          this.$set(target, 'deviceSN', selectionDataInfo.DeviceSN || '')
          this.$set(target, 'DeviceName', selectionDataInfo.DeviceName || '')
        } else {
          this.$set(target, 'deviceSN', '')
          this.$set(target, 'DeviceName', '')
        }
        this.$emit('prop-change')
        return
      }

      if (this.dataSelectMode === 'gltfAnimationGroup') {
        const groups = ensureGLTFAnimationGroups(this.currentObject, { defaultAnimationKey: this.firstGLTFAnimation })
        const targetGroup = groups.find(group => group.id === this.editingGLTFAnimationGroupId) || groups[0]
        if (!targetGroup) return
        targetGroup.condition = createGLTFAnimationCondition(targetGroup.condition || {})
        targetGroup.condition.dataID = selectionDataInfo.uuid || ''
        targetGroup.condition.dataName = selectionDataInfo.name || ''
        targetGroup.condition.isBandDevice = !!selectionDataInfo.IsDevice
        if (selectionDataInfo.IsDevice) {
          targetGroup.condition.deviceSN = selectionDataInfo.DeviceSN || ''
          targetGroup.condition.DeviceName = selectionDataInfo.DeviceName || ''
        } else {
          targetGroup.condition.deviceSN = ''
          targetGroup.condition.DeviceName = ''
        }
        this.commitGLTFAnimationGroups(groups)
        return
      }

      if (!this.currentObject.animate) return
      this.$set(this.currentObject.animate.condition, 'dataID', selectionDataInfo.uuid || '')
      this.$set(this.currentObject.animate.condition, 'dataName', selectionDataInfo.name || '')
      this.$set(this.currentObject.animate.condition, 'isBandDevice', !!selectionDataInfo.IsDevice)
      if (selectionDataInfo.IsDevice) {
        this.$set(this.currentObject.animate.condition, 'deviceSN', selectionDataInfo.DeviceSN || '')
        this.$set(this.currentObject.animate.condition, 'DeviceName', selectionDataInfo.DeviceName || '')
      } else {
        this.$set(this.currentObject.animate.condition, 'deviceSN', '')
        this.$set(this.currentObject.animate.condition, 'DeviceName', '')
      }
      this.$emit('prop-change')
    }
  }
}
</script>

<style scoped>
.right-panel { width:100%;height:100%;background:#fafafa;border-left:1px solid #e8e8e8;display:flex;flex-direction:column;overflow-y:auto; }
.no-select-tip { padding:24px 16px;text-align:center;color:#bbb;font-size:13px; }
.no-select-tip i { font-size:32px;margin-bottom:10px;display:block; }
.prop-tabs { flex:1;display:flex;flex-direction:column;min-height:0;overflow-y:auto; }
.tab-section { margin-bottom:8px;background:#fff;border:1px solid #e8e8e8;border-radius:4px; }
.tab-section-title { padding:6px 10px;font-size:11px;font-weight:600;color:#666;background:#f5f5f5;border-bottom:1px solid #e8e8e8;letter-spacing:0.6px; }
.cip { width:32px;height:24px;border:1px solid #3d3d7a;border-radius:3px;cursor:pointer;background:transparent;padding:1px;flex-shrink:0; }
.color-row { display:flex;align-items:center;gap:6px;width:100%; }
.color-input-wrap { position:relative;flex:1;min-width:0; }
.color-input-wrap ::v-deep(.ant-input) { padding-right:28px; }
.clear-btn { position:absolute;right:6px;top:50%;transform:translateY(-50%);width:16px;height:16px;border:none;border-radius:50%;background:transparent;color:#bfbfbf;font-size:10px;line-height:16px;text-align:center;cursor:pointer;padding:0;z-index:2;display:flex;align-items:center;justify-content:center; }
.clear-btn:hover { background:#f5f5f5;color:#666; }
.clear-btn:focus { outline:none; }
.texture-row { display:flex;gap:8px;align-items:flex-start; }
.texture-preview { width:56px;height:56px;border:1px dashed #d9d9d9;border-radius:4px;display:flex;align-items:center;justify-content:center;cursor:pointer;overflow:hidden;flex-shrink:0;background:#fafafa; }
.texture-preview:hover { border-color:#13c2c2; }
.texture-thumb { width:100%;height:100%;object-fit:cover; }
.texture-actions { display:flex;flex-direction:column;gap:4px; }
.tab-section ::v-deep(.ant-slider) { margin: 4px 0; }
.tab-section ::v-deep(.ant-slider-rail) { height: 3px; background: #e8e8e8; border-radius: 2px; }
.tab-section ::v-deep(.ant-slider-track) { height: 3px; background: #13c2c2; border-radius: 2px; }
.tab-section ::v-deep(.ant-slider-handle) { width: 12px; height: 12px; margin-top: -4.5px; border: 2px solid #13c2c2; background: #fff; }
.tab-section ::v-deep(.ant-slider-handle:hover) { border-color: #36cfc9; box-shadow: 0 0 0 2px rgba(19, 194, 194, 0.2); }
.tab-section ::v-deep(.ant-form-item) { margin-bottom: 6px; }
.tab-section ::v-deep(.ant-form-item-label) { padding-bottom: 2px; }
.tab-section ::v-deep(.ant-form-item-label > label) { font-size: 12px; }
.tab-section ::v-deep(.ant-input) { height: 26px; font-size: 12px; }
.tab-section ::v-deep(.ant-input-number) { height: 26px; font-size: 12px; }
.tab-section ::v-deep(.ant-input-number-input) { font-size: 12px; }
.tab-section ::v-deep(.ant-switch) { margin-top: 0; }
.tab-section ::v-deep(.ant-select) { width: 100%; font-size: 12px; }
.tab-section ::v-deep(.ant-select-selection--single) { height: 26px; }
.tab-section ::v-deep(.ant-select-selection__rendered) { line-height: 24px; }
.label-form {
  padding: 8px 12px 6px 8px;
}
.label-form ::v-deep(.ant-form-item) {
  margin-bottom: 8px;
}
.label-form ::v-deep(.ant-form-item-label) {
  padding-right: 6px;
  line-height: 26px;
}
.label-form ::v-deep(.ant-form-item-label > label) {
  height: 26px;
  line-height: 26px;
}
.label-form ::v-deep(.ant-form-item-control) {
  line-height: 26px;
}
.label-form ::v-deep(.ant-form-item-children) {
  display: block;
  min-height: 26px;
  line-height: 26px;
}
.label-form ::v-deep(.ant-input),
.label-form ::v-deep(.ant-input-number),
.label-form ::v-deep(.ant-select-selection--single) {
  height: 26px;
}
.label-form ::v-deep(.ant-input-number-input),
.label-form ::v-deep(.ant-select-selection__rendered) {
  height: 24px;
  line-height: 24px;
}
.label-form ::v-deep(textarea.ant-input) {
  height: auto;
  min-height: 52px;
  line-height: 18px;
}
.label-form ::v-deep(.ant-slider) {
  margin: 7px 0 5px;
}
.label-form .color-row {
  height: 26px;
  align-items: center;
}
.label-form .cip {
  width: 30px;
  height: 24px;
}
.label-form ::v-deep(.ant-switch) {
  transform: scale(0.88);
  transform-origin: left center;
}
.animate-section {
  margin-bottom: 6px;
}
.animate-section .tab-section-title {
  padding: 4px 10px;
}
.animate-form {
  padding: 8px 10px 6px 8px;
}
.animate-form ::v-deep(.ant-form-item) {
  margin-bottom: 8px;
}
.animate-form ::v-deep(.ant-form-item-label) {
  padding-right: 6px;
  line-height: 24px;
}
.animate-form ::v-deep(.ant-form-item-control) {
  line-height: 24px;
}
.animate-form ::v-deep(.ant-form-item-children) {
  min-height: 24px;
}
.animate-form ::v-deep(.ant-switch) {
  transform: scale(0.88);
  transform-origin: left center;
}
.animate-form ::v-deep(.ant-slider) {
  margin: 5px 0 3px;
}
.animate-form ::v-deep(.ant-input),
.animate-form ::v-deep(.ant-input-number),
.animate-form ::v-deep(.ant-select-selection--single) {
  height: 24px;
}
.animate-form ::v-deep(.ant-input-number-input),
.animate-form ::v-deep(.ant-select-selection__rendered) {
  line-height: 22px;
}
.animate-bind-form ::v-deep(.ant-form-item) { margin-bottom: 8px; }
.data-bind-section {
  margin: 0 0 6px 0;
}
.data-bind-form {
  padding: 8px 12px 2px 8px;
}
.data-bind-form ::v-deep(.ant-form-item) {
  margin-bottom: 8px;
}
.data-bind-form ::v-deep(.ant-form-item-label) {
  padding-right: 6px;
  line-height: 24px;
}
.data-bind-form ::v-deep(.ant-form-item-control) {
  line-height: 24px;
}
.data-bind-form ::v-deep(.ant-form-item-children) {
  min-height: 24px;
}
.data-bind-empty {
  padding: 12px;
  color: #999;
  font-size: 12px;
  text-align: center;
}
.material-section {
  margin: 0;
  border: none;
}
.material-form {
  padding: 8px 12px 6px 8px;
}
.material-form ::v-deep(.ant-form-item) {
  margin-bottom: 8px;
}
.material-form ::v-deep(.ant-form-item-label) {
  padding-right: 6px;
  line-height: 26px;
}
.material-form ::v-deep(.ant-form-item-control) {
  line-height: 26px;
}
.material-form ::v-deep(.ant-form-item-children) {
  min-height: 24px;
}
.material-form ::v-deep(.ant-input) {
  height: 24px;
}
.material-form ::v-deep(.ant-switch) {
  transform: scale(0.88);
  transform-origin: left center;
}
.material-form .color-row {
  gap: 6px;
}
.material-form .cip {
  width: 30px;
  height: 24px;
}
.material-form ::v-deep(.ant-slider) {
  margin: 5px 0 3px;
}
.material-form .texture-row {
  gap: 8px;
  align-items: center;
}
.material-form .texture-preview {
  width: 42px;
  height: 42px;
}
.material-form .texture-actions {
  gap: 4px;
}
.material-form ::v-deep(.ant-btn) {
  height: 24px;
  padding: 0 8px;
  font-size: 12px;
}
.event-section {
  margin: 0;
  padding: 8px;
  border: none;
  background: #f7f8fa;
}
.action-card {
  margin-bottom: 8px;
  border: 1px solid #e5e9ef;
  border-radius: 4px;
  background: #fff;
  overflow: hidden;
}
.action-card-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 30px;
  padding: 0 10px;
  color: #333;
  font-size: 12px;
  font-weight: 600;
  background: #f6f7f9;
  border-bottom: 1px solid #edf0f3;
}
.action-delete {
  cursor: pointer;
  font-size: 14px;
}
.action-card > .ant-form,
.action-card > form {
  padding: 8px 10px 6px 6px;
}
.action-card ::v-deep(.ant-form-item) {
  margin-bottom: 8px;
}
.action-card ::v-deep(.ant-form-item-label) {
  padding-right: 6px;
  line-height: 24px;
}
.action-card ::v-deep(.ant-form-item-control) {
  line-height: 24px;
}
.action-card ::v-deep(.ant-form-item-children) {
  min-height: 24px;
}
.action-card ::v-deep(.ant-input),
.action-card ::v-deep(.ant-select-selection--single) {
  height: 24px;
}
.action-card ::v-deep(.ant-select-selection__rendered) {
  line-height: 22px;
}
.action-card ::v-deep(.ant-select-selection--multiple) {
  min-height: 24px;
}
.action-card ::v-deep(.ant-switch) {
  transform: scale(0.88);
  transform-origin: left center;
}
.event-add {
  padding-top: 2px;
}
.event-add ::v-deep(.ant-btn) {
  height: 28px;
  border: none;
  border-radius: 3px;
  background: #13c2c2;
}
.set-value-actions {
  margin: -2px 0 8px 96px;
}
.set-value-card {
  margin: 0 0 8px 18px;
  padding: 8px 8px 2px;
  border: 1px solid #edf0f3;
  border-radius: 4px;
  background: #fafbfc;
}
.set-value-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 6px;
  font-size: 12px;
  color: #666;
}
.bind-animation-group {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  column-gap: 12px;
  row-gap: 6px;
  line-height: 22px;
}
.bind-animation-group ::v-deep(.ant-checkbox-wrapper) {
  margin-left: 0;
  white-space: nowrap;
  font-size: 12px;
}
.model-animation-empty {
  padding: 0 12px 8px;
  color: #999;
  font-size: 12px;
}
.model-animation-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.model-animation-group {
  margin: 8px 8px 10px;
  padding: 8px 8px 4px;
  border: 1px solid #e8e8e8;
  border-radius: 4px;
  background: #fff;
}
.model-animation-group-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}
.model-animation-group-name {
  flex: 1 1 auto;
  min-width: 0;
}
.data-select-group {
  display: flex !important;
  width: 100%;
}
.data-select-group ::v-deep(.ant-input) {
  flex: 1 1 auto;
  width: auto !important;
  min-width: 0;
  height: 24px;
  border-right: 0;
  border-radius: 4px 0 0 4px;
}
.data-select-group ::v-deep(.ant-btn) {
  flex: 0 0 88px;
  width: 88px !important;
  height: 24px;
  padding: 0 8px;
  font-size: 12px;
  color: #fff;
  background: #13c2c2;
  border-color: #13c2c2;
  border-left: 0;
  border-radius: 0 4px 4px 0;
  box-shadow: none;
}
.data-select-group ::v-deep(.ant-btn:hover),
.data-select-group ::v-deep(.ant-btn:focus) {
  color: #fff;
  background: #10adad;
  border-color: #10adad;
}
.data-select-group ::v-deep(.ant-btn i) {
  margin-right: 3px;
}
</style>
