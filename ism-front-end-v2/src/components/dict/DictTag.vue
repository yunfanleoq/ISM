<template>
  <a-tag v-if="shown" :style="tagStyle">{{ shown }}</a-tag>
</template>

<script>
import {dictColor, dictLabel} from '@/utils/systemDict'

export default {
  name: 'DictTag',
  props: {
    dictCode: {type: String, required: true},
    value: {type: [String, Number], default: ''},
  },
  computed: {
    shown() {
      return dictLabel(this.dictCode, this.value)
    },
    tagStyle() {
      const bg = dictColor(this.dictCode, this.value) || '#0099FF'
      return {
        backgroundColor: bg,
        color: textOn(bg),
        borderColor: bg,
      }
    },
  },
}

function textOn(color) {
  const hex = String(color || '').trim()
  const named = {yellow: '#FFFF00', orange: '#FFA500', red: '#FF0000'}
  const raw = named[hex.toLowerCase()] || hex
  const m = raw.replace('#', '')
  if (!/^[0-9a-fA-F]{6}$/.test(m)) return '#fff'
  const r = parseInt(m.slice(0, 2), 16)
  const g = parseInt(m.slice(2, 4), 16)
  const b = parseInt(m.slice(4, 6), 16)
  return (0.299 * r + 0.587 * g + 0.114 * b) > 180 ? '#222' : '#fff'
}
</script>
