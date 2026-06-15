<template>
  <div :class="['history-theme-shell', scrollbarThemeClass]" :style="styleVar">
    <div class="toolbar-row">
      <a-range-picker
        :defaultValue="[moment().subtract(1, 'day'), moment()]"
        @change="onDateRangeChange"
        size="small"
        :dropdownClassName="overlayThemeClass"
        :placeholder="['开始时间', '结束时间']"
        :show-time="{ format: 'HH:mm', defaultValue: [moment('00:00:00', 'HH:mm:ss'), moment('23:59:59', 'HH:mm:ss')] }"
        :format="dateTimeDisplayFormat"
        class="date-picker"
      />
      <div class="toolbar-actions">
        <a-button
          :disabled="messageShowLoad"
          type="primary"
          size="small"
          @click="QueryHistoryDataDifferent"
          class="btn-query"
        >
          <a-icon type="search" />
          {{ $t('reporting.AlarmHistory.Query') }}
        </a-button>
        <a-button
          :disabled="isLoadExecl || !hasQueriedData"
          type="default"
          size="small"
          @click="handleExport"
          class="btn-export"
        >
          <a-icon type="export" />
          {{ $t('reporting.AlarmHistory.Export') }}
        </a-button>
      </div>
    </div>

    <div class="table-container">
      <div
        class="table-scroll"
        ref="tableScroll"
      >
        <a-table
          :key="tableKey"
          :columns="dynamicColumns"
          :data-source="pagedData"
          :rowKey="(record, index) => `${record.rowName || `${record.deviceName}-${record.dataName}`}-${index}`"
          :rowClassName="record => record.isSummaryRow ? 'summary-row' : ''"
          :pagination="false"
          :style="{ minWidth: `${tableScrollWidth}px` }"
        />
      </div>
      <div class="table-pagination-bar" :class="{ visible: showPaginationBar }" v-if="paginationTotal > 0">
        <a-pagination
          :current="pagination.current"
          :pageSize="pagination.pageSize"
          :total="paginationTotal"
          :showSizeChanger="false"
          :hideOnSinglePage="false"
          :showLessItems="true"
          simple
          @change="handlePageChange"
        />
      </div>
    </div>
  </div>
</template>

<script>
import moment from 'moment';
import { formatDate } from '@/utils/common';
import 'moment/locale/zh-cn';
import 'moment/locale/en-ie';
import 'moment/locale/zh-tw';
import { GetHistoryDiyDifferenceReport } from '@/services/report';
import { exportExcelWithStyle } from '@/services/excelExport';
import ISMChildAutoMixin from '@/mixins/ISMChildAutoMixin';

const VNodes = {
  functional: true,
  render: (h, ctx) => ctx.props.vnodes
};

const THEME_OPTIONS = [
  { label: '极简亮色', value: 'light' },
  { label: '深空夜幕', value: 'dark' },
  { label: '海岸蓝调', value: 'ocean' },
  { label: '琥珀暖光', value: 'amber' },
  { label: '森林翠影', value: 'emerald' }
];

const DATE_TIME_FORMAT = 'YYYY-MM-DD HH:mm:ss';
const DATE_TIME_DISPLAY_FORMAT = 'YYYY-MM-DD HH:mm';
const EXPORT_DATE_TIME_FORMAT = 'YYYYMMDDHHmmss';

const THEME_MAP = {
  light: {
    panelBg: 'linear-gradient(180deg, #f8fbff 0%, #eef4fb 100%)',
    panelBorder: '#d9e7f5',
    panelShadow: '0 10px 30px rgba(48, 86, 132, 0.10)',
    toolbarBg: 'rgba(255, 255, 255, 0.78)',
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
    toolbarBg: 'rgba(15, 23, 42, 0.82)',
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
    toolbarBg: 'rgba(240, 252, 255, 0.86)',
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
    toolbarBg: 'rgba(255, 250, 240, 0.86)',
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
    toolbarBg: 'rgba(244, 255, 249, 0.84)',
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
  name: 'ism-view-history-diy-difference',
  inject: ['getNode'],
  props: {},
  i18n: require('@/i18n/language'),
  components: {
    VNodes
  },
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
    }
  },
  computed: {
    tableScrollWidth() {
      return this.dynamicColumns.reduce((total, column) => total + (Number(column.width) || 120), 0);
    },
    currentTheme() {
      return THEME_MAP[this.selectedTheme] || THEME_MAP.light;
    },
    scrollbarThemeClass() {
      const id = this.detail && this.detail.identifier ? this.detail.identifier : 'default';
      return `history-scrollbar-${id}`;
    },
    overlayThemeClass() {
      const id = this.detail && this.detail.identifier ? this.detail.identifier : 'default';
      return `history-overlay-${id}`;
    },
    styleVar() {
      return {
        '--panelBg': this.currentTheme.panelBg,
        '--panelBorder': this.currentTheme.panelBorder,
        '--panelShadow': this.currentTheme.panelShadow,
        '--toolbarBg': this.currentTheme.toolbarBg,
        '--toolbarText': this.currentTheme.toolbarText,
        '--toolbarAccent': this.currentTheme.toolbarAccent,
        '--toolbarAccentSoft': this.currentTheme.toolbarAccentSoft,
        '--toolbarAccentBorder': this.currentTheme.toolbarAccentBorder,
        '--toolbarAccentText': this.currentTheme.toolbarAccentText,
        '--tableHeaderColor': this.tableHeaderColor || this.currentTheme.tableHeaderColor,
        '--tableHeaderBackColor': this.tableHeaderBackColor || this.currentTheme.tableHeaderBackColor,
        '--tableHeaderFont': this.tableHeaderFont || this.currentTheme.tableHeaderFont,
        '--tableHeaderFontSize': `${this.tableHeaderFontSize || 12}px`,
        '--tableSplitColor': this.tableSplitColor || this.currentTheme.tableSplitColor,
        '--tableColumnSplitColor': this.currentTheme.tableColumnSplitColor,
        '--tableHoverColor': this.tableHoverColor || this.currentTheme.tableHoverColor,
        '--tableRowOddBg': this.currentTheme.tableRowOddBg,
        '--tableRowEvenBg': this.currentTheme.tableRowEvenBg,
        '--SearchColor': this.SearchColor || this.currentTheme.searchColor,
        '--SearchBackColor': this.SearchBackColor || this.currentTheme.searchBackColor,
        '--SearchBorderColor': this.SearchBorderColor || this.currentTheme.searchBorderColor,
        '--scrollBgColor': this.scrollBgColor || this.currentTheme.scrollBgColor,
        '--scrollFrColor': this.scrollFrColor || this.currentTheme.scrollFrColor,
        '--scrollHdColor': this.scrollHdColor || this.currentTheme.scrollHdColor,
        '--fontFamily': this.fontFamily,
        '--fontSize': `${this.fontSize || 12}px`,
        '--foreColor': this.foreColor || this.currentTheme.foreColor,
        '--backColor': this.backColor || this.currentTheme.backColor,
        height: `calc(${this.height}px - 16px)`,
        overflow: 'hidden'
      };
    },
    pagedData() {
      const current = this.pagination && this.pagination.current ? this.pagination.current : 1;
      const pageSize = this.pagination && this.pagination.pageSize ? this.pagination.pageSize : 10;
      const start = (current - 1) * pageSize;
      const end = start + pageSize;
      const pageRows = this.reportRows.slice(start, end);
      return this.isLastPage ? pageRows.concat(this.summaryRows) : pageRows;
    },
    isLastPage() {
      const pageSize = this.pagination && this.pagination.pageSize ? this.pagination.pageSize : 10;
      const maxPage = Math.max(1, Math.ceil(this.reportRows.length / pageSize));
      const current = this.pagination && this.pagination.current ? this.pagination.current : 1;
      return this.reportRows.length > 0 && current >= maxPage;
    },
    paginationTotal() {
      return this.reportRows.length;
    },
    summaryRows() {
      if (!this.reportRows.length) {
        return [];
      }
      const fields = ['startValue', 'endValue', 'difference'];
      const summaryRow = {
        rowName: '__summary_total__',
        isSummaryRow: true,
        deviceName: '----',
        dataName: '汇总'
      };
      const avgRow = {
        rowName: '__summary_avg__',
        isSummaryRow: true,
        deviceName: '----',
        dataName: '平均值'
      };
      fields.forEach(field => {
        const values = this.reportRows
          .map(row => Number(row[field]))
          .filter(value => Number.isFinite(value));
        const sum = values.reduce((total, value) => total + value, 0);
        summaryRow[field] = sum.toFixed(2);
        avgRow[field] = values.length ? (sum / values.length).toFixed(2) : '0.00';
      });
      return [summaryRow, avgRow];
    },
    reportRows() {
      if (!this.rowNames.length) {
        return this.dataSource.map(row => this.normalizeReportRow(row));
      }
      return this.rowNames.map((rowName, index) => {
        const parsed = this.parseDeviceDataName(rowName);
        const matchedRow = this.findReportRow(parsed, index);
        return this.normalizeReportRow(matchedRow, {
          rowName,
          // 优先用接口返回的名称，否则用配置里的别名（中文）
          // 设备名称/数据名称始终用配置别名显示，不随接口数据变化
          deviceName: parsed.deviceDisplayName || parsed.deviceSN,
          dataName: parsed.dataDisplayName || parsed.dataSN
        });
      });
    },
    dynamicColumns() {
      return [
        {
          title: '序号',
          key: 'index',
          width: 60,
          align: 'center',
          customRender: (t, r, index) => {
            if (r && r.isSummaryRow) {
              return '';
            }
            const { current, pageSize } = this.pagination || { current: 1, pageSize: 10 };
            return (current - 1) * pageSize + index + 1;
          }
        },
        {
          title: '设备名称',
          dataIndex: 'deviceName',
          key: 'deviceName',
          width: 150,
          align: 'center'
        },
        {
          title: '数据名称',
          dataIndex: 'dataName',
          key: 'dataName',
          width: 150,
          align: 'center'
        },
        {
          title: '起始值',
          dataIndex: 'startValue',
          key: 'startValue',
          width: 120,
          align: 'center'
        },
        {
          title: '结束值',
          dataIndex: 'endValue',
          key: 'endValue',
          width: 120,
          align: 'center'
        },
        {
          title: '差值',
          dataIndex: 'difference',
          key: 'difference',
          width: 100,
          align: 'center'
        }
      ];
    },
    exportTableData() {
      return this.reportRows.map((row, index) => ({
        序号: index + 1,
        设备名称: row.deviceName,
        数据名称: row.dataName,
        起始值: row.startValue,
        结束值: row.endValue,
        差值: row.difference
      }));
    }
  },
  data() {
    return {
      moment,
      dateTimeDisplayFormat: DATE_TIME_DISPLAY_FORMAT,
      tableKey: 0,
      dataSource: [],
      exportName: '',
      loadExecl: null,
      isLoadExecl: false,
      hasQueriedData: false,
      json_fields_cn: {
        记录时间: {
          field: 'MainTime',
          callback: value => formatDate(new Date(value), 'yyyy-MM-dd hh:mm:ss')
        }
      },
      json_fields_en: {
        'Record Time': {
          field: 'MainTime',
          callback: value => formatDate(new Date(value), 'yyyy-MM-dd hh:mm:ss')
        }
      },
      json_fields: {},
      json_meta: [[{ key: 'charset', value: 'utf-8' }]],
      SelectDateRange: [moment().subtract(1, 'day'), moment()],
      messageShowLoad: false,
      advanced: true,
      selectedTheme: 'light',
      themeOptions: THEME_OPTIONS,
      pagination: {
        current: 1,
        pageSize: 10,
        total: 0,
        showSizeChanger: false,
        hideOnSinglePage: false,
        showLessItems: true,
        simple: true
      },
      headerHeight: 48,
      paginationHeight: 40,
      showPaginationBar: true,
      tableSplitColor: '#000',
      tableHoverColor: '#fff',
      tableHeaderColor: '',
      tableHeaderBackColor: '',
      tableHeaderFont: 'Arial',
      scrollBgColor: '#f0f0f0',
      scrollFrColor: '#c1c1c1',
      scrollHdColor: '#a8a8a8',
      tableHeaderFontSize: 17,
      SearchColor: '',
      SearchBackColor: '',
      SearchBorderColor: '',
      rowNames: [
        '空调-AC001-温度-TEMP',
        '灯光-LT002-亮度-LIGHT',
        '门禁-AC003-状态-DOOR',
        '监控-CAM004-图像-CAM'
      ],
      columnNames: [],
      cellData: [
        [100, 200, 300, 400],
        [150, 250, 350, 450],
        [120, 220, 320, 420],
        [180, 280, 380, 480]
      ],
      detail: null,
      IsToolBox: false,
      editMode: true,
      Url: '',
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
      isShowDocx: false,
      docx: '',
      AlarmTimer: null,
      scrollbarStyleTagId: '',
      fontFamily: 'Arial',
      fontSize: 14,
      backColor: '',
      foreColor: '',
      base: {
        text: '差值表',
        icon: 'icon-gongzuohuibao',
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
            position: { x: 0, y: 0, w: 520, h: 200 },
            backColor: 'transparent',
            foreColor: '#000000',
            fontWeight: 400,
            fontSize: 12,
            fontFamily: 'Arial',
            visible: 1,
            zIndex: -1,
            transform: 0,
            diy: [
              {
                name: 'configComponent.viewRealTable.deviceList',
                type: 9,
                value: '空调-AC001-温度-TEMP,灯光-LT002-亮度-LIGHT,门禁-AC003-状态-DOOR,监控-CAM004-图像-CAM',
                key: 'deviceList'
              },
              { name: 'configComponent.AlarmList.waitTime', type: 7, value: 1000, min: 100, max: 10000, key: 'waitTime' },
              { name: 'configComponent.DeviceTree.ShowCount', type: 1, value: 5, min: 1, max: 100, key: 'ShowCount' },
              { name: 'configComponent.DeviceTree.SearchColor', type: 2, value: '#000000', key: 'SearchColor' },
              { name: 'configComponent.DeviceTree.SearchBackColor', type: 2, value: '#ffffff', key: 'SearchBackColor' },
              { name: 'configComponent.DeviceTree.SearchBorderColor', type: 2, value: '#cbc6c6', key: 'SearchBorderColor' },
              { name: 'configComponent.DataHistoryList.tableHeaderColor', type: 2, value: '#000000', key: 'tableHeaderColor' },
              { name: 'configComponent.DataHistoryList.tableHeaderBackColor', type: 2, value: '#fafafa', key: 'tableHeaderBackColor' },
              { name: 'configComponent.viewRealTable.tableHeaderFont', type: 3, value: 'Arial', key: 'tableHeaderFont' },
              { name: 'configComponent.viewRealTable.tableHeaderFontSize', type: 1, value: 14, key: 'tableHeaderFontSize' },
              { name: 'configComponent.DataHistoryList.tableSplitColor', type: 2, value: '#ebedf0', key: 'tableSplitColor' },
              { name: 'configComponent.DataHistoryList.tableHoverColor', type: 2, value: '#ffffff', key: 'tableHoverColor' },
              {
                name: '主题风格',
                type: 6,
                value: 'light',
                enumList: THEME_OPTIONS.map(item => ({ value: item.value, option: item.label })),
                key: 'themeName'
              }
            ]
          }
        }
      }
    };
  },
  methods: {
    handlePageChange(page) {
      this.pagination = {
        ...this.pagination,
        current: page
      };
    },
    syncPagination(resetCurrent = false) {
      const total = this.reportRows.length;
      const pageSize = Number(this.pagination.pageSize) || 10;
      const maxPage = Math.max(1, Math.ceil(total / pageSize));
      const current = resetCurrent ? 1 : Math.min(this.pagination.current || 1, maxPage);
      this.pagination = {
        ...this.pagination,
        current,
        pageSize,
        total
      };
    },
    getTableBodyElement() {
      const wrapper = this.$refs.tableScroll;
      if (!wrapper) return null;
      return wrapper.querySelector('.ant-table-body') || wrapper;
    },
    getScrollTargets() {
      const wrapper = this.$refs.tableScroll;
      if (!wrapper) return [];
      const targets = [
        wrapper,
        wrapper.querySelector('.ant-table-body'),
        wrapper.querySelector('.ant-table-content'),
        wrapper.querySelector('.ant-table-body-outer'),
        wrapper.querySelector('.ant-table-body-inner')
      ].filter(Boolean);
      return Array.from(new Set(targets));
    },
    onDateRangeChange(dates, dateStrings) {
      this.SelectDateRange = dates;
      this.hasQueriedData = false;
    },
    startDownload() {
      this.isLoadExecl = true;
      this.loadExecl = this.$message.loading(this.$t('reporting.AlarmHistory.LoadingExecl'), 0);
    },
    finishDownload() {
      if (typeof this.loadExecl === 'function') this.loadExecl();
      this.loadExecl = null;
      this.isLoadExecl = false;
    },
    handleExport() {
      this.startDownload();
      exportExcelWithStyle(this.exportTableData, this.json_fields, this.exportName, this.formatExportQueryDate()).then(() => {
        this.finishDownload();
      }).catch(() => {
        this.finishDownload();
        this.$message.error(this.$t('loginPage.serverError'), 3);
      });
    },
    formatExportQueryDate() {
      if (!this.SelectDateRange || !this.SelectDateRange[0]) {
        return '';
      }
      const start = this.SelectDateRange[0].format(DATE_TIME_DISPLAY_FORMAT);
      const end = this.SelectDateRange[1] ? this.SelectDateRange[1].format(DATE_TIME_DISPLAY_FORMAT) : '';
      return [start, end];
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
        .${className} .ant-table-body::-webkit-scrollbar,
        .${className} .ant-table-content::-webkit-scrollbar,
        .${className} ::-webkit-scrollbar {
          width: 10px !important;
          height: 10px !important;
        }
        .${className} .ant-table-body::-webkit-scrollbar-track,
        .${className} .ant-table-content::-webkit-scrollbar-track,
        .${className} ::-webkit-scrollbar-track {
          background: ${theme.scrollBgColor} !important;
          border-radius: 999px !important;
          box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.18) !important;
        }
        .${className} .ant-table-body::-webkit-scrollbar-thumb,
        .${className} .ant-table-content::-webkit-scrollbar-thumb,
        .${className} ::-webkit-scrollbar-thumb {
          background: ${theme.scrollFrColor} !important;
          border-radius: 999px !important;
          border: 2px solid ${theme.scrollBgColor} !important;
          background-image: none !important;
        }
        .${className} .ant-table-body::-webkit-scrollbar-thumb:hover,
        .${className} .ant-table-content::-webkit-scrollbar-thumb:hover,
        .${className} ::-webkit-scrollbar-thumb:hover {
          background: ${theme.scrollHdColor} !important;
        }
      `;
      this.applyOverlayTheme();
    },
    applyOverlayTheme() {
      if (typeof document === 'undefined') {
        return;
      }
      const theme = this.currentTheme;
      const className = this.overlayThemeClass;
      const styleId = `overlay-style-${className}`;
      let styleTag = document.getElementById(styleId);
      if (!styleTag) {
        styleTag = document.createElement('style');
        styleTag.id = styleId;
        document.head.appendChild(styleTag);
      }
      styleTag.textContent = `
        .${className}.ant-calendar-picker-container,
        .${className} .ant-calendar,
        .${className}.ant-select-dropdown,
        .${className} .ant-select-dropdown-menu {
          background: ${theme.searchBackColor} !important;
          color: ${theme.foreColor} !important;
        }
        .${className} .ant-calendar {
          border: 1px solid ${theme.searchBorderColor} !important;
          box-shadow: 0 10px 28px ${theme.toolbarAccentSoft} !important;
        }
        .${className} .ant-calendar-input,
        .${className} .ant-calendar-header,
        .${className} .ant-calendar-footer,
        .${className} .ant-calendar-month-panel,
        .${className} .ant-calendar-year-panel,
        .${className} .ant-calendar-decade-panel {
          background: ${theme.searchBackColor} !important;
          color: ${theme.foreColor} !important;
          border-color: ${theme.searchBorderColor} !important;
        }
        .${className} .ant-calendar-column-header,
        .${className} .ant-calendar-week-number,
        .${className} .ant-calendar-date,
        .${className} .ant-calendar-month-panel-cell,
        .${className} .ant-calendar-year-panel-cell,
        .${className} .ant-calendar-decade-panel-cell,
        .${className} .ant-calendar-month-panel-header,
        .${className} .ant-calendar-year-panel-header,
        .${className} .ant-calendar-decade-panel-header,
        .${className} .ant-calendar-header a,
        .${className} .ant-calendar-next-month-btn,
        .${className} .ant-calendar-prev-month-btn,
        .${className} .ant-calendar-next-year-btn,
        .${className} .ant-calendar-prev-year-btn,
        .${className} .ant-calendar-next-decade-btn,
        .${className} .ant-calendar-prev-decade-btn,
        .${className} .ant-calendar-today-btn,
        .${className} .ant-select-dropdown-menu-item {
          color: ${theme.foreColor} !important;
        }
        .${className} .ant-calendar-disabled-cell .ant-calendar-date,
        .${className} .ant-calendar-last-month-cell .ant-calendar-date {
          color: ${theme.scrollFrColor} !important;
        }
        .${className} .ant-calendar-selected-day .ant-calendar-date,
        .${className} .ant-calendar-selected-date .ant-calendar-date,
        .${className} .ant-calendar-today .ant-calendar-date,
        .${className} .ant-calendar-month-panel-selected-cell .ant-calendar-month-panel-month,
        .${className} .ant-calendar-year-panel-selected-cell .ant-calendar-year-panel-year,
        .${className} .ant-calendar-decade-panel-selected-cell .ant-calendar-decade-panel-decade,
        .${className} .ant-select-dropdown-menu-item-selected {
          background: ${theme.toolbarAccentSoft} !important;
          border-color: ${theme.toolbarAccentBorder} !important;
          color: ${theme.toolbarAccentText} !important;
        }
        .${className} .ant-calendar-date:hover,
        .${className} .ant-calendar-month-panel-month:hover,
        .${className} .ant-calendar-year-panel-year:hover,
        .${className} .ant-calendar-decade-panel-decade:hover,
        .${className} .ant-select-dropdown-menu-item:hover {
          background: ${theme.tableHoverColor} !important;
        }
      `;
    },
    batchUpdateConfig(newCellData) {
      this.cellData = newCellData;
      if (this.pagination) {
        this.pagination = {
          ...this.pagination,
          current: 1,
          total: newCellData ? newCellData.length : 0
        };
      }
      this.$nextTick(() => {
        this.tableKey += 1;
      });
    },
    normalizeReportName(value) {
      return String(value || '').trim().toLowerCase();
    },
    pickDisplayValue(value, fallback) {
      return value === undefined || value === null || value === '' ? fallback : value;
    },
    formatNumberValue(value, fallback) {
      const displayValue = this.pickDisplayValue(value, fallback);
      if (displayValue === fallback) {
        return displayValue;
      }
      const numberValue = Number(displayValue);
      return Number.isFinite(numberValue) ? numberValue.toFixed(2) : displayValue;
    },
    getReportValue(row, keys) {
      if (!row) {
        return undefined;
      }
      const matchedKey = keys.find(key => row[key] !== undefined && row[key] !== null && row[key] !== '');
      return matchedKey ? row[matchedKey] : undefined;
    },
    normalizeReportRow(row, fallback = {}) {
      return {
        ...fallback,
        deviceName: fallback.deviceName || this.getReportValue(row, ['deviceName', 'DeviceName', 'deviceRealName', 'deviceAlias']) || '',
        dataName: fallback.dataName || this.getReportValue(row, ['dataName', 'DataName', 'dataRealName', 'dataAlias']) || '',
        startValue: this.formatNumberValue(this.getReportValue(row, ['startValue', 'StartValue', 'start', 'Start']), '-'),
        endValue: this.formatNumberValue(this.getReportValue(row, ['endValue', 'EndValue', 'end', 'End']), '-'),
        difference: this.formatNumberValue(this.getReportValue(row, ['difference', 'Difference', 'diff', 'Diff']), '-')
      };
    },
    parseDeviceDataName(fullName) {
      // 格式：设备别名-设备编号-数据别名-数据编号（如 空调-AC001-温度-TEMP）
      const parts = String(fullName || '').split('-').map(value => value.trim());
      const deviceDisplayName = parts[0] || '';  // 显示用别名（中文）
      const deviceSN = parts[1] || parts[0] || '';  // 实际设备编号（发请求用）
      const dataDisplayName = parts[2] || '';  // 显示用数据别名（中文）
      const dataSN = parts[3] || parts[2] || '';  // 实际数据编号（发请求用）
      return {
        raw: fullName,
        deviceDisplayName,
        deviceSN,
        dataDisplayName,
        dataSN,
        // 匹配接口返回数据时的候选集（编号+别名都试）
        deviceCandidates: [deviceSN, deviceDisplayName].filter(Boolean),
        dataCandidates: [dataSN, dataDisplayName].filter(Boolean)
      };
    },
    isReportRowMatch(row, parsed) {
      if (!row) {
        return false;
      }
      const rowDeviceName = this.normalizeReportName(row.deviceName || row.DeviceName || row.deviceRealName || row.deviceAlias);
      const rowDataName = this.normalizeReportName(row.dataName || row.DataName || row.dataRealName || row.dataAlias);
      // 两个字段都必须非空才做匹配，避免空字段导致误匹配
      if (!rowDeviceName || !rowDataName) {
        return false;
      }
      const deviceCandidates = parsed.deviceCandidates.map(this.normalizeReportName);
      const dataCandidates = parsed.dataCandidates.map(this.normalizeReportName);
      return deviceCandidates.includes(rowDeviceName) && dataCandidates.includes(rowDataName);
    },
    findReportRow(parsed, index) {
      const matchedByName = this.dataSource.find(row => this.isReportRowMatch(row, parsed));
      if (matchedByName) {
        return matchedByName;
      }
      const rowByIndex = this.dataSource[index];
      if (rowByIndex && !rowByIndex.deviceName && !rowByIndex.dataName) {
        return rowByIndex;
      }
      return null;
    },
    getRequestPairs() {
      const pairs = [];
      this.rowNames.forEach(rowName => {
        const parsed = this.parseDeviceDataName(rowName);
        // 只用设备编号(deviceSN)和数据编号(dataSN)发请求，别名不发
        if (parsed.deviceSN && parsed.dataSN) {
          const key = `${parsed.deviceSN}::${parsed.dataSN}`;
          if (!pairs.some(item => item.key === key)) {
            pairs.push({ key, deviceName: parsed.deviceSN, dataName: parsed.dataSN });
          }
        }
      });
      return pairs;
    },
    QueryHistoryDataDifferent() {
      const startTime = this.SelectDateRange[0] ? this.SelectDateRange[0].format(DATE_TIME_FORMAT) : moment().subtract(1, 'day').format(DATE_TIME_FORMAT);
      const endTime = this.SelectDateRange[1] ? this.SelectDateRange[1].format(DATE_TIME_FORMAT) : moment().format(DATE_TIME_FORMAT);
      
      if (new Date(startTime) > new Date(endTime)) {
        this.$message.warning('开始时间不能大于结束时间', 3);
        return;
      }
      
      this.$message.info(`正在查询...`, 2);
      this.messageShowLoad = true;
      GetHistoryDiyDifferenceReport(this.buildRequestParams())
        .then(res => {
          if (res.data.code === 0||res.data.code === -1) {
            this.dataSource = Array.isArray(res.data.realData) ? res.data.realData : [];
            this.hasQueriedData = true;
            this.pagination = {
              ...this.pagination,
              current: 1
            };
            this.syncPagination(true);
            this.tableKey += 1;
            const dateRangeStr = this.formatDateRange();
            this.updateExportData('差值报表_', dateRangeStr);
          } else {
            this.dataSource = [];
            this.hasQueriedData = false;
            this.$message.warning('查询失败，请检查参数或联系管理员', 3);
          }
          this.messageShowLoad = false;
        })
        .catch(() => {
          this.messageShowLoad = false;
          this.hasQueriedData = false;
          this.$message.error(this.$t('loginPage.serverError'), 3);
        });
    },
    buildRequestParams() {
      const startTime = this.SelectDateRange[0] ? this.SelectDateRange[0].format(DATE_TIME_FORMAT) : moment().subtract(1, 'day').format(DATE_TIME_FORMAT);
      const endTime = this.SelectDateRange[1] ? this.SelectDateRange[1].format(DATE_TIME_FORMAT) : moment().format(DATE_TIME_FORMAT);
      const requestPairs = this.getRequestPairs();
      return {
        QueryDate: [startTime, endTime],
        DeviceList: requestPairs.map(item => item.deviceName),
        DataList: requestPairs.map(item => item.dataName)
      };
    },
    formatDateRange() {
      if (!this.SelectDateRange || !this.SelectDateRange[0]) return '';
      const start = this.SelectDateRange[0].format(EXPORT_DATE_TIME_FORMAT);
      const end = this.SelectDateRange[1] ? this.SelectDateRange[1].format(EXPORT_DATE_TIME_FORMAT) : '';
      return `${start}_${end}`;
    },
    updateExportData(prefix, exportDate) {
      this.json_fields = {
        序号: '序号',
        设备名称: '设备名称',
        数据名称: '数据名称',
        起始值: '起始值',
        结束值: '结束值',
        差值: '差值'
      };
      this.exportName = `${prefix}${String(exportDate).replace(/-/g, '')}`;
    },
    initComponents(option) {
      if (!option || !option.style) {
        return;
      }
      this.isShowDocx = false;
      this.width = option.style.position.w;
      this.height = option.style.position.h;
      this.foreColor = option.style.foreColor;
      this.backColor = option.style.backColor;
      this.fontSize = option.style.fontSize;
      this.fontFamily = option.style.fontFamily;
      let deviceListChanged = false;
      for (let i = 0; i < option.style.diy.length; i += 1) {
        const item = option.style.diy[i];
        if (item.key === 'deviceList' && item.value) {
          const nextRowNames = item.value.split(',').map(v => v.trim()).filter(Boolean);
          deviceListChanged = nextRowNames.join(',') !== this.rowNames.join(',');
          this.rowNames = nextRowNames;
        } else if (item.key === 'tableHeaderColor') {
          this.tableHeaderColor = item.value;
        } else if (item.key === 'tableHeaderBackColor') {
          this.tableHeaderBackColor = item.value;
        } else if (item.key === 'tableSplitColor') {
          this.tableSplitColor = item.value;
        } else if (item.key === 'tableHoverColor') {
          this.tableHoverColor = item.value;
        } else if (item.key === 'ShowCount') {
          this.pagination.pageSize = parseInt(item.value, 10);
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
      if (deviceListChanged) {
        this.dataSource = [];
        this.hasQueriedData = false;
        this.syncPagination(true);
        this.tableKey += 1;
      }
      this.syncPagination();
      this.animateType = option.animate.selected;
      this.isStart = !option.animate.isExpression;
      this.$nextTick(() => {
        this.applyScrollbarTheme();
      });
    }
  },
  beforeDestroy() {
    if (this._tableScrollTargets && this._boundTableBodyScroll) {
      this._tableScrollTargets.forEach((target) => {
        target.removeEventListener('scroll', this._boundTableBodyScroll, true);
      });
    }
    clearInterval(this.AlarmTimer);
    if (typeof document !== 'undefined' && this.scrollbarStyleTagId) {
      const styleTag = document.getElementById(this.scrollbarStyleTagId);
      if (styleTag) {
        styleTag.remove();
      }
    }
    if (typeof document !== 'undefined') {
      const overlayStyleTag = document.getElementById(`overlay-style-${this.overlayThemeClass}`);
      if (overlayStyleTag) {
        overlayStyleTag.remove();
      }
    }
  },
  mounted() {
    this.$nextTick(() => {
      this.initComponents(this.detail);
      this.applyScrollbarTheme();
      const id = this.detail && this.detail.identifier ? this.detail.identifier : '';
      if (!id) return;
      const activeEvent = `${id}activeEvent`;
      const animateEvent = `${id}animateEvent`;
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
    this.initComponents(this.detail);
  }
};
</script>

<style lang="less" scoped>
.history-theme-shell {
  display: flex;
  flex-direction: column;
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

.toolbar-row {
  flex: 0 0 auto;
  display: flex;
  align-items: center;
  justify-content: flex-start;
  gap: 10px;
  min-height: 48px;
  padding: 8px 12px;
  margin-bottom: 8px;
  background: var(--tableHeaderBackColor);
  border: 1px solid var(--panelBorder);
  border-radius: 10px;
  backdrop-filter: blur(6px);
}

.report-title {
  flex: 0 0 auto;
  font-size: 13px;
  font-weight: 600;
  color: var(--tableHeaderColor);
  letter-spacing: 0.5px;
  white-space: nowrap;
  margin-right: 4px;
}

.date-picker {
  flex: 0 1 50%;
  width: 50%;
  min-width: 320px;
  max-width: 50%;
}

.toolbar-actions {
  flex: 0 0 auto;
  display: flex;
  align-items: center;
  justify-content: flex-start;
  gap: 10px;
}

.btn-query {
  flex: 0 0 auto;
}

.btn-export {
  flex: 0 0 auto;
}

@media (max-width: 760px) {
  .toolbar-row {
    align-items: stretch;
    flex-wrap: wrap;
  }

  .date-picker,
  .toolbar-actions {
    flex-basis: 100%;
    width: 100%;
    max-width: none;
    min-width: 0;
  }
}

.table-container {
  display: flex;
  flex-direction: column;
  flex: 1 1 auto;
  min-height: 0;
  width: 100%;
  overflow: hidden;
  box-sizing: border-box;
  border: 1px solid var(--panelBorder);
  border-radius: 10px;
  background: var(--panelBg);
}

.table-scroll {
  flex: 1 1 auto;
  min-height: 0;
  overflow: auto;
  box-sizing: border-box;
  padding-bottom: 20px;
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
  width: fit-content;
  border: none;
  border-radius: 0;
  background: transparent;
  backdrop-filter: none;
  box-shadow: none;
  opacity: 0;
  transform: translateY(4px);
  pointer-events: none;
  overflow: hidden;
  transition: opacity 0.18s ease, transform 0.18s ease;
}

.table-pagination-bar.visible {
  min-height: 18px;
  height: 24px;
  padding: 0 8px;
  margin: 0 auto 2px;
  width: fit-content;
  border-width: 1px;
  opacity: 1;
  transform: translateY(0);
  pointer-events: auto;
}

::v-deep .ant-table-pagination.ant-pagination {
  margin: 0 !important;
  float: none !important;
  text-align: left !important;
  padding: 0 !important;
  line-height: 1 !important;
}

::v-deep .table-pagination-bar .ant-pagination-simple {
  font-size: 11px !important;
}

::v-deep .table-pagination-bar .ant-pagination-prev,
::v-deep .table-pagination-bar .ant-pagination-next,
::v-deep .table-pagination-bar .ant-pagination-simple-pager,
::v-deep .table-pagination-bar .ant-pagination-item-link {
  min-width: 18px !important;
  height: 18px !important;
  line-height: 16px !important;
}

::v-deep .table-pagination-bar .ant-pagination-simple-pager input {
  height: 16px !important;
  min-width: 36px !important;
  padding: 0 4px !important;
}

::v-deep .ant-form-item-label > label,
::v-deep .ant-radio-wrapper,
::v-deep .ant-select,
::v-deep .ant-calendar-picker,
::v-deep .ant-btn {
  color: var(--toolbarText);
}

::v-deep .ant-radio-button-wrapper {
  color: var(--toolbarText);
  background: var(--SearchBackColor);
  border-color: var(--SearchBorderColor);
}

::v-deep .ant-radio-button-wrapper-checked {
  color: #fff;
  background: var(--toolbarAccent);
  border-color: var(--toolbarAccent);
  box-shadow: none;
}

::v-deep .ant-input,
::v-deep .ant-select-selection,
::v-deep .ant-calendar-picker-input.ant-input {
  color: var(--SearchColor);
  background: var(--SearchBackColor);
  border-color: var(--SearchBorderColor);
  border-radius: 10px;
}

::v-deep .ant-calendar-picker {
  display: block;
}

::v-deep .ant-calendar-picker-icon,
::v-deep .ant-calendar-picker-clear {
  color: var(--toolbarAccentText) !important;
  background: var(--toolbarAccentSoft) !important;
  border-radius: 6px;
}

::v-deep .ant-calendar-picker-clear:hover,
::v-deep .ant-calendar-picker-icon:hover {
  color: var(--toolbarAccentText) !important;
  background: var(--tableHoverColor) !important;
}

::v-deep .ant-btn-primary {
  background: var(--toolbarAccent);
  border-color: var(--toolbarAccent);
  box-shadow: 0 2px 8px var(--toolbarAccentSoft);
}

::v-deep .btn-query.ant-btn-primary {
  padding: 0 16px;
  height: 32px;
  font-size: 12px;
  border-radius: 7px;
}

::v-deep .btn-export.ant-btn:not(.ant-btn-primary) {
  padding: 0 14px;
  height: 32px;
  font-size: 12px;
  border-radius: 7px;
}

::v-deep .ant-btn:not(.ant-btn-primary) {
  color: var(--toolbarAccentText);
  background: var(--toolbarAccentSoft);
  border-color: var(--toolbarAccentBorder);
}

::v-deep .ant-btn:not(.ant-btn-primary):hover,
::v-deep .ant-btn:not(.ant-btn-primary):focus {
  color: var(--toolbarAccentText);
  background: var(--toolbarAccentSoft);
  border-color: var(--toolbarAccent);
}

::v-deep .toolbar-row .ant-calendar-picker-input.ant-input {
  height: 32px;
  font-size: 12px;
  padding: 0 10px;
}

::v-deep .ant-table {
  height: 100%;
  color: var(--foreColor);
  background: transparent;
  font-size: 12px;
}

::v-deep .ant-table-wrapper,
::v-deep .ant-spin-nested-loading,
::v-deep .ant-spin-container {
  height: auto;
}

::v-deep .ant-table-thead > tr > th {
  position: sticky !important;
  top: 0;
  z-index: 6;
  color: var(--tableHeaderColor) !important;
  font-size: 12px !important;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'PingFang SC', 'Microsoft YaHei', sans-serif !important;
  font-weight: 600 !important;
  background: var(--tableHeaderBackColor) !important;
  border-bottom: 2px solid var(--tableSplitColor) !important;
  border-right: none !important;
  white-space: nowrap !important;
  padding: 8px 10px !important;
}

::v-deep .ant-table-thead > tr > th:last-child {
  border-right: none !important;
}

::v-deep .ant-table-tbody > tr > td {
  color: var(--foreColor) !important;
  font-size: 12px !important;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'PingFang SC', 'Microsoft YaHei', sans-serif !important;
  background: transparent !important;
  border-bottom: 1px solid var(--tableSplitColor) !important;
  border-right: none !important;
  white-space: nowrap !important;
  padding: 6px 10px !important;
  transition: background 0.15s ease;
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

::v-deep .ant-table-tbody > tr.summary-row > td {
  font-weight: 600 !important;
  background: var(--tableHeaderBackColor) !important;
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
  left: 210px;
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
