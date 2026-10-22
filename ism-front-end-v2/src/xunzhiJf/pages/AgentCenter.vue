<template>
  <div class="agent-center">
    <div class="agent-bar">
      <strong>智能体</strong>
      <span>按工作流看阶段和环节。技能、工具单独列在右边。模型步骤会跳过，工具只读。</span>
      <span class="spacer" />
      <a-button size="small" @click="$router.push('/xunzhi/jf/ontology')">本体中心</a-button>
      <a-button size="small" type="primary" :disabled="!current || !current.runnable" :loading="running" @click="runCurrent">跑告警研判</a-button>
    </div>
    <div class="agent-body">
      <aside>
        <div class="sec">工作流</div>
        <button v-for="flow in catalog.workflows" :key="flow.code" type="button" class="row" :class="{ on: current && current.code === flow.code }" @click="current = flow">
          <b>{{ flow.name }}</b>
          <em>{{ flow.runnable ? '可运行' : '已登记' }}</em>
        </button>
      </aside>
      <main>
        <p v-if="!current" class="muted">选一个工作流。</p>
        <template v-else>
          <h3>{{ current.name }}</h3>
          <p v-if="!current.runnable" class="muted">阶段已经登记。运行器目前只接入告警生命周期。</p>
          <ol>
            <li v-for="phase in current.phases" :key="phase.key" :class="{ on: picked && picked.key === phase.key }" @click="picked = phase">
              <div class="phase">阶段 {{ phase.order }} · {{ phase.name }}</div>
              <div v-if="phase.agent" class="step">
                环节 · {{ phase.agent.label }}
                <em v-if="phase.agent.uses_llm">模型步骤，当前跳过</em>
                <em v-else>确定性环节</em>
              </div>
              <div v-else class="step quiet">环节 · {{ phase.name.indexOf('审批') >= 0 ? '人工审批' : '没有智能体' }}</div>
              <div v-if="phase.agent && phase.agent.skills.length" class="chips">
                <span v-for="sk in phase.agent.skills" :key="sk.code">技能 {{ sk.name }}</span>
              </div>
              <div v-if="phase.agent && phase.agent.tools.length" class="chips">
                <span v-for="tool in phase.agent.tools" :key="tool.key">工具 {{ tool.label }}</span>
              </div>
            </li>
          </ol>
        </template>
      </main>
      <aside>
        <div class="sec">技能</div>
        <div v-for="sk in catalog.skills" :key="sk.code" class="card">
          <b>{{ sk.name }}</b>
          <p>{{ sk.body }}</p>
        </div>
        <div class="sec">工具</div>
        <div v-for="tool in catalog.tools" :key="tool.key" class="card">
          <b>{{ tool.label }}</b>
          <p>只读。不能写测点，不能下发。</p>
        </div>
        <div class="sec">最近运行</div>
        <p v-if="!(catalog.runs || []).length" class="muted">还没有运行记录。</p>
        <div v-for="item in catalog.runs || []" :key="item.run.ID" class="card">
          <b>{{ item.run.workflow_code }} · {{ statusText(item.run.status) }}</b>
          <p>{{ item.run.summary }}</p>
          <p v-for="step in item.steps" :key="step.ID">{{ step.step_key }} · {{ statusText(step.status) }} · {{ step.message }}</p>
        </div>
        <pre v-if="msg" class="msg">{{ msg }}</pre>
      </aside>
    </div>
  </div>
</template>

<script>
import { METHOD, request } from '@/utils/request'
import { resolveXunzhiJfProjectUuid } from '@/xunzhiJf/projectContext'

export default {
  name: 'AgentCenter',
  data() {
    return {
      catalog: { workflows: [], skills: [], tools: [], runs: [] },
      current: null,
      picked: null,
      projectUuid: '',
      running: false,
      msg: '',
    }
  },
  mounted() {
    this.boot()
  },
  methods: {
    statusText(status) {
      return { succeeded: '完成', failed: '失败', skipped: '跳过', running: '进行中' }[status] || status || ''
    },
    async boot() {
      try {
        this.projectUuid = (await resolveXunzhiJfProjectUuid()) || ''
      } catch (e) { /* 没有项目时目录仍然可看 */ }
      await this.load()
    },
    async load() {
      try {
        const res = await request('/api/agent/catalog', METHOD.GET, { project_uuid: this.projectUuid })
        this.catalog = (res.data && res.data.catalog) || this.catalog
        if (!this.current && this.catalog.workflows && this.catalog.workflows.length) {
          this.current = this.catalog.workflows[0]
          this.picked = this.current.phases[0]
        }
      } catch (e) {
        this.msg = '智能体目录没有读到。'
      }
    },
    async runCurrent() {
      if (!this.current || !this.current.runnable || !this.projectUuid) {
        this.msg = this.projectUuid ? '这个工作流还不能跑。' : '没有项目，先不跑。'
        return
      }
      this.running = true
      try {
        const res = await request('/api/agent/runs', METHOD.POST, {
          workflow_code: this.current.code,
          project_uuid: this.projectUuid,
        })
        this.msg = res.data && res.data.code === 1000 ? '研判完成，模型步骤已跳过' : JSON.stringify(res.data)
        await this.load()
      } catch (e) {
        this.msg = (e && e.message) || '研判失败'
      } finally {
        this.running = false
      }
    },
  },
}
</script>

<style scoped>
.agent-center { display: flex; flex-direction: column; height: calc(100vh - 140px); min-height: 640px; padding: 12px; background: #f4f7fb; }
.agent-bar { display: flex; align-items: center; gap: 8px; color: #52606d; font-size: 12px; }
.agent-bar strong { color: #102a43; font-size: 16px; }
.spacer { flex: 1; }
.agent-body { display: flex; gap: 12px; margin-top: 12px; min-height: 0; flex: 1; }
aside, main { background: #fff; border: 1px solid #d9e2ec; border-radius: 8px; overflow: auto; }
aside { width: 260px; flex: none; padding: 8px; }
main { flex: 1; padding: 16px; }
.sec { margin: 10px 4px 4px; font-weight: 600; font-size: 12px; color: #334e68; }
.row, li { display: block; width: 100%; text-align: left; border: 0; background: transparent; }
.row { padding: 8px; border-radius: 6px; cursor: pointer; }
.row.on, .row:hover, li.on { background: #e6f6ff; }
.row b { display: block; color: #102a43; }
.row em, .step em { color: #829ab1; font-style: normal; font-size: 12px; }
ol { list-style: none; margin: 0; padding: 0; }
li { border: 1px solid #e4e7eb; border-radius: 8px; padding: 10px 12px; margin-bottom: 8px; cursor: pointer; }
.phase { font-weight: 600; color: #102a43; }
.step { margin-top: 4px; color: #334e68; font-size: 13px; }
.step.quiet { color: #829ab1; }
.chips { display: flex; flex-wrap: wrap; gap: 6px; margin-top: 6px; }
.chips span { background: #f0f4f8; color: #334e68; border-radius: 999px; padding: 2px 8px; font-size: 12px; }
.card { padding: 8px; border-top: 1px solid #f0f4f8; font-size: 12px; color: #334e68; }
.card b { display: block; color: #102a43; }
.card p { margin: 4px 0; }
.muted { color: #829ab1; }
.msg { white-space: pre-wrap; font-size: 12px; }
h3 { margin-top: 0; }
</style>
