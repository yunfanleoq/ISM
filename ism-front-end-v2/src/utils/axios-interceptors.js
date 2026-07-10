import Cookie from 'js-cookie'
// 401拦截

import {AUTH_TYPE, removeAuthorization, getAuthorization} from "@/utils/request";

const resp401 = {
  /**
   * 响应数据之前做点什么
   * @param response 响应对象
   * @param options 应用配置 包含: {router, i18n, store, message}
   * @returns {*}
   */
  onFulfilled(response, options) {
    const {message} = options
    if (!response) {
      return response
    }
    const code = response.data?.code ?? response.code
    if (code === 401) {
      message.error('无此权限')
    }
    return response
  },
  /**
   * 响应出错时执行
   * @param error 错误对象
   * @param options 应用配置 包含: {router, i18n, store, message}
   * @returns {Promise<never>}
   */
  onRejected(error, options) {
    const {message} = options
    if (!error) {
      return Promise.reject(error)
    }
    const status = error.response?.status
    if (status === 401) {
      message.error('无此权限')
    }
    return Promise.reject(error)
  }
}

const resp403 = {
  onFulfilled(response, options) {
    const {message} = options
    if (!response) {
      return response
    }
    const code = response.data?.code ?? response.code
    if (code === 403) {
      message.error('请求被拒绝')
    }
    return response
  },
  onRejected(error, options) {
    const {message} = options
    if (!error) {
      return Promise.reject(error)
    }
    const status = error.response?.status
    if (status === 403) {
      message.error('请求被拒绝')
    }
    return Promise.reject(error)
  }
}

/** 调用方在 config.headers 里显式传入的 ProjectUuid（大屏/AppRun 会按路由指定项目） */
function getExplicitProjectUuid(config, headerName) {
  const h = config.headers || {}
  return h[headerName] || (h.common && h.common[headerName]) || null
}

const reqCommon = {
  /**
   * 发送请求之前做些什么
   * @param config axios config
   * @param options 应用配置 包含: {router, i18n, store, message}
   * @returns {*}
   */
  onFulfilled(config, options) {
    const {message} = options
    const projectHeaderName = 'ProjectUuid'
    const ShareAppHeaderName = 'ShareAppToken'
    const authHeaderName = 'Authorization'
    const {url, xsrfCookieName} = config
    if (!config.headers) config.headers = {}
    if (!config.headers.common) config.headers.common = {}

    // axios xsrf 只读 Cookie；Cookie 失效时从 sessionStorage 显式注入 Authorization
    const bearerToken = getAuthorization(AUTH_TYPE.BEARER)
    if (bearerToken && url.indexOf('login') === -1) {
      config.headers.common[authHeaderName] = bearerToken
    }

    const explicitProjectUuid = getExplicitProjectUuid(config, projectHeaderName)

    if (!explicitProjectUuid && getAuthorization(AUTH_TYPE.AUTH1)) {
      config.headers.common[projectHeaderName] =  getAuthorization(AUTH_TYPE.AUTH1)
    }
    if (getAuthorization(AUTH_TYPE.AUTH3)) {
      config.headers.common[ShareAppHeaderName] =  getAuthorization(AUTH_TYPE.AUTH3)
    }

    if(url.indexOf('static/company/license.lic')!==-1)
    {
      return config
    }
    if(url.indexOf('GetAuthLicenseInfo')!==-1)
    {
      return config
    }
    if(url.indexOf('DisplayLoginPage')!==-1)
    {
      return config
    }
    if(((url.indexOf('ImportProject')!==-1)||(url.indexOf('ExportProject')!==-1)||(url.indexOf('WitePhysicalID')!==-1)||(url.indexOf('GetPhysicalIDCheck')!==-1)||url.indexOf('setData')!==-1)||(url.indexOf('GetSystemParams')!==-1)||(url.indexOf('GetSystemHomeDashboard')!==-1)||(url.indexOf('SetSystemHomeDashboard')!==-1)|| (url.indexOf('GetSystemDeviceInfo')!==-1)||(url.indexOf('getDisplayModelLayerData')!==-1)||(url.indexOf('getDisplayModelPagerLayerData')!==-1)||(url.indexOf('getDisplayModelLayerDataByToken')!==-1)||(url.indexOf('GetCustomPel')!==-1)||(url.indexOf('getRealDataByUuid')!==-1)||(url.indexOf('GetSystemMonitorList')!==-1)||(url.indexOf('monitortree')!==-1))
    {
      return config
    }
    if (url.indexOf('login') === -1 && xsrfCookieName && !getAuthorization(AUTH_TYPE.BEARER)) {
      message.warning('认证 token 已过期，请重新登录')
    }
    if((url.indexOf('login') !== -1)||(url.indexOf('login') !== -1)||(url.indexOf('ProjectList') !== -1)||((url.indexOf('ProjectAdd') !== -1))||(url.indexOf('ProjectEdit') !== -1)||((url.indexOf('ProjectDel') !== -1)))
    {
      delete config.headers.common[projectHeaderName]
      removeAuthorization(AUTH_TYPE.AUTH1)
    }
    else{
      if (explicitProjectUuid) {
        config.headers.common[projectHeaderName] = explicitProjectUuid
      } else if (getAuthorization(AUTH_TYPE.AUTH1)) {
        config.headers.common[projectHeaderName] =  getAuthorization(AUTH_TYPE.AUTH1)
      }
      else{
        message.warning('项目ID错误')
      }
    }
    return config
  },
  /**
   * 请求出错时做点什么
   * @param error 错误对象
   * @param options 应用配置 包含: {router, i18n, store, message}
   * @returns {Promise<never>}
   */
  onRejected(error, options) {
    const {message} = options
    message.error(error.message)
    return Promise.reject(error)
  }
}

export default {
  request: [reqCommon], // 请求拦截
  response: [resp401, resp403] // 响应拦截
}
