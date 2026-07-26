<template>
  <svg
    xmlns="http://www.w3.org/2000/svg"
    preserveAspectRatio="none"
    x="0px"
    y="0px"
    xml:space="preserve"
    :style="{ overflow: 'visible', width: (detail && detail.style && detail.style.position && detail.style.position.w) || width || 520, height: (detail && detail.style && detail.style.position && detail.style.position.h) || height || 200 }"
  >
    <g
      class="svg-el"
      :class="{ animated: true, [`${detail && detail.style && detail.style.animate}`]: true }"
      :style="{
        opacity: fillOpacity,
        'stroke-opacity': strokeOpacity,
        stroke: strokeColor,
        'stroke-width': strokeWidth,
        'stroke-linecap': 'round',
        'stroke-linejoin': 'round',
        fill
      }"
    >
      <foreignObject style="overflow:visible;" pointer-events="all" :width="(detail && detail.style && detail.style.position && detail.style.position.w) || width || 520" :height="(detail && detail.style && detail.style.position && detail.style.position.h) || height || 200">
        <runtime-data-card-grid
          v-if="isRuntimeCardMode"
          :style="styleVar"
          :mode="isDeviceCardMode ? 'device' : 'point'"
          :items="runtimeCards"
          @select="openDeviceDatapoints"
        />
        <div v-else :class="['history-theme-shell', scrollbarThemeClass]" :style="styleVar">
          <div class="table-container">
            <div
              class="table-scroll"
              ref="tableScroll"
            >
              <a-table
                :columns="dynamicColumns"
                :data-source="pagedData"
                :rowKey="(record, index) => `${record.rowName}-${index}`"
                :pagination="false"
                :style="{ minWidth: `${tableScrollWidth}px` }"
              />
            </div>
          </div>
        </div>
      </foreignObject>
      <animate
        v-if="isStart && animateType.includes('blink') && !IsToolBox"
        attributeName="opacity"
        values="0.1;1;0.1"
        :dur="blinkSpeed + 's'"
        repeatCount="indefinite"
      />
      <animate
        v-if="isStart && animateType.includes('millcolorGrad') && !IsToolBox"
        attributeName="fill"
        :values="startColor + ';' + stopColor + ';' + startColor"
        :dur="animateSpeed + 's'"
        repeatCount="indefinite"
      />
      <animateTransform
        v-if="isStart && animateType.includes('Zoom') && !IsToolBox"
        attributeName="transform"
        begin="0s"
        dur="0.6s"
        type="scale"
        values="0.9;1;0.9"
        repeatCount="indefinite"
      />
      <animateTransform
        v-if="isStart && animateType.includes('animateSpin') && !IsToolBox && spinDirection == 0"
        attributeType="XML"
        attributeName="transform"
        :dur="animateSpinSpeed + 's'"
        type="rotate"
        from="0 0 0"
        to="360 0 0"
        repeatCount="indefinite"
      />
      <animateTransform
        v-if="isStart && animateType.includes('animateSpin') && !IsToolBox && spinDirection == 1"
        attributeType="XML"
        attributeName="transform"
        :dur="animateSpinSpeed + 's'"
        type="rotate"
        from="360 0 0"
        to="0 0 0"
        repeatCount="indefinite"
      />
    </g>
  </svg>
</template>

<script>
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin';
import RuntimeDataCardGrid from './RuntimeDataCardGrid.vue'
import { formatPointDisplayValue, isDeviceOnlineFromStatusValue } from '@/pages/ISMDisPlay/utils/pointValueDisplay'

const DEVICE_ONLINE_POLL_MS = 8000

// 勿静态 import @/store / @/services/device / @/utils/realDataBatch ——
// 主 chunk 已引用它们；ism-render 异步 chunk 再静态引用会触发 webpack4 scope-hoisting
// unused reexport bug，令本组件 export.default.data 丢失。
const REAL_DATA_DEFAULT_PAGE_SIZE = 30
const REAL_DATA_MAX_PAGE_SIZE = 100

function clampPageSize(size, fallback = REAL_DATA_DEFAULT_PAGE_SIZE) {
  const n = parseInt(size, 10)
  if (!Number.isFinite(n) || n < 1) {
    return fallback
  }
  return Math.min(n, REAL_DATA_MAX_PAGE_SIZE)
}

function chunkArray(list, size = REAL_DATA_DEFAULT_PAGE_SIZE) {
  const arr = Array.isArray(list) ? list : []
  const chunkSize = clampPageSize(size)
  const out = []
  for (let i = 0; i < arr.length; i += chunkSize) {
    out.push(arr.slice(i, i + chunkSize))
  }
  return out
}

function readDiyValue(detail, key) {
  const diy = (detail && detail.style && detail.style.diy) || []
  const item = diy.find(d => d && d.key === key)
  return item ? item.value : undefined
}

function stripDevicePrefix(pointName, deviceName) {
  const point = String(pointName || '').trim()
  const device = String(deviceName || '').trim()
  if (!point) return point
  if (device && point !== device && point.startsWith(device)) {
    const suffix = point.slice(device.length)
    // 仅当前缀后紧跟明确分隔符时裁剪，避免误伤“设备A区”这类相似名称。
    if (/^[\s_\-.:：/\\|·]+/u.test(suffix)) {
      return suffix.replace(/^[\s_\-.:：/\\|·]+/u, '').trim() || point
    }
  }
  // RC08bate：最后一个 '_' 前为设备名、后为测点名
  const idx = point.lastIndexOf('_')
  if (idx > 0 && idx < point.length - 1) {
    return point.slice(idx + 1).trim() || point
  }
  return point
}

/** 惰性取 vuex：只用 window 挂载，禁止 require('@/store')（会触发跨 chunk 导出残缺） */
function getEditorNavContext() {
  try {
    const st = typeof window !== 'undefined' ? window.__ISM_STORE__ : null
    return (st && st.state && st.state.ISMDisPlayEditorTool
      && st.state.ISMDisPlayEditorTool.navContext) || null
  } catch (e) {
    return null
  }
}

function commitEditorNav(partial) {
  try {
    const st = typeof window !== 'undefined' ? window.__ISM_STORE__ : null
    if (!st) return
    const cur = st.state.ISMDisPlayEditorTool && st.state.ISMDisPlayEditorTool.navContext
    // 离开设备页后首页会主动清空 navContext；晚到的实时请求不得重新创建旧设备上下文。
    if (!cur) return
    st.commit('ISMDisPlayEditorTool/setNavContext', { ...cur, ...partial })
  } catch (e) { /* ignore */ }
}

function postGetRealData(params) {
  // 与设备管理同源：/api/GetRealData；禁止 require @/services/device
  const headers = { 'Content-Type': 'application/json' }
  try {
    const raw = document.cookie || ''
    const pickCookie = (name) => {
      const m = raw.match(new RegExp('(?:(?:^|;\\s*)' + name + '=([^;]*))'))
      return m ? decodeURIComponent(m[1]) : ''
    }
    const token = pickCookie('Authorization') || pickCookie('authorization')
    if (token) {
      headers.Authorization = String(token).indexOf('Bearer') === 0 ? token : `Bearer ${token}`
    }
    const projectUuid = pickCookie('ProjectUuid') || pickCookie('projectUuid')
      || (typeof sessionStorage !== 'undefined' && (sessionStorage.getItem('ProjectUuid') || sessionStorage.getItem('projectUuid')))
      || ''
    if (projectUuid) {
      headers.ProjectUuid = projectUuid
    }
  } catch (e) { /* ignore */ }
  return fetch('/api/getRealData', {
    method: 'POST',
    headers,
    body: JSON.stringify(params || {}),
    credentials: 'same-origin',
  }).then((resp) => resp.json().then((data) => ({ data })))
}

function postRealDataByBindings(bindings) {
  // 禁止 require @/services/device / axios / js-cookie（避免再引入主 chunk 依赖边）
  const headers = { 'Content-Type': 'application/json' }
  try {
    const raw = document.cookie || ''
    const pickCookie = (name) => {
      const m = raw.match(new RegExp('(?:(?:^|;\\s*)' + name + '=([^;]*))'))
      return m ? decodeURIComponent(m[1]) : ''
    }
    const token = pickCookie('Authorization') || pickCookie('authorization')
    if (token) {
      headers.Authorization = String(token).indexOf('Bearer') === 0 ? token : `Bearer ${token}`
    }
    const projectUuid = pickCookie('ProjectUuid') || pickCookie('projectUuid')
      || (typeof sessionStorage !== 'undefined' && (sessionStorage.getItem('ProjectUuid') || sessionStorage.getItem('projectUuid')))
      || ''
    if (projectUuid) {
      headers.ProjectUuid = projectUuid
    }
  } catch (e) { /* ignore */ }
  return fetch('/api/GetRealDataByBindings', {
    method: 'POST',
    headers,
    body: JSON.stringify({ bindings }),
    credentials: 'same-origin',
  }).then((resp) => {
    // 后端可能 gzip；浏览器会自动解压 Content-Encoding
    return resp.json().then((data) => ({ data }))
  })
}

const THEME_MAP = {
  light: {
    panelBg: 'linear-gradient(180deg, #f8fbff 0%, #eef4fb 100%)',
    panelBorder: '#d9e7f5',
    panelShadow: '0 10px 30px rgba(48, 86, 132, 0.10)',
    toolbarText: '#1f2a37',
    toolbarAccent: '#2f6fed',
    toolbarAccentSoft: 'rgba(47, 111, 237, 0.12)',
    toolbarAccentBorder: '#9dbcf7',
    toolbarAccentText: '#1f4fc9',
    tableColumnSplitColor: 'rgba(15, 23, 42, 0.08)',
    tableHeaderColor: '#16324f',
    tableHeaderBackColor: '#e7f0fb',
    tableHeaderFont: 'Arial',
    tableSplitColor: '#d7e2ee',
    tableHoverColor: '#dbeafe',
    tableRowOddBg: 'rgba(255, 255, 255, 0.94)',
    tableRowEvenBg: 'rgba(243, 248, 253, 0.96)',
    searchColor: '#16324f',
    searchBackColor: '#ffffff',
    searchBorderColor: '#bfd1e5',
    scrollBgColor: '#dce8f5',
    scrollFrColor: '#8eb3da',
    scrollHdColor: '#5e8dbf',
    foreColor: '#1f2937',
    backColor: 'transparent'
  },
  dark: {
    panelBg: 'linear-gradient(180deg, #111827 0%, #0f172a 100%)',
    panelBorder: '#263449',
    panelShadow: '0 14px 40px rgba(0, 0, 0, 0.35)',
    toolbarText: '#e5eefc',
    toolbarAccent: '#60a5fa',
    toolbarAccentSoft: 'rgba(96, 165, 250, 0.18)',
    toolbarAccentBorder: '#4c7fbe',
    toolbarAccentText: '#d6e8ff',
    tableColumnSplitColor: 'rgba(255, 255, 255, 0.18)',
    tableHeaderColor: '#f8fbff',
    tableHeaderBackColor: '#1d3557',
    tableHeaderFont: 'Arial',
    tableSplitColor: '#263449',
    tableHoverColor: '#1e3a5f',
    tableRowOddBg: 'rgba(17, 24, 39, 0.94)',
    tableRowEvenBg: 'rgba(15, 23, 42, 0.98)',
    searchColor: '#e5eefc',
    searchBackColor: '#162033',
    searchBorderColor: '#314158',
    scrollBgColor: '#172132',
    scrollFrColor: '#3b82f6',
    scrollHdColor: '#60a5fa',
    foreColor: '#e5eefc',
    backColor: 'transparent'
  },
  ocean: {
    panelBg: 'linear-gradient(180deg, #f2fbff 0%, #dcf3fb 100%)',
    panelBorder: '#9fd4e2',
    panelShadow: '0 12px 32px rgba(20, 102, 128, 0.16)',
    toolbarText: '#0f3d4c',
    toolbarAccent: '#0891b2',
    toolbarAccentSoft: 'rgba(8, 145, 178, 0.14)',
    toolbarAccentBorder: '#86d1e3',
    toolbarAccentText: '#0b6b84',
    tableColumnSplitColor: 'rgba(15, 23, 42, 0.08)',
    tableHeaderColor: '#ffffff',
    tableHeaderBackColor: '#0e7490',
    tableHeaderFont: 'Arial',
    tableSplitColor: '#a8d8e4',
    tableHoverColor: '#c7eef7',
    tableRowOddBg: 'rgba(255, 255, 255, 0.92)',
    tableRowEvenBg: 'rgba(232, 248, 252, 0.96)',
    searchColor: '#0f3d4c',
    searchBackColor: '#ffffff',
    searchBorderColor: '#8fcad9',
    scrollBgColor: '#d5eef4',
    scrollFrColor: '#22a6c3',
    scrollHdColor: '#0e7490',
    foreColor: '#134152',
    backColor: 'transparent'
  },
  amber: {
    panelBg: 'linear-gradient(180deg, #fff9ed 0%, #fff1d6 100%)',
    panelBorder: '#efc98d',
    panelShadow: '0 12px 28px rgba(180, 113, 20, 0.18)',
    toolbarText: '#6b3f12',
    toolbarAccent: '#d97706',
    toolbarAccentSoft: 'rgba(217, 119, 6, 0.14)',
    toolbarAccentBorder: '#efbf79',
    toolbarAccentText: '#a45705',
    tableColumnSplitColor: 'rgba(91, 55, 20, 0.10)',
    tableHeaderColor: '#fffdf8',
    tableHeaderBackColor: '#b45309',
    tableHeaderFont: 'Arial',
    tableSplitColor: '#ecd6b1',
    tableHoverColor: '#fde6bf',
    tableRowOddBg: 'rgba(255, 253, 247, 0.95)',
    tableRowEvenBg: 'rgba(255, 247, 230, 0.96)',
    searchColor: '#6b3f12',
    searchBackColor: '#fffdf8',
    searchBorderColor: '#e7c48f',
    scrollBgColor: '#fae7c5',
    scrollFrColor: '#e19a2b',
    scrollHdColor: '#b45309',
    foreColor: '#5b3714',
    backColor: 'transparent'
  },
  emerald: {
    panelBg: 'linear-gradient(180deg, #f2fdf7 0%, #e0f7ea 100%)',
    panelBorder: '#9dd9b4',
    panelShadow: '0 12px 30px rgba(17, 94, 67, 0.16)',
    toolbarText: '#114b36',
    toolbarAccent: '#059669',
    toolbarAccentSoft: 'rgba(5, 150, 105, 0.14)',
    toolbarAccentBorder: '#8fd5bb',
    toolbarAccentText: '#0d6d4d',
    tableColumnSplitColor: 'rgba(17, 75, 54, 0.10)',
    tableHeaderColor: '#f7fffb',
    tableHeaderBackColor: '#047857',
    tableHeaderFont: 'Arial',
    tableSplitColor: '#b5e0c4',
    tableHoverColor: '#d4f5df',
    tableRowOddBg: 'rgba(255, 255, 255, 0.92)',
    tableRowEvenBg: 'rgba(236, 251, 241, 0.96)',
    searchColor: '#114b36',
    searchBackColor: '#ffffff',
    searchBorderColor: '#9fd2b1',
    scrollBgColor: '#d7f0df',
    scrollFrColor: '#21a37a',
    scrollHdColor: '#047857',
    foreColor: '#184936',
    backColor: 'transparent'
  }
};

export default {
  mixins: [ISMChildAutoMixin],
  name: 'ism-view-real-table',
  components: { RuntimeDataCardGrid },
  inject: ['getNode'],
  props: {},
  watch: {
    detail: {
      handler(newVal) {
        if (this.editMode) {
          this.initComponents(newVal);
        }
      },
      deep: true
    },
    selectedTheme() {
      this.applyScrollbarTheme();
    },
    navDatapointPagination: {
      handler(val, oldVal) {
        if (!val || !oldVal || this.editMode || this.IsToolBox) return
        if ((val.datapointPageIndex || 0) === (oldVal.datapointPageIndex || 0)) return
        this.initComponents(this.detail)
        this.$nextTick(() => {
          if (!this.editMode && !this.IsToolBox) {
            this.QueryRealData()
          }
        })
      },
      deep: true
    },
    navListPagination: {
      handler(val, oldVal) {
        if (!val || !oldVal || this.editMode || this.IsToolBox) return
        if ((val.pageIndex || 0) === (oldVal.pageIndex || 0)) return
        this.initComponents(this.detail)
        this.$nextTick(() => {
          if (!this.editMode && !this.IsToolBox) {
            this.QueryRealData()
            this.refreshDeviceOnlineStatus()
          }
        })
      },
      deep: true
    },
    isDeviceCardMode(val) {
      if (val && !this.editMode && !this.IsToolBox) this.startDeviceOnlinePolling()
      else this.stopDeviceOnlinePolling()
    },
  },
  computed: {
    currentTheme() {
      return THEME_MAP[this.selectedTheme] || THEME_MAP.light;
    },
    autoPageEnabled() {
      // 自动翻页已下线：页码 chrome 与总数在自动翻页时不同步（见 ViewSvgText 根因）。
      // 保留 diy 配置与 timer 代码便于后续修复后重开；运行态一律关闭。
      return false
    },
    autoPageIntervalMs() {
      const seconds = Number(readDiyValue(this.detail, 'AutoPageInterval'))
      return Math.max(1, Number.isFinite(seconds) ? seconds : 5) * 1000
    },
    autoPageResumeDelayMs() {
      const seconds = Number(readDiyValue(this.detail, 'AutoPageResumeDelay'))
      return Math.max(1, Number.isFinite(seconds) ? seconds : 60) * 1000
    },
    scrollbarThemeClass() {
      const id = this.detail && this.detail.identifier ? this.detail.identifier : 'default';
      return `real-table-scrollbar-${id}`;
    },
    tableScrollWidth() {
      return this.dynamicColumns.reduce((total, column) => total + (Number(column.width) || 120), 0);
    },
    styleVar() {
      const style = (this.detail && this.detail.style) || {}
      const configuredBg = style.backColor && style.backColor !== 'transparent'
        ? style.backColor
        : this.currentTheme.panelBg
      const borderWidth = Number(style.borderWidth)
      const borderRadius = Number(style.borderRadius)
      const accent = readDiyValue(this.detail, 'panelAccentColor') || this.currentTheme.toolbarAccent
      return {
        '--panelBg': configuredBg,
        '--panelBorder': style.borderColor || this.currentTheme.panelBorder,
        '--panelBorderWidth': `${Number.isFinite(borderWidth) ? borderWidth : 1}px`,
        '--panelBorderRadius': `${Number.isFinite(borderRadius) ? borderRadius : 10}px`,
        '--panelAccent': accent,
        '--deviceIconAccent': readDiyValue(this.detail, 'deviceIconAccent') || '#52e8ff',
        '--pointIconAccent': readDiyValue(this.detail, 'pointIconAccent') || '#bca5ff',
        '--panelShadow': this.currentTheme.panelShadow,
        '--toolbarText': this.currentTheme.toolbarText,
        '--toolbarAccent': this.currentTheme.toolbarAccent,
        '--toolbarAccentSoft': this.currentTheme.toolbarAccentSoft,
        '--toolbarAccentBorder': this.currentTheme.toolbarAccentBorder,
        '--toolbarAccentText': this.currentTheme.toolbarAccentText,
        '--tableHeaderColor': this.currentTheme.tableHeaderColor,
        '--tableHeaderBackColor': this.currentTheme.tableHeaderBackColor,
        '--tableHeaderFont': this.tableHeaderFont || this.currentTheme.tableHeaderFont,
        '--tableHeaderFontSize': `${this.tableHeaderFontSize || 12}px`,
        '--tableSplitColor': this.currentTheme.tableSplitColor,
        '--tableColumnSplitColor': this.currentTheme.tableColumnSplitColor,
        '--tableHoverColor': this.currentTheme.tableHoverColor,
        '--tableRowOddBg': this.currentTheme.tableRowOddBg,
        '--tableRowEvenBg': this.currentTheme.tableRowEvenBg,
        '--SearchColor': this.currentTheme.searchColor,
        '--SearchBackColor': this.currentTheme.searchBackColor,
        '--SearchBorderColor': this.currentTheme.searchBorderColor,
        '--scrollBgColor': this.currentTheme.scrollBgColor,
        '--scrollFrColor': this.currentTheme.scrollFrColor,
        '--scrollHdColor': this.currentTheme.scrollHdColor,
        '--fontFamily': this.fontFamily,
        '--fontSize': `${this.fontSize || 12}px`,
        '--foreColor': this.foreColor || this.currentTheme.foreColor,
        '--backColor': this.backColor || this.currentTheme.backColor,
        width: `${(this.detail && this.detail.style && this.detail.style.position && this.detail.style.position.w) || this.width || 520}px`,
        height: `${(this.detail && this.detail.style && this.detail.style.position && this.detail.style.position.h) || this.height || 200}px`
      };
    },
    dynamicData() {
      const data = [];
      const rowCount = Math.max(
        this.rowDeviceNames.length,
        this.rowDeviceCodes.length,
        this.bindingMatrix.length,
        this.cellData.length
      );
      for (let rowIndex = 0; rowIndex < rowCount; rowIndex += 1) {
        const rowData = {
          rowName: this.rowDeviceNames[rowIndex] || '',
          deviceCode: this.rowDeviceCodes[rowIndex] || ''
        };
        this.columnHeaders.forEach((_, colIndex) => {
          const rowCellData = this.cellData[rowIndex];
          let cellValue = '-';
          if (rowCellData && rowCellData[colIndex] !== undefined && rowCellData[colIndex] !== null) {
            cellValue = formatPointDisplayValue(rowData.rowName, rowCellData[colIndex]);
          }
          rowData[`col_${colIndex}`] = cellValue;
        });
        data.push(rowData);
      }
      return data;
    },
    pagedData() {
      const navDp = this.navDatapointPagination
      if (navDp) {
        return this.dynamicData
      }
      const nav = this.navListPagination
      if (nav) {
        return this.dynamicData
      }
      const current = this.pagination && this.pagination.current ? this.pagination.current : 1;
      const pageSize = this.pagination && this.pagination.pageSize ? this.pagination.pageSize : REAL_DATA_DEFAULT_PAGE_SIZE;
      const start = (current - 1) * pageSize;
      const end = start + pageSize;
      return this.dynamicData.slice(start, end);
    },
    paginationTotal() {
      if (this.isDeviceCardMode) {
        return this.cardDeviceTotal
      }
      const navDp = this.navDatapointPagination
      if (navDp && navDp.totalDatapoints) {
        return navDp.totalDatapoints
      }
      if (this.isNavDatapointsSource) {
        const fromDiy = Number(readDiyValue(this.detail, 'navTotalDatapoints'))
        if (fromDiy > 0) return fromDiy
      }
      const nav = this.navListPagination
      if (nav && nav.totalDevices) {
        return nav.totalDevices
      }
      const explicitTotal = this.pagination && Number(this.pagination.total);
      if (explicitTotal && explicitTotal > 0) {
        return explicitTotal;
      }
      return this.dynamicData.length;
    },
    pagerPageSize() {
      if (this.isDeviceCardMode) {
        return this.cardPageSize
      }
      const navDp = this.navDatapointPagination
      if (navDp) {
        return navDp.datapointPageSize || REAL_DATA_DEFAULT_PAGE_SIZE
      }
      if (this.isNavDatapointsSource) {
        const fromDiy = Number(readDiyValue(this.detail, 'navDatapointPageSize'))
        if (fromDiy > 0) return fromDiy
        const showCount = Number(readDiyValue(this.detail, 'ShowCount'))
        if (showCount > 0) return showCount
      }
      const nav = this.navListPagination
      if (nav) {
        return nav.pageSize || REAL_DATA_DEFAULT_PAGE_SIZE
      }
      return (this.pagination && this.pagination.pageSize) || REAL_DATA_DEFAULT_PAGE_SIZE
    },
    pagerCurrent() {
      if (this.isDeviceCardMode) {
        return this.cardCurrentIndex + 1
      }
      const navDp = this.navDatapointPagination
      if (navDp) {
        return (navDp.datapointPageIndex || 0) + 1
      }
      if (this.isNavDatapointsSource) {
        const fromDiy = Number(readDiyValue(this.detail, 'navDatapointPageIndex'))
        if (Number.isFinite(fromDiy) && fromDiy >= 0) return fromDiy + 1
      }
      const nav = this.navListPagination
      if (nav) {
        return (nav.pageIndex || 0) + 1
      }
      return (this.pagination && this.pagination.current) || 1
    },
    pagerTotalPages() {
      // 以总数/页大小为准，避免 store/diy 里陈旧的 datapointTotalPages=1 把下一页锁死
      const total = Number(this.paginationTotal) || 0
      const size = Math.max(1, Number(this.pagerPageSize) || 20)
      if (total > 0) {
        return Math.max(1, Math.ceil(total / size))
      }
      const navDp = this.navDatapointPagination
      if (navDp) {
        const explicit = Number(navDp.datapointTotalPages)
        if (explicit > 0) return explicit
      }
      if (this.isNavDatapointsSource) {
        const fromDiy = Number(readDiyValue(this.detail, 'navDatapointTotalPages'))
        if (fromDiy > 0) return fromDiy
      }
      const nav = this.navListPagination
      if (nav) {
        const explicit = Number(nav.totalPages)
        if (explicit > 0) return explicit
      }
      return 1
    },
    pagerCanPrev() {
      return this.pagerCurrent > 1
    },
    pagerCanNext() {
      return this.pagerCurrent < this.pagerTotalPages
    },
    pagerInfoText() {
      const cur = this.pagerCurrent
      const total = this.pagerTotalPages
      const n = this.paginationTotal || 0
      if (this.isNavDatapointsSource || this.navDatapointPagination) {
        return `第 ${cur}/${total} 页 · 共 ${n} 个测点`
      }
      if (this.navListPagination) {
        return `第 ${cur}/${total} 页 · 共 ${n} 台设备`
      }
      return `第 ${cur}/${total} 页 · 共 ${n} 条`
    },
    navListPagination() {
      const nav = getEditorNavContext()
      return nav && nav.deviceListMode ? nav : null
    },
    isDeviceCardMode() {
      // 页面配置先于 navContext 到达；优先按 rowSource 决定卡片模式，
      // 避免切页首帧短暂挂载 Ant Table 后再切换成设备卡片。
      return this.isNavChildrenSource || !!this.navListPagination
    },
    isDatapointCardMode() {
      // 设备点位页切换时，navContext 会先到、rowSource 稍后才注入。
      // 直接识别 signal 路由可避免中间态短暂渲染旧 Ant Table。
      const nav = getEditorNavContext()
      return this.isNavDatapointsSource || !!(nav && (nav.signalMode || nav.routeMode === 'signal'))
    },
    isRuntimeCardMode() {
      // 模板编辑态与运行态共用同一卡片矩阵，保证设备/点位列表所见即所得。
      // 工具箱缩略预览仍保留轻量模式，避免拖拽源节点创建完整卡片网格。
      return !this.IsToolBox
        && (this.isDeviceCardMode || this.isDatapointCardMode)
    },
    allCardDevices() {
      const nav = this.navListPagination
      return nav ? (nav.allChildDevices || nav.allChildNodes || nav.childDevices || nav.childNodes || []) : []
    },
    cardDeviceTotal() {
      return this.allCardDevices.length
    },
    cardPageSize() {
      return Math.max(1, Number((this.navListPagination || {}).pageSize) || 49)
    },
    cardCurrentIndex() {
      const totalPages = Math.max(1, Math.ceil(this.cardDeviceTotal / this.cardPageSize))
      const navIndex = Number((this.navListPagination || {}).pageIndex) || 0
      return Math.min(navIndex, totalPages - 1)
    },
    cardDevices() {
      const start = this.cardCurrentIndex * this.cardPageSize
      return this.allCardDevices.slice(start, start + this.cardPageSize)
    },
    runtimeCards() {
      if (this.isDeviceCardMode) {
        // 物理模板编辑态没有 navContext，使用模板内的设备行生成同版式卡片预览。
        // 运行态存在真实设备树时仍以 cardDevices 为唯一数据源。
        const devices = this.cardDevices.length
          ? this.cardDevices
          : this.dynamicData.slice(0, this.cardPageSize).map(row => ({
            name: row.rowName || '',
            label: row.rowName || '',
            code: row.deviceCode || '',
            status: 'off',
          }))
        return devices.map((device, index) => {
          const name = device.name || device.label || ''
          const isCabinet = device.kind === 'virtualCabinet' || !!device.virtualCabinet
          const cabinet = device.virtualCabinet || (isCabinet ? name : '')
          return {
            key: isCabinet
              ? `vc-${device.parentDeviceUuid || device.uuid || ''}-${cabinet}-${index}`
              : (device.sid || device.uuid || `${name || 'device'}-${index}`),
            name,
            displayName: name,
            online: this.resolveDeviceOnline(device, index),
            title: name,
            source: device,
          }
        })
      }
      const nav = getEditorNavContext() || {}
      const deviceName = nav.name || nav.label || ''
      return this.dynamicData.map((row, index) => {
        const name = row.rowName || ''
        const pointUuid = (this.rowPointUuids && this.rowPointUuids[index]) || ''
        return {
          key: pointUuid
            || (this.bindingMatrix[index] && this.bindingMatrix[index][0])
            || `${name || 'point'}-${index}`,
          name,
          displayName: stripDevicePrefix(name, deviceName),
          value: row.col_0,
          unit: row.deviceCode || '',
          title: name,
          pointUuid,
        }
      })
    },
    navDatapointPagination() {
      if (!this.isNavDatapointsSource) return null
      const nav = getEditorNavContext()
      return nav && (nav.signalMode || nav.routeMode === 'signal') ? nav : null
    },
    isNavChildrenSource() {
      return String(readDiyValue(this.detail, 'rowSource') || '') === 'navChildren'
    },
    isNavDatapointsSource() {
      return String(readDiyValue(this.detail, 'rowSource') || '') === 'navDatapoints'
    },
    dynamicColumns() {
      const indexCol = {
        title: '序号',
        key: 'index',
        width: 60,
        align: 'center',
        className: 'sticky-col sticky-col-0',
        customHeaderCell: () => ({ class: 'sticky-col sticky-col-0' }),
        customRender: (t, r, index) => {
          const navDp = this.navDatapointPagination
          if (navDp) {
            return (navDp.datapointPageIndex || 0) * (navDp.datapointPageSize || REAL_DATA_DEFAULT_PAGE_SIZE) + index + 1
          }
          const nav = this.navListPagination
          if (nav) {
            return (nav.pageIndex || 0) * (nav.pageSize || REAL_DATA_DEFAULT_PAGE_SIZE) + index + 1
          }
          const { current, pageSize } = this.pagination || { current: 1, pageSize: REAL_DATA_DEFAULT_PAGE_SIZE }
          return (current - 1) * pageSize + index + 1
        }
      }

      if (this.isNavDatapointsSource) {
        const columns = [
          indexCol,
          {
            title: '测点名称',
            key: 'pointName',
            width: 220,
            align: 'left',
            className: 'sticky-col sticky-col-1',
            customHeaderCell: () => ({ class: 'sticky-col sticky-col-1' }),
            customRender: (text, record) => record.rowName || '-'
          }
        ]
        this.columnHeaders.forEach((colName, colIndex) => {
          columns.push({
            title: colName,
            dataIndex: `col_${colIndex}`,
            key: `col_${colIndex}`,
            width: 120,
            align: 'center'
          })
        })
        columns.push({
          title: '单位',
          key: 'unit',
          width: 80,
          align: 'center',
          customRender: (text, record) => record.deviceCode || '-'
        })
        return columns
      }

      const columns = [
        indexCol,
        {
          title: '设备名称',
          key: 'alias',
          width: 120,
          align: 'center',
          className: 'sticky-col sticky-col-1',
          customHeaderCell: () => ({ class: 'sticky-col sticky-col-1' }),
          customRender: (text, record) => record.rowName || '-'
        },
        {
          title: '设备编号',
          key: 'realName',
          width: 120,
          align: 'center',
          className: 'sticky-col sticky-col-2',
          customHeaderCell: () => ({ class: 'sticky-col sticky-col-2' }),
          customRender: (text, record) => record.deviceCode || '-'
        }
      ]
      this.columnHeaders.forEach((colName, colIndex) => {
        columns.push({
          title: colName,
          dataIndex: `col_${colIndex}`,
          key: `col_${colIndex}`,
          width: 100,
          align: 'center'
        })
      })
      return columns
    }
  },
  data() {
    return {
      pagination: {
        pageSize: REAL_DATA_DEFAULT_PAGE_SIZE,
        showSizeChanger: false,
        hideOnSinglePage: false,
        showLessItems: true,
        simple: true,
        current: 1,
        total: 0,
        align: 'center'
      },
      tableSplitColor: '#000',
      tableHoverColor: '#fff',
      tableHeaderColor: '',
      tableHeaderBackColor: '',
      tableHeaderFont: 'Arial',
      scrollBgColor: '#f0f0f0',
      scrollFrColor: '#c1c1c1',
      scrollHdColor: '#a8a8a8',
      tableHeaderFontSize: '17px',
      selectedTheme: 'light',
      scrollbarStyleTagId: '',
      rowDeviceNames: [],
      rowDeviceCodes: [],
      columnHeaders: [],
      bindingMatrix: [],
      // 与当前页行一一对应的点位实例 uuid（device_real_data.uuid），刷新时按此匹配，禁止按下标串值
      rowPointUuids: [],
      cellData: [],
      detail: {
        identifier: '',
        style: {
          position: { x: 0, y: 0, w: 520, h: 200 },
          diy: [],
          visible: 1,
          animate: '',
        },
        animate: { selected: [], animateElement: [], isExpression: false },
      },
      IsToolBox: false,
      editMode: true,
      width: 600,
      height: 600,
      strokeColor: '#000000',
      fill: '#A1BFE2',
      strokeWidth: 0.3,
      fillOpacity: 1,
      strokeOpacity: 1,
      animateType: 'blink',
      startColor: '#74f808',
      stopColor: '#74f808',
      animateSpeed: 0.5,
      animateSpinSpeed: 0.5,
      spinDirection: 0,
      blinkSpeed: 0.5,
      isStart: false,
      AlarmTimer: null,
      autoPageTimer: null,
      autoPageResumeTimer: null,
      autoPagePendingUntil: 0,
      navRefreshRequestId: 0,
      // 设备卡在线态：优先 device.DeviceStatus 实时点，键为 uuid/sid/name
      deviceOnlineMap: Object.create(null),
      deviceOnlineTimer: null,
      deviceOnlineReqSeq: 0,
      waitTime: 8000,
      _navHydrateFingerprint: '',
      _navHydrateAt: 0,
      _navHydrateScheduled: false,
      _applyingSignalPage: false,
      _lastNavPageFingerprint: '',
      fontFamily: 'Arial',
      fontSize: '14',
      backColor: '',
      foreColor: '',
      base: {
        text: 'configComponent.viewRealTable.title',
        icon: 'icon-biaoge3',
        isFontIcon: true,
        info: {
          type: 'image',
          action: [],
          dataBind: [],
          animate: {
            selected: [],
            condition: {
              deviceSN: '',
              selectVideoType: 0,
              isBandDevice: false,
              bandType: 1,
              dataID: '',
              dataName: '',
              operator: '',
              OperatorValue: '',
              OperatorMaxValue: ''
            },
            isExpression: false,
            animateList: [
              { id: 'blink', name: 'component.public.animateBlink' },
              { id: 'Zoom', name: 'component.public.Zoom' },
              { id: 'animateSpin', name: 'component.public.animateSpin' }
            ],
            animateElement: [
              {
                id: 'blink',
                elementList: [
                  {
                    name: 'component.public.animateSpeed',
                    type: 7,
                    value: 1,
                    min: 0.1,
                    key: 'blinkSpeed'
                  }
                ]
              },
              {
                id: 'millcolorGrad',
                elementList: [
                  { name: 'component.public.startColor', type: 2, value: '#74f808', key: 'startColor' },
                  { name: 'component.public.stopColor', type: 2, value: '#f30b0b', key: 'stopColor' },
                  { name: 'component.public.animateSpeed', type: 7, value: 1, min: 0.1, key: 'animateSpeed' }
                ]
              },
              {
                id: 'animateSpin',
                elementList: [
                  { name: 'component.public.animateSpinSpeed', type: 7, value: 1, min: 0.1, key: 'spinSpeed' },
                  {
                    name: 'configComponent.bigScreen.border.border89Direction',
                    type: 6,
                    value: 0,
                    enumList: [
                      { value: 0, option: 'configComponent.bigScreen.border.border89DirectionForward' },
                      { value: 1, option: 'configComponent.bigScreen.border.border89DirectionNegative' }
                    ],
                    min: 1,
                    key: 'spinDirection'
                  }
                ]
              }
            ]
          },
          style: {
            position: {
              x: 0,
              y: 0,
              w: 520,
              h: 200
            },
            backColor: 'transparent',
            foreColor: '#000000',
            fontWeight: 400,
            fontSize: 15,
            fontFamily: 'Arial',
            visible: 1,
            zIndex: -1,
            transform: 0,
            diy: [
              {
                name: 'configComponent.viewRealTable.columnHeaders',
                type: 9,
                value: '运行状态, 功率, 温度, 在线时长',
                key: 'columnHeaders'
              },
              {
                name: 'configComponent.viewRealTable.rowDeviceNames',
                type: 9,
                value: '空调\n灯光\n门禁\n监控',
                key: 'rowDeviceNames'
              },
              {
                name: 'configComponent.viewRealTable.rowDeviceCodes',
                type: 9,
                value: 'AC001\nLT002\nAC003\nCAM004',
                key: 'rowDeviceCodes'
              },
              {
                name: 'configComponent.viewRealTable.rowBindings',
                type: 9,
                value: 'AC001->status,AC001->power,AC001->temp,AC001->online_time;LT002->status,LT002->power,LT002->temp,LT002->online_time;AC003->status,AC003->power,AC003->temp,AC003->online_time;CAM004->status,CAM004->power,CAM004->temp,CAM004->online_time',
                key: 'rowBindings'
              },
              {
                name: 'configComponent.AlarmList.waitTime',
                type: 7,
                value: 1000,
                min: 100,
                max: 10000,
                key: 'waitTime'
              },
              {
                name: 'configComponent.DeviceTree.ShowCount',
                type: 1,
                value: 5,
                min: 1,
                max: 100,
                key: 'ShowCount'
              },
              {
                name: '自动翻页（已下线）',
                type: 6,
                value: 0,
                enumList: [
                  { value: 0, option: '关闭' },
                  { value: 1, option: '开启（暂不可用）' }
                ],
                key: 'AutoPageEnabled'
              },
              {
                name: '自动翻页间隔（秒）',
                type: 1,
                value: 5,
                min: 1,
                max: 3600,
                key: 'AutoPageInterval'
              },
              {
                name: '手动翻页暂停（秒）',
                type: 1,
                value: 60,
                min: 1,
                max: 3600,
                key: 'AutoPageResumeDelay'
              },
              {
                name: '卡片外框强调色',
                type: 2,
                value: '#4ae6ff',
                key: 'panelAccentColor'
              },
              {
                name: '设备图标流光色',
                type: 2,
                value: '#52e8ff',
                key: 'deviceIconAccent'
              },
              {
                name: '点位图标流光色',
                type: 2,
                value: '#bca5ff',
                key: 'pointIconAccent'
              },
              {
                name: 'configComponent.DeviceTree.SearchColor',
                type: 2,
                value: '#000000',
                key: 'SearchColor'
              },
              {
                name: 'configComponent.DeviceTree.SearchBackColor',
                type: 2,
                value: '#ffffff',
                key: 'SearchBackColor'
              },
              {
                name: 'configComponent.DeviceTree.SearchBorderColor',
                type: 2,
                value: '#cbc6c6',
                key: 'SearchBorderColor'
              },
              {
                name: 'configComponent.DataHistoryList.tableHeaderColor',
                type: 2,
                value: '#000000',
                key: 'tableHeaderColor'
              },
              {
                name: 'configComponent.DataHistoryList.tableHeaderBackColor',
                type: 2,
                value: '#fafafa',
                key: 'tableHeaderBackColor'
              },
              {
                name: 'configComponent.viewRealTable.tableHeaderFont',
                type: 3,
                value: 'Arial',
                key: 'tableHeaderFont'
              },
              {
                name: 'configComponent.viewRealTable.tableHeaderFontSize',
                type: 1,
                value: 14,
                key: 'tableHeaderFontSize'
              },
              {
                name: 'configComponent.DataHistoryList.tableSplitColor',
                type: 2,
                value: '#ebedf0',
                key: 'tableSplitColor'
              },
              {
                name: 'configComponent.DataHistoryList.tableHoverColor',
                type: 2,
                value: '#ffffff',
                key: 'tableHoverColor'
              },
              {
                name: '主题风格',
                type: 6,
                value: 'light',
                enumList: [
                  { value: 'light', option: '极简亮色' },
                  { value: 'dark', option: '深空夜幕' },
                  { value: 'ocean', option: '海岸蓝调' },
                  { value: 'amber', option: '琥珀暖光' },
                  { value: 'emerald', option: '森林翠影' }
                ],
                key: 'themeName'
              }
            ]
          }
        }
      }
    };
  },
  methods: {
    openDeviceDatapoints(device) {
      if (!device) return
      const currentNav = getEditorNavContext()
      let payload = { ...device }
      // 虚拟列头柜列表页：卡片可能只剩 name，从列表 nav 补齐父设备与前缀
      if (currentNav && currentNav.virtualCabinetListMode) {
        try {
          const { normalizeVirtualCabinetClick } = require('@/pages/ISMDisPlay/utils/virtualCabinet')
          const normalized = normalizeVirtualCabinetClick(device, currentNav)
          if (normalized) payload = normalized
        } catch (e) {
          payload = {
            ...device,
            kind: 'virtualCabinet',
            virtualCabinet: device.virtualCabinet || device.name || device.label || '',
            parentDeviceUuid: device.parentDeviceUuid || currentNav.deviceUuid || currentNav.uuid || '',
            parentDeviceLabel: device.parentDeviceLabel || currentNav.name || currentNav.label || '',
            uuid: device.parentDeviceUuid || currentNav.deviceUuid || currentNav.uuid || device.uuid || '',
            deviceUuid: device.parentDeviceUuid || currentNav.deviceUuid || currentNav.uuid || device.deviceUuid || '',
            modelUuid: device.modelUuid || device.muid || currentNav.modelUuid || currentNav.muid || '',
            muid: device.muid || device.modelUuid || currentNav.muid || currentNav.modelUuid || '',
          }
        }
      }
      this.$EventBus.$emit('OpenDeviceDatapoints', {
        ...payload,
        deviceListReturnContext: currentNav && (currentNav.deviceListMode || currentNav.virtualCabinetListMode)
          ? { ...currentNav }
          : null,
      })
    },
    shortDeviceCode(code) {
      const value = String(code || '')
      return value.length > 14 ? `${value.slice(0, 8)}…${value.slice(-5)}` : value
    },
    deviceOnlineKey(device, index) {
      if (!device) return `idx-${index}`
      return String(device.uuid || device.sid || device.name || device.label || `idx-${index}`)
    },
    resolveDeviceOnline(device, index) {
      const key = this.deviceOnlineKey(device, index)
      if (Object.prototype.hasOwnProperty.call(this.deviceOnlineMap, key)) {
        return !!this.deviceOnlineMap[key]
      }
      // 尚未拉到状态点时回退设备树 Status（on/off）
      return !!(device && device.status === 'on')
    },
    refreshDeviceOnlineStatus() {
      if (!this.isDeviceCardMode || this.editMode || this.IsToolBox) return
      const devices = this.cardDevices
      if (!devices.length) {
        this.deviceOnlineMap = Object.create(null)
        return
      }
      const bindings = devices.map((d) => {
        const name = String((d && (d.name || d.label)) || '').trim()
        // 系统内置在线状态点；无设备名时跳过该行
        return name ? [`${name}->device.DeviceStatus`] : ['']
      })
      const seq = ++this.deviceOnlineReqSeq
      postRealDataByBindings(bindings)
        .then((res) => {
          if (seq !== this.deviceOnlineReqSeq || this._isBeingDestroyed || this._isDestroyed) return
          const rows = res && res.data && res.data.code === 0 && Array.isArray(res.data.realData)
            ? res.data.realData
            : null
          if (!rows) return
          const next = Object.create(null)
          for (let i = 0; i < devices.length; i++) {
            const device = devices[i]
            const key = this.deviceOnlineKey(device, i)
            const raw = Array.isArray(rows[i]) ? rows[i][0] : rows[i]
            const fromPoint = isDeviceOnlineFromStatusValue(raw)
            next[key] = fromPoint == null ? !!(device && device.status === 'on') : fromPoint
          }
          this.deviceOnlineMap = next
        })
        .catch(() => { /* 保持上一轮/树状态，避免闪全绿 */ })
    },
    startDeviceOnlinePolling() {
      this.stopDeviceOnlinePolling()
      if (!this.isDeviceCardMode || this.editMode || this.IsToolBox) return
      this.refreshDeviceOnlineStatus()
      this.deviceOnlineTimer = setInterval(() => {
        this.refreshDeviceOnlineStatus()
      }, DEVICE_ONLINE_POLL_MS)
    },
    stopDeviceOnlinePolling() {
      if (this.deviceOnlineTimer) {
        clearInterval(this.deviceOnlineTimer)
        this.deviceOnlineTimer = null
      }
    },
    applyScrollbarTheme() {
      if (typeof document === 'undefined') {
        return;
      }
      const theme = this.currentTheme;
      const className = this.scrollbarThemeClass;
      const styleId = `scrollbar-style-${className}`;
      this.scrollbarStyleTagId = styleId;
      let styleTag = document.getElementById(styleId);
      if (!styleTag) {
        styleTag = document.createElement('style');
        styleTag.id = styleId;
        document.head.appendChild(styleTag);
      }
      styleTag.textContent = `
        .${className} ::-webkit-scrollbar {
          width: 10px !important;
          height: 10px !important;
        }
        .${className} ::-webkit-scrollbar-track {
          background: ${theme.scrollBgColor} !important;
          border-radius: 999px !important;
        }
        .${className} ::-webkit-scrollbar-thumb {
          background: ${theme.scrollFrColor} !important;
          border-radius: 999px !important;
          border: 2px solid ${theme.scrollBgColor} !important;
          background-image: none !important;
        }
        .${className} ::-webkit-scrollbar-thumb:hover {
          background: ${theme.scrollHdColor} !important;
        }
      `;
    },
    isEmptyCellValue(value) {
      return value === undefined || value === null || value === '' || value === '-' || value === '—'
    },
    /** 点位 uuid 统一小写，避免 HTTP/WS 大小写不一致导致匹配失败后掉进短名串值 */
    normalizePointUuid(uuid) {
      return String(uuid || '').trim().toLowerCase()
    },
    /** 点位身份键：优先实例 uuid，其次绑点串，禁止仅用行号 */
    pointRowKey(uuid, binding, fallbackName, index) {
      const id = this.normalizePointUuid(uuid)
      if (id) return `u:${id}`
      const bind = String(binding || '').trim()
      if (bind) return `b:${bind}`
      const name = String(fallbackName || '').trim()
      if (name) return `n:${name}`
      return `i:${index}`
    },
    currentRowKeys() {
      const n = Math.max(
        (this.rowPointUuids && this.rowPointUuids.length) || 0,
        (this.bindingMatrix && this.bindingMatrix.length) || 0,
        (this.rowDeviceNames && this.rowDeviceNames.length) || 0,
        (this.cellData && this.cellData.length) || 0,
      )
      const keys = []
      for (let i = 0; i < n; i += 1) {
        keys.push(this.pointRowKey(
          this.rowPointUuids && this.rowPointUuids[i],
          this.bindingMatrix && this.bindingMatrix[i] && this.bindingMatrix[i][0],
          this.rowDeviceNames && this.rowDeviceNames[i],
          i,
        ))
      }
      return keys
    },
    /**
     * 空值不覆盖已有有效读数，但必须按点位身份键对齐。
     * 禁止「同下标保留上一行值」——翻页/过滤/重排时会把 A 点的值显示到 B 点卡片上。
     * @param {Array} prevRows
     * @param {Array} nextRows
     * @param {string[]} [nextKeys]
     * @param {string[]} [prevKeys] 须在改写 rowPointUuids 之前传入，否则会串键
     */
    mergeCellDataPreserve(prevRows, nextRows, nextKeys, prevKeys) {
      const prev = Array.isArray(prevRows) ? prevRows : []
      const next = Array.isArray(nextRows) ? nextRows : []
      const keys = Array.isArray(nextKeys) ? nextKeys : null
      const oldKeys = Array.isArray(prevKeys) ? prevKeys : this.currentRowKeys()
      const prevByKey = Object.create(null)
      for (let i = 0; i < oldKeys.length; i += 1) {
        const k = oldKeys[i]
        if (k && Array.isArray(prev[i])) prevByKey[k] = prev[i]
      }
      return next.map((row, rowIndex) => {
        const nextRow = Array.isArray(row) ? row : [row]
        const key = keys ? keys[rowIndex] : null
        // 有身份键时只取同点位旧值；无键时才退回同下标（非测点表兼容）
        let prevRow = []
        if (key && prevByKey[key]) {
          prevRow = prevByKey[key]
        } else if (!keys || !keys.length) {
          prevRow = Array.isArray(prev[rowIndex]) ? prev[rowIndex] : []
        }
        return nextRow.map((cell, colIndex) => {
          if (!this.isEmptyCellValue(cell)) return cell
          const kept = prevRow[colIndex]
          return this.isEmptyCellValue(kept) ? (cell === '—' ? '—' : '-') : kept
        })
      })
    },
    batchUpdateConfig(newCellData, nextKeys) {
      const keys = Array.isArray(nextKeys) ? nextKeys : this.currentRowKeys()
      const merged = this.mergeCellDataPreserve(this.cellData, newCellData, keys, keys)
      this.cellData = merged;
      if (this.pagination) {
        const total = newCellData ? newCellData.length : 0;
        const current = Math.min(this.pagination.current, Math.max(1, Math.ceil(total / this.pagination.pageSize)) || 1);
        this.pagination = {
          ...this.pagination,
          current,
          total
        };
      }
    },
    clearAutoPaging() {
      clearInterval(this.autoPageTimer)
      clearTimeout(this.autoPageResumeTimer)
      this.autoPageTimer = null
      this.autoPageResumeTimer = null
    },
    startAutoPaging() {
      clearInterval(this.autoPageTimer)
      this.autoPageTimer = null
      if (this.editMode || this.IsToolBox || !this.isRuntimeCardMode || !this.autoPageEnabled) return
      this.autoPageTimer = setInterval(() => {
        if (this.pagerTotalPages <= 1 || Date.now() < this.autoPagePendingUntil) return
        const nextPage = this.pagerCurrent >= this.pagerTotalPages ? 1 : this.pagerCurrent + 1
        this.autoPagePendingUntil = Date.now() + Math.min(this.autoPageIntervalMs, 4000)
        this.handlePageChange(nextPage, { automatic: true })
      }, this.autoPageIntervalMs)
    },
    pauseAutoPaging() {
      if (this.editMode || this.IsToolBox || !this.autoPageEnabled) return
      clearInterval(this.autoPageTimer)
      clearTimeout(this.autoPageResumeTimer)
      this.autoPageTimer = null
      this.autoPageResumeTimer = setTimeout(() => {
        this.autoPageResumeTimer = null
        this.startAutoPaging()
      }, this.autoPageResumeDelayMs)
    },
    handlePageChange(page, options = {}) {
      const nextPage = Math.max(1, Number(page) || 1)
      const totalPages = this.pagerTotalPages
      if (nextPage < 1 || nextPage > totalPages) return
      if (!options.automatic) this.pauseAutoPaging()
      const eventMeta = options.automatic ? { autoPage: true } : {}

      if (this.isDeviceCardMode) {
        const next = nextPage - 1
        if (next === this.cardCurrentIndex) return
        this.$EventBus.$emit('NavPageChange', { pageIndex: next, ...eventMeta })
        return
      }
      if (this.isNavDatapointsSource) {
        const next = nextPage - 1
        const nav = getEditorNavContext()
        const cur = (nav && (nav.signalMode || nav.routeMode === 'signal'))
          ? (nav.datapointPageIndex || 0)
          : Number(readDiyValue(this.detail, 'navDatapointPageIndex') || 0)
        if (next === cur) return
        // 确保 store 有 signalMode，否则 ISMRender 会忽略翻页事件
        if (!nav || !(nav.signalMode || nav.routeMode === 'signal')) {
          commitEditorNav({
            signalMode: true,
            routeMode: 'signal',
            datapointPageIndex: next,
          })
        }
        this.$EventBus.$emit('NavPageChange', { datapointPageIndex: next, ...eventMeta })
        return
      }
      if (this.isNavChildrenSource || this.navListPagination) {
        const next = nextPage - 1
        const nav = this.navListPagination || getEditorNavContext()
        if (next === ((nav && nav.pageIndex) || 0)) return
        this.$EventBus.$emit('NavPageChange', { pageIndex: next, ...eventMeta })
        return
      }
      this.pagination = {
        ...this.pagination,
        current: nextPage
      };
    },
    goPagerPrev(e) {
      if (e) {
        e.preventDefault()
        e.stopPropagation()
      }
      if (!this.pagerCanPrev) return
      // 防抖：pointer 合成事件偶发连点
      const now = Date.now()
      if (this._pagerGuardAt && now - this._pagerGuardAt < 280) return
      this._pagerGuardAt = now
      this.handlePageChange(this.pagerCurrent - 1)
    },
    goPagerNext(e) {
      if (e) {
        e.preventDefault()
        e.stopPropagation()
      }
      if (!this.pagerCanNext) return
      const now = Date.now()
      if (this._pagerGuardAt && now - this._pagerGuardAt < 280) return
      this._pagerGuardAt = now
      this.handlePageChange(this.pagerCurrent + 1)
    },
    setDiyOnDetail(key, value) {
      if (!this.detail || !this.detail.style) return
      if (!Array.isArray(this.detail.style.diy)) {
        this.$set(this.detail.style, 'diy', [])
      }
      const diy = this.detail.style.diy
      const item = diy.find(d => d && d.key === key)
      if (item) {
        item.value = value
      } else {
        diy.push({ name: key, type: 9, value, key })
      }
    },
    /** 测点页指纹：用于打断 NeedHydrate ↔ PageUpdate 同步重入 */
    navPageFingerprint(nav) {
      if (!nav) return ''
      const deviceUuid = nav.deviceUuid || nav.uuid || ''
      const idx = Number(nav.datapointPageIndex) || 0
      const size = Number(nav.datapointPageSize) || 0
      const total = Number(nav.totalDatapoints) || 0
      const points = Array.isArray(nav.datapoints) ? nav.datapoints : []
      const uuids = points.slice(0, 40).map(p => String((p && (p.uuid || p.Uuid)) || '').trim()).join(',')
      return [deviceUuid, idx, size, total, points.length, uuids, nav.serverPaged ? 1 : 0].join('|')
    },
    /** 同指纹短时间内只 hydrate 一次；且必须异步发出，打断同步事件环 */
    requestNavHydrate(nav) {
      if (!this.$EventBus) return
      const fp = this.navPageFingerprint(nav) || 'empty'
      const now = Date.now()
      if (this._navHydrateFingerprint === fp && now - this._navHydrateAt < 3000) {
        return
      }
      if (this._navHydrateScheduled) return
      this._navHydrateFingerprint = fp
      this._navHydrateAt = now
      this._navHydrateScheduled = true
      const emit = () => {
        this._navHydrateScheduled = false
        if (!this.$EventBus || this._isDestroyed) return
        this.$EventBus.$emit('NavDatapointNeedHydrate')
      }
      // 强制异步：禁止在 PageUpdate/applySignalPage 同步栈里再进 NeedHydrate（现场栈溢出根因）
      this.$nextTick(() => {
        setTimeout(emit, 0)
      })
    },
    /** 页内翻页：只换当前页测点行并拉实时值，不整页重载 */
    applySignalPageFromNav(nav) {
      if (!nav || this.editMode || this.IsToolBox) return
      if (!this.isNavDatapointsSource && String(readDiyValue(this.detail, 'rowSource') || '') !== 'navDatapoints') {
        return
      }
      if (this._applyingSignalPage) return
      this._applyingSignalPage = true
      try {
        this._applySignalPageFromNavInner(nav)
      } finally {
        this._applyingSignalPage = false
      }
    },
    _applySignalPageFromNavInner(nav) {
      const serverPaged = !!nav.serverPaged
      const declaredTotal = Number(nav.totalDatapoints) || 0
      const size = Math.max(1, Number(nav.datapointPageSize) || 20)
      // 服务端分页：datapoints 已是当前页，禁止再按 pageIndex 本地切片，
      // 也禁止 fetchAll 补灌（后端硬上限 5000 → 80 条/页时总页数锁死在 63）。
      let pagePoints
      let total
      let idx
      if (serverPaged) {
        pagePoints = Array.isArray(nav.datapoints) && nav.datapoints.length
          ? nav.datapoints
          : (Array.isArray(nav.childNodes) ? nav.childNodes : [])
        if (!pagePoints.length) {
          // 空页只异步请求一次；禁止在 PageUpdate 同步栈里再炸一次
          this.requestNavHydrate(nav)
          return
        }
        total = Math.max(declaredTotal, pagePoints.length)
        const totalPages = Math.max(1, Math.ceil(total / size) || 1)
        idx = Math.max(0, Math.min(Number(nav.datapointPageIndex) || 0, totalPages - 1))
        this._applySignalPageRows(nav, pagePoints, size, total, totalPages, idx, true)
        return
      }
      // 必须用全量 allDatapoints；datapoints 可能已是当前页切片
      let all = Array.isArray(nav.allDatapoints) && nav.allDatapoints.length
        ? nav.allDatapoints
        : (nav.datapoints || nav.childNodes || [])
      // store 总数大于当前 all 长度时，说明 all 被截成页切片，需回源补齐
      if (declaredTotal > all.length) {
        this.requestNavHydrate(nav)
        // 仍用现有行先渲染，避免空白
      }
      if (!all.length) {
        this.requestNavHydrate(nav)
        return
      }
      total = Math.max(all.length, declaredTotal)
      const totalPages = Math.max(1, Math.ceil(total / size) || 1)
      idx = Math.max(0, Math.min(Number(nav.datapointPageIndex) || 0, totalPages - 1))
      pagePoints = all.slice(idx * size, idx * size + size)
      this._applySignalPageRows(nav, pagePoints, size, total, totalPages, idx, false, all)
    },
    _applySignalPageRows(nav, pagePoints, size, total, totalPages, idx, serverPaged, all) {
      const deviceLabel = nav.name || nav.label || ''
      const prefix = deviceLabel ? `${deviceLabel}_` : ''
      const rowDeviceNames = pagePoints.map(p => {
        const n = String(p.name || p.label || '')
        if (prefix && n.startsWith(prefix)) return n.slice(prefix.length)
        return n
      })
      const rowDeviceCodes = pagePoints.map(p => p.unit || '')
      const rowPointUuids = pagePoints.map(p => this.normalizePointUuid(p.uuid || p.Uuid))
      // 绑点：优先网关设备名->测点名（与内存 Map key 一致）
      const bindingMatrix = pagePoints
        .map(p => {
          const pointName = p.name || p.label || ''
          if (!pointName) return null
          const owner = p.deviceName || p.device_name || ''
          if (owner) return [`${owner}->${pointName}`]
          return [pointName]
        })
        .filter(Boolean)

      const prevKeys = this.currentRowKeys()
      const prevCellData = this.cellData
      const nextKeys = pagePoints.map((p, i) => this.pointRowKey(
        rowPointUuids[i],
        bindingMatrix[i] && bindingMatrix[i][0],
        rowDeviceNames[i],
        i,
      ))

      this.applyTableConfig({
        rowDeviceNames,
        rowDeviceCodes,
        columnHeaders: ['实时值'],
        bindingMatrix,
        rowPointUuids,
      })
      this.setDiyOnDetail('rowSource', 'navDatapoints')
      this.setDiyOnDetail('columnHeaders', '实时值')
      this.setDiyOnDetail('rowDeviceNames', rowDeviceNames.join('\n'))
      this.setDiyOnDetail('rowDeviceCodes', rowDeviceCodes.join('\n'))
      this.setDiyOnDetail('rowBindings', bindingMatrix.map(r => r[0]).join(';'))
      this.setDiyOnDetail('navTotalDatapoints', String(total))
      this.setDiyOnDetail('navDatapointPageIndex', String(idx))
      this.setDiyOnDetail('navDatapointPageSize', String(size))
      this.setDiyOnDetail('navDatapointTotalPages', String(totalPages))
      // 回写 store，避免 datapointTotalPages 陈旧为 1
      const navPatch = {
        signalMode: true,
        routeMode: 'signal',
        datapointPageIndex: idx,
        datapointPageSize: size,
        totalDatapoints: total,
        datapointTotalPages: totalPages,
        datapoints: pagePoints,
      }
      if (serverPaged) {
        navPatch.serverPaged = true
        // 服务端分页禁止把当前页误写入 allDatapoints，否则后续会按本地全量切片卡死
        navPatch.allDatapoints = []
      } else {
        navPatch.allDatapoints = all
      }
      commitEditorNav(navPatch)
      // 先用 GetRealData 已覆盖的内存值填表（与设备管理同源），再异步刷新。
      // 空值仅保留「同一点位 uuid」的上一帧，绝不按下标串值。
      this.cellData = this.mergeCellDataPreserve(
        prevCellData,
        pagePoints.map(p => {
          const v = p.value
          if (v === undefined || v === null || v === '') return ['-']
          return [v]
        }),
        nextKeys,
        prevKeys,
      )
      this._navPagePointUuids = rowPointUuids.filter(Boolean)
      this._lastNavPageFingerprint = this.navPageFingerprint({
        ...nav,
        datapointPageIndex: idx,
        datapointPageSize: size,
        totalDatapoints: total,
        datapoints: pagePoints,
        serverPaged: !!serverPaged,
      })
      this.$nextTick(() => {
        this.QueryRealData()
      })
    },
    onNavDatapointPageUpdate(nav) {
      // 同一页仅总数/页码元数据变更时，不要整表重灌（会把值闪成 —）
      if (nav && (this.isNavDatapointsSource || String(readDiyValue(this.detail, 'rowSource') || '') === 'navDatapoints')) {
        if (nav._metaOnly) {
          if (nav.totalDatapoints) {
            this.setDiyOnDetail('navTotalDatapoints', String(nav.totalDatapoints))
            const size = Math.max(1, Number(nav.datapointPageSize) || this.pagerPageSize)
            this.setDiyOnDetail('navDatapointTotalPages', String(Math.max(1, Math.ceil(Number(nav.totalDatapoints) / size))))
          }
          return
        }
        const fp = this.navPageFingerprint(nav)
        if (fp && fp === this._lastNavPageFingerprint && !nav._forcePageReload) {
          if (nav.totalDatapoints) {
            this.setDiyOnDetail('navTotalDatapoints', String(nav.totalDatapoints))
            const size = Math.max(1, Number(nav.datapointPageSize) || this.pagerPageSize)
            this.setDiyOnDetail('navDatapointTotalPages', String(Math.max(1, Math.ceil(Number(nav.totalDatapoints) / size))))
          }
          return
        }
        const nextIdx = Number(nav.datapointPageIndex) || 0
        const curIdx = Number(readDiyValue(this.detail, 'navDatapointPageIndex') || 0)
        const hasRows = (this.rowDeviceNames && this.rowDeviceNames.length)
          || (this.bindingMatrix && this.bindingMatrix.length)
        const incoming = Array.isArray(nav.datapoints) ? nav.datapoints : []
        let namesChanged = false
        if (incoming.length && this.rowDeviceNames && this.rowDeviceNames.length) {
          const deviceLabel = nav.name || nav.label || ''
          const prefix = deviceLabel ? `${deviceLabel}_` : ''
          const n = Math.min(incoming.length, this.rowDeviceNames.length)
          for (let i = 0; i < n; i += 1) {
            const incomingUuid = String(incoming[i].uuid || incoming[i].Uuid || '').trim()
            const rowUuid = String((this.rowPointUuids && this.rowPointUuids[i]) || '').trim()
            if (incomingUuid && rowUuid && incomingUuid !== rowUuid) {
              namesChanged = true
              break
            }
            let name = String(incoming[i].name || incoming[i].label || '')
            if (prefix && name.startsWith(prefix)) name = name.slice(prefix.length)
            if (name !== this.rowDeviceNames[i]) {
              namesChanged = true
              break
            }
          }
          if (incoming.length !== this.rowDeviceNames.length) namesChanged = true
        }
        if (hasRows && nextIdx === curIdx && !namesChanged && !nav._forcePageReload) {
          if (nav.totalDatapoints) {
            this.setDiyOnDetail('navTotalDatapoints', String(nav.totalDatapoints))
            const size = Math.max(1, Number(nav.datapointPageSize) || this.pagerPageSize)
            this.setDiyOnDetail('navDatapointTotalPages', String(Math.max(1, Math.ceil(Number(nav.totalDatapoints) / size))))
          }
          if (fp) this._lastNavPageFingerprint = fp
          return
        }
        if (fp) this._lastNavPageFingerprint = fp
      }
      this.applySignalPageFromNav(nav)
    },
    /** 挂载后若表格无行，尝试用 store 中的 allDatapoints 补齐 */
    ensureNavDatapointRows() {
      if (!this.isNavDatapointsSource && String(readDiyValue(this.detail, 'rowSource') || '') !== 'navDatapoints') {
        return
      }
      const hasRows = (this.rowDeviceNames && this.rowDeviceNames.length)
        || (this.bindingMatrix && this.bindingMatrix.length)
      if (hasRows) return
      const nav = getEditorNavContext()
      if (nav && ((nav.allDatapoints && nav.allDatapoints.length) || (nav.datapoints && nav.datapoints.length))) {
        this.applySignalPageFromNav(nav)
        return
      }
      // 禁止在本组件 require navContext；交给 ISMRender 重拉后发 NavDatapointPageUpdate
      this.requestNavHydrate(nav)
    },
    parseBindings(rowBindings) {
      if (!rowBindings || !String(rowBindings).trim()) {
        return [];
      }
      return String(rowBindings)
        .split(';')
        .map(row => row.split(',').map(cell => cell.trim()).filter(Boolean));
    },
    applyTableConfig({ rowDeviceNames, rowDeviceCodes, columnHeaders, bindingMatrix, rowPointUuids }) {
      this.rowDeviceNames = rowDeviceNames || [];
      this.rowDeviceCodes = rowDeviceCodes || [];
      this.columnHeaders = columnHeaders || [];
      this.bindingMatrix = bindingMatrix || [];
      if (rowPointUuids !== undefined) {
        this.rowPointUuids = Array.isArray(rowPointUuids)
          ? rowPointUuids.map(u => this.normalizePointUuid(u))
          : []
      }
    },
    QueryRealData() {
      // 信号层：与设备管理相同，走 GetRealData（库元数据 + 内存实时值）
      if (this.isNavDatapointsSource || String(readDiyValue(this.detail, 'rowSource') || '') === 'navDatapoints') {
        this.refreshNavPageFromGetRealData()
        return
      }
      if (!this.bindingMatrix || !this.bindingMatrix.length) {
        return;
      }
      if (this.isNavChildrenSource) {
        const rowCount = this.bindingMatrix.length
        if (rowCount > REAL_DATA_DEFAULT_PAGE_SIZE * 2) {
          console.warn('[ViewRealTable] unexpected binding rows for nav mode:', rowCount)
        }
      }
      // 按行分批请求（每批不超过默认页大小），降低单次响应体积与浏览器内存
      const rowChunks = chunkArray(this.bindingMatrix, REAL_DATA_DEFAULT_PAGE_SIZE);
      this.messageShowLoad = true;
      const merged = [];
      const run = (idx) => {
        if (idx >= rowChunks.length) {
          this.batchUpdateConfig(merged);
          this.messageShowLoad = false;
          return;
        }
        postRealDataByBindings(rowChunks[idx])
          .then(res => {
            if (res.data && res.data.code === 0 && Array.isArray(res.data.realData)) {
              for (let i = 0; i < res.data.realData.length; i++) {
                merged.push(res.data.realData[i]);
              }
            }
            run(idx + 1);
          })
          .catch(() => {
            this.messageShowLoad = false;
            if (this._isBeingDestroyed || this._isDestroyed) return;
            const errorText = this.$i18n && typeof this.$t === 'function' ? this.$t('loginPage.serverError') : 'Server error';
            this.$message.error(errorText, 3);
          });
      };
      run(0);
    },
    /** 信号层刷新：GetRealData，与数据仓库 monitor.vue / last `_` 拆分同源 */
    refreshNavPageFromGetRealData() {
      const nav = getEditorNavContext() || {}
      const label = nav.label || nav.name || ''
      const muid = nav.modelUuid || nav.muid || ''
      const uuid = nav.uuid || nav.deviceUuid || ''
      // 虚拟设备：禁止 uuid+namePrefix（后端是 OR，会拉回整台设备测点）
      const virtualCabinet = String(nav.virtualCabinet || '').trim()
      const parentLabel = String(nav.parentDeviceLabel || '').trim()
      const isFallback = !!(
        nav.virtualCabinetFallback
        || nav.isFallbackGroup
        || (virtualCabinet && parentLabel && virtualCabinet === parentLabel)
      )
      if (!label && !uuid) {
        if (this.bindingMatrix && this.bindingMatrix.length) {
          this._queryRealDataByBindingsOnly()
        }
        return
      }
      const requestedPageIndex = Number(nav.datapointPageIndex) || 0
      const page = requestedPageIndex + 1
      const pageSize = Math.max(1, Number(nav.datapointPageSize) || 20)
      const requestId = ++this.navRefreshRequestId
      this.messageShowLoad = true
      const query = String(nav.datapointQuery || '').trim() || undefined
      const stripLabel = virtualCabinet || label
      let pointBelongsToVirtualDevice = null
      let displayPointNameForVirtualDevice = null
      try {
        const vcUtil = require('@/pages/ISMDisPlay/utils/virtualCabinet')
        pointBelongsToVirtualDevice = vcUtil.pointBelongsToVirtualDevice
        displayPointNameForVirtualDevice = vcUtil.displayPointNameForVirtualDevice
      } catch (e) { /* ignore */ }

      const finishWithRows = (rawRows, bodyTotal) => {
        if (this._isBeingDestroyed || this._isDestroyed || requestId !== this.navRefreshRequestId) return
        const currentNav = getEditorNavContext() || {}
        if ((Number(currentNav.datapointPageIndex) || 0) !== requestedPageIndex) return
        this.messageShowLoad = false
        const rows = (rawRows || []).filter(r => {
          const n = String(r.name || '').trim()
          if (!n || /^device\./i.test(n) || /^system\./i.test(n)) return false
          if (virtualCabinet && typeof pointBelongsToVirtualDevice === 'function') {
            return pointBelongsToVirtualDevice(n, virtualCabinet, isFallback)
          }
          if (virtualCabinet) {
            const vcPrefix = virtualCabinet.endsWith('_') ? virtualCabinet : `${virtualCabinet}_`
            if (!n.startsWith(vcPrefix) && n !== virtualCabinet) {
              if (!(isFallback && n.indexOf('_') < 0)) return false
            }
          }
          return true
        })
        const units = rows.map(r => r.unit || r.DataUnit || '')
        const names = rows.map(r => {
          const n = String(r.name || '')
          if (virtualCabinet && typeof displayPointNameForVirtualDevice === 'function') {
            return displayPointNameForVirtualDevice(n, stripLabel)
          }
          const prefix = stripLabel ? `${stripLabel}_` : ''
          if (prefix && n.startsWith(prefix)) return n.slice(prefix.length)
          return n
        })
        const rowPointUuids = rows.map(r => this.normalizePointUuid(r.uuid || r.Uuid))
        const bindings = rows.map(r => {
          const owner = r.DeviceName || r.device_name || ''
          const n = r.name || ''
          return owner ? [`${owner}->${n}`] : [n]
        })
        const values = rows.map(r => {
          const v = r.value
          return (v === undefined || v === null || v === '') ? ['-'] : [v]
        })
        const prevKeys = this.currentRowKeys()
        const prevCellData = this.cellData
        const nextKeys = rows.map((r, i) => this.pointRowKey(
          rowPointUuids[i],
          bindings[i] && bindings[i][0],
          names[i],
          i,
        ))
        if (names.length) {
          this.applyTableConfig({
            rowDeviceNames: names,
            rowDeviceCodes: units,
            columnHeaders: ['实时值'],
            bindingMatrix: bindings,
            rowPointUuids,
          })
          this.setDiyOnDetail('rowDeviceNames', names.join('\n'))
          this.setDiyOnDetail('rowDeviceCodes', units.join('\n'))
        }
        const navTotal = Number(bodyTotal) || Number(nav.totalDatapoints) || 0
        if (navTotal > 0) {
          this.setDiyOnDetail('navTotalDatapoints', String(navTotal))
          const size = pageSize
          const totalPages = Math.max(1, Math.ceil(navTotal / size))
          this.setDiyOnDetail('navDatapointTotalPages', String(totalPages))
          const cachedLen = Array.isArray(nav.allDatapoints) ? nav.allDatapoints.length : 0
          const shouldServerPage = !!nav.serverPaged || !cachedLen || navTotal > cachedLen
          const navPatch = {
            totalDatapoints: navTotal,
            datapointTotalPages: totalPages,
            datapointPageSize: size,
            datapointPageIndex: requestedPageIndex,
          }
          if (shouldServerPage) {
            navPatch.serverPaged = true
            navPatch.allDatapoints = []
          }
          commitEditorNav(navPatch)
          this.$EventBus.$emit('NavDatapointPageUpdate', {
            ...(getEditorNavContext() || {}),
            ...navPatch,
            _metaOnly: true,
          })
        }
        this.cellData = this.mergeCellDataPreserve(prevCellData, values, nextKeys, prevKeys)
        this._navPagePointUuids = rowPointUuids.filter(Boolean)
      }

      const fail = () => {
        if (requestId !== this.navRefreshRequestId) return
        const currentNav = getEditorNavContext() || {}
        if ((Number(currentNav.datapointPageIndex) || 0) !== requestedPageIndex) return
        this.messageShowLoad = false
        if (this.bindingMatrix && this.bindingMatrix.length) {
          this._queryRealDataByBindingsOnly()
        }
      }

      // fallback（无前缀点归属真设备名）：走 navContext 扫描过滤，避免 category 漏点
      if (isFallback && virtualCabinet && uuid) {
        try {
          const { fetchDeviceDatapointPage } = require('@/pages/ISMDisPlay/utils/navContext')
          fetchDeviceDatapointPage({
            muid,
            deviceUuid: uuid,
            pointNamePrefix: virtualCabinet,
            isFallbackGroup: true,
            page,
            pageSize,
            query: query || '',
          }).then((pointPage) => {
            finishWithRows(pointPage.points || [], pointPage.total)
          }).catch(fail)
        } catch (e) {
          fail()
        }
        return
      }

      // 虚拟设备：uuid + category AND；普通设备：uuid + namePrefix
      const payload = (virtualCabinet && uuid)
        ? {
          uuid,
          muid: muid || undefined,
          page,
          pageSize,
          query,
          category: virtualCabinet.endsWith('_') ? virtualCabinet : `${virtualCabinet}_`,
          IsRemoveGW: false,
        }
        : {
          uuid: uuid || undefined,
          muid: muid || undefined,
          namePrefix: label || undefined,
          deviceLabel: label || undefined,
          page,
          pageSize,
          query,
          category: String(nav.datapointCategory || '').trim() || undefined,
          IsRemoveGW: false,
        }
      postGetRealData(payload).then((res) => {
        const body = res && res.data
        if (!body || body.code !== 0 || !Array.isArray(body.realData)) {
          if (requestId === this.navRefreshRequestId) this.messageShowLoad = false
          return
        }
        finishWithRows(body.realData, body.total)
      }).catch(fail)
    },
    _queryRealDataByBindingsOnly() {
      if (!this.bindingMatrix || !this.bindingMatrix.length) return
      const rowChunks = chunkArray(this.bindingMatrix, REAL_DATA_DEFAULT_PAGE_SIZE)
      this.messageShowLoad = true
      const merged = []
      const run = (idx) => {
        if (idx >= rowChunks.length) {
          this.batchUpdateConfig(merged)
          this.messageShowLoad = false
          return
        }
        postRealDataByBindings(rowChunks[idx])
          .then(res => {
            if (res.data && res.data.code === 0 && Array.isArray(res.data.realData)) {
              for (let i = 0; i < res.data.realData.length; i++) {
                merged.push(res.data.realData[i])
              }
            }
            run(idx + 1)
          })
          .catch(() => {
            this.messageShowLoad = false
          })
      }
      run(0)
    },
    initComponents(option) {
      if (!option || !option.style || !option.style.position) {
        return
      }
      this.width = option.style.position.w;
      this.height = option.style.position.h;
      this.foreColor = option.style.foreColor;
      this.backColor = option.style.backColor;
      this.fontSize = option.style.fontSize;
      this.fontFamily = option.style.fontFamily;

      let columnHeadersText = '';
      let rowDeviceNamesText = '';
      let rowDeviceCodesText = '';
      let rowBindingsText = '';

      for (let i = 0; i < (option.style.diy || []).length; i += 1) {
        const item = option.style.diy[i];
        if (item.key === 'columnHeaders' && item.value) {
          columnHeadersText = item.value;
        } else if (item.key === 'rowDeviceNames' && item.value) {
          rowDeviceNamesText = item.value;
        } else if (item.key === 'rowDeviceCodes' && item.value) {
          rowDeviceCodesText = item.value;
        } else if (item.key === 'rowBindings' && item.value) {
          rowBindingsText = item.value;
        } else if (item.key === 'tableHeaderColor') {
          this.tableHeaderColor = item.value;
        } else if (item.key === 'tableHeaderBackColor') {
          this.tableHeaderBackColor = item.value;
        } else if (item.key === 'tableSplitColor') {
          this.tableSplitColor = item.value;
        } else if (item.key === 'tableHoverColor') {
          this.tableHoverColor = item.value;
        } else if (item.key === 'ShowCount') {
          this.pagination.pageSize = clampPageSize(item.value, REAL_DATA_DEFAULT_PAGE_SIZE);
        } else if (item.key === 'waitTime') {
          this.waitTime = item.value;
        } else if (item.key === 'SearchColor') {
          this.SearchColor = item.value;
        } else if (item.key === 'SearchBackColor') {
          this.SearchBackColor = item.value;
        } else if (item.key === 'SearchBorderColor') {
          this.SearchBorderColor = item.value;
        } else if (item.key === 'tableHeaderFont') {
          this.tableHeaderFont = item.value;
        } else if (item.key === 'tableHeaderFontSize') {
          this.tableHeaderFontSize = item.value;
        } else if (item.key === 'scrollBgColor') {
          this.scrollBgColor = item.value;
        } else if (item.key === 'scrollFrColor') {
          this.scrollFrColor = item.value;
        } else if (item.key === 'scrollHdColor') {
          this.scrollHdColor = item.value;
        } else if (item.key === 'themeName') {
          this.selectedTheme = item.value || 'light';
        }
      }

      this.applyTableConfig({
        rowDeviceNames: rowDeviceNamesText
          .split('\n')
          .map(v => v.trim())
          .filter(Boolean),
        rowDeviceCodes: rowDeviceCodesText
          .split('\n')
          .map(v => v.trim())
          .filter(Boolean),
        columnHeaders: columnHeadersText
          .split(',')
          .map(v => v.trim())
          .filter(Boolean),
        bindingMatrix: this.parseBindings(rowBindingsText)
      })

      const rowSourceItem = (option.style.diy || []).find(d => d && d.key === 'rowSource')
      const navRowSource = rowSourceItem && String(rowSourceItem.value)
      if (navRowSource === 'navDatapoints' || navRowSource === 'navChildren') {
        const rowCount = Math.max(this.rowDeviceNames.length, this.bindingMatrix.length)
        const colCount = Math.max(1, (this.columnHeaders && this.columnHeaders.length) || 1)
        const prev = Array.isArray(this.cellData) ? this.cellData : []
        // 结构变化时扩/缩行，但绝不把已有有效值整表刷成 '-'（录屏里一闪而过的主因之一）
        this.cellData = Array.from({ length: rowCount }, (_, rowIndex) => {
          const prevRow = Array.isArray(prev[rowIndex]) ? prev[rowIndex] : []
          return Array.from({ length: colCount }, (__, colIndex) => {
            const kept = prevRow[colIndex]
            return this.isEmptyCellValue(kept) ? '-' : kept
          })
        })
      }

      this.$nextTick(() => {
        this.applyScrollbarTheme();
      });
      this.animateType = option.animate.selected || [];
      this.isStart = !option.animate.isExpression;
      clearInterval(this.AlarmTimer);
      if (!this.editMode && !this.IsToolBox) {
        // 推送为主，轮询作兜底（默认至少 8s，降低大屏多表叠加压力）
        if (this.waitTime < 8000) {
          this.waitTime = 8000;
        }
        this.AlarmTimer = setInterval(this.QueryRealData, this.waitTime);
      }
    },
    applyRealtimePush(pushData) {
      if (this.editMode || this.IsToolBox || !pushData || !Array.isArray(pushData.Data) || !pushData.Data.length) {
        return
      }
      if (!Array.isArray(this.cellData) || !this.cellData.length) {
        return
      }
      // 信号层：只吃当前设备的推送，避免同型号其它设备短名串写到本页卡片
      const nav = getEditorNavContext() || {}
      const currentDeviceUuid = this.normalizePointUuid(nav.deviceUuid || nav.uuid || '')
      const pushDeviceUuid = this.normalizePointUuid(pushData.DeviceUuid || pushData.deviceUuid || '')
      if (this.isNavDatapointsSource && currentDeviceUuid && pushDeviceUuid && currentDeviceUuid !== pushDeviceUuid) {
        return
      }
      const virtualCabinet = String(nav.virtualCabinet || '').trim()
      const vcPrefix = virtualCabinet
        ? (virtualCabinet.endsWith('_') ? virtualCabinet : `${virtualCabinet}_`)
        : ''

      // 一律按点位 uuid（小写）匹配；名称仅在「本行无 uuid」时用全限定名兜底，禁止裸短名
      const byUuid = Object.create(null)
      const byQualifiedName = Object.create(null)
      for (let i = 0; i < pushData.Data.length; i++) {
        const d = pushData.Data[i]
        const uuid = this.normalizePointUuid(d.Uuid || d.uuid)
        if (uuid) {
          byUuid[uuid] = d.Value
        }
        const name = String(d.DataName || d.Name || d.name || '').trim()
        if (name) {
          byQualifiedName[name] = d.Value
        }
      }
      if (!Object.keys(byUuid).length && !Object.keys(byQualifiedName).length) {
        return
      }
      let changed = false
      const deviceName = String(pushData.DeviceName || '').trim()
      const next = this.cellData.map((row, rowIndex) => {
        const rowCopy = Array.isArray(row) ? row.slice() : [row]
        const rowUuid = this.normalizePointUuid(this.rowPointUuids && this.rowPointUuids[rowIndex])
        const rowName = String((this.rowDeviceNames && this.rowDeviceNames[rowIndex]) || '').trim()
        const bindingCell = this.bindingMatrix && this.bindingMatrix[rowIndex] && this.bindingMatrix[rowIndex][0]
        let newVal
        if (rowUuid && byUuid[rowUuid] !== undefined) {
          newVal = byUuid[rowUuid]
        }
        // 有 uuid 却推送里对不上：宁可保持 HTTP 真值，也不用短名猜（现场乱跳主因）
        if (newVal === undefined && !rowUuid) {
          const candidates = []
          if (bindingCell) {
            const bindStr = String(bindingCell)
            const parts = bindStr.split('->')
            const pointName = parts.length > 1 ? parts[parts.length - 1] : bindStr
            candidates.push(pointName, bindStr)
          }
          if (vcPrefix && rowName) {
            candidates.push(`${vcPrefix}${rowName}`, `${virtualCabinet}_${rowName}`)
          }
          if (deviceName && rowName) {
            candidates.push(`${deviceName}_${rowName}`, `${deviceName}->${rowName}`)
          }
          // 非信号层旧表才允许裸行名；信号层禁止，避免多设备同短名串值
          if (!this.isNavDatapointsSource && rowName) {
            candidates.push(rowName)
          }
          for (let c = 0; c < candidates.length; c += 1) {
            const key = candidates[c]
            if (key && byQualifiedName[key] !== undefined) {
              newVal = byQualifiedName[key]
              break
            }
          }
        }
        if (newVal !== undefined && rowCopy[0] !== newVal) {
          rowCopy[0] = newVal
          changed = true
        }
        return rowCopy
      })
      if (changed) {
        this.cellData = next
      }
    },
  },
  beforeDestroy() {
    clearInterval(this.AlarmTimer);
    this.stopDeviceOnlinePolling()
    this.clearAutoPaging()
    if (this._onReadDataPush && this.$EventBus) {
      this.$EventBus.$off('readDataPush', this._onReadDataPush)
      this._onReadDataPush = null
    }
    if (this._onAnyNavPageChange && this.$EventBus) {
      this.$EventBus.$off('NavPageChange', this._onAnyNavPageChange)
      this._onAnyNavPageChange = null
    }
    if (this._onNavDatapointPageUpdate && this.$EventBus) {
      this.$EventBus.$off('NavDatapointPageUpdate', this._onNavDatapointPageUpdate)
      this._onNavDatapointPageUpdate = null
    }
    if (typeof document !== 'undefined' && this.scrollbarStyleTagId) {
      const styleTag = document.getElementById(this.scrollbarStyleTagId);
      if (styleTag) {
        styleTag.remove();
      }
    }
  },
  mounted() {
    this.$nextTick(() => {
      this.initComponents(this.detail);
      this.applyScrollbarTheme();
      // 测点表若 diy 行为空，从 store navContext 页内补灌（避免 No Data）
      this.ensureNavDatapointRows()
      if (!this.editMode && !this.IsToolBox) {
        this.QueryRealData();
        this.startDeviceOnlinePolling()
        this._onReadDataPush = (data) => {
          this.applyRealtimePush(data)
        }
        this.$EventBus.$on('readDataPush', this._onReadDataPush)
      }
      this.startAutoPaging()
      const activeEvent = `${this.detail.identifier}activeEvent`;
      const animateEvent = `${this.detail.identifier}animateEvent`;
      this.$EventBus.$on(activeEvent, () => {});
      this.$EventBus.$on(animateEvent, data => {
        this.isStart = data;
      });
    });
  },
  created() {
    this.GetNodeObj = this.getNode();
    this.GetNodeObj.on('change:data', ({ current }) => {
      if (current) {
        this.detail = current.detail;
        if (!this.editMode && !this.IsToolBox) {
          this.initComponents(this.detail);
          this.$nextTick(() => {
            this.QueryRealData();
          });
        }
      }
    });
    this.GetNodeObj.on('change:size', ({ current }) => {
      this.detail.style.position.w = current.width;
      this.detail.style.position.h = current.height;
    });
    this.detail = this.GetNodeObj.getData().detail;
    this.editMode = this.GetNodeObj.getData().editMode;
    this.showDeviceUuid = this.GetNodeObj.getData().showDeviceUuid;
    this.IsToolBox = this.GetNodeObj.getData().IsToolBox;
    this.$EventBus.$on('cell-editMode', data => {
      this.editMode = data.edit;
      this.IsToolBox = data.toolbox;
      this.initComponents(this.detail);
      this.$nextTick(() => {
        if (this.editMode || this.IsToolBox) this.clearAutoPaging()
        else this.startAutoPaging()
      })
    });
    this._onAnyNavPageChange = payload => {
      if (!payload || !payload.autoPage) this.pauseAutoPaging()
    }
    this.$EventBus.$on('NavPageChange', this._onAnyNavPageChange)
    this._onNavDatapointPageUpdate = (nav) => this.onNavDatapointPageUpdate(nav)
    this.$EventBus.$on('NavDatapointPageUpdate', this._onNavDatapointPageUpdate)
    this.initComponents(this.detail);
  }
};
</script>

<style lang="less" scoped>
.history-theme-shell {
  box-sizing: border-box;
  width: 100%;
  height: 100%;
  padding: 12px;
  color: var(--toolbarText);
  background: var(--panelBg);
  border: 1px solid var(--panelBorder);
  border-radius: 16px;
  box-shadow: var(--panelShadow);
}

.table-container {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-sizing: border-box;
  border: 1px solid var(--panelBorder);
  border-radius: 12px;
  background: var(--backColor, transparent);
}

.table-scroll {
  flex: 1 1 auto;
  min-height: 0;
  overflow: auto;
  box-sizing: border-box;
}

::v-deep .ant-table-pagination.ant-pagination {
  margin: 0 !important;
  float: none !important;
  text-align: left !important;
  padding: 0 !important;
  line-height: 1 !important;
}

::v-deep .ant-table {
  color: var(--foreColor);
  background: transparent;
}

::v-deep .ant-table table,
::v-deep .ant-table-content table,
::v-deep .ant-table-header table,
::v-deep .ant-table-body table {
  table-layout: auto !important;
  min-width: max-content !important;
}

::v-deep .ant-table-header {
  overflow: hidden !important;
  margin-bottom: -1px !important;
}

::v-deep .ant-table-thead > tr > th {
  position: sticky !important;
  top: 0;
  z-index: 6;
}

::v-deep .ant-table-thead > tr > th {
  color: var(--tableHeaderColor) !important;
  font-size: var(--tableHeaderFontSize) !important;
  font-family: var(--tableHeaderFont) !important;
  background: var(--tableHeaderBackColor) !important;
  border-bottom: 1px solid var(--tableSplitColor) !important;
  border-right: 1px solid var(--tableColumnSplitColor) !important;
  white-space: nowrap !important;
  padding: 10px 10px;
}

::v-deep .ant-table-tbody > tr > td {
  color: var(--foreColor) !important;
  font-size: var(--fontSize) !important;
  font-family: var(--fontFamily) !important;
  background: var(--backColor) !important;
  border-bottom: 1px solid var(--tableSplitColor) !important;
  border-right: 1px solid var(--tableColumnSplitColor) !important;
  white-space: nowrap !important;
  padding: 7px 5px;
}

::v-deep .ant-table-thead > tr > th:last-child,
::v-deep .ant-table-tbody > tr > td:last-child {
  border-right: none !important;
}

::v-deep .ant-table-thead > tr > th.sticky-col,
::v-deep .ant-table-tbody > tr > td.sticky-col {
  position: sticky !important;
  z-index: 3;
}

::v-deep .ant-table-thead > tr > th.sticky-col {
  z-index: 8;
}

::v-deep .sticky-col-0 {
  left: 0;
}

::v-deep .sticky-col-1 {
  left: 60px;
}

::v-deep .sticky-col-2 {
  left: 180px;
}

::v-deep .ant-table-tbody > tr:nth-child(odd) > td {
  background: var(--tableRowOddBg) !important;
}

::v-deep .ant-table-tbody > tr:nth-child(even) > td {
  background: var(--tableRowEvenBg) !important;
}

::v-deep .ant-table-tbody > tr:hover > td {
  background: var(--tableHoverColor) !important;
}

::v-deep .ant-pagination-item,
::v-deep .ant-pagination-prev,
::v-deep .ant-pagination-next,
::v-deep .ant-pagination-item-link {
  color: var(--SearchColor) !important;
  background: var(--SearchBackColor) !important;
  border-color: var(--SearchBorderColor) !important;
}

::v-deep .ant-pagination-item-active {
  background: var(--toolbarAccentSoft) !important;
  border-color: var(--toolbarAccent) !important;
}
</style>
