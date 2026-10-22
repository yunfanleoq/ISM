<template>
  <div class="xunzhi-viz-host">
    <a-spin v-if="loading" tip="正在解析组态页…" />
    <template v-else-if="error">
      <a-alert type="error" show-icon :message="error" />
      <p style="margin-top: 12px">
        请确认已导入机房环境项目，且存在 pageName
        <code>{{ pageName }}</code>。亦可从「数据中台」进入组态大屏。
      </p>
    </template>
    <iframe v-else class="viz-frame" :src="iframeSrc" title="xunzhi-jf-viz" />
  </div>
</template>

<script>
import { resolveDisplayPageByName, XUNZHI_JF_ENV_PLANE_PAGE } from '@/xunzhiJf/resolveDisplayPage'
import { primeXunzhiJfRouteContext } from '@/xunzhiJf/projectContext'
import {
  buildAppRunIframeSrc,
  buildEnvPlaneAppRunQuery,
  markEnvPlaneAppRunSession,
} from '@/xunzhiJf/appRunEmbed'

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
      const projectUuid = await primeXunzhiJfRouteContext(this.$store)
      const resolved = await resolveDisplayPageByName(this.pageName)
      if (!resolved) {
        const used = projectUuid || '空'
        this.error = `未找到组态页 ${this.pageName}（本次 ProjectUuid: ${used}。请确认已登录，且机房环境项目已导入，勿用历史 Cookie）`
        return
      }
      if (this.pageName === XUNZHI_JF_ENV_PLANE_PAGE) {
        markEnvPlaneAppRunSession(this.pageName)
      }
      const query =
        this.pageName === XUNZHI_JF_ENV_PLANE_PAGE
          ? buildEnvPlaneAppRunQuery(resolved.pageId)
          : resolved.pageId
            ? { pageId: resolved.pageId }
            : {}
      this.iframeSrc = buildAppRunIframeSrc(resolved.displayUid, query)
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
