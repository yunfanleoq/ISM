<template>
  <div class="xiaoan">
    <section v-if="open" class="panel">
      <header>
        <b>小安</b>
        <span>文字或语音，问当前页面上的问题</span>
        <button type="button" @click="open = false">关闭</button>
      </header>
      <div ref="log" class="log">
        <p v-for="(item, i) in messages" :key="i" :class="item.role">{{ item.text }}</p>
      </div>
      <form @submit.prevent="send">
        <input v-model="text" :placeholder="listening ? '正在听…' : '问小安'" />
        <button type="button" :class="{ on: listening }" @click="toggleVoice">语音</button>
        <button type="submit" :disabled="busy">发送</button>
      </form>
    </section>
    <button type="button" class="sprite" :class="{ talk: speaking, listen: listening }" title="小安" @click="toggle">
      <svg viewBox="0 0 88 108" aria-hidden="true">
        <defs>
          <linearGradient id="xiaoan-cap" x1="44" y1="18" x2="44" y2="52" gradientUnits="userSpaceOnUse">
            <stop offset="0" stop-color="#1d6ec4" />
            <stop offset="1" stop-color="#0b3164" />
          </linearGradient>
          <linearGradient id="xiaoan-gold" x1="44" y1="28" x2="44" y2="96" gradientUnits="userSpaceOnUse">
            <stop offset="0" stop-color="#ffe7a3" />
            <stop offset="1" stop-color="#c89212" />
          </linearGradient>
        </defs>
        <path fill="#0c3d73" d="M16 92c8-18 48-18 56 0 2 8-8 14-28 14S14 100 16 92z" />
        <circle cx="44" cy="64" r="24" fill="#f4f8ff" stroke="#0c3d73" stroke-width="2" />
        <circle cx="30" cy="70" r="3.2" fill="#f3b7b0" />
        <circle cx="58" cy="70" r="3.2" fill="#f3b7b0" />
        <ellipse cx="35" cy="62" rx="3.2" ry="3.6" fill="#102a43" />
        <ellipse cx="53" cy="62" rx="3.2" ry="3.6" fill="#102a43" />
        <circle cx="36.2" cy="60.8" r="1.05" fill="#fff" />
        <circle cx="54.2" cy="60.8" r="1.05" fill="#fff" />
        <path d="M36 72c2.6 3.2 9.4 3.2 12 0" fill="none" stroke="#0c3d73" stroke-width="1.8" stroke-linecap="round" />
        <path fill="url(#xiaoan-cap)" d="M22 48C22 26 66 26 66 48c-6-8-38-8-44 0z" />
        <path fill="#071a33" d="M14 49c10 6 50 6 60 0-7 3-53 3-60 0z" />
        <path d="M26 46h36" stroke="url(#xiaoan-gold)" stroke-width="3.2" stroke-linecap="round" />
        <path fill="url(#xiaoan-gold)" d="M44 28l7 2.4v7c0 4.8-7 8.6-7 8.6s-7-3.8-7-8.6v-7z" />
        <circle cx="44" cy="35.5" r="2.3" fill="#0c3d73" />
        <path fill="url(#xiaoan-gold)" d="M44 90l8 2.8v8.4L44 105l-8-3.8v-8.4z" />
        <path fill="#0c3d73" d="M44 94.2l2.2 1.6 2.5-.4-1.6 2 .6 2.4-2.2-1.2-2.2 1.2.6-2.4-1.6-2 2.5.4z" />
      </svg>
      <em>小安</em>
    </button>
  </div>
</template>

<script>
import { METHOD, request } from '@/utils/request'
import { resolveXunzhiJfProjectUuid } from '@/xunzhiJf/projectContext'

export default {
  name: 'XiaoAn',
  data() {
    return {
      open: false,
      text: '',
      busy: false,
      speaking: false,
      listening: false,
      projectUuid: '',
      recognition: null,
      messages: [],
    }
  },
  watch: {
    '$route.fullPath'() {
      if (this.messages.length <= 1) this.greet()
    },
  },
  mounted() {
    this.$nextTick(() => this.greet())
    resolveXunzhiJfProjectUuid().then(uuid => { this.projectUuid = uuid || '' }).catch(() => {})
  },
  beforeDestroy() {
    this.stopVoice()
    if (window.speechSynthesis) window.speechSynthesis.cancel()
  },
  methods: {
    page() {
      return (this.$route && (this.$route.fullPath || this.$route.path)) || ''
    },
    greet() {
      const page = this.page()
      let text = '我是小安。页面上有问题，可以直接用文字或语音问我。'
      if (page.indexOf('ontology') >= 0) text = '我是小安。这里是本体中心，点一个节点就能看到它的生产数据和日志。'
      if (page.indexOf('/agent') >= 0) text = '我是小安。这里按阶段和环节看智能体，技能和工具在右边。'
      this.messages = [{ role: 'an', text }]
    },
    toggle() {
      this.open = !this.open
      this.$nextTick(this.scroll)
    },
    scroll() {
      const el = this.$refs.log
      if (el) el.scrollTop = el.scrollHeight
    },
    async send() {
      const text = (this.text || '').trim()
      if (!text || this.busy) return
      this.messages.push({ role: 'me', text })
      this.text = ''
      this.busy = true
      try {
        const res = await request('/api/agent/talk', METHOD.POST, {
          project_uuid: this.projectUuid,
          page: this.page(),
          text,
        })
        const turn = (res.data && res.data.turn) || {}
        const reply = turn.reply || '我这次没有答上来。'
        this.messages.push({ role: 'an', text: reply })
        this.speak(turn.speak || reply)
      } catch (e) {
        this.messages.push({ role: 'an', text: '我现在连不上，你稍后再问。' })
      } finally {
        this.busy = false
        this.$nextTick(this.scroll)
      }
    },
    speak(text) {
      if (!window.speechSynthesis || !text) return
      window.speechSynthesis.cancel()
      const utter = new SpeechSynthesisUtterance(String(text).slice(0, 180))
      utter.lang = 'zh-CN'
      utter.rate = 0.9
      utter.pitch = 1
      const voices = window.speechSynthesis.getVoices() || []
      const zh = voices.find(v => /zh|chinese/i.test((v.lang || '') + (v.name || '')))
      if (zh) utter.voice = zh
      utter.onstart = () => { this.speaking = true }
      utter.onend = () => { this.speaking = false }
      utter.onerror = () => { this.speaking = false }
      window.speechSynthesis.speak(utter)
    },
    toggleVoice() {
      if (this.listening) {
        this.stopVoice()
        return
      }
      const Rec = window.SpeechRecognition || window.webkitSpeechRecognition
      if (!Rec) {
        this.messages.push({ role: 'an', text: '这个浏览器不能听写，直接打字就行。' })
        this.open = true
        return
      }
      const rec = new Rec()
      rec.lang = 'zh-CN'
      rec.interimResults = false
      rec.onresult = (event) => {
        const said = event.results && event.results[0] && event.results[0][0] ? event.results[0][0].transcript : ''
        this.text = said
        this.listening = false
        if (said) this.send()
      }
      rec.onerror = () => { this.listening = false }
      rec.onend = () => { this.listening = false }
      this.recognition = rec
      this.listening = true
      this.open = true
      rec.start()
    },
    stopVoice() {
      this.listening = false
      if (this.recognition) {
        try { this.recognition.stop() } catch (e) { /* 已经停了 */ }
        this.recognition = null
      }
    },
  },
}
</script>

<style scoped>
.xiaoan { position: fixed; right: 18px; bottom: 18px; z-index: 30; }
.sprite { width: 88px; height: 118px; border: 0; background: transparent; cursor: pointer; padding: 0; }
.sprite svg { width: 84px; height: 102px; display: block; margin: 0 auto; filter: drop-shadow(0 8px 14px rgba(11, 49, 100, 0.28)); }
.sprite em { display: block; margin-top: -2px; text-align: center; font-style: normal; font-size: 12px; color: #0c3d73; font-weight: 700; letter-spacing: 1px; }
.sprite.talk svg, .sprite.listen svg { animation: bob 1s ease-in-out infinite; }
.sprite.listen svg { filter: drop-shadow(0 0 8px rgba(232, 193, 90, 0.85)); }
.panel { position: absolute; right: 0; bottom: 124px; width: 320px; height: 420px; background: #fff; border: 1px solid #d5e2f0; border-radius: 12px; box-shadow: 0 12px 32px rgba(11, 49, 100, 0.16); display: flex; flex-direction: column; }
header { display: flex; align-items: center; gap: 8px; padding: 10px 12px; border-bottom: 1px solid #f0f4f8; }
header b { color: #102a43; }
header span { flex: 1; color: #829ab1; font-size: 12px; }
header button, form button { border: 0; background: #e6f6ff; color: #186faf; border-radius: 6px; padding: 4px 8px; cursor: pointer; }
.log { flex: 1; overflow: auto; padding: 10px 12px; }
.log p { margin: 0 0 8px; padding: 8px 10px; border-radius: 8px; font-size: 13px; white-space: pre-wrap; }
.log .an { background: #f0f7fb; color: #102a43; }
.log .me { background: #186faf; color: #fff; margin-left: 32px; }
form { display: flex; gap: 6px; padding: 10px; border-top: 1px solid #f0f4f8; }
input { flex: 1; border: 1px solid #d9e2ec; border-radius: 6px; padding: 6px 8px; }
form button.on { background: #186faf; color: #fff; }
@keyframes bob { 50% { transform: translateY(-4px); } }
</style>
