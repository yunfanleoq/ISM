/**
 * 单设备测点表模板页（ViewRealTable rowSource=navDatapoints）
 * 一行 = 一个测点；整表 = 当前选中设备的所有点位（分页）
 */

/** 与 scripts/bootstrap_device_signal_template_mysql.sh 一致 */
export const DEVICE_SIGNAL_TEMPLATE_PAGE_ID = '7c3e8f92a1b04d6e9f3c2a1b0d8e7f65'

export function resolveDeviceSignalTemplateId(templateMap, pageList, modelUuid) {
  const map = templateMap || {}
  const byModel = map.deviceByModel || {}
  if (modelUuid && byModel[modelUuid]) return byModel[modelUuid]
  if (map.deviceDefault) return map.deviceDefault
  if (map.device) return map.device

  const pages = pageList || []
  const byStableId = pages.find(p => p && (p.pageUuid || p.pageUUID) === DEVICE_SIGNAL_TEMPLATE_PAGE_ID)
  if (byStableId) return byStableId.pageUuid || byStableId.pageUUID

  const byName = pages.find(p => {
    const t = String((p && (p.title || p.pageName)) || '')
    return /模板[-_].*测点|模板[-_].*设备信号|设备测点表|信号层模板/.test(t)
  })
  if (byName) return byName.pageUuid || byName.pageUUID

  const byKind = pages.find(p => p && (p.templateKind === 'device' || p.TemplateKind === 'device'))
  if (byKind) return byKind.pageUuid || byKind.pageUUID

  return DEVICE_SIGNAL_TEMPLATE_PAGE_ID
}
