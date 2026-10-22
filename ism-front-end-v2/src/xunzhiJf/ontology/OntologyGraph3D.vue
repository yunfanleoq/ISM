<template>
  <div ref="root" class="ont-star">
    <canvas ref="canvas" />
    <div v-if="!nodes.length" class="ont-star-empty">这一层还没有节点。先编译，或换一层看。</div>
  </div>
</template>

<script>
import * as THREE from 'three'
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls'

const EDGE_COLOR = {
  physical: 0x3ec6e0,
  business: 0xf0a202,
  system: 0x3ddc97,
  bridge: 0xd2b4ff,
  gap: 0xff6b7a,
}

function ringPlace(nodes) {
  const pos = {}
  const focus = nodes.find(n => n.focus)
  const rest = focus ? nodes.filter(n => n.id !== focus.id) : nodes
  if (focus) pos[focus.id] = new THREE.Vector3(0, 0, 0)
  const radius = Math.max(6, Math.min(14, 4.2 + rest.length * 0.38))
  rest.forEach((n, i) => {
    const turn = (i / Math.max(rest.length, 1)) * Math.PI * 2 - Math.PI / 2
    pos[n.id] = new THREE.Vector3(Math.cos(turn) * radius, (i % 2 === 0 ? 0.7 : -0.7), Math.sin(turn) * radius)
  })
  return pos
}

function layerPlace(nodes) {
  const bands = {
    physical: { x: -16, y: 0 },
    business: { x: 0, y: 0 },
    system: { x: 16, y: 0 },
    bridge: { x: 0, y: 8 },
    gap: { x: 0, y: -8 },
  }
  const groups = {}
  nodes.forEach(n => {
    const layer = n.layer || 'system'
    if (!groups[layer]) groups[layer] = []
    groups[layer].push(n)
  })
  const pos = {}
  Object.keys(groups).forEach(layer => {
    const list = groups[layer]
    const band = bands[layer] || bands.system
    const cores = list.filter(n => n.kind === 'layer' || n.core)
    const rest = list.filter(n => !n.core && n.kind !== 'layer')
    cores.forEach(n => { pos[n.id] = new THREE.Vector3(band.x, band.y, 0) })
    const radius = Math.max(4.2, Math.min(11, 2.6 + rest.length * 0.2))
    rest.forEach((n, i) => {
      const turn = (i / Math.max(rest.length, 1)) * Math.PI * 2 - Math.PI / 2
      const lift = n.kind === 'relation' ? 1.4 : 0
      pos[n.id] = new THREE.Vector3(
        band.x + Math.cos(turn) * radius,
        band.y + lift + (i % 3 - 1) * 0.35,
        Math.sin(turn) * radius * 0.72,
      )
    })
  })
  return pos
}

const ORG_LEAF = 2.8
const ORG_DY = 6.2

function orgWidth(id, children, guard) {
  if (guard[id]) return 1
  const next = Object.assign({ [id]: true }, guard)
  const kids = children[id] || []
  if (!kids.length) return 1
  return kids.reduce((sum, kid) => sum + orgWidth(kid, children, next), 0)
}

function layoutForest(roots, children, pos) {
  function place(id, depth, left, guard) {
    if (guard[id]) return
    const next = Object.assign({ [id]: true }, guard)
    const kids = (children[id] || []).filter(kid => !next[kid])
    if (!kids.length) {
      pos[id] = { x: (left + 0.5) * ORG_LEAF, y: -depth * ORG_DY, z: 0 }
      return
    }
    let cursor = left
    kids.forEach(kid => {
      const width = orgWidth(kid, children, next)
      place(kid, depth + 1, cursor, next)
      cursor += width
    })
    const first = pos[kids[0]]
    const last = pos[kids[kids.length - 1]]
    pos[id] = {
      x: first && last ? (first.x + last.x) / 2 : (left + 0.5) * ORG_LEAF,
      y: -depth * ORG_DY,
      z: 0,
    }
  }
  let cursor = 0
  roots.forEach(id => {
    const width = orgWidth(id, children, {})
    place(id, 0, cursor, {})
    cursor += width
  })
  return cursor * ORG_LEAF
}

function orgPlace(nodes, edges) {
  const pos = {}
  const forests = []
  ;['physical', 'business', 'system', 'bridge', 'gap'].forEach(layer => {
    const list = nodes.filter(n => (n.layer || 'system') === layer)
    if (!list.length) return
    const ids = {}
    list.forEach(n => { ids[n.id] = true })
    const children = {}
    const hasParent = {}
    ;(edges || []).forEach(e => {
      if (!e.org || !ids[e.from] || !ids[e.to]) return
      if (!children[e.from]) children[e.from] = []
      if (children[e.from].indexOf(e.to) < 0) children[e.from].push(e.to)
      hasParent[e.to] = true
    })
    const roots = list.filter(n => !hasParent[n.id]).map(n => n.id)
    const local = {}
    const width = layoutForest(roots, children, local)
    forests.push({ local, width })
  })
  let shift = 0
  forests.forEach(forest => {
    Object.keys(forest.local).forEach(id => {
      const p = forest.local[id]
      pos[id] = new THREE.Vector3(shift + p.x, p.y, p.z || 0)
    })
    shift += forest.width + 10
  })
  let minX = Infinity
  let maxX = -Infinity
  let minY = Infinity
  let maxY = -Infinity
  Object.keys(pos).forEach(id => {
    const p = pos[id]
    minX = Math.min(minX, p.x)
    maxX = Math.max(maxX, p.x)
    minY = Math.min(minY, p.y)
    maxY = Math.max(maxY, p.y)
  })
  if (minX < Infinity) {
    const cx = (minX + maxX) / 2
    const cy = (minY + maxY) / 2
    Object.keys(pos).forEach(id => {
      pos[id].x -= cx
      pos[id].y -= cy
    })
    minX -= cx
    maxX -= cx
    minY -= cy
    maxY -= cy
  }
  return { pos, bounds: { minX, maxX, minY, maxY } }
}

function labelSprite(text, color, compact) {
  const canvas = document.createElement('canvas')
  canvas.width = 320
  canvas.height = 72
  const ctx = canvas.getContext('2d')
  ctx.clearRect(0, 0, 320, 72)
  ctx.fillStyle = 'rgba(5, 12, 24, 0.72)'
  ctx.fillRect(8, 14, 304, 44)
  ctx.font = '26px sans-serif'
  ctx.fillStyle = color || '#d6f5ff'
  ctx.textAlign = 'center'
  ctx.textBaseline = 'middle'
  const shown = text.length > 12 ? text.slice(0, 12) + '…' : text
  ctx.fillText(shown, 160, 36)
  const tex = new THREE.CanvasTexture(canvas)
  tex.needsUpdate = true
  const mat = new THREE.SpriteMaterial({ map: tex, transparent: true, depthWrite: false })
  const sprite = new THREE.Sprite(mat)
  sprite.scale.set(compact ? 3.1 : 4.4, compact ? 0.7 : 1, 1)
  return sprite
}

function arcBetween(a, b, lift) {
  const mid = new THREE.Vector3().addVectors(a, b).multiplyScalar(0.5)
  const dist = a.distanceTo(b)
  mid.y += lift + Math.min(3.2, dist * 0.08)
  const curve = new THREE.QuadraticBezierCurve3(a, mid, b)
  return new THREE.BufferGeometry().setFromPoints(curve.getPoints(22))
}

export default {
  name: 'OntologyGraph3D',
  props: {
    nodes: { type: Array, default: () => [] },
    edges: { type: Array, default: () => [] },
    layout: { type: String, default: 'ring' },
    focusId: { type: String, default: '' },
  },
  watch: {
    nodes: { deep: true, handler() { this.rebuild() } },
    edges: { deep: true, handler() { this.rebuild() } },
    layout() { this.rebuild() },
    focusId() { this.rebuild() },
  },
  mounted() {
    this.init()
    this.rebuild()
    window.addEventListener('resize', this.onResize)
    if (window.ResizeObserver) {
      this._ro = new ResizeObserver(() => this.onResize())
      this._ro.observe(this.$refs.root)
    }
  },
  beforeDestroy() {
    window.removeEventListener('resize', this.onResize)
    if (this._ro) this._ro.disconnect()
    this.disposeScene()
  },
  methods: {
    init() {
      const root = this.$refs.root
      const canvas = this.$refs.canvas
      const w = root.clientWidth || 800
      const h = root.clientHeight || 560
      this._renderer = new THREE.WebGLRenderer({ canvas, antialias: true })
      this._renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, 2))
      this._renderer.setSize(w, h, false)
      this._renderer.setClearColor(0x07111f, 1)
      this._scene = new THREE.Scene()
      this._scene.fog = new THREE.Fog(0x07111f, 42, 96)
      this._camera = new THREE.PerspectiveCamera(48, w / Math.max(h, 1), 0.1, 280)
      this._camera.position.set(0, 12, 42)
      this._ortho = new THREE.OrthographicCamera(-20, 20, 12, -12, 0.1, 500)
      this._ortho.position.set(0, 0, 80)
      this._controls = new OrbitControls(this._camera, canvas)
      this._controls.enableDamping = true
      this._controls.target.set(0, 0, 0)
      this._scene.add(new THREE.AmbientLight(0xd6e8ff, 0.85))
      const key = new THREE.DirectionalLight(0xffffff, 1.1)
      key.position.set(12, 20, 16)
      this._scene.add(key)
      const rim = new THREE.PointLight(0x3ec6e0, 0.7, 90)
      rim.position.set(-16, 8, -12)
      this._scene.add(rim)
      this._stars = this.makeStars()
      this._scene.add(this._stars)
      this._group = new THREE.Group()
      this._scene.add(this._group)
      this._ray = new THREE.Raycaster()
      this._pointer = new THREE.Vector2()
      this._flows = []
      canvas.addEventListener('pointerdown', this.onPointerDown)
      canvas.addEventListener('pointerup', this.onPointerUp)
      this.loop()
    },
    makeStars() {
      const arr = new Float32Array(160 * 3)
      for (let i = 0; i < 160; i += 1) {
        const r = 30 + Math.random() * 40
        const t = Math.random() * Math.PI * 2
        const p = Math.acos(2 * Math.random() - 1)
        arr[i * 3] = r * Math.sin(p) * Math.cos(t)
        arr[i * 3 + 1] = r * Math.cos(p) * 0.55
        arr[i * 3 + 2] = r * Math.sin(p) * Math.sin(t)
      }
      const geo = new THREE.BufferGeometry()
      geo.setAttribute('position', new THREE.BufferAttribute(arr, 3))
      return new THREE.Points(geo, new THREE.PointsMaterial({
        color: 0xd5ecff,
        size: 0.14,
        transparent: true,
        opacity: 0.75,
      }))
    },
    disposeScene() {
      if (this._anim) cancelAnimationFrame(this._anim)
      const canvas = this.$refs.canvas
      if (canvas) {
        canvas.removeEventListener('pointerdown', this.onPointerDown)
        canvas.removeEventListener('pointerup', this.onPointerUp)
      }
      if (this._controls) this._controls.dispose()
      this.clearGroup()
      if (this._stars) {
        this._stars.geometry.dispose()
        this._stars.material.dispose()
      }
      if (this._renderer) this._renderer.dispose()
    },
    clearGroup() {
      if (!this._group) return
      const old = this._group.children.slice()
      old.forEach(obj => {
        this._group.remove(obj)
        if (obj.geometry) obj.geometry.dispose()
        const mats = Array.isArray(obj.material) ? obj.material : [obj.material]
        mats.forEach(mat => {
          if (!mat) return
          if (mat.map) mat.map.dispose()
          mat.dispose()
        })
      })
      this._meshes = []
      this._flows = []
    },
    rebuild() {
      if (!this._group) return
      this.clearGroup()
      const nodes = this.nodes || []
      if (!nodes.length) return
      let pos
      this._bounds = null
      if (this.layout === 'org') {
        const placed = orgPlace(nodes, this.edges)
        pos = placed.pos
        this._bounds = placed.bounds
        if (this._scene && this._scene.fog) this._scene.fog.far = 240
      } else {
        pos = this.layout === 'layers' ? layerPlace(nodes) : ringPlace(nodes)
        if (this._scene && this._scene.fog) this._scene.fog.far = 96
      }
      const byId = {}
      const compact = nodes.length > 36
      const linked = {}
      if (this.focusId) {
        ;(this.edges || []).forEach(e => {
          if (e.from === this.focusId || e.to === this.focusId) {
            linked[e.from] = true
            linked[e.to] = true
          }
        })
        linked[this.focusId] = true
      }
      nodes.forEach(n => {
        const p = pos[n.id]
        if (!p) return
        byId[n.id] = p
        const dim = this.focusId && !linked[n.id]
        const mesh = new THREE.Mesh(
          new THREE.SphereGeometry(n.kind === 'layer' || n.core ? 1.28 : (n.focus || n.id === this.focusId ? 1.12 : (n.kind === 'relation' ? 0.5 : 0.72)), 24, 18),
          new THREE.MeshStandardMaterial({
            color: n.color || '#8aa0b8',
            emissive: n.color || '#8aa0b8',
            emissiveIntensity: dim ? 0.08 : (n.focus || n.kind === 'layer' ? 0.72 : 0.42),
            roughness: 0.28,
            metalness: 0.18,
            transparent: !!dim,
            opacity: dim ? 0.28 : 1,
          }),
        )
        mesh.position.copy(p)
        mesh.userData.node = n
        this._group.add(mesh)
        this._meshes.push(mesh)
        const showLabel = !dim && (nodes.length <= 110 || n.kind === 'layer' || n.kind === 'relation' || n.focus || n.id === this.focusId)
        if (showLabel) {
          const tag = labelSprite(n.label || n.id, '#e7f6ff', compact && n.kind === 'entity')
          tag.position.copy(p)
          if (this.layout === 'org' && n.kind !== 'layer') {
            tag.position.x += 1.45
            tag.position.y += 0.05
            tag.scale.set(2.5, 0.56, 1)
          } else {
            tag.position.y += n.kind === 'layer' ? 1.9 : 1.25
          }
          this._group.add(tag)
        }
      })
      ;(this.edges || []).forEach(e => {
        const a = byId[e.from]
        const b = byId[e.to]
        if (!a || !b) return
        const hot = !this.focusId || e.from === this.focusId || e.to === this.focusId
        const color = EDGE_COLOR[e.layer] || 0x8ec8ff
        const lift = e.org ? 0.12 : (e.layer === 'bridge' ? 3.6 : (e.faint ? 0.35 : 1.15))
        const geo = arcBetween(a, b, lift)
        const glow = new THREE.Line(geo, new THREE.LineBasicMaterial({
          color,
          transparent: true,
          opacity: hot ? (e.org ? 0.9 : (e.faint ? 0.22 : 0.55)) : 0.05,
          blending: THREE.AdditiveBlending,
          depthWrite: false,
        }))
        this._group.add(glow)
        if (!e.faint) {
          const dashGeo = geo.clone()
          const dash = new THREE.Line(dashGeo, new THREE.LineDashedMaterial({
            color,
            transparent: true,
            opacity: hot ? 0.95 : 0.12,
            dashSize: e.org ? 0.28 : 0.45,
            gapSize: e.org ? 0.16 : 0.32,
            blending: THREE.AdditiveBlending,
            depthWrite: false,
          }))
          dash.computeLineDistances()
          this._group.add(dash)
          this._flows.push(dash.material)
        }
        if (e.org && e.label) {
          const mid = new THREE.Vector3().addVectors(a, b).multiplyScalar(0.5)
          mid.y += 0.32
          const rel = labelSprite(e.label, '#d7f3ff', true)
          rel.position.copy(mid)
          rel.scale.set(1.9, 0.42, 1)
          this._group.add(rel)
        }
      })
      this.frame()
    },
    activeCamera() {
      return this.layout === 'org' && this._ortho ? this._ortho : this._camera
    },
    frame() {
      if (!this._camera || !this._controls) return
      if (this.layout === 'org' && this._ortho && this._bounds && isFinite(this._bounds.minX)) {
        const b = this._bounds
        const aspect = this._camera.aspect || 1.6
        const treeW = (b.maxX - b.minX) / 2 + 6
        const treeH = (b.maxY - b.minY) / 2 + 4
        const halfW = Math.max(treeW, treeH * aspect)
        const halfH = halfW / aspect
        this._ortho.left = -halfW
        this._ortho.right = halfW
        this._ortho.top = halfH
        this._ortho.bottom = -halfH
        this._ortho.position.set(0, 0, 80)
        this._ortho.lookAt(0, 0, 0)
        this._ortho.updateProjectionMatrix()
        this._controls.object = this._ortho
        this._controls.enableRotate = false
        this._controls.screenSpacePanning = true
        this._controls.mouseButtons.LEFT = THREE.MOUSE.PAN
        this._controls.target.set(0, 0, 0)
        this._controls.update()
        return
      }
      this._controls.object = this._camera
      this._controls.enableRotate = true
      this._controls.mouseButtons.LEFT = THREE.MOUSE.ROTATE
      const count = (this.nodes || []).length
      const wide = this.layout === 'layers'
      const dist = wide ? Math.min(58, 28 + count * 0.08) : (count <= 4 ? 16 : Math.min(34, 14 + count * 0.28))
      this._camera.position.set(0, dist * (wide ? 0.38 : 0.42), dist)
      this._controls.target.set(0, wide ? 1 : 0, 0)
      this._controls.update()
    },
    onPointerDown(ev) {
      this._down = { x: ev.clientX, y: ev.clientY }
    },
    onPointerUp(ev) {
      if (!this._down || !this._meshes || !this._meshes.length) return
      const dx = ev.clientX - this._down.x
      const dy = ev.clientY - this._down.y
      this._down = null
      if (dx * dx + dy * dy > 16) return
      const rect = this.$refs.canvas.getBoundingClientRect()
      this._pointer.x = ((ev.clientX - rect.left) / rect.width) * 2 - 1
      this._pointer.y = -((ev.clientY - rect.top) / rect.height) * 2 + 1
      this._ray.setFromCamera(this._pointer, this.activeCamera())
      const hit = this._ray.intersectObjects(this._meshes, false)[0]
      if (hit && hit.object.userData.node) this.$emit('select', hit.object.userData.node)
    },
    onResize() {
      const root = this.$refs.root
      if (!root || !this._renderer) return
      const w = root.clientWidth || 800
      const h = root.clientHeight || 560
      this._camera.aspect = w / Math.max(h, 1)
      this._camera.updateProjectionMatrix()
      this._renderer.setSize(w, h, false)
      this.frame()
    },
    loop() {
      this._anim = requestAnimationFrame(this.loop)
      if (this._stars) this._stars.rotation.y += 0.00035
      if (this._flows) this._flows.forEach(mat => { mat.dashOffset -= 0.012 })
      if (this._controls) this._controls.update()
      if (this._renderer) this._renderer.render(this._scene, this.activeCamera())
    },
  },
}
</script>

<style scoped>
.ont-star {
  position: relative;
  width: 100%;
  height: 100%;
  min-height: 520px;
  background: radial-gradient(circle at 50% 42%, #14324f 0%, #07111f 62%);
  border-radius: 8px;
  overflow: hidden;
}
.ont-star canvas {
  width: 100%;
  height: 100%;
  display: block;
  cursor: pointer;
}
.ont-star-empty {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #9fb3c8;
  pointer-events: none;
}
</style>
