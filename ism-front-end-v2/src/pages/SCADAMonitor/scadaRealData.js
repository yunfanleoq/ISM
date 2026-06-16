import { DATA_KEY_MAP } from './config'

/** 实时数据缓存键：设备 UUID + 数据模型点 UUID */
export function realDataKey(deviceUid, modelDataUuid) {
  return `${deviceUid}:${modelDataUuid}`
}

/** 从 GetDeviceModelDataList 结果构建 per-device 点位目录 */
export function buildPointCatalog(buildings, rawList) {
  const catalog = {}
  const modelDataUuids = new Set()
  const deviceUuids = new Set()

  for (const modelItem of rawList || []) {
    const dataPoints = modelItem?.DataList || []
    if (!dataPoints.length) continue
    const modelUuid = dataPoints[0].muid || dataPoints[0].Muid || ''
    if (!modelUuid) continue

    for (const bldg of buildings) {
      for (const flr of bldg.floors || []) {
        for (const dev of flr.devices || []) {
          if (dev.muid !== modelUuid) continue
          for (const dp of dataPoints) {
            const modelDataUuid = dp.uuid || dp.mduid || dp.Muid
            if (!modelDataUuid || !dev.uid) continue
            catalog[realDataKey(dev.uid, modelDataUuid)] = {
              deviceUid: dev.uid,
              dataName: dp.name || dp.Name || '',
              modelDataUuid,
            }
            modelDataUuids.add(modelDataUuid)
            deviceUuids.add(dev.uid)
          }
        }
      }
    }
  }

  return {
    catalog,
    modelDataUuids: [...modelDataUuids],
    deviceUuids: [...deviceUuids],
  }
}

export function findModelDataUuid(catalog, deviceUid, logicalKey) {
  const candidates = DATA_KEY_MAP[logicalKey]
  if (!candidates || !deviceUid) return null
  for (const entry of Object.values(catalog)) {
    if (entry.deviceUid !== deviceUid) continue
    if (candidates.some(name => entry.dataName === name || entry.dataName.includes(name))) {
      return entry.modelDataUuid
    }
  }
  return null
}

/** 将 getRealDataByUuid / WebSocket 原始条目写入 pending 缓冲 */
export function ingestRealDataRows(pending, rows) {
  if (!Array.isArray(rows)) return
  for (const row of rows) {
    if (!row || row.value === '' || row.value === null || row.value === undefined) continue
    const deviceUid = row.duid || row.device_uuid || row.DeviceUuid
    const modelDataUuid = row.mduid || row.model_data_uuid || row.ModelDataUuid
    if (!deviceUid || !modelDataUuid) continue
    pending[realDataKey(deviceUid, modelDataUuid)] = {
      Value: row.value,
      Uuid: row.uuid,
      ModelDataUuid: modelDataUuid,
      DeviceUuid: deviceUid,
    }
  }
}

/** 将 WebSocket readDataPush 包写入 pending 缓冲 */
export function ingestWsPush(pending, wsData) {
  const deviceUid = wsData?.DeviceUuid
  if (!deviceUid || !Array.isArray(wsData.Data)) return
  for (const item of wsData.Data) {
    const modelDataUuid = item.ModelDataUuid || item.mduid || item.Uuid || item.uuid
    if (!modelDataUuid) continue
    if (item.Value === '' || item.Value === null || item.Value === undefined) continue
    pending[realDataKey(deviceUid, modelDataUuid)] = {
      ...item,
      DeviceUuid: deviceUid,
      ModelDataUuid: modelDataUuid,
    }
  }
}

/** 按 ISMRender 相同格式分组 getRealDataByUuid 响应并 ingest */
export function ingestRealDataResponse(pending, realData) {
  ingestRealDataRows(pending, realData)
}

export function readCachedValue(realTimeData, deviceUid, logicalKey, catalog) {
  const modelDataUuid = findModelDataUuid(catalog, deviceUid, logicalKey)
  if (!modelDataUuid) return null
  const data = realTimeData[realDataKey(deviceUid, modelDataUuid)]
  if (!data || data.Value === undefined || data.Value === null || data.Value === '') return null
  const v = Number(data.Value)
  if (Number.isNaN(v)) return null
  return logicalKey === 'pf' ? v : v
}
