#!/usr/bin/env node
/**
 * 按 5A1 CAD 网表生成显式布局 data/hx-5a1-cad-layout.json。
 * 全部竖刀 60×98。并联组：-1 处一根竖线挂下 + 小横母线 + 多把竖刀。
 */
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const root = path.join(path.dirname(fileURLToPath(import.meta.url)), '..');
const outPath = path.join(root, 'data/hx-5a1-cad-layout.json');

const KW = 60;
const KH = 98;
const OUT_DY = (114 * KH) / 124;
const PITCH = 76;
const RED = '#ce0c18';
const BLUE = '#2f7bff';
const GREEN = '#00e676';
const PURPLE = '#d24dff';

const BUS = 44;
const TOP = 88;
const SMALL = TOP - 18;

const UPS_END = TOP + OUT_DY + 78;
const UPS_JOIN = UPS_END + 22;
const XROW = UPS_JOIN + 28;
const UBUS = XROW + OUT_DY + 26;
const XFEED = UBUS + 18;
const Y02 = XFEED + OUT_DY + 28;
const MIDBUS = Y02 + OUT_DY + 18;
const YSMALL = MIDBUS + 16;
const YMID = YSMALL + 14;
const Y01 = YMID + KH + 52;
const BOTBUS = Y01 + OUT_DY + 20;
const YBSMALL = BOTBUS + 16;
const YBOT = YBSMALL + 14;

const knives = [];
const wires = [];
const texts = [];
const boxes = [];
const openEnds = [];
const by = {};

function cin(k) {
  return { x: k.x + KW / 2, y: k.y };
}
function cout(k) {
  return { x: k.x + KW / 2, y: k.y + OUT_DY };
}

function addKnife(cad, x, y, out, outWire, amp, use) {
  const k = { cad, x, y, out: !!out, outWire: !!outWire, amp, use };
  knives.push(k);
  by[cad] = k;
  return k;
}

function wire(net, pts, c, extra) {
  const w = { net, pts, c: c || RED, sw: (extra && extra.sw) || (String(net).indexOf('bus-') === 0 ? 3 : 2) };
  if (extra && extra.open) w.open = extra.open;
  if (extra && extra.note) w.note = extra.note;
  wires.push(w);
  const opens = !extra || !extra.open ? [] : extra.open === 'both' ? ['start', 'end'] : [extra.open];
  if (opens.indexOf('start') >= 0) openEnds.push({ x: pts[0][0], y: pts[0][1], why: net });
  if (opens.indexOf('end') >= 0) openEnds.push({ x: pts[pts.length - 1][0], y: pts[pts.length - 1][1], why: net });
  return w;
}
function V(net, x, y0, y1, c, extra) {
  if (Math.abs(y0 - y1) < 0.5) return;
  wire(net, [[x, y0], [x, y1]], c, extra);
}
function H(net, x0, x1, y, c, extra) {
  if (Math.abs(x0 - x1) < 0.5) return;
  const a = Math.min(x0, x1);
  const b = Math.max(x0, x1);
  wire(net, [[a, y], [b, y]], c, extra);
}

/* 411 仍从 212 起，中间间距不动。401 左移，给 410 留出空隙。 */
let cursor = 212;
function put(cad, opt) {
  const o = opt || {};
  const k = addKnife(cad, cursor, TOP, o.out, o.outWire, o.amp, o.use);
  cursor += PITCH + (o.gap || 0);
  return k;
}

const amp411 = ['160A', '250A', '250A', '630A', '630A'];
const use411 = ['5A1市电测试箱', '备用', '备用', '备用', '备用'];
addKnife('T1-401', 0, TOP, false, false, '4000A', '总输出柜');
for (let i = 1; i <= 5; i++) {
  put('T1-411-' + i, { out: true, amp: amp411[i - 1], use: use411[i - 1], gap: i === 5 ? 20 : 0 });
}
put('T1-412', { out: true, outWire: true, amp: '4000A', use: 'U1大旁路', gap: 22 });
for (let n = 413; n <= 416; n++) {
  const ups = 'U1' + (n - 412);
  put('T1-' + n + '-1', { out: true, outWire: true, amp: '1250A', use: ups + '主路' });
  put('T1-' + n + '-2', { out: true, outWire: true, amp: '1250A', use: ups + '旁路', gap: 20 });
}
put('445', { outWire: true, amp: '3200A', use: '母联', gap: 260 });

const t2pairs = [
  ['426', 'U24'],
  ['425', 'U23'],
  ['424', 'U22'],
  ['423', 'U21']
];
t2pairs.forEach(([n, ups]) => {
  put('T2-' + n + '-2', { out: true, outWire: true, amp: '1600A', use: ups + '旁路' });
  put('T2-' + n + '-1', {
    out: true,
    outWire: true,
    amp: '1600A',
    use: ups + '主路',
    gap: 20
  });
});
put('T2-422', { out: true, outWire: true, amp: '4000A', use: 'U2大旁路', gap: 22 });
const amp421 = { 5: '630A', 4: '630A', 3: '250A', 2: '250A', 1: '160A' };
const use421 = { 5: '备用', 4: '备用', 3: '备用', 2: '备用', 1: '5A3市电测试箱' };
for (let i = 5; i >= 1; i--) {
  put('T2-421-' + i, { out: true, amp: amp421[i], use: use421[i], gap: i === 1 ? 136 : 0 });
}
put('T2-402', { amp: '4000A', use: '总输出柜' });

function pairCenterX(a, b) {
  return (cin(by[a]).x + cin(by[b]).x) / 2;
}
function atCenter(cx, y, cad, out, outWire, amp, use) {
  return addKnife(cad, cx - KW / 2, y, out, outWire, amp, use);
}

atCenter(cin(by['T1-412']).x, XROW, 'U1-XBP', false, true, '4000A', '维修旁路');
atCenter(pairCenterX('T1-413-1', 'T1-413-2'), XROW, 'U1-X01-1', false, true, '1000A', 'UPS1进线');
atCenter(pairCenterX('T1-414-1', 'T1-414-2'), XROW, 'U1-X01-2', false, true, '1000A', 'UPS2进线');
atCenter(pairCenterX('T1-415-1', 'T1-415-2'), XROW, 'U1-X02-1', false, true, '1000A', 'UPS3进线');
atCenter(pairCenterX('T1-416-1', 'T1-416-2'), XROW, 'U1-X02-2', false, true, '1000A', 'UPS4进线');
atCenter((cin(by['U1-X01-1']).x + cin(by['U1-X01-2']).x) / 2, XFEED, 'U1-X12', true, true, '2000A', 'UPS出线');
atCenter(48, XFEED, 'U1-X11', true, true, '2000A', 'UPS出线');

atCenter(pairCenterX('T2-426-2', 'T2-426-1'), XROW, 'U2-X02-2', false, true, '1000A', 'UPS4进线');
atCenter(pairCenterX('T2-425-2', 'T2-425-1'), XROW, 'U2-X02-1', false, true, '1000A', 'UPS3进线');
atCenter(pairCenterX('T2-424-2', 'T2-424-1'), XROW, 'U2-X01-2', false, true, '1000A', 'UPS2进线');
atCenter(pairCenterX('T2-423-2', 'T2-423-1'), XROW, 'U2-X01-1', false, true, '1000A', 'UPS1进线');
atCenter(cin(by['T2-422']).x, XROW, 'U2-XBP', false, true, '4000A', '维修旁路');
atCenter((cin(by['U2-X02-2']).x + cin(by['U2-X02-1']).x) / 2, XFEED, 'U2-X11', true, true, '2000A', 'UPS出线');
atCenter((cin(by['U2-X02-1']).x + cin(by['U2-X01-2']).x) / 2, XFEED, 'U2-X12', true, true, '2000A', 'UPS出线');

const YUSE = {
  'U1-Y14': ['5A2-A列尾', '5A2-B列尾', '5A2-C列尾', '5A2-D列尾'],
  'U1-Y15': ['5A2-E列尾', '5A2-F列尾', '5A2-G列尾', '备用'],
  'U1-Y16': ['5A2-H列尾', '5A2-J列尾', '5A2-K列尾', '备用'],
  'U1-Y11': ['5A1-A列头', '5A1-B列头', '5A1-C列头', '5A1-D列头'],
  'U1-Y12': ['5A1-E列头', '5A1-F列头', '5A1-G列头', '备用'],
  'U1-Y13': ['5A1-H列头', '5A1-J列头', '5A1-K列头', '备用'],
  'U2-Y14': ['5A3-K列头', '5A3-J列头', '5A3-H列头', '5A3-G列头'],
  'U2-Y15': ['5A3-F列头', '5A3-E列头', '5A3-D列头', '备用'],
  'U2-Y16': ['5A3-C列头', '5A3-B列头', '5A3-A列头', '电力监测电源'],
  'U2-Y11': ['5A4-K列尾', '5A4-J列尾', '5A4-H列尾', '5A4-G列尾'],
  'U2-Y12': ['5A4-F列尾', '5A4-E列尾', '5A4-D列尾', '备用'],
  'U2-Y13': ['5A4-C列尾', '5A4-B列尾', '5A4-A列尾', '备用']
};

const delta = cin(by['U2-X12']).x - cin(by['U1-X12']).x;
const yPitch = Math.max(64, Math.min(72, Math.floor((delta - 210) / 12)));
if (delta - 210 < 12 * 64) {
  console.error('Y pitch too tight', { delta, yPitch });
  process.exit(1);
}

function placeFeederGroups(side, feederCad, names, yKnife, pitch) {
  const groups = [];
  let x = by[feederCad].x + PITCH;
  names.forEach((name) => {
    const cads = [1, 2, 3, 4].map((i) => side + '-' + name + '-' + i);
    const uses = YUSE[side + '-' + name];
    cads.forEach((cad, i) => {
      addKnife(cad, x + i * pitch, yKnife, true, false, '400A', uses[i]);
    });
    groups.push(cads);
    x += 4 * pitch + 18;
  });
  return groups;
}

addKnife('U1-Y02', by['U1-X12'].x, Y02, false, true, '2000A', '数据机房5A2');
addKnife('U1-Y01', by['U1-X11'].x, Y01, false, true, '2000A', '数据机房5A1');
addKnife('U2-Y02', by['U2-X12'].x, Y02, false, true, '2000A', '数据机房5A3');
addKnife('U2-Y01', by['U2-X11'].x, Y01, false, true, '2000A', '数据机房5A4');

const gU1mid = placeFeederGroups('U1', 'U1-Y02', ['Y14', 'Y15', 'Y16'], YMID, yPitch);
const gU1bot = placeFeederGroups('U1', 'U1-Y01', ['Y11', 'Y12', 'Y13'], YBOT, yPitch);
const gU2mid = placeFeederGroups('U2', 'U2-Y02', ['Y14', 'Y15', 'Y16'], YMID, yPitch);
const gU2bot = placeFeederGroups('U2', 'U2-Y01', ['Y11', 'Y12', 'Y13'], YBOT, yPitch);

function hang(id, cads, busY, smallY, note) {
  const ks = cads.map((c) => by[c]);
  const drop = ks.find((k) => /-1$/.test(k.cad));
  if (!drop) throw new Error('no -1 ' + id);
  const xs = ks.map((k) => cin(k).x);
  H('small-' + id, Math.min.apply(null, xs), Math.max.apply(null, xs), smallY, RED, note ? { note } : null);
  V('drop-' + id, cin(drop).x, busY, smallY, RED, note ? { note } : null);
  ks.forEach((k) => V('stub-' + k.cad, cin(k).x, smallY, cin(k).y));
}

const k401 = by['T1-401'];
const k402 = by['T2-402'];
const k445 = by['445'];
const riser401 = (k401.x + KW + by['T1-411-1'].x) / 2;
const riser402 = (by['T2-421-1'].x + KW + k402.x) / 2;
const riser445 = (k445.x + KW + by['T2-426-2'].x) / 2;
const BUS_OVER = 18;

H('bus-I', cin(k401).x - BUS_OVER, cin(k445).x, BUS);
H('bus-II', riser445, cin(k402).x + BUS_OVER, BUS);
V('drop-401', cin(k401).x, BUS, cin(k401).y);
V('drop-402', cin(k402).x, BUS, cin(k402).y);
openEnds.push({ x: cin(k401).x - BUS_OVER, y: BUS, why: 'bus-I外端' });
openEnds.push({ x: cin(k402).x + BUS_OVER, y: BUS, why: 'bus-II外端' });

const XF_W = 64;
const XF_H = 88;
const INK = '#d5deea';
function addTransformer(name, side, knife) {
  const x = side === 'left' ? knife.x - 14 - XF_W : knife.x + KW + 14;
  const cx = x + 26;
  const y = 70;
  const out = cout(knife);
  /* 两圆相切，竖线不穿过绕组：上引线止于上圆顶，下引线起于下圆底。viewBox 与 64×88 框 1:1。 */
  wire('xfmr-' + name, [[cx, BUS + 12], [cx, y + 1]], INK, { open: 'start' });
  wire('xfmr-' + name, [[cx, y + 61], [cx, out.y], [out.x, out.y]], INK);
  boxes.push({ sym: 'xfmr', id: name, x: x, y: y, w: XF_W, h: XF_H });
  texts.push({
    s: name + '\n2000kVA',
    x: x - 4,
    y: out.y + 10,
    w: 72,
    h: 28,
    fs: 11,
    c: '#e7eef8'
  });
}
addTransformer('5A1-T1', 'left', k401);
addTransformer('5A1-T2', 'right', k402);

const CAP_W = 120;
/* 与 hx-draw capCabinetSvg viewBox 高度一致：顶上 60×98 刀，刀下互感器/电抗/电容/接地 */
const CAP_H = 252;
const CAP_GAP = 16;
function addCap(id, leftKnife, rightKnife) {
  const x = leftKnife.x + KW + CAP_GAP;
  if (x + CAP_W + CAP_GAP > rightKnife.x) {
    console.error('cap crowded', id, { x, right: rightKnife.x });
    process.exit(1);
  }
  const cx = x + CAP_W / 2;
  V('drop-' + id, cx, BUS, TOP, RED, { open: 'end' });
  boxes.push({ sym: 'cap', id: id, x: x, y: TOP, w: CAP_W, h: CAP_H });
  texts.push({
    s: id + '\n630A\n电容补偿柜',
    x: x + 8,
    y: TOP + CAP_H + 2,
    w: CAP_W - 8,
    h: 42,
    fs: 11,
    c: '#e7eef8',
    capDraw: 1
  });
}
addCap('T1-410', k401, by['T1-411-1']);
addCap('T2-420', by['T2-421-1'], k402);

hang('T1-411', [1, 2, 3, 4, 5].map((i) => 'T1-411-' + i), BUS, SMALL);
V('in-T1-412', cin(by['T1-412']).x, BUS, cin(by['T1-412']).y);
V('blue-412', cout(by['T1-412']).x, cout(by['T1-412']).y, cin(by['U1-XBP']).y, BLUE);

for (let n = 413; n <= 416; n++) {
  const id = 'T1-' + n;
  hang(id, [id + '-1', id + '-2'], BUS, SMALL);
}
[
  ['T1-413-1', 'T1-413-2', 'U1-X01-1'],
  ['T1-414-1', 'T1-414-2', 'U1-X01-2'],
  ['T1-415-1', 'T1-415-2', 'U1-X02-1'],
  ['T1-416-1', 'T1-416-2', 'U1-X02-2']
].forEach(([main, byp, out]) => landUps(main, byp, out));

V('in-445', cin(k445).x, BUS, cin(k445).y);
const o445 = cout(k445);
const jog445 = o445.y + 16;
wire('tie-445', [[o445.x, o445.y], [o445.x, jog445], [riser445, jog445], [riser445, BUS]]);

['426', '425', '424', '423'].forEach((n) => {
  const id = 'T2-' + n;
  const note = n === '425' ? 'CAD漏画，已补' : null;
  hang(id, [id + '-2', id + '-1'], BUS, SMALL, note);
});
[
  ['T2-426-1', 'T2-426-2', 'U2-X02-2'],
  ['T2-425-1', 'T2-425-2', 'U2-X02-1'],
  ['T2-424-1', 'T2-424-2', 'U2-X01-2'],
  ['T2-423-1', 'T2-423-2', 'U2-X01-1']
].forEach(([main, byp, out]) => landUps(main, byp, out));
V('in-T2-422', cin(by['T2-422']).x, BUS, cin(by['T2-422']).y);
V('blue-422', cout(by['T2-422']).x, cout(by['T2-422']).y, cin(by['U2-XBP']).y, BLUE);
hang('T2-421', [5, 4, 3, 2, 1].map((i) => 'T2-421-' + i), BUS, SMALL);

function landUps(mainCad, bypCad, outCad) {
  const mouth = cin(by[outCad]);
  const main = cout(by[mainCad]);
  const byp = cout(by[bypCad]);
  const sw = 26;
  const sh = 22;
  V('purple-' + mainCad, main.x, main.y, UPS_JOIN, PURPLE);
  H('purple-' + mainCad, main.x, mouth.x, UPS_JOIN, PURPLE);
  V('purple-' + mainCad, mouth.x, UPS_JOIN, mouth.y, PURPLE);
  V('green-' + bypCad, byp.x, byp.y, UPS_JOIN, GREEN);
  H('green-' + bypCad, byp.x, mouth.x, UPS_JOIN, GREEN);
  boxes.push({ sym: 'rect', x: main.x - sw / 2, y: UPS_JOIN - sh * 2 - 2, w: sw, h: sh });
  boxes.push({ sym: 'inv', x: main.x - sw / 2, y: UPS_JOIN - sh, w: sw, h: sh });
  boxes.push({ sym: 'ssw', x: byp.x - 14, y: UPS_JOIN - 26, w: 28, h: 26 });
}

function joinDrop(a, b, dropCad, busY, net) {
  H('join-' + net, cout(by[a]).x, cout(by[b]).x, cout(by[a]).y);
  V('drop-' + net, cout(by[dropCad]).x, cout(by[dropCad]).y, busY);
}
joinDrop('U1-X01-1', 'U1-X01-2', 'U1-X01-1', UBUS, 'U1-X01');
joinDrop('U1-X02-1', 'U1-X02-2', 'U1-X02-1', UBUS, 'U1-X02');
joinDrop('U2-X02-2', 'U2-X02-1', 'U2-X02-2', UBUS, 'U2-X02');
joinDrop('U2-X01-2', 'U2-X01-1', 'U2-X01-2', UBUS, 'U2-X01');
V('xbp-U1', cout(by['U1-XBP']).x, cout(by['U1-XBP']).y, UBUS);
V('xbp-U2', cout(by['U2-XBP']).x, cout(by['U2-XBP']).y, UBUS);

H('bus-U1', cin(by['U1-X11']).x, cin(by['U1-X02-2']).x, UBUS);
H('bus-U2', cin(by['U2-X02-2']).x, cin(by['U2-XBP']).x, UBUS);
openEnds.push({ x: cin(by['U1-X02-2']).x, y: UBUS, why: 'bus-U1端部' });

['U1', 'U2'].forEach((side) => {
  const x11 = side + '-X11';
  const x12 = side + '-X12';
  const y01 = side + '-Y01';
  const y02 = side + '-Y02';
  V('in-' + x11, cin(by[x11]).x, UBUS, cin(by[x11]).y);
  V('in-' + x12, cin(by[x12]).x, UBUS, cin(by[x12]).y);
  V('link-' + x11, cout(by[x11]).x, cout(by[x11]).y, cin(by[y01]).y);
  V('link-' + x12, cout(by[x12]).x, cout(by[x12]).y, cin(by[y02]).y);
  V('out-' + y02, cout(by[y02]).x, cout(by[y02]).y, MIDBUS);
  V('out-' + y01, cout(by[y01]).x, cout(by[y01]).y, BOTBUS);
});

function busThrough(net, originCad, groups, busY) {
  const xs = [cin(by[originCad]).x];
  groups.forEach((g) => g.forEach((c) => xs.push(cin(by[c]).x)));
  const left = Math.min.apply(null, xs);
  const right = Math.max.apply(null, xs);
  H(net, left, right, busY);
  const origin = cin(by[originCad]).x;
  if (Math.abs(origin - left) < 1) openEnds.push({ x: right, y: busY, why: net + '端部' });
  else if (Math.abs(origin - right) < 1) openEnds.push({ x: left, y: busY, why: net + '端部' });
}
busThrough('bus-U1-mid', 'U1-Y02', gU1mid, MIDBUS);
busThrough('bus-U1-bot', 'U1-Y01', gU1bot, BOTBUS);
busThrough('bus-U2-mid', 'U2-Y02', gU2mid, MIDBUS);
busThrough('bus-U2-bot', 'U2-Y01', gU2bot, BOTBUS);
gU1mid.forEach((g) => hang(g[0].replace(/-1$/, ''), g, MIDBUS, YSMALL));
gU1bot.forEach((g) => hang(g[0].replace(/-1$/, ''), g, BOTBUS, YBSMALL));
gU2mid.forEach((g) => hang(g[0].replace(/-1$/, ''), g, MIDBUS, YSMALL));
gU2bot.forEach((g) => hang(g[0].replace(/-1$/, ''), g, BOTBUS, YBSMALL));

function label(k) {
  const head = k.cad === 'T1-401' || k.cad === 'T2-402' ? k.cad : k.cad.replace(/^T1-|^T2-|^U1-|^U2-/, '');
  const lines = [head, k.amp || '', k.use || ''].filter(Boolean);
  texts.push({
    s: lines.join('\n'),
    x: k.x - 8,
    y: k.y + KH + 2,
    w: 76,
    h: 38,
    fs: 10,
    c: '#e7eef8',
    k: k.cad
  });
}
knives.forEach(label);

texts.push({ s: '配电室 5A1', x: 16, y: 8, w: 160, h: 22, fs: 16, c: '#ffffff', al: 'l' });
texts.push({ s: 'I 段', x: riser401, y: BUS - 22, w: 48, h: 16, fs: 12, c: '#ffb4b4' });
texts.push({
  s: 'II 段',
  x: (riser445 + riser402) / 2 - 24,
  y: BUS - 22,
  w: 48,
  h: 16,
  fs: 12,
  c: '#ffb4b4'
});
['413', '414', '415', '416'].forEach((n, i) => {
  const cx = pairCenterX('T1-' + n + '-1', 'T1-' + n + '-2');
  texts.push({ s: 'U1' + (i + 1), x: cx - 16, y: UPS_JOIN - 46, w: 32, h: 14, fs: 11, c: '#d7c4ff' });
});
['426', '425', '424', '423'].forEach((n, i) => {
  const cx = pairCenterX('T2-' + n + '-2', 'T2-' + n + '-1');
  texts.push({ s: 'U2' + (4 - i), x: cx - 16, y: UPS_JOIN - 46, w: 32, h: 14, fs: 11, c: '#d7c4ff' });
});
let minX = 0;
let maxX = knives.reduce((m, k) => Math.max(m, k.x + KW), 0);
let maxY = knives.reduce((m, k) => Math.max(m, k.y + KH + 44), 0);
boxes.forEach((b) => {
  minX = Math.min(minX, b.x);
  maxX = Math.max(maxX, b.x + b.w);
  maxY = Math.max(maxY, b.y + b.h);
});
texts.forEach((t) => {
  minX = Math.min(minX, t.x);
  maxX = Math.max(maxX, t.x + (t.w || 0));
});
const layout = {
  version: '20261027cad5a1h',
  canvas: { w: Math.ceil(maxX + 48), h: Math.ceil(maxY + 16), x0: Math.floor(Math.min(0, minX - 16)) },
  knives: knives.map((k) => ({
    cad: k.cad,
    x: k.x,
    y: k.y,
    out: k.out,
    outWire: k.outWire
  })),
  wires,
  texts,
  boxes,
  openEnds,
  notes: [
    '全部竖刀 60×98，并联组在 -1 挂下',
    'T2-425 CAD 漏画挂线，已补',
    'U2-X11/X12 下短红线不画',
    '410/420 静态补偿柜，母线竖线接到刀上口',
    '紫线经整流逆变、绿线经静态开关，汇合后落到输出刀上口',
    '401/402 母线竖线接下刀上口；变压器出线进刀下口，竖线不搭母线'
  ]
};

function boxesOverlap(a, b) {
  return a.x < b.x + KW && a.x + KW > b.x && a.y < b.y + KH && a.y + KH > b.y;
}
const hits = [];
for (let i = 0; i < knives.length; i++) {
  for (let j = i + 1; j < knives.length; j++) {
    if (boxesOverlap(knives[i], knives[j])) hits.push(knives[i].cad + '×' + knives[j].cad);
  }
}
if (hits.length) {
  console.error('overlap', hits.slice(0, 12));
  process.exit(1);
}
const boxHits = [];
boxes.forEach((b) => {
  knives.forEach((k) => {
    if (b.x < k.x + KW && b.x + b.w > k.x && b.y < k.y + KH && b.y + b.h > k.y) {
      boxHits.push((b.id || b.sym) + '×' + k.cad);
    }
  });
});
if (boxHits.length) {
  console.error('box overlap', boxHits.slice(0, 8));
  process.exit(1);
}
const busI = wires.find((w) => w.net === 'bus-I');
const busII = wires.find((w) => w.net === 'bus-II');
const gap = busII.pts[0][0] - busI.pts[1][0];
if (gap < 24) {
  console.error('bus gap', gap);
  process.exit(1);
}
const scenes = JSON.parse(fs.readFileSync(path.join(root, 'data/hx-scenes.json'), 'utf8'));
const pg = scenes.pages.find((p) => p.name === '配电室5A1');
const sceneKeys = {};
pg.nodes.forEach((n) => {
  if ((n.t === 'I' || n.t === 'S') && n.k) {
    sceneKeys[String(n.k).replace(/^配电室5A1_/, '').replace(/_+$/, '')] = 1;
  }
});
const missing = knives.filter((k) => !sceneKeys[k.cad.replace(/-/g, '_')]).map((k) => k.cad);
const cadSet = {};
knives.forEach((k) => {
  cadSet[k.cad.replace(/-/g, '_')] = 1;
});
const extra = Object.keys(sceneKeys).filter((k) => !cadSet[k]);
if (missing.length || extra.length || knives.length !== 97) {
  console.error({ missing, extra, n: knives.length });
  process.exit(1);
}

fs.writeFileSync(outPath, JSON.stringify(layout));
const u1Right = Math.max.apply(null, gU1mid.flat().map((c) => by[c].x + KW));
const u2Left = Math.min.apply(null, gU2mid.flat().map((c) => by[c].x));
if (u1Right > by['U2-X11'].x - 8) {
  console.error('U1 Y band crosses U2-X11 column', { u1Right, x11: by['U2-X11'].x });
  process.exit(1);
}
console.log(
  JSON.stringify(
    {
      knives: knives.length,
      wires: wires.length,
      canvas: layout.canvas,
      yPitch,
      delta,
      busGap: gap,
      yBandGap: u2Left - u1Right,
      out: outPath
    },
    null,
    2
  )
);
