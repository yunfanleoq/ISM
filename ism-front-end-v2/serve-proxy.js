const http = require('http')
const fs = require('fs')
const path = require('path')
const url = require('url')

const DIST = __dirname + '/dist'
const BACKEND = 'http://127.0.0.1:8081'
const PORT = 7080

const MIME = {
  '.html': 'text/html',
  '.js': 'application/javascript',
  '.css': 'text/css',
  '.png': 'image/png',
  '.jpg': 'image/jpeg',
  '.svg': 'image/svg+xml',
  '.woff': 'font/woff',
  '.woff2': 'font/woff2',
  '.ttf': 'font/ttf',
  '.json': 'application/json',
}

function serveFile(res, filePath) {
  const ext = path.extname(filePath)
  const ct = MIME[ext] || 'application/octet-stream'
  try {
    const data = fs.readFileSync(filePath)
    res.writeHead(200, { 'Content-Type': ct })
    res.end(data)
  } catch (e) {
    res.writeHead(404)
    res.end('Not found')
  }
}

function proxyToBackend(req, res) {
  let proxyPath = req.url
  // vue.config.js 里 pathRewrite 去掉 /api，这里也一样
  if (proxyPath.startsWith('/api')) proxyPath = proxyPath.replace(/^\/api/, '')
  const opts = {
    hostname: '127.0.0.1',
    port: 8081,
    path: req.url,
    method: req.method,
    path: proxyPath,
    headers: { ...req.headers, host: '127.0.0.1:8081' },
  }
  const proxy = http.request(opts, (pRes) => {
    res.writeHead(pRes.statusCode, pRes.headers)
    pRes.pipe(res)
  })
  proxy.on('error', () => { res.writeHead(502); res.end('Bad Gateway') })
  req.pipe(proxy)
}

const server = http.createServer((req, res) => {
  const parsed = url.parse(req.url)
  if (parsed.pathname.startsWith('/api')) {
    return proxyToBackend(req, res)
  }
  // SPA fallback: non-file requests -> index.html
  let filePath = path.join(DIST, parsed.pathname === '/' ? '/index.html' : parsed.pathname)
  if (!path.extname(filePath)) filePath = path.join(DIST, 'index.html')
  serveFile(res, filePath)
})

server.listen(PORT, () => console.log(`Proxy server on http://localhost:${PORT} -> ${BACKEND}`))
