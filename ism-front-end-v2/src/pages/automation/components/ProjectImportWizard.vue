<template>
  <div>
    <a-steps :current="currentStep" style="margin-bottom: 24px;">
      <a-step title="上传配置" />
      <a-step title="解析预览" />
      <a-step title="确认执行" />
      <a-step title="执行结果" />
    </a-steps>

    <!-- Step 0: 上传配置 -->
    <div v-if="currentStep === 0">
      <a-row :gutter="24">
        <a-col :span="12">
          <a-card title="Excel 导入" :bordered="true">
            <a-upload-dragger
              name="file"
              :multiple="false"
              :action="uploadUrl"
              :beforeUpload="beforeUpload"
              @change="handleUploadChange"
            >
              <p class="ant-upload-drag-icon">
                <a-icon type="inbox" />
              </p>
              <p class="ant-upload-text">点击或拖拽 Excel 文件到此处</p>
              <p class="ant-upload-hint">
                支持 .xlsx/.xls 格式，需包含模板、主数据、设备清单三个 Sheet
              </p>
            </a-upload-dragger>
          </a-card>
        </a-col>
        <a-col :span="12">
          <a-card title="JSON 配置包导入" :bordered="true">
            <a-upload-dragger
              name="file"
              :multiple="false"
              :action="uploadJsonUrl"
              :beforeUpload="beforeJsonUpload"
              @change="handleJsonUploadChange"
            >
              <p class="ant-upload-drag-icon">
                <a-icon type="file-text" />
              </p>
              <p class="ant-upload-text">点击或拖拽 JSON 配置包到此处</p>
              <p class="ant-upload-hint">
                支持由 Python CLI 生成的项目配置包 JSON
              </p>
            </a-upload-dragger>
          </a-card>
        </a-col>
      </a-row>

      <a-divider />

      <a-form layout="vertical">
        <a-form-item label="导入步骤">
          <a-checkbox-group v-model="importSteps" :options="stepOptions" />
        </a-form-item>
        <a-form-item label="执行模式">
          <a-radio-group v-model="dryRun">
            <a-radio :value="true">仅预览（dry-run）</a-radio>
            <a-radio :value="false">立即执行</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>

      <div style="margin-top: 24px; text-align: right;">
        <a-button type="primary" @click="nextStep" :disabled="!hasFile">
          下一步 <a-icon type="right" />
        </a-button>
      </div>
    </div>

    <!-- Step 1: 解析预览 -->
    <div v-if="currentStep === 1">
      <a-spin :spinning="parsing">
        <a-alert
          v-if="previewData"
          message="解析完成"
          :description="`发现 ${previewData.modelCount} 个数据模型, ${previewData.deviceCount} 台设备, ${previewData.pageCount} 个组态页面`"
          type="success"
          show-icon
          style="margin-bottom: 16px;"
        />
        
        <a-descriptions title="数据模型预览" bordered v-if="previewData && previewData.models">
          <a-descriptions-item v-for="(model, idx) in previewData.models" :key="idx" :label="model.name">
            AI: {{ model.aiCount }} 点 / DI: {{ model.diCount }} 点 / 设备: {{ model.deviceCount }} 台
          </a-descriptions-item>
        </a-descriptions>

        <a-divider />

        <a-descriptions title="设备预览" bordered v-if="previewData && previewData.devices">
          <a-descriptions-item label="设备总数">{{ previewData.devices.length }} 台</a-descriptions-item>
          <a-descriptions-item label="设备类型">{{ deviceTypeSummary }}</a-descriptions-item>
        </a-descriptions>
      </a-spin>

      <div style="margin-top: 24px; text-align: right;">
        <a-button style="margin-right: 8px;" @click="prevStep">上一步</a-button>
        <a-button type="primary" @click="nextStep" :loading="parsing">
          下一步 <a-icon type="right" />
        </a-button>
      </div>
    </div>

    <!-- Step 2: 确认执行 -->
    <div v-if="currentStep === 2">
      <a-alert
        message="确认执行导入"
        :description="dryRun ? '当前为预览模式，不会修改任何数据' : '即将执行实际导入，请确认数据正确'"
        :type="dryRun ? 'info' : 'warning'"
        show-icon
        style="margin-bottom: 16px;"
      />

      <a-card title="执行摘要" :bordered="true">
        <a-row :gutter="16">
          <a-col :span="6">
            <a-statistic title="数据模型" :value="previewData ? previewData.modelCount : 0" />
          </a-col>
          <a-col :span="6">
            <a-statistic title="设备" :value="previewData ? previewData.deviceCount : 0" />
          </a-col>
          <a-col :span="6">
            <a-statistic title="组态页面" :value="previewData ? previewData.pageCount : 0" />
          </a-col>
          <a-col :span="6">
            <a-statistic title="预计耗时" :value="estimatedTime" suffix="秒" />
          </a-col>
        </a-row>
      </a-card>

      <div style="margin-top: 24px; text-align: right;">
        <a-button style="margin-right: 8px;" @click="prevStep">上一步</a-button>
        <a-button type="primary" @click="doImport" :loading="executing">
          {{ dryRun ? '开始预览' : '确认导入' }}
        </a-button>
      </div>
    </div>

    <!-- Step 3: 执行结果 -->
    <div v-if="currentStep === 3">
      <a-spin :spinning="executing" tip="正在执行...">
        <div v-if="taskResult">
          <a-result
            :status="taskResult.status === 'success' ? 'success' : taskResult.status === 'failed' ? 'error' : 'info'"
            :title="taskResult.status === 'success' ? '导入成功' : taskResult.status === 'failed' ? '导入失败' : '执行中...'"
            :sub-title="taskResult.message || taskResult.currentStep"
          >
            <template slot="extra">
              <a-button v-if="taskResult.taskId" @click="viewTaskDetail(taskResult.taskId)">
                查看任务详情
              </a-button>
              <a-button type="primary" @click="resetWizard">
                重新导入
              </a-button>
            </template>
          </a-result>

          <!-- 进度条 -->
          <a-progress
            v-if="taskResult.progress !== undefined"
            :percent="taskResult.progress"
            :status="taskResult.status === 'failed' ? 'exception' : taskResult.progress === 100 ? 'success' : 'active'"
          />

          <!-- 日志 -->
          <a-card title="执行日志" style="margin-top: 16px;" v-if="taskResult.logs && taskResult.logs.length">
            <a-timeline>
              <a-timeline-item
                v-for="(log, idx) in taskResult.logs"
                :key="idx"
                :color="log.status === 'success' ? 'green' : log.status === 'failed' ? 'red' : 'blue'"
              >
                {{ log.step }} - {{ log.name }} - {{ log.status }}
              </a-timeline-item>
            </a-timeline>
          </a-card>
        </div>
      </a-spin>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ProjectImportWizard',
  data() {
    return {
      currentStep: 0,
      hasFile: false,
      parsing: false,
      executing: false,
      dryRun: true,
      importSteps: ['model', 'device', 'dashboard'],
      stepOptions: [
        { label: '数据模型', value: 'model' },
        { label: '设备', value: 'device' },
        { label: '组态大屏', value: 'dashboard' },
      ],
      uploadUrl: '/api/autoGen/uploadExcel',
      uploadJsonUrl: '/api/autoGen/uploadJson',
      previewData: null,
      taskResult: null,
      estimatedTime: 30,
      pollInterval: null,
    }
  },
  computed: {
    deviceTypeSummary() {
      if (!this.previewData || !this.previewData.devices) return ''
      const types = {}
      this.previewData.devices.forEach(d => {
        const type = d.type || '未知'
        types[type] = (types[type] || 0) + 1
      })
      return Object.entries(types).map(([k, v]) => `${k}: ${v}台`).join(' / ')
    },
  },
  beforeDestroy() {
    if (this.pollInterval) {
      clearInterval(this.pollInterval)
    }
  },
  methods: {
    beforeUpload(file) {
      const isExcel = file.type === 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' ||
                      file.type === 'application/vnd.ms-excel'
      if (!isExcel) {
        this.$message.error('只支持 Excel 文件!')
      }
      return isExcel
    },
    beforeJsonUpload(file) {
      const isJson = file.type === 'application/json' || file.name.endsWith('.json')
      if (!isJson) {
        this.$message.error('只支持 JSON 文件!')
      }
      return isJson
    },
    handleUploadChange(info) {
      if (info.file.status === 'done') {
        this.hasFile = true
        this.$message.success(`${info.file.name} 上传成功`)
        // 解析返回的预览数据
        if (info.file.response && info.file.response.code === 200) {
          this.previewData = info.file.response.data
        }
      } else if (info.file.status === 'error') {
        this.$message.error(`${info.file.name} 上传失败`)
      }
    },
    handleJsonUploadChange(info) {
      if (info.file.status === 'done') {
        this.hasFile = true
        this.$message.success(`${info.file.name} 上传成功`)
      }
    },
    nextStep() {
      if (this.currentStep === 1) {
        this.parsePreview()
      }
      this.currentStep++
    },
    prevStep() {
      this.currentStep--
    },
    parsePreview() {
      this.parsing = true
      // 模拟解析请求
      setTimeout(() => {
        this.parsing = false
        this.previewData = {
          modelCount: 3,
          deviceCount: 76,
          pageCount: 26,
          models: [
            { name: 'A20电力仪表', aiCount: 20, diCount: 8, deviceCount: 45 },
            { name: 'A40电力仪表', aiCount: 40, diCount: 16, deviceCount: 18 },
            { name: '施耐德UPS', aiCount: 10, diCount: 4, deviceCount: 13 },
          ],
          devices: Array.from({ length: 76 }, (_, i) => ({
            name: `1A1_U11_S18_${i + 1}`,
            type: i < 45 ? 'A20' : i < 63 ? 'A40' : 'UPS',
          })),
        }
        this.estimatedTime = Math.ceil(this.previewData.deviceCount * 0.5 + this.previewData.pageCount * 2)
      }, 1000)
    },
    doImport() {
      this.executing = true
      this.currentStep = 3
      
      // 模拟执行请求
      const reqBody = {
        dryRun: this.dryRun,
        steps: this.importSteps,
      }
      
      console.log('执行导入:', reqBody)
      
      // 模拟后端响应
      setTimeout(() => {
        this.taskResult = {
          status: 'success',
          taskId: `T-${Date.now()}`,
          progress: 100,
          message: '导入成功',
          currentStep: '全部完成',
          logs: [
            { step: 1, name: '创建A20电力仪表', status: 'success' },
            { step: 2, name: '创建A40电力仪表', status: 'success' },
            { step: 3, name: '创建施耐德UPS', status: 'success' },
            { step: 4, name: '添加76台设备', status: 'success' },
            { step: 5, name: '生成组态大屏', status: 'success' },
          ],
        }
        this.executing = false
      }, 3000)
    },
    viewTaskDetail(taskId) {
      this.$emit('viewTask', taskId)
    },
    resetWizard() {
      this.currentStep = 0
      this.hasFile = false
      this.previewData = null
      this.taskResult = null
      this.dryRun = true
    },
  },
}
</script>
