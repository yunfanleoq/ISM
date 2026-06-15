#!/usr/bin/env python3
"""
ISM Modbus TCP 模拟器 v8 — 带轻量 HTTP 状态 API
Modbus TCP (502) + HTTP API (5040) 用于前端监控面板。
每个连接独立线程，全局锁保护寄存器数据。
"""
import math, struct, time, sys, threading, socketserver, traceback, signal, json
from http.server import HTTPServer, BaseHTTPRequestHandler

MAX_REGS = 200
MBAP_LEN = 7


def f2cdab(v):
    p = struct.pack('<f', v)
    return (p[1] << 8) | p[0], (p[3] << 8) | p[2]


def gen(base, noise, t, phase):
    return max(0.0, base + noise * (
        math.sin(t * 0.4 + phase) * 0.6 +
        math.sin(t * 1.3 + phase * 1.7) * 0.3 +
        math.sin(t * 2.7 + phase * 2.3) * 0.1))


# ─── 寄存器定义（硬编码，彻底避开 DB 依赖）─────────────

A20_REGS = [
    (0, "AB线电压", 380.0, 0.8), (2, "BC线电压", 380.0, 0.8),
    (4, "CA线电压", 379.8, 0.8), (6, "频率", 50.0, 0.05),
    (8, "A相电流", 8.0, 0.15), (10, "B相电流", 7.9, 0.15),
    (12, "C相电流", 8.2, 0.15), (14, "中性线电流", 0.62, 0.08),
    (16, "总有功功率", 51.8, 1.5), (18, "总无功功率", 16.2, 0.8),
    (20, "总视在功率", 54.3, 1.5), (22, "总功率因数", 0.952, 0.008),
    (24, "正有功电度", 1285.0, 0.01), (26, "A相电流谐波畸变率", 2.8, 0.4)]

A40_REGS = [
    (0, "AB线电压", 382.0, 0.8), (2, "BC线电压", 381.5, 0.8),
    (4, "CA线电压", 380.2, 0.8), (6, "频率", 50.01, 0.05),
    (8, "A相电流", 8.95, 0.18), (10, "B相电流", 8.72, 0.18),
    (12, "C相电流", 9.18, 0.18), (14, "中性线电流", 0.71, 0.08),
    (16, "总有功功率", 56.3, 1.6), (18, "总无功功率", 18.5, 0.9),
    (20, "总视在功率", 59.2, 1.6), (22, "总功率因数", 0.928, 0.008),
    (24, "A相电压", 221.0, 2.0), (26, "B相电压", 219.8, 2.0),
    (28, "C相电压", 220.5, 2.0), (30, "A相电流谐波畸变率", 3.5, 0.5)]

# UPS Short 模式：单寄存器 Short 值，匹配 UPS 数据模型 (muid=3e665b0e)
UPS_REGS = [
    (0, "UPS使用模式",       1, 0),    (1, "UPS逆变器内部故障",  0, 0),
    (2, "UPS电池开关状态",    1, 0),    (3, "UPS主路开关状态",   1, 0),
    (4, "UPS旁路状态",       0, 0),    (5, "UPS旁路内部硬件错误",0, 0),
    (6, "UPS充电状态",       1, 0),
    (7, "主路A相电流",    315, 25),   (8,  "主路B相电流",  298, 25),
    (9, "主路C相电流",    332, 25),   (10, "旁路A相电流",    5,  1),
    (11,"旁路B相电流",      4,  1),   (12, "旁路C相电流",    6,  1),
    (13,"输出A相电流",    315, 25),   (14, "输出B相电流",  298, 25),
    (15,"输出C相电流",    332, 25),   (16, "电池电流",     108, 20),
    (17,"主路输入AB线电压",2215,30),   (18, "主路输入BC线电压",2208,30),
    (19,"主路输入CA线电压",2195,30),   (20, "旁路输入A相电压",2200,15),
    (21,"旁路输入B相电压",2201,15),   (22, "旁路输入C相电压",2199,15),
    (23,"旁路输入AB线电压",2215,30),   (24, "旁路输入BC线电压",2208,30),
    (25,"旁路输入CA线电压",2195,30),   (26, "输出A相电压",  2200,15),
    (27,"输出B相电压",    2201,15),   (28, "输出C相电压",  2199,15),
    (29,"输出AB线电压",   2215,30),   (30, "输出BC线电压",  2208,30),
    (31,"输出CA线电压",   2195,30),   (32, "电池电压",      541, 5),
    (33,"输出总有功功率",  208,10),   (34, "输出视在功率",   210,10),
    (35,"输出功率因数",     92, 5),   (36, "主路输入频率",  5000,20),
    (37,"旁路输入频率",    5000,20),   (38, "输出频率",     5000,20),
    (39,"电池剩余运行时间", 125, 5),
]


# ─── 全局数据（线程安全）───────────────────────────────

LOCK = threading.Lock()
HOLDING = {}   # {sid: [reg0, reg1, ...]}
REG_MAPS = {}  # {sid: (regs, mode)}  mode='float' or 'short'


def init_data():
    global HOLDING, REG_MAPS
    with LOCK:
        for sid in range(1, 77):
            if sid <= 60:
                regs, mode = A20_REGS[:], 'float'
            elif sid <= 69:
                regs, mode = A40_REGS[:], 'float'
            else:
                regs, mode = UPS_REGS[:], 'short'
            REG_MAPS[sid] = (regs, mode)
            hr = [0] * MAX_REGS
            if mode == 'float':
                for a, n, b, _ in regs:
                    lo, hi = f2cdab(b)
                    hr[a] = lo
                    hr[a + 1] = hi
            else:  # short mode
                for a, n, b, _ in regs:
                    hr[a] = int(b)
            HOLDING[sid] = hr


def update_data():
    t = time.monotonic()
    with LOCK:
        for sid, (regs, mode) in REG_MAPS.items():
            hr = HOLDING[sid]
            if mode == 'float':
                for a, n, b, noise in regs:
                    phase = hash(n + str(sid)) % 628 / 100.0
                    val = gen(b, noise, t, phase)
                    lo, hi = f2cdab(val)
                    hr[a] = lo
                    hr[a + 1] = hi
            else:  # short mode
                for a, n, b, noise in regs:
                    if noise == 0:
                        hr[a] = int(b)
                    else:
                        phase = hash(n + str(sid)) % 628 / 100.0
                        val = max(0, int(b + noise * (
                            math.sin(t * 0.4 + phase) * 0.6 +
                            math.sin(t * 1.3 + phase * 1.7) * 0.3 +
                            math.sin(t * 2.7 + phase * 2.3) * 0.1)))
                        hr[a] = val


def read_hr(sid, start, count):
    with LOCK:
        hr = HOLDING.get(sid, [0] * MAX_REGS)
        end = min(start + count, len(hr))
        return list(hr[start:end])


def read_di(sid, start, count):
    return [1 if i in (0, 2) else 0 for i in range(start, min(start + count, 50))]


# ─── Modbus 连接处理 ──────────────────────────────────


class ModbusHandler(socketserver.StreamRequestHandler):

    def handle(self):
        addr = self.client_address
        print(f"  [+] client connected: {addr}", flush=True)
        try:
            while True:
                # 读 MBAP 头 (7 字节)
                hdr = self.rfile.read(MBAP_LEN)
                if len(hdr) < MBAP_LEN:
                    break
                tid = (hdr[0] << 8) | hdr[1]
                pid = (hdr[2] << 8) | hdr[3]
                length = (hdr[4] << 8) | hdr[5]
                uid = hdr[6]
                pdu_len = length - 1
                if pdu_len <= 0 or pdu_len > 253:
                    print(f"  [!] {addr} bad pdu_len={pdu_len}", flush=True)
                    break

                pdu = self.rfile.read(pdu_len)
                if len(pdu) < pdu_len:
                    print(f"  [!] {addr} short pdu", flush=True)
                    break

                fc = pdu[0]
                resp_pdu = b''

                # Debug: log incoming request
                if fc in (1,2,3,4):
                    start_addr = (pdu[1] << 8) | pdu[2]
                    cnt = (pdu[3] << 8) | pdu[4]
                    print(f"  [>] req slave={uid} fc={fc} addr={start_addr} cnt={cnt}", flush=True)

                if uid not in REG_MAPS:
                    resp_pdu = bytes([fc | 0x80, 0x0A])
                elif fc == 3 and len(pdu) >= 5:  # Read HR
                    start = (pdu[1] << 8) | pdu[2]
                    count = (pdu[3] << 8) | pdu[4]
                    if 1 <= count <= 125:
                        regs = read_hr(uid, start, count)
                        data = bytes([count * 2]) + b''.join(
                            struct.pack('>H', r) for r in regs)
                        resp_pdu = bytes([3]) + data
                    else:
                        resp_pdu = bytes([0x83, 0x03])
                elif fc == 4 and len(pdu) >= 5:  # Read Input Registers (same as HR)
                    start = (pdu[1] << 8) | pdu[2]
                    count = (pdu[3] << 8) | pdu[4]
                    if 1 <= count <= 125:
                        regs = read_hr(uid, start, count)
                        data = bytes([count * 2]) + b''.join(
                            struct.pack('>H', r) for r in regs)
                        resp_pdu = bytes([4]) + data
                    else:
                        resp_pdu = bytes([0x84, 0x03])
                elif fc == 1 and len(pdu) >= 5:  # Read Coils
                    start = (pdu[1] << 8) | pdu[2]
                    count = (pdu[3] << 8) | pdu[4]
                    if 1 <= count <= 2000:
                        vals = read_di(uid, start, count)
                        bc = (len(vals) + 7) // 8
                        bits = bytearray(bc)
                        for i, v in enumerate(vals):
                            if v:
                                bits[i // 8] |= 1 << (i % 8)
                        resp_pdu = bytes([1]) + bytes([bc]) + bytes(bits)
                    else:
                        resp_pdu = bytes([0x81, 0x03])
                elif fc == 2 and len(pdu) >= 5:  # Read DI
                    start = (pdu[1] << 8) | pdu[2]
                    count = (pdu[3] << 8) | pdu[4]
                    if 1 <= count <= 2000:
                        vals = read_di(uid, start, count)
                        bc = (len(vals) + 7) // 8
                        bits = bytearray(bc)
                        for i, v in enumerate(vals):
                            if v:
                                bits[i // 8] |= 1 << (i % 8)
                        resp_pdu = bytes([2]) + bytes([bc]) + bytes(bits)
                    else:
                        resp_pdu = bytes([0x82, 0x03])
                else:
                    resp_pdu = bytes([fc | 0x80, 0x01])

                # --- Write functions (accept silently, no data change) ---
                if fc == 5 and len(pdu) >= 5:  # Write Single Coil
                    resp_pdu = pdu  # echo back
                elif fc == 6 and len(pdu) >= 5:  # Write Single Register
                    resp_pdu = pdu  # echo back
                elif fc == 15 and len(pdu) >= 6:  # Write Multiple Coils
                    resp_pdu = pdu[:4]  # echo address + quantity only
                elif fc == 16 and len(pdu) >= 6:  # Write Multiple Registers
                    resp_pdu = pdu[:4]  # echo address + quantity only

                # 构造 MBAP 响应头
                resp_hdr = struct.pack('>HHHB', tid, pid, 1 + len(resp_pdu), uid)
                self.wfile.write(resp_hdr + resp_pdu)
                self.wfile.flush()
        except Exception as e:
            print(f"  [-] {addr} error: {e}", flush=True)
        finally:
            print(f"  [-] client disconnected: {addr}", flush=True)


class ThreadedModbusServer(socketserver.ThreadingMixIn, socketserver.TCPServer):
    allow_reuse_address = True
    daemon_threads = True
    request_queue_size = 256


# ─── HTTP API（轻量级，用于前端监控面板）─────────────────


class SimAPIHandler(BaseHTTPRequestHandler):

    def log_message(self, fmt, *args):
        pass  # 静默日志

    def _send(self, code, data):
        body = json.dumps(data, ensure_ascii=False).encode()
        self.send_response(code)
        self.send_header('Content-Type', 'application/json; charset=utf-8')
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Content-Length', str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    def do_GET(self):
        path = self.path.rstrip('/')
        if path == '/api/summary':
            self._summary()
        elif path == '/api/slaves':
            self._slaves()
        elif path.startswith('/api/slave/'):
            sid = path.split('/')[-1]
            self._slave_detail(sid)
        else:
            self._send(200, {'status': 'ok', 'version': 'v8'})

    def _summary(self):
        self._send(200, {
            'total_slaves': 76,
            'a20_range': [1, 60],
            'a40_range': [61, 69],
            'ups_range': [70, 76],
            'register_count': sum(len(regs) for regs, _ in REG_MAPS.values()),
        })

    def _slaves(self):
        result = []
        for sid in sorted(REG_MAPS.keys()):
            regs, mode = REG_MAPS[sid]
            hr = HOLDING.get(sid, [])
            reg_dict = {}
            if mode == 'float':
                for addr, name, _, _ in regs:
                    lo, hi = hr[addr], hr[addr + 1]
                    try:
                        val = struct.unpack('<f', struct.pack('<HH', lo, hi))[0]
                    except:
                        val = 0
                    reg_dict[addr] = {'name': name, 'value': round(val, 3)}
            else:  # short mode
                for addr, name, _, _ in regs:
                    reg_dict[addr] = {'name': name, 'value': hr[addr]}
            result.append({
                'slave': sid,
                'device_type': 'A20' if sid <= 60 else ('A40' if sid <= 69 else 'UPS'),
                'holding_registers': reg_dict,
                'discrete_inputs': {0: 1, 2: 1},
            })
        self._send(200, result)

    def _slave_detail(self, sid_str):
        try:
            sid = int(sid_str)
        except ValueError:
            self._send(400, {'error': 'invalid sid'})
            return
        if sid not in REG_MAPS:
            self._send(404, {'error': f'slave {sid} not found'})
            return
        regs, mode = REG_MAPS[sid]
        hr = HOLDING.get(sid, [])
        registers = {}
        if mode == 'float':
            for addr, name, base, noise in regs:
                lo, hi = hr[addr], hr[addr + 1]
                try:
                    val = struct.unpack('<f', struct.pack('<HH', lo, hi))[0]
                except:
                    val = 0
                registers[addr] = {'name': name, 'value': round(val, 3), 'raw': (hi << 8) | lo}
        else:  # short mode
            for addr, name, base, noise in regs:
                registers[addr] = {'name': name, 'value': hr[addr]}
        self._send(200, {
            'slave': sid,
            'device_type': 'A20' if sid <= 60 else ('A40' if sid <= 69 else 'UPS'),
            'holding_registers': registers,
            'discrete_inputs': {i: 1 if i in (0, 2) else 0 for i in range(4)},
        })


def start_http():
    httpd = HTTPServer(('0.0.0.0', 5040), SimAPIHandler)
    print(f"  HTTP API on 0.0.0.0:5040", flush=True)
    httpd.serve_forever()


# ─── 主程序 ───────────────────────────────────────────


def main():
    print("ISM Modbus TCP Simulator v8", flush=True)
    init_data()
    total = sum(len(regs) for regs, _ in REG_MAPS.values())
    print(f"  Slaves: 76  Points: {total}  Port: 0.0.0.0:502", flush=True)

    # 数据更新线程
    def updater():
        while True:
            update_data()
            time.sleep(5)

    threading.Thread(target=updater, daemon=True).start()

    # HTTP API 线程
    threading.Thread(target=start_http, daemon=True).start()

    # Modbus TCP 服务器
    server = ThreadedModbusServer(("0.0.0.0", 502), ModbusHandler)
    print("  [Running, Ctrl+C to stop]", flush=True)

    try:
        server.serve_forever()
    except KeyboardInterrupt:
        print("\n  Stopped", flush=True)


if __name__ == "__main__":
    # 忽略 SIGPIPE，防止写入已关闭连接时崩溃
    try:
        signal.signal(signal.SIGPIPE, signal.SIG_IGN)
    except:
        pass
    main()
