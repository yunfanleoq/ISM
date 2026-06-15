export function createGLTFAnimationCondition(source = {}) {
  return {
    isBandDevice: !!source.isBandDevice,
    deviceSN: source.deviceSN || '',
    DeviceName: source.DeviceName || '',
    dataID: source.dataID || '',
    dataName: source.dataName || '',
    operator: source.operator || '',
    OperatorValue: source.OperatorValue !== undefined ? source.OperatorValue : '',
    OperatorMaxValue: source.OperatorMaxValue !== undefined ? source.OperatorMaxValue : ''
  }
}

export function cloneGLTFAnimationGroups(groups) {
  return Array.isArray(groups)
    ? groups.map((group, index) => normalizeGLTFAnimationGroup(group, index))
    : []
}

export function normalizeGLTFAnimationGroup(group = {}, index = 0) {
  const names = Array.isArray(group.animationNames)
    ? group.animationNames.slice()
    : (group.animationName ? [group.animationName] : [])
  return {
    id: group.id || ('gltf_anim_group_' + index),
    name: group.name || ('动画组' + (index + 1)),
    playing: !!group.playing,
    animationNames: names,
    speed: group.speed !== undefined ? group.speed : 1,
    loop: group.loop !== false,
    conditionEnabled: !!group.conditionEnabled,
    condition: createGLTFAnimationCondition(group.condition || {})
  }
}

export function legacyToGLTFAnimationGroup(obj = {}, index = 0) {
  const names = Array.isArray(obj.gltfAnimationNames)
    ? obj.gltfAnimationNames.slice()
    : (obj.gltfAnimationName ? [obj.gltfAnimationName] : [])
  return normalizeGLTFAnimationGroup({
    id: 'gltf_anim_group_legacy_' + index,
    name: '动画组' + (index + 1),
    playing: !!obj.gltfAnimationPlaying,
    animationNames: names,
    speed: obj.gltfAnimationSpeed !== undefined ? obj.gltfAnimationSpeed : 1,
    loop: obj.gltfAnimationLoop !== false,
    conditionEnabled: !!obj.gltfAnimationConditionEnabled,
    condition: obj.gltfAnimationCondition || {}
  }, index)
}

export function ensureGLTFAnimationGroups(obj, options = {}) {
  if (!obj) return []
  const createDefault = options.createDefault !== false
  const hasGroups = Array.isArray(obj.gltfAnimationGroups) && obj.gltfAnimationGroups.length > 0
  let groups = hasGroups ? cloneGLTFAnimationGroups(obj.gltfAnimationGroups) : []
  if (!groups.length && createDefault) {
    groups = [legacyToGLTFAnimationGroup(obj, 0)]
  }
  if (options.defaultAnimationKey && groups.length && !groups[0].animationNames.length) {
    groups[0].animationNames = [options.defaultAnimationKey]
  }
  return groups
}

export function syncLegacyGLTFAnimationFields(obj, groups) {
  if (!obj) return
  const first = Array.isArray(groups) && groups.length ? groups[0] : null
  const names = first && Array.isArray(first.animationNames) ? first.animationNames.slice() : []
  obj.gltfAnimationPlaying = first ? !!first.playing : !!obj.gltfAnimationPlaying
  obj.gltfAnimationNames = names
  obj.gltfAnimationName = names[0] || ''
  obj.gltfAnimationSpeed = first && first.speed !== undefined ? first.speed : 1
  obj.gltfAnimationLoop = first ? first.loop !== false : true
  obj.gltfAnimationConditionEnabled = first ? !!first.conditionEnabled : false
  obj.gltfAnimationCondition = first ? createGLTFAnimationCondition(first.condition || {}) : createGLTFAnimationCondition()
}

export function remapGLTFAnimationGroupNames(groups, entries, actions, defaultAnimationKey) {
  const list = Array.isArray(groups) ? groups : []
  return list.map((group, index) => {
    const nextGroup = normalizeGLTFAnimationGroup(group, index)
    const nextNames = nextGroup.animationNames.map(savedName => {
      if (actions && actions[savedName]) return savedName
      const matchedEntry = entries.find(entry => entry.name === savedName || entry.label === savedName || entry.key === savedName)
      return matchedEntry ? matchedEntry.key : ''
    }).filter(Boolean)
    if (!nextNames.length && index === 0 && defaultAnimationKey) nextNames.push(defaultAnimationKey)
    nextGroup.animationNames = nextNames
    return nextGroup
  })
}
