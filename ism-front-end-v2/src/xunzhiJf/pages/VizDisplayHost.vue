<template>
  <div class="xunzhi-viz-host">
    <a-spin v-if="loading" tip="正在解析组态页…" />
    <template v-else-if="error">
      <a-alert type="error" show-icon :message="error" />
      <p style="margin-top: 12px">
        请确认已导入项目 <code>寻智机房环境-xunzhi-jf</code>，且存在 pageName
        <code>{{ pageName }}</code>。亦可从「数据中台」进入组态大屏。
      </p>
    </template>
    <iframe v-else class="viz-frame" :src="iframeSrc" title="xunzhi-jf-viz" />
  </div>
</template>

<script>
import { resolveDisplayPageByName, XUNZHI_JF_ENV_PLANE_PAGE } from '@/xunzhiJf/resolveDisplayPage'

export default {
  name: 'XunzhiJfVizDisplayHost',
  props: {
    pageName: { type: String, required: true },
  },
  data() {
    return {
      loading: true,
      error: '',
      iframeSrc: '',
    }
  },
  async mounted() {
    try {
      const resolved = await resolveDisplayPageByName(this.pageName)
      if (!resolved) {
        this.error = `未找到组态页 ${this.pageName}`
        return
      }
      const params = new URLSearchParams()
      if (resolved.pageId) params.set('pageId', resolved.pageId)
      if (this.pageName === XUNZHI_JF_ENV_PLANE_PAGE) {
        params.set('pageName', XUNZHI_JF_ENV_PLANE_PAGE)
        params.set('noRunTree', '1')
      }
      const q = params.toString() ? `?${params.toString()}` : ''
      const base = `${window.location.origin}${window.location.pathname}`
      this.iframeSrc = `${base}#/AppRun/${resolved.displayUid}${q}`
    } catch (e) {
      this.error = (e && e.message) || '解析组态失败'
    } finally {
      this.loading = false
    }
  },
}
</script>

<style scoped>
.xunzhi-viz-host {
  height: calc(100vh - 120px);
  min-height: 400px;
}
.viz-frame {
  width: 100%;
  height: 100%;
  border: 0;
}
</style>
