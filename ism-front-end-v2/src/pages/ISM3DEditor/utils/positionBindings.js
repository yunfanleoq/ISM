export const POSITION_BINDING_DEFAULTS = {
  dataID: '',
  dataName: '',
  dataUnit: '',
  isBandDevice: false,
  deviceSN: '',
  DeviceName: '',
  wsKey: '',
  transform: 'direct',
  scale: 1,
  offset: 0
}

export const POSITION_BINDING_KEYS = Object.keys(POSITION_BINDING_DEFAULTS)
export const POSITION_BINDING_AXES = ['x', 'y', 'z']

export function createDefaultPositionBinding() {
  return Object.assign({}, POSITION_BINDING_DEFAULTS)
}

export function createDefaultPositionBindings() {
  return {
    x: createDefaultPositionBinding(),
    y: createDefaultPositionBinding(),
    z: createDefaultPositionBinding()
  }
}

export function normalizePositionBinding(binding) {
  const next = binding || {}
  POSITION_BINDING_KEYS.forEach((key) => {
    if (next[key] === undefined) {
      next[key] = POSITION_BINDING_DEFAULTS[key]
    }
  })
  return next
}

export function ensurePositionBindings(target) {
  if (!target.positionBindings) {
    target.positionBindings = {}
  }
  POSITION_BINDING_AXES.forEach((axis) => {
    target.positionBindings[axis] = normalizePositionBinding(target.positionBindings[axis])
  })
  return target.positionBindings
}
