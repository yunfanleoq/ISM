// 登录后默认大屏固定为公安局动环。接口里若仍指向寻智机房环境平面图 W1，不能再当默认。
import { AUTH_TYPE, setAuthorization } from '@/utils/request'
import { isNoRunTreeMode } from '@/xunzhiJf/appRunEmbed'

export const HOME_DASHBOARD_UUID = 'a7b91292-dc10-4f00-8000-0000000000b1'
export const HOME_PROJECT_UUID = 'd96005d5-dd28-3ad2-bc4d-983316c3f218'
export const HOME_DASHBOARD_PATH = `/AppRun/${HOME_DASHBOARD_UUID}`
export const LEGACY_NCC_DASHBOARD_UUID = ''
export const LEGACY_NCC_PROJECT_UUID = ''
export const XUNAN_DASHBOARD_UUID_ALT = ''
export const RUN_TREE_DASHBOARD_UUIDS = []
/** 寻智机房环境平面图 W1，已停用 */
export const RETIRED_W1_ENV_PLANE_UUID = 'c7a8b8b6-c5c2-3d05-2971-281d65c442ca'

export function getHomeDashboardConfig() {
  return {
    dashboardUuid: HOME_DASHBOARD_UUID,
    projectUuid: HOME_PROJECT_UUID,
    dashboardName: '公安局动环',
  }
}

export function getHomeDashboardUuid(store) {
  return getHomeDashboardConfig(store).dashboardUuid
}

export function getHomeProjectUuid(store) {
  return getHomeDashboardConfig(store).projectUuid
}

export function getHomeDashboardPath(store) {
  return `/AppRun/${getHomeDashboardUuid(store)}`
}

/** 进入项目默认大屏前，将 ProjectUuid 请求头同步为该大屏所属项目（避免组态 API 因项目不一致一直 loading） */
export function applyHomeProjectAuth(store) {
  const projectUuid = getHomeProjectUuid(store)
  if (projectUuid) {
    setAuthorization({ token: projectUuid }, AUTH_TYPE.AUTH1)
  }
}

export function isConfiguredHomeDashboard(uid, store) {
  const configured = getHomeDashboardUuid(store)
  return !!configured && uid === configured
}

/** 寻智动环平面页；noRunTree=1 时隐藏旁路 RunTree 与组态内嵌机房监视树 */
export const XUNZHI_JF_ENV_PLANE_PAGE = 'xunzhi-jf-env-plane'

/** 当前 AppRun 是否为寻智平面 noRunTree 模式（旁路树 + 内嵌监视树均隐藏） */
export function isXunzhiJfEnvPlaneAppRun(route, store) {
  return isNoRunTreeMode(route, store)
}

/** AppRun 是否显示左侧设备导航树（避免 AppRunShell 引用 xunanDashboardPages 引发循环依赖） */
export function shouldShowRunTreeNav(modelId, store, route) {
  if (!modelId) return false
  if (route && isXunzhiJfEnvPlaneAppRun(route, store)) return false
  return isConfiguredHomeDashboard(modelId, store)
}

export function resolveHomeProjectUuid(dashboardUid, store) {
  if (dashboardUid === getHomeDashboardUuid(store)) {
    return getHomeProjectUuid(store)
  }
  return ''
}

function routePath(to) {
  const raw = (to && (to.path || to.fullPath)) || ''
  return String(raw).split('?')[0].replace(/\/+$/, '') || '/'
}

/** 寻智 W1 平面：AppRun 旧 uuid，或环境平面路由 / pageName */
export function isRetiredW1Target(to) {
  if (!to) return false
  const path = routePath(to)
  const uid = (to.params && to.params.uid) || ''
  if (uid === RETIRED_W1_ENV_PLANE_UUID) return true
  if (path.indexOf('/AppRun/' + RETIRED_W1_ENV_PLANE_UUID) === 0) return true
  if (path === '/xunzhi/jf/viz/plane') return true
  const pageName = (to.query && (to.query.pageName || to.query.page)) || ''
  if (pageName === XUNZHI_JF_ENV_PLANE_PAGE && path.indexOf(RETIRED_W1_ENV_PLANE_UUID) >= 0) {
    return true
  }
  return false
}

/** 菜单里指向 W1 的项隐藏，避免还能点进去 */
export function hideRetiredHomeMenu(routes) {
  if (!Array.isArray(routes)) return routes
  return routes.map(route => {
    if (!route) return route
    const children = route.children ? hideRetiredHomeMenu(route.children) : route.children
    const path = String(route.path || '')
    const fullPath = String(route.fullPath || '')
    const name = String(route.name || '')
    const menuCode = route.meta && route.meta.menuCode
    const retired = path.indexOf(RETIRED_W1_ENV_PLANE_UUID) >= 0
      || fullPath.indexOf(RETIRED_W1_ENV_PLANE_UUID) >= 0
      || path === '/xunzhi/jf/viz/plane'
      || fullPath === '/xunzhi/jf/viz/plane'
      || menuCode === 'jf.viz.plane'
      || name === 'XunzhiJf_jf_viz_plane'
      || name === 'XunzhiJfVizPlaneAlias'
    if (!retired && children === route.children) return route
    return {
      ...route,
      children,
      meta: retired ? { ...(route.meta || {}), invisible: true } : route.meta,
    }
  })
}

/** 侧栏「监控大屏」/ 登录后落地：固定公安局动环 */
export function resolveScadaMonitorRoute() {
  return HOME_DASHBOARD_PATH
}

export function createHomeDashboardRedirect() {
  return () => resolveScadaMonitorRoute()
}
