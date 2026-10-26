#!/usr/bin/env python3
"""扫描 SQL/SQLite 中未删除的 display_model_layer 组态 JSON 完整性。"""
import argparse
import base64
import json
import re
import sqlite3
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
SQL_DEFAULT = ROOT / "Mysql_Backup_2026-07-06_19-58-16.sql"
SQLITE_DEFAULT = ROOT / "ism_server_user" / "data" / "db" / "ism.db"
XUNAN_UUID = "3ec5821f-b512-2adb-3e1c-473720d0a93e"
MODEL_ID = "b8b4c094-faa9-a22a-1d0d-037539b27a6c"


def check_cell(cell, page_name):
    issues = []
    data = cell.get("data", {})
    detail = data.get("detail", {})
    shape = cell.get("shape", data.get("shape", "?"))
    animate = detail.get("animate")
    style = detail.get("style", {})

    if animate is None:
        issues.append(f"{page_name}/{shape}: detail.animate missing")
    elif animate.get("selected") is None:
        issues.append(f"{page_name}/{shape}: animate.selected missing")
    if animate is not None:
        for key in ("animateElement", "animateList"):
            if animate.get(key) is None:
                issues.append(f"{page_name}/{shape}: animate.{key} missing")

    if shape == "view-svg-text":
        if not style.get("text"):
            issues.append(f"{page_name}/{shape}: style.text missing")
        if style.get("visible") is None:
            issues.append(f"{page_name}/{shape}: style.visible missing")
        if style.get("diy") is None:
            issues.append(f"{page_name}/{shape}: style.diy missing")

    return issues


def scan_pages(pages):
    stats = {
        "pages": 0,
        "cells": 0,
        "text_cells": 0,
        "issues": [],
    }
    for page_name, page_id, components_b64 in pages:
        stats["pages"] += 1
        try:
            comp = json.loads(base64.b64decode(components_b64))
        except Exception as e:
            stats["issues"].append(f"{page_name}: decode fail {e}")
            continue
        for cell in comp.get("cells", []):
            stats["cells"] += 1
            shape = cell.get("shape", cell.get("data", {}).get("shape", ""))
            if shape == "view-svg-text":
                stats["text_cells"] += 1
            stats["issues"].extend(check_cell(cell, page_name))
    return stats


def load_from_sqlite(db_path):
    conn = sqlite3.connect(db_path)
    cur = conn.cursor()
    cur.execute(
        """
        SELECT page_name, page_id, components
        FROM display_model_layer
        WHERE model_id=? AND deleted_at IS NULL
        """,
        (MODEL_ID,),
    )
    return cur.fetchall()


def load_from_sql(sql_path):
    """逐行解析 display_model_layer INSERT（每行一条完整记录）。"""
    pages = []
    prefix = "INSERT INTO `display_model_layer`"
    for line in sql_path.read_text(encoding="utf-8", errors="replace").splitlines():
        if not line.startswith(prefix) or MODEL_ID not in line:
            continue
        # VALUES (id,'ts','ts',deleted_at,'model_id','page_name','page_id',is_home,page_type,'layer','components',is_login)
        m = re.search(
            r"VALUES\s*\(\d+,'[^']*','[^']*',(NULL|'[^']*'),'"
            + re.escape(MODEL_ID)
            + r"','([^']+)','([^']+)',[^,]*,[^,]*,'[^']*','([^']*)'",
            line,
        )
        if m and m.group(1) == "NULL":
            pages.append((m.group(2), m.group(3), m.group(4)))
    return pages


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--sql", type=Path, default=SQL_DEFAULT)
    parser.add_argument("--sqlite", type=Path, default=SQLITE_DEFAULT)
    parser.add_argument("--source", choices=["sql", "sqlite", "both"], default="both")
    args = parser.parse_args()

    ok = True
    for source, loader, path in [
        ("SQL", load_from_sql, args.sql),
        ("SQLite", load_from_sqlite, args.sqlite),
    ]:
        if args.source == "sql" and source != "SQL":
            continue
        if args.source == "sqlite" and source != "SQLite":
            continue
        if not path.exists():
            print(f"[{source}] SKIP: {path} not found")
            continue
        pages = loader(path)
        stats = scan_pages(pages)
        issue_count = len(stats["issues"])
        passed = issue_count == 0
        ok = ok and passed
        print(f"\n=== {source} ({path.name}) ===")
        print(f"  pages: {stats['pages']}")
        print(f"  cells: {stats['cells']}")
        print(f"  text cells: {stats['text_cells']}")
        print(f"  issues: {issue_count}")
        print(f"  result: {'PASS' if passed else 'FAIL'}")
        for issue in stats["issues"][:20]:
            print(f"    - {issue}")
        if issue_count > 20:
            print(f"    ... and {issue_count - 20} more")

    return 0 if ok else 1


if __name__ == "__main__":
    sys.exit(main())
