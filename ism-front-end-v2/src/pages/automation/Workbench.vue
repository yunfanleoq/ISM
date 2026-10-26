<template>
  <div style="padding: 10px; min-height: 600px;">
    <a-card :bordered="false" :title="$t('automation.workbench') || '自动化工作台'" style="min-height: 700px;">
      <a-tabs default-active-key="import" @change="handleTabChange">
        
        <!-- Tab 1: 项目导入向导 -->
        <a-tab-pane key="import" :tab="$t('automation.projectImport') || '项目导入'">
          <ProjectImportWizard />
        </a-tab-pane>
        
        <!-- Tab 2: 模板市场 -->
        <a-tab-pane key="template" :tab="$t('automation.templateMarket') || '模板市场'">
          <TemplateMarket />
        </a-tab-pane>
        
        <!-- Tab 3: 任务历史 -->
        <a-tab-pane key="history" :tab="$t('automation.taskHistory') || '任务历史'">
          <TaskHistory />
        </a-tab-pane>
        
        <!-- Tab 4: 配置管理 -->
        <a-tab-pane key="config" :tab="$t('automation.configManager') || '配置管理'">
          <ConfigManager />
        </a-tab-pane>
        
      </a-tabs>
    </a-card>
  </div>
</template>

<script>
import ProjectImportWizard from './components/ProjectImportWizard.vue'
import TemplateMarket from './components/TemplateMarket.vue'
import TaskHistory from './components/TaskHistory.vue'
import ConfigManager from './components/ConfigManager.vue'

export default {
  name: 'AutoGenWorkbench',
  components: {
    ProjectImportWizard,
    TemplateMarket,
    TaskHistory,
    ConfigManager,
  },
  data() {
    return {
      activeTab: 'import',
    }
  },
  methods: {
    handleTabChange(key) {
      this.activeTab = key
    },
  },
}
</script>
