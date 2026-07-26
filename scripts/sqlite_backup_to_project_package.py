#!/usr/bin/env python3
"""
将 ISM Sqlite3 全库备份还原后，导出为 /ImportProject 可用的项目包 JSON。

用法:
  python3 scripts/sqlite_backup_to_project_package.py \\
      --sql Sqlite3_Backup_xxx.sql \\
      --out projects-import/柴发监控/后沙峪改造-柴发部分_ISM项目包.json \\
      [--project-name "后沙峪改造-柴发部分"]

导入侧（客户机）:
  curl -X POST http://127.0.0.1:8081/ImportProject \\
    -H "Authorization: <admin_token>" \\
    -H "Content-Type: application/json" \\
    --data-binary @后沙峪改造-柴发部分_ISM项目包.json
"""

from __future__ import annotations

import argparse
import json
import os
import re
import sqlite3
import subprocess
import sys
import tempfile
from pathlib import Path
from typing import Any


def restore_sql_to_db(sql_path: Path, db_path: Path) -> None:
    if db_path.exists():
        db_path.unlink()
    # sqlite3 CLI 比 Python 执行大 SQL 更稳
    r = subprocess.run(
        ["sqlite3", str(db_path)],
        input=sql_path.read_bytes(),
        capture_output=True,
    )
    if r.returncode != 0:
        raise RuntimeError(
            f"sqlite3 restore failed: {r.stderr.decode('utf-8', errors='replace')[:2000]}"
        )


def q_all(conn: sqlite3.Connection, sql: str, args: tuple = ()) -> list[sqlite3.Row]:
    cur = conn.execute(sql, args)
    return cur.fetchall()


def q_one(conn: sqlite3.Connection, sql: str, args: tuple = ()) -> sqlite3.Row | None:
    cur = conn.execute(sql, args)
    return cur.fetchone()


def to_int(v: Any, default: int = 0) -> int:
    if v is None or v == "":
        return default
    try:
        return int(float(str(v)))
    except (TypeError, ValueError):
        return default


def normalize_modbus_extra(raw: str | None) -> str:
    """备份里常见 lowercase 'modbus'，导入引擎期望 'Modbus'。"""
    if not raw:
        return "{}"
    try:
        obj = json.loads(raw)
    except json.JSONDecodeError:
        return raw
    if not isinstance(obj, dict):
        return raw

    if "Modbus" in obj and isinstance(obj["Modbus"], dict):
        return json.dumps(obj, ensure_ascii=False, separators=(",", ":"))

    src = None
    if isinstance(obj.get("modbus"), dict):
        src = obj["modbus"]
    elif isinstance(obj.get("Modbus"), dict):
        src = obj["Modbus"]

    if src is None:
        return raw

    modbus = {
        "IPAddress": str(src.get("IPAddress") or src.get("ip") or "127.0.0.1"),
        "Port": str(src.get("Port") or src.get("port") or "502"),
        "address": str(src.get("address") or src.get("slaveId") or "1"),
        "packTime": int(src.get("packTime") or 200),
        "RegisterPack": int(src.get("RegisterPack") if src.get("RegisterPack") is not None else -1),
    }
    out = {k: v for k, v in obj.items() if k not in ("modbus", "Modbus")}
    out["Modbus"] = modbus
    return json.dumps(out, ensure_ascii=False, separators=(",", ":"))


def infer_gateway(conn: sqlite3.Connection) -> dict[str, Any]:
    rows = q_all(
        conn,
        """
        SELECT extra_data FROM monitor_list
        WHERE type=1 AND deleted_at IS NULL AND extra_data IS NOT NULL AND extra_data != ''
        LIMIT 20
        """,
    )
    ips, ports = [], []
    for r in rows:
        try:
            extra = json.loads(r["extra_data"] or "{}")
        except json.JSONDecodeError:
            continue
        m = extra.get("Modbus") or extra.get("modbus") or {}
        if m.get("IPAddress"):
            ips.append(str(m["IPAddress"]))
        if m.get("Port"):
            ports.append(to_int(m["Port"], 502))
    ip = ips[0] if ips else "127.0.0.1"
    port = ports[0] if ports else 502
    return {"ip": ip, "port": port}


def build_package(conn: sqlite3.Connection, project_name: str | None) -> dict[str, Any]:
    proj = q_one(
        conn,
        """
        SELECT uuid, name, description, industry, creator_uuid
        FROM project_lists
        WHERE deleted_at IS NULL
        ORDER BY id
        LIMIT 1
        """,
    )
    if not proj:
        raise RuntimeError("备份中未找到 project_lists 记录")

    name = (project_name or proj["name"] or "柴发监控").strip()
    gateway = infer_gateway(conn)

    models = q_all(
        conn,
        """
        SELECT uuid, name, described, type, gather_number, port, timeout,
               data_format, modbus_connect_type, modbus_connect_mode
        FROM devices_model
        WHERE deleted_at IS NULL AND project_uuid = ?
        ORDER BY id
        """,
        (proj["uuid"],),
    )
    model_uuids = {m["uuid"] for m in models}

    device_models = []
    for m in models:
        device_models.append(
            {
                "uuid": m["uuid"],
                "name": m["name"],
                "dec": m["described"] or m["name"],
                "type": to_int(m["type"], 2),
                "modbusConnectType": m["modbus_connect_type"] or "TCPClient",
                "modbusConnectMode": m["modbus_connect_mode"] or "TCP/IP",
                "DataFormat": m["data_format"] or "BigEndian",
                "timeout": to_int(m["timeout"], 2000),
                "port": to_int(m["port"], 502),
                "gatherNumber": to_int(m["gather_number"], 30),
            }
        )

    groups = q_all(
        conn,
        """
        SELECT uuid, muid, name, function, register_start, register_count
        FROM modbus_devices_register_group
        WHERE deleted_at IS NULL AND muid IN ({})
        ORDER BY id
        """.format(",".join("?" * len(model_uuids))),
        tuple(model_uuids),
    )
    register_groups = []
    group_uuids = set()
    for g in groups:
        group_uuids.add(g["uuid"])
        register_groups.append(
            {
                "uuid": g["uuid"],
                "muid": g["muid"],
                "name": g["name"],
                "function": to_int(g["function"], 3),
                "registerStart": to_int(g["register_start"], 0),
                "registerCount": to_int(g["register_count"], 1),
            }
        )

    points = q_all(
        conn,
        """
        SELECT uuid, muid, name, register_address, register_group_uuid, auth, type,
               byte_order, model_type, data_unit, conversion_expression,
               is_alarm, alarm_level, alarm_message, alarm_clear_message,
               is_record, record_type, record_interval, record_data_charge,
               record_data_timely, float_accuracy
        FROM modbus_devices_data_model
        WHERE deleted_at IS NULL AND muid IN ({})
        ORDER BY id
        """.format(",".join("?" * len(model_uuids))),
        tuple(model_uuids),
    )
    register_points = []
    for p in points:
        register_points.append(
            {
                "uuid": p["uuid"],
                "muid": p["muid"],
                "name": p["name"],
                "registerAddress": to_int(p["register_address"], 0),
                "registerGroupUuid": p["register_group_uuid"],
                "auth": p["auth"] or "ReadWrite",
                "type": p["type"] or "Float",
                "ByteOrder": p["byte_order"] or "",
                "modeltype": to_int(p["model_type"], 2),
                "unit": p["data_unit"] or "",
                "conversionExpression": p["conversion_expression"] or "",
                "alarm": to_int(p["is_alarm"], 0),
                "alarmLevel": to_int(p["alarm_level"], 0),
                "AlarmMessage": p["alarm_message"] or "",
                "AlarmClearMessage": p["alarm_clear_message"] or "",
                "record": to_int(p["is_record"], 0),
                "RecordType": to_int(p["record_type"], 0),
                "recordInterval": to_int(p["record_interval"], 0),
                "RecordDataCharge": p["record_data_charge"] or "",
                "RecordDataTimely": p["record_data_timely"] or "",
                "FloatAccuracy": p["float_accuracy"] or "",
            }
        )

    monitors = q_all(
        conn,
        """
        SELECT uuid, sid, pid, name, type, timeout, is_enable, interval, failed_times,
               described, offline_clear, offline_default_value, device_type, muid,
               configuration_uid, page_uuid, extra_data, status, longitude, latitude
        FROM monitor_list
        WHERE deleted_at IS NULL AND project_uuid = ?
        ORDER BY CASE WHEN pid=0 THEN 0 ELSE 1 END, sid
        """,
        (proj["uuid"],),
    )
    monitor_tree = []
    for n in monitors:
        monitor_tree.append(
            {
                "uuid": n["uuid"],
                "sid": to_int(n["sid"], 0),
                "pid": to_int(n["pid"], 0),
                "name": n["name"],
                "type": to_int(n["type"], 0),
                "timeout": to_int(n["timeout"], 5),
                "IsEnable": to_int(n["is_enable"], 1),
                "interval": to_int(n["interval"], 5),
                "failedTimes": to_int(n["failed_times"], 5),
                "description": n["described"] or "",
                "offlineClear": to_int(n["offline_clear"], 0),
                "offlineDefaultValue": n["offline_default_value"] or "0",
                "deviceType": to_int(n["device_type"], 0),
                "muid": n["muid"] or "",
                "configUid": n["configuration_uid"] or "",
                "PageUUID": n["page_uuid"] or "",
                "extra": normalize_modbus_extra(n["extra_data"]),
                "Status": to_int(n["status"], 0),
                "longitude": n["longitude"] or "",
                "latitude": n["latitude"] or "",
            }
        )

    # 告警触发器（若有）
    alarm_triggers = []
    try:
        triggers = q_all(
            conn,
            """
            SELECT * FROM alarm_trigger
            WHERE deleted_at IS NULL AND project_uuid = ?
            """,
            (proj["uuid"],),
        )
        for t in triggers:
            # 列名随版本可能不同，尽量宽松映射
            keys = set(t.keys())
            alarm_triggers.append(
                {
                    "Uuid": t["uuid"] if "uuid" in keys else "",
                    "TriggerName": t["trigger_name"] if "trigger_name" in keys else "",
                    "TriggerDeviceUuid": t["trigger_device_uuid"] if "trigger_device_uuid" in keys else "",
                    "TriggerDeviceName": t["trigger_device_name"] if "trigger_device_name" in keys else "",
                    "TriggerDataUuid": t["trigger_data_uuid"] if "trigger_data_uuid" in keys else "",
                    "TriggerDeviceType": to_int(t["trigger_device_type"] if "trigger_device_type" in keys else 0),
                    "TriggerDeviceModelUuid": t["trigger_device_model_uuid"]
                    if "trigger_device_model_uuid" in keys
                    else "",
                    "TriggerModelDataUuid": t["trigger_model_data_uuid"]
                    if "trigger_model_data_uuid" in keys
                    else "",
                    "TriggerAlarmHideText": t["trigger_alarm_hide_text"]
                    if "trigger_alarm_hide_text" in keys
                    else "",
                    "TriggerAlarmShowText": t["trigger_alarm_show_text"]
                    if "trigger_alarm_show_text" in keys
                    else "",
                    "TriggerCondition": t["trigger_condition"] if "trigger_condition" in keys else "",
                    "TriggerXValue": t["trigger_x_value"] if "trigger_x_value" in keys else "",
                    "TriggerYValue": t["trigger_y_value"] if "trigger_y_value" in keys else "",
                    "TriggerAlarmLevel": to_int(t["trigger_alarm_level"] if "trigger_alarm_level" in keys else 0),
                    "TriggerKeepTime": to_int(t["trigger_keep_time"] if "trigger_keep_time" in keys else 0),
                    "TriggerLinkDeviceType": to_int(
                        t["trigger_link_device_type"] if "trigger_link_device_type" in keys else 0
                    ),
                    "TriggerLinkdeviceModelUuid": t["trigger_linkdevice_model_uuid"]
                    if "trigger_linkdevice_model_uuid" in keys
                    else "",
                    "TriggerLinkModelDataUuid": t["trigger_link_model_data_uuid"]
                    if "trigger_link_model_data_uuid" in keys
                    else "",
                    "TriggerLinkageAlarmValue": t["trigger_linkage_alarm_value"]
                    if "trigger_linkage_alarm_value" in keys
                    else "",
                    "TriggerLinkageAlarmClearValue": t["trigger_linkage_alarm_clear_value"]
                    if "trigger_linkage_alarm_clear_value" in keys
                    else "",
                    "TriggerType": to_int(t["trigger_type"] if "trigger_type" in keys else 0),
                }
            )
    except sqlite3.Error:
        pass

    device_count = sum(1 for n in monitor_tree if n["type"] == 1)
    zone_count = sum(1 for n in monitor_tree if n["type"] == 0)

    pkg = {
        "_meta": {
            "format_version": "1.0",
            "generator": "sqlite_backup_to_project_package.py",
            "source_project_uuid": proj["uuid"],
            "source_project_name": proj["name"],
            "description": proj["description"] or "",
            "note": "通过 POST /ImportProject 导入，会新建独立项目（UUID 全部重新生成）",
        },
        "project": {
            "name": name,
            "gateway": gateway,
        },
        "deviceModels": device_models,
        "registerGroups": register_groups,
        "registerPoints": register_points,
        "monitorTree": monitor_tree,
        "alarmTriggers": alarm_triggers,
        "statistics": {
            "deviceModels": len(device_models),
            "registerGroups": len(register_groups),
            "registerPoints": len(register_points),
            "monitorTree": len(monitor_tree),
            "devices": device_count,
            "zones": zone_count,
            "alarmTriggers": len(alarm_triggers),
        },
    }
    return pkg


def validate_package(pkg: dict[str, Any]) -> list[str]:
    errors = []
    model_ids = {m["uuid"] for m in pkg["deviceModels"]}
    group_ids = {g["uuid"] for g in pkg["registerGroups"]}
    for g in pkg["registerGroups"]:
        if g["muid"] not in model_ids:
            errors.append(f"registerGroup {g['name']} muid 不在 deviceModels 中: {g['muid']}")
    orphan_pts = 0
    for p in pkg["registerPoints"]:
        if p["muid"] not in model_ids:
            orphan_pts += 1
        if p["registerGroupUuid"] not in group_ids:
            orphan_pts += 1
    if orphan_pts:
        errors.append(f"registerPoints 存在 {orphan_pts} 条孤儿引用")
    for n in pkg["monitorTree"]:
        if n["type"] == 1 and n["muid"] and n["muid"] not in model_ids:
            errors.append(f"设备 {n['name']} muid 不在 deviceModels 中")
    if not pkg["project"]["name"]:
        errors.append("project.name 为空")
    if not pkg["deviceModels"]:
        errors.append("deviceModels 为空")
    if not any(n["type"] == 1 for n in pkg["monitorTree"]):
        errors.append("monitorTree 中无设备节点")
    return errors


def main() -> int:
    ap = argparse.ArgumentParser(description="Sqlite3 备份 → ISM ImportProject 项目包")
    ap.add_argument("--sql", required=True, help="Sqlite3_Backup_*.sql 或 .zip")
    ap.add_argument("--out", required=True, help="输出 JSON 路径")
    ap.add_argument("--project-name", default="", help="覆盖项目名称")
    ap.add_argument("--pretty", action="store_true", help="美化缩进（文件更大）")
    args = ap.parse_args()

    sql_path = Path(args.sql).expanduser().resolve()
    out_path = Path(args.out).expanduser().resolve()
    out_path.parent.mkdir(parents=True, exist_ok=True)

    with tempfile.TemporaryDirectory(prefix="ism_chaifa_") as tmp:
        tmpdir = Path(tmp)
        work_sql = sql_path
        if sql_path.suffix.lower() == ".zip":
            import zipfile

            with zipfile.ZipFile(sql_path) as zf:
                names = [n for n in zf.namelist() if n.endswith(".sql")]
                if not names:
                    print("zip 内未找到 .sql", file=sys.stderr)
                    return 1
                zf.extract(names[0], tmpdir)
                work_sql = tmpdir / names[0]

        db_path = tmpdir / "source.db"
        print(f"[1/4] 还原 SQL → {db_path}")
        restore_sql_to_db(work_sql, db_path)

        print("[2/4] 读取并组装项目包")
        conn = sqlite3.connect(str(db_path))
        conn.row_factory = sqlite3.Row
        try:
            pkg = build_package(conn, args.project_name or None)
        finally:
            conn.close()

        print("[3/4] 校验引用完整性")
        errs = validate_package(pkg)
        if errs:
            for e in errs:
                print(f"  ERROR: {e}", file=sys.stderr)
            return 2

        print("[4/4] 写入", out_path)
        indent = 2 if args.pretty else None
        with open(out_path, "w", encoding="utf-8") as f:
            json.dump(pkg, f, ensure_ascii=False, indent=indent, separators=(",", ":") if not indent else None)

        st = pkg["statistics"]
        size_mb = out_path.stat().st_size / (1024 * 1024)
        print(
            f"完成: 模型={st['deviceModels']} 寄存器组={st['registerGroups']} "
            f"数据点={st['registerPoints']} 设备={st['devices']} "
            f"树节点={st['monitorTree']} 文件={size_mb:.2f}MB"
        )
        print(f"项目名: {pkg['project']['name']}")
        print(f"网关: {pkg['project']['gateway']}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
