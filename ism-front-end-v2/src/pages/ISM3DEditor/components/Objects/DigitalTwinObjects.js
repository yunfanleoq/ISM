/**
 * DigitalTwinObjects - 高端数字孪生专用3D组件
 * 为 SceneTemplates 提供具有国际一流水准的视觉效果
 * 参照翠鸟智慧园区数字孪生大屏标准
 *
 * 使用方式：
 *   在 IndustrialObjects.js 的 createThreeMesh 里引入此文件，
 *   在 switch 里调用 createDTObject(type, color)
 */

import * as THREE from 'three'

// ===== 翠鸟主题色 =====
const CUI = {
  cyan:    '#13c2c2',
  cyanGlow:'#4ddada',
  darkBg:  '#080e1a',
  gridLine:'#13c2c2',
  windowOn:'#a0e8ff',
  windowOff:'#0a2840',
  panelBg: '#0a1628',
  panelBorder: '#13c2c2',
  scanCol:  '#13c2c288',
  flowCol:  '#00e5ff',
}

/**
 * 主入口：根据 type 创建高端数字孪生对象
 * @returns {THREE.Group|null}
 */
export function createDTObject(type, color) {
  const c = color || CUI.cyan
  switch (type) {

    // ============================================================
    //  玻璃幕墙发光建筑（数字孪生专用）
    //  参照翠鸟截图：深蓝夜景 + 建筑玻璃幕墙自发光
    // ============================================================
    case 'dtBuilding': return buildDTBuilding(c)
    case 'dtBuildingTall': return buildDTBuildingTall(c)
    case 'dtBuildingWide': return buildDTBuildingWide(c)
    case 'dtBuildingComplex': return buildDTBuildingComplex(c)

    // ============================================================
    //  地面效果
    // ============================================================
    case 'dtGroundGrid': return buildDTGroundGrid(c)
    case 'dtScanRing': return buildDTScanRing(c)
    case 'dtScanLine': return buildDTScanLine(c)

    // ============================================================
    //  数据可视化
    // ============================================================
    case 'dtDataPanel': return buildDTDataPanel(c)
    case 'dtFlowLine': return buildDTFlowLine(c)
    case 'dtParticleField': return buildDTParticleField(c)
    case 'dtHologram': return buildDTHologram(c)

    // ============================================================
    //  园区/城市元素
    // ============================================================
    case 'dtRoad': return buildDTRoad(c)
    case 'dtPark': return buildDTPark(c)
    case 'dtStreetLightDT': return buildDTStreetLight(c)
    case 'dtBaseStation': return buildDTBaseStation(c)

    default: return null
  }
}

// ================================================================
//  建筑系列
// ================================================================

/**
 * dtBuilding — 标准玻璃幕墙建筑（程序化窗户纹理版）
 * 特点：DataTexture 生成窗户纹理 + 物理玻璃材质 + 随机亮灯动画
 */
function buildDTBuilding(c) {
  const group = new THREE.Group()
  const W = 1.4, D = 1.4, H = 3.2
  const floorCount = 12
  const winPerFloor = 5

  // ===== 主墙体：物理玻璃材质 =====
  const bodyMat = new THREE.MeshPhysicalMaterial({
    color: new THREE.Color(0x061428),
    metalness: 0.95,
    roughness: 0.08,
    transparent: true,
    opacity: 0.78,
    envMapIntensity: 1.2,
    clearcoat: 0.4,
    clearcoatRoughness: 0.1,
  })
  const body = new THREE.Mesh(new THREE.BoxGeometry(W, H, D), bodyMat)
  body.position.y = H / 2
  group.add(body)

  // ===== 程序化窗户纹理（DataTexture）=====
  const texSize = 256
  const texData = new Uint8Array(texSize * texSize * 4)
  const winW = Math.floor(texSize / winPerFloor)
  const winH = Math.floor(texSize / floorCount)

  for (let f = 0; f < floorCount; f++) {
    for (let w = 0; w < winPerFloor; w++) {
      const on = Math.random() > 0.35
      const r = on ? 0xa0 : 0x05, g = on ? 0xe8 : 0x10, b = on ? 0xff : 0x20
      const x0 = w * winW, x1 = x0 + winW - 2
      const y0 = (floorCount - 1 - f) * winH, y1 = y0 + winH - 2
      for (let y = y0; y < y1; y++) {
        for (let x = x0; x < x1; x++) {
          const i = (y * texSize + x) * 4
          texData[i] = r; texData[i+1] = g; texData[i+2] = b; texData[i+3] = 0xff
        }
      }
    }
  }
  const winTexture = new THREE.DataTexture(texData, texSize, texSize, THREE.RGBAFormat)
  winTexture.needsUpdate = true
  winTexture.wrapS = THREE.RepeatWrapping
  winTexture.wrapT = THREE.RepeatWrapping

  // 窗户平面（贴在四面墙上）
  const winMat = new THREE.MeshBasicMaterial({
    map: winTexture,
    emissive: new THREE.Color(0xa0e8ff),
    emissiveIntensity: 0.7,
    transparent: true,
    opacity: 0.82,
    side: THREE.FrontSide,
  })
  const winPlane = new THREE.PlaneGeometry(W * 0.98, H * 0.94)
  for (const [z, rotY] of [[D/2 + 0.04, 0], [-D/2 - 0.04, Math.PI], [W/2 + 0.04, Math.PI/2], [-W/2 - 0.04, -Math.PI/2]]) {
    const plane = new THREE.Mesh(winPlane, winMat)
    plane.position.set(0, H/2, z)
    plane.rotation.y = rotY
    group.add(plane)
  }

  // ===== 窗框线条（金属质感）=====
  const frameMat = new THREE.MeshStandardMaterial({ color: 0x334455, metalness: 0.9, roughness: 0.2 })
  for (let f = 0; f <= floorCount; f++) {
    const y = (f / floorCount) * H
    const bar = new THREE.Mesh(new THREE.BoxGeometry(W * 1.01, 0.012, 0.015), frameMat)
    bar.position.set(0, y, D/2 + 0.02)
    group.add(bar)
    const bar2 = bar.clone()
    bar2.position.z = -D/2 - 0.02
    group.add(bar2)
  }
  for (let w = 0; w <= winPerFloor; w++) {
    const x = -W/2 + (w / winPerFloor) * W
    const bar = new THREE.Mesh(new THREE.BoxGeometry(0.012, H * 0.96, 0.015), frameMat)
    bar.position.set(x, H/2, D/2 + 0.02)
    group.add(bar)
    const bar2 = bar.clone()
    bar2.position.z = -D/2 - 0.02
    group.add(bar2)
  }

  // ===== 顶部发光环 + 天线 =====
  const topRing = new THREE.Mesh(
    new THREE.TorusGeometry(0.22, 0.028, 12, 48),
    new THREE.MeshBasicMaterial({ color: 0x13c2c2, emissive: 0x13c2c2, emissiveIntensity: 1.8 })
  )
  topRing.position.y = H + 0.06
  topRing.rotation.x = Math.PI / 2
  group.add(topRing)

  const antenna = new THREE.Mesh(
    new THREE.CylinderGeometry(0.012, 0.012, 0.6, 8),
    new THREE.MeshStandardMaterial({ color: 0x888888, metalness: 0.95 })
  )
  antenna.position.y = H + 0.3
  group.add(antenna)
  const antennaTop = new THREE.Mesh(
    new THREE.SphereGeometry(0.05, 16, 16),
    new THREE.MeshBasicMaterial({ color: 0xff3333, emissive: 0xff33333, emissiveIntensity: 2.2 })
  )
  antennaTop.position.y = H + 0.6
  group.add(antennaTop)

  // ===== 底部发光底座 =====
  const baseGlow = new THREE.Mesh(
    new THREE.BoxGeometry(W + 0.1, 0.06, D + 0.1),
    new THREE.MeshBasicMaterial({ color: 0x13c2c2, emissive: 0x13c2c2, emissiveIntensity: 0.7, transparent: true, opacity: 0.65 })
  )
  baseGlow.position.y = 0.03
  group.add(baseGlow)

  // ===== 动画 =====
  group.userData.isMeshGroup = true
  group.userData.type = 'dtBuilding'
  group.userData.winTexture = winTexture
  group.userData.update = function(elapsed) {
    // 窗户纹理动画：随机窗户闪烁
    const data = winTexture.image.data
    for (let i = 0; i < data.length; i += 4) {
      if (Math.random() < 0.0003) {
        const on = data[i+1] > 100  // 当前是亮灯？
        data[i]   = on ? 0x05 : 0xa0
        data[i+1] = on ? 0x10 : 0xe8
        data[i+2] = on ? 0x20 : 0xff
      }
    }
    winTexture.needsUpdate = true
    // 顶部发光环呼吸
    const s = 0.85 + 0.35 * Math.sin(elapsed * 1.8)
    topRing.material.emissiveIntensity = 1.8 * s
    topRing.scale.set(s, s, s)
    // 天线灯闪烁
    antennaTop.material.emissiveIntensity = 1.8 + 1.8 * Math.sin(elapsed * 4.5)
  }
  return group
}

/**
 * dtBuildingTall — 高层玻璃幕墙建筑（H=5.0，程序化窗户纹理）
 */
function buildDTBuildingTall(c) {
  const group = new THREE.Group()
  const W = 1.2, D = 1.2, H = 5.2
  const floorCount = 18
  const winPerFloor = 5

  // ===== 主墙体：物理玻璃材质 =====
  const bodyMat = new THREE.MeshPhysicalMaterial({
    color: new THREE.Color(0x040a1c),
    metalness: 0.95,
    roughness: 0.06,
    transparent: true,
    opacity: 0.76,
    envMapIntensity: 1.4,
    clearcoat: 0.5,
    clearcoatRoughness: 0.08,
    reflectivity: 0.9,
  })
  const body = new THREE.Mesh(new THREE.BoxGeometry(W, H, D), bodyMat)
  body.position.y = H / 2
  group.add(body)

  // ===== 程序化窗户纹理 =====
  const texSize = 256
  const texData = new Uint8Array(texSize * texSize * 4)
  const wWin = Math.floor(texSize / winPerFloor)
  const hWin = Math.floor(texSize / floorCount)

  for (let f = 0; f < floorCount; f++) {
    for (let w = 0; w < winPerFloor; w++) {
      const on = Math.random() > 0.3
      const r = on ? 0xa0 : 0x05, g = on ? 0xe8 : 0x10, b = on ? 0xff : 0x20
      const x0 = w * wWin, x1 = x0 + wWin - 2
      const y0 = (floorCount - 1 - f) * hWin, y1 = y0 + hWin - 2
      for (let y = y0; y < y1; y++) {
        for (let x = x0; x < x1; x++) {
          const i = (y * texSize + x) * 4
          texData[i] = r, texData[i+1] = g, texData[i+2] = b, texData[i+3] = 0xff
        }
      }
    }
  }
  const winTexture = new THREE.DataTexture(texData, texSize, texSize, THREE.RGBAFormat)
  winTexture.needsUpdate = true
  winTexture.wrapS = THREE.RepeatWrapping
  winTexture.wrapT = THREE.RepeatWrapping

  const winMat = new THREE.MeshBasicMaterial({
    map: winTexture,
    emissive: new THREE.Color(0xa0e8ff),
    emissiveIntensity: 0.8,
    transparent: true,
    opacity: 0.85,
    side: THREE.FrontSide,
  })
  const winPlane = new THREE.PlaneGeometry(W * 0.96, H * 0.96)
  for (const [z, rotY] of [[D/2 + 0.03, 0], [-D/2 - 0.03, Math.PI], [W/2 + 0.03, Math.PI/2], [-W/2 - 0.03, -Math.PI/2]]) {
    const plane = new THREE.Mesh(winPlane, winMat)
    plane.position.set(0, H/2, z)
    plane.rotation.y = rotY
    group.add(plane)
  }

  // ===== 窗框线条 =====
  const frameMat = new THREE.MeshStandardMaterial({ color: 0x223344, metalness: 0.95, roughness: 0.15 })
  for (let f = 0; f <= floorCount; f++) {
    const bar = new THREE.Mesh(new THREE.BoxGeometry(W * 1.01, 0.01, 0.012), frameMat)
    bar.position.set(0, (f / floorCount) * H, D/2 + 0.04)
    group.add(bar)
    const bar2 = bar.clone()
    bar2.position.z = -D/2 - 0.04
    group.add(bar2)
  }
  for (let w = 0; w <= winPerFloor; w++) {
    const bar = new THREE.Mesh(new THREE.BoxGeometry(0.01, H * 0.96, 0.012), frameMat)
    bar.position.set(-W/2 + (w / winPerFloor) * W, H/2, D/2 + 0.04)
    group.add(bar)
    const bar2 = bar.clone()
    bar2.position.z = -D/2 - 0.04
    group.add(bar2)
  }

  // ===== 顶部双层发光环 + 天线 =====
  for (let i = 0; i < 2; i++) {
    const ring = new THREE.Mesh(
      new THREE.TorusGeometry(0.28 + i * 0.12, 0.02, 12, 48),
      new THREE.MeshBasicMaterial({ color: 0x13c2c2, emissive: 0x13c2c2, emissiveIntensity: 2.0 })
    )
    ring.position.y = H + 0.06 + i * 0.15
    ring.rotation.x = Math.PI / 2
    group.add(ring)
  }
  const antenna = new THREE.Mesh(
    new THREE.CylinderGeometry(0.01, 0.01, 0.7, 8),
    new THREE.MeshStandardMaterial({ color: 0x888888, metalness: 0.95 })
  )
  antenna.position.y = H + 0.35
  group.add(antenna)
  const antennaTop = new THREE.Mesh(
    new THREE.SphereGeometry(0.035, 16, 16),
    new THREE.MeshBasicMaterial({ color: 0xff3333, emissive: 0xff3333, emissiveIntensity: 2.5 })
  )
  antennaTop.position.y = H + 0.7
  group.add(antennaTop)

  // ===== 底部发光底座 =====
  const baseGlow = new THREE.Mesh(
    new THREE.BoxGeometry(W + 0.12, 0.05, D + 0.12),
    new THREE.MeshBasicMaterial({ color: 0x13c2c2, emissive: 0x13c2c2, emissiveIntensity: 0.7, transparent: true, opacity: 0.6 })
  )
  baseGlow.position.y = 0.025
  group.add(baseGlow)

  // ===== 动画 =====
  group.userData.isMeshGroup = true
  group.userData.type = 'dtBuildingTall'
  group.userData.winTexture = winTexture
  group.userData.update = function(elapsed) {
    // 窗户纹理随机闪烁
    const data = winTexture.image.data
    for (let i = 0; i < data.length; i += 4) {
      if (Math.random() < 0.0002) {
        const on = data[i+1] > 100
        data[i]   = on ? 0x05 : 0xa0
        data[i+1] = on ? 0x10 : 0xe8
        data[i+2] = on ? 0x20 : 0xff
      }
    }
    winTexture.needsUpdate = true
    // 顶部发光环呼吸
    const s = 0.82 + 0.38 * Math.sin(elapsed * 1.8)
    group.children.forEach(child => {
      if (child.type === 'TorusGeometry' || (child.geometry && child.geometry.type === 'TorusGeometry')) {
        child.material.emissiveIntensity = 2.0 * s
        child.scale.set(s, s, s)
      }
    })
    // 天线灯闪烁
    antennaTop.material.emissiveIntensity = 1.8 + 1.8 * Math.sin(elapsed * 4.5)
  }
  return group
}

/**
 * dtBuildingWide — 扁平厂房/数据中心（高端版）
 * DataTexture 窗户纹理 + 屋顶设备细节 + 底部发光底座
 */
function buildDTBuildingWide(c) {
  const group = new THREE.Group()
  const W = 2.8, D = 1.4, H = 2.0
  const floorCount = 3  // 扁平，只有3层

  // ===== 主墙体：物理玻璃材质 =====
  const bodyMat = new THREE.MeshPhysicalMaterial({
    color: new THREE.Color(0x040a18),
    metalness: 0.92,
    roughness: 0.10,
    transparent: true,
    opacity: 0.80,
    envMapIntensity: 1.2,
    clearcoat: 0.45,
    clearcoatRoughness: 0.12,
    reflectivity: 0.85,
  })
  const body = new THREE.Mesh(new THREE.BoxGeometry(W, H, D), bodyMat)
  body.position.y = H / 2
  group.add(body)

  // ===== 程序化窗户纹理（宽扁比例）=====
  const texSize = 256
  const texData = new Uint8Array(texSize * texSize * 4)
  const winPerRow = 8
  const winPerCol = floorCount
  const winW = Math.floor(texSize / winPerRow)
  const winH = Math.floor(texSize / winPerCol)

  for (let row = 0; row < winPerCol; row++) {
    for (let col = 0; col < winPerRow; col++) {
      const on = Math.random() > 0.35
      const r = on ? 0xa0 : 0x05, g = on ? 0xe8 : 0x10, b = on ? 0xff : 0x20
      const x0 = col * winW, x1 = x0 + winW - 3
      const y0 = (winPerCol - 1 - row) * winH, y1 = y0 + winH - 3
      for (let y = y0; y < y1; y++) {
        for (let x = x0; x < x1; x++) {
          const i = (y * texSize + x) * 4
          texData[i] = r; texData[i+1] = g; texData[i+2] = b; texData[i+3] = 0xff
        }
      }
    }
  }
  const winTexture = new THREE.DataTexture(texData, texSize, texSize, THREE.RGBAFormat)
  winTexture.needsUpdate = true
  winTexture.wrapS = THREE.RepeatWrapping
  winTexture.wrapT = THREE.RepeatWrapping

  // 窗户平面（贴在四面墙上）
  const winMat = new THREE.MeshBasicMaterial({
    map: winTexture,
    emissive: new THREE.Color(0xa0e8ff),
    emissiveIntensity: 0.75,
    transparent: true,
    opacity: 0.80,
    side: THREE.FrontSide,
  })
  const winPlane = new THREE.PlaneGeometry(W * 0.96, H * 0.92)
  for (const [z, rotY] of [[D/2 + 0.04, 0], [-D/2 - 0.04, Math.PI], [W/2 + 0.04, Math.PI/2], [-W/2 - 0.04, -Math.PI/2]]) {
    const plane = new THREE.Mesh(winPlane, winMat)
    plane.position.set(0, H/2, z)
    plane.rotation.y = rotY
    group.add(plane)
  }

  // ===== 窗框线条 =====
  const frameMat = new THREE.MeshStandardMaterial({ color: 0x223344, metalness: 0.95, roughness: 0.18 })
  // 水平线条
  for (let f = 0; f <= floorCount; f++) {
    const bar = new THREE.Mesh(new THREE.BoxGeometry(W * 1.01, 0.010, 0.014), frameMat)
    bar.position.set(0, (f / floorCount) * H, D/2 + 0.05)
    group.add(bar)
    const bar2 = bar.clone()
    bar2.position.z = -D/2 - 0.05
    group.add(bar2)
  }
  // 垂直线条
  for (let w = 0; w <= winPerRow; w++) {
    const bar = new THREE.Mesh(new THREE.BoxGeometry(0.010, H * 0.94, 0.014), frameMat)
    bar.position.set(-W/2 + (w / winPerRow) * W, H/2, D/2 + 0.05)
    group.add(bar)
    const bar2 = bar.clone()
    bar2.position.z = -D/2 - 0.05
    group.add(bar2)
  }

  // ===== 水平发光装饰条 =====
  for (let i = 0; i < 5; i++) {
    const strip = new THREE.Mesh(
      new THREE.BoxGeometry(W + 0.06, 0.035, D + 0.06),
      new THREE.MeshBasicMaterial({ color: 0x13c2c2, emissive: 0x13c2c2, emissiveIntensity: 0.9, transparent: true, opacity: 0.7 })
    )
    strip.position.y = 0.08 + (i + 1) * (H / 6)
    group.add(strip)
  }

  // ===== 正面进门 + 标识灯 =====
  const doorLight = new THREE.Mesh(
    new THREE.BoxGeometry(0.6, 0.10, 0.04),
    new THREE.MeshBasicMaterial({ color: 0x13c2c2, emissive: 0x13c2c2, emissiveIntensity: 2.2, transparent: true, opacity: 0.85 })
  )
  doorLight.position.set(0, 0.05, D / 2 + 0.02)
  group.add(doorLight)

  // ===== 屋顶设备（多个，细节丰富）=====
  const roofMat = new THREE.MeshStandardMaterial({ color: 0x555555, metalness: 0.7, roughness: 0.3 })
  // 空调机组
  const hvac = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.25, 0.35), roofMat)
  hvac.position.set(-0.6, H + 0.125, 0)
  group.add(hvac)
  // 通风管道
  const pipe = new THREE.Mesh(new THREE.CylinderGeometry(0.04, 0.04, 0.8, 8), roofMat)
  pipe.rotation.z = Math.PI / 2
  pipe.position.set(0.0, H + 0.22, 0.15)
  group.add(pipe)
  // 天线
  const antenna = new THREE.Mesh(new THREE.CylinderGeometry(0.015, 0.015, 0.7, 8), new THREE.MeshStandardMaterial({ color: 0x888888, metalness: 0.95 }))
  antenna.position.set(0.8, H + 0.35, -0.2)
  group.add(antenna)
  const antennaTop = new THREE.Mesh(
    new THREE.SphereGeometry(0.05, 16, 16),
    new THREE.MeshBasicMaterial({ color: 0xff3333, emissive: 0xff3333, emissiveIntensity: 2.5 })
  )
  antennaTop.position.set(0.8, H + 0.7, -0.2)
  group.add(antennaTop)

  // ===== 底部发光底座 =====
  const baseGlow = new THREE.Mesh(
    new THREE.BoxGeometry(W + 0.15, 0.06, D + 0.15),
    new THREE.MeshBasicMaterial({ color: 0x13c2c2, emissive: 0x13c2c2, emissiveIntensity: 0.7, transparent: true, opacity: 0.65 })
  )
  baseGlow.position.y = 0.03
  group.add(baseGlow)

  // ===== 动画 =====
  group.userData.isMeshGroup = true
  group.userData.type = 'dtBuildingWide'
  group.userData.winTexture = winTexture
  group.userData.antennaTop = antennaTop
  group.userData.update = function(elapsed) {
    // 窗户纹理随机闪烁
    const data = winTexture.image.data
    for (let i = 0; i < data.length; i += 4) {
      if (Math.random() < 0.0002) {
        const on = data[i+1] > 100
        data[i]   = on ? 0x05 : 0xa0
        data[i+1] = on ? 0x10 : 0xe8
        data[i+2] = on ? 0x20 : 0xff
      }
    }
    winTexture.needsUpdate = true
    // 天线灯闪烁
    antennaTop.material.emissiveIntensity = 2.0 + 2.0 * Math.sin(elapsed * 5.0)
    // 发光装饰条呼吸
    group.children.forEach(child => {
      if (child.material && child.material.emissive && child.material.color.getHex() === 0x13c2c2) {
        child.material.emissiveIntensity = 0.7 + 0.4 * Math.sin(elapsed * 2.0 + child.position.y)
      }
    })
  }
  return group
}

/**
 * dtBuildingComplex — 复合建筑群（主楼+裙楼）
 */
function buildDTBuildingComplex(c) {
  const group = new THREE.Group()
  // 主楼
  const main = buildDTBuildingTall(c)
  main.position.set(0, 0, 0)
  group.add(main)
  // 裙楼
  const wing1 = buildDTBuildingWide(c)
  wing1.position.set(-1.8, 0, 0)
  wing1.scale.set(0.7, 0.7, 0.7)
  group.add(wing1)
  const wing2 = buildDTBuilding('#1890ff')
  wing2.position.set(1.5, 0, 0.3)
  wing2.scale.set(0.6, 0.6, 0.6)
  group.add(wing2)

  group.userData.isMeshGroup = true
  group.userData.type = 'dtBuildingComplex'
  return group
}

// ================================================================
//  地面效果
// ================================================================

/**
 * dtGroundGrid — 发光地面网格（高端版）
 * CanvasTexture 程序化网格 + 动态流动 + 多层发光 + 中心脉冲
 */
function buildDTGroundGrid(c) {
  const group = new THREE.Group()
  const SIZE = 32

  // ===== 地面平面（深蓝黑半透明）=====
  const groundMat = new THREE.MeshPhysicalMaterial({
    color: 0x020818,
    metalness: 0.6,
    roughness: 0.12,
    transparent: true,
    opacity: 0.72,
    clearcoat: 0.5,
    clearcoatRoughness: 0.08,
    side: THREE.DoubleSide,
  })
  const ground = new THREE.Mesh(new THREE.PlaneGeometry(SIZE * 2, SIZE * 2), groundMat)
  ground.rotation.x = -Math.PI / 2
  ground.position.y = -0.01
  group.add(ground)

  // ===== 程序化网格纹理 =====
  const makeGridTexture = (lineColor, bgColor, lineWidth, spacing) => {
    const res = 512
    const cvs = document.createElement('canvas')
    cvs.width = res; cvs.height = res
    const ctx = cvs.getContext('2d')
    ctx.fillStyle = bgColor
    ctx.fillRect(0, 0, res, res)
    ctx.strokeStyle = lineColor
    ctx.lineWidth = lineWidth
    for (let i = 0; i <= res; i += spacing) {
      ctx.beginPath(); ctx.moveTo(i, 0); ctx.lineTo(i, res); ctx.stroke()
      ctx.beginPath(); ctx.moveTo(0, i); ctx.lineTo(res, i); ctx.stroke()
    }
    return new THREE.CanvasTexture(cvs)
  }

  // 主网格纹理（细线）
  const gridTex = makeGridTexture('rgba(19,194,194,0.35)', 'rgba(0,0,0,0)', 1, 16)
  gridTex.wrapS = THREE.RepeatWrapping
  gridTex.wrapT = THREE.RepeatWrapping
  gridTex.repeat.set(SIZE / 4, SIZE / 4)

  // 副网格纹理（粗线，亮色）
  const gridTexBold = makeGridTexture('rgba(19,194,194,0.55)', 'rgba(0,0,0,0)', 2, 64)
  gridTexBold.wrapS = THREE.RepeatWrapping
  gridTexBold.wrapT = THREE.RepeatWrapping
  gridTexBold.repeat.set(SIZE / 8, SIZE / 8)

  // 网格平面（贴在地面上）
  const gridMat = new THREE.MeshBasicMaterial({
    map: gridTex,
    transparent: true,
    opacity: 0.38,
    blending: THREE.AdditiveBlending,
    side: THREE.DoubleSide,
    depthWrite: false,
  })
  const gridMesh = new THREE.Mesh(new THREE.PlaneGeometry(SIZE * 2, SIZE * 2), gridMat)
  gridMesh.rotation.x = -Math.PI / 2
  gridMesh.position.y = 0.005
  group.add(gridMesh)

  // 粗网格平面
  const gridBoldMat = new THREE.MeshBasicMaterial({
    map: gridTexBold,
    transparent: true,
    opacity: 0.25,
    blending: THREE.AdditiveBlending,
    side: THREE.DoubleSide,
    depthWrite: false,
  })
  const gridBoldMesh = new THREE.Mesh(new THREE.PlaneGeometry(SIZE * 2, SIZE * 2), gridBoldMat)
  gridBoldMesh.rotation.x = -Math.PI / 2
  gridBoldMesh.position.y = 0.006
  group.add(gridBoldMesh)

  // ===== 轴向粗线（X/Z 轴，亮青）=====
  const axisMat = new THREE.LineBasicMaterial({
    color: 0x13c2c2,
    emissive: new THREE.Color(0x13c2c2),
    emissiveIntensity: 1.8,
    transparent: true,
    opacity: 0.6,
    linewidth: 2,
  })
  const axisGeoX = new THREE.BufferGeometry().setFromPoints([
    new THREE.Vector3(-SIZE, 0.007, 0), new THREE.Vector3(SIZE, 0.007, 0)
  ])
  const axisGeoZ = new THREE.BufferGeometry().setFromPoints([
    new THREE.Vector3(0, 0.007, -SIZE), new THREE.Vector3(0, 0.007, SIZE)
  ])
  group.add(new THREE.Line(axisGeoX, axisMat))
  group.add(new THREE.Line(axisGeoZ, axisMat))

  // ===== 中心发光脉冲圆（多层）=====
  const centerRings = []
  for (let i = 0; i < 3; i++) {
    const r = 1.5 + i * 1.8
    const ring = new THREE.Mesh(
      new THREE.RingGeometry(r, r + 0.06 + i * 0.03, 64),
      new THREE.MeshBasicMaterial({
        color: 0x13c2c2,
        emissive: new THREE.Color(0x13c2c2),
        emissiveIntensity: 2.0 - i * 0.5,
        transparent: true,
        opacity: 0.35 - i * 0.08,
        side: THREE.DoubleSide,
        blending: THREE.AdditiveBlending,
        depthWrite: false,
      })
    )
    ring.rotation.x = -Math.PI / 2
    ring.position.y = 0.008 + i * 0.001
    group.add(ring)
    centerRings.push({ mesh: ring, baseRadius: r, idx: i })
  }

  // ===== 中心辉光球 =====
  const glowSphere = new THREE.Mesh(
    new THREE.SphereGeometry(0.6, 32, 32),
    new THREE.MeshBasicMaterial({
      color: 0x13c2c2,
      emissive: new THREE.Color(0x13c2c2),
      emissiveIntensity: 3.0,
      transparent: true,
      opacity: 0.22,
      side: THREE.DoubleSide,
      blending: THREE.AdditiveBlending,
      depthWrite: false,
    })
  )
  glowSphere.position.y = 0.6
  group.add(glowSphere)

  // ===== 动画 =====
  group.userData.isMeshGroup = true
  group.userData.type = 'dtGroundGrid'
  group.userData.gridTex = gridTex
  group.userData.gridTexBold = gridTexBold
  group.userData.centerRings = centerRings
  group.userData.glowSphere = glowSphere
  group.userData.spawnTime = performance.now() / 1000

  group.userData.update = function(elapsed) {
    // 网格纹理流动
    gridTex.offset.x = elapsed * 0.015
    gridTex.offset.y = elapsed * 0.012
    gridTexBold.offset.x = elapsed * -0.008
    gridTexBold.offset.y = elapsed * 0.01

    // 中心环脉冲扩展
    centerRings.forEach(({ mesh, baseRadius, idx }) => {
      const pulse = 1.0 + 0.12 * Math.sin(elapsed * (1.8 + idx * 0.3) + idx * 1.5)
      mesh.scale.set(pulse, pulse, pulse)
      mesh.material.opacity = (0.35 - idx * 0.08) * (0.7 + 0.3 * Math.sin(elapsed * 2.2 + idx))
    })

    // 辉光球呼吸
    const s = 0.9 + 0.25 * Math.sin(elapsed * 2.0)
    glowSphere.scale.set(s, s, s)
    glowSphere.material.emissiveIntensity = 2.5 + 1.0 * Math.sin(elapsed * 3.0)
    glowSphere.position.y = 0.6 + 0.08 * Math.sin(elapsed * 1.5)
  }

  return group
}

/**
 * dtScanRing — 地面扫描环（高端版）
 * 多层冲击波环 + 地面光晕 + 扩展淡出动画
 */
function buildDTScanRing(c) {
  const group = new THREE.Group()
  const RING_COUNT = 3
  const rings = []

  for (let i = 0; i < RING_COUNT; i++) {
    const baseRadius = 0.5 + i * 0.3
    const ring = new THREE.Mesh(
      new THREE.TorusGeometry(baseRadius, 0.025 + i * 0.008, 12, 64),
      new THREE.MeshBasicMaterial({
        color: 0x13c2c2,
        emissive: new THREE.Color(0x13c2c2),
        emissiveIntensity: 2.0 - i * 0.4,
        transparent: true,
        opacity: 0.85 - i * 0.15,
      })
    )
    ring.rotation.x = Math.PI / 2
    ring.position.y = 0.03
    ring.userData.baseRadius = baseRadius
    ring.userData.speed = 1.2 + i * 0.35
    ring.userData.delay = i * 0.6  // 错开触发
    group.add(ring)
    rings.push(ring)

    // 地面光晕（每个环对应一个）
    const glow = new THREE.Mesh(
      new THREE.CircleGeometry(baseRadius + 0.3, 48),
      new THREE.MeshBasicMaterial({
        color: 0x13c2c2,
        emissive: new THREE.Color(0x13c2c2),
        emissiveIntensity: 0.4 - i * 0.08,
        transparent: true,
        opacity: 0.12 - i * 0.02,
        side: THREE.DoubleSide,
        depthWrite: false,
      })
    )
    glow.rotation.x = -Math.PI / 2
    glow.position.y = 0.005
    glow.userData.baseRadius = baseRadius + 0.3
    group.add(glow)
    rings.push(glow)  // 存到同一个数组方便动画
  }

  // 中心闪光
  const centerFlash = new THREE.Mesh(
    new THREE.SphereGeometry(0.15, 16, 16),
    new THREE.MeshBasicMaterial({
      color: 0xffffff,
      emissive: new THREE.Color(0xffffff),
      emissiveIntensity: 3.0,
      transparent: true,
      opacity: 0.8,
    })
  )
  centerFlash.position.y = 0.05
  group.add(centerFlash)

  // 动画
  group.userData.isMeshGroup = true
  group.userData.type = 'dtScanRing'
  group.userData.rings = rings
  group.userData.centerFlash = centerFlash
  group.userData.spawnTime = performance.now() / 1000

  group.userData.update = function(elapsed) {
    const t = elapsed - group.userData.spawnTime
    // 冲击波环扩展
    rings.forEach((ring, idx) => {
      if (ring.geometry && ring.geometry.type === 'TorusGeometry') {
        const delay = ring.userData.delay
        if (t < delay) return
        const expand = 1.0 + (t - delay) * ring.userData.speed
        const fade = Math.max(0, 1.0 - (t - delay) * 0.35)
        ring.scale.set(expand, expand, expand)
        ring.material.opacity = (0.85 - idx * 0.15) * fade
        ring.material.emissiveIntensity = (2.0 - idx * 0.4) * fade
        if (fade <= 0) ring.visible = false
      }
      // 地面光晕跟随扩展
      if (ring.geometry && ring.geometry.type === 'CircleGeometry') {
        const delay = ring.userData.baseRadius * 0.4
        if (t < delay) return
        const expand = 1.0 + (t - delay) * 0.9
        const fade = Math.max(0, 1.0 - (t - delay) * 0.28)
        ring.scale.set(expand, expand, expand)
        ring.material.opacity = (0.12 - idx * 0.02) * fade
        if (fade <= 0) ring.visible = false
      }
    })

    // 中心闪光脉冲
    const flashPulse = 0.6 + 0.4 * Math.sin(elapsed * 8.0)
    centerFlash.material.emissiveIntensity = 2.5 + 2.5 * flashPulse
    centerFlash.scale.setScalar(0.9 + 0.2 * flashPulse)
    centerFlash.material.opacity = 0.5 + 0.3 * flashPulse
  }

  return group
}

/**
 * dtScanLine — 垂直扫描光束（高端版）
 * 多层长方体体积光 + 噪声扫描 + 光头辉光
 */
function buildDTScanLine(c) {
  const group = new THREE.Group()
  const BEAM_H = 10
  const BEAM_W = 0.15  // 核心宽度

  // ===== 第1层：粗光晕（外光）=====
  const glowMat = new THREE.MeshBasicMaterial({
    color: 0x064545,
    emissive: new THREE.Color(0x064545),
    emissiveIntensity: 0.3,
    transparent: true,
    opacity: 0.12,
    side: THREE.DoubleSide,
    depthWrite: false,
  })
  const glowBeam = new THREE.Mesh(
    new THREE.BoxGeometry(BEAM_W * 6, BEAM_H, BEAM_W * 6),
    glowMat
  )
  glowBeam.position.y = BEAM_H / 2
  group.add(glowBeam)

  // ===== 第2层：中光带 =====
  const midMat = new THREE.MeshBasicMaterial({
    color: 0x00bcd4,
    emissive: new THREE.Color(0x00bcd4),
    emissiveIntensity: 1.0,
    transparent: true,
    opacity: 0.25,
    side: THREE.DoubleSide,
    depthWrite: false,
  })
  const midBeam = new THREE.Mesh(
    new THREE.BoxGeometry(BEAM_W * 3, BEAM_H, BEAM_W * 3),
    midMat
  )
  midBeam.position.y = BEAM_H / 2
  group.add(midBeam)

  // ===== 第3层：核心光束 =====
  const coreMat = new THREE.MeshBasicMaterial({
    color: 0x00e5ff,
    emissive: new THREE.Color(0x80ffff),
    emissiveIntensity: 2.5,
    transparent: true,
    opacity: 0.65,
    side: THREE.DoubleSide,
    depthWrite: false,
  })
  const coreBeam = new THREE.Mesh(
    new THREE.BoxGeometry(BEAM_W, BEAM_H, BEAM_W),
    coreMat
  )
  coreBeam.position.y = BEAM_H / 2
  group.add(coreBeam)

  // ===== 光头辉光球 =====
  const headGlow = new THREE.Mesh(
    new THREE.SphereGeometry(0.18, 16, 16),
    new THREE.MeshBasicMaterial({
      color: 0x00e5ff,
      emissive: new THREE.Color(0x00e5ff),
      emissiveIntensity: 2.0,
      transparent: true,
      opacity: 0.35,
      depthWrite: false,
    })
  )
  group.add(headGlow)

  const headCore = new THREE.Mesh(
    new THREE.SphereGeometry(0.06, 12, 12),
    new THREE.MeshBasicMaterial({
      color: 0xffffff,
      emissive: new THREE.Color(0xffffff),
      emissiveIntensity: 3.5,
    })
  )
  group.add(headCore)

  // ===== 底部接触光晕 =====
  const baseGlow = new THREE.Mesh(
    new THREE.CircleGeometry(0.4, 32),
    new THREE.MeshBasicMaterial({
      color: 0x13c2c2,
      emissive: new THREE.Color(0x13c2c2),
      emissiveIntensity: 1.2,
      transparent: true,
      opacity: 0.25,
      side: THREE.DoubleSide,
      depthWrite: false,
    })
  )
  baseGlow.rotation.x = -Math.PI / 2
  baseGlow.position.y = 0.01
  group.add(baseGlow)

  // ===== 动画 =====
  group.userData.isMeshGroup = true
  group.userData.type = 'dtScanLine'
  group.userData.glowMat  = glowMat
  group.userData.midMat    = midMat
  group.userData.coreMat  = coreMat
  group.userData.headGlow = headGlow
  group.userData.headCore = headCore
  group.userData.baseGlow = baseGlow
  group.userData.spawnTime = performance.now() / 1000
  group.userData.update = function(elapsed) {
    const t = elapsed - group.userData.spawnTime
    // 噪声扫描（sin+cos 模拟噪声）
    const noise = Math.sin(t * 0.8) * 0.5 + Math.cos(t * 1.3) * 0.3
    const sweep = noise * 6.0  // 扫描范围 ±6
    group.position.z = sweep
    // 光束呼吸
    const breath = 0.5 + 0.35 * Math.sin(elapsed * 2.0)
    glowMat.opacity = 0.08 + 0.06 * breath
    midMat.opacity = 0.18 + 0.12 * breath
    coreMat.opacity = 0.55 + 0.25 * breath
    glowMat.emissiveIntensity = 0.2 + 0.15 * breath
    midMat.emissiveIntensity = 0.8 + 0.4 * breath
    coreMat.emissiveIntensity = 2.0 + 0.8 * Math.sin(elapsed * 3.0)
    // 光头跟随
    headGlow.position.y = BEAM_H * (0.3 + 0.15 * Math.sin(elapsed * 1.5))
    headCore.position.y = headGlow.position.y
    const headPulse = 1.0 + 0.3 * Math.sin(elapsed * 5.0)
    headGlow.scale.setScalar(headPulse)
    // 底部光晕脉冲
    const basePulse = 1.0 + 0.25 * Math.sin(elapsed * 2.5)
    baseGlow.scale.setScalar(basePulse)
    baseGlow.material.opacity = 0.2 + 0.12 * Math.sin(elapsed * 3.0)
  }

  return group
}

// ================================================================
//  数据可视化
// ================================================================

/**
 * dtDataPanel — 3D悬浮数据面板（高端版）
 * 带发光边框 + 动态数据线条 + 模拟KPI数字 + 扫描线动画
 */
function buildDTDataPanel(c) {
  const group = new THREE.Group()
  const W = 2.2, H = 3.0

  // ===== 面板底板（深色半透明玻璃感）=====
  const panelMat = new THREE.MeshPhysicalMaterial({
    color: 0x050c1a,
    metalness: 0.3,
    roughness: 0.15,
    transparent: true,
    opacity: 0.72,
    clearcoat: 0.6,
    clearcoatRoughness: 0.1,
    side: THREE.DoubleSide,
  })
  const panel = new THREE.Mesh(new THREE.PlaneGeometry(W, H), panelMat)
  group.add(panel)

  // ===== 发光边框（外）=====
  const borderMat = new THREE.LineBasicMaterial({
    color: 0x13c2c2,
    emissive: new THREE.Color(0x13c2c2),
    emissiveIntensity: 1.5,
    transparent: true,
    opacity: 0.92,
  })
  const borderGeo = new THREE.BufferGeometry().setFromPoints([
    new THREE.Vector3(-W/2, -H/2, 0.012),
    new THREE.Vector3( W/2, -H/2, 0.012),
    new THREE.Vector3( W/2,  H/2, 0.012),
    new THREE.Vector3(-W/2,  H/2, 0.012),
    new THREE.Vector3(-W/2, -H/2, 0.012),
  ])
  const border = new THREE.Line(borderGeo, borderMat)
  group.add(border)

  // ===== 内边框（淡）=====
  const innerBorderMat = new THREE.LineBasicMaterial({
    color: 0x13c2c2,
    transparent: true,
    opacity: 0.25,
  })
  const innerBorderGeo = new THREE.BufferGeometry().setFromPoints([
    new THREE.Vector3(-W/2 + 0.06, -H/2 + 0.06, 0.014),
    new THREE.Vector3( W/2 - 0.06, -H/2 + 0.06, 0.014),
    new THREE.Vector3( W/2 - 0.06,  H/2 - 0.06, 0.014),
    new THREE.Vector3(-W/2 + 0.06,  H/2 - 0.06, 0.014),
    new THREE.Vector3(-W/2 + 0.06, -H/2 + 0.06, 0.014),
  ])
  group.add(new THREE.Line(innerBorderGeo, innerBorderMat))

  // ===== 四角发光块 =====
  const cornerMat = new THREE.MeshBasicMaterial({ color: 0x13c2c2, emissive: 0x13c2c2, emissiveIntensity: 2.2 })
  const cornerSize = 0.1
  [[-1,-1],[1,-1],[1,1],[-1,1]].forEach(([sx,sy]) => {
    const c = new THREE.Mesh(new THREE.BoxGeometry(cornerSize, 0.025, 0.025), cornerMat)
    c.position.set(sx * (W/2 - 0.05), sy * (H/2 - 0.025), 0.016)
    group.add(c)
  })

  // ===== 标题栏 =====
  const titleBar = new THREE.Mesh(
    new THREE.PlaneGeometry(W * 0.82, 0.28),
    new THREE.MeshBasicMaterial({ color: 0x13c2c2, emissive: 0x13c2c2, emissiveIntensity: 0.35, transparent: true, opacity: 0.45, side: THREE.DoubleSide })
  )
  titleBar.position.set(0, H/2 - 0.32, 0.016)
  group.add(titleBar)

  // ===== 动态数据线条（模拟图表）=====
  const lineColors = [0x00e5ff, 0x36cfc9, 0xffa940, 0xff4d4f, 0x9254de]
  const dataLines = []
  for (let i = 0; i < 5; i++) {
    const pts = []
    const seg = 20
    for (let j = 0; j <= seg; j++) {
      const t = j / seg
      const y = (Math.sin(t * Math.PI * 2 + i) * 0.3 + (i - 2) * 0.32) * 0.6
      pts.push(new THREE.Vector3(-W/2 + 0.25 + t * (W * 0.7), y, 0.015))
    }
    const lg = new THREE.BufferGeometry().setFromPoints(pts)
    const lm = new THREE.LineBasicMaterial({
      color: lineColors[i],
      emissive: new THREE.Color(lineColors[i]),
      emissiveIntensity: 1.2,
      transparent: true,
      opacity: 0.82,
    })
    const line = new THREE.Line(lg, lm)
    group.add(line)
    dataLines.push({ line, lm, phase: i * 0.8 })
  }

  // ===== 模拟KPI数字（发光小方块）=====
  const kpiMats = []
  for (let i = 0; i < 6; i++) {
    const kw = 0.18 + Math.random() * 0.25
    const kh = 0.08 + Math.random() * 0.06
    const ky = -H/2 + 0.55 + i * 0.32
    const kmat = new THREE.MeshBasicMaterial({
      color: i < 3 ? 0x00e5ff : 0xffa940,
      emissive: new THREE.Color(i < 3 ? 0x00e5ff : 0xffa940),
      emissiveIntensity: 0.6,
      transparent: true,
      opacity: 0.55,
    })
    kpiMats.push(kmat)
    const kbar = new THREE.Mesh(new THREE.PlaneGeometry(kw, kh), kmat)
    kbar.position.set(-W/2 + 0.25 + kw/2, ky, 0.015)
    group.add(kbar)
  }

  // ===== 中心扫描线（水平来回）=====
  const scanLineMat = new THREE.LineBasicMaterial({
    color: 0x13c2c2,
    emissive: 0x13c2c2,
    emissiveIntensity: 2.0,
    transparent: true,
    opacity: 0.7,
  })
  const scanLineGeo = new THREE.BufferGeometry().setFromPoints([
    new THREE.Vector3(-W/2 + 0.12, 0, 0.016),
    new THREE.Vector3( W/2 - 0.12, 0, 0.016),
  ])
  const scanLine = new THREE.Line(scanLineGeo, scanLineMat)
  group.add(scanLine)

  // ===== 动画 =====
  group.userData.isMeshGroup = true
  group.userData.type = 'dtDataPanel'
  group.userData.dataLines = dataLines
  group.userData.kpiMats = kpiMats
  group.userData.scanLine = scanLine
  group.userData.update = function(elapsed) {
    // 边框呼吸
    borderMat.emissiveIntensity = 1.5 + 0.5 * Math.sin(elapsed * 2.0)
    // 数据线条相位动画
    dataLines.forEach(({ lm, phase }) => {
      lm.opacity = 0.6 + 0.3 * Math.sin(elapsed * 1.5 + phase)
      lm.emissiveIntensity = 1.0 + 0.5 * Math.sin(elapsed * 1.8 + phase)
    })
    // KPI块闪烁
    kpiMats.forEach((km, i) => {
      km.emissiveIntensity = 0.4 + 0.4 * Math.sin(elapsed * 2.2 + i * 0.7)
    })
    // 扫描线上下扫
    const sweep = Math.sin(elapsed * 1.2) * (H * 0.35)
    scanLine.position.y = sweep
    scanLineMat.opacity = 0.5 + 0.3 * Math.sin(elapsed * 3.0)
  }
  return group
}

/**
 * dtFlowLine — 流光数据线（高端版）
 * 多层Tube模拟体积发光 + 多个流动光点 + 粒子拖尾
 */
function buildDTFlowLine(c) {
  const group = new THREE.Group()
  const LEN = 8
  const SEG = 60

  // ===== 有机曲线（噪声扰动）=====
  const pts = []
  for (let i = 0; i <= SEG; i++) {
    const t = i / SEG
    const x = -LEN/2 + t * LEN
    const noise = Math.sin(t * Math.PI * 4) * 0.5 + Math.cos(t * Math.PI * 7) * 0.3
    const z = noise + Math.sin(t * Math.PI * 2.5) * 0.4
    const y = 0.04 + Math.sin(t * Math.PI * 3) * 0.06
    pts.push(new THREE.Vector3(x, y, z))
  }
  const curve = new THREE.CatmullRomCurve3(pts)
  const curvePoints = curve.getPoints(80)

  // ===== 第1层：粗发光晕（外光）=====
  const glowTubeGeo = new THREE.TubeGeometry(curve, 60, 0.07, 8, false)
  const glowMesh = new THREE.Mesh(glowTubeGeo, new THREE.MeshBasicMaterial({
    color: 0x001a3a,
    emissive: new THREE.Color(0x001a3a),
    emissiveIntensity: 0.6,
    transparent: true,
    opacity: 0.22,
    side: THREE.DoubleSide,
  }))
  group.add(glowMesh)

  // ===== 第2层：中光带 =====
  const midTubeGeo = new THREE.TubeGeometry(curve, 60, 0.035, 8, false)
  const midMesh = new THREE.Mesh(midTubeGeo, new THREE.MeshBasicMaterial({
    color: 0x00bcd4,
    emissive: new THREE.Color(0x00bcd4),
    emissiveIntensity: 1.5,
    transparent: true,
    opacity: 0.65,
    side: THREE.DoubleSide,
  }))
  group.add(midMesh)

  // ===== 第3层：亮芯线 =====
  const coreGeo = new THREE.BufferGeometry().setFromPoints(curvePoints)
  const coreMat = new THREE.LineBasicMaterial({
    color: 0x80ffff,
    emissive: new THREE.Color(0x80ffff),
    emissiveIntensity: 3.0,
    transparent: true,
    opacity: 0.95,
  })
  const coreLine = new THREE.Line(coreGeo, coreMat)
  group.add(coreLine)

  // ===== 流动光点（3个，不同相位）=====
  const heads = []
  for (let h = 0; h < 3; h++) {
    const glowSphere = new THREE.Mesh(
      new THREE.SphereGeometry(0.14, 16, 16),
      new THREE.MeshBasicMaterial({
        color: 0x00e5ff,
        emissive: new THREE.Color(0x00e5ff),
        emissiveIntensity: 2.0,
        transparent: true,
        opacity: 0.32,
      })
    )
    const coreSphere = new THREE.Mesh(
      new THREE.SphereGeometry(0.045, 12, 12),
      new THREE.MeshBasicMaterial({
        color: 0xffffff,
        emissive: new THREE.Color(0xffffff),
        emissiveIntensity: 3.5,
      })
    )
    const hg = new THREE.Group()
    hg.add(glowSphere)
    hg.add(coreSphere)
    hg.userData.phase = h / 3.0
    group.add(hg)
    heads.push(hg)
  }

  // ===== 尾部粒子拖尾 =====
  const trailCount = 15
  const trails = []
  for (let i = 0; i < trailCount; i++) {
    const tp = new THREE.Mesh(
      new THREE.SphereGeometry(0.012 + Math.random() * 0.018, 6, 6),
      new THREE.MeshBasicMaterial({
        color: 0x00e5ff,
        emissive: new THREE.Color(0x00e5ff),
        emissiveIntensity: 1.8,
        transparent: true,
        opacity: 0.5 + Math.random() * 0.3,
      })
    )
    tp.userData.idx = i
    tp.userData.phaseOffset = Math.random() * 0.3
    group.add(tp)
    trails.push(tp)
  }

  // ===== 动画 =====
  group.userData.isMeshGroup = true
  group.userData.type = 'dtFlowLine'
  group.userData.curve = curve
  group.userData.heads = heads
  group.userData.trails = trails
  group.userData.glowMesh = glowMesh
  group.userData.midMesh = midMesh
  group.userData.coreMat = coreMat

  group.userData.update = function(elapsed) {
    // 光带呼吸
    const breath = 0.6 + 0.35 * Math.sin(elapsed * 2.0)
    glowMesh.material.emissiveIntensity = 0.4 + 0.25 * breath
    midMesh.material.emissiveIntensity = 1.2 + 0.5 * breath
    coreMat.emissiveIntensity = 2.5 + 0.9 * Math.sin(elapsed * 3.0)

    // 流动光点沿曲线运动
    heads.forEach((hg) => {
      hg.userData.phase += 0.0035
      if (hg.userData.phase > 1.0) hg.userData.phase -= 1.0
      const pt = curve.getPoint(hg.userData.phase)
      hg.position.copy(pt)
      hg.position.y += 0.05
      // 大小脉冲
      const pulse = 1.0 + 0.35 * Math.sin(elapsed * 5.0 + hg.userData.phase * 10)
      hg.children[0].scale.setScalar(pulse)
    })

    // 尾部粒子拖尾
    trails.forEach((tp) => {
      const phase = ((elapsed * 0.28 + tp.userData.idx * 0.065 + tp.userData.phaseOffset) % 1.0)
      const pt = curve.getPoint(phase)
      tp.position.copy(pt)
      tp.position.y -= 0.06 + tp.userData.idx * 0.008
      tp.material.opacity = 0.6 - tp.userData.idx * 0.03
      tp.scale.setScalar(1.0 - tp.userData.idx * 0.04)
    })
  }

  return group
}

/**
 * dtParticleField — 氛围粒子场（高端版）
 * 圆形光晕纹理 + 噪声运动 + 颜色/大小脉动
 */
function buildDTParticleField(c) {
  const group = new THREE.Group()
  const COUNT = 280

  // ===== 生成圆形光晕纹理 =====
  const canvas = document.createElement('canvas')
  canvas.width = 64
  canvas.height = 64
  const ctx = canvas.getContext('2d')
  const gradient = ctx.createRadialGradient(32, 32, 0, 32, 32, 32)
  gradient.addColorStop(0, 'rgba(160, 232, 194, 1.0)')
  gradient.addColorStop(0.25, 'rgba(100, 220, 255, 0.6)')
  gradient.addColorStop(0.6, 'rgba(50, 180, 255, 0.15)')
  gradient.addColorStop(1, 'rgba(0, 0, 0, 0.0)')
  ctx.fillStyle = gradient
  ctx.fillRect(0, 0, 64, 64)
  const particleTexture = new THREE.CanvasTexture(canvas)

  // ===== 粒子位置 + 噪声种子 =====
  const positions = new Float32Array(COUNT * 3)
  const colors    = new Float32Array(COUNT * 3)
  const sizes     = new Float32Array(COUNT)    // 每个粒子的基础大小
  const phases    = new Float32Array(COUNT * 3) // noise 偏移种子

  const colorPalette = [
    new THREE.Color(0x13c2c2),  // cyan
    new THREE.Color(0x36cfc9),  // teal
    new THREE.Color(0x4d8af),   // blue
    new THREE.Color(0x9254de),  // purple
    new THREE.Color(0x00e5ff),  // bright cyan
  ]

  for (let i = 0; i < COUNT; i++) {
    positions[i*3]   = (Math.random() - 0.5) * 22
    positions[i*3+1] = Math.random() * 8
    positions[i*3+2] = (Math.random() - 0.5) * 22

    const ci = Math.floor(Math.random() * colorPalette.length)
    const c = colorPalette[ci]
    colors[i*3]   = c.r
    colors[i*3+1] = c.g
    colors[i*3+2] = c.b

    sizes[i] = 0.03 + Math.random() * 0.06

    phases[i*3]   = Math.random() * 1000  // noise x 种子
    phases[i*3+1] = Math.random() * 1000  // noise y 种子
    phases[i*3+2] = Math.random() * 1000  // noise z 种子
  }

  const geo = new THREE.BufferGeometry()
  geo.setAttribute('position', new THREE.BufferAttribute(positions, 3))
  geo.setAttribute('color',    new THREE.BufferAttribute(colors, 3))
  geo.setAttribute('size',      new THREE.BufferAttribute(sizes, 1))

  // ===== 自定义着色粒子材质 =====
  // 用 ShaderMaterial 让每个粒子根据 `size` 属性和时间缩放
  const mat = new THREE.ShaderMaterial({
    uniforms: {
      uTime:    { value: 0 },
      uTexture: { value: particleTexture },
    },
    vertexShader: `
      attribute float size;
      varying vec3 vColor;
      varying float vAlpha;
      uniform float uTime;
      void main() {
        vColor = color;
        // 脉动缩放
        float pulse = 1.0 + 0.35 * sin(uTime * 2.0 + position.x * 3.0 + position.z * 2.0);
        vAlpha = 0.45 + 0.35 * sin(uTime * 1.8 + position.y * 2.5);
        vec4 mvPosition = modelViewMatrix * vec4(position, 1.0);
        gl_PointSize = size * pulse * (300.0 / -mvPosition.z);
        gl_Position = projectionMatrix * mvPosition;
      }
    `,
    fragmentShader: `
      uniform sampler2D uTexture;
      varying vec3 vColor;
      varying float vAlpha;
      void main() {
        vec4 texColor = texture2D(uTexture, gl_PointCoord);
        gl_FragColor = vec4(vColor * 1.8, texColor.a * vAlpha);
        if (gl_FragColor.a < 0.01) discard;
      }
    `,
    transparent: true,
    blending: THREE.AdditiveBlending,
    depthWrite: false,
    vertexColors: true,
  })

  const points = new THREE.Points(geo, mat)
  group.add(points)

  // ===== 动画 =====
  group.userData.isMeshGroup = true
  group.userData.type = 'dtParticleField'
  group.userData.positions = positions
  group.userData.phases    = phases
  group.userData.mat        = mat
  group.userData.update = function(elapsed) {
    mat.uniforms.uTime.value = elapsed
    const pos = points.geometry.attributes.position.array
    for (let i = 0; i < COUNT; i++) {
      // 噪声运动（用 sin/cos 模拟 simplex noise）
      const nx = Math.sin(pos[i*3]   * 0.3 + elapsed * 0.15 + phases[i*3])
      const ny = Math.cos(pos[i*3+1] * 0.4 + elapsed * 0.2  + phases[i*3+1])
      const nz = Math.sin(pos[i*3+2] * 0.3 + elapsed * 0.12 + phases[i*3+2])
      pos[i*3]   += nx * 0.003
      pos[i*3+1] += 0.004 + ny * 0.002  // 缓慢上浮
      pos[i*3+2] += nz * 0.003
      // 边界重置
      if (pos[i*3+1] > 9)  pos[i*3+1] = 0
      if (pos[i*3]   < -12) pos[i*3]   = 12
      if (pos[i*3]   > 12)  pos[i*3]   = -12
      if (pos[i*3+2] < -12) pos[i*3+2] = 12
      if (pos[i*3+2] > 12)  pos[i*3+2] = -12
    }
    points.geometry.attributes.position.needsUpdate = true
  }
  return group
}

/**
 * dtHologram — 全息投影（高端版）
 * 多层锥体体积感 + 3条扫描线 + 全息图像 + 噪声抖动
 */
function buildDTHologram(c) {
  const group = new THREE.Group()
  const H_CION = 2.4  // 主锥体高度

  // ===== 多层锥体（大/中/小）=====
  const coneMaterials = [
    { color: 0x031a2e, opacity: 0.10, scale: 1.0  },
    { color: 0x0a3a5a, opacity: 0.18, scale: 0.78 },
    { color: 0x13c2c2, opacity: 0.30, scale: 0.52 },
  ]
  const cones = []
  coneMaterials.forEach(({ color, opacity, scale }, idx) => {
    const cone = new THREE.Mesh(
      new THREE.ConeGeometry(0.9 * scale, H_CION * scale, 48, 1, true),
      new THREE.MeshBasicMaterial({
        color: color,
        emissive: new THREE.Color(color),
        emissiveIntensity: idx === 2 ? 0.5 : 0.2,
        transparent: true,
        opacity: opacity,
        side: THREE.DoubleSide,
      })
    )
    cone.position.y = (H_CION * scale) / 2
    group.add(cone)
    cones.push(cone)
  })

  // ===== 扫描线（3条，不同大小/速度/相位）=====
  const scanLines = []
  for (let i = 0; i < 3; i++) {
    const radius = 0.15 + i * 0.28
    const ring = new THREE.Mesh(
      new THREE.RingGeometry(radius - 0.012, radius + 0.012, 48),
      new THREE.MeshBasicMaterial({
        color: 0x13c2c2,
        emissive: new THREE.Color(0x13c2c2),
        emissiveIntensity: 2.0 - i * 0.4,
        transparent: true,
        opacity: 0.75 - i * 0.15,
        side: THREE.DoubleSide,
      })
    )
    ring.rotation.x = -Math.PI / 2
    ring.position.y = 0.2 + i * 0.01  // 初始高度不同
    ring.userData.speed   = 0.6 + i * 0.25
    ring.userData.phase  = i * 2.1
    ring.userData.amplitude = 0.9 - i * 0.15
    group.add(ring)
    scanLines.push(ring)
  }

  // ===== 全息图像模拟（中心简笔人形）=====
  const holoCanvas = document.createElement('canvas')
  holoCanvas.width = 128
  holoCanvas.height = 200
  const hCtx = holoCanvas.getContext('2d')
  // 半透明背景
  hCtx.fillStyle = 'rgba(0, 0, 0, 0)'
  hCtx.fillRect(0, 0, 128, 200)
  // 简笔人形（青色发光线）
  hCtx.strokeStyle = 'rgba(19, 194, 194, 0.7)'
  hCtx.lineWidth = 2
  hCtx.beginPath()
  // 头
  hCtx.arc(64, 28, 12, 0, Math.PI * 2)
  hCtx.stroke()
  // 身体
  hCtx.moveTo(64, 40)
  hCtx.lineTo(64, 110)
  hCtx.stroke()
  // 手臂
  hCtx.moveTo(64, 65)
  hCtx.lineTo(36, 95)
  hCtx.stroke()
  hCtx.moveTo(64, 65)
  hCtx.lineTo(92, 95)
  hCtx.stroke()
  // 腿
  hCtx.moveTo(64, 110)
  hCtx.lineTo(48, 170)
  hCtx.stroke()
  hCtx.moveTo(64, 110)
  hCtx.lineTo(80, 170)
  hCtx.stroke()
  const holoTexture = new THREE.CanvasTexture(holoCanvas)
  const holoMat = new THREE.MeshBasicMaterial({
    map: holoTexture,
    transparent: true,
    opacity: 0.55,
    side: THREE.DoubleSide,
    depthWrite: false,
  })
  const holoPlane = new THREE.Mesh(new THREE.PlaneGeometry(0.7, 1.6), holoMat)
  holoPlane.position.y = 0.95
  holoPlane.position.z = -0.08
  group.add(holoPlane)

  // ===== 底座 =====
  // 扁圆柱
  const baseDisk = new THREE.Mesh(
    new THREE.CylinderGeometry(0.9, 0.95, 0.06, 48),
    new THREE.MeshStandardMaterial({
      color: 0x112233,
      metalness: 0.9,
      roughness: 0.15,
    })
  )
  baseDisk.position.y = 0.03
  group.add(baseDisk)
  // 发光环
  const baseRing = new THREE.Mesh(
    new THREE.TorusGeometry(0.92, 0.025, 12, 64),
    new THREE.MeshBasicMaterial({
      color: 0x13c2c2,
      emissive: new THREE.Color(0x13c2c2),
      emissiveIntensity: 2.0,
    })
  )
  baseRing.rotation.x = Math.PI / 2
  baseRing.position.y = 0.06
  group.add(baseRing)
  // 底部扩散光
  const baseGlow = new THREE.Mesh(
    new THREE.CircleGeometry(1.1, 48),
    new THREE.MeshBasicMaterial({
      color: 0x13c2c2,
      emissive: new THREE.Color(0x13c2c2),
      emissiveIntensity: 0.3,
      transparent: true,
      opacity: 0.12,
      side: THREE.DoubleSide,
    })
  )
  baseGlow.rotation.x = -Math.PI / 2
  baseGlow.position.y = 0.005
  group.add(baseGlow)

  // ===== 动画 =====
  group.userData.isMeshGroup = true
  group.userData.type = 'dtHologram'
  group.userData.cones      = cones
  group.userData.scanLines  = scanLines
  group.userData.holoPlane  = holoPlane
  group.userData.baseRing  = baseRing
  group.userData.baseGlow  = baseGlow

  group.userData.update = function(elapsed) {
    // 锥体呼吸
    cones.forEach((cone, idx) => {
      const s = 0.92 + 0.12 * Math.sin(elapsed * (1.0 + idx * 0.3) + idx * 1.5)
      cone.scale.setScalar(s)
      cone.material.opacity = cone.material.userData_baseOpacity * (0.8 + 0.2 * Math.sin(elapsed * 1.8 + idx))
    })

    // 扫描线：有机噪声运动
    scanLines.forEach((ring) => {
      const { speed, phase, amplitude } = ring.userData
      // 用 sin+cos 模拟噪声
      const raw = Math.sin(elapsed * speed + phase) * amplitude
      const noise = raw + Math.sin(elapsed * speed * 2.7 + phase * 1.3) * 0.12
      ring.position.y = 0.15 + (noise + 1.0) * (H_CION * 0.42)
      ring.material.opacity = (0.55 + 0.25 * Math.sin(elapsed * 3.0 + phase))
    })

    // 全息图像上下浮动 + 透明度脉动
    holoPlane.position.y = 0.95 + 0.04 * Math.sin(elapsed * 1.2)
    holoMat.opacity = 0.45 + 0.15 * Math.sin(elapsed * 2.0)

    // 底座发光环呼吸
    const breath = 1.8 + 0.6 * Math.sin(elapsed * 2.2)
    baseRing.material.emissiveIntensity = breath
    baseRing.scale.setScalar(0.97 + 0.06 * Math.sin(elapsed * 1.8))

    // 底部扩散光呼吸
    baseGlow.material.opacity = 0.10 + 0.06 * Math.sin(elapsed * 1.5)
    baseGlow.scale.setScalar(1.0 + 0.08 * Math.sin(elapsed * 1.3))
  }

  // 保存基础透明度供动画使用
  cones.forEach(c => { c.material.userData_baseOpacity = c.material.opacity })

  return group
}

// ================================================================
//  园区/城市元素
// ================================================================

/**
 * dtRoad — 道路片段（带车道线和路边发光）
 */
function buildDTRoad(c) {
  const group = new THREE.Group()
  const W = 6, L = 8

  // 路面
  const road = new THREE.Mesh(
    new THREE.PlaneGeometry(W, L),
    new THREE.MeshStandardMaterial({ color: 0x1a1a2e, roughness: 0.9 })
  )
  road.rotation.x = -Math.PI / 2
  group.add(road)

  // 车道线（虚线）
  const lineMat = new THREE.LineBasicMaterial({ color: 0xddddaa, transparent: true, opacity: 0.5 })
  for (let i = 0; i < L; i += 1.5) {
    const seg = new THREE.BufferGeometry().setFromPoints([
      new THREE.Vector3(0, 0.01, -L/2 + i),
      new THREE.Vector3(0, 0.01, -L/2 + i + 0.6),
    ])
    const line = new THREE.Line(seg, lineMat)
    group.add(line)
  }

  // 路边发光边线
  const edgeMat = new THREE.LineBasicMaterial({ color: 0x13c2c2, emissive: 0x13c2c2, emissiveIntensity: 0.8, transparent: true, opacity: 0.6 })
  const leftEdge = new THREE.BufferGeometry().setFromPoints([
    new THREE.Vector3(-W/2 + 0.1, 0.01, -L/2),
    new THREE.Vector3(-W/2 + 0.1, 0.01, L/2),
  ])
  const rightEdge = new THREE.BufferGeometry().setFromPoints([
    new THREE.Vector3(W/2 - 0.1, 0.01, -L/2),
    new THREE.Vector3(W/2 - 0.1, 0.01, L/2),
  ])
  group.add(new THREE.Line(leftEdge, edgeMat))
  group.add(new THREE.Line(rightEdge, edgeMat))

  group.userData.isMeshGroup = true
  group.userData.type = 'dtRoad'
  return group
}

/**
 * dtPark — 绿地+树（数字孪生风格，树顶部有发光点）
 */
function buildDTPark(c) {
  const group = new THREE.Group()

  // 草地
  const grass = new THREE.Mesh(
    new THREE.PlaneGeometry(2, 2),
    new THREE.MeshStandardMaterial({ color: 0x0a331a, roughness: 1.0 })
  )
  grass.rotation.x = -Math.PI / 2
  grass.position.y = 0.005
  group.add(grass)

  // 树（简化：圆柱+球）
  const trunk = new THREE.Mesh(
    new THREE.CylinderGeometry(0.06, 0.08, 0.6, 8),
    new THREE.MeshStandardMaterial({ color: 0x553311 })
  )
  trunk.position.y = 0.3
  group.add(trunk)

  const crown = new THREE.Mesh(
    new THREE.SphereGeometry(0.35, 16, 16),
    new THREE.MeshStandardMaterial({ color: 0x116633, emissive: 0x0a331a, emissiveIntensity: 0.2 })
  )
  crown.position.y = 0.9
  group.add(crown)

  // 树顶传感器/发光点（数字孪生感）
  const sensor = new THREE.Mesh(
    new THREE.SphereGeometry(0.04, 8, 8),
    new THREE.MeshBasicMaterial({ color: 0x00ff88, emissive: 0x00ff88, emissiveIntensity: 2 })
  )
  sensor.position.y = 1.35
  group.add(sensor)

  group.userData.isMeshGroup = true
  group.userData.type = 'dtPark'
  group.userData.update = function(elapsed) {
    sensor.material.emissiveIntensity = 1.5 + 0.8 * Math.sin(elapsed * 2.5)
  }
  return group
}

/**
 * dtStreetLightDT — 智慧路灯（带顶部传感器）
 */
function buildDTStreetLight(c) {
  const group = new THREE.Group()

  const pole = new THREE.Mesh(
    new THREE.CylinderGeometry(0.03, 0.05, 3.0, 8),
    new THREE.MeshStandardMaterial({ color: 0x444444, metalness: 0.7 })
  )
  pole.position.y = 1.5
  group.add(pole)

  // 灯头
  const lamp = new THREE.Mesh(
    new THREE.BoxGeometry(0.3, 0.08, 0.15),
    new THREE.MeshBasicMaterial({ color: 0xffffcc, emissive: 0xffeeaa, emissiveIntensity: 1.0 })
  )
  lamp.position.set(0.15, 3.0, 0)
  group.add(lamp)

  // 顶部传感器
  const sensor = new THREE.Mesh(
    new THREE.BoxGeometry(0.08, 0.04, 0.08),
    new THREE.MeshBasicMaterial({ color: 0x13c2c2, emissive: 0x13c2c2, emissiveIntensity: 1.5 })
  )
  sensor.position.y = 3.05
  group.add(sensor)

  group.userData.isMeshGroup = true
  group.userData.type = 'dtStreetLightDT'
  return group
}

/**
 * dtBaseStation — 基站/通信塔
 */
function buildDTBaseStation(c) {
  const group = new THREE.Group()

  // 主塔
  const tower = new THREE.Mesh(
    new THREE.CylinderGeometry(0.05, 0.08, 5.0, 8),
    new THREE.MeshStandardMaterial({ color: 0x666666, metalness: 0.6 })
  )
  tower.position.y = 2.5
  group.add(tower)

  // 天线阵列（顶部）
  for (let i = 0; i < 3; i++) {
    const ant = new THREE.Mesh(
      new THREE.BoxGeometry(0.6, 0.03, 0.08),
      new THREE.MeshStandardMaterial({ color: 0x888888, metalness: 0.7 })
    )
    ant.position.y = 4.8 + i * 0.15
    ant.rotation.y = (i * Math.PI) / 3
    group.add(ant)
  }

  // 顶部警示灯
  const beacon = new THREE.Mesh(
    new THREE.SphereGeometry(0.06, 16, 16),
    new THREE.MeshBasicMaterial({ color: 0xff3333, emissive: 0xff3333, emissiveIntensity: 2 })
  )
  beacon.position.y = 5.1
  group.add(beacon)

  // 底座
  const base = new THREE.Mesh(
    new THREE.BoxGeometry(0.8, 0.15, 0.8),
    new THREE.MeshStandardMaterial({ color: 0x444444 })
  )
  base.position.y = 0.075
  group.add(base)

  group.userData.isMeshGroup = true
  group.userData.type = 'dtBaseStation'
  group.userData.update = function(elapsed) {
    beacon.material.emissiveIntensity = 1.5 + 1.5 * Math.sin(elapsed * 4.0)
  }
  return group
}

export default {
  createDTObject,
  CUI,
}
