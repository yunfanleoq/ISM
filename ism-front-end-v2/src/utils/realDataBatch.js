/**
 * 大屏/数据仓库实时点位分批请求工具。
 * 单批默认 30、硬上限 100，避免一次拉上万点撑爆浏览器内存。
 */

export const REAL_DATA_DEFAULT_PAGE_SIZE = 30
export const REAL_DATA_MAX_PAGE_SIZE = 100
export const REAL_DATA_PAGE_SIZE_OPTIONS = ['20', '30', '50', '100']

export function clampPageSize(size, fallback = REAL_DATA_DEFAULT_PAGE_SIZE) {
  const n = parseInt(size, 10)
  if (!Number.isFinite(n) || n < 1) {
    return fallback
  }
  return Math.min(n, REAL_DATA_MAX_PAGE_SIZE)
}

/**
 * 将数组按 size 切块
 */
export function chunkArray(list, size = REAL_DATA_DEFAULT_PAGE_SIZE) {
  const arr = Array.isArray(list) ? list : []
  const chunkSize = clampPageSize(size)
  const out = []
  for (let i = 0; i < arr.length; i += chunkSize) {
    out.push(arr.slice(i, i + chunkSize))
  }
  return out
}

/**
 * 分批调用 getRealDataByUuid，合并 realData。
 * requestFn(params) 应返回 axios Promise。
 */
export async function fetchRealDataByUuidBatched(requestFn, { uuid = [], devices = [], batchSize = REAL_DATA_DEFAULT_PAGE_SIZE } = {}) {
  const uuids = (uuid || []).filter(u => typeof u === 'string' && u.length > 0)
  const deviceList = (devices || []).filter(d => typeof d === 'string' && d.length > 0)
  if (!uuids.length || !deviceList.length || typeof requestFn !== 'function') {
    return { code: 0, realData: [] }
  }

  const size = clampPageSize(batchSize)
  const uuidChunks = chunkArray(uuids, size)
  const merged = []
  let lastCode = 0

  for (let i = 0; i < uuidChunks.length; i++) {
    const res = await requestFn({
      uuid: uuidChunks[i],
      devices: deviceList,
    })
    const body = res && res.data
    if (!body) {
      continue
    }
    lastCode = body.code
    if (body.code === 0 && Array.isArray(body.realData)) {
      for (let k = 0; k < body.realData.length; k++) {
        merged.push(body.realData[k])
      }
    } else if (body.code !== 0) {
      return { code: body.code, realData: merged }
    }
  }

  return { code: lastCode, realData: merged }
}
