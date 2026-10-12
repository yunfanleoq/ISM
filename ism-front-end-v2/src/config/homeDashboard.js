// 项目默认监控大屏：唯一来源是 GET /SystemHomeDashboard（按当前 ProjectUuid 隔离），禁止绑定某套种子数据 UUID。
import { AUTH_TYPE, setAuthorization } from '@/utils/request'
import {
  buildEnvPlaneAppRunQuery,
  isNoRunTreeMode,
} from '@/xunzhiJf/appRunEmbed'
import { isXunzhiJfActiveProject } from '@/xunzhiJf/projectContext'

export const HOME_DASHBOARD_UUID = ''
export const HOME_PROJECT_UUID = ''
export const HOME_DASHBOARD_PATH = `/AppRun/${HOME_DASHBOARD_UUID}`
export const LEGACY_NCC_DASHBOARD_UUID = ''
export const LEGACY_NCC_PROJECT_UUID = ''
export const XUNAN_DASHBOARD_UUID_ALT = ''
export const RUN_TREE_DASHBOARD_UUIDS = []

function readStoreConfig(store) {
  return store && store.state && store.state.setting
    ? store.state.setting.systemHomeDashboard
    : null
}

export function getHomeDashboardConfig(store) {
  const cfg = readStoreConfig(store)
  return {
    dashboardUuid: (cfg && cfg.dashboardUuid) || HOME_DASHBOARD_UUID,
    projectUuid: (cfg && cfg.projectUuid) || HOME_PROJECT_UUID,
    dashboardName: (cfg && cfg.dashboardName) || '',
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

function shouldUseXunzhiJfVizEmbed(store) {
  if (isXunzhiJfActiveProject()) return true
  const name = getHomeDashboardConfig(store).dashboardName || ''
  return /寻智|xunzhi-jf|平面监视|env-plane/i.test(name)
}

/** 侧栏「监控大屏」/ SCADAMonitor：寻智平面走 embed；否则 AppRun 带 noRunTree */
export function resolveScadaMonitorRoute() {
  const store = require('@/store').default
  const uid = getHomeDashboardUuid(store)
  if (!uid) {
    return '/dashboard'
  }
  if (shouldUseXunzhiJfVizEmbed(store)) {
    return '/xunzhi/jf/viz/plane'
  }
  const name = getHomeDashboardConfig(store).dashboardName || ''
  if (/平面|env-plane|寻智|xunzhi/i.test(name)) {
    return {
      path: `/AppRun/${uid}`,
      query: buildEnvPlaneAppRunQuery(),
    }
  }
  return `/AppRun/${uid}`
}

export function createHomeDashboardRedirect() {
  return () => resolveScadaMonitorRoute()
}
