/* 报警播报：优先用穹盾项目里单宇辰的声纹。合成失败时才退回普通中文语音，不再把音调拉到 2。 */

let current = null
let ticket = 0

function stripAlarmPreamble(text) {
  return String(text || '')
    .replace(/你有一条新的告警提醒/g, '')
    .replace(/你有一條新的告警提醒/g, '')
    .replace(/有一条新的报警/g, '')
    .replace(/有一条新的告警/g, '')
    .replace(/有一條新的報警/g, '')
    .replace(/有一條新的告警/g, '')
    .replace(/^[，,\s]*注意[，,]?/g, '')
}

const VALUE_LABELS = ['设备名称', '設備名稱', '数据名称', '數據名稱', '告警等级', '告警等級', '发生时间', '發生時間', '消息']

function speakClock(raw) {
  const m = String(raw || '').match(/(\d{4})-(\d{1,2})-(\d{1,2})(\d{1,2}):(\d{2})/)
  if (!m) return raw
  return Number(m[2]) + '月' + Number(m[3]) + '日' + Number(m[4]) + '点' + Number(m[5]) + '分'
}

/* 不念开头提醒，不念字段名。每个字段的值说一遍。 */
function shortenAlarm(text) {
  const raw = stripAlarmPreamble(text).replace(/\s+/g, '')
  const marks = []
  VALUE_LABELS.forEach(function (label) {
    const i = raw.indexOf(label)
    if (i >= 0) marks.push({ i: i, label: label })
  })
  marks.sort(function (a, b) { return a.i - b.i })
  const values = []
  marks.forEach(function (mark, n) {
    const start = mark.i + mark.label.length
    const end = n + 1 < marks.length ? marks[n + 1].i : raw.length
    let val = raw.slice(start, end).replace(/[，,。]/g, '')
    if (mark.label.indexOf('时间') >= 0 || mark.label.indexOf('時間') >= 0) val = speakClock(val)
    if (val) values.push(val)
  })
  let core = values.join('，')
  if (!core) {
    core = raw.replace(/设备名称|設備名稱|数据名称|數據名稱|告警等级|告警等級|发生时间|發生時間|消息/g, '')
  }
  return String(core || '').replace(/^[，,\s]+/, '').replace(/[，,\s]+$/, '')
}

function fallbackSpeak(text, rate) {
  if (!window.speechSynthesis) return
  const utter = new SpeechSynthesisUtterance(text)
  utter.lang = 'zh-CN'
  utter.pitch = 1
  utter.rate = 0.72
  const voices = window.speechSynthesis.getVoices() || []
  const zh = voices.find(function (v) {
    return /zh|chinese/i.test((v.lang || '') + (v.name || ''))
  })
  if (zh) utter.voice = zh
  window.speechSynthesis.cancel()
  window.speechSynthesis.speak(utter)
}

export function playAlarmVoice(text, rate) {
  const content = shortenAlarm(text)
  if (!content) return
  const mine = ++ticket
  if (current) {
    current.pause()
    current = null
  }
  if (window.speechSynthesis) window.speechSynthesis.cancel()
  const speed = rate || 1
  const ctrl = new AbortController()
  const timer = setTimeout(function () { ctrl.abort() }, 90000)
  fetch('/api/alarmVoice', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ text: content }),
    signal: ctrl.signal
  }).then(function (res) {
    if (!res.ok) throw new Error('voice')
    return res.blob()
  }).then(function (blob) {
    if (mine !== ticket) return
    const url = URL.createObjectURL(blob)
    const audio = new Audio(url)
    const pace = Math.min(1, Math.max(0.75, (Number(speed) || 1) * 0.9))
    audio.playbackRate = pace
    current = audio
    audio.onended = function () {
      URL.revokeObjectURL(url)
      if (current === audio) current = null
    }
    audio.play().catch(function () { fallbackSpeak(content, speed) })
  }).catch(function () {
    if (mine === ticket) fallbackSpeak(content, speed)
  }).finally(function () {
    clearTimeout(timer)
  })
}

if (typeof window !== 'undefined') window.playAlarmVoice = playAlarmVoice
