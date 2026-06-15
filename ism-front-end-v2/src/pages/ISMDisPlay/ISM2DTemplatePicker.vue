<template>
  <transition name="tpl-fade">
    <div v-if="visible" class="tpl-overlay" @click.self="onClose">
      <div class="tpl-panel">
        <div class="tpl-header">
          <div class="tpl-window-dots">
            <span></span>
            <span></span>
            <span></span>
          </div>
          <div class="tpl-title">
            <i class="fas fa-layer-group"></i>
            <span>{{ $t('displayConfig.Template2D.title') }}</span>
            <b>{{ templates.length }}</b>
          </div>
          <button class="tpl-close" :title="$t('displayConfig.Template2D.close')" @click="onClose">
            <i class="fas fa-times"></i>
          </button>
        </div>

        <div class="tpl-toolbar">
          <div class="tpl-search">
            <i class="fas fa-search"></i>
            <input v-model="keyword" :placeholder="$t('displayConfig.Template2D.searchPlaceholder')" />
          </div>
          <div class="tpl-cats">
            <button
              v-for="cat in categories"
              :key="cat"
              :class="{ active: currentCategory === cat }"
              @click="currentCategory = cat"
            >
              {{ cat }}
            </button>
          </div>
        </div>

        <div class="tpl-grid-wrap">
          <div class="tpl-grid">
            <div v-for="tpl in filteredTemplates" :key="tpl.id" class="tpl-card" @click="onSelect(tpl)">
              <div class="tpl-preview" :style="previewStyle(tpl)">
                <div class="tpl-preview-layout" :class="`tpl-preview-${tpl.id}`">
                  <span class="preview-panel preview-panel-left"></span>
                  <span class="preview-panel preview-panel-main"></span>
                  <span class="preview-panel preview-panel-right"></span>
                  <span
                    v-for="line in previewLines(tpl.id)"
                    :key="line"
                    class="preview-line"
                    :style="{ '--line-index': line }"
                  ></span>
                  <span
                    v-for="bar in previewBars(tpl.id)"
                    :key="bar.key"
                    class="preview-bar"
                    :style="{ left: bar.left, height: bar.height }"
                  ></span>
                </div>
                <i :class="tpl.icon"></i>
                <span class="tpl-preview-count">{{ getTemplateScene(tpl).components.cells.length }} {{ $t('displayConfig.Template2D.components') }}</span>
              </div>
              <div class="tpl-body">
                <div class="tpl-name">{{ localText(tpl.name) }}</div>
                <div class="tpl-desc">{{ localText(tpl.description) }}</div>
                <div class="tpl-highlights">
                  <span v-for="item in previewHighlights(tpl.id)" :key="localText(item)">{{ localText(item) }}</span>
                </div>
                <div class="tpl-meta">
                  <span>{{ localText(tpl.category) }}</span>
                  <button type="button">
                    <i class="fas fa-play"></i>
                    {{ $t('displayConfig.Template2D.apply') }}
                  </button>
                </div>
              </div>
            </div>
            <div v-if="filteredTemplates.length === 0" class="tpl-empty">
              <i class="fas fa-search"></i>
              <span>{{ $t('displayConfig.Template2D.empty') }}</span>
            </div>
          </div>
        </div>

        <div class="tpl-footer">
          <span>{{ $t('displayConfig.Template2D.footerHint') }}</span>
          <b>{{ filteredTemplates.length }} / {{ templates.length }}</b>
        </div>
      </div>
    </div>
  </transition>
</template>

<script>
import { mapState } from 'vuex'
import { TEMPLATE_2D_SCENES } from './utils/Template2DScenes'

const CATEGORY_COLORS = [
  ['#061321', '#0f4c75'],
  ['#07101d', '#7c4d00'],
  ['#07151e', '#0b6b5a']
]

export default {
  name: 'ISM2DTemplatePicker',
  i18n: require('../../i18n/language'),
  props: {
    visible: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      templates: TEMPLATE_2D_SCENES,
      keyword: '',
      currentCategory: ''
    }
  },
  computed: {
    ...mapState('setting', ['lang']),
    allText() {
      return this.$t('displayConfig.Template2D.all')
    },
    categories() {
      const cats = this.templates.map(tpl => this.localText(tpl.category))
      return [this.allText, ...Array.from(new Set(cats))]
    },
    filteredTemplates() {
      const key = this.keyword.trim().toLowerCase()
      return this.templates.filter(tpl => {
        const name = this.localText(tpl.name)
        const desc = this.localText(tpl.description)
        const cat = this.localText(tpl.category)
        const matchCat = !this.currentCategory || this.currentCategory === this.allText || cat === this.currentCategory
        const matchKey = !key || [name, desc, cat].some(text => String(text).toLowerCase().includes(key))
        return matchCat && matchKey
      })
    }
  },
  watch: {
    allText: {
      immediate: true,
      handler(value) {
        this.currentCategory = value
      }
    },
    visible(value) {
      if (value) {
        this.keyword = ''
        this.currentCategory = this.allText
      }
    }
  },
  methods: {
    localText(value) {
      if (!value || typeof value !== 'object') return value || ''
      return value[this.lang] || value.CN || value.US || ''
    },
    getTemplateScene(tpl) {
      return tpl.build()
    },
    previewStyle(tpl) {
      const index = Math.max(0, this.templates.findIndex(item => item.id === tpl.id))
      const colors = CATEGORY_COLORS[index % CATEGORY_COLORS.length]
      return {
        background: `radial-gradient(circle at 70% 30%, ${colors[1]} 0%, ${colors[0]} 68%)`
      }
    },
    previewLines(id) {
      return id === 'power-room' ? [1, 2, 3, 4] : [1, 2, 3]
    },
    previewBars(id) {
      const maps = {
        'factory-dashboard': [
          { key: 'a', left: '42%', height: '24px' },
          { key: 'b', left: '50%', height: '36px' },
          { key: 'c', left: '58%', height: '29px' }
        ],
        'power-room': [
          { key: 'a', left: '46%', height: '34px' },
          { key: 'b', left: '55%', height: '22px' }
        ],
        'energy-overview': [
          { key: 'a', left: '42%', height: '20px' },
          { key: 'b', left: '50%', height: '32px' },
          { key: 'c', left: '58%', height: '44px' },
          { key: 'd', left: '66%', height: '27px' }
        ]
      }
      return maps[id] || maps['factory-dashboard']
    },
    previewHighlights(id) {
      const maps = {
        'factory-dashboard': [
          { CN: '设备总览', US: 'Device status', HK: '設備總覽' },
          { CN: '产线趋势', US: 'Line trend', HK: '產線趨勢' },
          { CN: '告警待办', US: 'Alarms', HK: '告警待辦' }
        ],
        'power-room': [
          { CN: '一次系统', US: 'One-line', HK: '一次系統' },
          { CN: '环境监测', US: 'Environment', HK: '環境監測' },
          { CN: '回路告警', US: 'Circuit alarms', HK: '迴路告警' }
        ],
        'energy-overview': [
          { CN: '分项能耗', US: 'Energy split', HK: '分項能耗' },
          { CN: '能耗排行', US: 'Ranking', HK: '能耗排行' },
          { CN: '节能建议', US: 'Advice', HK: '節能建議' }
        ]
      }
      return maps[id] || maps['factory-dashboard']
    },
    onClose() {
      this.$emit('close')
    },
    onSelect(tpl) {
      this.$emit('select', tpl)
    }
  }
}
</script>

<style scoped>
.tpl-overlay {
  position: fixed;
  inset: 0;
  z-index: 9999;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding-top: 6vh;
  background: rgba(240, 242, 245, 0.45);
  backdrop-filter: blur(2px);
  overflow-y: auto;
}
.tpl-panel {
  width: min(92vw, 1080px);
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: #fff;
  border: 1px solid #e8e8e8;
  border-radius: 10px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.08);
}
.tpl-header,
.tpl-footer {
  flex: none;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 20px;
  background: #fafbfc;
  border-bottom: 1px solid #f0f0f0;
}
.tpl-footer {
  border-top: 1px solid #f0f0f0;
  border-bottom: 0;
  color: #8c8c8c;
  font-size: 12px;
}
.tpl-window-dots {
  display: flex;
  gap: 6px;
}
.tpl-window-dots span {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}
.tpl-window-dots span:nth-child(1) { background: #ff5f57; }
.tpl-window-dots span:nth-child(2) { background: #ffbd2e; }
.tpl-window-dots span:nth-child(3) { background: #28ca42; }
.tpl-title {
  display: flex;
  align-items: center;
  gap: 8px;
  height: 32px;
  color: #333;
  font-size: 15px;
  font-weight: 600;
  line-height: 20px;
  white-space: nowrap;
}
.tpl-title i {
  flex: none;
  color: #13c2c2;
}
.tpl-title span {
  line-height: 20px;
}
.tpl-title b {
  display: inline-flex;
  flex: none;
  align-items: center;
  justify-content: center;
  box-sizing: border-box;
  min-width: 22px;
  height: 20px;
  padding: 0 7px;
  color: #13c2c2;
  font-size: 11px;
  font-weight: 400;
  line-height: 18px;
  text-orientation: mixed;
  vertical-align: middle;
  white-space: nowrap;
  writing-mode: horizontal-tb;
  background: #e6fffb;
  border: 1px solid #b5f5ec;
  border-radius: 10px;
}
.tpl-close {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
  text-orientation: mixed;
  writing-mode: horizontal-tb;
  cursor: pointer;
  color: #888;
  background: #f5f5f5;
  border: 1px solid #e8e8e8;
  border-radius: 6px;
}
.tpl-close:hover {
  color: #fff;
  background: #ff4d4f;
  border-color: #ff4d4f;
}
.tpl-toolbar {
  flex: none;
  display: flex;
  gap: 12px;
  align-items: center;
  padding: 12px 20px 8px;
  background: #fafbfc;
  border-bottom: 1px solid #f6f6f6;
}
.tpl-search {
  position: relative;
  width: 240px;
  flex: none;
}
.tpl-search i {
  position: absolute;
  left: 10px;
  top: 50%;
  color: #bfbfbf;
  transform: translateY(-50%);
}
.tpl-search input {
  width: 100%;
  height: 32px;
  padding: 0 10px 0 30px;
  line-height: 30px;
  border: 1px solid #d9d9d9;
  border-radius: 6px;
  outline: none;
}
.tpl-search input:focus {
  border-color: #13c2c2;
  box-shadow: 0 0 0 2px rgba(19, 194, 194, 0.12);
}
.tpl-cats {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 6px;
  min-height: 32px;
}
.tpl-cats button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  box-sizing: border-box;
  min-width: 52px;
  height: 28px;
  padding: 0 12px;
  color: #666;
  font-size: 12px;
  line-height: 26px;
  text-align: center;
  text-orientation: mixed;
  white-space: nowrap;
  writing-mode: horizontal-tb;
  background: #fff;
  border: 1px solid #e8e8e8;
  border-radius: 14px;
  cursor: pointer;
}
.tpl-cats button.active,
.tpl-cats button:hover {
  color: #fff;
  background: #13c2c2;
  border-color: #13c2c2;
}
.tpl-grid-wrap {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 16px 20px 24px;
  background: #f5f7f8;
}
.tpl-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 12px;
}
.tpl-card {
  overflow: hidden;
  cursor: pointer;
  background: #fff;
  border: 1px solid #eee;
  border-radius: 8px;
  transition: border-color 0.2s, box-shadow 0.2s;
}
.tpl-card:hover {
  border-color: #13c2c2;
  box-shadow: 0 4px 14px rgba(19, 194, 194, 0.12);
}
.tpl-preview {
  position: relative;
  height: 104px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}
.tpl-preview::before {
  content: '';
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(19, 194, 194, 0.08) 1px, transparent 1px),
    linear-gradient(90deg, rgba(19, 194, 194, 0.08) 1px, transparent 1px);
  background-size: 22px 22px;
}
.tpl-preview i {
  position: relative;
  z-index: 3;
  color: rgba(127, 252, 255, 0.78);
  font-size: 28px;
}
.tpl-preview-count {
  position: absolute;
  right: 8px;
  bottom: 8px;
  z-index: 4;
  padding: 2px 6px;
  color: #fff;
  font-size: 10px;
  background: rgba(0, 0, 0, 0.45);
  border-radius: 3px;
}
.tpl-preview-layout {
  position: absolute;
  inset: 12px;
  z-index: 2;
}
.preview-panel {
  position: absolute;
  display: block;
  border: 1px solid rgba(127, 252, 255, 0.22);
  background: rgba(6, 24, 38, 0.46);
  border-radius: 4px;
}
.preview-panel-left {
  left: 0;
  top: 7px;
  width: 22%;
  height: 68px;
}
.preview-panel-main {
  left: 28%;
  top: 0;
  width: 44%;
  height: 76px;
}
.preview-panel-right {
  right: 0;
  top: 7px;
  width: 22%;
  height: 68px;
}
.preview-line {
  position: absolute;
  left: 31%;
  top: calc(18px + var(--line-index) * 13px);
  width: 38%;
  height: 3px;
  overflow: hidden;
  background: rgba(127, 252, 255, 0.12);
  border-radius: 3px;
}
.preview-line::after {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  width: 68%;
  height: 100%;
  background: rgba(127, 252, 255, 0.55);
  border-radius: inherit;
}
.preview-bar {
  position: absolute;
  bottom: 12px;
  width: 10px;
  background: rgba(127, 252, 255, 0.62);
  border-radius: 6px 6px 2px 2px;
  opacity: 0.86;
}
.tpl-preview-power-room .preview-line {
  left: 24%;
  width: 52%;
}
.tpl-preview-energy-overview .preview-bar {
  background: rgba(255, 214, 102, 0.72);
}
.tpl-body {
  padding: 10px 12px 12px;
}
.tpl-name {
  margin-bottom: 6px;
  overflow: hidden;
  color: #333;
  font-size: 13px;
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.tpl-desc {
  min-height: 32px;
  margin-bottom: 8px;
  overflow: hidden;
  color: #888;
  font-size: 11px;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}
.tpl-highlights {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
  min-height: 22px;
  margin-bottom: 8px;
}
.tpl-highlights span {
  height: 20px;
  padding: 0 6px;
  overflow: hidden;
  color: #5b7c8f;
  font-size: 10px;
  line-height: 18px;
  text-overflow: ellipsis;
  white-space: nowrap;
  background: #f3fafb;
  border: 1px solid #d8eef0;
  border-radius: 3px;
}
.tpl-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  min-height: 24px;
}
.tpl-meta span {
  max-width: 68%;
  height: 24px;
  padding: 0 8px;
  overflow: hidden;
  color: #13c2c2;
  font-size: 10px;
  line-height: 22px;
  text-orientation: mixed;
  text-overflow: ellipsis;
  white-space: nowrap;
  writing-mode: horizontal-tb;
  background: #e6fffb;
  border: 1px solid #b5f5ec;
  border-radius: 4px;
}
.tpl-meta button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex: none;
  height: 24px;
  padding: 0 10px;
  color: #fff;
  font-size: 11px;
  line-height: 24px;
  text-orientation: mixed;
  white-space: nowrap;
  writing-mode: horizontal-tb;
  background: #13c2c2;
  border: 0;
  border-radius: 4px;
  cursor: pointer;
}
.tpl-meta button i {
  margin-right: 4px;
  font-size: 9px;
}
.tpl-empty {
  grid-column: 1 / -1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 56px 0;
  color: #bfbfbf;
  font-size: 13px;
}
.tpl-empty i {
  font-size: 30px;
  opacity: 0.45;
}
.tpl-fade-enter-active { transition: opacity 0.2s; }
.tpl-fade-leave-active { transition: opacity 0.15s; }
.tpl-fade-enter,
.tpl-fade-leave-to { opacity: 0; }
</style>
