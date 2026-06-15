<template>
  <transition name="tpl-fade">
    <div v-if="visible" class="tpl-overlay" @click.self="onClose">
      <!-- 背景扫描线 -->
      <div class="tpl-scanlines" aria-hidden="true"></div>

      <!-- 主容器 -->
      <div class="tpl-panel">
        <!-- 顶部标题栏 -->
        <div class="tpl-header">
          <div class="tpl-header-left">
            <span class="tpl-header-dot tpl-dot-red"></span>
            <span class="tpl-header-dot tpl-dot-yellow"></span>
            <span class="tpl-header-dot tpl-dot-green"></span>
          </div>
          <div class="tpl-header-title">
            <i class="fas fa-layer-group tpl-header-icon"></i>
            <span>场景模板库</span>
            <span class="tpl-header-badge">{{ SCENE_TEMPLATES.length }}</span>
          </div>
          <div class="tpl-header-right">
            <button class="tpl-close-btn" @click="onClose" title="关闭">
              <i class="fas fa-times"></i>
            </button>
          </div>
        </div>

        <!-- 搜索 + 分类 -->
        <div class="tpl-toolbar">
          <div class="tpl-search-wrap">
            <i class="fas fa-search tpl-search-icon"></i>
            <input
              v-model="search"
              class="tpl-search"
              placeholder="搜索场景模板..."
              autocomplete="off"
              @input="onSearch"
            />
          </div>
          <div class="tpl-cats">
            <button
              v-for="cat in allCats"
              :key="cat"
              class="tpl-cat-btn"
              :class="{ active: currentCat === cat }"
              @click="currentCat = cat"
            >
              {{ cat }}
            </button>
          </div>
        </div>

        <!-- 模板网格 -->
        <div class="tpl-grid-wrap scrollbar">
          <div class="tpl-grid">
            <div
              v-for="tpl in filtered"
              :key="tpl.id"
              class="tpl-card"
              :class="{ 'tpl-card--hover': hoverId === tpl.id }"
              @mouseenter="hoverId = tpl.id"
              @mouseleave="hoverId = null"
              @click="onSelect(tpl)"
            >
              <!-- 扫描动效边框 -->
              <span class="tpl-card-border tpl-b-tl"></span>
              <span class="tpl-card-border tpl-b-tr"></span>
              <span class="tpl-card-border tpl-b-bl"></span>
              <span class="tpl-card-border tpl-b-br"></span>

              <!-- 预览区：伪3D渐变背景 -->
              <div class="tpl-card-preview" :style="previewBg(tpl)">
                <i :class="tpl.icon" class="tpl-card-preview-icon"></i>
                <!-- 对象数量标签 -->
                <div class="tpl-card-count">
                  <i class="fas fa-cubes" style="margin-right:4px;"></i>
                  {{ tpl.objects.length }} 个对象
                </div>
              </div>

              <!-- 信息区 -->
              <div class="tpl-card-body">
                <div class="tpl-card-name">{{ tpl.name }}</div>
                <div class="tpl-card-desc">{{ tpl.description }}</div>
                <div class="tpl-card-meta">
                  <span class="tpl-card-category">
                    <i class="fas fa-tag" style="margin-right:3px;"></i>{{ tpl.category }}
                  </span>
                  <button class="tpl-card-apply-btn">
                    <i class="fas fa-play" style="margin-right:4px;font-size:9px;"></i>应用
                  </button>
                </div>
              </div>
            </div>

            <!-- 空状态 -->
            <div v-if="filtered.length === 0" class="tpl-empty">
              <i class="fas fa-search" style="font-size:32px;margin-bottom:12px;opacity:0.3;"></i>
              <div>无匹配模板</div>
            </div>
          </div>
        </div>

        <!-- 底部状态栏 -->
        <div class="tpl-footer">
          <span class="tpl-footer-hint">
            <i class="fas fa-info-circle" style="margin-right:5px;"></i>
            点击模板卡片或「应用」按钮将替换当前场景
          </span>
          <span class="tpl-footer-count">{{ filtered.length }} / {{ SCENE_TEMPLATES.length }} 个模板</span>
        </div>
      </div>
    </div>
  </transition>
</template>

<script>
import { SCENE_TEMPLATES } from '../utils/SceneTemplates'

// 每个分类对应的预览渐变配色
const CAT_COLORS = {
  '工业场景':    ['#f0f8f8', '#e0f5f5'],
  '智慧城市':    ['#f0f5f8', '#e0eef5'],
  '工业流水线':  ['#f5f0f5', '#f0e8f0'],
  '数字孪生':    ['#f0f8f8', '#e6f2f2'],
  '能源电力':    ['#f8f5f0', '#f5efe0'],
  '仓储物流':    ['#f0f2f5', '#e8eef2'],
  '交通设施':    ['#f0f2f8', '#e8ecf5'],
  '室内设施':    ['#f8f5f0', '#f5f0e8'],
  '消防安防':    ['#f8f0f0', '#f5e8e8'],
  '绿化景观':    ['#f0f8f0', '#e8f5e8'],
  '医疗设施':    ['#f0f5f8', '#e8f0f5'],
  '农业设施':    ['#f5f8f0', '#eef5e0'],
  '水利水务':    ['#f0f5f8', '#e0eef5'],
  '运动休闲':    ['#f5f0f8', '#f0e8f5'],
  '实验室':      ['#f0f2f5', '#e8eef2'],
}

export default {
  name: 'SceneTemplatePicker',
  props: {
    visible: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      SCENE_TEMPLATES,
      search: '',
      currentCat: '全部',
      hoverId: null,
    }
  },
  computed: {
    allCats() {
      const cats = ['全部', ...new Set(SCENE_TEMPLATES.map(t => t.category))]
      return cats
    },
    filtered() {
      let list = SCENE_TEMPLATES
      if (this.currentCat !== '全部') {
        list = list.filter(t => t.category === this.currentCat)
      }
      if (this.search.trim()) {
        const q = this.search.trim().toLowerCase()
        list = list.filter(t =>
          (t.name || '').toLowerCase().includes(q) ||
          (t.description || '').toLowerCase().includes(q) ||
          (t.category || '').toLowerCase().includes(q)
        )
      }
      return list
    },
  },
  methods: {
    onClose() {
      this.$emit('close')
    },
    onSearch() {
      // 实时过滤，由 computed 自动处理
    },
    onSelect(tpl) {
      this.$emit('select', tpl)
    },
    previewBg(tpl) {
      const colors = CAT_COLORS[tpl.category] || ['#0d1b2a', '#1a3040']
      return {
        background: `radial-gradient(ellipse at 60% 40%, ${colors[1]}ee 0%, ${colors[0]}ff 100%)`,
      }
    },
  },
}
</script>

<style scoped>
/* ===== 遮罩：透明，只加少量磨砂 ===== */
.tpl-overlay {
  position: fixed;
  inset: 0;
  z-index: 9999;
  background: rgba(240, 242, 245, 0.45);
  backdrop-filter: blur(2px);
  display: flex;
  align-items: flex-start;   /* 顶部对齐，不要居中 */
  justify-content: center;
  padding-top: 6vh;           /* 留一点顶部空间 */
  overflow-y: auto;            /* 遮罩自身可滚，防止面板超高时看不到底部 */
}

/* ===== 主面板：白底 + 极浅阴影 ===== */
.tpl-panel {
  position: relative;
  width: min(92vw, 1160px);
  max-height: 90vh;
  background: #fff;
  border: 1px solid #e8e8e8;
  border-radius: 10px;
  box-shadow:
    0 4px 24px rgba(0, 0, 0, 0.06),
    0 1px 2px rgba(0, 0, 0, 0.04);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* ===== Header ===== */
.tpl-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 20px 10px;
  border-bottom: 1px solid #f0f0f0;
  flex-shrink: 0;
}
.tpl-header-left {
  display: flex;
  gap: 6px;
  align-items: center;
}
.tpl-header-dot {
  width: 10px; height: 10px;
  border-radius: 50%;
  display: inline-block;
}
.tpl-dot-red    { background: #ff5f57; }
.tpl-dot-yellow { background: #ffbd2e; }
.tpl-dot-green  { background: #28ca42; }

.tpl-header-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 15px;
  font-weight: 600;
  color: #333;
  letter-spacing: 0.02em;
}
.tpl-header-icon {
  color: #13c2c2;
  font-size: 16px;
}
.tpl-header-badge {
  background: #e6fffb;
  border: 1px solid #b5f5ec;
  border-radius: 10px;
  font-size: 11px;
  color: #13c2c2;
  padding: 1px 8px;
  font-weight: 400;
}
.tpl-header-right {
  display: flex;
  align-items: center;
}
.tpl-close-btn {
  width: 28px; height: 28px;
  border-radius: 6px;
  border: 1px solid #e8e8e8;
  background: #f5f5f5;
  color: #888;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  outline: none;
}
.tpl-close-btn:hover {
  background: #ff4d4f;
  border-color: #ff4d4f;
  color: #fff;
}

/* ===== Toolbar ===== */
.tpl-toolbar {
  padding: 12px 20px 8px;
  display: flex;
  gap: 12px;
  align-items: flex-start;
  flex-shrink: 0;
  border-bottom: 1px solid #f6f6f6;
  background: #fafbfc;
}
.tpl-search-wrap {
  position: relative;
  flex-shrink: 0;
  width: 220px;
}
.tpl-search-icon {
  position: absolute;
  left: 10px;
  top: 50%;
  transform: translateY(-50%);
  color: #bfbfbf;
  font-size: 12px;
}
.tpl-search {
  width: 100%;
  height: 32px;
  padding: 0 10px 0 30px;
  background: #fff;
  border: 1px solid #d9d9d9;
  border-radius: 6px;
  color: #333;
  font-size: 13px;
  outline: none;
  transition: all 0.2s;
  box-sizing: border-box;
}
.tpl-search::placeholder { color: #bfbfbf; }
.tpl-search:focus {
  border-color: #13c2c2;
  box-shadow: 0 0 0 2px rgba(19, 194, 194, 0.12);
}

/* 分类胶囊 */
.tpl-cats {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  flex: 1;
}
.tpl-cat-btn {
  height: 28px;
  padding: 0 12px;
  border-radius: 14px;
  border: 1px solid #e8e8e8;
  background: #fff;
  color: #666;
  font-size: 12px;
  cursor: pointer;
  outline: none;
  transition: all 0.2s;
  white-space: nowrap;
}
.tpl-cat-btn:hover {
  border-color: #13c2c2;
  color: #13c2c2;
  background: #e6fffb;
}
.tpl-cat-btn.active {
  border-color: #13c2c2;
  background: #13c2c2;
  color: #fff;
}

/* ===== 网格 ===== */
.tpl-grid-wrap {
  flex: 1;
  overflow-y: auto;
  min-height: 0;          /* 关键：让 flex 子项可以收缩，否则溢出不可滚动 */
  padding: 16px 20px 24px; /* 底部多加一点，避免最后一张卡片贴边 */
  background: #f5f7f8;
}
.tpl-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(210px, 1fr));
  gap: 14px;
}

/* ===== 卡片：白底 + 青色 hover ===== */
.tpl-card {
  position: relative;
  border-radius: 8px;
  background: #fff;
  border: 1px solid #eee;
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.2s, border-color 0.2s, box-shadow 0.2s;
}
.tpl-card:hover, .tpl-card--hover {
  border-color: #13c2c2;
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(19, 194, 194, 0.15);
}

/* L 形角标：hover 时出现 */
.tpl-card-border {
  position: absolute;
  width: 10px;
  height: 10px;
  border-color: #13c2c2;
  border-style: solid;
  opacity: 0;
  transition: opacity 0.2s;
  z-index: 2;
}
.tpl-card:hover .tpl-card-border { opacity: 1; }
.tpl-b-tl { top: -1px; left: -1px; border-width: 1px 0 0 1px; }
.tpl-b-tr { top: -1px; right: -1px; border-width: 1px 1px 0 0; }
.tpl-b-bl { bottom: -1px; left: -1px; border-width: 0 0 1px 1px; }
.tpl-b-br { bottom: -1px; right: -1px; border-width: 0 1px 1px 0; }

/* 预览区：极浅灰底 + 网格线 */
.tpl-card-preview {
  height: 110px;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  background: #f8fafb;
}
.tpl-card-preview::before {
  content: '';
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(19, 194, 194, 0.05) 1px, transparent 1px),
    linear-gradient(90deg, rgba(19, 194, 194, 0.05) 1px, transparent 1px);
  background-size: 20px 20px;
}
.tpl-card-preview::after {
  content: '';
  position: absolute;
  width: 70px; height: 70px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(19, 194, 194, 0.08) 0%, transparent 70%);
}
.tpl-card-preview-icon {
  font-size: 36px;
  color: #13c2c2;
  position: relative;
  z-index: 1;
  transition: transform 0.3s;
  opacity: 0.7;
}
.tpl-card:hover .tpl-card-preview-icon {
  transform: scale(1.1);
  opacity: 1;
}

/* 对象数量 */
.tpl-card-count {
  position: absolute;
  bottom: 6px; right: 7px;
  background: rgba(0, 0, 0, 0.45);
  border-radius: 3px;
  font-size: 10px;
  color: #fff;
  padding: 1px 5px;
  z-index: 1;
}

/* 卡片内容 */
.tpl-card-body {
  padding: 10px 12px;
  background: #fff;
}
.tpl-card-name {
  font-size: 13px;
  font-weight: 600;
  color: #333;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.tpl-card-desc {
  font-size: 11px;
  color: #888;
  line-height: 1.45;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  margin-bottom: 8px;
  min-height: 30px;
}
.tpl-card-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.tpl-card-category {
  font-size: 10px;
  color: #13c2c2;
  border: 1px solid #b5f5ec;
  background: #e6fffb;
  border-radius: 3px;
  padding: 1px 5px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 70%;
}
.tpl-card-apply-btn {
  font-size: 11px;
  color: #fff;
  background: #13c2c2;
  border: none;
  border-radius: 4px;
  padding: 3px 10px;
  cursor: pointer;
  outline: none;
  transition: all 0.2s;
  flex-shrink: 0;
}
.tpl-card-apply-btn:hover {
  background: #36cfc9;
  box-shadow: 0 2px 6px rgba(19, 194, 194, 0.3);
}

/* ===== 空状态 ===== */
.tpl-empty {
  grid-column: 1 / -1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 0;
  color: #bfbfbf;
  font-size: 13px;
}

/* ===== 底部栏 ===== */
.tpl-footer {
  padding: 10px 20px;
  border-top: 1px solid #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-shrink: 0;
  background: #fafbfc;
}
.tpl-footer-hint {
  font-size: 12px;
  color: #aaa;
}
.tpl-footer-count {
  font-size: 12px;
  color: #999;
  font-family: monospace;
}

/* ===== 过渡动画 ===== */
.tpl-fade-enter-active { transition: opacity 0.2s, transform 0.2s; }
.tpl-fade-leave-active { transition: opacity 0.15s, transform 0.15s; }
.tpl-fade-enter { opacity: 0; }
.tpl-fade-leave-to { opacity: 0; }
.tpl-fade-enter .tpl-panel,
.tpl-fade-leave-to .tpl-panel { transform: scale(0.97) translateY(6px); }

/* ===== 滚动条（翠鸟风格） ===== */
.scrollbar::-webkit-scrollbar { width: 5px; }
.scrollbar::-webkit-scrollbar-thumb {
  border-radius: 10px;
  background: #91d5ff;
}
.scrollbar::-webkit-scrollbar-track {
  background: #f0f0f0;
}
</style>
