#!/usr/bin/env python3
"""硬删除大屏旧预生成页，只保留三模板运行链路。

保留：
  - 首页模板     template_kind=home
  - 设备列表模板 template_kind=deviceList
  - 点位列表模板 template_kind=datapointList

用法：
  python3 scripts/prune_legacy_dashboard_pages.py                      # dry-run 默认 model
  python3 scripts/prune_legacy_dashboard_pages.py --apply              # 硬删除 SQLite
  python3 scripts/prune_legacy_dashboard_pages.py --apply --oceanbase
  python3 scripts/prune_legacy_dashboard_pages.py --apply --all-template-models
"""

from __future__ import annotations

import argparse
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
DEFAULT_MODEL = "b8b4c094-faa9-a22a-1d0d-037539b27a6c"
KEEP_NAMES = ("首页模板", "设备列表模板", "点位列表模板")
KEEP_KINDS = ("home", "deviceList", "datapointList")
LEGACY_EXACT = ("device-detail", "oneline", "main", "building-detail", "floor-detail")


def load_app_conf():
    conf = ROOT / "ism_server_user" / "conf" / "app.conf"
    if not conf.exists():
        # 部署包内：脚本在 <pkg>/scripts/，conf 在 <pkg>/ism_server_user/conf/
        conf = Path(__file__).resolve().parents[1] / "ism_server_user" / "conf" / "app.conf"
    raw = conf.read_text(encoding="utf-8", errors="ignore")
    kv = {}
    for line in raw.splitlines():
        line = line.strip()
        if not line or line.startswith("#") or line.startswith(";") or "=" not in line or line.startswith("["):
            continue
        k, v = line.split("=", 1)
        kv[k.strip()] = v.strip()
    return kv


def connect(use_ob: bool):
    kv = load_app_conf()
    dbtype = int(kv.get("dbtype", "1"))
    if use_ob or dbtype == 4:
        import pymysql

        return (
            pymysql.connect(
                host=kv.get("oceanbasehost", "127.0.0.1"),
                port=int(kv.get("oceanbaseport", "2881")),
                user=kv.get("oceanbaseuser", "root@ism_tenant"),
                password=kv.get("oceanbasepwd", ""),
                database=kv.get("oceanbasedbname", "ism"),
                charset="utf8mb4",
                autocommit=False,
            ),
            "mysql",
        )
    import sqlite3

    db = ROOT / "ism_server_user" / "data" / "db" / "ism.db"
    if not db.exists():
        raise SystemExit(f"SQLite not found: {db}")
    return sqlite3.connect(str(db)), "sqlite"


def ph(dialect: str, n: int) -> str:
    mark = "%s" if dialect == "mysql" else "?"
    return ",".join([mark] * n)


def list_template_models(cur, dialect: str):
    qmark = "%s" if dialect == "mysql" else "?"
    cur.execute(
        f"""
        SELECT DISTINCT model_id FROM display_model_layer
        WHERE COALESCE(template_kind,'') IN ({ph(dialect, 3)})
           OR page_name IN ({ph(dialect, 3)})
        """,
        (*KEEP_KINDS, *KEEP_NAMES),
    )
    return [r[0] for r in cur.fetchall() if r and r[0]]


def prune_one(cur, dialect: str, mid: str, apply: bool) -> int:
    qmark = "%s" if dialect == "mysql" else "?"
    cur.execute(
        f"""
        SELECT page_name, page_id, is_home, COALESCE(template_kind,'')
        FROM display_model_layer
        WHERE model_id={qmark}
          AND (
            page_name IN ({ph(dialect, 3)})
            OR COALESCE(template_kind,'') IN ({ph(dialect, 3)})
            OR COALESCE(is_home,0)=1
          )
        ORDER BY is_home DESC, page_name
        """,
        (mid, *KEEP_NAMES, *KEEP_KINDS),
    )
    keep_rows = cur.fetchall()
    print(f"\nmodel={mid}")
    print("=== KEEP ===")
    for name, pid, home, kind in keep_rows:
        print(f"  home={home} kind={kind or '-':14s} {name}")
    if not keep_rows:
        print("  skip: 无保留页")
        return 0

    cur.execute(
        f"""
        SELECT page_name, page_id, LENGTH(COALESCE(components,''))
        FROM display_model_layer
        WHERE model_id={qmark}
          AND COALESCE(is_home,0)<>1
          AND COALESCE(page_name,'') NOT IN ({ph(dialect, 3)})
          AND COALESCE(template_kind,'') NOT IN ({ph(dialect, 3)})
          AND (
            page_name LIKE 'building-%'
            OR page_name LIKE 'floor-%'
            OR page_name LIKE 'zone-%'
            OR page_name LIKE 'room-%'
            OR page_name LIKE 'oneline-%'
            OR page_name LIKE 'device-%'
            OR page_name IN ({ph(dialect, len(LEGACY_EXACT))})
          )
        ORDER BY page_name
        """,
        (mid, *KEEP_NAMES, *KEEP_KINDS, *LEGACY_EXACT),
    )
    legacy = cur.fetchall()
    print(f"=== HARD-DELETE candidates: {len(legacy)} ===")
    for name, pid, sz in legacy[:30]:
        print(f"  size={sz:7d}  {name}")
    if len(legacy) > 30:
        print(f"  ... +{len(legacy) - 30} more")
    if not apply or not legacy:
        return 0
    ids = [pid for _, pid, _ in legacy]
    cur.execute(
        f"""
        DELETE FROM display_model_layer
        WHERE model_id={qmark}
          AND page_id IN ({ph(dialect, len(ids))})
        """,
        (mid, *ids),
    )
    print(f"[apply] HARD-deleted {cur.rowcount} pages")
    return cur.rowcount or 0


def main() -> int:
    ap = argparse.ArgumentParser(description="硬删除大屏旧预生成页，只留 3 模板")
    ap.add_argument("--model-id", default=DEFAULT_MODEL)
    ap.add_argument("--all-template-models", action="store_true",
                    help="处理所有已挂三模板的大屏（部署启动推荐）")
    ap.add_argument("--apply", action="store_true", help="真正硬删除；默认 dry-run")
    ap.add_argument("--oceanbase", action="store_true")
    args = ap.parse_args()

    conn, dialect = connect(args.oceanbase)
    cur = conn.cursor()
    print(f"dialect={dialect} mode=HARD-DELETE apply={args.apply}")

    if args.all_template_models:
        mids = list_template_models(cur, dialect)
        print(f"template models: {len(mids)}")
    else:
        mids = [args.model_id]

    total = 0
    for mid in mids:
        total += prune_one(cur, dialect, mid, args.apply)
    if args.apply:
        conn.commit()
        print(f"\n[apply] total HARD-deleted {total}")
    else:
        print("\n[dry-run] 未写库。加 --apply 执行硬删除。")
    conn.close()
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except Exception as e:
        print("ERROR:", e, file=sys.stderr)
        raise
