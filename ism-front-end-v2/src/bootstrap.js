import {loadRoutes, loadGuards, setAppOptions} from '@/utils/routerUtil'
import {loadInterceptors} from '@/utils/request'
import guards from '@/router/guards'
import interceptors from '@/utils/axios-interceptors'

/**
 * 启动引导方法
 * 应用启动时需要执行的操作放在这里
 * @param router 应用的路由实例
 * @param store 应用的 vuex.store 实例
 * @param i18n 应用的 vue-i18n 实例
 * @param i18n 应用的 message 实例
 */
async function bootstrap({router, store, i18n, message}) {
  // 设置应用配置
  setAppOptions({router, store, i18n})
  // 加载 axios 拦截器
  loadInterceptors(interceptors, {router, store, i18n, message})
  // 先注册动态路由（从 localStorage Menu），避免深链刷新先命中 * → 404。
  // 首页配置接口可能很慢/超时，绝不能阻塞 loadRoutes。
  loadRoutes()
  // 加载路由守卫
  loadGuards(guards, {router, store, i18n, message})
  // 首页配置并行拉取；成功后仅补丁侧栏名称，失败不阻塞启动
  store.dispatch('setting/fetchSystemHomeDashboard').then(() => {
    const homeName = store.state.setting.systemHomeDashboard &&
      store.state.setting.systemHomeDashboard.dashboardName
    if (homeName) {
      store.commit('setting/patchScadaMonitorMenuName', homeName)
    }
  }).catch(e => {
    console.warn('[bootstrap] fetchSystemHomeDashboard failed, using defaults', e)
  })
}

export default bootstrap
