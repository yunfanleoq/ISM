/**
 * ISM3DEditor Vuex Store
 * 状态管理入口
 */
import scene from './modules/scene'
import selection from './modules/selection'
import history from './modules/history'
import dataBinding from './modules/dataBinding'

export default {
  namespaced: true,

  modules: {
    scene,
    selection,
    history,
    dataBinding
  },

  state: {
    initialized: false,
    sceneSettings: {
      background: '#1a1a2e',
      gridVisible: true,
      shadowEnabled: true,
      ambientIntensity: 0.5
    }
  },

  mutations: {
    SET_INITIALIZED(state, value) {
      state.initialized = value
    },

    UPDATE_SCENE_SETTINGS(state, settings) {
      Object.assign(state.sceneSettings, settings)
    }
  },

  actions: {
    initScene({ commit, dispatch }) {
      dispatch('scene/init', null, { root: true })
      dispatch('selection/reset', null, { root: true })
      dispatch('history/clear', null, { root: true })
      commit('SET_INITIALIZED', true)
    },

    cleanupScene({ commit, state }) {
      if (!state.initialized) return

      commit('SET_INITIALIZED', false)
    }
  },

  getters: {
    selectedObject: (state) => {
      if (!state.selection.objectId) return null
      return state.scene.objects.find(obj => obj.id === state.selection.objectId)
    },

    canUndo: (state) => state.history.past.length > 0,
    canRedo: (state) => state.history.future.length > 0
  }
}
