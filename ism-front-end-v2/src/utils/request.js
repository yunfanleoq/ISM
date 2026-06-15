import axios from 'axios'
import Cookie from 'js-cookie'

// 跨域认证信息 header 名
const xsrfHeaderName = 'Authorization'
const projectHeaderName = 'ProjectUuid'
const ShareAppHeaderName = 'ShareAppToken'

axios.defaults.timeout = 30000
axios.defaults.withCredentials= true
axios.defaults.xsrfHeaderName= xsrfHeaderName
axios.defaults.xsrfCookieName= xsrfHeaderName

// 认证类型
const AUTH_TYPE = {
  BEARER: 'Bearer',
  BASIC: 'basic',
  AUTH1: 'auth1',
  AUTH2: 'auth2',
  AUTH3: 'auth3',
}

// http method
const METHOD = {
  GET: 'get',
  POST: 'post'
}

/**
 * axios请求
 * @param url 请求地址
 * @param method {METHOD} http method
 * @param params 请求参数
 * @returns {Promise<AxiosResponse<T>>}
 */
async function request(url, method, params, config) {
  switch (method) {
    case METHOD.GET:
      return axios.get(url, {params, ...config})
    case METHOD.POST:
      return axios.post(url, params, config)
    default:
      return axios.get(url, {params, ...config})
  }
}

/**
 * 设置认证信息
 * @param auth {Object}
 * @param authType {AUTH_TYPE} 认证类型，默认：{AUTH_TYPE.BEARER}
 */
function setAuthorization(auth, authType = AUTH_TYPE.BEARER) {
  switch (authType) {
    case AUTH_TYPE.BEARER:
      Cookie.set(xsrfHeaderName, auth.token, {expires: auth.expireAt})
      sessionStorage.setItem(xsrfHeaderName,auth.token)
      break
    case AUTH_TYPE.BASIC:
      break
    case AUTH_TYPE.AUTH1: {
      Cookie.set(projectHeaderName, auth.token)
      sessionStorage.setItem(projectHeaderName,auth.token)
      break
    }
    case AUTH_TYPE.AUTH3: {
      Cookie.set(ShareAppHeaderName, auth.token)
      sessionStorage.setItem(ShareAppHeaderName,auth.token)
      break
    }
    case AUTH_TYPE.AUTH2:
    default:
      break
  }
}

/**
 * 移出认证信息
 * @param authType {AUTH_TYPE} 认证类型
 */
function removeAuthorization(authType = AUTH_TYPE.BEARER) {
  switch (authType) {
    case AUTH_TYPE.BEARER: {
      Cookie.remove(xsrfHeaderName)
      Cookie.remove(projectHeaderName)
      Cookie.remove(ShareAppHeaderName)
      sessionStorage.removeItem(xsrfHeaderName)
      sessionStorage.removeItem(projectHeaderName)
      sessionStorage.removeItem(ShareAppHeaderName)
      break
    }
    case AUTH_TYPE.BASIC:
    {
      break
    }
    case AUTH_TYPE.AUTH1: {
      Cookie.remove(projectHeaderName)
      sessionStorage.removeItem(projectHeaderName)
      break
    }
    case AUTH_TYPE.AUTH3: {
      Cookie.remove(ShareAppHeaderName)
      sessionStorage.removeItem(ShareAppHeaderName)
      break
    }
    case AUTH_TYPE.AUTH2:
    default:
      break
  }
}
/**
 * 检查认证信息
 * @param authType
 * @returns {boolean}
 */
function getAuthorization(authType = AUTH_TYPE.BEARER) {
  switch (authType) {
    case AUTH_TYPE.BEARER: {
      return Cookie.get(xsrfHeaderName)?Cookie.get(xsrfHeaderName):sessionStorage.getItem(xsrfHeaderName)
    }
    case AUTH_TYPE.BASIC:
    case AUTH_TYPE.AUTH1: {
      return Cookie.get(projectHeaderName)?Cookie.get(projectHeaderName):sessionStorage.getItem(projectHeaderName)
    }
    case AUTH_TYPE.AUTH2:
    {
      break
    }
    case AUTH_TYPE.AUTH3: {
      return Cookie.get(ShareAppHeaderName)?Cookie.get(ShareAppHeaderName):sessionStorage.getItem(ShareAppHeaderName)
    }
    default:
      break
  }
  return false
}

/**
 * 检查认证信息
 * @param authType
 * @returns {boolean}
 */
function checkAuthorization(authType = AUTH_TYPE.BEARER) {
  switch (authType) {
    case AUTH_TYPE.BEARER:
      if (Cookie.get(xsrfHeaderName)||sessionStorage.getItem(xsrfHeaderName)) {
        return true
      }
      break
    case AUTH_TYPE.BASIC:
      break;
    case AUTH_TYPE.AUTH1: {
      if (Cookie.get(projectHeaderName)||sessionStorage.getItem(projectHeaderName)) {
        return true
      }
      break
    }
    case AUTH_TYPE.AUTH2:
      break;
    case AUTH_TYPE.AUTH3: {
      if (Cookie.get(ShareAppHeaderName)||sessionStorage.getItem(ShareAppHeaderName)) {
        return true
      }
      break
    }
    default:
      break
  }
  return false
}

/**
 * 加载 axios 拦截器
 * @param interceptors
 * @param options
 */
function loadInterceptors(interceptors, options) {
  const {request, response} = interceptors
  // 加载请求拦截器
  request.forEach(item => {
    let {onFulfilled, onRejected} = item
    if (!onFulfilled || typeof onFulfilled !== 'function') {
      onFulfilled = config => config
    }
    if (!onRejected || typeof onRejected !== 'function') {
      onRejected = error => Promise.reject(error)
    }
    axios.interceptors.request.use(
      config => onFulfilled(config, options),
      error => onRejected(error, options)
    )
  })
  // 加载响应拦截器
  response.forEach(item => {
    let {onFulfilled, onRejected} = item
    if (!onFulfilled || typeof onFulfilled !== 'function') {
      onFulfilled = response => response
    }
    if (!onRejected || typeof onRejected !== 'function') {
      onRejected = error => Promise.reject(error)
    }
    axios.interceptors.response.use(
      response => onFulfilled(response, options),
      error => onRejected(error, options)
    )
  })
}

/**
 * 解析 url 中的参数
 * @param url
 * @returns {Object}
 */
function parseUrlParams(url) {
  const params = {}
  if (!url || url === '' || typeof url !== 'string') {
    return params
  }
  const paramsStr = url.split('?')[1]
  if (!paramsStr) {
    return params
  }
  const paramsArr = paramsStr.replace(/&|=/g, ' ').split(' ')
  for (let i = 0; i < paramsArr.length / 2; i++) {
    const value = paramsArr[i * 2 + 1]
    params[paramsArr[i * 2]] = value === 'true' ? true : (value === 'false' ? false : value)
  }
  return params
}

export {
  METHOD,
  AUTH_TYPE,
  request,
  setAuthorization,
  removeAuthorization,
  checkAuthorization,
  getAuthorization,
  loadInterceptors,
  parseUrlParams
}
