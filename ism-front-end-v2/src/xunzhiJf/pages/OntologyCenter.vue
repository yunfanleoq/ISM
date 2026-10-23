<template>
  <div ref="center" class="ont-center">
    <div class="ont-bar">
      <a-radio-group v-model="mainTab" button-style="solid" @change="onTab">
        <a-radio-button value="schema">类型星图</a-radio-button>
        <a-radio-button value="system">系统实例</a-radio-button>
        <a-radio-button value="world">世界实例</a-radio-button>
      </a-radio-group>
      <a-radio-group v-model="viewMode" size="small" @change="onViewMode">
        <a-radio-button value="whole">整体关系</a-radio-button>
        <a-radio-button value="drill">逐层展开</a-radio-button>
      </a-radio-group>
      <a-radio-group v-model="layer" size="small" style="margin-left: 12px" @change="onLayer">
        <a-radio-button value="all">全部</a-radio-button>
        <a-radio-button value="physical">物理</a-radio-button>
        <a-radio-button value="business">业务</a-radio-button>
        <a-radio-button value="system">系统</a-radio-button>
        <a-radio-button value="bridge">桥接</a-radio-button>
        <a-radio-button value="gap">缺口</a-radio-button>
      </a-radio-group>
      <span class="ont-hint">{{ hint }}</span>
      <a-button v-if="mainTab !== 'schema'" type="primary" size="small" :loading="busy" @click="compileCurrent">编译</a-button>
      <a-checkbox v-if="mainTab === 'world'" v-model="showPoints" style="margin-left: 8px">显示测点</a-checkbox>
      <span class="ont-spacer" />
      <a-button size="small" @click="$router.push('/xunzhi/jf/agent')">智能体</a-button>
      <a-button size="small" @click="toggleFull">{{ isFull ? '退出全屏' : '全屏' }}</a-button>
      <a-button size="small" @click="goMain">返回主界面</a-button>
      <a-button size="small" @click="goAdmin">返回后台</a-button>
    </div>
    <div class="ont-body">
      <button v-if="!sideOpen" type="button" class="ont-rail" title="展开目录" @click="sideOpen = true">目录</button>
      <aside v-show="sideOpen" class="ont-side">
        <div class="ont-foldbar">
          <span>目录</span>
          <button type="button" class="ont-fold" title="收起目录" @click="sideOpen = false">‹</button>
        </div>
        <div class="ont-sec">实体 {{ catalogEntities.length }}</div>
        <button v-for="ent in catalogEntities" :key="ent.type" type="button" class="ont-row" :class="{ on: selected && selected.id === ent.type }" @click="pickEntity(ent)">
          <i :style="{ background: colorOf(ent.layer) }" />
          <span>{{ ent.label }}</span>
          <em>{{ ent.maturity }}</em>
        </button>
        <div class="ont-sec">关系 {{ catalogRelations.length }}</div>
        <button v-for="rel in catalogRelations" :key="rel.relation" type="button" class="ont-row" @click="pickRelation(rel)">
          <i :style="{ background: colorOf(rel.layer || 'bridge') }" />
          <span>{{ rel.label }}</span>
          <em>{{ rel.relation }}</em>
        </button>
        <div class="ont-sec">动作 {{ (schema.action_types || []).length }}</div>
        <button v-for="act in schema.action_types || []" :key="act.id" type="button" class="ont-row" @click="pickAction(act)">
          <i style="background:#3ddc97" />
          <span>{{ act.label }}</span>
          <em>{{ act.implemented ? '可执行' : '未落地' }}</em>
        </button>
        <div class="ont-sec">函数 / 界面</div>
        <div v-for="fn in schema.function_types || []" :key="fn.id" class="ont-mini">{{ fn.label }} · {{ fn.id }}</div>
        <div v-for="iface in schema.interface_types || []" :key="iface.id" class="ont-mini">{{ iface.label }} · 允许 {{ (iface.allowed_actions || []).length }} 个动作</div>
      </aside>
      <main class="ont-main">
        <div class="ont-crumb">
          <button type="button" @click="setStack([])">总览</button>
          <button v-for="(crumb, index) in stack" :key="crumb.id + index" type="button" @click="setStack(stack.slice(0, index + 1))">{{ crumb.label }}</button>
          <span>{{ viewMode === 'whole' && !stack.length ? '自上而下是组织层级。点图例可开关某一层，点球看它的上下级。' : '点一个球，从这里再展开一层。点路径可以退回。' }}</span>
        </div>
        <div class="ont-stage">
          <OntologyGraph3D
            :nodes="graph.nodes"
            :edges="graph.edges"
            :layout="graph.layout || 'ring'"
            :focus-id="focusId"
            @select="openNode"
          />
          <div class="ont-legend">
            <button
              v-for="item in legendItems"
              :key="item.id"
              type="button"
              :class="{ off: !showLayer[item.id] }"
              @click="toggleLayer(item.id)"
            >
              <i :style="{ background: item.color }" />
              {{ item.label }}
            </button>
          </div>
        </div>
      </main>
      <aside v-show="detailOpen" class="ont-detail">
        <div class="ont-foldbar">
          <button type="button" class="ont-fold" title="收起详情" @click="detailOpen = false">›</button>
          <span>详情</span>
        </div>
        <a-card size="small" :title="detailTitle">
          <template v-if="detail">
            <p>{{ detail.label }} <code>{{ detail.type || detail.id || detail.relation }}</code></p>
            <p v-if="detail.layer">层：{{ detail.layer }} · 成熟度：{{ detail.maturity || '—' }}</p>
            <p v-if="detail.source">来源：{{ detail.source }}</p>
            <p v-if="detail.description">{{ detail.description }}</p>
            <p v-if="detail.workflow">工作流：{{ detail.workflow }}</p>
            <p v-if="detail.raw">{{ detail.raw.node_type }} · {{ detail.raw.node_key }}</p>
          </template>
          <p v-else class="muted">从总览点一个世界，再一层层点进去。</p>
        </a-card>
        <node-live :live="live" :loading="liveBusy" />
        <a-card size="small" title="从这里展开" style="margin-top: 8px">
          <button v-for="item in around" :key="item.id" type="button" class="ont-row" @click="openNode(item)">
            <i :style="{ background: item.color || colorOf(item.layer) }" />
            <span>{{ item.label }}</span>
          </button>
          <p v-if="!around.length" class="muted">这一层没有可继续展开的对象。</p>
        </a-card>
        <a-card size="small" title="公理" style="margin-top: 8px">
          <div v-if="!(schema.axioms || []).length" class="muted">无</div>
          <div v-for="ax in shownAxioms" :key="ax.id" class="hit" :class="ax.severity">
            {{ ax.severity }} · {{ ax.label }}
          </div>
        </a-card>
        <a-card v-if="mainTab !== 'schema'" size="small" title="本次检查" style="margin-top: 8px">
          <div v-if="!lint.hits || !lint.hits.length" class="muted">还没有命中。编译世界图或系统图后再看。</div>
          <div v-for="(hit, i) in (lint.hits || []).slice(0, 12)" :key="i" class="hit" :class="hit.severity">
            {{ hit.severity }} · {{ hit.message || hit.label || hit.axiom_id }}
          </div>
        </a-card>
        <a-card v-if="mainTab === 'world'" size="small" title="动作" style="margin-top: 8px">
          <a-select v-model="actionId" style="width: 100%" :options="actionOptions" />
          <a-input v-model="alarmId" placeholder="alarm_id" style="margin-top: 8px" />
          <a-input v-model="titleText" placeholder="标题" style="margin-top: 8px" />
          <a-button type="primary" block style="margin-top: 8px" :loading="busy" @click="runAction">执行</a-button>
          <pre v-if="actionMsg" class="msg">{{ actionMsg }}</pre>
        </a-card>
        <a-card v-if="run" size="small" title="告警研判" style="margin-top: 8px">
          <div v-for="step in run.steps || []" :key="step.step_key">{{ step.step_key }} · {{ step.status }} · {{ step.message }}</div>
        </a-card>
        <a-button v-if="mainTab === 'world'" block style="margin-top: 8px" :loading="busy" @click="runAlarm">跑告警研判</a-button>
      </aside>
      <button v-if="!detailOpen" type="button" class="ont-rail" title="展开详情" @click="detailOpen = true">详情</button>
    </div>
  </div>
</template>

<script>
import OntologyGraph3D from '@/xunzhiJf/ontology/OntologyGraph3D.vue'
import NodeLive from '@/xunzhiJf/ontology/NodeLive.vue'
import { applyLegend, drillInstances, drillSchema, LAYER_COLOR, LAYER_LABEL, layerColor, wholeInstances, wholeSchema } from '@/xunzhiJf/ontology/graphModel'
import { getHomeDashboardPath } from '@/config/homeDashboard'
import { resolveXunzhiJfProjectUuid } from '@/xunzhiJf/projectContext'
import { METHOD, request } from '@/utils/request'

export default {
  name: 'OntologyCenter',
  components: { OntologyGraph3D, NodeLive },
  data() {
    return {
      mainTab: 'schema',
      viewMode: 'whole',
      layer: 'all',
      focusId: '',
      sideOpen: true,
      detailOpen: true,
      isFull: false,
      showLayer: { physical: true, business: true, system: true, bridge: true, gap: false },
      schema: {},
      nodes: [],
      edges: [],
      lint: {},
      graph: { nodes: [], edges: [], hidden: 0, total: 0 },
      stack: [],
      detail: null,
      live: { fields: [], rows: [], logs: [], note: '点具体的设备、测点或告警，这里列出它的生产数据和日志。' },
      liveBusy: false,
      liveTicket: 0,
      busy: false,
      showPoints: false,
      projectUuid: '',
      projectName: '',
      actionId: 'AckAlarm',
      alarmId: '',
      titleText: '',
      actionMsg: '',
      run: null,
      actionOptions: [
        { value: 'AckAlarm', label: '确认告警' },
        { value: 'CreateWorkOrder', label: '由告警开单' },
        { value: 'VerifyEdge', label: '确认边' },
        { value: 'ApproveControl', label: '批准控制（应拒绝）' },
      ],
    }
  },
  computed: {
    hint() {
      if (this.viewMode === 'whole' && !this.stack.length) {
        return '三个世界各自是一棵组织树，上面管下面。左右拖动能看全。点图例可开关，点球看上下级。'
      }
      if (!this.stack.length) return '先点一个世界。物理、业务、系统各是一个入口。'
      return '当前从「' + this.stack[this.stack.length - 1].label + '」展开。再点周围的球，往下看一层。'
    },
    legendItems() {
      return [
        { id: 'physical', label: '物理', color: LAYER_COLOR.physical },
        { id: 'business', label: '业务', color: LAYER_COLOR.business },
        { id: 'system', label: '系统', color: LAYER_COLOR.system },
        { id: 'bridge', label: '桥接', color: LAYER_COLOR.bridge },
        { id: 'gap', label: '缺口', color: LAYER_COLOR.gap },
      ]
    },
    detailTitle() {
      if (!this.detail) return '详情'
      if (this.detail.kind === 'instance') return '实例'
      if (this.detail.relation) return '关系'
      if (this.detail.workflow) return '动作'
      return '类型'
    },
    catalogEntities() {
      return (this.schema.entity_types || []).filter(e => this.layer === 'all' || e.layer === this.layer)
    },
    catalogRelations() {
      return (this.schema.relation_types || []).filter(r => this.layer === 'all' || r.layer === this.layer)
    },
    shownAxioms() {
      return (this.schema.axioms || []).filter(a => this.layer === 'all' || a.layer === this.layer).slice(0, 16)
    },
    around() {
      const id = this.detail && (this.detail.id || this.detail.type)
      const nodes = (this.graph.nodes || []).filter(n => n.id !== id)
      if (this.viewMode === 'whole' && !this.stack.length) {
        return nodes.filter(n => n.kind === 'layer' || n.kind === 'relation').slice(0, 20)
      }
      return nodes.slice(0, 36)
    },
  },
  mounted() {
    this.boot()
    document.addEventListener('fullscreenchange', this.onFullChange)
  },
  beforeDestroy() {
    document.removeEventListener('fullscreenchange', this.onFullChange)
  },
  methods: {
    colorOf(layer) { return layerColor(layer) },
    goMain() {
      const path = getHomeDashboardPath(this.$store)
      this.$router.push(path && path !== '/AppRun/' ? path : '/xunzhi/jf/viz/bjga')
    },
    goAdmin() {
      this.$router.push('/Project')
    },
    async boot() {
      await this.loadSchema()
      this.paintSchema()
      try {
        const uuid = await resolveXunzhiJfProjectUuid()
        this.projectUuid = uuid || ''
        this.projectName = sessionStorage.getItem('xunzhiJf.projectName') || ''
      } catch (e) { /* 没有项目时类型星图仍然可看 */ }
      this.compileCurrent('system', true)
    },
    async loadSchema() {
      try {
        const res = await request('/api/ontology/schema', METHOD.GET, {})
        this.schema = (res.data && res.data.schema) || {}
      } catch (e) {
        this.actionMsg = '类型定义没有读到。请先登录，再打开本体中心。'
      }
    },
    paintSchema() { this.paintVisible() },
    paintInstances() { this.paintVisible() },
    paintVisible() {
      let graph
      if (this.viewMode === 'whole' && !this.stack.length) {
        graph = this.mainTab === 'schema'
          ? wholeSchema(this.schema, this.showLayer)
          : wholeInstances(this.nodes, this.edges, this.showPoints, this.showLayer)
      } else {
        graph = this.mainTab === 'schema'
          ? drillSchema(this.schema, this.stack)
          : drillInstances(this.nodes, this.edges, this.stack, this.showPoints, this.schema)
        graph = applyLegend(graph, this.showLayer)
      }
      this.graph = graph
    },
    setStack(next) {
      this.stack = next
      this.viewMode = next.length ? 'drill' : 'whole'
      this.focusId = next.length ? next[next.length - 1].id : ''
      const last = next[next.length - 1]
      if (last) this.detail = last
      this.paintVisible()
    },
    onViewMode() {
      if (this.viewMode === 'whole') this.stack = []
      this.focusId = ''
      this.paintVisible()
    },
    toggleLayer(id) {
      this.$set(this.showLayer, id, !this.showLayer[id])
      const keys = ['physical', 'business', 'system', 'bridge', 'gap']
      const on = keys.filter(k => this.showLayer[k])
      this.layer = on.length === 1 ? on[0] : 'all'
      this.paintVisible()
    },
    toggleFull() {
      const el = this.$refs.center
      if (!el) return
      if (document.fullscreenElement) document.exitFullscreen()
      else if (el.requestFullscreen) el.requestFullscreen()
    },
    onFullChange() {
      this.isFull = !!document.fullscreenElement
      this.$nextTick(() => window.dispatchEvent(new Event('resize')))
    },
    openNode(node) {
      if (!node) return
      this.viewMode = 'drill'
      this.focusId = node.id || node.type || ''
      if (node.kind === 'type' && node.nodeType === 'Point') {
        this.showPoints = true
        this.paintVisible()
        return
      }
      if (node.kind === 'layer') {
        this.stack = [{ kind: 'layer', id: node.id, label: LAYER_LABEL[node.layer] || node.label, layer: node.layer }]
      } else if (node.kind === 'type') {
        const layer = node.layer
        this.stack = [
          { kind: 'layer', id: 'layer:' + layer, label: LAYER_LABEL[layer] || layer, layer },
          { kind: 'type', id: node.id, label: String(node.label || '').split(' · ')[0], layer, nodeType: node.nodeType },
        ]
      } else if (node.kind === 'entity' || node.kind === 'relation' || (node.type && !node.kind)) {
        const id = node.id || node.type
        const layer = node.layer
        const crumb = { kind: node.kind || 'entity', id, label: node.label, layer }
        const last = this.stack[this.stack.length - 1]
        if (last && (last.kind === 'entity' || last.kind === 'relation') && last.id !== id) {
          this.stack = this.stack.concat([crumb]).slice(-8)
        } else {
          this.stack = [
            { kind: 'layer', id: 'layer:' + layer, label: LAYER_LABEL[layer] || layer, layer },
            crumb,
          ]
        }
      } else {
        const layer = node.layer
        const nodeType = node.nodeType || (node.raw && node.raw.node_type)
        const crumb = { kind: 'instance', id: node.id, label: node.label, layer, nodeType }
        const last = this.stack[this.stack.length - 1]
        if (last && last.kind === 'instance' && last.id !== crumb.id) this.stack = this.stack.concat([crumb]).slice(-8)
        else {
          this.stack = [
            { kind: 'layer', id: 'layer:' + layer, label: LAYER_LABEL[layer] || layer, layer },
            { kind: 'type', id: 'type:' + layer + ':' + nodeType, label: nodeType, layer, nodeType },
            crumb,
          ]
        }
      }
      this.detail = node.raw
        ? Object.assign({ kind: 'instance', id: node.id, label: node.label }, node.raw, node)
        : ((this.schema.entity_types || []).find(e => e.type === (node.id || node.type)) || node)
      this.paintVisible()
      this.loadLive(this.detail)
    },
    onTab() {
      this.stack = []
      if (this.mainTab === 'schema') this.paintVisible()
      else this.loadGraph()
    },
    onLayer() {
      const keys = ['physical', 'business', 'system', 'bridge', 'gap']
      if (this.viewMode === 'whole') {
        if (this.layer === 'all') keys.forEach(k => this.$set(this.showLayer, k, k !== 'gap'))
        else keys.forEach(k => this.$set(this.showLayer, k, k === this.layer))
        this.stack = []
        this.paintVisible()
        return
      }
      if (this.layer === 'all') this.stack = []
      else this.stack = [{ kind: 'layer', id: 'layer:' + this.layer, label: LAYER_LABEL[this.layer] || this.layer, layer: this.layer }]
      if (this.mainTab === 'schema' || this.nodes.length) this.paintVisible()
    },
    pickEntity(ent) {
      this.openNode(Object.assign({ kind: 'entity', id: ent.type }, ent))
    },
    async loadLive(node) {
      const ticket = ++this.liveTicket
      const raw = (node && node.raw) || {}
      const id = String((node && (node.id || node.type)) || raw.node_key || '')
      const typed = raw.node_type || (node && (node.nodeType || node.type)) || ''
      const key = raw.node_key || (id.indexOf(':') > 0 && id.indexOf('type:') !== 0 && id.indexOf('layer:') !== 0 ? id : '')
      if (!typed || (node && (node.kind === 'layer' || node.kind === 'relation' || node.relation))) {
        this.live = { fields: [], rows: [], logs: [], note: '点具体的设备、测点或告警，这里列出它的生产数据和日志。' }
        this.liveBusy = false
        return
      }
      this.liveBusy = true
      try {
        const res = await request('/api/ontology/live', METHOD.GET, {
          project_uuid: this.projectUuid,
          node_type: typed,
          node_key: key,
        })
        if (ticket !== this.liveTicket) return
        this.live = (res.data && res.data.live) || { fields: [], rows: [], logs: [], note: '没有读到生产数据' }
      } catch (e) {
        if (ticket === this.liveTicket) this.live = { fields: [], rows: [], logs: [], note: '生产数据暂时读不到' }
      } finally {
        if (ticket === this.liveTicket) this.liveBusy = false
      }
    },
    pickRelation(rel) {
      this.detail = rel
    },
    pickAction(act) {
      this.detail = act
      this.actionId = act.id
    },
    onGraphSelect(node) {
      this.openNode(node)
    },
    async compileCurrent(kind, silent) {
      const graphKind = kind || this.mainTab
      if (graphKind === 'schema') return
      if (graphKind === 'world' && !this.projectUuid) {
        this.actionMsg = '没有项目 uuid，世界图先不编译。类型星图可以看。'
        return
      }
      this.busy = true
      try {
        const res = await request('/api/ontology/compile', METHOD.POST, {
          graph_kind: graphKind,
          project_uuid: this.projectUuid,
          focus_key: '',
        })
        if (!res.data || res.data.code !== 1000) {
          this.actionMsg = (res.data && (res.data.message || res.data.detail)) || '编译没有成功'
          return
        }
        if (!silent && this.mainTab === graphKind) await this.loadGraph()
        if (silent && graphKind === 'system') {
          /* 后台编好，切到系统实例时再取 */
        }
      } catch (e) {
        this.actionMsg = (e && e.message) || '编译失败'
      } finally {
        this.busy = false
      }
    },
    async loadGraph() {
      if (this.mainTab === 'schema') {
        this.paintSchema()
        return
      }
      if (this.mainTab === 'world' && !this.projectUuid) return
      this.busy = true
      try {
        const res = await request('/api/ontology/graph', METHOD.GET, {
          graph_kind: this.mainTab,
          project_uuid: this.projectUuid,
          view: 'all',
        })
        const data = res.data || {}
        this.nodes = data.nodes || []
        this.edges = data.edges || []
        this.lint = data.lint || {}
        if (!this.nodes.length) await this.compileCurrent(this.mainTab, true)
        if (!this.nodes.length) {
          const again = await request('/api/ontology/graph', METHOD.GET, {
            graph_kind: this.mainTab,
            project_uuid: this.projectUuid,
            view: 'all',
          })
          this.nodes = (again.data && again.data.nodes) || []
          this.edges = (again.data && again.data.edges) || []
          this.lint = (again.data && again.data.lint) || {}
        }
        this.paintVisible()
      } catch (e) {
        this.actionMsg = (e && e.message) || '读取图失败'
      } finally {
        this.busy = false
      }
    },
    async runAction() {
      this.busy = true
      try {
        const res = await request('/api/ontology/actions/' + this.actionId, METHOD.POST, {
          project_uuid: this.projectUuid,
          alarm_id: this.alarmId,
          title: this.titleText,
          interface_id: 'OntologyCenter',
        })
        this.actionMsg = JSON.stringify(res.data, null, 2)
        if (res.data && res.data.code === 1000) await this.loadGraph()
      } catch (e) {
        this.actionMsg = (e && e.message) || '动作失败'
      } finally {
        this.busy = false
      }
    },
    async runAlarm() {
      this.busy = true
      try {
        const res = await request('/api/agent/runs', METHOD.POST, {
          workflow_code: 'alarm_lifecycle',
          project_uuid: this.projectUuid,
        })
        this.run = (res.data && res.data.run) || null
        this.actionMsg = res.data && res.data.code === 1000 ? '研判完成，模型步骤已跳过' : JSON.stringify(res.data)
      } catch (e) {
        this.actionMsg = (e && e.message) || '研判失败'
      } finally {
        this.busy = false
      }
    },
  },
  watch: {
    showPoints() { if (this.mainTab !== 'schema') this.paintInstances() },
  },
}
</script>

<style scoped>
.ont-center { display: flex; flex-direction: column; height: calc(100vh - 140px); min-height: 640px; padding: 12px; background: #f4f7fb; }
.ont-bar { display: flex; align-items: center; flex-wrap: wrap; gap: 8px; }
.ont-hint { color: #52606d; font-size: 12px; margin-left: 8px; }
.ont-spacer { flex: 1; }
.ont-body { display: flex; gap: 12px; margin-top: 12px; min-height: 0; flex: 1; }
.ont-side { width: 260px; flex: none; overflow: auto; background: #fff; border: 1px solid #d9e2ec; border-radius: 8px; padding: 8px; }
.ont-detail { width: 300px; flex: none; overflow: auto; background: #fff; border: 1px solid #d9e2ec; border-radius: 8px; padding: 8px; }
.ont-main { flex: 1; min-width: 0; display: flex; flex-direction: column; min-height: 0; }
.ont-stage { position: relative; flex: 1; min-height: 0; }
.ont-stage >>> .ont-star { min-height: 0; height: 100%; }
.ont-legend { position: absolute; left: 10px; bottom: 10px; display: flex; flex-wrap: wrap; gap: 6px; z-index: 2; max-width: 70%; }
.ont-legend button { display: flex; align-items: center; gap: 6px; border: 1px solid rgba(255,255,255,0.16); background: rgba(7, 17, 31, 0.78); color: #e7f6ff; border-radius: 999px; padding: 4px 10px; cursor: pointer; font-size: 12px; }
.ont-legend button.off { opacity: 0.38; }
.ont-legend i { width: 8px; height: 8px; border-radius: 50%; display: inline-block; }
.ont-rail { width: 28px; flex: none; border: 1px solid #d9e2ec; background: #fff; border-radius: 8px; cursor: pointer; writing-mode: vertical-rl; letter-spacing: 4px; color: #334e68; font-size: 12px; }
.ont-foldbar { display: flex; align-items: center; justify-content: space-between; gap: 8px; font-size: 12px; font-weight: 600; color: #334e68; margin-bottom: 4px; }
.ont-fold { border: 0; background: #eef4f8; color: #334e68; width: 22px; height: 22px; border-radius: 4px; cursor: pointer; line-height: 22px; padding: 0; }
.ont-center:fullscreen { height: 100vh; background: #eef3f8; }
.ont-center:fullscreen .ont-side,
.ont-center:fullscreen .ont-detail,
.ont-center:fullscreen .ont-rail { background: #fff; }
.ont-stats { color: #52606d; font-size: 12px; margin-bottom: 6px; }
.ont-crumb { display: flex; flex-wrap: wrap; align-items: center; gap: 6px; margin-bottom: 6px; color: #52606d; font-size: 12px; }
.ont-crumb button { border: 0; background: #e6f6ff; color: #186faf; border-radius: 4px; padding: 2px 8px; cursor: pointer; }
.ont-sec { margin: 10px 4px 4px; font-weight: 600; font-size: 12px; color: #334e68; }
.ont-row { display: flex; align-items: center; width: 100%; border: 0; background: transparent; text-align: left; padding: 4px; cursor: pointer; border-radius: 4px; }
.ont-row.on, .ont-row:hover { background: #e6f6ff; }
.ont-row i { width: 8px; height: 8px; border-radius: 50%; margin-right: 6px; flex: none; }
.ont-row span { flex: 1; }
.ont-row em { color: #829ab1; font-style: normal; font-size: 11px; }
.ont-mini { font-size: 12px; color: #486581; padding: 2px 4px; }
.muted { color: #829ab1; }
.hit { margin: 4px 0; font-size: 12px; }
.hit.error { color: #9b1c1c; }
.hit.warning { color: #9a6700; }
.msg { white-space: pre-wrap; font-size: 12px; max-height: 140px; overflow: auto; }
code { font-size: 12px; }
</style>
