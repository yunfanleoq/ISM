export function defaultGridSettings() {
  return {
    size: 20,
    divisions: 20,
    colorCenterLine: '#111111',
    colorGrid: '#cccccc',
    backgroundColor: '#ffffff',
    backgroundColor2: '#dbefff',
    backgroundImage: '',
    backgroundMode: 'solid',
    sceneEnvironmentPreset: 'clearSky',
    environmentPreset: 'sky',
    skyboxEnabled: true,
    skyboxPreset: 'horizon',
    skyboxImage: '',
    skyboxHdri: '',
    groundEnabled: true,
    groundStyle: 'slate',
    sectionLength: 220,
    sectionWidth: 36,
    sectionHeight: 18,
    tunnelRadius: 18,
    lightingPreset: 'day',
    floorReflection: 'none',
    enhanceDepth: false,
    modelOptimize: false,
    lightSettings: {
      envIntensity: 1.0,
      ambientColor: '#ffffff',
      ambientIntensity: 1.0,
      directionalColor: '#ffffff',
      directionalIntensity: 0.9
    }
  }
}

export function mergeSceneSettings(base, incoming) {
  const defaults = defaultGridSettings()
  const source = incoming || {}
  const target = Object.assign({}, defaults, base || {})

  if (source.backgroundColor) target.backgroundColor = source.backgroundColor
  if (source.backgroundColor2) target.backgroundColor2 = source.backgroundColor2
  if (source.backgroundImage !== undefined) target.backgroundImage = source.backgroundImage
  if (source.backgroundMode) target.backgroundMode = source.backgroundMode
  if (source.sceneEnvironmentPreset) target.sceneEnvironmentPreset = source.sceneEnvironmentPreset
  if (source.environmentPreset) target.environmentPreset = source.environmentPreset
  if (source.skyboxEnabled !== undefined) target.skyboxEnabled = source.skyboxEnabled
  if (source.skyboxPreset) target.skyboxPreset = source.skyboxPreset
  if (source.skyboxImage !== undefined) target.skyboxImage = source.skyboxImage
  if (source.skyboxHdri !== undefined) target.skyboxHdri = source.skyboxHdri
  if (source.groundEnabled !== undefined) target.groundEnabled = source.groundEnabled
  if (source.groundStyle) target.groundStyle = source.groundStyle
  if (source.sectionLength !== undefined) target.sectionLength = source.sectionLength
  if (source.sectionWidth !== undefined) target.sectionWidth = source.sectionWidth
  if (source.sectionHeight !== undefined) target.sectionHeight = source.sectionHeight
  if (source.tunnelRadius !== undefined) target.tunnelRadius = source.tunnelRadius
  if (source.lightingPreset) target.lightingPreset = source.lightingPreset
  if (source.floorReflection) target.floorReflection = source.floorReflection
  if (source.enhanceDepth !== undefined) target.enhanceDepth = source.enhanceDepth
  if (source.modelOptimize !== undefined) target.modelOptimize = source.modelOptimize
  if (source.gridSize !== undefined) target.size = source.gridSize
  if (source.gridDivisions !== undefined) target.divisions = source.gridDivisions
  if (source.gridColorCenterLine) target.colorCenterLine = source.gridColorCenterLine
  if (source.gridColorGrid) target.colorGrid = source.gridColorGrid
  if (source.lightSettings) target.lightSettings = source.lightSettings

  return target
}

export function getEnvironmentPresetColors(preset) {
  const colors = {
    sky: { backgroundColor: '#e8f5ff', backgroundColor2: '#b7ddff' },
    sunset: { backgroundColor: '#fff0d9', backgroundColor2: '#ffb36f' },
    ocean: { backgroundColor: '#e5fbff', backgroundColor2: '#8fd8ef' },
    forest: { backgroundColor: '#ecf8ed', backgroundColor2: '#9fd6ad' },
    twilight: { backgroundColor: '#f0ecff', backgroundColor2: '#b6a7ff' },
    night: { backgroundColor: '#101827', backgroundColor2: '#273452' }
  }
  return colors[preset] || colors.sky
}
