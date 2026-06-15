import Vue from 'vue'

const registry = {}
const toolboxGroups = []

function addContext(context, groupTitle) {
  const group = {
    title: groupTitle,
    items: []
  }
  context.keys().forEach((filePath) => {
    const mod = context(filePath)
    const component = mod && mod.default ? mod.default : mod
    if (!component || !component.name) return
    registry[component.name] = component
    if (!Vue.options.components[component.name]) {
      Vue.component(component.name, component)
    }
    try {
      const base = component.data && typeof component.data === 'function' ? component.data().base : null
      if (base && base.info) {
        const normalizedBase = JSON.parse(JSON.stringify(base))
        normalizedBase.info.type = component.name
        group.items.push(normalizedBase)
      }
    } catch (err) {
      // ignore component toolbox extraction failures
    }
  })
  if (group.items.length > 0) {
    toolboxGroups.push(group)
  }
}

addContext(require.context('./ISMComponents/standard/', true, /\.vue$/), 'displayConfig.ToolBox.Base.title')
addContext(require.context('./ISMComponents/video/', true, /\.vue$/), 'displayConfig.ToolBox.Video.title')
addContext(require.context('./ISMComponents/login/', true, /\.vue$/), 'displayConfig.ToolBox.login.title')
addContext(require.context('./ISMComponents/canvas/', true, /\.vue$/), 'displayConfig.ToolBox.Diagram.title')
addContext(require.context('./ISMComponents/charts/', true, /\.vue$/), 'displayConfig.ToolBox.Charts.title')
addContext(require.context('./ISMComponents/bigScreen/', true, /\.vue$/), 'displayConfig.ToolBox.bigScreen.Container')
addContext(require.context('./ISMComponents/svg/arrows/', true, /\.vue$/), 'displayConfig.ToolBox.Arrows')
addContext(require.context('./ISMComponents/ComponentClassification/electric/', true, /\.vue$/), 'displayConfig.ToolBox.Electric')
addContext(require.context('./ISMComponents/Images/', true, /\.vue$/), 'configComponent.image.Text')
addContext(require.context('./ISMComponents/device/', true, /\.vue$/), 'displayConfig.ToolBox.device.Container')
addContext(require.context('./ISMComponents/map/', true, /\.vue$/), 'displayConfig.ToolBox.Map')
addContext(require.context('./ISMComponents/historyCharts/', true, /\.vue$/), 'displayConfig.ToolBox.HistoryCharts')
addContext(require.context('./ISMComponents/Mes/standard/', true, /\.vue$/), 'displayConfig.ToolBox.MesStandard.title')

export const COMPONENT_REGISTRY = registry
export const TOOLBOX_GROUPS = toolboxGroups

export function getISMComponent(type) {
  return registry[type] || null
}
