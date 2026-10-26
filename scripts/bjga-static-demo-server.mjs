#!/usr/bin/env node
/** 本地截图用：静态 bjga + /api/bjga/bootstrap 反代到 8082 */
import http from 'http';
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const ROOT = path.join(__dirname, '../ism-front-end-v2/public/bjga');
const BE = '127.0.0.1';
const BE_PORT = 8082;
const PORT = Number(process.env.BJGA_DEMO_PORT || 7099);

function loginToken() {
  return fetch(`http://${BE}:${BE_PORT}/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ Username: 'admin', password: 'e10adc3949ba59abbe56e057f20f883e' }),
  }).then((r) => r.json()).then((j) => j.data?.token || '');
}

const mime = {
  '.html': 'text/html',
  '.js': 'text/javascript',
  '.css': 'text/css',
  '.png': 'image/png',
  '.svg': 'image/svg+xml',
};

const server = http.createServer(async (req, res) => {
  const url = req.url || '/';
  if (url.startsWith('/api/')) {
    const target = `http://${BE}:${BE_PORT}${url.replace(/^\/api/, '')}`;
    const chunks = [];
    for await (const c of req) chunks.push(c);
    const body = Buffer.concat(chunks);
    const headers = { 'Content-Type': req.headers['content-type'] || 'application/json' };
    if (!url.startsWith('/api/login')) {
      headers.Authorization = await loginToken();
    } else if (req.headers.authorization) {
      headers.Authorization = req.headers.authorization;
    }
    const r = await fetch(target, { method: req.method || 'GET', headers, body: body.length ? body : undefined });
    const buf = Buffer.from(await r.arrayBuffer());
    res.writeHead(r.status, { 'Content-Type': r.headers.get('content-type') || 'application/json' });
    res.end(buf);
    return;
  }
  if (url === '/fixture/dcim.json') {
    const dcimPath = path.join(__dirname, '../ism_server_user/data/bjga/dcim.json');
    res.writeHead(200, { 'Content-Type': 'application/json' });
    fs.createReadStream(dcimPath).pipe(res);
    return;
  }
  const cleanUrl = (url.split('?')[0] || '/');
  let p = path.join(ROOT, cleanUrl === '/' ? 'index.html' : cleanUrl.replace(/^\//, ''));
  if (!p.startsWith(ROOT)) {
    res.writeHead(403);
    res.end('forbidden');
    return;
  }
  if (!fs.existsSync(p) || fs.statSync(p).isDirectory()) {
    res.writeHead(404);
    res.end('not found');
    return;
  }
  const ext = path.extname(p);
  res.writeHead(200, { 'Content-Type': mime[ext] || 'application/octet-stream' });
  fs.createReadStream(p).pipe(res);
});

server.listen(PORT, () => {
  console.log(`bjga demo http://127.0.0.1:${PORT}/index.html#/building/drill5`);
});
