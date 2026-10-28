#!/usr/bin/env node
/** 5A1 显式布局 vs CAD 网表。通过则 exit 0。 */
import path from 'path';
import { fileURLToPath } from 'url';
import { bootHx, prepRoom, segments, cluster, linked, idsTouching, touches } from './lib_5a1.mjs';

const root = path.join(path.dirname(fileURLToPath(import.meta.url)), '..');
const TOL = 2.5;
const ctx = bootHx(root);
const pg = prepRoom(ctx, '配电室5A1');
const segs = segments(pg.nodes);
const uf = cluster(segs, TOL);
const taps = ctx.sandbox.window.HX5A1.knifeTaps;
const fail = [];

function knife(cad) {
  return pg.nodes.find((n) => n.cad === cad && (n.t === 'I' || n.t === 'S'));
}
function must(cond, msg) {
  if (!cond) fail.push(msg);
}
function reach(a, b, msg) {
  must(linked(segs, uf, a.x, a.y, b.x, b.y, TOL), msg);
}
function apart(a, b, msg) {
  must(!linked(segs, uf, a.x, a.y, b.x, b.y, TOL), msg);
}
function withLink(ax, ay, bx, by) {
  const extra = {
    horiz: Math.abs(ay - by) < 1,
    x1: ax,
    y1: ay,
    x2: bx,
    y2: by,
    net: 'knife',
    c: '#fff'
  };
  const s2 = segs.concat([extra]);
  return cluster(s2, TOL);
}

const knives = pg.nodes.filter((n) => n.cad && (n.t === 'I' || n.t === 'S'));
must(knives.length === 97, 'knife count ' + knives.length);
must(!knives.some((n) => n.orient === 'h' || Number(n.rot) === 90), 'horizontal knife');
must(knives.every((n) => n.w === 60 && n.h === 98), 'knife size');
must(knives.every((n) => n.k && n.p && n.dev), 'scene k/p/dev');
must(!pg._lay5a1Missing || pg._lay5a1Missing.length === 0, 'missing scene ' + (pg._lay5a1Missing || []).join(','));
must(!knives.some((n) => /410|420/.test(n.cad)), 'capacitor cabinet drawn');

const y411 = [1, 2, 3, 4, 5].map((i) => knife('T1-411-' + i).y);
must(Math.max(...y411) - Math.min(...y411) < 1, '411 not side by side');
const x421 = [5, 4, 3, 2, 1].map((i) => knife('T2-421-' + i).x);
must(x421.every((x, i) => i === 0 || x > x421[i - 1]), '421 order');

const k401 = knife('T1-401');
const k402 = knife('T2-402');
const k445 = knife('445');
const t401 = taps(k401);
const t402 = taps(k402);
const t445 = taps(k445);
const busI = segs.find((s) => s.net === 'bus-I');
const busII = segs.find((s) => s.net === 'bus-II');
must(busI && busII, 'two top buses');
const busIy = (busI.y1 + busI.y2) / 2;
const gap = Math.min(busII.x1, busII.x2) - Math.max(busI.x1, busI.x2);
must(gap > 24, 'bus gap ' + gap);
must(t401.in.y > busIy + 8, '401 knife overlaps bus');
reach(t401.in, { x: t401.in.x, y: busIy }, '401 bus not down to knife top');
apart(t401.out, { x: busI.cx, y: busIy }, '401 bottom still jogged up to bus');
const busIIy = (busII.y1 + busII.y2) / 2;
reach(t402.in, { x: t402.in.x, y: busIIy }, '402 bus not down to knife top');
apart(t402.out, { x: busII.cx, y: busIIy }, '402 bottom still jogged up to bus');
must(!segs.some((s) => s.net === 'feed-401' || s.net === 'feed-402'), 'old feed jog remains');
function riserOff(net, busSeg, busY) {
  const r = segs.filter((s) => s.net === net);
  must(r.length >= 1, net + ' missing');
  r.forEach((seg) => {
    apart({ x: seg.x1, y: seg.y1 }, { x: busSeg.cx, y: busY }, net + ' hits bus');
    apart({ x: seg.x2, y: seg.y2 }, { x: busSeg.cx, y: busY }, net + ' hits bus');
  });
  if (!r.length) return;
  const top = Math.min.apply(null, r.map((s) => Math.min(s.y1, s.y2)));
  must(top > busY + 4, net + ' reaches bus y');
}
riserOff('xfmr-5A1-T1', busI, busIy);
riserOff('xfmr-5A1-T2', busII, busIIy);
must(segs.some((s) => s.net === 'xfmr-5A1-T1' && touches(s, t401.out.x, t401.out.y, TOL)), 'T1 not wired into 401');
must(segs.some((s) => s.net === 'xfmr-5A1-T2' && touches(s, t402.out.x, t402.out.y, TOL)), 'T2 not wired into 402');
['drop-T1-410', 'drop-T2-420'].forEach((net) => {
  const d = segs.find((s) => s.net === net);
  must(!!d, net + ' missing');
  if (!d) return;
  const top = Math.min(d.y1, d.y2);
  const bus = net === 'drop-T1-410' ? busI : busII;
  const by = net === 'drop-T1-410' ? busIy : busIIy;
  must(Math.abs(top - by) < 1.5, net + ' not from bus');
  reach({ x: d.cx, y: top }, { x: bus.cx, y: by }, net + ' not on bus');
});
must(pg.nodes.filter((n) => n.sym === 'cap').length === 2, 'cap symbols');
must(pg.nodes.some((n) => n.t === 'T' && /T1-410/.test(n.s || '')), '410 label');
must(pg.nodes.some((n) => n.t === 'T' && /T2-420/.test(n.s || '')), '420 label');
apart({ x: busI.cx, y: busIy }, { x: busII.cx, y: (busII.y1 + busII.y2) / 2 }, 'bus I directly tied to bus II');
const uf445 = withLink(t445.in.x, t445.in.y, t445.out.x, t445.out.y);
const s445 = segs.concat([{ horiz: false, x1: t445.in.x, y1: t445.in.y, x2: t445.out.x, y2: t445.out.y }]);
must(
  linked(s445, uf445, busI.x2, busIy, busII.x1, (busII.y1 + busII.y2) / 2, TOL) ||
    linked(s445, uf445, busI.x1, busIy, Math.max(busII.x1, busII.x2), (busII.y1 + busII.y2) / 2, TOL),
  '445 should be the only tie'
);
reach(t445.in, { x: (busI.x1 + busI.x2) / 2, y: busIy }, '445 top not on bus I');

reach(taps(knife('T1-411-5')).in, { x: busI.cx, y: busIy }, '411-5 not on bus I');
reach(taps(knife('T1-411-1')).in, taps(knife('T1-411-5')).in, '411 small bus');
must(segs.filter((s) => s.net === 'drop-T1-411').length === 1, '411 drops');
const d425 = segs.filter((s) => s.net === 'drop-T2-425');
must(d425.length === 1, '425 drop missing');
must(Math.abs(d425[0].cx - taps(knife('T2-425-1')).in.x) < 1, '425 drop not on -1');
must(segs.some((s) => s.net === 'drop-T2-425' && /漏画/.test(String(ctx.layout.wires.find((w) => w.net === 'drop-T2-425').note || ''))), '425 note');
reach(taps(knife('T2-421-1')).in, { x: busII.cx, y: (busII.y1 + busII.y2) / 2 }, '421-1 not on bus II');
must(segs.filter((s) => s.net === 'drop-T2-421').length === 1, '421 drops');

function colorNet(net, hex) {
  const hit = segs.filter((s) => s.net === net);
  must(hit.length >= 1 && hit.every((s) => s.c === hex), net + ' color ' + (hit[0] && hit[0].c));
}
colorNet('blue-412', '#2f7bff');
colorNet('blue-422', '#2f7bff');
colorNet('purple-T1-413-1', '#d24dff');
colorNet('green-T1-413-2', '#00e676');
colorNet('purple-T2-425-1', '#d24dff');
colorNet('green-T2-426-2', '#00e676');
reach(taps(knife('T1-412')).out, taps(knife('U1-XBP')).in, '412 blue to XBP');
reach(taps(knife('T2-422')).out, taps(knife('U2-XBP')).in, '422 blue to XBP');
[
  ['T1-413-1', 'U1-X01-1'],
  ['T1-413-2', 'U1-X01-1'],
  ['T1-414-1', 'U1-X01-2'],
  ['T1-414-2', 'U1-X01-2'],
  ['T1-415-1', 'U1-X02-1'],
  ['T1-415-2', 'U1-X02-1'],
  ['T1-416-1', 'U1-X02-2'],
  ['T1-416-2', 'U1-X02-2'],
  ['T2-426-1', 'U2-X02-2'],
  ['T2-426-2', 'U2-X02-2'],
  ['T2-425-1', 'U2-X02-1'],
  ['T2-425-2', 'U2-X02-1'],
  ['T2-424-1', 'U2-X01-2'],
  ['T2-424-2', 'U2-X01-2'],
  ['T2-423-1', 'U2-X01-1'],
  ['T2-423-2', 'U2-X01-1']
].forEach(([src, dst]) => {
  reach(taps(knife(src)).out, taps(knife(dst)).in, src + ' should reach ' + dst + ' inlet');
});
apart(taps(knife('T1-412')).out, taps(knife('U1-X01-1')).in, '412 blue merged into U11');
apart(taps(knife('T2-422')).out, taps(knife('U2-X01-1')).in, '422 blue merged into UPS');
const u11in = taps(knife('U1-X01-1')).in;
must(
  segs.some((s) => !s.horiz && s.c === '#d24dff' && touches(s, u11in.x, u11in.y, TOL)),
  'U11 inlet drop not purple'
);
must(
  segs.some((s) => s.c === '#00e676' && (touches(s, u11in.x, u11in.y, TOL) || segs.some((p) => p.c === '#d24dff' && !p.horiz && touches(p, u11in.x, u11in.y, TOL) && (touches(s, p.x1, p.y1, TOL) || touches(s, p.x2, p.y2, TOL))))),
  'U11 green not joined to inlet'
);
must(!segs.some((s) => String(s.net || '').indexOf('ups-out') === 0), 'ups-out red stub');

reach(taps(knife('U1-X01-1')).out, taps(knife('U1-X01-2')).out, 'U1 X01 join');
reach(taps(knife('U1-X01-1')).out, { x: segs.find((s) => s.net === 'bus-U1').cx, y: (segs.find((s) => s.net === 'bus-U1').y1 + segs.find((s) => s.net === 'bus-U1').y2) / 2 }, 'U1 X01 to bus');
const uBusY = (segs.find((s) => s.net === 'bus-U1').y1 + segs.find((s) => s.net === 'bus-U1').y2) / 2;
const x01y = taps(knife('U1-X01-1')).out.y;
const x01a = taps(knife('U1-X01-1')).out.x;
const x01b = taps(knife('U1-X01-2')).out.x;
const drops = segs.filter(
  (s) =>
    !s.horiz &&
    s.cx >= Math.min(x01a, x01b) - 2 &&
    s.cx <= Math.max(x01a, x01b) + 2 &&
    Math.min(s.y1, s.y2) <= x01y + 2 &&
    Math.max(s.y1, s.y2) >= uBusY - 2
);
must(drops.length === 1 && Math.abs(drops[0].cx - x01a) < 1, 'U1 X01 single drop at -1');

const u2bus = segs.find((s) => s.net === 'bus-U2');
const u2y = (u2bus.y1 + u2bus.y2) / 2;
const p22 = taps(knife('U2-X02-2')).out;
const p21 = taps(knife('U2-X02-1')).out;
const d2 = segs.filter(
  (s) =>
    !s.horiz &&
    s.cx >= Math.min(p22.x, p21.x) - 2 &&
    s.cx <= Math.max(p22.x, p21.x) + 2 &&
    Math.min(s.y1, s.y2) <= p22.y + 2 &&
    Math.max(s.y1, s.y2) >= u2y - 2
);
must(d2.length === 1 && Math.abs(d2[0].cx - p22.x) < 1, 'U2 X02 drop at X02-2');

reach(taps(knife('U1-X11')).out, taps(knife('U1-Y01')).in, 'X11 to Y01');
reach(taps(knife('U1-X12')).out, taps(knife('U1-Y02')).in, 'X12 to Y02');
reach(taps(knife('U2-X11')).out, taps(knife('U2-Y01')).in, 'U2 X11 to Y01');
reach(taps(knife('U2-X12')).out, taps(knife('U2-Y02')).in, 'U2 X12 to Y02');
const mid = segs.find((s) => s.net === 'bus-U1-mid');
const bot = segs.find((s) => s.net === 'bus-U1-bot');
const mid2 = segs.find((s) => s.net === 'bus-U2-mid');
const bot2 = segs.find((s) => s.net === 'bus-U2-bot');
const midY = (mid.y1 + mid.y2) / 2;
const botY = (bot.y1 + bot.y2) / 2;
reach(taps(knife('U1-Y02')).out, { x: mid.cx, y: midY }, 'U1 Y02 to mid bus');
reach(taps(knife('U1-Y01')).out, { x: bot.cx, y: botY }, 'U1 Y01 to bottom bus');
reach(taps(knife('U2-Y02')).out, { x: mid2.cx, y: (mid2.y1 + mid2.y2) / 2 }, 'U2 Y02 to mid bus');
reach(taps(knife('U2-Y01')).out, { x: bot2.cx, y: (bot2.y1 + bot2.y2) / 2 }, 'U2 Y01 to bottom bus');
apart(taps(knife('U1-X11')).out, { x: mid.cx, y: midY }, 'X11 should not tie the mid bus');
apart(taps(knife('U1-X11')).out, { x: mid2.cx, y: (mid2.y1 + mid2.y2) / 2 }, 'U1 X11 should not tie U2 mid bus');
apart(taps(knife('U1-X11')).out, { x: bot2.cx, y: (bot2.y1 + bot2.y2) / 2 }, 'U1 X11 should not tie U2 bottom bus');
apart(
  { x: segs.find((s) => s.net === 'bus-U1').cx, y: uBusY },
  { x: u2bus.cx, y: u2y },
  'U1 bus tied to U2 bus'
);
reach(taps(knife('U1-Y14-4')).in, { x: mid.cx, y: (mid.y1 + mid.y2) / 2 }, 'Y14 not on mid bus');
must(segs.filter((s) => s.net === 'drop-U1-Y14').length === 1, 'Y14 drops');
must(Math.abs(segs.find((s) => s.net === 'drop-U1-Y14').cx - taps(knife('U1-Y14-1')).in.x) < 1, 'Y14 drop not on -1');
must(Math.abs(segs.find((s) => s.net === 'drop-U2-Y14').cx - taps(knife('U2-Y14-1')).in.x) < 1, 'U2 Y14 drop not on -1');

['U2-X11', 'U2-X12'].forEach((cad) => {
  const o = taps(knife(cad)).out;
  const stub = segs.filter(
    (s) =>
      s.horiz &&
      s.c === '#ce0c18' &&
      Math.abs((s.y1 + s.y2) / 2 - o.y) > 3 &&
      Math.abs((s.y1 + s.y2) / 2 - o.y) < 40 &&
      o.x >= Math.min(s.x1, s.x2) - 2 &&
      o.x <= Math.max(s.x1, s.x2) + 2 &&
      Math.abs(s.x1 - s.x2) < 80
  );
  must(stub.length === 0, cad + ' short red under knife');
});

const report = { ok: fail.length === 0, fail, knives: knives.length, busGap: gap, missing: pg._lay5a1Missing || [] };
console.log(JSON.stringify(report, null, 2));
process.exit(report.ok ? 0 : 1);
