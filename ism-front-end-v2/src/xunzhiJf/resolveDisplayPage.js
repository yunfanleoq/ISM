import { displayModelList } from '@/services/displayModel'
import { getDisplayModelLayerData } from '@/services/displayModel'
import { ensureXunzhiJfProjectAuth } from '@/xunzhiJf/projectContext'

export const XUNZHI_JF_PROJECT_NAME_HINT = '寻智机房环境-xunzhi-jf'
export const XUNZHI_JF_ENV_PLANE_PAGE = 'xunzhi-jf-env-plane'
export const XUNZHI_JF_NAV_25D_PAGE = 'xunzhi-jf-25d-shell'

/**
 * 按组态 pageName 解析 displayUid + pageId（用于环境平面 / 2.5D 壳）
 */
export async function resolveDisplayPageByName(pageName) {
  const projectUuid = await ensureXunzhiJfProjectAuth()
  for (const displayType of [1, 2]) {
    const res = await displayModelList({ DisplayType: displayType })
    const list = res && res.data && res.data.list
    if (!list) continue
    for (const item of list) {
      const displayUid = item.displayUid
      if (!displayUid) continue
      const layerRes = await getDisplayModelLayerData({ muid: displayUid })
      if (!layerRes || !layerRes.data || layerRes.data.code !== 0) continue
      const layers = layerRes.data.layer || []
      const page = layers.find(l => l.PageName === pageName && l.IsLogin !== 1)
      if (page) {
        return {
          displayUid,
          pageId: page.PageId,
          projectUuid,
          displayName: item.name,
        }
      }
    }
  }
  return null
}
