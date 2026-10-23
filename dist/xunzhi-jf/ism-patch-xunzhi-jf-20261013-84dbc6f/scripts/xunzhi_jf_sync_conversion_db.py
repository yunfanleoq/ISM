#!/usr/bin/env python3
"""将 Modbus 点表 conversion_expression 同步到 device_real_data（寻智 xunzhi-jf 温湿度工程换算）。"""
import sqlite3
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(ROOT / "scripts"))
from xunzhi_jf_lib import login_token, project_uuid, resolve_db  # noqa: E402

NAMES = ("温度", "湿度", "漏水位置米数")


def main():
    db = resolve_db(ROOT)
    token = login_token()
    proj = project_uuid(token)
    if not proj:
        raise SystemExit("无法解析寻智 project_uuid（请先启动 ism_server 并 ImportProject）")
    conn = sqlite3.connect(db)
    cur = conn.cursor()
    updated = 0
    for name in NAMES:
        cur.execute(
            """
            SELECT m.uuid, m.conversion_expression
            FROM modbus_devices_data_model m
            JOIN device_real_data d ON d.model_data_uuid = m.uuid
            WHERE d.project_uuid = ? AND d.name = ?
            LIMIT 1
            """,
            (proj, name),
        )
        row = cur.fetchone()
        if not row:
            print(f"skip: no row for {name}")
            continue
        _, expr = row
        cur.execute(
            """
            UPDATE device_real_data
            SET conversion_expression = ?
            WHERE project_uuid = ? AND name = ?
            """,
            (expr, proj, name),
        )
        updated += cur.rowcount
        print(f"sync {name} -> {expr} ({cur.rowcount} rows)")
    conn.commit()
    conn.close()
    print(f"OK: updated {updated} device_real_data rows in {db}")


if __name__ == "__main__":
    main()
