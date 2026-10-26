#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
ISM Modbus TCP 模拟器 v9 — DB 模型驱动（航信机房）
================================================================
为什么是 v9（彻底重写 v8）：
  v8 把寄存器布局硬编码成 Float（A20/A40）+ 固定地址，但 OceanBase 里的
  `modbus_devices_data_model` 实际是 **Unsigned short / Long / Bool**，且地址、
  系数(conversion_expression)、字节序(byte_order / data_format) 都不一样 →
  v8 产的原始 uint16 既没按系数换算、又没按后端的字节序编码，导致大屏数值离谱
  （UPS 输出电压 44296V、有功为负、电池剩余 32512min；标准电表同样）。

第一性原理：后端 modbusPthread.go 的解码链是确定的，模拟器只要 **反演** 它即可：
  最终显示值 final = decode(raw_registers) * K        （K 来自 conversion_expression `{val}*K`）
  - Short/Unsigned short(1寄存器)：device.data_format ∉ {BigEndian,ABCD} ⇒ 后端按
    LittleEndian 读 ⇒ decode(raw)=byteswap16(raw)。故 raw = byteswap16(D)，D=round(物理量/K)。
  - Long/Float(2寄存器, byte_order=CDAB)：decode = (reg[A+1]<<16)|reg[A]。
    故 reg[A]=D&0xFFFF（低字）、reg[A+1]=(D>>16)&0xFFFF（高字）。
  - Bool：经离散量(function 2/DI) 读取 → di_normal()（状态/合闸/通信=1，故障/报警=0）。

因此本模拟器从 DB 读每个数据点的 (address, type, byte_order, conversion_expression)
与设备的 data_format/slave，按 **物理量→系数反演→字节序编码** 产值（DRY：单位/系数字段驱动，
不逐点硬编码）。每个点落到合理量程（电压~220/380V、电流合理安培、功率为正、功率因数0~1、
电池剩余分钟/温度合理）。重启只影响模拟数据，不动后端/前端/DB。

用法:
  python3 scripts/modbus_simulator.py                 # 连 OceanBase 读模型, 监听 0.0.0.0:502
  python3 scripts/modbus_simulator.py --selfcheck     # 不起服务, 打印反演解码后的样例值核对量程
"""
import math, struct, time, threading, socketserver, signal, json, argparse, re, sys
from http.server import HTTPServer, BaseHTTPRequestHandler
from collections import defaultdict

import pymysql

# ── DB（与 build_ncc_dashboard.py 同源，OceanBase dbtype=4）──
OB = dict(host='127.0.0.1', port=2881, user='root@ism_tenant',
          password='ism2024!', database='ism')
PROJECT_UUID = '31bc90be-ebc4-dd61-ba9d-ce6e075e40e2'

MAX_REGS = 256
MBAP_LEN = 7
LOCK = threading.Lock()

# slave(int) -> {
#   'name': str, 'data_format': str,
#   'hr_points': [(addr, name, typ, K, byte_order)],   # 模拟量 (function 3/4)
#   'di_points': [(addr, bit)],                          # 离散量 (function 2)
#   'hr': [int]*MAX_REGS, 'di': {addr:0/1},
# }
SLAVES = {}


# ─────────────────────── 编解码反演 ───────────────────────

def byteswap16(v):
    v &= 0xFFFF
    return ((v & 0xFF) << 8) | (v >> 8)


def parse_factor(expr):
    """从 conversion_expression `{val}*K` 解析系数 K，默认 1.0。"""
    if not expr:
        return 1.0
    m = re.search(r'\*\s*([0-9]*\.?[0-9]+)', expr)
    if m:
        try:
            return float(m.group(1))
        except ValueError:
            return 1.0
    return 1.0


def short_is_byteswap(data_format):
    """单寄存器 Short/Unsigned short：后端 data_format∉{BigEndian,ABCD} 时按 LittleEndian 读。"""
    return data_format not in ('BigEndian', 'ABCD')


# ─────────────────────── 物理量基准（DRY：按名/单位） ───────────────────────

def phys_target(name, unit):
    """返回 (基准值, 噪声幅度)。落到该物理量的合理量程；噪声=0 表示状态/码值恒定。"""
    n = name
    if '畸变' in n:                       # 谐波畸变率(含"电流/电压"字样)须先于电压/电流判断
        return 3.0, 0.5
    if '线电压' in n:
        return 380.0, 2.0
    if '相电压' in n:
        return 220.0, 1.5
    if '电池电压' in n:
        return 540.0, 3.0
    if '电压' in n:
        return 220.0, 1.5
    if '电流' in n:
        if '中性' in n or '零序' in n:
            return 0.8, 0.1
        if '旁路' in n:
            return 6.0, 0.6
        if '电池' in n:
            return 30.0, 3.0
        if '输出' in n or '主路' in n:
            return 90.0, 5.0
        return 8.0, 0.4
    if '频率' in n:
        return 50.0, 0.05
    if '功率因数' in n:
        # 后端对 Short/Long 类型把换算结果强制 int32 截断 → 0.95 会变 0。
        # 目标取单位功率因数 1.0（合理、落在 0~1、截断后显示 "1" 而非误导性的 "0"）。
        return 1.0, 0.0
    if ('有功' in n or '无功' in n or '视在' in n) and '电度' not in n and '电能' not in n:
        if '无功' in n:
            return 18.0, 1.0
        if '视在' in n:
            return 60.0, 2.0
        return 58.0, 2.0          # 有功功率：正值
    if '电度' in n or '电能' in n:
        return 12000.0, 2.0
    if '畸变' in n:
        return 3.0, 0.5
    if '温度' in n:
        return 30.0, 1.5
    if '湿度' in n:
        return 45.0, 3.0
    if '剩余' in n and ('时间' in n or '运行' in n):
        return 120.0, 5.0
    if '使用模式' in n or '模式' in n:
        return 2.0, 0.0           # 运行模式码（2=正常逆变示意）
    if '状态' in n or '故障' in n or '错误' in n or '充电' in n:
        return 1.0, 0.0           # 数值化状态（DI 才是真正的开关量）
    if unit == 'V':
        return 220.0, 1.5
    if unit == 'A':
        return 8.0, 0.4
    if unit in ('kW', 'kvar', 'kVA'):
        return 50.0, 2.0
    if unit == 'Hz':
        return 50.0, 0.05
    if unit == 'min':
        return 120.0, 5.0
    if unit in ('℃', '°C'):
        return 30.0, 1.5
    return 1.0, 0.0


def di_normal(name):
    bad = ('故障', '报警', '告警', '欠压', '过载', '异常', '硬件', '错误')
    good = ('合闸', '通信', '通讯', '正常', '状态', '运行', '开关', '充电')
    if any(k in name for k in bad):
        return 0
    if any(k in name for k in good):
        return 1
    return 0


def gen_phys(base, noise, t, phase):
    """围绕 base 缓慢起伏（功率为正等由 base 决定，这里不强制非负，调用方按需 clamp）。"""
    if noise == 0:
        return base
    return base + noise * (math.sin(t * 0.4 + phase) * 0.6 +
                           math.sin(t * 1.3 + phase * 1.7) * 0.3 +
                           math.sin(t * 2.7 + phase * 2.3) * 0.1)


# ─────────────────────── 模型加载（DB 驱动） ───────────────────────

def load_model():
    conn = pymysql.connect(**OB)
    cur = conn.cursor()
    # 设备：本项目 type=1，取 slave / muid / data_format
    cur.execute("""
        SELECT m.name, m.muid, m.extra_data, dm.data_format
        FROM monitor_list m
        LEFT JOIN devices_model dm ON m.muid = dm.uuid
        WHERE m.project_uuid = %s AND m.type = 1 AND m.deleted_at IS NULL
    """, (PROJECT_UUID,))
    devices = cur.fetchall()

    # 每个 muid 的点表（缓存，避免重复查询）
    points_cache = {}

    def points_for(muid):
        if muid in points_cache:
            return points_cache[muid]
        cur.execute("""
            SELECT register_address, name, type, byte_order, conversion_expression
            FROM modbus_devices_data_model WHERE muid = %s ORDER BY register_address
        """, (muid,))
        pts = cur.fetchall()
        points_cache[muid] = pts
        return pts

    for name, muid, extra, data_format in devices:
        try:
            slave = int(json.loads(extra or '{}')['modbus']['address'])
        except (KeyError, ValueError, TypeError):
            continue
        df = data_format or 'CDAB'
        hr_points, di_points = [], []
        for addr, pname, ptype, bo, conv in points_for(muid):
            if ptype == 'Bool':
                di_points.append((addr, di_normal(pname)))
            else:
                hr_points.append((addr, pname, ptype, parse_factor(conv), bo or df))
        st = {'name': name, 'data_format': df,
              'hr_points': hr_points, 'di_points': di_points,
              'hr': [0] * MAX_REGS, 'di': {a: b for a, b in di_points}}
        SLAVES[slave] = st
    conn.close()


def encode_point(st, addr, name, typ, K, byte_order, t):
    """把一个数据点的物理量编码进 st['hr']（反演后端解码）。"""
    base, noise = phys_target(name, '')
    phase = (hash(name + str(addr)) % 628) / 100.0
    phys = gen_phys(base, noise, t, phase)
    if K == 0:
        K = 1.0
    D = phys / K                                   # 后端 getValue 应得到的整数
    hr = st['hr']
    if typ in ('Long', 'Float'):
        if typ == 'Float':
            bits = struct.unpack('>I', struct.pack('>f', float(phys / K)))[0]
            D = bits
        Di = int(round(D)) & 0xFFFFFFFF
        # CDAB: 低字在前(reg A), 高字在后(reg A+1)
        if addr + 1 < len(hr):
            hr[addr] = Di & 0xFFFF
            hr[addr + 1] = (Di >> 16) & 0xFFFF
    else:  # Short / Unsigned short（单寄存器）
        Di = int(round(D))
        if typ == 'Short':
            Di = max(-32768, min(32767, Di)) & 0xFFFF
        else:
            Di = max(0, min(65535, Di))
        hr[addr] = byteswap16(Di) if short_is_byteswap(st['data_format']) else Di


def update_all():
    t = time.monotonic()
    with LOCK:
        for st in SLAVES.values():
            for addr, name, typ, K, bo in st['hr_points']:
                encode_point(st, addr, name, typ, K, bo, t)


# ─────────────────────── Modbus TCP 服务 ───────────────────────

def read_hr(slave, start, count):
    st = SLAVES.get(slave)
    if not st:
        return None
    with LOCK:
        hr = st['hr']
        end = min(start + count, len(hr))
        regs = list(hr[start:end])
    regs += [0] * (count - len(regs))
    return regs


def read_di(slave, start, count):
    st = SLAVES.get(slave)
    if not st:
        return None
    with LOCK:
        return [st['di'].get(start + i, 0) for i in range(count)]


class ModbusHandler(socketserver.StreamRequestHandler):
    def handle(self):
        try:
            while True:
                hdr = self.rfile.read(MBAP_LEN)
                if len(hdr) < MBAP_LEN:
                    break
                tid = (hdr[0] << 8) | hdr[1]
                pid = (hdr[2] << 8) | hdr[3]
                length = (hdr[4] << 8) | hdr[5]
                uid = hdr[6]
                pdu_len = length - 1
                if pdu_len <= 0 or pdu_len > 253:
                    break
                pdu = self.rfile.read(pdu_len)
                if len(pdu) < pdu_len:
                    break
                fc = pdu[0]
                resp = b''
                if uid not in SLAVES:
                    resp = bytes([fc | 0x80, 0x0A])
                elif fc in (3, 4) and len(pdu) >= 5:
                    start = (pdu[1] << 8) | pdu[2]
                    cnt = (pdu[3] << 8) | pdu[4]
                    regs = read_hr(uid, start, cnt) if 1 <= cnt <= 125 else None
                    if regs is None:
                        resp = bytes([fc | 0x80, 0x03])
                    else:
                        resp = bytes([fc, cnt * 2]) + b''.join(struct.pack('>H', r & 0xFFFF) for r in regs)
                elif fc in (1, 2) and len(pdu) >= 5:
                    start = (pdu[1] << 8) | pdu[2]
                    cnt = (pdu[3] << 8) | pdu[4]
                    vals = read_di(uid, start, cnt) if 1 <= cnt <= 2000 else None
                    if vals is None:
                        resp = bytes([fc | 0x80, 0x03])
                    else:
                        bc = (cnt + 7) // 8
                        bits = bytearray(bc)
                        for i, v in enumerate(vals):
                            if v:
                                bits[i // 8] |= 1 << (i % 8)
                        resp = bytes([fc, bc]) + bytes(bits)
                elif fc in (5, 6) and len(pdu) >= 5:
                    resp = pdu                       # echo write-single
                elif fc in (15, 16) and len(pdu) >= 6:
                    resp = pdu[:5]                   # echo write-multiple
                else:
                    resp = bytes([fc | 0x80, 0x01])
                self.wfile.write(struct.pack('>HHHB', tid, pid, 1 + len(resp), uid) + resp)
                self.wfile.flush()
        except Exception:
            pass


class ThreadedServer(socketserver.ThreadingMixIn, socketserver.TCPServer):
    allow_reuse_address = True
    daemon_threads = True
    request_queue_size = 256


# ─────────────────────── 轻量 HTTP 状态 API（兼容 v8 监控面板） ───────────────────────

class APIHandler(BaseHTTPRequestHandler):
    def log_message(self, *a):
        pass

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
            self._send(200, {'total_slaves': len(SLAVES),
                             'register_count': sum(len(s['hr_points']) for s in SLAVES.values())})
        elif path.startswith('/api/slave/'):
            try:
                sid = int(path.split('/')[-1])
            except ValueError:
                return self._send(400, {'error': 'invalid sid'})
            st = SLAVES.get(sid)
            if not st:
                return self._send(404, {'error': 'not found'})
            self._send(200, {'slave': sid, 'name': st['name'],
                             'holding_registers': {a: {'name': n} for a, n, *_ in st['hr_points']}})
        else:
            self._send(200, {'status': 'ok', 'version': 'v9-db-driven'})


def start_http(port=5040):
    HTTPServer(('0.0.0.0', port), APIHandler).serve_forever()


def selfcheck():
    """不起服务：反演解码后打印样例点，核对落在合理量程（与后端解码逻辑一致）。"""
    update_all()

    def decode(st, addr, typ, K, bo):
        hr = st['hr']
        if typ in ('Long', 'Float'):
            lo, hi = hr[addr], hr[addr + 1]
            bits = (hi << 16) | lo
            if typ == 'Float':
                return round(struct.unpack('>f', struct.pack('>I', bits))[0] * K, 3)
            if bits >= 0x80000000:
                bits -= 0x100000000
            return round(bits * K, 3)
        raw = hr[addr]
        g = byteswap16(raw) if short_is_byteswap(st['data_format']) else raw
        if typ == 'Short' and g >= 0x8000:
            g -= 0x10000
        return round(g * K, 3)

    for slave in sorted(SLAVES)[:1] + sorted(SLAVES)[-1:]:
        st = SLAVES[slave]
        print(f"\n[slave {slave}] {st['name']} data_format={st['data_format']}")
        for addr, name, typ, K, bo in st['hr_points'][:60]:
            print(f"   {name:24s} type={typ:14s} K={K:<6} → decode≈{decode(st, addr, typ, K, bo)}")
        if st['di_points']:
            print("   [DI]", ', '.join(f"{a}={b}" for a, b in st['di_points']))


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument('--host', default='0.0.0.0')
    ap.add_argument('--port', type=int, default=502)
    ap.add_argument('--http-port', type=int, default=5040)
    ap.add_argument('--interval', type=float, default=3.0)
    ap.add_argument('--selfcheck', action='store_true')
    args = ap.parse_args()

    print('ISM Modbus TCP Simulator v9 (DB-driven)', flush=True)
    load_model()
    total_pts = sum(len(s['hr_points']) for s in SLAVES.values())
    print(f'  从 OceanBase 加载 slaves={len(SLAVES)}  模拟量点={total_pts}', flush=True)
    if not SLAVES:
        print('  [!] 未加载到任何设备，检查 PROJECT_UUID / DB 连接', flush=True)
        sys.exit(1)

    if args.selfcheck:
        selfcheck()
        return

    def updater():
        while True:
            update_all()
            time.sleep(args.interval)

    update_all()
    threading.Thread(target=updater, daemon=True).start()
    threading.Thread(target=lambda: start_http(args.http_port), daemon=True).start()

    srv = ThreadedServer((args.host, args.port), ModbusHandler)
    print(f'  Modbus TCP on {args.host}:{args.port}  HTTP on :{args.http_port}', flush=True)
    print('  [Running, Ctrl+C to stop]', flush=True)
    try:
        srv.serve_forever()
    except KeyboardInterrupt:
        print('\n  stopped')


if __name__ == '__main__':
    try:
        signal.signal(signal.SIGPIPE, signal.SIG_IGN)
    except Exception:
        pass
    main()
