#!/usr/bin/env python3
"""
层级模板页补丁（配合前端运行时槽位重映射）:
  1. 备份 SQLite
  2. 从软删除样本 floor-1275279789-default 新建「模板-设备组」(template_kind='floor')
  3. 四个容器模板页的页级静态文本改为 {{nav.*}} 占位符（运行时按 navContext 替换）

幂等：可反复执行；已含 {{nav. 的文本跳过。

用法:
  python3 scripts/patch_dashboard_floor_template.py            # 预览
  python3 scripts/patch_dashboard_floor_template.py --apply    # 写库
"""
from __future__ import annotations

import argparse
import base64
import json
import re
import shutil
import sqlite3
import sys
import time
import uuid
from datetime import datetime, timezone
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
DB = ROOT / "ism_server_user" / "data" / "db" / "ism.db"
MODEL = "b8b4c094-faa9-a22a-1d0d-037539b27a6c"
FLOOR_TPL_PAGE_ID = uuid.uuid5(uuid.NAMESPACE_DNS, "ism-tpl-floor").hex
FLOOR_SAMPLE_NAME = "floor-1275279789-default"


def now_sql():
    return datetime.now(timezone.utc).strftime("%Y-%m-%d %H:%M:%S")


def decode_components(raw: str):
    s = raw or ""
    if not s:
        return {"cells": []}
    if not s.lstrip().startswith("{"):
        s = base64.b64decode(s).decode("utf-8", errors="replace")
    return json.loads(s)


def encode_components(obj) -> str:
    raw = json.dumps(obj, ensure_ascii=False, separators=(",", ":"))
    return base64.b64encode(raw.encode("utf-8")).decode("ascii")


# 每个模板页的「精确文本 → 占位符」规则（正则全匹配）
PAGE_TEXT_RULES = {
    "模板-机柜": [
        (r"^配电室_机房模块3A1$", "{{nav.name}}"),
        (r"^🏢 配电室_机房模块3A1$", "🏢 {{nav.name}}"),
        (r"^\d+台设备 · \d+条异常$",
         "{{nav.deviceCount}}台设备 · {{nav.abnormalCount}}条异常"),
    ],
    "模板-区域": [
        (r"^配电室$", "{{nav.name}}"),
        (r"^\d+子区域 · \d+台设备 · 在线\d* · 异常\d*$",
         "{{nav.childCount}}子区域 · {{nav.deviceCount}}台设备 · 在线{{nav.onlineCount}} · 异常{{nav.abnormalCount}}"),
    ],
    "模板-机房": [
        (r"^配电室$", "{{nav.name}}"),
        (r"^🏛 配电室$", "🏛 {{nav.name}}"),
        (r"^\d+个机柜 · \d+台设备 · 在线\d* · 异常\d*( · 表格展示 \d+ 台)?$",
         "{{nav.childCount}}个机柜 · {{nav.deviceCount}}台设备 · 在线{{nav.onlineCount}} · 异常{{nav.abnormalCount}}"),
    ],
    "模板-设备组": [
        (r"^default设备组$", "{{nav.name}}"),
        (r"^📋 default设备组$", "📋 {{nav.name}}"),
        (r"^\d+台设备$", "{{nav.deviceCount}}台设备"),
        (r"^共 \d+ 台设备 \| 运行: \d+台 \| 离线: \d+台$",
         "共 {{nav.deviceCount}} 台设备 | 运行: {{nav.onlineCount}}台 | 离线: {{nav.offlineCount}}台"),
    ],
}


def apply_text_rules(comps: dict, rules) -> int:
    n = 0
    for cell in comps.get("cells") or []:
        detail = (((cell or {}).get("data") or {}).get("detail")) or {}
        style = detail.get("style") or {}
        text = style.get("text")
        if not isinstance(text, str) or not text or "{{nav." in text:
            continue
        for pat, repl in rules:
            if re.fullmatch(pat, text):
                style["text"] = repl
                # name 字段仅编辑器展示用，保持同步
                if detail.get("name") == text:
                    detail["name"] = repl
                n += 1
                break
    return n


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--apply", action="store_true")
    args = ap.parse_args()

    if not DB.exists():
        print(f"DB not found: {DB}", file=sys.stderr)
        sys.exit(1)

    if args.apply:
        bak = DB.with_name(f"ism.db.bak-{int(time.time())}")
        shutil.copy2(DB, bak)
        print(f"+ backup {bak.name}")

    conn = sqlite3.connect(str(DB))

    # ---- 1. floor 模板 ----
    src = conn.execute(
        "SELECT layer, components, page_type FROM display_model_layer "
        "WHERE model_id=? AND page_name=? ORDER BY id DESC LIMIT 1",
        (MODEL, FLOOR_SAMPLE_NAME),
    ).fetchone()
    if not src:
        print(f"! floor sample {FLOOR_SAMPLE_NAME} not found", file=sys.stderr)
        sys.exit(1)
    layer, comp_raw, page_type = src
    comps = decode_components(comp_raw)
    n_floor_rules = apply_text_rules(comps, PAGE_TEXT_RULES["模板-设备组"])
    print(f"floor tpl: cells={len(comps.get('cells') or [])} placeholder-edits={n_floor_rules} "
          f"page_id={FLOOR_TPL_PAGE_ID}")

    if args.apply:
        ts = now_sql()
        exist = conn.execute(
            "SELECT id FROM display_model_layer WHERE model_id=? AND page_id=?",
            (MODEL, FLOOR_TPL_PAGE_ID),
        ).fetchone()
        b64 = encode_components(comps)
        if exist:
            conn.execute(
                """UPDATE display_model_layer
                   SET page_name=?, layer=?, components=?, is_home=0, page_type=?,
                       template_kind='floor', template_model_uuid='',
                       updated_at=?, deleted_at=NULL
                   WHERE id=?""",
                ("模板-设备组", layer or "", b64, page_type, ts, exist[0]),
            )
            print("  [update] 模板-设备组")
        else:
            conn.execute(
                """INSERT INTO display_model_layer
                   (created_at, updated_at, deleted_at, model_id, page_name, page_id,
                    is_home, is_login, page_type, layer, components,
                    template_kind, template_model_uuid)
                   VALUES (?,?,NULL,?,?,?,0,0,?,?,?,'floor','')""",
                (ts, ts, MODEL, "模板-设备组", FLOOR_TPL_PAGE_ID,
                 page_type, layer or "", b64),
            )
            print("  [insert] 模板-设备组")

    # ---- 2. 现有容器模板占位符 ----
    for page_name in ("模板-机柜", "模板-区域", "模板-机房"):
        row = conn.execute(
            "SELECT id, components FROM display_model_layer "
            "WHERE model_id=? AND page_name=? AND deleted_at IS NULL",
            (MODEL, page_name),
        ).fetchone()
        if not row:
            print(f"! {page_name} not found, skip")
            continue
        comps = decode_components(row[1])
        n = apply_text_rules(comps, PAGE_TEXT_RULES[page_name])
        print(f"{page_name}: placeholder-edits={n}")
        if args.apply and n:
            conn.execute(
                "UPDATE display_model_layer SET components=?, updated_at=? WHERE id=?",
                (encode_components(comps), now_sql(), row[0]),
            )

    if args.apply:
        conn.commit()
        print("committed.")
    else:
        print("(dry-run) 加 --apply 写库")
    conn.close()


if __name__ == "__main__":
    main()
