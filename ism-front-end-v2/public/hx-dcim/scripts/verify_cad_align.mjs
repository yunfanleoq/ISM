#!/usr/bin/env node
/** 离线验证四配电室进线。5A1 仍走已验证竖刀布局；3A1/3A2/4B1 改为各室 CAD 网表。 */
import path from 'path';
import { fileURLToPath } from 'url';
import { bootHx, prepRoom, segments, touches, cluster, linked } from './lib_5a1.mjs';

const root = path.join(path.dirname(fileURLToPath(import.meta.url)), '..');
const ctx = bootHx(root);
const sandbox = ctx.sandbox;
const rooms = ['配电室3A1', '配电室3A2', '配电室4B1', '配电室5A1'];
const report = [];

function checkCadRoom(name) {
  const code = name.replace('配电室', '');
  const lay = ctx.rooms[code];
  const pg = prepRoom(ctx, name);
  const taps = sandbox.window.HXRooms.knifeTaps;
  const segs = segments(pg.nodes);
  const ends = [].concat(lay.meta.incomers || [], lay.meta.caps || []);
  const missed = [];
  ends.forEach((cad) => {
    const kn = pg.nodes.find((n) => n.cad === cad);
    if (!kn || kn.w !== 60 || kn.h !== 98) {
      missed.push(cad + ':size');
      return;
    }
    const t = taps(kn);
    if (!segs.some((s) => touches(s, t.in.x, t.in.y, 4))) missed.push(cad + ':top');
  });
  const buses = (lay.wires || []).filter((w) => w.net === 'bus').length;
  return {
    name,
    incomers: lay.meta.incomers,
    caps: lay.meta.caps,
    buses,
    missed,
    ok: ends.length >= 4 && missed.length === 0 && buses >= 2
  };
}

function check5A1() {
  const pg = prepRoom(ctx, '配电室5A1');
  const taps = sandbox.window.HX5A1.knifeTaps;
  const kn = pg.nodes.find((n) => n.cad === 'T1-401');
  const t = taps(kn);
  const segs = segments(pg.nodes);
  const uf = cluster(segs, 2.5);
  const bus = segs.find((s) => s.net === 'bus-I');
  const bus2 = segs.find((s) => s.net === 'bus-II');
  const busY = (bus.y1 + bus.y2) / 2;
  const gap = Math.min(bus2.x1, bus2.x2) - Math.max(bus.x1, bus.x2);
  const feed = linked(segs, uf, t.in.x, t.in.y, bus.cx, busY, 2.5);
  const bottomOff = !linked(segs, uf, t.out.x, t.out.y, bus.cx, busY, 2.5);
  return {
    name: '配电室5A1',
    busY,
    tap401x: t.in.x,
    busGap: gap,
    feed401: feed,
    bottom401OffBus: bottomOff,
    knife401: { w: kn.w, h: kn.h },
    ok: kn.w === 60 && kn.h === 98 && gap > 24 && feed && bottomOff && t.in.y > busY + 8
  };
}

for (const name of rooms) {
  report.push(name === '配电室5A1' ? check5A1() : checkCadRoom(name));
}

console.log(JSON.stringify(report, null, 2));
process.exit(report.every((r) => r.ok) ? 0 : 1);
