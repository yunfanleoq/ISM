/**
 * GisProjection — WGS84经纬度 ↔ 本地Three.js坐标 转换工具
 *
 * 用法：
 *   const proj = new GisProjection({ originLng: 120.5, originLat: 31.2 })
 *   const { x, z } = proj.lngLatToWorld(120.51, 31.21)
 *   const { lng, lat } = proj.worldToLngLat(x, z)
 */
import * as THREE from 'three'

const EARTH_RADIUS = 6378137 // WGS84 长半轴（米）
const DEG_TO_RAD = Math.PI / 180
const RAD_TO_DEG = 180 / Math.PI

export class GisProjection {
  /**
   * @param {Object} options
   * @param {number} options.originLng - 原点经度 (WGS84)
   * @param {number} options.originLat - 原点纬度 (WGS84)
   * @param {number} [options.scale=1]  - 缩放比例（用于大屏适配）
   */
  constructor(options = {}) {
    this.originLng = options.originLng || 0
    this.originLat = options.originLat || 0
    this.scale = options.scale || 1

    // 原点世界坐标（使用 Web Mercator 投影近似）
    this._originWorld = this._lngLatToWorld(this.originLng, this.originLat)
  }

  /**
   * 设置原点
   */
  setOrigin(lng, lat) {
    this.originLng = lng
    this.originLat = lat
    this._originWorld = this._lngLatToWorld(lng, lat)
  }

  /**
   * 设置缩放
   */
  setScale(s) {
    this.scale = s
  }

  // ---- 公开 API ----

  /**
   * WGS84经纬度 → Three.js 本地坐标系（y为高度轴，xz为水平面）
   * @returns {{ x: number, z: number }}
   */
  lngLatToWorld(lng, lat) {
    const world = this._lngLatToWorld(lng, lat)
    return {
      x: (world.x - this._originWorld.x) * this.scale,
      z: (world.z - this._originWorld.z) * this.scale
    }
  }

  /**
   * Three.js 本地坐标 → WGS84经纬度
   * @returns {{ lng: number, lat: number }}
   */
  worldToLngLat(x, z) {
    const worldX = x / this.scale + this._originWorld.x
    const worldZ = z / this.scale + this._originWorld.z
    return this._worldToLngLat(worldX, worldZ)
  }

  /**
   * 将目标点世界坐标转为相对于原点的 Three.js Vector3
   */
  lngLatToVector3(lng, lat, height = 0) {
    const { x, z } = this.lngLatToWorld(lng, lat)
    return new THREE.Vector3(x, height, z)
  }

  /**
   * 将 Three.js Vector3 转回经纬度
   */
  vector3ToLngLat(vec) {
    return this.worldToLngLat(vec.x, vec.z)
  }

  /**
   * 计算两个经纬度点之间的世界距离（米）
   */
  distance(lng1, lat1, lng2, lat2) {
    const p1 = this._lngLatToWorld(lng1, lat1)
    const p2 = this._lngLatToWorld(lng2, lat2)
    return Math.sqrt((p2.x-p1.x)**2 + (p2.z-p1.z)**2)
  }

  // ---- 内部方法 ----

  /**
   * WGS84 → Web Mercator 世界坐标
   */
  _lngLatToWorld(lng, lat) {
    const x = lng * DEG_TO_RAD * EARTH_RADIUS
    const y = Math.log(Math.tan((90 + lat) * DEG_TO_RAD / 2)) * EARTH_RADIUS
    return { x: x, z: -y } // z 轴镜像使 Three.js 坐标直观
  }

  /**
   * Web Mercator 世界坐标 → WGS84
   */
  _worldToLngLat(x, z) {
    const lng = (x / EARTH_RADIUS) * RAD_TO_DEG
    const lat = (Math.atan(Math.exp(-z / EARTH_RADIUS)) * 2 - Math.PI/2) * RAD_TO_DEG
    return { lng, lat }
  }
}

export default GisProjection
