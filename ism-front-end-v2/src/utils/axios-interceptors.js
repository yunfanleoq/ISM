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
    if (response.code === 401) {
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
    const {response} = error
    if (response.status === 401) {
      message.error('无此权限')
    }
    return Promise.reject(error)
  }
}

const resp403 = {
  onFulfilled(response, options) {
    const {message} = options
    if (response.code === 403) {
      message.error('请求被拒绝')
    }
    return response
  },
  onRejected(error, options) {
    const {message} = options
    const {response} = error
    if (response.status === 403) {
      message.error('请求被拒绝')
    }
    return Promise.reject(error)
  }
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
    const {url, xsrfCookieName} = config

    if (getAuthorization(AUTH_TYPE.AUTH1)) {
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
    if(((url.indexOf('ImportProject')!==-1)||(url.indexOf('ExportProject')!==-1)||(url.indexOf('WitePhysicalID')!==-1)||(url.indexOf('GetPhysicalIDCheck')!==-1)||url.indexOf('setData')!==-1)||(url.indexOf('GetSystemParams')!==-1)|| (url.indexOf('GetSystemDeviceInfo')!==-1)||(url.indexOf('getDisplayModelLayerData')!==-1)||(url.indexOf('getDisplayModelLayerDataByToken')!==-1)||(url.indexOf('GetCustomPel')!==-1)||(url.indexOf('getRealDataByUuid')!==-1)||(url.indexOf('GetSystemMonitorList')!==-1))
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
      if (getAuthorization(AUTH_TYPE.AUTH1)) {
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
