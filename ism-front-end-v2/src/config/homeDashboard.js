// 项目默认监控大屏：唯一来源是 GET /SystemHomeDashboard（按当前 ProjectUuid 隔离），禁止绑定某套种子数据 UUID。
import { AUTH_TYPE, setAuthorization } from '@/utils/request'

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

/** AppRun 是否显示左侧设备导航树（避免 AppRunShell 引用 xunanDashboardPages 引发循环依赖） */
export function shouldShowRunTreeNav(modelId, store) {
  if (!modelId) return false
  return isConfiguredHomeDashboard(modelId, store)
}

export function resolveHomeProjectUuid(dashboardUid, store) {
  if (dashboardUid === getHomeDashboardUuid(store)) {
    return getHomeProjectUuid(store)
  }
  return ''
}

export function createHomeDashboardRedirect() {
  return () => {
    const store = require('@/store').default
    // 未配置当前项目默认大屏时，回到项目管理页，勿跳到别的项目大屏
    if (!getHomeDashboardUuid(store)) {
      return '/dashboard'
    }
    return getHomeDashboardPath(store)
  }
}
