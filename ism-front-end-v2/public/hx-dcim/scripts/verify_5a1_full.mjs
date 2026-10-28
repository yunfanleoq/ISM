#!/usr/bin/env node
/**
 * 5A1 全量验收：竖刀 60×98、触点压线、红/彩线端点不悬空（开放端除外）。
 * 用法：node verify_5a1_full.mjs [--dom]
 */
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';
import { bootHx, prepRoom, segments, touches } from './lib_5a1.mjs';

const root = path.join(path.dirname(fileURLToPath(import.meta.url)), '..');
const TOL = 2.5;
const ctx = bootHx(root);
const pg = prepRoom(ctx, '配电室5A1');
const tapsOf = ctx.sandbox.window.HX5A1.knifeTaps;
const knives = pg.nodes.filter((n) => (n.t === 'I' || n.t === 'S') && n.cad);
const segs = segments(pg.nodes);

const sizeIssues = knives.filter((n) => n.w !== 60 || n.h !== 98 || n.orient === 'h' || Number(n.rot) === 90);

const missIn = [];
const missOut = [];
const tapPoints = [];
knives.forEach((kn) => {
  const t = tapsOf(kn);
  tapPoints.push(t.in, t.out);
  if (!segs.some((s) => touches(s, t.in.x, t.in.y, TOL))) missIn.push(kn.cad);
  if (kn.outWire && !segs.some((s) => touches(s, t.out.x, t.out.y, TOL))) missOut.push(kn.cad);
});

const open = pg._openEnds || [];
function near(x, y, pts, tol) {
  return pts.some((p) => Math.hypot(p.x - x, p.y - y) <= tol);
}
const dangling = [];
segs.forEach((s) => {
  [ { x: s.x1, y: s.y1 }, { x: s.x2, y: s.y2 } ].forEach((pt) => {
    if (near(pt.x, pt.y, tapPoints, TOL)) return;
    if (near(pt.x, pt.y, open, TOL)) return;
    if (segs.some((o) => o !== s && touches(o, pt.x, pt.y, TOL))) return;
    dangling.push({ x: Math.round(pt.x), y: Math.round(pt.y), net: s.net });
  });
});

const report = {
  knives: knives.length,
  sizeIssues: sizeIssues.length,
  sizeIssueSample: sizeIssues.slice(0, 5).map((n) => n.cad),
  missIn: missIn.length,
  missOut: missOut.length,
  missInSample: missIn.slice(0, 8),
  missOutSample: missOut.slice(0, 8),
  danglingEndpoints: dangling.length,
  danglingSample: dangling.slice(0, 8),
  ok: sizeIssues.length === 0 && missIn.length === 0 && missOut.length === 0 && dangling.length === 0
};

let domReport = null;
if (process.argv.includes('--dom')) {
  try {
    const puppeteer = await import('puppeteer');
    const browser = await puppeteer.default.launch({ headless: true });
    const page = await browser.newPage();
    await page.setViewport({ width: 1400, height: 900 });
    await page.goto('http://127.0.0.1:7080/hx-dcim/index.html?v=20261027cad5a1f#/power', {
      waitUntil: 'networkidle2',
      timeout: 120000
    });
    await page.waitForTimeout(8000);
    domReport = await page.evaluate(() => {
      const knives = [...document.querySelectorAll('[data-hx-knife]')];
      const horiz = knives.filter((el) => el.classList.contains('hx-sym-h')).length;
      return { count: knives.length, horiz };
    });
    await browser.close();
    report.dom = domReport;
    report.domOk = domReport.count >= 97 && domReport.horiz === 0;
    report.ok = report.ok && report.domOk;
  } catch (e) {
    report.domError = String(e.message || e);
  }
}

console.log(JSON.stringify(report, null, 2));
process.exit(report.ok ? 0 : 1);
