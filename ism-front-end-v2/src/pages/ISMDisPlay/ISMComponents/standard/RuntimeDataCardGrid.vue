<template>
  <div class="runtime-card-grid" :class="`runtime-card-grid--${mode}`">
    <button
      v-for="item in items"
      :key="item.key"
      type="button"
      class="runtime-card"
      :class="{
        'runtime-card--device': mode === 'device',
        'runtime-card--point': mode === 'point',
        'runtime-card--active': mode === 'point' && isPointActive(item)
      }"
      :disabled="mode !== 'device'"
      :title="item.title || item.name"
      @mousedown="selectFromPointer(item.source)"
      @click="selectFromClick(item.source)"
    >
      <span class="runtime-card__accent" aria-hidden="true"></span>
      <span
        class="runtime-card__icon"
        :class="{ 'runtime-card__icon--active': mode === 'point' && isPointActive(item) }"
        aria-hidden="true"
      >
        <span class="runtime-card__icon-glyph">
          <i></i><i></i><i></i>
        </span>
      </span>
      <span class="runtime-card__content">
        <strong :class="['runtime-card__name', nameSizeClass(item)]">{{ formatName(item) }}</strong>
        <template v-if="mode !== 'device'">
          <span class="runtime-card__reading">
            <b>{{ formatValue(item.value, item) }}</b>
            <em v-if="hasUnit(item)">{{ item.unit }}</em>
          </span>
        </template>
      </span>
      <!-- 20260729：取消设备在线态文案与变色，卡片样式统一 -->
      <span v-if="mode === 'device'" class="runtime-card__enter">查看点位 ›</span>
    </button>
  </div>
</template>

<script>
import { formatPointDisplayValue } from '@/pages/ISMDisPlay/utils/pointValueDisplay'

const ACTIVE_MS = 3000

export default {
  name: 'RuntimeDataCardGrid',
  props: {
    mode: {
      type: String,
      default: 'device',
      validator: value => ['device', 'point'].includes(value),
    },
    items: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      activeMap: Object.create(null),
      _prevValues: Object.create(null),
      _activeTimers: Object.create(null),
    }
  },
  watch: {
    // 避免 deep watch 在每帧数值刷新时遍历全部卡片属性（麒麟现场卡顿主因之一）
    items: {
      handler(items) {
        const list = Array.isArray(items) ? items : []
        let fingerprint = `${list.length}`
        for (let i = 0; i < list.length; i++) {
          const it = list[i]
          if (!it) continue
          fingerprint += `|${it.key}:${it.value}`
        }
        if (fingerprint === this._itemsFingerprint) return
        this._itemsFingerprint = fingerprint
        this.syncActiveFromItems(list)
      },
      immediate: true,
    },
    mode() {
      this.clearAllActive()
      this._prevValues = Object.create(null)
      this._itemsFingerprint = ''
      this.syncActiveFromItems(this.items)
    },
  },
  beforeDestroy() {
    this.clearAllActive()
  },
  methods: {
    hasUnit(item) {
      const unit = item && item.unit
      if (unit == null) return false
      const text = String(unit).trim()
      return !!text && text !== '—' && text !== '-'
    },
    isPointActive(item) {
      return !!(item && item.key != null && this.activeMap[item.key])
    },
    markActive(key) {
      if (this._activeTimers[key]) clearTimeout(this._activeTimers[key])
      this.$set(this.activeMap, key, true)
      this._activeTimers[key] = setTimeout(() => {
        this.$delete(this.activeMap, key)
        delete this._activeTimers[key]
      }, ACTIVE_MS)
    },
    clearActive(key) {
      if (this._activeTimers[key]) {
        clearTimeout(this._activeTimers[key])
        delete this._activeTimers[key]
      }
      if (this.activeMap[key]) this.$delete(this.activeMap, key)
    },
    clearAllActive() {
      Object.keys(this._activeTimers).forEach(key => {
        clearTimeout(this._activeTimers[key])
      })
      this._activeTimers = Object.create(null)
      this.activeMap = Object.create(null)
    },
    syncActiveFromItems(items) {
      if (this.mode !== 'point') return
      const list = Array.isArray(items) ? items : []
      const nextKeys = new Set()
      list.forEach(item => {
        if (!item || item.key == null) return
        const key = item.key
        nextKeys.add(key)
        const val = item.value
        if (!(key in this._prevValues)) {
          // 首次出现只记基线，不算变化（避免整页同时亮）
          this._prevValues[key] = val
          return
        }
        if (String(this._prevValues[key]) !== String(val)) {
          this._prevValues[key] = val
          this.markActive(key)
        }
      })
      Object.keys(this._prevValues).forEach(key => {
        if (!nextKeys.has(key)) delete this._prevValues[key]
      })
      Object.keys(this.activeMap).forEach(key => {
        if (!nextKeys.has(key)) this.clearActive(key)
      })
    },
    rawDisplayName(item) {
      return String(
        (item && (item.displayName || item.name))
        || (this.mode === 'device' ? '未命名设备' : '未命名测点')
      ).trim()
    },
    nameSizeClass(item) {
      const length = Array.from(this.rawDisplayName(item)).length
      const longAt = this.mode === 'device' ? 16 : 13
      const compactAt = this.mode === 'device' ? 23 : 19
      return {
        'runtime-card__name--long': length > longAt,
        'runtime-card__name--compact': length > compactAt,
      }
    },
    formatName(item) {
      const name = this.rawDisplayName(item)
      const chars = Array.from(name)
      const limit = this.mode === 'device' ? 27 : 23
      if (chars.length <= limit) return name
      const tailLength = this.mode === 'device' ? 9 : 7
      const headLength = limit - tailLength - 1
      return `${chars.slice(0, headLength).join('')}…${chars.slice(-tailLength).join('')}`
    },
    formatValue(value, item) {
      const name = item && (item.name || item.displayName || item.title)
      const shown = formatPointDisplayValue(name, value)
      return shown === undefined || shown === null || shown === '' ? '—' : shown
    },
    selectFromPointer(source) {
      if (this.mode !== 'device') return
      this._selectedByPointer = true
      this.$emit('select', source)
    },
    selectFromClick(source) {
      if (this.mode !== 'device') return
      if (this._selectedByPointer) {
        this._selectedByPointer = false
        return
      }
      this.$emit('select', source)
    },
  },
}
</script>

<style scoped>
.runtime-card-grid {
  position: relative;
  flex: 1 1 auto;
  min-height: 0;
  display: grid;
  gap: 12px;
  padding: 14px;
  overflow: hidden;
  box-sizing: border-box;
  border: var(--panelBorderWidth, 1px) solid var(--panelBorder, #263449);
  border-radius: var(--panelBorderRadius, 10px);
  background: var(--panelBg, linear-gradient(180deg, #111827 0%, #0f172a 100%));
  box-shadow:
    var(--panelShadow, 0 14px 40px rgba(0, 0, 0, 0.35)),
    inset 0 0 0 1px rgba(54, 211, 243, 0.2),
    inset 0 0 32px rgba(20, 184, 220, 0.055);
}

.runtime-card-grid::before,
.runtime-card-grid::after {
  content: "";
  position: absolute;
  z-index: 3;
  width: 38px;
  height: 13px;
  pointer-events: none;
}

.runtime-card-grid::before {
  top: -1px;
  left: -1px;
  border-top: 2px solid var(--panelAccent, #4ae6ff);
  border-left: 2px solid var(--panelAccent, #4ae6ff);
  border-radius: var(--panelBorderRadius, 10px) 0 0;
}

.runtime-card-grid::after {
  right: -1px;
  bottom: -1px;
  border-right: 2px solid var(--panelAccent, #4ae6ff);
  border-bottom: 2px solid var(--panelAccent, #4ae6ff);
  border-radius: 0 0 var(--panelBorderRadius, 10px);
}

.runtime-card-grid--device {
  grid-template-columns: repeat(7, minmax(0, 1fr));
  grid-template-rows: repeat(7, minmax(0, 1fr));
  /* 透明底：只保留卡片本身，外框交给页面 box13，避免网格边框溢出压边 */
  border: none;
  border-radius: 0;
  background: transparent;
  box-shadow: none;
  padding: 8px;
}

.runtime-card-grid--device::before,
.runtime-card-grid--device::after {
  display: none;
}

.runtime-card-grid--point {
  grid-template-columns: repeat(10, minmax(0, 1fr));
  grid-template-rows: repeat(8, minmax(0, 1fr));
  gap: 8px;
  padding: 10px;
  border: var(--panelBorderWidth, 1px) solid var(--panelBorder, #2fd5f2);
  box-shadow:
    var(--panelShadow, 0 14px 40px rgba(0, 0, 0, 0.35)),
    inset 0 0 0 4px rgba(3, 18, 29, 0.72),
    inset 0 0 0 5px rgba(74, 230, 255, 0.09),
    inset 0 0 32px rgba(20, 184, 220, 0.055);
}

.runtime-card-grid--point::before,
.runtime-card-grid--point::after {
  display: block;
  width: 82px;
  height: 25px;
  filter:
    drop-shadow(0 0 3px rgba(74, 230, 255, 0.5))
    drop-shadow(0 0 8px rgba(47, 213, 242, 0.18));
  opacity: 1;
}

.runtime-card-grid--point::before {
  top: -1px;
  left: -1px;
  border-top-width: 3px;
  border-left-width: 3px;
  background:
    linear-gradient(90deg, var(--panelAccent, #4ae6ff) 0 22px, transparent 22px) 8px 7px / 100% 2px no-repeat,
    linear-gradient(180deg, var(--panelAccent, #4ae6ff) 0 10px, transparent 10px) 8px 7px / 2px 100% no-repeat;
}

.runtime-card-grid--point::after {
  right: -1px;
  bottom: -1px;
  border-right-width: 3px;
  border-bottom-width: 3px;
  background:
    linear-gradient(270deg, var(--panelAccent, #4ae6ff) 0 22px, transparent 22px) right 8px bottom 7px / 100% 2px no-repeat,
    linear-gradient(0deg, var(--panelAccent, #4ae6ff) 0 10px, transparent 10px) right 8px bottom 7px / 2px 100% no-repeat;
}

.runtime-card {
  position: relative;
  min-width: 0;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  overflow: hidden;
  border: 1px solid rgba(44, 190, 231, 0.38);
  border-radius: 8px;
  background: linear-gradient(135deg, rgba(9, 47, 68, 0.96), rgba(6, 25, 40, 0.98));
  color: #dffaff;
  text-align: left;
  box-shadow: inset 0 0 20px rgba(0, 195, 235, 0.06), 0 5px 12px rgba(0, 0, 0, 0.2);
  transition: transform 0.16s ease, border-color 0.16s ease, box-shadow 0.16s ease;
}

.runtime-card:not(:disabled) {
  cursor: pointer;
}

.runtime-card:disabled {
  cursor: default;
}

.runtime-card > * {
  pointer-events: none;
}

.runtime-card:not(:disabled):hover,
.runtime-card:not(:disabled):focus-visible {
  outline: none;
  transform: translateY(-2px);
  border-color: rgba(83, 235, 255, 0.82);
  box-shadow: inset 0 0 24px rgba(0, 210, 255, 0.11), 0 0 16px rgba(0, 210, 255, 0.16);
}

.runtime-card__accent {
  position: absolute;
  left: 58px;
  right: 12px;
  bottom: 0;
  height: 1px;
  background: linear-gradient(90deg, #19dfff, transparent);
}

.runtime-card__icon {
  flex: none;
  position: relative;
  display: grid;
  place-items: center;
  width: 42px;
  height: 42px;
  overflow: hidden;
  color: var(--deviceIconAccent, #52e8ff);
  background: rgba(0, 28, 42, 0.92);
  clip-path: polygon(50% 0, 100% 24%, 100% 76%, 50% 100%, 0 76%, 0 24%);
}

.runtime-card__icon::before {
  content: "";
  position: absolute;
  inset: -55%;
  background: conic-gradient(
    from 0deg,
    transparent 0deg 235deg,
    currentColor 270deg,
    rgba(255, 255, 255, 0.95) 300deg,
    transparent 332deg
  );
  animation: iconBorderFlow 2.4s linear infinite;
}

.runtime-card__icon::after {
  content: "";
  position: absolute;
  inset: 2px;
  z-index: 1;
  background: linear-gradient(145deg, rgba(5, 34, 48, 0.98), rgba(4, 17, 28, 0.98));
  clip-path: inherit;
}

.runtime-card__icon-glyph {
  position: relative;
  z-index: 2;
  width: 19px;
  height: 17px;
  box-sizing: border-box;
  border: 1px solid currentColor;
  border-radius: 3px;
  box-shadow: 0 0 9px currentColor;
}

.runtime-card__icon-glyph i {
  position: absolute;
  left: 3px;
  right: 3px;
  height: 1px;
  background: currentColor;
  box-shadow: 0 0 4px currentColor;
}

.runtime-card__icon-glyph i:nth-child(1) { top: 4px; }
.runtime-card__icon-glyph i:nth-child(2) { top: 8px; }
.runtime-card__icon-glyph i:nth-child(3) { top: 12px; }

.runtime-card--device .runtime-card__icon {
  width: 44px;
  height: 44px;
  overflow: hidden;
  border: 1px solid rgba(82, 232, 255, 0.42);
  border-radius: 5px;
  background:
    linear-gradient(90deg, transparent 49%, rgba(82, 232, 255, 0.07) 50%, transparent 51%),
    linear-gradient(0deg, transparent 49%, rgba(82, 232, 255, 0.07) 50%, transparent 51%),
    linear-gradient(145deg, rgba(8, 43, 58, 0.98), rgba(3, 18, 29, 0.98));
  clip-path: polygon(7px 0, 100% 0, 100% calc(100% - 7px), calc(100% - 7px) 100%, 0 100%, 0 7px);
  box-shadow:
    inset 0 0 0 1px rgba(4, 15, 24, 0.85),
    inset 0 0 12px rgba(39, 203, 235, 0.06);
}

.runtime-card--device .runtime-card__icon::before {
  inset: 4px;
  border: 1px solid rgba(82, 232, 255, 0.14);
  background:
    linear-gradient(90deg, currentColor 0 4px, transparent 4px) top left / 8px 1px no-repeat,
    linear-gradient(180deg, currentColor 0 4px, transparent 4px) top left / 1px 8px no-repeat,
    linear-gradient(270deg, currentColor 0 4px, transparent 4px) bottom right / 8px 1px no-repeat,
    linear-gradient(0deg, currentColor 0 4px, transparent 4px) bottom right / 1px 8px no-repeat;
  opacity: 0.55;
  animation: none;
}

.runtime-card--device .runtime-card__icon::after {
  inset: 0;
  border: 0;
  background: linear-gradient(
    110deg,
    transparent 0 46%,
    rgba(111, 240, 255, 0.055) 46% 52%,
    transparent 52% 100%
  );
  clip-path: inherit;
}

.runtime-card--device .runtime-card__icon-glyph {
  width: 20px;
  height: 24px;
  border-color: rgba(133, 241, 255, 0.88);
  border-radius: 2px;
  background: linear-gradient(180deg, rgba(30, 111, 132, 0.22), rgba(3, 20, 31, 0.65));
  box-shadow:
    inset 0 0 0 2px rgba(2, 20, 30, 0.8),
    0 0 5px rgba(82, 232, 255, 0.16);
}

.runtime-card--device .runtime-card__icon-glyph::before {
  content: "";
  position: absolute;
  top: 4px;
  left: 4px;
  width: 3px;
  height: 3px;
  border-radius: 50%;
  background: #64727c;
  box-shadow: 5px 0 0 rgba(100, 114, 124, 0.38), 10px 0 0 rgba(100, 114, 124, 0.2);
}

.runtime-card--device.runtime-card--online .runtime-card__icon-glyph::before {
  background: #70f0b5;
  box-shadow:
    0 0 4px rgba(112, 240, 181, 0.52),
    5px 0 0 rgba(112, 240, 181, 0.32),
    10px 0 0 rgba(112, 240, 181, 0.16);
  animation: deviceStatusPulse 2.2s ease-in-out infinite;
}

.runtime-card--device .runtime-card__icon-glyph::after {
  content: "";
  position: absolute;
  right: 4px;
  bottom: 3px;
  width: 2px;
  height: 2px;
  background: currentColor;
  box-shadow: -4px 0 0 rgba(82, 232, 255, 0.4);
}

.runtime-card--device .runtime-card__icon-glyph i {
  left: 4px;
  right: 4px;
  height: 2px;
  border: 1px solid rgba(82, 232, 255, 0.48);
  background: rgba(82, 232, 255, 0.1);
  box-shadow: none;
}

.runtime-card--device .runtime-card__icon-glyph i:nth-child(1) { top: 10px; }
.runtime-card--device .runtime-card__icon-glyph i:nth-child(2) { top: 14px; }
.runtime-card--device .runtime-card__icon-glyph i:nth-child(3) { display: none; }

.runtime-card__content {
  min-width: 0;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.runtime-card__content small {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.runtime-card__name {
  display: -webkit-box;
  overflow: hidden;
  overflow-wrap: anywhere;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  white-space: normal;
  font-size: 14px;
  line-height: 1.25;
  letter-spacing: 0.3px;
}

.runtime-card--device .runtime-card__name {
  padding-right: 38px;
}

.runtime-card__name--long {
  font-size: 12px;
  letter-spacing: 0.15px;
}

.runtime-card__name--compact {
  font-size: 11px;
  line-height: 1.18;
  letter-spacing: 0;
}

.runtime-card__content small {
  color: #6f98ad;
  font-size: 10px;
}

.runtime-card__enter {
  position: absolute;
  right: 12px;
  bottom: 9px;
  color: #4caec9;
  font-size: 10px;
}

.runtime-card--point {
  align-items: center;
  gap: 7px;
  min-height: 0;
  padding: 8px 9px;
  border-radius: 8px;
  box-shadow: inset 0 0 14px rgba(0, 195, 235, 0.045), 0 3px 8px rgba(0, 0, 0, 0.18);
}

.runtime-card--point .runtime-card__icon {
  width: 28px;
  height: 28px;
  overflow: hidden;
  color: var(--pointIconAccent, #7a8a9a);
  border: 1px solid rgba(82, 232, 255, 0.22);
  border-radius: 4px;
  background:
    linear-gradient(90deg, transparent 49%, rgba(82, 232, 255, 0.04) 50%, transparent 51%),
    linear-gradient(0deg, transparent 49%, rgba(82, 232, 255, 0.04) 50%, transparent 51%),
    linear-gradient(145deg, rgba(10, 39, 57, 0.98), rgba(5, 19, 32, 0.98));
  clip-path: polygon(5px 0, 100% 0, 100% calc(100% - 5px), calc(100% - 5px) 100%, 0 100%, 0 5px);
  box-shadow:
    inset 0 0 0 1px rgba(4, 15, 24, 0.82),
    inset 0 0 9px rgba(82, 232, 255, 0.05);
  transition: color 0.2s ease, border-color 0.2s ease, box-shadow 0.2s ease;
}

.runtime-card--point .runtime-card__icon::before {
  inset: 3px;
  border: 1px solid rgba(82, 232, 255, 0.11);
  background:
    linear-gradient(90deg, currentColor 0 3px, transparent 3px) top left / 6px 1px no-repeat,
    linear-gradient(180deg, currentColor 0 3px, transparent 3px) top left / 1px 6px no-repeat,
    linear-gradient(270deg, currentColor 0 3px, transparent 3px) bottom right / 6px 1px no-repeat,
    linear-gradient(0deg, currentColor 0 3px, transparent 3px) bottom right / 1px 6px no-repeat;
  opacity: 0.35;
  animation: none;
}

.runtime-card--point .runtime-card__icon::after {
  inset: 0;
  border-radius: 0;
  background: linear-gradient(
    110deg,
    transparent 0 46%,
    rgba(188, 165, 255, 0.03) 46% 53%,
    transparent 53% 100%
  );
  clip-path: inherit;
}

.runtime-card--point .runtime-card__icon-glyph {
  width: 16px;
  height: 16px;
  border: 0;
  border-radius: 0;
  box-shadow: none;
}

.runtime-card--point .runtime-card__icon-glyph::before {
  content: "";
  position: absolute;
  top: 50%;
  left: 50%;
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: currentColor;
  box-shadow: 0 0 4px rgba(122, 138, 154, 0.28);
  transform: translate(-50%, -50%);
  opacity: 0.55;
  animation: none;
}

.runtime-card--point .runtime-card__icon-glyph::after {
  content: "";
  position: absolute;
  top: 50%;
  left: 1px;
  right: 1px;
  height: 1px;
  background: linear-gradient(90deg, transparent, currentColor 35%, transparent 35% 65%, currentColor 65%, transparent);
  opacity: 0.22;
  transform: translateY(-50%);
}

.runtime-card--point .runtime-card__icon-glyph i {
  top: 50%;
  bottom: auto;
  height: 12px;
  background: transparent;
  box-shadow: none;
  opacity: 0.55;
}

.runtime-card--point .runtime-card__icon-glyph i:nth-child(1) {
  left: 4px;
  right: 4px;
  width: auto;
  height: 8px;
  border: 1px solid currentColor;
  border-radius: 50%;
  opacity: 0.4;
  transform: translateY(-50%);
}

.runtime-card--point .runtime-card__icon-glyph i:nth-child(2) {
  left: 0;
  right: auto;
  width: 5px;
  border: 1px solid currentColor;
  border-right: 0;
  border-radius: 10px 0 0 10px;
  transform: translateY(-50%);
}

.runtime-card--point .runtime-card__icon-glyph i:nth-child(3) {
  right: 0;
  left: auto;
  width: 5px;
  border: 1px solid currentColor;
  border-left: 0;
  border-radius: 0 10px 10px 0;
  transform: translateY(-50%);
}

/* 值变化活跃：左侧图标高亮脉冲约 3s */
.runtime-card--point .runtime-card__icon--active {
  color: #58e0ad;
  border-color: rgba(88, 224, 173, 0.72);
  box-shadow:
    inset 0 0 0 1px rgba(4, 15, 24, 0.82),
    inset 0 0 12px rgba(88, 224, 173, 0.18),
    0 0 10px rgba(88, 224, 173, 0.28);
}

.runtime-card--point .runtime-card__icon--active::before {
  opacity: 0.85;
  border-color: rgba(88, 224, 173, 0.45);
}

.runtime-card--point .runtime-card__icon--active .runtime-card__icon-glyph::before {
  background: #70f0b5;
  box-shadow: 0 0 8px rgba(112, 240, 181, 0.65);
  opacity: 1;
  animation: telemetryPulse 1.1s ease-in-out infinite;
}

.runtime-card--point .runtime-card__content {
  gap: 4px;
}

.runtime-card--point .runtime-card__content strong {
  padding-right: 0;
  font-size: 11px;
  line-height: 1.2;
}

.runtime-card--point .runtime-card__name--long {
  font-size: 10px;
}

.runtime-card--point .runtime-card__name--compact {
  font-size: 9px;
  line-height: 1.15;
}

.runtime-card__reading {
  display: flex;
  align-items: baseline;
  gap: 5px;
  min-width: 0;
}

.runtime-card__reading b {
  max-width: 100%;
  overflow: hidden;
  color: #69f2ff;
  font-size: clamp(15px, 0.95vw, 20px);
  line-height: 1.1;
  text-overflow: ellipsis;
  text-shadow: 0 0 12px rgba(53, 225, 255, 0.28);
  white-space: nowrap;
  transition: opacity 0.12s ease;
}

.runtime-card__reading em {
  color: #79a9ba;
  font-size: 9px;
  font-style: normal;
}

@keyframes iconBorderFlow {
  to { transform: rotate(360deg); }
}

@keyframes telemetryPulse {
  0%, 100% {
    opacity: 0.62;
    transform: translate(-50%, -50%) scale(0.86);
  }
  50% {
    opacity: 1;
    transform: translate(-50%, -50%) scale(1);
  }
}

@keyframes deviceStatusPulse {
  0%, 100% {
    opacity: 0.5;
  }
  50% {
    opacity: 1;
  }
}

@media (prefers-reduced-motion: reduce) {
  .runtime-card {
    transition: none;
  }

  .runtime-card__reading b {
    transition: none;
  }

  .runtime-card__icon::before,
  .runtime-card--device.runtime-card--online .runtime-card__icon-glyph::before,
  .runtime-card--point .runtime-card__icon--active .runtime-card__icon-glyph::before {
    animation: none;
  }
}
</style>
