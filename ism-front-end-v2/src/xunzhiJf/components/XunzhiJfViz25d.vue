<template>
  <div ref="root" class="xunzhi-jf-viz25d">
    <canvas ref="canvas" class="xunzhi-jf-viz25d__canvas" />
    <div class="xunzhi-jf-viz25d__hud">
      <span class="xunzhi-jf-viz25d__hint">拖拽旋转 · 滚轮缩放 · 冷/热通道 · 程序化机柜（无 BIM）</span>
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
const GAP_COL = 0.32
const GAP_ROW_HOT = 0.22
const GAP_ROW_COLD = 1.05

/** 假数据：每层若干机柜标红 */
const FLOOR_DEFS = [
  { id: 'hall-a', label: 'A 区主机房', rows: 3, cols: 5, alarms: [[0, 2], [1, 4], [2, 1]] },
  { id: 'hall-b', label: 'B 区配电间', rows: 2, cols: 4, alarms: [[0, 1], [1, 3]] },
  { id: 'hall-c', label: 'C 区电池间', rows: 2, cols: 3, alarms: [[1, 0]] },
]

function rowGapBefore(rowIndex) {
  if (rowIndex <= 0) return 0
  return rowIndex % 2 === 1 ? GAP_ROW_COLD : GAP_ROW_HOT
}

function computeRowZs(rows) {
  const zs = []
  let z = 0
  for (let r = 0; r < rows; r++) {
    zs.push(z)
    if (r < rows - 1) z += CAB_D + rowGapBefore(r + 1)
  }
  const mid = (zs[0] + zs[rows - 1]) / 2
  return zs.map(v => v - mid)
}

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

      const frameMat = new THREE.MeshStandardMaterial({
        color: 0x3d5568,
        metalness: 0.5,
        roughness: 0.48,
      })
      const doorMat = new THREE.MeshStandardMaterial({
        color: 0x5a7a98,
        metalness: 0.42,
        roughness: 0.44,
      })
      const seamMat = new THREE.MeshStandardMaterial({
        color: 0x2a3848,
        metalness: 0.35,
        roughness: 0.65,
      })
      const postMat = new THREE.MeshStandardMaterial({
        color: 0x6a8498,
        metalness: 0.62,
        roughness: 0.35,
      })
      const ventMat = new THREE.MeshStandardMaterial({
        color: 0x6a90c0,
        metalness: 0.55,
        roughness: 0.38,
      })
      const ledOkMat = new THREE.MeshStandardMaterial({
        color: 0x33dd66,
        emissive: 0x22cc55,
        emissiveIntensity: 0.55,
        metalness: 0.1,
        roughness: 0.4,
      })
      const alarmMat = new THREE.MeshStandardMaterial({
        color: 0xff4455,
        emissive: 0xff2233,
        emissiveIntensity: 1.05,
        metalness: 0.1,
        roughness: 0.3,
      })
      const coldAisleMat = new THREE.MeshStandardMaterial({
        color: 0x9ec8e8,
        metalness: 0.05,
        roughness: 0.92,
        emissive: 0x4488aa,
        emissiveIntensity: 0.08,
      })
      const hotAisleMat = new THREE.MeshStandardMaterial({
        color: 0x2e4050,
        metalness: 0.12,
        roughness: 0.88,
      })
      const pillarMat = new THREE.MeshStandardMaterial({
        color: 0x7a94a8,
        metalness: 0.45,
        roughness: 0.5,
      })
      const lightStripMat = new THREE.MeshStandardMaterial({
        color: 0xe8f4ff,
        emissive: 0xc8e8ff,
        emissiveIntensity: 0.65,
        metalness: 0.1,
        roughness: 0.35,
      })

      const rowZs = computeRowZs(def.rows)
      const offsetX = ((def.cols - 1) * (CAB_W + GAP_COL)) / 2
      const alarmSet = new Set(def.alarms.map(([r, c]) => `${r},${c}`))
      const hallWidth = def.cols * (CAB_W + GAP_COL) - GAP_COL + 1.2
      const hallDepth =
        rowZs[def.rows - 1] - rowZs[0] + CAB_D + 1.2

      this.addAisleFloors(group, def, rowZs, hallWidth, coldAisleMat, hotAisleMat)
      this.addCeilingStrips(group, def, rowZs, hallWidth, lightStripMat)
      this.addRoomPillars(group, hallWidth, hallDepth, pillarMat)

      for (let r = 0; r < def.rows; r++) {
        const faceSign = r % 2 === 0 ? 1 : -1
        for (let c = 0; c < def.cols; c++) {
          const x = c * (CAB_W + GAP_COL) - offsetX
          const z = rowZs[r]
          const isAlarm = alarmSet.has(`${r},${c}`)
          this.addCabinet(group, x, z, faceSign, isAlarm, {
            frameMat,
            doorMat,
            seamMat,
            postMat,
            ventMat,
            ledOkMat,
            alarmMat,
          })
        }
      }

      return group
    },

    addCabinet(group, x, z, faceSign, isAlarm, mats) {
      const yBase = CAB_H / 2
      const postT = 0.045
      const innerW = CAB_W - postT * 2
      const innerD = CAB_D - postT * 2

      const shellGeo = new THREE.BoxGeometry(innerW, CAB_H * 0.98, innerD)
      const shell = new THREE.Mesh(shellGeo, mats.frameMat)
      shell.position.set(x, yBase, z)
      group.add(shell)

      const posts = [
        [-innerW / 2, innerD / 2],
        [innerW / 2, innerD / 2],
        [-innerW / 2, -innerD / 2],
        [innerW / 2, -innerD / 2],
      ]
      posts.forEach(([px, pz]) => {
        const geo = new THREE.BoxGeometry(postT, CAB_H, postT)
        const post = new THREE.Mesh(geo, mats.postMat)
        post.position.set(x + px, yBase, z + pz)
        group.add(post)
      })

      const doorTh = 0.035
      const doorGeo = new THREE.BoxGeometry(innerW * 0.9, CAB_H * 0.92, doorTh)
      const door = new THREE.Mesh(doorGeo, mats.doorMat)
      door.position.set(x, yBase, z + faceSign * (innerD / 2 + doorTh * 0.45))
      group.add(door)

      const seamW = 0.018
      const seamH = CAB_H * 0.86
      ;[-innerW * 0.22, innerW * 0.22].forEach(dx => {
        const seamGeo = new THREE.BoxGeometry(seamW, seamH, doorTh * 0.6)
        const seam = new THREE.Mesh(seamGeo, mats.seamMat)
        seam.position.set(x + dx, yBase, door.position.z + faceSign * 0.008)
        group.add(seam)
      })

      const ventGeo = new THREE.BoxGeometry(innerW * 0.88, 0.05, innerD * 0.2)
      const vent = new THREE.Mesh(ventGeo, mats.ventMat)
      vent.position.set(x, CAB_H * 0.74, z + faceSign * (innerD * 0.35))
      group.add(vent)

      const ledGeo = new THREE.BoxGeometry(innerW * 0.55, 0.04, 0.06)
      const ledMat = isAlarm ? mats.alarmMat : mats.ledOkMat
      const led = new THREE.Mesh(ledGeo, ledMat)
      led.position.set(x, CAB_H * 0.9, z + faceSign * (innerD / 2 + 0.05))
      group.add(led)

      if (isAlarm) {
        const lampGeo = new THREE.BoxGeometry(innerW * 0.2, CAB_H * 0.1, 0.07)
        const lamp = new THREE.Mesh(lampGeo, mats.alarmMat)
        lamp.position.set(x + innerW * 0.28, CAB_H * 0.82, door.position.z + faceSign * 0.02)
        group.add(lamp)
      }
    },

    addAisleFloors(group, def, rowZs, hallWidth, coldMat, hotMat) {
      for (let r = 0; r < def.rows - 1; r++) {
        const gap = rowGapBefore(r + 1)
        const zMid = (rowZs[r] + rowZs[r + 1]) / 2
        const isCold = r % 2 === 0
        const mat = isCold ? coldMat : hotMat
        const depth = isCold ? gap * 0.92 : gap * 0.85
        const geo = new THREE.PlaneGeometry(hallWidth, depth)
        const strip = new THREE.Mesh(geo, mat)
        strip.rotation.x = -Math.PI / 2
        strip.position.set(0, 0.008, zMid)
        group.add(strip)
      }
    },

    addCeilingStrips(group, def, rowZs, offsetX, hallWidth, lightMat) {
      const ceilY = CAB_H + 0.55
      const geo = new THREE.PlaneGeometry(hallWidth + 0.4, 0.35)
      const ceil = new THREE.Mesh(geo, new THREE.MeshStandardMaterial({
        color: 0x4a6888,
        metalness: 0.2,
        roughness: 0.85,
        transparent: true,
        opacity: 0.35,
        side: THREE.DoubleSide,
      }))
      ceil.rotation.x = Math.PI / 2
      ceil.position.set(0, ceilY, rowZs[Math.floor(def.rows / 2)])
      group.add(ceil)

      for (let r = 0; r < def.rows - 1; r++) {
        if (r % 2 !== 0) continue
        const gap = rowGapBefore(r + 1)
        const zMid = (rowZs[r] + rowZs[r + 1]) / 2
        const stripGeo = new THREE.BoxGeometry(hallWidth * 0.85, 0.06, gap * 0.35)
        const strip = new THREE.Mesh(stripGeo, lightMat)
        strip.position.set(0, ceilY - 0.12, zMid)
        group.add(strip)
      }
    },

    addRoomPillars(group, hallWidth, hallDepth, mat) {
      const h = CAB_H + 0.35
      const hw = hallWidth / 2 + 0.35
      const hd = hallDepth / 2 + 0.25
      const corners = [
        [-hw, -hd],
        [hw, -hd],
        [-hw, hd],
        [hw, hd],
      ]
      corners.forEach(([px, pz]) => {
        const geo = new THREE.BoxGeometry(0.14, h, 0.14)
        const pillar = new THREE.Mesh(geo, mat)
        pillar.position.set(px, h / 2, pz)
        group.add(pillar)
      })
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
