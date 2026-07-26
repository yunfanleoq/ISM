<template>
  <div class="edt-wrap">
    <div class="edt-head">
      <span class="edt-title">{{ $t('displayConfig.RuntimePagesTitle') }}</span>
      <a-button size="small" type="link" @click="refresh" :loading="loading">刷新</a-button>
    </div>
    <div v-if="previewHint" class="edt-hint">{{ $t('displayConfig.RuntimePreviewContext') }}：{{ previewHint }}</div>
    <div v-if="tplSummary" class="edt-tpl">
      <div class="edt-tpl-row" v-for="row in tplSummary" :key="row.key">
        <span class="edt-tpl-k">{{ row.label }}</span>
        <span class="edt-tpl-v" :title="row.pageName">{{ row.pageName || '未配置' }}</span>
      </div>
    </div>
    <div class="edt-body">
      <div v-if="loading" class="edt-empty">加载中…</div>
      <div v-else-if="!treeData.length" class="edt-empty">暂无设备树</div>
      <a-tree
        v-else
        :tree-data="treeData"
        :selected-keys="selectedKeys"
        :expanded-keys="expandedKeys"
        @expand="onExpand"
        @select="onSelect"
      />
    </div>
    <div class="edt-actions" v-if="selectedNode">
      <div class="edt-sel">当前：{{ selectedNode.title }}（{{ kindLabel(selectedNode.kind) }}）</div>
      <a-space direction="vertical" style="width:100%">
        <a-button type="primary" block size="small" @click="createOrOpenTemplate" :loading="acting">
          {{ hasTemplateForSelected ? $t('displayConfig.EditRuntimeTemplate') : '创建该层级模板' }}
        </a-button>
        <a-button block size="small" @click="bindCurrentAsTemplate" :loading="acting" :disabled="!selectPageUuid">
          绑定当前页为该层级模板
        </a-button>
        <a-button block size="small" @click="unbindTemplate" :loading="acting" :disabled="!hasTemplateForSelected">
          解除该层级模板角色
        </a-button>
      </a-space>
    </div>
  </div>
</template>

<script>
import { getMonitorTree } from '@/services/device'
import {
  DisplayModelPageAdd,
  DisplayModelPageBindTemplate,
  displayModelTemplateMap,
} from '@/services/displayModel'
import { mapActions, mapState } from 'vuex'
import store from '@/store'
import Cookie from 'js-cookie'
import {
  buildNavContextForNode,
  buildNavTreeIndex,
  resolveTemplatePageIdForKind,
} from '@/pages/ISMDisPlay/utils/navTreeIndex'
import {
  normalizeRootNodes,
  countDevicesInSubtree,
} from '@/pages/ISMDisPlay/utils/monitorTreeTransform'

const KIND_LABEL = {
  home: '主页',
  root: '主页',
  organization: '组织',
  device: '设备',
}

export default {
  name: 'ISMEditorDeviceTree',
  i18n: require('../../i18n/language'),
  data() {
    return {
      loading: false,
      acting: false,
      treeData: [],
      selectedKeys: [],
      expandedKeys: [],
      selectedNode: null,
      templateMap: null,
      navIndex: null,
      previewHint: '',
    }
  },
  computed: {
    ...mapState({
      selectPageUuid: state => store.state.ISMDisPlayEditorTool.selectPageUuid,
      PCPageList: state => store.state.ISMDisPlayEditorTool.PCPageList,
      editorRuntimePreview: state => store.state.ISMDisPlayEditorTool.editorRuntimePreview,
    }),
    modelId() {
      return this.$route.params.uid || ''
    },
    projectUuid() {
      return Cookie.get('ProjectUuid') || sessionStorage.getItem('ProjectUuid') || ''
    },
    hasTemplateForSelected() {
      if (!this.selectedNode || !this.templateMap) return false
      return !!this.resolveTemplatePageId(this.selectedNode)
    },
    tplSummary() {
      const m = this.templateMap
      if (!m) return []
      const findName = pid => {
        if (!pid) return ''
        const p = (this.PCPageList || []).find(x => x.pageUuid === pid)
        return (p && p.title) || pid.slice(0, 8) + '…'
      }
      const rows = [
        { key: 'home', label: '首页模板', pageName: findName(m.home) },
        { key: 'deviceList', label: '设备列表模板', pageName: findName(m.deviceList) },
        { key: 'datapointList', label: '点位列表模板', pageName: findName(m.datapointList) },
      ]
      return rows
    },
  },
  mounted() {
    this.refresh()
  },
  methods: {
    ...mapActions('ISMDisPlayEditorTool', [
      'selectLayerDataStruct',
      'selectEditorRuntimePreview',
      'getLayerDataStruct',
    ]),
    kindLabel(k) {
      return KIND_LABEL[k] || k || '-'
    },
    async refresh() {
      this.loading = true
      try {
        // 先取得映射再生成树，确保虚拟节点首次渲染即可显示模板状态。
        await this.fetchTemplateMap()
        await this.fetchTree()
      } finally {
        this.loading = false
      }
    },
    async fetchTemplateMap() {
      if (!this.modelId) return
      try {
        const res = await displayModelTemplateMap({ muid: this.modelId })
        if (res && res.data && res.data.code === 0) {
          this.templateMap = res.data.map || null
          store.commit('ISMDisPlayEditorTool/setNavTemplateMap', this.templateMap)
        }
      } catch (e) {
        console.warn('[EditorDeviceTree] templateMap', e && e.message)
      }
    },
    async fetchTree() {
      try {
        const res = await getMonitorTree({}, { headers: { ProjectUuid: this.projectUuid } })
        if (res && res.data && res.data.code === 0 && Array.isArray(res.data.list)) {
          this.treeData = this.buildTree(res.data.list)
          this.navIndex = buildNavTreeIndex(this.treeData)
          store.commit('ISMDisPlayEditorTool/setNavTreeIndex', this.navIndex)
          if (this.treeData.length) {
            this.expandedKeys = [this.treeData[0].key]
          }
        } else {
          this.treeData = []
        }
      } catch (e) {
        this.treeData = []
        console.warn('[EditorDeviceTree] tree', e && e.message)
      }
    },
    buildTree(nodes) {
      const mapNode = (node, depth = 1, fallbackKey = '') => {
        const v = node.value || {}
        const name = node.text || v.Name || '未命名'
        const type = Number(v.type)
        const sid = v.sid
        const businessKey = String(v.uuid || node.key || (sid != null ? sid : fallbackKey))
        const children = (node.children || []).map((child, index) =>
          mapNode(child, depth + 1, `${businessKey}-${index}`),
        ).filter(Boolean)
        const kind = type === 1 ? 'device' : (depth === 1 ? 'root' : 'organization')
        const key = `${kind}-${businessKey}`
        const tplId = this.resolveTemplatePageId({ kind })
        const title = tplId ? `${name} ✓` : name
        return {
          key,
          title,
          label: name,
          name,
          kind,
          layer: kind === 'root' ? 'home' : (kind === 'device' ? 'device' : 'organization'),
          sid,
          uuid: v.uuid || node.key || '',
          modelUuid: v.muid || '',
          muid: v.muid || '',
          status: v.status || v.Status || 'off',
          treeDepth: depth,
          type,
          count: countDevicesInSubtree(children) || undefined,
          raw: v,
          children: children.length ? children : undefined,
          isLeaf: type === 1 || !children.length,
        }
      }
      return normalizeRootNodes(nodes || []).map((node, index) =>
        mapNode(node, 1, `root-${index}`),
      ).filter(Boolean)
    },
    resolveTemplatePageId(node) {
      const map = this.templateMap
      if (!map || !node) return ''
      const kind = node.kind === 'root' ? 'home' : node.kind
      if (kind === 'home') return map.home || ''
      return kind === 'device' ? (map.datapointList || '') : (map.deviceList || '')
    },
    templateKindOf(node) {
      if (!node) return ''
      if (node.kind === 'root') return 'home'
      return node.kind === 'device' ? 'datapointList' : 'deviceList'
    },
    templateModelUuidOf(node) {
      return ''
    },
    onExpand(keys) {
      this.expandedKeys = keys
    },
    async onSelect(keys, info) {
      this.selectedKeys = keys
      const n = info && info.node && info.node.dataRef
      this.selectedNode = n || null
      if (!n) return
      this.previewHint = `${n.title} / ${this.kindLabel(n.kind)}`
      const navContext = buildNavContextForNode(n, this.navIndex)
      const pageId = resolveTemplatePageIdForKind(
        this.templateMap,
        navContext.kind || n.kind,
        navContext.modelUuid || n.modelUuid || '',
      )
      if (!pageId) return
      this.acting = true
      try {
        const resolved = await this.selectEditorRuntimePreview({
          pageUuid: pageId,
          virtualKey: keys[0] || n.key,
          virtualTitle: n.label || n.title,
          navContext,
        })
        if (!resolved) {
          this.$message.error('运行态页面预览加载失败')
        } else {
          document.title = `${resolved.name || n.title} | ${this.$t('displayConfig.RuntimePreviewReadonly')}`
        }
      } finally {
        this.acting = false
      }
    },
    openPage(pageUuid) {
      if (!pageUuid) return
      store.commit('ISMDisPlayEditorTool/setEditorRuntimePreview', null)
      store.commit('ISMDisPlayEditorTool/setNavContext', null)
      const page = (this.PCPageList || []).find(item => item.pageUuid === pageUuid)
      this.selectLayerDataStruct({
        pageType: 1,
        pageUuid,
        title: page ? page.title : '',
      })
      if (page) document.title = `${page.AppName || ''} | ${page.title}`
    },
    defaultPageName(node) {
      const kind = this.templateKindOf(node)
      const labels = {
        home: '首页模板',
        deviceList: '设备列表模板',
        datapointList: '点位列表模板',
      }
      return labels[kind] || '模板页'
    },
    async createOrOpenTemplate() {
      const node = this.selectedNode
      if (!node) return
      const exist = this.resolveTemplatePageId(node)
      if (exist) {
        this.openPage(exist)
        this.$message.success('已打开该层级模板页')
        return
      }
      this.acting = true
      try {
        const kind = this.templateKindOf(node)
        const modelUuid = this.templateModelUuidOf(node)
        const params = {
          modelUuid: this.modelId,
          name: this.defaultPageName(node),
          size: '1',
          pageType: 1,
          isLogin: 0,
          templateKind: kind,
          templateModelUuid: modelUuid,
        }
        const res = await DisplayModelPageAdd(params)
        const code = res && res.data && res.data.code
        const pageId = res && res.data && res.data.pageId
        if ((code === 4002 || code === 4001) && pageId) {
          await this.reloadPagesAndOpen(pageId)
          this.$message.success(code === 4001 ? '模板已存在，已打开' : '模板页已创建')
        } else if (code === 4002 || code === 200) {
          await this.reloadPagesAndOpen(pageId)
          this.$message.success('模板页已创建')
        } else {
          this.$message.error('创建模板失败 code=' + code)
        }
        await this.fetchTemplateMap()
      } catch (e) {
        this.$message.error('创建模板失败')
      } finally {
        this.acting = false
      }
    },
    async bindCurrentAsTemplate() {
      const node = this.selectedNode
      if (!node || !this.selectPageUuid) return
      this.acting = true
      try {
        const kind = this.templateKindOf(node)
        const modelUuid = this.templateModelUuidOf(node)
        let res = await DisplayModelPageBindTemplate({
          modelUuid: this.modelId,
          pageId: this.selectPageUuid,
          templateKind: kind,
          templateModelUuid: modelUuid,
          force: false,
        })
        let code = res && res.data && res.data.code
        if (code === 4001) {
          const ok = await new Promise(resolve => {
            this.$confirm({
              title: '该层级已有模板页',
              content: '是否覆盖为当前页？',
              onOk: () => resolve(true),
              onCancel: () => resolve(false),
            })
          })
          if (!ok) return
          res = await DisplayModelPageBindTemplate({
            modelUuid: this.modelId,
            pageId: this.selectPageUuid,
            templateKind: kind,
            templateModelUuid: modelUuid,
            force: true,
          })
          code = res && res.data && res.data.code
        }
        if (code === 200 || code === 0) {
          this.$message.success('已绑定为层级模板')
          await this.fetchTemplateMap()
        } else {
          this.$message.error('绑定失败 code=' + code)
        }
      } catch (e) {
        this.$message.error('绑定失败')
      } finally {
        this.acting = false
      }
    },
    async unbindTemplate() {
      const node = this.selectedNode
      const pageId = this.resolveTemplatePageId(node)
      if (!pageId) return
      this.acting = true
      try {
        const res = await DisplayModelPageBindTemplate({
          modelUuid: this.modelId,
          pageId,
          templateKind: '',
          templateModelUuid: '',
          force: false,
        })
        const code = res && res.data && res.data.code
        if (code === 200 || code === 0) {
          this.$message.success('已解除模板角色')
          await this.fetchTemplateMap()
        } else {
          this.$message.error('解绑失败')
        }
      } catch (e) {
        this.$message.error('解绑失败')
      } finally {
        this.acting = false
      }
    },
    reloadPagesAndOpen(pageId) {
      const uid = this.modelId
      return new Promise(resolve => {
        this.getLayerDataStruct({
          uuid: uid,
          metaOnly: true,
          cb: () => {
            if (pageId) this.openPage(pageId)
            resolve()
          },
        })
      })
    },
  },
}
</script>

<style scoped>
.edt-wrap {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 8px;
  background: #fff;
}
.edt-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 6px;
}
.edt-title {
  font-weight: 600;
  font-size: 13px;
}
.edt-hint {
  font-size: 12px;
  color: #1890ff;
  margin-bottom: 6px;
  word-break: break-all;
}
.edt-tpl {
  border: 1px solid #e8e8e8;
  border-radius: 4px;
  padding: 6px 8px;
  margin-bottom: 8px;
  font-size: 12px;
  max-height: 120px;
  overflow: auto;
}
.edt-tpl-row {
  display: flex;
  justify-content: space-between;
  gap: 8px;
  line-height: 1.6;
}
.edt-tpl-k { color: #666; flex-shrink: 0; }
.edt-tpl-v { color: #333; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.edt-body {
  flex: 1;
  overflow: auto;
  border: 1px solid #f0f0f0;
  border-radius: 4px;
  padding: 4px;
  min-height: 200px;
}
.edt-empty {
  color: #999;
  padding: 16px;
  text-align: center;
}
.edt-actions {
  margin-top: 8px;
  border-top: 1px solid #f0f0f0;
  padding-top: 8px;
}
.edt-sel {
  font-size: 12px;
  margin-bottom: 6px;
  color: #333;
  word-break: break-all;
}
</style>
