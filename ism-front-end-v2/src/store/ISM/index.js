import state from './state'
import * as getters from './getters'
import * as mutations from './mutations'
import * as actions from './actions'
const LockState = localStorage.getItem("LockState")
if(LockState&&LockState=='true')
{
  state.isLocked = true
}
else{
  state.isLocked = false
}
export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
}
