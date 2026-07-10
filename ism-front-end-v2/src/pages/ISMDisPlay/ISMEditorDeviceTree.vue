<template>
  <div class="edt-wrap">
    <div class="edt-head">
      <span class="edt-title">设备树 / 层级模板</span>
      <a-button size="small" type="link" @click="refresh" :loading="loading">刷新</a-button>
    </div>
    <div v-if="previewHint" class="edt-hint">预览上下文：{{ previewHint }}</div>
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
          创建/打开该层级模板
        </a-button>
        <a-button block size="small" @click="bindCurrentAsTemplate" :loading="acting" :disabled="!selectPageUuid">
          绑定当前页为该层级模板
        </a-button>
        <a-button block size="small" @click="unbindTemplate" :loading="acting" :disabled="!hasTemplateForSelected">
          解除该层级模板角色
        </a-button>
        <a-checkbox v-model="useModelOverride" v-if="selectedNode.kind==='device'">
          按物模型覆盖（{{ shortMuid }}）
        </a-checkbox>
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

const KIND_LABEL = {
  home: '主页',
  root: '主页',
  zone: '区域',
  room: '机房',
  cabinet: '机柜',
  device: '设备',
}

export default {
  name: 'ISMEditorDeviceTree',
  data() {
    return {
      loading: false,
      acting: false,
      treeData: [],
      selectedKeys: [],
      expandedKeys: [],
      selectedNode: null,
      templateMap: null,
      useModelOverride: false,
      previewHint: '',
    }
  },
  computed: {
    ...mapState({
      selectPageUuid: state => store.state.ISMDisPlayEditorTool.selectPageUuid,
      PCPageList: state => store.state.ISMDisPlayEditorTool.PCPageList,
    }),
    modelId() {
      return this.$route.params.uid || ''
    },
    projectUuid() {
      return Cookie.get('ProjectUuid') || sessionStorage.getItem('ProjectUuid') || ''
    },
    shortMuid() {
      const m = (this.selectedNode && this.selectedNode.modelUuid) || ''
      return m ? m.slice(0, 8) + '…' : '无模型'
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
        { key: 'home', label: '主页', pageName: findName(m.home) },
        { key: 'zone', label: '区域', pageName: findName(m.zone) },
        { key: 'room', label: '机房', pageName: findName(m.room) },
        { key: 'cabinet', label: '机柜', pageName: findName(m.cabinet) },
        { key: 'device', label: '设备(通用)', pageName: findName(m.deviceDefault) },
      ]
      const by = m.deviceByModel || {}
      Object.keys(by).forEach(muid => {
        rows.push({
          key: 'dev-' + muid,
          label: '设备覆盖 ' + muid.slice(0, 6),
          pageName: findName(by[muid]),
        })
      })
      return rows
    },
  },
  mounted() {
    this.refresh()
  },
  methods: {
    ...mapActions('ISMDisPlayEditorTool', ['selectLayerDataStruct', 'getLayerDataStruct']),
    kindLabel(k) {
      return KIND_LABEL[k] || k || '-'
    },
    async refresh() {
      this.loading = true
      try {
        await Promise.all([this.fetchTemplateMap(), this.fetchTree()])
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
        }
      } catch (e) {
        console.warn('[EditorDeviceTree] templateMap', e && e.message)
      }
    },
    async fetchTree() {
      try {
        const res = await getMonitorTree({ headers: { ProjectUuid: this.projectUuid } })
        if (res && res.data && res.data.code === 0 && Array.isArray(res.data.list)) {
          this.treeData = this.buildTree(res.data.list)
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
      const mapNode = (node) => {
        const v = node.value || {}
        const name = node.text || v.Name || '未命名'
        const type = v.type
        const sid = v.sid
        const pid = v.pid
        const children = (node.children || []).map(mapNode).filter(Boolean)
        let kind = 'zone'
        if (type === 1) kind = 'device'
        else if (sid === 1) kind = 'root'
        else if (pid === 1) kind = 'zone'
        else if (children.length && children.every(c => c.kind === 'device')) kind = 'cabinet'
        else if (/配电室|机房|房间/.test(name)) kind = 'room'
        else kind = 'zone'
        const key = node.key || `${kind}-${sid}`
        const tplId = this.resolveTemplatePageId({ kind, modelUuid: v.muid || '' })
        const title = tplId ? `${name} ✓` : name
        return {
          key,
          title,
          kind,
          sid,
          uuid: v.uuid || node.key || '',
          modelUuid: v.muid || '',
          children: children.length ? children : undefined,
          isLeaf: type === 1 || !children.length,
        }
      }
      return (nodes || []).map(mapNode).filter(Boolean)
    },
    resolveTemplatePageId(node) {
      const map = this.templateMap
      if (!map || !node) return ''
      const kind = node.kind === 'root' ? 'home' : node.kind
      if (kind === 'home') return map.home || ''
      if (kind === 'zone') return map.zone || ''
      if (kind === 'room') return map.room || map.zone || ''
      if (kind === 'cabinet') return map.cabinet || ''
      if (kind === 'device') {
        const muid = node.modelUuid || ''
        if (this.useModelOverride && muid && map.deviceByModel && map.deviceByModel[muid]) {
          return map.deviceByModel[muid]
        }
        if (muid && map.deviceByModel && map.deviceByModel[muid]) return map.deviceByModel[muid]
        return map.deviceDefault || ''
      }
      return ''
    },
    templateKindOf(node) {
      if (!node) return ''
      if (node.kind === 'root') return 'home'
      return node.kind
    },
    templateModelUuidOf(node) {
      if (!node || node.kind !== 'device') return ''
      if (this.useModelOverride) return node.modelUuid || ''
      return ''
    },
    onExpand(keys) {
      this.expandedKeys = keys
    },
    onSelect(keys, info) {
      this.selectedKeys = keys
      const n = info && info.node && info.node.dataRef
      this.selectedNode = n || null
      if (!n) return
      this.previewHint = `${n.title} / ${this.kindLabel(n.kind)}`
      try {
        store.commit('ISMDisPlayEditorTool/setNavContext', {
          sid: n.sid,
          uuid: n.uuid,
          name: n.title,
          kind: n.kind === 'root' ? 'home' : n.kind,
          modelUuid: n.modelUuid || '',
          childDevices: [],
        })
      } catch (e) { /* ignore */ }
      const pageId = this.resolveTemplatePageId(n)
      if (pageId) this.openPage(pageId)
    },
    openPage(pageUuid) {
      if (!pageUuid) return
      this.selectLayerDataStruct({
        pageType: 1,
        pageUuid,
      })
    },
    defaultPageName(node) {
      const kind = this.templateKindOf(node)
      const labels = { home: '主页模板', zone: '区域模板', room: '机房模板', cabinet: '机柜模板', device: '设备模板' }
      let name = labels[kind] || '模板页'
      const muid = this.templateModelUuidOf(node)
      if (kind === 'device' && muid) name = `设备模板-${muid.slice(0, 8)}`
      return name
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
