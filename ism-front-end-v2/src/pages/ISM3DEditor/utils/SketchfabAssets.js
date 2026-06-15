/**
 * SketchfabAssets — 外部资产清单加载 + Sketchfab API 搜索 + 模型加载辅助
 *
 * 用法：
 *   import { getAssetList, searchSketchfab, loadAssetIntoScene } from './SketchfabAssets'
 *   const list = await getAssetList()
 *   const results = await searchSketchfab('building', { token: 'xxx' })
 *   const group = await loadAssetIntoScene(asset, scene, camera, renderer)
 */
import * as THREE from 'three'
import { GLTFLoader } from 'three/examples/jsm/loaders/GLTFLoader.js'
import { DRACOLoader } from 'three/examples/jsm/loaders/DRACOLoader.js'

const MANIFEST_URL = '/static/ISM/sketchfab-assets/manifest.json'

let _manifest = null
let _modelCache = new Map()
let _searchCache = new Map()

function _stripJsonComments(text) {
  return String(text || '')
    .replace(/\/\*[\s\S]*?\*\//g, '')
    .replace(/^\s*\/\/.*$/gm, '')
}

function _normalizeManifest(data) {
  const models = Array.isArray(data.models)
    ? data.models
    : (Array.isArray(data.assets) ? data.assets : [])
  let categories = Array.isArray(data.categories) ? data.categories : []
  if (!categories.length) {
    const seen = {}
    categories = models
      .map(asset => asset && asset.category)
      .filter(Boolean)
      .filter(key => {
        if (seen[key]) return false
        seen[key] = true
        return true
      })
      .map(key => ({ key, label: key }))
  }
  return {
    ...data,
    models,
    assets: models,
    categories
  }
}

function _loadManifest() {
  if (_manifest) return Promise.resolve(_manifest)
  return fetch(MANIFEST_URL + '?t=' + Date.now())
    .then(r => {
      if (!r.ok) {
        console.warn('[SketchfabAssets] manifest.json not found (' + r.status + '), using empty manifest')
        return '{"models":[],"categories":[]}'
      }
      return r.text()
    })
    .then(text => {
      try {
        return JSON.parse(_stripJsonComments(text))
      } catch (e) {
        console.warn('[SketchfabAssets] manifest.json parse error, using empty manifest')
        return { models: [], categories: [] }
      }
    })
    .then(data => {
      _manifest = _normalizeManifest(data)
      return _manifest
    })
    .catch(err => {
      console.warn('[SketchfabAssets] manifest load failed', err)
      _manifest = { models: [], categories: [] }
      return _manifest
    })
}

/**
 * 获取全部资产列表
 * @returns {Promise<Array>}
 */
export function getAssetList() {
  return _loadManifest().then(m => (m.models || []))
}

/**
 * 获取分类列表
 * @returns {Promise<Array>}
 */
export function getCategories() {
  return _loadManifest().then(m => (m.categories || []))
}

/**
 * 按分类/关键字过滤资产
 * @param {Object} opts
 * @param {string} [opts.category]  - 分类 key
 * @param {string} [opts.q]        - 搜索关键字
 * @returns {Promise<Array>}
 */
export function filterAssets(opts = {}) {
  return getAssetList().then(list => {
    let res = list.slice()
    if (opts.category && opts.category !== 'all') {
      res = res.filter(a => a.category === opts.category)
    }
    if (opts.q) {
      const q = opts.q.toLowerCase()
      res = res.filter(a => {
        const hay = (a.name + ' ' + (a.tags || []).join(' ') + ' ' + (a.description || '')).toLowerCase()
        return hay.includes(q)
      })
    }
    return res
  })
}

/**
 * 从 manifest 的 modelUrl 加载 GLTF 模型
 * @param {Object}   asset      - manifest 中的单个资产对象
 * @param {THREE.Scene} scene  - THREE.Scene（可选，传入则自动加入）
 * @param {Function} [onProgress] - 进度回调 (0-1)
 * @returns {Promise<THREE.Group>}
 */
export function loadAssetModel(asset, scene, onProgress) {
  if (!asset || !asset.modelUrl) {
    return Promise.reject(new Error('invalid asset: missing modelUrl'))
  }
  const url = asset.modelUrl

  // 友好检测：如果 modelUrl 不是 .glb/.gltf 文件，提示用户需要手动下载
  if (url && !/\.(glb|gltf)(\?|$)/i.test(url)) {
    const hint = [
      '[SketchfabAssets] 模型 URL 不是可直接加载的 .glb/.gltf 文件：',
      url,
      '',
      '请按以下步骤操作：',
      '1. 点击编辑器左侧"外部资产"里的下载按钮（外部链接图标），在 Sketchfab 页面登录并下载 .glb 文件',
      '2. 将下载的文件放到：public/static/ISM/sketchfab-assets/models/ 目录',
      '3. 更新 manifest.json 里对应条目的 modelUrl 为本地路径，例如：/static/ISM/sketchfab-assets/models/xxx.glb',
    ].join('\n')
    console.warn(hint)
    return Promise.reject(new Error(
      '模型文件未就绪，请先从 Sketchfab 下载 .glb 文件（详见控制台提示）'
    ))
  }

  // 缓存
  if (_modelCache.has(url)) {
    const clone = _modelCache.get(url).clone()
    if (scene) scene.add(clone)
    return Promise.resolve(clone)
  }

  return new Promise((resolve, reject) => {
    const dracoLoader = new DRACOLoader()
    dracoLoader.setDecoderPath('/draco/')
    dracoLoader.setDecoderConfig({ type: 'js' })

    const loader = new GLTFLoader()
    loader.setDRACOLoader(dracoLoader)

    loader.load(
      url,
      (gltf) => {
        const model = gltf.scene
        // 应用 manifest 中的 transform 提示
        if (asset.scale && asset.scale !== 1.0) {
          model.scale.setScalar(asset.scale)
        }
        if (asset.rotationY) {
          model.rotation.y = THREE.MathUtils.degToRad(asset.rotationY)
        }
        // 标准化：计算包围盒并居中到原点
        _normalizeModel(model)

        // 缓存原始 scene（注意：不能直接缓存 model，因为它是 Group，clone 更安全）
        const clonedForCache = model.clone()
        _modelCache.set(url, clonedForCache)

        if (scene) scene.add(model)
        resolve(model)
      },
      (xhr) => {
        if (onProgress && xhr.total > 0) {
          onProgress(xhr.loaded / xhr.total)
        }
      },
      (err) => {
        reject(err)
      }
    )
  })
}

/**
 * 将外部模型封装成 ISM3DEditor 能识别的 obj 结构
 * （用于把加载结果写回 vm.__3d.objects）
 *
 * @param {THREE.Group} model      - loadAssetModel 返回的 model
 * @param {Object}       asset      - manifest 资产条目
 * @param {string}       [id]       - 可选，指定对象 ID
 * @returns {Object} 符合 ISM3DEditor 格式的 obj
 */
export function wrapAssetToObject(model, asset, id) {
  const objId = id || ('skf_' + asset.id + '_' + Date.now())
  // 记录 model 的 world matrix，方便序列化（我们只记 position/rotation/scale）
  const pos = model.position
  const rot = model.rotation
  const s = model.scale

  return {
    id: objId,
    name: asset.name || '外部模型',
    type: 'sketchfab-model',
    assetId: asset.id,
    modelUrl: asset.modelUrl,
    thumbnail: asset.thumbnail || '',
    visible: true,
    x: +pos.x.toFixed(4),
    y: +pos.y.toFixed(4),
    z: +pos.z.toFixed(4),
    rx: +(THREE.MathUtils.radToDeg(rot.x)).toFixed(2),
    ry: +(THREE.MathUtils.radToDeg(rot.y)).toFixed(2),
    rz: +(THREE.MathUtils.radToDeg(rot.z)).toFixed(2),
    sx: +s.x.toFixed(4),
    sy: +s.y.toFixed(4),
    sz: +s.z.toFixed(4),
    opacity: 1,
    castShadow: true,
    receiveShadow: true,
    userData: {
      isMeshGroup: true,
      isSketchfabModel: true,
      assetId: asset.id
    }
  }
}

function _normalizeModel(model) {
  // 计算包围盒中心，把模型中心移到原点
  const box = new THREE.Box3().setFromObject(model)
  const center = new THREE.Vector3()
  box.getCenter(center)
  // 把模型底部放在 y=0
  const size = new THREE.Vector3()
  box.getSize(size)
  model.position.x -= center.x
  model.position.y -= box.min.y
  model.position.z -= center.z
}

/**
 * 搜索 Sketchfab 公开模型（不需要 API token）
 * @param {string} query - 搜索关键词
 * @param {Object} [opts]
 * @param {number} [opts.limit=24] - 返回结果数量
 * @param {string} [opts.token] - Sketchfab API token（可选，提供后可获取下载链接）
 * @returns {Promise<Array>} 搜索结果数组
 *
 * 返回每项结构：
 * {
 *   uid, name, description,
 *   thumbnails: { url: string }[],
 *   viewerUrl: string,
 *   downloadUrl: string|null,  // 仅当传入 token 时可能非空
 *   polygonCount: number,
 *   author: { username, avatar },
 *   license: string,
 *   animated: boolean,
 *   hasEmbedded: boolean,
 * }
 */
export function searchSketchfab(query, opts = {}) {
  if (!query || !query.trim()) {
    return Promise.resolve([])
  }
  const limit = opts.limit || 24
  const cacheKey = 'sf_search_' + query.trim().toLowerCase() + '_' + limit
  if (_searchCache.has(cacheKey)) {
    return Promise.resolve(_searchCache.get(cacheKey))
  }

  const url = 'https://api.sketchfab.com/v3/search?type=models&downloadable=true&staffpicked=true&q='
    + encodeURIComponent(query.trim())
    + '&count=' + limit

  return fetch(url, {
    headers: opts.token ? { 'Authorization': 'Token ' + opts.token } : {}
  })
    .then(r => r.json())
    .then(data => {
      const results = (data.results || []).map(item => {
        // 取最大尺寸缩略图
        let thumbUrl = ''
        if (Array.isArray(item.thumbnails) && item.thumbnails.length) {
          const thumbs = item.thumbnails.slice().sort((a, b) => (b.width || 0) - (a.width || 0))
          thumbUrl = thumbs[0].url || ''
        }

        return {
          uid: item.uid || '',
          name: item.name || 'Untitled',
          description: item.description || '',
          thumbnail: thumbUrl,
          viewerUrl: item.viewerUrl || ('https://sketchfab.com/models/' + (item.uid || '')),
          downloadUrl: '',   // 需要单独调用下载 API
          polygonCount: item.faceCount || item.vertexCount || 0,
          author: {
            username: (item.user && item.user.username) || '',
            avatar: (item.user && item.user.avatar) || ''
          },
          license: (item.license && item.license.slug) || 'unknown',
          animated: !!item.animated,
          hasEmbedded: !!item.hasEmbedded,
          _raw: item  // 保留原始数据以供调试
        }
      })
      _searchCache.set(cacheKey, results)
      return results
    })
    .catch(err => {
      console.warn('[SketchfabAssets] 搜索失败', err)
      return []
    })
}

/**
 * 获取 Sketchfab 模型的 glTF 下载链接（需要 API token）
 * @param {string} uid - Sketchfab 模型 UID
 * @param {string} apiToken - Sketchfab API token
 * @returns {Promise<string>} glTF 下载 URL
 */
export function getSketchfabDownloadUrl(uid, apiToken) {
  if (!uid || !apiToken) {
    return Promise.reject(new Error('uid 和 apiToken 均为必填'))
  }
  const url = 'https://api.sketchfab.com/v3/models/' + uid + '/download'
  return fetch(url, {
    headers: { 'Authorization': 'Token ' + apiToken }
  })
    .then(r => r.json())
    .then(data => {
      // 返回 glTF 格式的下载 URL
      const gltf = data.gltf || data.glb || data
      const dlUrl = (gltf && gltf.url) || ''
      if (!dlUrl) {
        return Promise.reject(new Error('无法获取下载链接，模型可能未开启下载权限'))
      }
      return dlUrl
    })
}

/**
 * 从 URL 加载 glTF/GLB 模型（支持 Sketchfab CDN 或任意可访问 URL）
 * @param {string}   url             - 模型 URL（.gltf 或 .glb）
 * @param {THREE.Scene} [scene]    - 传入则自动添加到场景
 * @param {Function}      [onProgress] - 进度回调 (0-1)
 * @param {Object}       [opts]
 * @param {string}       [opts.dracoPath] - Draco 解码器路径，默认 '/draco/'
 * @returns {Promise<THREE.Group>}
 */
export function loadModelFromUrl(url, scene, onProgress, opts = {}) {
  if (!url) {
    return Promise.reject(new Error('url 不能为空'))
  }
  if (_modelCache.has(url)) {
    const clone = _modelCache.get(url).clone()
    if (scene) scene.add(clone)
    return Promise.resolve(clone)
  }

  return new Promise((resolve, reject) => {
    const dracoLoader = new DRACOLoader()
    dracoLoader.setDecoderPath(opts.dracoPath || '/draco/')
    dracoLoader.setDecoderConfig({ type: 'js' })

    const loader = new GLTFLoader()
    loader.setDRACOLoader(dracoLoader)

    loader.load(
      url,
      (gltf) => {
        const model = gltf.scene
        _normalizeModel(model)

        // 缓存原始 scene 的克隆
        const clonedForCache = model.clone()
        _modelCache.set(url, clonedForCache)

        // 确保材质支持 Bloom 辉光
        model.traverse(child => {
          if (child.isMesh && child.material) {
            const mat = child.material
            if (mat.emissive && mat.emissiveIntensity !== undefined) {
              // 已支持 Bloom
            } else if (mat.emissive) {
              mat.emissiveIntensity = 1.0
            }
            // 透明对象
            if (mat.transparent) {
              mat.depthWrite = false
            }
          }
        })

        if (scene) scene.add(model)
        resolve(model)
      },
      (xhr) => {
        if (onProgress && xhr.total > 0) {
          onProgress(xhr.loaded / xhr.total)
        }
      },
      (err) => {
        reject(err)
      }
    )
  })
}
