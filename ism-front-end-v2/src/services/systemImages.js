import {
  SYSTEMIMAGELIST, SYSTEMIMAGEDEL
} from '@/services/api'
import {request, METHOD} from '@/utils/request'


/**
 * snmp模型列表
 */
export async function systemImageList() {
  return request(SYSTEMIMAGELIST, METHOD.POST)
}

/**
 * snmp模型删除
 */
export async function systemImageDel(params) {
  return request(SYSTEMIMAGEDEL, METHOD.POST,params)
}

export default {
  systemImageDel,
  systemImageList
}
