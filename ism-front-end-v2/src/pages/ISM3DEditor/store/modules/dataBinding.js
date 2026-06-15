/**
 * DataBinding Module - 数据绑定状态管理
 */
const generateBindingId = () => {
  return 'binding_' + Math.random().toString(36).substr(2, 9)
}

const state = {
  bindings: [],
  dataSources: [],
  connectionStatus: 'disconnected'
}

const mutations = {
  ADD_BINDING(state, binding) {
    state.bindings.push(binding)
  },

  REMOVE_BINDING(state, bindingId) {
    const index = state.bindings.findIndex(b => b.id === bindingId)
    if (index !== -1) {
      state.bindings.splice(index, 1)
    }
  },

  UPDATE_BINDING(state, { id, changes }) {
    const binding = state.bindings.find(b => b.id === id)
    if (binding) {
      Object.assign(binding, changes)
    }
  },

  CLEAR_OBJECT_BINDINGS(state, objectId) {
    state.bindings = state.bindings.filter(b => b.objectId !== objectId)
  },

  ADD_DATA_SOURCE(state, source) {
    state.dataSources.push(source)
  },

  REMOVE_DATA_SOURCE(state, sourceId) {
    const index = state.dataSources.findIndex(s => s.id === sourceId)
    if (index !== -1) {
      state.dataSources.splice(index, 1)
    }
  },

  SET_CONNECTION_STATUS(state, status) {
    state.connectionStatus = status
  }
}

const actions = {
  addBinding({ commit }, { objectId, binding }) {
    const newBinding = {
      id: generateBindingId(),
      objectId,
      property: binding.property,
      dataSource: binding.dataSource,
      dataPath: binding.dataPath,
      transform: binding.transform || {
        type: 'direct',
        scale: 1,
        offset: 0
      },
      enabled: true,
      ...binding
    }
    commit('ADD_BINDING', newBinding)
    return newBinding
  },

  removeBinding({ commit }, bindingId) {
    commit('REMOVE_BINDING', bindingId)
  },

  updateBinding({ commit }, { id, changes }) {
    commit('UPDATE_BINDING', { id, changes })
  },

  clearBindings({ commit }, objectId) {
    commit('CLEAR_OBJECT_BINDINGS', objectId)
  },

  addDataSource({ commit }, source) {
    const newSource = {
      id: generateBindingId(),
      name: source.name,
      type: source.type,
      config: source.config || {},
      ...source
    }
    commit('ADD_DATA_SOURCE', newSource)
    return newSource
  },

  removeDataSource({ commit }, sourceId) {
    commit('REMOVE_DATA_SOURCE', sourceId)
  }
}

const getters = {
  bindingsByObject: (state) => (objectId) => {
    return state.bindings.filter(b => b.objectId === objectId)
  },

  bindingByProperty: (state) => (objectId, property) => {
    return state.bindings.find(b => b.objectId === objectId && b.property === property)
  },

  enabledBindings: (state) => {
    return state.bindings.filter(b => b.enabled)
  }
}

export default {
  namespaced: true,

  state,
  mutations,
  actions,
  getters
}
