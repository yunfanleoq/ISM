<template>
  <div>
    <a-row :gutter="16">
      <a-col :span="12">
        <a-card title="项目配置" :bordered="true" style="margin-bottom: 16px;">
          <a-table
            :columns="projectColumns"
            :data-source="projectConfigs"
            :pagination="false"
            row-key="name"
            size="small"
          >
            <template #action="{ record }">
              <a-button type="link" size="small" @click="editProject(record)">
                <a-icon type="edit" /> 编辑
              </a-button>
            </template>
          </a-table>
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card title="Excel 列映射规则" :bordered="true">
          <a-form layout="vertical" :model="excelMapping" size="small">
            <a-row :gutter="8">
              <a-col :span="12">
                <a-form-item label="模板Sheet名称">
                  <a-input v-model="excelMapping.templateSheetName" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="设备名称列">
                  <a-input v-model="excelMapping.deviceNameCol" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="AI 起始地址列">
                  <a-input v-model="excelMapping.aiStartCol" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="DI 起始地址列">
                  <a-input v-model="excelMapping.diStartCol" />
                </a-form-item>
              </a-col>
            </a-row>
            <a-form-item>
              <a-button type="primary" size="small" @click="saveMapping">
                <a-icon type="save" /> 保存映射规则
              </a-button>
            </a-form-item>
          </a-form>
        </a-card>
      </a-col>
    </a-row>

    <a-card title="系统设置" :bordered="true" style="margin-top: 16px;">
      <a-form layout="vertical">
        <a-row :gutter="16">
          <a-col :span="8">
            <a-form-item label="API 基础地址">
              <a-input v-model="systemConfig.apiBaseUrl" placeholder="http://127.0.0.1:8081" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item label="默认超时（秒）">
              <a-input-number v-model="systemConfig.timeout" :min="5" :max="300" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item label="最大重试次数">
              <a-input-number v-model="systemConfig.maxRetries" :min="1" :max="10" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="8">
            <a-form-item label="默认画布宽度">
              <a-input-number v-model="systemConfig.defaultCanvasWidth" :min="800" :max="3840" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item label="默认画布高度">
              <a-input-number v-model="systemConfig.defaultCanvasHeight" :min="600" :max="2160" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item label="默认主题">
              <a-select v-model="systemConfig.defaultTheme" style="width: 100%">
                <a-select-option value="dark">深色</a-select-option>
                <a-select-option value="light">浅色</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item>
          <a-button type="primary" @click="saveSystemConfig">
            <a-icon type="save" /> 保存系统设置
          </a-button>
        </a-form-item>
      </a-form>
    </a-card>

    <!-- 编辑项目配置弹窗 -->
    <a-modal
      v-model="editVisible"
      title="编辑项目配置"
      @ok="saveProjectConfig"
    >
      <a-form layout="vertical" v-if="editingProject">
        <a-form-item label="项目名称">
          <a-input v-model="editingProject.name" />
        </a-form-item>
        <a-form-item label="项目UUID">
          <a-input v-model="editingProject.uuid" disabled />
        </a-form-item>
        <a-form-item label="Excel 路径">
          <a-input v-model="editingProject.excelPath" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script>
export default {
  name: 'ConfigManager',
  data() {
    return {
      projectConfigs: [],
      excelMapping: {
        templateSheetName: '模板',
        deviceNameCol: 'A',
        aiStartCol: 'O',
        diStartCol: 'P',
      },
      systemConfig: {
        apiBaseUrl: 'http://127.0.0.1:8081',
        timeout: 30,
        maxRetries: 3,
        defaultCanvasWidth: 1920,
        defaultCanvasHeight: 1080,
        defaultTheme: 'dark',
      },
      projectColumns: [
        { title: '项目名称', dataIndex: 'name', key: 'name' },
        { title: 'Excel文件', dataIndex: 'excelPath', key: 'excelPath', ellipsis: true },
        { title: '模型数', dataIndex: 'modelCount', key: 'modelCount', width: 80 },
        { title: '设备数', dataIndex: 'deviceCount', key: 'deviceCount', width: 80 },
        { title: '操作', key: 'action', scopedSlots: { customRender: 'action' }, width: 80 },
      ],
      editVisible: false,
      editingProject: null,
    }
  },
  mounted() {
    this.loadProjectConfigs()
    this.loadExcelMapping()
    this.loadSystemConfig()
  },
  methods: {
    loadProjectConfigs() {
      // 模拟数据
      this.projectConfigs = [
        {
          name: '航信机房',
          uuid: '31bc90be-ebc4-dd61-ba9d-ce6e075e40e2',
          excelPath: '1A配电室 172.31.4.14 172.20.255.14.xlsx',
          modelCount: 3,
          deviceCount: 76,
        },
      ]
    },
    loadExcelMapping() {
      // 从本地存储或后端加载
      const saved = localStorage.getItem('ism_automation_excel_mapping')
      if (saved) {
        this.excelMapping = JSON.parse(saved)
      }
    },
    loadSystemConfig() {
      const saved = localStorage.getItem('ism_automation_system_config')
      if (saved) {
        this.systemConfig = JSON.parse(saved)
      }
    },
    editProject(record) {
      this.editingProject = { ...record }
      this.editVisible = true
    },
    saveProjectConfig() {
      const idx = this.projectConfigs.findIndex(p => p.uuid === this.editingProject.uuid)
      if (idx >= 0) {
        this.projectConfigs.splice(idx, 1, this.editingProject)
      }
      this.editVisible = false
      this.$message.success('项目配置已保存')
    },
    saveMapping() {
      localStorage.setItem('ism_automation_excel_mapping', JSON.stringify(this.excelMapping))
      this.$message.success('Excel 映射规则已保存')
    },
    saveSystemConfig() {
      localStorage.setItem('ism_automation_system_config', JSON.stringify(this.systemConfig))
      this.$message.success('系统设置已保存')
    },
  },
}
</script>
