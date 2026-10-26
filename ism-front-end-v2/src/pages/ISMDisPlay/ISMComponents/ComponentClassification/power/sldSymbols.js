/**
 * 机房低压配电一次图。原创几何，对应常见单线图画法，不是从标准文本描下来的。
 * 合闸 when:1 ，分闸 when:0 。默认合闸。
 */

function L(x1, y1, x2, y2, extra) {
  return ['line', Object.assign({ x1: x1, y1: y1, x2: x2, y2: y2 }, extra || {})]
}
function C(cx, cy, r, extra) {
  return ['circle', Object.assign({ cx: cx, cy: cy, r: r, fill: 'none' }, extra || {})]
}
function R(x, y, w, h, extra) {
  return ['rect', Object.assign({ x: x, y: y, width: w, height: h, rx: 0, fill: 'none' }, extra || {})]
}
function T(x, y, text, size) {
  return ['text', { x: x, y: y, text: text, size: size || 9 }]
}

function spine() {
  return [
    C(32, 5, 2.1, { fill: 'stroke' }),
    C(32, 83, 2.1, { fill: 'stroke' }),
    L(32, 5, 32, 24),
    L(32, 64, 32, 83)
  ]
}

const OPEN_CLOSE = {
  default: 1,
  name: 'component.Power.state',
  options: [
    { value: 0, option: 'component.Power.open' },
    { value: 1, option: 'component.Power.closed' }
  ]
}

const SLD_SYMBOLS = [
  {
    id: 'sld-busbar',
    name: 'view-svg-sld-busbar',
    text: 'component.Sld.busbar',
    viewBox: '0 0 140 40',
    w: 140,
    h: 40,
    elements: [
      R(8, 16, 124, 6, { rx: 1, fill: 'stroke' }),
      L(32, 22, 32, 36),
      L(70, 22, 70, 36),
      L(108, 22, 108, 36),
      C(32, 36, 2, { fill: 'stroke' }),
      C(70, 36, 2, { fill: 'stroke' }),
      C(108, 36, 2, { fill: 'stroke' })
    ]
  },
  {
    id: 'sld-breaker',
    name: 'view-svg-sld-breaker',
    text: 'component.Sld.breaker',
    viewBox: '0 0 64 88',
    w: 64,
    h: 88,
    state: OPEN_CLOSE,
    elements: spine().concat([
      R(22, 28, 20, 28),
      L(26, 42, 38, 42, { when: 1 }),
      L(26, 34, 38, 50, { when: 0 })
    ])
  },
  {
    id: 'sld-disconnector',
    name: 'view-svg-sld-disconnector',
    text: 'component.Sld.disconnector',
    viewBox: '0 0 64 88',
    w: 64,
    h: 88,
    state: OPEN_CLOSE,
    elements: [
      C(32, 5, 2.1, { fill: 'stroke' }),
      C(32, 83, 2.1, { fill: 'stroke' }),
      L(32, 5, 32, 28),
      L(32, 60, 32, 83),
      C(32, 30, 2.4, { fill: 'none' }),
      C(32, 58, 2.4, { fill: 'none' }),
      L(32, 32, 32, 56, { when: 1 }),
      L(32, 32, 48, 50, { when: 0 })
    ]
  },
  {
    id: 'sld-load-switch',
    name: 'view-svg-sld-load-switch',
    text: 'component.Sld.loadSwitch',
    viewBox: '0 0 64 88',
    w: 64,
    h: 88,
    state: OPEN_CLOSE,
    elements: [
      C(32, 5, 2.1, { fill: 'stroke' }),
      C(32, 83, 2.1, { fill: 'stroke' }),
      L(32, 5, 32, 26),
      L(32, 62, 32, 83),
      C(32, 28, 2.2, { fill: 'none' }),
      C(32, 60, 2.2, { fill: 'none' }),
      L(32, 30, 32, 58, { when: 1 }),
      L(32, 30, 46, 50, { when: 0 }),
      ['path', { d: 'M40,34 A10,10 0 0 1 40,54', fill: 'none' }]
    ]
  },
  {
    id: 'sld-earth',
    name: 'view-svg-sld-earth',
    text: 'component.Sld.earth',
    viewBox: '0 0 64 88',
    w: 64,
    h: 72,
    elements: [
      C(32, 6, 2.1, { fill: 'stroke' }),
      L(32, 6, 32, 28),
      L(16, 28, 48, 28, { sw: 2.2 }),
      L(20, 36, 44, 36),
      L(24, 44, 40, 44),
      L(28, 52, 36, 52)
    ]
  },
  {
    id: 'sld-transformer',
    name: 'view-svg-sld-transformer',
    text: 'component.Sld.transformer',
    viewBox: '0 0 64 88',
    w: 64,
    h: 88,
    elements: [
      C(32, 5, 2.1, { fill: 'stroke' }),
      C(32, 83, 2.1, { fill: 'stroke' }),
      L(32, 5, 32, 22),
      L(32, 66, 32, 83),
      C(32, 34, 12, { fill: 'none' }),
      C(32, 54, 12, { fill: 'none' })
    ]
  },
  {
    id: 'sld-ct',
    name: 'view-svg-sld-ct',
    text: 'component.Sld.ct',
    viewBox: '0 0 64 88',
    w: 64,
    h: 88,
    elements: spine().concat([
      C(32, 44, 14, { fill: 'none' }),
      L(32, 24, 32, 64)
    ])
  },
  {
    id: 'sld-vt',
    name: 'view-svg-sld-vt',
    text: 'component.Sld.vt',
    viewBox: '0 0 64 88',
    w: 64,
    h: 88,
    elements: [
      C(32, 5, 2.1, { fill: 'stroke' }),
      L(32, 5, 32, 24),
      C(32, 40, 14, { fill: 'none' }),
      C(32, 56, 14, { fill: 'none' }),
      L(18, 78, 46, 78, { sw: 2 })
    ]
  },
  {
    id: 'sld-arrester',
    name: 'view-svg-sld-arrester',
    text: 'component.Sld.arrester',
    viewBox: '0 0 64 88',
    w: 64,
    h: 88,
    elements: [
      C(32, 5, 2.1, { fill: 'stroke' }),
      L(32, 5, 32, 22),
      R(24, 22, 16, 28, { fill: 'none' }),
      L(27, 30, 37, 26),
      L(27, 38, 37, 34),
      L(27, 46, 37, 42),
      L(32, 50, 32, 62),
      L(18, 62, 46, 62, { sw: 2 }),
      L(22, 68, 42, 68),
      L(26, 74, 38, 74)
    ]
  },
  {
    id: 'sld-fuse',
    name: 'view-svg-sld-fuse',
    text: 'component.Sld.fuse',
    viewBox: '0 0 64 88',
    w: 64,
    h: 88,
    elements: spine().concat([
      R(26, 30, 12, 28, { rx: 6 }),
      L(32, 34, 32, 54)
    ])
  },
  {
    id: 'sld-generator',
    name: 'view-svg-sld-generator',
    text: 'component.Sld.generator',
    viewBox: '0 0 64 88',
    w: 64,
    h: 88,
    elements: [
      C(32, 83, 2.1, { fill: 'stroke' }),
      L(32, 58, 32, 83),
      C(32, 40, 18, { fill: 'none' }),
      T(32, 40, 'G', 14)
    ]
  },
  {
    id: 'sld-ups',
    name: 'view-svg-sld-ups',
    text: 'component.Sld.ups',
    viewBox: '0 0 64 88',
    w: 64,
    h: 88,
    elements: spine().concat([
      R(18, 28, 28, 32, { rx: 2 }),
      T(32, 44, 'UPS', 8)
    ])
  },
  {
    id: 'sld-battery',
    name: 'view-svg-sld-battery',
    text: 'component.Sld.battery',
    viewBox: '0 0 64 88',
    w: 64,
    h: 88,
    elements: [
      C(32, 5, 2.1, { fill: 'stroke' }),
      C(32, 83, 2.1, { fill: 'stroke' }),
      L(32, 5, 32, 24),
      L(32, 64, 32, 83),
      L(18, 28, 46, 28, { sw: 2.4 }),
      L(22, 38, 42, 38),
      L(18, 48, 46, 48, { sw: 2.4 }),
      L(22, 58, 42, 58),
      T(50, 30, '+', 8),
      T(50, 58, '-', 8)
    ]
  },
  {
    id: 'sld-ats',
    name: 'view-svg-sld-ats',
    text: 'component.Sld.ats',
    viewBox: '0 0 80 88',
    w: 80,
    h: 88,
    state: {
      default: 0,
      name: 'component.Power.state',
      options: [
        { value: 0, option: 'component.Power.sourceA' },
        { value: 1, option: 'component.Power.sourceB' }
      ]
    },
    elements: [
      C(16, 6, 2.1, { fill: 'stroke' }),
      C(64, 6, 2.1, { fill: 'stroke' }),
      C(40, 82, 2.1, { fill: 'stroke' }),
      L(16, 6, 16, 28),
      L(64, 6, 64, 28),
      L(40, 58, 40, 82),
      C(40, 56, 2.4, { fill: 'stroke' }),
      L(16, 28, 40, 54, { when: 0 }),
      L(64, 28, 40, 54, { when: 1 }),
      T(16, 20, 'A', 8),
      T(64, 20, 'B', 8)
    ]
  },
  {
    id: 'sld-in',
    name: 'view-svg-sld-in',
    text: 'component.Sld.feederIn',
    viewBox: '0 0 64 88',
    w: 64,
    h: 88,
    elements: [
      C(32, 83, 2.1, { fill: 'stroke' }),
      L(32, 18, 32, 83),
      ['polygon', { points: '32,8 24,24 40,24', fill: 'stroke' }],
      T(32, 36, 'IN', 8)
    ]
  },
  {
    id: 'sld-out',
    name: 'view-svg-sld-out',
    text: 'component.Sld.feederOut',
    viewBox: '0 0 64 88',
    w: 64,
    h: 88,
    elements: [
      C(32, 5, 2.1, { fill: 'stroke' }),
      L(32, 5, 32, 58),
      ['polygon', { points: '32,80 24,60 40,60', fill: 'stroke' }],
      T(32, 40, 'OUT', 8)
    ]
  },
  {
    id: 'sld-capacitor',
    name: 'view-svg-sld-capacitor',
    text: 'component.Sld.capacitor',
    viewBox: '0 0 64 88',
    w: 64,
    h: 88,
    elements: [
      C(32, 5, 2.1, { fill: 'stroke' }),
      L(32, 5, 32, 34),
      L(16, 34, 48, 34, { sw: 2.4 }),
      L(16, 46, 48, 46, { sw: 2.4 }),
      L(32, 46, 32, 70),
      L(18, 70, 46, 70, { sw: 2 }),
      L(22, 76, 42, 76),
      L(26, 82, 38, 82)
    ]
  },
  {
    id: 'sld-motor',
    name: 'view-svg-sld-motor',
    text: 'component.Sld.motor',
    viewBox: '0 0 64 88',
    w: 64,
    h: 88,
    elements: [
      C(32, 5, 2.1, { fill: 'stroke' }),
      L(32, 5, 32, 28),
      C(32, 50, 20, { fill: 'none' }),
      T(32, 50, 'M', 14)
    ]
  }
]

module.exports = {
  SLD_SYMBOLS: SLD_SYMBOLS
}
