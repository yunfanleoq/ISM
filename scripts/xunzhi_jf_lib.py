#!/usr/bin/env python3
"""寻智 xunzhi-jf 脚本公共：路径、DB 预检、API、触发器热加载。"""
from __future__ import annotations

import hashlib
import json
import os
import sqlite3
import subprocess
import sys
from pathlib import Path

PROJECT_NAME = "寻智机房环境-xunzhi-jf"
LEAK_DEVICE = "xunzhi-jf-LEAK-01"
LEAK_TRIGGER = "漏水信号告警"
BASE = os.environ.get("BASE_URL", "http://127.0.0.1:8081").rstrip("/")


def resolve_root() -> Path:
    env = os.environ.get("XUNZHI_ISM_ROOT") or os.environ.get("ISM_ROOT")
    if env:
        return Path(env).resolve()
    return Path(__file__).resolve().parents[1]


def resolve_db(root: Path | None = None) -> Path:
    root = root or resolve_root()
    if os.environ.get("ISM_DB"):
        return Path(os.environ["ISM_DB"]).resolve()
    return (root / "ism_server_user" / "data" / "db" / "ism.db").resolve()


def curl_json(path: str, body=None, token: str | None = None, project_uuid: str | None = None) -> dict:
    cmd = ["curl", "-s", "-X", "POST", f"{BASE}{path}"]
    if token:
        cmd += ["-H", f"Authorization:{token}"]
    if project_uuid:
        cmd += ["-H", f"ProjectUuid:{project_uuid}"]
    if body is not None:
        cmd += ["-H", "Content-Type:application/json", "-d", json.dumps(body, ensure_ascii=False)]
    r = subprocess.run(cmd, capture_output=True, text=True, timeout=60)
    try:
        return json.loads(r.stdout or "{}")
    except json.JSONDecodeError:
        return {"code": -1, "raw": r.stdout, "stderr": r.stderr}


def login_token() -> str:
    pwd = hashlib.md5(b"123456").hexdigest()
    d = curl_json("/login", {"Username": "admin", "password": pwd})
    if d.get("code") != 1000:
        raise SystemExit(f"login failed: {d}")
    return d["data"]["token"]


def project_uuid(token: str) -> str | None:
    d = curl_json("/ProjectList", {}, token)
    for p in d.get("list", []):
        pi = p.get("ProjectInfo", p)
        if pi.get("name") == PROJECT_NAME:
            return pi.get("uuid")
    return None


def ism_server_cwd() -> str | None:
    try:
        r = subprocess.run(
            ["pgrep", "-f", "ism_server_user/ism_server"],
            capture_output=True,
            text=True,
            timeout=5,
        )
        pids = [x for x in (r.stdout or "").split() if x.strip()]
        if not pids:
            return None
        pid = pids[0]
        cwd = Path(f"/proc/{pid}/cwd")
        if cwd.is_symlink():
            return str(cwd.resolve())
    except (OSError, subprocess.SubprocessError):
        pass
    return None


def print_env(root: Path | None = None) -> dict:
    root = root or resolve_root()
    db = resolve_db(root)
    server_cwd = ism_server_cwd()
    expected_db = (root / "ism_server_user" / "data" / "db" / "ism.db").resolve()
    info = {
        "script_root": str(root),
        "ism_db": str(db),
        "ism_db_exists": db.is_file(),
        "ism_db_bytes": db.stat().st_size if db.is_file() else 0,
        "ism_server_cwd": server_cwd,
        "db_matches_script_tree": db == expected_db if db.is_file() else None,
        "server_cwd_matches_script_tree": (
            server_cwd is not None and Path(server_cwd).resolve() == (root / "ism_server_user").resolve()
        ),
    }
    print(json.dumps(info, indent=2, ensure_ascii=False))
    if server_cwd and Path(server_cwd).resolve() != (root / "ism_server_user").resolve():
        print(
            f"WARN: ism_server cwd={server_cwd} 与脚本树 {root}/ism_server_user 不一致（双目录串库风险）",
            file=sys.stderr,
        )
    if db.is_file() and db != expected_db:
        print(f"WARN: ISM_DB={db} 非本树默认库 {expected_db}", file=sys.stderr)
    return info


def db_preflight(db: Path, proj_uuid: str) -> dict:
    out = {
        "project_uuid": proj_uuid,
        "leak_device_rows": 0,
        "leak_trigger_rows": 0,
        "leak_point_is_alarm": None,
        "leak_point_alarm_on_value": None,
        "leak_point_value": None,
        "max_devices_alarm_list_id": 0,
    }
    if not db.is_file():
        return out
    conn = sqlite3.connect(db)
    cur = conn.cursor()
    cur.execute(
        "SELECT COUNT(*) FROM monitor_list WHERE project_uuid=? AND name=? AND type=1 AND deleted_at IS NULL",
        (proj_uuid, LEAK_DEVICE),
    )
    out["leak_device_rows"] = cur.fetchone()[0]
    cur.execute(
        "SELECT COUNT(*) FROM alarm_trigger WHERE project_uuid=? AND trigger_name=? AND deleted_at IS NULL",
        (proj_uuid, LEAK_TRIGGER),
    )
    out["leak_trigger_rows"] = cur.fetchone()[0]
    cur.execute(
        """
        SELECT d.is_alarm, d.alarm_on_value, d.value
        FROM device_real_data d
        JOIN monitor_list m ON m.uuid = d.device_uuid
        WHERE m.project_uuid=? AND m.name=? AND d.name='漏水信号' AND m.deleted_at IS NULL
        LIMIT 1
        """,
        (proj_uuid, LEAK_DEVICE),
    )
    row = cur.fetchone()
    if row:
        out["leak_point_is_alarm"], out["leak_point_alarm_on_value"], out["leak_point_value"] = row
    cur.execute("SELECT COALESCE(MAX(id),0) FROM devices_alarm_list")
    out["max_devices_alarm_list_id"] = cur.fetchone()[0]
    conn.close()
    return out


def reload_alarm_triggers(token: str, proj_uuid: str) -> int:
    """对每个触发器调用 AlarmTriggerEdit（等价于 UI 保存），触发 AlarmTriggerCloseChan 热加载。"""
    listed = curl_json("/GetAlarmTriggerList", {}, token, proj_uuid)
    triggers = listed.get("list") or []
    n = 0
    for t in triggers:
        body = dict(t)
        tid = body.get("ID") or body.get("id")
        if tid is None:
            continue
        body["ID"] = int(tid)
        body["project_uuid"] = proj_uuid
        resp = curl_json("/AlarmTriggerEdit", body, token, proj_uuid)
        if resp.get("code") == 0:
            n += 1
    return n


def diagnose_failure(root: Path | None, proj_uuid: str) -> None:
    root = root or resolve_root()
    db = resolve_db(root)
    print("=== xunzhi-jf 告警诊断 ===", file=sys.stderr)
    print_env(root)
    pf = db_preflight(db, proj_uuid)
    print(json.dumps(pf, indent=2, ensure_ascii=False), file=sys.stderr)
    if pf["leak_device_rows"] == 0:
        print(
            f"FAIL: 库中无设备 {LEAK_DEVICE}；请在本树执行 python3 scripts/import_xunzhi_jf_project.py",
            file=sys.stderr,
        )
    if pf["leak_trigger_rows"] == 0:
        print(f"FAIL: 库中无触发器「{LEAK_TRIGGER}」；请运行 xunzhi_jf_sync_alarm_db.py", file=sys.stderr)
    if pf.get("leak_point_value") == "1":
        print(
            "HINT: 漏水 DI 已为 1；若刚重启 ism_server，StartupAlarmDelay 默认约 10min 内抑警，"
            "须先 Mock 漏水=0 再打 0→1 边沿",
            file=sys.stderr,
        )
