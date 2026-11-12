import Vue from 'vue'
import Cookie from 'js-cookie'
import {GETSYSTEMDICTALL} from '@/services/api'
import {METHOD, request} from '@/utils/request'

const FALLBACK = {
  alarm_level: [
    {value: '0', label: '提示', color: '#0099FF', status: 1, sort: 0},
    {value: '1', label: '次要', color: '#0099FF', status: 1, sort: 1},
    {value: '2', label: '重要', color: '#FFFF00', status: 1, sort: 2},
    {value: '3', label: '紧急', color: '#FFA500', status: 1, sort: 3},
    {value: '4', label: '致命', color: '#FF0000', status: 1, sort: 4},
  ],
}

export const dictState = Vue.observable({
  revision: 0,
  byCode: JSON.parse(JSON.stringify(FALLBACK)),
  typeStatus: {alarm_level: 1},
})

let loaded = false
let pending = null

function hasToken() {
  return !!(sessionStorage.getItem('Authorization') || Cookie.get('Authorization'))
}

export function applySystemDict(list) {
  const byCode = {}
  const typeStatus = {}
  ;(list || []).forEach(row => {
    const code = row.DictCode
    if (!code) return
    typeStatus[code] = row.Status === 0 ? 0 : 1
    byCode[code] = (row.Items || []).map(item => ({
      value: String(item.ItemValue),
      label: item.ItemLabel,
      color: item.ItemColor || '',
      status: item.Status === 0 ? 0 : 1,
      sort: Number(item.Sort) || 0,
    })).sort((a, b) => a.sort - b.sort)
  })
  if (!byCode.alarm_level || !byCode.alarm_level.length) {
    byCode.alarm_level = FALLBACK.alarm_level
    typeStatus.alarm_level = 1
  }
  dictState.byCode = byCode
  dictState.typeStatus = typeStatus
  dictState.revision += 1
}

export function loadSystemDict(force) {
  if (!hasToken()) return Promise.resolve()
  if (loaded && !force) return Promise.resolve()
  if (pending && !force) return pending
  pending = request(GETSYSTEMDICTALL, METHOD.POST, {}).then(res => {
    const body = res && res.data
    if (body && body.code === 200 && Array.isArray(body.list)) {
      applySystemDict(body.list)
      loaded = true
    }
  }).catch(() => {}).finally(() => {
    pending = null
  })
  return pending
}

function itemsOf(code) {
  const _rev = dictState.revision
  void _rev
  return dictState.byCode[code] || []
}

export function dictOptions(code) {
  if (dictState.typeStatus[code] === 0) return []
  return itemsOf(code).filter(item => item.status !== 0)
}

export function dictItem(code, value) {
  const key = value === undefined || value === null ? '' : String(value)
  return itemsOf(code).find(item => item.value === key) || null
}

export function dictLabel(code, value) {
  const item = dictItem(code, value)
  if (item && item.label) return item.label
  if (value === undefined || value === null || value === '') return ''
  return String(value)
}

export function dictColor(code, value) {
  const item = dictItem(code, value)
  return (item && item.color) || ''
}
