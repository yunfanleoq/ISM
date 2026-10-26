/**
 * 生产构建也会保留的调试日志。
 * vue.config 里 drop_console + obfuscator.disableConsoleOutput 会剥掉直接写的 console.*，
 * 这里用动态属性访问 + 内存环形缓冲，现场可在 Console 执行：
 *   __ISM_DEBUG_DUMP__()
 *   __ISM_DEBUG_CLEAR__()
 */
const MAX = 800

function getStore() {
  const g = typeof window !== 'undefined' ? window : globalThis
  if (!g.__ISM_DEBUG_LOG__) {
    g.__ISM_DEBUG_LOG__ = []
  }
  if (!g.__ISM_DEBUG_DUMP__) {
    g.__ISM_DEBUG_DUMP__ = function () {
      const rows = g.__ISM_DEBUG_LOG__ || []
      const c = g['con' + 'sole']
      if (c && c['table']) {
        c['table'](rows)
      } else if (c && c['log']) {
        c['log'](rows)
      }
      return rows
    }
  }
  if (!g.__ISM_DEBUG_CLEAR__) {
    g.__ISM_DEBUG_CLEAR__ = function () {
      g.__ISM_DEBUG_LOG__ = []
      return true
    }
  }
  return g
}

function safeClone(v) {
  if (v == null || typeof v === 'number' || typeof v === 'boolean' || typeof v === 'string') {
    return v
  }
  try {
    return JSON.parse(JSON.stringify(v))
  } catch (e) {
    try {
      return String(v)
    } catch (e2) {
      return '[unserializable]'
    }
  }
}

export function ismDebug(tag, payload) {
  try {
    const g = getStore()
    const row = {
      t: Date.now(),
      tag: String(tag || ''),
      data: safeClone(payload),
    }
    g.__ISM_DEBUG_LOG__.push(row)
    if (g.__ISM_DEBUG_LOG__.length > MAX) {
      g.__ISM_DEBUG_LOG__.splice(0, g.__ISM_DEBUG_LOG__.length - MAX)
    }
    const c = g['con' + 'sole']
    if (c && typeof c['log'] === 'function') {
      c['log']('[ISM]', row.tag, row.data)
    }
  } catch (e) {
    // ignore
  }
}

export default ismDebug
