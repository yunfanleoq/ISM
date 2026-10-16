<template>
  <div ref="root" class="xunzhi-jf-viz25d">
    <canvas ref="canvas" class="xunzhi-jf-viz25d__canvas" />
    <div class="xunzhi-jf-viz25d__hud">
      <span class="xunzhi-jf-viz25d__hint">拖拽旋转 · 滚轮缩放 · 程序化机柜（无 BIM）</span>
      <span v-if="alarmCount" class="xunzhi-jf-viz25d__alarm">告警示意 {{ alarmCount }} 柜</span>
    </div>
  </div>
</template>

<script>
import * as THREE from 'three'
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls'

const CAB_W = 0.6
const CAB_H = 2
const CAB_D = 1
const GAP = 0.35

/** 假数据：每层若干机柜标红 */
const FLOOR_DEFS = [
  { id: 'hall-a', label: 'A 区主机房', rows: 3, cols: 5, alarms: [[0, 2], [1, 4], [2, 1]] },
  { id: 'hall-b', label: 'B 区配电间', rows: 2, cols: 4, alarms: [[0, 1], [1, 3]] },
  { id: 'hall-c', label: 'C 区电池间', rows: 2, cols: 3, alarms: [[1, 0]] },
]

export default {
  name: 'XunzhiJfViz25d',
  props: {
    floorId: {
      type: String,
      default: 'hall-a',
    },
  },
  data() {
    return {
      alarmCount: 0,
    }
  },
  watch: {
    floorId(id) {
      this.applyFloor(id)
    },
  },
  mounted() {
    this.initScene()
    this.applyFloor(this.floorId)
    window.addEventListener('resize', this.onResize)
  },
  beforeDestroy() {
    this.teardown()
  },
  methods: {
    initScene() {
      const root = this.$refs.root
      const canvas = this.$refs.canvas
      if (!root || !canvas) return

      const w = root.clientWidth || 800
      const h = root.clientHeight || 480

      this._renderer = new THREE.WebGLRenderer({ canvas, antialias: true, alpha: true })
      this._renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, 2))
      this._renderer.setSize(w, h, false)
      const sky = 0x1a3a5c
      this._renderer.setClearColor(sky, 1)

      this._scene = new THREE.Scene()
      this._scene.fog = new THREE.Fog(sky, 22, 50)

      this._camera = new THREE.PerspectiveCamera(48, w / h, 0.1, 200)
      this._camera.position.set(5.5, 5.2, 6.5)

      this._controls = new OrbitControls(this._camera, canvas)
      this._controls.enableDamping = true
      this._controls.dampingFactor = 0.06
      this._controls.target.set(0, 0.8, 0)
      this._controls.maxPolarAngle = Math.PI * 0.48
      this._controls.minDistance = 2.5
      this._controls.maxDistance = 22
      this._controls.update()

      const amb = new THREE.AmbientLight(0xf0f4ff, 0.92)
      this._scene.add(amb)
      const dir = new THREE.DirectionalLight(0xffffff, 1.2)
      dir.position.set(6, 14, 8)
      this._scene.add(dir)
      const fill = new THREE.DirectionalLight(0x88b8ff, 0.42)
      fill.position.set(-8, 6, -6)
      this._scene.add(fill)

      const floorMat = new THREE.MeshStandardMaterial({
        color: 0x3a5a7a,
        metalness: 0.18,
        roughness: 0.78,
      })
      const floorGeo = new THREE.PlaneGeometry(24, 18)
      const floor = new THREE.Mesh(floorGeo, floorMat)
      floor.rotation.x = -Math.PI / 2
      floor.position.y = -0.01
      this._scene.add(floor)

      const grid = new THREE.GridHelper(22, 22, 0x5a8ab8, 0x3a6088)
      grid.position.y = 0.002
      this._scene.add(grid)

      this._floorGroups = {}
      FLOOR_DEFS.forEach(def => {
        this._floorGroups[def.id] = this.buildFloorGroup(def)
        this._scene.add(this._floorGroups[def.id])
      })

      this._animId = requestAnimationFrame(this.tick)
    },

    buildFloorGroup(def) {
      const group = new THREE.Group()
      group.visible = false
      const cabBodyMat = new THREE.MeshStandardMaterial({
        color: 0x4a6580,
        metalness: 0.4,
        roughness: 0.5,
      })
      const cabEdgeMat = new THREE.MeshStandardMaterial({
        color: 0x6a90c0,
        metalness: 0.55,
        roughness: 0.38,
      })
      const alarmMat = new THREE.MeshStandardMaterial({
        color: 0xff4455,
        emissive: 0xff2233,
        emissiveIntensity: 1.05,
        metalness: 0.1,
        roughness: 0.3,
      })

      const alarmSet = new Set(def.alarms.map(([r, c]) => `${r},${c}`))
      const offsetX = ((def.cols - 1) * (CAB_W + GAP)) / 2
      const offsetZ = ((def.rows - 1) * (CAB_D + GAP)) / 2

      for (let r = 0; r < def.rows; r++) {
        for (let c = 0; c < def.cols; c++) {
          const bodyGeo = new THREE.BoxGeometry(CAB_W, CAB_H, CAB_D)
          const mesh = new THREE.Mesh(bodyGeo, cabBodyMat)
          mesh.position.set(
            c * (CAB_W + GAP) - offsetX,
            CAB_H / 2,
            r * (CAB_D + GAP) - offsetZ
          )
          group.add(mesh)

          const stripeGeo = new THREE.BoxGeometry(CAB_W * 0.92, 0.06, CAB_D * 0.15)
          const stripe = new THREE.Mesh(stripeGeo, cabEdgeMat)
          stripe.position.copy(mesh.position)
          stripe.position.y = CAB_H * 0.72
          group.add(stripe)

          if (alarmSet.has(`${r},${c}`)) {
            const lampGeo = new THREE.BoxGeometry(CAB_W * 0.25, CAB_H * 0.12, 0.08)
            const lamp = new THREE.Mesh(lampGeo, alarmMat)
            lamp.position.copy(mesh.position)
            lamp.position.y = CAB_H * 0.88
            lamp.position.z += CAB_D * 0.52
            group.add(lamp)
          }
        }
      }
      return group
    },

    applyFloor(id) {
      const def = FLOOR_DEFS.find(f => f.id === id) || FLOOR_DEFS[0]
      Object.keys(this._floorGroups || {}).forEach(key => {
        const g = this._floorGroups[key]
        if (g) g.visible = key === def.id
      })
      this.alarmCount = def.alarms.length
      if (this._controls) {
        this._controls.target.set(0, 0.8, 0)
        this._controls.update()
      }
    },

    tick() {
      if (!this._renderer || !this._scene || !this._camera) return
      if (this._controls) this._controls.update()
      this._renderer.render(this._scene, this._camera)
      this._animId = requestAnimationFrame(this.tick)
    },

    onResize() {
      const root = this.$refs.root
      if (!root || !this._renderer || !this._camera) return
      const w = root.clientWidth
      const h = root.clientHeight
      if (w < 1 || h < 1) return
      this._camera.aspect = w / h
      this._camera.updateProjectionMatrix()
      this._renderer.setSize(w, h, false)
    },

    teardown() {
      window.removeEventListener('resize', this.onResize)
      if (this._animId) {
        cancelAnimationFrame(this._animId)
        this._animId = null
      }
      if (this._controls) {
        this._controls.dispose()
        this._controls = null
      }
      if (this._scene) {
        const geos = new Set()
        const mats = new Set()
        this._scene.traverse(obj => {
          if (obj.geometry) geos.add(obj.geometry)
          if (obj.material) {
            const list = Array.isArray(obj.material) ? obj.material : [obj.material]
            list.forEach(m => mats.add(m))
          }
        })
        geos.forEach(g => g.dispose())
        mats.forEach(m => m.dispose())
        this._scene = null
      }
      if (this._renderer) {
        this._renderer.dispose()
        this._renderer.forceContextLoss()
        this._renderer = null
      }
      this._floorGroups = null
      this._camera = null
    },
  },
}
</script>

<style scoped>
.xunzhi-jf-viz25d {
  position: relative;
  width: 100%;
  height: 100%;
  min-height: 420px;
  background: #1a3a5c;
  border-radius: 4px;
  overflow: hidden;
}
.xunzhi-jf-viz25d__canvas {
  display: block;
  width: 100%;
  height: 100%;
}
.xunzhi-jf-viz25d__hud {
  position: absolute;
  left: 12px;
  bottom: 10px;
  display: flex;
  gap: 16px;
  font-size: 12px;
  color: rgba(230, 240, 255, 0.88);
  pointer-events: none;
}
.xunzhi-jf-viz25d__alarm {
  color: #ff6b7a;
}
</style>
