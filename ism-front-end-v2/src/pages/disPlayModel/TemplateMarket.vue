<template>
  <div class="template-market">
    <div class="market-header">
      <div class="header-left">
        <h2 class="market-title">
          <a-icon type="appstore" style="margin-right: 8px;" />
          模板市场
        </h2>
        <p class="market-subtitle">50+ 电力行业专业模板，一键创建大屏应用</p>
      </div>
      <div class="header-right">
        <a-input-search
          v-model="searchText"
          placeholder="搜索模板名称或关键词..."
          class="search-input"
          @search="onSearch"
          @change="onSearch"
        >
          <a-icon slot="prefix" type="search" />
        </a-input-search>
      </div>
    </div>

    <div class="category-bar">
      <a-tag
        v-for="cat in categoryList"
        :key="cat.key"
        :color="selectedCategory === cat.key ? cat.color : undefined"
        :class="{ 'cat-tag': true, 'cat-active': selectedCategory === cat.key }"
        @click="selectCategory(cat.key)"
      >
        {{ cat.label }}
        <span class="cat-count">{{ cat.count }}</span>
      </a-tag>
      <a-tag
        :class="{ 'cat-tag': true, 'cat-active': selectedCategory === '' }"
        @click="selectCategory('')"
        color=""
      >
        全部
        <span class="cat-count">{{ totalCount }}</span>
      </a-tag>
    </div>

    <a-spin :spinning="loading" tip="加载模板中...">
      <div class="template-grid" v-if="filteredTemplates.length > 0">
        <a-row :gutter="[16, 16]">
          <a-col
            v-for="item in filteredTemplates"
            :key="item.id"
            :xs="24" :sm="12" :md="8" :lg="6" :xl="6"
          >
            <a-card
              class="template-card"
              hoverable
              @click="showTemplateDetail(item)"
            >
              <div class="card-cover" :class="'cover-' + item.category">
                <div class="cover-icon">
                  <a-icon :type="getCategoryIcon(item.category)" />
                </div>
                <div class="cover-overlay">
                  <span class="overlay-text">点击查看详情</span>
                </div>
                <a-tag class="cover-category-tag" :color="getCategoryColor(item.category)">
                  {{ getCategoryLabel(item.category) }}
                </a-tag>
              </div>

              <a-card-meta>
                <template #title>
                  <span class="card-title">{{ item.name }}</span>
                </template>
                <template #description>
                  <div class="card-keywords">
                    <a-tag
                      v-for="kw in item.keywords.slice(0, 3)"
                      :key="kw"
                      class="kw-tag"
                    >{{ kw }}</a-tag>
                    <a-tag v-if="item.keywords.length > 3" class="kw-tag kw-more">
                      +{{ item.keywords.length - 3 }}
                    </a-tag>
                  </div>
                </template>
              </a-card-meta>

              <template #actions>
                <a-tooltip title="预览模板">
                  <span @click.stop="previewTemplate(item)">
                    <a-icon type="eye" /> 预览
                  </span>
                </a-tooltip>
                <a-tooltip title="使用此模板创建应用">
                  <span class="action-use" @click.stop="useTemplate(item)">
                    <a-icon type="plus-circle" /> 使用
                  </span>
                </a-tooltip>
                <a-tooltip title="直接创建并进入编辑器">
                  <span class="action-edit" @click.stop="editInEditor(item)">
                    <a-icon type="edit" /> 编辑
                  </span>
                </a-tooltip>
              </template>
            </a-card>
          </a-col>
        </a-row>
      </div>

      <a-empty
        v-else
        description="未找到匹配的模板"
        class="empty-state"
      >
        <a-button type="primary" @click="resetFilters">重置筛选</a-button>
      </a-empty>
    </a-spin>

    <!-- 模板详情弹窗 -->
    <a-modal
      v-model="detailVisible"
      :title="currentTemplate ? currentTemplate.name : '模板详情'"
      :footer="null"
      width="680px"
      class="detail-modal"
      :bodyStyle="{ padding: '24px', background: '#141a2e' }"
    >
      <div v-if="currentTemplate" class="detail-content">
        <div class="detail-header">
          <div class="detail-icon-wrap" :class="'icon-' + currentTemplate.category">
            <a-icon :type="getCategoryIcon(currentTemplate.category)" />
          </div>
          <div class="detail-info">
            <h3>{{ currentTemplate.name }}</h3>
            <a-tag :color="getCategoryColor(currentTemplate.category)">
              {{ getCategoryLabel(currentTemplate.category) }}
            </a-tag>
            <span class="detail-id">ID: {{ currentTemplate.id }}</span>
          </div>
        </div>

        <div class="detail-section">
          <h4><a-icon type="tags" /> 关键词</h4>
          <div class="detail-tags">
            <a-tag v-for="kw in currentTemplate.keywords" :key="kw" color="blue">
              {{ kw }}
            </a-tag>
          </div>
        </div>

        <div class="detail-section">
          <h4><a-icon type="info-circle" /> 模板描述</h4>
          <p class="detail-desc">
            该模板为"{{ currentTemplate.name }}"，属于{{ getCategoryLabel(currentTemplate.category) }}类模板，
            包含预配置的组态页面结构、数据绑定点位和可视化组件布局。
            使用此模板可快速创建专业的电力监控大屏应用。
          </p>
        </div>

        <div class="detail-actions">
          <a-button type="primary" size="large" icon="plus-circle" @click="useTemplate(currentTemplate)">
            使用此模板创建应用
          </a-button>
          <a-button size="large" icon="eye" @click="previewTemplate(currentTemplate)">
            预览模板效果
          </a-button>
          <a-button type="dashed" size="large" icon="edit" @click="editInEditor(currentTemplate)">
            直接编辑此模板
          </a-button>
        </div>
      </div>
    </a-modal>

    <!-- 创建应用确认弹窗 -->
    <a-modal
      v-model="createVisible"
      title="从模板创建应用"
      :confirmLoading="creating"
      @ok="doCreateApp"
      @cancel="createVisible = false"
    >
      <a-form :form="createForm">
        <a-form-item label="模板名称">
          <a-input :value="selectedTemplate ? selectedTemplate.name : ''" disabled />
        </a-form-item>
        <a-form-item label="应用名称">
          <a-input
            v-decorator="['appName', {
              rules: [{ required: true, message: '请输入应用名称' }],
              initialValue: selectedTemplate ? selectedTemplate.name : ''
            }]"
            placeholder="请输入新应用的名称"
          />
        </a-form-item>
        <a-form-item label="应用描述">
          <a-textarea
            v-decorator="['appDescription', {
              rules: [{ required: true, message: '请输入应用描述' }],
              initialValue: selectedTemplate ? '基于模板 ' + selectedTemplate.name + ' 创建的大屏应用' : ''
            }]"
            placeholder="请输入应用描述"
            :rows="3"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script>
import templateData from '@/assets/templates/index.json'
import { displayModelAdd, displayModelList, getDisplayModelLayerData, setDisplayModelLayerData } from '@/services/displayModel'

export default {
  name: 'TemplateMarket',
  i18n: require('../../i18n/language'),
  data() {
    return {
      loading: false,
      creating: false,
      editorMode: false,
      searchText: '',
      selectedCategory: '',
      detailVisible: false,
      createVisible: false,
      currentTemplate: null,
      selectedTemplate: null,
      createForm: this.$form.createForm(this),
      categories: templateData.categories,
      templates: templateData.templates,
      categoryColors: {
        substation: '#00d4ff',
        generation: '#52c41a',
        industrial: '#fa8c16',
        integrated: '#722ed1',
        cabinet: '#eb2f96',
        mobile: '#13c2c2'
      },
      categoryIcons: {
        substation: 'thunderbolt',
        generation: 'environment',
        industrial: 'setting',
        integrated: 'dashboard',
        cabinet: 'hdd',
        mobile: 'mobile'
      }
    }
  },
  computed: {
    categoryList() {
      return Object.keys(this.categories).map(key => ({
        key,
        label: this.categories[key].label,
        count: this.categories[key].count,
        color: this.categoryColors[key]
      }))
    },
    totalCount() {
      return this.templates.length
    },
    filteredTemplates() {
      let list = this.templates
      if (this.selectedCategory) {
        list = list.filter(t => t.category === this.selectedCategory)
      }
      if (this.searchText.trim()) {
        const keyword = this.searchText.trim().toLowerCase()
        list = list.filter(t =>
          t.name.toLowerCase().includes(keyword) ||
          t.keywords.some(kw => kw.toLowerCase().includes(keyword))
        )
      }
      return list
    }
  },
  methods: {
    selectCategory(key) {
      this.selectedCategory = key
    },
    onSearch() {},
    resetFilters() {
      this.searchText = ''
      this.selectedCategory = ''
    },
    getCategoryColor(category) {
      return this.categoryColors[category] || '#1890ff'
    },
    getCategoryLabel(category) {
      const cat = this.categories[category]
      return cat ? cat.label : category
    },
    getCategoryIcon(category) {
      return this.categoryIcons[category] || 'project'
    },
    showTemplateDetail(item) {
      this.currentTemplate = item
      this.detailVisible = true
    },
    previewTemplate(item) {
      this.$message.info(`预览功能: ${item.name} - 模板将在组态编辑器中打开`)
      this.$router.push({ path: '/DisPlayEditor/preview', query: { template: item.id } })
    },
    useTemplate(item) {
      this.selectedTemplate = item
      this.editorMode = false
      this.createVisible = true
      this.$nextTick(() => {
        this.createForm.setFieldsValue({
          appName: item.name,
          appDescription: `基于模板「${item.name}」创建的大屏应用`
        })
      })
    },
    editInEditor(item) {
      this.selectedTemplate = item
      this.editorMode = true
      this.createVisible = true
      const appNameSuffix = item.ComponentData ? '（编辑模式）' : ''
      this.$nextTick(() => {
        this.createForm.setFieldsValue({
          appName: item.name + appNameSuffix,
          appDescription: `基于模板「${item.name}」创建，将直接进入编辑器`
        })
      })
    },
    doCreateApp() {
      this.createForm.validateFields((err, values) => {
        if (err) return
        this.creating = true
        const isEditorMode = this.editorMode
        const templateData = this.selectedTemplate
        const params = {
          name: values.appName,
          description: values.appDescription,
          DisplayType: 1
        }
        displayModelAdd(params).then(res => {
          if (res.data.code === 4002) {
            this.createVisible = false
            if (isEditorMode) {
              // 编辑模式：查找新创建的应用并直接进入编辑器
              this.$message.loading({ content: '正在初始化编辑器...', key: 'editLoading', duration: 0 })
              this.findAndOpenEditor(values.appName, templateData)
            } else {
              // 使用模式：返回应用管理列表
              this.creating = false
              this.$message.success(`应用「${values.appName}」创建成功！请在应用管理中编辑和发布`, 4)
              this.$router.push({ path: '/Application' })
            }
          } else if (res.data.code === 4001) {
            this.creating = false
            this.$message.error('应用名称已存在，请修改名称')
          } else if (res.data.code === 4005) {
            this.creating = false
            this.$message.error('应用数量已达上限')
          } else {
            this.creating = false
            this.$message.error('创建失败，请稍后重试')
          }
        }).catch(() => {
          this.creating = false
          this.$message.error('创建失败，请检查网络连接')
        })
      })
    },
    findAndOpenEditor(appName, templateItem) {
      const params = { DisplayType: 1 }
      displayModelList(params).then(res => {
        this.creating = false
        if (res.data.list) {
          const found = res.data.list.find(item => item.name === appName)
          if (found) {
            const uid = found.displayUid
            // 如果模板有 ComponentData，则更新 demo 页面的图层和组件
            if (templateItem && templateItem.ComponentData) {
              this.applyTemplateData(uid, templateItem.ComponentData)
            } else {
              this.$message.destroy('editLoading')
              this.$message.success(`应用「${appName}」创建成功，正在进入编辑器...`, 2)
              this.$router.push({ path: `/DisPlayEditor/${uid}` })
            }
          } else {
            this.$message.destroy('editLoading')
            this.$message.warning('应用创建成功，但未能找到记录，请手动进入应用管理查看', 4)
            this.$router.push({ path: '/Application' })
          }
        } else {
          this.$message.destroy('editLoading')
          this.$message.warning('应用创建成功，请手动进入应用管理查看', 4)
          this.$router.push({ path: '/Application' })
        }
      }).catch(() => {
        this.creating = false
        this.$message.destroy('editLoading')
        this.$message.warning('应用创建成功，请手动进入应用管理查看', 4)
        this.$router.push({ path: '/Application' })
      })
    },
    applyTemplateData(uid, componentDataB64) {
      let pageConfig = null
      try {
        const decoded = atob(componentDataB64)
        pageConfig = JSON.parse(decoded)
      } catch (e) {
        this.$message.destroy('editLoading')
        this.$message.warning('模板数据解析失败，将以空白页面进入编辑器', 3)
        this.$router.push({ path: `/DisPlayEditor/${uid}` })
        return
      }
      // 获取当前模型的图层数据（找到 demo 页面）
      getDisplayModelLayerData({ muid: uid }).then(layerRes => {
        if (layerRes.data && layerRes.data.code === 200 && layerRes.data.list) {
          const layers = layerRes.data.list
          const demoPage = layers.find(l => l.PageName === 'demo')
          if (demoPage && pageConfig.pageConfig) {
            const { layer, components } = pageConfig.pageConfig
            const saveParams = {
              muid: uid,
              pageid: demoPage.PageId,
              saveData: {
                layer: layer || { backColor: '#0a1a3a', backgroundImage: '', widthHeightRatio: '', width: 1920, height: 1080 },
                components: components || { cells: [] }
              }
            }
            setDisplayModelLayerData(saveParams).then(saveRes => {
              this.$message.destroy('editLoading')
              if (saveRes.data && saveRes.data.code === 200) {
                this.$message.success('模板页面已初始化，正在进入编辑器...', 2)
              } else {
                this.$message.warning('页面初始化部分失败，将进入编辑器', 3)
              }
              this.$router.push({ path: `/DisPlayEditor/${uid}` })
            }).catch(() => {
              this.$message.destroy('editLoading')
              this.$message.warning('页面初始化失败，将以默认页面进入编辑器', 3)
              this.$router.push({ path: `/DisPlayEditor/${uid}` })
            })
          } else {
            // 没有 demo 页面或没有 pageConfig，直接进入编辑器
            this.$message.destroy('editLoading')
            this.$message.success('应用已创建，正在进入编辑器...', 2)
            this.$router.push({ path: `/DisPlayEditor/${uid}` })
          }
        } else {
          this.$message.destroy('editLoading')
          this.$message.success('应用已创建，正在进入编辑器...', 2)
          this.$router.push({ path: `/DisPlayEditor/${uid}` })
        }
      }).catch(() => {
        this.$message.destroy('editLoading')
        this.$message.warning('页面初始化失败，将以默认页面进入编辑器', 3)
        this.$router.push({ path: `/DisPlayEditor/${uid}` })
      })
    },
  }
}
</script>

<style lang="less" scoped>
.template-market {
  min-height: 100vh;
  padding: 24px;
  background: #0a0e1a;
}

.market-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;

  .header-left {
    .market-title {
      color: #e8eaed;
      font-size: 24px;
      font-weight: 600;
      margin: 0 0 8px 0;
    }
    .market-subtitle {
      color: #8892a4;
      font-size: 14px;
      margin: 0;
    }
  }

  .header-right {
    .search-input {
      width: 320px;

      ::v-deep .ant-input {
        background: #141a2e;
        border-color: #2a3052;
        color: #e8eaed;
        border-radius: 8px;
        &:focus, &:hover {
          border-color: #00d4ff;
          box-shadow: 0 0 0 2px rgba(0, 212, 255, 0.1);
        }
      }
      ::v-deep .ant-input-prefix {
        color: #8892a4;
      }
      ::v-deep .ant-input-search-icon {
        color: #00d4ff;
      }
    }
  }
}

.category-bar {
  margin-bottom: 24px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;

  .cat-tag {
    cursor: pointer;
    padding: 4px 14px;
    font-size: 13px;
    border-radius: 16px;
    background: #141a2e;
    border: 1px solid #2a3052;
    color: #8892a4;
    transition: all 0.3s ease;

    &:hover {
      border-color: #00d4ff;
      color: #00d4ff;
    }

    &.cat-active {
      background: rgba(0, 212, 255, 0.1);
      border-color: #00d4ff;
      color: #00d4ff;
      font-weight: 500;
    }

    .cat-count {
      margin-left: 4px;
      opacity: 0.7;
      font-size: 11px;
    }

    ::v-deep .ant-tag {
      margin: 0;
    }
  }
}

.template-grid {
  min-height: 400px;
}

.template-card {
  background: #141a2e;
  border: 1px solid #2a3052;
  border-radius: 10px;
  overflow: hidden;
  transition: all 0.3s ease;

  &:hover {
    border-color: #00d4ff;
    transform: translateY(-4px);
    box-shadow: 0 8px 25px rgba(0, 212, 255, 0.15);
  }

  ::v-deep .ant-card-body {
    padding: 0;
    background: #141a2e;
  }

  ::v-deep .ant-card-meta {
    padding: 14px 16px 8px;
  }

  ::v-deep .ant-card-meta-title {
    color: #e8eaed;
    font-size: 15px;
    font-weight: 500;
    margin-bottom: 8px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  ::v-deep .ant-card-meta-description {
    color: #8892a4;
  }

  ::v-deep .ant-card-actions {
    background: #141a2e;
    border-top: 1px solid #2a3052;
    margin: 0;

    > li {
      margin: 8px 0;

      > span {
        color: #8892a4;
        font-size: 13px;
        transition: color 0.3s ease;

        &:hover {
          color: #00d4ff;
        }

        &.action-use {
          color: #00d4ff;
          &:hover {
            color: #33ddff;
          }
        }
        &.action-edit {
          color: #722ed1;
          &:hover {
            color: #9254de;
          }
        }
      }
    }
  }
}

.card-cover {
  height: 140px;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;

  &.cover-substation { background: linear-gradient(135deg, #0a2a4a 0%, #0d3b6e 100%); }
  &.cover-generation { background: linear-gradient(135deg, #0a3a1a 0%, #0d5e2e 100%); }
  &.cover-industrial { background: linear-gradient(135deg, #3a2a0a 0%, #5e3d0d 100%); }
  &.cover-integrated { background: linear-gradient(135deg, #1a0a3a 0%, #2e0d5e 100%); }
  &.cover-cabinet { background: linear-gradient(135deg, #3a0a2a 0%, #5e0d3e 100%); }
  &.cover-mobile { background: linear-gradient(135deg, #0a2a3a 0%, #0d4e5e 100%); }

  .cover-icon {
    font-size: 48px;
    color: rgba(255, 255, 255, 0.2);
    transition: all 0.3s ease;
  }

  .cover-overlay {
    position: absolute;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(0, 0, 0, 0.5);
    opacity: 0;
    transition: opacity 0.3s ease;

    .overlay-text {
      color: #fff;
      font-size: 14px;
      font-weight: 500;
    }
  }

  .cover-category-tag {
    position: absolute;
    top: 10px;
    right: 10px;
    border-radius: 10px;
    font-size: 11px;
  }
}

.template-card:hover .card-cover {
  .cover-icon {
    color: rgba(255, 255, 255, 0.6);
    transform: scale(1.1);
  }
  .cover-overlay {
    opacity: 1;
  }
}

.card-title {
  color: #e8eaed;
  font-size: 15px;
}

.card-keywords {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;

  .kw-tag {
    font-size: 11px;
    background: rgba(0, 212, 255, 0.1);
    border: 1px solid rgba(0, 212, 255, 0.2);
    color: #00d4ff;
    border-radius: 8px;
    margin: 0;
    padding: 0 6px;
    line-height: 20px;

    &.kw-more {
      background: rgba(136, 146, 164, 0.1);
      border-color: rgba(136, 146, 164, 0.2);
      color: #8892a4;
    }
  }
}

.empty-state {
  margin-top: 80px;
  ::v-deep .ant-empty-description {
    color: #8892a4;
  }
}

.detail-modal {
  ::v-deep .ant-modal-header {
    background: #141a2e;
    border-bottom: 1px solid #2a3052;
    .ant-modal-title {
      color: #e8eaed;
    }
  }
  ::v-deep .ant-modal-close-x {
    color: #8892a4;
    &:hover { color: #e8eaed; }
  }
}

.detail-content {
  color: #c4c9d4;

  .detail-header {
    display: flex;
    align-items: center;
    gap: 16px;
    margin-bottom: 24px;
    padding-bottom: 20px;
    border-bottom: 1px solid #2a3052;

    .detail-icon-wrap {
      width: 64px;
      height: 64px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 30px;
      color: rgba(255, 255, 255, 0.4);

      &.icon-substation { background: linear-gradient(135deg, #0a2a4a, #0d3b6e); }
      &.icon-generation { background: linear-gradient(135deg, #0a3a1a, #0d5e2e); }
      &.icon-industrial { background: linear-gradient(135deg, #3a2a0a, #5e3d0d); }
      &.icon-integrated { background: linear-gradient(135deg, #1a0a3a, #2e0d5e); }
      &.icon-cabinet { background: linear-gradient(135deg, #3a0a2a, #5e0d3e); }
      &.icon-mobile { background: linear-gradient(135deg, #0a2a3a, #0d4e5e); }
    }

    .detail-info {
      h3 {
        color: #e8eaed;
        margin: 0 0 8px 0;
        font-size: 18px;
      }
      .detail-id {
        color: #8892a4;
        font-size: 12px;
        margin-left: 12px;
      }
    }
  }

  .detail-section {
    margin-bottom: 20px;
    h4 {
      color: #8892a4;
      font-size: 13px;
      margin-bottom: 12px;
      i { margin-right: 6px; }
    }
    .detail-tags {
      display: flex;
      flex-wrap: wrap;
      gap: 6px;
    }
    .detail-desc {
      color: #c4c9d4;
      font-size: 14px;
      line-height: 1.8;
      margin: 0;
    }
  }

  .detail-actions {
    display: flex;
    gap: 12px;
    margin-top: 24px;
    padding-top: 20px;
    border-top: 1px solid #2a3052;
  }
}
</style>
