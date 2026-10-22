<template>
  <div class="boxer xunan-status">
    <div class="box1">
      <div class="middle1">
        <div class="middle_left">
          <div class="brand-panel">
            <img class="brand-logo" src="/static/branding/logo-xunan-hexagon.png" alt="数据中心动环监控平台">
            <h1 class="brand-title">数据中心动环监控平台</h1>
            <p class="brand-subtitle">POWER MONITORING &amp; INTELLIGENT MANAGEMENT</p>
            <ul class="brand-features">
              <li>实时数据采集与监控</li>
              <li>智能告警与故障预警</li>
              <li>能耗分析与报表统计</li>
            </ul>
          </div>
        </div>
        <div class="middle_right">
          <div class="mr_box">
            <div class="status-badge">{{ view.badge }}</div>
            <h2 class="status-title">{{ view.title }}</h2>
            <p class="status-lead">{{ view.lead }}</p>
            <p class="status-detail">{{ view.detail }}</p>
            <p v-if="extra" class="status-extra">{{ extra }}</p>
            <button class="status-btn" type="button" :disabled="checking" @click="retry">
              {{ checking ? '正在检查…' : '重新检查' }}
            </button>
          </div>
          <div class="login-footer">© 数据中心动环监控平台</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { GetSystemCodeCheck } from '@/services/system'

const REASONS = {
  backend: {
    badge: '服务状态',
    title: '后端没有连上',
    lead: '后端没有启动。',
    detail: '当前页面访问不到本机服务。请先启动本机后端（端口 8082），确认进程在监听后再重新检查。这不是授权页面。'
  },
  license: {
    badge: '授权状态',
    title: '授权信息读不到',
    lead: '后端已经连上，但授权文件没有生效。',
    detail: '请联系管理员检查授权文件。这里不再进入厂商授权页。'
  },
  expired: {
    badge: '授权状态',
    title: '授权已过期',
    lead: '后端已经连上，软件授权期限已到。',
    detail: '请联系管理员续期。'
  }
}

export default {
  name: 'Auth',
  data () {
    return {
      checking: false,
      extra: ''
    }
  },
  computed: {
    view () {
      const key = this.$route.query.reason
      return REASONS[key] || REASONS.backend
    }
  },
  watch: {
    '$route.query.reason' () {
      this.extra = ''
    }
  },
  methods: {
    retry () {
      const _t = this
      this.checking = true
      this.extra = ''
      GetSystemCodeCheck().then(function (res) {
        const code = Number(res && res.data ? res.data.code : NaN)
        if (code === 0) {
          _t.$router.replace('/login')
          return
        }
        const reason = code === -3 ? 'expired' : 'license'
        _t.$router.replace({ path: '/auth', query: { reason: reason } }).catch(function () {})
      }).catch(function () {
        _t.extra = '刚刚又检查了一次，后端仍然没有连上。'
        if (_t.$route.query.reason !== 'backend') {
          _t.$router.replace({ path: '/auth', query: { reason: 'backend' } }).catch(function () {})
        }
      }).finally(function () {
        _t.checking = false
      })
    }
  }
}
</script>

<style lang="less" scoped>
.boxer {
  background-image: linear-gradient(135deg, #0a1a3a 0%, #0d2b55 45%, #0a3d5c 100%);
  height: 100%;
  min-height: 100vh;
  position: relative;
  overflow: hidden;
}
.box1 {
  width: 1080px;
  max-width: calc(100% - 48px);
  margin: 6vh auto 0;
  overflow: hidden;
  box-shadow: 0 0 24px 8px rgba(0, 40, 80, 0.35);
  border-radius: 12px;
  background: rgba(8, 28, 58, 0.55);
  border: 1px solid rgba(19, 194, 194, 0.25);
}
.middle1 {
  display: flex;
  min-height: 560px;
}
.middle_left {
  flex: 0 0 52%;
  display: flex;
  align-items: center;
  background: radial-gradient(ellipse at 30% 20%, rgba(19, 194, 194, 0.18), transparent 55%),
    linear-gradient(160deg, rgba(10, 40, 80, 0.9), rgba(8, 24, 50, 0.95));
}
.brand-panel {
  padding: 40px 48px;
  color: #e8f7f7;
}
.brand-logo {
  width: 88px;
  height: 88px;
  object-fit: contain;
  margin-bottom: 20px;
}
.brand-title {
  margin: 0 0 10px;
  font-size: 32px;
  font-weight: 700;
  letter-spacing: 2px;
  color: #fff;
  line-height: 1.3;
}
.brand-subtitle {
  margin: 0 0 28px;
  font-size: 12px;
  letter-spacing: 1.5px;
  color: rgba(160, 220, 230, 0.85);
}
.brand-features {
  list-style: none;
  padding: 0;
  margin: 0;
}
.brand-features li {
  margin: 12px 0;
  padding: 12px 16px;
  border-radius: 8px;
  background: rgba(12, 40, 72, 0.65);
  border: 1px solid rgba(19, 194, 194, 0.35);
  color: #d7f3f3;
  font-size: 15px;
}
.middle_right {
  flex: 1;
  background: rgba(255, 255, 255, 0.96);
  position: relative;
  display: flex;
  flex-direction: column;
}
.mr_box {
  width: 100%;
  max-width: 420px;
  margin: auto;
  padding: 48px 36px 72px;
  box-sizing: border-box;
}
.status-badge {
  display: inline-block;
  padding: 4px 10px;
  border-radius: 999px;
  background: rgba(19, 194, 194, 0.12);
  color: #0e7490;
  font-size: 13px;
  letter-spacing: 1px;
}
.status-title {
  margin: 18px 0 12px;
  font-size: 32px;
  line-height: 1.25;
  color: #0f2744;
}
.status-lead {
  margin: 0 0 12px;
  font-size: 18px;
  line-height: 1.6;
  color: #1f3b57;
}
.status-detail,
.status-extra {
  margin: 0 0 12px;
  font-size: 15px;
  line-height: 1.7;
  color: #4b6478;
}
.status-extra {
  color: #b45309;
}
.status-btn {
  margin-top: 18px;
  min-width: 160px;
  height: 44px;
  padding: 0 22px;
  border: 0;
  border-radius: 6px;
  background: linear-gradient(90deg, #0e7490, #13c2c2);
  color: #fff;
  font-size: 16px;
  cursor: pointer;
}
.status-btn:disabled {
  opacity: 0.7;
  cursor: wait;
}
.login-footer {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 28px;
  text-align: center;
  color: #7d93a6;
  font-size: 13px;
}
</style>
