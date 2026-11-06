/** dev 代理 / 直连后端不可达时的统一提示 */
export const BACKEND_UNREACHABLE_MESSAGE = '后端服务未启动或不可达(8081)'

/**
 * 判断 axios 错误是否由后端未启动、网络中断或 dev-proxy ECONNREFUSED 引起
 * @param {*} error axios 拒绝对象
 * @returns {boolean}
 */
export function isBackendUnreachableError(error) {
  if (!error) {
    return false
  }
  const message = String(error.message || '')
  const code = String(error.code || '')
  if (
    message === 'Network Error' ||
    /ECONNREFUSED|ETIMEDOUT|ECONNRESET|EHOSTUNREACH|ERR_NETWORK/i.test(message + code)
  ) {
    return true
  }
  let body = ''
  const data = error.response && error.response.data
  if (typeof data === 'string') {
    body = data
  } else if (data != null) {
    try {
      body = JSON.stringify(data)
    } catch (e) {
      body = String(data)
    }
  }
  const combined = message + body
  if (/Proxy error|Could not proxy|ECONNREFUSED/i.test(combined)) {
    return true
  }
  const status = error.response && error.response.status
  if (status === 502 || status === 503) {
    return true
  }
  if (status === 500 && /proxy|ECONNREFUSED/i.test(combined)) {
    return true
  }
  return false
}

export function backendUnreachableMessage() {
  return BACKEND_UNREACHABLE_MESSAGE
}
