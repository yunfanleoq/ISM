import { ALARM_COL, ALARM_ROW, rackWorld } from './layout'

export const FLY_MS = 1100

export const NAV_ITEMS = [
  { id: 'cockpit', label: '驾驶舱' },
  { id: 'building', label: '动环概要' },
  { id: 'room', label: '环境' },
  { id: 'power', label: '电力' },
  { id: 'battery', label: '电池' },
  { id: 'asset', label: '资产' },
]

const SCREEN_TITLE = {
  cockpit: '驾驶舱',
  building: '动环概要',
  floor: '四层概要',
  room: '环境总览',
  follow: '设备跟随',
  power: '电力监控',
  battery: '电池监控',
  asset: '资产',
}

export function createSession() {
  return {
    state: 'cockpit',
    busy: false,
    floorId: 'f4',
    roomId: 'r4-1',
    rackId: 'rk-c2',
    flyMs: FLY_MS,
  }
}

export function navKey(state) {
  if (state === 'floor') return 'building'
  if (state === 'follow') return 'room'
  return state
}

export function screenTitle(state) {
  return SCREEN_TITLE[state] || '驾驶舱'
}

export function backLabel(state) {
  const map = {
    building: '返回驾驶舱',
    floor: '返回楼宇',
    room: '返回楼层',
    follow: '退出跟随',
    power: '返回楼层',
    battery: '返回楼层',
    asset: '返回机房',
  }
  return map[state] || ''
}

export function poseFor(state) {
  const alarm = rackWorld(ALARM_COL, ALARM_ROW)
  const poses = {
    cockpit: { pos: [16, 11, 20], target: [0, 3.2, 0] },
    building: { pos: [10.5, 7.6, 13.5], target: [0, 4, 0] },
    floor: { pos: [0.4, 16.5, 11.5], target: [0, 0, 0] },
    room: { pos: [6.8, 5.8, 7.6], target: [0, 0.9, 0] },
    follow: {
      pos: [alarm.x + 2.3, 2.55, alarm.z + 2.5],
      target: [alarm.x, alarm.y, alarm.z],
    },
    power: { pos: [0.2, 13, 7.5], target: [0, 0, 0] },
    battery: { pos: [0.2, 12, 8], target: [0, 0, 0] },
    asset: { pos: [0.15, 1.85, 4.2], target: [0, 1.15, 0] },
  }
  return poses[state] || poses.cockpit
}

function nextOf(session, intent, payload) {
  const p = payload || {}
  if (intent === 'nav') {
    if (p.id === 'cockpit') return { state: 'cockpit' }
    if (p.id === 'building') return { state: 'building' }
    if (p.id === 'room') return { state: 'room', floorId: 'f4', roomId: 'r4-1' }
    if (p.id === 'power') return { state: 'power', floorId: 'f4' }
    if (p.id === 'battery') return { state: 'battery', floorId: 'f4' }
    if (p.id === 'asset') return { state: 'asset' }
    return null
  }
  if (intent === 'pick-building') {
    if (session.state !== 'cockpit') return null
    return { state: 'building' }
  }
  if (intent === 'pick-floor') {
    if (session.state !== 'building' && session.state !== 'floor') return null
    return { state: 'floor', floorId: p.floorId || 'f4' }
  }
  if (intent === 'pick-room') {
    if (session.state !== 'floor') return null
    return { state: 'room', roomId: p.roomId, floorId: session.floorId }
  }
  if (intent === 'pick-power-room') {
    if (session.state !== 'floor') return null
    return { state: 'power', floorId: session.floorId }
  }
  if (intent === 'pick-battery-room') {
    if (session.state !== 'floor') return null
    return { state: 'battery', floorId: session.floorId }
  }
  if (intent === 'pick-rack') {
    if (session.state !== 'room') return null
    return { state: 'follow', rackId: p.rackId || 'rk-c2' }
  }
  if (intent === 'back') {
    const back = {
      building: 'cockpit',
      floor: 'building',
      room: 'floor',
      follow: 'room',
      power: 'floor',
      battery: 'floor',
      asset: 'room',
    }
    const state = back[session.state]
    if (!state) return null
    return { state }
  }
  if (intent === 'select-asset') return { state: 'asset' }
  return null
}

export function begin(session, intent, payload) {
  if (session.busy) return { ok: false, session, reason: 'busy' }
  const patch = nextOf(session, intent, payload)
  if (!patch) return { ok: false, session, reason: 'illegal' }
  const same = patch.state === session.state
    && (patch.floorId || session.floorId) === session.floorId
    && intent !== 'select-asset'
  if (same) return { ok: false, session, reason: 'same' }
  return {
    ok: true,
    session: {
      ...session,
      ...patch,
      busy: true,
      flyMs: FLY_MS,
    },
  }
}

export function land(session) {
  return { ...session, busy: false }
}
