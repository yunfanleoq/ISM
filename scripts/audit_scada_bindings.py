#!/usr/bin/env python3
"""审计指定 ISM 大屏页面的实时绑定及其 device_real_data 可用性。"""

from __future__ import annotations

import argparse
import base64
import json
import sqlite3
from collections import Counter, defaultdict
from pathlib import Path


REPO_ROOT = Path(__file__).resolve().parent.parent
DEFAULT_DB = REPO_ROOT / "ism_server_user" / "data" / "db" / "ism.db"
DEFAULT_PROJECT = "3ec5821f-b512-2adb-3e1c-473720d0a93e"
DEFAULT_MODEL = "b8b4c094-faa9-a22a-1d0d-037539b27a6c"


def decode_components(raw: str) -> list[dict]:
    try:
        return json.loads(base64.b64decode(raw).decode("utf-8")).get("cells", [])
    except (ValueError, UnicodeDecodeError) as exc:
        raise ValueError(f"components 解码失败: {exc}") from exc


def binding_conditions(detail: dict) -> list[dict]:
    result = []
    for bindings in (detail.get("active") or [], detail.get("dataBind") or []):
        for binding in bindings if isinstance(bindings, list) else []:
            condition = binding.get("condition") if isinstance(binding, dict) else None
            if condition and condition.get("deviceSN") and condition.get("dataID"):
                result.append(condition)
    return result


def real_data_exists(conn: sqlite3.Connection, device_uuid: str, data_uuid: str) -> bool:
    row = conn.execute(
        """SELECT 1 FROM device_real_data
           WHERE device_uuid=? AND model_data_uuid=? LIMIT 1""",
        (device_uuid, data_uuid),
    ).fetchone()
    return row is not None


def main() -> int:
    parser = argparse.ArgumentParser(description="ISM SCADA 绑定审计（只读）")
    parser.add_argument("--db", default=str(DEFAULT_DB), help="SQLite 数据库路径")
    parser.add_argument("--project", default=DEFAULT_PROJECT, help="项目 UUID")
    parser.add_argument("--model", default=DEFAULT_MODEL, help="大屏模型 UUID")
    parser.add_argument("--json-out", help="将完整报告写入 JSON 文件")
    args = parser.parse_args()

    conn = sqlite3.connect(f"file:{Path(args.db).resolve()}?mode=ro", uri=True)
    conn.row_factory = sqlite3.Row
    pages = conn.execute(
        """SELECT page_name, page_id, components FROM display_model_layer
           WHERE model_id=? AND deleted_at IS NULL ORDER BY page_name""",
        (args.model,),
    ).fetchall()
    if not pages:
        raise SystemExit(f"未找到 model={args.model} 的活动页面")

    device_rows = conn.execute(
        """SELECT COUNT(*) AS total,
                  SUM(CASE WHEN status=1 THEN 1 ELSE 0 END) AS enabled
           FROM monitor_list
           WHERE project_uuid=? AND type=1 AND deleted_at IS NULL""",
        (args.project,),
    ).fetchone()

    report = {
        "project_uuid": args.project,
        "model_id": args.model,
        "devices": {"total": device_rows["total"], "enabled": device_rows["enabled"] or 0},
        "pages": [],
        "summary": Counter(),
        "invalid_bindings": [],
    }
    seen_bindings: set[tuple[str, str]] = set()

    for page in pages:
        cells = decode_components(page["components"])
        page_report = Counter(cells=len(cells))
        shapes = Counter()
        for cell in cells:
            detail = ((cell.get("data") or {}).get("detail") or {})
            shape = cell.get("shape") or detail.get("type") or "unknown"
            shapes[shape] += 1
            active = detail.get("active") or []
            data_bind = detail.get("dataBind") or []
            conditions = binding_conditions(detail)
            if active:
                page_report["active_entries"] += len(active)
            if data_bind:
                page_report["dataBind_entries"] += len(data_bind)
            if conditions:
                page_report["bound_cells"] += 1
                if "chart" in shape:
                    page_report["bound_charts"] += 1
                if "electric" in shape.lower():
                    page_report["bound_electric"] += 1
                if "arrow" in shape.lower():
                    page_report["bound_arrows"] += 1
            elif "electric" in shape.lower():
                page_report["unbound_electric"] += 1
            elif "arrow" in shape.lower():
                page_report["unbound_arrows"] += 1
            elif "chart" in shape:
                page_report["unbound_charts"] += 1

            for condition in conditions:
                key = (str(condition["deviceSN"]), str(condition["dataID"]))
                if key in seen_bindings:
                    continue
                seen_bindings.add(key)
                if not real_data_exists(conn, *key):
                    report["invalid_bindings"].append({
                        "page": page["page_name"],
                        "shape": shape,
                        "device_uuid": key[0],
                        "data_uuid": key[1],
                        "data_name": condition.get("dataName", ""),
                    })
        report["summary"].update(page_report)
        page_report["shapes"] = dict(shapes)
        page_report["page_name"] = page["page_name"]
        report["pages"].append(dict(page_report))

    report["summary"] = dict(report["summary"])
    report["summary"]["unique_bindings"] = len(seen_bindings)
    report["summary"]["invalid_binding_count"] = len(report["invalid_bindings"])
    print(json.dumps(report, ensure_ascii=False, indent=2))
    if args.json_out:
        Path(args.json_out).write_text(
            json.dumps(report, ensure_ascii=False, indent=2),
            encoding="utf-8",
        )
    conn.close()
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
