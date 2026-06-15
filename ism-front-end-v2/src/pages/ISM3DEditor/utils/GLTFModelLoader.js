/**
 * GLTFModelLoader — 增强型GLTF模型加载器
 *
 * 功能：
 * - 批量加载多个GLTF/GLB模型
 * - 加载进度回调
 * - Draco 压缩模型解码
 * - 自动 LOD 生成（简化模型层级）
 * - 模型缓存避免重复加载
 *
 * 用法：
 *   const loader = new GLTFModelLoader({ dracoPath: '/draco/' })
 *   const result = await loader.load('models/building.glb', (pct) => console.log(pct))
 *   // result: { scene, animations, ... }
 *
 *   // 批量加载
 *   const results = await loader.loadBatch(
 *     [{ url: 'a.glb' }, { url: 'b.glb' }],
 *     (pct, current, total) => console.log(pct)
 *   )
 */
import * as THREE from 'three'
import { GLTFLoader } from 'three/examples/jsm/loaders/GLTFLoader.js'
import { DRACOLoader } from 'three/examples/jsm/loaders/DRACOLoader.js'

export class GLTFModelLoader {
  /**
   * @param {Object} options
   * @param {string} [options.dracoPath='/draco/'] - Draco 解码器目录
   * @param {boolean} [options.enableLOD=false] - 是否自动生成 LOD
   * @param {number[]} [options.lodLevels=[0.5,0.25,0.1]] - LOD 简化比率
   */
  constructor(options = {}) {
    this.dracoPath = options.dracoPath || '/draco/'
    this.enableLOD = options.enableLOD || false
    this.lodLevels = options.lodLevels || [0.5, 0.25, 0.1]

    // 缓存
    this._cache = new Map()

    // GLTF Loader
    this._gltfLoader = new GLTFLoader()

    // Draco Loader
    this._dracoLoader = new DRACOLoader()
    this._dracoLoader.setDecoderPath(this.dracoPath)
    this._dracoLoader.setDecoderConfig({ type: 'js' })
    this._gltfLoader.setDRACOLoader(this._dracoLoader)
  }

  /**
   * 加载单个模型
   * @param {string} url
   * @param {Function} [onProgress] - (percent: 0-1)
   * @returns {Promise<GLTF>}
   */
  load(url, onProgress) {
    // 检查缓存
    if (this._cache.has(url)) {
      if (onProgress) onProgress(1)
      return Promise.resolve(this._cache.get(url))
    }

    return new Promise((resolve, reject) => {
      this._gltfLoader.load(
        url,
        (gltf) => {
          if (this.enableLOD) {
            this._addLOD(gltf.scene)
          }
          this._cache.set(url, gltf)
          resolve(gltf)
        },
        (xhr) => {
          if (onProgress && xhr.total > 0) {
            onProgress(xhr.loaded / xhr.total)
          }
        },
        (err) => reject(err)
      )
    })
  }

  /**
   * 批量加载模型
   * @param {Array<{url:string}>} items
   * @param {Function} [onProgress] - (overallPercent, currentIndex, total)
   * @returns {Promise<Array<GLTF>>}
   */
  async loadBatch(items, onProgress) {
    const results = []
    const total = items.length

    for (let i = 0; i < total; i++) {
      const gltf = await this.load(items[i].url, (itemPct) => {
        if (onProgress) {
          const overall = (i + itemPct) / total
          onProgress(overall, i + 1, total)
        }
      })
      results.push(gltf)
    }

    return results
  }

  /**
   * 为场景中的每个 Mesh 生成简化 LOD
   */
  _addLOD(scene) {
    scene.traverse((child) => {
      if (child.isMesh && child.geometry) {
        const lod = new THREE.LOD()
        // 原始精度 (0-10)
        lod.addLevel(child.clone(), 0)
        // 各级简化
        for (let i = 0; i < this.lodLevels.length; i++) {
          const simplified = this._simplifyGeometry(child.geometry, this.lodLevels[i])
          const mesh = new THREE.Mesh(simplified, child.material.clone())
          lod.addLevel(mesh, (i + 1) * 15) // 距离阈值递增
        }
        // 替换原mesh
        child.parent.add(lod)
        child.parent.remove(child)
      }
    })
  }

  /**
   * 简化几何体（使用 SimplifyModifier）
   */
  _simplifyGeometry(geometry, ratio) {
    if (!geometry.index) {
      // 非索引几何体直接返回低面数版本
      const targetCount = Math.floor(geometry.attributes.position.count * ratio)
      if (targetCount < 3) return geometry.clone()
      // 简单采样
      const newGeo = new THREE.BufferGeometry()
      const stride = Math.ceil(geometry.attributes.position.count / targetCount)
      const posArr = geometry.attributes.position.array
      const newPos = new Float32Array(targetCount * 3)
      for (let i = 0; i < targetCount; i++) {
        const src = i * stride * 3
        newPos[i*3] = posArr[src] || posArr[0]
        newPos[i*3+1] = posArr[src+1] || posArr[1]
        newPos[i*3+2] = posArr[src+2] || posArr[2]
      }
      newGeo.setAttribute('position', new THREE.BufferAttribute(newPos, 3))
      return newGeo
    }
    // 有索引的几何体：按比例减少三角形
    const idxArr = geometry.index.array
    const targetTriangles = Math.floor(idxArr.length / 3 * ratio)
    if (targetTriangles < 1) return geometry.clone()

    const newIdx = new Uint16Array(targetTriangles * 3)
    const stride = Math.ceil(idxArr.length / 3 / targetTriangles)
    for (let i = 0; i < targetTriangles; i++) {
      const src = i * stride * 3
      newIdx[i*3] = idxArr[src]
      newIdx[i*3+1] = idxArr[src+1]
      newIdx[i*3+2] = idxArr[src+2]
    }

    const newGeo = geometry.clone()
    newGeo.setIndex(new THREE.BufferAttribute(newIdx, 1))
    return newGeo
  }

  /**
   * 清除指定URL的缓存
   */
  clearCache(url) {
    if (url) {
      this._cache.delete(url)
    } else {
      this._cache.clear()
    }
  }

  /**
   * 销毁加载器
   */
  dispose() {
    this._dracoLoader.dispose()
    this._cache.clear()
  }
}

export default GLTFModelLoader
