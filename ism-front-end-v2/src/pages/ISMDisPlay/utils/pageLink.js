/**
 * 组态页跳转 / 弹窗：统一「是否弹窗」判定，避免 "false"/1/"true" 被当成不同语义。
 * 仅 true / 1 / "1" / "true" 视为弹窗。
 */
export function isPopUpEnabled(value) {
  return value === true || value === 1 || value === '1' || value === 'true'
}

export function resolveHomePageUuidFromList(pages, displayUUID) {
  if (!Array.isArray(pages) || !pages.length) {
    return ''
  }
  const ofModel = displayUUID
    ? pages.filter(function (p) {
      return p.pageModelUuid == displayUUID || p.modelId == displayUUID
    })
    : pages
  const pool = ofModel.length ? ofModel : pages
  const home = pool.find(function (p) { return p.IsHome == 1 })
  const hit = home || pool[0]
  return (hit && (hit.value || hit.pageUuid || hit.PageId)) || ''
}

/**
 * 菜单保存：若选中的是模型 UUID，换成该模型真实首页 pageUuid。
 */
export function resolveMenuPagePath(selectPage, displayUUID, pages) {
  if (!selectPage) {
    return selectPage
  }
  if (displayUUID && selectPage === displayUUID) {
    const homeId = resolveHomePageUuidFromList(pages, displayUUID)
    if (homeId) {
      return homeId
    }
  }
  return selectPage
}
