#!/usr/bin/env node
/** 只读。UPS 进线刀：同室紫/绿竖线差一截；红横线盖过 U1/U2 主列空档。不改布局。 */
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const dir = path.join(path.dirname(fileURLToPath(import.meta.url)), '../data/hx-power-rooms');
const PURPLE = '#d24dff';
const GREEN = '#00e676';
const RED = '#ce0c18';
const files = fs.readdirSync(dir).filter((f) => f.endsWith('.json') && f !== 'index.json').sort();

function upsIn(use) {
  return /UPS\d*进线|UPS进线\d+/.test(use || '');
}

function distSeg(px, py, a, b) {
  const dx = b[0] - a[0];
  const dy = b[1] - a[1];
  const l2 = dx * dx + dy * dy;
  if (l2 < 1e-9) return Math.hypot(px - a[0], py - a[1]);
  const t = Math.max(0, Math.min(1, ((px - a[0]) * dx + (py - a[1]) * dy) / l2));
  return Math.hypot(px - (a[0] + t * dx), py - (a[1] + t * dy));
}

function colGap(knives) {
  const cx = (k) => k.x + 30;
  const u1 = knives.filter((k) => /^U1-X/.test(k.cad || '')).map(cx);
  const u2 = knives.filter((k) => /^U2-X/.test(k.cad || '')).map(cx);
  if (!u1.length || !u2.length) return null;
  const a1 = Math.max(...u1);
  const b0 = Math.min(...u2);
  const b1 = Math.max(...u2);
  const a0 = Math.min(...u1);
  if (a1 < b0) return [a1, b0];
  if (b1 < a0) return [b1, a0];
  return null;
}

const short = [];
const span = [];
let rooms = 0;
let ups = 0;

for (const file of files) {
  const lay = JSON.parse(fs.readFileSync(path.join(dir, file), 'utf8'));
  const room = lay.room || file.replace('.json', '');
  rooms += 1;
  const wires = lay.wires || [];
  const knives = (lay.knives || []).filter((k) => upsIn(k.use));
  ups += knives.length;
  for (const k of knives) {
    const mx = k.x + 30;
    const my = k.y;
    const landed = { [PURPLE]: false, [GREEN]: false };
    wires.forEach((w) => {
      const c = (w.c || '').toLowerCase();
      if (c !== PURPLE && c !== GREEN) return;
      const pts = w.pts || [];
      for (let i = 1; i < pts.length; i++) {
        if (distSeg(mx, my, pts[i - 1], pts[i]) <= 2) landed[c] = true;
      }
    });
    wires.forEach((w, wi) => {
      const c = (w.c || '').toLowerCase();
      if (c !== PURPLE && c !== GREEN || landed[c]) return;
      const pts = w.pts || [];
      for (let i = 1; i < pts.length; i++) {
        const a = pts[i - 1];
        const b = pts[i];
        if (Math.abs(a[0] - b[0]) > 15 || Math.abs(a[1] - b[1]) <= Math.abs(a[0] - b[0])) continue;
        const near = Math.abs(a[1] - my) <= Math.abs(b[1] - my) ? a : b;
        const far = near === a ? b : a;
        const dy = Math.abs(near[1] - my);
        if (Math.abs(near[0] - mx) > 15 || dy <= 2 || dy >= 200) continue;
        if (Math.abs(far[1] - my) <= dy) continue;
        short.push({
          room, cad: k.cad, color: c === PURPLE ? '紫' : '绿', wi,
          dx: Math.round((near[0] - mx) * 10) / 10,
          gap: Math.round(dy * 10) / 10,
        });
      }
    });
  }
  const gap = colGap(lay.knives || []);
  if (!gap) continue;
  wires.forEach((w, wi) => {
    if ((w.c || '').toLowerCase() !== RED) return;
    const pts = w.pts || [];
    for (let i = 1; i < pts.length; i++) {
      const a = pts[i - 1];
      const b = pts[i];
      if (Math.abs(a[1] - b[1]) > 2.5) continue;
      const x0 = Math.min(a[0], b[0]);
      const x1 = Math.max(a[0], b[0]);
      if (x0 < gap[0] && x1 > gap[1]) {
        span.push({
          room, wi, y: Math.round(((a[1] + b[1]) / 2) * 100) / 100,
          x0: Math.round(x0 * 100) / 100, x1: Math.round(x1 * 100) / 100,
          gap0: Math.round(gap[0] * 100) / 100, gap1: Math.round(gap[1] * 100) / 100,
        });
      }
    }
  });
}

console.log('rooms', rooms, 'upsIn', ups, 'short', short.length, 'span', span.length);
for (const s of short) console.log('SHORT', s.room, s.cad, s.color, 'gap', s.gap, 'dx', s.dx, 'w', s.wi);
for (const s of span) console.log('SPAN', s.room, 'y', s.y, 'x', s.x0, s.x1, 'gap', s.gap0, s.gap1, 'w', s.wi);
