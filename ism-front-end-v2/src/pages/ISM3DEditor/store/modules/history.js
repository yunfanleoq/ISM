/**
 * History Module - 撤销/重做管理
 */
const MAX_HISTORY = 50

const state = {
  past: [],
  future: []
}

const mutations = {
  PUSH_HISTORY(state, snapshot) {
    state.past.push(snapshot)
    state.future = []

    if (state.past.length > MAX_HISTORY) {
      state.past.shift()
    }
  },

  UNDO(state) {
    if (state.past.length === 0) return null
    const snapshot = state.past.pop()
    return snapshot
  },

  REDO(state) {
    if (state.future.length === 0) return null
    const snapshot = state.future.pop()
    return snapshot
  },

  SAVE_TO_FUTURE(state, snapshot) {
    state.future.push(snapshot)
  },

  CLEAR(state) {
    state.past = []
    state.future = []
  }
}

const actions = {
  record({ commit }, snapshot) {
    commit('PUSH_HISTORY', JSON.parse(JSON.stringify(snapshot)))
  },

  undo({ commit, state }) {
    if (state.past.length === 0) return null
    return commit('UNDO')
  },

  redo({ commit, state }) {
    if (state.future.length === 0) return null
    return commit('REDO')
  },

  saveCurrent({ commit }, currentState) {
    commit('SAVE_TO_FUTURE', JSON.parse(JSON.stringify(currentState)))
  },

  clear({ commit }) {
    commit('CLEAR')
  }
}

const getters = {
  canUndo: (state) => state.past.length > 0,
  canRedo: (state) => state.future.length > 0
}

export default {
  namespaced: true,

  state,
  mutations,
  actions,
  getters
}
