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

export function isOnXunzhiJfRoute() {
  if (typeof window === 'undefined') return false
  const hash = window.location.hash || ''
  const path = window.location.pathname || ''
  return hash.indexOf('/xunzhi/jf') >= 0 || path.indexOf('/xunzhi/jf') >= 0
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
 * 以 project_lists / ProjectList API 为准（禁止钉死 uuid；重导后 uuid 会变）。
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
        const active = getActiveProjectUuid()
        if (!active || active !== uuid || isOnXunzhiJfRoute()) {
          syncXunzhiJfProjectToClient(uuid, hit.name || hit.Name || '')
        }
        return uuid
      }
    }
    const cached = sessionStorage.getItem(CACHE_UUID)
    if (cached && list.some(p => (p.uuid || p.UUID) === cached)) {
      syncXunzhiJfProjectToClient(cached, sessionStorage.getItem(CACHE_NAME) || '')
      return cached
    }
    invalidateXunzhiJfProjectCache()
  } catch (e) {
    try {
      const cached = sessionStorage.getItem(CACHE_UUID)
      if (cached) {
        syncXunzhiJfProjectToClient(cached, sessionStorage.getItem(CACHE_NAME) || '')
        return cached
      }
    } catch (err) { /* ignore */ }
  }

  const fallback = getActiveProjectUuid()
  if (fallback) {
    syncXunzhiJfProjectToClient(fallback, sessionStorage.getItem(CACHE_NAME) || '')
  }
  return fallback
}

export async function ensureXunzhiJfProjectAuth() {
  const uuid = await resolveXunzhiJfProjectUuid()
  if (uuid) {
    setAuthorization({ token: uuid }, AUTH_TYPE.AUTH1)
  }
  return uuid
}

/**
 * 直链 /xunzhi/jf/viz/plane 等：先对齐 ProjectUuid + 拉 HomeDashboard，再解析组态页。
 */
export async function primeXunzhiJfRouteContext(store) {
  const uuid = await resolveXunzhiJfProjectUuid()
  await ensureXunzhiJfProjectAuth()
  const vuex = store || (typeof require !== 'undefined' ? require('@/store').default : null)
  if (vuex && vuex.dispatch) {
    try {
      await vuex.dispatch('setting/fetchSystemHomeDashboard')
    } catch (e) { /* ignore */ }
  }
  return uuid
}

export function xunzhiJfVizRouteBeforeEnter(to, from, next) {
  primeXunzhiJfRouteContext()
    .then(() => next())
    .catch(() => next())
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

export function alignHomeDashboardProjectUuid(apiProjectUuid) {
  const active = getActiveProjectUuid()
  if (!active) return apiProjectUuid || ''
  if (!apiProjectUuid || apiProjectUuid === active) {
    return apiProjectUuid || active
  }
  if (isXunzhiJfActiveProject()) return active
  if (isOnXunzhiJfRoute()) {
    return active
  }
  return apiProjectUuid
}

export async function alignHomeDashboardProjectUuidAsync(apiProjectUuid, dashboardName) {
  const xunzhiUuid = await resolveXunzhiJfProjectUuid()
  if (xunzhiUuid && (isXunzhiJfActiveProject() || isXunzhiJfDashboardName(dashboardName) || isOnXunzhiJfRoute())) {
    return xunzhiUuid
  }
  return alignHomeDashboardProjectUuid(apiProjectUuid)
}

export function rememberProjectSelection(uuid, name) {
  if (!uuid) return
  syncXunzhiJfProjectToClient(uuid, name)
}
