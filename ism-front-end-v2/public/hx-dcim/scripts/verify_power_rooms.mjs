#!/usr/bin/env node
/** 新配电室：能被架构树索引，刀数与网表一致，进线/补偿柜上口接到母线。5A1 不在这套布局里。 */
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const root = path.join(path.dirname(fileURLToPath(import.meta.url)), '..');
const dir = path.join(root, 'data/hx-power-rooms');
const index = JSON.parse(fs.readFileSync(path.join(dir, 'index.json'), 'utf8'));
const arch = JSON.parse(fs.readFileSync(path.join(root, 'data/hx-arch-tree.json'), 'utf8'));
function leaves() {
  const out = {};
  function walk(node) {
    const kids = node.children || [];
    if (!kids.length && node.name) out[node.name] = node;
    kids.forEach(walk);
  }
  (arch.branches || []).forEach(walk);
  return out;
}

function touches(wires, x, y) {
  return wires.some((w) => {
    const pts = w.pts || [];
    for (let i = 1; i < pts.length; i++) {
      const a = pts[i - 1];
      const b = pts[i];
      if (Math.abs(a[0] - b[0]) <= 1.6) {
        const cx = (a[0] + b[0]) / 2;
        if (Math.abs(cx - x) <= 4 && Math.min(a[1], b[1]) - 3 <= y && y <= Math.max(a[1], b[1]) + 3) return true;
      }
      if (Math.abs(a[1] - b[1]) <= 1.6) {
        const cy = (a[1] + b[1]) / 2;
        if (Math.abs(cy - y) <= 4 && Math.min(a[0], b[0]) - 3 <= x && x <= Math.max(a[0], b[0]) + 3) return true;
      }
    }
    return false;
  });
}

const tree = leaves();
const fails = [];
const gaps = [];
if ((index.rooms || []).some((r) => r.code === '5A1')) fails.push('index 含 5A1，不能覆盖基线');

for (const room of index.rooms || []) {
  const name = room.name || ('配电室' + room.code);
  const leafName = room.code.startsWith('ECC') ? room.code : name;
  const leaf = tree[leafName];
  if (!leaf || !leaf.pageId) fails.push(leafName + ' 未进架构树');
  const lay = JSON.parse(fs.readFileSync(path.join(dir, room.file), 'utf8'));
  const n = (lay.knives || []).length;
  if (n !== lay.meta.knifeCount) fails.push(room.code + ' 刀数 ' + n + ' != meta ' + lay.meta.knifeCount);
  if (room.knives != null && room.knives !== n) fails.push(room.code + ' 索引刀数不符');
  const wires = lay.wires || [];
  const ends = []
    .concat(lay.meta.incomers || [])
    .concat(lay.meta.caps || []);
  ends.forEach((cad) => {
    const k = (lay.knives || []).find((x) => x.cad === cad);
    if (!k) {
      fails.push(room.code + ' 缺刀 ' + cad);
      return;
    }
    if (!touches(wires, k.x + 30, k.y)) fails.push(room.code + ' ' + cad + ' 上口未接母线');
  });
  if ((lay.meta.incomers || []).length >= 2 && wires.filter((w) => w.net === 'bus').length < 2) {
    fails.push(room.code + ' 两段进线但红色母线少于 2');
  }
  if (!(lay.meta.incomers || []).length || !(lay.meta.caps || []).length) {
    gaps.push(room.code + ' 进线' + (lay.meta.incomers || []).length + ' 补偿' + (lay.meta.caps || []).length);
  }
}

if (fails.length) {
  console.log(fails.join('\n'));
  process.exit(1);
}
console.log('rooms', (index.rooms || []).length, 'ok');
if (gaps.length) console.log('CAD 未标全进线或补偿:\n' + gaps.join('\n'));
