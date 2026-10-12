#!/usr/bin/env python3
"""已导入项目：对齐 alarm_on_value / 补 alarm_trigger，并 API 热加载触发器 map。

演示须固定单一仓库根（XUNZHI_ISM_ROOT 或脚本上级目录），与 ism_server cwd 一致。
"""
import json
import sqlite3
import sys
import uuid as uuidlib
from pathlib import Path

_SCRIPTS = Path(__file__).resolve().parent
sys.path.insert(0, str(_SCRIPTS))

from xunzhi_jf_lib import (  # noqa: E402
    LEAK_TRIGGER,
    PROJECT_NAME,
    db_preflight,
    diagnose_failure,
    ensure_startup_delay_zero,
    login_token,
    print_env,
    project_uuid,
    reload_alarm_triggers,
    resolve_db,
    resolve_root,
    ism_server_uptime_sec,
    startup_suppression_remaining_sec,
)

def fix_sqlite(db: Path, proj_uuid: str) -> int:
    if not db.is_file():
        raise SystemExit(f"ism.db 不存在: {db}（请确认 XUNZHI_ISM_ROOT / ISM_DB 与 ism_server cwd）")
    conn = sqlite3.connect(db)
    cur = conn.cursor()
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
    n = cur.rowcount
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


def upsert_triggers(db: Path, proj_uuid: str) -> int:
    specs = [
        (LEAK_TRIGGER, "xunzhi-jf-LEAK-01", "漏水信号", "=", "1", 1, "漏水报警", "漏水报警恢复"),
        ("温度过高告警", "xunzhi-jf-TH-01", "温度", ">", "35", 2, "温度超过35℃", "温度恢复正常"),
    ]
    conn = sqlite3.connect(db)
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
    conn.close()
    return added


def main():
    import argparse

    ap = argparse.ArgumentParser()
    ap.add_argument("--print-env", action="store_true", help="打印 cwd / ism.db / Mock:1502 等后退出")
    ap.add_argument(
        "--strict-mock",
        action="store_true",
        help="与 --print-env 合用：Mock 未监听 1502 时 exit 1",
    )
    ap.add_argument("--skip-reload", action="store_true", help="跳过 AlarmTriggerEdit 热加载（改库后须重启 ism_server）")
    ap.add_argument("--skip-startup-delay-zero", action="store_true", help="不写入 StartupAlarmDelay DelayMinutes=0")
    args = ap.parse_args()

    root = resolve_root()
    db = resolve_db(root)
    if args.print_env:
        proj = None
        try:
            token = login_token()
            proj = project_uuid(token)
        except SystemExit:
            pass
        print_env(root, proj, strict_mock=args.strict_mock)
        return

    print_env(root, strict_mock=args.strict_mock)
    token = login_token()
    proj = project_uuid(token)
    if not proj:
        raise SystemExit(
            f"API 未找到项目「{PROJECT_NAME}」；请在本树执行: python3 scripts/import_xunzhi_jf_project.py"
        )

    startup_delay = None
    if not args.skip_startup_delay_zero:
        startup_delay = ensure_startup_delay_zero(db, proj, token)

    before = db_preflight(db, proj)
    rows = fix_sqlite(db, proj)
    triggers_added = upsert_triggers(db, proj)
    after = db_preflight(db, proj)
    remaining = startup_suppression_remaining_sec(db, proj, ism_server_uptime_sec())

    reloaded = 0
    reload_note = ""
    if not args.skip_reload:
        reloaded = reload_alarm_triggers(token, proj)
        if reloaded == 0 and after["leak_trigger_rows"] > 0:
            reload_note = "reload_api_returned_0; 若仍不入列请从 ism_server_user 目录重启 ism_server"
    else:
        reload_note = "skipped_reload; 必须从 ism_server_user 目录重启 ism_server"

    out = {
        "project_uuid": proj,
        "ism_db": str(db),
        "sqlite_alarm_on_value_rows": rows,
        "alarm_triggers_added": triggers_added,
        "alarm_triggers_reloaded_via_api": reloaded,
        "startup_delay_zero": startup_delay,
        "startup_suppression_remaining_sec": remaining,
        "preflight_before": before,
        "preflight_after": after,
        "note": reload_note,
    }
    print(json.dumps(out, indent=2, ensure_ascii=False))

    if after["leak_device_rows"] == 0:
        diagnose_failure(root, proj)
        raise SystemExit(2)
    if triggers_added == 0 and rows == 0 and after["leak_trigger_rows"] == 0:
        print("提示：无匹配设备/触发器；新库请 ImportProject。", file=sys.stderr)
        raise SystemExit(2)


if __name__ == "__main__":
    main()
