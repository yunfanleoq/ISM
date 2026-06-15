<template>
  <!-- 完全保留原有模板，不做任何修改 -->
  <a-card>
    <div>
      <a-form layout="horizontal">
        <div :class="advanced ? null : 'fold'">
          <a-row>
            <a-col :md="8" :sm="24">
              <a-form-item
                :label="$t('reporting.AlarmHistory.DeviceList')"
                :labelCol="{span: 5}"
                :wrapperCol="{span: 18, offset: 1}"
              >
                <a-tree-select
                  show-search
                  tree-node-filter-prop="title"
                  v-model="SelectDevice"
                  @change="SelectTreeDevice"
                  style="width: 100%"
                  tree-checkable
                  allow-clear
                  :dropdown-style="{ maxHeight: '400px', overflow: 'auto' }"
                  :tree-data="deviceTreeData"
                  :replace-fields="{ value: 'key', title: 'text' }"
                  :placeholder="$t('reporting.AlarmHistory.DeviceList')"
                  tree-default-expand-all
                >
                </a-tree-select>
              </a-form-item>
            </a-col>
            <a-col :md="8" :sm="24">
              <a-form-item
                :label="$t('reporting.AlarmHistory.DataList')"
                :labelCol="{span: 5}"
                :wrapperCol="{span: 18, offset: 1}"
              >
                <a-select
                  @dropdownVisibleChange="GetDeviceModelDataList"
                  @popupScroll="handlePopupScroll"
                  @search="handleSearch"
                  allowClear
                  show-search
                  optionFilterProp="children"
                  mode="multiple"
                  style="width: 100%"
                  :token-separators="[',']"
                  v-model="SelectAlarmData"
                >
                  <a-select-option v-for="(alarmItem, itemIndex) in frontDataZ" :key="itemIndex" :value="alarmItem.uuid">
                    {{ $t(alarmItem.name) }}
                  </a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :md="8" :sm="24">
              <a-form-item
                :label="$t('reporting.AlarmHistory.DateType')"
                :labelCol="{span: 5}"
                :wrapperCol="{span: 18, offset: 1}"
              >
                <a-radio-group v-model="SelectDateType" @change="chargeDateType">
                  <a-radio-button value="Day">
                    {{ $t('reporting.AlarmHistory.DateDayType') }}
                  </a-radio-button>
                  <a-radio-button value="Weekly">
                    {{ $t('reporting.AlarmHistory.DateWeeklyType') }}
                  </a-radio-button>
                  <a-radio-button value="Month">
                    {{ $t('reporting.AlarmHistory.DateMonthType') }}
                  </a-radio-button>
                  <a-radio-button value="Diy">
                    {{ $t('reporting.AlarmHistory.DateDiyType') }}
                  </a-radio-button>
                </a-radio-group>
              </a-form-item>
            </a-col>
          </a-row>
          <a-row>
            <a-col :md="8" :sm="24">
              <a-form-item
                :label="$t('reporting.AlarmHistory.SelectDate')"
                :labelCol="{span: 5}"
                :wrapperCol="{span: 18, offset: 1}"
              >
                <a-date-picker
                  :defaultValue="moment()"
                  style="width: 100%"
                  @change="onDateChange"
                  size="default"
                  :placeholder="$t('reporting.AlarmHistory.DateDayType')"
                  v-if="SelectDateType=='Day'"
                />
                <a-month-picker
                  :defaultValue="moment()"
                  style="width: 100%"
                  @change="onDateChange"
                  size="default"
                  :placeholder="$t('reporting.AlarmHistory.DateMonthType')"
                  v-if="SelectDateType=='Month'"
                />
                <a-week-picker
                  :defaultValue="moment()"
                  style="width: 100%"
                  @change="onWeeklyDateChange"
                  size="default"
                  :placeholder="$t('reporting.AlarmHistory.DateWeeklyType')"
                  v-if="SelectDateType=='Weekly'"
                />
                <a-range-picker
                  :default-value="[moment().add(-1, 'day'), moment()]"
                  :showTime="true"
                  @change="onDateChange"
                  size="default"
                  v-if="SelectDateType=='Diy'"
                />
              </a-form-item>
            </a-col>
            <a-col :md="getWindowColMd" :sm="24" style="margin-left:10px">
              <a-form-item
                  :label="$t('reporting.AlarmHistory.TimeWindows')"
                  :labelCol="{span: 6}"
                  :wrapperCol="{span: 10, offset: 1}"
              >
              <a-input v-model="dataWindows"></a-input>
              </a-form-item>
            </a-col>
            <a-col :md="3" :sm="24" v-if="SelectDateType=='Diy'">
              <span style="float: right; margin-top: 3px;">
                <a-button :disabled="messageShowLoad" type="primary" @click="QueryHistoryDataList">{{ $t('reporting.AlarmHistory.Query') }}</a-button>
              </span>
            </a-col>
            <a-col v-else :md="2" :sm="24">
              <span style="float: right; margin-top: 3px;">
                <a-button :disabled="messageShowLoad" type="primary" @click="QueryHistoryDataList">{{ $t('reporting.AlarmHistory.Query') }}</a-button>
              </span>
            </a-col>
            <a-col :md="2" :sm="24">
              <span style="float: left; margin-top: 3px;">
                <a-button :disabled="isLoadExecl || !dataSource.length" type="default" style="margin-left: 5px" @click="handleExport">{{ $t('reporting.AlarmHistory.Export') }}</a-button>
              </span>
            </a-col>
          </a-row>
        </div>
      </a-form>
    </div>
    <a-spin style="padding: 1px;" :spinning="messageShowLoad" tip="Loading...">
      <div v-if="!dataSource.length && !messageShowLoad" style="text-align: center; padding: 40px; color: #999;">
        {{ $t('reporting.AlarmHistory.NoData') }}
      </div>
      <a-table
        v-else
        :pagination="tablePagination"
        :columns="columns"
        :data-source="dataSource"
        :rowKey="record => record.key"
        :scroll="{ x: 'max-content',scrollToFirstRowOnChange:true }"
        bo
      >
        <template v-for="(item, index) in columns" :slot="item.slotName" slot-scope="text">
          <span :key="index">{{ $t(item.slotName) }}</span>
        </template>
        <span slot="MainTime" slot-scope="text">
          {{ text | formatDate }}
        </span>
      </a-table>
    </a-spin>
  </a-card>
</template>

<script>
// 完全保留原有导入，不做任何修改
import { getMonitorTree } from "@/services/device";
import { GetDeviceModelDataList } from "../../../services/device";
import { GetDataHistoryList, GetDataHistoryReport } from "../../../services/report";
import moment from 'moment';
import { formatDate } from '@/utils/common';
import {exportExcel} from '@/services/excelExport';
import 'moment/locale/zh-cn';
import 'moment/locale/en-ie';
import 'moment/locale/zh-tw';

export default {
  name: 'DataHistoryQueryReport',
  i18n: require('../../../i18n/language'),
  components: {
  },
  props: {
    tablePagination: {
      type: Object,
      default:{
        pageSize: 15,
        showSizeChanger: true,
        showTotal: (total) => `${this.$t('reporting.AlarmHistory.Total')} ${total} ${this.$t('reporting.AlarmHistory.Records')}`
      },
      required: true
    },
  },
  computed: {
    /**
     * 动态计算 a-col 的 md 宽度（基于日期类型）
     * Diy 类型时宽度变小，其他类型宽度正常
     */
    getWindowColMd() {
      return this.SelectDateType === 'Diy' ? 10 : 8; // Diy 时占 6 格，其他占 8 格
    },
    /**
     * 动态计算输入框宽度（基于日期类型）
     * 适配不同屏幕和日期类型的布局
     */
    getWindowInputWidth() {
      // 屏幕宽度判断（适配响应式）
      const isSmallScreen = window.innerWidth < 992;
      if (isSmallScreen) {
        return '100%'; // 小屏幕始终占满
      }
      // 大屏时：Diy 类型输入框变窄，其他类型正常宽度
      return this.SelectDateType === 'Diy' ? '180px' : '100%';
    }
  },
  data() {
    // 完全保留原有 data 定义，仅新增 3 个缓存字段（不影响原有逻辑）
    return {
      moment,
      dataWindows:30,
      pagination: {
        pageSize: 15,
        showSizeChanger: true,
        showTotal: (total) => `${this.$t('reporting.AlarmHistory.Total')} ${total} ${this.$t('reporting.AlarmHistory.Records')}`
      },
      dataZ: [],
      valueData: '',
      treePageSize: 100,
      scrollPage: 1,
      frontDataZ: [],
      loadExecl: null,
      isLoadExecl: false,
      exportName: "",
      json_fields_cn: {
        "记录时间": {
          field: "MainTime",
          callback: value => formatDate(new Date(value), 'yyyy-MM-dd hh:mm:ss')
        }
      },
      json_fields_en: {
        "Record Time": {
          field: "MainTime",
          callback: value => formatDate(new Date(value), 'yyyy-MM-dd hh:mm:ss')
        }
      },
      json_fields: {},
      json_meta: [
        [{ "key": "charset", "value": "utf-8" }]
      ],
      SelectDateType: 'Day',
      SelectDevice: [],
      SelectDateRange: moment().format("YYYY-MM-DD"),
      SelectAlarmData: [],
      deviceTreeData: [],
      AlarmDataTree: [],
      form: this.$form.createForm(this),
      messageShowLoad: false,
      advanced: true,
      refIconLoading: false,
      baseColumns: [
        {
          width: '180px', // 固定时间列宽度（足够容纳时间字符串）
          slotName: 'reporting.DataHistory.RecordTime',
          scopedSlots: { customRender: 'MainTime', title: 'reporting.DataHistory.RecordTime' },
          dataIndex: 'MainTime',
          fixed: 'left',
          align: 'center',
          customCell: () => ({
            style: {
              whiteSpace: 'nowrap', // 强制单行
              overflow: 'hidden',
              textOverflow: 'ellipsis'
            }
          })
        }
      ],
      columns: [],
      dataSource: [],
      conditionExpress: "",
      selectedRows: [],
      // 新增 3 个缓存字段（仅用于数据处理，不影响原有逻辑）
      rawApiData: [], // 存储接口返回的原始嵌套数据
      cacheDeviceNames: [], // 缓存提取的设备名
      cacheDataNames: [] // 缓存提取的寄存器名
    };
  },
  authorize: {
    // deleteRecord: 'delete'
  },
  mounted() {
    // 完全保留原有逻辑
    this.json_fields = this.$i18n.locale === "CN" ? { ...this.json_fields_cn } : { ...this.json_fields_en };
    this.exportName = `${this.$t('reporting.AlarmHistory.exportName')}_${formatDate(new Date(), 'yyyy-MM-dd hh:mm:ss')}.xls`;
  },
  activated() { },
  filters: {
    formatDate(time) {
      let date = new Date(time);
      return formatDate(date, 'yyyy-MM-dd hh:mm:ss');
    },
  },
  created() {
    // 完全保留原有逻辑
    this.getMonitorTree();
    this.columns = [...this.baseColumns];
  },
  watch: {
    '$route'() {
      this.getMonitorTree();
    },
    '$i18n.locale'(val) {
      this.json_fields = val === "CN" ? { ...this.json_fields_cn } : { ...this.json_fields_en };
      // 新增：语言切换时同步更新导出字段（不影响原有逻辑）
      this.syncExportFields();
    }
  },
  methods: {
    // 完全保留原有所有方法，不做任何修改
    chargeDateType(e) {
      let type = e.target.value;
      if (type === "Day") {
        this.SelectDateRange = moment().format("YYYY-MM-DD");
      } else if (type === "Weekly") {
        const startDate = moment().day(1).format('YYYY-MM-DD');
        const endDate = moment().day(7).format('YYYY-MM-DD');
        this.SelectDateRange = [startDate, endDate];
      } else if (type === "Month") {
        this.SelectDateRange = moment().format("YYYY-MM");
      } else {
        this.SelectDateRange = [moment().add(-1, 'day'), moment()];
      }
    },
    SelectTreeDevice(value, node, extera) {
      this.GetDeviceModelDataList();
    },
    handleSearch(val) {
      this.valueData = val;
      if (!val) {
        this.GetDeviceModelDataList();
      } else {
        this.frontDataZ = [];
        this.scrollPage = 1;
        this.dataZ.forEach(item => {
          if (item.name.indexOf(val) >= 0) {
            this.frontDataZ.push(item);
          }
        });
        this.allDataZ = this.frontDataZ;
        this.frontDataZ = this.frontDataZ.slice(0, 100);
      }
    },
    handlePopupScroll(e) {
      const { target } = e;
      const scrollHeight = target.scrollHeight - target.scrollTop;
      const clientHeight = target.clientHeight;
      if (scrollHeight === 0 && clientHeight === 0) {
        this.scrollPage = 1;
      } else {
        if (scrollHeight < clientHeight + 5) {
          this.scrollPage += 1;
          const max = Math.min(this.treePageSize * this.scrollPage, this.dataZ.length);
          const newData = this.valueData 
            ? this.allDataZ.slice(0, max) 
            : this.dataZ.slice(0, max);
          this.frontDataZ = newData;
        }
      }
    },
    async handleExport() {
      this.isLoadExecl = true;
      this.loadExecl = this.$message.loading(this.$t("reporting.AlarmHistory.LoadingExecl"), 0);
      try {
        await exportExcel(this.dataSource, this.json_fields, this.exportName.replace('.xls', ''));
      } finally {
        this.$message.destroy(this.loadExecl);
        this.isLoadExecl = false;
      }
    },
    filterOption(input, option) {
      return option.componentOptions.children[0].text.toLowerCase().indexOf(input.toLowerCase()) >= 0;
    },
    GetDeviceModelDataList(open) {
      let _t = this;
      this.AlarmDataTree = [];
      _t.dataZ = [];
      _t.frontDataZ = [];
      const params = {
        SelectDevice: this.SelectDevice,
        getType: 2
      };
      GetDeviceModelDataList(params).then(function (res) {
        if (res.data.code === 0) {
          res.data.list.forEach(device => {
            if (device.DataList && device.DataList.length) {
              device.DataList.forEach(dataItem => {
                _t.dataZ.push(dataItem);
                _t.AlarmDataTree.push(dataItem);
              });
            }
          });
          _t.frontDataZ = _t.dataZ.slice(0, 100);
        }
      });
    },
    onDateChange(date, dateString) {
      this.SelectDateRange = dateString;
    },
    onWeeklyDateChange(date, dateString) {
      const startDate = moment(date).day(1).format('YYYY-MM-DD');
      const endDate = moment(date).day(7).format('YYYY-MM-DD');
      this.SelectDateRange = [startDate, endDate];
    },
    getMonitorTree() {
      let _t = this;
      this.deviceTreeData = [];
      getMonitorTree().then(function (res) {
        if (res.data.code === 0) {
          _t.deviceTreeData = res.data.list;
        }
      });
    },

    // -------------------------- 新增以下 3 个方法（不影响原有逻辑）--------------------------
    /**
     * 1. 处理接口返回的嵌套数据：平化为表格结构 + 提取设备名/寄存器名
     * @param {Array} apiData - 接口返回的原始数据（[{ MainTime: "...", Items: [...] }, ...]）
     * @returns {Object} { tableData: 平化后的表格数据, deviceNames: 设备名数组, dataNames: 寄存器名数组 }
     */
    handleApiData(apiData) {
      if (!apiData || !apiData.length) return { tableData: [], deviceNames: [], dataNames: [] };

      // 提取所有设备名（去重）
      const deviceNames = Array.from(
        new Set(apiData.flatMap(item => item.Items.map(sub => sub.DeviceName)))
      );

      // 提取所有寄存器名（去重 + 排序）
      const dataNames = Array.from(
        new Set(apiData.flatMap(item => item.Items.map(sub => sub.DataName)))
      ).sort();

      // 平化数据：将 Items 转为表格行的扁平字段
      const tableData = apiData.map((item, rowIndex) => {
        const flatRow = {
          key: `row_${rowIndex}_${item.MainTime}`, // 生成唯一 key（兼容原有 rowKey 逻辑）
          MainTime: item.MainTime
        };
        deviceNames.forEach(device => {
          dataNames.forEach(dataName => {
            const targetItem = item.Items.find(
              sub => sub.DeviceName === device && sub.DataName === dataName
            );
            flatRow[`${device}_${dataName}`] = targetItem ? targetItem.DataValue : '-';
          });
        });
        return flatRow;
      });

      return { tableData, deviceNames, dataNames };
    },

    /**
     * 2. 生成动态列配置（设备-寄存器列）
     * @param {Array} deviceNames - 设备名数组
     * @param {Array} dataNames - 寄存器名数组
     * @returns {Array} 动态列配置
     */
    generateDynamicColumns(deviceNames, dataNames) {
      if (!deviceNames.length || !dataNames.length) return [];

      return deviceNames.flatMap(device => {
        return dataNames.map(dataName => ({
          width: `${(100 - 12) / (deviceNames.length * dataNames.length)}%`,
          dataIndex: `${device}_${dataName}`,
          key: `${device}_${dataName}`,
          title: `${device}(${dataName})`,
          ellipsis: true,
          tooltip: true,
          sortable: true,
          align: 'center',
          sortMethod: (a, b) => {
            const numA = a === '-' ? 0 : Number(a);
            const numB = b === '-' ? 0 : Number(b);
            return numA - numB;
          }
        }));
      });
    },

    /**
     * 3. 同步导出字段（追加动态列）
     */
    syncExportFields() {
      if (!this.cacheDeviceNames.length || !this.cacheDataNames.length) return;

      const isCN = this.$i18n.locale === "CN";
      // 重置为原始导出字段（避免重复追加）
      this.json_fields = isCN ? { ...this.json_fields_cn } : { ...this.json_fields_en };

      // 追加动态列到导出字段
      this.cacheDeviceNames.forEach(device => {
        this.cacheDataNames.forEach(dataName => {
          const fieldKey = `${device}_${dataName}`;
          const fieldTitle = isCN ? `${device}(${dataName})` : `${device}(${dataName})`;
          this.json_fields[fieldTitle] = fieldKey;
        });
      });
    },

    // -------------------------- 原有 QueryHistoryDataList 方法（仅新增 4 行处理逻辑）--------------------------
    QueryHistoryDataList() {
      let _t = this;
      _t.dataSource = [];
      _t.columns = [..._t.baseColumns];

      const params = {
        deviceList: this.SelectDevice,
        dataList: this.SelectAlarmData,
        dateType: this.SelectDateType,
        dateRange: this.SelectDateRange,
		dataWindows:parseInt(this.dataWindows)
      };

      if (params.dateRange === "" || (Array.isArray(params.dateRange) && !params.dateRange[0])) {
        this.$message.error(this.$t("reporting.AlarmHistory.SelectDateError"));
        return;
      }

      this.messageShowLoad = true;
      GetDataHistoryReport(params).then(function (res) {
        if (res.data.code === 0) {
          // 原有逻辑：直接赋值
          // _t.dataSource = res.data.list;

          // 新增逻辑开始（仅 4 行，不修改原有逻辑）
          _t.rawApiData = res.data.list; // 存储原始接口数据
          const { tableData, deviceNames, dataNames } = _t.handleApiData(_t.rawApiData); // 处理数据
          _t.cacheDeviceNames = deviceNames; // 缓存设备名
          _t.cacheDataNames = dataNames; // 缓存寄存器名
          // 新增逻辑结束

          // 替换原有赋值：使用处理后的表格数据
          _t.dataSource = tableData;

          // 新增：生成动态列并合并到原有列
          const dynamicColumns = _t.generateDynamicColumns(deviceNames, dataNames);
          _t.columns = [..._t.baseColumns, ...dynamicColumns];

          // 新增：同步导出字段
          _t.syncExportFields();

          // 原有逻辑：更新分页和导出文件名
          _t.pagination.total = _t.dataSource.length;
          _t.exportName = `${_t.$t('reporting.DataHistory.exportName')}_${formatDate(new Date(), 'yyyy-MM-dd hh:mm:ss')}.xls`;
        } else {
          _t.$message.error(res.data.msg || _t.$t('reporting.AlarmHistory.QueryFail'));
        }
        _t.messageShowLoad = false;
      }).catch(err => {
        _t.messageShowLoad = false;
        _t.$message.error(_t.$t('reporting.AlarmHistory.QueryFail'));
        console.error('接口请求失败：', err);
      });
    },

    /**
     * 格式化表格数据为AI易理解的结构
     * @returns {Object} 格式化后的分析数据
     */
    formatDataForAI() {
      if (!this.dataSource.length) return null;

      // 提取核心分析信息
      const analysisData = {
        // 基础查询条件
        queryCondition: {
          deviceList: this.SelectDevice.map(item => this.getDeviceNameById(item)), // 设备名称（需实现ID转名称）
          dataList: this.SelectAlarmData.map(item => this.getDataNameById(item)), // 数据项名称
          dateType: this.SelectDateType,
          dateRange: this.SelectDateRange,
          timeWindow: this.dataWindows
        },
        // 表格数据（简化版，保留核心值）
        tableData: this.dataSource.map(item => {
          const simplified = {
            recordTime: item.MainTime,
            values: {}
          };
          // 提取设备-寄存器值
          Object.keys(item).forEach(key => {
            if (key.includes('_') && key !== 'key' && key !== 'MainTime') {
              simplified.values[key] = item[key];
            }
          });
          return simplified;
        }),
        // 数据统计信息
        statistics: {
          totalRecords: this.dataSource.length,
          deviceCount: this.cacheDeviceNames.length,
          dataItemCount: this.cacheDataNames.length,
          timeRange: `${this.dataSource[0]?.MainTime || ''} 至 ${this.dataSource[this.dataSource.length-1]?.MainTime || ''}`
        }
      };

      return analysisData;
    },

    /**
     * 获取设备名称（根据ID，需根据你的设备树数据实现）
     */
    getDeviceNameById(deviceId) {
      // 从deviceTreeData中查找设备名称，示例逻辑：
      const findName = (tree, id) => {
        for (const node of tree) {
          if (node.key === id) return node.text;
          if (node.children && node.children.length) {
            const res = findName(node.children, id);
            if (res) return res;
          }
        }
        return id; // 找不到返回ID
      };
      return findName(this.deviceTreeData, deviceId);
    },

    /**
     * 获取数据项名称（根据ID）
     */
    getDataNameById(dataId) {
      const dataItem = this.dataZ.find(item => item.uuid === dataId);
      return dataItem ? dataItem.name : dataId;
    },

    /**
     * 构造AI分析提示词
     */
    buildAIPrompt() {
      const formattedData = this.formatDataForAI();
      if (!formattedData) return '';

      // 多语言提示词（适配国际化）
      const isCN = this.$i18n.locale === "CN";
      const prompt = isCN ?
          `请分析以下设备历史数据：
1. 查询条件：${JSON.stringify(formattedData.queryCondition, null, 2)}
2. 数据统计：共${formattedData.statistics.totalRecords}条记录，涉及${formattedData.statistics.deviceCount}台设备，${formattedData.statistics.dataItemCount}个数据项
3. 具体数据：${JSON.stringify(formattedData.tableData, null, 2)}

请完成以下分析：
- 数据整体趋势（如数值上升/下降/波动）
- 异常值识别（如数值超出正常范围、突变）
- 设备性能分析（不同设备同一指标对比）
- 关键结论和优化建议
- 分析结果用Markdown格式返回，语言简洁易懂` :
          `Please analyze the following device historical data:
1. Query conditions: ${JSON.stringify(formattedData.queryCondition, null, 2)}
2. Data statistics: Total ${formattedData.statistics.totalRecords} records, involving ${formattedData.statistics.deviceCount} devices, ${formattedData.statistics.dataItemCount} data items
3. Specific data: ${JSON.stringify(formattedData.tableData, null, 2)}

Please complete the following analysis:
- Overall data trends (e.g., rising/falling/fluctuating values)
- Abnormal value identification (e.g., values outside normal range, sudden changes)
- Device performance analysis (comparison of the same indicator across different devices)
- Key conclusions and optimization suggestions
- Return analysis results in Markdown format with concise and understandable language`;

      return prompt;
    },

    /**
     * 调用AI接口进行分析
     */
    async callAIAnalysis() {
      const prompt = this.buildAIPrompt();
      if (!prompt) {
        this.$message.warning(this.$t('reporting.AlarmHistory.NoDataForAnalysis'));
        return '';
      }

      try {
        // 调用你之前的Go AI服务（或直接调用火山引擎API）
        const response = await axios({
          method: 'POST',
          url: `${this.aiConfig.baseURL}/chat/completions`,
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${this.aiConfig.apiKey}`
          },
          data: {
            model: this.aiConfig.model,
            messages: [{ role: 'user', content: prompt }],
            stream: false
          },
          timeout: 60000 // 60秒超时
        });

        if (response.data.error && response.data.error.code !== 0) {
          throw new Error(`${response.data.error.code}: ${response.data.error.message}`);
        }

        return response.data.choices[0]?.message?.content || '';
      } catch (error) {
        console.error('AI分析接口调用失败：', error);
        throw new Error(this.$t('reporting.AlarmHistory.AIAnalysisFail') + (error.message || ''));
      }
    },

    /**
     * 启动AI分析
     */
    async startAIAnalysis() {
      if (!this.dataSource.length) {
        this.$message.warning(this.$t('reporting.AlarmHistory.NoDataForAnalysis'));
        return;
      }

      this.isLoadAI = true;
      this.showAIResult = true;
      this.aiAnalysisResult = '';

      try {
        // 调用AI分析接口
        const result = await this.callAIAnalysis();
        // 将Markdown转为HTML（可选：使用marked库）
        this.aiAnalysisResult = result
            // 简单的Markdown转HTML（如需完整支持可引入marked库）
            .replace(/### (.*)/g, '<h3>$1</h3>')
            .replace(/## (.*)/g, '<h2>$1</h2>')
            .replace(/### (.*)/g, '<h3>$1</h3>')
            .replace(/\* (.*)/g, '<li>$1</li>')
            .replace(/\n/g, '<br>');
      } catch (error) {
        this.$message.error(error.message);
        this.aiAnalysisResult = `<div style="color: #f5222d;">${error.message}</div>`;
      } finally {
        this.isLoadAI = false;
      }
    }
  }
};
</script>

<style lang="less">
// 新增表格滚动容器样式
.ant-table-container {
  height: 100%;
}
// 修复表头固定后的样式问题
.ant-table-header {
  position: sticky !important;
  top: 0 !important;
  z-index: 1;
  background: #fff;
}
// 分页组件样式优化
.ant-pagination {
  margin: 16px 0 !important;
}
// AI分析结果样式
.ai-analysis-result {
  padding: 16px;
  line-height: 1.8;
  font-size: 14px;

  h2 {
    font-size: 16px;
    color: #1f2937;
    margin: 16px 0 8px;
    font-weight: 600;
  }

  h3 {
    font-size: 15px;
    color: #374151;
    margin: 12px 0 6px;
    font-weight: 600;
  }

  li {
    margin: 4px 0 4px 20px;
  }

  br {
    margin-bottom: 4px;
  }
}
/* 完全保留原有样式，不做任何修改 */
.plus-icon-enter-active {
  transition: opacity .5s;
}
.plus-icon-enter {
  opacity: 0;
}
.plus-icon-leave-active {
  transition: opacity .5s;
}
.plus-icon-leave-to {
  opacity: 0;
}
.plus-icon-enter-to {
  opacity: 1;
}

.code-box-actions {
  padding-top: 12px;
  text-align: center;
  opacity: .7;
  transition: opacity .3s;
}
.code-box-meta .demo-description>h4, .code-box-meta>h4 {
  position: absolute;
  top: -14px;
  padding: 1px 8px;
  margin-left: 16px;
  color: #777;
  border-radius: 2px 2px 0 0;
  background: #fff;
  font-size: 14px;
  width: auto;
}
.code-box {
  border: 1px solid #ebedf0;
  border-radius: 2px;
  display: inline-block;
  width: 100%;
  position: relative;
  margin: 0 0 16px;
  transition: all .2s;
}

.search {
  margin-bottom: 54px;
}
.fold {
  width: calc(100% - 216px);
  display: inline-block;
}
.operator {
  margin-bottom: 18px;
}
@media screen and (max-width: 900px) {
  .fold {
    width: 100%;
  }
}

.ant-table-thead > tr > th {
  padding: 10px 10px;
  overflow-wrap: break-word;
}
.ant-table-tbody > tr > td {
  padding: 7px 5px;
  overflow-wrap: break-word;
}

.ant-table-thead>tr>th {
  color: #909399;
  font-weight: 500;
  text-align: left;
  transition: background .3s ease;
}
</style>