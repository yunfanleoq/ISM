/**
 * Text3DFactory - 翠鸟风格 3D 文字几何体工厂
 * 使用 opentype.js 直接解析中文 TTF 字体，生成 THREE.Shape → ExtrudeGeometry
 * 替代原有的 Canvas+Sprite 方案，实现真正的 3D 立体文字
 */
import * as THREE from 'three'

// opentype.js 通过 import-from 加载 ESM 模块
let opentypeModule = null
function getOpentype() {
  if (!opentypeModule) {
    opentypeModule = require('opentype.js')
  }
  return opentypeModule
}

// ===== 字体缓存 =====
var fontCache = {}

/**
 * 加载 TTF 字体 (异步)
 * @param {string} url - TTF 文件路径
 * @returns {Promise<opentype.Font>}
 */
export function loadFont(url) {
  if (fontCache[url]) {
    return Promise.resolve(fontCache[url])
  }
  var opentype = getOpentype()
  return new Promise(function(resolve, reject) {
    opentype.load(url, function(err, font) {
      if (err) {
        console.error('[Text3DFactory] 字体加载失败:', url, err)
        reject(err)
        return
      }
      fontCache[url] = font
      console.log('[Text3DFactory] 字体已加载:', font.names.fontFamily.en, url)
      resolve(font)
    })
  })
}

/**
 * 将 opentype.js glyph path 转换为 THREE.Shape
 * @param {opentype.Font} font
 * @param {opentype.Glyph} glyph
 * @param {number} scale - 缩放
 * @param {number} offsetX - X偏移
 * @param {number} offsetY - Y偏移
 * @returns {THREE.Shape}
 */
function glyphToShape(font, glyph, scale, offsetX, offsetY) {
  var path = glyph.getPath(0, 0, font.unitsPerEm)
  var shape = new THREE.Shape()
  var holeShapes = []
  var currentShape = shape
  var hasHoles = false

  var commands = path.commands || []
  if (commands.length === 0) return shape

  for (var i = 0; i < commands.length; i++) {
    var cmd = commands[i]
    var cmdType = cmd.type

    if (cmdType === 'M') {
      // 如果当前 shape 已有路径，新 'M' 表示开始一个洞
      if (i > 0 && currentShape.currentPoint) {
        hasHoles = true
        currentShape = new THREE.Path()
        holeShapes.push(currentShape)
      }
      currentShape.moveTo(
        cmd.x * scale + offsetX,
        (font.unitsPerEm - cmd.y) * scale + offsetY
      )
    } else if (cmdType === 'L') {
      currentShape.lineTo(
        cmd.x * scale + offsetX,
        (font.unitsPerEm - cmd.y) * scale + offsetY
      )
    } else if (cmdType === 'Q') {
      currentShape.quadraticCurveTo(
        cmd.x1 * scale + offsetX,
        (font.unitsPerEm - cmd.y1) * scale + offsetY,
        cmd.x * scale + offsetX,
        (font.unitsPerEm - cmd.y) * scale + offsetY
      )
    } else if (cmdType === 'C') {
      currentShape.bezierCurveTo(
        cmd.x1 * scale + offsetX,
        (font.unitsPerEm - cmd.y1) * scale + offsetY,
        cmd.x2 * scale + offsetX,
        (font.unitsPerEm - cmd.y2) * scale + offsetY,
        cmd.x * scale + offsetX,
        (font.unitsPerEm - cmd.y) * scale + offsetY
      )
    } else if (cmdType === 'Z') {
      currentShape.closePath()
    }
  }

  // 将洞形状合并到主形状
  if (hasHoles) {
    for (var h = 0; h < holeShapes.length; h++) {
      shape.holes.push(holeShapes[h])
    }
  }

  return shape
}

/**
 * 将文本转换为 THREE.Shape 数组
 * @param {string} text
 * @param {opentype.Font} font
 * @param {number} size - 文字大小
 * @returns {THREE.Shape[]}
 */
export function textToShapes(text, font, size) {
  if (!text || !font) return []

  size = size || 2
  var scale = size / font.unitsPerEm
  var shapes = []
  var offsetX = 0

  // 获取行高
  var ascender = font.ascender || font.unitsPerEm * 0.8
  var descender = font.descender || -font.unitsPerEm * 0.2
  var lineHeight = (ascender - descender) * scale * 1.2

  var chars = String(text).split('')

  for (var i = 0; i < chars.length; i++) {
    var ch = chars[i]

    if (ch === '\n') {
      // 换行（暂不处理，后续可扩展）
      continue
    }

    var glyphIndex = font.charToGlyphIndex(ch)
    if (glyphIndex === 0) {
      // 字符不存在于字体中，跳过
      continue
    }

    var glyph = font.glyphs.get(glyphIndex)
    if (!glyph) continue

    var shape = glyphToShape(font, glyph, scale, offsetX, 0)
    shapes.push(shape)

    // 前进量
    var advanceWidth = glyph.advanceWidth || 0
    offsetX += advanceWidth * scale
  }

  return shapes
}

/**
 * 创建翠鸟风格 3D 文字 Mesh
 * @param {string} text - 文字内容
 * @param {opentype.Font} font - 已加载的字体
 * @param {object} options
 * @param {number} options.size - 文字大小 (默认2)
 * @param {number} options.depth - 挤出深度 (默认0.15)
 * @param {number} options.bevelThickness - 倒角厚度 (默认0.02)
 * @param {number} options.bevelSize - 倒角大小 (默认0.02)
 * @param {number} options.curveSegments - 曲线段数 (默认6)
 * @param {string|number} options.color - 材质颜色 (默认 '#8899aa')
 * @param {number} options.metalness - 金属度 (默认0.8)
 * @param {number} options.roughness - 粗糙度 (默认0.3)
 * @returns {THREE.Mesh|null}
 */
export function createText3DMesh(text, font, options) {
  if (!text || !font) return null

  options = options || {}
  var size = options.size || 2
  var depth = options.depth !== undefined ? options.depth : 0.15
  var bevelThickness = options.bevelThickness || 0.02
  var bevelSize = options.bevelSize || 0.02
  var curveSegments = options.curveSegments || 6
  var color = options.color || '#8899aa'
  var metalness = options.metalness !== undefined ? options.metalness : 0.8
  var roughness = options.roughness !== undefined ? options.roughness : 0.3

  // 1. 生成 SHAPE
  var shapes = textToShapes(text, font, size)
  if (shapes.length === 0) return null

  // 2. 创建 ExtrudeGeometry
  var geometry = new THREE.ExtrudeGeometry(shapes, {
    depth: depth,
    bevelEnabled: true,
    bevelThickness: bevelThickness,
    bevelSize: bevelSize,
    bevelOffset: 0,
    bevelSegments: 3,
    curveSegments: curveSegments,
    steps: 1
  })

  // 3. 居中几何体
  geometry.computeBoundingBox()
  if (geometry.boundingBox) {
    var centerX = -(geometry.boundingBox.max.x + geometry.boundingBox.min.x) / 2
    var centerZ = -(geometry.boundingBox.max.z + geometry.boundingBox.min.z) / 2
    geometry.translate(centerX, 0, centerZ - depth / 2)
  }

  // 4. 翠鸟风格金属材质
  var colorInt = typeof color === 'string' ? parseInt(color.replace('#', ''), 16) : color
  var material = new THREE.MeshStandardMaterial({
    color: colorInt,
    metalness: metalness,
    roughness: roughness,
    flatShading: false,
    envMapIntensity: 0.5
  })

  // 5. 创建 Mesh
  var mesh = new THREE.Mesh(geometry, material)
  mesh.castShadow = true
  mesh.receiveShadow = true
  mesh.userData.is3DText = true
  mesh.userData.textData = {
    text: text,
    size: size,
    depth: depth
  }
  mesh.name = '3DText_' + text

  return mesh
}

/**
 * 更新已有 3D 文字 Mesh 的内容
 * @param {THREE.Mesh} mesh
 * @param {string} text
 * @param {opentype.Font} font
 * @param {object} options
 */
export function updateText3DMesh(mesh, text, font, options) {
  if (!mesh || !mesh.userData.is3DText || !text || !font) return

  options = options || {}
  var size = options.size || mesh.userData.textData.size || 2
  var depth = options.depth !== undefined ? options.depth : (mesh.userData.textData.depth || 0.15)

  // 重新生成几何体
  var shapes = textToShapes(text, font, size)
  if (shapes.length === 0) return

  var oldGeo = mesh.geometry
  var geometry = new THREE.ExtrudeGeometry(shapes, {
    depth: depth,
    bevelEnabled: true,
    bevelThickness: options.bevelThickness || 0.02,
    bevelSize: options.bevelSize || 0.02,
    bevelOffset: 0,
    bevelSegments: 3,
    curveSegments: options.curveSegments || 6,
    steps: 1
  })

  geometry.computeBoundingBox()
  if (geometry.boundingBox) {
    var centerX = -(geometry.boundingBox.max.x + geometry.boundingBox.min.x) / 2
    var centerZ = -(geometry.boundingBox.max.z + geometry.boundingBox.min.z) / 2
    geometry.translate(centerX, 0, centerZ - depth / 2)
  }

  mesh.geometry = geometry
  if (oldGeo) oldGeo.dispose()

  mesh.userData.textData = {
    text: text,
    size: size,
    depth: depth
  }
}
