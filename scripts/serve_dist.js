/**
 * ISM 前端静态服务器 + /api 代理
 * 用于 serve dist/ 目录，代理 /api 请求到后端 8081
 */
const http = require('http');
const httpProxy = require('http-proxy');
const fs = require('fs');
const path = require('path');
const url = require('url');

const PORT = 7080;
const API_TARGET = 'http://127.0.0.1:8081';
const DIST_DIR = path.join(__dirname, 'dist');

const proxy = httpProxy.createProxyServer({
  target: API_TARGET,
  changeOrigin: true,
});

proxy.on('error', (err, req, res) => {
  console.error('Proxy error:', err.message);
  if (!res.headersSent) {
    res.writeHead(502, { 'Content-Type': 'text/plain' });
  }
  res.end('Backend unavailable');
});

const MIME_TYPES = {
  '.html': 'text/html; charset=utf-8',
  '.js': 'application/javascript; charset=utf-8',
  '.css': 'text/css; charset=utf-8',
  '.json': 'application/json; charset=utf-8',
  '.png': 'image/png',
  '.jpg': 'image/jpeg',
  '.jpeg': 'image/jpeg',
  '.gif': 'image/gif',
  '.svg': 'image/svg+xml',
  '.ico': 'image/x-icon',
  '.woff': 'font/woff',
  '.woff2': 'font/woff2',
  '.ttf': 'font/ttf',
  '.map': 'application/json',
};

const server = http.createServer((req, res) => {
  const parsed = url.parse(req.url);

  // Proxy /api/* 和 /login 到后端
  if (parsed.pathname.startsWith('/api') || parsed.pathname === '/login' || parsed.pathname.startsWith('/ws')) {
    // Strip /api prefix for backend
    if (parsed.pathname.startsWith('/api')) {
      req.url = req.url.replace(/^\/api/, '');
    }
    return proxy.web(req, res);
  }

  // Serve static files
  let filePath = path.join(DIST_DIR, parsed.pathname === '/' ? 'index.html' : parsed.pathname);

  // SPA fallback: if file doesn't exist, serve index.html
  if (!fs.existsSync(filePath)) {
    filePath = path.join(DIST_DIR, 'index.html');
  }

  const ext = path.extname(filePath);
  const contentType = MIME_TYPES[ext] || 'application/octet-stream';

  fs.readFile(filePath, (err, data) => {
    if (err) {
      res.writeHead(404, { 'Content-Type': 'text/plain' });
      res.end('Not Found');
      return;
    }
    res.writeHead(200, { 'Content-Type': contentType });
    res.end(data);
  });
});

server.listen(PORT, () => {
  console.log(`ISM Frontend serving dist/ on http://localhost:${PORT}`);
  console.log(`API proxy: ${API_TARGET}`);
});
