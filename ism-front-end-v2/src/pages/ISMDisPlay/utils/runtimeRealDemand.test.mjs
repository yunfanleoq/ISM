import test from 'node:test'
import assert from 'node:assert/strict'
import * as demand from './runtimeRealDemand.js'
import { groupDemandRequests } from '../../../utils/realDataBatch.js'

function cell(partial) {
  return Object.assign({
    id: 'c1',
    x: 0,
    y: 0,
    width: 80,
    height: 40,
    shape: 'ism-view-svg-text',
    data: { detail: { identifier: 'c1', active: [], style: { visible: 1 } } },
  }, partial)
}

test('当前画面外的对象不取实时值，鼠标放上去才取', () => {
  const onScreen = cell({
    id: 'on',
    x: 10,
    y: 10,
    data: { detail: { identifier: 'on', active: [{ condition: { dataID: 'p-on', deviceSN: 'dev-a' } }], style: { visible: 1 } } },
  })
  const offScreen = cell({
    id: 'off',
    x: 2000,
    y: 10,
    data: { detail: { identifier: 'off', active: [{ condition: { dataID: 'p-off', deviceSN: 'dev-b' } }], style: { visible: 1 } } },
  })
  const view = { x: 0, y: 0, w: 400, h: 300 }
  const visible = demand.selectRuntimeRealDemand([onScreen, offScreen], { viewport: view })
  assert.deepEqual(visible.pairs, [{ uuid: 'p-on', device: 'dev-a' }])

  const hovered = demand.selectRuntimeRealDemand([onScreen, offScreen], { viewport: view, hoverId: 'off' })
  assert.deepEqual(hovered.pairs, [
    { uuid: 'p-on', device: 'dev-a' },
    { uuid: 'p-off', device: 'dev-b' },
  ])
})

test('翻到另一页后，上一页的测点不在名单里', () => {
  const pageA = [cell({
    id: 'a',
    data: { detail: { identifier: 'a', active: [{ condition: { dataID: 'pa', deviceSN: 'da' } }], style: { visible: 1 } } },
  })]
  const pageB = [cell({
    id: 'b',
    data: { detail: { identifier: 'b', active: [{ condition: { dataID: 'pb', deviceSN: 'db' } }], style: { visible: 1 } } },
  })]
  const a = demand.selectRuntimeRealDemand(pageA, {})
  const b = demand.selectRuntimeRealDemand(pageB, {})
  assert.notEqual(demand.demandSignature(a), demand.demandSignature(b))
  assert.deepEqual(b.pairs.map(p => p.uuid), ['pb'])
})

test('隐藏对象不订，系统点没有设备时单独列出', () => {
  const hidden = cell({
    id: 'h',
    visible: 0,
    data: { detail: { identifier: 'h', animate: { condition: { dataID: 'ph', deviceSN: 'dh' } }, style: { visible: 0 } } },
  })
  const system = cell({
    id: 's',
    data: { detail: { identifier: 's', active: [{ condition: { dataID: 'sys-1', deviceSN: '' } }], style: { visible: 1 } } },
  })
  const idle = demand.selectRuntimeRealDemand([hidden, system], {})
  assert.deepEqual(idle.pairs, [])
  assert.deepEqual(idle.systemUuids, ['sys-1'])
  const hoverHidden = demand.selectRuntimeRealDemand([hidden], { hoverId: 'h' })
  assert.deepEqual(hoverHidden.pairs, [{ uuid: 'ph', device: 'dh' }])
})

test('按设备拆请求，不同设备的测点不会交叉', () => {
  const requests = groupDemandRequests({
    pairs: [
      { uuid: 'p1', device: 'd1' },
      { uuid: 'p2', device: 'd1' },
      { uuid: 'p3', device: 'd2' },
    ],
    systemUuids: ['sys'],
  }, 30)
  assert.deepEqual(requests, [
    { uuid: ['p1', 'p2'], devices: ['d1'] },
    { uuid: ['p3'], devices: ['d2'] },
    { uuid: ['sys'], devices: [] },
  ])
})
