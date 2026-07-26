/**
 * 单设备测点表模板页（ViewRealTable rowSource=navDatapoints）
 * 一行 = 一个测点；整表 = 当前选中设备的所有点位（分页）
 */

/** 点位列表只允许使用后端声明的 datapointList 模板；muid 仅用于查点位。 */
export function resolveDeviceSignalTemplateId(templateMap) {
  const map = templateMap || {}
  return map.datapointList || ''
}
