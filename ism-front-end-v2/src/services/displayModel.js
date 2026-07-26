import {
  DISPLAYMODELADD,DISPLAYMODELSINGLE,DISPLAYMODELEDIT, DISPLAYMODELLIST,DISPLAYMODELDELETE,
  DISPLAYMODELDELETEDLIST,DISPLAYMODELRESTORE,DISPLAYMODELFORCEDEL,
  GETDISPLAYMODELLAYERDATA,SAVEDISPLAYMODELLAYERDATA,DISPLAYMODELPAGEADD,DISPLAYMODELPAGEDEL,DISPLAYMODELPAGEEDIT,GETDISPLAYMODELPAGERLAYERDATA,
  DISPLAYMODELPAGESETHOME,GETUSERDISPLAYLIST,DISPLAYMODELPAGECOPY,GETDISPLAYMODELLAYERDATABYTOKEN,MODELADDUSER,MODELDELUSER,GETMODELUSERS,
  DISPLAYMODELPAGEBINDTEMPLATE,DISPLAYMODELTEMPLATEMAP
} from '@/services/api'
import {request, METHOD} from '@/utils/request'

/**
 * 模型添加
 */
export async function displayModelAdd(params) {
  return request(DISPLAYMODELADD, METHOD.POST,params)
}

/**
 * 单个模型获取
 */
export async function getDisplayModelDetail(params) {
  return request(DISPLAYMODELSINGLE, METHOD.POST,params)
}

/**
 * 模型修改
 */
export async function displayModelEdit(params) {
  return request(DISPLAYMODELEDIT, METHOD.POST,params)
}

/**
 * 模型列表
 */
export async function displayModelList(params) {
  return request(DISPLAYMODELLIST, METHOD.POST,params)
}

/**
 * 模型删除（软删除，可在回收站恢复）
 */
export async function displayModelDelete(params) {
  return request(DISPLAYMODELDELETE, METHOD.POST,params)
}

/**
 * 回收站：已删除模型列表
 */
export async function displayModelDeletedList(params) {
  return request(DISPLAYMODELDELETEDLIST, METHOD.POST,params)
}

/**
 * 回收站：恢复已删除模型
 */
export async function displayModelRestore(params) {
  return request(DISPLAYMODELRESTORE, METHOD.POST,params)
}

/**
 * 回收站：彻底删除（物理删除，不可恢复）
 */
export async function displayModelForceDelete(params) {
  return request(DISPLAYMODELFORCEDEL, METHOD.POST,params)
}

/**
 * 模型图层数据
 */
const DISPLAY_META_TIMEOUT = 60000
const DISPLAY_PAGE_TIMEOUT = 30000
// 编辑态显式全量读取保留兼容窗口；运行态必须传 metaOnly=true。
const DISPLAY_FULL_TIMEOUT = 300000
const DISPLAY_LAYER_RETRY_DELAYS = [1000, 3000]

function isRetryableDisplayLayerError(error) {
  if (error && error.response) return false
  if (error && (error.code === 'ERR_CANCELED' || error.code === 'ECONNABORTED')) return false
  return typeof navigator === 'undefined' || navigator.onLine
}

async function requestDisplayLayer(url, params, timeout) {
  for (let attempt = 0; ; attempt += 1) {
    try {
      return await request(url, METHOD.POST, params, {
        timeout
      })
    } catch (error) {
      if (attempt >= DISPLAY_LAYER_RETRY_DELAYS.length ||
        !isRetryableDisplayLayerError(error)) {
        throw error
      }
      await new Promise(resolve => setTimeout(resolve, DISPLAY_LAYER_RETRY_DELAYS[attempt]))
    }
  }
}

export async function getDisplayModelPagerLayerData(params) {
  return requestDisplayLayer(GETDISPLAYMODELPAGERLAYERDATA, params, DISPLAY_PAGE_TIMEOUT)
}
/**
 * 模型图层数据
 */
export async function getDisplayModelLayerData(params) {
  const timeout = params && params.metaOnly === true
    ? DISPLAY_META_TIMEOUT
    : DISPLAY_FULL_TIMEOUT
  return requestDisplayLayer(GETDISPLAYMODELLAYERDATA, params, timeout)
}
/**
 * 模型图层数据
 */
export async function getLayerDataStructByToken(params) {
  const timeout = params && params.metaOnly === false
    ? DISPLAY_FULL_TIMEOUT
    : DISPLAY_META_TIMEOUT
  return requestDisplayLayer(GETDISPLAYMODELLAYERDATABYTOKEN, params, timeout)
}
/**
 * 模型图层数据
 */
export async function setDisplayModelLayerData(params) {
  return request(SAVEDISPLAYMODELLAYERDATA, METHOD.POST,params,{
    timeout:10000
  })
}

/**
 * 模型图层页面添加
 */
export async function DisplayModelPageAdd(params) {
  return request(DISPLAYMODELPAGEADD, METHOD.POST,params)
}
/**
 * 模型图层页面删除
 */
export async function DisplayModelPageDel(params) {
  return request(DISPLAYMODELPAGEDEL, METHOD.POST,params)
}
/**
 * 模型图层页面编辑
 */
export async function DisplayModelPageEdit(params) {
  return request(DISPLAYMODELPAGEEDIT, METHOD.POST,params)
}

/**
 * 模型图层页面首页
 */
export async function DisplayModelPageSetHome(params) {
  return request(DISPLAYMODELPAGESETHOME, METHOD.POST,params)
}
/**
 * 复制页面
 */
export async function DisplayModelPageCopy(params) {
  return request(DISPLAYMODELPAGECOPY, METHOD.POST,params)
}

/**
 * 绑定/解绑层级模板角色
 */
export async function DisplayModelPageBindTemplate(params) {
  return request(DISPLAYMODELPAGEBINDTEMPLATE, METHOD.POST, params)
}

/**
 * 查询大屏模板页映射
 */
export async function displayModelTemplateMap(params) {
  return request(DISPLAYMODELTEMPLATEMAP, METHOD.POST, params)
}

//获取用户模型
export async function GetUserDisplayList(params) {
  return request(GETUSERDISPLAYLIST, METHOD.POST,params)
}
//获取用户模型
export async function DisplayModelAddUser(params) {
  return request(MODELADDUSER, METHOD.POST,params)
}
//获取用户模型
export async function DisplayModelDelUser(params) {
  return request(MODELDELUSER, METHOD.POST,params)
}

//获取用户模型
export async function GetDisplayModelUser(params) {
  return request(GETMODELUSERS, METHOD.POST,params)
}

export default {
  displayModelAdd,
  getDisplayModelDetail,
  displayModelList,
  getDisplayModelLayerData,
  getDisplayModelPagerLayerData,
  setDisplayModelLayerData,
  displayModelDelete,
  displayModelDeletedList,
  displayModelRestore,
  displayModelForceDelete,
  DisplayModelPageAdd,
  DisplayModelPageDel,
  DisplayModelPageEdit,
  DisplayModelPageSetHome,
  DisplayModelPageCopy,
  DisplayModelPageBindTemplate,
  displayModelTemplateMap,
  getLayerDataStructByToken,
  DisplayModelAddUser,
  DisplayModelDelUser,
  GetDisplayModelUser
}
