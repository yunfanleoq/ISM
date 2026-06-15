#!/usr/bin/env python3
"""
EMQX v5 API Mock 服务
模拟 EMQX Dashboard API，供 ISM 后端本地开发使用。
监听 18083 端口，处理 /api/v5/clients 请求。
"""

import json
import base64
from http.server import HTTPServer, BaseHTTPRequestHandler

# 与 mqtt.conf 中配置一致
EXPECTED_API_KEY = "23a808a53adc53eb"
EXPECTED_SECRET_KEY = "ODjOMy2SidBDYbwj2hVSZd5yFQksNQjEwS7LZYUPeXD"
EXPECTED_AUTH = base64.b64encode(
    f"{EXPECTED_API_KEY}:{EXPECTED_SECRET_KEY}".encode()
).decode()


class EmqxMockHandler(BaseHTTPRequestHandler):
    """EMQX v5 API Mock 请求处理器"""

    def _check_auth(self):
        """验证 Basic Auth"""
        auth_header = self.headers.get("Authorization", "")
        if not auth_header.startswith("Basic "):
            return False
        token = auth_header[6:]
        return token == EXPECTED_AUTH

    def _send_json(self, code, data):
        """发送 JSON 响应"""
        body = json.dumps(data, ensure_ascii=False).encode("utf-8")
        self.send_response(code)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    def log_message(self, format, *args):
        """自定义日志格式"""
        print(f"[EMQX Mock] {self.client_address[0]} - {format % args}")

    def do_GET(self):
        # 验证认证
        if not self._check_auth():
            self._send_json(401, {"code": "UNAUTHORIZED", "message": "Auth failed"})
            return

        # /api/v5/clients?page=1&size=2000
        if self.path.startswith("/api/v5/clients"):
            self._send_json(200, {
                "data": [],
                "meta": {"page": 1, "limit": 2000, "count": 0}
            })
            return

        # 未匹配的路由
        self._send_json(404, {"code": "NOT_FOUND", "message": f"Unknown path: {self.path}"})

    def do_POST(self):
        self._check_auth()
        self._send_json(200, {"data": {}, "meta": {}})


def main():
    host = "127.0.0.1"
    port = 18083

    server = HTTPServer((host, port), EmqxMockHandler)
    print(f"[EMQX Mock] 服务已启动 → http://{host}:{port}")
    print(f"[EMQX Mock] API:     http://{host}:{port}/api/v5/clients")
    print(f"[EMQX Mock] Auth:    Basic {EXPECTED_API_KEY}:{EXPECTED_SECRET_KEY[:8]}...")
    print(f"[EMQX Mock] 按 Ctrl+C 停止服务")

    try:
        server.serve_forever()
    except KeyboardInterrupt:
        print("\n[EMQX Mock] 服务已停止")
        server.shutdown()


if __name__ == "__main__":
    main()
