#!/usr/bin/env python3
"""测试环境：静态托管 dist/ 并将 /api/* 及 SSE 等后端根路径反向代理到 ism_server。"""

from __future__ import annotations

import argparse
import http.client
import http.server
import json
import os
import socket
import sys
import urllib.parse
from pathlib import Path

# EventSource 等不走 /api 前缀的后端根路径
BACKEND_ROOT_PREFIXES = ("/SSEPushData",)

# 大组态 JSON 经 cpolar 外网可达 90s+；SSE 需长连接
PROXY_TIMEOUT = int(os.environ.get("ISM_PROXY_TIMEOUT", "300"))
SSE_TIMEOUT = int(os.environ.get("ISM_SSE_TIMEOUT", "3600"))
STREAM_CHUNK = 65536


def _default_fe_port() -> int:
    raw = os.environ.get("ISM_FE_PORT") or os.environ.get("ISM_FRONTEND_PORT") or "7080"
    return int(raw)


def _default_be_url() -> str:
    if url := os.environ.get("ISM_BE_URL"):
        return url
    be_port = os.environ.get("ISM_BE_PORT") or "8091"
    return f"http://127.0.0.1:{be_port}"


class ProxyHandler(http.server.SimpleHTTPRequestHandler):
    protocol_version = "HTTP/1.1"
    api_backend: str = _default_be_url()
    dist_dir: Path = Path(".")

    def __init__(self, *args, **kwargs):
        super().__init__(*args, directory=str(self.dist_dir), **kwargs)

    def _path_only(self) -> str:
        return self.path.split("?", 1)[0]

    def _needs_proxy(self) -> bool:
        path = self._path_only()
        if path.startswith("/api/") or path == "/api":
            return True
        return any(path.startswith(prefix) for prefix in BACKEND_ROOT_PREFIXES)

    def _is_sse_path(self) -> bool:
        return self._path_only().startswith("/SSEPushData")

    def do_GET(self):
        if self._needs_proxy():
            self._proxy_backend(stream_only=self._is_sse_path())
        else:
            super().do_GET()

    def do_POST(self):
        if self._needs_proxy():
            self._proxy_backend(stream_only=False)
        else:
            self.send_error(405, "Method Not Allowed")

    def do_PUT(self):
        if self._needs_proxy():
            self._proxy_backend(stream_only=False)
        else:
            self.send_error(405, "Method Not Allowed")

    def do_DELETE(self):
        if self._needs_proxy():
            self._proxy_backend(stream_only=False)
        else:
            self.send_error(405, "Method Not Allowed")

    def _proxy_headers(self, body_len: int) -> dict[str, str]:
        skip = {"host", "connection", "content-length", "proxy-connection", "keep-alive"}
        headers = {
            k: v for k, v in self.headers.items() if k.lower() not in skip
        }
        if body_len:
            headers["Content-Length"] = str(body_len)
        # 保留后端 gzip，避免代理层解压后体积膨胀（26MB gzip → 420MB 明文）
        if "Accept-Encoding" not in headers and "accept-encoding" not in headers:
            headers["Accept-Encoding"] = "gzip"
        return headers

    def _backend_target_path(self) -> str:
        target_path = self.path
        if target_path.startswith("/api"):
            target_path = target_path[4:] or "/"
        return target_path

    def _backend_host_port(self) -> tuple[str, int]:
        parsed = urllib.parse.urlparse(self.api_backend)
        host = parsed.hostname or "127.0.0.1"
        port = parsed.port or (443 if parsed.scheme == "https" else 80)
        return host, port

    def _proxy_backend(self, stream_only: bool):
        target_path = self._backend_target_path()
        if os.environ.get("ISM_TEST_SKIP_LICENSE", "1") == "1":
            stub_path = target_path.split("?", 1)[0].rstrip("/")
            if stub_path.endswith("GetPhysicalIDCheck"):
                payload = json.dumps({"code": 0, "id": "3B93406239702DCE"}).encode()
                self.send_response(200)
                self.send_header("Content-Type", "application/json")
                self.send_header("Content-Length", str(len(payload)))
                self.end_headers()
                self.wfile.write(payload)
                return

        host, port = self._backend_host_port()
        length = int(self.headers.get("Content-Length", 0))
        body = self.rfile.read(length) if length else None
        timeout = SSE_TIMEOUT if stream_only else PROXY_TIMEOUT
        conn = http.client.HTTPConnection(host, port, timeout=timeout)
        try:
            conn.request(
                self.command,
                target_path,
                body=body,
                headers=self._proxy_headers(length),
            )
            resp = conn.getresponse()
            self.send_response(resp.status)
            skip_resp = {"transfer-encoding", "connection"}
            for k, v in resp.getheaders():
                if k.lower() not in skip_resp:
                    self.send_header(k, v)
            if stream_only:
                # cpolar / nginx 反代：禁用缓冲，避免 SSE 502
                self.send_header("Cache-Control", "no-cache, no-transform")
                self.send_header("Connection", "keep-alive")
                self.send_header("X-Accel-Buffering", "no")
            self.end_headers()
            while True:
                chunk = resp.read(STREAM_CHUNK)
                if not chunk:
                    break
                self.wfile.write(chunk)
                self.wfile.flush()
        except (BrokenPipeError, ConnectionResetError, socket.timeout) as e:
            if stream_only:
                sys.stderr.write(f"SSE client disconnected: {e}\n")
            else:
                self.send_error(502, f"Backend proxy error: {e}")
        except Exception as e:
            self.send_error(502, f"Backend proxy error: {e}")
        finally:
            conn.close()

    def log_message(self, fmt, *args):
        sys.stdout.write("%s - %s\n" % (self.log_date_time_string(), fmt % args))


def main():
    parser = argparse.ArgumentParser(description="ISM 测试环境前端静态服务 + API 代理")
    parser.add_argument("--port", type=int, default=_default_fe_port())
    parser.add_argument("--dist", type=Path, default=Path("web/dist"))
    parser.add_argument("--backend", default=_default_be_url())
    args = parser.parse_args()
    dist = args.dist.resolve()
    if not dist.is_dir():
        print(f"错误: dist 目录不存在: {dist}", file=sys.stderr)
        sys.exit(1)
    ProxyHandler.dist_dir = dist
    ProxyHandler.api_backend = args.backend
    os.chdir(dist)
    server = http.server.ThreadingHTTPServer(("0.0.0.0", args.port), ProxyHandler)
    print(f"前端: http://0.0.0.0:{args.port}/  (dist={dist})")
    print(f"API 代理: /api/* -> {args.backend} (timeout={PROXY_TIMEOUT}s, stream chunk={STREAM_CHUNK})")
    print(f"根路径代理: {', '.join(BACKEND_ROOT_PREFIXES)} -> {args.backend} (timeout={SSE_TIMEOUT}s)")
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        print("\n已停止")


if __name__ == "__main__":
    main()
