#!/usr/bin/env python3
"""将指定大屏从预生成页面迁移为运行时轻量模板（默认只演练）。"""

from __future__ import annotations

import argparse
import sqlite3
from pathlib import Path


ROOT = Path(__file__).resolve().parent.parent
DEFAULT_DB = ROOT / "ism_server_user" / "data" / "db" / "ism.db"
DEFAULT_MODEL = "b8b4c094-faa9-a22a-1d0d-037539b27a6c"
LEGACY_NAMES = (
    "page_name = 'device-detail'"
    " OR page_name LIKE 'device-%'"
    " OR page_name LIKE 'building-%'"
    " OR page_name LIKE 'floor-%'"
    " OR page_name LIKE 'room-%'"
    " OR page_name LIKE 'oneline%'"
)


def main() -> int:
    parser = argparse.ArgumentParser(description="ISM 大屏轻量模板迁移")
    parser.add_argument("--db", default=str(DEFAULT_DB))
    parser.add_argument("--model", default=DEFAULT_MODEL)
    parser.add_argument("--apply", action="store_true", help="确认后执行软删除")
    args = parser.parse_args()

    db = sqlite3.connect(args.db)
    template_count = db.execute(
        """SELECT COUNT(*) FROM display_model_layer
           WHERE model_id=? AND deleted_at IS NULL AND page_name LIKE '模板-%'""",
        (args.model,),
    ).fetchone()[0]
    legacy_count = db.execute(
        f"""SELECT COUNT(*) FROM display_model_layer
            WHERE model_id=? AND deleted_at IS NULL AND ({LEGACY_NAMES})""",
        (args.model,),
    ).fetchone()[0]
    if template_count < 2:
        raise SystemExit("未发现完整运行时模板页，已拒绝迁移")

    action = "将软删除" if args.apply else "将演练软删除"
    print(f"model={args.model}，模板页={template_count}，{action}旧页面={legacy_count}")
    if not args.apply:
        print("演练完成；确认备份后使用 --apply 执行。")
        return 0

    db.execute(
        f"""UPDATE display_model_layer
            SET deleted_at=datetime('now'), updated_at=datetime('now')
            WHERE model_id=? AND deleted_at IS NULL AND ({LEGACY_NAMES})""",
        (args.model,),
    )
    db.commit()
    print("迁移完成：旧页面已软删除，运行时导航将仅使用轻量模板。")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
