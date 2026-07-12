/**
 * 组织层级运行态大屏的性能边界。
 * 所有海量数据入口都必须以这些上限分页，禁止在浏览器保留完整测点列表。
 */
export const DASHBOARD_PERFORMANCE = Object.freeze({
  // 1888×896 运行区按 10 列×8 行展示高密度实时卡片；始终服务端分页。
  datapointPageSize: 80,
  datapointMaxPageSize: 100,
  deviceSummaryPageSize: 30,
  historySeriesLimit: 1,
  historyPointLimit: 300,
  organizationCacheTtlMs: 15 * 1000,
})

export function normalizePageSize(value, fallback = DASHBOARD_PERFORMANCE.datapointPageSize) {
  const size = Number(value)
  if (!Number.isFinite(size) || size < 1) return fallback
  return Math.min(Math.floor(size), DASHBOARD_PERFORMANCE.datapointMaxPageSize)
}
