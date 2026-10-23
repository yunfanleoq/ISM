<template>
  <div
    ref="root"
    class="demo-root"
    :data-state="session.state"
    :data-busy="session.busy ? '1' : '0'"
    :data-fly-ms="session.flyMs"
  >
    <header class="demo-head">
      <h1>
        <span>3号楼 · 动环演示</span>
        <small>{{ title }}</small>
        <em class="demo-badge" data-demo-badge>演示数据</em>
      </h1>
      <nav class="demo-nav">
        <button
          v-for="item in navItems"
          :key="item.id"
          type="button"
          :data-nav="item.id"
          :class="{ 'is-on': activeNav === item.id }"
          @click="go('nav', { id: item.id })"
        >{{ item.label }}</button>
      </nav>
      <div class="demo-head-actions">
        <button
          v-if="backText"
          type="button"
          class="demo-textbtn"
          :data-action="session.state === 'follow' ? 'exit-follow' : 'back'"
          @click="go('back')"
        >{{ backText }}</button>
        <button type="button" class="demo-textbtn" data-action="fullscreen" @click="toggleFull">全屏</button>
      </div>
    </header>

    <div class="demo-body" :class="{ 'is-flying': session.busy }">
      <div class="demo-stage-wrap" :class="{ 'is-dim': dimStage }">
        <demo-stage
          ref="stage"
          :state="session.state"
          :racks="model.racks"
          :interactive="stageInteractive"
          @pick-building="go('pick-building')"
          @pick-rack="onRack"
        />
      </div>

      <aside class="glass demo-side">
        <template v-if="session.state === 'cockpit'">
          <h2>进线功率</h2>
          <div v-for="feed in model.feeds" :key="feed.name">
            <div class="row-line"><span>{{ feed.name }}</span><b>{{ feed.kw }} kW</b></div>
            <div class="bar"><i :style="{ width: feed.loadPct + '%' }" /></div>
          </div>
          <h2>室外</h2>
          <div class="row-line"><span>温度</span><b>{{ model.outdoor.temp }} °C</b></div>
          <div class="row-line"><span>湿度</span><b>{{ model.outdoor.humidity }} %</b></div>
          <h2>机房 IT 功率</h2>
          <div v-for="room in itRooms" :key="room.id">
            <div class="row-line"><span>{{ room.label }}</span><span>{{ room.itKw }} kW</span></div>
            <div class="bar"><i :style="{ width: (room.itKw / 420 * 100) + '%' }" /></div>
          </div>
        </template>

        <template v-else-if="session.state === 'building'">
          <h2>动环系统</h2>
          <p>一层到四层都在。点右侧楼层标签落入该层，演示主路径是四层。</p>
        </template>
        <template v-else-if="session.state === 'floor'">
          <h2>楼层</h2>
          <div class="stack">
            <button
              v-for="floor in model.floors"
              :key="floor.id"
              type="button"
              class="floor-tag"
              :data-hotspot="'floor-tree-' + floor.id"
              @click="go('pick-floor', { floorId: floor.id })"
            >{{ floor.label }}</button>
          </div>
        </template>

        <template v-else-if="session.state === 'room' || session.state === 'follow'">
          <h2>{{ currentRoomLabel }}</h2>
          <div class="row-line"><span>最高温度</span><b>{{ model.roomEnv.maxTemp }} °C</b></div>
          <div class="row-line"><span>最低温度</span><b>{{ model.roomEnv.minTemp }} °C</b></div>
          <div class="row-line"><span>漏水</span><b>{{ model.roomEnv.leak }}</b></div>
          <div class="row-line"><span>空调运行</span><b>{{ model.roomEnv.acOn }}/{{ model.roomEnv.acTotal }}</b></div>
        </template>

        <template v-else-if="session.state === 'power'">
          <h2>主路电气量</h2>
          <div class="metric">{{ model.power.voltage }}<span>V</span></div>
          <p class="unit">线电压</p>
          <div class="metric">{{ model.power.current }}<span>A</span></div>
          <p class="unit">相电流</p>
          <div class="metric">{{ model.power.kw }}<span>kW</span></div>
          <p class="unit">有功功率 · 负载 {{ model.power.loadPct }}%</p>
        </template>

        <template v-else-if="session.state === 'battery'">
          <h2>UPS</h2>
          <div class="row-line"><span>电池电压</span><b>{{ model.ups.batteryVoltage }} V</b></div>
          <div class="row-line"><span>电池电流</span><b>{{ model.ups.batteryCurrent }} A</b></div>
          <div class="row-line"><span>后备</span><b>{{ model.ups.backupMin }} min</b></div>
          <div class="row-line"><span>剩余容量</span><b>{{ model.ups.capacityPct }}%</b></div>
          <div class="row-line"><span>输出负载</span><b>{{ model.ups.outputLoadPct }}%</b></div>
          <template v-if="activeCell">
            <h2>单体</h2>
            <p :data-cell-id="activeCell.id">{{ activeCell.id }}</p>
            <div class="metric" data-cell-voltage>{{ activeCell.voltage }}<span>V</span></div>
          </template>
        </template>

        <template v-else>
          <h2>资产</h2>
          <p>位置和状态在右侧台账。流程不在这次演示。</p>
          <p v-if="assetPath" data-asset-path>{{ assetPath }}</p>
        </template>
      </aside>

      <div class="center-ui">
        <button
          v-if="session.state === 'cockpit'"
          type="button"
          class="pin hud"
          data-hotspot="building"
          @click="go('pick-building')"
        >3号楼</button>

        <div v-else-if="session.state === 'building'" class="floor-tags hud">
          <button
            v-for="floor in model.floors"
            :key="floor.id"
            type="button"
            class="floor-tag"
            :data-hotspot="'floor-' + floor.id"
            @click="go('pick-floor', { floorId: floor.id })"
          >{{ floor.overviewLabel }}</button>
        </div>

        <div v-else-if="session.state === 'floor'" class="room-grid hud">
          <button
            v-for="room in floorRooms"
            :key="room.id"
            type="button"
            class="room-card"
            :class="room.kind"
            :data-hotspot="'room-' + room.id"
            @click="onRoom(room)"
          >
            <b>{{ room.label }}</b>
            <span>{{ roomKindText(room.kind) }}</span>
          </button>
        </div>

        <div v-else-if="session.state === 'room' || session.state === 'follow'" class="hud">
          <div class="rack-row">
            <button
              v-for="rack in roomRacks"
              :key="rack.id"
              type="button"
              class="rack-chip"
              :class="{ alarm: rack.level === 'alarm' }"
              :data-hotspot="'rack-' + rack.id"
              @click="onRack(rack.id)"
            >
              {{ rack.name }}
              <small v-if="rack.level === 'alarm'"> {{ rack.temp }}°C</small>
            </button>
          </div>
          <p v-if="session.state === 'follow'" data-path>{{ followPath }}</p>
        </div>

        <div v-else-if="session.state === 'power'" class="flow hud">
          <div class="flow-row" role="img" aria-label="市电经主路送入机房">
            <span class="flow-node">市电</span>
            <i class="flow-line" />
            <span class="flow-node">UPS主路</span>
            <i class="flow-line" />
            <span class="flow-node">机房IT</span>
          </div>
          <div class="flow-bypass-row">
            <span>旁路</span>
            <i class="flow-line is-dead" />
            <span data-bypass-current="0">{{ model.power.bypassCurrent }} A</span>
          </div>
          <button type="button" class="demo-textbtn" data-action="bypass" @click="bypassNote = !bypassNote">旁路</button>
          <p v-if="bypassNote" class="bypass-note">当前未走旁路</p>
        </div>

        <div v-else-if="session.state === 'battery'" class="hud" style="width:min(720px,100%)">
          <div class="ups-strip">
            <span class="live">主路</span>
            <span class="dead">旁路</span>
            <span class="live">电池浮充</span>
            <span class="live">{{ model.ups.mode }}</span>
          </div>
          <div class="cell-grid" :data-watch-count="watchCount">
            <button
              v-for="cell in model.batteryRoom.cells"
              :key="cell.id"
              type="button"
              class="cell"
              :class="{ watch: cell.level === 'watch', 'is-pick': activeCell && activeCell.id === cell.id }"
              :data-level="cell.level"
              :data-cell="cell.id"
              @click="activeCell = cell"
            >
              {{ cell.voltage }}V
              <small>{{ cell.id.slice(-2) }}</small>
            </button>
          </div>
        </div>

        <div v-else-if="session.state === 'asset'" class="cabinet hud" data-cabinet="A01">
          <div
            v-for="u in uSlots"
            :key="u"
            class="u-slot"
            :class="slotClass(u)"
            :data-u="u"
          >
            <span>U{{ u }}</span>
            <span>{{ slotName(u) }}</span>
          </div>
        </div>
      </div>

      <aside class="glass demo-side">
        <template v-if="session.state === 'cockpit'">
          <h2>PUE</h2>
          <div class="metric" data-pue="1.15">{{ model.pue }}</div>
          <p class="unit">WUE {{ model.wue }} · 演示</p>
          <h2>能耗</h2>
          <div class="metric">{{ model.energy.kw }}<span>kW</span></div>
          <p class="unit">总电能 {{ model.energy.kwh }} kWh</p>
          <div class="donuts">
            <div class="donut" :style="donut(model.energy.itPct, '#5cc8ff')"><span>IT {{ model.energy.itPct }}%</span></div>
            <div class="donut" :style="donut(model.energy.coolPct, '#7ed0e8')"><span>冷</span></div>
            <div class="donut" :style="donut(model.energy.otherPct, '#8eadc4')"><span>其他</span></div>
          </div>
        </template>
        <template v-else-if="session.state === 'floor'">
          <h2>从这一层进入</h2>
          <div class="stack subbar">
            <button type="button" data-hotspot="enter-env" @click="go('pick-room', { roomId: 'r4-1' })">环境</button>
            <button type="button" data-hotspot="enter-power" @click="go('pick-power-room')">电力</button>
            <button type="button" data-hotspot="enter-battery" @click="go('pick-battery-room')">电池</button>
          </div>
        </template>
        <template v-else-if="session.state === 'follow'">
          <h2>{{ alarmRack.name }}</h2>
          <div class="metric">{{ alarmRack.temp }}<span>°C</span></div>
          <p data-path>{{ followPath }}</p>
        </template>
        <template v-else-if="session.state === 'power'">
          <h2>路径</h2>
          <p>电沿主路送进机房。旁路电流为 0。</p>
        </template>
        <template v-else-if="session.state === 'asset'">
          <h2>台账</h2>
          <button
            v-for="row in assets"
            :key="row.id"
            type="button"
            class="asset-row"
            :class="{ 'is-on': selectedAssetId === row.id }"
            :data-asset-id="row.id"
            :data-asset-status="row.status"
            :data-u-start="row.uStart || ''"
            @click="onAssetRow(row)"
          >
            <b>{{ row.name }}</b>
            <div>{{ row.status }} · {{ row.place }}</div>
          </button>
          <button
            type="button"
            class="demo-textbtn"
            data-action="shelve"
            :disabled="shelving || shelved"
            @click="shelve"
          >演示上架</button>
        </template>
        <template v-else>
          <h2>同一栋楼</h2>
          <p>镜头从楼到层，再进机房、配电和电池。</p>
        </template>
      </aside>

      <div v-if="session.state === 'cockpit'" class="glass month">
        <i v-for="(v, i) in model.monthKwh" :key="i" :style="{ height: v + '%' }" />
      </div>
    </div>
  </div>
</template>

<script>
import DemoStage from '@/xunzhiJf/demoShow/DemoStage.vue'
import { loadDemoSource } from '@/xunzhiJf/demoShow/demoSource'
import {
  NAV_ITEMS,
  backLabel,
  begin,
  createSession,
  land,
  navKey,
  poseFor,
  screenTitle,
} from '@/xunzhiJf/demoShow/cameraFsm'
export default {
  name: 'XunzhiJfDemoShowPage',
  components: { DemoStage },
  data() {
    const model = loadDemoSource()
    return {
      model,
      navItems: NAV_ITEMS,
      session: createSession(),
      assets: model.assets.map(row => ({ ...row })),
      selectedAssetId: 'a-server',
      activeCell: null,
      bypassNote: false,
      shelving: false,
      _landTimer: null,
    }
  },
  computed: {
    title() { return screenTitle(this.session.state) },
    activeNav() { return navKey(this.session.state) },
    backText() { return backLabel(this.session.state) },
    dimStage() {
      return this.session.state === 'power' || this.session.state === 'battery' || this.session.state === 'asset'
    },
    stageInteractive() {
      const s = this.session.state
      return !this.session.busy && (s === 'room' || s === 'asset')
    },
    itRooms() {
      return this.model.rooms.filter(r => r.floorId === 'f4' && r.kind === 'it')
    },
    floorRooms() {
      return this.model.rooms.filter(r => r.floorId === this.session.floorId)
    },
    roomRacks() {
      return this.model.racks.filter(r => r.roomId === 'r4-1')
    },
    alarmRack() {
      return this.model.racks.find(r => r.id === 'rk-c2') || { name: 'C2', temp: 31.6 }
    },
    followPath() {
      return '3号楼 / 4层 / 机房4-1 / ' + this.alarmRack.name
    },
    currentRoomLabel() {
      const room = this.model.rooms.find(r => r.id === this.session.roomId)
      return room ? room.label : '机房4-1'
    },
    watchCount() {
      return this.model.batteryRoom.cells.filter(c => c.level === 'watch').length
    },
    assetPath() {
      const row = this.assets.find(a => a.id === this.selectedAssetId)
      return row ? row.place : ''
    },
    shelved() {
      const row = this.assets.find(a => a.id === 'a-server')
      return row && row.status === '在用'
    },
    uSlots() {
      const list = []
      for (let u = 42; u >= 1; u--) list.push(u)
      return list
    },
  },
  mounted() {
    this.$nextTick(() => {
      if (this.$refs.stage) this.$refs.stage.setPose(poseFor('cockpit'))
    })
  },
  beforeDestroy() {
    clearTimeout(this._landTimer)
  },
  methods: {
    go(intent, payload) {
      const result = begin(this.session, intent, payload)
      if (!result.ok) return
      this.session = result.session
      this.$nextTick(() => {
        if (this.$refs.stage) this.$refs.stage.flyTo(poseFor(this.session.state), this.session.flyMs)
      })
      clearTimeout(this._landTimer)
      this._landTimer = setTimeout(() => {
        this.session = land(this.session)
      }, result.session.flyMs)
    },
    onRoom(room) {
      if (room.kind === 'power') this.go('pick-power-room')
      else if (room.kind === 'battery') this.go('pick-battery-room')
      else this.go('pick-room', { roomId: room.id })
    },
    onRack(rackId) {
      const rack = this.model.racks.find(r => r.id === rackId)
      if (!rack || rack.level !== 'alarm') return
      this.go('pick-rack', { rackId })
    },
    onAssetRow(row) {
      this.selectedAssetId = row.id
      this.go('select-asset')
    },
    shelve() {
      if (this.session.busy || this.shelving || this.shelved) return
      this.shelving = true
      setTimeout(() => {
        const row = this.assets.find(a => a.id === 'a-server')
        if (row) {
          row.status = '在用'
          row.uStart = row.shelveU
        }
        this.shelving = false
      }, 1200)
    },
    slotClass(u) {
      if (this.shelving && u === 2) return 'dropping'
      const row = this.occupant(u)
      if (!row) return ''
      if (row.status === '借用到期') return 'borrow'
      return 'occupied'
    },
    occupant(u) {
      return this.assets.find(a => a.uStart && u >= a.uStart && u < a.uStart + a.uHeight)
    },
    slotName(u) {
      if (this.shelving && u === 2) return '计算节点-07'
      const row = this.occupant(u)
      return row ? row.name : ''
    },
    roomKindText(kind) {
      if (kind === 'power') return '配电'
      if (kind === 'battery') return '电池间'
      return '机房'
    },
    donut(pct, color) {
      return {
        background: 'conic-gradient(' + color + ' ' + pct + '%, rgba(255,255,255,0.08) 0)',
      }
    },
    toggleFull() {
      const el = this.$refs.root
      if (!document.fullscreenElement) {
        if (el.requestFullscreen) el.requestFullscreen()
      } else if (document.exitFullscreen) {
        document.exitFullscreen()
      }
    },
  },
}
</script>

<style src="../demoShow/demo-show.css"></style>
