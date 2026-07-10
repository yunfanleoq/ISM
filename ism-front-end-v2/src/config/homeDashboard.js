// 系统首页大屏配置（默认值 + 动态读取 store）
// 运行时优先使用 setting.systemHomeDashboard（来自 GET /SystemHomeDashboard）
import { AUTH_TYPE, setAuthorization } from '@/utils/request'

export const HOME_DASHBOARD_UUID = 'b8b4c094-faa9-a22a-1d0d-037539b27a6c'
export const HOME_PROJECT_UUID = '3ec5821f-b512-2adb-3e1c-473720d0a93e'
export const HOME_DASHBOARD_PATH = `/AppRun/${HOME_DASHBOARD_UUID}`
export const LEGACY_NCC_DASHBOARD_UUID = '043135ad-44be-e5d8-89be-3e54883c23a8'
export const LEGACY_NCC_PROJECT_UUID = '31bc90be-ebc4-dd61-ba9d-ce6e075e40e2'
/** 应用管理里的「电力监控大屏」副本（同项目同架构） */
export const XUNAN_DASHBOARD_UUID_ALT = '8278e8be-eb42-a30e-231b-1850fde94ca5'

/** 需挂载左侧 ISMRunTreeNav 的 display_model_uid */
export const RUN_TREE_DASHBOARD_UUIDS = [
  HOME_DASHBOARD_UUID,
  XUNAN_DASHBOARD_UUID_ALT,
  LEGACY_NCC_DASHBOARD_UUID,
]

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

/** 进入系统首页大屏前，将 ProjectUuid 请求头同步为首页所属项目（避免组态 API 因项目不一致一直 loading） */
export function applyHomeProjectAuth(store) {
  const projectUuid = getHomeProjectUuid(store)
  if (projectUuid) {
    setAuthorization({ token: projectUuid }, AUTH_TYPE.AUTH1)
  }
}

export function isConfiguredHomeDashboard(uid, store) {
  return uid === getHomeDashboardUuid(store)
}

/** AppRun 是否显示左侧设备导航树（避免 AppRunShell 引用 xunanDashboardPages 引发循环依赖） */
export function shouldShowRunTreeNav(modelId, store) {
  if (!modelId) return false
  if (RUN_TREE_DASHBOARD_UUIDS.includes(modelId)) return true
  return isConfiguredHomeDashboard(modelId, store)
}

export function resolveHomeProjectUuid(dashboardUid, store) {
  if (dashboardUid === getHomeDashboardUuid(store)) {
    return getHomeProjectUuid(store)
  }
  if (dashboardUid === XUNAN_DASHBOARD_UUID_ALT) {
    return HOME_PROJECT_UUID
  }
  if (dashboardUid === LEGACY_NCC_DASHBOARD_UUID) {
    return LEGACY_NCC_PROJECT_UUID
  }
  return ''
}

export function createHomeDashboardRedirect() {
  return () => {
    const store = require('@/store').default
    return getHomeDashboardPath(store)
  }
}
