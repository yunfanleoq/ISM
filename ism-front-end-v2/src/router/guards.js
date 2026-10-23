import {hasAuthority} from '@/utils/authority-utils'
import {loginIgnore} from '@/router/index'
import {checkAuthorization} from '@/utils/request'
import { HOME_DASHBOARD_PATH, isRetiredW1Target } from '@/config/homeDashboard'
import NProgress from 'nprogress'

NProgress.configure({ showSpinner: false })

/**
 * 进度条开始
 * @param to
 * @param form
 * @param next
 */
const progressStart = (to, from, next) => {
  // start progress bar
  if (!NProgress.isStarted()) {
    NProgress.start()
  }
  next()
}

/**
 * 登录守卫
 * @param to
 * @param form
 * @param next
 * @param options
 */
const loginGuard = (to, from, next, options) => {
  const {message} = options
  if ((to.path.indexOf('/AppRun')>=0) && !checkAuthorization())
  {

    message.warning('登录已失效，请重新登录')
    next({path: '/AppLogin/'+to.params.uid})
  }
  else
  {
    if (!loginIgnore.includes(to) && !checkAuthorization()) {
      message.warning('登录已失效，请重新登录')
      next({path: '/login'})
    } else {
      next()
    }
  }

}

/**
 * 权限守卫
 * @param to
 * @param form
 * @param next
 * @param options
 */
const authorityGuard = (to, from, next, options) => {
  const {store, message} = options
  const permissions = store.getters['account/permissions']
  const roles = store.getters['account/roles']
  if (!hasAuthority(to, permissions, roles)) {
    message.warning(`对不起，您无权访问页面: ${to.fullPath}，请联系管理员`)
    next({path: '/403'})
    // NProgress.done()
  } else {
    next()
  }
}

/** 寻智 W1 平面已停用：直接打开旧地址或环境平面菜单时改去公安局动环 */
const retiredW1Guard = (to, from, next) => {
  if (!isRetiredW1Target(to) || to.path === HOME_DASHBOARD_PATH) {
    next()
    return
  }
  next({ path: HOME_DASHBOARD_PATH, replace: true })
}

/**
 * 混合导航模式下一级菜单跳转重定向
 * @param to
 * @param from
 * @param next
 * @param options
 * @returns {*}
 */
const redirectGuard = (to, from, next, options) => {
  const {store} = options
  const getFirstChild = (routes) => {
    const route = routes[0]
    if (!route.children || route.children.length === 0) {
      return route
    }
    return getFirstChild(route.children)
  }
  if (store.state.setting.layout === 'mix') {
    const firstMenu = store.getters['setting/firstMenu']
    if (firstMenu.find(item => item.fullPath === to.fullPath)) {
      store.commit('setting/setActivatedFirst', to.fullPath)
      const subMenu = store.getters['setting/subMenu']
      if (subMenu.length > 0) {
        const redirect = getFirstChild(subMenu)
        return next({path: redirect.fullPath})
      }
    }
  }
  next()
}

/**
 * 进度条结束
 * @param to
 * @param form
 * @param options
 */
const progressDone = () => {
  // finish progress bar
  NProgress.done()
}

export default {
  beforeEach: [progressStart, retiredW1Guard, loginGuard, authorityGuard, redirectGuard],
  afterEach: [progressDone]
}
