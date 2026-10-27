<template>
  <div class="hx-dcim-frame">
    <div v-if="frameState !== 'ready'" class="hx-dcim-shell-loading">
      <div class="hx-dcim-shell-panel">
        <div class="hx-dcim-shell-title">{{ frameState === 'error' ? '新大屏未能打开' : '新大屏加载中' }}</div>
        <div class="hx-dcim-shell-sub">{{ shellHint }}</div>
      </div>
    </div>
    <iframe
      ref="frame"
      title="新大屏"
      src="/hx-dcim/index.html?v=20261027load#/power"
      @load="onFrameLoad"
      @error="onFrameError"
    />
  </div>
</template>

<script>
export default {
  name: 'HxDcimScreen',
  data () {
    return {
      frameState: 'loading',
      shellHint: '正在加载 hx-dcim 静态壳（测点索引约 25MB，首次可能 1～2 分钟）…',
      loadTimer: null
    }
  },
  mounted () {
    this.loadTimer = setTimeout(() => {
      if (this.frameState === 'loading') {
        this.shellHint = '加载时间较长：请确认已登录；在浏览器 Network 中查看 /hx-dcim/data/hx-field-index.json 是否 200。'
      }
    }, 90000)
  },
  beforeDestroy () {
    if (this.loadTimer) clearTimeout(this.loadTimer)
  },
  methods: {
    onFrameLoad () {
      if (this.loadTimer) clearTimeout(this.loadTimer)
      this.frameState = 'ready'
    },
    onFrameError () {
      if (this.loadTimer) clearTimeout(this.loadTimer)
      this.frameState = 'error'
      this.shellHint = 'iframe 加载失败。请检查 web/dist/hx-dcim/ 是否已部署，或直连 /hx-dcim/index.html#/power 排查。'
    }
  }
}
</script>

<style scoped>
.hx-dcim-frame {
  width: 100%;
  height: calc(100vh - 64px);
  background: #040c1c;
}
.hx-dcim-frame iframe {
  width: 100%;
  height: 100%;
  border: 0;
  background: #040c1c;
}
</style>
