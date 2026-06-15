/**
 * ISM 电力行业 SCADA 监控配置
 * 设备类型图标 & 颜色映射表 + 开发用 fallback mock 数据
 */

// ====== 设备类型 → SVG 组件标识 ======
export const DEVICE_TYPE_COMPONENT = {
  'transformer':      'electric_2',    // 变压器
  'switchgear':       'electric_5',    // 开关/刀闸（有开关状态）
  'meter-cabinet':    'electric_7',    // 仪表/互感器
  'pt-cabinet':       'electric_7',    // PT柜 = 互感器
  'bus-coupler':      'electric_3',    // 母联柜
  'capacitor':        'electric_1',    // 电容柜
  'feeder':           'electric_8',    // 馈线柜
  'ats':              'electric_4',    // ATS双电源
  'dc-panel':         'electric_8',    // 直流屏
  'ups':              'electric_8',    // UPS
  'comm-manager':     'electric_8',    // 通信管理机
  'sensor':           'electric_8',    // 传感器
}

// ====== 设备类型 Unicode 图标（fallback） ======
export const DEVICE_TYPE_ICONS = {
  'transformer':      '\u23F7',  // ⏷
  'switchgear':       '\u26A1',  // ⚡
  'meter-cabinet':    '\u2302',  // ⌂
  'pt-cabinet':       '\u238D',  // ⎍
  'bus-coupler':      '\u21C4',  // ⇄
  'capacitor':        '\u2395',  // ⎕
  'feeder':           '\u21E8',  // ⇨
  'ats':              '\u21C5',  // ⇅
  'dc-panel':         '\u2393',  // ⎓
  'ups':              '\u23FB',  // ⏻
  'comm-manager':     '\u2302',  // ⌂
  'sensor':           '\u2316',  // ⌖
}

// ====== 设备类型颜色 ======
export const DEVICE_TYPE_COLORS = {
  'transformer':   '#e74c3c',
  'switchgear':    '#3498db',
  'meter-cabinet': '#2ecc71',
  'pt-cabinet':    '#9b59b6',
  'bus-coupler':   '#f39c12',
  'capacitor':     '#1abc9c',
  'feeder':        '#2980b9',
  'ats':           '#e67e22',
  'dc-panel':      '#8e44ad',
  'ups':           '#16a085',
  'comm-manager':  '#7f8c8d',
  'sensor':        '#95a5a6',
}

// ====== Fallback mock 数据（API 无数据时使用） ======
export const FALLBACK_BUILDINGS = [
  {
    id: 'bldg-1',
    name: '1号配电楼',
    floors: [
      {
        id: 'floor-1-1',
        name: 'B1层·高压配电室',
        devices: [
          { uid: 'device-uuid-1', name: 'T1变压器', type: 'transformer', protocol: 'modbus', status: 'running' },
          { uid: 'device-uuid-2', name: '高压进线柜 AH01', type: 'switchgear', protocol: 'modbus', status: 'running' },
          { uid: 'device-uuid-3', name: '高压计量柜 AH02', type: 'meter-cabinet', protocol: 'dl645', status: 'running' },
          { uid: 'device-uuid-4', name: '高压PT柜 AH03', type: 'pt-cabinet', protocol: 'modbus', status: 'running' },
          { uid: 'device-uuid-5', name: '高压出线柜 AH04', type: 'switchgear', protocol: 'modbus', status: 'alarm' },
          { uid: 'device-uuid-6', name: '高压母联柜 AH05', type: 'bus-coupler', protocol: 'modbus', status: 'stopped' },
        ]
      },
      {
        id: 'floor-1-2',
        name: '1层·低压配电室',
        devices: [
          { uid: 'device-uuid-7', name: 'T2变压器', type: 'transformer', protocol: 'modbus', status: 'running' },
          { uid: 'device-uuid-8', name: '低压进线柜 AA01', type: 'switchgear', protocol: 'modbus', status: 'running' },
          { uid: 'device-uuid-9', name: '低压电容柜 AA02', type: 'capacitor', protocol: 'modbus', status: 'running' },
          { uid: 'device-uuid-10', name: '低压馈线柜 AA03', type: 'feeder', protocol: 'modbus', status: 'running' },
          { uid: 'device-uuid-11', name: '低压馈线柜 AA04', type: 'feeder', protocol: 'modbus', status: 'running' },
          { uid: 'device-uuid-12', name: '低压母联柜 AA05', type: 'bus-coupler', protocol: 'modbus', status: 'stopped' },
          { uid: 'device-uuid-13', name: 'ATS双电源柜 AT01', type: 'ats', protocol: 'modbus', status: 'running' },
        ]
      },
      {
        id: 'floor-1-3',
        name: '2层·中控室',
        devices: [
          { uid: 'device-uuid-14', name: '直流屏 DC01', type: 'dc-panel', protocol: 'modbus', status: 'running' },
          { uid: 'device-uuid-15', name: 'UPS电源 UPS01', type: 'ups', protocol: 'snmp', status: 'running' },
          { uid: 'device-uuid-16', name: '通信管理机 CM01', type: 'comm-manager', protocol: 'iec104', status: 'running' },
          { uid: 'device-uuid-17', name: '温湿度传感器 TH01', type: 'sensor', protocol: 'modbus', status: 'running' },
        ]
      }
    ]
  },
  {
    id: 'bldg-2',
    name: '2号配电楼',
    floors: [
      {
        id: 'floor-2-1',
        name: 'B1层·高压配电室',
        devices: [
          { uid: 'device-uuid-20', name: 'T3变压器', type: 'transformer', protocol: 'modbus', status: 'running' },
          { uid: 'device-uuid-21', name: '高压进线柜 BH01', type: 'switchgear', protocol: 'modbus', status: 'running' },
          { uid: 'device-uuid-22', name: '高压出线柜 BH02', type: 'switchgear', protocol: 'modbus', status: 'running' },
        ]
      },
      {
        id: 'floor-2-2',
        name: '1层·低压配电室',
        devices: [
          { uid: 'device-uuid-23', name: 'T4变压器', type: 'transformer', protocol: 'modbus', status: 'running' },
          { uid: 'device-uuid-24', name: '低压进线柜 BA01', type: 'switchgear', protocol: 'modbus', status: 'running' },
          { uid: 'device-uuid-25', name: '低压馈线柜 BA02', type: 'feeder', protocol: 'modbus', status: 'running' },
        ]
      }
    ]
  }
]

// 设备状态映射
export const STATUS_LABELS = {
  running: '运行中',
  stopped: '已停止',
  alarm: '告警',
  offline: '离线',
}

// 数据点中文名 → 英文 key 映射（用于 SCADA 数据查询）
export const DATA_KEY_MAP = {
  'power': ['有功功率', '功率', '总功率', '有功总功率'],
  'Ia': ['A相电流', 'Ia', 'A相电流Ia', '电流A'],
  'Ib': ['B相电流', 'Ib', 'B相电流Ib', '电流B'],
  'Ic': ['C相电流', 'Ic', 'C相电流Ic', '电流C'],
  'Ua': ['A相电压', 'Ua', 'A相电压Ua', '电压A'],
  'Ub': ['B相电压', 'Ub', 'B相电压Ub', '电压B'],
  'Uc': ['C相电压', 'Uc', 'C相电压Uc', '电压C'],
  'temp': ['温度', 'A相温度', 'B相温度', 'C相温度', 'Temperature', '环境温度'],
  'pf': ['功率因数', 'PF', '功率因素', '总功率因数'],
}
