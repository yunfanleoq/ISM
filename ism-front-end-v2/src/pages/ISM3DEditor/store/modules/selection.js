/**
 * Selection Module - 选择状态管理
 */
const state = {
  objectId: null,
  multiple: false,
  selectedIds: []
}

const mutations = {
  SELECT_OBJECT(state, objectId) {
    state.objectId = objectId
  },

  ADD_TO_SELECTION(state, objectId) {
    if (!state.selectedIds.includes(objectId)) {
      state.selectedIds.push(objectId)
    }
  },

  REMOVE_FROM_SELECTION(state, objectId) {
    const index = state.selectedIds.indexOf(objectId)
    if (index !== -1) {
      state.selectedIds.splice(index, 1)
    }
  },

  CLEAR_SELECTION(state) {
    state.objectId = null
    state.selectedIds = []
  },

  SET_MULTIPLE(state, value) {
    state.multiple = value
  }
}

const actions = {
  selectObject({ commit }, objectId) {
    commit('SELECT_OBJECT', objectId)
  },

  toggleSelection({ commit, state }, objectId) {
    if (state.selectedIds.includes(objectId)) {
      commit('REMOVE_FROM_SELECTION', objectId)
    } else {
      commit('ADD_TO_SELECTION', objectId)
    }
  },

  clearSelection({ commit }) {
    commit('CLEAR_SELECTION')
  },

  reset({ commit }) {
    commit('CLEAR_SELECTION')
    commit('SET_MULTIPLE', false)
  }
}

const getters = {
  hasSelection: (state) => !!state.objectId,
  isSelected: (state) => (id) => state.objectId === id || state.selectedIds.includes(id)
}

export default {
  namespaced: true,

  state,
  mutations,
  actions,
  getters
}
