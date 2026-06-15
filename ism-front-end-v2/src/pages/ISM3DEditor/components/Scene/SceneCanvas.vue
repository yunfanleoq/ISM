<template>
  <div class="center-panel" ref="canvasContainer">
    <div class="canvas-toolbar">
      <button class="toolbar-btn" @click="setCameraView('perspective')" :class="{active:cameraView==='perspective'}"><i class="fas fa-camera"></i> {{ $t('ISM3DEditor.perspective') }}</button>
      <button class="toolbar-btn" @click="setCameraView('top')" :class="{active:cameraView==='top'}"><i class="fas fa-arrow-down"></i> {{ $t('ISM3DEditor.top') }}</button>
      <button class="toolbar-btn" @click="setCameraView('front')" :class="{active:cameraView==='front'}"><i class="fas fa-arrow-right"></i> {{ $t('ISM3DEditor.front') }}</button>
      <button class="toolbar-btn" @click="setCameraView('right')" :class="{active:cameraView==='right'}"><i class="fas fa-arrow-left"></i> {{ $t('ISM3DEditor.side') }}</button>
      <div class="toolbar-divider"></div>
      <button class="toolbar-btn" @click="$emit('focus-selected')" :disabled="!selectedId"><i class="fas fa-crosshairs"></i> {{ $t('ISM3DEditor.focusObject') }}</button>
      <button class="toolbar-btn" @click="$emit('frame-all')"><i class="fas fa-compress-arrows-alt"></i> {{ $t('ISM3DEditor.frame') }}</button>
      <div class="canvas-info">
        <span :title="$t('ISM3DEditor.fps')">FPS: {{ fps }}</span>
        <span :title="$t('ISM3DEditor.objectCountTitle')">{{ $t('ISM3DEditor.objects') }}: {{ objects.length }}</span>
        <span v-if="hoveredName" style="color:#13c2c2">{{ hoveredName }}</span>
      </div>
    </div>

    <div id="three-canvas-wrapper" ref="canvasWrapper"
      @dragover.prevent
      @drop="onCanvasDrop"
    >
      <canvas id="three-canvas" ref="threeCanvas"></canvas>
      <div class="overlay-2d-layer">
        <div
          v-for="overlay in overlay2DItems"
          :key="overlay.id"
          class="overlay-2d-item"
          :class="{ selected: selectedId === overlay.id, 'has-action': overlay.hasAction }"
          :style="overlay.style"
          @mousedown.stop="onOverlayMouseDown($event, overlay)"
          @click.stop="onOverlayClick(overlay)"
          @contextmenu.prevent.stop="onOverlayContextMenu($event, overlay)"
        >
          <div class="overlay-2d-content">
            <ISM2DNodeRenderer
              v-if="overlay.kind === '2dComponent'"
              :object-data="overlay.objectData"
              :showDeviceUuid="''"
              :editMode="true"
              :selected="selectedId === overlay.id"
            />
            <div v-else-if="overlay.kind === 'uiLabel'" class="basic-ui-label" :style="overlay.contentStyle">
              {{ overlay.text }}
            </div>
            <img v-if="overlay.kind === 'uiImage' && overlay.src" class="basic-ui-image" :src="overlay.src" />
            <div v-else-if="overlay.kind === 'uiImage'" class="basic-ui-image basic-ui-image-empty">
              <i class="far fa-image"></i>
            </div>
            <iframe v-else-if="overlay.kind === 'webEmbed'" class="basic-ui-frame" :src="overlay.src"></iframe>
          </div>
          <template v-if="selectedId === overlay.id">
            <div class="overlay-2d-selection-box"></div>
            <div class="overlay-2d-resize-handle" @mousedown.stop="onOverlayHandleMouseDown($event, overlay, 'resize')"></div>
            <div class="overlay-2d-resize-handle overlay-2d-resize-handle-tl" @mousedown.stop="onOverlayHandleMouseDown($event, overlay, 'resize-tl')"></div>
            <div class="overlay-2d-resize-handle overlay-2d-resize-handle-tr" @mousedown.stop="onOverlayHandleMouseDown($event, overlay, 'resize-tr')"></div>
            <div class="overlay-2d-resize-handle overlay-2d-resize-handle-bl" @mousedown.stop="onOverlayHandleMouseDown($event, overlay, 'resize-bl')"></div>
            <div
              class="overlay-2d-rotate-handle"
              @mousedown.stop="onOverlayHandleMouseDown($event, overlay, 'rotate')"
              :style="getRotateHandleStyle(overlay)"
            ></div>
          </template>
        </div>
      </div>

      <div class="selected-indicator" v-if="selectedObj" :key="selectedId">
        <i class="fas fa-hand-pointer" style="margin-right:6px"></i>{{ selectedObj.name }}
      </div>

      <div
        v-if="selectionFrame.visible"
        class="transform-selection-frame"
        :style="{
          left: selectionFrame.left + 'px',
          top: selectionFrame.top + 'px',
          width: selectionFrame.width + 'px',
          height: selectionFrame.height + 'px'
        }"
      >
        <div class="kf-border"></div>
        <span class="kf-c kf-tl" @mousedown.stop.prevent="onSelectionFrameHandleMouseDown($event, 'tl')"></span>
        <span class="kf-c kf-tr" @mousedown.stop.prevent="onSelectionFrameHandleMouseDown($event, 'tr')"></span>
        <span class="kf-c kf-bl" @mousedown.stop.prevent="onSelectionFrameHandleMouseDown($event, 'bl')"></span>
        <span class="kf-c kf-br" @mousedown.stop.prevent="onSelectionFrameHandleMouseDown($event, 'br')"></span>
      </div>

      <div class="drop-hint" v-if="objects.length===0">
        <i class="fas fa-cube"></i>
        <p>{{ $t('ISM3DEditor.dragHint') }}</p>
      </div>

      <svg class="axis-indicator" viewBox="0 0 70 70" xmlns="http://www.w3.org/2000/svg">
        <line x1="35" y1="35" x2="62" y2="20" stroke="#ff4d4f" stroke-width="2"/>
        <text x="64" y="18" fill="#ff4d4f" font-size="10" font-family="Arial">X</text>
        <line x1="35" y1="35" x2="10" y2="20" stroke="#52c41a" stroke-width="2"/>
        <text x="3" y="18" fill="#52c41a" font-size="10" font-family="Arial">Y</text>
        <line x1="35" y1="35" x2="35" y2="60" stroke="#13c2c2" stroke-width="2"/>
        <text x="30" y="68" fill="#13c2c2" font-size="10" font-family="Arial">Z</text>
      </svg>
    </div>

    <div class="canvas-overlay-hint">
      {{ $t('ISM3DEditor.canvasHint') }}
    </div>

    <div class="status-bar">
      <span><i class="fas fa-circle" style="font-size:7px;color:#13c2c2"></i> {{ $t('ISM3DEditor.ready') }}</span>
      <span v-if="selectedObj">{{ $t('ISM3DEditor.selectedPrefix') }}{{ selectedObj.name }}</span>
      <span v-if="selectedObj">{{ $t('ISM3DEditor.positionPrefix') }}({{ fmt(selectedObj.x) }}, {{ fmt(selectedObj.y) }}, {{ fmt(selectedObj.z) }})</span>
      <span style="margin-left:auto;opacity:.7">{{ modeLabel }}</span>
    </div>
  </div>
</template>

<script>
import * as THREE from 'three'
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls'
import { TransformControls } from 'three/examples/jsm/controls/TransformControls'
import { GLTFLoader } from 'three/examples/jsm/loaders/GLTFLoader'
import { RGBELoader } from 'three/examples/jsm/loaders/RGBELoader.js'
import { createThreeMesh as createThreeMeshFromLib, create2DComponentPlane, createMediaPlane, updateMediaPlane, disposeMediaPlane, createTextSprite, updateTextSprite, createText3DMesh, updateText3DMesh, createFlowPipe, updateFlowPipeAnimation, updateFlowPipe, isFlowPipeSegmented, tryLoadChineseFontFor3D } from '../Objects/IndustrialObjects'
import ISM2DNodeRenderer from './ISM2DNodeRenderer.vue'
import { ensureGLTFAnimationGroups, remapGLTFAnimationGroupNames, syncLegacyGLTFAnimationFields } from '../../utils/GLTFAnimationGroups'
import Vue from 'vue'
import Contextmenu from "vue-contextmenujs"

Vue.use(Contextmenu)

/** 将任意颜色值转为 '#rrggbb' 或 '#rrggbbaa' 格式 */
function normColorValue(val, fallback) {
  fallback = fallback || '#000000'
  if (!val || typeof val !== 'string') return fallback
  var s = val.trim()
  var COLOR_KEYS = {
    blue: '#4da6ff', green: '#4dffa6', orange: '#ffaa4d',
    purple: '#c17bff', pink: '#ff7bb5', cyan: '#4dfff0',
    yellow: '#fffb4d', red: '#ff4d4f', gray: '#8c8c8c', brown: '#a67c52'
  }
  if (COLOR_KEYS[s]) return COLOR_KEYS[s]
  if (s.charAt(0) === '#') {
    if (/^#[0-9a-fA-F]{3}$/.test(s)) {
      return '#' + s[1] + s[1] + s[2] + s[2] + s[3] + s[3]
    }
    if (/^#[0-9a-fA-F]{6}$/.test(s)) return s.toLowerCase()
    if (/^#[0-9a-fA-F]{8}$/.test(s)) return s.toLowerCase()
  }
  return fallback
}

function buildTextLabelOptions(obj) {
  obj = obj || {}
  var bgOpacity = obj.textBgOpacity !== undefined ? Number(obj.textBgOpacity) : undefined
  if (obj.type === 'text3d' || obj.type === 'textPlain3d' || obj.type === 'dataText') {
    bgOpacity = bgOpacity !== undefined && isFinite(bgOpacity) ? bgOpacity : undefined
  }
  return {
    opacity: obj.labelOpacity !== undefined ? Number(obj.labelOpacity) : 1,
    backgroundOpacity: bgOpacity,
    faceCamera: obj.labelFaceCamera !== false,
    fixedSize: obj.labelFixedSize !== false,
    fontFamily: obj.labelFontFamily || '系统默认',
    renderMode: obj.labelRenderMode || 'component',
    true3D: !(bgOpacity > 0),
    showBorder: !!obj.labelShowBorder,
    borderWidth: obj.labelBorderWidth !== undefined ? Number(obj.labelBorderWidth) : 1,
    borderColor: obj.labelBorderColor || '#ffffff'
  }
}

function isTextLabelObject(obj) {
  return obj && (obj.type === 'text3d' || obj.type === 'dataText' || obj.type === 'textPlain3d' || obj.type === 'uiLabel')
}

function normTextLabelColor(value) {
  if (typeof value === 'number') return '#' + value.toString(16).padStart(6, '0')
  if (typeof value === 'string' && value.trim().charAt(0) === '#') {
    return normColorValue(value, '#ffffff')
  }
  return '#ffffff'
}

function applyTextSpriteObjectScale(mesh, obj) {
  if (!mesh || !obj) return
  var sx = obj.sx !== undefined ? Number(obj.sx) : 1
  var sy = obj.sy !== undefined ? Number(obj.sy) : 1
  var sz = obj.sz !== undefined ? Number(obj.sz) : 1
  if (!isFinite(sx)) sx = 1
  if (!isFinite(sy)) sy = 1
  if (!isFinite(sz)) sz = 1
  // 3D 文字 Mesh: 绝对缩放 (geometry 已含物理尺寸)
  if (mesh.userData && mesh.userData.is3DText) {
    mesh.scale.set(sx, sy, sz)
  } else {
    mesh.scale.set(mesh.scale.x * sx, mesh.scale.y * sy, mesh.scale.z * sz)
  }
}

function updateTextPickArea(mesh, obj) {
  if (!mesh || !obj || !mesh.userData || !mesh.userData.isTextSprite) return
  var localBox = mesh.geometry && mesh.geometry.boundingBox ? mesh.geometry.boundingBox.clone() : null
  if (!localBox && mesh.geometry && mesh.geometry.computeBoundingBox) {
    mesh.geometry.computeBoundingBox()
    localBox = mesh.geometry.boundingBox ? mesh.geometry.boundingBox.clone() : null
  }
  if (!localBox) localBox = new THREE.Box3().setFromObject(mesh)
  var size = localBox.getSize(new THREE.Vector3())
  var center = localBox.getCenter(new THREE.Vector3())
  var width = Math.max(Math.abs(size.x), 0.2)
  var height = Math.max(Math.abs(size.y), 0.12)
  var depth = Math.max(Math.abs(size.z), 0.02)
  var pick = mesh.userData.textPickArea
  if (!pick) {
    pick = new THREE.Mesh(
      new THREE.PlaneGeometry(1, 1),
      new THREE.MeshBasicMaterial({
        transparent: true,
        opacity: 0,
        depthWrite: false,
        side: THREE.DoubleSide
      })
    )
    pick.name = 'TextPickArea_' + (obj.id || '')
    pick.userData = { id: obj.id, isTextPickArea: true }
    mesh.add(pick)
    mesh.userData.textPickArea = pick
  }
  pick.userData.id = obj.id
  pick.scale.set(width, height, 1)
  pick.position.set(center.x, center.y, center.z + depth / 2 + 0.01)
}

function getGLTFTargetSize(objData) {
  var value = objData && objData.fitSize !== undefined ? Number(objData.fitSize) : 2
  if (!isFinite(value) || value <= 0) value = 2
  return value
}

function centerGLTFModelGeometry(model) {
  if (!model) return
  model.updateMatrixWorld(true)
  var box = new THREE.Box3().setFromObject(model)
  if (box.isEmpty()) return
  var centerWorld = box.getCenter(new THREE.Vector3())
  var center = model.worldToLocal(centerWorld.clone())
  model.children.forEach(function(child) {
    child.position.sub(center)
  })
  model.updateMatrixWorld(true)
}

function isUiOverlayObject(obj) {
  return obj && (obj.type === 'uiLabel' || obj.type === 'uiImage' || obj.type === 'webEmbed')
}

function getTextLabelContent(obj) {
  if (!obj) return ''
  if (obj.type === 'dataText') {
    var value = obj.realTimeValue !== undefined && obj.realTimeValue !== '' ? obj.realTimeValue : '--'
    var format = obj.dataFormat && obj.dataFormat !== '{value}' ? obj.dataFormat : '{value}'
    return format.replace('{value}', value)
  }
  if (obj.type === 'uiLabel') return obj.textContent || 'UI标签'
  if (obj.type === 'textPlain3d') return obj.textContent || '3D文字'
  return obj.textContent || '标签'
}

export default {
  name: 'SceneCanvas',
  components: {
    ISM2DNodeRenderer
  },
  props: {
    selectedId: { type: String, default: null },
    selectedObj: { type: Object, default: null },
    objects: { type: Array, default: () => [] },
    mode: { type: String, default: 'select' },
    gridSize: { type: Number, default: 10 },
    gridColor: { type: String, default: '#444444' },
  },
  i18n: require('@/i18n/language'),
  data() {
    return {
      cameraView: 'perspective',
      fps: 0,
      hoveredName: '',
      _fpsFrames: 0,
      _fpsLast: 0,
      _animId: null,
      _handleCanvasClick: null,
      _handleCanvasContextMenu: null,
      _handleCanvasMouseMove: null,
      _handleWindowResize: null,
      _handleWindowMouseMove: null,
      _handleWindowMouseUp: null,
      _handleWindowKeyDown: null,
      _handleWindowKeyUp: null,
      _overlayDrag: null,
      _selectionResizeDrag: null,
      _transformAutoHideTimer: null,
      _ignoreNextCanvasClickAfterTransform: false,
      _spacePanActive: false,
      overlay2DItems: [],
      isDrawingPipe: false,
      pipePoints: [],
      pipePreviewLine: null,      // 旧预览，保留兼容
      pipePreviewMesh: null,      // 实时管道预览 Mesh
      pipePreviewPoints: [],
      dragGuide: {
        visible: false,
        left: 0,
        top: 0,
        x: '0.00',
        y: '0.00',
        z: '0.00',
        name: '',
        alignText: ''
      },
      selectionFrame: {
        visible: false,
        left: 0,
        top: 0,
        width: 0,
        height: 0
      }
    }
  },
  i18n: require('@/i18n/language'),
  computed: {
    modeLabel() {
      var map = { select: '选择模式', move: '移动模式', rotate: '旋转模式', scale: '缩放模式' }
      return map[this.mode] || ''
    }
  },
  watch: {
    objects: {
      handler: function(newList) {
        this.syncMeshes(newList)
        this.update2DOverlays()
      },
      deep: true
    },
    selectedId: function(id) {
      this.updateSelection(id)
      if (!id) {
        this.dragGuide.visible = false
      }
    },
    mode: function(val) {
      this.setTransformMode(val)
    },
  },
  mounted() {
    var self = this
    this.$nextTick(function() {
      self.initThree()
      self.startRenderLoop()
      // 预加载中文字体，用于真正的3D几何体文字
      if (typeof window !== 'undefined' && window.opentype && tryLoadChineseFontFor3D) {
        tryLoadChineseFontFor3D()
      }
      // 字体加载完成后，升级场景中已有的 Canvas fallback 文字为 true 3D
      self._onFontReady = function() {
        self._upgradeCanvasFallbackTexts()
      }
      window.addEventListener('ism3d-font-ready', self._onFontReady)
    })
  },
  beforeDestroy() {
    if (this._onFontReady) {
      window.removeEventListener('ism3d-font-ready', this._onFontReady)
    }
    if (this._transformAutoHideTimer) {
      clearTimeout(this._transformAutoHideTimer)
      this._transformAutoHideTimer = null
    }
    if (this._handleWindowResize) {
      window.removeEventListener('resize', this._handleWindowResize)
    }
    if (this._handleWindowMouseMove) {
      window.removeEventListener('mousemove', this._handleWindowMouseMove)
    }
    if (this._handleWindowMouseUp) {
      window.removeEventListener('mouseup', this._handleWindowMouseUp)
    }
    if (this._handleWindowKeyDown) {
      window.removeEventListener('keydown', this._handleWindowKeyDown)
    }
    if (this._handleWindowKeyUp) {
      window.removeEventListener('keyup', this._handleWindowKeyUp)
    }
    if (this.$refs.threeCanvas) {
      if (this._handleCanvasClick) {
        this.$refs.threeCanvas.removeEventListener('click', this._handleCanvasClick)
      }
      if (this._handleCanvasContextMenu) {
        this.$refs.threeCanvas.removeEventListener('contextmenu', this._handleCanvasContextMenu)
      }
      if (this._handleCanvasMouseMove) {
        this.$refs.threeCanvas.removeEventListener('mousemove', this._handleCanvasMouseMove)
      }
      if (this._handleCanvasDoubleClick) {
        this.$refs.threeCanvas.removeEventListener('dblclick', this._handleCanvasDoubleClick)
      }
      if (this._handleCanvasMouseDown) {
        this.$refs.threeCanvas.removeEventListener('mousedown', this._handleCanvasMouseDown)
      }
    }
    if (this._animId) cancelAnimationFrame(this._animId)
    this.disposeThree()
  },
  methods: {
    /**
     * 字体加载完成后，将场景中还是 Canvas fallback（is3DText=true, isTrue3DText 不存在）
     * 的文字升级为 true 3D 几何体文字，并重建已有 True3D 文字以应用新加载的字体
     */
    _upgradeCanvasFallbackTexts: function() {
      if (!this.__3d || !this.__3d.scene) return
      var scene = this.__3d.scene
      var toUpgrade = []
      var toRebuild = []
      scene.traverse(function(obj) {
        if (obj.userData && obj.userData.is3DText) {
          if (!obj.userData.isTrue3DText) {
            toUpgrade.push(obj)
          } else {
            toRebuild.push(obj)
          }
        }
      })
      if (toUpgrade.length === 0 && toRebuild.length === 0) return
      if (toUpgrade.length) {
        console.log('[SceneCanvas] Upgrading', toUpgrade.length, 'Canvas fallback text(s) to true 3D')
      }
      if (toRebuild.length) {
        console.log('[SceneCanvas] Rebuilding', toRebuild.length, 'True3D text(s) for font update')
      }
      var all = toUpgrade.concat(toRebuild)
      for (var i = 0; i < all.length; i++) {
        var mesh = all[i]
        var td = mesh.userData.textData || {}
        updateText3DMesh(
          mesh,
          td.text || '',
          td.color,
          td.fontSize,
          td.bgColor,
          td.options
        )
      }
    },
    applyMaterialProps: function(material, obj) {
      if (!material || !obj) return

      // 只有明确设置了颜色才修改，否则保留原始材质颜色
      if (obj.color) {
        var mc = normColorValue(obj.color, '#4a90d9')
        var c = parseInt(mc.replace('#', ''), 16)
        if (material.color && material.color.setHex) {
          material.color.setHex(c)
        }
      }

      // 只有明确设置了 opacity 才修改
      // 闪烁/浮动/自旋动画需要透明材质支持，避免 syncMeshes 重置 transparent = false
      if (obj.opacity !== undefined) {
        var needsTransparent = obj.opacity < 1 || obj.blink || obj.floatAnim || obj.animateVisible
        material.opacity = obj.opacity
        material.transparent = needsTransparent
      }

      // 只有明确设置了 wireframe 才修改
      if (obj.wireframe !== undefined) {
        material.wireframe = !!obj.wireframe
      }

      if (material.metalness !== undefined && obj.metalness !== undefined) {
        material.metalness = obj.metalness
      }

      if (material.roughness !== undefined && obj.roughness !== undefined) {
        material.roughness = obj.roughness
      }

      // 只有明确设置了 emissive 才修改
      if (obj.emissive && material.emissive && material.emissive.setHex) {
        var emissive = normColorValue(obj.emissive, '#000000')
        material.emissive.setHex(parseInt(emissive.replace('#', ''), 16))
      }

      // 贴图处理：textureData 为图片 URL，来自项目图片组件或上传
      if (obj.textureData && 'map' in material) {
        // 防止重复加载相同贴图
        if (material.userData && material.userData._textureKey === obj.textureData) {
          // 相同贴图，跳过
          console.log('[贴图] 跳过重复贴图:', obj.id, obj.textureData)
        } else {
          var loader = new THREE.TextureLoader()
          var self = this
          // 换贴图前，先释放之前由用户设置的旧贴图
          if (material.userData && material.userData._userTexture) {
            console.log('[贴图] 释放旧贴图:', obj.id)
            material.userData._userTexture.dispose()
          }
          // 如果原始 GLTF 材质有默认贴图，保留它但允许覆盖
          var originalMap = material.map
          console.log('[贴图] 开始加载新贴图:', obj.id, 'URL:', obj.textureData, '原始map:', !!originalMap)
          var tex = loader.load(
            obj.textureData,
            function(texture) {
              console.log('[贴图] 贴图加载成功:', obj.id)
              material.map = texture
              material.needsUpdate = true
              if (self && self.__3d) {
                self.__3d.renderer.render(self.__3d.scene, self.__3d.camera)
              }
            },
            undefined,
            function(err) {
              console.error('[贴图] 贴图加载失败:', obj.id, 'URL:', obj.textureData, err)
            }
          )
          tex.colorSpace = THREE.SRGBColorSpace
          tex.wrapS = THREE.RepeatWrapping
          tex.wrapT = THREE.RepeatWrapping
          // 先设置到 material.map，即使纹理还未加载完毕也会显示占位
          material.map = tex
          if (!material.userData) material.userData = {}
          material.userData._userTexture = tex
          material.userData._textureKey = obj.textureData
          // 如果存在原始贴图，标记以便后续可恢复
          if (originalMap && !material.userData._originalTexture) {
            material.userData._originalTexture = originalMap
          }
          console.log('[贴图] material.map 已设置:', obj.id, '材质类型:', material.type)
        }
      } else if (obj.textureData === '' && material.userData && material.userData._userTexture) {
        // 用户清除了贴图
        console.log('[贴图] 清除贴图:', obj.id)
        material.userData._userTexture.dispose()
        material.userData._userTexture = null
        material.userData._textureKey = null
        material.map = material.userData._originalTexture || null
        material.needsUpdate = true
      }

      material.needsUpdate = true
    },
    /** 同步 3D mesh */
    applyObjectMaterialProps: function(mesh, obj) {
      if (!mesh || !obj || (mesh.userData && mesh.userData.isTextSprite)) return
      var vm = this
      var applyOne = function(material) {
        if (!material || material.isShaderMaterial) return
        vm.applyMaterialProps(material, obj)
      }
      if (mesh.material) {
        if (Array.isArray(mesh.material)) {
          mesh.material.forEach(applyOne)
        } else {
          applyOne(mesh.material)
        }
      }
      if (mesh.traverse) {
        mesh.traverse(function(child) {
          if ((!child.isMesh && !child.isPoints && !child.isLineSegments) || !child.material || child === mesh) return
          if (Array.isArray(child.material)) {
            child.material.forEach(applyOne)
          } else {
            applyOne(child.material)
          }
        })
      }
    },

    applyLightProps: function(mesh, obj) {
      if (!mesh || !obj || !obj.type || obj.type.indexOf('light') === -1) return
      var color = normColorValue(obj.color, '#ffffff')
      var colorInt = parseInt(color.replace('#', ''), 16)
      var intensity = obj.intensity !== undefined ? Number(obj.intensity) : 1
      var distance = obj.distance !== undefined ? Number(obj.distance) : undefined
      var angle = obj.angle !== undefined ? Number(obj.angle) : undefined
      var penumbra = obj.penumbra !== undefined ? Number(obj.penumbra) : undefined
      mesh.traverse(function(child) {
        if (child.isLight) {
          if (child.color && child.color.set) child.color.set(colorInt)
          if (!isNaN(intensity)) child.intensity = intensity
          if ((child.isPointLight || child.isSpotLight) && distance !== undefined && !isNaN(distance)) child.distance = distance
          if (child.isSpotLight) {
            if (angle !== undefined && !isNaN(angle)) child.angle = angle
            if (penumbra !== undefined && !isNaN(penumbra)) child.penumbra = penumbra
          }
        } else if (child.isMesh && child.material && child.material.color) {
          child.material.color.set(colorInt)
          if (child.material.emissive) child.material.emissive.set(colorInt)
          child.material.needsUpdate = true
        } else if (child.type === 'ArrowHelper' && child.setColor) {
          child.setColor(colorInt)
        }
      })
    },

    syncMeshes: function(objects) {
      if (!this.__3d) return
      var vm = this
      var scene = this.__3d.scene
      var meshMap = this.__3d.meshMap
      var reflectionMap = this.__3d.reflectionMap || {}
      this.__3d.reflectionMap = reflectionMap

      // 删除已移除的 mesh
      var objIds = new Set(objects.map(function(o) { return o.id }))
      Object.keys(meshMap).forEach(function(id) {
        if (!objIds.has(id)) {
          var mesh = meshMap[id]
          scene.remove(mesh)
          if (mesh.traverse) {
            mesh.traverse(function(child) {
              if (child.userData && child.userData.isMediaPlane) disposeMediaPlane(child)
              if (child.geometry) child.geometry.dispose()
              if (child.material) {
                if (Array.isArray(child.material)) child.material.forEach(function(m) { m.dispose() })
                else child.material.dispose()
              }
            })
          } else {
            if (mesh.userData && mesh.userData.isMediaPlane) disposeMediaPlane(mesh)
            if (mesh.geometry) mesh.geometry.dispose()
            if (mesh.material) {
              if (Array.isArray(mesh.material)) mesh.material.forEach(function(m) { m.dispose() })
              else mesh.material.dispose()
            }
          }
          delete meshMap[id]
          vm.removeObjectReflection(id)
        }
      })

      // 创建或更新 mesh
      objects.forEach(function(obj) {
        var mesh = meshMap[obj.id]
        if (isUiOverlayObject(obj)) {
          if (mesh) {
            scene.remove(mesh)
            if (mesh.traverse) {
              mesh.traverse(function(child) {
                if (child.userData && child.userData.isMediaPlane) disposeMediaPlane(child)
                if (child.geometry) child.geometry.dispose()
                if (child.material) {
                  if (Array.isArray(child.material)) child.material.forEach(function(m) { m.dispose() })
                  else child.material.dispose()
                }
              })
            } else {
              if (mesh.userData && mesh.userData.isMediaPlane) disposeMediaPlane(mesh)
              if (mesh.geometry) mesh.geometry.dispose()
              if (mesh.material) {
                if (Array.isArray(mesh.material)) mesh.material.forEach(function(m) { m.dispose() })
                else mesh.material.dispose()
              }
            }
            delete meshMap[obj.id]
          }
          return
        }
        if (obj.type === 'flowPipe' && mesh && !isFlowPipeSegmented(mesh)) {
          scene.remove(mesh)
          if (mesh.traverse) {
            mesh.traverse(function(child) {
              if (child.geometry) child.geometry.dispose()
              if (child.material) {
                if (Array.isArray(child.material)) child.material.forEach(function(m) { m.dispose() })
                else child.material.dispose()
              }
            })
          } else {
            if (mesh.geometry) mesh.geometry.dispose()
            if (mesh.material) mesh.material.dispose()
          }
          delete meshMap[obj.id]
          mesh = null
        }
        // GLTF 模型的 mesh 已在加载时创建，不需要重复创建
        if (!mesh && obj.type !== 'gltf') {
          var colorStr = normColorValue(obj.color, '#4a90d9')
          if (isTextLabelObject(obj)) colorStr = normTextLabelColor(obj.color)
          var colorInt = parseInt(colorStr.replace('#', ''), 16)
          if (obj.type === 'text3d') {
            var result = createText3DMesh(obj.textContent || '标签', colorInt, obj.fontSize || 16, obj.textBgColor || '#00000000', buildTextLabelOptions(obj))
            mesh = result ? result.mesh : null
          } else if (obj.type === 'dataText') {
            var result2 = createText3DMesh(getTextLabelContent(obj), colorInt, obj.fontSize || 16, obj.textBgColor || '#00000000', buildTextLabelOptions(obj))
            mesh = result2 ? result2.mesh : null
          } else if (obj.type === 'textPlain3d' || obj.type === 'uiLabel') {
            var result3 = createText3DMesh(getTextLabelContent(obj), colorInt, obj.fontSize || 16, obj.textBgColor || '#00000000', buildTextLabelOptions(obj))
            mesh = result3 ? result3.mesh : null
          } else if (obj.type === 'image3d' || obj.type === 'video3d') {
            mesh = createMediaPlane(obj, obj.type)
          } else if (obj.type === 'flowPipe') {
            mesh = createFlowPipe(obj.points, obj.color, obj.radius, obj.flowSpeed, obj.highlightColor, obj.flowDashLength)
          } else if (obj.type === '2dComponent') {
            mesh = create2DComponentPlane(obj)
          } else {
            mesh = createThreeMeshFromLib(obj.type, colorInt)
          }
          if (mesh) {
            if (obj.type === 'flowRibbon' || obj.type === 'energyWall') {
              console.log('[syncMeshes] 创建特效mesh:', obj.id, 'type=' + obj.type, 'hasAnimation=' + mesh.userData.hasAnimation)
            }
            mesh.userData = mesh.userData || {}
            mesh.userData.id = obj.id
            scene.add(mesh)
            meshMap[obj.id] = mesh
            if (obj.id === vm.selectedId) {
              vm.attachTransform(mesh)
              vm.updateActiveSelectionFrame()
            }
          }
        }
        if (mesh) {
          vm.applyObjectShadowState(mesh, obj)
          var isTransforming = vm.__3d.transform.dragging && vm.__3d.transform.object === mesh
          if (!isTransforming) {
            if (!obj.floatAnim) {
              mesh.position.set(obj.x || 0, obj.y || 0, obj.z || 0)
            } else {
              mesh.position.x = obj.x || 0
              mesh.position.z = obj.z || 0
            }
            if (!obj.autoRotate || (mesh.userData && mesh.userData.hasRotorAnimation)) {
              mesh.rotation.set(obj.rx || 0, obj.ry || 0, obj.rz || 0)
            } else {
              var axis = obj.rotateAxis || 'y'
              mesh.rotation.x = obj.rx || 0
              mesh.rotation.z = obj.rz || 0
              if (axis !== 'y') mesh.rotation.y = obj.ry || 0
            }
            mesh.scale.set(obj.sx || 1, obj.sy || 1, obj.sz || 1)
          }
          mesh.visible = obj.type === '2dComponent' ? false : (obj.visible !== false)

          if (obj.type === 'flowPipe' && mesh.userData && mesh.userData.isFlowPipe) {
            var flowDelta = 0.016
            updateFlowPipeAnimation(mesh, flowDelta)
            updateFlowPipe(mesh, obj)
            if (obj.materialOverridden === true) {
              vm.applyObjectMaterialProps(mesh, obj)
            }
          } else if (mesh.material && !mesh.userData?.isTextSprite) {
            if (Array.isArray(mesh.material)) {
              mesh.material.forEach(function(material) {
                vm.applyMaterialProps(material, obj)
              })
            } else {
              vm.applyMaterialProps(mesh.material, obj)
            }
          }

          if (isTextLabelObject(obj) && mesh.userData && mesh.userData.isTextSprite) {
            var tcStr = normTextLabelColor(obj.color)
            var tc = parseInt(tcStr.replace('#', ''), 16)
            var displayText = getTextLabelContent(obj)
            console.log('[syncMeshes] text update objId=' + obj.id + ' type=' + obj.type + ' is3D=' + mesh.userData.is3DText + ' isTrue3D=' + mesh.userData.isTrue3DText + ' color=' + obj.color + ' tc=' + tc + ' tcHex=' + tc.toString(16) + ' text=' + displayText + ' fontSize=' + (obj.fontSize || 16) + ' fontFamily=' + (buildTextLabelOptions(obj) && buildTextLabelOptions(obj).fontFamily))
            if (mesh.userData.is3DText) {
              updateText3DMesh(mesh, displayText, tc, obj.fontSize || 16, obj.textBgColor || '#00000000', buildTextLabelOptions(obj))
            } else {
              updateTextSprite(mesh, displayText, tc, obj.fontSize || 16, obj.textBgColor || '#00000000', buildTextLabelOptions(obj))
            }
            applyTextSpriteObjectScale(mesh, obj)
            updateTextPickArea(mesh, obj)
          }

          if ((obj.type === 'image3d' || obj.type === 'video3d') && mesh.userData && mesh.userData.isMediaPlane) {
            updateMediaPlane(mesh, obj, obj.type)
          }

          if (mesh.userData && mesh.userData.isLight) {
            vm.applyLightProps(mesh, obj)
          }

          if (mesh.userData && mesh.userData.isMeshGroup && mesh.traverse) {
            var hasVisualFx = mesh.userData.hasAnimation === true || mesh.userData.hasParticleAnimation === true
            var shouldApplyMaterial = obj.materialOverridden === true || (typeof obj.textureData === 'string' && obj.textureData.length > 0) || hasVisualFx
            if (shouldApplyMaterial && (obj.color || obj.opacity !== undefined || obj.wireframe !== undefined || obj.metalness !== undefined || obj.roughness !== undefined || obj.emissive || obj.textureData !== undefined)) {
              vm.applyObjectMaterialProps(mesh, obj)
            }
          }          vm.updateObjectReflection(obj, mesh)
        }
      })

      this.updateSelection(this.selectedId)
      this.keepSelectedTransformVisible()
    },

    shouldShowObjectReflection: function(obj, mesh) {
      if (!obj || !mesh) return false
      if (!obj.showShadow) return false
      if (obj.visible === false) return false
      if (obj.type === '2dComponent') return false
      if (obj.type && obj.type.indexOf('light') !== -1) return false
      if (mesh.userData && (mesh.userData.isDragShadow || mesh.userData.isObjectReflection)) return false
      return true
    },

    applyObjectShadowState: function(mesh, obj) {
      if (!mesh || !obj) return
      var enabled = !!obj.showShadow
      if (mesh.isMesh) {
        mesh.castShadow = enabled
        mesh.receiveShadow = enabled
      }
      if (mesh.traverse) {
        mesh.traverse(function(child) {
          if (!child.isMesh) return
          if (child.userData && (child.userData.isDragShadow || child.userData.isObjectReflection)) return
          child.castShadow = enabled
          child.receiveShadow = enabled
        })
      }
    },

    createReflectionMaterial: function(material, isSprite) {
      if (!material) return material
      var originalUserData = material.userData
      material.userData = this.getSafeReflectionUserData(originalUserData)
      var cloned = material.clone ? material.clone() : material
      material.userData = originalUserData
      if (cloned && cloned !== material) {
        cloned.userData = {}
      }
      cloned.transparent = true
      cloned.opacity = Math.min(material.opacity !== undefined ? material.opacity : 1, isSprite ? 0.2 : 0.24)
      cloned.depthWrite = false
      cloned.depthTest = false
      cloned.side = THREE.DoubleSide
      if (cloned.color && cloned.color.multiplyScalar) {
        cloned.color.multiplyScalar(isSprite ? 0.76 : 0.58)
      }
      if (cloned.emissive && cloned.emissive.setHex) {
        cloned.emissive.setHex(0x000000)
      }
      cloned.needsUpdate = true
      return cloned
    },

    getSafeReflectionUserData: function(userData) {
      var safe = {}
      if (!userData) return safe
      Object.keys(userData).forEach(function(key) {
        var value = userData[key]
        if (value === null || value === undefined) return
        var type = typeof value
        if (type === 'string' || type === 'number' || type === 'boolean') {
          safe[key] = value
        } else if (Array.isArray(value) && value.every(function(item) {
          var itemType = typeof item
          return item === null || itemType === 'string' || itemType === 'number' || itemType === 'boolean'
        })) {
          safe[key] = value.slice()
        }
      })
      return safe
    },

    withSafeReflectionUserData: function(source, callback) {
      var records = []
      var vm = this
      source.traverse(function(child) {
        records.push({
          target: child,
          userData: child.userData
        })
        child.userData = vm.getSafeReflectionUserData(child.userData)
      })

      try {
        return callback()
      } finally {
        records.forEach(function(record) {
          record.target.userData = record.userData
        })
      }
    },

    prepareReflectionClone: function(source, objectId) {
      if (!source) return null
      var vm = this
      var clone = this.withSafeReflectionUserData(source, function() {
        return source.clone(true)
      })
      clone.userData = Object.assign({}, clone.userData || {}, {
        id: objectId + '_reflection',
        sourceId: objectId,
        isObjectReflection: true
      })
      clone.traverse(function(child) {
        child.userData = Object.assign({}, child.userData || {}, {
          id: objectId + '_reflection',
          sourceId: objectId,
          isObjectReflection: true
        })
        child.castShadow = false
        child.receiveShadow = false
        child.renderOrder = 2
        if (child.isLight || child.type === 'ArrowHelper') {
          child.visible = false
          return
        }
        if (child.material) {
          var isSprite = child.isSprite
          if (Array.isArray(child.material)) {
            child.material = child.material.map(function(material) {
              return vm.createReflectionMaterial(material, isSprite)
            })
          } else {
            child.material = vm.createReflectionMaterial(child.material, isSprite)
          }
        }
      })
      return clone
    },

    updateObjectReflection: function(obj, mesh) {
      if (!this.__3d || !this.__3d.scene) return
      if (!this.__3d.reflectionMap) this.__3d.reflectionMap = {}
      if (!this.shouldShowObjectReflection(obj, mesh)) {
        this.removeObjectReflection(obj && obj.id)
        return
      }
      var reflection = this.__3d.reflectionMap[obj.id]
      if (!reflection) {
        reflection = this.prepareReflectionClone(mesh, obj.id)
        if (!reflection) return
        this.__3d.scene.add(reflection)
        this.__3d.reflectionMap[obj.id] = reflection
      }
      reflection.visible = mesh.visible !== false
      reflection.position.set(mesh.position.x, -mesh.position.y + 0.018, mesh.position.z)
      reflection.rotation.copy(mesh.rotation)
      reflection.scale.set(mesh.scale.x, -mesh.scale.y, mesh.scale.z)
    },

    removeObjectReflection: function(objectId) {
      if (!objectId || !this.__3d || !this.__3d.reflectionMap) return
      var reflection = this.__3d.reflectionMap[objectId]
      if (!reflection) return
      this.__3d.scene.remove(reflection)
      reflection.traverse(function(child) {
        if (!child.material) return
        if (Array.isArray(child.material)) {
          child.material.forEach(function(material) {
            if (material && material.dispose) material.dispose()
          })
        } else if (child.material.dispose) {
          child.material.dispose()
        }
      })
      delete this.__3d.reflectionMap[objectId]
    },

    updateSelection: function(id) {
      if (!this.__3d) return
      if (!id) {
        this.detachTransform()
        this.dragGuide.visible = false
        return
      }
      var obj = this.objects ? this.objects.find(function(item) { return item.id === id }) : null
      if (obj && obj.locked) {
        this.detachTransform()
        return
      }
      var mesh = this.__3d.meshMap[id]
      if (mesh) {
        if (this.__3d.transform.object === mesh) {
          this.setTransformControlVisible(true)
          this.updateSelectionFrame(mesh)
          return
        }
        this.detachTransform()
        this.attachTransform(mesh)
        return
      }
      this.detachTransform()
      this.retrySelectedMeshAttach(id)
    },

    retrySelectedMeshAttach: function(id) {
      var vm = this
      setTimeout(function() {
        if (!vm.__3d || vm.selectedId !== id) return
        var mesh = vm.__3d.meshMap && vm.__3d.meshMap[id]
        if (!mesh) return
        vm.attachTransform(mesh)
        vm.updateActiveSelectionFrame()
      }, 0)
    },

    refreshLoadedSelection: function(objectId) {
      if (!this.__3d || this.selectedId !== objectId) return
      var mesh = this.__3d.meshMap && this.__3d.meshMap[objectId]
      if (!mesh) return
      this.attachTransform(mesh)
      this.updateActiveSelectionFrame()
    },

    updateTextObject: function(obj) {
      if (!this.__3d || !obj || !obj.id) return
      var mesh = this.__3d.meshMap[obj.id]
      if (!mesh) return

      if (isTextLabelObject(obj) && mesh.userData && mesh.userData.isTextSprite) {
        var tcStr = normTextLabelColor(obj.color)
        var tc = parseInt(tcStr.replace('#', ''), 16)
        var displayText = getTextLabelContent(obj)
        if (mesh.userData.is3DText) {
          updateText3DMesh(mesh, displayText, tc, obj.fontSize || 16, obj.textBgColor || '#00000000', buildTextLabelOptions(obj))
        } else {
          updateTextSprite(mesh, displayText, tc, obj.fontSize || 16, obj.textBgColor || '#00000000', buildTextLabelOptions(obj))
        }
        applyTextSpriteObjectScale(mesh, obj)
      }
    },

    initThree: function() {
      var vm = this
      var canvas = vm.$refs.threeCanvas
      var wrapper = vm.$refs.canvasWrapper
      if (!wrapper) return
      var rect = wrapper.getBoundingClientRect()
      var w = rect.width || 800
      var h = rect.height || 600

      var scene = new THREE.Scene()
      scene.background = new THREE.Color(0xffffff)
      scene.fog = new THREE.Fog(0xffffff, 40, 100)

      var camera = new THREE.PerspectiveCamera(45, w / h, 0.1, 1000)
      camera.position.set(6, 5, 8)
      camera.lookAt(0, 0, 0)

      var renderer = new THREE.WebGLRenderer({ canvas: canvas, antialias: true })
      renderer.setSize(w, h)
      renderer.setPixelRatio(window.devicePixelRatio)
      renderer.shadowMap.enabled = true
      renderer.shadowMap.type = THREE.PCFSoftShadowMap
      renderer.outputEncoding = THREE.sRGBEncoding
      renderer.toneMapping = THREE.ACESFilmicToneMapping
      renderer.toneMappingExposure = 1.0

      var ambientLight = new THREE.AmbientLight(0xffffff, 1.0)
      ambientLight.name = 'ambientLight'
      scene.add(ambientLight)
      var dirLight = new THREE.DirectionalLight(0xffffff, 0.9)
      dirLight.name = 'directionalLight'
      dirLight.position.set(8, 12, 6)
      dirLight.castShadow = true
      dirLight.shadow.mapSize.set(2048, 2048)
      scene.add(dirLight)
      var fillLight = new THREE.DirectionalLight(0x4466aa, 0.3)
      fillLight.name = 'fillLight'
      scene.add(fillLight)

      var _auxObjects = []

      var gridHelper = new THREE.GridHelper(40, 40, 0x4488aa, 0x334455)
      gridHelper.material.opacity = 0.18
      gridHelper.material.transparent = true
      gridHelper.userData.isAux = true
      scene.add(gridHelper)
      _auxObjects.push(gridHelper)
      var axesHelper = new THREE.AxesHelper(1.5)
      axesHelper.position.set(-9.5, 0.01, -9.5)
      axesHelper.userData.isAux = true
      scene.add(axesHelper)
      _auxObjects.push(axesHelper)

      var orbit = new OrbitControls(camera, canvas)
      orbit.enableDamping = true
      orbit.enablePan = true
      orbit.dampingFactor = 0.06
      orbit.minDistance = 1
      orbit.maxDistance = 100
      orbit.mouseButtons = { LEFT: THREE.MOUSE.ROTATE, MIDDLE: THREE.MOUSE.DOLLY, RIGHT: THREE.MOUSE.PAN }

      var transform = new TransformControls(camera, canvas)
      transform.setTranslationSnap(null)
      transform.setRotationSnap(THREE.MathUtils.degToRad(15))
      transform.setScaleSnap(0.1)
      transform.visible = false
      transform.userData.isAux = true
      scene.add(transform)
      vm.applyTransformControlContrast(transform)
      _auxObjects.push(transform)

      // Tech center ring (翠鸟 style: cyan glow ring + translucent core cube)
      var techRingGroup = new THREE.Group()
      techRingGroup.visible = false
      techRingGroup.userData.isAux = true

      // Outer cyan glow ring
      var ringGeo = new THREE.TorusGeometry(0.22, 0.018, 12, 40)
      var ringMat = new THREE.MeshBasicMaterial({
        color: 0x00f0ff,
        transparent: true,
        opacity: 0.9,
        depthTest: false,
        depthWrite: false
      })
      var ringMesh = new THREE.Mesh(ringGeo, ringMat)
      ringMesh.userData.isAux = true
      ringMesh.renderOrder = 1002
      techRingGroup.add(ringMesh)

      // Inner accent ring (smaller, magenta)
      var ring2Geo = new THREE.TorusGeometry(0.14, 0.01, 8, 32)
      var ring2Mat = new THREE.MeshBasicMaterial({
        color: 0xff00cc,
        transparent: true,
        opacity: 0.6,
        depthTest: false,
        depthWrite: false
      })
      var ring2Mesh = new THREE.Mesh(ring2Geo, ring2Mat)
      ring2Mesh.userData.isAux = true
      ring2Mesh.renderOrder = 1002
      techRingGroup.add(ring2Mesh)

      scene.add(techRingGroup)
      _auxObjects.push(techRingGroup)

      var clock = new THREE.Clock()
      var raycaster = new THREE.Raycaster()
      var mouse = new THREE.Vector2()
      var meshMap = {}

      vm.__3d = { scene: scene, camera: camera, renderer: renderer, orbit: orbit, transform: transform, gridHelper: gridHelper, axesHelper: axesHelper, clock: clock, raycaster: raycaster, mouse: mouse, meshMap: meshMap, _auxObjects: _auxObjects, techRing: techRingGroup }

      transform.addEventListener('dragging-changed', function(e) {
        if (!vm.__3d) return
        orbit.enabled = !e.value
        if (e.value && transform.object) {
          vm.setTransformControlVisible(true)
          vm.showTransformGuides(transform.object)
          vm.updateTransformGuides(transform.object)
        }
        if (!e.value && transform.object) {
          var mesh = transform.object
          var id = mesh.userData.id
          vm._ignoreNextCanvasClickAfterTransform = true
          var sourceObj = vm.objects && vm.objects.find(function(item) { return item && item.id === id })
          vm.updateTransformGuides(mesh)
          vm.hideTransformGuides()
          vm.clearTransformAutoHide()
          vm.setTransformControlVisible(true)
          var sx = parseFloat(mesh.scale.x.toFixed(4))
          var sy = parseFloat(mesh.scale.y.toFixed(4))
          var sz = parseFloat(mesh.scale.z.toFixed(4))
          if (transform.mode !== 'scale' && sourceObj) {
            sx = sourceObj.sx !== undefined ? Number(sourceObj.sx) : 1
            sy = sourceObj.sy !== undefined ? Number(sourceObj.sy) : 1
            sz = sourceObj.sz !== undefined ? Number(sourceObj.sz) : 1
          } else if (mesh.userData && mesh.userData.isTextSprite && !mesh.userData.is3DText && mesh.userData.baseTextScale) {
            var base = mesh.userData.baseTextScale
            sx = base.x ? parseFloat((mesh.scale.x / base.x).toFixed(4)) : 1
            sy = base.y ? parseFloat((mesh.scale.y / base.y).toFixed(4)) : 1
            sz = base.z ? parseFloat((mesh.scale.z / base.z).toFixed(4)) : 1
          }
          vm.$emit('object-transformed', {
            objectId: id,
            transform: {
              x: parseFloat(mesh.position.x.toFixed(4)),
              y: parseFloat(mesh.position.y.toFixed(4)),
              z: parseFloat(mesh.position.z.toFixed(4)),
              rx: parseFloat(mesh.rotation.x.toFixed(4)),
              ry: parseFloat(mesh.rotation.y.toFixed(4)),
              rz: parseFloat(mesh.rotation.z.toFixed(4)),
              sx: sx,
              sy: sy,
              sz: sz
            }
          })
          vm.setTransformControlVisible(true)
          vm.$nextTick(function() {
            vm.keepSelectedTransformVisible()
          })
        }
      })
      transform.addEventListener('objectChange', function() {
        if (!vm.__3d || !transform.object) return
        vm.applyTransformSnapAndAlignment(transform.object)
        vm.updateTransformGuides(transform.object)
      })

      vm._handleCanvasClick = function(e) { vm.onCanvasClick(e) }
      vm._handleCanvasContextMenu = function(e) { vm.onCanvasContextMenu(e) }
      vm._handleCanvasMouseMove = function(e) { vm.onCanvasMouseMove(e) }
      vm._handleCanvasDoubleClick = function(e) { vm.onCanvasDoubleClick(e) }
      vm._handleCanvasMouseDown = function() { if (vm.__3d) vm.__3d._cameraAnim = null }
      vm._handleWindowResize = function() { vm.onResize() }
      vm._handleWindowMouseMove = function(e) {
        vm.onSelectionFrameResizeMove(e)
        vm.onOverlayMouseMove(e)
      }
      vm._handleWindowMouseUp = function() {
        vm.onSelectionFrameResizeUp()
        vm.onOverlayMouseUp()
      }
      vm._handleWindowKeyDown = function(e) { vm.onWindowKeyDown(e) }
      vm._handleWindowKeyUp = function(e) { vm.onWindowKeyUp(e) }

      canvas.addEventListener('click', vm._handleCanvasClick)
      canvas.addEventListener('contextmenu', vm._handleCanvasContextMenu)
      canvas.addEventListener('mousemove', vm._handleCanvasMouseMove)
      canvas.addEventListener('dblclick', vm._handleCanvasDoubleClick)
      canvas.addEventListener('mousedown', vm._handleCanvasMouseDown)

      window.addEventListener('resize', vm._handleWindowResize)
      window.addEventListener('mousemove', vm._handleWindowMouseMove)
      window.addEventListener('mouseup', vm._handleWindowMouseUp)
      window.addEventListener('keydown', vm._handleWindowKeyDown)
      window.addEventListener('keyup', vm._handleWindowKeyUp)

      vm.$emit('initialized')
    },

    isEditableTarget: function(target) {
      if (!target || !target.tagName) return false
      var tag = target.tagName.toUpperCase()
      return tag === 'INPUT' || tag === 'TEXTAREA' || tag === 'SELECT' || target.isContentEditable
    },

    setSpacePanActive: function(active) {
      if (!this.__3d || !this.__3d.orbit) return
      this._spacePanActive = !!active
      this.__3d.orbit.mouseButtons = this._spacePanActive
        ? { LEFT: THREE.MOUSE.PAN, MIDDLE: THREE.MOUSE.DOLLY, RIGHT: THREE.MOUSE.ROTATE }
        : { LEFT: THREE.MOUSE.ROTATE, MIDDLE: THREE.MOUSE.DOLLY, RIGHT: THREE.MOUSE.PAN }
      if (this.$refs.threeCanvas) {
        this.$refs.threeCanvas.style.cursor = this._spacePanActive ? 'grab' : 'default'
      }
    },

    onWindowKeyDown: function(e) {
      if (e.code !== 'Space' || this.isEditableTarget(e.target) || this.isDrawingPipe) return
      e.preventDefault()
      if (!this._spacePanActive) {
        this.setSpacePanActive(true)
      }
    },

    onWindowKeyUp: function(e) {
      if (e.code !== 'Space') return
      this.setSpacePanActive(false)
    },

    startRenderLoop: function() {
      var vm = this
      var render = function() {
        vm._animId = requestAnimationFrame(render)
        if (!vm.__3d) return
        var delta = vm.__3d.clock.getDelta()
        var elapsed = vm.__3d.clock.getElapsedTime()
        vm.__3d.orbit.update()

        var anim = vm.__3d._cameraAnim
        if (anim) {
          var t = Math.min((performance.now() - anim.startTime) / anim.duration, 1)
          var ease = t < 0.5 ? 2 * t * t : -1 + (4 - 2 * t) * t
          vm.__3d.orbit.target.lerpVectors(anim.startTarget, anim.endTarget, ease)
          vm.__3d.camera.position.lerpVectors(anim.startCam, anim.endCam, ease)
          if (t >= 1) {
            vm.__3d._cameraAnim = null
          }
        }

        if (vm.pipePreviewMesh && vm.pipePreviewMesh.userData && vm.pipePreviewMesh.userData.isFlowPipe) {
          updateFlowPipeAnimation(vm.pipePreviewMesh, delta)
        }

        // 动画
        var objs = vm.objects
        var meshMap = vm.__3d.meshMap
        for (var i = 0; i < objs.length; i++) {
          var obj = objs[i]
          var mesh = meshMap[obj.id]
          if (!mesh || !obj) continue
          // 自定义更新函数（数字孪生特效动画）
          if (mesh.userData && mesh.userData.update) mesh.userData.update(elapsed)
          var isTransforming = vm.__3d.transform.dragging && vm.__3d.transform.object === mesh
          if (mesh.userData && mesh.userData.isTextSprite && !mesh.userData.is3DText && obj.labelFaceCamera !== false) {
            if (vm.__3d.transform.object !== mesh) {
              mesh.quaternion.copy(vm.__3d.camera.quaternion)
            }
          }
          if (isTransforming) continue

          if (obj.autoRotate) {
            var speed = obj.rotateSpeed !== undefined ? obj.rotateSpeed : 1
            var axis = obj.rotateAxis || 'y'
            if (mesh.userData && mesh.userData.hasRotorAnimation && Array.isArray(mesh.userData.rotors)) {
              mesh.userData.rotors.forEach(function(rotor) {
                rotor.rotation.y -= speed * delta * 8
              })
            } else {
              mesh.rotation[axis] += speed * delta
            }
          }

          vm.updateGLTFAnimation(mesh, obj, delta)

          if (obj.floatAnim) {
            var floatRange = obj.floatRange !== undefined ? obj.floatRange : 0.15
            var floatSpeed = obj.floatSpeed !== undefined ? obj.floatSpeed : 2
            var floatOffset = Math.sin(elapsed * floatSpeed + i * 1.5) * floatRange
            var baseY = obj._baseY !== undefined ? obj._baseY : obj.y
            if (obj._baseY === undefined) obj._baseY = obj.y
            mesh.position.y = baseY + floatOffset
          } else {
            if (obj._baseY !== undefined) {
              mesh.position.y = obj._baseY
              delete obj._baseY
            }
          }

          // 闪烁动画，支持 Mesh 和 MeshGroup
          if (obj.blink) {
            var blinkSpeed = obj.blinkSpeed !== undefined ? obj.blinkSpeed : 6
            var blinkMin = obj.blinkMin !== undefined ? obj.blinkMin : 0.2
            var blinkVal = (Math.sin(elapsed * blinkSpeed) + 1) * 0.5
            var targetOpacity = blinkMin + blinkVal * (1 - blinkMin)

            if (mesh.material) {
              // 普通 Mesh：直接操作顶层 material
              var mats = Array.isArray(mesh.material) ? mesh.material : [mesh.material]
              mats.forEach(function(mat) {
                mat.opacity = targetOpacity
                mat.transparent = true
                mat.needsUpdate = true
              })
            }
            // MeshGroup：遍历所有子Mesh
            if (mesh.userData && mesh.userData.isMeshGroup && mesh.traverse) {
              mesh.traverse(function(child) {
                if (child.isMesh && child.material) {
                  var childMats = Array.isArray(child.material) ? child.material : [child.material]
                  childMats.forEach(function(mat) {
                    // 跳过 ShaderMaterial，视觉特效组件自己管理材质
                    if (mat.isShaderMaterial) return
                    mat.opacity = targetOpacity
                    mat.transparent = true
                    mat.needsUpdate = true
                  })
                }
              })
            }
          }

          // 流动管道动画
          if (obj.type === 'flowPipe' && mesh.userData && mesh.userData.isFlowPipe) {
            var flowCondition = obj.flowAnim !== false
            if (obj.animate && obj.animate.condition && obj.animate.condition.isBandDevice) {
              flowCondition = flowCondition && vm.checkAnimateCondition(obj)
            }
            updateFlowPipeAnimation(mesh, delta, flowCondition)
          }

          // Shader 特效组件：更新 uniforms.time
          if (mesh.userData && mesh.userData.hasAnimation) {
            var frameCount = vm._fpsFrames
            var foundTime = false
            mesh.traverse(function(child) {
              if (child.isMesh && child.material && child.material.uniforms && child.material.uniforms.time !== undefined) {
                child.material.uniforms.time.value = elapsed
                foundTime = true
              }
            })
            // 每秒打印一次调试日志
            if (frameCount === 0 && (obj.type === 'flowRibbon' || obj.type === 'energyWall' || obj.type === 'pulseColumn')) {
              console.log('[动画]', obj.type, 'id=' + obj.id, 'hasAnimation=' + mesh.userData.hasAnimation, 'time更新=' + foundTime, 'elapsed=' + elapsed.toFixed(2))
            }
          }

          // 粒子特效动画（翠鸟风格新增）
          if (mesh.userData && mesh.userData.hasAnimation) {
            var fxType = obj.type
            var speed = obj.animateSpeed || 1
            // 火焰/烟雾/火花：粒子上浮
            if (fxType === 'fireFX' || fxType === 'smokeFX' || fxType === 'sparkFX') {
              var pts = mesh.children[0]
              if (pts && pts.isPoints && pts.geometry.attributes.position) {
                var posArr = pts.geometry.attributes.position.array
                var count = posArr.length / 3
                for (var pi = 0; pi < count; pi++) {
                  posArr[pi * 3 + 1] += delta * speed * (fxType === 'smokeFX' ? 0.25 : 0.6)
                  if (fxType === 'smokeFX') {
                    posArr[pi * 3] += (Math.sin(elapsed * 1.5 + pi) * delta * 0.12)
                    posArr[pi * 3 + 2] += (Math.cos(elapsed * 1.3 + pi) * delta * 0.12)
                  }
                  if (posArr[pi * 3 + 1] > (fxType === 'sparkFX' ? 0.7 : 0.9)) {
                    posArr[pi * 3 + 1] = fxType === 'sparkFX' ? 0 : 0.04
                    posArr[pi * 3] = (Math.random() - 0.5) * (fxType === 'fireFX' ? 0.2 : 0.3)
                    posArr[pi * 3 + 2] = (Math.random() - 0.5) * (fxType === 'fireFX' ? 0.2 : 0.3)
                  }
                }
                pts.geometry.attributes.position.needsUpdate = true
              }
            }
            // 灰尘：微动
            if (fxType === 'dustFX') {
              var dPts = mesh.children[0]
              if (dPts && dPts.isPoints) { dPts.rotation.y += delta * 0.15 }
            }
            // 萤火虫：脉冲明灭 + 微小漂浮
            if (fxType === 'fireflyFX') {
              var ffPts = mesh.children[0]
              if (ffPts && ffPts.isPoints && mesh.userData.fireflyPhases) {
                var phases = mesh.userData.fireflyPhases
                var spds = mesh.userData.fireflySpeeds
                ffPts.material.opacity = 0.3 + 0.6 * (0.5 + 0.5 * Math.sin(elapsed * 2.0))
                var ffPos = ffPts.geometry.attributes.position.array
                var ffCount = ffPos.length / 3
                for (var fi = 0; fi < ffCount; fi++) {
                  ffPos[fi * 3 + 1] += Math.sin(elapsed * spds[fi] + phases[fi]) * delta * 0.08
                }
                ffPts.geometry.attributes.position.needsUpdate = true
              }
            }
            // 雪花：飘落 + 水平摆动
            if (fxType === 'snowFX') {
              var sPts = mesh.children[0]
              if (sPts && sPts.isPoints && sPts.geometry.attributes.position) {
                var sArr = sPts.geometry.attributes.position.array
                var sCnt = sArr.length / 3
                for (var si = 0; si < sCnt; si++) {
                  sArr[si * 3 + 1] -= delta * speed * 0.2
                  sArr[si * 3] += Math.sin(elapsed * 2.5 + si) * delta * 0.1
                  if (sArr[si * 3 + 1] < 0) {
                    sArr[si * 3 + 1] = 1.2 + Math.random() * 0.3
                    sArr[si * 3] = (Math.random() - 0.5) * 1.6
                    sArr[si * 3 + 2] = (Math.random() - 0.5) * 1.6
                  }
                }
                sPts.geometry.attributes.position.needsUpdate = true
              }
            }
            // 雨滴：下落
            if (fxType === 'rainFX') {
              var rLines = mesh.children[0]
              if (rLines && rLines.isLineSegments && rLines.geometry.attributes.position) {
                var rArr = rLines.geometry.attributes.position.array
                var rCnt = rArr.length / 3
                for (var ri = 0; ri < rCnt; ri += 2) {
                  rArr[ri * 3 + 1] -= delta * speed * 1.2
                  rArr[ri * 3 + 4] -= delta * speed * 1.2
                  if (rArr[ri * 3 + 1] < -0.1) {
                    rArr[ri * 3 + 1] = 1.3 + Math.random() * 0.3
                    rArr[ri * 3 + 4] = rArr[ri * 3 + 1] - 0.12
                  }
                }
                rLines.geometry.attributes.position.needsUpdate = true
              }
            }
            // 气泡：上浮晃动
            if (fxType === 'bubbleFX') {
              mesh.children.forEach(function(b, bi) {
                if (b.isMesh && b.name && b.name.indexOf('bubble_') === 0) {
                  b.position.y += delta * speed * 0.18
                  b.position.x += Math.sin(elapsed * 1.8 + bi) * delta * 0.06
                  if (b.position.y > 1.0) { b.position.y = 0.05; b.position.x = (Math.random()-0.5)*1.2; b.position.z = (Math.random()-0.5)*1.2 }
                }
              })
            }
            // 落叶：旋转下落
            if (fxType === 'leafFX') {
              mesh.children.forEach(function(lf, li) {
                if (lf.isMesh && lf.name && lf.name.indexOf('leaf_') === 0) {
                  lf.position.y -= delta * speed * 0.15
                  lf.position.x += Math.sin(elapsed * 1.2 + li * 1.7) * delta * 0.1
                  lf.position.z += Math.cos(elapsed * 1.0 + li * 1.3) * delta * 0.08
                  lf.rotation.z += delta * 1.5
                  lf.rotation.x += delta * 0.7
                  if (lf.position.y < -0.1) {
                    lf.position.y = 1.1 + Math.random() * 0.3
                    lf.position.x = (Math.random()-0.5)*1.4
                    lf.position.z = (Math.random()-0.5)*1.4
                  }
                }
              })
            }
            // 爆炸：环膨胀 + 粒子散射
            if (fxType === 'explosionFX') {
              var scale = 0.6 + 0.4 * Math.abs(Math.sin(elapsed * 3.0))
              mesh.children.forEach(function(c) {
                if (c.name === 'explosionRing') { c.scale.setScalar(0.6 + Math.sin(elapsed * 5.0) * 0.5); c.material.opacity = 0.3 + scale * 0.5 }
                if (c.name === 'explosionInnerRing') { c.scale.setScalar(0.5 + Math.sin(elapsed * 4.5 + 0.5) * 0.4) }
              })
            }
            // 水花：喷射
            if (fxType === 'splashFX') {
              var spPts = mesh.children[0]
              if (spPts && spPts.isPoints && spPts.geometry.attributes.position) {
                var spArr = spPts.geometry.attributes.position.array
                var spCnt = spArr.length / 3
                for (var spi = 0; spi < spCnt; spi++) {
                  spArr[spi * 3 + 1] -= delta * speed * 0.15
                  if (spArr[spi * 3 + 1] < 0) { spArr[spi * 3 + 1] = Math.random() * 0.6; spArr[spi * 3] = (Math.random()-0.5)*0.45; spArr[spi * 3 + 2] = (Math.random()-0.5)*0.45 }
                }
                spPts.geometry.attributes.position.needsUpdate = true
              }
            }
            // 星星：闪烁 + 旋转
            if (fxType === 'starFX') {
              var stPts = mesh.children[0]
              if (stPts && stPts.isPoints) {
                stPts.material.opacity = 0.5 + 0.4 * (0.5 + 0.5 * Math.sin(elapsed * 3.5))
                stPts.rotation.y += delta * 0.1
              }
            }
            // 魔法光晕：粒子缓慢浮动 + 明灭
            if (fxType === 'magicGlow') {
              var mgPts = mesh.children.find(function(c){return c.name==='magicGlowParticles'})
              if (mgPts && mgPts.geometry) {
                var mgPos = mgPts.geometry.attributes.position
                var mgBase = mesh.userData.magicBasePos
                var mgPhase = mesh.userData.magicPhases
                if (mgPos && mgBase && mgPhase) {
                  for (var mi = 0; mi < mgPos.count; mi++) {
                    var mphase = mgPhase[mi]
                    var my = mgBase[mi*3+1] + Math.sin(elapsed * 0.4 + mphase) * 0.08
                    var mx = mgBase[mi*3]   + Math.sin(elapsed * 0.25 + mphase * 1.3) * 0.03
                    var mz = mgBase[mi*3+2] + Math.cos(elapsed * 0.3 + mphase * 0.7) * 0.03
                    mgPos.setXYZ(mi, mx, my, mz)
                  }
                  mgPos.needsUpdate = true
                }
                mgPts.material.opacity = 0.45 + 0.35 * Math.sin(elapsed * 0.6)
              }
            }
          }

          vm.updateObjectReflection(obj, mesh)
        }

        vm.updateActiveSelectionFrame()
        vm.updateSelectedCoordinateTip()

        // Tech ring animation (翠鸟 style)
        if (vm.__3d.techRing && vm.__3d.techRing.visible && vm.__3d.transform.object) {
          var tr = vm.__3d.techRing
          var target = vm.__3d.transform.object
          target.getWorldPosition(tr.position)

          // Keep constant screen size
          var camPos = vm.__3d.camera.position
          var dist = tr.position.distanceTo(camPos)
          var fov = vm.__3d.camera.fov * (Math.PI / 180)
          var vh = 2 * Math.tan(fov / 2) * dist
          var base = vh * 0.012
          var pulse = 1 + 0.08 * Math.sin(elapsed * 4)
          tr.scale.setScalar(base * pulse)

          // Opacity pulse for rings
          if (tr.children[0] && tr.children[0].material) {
            tr.children[0].material.opacity = 0.7 + 0.25 * Math.sin(elapsed * 3)
          }
          if (tr.children[1] && tr.children[1].material) {
            tr.children[1].material.opacity = 0.4 + 0.3 * Math.sin(elapsed * 4 + 1.2)
          }
        }

        // 简化渲染：暂时禁用 Bloom 两遍渲染，确保场景正常显示
        if (vm.__3d && vm.__3d.renderer) {
          vm.__3d.renderer.render(vm.__3d.scene, vm.__3d.camera)
        }
        vm.update2DOverlays()
        vm._fpsFrames++
        var now = performance.now()
        if (now - vm._fpsLast >= 1000) {
          vm.fps = vm._fpsFrames
          vm._fpsFrames = 0
          vm._fpsLast = now
        }
      }
      render()
    },

    onCanvasClick: function(e) {
      if (e.button !== 0) return

      if (this._ignoreNextCanvasClickAfterTransform) {
        this._ignoreNextCanvasClickAfterTransform = false
        this.keepSelectedTransformVisible()
        return
      }

      if (this.isDrawingPipe) {
        this.addPipePoint(e)
        return
      }

      var vm = this
      var rect = (vm.$refs.canvasWrapper || vm.$refs.threeCanvas).getBoundingClientRect()
      vm.__3d.mouse.x = ((e.clientX - rect.left) / rect.width) * 2 - 1
      vm.__3d.mouse.y = -((e.clientY - rect.top) / rect.height) * 2 + 1
      vm.__3d.raycaster.setFromCamera(vm.__3d.mouse, vm.__3d.camera)

      var allMeshes = []
      Object.values(vm.__3d.meshMap).forEach(function(m) {
        if (m && m.traverse) {
          m.traverse(function(c) { if (c.isMesh || c.isPoints || c.isLineSegments) allMeshes.push(c) })
        } else if (m) {
          allMeshes.push(m)
        }
      })

      var hits = vm.__3d.raycaster.intersectObjects(allMeshes, false)
      if (hits.length > 0) {
        var hit = hits[0].object
        var id = hit.userData.id
        if (!id) {
          var parent = hit.parent
          while (parent) {
            if (parent.userData && parent.userData.id) { id = parent.userData.id; break }
            parent = parent.parent
          }
        }
        if (id) {
          vm.$emit('object-selected', id)
          if (id === vm.selectedId && vm.mode === 'select') {
            var selectedMesh = vm.__3d.meshMap && vm.__3d.meshMap[id]
            if (selectedMesh) {
              vm.attachTransform(selectedMesh)
              vm.setTransformControlVisible(true)
          vm.clearTransformAutoHide()
            }
          }
        }
        // 点击空白处不再取消选中，坐标始终保持可见
      }
    },

    onCanvasDoubleClick: function(e) {
      if (this.isDrawingPipe) {
        e.preventDefault()
        e.stopPropagation()
        if (this.pipePoints.length > 2) {
          this.pipePoints.pop()
        } else if (this.pipePoints.length > 1) {
          var last = this.pipePoints[this.pipePoints.length - 1]
          var prev = this.pipePoints[this.pipePoints.length - 2]
          var dx = last.x - prev.x
          var dy = last.y - prev.y
          var dz = last.z - prev.z
          if (Math.abs(dx) < 0.001 && Math.abs(dy) < 0.001 && Math.abs(dz) < 0.001) {
            this.pipePoints.pop()
          }
        }
        this.finishDrawingPipe()
        return
      }

      if (!this.__3d) return
      var vm = this
      var rect = vm.$refs.threeCanvas.getBoundingClientRect()
      vm.__3d.mouse.x = ((e.clientX - rect.left) / rect.width) * 2 - 1
      vm.__3d.mouse.y = -((e.clientY - rect.top) / rect.height) * 2 + 1
      vm.__3d.raycaster.setFromCamera(vm.__3d.mouse, vm.__3d.camera)

      var allMeshes = []
      Object.values(vm.__3d.meshMap).forEach(function(m) {
        if (m && m.traverse) {
          m.traverse(function(c) { if (c.isMesh || c.isPoints || c.isLineSegments) allMeshes.push(c) })
        } else if (m) {
          allMeshes.push(m)
        }
      })

      var hits = vm.__3d.raycaster.intersectObjects(allMeshes, false)
      if (hits.length > 0) {
        var hit = hits[0].object
        var id = hit.userData.id
        if (!id) {
          var parent = hit.parent
          while (parent) {
            if (parent.userData && parent.userData.id) { id = parent.userData.id; break }
            parent = parent.parent
          }
        }
        if (id) {
          var mesh = vm.__3d.meshMap[id]
          if (mesh) {
            var box = new THREE.Box3().setFromObject(mesh)
            var center = box.getCenter(new THREE.Vector3())
            var size = box.getSize(new THREE.Vector3())
            var maxDim = Math.max(size.x, size.y, size.z, 0.5)
            var dist = Math.max(maxDim * 4, 3)

            var startTarget = vm.__3d.orbit.target.clone()
            var endTarget = center.clone()
            var startCam = vm.__3d.camera.position.clone()
            var dir = startCam.clone().sub(startTarget).normalize()
            if (dir.length() < 0.01) dir.set(0.5, 0.4, 0.7).normalize()
            var endCam = endTarget.clone().add(dir.multiplyScalar(dist))

            vm.__3d._cameraAnim = {
              startTime: performance.now(),
              duration: 600,
              startTarget: startTarget,
              endTarget: endTarget,
              startCam: startCam,
              endCam: endCam
            }
          }
        }
      }
    },

    onCanvasContextMenu: function(e) {
      e.preventDefault()

      var _t = this
      var menuItems = []

      if (this.selectedId) {
        var obj = this.selectedObj
        if (!obj) return

        menuItems = [
          {
            label: '删除',
            icon: 'fas fa-trash',
            onClick: function() {
              _t.$emit('delete-selected')
            }
          },
          {
            label: '复制',
            icon: 'fas fa-copy',
            onClick: function() {
              _t.$emit('duplicate', _t.selectedId)
            }
          },
          {
            label: '上移',
            icon: 'fas fa-arrow-up',
            onClick: function() {
              _t.$emit('move-up', _t.selectedId)
            }
          },
          {
            label: '下移',
            icon: 'fas fa-arrow-down',
            onClick: function() {
              _t.$emit('move-down', _t.selectedId)
            }
          },
          {
            label: '聚焦',
            icon: 'fas fa-crosshairs',
            onClick: function() {
              _t.$emit('focus-selected')
            }
          }
        ]

        if (obj.locked) {
          menuItems.push({
            label: '解锁',
            icon: 'fas fa-unlock',
            onClick: function() {
              _t.$emit('toggle-lock')
            }
          })
        } else {
          menuItems.push({
            label: '锁定',
            icon: 'fas fa-lock',
            onClick: function() {
              _t.$emit('toggle-lock')
            }
          })
        }
      }

      menuItems.push({
        label: '加载模型',
        icon: 'fas fa-cube',
        onClick: function() {
          _t.$emit('load-model')
        }
      })

      menuItems.push({
        label: '画管道',
        icon: 'fas fa-grip-lines',
        onClick: function() {
          _t.startDrawingPipe()
        }
      })

      this.$contextmenu({
        items: menuItems,
        event: e,
        zIndex: 10000
      })
    },

    startDrawingPipe: function() {
      console.log('start drawing pipe')
      this.isDrawingPipe = true
      this.pipePoints = []
      this.pipePreviewPoints = []
      if (this.__3d && this.__3d.orbit) {
        this.__3d.orbit.enabled = false
      }
      this.$emit('drawing-pipe-start')
    },

    addPipePoint: function(e) {
      var vm = this
      var point = vm.getPipeDrawingPoint(e)
      if (point) {
        point = vm.snapPipePoint(point)
        vm.pipePoints.push(point)
        // 添加新点后立即更新实时预览
        vm.updateLivePipePreview(e)
      } else {
        console.log('cannot get pipe point')
      }
    },

    updateRaycasterFromMouseEvent: function(e) {
      var vm = this
      if (!vm.__3d || !vm.$refs.threeCanvas) return false
      var rect = vm.$refs.threeCanvas.getBoundingClientRect()
      vm.__3d.mouse.x = ((e.clientX - rect.left) / rect.width) * 2 - 1
      vm.__3d.mouse.y = -((e.clientY - rect.top) / rect.height) * 2 + 1
      vm.__3d.raycaster.setFromCamera(vm.__3d.mouse, vm.__3d.camera)
      return true
    },

    getScenePickMeshes: function() {
      var vm = this
      var allMeshes = []
      if (!vm.__3d || !vm.__3d.meshMap) return allMeshes
      Object.values(vm.__3d.meshMap).forEach(function(m) {
        if (!m || m.visible === false) return
        if (m.userData && m.userData.isMeshGroup) {
          m.traverse(function(c) {
            if ((c.isMesh || c.isPoints || c.isLineSegments) && c.visible !== false) allMeshes.push(c)
          })
        } else if (m.isMesh || m.isPoints || m.isLineSegments) {
          allMeshes.push(m)
        } else if (m.traverse) {
          m.traverse(function(c) {
            if ((c.isMesh || c.isPoints || c.isLineSegments) && c.visible !== false) allMeshes.push(c)
          })
        }
      })

      return allMeshes
    },

    getPipeDrawingPoint: function(e) {
      var vm = this
      if (!vm.updateRaycasterFromMouseEvent(e)) return null

      var allMeshes = vm.getScenePickMeshes()
      var hits = vm.__3d.raycaster.intersectObjects(allMeshes, false)
      if (hits.length > 0 && hits[0].point) {
        return {
          x: hits[0].point.x,
          y: hits[0].point.y,
          z: hits[0].point.z
        }
      }

      if (vm.pipePoints && vm.pipePoints.length > 0) {
        var axisPoint = vm.getPipeAxisPointFromRay(vm.pipePoints[vm.pipePoints.length - 1])
        if (axisPoint) return axisPoint
      }

      return vm.getPipeGroundPoint()
    },

    getPipeGroundPoint: function() {
      var vm = this
      if (!vm.__3d) return null
      var plane = new THREE.Plane(new THREE.Vector3(0, 1, 0), 0)
      var intersection = new THREE.Vector3()
      var hasIntersection = vm.__3d.raycaster.ray.intersectPlane(plane, intersection)
      if (!hasIntersection || isNaN(intersection.x)) return null
      return {
        x: intersection.x,
        y: intersection.y,
        z: intersection.z
      }
    },

    getPipeAxisPointFromRay: function(lastPoint) {
      var vm = this
      if (!vm.__3d || !lastPoint) return null
      var ray = vm.__3d.raycaster.ray
      var origin = ray.origin
      var rayDir = ray.direction.clone().normalize()
      var base = new THREE.Vector3(lastPoint.x, lastPoint.y, lastPoint.z)
      var axes = [
        new THREE.Vector3(1, 0, 0),
        new THREE.Vector3(0, 1, 0),
        new THREE.Vector3(0, 0, 1)
      ]
      var best = null

      for (var i = 0; i < axes.length; i++) {
        var axis = axes[i]
        var w0 = origin.clone().sub(base)
        var b = rayDir.dot(axis)
        var d = rayDir.dot(w0)
        var e = axis.dot(w0)
        var denom = 1 - b * b
        if (Math.abs(denom) < 0.000001) continue

        var rayT = (b * e - d) / denom
        if (rayT < 0) continue
        var axisT = e + b * rayT
        var rayPoint = origin.clone().add(rayDir.clone().multiplyScalar(rayT))
        var axisPoint = base.clone().add(axis.clone().multiplyScalar(axisT))
        var distance = rayPoint.distanceTo(axisPoint)
        if (!best || distance < best.distance) {
          best = { point: axisPoint, distance: distance }
        }
      }

      if (!best) return null
      return {
        x: best.point.x,
        y: best.point.y,
        z: best.point.z
      }
    },

    snapPipePoint: function(point) {
      if (!point) return null
      if (!this.pipePoints || this.pipePoints.length === 0) return point
      var last = this.pipePoints[this.pipePoints.length - 1]
      if (!last) return point
      var dx = Math.abs(point.x - last.x)
      var dy = Math.abs(point.y - last.y)
      var dz = Math.abs(point.z - last.z)
      if (dy >= dx && dy >= dz) {
        return { x: last.x, y: point.y, z: last.z }
      }
      if (dx >= dz) {
        return { x: point.x, y: last.y, z: last.z }
      }
      return { x: last.x, y: last.y, z: point.z }
    },

    createRoundedPipeCurve: function(pipePoints, pipeRadius) {
      var curve = new THREE.CurvePath()
      if (!pipePoints || pipePoints.length < 2) return curve
      pipePoints = pipePoints.filter(function(point, index) {
        if (!point) return false
        if (index === 0) return true
        var prev = pipePoints[index - 1]
        return prev && point.distanceTo(prev) > 0.0001
      })
      if (pipePoints.length < 2) return curve
      var current = pipePoints[0].clone()
      for (var i = 1; i < pipePoints.length - 1; i++) {
        var prev = pipePoints[i - 1]
        var corner = pipePoints[i]
        var next = pipePoints[i + 1]
        if (!prev || !corner || !next) continue
        var prevLen = corner.distanceTo(prev)
        var nextLen = corner.distanceTo(next)
        if (prevLen <= 0.0001 || nextLen <= 0.0001) continue
        var bendRadius = Math.min(Math.max(pipeRadius * 3.2, 0.18), prevLen * 0.35, nextLen * 0.35)
        var before = corner.clone().add(prev.clone().sub(corner).normalize().multiplyScalar(bendRadius))
        var after = corner.clone().add(next.clone().sub(corner).normalize().multiplyScalar(bendRadius))
        if (current.distanceTo(before) > 0.0001) {
          curve.add(new THREE.LineCurve3(current, before))
        }
        curve.add(new THREE.QuadraticBezierCurve3(before, corner, after))
        current = after
      }
      var last = pipePoints[pipePoints.length - 1]
      if (current.distanceTo(last) > 0.0001) {
        curve.add(new THREE.LineCurve3(current, last))
      }
      return curve
    },

    updateLivePipePreview: function(e) {
      var vm = this
      if (!vm.isDrawingPipe || vm.pipePoints.length === 0) return

      // 移除旧的预览管道
      if (vm.pipePreviewMesh) {
        vm.__3d.scene.remove(vm.pipePreviewMesh)
        vm.pipePreviewMesh.traverse(function(child) {
          if (child.geometry) child.geometry.dispose()
          if (child.material) child.material.dispose()
        })
        vm.pipePreviewMesh = null
      }

      // Prefer a real scene hit so pipe points can land on ports at different heights.
      var drawingPoint = vm.getPipeDrawingPoint(e)
      if (!drawingPoint) return

      var mousePoint = vm.snapPipePoint(drawingPoint)
      if (!mousePoint) return

      // 构建预览点：已有定点 + 当前鼠标位置
      var previewPoints = vm.pipePoints.concat([mousePoint]).filter(function(point, index, list) {
        if (!point || typeof point.x !== 'number' || typeof point.y !== 'number' || typeof point.z !== 'number') return false
        if (index === 0) return true
        var prev = list[index - 1]
        if (!prev) return true
        return Math.abs(point.x - prev.x) > 0.0001 || Math.abs(point.y - prev.y) > 0.0001 || Math.abs(point.z - prev.z) > 0.0001
      })
      if (previewPoints.length < 2) return

      // Use the same point format as createFlowPipe for the live preview.
      var pipePoints = previewPoints.map(function(p) { return new THREE.Vector3(p.x, p.y, p.z) })
      var tubeRadius = 0.1
      var curve = vm.createRoundedPipeCurve(pipePoints, tubeRadius)
      if (!curve || !curve.curves || curve.curves.length === 0) return

      vm.pipePreviewMesh = createFlowPipe(previewPoints, '#f4ead8', tubeRadius, 1.0, '#ff6a00', 3)
      if (!vm.pipePreviewMesh) return
      vm.pipePreviewMesh.renderOrder = 999
      vm.__3d.scene.add(vm.pipePreviewMesh)
    },

    createPipePreview: function() {
      // 已废弃，使用 updateLivePipePreview 代替
    },

    updatePipePreview: function() {
      // 已废弃，使用 updateLivePipePreview 代替
    },

    finishDrawingPipe: function() {
      var vm = this
      if (vm.pipePoints.length < 2) {
        vm.cancelDrawingPipe()
        return
      }

      var pipeData = {
        type: 'flowPipe',
        points: vm.pipePoints,
        color: '#f4ead8',
        radius: 0.1,
        flowSpeed: 1.0,
        highlightColor: '#ff6a00',
        flowDashLength: 3
      }

      vm.$emit('create-flow-pipe', pipeData)
      vm.cancelDrawingPipe()
    },

    cancelDrawingPipe: function() {
      var vm = this
      vm.isDrawingPipe = false
      vm.pipePoints = []
      if (vm.__3d && vm.__3d.orbit) {
        vm.__3d.orbit.enabled = true
      }

      // 清理旧的线预览
      if (vm.pipePreviewLine) {
        vm.__3d.scene.remove(vm.pipePreviewLine)
        vm.pipePreviewLine.geometry.dispose()
        vm.pipePreviewLine.material.dispose()
        vm.pipePreviewLine = null
      }

      // 清理新的管道预览 Mesh
      if (vm.pipePreviewMesh) {
        vm.__3d.scene.remove(vm.pipePreviewMesh)
        vm.pipePreviewMesh.traverse(function(child) {
          if (child.geometry) child.geometry.dispose()
          if (child.material) child.material.dispose()
        })
        vm.pipePreviewMesh = null
      }

      vm.pipePreviewPoints = []

      vm.$emit('drawing-pipe-cancel')
    },

    onOverlayContextMenu: function(e, overlay) {
      e.preventDefault()
      e.stopPropagation()

      var _t = this
      var obj = overlay.objectData
      if (!obj) return

      this.$emit('object-selected', overlay.id)

      var menuItems = [
        {
          label: '删除',
          icon: 'fas fa-trash',
          onClick: function() {
            _t.$emit('delete-selected')
          }
        },
        {
          label: '复制',
          icon: 'fas fa-copy',
          onClick: function() {
            _t.$emit('duplicate', overlay.id)
          }
        },
        {
          label: '上移',
          icon: 'fas fa-arrow-up',
          onClick: function() {
            _t.$emit('move-up', overlay.id)
          }
        },
        {
          label: '下移',
          icon: 'fas fa-arrow-down',
          onClick: function() {
            _t.$emit('move-down', overlay.id)
          }
        },
        {
          label: '聚焦',
          icon: 'fas fa-crosshairs',
          onClick: function() {
            _t.$emit('focus-selected')
          }
        }
      ]

      if (obj.locked) {
        menuItems.push({
          label: '解锁',
          icon: 'fas fa-unlock',
          onClick: function() {
            _t.$emit('toggle-lock')
          }
        })
      } else {
        menuItems.push({
          label: '锁定',
          icon: 'fas fa-lock',
          onClick: function() {
            _t.$emit('toggle-lock')
          }
        })
      }

      this.$contextmenu({
        items: menuItems,
        event: e,
        zIndex: 10000
      })
    },

    onCanvasMouseMove: function(e) {
      var vm = this
      if (!vm.__3d) return
      if (vm._spacePanActive) {
        vm.hoveredName = ''
        if (vm.$refs.threeCanvas) vm.$refs.threeCanvas.style.cursor = 'grab'
        return
      }

      // 管道绘制模式：实时更新管道预览，让它跟随鼠标
      if (vm.isDrawingPipe && vm.pipePoints.length >= 1) {
        vm.updateLivePipePreview(e)
        return
      }

      var rect = vm.$refs.threeCanvas.getBoundingClientRect()
      vm.__3d.mouse.x = ((e.clientX - rect.left) / rect.width) * 2 - 1
      vm.__3d.mouse.y = -((e.clientY - rect.top) / rect.height) * 2 + 1
      vm.__3d.raycaster.setFromCamera(vm.__3d.mouse, vm.__3d.camera)

      var allMeshes = []
      Object.values(vm.__3d.meshMap).forEach(function(m) {
        if (m && m.traverse) {
          m.traverse(function(c) { if (c.isMesh || c.isPoints || c.isLineSegments) allMeshes.push(c) })
        } else if (m) {
          allMeshes.push(m)
        }
      })

      var hits = vm.__3d.raycaster.intersectObjects(allMeshes, false)
      if (hits.length > 0) {
        var hit = hits[0].object
        var id = hit.userData.id
        if (!id) {
          var p = hit.parent
          while (p) {
            if (p.userData && p.userData.id) { id = p.userData.id; break }
            p = p.parent
          }
        }
        var obj = vm.objects ? vm.objects.find(function(o) { return o.id === id }) : null
        vm.hoveredName = obj ? obj.name : ''
        vm.$refs.threeCanvas.style.cursor = vm.has2DAction(obj) ? 'pointer' : 'default'
        if (id && id === vm.selectedId && vm.mode === 'select') {
          var selectedMesh = vm.__3d.meshMap && vm.__3d.meshMap[id]
          if (selectedMesh && vm.__3d.transform.object !== selectedMesh) {
            vm.attachTransform(selectedMesh)
          }
          vm.setTransformControlVisible(true)
          vm.showTransformGuides(selectedMesh)
          vm.clearTransformAutoHide()
        } else if (vm.mode === 'select') {
          vm.clearTransformAutoHide()
        }
      } else {
        vm.hoveredName = ''
        vm.$refs.threeCanvas.style.cursor = 'default'
        if (vm.mode === 'select') {
          vm.clearTransformAutoHide()
        }
      }
    },

    getCanvasDropTarget: function(e) {
      var vm = this
      var rect = vm.$refs.threeCanvas.getBoundingClientRect()
      var mouseX = ((e.clientX - rect.left) / rect.width) * 2 - 1
      var mouseY = -((e.clientY - rect.top) / rect.height) * 2 + 1

      vm.__3d.raycaster.setFromCamera(new THREE.Vector2(mouseX, mouseY), vm.__3d.camera)

      var meshes = vm.getScenePickMeshes()
      var hits = vm.__3d.raycaster.intersectObjects(meshes, false)
      if (hits.length > 0 && hits[0].point && isFinite(hits[0].point.x)) {
        return hits[0].point.clone()
      }

      var plane = new THREE.Plane(new THREE.Vector3(0, 1, 0), 0)
      var target = new THREE.Vector3()
      var hasIntersection = vm.__3d.raycaster.ray.intersectPlane(plane, target)
      if (hasIntersection && isFinite(target.x) && isFinite(target.y) && isFinite(target.z)) {
        return target
      }

      var fallbackDistance = 10
      if (vm.__3d.orbit && vm.__3d.orbit.target) {
        fallbackDistance = Math.max(1, vm.__3d.camera.position.distanceTo(vm.__3d.orbit.target))
      }
      target.copy(vm.__3d.raycaster.ray.origin).add(vm.__3d.raycaster.ray.direction.clone().multiplyScalar(fallbackDistance))
      if (isFinite(target.x) && isFinite(target.y) && isFinite(target.z)) {
        return target
      }

      return null
    },

    onCanvasDrop: function(e) {
      e.preventDefault()
      var vm = this
      if (!vm.__3d) return

      var itemType = e.dataTransfer.getData('text/plain')
      var itemPayload = null
      var payloadText = e.dataTransfer.getData('application/x-ism3d-object') || e.dataTransfer.getData('application/x-ism3d-asset')
      if (payloadText) {
        try {
          itemPayload = JSON.parse(payloadText)
          itemType = itemPayload.type || itemType
        } catch (err) {
          itemPayload = null
        }
      }
      if (!itemType) return

      var rect = vm.$refs.threeCanvas.getBoundingClientRect()
      var target = vm.getCanvasDropTarget(e)
      if (!target) return

      vm.$emit('drop-object', {
        type: itemType,
        item: itemPayload,
        x: target.x,
        y: target.y,
        z: target.z,
        screenX: e.clientX - rect.left,
        screenY: e.clientY - rect.top
      })
    },

    onOverlayClick: function(overlay) {
      if (!overlay) return
      this.$emit('object-selected', overlay.id)
    },

    getOverlayPosition: function(obj) {
      if (!obj) return { x: 0, y: 0, w: 160, h: 80 }
      if (obj.type === '2dComponent') {
        return obj.source2D && obj.source2D.style && obj.source2D.style.position
          ? obj.source2D.style.position
          : { x: 0, y: 0, w: 160, h: 80 }
      }
      return {
        x: obj.uiX !== undefined ? obj.uiX : 40,
        y: obj.uiY !== undefined ? obj.uiY : 40,
        w: obj.uiWidth !== undefined ? obj.uiWidth : (obj.type === 'webEmbed' ? 420 : 180),
        h: obj.uiHeight !== undefined ? obj.uiHeight : (obj.type === 'webEmbed' ? 260 : 100)
      }
    },

    setOverlayPositionValue: function(obj, key, value) {
      if (!obj) return
      if (obj.type === '2dComponent') {
        if (!obj.source2D) this.$set(obj, 'source2D', { style: { position: {} } })
        if (!obj.source2D.style) this.$set(obj.source2D, 'style', { position: {} })
        if (!obj.source2D.style.position) this.$set(obj.source2D.style, 'position', {})
        this.$set(obj.source2D.style.position, key, value)
      } else {
        var map = { x: 'uiX', y: 'uiY', w: 'uiWidth', h: 'uiHeight' }
        this.$set(obj, map[key] || key, value)
      }
    },

    getOverlayRotation: function(obj) {
      if (!obj) return 0
      if (obj.type === '2dComponent') {
        return obj.source2D && obj.source2D.style ? obj.source2D.style.transform : 0
      }
      return obj.uiRotation || 0
    },

    setOverlayRotation: function(obj, value) {
      if (!obj) return
      if (obj.type === '2dComponent') {
        if (!obj.source2D) this.$set(obj, 'source2D', { style: { position: {} } })
        if (!obj.source2D.style) this.$set(obj.source2D, 'style', { position: {} })
        this.$set(obj.source2D.style, 'transform', value)
      } else {
        this.$set(obj, 'uiRotation', value)
      }
    },

    onOverlayMouseDown: function(e, overlay) {
      if (!overlay || !overlay.objectData) return
      if (overlay.objectData.locked) {
        this.$emit('object-selected', overlay.id)
        return
      }
      var position = this.getOverlayPosition(overlay.objectData)
      this.$emit('object-selected', overlay.id)
      this.$emit('overlay-edit-start', { objectId: overlay.id, mode: 'move' })
      this._overlayDrag = {
        mode: 'move',
        id: overlay.id,
        startMouseX: e.clientX,
        startMouseY: e.clientY,
        startX: parseFloat(position.x) || 0,
        startY: parseFloat(position.y) || 0
      }
    },

    onOverlayHandleMouseDown: function(e, overlay, mode) {
      if (!overlay || !overlay.objectData) return
      if (overlay.objectData.locked) {
        this.$emit('object-selected', overlay.id)
        return
      }
      var position = this.getOverlayPosition(overlay.objectData)
      var width = parseFloat(position.w) || 160
      var height = parseFloat(position.h) || 80
      var transform = parseFloat(this.getOverlayRotation(overlay.objectData))
      if (isNaN(transform) || transform === -1098 || transform === -1099) transform = 0
      this.$emit('object-selected', overlay.id)
      this.$emit('overlay-edit-start', { objectId: overlay.id, mode: mode })
      this._overlayDrag = {
        mode: mode,
        id: overlay.id,
        startMouseX: e.clientX,
        startMouseY: e.clientY,
        startX: parseFloat(position.x) || 0,
        startY: parseFloat(position.y) || 0,
        startW: width,
        startH: height,
        startTransform: transform
      }
    },

    onOverlayMouseMove: function(e) {
      if (!this._overlayDrag) return
      var drag = this._overlayDrag
      var obj = this.objects.find(function(item) { return item.id === drag.id })
      if (!obj) return
      if (drag.mode === 'move') {
        var nextX = drag.startX + (e.clientX - drag.startMouseX)
        var nextY = drag.startY + (e.clientY - drag.startMouseY)
        this.setOverlayPositionValue(obj, 'x', Math.max(0, Math.round(nextX)))
        this.setOverlayPositionValue(obj, 'y', Math.max(0, Math.round(nextY)))
      } else if (drag.mode === 'resize' || drag.mode === 'resize-tl' || drag.mode === 'resize-tr' || drag.mode === 'resize-bl') {
        var deltaX = e.clientX - drag.startMouseX
        var deltaY = e.clientY - drag.startMouseY

        if (drag.mode === 'resize') {
          var nextW = drag.startW + deltaX
          var nextH = drag.startH + deltaY
          this.setOverlayPositionValue(obj, 'w', Math.max(24, Math.round(nextW)))
          this.setOverlayPositionValue(obj, 'h', Math.max(24, Math.round(nextH)))
        } else if (drag.mode === 'resize-tl') {
          var nextW = drag.startW - deltaX
          var nextH = drag.startH - deltaY
          this.setOverlayPositionValue(obj, 'w', Math.max(24, Math.round(nextW)))
          this.setOverlayPositionValue(obj, 'h', Math.max(24, Math.round(nextH)))
          this.setOverlayPositionValue(obj, 'x', Math.max(0, Math.round(drag.startX + deltaX)))
          this.setOverlayPositionValue(obj, 'y', Math.max(0, Math.round(drag.startY + deltaY)))
        } else if (drag.mode === 'resize-tr') {
          var nextW = drag.startW + deltaX
          var nextH = drag.startH - deltaY
          this.setOverlayPositionValue(obj, 'w', Math.max(24, Math.round(nextW)))
          this.setOverlayPositionValue(obj, 'h', Math.max(24, Math.round(nextH)))
          this.setOverlayPositionValue(obj, 'y', Math.max(0, Math.round(drag.startY + deltaY)))
        } else if (drag.mode === 'resize-bl') {
          var nextW = drag.startW - deltaX
          var nextH = drag.startH + deltaY
          this.setOverlayPositionValue(obj, 'w', Math.max(24, Math.round(nextW)))
          this.setOverlayPositionValue(obj, 'h', Math.max(24, Math.round(nextH)))
          this.setOverlayPositionValue(obj, 'x', Math.max(0, Math.round(drag.startX + deltaX)))
        }
      } else if (drag.mode === 'rotate') {
        var centerX = drag.startX + drag.startW / 2
        var centerY = drag.startY + drag.startH / 2
        var startAngle = Math.atan2(drag.startMouseY - centerY, drag.startMouseX - centerX)
        var currentAngle = Math.atan2(e.clientY - centerY, e.clientX - centerX)
        var deltaAngle = currentAngle - startAngle
        this.setOverlayRotation(obj, Math.round(drag.startTransform + (deltaAngle * 180 / Math.PI)))
      }
      this.update2DOverlays()
    },

    onOverlayMouseUp: function() {
      if (this._overlayDrag) {
        this.$emit('overlay-updated', { objectId: this._overlayDrag.id, mode: this._overlayDrag.mode })
      }
      this._overlayDrag = null
    },

    onSelectionFrameHandleMouseDown: function(e, handle) {
      if (!this.__3d || !this.selectedId || !this.selectionFrame.visible) return
      var obj = this.objects ? this.objects.find(function(item) { return item.id === this.selectedId }, this) : null
      if (obj && obj.locked) return
      var mesh = this.__3d.meshMap && this.__3d.meshMap[this.selectedId]
      if (!mesh) return
      this.$emit('object-selected', this.selectedId)
      var frame = this.selectionFrame
      var centerX = frame.left + frame.width / 2
      var centerY = frame.top + frame.height / 2
      var wrapperRect = this.$refs.canvasWrapper.getBoundingClientRect()
      var startX = e.clientX - wrapperRect.left
      var startY = e.clientY - wrapperRect.top
      var startDistance = Math.sqrt(Math.pow(startX - centerX, 2) + Math.pow(startY - centerY, 2))
      this._selectionResizeDrag = {
        id: this.selectedId,
        handle: handle,
        centerX: centerX,
        centerY: centerY,
        startDistance: Math.max(startDistance, 8),
        startScaleX: mesh.scale.x || 1,
        startScaleY: mesh.scale.y || 1,
        startScaleZ: mesh.scale.z || 1
      }
      this.setTransformControlVisible(true)
      this.keepSelectedTransformVisible()
    },

    onSelectionFrameResizeMove: function(e) {
      if (!this._selectionResizeDrag || !this.__3d || !this.$refs.canvasWrapper) return
      var drag = this._selectionResizeDrag
      var mesh = this.__3d.meshMap && this.__3d.meshMap[drag.id]
      if (!mesh) return
      var wrapperRect = this.$refs.canvasWrapper.getBoundingClientRect()
      var currentX = e.clientX - wrapperRect.left
      var currentY = e.clientY - wrapperRect.top
      var distance = Math.sqrt(Math.pow(currentX - drag.centerX, 2) + Math.pow(currentY - drag.centerY, 2))
      var ratio = Math.max(0.08, Math.min(20, distance / drag.startDistance))
      mesh.scale.set(
        drag.startScaleX * ratio,
        drag.startScaleY * ratio,
        drag.startScaleZ * ratio
      )
      this.setTransformControlVisible(true)
      this.updateSelectionFrame(mesh)
      var obj = this.objects ? this.objects.find(function(item) { return item.id === drag.id }) : null
      this.updateObjectReflection(obj || { id: drag.id, visible: true }, mesh)
    },

    onSelectionFrameResizeUp: function() {
      if (!this._selectionResizeDrag || !this.__3d) return
      var drag = this._selectionResizeDrag
      var mesh = this.__3d.meshMap && this.__3d.meshMap[drag.id]
      this._selectionResizeDrag = null
      if (!mesh) return
      this.$emit('object-transformed', {
        objectId: drag.id,
        transform: {
          x: parseFloat(mesh.position.x.toFixed(4)),
          y: parseFloat(mesh.position.y.toFixed(4)),
          z: parseFloat(mesh.position.z.toFixed(4)),
          rx: parseFloat(mesh.rotation.x.toFixed(4)),
          ry: parseFloat(mesh.rotation.y.toFixed(4)),
          rz: parseFloat(mesh.rotation.z.toFixed(4)),
          sx: parseFloat(mesh.scale.x.toFixed(4)),
          sy: parseFloat(mesh.scale.y.toFixed(4)),
          sz: parseFloat(mesh.scale.z.toFixed(4))
        }
      })
      this.keepSelectedTransformVisible()
    },

    onResize: function() {
      var vm = this
      if (!vm.__3d) return
      var wrapper = vm.$refs.canvasWrapper
      if (!wrapper) return
      var rect = wrapper.getBoundingClientRect()
      var w = rect.width
      var h = rect.height
      if (w <= 0 || h <= 0) return
      w = Math.max(w, 100)
      h = Math.max(h, 100)
      vm.__3d.camera.aspect = w / h
      vm.__3d.camera.updateProjectionMatrix()
      vm.__3d.renderer.setSize(w, h)
      vm.update2DOverlays()
    },

    setCameraView: function(view) {
      if (!this.__3d) return
      this.cameraView = view
      var cam = this.__3d.camera
      var orb = this.__3d.orbit
      var d = 10
      if (view === 'perspective') {
        cam.position.set(6, 5, 8); cam.fov = 45
      } else if (view === 'top') {
        cam.position.set(0, d, 0.001); cam.fov = 50
      } else if (view === 'front') {
        cam.position.set(0, 0, d); cam.fov = 45
      } else if (view === 'right') {
        cam.position.set(d, 0, 0); cam.fov = 45
      }
      cam.updateProjectionMatrix()
      cam.lookAt(0, 0, 0)
      orb.target.set(0, 0, 0)
      orb.update()
    },

    setTransformMode: function(mode) {
      if (!this.__3d) return
      if (this.selectedObj && this.selectedObj.locked) {
        this.detachTransform()
        return
      }
      var modeMap = { select: 'translate', move: 'translate', rotate: 'rotate', scale: 'scale' }
      this.__3d.transform.setMode(modeMap[mode] || 'translate')
      var mesh = this.selectedId && this.__3d.meshMap ? this.__3d.meshMap[this.selectedId] : null
      if (mesh && this.__3d.transform.object !== mesh) {
        this.attachTransform(mesh)
      } else if (mesh) {
        // mesh 已 attach，setMode 重建了 gizmo，需重新应用颜色/样式
        this.applyTransformControlContrast(this.__3d.transform)
      }
      this.setTransformControlVisible(!!mesh)
    },

    attachTransform: function(mesh) {
      if (!this.__3d || !mesh) return
      var obj = this.objects ? this.objects.find(function(item) { return item.id === mesh.userData.id }) : null
      if (obj && obj.locked) return
      if (mesh.userData && mesh.userData.isTextSprite && !mesh.userData.is3DText) {
        mesh.rotation.set(obj ? obj.rx || 0 : 0, obj ? obj.ry || 0 : 0, obj ? obj.rz || 0 : 0)
      }
      var modeMap = { select: 'translate', move: 'translate', rotate: 'rotate', scale: 'scale' }
      this.__3d.transform.setMode(modeMap[this.mode] || 'translate')
      this.__3d.transform.attach(mesh)
      this.applyTransformControlContrast(this.__3d.transform)
      this.setTransformControlVisible(true)
      if (this.__3d.techRing) this.__3d.techRing.visible = true
      this.updateSelectionFrame(mesh)
    },

    detachTransform: function() {
      if (!this.__3d) return
      this.clearTransformAutoHide()
      this.setTransformControlVisible(false)
      if (this.__3d.techRing) this.__3d.techRing.visible = false
      this.__3d.transform.detach()
      this.hideTransformGuides(false)
      this.hideSelectionFrame()
    },

    clearTransformAutoHide: function() {
      if (!this._transformAutoHideTimer) return
      clearTimeout(this._transformAutoHideTimer)
      this._transformAutoHideTimer = null
    },

    setTransformControlVisible: function(visible) {
      if (!this.__3d || !this.__3d.transform) return
      if (visible) this.clearTransformAutoHide()
      this.__3d.transform.visible = !!visible
      if (this.__3d.techRing) this.__3d.techRing.visible = !!visible
    },

    keepSelectedTransformVisible: function() {
      if (!this.__3d || !this.selectedId) return
      var mesh = this.__3d.meshMap && this.__3d.meshMap[this.selectedId]
      if (!mesh) return
      if (this.__3d.transform.object !== mesh) {
        this.attachTransform(mesh)
      }
      this.setTransformControlVisible(true)
    },

    scheduleTransformAutoHide: function(delay) {
      return
    },

    applyTransformControlContrast: function(transform) {
      if (!transform || !transform.traverse) return
      transform.size = 1.35

      // TransformControls 的 gizmo 子 mesh 名字就是轴名 "X"/"Y"/"Z"/"XYZ" 等
      var colorForAxis = function(name) {
        if (name === 'XZ') return 0xff7733
        if (name === 'XY') return 0xff22aa
        if (name === 'YZ') return 0x00ccff
        if (name === 'X') return 0xff3399
        if (name === 'Y') return 0x00ffcc
        if (name === 'Z') return 0x6655ff
        if (name === 'E' || name === 'XYZE') return 0xffaa00
        return undefined
      }

      var isPlane = function(name) {
        return name === 'XY' || name === 'YZ' || name === 'XZ'
      }

      transform.traverse(function(child) {
        if (!child) return

        var name = child.name || ''
        var gType = child.geometry ? child.geometry.type : ''
        var plane = isPlane(name)

        if (name === 'XYZ' && (gType === 'OctahedronGeometry' || gType === 'BoxGeometry')) {
          child.visible = true
        }

        if (!child.material) return

        var color = colorForAxis(name)

        var materials = Array.isArray(child.material) ? child.material : [child.material]
        materials.forEach(function(mat) {
          if (!mat || mat.name === 'matInvisible') return
          if (color !== undefined && mat.color) {
            mat.color.setHex(color)
          } else if (name === 'XYZ' && mat.color) {
            mat.color.setHex(0x13c2c2)
          }
          if (mat.linewidth !== undefined) mat.linewidth = 2.5
          mat.transparent = true
          if (name === 'XYZ') {
            mat.opacity = 0.55
          } else {
            mat.opacity = plane ? 0.25 : 0.95
          }
          mat.depthTest = false
          mat.depthWrite = false
          mat.needsUpdate = true
          if (plane && mat.emissive) {
            mat.emissive.setHex(color)
            mat.emissiveIntensity = 0.4
          } else if (name === 'XYZ' && mat.emissive) {
            mat.emissive.setHex(0x13c2c2)
            mat.emissiveIntensity = 0.25
          }
        })
        child.renderOrder = 1000
        child.userData.isAux = true
      })
    },

    projectWorldToCanvas: function(position) {
      if (!this.__3d || !this.$refs.canvasWrapper) return { left: 0, top: 0 }
      var projected = position.clone().project(this.__3d.camera)
      var rect = this.$refs.canvasWrapper.getBoundingClientRect()
      return {
        left: (projected.x * 0.5 + 0.5) * rect.width,
        top: (-projected.y * 0.5 + 0.5) * rect.height
      }
    },

    updateActiveSelectionFrame: function() {
      if (!this.__3d || !this.selectedId) {
        this.hideSelectionFrame()
        this.dragGuide.visible = false
        return
      }
      var mesh = this.__3d.meshMap && this.__3d.meshMap[this.selectedId]
      if (!mesh) {
        this.hideSelectionFrame()
        this.dragGuide.visible = false
        return
      }
      this.updateSelectionFrame(mesh)
      this.updateSelectedCoordinateTip()
    },

    updateSelectedCoordinateTip: function(alignText) {
      if (!this.__3d || !this.selectedId) {
        this.dragGuide.visible = false
        return
      }
      var mesh = this.__3d.meshMap && this.__3d.meshMap[this.selectedId]
      if (!mesh) {
        this.dragGuide.visible = false
        return
      }
      var obj = this.objects ? this.objects.find(function(item) { return item.id === mesh.userData.id }) : null
      var box = new THREE.Box3().setFromObject(mesh)
      var anchor = mesh.position.clone()
      if (box && !box.isEmpty()) {
        anchor.set((box.min.x + box.max.x) / 2, box.max.y, (box.min.z + box.max.z) / 2)
      }
      var screen = this.projectWorldToCanvas(anchor)
      var x = mesh.position.x
      var y = mesh.position.y
      var z = mesh.position.z
      this.dragGuide = {
        visible: true,
        left: Math.round(screen.left + 14),
        top: Math.round(screen.top - 36),
        x: (isFinite(x) ? x : 0).toFixed(2),
        y: (isFinite(y) ? y : 0).toFixed(2),
        z: (isFinite(z) ? z : 0).toFixed(2),
        name: obj && obj.name ? obj.name : 'Selected',
        alignText: alignText || ''
      }
    },

    updateSelectionFrame: function(mesh) {
      if (!this.__3d || !mesh || !this.$refs.canvasWrapper) return
      if (mesh.userData && mesh.userData.isTextSprite && !mesh.userData.is3DText) {
        this.hideSelectionFrame()
        return
      }
      var box = new THREE.Box3().setFromObject(mesh)
      if (!box || box.isEmpty()) {
        this.hideSelectionFrame()
        return
      }
      var min = box.min
      var max = box.max
      var points = [
        new THREE.Vector3(min.x, min.y, min.z),
        new THREE.Vector3(min.x, min.y, max.z),
        new THREE.Vector3(min.x, max.y, min.z),
        new THREE.Vector3(min.x, max.y, max.z),
        new THREE.Vector3(max.x, min.y, min.z),
        new THREE.Vector3(max.x, min.y, max.z),
        new THREE.Vector3(max.x, max.y, min.z),
        new THREE.Vector3(max.x, max.y, max.z)
      ]
      var left = Infinity
      var top = Infinity
      var right = -Infinity
      var bottom = -Infinity
      points.forEach(function(point) {
        var projected = this.projectWorldToCanvas(point)
        left = Math.min(left, projected.left)
        top = Math.min(top, projected.top)
        right = Math.max(right, projected.left)
        bottom = Math.max(bottom, projected.top)
      }, this)

      if (!isFinite(left) || !isFinite(top) || !isFinite(right) || !isFinite(bottom)) {
        this.hideSelectionFrame()
        return
      }
      var pad = 7
      this.selectionFrame = {
        visible: true,
        left: Math.round(left - pad),
        top: Math.round(top - pad),
        width: Math.max(36, Math.round(right - left + pad * 2)),
        height: Math.max(28, Math.round(bottom - top + pad * 2))
      }
    },

    hideSelectionFrame: function() {
      if (!this.selectionFrame.visible) return
      this.selectionFrame = {
        visible: false,
        left: 0,
        top: 0,
        width: 0,
        height: 0
      }
    },

    createGuideLine: function(points, color) {
      var geometry = new THREE.BufferGeometry().setFromPoints(points)
      var material = new THREE.LineDashedMaterial({
        color: color,
        transparent: true,
        opacity: 0.55,
        depthTest: false,
        dashSize: 0.6,
        gapSize: 0.3
      })
      var line = new THREE.Line(geometry, material)
      line.computeLineDistances()
      line.renderOrder = 999
      line.userData.isAux = true
      return line
    },

    createDragShadowTexture: function() {
      var canvas = document.createElement('canvas')
      canvas.width = 256
      canvas.height = 256
      var ctx = canvas.getContext('2d')
      var gradient = ctx.createRadialGradient(128, 128, 12, 128, 128, 118)
      gradient.addColorStop(0, 'rgba(0, 0, 0, 0.34)')
      gradient.addColorStop(0.45, 'rgba(0, 0, 0, 0.2)')
      gradient.addColorStop(0.75, 'rgba(0, 0, 0, 0.08)')
      gradient.addColorStop(1, 'rgba(0, 0, 0, 0)')
      ctx.fillStyle = gradient
      ctx.fillRect(0, 0, canvas.width, canvas.height)
      var texture = new THREE.CanvasTexture(canvas)
      texture.needsUpdate = true
      return texture
    },

    showDragShadow: function(mesh) {
      if (!this.__3d || !this.__3d.scene || !mesh) return
      this.hideDragShadow()
      var texture = this.createDragShadowTexture()
      var geometry = new THREE.PlaneGeometry(1, 1)
      var material = new THREE.MeshBasicMaterial({
        map: texture,
        transparent: true,
        opacity: 0.72,
        depthWrite: false,
        depthTest: false,
        side: THREE.DoubleSide
      })
      var shadow = new THREE.Mesh(geometry, material)
      shadow.rotation.x = -Math.PI / 2
      shadow.renderOrder = 998
      shadow.userData.isDragShadow = true
      shadow.userData.isAux = true
      this.__3d.scene.add(shadow)
      this.__3d.dragShadowMesh = shadow
      if (this.__3d._auxObjects) this.__3d._auxObjects.push(shadow)
      this.updateDragShadow(mesh)
    },

    updateDragShadow: function(mesh) {
      if (!this.__3d || !mesh || !this.__3d.dragShadowMesh) return
      var shadow = this.__3d.dragShadowMesh
      var box = new THREE.Box3().setFromObject(mesh)
      var size = new THREE.Vector3()
      if (box && !box.isEmpty()) {
        box.getSize(size)
      }
      var radiusX = Math.max(0.55, Math.min(12, size.x || 1))
      var radiusZ = Math.max(0.55, Math.min(12, size.z || 1))
      var scaleX = Math.max(1.1, radiusX * 1.5)
      var scaleZ = Math.max(1.1, radiusZ * 1.5)
      shadow.position.set(mesh.position.x, 0.018, mesh.position.z)
      shadow.scale.set(scaleX, scaleZ, 1)
      shadow.visible = true
    },

    hideDragShadow: function() {
      if (!this.__3d || !this.__3d.dragShadowMesh) return
      var shadow = this.__3d.dragShadowMesh
      this.__3d.scene.remove(shadow)
      if (shadow.geometry) shadow.geometry.dispose()
      if (shadow.material) {
        if (shadow.material.map) shadow.material.map.dispose()
        shadow.material.dispose()
      }
      this.__3d.dragShadowMesh = null
    },

    showTransformGuides: function(mesh) {
      if (!this.__3d || !this.__3d.scene || !mesh) return
      this.hideTransformGuides(false)
      var scene = this.__3d.scene
      var p = mesh.position
      var xLine = this.createGuideLine([new THREE.Vector3(-80, p.y, p.z), new THREE.Vector3(80, p.y, p.z)], 0xdc3b5a)
      var zLine = this.createGuideLine([new THREE.Vector3(p.x, p.y, -80), new THREE.Vector3(p.x, p.y, 80)], 0x1a7bff)
      scene.add(xLine)
      scene.add(zLine)
      this.__3d.dragGuideLines = { xLine: xLine, zLine: zLine }
      if (this.__3d._auxObjects) {
        this.__3d._auxObjects.push(xLine)
        this.__3d._auxObjects.push(zLine)
      }
      this.dragGuide.visible = true
    },

    hideTransformGuides: function(hideTip) {
      if (hideTip !== false) {
        this.dragGuide.visible = false
      }
      if (!this.__3d || !this.__3d.dragGuideLines) return
      var scene = this.__3d.scene
      Object.keys(this.__3d.dragGuideLines).forEach(function(key) {
        var line = this.__3d.dragGuideLines[key]
        if (!line) return
        scene.remove(line)
        if (line.geometry) line.geometry.dispose()
        if (line.material) line.material.dispose()
      }, this)
      this.__3d.dragGuideLines = null
    },

    getSnapStep: function() {
      return 0.05
    },

    applyTransformSnapAndAlignment: function(mesh) {
      if (!this.__3d || !mesh || !this.__3d.transform || this.__3d.transform.mode !== 'translate') return ''
      var p = mesh.position
      var threshold = 0.08
      var alignments = []
      var currentId = mesh.userData && mesh.userData.id
      ;(this.objects || []).forEach(function(obj) {
        if (!obj || obj.id === currentId) return
        if (typeof obj.x === 'number' && Math.abs(p.x - obj.x) <= threshold) {
          p.x = obj.x
          alignments.push('X 对齐 ' + (obj.name || obj.id))
        }
        if (typeof obj.z === 'number' && Math.abs(p.z - obj.z) <= threshold) {
          p.z = obj.z
          alignments.push('Z 对齐 ' + (obj.name || obj.id))
        }
        if (typeof obj.y === 'number' && Math.abs(p.y - obj.y) <= threshold && Math.abs(p.y - obj.y) > 0) {
          p.y = obj.y
          alignments.push('Y 对齐 ' + (obj.name || obj.id))
        }
      })
      return alignments.slice(0, 2).join(' / ')
    },

    updateTransformGuides: function(mesh) {
      if (!this.__3d || !mesh) return
      var p = mesh.position
      var alignText = this.applyTransformSnapAndAlignment(mesh)
      if (this.__3d.dragGuideLines) {
        var xLine = this.__3d.dragGuideLines.xLine
        var zLine = this.__3d.dragGuideLines.zLine
        if (xLine) {
          xLine.geometry.dispose()
          xLine.geometry = new THREE.BufferGeometry().setFromPoints([
            new THREE.Vector3(-80, p.y, p.z),
            new THREE.Vector3(80, p.y, p.z)
          ])
          xLine.computeLineDistances()
        }
        if (zLine) {
          zLine.geometry.dispose()
          zLine.geometry = new THREE.BufferGeometry().setFromPoints([
            new THREE.Vector3(p.x, p.y, -80),
            new THREE.Vector3(p.x, p.y, 80)
          ])
          zLine.computeLineDistances()
        }
      }
      this.updateSelectedCoordinateTip(alignText)
      return
      var screen = this.projectWorldToCanvas(new THREE.Vector3(p.x, p.y + 0.85, p.z))
      var obj = this.objects ? this.objects.find(function(item) { return item.id === mesh.userData.id }) : null
      this.dragGuide = {
        visible: true,
        left: Math.round(screen.left + 14),
        top: Math.round(screen.top - 36),
        x: p.x.toFixed(2),
        y: p.y.toFixed(2),
        z: p.z.toFixed(2),
        name: obj && obj.name ? obj.name : '选中组件',
        alignText: alignText
      }
    },

    focusOn: function(obj) {
      if (!obj || !this.__3d) return
      this.__3d.orbit.target.set(obj.x, obj.y, obj.z)
      this.__3d.camera.position.set(obj.x + 4, obj.y + 3, obj.z + 5)
      this.__3d.orbit.update()
    },

    restoreCamera: function(sceneSettings) {
      if (!this.__3d || !sceneSettings) return
      var cp = sceneSettings.cameraPosition
      var ct = sceneSettings.cameraTarget
      var fov = sceneSettings.cameraFov
      if (cp && typeof cp.x === 'number') {
        this.__3d.camera.position.set(cp.x, cp.y, cp.z)
      }
      if (ct && typeof ct.x === 'number') {
        this.__3d.orbit.target.set(ct.x, ct.y, ct.z)
      }
      if (typeof fov === 'number' && fov > 0) {
        this.__3d.camera.fov = fov
      }
      this.__3d.camera.updateProjectionMatrix()
      this.__3d.orbit.update()
    },

    frameAll: function(sceneObjects, meshMap) {
      if (!sceneObjects || sceneObjects.length === 0 || !this.__3d) return
      var box = new THREE.Box3()
      Object.values(meshMap).forEach(function(m) { box.expandByObject(m) })
      var center = box.getCenter(new THREE.Vector3())
      var size = box.getSize(new THREE.Vector3())
      var maxDim = Math.max(size.x, size.y, size.z)
      this.__3d.orbit.target.copy(center)
      this.__3d.camera.position.set(center.x + maxDim * 1.5, center.y + maxDim, center.z + maxDim * 1.5)
      this.__3d.orbit.update()
    },

    /** 聚焦当前选中对象，由父组件调用 */
    focusSelected: function() {
      if (!this.selectedId || !this.__3d) return
      var mesh = this.__3d.meshMap[this.selectedId]
      if (mesh) {
        var pos = new THREE.Vector3()
        mesh.getWorldPosition(pos)
        this.focusOn({ x: pos.x, y: pos.y, z: pos.z })
      }
    },

    toggleGrid: function(visible) {
      if (this.__3d && this.__3d.gridHelper) this.__3d.gridHelper.visible = visible
    },

    createGradientTexture: function(colorTop, colorBottom) {
      var canvas = document.createElement('canvas')
      canvas.width = 16
      canvas.height = 256
      var ctx = canvas.getContext('2d')
      var gradient = ctx.createLinearGradient(0, 0, 0, canvas.height)
      gradient.addColorStop(0, colorTop || '#ffffff')
      gradient.addColorStop(1, colorBottom || '#dbefff')
      ctx.fillStyle = gradient
      ctx.fillRect(0, 0, canvas.width, canvas.height)
      var texture = new THREE.CanvasTexture(canvas)
      texture.needsUpdate = true
      return texture
    },

    getEnvironmentStyle: function(preset) {
      var styles = {
        sky: { top: '#1fa4bd', bottom: '#59d1e4', floor: '#8ea0a3', fog: '#cceff5' },
        sunset: { top: '#f39a52', bottom: '#ffe1b8', floor: '#9d8f82', fog: '#ffe0bf' },
        ocean: { top: '#117da2', bottom: '#65d7e8', floor: '#6f9ca6', fog: '#c7f4fb' },
        forest: { top: '#2d8f71', bottom: '#bfe9c8', floor: '#6e8a76', fog: '#d6f0dc' },
        twilight: { top: '#4855a9', bottom: '#c9bcff', floor: '#7f7c94', fog: '#ded8ff' },
        night: { top: '#0d1630', bottom: '#283a66', floor: '#394657', fog: '#283852' }
      }
      return styles[preset || 'sky'] || styles.sky
    },

    getSkyboxStyle: function(preset, envStyle) {
      var styleMap = {
        horizon: {
          top: envStyle.top,
          middle: envStyle.bottom,
          bottom: '#d7e2ea',
          horizon: '#8ba0af',
          glow: '#ffffff'
        },
        engineeringSky: {
          top: '#b9dce7',
          middle: '#d8eef5',
          bottom: '#aab8c2',
          horizon: '#667782',
          glow: '#f7fbff'
        },
        blueSky: {
          top: '#67b7ff',
          middle: '#b8ecff',
          bottom: '#eef8ff',
          horizon: '#a8c7d8',
          glow: '#ffffff'
        },
        sunsetGlow: {
          top: '#465f9e',
          middle: '#ffb06a',
          bottom: '#ffe0bd',
          horizon: '#8d7f91',
          glow: '#fff2cf'
        },
        deepNight: {
          top: '#07111f',
          middle: '#183257',
          bottom: '#3a4f68',
          horizon: '#4e6276',
          glow: '#a9c7ff'
        }
      }
      return styleMap[preset || 'horizon'] || styleMap.horizon
    },

    createSkyboxTexture: function(style) {
      var canvas = document.createElement('canvas')
      canvas.width = 1024
      canvas.height = 512
      var ctx = canvas.getContext('2d')
      var gradient = ctx.createLinearGradient(0, 0, 0, canvas.height)
      gradient.addColorStop(0, style.top)
      gradient.addColorStop(0.45, style.middle)
      gradient.addColorStop(0.7, style.bottom)
      gradient.addColorStop(1, style.bottom)
      ctx.fillStyle = gradient
      ctx.fillRect(0, 0, canvas.width, canvas.height)

      var horizonY = Math.round(canvas.height * 0.58)
      var horizonGradient = ctx.createLinearGradient(0, horizonY - 18, 0, horizonY + 18)
      horizonGradient.addColorStop(0, 'rgba(255,255,255,0)')
      horizonGradient.addColorStop(0.5, this.hexToCanvasRgba(style.horizon, 0.75))
      horizonGradient.addColorStop(1, 'rgba(255,255,255,0)')
      ctx.fillStyle = horizonGradient
      ctx.fillRect(0, horizonY - 18, canvas.width, 36)

      for (var i = 0; i < 12; i++) {
        var x = (canvas.width / 12) * i + (i % 2 === 0 ? 20 : 48)
        var y = 38 + (i % 3) * 22
        var w = 130 + (i % 4) * 35
        var h = 6 + (i % 2) * 3
        ctx.fillStyle = 'rgba(255,255,255,0.22)'
        ctx.beginPath()
        ctx.roundRect ? ctx.roundRect(x, y, w, h, 6) : roundedRectPath(ctx, x, y, w, h, 6)
        ctx.fill()
      }

      var glow = ctx.createRadialGradient(canvas.width * 0.7, canvas.height * 0.22, 10, canvas.width * 0.7, canvas.height * 0.22, 120)
      glow.addColorStop(0, this.hexToCanvasRgba(style.glow, 0.9))
      glow.addColorStop(0.4, this.hexToCanvasRgba(style.glow, 0.22))
      glow.addColorStop(1, 'rgba(255,255,255,0)')
      ctx.fillStyle = glow
      ctx.beginPath()
      ctx.arc(canvas.width * 0.7, canvas.height * 0.22, 120, 0, Math.PI * 2)
      ctx.fill()

      var texture = new THREE.CanvasTexture(canvas)
      texture.colorSpace = THREE.SRGBColorSpace
      texture.needsUpdate = true
      return texture
    },

    hexToCanvasRgba: function(hex, alpha) {
      var color = normColorValue(hex, '#ffffff')
      var r = parseInt(color.slice(1, 3), 16)
      var g = parseInt(color.slice(3, 5), 16)
      var b = parseInt(color.slice(5, 7), 16)
      return 'rgba(' + r + ',' + g + ',' + b + ',' + alpha + ')'
    },

    applySkybox: function(gs, envStyle) {
      if (!this.__3d || !this.__3d.scene) return
      var scene = this.__3d.scene
      if (this.__3d.skyboxMesh) {
        scene.remove(this.__3d.skyboxMesh)
        if (this.__3d.skyboxMesh.geometry) this.__3d.skyboxMesh.geometry.dispose()
        if (this.__3d.skyboxMesh.material) {
          if (this.__3d.skyboxMesh.material.map) this.__3d.skyboxMesh.material.map.dispose()
          this.__3d.skyboxMesh.material.dispose()
        }
        this.__3d.skyboxMesh = null
      }
      if (this.__3d.skyboxEnvTexture) {
        if (this.__3d.scene.environment === this.__3d.skyboxEnvTexture) {
          this.__3d.scene.environment = null
        }
        this.__3d.skyboxEnvTexture.dispose()
        this.__3d.skyboxEnvTexture = null
      }
      if (gs.skyboxEnabled === false) return
      if (gs.skyboxPreset === 'customHdri' && gs.skyboxHdri) {
        var hdriLoader = new RGBELoader()
        var self = this
        hdriLoader.load(gs.skyboxHdri, function(texture) {
          texture.mapping = THREE.EquirectangularReflectionMapping
          self.__3d.scene.background = texture
          self.__3d.scene.environment = texture
          self.__3d.skyboxEnvTexture = texture
        })
        return
      }
      if (gs.skyboxPreset === 'customImage' && gs.skyboxImage) {
        var panoLoader = new THREE.TextureLoader()
        var vm = this
        panoLoader.load(gs.skyboxImage, function(texture) {
          texture.mapping = THREE.EquirectangularReflectionMapping
          texture.colorSpace = THREE.SRGBColorSpace
          vm.__3d.scene.background = texture
          var geometry = new THREE.SphereGeometry(180, 48, 32)
          var material = new THREE.MeshBasicMaterial({
            map: texture,
            side: THREE.BackSide,
            depthWrite: false,
            fog: false
          })
          var mesh = new THREE.Mesh(geometry, material)
          mesh.userData.isAux = true
          vm.__3d.scene.add(mesh)
          vm.__3d.skyboxMesh = mesh
        })
        return
      }
      var skyStyle = this.getSkyboxStyle(gs.skyboxPreset || 'horizon', envStyle)
      var texture = this.createSkyboxTexture(skyStyle)
      var geometry = new THREE.SphereGeometry(180, 48, 32)
      var material = new THREE.MeshBasicMaterial({
        map: texture,
        side: THREE.BackSide,
        depthWrite: false,
        fog: false
      })
      var mesh = new THREE.Mesh(geometry, material)
      mesh.userData.isAux = true
      scene.add(mesh)
      this.__3d.skyboxMesh = mesh
    },

    getGroundStyle: function(style, envStyle) {
      var styles = {
        slate: { color: '#4f5961', stripe: '#626d77', stripeOpacity: 0.28 },
        sand: { color: '#bca07b', stripe: '#d4c2a5', stripeOpacity: 0.18 },
        concrete: { color: '#8b9095', stripe: '#b0b4b8', stripeOpacity: 0.16 },
        oceanDark: { color: '#2e475d', stripe: '#4f6b80', stripeOpacity: 0.22 },
        tunnelSection: { color: '#55585d', stripe: '#7b8087', stripeOpacity: 0.2 },
        geologySection: { color: '#6e5c4a', stripe: '#9b8262', stripeOpacity: 0.18 }
      }
      return styles[style || 'slate'] || { color: envStyle.floor || '#8ea0a3', stripe: '#ffffff', stripeOpacity: 0.12 }
    },

    createGroundTexture: function(style) {
      var canvas = document.createElement('canvas')
      canvas.width = 1024
      canvas.height = 1024
      var ctx = canvas.getContext('2d')
      ctx.fillStyle = style.color
      ctx.fillRect(0, 0, canvas.width, canvas.height)
      for (var y = 0; y < canvas.height; y += 56) {
        ctx.fillStyle = this.hexToCanvasRgba(style.stripe, style.stripeOpacity)
        ctx.fillRect(0, y, canvas.width, 10)
      }
      for (var x = 0; x < canvas.width; x += 90) {
        ctx.fillStyle = this.hexToCanvasRgba('#000000', 0.05)
        ctx.fillRect(x, 0, 8, canvas.height)
      }
      for (var i = 0; i < 220; i++) {
        ctx.fillStyle = this.hexToCanvasRgba('#ffffff', 0.03 + (i % 3) * 0.01)
        ctx.beginPath()
        ctx.arc(Math.random() * canvas.width, Math.random() * canvas.height, 1 + Math.random() * 2.5, 0, Math.PI * 2)
        ctx.fill()
      }
      var texture = new THREE.CanvasTexture(canvas)
      texture.wrapS = THREE.RepeatWrapping
      texture.wrapT = THREE.RepeatWrapping
      texture.repeat.set(18, 18)
      texture.colorSpace = THREE.SRGBColorSpace
      texture.needsUpdate = true
      return texture
    },

    createTunnelSectionMesh: function(style, gs) {
      var group = new THREE.Group()
      var tunnelRadius = Math.max(2, Number(gs.tunnelRadius || 18))
      var sectionLength = Math.max(20, Number(gs.sectionLength || 220))
      var sectionWidth = Math.max(6, Number(gs.sectionWidth || 36))
      var sectionHeight = Math.max(4, Number(gs.sectionHeight || tunnelRadius))
      var wallMat = new THREE.MeshStandardMaterial({
        color: new THREE.Color(style.color),
        roughness: 0.95,
        metalness: 0.02
      })
      var shell = new THREE.Mesh(new THREE.CylinderGeometry(tunnelRadius, tunnelRadius, sectionLength, 48, 1, false, Math.PI, Math.PI), wallMat)
      shell.rotation.z = Math.PI / 2
      shell.scale.z = sectionWidth / (tunnelRadius * 2)
      shell.position.y = sectionHeight * 0.49
      shell.userData.isEnvironmentFloor = true
      group.add(shell)

      var floor = new THREE.Mesh(
        new THREE.BoxGeometry(sectionLength, 0.4, sectionWidth),
        new THREE.MeshStandardMaterial({ color: new THREE.Color('#4a4e54'), roughness: 0.92, metalness: 0.03 })
      )
      floor.position.set(0, -0.2, 0)
      floor.userData.isEnvironmentFloor = true
      group.add(floor)
      return group
    },

    createGeologySectionMesh: function(gs) {
      var group = new THREE.Group()
      var sectionLength = Math.max(20, Number(gs.sectionLength || 220))
      var sectionWidth = Math.max(6, Number(gs.sectionWidth || 36))
      var sectionHeight = Math.max(4, Number(gs.sectionHeight || 18))
      var layers = [
        { color: '#7a6248', ratio: 0.22 },
        { color: '#8b7153', ratio: 0.28 },
        { color: '#5d4d3c', ratio: 0.5 }
      ]
      var accumulated = 0
      layers.forEach(function(layer) {
        var height = sectionHeight * layer.ratio
        var centerY = -(accumulated + height / 2)
        var mesh = new THREE.Mesh(
          new THREE.BoxGeometry(sectionLength, height, sectionWidth),
          new THREE.MeshStandardMaterial({ color: new THREE.Color(layer.color), roughness: 0.98, metalness: 0.01 })
        )
        mesh.position.set(0, centerY, 0)
        mesh.userData.isEnvironmentFloor = true
        group.add(mesh)
        accumulated += height
      })
      return group
    },

    applyBackgroundSettings: function(gs) {
      if (!this.__3d || !this.__3d.scene) return
      gs = gs || {}
      var envStyle = this.getEnvironmentStyle(gs.environmentPreset)
      var mode = gs.backgroundMode || (gs.backgroundImage ? 'image' : 'solid')
      if (mode === 'image' && gs.backgroundImage) {
        var loader = new THREE.TextureLoader()
        var self = this
        loader.load(gs.backgroundImage, function(texture) {
          texture.wrapS = THREE.RepeatWrapping
          texture.wrapT = THREE.RepeatWrapping
          self.__3d.scene.background = texture
        }, undefined, function() {
          var fallbackColor = parseInt((gs.backgroundColor || '#ffffff').replace('#', '0x'), 16)
          self.__3d.scene.background = new THREE.Color(fallbackColor)
        })
      } else if (mode === 'gradient') {
        this.__3d.scene.background = this.createGradientTexture(gs.backgroundColor || envStyle.top, gs.backgroundColor2 || envStyle.bottom)
      } else {
        this.__3d.scene.background = this.createGradientTexture(envStyle.top, gs.backgroundColor || envStyle.bottom)
      }
      this.applySkybox(gs, envStyle)
      this.applyEnvironmentFloor(gs, envStyle)
    },

    applyEnvironmentFloor: function(gs, envStyle) {
      if (!this.__3d || !this.__3d.scene) return
      var scene = this.__3d.scene
      if (this.__3d.environmentFloorMesh) {
        scene.remove(this.__3d.environmentFloorMesh)
        if (this.__3d.environmentFloorMesh.geometry) this.__3d.environmentFloorMesh.geometry.dispose()
        if (this.__3d.environmentFloorMesh.material) {
          if (this.__3d.environmentFloorMesh.material.map) this.__3d.environmentFloorMesh.material.map.dispose()
          this.__3d.environmentFloorMesh.material.dispose()
        }
        this.__3d.environmentFloorMesh = null
      }
      if (gs.backgroundMode === 'image' || gs.groundEnabled === false) return
      var groundStyle = this.getGroundStyle(gs.groundStyle || 'slate', envStyle)
      if (gs.groundStyle === 'tunnelSection') {
        var tunnelGroup = this.createTunnelSectionMesh(groundStyle, gs)
        scene.add(tunnelGroup)
        this.__3d.environmentFloorMesh = tunnelGroup
        return
      }
      if (gs.groundStyle === 'geologySection') {
        var geologyGroup = this.createGeologySectionMesh(gs)
        scene.add(geologyGroup)
        this.__3d.environmentFloorMesh = geologyGroup
        return
      }
      var groundTexture = this.createGroundTexture(groundStyle)
      var geometry = new THREE.PlaneGeometry(260, 260)
      var material = new THREE.MeshStandardMaterial({
        color: new THREE.Color(groundStyle.color || envStyle.floor || '#8ea0a3'),
        map: groundTexture,
        roughness: 0.9,
        metalness: 0.04,
        transparent: true,
        opacity: 0.96,
        side: THREE.DoubleSide
      })
      var mesh = new THREE.Mesh(geometry, material)
      mesh.rotation.x = -Math.PI / 2
      mesh.position.y = -0.035
      mesh.receiveShadow = true
      mesh.userData.isEnvironmentFloor = true
      scene.add(mesh)
      this.__3d.environmentFloorMesh = mesh
    },

    updateGrid: function(settings) {
      if (!this.__3d) return
      var gs = settings
      this.gridSettings = Object.assign({}, gs)
      if (this.__3d.gridHelper) {
        this.__3d.scene.remove(this.__3d.gridHelper)
      }
      var centerColor = parseInt((gs.colorCenterLine || '#111111').replace('#', '0x'), 16)
      var gridColor = parseInt((gs.colorGrid || '#cccccc').replace('#', '0x'), 16)
      this.__3d.gridHelper = new THREE.GridHelper(gs.size || 20, gs.divisions || 20, centerColor, gridColor)
      this.__3d.gridHelper.visible = true
      this.__3d.scene.add(this.__3d.gridHelper)

      this.applyBackgroundSettings(gs)
      if (this.__3d.scene.fog) {
        var envStyle = this.getEnvironmentStyle(gs.environmentPreset)
        var fogColor = parseInt((envStyle.fog || gs.backgroundColor || '#ffffff').replace('#', '0x'), 16)
        this.__3d.scene.fog = new THREE.Fog(fogColor, 40, 100)
      }
      this.applySceneVisualSettings(gs)
    },

    applySceneVisualSettings: function(settings) {
      if (!this.__3d || !this.__3d.scene) return
      var gs = settings || this.gridSettings || {}
      this.applyLightingPreset(gs.lightingPreset || 'day', !!gs.enhanceDepth)
      this.applyFloorReflection(gs.floorReflection || 'none')
      this.applyDepthEnhancement(gs)
      // 应用自定义灯光设置（覆盖预设）
      if (gs.lightSettings) {
        this.applyLightSettings(gs.lightSettings)
      }
      if (this.__3d.renderer) {
        this.__3d.renderer.setPixelRatio(gs.modelOptimize ? Math.min(window.devicePixelRatio || 1, 1.25) : (window.devicePixelRatio || 1))
        this.__3d.renderer.shadowMap.enabled = !gs.modelOptimize
        // 如果没有自定义灯光设置，使用默认 exposure
        if (!gs.lightSettings || gs.lightSettings.envIntensity === undefined) {
          this.__3d.renderer.toneMappingExposure = gs.enhanceDepth ? 1.12 : 1.0
        }
      }
    },

    applyDepthEnhancement: function(gs) {
      if (!this.__3d || !this.__3d.scene) return
      var scene = this.__3d.scene
      if (gs.enhanceDepth) {
        var fogColor = parseInt((gs.backgroundColor || '#ffffff').replace('#', '0x'), 16)
        scene.fog = new THREE.Fog(fogColor, 36, 95)
      } else if (scene.fog && !scene.userData.fogEnabledByExtras) {
        scene.fog = null
      }
      var objectsById = {}
      ;(this.objects || []).forEach(function(obj) {
        if (obj && obj.id) objectsById[obj.id] = obj
      })
      scene.traverse(function(child) {
        if (!child.isMesh || !child.material) return
        var objectId = child.userData && child.userData.id
        var obj = objectId ? objectsById[objectId] : null
        if (!obj) {
          var parent = child.parent
          while (parent && !obj) {
            var parentId = parent.userData && parent.userData.id
            if (parentId) obj = objectsById[parentId]
            parent = parent.parent
          }
        }
        var shadowEnabled = !!(obj && obj.showShadow && !gs.modelOptimize)
        child.castShadow = shadowEnabled
        child.receiveShadow = shadowEnabled
        var materials = Array.isArray(child.material) ? child.material : [child.material]
        materials.forEach(function(material) {
          if (!material) return
          if (material.roughness !== undefined && gs.enhanceDepth) {
            material.roughness = Math.max(0.22, material.roughness * 0.92)
          }
          material.needsUpdate = true
        })
      })
    },

    applyLightingPreset: function(preset, enhanceDepth) {
      if (!this.__3d || !this.__3d.scene) return
      var scene = this.__3d.scene
      if (!this.__3d.sceneLights) {
        this.__3d.sceneLights = scene.children.filter(function(child) { return child.isLight })
      }
      var configs = {
        day: { ambient: 1.0, directional: 0.95, color: 0xffffff, fill: 0x4466aa, position: [8, 12, 6] },
        evening: { ambient: 0.58, directional: 1.05, color: 0xffb36f, fill: 0x6e7fb8, position: [-8, 5, 8] },
        night: { ambient: 0.28, directional: 0.55, color: 0x9ec8ff, fill: 0x1b2a55, position: [4, 9, -7] },
        studio: { ambient: 1.18, directional: 0.8, color: 0xffffff, fill: 0xdde8ff, position: [0, 10, 8] },
        industrial: { ambient: 0.78, directional: 1.18, color: 0xf2f7ff, fill: 0x9fb7c8, position: [10, 10, -4] }
      }
      var cfg = configs[preset] || configs.day
      var lights = this.__3d.sceneLights || []
      if (lights[0]) {
        lights[0].intensity = cfg.ambient + (enhanceDepth ? 0.1 : 0)
        if (lights[0].color) lights[0].color.setHex(0xffffff)
      }
      if (lights[1]) {
        lights[1].intensity = cfg.directional + (enhanceDepth ? 0.25 : 0)
        if (lights[1].color) lights[1].color.setHex(cfg.color)
        if (lights[1].position && cfg.position) lights[1].position.set(cfg.position[0], cfg.position[1], cfg.position[2])
      }
      if (lights[2]) {
        lights[2].intensity = enhanceDepth ? 0.45 : 0.3
        if (lights[2].color) lights[2].color.setHex(cfg.fill)
      }
    },

    applyLightSettings: function(lightSettings) {
      if (!this.__3d || !this.__3d.scene) return
      var scene = this.__3d.scene
      var ls = lightSettings || {}
      // 查找命名灯光，回退到按类型过滤
      var ambient = scene.getObjectByName('ambientLight')
      var directional = scene.getObjectByName('directionalLight')
      var fill = scene.getObjectByName('fillLight')
      if (!ambient || !directional) {
        var allLights = scene.children.filter(function(c) { return c.isLight })
        if (!ambient) ambient = allLights[0]
        if (!directional) directional = allLights[1]
        if (!fill) fill = allLights[2]
      }
      // 辅助：CSS #rrggbb 字符串 → hex number
      function parseColorHex(val) {
        if (typeof val === 'string' && val.charAt(0) === '#') {
          return parseInt(val.slice(1), 16)
        }
        return val
      }
      // 环境光
      if (ambient) {
        if (ls.ambientIntensity !== undefined) ambient.intensity = ls.ambientIntensity
        if (ls.ambientColor && ambient.color) ambient.color.setHex(parseColorHex(ls.ambientColor))
      }
      // 方向光
      if (directional) {
        if (ls.directionalIntensity !== undefined) directional.intensity = ls.directionalIntensity
        if (ls.directionalColor && directional.color) directional.color.setHex(parseColorHex(ls.directionalColor))
      }
      // HDRI 环境强度 -> 映射到 fillLight + Bloom + toneMappingExposure
      if (ls.envIntensity !== undefined) {
        if (fill) fill.intensity = 0.3 * ls.envIntensity
        if (this.__3d.renderer) {
          this.__3d.renderer.toneMappingExposure = 0.8 + ls.envIntensity * 0.4
        }
        // Bloom 强度随环境强度联动
        if (this.__3d.composer) {
          var bloomPass = this.__3d.composer.passes.find(function(p) { return p.name === 'bloomPass' })
          if (bloomPass) {
            bloomPass.strength = 0.5 + ls.envIntensity * 1.5
          }
        }
      }
    },

    applyFloorReflection: function(level) {
      if (!this.__3d || !this.__3d.scene) return
      var scene = this.__3d.scene
      if (this.__3d.floorReflectionMesh) {
        scene.remove(this.__3d.floorReflectionMesh)
        if (this.__3d.floorReflectionMesh.geometry) this.__3d.floorReflectionMesh.geometry.dispose()
        if (this.__3d.floorReflectionMesh.material) this.__3d.floorReflectionMesh.material.dispose()
        this.__3d.floorReflectionMesh = null
      }
      if (!level || level === 'none') return
      var opacityMap = { matte: 0.12, low: 0.2, medium: 0.32, strong: 0.48 }
      var roughnessMap = { matte: 0.95, low: 0.72, medium: 0.48, strong: 0.24 }
      var geometry = new THREE.PlaneGeometry(200, 200)
      var material = new THREE.MeshStandardMaterial({
        color: 0xeef7ff,
        transparent: true,
        opacity: opacityMap[level] || 0.12,
        metalness: level === 'strong' ? 0.75 : 0.35,
        roughness: roughnessMap[level] || 0.72,
        side: THREE.DoubleSide
      })
      var mesh = new THREE.Mesh(geometry, material)
      mesh.rotation.x = -Math.PI / 2
      mesh.position.y = -0.012
      mesh.receiveShadow = true
      mesh.userData.isSceneFloor = true
      scene.add(mesh)
      this.__3d.floorReflectionMesh = mesh
    },

    applySceneExtras: function(sceneExtras, gridSettings) {
      if (!this.__3d || !this.__3d.scene || !sceneExtras) return
      var envColors = {
        studio: '#ffffff',
        industrial: '#eef2f5',
        night: '#101827',
        outdoor: '#dbefff',
        sky: '#dbefff',
        sunset: '#ffd5a5',
        ocean: '#c9f0ff',
        forest: '#d9f3df',
        twilight: '#ded6ff'
      }
      var gs = gridSettings || this.gridSettings || {}
      if (!gs.backgroundMode && !gs.backgroundImage) {
        var envKey = gs.environmentPreset || sceneExtras.environment
        var envColor = envColors[envKey] || envColors.studio
        this.__3d.scene.background = new THREE.Color(parseInt(envColor.replace('#', '0x'), 16))
      }
      if (sceneExtras.fogEnabled) {
        var fogColor = sceneExtras.fogColor || (this.gridSettings && this.gridSettings.backgroundColor) || '#d8e4f0'
        this.__3d.scene.fog = new THREE.Fog(parseInt(fogColor.replace('#', '0x'), 16), 30, 120)
        this.__3d.scene.userData.fogEnabledByExtras = true
      } else {
        this.__3d.scene.userData.fogEnabledByExtras = false
        this.__3d.scene.fog = null
      }
      this.applySceneVisualSettings(gs)
    },

    toggleWireframe: function(wireframe, meshMap) {
      if (!this.__3d) return
      Object.values(meshMap).forEach(function(m) {
        if (m.userData && m.userData.isMeshGroup) {
          m.traverse(function(c) { if (c.isMesh && c.material) { c.material.wireframe = wireframe; c.material.needsUpdate = true } })
        } else if (m.material) {
          m.material.wireframe = wireframe
          m.material.needsUpdate = true
        }
      })
    },

    disposeThree: function() {
      if (!this.__3d) return
      var scene = this.__3d.scene
      var meshMap = this.__3d.meshMap || {}

      Object.keys(meshMap).forEach(function(id) {
        var mesh = meshMap[id]
        if (scene && mesh) {
          scene.remove(mesh)
        }
        if (mesh && mesh.traverse) {
          mesh.traverse(function(child) {
            if (child.geometry) child.geometry.dispose()
            if (child.material) {
              if (Array.isArray(child.material)) child.material.forEach(function(m) { m.dispose() })
              else child.material.dispose()
            }
          })
        } else if (mesh) {
          if (mesh.geometry) mesh.geometry.dispose()
          if (mesh.material) {
            if (Array.isArray(mesh.material)) mesh.material.forEach(function(m) { m.dispose() })
            else mesh.material.dispose()
          }
        }
      })

      if (this.__3d.transform) {
        this.__3d.transform.detach()
      }
      if (this.__3d.orbit && this.__3d.orbit.dispose) {
        this.__3d.orbit.dispose()
      }
      if (this.__3d.transform && this.__3d.transform.dispose) {
        this.__3d.transform.dispose()
      }
      if (this.__3d.renderer) {
        this.__3d.renderer.dispose()
      }
      this.overlay2DItems = []
      this.__3d = null
    },

    update2DOverlays: function() {
      if (!this.__3d || !this.$refs.canvasWrapper) {
        this.overlay2DItems = []
        return
      }
      var camera = this.__3d.camera
      var wrapper = this.$refs.canvasWrapper
      var width = wrapper.clientWidth || 0
      var height = wrapper.clientHeight || 0
      var overlays = []

      for (var i = 0; i < this.objects.length; i++) {
        var obj = this.objects[i]
        if (!obj || (obj.type !== '2dComponent' && !isUiOverlayObject(obj))) continue
        var sourceStyle = obj.source2D && obj.source2D.style ? obj.source2D.style : {}
        var styleVisible = obj.type === '2dComponent' ? sourceStyle.visible : obj.visible
        var visible = obj.visible !== false && styleVisible !== 0 && styleVisible !== false
        if (!visible) continue
        var sourcePos = this.getOverlayPosition(obj)
        var x = parseFloat(sourcePos.x) || 0
        var y = parseFloat(sourcePos.y) || 0
        var overlayWidth = Math.max(1, parseFloat(sourcePos.w) || 160)
        var overlayHeight = Math.max(1, parseFloat(sourcePos.h) || 80)

        // position.x/y 是左上角坐标，计算中心点坐标
        var centerX = x + overlayWidth / 2
        var centerY = y + overlayHeight / 2

        if (isNaN(centerX)) centerX = width / 2
        if (isNaN(centerY)) centerY = height / 2

        var transform = parseFloat(this.getOverlayRotation(obj))
        if (isNaN(transform) || transform === -1098 || transform === -1099) transform = 0
        var zIndex = parseInt(sourceStyle.zIndex, 10)
        if (isNaN(zIndex)) zIndex = i + 1
        var hasAction = this.has2DAction(obj)

        var overlay = {
          id: obj.id,
          objectData: obj,
          kind: obj.type,
          style: {
            width: overlayWidth + 'px',
            height: overlayHeight + 'px',
            transform: `translate(-50%, -50%) rotate(${transform}deg)`,
            left: centerX + 'px',
            top: centerY + 'px',
            zIndex: zIndex,
            cursor: 'move',
            display: 'block'
          },
          hasAction: hasAction,
          handlesStyle: {
            width: overlayWidth + 'px',
            height: overlayHeight + 'px',
            left: centerX + 'px',
            top: centerY + 'px',
            zIndex: zIndex,
            transform: `translate(-50%, -50%) rotate(${transform}deg)`
          }
        }
        if (obj.type === 'uiLabel') {
          overlay.text = getTextLabelContent(obj)
          overlay.contentStyle = {
            color: obj.color || '#ffffff',
            background: obj.textBgColor || 'rgba(0,0,0,0.35)',
            opacity: obj.labelOpacity !== undefined ? obj.labelOpacity : 1,
            fontSize: (obj.fontSize || 16) + 'px',
            border: obj.labelShowBorder ? ((obj.labelBorderWidth || 1) + 'px solid ' + (obj.labelBorderColor || '#ffffff')) : 'none'
          }
        } else if (obj.type === 'uiImage') {
          overlay.src = obj.mediaUrl || obj.imageUrl || obj.textureData || ''
        } else if (obj.type === 'webEmbed') {
          overlay.src = obj.webUrl || obj.mediaUrl || 'about:blank'
        }
        overlays.push(overlay)
      }

      this.overlay2DItems = overlays
    },

    has2DAction: function(obj) {
      if (!obj) return false
      var actions = obj.source2D && Array.isArray(obj.source2D.action) ? obj.source2D.action : obj.action
      var activeList = obj.source2D && Array.isArray(obj.source2D.active) ? obj.source2D.active : obj.active
      var hasAction = Array.isArray(actions) && actions.some(function(action) {
        if (!action) return false
        var type = action.type || action.eventType || action.EventType || action.trigger
        var actionName = action.action || action.type
        return !!type && !!actionName && actionName !== 'none'
      })
      var hasActive = Array.isArray(activeList) && activeList.some(function(active) {
        return active && active.isSwitch === true
      })
      return hasAction || hasActive
    },

    getRotateHandleStyle: function(overlay) {
      var obj = overlay.objectData
      var sourceStyle = obj && obj.source2D && obj.source2D.style ? obj.source2D.style : {}
      var transform = parseFloat(sourceStyle.transform)
      if (isNaN(transform) || transform === -1098 || transform === -1099) transform = 0

      // 返回反向旋转变换，使手柄保持水平
      return {
        transform: 'rotate(' + (-transform) + 'deg)'
      }
    },

    /**
     * 加载 GLTF/GLB 模型
     * @param {File} file - 模型文件
     * @param {String} objectId - 对象ID
     * @param {Function} callback - 回调函数，参数是加载后的mesh
     */
    setupGLTFAnimations: function(model, objData, clips) {
      if (!model) return
      model.userData = model.userData || {}
      clips = Array.isArray(clips) ? clips : []
      var nameCounts = {}
      var entries = clips.map(function(clip, index) {
        var name = clip.name || ('Animation ' + (index + 1))
        nameCounts[name] = (nameCounts[name] || 0) + 1
        return {
          key: name + '__' + (index + 1),
          name: name,
          label: name
        }
      })
      var duplicateTotals = {}
      entries.forEach(function(entry) {
        duplicateTotals[entry.name] = nameCounts[entry.name]
      })
      entries.forEach(function(entry, index) {
        if (duplicateTotals[entry.name] > 1) {
          entry.label = entry.name + ' #' + (index + 1)
        }
      })
      var names = entries.map(function(entry) { return entry.key })

      if (!clips.length) {
        model.userData.animationMixer = null
        model.userData.animationActions = {}
        model.userData.animationNames = []
        model.userData.animationEntries = []
        if (objData && !Array.isArray(objData.gltfAnimations)) objData.gltfAnimations = []
        return
      }

      var mixer = new THREE.AnimationMixer(model)
      var actions = {}
      clips.forEach(function(clip, index) {
        var name = entries[index].key
        var action = mixer.clipAction(clip)
        action.enabled = true
        action.paused = true
        actions[name] = action
      })

      model.userData.animationMixer = mixer
      model.userData.animationActions = actions
      model.userData.animationNames = names
      model.userData.animationEntries = entries
      model.userData.currentAnimationName = ''

      if (objData) {
        var savedAnimationName = objData.gltfAnimationName
        var savedAnimationNames = Array.isArray(objData.gltfAnimationNames) ? objData.gltfAnimationNames.slice() : (savedAnimationName ? [savedAnimationName] : [])
        var savedGroups = ensureGLTFAnimationGroups(objData)
        objData.gltfAnimations = entries
        objData.gltfAnimationNames = savedAnimationNames.map(function(savedName) {
          if (actions[savedName]) return savedName
          var matchedEntry = entries.find(function(entry) { return entry.name === savedName || entry.label === savedName })
          return matchedEntry ? matchedEntry.key : ''
        }).filter(function(savedName) { return !!savedName })
        if (!objData.gltfAnimationNames.length && names.length) objData.gltfAnimationNames = [names[0]]
        objData.gltfAnimationName = objData.gltfAnimationNames[0] || ''
        objData.gltfAnimationGroups = remapGLTFAnimationGroupNames(savedGroups, entries, actions, names[0])
        syncLegacyGLTFAnimationFields(objData, objData.gltfAnimationGroups)
      }
      this.$emit('gltf-animations-loaded', {
        objectId: objData && objData.id ? objData.id : model.userData.id,
        animations: entries
      })
    },

    updateGLTFAnimation: function(mesh, obj, delta) {
      if (!mesh || !obj || !mesh.userData || !mesh.userData.animationMixer) return
      var actions = mesh.userData.animationActions || {}
      var names = mesh.userData.animationNames || Object.keys(actions)
      if (!names.length) return

      var groups = ensureGLTFAnimationGroups(obj, { defaultAnimationKey: names[0] })
      var activeMap = {}
      groups.forEach(function(group) {
        if (!group.playing) return
        var selectedNames = Array.isArray(group.animationNames) && group.animationNames.length ? group.animationNames : []
        selectedNames.forEach(function(name) {
          if (!actions[name]) return
          activeMap[name] = group
        })
      })
      Object.keys(actions).forEach(function(name) {
        var group = activeMap[name]
        var action = actions[name]
        if (!group) {
          action.stop()
          return
        }
        var loop = group.loop !== false
        var loopMode = loop ? THREE.LoopRepeat : THREE.LoopOnce
        var speed = group.speed !== undefined && !isNaN(Number(group.speed)) ? Number(group.speed) : 1
        if (!action.isRunning()) {
          action.reset()
          action.play()
        }
        action.setLoop(loopMode, loop ? Infinity : 1)
        action.clampWhenFinished = !loop
        action.timeScale = speed
        action.paused = false
      })
      if (Object.keys(activeMap).length) {
        mesh.userData.animationMixer.update(delta)
      }
    },

    loadGLTFModel: function(file, objectId, callback) {
      if (!this.__3d || !file || !objectId) return
      var loader = new GLTFLoader()
      var reader = new FileReader()
      var self = this

      reader.onload = function(e) {
        var contents = e.target.result
        loader.parse(contents, '', function(gltf) {
          var model = gltf.scene
          var objects = self.objects || []
          var objData = objects.find(function(o) { return o.id === objectId })

          // 计算包围盒，缩放模型到合适大小
          var box = new THREE.Box3().setFromObject(model)
          var size = box.getSize(new THREE.Vector3())
          var maxDim = Math.max(size.x, size.y, size.z)
          var autoScale = maxDim > 0 ? (getGLTFTargetSize(objData) / maxDim) : 1
          model.scale.setScalar(autoScale)

          centerGLTFModelGeometry(model)

          // 遍历模型，设置材质可响应
          model.traverse(function(child) {
            if (child.isMesh) {
              child.castShadow = !!(objData && objData.showShadow)
              child.receiveShadow = !!(objData && objData.showShadow)
              if (child.material) {
                if (Array.isArray(child.material)) {
                  child.material.forEach(function(material) {
                    material.transparent = true
                    material.opacity = 1
                  })
                } else {
                  child.material.transparent = true
                  child.material.opacity = 1
                }
              }
            }
          })

          self.__3d.scene.add(model)

          // 设置ID并存储到meshMap
          model.userData = model.userData || {}
          model.userData.id = objectId
          model.userData.isMeshGroup = true
          self.__3d.meshMap[objectId] = model

          // 将自动缩放值写回 sceneObjects，确保保存/重载时不丢失
          if (objData) {
            objData.sx = autoScale
            objData.sy = autoScale
            objData.sz = autoScale
            model.position.set(objData.x || 0, objData.y || 0, objData.z || 0)
          }

          // 聚焦到模型
          self.setupGLTFAnimations(model, objData, gltf.animations)
          self.updateObjectReflection(objData || { id: objectId, visible: true }, model)
          self.refreshLoadedSelection(objectId)
          self.frameAll(self.objects, self.__3d.meshMap)

          if (callback && typeof callback === 'function') {
            callback(objectId, model)
          }
        }, function(error) {
          console.error('GLTF模型解析失败:', error)
        })
      }

      reader.readAsArrayBuffer(file)
    },

    /**
     * 直接从 ArrayBuffer 加载 GLTF/GLB 模型，用于 base64 恢复场景
     * @param {ArrayBuffer} arrayBuffer - 模型二进制数据
     * @param {String} objectId - 对象ID
     * @param {Function} callback - 回调函数
     */
    loadGLTFModelBuffer: function(arrayBuffer, objectId, callback) {
      if (!this.__3d || !arrayBuffer || !objectId) return
      var loader = new GLTFLoader()
      var self = this
      loader.parse(arrayBuffer, '', function(gltf) {
        var model = gltf.scene

        // 先查 sceneObjects 中是否有已保存的缩放值
        var objects = self.objects || []
        var objData = objects.find(function(o) { return o.id === objectId })

        // 计算自动缩放比例
        var box = new THREE.Box3().setFromObject(model)
        var size = box.getSize(new THREE.Vector3())
        var maxDim = Math.max(size.x, size.y, size.z)
        var autoScale = maxDim > 0 ? (getGLTFTargetSize(objData) / maxDim) : 1

        // 如果 objData 中已有非默认缩放值，说明是重载场景，使用保存的缩放；否则使用自动缩放并写回
        var hasSavedScale = objData && (objData.sx !== undefined && objData.sx !== 1)
        if (hasSavedScale) {
          model.scale.set(objData.sx, objData.sy, objData.sz)
          centerGLTFModelGeometry(model)
          model.position.set(objData.x || 0, objData.y || 0, objData.z || 0)
          model.rotation.set(objData.rx || 0, objData.ry || 0, objData.rz || 0)
        } else {
          model.scale.setScalar(autoScale)
          centerGLTFModelGeometry(model)
          // 写回 sceneObjects
          if (objData) {
            objData.sx = autoScale
            objData.sy = autoScale
            objData.sz = autoScale
            model.position.set(objData.x || 0, objData.y || 0, objData.z || 0)
          }
        }

        // 遍历模型，设置材质
        model.traverse(function(child) {
          if (child.isMesh) {
            child.castShadow = !!(objData && objData.showShadow)
            child.receiveShadow = !!(objData && objData.showShadow)
            // 只有当对象有自定义材质属性时才应用修改，否则保留原始材质
            if (objData && objData.materialOverridden === true && (objData.color || objData.opacity !== undefined || objData.wireframe || objData.metalness !== undefined || objData.roughness !== undefined || objData.emissive || objData.textureData)) {
              if (child.material) {
                if (Array.isArray(child.material)) {
                  child.material.forEach(function(material) {
                    self.applyMaterialProps(material, objData)
                  })
                } else {
                  self.applyMaterialProps(child.material, objData)
                }
              }
            }
          }
        })
        // 若 meshMap 中已有同 ID 的旧 mesh，先移除
        if (self.__3d.meshMap[objectId]) {
          self.removeObjectReflection(objectId)
          self.__3d.scene.remove(self.__3d.meshMap[objectId])
          delete self.__3d.meshMap[objectId]
        }
        self.__3d.scene.add(model)
        model.userData = model.userData || {}
        model.userData.id = objectId
        model.userData.isMeshGroup = true
        self.__3d.meshMap[objectId] = model
        model.visible = objData ? (objData.visible !== false) : true
        self.setupGLTFAnimations(model, objData, gltf.animations)
        self.updateObjectReflection(objData || { id: objectId, visible: true }, model)
        self.refreshLoadedSelection(objectId)

        if (callback && typeof callback === 'function') {
          callback(objectId, model)
        }
      }, function(error) {
        console.error('GLTF模型(buffer)解析失败:', error)
      })
    },

    loadGLTFModelUrl: function(url, objectId, callback, options) {
      if (!this.__3d || !url || !objectId) return
      options = options || {}
      var loader = new GLTFLoader()
      var self = this
      loader.load(url, function(gltf) {
        var model = gltf.scene
        var objects = self.objects || []
        var objData = objects.find(function(o) { return o.id === objectId })

        var box = new THREE.Box3().setFromObject(model)
        var size = box.getSize(new THREE.Vector3())
        var maxDim = Math.max(size.x, size.y, size.z)
        var autoScale = maxDim > 0 ? (getGLTFTargetSize(objData) / maxDim) : 1

        var hasSavedScale = objData && (objData.sx !== undefined && objData.sx !== 1)
        if (hasSavedScale) {
          model.scale.set(objData.sx, objData.sy, objData.sz)
          centerGLTFModelGeometry(model)
          model.position.set(objData.x || 0, objData.y || 0, objData.z || 0)
          model.rotation.set(objData.rx || 0, objData.ry || 0, objData.rz || 0)
        } else {
          model.scale.setScalar(autoScale)
          centerGLTFModelGeometry(model)
          if (objData) {
            objData.sx = autoScale
            objData.sy = autoScale
            objData.sz = autoScale
            model.position.set(objData.x || 0, objData.y || 0, objData.z || 0)
          }
        }

        model.traverse(function(child) {
          if (child.isMesh) {
            child.castShadow = !!(objData && objData.showShadow)
            child.receiveShadow = !!(objData && objData.showShadow)
            // 只有当对象有自定义材质属性时才应用修改，否则保留原始材质
            if (objData && objData.materialOverridden === true && (objData.color || objData.opacity !== undefined || objData.wireframe || objData.metalness !== undefined || objData.roughness !== undefined || objData.emissive || objData.textureData)) {
              if (child.material) {
                if (Array.isArray(child.material)) {
                  child.material.forEach(function(material) {
                    self.applyMaterialProps(material, objData)
                  })
                } else {
                  self.applyMaterialProps(child.material, objData)
                }
              }
            }
          }
        })
        if (self.__3d.meshMap[objectId]) {
          self.removeObjectReflection(objectId)
          self.__3d.scene.remove(self.__3d.meshMap[objectId])
          delete self.__3d.meshMap[objectId]
        }
        self.__3d.scene.add(model)
        model.userData = model.userData || {}
        model.userData.id = objectId
        model.userData.isMeshGroup = true
        self.__3d.meshMap[objectId] = model
        model.visible = objData ? (objData.visible !== false) : true
        self.setupGLTFAnimations(model, objData, gltf.animations)
        self.updateObjectReflection(objData || { id: objectId, visible: true }, model)
        self.refreshLoadedSelection(objectId)

        if (callback && typeof callback === 'function') {
          callback(objectId, model)
        }
      }, undefined, function(error) {
        if (!options.silent) {
          console.error('GLTF模型(URL)加载失败:', error)
        }
      })
    },

    fmt: function(v) { return typeof v === 'number' ? v.toFixed(2) : v }
  }
}
</script>

<style scoped>
.center-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: #ffffff;
  position: relative;
  overflow: hidden;
  width: 100%;
  min-height: 0;
}
.canvas-toolbar {
  height: 36px;
  background: #fff;
  border-bottom: 1px solid #e8e8e8;
  display: flex;
  align-items: center;
  padding: 0 10px;
  gap: 4px;
  flex-shrink: 0;
}
.canvas-toolbar .toolbar-btn { padding: 3px 8px; font-size: 12px; }
.canvas-info {
  margin-left: auto;
  font-size: 11px;
  color: #999;
  display: flex;
  gap: 12px;
}
#three-canvas-wrapper {
  flex: 1;
  display: flex;
  min-height: 0;
  overflow: hidden;
  position: relative;
}
.overlay-2d-layer {
  position: absolute;
  inset: 0;
  z-index: 2;
  pointer-events: none;
  overflow: hidden;
}

.overlay-2d-item {
  position: absolute;
  transform-origin: center center;
  pointer-events: auto;
  cursor: move;
  box-sizing: border-box;
}

.overlay-2d-content {
  width: 100%;
  height: 100%;
}

.overlay-2d-item.has-action .overlay-2d-content {
  cursor: pointer;
}

.basic-ui-label {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-sizing: border-box;
  padding: 6px 10px;
  overflow: hidden;
  word-break: break-word;
}

.basic-ui-image {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: contain;
  background: rgba(17, 24, 39, 0.18);
  border: 1px solid rgba(255, 255, 255, 0.55);
  box-sizing: border-box;
}

.basic-ui-image-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(255, 255, 255, 0.82);
  font-size: 22px;
}

.basic-ui-frame {
  width: 100%;
  height: 100%;
  display: block;
  border: 1px solid rgba(17, 24, 39, 0.2);
  background: #fff;
  box-sizing: border-box;
  pointer-events: none;
}

.overlay-2d-item.selected {
}

.overlay-2d-selection-box {
  position: absolute;
  top: -4px;
  left: -4px;
  right: -4px;
  bottom: -4px;
  border: 2px solid #13c2c2;
  pointer-events: none;
}

.overlay-2d-rotate-handle,
.overlay-2d-resize-handle {
  position: absolute;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: #fff;
  border: 1px solid #13c2c2;
  box-shadow: 0 0 0 1px rgba(19, 194, 194, 0.3);
}

.overlay-2d-rotate-handle {
  position: absolute;
  left: 50%;
  top: -16px;
  transform: translateX(-50%);
  cursor: grab;
  width: 14px;
  height: 14px;
  background: #13c2c2;
  border: 2px solid #fff;
  box-shadow: 0 0 0 2px rgba(19, 194, 194, 0.5);
  z-index: 10;
}

.overlay-2d-resize-handle {
  right: -5px;
  bottom: -5px;
  cursor: nwse-resize;
}

.overlay-2d-resize-handle-tl {
  top: -5px;
  left: -5px;
  cursor: nwse-resize;
}

.overlay-2d-resize-handle-tr {
  top: -5px;
  right: -5px;
  cursor: nesw-resize;
}

.overlay-2d-resize-handle-bl {
  bottom: -5px;
  left: -5px;
  cursor: nesw-resize;
}
#three-canvas {
  flex: 1;
  display: block;
  width: 100%;
  min-height: 0;
  cursor: default;
}
.canvas-overlay-hint { display: none; }
.drop-hint {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  text-align: center;
  pointer-events: none;
}
.drop-hint i { font-size: 48px; color: #d9d9d9; }
.drop-hint p { color: #bbb; font-size: 14px; margin-top: 10px; }
.axis-indicator {
  position: absolute;
  bottom: 50px;
  right: 14px;
  width: 70px;
  height: 70px;
  pointer-events: none;
}
.selected-indicator {
  position: absolute;
  top: 8px;
  left: 50%;
  transform: translateX(-50%);
  background: linear-gradient(135deg, #13c2c2 0%, #36cfc9 100%);
  color: white;
  padding: 4px 16px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
  pointer-events: none;
  box-shadow: 0 2px 12px rgba(19, 194, 194, 0.35);
  z-index: 10;
}

.transform-selection-frame {
  position: absolute;
  z-index: 11;
  pointer-events: none;
}

.kf-border {
  position: absolute;
  inset: 0;
  border: 1px dashed #00bcd4;
}

.kf-c {
  position: absolute;
  width: 12px;
  height: 12px;
  border-color: #00bcd4;
  border-style: solid;
  pointer-events: auto;
  background: rgba(0, 188, 212, 0.08);
}

.kf-c.kf-tl { top: -5px; left: -5px; border-width: 2px 0 0 2px; cursor: nwse-resize; }
.kf-c.kf-tr { top: -5px; right: -5px; border-width: 2px 2px 0 0; cursor: nesw-resize; }
.kf-c.kf-bl { bottom: -5px; left: -5px; border-width: 0 0 2px 2px; cursor: nesw-resize; }
.kf-c.kf-br { bottom: -5px; right: -5px; border-width: 0 2px 2px 0; cursor: nwse-resize; }
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
  transition: all 0.15s;
  white-space: nowrap;
}
.toolbar-btn:hover { background: #e6fffb; border-color: #87e8de; color: #13c2c2; }
.toolbar-btn.active { background: #e6fffb; border-color: #13c2c2; color: #13c2c2; }
.toolbar-btn[disabled] { opacity: 0.4; cursor: not-allowed; }
.toolbar-btn[disabled]:hover { background: transparent; border-color: transparent; color: #333; }
.toolbar-divider {
  width: 1px;
  height: 24px;
  background: #e8e8e8;
  margin: 0 6px;
}
.status-bar { display: none; }
.status-bar span { display: flex; align-items: center; gap: 4px; }
</style>
