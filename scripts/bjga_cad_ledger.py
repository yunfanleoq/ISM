#!/usr/bin/env python3
"""从 10 月 23 日平面图和动环系统图生成公安局动环台账，并写回 dcim.json。"""
import hashlib
import json
import math
import os
import re
from collections import defaultdict

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
PING = "/tmp/bjga-cad-read/pingmian.json"
SYS = "/tmp/bjga-cad-read/sys.json"
DCIM = os.path.join(ROOT, "ism_server_user/data/bjga/dcim.json")
LEDGER = os.path.join(ROOT, "docs/北京公安-现场资料/台账.md")

BANDS = {
    "3F": (-231247, -171847),
    "4F": (-171847, -112447),
    "5F": (-112447, -53047),
}
ELEV = {"3F": 11.4, "4F": 16.8, "5F": 22.2}
ORDER = ["3F", "4F", "5F"]
X0, X1 = 820000, 880000
# 原型列头柜曲线。重跑时固定用这份底数，不能再乘上一轮已经缩放过的输出。
SERIES_DATES = [
    "2026-10-01 00:00", "2026-10-02 00:00", "2026-10-03 00:00", "2026-10-04 00:00",
    "2026-10-05 00:00", "2026-10-06 00:00", "2026-10-07 00:00", "2026-10-08 00:00",
    "2026-10-09 00:00", "2026-10-10 00:00", "2026-10-11 00:00", "2026-10-12 00:00",
]
SERIES_BASE = [384.02, 384.18, 384.14, 384.04, 384.06, 384.09, 384.13, 384.24, 383.97, 384.22, 384.02, 384.25]


def hkey(h):
    return h[-1] if isinstance(h, list) and h else None


def floor_of(y):
    # 三张平面在图上上下排开，机柜 Y 之间有空档，按空档切开，避免标题线把同一层切成两半。
    if y < -180000:
        return "3F"
    if y < -120000:
        return "4F"
    if -120000 <= y < -60000:
        return "5F"
    return None


def clean_text(t):
    t = re.sub(r"\\[A-Za-z][^;]*;", "", t or "")
    return t.replace("\\P", " ").strip()


def load_dwg(path):
    objs = json.load(open(path))["OBJECTS"]
    layers, blocks = {}, {}
    for o in objs:
        h = o.get("handle") or []
        key = h[-1] if h else None
        if o.get("object") == "LAYER":
            layers[key] = o.get("name")
        if o.get("object") == "BLOCK_HEADER":
            blocks[key] = o.get("name")
    return objs, layers, blocks


def poly_area(pts):
    if len(pts) < 3:
        return 0
    s = 0
    for i, p in enumerate(pts):
        q = pts[(i + 1) % len(pts)]
        s += p[0] * q[1] - q[0] * p[1]
    return abs(s) / 2


def outer_polygon(segments):
    adj = defaultdict(set)
    for a, b in segments:
        if a == b:
            continue
        adj[a].add(b)
        adj[b].add(a)
    if len(adj) < 3:
        return []

    def ang(p, q):
        return math.atan2(q[1] - p[1], q[0] - p[0])

    start = min(adj)
    nbrs = list(adj[start])
    cur = min(nbrs, key=lambda q: ang(start, q))
    poly = [start]
    prev = start
    seen = set()
    for _ in range(len(adj) * 3):
        if cur == start and len(poly) > 2:
            break
        poly.append(cur)
        mark = (prev, cur)
        if mark in seen:
            break
        seen.add(mark)
        cands = [n for n in adj[cur] if n != prev]
        if not cands:
            break
        heading = ang(prev, cur)

        def turn(n):
            d = ang(cur, n) - heading
            while d <= 1e-6:
                d += 2 * math.pi
            return d

        nxt = min(cands, key=turn)
        prev, cur = cur, nxt
    if len(poly) > 2 and poly[-1] != poly[0]:
        if math.hypot(poly[-1][0] - poly[0][0], poly[-1][1] - poly[0][1]) < 80:
            poly[-1] = poly[0]
    if poly and poly[0] == poly[-1]:
        poly = poly[:-1]
    return poly


def point_in(poly, x, y):
    inside = False
    j = len(poly) - 1
    for i in range(len(poly)):
        xi, yi = poly[i]
        xj, yj = poly[j]
        if ((yi > y) != (yj > y)) and (x < (xj - xi) * (y - yi) / ((yj - yi) or 1e-9) + xi):
            inside = not inside
        j = i
    return inside


def r2(n):
    return round(n, 2)


def shash(key):
    return int(hashlib.md5(key.encode()).hexdigest()[:8], 16)


def sim_kw(key):
    return round(1.2 + (shash(key) % 280) / 100, 2)


SRV_MODELS = (
    ("浪潮", "CS5280H"),
    ("华为", "2288H V5"),
    ("华三", "R4900 G3"),
)
SIM_ORGS = ("未登记", "模拟单位甲", "模拟单位乙", "模拟单位丙")


def sim_servers_for(cabs, room_name, floor):
    """图纸没有服务器清单。有号柜 6–12 台 2U，无号柜 2–4 台，U 位不重叠。"""
    out = []
    for c in cabs:
        name = c["name"]
        blank = str(c.get("id") or "").startswith("无号")
        n = (2 + (shash(name + ":n") % 3)) if blank else (6 + (shash(name + ":n") % 7))
        free = 42 - n * 2
        gaps = [0] * (n + 1)
        for i in range(free):
            gaps[shash("%s:g%d" % (name, i)) % (n + 1)] += 1
        u = 1 + gaps[0]
        kw_sum = 0.0
        for i in range(n):
            brand, model = SRV_MODELS[shash("%s:m%d" % (name, i)) % 3]
            sn = "SIM-%s-%02d" % (name, i + 1)
            kw = round(0.28 + (shash(sn + ":k") % 25) / 100.0, 2)
            u_end = u + 1
            out.append({
                "room": room_name,
                "floor": floor,
                "cabinet": name,
                "sn": sn,
                "status": "使用中",
                "uStart": u,
                "uEnd": u_end,
                "uH": 2,
                "cat": "服务器",
                "type": "机架服务器",
                "brand": brand,
                "model": model,
                "kw": kw,
                "psu": 2,
                "weight": 20 + (shash(sn + ":w") % 9),
                "net": "未登记",
                "biz": "模拟业务",
                "org": SIM_ORGS[shash(sn + ":o") % 4],
                "owner": "",
                "source": "sim",
            })
            kw_sum += kw
            u = u_end + 1 + gaps[i + 1]
        used = n * 2
        rated = 10.0
        space = round(min(100.0, (used + n) / 42.0 * 100), 1)
        powr = round(kw_sum / rated * 100, 1)
        if_used = n if n < 20 else 20
        ifr = round(if_used / 20.0 * 100, 1)
        comp = round(max(space, powr, ifr), 1)
        c["dev"] = n
        c["usedU"] = used
        c["kw"] = round(kw_sum, 2)
        c["rtKw"] = round(kw_sum, 2)
        c["rated"] = rated
        c["coef"] = 1
        c["spaceRate"] = space
        c["powRate"] = powr
        c["ifUsed"] = if_used
        c["ifRate"] = ifr
        c["comp"] = comp
        c["avgRate"] = round(comp * 0.9, 1)
        c["peakRate"] = round(min(99.0, comp * 1.08), 1)
        c["topCat"] = "服务器"
        c["lv"] = "red" if comp > 90 else ("yellow" if comp >= 70 else "blue")
    return out


def verify_sim_servers(assets, rooms):
    by = defaultdict(list)
    for a in assets:
        if a.get("cat") != "服务器":
            continue
        if a.get("source") != "sim":
            raise SystemExit("服务器 source 不是 sim: " + a.get("sn", ""))
        by[a["cabinet"]].append(a)
    names = {}
    for rm in rooms:
        for c in rm.get("cabinets") or []:
            names[c["name"]] = c
            got = by.get(c["name"], [])
            if c.get("dev") != len(got) or c.get("usedU") != sum(x["uH"] for x in got):
                raise SystemExit("机柜 dev/usedU 与服务器不一致: " + c["name"])
            used = set()
            for x in got:
                if x["uH"] != 2 or x["uEnd"] != x["uStart"] + 1:
                    raise SystemExit("不是 2U: " + x["sn"])
                for u in range(x["uStart"], x["uEnd"] + 1):
                    if u < 1 or u > 42 or u in used:
                        raise SystemExit("U 位重叠或越界: " + c["name"])
                    used.add(u)
            if not got:
                raise SystemExit("空柜: " + c["name"])
    orphan = [k for k in by if k not in names]
    if orphan:
        raise SystemExit("服务器挂在未知柜: " + orphan[0])


def main():
    objs, layers, blocks = load_dwg(PING)

    def lay(o):
        return layers.get(hkey(o.get("layer")), "")

    walls = defaultdict(list)
    doors = defaultdict(list)
    cabs = defaultdict(list)
    texts = defaultdict(list)
    for o in objs:
        e = o.get("entity")
        name = lay(o)
        if e == "LINE":
            s, t = o.get("start"), o.get("end")
            if not s or not t:
                continue
            mx, my = (s[0] + t[0]) / 2, (s[1] + t[1]) / 2
            fl = floor_of(my)
            if not fl or not (X0 < mx < X1):
                continue
            a = (round(s[0] / 50) * 50, round(s[1] / 50) * 50)
            b = (round(t[0] / 50) * 50, round(t[1] / 50) * 50)
            if a == b:
                continue
            if name == "WALL":
                walls[fl].append((a, b))
            elif name == "DOOR_FIRE":
                length = math.hypot(b[0] - a[0], b[1] - a[1])
                if 600 <= length <= 2400:
                    doors[fl].append((a, b, length))
            continue
        pt = o.get("ins_pt")
        if not pt:
            continue
        fl = floor_of(pt[1])
        if not fl or not (X0 < pt[0] < X1):
            continue
        if e == "INSERT":
            bn = blocks.get(hkey(o.get("block_header")), "")
            if bn in ("jigui-1200", "jigui-1400"):
                cabs[fl].append({
                    "block": bn,
                    "x": pt[0],
                    "y": pt[1],
                    "rot": o.get("rotation") or 0,
                })
        elif e in ("TEXT", "MTEXT"):
            tx = clean_text(o.get("text_value") or o.get("text") or "")
            if tx:
                texts[fl].append((tx, pt[0], pt[1]))

    floors_out = []
    rooms_out = []
    assets = []
    sim_srv = []
    md = ["# 北京公安项目设备台账", "", "来源：2026-10-23 平面定位图与动环监控系统图。机柜和房间是图纸上的对象。图纸没有服务器清单，柜内服务器、功率、温度、告警都是开发环境里的模拟数，来源标为模拟。", ""]

    for fl in ORDER:
        segs = list({tuple(sorted((a, b))): (a, b) for a, b in walls[fl]}.values())
        xs, ys = [], []
        for a, b in segs:
            xs += [a[0], b[0]]
            ys += [a[1], b[1]]
        ox, oy = min(xs), min(ys)
        span_x = (max(xs) - ox) / 1000
        span_y = (max(ys) - oy) / 1000

        def mxy(x, y, ox=ox, oy=oy):
            return ((x - ox) / 1000, (y - oy) / 1000)

        wall_m = []
        for a, b in segs:
            p, q = mxy(*a), mxy(*b)
            if math.hypot(q[0] - p[0], q[1] - p[1]) < 0.2:
                continue
            wall_m.append([r2(p[0]), r2(p[1]), r2(q[0]), r2(q[1])])
        outline = outer_polygon(segs)
        area = poly_area(outline) / 1e6
        bbox_area = span_x * span_y
        if area < bbox_area * 0.45 or len(outline) < 4:
            outline = [(min(xs), min(ys)), (max(xs), min(ys)), (max(xs), max(ys)), (min(xs), max(ys))]
        outline_m = [[r2(p[0]), r2(p[1])] for p in (mxy(*p) for p in outline)]

        # 机柜对编号
        labels = [(t, x, y) for t, x, y in texts[fl] if re.match(r"^[A-J]-\d{2}$", t)]
        hall_cabs = {"1#": [], "2#": []}
        unlabeled = []
        # 同一侧、同一列：机柜归到 Y 最近的列号，再按 X 与 A-01 起的文字一一对齐。
        assign = {}
        for side_name, side_ok in (("1#", lambda x: x < 851000), ("2#", lambda x: x >= 851000)):
            idxs = [i for i, c in enumerate(cabs[fl]) if side_ok(c["x"])]
            rows = defaultdict(list)
            for t, x, y in labels:
                if side_ok(x):
                    rows[t[0]].append((int(t[2:]), x, y, t))
            if not rows or not idxs:
                continue
            means = {L: sum(p[2] for p in v) / len(v) for L, v in rows.items()}
            groups = defaultdict(list)
            for i in idxs:
                cy = cabs[fl][i]["y"]
                letter = min(means, key=lambda L: abs(cy - means[L]))
                if abs(cy - means[letter]) > 1600:
                    continue
                groups[letter].append(i)
            for letter, members in groups.items():
                labs = sorted(rows[letter])
                members.sort(key=lambda i: cabs[fl][i]["x"])
                if len(members) == len(labs) + 1:
                    def far(i, labs=labs):
                        return min(math.hypot(cabs[fl][i]["x"] - x, cabs[fl][i]["y"] - y) for _, x, y, _ in labs)
                    members = [i for i in members if i != max(members, key=far)]
                if len(members) != len(labs):
                    continue
                for i, lab in zip(members, labs):
                    assign[i] = lab[3]
        for ci, c in enumerate(cabs[fl]):
            side = "1#" if c["x"] < 851000 else "2#"
            best = assign.get(ci)
            w = 1.4 if c["block"] == "jigui-1400" else 1.2
            along = abs(math.sin(c["rot"])) < 0.5
            ww, dd = (w, 0.6) if along else (0.6, w)
            px, py = mxy(c["x"], c["y"])
            item = {
                "x": r2(px - ww / 2),
                "y": r2(py - dd / 2),
                "w": w if along else 0.6,
                "d": 0.6 if along else w,
                "kind": "1400" if c["block"] == "jigui-1400" else "1200",
                "cx": c["x"],
                "cy": c["y"],
            }
            if not best:
                item["id"] = ""
                unlabeled.append(item)
            else:
                item["id"] = best
                col = ord(item["id"][0]) - ord("A") + 1
                seq = int(item["id"][2:])
                item["col"] = col
                item["seq"] = seq
            hall_cabs[side].append(item)

        # 房间名
        room_labels = []
        seen_lab = set()
        for t, x, y in texts[fl]:
            kind = None
            if "数据机房" in t:
                kind = "主机房"
            elif t in ("变配电室", "操作间", "及电池间") or t.startswith("空调"):
                kind = {"变配电室": "配电室", "操作间": "操作间", "及电池间": "电池室"}.get(t, "空调间")
            elif t in ("强电井", "弱电", "机房弱电井", "风井", "水管井", "新风", "排风", "排烟", "钢瓶间", "恒湿机"):
                kind = "井道" if "井" in t else "配套"
            if not kind:
                continue
            key = (t, round(x / 2000), round(y / 2000))
            if key in seen_lab:
                continue
            seen_lab.add(key)
            room_labels.append((t, kind, x, y))

        code = {"3F": "3", "4F": "4", "5F": "5"}[fl]
        hall_name = {
            "1#": "%s层1#数据机房" % code,
            "2#": "%s层2#数据机房" % code,
        }
        hall_code = {"1#": code + "M1", "2#": code + "M2"}

        def room_shell(name, typ, code_name, cab_list, anchor):
            if cab_list:
                xs2 = [c["x"] for c in cab_list] + [c["x"] + c["w"] for c in cab_list]
                ys2 = [c["y"] for c in cab_list] + [c["y"] + c["d"] for c in cab_list]
                pad = 1.2
                poly = [
                    [r2(min(xs2) - pad), r2(min(ys2) - pad)],
                    [r2(max(xs2) + pad), r2(min(ys2) - pad)],
                    [r2(max(xs2) + pad), r2(max(ys2) + pad)],
                    [r2(min(xs2) - pad), r2(max(ys2) + pad)],
                ]
            else:
                ax, ay = mxy(*anchor)
                poly = [[r2(ax - 4), r2(ay - 3)], [r2(ax + 4), r2(ay - 3)], [r2(ax + 4), r2(ay + 3)], [r2(min(ax + 4, span_x)), r2(ay + 3)]]
                poly = [[r2(ax - 4), r2(ay - 3)], [r2(ax + 4), r2(ay - 3)], [r2(ax + 4), r2(ay + 3)], [r2(ax - 4), r2(ay + 3)]]
            cabs_json = []
            blank_n = 0
            for c in cab_list:
                cid = c.get("id") or ""
                if not cid:
                    blank_n += 1
                    cid = "无号" if blank_n == 1 else ("无号-%d" % blank_n)
                comp = 40 + (shash(fl + name + cid) % 45)
                cabs_json.append({
                    "name": (code_name + "-" + cid) if code_name else (name + "-" + cid),
                    "id": cid,
                    "col": c.get("col") or 0,
                    "seq": c.get("seq") or 0,
                    "x": c["x"], "y": c["y"], "w": c["w"], "d": c["d"],
                    "kind": c["kind"],
                    "comp": comp,
                    "dev": 0,
                    "source": "sim",
                })
            sim_srv.extend(sim_servers_for(cabs_json, name, fl))
            temp = round(22 + (shash(name) % 30) / 10, 1)
            humi = round(45 + (shash(name + "h") % 120) / 10, 1)
            alarm_n = shash(name + "a") % 6
            return {
                "name": name,
                "floor": fl,
                "type": typ,
                "code": code_name,
                "zone": "N" if anchor[0] < 851000 else "S",
                "cabCount": len(cabs_json),
                "dev": len(cabs_json),
                "kw": round(sum(c.get("rtKw") or 0 for c in cabs_json), 1) if cabs_json else round(sim_kw(name), 1),
                "usedU": sum(c.get("usedU") or 0 for c in cabs_json),
                "usedUEff": sum((c.get("usedU") or 0) + (c.get("dev") or 0) * (c.get("coef") or 1) for c in cabs_json),
                "capU": len(cabs_json) * 42,
                "rate": round(sum(c["comp"] for c in cabs_json) / len(cabs_json), 1) if cabs_json else 0,
                "comp": round(sum(c["comp"] for c in cabs_json) / len(cabs_json), 1) if cabs_json else 0,
                "alarmN": alarm_n,
                "alarmLv": [{"name": "次要", "value": alarm_n}] if alarm_n else [],
                "temp": temp,
                "humi": humi,
                "source": "sim",
                "cad": True,
                "poly": poly,
                "cabinets": cabs_json,
                "layout": None,
            }

        floor_rooms = []
        for side in ("1#", "2#"):
            floor_rooms.append(room_shell(hall_name[side], "主机房", hall_code[side], hall_cabs[side], (830000 if side == "1#" else 860000, (BANDS[fl][0] + BANDS[fl][1]) / 2)))
        # 其它房间：同名靠太近的只留一个
        kept = []
        for t, kind, x, y in room_labels:
            if "数据机房" in t:
                continue
            if any(math.hypot(x - x2, y - y2) < 2500 and t == t2 for t2, _, x2, y2 in kept):
                continue
            kept.append((t, kind, x, y))
        name_n = defaultdict(int)
        for t, kind, x, y in kept:
            name_n[t] += 1
        seen_n = defaultdict(int)
        for t, kind, x, y in kept:
            seen_n[t] += 1
            nm = ("%s层%s" % (code, t)) if name_n[t] == 1 else ("%s层%s-%d" % (code, t, seen_n[t]))
            floor_rooms.append(room_shell(nm, kind, "", [], (x, y)))

        door_m = []
        for a, b, length in doors[fl]:
            mx, my = (a[0] + b[0]) / 2, (a[1] + b[1]) / 2
            px, py = mxy(mx, my)
            dx, dy = (b[0] - a[0]) / 1000, (b[1] - a[1]) / 1000
            nearest, nd = floor_rooms[0]["name"], 1e18
            for rm in floor_rooms:
                cx = sum(p[0] for p in rm["poly"]) / 4
                cy = sum(p[1] for p in rm["poly"]) / 4
                d = math.hypot(px - cx, py - cy)
                if d < nd:
                    nd, nearest = d, rm["name"]
            door_m.append({
                "x": r2(px), "y": r2(py),
                "dx": r2(dx), "dy": r2(dy),
                "len": r2(length / 1000),
                "room": nearest,
            })

        rooms_out.extend(floor_rooms)
        cab_n = sum(r["cabCount"] for r in floor_rooms)
        floors_out.append({
            "floor": fl,
            "name": code + "层",
            "elev": ELEV[fl],
            "rooms": len(floor_rooms),
            "cabCount": cab_n,
            "dev": cab_n,
            "kw": round(sum(r["kw"] for r in floor_rooms), 1),
            "rate": round(sum(r["comp"] for r in floor_rooms if r["cabCount"]) / max(1, sum(1 for r in floor_rooms if r["cabCount"])), 1),
            "alarmN": sum(r["alarmN"] for r in floor_rooms),
            "w": r2(span_x),
            "d": r2(span_y),
            "outline": outline_m,
            "walls": wall_m,
            "doors": door_m,
            "source": "cad",
        })

        md.append("## %s层" % code)
        md.append("")
        md.append("外轮廓宽 %.3f m，墙线包住的进深 %.3f m，标高 %.1f m。墙段 %d，防火门 %d。" % (span_x, span_y, ELEV[fl], len(wall_m), len(door_m)))
        md.append("")
        for rm in floor_rooms:
            if rm["type"] != "主机房":
                continue
            ids = [c["id"] for c in rm["cabinets"] if c["id"] and c["id"] != "无号"]
            by = defaultdict(list)
            for i in ids:
                by[i[0]].append(int(i[2:]))
            parts = []
            for col in sorted(by):
                ns = sorted(set(by[col]))
                parts.append("%s-%02d–%s-%02d（%d）" % (col, ns[0], col, ns[-1], len(ns)))
            blank = sum(1 for c in rm["cabinets"] if not c["id"] or c["id"] == "无号")
            wide = sum(1 for c in rm["cabinets"] if c["kind"] == "1400")
            md.append("- **%s**（%s）%d 台，其中 1400 mm %d 台%s" % (
                rm["name"], rm["code"], rm["cabCount"], wide,
                ("，有柜无号 %d 台" % blank) if blank else ""))
            md.append("  - " + "，".join(parts))
        others = [rm["name"] + "（" + rm["type"] + "）" for rm in floor_rooms if rm["type"] != "主机房"]
        md.append("- 其他房间：" + "、".join(others))
        md.append("")

        for rm in floor_rooms:
            for c in rm["cabinets"]:
                assets.append({
                    "room": rm["name"],
                    "floor": fl,
                    "cabinet": rm["code"] + "-" + (c["id"] or "无号"),
                    "sn": rm["code"] + "-" + (c["id"] or ("X%d" % (shash("%s-%s" % (c["x"], c["y"])) % 10000))),
                    "status": "在架",
                    "uStart": 0, "uEnd": 0, "uH": 0,
                    "cat": "机柜",
                    "type": "服务器机柜",
                    "brand": "华三",
                    "model": "机柜" + c["kind"],
                    "kw": sim_kw(rm["name"] + c["name"]),
                    "psu": 0, "weight": 0,
                    "net": "", "biz": "机柜",
                    "org": "未登记", "owner": "",
                    "source": "sim",
                })

    assets.extend(sim_srv)
    md.append("## 模拟服务器")
    md.append("")
    md.append("平面图没有服务器清单。有号机柜按稳定随机放 6–12 台 2U 服务器，无号柜 2–4 台，U 位不重叠。source=sim，不编造图上没有的柜号。共 %d 台。" % len(sim_srv))
    md.append("")

    # 系统图上有名字的设备
    named = [
        ("楼", "三层ECC", "汇聚交换机A", "24口动环监控汇聚交换机A"),
        ("楼", "三层ECC", "汇聚交换机B", "24口动环监控汇聚交换机B"),
        ("楼", "三层ECC", "管理服务器A", "动环监控系统管理服务器主机A"),
        ("楼", "三层ECC", "管理服务器B", "动环监控系统管理服务器主机B"),
    ]
    for fl, code_n in (("3F", "3"), ("4F", "4"), ("5F", "5")):
        room = "%s层及电池间" % code_n
        for side in ("A", "B"):
            named.append((fl, room, "交换机", "%sF-UPS-动环接入交换机%s-1" % (code_n, side)))
            for k in ("1", "2"):
                named.append((fl, room, "采集器", "%sF-UPS-动环采集器%s-%s" % (code_n, side, k)))
        for hall, title in (("M1", "1#数据机房"), ("M2", "2#数据机房")):
            room = "%s层%s" % (code_n, title)
            named.append((fl, room, "交换机", "%sF-%s%s-动环接入交换机A-1" % (code_n, code_n, hall)))
            for k in ("1", "2"):
                named.append((fl, room, "采集器", "%sF-%s%s-动环采集器A-%s" % (code_n, code_n, hall, k)))

    branches = []
    if os.path.exists(SYS):
        sobjs, _, _ = load_dwg(SYS)
        notes = []
        for o in sobjs:
            if o.get("entity") not in ("TEXT", "MTEXT"):
                continue
            tx = clean_text(o.get("text_value") or o.get("text") or "")
            if tx:
                pt = o.get("ins_pt") or [0, 0, 0]
                notes.append((tx, pt[0], pt[1]))
        anchors = [(t, x, y) for t, x, y in notes if "采集器" in t or "数据机房" in t or t == "UPS"]
        for t, x, y in notes:
            m = re.search(r"×\s*(\d+)|\*(\d+)", t)
            if not m:
                continue
            n = int(m.group(1) or m.group(2))
            if n <= 1 or n > 2000:
                continue
            host = ""
            if anchors:
                host = min(anchors, key=lambda a: (a[1] - x) ** 2 + (a[2] - y) ** 2)[0]
            branches.append({"host": host, "item": re.sub(r"[×\*]\d+", "", t).strip() or t, "n": n, "source": "sim"})

    # 去重分支
    uniq_b, seen_b = [], set()
    for b in branches:
        k = (b["host"], b["item"], b["n"])
        if k in seen_b:
            continue
        seen_b.add(k)
        uniq_b.append(b)

    for fl, room, cat, name in named:
        assets.append({
            "room": room if room != "三层ECC" else "3层操作间",
            "floor": "3F" if fl == "楼" else fl,
            "cabinet": "",
            "sn": name,
            "status": "在架",
            "uStart": 0, "uEnd": 0, "uH": 0,
            "cat": cat,
            "type": cat,
            "brand": "循安科技",
            "model": name,
            "kw": 0.2,
            "psu": 1, "weight": 0,
            "net": "", "biz": "动环",
            "org": "未登记", "owner": "",
            "source": "sim",
        })

    md.append("## 有名字的监测设备")
    md.append("")
    for _, room, cat, name in named:
        md.append("- %s · %s · %s" % (room, cat, name))
    md.append("")
    md.append("## 系统图分支倍数")
    md.append("")
    md.append("这些只写了数量，没有单独坐标，挂在就近的采集器或机房名下。读数是模拟。")
    md.append("")
    for b in uniq_b:
        md.append("- %s：%s × %d" % (b["host"] or "未归类", b["item"], b["n"]))
    md.append("")
    open(LEDGER, "w").write("\n".join(md) + "\n")

    dates = SERIES_DATES
    base_vals = SERIES_BASE

    def series_dev(name, loc, scale):
        return {
            "name": name,
            "loc": loc,
            "source": "sim",
            "points": [{
                "name": "模拟值",
                "source": "sim",
                "vals": [round(v * scale, 2) for v in base_vals],
            }],
        }

    halls = [r for r in rooms_out if r["type"] == "主机房"]
    series = {
        "列头柜": {"dates": dates, "source": "sim", "devices": [
            series_dev(r["code"] + "-列头柜", r["name"], 0.8 + i * 0.05) for i, r in enumerate(halls)
        ]},
        "蓄电池": {"dates": dates, "source": "sim", "devices": [
            series_dev(fl + "-电池", fl + " 电池间", 0.9) for fl in ORDER
        ]},
        "空调": {"dates": dates, "source": "sim", "devices": [
            series_dev(r["code"] + "-列间空调", r["name"], 0.7) for r in halls
        ]},
        "氢气": {"dates": dates, "source": "sim", "devices": [
            series_dev(fl + "-氢气", fl + " 电池间", 0.2) for fl in ORDER
        ]},
    }
    th = []
    for r in halls:
        for aisle, temp in (("冷通道", r["temp"] - 1.5), ("热通道", r["temp"] + 2.5)):
            th.append({
                "name": r["code"] + aisle,
                "room": r["name"],
                "floor": r["floor"],
                "aisle": "cold" if aisle == "冷通道" else "hot",
                "temp": round(temp, 1),
                "humi": r["humi"],
                "source": "sim",
            })
    alarms = []
    for i, r in enumerate(rooms_out):
        if r["alarmN"] <= 0:
            continue
        alarms.append({
            "level": "次要" if i % 4 else "预警",
            "what": "温度偏高" if i % 2 == 0 else "湿度偏高",
            "floor": r["floor"],
            "room": r["name"],
            "zone": r["zone"],
            "obj": r["code"] or r["name"],
            "time": "2026-10-23 09:%02d:00" % (i % 60),
            "state": "未确认",
            "accept": "",
            "confirm": "",
            "recover": "",
            "handler": "",
            "kind": "过高报警",
            "cat": "环境",
            "dur": "2小时",
            "rtime": "",
            "tv": r["temp"] + 3,
            "cv": r["temp"],
            "source": "sim",
        })
    by_floor = defaultdict(int)
    by_room = defaultdict(int)
    by_cat = defaultdict(int)
    by_brand = defaultdict(int)
    by_status = defaultdict(int)
    by_org = defaultdict(int)
    by_biz = defaultdict(int)
    power = 0
    weight = 0
    for a in assets:
        by_floor[a["floor"]] += 1
        by_room[a["room"]] += 1
        by_cat[a["cat"]] += 1
        by_brand[a["brand"] or "未登记"] += 1
        by_status[a["status"] or "未登记"] += 1
        by_org[a["org"] or "未登记"] += 1
        by_biz[a["biz"] or "未登记"] += 1
        power += a["kw"]
        weight += a.get("weight") or 0
    cab_total = sum(f["cabCount"] for f in floors_out)
    old = json.load(open(DCIM))
    doc = old
    doc["meta"] = {
        "system": "数据中心动环管理系统",
        "systemEn": "DATA CENTER DCIM PLATFORM",
        "project": "北京公安局数据中心",
        "exportDate": "2026-10-23",
        "assetTotal": len(assets),
        "alarmTotal": len(alarms),
        "roomTotal": len(rooms_out),
        "floorTotal": 3,
        "cabTotal": cab_total,
        "cabUsed": cab_total,
        "cabFree": 0,
        "realSeed": len(assets),
        "dataMode": "sim",
        "area": round(sum(f["w"] * f["d"] for f in floors_out), 0),
        "note": "楼层、房间、机柜和监测点来自 2026-10-23 平面图与动环系统图。图纸没有服务器清单，柜内服务器是模拟，source=sim。",
    }
    doc["floors"] = floors_out
    doc["rooms"] = rooms_out
    doc["assets"] = assets
    doc["assetStats"] = {
        "total": len(assets),
        "byRoom": [{"name": k, "value": v} for k, v in sorted(by_room.items())],
        "byCat": [{"name": k, "value": v} for k, v in sorted(by_cat.items())],
        "byType": [{"name": k, "value": v} for k, v in sorted(by_cat.items())],
        "byBrand": [{"name": k, "value": v} for k, v in sorted(by_brand.items())],
        "byStatus": [{"name": k, "value": v} for k, v in sorted(by_status.items())],
        "byOrg": [{"name": k, "value": v} for k, v in sorted(by_org.items())],
        "byFloor": [{"name": k, "value": v} for k, v in sorted(by_floor.items())],
        "byNet": [],
        "byBiz": [{"name": k, "value": v} for k, v in sorted(by_biz.items())],
        "powerSum": round(power, 1),
        "weightSum": weight,
        "bizCount": len(by_biz),
    }
    doc["alarms"] = alarms
    doc["alarmStats"] = {
        "total": len(alarms),
        "byLevel": [{"name": "次要", "value": sum(1 for a in alarms if a["level"] == "次要")},
                    {"name": "预警", "value": sum(1 for a in alarms if a["level"] == "预警")}],
        "source": "sim",
    }
    doc["series"] = series
    doc["thDevices"] = th
    doc["branches"] = uniq_b
    doc["capacity"]["cabU"] = cab_total * 42
    doc["capacity"]["source"] = "sim"
    verify_sim_servers(assets, rooms_out)
    text = json.dumps(doc, ensure_ascii=False, separators=(",", ":"))
    for bad in ("205机房", "407机房", "BJDC-205"):
        if bad in text:
            raise SystemExit("出现了不该恢复的旧房间: " + bad)
    open(DCIM, "w").write(text)
    srv_n = [a for a in assets if a["cat"] == "服务器"]
    per = defaultdict(int)
    for a in srv_n:
        per[a["cabinet"]] += 1
    counts = list(per.values())
    print("floors", [(f["floor"], f["rooms"], f["cabCount"], f["w"], f["d"], len(f["walls"]), len(f["outline"])) for f in floors_out])
    print("rooms", len(rooms_out), "assets", len(assets), "servers", len(srv_n), "alarms", len(alarms), "branches", len(uniq_b))
    print("per-cab", min(counts), max(counts), "cabs", len(counts))
    print("ledger", LEDGER)


if __name__ == "__main__":
    main()
