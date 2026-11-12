import {
  DELSYSTEMDICTITEM,
  DELSYSTEMDICTTYPE,
  GETSYSTEMDICTALL,
  GETSYSTEMDICTITEMLIST,
  GETSYSTEMDICTTYPELIST,
  SAVESYSTEMDICTITEM,
  SAVESYSTEMDICTTYPE,
} from '@/services/api'
import {METHOD, request} from '@/utils/request'

export function getSystemDictAll() {
  return request(GETSYSTEMDICTALL, METHOD.POST, {})
}
export function getSystemDictTypeList() {
  return request(GETSYSTEMDICTTYPELIST, METHOD.POST, {})
}
export function saveSystemDictType(params) {
  return request(SAVESYSTEMDICTTYPE, METHOD.POST, params)
}
export function delSystemDictType(params) {
  return request(DELSYSTEMDICTTYPE, METHOD.POST, params)
}
export function getSystemDictItemList(params) {
  return request(GETSYSTEMDICTITEMLIST, METHOD.POST, params)
}
export function saveSystemDictItem(params) {
  return request(SAVESYSTEMDICTITEM, METHOD.POST, params)
}
export function delSystemDictItem(params) {
  return request(DELSYSTEMDICTITEM, METHOD.POST, params)
}
