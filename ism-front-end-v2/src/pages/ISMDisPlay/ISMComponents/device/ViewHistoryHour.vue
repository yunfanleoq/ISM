<template>
  <div :class="['history-theme-shell', scrollbarThemeClass]" :style="styleVar">
    <div class="history-toolbar">
      <a-form layout="horizontal">
        <div :class="advanced ? null : 'fold'">
          <a-row>
            <a-col :md="7" :sm="24">
              <a-form-item :labelCol="{ span: 5 }" :wrapperCol="{ span: 18, offset: 1 }">
                <a-date-picker
                  :defaultValue="moment()"
                  style="width: 100%"
                  @change="onDateChange"
                  size="default"
                  :dropdownClassName="overlayThemeClass"
                  :placeholder="$t('reporting.AlarmHistory.DateDayType')"
                />
              </a-form-item>
            </a-col>

            <a-col :md="3" :sm="24">
              <span style="display: inline-block; margin-top: 3px;">
                <a-button :disabled="messageShowLoad" type="primary" @click="QueryHistoryDataHour">
                  {{ $t('reporting.AlarmHistory.Query') }}
                </a-button>
              </span>
            </a-col>

            <a-col :md="2" :sm="24">
              <span style="display: inline-block; margin-top: 3px;">
                <a-button :disabled="isLoadExecl || !dataSource.length" type="default" style="margin-left: 5px" @click="handleExport">
                  {{ $t('reporting.AlarmHistory.Export') }}
                </a-button>
              </span>
            </a-col>
          </a-row>
        </div>
      </a-form>
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
          :rowKey="(record, index) => `${record.rowName}-${index}`"
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
import { GetHistoryHour } from '@/services/report';
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
  name: 'ism-view-history-hour',
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
    dynamicData() {
      const data = [];
      this.rowNames.forEach((fullRowName, rowIndex) => {
        const rowData = { rowName: fullRowName };
        this.columnNames.forEach((_, colIndex) => {
          const rowCellData = this.cellData[rowIndex] || [];
          rowData[`col_${colIndex}`] = rowCellData[colIndex] ?? '-';
        });
        data.push(rowData);
      });
      return data;
    },
    dynamicColumns() {
      const columns = [
        {
          title: '序号',
          key: 'index',
          width: 60,
          align: 'center',
          className: 'sticky-col sticky-col-0',
          customHeaderCell: () => ({ class: 'sticky-col sticky-col-0' }),
          customRender: (t, r, index) => {
            const { current = 1, pageSize = 10 } = this.pagination || {};
            return ((current || 1) - 1) * pageSize + index + 1;
          }
        },
        {
          title: '设备名称',
          key: 'alias',
          width: 120,
          align: 'center',
          className: 'sticky-col sticky-col-1',
          customHeaderCell: () => ({ class: 'sticky-col sticky-col-1' }),
          customRender: (text, record) => this.getAlias(record.rowName)
        },
        {
          title: '数据名称',
          key: 'realName',
          width: 120,
          align: 'center',
          className: 'sticky-col sticky-col-2',
          customHeaderCell: () => ({ class: 'sticky-col sticky-col-2' }),
          customRender: (text, record) => this.getRealName(record.rowName)
        }
      ];

      this.columnNames.forEach((colName, colIndex) => {
        columns.push({
          title: colName,
          dataIndex: `col_${colIndex}`,
          key: `col_${colIndex}`,
          width: 100,
          align: 'center'
        });
      });

      return columns;
    },
    pagedData() {
      const current = this.pagination && this.pagination.current ? this.pagination.current : 1;
      const pageSize = this.pagination && this.pagination.pageSize ? this.pagination.pageSize : 10;
      const start = (current - 1) * pageSize;
      const end = start + pageSize;
      return this.dynamicData.slice(start, end);
    },
    paginationTotal() {
      const explicitTotal = this.pagination && Number(this.pagination.total);
      if (explicitTotal && explicitTotal > 0) {
        return explicitTotal;
      }
      return this.dynamicData.length;
    },
    exportTableData() {
      return this.dynamicData.map((row, index) => {
        return {
          序号: index + 1,
          设备名称: this.getAlias(row.rowName),
          数据名称: this.getRealName(row.rowName),
          ...this.columnNames.reduce((acc, colName, idx) => {
            acc[colName] = row[`col_${idx}`];
            return acc;
          }, {})
        };
      });
    }
  },
  data() {
    return {
      moment,
      dataSource: [],
      exportName: '',
      loadExecl: null,
      isLoadExecl: false,
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
      SelectDateRange: moment().format('YYYY-MM-DD'),
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
      tableKey: 0,
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
      columnNames: Array.from({ length: 24 }, (_, i) => `${i}时`),
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
        text: 'configComponent.viewHistroyHour.title',
        icon: 'icon-yinhuanpaichafenxibiao',
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
    onDateChange(date, dateString) {
      this.SelectDateRange = dateString;
    },
    startDownload() {
      this.isLoadExecl = true;
      this.loadExecl = this.$message.loading(this.$t('reporting.AlarmHistory.LoadingExecl'), 0);
    },
    finishDownload() {
      this.$message.destroy(this.loadExecl);
      this.isLoadExecl = false;
    },
    handleExport() {
      this.startDownload();
      exportExcelWithStyle(this.dataSource, this.json_fields, this.exportName, this.SelectDateRange).then(() => {
        this.finishDownload();
      }).catch(() => {
        this.finishDownload();
        this.$message.error(this.$t('loginPage.serverError'), 3);
      });
    },
    handlePageChange(page) {
      this.pagination = {
        ...this.pagination,
        current: page
      };
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
        .${className} .ant-calendar-last-month-cell .ant-calendar-date,
        .${className} .ant-calendar-next-month-btn,
        .${className} .ant-calendar-prev-month-btn {
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
    getRealName(fullName) {
      const arr = (fullName || '').split('-');
      return arr[2] ? arr[2].trim() : '-';
    },
    getAlias(fullName) {
      const arr = (fullName || '').split('-');
      return arr[0] ? arr[0].trim() : '-';
    },
    QueryHistoryDataHour() {
      const params = {
        QueryDate: this.SelectDateRange,
        DeviceList: this.rowNames.map(rowName => {
          const parts = rowName.split('-');
          return parts[1] ? parts[1].trim() : '';
        }),
        DataList: this.rowNames.map(rowName => {
          const parts = rowName.split('-');
          return parts[3] ? parts[3].trim() : '';
        })
      };
      this.messageShowLoad = true;
      GetHistoryHour(params)
        .then(res => {
          if (res.data.code === 0) {
            this.batchUpdateConfig(res.data.realData);
            this.dataSource = this.exportTableData;
            this.json_fields = {
              设备名称: '设备名称',
              数据名称: '数据名称'
            };
            this.columnNames.forEach(h => {
              this.json_fields[h] = h;
            });
            this.exportName = `历史小时报表_${this.SelectDateRange.replace(/-/g, '')}`;
          }
          this.messageShowLoad = false;
        })
        .catch(() => {
          this.messageShowLoad = false;
          this.$message.error(this.$t('loginPage.serverError'), 3);
        });
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
      for (let i = 0; i < option.style.diy.length; i += 1) {
        const item = option.style.diy[i];
        if (item.key === 'deviceList' && item.value) {
          this.rowNames = item.value.split(',').map(v => v.trim()).filter(Boolean);
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
      this.animateType = option.animate.selected;
      this.isStart = !option.animate.isExpression;
      this.$nextTick(() => {
        this.applyScrollbarTheme();
      });
    }
  },
  beforeDestroy() {
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
    if (!this.editMode) {
      this.QueryHistoryDataHour();
    }
    this.$nextTick(() => {
      this.initComponents(this.detail);
      this.applyScrollbarTheme();
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

.history-toolbar {
  padding: 12px 14px 0;
  margin-bottom: 12px;
  background: var(--toolbarBg);
  border: 1px solid var(--panelBorder);
  border-radius: 14px;
  backdrop-filter: blur(8px);
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
  border-radius: 12px;
  background: #fff;
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
::v-deep .ant-select,
::v-deep .ant-calendar-picker,
::v-deep .ant-btn {
  color: var(--toolbarText);
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
  box-shadow: 0 8px 18px var(--toolbarAccentSoft);
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

::v-deep .ant-table {
  height: 100%;
  color: var(--foreColor);
  background: transparent;
}

::v-deep .ant-table-wrapper,
::v-deep .ant-spin-nested-loading,
::v-deep .ant-spin-container {
  height: auto;
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

::v-deep .ant-table-fixed-left,
::v-deep .ant-table-fixed-right {
  bottom: 14px !important;
  height: calc(100% - 14px) !important;
  overflow: hidden !important;
}

::v-deep .ant-table-fixed-left .ant-table-body-outer,
::v-deep .ant-table-fixed-right .ant-table-body-outer,
::v-deep .ant-table-fixed-left .ant-table-body-inner,
::v-deep .ant-table-fixed-right .ant-table-body-inner {
  margin-bottom: 14px !important;
  padding-bottom: 0 !important;
  background: transparent !important;
}

::v-deep .ant-table-fixed-left::after,
::v-deep .ant-table-fixed-right::after {
  display: none !important;
}

::v-deep .ant-table-body {
  overflow: visible !important;
  padding-bottom: 0;
}

::v-deep .ant-table-thead > tr > th {
  color: var(--tableHeaderColor) !important;
  font-size: var(--tableHeaderFontSize) !important;
  font-family: var(--tableHeaderFont) !important;
  background: var(--tableHeaderBackColor) !important;
  border-bottom: 1px solid var(--tableSplitColor) !important;
  border-right: 1px solid var(--tableColumnSplitColor) !important;
  white-space: nowrap !important;
}

::v-deep .ant-table-tbody > tr > td {
  color: var(--foreColor) !important;
  font-size: var(--fontSize) !important;
  font-family: var(--fontFamily) !important;
  background: var(--backColor) !important;
  border-bottom: 1px solid var(--tableSplitColor) !important;
  border-right: 1px solid var(--tableColumnSplitColor) !important;
  white-space: nowrap !important;
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