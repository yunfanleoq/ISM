/**
 * Serializer - 场景序列化/反序列化
 */
import * as THREE from 'three'

/**
 * 序列化场景
 */
export function serializeScene(scene, options = {}) {
  const { includeMetadata = true } = options

  const data = {
    version: '1.0',
    metadata: includeMetadata ? {
      created: new Date().toISOString(),
      generator: 'ISM3DEditor'
    } : undefined,
    scene: serializeObject3D(scene)
  }

  return data
}

/**
 * 序列化 Object3D
 */
function serializeObject3D(object) {
  const data = {
    type: object.type || 'Object3D',
    name: object.name,
    uuid: object.uuid
  }

  // 变换
  data.position = object.position.toArray()
  data.rotation = [object.rotation.x, object.rotation.y, object.rotation.z]
  data.scale = object.scale.toArray()

  // 可见性
  data.visible = object.visible
  data.castShadow = object.castShadow
  data.receiveShadow = object.receiveShadow

  // 用户数据
  if (Object.keys(object.userData).length > 0) {
    data.userData = { ...object.userData }
  }

  // 材质
  if (object.material) {
    data.material = serializeMaterial(object.material)
  }

  // 几何体
  if (object.geometry) {
    data.geometry = serializeGeometry(object.geometry)
  }

  // 子对象
  if (object.children && object.children.length > 0) {
    data.children = object.children.map(child => serializeObject3D(child))
  }

  return data
}

/**
 * 序列化材质
 */
function serializeMaterial(material) {
  if (Array.isArray(material)) {
    return material.map(m => serializeMaterial(m))
  }

  return {
    type: material.type || 'MeshStandardMaterial',
    color: material.color ? '#' + material.color.getHexString() : '#ffffff',
    metalness: material.metalness ?? 0.3,
    roughness: material.roughness ?? 0.7,
    opacity: material.opacity ?? 1,
    transparent: material.transparent ?? false,
    wireframe: material.wireframe ?? false,
    emissive: material.emissive ? '#' + material.emissive.getHexString() : '#000000',
    emissiveIntensity: material.emissiveIntensity ?? 0
  }
}

/**
 * 序列化几何体
 */
function serializeGeometry(geometry) {
  const data = {
    type: geometry.type
  }

  // 根据几何体类型保存参数
  switch (geometry.type) {
    case 'BoxGeometry':
    case 'BoxBufferGeometry':
      data.parameters = {
        width: geometry.parameters.width,
        height: geometry.parameters.height,
        depth: geometry.parameters.depth
      }
      break

    case 'SphereGeometry':
    case 'SphereBufferGeometry':
      data.parameters = {
        radius: geometry.parameters.radius,
        widthSegments: geometry.parameters.widthSegments,
        heightSegments: geometry.parameters.heightSegments
      }
      break

    case 'CylinderGeometry':
    case 'CylinderBufferGeometry':
      data.parameters = {
        radiusTop: geometry.parameters.radiusTop,
        radiusBottom: geometry.parameters.radiusBottom,
        height: geometry.parameters.height,
        radialSegments: geometry.parameters.radialSegments
      }
      break

    case 'ConeGeometry':
    case 'ConeBufferGeometry':
      data.parameters = {
        radius: geometry.parameters.radius,
        height: geometry.parameters.height,
        radialSegments: geometry.parameters.radialSegments
      }
      break

    case 'TorusGeometry':
    case 'TorusBufferGeometry':
      data.parameters = {
        radius: geometry.parameters.radius,
        tube: geometry.parameters.tube,
        radialSegments: geometry.parameters.radialSegments,
        tubularSegments: geometry.parameters.tubularSegments
      }
      break

    case 'PlaneGeometry':
    case 'PlaneBufferGeometry':
      data.parameters = {
        width: geometry.parameters.width,
        height: geometry.parameters.height
      }
      break
  }

  return data
}

/**
 * 反序列化场景
 */
export function deserializeScene(data, scene) {
  if (!data || !data.scene) {
    throw new Error('Invalid scene data')
  }

  // 验证版本
  if (data.version && data.version !== '1.0') {
    console.warn(`Scene version mismatch: ${data.version}`)
  }

  // 清空现有场景
  while (scene.children.length > 0) {
    scene.remove(scene.children[0])
  }

  // 解析场景
  deserializeObject3D(data.scene, scene)

  return scene
}

/**
 * 反序列化 Object3D
 */
function deserializeObject3D(data, object) {
  // 名称
  if (data.name) {
    object.name = data.name
  }

  // 变换
  if (data.position) {
    object.position.fromArray(data.position)
  }
  if (data.rotation) {
    object.rotation.set(data.rotation[0], data.rotation[1], data.rotation[2])
  }
  if (data.scale) {
    object.scale.fromArray(data.scale)
  }

  // 可见性
  if (data.visible !== undefined) {
    object.visible = data.visible
  }
  if (data.castShadow !== undefined) {
    object.castShadow = data.castShadow
  }
  if (data.receiveShadow !== undefined) {
    object.receiveShadow = data.receiveShadow
  }

  // 用户数据
  if (data.userData) {
    Object.assign(object.userData, data.userData)
  }

  // 材质
  if (data.material) {
    object.material = deserializeMaterial(data.material)
  }

  // 几何体
  if (data.geometry) {
    object.geometry = deserializeGeometry(data.geometry)
  }

  // 子对象
  if (data.children && Array.isArray(data.children)) {
    data.children.forEach(childData => {
      const child = new THREE.Object3D()
      deserializeObject3D(childData, child)
      object.add(child)
    })
  }

  return object
}

/**
 * 反序列化材质
 */
function deserializeMaterial(data) {
  if (Array.isArray(data)) {
    return data.map(d => deserializeMaterial(d))
  }

  const material = new THREE.MeshStandardMaterial({
    color: new THREE.Color(data.color || '#ffffff'),
    metalness: data.metalness ?? 0.3,
    roughness: data.roughness ?? 0.7,
    opacity: data.opacity ?? 1,
    transparent: data.transparent ?? false,
    wireframe: data.wireframe ?? false
  })

  if (data.emissive) {
    material.emissive = new THREE.Color(data.emissive)
    material.emissiveIntensity = data.emissiveIntensity ?? 0
  }

  return material
}

/**
 * 反序列化几何体
 */
function deserializeGeometry(data) {
  const { type, parameters } = data

  if (!parameters) {
    return new THREE.BoxGeometry(1, 1, 1)
  }

  switch (type) {
    case 'BoxGeometry':
    case 'BoxBufferGeometry':
      return new THREE.BoxGeometry(
        parameters.width || 1,
        parameters.height || 1,
        parameters.depth || 1
      )

    case 'SphereGeometry':
    case 'SphereBufferGeometry':
      return new THREE.SphereGeometry(
        parameters.radius || 0.5,
        parameters.widthSegments || 32,
        parameters.heightSegments || 32
      )

    case 'CylinderGeometry':
    case 'CylinderBufferGeometry':
      return new THREE.CylinderGeometry(
        parameters.radiusTop ?? 0.5,
        parameters.radiusBottom ?? 0.5,
        parameters.height || 2,
        parameters.radialSegments || 32
      )

    case 'ConeGeometry':
    case 'ConeBufferGeometry':
      return new THREE.ConeGeometry(
        parameters.radius || 0.5,
        parameters.height || 2,
        parameters.radialSegments || 32
      )

    case 'TorusGeometry':
    case 'TorusBufferGeometry':
      return new THREE.TorusGeometry(
        parameters.radius || 1,
        parameters.tube || 0.3,
        parameters.radialSegments || 16,
        parameters.tubularSegments || 100
      )

    case 'PlaneGeometry':
    case 'PlaneBufferGeometry':
      return new THREE.PlaneGeometry(
        parameters.width || 1,
        parameters.height || 1
      )

    default:
      return new THREE.BoxGeometry(1, 1, 1)
  }
}

/**
 * 导出为 JSON 字符串
 */
export function exportToJSON(scene, options) {
  const data = serializeScene(scene, options)
  return JSON.stringify(data, null, 2)
}

/**
 * 从 JSON 字符串导入
 */
export function importFromJSON(jsonString, scene) {
  const data = JSON.parse(jsonString)
  return deserializeScene(data, scene)
}
