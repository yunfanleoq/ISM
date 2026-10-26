import config from '@/config'
import {ADMIN} from '@/config/default'
import {
  getHomeDashboardConfig,
  getHomeDashboardPath,
} from '@/config/homeDashboard'
import {formatFullPath} from '@/utils/i18n'
import {filterMenu} from '@/utils/authority-utils'
import {getLocalSetting} from '@/utils/themeUtil'
import deepClone from 'lodash.clonedeep'
import {GetSystemHomeDashboard, SetSystemHomeDashboard} from '@/services/system'

const localSetting = getLocalSetting(true)
const customTitlesStr = sessionStorage.getItem(process.env.VUE_APP_TBAS_TITLES_KEY)
const customTitles = (customTitlesStr && JSON.parse(customTitlesStr)) || []

export default {
  namespaced: true,
  state: {
    isMobile: false,
    animates: ADMIN.animates,
    componentAnimates: ADMIN.Canimates,
    palettes: ADMIN.palettes,
    pageMinHeight: 0,
    menuData: [],
    videoServer:location.host,
    langList: [
      {key: 'CN', name: '简体中文', alias: '简体'},
      {key: 'HK', name: '繁體中文', alias: '繁體'},
      {key: 'US', name: 'English', alias: 'English'}
    ],
    activatedFirst: undefined,
    customTitles,
    systemHomeDashboard: null,
    ...config,
    ...localSetting
  },
  getters: {
    menuData(state, getters, rootState, rootGetters) {
      if (state.filterMenu) {
        const permissions = rootGetters['account/permissions'] || []
        const roles = rootGetters['account/roles'] || []
        return filterMenu(deepClone(state.menuData), permissions, roles)
      }
      return state.menuData
    },
    firstMenu(state, getters) {
      const {menuData} = getters
      if (menuData.length > 0 && !menuData[0].fullPath) {
        formatFullPath(menuData)
      }
      return menuData.map(item => {
        const menuItem = {...item}
        delete menuItem.children
        return menuItem
      })
    },
    subMenu(state) {
      const {menuData, activatedFirst} = state
      if (menuData.length > 0 && !menuData[0].fullPath) {
        formatFullPath(menuData)
      }
      const current = menuData.find(menu => menu.fullPath === activatedFirst)
      return current && current.children || []
    },
    homeDashboardConfig(state) {
      return getHomeDashboardConfig({state: {setting: state}})
    },
    homeDashboardUuid(state, getters) {
      return getters.homeDashboardConfig.dashboardUuid
    },
    homeDashboardProjectUuid(state, getters) {
      return getters.homeDashboardConfig.projectUuid
    },
    homeDashboardPath(state, getters) {
      return getHomeDashboardPath({state: {setting: state}})
    },
  },
  mutations: {
    setDevice (state, isMobile) {
      state.isMobile = isMobile
    },
    setTheme (state, theme) {
      state.theme = theme
    },
    setLayout (state, layout) {
      state.layout = layout
    },
    setMultiPage (state, multiPage) {
      state.multiPage = multiPage
    },
    setAnimate (state, animate) {
      state.animate = animate
    },
    setWeekMode(state, weekMode) {
      state.weekMode = weekMode
    },
    setFixedHeader(state, fixedHeader) {
      state.fixedHeader = fixedHeader
    },
    setFixedSideBar(state, fixedSideBar) {
      state.fixedSideBar = fixedSideBar
    },
    setLang(state, lang) {
      state.lang = lang
      localStorage.setItem("lang", lang)
    },
    setHideSetting(state, hideSetting) {
      state.hideSetting = hideSetting
    },
    correctPageMinHeight(state, minHeight) {
      state.pageMinHeight += minHeight
    },
    setMenuData(state, menuData) {
      state.menuData = menuData
    },
    setAsyncRoutes(state, asyncRoutes) {
      state.asyncRoutes = asyncRoutes
    },
    setPageWidth(state, pageWidth) {
      state.pageWidth = pageWidth
    },
    setActivatedFirst(state, activatedFirst) {
      state.activatedFirst = activatedFirst
    },
    setFixedTabs(state, fixedTabs) {
      state.fixedTabs = fixedTabs
    },
    setCustomTitle(state, {path, title}) {
      if (title) {
        const obj = state.customTitles.find(item => item.path === path)
        if (obj) {
          obj.title = title
        } else {
          state.customTitles.push({path, title})
        }
        sessionStorage.setItem(process.env.VUE_APP_TBAS_TITLES_KEY, JSON.stringify(state.customTitles))
      }
    },
    setSystemHomeDashboard(state, payload) {
      state.systemHomeDashboard = payload || null
    },
    patchScadaMonitorMenuName(state, name) {
      if (!name) return
      const patch = (routes) => {
        if (!routes || !routes.length) return
        routes.forEach(route => {
          if (route.path === '/SCADAMonitor') {
            route.name = name
          }
          if (route.children && route.children.length) {
            patch(route.children)
          }
        })
      }
      patch(state.menuData)
    }
  },
  actions: {
    fetchSystemHomeDashboard({commit}) {
      return GetSystemHomeDashboard().then(res => {
        if (res.data && res.data.code === 0) {
          commit('setSystemHomeDashboard', {
            dashboardUuid: res.data.dashboardUuid || '',
            projectUuid: res.data.projectUuid || '',
            dashboardName: res.data.dashboardName || '',
          })
          // 侧栏固定「监控大屏」，不按应用名覆盖
          commit('patchScadaMonitorMenuName', '监控大屏')
        }
        return res
      })
    },
    saveSystemHomeDashboard({commit, dispatch}, payload) {
      return SetSystemHomeDashboard(payload).then(res => {
        if (res.data && res.data.code === 0) {
          commit('setSystemHomeDashboard', {
            dashboardUuid: payload.dashboardUuid,
            projectUuid: payload.projectUuid,
            dashboardName: '',
          })
          return dispatch('fetchSystemHomeDashboard').then(fetchRes => {
            if (fetchRes && fetchRes.data && fetchRes.data.code === 0) {
              return fetchRes
            }
            return res
          })
        }
        return res
      })
    }
  }
}
