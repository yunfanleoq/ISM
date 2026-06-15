﻿﻿<template>
  <div class="scene-preview-page">
    <PreviewWatermark />


    <!-- 3D 画布 -->
    <div class="preview-canvas-wrapper" ref="canvasWrapper">
      <canvas id="preview-canvas" ref="previewCanvas"></canvas>
      <div class="overlay-2d-layer">
        <div
          v-for="overlay in overlay2DItems"
          :key="overlay.id"
          class="overlay-2d-item"
          :class="{ 'has-action': overlay.hasAction }"
          :style="overlay.style"
        >
          <div
            class="overlay-2d-content"
            @click.capture="executeInteraction(overlay.objectData, 'click', $event)"
            @dblclick.capture="executeInteraction(overlay.objectData, 'dblclick', $event)"
            @mousedown.capture="executeInteraction(overlay.objectData, 'mousedown', $event)"
            @mouseup.capture="executeInteraction(overlay.objectData, 'mouseup', $event)"
            @mouseenter="executeInteraction(overlay.objectData, 'mouseenter', $event)"
            @mouseleave="executeInteraction(overlay.objectData, 'mouseleave', $event)"
          >
            <ISM2DNodeRenderer
              v-if="overlay.kind === '2dComponent'"
              :object-data="overlay.objectData"
              :showDeviceUuid="$route.query.deviceKey || ''"
              :editMode="false"
              :selected="false"
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
        </div>
      </div>

      <!-- 空场景提示 -->
      <div class="empty-hint" v-if="objectCount === 0">
        <i class="fas fa-cube" style="font-size:48px;color:#d9d9d9"></i>
        <p style="color:#bbb;margin-top:12px">{{ $t('ISM3DEditor.sceneEmpty') }}</p>
      </div>

      <!-- 加载中 -->
      <div class="loading-mask" v-if="loading">
        <a-spin size="large" />
        <p style="margin-top:12px;color:#666">{{ $t('ISM3DEditor.loadingScene') }}</p>
      </div>

      <div v-if="popupInteraction.visible" class="preview-popup-mask" @click="handlePopupMaskClick">
        <div class="preview-popup-window" :style="{ width: popupInteraction.width + 'px', height: popupInteraction.height + 'px' }" @click.stop>
          <div class="preview-popup-header">
            <span>{{ popupInteraction.title || $t('ISM3DEditor.popupPage') }}</span>
            <button class="preview-popup-close" @click="closePopupInteraction">x</button>
          </div>
          <iframe class="preview-popup-frame" :src="popupInteraction.url"></iframe>
        </div>
      </div>

      <!-- 右键菜单 -->
      <div
        v-if="contextMenu.visible"
        class="context-menu"
        :style="{ left: contextMenu.x + 'px', top: contextMenu.y + 'px' }"
        @click.stop
      >
        <div class="context-menu-title">{{ $t('ISM3DEditor.cameraView') }}</div>
        <button class="context-menu-item" :class="{active:cameraView==='perspective'}" @click="handleContextMenuClick('setCameraView', 'perspective')">
          <i class="fas fa-camera"></i> {{ $t('ISM3DEditor.perspective') }}
        </button>
        <button class="context-menu-item" :class="{active:cameraView==='top'}" @click="handleContextMenuClick('setCameraView', 'top')">
          <i class="fas fa-arrow-down"></i> {{ $t('ISM3DEditor.top') }}
        </button>
        <button class="context-menu-item" :class="{active:cameraView==='front'}" @click="handleContextMenuClick('setCameraView', 'front')">
          <i class="fas fa-arrow-right"></i> {{ $t('ISM3DEditor.front') }}
        </button>
        <button class="context-menu-item" :class="{active:cameraView==='right'}" @click="handleContextMenuClick('setCameraView', 'right')">
          <i class="fas fa-arrow-left"></i> {{ $t('ISM3DEditor.side') }}
        </button>
        <div class="context-menu-divider"></div>
        <button class="context-menu-item" @click="handleContextMenuClick('resetCamera')">
          <i class="fas fa-home"></i> {{ $t('ISM3DEditor.reset') }}
        </button>
        <button class="context-menu-item" @click="handleContextMenuClick('frameAll')">
          <i class="fas fa-compress-arrows-alt"></i> {{ $t('ISM3DEditor.frame') }}
        </button>
        <div class="context-menu-divider"></div>
        <button class="context-menu-item" :class="{active:showGrid}" @click="handleContextMenuClick('toggleGrid')">
          <i class="fas fa-border-all"></i> {{ $t('ISM3DEditor.grid') }}
        </button>
        <button class="context-menu-item" :class="{active:autoRotate}" @click="handleContextMenuClick('toggleAutoRotate')">
          <i class="fas fa-sync"></i> {{ $t('ISM3DEditor.autoRotate') }}
        </button>
        <div class="context-menu-divider"></div>
        <button class="context-menu-item" @click="handleContextMenuClick('toggleFullscreen')">
          <i class="fas fa-expand"></i> {{ isFullscreen ? $t('ISM3DEditor.exitFullscreen') : $t('ISM3DEditor.fullscreen') }}
        </button>
        <button class="context-menu-item" @click="handleContextMenuClick('fitPage')">
          <i class="fas fa-compress"></i> {{ $t('ISM3DEditor.fitPage') }}
        </button>
      </div>
    </div>
    <a-modal :visible="settingDialog" :footer="null" :maskClosable="false" :destroyOnClose="true" @cancel="closeSetDialogs" :title="$t('monitor.Set')" :width="320">
      <div style="padding: 5px;">
        <a-form>
          <a-form-item v-if="SetPassword" :label="$t('displayConfig.Properties.SetPassword')">
            <a-input v-model="SetPasswordFormValue" :style="{ 'text-security': 'disc', '-webkit-text-security': 'disc' }" />
          </a-form-item>
          <a-form-item :label="$t('monitor.SetValue')">
            <a-textarea v-model="SetValueFormValue" auto-size />
          </a-form-item>
        </a-form>
      </div>
      <a-divider />
      <div class="dialog-button">
        <a-button type="primary" :loading="settingLoading" @click="ManualSetData">{{ $t('component.deviceDataModel.submit') }}</a-button>
        <a-button style="margin-left: 10px" @click="closeSetDialogs">{{ $t('component.deviceDataModel.cancel') }}</a-button>
      </div>
    </a-modal>

    <a-modal :visible="setPasswordDialog" :footer="null" :maskClosable="false" :destroyOnClose="true" @cancel="closeSetDialogs" :title="$t('displayConfig.Properties.SetPasswordTips')" :width="320">
      <div style="padding: 5px;">
        <a-form>
          <a-form-item :label="$t('displayConfig.Properties.SetPassword')">
            <a-input v-model="SetAutoPasswordFormValue" :style="{ 'text-security': 'disc', '-webkit-text-security': 'disc' }" />
          </a-form-item>
        </a-form>
      </div>
      <a-divider />
      <div class="dialog-button">
        <a-button type="primary" :loading="settingLoading" @click="PasswordSetData">{{ $t('component.deviceDataModel.submit') }}</a-button>
        <a-button style="margin-left: 10px" @click="closeSetDialogs">{{ $t('component.deviceDataModel.cancel') }}</a-button>
      </div>
    </a-modal>

    <a-modal :visible="actionPasswordDialog" :footer="null" :maskClosable="false" :destroyOnClose="true" @cancel="closeActionPasswordDialog" :title="$t('displayConfig.Properties.SetPasswordTips')" :width="320">
      <div style="padding: 5px;">
        <a-form>
          <a-form-item :label="$t('displayConfig.Properties.SetPassword')">
            <a-input v-model="ActionPasswordValue" :style="{ 'text-security': 'disc', '-webkit-text-security': 'disc' }" />
          </a-form-item>
        </a-form>
      </div>
      <a-divider />
      <div class="dialog-button">
        <a-button type="primary" :loading="settingLoading" @click="PasswordSetAction">{{ $t('component.deviceDataModel.submit') }}</a-button>
        <a-button style="margin-left: 10px" @click="closeActionPasswordDialog">{{ $t('component.deviceDataModel.cancel') }}</a-button>
      </div>
    </a-modal>
  </div>
</template>

<script>
import * as THREE from 'three'
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls'
import { GLTFLoader } from 'three/examples/jsm/loaders/GLTFLoader'
import { request } from '@/utils/request'
import { message } from 'ant-design-vue'
import { mapGetters } from 'vuex'
import { setData } from '@/services/device'
import { ComponentRestApi } from '@/services/RestApi'
import { ExecSysScript } from '@/services/system'
import ISM2DNodeRenderer from '../Scene/ISM2DNodeRenderer.vue'
import PreviewWatermark from '@/components/PreviewWatermark.vue'
import {
  createThreeMesh as createThreeMeshFromLib,
  create2DComponentPlane,
  createMediaPlane,
  updateMediaPlane,
  disposeMediaPlane,
  createTextSprite,
  updateTextSprite,
  createText3DMesh,
  updateText3DMesh,
  createFlowPipe,
  updateFlowPipeAnimation,
  updateFlowPipe,
  isFlowPipeSegmented,
  tryLoadChineseFontFor3D,
  COMP_COLORS
} from '../Objects/IndustrialObjects'
import { parseScenePayloadOrEmpty } from '../../utils/scenePayload'
import { ensureGLTFAnimationGroups, remapGLTFAnimationGroupNames, syncLegacyGLTFAnimationFields } from '../../utils/GLTFAnimationGroups'
import { ensurePositionBindings } from '../../utils/positionBindings'
import { defaultGridSettings } from '../../utils/sceneSettings'
import {getDisplayModelLayerData} from "@/services/displayModel";
import { getISMComponent } from '@/pages/ISMDisPlay/componentRegistry'

/** 将任意颜色值转为 '#rrggbb' 或 '#rrggbbaa' 格式 */
function normColorValue(val, fallback) {
  fallback = fallback || '#000000'
  if (!val || typeof val !== 'string') return fallback
  var s = val.trim()
  var COLOR_KEYS = {
    blue: '#4da6ff', green: '#4dffa6', orange: '#ffaa4d',
    purple: '#c17bff', pink: '#ff7bb5', cyan: '#4dfff0', yellow: '#fffb4d',
    red: '#ff4d4f', gray: '#8c8c8c', brown: '#a67c52'
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
  const bgOpacity = obj.textBgOpacity !== undefined ? Number(obj.textBgOpacity) : undefined
  const hasBackground = isFinite(bgOpacity) && bgOpacity > 0
  return {
    opacity: obj.labelOpacity !== undefined ? Number(obj.labelOpacity) : 1,
    backgroundOpacity: isFinite(bgOpacity) ? bgOpacity : undefined,
    faceCamera: obj.labelFaceCamera !== false,
    fixedSize: obj.labelFixedSize !== false,
    fontFamily: obj.labelFontFamily || '系统默认',
    renderMode: obj.labelRenderMode || 'component',
    true3D: !hasBackground,
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
  let sx = obj.sx !== undefined ? Number(obj.sx) : 1
  let sy = obj.sy !== undefined ? Number(obj.sy) : 1
  let sz = obj.sz !== undefined ? Number(obj.sz) : 1
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
  let localBox = mesh.geometry && mesh.geometry.boundingBox ? mesh.geometry.boundingBox.clone() : null
  if (!localBox && mesh.geometry && mesh.geometry.computeBoundingBox) {
    mesh.geometry.computeBoundingBox()
    localBox = mesh.geometry.boundingBox ? mesh.geometry.boundingBox.clone() : null
  }
  if (!localBox) localBox = new THREE.Box3().setFromObject(mesh)
  const size = localBox.getSize(new THREE.Vector3())
  const center = localBox.getCenter(new THREE.Vector3())
  const width = Math.max(Math.abs(size.x), 0.2)
  const height = Math.max(Math.abs(size.y), 0.12)
  const depth = Math.max(Math.abs(size.z), 0.02)
  let pick = mesh.userData.textPickArea
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

function isUiOverlayObject(obj) {
  return obj && (obj.type === 'uiLabel' || obj.type === 'uiImage' || obj.type === 'webEmbed')
}

function getTextLabelContent(obj) {
  if (!obj) return ''
  if (obj.type === 'dataText') {
    const value = obj.realTimeValue !== undefined && obj.realTimeValue !== '' ? obj.realTimeValue : '--'
    const format = obj.dataFormat && obj.dataFormat !== '{value}' ? obj.dataFormat : '{value}'
    return format.replace('{value}', value)
  }
  if (obj.type === 'uiLabel') return obj.textContent || 'UI标签'
  if (obj.type === 'textPlain3d') return obj.textContent || '3D文字'
  return obj.textContent || '标签'
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
  name: 'ScenePreview',
  components: {
    ISM2DNodeRenderer,
    PreviewWatermark
  },

  data() {
    return {
      loading: false,
      sceneObjects: [],
      objectCount: 0,
      fps: 0,
      cameraView: 'perspective',
      showGrid: true,
      autoRotate: false,
      popupInteraction: {
        visible: false,
        url: '',
        title: '',
        width: 960,
        height: 640,
        autoClose: false,
      },
      _savedSceneSettings: defaultGridSettings(),
      sceneExtras: {
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
      },
      _fpsFrames: 0,
      _fpsLast: 0,
      _animId: null,
      overlay2DItems: [],
      contextMenu: {
        visible: false,
        x: 0,
        y: 0
      },
      isFullscreen: false,
      settingDialog: false,
      setPasswordDialog: false,
      actionPasswordDialog: false,
      settingLoading: false,
      SetPassword: '',
      SetPasswordFormValue: '',
      SetAutoPasswordFormValue: '',
      SetValueFormValue: '',
      ActionPasswordValue: '',
      ActionPasswordSet: '',
      _pendingActionRunner: null,
      _pendingSetValueResolve: null,
      _pendingSetValueItem: null,
      _pendingSetValueDeviceUuid: '',
    }
  },
  i18n: require('@/i18n/language'),
  computed: {
    ...mapGetters('account', ['user'])
  },
  watch: {
    '$route.fullPath'() {
      this.loadScene()
    }
  },
  mounted() {
    this.$nextTick(() => {
      this.initThree()
      this.startRenderLoop()
      // 等待中文字体加载完成后再渲染场景，保证 true 3D 文字能正确生成
      const fontPromise = (typeof window !== 'undefined' && window.opentype && tryLoadChineseFontFor3D)
        ? tryLoadChineseFontFor3D()
        : Promise.resolve(null)
      fontPromise.then(() => {
        this.loadScene()
      }).catch(() => {
        this.loadScene()
      })
      // 监听字体异步加载完成事件，重建 True3D 文字以应用新字体
      this._onFontReady = this._upgradeCanvasFallbackTexts.bind(this)
      if (typeof window !== 'undefined') {
        window.addEventListener('ism3d-font-ready', this._onFontReady)
      }
    })
    // 监听数据推送，驱动动画
    this.$EventBus.$on("readDataPush", this.DealWithUpdateData)
    this.$EventBus.$on("StaticData", this.DealWithUpdateData)
    this.$EventBus.$on("SystemData", this.DealWithUpdateData)
  },

  beforeDestroy() {
    this.$EventBus.$off("readDataPush", this.DealWithUpdateData)
    this.$EventBus.$off("StaticData", this.DealWithUpdateData)
    this.$EventBus.$off("SystemData", this.DealWithUpdateData)
    window.removeEventListener('resize', this.onResize)
    document.removeEventListener('click', this.hideContextMenu)
    document.removeEventListener('fullscreenchange', this.onFullscreenChange)
    if (typeof window !== 'undefined' && this._onFontReady) {
      window.removeEventListener('ism3d-font-ready', this._onFontReady)
    }
    if (this.$refs.previewCanvas) {
      this.$refs.previewCanvas.removeEventListener('click', this.onCanvasClick)
      this.$refs.previewCanvas.removeEventListener('dblclick', this.onCanvasDblClick)
      this.$refs.previewCanvas.removeEventListener('mousedown', this._handleMousedown)
      this.$refs.previewCanvas.removeEventListener('mousemove', this.onCanvasMouseMove)
      this.$refs.previewCanvas.removeEventListener('contextmenu', this.showContextMenu)
    }
    if (this._animId) cancelAnimationFrame(this._animId)
    if (this.__3d) {
      Object.keys(this.__3d.meshMap).forEach(id => {
        const mesh = this.__3d.meshMap[id]
        this.__3d.scene.remove(mesh)
        if (mesh.traverse) {
          mesh.traverse(child => {
            if (child.geometry) child.geometry.dispose()
            if (child.material) {
              if (Array.isArray(child.material)) child.material.forEach(m => m.dispose())
              else child.material.dispose()
            }
          })
        } else {
          if (mesh.geometry) mesh.geometry.dispose()
          if (mesh.material) mesh.material.dispose()
        }
      })
      if (this.__3d.orbit && this.__3d.orbit.dispose) {
        this.__3d.orbit.dispose()
      }
      this.__3d.renderer.dispose()
      this.__3d = null
    }
  },

  methods: {
    goBack() {
      this.$router.go(-1)
    },

    initThree() {
      const vm = this
      const canvas = vm.$refs.previewCanvas
      const wrapper = vm.$refs.canvasWrapper
      if (!wrapper) return
      const rect = wrapper.getBoundingClientRect()
      const w = rect.width || 800
      const h = rect.height || 600

      const scene = new THREE.Scene()
      scene.background = new THREE.Color(0xffffff)
      scene.fog = new THREE.Fog(0xffffff, 40, 100)

      const camera = new THREE.PerspectiveCamera(45, w / h, 0.1, 1000)
      camera.position.set(6, 5, 8)
      camera.lookAt(0, 0, 0)

      const renderer = new THREE.WebGLRenderer({ canvas, antialias: true })
      renderer.setSize(w, h)
      renderer.setPixelRatio(window.devicePixelRatio)
      renderer.shadowMap.enabled = true
      renderer.shadowMap.type = THREE.PCFSoftShadowMap
      renderer.outputEncoding = THREE.sRGBEncoding
      renderer.toneMapping = THREE.ACESFilmicToneMapping
      renderer.toneMappingExposure = 1.0

      scene.add(new THREE.AmbientLight(0xffffff, 1.0))
      const dirLight = new THREE.DirectionalLight(0xffffff, 0.9)
      dirLight.position.set(8, 12, 6)
      dirLight.castShadow = true
      dirLight.shadow.mapSize.set(2048, 2048)
      scene.add(dirLight)
      scene.add(new THREE.DirectionalLight(0x4466aa, 0.3))

      const gridHelper = new THREE.GridHelper(20, 20, 0x1890ff, 0xcccccc)
      scene.add(gridHelper)

      const orbit = new OrbitControls(camera, canvas)
      orbit.enableDamping = true
      orbit.dampingFactor = 0.06
      orbit.minDistance = 1
      orbit.maxDistance = 100
      orbit.mouseButtons = { LEFT: THREE.MOUSE.ROTATE, MIDDLE: THREE.MOUSE.DOLLY, RIGHT: THREE.MOUSE.PAN }

      const clock = new THREE.Clock()
      const raycaster = new THREE.Raycaster()
      const mouse = new THREE.Vector2()
      const meshMap = {}

      vm.__3d = { scene, camera, renderer, orbit, gridHelper, clock, meshMap, raycaster, mouse }
      vm.__3d.sceneLights = scene.children.filter(child => child.isLight)

      vm._handleMousedown = function() { if (vm.__3d) vm.__3d._cameraAnim = null }

      window.addEventListener('resize', vm.onResize)
      canvas.addEventListener('click', vm.onCanvasClick)
      canvas.addEventListener('dblclick', vm.onCanvasDblClick)
      canvas.addEventListener('mousedown', vm._handleMousedown)
      canvas.addEventListener('mousemove', vm.onCanvasMouseMove)
      canvas.addEventListener('contextmenu', vm.showContextMenu)
      document.addEventListener('click', vm.hideContextMenu)
      document.addEventListener('fullscreenchange', vm.onFullscreenChange)
    },

    startRenderLoop() {
      const vm = this
      let elapsed = 0
      const render = function() {
        vm._animId = requestAnimationFrame(render)
        if (!vm.__3d) return

        const delta = vm.__3d.clock.getDelta()
        elapsed += delta

        // ===== 移植 SceneCanvas.vue 的动画逻辑 =====
        const objs = vm.sceneObjects
        const meshMap = vm.__3d.meshMap
        for (let i = 0; i < objs.length; i++) {
          const obj = objs[i]
          const mesh = meshMap[obj.id]
          if (!mesh || !obj) continue

          // autoRotate：自旋动画（用 delta，帧率无关）
          if (obj.autoRotate) {
            const speed = obj.rotateSpeed !== undefined ? obj.rotateSpeed : 1
            const axis = obj.rotateAxis || 'y'
            if (mesh.userData && mesh.userData.hasRotorAnimation && Array.isArray(mesh.userData.rotors)) {
              mesh.userData.rotors.forEach(function(rotor) {
                rotor.rotation.y -= speed * delta * 8
              })
            } else {
              mesh.rotation[axis] += speed * delta
            }
          }

          vm.updateGLTFAnimation(mesh, obj, delta)

          // floatAnim：浮动动画（用 elapsed，持续振荡）
          if (obj.floatAnim) {
            const floatRange = obj.floatRange !== undefined ? obj.floatRange : 0.15
            const floatSpeed = obj.floatSpeed !== undefined ? obj.floatSpeed : 2
            const floatOffset = Math.sin(elapsed * floatSpeed + i * 1.5) * floatRange
            if (obj._baseY === undefined) obj._baseY = obj.y
            mesh.position.y = obj._baseY + floatOffset
          } else {
            if (obj._baseY !== undefined) {
              mesh.position.y = obj._baseY
              delete obj._baseY
            }
          }

          // blink：闪烁动画（支持单Mesh和MeshGroup）
          if (obj.blink) {
            const blinkSpeed = obj.blinkSpeed !== undefined ? obj.blinkSpeed : 6
            const blinkMin = obj.blinkMin !== undefined ? obj.blinkMin : 0.2
            const blinkVal = (Math.sin(elapsed * blinkSpeed) + 1) * 0.5
            const targetOpacity = blinkMin + blinkVal * (1 - blinkMin)

            if (mesh.material) {
              const mats = Array.isArray(mesh.material) ? mesh.material : [mesh.material]
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
                  const childMats = Array.isArray(child.material) ? child.material : [child.material]
                  childMats.forEach(function(mat) {
                    if (mat.isShaderMaterial) return
                    mat.opacity = targetOpacity
                    mat.transparent = true
                    mat.needsUpdate = true
                  })
                }
              })
            }
          }

          // Shader 特效组件：更新 uniforms.time（flowRibbon/energyWall/pulseColumn 等）
          if (mesh.userData && mesh.userData.hasAnimation) {
            mesh.traverse(function(child) {
              if (child.isMesh && child.material && child.material.uniforms && child.material.uniforms.time !== undefined) {
                child.material.uniforms.time.value = elapsed
              }
            })
          }

          // flowPipe：流动管道动画
          if (obj.type === 'flowPipe' && mesh.userData && mesh.userData.isFlowPipe) {
            var flowCondition = obj.flowAnim !== false
            if (obj.animate && obj.animate.condition && obj.animate.condition.isBandDevice) {
              flowCondition = flowCondition && vm.checkAnimateCondition(obj)
            }
            updateFlowPipeAnimation(mesh, delta, flowCondition)
            updateFlowPipe(mesh, obj)
          }
        }
        // ===== 动画逻辑结束 =====

        vm.applyTimelineTracks(elapsed)

        // 自动旋转（相机自动旋转，非物体自旋）
        if (vm.autoRotate) {
          vm.__3d.orbit.autoRotate = true
          vm.__3d.orbit.autoRotateSpeed = 2.0
        } else {
          vm.__3d.orbit.autoRotate = false
        }

        // 相机移动动画
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

        vm.__3d.orbit.update()
        vm.__3d.renderer.render(vm.__3d.scene, vm.__3d.camera)
        vm.update2DOverlays()

        vm._fpsFrames++
        const now = performance.now()
        if (now - vm._fpsLast >= 1000) {
          vm.fps = vm._fpsFrames
          vm._fpsFrames = 0
          vm._fpsLast = now
        }
      }
      render()
    },

    async loadScene() {
      const uid = this.$route && this.$route.params && this.$route.params.uid
      if (!uid) return

      this.loading = true
      try {
        let params = {
          muid:uid,
        }
        const res = await getDisplayModelLayerData(params)
        if(res.data.code==0)
        {
          const layers = Array.isArray(res.data.layer) ? res.data.layer : []
          const query = (this.$route && this.$route.query) || {}
          const pageId = query.pageId || query.pageid || query.page || query.pageUuid
          let pageLayer = null
          if (pageId) {
            pageLayer = layers.find(layer => String(layer.PageId || layer.pageUuid || layer.PageUuid || layer.id) === String(pageId))
          }
          pageLayer = pageLayer || layers.find(layer => layer.IsHome === 1 || layer.IsHome === '1' || layer.IsHome === true || layer.isHome === true) || layers[0]
          if (!pageLayer) return
          try {
            let dataJson = parseScenePayloadOrEmpty(pageLayer && pageLayer.components)
            // 校验文件格式
            if (!dataJson.objects || !Array.isArray(dataJson.objects)) {
              message.error(this.$t('ISM3DEditor.invalidSceneFormat'))
              return
            }
            // 校验每个对象必须包含 id 和 type
            const valid = dataJson.objects.every(obj => obj.id && obj.type)
            if (!valid) {
              message.error(this.$t('ISM3DEditor.invalidSceneObjectFields'))
              return
            }
            this.normalizeSceneObjects(dataJson.objects)
            this.sceneObjects = dataJson.objects
            this._savedSceneSettings = dataJson.sceneSettings || null
            this.sceneExtras = Object.assign({}, this.sceneExtras, dataJson.sceneExtras || {})
            if (dataJson.sceneSettings) {
              this.applySceneSettings(dataJson.sceneSettings)
            }
            this.applySceneExtras(this.sceneExtras)
            this.$nextTick(() => {
              this.syncMeshes(this.sceneObjects)
              this.restoreGLTFModels(this.sceneObjects)
            })
          } catch (parseErr) {
            console.error("场景数据解析失败:", parseErr)
            message.error(this.$t('ISM3DEditor.sceneParseFailed') + parseErr.message)
          }
        }
      } catch (err) {
        console.error("加载场景失败:", err)
        message.error(this.$t('ISM3DEditor.sceneLoadFailed'))
      } finally {
        this.loading = false
      }
    },

    /** 将任意颜色值转为 '#rrggbb' 格式，失败返回 fallback */
    normColorValue(val, fallback) {
      fallback = fallback || '#000000'
      if (!val || typeof val !== 'string') return fallback
      var s = val.trim()
      // colorKey → 查表
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
    /** 合并 source2D 对象与基础配置 */
    mergeSource2DWithBase(target, base) {
      if (!target || !base) return
      Object.keys(base).forEach(key => {
        const targetValue = target[key]
        const baseValue = base[key]
        if (targetValue === undefined || targetValue === null) {
          target[key] = JSON.parse(JSON.stringify(baseValue))
        } else if (typeof baseValue === 'object' && baseValue !== null && !Array.isArray(baseValue)) {
          if (typeof targetValue === 'object' && targetValue !== null && !Array.isArray(targetValue)) {
            this.mergeSource2DWithBase(targetValue, baseValue)
          } else {
            target[key] = JSON.parse(JSON.stringify(baseValue))
          }
        } else if (Array.isArray(baseValue)) {
          if (!Array.isArray(targetValue)) {
            target[key] = JSON.parse(JSON.stringify(baseValue))
          } else if (targetValue.length === 0 && baseValue.length > 0) {
            target[key] = JSON.parse(JSON.stringify(baseValue))
          }
        }
      })
    },
    /** 统一规范化场景对象中的颜色字段，确保 <input type="color"> 永不报错 */
    normalizeSceneObjects(objects) {
      var self = this
      if (!Array.isArray(objects)) return
      objects.forEach(function(obj) {
        if (!obj) return
        // 规范化 color（GLTF 模型未修改材质时保留 undefined，不覆盖原始材质/贴图）
        if (isTextLabelObject(obj)) {
          obj.color = normTextLabelColor(obj.color)
        } else if (obj.type === 'gltf' && !obj.materialOverridden) {
          // 不设置默认值，保留 undefined 让 GLTF 原始材质生效
        } else {
          obj.color = self.normColorValue(obj.color)
        }
        // 规范化 textBgColor
        if (obj.type === 'text3d' || obj.type === 'dataText' || obj.type === 'textPlain3d' || obj.type === 'uiLabel') {
          obj.textBgColor = self.normColorValue(obj.textBgColor, '#00000000')
          if (obj.fontSize === undefined) obj.fontSize = 16
          if (obj.labelOpacity === undefined) obj.labelOpacity = 1
          if (obj.textBgOpacity === undefined) obj.textBgOpacity = obj.textBgColor === '#00000000' ? 0 : 0.2
          if (obj.labelFaceCamera === undefined) obj.labelFaceCamera = true
          if (obj.labelFixedSize === undefined) obj.labelFixedSize = obj.type !== 'dataText'
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
          // 不设置默认值
        } else if (!obj.emissive) {
          obj.emissive = '#000000'
        }
        // 确保动画属性存在
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
        // 流动管道属性
        if (obj.type === 'flowPipe') {
          if (obj.points === undefined || !Array.isArray(obj.points)) obj.points = []
          if (obj.radius === undefined) obj.radius = 0.1
          if (obj.flowSpeed === undefined) obj.flowSpeed = 1.0
          if (obj.highlightColor === undefined) obj.highlightColor = '#ff6a00'
          if (obj.flowDashLength === undefined) obj.flowDashLength = 3
        }
        // 数据绑定属性
        if (obj.wsKey === undefined) obj.wsKey = ''
        if (obj.bindProp === undefined) obj.bindProp = ''
        if (obj.bindTransform === undefined) obj.bindTransform = 'direct'
        if (obj.bindScale === undefined) obj.bindScale = 1
        if (obj.bindOffset === undefined) obj.bindOffset = 0
        ensurePositionBindings(obj)
        if (obj.type === '2dComponent') {
          obj.is2DComponent = true
          if (obj.label2D === undefined) obj.label2D = obj.name || obj.typeName || '2D Component'
          if (obj.source2D === undefined) obj.source2D = null
          if (obj.source2DMeta === undefined) obj.source2DMeta = null

          // 获取组件的基础配置（从组件定义中获取默认值）
          let baseInfo = null
          if (obj.source2D && obj.source2D.type) {
            const component = getISMComponent(obj.source2D.type)
            if (component && typeof component.data === 'function') {
              const compData = component.data()
              if (compData.base && compData.base.info) {
                baseInfo = JSON.parse(JSON.stringify(compData.base.info))
              }
            }
          }

          // 确保 source2D 对象存在
          if (!obj.source2D || typeof obj.source2D !== 'object') {
            obj.source2D = baseInfo || {
              type: obj.typeName || '2DComponent',
              action: [],
              active: [],
              animate: {
                selected: [],
                condition: {
                  deviceSN: '',
                  deviceName: '',
                  selectVideoType: 0,
                  isBandDevice: false,
                  bandType: 1,
                  dataID: '',
                  dataName: '',
                  operator: '',
                  OperatorValue: '',
                  OperatorMaxValue: ''
                },
                isExpression: false,
                animateList: [],
                animateElement: []
              },
              style: {
                position: { x: 0, y: 0, w: 200, h: 200 },
                backColor: 'transparent',
                foreColor: '#ffffff',
                fontSize: 14,
                fontFamily: 'Arial',
                zIndex: 1,
                transform: 0,
                diy: []
              }
            }
          }

          // 如果有基础配置，合并缺失的属性
          if (baseInfo) {
            self.mergeSource2DWithBase(obj.source2D, baseInfo)
          }
          if (obj.source2D) {
            obj.action = Array.isArray(obj.source2D.action) ? JSON.parse(JSON.stringify(obj.source2D.action)) : []
            obj.active = Array.isArray(obj.source2D.active) ? JSON.parse(JSON.stringify(obj.source2D.active)) : []
            obj.animate = obj.source2D.animate ? JSON.parse(JSON.stringify(obj.source2D.animate)) : obj.animate
          }

          // 确保 source2D.active 数组存在（用于数据绑定）
          if (!Array.isArray(obj.source2D.active)) {
            obj.source2D.active = []
          }

          // 如果没有 identifier，生成一个基于对象 id 的唯一标识
          if (!obj.source2D.identifier) {
            obj.source2D.identifier = 'scene_2d_' + obj.id
          }

          // 确保 animate 对象存在
          if (!obj.source2D.animate || typeof obj.source2D.animate !== 'object') {
            obj.source2D.animate = {
              selected: [],
              condition: {
                deviceSN: '',
                deviceName: '',
                selectVideoType: 0,
                isBandDevice: false,
                bandType: 1,
                dataID: '',
                dataName: '',
                operator: '',
                OperatorValue: '',
                OperatorMaxValue: ''
              },
              isExpression: false,
              animateList: [],
              animateElement: []
            }
          }

          // 确保 style 对象存在
          if (!obj.source2D.style || typeof obj.source2D.style !== 'object') {
            obj.source2D.style = {
              position: { x: 0, y: 0, w: 200, h: 200 },
              backColor: 'transparent',
              foreColor: '#ffffff',
              fontSize: 14,
              fontFamily: 'Arial',
              zIndex: 1,
              transform: 0,
              diy: []
            }
          }

          // 确保 position 对象存在
          if (!obj.source2D.style.position || typeof obj.source2D.style.position !== 'object') {
            obj.source2D.style.position = { x: 0, y: 0, w: 200, h: 200 }
          }
        }
        // 初始化 animate 对象（数据驱动动画）
        if (!obj.animate || typeof obj.animate !== 'object') {
          obj.animate = {
            isExpression: false,
            animateList: [],
            animateElement: '',
            condition: { dataID: '', operator: '', value: '', isBandDevice: false },
            selected: ''
          }
        } else {
          if (!obj.animate.condition) {
            obj.animate.condition = { dataID: '', operator: '', OperatorValue: '', OperatorMaxValue: '', isBandDevice: false, deviceSN: '' }
          } else {
            if (obj.animate.condition.dataID === undefined) obj.animate.condition.dataID = ''
            if (obj.animate.condition.operator === undefined) obj.animate.condition.operator = ''
            if (obj.animate.condition.OperatorValue === undefined) obj.animate.condition.OperatorValue = ''
            if (obj.animate.condition.OperatorMaxValue === undefined) obj.animate.condition.OperatorMaxValue = ''
            if (obj.animate.condition.isBandDevice === undefined) obj.animate.condition.isBandDevice = false
            if (obj.animate.condition.deviceSN === undefined) obj.animate.condition.deviceSN = ''
          }
          if (obj.animate.selected === undefined || !Array.isArray(obj.animate.selected)) obj.animate.selected = []
          if (obj.animate.isExpression === undefined) obj.animate.isExpression = false
          if (!obj.animate.animateList) obj.animate.animateList = []
        }
      })
    },

    applyMaterialProps(material, obj) {
      if (!material || !obj) return

      // 只有明确设置了颜色才修改，否则保留原始材质颜色
      if (obj.color) {
        const cStr = normColorValue(obj.color, '#4a90d9')
        const c = parseInt(cStr.replace('#', ''), 16)
        if (material.color && material.color.setHex) {
          material.color.setHex(c)
        }
      }

      // 只有明确设置了 opacity 才修改
      if (obj.opacity !== undefined) {
        material.opacity = obj.opacity
        material.transparent = obj.opacity < 1
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
        const emissive = normColorValue(obj.emissive, '#000000')
        material.emissive.setHex(parseInt(emissive.replace('#', ''), 16))
      }

      material.needsUpdate = true
    },

    applyObjectMaterialProps(mesh, obj) {
      if (!mesh || !obj || (mesh.userData && mesh.userData.isTextSprite)) return
      const applyOne = material => {
        if (!material || material.isShaderMaterial) return
        this.applyMaterialProps(material, obj)
      }
      if (mesh.material) {
        if (Array.isArray(mesh.material)) {
          mesh.material.forEach(applyOne)
        } else {
          applyOne(mesh.material)
        }
      }
      if (mesh.traverse) {
        mesh.traverse(child => {
          if (!child.isMesh || !child.material || child === mesh) return
          if (Array.isArray(child.material)) {
            child.material.forEach(applyOne)
          } else {
            applyOne(child.material)
          }
        })
      }
    },

    applyLightProps(mesh, obj) {
      if (!mesh || !obj || !obj.type || obj.type.indexOf('light') === -1) return
      const color = normColorValue(obj.color, '#ffffff')
      const colorInt = parseInt(color.replace('#', ''), 16)
      const intensity = obj.intensity !== undefined ? Number(obj.intensity) : 1
      const distance = obj.distance !== undefined ? Number(obj.distance) : undefined
      const angle = obj.angle !== undefined ? Number(obj.angle) : undefined
      const penumbra = obj.penumbra !== undefined ? Number(obj.penumbra) : undefined
      mesh.traverse(child => {
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

    _upgradeCanvasFallbackTexts() {
      if (!this.__3d || !this.__3d.scene) return
      const scene = this.__3d.scene
      const toUpgrade = []
      const toRebuild = []
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
        console.log('[ScenePreview] Upgrading', toUpgrade.length, 'Canvas fallback text(s) to true 3D')
      }
      if (toRebuild.length) {
        console.log('[ScenePreview] Rebuilding', toRebuild.length, 'True3D text(s) for font update')
      }
      const all = toUpgrade.concat(toRebuild)
      for (let i = 0; i < all.length; i++) {
        const mesh = all[i]
        const td = mesh.userData.textData || {}
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

    syncMeshes(objects) {
      if (!this.__3d) return
      const scene = this.__3d.scene
      const meshMap = this.__3d.meshMap

      // 清理旧 mesh
      const objIds = new Set(objects.map(o => o.id))
      Object.keys(meshMap).forEach(id => {
        if (!objIds.has(id)) {
          const mesh = meshMap[id]
          scene.remove(mesh)
          if (mesh.traverse) {
            mesh.traverse(child => {
              if (child.userData && child.userData.isMediaPlane) disposeMediaPlane(child)
              if (child.geometry) child.geometry.dispose()
              if (child.material) {
                if (Array.isArray(child.material)) child.material.forEach(m => m.dispose())
                else child.material.dispose()
              }
            })
          } else {
            if (mesh.userData && mesh.userData.isMediaPlane) disposeMediaPlane(mesh)
            if (mesh.geometry) mesh.geometry.dispose()
            if (mesh.material) mesh.material.dispose()
          }
          delete meshMap[id]
        }
      })

      // 创建或更新 mesh
      objects.forEach(obj => {
        let mesh = meshMap[obj.id]
        if (isUiOverlayObject(obj)) {
          if (mesh) {
            scene.remove(mesh)
            if (mesh.traverse) {
              mesh.traverse(child => {
                if (child.userData && child.userData.isMediaPlane) disposeMediaPlane(child)
                if (child.geometry) child.geometry.dispose()
                if (child.material) {
                  if (Array.isArray(child.material)) child.material.forEach(m => m.dispose())
                  else child.material.dispose()
                }
              })
            } else {
              if (mesh.userData && mesh.userData.isMediaPlane) disposeMediaPlane(mesh)
              if (mesh.geometry) mesh.geometry.dispose()
              if (mesh.material) mesh.material.dispose()
            }
            delete meshMap[obj.id]
          }
          return
        }
        if (obj.type === 'flowPipe' && mesh && !isFlowPipeSegmented(mesh)) {
          scene.remove(mesh)
          if (mesh.traverse) {
            mesh.traverse(child => {
              if (child.geometry) child.geometry.dispose()
              if (child.material) {
                if (Array.isArray(child.material)) child.material.forEach(m => m.dispose())
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
        // GLTF 模型跳过，由 restoreGLTFModels 单独处理
        if (!mesh && obj.type === 'gltf') return
        if (!mesh) {
          const colorStr = isTextLabelObject(obj) ? normTextLabelColor(obj.color) : normColorValue(obj.color, '#4a90d9')
          const colorInt = parseInt(colorStr.replace('#', ''), 16)
          if (obj.type === 'text3d') {
            const result = createText3DMesh(obj.textContent || '标签', colorInt, obj.fontSize || 16, obj.textBgColor || '#00000000', buildTextLabelOptions(obj))
            mesh = result ? result.mesh : null
          } else if (obj.type === 'dataText') {
            const result = createText3DMesh(getTextLabelContent(obj), colorInt, obj.fontSize || 16, obj.textBgColor || '#00000000', buildTextLabelOptions(obj))
            mesh = result ? result.mesh : null
          } else if (obj.type === 'textPlain3d' || obj.type === 'uiLabel') {
            const result = createText3DMesh(getTextLabelContent(obj), colorInt, obj.fontSize || 16, obj.textBgColor || '#00000000', buildTextLabelOptions(obj))
            mesh = result ? result.mesh : null
          } else if (obj.type === 'image3d' || obj.type === 'video3d') {
            mesh = createMediaPlane(obj, obj.type)
          } else if (obj.type === 'flowPipe') {
            console.log('创建流动管道:', obj)
            mesh = createFlowPipe(obj.points, obj.color, obj.radius, obj.flowSpeed, obj.highlightColor, obj.flowDashLength)
            console.log('管道创建结果:', mesh)
          } else if (obj.type === '2dComponent') {
            mesh = create2DComponentPlane(obj)
          } else {
            mesh = createThreeMeshFromLib(obj.type, colorInt)
          }
          if (mesh) {
            mesh.userData.id = obj.id
            scene.add(mesh)
            meshMap[obj.id] = mesh
          }
        }
        if (mesh) {
          mesh.position.set(obj.x || 0, obj.y || 0, obj.z || 0)
          mesh.rotation.set(obj.rx || 0, obj.ry || 0, obj.rz || 0)
          mesh.scale.set(obj.sx || 1, obj.sy || 1, obj.sz || 1)
          mesh.visible = obj.type === '2dComponent' ? false : (obj.visible !== false)

          if (obj.type === 'flowPipe' && mesh.userData && mesh.userData.isFlowPipe) {
            updateFlowPipe(mesh, obj)
            if (obj.materialOverridden === true) {
              this.applyObjectMaterialProps(mesh, obj)
            }
          } else if (mesh.material && !(mesh.userData && mesh.userData.isTextSprite)) {
            if (Array.isArray(mesh.material)) {
              mesh.material.forEach(material => this.applyMaterialProps(material, obj))
            } else {
              this.applyMaterialProps(mesh.material, obj)
            }
          }

          // text3d 类型特殊更新
          if (isTextLabelObject(obj) && mesh.userData && mesh.userData.isTextSprite) {
            const tcStr = normTextLabelColor(obj.color)
            const tc = parseInt(tcStr.replace('#', ''), 16)
            if (mesh.userData.is3DText) {
              updateText3DMesh(mesh, getTextLabelContent(obj), tc, obj.fontSize || 16, obj.textBgColor || '#00000000', buildTextLabelOptions(obj))
            } else {
              updateTextSprite(mesh, getTextLabelContent(obj), tc, obj.fontSize || 16, obj.textBgColor || '#00000000', buildTextLabelOptions(obj))
            }
            applyTextSpriteObjectScale(mesh, obj)
            updateTextPickArea(mesh, obj)
          }
          if ((obj.type === 'image3d' || obj.type === 'video3d') && mesh.userData && mesh.userData.isMediaPlane) {
            updateMediaPlane(mesh, obj, obj.type)
          }
          if (mesh.userData && mesh.userData.isLight) {
            this.applyLightProps(mesh, obj)
          }

          if (mesh.userData && mesh.userData.isMeshGroup && mesh.traverse) {
            // 只有当对象有自定义材质属性时才应用修改，否则保留原始材质
            if (obj.materialOverridden === true && (obj.color || obj.opacity !== undefined || obj.wireframe !== undefined || obj.metalness !== undefined || obj.roughness !== undefined || obj.emissive || obj.textureData !== undefined)) {
              this.applyObjectMaterialProps(mesh, obj)
            }
          }
        }
      })

      this.objectCount = objects.length
    },

    applySceneSettings(sceneSettings) {
      if (!sceneSettings || !this.__3d || !this.__3d.scene) return

      var ss = sceneSettings
      var envStyle = this.getEnvironmentStyle(ss.environmentPreset)
      var bgColor = parseInt((ss.backgroundColor || '#ffffff').replace('#', '0x'), 16)

      if (ss.backgroundMode === 'gradient') {
        this.__3d.scene.background = this.createGradientTexture(ss.backgroundColor || envStyle.top, ss.backgroundColor2 || envStyle.bottom)
      } else if ((ss.backgroundMode === 'image' || !ss.backgroundMode) && ss.backgroundImage) {
        var loader = new THREE.TextureLoader()
        var self = this
        loader.load(
          ss.backgroundImage,
          function(texture) {
            texture.wrapS = THREE.RepeatWrapping
            texture.wrapT = THREE.RepeatWrapping
            self.__3d.scene.background = texture
          },
          undefined,
          function() {
            self.__3d.scene.background = new THREE.Color(bgColor)
          }
        )
      } else {
        this.__3d.scene.background = this.createGradientTexture(envStyle.top, ss.backgroundColor || envStyle.bottom)
      }
      this.applyEnvironmentFloor(ss, envStyle)

      if (this.__3d.scene.fog) {
        this.__3d.scene.fog.color = new THREE.Color(bgColor)
      }

      if (ss.showGrid !== undefined && this.__3d.gridHelper) {
        this.showGrid = ss.showGrid
        this.__3d.gridHelper.visible = ss.showGrid
      }
      this.applySceneVisualSettings(ss)
    },

    createGradientTexture(colorTop, colorBottom) {
      const canvas = document.createElement('canvas')
      canvas.width = 16
      canvas.height = 256
      const ctx = canvas.getContext('2d')
      const gradient = ctx.createLinearGradient(0, 0, 0, canvas.height)
      gradient.addColorStop(0, colorTop || '#ffffff')
      gradient.addColorStop(1, colorBottom || '#dbefff')
      ctx.fillStyle = gradient
      ctx.fillRect(0, 0, canvas.width, canvas.height)
      const texture = new THREE.CanvasTexture(canvas)
      texture.needsUpdate = true
      return texture
    },

    getEnvironmentStyle(preset) {
      const styles = {
        sky: { top: '#1fa4bd', bottom: '#59d1e4', floor: '#8ea0a3', fog: '#cceff5' },
        sunset: { top: '#f39a52', bottom: '#ffe1b8', floor: '#9d8f82', fog: '#ffe0bf' },
        ocean: { top: '#117da2', bottom: '#65d7e8', floor: '#6f9ca6', fog: '#c7f4fb' },
        forest: { top: '#2d8f71', bottom: '#bfe9c8', floor: '#6e8a76', fog: '#d6f0dc' },
        twilight: { top: '#4855a9', bottom: '#c9bcff', floor: '#7f7c94', fog: '#ded8ff' },
        night: { top: '#0d1630', bottom: '#283a66', floor: '#394657', fog: '#283852' }
      }
      return styles[preset || 'sky'] || styles.sky
    },

    applyEnvironmentFloor(ss, envStyle) {
      if (!this.__3d || !this.__3d.scene) return
      const scene = this.__3d.scene
      if (this.__3d.environmentFloorMesh) {
        scene.remove(this.__3d.environmentFloorMesh)
        if (this.__3d.environmentFloorMesh.geometry) this.__3d.environmentFloorMesh.geometry.dispose()
        if (this.__3d.environmentFloorMesh.material) this.__3d.environmentFloorMesh.material.dispose()
        this.__3d.environmentFloorMesh = null
      }
      if (ss.backgroundMode === 'image') return
      const geometry = new THREE.PlaneGeometry(260, 260)
      const material = new THREE.MeshStandardMaterial({
        color: new THREE.Color(envStyle.floor || '#8ea0a3'),
        roughness: 0.86,
        metalness: 0.04,
        transparent: true,
        opacity: 0.82,
        side: THREE.DoubleSide
      })
      const mesh = new THREE.Mesh(geometry, material)
      mesh.rotation.x = -Math.PI / 2
      mesh.position.y = -0.035
      mesh.receiveShadow = true
      mesh.userData.isEnvironmentFloor = true
      scene.add(mesh)
      this.__3d.environmentFloorMesh = mesh
    },

    applySceneVisualSettings(settings) {
      if (!this.__3d || !this.__3d.scene) return
      const ss = settings || {}
      this.applyLightingPreset(ss.lightingPreset || 'day', !!ss.enhanceDepth)
      this.applyFloorReflection(ss.floorReflection || 'none')
      this.applyDepthEnhancement(ss)
      if (this.__3d.renderer) {
        this.__3d.renderer.setPixelRatio(ss.modelOptimize ? Math.min(window.devicePixelRatio || 1, 1.25) : (window.devicePixelRatio || 1))
        this.__3d.renderer.shadowMap.enabled = !ss.modelOptimize
        this.__3d.renderer.toneMappingExposure = ss.enhanceDepth ? 1.12 : 1.0
      }
    },

    applyDepthEnhancement(ss) {
      if (!this.__3d || !this.__3d.scene) return
      const scene = this.__3d.scene
      if (ss.enhanceDepth) {
        const fogColor = parseInt((ss.backgroundColor || '#ffffff').replace('#', '0x'), 16)
        scene.fog = new THREE.Fog(fogColor, 36, 95)
      } else if (scene.fog && !scene.userData.fogEnabledByExtras) {
        scene.fog = null
      }
      scene.traverse((child) => {
        if (!child.isMesh || !child.material) return
        child.castShadow = !ss.modelOptimize
        child.receiveShadow = !ss.modelOptimize
        const materials = Array.isArray(child.material) ? child.material : [child.material]
        materials.forEach((material) => {
          if (!material) return
          if (material.roughness !== undefined && ss.enhanceDepth) {
            material.roughness = Math.max(0.22, material.roughness * 0.92)
          }
          material.needsUpdate = true
        })
      })
    },

    applyLightingPreset(preset, enhanceDepth) {
      if (!this.__3d || !this.__3d.scene) return
      const configs = {
        day: { ambient: 1.0, directional: 0.95, color: 0xffffff, fill: 0x4466aa, position: [8, 12, 6] },
        evening: { ambient: 0.58, directional: 1.05, color: 0xffb36f, fill: 0x6e7fb8, position: [-8, 5, 8] },
        night: { ambient: 0.28, directional: 0.55, color: 0x9ec8ff, fill: 0x1b2a55, position: [4, 9, -7] },
        studio: { ambient: 1.18, directional: 0.8, color: 0xffffff, fill: 0xdde8ff, position: [0, 10, 8] },
        industrial: { ambient: 0.78, directional: 1.18, color: 0xf2f7ff, fill: 0x9fb7c8, position: [10, 10, -4] }
      }
      const cfg = configs[preset] || configs.day
      const lights = this.__3d.sceneLights || []
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

    applyFloorReflection(level) {
      if (!this.__3d || !this.__3d.scene) return
      const scene = this.__3d.scene
      if (this.__3d.floorReflectionMesh) {
        scene.remove(this.__3d.floorReflectionMesh)
        if (this.__3d.floorReflectionMesh.geometry) this.__3d.floorReflectionMesh.geometry.dispose()
        if (this.__3d.floorReflectionMesh.material) this.__3d.floorReflectionMesh.material.dispose()
        this.__3d.floorReflectionMesh = null
      }
      if (!level || level === 'none') return
      const opacityMap = { matte: 0.12, low: 0.2, medium: 0.32, strong: 0.48 }
      const roughnessMap = { matte: 0.95, low: 0.72, medium: 0.48, strong: 0.24 }
      const geometry = new THREE.PlaneGeometry(200, 200)
      const material = new THREE.MeshStandardMaterial({
        color: 0xeef7ff,
        transparent: true,
        opacity: opacityMap[level] || 0.12,
        metalness: level === 'strong' ? 0.75 : 0.35,
        roughness: roughnessMap[level] || 0.72,
        side: THREE.DoubleSide
      })
      const mesh = new THREE.Mesh(geometry, material)
      mesh.rotation.x = -Math.PI / 2
      mesh.position.y = -0.012
      mesh.receiveShadow = true
      scene.add(mesh)
      this.__3d.floorReflectionMesh = mesh
    },

    applySceneExtras(sceneExtras) {
      if (!sceneExtras || !this.__3d || !this.__3d.scene) return
      const scene = this.__3d.scene
      const envColors = {
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
      const ss = this._savedSceneSettings || {}
      if (!ss.backgroundMode && !ss.backgroundImage) {
        const color = envColors[ss.environmentPreset || sceneExtras.environment] || envColors.studio
        scene.background = new THREE.Color(color)
      }
      if (sceneExtras.fogEnabled) {
        scene.fog = new THREE.Fog(sceneExtras.fogColor || '#d8e4f0', 30, 120)
        scene.userData.fogEnabledByExtras = true
      } else {
        scene.userData.fogEnabledByExtras = false
        scene.fog = null
      }
      if (sceneExtras.defaultCameraId && Array.isArray(sceneExtras.cameraViews)) {
        const cameraView = sceneExtras.cameraViews.find(item => item.id === sceneExtras.defaultCameraId)
        if (cameraView) {
          this.restoreCamera({
            cameraPosition: cameraView.position,
            cameraTarget: cameraView.target,
            cameraFov: cameraView.fov
          })
        }
      }
      this.applySceneVisualSettings(ss)
    },

    applyTimelineTracks(elapsed) {
      if (!this.sceneExtras || !this.sceneExtras.timelineEnabled || !Array.isArray(this.sceneExtras.timelineTracks)) return
      const tracks = this.sceneExtras.timelineTracks.filter(item => item && item.enabled !== false && item.objectId && item.property)
      if (!tracks.length || !this.__3d || !this.__3d.meshMap) return

      const groups = {}
      let duration = 0
      tracks.forEach(item => {
        const key = item.objectId + '::' + item.property
        if (!groups[key]) groups[key] = []
        groups[key].push(item)
        duration = Math.max(duration, Number(item.time) || 0)
      })
      const configuredDuration = Number(this.sceneExtras.timelineDuration)
      if (!isNaN(configuredDuration) && configuredDuration > 0) {
        duration = Math.max(duration, configuredDuration)
      }
      if (duration <= 0) duration = 1
      const speed = Number(this.sceneExtras.timelineSpeed)
      const rawTime = elapsed * (isNaN(speed) || speed <= 0 ? 1 : speed)
      const time = this.sceneExtras.timelineLoop === false ? Math.min(rawTime, duration) : (rawTime % duration)

      Object.keys(groups).forEach(key => {
        const list = groups[key].sort((a, b) => (Number(a.time) || 0) - (Number(b.time) || 0))
        const first = list[0]
        const mesh = this.__3d.meshMap[first.objectId]
        if (!mesh) return
        const sourceObject = this.sceneObjects.find(item => item && item.id === first.objectId) || {}
        const normalizedList = list.slice()
        if ((Number(normalizedList[0].time) || 0) > 0) {
          normalizedList.unshift({
            id: first.objectId + '_timeline_base',
            objectId: first.objectId,
            property: first.property,
            time: 0,
            value: this.getTimelineBaseValue(sourceObject, first.property)
          })
        }
        let prev = normalizedList[0]
        let next = normalizedList[normalizedList.length - 1]
        for (let i = 0; i < normalizedList.length; i++) {
          if ((Number(normalizedList[i].time) || 0) <= time) prev = normalizedList[i]
          if ((Number(normalizedList[i].time) || 0) >= time) {
            next = normalizedList[i]
            break
          }
        }
        if (next === prev) {
          this.applyTimelineValue(mesh, first.property, prev.value)
          this.applyTimelineObjectValue(sourceObject, first.property, prev.value)
          return
        }
        const start = Number(prev.time) || 0
        const end = Number(next.time) || duration
        const ratio = end === start ? 0 : (time - start) / (end - start)
        const nextValue = this.interpolateVectorValue(prev.value, next.value, ratio)
        this.applyTimelineValue(mesh, first.property, nextValue)
        this.applyTimelineObjectValue(sourceObject, first.property, nextValue)
      })
    },
    getTimelineBaseValue(obj, property) {
      obj = obj || {}
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
    interpolateVectorValue(fromValue, toValue, ratio) {
      const from = fromValue || {}
      const to = toValue || {}
      const nextValue = {
        x: (Number(from.x) || 0) + ((Number(to.x) || 0) - (Number(from.x) || 0)) * ratio,
        y: (Number(from.y) || 0) + ((Number(to.y) || 0) - (Number(from.y) || 0)) * ratio,
        z: (Number(from.z) || 0) + ((Number(to.z) || 0) - (Number(from.z) || 0)) * ratio
      }
      if (from.rx !== undefined || to.rx !== undefined) {
        nextValue.rx = (Number(from.rx) || 0) + ((Number(to.rx) || 0) - (Number(from.rx) || 0)) * ratio
        nextValue.ry = (Number(from.ry) || 0) + ((Number(to.ry) || 0) - (Number(from.ry) || 0)) * ratio
        nextValue.rz = (Number(from.rz) || 0) + ((Number(to.rz) || 0) - (Number(from.rz) || 0)) * ratio
      }
      if (from.sx !== undefined || to.sx !== undefined) {
        nextValue.sx = (Number(from.sx) || 1) + ((Number(to.sx) || 1) - (Number(from.sx) || 1)) * ratio
        nextValue.sy = (Number(from.sy) || 1) + ((Number(to.sy) || 1) - (Number(from.sy) || 1)) * ratio
        nextValue.sz = (Number(from.sz) || 1) + ((Number(to.sz) || 1) - (Number(from.sz) || 1)) * ratio
      }
      return nextValue
    },
    applyTimelineValue(mesh, property, value) {
      if (!mesh || !value) return
      if (property === 'position') {
        mesh.position.set(Number(value.x) || 0, Number(value.y) || 0, Number(value.z) || 0)
      } else if (property === 'rotation') {
        mesh.rotation.set(Number(value.x) || 0, Number(value.y) || 0, Number(value.z) || 0)
      } else if (property === 'scale') {
        mesh.scale.set(Number(value.x) || 1, Number(value.y) || 1, Number(value.z) || 1)
      } else if (property === 'transform') {
        mesh.position.set(Number(value.x) || 0, Number(value.y) || 0, Number(value.z) || 0)
        mesh.rotation.set(Number(value.rx) || 0, Number(value.ry) || 0, Number(value.rz) || 0)
        mesh.scale.set(Number(value.sx) || 1, Number(value.sy) || 1, Number(value.sz) || 1)
      }
    },
    applyTimelineObjectValue(obj, property, value) {
      if (!obj || !value) return
      if (property === 'position') {
        obj.x = Number(value.x) || 0
        obj.y = Number(value.y) || 0
        obj.z = Number(value.z) || 0
      } else if (property === 'rotation') {
        obj.rx = Number(value.x) || 0
        obj.ry = Number(value.y) || 0
        obj.rz = Number(value.z) || 0
      } else if (property === 'scale') {
        obj.sx = Number(value.x) || 1
        obj.sy = Number(value.y) || 1
        obj.sz = Number(value.z) || 1
      } else if (property === 'transform') {
        obj.x = Number(value.x) || 0
        obj.y = Number(value.y) || 0
        obj.z = Number(value.z) || 0
        obj.rx = Number(value.rx) || 0
        obj.ry = Number(value.ry) || 0
        obj.rz = Number(value.rz) || 0
        obj.sx = Number(value.sx) || 1
        obj.sy = Number(value.sy) || 1
        obj.sz = Number(value.sz) || 1
      }
    },

    getIntersectedObject(event) {
      if (!this.__3d || !this.$refs.previewCanvas) return null
      const rect = this.$refs.previewCanvas.getBoundingClientRect()
      this.__3d.mouse.x = ((event.clientX - rect.left) / rect.width) * 2 - 1
      this.__3d.mouse.y = -((event.clientY - rect.top) / rect.height) * 2 + 1
      this.__3d.raycaster.setFromCamera(this.__3d.mouse, this.__3d.camera)

      const allMeshes = []
      Object.values(this.__3d.meshMap).forEach(m => {
        if (m && m.traverse) {
          m.traverse(c => { if (c.isMesh) allMeshes.push(c) })
        } else if (m) {
          allMeshes.push(m)
        }
      })
      const hits = this.__3d.raycaster.intersectObjects(allMeshes, false)
      if (!hits.length) return null
      let target = hits[0].object
      let id = target.userData && target.userData.id
      while (!id && target.parent) {
        target = target.parent
        id = target.userData && target.userData.id
      }
      return id ? this.sceneObjects.find(obj => obj.id === id) || null : null
    },
    hasPreviewActions(obj) {
      if (!obj) return false
      const actionSource = obj.type === '2dComponent' && obj.source2D ? obj.source2D.action : obj.action
      const activeSource = obj.type === '2dComponent' && obj.source2D ? obj.source2D.active : obj.active
      const hasAction = Array.isArray(actionSource) && actionSource.some(item => item && this.match2DActionTrigger(item, item.type || item.eventType || item.EventType || item.trigger || 'click'))
      const hasActive = Array.isArray(activeSource) && activeSource.some(item => item && item.isSwitch === true)
      return hasAction || hasActive
    },
    executeInteraction(obj, eventName, event) {
      if (!obj) return
      const actionSource = obj.type === '2dComponent' && obj.source2D ? obj.source2D.action : obj.action
      const activeSource = obj.type === '2dComponent' && obj.source2D ? obj.source2D.active : obj.active
      const actionItems = Array.isArray(actionSource) ? actionSource.filter(item => this.match2DActionTrigger(item, eventName)) : []
      const workbenchEvents = obj.interactions && Array.isArray(obj.interactions.events)
        ? obj.interactions.events.filter(item => item && item.enabled !== false && this.matchWorkbenchTrigger(item, eventName))
        : []
      actionItems.forEach(actionItem => {
        this.executeActionItem(obj, actionItem)
      })
      if (!actionItems.length && eventName === 'click' && Array.isArray(activeSource)) {
        activeSource.filter(item => item && item.isSwitch === true).forEach(activeItem => {
          this.executeActiveSwitch(activeItem)
        })
      }
      workbenchEvents.forEach(eventItem => {
        this.executeWorkbenchEvent(obj, eventItem)
      })
    },
    match2DActionTrigger(actionItem, eventName) {
      if (!actionItem) return false
      const type = actionItem.type || actionItem.eventType || actionItem.EventType || actionItem.trigger
      if (type === eventName) return true
      if (type === 'dbclick' && eventName === 'dblclick') return true
      return false
    },
    matchWorkbenchTrigger(eventItem, eventName) {
      if (!eventItem || !eventItem.trigger) return false
      if (eventItem.trigger === eventName) return true
      if (eventItem.trigger === 'hover' && eventName === 'mouseenter') return true
      return false
    },
    executeActiveSwitch(activeItem) {
      const condition = activeItem && activeItem.condition ? activeItem.condition : null
      if (!condition) return
      const run = () => {
        const nextCondition = Object.assign({}, condition)
        if (!nextCondition.IsManual && Array.isArray(nextCondition.AutoSet) && nextCondition.AutoSet.length > 0) {
          const currentIndex = nextCondition.AutoSet.findIndex(item => item && item.value === activeItem.result)
          const nextIndex = currentIndex === 0 && nextCondition.AutoSet.length > 1 ? 1 : 0
          nextCondition.AutoSetValue = nextCondition.AutoSet[nextIndex] ? nextCondition.AutoSet[nextIndex].value : nextCondition.AutoSetValue
        }
        this.executeSetValueItem(nextCondition)
      }
      if (condition.ConfirmationDialog === 1 || condition.actionConfirm) {
        this.$confirm({
          title: '二次确认',
          content: '确认执行该组件行为？',
          onOk: () => run(),
        })
      } else {
        run()
      }
    },
    executeWorkbenchEvent(obj, eventItem) {
      if (!obj || !eventItem || eventItem.enabled === false) return
      const value = eventItem.value || ''
      switch (eventItem.action) {
        case 'focus':
          this.focusObject(obj)
          break
        case 'camera':
          this.applyCameraViewById(value)
          break
        case 'toggleVisible':
          this.toggleTargetVisible(value || obj.id)
          break
        case 'openPanel':
          if (value) {
            this.popupInteraction = {
              visible: true,
              url: value,
              title: obj.name || '页面弹窗',
              width: 960,
              height: 640,
              autoClose: false
            }
          }
          break
        case 'callApi':
          if (value) {
            request(value, String(eventItem.method || 'get').toLowerCase()).catch(() => {
              this.$message.error(this.$t ? this.$t('readData.SetFailed') : '执行失败')
            })
          }
          break
      }
    },
    executeDataWorkbenchEvents(obj, currentData) {
      if (!obj || !currentData || !obj.interactions || !Array.isArray(obj.interactions.events)) return
      const events = obj.interactions.events.filter(item => item && item.enabled !== false && item.trigger === 'data')
      if (!events.length) return
      events.forEach(eventItem => {
        const value = (eventItem.value || '').trim()
        const matched = !value ||
          value === currentData.ModelDataUuid ||
          value === currentData.Uuid ||
          value === currentData.Name ||
          value === currentData.Identifier
        if (matched) this.executeWorkbenchEvent(obj, eventItem)
      })
    },
    focusObject(obj) {
      if (!this.__3d || !obj) return
      const mesh = this.__3d.meshMap[obj.id]
      if (!mesh) return
      const box = new THREE.Box3().setFromObject(mesh)
      const center = box.getCenter(new THREE.Vector3())
      const size = box.getSize(new THREE.Vector3())
      const maxDim = Math.max(size.x, size.y, size.z, 1)
      this.__3d.orbit.target.copy(center)
      this.__3d.camera.position.set(center.x + maxDim * 2, center.y + maxDim * 1.4, center.z + maxDim * 2)
      this.__3d.camera.lookAt(center)
      this.__3d.orbit.update()
    },
    applyCameraViewById(cameraId) {
      if (!cameraId || !this.sceneExtras || !Array.isArray(this.sceneExtras.cameraViews)) return
      const cameraView = this.sceneExtras.cameraViews.find(item => item.id === cameraId)
      if (!cameraView) return
      this.restoreCamera({
        cameraPosition: cameraView.position,
        cameraTarget: cameraView.target,
        cameraFov: cameraView.fov
      })
    },
    toggleTargetVisible(target) {
      const tokens = Array.isArray(target) ? target : String(target || '').split(',').map(item => item.trim()).filter(Boolean)
      if (!tokens.length) return
      this.sceneObjects.forEach(obj => {
        if (!obj) return
        const source2D = obj.source2D || {}
        const matches = [obj.id, obj.name, source2D.identifier, source2D.groupID, source2D.GroupName, source2D.name].filter(Boolean)
        if (matches.some(token => tokens.includes(token))) {
          this.applyObjectChanges(obj, { visible: obj.visible === false })
        }
      })
    },
    executeActionItem(obj, actionItem) {
      if (!actionItem) return

      try {
        if ((typeof this.user !== 'undefined') && this.user && this.user.Role && this.user.Role !== 'Admin') {
          if (actionItem.actionAuth && actionItem.actionAuth.length > 0) {
            const auth = actionItem.actionAuth.find(item => item === this.user.Role)
            if (auth !== this.user.Role) {
              this.$message.error(this.$t ? this.$t('displayModel.NoAuthentication') : 'No authentication')
              return
            }
          }
        }
      } catch (e) {
        this.$message.error(this.$t ? this.$t('displayModel.NoAuthentication') : 'No authentication')
        return
      }

      if (actionItem.ActionPassword) {
        this.ActionPasswordSet = actionItem.ActionPassword
        this.ActionPasswordValue = ''
        this._pendingActionRunner = () => this.executeActionItem(obj, Object.assign({}, actionItem, { ActionPassword: '' }))
        this.actionPasswordDialog = true
        return
      }


      const runAction = () => {
        const actionVoice = actionItem.actionVoice || actionItem.ActionVoice || ''
        if (actionVoice) {
          speechSynthesis.cancel()
          const speech = new SpeechSynthesisUtterance(actionVoice)
          speech.lang = 'zh-CN'
          speechSynthesis.speak(speech)
        }

        if (actionItem.action === 'popupText') {
          this.$info({
            title: obj.name || '组件行为',
            content: actionItem.popupText || '',
          })
        } else if (actionItem.action === 'visible') {
          this.applyVisibilityTargets(actionItem.showItems || [], true)
          this.applyVisibilityTargets(actionItem.hideItems || [], false)
        } else if ((actionItem.action === 'link' || actionItem.action === 'route') && actionItem.link) {
          if (actionItem.link.linkType === 'Inside') {
            const route = this.resolveInsideLinkRoute(actionItem.link)
            if (actionItem.link.isPopUp) {
              this.popupInteraction = {
                visible: true,
                url: route.href,
                title: actionItem.link.title || obj.name || '页面弹窗',
                width: actionItem.link.width || 960,
                height: actionItem.link.height || 640,
                autoClose: !!actionItem.link.autoClose,
              }
            } else {
              const targetLocation = route.location || {
                path: route.route && route.route.path ? route.route.path : this.$route.path,
                query: route.route && route.route.query ? route.route.query : {}
              }
              this.$router.push(targetLocation).catch(() => {})
            }
          } else {
            const target = actionItem.link.OpenExternalType === 'self' ? '_self' : '_blank'
            if (target === '_self') {
              window.location.href = actionItem.link.External || ''
            } else {
              const features = 'width=' + (actionItem.link.width || 960) + ',height=' + (actionItem.link.height || 640)
              window.open(actionItem.link.External || '', actionItem.link.title || '_blank', features)
            }
          }
        } else if (actionItem.action === 'DeviceView' && actionItem.DeviceView) {
          const onSelectData = {
            key: actionItem.DeviceView.key,
            showUUID: actionItem.DeviceView.showUUID,
            showPageUUID: actionItem.DeviceView.showPageUUID,
            type: actionItem.DeviceView.type,
            isPopUp: !!actionItem.DeviceView.isPopUp,
            selectKey: actionItem.DeviceView.key,
            isContainer: !!actionItem.DeviceView.isContainer,
          }
          if (onSelectData.key || onSelectData.showUUID || onSelectData.showPageUUID) {
            if (onSelectData.isContainer) {
              this.$EventBus.$emit('onContainerSelectDevice', onSelectData)
            } else {
              this.$EventBus.$emit('onSelectDevice', onSelectData)
            }
          } else if (actionItem.DeviceView.routePath) {
            const route = this.$router.resolve({
              path: actionItem.DeviceView.routePath,
              query: actionItem.DeviceView.key ? { deviceKey: actionItem.DeviceView.key } : {}
            })
            window.open(route.href, obj.name || '_blank')
          }
        } else if (actionItem.action === 'RestApi' && actionItem.RestApi) {
          ComponentRestApi(actionItem.RestApi.Type || 'Post', actionItem.RestApi.Url || '', actionItem.RestApi.Params || '{}').catch(() => {
            this.$message.error(this.$t ? this.$t('readData.SetFailed') : '执行失败')
          })
        } else if (actionItem.action === 'SysScript') {
          this.doSysScript(actionItem.ScriptList || (actionItem.SysScript && actionItem.SysScript.ScriptList) || [])
        } else if (actionItem.action === 'SetValue') {
          this.doLoopSetValue(0, actionItem.setValue || [], actionItem.SetDelay)
        } else if (actionItem.action === 'Animation') {
          const animateIdentifier = obj && obj.type === '2dComponent' && obj.source2D ? obj.source2D.identifier : obj.identifier
          const animationStatus = actionItem.animationStatus || (actionItem.Animation && actionItem.Animation.animationStatus) || ''
          if (animateIdentifier) {
            this.$EventBus.$emit(animateIdentifier + 'animateEvent', animationStatus === 'start' ? 1 : 0)
          }
        }
      }

      if (actionItem.actionConfirm || actionItem.SecondConfirm) {
        this.$confirm({
          title: '二次确认',
          content: '确认执行该组件行为？',
          onOk: () => runAction(),
        })
      } else {
        runAction()
      }
    },
    applyVisibilityTargets(targetsText, visible) {
      let ids = []
      if (Array.isArray(targetsText)) {
        ids = targetsText
      } else if (targetsText) {
        ids = targetsText.split(',').map(item => item.trim()).filter(Boolean)
      }
      if (!ids.length) return
      this.sceneObjects.forEach(obj => {
        const source2D = obj && obj.source2D ? obj.source2D : null
        const matches = [
          obj.id,
          obj.name,
          source2D && source2D.identifier,
          source2D && source2D.groupID,
          source2D && source2D.GroupName,
          source2D && source2D.name
        ].filter(Boolean)
        if (matches.some(token => ids.includes(token))) {
          this.applyObjectChanges(obj, { visible: visible })
        }
      })
    },
    resolveInsideLinkRoute(linkInfo) {
      const inside = linkInfo && linkInfo.Inside ? linkInfo.Inside : {}
      if (linkInfo && linkInfo.routePath) {
        return this.$router.resolve({ path: linkInfo.routePath })
      }
      if (inside.displayUUID) {
        const is3DDisplay = inside.displayType === 2 || inside.displayType === '2'
        return this.$router.resolve({
          path: (is3DDisplay ? '/DisPlay3DRunApp/' : '/AppRun/') + inside.displayUUID,
          query: inside.pageUUID ? { pageId: inside.pageUUID } : {}
        })
      }
      return this.$router.resolve({ path: this.$route.path })
    },
    closePopupInteraction() {
      this.popupInteraction.visible = false
      this.popupInteraction.url = ''
    },
    doSysScript(scriptList) {
      ExecSysScript({ Script: scriptList }).then((res) => {
        if (res && res.data && res.data.code === 0) {
          this.$message.success(this.$t ? this.$t(res.data.msg) : '执行成功')
        } else {
          this.$message.error(this.$t ? this.$t((res && res.data && res.data.msg) || 'readData.SetFailed') : '执行失败')
        }
      }).catch(() => {
        this.$message.error(this.$t ? this.$t('readData.SetFailed') : '执行失败')
      })
    },
    doLoopSetValue(counter, setValueList, delay) {
      const setValueItem = setValueList[counter]
      if (!setValueItem) return
      const wait = parseInt(delay, 10)
      this.executeSetValueItem(setValueItem).finally(() => {
        if (counter + 1 < setValueList.length) {
          setTimeout(() => this.doLoopSetValue(counter + 1, setValueList, isNaN(wait) ? 1000 : wait), isNaN(wait) ? 1000 : wait)
        }
      })
    },
    executeSetValueItem(setValueItem) {
      const deviceUuid = setValueItem.isBandDevice ? setValueItem.deviceSN : (this.$route.query.deviceKey || '')
      if (!deviceUuid || !setValueItem.dataID) {
        this.$message.error(this.$t ? this.$t('readData.SetFailed') : '设置失败')
        return Promise.resolve()
      }

      if (setValueItem.IsManual) {
        this._pendingSetValueItem = setValueItem
        this._pendingSetValueDeviceUuid = deviceUuid
        this._pendingSetValueResolve = null
        this.SetPassword = setValueItem.SetPassword || ''
        this.SetPasswordFormValue = ''
        this.SetValueFormValue = ''
        this.settingDialog = true
        return new Promise((resolve) => {
          this._pendingSetValueResolve = resolve
        })
      }

      if (setValueItem.SetPassword) {
        this._pendingSetValueItem = setValueItem
        this._pendingSetValueDeviceUuid = deviceUuid
        this._pendingSetValueResolve = null
        this.SetPassword = setValueItem.SetPassword
        this.SetAutoPasswordFormValue = ''
        this.setPasswordDialog = true
        return new Promise((resolve) => {
          this._pendingSetValueResolve = resolve
        })
      }


      return this.submitSetData(deviceUuid, setValueItem.dataID, setValueItem.AutoSetValue)
    },
    submitSetData(deviceUuid, dataUuid, value) {
      return setData({
        deviceUuid: deviceUuid,
        dataUuid: dataUuid,
        value: value === undefined || value === null ? '' : String(value)
      }).then((res) => {
        if (res && res.data && res.data.code === 0) {
          this.$message.success(this.$t ? this.$t('readData.SetSuccess') : '设置成功')
        } else {
          this.$message.error(this.$t ? this.$t('readData.SetFailed') : '设置失败')
        }
      }).catch(() => {
        this.$message.error(this.$t ? this.$t('readData.SetFailed') : '设置失败')
      })
    },
    closeSetDialogs() {
      this.settingDialog = false
      this.setPasswordDialog = false
      if (this._pendingSetValueResolve) {
        const resolve = this._pendingSetValueResolve
        this._pendingSetValueResolve = null
        resolve()
      }
    },
    closeActionPasswordDialog() {
      this.actionPasswordDialog = false
      this._pendingActionRunner = null
    },
    PasswordSetAction() {
      if (this.ActionPasswordSet !== this.ActionPasswordValue) {
        this.$message.error(this.$t ? this.$t('readData.SetPasswordError') : '密码错误')
        return
      }
      const runner = this._pendingActionRunner
      this.actionPasswordDialog = false
      this._pendingActionRunner = null
      this.ActionPasswordValue = ''
      this.ActionPasswordSet = ''
      if (runner) runner()
    },
    ManualSetData() {
      if (!this._pendingSetValueItem) return
      if (this.SetPassword && this.SetPasswordFormValue !== this.SetPassword) {
        this.$message.error(this.$t ? this.$t('readData.SetPasswordError') : '密码错误')
        return
      }
      this.settingLoading = true
      this.submitSetData(this._pendingSetValueDeviceUuid, this._pendingSetValueItem.dataID, this.SetValueFormValue).finally(() => {
        this.settingLoading = false
        const resolve = this._pendingSetValueResolve
        this._pendingSetValueResolve = null
        this.settingDialog = false
        this.SetPasswordFormValue = ''
        this.SetValueFormValue = ''
        if (resolve) resolve()
      })
    },
    PasswordSetData() {
      if (!this._pendingSetValueItem) return
      if (this.SetAutoPasswordFormValue !== this.SetPassword) {
        this.$message.error(this.$t ? this.$t('readData.SetPasswordError') : '密码错误')
        return
      }
      this.settingLoading = true
      this.submitSetData(this._pendingSetValueDeviceUuid, this._pendingSetValueItem.dataID, this._pendingSetValueItem.AutoSetValue).finally(() => {
        this.settingLoading = false
        const resolve = this._pendingSetValueResolve
        this._pendingSetValueResolve = null
        this.setPasswordDialog = false
        this.SetAutoPasswordFormValue = ''
        if (resolve) resolve()
      })
    },
    handlePopupMaskClick() {
      if (this.popupInteraction.autoClose) {
        this.closePopupInteraction()
      }
    },
    onCanvasClick(event) {
      const obj = this.getIntersectedObject(event)
      this.executeInteraction(obj, 'click')
    },
    onCanvasDblClick(event) {
      const obj = this.getIntersectedObject(event)
      this.executeInteraction(obj, 'dblclick')

      // 双击聚焦：移动相机到双击的组件
      if (obj && this.__3d && this.__3d.meshMap[obj.id]) {
        const mesh = this.__3d.meshMap[obj.id]
        const box = new THREE.Box3().setFromObject(mesh)
        const center = box.getCenter(new THREE.Vector3())
        const size = box.getSize(new THREE.Vector3())
        const maxDim = Math.max(size.x, size.y, size.z, 0.5)
        const dist = Math.max(maxDim * 4, 3)

        const startTarget = this.__3d.orbit.target.clone()
        const endTarget = center.clone()
        const startCam = this.__3d.camera.position.clone()
        const dir = startCam.clone().sub(startTarget).normalize()
        if (dir.length() < 0.01) dir.set(0.5, 0.4, 0.7).normalize()
        const endCam = endTarget.clone().add(dir.multiplyScalar(dist))

        this.__3d._cameraAnim = {
          startTime: performance.now(),
          duration: 600,
          startTarget: startTarget,
          endTarget: endTarget,
          startCam: startCam,
          endCam: endCam
        }
      }
    },
    onCanvasMouseMove(event) {
      const nextObj = this.getIntersectedObject(event)
      const prevId = this._hoveredObjectId
      const nextId = nextObj ? nextObj.id : null
      if (this.$refs.previewCanvas) {
        this.$refs.previewCanvas.style.cursor = this.hasPreviewActions(nextObj) ? 'pointer' : 'default'
      }

      if (prevId !== nextId) {
        if (prevId) {
          const prevObj = this.sceneObjects.find(obj => obj.id === prevId)
          this.executeInteraction(prevObj, 'mouseleave')
        }
        if (nextObj) {
          this.executeInteraction(nextObj, 'mouseenter')
        }
        this._hoveredObjectId = nextId
      }
    },

    setCameraView(view) {
      if (!this.__3d) return
      this.cameraView = view
      const cam = this.__3d.camera
      const orb = this.__3d.orbit
      const d = 10
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

    resetCamera() {
      this.setCameraView('perspective')
    },

    frameAll() {
      if (!this.__3d || this.objectCount === 0) return
      const box = new THREE.Box3()
      Object.values(this.__3d.meshMap).forEach(m => box.expandByObject(m))
      const center = box.getCenter(new THREE.Vector3())
      const size = box.getSize(new THREE.Vector3())
      const maxDim = Math.max(size.x, size.y, size.z)
      this.__3d.orbit.target.copy(center)
      this.__3d.camera.position.set(center.x + maxDim * 1.5, center.y + maxDim, center.z + maxDim * 1.5)
      this.__3d.orbit.update()
    },

    restoreCamera(sceneSettings) {
      if (!this.__3d || !sceneSettings) return
      const cp = sceneSettings.cameraPosition
      const ct = sceneSettings.cameraTarget
      const fov = sceneSettings.cameraFov
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

    _finishSceneLoad() {
      if (this._savedSceneSettings && this._savedSceneSettings.cameraPosition) {
        this.restoreCamera(this._savedSceneSettings)
      } else {
        this.frameAll()
      }
    },

    setupGLTFAnimations(model, objData, clips) {
      if (!model) return
      model.userData = model.userData || {}
      clips = Array.isArray(clips) ? clips : []
      const nameCounts = {}
      const entries = clips.map((clip, index) => {
        const name = clip.name || ('Animation ' + (index + 1))
        nameCounts[name] = (nameCounts[name] || 0) + 1
        return {
          key: name + '__' + (index + 1),
          name,
          label: name
        }
      })
      entries.forEach((entry, index) => {
        if (nameCounts[entry.name] > 1) {
          entry.label = entry.name + ' #' + (index + 1)
        }
      })
      const names = entries.map(entry => entry.key)
      if (!clips.length) {
        model.userData.animationMixer = null
        model.userData.animationActions = {}
        model.userData.animationNames = []
        model.userData.animationEntries = []
        return
      }
      const mixer = new THREE.AnimationMixer(model)
      const actions = {}
      clips.forEach((clip, index) => {
        const name = entries[index].key
        const action = mixer.clipAction(clip)
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
        const savedAnimationName = objData.gltfAnimationName
        const savedAnimationNames = Array.isArray(objData.gltfAnimationNames) ? objData.gltfAnimationNames.slice() : (savedAnimationName ? [savedAnimationName] : [])
        const savedGroups = ensureGLTFAnimationGroups(objData)
        this.$set(objData, 'gltfAnimations', entries)
        const nextAnimationNames = savedAnimationNames.map(savedName => {
          if (actions[savedName]) return savedName
          const matchedEntry = entries.find(entry => entry.name === savedName || entry.label === savedName)
          return matchedEntry ? matchedEntry.key : ''
        }).filter(Boolean)
        if (!nextAnimationNames.length && names.length) nextAnimationNames.push(names[0])
        this.$set(objData, 'gltfAnimationNames', nextAnimationNames)
        this.$set(objData, 'gltfAnimationName', nextAnimationNames[0] || '')
        const nextGroups = remapGLTFAnimationGroupNames(savedGroups, entries, actions, names[0])
        this.$set(objData, 'gltfAnimationGroups', nextGroups)
        syncLegacyGLTFAnimationFields(objData, nextGroups)
      }
    },

    updateGLTFAnimation(mesh, obj, delta) {
      if (!mesh || !obj || !mesh.userData || !mesh.userData.animationMixer) return
      const actions = mesh.userData.animationActions || {}
      const names = mesh.userData.animationNames || Object.keys(actions)
      if (!names.length) return

      const groups = ensureGLTFAnimationGroups(obj, { defaultAnimationKey: names[0] })
      const activeMap = {}
      groups.forEach(group => {
        if (!group.playing) return
        const selectedNames = Array.isArray(group.animationNames) && group.animationNames.length ? group.animationNames : []
        selectedNames.forEach(name => {
          if (actions[name]) activeMap[name] = group
        })
      })
      Object.keys(actions).forEach(name => {
        const group = activeMap[name]
        const action = actions[name]
        if (!group) {
          action.stop()
          return
        }
        const loop = group.loop !== false
        const loopMode = loop ? THREE.LoopRepeat : THREE.LoopOnce
        const speed = group.speed !== undefined && !isNaN(Number(group.speed)) ? Number(group.speed) : 1
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
      this.applyObjectChanges(obj, {
        gltfAnimationGroups: nextGroups,
        ...legacy
      })
    },

    /**
     * 从 base64 恢复 GLTF 模型
     */
    restoreGLTFModels(objects) {
      var self = this
      var gltfCount = 0
      var hasBackgroundModel = false
      objects.forEach(function(obj) {
        if (obj.type !== 'gltf') return
        gltfCount++
        if (obj.isBackground) hasBackgroundModel = true
        if (obj.modelPath || obj.modelUrl) {
          var remoteModelUrl = obj.modelPath || obj.modelUrl
          console.log('预览页面加载GLTF模型:', remoteModelUrl, 'ID:', obj.id)

          var loaderByUrl = new GLTFLoader()
          loaderByUrl.load(remoteModelUrl, function(gltf) {
            console.log('预览页面GLTF模型加载成功:', obj.id)
            var model = gltf.scene
            var box = new THREE.Box3().setFromObject(model)
            var size = box.getSize(new THREE.Vector3())
            var maxDim = Math.max(size.x, size.y, size.z)
            var autoScale = maxDim > 0 ? (10 / maxDim) : 1
            var hasSavedScale = obj.sx !== undefined && obj.sx !== 1
            if (hasSavedScale) {
              model.scale.set(obj.sx, obj.sy, obj.sz)
              model.position.set(obj.x || 0, obj.y || 0, obj.z || 0)
              model.rotation.set(obj.rx || 0, obj.ry || 0, obj.rz || 0)
            } else {
              model.scale.setScalar(autoScale)
              box.setFromObject(model)
              var centerByUrl = box.getCenter(new THREE.Vector3())
              model.position.sub(centerByUrl)
            }
            model.traverse(function(child) {
              if (child.isMesh) {
                child.castShadow = true
                child.receiveShadow = true
                // 只有当对象有自定义材质属性时才应用修改，否则保留原始材质
                if (obj.materialOverridden === true && (obj.color || obj.opacity !== undefined || obj.wireframe || obj.metalness !== undefined || obj.roughness !== undefined || obj.emissive)) {
                  if (child.material) {
                    if (Array.isArray(child.material)) {
                      child.material.forEach(material => self.applyMaterialProps(material, obj))
                    } else {
                      self.applyMaterialProps(child.material, obj)
                    }
                  }
                }
              }
            })
            if (obj.isBackground && model) {
              model.traverse(function(child) {
                if (child.isMesh && child.material) {
                  child.material.transparent = true
                  child.material.opacity = obj.opacity !== undefined ? obj.opacity : 0.5
                  child.material.needsUpdate = true
                }
              })
            }
            self.__3d.scene.add(model)
            model.userData = model.userData || {}
            model.userData.id = obj.id
            model.userData.isMeshGroup = true
            self.__3d.meshMap[obj.id] = model
            model.visible = obj.visible !== false
            self.setupGLTFAnimations(model, obj, gltf.animations)
            gltfCount--
            if (gltfCount <= 0) {
              if (hasBackgroundModel && self.__3d.gridHelper) {
                self.__3d.gridHelper.visible = false
                self.showGrid = false
              }
              self.$nextTick(function() { self._finishSceneLoad() })
            }
          }, undefined, function(error) {
            console.error('预览页面GLTF模型(URL)加载失败:', error, 'URL:', remoteModelUrl, 'ID:', obj.id)
            if (error && error.message && error.message.includes('<!DOCTYPE')) {
              console.error('服务器返回了HTML页面而不是模型文件，请检查URL是否正确:', remoteModelUrl)
            }
            gltfCount--
            if (gltfCount <= 0 && self.objectCount > 0) {
              self.$nextTick(function() { self._finishSceneLoad() })
            }
          })
          return
        } else {
          console.warn('预览页面GLTF模型缺少路径:', obj.id, '名称:', obj.name)
        }
      })

      if (hasBackgroundModel && self.__3d.gridHelper) {
        self.__3d.gridHelper.visible = false
        self.showGrid = false
      }
      if (gltfCount === 0 && self.objectCount > 0) {
        self.$nextTick(function() { self._finishSceneLoad() })
      }
    },

    toggleGrid() {
      if (this.__3d && this.__3d.gridHelper) {
        this.__3d.gridHelper.visible = this.showGrid
      }
    },

    toggleAutoRotate() {
      this.autoRotate = !this.autoRotate
    },

    showContextMenu(event) {
      event.preventDefault()
      const rect = this.$refs.canvasWrapper.getBoundingClientRect()
      this.contextMenu = {
        visible: true,
        x: event.clientX - rect.left,
        y: event.clientY - rect.top
      }
    },

    hideContextMenu() {
      this.contextMenu.visible = false
    },

    handleContextMenuClick(action, param) {
      this.hideContextMenu()
      switch (action) {
        case 'setCameraView':
          this.setCameraView(param)
          break
        case 'resetCamera':
          this.resetCamera()
          break
        case 'frameAll':
          this.frameAll()
          break
        case 'toggleGrid':
          this.showGrid = !this.showGrid
          this.toggleGrid()
          break
        case 'toggleAutoRotate':
          this.toggleAutoRotate()
          break
        case 'toggleFullscreen':
          this.toggleFullscreen()
          break
        case 'fitPage':
          this.fitPage()
          break
      }
    },

    toggleFullscreen() {
      const container = this.$refs.canvasWrapper
      if (!container) return

      if (!document.fullscreenElement) {
        container.requestFullscreen().then(() => {
          this.isFullscreen = true
        }).catch(err => {
          console.error('全屏失败:', err)
        })
      } else {
        document.exitFullscreen().then(() => {
          this.isFullscreen = false
        }).catch(err => {
          console.error('退出全屏失败:', err)
        })
      }
    },

    fitPage() {
      if (!this.__3d || !this.$refs.canvasWrapper) return

      const rect = this.$refs.canvasWrapper.getBoundingClientRect()
      const w = rect.width
      const h = rect.height

      this.__3d.camera.aspect = w / h
      this.__3d.camera.updateProjectionMatrix()
      this.__3d.renderer.setSize(w, h)

      this.frameAll()
    },

    onFullscreenChange() {
      this.isFullscreen = !!document.fullscreenElement
    },

    onResize() {
      const vm = this
      if (!vm.__3d) return
      const wrapper = vm.$refs.canvasWrapper
      if (!wrapper) return
      const rect = wrapper.getBoundingClientRect()
      let w = rect.width
      let h = rect.height
      if (w <= 0 || h <= 0) return
      w = Math.max(w, 100)
      h = Math.max(h, 100)
      vm.__3d.camera.aspect = w / h
      vm.__3d.camera.updateProjectionMatrix()
      vm.__3d.renderer.setSize(w, h)
      vm.update2DOverlays()
    },

    update2DOverlays() {
      if (!this.__3d || !this.$refs.canvasWrapper) {
        this.overlay2DItems = []
        return
      }
      const wrapper = this.$refs.canvasWrapper
      const overlays = []

      for (let i = 0; i < this.sceneObjects.length; i++) {
        const obj = this.sceneObjects[i]
        if (!obj || (obj.type !== '2dComponent' && !isUiOverlayObject(obj)) || obj.visible === false) continue
        const sourceStyle = obj.source2D && obj.source2D.style ? obj.source2D.style : {}
        if (obj.type === '2dComponent' && (sourceStyle.visible === 0 || sourceStyle.visible === false)) continue
        const sourcePos = obj.type === '2dComponent'
          ? (sourceStyle.position || {})
          : {
            x: obj.uiX !== undefined ? obj.uiX : 40,
            y: obj.uiY !== undefined ? obj.uiY : 40,
            w: obj.uiWidth !== undefined ? obj.uiWidth : (obj.type === 'webEmbed' ? 420 : 180),
            h: obj.uiHeight !== undefined ? obj.uiHeight : (obj.type === 'webEmbed' ? 260 : 100)
          }
        let screenX = parseFloat(sourcePos.x) || 0
        let screenY = parseFloat(sourcePos.y) || 0
        const overlayWidth = Math.max(1, parseFloat(sourcePos.w) || 160)
        const overlayHeight = Math.max(1, parseFloat(sourcePos.h) || 80)

        // 计算中心点坐标（左上角 + 半宽高）
        const centerX = screenX + overlayWidth / 2
        const centerY = screenY + overlayHeight / 2

        // 获取旋转角度
        const transform = parseFloat(obj.type === '2dComponent' ? sourceStyle.transform : obj.uiRotation)
        const rotation = (isNaN(transform) || transform === -1098 || transform === -1099) ? 0 : transform

        const overlay = {
          id: obj.id,
          objectData: obj,
          kind: obj.type,
          style: {
            width: overlayWidth + 'px',
            height: overlayHeight + 'px',
            transform: `translate(-50%, -50%) rotate(${rotation}deg)`,
            left: centerX + 'px',
            top: centerY + 'px',
            display: obj.visible === false ? 'none' : 'block'
          },
          hasAction: this.hasPreviewActions(obj)
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
        this.applyObjectChanges(obj, changes)
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
      this.applyObjectChanges(obj, changes)
    },

    /**
     * 数据推送驱动动画
     * 响应 readDataPush / StaticData / SystemData 事件
     * 与 ISM3DEditor.vue 中的 DealWithUpdateData 逻辑保持一致
     * realData 结构：{ Data: [{ ModelDataUuid, Uuid, Value }], DeviceUuid }
     */
    DealWithUpdateData(realData) {
      if (!realData || !realData.Data || !Array.isArray(realData.Data)) return
      const objs = this.sceneObjects
      if (!objs || !objs.length) return

      for (let k = 0, Datalen = realData.Data.length; k < Datalen; k++) {
        for (let j = 0, objLen = objs.length; j < objLen; j++) {
          const obj = objs[j]
          if (!obj) continue

          const currentData = realData.Data[k]
          const identifier = obj.source2D && obj.source2D.identifier ? obj.source2D.identifier : ''
          const realValue = parseFloat(currentData.Value)
          this.executeDataWorkbenchEvents(obj, currentData)
          this.applyGLTFAnimationCondition(obj, realData, currentData)

          // ========== dataText 类型单独处理（不需要 animate.condition） ==========
          if (obj.type === 'dataText') {
            const wsMatched = obj.wsKey && (obj.wsKey === currentData.ModelDataUuid || obj.wsKey === currentData.Uuid || obj.wsKey === currentData.Name)
            const condition = obj.animate && obj.animate.condition
            const dataMatched = condition && condition.dataID && (condition.dataID === currentData.ModelDataUuid || condition.dataID === currentData.Uuid)
            if (wsMatched || dataMatched) {
              this.applyObjectChanges(obj, { realTimeValue: currentData.Value })
            }
            continue
          }

          this.applyPositionBindings(obj, realData, currentData)
          this.applyGenericDataBinding(obj, realData, currentData)

          // ========== 2D组件数据推送处理（独立于 animate.condition） ==========
          if (obj.type === '2dComponent' && obj.source2D) {
            const identifier = obj.source2D.identifier || ''
            if (identifier) {
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

                  // 处理2D组件的动画事件（animateEvent）
              const animateCondition = obj.source2D.animate && obj.source2D.animate.condition ? obj.source2D.animate.condition : null
              const animateSelected = obj.source2D.animate && Array.isArray(obj.source2D.animate.selected) ? obj.source2D.animate.selected : []
              if (animateCondition && animateCondition.dataID && animateSelected.length > 0) {
                const deviceMatched = animateCondition.isBandDevice ? realData.DeviceUuid === animateCondition.deviceSN : true
                const dataMatched = animateCondition.dataID === currentData.ModelDataUuid || animateCondition.dataID === currentData.Uuid
                if (deviceMatched && dataMatched) {
                  let isStart = false
                  if (obj.source2D.animate.isExpression && animateCondition.operator) {
                    const RealValue = parseFloat(currentData.Value)
                    const OperatorValue = parseFloat(animateCondition.OperatorValue)
                    if (!isNaN(OperatorValue) && !isNaN(RealValue)) {
                      switch (animateCondition.operator) {
                        case "==":  isStart = (RealValue == OperatorValue); break
                        case "!=":  isStart = (RealValue != OperatorValue); break
                        case ">":   isStart = (RealValue > OperatorValue); break
                        case "<":   isStart = (RealValue < OperatorValue); break
                        case ">=":  isStart = (RealValue >= OperatorValue); break
                        case "<=":  isStart = (RealValue <= OperatorValue); break
                        case "<>": {
                          const OperatorMaxValue = parseFloat(animateCondition.OperatorMaxValue)
                          if (!isNaN(OperatorMaxValue)) {
                            isStart = (RealValue >= OperatorValue && RealValue <= OperatorMaxValue)
                          }
                        }
                        break
                        case "<!>": {
                          const OperatorMaxValue = parseFloat(animateCondition.OperatorMaxValue)
                          if (!isNaN(OperatorMaxValue)) {
                            isStart = (RealValue < OperatorValue || RealValue > OperatorMaxValue)
                          }
                        }
                        break
                      }
                    }
                  } else {
                    isStart = true
                  }
                  this.$EventBus.$emit(identifier + 'animateEvent', isStart)
                }
              }
            }
          }

          // 其他类型需要 animate.condition
          if (!obj.animate || !obj.animate.condition) continue

          const condition = obj.animate.condition
          const selectAnimate = obj.animate.selected || []
          // ========== 情况 1：绑定设备，优先匹配设备 ==========
          if (condition.isBandDevice) {
            // 设备SN不匹配则跳过
            if (realData.DeviceUuid !== condition.deviceSN) continue
            // 如果配置了 dataID，则还需 dataID 匹配；否则只要设备匹配就处理
            if (condition.dataID) {
              if (condition.dataID !== currentData.ModelDataUuid) continue
            }
            let isStart = true
            if (condition.operator) {
              const RealValue = parseFloat(currentData.Value)
              const OperatorValue = parseFloat(condition.OperatorValue)
              if (isNaN(OperatorValue)) continue
              if (isNaN(RealValue)) continue

              isStart = false
              switch (condition.operator) {
                case "==":  isStart = (RealValue == OperatorValue); break
                case "!=":  isStart = (RealValue != OperatorValue); break
                case ">":   isStart = (RealValue > OperatorValue); break
                case "<":   isStart = (RealValue < OperatorValue); break
                case ">=":  isStart = (RealValue >= OperatorValue); break
                case "<=":  isStart = (RealValue <= OperatorValue); break
                case "<>":
                  {
                    const OperatorMaxValue = parseFloat(condition.OperatorMaxValue)
                    if (!isNaN(OperatorMaxValue)) {
                      isStart = (RealValue >= OperatorValue && RealValue <= OperatorMaxValue)
                    }
                  }
                  break
                case "<!>":
                  {
                    const OperatorMaxValue = parseFloat(condition.OperatorMaxValue)
                    if (!isNaN(OperatorMaxValue)) {
                      isStart = (RealValue < OperatorValue || RealValue > OperatorMaxValue)
                    }
                  }
                  break
              }
            }
            if (selectAnimate && selectAnimate.length > 0) {
              const changes = {}
              selectAnimate.forEach(animId => {
                switch (animId) {
                    case 'autoRotate': changes.autoRotate = isStart; break
                    case 'floatAnim':  changes.floatAnim = isStart; break
                    case 'blink':      changes.blink = isStart; break
                    case 'visible':    changes.visible = isStart; break
                }
              })
              this.applyObjectChanges(obj, changes)
              if (obj.type === '2dComponent' && identifier) {
                this.$EventBus.$emit(identifier + 'animateEvent', isStart)
              }
            }
          }

          // ========== 情况2：未绑定设备，直接匹配数据ID ==========
          if (!condition.isBandDevice) {
            if (!condition.dataID) continue
            if (condition.dataID !== currentData.ModelDataUuid &&
                condition.dataID !== currentData.Uuid) continue

            let isStart = true
            if (condition.operator) {
              const RealValue = parseFloat(currentData.Value)
              const OperatorValue = parseFloat(condition.OperatorValue)
              if (isNaN(OperatorValue)) continue
              if (isNaN(RealValue)) continue

              isStart = false
              switch (condition.operator) {
                case "==":  isStart = (RealValue == OperatorValue); break
                case "!=":  isStart = (RealValue != OperatorValue); break
                case ">":   isStart = (RealValue > OperatorValue); break
                case "<":   isStart = (RealValue < OperatorValue); break
                case ">=":  isStart = (RealValue >= OperatorValue); break
                case "<=":  isStart = (RealValue <= OperatorValue); break
                case "<>":
                  {
                    const OperatorMaxValue = parseFloat(condition.OperatorMaxValue)
                    if (!isNaN(OperatorMaxValue)) {
                      isStart = (RealValue >= OperatorValue && RealValue <= OperatorMaxValue)
                    }
                  }
                  break
                case "<!>":
                  {
                    const OperatorMaxValue = parseFloat(condition.OperatorMaxValue)
                    if (!isNaN(OperatorMaxValue)) {
                      isStart = (RealValue < OperatorValue || RealValue > OperatorMaxValue)
                    }
                  }
                  break
              }
            }

            if (selectAnimate && selectAnimate.length > 0) {
              const changes = {}
              selectAnimate.forEach(animId => {
                switch (animId) {
                case 'autoRotate': changes.autoRotate = isStart; break
                case 'floatAnim':  changes.floatAnim = isStart; break
                case 'blink':      changes.blink = isStart; break
                case 'visible':    changes.visible = isStart; break
                }
              })
              this.applyObjectChanges(obj, changes)
              if (obj.type === '2dComponent' && identifier) {
                this.$EventBus.$emit(identifier + 'animateEvent', isStart)
              }
            }
          }
        }
      }
      // 触发图表组件的数据更新完成事件（在所有数据处理完成后触发）
      this.$EventBus.$emit('DealWithRealDataFinish', realData)
    },

    /**
     * 应用对象属性变更，触发 Vue 响应式更新
     * 与 ISM3DEditor.vue 中的 updateObject 逻辑保持一致
     */
    applyObjectChanges(obj, changes) {
      if (!obj || !changes) return
      const mesh = this.__3d && this.__3d.meshMap && this.__3d.meshMap[obj.id]
      Object.keys(changes).forEach(key => {
        this.$set(obj, key, changes[key])
      })
      if (obj.type === '2dComponent' && obj.source2D && obj.source2D.style && changes.visible !== undefined) {
        this.$set(obj.source2D.style, 'visible', changes.visible ? 1 : 0)
      }
      // 替换数组引用，确保 Vue 响应式更新，与 ISM3DEditor.vue 中的 updateObject 保持一致
      this.sceneObjects = [...this.sceneObjects]
      if (mesh) {
        if (changes.blink !== undefined && mesh.material) {
          mesh.material.transparent = obj.blink
          mesh.material.needsUpdate = true
        }
        if (changes.visible !== undefined) {
          mesh.visible = obj.visible
        }
        if (obj.type === 'dataText' && mesh.userData && mesh.userData.isTextSprite) {
          const textColor = parseInt(normTextLabelColor(obj.color).replace('#', ''), 16)
          if (mesh.userData.is3DText) {
            updateText3DMesh(mesh, getTextLabelContent(obj), textColor, obj.fontSize || 16, obj.textBgColor || '#00000000', buildTextLabelOptions(obj))
          } else {
            updateTextSprite(mesh, getTextLabelContent(obj), textColor, obj.fontSize || 16, obj.textBgColor || '#00000000', buildTextLabelOptions(obj))
          }
          applyTextSpriteObjectScale(mesh, obj)
          updateTextPickArea(mesh, obj)
        }
      }
    }
  }
}
</script>

<style scoped>
.scene-preview-page {
  width: 100%;
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #ffffff;
  overflow: hidden;
  margin: 0;
  padding: 0;
}

.preview-toolbar {
  height: 48px;
  background: #fff;
  border-bottom: 1px solid #e8e8e8;
  display: flex;
  align-items: center;
  padding: 0 12px;
  gap: 8px;
  flex-shrink: 0;
  z-index: 100;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.toolbar-title {
  font-size: 15px;
  font-weight: 600;
  color: #333;
  white-space: nowrap;
}

.toolbar-center {
  display: flex;
  align-items: center;
  gap: 4px;
  flex: 1;
  justify-content: center;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.info-item {
  font-size: 12px;
  color: #999;
}

.toolbar-divider {
  width: 1px;
  height: 24px;
  background: #e8e8e8;
  margin: 0 6px;
}

.tb-btn {
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
.tb-btn:hover { background: #e6fffb; border-color: #87e8de; color: #13c2c2; }
.tb-btn.active { background: #e6fffb; border-color: #13c2c2; color: #13c2c2; }
.tb-btn i { font-size: 13px; }

.preview-canvas-wrapper {
  flex: 1;
  position: relative;
  overflow: hidden;
  width: 100%;
  height: 100%;
  margin: 0;
  padding: 0;
}

.context-menu {
  position: absolute;
  background: #fff;
  border-radius: 6px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  border: 1px solid #e8e8e8;
  padding: 4px 0;
  min-width: 140px;
  z-index: 1000;
}

.context-menu-title {
  padding: 6px 12px;
  font-size: 12px;
  color: #999;
  border-bottom: 1px solid #f0f0f0;
  margin-bottom: 4px;
}

.context-menu-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  width: 100%;
  border: none;
  background: transparent;
  cursor: pointer;
  font-size: 13px;
  color: #333;
  text-align: left;
  transition: background-color 0.15s;
}

.context-menu-item:hover {
  background: #f5f5f5;
}

.context-menu-item.active {
  background: #e6fffb;
  color: #13c2c2;
}

.context-menu-item i {
  font-size: 13px;
  width: 16px;
}

.context-menu-divider {
  height: 1px;
  background: #f0f0f0;
  margin: 4px 0;
}

.overlay-2d-layer {
  position: absolute;
  inset: 0;
  pointer-events: none;
  overflow: hidden;
  z-index: 5;
}

.overlay-2d-item {
  position: absolute;
  transform-origin: center center;
  pointer-events: none;
}

.overlay-2d-content {
  width: 100%;
  height: 100%;
  pointer-events: auto;
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
}

#preview-canvas {
  display: block;
  width: 100%;
  height: 100%;
  margin: 0;
  padding: 0;
}

.empty-hint {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  text-align: center;
  pointer-events: none;
}

.loading-mask {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.85);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 50;
}

.axis-indicator {
  position: absolute;
  bottom: 50px;
  right: 14px;
  width: 70px;
  height: 70px;
  pointer-events: none;
}

.preview-popup-mask {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.35);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 80;
}

.preview-popup-window {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.25);
  display: flex;
  flex-direction: column;
  max-width: calc(100% - 40px);
  max-height: calc(100% - 40px);
}

.preview-popup-header {
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 12px;
  background: #fafafa;
  border-bottom: 1px solid #e8e8e8;
  color: #333;
  font-size: 13px;
  flex-shrink: 0;
}

.preview-popup-close {
  border: none;
  background: transparent;
  cursor: pointer;
  color: #999;
  font-size: 16px;
}

.preview-popup-frame {
  width: 100%;
  flex: 1;
  border: none;
  background: #fff;
}

.preview-statusbar {
  height: 22px;
  background: #fafafa;
  display: flex;
  align-items: center;
  padding: 0 10px;
  font-size: 11px;
  color: #333;
  gap: 16px;
  border-top: 1px solid #e8e8e8;
}
</style>

<style>
body {
  margin: 0;
  padding: 0;
  overflow: hidden;
}

html {
  margin: 0;
  padding: 0;
  overflow: hidden;
}

#app {
  margin: 0;
  padding: 0;
  overflow: hidden;
}

/* Edge浏览器滚动条隐藏 */
::-webkit-scrollbar {
  display: none;
}

/* 针对Edge浏览器的特殊处理 */
@supports (-ms-ime-align: auto) {
  body {
    overflow: hidden;
  }

  html {
    overflow: hidden;
  }
}
</style>
