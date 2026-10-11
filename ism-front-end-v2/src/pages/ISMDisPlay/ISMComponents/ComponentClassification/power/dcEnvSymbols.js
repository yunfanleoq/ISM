/**
 * 机房动环设备图。几何是原创线框，线条色和填充色可在组态里改。
 * viewBox 默认 0 0 80 96。
 */

function R(x, y, w, h, extra) {
  return ['rect', Object.assign({ x: x, y: y, width: w, height: h, rx: 2, fill: 'body' }, extra || {})]
}
function L(x1, y1, x2, y2, extra) {
  return ['line', Object.assign({ x1: x1, y1: y1, x2: x2, y2: y2 }, extra || {})]
}
function C(cx, cy, r, extra) {
  return ['circle', Object.assign({ cx: cx, cy: cy, r: r, fill: 'none' }, extra || {})]
}
function T(x, y, text, size) {
  return ['text', { x: x, y: y, text: text, size: size || 8 }]
}
function led(x, y) {
  return C(x, y, 2.3, { fill: 'stroke' })
}
function cabinet(x, y, w, h) {
  return [
    R(x, y, w, h, { rx: 3 }),
    L(x, y + 10, x + w, y + 10),
    L(x, y + h - 8, x + w, y + h - 8)
  ]
}

const DC_SYMBOLS = [
  {
    id: 'dc-utility',
    name: 'view-svg-dc-utility',
    text: 'component.DcEnv.utility',
    viewBox: '0 0 80 96',
    w: 80,
    h: 96,
    elements: [
      L(14, 18, 66, 18, { sw: 2.4 }),
      L(22, 18, 22, 30),
      L(40, 18, 40, 30),
      L(58, 18, 58, 30),
      L(22, 30, 58, 30),
      L(40, 30, 40, 46),
      R(24, 46, 32, 28, { rx: 3 }),
      T(40, 60, '市电', 9),
      L(40, 74, 40, 86),
      ['polygon', { points: '40,92 34,82 46,82', fill: 'stroke' }]
    ]
  },
  {
    id: 'dc-transformer',
    name: 'view-svg-dc-transformer',
    text: 'component.DcEnv.transformer',
    viewBox: '0 0 80 96',
    w: 80,
    h: 96,
    elements: [
      C(24, 16, 4, { fill: 'body' }),
      C(40, 16, 4, { fill: 'body' }),
      C(56, 16, 4, { fill: 'body' }),
      L(24, 20, 24, 28),
      L(40, 20, 40, 28),
      L(56, 20, 56, 28),
      R(16, 28, 48, 52, { rx: 3 }),
      L(16, 36, 10, 36), L(16, 44, 10, 44), L(16, 52, 10, 52), L(16, 60, 10, 60), L(16, 68, 10, 68),
      L(64, 36, 70, 36), L(64, 44, 70, 44), L(64, 52, 70, 52), L(64, 60, 70, 60), L(64, 68, 70, 68),
      T(40, 54, 'TR', 11)
    ]
  },
  {
    id: 'dc-lv-panel',
    name: 'view-svg-dc-lv-panel',
    text: 'component.DcEnv.lvPanel',
    viewBox: '0 0 80 96',
    w: 80,
    h: 96,
    elements: cabinet(14, 8, 52, 80).concat([
      led(58, 16),
      R(20, 26, 12, 28, { rx: 1 }),
      R(34, 26, 12, 28, { rx: 1 }),
      R(48, 26, 12, 28, { rx: 1 }),
      L(26, 32, 26, 46),
      L(40, 32, 40, 46),
      L(54, 32, 54, 46),
      C(26, 62, 3, { fill: 'none' }),
      C(40, 62, 3, { fill: 'none' }),
      C(54, 62, 3, { fill: 'none' }),
      T(40, 78, 'LV', 8)
    ])
  },
  {
    id: 'dc-ats',
    name: 'view-svg-dc-ats',
    text: 'component.DcEnv.ats',
    viewBox: '0 0 80 96',
    w: 80,
    h: 96,
    state: {
      default: 1,
      name: 'component.Power.state',
      options: [
        { value: 0, option: 'component.Power.sourceA' },
        { value: 1, option: 'component.Power.sourceB' }
      ]
    },
    elements: cabinet(16, 8, 48, 80).concat([
      led(56, 16),
      T(40, 28, 'ATS', 9),
      L(28, 40, 28, 52),
      L(52, 40, 52, 52),
      L(28, 70, 40, 70),
      L(40, 70, 40, 78),
      L(28, 52, 46, 64, { when: 0 }),
      L(52, 52, 34, 64, { when: 1 }),
      C(40, 66, 2.2, { fill: 'stroke' })
    ])
  },
  {
    id: 'dc-ups',
    name: 'view-svg-dc-ups',
    text: 'component.DcEnv.ups',
    viewBox: '0 0 80 96',
    w: 80,
    h: 96,
    elements: cabinet(16, 8, 48, 80).concat([
      led(56, 16),
      T(40, 28, 'UPS', 10),
      ['polyline', { points: '24,46 30,46 34,38 42,54 46,46 56,46', fill: 'none' }],
      R(26, 62, 8, 12, { rx: 1, fill: 'none' }),
      R(36, 62, 8, 12, { rx: 1, fill: 'none' }),
      R(46, 62, 8, 12, { rx: 1, fill: 'none' }),
      L(28, 66, 32, 66), L(38, 66, 42, 66), L(48, 66, 52, 66)
    ])
  },
  {
    id: 'dc-battery',
    name: 'view-svg-dc-battery',
    text: 'component.DcEnv.battery',
    viewBox: '0 0 80 96',
    w: 80,
    h: 96,
    elements: [
      R(10, 18, 60, 64, { rx: 3 }),
      R(18, 28, 10, 36, { rx: 1, fill: 'none' }),
      R(32, 28, 10, 36, { rx: 1, fill: 'none' }),
      R(46, 28, 10, 36, { rx: 1, fill: 'none' }),
      L(20, 40, 26, 40), L(23, 37, 23, 43),
      L(34, 40, 40, 40), L(37, 37, 37, 43),
      L(50, 40, 54, 40),
      T(40, 74, 'BAT', 8),
      L(34, 10, 34, 18), L(46, 10, 46, 18), L(34, 10, 46, 10)
    ]
  },
  {
    id: 'dc-rpp',
    name: 'view-svg-dc-rpp',
    text: 'component.DcEnv.rpp',
    viewBox: '0 0 80 96',
    w: 64,
    h: 96,
    elements: cabinet(22, 6, 36, 84).concat([
      led(50, 14),
      C(40, 32, 8, { fill: 'none' }),
      L(40, 32, 46, 28),
      R(28, 48, 24, 6, { rx: 1, fill: 'none' }),
      R(28, 58, 24, 6, { rx: 1, fill: 'none' }),
      R(28, 68, 24, 6, { rx: 1, fill: 'none' }),
      T(40, 82, 'RPP', 7)
    ])
  },
  {
    id: 'dc-pdu',
    name: 'view-svg-dc-pdu',
    text: 'component.DcEnv.pdu',
    viewBox: '0 0 80 96',
    w: 48,
    h: 96,
    elements: [
      R(28, 6, 24, 84, { rx: 3 }),
      T(40, 16, 'PDU', 6),
      C(40, 28, 3, { fill: 'none' }),
      C(40, 40, 3, { fill: 'none' }),
      C(40, 52, 3, { fill: 'none' }),
      C(40, 64, 3, { fill: 'none' }),
      C(40, 76, 3, { fill: 'none' }),
      L(37, 28, 43, 28), L(37, 40, 43, 40), L(37, 52, 43, 52), L(37, 64, 43, 64), L(37, 76, 43, 76)
    ]
  },
  {
    id: 'dc-generator',
    name: 'view-svg-dc-generator',
    text: 'component.DcEnv.generator',
    viewBox: '0 0 140 80',
    w: 140,
    h: 80,
    elements: [
      L(10, 64, 130, 64, { sw: 2.4 }),
      L(22, 64, 22, 72),
      L(118, 64, 118, 72),
      R(12, 30, 22, 34, { rx: 2 }),
      L(17, 34, 17, 60),
      L(23, 34, 23, 60),
      L(29, 34, 29, 60),
      R(34, 34, 46, 30, { rx: 2 }),
      C(100, 48, 18, { fill: 'body' }),
      T(100, 48, 'G', 14),
      L(50, 34, 50, 16),
      L(50, 16, 68, 16),
      C(76, 14, 5, { fill: 'none' }),
      C(88, 10, 3.5, { fill: 'none' })
    ]
  },
  {
    id: 'dc-busway',
    name: 'view-svg-dc-busway',
    text: 'component.DcEnv.busway',
    viewBox: '0 0 120 48',
    w: 120,
    h: 48,
    elements: [
      R(6, 16, 108, 16, { rx: 2 }),
      L(30, 16, 30, 32),
      L(60, 16, 60, 32),
      L(90, 16, 90, 32),
      L(30, 32, 30, 42),
      L(60, 32, 60, 42),
      L(90, 32, 90, 42)
    ]
  },
  {
    id: 'dc-crac',
    name: 'view-svg-dc-crac',
    text: 'component.DcEnv.crac',
    viewBox: '0 0 80 96',
    w: 80,
    h: 96,
    elements: cabinet(16, 8, 48, 80).concat([
      led(56, 16),
      C(40, 42, 14, { fill: 'none' }),
      L(40, 30, 40, 54),
      L(28, 42, 52, 42),
      L(31, 33, 49, 51),
      L(31, 51, 49, 33),
      L(24, 68, 56, 68),
      L(24, 74, 56, 74),
      ['polyline', { points: '22,64 28,60 34,64 40,60 46,64 52,60 58,64', fill: 'none' }]
    ])
  },
  {
    id: 'dc-inrow',
    name: 'view-svg-dc-inrow',
    text: 'component.DcEnv.inrow',
    viewBox: '0 0 80 96',
    w: 56,
    h: 96,
    elements: cabinet(26, 6, 28, 84).concat([
      led(48, 14),
      C(40, 36, 8, { fill: 'none' }),
      L(40, 30, 40, 42),
      L(34, 36, 46, 36),
      ['polyline', { points: '32,58 40,50 48,58', fill: 'none' }],
      ['polyline', { points: '32,66 40,58 48,66', fill: 'none' }],
      L(32, 78, 48, 78)
    ])
  },
  {
    id: 'dc-th',
    name: 'view-svg-dc-th',
    text: 'component.DcEnv.thSensor',
    viewBox: '0 0 80 96',
    w: 72,
    h: 80,
    elements: [
      R(8, 16, 64, 64, { rx: 6 }),
      led(62, 26),
      R(18, 28, 8, 36, { rx: 4, fill: 'none' }),
      C(22, 58, 4, { fill: 'stroke' }),
      L(22, 36, 22, 52),
      ['path', { d: 'M46,34 C46,28 58,28 58,36 C58,42 46,44 46,52 C46,58 58,58 58,52', fill: 'none' }],
      T(40, 70, 'T/H', 8)
    ]
  },
  {
    id: 'dc-leak',
    name: 'view-svg-dc-leak',
    text: 'component.DcEnv.leak',
    viewBox: '0 0 80 96',
    w: 80,
    h: 72,
    elements: [
      ['path', { d: 'M40,10 C40,10 22,36 22,48 A18,18 0 1 0 58,48 C58,36 40,10 40,10 Z', fill: 'body' }],
      ['polyline', { points: '8,78 20,70 32,78 44,70 56,78 72,70', fill: 'none' }],
      L(40, 66, 40, 74)
    ]
  },
  {
    id: 'dc-smoke',
    name: 'view-svg-dc-smoke',
    text: 'component.DcEnv.smoke',
    viewBox: '0 0 80 96',
    w: 72,
    h: 72,
    elements: [
      R(16, 10, 48, 10, { rx: 2 }),
      C(40, 46, 22, { fill: 'body' }),
      C(40, 46, 10, { fill: 'none' }),
      C(40, 46, 3, { fill: 'stroke' }),
      L(24, 34, 18, 28),
      L(56, 34, 62, 28),
      L(24, 58, 18, 66),
      L(56, 58, 62, 66)
    ]
  },
  {
    id: 'dc-access',
    name: 'view-svg-dc-access',
    text: 'component.DcEnv.access',
    viewBox: '0 0 80 96',
    w: 72,
    h: 96,
    elements: [
      R(14, 10, 36, 76, { rx: 2 }),
      L(32, 10, 32, 86),
      C(26, 48, 2, { fill: 'stroke' }),
      R(52, 28, 16, 24, { rx: 2, fill: 'body' }),
      L(56, 36, 64, 36),
      L(56, 42, 62, 42),
      led(60, 22)
    ]
  },
  {
    id: 'dc-camera',
    name: 'view-svg-dc-camera',
    text: 'component.DcEnv.camera',
    viewBox: '0 0 80 96',
    w: 80,
    h: 64,
    elements: [
      R(8, 28, 40, 28, { rx: 4 }),
      ['polygon', { points: '48,34 70,24 70,60 48,50', fill: 'body' }],
      C(26, 42, 8, { fill: 'none' }),
      C(26, 42, 3, { fill: 'stroke' }),
      L(28, 18, 28, 28),
      R(18, 12, 20, 8, { rx: 2, fill: 'body' })
    ]
  },
  {
    id: 'dc-fresh-air',
    name: 'view-svg-dc-fresh-air',
    text: 'component.DcEnv.freshAir',
    viewBox: '0 0 80 96',
    w: 80,
    h: 80,
    elements: [
      C(40, 48, 26, { fill: 'body' }),
      C(40, 48, 6, { fill: 'stroke' }),
      ['path', { d: 'M40,28 C52,36 52,44 40,48 C28,44 28,36 40,28 Z', fill: 'none' }],
      ['path', { d: 'M60,48 C52,60 44,60 40,48 C44,36 52,36 60,48 Z', fill: 'none' }],
      ['path', { d: 'M40,68 C28,60 28,52 40,48 C52,52 52,60 40,68 Z', fill: 'none' }],
      ['polyline', { points: '12,18 24,18 18,12', fill: 'none' }],
      L(8, 18, 20, 18)
    ]
  },
  {
    id: 'dc-rack',
    name: 'view-svg-dc-rack',
    text: 'component.DcEnv.rack',
    viewBox: '0 0 80 96',
    w: 64,
    h: 96,
    elements: [
      R(18, 6, 44, 84, { rx: 2 }),
      L(24, 6, 24, 90),
      L(56, 6, 56, 90),
      L(24, 18, 56, 18),
      L(24, 30, 56, 30),
      L(24, 42, 56, 42),
      L(24, 54, 56, 54),
      L(24, 66, 56, 66),
      L(24, 78, 56, 78),
      C(22, 22, 1.3, { fill: 'stroke' }),
      C(22, 46, 1.3, { fill: 'stroke' }),
      C(22, 70, 1.3, { fill: 'stroke' }),
      led(52, 12)
    ]
  },
  {
    id: 'dc-cold-aisle',
    name: 'view-svg-dc-cold-aisle',
    text: 'component.DcEnv.coldAisle',
    viewBox: '0 0 140 96',
    w: 140,
    h: 96,
    elements: [
      R(8, 20, 36, 64, { rx: 2 }),
      R(96, 20, 36, 64, { rx: 2 }),
      L(14, 32, 38, 32), L(14, 44, 38, 44), L(14, 56, 38, 56), L(14, 68, 38, 68),
      L(102, 32, 126, 32), L(102, 44, 126, 44), L(102, 56, 126, 56), L(102, 68, 126, 68),
      ['path', { d: 'M8,20 L44,20 L70,8 L96,20 L132,20', fill: 'none' }],
      L(70, 8, 70, 20),
      T(70, 52, 'COLD', 8)
    ]
  },
  {
    id: 'dc-tray',
    name: 'view-svg-dc-tray',
    text: 'component.DcEnv.tray',
    viewBox: '0 0 120 48',
    w: 120,
    h: 40,
    elements: [
      L(6, 10, 114, 10),
      L(6, 34, 114, 34),
      L(18, 10, 18, 34),
      L(42, 10, 42, 34),
      L(66, 10, 66, 34),
      L(90, 10, 90, 34),
      L(6, 10, 6, 34),
      L(114, 10, 114, 34)
    ]
  },
  {
    id: 'dc-meter',
    name: 'view-svg-dc-meter',
    text: 'component.DcEnv.meter',
    viewBox: '0 0 80 96',
    w: 72,
    h: 80,
    elements: [
      R(12, 12, 56, 72, { rx: 4 }),
      C(40, 46, 20, { fill: 'none' }),
      L(40, 46, 52, 36),
      L(26, 46, 28, 46),
      L(52, 46, 54, 46),
      L(40, 30, 40, 32),
      T(40, 74, 'kWh', 8),
      led(60, 20)
    ]
  }
]

module.exports = {
  DC_SYMBOLS: DC_SYMBOLS
}
