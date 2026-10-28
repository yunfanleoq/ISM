#!/usr/bin/env python3
"""从航信组态库 display_model_layer 导出 hx-dcim 用的 hx-scenes.json。

数据源默认：ism_server_user/data/db/ism-hangxin-20261023.db
模型 uuid 与 public/hx-dcim/data/hx-scenes.json 中 model 字段一致。

用法：
  python3 scripts/hangxin_export_hx_scenes.py --list
  python3 scripts/hangxin_export_hx_scenes.py --export-power --export-columns
  python3 scripts/hangxin_export_hx_scenes.py --export-power -o ism-front-end-v2/public/hx-dcim/data/hx-scenes.json
  python3 scripts/hangxin_export_hx_scenes.py --manifest docs/hx-scenes-pending-power.json

components / layer 列为 base64(JSON)。配电室页走 ism-view-png-image；列头柜走 view-svg-sld-load-switch。
"""
from __future__ import annotations

import argparse
import base64
import json
import os
import re
import sqlite3
import sys
from collections import Counter
from datetime import datetime, timezone
from typing import Any

ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
DEFAULT_DB = os.path.join(ROOT, "ism_server_user", "data", "db", "ism-hangxin-20261023.db")
DEFAULT_OUT = os.path.join(ROOT, "ism-front-end-v2", "public", "hx-dcim", "data", "hx-scenes.json")
DEFAULT_ARCH = os.path.join(ROOT, "ism-front-end-v2", "public", "hx-dcim", "data", "hx-arch-tree.json")
DEFAULT_POWER_INDEX = os.path.join(
    ROOT, "ism-front-end-v2", "public", "hx-dcim", "data", "hx-power-page-index.json"
)
MODEL_ID = "b620dc75-1da2-351b-b89e-05b5cb608886"

POWER_PAGE_RE = re.compile(r"^配电室\d[AB]\d$")
COLUMN_PAGE_RE = re.compile(r"^(\d[AB]\d)数据机房_([A-K]列头)$")


def decode_blob(raw: str) -> Any:
    if not raw:
        return {}
    text = raw.strip()
    if not text:
        return {}
    if text[0] in "{[":
        return json.loads(text)
    return json.loads(base64.b64decode(text))


def key_from_point(point_name: str, node_type: str) -> str:
    if not point_name:
        return ""
    if node_type in ("S", "I"):
        if point_name.endswith("开关状态"):
            return point_name[: -len("开关状态")]
        if "_" in point_name:
            return point_name.rsplit("_", 1)[0]
    return point_name


def skip_capacitor(det: dict) -> bool:
    return "电容补偿" in json.dumps(det, ensure_ascii=False)


def node_pos(c: dict, style: dict) -> tuple[float, float, float, float, int]:
    pos = c.get("position") or style.get("position") or {}
    x = float(pos.get("x") or 0)
    y = float(pos.get("y") or 0)
    w = float((c.get("size") or {}).get("width") or pos.get("w") or 0)
    h = float((c.get("size") or {}).get("height") or pos.get("h") or 0)
    z = int(style.get("zIndex") or 0)
    return x, y, w, h, z


def convert_cells(cells: list[dict], layer: dict) -> list[dict]:
    nodes: list[dict] = []
    for c in cells:
        shape = c.get("shape") or ""
        det = (c.get("data") or {}).get("detail") or {}
        if skip_capacitor(det):
            continue
        style = det.get("style") or {}
        x, y, w, h, z = node_pos(c, style)
        cond = (det.get("animate") or {}).get("condition") or {}

        if shape == "view-svg-line":
            nodes.append(
                {
                    "t": "L",
                    "x": x,
                    "y": y,
                    "w": w,
                    "h": h,
                    "a": float(c.get("angle") or 0),
                    "c": style.get("foreColor") or "#ffffff",
                    "sw": float(style.get("borderWidth") or 2),
                    "z": z,
                    "k": "",
                }
            )
        elif shape == "view-svg-text":
            nodes.append(
                {
                    "t": "T",
                    "x": x,
                    "y": y,
                    "w": w,
                    "h": h,
                    "s": style.get("text") or "",
                    "c": style.get("foreColor") or "#fff",
                    "fs": style.get("fontSize") or 18,
                    "z": z,
                }
            )
        elif shape == "dv-border-box12":
            nodes.append({"t": "B", "x": x, "y": y, "w": w, "h": h, "z": z})
        elif shape == "view-svg-sld-out":
            nodes.append({"t": "O", "x": x, "y": y, "w": w, "h": h, "z": z})
        elif shape == "view-svg-sld-load-switch":
            p = cond.get("dataName") or ""
            if not p:
                continue
            k = key_from_point(p, "S")
            nodes.append(
                {
                    "t": "S",
                    "x": x,
                    "y": y,
                    "w": w,
                    "h": h,
                    "c": "#00ff00",
                    "k": k,
                    "p": p,
                    "dev": cond.get("DeviceName") or "",
                    "op": cond.get("operator") or "==",
                    "val": str(cond.get("OperatorValue") if cond.get("OperatorValue") is not None else "0"),
                    "z": z,
                    "decor": 1,
                }
            )
        elif shape == "ism-view-png-image":
            p = cond.get("dataName") or ""
            if not p:
                continue
            img = style.get("imageURL") or ""
            st = "0" if ".state0" in img else "1"
            k = key_from_point(p, "I")
            nodes.append(
                {
                    "t": "I",
                    "x": x,
                    "y": y,
                    "w": w,
                    "h": h,
                    "st": st,
                    "k": k,
                    "p": p,
                    "dev": cond.get("DeviceName") or "",
                    "op": cond.get("operator") or "==",
                    "val": str(cond.get("OperatorValue") if cond.get("OperatorValue") is not None else ""),
                    "z": z,
                }
            )
    return nodes


def page_bounds(nodes: list[dict], layer: dict) -> tuple[int, int, int]:
    if not nodes:
        return 3150, 1400, 0
    max_x = max(n["x"] + n["w"] for n in nodes)
    max_y = max(n["y"] + n["h"] for n in nodes)
    min_x = min(n["x"] for n in nodes)
    x0 = max(0, int(min_x) - 20) if min_x > 50 else 0
    w = max(int(max_x) + 40, int(layer.get("width") or 0) or 3150)
    h = max(int(max_y) + 40, int(layer.get("height") or 0) or 1400)
    return w, h, x0


def classify_page(page_name: str) -> tuple[str | None, dict]:
    m = COLUMN_PAGE_RE.match(page_name)
    if m:
        code, cab = m.group(1), m.group(2)
        return "column", {"group": f"{code}数据机房", "cab": cab}
    if POWER_PAGE_RE.match(page_name):
        return "power", {"group": "配电一次图", "cab": ""}
    return None, {}


def load_arch_power_leaves(path: str) -> list[str]:
    if not os.path.isfile(path):
        return []
    with open(path, encoding="utf-8") as f:
        spec = json.load(f)
    leaves: list[str] = []

    def walk(node: dict) -> None:
        name = node.get("name") or ""
        kids = node.get("children") or []
        if not kids and POWER_PAGE_RE.match(name):
            leaves.append(name)
            return
        for ch in kids:
            walk(ch)

    for b in spec.get("branches") or []:
        if b.get("name") == "配电室":
            walk(b)
    return sorted(set(leaves))


def fetch_layers(conn: sqlite3.Connection, model_id: str) -> list[sqlite3.Row]:
    conn.row_factory = sqlite3.Row
    return conn.execute(
        """
        SELECT page_name, page_id, is_home, components, layer
        FROM display_model_layer
        WHERE model_id = ? AND deleted_at IS NULL
        ORDER BY page_name
        """,
        (model_id,),
    ).fetchall()


def build_page(row: sqlite3.Row) -> dict | None:
    kind, extra = classify_page(row["page_name"])
    if not kind:
        return None
    comp = decode_blob(row["components"])
    layer = decode_blob(row["layer"])
    cells = comp.get("cells") if isinstance(comp, dict) else []
    if not isinstance(cells, list):
        return None
    nodes = convert_cells(cells, layer)
    if not nodes:
        return None
    w, h, x0 = page_bounds(nodes, layer)
    page = {
        "id": row["page_id"],
        "name": row["page_name"],
        "kind": kind,
        "group": extra.get("group") or "",
        "cab": extra.get("cab") or "",
        "home": 1 if row["is_home"] else 0,
        "w": w,
        "h": h,
        "x0": x0,
        "bg": layer.get("backgroundImage") or "/static/ISM/page/background/background07.state0.jpg",
        "nodes": nodes,
    }
    return page


def merge_pages(existing: list[dict], fresh: list[dict]) -> list[dict]:
    by_id = {p["id"]: p for p in existing}
    for p in fresh:
        by_id[p["id"]] = p
    pages = list(by_id.values())
    pages.sort(key=lambda p: (p.get("kind") or "", p.get("name") or ""))
    return pages


def cmd_list(conn: sqlite3.Connection, model_id: str, arch_path: str) -> int:
    rows = fetch_layers(conn, model_id)
    power_in_db = [r["page_name"] for r in rows if POWER_PAGE_RE.match(r["page_name"])]
    col_in_db = [r["page_name"] for r in rows if COLUMN_PAGE_RE.match(r["page_name"])]
    arch_leaves = load_arch_power_leaves(arch_path)
    pending = sorted(set(arch_leaves) - set(power_in_db))

    print(f"model_id: {model_id}")
    print(f"DB 配电室一次图页 ({len(power_in_db)}): {', '.join(power_in_db) or '（无）'}")
    print(f"DB 列头柜页 ({len(col_in_db)}): {len(col_in_db)} 页")
    print(f"架构树配电室叶子 ({len(arch_leaves)}), 尚无 DB 页 ({len(pending)}):")
    for name in pending:
        print(f"  - {name}")
    other = [
        r["page_name"]
        for r in rows
        if not POWER_PAGE_RE.match(r["page_name"]) and not COLUMN_PAGE_RE.match(r["page_name"])
    ]
    if other:
        print(f"其它组态页 ({len(other)}): {', '.join(other[:8])}{'…' if len(other) > 8 else ''}")
    return 0


def cmd_manifest(conn: sqlite3.Connection, model_id: str, arch_path: str, out_path: str) -> int:
    rows = fetch_layers(conn, model_id)
    power_in_db = sorted(r["page_name"] for r in rows if POWER_PAGE_RE.match(r["page_name"]))
    arch_leaves = load_arch_power_leaves(arch_path)
    doc = {
        "generatedAt": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "model": model_id,
        "db": DEFAULT_DB,
        "exportedPower": power_in_db,
        "archPowerLeaves": arch_leaves,
        "pendingPower": sorted(set(arch_leaves) - set(power_in_db)),
        "note": "pendingPower 需从现场 MySQL/制作画面文件.zip 补 display_model_layer 后再跑 --export-power",
    }
    os.makedirs(os.path.dirname(out_path) or ".", exist_ok=True)
    with open(out_path, "w", encoding="utf-8") as f:
        json.dump(doc, f, ensure_ascii=False, indent=2)
    print(f"Wrote manifest: {out_path}")
    print(f"exported {len(power_in_db)}, pending {len(doc['pendingPower'])}")
    return 0


def cmd_power_index(
    conn: sqlite3.Connection,
    model_id: str,
    out_path: str,
    arch_path: str,
    patch_arch: bool,
    db_path: str,
) -> int:
    rows = fetch_layers(conn, model_id)
    pages = [
        {"name": r["page_name"], "pageId": r["page_id"]}
        for r in rows
        if POWER_PAGE_RE.match(r["page_name"])
    ]
    pages.sort(key=lambda p: p["name"])
    doc = {
        "generatedAt": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "model": model_id,
        "db": os.path.abspath(db_path),
        "pages": pages,
    }
    os.makedirs(os.path.dirname(out_path) or ".", exist_ok=True)
    with open(out_path, "w", encoding="utf-8") as f:
        json.dump(doc, f, ensure_ascii=False, indent=2)
    print(f"Wrote power index ({len(pages)} pages): {out_path}")

    if patch_arch and os.path.isfile(arch_path):
        by_name = {p["name"]: p["pageId"] for p in pages}
        with open(arch_path, encoding="utf-8") as f:
            spec = json.load(f)

        def walk(node: dict) -> None:
            name = node.get("name") or ""
            kids = node.get("children") or []
            if not kids and POWER_PAGE_RE.match(name):
                pid = by_name.get(name)
                if pid:
                    node["pageId"] = pid
                elif "pageId" in node:
                    del node["pageId"]
                return
            for ch in kids:
                walk(ch)

        for b in spec.get("branches") or []:
            if b.get("name") == "配电室":
                walk(b)
        with open(arch_path, "w", encoding="utf-8") as f:
            json.dump(spec, f, ensure_ascii=False, indent=2)
        print(f"Patched arch tree pageId on {len(by_name)} power leaves: {arch_path}")
    return 0


def cmd_export(
    conn: sqlite3.Connection,
    model_id: str,
    out_path: str,
    export_power: bool,
    export_columns: bool,
    merge: bool,
) -> int:
    rows = fetch_layers(conn, model_id)
    fresh: list[dict] = []
    for row in rows:
        kind, _ = classify_page(row["page_name"])
        if kind == "power" and not export_power:
            continue
        if kind == "column" and not export_columns:
            continue
        page = build_page(row)
        if page:
            fresh.append(page)

    if merge and os.path.isfile(out_path):
        with open(out_path, encoding="utf-8") as f:
            doc = json.load(f)
    else:
        doc = {"model": model_id, "title": "一次图界面", "pages": []}

    doc["model"] = model_id
    doc["pages"] = merge_pages(doc.get("pages") or [], fresh)

    os.makedirs(os.path.dirname(out_path) or ".", exist_ok=True)
    with open(out_path, "w", encoding="utf-8") as f:
        json.dump(doc, f, ensure_ascii=False, separators=(",", ":"))

    kinds = Counter(p["kind"] for p in doc["pages"])
    print(f"Wrote {out_path}")
    print(f"pages: power={kinds.get('power', 0)} column={kinds.get('column', 0)} (this run +{len(fresh)})")
    return 0


def main() -> int:
    parser = argparse.ArgumentParser(description="Export hangxin display_model_layer → hx-scenes.json")
    parser.add_argument("--db", default=DEFAULT_DB, help="SQLite 航信库路径")
    parser.add_argument("-o", "--output", default=DEFAULT_OUT, help="输出 hx-scenes.json")
    parser.add_argument("--model", default=MODEL_ID, help="display_model uid")
    parser.add_argument("--arch", default=DEFAULT_ARCH, help="hx-arch-tree.json（待导出清单）")
    parser.add_argument("--list", action="store_true", help="列出 DB/架构树配电室覆盖情况")
    parser.add_argument("--manifest", metavar="PATH", help="写出 pending 清单 JSON")
    parser.add_argument("--export-power", action="store_true", help="导出配电室一次图")
    parser.add_argument("--export-columns", action="store_true", help="导出列头柜页")
    parser.add_argument(
        "--export-power-index",
        action="store_true",
        help="写出 hx-power-page-index.json，并给架构树配电室叶子补 pageId",
    )
    parser.add_argument("--power-index-out", default=DEFAULT_POWER_INDEX, help="配电页索引 JSON")
    parser.add_argument(
        "--patch-arch",
        action="store_true",
        default=True,
        help="--export-power-index 时同步 patch hx-arch-tree.json（默认开）",
    )
    parser.add_argument("--merge", action="store_true", default=True, help="与已有 hx-scenes 按 page id 合并（默认开）")
    parser.add_argument("--replace", action="store_true", help="不合并，仅用本次导出覆盖 pages")
    args = parser.parse_args()

    if not os.path.isfile(args.db):
        print(f"DB not found: {args.db}", file=sys.stderr)
        return 1

    conn = sqlite3.connect(args.db)
    try:
        if args.list:
            return cmd_list(conn, args.model, args.arch)
        if args.manifest:
            return cmd_manifest(conn, args.model, args.arch, args.manifest)
        if args.export_power_index:
            return cmd_power_index(
                conn,
                args.model,
                args.power_index_out,
                args.arch,
                patch_arch=args.patch_arch,
                db_path=args.db,
            )
        if not args.export_power and not args.export_columns:
            parser.error("指定 --list、--manifest、--export-power-index，或 --export-power / --export-columns")
        return cmd_export(
            conn,
            args.model,
            args.output,
            export_power=args.export_power,
            export_columns=args.export_columns,
            merge=not args.replace,
        )
    finally:
        conn.close()


if __name__ == "__main__":
    raise SystemExit(main())
