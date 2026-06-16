const http = require('http');
const fs = require('fs');
const path = require('path');

const DIST = path.join(__dirname, 'ism-front-end-v2/dist');
const PORT = 7080;
const BACKEND = { host: '127.0.0.1', port: 8081 };

const MIME = {
  '.html': 'text/html', '.js': 'application/javascript', '.css': 'text/css',
  '.json': 'application/json', '.png': 'image/png', '.jpg': 'image/jpeg',
  '.svg': 'image/svg+xml', '.ico': 'image/x-icon', '.woff': 'font/woff',
  '.woff2': 'font/woff2', '.ttf': 'font/ttf', '.map': 'application/json',
};

http.createServer((req, res) => {
  // Proxy /api -> backend
  if (req.url.startsWith('/api')) {
    const opts = {
      hostname: BACKEND.host, port: BACKEND.port,
      path: req.url.replace('/api', '') || '/',
      method: req.method, headers: req.headers,
    };
    const proxy = http.request(opts, (pres) => {
      res.writeHead(pres.statusCode, pres.headers);
      pres.pipe(res);
    });
    proxy.on('error', () => { res.writeHead(502); res.end('Bad Gateway'); });
    req.pipe(proxy);
    return;
  }

  // Serve static files
  let filePath = path.join(DIST, req.url === '/' ? 'index.html' : req.url.split('?')[0]);
  const ext = path.extname(filePath);

  fs.readFile(filePath, (err, data) => {
    if (err) {
      // SPA fallback
      fs.readFile(path.join(DIST, 'index.html'), (e2, d2) => {
        if (e2) { res.writeHead(404); res.end('Not Found'); return; }
        res.writeHead(200, { 'Content-Type': 'text/html' });
        res.end(d2);
      });
      return;
    }
    res.writeHead(200, { 'Content-Type': MIME[ext] || 'application/octet-stream' });
    res.end(data);
  });
}).listen(PORT, () => {
  console.log(`Frontend on http://localhost:${PORT}, /api -> ${BACKEND.host}:${BACKEND.port}`);
});
