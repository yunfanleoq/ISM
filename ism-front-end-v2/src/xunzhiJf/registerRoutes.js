import BlankView from '@/layouts/BlankView'
import { xunzhiJfVizRouteBeforeEnter } from '@/xunzhiJf/projectContext'
import {
  XUNZHI_JF_MENU_FLAT,
  menuCodeToRouterKey,
  XUNZHI_JF_SECTION_ROUTER,
} from '@/xunzhiJf/menuFlat'

const GROUP_ICON = {
  'jf.monitor': 'dashboard',
  'jf.alarm': 'bell',
  'jf.viz': 'picture',
  'jf.asset': 'database',
  'jf.capacity': 'pie-chart',
  'jf.workorder': 'solution',
  'jf.energy': 'thunderbolt',
  'jf.security': 'safety',
  'jf.report': 'file-text',
  'jf.system': 'setting',
}

const WIRED = new Set([
  'jf.home',
  'jf.monitor.overview',
  'jf.alarm.active',
  'jf.alarm.history',
  'jf.viz.plane',
  'jf.viz.dashboard',
  'jf.viz.nav25d',
  'jf.viz.show',
  'jf.viz.bjga',
])

const BJGA_HASH = {
  'jf.viz.bjga': 'dashboard',
  'jf.asset.overview': 'assetmgr',
  'jf.asset.tree': 'building',
  'jf.asset.u': 'assetmgr/cabinet',
  'jf.asset.lifecycle': 'assetmgr/flow',
  'jf.asset.retire': 'assetmgr/inout',
  'jf.asset.multi': 'building',
  'jf.capacity.overview': 'capacity',
  'jf.capacity.reserve': 'capacity',
  'jf.capacity.suggest': 'capacity',
  'jf.workorder.list': 'ops',
  'jf.workorder.todo': 'ops/duty',
  'jf.workorder.kb': 'ops/kb',
  'jf.energy.power': 'energy/pue',
  'jf.energy.pue': 'energy/pue',
  'jf.energy.carbon': 'energy/carbon',
  'jf.security.video': 'facility/video',
  'jf.security.access': 'facility/door',
  'jf.security.fire': 'facility/fire',
  'jf.report.ops': 'report',
  'jf.report.custom': 'report',
  'jf.system.users': 'admin',
  'jf.system.scope': 'audit/policy',
  'jf.system.north': 'govern',
  'jf.system.audit': 'audit',
}

/** 菜单嵌套下的单段相对 path（禁止 leading `/`，否则脱离 root/TabsView） */
export function routeSegmentForRow(row) {
  if (row.route_or_page === '—') {
    const parts = row.menu_code.split('.')
    return parts[parts.length - 1]
  }
  if (row.route_or_page.startsWith('page:')) {
    if (row.menu_code === 'jf.viz.plane') return 'plane'
    if (row.menu_code === 'jf.viz.nav25d') return 'nav25d'
    return 'plane'
  }
  const rest = row.route_or_page.replace(/^\/xunzhi\/jf\//, '')
  const segs = rest.split('/').filter(Boolean)
  return segs[segs.length - 1] || row.menu_code.split('.').pop()
}

export function fullPathForRow(row) {
  if (row.route_or_page && row.route_or_page.startsWith('/')) {
    return row.route_or_page
  }
  if (row.route_or_page && row.route_or_page.startsWith('page:')) {
    if (row.menu_code === 'jf.viz.plane') return '/xunzhi/jf/viz/plane'
    if (row.menu_code === 'jf.viz.nav25d') return '/xunzhi/jf/viz/nav25d'
  }
  const segs = []
  let code = row.menu_code
  while (code && code !== 'jf') {
    const rowAt = XUNZHI_JF_MENU_FLAT.find(r => r.menu_code === code)
    if (rowAt) segs.unshift(routeSegmentForRow(rowAt))
    const dot = code.lastIndexOf('.')
    code = dot > 0 ? code.slice(0, dot) : ''
  }
  return `/xunzhi/jf/${segs.join('/')}`
}

function firstLeafSegment(menuCode) {
  const child = XUNZHI_JF_MENU_FLAT
    .filter(r => r.parent_code === menuCode && r.route_or_page && r.route_or_page !== '—')
    .sort((a, b) => a.sort - b.sort)[0]
  if (!child) return ''
  return routeSegmentForRow(child)
}

function skeletonMeta(row) {
  return {
    title: row.title,
    menuCode: row.menu_code,
    moduleId: row.module_id,
    perm: row.perm,
    routePath: fullPathForRow(row),
    sliceHint: row.sort < 500 ? 'S0/S1' : row.sort < 700 ? 'S2–S4' : 'S5/S6',
  }
}

function componentForRow(row) {
  if (BJGA_HASH[row.menu_code]) {
    return () => import('@/xunzhiJf/pages/BjgaScreenPage.vue')
  }
  switch (row.menu_code) {
    case 'jf.home':
      return () => import('@/xunzhiJf/pages/HomeWorkbench.vue')
    case 'jf.monitor.overview':
      return () => import('@/xunzhiJf/pages/MonitorOverview.vue')
    case 'jf.alarm.active':
      return () => import('@/xunzhiJf/pages/AlarmActiveHost.vue')
    case 'jf.alarm.history':
      return () => import('@/pages/reporting/alarmReport/alarmHistory')
    case 'jf.viz.plane':
      return () => import('@/xunzhiJf/pages/VizPlanePage.vue')
    case 'jf.viz.dashboard':
      return () => import('@/xunzhiJf/pages/VizDashboardPage.vue')
    case 'jf.viz.nav25d':
      return () => import('@/xunzhiJf/pages/VizNav25dPage.vue')
    case 'jf.viz.show':
      return () => import('@/xunzhiJf/pages/DemoShowPage.vue')
    default:
      return () => import('@/xunzhiJf/pages/SkeletonPage.vue')
  }
}

export function buildXunzhiJfRouterMap() {
  const map = {}

  map[XUNZHI_JF_SECTION_ROUTER] = {
    path: 'xunzhi/jf',
    name: XUNZHI_JF_SECTION_ROUTER,
    redirect: 'home',
    meta: {
      icon: 'cloud-server',
      title: '循智元晟',
      menuCode: 'jf.section',
      moduleId: 'envmon',
      authority: { role: ['Admin', 'Operator', 'User'] },
    },
    component: BlankView,
  }

  XUNZHI_JF_MENU_FLAT.forEach(row => {
    const key = menuCodeToRouterKey(row.menu_code)
    const isGroup = row.route_or_page === '—'
    const entry = {
      path: routeSegmentForRow(row),
      name: key,
      meta: {
        title: row.title,
        icon: GROUP_ICON[row.menu_code] || 'file',
        menuCode: row.menu_code,
        moduleId: row.module_id,
        xunzhiPerm: row.perm,
        authority: { role: ['Admin', 'Operator', 'User'] },
      },
    }
    if (isGroup) {
      entry.component = BlankView
      const seg = firstLeafSegment(row.menu_code)
      if (seg) entry.redirect = seg
    } else {
      entry.component = componentForRow(row)
      if (BJGA_HASH[row.menu_code]) {
        entry.meta.bjgaHash = BJGA_HASH[row.menu_code]
      }
      if (row.menu_code === 'jf.viz.plane' || row.menu_code === 'jf.viz.dashboard') {
        entry.beforeEnter = xunzhiJfVizRouteBeforeEnter
      }
      if (!WIRED.has(row.menu_code) && !BJGA_HASH[row.menu_code]) {
        entry.meta.skeletonProps = skeletonMeta(row)
      }
    }
    map[key] = entry
  })

  return map
}

export function applyXunzhiJfRouterMap(routerMap) {
  Object.assign(routerMap, buildXunzhiJfRouterMap())
}

export function buildXunzhiJfMenuConfigChildren() {
  const byParent = new Map()
  XUNZHI_JF_MENU_FLAT.forEach(row => {
    const p = row.parent_code
    if (!byParent.has(p)) byParent.set(p, [])
    byParent.get(p).push(row)
  })
  byParent.forEach(list => list.sort((a, b) => a.sort - b.sort))

  function build(parentCode) {
    return (byParent.get(parentCode) || []).map(row => {
      const node = {
        router: menuCodeToRouterKey(row.menu_code),
        name: row.title,
      }
      const kids = build(row.menu_code)
      if (kids.length) node.children = kids
      return node
    })
  }

  return {
    router: XUNZHI_JF_SECTION_ROUTER,
    name: '循智元晟',
    icon: 'cloud-server',
    children: build('(root)'),
  }
}
