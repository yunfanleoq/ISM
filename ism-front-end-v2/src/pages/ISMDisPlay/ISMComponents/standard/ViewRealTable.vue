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
        <div :class="['history-theme-shell', scrollbarThemeClass]" :style="styleVar">
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
            <div
              class="table-pagination-bar"
              :class="{ visible: showPaginationBar }"
              v-if="paginationTotal > 0"
              @mousedown.stop
              @click.stop
            >
              <div class="rt-pager" role="navigation" aria-label="表格分页">
                <button
                  type="button"
                  class="rt-pager-btn"
                  :disabled="!pagerCanPrev"
                  @click.stop.prevent="goPagerPrev"
                >
                  上一页
                </button>
                <span class="rt-pager-info">{{ pagerInfoText }}</span>
                <button
                  type="button"
                  class="rt-pager-btn"
                  :disabled="!pagerCanNext"
                  @click.stop.prevent="goPagerNext"
                >
                  下一页
                </button>
              </div>
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
    const cur = (st.state.ISMDisPlayEditorTool && st.state.ISMDisPlayEditorTool.navContext) || {}
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
          }
        })
      },
      deep: true
    }
  },
  computed: {
    currentTheme() {
      return THEME_MAP[this.selectedTheme] || THEME_MAP.light;
    },
    scrollbarThemeClass() {
      const id = this.detail && this.detail.identifier ? this.detail.identifier : 'default';
      return `real-table-scrollbar-${id}`;
    },
    tableScrollWidth() {
      return this.dynamicColumns.reduce((total, column) => total + (Number(column.width) || 120), 0);
    },
    styleVar() {
      return {
        '--panelBg': this.currentTheme.panelBg,
        '--panelBorder': this.currentTheme.panelBorder,
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
            cellValue = rowCellData[colIndex];
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
      if (!this.isNavChildrenSource) return null
      const nav = getEditorNavContext()
      return nav && nav.deviceListMode ? nav : null
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
      showPaginationBar: true,
      rowDeviceNames: ['空调', '灯光', '门禁', '监控'],
      rowDeviceCodes: ['AC001', 'LT002', 'AC003', 'CAM004'],
      columnHeaders: ['运行状态', '功率', '温度', '在线时长'],
      bindingMatrix: [
        ['AC001->status', 'AC001->power', 'AC001->temp', 'AC001->online_time'],
        ['LT002->status', 'LT002->power', 'LT002->temp', 'LT002->online_time'],
        ['AC003->status', 'AC003->power', 'AC003->temp', 'AC003->online_time'],
        ['CAM004->status', 'CAM004->power', 'CAM004->temp', 'CAM004->online_time']
      ],
      cellData: [
        [100, 200, 300, 400],
        [150, 250, 350, 450],
        [120, 220, 320, 420],
        [180, 280, 380, 480]
      ],
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
      waitTime: 1000,
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
    batchUpdateConfig(newCellData) {
      this.cellData = newCellData;
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
    handlePageChange(page) {
      const nextPage = Math.max(1, Number(page) || 1)
      const totalPages = this.pagerTotalPages
      if (nextPage < 1 || nextPage > totalPages) return

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
        this.$EventBus.$emit('NavPageChange', { datapointPageIndex: next })
        return
      }
      if (this.isNavChildrenSource || this.navListPagination) {
        const next = nextPage - 1
        const nav = this.navListPagination || getEditorNavContext()
        if (next === ((nav && nav.pageIndex) || 0)) return
        this.$EventBus.$emit('NavPageChange', { pageIndex: next })
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
    /** 页内翻页：只换当前页测点行并拉实时值，不整页重载 */
    applySignalPageFromNav(nav) {
      if (!nav || this.editMode || this.IsToolBox) return
      if (!this.isNavDatapointsSource && String(readDiyValue(this.detail, 'rowSource') || '') !== 'navDatapoints') {
        return
      }
      // 必须用全量 allDatapoints；datapoints 可能已是当前页切片
      let all = Array.isArray(nav.allDatapoints) && nav.allDatapoints.length
        ? nav.allDatapoints
        : (nav.datapoints || nav.childNodes || [])
      const declaredTotal = Number(nav.totalDatapoints) || 0
      // store 总数大于当前 all 长度时，说明 all 被截成页切片，需回源补齐
      if (declaredTotal > all.length) {
        this.$EventBus.$emit('NavDatapointNeedHydrate')
        // 仍用现有行先渲染，避免空白
      }
      if (!all.length) {
        return
      }
      const size = Math.max(1, Number(nav.datapointPageSize) || 20)
      const total = Math.max(all.length, declaredTotal)
      const totalPages = Math.max(1, Math.ceil(total / size) || 1)
      const idx = Math.max(0, Math.min(Number(nav.datapointPageIndex) || 0, totalPages - 1))
      const pagePoints = all.slice(idx * size, idx * size + size)
      const deviceLabel = nav.name || nav.label || ''
      const prefix = deviceLabel ? `${deviceLabel}_` : ''
      const rowDeviceNames = pagePoints.map(p => {
        const n = String(p.name || p.label || '')
        if (prefix && n.startsWith(prefix)) return n.slice(prefix.length)
        return n
      })
      const rowDeviceCodes = pagePoints.map(p => p.unit || '')
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

      this.applyTableConfig({
        rowDeviceNames,
        rowDeviceCodes,
        columnHeaders: ['实时值'],
        bindingMatrix,
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
      commitEditorNav({
        signalMode: true,
        routeMode: 'signal',
        datapointPageIndex: idx,
        datapointPageSize: size,
        totalDatapoints: total,
        datapointTotalPages: totalPages,
        allDatapoints: all,
        datapoints: pagePoints,
      })
      // 先用 GetRealData 已覆盖的内存值填表（与设备管理同源），再异步刷新
      this.cellData = pagePoints.map(p => {
        const v = p.value
        if (v === undefined || v === null || v === '') return ['-']
        return [v]
      })
      this._navPagePointUuids = pagePoints.map(p => p.uuid || '').filter(Boolean)
      this.$nextTick(() => {
        this.QueryRealData()
      })
    },
    onNavDatapointPageUpdate(nav) {
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
      this.$EventBus.$emit('NavDatapointNeedHydrate')
    },
    parseBindings(rowBindings) {
      if (!rowBindings || !String(rowBindings).trim()) {
        return [];
      }
      return String(rowBindings)
        .split(';')
        .map(row => row.split(',').map(cell => cell.trim()).filter(Boolean));
    },
    applyTableConfig({ rowDeviceNames, rowDeviceCodes, columnHeaders, bindingMatrix }) {
      this.rowDeviceNames = rowDeviceNames || [];
      this.rowDeviceCodes = rowDeviceCodes || [];
      this.columnHeaders = columnHeaders || [];
      this.bindingMatrix = bindingMatrix || [];
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
    /** 信号层刷新：GetRealData + namePrefix，与 monitor.vue 同源 */
    refreshNavPageFromGetRealData() {
      const nav = getEditorNavContext() || {}
      const label = nav.label || nav.name || ''
      const muid = nav.modelUuid || nav.muid || ''
      const uuid = nav.uuid || nav.deviceUuid || ''
      if (!label && !uuid) {
        // 回退绑点解析
        if (this.bindingMatrix && this.bindingMatrix.length) {
          this._queryRealDataByBindingsOnly()
        }
        return
      }
      const page = (Number(nav.datapointPageIndex) || 0) + 1
      const pageSize = Math.max(1, Number(nav.datapointPageSize) || 20)
      this.messageShowLoad = true
      postGetRealData({
        // 有逻辑设备名时只按前缀查，避免 OR uuid 混入 device.DeviceStatus 等系统点
        uuid: label ? undefined : (uuid || undefined),
        muid: muid || undefined,
        namePrefix: label || undefined,
        deviceLabel: label || undefined,
        page,
        pageSize,
        IsRemoveGW: false,
      }).then((res) => {
        this.messageShowLoad = false
        if (this._isBeingDestroyed || this._isDestroyed) return
        const body = res && res.data
        if (!body || body.code !== 0 || !Array.isArray(body.realData)) return
        const rows = (body.realData || []).filter(r => {
          const n = String(r.name || '').trim()
          return n && !/^device\./i.test(n) && !/^system\./i.test(n)
        })
        // 同步单位（库）与实时值（内存已覆盖）
        const units = rows.map(r => r.unit || r.DataUnit || '')
        const values = rows.map(r => {
          const v = r.value
          return (v === undefined || v === null || v === '') ? ['-'] : [v]
        })
        const names = rows.map(r => {
          const n = String(r.name || '')
          const prefix = label ? `${label}_` : ''
          if (prefix && n.startsWith(prefix)) return n.slice(prefix.length)
          return n
        })
        if (names.length) {
          this.applyTableConfig({
            rowDeviceNames: names,
            rowDeviceCodes: units,
            columnHeaders: ['实时值'],
            bindingMatrix: rows.map(r => {
              const owner = r.DeviceName || r.device_name || ''
              const n = r.name || ''
              return owner ? [`${owner}->${n}`] : [n]
            }),
          })
          this.setDiyOnDetail('rowDeviceNames', names.join('\n'))
          this.setDiyOnDetail('rowDeviceCodes', units.join('\n'))
        }
        // 总数以客户端 allDatapoints / 已声明为准，避免服务端分页与本地切片不一致改写页码
        const navTotal = Number(nav.totalDatapoints) || Number(body.total) || 0
        if (navTotal > 0) {
          this.setDiyOnDetail('navTotalDatapoints', String(navTotal))
          const size = pageSize
          const totalPages = Math.max(1, Math.ceil(navTotal / size))
          this.setDiyOnDetail('navDatapointTotalPages', String(totalPages))
          commitEditorNav({
            totalDatapoints: navTotal,
            datapointTotalPages: totalPages,
            datapointPageSize: size,
            datapointPageIndex: page - 1,
          })
          // 通知顶部页码同步
          this.$EventBus.$emit('NavDatapointPageUpdate', {
            ...(getEditorNavContext() || {}),
            totalDatapoints: navTotal,
            datapointTotalPages: totalPages,
            datapointPageSize: size,
            datapointPageIndex: page - 1,
          })
        }
        this.batchUpdateConfig(values)
      }).catch(() => {
        this.messageShowLoad = false
        // 失败时尝试绑点回退
        if (this.bindingMatrix && this.bindingMatrix.length) {
          this._queryRealDataByBindingsOnly()
        }
      })
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
        this.cellData = Array.from({ length: rowCount }, () =>
          this.columnHeaders.map(() => '-'),
        )
      }

      this.$nextTick(() => {
        this.applyScrollbarTheme();
      });
      this.animateType = option.animate.selected || [];
      this.isStart = !option.animate.isExpression;
      clearInterval(this.AlarmTimer);
      if (!this.editMode && !this.IsToolBox) {
        if (this.waitTime < 1000) {
          this.waitTime = 1000;
        }
        this.AlarmTimer = setInterval(this.QueryRealData, this.waitTime);
      }
    }
  },
  beforeDestroy() {
    clearInterval(this.AlarmTimer);
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
      }
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
    });
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

.table-pagination-bar {
  flex: 0 0 auto;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 0;
  height: 0;
  padding: 0;
  margin: 0 auto;
  width: 100%;
  border: none;
  border-radius: 0;
  background: transparent;
  backdrop-filter: none;
  box-shadow: none;
  opacity: 0;
  transform: translateY(4px);
  pointer-events: none;
  overflow: visible;
  transition: opacity 0.18s ease, transform 0.18s ease;
  z-index: 5;
}

.table-pagination-bar.visible {
  min-height: 48px;
  height: auto;
  padding: 8px 12px 10px;
  margin: 0 auto;
  width: 100%;
  opacity: 1;
  transform: translateY(0);
  pointer-events: auto;
}

.rt-pager {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  min-height: 44px;
  padding: 6px 12px;
  border-radius: 10px;
  background: rgba(15, 35, 58, 0.85);
  border: 1px solid rgba(0, 200, 255, 0.35);
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.22);
  pointer-events: auto;
  user-select: none;
  transform: translateZ(0);
}

.rt-pager-btn {
  min-width: 96px;
  height: 40px;
  padding: 0 18px;
  border-radius: 8px;
  border: 1px solid rgba(0, 229, 255, 0.45);
  background: linear-gradient(180deg, rgba(20, 90, 140, 0.95) 0%, rgba(12, 58, 98, 0.95) 100%);
  color: #9ef0ff;
  font-size: 15px;
  font-weight: 600;
  letter-spacing: 0.5px;
  cursor: pointer;
  line-height: 38px;
  outline: none;
  pointer-events: auto;
  touch-action: manipulation;
  transition: background 0.15s ease, border-color 0.15s ease, opacity 0.15s ease;
}

.rt-pager-btn:hover:not(:disabled) {
  border-color: rgba(120, 240, 255, 0.85);
  background: linear-gradient(180deg, rgba(28, 120, 180, 0.98) 0%, rgba(16, 72, 120, 0.98) 100%);
  color: #e8fbff;
}

.rt-pager-btn:active:not(:disabled) {
  transform: translateY(1px);
}

.rt-pager-btn:disabled {
  opacity: 0.38;
  cursor: not-allowed;
  border-color: rgba(90, 120, 150, 0.35);
  background: rgba(20, 40, 60, 0.55);
  color: #7a93ad;
}

.rt-pager-info {
  min-width: 160px;
  text-align: center;
  color: #c8ddf5;
  font-size: 13px;
  font-weight: 500;
  line-height: 1.3;
  white-space: nowrap;
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
