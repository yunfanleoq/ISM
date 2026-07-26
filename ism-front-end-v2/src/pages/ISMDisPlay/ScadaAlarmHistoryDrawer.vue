<template>
  <!--
    大屏运行态：活跃告警面板「历史查询」入口。
    全屏抽屉复用报表页告警历史查询，避免侧栏 foreignObject 塞 Ant Design 表单。
  -->
  <a-drawer
    :visible="visible"
    placement="right"
    :width="'92%'"
    :destroyOnClose="true"
    :bodyStyle="bodyStyle"
    :drawerStyle="drawerStyle"
    :headerStyle="headerStyle"
    wrapClassName="scada-alarm-history-drawer"
    @close="onClose"
  >
    <template slot="title">
      <div class="sahd-title-row">
        <span class="sahd-title">告警历史查询</span>
        <a
          class="sahd-link"
          href="javascript:;"
          title="打开完整报表页"
          @click.prevent="goFullReport"
        >打开完整报表页 →</a>
      </div>
    </template>

    <div class="sahd-body">
      <AlarmHistoryReport v-if="visible" />
    </div>
  </a-drawer>
</template>

<script>
const AlarmHistoryReport = () => import(
  /* webpackChunkName: "scada-alarm-history-report" */
  '@/pages/reporting/alarmReport/alarmHistory.vue'
)

export default {
  name: 'ScadaAlarmHistoryDrawer',
  components: { AlarmHistoryReport },
  props: {
    visible: { type: Boolean, default: false },
    projectUuid: { type: String, default: '' },
  },
  computed: {
    bodyStyle() {
      return {
        padding: '12px 16px 20px',
        background: '#0b1420',
        height: 'calc(100% - 55px)',
        overflow: 'auto',
      }
    },
    drawerStyle() {
      return { background: '#0b1420' }
    },
    headerStyle() {
      return {
        background: 'linear-gradient(180deg, #122436, #0b1c2b)',
        borderBottom: '1px solid rgba(0, 229, 255, 0.22)',
        color: '#e8f1ff',
        padding: '12px 16px',
      }
    },
  },
  methods: {
    onClose() {
      this.$emit('update:visible', false)
      this.$emit('close')
    },
    goFullReport() {
      // 报表路由在管理端菜单下；新开页不打断大屏盯屏
      const { href } = this.$router.resolve({ path: '/Reporting/AlarmHistory' })
      window.open(href, '_blank')
    },
  },
}
</script>

<style>
/* wrapClassName 挂在 portal 上，需非 scoped */
.scada-alarm-history-drawer .ant-drawer-content {
  background: #0b1420;
}
.scada-alarm-history-drawer .ant-drawer-close {
  color: #9fefff;
}
.scada-alarm-history-drawer .ant-drawer-close:hover {
  color: #00e5ff;
}
.scada-alarm-history-drawer .sahd-title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding-right: 28px;
}
.scada-alarm-history-drawer .sahd-title {
  font-size: 16px;
  font-weight: 600;
  color: #ff6b35;
  letter-spacing: 1px;
}
.scada-alarm-history-drawer .sahd-link {
  font-size: 13px;
  color: #00e5ff;
  white-space: nowrap;
}
.scada-alarm-history-drawer .sahd-link:hover {
  color: #9fefff;
}
.scada-alarm-history-drawer .sahd-body {
  min-height: 100%;
}
.scada-alarm-history-drawer .sahd-body .ant-card {
  background: transparent;
  border: 0;
  color: #e8f1ff;
}
.scada-alarm-history-drawer .sahd-body .ant-card-body {
  padding: 8px 0;
}
.scada-alarm-history-drawer .sahd-body .ant-form-item-label > label,
.scada-alarm-history-drawer .sahd-body .ant-table,
.scada-alarm-history-drawer .sahd-body .ant-table-thead > tr > th,
.scada-alarm-history-drawer .sahd-body .ant-table-tbody > tr > td {
  color: #c8d6e8;
}
.scada-alarm-history-drawer .sahd-body .ant-table-thead > tr > th {
  background: #0d2438;
  border-bottom-color: #1e3a5f;
}
.scada-alarm-history-drawer .sahd-body .ant-table-tbody > tr > td {
  border-bottom-color: #1e3a5f;
  background: transparent;
}
.scada-alarm-history-drawer .sahd-body .ant-table-tbody > tr:hover > td {
  background: #12304a !important;
}
</style>
