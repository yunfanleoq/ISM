<template>
  <!-- 锁定屏幕弹窗：科技感风格 -->
  <a-modal
      v-model="visible"
      :width="360"
      :closable="false"
      :mask-closable="false"
      :footer="null"
      title=""
      class="lock-screen-modal"
      :mask-style="{ backgroundColor: 'rgba(10, 15, 30, 0.95)' }"
  >
    <!-- 动态网格背景（底层） -->
    <div class="grid-bg"></div>

    <!-- 光晕效果（中层） -->
    <div class="glow-effect"></div>

    <!-- 粒子装饰（顶层） -->
    <div class="particles-container">
      <div v-for="i in 12" :key="i" class="particle" :style="{
        top: `${Math.random() * 100}%`,
        left: `${Math.random() * 100}%`,
        animationDelay: `${Math.random() * 5}s`,
        animationDuration: `${5 + Math.random() * 10}s`
      }"></div>
    </div>

    <!-- 科技感边框容器 -->
    <div class="tech-container">
      <!-- 顶部装饰条 -->
      <div class="tech-header">
        <div class="glow-line"></div>
      </div>

      <!-- 锁定图标（带呼吸动画） -->
      <div class="lock-icon">
        <a-icon
            type="lock"
            theme="filled"
            class="pulse-animation"
            style="font-size: 56px; color: #00e5ff;"
        />
      </div>

      <!-- 提示文本 -->
      <div class="lock-title">
        <span class="title-glow">屏幕已锁定</span>
        <p>请输入密码解锁系统</p>
      </div>

      <!-- 密码输入框（带焦点效果） -->
      <div class="password-container">
        <a-input-password
            v-model="inputPassword"
            placeholder="输入解锁密码"
            @press-enter="handleUnlock"
            class="tech-input"
            :class="{ 'tech-input-focused': focused }"
            @focus="focused = true"
            @blur="focused = false"
        />
      </div>

      <!-- 解锁按钮（带悬停动画） -->
      <a-button
          type="primary"
          block
          @click="handleUnlock"
          class="unlock-btn"
      >
        <span>解锁系统</span>
        <a-icon type="arrow-right" style="margin-left: 8px;" />
      </a-button>

      <!-- 错误提示（带滑入动画） -->
      <a-alert
          v-if="errorMsg"
          type="error"
          message="验证失败"
          description="密码错误，请重新输入"
          show-icon
          class="error-alert"
      />

      <!-- 底部状态指示 -->
      <div class="status-bar">
        <span class="status-indicator"></span>
        <span class="system-status">SYSTEM LOCKED - SECURE MODE</span>
      </div>
    </div>
  </a-modal>
</template>

<script>
import { mapState, mapActions } from 'vuex';
import {login, getRoutesConfig, UserUnlockScreen} from '@/services/user'
import md5 from "js-md5";
import {unLockScreen} from "@/store/ISM/actions";
export default {
  name: 'LockScreen',
  i18n: require('@/i18n/language'),
  data() {
    return {
      inputPassword: '',
      errorMsg: '',
      focused: false // 密码框焦点状态
    };
  },
  computed: {
    visible() {
      return this.$store.state.ISMDisPlayEditorTool.isLocked;
    }
  },
  methods: {
    ...mapActions('ISMDisPlayEditorTool', ['unLockScreen']),

    async handleUnlock() {
      if (!this.inputPassword) {
        this.$message.warning({
          content: '请输入密码',
          style: { backgroundColor: '#1e293b', color: '#e2e8f0' }
        });
        return;
      }

      try {
        let _t = this
        this.errorMsg = '';
        const params = {
          Password:md5(this.inputPassword)
        }
        UserUnlockScreen(params).then(function (res){
          if(res.data.code==200)
          {
            _t.$message.success({
              content: '解锁成功，系统正在恢复',
              style: { backgroundColor: '#0f172a', color: '#00e5ff' }
            });
            _t.unLockScreen()
          }
          else
          {
            _t.$message.error({
              content: '解锁失败',
              style: { backgroundColor: '#0f172a', color: '#00e5ff' }
            });
          }
        }).catch(function(){
          _t.logging = false
          _t.$message.error(_t.$t('loginPage.serverError'), 3)
        })

      } catch (err) {
        this.errorMsg = err.message;
        this.inputPassword = '';
        // 错误时密码框抖动效果
        const input = _t.$el.querySelector('.tech-input');
        input.classList.add('shake-animation');
        setTimeout(() => input.classList.remove('shake-animation'), 600);
      }
    },

  }
};
</script>

<style scoped>
/* 动态网格背景 */
.grid-bg {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image:
      linear-gradient(rgba(0, 229, 255, 0.05) 1px, transparent 1px),
      linear-gradient(90deg, rgba(0, 229, 255, 0.05) 1px, transparent 1px);
  background-size: 30px 30px;
  animation: gridMove 40s linear infinite;
  z-index: 0;
}

/* 网格移动动画 */
@keyframes gridMove {
  0% {
    background-position: 0 0;
  }
  100% {
    background-position: 300px 300px;
  }
}

/* 光晕效果 */
.glow-effect {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 600px;
  height: 600px;
  transform: translate(-50%, -50%);
  background: radial-gradient(circle, rgba(0, 229, 255, 0.15) 0%, transparent 70%);
  filter: blur(40px);
  z-index: 1;
  animation: glowPulse 15s ease-in-out infinite;
}

/* 光晕呼吸动画 */
@keyframes glowPulse {
  0%, 100% {
    transform: translate(-50%, -50%) scale(1);
    opacity: 0.8;
  }
  50% {
    transform: translate(-50%, -50%) scale(1.2);
    opacity: 0.5;
  }
}

/* 粒子装饰 */
.particles-container {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 1;
  overflow: hidden;
}

.particle {
  position: absolute;
  width: 2px;
  height: 2px;
  background: #00e5ff;
  border-radius: 50%;
  box-shadow: 0 0 6px rgba(0, 229, 255, 0.8);
  opacity: 0;
  animation: floatParticle linear infinite;
}

/* 粒子漂浮动画 */
@keyframes floatParticle {
  0% {
    transform: translateY(0) translateX(0);
    opacity: 0;
  }
  20% {
    opacity: 0.8;
  }
  80% {
    opacity: 0.8;
  }
  100% {
    transform: translateY(-100px) translateX(50px);
    opacity: 0;
  }
}

/* 解决样式穿透问题的工具类 */
::v-deep .tech-input {
  background: rgba(30, 41, 59, 0.7) !important;
  backdrop-filter: blur(8px);
  border-radius: 6px !important;
  color: #e2e8f0 !important;
  height: 48px !important;
  transition: all 0.3s ease !important;
  border: 1px solid #2d3a59 !important;
  position: relative;
  z-index: 10;
}

/* 焦点状态样式 */
::v-deep .tech-input-focused {
  border-color: #00e5ff !important;
  box-shadow: 0 0 8px rgba(0, 229, 255, 0.5) !important;
}

/* 输入框内部文本样式 */
::v-deep .tech-input .ant-input {
  background: transparent !important;
  color: #e2e8f0 !important;
  padding: 0 12px !important;
  height: 100% !important;
  font-size: 14px !important;
}

/* 密码框图标样式 */
::v-deep .tech-input .ant-input-suffix {
  color: #64748b !important;
}

/* 占位符样式 */
::v-deep .tech-input input::placeholder {
  color: #64748b !important;
  opacity: 1 !important;
}

/* 主容器样式 */
.lock-screen-modal .ant-modal-content {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 360px;
  background: rgba(18, 26, 46, 0.6) !important;
  backdrop-filter: blur(12px);
  border: 1px solid rgba(0, 229, 255, 0.3);
  border-radius: 12px;
  box-shadow: 0 0 30px rgba(0, 229, 255, 0.15);
  overflow: hidden;
  z-index: 10;
}

/* 科技感容器 */
.tech-container {
  padding: 30px 25px;
  position: relative;
  z-index: 10;
}

/* 顶部装饰条 */
.tech-header {
  display: flex;
  justify-content: center;
  margin-bottom: 30px;
}

.glow-line {
  width: 60px;
  height: 3px;
  background: #00e5ff;
  box-shadow: 0 0 10px rgba(0, 229, 255, 0.8);
  border-radius: 3px;
}

/* 锁定图标 */
.lock-icon {
  text-align: center;
  margin-bottom: 30px;
}

/* 呼吸动画 */
.pulse-animation {
  animation: pulse 2s infinite ease-in-out;
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
    filter: drop-shadow(0 0 8px rgba(0, 229, 255, 0.6));
  }
  50% {
    transform: scale(1.05);
    filter: drop-shadow(0 0 12px rgba(0, 229, 255, 0.8));
  }
}

/* 标题样式 */
.lock-title {
  text-align: center;
  margin-bottom: 30px;
}

.lock-title .title-glow {
  font-size: 18px;
  font-weight: 600;
  color: #00e5ff;
  text-shadow: 0 0 8px rgba(0, 229, 255, 0.6);
}

.lock-title p {
  margin: 8px 0 0;
  color: #94a3b8;
  font-size: 14px;
}

/* 密码输入框容器 */
.password-container {
  margin-bottom: 20px;
}

/* 错误抖动动画 */
::v-deep .shake-animation {
  animation: shake 0.6s cubic-bezier(.36,.07,.19,.97) both;
}

@keyframes shake {
  10%, 90% { transform: translateX(-1px); }
  20%, 80% { transform: translateX(2px); }
  30%, 50%, 70% { transform: translateX(-3px); }
  40%, 60% { transform: translateX(3px); }
}

/* 解锁按钮 */
.unlock-btn {
  background: linear-gradient(90deg, #00b4d8, #00e5ff) !important;
  border: none !important;
  height: 48px !important;
  border-radius: 6px !important;
  font-size: 16px !important;
  font-weight: 500 !important;
  color: #0f172a !important;
  transition: all 0.3s ease !important;
  overflow: hidden;
  position: relative;
  z-index: 10;
}

.unlock-btn::after {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: linear-gradient(
      to right,
      rgba(255, 255, 255, 0) 0%,
      rgba(255, 255, 255, 0.3) 50%,
      rgba(255, 255, 255, 0) 100%
  );
  transform: rotate(30deg);
  animation: shine 3s infinite;
}

.unlock-btn:hover {
  box-shadow: 0 0 15px rgba(0, 229, 255, 0.6) !important;
  transform: translateY(-2px);
}

@keyframes shine {
  0% { transform: translateX(-100%) rotate(30deg); }
  100% { transform: translateX(100%) rotate(30deg); }
}

/* 错误提示 */
::v-deep .error-alert {
  margin-top: 16px !important;
  background-color: rgba(239, 68, 68, 0.1) !important;
  border: 1px solid rgba(239, 68, 68, 0.3) !important;
  color: #f87171 !important;
  animation: slideIn 0.3s ease forwards;
  position: relative;
  z-index: 10;
}

::v-deep .error-alert .ant-alert-message {
  color: #f87171 !important;
}

::v-deep .error-alert .ant-alert-description {
  color: #fecaca !important;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 底部状态条 */
.status-bar {
  display: flex;
  align-items: center;
  margin-top: 25px;
  padding-top: 15px;
  border-top: 1px solid rgba(148, 163, 184, 0.1);
  position: relative;
  z-index: 10;
}

.status-indicator {
  display: inline-block;
  width: 8px;
  height: 8px;
  background: #ef4444;
  border-radius: 50%;
  margin-right: 8px;
  box-shadow: 0 0 6px rgba(239, 68, 68, 0.8);
}

.system-status {
  font-size: 12px;
  color: #64748b;
  letter-spacing: 0.5px;
}
</style>