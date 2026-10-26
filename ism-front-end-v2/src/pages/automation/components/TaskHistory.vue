<template>
  <div>
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-statistic title="总任务" :value="totalCount" />
      </a-col>
      <a-col :span="6">
        <a-statistic title="成功" :value="successCount" :value-style="{ color: '#3f8600' }" />
      </a-col>
      <a-col :span="6">
        <a-statistic title="失败" :value="failedCount" :value-style="{ color: '#cf1322' }" />
      </a-col>
      <a-col :span="6">
        <a-statistic title="回滚" :value="rollbackCount" />
      </a-col>
    </a-row>

    <a-table
      :columns="columns"
      :data-source="taskList"
      :loading="loading"
      :pagination="pagination"
      row-key="taskId"
    >
      <template slot="status" slot-scope="text, record">
        <a-badge
          :status="getStatusType(record.status)"
          :text="getStatusText(record.status)"
        />
      </template>
      <template slot="progress" slot-scope="text, record">
        <a-progress
          :percent="record.progress"
          :size="'small'"
          :status="record.status === 'failed' ? 'exception' : record.progress === 100 ? 'success' : 'active'"
        />
      </template>
      <template slot="type" slot-scope="text, record">
        <a-tag :color="getTypeColor(record.type)">{{ getTypeText(record.type) }}</a-tag>
      </template>
      <template slot="action" slot-scope="text, record">
        <a-space>
          <a-button type="link" size="small" @click="viewDetail(record)">
            <a-icon type="eye" /> 详情
          </a-button>
          <a-button
            v-if="record.status === 'success'"
            type="link"
            size="small"
            danger
            @click="rollbackTask(record)"
          >
            <a-icon type="rollback" /> 回滚
          </a-button>
          <a-button
            v-if="record.status === 'failed'"
            type="link"
            size="small"
            @click="retryTask(record)"
          >
            <a-icon type="reload" /> 重试
          </a-button>
        </a-space>
      </template>
    </a-table>

    <!-- 任务详情弹窗 -->
    <a-modal
      v-model="detailVisible"
      :title="`任务详情: ${selectedTask ? selectedTask.taskId : ''}`"
      width="800px"
      :footer="null"
    >
      <a-descriptions bordered :column="2" v-if="selectedTask">
        <a-descriptions-item label="任务ID">{{ selectedTask.taskId }}</a-descriptions-item>
        <a-descriptions-item label="类型">{{ getTypeText(selectedTask.type) }}</a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-badge :status="getStatusType(selectedTask.status)" :text="getStatusText(selectedTask.status)" />
        </a-descriptions-item>
        <a-descriptions-item label="进度">{{ selectedTask.progress }}%</a-descriptions-item>
        <a-descriptions-item label="当前步骤">{{ selectedTask.currentStep }}</a-descriptions-item>
        <a-descriptions-item label="创建时间">{{ formatDate(selectedTask.createdAt) }}</a-descriptions-item>
        <a-descriptions-item label="更新时间">{{ formatDate(selectedTask.updatedAt) }}</a-descriptions-item>
        <a-descriptions-item label="错误信息" v-if="selectedTask.error">
          <span style="color: #cf1322;">{{ selectedTask.error }}</span>
        </a-descriptions-item>
      </a-descriptions>

      <a-divider orientation="left">操作日志</a-divider>
      <a-timeline v-if="selectedTask && selectedTask.operations">
        <a-timeline-item
          v-for="(op, idx) in selectedTask.operations"
          :key="idx"
          :color="op.status === 'success' ? 'green' : op.status === 'failed' ? 'red' : 'blue'"
        >
          <p><strong>{{ op.type }}</strong> - {{ op.name }}</p>
          <p style="color: #999; font-size: 12px;">{{ op.action }} - {{ op.status }}</p>
        </a-timeline-item>
      </a-timeline>
    </a-modal>
  </div>
</template>

<script>
export default {
  name: 'TaskHistory',
  data() {
    return {
      loading: false,
      taskList: [],
      totalCount: 0,
      successCount: 0,
      failedCount: 0,
      rollbackCount: 0,
      detailVisible: false,
      selectedTask: null,
      columns: [
        { title: '任务ID', dataIndex: 'taskId', key: 'taskId', width: 180 },
        { title: '类型', dataIndex: 'type', key: 'type', scopedSlots: { customRender: 'type' } },
        { title: '状态', dataIndex: 'status', key: 'status', scopedSlots: { customRender: 'status' } },
        { title: '进度', dataIndex: 'progress', key: 'progress', scopedSlots: { customRender: 'progress' } },
        { title: '当前步骤', dataIndex: 'currentStep', key: 'currentStep' },
        { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt' },
        { title: '操作', key: 'action', scopedSlots: { customRender: 'action' }, width: 180 },
      ],
      pagination: {
        pageSize: 10,
        showSizeChanger: true,
        showTotal: (total) => `共 ${total} 条`,
      },
      typeMap: {
        project_import: '项目导入',
        model_import: '模型导入',
        device_import: '设备导入',
        dashboard_generate: '大屏生成',
      },
      typeColors: {
        project_import: 'blue',
        model_import: 'cyan',
        device_import: 'green',
        dashboard_generate: 'purple',
      },
      statusMap: {
        pending: '等待中',
        running: '执行中',
        success: '成功',
        failed: '失败',
        rolled_back: '已回滚',
      },
      statusTypes: {
        pending: 'default',
        running: 'processing',
        success: 'success',
        failed: 'error',
        rolled_back: 'warning',
      },
    }
  },
  mounted() {
    this.loadTasks()
  },
  methods: {
    loadTasks() {
      this.loading = true
      // 模拟数据
      setTimeout(() => {
        this.taskList = [
          {
            taskId: 'T-2026062401',
            type: 'project_import',
            status: 'success',
            progress: 100,
            currentStep: '全部完成',
            createdAt: '2026-06-24 10:00:00',
            updatedAt: '2026-06-24 10:05:30',
            operations: [
              { type: 'model', name: 'A20电力仪表', action: 'create', status: 'success' },
              { type: 'device', name: '76台设备', action: 'create', status: 'success' },
            ],
          },
          {
            taskId: 'T-2026062402',
            type: 'dashboard_generate',
            status: 'success',
            progress: 100,
            currentStep: '全部完成',
            createdAt: '2026-06-24 11:00:00',
            updatedAt: '2026-06-24 11:02:15',
            operations: [
              { type: 'display_model', name: '航信机房大屏', action: 'create', status: 'success' },
            ],
          },
          {
            taskId: 'T-2026062403',
            type: 'model_import',
            status: 'failed',
            progress: 45,
            currentStep: '创建A40电力仪表',
            error: '模型名称已存在',
            createdAt: '2026-06-24 12:00:00',
            updatedAt: '2026-06-24 12:00:30',
            operations: [
              { type: 'model', name: 'A20电力仪表', action: 'create', status: 'success' },
              { type: 'model', name: 'A40电力仪表', action: 'create', status: 'failed' },
            ],
          },
        ]
        this.totalCount = this.taskList.length
        this.successCount = this.taskList.filter(t => t.status === 'success').length
        this.failedCount = this.taskList.filter(t => t.status === 'failed').length
        this.rollbackCount = this.taskList.filter(t => t.status === 'rolled_back').length
        this.loading = false
      }, 1000)
    },
    getTypeText(type) {
      return this.typeMap[type] || type
    },
    getTypeColor(type) {
      return this.typeColors[type] || 'default'
    },
    getStatusText(status) {
      return this.statusMap[status] || status
    },
    getStatusType(status) {
      return this.statusTypes[status] || 'default'
    },
    formatDate(dateStr) {
      if (!dateStr) return '-'
      return new Date(dateStr).toLocaleString('zh-CN')
    },
    viewDetail(record) {
      this.selectedTask = record
      this.detailVisible = true
    },
    rollbackTask(record) {
      this.$confirm({
        title: '确认回滚',
        content: `回滚任务 ${record.taskId} 将删除该任务创建的所有数据，是否继续？`,
        onOk: () => {
          this.$message.success(`任务 ${record.taskId} 已回滚`)
          record.status = 'rolled_back'
          this.rollbackCount++
        },
      })
    },
    retryTask(record) {
      this.$message.info(`重试任务 ${record.taskId}`)
      record.status = 'running'
      record.progress = 0
      setTimeout(() => {
        record.status = 'success'
        record.progress = 100
        this.$message.success('重试成功')
      }, 3000)
    },
  },
}
</script>
