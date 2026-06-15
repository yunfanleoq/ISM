/**
 * ThreeHelper - Three.js 辅助工具
 */
import * as THREE from 'three'

/**
 * 创建材质
 */
export function createMaterial(options = {}) {
  const {
    color = 0x4a90d9,
    metalness = 0.3,
    roughness = 0.7,
    opacity = 1,
    transparent = false,
    wireframe = false
  } = options

  return new THREE.MeshStandardMaterial({
    color,
    metalness,
    roughness,
    opacity,
    transparent,
    wireframe
  })
}

/**
 * 创建几何体
 */
export function createGeometry(type, params = {}) {
  switch (type) {
    case 'box':
      return new THREE.BoxGeometry(
        params.width || 1,
        params.height || 1,
        params.depth || 1
      )
    case 'sphere':
      return new THREE.SphereGeometry(
        params.radius || 0.5,
        params.widthSegments || 32,
        params.heightSegments || 32
      )
    case 'cylinder':
      return new THREE.CylinderGeometry(
        params.radiusTop || 0.5,
        params.radiusBottom || 0.5,
        params.height || 2,
        params.radialSegments || 32
      )
    case 'cone':
      return new THREE.ConeGeometry(
        params.radius || 0.5,
        params.height || 2,
        params.radialSegments || 32
      )
    case 'torus':
      return new THREE.TorusGeometry(
        params.radius || 1,
        params.tube || 0.3,
        params.radialSegments || 16,
        params.tubularSegments || 100
      )
    case 'plane':
      return new THREE.PlaneGeometry(
        params.width || 1,
        params.height || 1
      )
    case 'circle':
      return new THREE.CircleGeometry(
        params.radius || 1,
        params.segments || 32
      )
    default:
      return new THREE.BoxGeometry(1, 1, 1)
  }
}

/**
 * 创建网格
 */
export function createMesh(type, params = {}) {
  const geometry = createGeometry(type, params)
  const material = createMaterial(params)
  const mesh = new THREE.Mesh(geometry, material)

  mesh.castShadow = true
  mesh.receiveShadow = true

  return mesh
}

/**
 * 计算包围盒
 */
export function getBoundingBox(object) {
  const box = new THREE.Box3()
  box.setFromObject(object)
  return box
}

/**
 * 获取中心点
 */
export function getCenter(object) {
  const box = getBoundingBox(object)
  return box.getCenter(new THREE.Vector3())
}

/**
 * 获取尺寸
 */
export function getSize(object) {
  const box = getBoundingBox(object)
  return box.getSize(new THREE.Vector3())
}

/**
 * 复制对象
 */
export function cloneObject(object) {
  const cloned = object.clone()

  if (object.geometry) {
    cloned.geometry = object.geometry.clone()
  }
  if (object.material) {
    if (Array.isArray(object.material)) {
      cloned.material = object.material.map(m => m.clone())
    } else {
      cloned.material = object.material.clone()
    }
  }

  return cloned
}

/**
 * 销毁对象
 */
export function disposeObject(object) {
  if (object.geometry) {
    object.geometry.dispose()
  }

  if (object.material) {
    if (Array.isArray(object.material)) {
      object.material.forEach(m => m.dispose())
    } else {
      object.material.dispose()
    }
  }
}

/**
 * 颜色转换
 */
export function hexToColor(hex) {
  return new THREE.Color(hex)
}

/**
 * 颜色转十六进制
 */
export function colorToHex(color) {
  return '#' + color.getHexString()
}

/**
 * 生成唯一ID
 */
export function generateId(prefix = 'obj') {
  return `${prefix}_${Math.random().toString(36).substr(2, 9)}`
}

/**
 * 深拷贝
 */
export function deepClone(obj) {
  return JSON.parse(JSON.stringify(obj))
}

/**
 * 对象合并
 */
export function merge(target, source) {
  return Object.assign({}, target, source)
}

/**
 * 验证数据
 */
export function validateObject(object) {
  const errors = []

  if (!object.type) {
    errors.push('缺少 type 属性')
  }

  if (!object.transform) {
    errors.push('缺少 transform 属性')
  } else {
    if (!object.transform.position) {
      errors.push('缺少 position 属性')
    }
    if (!object.transform.rotation) {
      errors.push('缺少 rotation 属性')
    }
    if (!object.transform.scale) {
      errors.push('缺少 scale 属性')
    }
  }

  return {
    valid: errors.length === 0,
    errors
  }
}

/**
 * 坐标转换：屏幕坐标转世界坐标
 */
export function screenToWorld(screenX, screenY, camera, width, height) {
  const vector = new THREE.Vector3(
    (screenX / width) * 2 - 1,
    -(screenY / height) * 2 + 1,
    0.5
  )
  vector.unproject(camera)
  return vector
}

/**
 * 坐标转换：世界坐标转屏幕坐标
 */
export function worldToScreen(worldPos, camera, width, height) {
  const vector = worldPos.clone()
  vector.project(camera)

  return {
    x: (vector.x + 1) / 2 * width,
    y: -(vector.y - 1) / 2 * height,
    z: vector.z
  }
}

/**
 * 插值
 */
export function lerp(a, b, t) {
  return a + (b - a) * t
}

/**
 * 颜色插值
 */
export function lerpColor(color1, color2, t) {
  const c1 = new THREE.Color(color1)
  const c2 = new THREE.Color(color2)
  return c1.lerp(c2, t)
}

/**
 * 限制数值范围
 */
export function clamp(value, min, max) {
  return Math.min(Math.max(value, min), max)
}

/**
 * 映射数值范围
 */
export function mapRange(value, inMin, inMax, outMin, outMax) {
  return ((value - inMin) * (outMax - outMin)) / (inMax - inMin) + outMin
}

/**
 * 平滑阻尼
 */
export function damp(current, target, alpha, dt) {
  return lerp(current, target, 1 - Math.pow(alpha, dt))
}
