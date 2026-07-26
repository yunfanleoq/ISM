import {
  GETENERGYOVERVIEWCONFIG,
  SAVEENERGYOVERVIEWCONFIG,
  GETENERGYOVERVIEWCANDIDATES,
  GETENERGYOVERVIEWSTATS,
} from '@/services/api'
import { request, METHOD } from '@/utils/request'

const CACHE_TTL = 30000
const REFRESH_INTERVAL = 60000

let statsCache = null
let statsCacheTime = 0
let statsPromise = null
let refreshTimer = null
const subscribers = new Set()

function responseResult(response) {
  return response && response.data ? response.data.result : null
}

function normalizeSeries(series) {
  if (!Array.isArray(series)) return []
  return series.map(item => {
    if (Array.isArray(item)) {
      return { time: item[0], value: item[1] }
    }
    return {
      time: item && (item.time || item.timestamp || item.date),
      value: item && item.value,
    }
  }).filter(item => item.time != null && item.value != null)
}

function powerSeriesKey(key, name) {
  const value = `${key || ''} ${name || ''}`.toLowerCase()
  if (/reactivepower|总无功功率|無功/.test(value)) return 'reactivePower'
  if (/apparentpower|总视在功率|視在/.test(value)) return 'apparentPower'
  if (/activepower|总有功功率|有功/.test(value)) return 'activePower'
  return ''
}

function normalizePowerSeries(powerSeries) {
  const normalized = {
    activePower: [],
    reactivePower: [],
    apparentPower: [],
  }
  if (!powerSeries) return normalized
  if (!Array.isArray(powerSeries)) {
    Object.keys(normalized).forEach(key => {
      normalized[key] = normalizeSeries(powerSeries[key])
    })
    return normalized
  }
  const grouped = powerSeries.some(item =>
    item && !Array.isArray(item) && (Array.isArray(item.data) || Array.isArray(item.series))
  )
  if (!grouped) {
    normalized.activePower = normalizeSeries(powerSeries)
    return normalized
  }
  powerSeries.forEach(item => {
    const key = powerSeriesKey(item && item.key, item && item.name)
    if (key) normalized[key] = normalizeSeries(item.data || item.series)
  })
  return normalized
}

function normalizeStats(result) {
  const value = result || {}
  const buckets = Array.isArray(value.series) ? value.series : []
  const bucketPowerSeries = {
    activePower: normalizeSeries(buckets.map(item => ({ time: item.time, value: item.activePower }))),
    reactivePower: normalizeSeries(buckets.map(item => ({ time: item.time, value: item.reactivePower }))),
    apparentPower: normalizeSeries(buckets.map(item => ({ time: item.time, value: item.apparentPower }))),
  }
  const current = { ...(value.current || {}) }
  if (value.todayEnergy != null) current.todayEnergy = value.todayEnergy
  return {
    configured: value.configured !== false,
    dataStatus: value.dataStatus || '',
    missingPoints: value.missingPoints || [],
    totalDevices: Number(value.totalDevices || 0),
    eligibleDevices: Number(value.eligibleDevices || 0),
    validDevices: Number(value.validDevices || 0),
    missingDevices: Number(value.missingDevices || 0),
    ambiguousDevices: Number(value.ambiguousDevices || 0),
    resetDevices: Number(value.resetDevices || 0),
    current,
    powerSeries: value.powerSeries ? normalizePowerSeries(value.powerSeries) : bucketPowerSeries,
    energySeries: value.energySeries
      ? normalizeSeries(value.energySeries)
      : normalizeSeries(buckets.map(item => ({ time: item.time, value: item.energy }))),
    units: value.units || {},
  }
}

export function getEnergyOverviewConfig() {
  return request(GETENERGYOVERVIEWCONFIG, METHOD.POST)
}

export function saveEnergyOverviewConfig(params) {
  return request(SAVEENERGYOVERVIEWCONFIG, METHOD.POST, params)
}

export function getEnergyOverviewCandidates() {
  return request(GETENERGYOVERVIEWCANDIDATES, METHOD.POST)
}

export function getEnergyOverviewStats() {
  return request(GETENERGYOVERVIEWSTATS, METHOD.POST)
}

export function invalidateEnergyOverviewStats() {
  statsCache = null
  statsCacheTime = 0
}

export function fetchEnergyOverviewStats(force = false) {
  const now = Date.now()
  if (!force && statsCache && now - statsCacheTime < CACHE_TTL) {
    return Promise.resolve(statsCache)
  }
  if (statsPromise) return statsPromise

  statsPromise = getEnergyOverviewStats()
    .then(response => {
      statsCache = normalizeStats(responseResult(response))
      statsCacheTime = Date.now()
      return statsCache
    })
    .finally(() => {
      statsPromise = null
    })
  return statsPromise
}

export function refreshEnergyOverviewStats() {
  const pending = statsPromise
  invalidateEnergyOverviewStats()
  const refresh = pending
    ? pending.catch(() => null).then(() => {
      invalidateEnergyOverviewStats()
      return fetchEnergyOverviewStats(true)
    })
    : fetchEnergyOverviewStats(true)
  return refresh
    .then(stats => {
      subscribers.forEach(callback => callback(stats))
      return stats
    })
    .catch(error => {
      subscribers.forEach(callback => callback(null, error))
      throw error
    })
}

function ensureRefreshTimer() {
  if (refreshTimer || !subscribers.size) return
  refreshTimer = setInterval(() => {
    refreshEnergyOverviewStats().catch(() => {})
  }, REFRESH_INTERVAL)
}

export function subscribeEnergyOverviewStats(callback, immediate = true) {
  subscribers.add(callback)
  ensureRefreshTimer()
  if (immediate) {
    fetchEnergyOverviewStats().then(
      stats => {
        if (subscribers.has(callback)) callback(stats)
      },
      error => {
        if (subscribers.has(callback)) callback(null, error)
      }
    )
  }
  return () => {
    subscribers.delete(callback)
    if (!subscribers.size && refreshTimer) {
      clearInterval(refreshTimer)
      refreshTimer = null
    }
  }
}

