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

export function isNoRunTreeFlag(value) {
  if (value === undefined || value === null || value === '') return false
  if (value === true || value === 1) return true
  const s = String(value).toLowerCase()
  return s === '1' || s === 'true' || s === 'yes'
}

/** 从 location.hash / search 解析 query（iframe AppRun 与父页路由解耦时 route.query 可能为空） */
export function readHashRouteQuery() {
  const out = {}
  try {
    const search = window.location.search || ''
    if (search.startsWith('?')) {
      new URLSearchParams(search.slice(1)).forEach((value, key) => {
        out[key] = value
      })
    }
    const hash = window.location.hash || ''
    const qIdx = hash.indexOf('?')
    if (qIdx >= 0) {
      const params = new URLSearchParams(hash.slice(qIdx + 1))
      params.forEach((value, key) => {
        out[key] = value
      })
    }
  } catch (e) { /* ignore */ }
  return out
}

export function mergedAppRunQuery(route) {
  const fromRoute = route && route.query ? route.query : {}
  const fromHash = readHashRouteQuery()
  return { ...fromHash, ...fromRoute }
}

/** AppRun / embed：须隐藏旁路 RunTree 与组态内嵌「机房监视树」 */
export function isNoRunTreeMode(route, store) {
  const q = mergedAppRunQuery(route)
  if (isNoRunTreeFlag(q.noRunTree)) return true
  if (q.pageName === XUNZHI_JF_ENV_PLANE_PAGE) return true
  if (
    isEnvPlaneNoRunTreeSession()
    && (isNoRunTreeFlag(q.noRunTree) || q.pageName === XUNZHI_JF_ENV_PLANE_PAGE)
  ) {
    return true
  }
  const pageId = q.pageId || ''
  if (pageId && store && store.state && store.state.ISM) {
    const lists = [
      ...(store.state.ISM.PCPageList || []),
      ...(store.state.ISM.PhonePageList || []),
    ]
    const hit = lists.find(p => (p.pageUuid || p.PageId || '') === pageId)
    const name = hit && (hit.name || hit.PageName || hit.title || '')
    if (name === XUNZHI_JF_ENV_PLANE_PAGE) return true
  }
  const ism = store && store.state && store.state.ISM
  if (
    ism
    && ism.LayerData
    && ism.LayerData.name === XUNZHI_JF_ENV_PLANE_PAGE
    && (isNoRunTreeFlag(q.noRunTree) || q.pageName === XUNZHI_JF_ENV_PLANE_PAGE || isEnvPlaneNoRunTreeSession())
  ) {
    return true
  }
  return false
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
