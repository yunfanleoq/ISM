/** 寻智平面页嵌入 AppRun：iframe URL 与无旁路树会话标记 */
import { XUNZHI_JF_ENV_PLANE_PAGE } from '@/xunzhiJf/resolveDisplayPage'

export const SESSION_ENV_PLANE_NO_RUN_TREE = 'xunzhiJf.envPlane.noRunTree'
export const SESSION_ENV_PLANE_PAGE_NAME = 'xunzhiJf.envPlane.pageName'

export function markEnvPlaneAppRunSession(pageName = XUNZHI_JF_ENV_PLANE_PAGE) {
  try {
    sessionStorage.setItem(SESSION_ENV_PLANE_NO_RUN_TREE, '1')
    sessionStorage.setItem(SESSION_ENV_PLANE_PAGE_NAME, pageName || XUNZHI_JF_ENV_PLANE_PAGE)
  } catch (e) { /* ignore */ }
}

export function clearEnvPlaneAppRunSession() {
  try {
    sessionStorage.removeItem(SESSION_ENV_PLANE_NO_RUN_TREE)
    sessionStorage.removeItem(SESSION_ENV_PLANE_PAGE_NAME)
  } catch (e) { /* ignore */ }
}

export function isEnvPlaneNoRunTreeSession() {
  try {
    return sessionStorage.getItem(SESSION_ENV_PLANE_NO_RUN_TREE) === '1'
  } catch (e) {
    return false
  }
}

/** 从 location.hash 解析 query（iframe 内 AppRun 与父页路由解耦时 route.query 可能为空） */
export function readHashRouteQuery() {
  const out = {}
  try {
    const hash = window.location.hash || ''
    const qIdx = hash.indexOf('?')
    if (qIdx < 0) return out
    const params = new URLSearchParams(hash.slice(qIdx + 1))
    params.forEach((value, key) => {
      out[key] = value
    })
  } catch (e) { /* ignore */ }
  return out
}

/** AppRun 须从 SPA 根加载，勿带 /xunzhi/jf/viz/plane 等嵌套 pathname */
export function buildAppRunIframeSrc(displayUid, query = {}) {
  const pp = (process.env.VUE_APP_PUBLIC_PATH || '/').replace(/\/?$/, '/')
  const root = `${window.location.origin}${pp}`
  const params = new URLSearchParams()
  Object.keys(query).forEach((key) => {
    const v = query[key]
    if (v !== undefined && v !== null && String(v) !== '') {
      params.set(key, String(v))
    }
  })
  const q = params.toString()
  return `${root}#/AppRun/${encodeURIComponent(displayUid)}${q ? `?${q}` : ''}`
}

export function buildEnvPlaneAppRunQuery(pageId) {
  const q = {
    pageName: XUNZHI_JF_ENV_PLANE_PAGE,
    noRunTree: '1',
  }
  if (pageId) q.pageId = pageId
  return q
}
