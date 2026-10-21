<template>
  <div ref="root" class="demo-stage">
    <canvas ref="canvas" />
  </div>
</template>

<script>
import * as THREE from 'three'
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls'
import { rackWorld } from './layout'

const EASE = [0.4, 0, 0.2, 1]

function bezier(t, p1, p2) {
  const u = 1 - t
  return 3 * u * u * t * p1 + 3 * u * t * t * p2 + t * t * t
}

function ease(t) {
  let lo = 0
  let hi = 1
  for (let i = 0; i < 8; i++) {
    const mid = (lo + hi) / 2
    if (bezier(mid, EASE[0], EASE[2]) < t) lo = mid
    else hi = mid
  }
  return bezier((lo + hi) / 2, EASE[1], EASE[3])
}

export default {
  name: 'DemoStage',
  props: {
    state: { type: String, default: 'cockpit' },
    racks: { type: Array, default: () => [] },
    interactive: { type: Boolean, default: false },
  },
  watch: {
    state() {
      this.applyVisibility()
    },
    interactive(on) {
      if (this._controls) this._controls.enabled = !!on
    },
  },
  mounted() {
    this.initScene()
    this.applyVisibility()
    window.addEventListener('resize', this.onResize)
  },
  beforeDestroy() {
    window.removeEventListener('resize', this.onResize)
    if (this._anim) cancelAnimationFrame(this._anim)
    if (this._controls) this._controls.dispose()
    if (this._renderer) this._renderer.dispose()
  },
  methods: {
    initScene() {
      const root = this.$refs.root
      const canvas = this.$refs.canvas
      const w = root.clientWidth || 960
      const h = root.clientHeight || 540
      this._renderer = new THREE.WebGLRenderer({ canvas, antialias: true, alpha: true })
      this._renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, 2))
      this._renderer.setSize(w, h, false)
      this._renderer.setClearColor(0x0e1c2e, 1)
      this._scene = new THREE.Scene()
      this._scene.fog = new THREE.Fog(0x0e1c2e, 28, 60)
      this._camera = new THREE.PerspectiveCamera(42, w / h, 0.1, 200)
      this._controls = new OrbitControls(this._camera, canvas)
      this._controls.enableDamping = true
      this._controls.dampingFactor = 0.08
      this._controls.enablePan = false
      this._controls.enabled = false
      this._controls.maxPolarAngle = Math.PI * 0.48
      this._scene.add(new THREE.AmbientLight(0xd6e8ff, 0.72))
      const sun = new THREE.DirectionalLight(0xffffff, 1.15)
      sun.position.set(8, 16, 10)
      this._scene.add(sun)
      const fill = new THREE.DirectionalLight(0x5cc8ff, 0.35)
      fill.position.set(-10, 6, -4)
      this._scene.add(fill)
      const ground = new THREE.Mesh(
        new THREE.CircleGeometry(26, 48),
        new THREE.MeshStandardMaterial({ color: 0x102033, metalness: 0.2, roughness: 0.85 })
      )
      ground.rotation.x = -Math.PI / 2
      ground.position.y = -0.02
      this._scene.add(ground)
      this._scene.add(new THREE.GridHelper(30, 30, 0x1d4668, 0x16344e))
      this._building = this.buildBuilding()
      this._floor = this.buildFloorPlates()
      this._room = this.buildRoom()
      this._asset = this.buildAssetRack()
      this._scene.add(this._building, this._floor, this._room, this._asset)
      canvas.addEventListener('pointerdown', this.onPointer)
      this.setPose({ pos: [16, 11, 20], target: [0, 3.2, 0] })
      this._anim = requestAnimationFrame(this.tick)
    },
    buildBuilding() {
      const group = new THREE.Group()
      const bodyMat = new THREE.MeshStandardMaterial({
        color: 0x1a3d58, metalness: 0.42, roughness: 0.38,
      })
      const body = new THREE.Mesh(new THREE.BoxGeometry(8.2, 9.2, 5.6), bodyMat)
      body.position.y = 4.6
      body.userData.pick = 'building'
      group.add(body)
      const edges = new THREE.LineSegments(
        new THREE.EdgesGeometry(body.geometry),
        new THREE.LineBasicMaterial({ color: 0x7ec8ea })
      )
      edges.position.copy(body.position)
      group.add(edges)
      for (let i = 0; i < 4; i++) {
        const band = new THREE.Mesh(
          new THREE.BoxGeometry(8.28, 0.06, 5.68),
          new THREE.MeshStandardMaterial({
            color: 0x5cc8ff, emissive: 0x1a6a90, emissiveIntensity: 0.6, roughness: 0.4,
          })
        )
        band.position.y = 1.4 + i * 2.15
        band.userData.pick = 'building'
        group.add(band)
      }
      for (let i = 0; i < 6; i++) {
        const ac = new THREE.Mesh(
          new THREE.BoxGeometry(0.9, 0.45, 0.7),
          new THREE.MeshStandardMaterial({ color: 0xc5d5e4, metalness: 0.55, roughness: 0.35 })
        )
        ac.position.set(-2.2 + (i % 3) * 2.2, 9.45, i < 3 ? -1.1 : 1.1)
        ac.userData.pick = 'building'
        group.add(ac)
      }
      ;[-14, 13].forEach((x, idx) => {
        const ghost = new THREE.Mesh(
          new THREE.BoxGeometry(5, 6 + idx, 4),
          new THREE.MeshStandardMaterial({ color: 0x0c1824, metalness: 0.2, roughness: 0.8 })
        )
        ghost.position.set(x, 3, -2)
        group.add(ghost)
      })
      return group
    },
    buildFloorPlates() {
      const group = new THREE.Group()
      const slab = new THREE.Mesh(
        new THREE.BoxGeometry(14, 0.12, 9),
        new THREE.MeshStandardMaterial({ color: 0x143044, metalness: 0.25, roughness: 0.7 })
      )
      slab.position.y = 0.06
      group.add(slab)
      const spots = [
        { x: -3.2, z: -1.8, c: 0x1c6b45, k: 'it' },
        { x: 0.2, z: -1.8, c: 0x1c6b45, k: 'it' },
        { x: 3.4, z: -1.8, c: 0x1c6b45, k: 'it' },
        { x: -3.2, z: 1.8, c: 0x1c6b45, k: 'it' },
        { x: 0.6, z: 1.8, c: 0x1a4e78, k: 'power' },
        { x: 3.6, z: 1.8, c: 0x24506a, k: 'battery' },
      ]
      spots.forEach(s => {
        const room = new THREE.Mesh(
          new THREE.BoxGeometry(2.6, 0.28, 2.2),
          new THREE.MeshStandardMaterial({
            color: s.c, emissive: s.c, emissiveIntensity: 0.25, roughness: 0.55, metalness: 0.15,
          })
        )
        room.position.set(s.x, 0.24, s.z)
        group.add(room)
      })
      group.visible = false
      return group
    },
    buildRoom() {
      const group = new THREE.Group()
      const cold = new THREE.Mesh(
        new THREE.PlaneGeometry(5.2, 0.7),
        new THREE.MeshStandardMaterial({
          color: 0x7ed0e8, emissive: 0x245868, emissiveIntensity: 0.2, roughness: 0.9,
        })
      )
      cold.rotation.x = -Math.PI / 2
      cold.position.y = 0.02
      group.add(cold)
      const frame = new THREE.MeshStandardMaterial({ color: 0x8aa0b4, metalness: 0.62, roughness: 0.32 })
      const door = new THREE.MeshStandardMaterial({ color: 0x243848, metalness: 0.4, roughness: 0.45 })
      const ok = new THREE.MeshStandardMaterial({
        color: 0x3ddc97, emissive: 0x146844, emissiveIntensity: 0.8,
      })
      this._alarmMat = new THREE.MeshStandardMaterial({
        color: 0xff5c7a, emissive: 0xff5c7a, emissiveIntensity: 0.4,
      })
      this.racks.forEach(rack => {
        if (rack.roomId !== 'r4-1') return
        const p = rackWorld(rack.col, rack.row)
        const cab = new THREE.Mesh(new THREE.BoxGeometry(0.62, 2.05, 0.92), frame)
        cab.position.set(p.x, 1.02, p.z)
        cab.userData.pick = 'rack'
        cab.userData.rackId = rack.id
        group.add(cab)
        const face = new THREE.Mesh(new THREE.BoxGeometry(0.5, 1.7, 0.04), door)
        face.position.set(p.x, 1.05, p.z + 0.48)
        face.userData.pick = 'rack'
        face.userData.rackId = rack.id
        group.add(face)
        const led = new THREE.Mesh(
          new THREE.BoxGeometry(0.28, 0.05, 0.05),
          rack.level === 'alarm' ? this._alarmMat : ok
        )
        led.position.set(p.x, 1.9, p.z + 0.5)
        group.add(led)
      })
      group.visible = false
      return group
    },
    buildAssetRack() {
      const group = new THREE.Group()
      const cab = new THREE.Mesh(
        new THREE.BoxGeometry(1.15, 2.4, 0.7),
        new THREE.MeshStandardMaterial({ color: 0x9aafc0, metalness: 0.7, roughness: 0.28 })
      )
      cab.position.y = 1.2
      group.add(cab)
      group.visible = false
      return group
    },
    applyVisibility() {
      if (!this._building) return
      const s = this.state
      const spatial = s === 'cockpit' || s === 'building'
      this._building.visible = spatial
      this._floor.visible = s === 'floor' || s === 'power' || s === 'battery'
      this._room.visible = s === 'room' || s === 'follow'
      this._asset.visible = s === 'asset'
    },
    setPose(pose) {
      if (!this._camera) return
      this._camera.position.set(pose.pos[0], pose.pos[1], pose.pos[2])
      this._controls.target.set(pose.target[0], pose.target[1], pose.target[2])
      this._controls.update()
      this._fly = null
    },
    flyTo(pose, ms) {
      if (!this._camera) return
      this._fly = {
        from: this._camera.position.clone(),
        to: new THREE.Vector3(pose.pos[0], pose.pos[1], pose.pos[2]),
        fromT: this._controls.target.clone(),
        toT: new THREE.Vector3(pose.target[0], pose.target[1], pose.target[2]),
        t0: performance.now(),
        ms: ms || 1100,
      }
      this._controls.enabled = false
    },
    tick(now) {
      this._anim = requestAnimationFrame(this.tick)
      if (this._fly) {
        const t = Math.min(1, (now - this._fly.t0) / this._fly.ms)
        const k = ease(t)
        this._camera.position.lerpVectors(this._fly.from, this._fly.to, k)
        this._controls.target.lerpVectors(this._fly.fromT, this._fly.toT, k)
        if (t >= 1) this._fly = null
      }
      if (this._controls && this._controls.enabled) this._controls.update()
      if (this._renderer && this._scene && this._camera) {
        this._renderer.render(this._scene, this._camera)
      }
    },
    onPointer(ev) {
      if (!this._camera || this._fly) return
      const rect = this.$refs.canvas.getBoundingClientRect()
      const x = ((ev.clientX - rect.left) / rect.width) * 2 - 1
      const y = -((ev.clientY - rect.top) / rect.height) * 2 + 1
      const ray = new THREE.Raycaster()
      ray.setFromCamera({ x, y }, this._camera)
      const hits = ray.intersectObjects(this._scene.children, true)
      const hit = hits.find(h => h.object.userData && h.object.userData.pick)
      if (!hit) return
      const data = hit.object.userData
      if (data.pick === 'building') this.$emit('pick-building')
      if (data.pick === 'rack') this.$emit('pick-rack', data.rackId)
    },
    onResize() {
      const root = this.$refs.root
      if (!root || !this._renderer || !this._camera) return
      const w = root.clientWidth || 1
      const h = root.clientHeight || 1
      this._camera.aspect = w / h
      this._camera.updateProjectionMatrix()
      this._renderer.setSize(w, h, false)
    },
  },
}
</script>

<style scoped>
.demo-stage {
  position: absolute;
  inset: 0;
}
.demo-stage canvas {
  width: 100%;
  height: 100%;
  display: block;
}
</style>
