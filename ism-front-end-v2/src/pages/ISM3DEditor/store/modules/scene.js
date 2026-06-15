/**
 * Scene Module - 场景对象状态管理
 */
import * as THREE from 'three'

const generateId = () => {
  return 'obj_' + Math.random().toString(36).substr(2, 9)
}

const state = {
  objects: [],
  groups: [],
  objectCounter: 0
}

const mutations = {
  ADD_OBJECT(state, object) {
    state.objects.push(object)
  },

  REMOVE_OBJECT(state, objectId) {
    const index = state.objects.findIndex(obj => obj.id === objectId)
    if (index !== -1) {
      state.objects.splice(index, 1)
    }
  },

  UPDATE_OBJECT(state, { id, changes }) {
    const object = state.objects.find(obj => obj.id === id)
    if (object) {
      Object.assign(object, changes)
    }
  },

  SET_OBJECTS(state, objects) {
    state.objects = objects
  },

  CLEAR_OBJECTS(state) {
    state.objects = []
    state.groups = []
  }
}

const actions = {
  init({ commit }) {
    commit('CLEAR_OBJECTS')
  },

  addObject({ commit, state }, { type, options }) {
    const id = generateId()
    const object = {
      id,
      type,
      name: `${type}_${++state.objectCounter}`,
      transform: {
        position: { x: 0, y: 0, z: 0 },
        rotation: { x: 0, y: 0, z: 0 },
        scale: { x: 1, y: 1, z: 1 }
      },
      appearance: {
        color: '#4a90d9',
        opacity: 1,
        metalness: 0.3,
        roughness: 0.7,
        wireframe: false
      },
      material: options?.material || 'standard',
      dataBindings: [],
      children: [],
      parentId: null,
      locked: false,
      visible: true,
      ...options
    }

    commit('ADD_OBJECT', object)
    return object
  },

  removeObject({ commit }, objectId) {
    commit('REMOVE_OBJECT', objectId)
  },

  updateObject({ commit }, { id, changes }) {
    commit('UPDATE_OBJECT', { id, changes })
  },

  setObjects({ commit }, objects) {
    commit('SET_OBJECTS', objects)
  },

  cloneObject({ state, dispatch }, objectId) {
    const source = state.objects.find(obj => obj.id === objectId)
    if (!source) return null

    return dispatch('addObject', {
      type: source.type,
      options: {
        ...source,
        name: `${source.name}_copy`,
        transform: {
          position: {
            x: source.transform.position.x + 2,
            y: source.transform.position.y,
            z: source.transform.position.z
          },
          rotation: { ...source.transform.rotation },
          scale: { ...source.transform.scale }
        },
        appearance: { ...source.appearance },
        dataBindings: []
      }
    })
  }
}

const getters = {
  objectById: (state) => (id) => {
    return state.objects.find(obj => obj.id === id)
  },

  objectsByType: (state) => (type) => {
    return state.objects.filter(obj => obj.type === type)
  },

  visibleObjects: (state) => {
    return state.objects.filter(obj => obj.visible)
  }
}

export default {
  namespaced: true,

  state,
  mutations,
  actions,
  getters
}
