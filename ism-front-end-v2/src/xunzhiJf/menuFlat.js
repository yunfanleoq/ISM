/**
 * §8 扁平表（架构 v1.0 正式锁）— 挂菜单唯一照抄源
 * @see docs/机房环境监测-菜单树与路由权限挂点.md
 */
export const XUNZHI_JF_MENU_FLAT = [
  { sort: 100, menu_code: 'jf.home', title: '工作台', parent_code: '(root)', module_id: 'envmon-viz', route_or_page: '/xunzhi/jf/home', perm: 'viz:read' },
  { sort: 200, menu_code: 'jf.monitor', title: '场地监测', parent_code: '(root)', module_id: 'envmon-collect', route_or_page: '—', perm: 'collect:read' },
  { sort: 210, menu_code: 'jf.monitor.overview', title: '实时总览', parent_code: 'jf.monitor', module_id: 'envmon-collect', route_or_page: '/xunzhi/jf/monitor/overview', perm: 'collect:read' },
  { sort: 220, menu_code: 'jf.monitor.power', title: '动力监测', parent_code: 'jf.monitor', module_id: 'envmon-collect', route_or_page: '/xunzhi/jf/monitor/power', perm: 'collect:read' },
  { sort: 230, menu_code: 'jf.monitor.env', title: '环境监测', parent_code: 'jf.monitor', module_id: 'envmon-collect', route_or_page: '/xunzhi/jf/monitor/env', perm: 'collect:read' },
  { sort: 240, menu_code: 'jf.monitor.battery', title: '蓄电池监测', parent_code: 'jf.monitor', module_id: 'envmon-collect', route_or_page: '/xunzhi/jf/monitor/battery', perm: 'collect:read' },
  { sort: 250, menu_code: 'jf.monitor.mm', title: '微模块监测', parent_code: 'jf.monitor', module_id: 'envmon-south', route_or_page: '/xunzhi/jf/monitor/micromodule', perm: 'collect:read' },
  { sort: 260, menu_code: 'jf.monitor.points', title: '测点/点表', parent_code: 'jf.monitor', module_id: 'envmon-collect', route_or_page: '/xunzhi/jf/monitor/points', perm: 'collect:config' },
  { sort: 300, menu_code: 'jf.alarm', title: '告警中心', parent_code: '(root)', module_id: 'envmon-alarm', route_or_page: '—', perm: 'alarm:read' },
  { sort: 310, menu_code: 'jf.alarm.active', title: '活动告警', parent_code: 'jf.alarm', module_id: 'envmon-alarm', route_or_page: '/xunzhi/jf/alarm/active', perm: 'alarm:read' },
  { sort: 320, menu_code: 'jf.alarm.history', title: '历史告警', parent_code: 'jf.alarm', module_id: 'envmon-alarm', route_or_page: '/xunzhi/jf/alarm/history', perm: 'alarm:read' },
  { sort: 330, menu_code: 'jf.alarm.policy', title: '告警策略', parent_code: 'jf.alarm', module_id: 'envmon-alarm', route_or_page: '/xunzhi/jf/alarm/policy', perm: 'alarm:config' },
  { sort: 400, menu_code: 'jf.viz', title: '可视化', parent_code: '(root)', module_id: 'envmon-viz', route_or_page: '—', perm: 'viz:read' },
  { sort: 410, menu_code: 'jf.viz.plane', title: '环境平面大屏', parent_code: 'jf.viz', module_id: 'envmon-viz', route_or_page: 'page:xunzhi-jf-env-plane', perm: 'viz:read' },
  { sort: 420, menu_code: 'jf.viz.dashboard', title: '平面组态大屏', parent_code: 'jf.viz', module_id: 'envmon-viz', route_or_page: '/xunzhi/jf/viz/dashboard', perm: 'viz:read' },
  { sort: 430, menu_code: 'jf.viz.nav25d', title: '2.5D 导航', parent_code: 'jf.viz', module_id: 'envmon-viz', route_or_page: '/xunzhi/jf/viz/nav25d', perm: 'viz:read' },
  { sort: 440, menu_code: 'jf.viz.editor', title: '组态编辑', parent_code: 'jf.viz', module_id: 'envmon-viz', route_or_page: '/xunzhi/jf/viz/editor', perm: 'viz:edit' },
  { sort: 450, menu_code: 'jf.viz.charts', title: '图表模板库', parent_code: 'jf.viz', module_id: 'envmon-viz', route_or_page: '/xunzhi/jf/viz/charts', perm: 'viz:read' },
  { sort: 500, menu_code: 'jf.asset', title: '资产管理', parent_code: '(root)', module_id: 'dcim-asset', route_or_page: '—', perm: 'asset:read' },
  { sort: 510, menu_code: 'jf.asset.overview', title: '资产总览', parent_code: 'jf.asset', module_id: 'dcim-asset', route_or_page: '/xunzhi/jf/asset/overview', perm: 'asset:read' },
  { sort: 520, menu_code: 'jf.asset.tree', title: '空间树', parent_code: 'jf.asset', module_id: 'dcim-asset', route_or_page: '/xunzhi/jf/asset/tree', perm: 'asset:read' },
  { sort: 530, menu_code: 'jf.asset.u', title: 'U位管理', parent_code: 'jf.asset', module_id: 'dcim-asset', route_or_page: '/xunzhi/jf/asset/u-slot', perm: 'asset:read' },
  { sort: 540, menu_code: 'jf.asset.lifecycle', title: '入库上架变更', parent_code: 'jf.asset', module_id: 'dcim-asset', route_or_page: '/xunzhi/jf/asset/lifecycle', perm: 'asset:write' },
  { sort: 550, menu_code: 'jf.asset.retire', title: '下架报废', parent_code: 'jf.asset', module_id: 'dcim-asset', route_or_page: '/xunzhi/jf/asset/retire', perm: 'asset:write' },
  { sort: 560, menu_code: 'jf.asset.multi', title: '机房2/3纳管', parent_code: 'jf.asset', module_id: 'dcim-asset', route_or_page: '/xunzhi/jf/asset/sites', perm: 'asset:read' },
  { sort: 600, menu_code: 'jf.capacity', title: '容量管理', parent_code: '(root)', module_id: 'dcim-capacity', route_or_page: '—', perm: 'capacity:read' },
  { sort: 610, menu_code: 'jf.capacity.overview', title: '容量水位', parent_code: 'jf.capacity', module_id: 'dcim-capacity', route_or_page: '/xunzhi/jf/capacity/overview', perm: 'capacity:read' },
  { sort: 620, menu_code: 'jf.capacity.reserve', title: '预占管理', parent_code: 'jf.capacity', module_id: 'dcim-capacity', route_or_page: '/xunzhi/jf/capacity/reserve', perm: 'capacity:reserve' },
  { sort: 630, menu_code: 'jf.capacity.suggest', title: '推荐上架', parent_code: 'jf.capacity', module_id: 'dcim-capacity', route_or_page: '/xunzhi/jf/capacity/suggest', perm: 'capacity:read' },
  { sort: 700, menu_code: 'jf.workorder', title: '运维工单', parent_code: '(root)', module_id: 'dcim-workorder', route_or_page: '—', perm: 'workorder:read' },
  { sort: 710, menu_code: 'jf.workorder.list', title: '工单列表', parent_code: 'jf.workorder', module_id: 'dcim-workorder', route_or_page: '/xunzhi/jf/workorder/list', perm: 'workorder:read' },
  { sort: 720, menu_code: 'jf.workorder.todo', title: '我的待办', parent_code: 'jf.workorder', module_id: 'dcim-workorder', route_or_page: '/xunzhi/jf/workorder/todo', perm: 'workorder:read' },
  { sort: 730, menu_code: 'jf.workorder.kb', title: '知识库', parent_code: 'jf.workorder', module_id: 'dcim-workorder', route_or_page: '/xunzhi/jf/workorder/kb', perm: 'workorder:read' },
  { sort: 800, menu_code: 'jf.energy', title: '能耗管理', parent_code: '(root)', module_id: 'dcim-energy', route_or_page: '—', perm: 'energy:read' },
  { sort: 810, menu_code: 'jf.energy.power', title: '用电功率分项', parent_code: 'jf.energy', module_id: 'dcim-energy', route_or_page: '/xunzhi/jf/energy/power', perm: 'energy:read' },
  { sort: 820, menu_code: 'jf.energy.pue', title: 'PUE分析', parent_code: 'jf.energy', module_id: 'dcim-energy', route_or_page: '/xunzhi/jf/energy/pue', perm: 'energy:read' },
  { sort: 830, menu_code: 'jf.energy.carbon', title: '碳排放', parent_code: 'jf.energy', module_id: 'dcim-energy', route_or_page: '/xunzhi/jf/energy/carbon', perm: 'energy:read' },
  { sort: 900, menu_code: 'jf.security', title: '安防消防', parent_code: '(root)', module_id: 'envmon-alarm', route_or_page: '—', perm: 'alarm:read' },
  { sort: 910, menu_code: 'jf.security.video', title: '视频点位', parent_code: 'jf.security', module_id: 'envmon-alarm', route_or_page: '/xunzhi/jf/security/video', perm: 'alarm:read' },
  { sort: 920, menu_code: 'jf.security.access', title: '门禁入侵', parent_code: 'jf.security', module_id: 'envmon-alarm', route_or_page: '/xunzhi/jf/security/access', perm: 'alarm:read' },
  { sort: 930, menu_code: 'jf.security.fire', title: '消防联动', parent_code: 'jf.security', module_id: 'envmon-alarm', route_or_page: '/xunzhi/jf/security/fire', perm: 'alarm:read' },
  { sort: 1000, menu_code: 'jf.report', title: '报表中心', parent_code: '(root)', module_id: 'envmon-viz', route_or_page: '—', perm: 'viz:read' },
  { sort: 1010, menu_code: 'jf.report.ops', title: '运行报表', parent_code: 'jf.report', module_id: 'envmon-viz', route_or_page: '/xunzhi/jf/report/ops', perm: 'viz:read' },
  { sort: 1020, menu_code: 'jf.report.custom', title: '自定义报表', parent_code: 'jf.report', module_id: 'envmon-viz', route_or_page: '/xunzhi/jf/report/custom', perm: 'viz:edit' },
  { sort: 1100, menu_code: 'jf.system', title: '系统管理', parent_code: '(root)', module_id: 'envmon-iam', route_or_page: '—', perm: 'iam:user' },
  { sort: 1110, menu_code: 'jf.system.users', title: '用户与角色', parent_code: 'jf.system', module_id: 'envmon-iam', route_or_page: '/xunzhi/jf/system/users', perm: 'iam:role' },
  { sort: 1120, menu_code: 'jf.system.scope', title: '数据范围', parent_code: 'jf.system', module_id: 'envmon-iam', route_or_page: '/xunzhi/jf/system/scope', perm: 'iam:scope' },
  { sort: 1130, menu_code: 'jf.system.north', title: '北向接口', parent_code: 'jf.system', module_id: 'envmon-north', route_or_page: '/xunzhi/jf/system/north', perm: 'collect:config' },
  { sort: 1140, menu_code: 'jf.system.audit', title: '审计日志', parent_code: 'jf.system', module_id: 'envmon-iam', route_or_page: '/xunzhi/jf/system/audit', perm: 'iam:user' },
]

export function menuCodeToRouterKey(menuCode) {
  return 'XunzhiJf_' + menuCode.replace(/\./g, '_')
}

export const XUNZHI_JF_SECTION_ROUTER = 'XunzhiJfSection'
