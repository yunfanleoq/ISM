/**
 * 真正的3D几何体文字 - 用 opentype.js 解析字体 + ExtrudeGeometry 挤出
 * 替换之前的 Canvas 贴图方案，实现翠鸟风格的无背景3D文字
 */

import * as THREE from 'three'

// opentype 实例（由外部加载后注入）
let _opentype = null
let _loadedFonts = {} // 缓存已加载的字体

/**
 * 设置 opentype.js 实例（在外部加载后调用一次）
 */
export function setOpentype(opentype) {
  _opentype = opentype
}

/**
 * 加载字体文件（TTF/OTF）
 * @param {string} url - 字体文件 URL
 * @returns {Promise} - opentype.Font 对象
 */
export function loadFont(url) {
  if (_loadedFonts[url]) {
    return Promise.resolve(_loadedFonts[url])
  }
  if (!_opentype) {
    return Promise.reject(new Error('opentype.js not loaded. Call setOpentype first.'))
  }
  return new Promise(function (resolve, reject) {
    _opentype.load(url, function (err, font) {
      if (err) {
        reject(err)
      } else {
        _loadedFonts[url] = font
        resolve(font)
      }
    })
  })
}

/**
 * 从 ArrayBuffer 解析字体（用于用户上传的字体文件）
 */
export function parseFontBuffer(buffer) {
  if (!_opentype) {
    throw new Error('opentype.js not loaded')
  }
  return _opentype.parse(buffer)
}

/**
 * 将 opentype.js 的路径命令转换为 Three.js Shape
 * @param {Array} pathCommands - opentype.js 路径命令
 * @param {number} scale - 缩放比例
 * @returns {THREE.Shape[]} - Three.js Shape 数组
 */
function convertPathToShapes(pathCommands, scale) {
  var shapes = []
  var currentShape = null
  var currentX = 0
  var currentY = 0

  for (var i = 0; i < pathCommands.length; i++) {
    var cmd = pathCommands[i]
    var type = cmd.type
    var x = cmd.x !== undefined ? cmd.x * scale : currentX
    var y = cmd.y !== undefined ? -cmd.y * scale : currentY // Y轴翻转

    if (type === 'M') {
      // 开始新路径
      if (currentShape && currentShape.getPoints().length > 0) {
        shapes.push(currentShape)
      }
      currentShape = new THREE.Shape()
      currentShape.moveTo(x, y)
      currentX = x
      currentY = y
    } else if (type === 'L') {
      // 直线
      if (currentShape) {
        currentShape.lineTo(x, y)
      }
      currentX = x
      currentY = y
    } else if (type === 'Q') {
      // 二次贝塞尔曲线
      if (currentShape && cmd.x1 !== undefined) {
        var x1 = cmd.x1 * scale
        var y1 = -cmd.y1 * scale
        currentShape.quadraticCurveTo(x1, y1, x, y)
      }
      currentX = x
      currentY = y
    } else if (type === 'C') {
      // 三次贝塞尔曲线
      if (currentShape && cmd.x1 !== undefined && cmd.x2 !== undefined) {
        var cx1 = cmd.x1 * scale
        var cy1 = -cmd.y1 * scale
        var cx2 = cmd.x2 * scale
        var cy2 = -cmd.y2 * scale
        currentShape.bezierCurveTo(cx1, cy1, cx2, cy2, x, y)
      }
      currentX = x
      currentY = y
    } else if (type === 'Z') {
      // 闭合路径
      if (currentShape) {
        currentShape.closePath()
      }
    }
  }

  if (currentShape && currentShape.getPoints().length > 0) {
    shapes.push(currentShape)
  }

  return shapes
}

/**
 * 用 opentype.js 字体生成 Three.js 文字 Shape
 * @param {opentype.Font} font - opentype.js 字体对象
 * @param {string} text - 文字内容
 * @param {number} fontSize - 字体大小（世界单位）
 * @returns {THREE.Shape[]} - Three.js Shape 数组
 */
export function textToShapes(font, text, fontSize) {
  if (!font || !text) return []

  // opentype.js 的单位是 em，需要缩放 to 世界单位
  // 通常 1em = 字体文件的 unitsPerEm（通常是1000或2048）
  var scale = fontSize / font.unitsPerEm

  var path = font.getPath(text, 0, 0, fontSize)
  // path.commands 包含路径命令
  var shapes = convertPathToShapes(path.commands, scale)

  return shapes
}

/**
 * 用 ExtrudeGeometry 生成3D文字网格
 * @param {opentype.Font} font - opentype.js 字体对象
 * @param {string} text - 文字内容
 * @param {object} options - 配置选项
 * @param {number} options.fontSize - 字体大小（世界单位，默认 0.8）
 * @param {number} options.depth - 挤出深度（默认 fontSize * 0.12）
 * @param {number} options.curveSegments - 曲线分段数（默认 12）
 * @param {string} options.frontColor - 正面颜色（默认 '#ffffff'）
 * @param {string} options.sideColor - 侧面颜色（默认 '#888888'）
 * @param {boolean} options.bevelEnabled - 是否启用倒角（默认 false）
 * @returns {THREE.Mesh} - 3D文字网格
 */
export function createTrue3DTextMesh(font, text, options) {
  options = options || {}

  var fontSize = options.fontSize || 0.8
  var depth = options.depth || fontSize * 0.18
  var curveSegments = options.curveSegments || 12
  var frontColor = new THREE.Color(options.frontColor || '#ffffff')
  var sideColor = new THREE.Color(options.sideColor || '#888888')

  // 生成文字 Shape
  var shapes = textToShapes(font, text, fontSize)
  if (!shapes || shapes.length === 0) {
    //  fallback: 返回一个空的 Group
    return new THREE.Group()
  }

  // ExtrudeGeometry 参数
  var extrudeSettings = {
    depth: depth,
    bevelEnabled: options.bevelEnabled || false,
    bevelThickness: options.bevelThickness || 0.02,
    bevelSize: options.bevelSize || 0.01,
    bevelOffset: 0,
    bevelSegments: 3,
    curveSegments: curveSegments,
    steps: 1
  }

  // 创建 ExtrudeGeometry
  var geo = new THREE.ExtrudeGeometry(shapes, extrudeSettings)

  // 计算边界框，居中几何体
  geo.computeBoundingBox()
  var box = geo.boundingBox
  var centerX = (box.min.x + box.max.x) / 2
  var centerY = (box.min.y + box.max.y) / 2
  geo.translate(-centerX, -centerY, -depth / 2)

  // 正面材质（双面，让文字正反都能看到）
  var matFront = new THREE.MeshStandardMaterial({
    color: frontColor,
    metalness: options.metalness || 0.1,
    roughness: options.roughness || 0.45,
    side: THREE.DoubleSide
  })

  // 侧面材质
  var matSide = new THREE.MeshStandardMaterial({
    color: sideColor,
    metalness: options.metalnessSide || 0.7,
    roughness: options.roughnessSide || 0.25,
    side: THREE.DoubleSide
  })

  // ExtrudeGeometry 会自动分组：正面/背面 + 侧面
  // 用 groups 给不同面分配不同材质
  var mesh = new THREE.Mesh(geo, [matFront, matSide])
  mesh.castShadow = true
  mesh.receiveShadow = true

  // 标记
  mesh.userData.is3DText = true
  mesh.userData.isTrue3DText = true // 真正的3D几何体文字
  mesh.userData.textData = {
    text: text,
    fontSize: fontSize,
    depth: depth
  }
  mesh.name = 'True3DText_' + text.slice(0, 10)

  return mesh
}

/**
 * 默认中文字体 URL（放在 public/fonts/ 下）
 * 用户需要自行放置中文字体文件，或配置此 URL
 */
export var DEFAULT_CHINESE_FONT_URL = '/fonts/NotoSansSC-Regular.ttf'

/**
 * 初始化默认字体（在应用启动时调用一次）
 */
export function initDefaultFont() {
  return loadFont(DEFAULT_CHINESE_FONT_URL).catch(function (err) {
    console.warn('[Text3DGeometry] Default Chinese font not loaded:', err.message)
    console.warn('[Text3DGeometry] Please place a Chinese TTF font at:', DEFAULT_CHINESE_FONT_URL)
    return null
  })
}
