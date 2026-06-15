<template>
  <svg xmlns="http://www.w3.org/2000/svg"
       :width="size" :height="size"
       :viewBox="viewBox"
       preserveAspectRatio="xMidYMid meet"
       :style="{ overflow: 'visible' }">
    <g :style="svgStyle">
      <!-- ====== electric_1: 电容器 ====== -->
      <template v-if="type === 'capacitor'">
        <circle cx="5" cy="16" r="1.5" />
        <line x1="15" y1="8" x2="15" y2="25" />
        <line x1="22" y1="8" x2="22" y2="25" />
        <line x1="5" y1="16" x2="15" y2="16" />
      </template>

      <!-- ====== electric_2: 变压器 ====== -->
      <template v-else-if="type === 'transformer'">
        <circle cx="33" cy="4" r="1.5" />
        <rect fill="none" x="9" y="13" width="8" height="19" />
        <line x1="9" y1="15" x2="17" y2="15" />
        <line x1="9" y1="17" x2="17" y2="17" />
        <line x1="9" y1="28" x2="17" y2="28" />
        <line x1="9" y1="30" x2="17" y2="30" />
        <line x1="13" y1="32" x2="13" y2="36" />
        <line x1="9" y1="36" x2="17" y2="36" />
        <line x1="16" y1="38" x2="10" y2="38" />
        <line x1="12" y1="40" x2="14" y2="40" />
        <circle fill="none" cx="33" cy="41" r="6" />
        <circle fill="none" cx="33" cy="31" r="6" />
        <circle fill="none" cx="25" cy="36" r="6" />
        <line x1="33" y1="25" x2="33" y2="4" />
        <polyline fill="none" points="13,13 13,10 33,10" />
      </template>

      <!-- ====== electric_3: 母联柜（两水平线+连接） ====== -->
      <template v-else-if="type === 'bus-coupler'">
        <circle cx="16" cy="32" r="1.5" />
        <circle cx="16" cy="0" r="1.5" />
        <line x1="8" y1="15" x2="25" y2="15" />
        <line x1="8" y1="22" x2="25" y2="22" />
        <line x1="16" y1="0" x2="16" y2="15" />
        <line x1="16" y1="22" x2="16" y2="32" />
      </template>

      <!-- ====== electric_4: ATS 双电源 ====== -->
      <template v-else-if="type === 'ats'">
        <circle cx="16" cy="0" r="1.5" />
        <circle cx="16" cy="20" r="1.5" />
        <circle cx="16" cy="32" r="1.5" />
        <rect fill="none" x="8" y="20" width="16" height="12" />
        <line x1="16" y1="0" x2="16" y2="20" />
        <line x1="16" y1="32" x2="16" y2="32" />
        <line x1="12" y1="26" x2="20" y2="26" />
      </template>

      <!-- ====== electric_5: 开关/刀闸 ====== -->
      <template v-else-if="type === 'switchgear'">
        <circle cx="3" cy="16" r="1.5" />
        <circle cx="30" cy="16" r="1.5" />
        <line v-if="status === 'running' || status === 'alarm'" x1="3" y1="16" x2="28" y2="6" />
        <line v-else x1="3" y1="16" x2="30" y2="16" />
      </template>

      <!-- ====== electric_6: 变压器绕组 ====== -->
      <template v-else-if="type === 'transformer-winding'">
        <circle cx="16" cy="3" r="1.5" />
        <circle cx="16" cy="30" r="1.5" />
        <line x1="16" y1="3" x2="16" y2="10" />
        <line x1="4" y1="10" x2="29" y2="10" />
        <line x1="4" y1="18" x2="29" y2="18" />
        <line x1="16" y1="30" x2="16" y2="18" />
        <line x1="16" y1="10" x2="16" y2="18" />
      </template>

      <!-- ====== electric_7: 仪表/互感器 ====== -->
      <template v-else-if="type === 'meter-cabinet' || type === 'pt-cabinet'">
        <circle cx="32" cy="3" r="1.5" />
        <line x1="32" y1="10" x2="32" y2="3" />
        <rect x="43" y="18" width="6" height="18" />
        <line x1="43" y1="20" x2="48" y2="20" />
        <line x1="43" y1="22" x2="48" y2="22" />
        <line x1="43" y1="32" x2="49" y2="32" />
        <line x1="43" y1="34" x2="48" y2="34" />
        <line x1="18" y1="22" x2="13" y2="29" />
        <line x1="32" y1="41" x2="32" y2="47" />
        <line x1="28" y1="47" x2="36" y2="47" />
        <line x1="29" y1="49" x2="35" y2="49" />
        <line x1="31" y1="51" x2="33" y2="51" />
        <line x1="18" y1="10" x2="46" y2="10" />
        <line x1="46" y1="10" x2="46" y2="18" />
        <line x1="18" y1="10" x2="18" y2="22" />
        <line x1="18" y1="41" x2="18" y2="47" />
      </template>

      <!-- ====== electric_8: 馈线柜 / 通用设备 ====== -->
      <template v-else>
        <circle cx="16" cy="3" r="1.5" />
        <circle cx="16" cy="30" r="1.5" />
        <rect fill="none" x="9" y="3" width="15" height="27" />
      </template>
    </g>
  </svg>
</template>

<script>
import { DEVICE_TYPE_COLORS } from './config'

export default {
  name: 'ElectricIcon',
  props: {
    type: { type: [String, Number], default: 'switchgear' },
    status: { type: String, default: 'running' },
    size: { type: Number, default: 32 },
  },
  computed: {
    viewBox() {
      if (this.type === 'transformer') return '0 0 48 48'
      if (this.type === 'meter-cabinet' || this.type === 'pt-cabinet') return '0 0 64 64'
      if (this.type === 'bus-coupler' || this.type === 'ats') return '0 0 32 32'
      return '0 0 32 32'
    },
    color() {
      return DEVICE_TYPE_COLORS[this.type] || '#3498db'
    },
    svgStyle() {
      const statusColor = this.status === 'alarm' ? '#ef4444'
        : this.status === 'stopped' ? '#64748b'
        : this.status === 'offline' ? '#475569'
        : this.color
      return {
        opacity: 1,
        'stroke-opacity': 1,
        stroke: statusColor,
        'stroke-width': 1.2,
        'stroke-linecap': 'round',
        'stroke-linejoin': 'round',
        fill: 'none',
      }
    },
  },
}
</script>
