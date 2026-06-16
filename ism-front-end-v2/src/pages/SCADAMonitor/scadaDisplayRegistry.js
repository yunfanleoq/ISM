/**
 * 组态 displayUUID → SCADA 原生大屏配置
 * 命中后 AppRun 路由渲染 SCADAMonitor，而非静态 ISMDisPlay cell 画布
 */
export const SCADA_DISPLAY_REGISTRY = {
  '043135ad-44be-e5d8-89be-3e54883c23a8': {
    title: '航信机房电力监控系统',
    subtitle: 'NCC ROOM POWER SCADA',
    projectUuid: '31bc90be-ebc4-dd61-ba9d-ce6e075e40e2',
  },
}

export function getScadaConfigForDisplay(displayUuid) {
  if (!displayUuid) return null
  return SCADA_DISPLAY_REGISTRY[displayUuid] || null
}
