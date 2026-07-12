<template>
  <div class="runtime-card-grid" :class="`runtime-card-grid--${mode}`">
    <button
      v-for="item in items"
      :key="item.key"
      type="button"
      class="runtime-card"
      :class="{ 'runtime-card--online': item.online, 'runtime-card--point': mode === 'point' }"
      :disabled="mode !== 'device'"
      :title="item.title || item.name"
      @mousedown="selectFromPointer(item.source)"
      @click="selectFromClick(item.source)"
    >
      <span class="runtime-card__accent" aria-hidden="true"></span>
      <span class="runtime-card__icon" aria-hidden="true">{{ mode === 'device' ? '▰' : '⌁' }}</span>
      <span class="runtime-card__content">
        <strong>{{ item.name || (mode === 'device' ? '未命名设备' : '未命名测点') }}</strong>
        <template v-if="mode === 'device'">
          <small>{{ item.code || '暂无设备编号' }}</small>
          <small v-if="item.model">模型 {{ item.model }}</small>
        </template>
        <template v-else>
          <span class="runtime-card__reading">
            <transition name="reading-update" mode="out-in">
              <b :key="String(item.value)">{{ formatValue(item.value) }}</b>
            </transition>
            <em>{{ item.unit || '—' }}</em>
          </span>
          <small>REALTIME DATA</small>
        </template>
      </span>
      <span v-if="mode === 'device'" class="runtime-card__status">
        {{ item.online ? '在线' : '离线' }}
      </span>
      <span v-if="mode === 'device'" class="runtime-card__enter">查看点位 ›</span>
      <span v-else class="runtime-card__live"><i></i> LIVE</span>
    </button>
  </div>
</template>

<script>
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
  methods: {
    formatValue(value) {
      return value === undefined || value === null || value === '' ? '—' : value
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
  flex: 1 1 auto;
  min-height: 0;
  display: grid;
  gap: 12px;
  padding: 14px;
  overflow: hidden;
  box-sizing: border-box;
}

.runtime-card-grid--device {
  grid-template-columns: repeat(7, minmax(0, 1fr));
  grid-template-rows: repeat(7, minmax(0, 1fr));
}

.runtime-card-grid--point {
  grid-template-columns: repeat(10, minmax(0, 1fr));
  grid-template-rows: repeat(8, minmax(0, 1fr));
  gap: 8px;
  padding: 10px;
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
  display: grid;
  place-items: center;
  width: 42px;
  height: 42px;
  border: 1px solid rgba(55, 226, 255, 0.52);
  color: #52e8ff;
  background: rgba(0, 190, 225, 0.1);
  clip-path: polygon(50% 0, 100% 24%, 100% 76%, 50% 100%, 0 76%, 0 24%);
}

.runtime-card__content {
  min-width: 0;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.runtime-card__content strong,
.runtime-card__content small {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.runtime-card__content strong {
  font-size: 14px;
  letter-spacing: 0.3px;
}

.runtime-card__content small {
  color: #6f98ad;
  font-size: 10px;
}

.runtime-card__status,
.runtime-card__live {
  position: absolute;
  top: 10px;
  right: 12px;
  color: #8497aa;
  font-size: 10px;
}

.runtime-card__status::before {
  content: "●";
  margin-right: 4px;
}

.runtime-card--online .runtime-card__status {
  color: #58e0ad;
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
  border-radius: 5px;
  box-shadow: inset 0 0 14px rgba(0, 195, 235, 0.05), 0 3px 8px rgba(0, 0, 0, 0.18);
}

.runtime-card--point .runtime-card__icon {
  width: 24px;
  height: 24px;
  font-size: 11px;
}

.runtime-card--point .runtime-card__content {
  gap: 3px;
}

.runtime-card--point .runtime-card__content strong {
  padding-right: 27px;
  font-size: 11px;
  line-height: 1.2;
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
}

.runtime-card__reading em {
  color: #79a9ba;
  font-size: 9px;
  font-style: normal;
}

.runtime-card__live {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #58e0ad;
  letter-spacing: 0.5px;
}

.runtime-card--point .runtime-card__live {
  top: 7px;
  right: 8px;
  font-size: 7px;
}

.runtime-card--point .runtime-card__content small {
  font-size: 7px;
  letter-spacing: 0.7px;
}

.runtime-card__live i {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: currentColor;
  box-shadow: 0 0 7px currentColor;
  animation: livePulse 2.4s ease-in-out infinite;
}

.reading-update-enter-active,
.reading-update-leave-active {
  transition: opacity 0.18s ease, transform 0.18s ease;
}

.reading-update-enter,
.reading-update-leave-to {
  opacity: 0.35;
  transform: translateY(3px);
}

@keyframes livePulse {
  0%, 100% { opacity: 0.45; }
  50% { opacity: 1; }
}

@media (prefers-reduced-motion: reduce) {
  .runtime-card,
  .reading-update-enter-active,
  .reading-update-leave-active {
    transition: none;
  }

  .runtime-card__live i {
    animation: none;
  }
}
</style>
