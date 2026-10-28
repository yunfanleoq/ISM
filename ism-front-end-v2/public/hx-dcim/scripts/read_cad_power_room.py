#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""从《配电系统图.dwg》读 1～5 层配电室。只读 DWG，DXF 只写 /tmp。

  python3 read_cad_power_room.py list
  python3 read_cad_power_room.py room 3A1
  python3 read_cad_power_room.py emit

标准文件 MD5 必须是 c58f0fb1dfacc088b8f6f470ffd2b6d5，不对就停。
"""
import hashlib
import json
import subprocess
import sys
from collections import defaultdict
from pathlib import Path

from cad_dxf import load_model
from cad_room_layout import OUT_DY, layout_room
from cad_rooms import build_rooms, summary_rows

MD5 = "c58f0fb1dfacc088b8f6f470ffd2b6d5"
VERSION = "20261028cad76"
DWG_REL = Path("docs/refs/gongan-dashboard/group-files/配电系统图.dwg")
DXF = Path("/tmp/hx-power-system.dxf")


def repo_root():
    for p in Path(__file__).resolve().parents:
        if (p / DWG_REL).exists():
            return p
    sys.exit("找不到配电系统图.dwg，停止")


def ensure_dxf(root):
    dwg = root / DWG_REL
    digest = hashlib.md5(dwg.read_bytes()).hexdigest()
    if digest != MD5:
        sys.exit("配电系统图.dwg MD5=%s，与标准 %s 不符，停止" % (digest, MD5))
    copy = root / "docs/refs/hangxin-screen-drawings/配电系统图.dwg"
    if copy.exists():
        digest2 = hashlib.md5(copy.read_bytes()).hexdigest()
        if digest2 != MD5:
            sys.exit("副本配电系统图.dwg MD5 不符，停止")
    if not DXF.exists() or DXF.stat().st_mtime < dwg.stat().st_mtime:
        subprocess.check_call(["dwg2dxf", "-y", "-o", str(DXF), str(dwg)], stdout=subprocess.DEVNULL)
    return dwg


def load_rooms(root):
    ensure_dxf(root)
    model = load_model(str(DXF))
    rooms = build_rooms(model)
    return [r for r in rooms if r.get("floor") in (1, 2, 3, 4, 5) or str(r["code"]).startswith("ECC")]


def _md(room, scene_keys):
    lines = [
        "# 配电室 %s 一次图网表" % room["code"],
        "",
        "- 标准：`docs/refs/gongan-dashboard/group-files/配电系统图.dwg`（MD5 `%s`）" % MD5,
        "- 读法：`read_cad_power_room.py` → LibreDWG `dwg2dxf` 到 `/tmp`，不入库。",
        "- 标题坐标：(%s, %s)，刀云范围 %s，楼层 %s" % (
            room["title"]["x"], room["title"]["y"], room["bbox"], room.get("floor")),
        "- 刀 %d，线段 %d，变压器符号团 %d" % (len(room["knives"]), len(room["segs"]), len(room["xfmrs"])),
        "",
        "柜号来自该室文字，不套用 5A1 的 401/411/U11。没有场景测点的刀画静态符号。",
        "",
        "| 刀 | 角色 | 电流 | 标注 | CAD x | 来源 | 场景 |",
        "|---|---|---|---|---|---|---|",
    ]
    bound = scene_keys or {}
    static = []
    for k in sorted(room["knives"], key=lambda n: (-n["y"], n["x"])):
        src = "文字补" if k.get("synthetic") else ("未命名" if k.get("unnamed") else "刀块")
        key = k["cad"].split("#")[0]
        hit = bound.get(k["cad"]) or bound.get(key)
        if not hit and key == "445":
            hit = bound.get("T1-445") or bound.get("T2-445")
        scene = "有测点" if hit else "静态"
        if not hit:
            static.append(k["cad"])
        lines.append("| %s | %s | %s | %s | %.1f | %s | %s |" % (
            k["cad"], k.get("role") or "", k.get("amp") or "", (k.get("use") or "").replace("|", "/"),
            k["x"], src, scene))
    buses = [s for s in room["segs"] if s.get("color") == 1 and s.get("horiz") and s.get("len", 0) > 15]
    lines += ["", "## 红色母线", ""]
    if not buses:
        lines.append("（本室没有长红母线，连接以默认色正交线为准）")
    for s in buses:
        lines.append("- y=%.1f  x %.1f → %.1f  长 %.1f" % (
            (s["y1"] + s["y2"]) / 2, min(s["x1"], s["x2"]), max(s["x1"], s["x2"]), s["len"]))
    lines += ["", "## 静态刀（场景无测点）", "", "、".join(static) if static else "无", ""]
    if scene_keys:
        cad_ids = set()
        for k in room["knives"]:
            cad_ids.add(k["cad"])
            cad_ids.add(k["cad"].split("#")[0])
        scene_ids = set(scene_keys)
        if "445" in cad_ids and ("T1-445" in scene_ids or "T2-445" in scene_ids):
            scene_ids.discard("T1-445")
            scene_ids.discard("T2-445")
            cad_ids.discard("445")
        only_scene = sorted(scene_ids - cad_ids)
        only_cad = sorted(cad_ids - scene_ids)
        lines += ["", "## 和原场景刀号", ""]
        lines.append("- 场景有、这间 CAD 没有：%s" % ("、".join(only_scene) if only_scene else "无"))
        lines.append("- CAD 有、场景没有（静态）：%s" % ("、".join(only_cad) if only_cad else "无"))
        lines.append("")
    for note in room.get("droppedNotes") or []:
        lines.append("- " + note)
    if room.get("extraKnivesDropped"):
        lines.append("- 同名另一团未并入的刀约 %d 把" % room["extraKnivesDropped"])
    lines.append("")
    return "\n".join(lines)


def _scene_index(root):
    path = root / "ism-front-end-v2/public/hx-dcim/data/hx-scenes.json"
    if not path.exists():
        return {}
    data = json.loads(path.read_text(encoding="utf-8"))
    out = {}
    for pg in data.get("pages") or []:
        name = pg.get("name") or ""
        if not name.startswith("配电室"):
            continue
        code = name.replace("配电室", "")
        keys = {}
        for n in pg.get("nodes") or []:
            k = n.get("k") or ""
            if not k.startswith(name + "_"):
                continue
            cad = k[len(name) + 1:].strip("_").replace("_", "-")
            if n.get("p") or n.get("dev") or n.get("k"):
                keys[cad] = {"k": k, "p": n.get("p") or "", "dev": n.get("dev") or ""}
        out[code] = keys
    return out


def _touch(layout, cad):
    k = next((x for x in layout["knives"] if x["cad"] == cad), None)
    if not k:
        return False
    tx = k["x"] + 30
    ty = k["y"]
    for w in layout["wires"]:
        pts = w["pts"]
        for i in range(1, len(pts)):
            a, b = pts[i - 1], pts[i]
            if abs(a[0] - b[0]) <= 1.6:
                cx = (a[0] + b[0]) / 2
                if abs(cx - tx) <= 4 and min(a[1], b[1]) - 3 <= ty <= max(a[1], b[1]) + 3:
                    return True
            if abs(a[1] - b[1]) <= 1.6:
                cy = (a[1] + b[1]) / 2
                if abs(cy - ty) <= 4 and min(a[0], b[0]) - 3 <= tx <= max(a[0], b[0]) + 3:
                    return True
    return False


def cmd_list(rooms):
    print("floor code knives segs in cap tie xfmr unnamed  title")
    for row in summary_rows(rooms):
        print("%s %s %s %s %s %s %s %s %s  %s" % (
            row["floor"], row["code"], row["knives"], row["segs"], row["in"], row["cap"],
            row["tie"], row["xfmr"], row["unnamed"], row["title"]))


def cmd_room(rooms, code):
    hit = [r for r in rooms if r["code"] == code]
    if not hit:
        sys.exit("没有这间: " + code)
    r = hit[0]
    print(_md(r, {}))
    lay = layout_room(r, VERSION)
    if not lay:
        return
    print("\n## 布局自检")
    for cad in lay["meta"]["incomers"] + lay["meta"]["caps"]:
        print(" ", cad, "上口接通" if _touch(lay, cad) else "上口未接")


def _patch_arch(root, codes):
    path = root / "ism-front-end-v2/public/hx-dcim/data/hx-arch-tree.json"
    tree = json.loads(path.read_text(encoding="utf-8"))
    want = {c: "配电室" + c for c in codes if not c.endswith("b") and not c.startswith("ECC")}
    existing = {}

    def walk(node, floor):
        name = node.get("name") or ""
        kids = node.get("children") or []
        fl = floor
        if name in ("F1", "F2", "F3", "F4", "F5"):
            fl = int(name[1])
        if name.startswith("配电室") and not kids:
            existing[name] = node
            if fl:
                node["_floor"] = fl
        for ch in kids:
            walk(ch, fl)

    for b in tree.get("branches") or []:
        walk(b, None)
    # 补楼层里有 CAD、树上没有的室
    floors = {}

    def find_floor(node):
        if node.get("name") in ("F1", "F2", "F3", "F4", "F5"):
            floors[node["name"]] = node
        for ch in node.get("children") or []:
            find_floor(ch)

    for b in tree.get("branches") or []:
        find_floor(b)
    for code in sorted(want):
        name = want[code]
        if name in existing:
            node = existing[name]
            if not node.get("pageId"):
                node["pageId"] = "cad-" + code
            continue
        fl = "F" + code[0]
        host = floors.get(fl)
        if not host:
            continue
        host.setdefault("children", []).append({"name": name, "children": [], "pageId": "cad-" + code})
    for fl, host in floors.items():
        host["children"] = sorted(host.get("children") or [], key=lambda n: n.get("name") or "")
    # 去掉临时字段
    def clean(node):
        node.pop("_floor", None)
        for ch in node.get("children") or []:
            clean(ch)
    for b in tree.get("branches") or []:
        clean(b)
    path.write_text(json.dumps(tree, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")


def _patch_ecc(root):
    path = root / "ism-front-end-v2/public/hx-dcim/data/hx-arch-tree.json"
    tree = json.loads(path.read_text(encoding="utf-8"))

    def walk(node):
        name = node.get("name") or ""
        if name in ("ECC1", "ECC2") and not (node.get("children") or []):
            node["pageId"] = "cad-" + name
        for ch in node.get("children") or []:
            walk(ch)

    for b in tree.get("branches") or []:
        walk(b)
    path.write_text(json.dumps(tree, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")


def cmd_emit_codes(root, rooms, codes):
    """只重写点名的室，不动 5A1，也不重画其余配电室。"""
    scenes = _scene_index(root)
    doc = root / "docs/refs/hangxin-screen-drawings/check/cad-rooms"
    data = root / "ism-front-end-v2/public/hx-dcim/data/hx-power-rooms"
    doc.mkdir(parents=True, exist_ok=True)
    by = {r["code"]: r for r in rooms}
    index_path = data / "index.json"
    index = json.loads(index_path.read_text(encoding="utf-8"))
    for code in codes:
        if code == "5A1":
            print("skip 5A1")
            continue
        room = by.get(code)
        if not room:
            sys.exit("没有这间: " + code)
        scene = scenes.get(code) or {}
        (doc / (code + ".md")).write_text(_md(room, scene), encoding="utf-8")
        lay = layout_room(room, VERSION)
        if not lay:
            sys.exit(code + " 未画出")
        lay["meta"]["sceneBound"] = [k["cad"] for k in lay["knives"] if k["cad"] in scene or k["cad"].split("#")[0] in scene]
        lay["meta"]["staticKnives"] = [k["cad"] for k in lay["knives"] if k["cad"] not in scene and k["cad"].split("#")[0] not in scene]
        (data / (code + ".json")).write_text(json.dumps(lay, ensure_ascii=False, separators=(",", ":")), encoding="utf-8")
        for item in index.get("rooms") or []:
            if item.get("code") == code:
                item["knives"] = lay["meta"]["knifeCount"]
        bad = [c for c in lay["meta"]["incomers"] + lay["meta"]["caps"] if not _touch(lay, c)]
        print(code, "knives", lay["meta"]["knifeCount"], "wires", len(lay["wires"]), "untouched", bad)
    index["version"] = VERSION
    index_path.write_text(json.dumps(index, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print("index version", VERSION)


def cmd_emit(root, rooms):
    scenes = _scene_index(root)
    doc = root / "docs/refs/hangxin-screen-drawings/check/cad-rooms"
    data = root / "ism-front-end-v2/public/hx-dcim/data/hx-power-rooms"
    doc.mkdir(parents=True, exist_ok=True)
    data.mkdir(parents=True, exist_ok=True)
    index = []
    report = []
    drawn = []
    for room in rooms:
        code = room["code"]
        scene = scenes.get(code) or {}
        (doc / (code + ".md")).write_text(_md(room, scene), encoding="utf-8")
        lay = layout_room(room, VERSION)
        if not lay or code == "5A1":
            status = "基线不重画" if code == "5A1" else "未画"
            report.append((room, status, []))
            if code == "5A1" and lay:
                bad = [c for c in lay["meta"]["incomers"] + lay["meta"]["caps"] if not _touch(lay, c)]
                (doc / "5A1-autolayout-check.json").write_text(json.dumps({
                    "knives": lay["meta"]["knifeCount"],
                    "incomers": lay["meta"]["incomers"],
                    "caps": lay["meta"]["caps"],
                    "untouched": bad,
                }, ensure_ascii=False, indent=2), encoding="utf-8")
            continue
        # 场景测点写进布局，运行时再按名绑定；这里只记哪些是静态
        lay["meta"]["sceneBound"] = [k["cad"] for k in lay["knives"] if k["cad"] in scene or k["cad"].split("#")[0] in scene]
        lay["meta"]["staticKnives"] = [k["cad"] for k in lay["knives"] if k["cad"] not in scene and k["cad"].split("#")[0] not in scene]
        (data / (code + ".json")).write_text(json.dumps(lay, ensure_ascii=False, separators=(",", ":")), encoding="utf-8")
        (doc / (code + ".layout.json")).write_text(json.dumps({
            "knifeCount": lay["meta"]["knifeCount"],
            "incomers": lay["meta"]["incomers"],
            "caps": lay["meta"]["caps"],
            "ties": lay["meta"]["ties"],
            "untouched": [c for c in lay["meta"]["incomers"] + lay["meta"]["caps"] if not _touch(lay, c)],
        }, ensure_ascii=False, indent=2), encoding="utf-8")
        index.append({
            "code": code,
            "name": "配电室" + code,
            "floor": room.get("floor"),
            "file": code + ".json",
            "knives": lay["meta"]["knifeCount"],
        })
        bad = [c for c in lay["meta"]["incomers"] + lay["meta"]["caps"] if not _touch(lay, c)]
        report.append((room, "已画", bad))
        drawn.append(code)
    index.sort(key=lambda r: (r["floor"] or 9, r["code"]))
    (data / "index.json").write_text(json.dumps({
        "version": VERSION, "rooms": index,
    }, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    _patch_arch(root, [r["code"] for r in rooms if r["code"] != "5A1"])
    _patch_ecc(root)
    # 已有场景页的室保留原 pageId
    lines = ["# CAD 配电室清单", "", "版本 `%s`。5A1 只出网表，画面仍用已验证布局。" % VERSION, "",
             "| 楼层 | 室 | 刀 | 进线 | 补偿 | 状态 | 上口未接 |",
             "|---|---|---|---|---|---|---|"]
    for room, status, bad in report:
        meta_in = sum(1 for k in room["knives"] if k.get("role") == "in")
        meta_cap = sum(1 for k in room["knives"] if k.get("role") == "cap")
        lines.append("| %s | %s | %d | %d | %d | %s | %s |" % (
            room.get("floor"), room["code"], len(room["knives"]), meta_in, meta_cap, status,
            "、".join(bad) if bad else ""))
    (doc / "index.md").write_text("\n".join(lines) + "\n", encoding="utf-8")
    print("\n".join(lines))
    print("drawn", len(drawn))


def main():
    root = repo_root()
    rooms = load_rooms(root)
    cmd = sys.argv[1] if len(sys.argv) > 1 else "list"
    if cmd == "list":
        cmd_list(rooms)
    elif cmd == "room":
        cmd_room(rooms, sys.argv[2])
    elif cmd == "emit":
        cmd_emit(root, rooms)
    elif cmd == "emit-codes":
        if len(sys.argv) < 3:
            sys.exit("用法: emit-codes 1A1 1A2 ...")
        cmd_emit_codes(root, rooms, sys.argv[2:])
    else:
        sys.exit("用法: list | room CODE | emit | emit-codes CODE...")


if __name__ == "__main__":
    main()
