<template>
  <div>
    <a-row :gutter="[16, 16]">
      <a-col
        v-for="template in templates"
        :key="template.templateId"
        :xs="24" :sm="12" :md="8" :lg="6"
      >
        <a-card
          hoverable
          class="template-card"
          :class="{ 'template-card-active': selectedTemplate === template.templateId }"
          @click="selectTemplate(template)"
        >
          <template slot="cover">
            <div class="template-cover" :style="{ background: getCoverColor(template.category) }">
              <a-icon :type="getIcon(template.category)" style="font-size: 48px; color: #fff;" />
              <div class="template-category-badge">{{ getCategoryName(template.category) }}</div>
            </div>
          </template>
          <a-card-meta
            :title="template.name"
            :description="template.description"
          />
          <template slot="actions">
            <a-button type="link" @click.stop="previewTemplate(template)">
              <a-icon type="eye" /> 预览
            </a-button>
            <a-button type="primary" size="small" @click.stop="applyTemplate(template)">
              <a-icon type="check" /> 应用
            </a-button>
          </template>
        </a-card>
      </a-col>
    </a-row>

    <!-- 模板预览弹窗 -->
    <a-modal
      v-model="previewVisible"
      :title="`模板预览: ${previewTemplateData ? previewTemplateData.name : ''}`"
      width="900px"
      :footer="null"
    >
      <a-descriptions bordered :column="2" v-if="previewTemplateData">
        <a-descriptions-item label="模板ID">{{ previewTemplateData.templateId }}</a-descriptions-item>
        <a-descriptions-item label="版本">{{ previewTemplateData.version }}</a-descriptions-item>
        <a-descriptions-item label="分类">{{ getCategoryName(previewTemplateData.category) }}</a-descriptions-item>
        <a-descriptions-item label="类型">{{ previewTemplateData.isBuiltin ? '内置' : '自定义' }}</a-descriptions-item>
        <a-descriptions-item label="描述" :span="2">{{ previewTemplateData.description }}</a-descriptions-item>
      </a-descriptions>

      <a-divider orientation="left">布局参数</a-divider>
      <a-row :gutter="16" v-if="previewTemplateData && previewTemplateData.params">
        <a-col :span="12">
          <a-statistic title="画布宽度" :value="previewTemplateData.params && previewTemplateData.params.canvas ? previewTemplateData.params.canvas.width : undefined" suffix="px" />
        </a-col>
        <a-col :span="12">
          <a-statistic title="画布高度" :value="previewTemplateData.params && previewTemplateData.params.canvas ? previewTemplateData.params.canvas.height : undefined" suffix="px" />
        </a-col>
      </a-row>

      <a-divider orientation="left">页面结构</a-divider>
      <a-timeline v-if="previewTemplateData && previewTemplateData.layouts">
        <a-timeline-item color="blue">
          <p>概览页面（overview）</p>
          <p style="color: #999; font-size: 12px;">包含统计卡片、趋势图、告警列表等</p>
        </a-timeline-item>
        <a-timeline-item color="green">
          <p>子页面（钻探层级）</p>
          <p style="color: #999; font-size: 12px;">根据设备层级自动生成</p>
        </a-timeline-item>
      </a-timeline>

      <div style="text-align: right; margin-top: 16px;">
        <a-button @click="previewVisible = false">关闭</a-button>
        <a-button type="primary" style="margin-left: 8px;" @click="applyTemplate(previewTemplateData)">
          应用此模板
        </a-button>
      </div>
    </a-modal>

    <!-- 应用模板确认弹窗 -->
    <a-modal
      v-model="applyVisible"
      title="应用模板"
      @ok="confirmApply"
      :confirmLoading="applying"
    >
      <a-form layout="vertical">
        <a-form-item label="目标项目">
          <a-select v-model="targetProject" placeholder="选择项目" style="width: 100%">
            <a-select-option v-for="proj in projects" :key="proj.uuid" :value="proj.uuid">
              {{ proj.name }}
            </a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="大屏名称">
          <a-input v-model="dashboardName" placeholder="输入组态大屏名称" />
        </a-form-item>
        <a-form-item label="主题">
          <a-radio-group v-model="theme">
            <a-radio-button value="dark">深色</a-radio-button>
            <a-radio-button value="light">浅色</a-radio-button>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script>
export default {
  name: 'TemplateMarket',
  data() {
    return {
      templates: [],
      selectedTemplate: null,
      previewVisible: false,
      applyVisible: false,
      applying: false,
      previewTemplateData: null,
      targetProject: '',
      dashboardName: '',
      theme: 'dark',
      projects: [],
      categoryColors: {
        industrial: '#1890ff',
        'data_center': '#52c41a',
        building: '#faad14',
        energy: '#eb2f96',
      },
      categoryIcons: {
        industrial: 'build',
        'data_center': 'cloud-server',
        building: 'bank',
        energy: 'thunderbolt',
      },
      categoryNames: {
        industrial: '工业',
        'data_center': '数据中心',
        building: '楼宇',
        energy: '能源',
      },
    }
  },
  mounted() {
    this.loadTemplates()
    this.loadProjects()
  },
  methods: {
    loadTemplates() {
      // 模拟从后端加载
      this.templates = [
        {
          templateId: 'industrial_4level',
          name: '工业配电室四级钻探大屏',
          description: '适用于多柜/多设备组场景，含概览、柜、设备组、设备详情四层',
          category: 'industrial',
          version: '1.0',
          isBuiltin: 1,
          params: {
            canvas: { width: 1920, height: 1080 },
          },
          layouts: {
            overview: {},
            building: {},
            floor: {},
            device: {},
          },
        },
        {
          templateId: 'data_center_2level',
          name: '数据中心二级监控大屏',
          description: '适用于简单场景：概览 + 设备列表两级',
          category: 'data_center',
          version: '1.0',
          isBuiltin: 1,
          params: {
            canvas: { width: 1920, height: 1080 },
          },
          layouts: {
            overview: {},
            device: {},
          },
        },
        {
          templateId: 'building_energy',
          name: '楼宇能耗监控大屏',
          description: '适用于楼宇能耗管理场景',
          category: 'building',
          version: '1.0',
          isBuiltin: 1,
          params: {
            canvas: { width: 1920, height: 1080 },
          },
          layouts: {
            overview: {},
          },
        },
      ]
    },
    loadProjects() {
      // 模拟项目列表
      this.projects = [
        { uuid: '31bc90be-ebc4-dd61-ba9d-ce6e075e40e2', name: '航信机房' },
      ]
    },
    getCoverColor(category) {
      return this.categoryColors[category] || '#1890ff'
    },
    getIcon(category) {
      return this.categoryIcons[category] || 'appstore'
    },
    getCategoryName(category) {
      return this.categoryNames[category] || category
    },
    selectTemplate(template) {
      this.selectedTemplate = template.templateId
    },
    previewTemplate(template) {
      this.previewTemplateData = template
      this.previewVisible = true
    },
    applyTemplate(template) {
      this.selectedTemplate = template.templateId
      this.previewTemplateData = template
      this.applyVisible = true
      this.dashboardName = template.name
    },
    confirmApply() {
      if (!this.targetProject) {
        this.$message.error('请选择目标项目')
        return
      }
      this.applying = true
      // 模拟应用模板
      setTimeout(() => {
        this.applying = false
        this.applyVisible = false
        this.$message.success('模板应用成功！')
        this.$emit('templateApplied', {
          templateId: this.selectedTemplate,
          projectUuid: this.targetProject,
          name: this.dashboardName,
        })
      }, 2000)
    },
  },
}
</script>

<style scoped>
.template-card {
  cursor: pointer;
  transition: all 0.3s;
}
.template-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}
.template-card-active {
  border: 2px solid #1890ff;
}
.template-cover {
  height: 160px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: relative;
}
.template-category-badge {
  position: absolute;
  top: 8px;
  right: 8px;
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}
</style>
