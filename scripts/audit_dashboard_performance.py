#!/usr/bin/env python3
"""轻量大屏迁移后的只读性能结构审计。"""

from __future__ import annotations

import argparse
import sqlite3
from pathlib import Path


ROOT = Path(__file__).resolve().parent.parent
DEFAULT_DB = ROOT / "ism_server_user" / "data" / "db" / "ism.db"
DEFAULT_MODEL = "b8b4c094-faa9-a22a-1d0d-037539b27a6c"


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--db", default=str(DEFAULT_DB))
    parser.add_argument("--model", default=DEFAULT_MODEL)
    args = parser.parse_args()
    db = sqlite3.connect(f"file:{Path(args.db).resolve()}?mode=ro", uri=True)
    active = db.execute(
        "SELECT COUNT(*) FROM display_model_layer WHERE model_id=? AND deleted_at IS NULL",
        (args.model,),
    ).fetchone()[0]
    templates = db.execute(
        """SELECT COUNT(*) FROM display_model_layer
           WHERE model_id=? AND deleted_at IS NULL AND page_name LIKE '模板-%'""",
        (args.model,),
    ).fetchone()[0]
    legacy = db.execute(
        """SELECT COUNT(*) FROM display_model_layer
           WHERE model_id=? AND deleted_at IS NULL
             AND (page_name='device-detail' OR page_name LIKE 'device-%'
                  OR page_name LIKE 'building-%' OR page_name LIKE 'floor-%'
                  OR page_name LIKE 'room-%' OR page_name LIKE 'oneline%')""",
        (args.model,),
    ).fetchone()[0]
    print(f"active_pages={active}")
    print(f"runtime_templates={templates}")
    print(f"legacy_prebuilt_pages={legacy}")
    if templates < 2 or legacy:
        raise SystemExit(1)
    print("PASS: 运行时轻模板已取代预生成组织/设备页面。")


if __name__ == "__main__":
    main()
