import {
  SNAPHISTORYIMAGESLIST,VIDEOCODECS,VIDERECVIVER,VIDEOLIST,VIDEODEL,VIDEOADD,GETVIDEOSTATUS,VIDEOEDIT,SETVIDEOSTOPORSTART,PTZCONTROL,GETMONIBUCA,HISTORYVIDEOLIST,HISTORYIMAGESLIST
} from '@/services/api'
import {request, METHOD} from '@/utils/request'


export async function videoCodec(uuid) {
  return request(VIDEOCODECS+uuid, METHOD.GET)
}

export async function videoRecviver(params) {
  return request(VIDERECVIVER, METHOD.POST,params)
}

export async function videoList(params) {
  return request(VIDEOLIST, METHOD.POST,params)
}
export async function historyListVideoList(params) {
  return request(HISTORYVIDEOLIST, METHOD.POST,params)
}
export async function snapHistoryListVideoList(params) {
  return request(SNAPHISTORYIMAGESLIST, METHOD.POST,params)
}
export async function historyImagesList(params) {
  return request(HISTORYIMAGESLIST, METHOD.POST,params)
}
export async function videoDel(params) {
  return request(VIDEODEL, METHOD.POST,params)
}

export async function videoAdd(params) {
  return request(VIDEOADD, METHOD.POST,params)
}
export async function videoEdit(params) {
  return request(VIDEOEDIT, METHOD.POST,params)
}
export async function getVideoStatus(params) {
  return request(GETVIDEOSTATUS, METHOD.POST,params)
}
export async function VideoStopOrStart(params) {
  return request(SETVIDEOSTOPORSTART, METHOD.POST,params)
}
export async function PtzControl(params) {
  return request(PTZCONTROL, METHOD.POST,params)
}
export async function GetMonibucaVideoList(params) {
  return request(GETMONIBUCA, METHOD.POST,params)
}

export default {
  videoCodec,
  videoRecviver,
  videoList,
  videoDel,
  videoAdd,
  videoEdit,
  getVideoStatus,
  PtzControl,
  VideoStopOrStart,
  GetMonibucaVideoList,
  historyListVideoList,
  historyImagesList,
  snapHistoryListVideoList
}
