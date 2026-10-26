#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""中航信数据中心 —— 包驱动 Modbus TCP 模拟器。

从契约 JSON 自动构建每台设备(=一个 slave id)的寄存器映射，按真实管理机 IP 分别监听:
  - 真实IP模式(默认): 为每个网关 IP 起一个 :502 服务(需 root + 回环别名)
  - 测试模式(--test): 仅在 127.0.0.1:<port> 起单服务(无需 root)，用于解码自检

编码与后端 modbusPthread.go 解码严格对齐:
  - 单寄存器 Short/Unsigned short: 大端(模型 DataFormat=ABCD) → struct '>H'
  - 双寄存器 Float: ByteOrder=CDAB → f2cdab()
  - DI: function 2 离散量

分时轮巡(staggered round-robin)更新:
  全量集设备/点位巨大(full≈2252 设备 / 21 万点)，若一个 tick 内 encode 所有点会造成
  周期性 CPU/锁尖峰并可能拖垮系统。因此 updater 把全部 (ip,slave) 均匀切成 --slots 片，
  每个短 tick(=interval/slots) 只 encode 当前一片，轮流推进 → 每台设备约每 --interval 秒
  被完整刷新一轮，但任意瞬间只算一小批，削平 CPU/内存峰值。

5040 监控 API(供前端 SimulatorMonitor.vue):
  /api/summary       轻量摘要(总设备/IP/点位/型号 + 设备清单, 不含点位负载)
  /api/slaves        与 summary.slaves 同源的轻量清单(兼容旧前端, 不返回寄存器明细)
  /api/slave/{key}   单台设备明细(按需解码当前寄存器), key = "ip#slave"

用法:
  python3 scripts/hx_simulator.py --set validation --test --port 1502   # 自检(无 sudo)
  python3 scripts/hx_simulator.py --set full --local                    # 全量多IP(无 sudo)
  sudo python3 scripts/hx_simulator.py --set full                       # 真实IP(需别名)
"""
import os, json, math, struct, time, threading, signal, argparse, asyncio
from collections import defaultdict
from http.server import HTTPServer, BaseHTTPRequestHandler

ROOT = os.path.join(os.path.dirname(__file__), "..")
MBAP_LEN = 7
MAX_REGS = 1024

LOCK = threading.Lock()
SLAVES = {}     # (ip, slave) -> {"ai":[(addr,type,encfn,base,noise,factor)], "di":[(addr,base)], "hr":[...], "di_vals":{addr:0/1}}
IP_SLAVES = defaultdict(dict)   # ip -> {slave: state}
REGISTRY = {}                   # "ip#slave" -> state  (供 5040 HTTP API 按唯一键检索)


def f2cdab(v):
    p = struct.pack("<f", float(v))
    return (p[1] << 8) | p[0], (p[3] << 8) | p[2]


# ---- 物理真值基准(按点位名/单位猜测) ----
def base_for(name, unit):
    n = name
    if "线电压" in n: return 380.0, 1.5
    if "相电压" in n: return 220.0, 1.5
    if "电压" in n: return 220.0, 1.5
    if "电流" in n:
        if "中性线" in n or "零" in n: return 0.6, 0.1
        if "旁路" in n: return 5.0, 1.0
        return 8.0, 0.4
    if "频率" in n: return 50.0, 0.05
    if "功率因数" in n: return 0.95, 0.01
    if "有功" in n and "电度" not in n and "电能" not in n: return 55.0, 2.0
    if "无功" in n: return 18.0, 1.0
    if "视在" in n: return 58.0, 2.0
    if "电度" in n or "电能" in n: return 12000.0, 5.0
    if "畸变率" in n: return 3.0, 0.6
    if "温度" in n: return 32.0, 1.5
    if "湿度" in n: return 45.0, 3.0
    if "剩余" in n and "时间" in n: return 120.0, 3.0
    if unit == "V": return 220.0, 1.5
    if unit == "A": return 8.0, 0.4
    if unit in ("kW", "kvar", "kVA"): return 50.0, 2.0
    if unit == "Hz": return 50.0, 0.05
    return 10.0, 1.0


# ---- DI 正常态(状态/合闸/通信=1，故障/报警=0) ----
def di_normal(name):
    bad = ("故障", "报警", "告警", "欠压", "过载", "异常", "硬件")
    good = ("合闸", "通信", "通讯", "正常", "状态", "运行", "开关")
    if any(k in name for k in bad):
        return 0
    if any(k in name for k in good):
        return 1
    return 0


def build_slaves(set_name):
    base = os.path.join(ROOT, "hx-data", set_name)
    pkg = json.load(open(os.path.join(base, "中航信_complete_project_package.json")))
    dm = json.load(open(os.path.join(base, "ism_data_models.json")))
    name2muid = {d["name"]: d["uuid"] for d in pkg["deviceModels"]}
    ai_groups = {g["uuid"] for g in pkg["registerGroups"] if g["function"] == 3}
    di_groups = {g["uuid"] for g in pkg["registerGroups"] if g["function"] == 2}
    model_ai = defaultdict(list)
    model_di = defaultdict(list)
    for p in pkg["registerPoints"]:
        if p["registerGroupUuid"] in ai_groups:
            b, noise = base_for(p["name"], p.get("unit", ""))
            model_ai[p["muid"]].append((p["registerAddress"], p["type"],
                                        float(p.get("simFactor", 1.0)), b, noise, p["name"]))
        elif p["registerGroupUuid"] in di_groups:
            model_di[p["muid"]].append((p["registerAddress"], di_normal(p["name"])))

    for d in dm["devices"]:
        muid = name2muid.get(d["modelName"])
        if not muid:
            continue
        ip = d.get("gatewayIP", "127.0.0.1")
        slave = int(d["slaveId"])
        st = {"ai": model_ai.get(muid, []), "di": model_di.get(muid, []),
              "hr": [0] * MAX_REGS, "di_vals": {},
              "ip": ip, "slave": slave, "type": d.get("modelName", "")}
        for addr, b in st["di"]:
            st["di_vals"][addr] = b
        IP_SLAVES[ip][slave] = st
        REGISTRY["%s#%d" % (ip, slave)] = st
    return pkg, dm


# 电度(累计用电量)按墙钟时间从当日 00:00 起单调累增：base + 累计速率 * 当日已过秒数。
# 这样历史趋势图拉到的电度曲线是一条随时间上升的折线(斜率≈用电增量)，而非恒定平线；
# 每日 0 点回到 base，符合「今日累计用电量」语义。速率取 ~55kWh/h(与有功功率量级一致)。
_ENERGY_RATE_PER_SEC = 55.0 / 3600.0


def _energy_accum_value(b, noise):
    lt = time.localtime()
    secs_since_midnight = lt.tm_hour * 3600 + lt.tm_min * 60 + lt.tm_sec
    jitter = (noise * 0.15) * math.sin(time.time() * 0.05)
    return b + _ENERGY_RATE_PER_SEC * secs_since_midnight + jitter


def encode_ai(st, t):
    hr = st["hr"]
    for addr, typ, fac, b, noise, name in st["ai"]:
        if "电度" in name or "电能" in name:
            v = _energy_accum_value(b, noise)
        else:
            phase = (hash(name) % 628) / 100.0
            v = b + noise * (math.sin(t * 0.4 + phase) * 0.6 +
                             math.sin(t * 1.3 + phase * 1.7) * 0.3 +
                             math.sin(t * 2.7 + phase * 2.3) * 0.1)
        if typ == "Float":
            lo, hi = f2cdab(v)
            if addr + 1 < len(hr):
                hr[addr] = lo; hr[addr + 1] = hi
        else:  # 单寄存器整型: raw = 真值/系数, 大端
            raw = int(round(v / fac)) if fac else int(round(v))
            if typ == "Short":
                raw = max(-32768, min(32767, raw)) & 0xFFFF
            else:
                raw = max(0, min(65535, raw))
            hr[addr] = raw


async def updater(interval, slots):
    """分时轮巡错峰更新(asyncio 协程版)：把全部设备均匀切成 slots 片，每 tick 只 encode 一片。

    - tick = interval / slots
    - 第 k 个 tick 处理 buckets[k]，k 循环 0..slots-1
    - 每台设备约每 interval 秒被完整刷新一轮，但任意瞬间只算 1/slots 的设备
    分片用「条带化」(items[i::slots])，使每片设备跨多个 IP 均匀分布，避免单 IP 扎堆。

    运行于与所有 Modbus 连接处理器同一个事件循环线程：encode 与后端读取在同一线程
    协作交替，绝不会同时执行 → 后端永远读不到「写一半」的浮点寄存器，且全程无需加锁。
    """
    items = [st for ip in IP_SLAVES for st in IP_SLAVES[ip].values()]
    slots = max(1, min(slots, len(items) or 1))
    buckets = [items[i::slots] for i in range(slots)]
    tick = interval / slots
    total_pts = sum(len(st["ai"]) for st in items)
    print("[updater] 分时轮巡(asyncio): 设备=%d 切成 %d 片, 每 %.2fs 更新 1 片, 每设备约每 %gs 一轮 (总AI点=%d)"
          % (len(items), slots, tick, interval, total_pts), flush=True)
    k = 0
    cycle = 0
    while True:
        t = time.monotonic()
        batch = buckets[k]
        pts = 0
        for st in batch:
            encode_ai(st, t)
            pts += len(st["ai"])
        pct = pts * 100 // total_pts if total_pts else 0
        print("[tick] slot %d/%d 更新 %d 台设备 / %d 个AI点 (占总量 %d%%)"
              % (k + 1, slots, len(batch), pts, pct), flush=True)
        k += 1
        if k >= slots:
            k = 0
            cycle += 1
            print("[cycle] 第 %d 轮全量刷新完成: %d 台设备全部更新一次 (≈%gs/轮)"
                  % (cycle, len(items), interval), flush=True)
        await asyncio.sleep(tick)


# ─────────────────────── 5040 监控 API（解码 + 轻量摘要 + 按需明细） ───────────────────────

def decode_ai_value(st, addr, typ, fac):
    """把当前寄存器值解码回物理量估计(与 encode_ai 对齐)，供监控页展示。"""
    hr = st["hr"]
    if typ == "Float":
        if addr + 1 >= len(hr):
            return None
        lo, hi = hr[addr], hr[addr + 1]
        raw = bytes([lo & 0xFF, (lo >> 8) & 0xFF, hi & 0xFF, (hi >> 8) & 0xFF])
        return round(struct.unpack("<f", raw)[0], 3)
    raw = hr[addr]
    if typ == "Short" and raw >= 0x8000:
        raw -= 0x10000
    return round(raw * fac, 3)


def slave_detail(st):
    hr = {}
    for addr, typ, fac, b, noise, name in st["ai"]:
        hr[addr] = {"name": name, "raw": st["hr"][addr],
                    "value": decode_ai_value(st, addr, typ, fac)}
    di = {addr: {"name": "DI", "value": st["di_vals"].get(addr, 0)} for addr, _ in st["di"]}
    return {"slave": "%s#%d" % (st["ip"], st["slave"]), "device_type": st["type"],
            "ip": st["ip"], "slave_id": st["slave"],
            "holding_registers": hr, "discrete_inputs": di}


def build_summary(meta):
    by_ip = sorted(({"ip": ip, "slaves": len(IP_SLAVES[ip])} for ip in IP_SLAVES),
                   key=lambda x: -x["slaves"])
    type_count = defaultdict(int)
    slaves = []
    total_pts = 0
    for ip in IP_SLAVES:
        for slave, st in IP_SLAVES[ip].items():
            type_count[st["type"]] += 1
            total_pts += len(st["ai"]) + len(st["di"])
            slaves.append({"id": "%s#%d" % (ip, slave), "slave_id": slave,
                           "ip": ip, "type": st["type"]})
    return {"total_slaves": len(REGISTRY), "ip_count": len(IP_SLAVES),
            "total_points": total_pts,
            "device_types": [{"type": t, "count": c}
                             for t, c in sorted(type_count.items(), key=lambda x: -x[1])],
            "by_ip": by_ip, "slaves": slaves, **meta}


class APIHandler(BaseHTTPRequestHandler):
    summary_meta = {}

    def log_message(self, *a):
        pass

    def _send(self, code, data):
        body = json.dumps(data, ensure_ascii=False).encode()
        self.send_response(code)
        self.send_header("Content-Type", "application/json; charset=utf-8")
        self.send_header("Access-Control-Allow-Origin", "*")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    def do_GET(self):
        from urllib.parse import unquote
        path = unquote(self.path.split("?")[0]).rstrip("/")
        if path == "/api/summary":
            with LOCK:
                self._send(200, build_summary(self.summary_meta))
        elif path == "/api/slaves":
            # 兼容旧前端：仅返回轻量清单(不含寄存器明细)，避免全量一次性返回拖垮前端
            with LOCK:
                self._send(200, build_summary(self.summary_meta)["slaves"])
        elif path.startswith("/api/slave/"):
            key = path[len("/api/slave/"):]
            st = REGISTRY.get(key)
            if st is None and key.isdigit():
                st = next((s for s in REGISTRY.values() if s["slave"] == int(key)), None)
            if st is None:
                return self._send(404, {"error": "not found", "key": key})
            with LOCK:
                self._send(200, slave_detail(st))
        else:
            self._send(200, {"status": "ok", "version": "hx-package-driven"})


def start_http(port, meta):
    APIHandler.summary_meta = meta
    HTTPServer(("0.0.0.0", port), APIHandler).serve_forever()


# 连接空闲超时(秒)：后端在全量(2252台/54网关)下会高频新建连接并丢弃旧连接。
# asyncio 模型下连接是协程而非线程，不存在线程泄漏；但仍给读取设超时，及时回收被丢弃的
# 半开连接(关闭 socket + 取消协程)，避免空闲连接无限堆积占用 fd。
CONN_IDLE_TIMEOUT = 30


def build_modbus_response(slaves, uid, pdu):
    """纯函数：根据 slave 表 + PDU 计算 Modbus 响应体(不含 MBAP 头)。无锁——
    运行于事件循环单线程，与 updater 协作交替，绝不并发访问寄存器。"""
    fc = pdu[0]
    st = slaves.get(uid)
    if st is None:
        return bytes([fc | 0x80, 0x0A])
    if fc in (3, 4) and len(pdu) >= 5:
        start = (pdu[1] << 8) | pdu[2]; cnt = (pdu[3] << 8) | pdu[4]
        if 1 <= cnt <= 125:
            regs = st["hr"][start:start + cnt]
            regs += [0] * (cnt - len(regs))
            return bytes([fc, cnt * 2]) + b"".join(struct.pack(">H", r & 0xFFFF) for r in regs)
        return bytes([fc | 0x80, 0x03])
    if fc in (1, 2) and len(pdu) >= 5:
        start = (pdu[1] << 8) | pdu[2]; cnt = (pdu[3] << 8) | pdu[4]
        if 1 <= cnt <= 2000:
            vals = [st["di_vals"].get(start + i, 0) for i in range(cnt)]
            bc = (cnt + 7) // 8
            bits = bytearray(bc)
            for i, v in enumerate(vals):
                if v:
                    bits[i // 8] |= 1 << (i % 8)
            return bytes([fc, bc]) + bytes(bits)
        return bytes([fc | 0x80, 0x03])
    if fc in (5, 6) and len(pdu) >= 5:
        return pdu
    if fc in (15, 16) and len(pdu) >= 6:
        return pdu[:5]
    return bytes([fc | 0x80, 0x01])


def make_handler(slaves):
    async def handle(reader, writer):
        try:
            while True:
                try:
                    hdr = await asyncio.wait_for(reader.readexactly(MBAP_LEN), CONN_IDLE_TIMEOUT)
                except (asyncio.TimeoutError, asyncio.IncompleteReadError, ConnectionError, OSError):
                    break
                tid = (hdr[0] << 8) | hdr[1]; pid = (hdr[2] << 8) | hdr[3]
                length = (hdr[4] << 8) | hdr[5]; uid = hdr[6]
                pdu_len = length - 1
                if pdu_len <= 0 or pdu_len > 253:
                    break
                try:
                    pdu = await asyncio.wait_for(reader.readexactly(pdu_len), CONN_IDLE_TIMEOUT)
                except (asyncio.TimeoutError, asyncio.IncompleteReadError, ConnectionError, OSError):
                    break
                resp = build_modbus_response(slaves, uid, pdu)
                writer.write(struct.pack(">HHHB", tid, pid, 1 + len(resp), uid) + resp)
                try:
                    await writer.drain()
                except (ConnectionError, OSError):
                    break
        except Exception:
            pass
        finally:
            try:
                writer.close()
            except Exception:
                pass
    return handle


async def run_event_loop(listen_specs):
    """listen_specs: list of (host, port, slaves). 在单事件循环里同时监听所有端口。"""
    servers = []
    for host, port, slaves in listen_specs:
        srv = await asyncio.start_server(make_handler(slaves), host, port, limit=2 ** 18)
        servers.append(srv)
        print(f"  [listen] {host}:{port}  slaves={len(slaves)}", flush=True)
    await asyncio.gather(*(s.serve_forever() for s in servers))


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--set", choices=["validation", "full"], default="validation")
    ap.add_argument("--test", action="store_true", help="测试模式: 仅 127.0.0.1 单服务")
    ap.add_argument("--local", action="store_true",
                    help="无 sudo 全量验证: 每网关在 127.0.0.1:(1502+序号) 起服务，与导入器 --local-sim 对齐")
    ap.add_argument("--host", default="127.0.0.1")
    ap.add_argument("--port", type=int, default=1502)
    ap.add_argument("--interval", type=float, default=10.0,
                    help="每台设备完整刷新一轮的周期(秒)，默认 10；怕崩溃可调大")
    ap.add_argument("--slots", type=int, default=20,
                    help="分时轮巡片数：每 interval/slots 秒只更新 1/slots 的设备")
    ap.add_argument("--http-port", type=int, default=5040,
                    help="5040 监控 API 端口(供前端 SimulatorMonitor)；0 关闭")
    args = ap.parse_args()

    _pkg, _dm = build_slaves(args.set)
    all_gw = sorted({d["gatewayIP"] for d in _dm["devices"]})   # 与导入器 --local-sim 同源同序
    total = sum(len(s) for s in IP_SLAVES.values())
    print(f"中航信模拟器  set={args.set}  IP数={len(IP_SLAVES)}  设备(slave)={total}", flush=True)
    for ip in IP_SLAVES:
        print(f"   {ip}: {len(IP_SLAVES[ip])} slaves")

    if args.http_port:
        meta = {"set": args.set, "interval": args.interval, "slots": args.slots}
        threading.Thread(target=lambda: start_http(args.http_port, meta), daemon=True).start()
        print(f"  [HTTP] 监控 API on 0.0.0.0:{args.http_port}", flush=True)

    # 组装监听规格 (host, port, slaves)，全部交给单个事件循环
    if args.test:
        # 合并所有 slave 到单个服务(测试解码用)
        merged = {}
        for ip in IP_SLAVES:
            merged.update(IP_SLAVES[ip])
        listen_specs = [(args.host, args.port, merged)]
        print(f"[TEST] {args.host}:{args.port} 合并 {len(merged)} slaves", flush=True)
    elif args.local:
        # 无 sudo: 每网关映射到 127.0.0.1:(1502+序号)，端口顺序与导入器 --local-sim 一致
        listen_specs = [("127.0.0.1", 1502 + i, IP_SLAVES.get(ip, {})) for i, ip in enumerate(all_gw)]
        print(f"  [LOCAL] 127.0.0.1:1502..{1502 + len(all_gw) - 1}  共 {len(all_gw)} 网关", flush=True)
    else:
        listen_specs = [(ip, 502, IP_SLAVES[ip]) for ip in IP_SLAVES]
        print("  [Running, Ctrl+C to stop]", flush=True)

    async def _amain():
        asyncio.create_task(updater(args.interval, args.slots))
        await run_event_loop(listen_specs)

    try:
        asyncio.run(_amain())
    except KeyboardInterrupt:
        print("\n stopped")


if __name__ == "__main__":
    try:
        signal.signal(signal.SIGPIPE, signal.SIG_IGN)
    except Exception:
        pass
    main()
