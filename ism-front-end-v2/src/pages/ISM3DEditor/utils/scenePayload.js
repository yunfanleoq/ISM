import { defaultGridSettings } from './sceneSettings'

function createDefaultSceneExtras() {
  return {
    environment: 'studio',
    fogEnabled: false,
    fogColor: '#d8e4f0',
    defaultCameraId: '',
    cameraViews: [],
    eventFlows: [],
    timelineEnabled: false,
    timelineSpeed: 1,
    timelineLoop: true,
    timelineDuration: 0,
    timelineTracks: [],
    resourceLibrary: []
  }
}

export function createEmptyScenePayload() {
  return {
    objects: [],
    sceneSettings: defaultGridSettings(),
    sceneExtras: createDefaultSceneExtras()
  }
}

/**
 * Parse scene payloads saved by different historical versions.
 *
 * Some records store `components` as JSON once, while older paths stored a
 * JSON string inside another JSON string. Keep reads compatible with both.
 */
export function parseScenePayload(raw) {
  if (raw === null || raw === undefined || raw === '') {
    throw new Error('Empty scene payload')
  }

  let data = raw
  for (let i = 0; i < 2 && typeof data === 'string'; i++) {
    data = JSON.parse(data)
  }

  if (!data || typeof data !== 'object') {
    throw new Error('Invalid scene payload')
  }

  return data
}

export function isBlankLayerPayload(data) {
  if (data === null || data === undefined || data === '') return true
  if (Array.isArray(data)) return data.length === 0
  if (!data || typeof data !== 'object') return false

  const keys = Object.keys(data)
  if (keys.length === 0) return true
  if (keys.every(key => key === 'layer')) return true
  if (Array.isArray(data.cells) && data.cells.length === 0 && keys.every(key => key === 'cells')) return true
  if (data.components && typeof data.components === 'object' && Array.isArray(data.components.cells) && data.components.cells.length === 0 && keys.every(key => key === 'components' || key === 'layer')) return true
  if (Array.isArray(data.components) && data.components.length === 0 && keys.every(key => key === 'components' || key === 'layer')) return true

  return false
}

export function parseScenePayloadOrEmpty(raw) {
  if (raw === null || raw === undefined || raw === '') {
    return createEmptyScenePayload()
  }

  const data = parseScenePayload(raw)
  return isBlankLayerPayload(data) ? createEmptyScenePayload() : data
}
