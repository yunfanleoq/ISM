#!/usr/bin/env node
/** 5A1：401 母线竖线接下刀上口，411 竖刀小母线，刀 60×98。 */
import path from 'path';
import { fileURLToPath } from 'url';
import { bootHx, prepRoom, segments, cluster, linked } from './lib_5a1.mjs';

const root = path.join(path.dirname(fileURLToPath(import.meta.url)), '..');
const ctx = bootHx(root);
const pg = prepRoom(ctx, '配电室5A1');
const taps = ctx.sandbox.window.HX5A1.knifeTaps;
const segs = segments(pg.nodes);
const uf = cluster(segs, 2.5);
const k401 = pg.nodes.find((n) => n.cad === 'T1-401');
const k411 = pg.nodes.find((n) => n.cad === 'T1-411-1');
const t401 = taps(k401);
const t411 = taps(k411);
const bus = segs.find((s) => s.net === 'bus-I');
const busY = (bus.y1 + bus.y2) / 2;
const xs = [1, 2, 3, 4, 5].map((i) => pg.nodes.find((n) => n.cad === 'T1-411-' + i).x).sort((a, b) => a - b);
let minPitch = 999;
for (let i = 1; i < xs.length; i++) minPitch = Math.min(minPitch, xs[i] - xs[i - 1]);
const ys = [1, 2, 3, 4, 5].map((i) => pg.nodes.find((n) => n.cad === 'T1-411-' + i).y);

const report = {
  knife401: { w: k401.w, h: k401.h, orient: k401.orient || 'v' },
  knife411: { w: k411.w, h: k411.h, orient: k411.orient || 'v', y: k411.y },
  busY,
  sameRow411: Math.max(...ys) - Math.min(...ys) < 1,
  parallelMinPitch: minPitch,
  feed401: linked(segs, uf, t401.in.x, t401.in.y, bus.cx, busY, 2.5),
  bottom401OffBus: !linked(segs, uf, t401.out.x, t401.out.y, bus.cx, busY, 2.5),
  tap411: linked(segs, uf, t411.in.x, t411.in.y, bus.cx, busY, 2.5),
  ok:
    k401.w === 60 &&
    k401.h === 98 &&
    k411.w === 60 &&
    k411.h === 98 &&
    !k401.orient &&
    !k411.orient &&
    Math.max(...ys) - Math.min(...ys) < 1 &&
    minPitch >= 64 &&
    linked(segs, uf, t401.in.x, t401.in.y, bus.cx, busY, 2.5) &&
    t401.in.y > busY + 8 &&
    !linked(segs, uf, t401.out.x, t401.out.y, bus.cx, busY, 2.5) &&
    linked(segs, uf, t411.in.x, t411.in.y, bus.cx, busY, 2.5)
};
console.log(JSON.stringify(report, null, 2));
process.exit(report.ok ? 0 : 1);
