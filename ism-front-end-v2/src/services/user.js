import {SNMPSET,LOGIN, ROUTES,GETUSERINFO,SETUSERINFO,SETUSERPASSWORD,
  SYSTEMUSERLIST,SYSTEMUSERDEL,SYSTEMUSERADD,GETTOKENLIST,DELTOKENLIST,CREATETOKENLIST,USERUNLOCK
} from '@/services/api'
import {request, METHOD, removeAuthorization} from '@/utils/request'

/**
 * 登录服务
 * @param name 账户名
 * @param password 账户密码
 * @returns {Promise<AxiosResponse<T>>}
 */
export async function login(name, password,appid) {
  return request(LOGIN, METHOD.POST, {
    Username: name,
    password: password,
    appid:appid?appid:""
  })
}
/**
 * 登录服务
 * @param name 账户名
 * @param password 账户密码
 * @returns {Promise<AxiosResponse<T>>}
 */
export async function setValue(value) {
  return request(SNMPSET, METHOD.POST, {
    value: value
  })
}

export async function getRoutesConfig() {
  return request(ROUTES, METHOD.GET)
}

export async function GetUserInfo() {
  return request(GETUSERINFO, METHOD.POST)
}

export async function SetUserInfo(params) {
  return request(SETUSERINFO, METHOD.POST,params)
}
export async function SetUserPassword(params) {
  return request(SETUSERPASSWORD, METHOD.POST,params)
}

/**
 * 添加用户
 */
export async function AddSystemUser(params) {
  return request(SYSTEMUSERADD, METHOD.POST,params)
}

/**
 * 删除用户
 */
export async function SystemUserDel(params) {
  return request(SYSTEMUSERDEL, METHOD.POST,params)
}
/**
 * 用户列表
 */
export async function SystemUserList() {
  return request(SYSTEMUSERLIST, METHOD.POST)
}

/**
 * 退出登录
 */
export function logout() {
  localStorage.removeItem(process.env.VUE_APP_ROUTES_KEY)
  localStorage.removeItem(process.env.VUE_APP_PERMISSIONS_KEY)
  localStorage.removeItem(process.env.VUE_APP_ROLES_KEY)
  localStorage.removeItem(process.env.VUE_APP_USER_KEY)
  localStorage.removeItem("phoneUser")
  localStorage.removeItem("phonePassword")
  localStorage.removeItem("autologin")
  localStorage.removeItem("User")
  localStorage.removeItem("Password")
  removeAuthorization()
}

export async function GetAccessTokenList(params) {
  return request(GETTOKENLIST, METHOD.POST,params)
}
export async function CreateAccessToken(params) {
  return request(CREATETOKENLIST, METHOD.POST,params)
}
export async function DelAccessToken(params) {
  return request(DELTOKENLIST, METHOD.POST,params)
}
export async function UserUnlockScreen(params) {
  return request(USERUNLOCK, METHOD.POST,params)
}
export default {
  login,
  logout,
  SetUserInfo,
  GetUserInfo,
  SetUserPassword,
  AddSystemUser,
  SystemUserDel,
  SystemUserList,
  getRoutesConfig,
  GetAccessTokenList,
  CreateAccessToken,
  DelAccessToken,
  UserUnlockScreen
}
