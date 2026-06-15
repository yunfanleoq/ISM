import {
  DISPLAYMODELADD,DISPLAYMODELSINGLE,DISPLAYMODELEDIT, DISPLAYMODELLIST,DISPLAYMODELDELETE,
  GETDISPLAYMODELLAYERDATA,SAVEDISPLAYMODELLAYERDATA,DISPLAYMODELPAGEADD,DISPLAYMODELPAGEDEL,DISPLAYMODELPAGEEDIT,GETDISPLAYMODELPAGERLAYERDATA,
  DISPLAYMODELPAGESETHOME,GETUSERDISPLAYLIST,DISPLAYMODELPAGECOPY,GETDISPLAYMODELLAYERDATABYTOKEN,MODELADDUSER,MODELDELUSER,GETMODELUSERS
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
 * 模型删除
 */
export async function displayModelDelete(params) {
  return request(DISPLAYMODELDELETE, METHOD.POST,params)
}

/**
 * 模型图层数据
 */
export async function getDisplayModelPagerLayerData(params) {
  return request(GETDISPLAYMODELPAGERLAYERDATA, METHOD.POST,params,{
    timeout:60000
  })
}
/**
 * 模型图层数据
 */
export async function getDisplayModelLayerData(params) {
  return request(GETDISPLAYMODELLAYERDATA, METHOD.POST,params,{
    timeout:60000
  })
}
/**
 * 模型图层数据
 */
export async function getLayerDataStructByToken(params) {
  return request(GETDISPLAYMODELLAYERDATABYTOKEN, METHOD.POST,params,{
    timeout:60000
  })
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
  DisplayModelPageAdd,
  DisplayModelPageDel,
  DisplayModelPageEdit,
  DisplayModelPageSetHome,
  DisplayModelPageCopy,
  getLayerDataStructByToken,
  DisplayModelAddUser,
  DisplayModelDelUser,
  GetDisplayModelUser
}
