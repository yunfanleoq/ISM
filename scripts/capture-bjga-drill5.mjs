#!/usr/bin/env node
import puppeteer from 'puppeteer-core';
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const OUT = path.join(__dirname, '../docs/assets/bjga-drill5');
const BASE = process.env.BJGA_DEMO_URL || 'http://127.0.0.1:7099/index.html';
const CHROME = process.env.CHROME_PATH || '/usr/local/bin/google-chrome';
const dcim = JSON.parse(fs.readFileSync(path.join(__dirname, '../ism_server_user/data/bjga/dcim.json'), 'utf8'));

async function shot(page, name) {
  await new Promise((r) => setTimeout(r, 600));
  await page.screenshot({ path: path.join(OUT, name), fullPage: false });
  console.log('saved', name);
}

const browser = await puppeteer.launch({
  executablePath: CHROME,
  headless: 'new',
  args: ['--no-sandbox', '--disable-setuid-sandbox', '--window-size=1920,1080'],
});
const page = await browser.newPage();
await page.setViewport({ width: 1920, height: 1080, deviceScaleFactor: 1 });

await page.goto(BASE, { waitUntil: 'networkidle0', timeout: 180000 });
await page.waitForFunction(() => window.App && App.pages && App.pages.building, { timeout: 120000 });
await page.evaluate((payload) => {
  window.__dcimPanicOff = true;
  window.DCIM = payload;
  if (window.App._ready) window.App._ready();
  if (window.App.boot) window.App.boot();
  window.location.hash = '#/building/drill5';
  if (window.App.render) window.App.render();
}, dcim);

await page.waitForSelector('.campus-bld', { timeout: 60000 });
await shot(page, '01-campus.png');

await page.click('.campus-bld');
await page.waitForSelector('.bld-stage2', { timeout: 30000 });
await shot(page, '02-building.png');

const demoRoom = '3层2#数据机房';
await page.evaluate((room) => {
  window.location.hash = '#/building/room/' + encodeURIComponent(room);
  if (window.App.render) window.App.render();
}, demoRoom);
await page.waitForSelector('.rm-plan, .rm-grid-wrap, .rm-wrap', { timeout: 60000 });
await shot(page, '03-room.png');

await page.evaluate((cab) => {
  var self = window.App.pages.building;
  if (self && self.viewCab) {
    self.state.cab = cab;
    self.viewCab(cab, null);
  }
}, '3M2-无号');
await page.waitForSelector('.rk-rack, .rk-wrap', { timeout: 60000 });
await shot(page, '04-cabinet.png');

await page.evaluate(() => {
  var sn = 'SIM-3M2-无号-01';
  var self = window.App.pages.building;
  if (self && self.viewDevice) self.viewDevice(sn, null);
});
await page.waitForSelector('.dev-wrap', { timeout: 60000 });
await shot(page, '05-device.png');

await browser.close();
