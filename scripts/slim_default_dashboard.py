#!/usr/bin/env python3
"""
默认大屏（组态界面）瘦身：17 页 / 42 组件 → 4 页 / ~25 组件。

目标结构
  Level0  主界面          保留标题/时间/用户/告警/导航（菜单精简）
  Level1  配电室总览      卡片入口（替代 13 个重复实时表页）
  Level2  统一实时数据    复用一份 ism-view-real-table
  Level3  历史报警        保留原页

用法
  # 仅从 SQL 备份生成补丁 SQL（不连库，适合离线交付）
  python3 scripts/slim_default_dashboard.py --from-sql Mysql_Backup_2026-07-06_19-58-16.sql

  # 直接写入本地/现场库（自动判 SQLite / OceanBase:2881）
  python3 scripts/slim_default_dashboard.py --apply

  # 指定 model
  NCC_MODEL_ID=b8b4c094-... python3 scripts/slim_default_dashboard.py --apply

安全
  - 默认 dry-run；--apply 才写库
  - 多余页软删除（deleted_at），不物理删
  - 只动指定 model_id
"""
from __future__ import annotations

import argparse
import base64
import json
import os
import re
import socket
import sqlite3
import sys
import time
import uuid
from copy import deepcopy
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
DEFAULT_SQL = ROOT / "Mysql_Backup_2026-07-06_19-58-16.sql"
DEFAULT_MODEL = "b8b4c094-faa9-a22a-1d0d-037539b27a6c"
DEFAULT_SQLITE = ROOT / "ism_server_user" / "data" / "db" / "ism.db"

# 稳定 page_id（可反复跑）
PAGE_ID_OVERVIEW = str(uuid.uuid5(uuid.NAMESPACE_DNS, "ism-slim-room-overview"))
PAGE_ID_UNIFIED = str(uuid.uuid5(uuid.NAMESPACE_DNS, "ism-slim-unified-realtime"))

KEEP_HOME = "主界面"
KEEP_ALARM = "历史报警"
NAME_OVERVIEW = "配电室总览"
NAME_UNIFIED = "统一实时数据"

# 原 13 个「单表配电室」页名（将被软删，入口合并到总览卡片）
ROOM_CARD_NAMES = [
    "1A配电室", "1B配电室",
    "2A1配电室", "2A2配电室", "2A3配电室", "2A4配电室",
    "2B1配电室", "2B2配电室", "2B3配电室", "2B4配电室",
    "3A1配电室", "3A2配电室", "3A3配电室", "3A4配电室",
    "4A1配电室",
]


def _ob_up(port=2881):
    s = socket.socket()
    try:
        s.settimeout(0.5)
        s.connect(("127.0.0.1", port))
        return True
    except OSError:
        return False
    finally:
        s.close()


def gen_uid(seed: str) -> str:
    return uuid.uuid5(uuid.NAMESPACE_DNS, f"ism-slim-{seed}").hex


def b64_components(cells):
    raw = json.dumps({"cells": cells}, ensure_ascii=False).encode("utf-8")
    return base64.b64encode(raw).decode("ascii")


def decode_components(comps_b64: str):
    if not comps_b64:
        return {"cells": []}
    return json.loads(base64.b64decode(comps_b64))


def _base_animate():
    return {
        "selected": [],
        "animateElement": [],
        "animateList": [],
        "condition": {
            "deviceSN": "", "selectVideoType": 0, "isBandDevice": False,
            "bandType": 1, "dataID": "", "dataName": "",
            "operator": "", "OperatorValue": "", "OperatorMaxValue": "",
        },
        "isExpression": False,
        "move": {
            "x": {"deviceSN": "", "selectVideoType": 0, "isBandDevice": False,
                  "bandType": 1, "dataID": "", "dataName": ""},
            "y": {"deviceSN": "", "selectVideoType": 0, "isBandDevice": False,
                  "bandType": 1, "dataID": "", "dataName": ""},
        },
    }


def make_text(seed, x, y, w, h, text, color="#e0faff", font_size=16, z=10, action=None):
    cell_id = gen_uid(seed)
    return {
        "shape": "view-svg-text",
        "id": cell_id, "x": x, "y": y, "width": w, "height": h,
        "zIndex": z, "visible": True,
        "position": {"x": x, "y": y},
        "size": {"width": w, "height": h},
        "data": {
            "detail": {
                "type": "view-svg-text",
                "identifier": cell_id,
                "name": text,
                "style": {
                    "position": {"x": x, "y": y, "w": w, "h": h},
                    "visible": 1, "opacity": 1, "diy": [],
                    "text": text, "fontSize": font_size, "foreColor": color,
                    "borderWidth": 0, "BorderEdges": 0,
                },
                "animate": _base_animate(),
                "action": action or [],
                "active": [],
                "dataBind": [],
            }
        },
    }


def make_box(seed, x, y, w, h, z=1):
    cell_id = gen_uid(seed)
    return {
        "shape": "dv-border-box13",
        "id": cell_id, "x": x, "y": y, "width": w, "height": h,
        "zIndex": z, "visible": True,
        "position": {"x": x, "y": y},
        "size": {"width": w, "height": h},
        "data": {
            "detail": {
                "type": "dv-border-box13",
                "identifier": cell_id,
                "name": seed,
                "style": {"position": {"x": x, "y": y, "w": w, "h": h}, "visible": 1},
                "animate": _base_animate(),
                "action": [], "active": [], "dataBind": [],
            }
        },
    }


def nav_action(model_id, page_uuid):
    return [{
        "type": "click",
        "action": "link",
        "isPopUp": False,
        "link": {
            "linkType": "Inside",
            "isPopUp": False,
            "Inside": {
                "displayUUID": model_id,
                "pageUUID": page_uuid,
                "displayType": 1,
            },
        },
    }]


def default_layer():
    return json.dumps({
        "backColor": "#0a0e17",
        "backgroundImage": "",
        "height": 1080,
        "width": 1920,
        "widthHeightRatio": "",
        "autoSize": 1,
        "Padding": 0,
    }, ensure_ascii=False)


# ── load from SQL dump ──────────────────────────────────────

def load_pages_from_sql(sql_path: Path, model_id: str):
    prefix = "INSERT INTO `display_model_layer`"
    pages = {}
    for line in sql_path.read_text(encoding="utf-8", errors="replace").splitlines():
        if not line.startswith(prefix) or model_id not in line:
            continue
        m = re.search(
            r"VALUES\s*\((\d+),'([^']*)','([^']*)',([^,]*),'"
            + re.escape(model_id)
            + r"','([^']+)','([^']+)',(\d+),(\d+),'([^']*)','([^']*)',(\d+)\)",
            line,
        )
        if not m:
            continue
        rid, _c, _u, _d, page_name, page_id, is_home, page_type, layer, comps, is_login = m.groups()
        pages[page_name] = {
            "id": int(rid),
            "page_name": page_name,
            "page_id": page_id,
            "is_home": int(is_home),
            "page_type": int(page_type),
            "is_login": int(is_login),
            "layer": layer,
            "components": comps,
            "cells": decode_components(comps).get("cells", []),
        }
    return pages


def slim_home_menu(home_cells, model_id, home_page_id, overview_id, unified_id, alarm_id):
    cells = deepcopy(home_cells)
    menu_cfg = [
        {
            "DisPlayID": model_id,
            "key": "2",
            "path": home_page_id,
            "title": "配电系统总览",
            "children": [
                {
                    "DisPlayID": model_id,
                    "IsPopUp": False,
                    "key": 1,
                    "path": overview_id,
                    "title": NAME_OVERVIEW,
                },
                {
                    "DisPlayID": model_id,
                    "IsPopUp": False,
                    "key": 2,
                    "path": unified_id,
                    "title": NAME_UNIFIED,
                },
            ],
        },
        {
            "DisPlayID": model_id,
            "IsPopUp": False,
            "key": 3,
            "path": alarm_id,
            "title": "历史报警查询",
        },
    ]
    for c in cells:
        if c.get("shape") == "view-menu-nav":
            c["data"]["detail"]["style"]["MenuConfig"] = menu_cfg
    return cells


def build_overview_cells(model_id, unified_page_id, room_names):
    """轻量总览：1 外框 + 标题/提示 + N 个可点击文字入口（不再每卡套一层边框）。"""
    cells = []
    cells.append(make_box("ov-frame", 40, 40, 1840, 1000, z=1))
    cells.append(make_text("ov-title", 80, 60, 700, 40, "配电室总览（轻量入口）",
                           color="#22d3ee", font_size=28, z=10))
    cells.append(make_text(
        "ov-hint", 80, 110, 1400, 28,
        "点击名称进入统一实时数据页（原 13 个重复表页已合并）",
        color="#94a3b8", font_size=14, z=10,
    ))
    cols, card_w, card_h, gap_x, gap_y = 5, 320, 56, 24, 18
    start_x, start_y = 80, 180
    for i, name in enumerate(room_names):
        r, c = divmod(i, cols)
        x = start_x + c * (card_w + gap_x)
        y = start_y + r * (card_h + gap_y)
        cells.append(make_text(
            f"card-txt-{i}", x, y, card_w, card_h, f"▸ {name}",
            color="#e0faff", font_size=18, z=12,
            action=nav_action(model_id, unified_page_id),
        ))
    return cells


def build_unified_cells(template_cells):
    """复用一份实时表模板；若无模板则放提示文字。"""
    if template_cells:
        cells = deepcopy(template_cells)
        # 补齐 animate.selected，避免渲染崩溃
        for c in cells:
            detail = ((c.get("data") or {}).get("detail") or {})
            animate = detail.get("animate")
            if not isinstance(animate, dict):
                detail["animate"] = _base_animate()
            else:
                animate.setdefault("selected", [])
                animate.setdefault("animateElement", [])
                animate.setdefault("animateList", [])
        return cells
    return [
        make_text("uni-empty", 80, 80, 800, 40,
                  "暂无实时表模板，请在编辑器中绑定测点", color="#f87171", font_size=20),
    ]


def plan_slim(pages: dict, model_id: str):
    if KEEP_HOME not in pages:
        raise SystemExit(f"找不到首页页名 {KEEP_HOME!r}，当前页: {list(pages)}")
    home = pages[KEEP_HOME]
    alarm = pages.get(KEEP_ALARM)
    if not alarm:
        raise SystemExit(f"找不到 {KEEP_ALARM!r}")

    # 选一个最轻的实时表页作模板（优先 2A1，否则任意 ism-view-real-table）
    template_cells = None
    for prefer in ("2A1配电室", "2B1配电室", "1B配电室"):
        if prefer in pages and pages[prefer]["cells"]:
            template_cells = pages[prefer]["cells"]
            break
    if template_cells is None:
        for p in pages.values():
            for c in p["cells"]:
                if c.get("shape") == "ism-view-real-table":
                    template_cells = p["cells"]
                    break
            if template_cells:
                break

    room_names = [n for n in ROOM_CARD_NAMES if n in pages] or ROOM_CARD_NAMES

    home_cells = slim_home_menu(
        home["cells"], model_id, home["page_id"],
        PAGE_ID_OVERVIEW, PAGE_ID_UNIFIED, alarm["page_id"],
    )
    overview_cells = build_overview_cells(model_id, PAGE_ID_UNIFIED, room_names)
    unified_cells = build_unified_cells(template_cells)

    keep_names = {KEEP_HOME, KEEP_ALARM, NAME_OVERVIEW, NAME_UNIFIED}
    soft_delete = [n for n in pages if n not in keep_names]

    upserts = [
        {
            "page_name": KEEP_HOME,
            "page_id": home["page_id"],
            "is_home": 1,
            "layer": home.get("layer") or default_layer(),
            "components": b64_components(home_cells),
            "cell_count": len(home_cells),
        },
        {
            "page_name": NAME_OVERVIEW,
            "page_id": PAGE_ID_OVERVIEW,
            "is_home": 0,
            "layer": default_layer(),
            "components": b64_components(overview_cells),
            "cell_count": len(overview_cells),
        },
        {
            "page_name": NAME_UNIFIED,
            "page_id": PAGE_ID_UNIFIED,
            "is_home": 0,
            "layer": default_layer(),
            "components": b64_components(unified_cells),
            "cell_count": len(unified_cells),
        },
        {
            "page_name": KEEP_ALARM,
            "page_id": alarm["page_id"],
            "is_home": 0,
            "layer": alarm.get("layer") or default_layer(),
            "components": alarm["components"] if isinstance(alarm["components"], str)
            else b64_components(alarm["cells"]),
            "cell_count": len(alarm["cells"]),
        },
    ]
    return upserts, soft_delete, {
        "before_pages": len(pages),
        "before_cells": sum(len(p["cells"]) for p in pages.values()),
        "after_pages": 4,
        "after_cells": sum(u["cell_count"] for u in upserts),
        "soft_delete": soft_delete,
    }


def emit_sql(model_id: str, upserts, soft_delete_names, out_path: Path):
    lines = [
        "-- ISM 默认大屏瘦身补丁（软删多余页 + upsert 4 页）",
        f"-- model_id={model_id}",
        "SET NAMES utf8mb4;",
        "",
    ]
    for name in soft_delete_names:
        safe = name.replace("'", "''")
        lines.append(
            f"UPDATE display_model_layer SET deleted_at=NOW() "
            f"WHERE model_id='{model_id}' AND page_name='{safe}' AND deleted_at IS NULL;"
        )
    lines.append("")
    for u in upserts:
        pn = u["page_name"].replace("'", "''")
        pid = u["page_id"]
        layer = u["layer"].replace("\\", "\\\\").replace("'", "\\'")
        comps = u["components"].replace("\\", "\\\\").replace("'", "\\'")
        lines.append(f"-- upsert {pn} cells={u['cell_count']}")
        lines.append(
            f"UPDATE display_model_layer SET deleted_at=NULL, page_name='{pn}', "
            f"is_home={u['is_home']}, layer='{layer}', components='{comps}', updated_at=NOW() "
            f"WHERE model_id='{model_id}' AND page_id='{pid}';"
        )
        lines.append(
            f"INSERT INTO display_model_layer "
            f"(model_id, page_name, page_id, is_home, is_login, page_type, layer, components, created_at, updated_at) "
            f"SELECT '{model_id}', '{pn}', '{pid}', {u['is_home']}, 0, 1, '{layer}', '{comps}', NOW(), NOW() "
            f"FROM DUAL WHERE NOT EXISTS ("
            f"SELECT 1 FROM display_model_layer WHERE model_id='{model_id}' AND page_id='{pid}' AND deleted_at IS NULL"
            f");"
        )
        lines.append("")
    out_path.parent.mkdir(parents=True, exist_ok=True)
    out_path.write_text("\n".join(lines) + "\n", encoding="utf-8")
    return out_path


def apply_db(model_id: str, upserts, soft_delete_names):
    use_sqlite = os.environ.get("ISM_FORCE_OB") != "1" and not _ob_up()
    if use_sqlite:
        if not DEFAULT_SQLITE.exists():
            raise SystemExit(f"SQLite 不存在: {DEFAULT_SQLITE}")
        conn = sqlite3.connect(str(DEFAULT_SQLITE))
        ph = "?"
        now_fn = "datetime('now')"
    else:
        import pymysql
        conn = pymysql.connect(
            host="127.0.0.1", port=2881,
            user="root@ism_tenant", password="ism2024!",
            database="ism", charset="utf8mb4",
        )
        ph = "%s"
        now_fn = "NOW()"

    cur = conn.cursor()
    # 按 page_name 软删（来自 SQL 审计名单）
    for name in soft_delete_names:
        cur.execute(
            f"UPDATE display_model_layer SET deleted_at={now_fn} "
            f"WHERE model_id={ph} AND page_name={ph} AND deleted_at IS NULL",
            (model_id, name),
        )
    # 本地若已被 NCC 重建成数百页：软删本 model 下所有「非保留 page_id」
    keep_ids = tuple(u["page_id"] for u in upserts)
    if keep_ids:
        placeholders = ",".join([ph] * len(keep_ids))
        cur.execute(
            f"UPDATE display_model_layer SET deleted_at={now_fn} "
            f"WHERE model_id={ph} AND deleted_at IS NULL AND page_id NOT IN ({placeholders})",
            (model_id, *keep_ids),
        )
    for u in upserts:
        cur.execute(
            f"SELECT id FROM display_model_layer WHERE model_id={ph} AND page_id={ph} AND deleted_at IS NULL",
            (model_id, u["page_id"]),
        )
        row = cur.fetchone()
        if row:
            cur.execute(
                f"UPDATE display_model_layer SET page_name={ph}, is_home={ph}, layer={ph}, "
                f"components={ph}, deleted_at=NULL, updated_at={now_fn} WHERE id={ph}",
                (u["page_name"], u["is_home"], u["layer"], u["components"], row[0]),
            )
        else:
            cur.execute(
                f"INSERT INTO display_model_layer "
                f"(model_id, page_name, page_id, is_home, is_login, page_type, layer, components, created_at, updated_at) "
                f"VALUES ({ph},{ph},{ph},{ph},0,1,{ph},{ph},{now_fn},{now_fn})",
                (model_id, u["page_name"], u["page_id"], u["is_home"], u["layer"], u["components"]),
            )
    conn.commit()
    cur.execute(
        f"SELECT page_name, page_id, is_home, LENGTH(components) FROM display_model_layer "
        f"WHERE model_id={ph} AND deleted_at IS NULL ORDER BY is_home DESC, id",
        (model_id,),
    )
    rows = cur.fetchall()
    conn.close()
    return rows


def main():
    ap = argparse.ArgumentParser(description="默认大屏瘦身 17→4 页")
    ap.add_argument("--from-sql", type=Path, default=DEFAULT_SQL, help="MySQL dump 路径")
    ap.add_argument("--model-id", default=os.environ.get("NCC_MODEL_ID", DEFAULT_MODEL))
    ap.add_argument("--apply", action="store_true", help="写入数据库（默认只生成 SQL）")
    ap.add_argument("--out-sql", type=Path,
                    default=ROOT / "releases" / "sql" / "slim_default_dashboard.sql")
    args = ap.parse_args()

    if not args.from_sql.exists():
        raise SystemExit(f"SQL 不存在: {args.from_sql}")

    pages = load_pages_from_sql(args.from_sql, args.model_id)
    if not pages:
        raise SystemExit(f"SQL 中未找到 model_id={args.model_id} 的图层")

    upserts, soft_delete, stats = plan_slim(pages, args.model_id)
    print("=== 瘦身计划 ===")
    print(f"  model_id     = {args.model_id}")
    print(f"  before       = {stats['before_pages']} 页 / {stats['before_cells']} 组件")
    print(f"  after        = {stats['after_pages']} 页 / {stats['after_cells']} 组件")
    print(f"  soft_delete  = {len(soft_delete)} 页: {soft_delete}")
    for u in upserts:
        print(f"  upsert {u['page_name']}: page_id={u['page_id'][:12]}... cells={u['cell_count']}")

    out = emit_sql(args.model_id, upserts, soft_delete, args.out_sql)
    print(f"\n已生成 SQL 补丁: {out}")

    if args.apply:
        rows = apply_db(args.model_id, upserts, soft_delete)
        print("\n=== 写库后页面 ===")
        for r in rows:
            print(f"  home={r[2]} name={r[0]!r} page_id={r[1]} comps_len={r[3]}")
        print("APPLY OK")
    else:
        print("\n(dry-run) 未写库。确认后加 --apply，或把 SQL 拿到现场执行。")


if __name__ == "__main__":
    main()
