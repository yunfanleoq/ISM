/**
 * SceneTemplates - 场景展示模板（翠鸟高端数字孪生风格）
 * 参照翠鸟智慧园区数字孪生大屏标准设计：
 * - 深色夜景背景 (#080e1a → #1a2a3a)
 * - 青蓝色发光建筑群 (emissive + high metalness)
 * - 地面发光网格 (gridFloor)
 * - 扫描线/流光/辉光球等视觉特效
 * - 数据标签面板 (text3d / dataText)
 * 
 * 所有 type 均来自 IndustrialObjects.js 真实存在的组件。
 */

const TPL = (id, type, name, typeName, x, y, z, color, opts = {}) => ({
  id, type, name, typeName, x, y, z,
  rx: opts.rx || 0, ry: opts.ry || 0, rz: opts.rz || 0,
  sx: opts.sx || 1, sy: opts.sy || 1, sz: opts.sz || 1,
  color: color || '#4a90d9',
  opacity: opts.opacity !== undefined ? opts.opacity : 1,
  metalness: opts.metalness !== undefined ? opts.metalness : 0.6,
  roughness: opts.roughness !== undefined ? opts.roughness : 0.3,
  wireframe: opts.wireframe || false,
  visible: true, locked: false,
  emissive: opts.emissive || '#000000',
  textureData: '',
  textContent: opts.textContent || '',
  fontSize: opts.fontSize || 48,
  textBgColor: opts.textBgColor || '#00000000',
  autoRotate: false, rotateSpeed: 1, rotateAxis: 'y',
  floatAnim: false, floatRange: 0.15, floatSpeed: 2,
  blink: false, blinkSpeed: 6, blinkMin: 0.2,
  animate: null, interactions: null,
  wsKey: '', bindProp: '', bindTransform: 'direct',
  bindScale: 1, bindOffset: 0, realTimeValue: '', dataFormat: '{value}',
  is2DComponent: false, label2D: '', source2D: null, source2DMeta: null,
  action: [],
  isExternalModel: false, isBackground: false, modelPath: '', materialOverridden: false
})

// 生成唯一ID
let _tplCounter = 0
const tid = (prefix) => `tpl_${prefix}_${++_tplCounter}`

// ===== 通用 sceneSettings：深色科技夜景 =====
const NIGHT_SCENE = (camPos, camTarget) => ({
  showGrid: false,
  backgroundMode: 'gradient',
  backgroundColor: '#080e1a',
  backgroundColor2: '#1a2a3a',
  environmentPreset: 'night',
  lightingPreset: 'night',
  floorReflection: 'soft',
  enhanceDepth: true,
  gridSize: 20,
  gridDivisions: 30,
  gridColorCenterLine: '#13c2c244',
  gridColorGrid: '#13c2c222',
  cameraPosition: camPos || { x: 5, y: 4, z: 6 },
  cameraTarget: camTarget || { x: 0, y: 0.5, z: 0 },
  cameraFov: 38,
  lightSettings: {
    envIntensity: 1.6,
    ambientColor: '#b0d8e8',
    ambientIntensity: 1.1,
    directionalColor: '#ffffff',
    directionalIntensity: 1.2
  }
})

// ===== 辅助函数：创建发光建筑 =====
const GLOW_BUILDING = (prefix, id, name, typeName, x, y, z, color, h, w, d, glowColor) =>
  TPL(tid(prefix), id, name, typeName, x, y, z, color, {
    sx: w || 1, sy: h || 1, sz: d || 1,
    emissive: glowColor || '#081828',
    metalness: 0.7, roughness: 0.2
  })

// ===== 辅助函数：创建数据标签 =====
const DATA_LABEL = (prefix, text, x, y, z, color, size) =>
  TPL(tid(prefix), 'text3d', text, 'Text3D', x, y, z, color || '#13c2c2', {
    textContent: text,
    fontSize: size || 24,
    textBgColor: '#080e1acc',
    sx: 0.8, sy: 0.8
  })

export const SCENE_TEMPLATES = [
  // ==================== 1. 智慧园区总览（翠鸟高端数字孪生旗舰模板）====================
  {
    id: 'smart-park-overview',
    name: '智慧园区总览',
    icon: 'fas fa-city',
    category: '智慧城市',
    description: '翠鸟高端数字孪生：玻璃幕墙发光建筑群+地面扫描环+流光数据线+悬浮数据面板+全息投影',
    sceneSettings: NIGHT_SCENE({ x: 12, y: 10, z: 12 }, { x: 0, y: 3.0, z: 0 }),
    objects: [
      // ========== 地面 ==========
      TPL(tid('pk'), 'dtGroundGrid', '地面发光网格', 'DTGroundGrid', 0, 0, 0, '#13c2c2', {}),
      // 中心扫描环（动画展开）
      TPL(tid('pk'), 'dtScanRing', '中心扫描环', 'DTScanRing', 0, 0.03, 0, '#13c2c2', {}),
      // 氛围粒子
      TPL(tid('pk'), 'dtParticleField', '氛围粒子场', 'DTParticleField', 0, 3, 0, '#13c2c2', {}),

      // ========== 中心建筑群（玻璃幕墙发光） ==========
      // 主楼（高层）
      TPL(tid('pk'), 'dtBuildingTall', '总部大楼A', 'DTBuildingTall', 0, 0, -4, '#0d1f2a', { sx: 1.5, sz: 1.5 }),
      // 副楼（标准）
      TPL(tid('pk'), 'dtBuilding', '研发大楼B', 'DTBuilding', -4, 0, -2, '#0d1f2a', { sx: 1.2, sz: 1.2 }),
      TPL(tid('pk'), 'dtBuilding', '行政大楼C', 'DTBuilding', 4, 0, -2, '#0d1f2a', { sx: 1.2, sz: 1.2 }),
      // 数据中心（扁平厂房）
      TPL(tid('pk'), 'dtBuildingWide', '数据中心', 'DTBuildingWide', 0, 0, 4, '#0a1a2e', { sx: 1.0, sz: 1.0 }),
      // 公寓
      TPL(tid('pk'), 'dtBuilding', '人才公寓D', 'DTBuilding', -3, 0, 4, '#1a0d2a', { sx: 0.9, sz: 0.9, sy: 0.7 }),
      TPL(tid('pk'), 'dtBuilding', '员工宿舍E', 'DTBuilding', 3, 0, 4, '#1a0d2a', { sx: 0.9, sz: 0.9, sy: 0.7 }),
      // 复合建筑（主楼+裙楼）
      TPL(tid('pk'), 'dtBuildingComplex', '综合商务楼', 'DTBuildingComplex', -6, 0, 4, '#0d1a2e', { sx: 1.0, sz: 1.0 }),

      // ========== 道路 ==========
      TPL(tid('pk'), 'dtRoad', '主干道-南北', 'DTRoad', 0, 0, 0, '#1a1a2e', { sx: 1, sz: 3, ry: 0 }),
      TPL(tid('pk'), 'dtRoad', '主干道-东西', 'DTRoad', 0, 0, 0, '#1a1a2e', { sx: 3, sz: 1, ry: 0 }),
      TPL(tid('pk'), 'dtRoad', '园区环路-内', 'DTRoad', 0, 0, 0, '#1a1a2e', { sx: 2.5, sz: 2.5, ry: 0 }),

      // ========== 智慧路灯 ==========
      TPL(tid('pk'), 'dtStreetLightDT', '智慧路灯1', 'DTStreetLight', -5, 0, -5, '#444444', {}),
      TPL(tid('pk'), 'dtStreetLightDT', '智慧路灯2', 'DTStreetLight', 5, 0, -5, '#444444', {}),
      TPL(tid('pk'), 'dtStreetLightDT', '智慧路灯3', 'DTStreetLight', -5, 0, 5, '#444444', {}),
      TPL(tid('pk'), 'dtStreetLightDT', '智慧路灯4', 'DTStreetLight', 5, 0, 5, '#444444', {}),

      // ========== 通信基站 ==========
      TPL(tid('pk'), 'dtBaseStation', '5G基站A', 'DTBaseStation', -8, 0, -8, '#666666', {}),
      TPL(tid('pk'), 'dtBaseStation', '5G基站B', 'DTBaseStation', 8, 0, -8, '#666666', {}),

      // ========== 绿地 ==========
      TPL(tid('pk'), 'dtPark', '智慧绿地A', 'DTPark', -8, 0, 6, '#0a331a', {}),
      TPL(tid('pk'), 'dtPark', '智慧绿地B', 'DTPark', 8, 0, 6, '#0a331a', {}),

      // ========== 扫描光束（垂直） ==========
      TPL(tid('pk'), 'dtScanLine', '扫描光束A', 'DTScanLine', -3, 0, 0, '#13c2c2', {}),
      TPL(tid('pk'), 'dtScanLine', '扫描光束B', 'DTScanLine', 3, 0, 0, '#13c2c2', {}),

      // ========== 流光数据线 ==========
      TPL(tid('pk'), 'dtFlowLine', '数据流-东西', 'DTFlowLine', -6, 0.5, 0, '#00e5ff', {}),
      TPL(tid('pk'), 'dtFlowLine', '数据流-南北', 'DTFlowLine', 0, 0.5, -6, '#00e5ff', {}),

      // ========== 全息投影（中心） ==========
      TPL(tid('pk'), 'dtHologram', '全息指挥台', 'DTHologram', 0, 0, 0, '#13c2c2', {}),

      // ========== 悬浮数据面板（左右侧） ==========
      TPL(tid('pk'), 'dtDataPanel', 'KPI面板-左', 'DTDataPanel', -10, 4, 0, '#0a1628', { sx: 1.0 }),
      TPL(tid('pk'), 'dtDataPanel', 'KPI面板-右', 'DTDataPanel', 10, 4, 0, '#0a1628', { sx: 1.0 }),
      TPL(tid('pk'), 'dtDataPanel', 'KPI面板-左上', 'DTDataPanel', -10, 6, -3, '#0a1628', { sx: 0.8 }),
      TPL(tid('pk'), 'dtDataPanel', 'KPI面板-右上', 'DTDataPanel', 10, 6, -3, '#0a1628', { sx: 0.8 }),

      // ========== 数据标签（大标题） ==========
      DATA_LABEL('pk', '深国際 · 智慧园区数字孪生平台', 0, 8, -8, '#13c2c2', 48),
      DATA_LABEL('pk', '服务企业: 3', -9, 4.5, 2, '#4dfff0', 28),
      DATA_LABEL('pk', '园区面积: 108,279 m²', -9, 4.0, 2, '#4dfff0', 24),
      DATA_LABEL('pk', '今日人流: 1,247', -9, 3.5, 2, '#4dffa6', 24),
      DATA_LABEL('pk', '在线设备: 3,073', 9, 4.5, 2, '#4dffa6', 28),
      DATA_LABEL('pk', '告警: 0', 9, 4.0, 2, '#ff4d4f', 24),
      DATA_LABEL('pk', '离线: 164', 9, 3.5, 2, '#ffaa4d', 24),
    ]
  },

  // ==================== 2. 台站总览园区 ====================
  {
    id: 'station-overview-yard',
    name: '台站总览园区',
    icon: 'fas fa-broadcast-tower',
    category: '工业场景',
    description: '深色科技风台站总览，冷蓝灯光+流光管线+悬浮数据牌',
    sceneSettings: NIGHT_SCENE({ x: 5, y: 4, z: 6 }, { x: 0, y: 0.5, z: 0 }),
    objects: [
      // 地面
      TPL(tid('stn'), 'gridFloor', '地面网格', 'GridFloor', 0, 0.01, 0, '#13c2c2', { sx: 2, sz: 2 }),
      TPL(tid('stn'), 'plane', '园区底盘', 'Plane', 0, 0, 0, '#0d1f2a', { sx: 4, sz: 3, emissive: '#040a10' }),
      // 建筑
      GLOW_BUILDING('stn', 'building', '综合业务楼', 'Building', -0.5, 0, -0.5, '#1a3a4a', 1.5, 1.5, 1.0, '#0a1a2a'),
      GLOW_BUILDING('stn', 'factorybld', '运维中心', 'FactoryBld', 2.5, 0, -0.5, '#1a2a3a', 1.0, 1.0, 0.8, '#081828'),
      // 通信塔
      TPL(tid('stn'), 'telecomtower', '北塔', 'TelecomTower', -3, 0, -2, '#e0f4ff', { sy: 1.5, emissive: '#081830' }),
      TPL(tid('stn'), 'telecomtower', '南塔', 'TelecomTower', 2, 0, -2, '#e0f4ff', { sy: 1.3, emissive: '#081830' }),
      // 控制设备
      TPL(tid('stn'), 'controlpanel', '控制面板', 'ControlPanel', 1.5, 0, 0.5, '#1a2a3a', { sy: 0.8, emissive: '#050f18' }),
      TPL(tid('stn'), 'plc', 'PLC控制器', 'PLC', 1.8, 0, 0.8, '#ffaa4d', { emissive: '#ffaa4d22' }),
      // 摄像头
      TPL(tid('stn'), 'surveilcam', '云台摄像机', 'SurveilCam', -1.5, 0, 1, '#13c2c2', { sy: 1.2 }),
      // 树木
      TPL(tid('stn'), 'citytree', '前景树A', 'CityTree', -2.5, 0, 1.5, '#0a2a1a', { sx: 0.6, sy: 0.7, sz: 0.6 }),
      TPL(tid('stn'), 'citytree', '前景树B', 'CityTree', 2.5, 0, 1.5, '#0a2a1a', { sx: 0.6, sy: 0.7, sz: 0.6 }),
      // 道路
      TPL(tid('stn'), 'road', '入口道路', 'Road', 0, 0.02, 2.5, '#0d1f2a', { ry: Math.PI/2, sx: 0.3, sz: 3 }),
      // 特效
      TPL(tid('stn'), 'scanLine', '扫描线', 'ScanLine', 0, 0.5, 0, '#13c2c2', { sy: 1.2 }),
      TPL(tid('stn'), 'flyLine', '流光A', 'FlyLine', -1, 0.2, 2.2, '#13c2c2', { sx: 1.2 }),
      TPL(tid('stn'), 'flyLine', '流光B', 'FlyLine', 1, 0.2, 2.2, '#36cfc9', { sx: 1.2 }),
      TPL(tid('stn'), 'glowSphere', '辉光球', 'GlowSphere', 0, 1.5, 0, '#13c2c2', { sx: 0.3, sy: 0.3, sz: 0.3 }),
      // 围栏
      TPL(tid('stn'), 'fence', '围栏N', 'Fence', 0, 0, -2, '#1a2a3a', { sx: 3, emissive: '#081020' }),
      TPL(tid('stn'), 'fence', '围栏S', 'Fence', 0, 0, 2, '#1a2a3a', { sx: 3, emissive: '#081020' }),
      // 数据标签
      DATA_LABEL('stn', '台站总览', 0, 3.5, -2.5, '#13c2c2', 34),
      DATA_LABEL('stn', '北塔在线', -2.5, 1, -0.5, '#4dffa6', 20),
      DATA_LABEL('stn', '南塔在线', 2, 1, -0.5, '#4dffa6', 20),
      DATA_LABEL('stn', '运行正常', 0, 0.5, 2.8, '#13c2c2', 22),
      // 车辆
      TPL(tid('stn'), 'car', '巡检车', 'Car', 1.5, 0, 2.5, '#1a3a4a', { ry: Math.PI/2, sx: 0.6, sy: 0.6, sz: 0.6 }),
    ]
  },

  // ==================== 3. 冷却机组直播间 ====================
  {
    id: 'cooling-room-live',
    name: '冷却机组直播间',
    icon: 'fas fa-snowflake',
    category: '工业场景',
    description: '深色工业风冷却机组监控场景，冷蓝主调+流光管线',
    sceneSettings: NIGHT_SCENE({ x: 5, y: 3.5, z: 5 }, { x: 0, y: 0.5, z: 0 }),
    objects: [
      TPL(tid('clr'), 'gridFloor', '地面网格', 'GridFloor', 0, 0.01, 0, '#13c2c2', { sx: 2, sz: 1.5 }),
      TPL(tid('clr'), 'plane', '设备地台', 'Plane', 0, 0, 0, '#0d1a22', { sx: 3, sz: 2, emissive: '#040810' }),
      // 冷却机组
      TPL(tid('clr'), 'coolingFanUnit', '双风机冷却机组', 'CoolingFanUnit', 0, 0, -0.5, '#0a2a1a', { sx: 1.5, sy: 1.2, sz: 1.0, emissive: '#051520' }),
      // 控制柜
      TPL(tid('clr'), 'controlpanel', 'PLC控制柜', 'ControlPanel', -2.5, 0, -0.5, '#1a1a2a', { sx: 1.0, sy: 1.2, emissive: '#050a10' }),
      TPL(tid('clr'), 'hmipanel', 'HMI面板', 'HMIPanel', -2.0, 0, 0.5, '#1a2a3a', { emissive: '#050a10' }),
      // 管道
      TPL(tid('clr'), 'pipe', '供水管', 'Pipe', 0, 0.3, 1.2, '#13c2c244', { ry: Math.PI/2, sx: 4, sy: 0.3, sz: 0.3, emissive: '#13c2c222' }),
      TPL(tid('clr'), 'pipe', '回水管', 'Pipe', 0, 0.35, -1.2, '#36cfc944', { ry: Math.PI/2, sx: 4, sy: 0.25, sz: 0.25, emissive: '#13c2c222' }),
      // 泵
      TPL(tid('clr'), 'centrifugalpump', '循环泵A', 'CentrifugalPump', -1.2, 0, 0.8, '#0a2a1a', { sx: 0.7, sy: 0.7, sz: 0.7, emissive: '#13c2c222' }),
      TPL(tid('clr'), 'centrifugalpump', '循环泵B', 'CentrifugalPump', 0, 0, 0.8, '#0a2a1a', { sx: 0.7, sy: 0.7, sz: 0.7, emissive: '#13c2c222' }),
      TPL(tid('clr'), 'centrifugalpump', '循环泵C', 'CentrifugalPump', 1.2, 0, 0.8, '#0a2a1a', { sx: 0.7, sy: 0.7, sz: 0.7, emissive: '#13c2c222' }),
      // 阀门
      TPL(tid('clr'), 'valve', '阀门A', 'Valve', -1.8, 0, 1.2, '#13c2c2', { sx: 0.6, sy: 0.6, sz: 0.6 }),
      TPL(tid('clr'), 'valve', '阀门B', 'Valve', 1.8, 0, 1.2, '#13c2c2', { sx: 0.6, sy: 0.6, sz: 0.6 }),
      // 传感器
      TPL(tid('clr'), 'sensorarray', '温湿度传感器', 'SensorArray', -2.2, 0, -1.0, '#13c2c2', { sy: 1.2 }),
      TPL(tid('clr'), 'cctv', '直播摄像头', 'CCTV', 2.5, 0, -1.0, '#13c2c2', { sy: 1.3 }),
      // 特效
      TPL(tid('clr'), 'scanLine', '扫描线', 'ScanLine', 0, 0.5, 0, '#13c2c2', { sy: 1 }),
      TPL(tid('clr'), 'flyLine', '供水流光', 'FlyLine', -1.5, 0.4, 1.5, '#13c2c2', { sx: 1.2 }),
      TPL(tid('clr'), 'flyLine', '回水流光', 'FlyLine', 1.5, 0.45, 1.5, '#36cfc9', { sx: 1.2 }),
      TPL(tid('clr'), 'glowSphere', '中心辉光', 'GlowSphere', 0, 1, 0, '#13c2c2', { sx: 0.25, sy: 0.25, sz: 0.25 }),
      // 数据标签
      DATA_LABEL('clr', '冷却机组直播间', 0, 3, -2, '#13c2c2', 30),
      DATA_LABEL('clr', '温度: 24.5C', -3, 1, 1.5, '#4dfff0', 20),
      DATA_LABEL('clr', '压力: 0.8MPa', 3, 1, 1.5, '#4dfff0', 20),
      DATA_LABEL('clr', '点击进入直播间', 0, 0.5, 2.5, '#13c2c2', 22),
    ]
  },

  // ==================== 4. 智慧城市街道 ====================
  {
    id: 'smart-city',
    name: '智慧城市街道',
    icon: 'fas fa-city',
    category: '智慧城市',
    description: '夜景智慧城市街道，冷蓝灯光+霓虹道路+发光建筑窗户',
    sceneSettings: NIGHT_SCENE({ x: 6, y: 5, z: 6 }, { x: 0, y: 1.0, z: 0 }),
    objects: [
      // 地面
      TPL(tid('city'), 'gridFloor', '地面网格', 'GridFloor', 0, 0.01, 0, '#13c2c2', { sx: 3, sz: 3 }),
      // 道路
      TPL(tid('city'), 'road', '主干道X', 'Road', 0, 0.02, -3.5, '#1a1a2a', { ry: Math.PI/2, sx: 0.4, sz: 8 }),
      TPL(tid('city'), 'road', '主干道Z', 'Road', -3.5, 0.02, 0, '#1a1a2a', { sx: 0.4, sz: 8 }),
      TPL(tid('city'), 'zebracrossing', '斑马线', 'ZebraCrossing', 0, 0.03, -2.5, '#2a2a3a', { sx: 2 }),
      // 建筑
      GLOW_BUILDING('city', 'building', '办公楼A', 'Building', -3, 0, -3, '#0d1f2a', 2.5, 1.2, 1.0, '#081020'),
      GLOW_BUILDING('city', 'building', '办公楼B', 'Building', 3, 0, -3, '#1a0d2a', 2.0, 1.0, 1.2, '#100818'),
      GLOW_BUILDING('city', 'apartment', '住宅楼A', 'Apartment', -3, 0, 3, '#0d2a1a', 3.0, 0.8, 0.8, '#081820'),
      GLOW_BUILDING('city', 'apartment', '住宅楼B', 'Apartment', 3, 0, 3, '#1a2a0d', 2.2, 0.9, 0.9, '#101808'),
      // 发光窗户
      TPL(tid('city'), 'box', '窗户A1', 'Box', -3.3, 1.5, -2.9, '#13c2c2', { sx: 0.12, sy: 0.2, sz: 0.05, emissive: '#13c2c288' }),
      TPL(tid('city'), 'box', '窗户A2', 'Box', -2.7, 1.5, -2.9, '#13c2c2', { sx: 0.12, sy: 0.2, sz: 0.05, emissive: '#13c2c288' }),
      TPL(tid('city'), 'box', '窗户B1', 'Box', 2.7, 1.2, -2.9, '#36cfc9', { sx: 0.12, sy: 0.2, sz: 0.05, emissive: '#36cfc988' }),
      TPL(tid('city'), 'box', '窗户B2', 'Box', 3.3, 1.2, -2.9, '#36cfc9', { sx: 0.12, sy: 0.2, sz: 0.05, emissive: '#36cfc988' }),
      // 路灯
      TPL(tid('city'), 'streetlight', '路灯1', 'StreetLight', -1.5, 0, -3.8, '#13c2c2', { emissive: '#13c2c244' }),
      TPL(tid('city'), 'streetlight', '路灯2', 'StreetLight', 1.5, 0, -3.8, '#13c2c2', { emissive: '#13c2c244' }),
      TPL(tid('city'), 'streetlight', '路灯3', 'StreetLight', -3.8, 0, -1.5, '#36cfc9', { emissive: '#36cfc944' }),
      TPL(tid('city'), 'streetlight', '路灯4', 'StreetLight', -3.8, 0, 1.5, '#36cfc9', { emissive: '#36cfc944' }),
      // 信号灯
      TPL(tid('city'), 'trafficlight', '信号灯A', 'TrafficLight', -0.3, 0, -3.2, '#ff4d4f', { emissive: '#ff4d4f33' }),
      TPL(tid('city'), 'trafficlight', '信号灯B', 'TrafficLight', -3.2, 0, 0.3, '#ff4d4f', { emissive: '#ff4d4f33' }),
      // 车辆
      TPL(tid('city'), 'car', '小汽车A', 'Car', 0.5, 0, -2, '#13c2c244', { ry: Math.PI/2 }),
      TPL(tid('city'), 'bus', '公交车', 'Bus', -2, 0, -2, '#ffaa4d44', { rx: 0, ry: Math.PI/2 }),
      TPL(tid('city'), 'suv', 'SUV', 'SUV', 2, 0, -1, '#4dffa644', { ry: Math.PI/2 }),
      // 绿化
      TPL(tid('city'), 'citytree', '行道树1', 'CityTree', -2.5, 0, -4.3, '#0a1a0a', { sx: 0.7, sy: 0.8, sz: 0.7 }),
      TPL(tid('city'), 'citytree', '行道树2', 'CityTree', 2.5, 0, -4.3, '#0a1a0a', { sx: 0.7, sy: 0.8, sz: 0.7 }),
      TPL(tid('city'), 'bush', '灌木A', 'Bush', -1, 0, -4.0, '#0a1a0a', { sx: 0.8, sz: 0.5 }),
      // 设施
      TPL(tid('city'), 'bench', '长椅', 'Bench', -1.5, 0, -4.0, '#1a100a', { sx: 0.8 }),
      TPL(tid('city'), 'billboard', '广告牌', 'Billboard', 0, 0, -3.5, '#13c2c2', { ry: Math.PI/2, sy: 0.8, emissive: '#13c2c222' }),
      TPL(tid('city'), 'flagpole', '旗杆', 'Flagpole', 0, 0, -1, '#555555', { sy: 1.2 }),
      TPL(tid('city'), 'surveilcam', '摄像头', 'SurveilCam', -1, 0, -3.2, '#13c2c2', { sy: 1.5 }),
      // 特效
      TPL(tid('city'), 'scanLine', '扫描线', 'ScanLine', 0, 0.5, 0, '#13c2c2', { sy: 1.2 }),
      TPL(tid('city'), 'ringScan', '雷达扫描', 'RingScan', 0, 0.1, 0, '#13c2c2', { sx: 1.5, sz: 1.5 }),
      TPL(tid('city'), 'flyLine', '道路流光A', 'FlyLine', -0.5, 0.1, -3.5, '#13c2c2', { sx: 2.0 }),
      TPL(tid('city'), 'flyLine', '道路流光B', 'FlyLine', -3.5, 0.1, 0.5, '#36cfc9', { sx: 2.0 }),
      TPL(tid('city'), 'glowSphere', '辉光球', 'GlowSphere', 0, 2, 0, '#13c2c2', { sx: 0.3, sy: 0.3, sz: 0.3 }),
      // 数据标签
      DATA_LABEL('city', '智慧城市街道', 0, 4.5, -3.5, '#13c2c2', 32),
      DATA_LABEL('city', '实时车流: 正常', -4, 2, 0, '#4dffa6', 20),
      DATA_LABEL('city', '空气质量: 优', 4, 2, 0, '#4dffa6', 20),
    ]
  },

  // ==================== 5. 智慧工厂流水线 ====================
  {
    id: 'factory-line',
    name: '智慧工厂流水线',
    icon: 'fas fa-industry',
    category: '工业流水线',
    description: '深色工业风智能产线，发光机械臂+流光传送带+数据监控面板',
    sceneSettings: NIGHT_SCENE({ x: 7, y: 4, z: 5 }, { x: 0, y: 1.0, z: 0 }),
    objects: [
      TPL(tid('fact'), 'gridFloor', '地面网格', 'GridFloor', 0, 0.01, 0, '#13c2c2', { sx: 2.5, sz: 2 }),
      TPL(tid('fact'), 'plane', '厂房底盘', 'Plane', 0, 0, 0, '#0d1a22', { sx: 5, sz: 4, emissive: '#040810' }),
      // 厂房
      GLOW_BUILDING('fact', 'factorybld', '主厂房', 'FactoryBld', 0, 0, 0, '#0d1a22', 1.5, 2.5, 2.5, '#040810'),
      // 传送带
      TPL(tid('fact'), 'conveyorbelt', '传送带1', 'ConveyorBelt', -1.5, 0.2, 0, '#1a1a2a', { sx: 2.5, emissive: '#080e16' }),
      TPL(tid('fact'), 'conveyorbelt', '传送带2', 'ConveyorBelt', 1.5, 0.2, 0, '#1a1a2a', { sx: 2.5, emissive: '#080e16' }),
      TPL(tid('fact'), 'conveyorbelt', '传送带3', 'ConveyorBelt', 0, 0.2, -1.5, '#1a1a2a', { ry: Math.PI/2, sx: 1.5, emissive: '#080e16' }),
      // 机械臂
      TPL(tid('fact'), 'robotarm', '机械臂A', 'RobotArm', -1.5, 0, 1.2, '#13c2c2', { emissive: '#13c2c222', metalness: 0.8 }),
      TPL(tid('fact'), 'robotarm', '机械臂B', 'RobotArm', 1.5, 0, 1.2, '#36cfc9', { emissive: '#36cfc922', metalness: 0.8 }),
      TPL(tid('fact'), 'robotarm', '机械臂C', 'RobotArm', 0, 0, -2.2, '#13c2c2', { emissive: '#13c2c222', metalness: 0.8 }),
      TPL(tid('fact'), 'weldingrobot', '焊接机器人', 'WeldingRobot', 0, 0, 1.5, '#ffaa4d', { emissive: '#ffaa4d22' }),
      // CNC
      TPL(tid('fact'), 'cncmachine', 'CNC机床A', 'CNCMachine', -2, 0, -1.5, '#0d1f2a', { emissive: '#081018', metalness: 0.6 }),
      TPL(tid('fact'), 'cncmachine', 'CNC机床B', 'CNCMachine', 2, 0, -1.5, '#0d1f2a', { emissive: '#081018', metalness: 0.6 }),
      // 装配
      TPL(tid('fact'), 'assemblystation', '装配站', 'AssemblyStation', 0, 0, 0, '#0d1f2a', { emissive: '#081018' }),
      TPL(tid('fact'), 'workbench', '工作台A', 'Workbench', -0.8, 0, -0.5, '#1a0d00', { sx: 0.8 }),
      TPL(tid('fact'), 'workbench', '工作台B', 'Workbench', 0.8, 0, -0.5, '#1a0d00', { sx: 0.8 }),
      // 包装质检
      TPL(tid('fact'), 'packmachine', '包装机', 'PackMachine', 2, 0, 1.5, '#0d1a0d', { emissive: '#081808' }),
      TPL(tid('fact'), 'qualitygate', '质检门', 'QualityGate', 1.5, 0, 2.2, '#13c2c2', { emissive: '#13c2c222' }),
      TPL(tid('fact'), 'inspectioncam', '检测相机', 'InspectionCam', 1.5, 0, 1.8, '#13c2c2', { sy: 1.5 }),
      // AGV
      TPL(tid('fact'), 'agv', 'AGV小车A', 'AGV', -2.2, 0, 1.8, '#13c2c244', { sx: 0.8 }),
      TPL(tid('fact'), 'agv', 'AGV小车B', 'AGV', 2.2, 0, 1.8, '#36cfc944', { sx: 0.8 }),
      // 安全围栏
      TPL(tid('fact'), 'safetyfence', '安全围栏L', 'SafetyFence', -2, 0, 1.2, '#1a1a0d', { ry: Math.PI/2, sx: 2.5, emissive: '#080800' }),
      TPL(tid('fact'), 'safetyfence', '安全围栏R', 'SafetyFence', 2, 0, 1.2, '#1a1a0d', { ry: Math.PI/2, sx: 2.5, emissive: '#080800' }),
      // 特效
      TPL(tid('fact'), 'scanLine', '扫描线', 'ScanLine', 0, 0.5, 0, '#13c2c2', { sy: 1.2 }),
      TPL(tid('fact'), 'flyLine', '传送流光A', 'FlyLine', -1.5, 0.35, 0, '#13c2c2', { sx: 1.5 }),
      TPL(tid('fact'), 'flyLine', '传送流光B', 'FlyLine', 1.5, 0.35, 0, '#36cfc9', { sx: 1.5 }),
      TPL(tid('fact'), 'glowSphere', '辉光球', 'GlowSphere', 0, 1.5, 0, '#13c2c2', { sx: 0.25, sy: 0.25, sz: 0.25 }),
      // 数据大屏
      TPL(tid('fact'), 'dashboardscreen', '设备状态大屏', 'DashboardScreen', 0, 0, -2.2, '#0d1f2a', { sx: 1.1, sy: 0.8 }),
      // 数据标签
      DATA_LABEL('fact', '智慧工厂', 0, 3.5, -2.5, '#13c2c2', 32),
      DATA_LABEL('fact', '产线效率: 97%', -3, 1.5, 2, '#4dffa6', 20),
      DATA_LABEL('fact', '设备状态: 正常', 3, 1.5, 2, '#4dffa6', 20),
    ]
  },

  // ==================== 6. 数据中心机房 ====================
  {
    id: 'data-center',
    name: '数据中心机房',
    icon: 'fas fa-server',
    category: '数字孪生',
    description: '冷蓝色调数据中心，发光机柜+光纤流光+监控大屏',
    sceneSettings: NIGHT_SCENE({ x: 7, y: 3.5, z: 5 }, { x: 0, y: 1.0, z: 0.5 }),
    objects: [
      TPL(tid('dc'), 'gridFloor', '地面网格', 'GridFloor', 0, 0.01, 0, '#13c2c2', { sx: 3, sz: 2.5 }),
      TPL(tid('dc'), 'plane', '机房地板', 'Plane', 0, 0, 0, '#080e16', { sx: 6, sz: 5, emissive: '#030508' }),
      // 机柜列A
      TPL(tid('dc'), 'serverrack', '机柜A1', 'ServerRack', -2.5, 0, -2, '#0d1f2a', { emissive: '#081018', metalness: 0.5 }),
      TPL(tid('dc'), 'serverrack', '机柜A2', 'ServerRack', -1.5, 0, -2, '#0d1f2a', { emissive: '#081018', metalness: 0.5 }),
      TPL(tid('dc'), 'serverrack', '机柜A3', 'ServerRack', -0.5, 0, -2, '#0d1f2a', { emissive: '#081018', metalness: 0.5 }),
      TPL(tid('dc'), 'serverrack', '机柜A4', 'ServerRack', 0.5, 0, -2, '#0d1f2a', { emissive: '#081018', metalness: 0.5 }),
      TPL(tid('dc'), 'serverrack', '机柜A5', 'ServerRack', 1.5, 0, -2, '#0d1f2a', { emissive: '#081018', metalness: 0.5 }),
      TPL(tid('dc'), 'serverrack', '机柜A6', 'ServerRack', 2.5, 0, -2, '#0d1f2a', { emissive: '#081018', metalness: 0.5 }),
      // 机柜列B
      TPL(tid('dc'), 'serverrack', '机柜B1', 'ServerRack', -2.5, 0, 1, '#0d1f2a', { emissive: '#081018', metalness: 0.5 }),
      TPL(tid('dc'), 'serverrack', '机柜B2', 'ServerRack', -1.5, 0, 1, '#0d1f2a', { emissive: '#081018', metalness: 0.5 }),
      TPL(tid('dc'), 'serverrack', '机柜B3', 'ServerRack', -0.5, 0, 1, '#0d1f2a', { emissive: '#081018', metalness: 0.5 }),
      TPL(tid('dc'), 'serverrack', '机柜B4', 'ServerRack', 0.5, 0, 1, '#0d1f2a', { emissive: '#081018', metalness: 0.5 }),
      TPL(tid('dc'), 'serverrack', '机柜B5', 'ServerRack', 1.5, 0, 1, '#0d1f2a', { emissive: '#081018', metalness: 0.5 }),
      TPL(tid('dc'), 'serverrack', '机柜B6', 'ServerRack', 2.5, 0, 1, '#0d1f2a', { emissive: '#081018', metalness: 0.5 }),
      // 网络设备
      TPL(tid('dc'), 'router', '核心路由器', 'Router', -1, 0, 2.2, '#13c2c2', { emissive: '#13c2c222' }),
      TPL(tid('dc'), 'switchgear', '核心交换机', 'Switchgear', 0, 0, 2.2, '#0d1f2a', { emissive: '#081018' }),
      TPL(tid('dc'), 'fiberbox', '光纤盒', 'FiberBox', -1.5, 0, 2.2, '#36cfc9', { emissive: '#36cfc922' }),
      // UPS
      TPL(tid('dc'), 'ups', 'UPS电源A', 'UPS', -2.8, 0, 2.2, '#0a2a1a', { emissive: '#051510' }),
      TPL(tid('dc'), 'ups', 'UPS电源B', 'UPS', 2.8, 0, 2.2, '#0a2a1a', { emissive: '#051510' }),
      // 监控
      TPL(tid('dc'), 'dashboardscreen', '监控大屏', 'DashboardScreen', 0, 0, 2.0, '#0d1f2a', { sy: 1.2 }),
      TPL(tid('dc'), 'cctv', '监控球机A', 'CCTV', -2.5, 0, 0, '#13c2c2', { sy: 2 }),
      TPL(tid('dc'), 'cctv', '监控球机B', 'CCTV', 2.5, 0, 0, '#13c2c2', { sy: 2 }),
      // 传感器
      TPL(tid('dc'), 'sensorarray', '温湿度传感器', 'SensorArray', -2, 0, -0.5, '#13c2c2', { sy: 1.5 }),
      // 特效
      TPL(tid('dc'), 'scanLine', '扫描线', 'ScanLine', 0, 0.5, 0, '#13c2c2', { sy: 1.5 }),
      TPL(tid('dc'), 'flyLine', '光纤流光A', 'FlyLine', -0.5, 0.15, 2.0, '#13c2c2', { sx: 1.2 }),
      TPL(tid('dc'), 'flyLine', '光纤流光B', 'FlyLine', 0.5, 0.15, 2.0, '#36cfc9', { sx: 1.2 }),
      TPL(tid('dc'), 'glowSphere', '辉光球', 'GlowSphere', 0, 1.5, 0, '#13c2c2', { sx: 0.2, sy: 0.2, sz: 0.2 }),
      TPL(tid('dc'), 'particleCloud', '冷气粒子', 'ParticleCloud', 0, 0.5, 0, '#13c2c211', { sx: 3, sy: 1, sz: 3 }),
      // 数据标签
      DATA_LABEL('dc', '数据中心机房', 0, 3, -2.5, '#13c2c2', 32),
      DATA_LABEL('dc', 'CPU: 45%', -4, 1.5, 2, '#4dffa6', 20),
      DATA_LABEL('dc', '内存: 62%', 4, 1.5, 2, '#ffaa4d', 20),
      DATA_LABEL('dc', '温度: 22C', 0, 0.5, 2.8, '#4dfff0', 20),
    ]
  },

  // ==================== 7. 能源变电站 ====================
  {
    id: 'energy-station',
    name: '能源变电站',
    icon: 'fas fa-bolt',
    category: '能源电力',
    description: '高压工业风变电站，橙色变压器+蓝色电弧光效+光纤通信',
    sceneSettings: NIGHT_SCENE({ x: 6, y: 4, z: 6 }, { x: 0, y: 1.0, z: 0 }),
    objects: [
      TPL(tid('en'), 'gridFloor', '地面网格', 'GridFloor', 0, 0.01, 0, '#13c2c2', { sx: 2.5, sz: 2.5 }),
      TPL(tid('en'), 'plane', '变电站底盘', 'Plane', 0, 0, 0, '#0d1a16', { sx: 5, sz: 5, emissive: '#040a08' }),
      // 变压器
      TPL(tid('en'), 'transformer', '主变压器A', 'Transformer', -1.5, 0, -1.5, '#ff8c00', { emissive: '#ff8c0022', metalness: 0.7 }),
      TPL(tid('en'), 'transformer', '主变压器B', 'Transformer', 1.5, 0, -1.5, '#ff8c00', { emissive: '#ff8c0022', metalness: 0.7 }),
      // 断路器
      TPL(tid('en'), 'circuitbreaker', '断路器A', 'CircuitBreaker', -1, 0, -2.2, '#ff4d4f', { emissive: '#ff4d4f22' }),
      TPL(tid('en'), 'circuitbreaker', '断路器B', 'CircuitBreaker', 1, 0, -2.2, '#ff4d4f', { emissive: '#ff4d4f22' }),
      // 电线杆
      TPL(tid('en'), 'powerpole', '电线杆A', 'PowerPole', -2.5, 0, -2.5, '#1a1a1a', { emissive: '#080808' }),
      TPL(tid('en'), 'powerpole', '电线杆B', 'PowerPole', 2.5, 0, -2.5, '#1a1a1a', { emissive: '#080808' }),
      TPL(tid('en'), 'insulator', '绝缘子A', 'Insulator', -2, 0, -3, '#e8e8e8', { emissive: '#e8e8e822' }),
      // 太阳能
      TPL(tid('en'), 'solarpanel', '光伏板A', 'SolarPanel', -2, 0, 1, '#13c2c2', { sx: 1.2, sz: 0.8, emissive: '#13c2c222' }),
      TPL(tid('en'), 'solarpanel', '光伏板B', 'SolarPanel', -2.8, 0, 1, '#13c2c2', { sx: 1.2, sz: 0.8, emissive: '#13c2c222' }),
      // 风力发电
      TPL(tid('en'), 'windturbine', '风力发电机A', 'WindTurbine', 2, 0, 2, '#e8e8e8', { emissive: '#88888811' }),
      TPL(tid('en'), 'windturbine', '风力发电机B', 'WindTurbine', 3, 0, 2, '#e8e8e8', { emissive: '#88888811' }),
      // 储能
      TPL(tid('en'), 'batterystorage', '储能电池组', 'BatteryStorage', 2, 0, -0.5, '#0d2a1a', { emissive: '#051510', sx: 1.5 }),
      // 电缆桥架
      TPL(tid('en'), 'cabletray', '电缆桥架', 'CableTray', 0, 0, -3, '#1a1a1a', { sx: 4, emissive: '#080808' }),
      // 电表
      TPL(tid('en'), 'powermeter', '智能电表', 'PowerMeter', 0, 0, -2, '#13c2c2', { emissive: '#13c2c222' }),
      // 特效
      TPL(tid('en'), 'scanLine', '扫描线', 'ScanLine', 0, 0.5, 0, '#13c2c2', { sy: 1.2 }),
      TPL(tid('en'), 'flyLine', '电流流光A', 'FlyLine', -0.8, 0.3, -2.5, '#ff8c00', { sx: 2.0 }),
      TPL(tid('en'), 'flyLine', '电流流光B', 'FlyLine', 0.8, 0.3, -2.5, '#ff4d4f', { sx: 2.0 }),
      TPL(tid('en'), 'electricArc', '电弧特效', 'ElectricArc', -1.5, 1, -1.5, '#ffaa4d', { sy: 1.2 }),
      TPL(tid('en'), 'glowSphere', '辉光球', 'GlowSphere', 0, 1.5, 0, '#ff8c00', { sx: 0.3, sy: 0.3, sz: 0.3 }),
      // 数据标签
      DATA_LABEL('en', '能源变电站', 0, 3.5, -3, '#13c2c2', 32),
      DATA_LABEL('en', '电压: 110kV', -3, 1.5, 2, '#ffaa4d', 20),
      DATA_LABEL('en', '负载: 78%', 3, 1.5, 2, '#4dffa6', 20),
    ]
  },

  // ==================== 8. 仓储物流中心 ====================
  {
    id: 'warehouse',
    name: '仓储物流中心',
    icon: 'fas fa-warehouse',
    category: '仓储物流',
    description: '大型仓储物流中心，发光货架+流光传送带+AGV路径线',
    sceneSettings: NIGHT_SCENE({ x: 8, y: 5, z: 7 }, { x: 0, y: 1.0, z: 0 }),
    objects: [
      TPL(tid('wh'), 'gridFloor', '地面网格', 'GridFloor', 0, 0.01, 0, '#13c2c2', { sx: 3, sz: 2.5 }),
      TPL(tid('wh'), 'plane', '仓库地板', 'Plane', 0, 0, 0, '#0d1a16', { sx: 6, sz: 5, emissive: '#040810' }),
      // 货架
      TPL(tid('wh'), 'warehouserack', '货架A1', 'WarehouseRack', -2, 0, -2, '#1a1000', { emissive: '#080500' }),
      TPL(tid('wh'), 'warehouserack', '货架A2', 'WarehouseRack', -2, 0, -1, '#1a1000', { emissive: '#080500' }),
      TPL(tid('wh'), 'warehouserack', '货架A3', 'WarehouseRack', -2, 0, 0, '#1a1000', { emissive: '#080500' }),
      TPL(tid('wh'), 'warehouserack', '货架A4', 'WarehouseRack', -2, 0, 1, '#1a1000', { emissive: '#080500' }),
      TPL(tid('wh'), 'warehouserack', '货架B1', 'WarehouseRack', 2, 0, -2, '#1a1000', { emissive: '#080500' }),
      TPL(tid('wh'), 'warehouserack', '货架B2', 'WarehouseRack', 2, 0, -1, '#1a1000', { emissive: '#080500' }),
      TPL(tid('wh'), 'warehouserack', '货架B3', 'WarehouseRack', 2, 0, 0, '#1a1000', { emissive: '#080500' }),
      TPL(tid('wh'), 'warehouserack', '货架B4', 'WarehouseRack', 2, 0, 1, '#1a1000', { emissive: '#080500' }),
      // 堆垛机
      TPL(tid('wh'), 'stackercrane', '堆垛机A', 'StackerCrane', -2, 0, -1.5, '#ff8c00', { emissive: '#ff8c0022' }),
      TPL(tid('wh'), 'stackercrane', '堆垛机B', 'StackerCrane', 2, 0, -1.5, '#ff8c00', { emissive: '#ff8c0022' }),
      // 传送分拣
      TPL(tid('wh'), 'conveyorjunc', '传送交汇A', 'ConveyorJunc', 0, 0.2, 0, '#1a1a1a', { emissive: '#080808' }),
      TPL(tid('wh'), 'conveyorjunc', '传送交汇B', 'ConveyorJunc', 0, 0.2, -1, '#1a1a1a', { emissive: '#080808' }),
      TPL(tid('wh'), 'sorter', '分拣机', 'Sorter', 0, 0, -0.5, '#0d1f2a', { emissive: '#081018' }),
      // AGV
      TPL(tid('wh'), 'agv', 'AGV小车A', 'AGV', -1, 0, 2, '#13c2c244', { sx: 0.8 }),
      TPL(tid('wh'), 'agv', 'AGV小车B', 'AGV', 1, 0, 2, '#36cfc944', { sx: 0.8 }),
      // 叉车
      TPL(tid('wh'), 'forklift', '叉车A', 'Forklift', -1, 0, -2.8, '#ff8c00', { emissive: '#ff8c0022' }),
      TPL(tid('wh'), 'forklift', '叉车B', 'Forklift', 1, 0, -2.8, '#ff8c00', { emissive: '#ff8c0022', ry: Math.PI }),
      // 装卸台
      TPL(tid('wh'), 'loadingdock', '装卸台A', 'LoadingDock', -1.5, 0, -3.5, '#1a1a1a', { emissive: '#080808' }),
      TPL(tid('wh'), 'loadingdock', '装卸台B', 'LoadingDock', 1.5, 0, -3.5, '#1a1a1a', { emissive: '#080808' }),
      // 集装箱
      TPL(tid('wh'), 'container', '集装箱A', 'Container', -3, 0, 3, '#13c2c2', { emissive: '#13c2c222' }),
      TPL(tid('wh'), 'container', '集装箱B', 'Container', -1.5, 0, 3, '#36cfc9', { emissive: '#36cfc922' }),
      // 特效
      TPL(tid('wh'), 'scanLine', '扫描线', 'ScanLine', 0, 0.5, 0, '#13c2c2', { sy: 1.2 }),
      TPL(tid('wh'), 'flyLine', '传送流光A', 'FlyLine', 0, 0.35, 0.3, '#13c2c2', { sx: 1.5 }),
      TPL(tid('wh'), 'flyLine', '传送流光B', 'FlyLine', 0, 0.35, -0.7, '#36cfc9', { sx: 1.5 }),
      TPL(tid('wh'), 'glowSphere', '辉光球', 'GlowSphere', 0, 1.5, 0, '#13c2c2', { sx: 0.25, sy: 0.25, sz: 0.25 }),
      // 数据标签
      DATA_LABEL('wh', '仓储物流中心', 0, 3.5, -3.5, '#13c2c2', 32),
      DATA_LABEL('wh', '库存: 12,580', -4, 1.5, 2, '#4dffa6', 20),
      DATA_LABEL('wh', '吞吐: 98%', 4, 1.5, 2, '#4dffa6', 20),
    ]
  },

  // ==================== 9. 城市交通枢纽 ====================
  {
    id: 'transport-hub',
    name: '城市交通枢纽',
    icon: 'fas fa-bus-alt',
    category: '交通设施',
    description: '夜景城市交通枢纽，发光道路+流光车流+悬浮信息牌',
    sceneSettings: NIGHT_SCENE({ x: 7, y: 5, z: 7 }, { x: 0, y: 1.0, z: 0 }),
    objects: [
      TPL(tid('tr'), 'gridFloor', '地面网格', 'GridFloor', 0, 0.01, 0, '#13c2c2', { sx: 3, sz: 3 }),
      // 道路
      TPL(tid('tr'), 'road', '主干道X', 'Road', 0, 0.02, -3.5, '#1a1a2a', { ry: Math.PI/2, sx: 0.5, sz: 10 }),
      TPL(tid('tr'), 'road', '主干道Z', 'Road', -3.5, 0.02, 0, '#1a1a2a', { sx: 0.5, sz: 10 }),
      TPL(tid('tr'), 'zebracrossing', '斑马线', 'ZebraCrossing', 0, 0.03, -2.2, '#2a2a3a', { sx: 3 }),
      // 火车/地铁
      TPL(tid('tr'), 'train', '高铁列车', 'Train', -2, 0, -2, '#13c2c244', { sx: 1.5 }),
      TPL(tid('tr'), 'metro', '地铁列车', 'Metro', -1.5, 0, 1.5, '#36cfc944', { sx: 1.3 }),
      // 巴士站
      TPL(tid('tr'), 'busstop', '公交站台', 'BusStop', 1.5, 0, -2, '#0d1f2a', { emissive: '#081018' }),
      TPL(tid('tr'), 'bus', '公交车A', 'Bus', 1, 0, -2.5, '#13c2c244', { ry: Math.PI/2 }),
      // 车辆
      TPL(tid('tr'), 'car', '轿车A', 'Car', -1, 0, -2.8, '#13c2c244', { ry: Math.PI/2 }),
      TPL(tid('tr'), 'suv', 'SUV', 'SUV', 0, 0, -2.8, '#4dffa644', { ry: Math.PI/2 }),
      TPL(tid('tr'), 'truck', '货车', 'Truck', 2, 0, -3, '#ffaa4d44', { ry: Math.PI/2 }),
      // 飞机
      TPL(tid('tr'), 'airplane', '客机', 'Airplane', 0, 0, 3, '#e8e8e8', { sy: 0.6, emissive: '#88888811' }),
      // 交通设施
      TPL(tid('tr'), 'trafficlight', '信号灯', 'TrafficLight', -0.3, 0, -3.2, '#ff4d4f', { emissive: '#ff4d4f33' }),
      TPL(tid('tr'), 'barrier', '护栏', 'Barrier', 0, 0, -3, '#1a1a0d', { sx: 1.5, emissive: '#080800' }),
      TPL(tid('tr'), 'trafficcone', '锥形桶A', 'TrafficCone', -0.8, 0, -3.5, '#ff8c00', { emissive: '#ff8c0022' }),
      TPL(tid('tr'), 'trafficcone', '锥形桶B', 'TrafficCone', 0.8, 0, -3.5, '#ff8c00', { emissive: '#ff8c0022' }),
      // 监控
      TPL(tid('tr'), 'surveilcam', '交通摄像头', 'SurveilCam', -1, 0, -3.2, '#13c2c2', { sy: 1.5 }),
      // 特效
      TPL(tid('tr'), 'scanLine', '扫描线', 'ScanLine', 0, 0.5, 0, '#13c2c2', { sy: 1.2 }),
      TPL(tid('tr'), 'flyLine', '车流光带A', 'FlyLine', 0, 0.1, -3.0, '#13c2c2', { sx: 3.0 }),
      TPL(tid('tr'), 'flyLine', '车流光带B', 'FlyLine', -3.0, 0.1, 0, '#36cfc9', { sx: 3.0 }),
      TPL(tid('tr'), 'glowSphere', '辉光球', 'GlowSphere', 0, 2, 0, '#13c2c2', { sx: 0.3, sy: 0.3, sz: 0.3 }),
      // 数据标签
      DATA_LABEL('tr', '城市交通枢纽', 0, 4.5, -3, '#13c2c2', 32),
      DATA_LABEL('tr', '实时客流: 2,847', -4, 2, 0, '#4dffa6', 20),
      DATA_LABEL('tr', '列车准点: 100%', 4, 2, 0, '#4dffa6', 20),
    ]
  },

  // ==================== 10. 现代办公空间 ====================
  {
    id: 'modern-office',
    name: '现代办公空间',
    icon: 'fas fa-building',
    category: '室内设施',
    description: '明亮现代办公空间，木质工位+大屏显示器+休闲区（夜间模式）',
    sceneSettings: NIGHT_SCENE({ x: 5, y: 3.5, z: 5 }, { x: 0, y: 1.0, z: 0 }),
    objects: [
      TPL(tid('of'), 'gridFloor', '地面网格', 'GridFloor', 0, 0.01, 0, '#13c2c2', { sx: 2, sz: 2 }),
      TPL(tid('of'), 'plane', '办公地板', 'Plane', 0, 0, 0, '#0d1a22', { sx: 4, sz: 4, emissive: '#040810' }),
      // 工位
      TPL(tid('of'), 'desk', '工位桌A', 'Desk', -1.5, 0, 1, '#1a0d00', { sx: 0.9 }),
      TPL(tid('of'), 'desk', '工位桌B', 'Desk', 1.5, 0, 1, '#1a0d00', { sx: 0.9 }),
      TPL(tid('of'), 'desk', '工位桌C', 'Desk', -1.5, 0, -1, '#1a0d00', { sx: 0.9 }),
      TPL(tid('of'), 'desk', '工位桌D', 'Desk', 1.5, 0, -1, '#1a0d00', { sx: 0.9 }),
      TPL(tid('of'), 'officechair', '办公椅A', 'OfficeChair', -1.5, 0, 1.4, '#1a1a1a', { sx: 0.8 }),
      TPL(tid('of'), 'officechair', '办公椅B', 'OfficeChair', 1.5, 0, 1.4, '#1a1a1a', { sx: 0.8 }),
      // 显示器（发光）
      TPL(tid('of'), 'monitor', '显示器A', 'Monitor', -1.5, 0, 0.9, '#13c2c2', { emissive: '#13c2c211' }),
      TPL(tid('of'), 'monitor', '显示器B', 'Monitor', 1.5, 0, 0.9, '#13c2c2', { emissive: '#13c2c211' }),
      // 会议室
      TPL(tid('of'), 'meetingtable', '会议桌', 'MeetingTable', 0, 0, -2, '#1a0d00', { sx: 1.2 }),
      TPL(tid('of'), 'officechair', '会议椅1', 'OfficeChair', -0.6, 0, -2.5, '#1a1a1a', { sx: 0.8 }),
      TPL(tid('of'), 'officechair', '会议椅2', 'OfficeChair', 0.6, 0, -2.5, '#1a1a1a', { sx: 0.8 }),
      TPL(tid('of'), 'whiteboard', '白板', 'Whiteboard', 0, 0, -3, '#e8e8e8', { sx: 0.8 }),
      // 休闲区
      TPL(tid('of'), 'sofa', '沙发', 'Sofa', -2, 0, -2.5, '#0d1f2a', { sx: 1.0, emissive: '#081018' }),
      TPL(tid('of'), 'coffeetable', '茶几', 'CoffeeTable', -2, 0, -2, '#1a0d00', { sx: 0.8 }),
      // 设备
      TPL(tid('of'), 'printer', '打印机', 'Printer', 2.5, 0, -2, '#1a1a1a', { emissive: '#080808' }),
      TPL(tid('of'), 'fridge', '冰箱', 'Fridge', 2.5, 0, -0.5, '#e8e8e8', { sx: 0.8 }),
      // 灯光
      TPL(tid('of'), 'ceilinglight', '吸顶灯A', 'CeilingLight', -1, 3, 0, '#fff5e0', { emissive: '#fff5e044' }),
      TPL(tid('of'), 'ceilinglight', '吸顶灯B', 'CeilingLight', 1, 3, 0, '#fff5e0', { emissive: '#fff5e044' }),
      // 特效
      TPL(tid('of'), 'scanLine', '扫描线', 'ScanLine', 0, 0.5, 0, '#13c2c2', { sy: 1 }),
      TPL(tid('of'), 'glowSphere', '辉光球', 'GlowSphere', 0, 1.5, 0, '#13c2c2', { sx: 0.2, sy: 0.2, sz: 0.2 }),
      // 数据标签
      DATA_LABEL('of', '现代办公空间', 0, 3.5, -3, '#13c2c2', 32),
      DATA_LABEL('of', '工位占用: 4/20', -3, 1.5, 0, '#4dffa6', 20),
      DATA_LABEL('of', '会议室: 空闲', 3, 1.5, 0, '#4dffa6', 20),
    ]
  },

  // ==================== 11. 消防安防指挥中心 ====================
  {
    id: 'fire-security',
    name: '消防安防指挥中心',
    icon: 'fas fa-shield-alt',
    category: '消防安防',
    description: '红色警报风消防安防中心，发光灭火器+监控球机+逃生指引光带',
    sceneSettings: NIGHT_SCENE({ x: 5, y: 4, z: 5 }, { x: 0, y: 1.0, z: 0 }),
    objects: [
      TPL(tid('fs'), 'gridFloor', '地面网格', 'GridFloor', 0, 0.01, 0, '#ff4d4f', { sx: 2, sz: 2 }),
      TPL(tid('fs'), 'plane', '指挥中心地板', 'Plane', 0, 0, 0, '#0d1a16', { sx: 4, sz: 4, emissive: '#040810' }),
      // 消防设备（红色 emissive）
      TPL(tid('fs'), 'fireextinguisher', '灭火器A', 'FireExtinguisher', -2, 0, 1, '#ff4d4f', { emissive: '#ff4d4f33' }),
      TPL(tid('fs'), 'fireextinguisher', '灭火器B', 'FireExtinguisher', -2, 0, 0, '#ff4d4f', { emissive: '#ff4d4f33' }),
      TPL(tid('fs'), 'firehydrantbox', '消防栓箱', 'FireHydrantBox', -2, 0, -1, '#ff4d4f', { emissive: '#ff4d4f22' }),
      TPL(tid('fs'), 'hydrant', '室外消防栓', 'Hydrant', -2, 0, -2.5, '#ff4d4f', { emissive: '#ff4d4f33' }),
      // 喷淋
      TPL(tid('fs'), 'sprinkler', '喷淋头A', 'Sprinkler', -1, 2.5, 1, '#ff4d4f', { emissive: '#ff4d4f22' }),
      TPL(tid('fs'), 'sprinkler', '喷淋头B', 'Sprinkler', 0, 2.5, 1, '#ff4d4f', { emissive: '#ff4d4f22' }),
      // 探测报警
      TPL(tid('fs'), 'smokedetector', '烟感A', 'SmokeDetector', -1, 2.5, 0, '#1a1a1a', { emissive: '#080808' }),
      TPL(tid('fs'), 'firealarm', '火警铃', 'FireAlarm', 0, 2.5, -2, '#ff4d4f', { emissive: '#ff4d4f44' }),
      // 安防设备
      TPL(tid('fs'), 'securitygate', '安检门', 'SecurityGate', 2, 0, 1, '#13c2c2', { emissive: '#13c2c222' }),
      TPL(tid('fs'), 'accesscontrol', '门禁', 'AccessControl', 2, 0, 0.5, '#13c2c2', { emissive: '#13c2c233' }),
      TPL(tid('fs'), 'cctv', '监控球机A', 'CCTV', -1.5, 2.5, -2, '#13c2c2', { sy: 2, emissive: '#13c2c222' }),
      TPL(tid('fs'), 'cctv', '监控球机B', 'CCTV', 1.5, 2.5, -2, '#13c2c2', { sy: 2, emissive: '#13c2c222' }),
      // 逃生指引
      TPL(tid('fs'), 'emergencysign', '应急出口A', 'EmergencySign', -1, 0, 2.5, '#4dffa6', { emissive: '#4dffa644' }),
      TPL(tid('fs'), 'emergencysign', '应急出口B', 'EmergencySign', 1, 0, 2.5, '#4dffa6', { emissive: '#4dffa644' }),
      // 指挥中心大屏
      TPL(tid('fs'), 'dashboardscreen', '安防监控大屏', 'DashboardScreen', 0, 0, -2.5, '#0d1f2a', { sx: 1.2, sy: 0.9 }),
      // 特效
      TPL(tid('fs'), 'scanLine', '扫描线', 'ScanLine', 0, 0.5, 0, '#ff4d4f', { sy: 1.2 }),
      TPL(tid('fs'), 'flyLine', '逃生指引光带', 'FlyLine', 0, 0.1, 2.0, '#4dffa6', { sx: 2.0 }),
      TPL(tid('fs'), 'glowSphere', '警示辉光', 'GlowSphere', 0, 1.5, 0, '#ff4d4f', { sx: 0.3, sy: 0.3, sz: 0.3 }),
      // 数据标签
      DATA_LABEL('fs', '消防安防中心', 0, 3.5, -2.5, '#ff4d4f', 32),
      DATA_LABEL('fs', '系统正常', -3, 1.5, 2, '#4dffa6', 20),
      DATA_LABEL('fs', '监控在线: 24', 3, 1.5, 2, '#4dffa6', 20),
    ]
  },

  // ==================== 12. 中式庭院景观 ====================
  {
    id: 'garden-landscape',
    name: '中式庭院景观',
    icon: 'fas fa-tree',
    category: '绿化景观',
    description: '夜间中式园林庭院，木质凉亭+荷花池+假山石+青石板路',
    sceneSettings: NIGHT_SCENE({ x: 6, y: 4, z: 6 }, { x: 0, y: 0.5, z: 0 }),
    objects: [
      TPL(tid('gl'), 'gridFloor', '地面网格', 'GridFloor', 0, 0.01, 0, '#13c2c2', { sx: 2.5, sz: 2.5 }),
      TPL(tid('gl'), 'plane', '庭院地面', 'Plane', 0, 0, 0, '#0d1a16', { sx: 5, sz: 5, emissive: '#040810' }),
      // 凉亭
      TPL(tid('gl'), 'pavilion', '凉亭', 'Pavilion', 0, 0, 0, '#1a0d00', { sx: 1.2, sy: 1.2, sz: 1.2 }),
      // 荷花池
      TPL(tid('gl'), 'cylinder', '荷花池', 'Cylinder', 1.5, 0.05, 0, '#13c2c222', { sx: 1.2, sy: 0.1, sz: 0.8 }),
      // 围墙
      TPL(tid('gl'), 'gardenwall', '围墙N', 'GardenWall', 0, 0, -2.5, '#1a100d', { sx: 3, emissive: '#080500' }),
      TPL(tid('gl'), 'gardenwall', '围墙W', 'GardenWall', -2.5, 0, 0, '#1a100d', { ry: Math.PI/2, sx: 3 }),
      TPL(tid('gl'), 'gate', '院门', 'Gate', 0, 0, -2.5, '#1a0d00', { sx: 0.8 }),
      TPL(tid('gl'), 'fence', '木栅栏S', 'Fence', 0, 0, 2.5, '#1a0d00', { sx: 2.5 }),
      // 树木
      TPL(tid('gl'), 'citytree', '松树A', 'CityTree', -1.5, 0, -1.5, '#0a1a0a', { sx: 0.8, sy: 0.9, sz: 0.8 }),
      TPL(tid('gl'), 'citytree', '松树B', 'CityTree', 1.5, 0, -1.5, '#0a1a0a', { sx: 0.8, sy: 0.9, sz: 0.8 }),
      TPL(tid('gl'), 'bush', '灌木A', 'Bush', -1, 0, 2, '#0a1a0a', { sx: 0.8 }),
      TPL(tid('gl'), 'hedge', '树篱', 'Hedge', 0, 0, 2.2, '#0a1a0a', { sx: 1.5 }),
      // 假山石
      TPL(tid('gl'), 'rock', '假山石A', 'Rock', -1.8, 0, 0.3, '#1a1a1a', { sx: 0.8 }),
      TPL(tid('gl'), 'rock', '假山石B', 'Rock', 1.8, 0, 0.5, '#1a1a1a', { sx: 0.6 }),
      // 草坪
      TPL(tid('gl'), 'plane', '草坪', 'Plane', 0, 0.01, 1, '#0a1a0a', { sx: 2, sz: 1 }),
      // 喷泉
      TPL(tid('gl'), 'fountain', '喷泉', 'Fountain', 0, 0, 0.8, '#13c2c222', { sy: 0.8 }),
      // 长椅
      TPL(tid('gl'), 'bench', '长椅A', 'Bench', -0.5, 0, 2, '#1a0d00', { sx: 0.8 }),
      // 庭院灯
      TPL(tid('gl'), 'lamppost', '庭院灯A', 'LampPost', -1.2, 0, -1.5, '#ffaa4d', { emissive: '#ffaa4d22' }),
      TPL(tid('gl'), 'lamppost', '庭院灯B', 'LampPost', 1.2, 0, -1.5, '#ffaa4d', { emissive: '#ffaa4d22' }),
      // 特效
      TPL(tid('gl'), 'scanLine', '扫描线', 'ScanLine', 0, 0.5, 0, '#13c2c2', { sy: 1 }),
      TPL(tid('gl'), 'glowSphere', '辉光球', 'GlowSphere', 0, 1, 0, '#ffaa4d', { sx: 0.2, sy: 0.2, sz: 0.2 }),
      TPL(tid('gl'), 'particleCloud', '氛围粒子', 'ParticleCloud', 0, 0.5, 0, '#13c2c211', { sx: 2, sy: 1, sz: 2 }),
      // 数据标签
      DATA_LABEL('gl', '中式庭院景观', 0, 3, -2.5, '#13c2c2', 32),
      DATA_LABEL('gl', '温度: 26C', -3, 1, 0, '#4dffa6', 20),
      DATA_LABEL('gl', '湿度: 65%', 3, 1, 0, '#4dffa6', 20),
    ]
  },

  // ==================== 13. 智慧医院病房 ====================
  {
    id: 'hospital-ward',
    name: '智慧医院病房',
    icon: 'fas fa-hospital',
    category: '医疗设施',
    description: '冷白洁净风医院病房，发光监护仪+输液架+MRI/CT设备',
    sceneSettings: NIGHT_SCENE({ x: 5, y: 3.5, z: 5 }, { x: 0, y: 1.0, z: 0 }),
    objects: [
      TPL(tid('hp'), 'gridFloor', '地面网格', 'GridFloor', 0, 0.01, 0, '#13c2c2', { sx: 2, sz: 2 }),
      TPL(tid('hp'), 'plane', '病房地板', 'Plane', 0, 0, 0, '#0d1a22', { sx: 4, sz: 4, emissive: '#040810' }),
      // 病床
      TPL(tid('hp'), 'hospitalbed', '病床A', 'HospitalBed', -1.2, 0, 0, '#e8e8e8', { sx: 0.9 }),
      TPL(tid('hp'), 'hospitalbed', '病床B', 'HospitalBed', 1.2, 0, 0, '#e8e8e8', { sx: 0.9 }),
      // 监护仪
      TPL(tid('hp'), 'hospitalmonitor', '监护仪A', 'HospitalMonitor', -1.8, 0, 0, '#13c2c2', { emissive: '#13c2c222' }),
      TPL(tid('hp'), 'hospitalmonitor', '监护仪B', 'HospitalMonitor', 1.8, 0, 0, '#13c2c2', { emissive: '#13c2c222' }),
      // 输液架
      TPL(tid('hp'), 'ivstand', '输液架A', 'IVStand', -1.2, 0, 0.5, '#c0c0c0', { sx: 0.8 }),
      TPL(tid('hp'), 'ivstand', '输液架B', 'IVStand', 1.2, 0, 0.5, '#c0c0c0', { sx: 0.8 }),
      // 大型设备
      TPL(tid('hp'), 'operatingtable', '手术台', 'OperatingTable', 0, 0, -2, '#e8e8e8', { sx: 0.9 }),
      TPL(tid('hp'), 'examinationlamp', '手术灯', 'ExaminationLamp', 0, 0, -2.3, '#fff5e0', { emissive: '#fff5e044', sy: 1.5 }),
      TPL(tid('hp'), 'mrimachine', 'MRI设备', 'MRIMachine', -2.5, 0, 2, '#1a1a1a', { emissive: '#080808' }),
      TPL(tid('hp'), 'ctscanner', 'CT扫描仪', 'CTScanner', 2.5, 0, 2, '#e8e8e8', { sx: 0.9 }),
      // 消防
      TPL(tid('hp'), 'fireextinguisher', '灭火器', 'FireExtinguisher', -2.5, 0, 1, '#ff4d4f', { emissive: '#ff4d4f22' }),
      TPL(tid('hp'), 'emergencysign', '应急出口', 'EmergencySign', 0, 0, 3, '#4dffa6', { emissive: '#4dffa644' }),
      // 特效
      TPL(tid('hp'), 'scanLine', '扫描线', 'ScanLine', 0, 0.5, 0, '#13c2c2', { sy: 1 }),
      TPL(tid('hp'), 'glowSphere', '辉光球', 'GlowSphere', 0, 1.5, 0, '#13c2c2', { sx: 0.2, sy: 0.2, sz: 0.2 }),
      // 数据标签
      DATA_LABEL('hp', '智慧医院病房', 0, 3, -2.5, '#13c2c2', 32),
      DATA_LABEL('hp', '心率: 72bpm', -3, 1, 0, '#4dffa6', 20),
      DATA_LABEL('hp', '血氧: 98%', 3, 1, 0, '#4dffa6', 20),
    ]
  },

  // ==================== 14. 智慧农场 ====================
  {
    id: 'smart-farm',
    name: '智慧农场',
    icon: 'fas fa-tractor',
    category: '农业设施',
    description: '绿色生态智慧农场，温室大棚+喷灌系统+气象站+无人机监控',
    sceneSettings: NIGHT_SCENE({ x: 6, y: 4, z: 6 }, { x: 0, y: 0.5, z: 0 }),
    objects: [
      TPL(tid('sf'), 'gridFloor', '地面网格', 'GridFloor', 0, 0.01, 0, '#13c2c2', { sx: 2.5, sz: 2.5 }),
      TPL(tid('sf'), 'plane', '农场地面', 'Plane', 0, 0, 0, '#0d1a16', { sx: 5, sz: 5, emissive: '#040810' }),
      // 温室
      TPL(tid('sf'), 'greenhouse', '温室大棚A', 'Greenhouse', -2, 0, 1.5, '#e8e8e844', { sx: 1.5, sz: 1.5 }),
      TPL(tid('sf'), 'greenhouse', '育苗温室', 'Greenhouse', 2, 0, 1.5, '#e8e8e844', { sx: 1.3, sz: 1.3 }),
      // 农机
      TPL(tid('sf'), 'tractor', '拖拉机', 'Tractor', 0, 0, -2, '#ff4d4f', { emissive: '#ff4d4f22' }),
      TPL(tid('sf'), 'harvester', '收割机', 'Harvester', 1.5, 0, -2, '#ff8c00', { emissive: '#ff8c0022' }),
      // 喷灌
      TPL(tid('sf'), 'irrigationpivot', '喷灌机', 'IrrigationPivot', 0, 0, 1, '#13c2c2', { sx: 2, emissive: '#13c2c222' }),
      // 粮仓
      TPL(tid('sf'), 'grainsilo', '粮仓A', 'GrainSilo', -2.5, 0, -1.5, '#1a0d00', { emissive: '#080500' }),
      TPL(tid('sf'), 'grainsilo', '粮仓B', 'GrainSilo', -2.5, 0, -0.5, '#1a0d00', { emissive: '#080500' }),
      // 种植床
      TPL(tid('sf'), 'plantingbed', '种植床A', 'PlantingBed', -1, 0, 3, '#0a1a0a', { sx: 0.8 }),
      TPL(tid('sf'), 'plantingbed', '种植床B', 'PlantingBed', 1, 0, 3, '#0a1a0a', { sx: 0.8 }),
      // 气象站
      TPL(tid('sf'), 'weatherstation', '气象站', 'WeatherStation', 2.5, 0, 2.5, '#13c2c2', { emissive: '#13c2c222' }),
      TPL(tid('sf'), 'surveilcam', '农场监控', 'SurveilCam', 2.8, 0, 2, '#13c2c2', { sy: 2 }),
      // 无人机
      TPL(tid('sf'), 'drone', '植保无人机', 'Drone', 0, 0, 2.0, '#1a1a1a', { emissive: '#080808' }),
      // 特效
      TPL(tid('sf'), 'scanLine', '扫描线', 'ScanLine', 0, 0.5, 0, '#13c2c2', { sy: 1 }),
      TPL(tid('sf'), 'glowSphere', '辉光球', 'GlowSphere', 0, 1.5, 0, '#4dffa6', { sx: 0.25, sy: 0.25, sz: 0.25 }),
      TPL(tid('sf'), 'particleCloud', '氛围粒子', 'ParticleCloud', 0, 0.5, 0, '#4dffa611', { sx: 2, sy: 1, sz: 2 }),
      // 数据标签
      DATA_LABEL('sf', '智慧农场', 0, 3, -2.5, '#13c2c2', 32),
      DATA_LABEL('sf', '温度: 28C', -3, 1, 2, '#4dffa6', 20),
      DATA_LABEL('sf', '土壤湿度: 72%', 3, 1, 2, '#4dffa6', 20),
    ]
  },

  // ==================== 15. 水处理厂 ====================
  {
    id: 'water-plant',
    name: '水处理厂',
    icon: 'fas fa-water',
    category: '水利水务',
    description: '工业蓝水务处理厂，蓝色水池+流光管道+泵站+监控数据面板',
    sceneSettings: NIGHT_SCENE({ x: 6, y: 4, z: 6 }, { x: 0, y: 0.8, z: 0 }),
    objects: [
      TPL(tid('wp'), 'gridFloor', '地面网格', 'GridFloor', 0, 0.01, 0, '#13c2c2', { sx: 2.5, sz: 2.5 }),
      TPL(tid('wp'), 'plane', '水厂地面', 'Plane', 0, 0, 0, '#0d1a22', { sx: 5, sz: 5, emissive: '#040810' }),
      // 水塔
      TPL(tid('wp'), 'watertower', '水塔', 'WaterTower', -2, 0, 2, '#13c2c2', { emissive: '#13c2c222' }),
      // 处理池
      TPL(tid('wp'), 'cylinder', '进水调节池', 'Cylinder', -1.5, 0.05, -1, '#0d1f2a44', { sx: 1.2, sy: 0.1, sz: 0.8 }),
      TPL(tid('wp'), 'cylinder', '沉淀池', 'Cylinder', 0, 0.05, -1, '#0d1f2a66', { sx: 1.2, sy: 0.1, sz: 0.8 }),
      TPL(tid('wp'), 'cylinder', '曝气池', 'Cylinder', 1.5, 0.05, -1, '#13c2c222', { sx: 1.2, sy: 0.1, sz: 0.8 }),
      // 管道
      TPL(tid('wp'), 'pipe', '主水管', 'Pipe', 0, 0.2, -1.8, '#13c2c244', { ry: Math.PI/2, sx: 3, emissive: '#13c2c211' }),
      // 泵站
      TPL(tid('wp'), 'pump', '提升泵', 'Pump', -2.5, 0, -1, '#0d1f2a', { emissive: '#081018' }),
      TPL(tid('wp'), 'valve', '控制阀A', 'Valve', -2, 0, -2.5, '#13c2c2', { sx: 0.6 }),
      TPL(tid('wp'), 'valve', '控制阀B', 'Valve', 2, 0, -2.5, '#13c2c2', { sx: 0.6 }),
      // 监控
      TPL(tid('wp'), 'cctv', '监控球机', 'CCTV', 0, 0, 2, '#13c2c2', { sy: 2, emissive: '#13c2c222' }),
      TPL(tid('wp'), 'sensorarray', '水质传感器', 'SensorArray', -1.5, 0, -0.5, '#13c2c2', { sy: 1.5 }),
      // 特效
      TPL(tid('wp'), 'scanLine', '扫描线', 'ScanLine', 0, 0.5, 0, '#13c2c2', { sy: 1.2 }),
      TPL(tid('wp'), 'flyLine', '水流光带A', 'FlyLine', -1.0, 0.15, -0.8, '#13c2c2', { sx: 1.5 }),
      TPL(tid('wp'), 'flyLine', '水流光带B', 'FlyLine', 0.8, 0.15, -0.8, '#36cfc9', { sx: 1.5 }),
      TPL(tid('wp'), 'glowSphere', '辉光球', 'GlowSphere', 0, 1.5, 0, '#13c2c2', { sx: 0.25, sy: 0.25, sz: 0.25 }),
      // 数据标签
      DATA_LABEL('wp', '水处理厂', 0, 3, -2.5, '#13c2c2', 32),
      DATA_LABEL('wp', 'PH: 7.2', -3, 1, 2, '#4dffa6', 20),
      DATA_LABEL('wp', '浊度: 0.5NTU', 3, 1, 2, '#4dffa6', 20),
    ]
  },

  // ==================== 16. 运动场馆 ====================
  {
    id: 'sports-venue',
    name: '运动场馆',
    icon: 'fas fa-futbol',
    category: '运动休闲',
    description: '夜间体育场，发光球场+看台+记分牌+泛光灯',
    sceneSettings: NIGHT_SCENE({ x: 7, y: 5, z: 7 }, { x: 0, y: 1.0, z: 0 }),
    objects: [
      TPL(tid('sp'), 'gridFloor', '地面网格', 'GridFloor', 0, 0.01, 0, '#13c2c2', { sx: 3, sz: 3 }),
      TPL(tid('sp'), 'plane', '体育场地板', 'Plane', 0, 0, 0, '#0d1a16', { sx: 6, sz: 6, emissive: '#040810' }),
      // 体育场
      TPL(tid('sp'), 'stadium', '体育场', 'Stadium', 0, 0, 0, '#0d1a16', { sx: 2, sz: 2, emissive: '#040a08' }),
      // 看台
      TPL(tid('sp'), 'bleacher', '看台A', 'Bleacher', 0, 0, -2, '#1a1a1a', { sx: 1.5, emissive: '#080808' }),
      // 记分牌
      TPL(tid('sp'), 'scoreboard', '记分牌', 'Scoreboard', 0, 0, -2.3, '#0d1f2a', { emissive: '#081018' }),
      // 球门
      TPL(tid('sp'), 'box', '球门A', 'Box', 0, 0, 1.5, '#e8e8e8', { sx: 1.2, sy: 0.8, sz: 0.1 }),
      TPL(tid('sp'), 'box', '球门B', 'Box', 0, 0, -1.5, '#e8e8e8', { sx: 1.2, sy: 0.8, sz: 0.1 }),
      // 篮球架
      TPL(tid('sp'), 'box', '篮球架A', 'Box', -2.5, 0, 0, '#ff8c00', { sx: 0.1, sy: 1.2, sz: 0.8, emissive: '#ff8c0022' }),
      // 游泳池
      TPL(tid('sp'), 'cylinder', '游泳池', 'Cylinder', 0, 0.05, 2.5, '#13c2c222', { sx: 1.5, sy: 0.1, sz: 0.8 }),
      // 泛光灯
      TPL(tid('sp'), 'lamppost', '泛光灯A', 'LampPost', -3, 0, -3, '#fff5e0', { emissive: '#fff5e055' }),
      TPL(tid('sp'), 'lamppost', '泛光灯B', 'LampPost', 3, 0, -3, '#fff5e0', { emissive: '#fff5e055' }),
      TPL(tid('sp'), 'lamppost', '泛光灯C', 'LampPost', -3, 0, 3, '#fff5e0', { emissive: '#fff5e055' }),
      TPL(tid('sp'), 'lamppost', '泛光灯D', 'LampPost', 3, 0, 3, '#fff5e0', { emissive: '#fff5e055' }),
      // 特效
      TPL(tid('sp'), 'scanLine', '扫描线', 'ScanLine', 0, 0.5, 0, '#13c2c2', { sy: 1.2 }),
      TPL(tid('sp'), 'glowSphere', '辉光球', 'GlowSphere', 0, 2, 0, '#fff5e0', { sx: 0.3, sy: 0.3, sz: 0.3 }),
      // 数据标签
      DATA_LABEL('sp', '运动场馆', 0, 4, -3, '#13c2c2', 32),
      DATA_LABEL('sp', '场内人数: 128', -3, 1.5, 2, '#4dffa6', 20),
      DATA_LABEL('sp', '温度: 24C', 3, 1.5, 2, '#4dffa6', 20),
    ]
  },

  // ==================== 17. 科研实验室 ====================
  {
    id: 'research-lab',
    name: '科研实验室',
    icon: 'fas fa-flask',
    category: '实验室',
    description: '冷蓝科技风科研实验室，发光分析仪+通风橱+培养箱+安全设备',
    sceneSettings: NIGHT_SCENE({ x: 5, y: 3.5, z: 5 }, { x: 0, y: 1.0, z: 0 }),
    objects: [
      TPL(tid('lb'), 'gridFloor', '地面网格', 'GridFloor', 0, 0.01, 0, '#13c2c2', { sx: 2, sz: 2 }),
      TPL(tid('lb'), 'plane', '实验室地板', 'Plane', 0, 0, 0, '#0d1a22', { sx: 4, sz: 4, emissive: '#040810' }),
      // 实验台
      TPL(tid('lb'), 'labbench', '实验台A', 'LabBench', -1.5, 0, 0, '#0d1a22', { emissive: '#040810' }),
      TPL(tid('lb'), 'labbench', '实验台B', 'LabBench', 1.5, 0, 0, '#0d1a22', { emissive: '#040810' }),
      TPL(tid('lb'), 'labbench', '实验台C', 'LabBench', 0, 0, -1.5, '#0d1a22', { emissive: '#040810' }),
      // 通风橱
      TPL(tid('lb'), 'fumehood', '通风橱A', 'FumeHood', -2.5, 0, -1.5, '#0d1f2a', { emissive: '#081018' }),
      TPL(tid('lb'), 'fumehood', '通风橱B', 'FumeHood', 2.5, 0, -1.5, '#0d1f2a', { emissive: '#081018' }),
      // 分析仪器
      TPL(tid('lb'), 'microscope', '显微镜A', 'Microscope', -1.8, 0, 0.3, '#13c2c2', { emissive: '#13c2c222' }),
      TPL(tid('lb'), 'box', '光谱仪', 'Box', 1.5, 0, -0.5, '#36cfc9', { sx: 0.6, sy: 0.5, sz: 0.4, emissive: '#36cfc922' }),
      TPL(tid('lb'), 'box', '离心机A', 'Box', -0.5, 0, -1, '#0d1f2a', { sx: 0.5, sy: 0.4, sz: 0.4, emissive: '#081018' }),
      // 培养/灭菌
      TPL(tid('lb'), 'box', '培养箱', 'Box', 2.5, 0, 0.5, '#0d2a1a', { sx: 0.5, sy: 0.6, sz: 0.4, emissive: '#051510' }),
      TPL(tid('lb'), 'box', '高压灭菌器', 'Box', 2.5, 0, 1.2, '#ff4d4f', { sx: 0.5, sy: 0.5, sz: 0.4, emissive: '#ff4d4f11' }),
      // 安全
      TPL(tid('lb'), 'fireextinguisher', '灭火器', 'FireExtinguisher', -2.5, 0, 1, '#ff4d4f', { emissive: '#ff4d4f33' }),
      TPL(tid('lb'), 'firstaidkit', '急救箱', 'FirstAidKit', 0, 0, 2.5, '#4dffa6', { emissive: '#4dffa644' }),
      TPL(tid('lb'), 'sinkstation', '洗手台', 'SinkStation', -1, 0, 2.5, '#e8e8e8', { sx: 0.8 }),
      // 特效
      TPL(tid('lb'), 'scanLine', '扫描线', 'ScanLine', 0, 0.5, 0, '#13c2c2', { sy: 1 }),
      TPL(tid('lb'), 'glowSphere', '辉光球', 'GlowSphere', 0, 1.5, 0, '#13c2c2', { sx: 0.2, sy: 0.2, sz: 0.2 }),
      // 数据标签
      DATA_LABEL('lb', '科研实验室', 0, 3, -2.5, '#13c2c2', 32),
      DATA_LABEL('lb', '室温: 22C', -3, 1, 0, '#4dffa6', 20),
      DATA_LABEL('lb', '负压: 正常', 3, 1, 0, '#4dffa6', 20),
    ]
  }
]

/**
 * 按分类获取模板列表
 */
export function getTemplatesByCategory(category) {
  if (!category) return SCENE_TEMPLATES
  return SCENE_TEMPLATES.filter(t => t.category === category)
}

/**
 * 获取所有模板分类
 */
export function getTemplateCategories() {
  return [...new Set(SCENE_TEMPLATES.map(t => t.category))]
}

/**
 * 通过ID查找模板
 */
export function getTemplateById(id) {
  return SCENE_TEMPLATES.find(t => t.id === id)
}
