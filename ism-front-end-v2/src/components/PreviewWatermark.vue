<template>
  <div v-if="showWatermark" class="preview-watermark" aria-hidden="true">
    <span>{{ text }}</span>
    <span>{{ subText }}</span>
  </div>
</template>

<script>
export default {
  name: 'PreviewWatermark',
  props: {
    text: {
      type: String,
      default: 'www.ismctl.com'
    },
    subText: {
      type: String,
      default: '免费试用版'
    }
  },
  computed: {
    showWatermark() {
      const setting = this.$store && this.$store.state ? this.$store.state.setting || {} : {}
      const isOEM = this.normalizeLicenseFlag(setting.IsOEM)
      if (isOEM === true) return false
      if (isOEM === false) return true
      return false
    }
  },
  methods: {
    normalizeLicenseFlag(value) {
      if (value === true || value === 1) return true
      if (value === false || value === 0) return false
      if (typeof value === 'string') {
        const normalized = value.trim().toLowerCase()
        if (['true', '1', 'yes', 'y', 'authorized', 'license', 'licensed'].indexOf(normalized) !== -1) return true
        if (['false', '0', 'no', 'n', 'unauthorized', 'unlicensed'].indexOf(normalized) !== -1) return false
      }
      return null
    }
  }
}
</script>

<style scoped>
.preview-watermark {
  position: fixed;
  right: 18px;
  bottom: 18px;
  z-index: 100000;
  pointer-events: none;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 2px;
  padding: 8px 12px;
  border-radius: 6px;
  border: 1px solid rgba(24, 144, 255, 0.28);
  background: rgba(255, 255, 255, 0.75);
  box-shadow: 0 6px 18px rgba(20, 35, 52, 0.12);
  backdrop-filter: blur(6px);
  user-select: none;
  /* 整体微呼吸 */
  animation: wm-breathe 3s ease-in-out infinite;
}

.preview-watermark span {
  color: rgba(46, 55, 70, 0.72);
  font-size: 12px;
  font-weight: 600;
  line-height: 1.35;
  letter-spacing: 0.3px;
  text-align: right;
  white-space: nowrap;
  /* 文字柔和发光扫过 — 不影响阅读 */
  animation: wm-text-glow 2.5s ease-in-out infinite;
}

/* 整体呼吸：更明显的透明度 + 边框光晕起伏 */
@keyframes wm-breathe {
  0%, 100% {
    opacity: 0.85;
    border-color: rgba(24, 144, 255, 0.10);
    box-shadow: 0 4px 12px rgba(20, 35, 52, 0.06);
  }
  50% {
    opacity: 1;
    border-color: rgba(24, 144, 255, 0.55);
    box-shadow: 0 6px 24px rgba(24, 144, 255, 0.18);
  }
}

/* 文字发光脉冲：蓝色光晕周期性扫过 */
@keyframes wm-text-glow {
  0%, 100% {
    text-shadow: 0 0 0 transparent, 0 1px 0 rgba(255, 255, 255, 0.50);
    color: rgba(46, 55, 70, 0.72);
  }
  50% {
    text-shadow: 0 0 14px rgba(24, 144, 255, 0.35), 0 0 4px rgba(24, 144, 255, 0.20), 0 1px 0 rgba(255, 255, 255, 0.55);
    color: rgba(24, 100, 220, 0.85);
  }
}
</style>
