<template>
  <!-- 仅具备后台管理权限（Admin 角色）的用户渲染；其他角色 DOM 不出现 -->
  <div
    v-if="isAdmin"
    class="back-to-admin-btn"
    @click="goAdmin"
    title="项目列表"
  >
    <svg class="bta-gear" viewBox="0 0 24 24" width="14" height="14" aria-hidden="true">
      <path
        fill="currentColor"
        d="M19.43 12.98c.04-.32.07-.64.07-.98s-.03-.66-.07-.98l2.11-1.65a.5.5 0 0 0 .12-.64l-2-3.46a.5.5 0 0 0-.61-.22l-2.49 1a7.3 7.3 0 0 0-1.69-.98l-.38-2.65A.49.49 0 0 0 14 1h-4a.49.49 0 0 0-.49.42l-.38 2.65c-.61.25-1.17.58-1.69.98l-2.49-1a.5.5 0 0 0-.61.22l-2 3.46a.5.5 0 0 0 .12.64l2.11 1.65c-.04.32-.07.65-.07.98s.03.66.07.98l-2.11 1.65a.5.5 0 0 0-.12.64l2 3.46c.14.24.43.33.61.22l2.49-1c.52.4 1.08.73 1.69.98l.38 2.65c.04.24.25.42.49.42h4c.24 0 .45-.18.49-.42l.38-2.65c.61-.25 1.17-.59 1.69-.98l2.49 1c.18.07.47-.02.61-.22l2-3.46a.5.5 0 0 0-.12-.64l-2.11-1.65ZM12 15.5A3.5 3.5 0 1 1 12 8.5a3.5 3.5 0 0 1 0 7Z"
      />
    </svg>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { AUTH_TYPE, getAuthorization, setAuthorization } from '@/utils/request'
import {
  resolveHomeProjectUuid,
} from '@/config/homeDashboard.js'

export default {
  name: 'IsmBackToAdminButton',
  props: {
    projectUuid: { type: String, default: '' },
  },
  computed: {
    ...mapGetters('account', ['roles']),
    isAdmin() {
      const roles = this.roles || []
      return roles.some(r => r && r.id === 'Admin')
    },
  },
  methods: {
    resolveProjectUuid() {
      if (this.projectUuid) return this.projectUuid
      const stored = getAuthorization(AUTH_TYPE.AUTH1)
      if (stored) return stored
      const uid = this.$route.params.uid
      const resolved = resolveHomeProjectUuid(uid, this.$store)
      if (resolved) return resolved
      return ''
    },
    goAdmin() {
      const projectUuid = this.resolveProjectUuid()
      if (projectUuid) {
        setAuthorization({ token: projectUuid }, AUTH_TYPE.AUTH1)
      }
      // 回到「我的项目」列表，便于多项目切换；再点具体项目进入其管理后台
      const target = '/Project'
      if (this.$route.path.toLowerCase() !== target.toLowerCase()) {
        this.$router.push(target)
      }
    },
  },
}
</script>

<style lang="less" scoped>
/* 齿轮贴近右边缘，避免与“在线”状态文字重叠。
   垂直居中对齐同一状态行，小尺寸、低调样式，不撑高顶栏。
   注：齿轮是前端浮层（需 Admin 权限 + 路由），生成脚本只负责让"在线/日期"右移腾位。 */
.back-to-admin-btn {
  position: fixed;
  top: 21px;
  right: 5px;
  z-index: 1200;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  color: rgba(174, 233, 255, 0.5);
  cursor: pointer;
  user-select: none;
  background: transparent;
  opacity: 0.6;
  transition: color 0.2s ease, opacity 0.2s ease, filter 0.2s ease;
}

.back-to-admin-btn:hover {
  color: #d6f5ff;
  opacity: 1;
  filter: drop-shadow(0 0 4px rgba(34, 211, 238, 0.55));
}

.back-to-admin-btn:active {
  transform: translateY(1px);
}

.bta-gear {
  display: block;
  transition: transform 0.4s ease;
}

.back-to-admin-btn:hover .bta-gear {
  transform: rotate(60deg);
}
</style>
