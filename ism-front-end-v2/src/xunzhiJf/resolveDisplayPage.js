import { displayModelList } from '@/services/displayModel'
import { getDisplayModelLayerData } from '@/services/displayModel'
import { GetSystemHomeDashboard } from '@/services/system'
import {
  ensureXunzhiJfProjectAuth,
  invalidateXunzhiJfProjectCache,
  primeXunzhiJfRouteContext,
  resolveXunzhiJfProjectUuid,
} from '@/xunzhiJf/projectContext'
import { getHomeDashboardConfig } from '@/config/homeDashboard'

export const XUNZHI_JF_PROJECT_NAME_HINT = '循智元晟'
export const XUNZHI_JF_PROJECT_NAME_LEGACY = '寻智机房环境-xunzhi-jf'
export const XUNZHI_JF_ENV_PLANE_PAGE = 'xunzhi-jf-env-plane'
export const XUNZHI_JF_NAV_25D_PAGE = 'xunzhi-jf-25d-shell'

function layerPageName(layer) {
  if (!layer) return ''
  return String(layer.PageName || layer.pageName || layer.Name || layer.name || '').trim()
}

function isLoginLayer(layer) {
  const v = layer.IsLogin ?? layer.isLogin
  return v === 1 || v === '1' || v === true
}

async function findPageInDisplay(displayUid, pageName, displayName) {
  if (!displayUid) return null
  const layerRes = await getDisplayModelLayerData({ muid: displayUid, metaOnly: true })
  if (!layerRes || !layerRes.data || layerRes.data.code !== 0) return null
  const layers = layerRes.data.layer || []
  const page = layers.find(l => layerPageName(l) === pageName && !isLoginLayer(l))
  if (!page) return null
  return {
    displayUid,
    pageId: page.PageId || page.pageId || page.PageUUID || page.pageUuid || '',
    displayName: displayName || '',
  }
}

async function scanAllDisplays(pageName, projectUuid) {
  for (const displayType of [1, 2]) {
    const res = await displayModelList({ DisplayType: displayType })
    const list = res && res.data && res.data.list
    if (!list || res.data.code !== 0) continue
    for (const item of list) {
      const displayUid = item.displayUid || item.DisplayUid
      if (!displayUid) continue
      const hit = await findPageInDisplay(displayUid, pageName, item.name || item.Name)
      if (hit) {
        return { ...hit, projectUuid }
      }
    }
  }
  return null
}

async function resolveDashboardCandidates(store) {
  const out = []
  try {
    const res = await GetSystemHomeDashboard()
    if (res && res.data && res.data.code === 0 && res.data.dashboardUuid) {
      out.push({
        displayUid: res.data.dashboardUuid,
        displayName: res.data.dashboardName || '',
        projectUuid: res.data.projectUuid || '',
      })
    }
  } catch (e) { /* ignore */ }
  try {
    const home = getHomeDashboardConfig(store)
    if (home.dashboardUuid && !out.some(c => c.displayUid === home.dashboardUuid)) {
      out.push({
        displayUid: home.dashboardUuid,
        displayName: home.dashboardName || '',
        projectUuid: home.projectUuid || '',
      })
    }
  } catch (e) { /* ignore */ }
  return out
}

/**
 * 按组态 pageName 解析 displayUid + pageId（用于环境平面 / 2.5D 壳）
 */
export async function resolveDisplayPageByName(pageName) {
  const store = require('@/store').default
  await primeXunzhiJfRouteContext(store)
  let projectUuid = await ensureXunzhiJfProjectAuth()

  const dashboards = await resolveDashboardCandidates(store)
  for (const dash of dashboards) {
    const fromHome = await findPageInDisplay(dash.displayUid, pageName, dash.displayName)
    if (fromHome) {
      return { ...fromHome, projectUuid: projectUuid || dash.projectUuid }
    }
  }

  let hit = await scanAllDisplays(pageName, projectUuid)
  if (hit) return hit

  invalidateXunzhiJfProjectCache()
  projectUuid = await resolveXunzhiJfProjectUuid()
  if (projectUuid) {
    await ensureXunzhiJfProjectAuth()
    await store.dispatch('setting/fetchSystemHomeDashboard').catch(() => {})
    const dashboardsRetry = await resolveDashboardCandidates(store)
    for (const dash of dashboardsRetry) {
      const fromHome = await findPageInDisplay(dash.displayUid, pageName, dash.displayName)
      if (fromHome) {
        return { ...fromHome, projectUuid }
      }
    }
    hit = await scanAllDisplays(pageName, projectUuid)
    if (hit) return hit
  }

  return null
}
