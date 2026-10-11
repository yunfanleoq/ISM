#!/usr/bin/env python3
"""已导入项目：将 xunzhi-jf 告警点 alarm_on_value 置 1，并补齐包内 alarmTriggers（项目已存在无法整包重导时用）。"""
import json
import sqlite3
import subprocess
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
DB = ROOT / "ism_server_user" / "data" / "db" / "ism.db"
PKG = ROOT / "projects-import" / "xunzhi-jf" / "xunzhi-jf_complete_project_package.json"
PROJECT_NAME = "寻智机房环境-xunzhi-jf"
BASE = "http://127.0.0.1:8081"


def login_token():
    import hashlib

    pwd = hashlib.md5(b"123456").hexdigest()
    r = subprocess.run(
        [
            "curl",
            "-s",
            "-X",
            "POST",
            f"{BASE}/login",
            "-H",
            "Content-Type:application/json",
            "-d",
            json.dumps({"Username": "admin", "password": pwd}),
        ],
        capture_output=True,
        text=True,
        timeout=30,
    )
    d = json.loads(r.stdout or "{}")
    if d.get("code") != 1000:
        raise SystemExit(f"login failed: {d}")
    return d["data"]["token"]


def project_uuid(token):
    r = subprocess.run(
        [
            "curl",
            "-s",
            "-X",
            "POST",
            f"{BASE}/ProjectList",
            "-H",
            "Content-Type:application/json",
            "-H",
            f"Authorization:{token}",
            "-d",
            "{}",
        ],
        capture_output=True,
        text=True,
    )
    d = json.loads(r.stdout or "{}")
    for p in d.get("list", []):
        pi = p.get("ProjectInfo", p)
        if pi.get("name") == PROJECT_NAME:
            return pi.get("uuid")
    return None


def fix_sqlite(proj_uuid: str) -> int:
    if not DB.is_file():
        print(f"skip sqlite: no {DB}", file=sys.stderr)
        return 0
    conn = sqlite3.connect(DB)
    cur = conn.cursor()
    n = 0
    cur.execute(
        """
        UPDATE modbus_devices_data_model
        SET alarm_on_value = 1
        WHERE is_alarm = 1
          AND muid IN (
            SELECT muid FROM monitor_list
            WHERE project_uuid = ? AND type = 1 AND deleted_at IS NULL
          )
        """,
        (proj_uuid,),
    )
    n += cur.rowcount
    cur.execute(
        """
        UPDATE device_real_data
        SET alarm_on_value = 1
        WHERE is_alarm = 1 AND project_uuid = ?
        """,
        (proj_uuid,),
    )
    n += cur.rowcount
    conn.commit()
    conn.close()
    return n


def upsert_triggers(conn, proj_uuid: str) -> int:
    """按设备名+点名将包内触发器落到 DB（已导入项目 UUID 与包内不一致）。"""
    import uuid as uuidlib

    specs = [
        ("漏水信号告警", "xunzhi-jf-LEAK-01", "漏水信号", "=", "1", 1, "漏水报警", "漏水报警恢复"),
        ("温度过高告警", "xunzhi-jf-TH-01", "温度", ">", "35", 2, "温度超过35℃", "温度恢复正常"),
    ]
    cur = conn.cursor()
    added = 0
    for tname, dev_name, pt_name, cond, xval, level, show, hide in specs:
        cur.execute(
            "SELECT uuid, muid FROM monitor_list WHERE project_uuid=? AND name=? AND type=1 AND deleted_at IS NULL LIMIT 1",
            (proj_uuid, dev_name),
        )
        row = cur.fetchone()
        if not row:
            continue
        _dev_uuid, muid = row
        cur.execute(
            "SELECT uuid FROM modbus_devices_data_model WHERE muid=? AND name=? AND deleted_at IS NULL LIMIT 1",
            (muid, pt_name),
        )
        pt = cur.fetchone()
        if not pt:
            continue
        model_data_uuid = pt[0]
        cur.execute(
            "SELECT id FROM alarm_trigger WHERE project_uuid=? AND trigger_name=? AND deleted_at IS NULL",
            (proj_uuid, tname),
        )
        if cur.fetchone():
            continue
        cur.execute(
            """
            INSERT INTO alarm_trigger (
              uuid, trigger_name, project_uuid, trigger_device_type,
              trigger_device_model_uuid, trigger_model_data_uuid,
              trigger_alarm_hide_text, trigger_alarm_show_text,
              trigger_condition, trigger_x_value, trigger_y_value,
              trigger_alarm_level, trigger_keep_time, trigger_type,
              created_at, updated_at
            ) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,datetime('now'),datetime('now'))
            """,
            (
                uuidlib.uuid4().hex,
                tname,
                proj_uuid,
                2,
                muid,
                model_data_uuid,
                hide,
                show,
                cond,
                xval,
                "",
                level,
                0,
                0,
            ),
        )
        added += 1
    conn.commit()
    return added


def main():
    token = login_token()
    proj = project_uuid(token)
    if not proj:
        raise SystemExit(f"project not found: {PROJECT_NAME}")
    rows = fix_sqlite(proj)
    triggers = 0
    if DB.is_file():
        conn = sqlite3.connect(DB)
        triggers = upsert_triggers(conn, proj)
        conn.close()
    print(
        json.dumps(
            {"project_uuid": proj, "sqlite_alarm_on_value_rows": rows, "alarm_triggers_added": triggers},
            indent=2,
        )
    )
    if triggers == 0 and rows == 0:
        print("提示：无 DB 或未匹配设备；新库请 ImportProject 整包导入。", file=sys.stderr)


if __name__ == "__main__":
    main()
