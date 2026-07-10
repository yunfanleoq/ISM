import Vue from 'vue'
import Vuex from 'vuex'
import storeModules from './modules'

Vue.use(Vuex)
const store = new Vuex.Store({ modules: storeModules })

export default store
