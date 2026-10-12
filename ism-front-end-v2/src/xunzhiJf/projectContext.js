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

export async function resolveXunzhiJfProjectUuid() {
  const fromEnv = process.env.VUE_APP_XUNZHI_JF_PROJECT_UUID
  if (fromEnv) return String(fromEnv).trim()

  try {
    const cached = sessionStorage.getItem(CACHE_UUID)
    if (cached) return cached
  } catch (e) { /* ignore */ }

  try {
    const res = await ProjectList()
    const list = (res && res.data && res.data.list) || []
    let hit = list.find(p => (p.name || p.Name) === XUNZHI_JF_PROJECT_NAME_HINT)
    if (!hit) {
      hit = list.find(
        p => (p.name || '').indexOf('xunzhi-jf') !== -1
          || (p.name || '').indexOf('寻智机房环境') !== -1
      )
    }
    const uuid = hit ? (hit.uuid || hit.UUID || '') : ''
    if (uuid) {
      sessionStorage.setItem(CACHE_UUID, uuid)
      sessionStorage.setItem(CACHE_NAME, hit.name || hit.Name || '')
    }
    return uuid
  } catch (e) {
    return ''
  }
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

/** GetSystemHomeDashboard 勿用航信等旧 scope 的 projectUuid 覆盖当前寻智 ProjectUuid 头 */
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

export function rememberProjectSelection(uuid, name) {
  if (!uuid) return
  try {
    sessionStorage.setItem(CACHE_UUID, uuid)
    if (name) sessionStorage.setItem(CACHE_NAME, name)
  } catch (e) { /* ignore */ }
}
