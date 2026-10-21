function cells() {
  const list = []
  for (let i = 1; i <= 18; i++) {
    const watch = i === 7
    const id = '4-1A1-UPS1-BAT-' + String(i).padStart(2, '0')
    list.push({
      id,
      voltage: watch ? 582.4 : 597.0,
      current: 2.0,
      level: watch ? 'watch' : 'normal',
    })
  }
  return list
}

export function createBuilding3Fixture() {
  return {
    site: { id: 'b3', name: '3号楼' },
    pue: 1.15,
    wue: 0,
    outdoor: { temp: 2.6, humidity: 57 },
    energy: { kwh: 12684504, kw: 3749, itPct: 86.95, coolPct: 6.08, otherPct: 9.13 },
    feeds: [
      { name: '1#进线', kw: 1860, loadPct: 62 },
      { name: '2#进线', kw: 1889, loadPct: 64 },
    ],
    monthKwh: [86, 92, 88, 101, 97, 110, 104, 99, 108, 112, 96],
    floors: [
      { id: 'f1', label: '1层', overviewLabel: '一层总览' },
      { id: 'f2', label: '2层', overviewLabel: '二层总览' },
      { id: 'f3', label: '3层', overviewLabel: '三层总览' },
      { id: 'f4', label: '4层', overviewLabel: '四层总览' },
    ],
    rooms: [
      { id: 'r2-1', floorId: 'f2', kind: 'it', label: '机房2-1' },
      { id: 'r4-1', floorId: 'f4', kind: 'it', label: '机房4-1', itKw: 420, coolKw: 86 },
      { id: 'r4-2', floorId: 'f4', kind: 'it', label: '机房4-2', itKw: 388, coolKw: 74 },
      { id: 'r4-3', floorId: 'f4', kind: 'it', label: '机房4-3', itKw: 356, coolKw: 69 },
      { id: 'r4-4', floorId: 'f4', kind: 'it', label: '机房4-4', itKw: 341, coolKw: 63 },
      { id: 'p4-1', floorId: 'f4', kind: 'power', label: '配电室4-1' },
      { id: 'b4-1', floorId: 'f4', kind: 'battery', label: '电池间4-1' },
    ],
    racks: [
      { id: 'rk-a1', name: 'A1', roomId: 'r4-1', col: 0, row: 0, level: 'normal', temp: 23.1 },
      { id: 'rk-b1', name: 'B1', roomId: 'r4-1', col: 1, row: 0, level: 'normal', temp: 23.4 },
      { id: 'rk-c2', name: 'C2', roomId: 'r4-1', col: 2, row: 0, level: 'alarm', temp: 31.6 },
      { id: 'rk-d1', name: 'D1', roomId: 'r4-1', col: 3, row: 0, level: 'normal', temp: 22.8 },
      { id: 'rk-a2', name: 'A2', roomId: 'r4-1', col: 0, row: 1, level: 'normal', temp: 22.4 },
      { id: 'rk-b2', name: 'B2', roomId: 'r4-1', col: 1, row: 1, level: 'normal', temp: 23.0 },
      { id: 'rk-c3', name: 'C3', roomId: 'r4-1', col: 2, row: 1, level: 'normal', temp: 22.9 },
      { id: 'rk-d2', name: 'D2', roomId: 'r4-1', col: 3, row: 1, level: 'normal', temp: 23.2 },
      { id: 'A01', name: 'A01', roomId: 'r2-1', col: 0, row: 0, level: 'normal', temp: 23.0 },
    ],
    roomEnv: { maxTemp: 31.6, minTemp: 22.4, leak: '正常', acOn: 6, acTotal: 8 },
    power: {
      roomId: 'p4-1',
      voltage: 398.6,
      current: 86.4,
      kw: 52.8,
      loadPct: 41,
      bypassCurrent: 0,
    },
    ups: {
      id: 'ups1',
      mode: '市电供电',
      batteryVoltage: 570.0,
      batteryCurrent: 0.2,
      backupMin: 365,
      capacityPct: 100,
      outputLoadPct: 8.7,
    },
    batteryRoom: { id: 'b4-1', cells: cells() },
    assets: [
      {
        id: 'a-server',
        name: '计算节点-07',
        status: '闲置',
        place: '3号楼 / 2层 / 机房2-1 / A01',
        rackId: 'A01',
        uStart: null,
        uHeight: 1,
        shelveU: 2,
      },
      {
        id: 'a-switch',
        name: '汇聚交换机-02',
        status: '在用',
        place: '3号楼 / 2层 / 机房2-1 / A01',
        rackId: 'A01',
        uStart: 10,
        uHeight: 1,
      },
      {
        id: 'a-borrow',
        name: '测试仪-借出',
        status: '借用到期',
        place: '3号楼 / 2层 / 机房2-1 / A01',
        rackId: 'A01',
        uStart: 20,
        uHeight: 1,
      },
    ],
  }
}
