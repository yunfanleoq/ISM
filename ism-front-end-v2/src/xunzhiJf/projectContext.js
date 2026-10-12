import Cookie from 'js-cookie'
import { ProjectList } from '@/services/project'
import { AUTH_TYPE, setAuthorization } from '@/utils/request'
import { XUNZHI_JF_PROJECT_NAME_HINT } from '@/xunzhiJf/resolveDisplayPage'

const CACHE_UUID = 'xunzhiJf.projectUuid'
const CACHE_NAME = 'xunzhiJf.projectName'

export function getActiveProjectUuid() {
  return (
    Cookie.get('ProjectUuid')
    || (typeof sessionStorage !== 'undefined' && sessionStorage.getItem('ProjectUuid'))
    || ''
  )
}

export function isXunzhiJfDashboardName(name) {
  const n = String(name || '')
  return /寻智|xunzhi-jf|平面监视|env-plane|机房环境/i.test(n)
}

export function syncXunzhiJfProjectToClient(uuid, name) {
  if (!uuid) return
  try {
    Cookie.set('ProjectUuid', uuid)
    sessionStorage.setItem('ProjectUuid', uuid)
    sessionStorage.setItem(CACHE_UUID, uuid)
    if (name) sessionStorage.setItem(CACHE_NAME, name)
  } catch (e) { /* ignore */ }
  setAuthorization({ token: uuid }, AUTH_TYPE.AUTH1)
}

export function invalidateXunzhiJfProjectCache() {
  try {
    sessionStorage.removeItem(CACHE_UUID)
    sessionStorage.removeItem(CACHE_NAME)
  } catch (e) { /* ignore */ }
}

function pickXunzhiJfFromList(list) {
  if (!Array.isArray(list) || !list.length) return null
  let hit = list.find(p => (p.name || p.Name) === XUNZHI_JF_PROJECT_NAME_HINT)
  if (!hit) {
    hit = list.find(
      p => (p.name || '').indexOf('xunzhi-jf') !== -1
        || (p.name || '').indexOf('寻智机房环境') !== -1
    )
  }
  return hit || null
}

/**
 * 以 project_lists / ProjectList API 为准解析寻智项目 uuid（禁止信任过期 session 缓存）。
 */
export async function resolveXunzhiJfProjectUuid() {
  const fromEnv = process.env.VUE_APP_XUNZHI_JF_PROJECT_UUID
  if (fromEnv) {
    const uuid = String(fromEnv).trim()
    syncXunzhiJfProjectToClient(uuid, XUNZHI_JF_PROJECT_NAME_HINT)
    return uuid
  }

  try {
    const res = await ProjectList()
    const list = (res && res.data && res.data.list) || []
    const hit = pickXunzhiJfFromList(list)
    if (hit) {
      const uuid = hit.uuid || hit.UUID || ''
      if (uuid) {
        syncXunzhiJfProjectToClient(uuid, hit.name || hit.Name || '')
        return uuid
      }
    }
    const cached = sessionStorage.getItem(CACHE_UUID)
    if (cached && list.some(p => (p.uuid || p.UUID) === cached)) {
      return cached
    }
    invalidateXunzhiJfProjectCache()
  } catch (e) {
    try {
      const cached = sessionStorage.getItem(CACHE_UUID)
      if (cached) return cached
    } catch (err) { /* ignore */ }
  }

  return getActiveProjectUuid()
}

export async function ensureXunzhiJfProjectAuth() {
  const uuid = await resolveXunzhiJfProjectUuid()
  if (uuid) {
    setAuthorization({ token: uuid }, AUTH_TYPE.AUTH1)
  }
  return uuid
}

export function isXunzhiJfActiveProject() {
  try {
    const name = sessionStorage.getItem(CACHE_NAME) || ''
    if (name === XUNZHI_JF_PROJECT_NAME_HINT || name.indexOf('xunzhi-jf') !== -1) {
      return true
    }
    const active = getActiveProjectUuid()
    const cached = sessionStorage.getItem(CACHE_UUID) || ''
    return !!(active && cached && active === cached)
  } catch (e) {
    return false
  }
}

/** GetSystemHomeDashboard 返回的 projectUuid 须与当前库寻智项目一致（重导后 uuid 会变） */
export function alignHomeDashboardProjectUuid(apiProjectUuid) {
  const active = getActiveProjectUuid()
  if (!active) return apiProjectUuid || ''
  if (!apiProjectUuid || apiProjectUuid === active) {
    return apiProjectUuid || active
  }
  if (isXunzhiJfActiveProject()) return active
  if (typeof window !== 'undefined' && window.location.hash.indexOf('/xunzhi/jf') >= 0) {
    return active
  }
  return apiProjectUuid
}

export async function alignHomeDashboardProjectUuidAsync(apiProjectUuid, dashboardName) {
  const xunzhiUuid = await resolveXunzhiJfProjectUuid()
  if (xunzhiUuid && (isXunzhiJfActiveProject() || isXunzhiJfDashboardName(dashboardName))) {
    return xunzhiUuid
  }
  return alignHomeDashboardProjectUuid(apiProjectUuid)
}

export function rememberProjectSelection(uuid, name) {
  if (!uuid) return
  syncXunzhiJfProjectToClient(uuid, name)
}
