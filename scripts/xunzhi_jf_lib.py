#!/usr/bin/env python3
"""寻智 xunzhi-jf 脚本公共：路径、DB 预检、API、触发器热加载、Mock/抑警窗。"""
from __future__ import annotations

import hashlib
import json
import os
import platform
import re
import socket
import sqlite3
import subprocess
import sys
import time
import uuid as uuidlib
from datetime import datetime, timezone
from pathlib import Path

PROJECT_NAME = "寻智机房环境-xunzhi-jf"
LEAK_DEVICE = "xunzhi-jf-LEAK-01"
LEAK_POINT = "漏水信号"
LEAK_TRIGGER = "漏水信号告警"
DEFAULT_STARTUP_DELAY_MIN = 10
BASE = os.environ.get("BASE_URL", "http://127.0.0.1:8081").rstrip("/")


def mock_port() -> int:
    return int(os.environ.get("MOCK_PORT", "1502"))


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


def tcp_port_listening(host: str = "127.0.0.1", port: int | None = None, timeout: float = 1.0) -> bool:
    port = port if port is not None else mock_port()
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.settimeout(timeout)
    try:
        sock.connect((host, port))
        return True
    except OSError:
        return False
    finally:
        sock.close()


def wait_tcp_port(
    host: str = "127.0.0.1",
    port: int | None = None,
    timeout_sec: float = 30.0,
    interval: float = 0.5,
) -> bool:
    port = port if port is not None else mock_port()
    deadline = time.time() + timeout_sec
    while time.time() < deadline:
        if tcp_port_listening(host, port):
            return True
        time.sleep(interval)
    return False


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


def ism_http_port() -> int:
    if os.environ.get("ISM_HTTP_PORT"):
        return int(os.environ["ISM_HTTP_PORT"])
    if BASE.endswith(":8081"):
        return 8081
    return 8081


def ism_server_pids() -> list[int]:
    pids: list[int] = []
    for cmd in (["pgrep", "-x", "ism_server"], ["pgrep", "-f", r"(^|/)ism_server( |$)"]):
        try:
            r = subprocess.run(cmd, capture_output=True, text=True, timeout=5)
            for line in (r.stdout or "").split():
                if line.isdigit():
                    pids.append(int(line))
        except (OSError, subprocess.SubprocessError):
            continue
    seen: set[int] = set()
    out: list[int] = []
    for p in pids:
        if p not in seen:
            seen.add(p)
            out.append(p)
    return out


def ism_server_listen_http() -> bool:
    return tcp_port_listening("127.0.0.1", ism_http_port())


def ism_server_running() -> bool:
    if ism_server_listen_http():
        return True
    return len(ism_server_pids()) > 0


def _pid_cwd(pid: int) -> str | None:
    try:
        proc_cwd = Path(f"/proc/{pid}/cwd")
        if proc_cwd.is_symlink():
            return str(proc_cwd.resolve())
    except OSError:
        pass
    try:
        r = subprocess.run(
            ["lsof", "-a", "-p", str(pid), "-d", "cwd", "-Fn"],
            capture_output=True,
            text=True,
            timeout=5,
        )
        if r.returncode == 0:
            for line in (r.stdout or "").splitlines():
                if line.startswith("n"):
                    return line[1:].strip()
    except (OSError, subprocess.SubprocessError, FileNotFoundError):
        pass
    return None


def ism_server_cwds() -> list[str]:
    cwds: list[str] = []
    for pid in ism_server_pids():
        cwd = _pid_cwd(pid)
        if cwd:
            cwds.append(cwd)
    return cwds


def ism_server_cwd() -> str | None:
    cwds = ism_server_cwds()
    return cwds[0] if cwds else None


def _parse_ps_etime_seconds(raw: str) -> int | None:
    s = raw.strip()
    if not s:
        return None
    days = 0
    if "-" in s:
        day_part, rest = s.split("-", 1)
        if not day_part.isdigit():
            return None
        days = int(day_part)
        s = rest
    parts = s.split(":")
    try:
        if len(parts) == 3:
            h, m, sec = int(parts[0]), int(parts[1]), int(parts[2])
        elif len(parts) == 2:
            h, m, sec = 0, int(parts[0]), int(parts[1])
        else:
            return None
        return days * 86400 + h * 3600 + m * 60 + sec
    except ValueError:
        return None


def ism_server_uptime_sec() -> int | None:
    pids = ism_server_pids()
    if not pids:
        return None
    pid = pids[0]
    attempts: list[list[str]] = []
    if platform.system() == "Darwin":
        attempts.append(["ps", "-o", "etime=", "-p", str(pid)])
        attempts.append(["ps", "-o", "etimes=", "-p", str(pid)])
    else:
        attempts.append(["ps", "-o", "etimes=", "-p", str(pid)])
        attempts.append(["ps", "-o", "etime=", "-p", str(pid)])
    for cmd in attempts:
        try:
            ps = subprocess.run(cmd, capture_output=True, text=True, timeout=5)
            if ps.returncode != 0:
                continue
            lines = [ln.strip() for ln in (ps.stdout or "").splitlines() if ln.strip()]
            if not lines:
                continue
            raw = lines[-1]
            if "etimes" in cmd[1]:
                if re.fullmatch(r"\d+", raw):
                    return int(raw)
                continue
            sec = _parse_ps_etime_seconds(raw)
            if sec is not None:
                return sec
        except (OSError, subprocess.SubprocessError, ValueError):
            continue
    return None


def mock_tcp_established_count(port: int | None = None) -> int:
    port = port if port is not None else mock_port()
    try:
        r = subprocess.run(
            ["ss", "-H", "-tn", "state", "established", f"( sport = :{port} or dport = :{port} )"],
            capture_output=True,
            text=True,
            timeout=5,
        )
        lines = [ln for ln in (r.stdout or "").splitlines() if ln.strip()]
        if lines:
            return len(lines)
    except (OSError, subprocess.SubprocessError, FileNotFoundError):
        pass
    try:
        r = subprocess.run(
            ["lsof", "-nP", f"-iTCP:{port}", "-sTCP:ESTABLISHED"],
            capture_output=True,
            text=True,
            timeout=5,
        )
        return max(0, len([ln for ln in (r.stdout or "").splitlines() if "ESTABLISHED" in ln]) - 1)
    except (OSError, subprocess.SubprocessError, FileNotFoundError):
        return 0


def read_startup_delay_minutes(db: Path, proj_uuid: str) -> int:
    """库内 StartupAlarmDelay；无记录时与后端默认一致为 10。"""
    if not db.is_file():
        return DEFAULT_STARTUP_DELAY_MIN
    conn = sqlite3.connect(db)
    cur = conn.cursor()
    cur.execute(
        """
        SELECT alarm_notice_params FROM alarm_notice
        WHERE project_uuid=? AND alarm_notice_type='StartupAlarmDelay' AND deleted_at IS NULL
        LIMIT 1
        """,
        (proj_uuid,),
    )
    row = cur.fetchone()
    conn.close()
    if not row or not row[0]:
        return DEFAULT_STARTUP_DELAY_MIN
    try:
        return int(json.loads(row[0]).get("DelayMinutes", DEFAULT_STARTUP_DELAY_MIN))
    except (json.JSONDecodeError, TypeError, ValueError):
        return DEFAULT_STARTUP_DELAY_MIN


def upsert_startup_delay_sqlite(db: Path, proj_uuid: str, delay_minutes: int) -> str:
    params = json.dumps({"DelayMinutes": delay_minutes}, ensure_ascii=False)
    conn = sqlite3.connect(db)
    cur = conn.cursor()
    cur.execute(
        """
        SELECT id FROM alarm_notice
        WHERE project_uuid=? AND alarm_notice_type='StartupAlarmDelay' AND deleted_at IS NULL
        LIMIT 1
        """,
        (proj_uuid,),
    )
    row = cur.fetchone()
    if row:
        cur.execute(
            "UPDATE alarm_notice SET alarm_notice_params=?, updated_at=datetime('now') WHERE id=?",
            (params, row[0]),
        )
        action = "updated"
    else:
        cur.execute(
            """
            INSERT INTO alarm_notice (uuid, alarm_notice_type, alarm_notice_params, project_uuid, created_at, updated_at)
            VALUES (?,?,?,?,datetime('now'),datetime('now'))
            """,
            (uuidlib.uuid4().hex, "StartupAlarmDelay", params, proj_uuid),
        )
        action = "inserted"
    conn.commit()
    conn.close()
    return action


def ensure_startup_delay_zero(db: Path, proj_uuid: str, token: str) -> dict:
    before = read_startup_delay_minutes(db, proj_uuid)
    sqlite_action = upsert_startup_delay_sqlite(db, proj_uuid, 0)
    params = json.dumps({"DelayMinutes": 0}, ensure_ascii=False)
    api = curl_json(
        "/UpdateAlarmNotice",
        {"type": "StartupAlarmDelay", "params": params},
        token,
        proj_uuid,
    )
    after = read_startup_delay_minutes(db, proj_uuid)
    return {
        "delay_minutes_before": before,
        "delay_minutes_after": after,
        "sqlite_action": sqlite_action,
        "api_code": api.get("code"),
        "note": (
            "已写入 DelayMinutes=0；须从本树重启 ism_server（Mock:1502 保持）使抑警窗关闭。"
            "批量 Modbus 重连也会 ReconfigureStartupAlarmWindows。"
        ),
    }


def startup_suppression_remaining_sec(db: Path, proj_uuid: str, server_uptime: int | None) -> int:
    """
    服务端启动时若仍为 10min 抑警窗，ObserveStartupAlarm 会吞入列。
    DB 已改为 0 但进程未重启时，按「启动时默认 10min」估算剩余秒数。
    """
    if server_uptime is None:
        return 0
    delay_db = read_startup_delay_minutes(db, proj_uuid)
    if delay_db <= 0 and server_uptime < 20:
        return 0
    if server_uptime >= DEFAULT_STARTUP_DELAY_MIN * 60:
        return 0
    effective = DEFAULT_STARTUP_DELAY_MIN
    return max(0, effective * 60 - server_uptime)


def leak_realdata_freshness(db: Path, proj_uuid: str, max_age_sec: int = 45) -> dict:
    out = {
        "value": None,
        "updated_at": None,
        "age_sec": None,
        "fresh": False,
        "monitor_status": None,
    }
    if not db.is_file():
        return out
    conn = sqlite3.connect(db)
    cur = conn.cursor()
    cur.execute(
        """
        SELECT d.value, d.updated_at, m.status
        FROM device_real_data d
        JOIN monitor_list m ON m.uuid = d.device_uuid
        WHERE m.project_uuid=? AND m.name=? AND d.name=? AND m.deleted_at IS NULL
        LIMIT 1
        """,
        (proj_uuid, LEAK_DEVICE, LEAK_POINT),
    )
    row = cur.fetchone()
    conn.close()
    if not row:
        return out
    out["value"], updated_at, out["monitor_status"] = row[0], row[1], row[2]
    out["updated_at"] = updated_at
    if updated_at:
        try:
            if "T" in updated_at:
                dt = datetime.fromisoformat(updated_at.replace("Z", "+00:00"))
            else:
                dt = datetime.strptime(updated_at, "%Y-%m-%d %H:%M:%S").replace(tzinfo=timezone.utc)
            age = (datetime.now(timezone.utc) - dt.astimezone(timezone.utc)).total_seconds()
            out["age_sec"] = int(age)
            out["fresh"] = age <= max_age_sec
        except (ValueError, TypeError):
            out["fresh"] = False
    return out


def wait_leak_realdata_fresh(
    db: Path,
    proj_uuid: str,
    max_age_sec: int = 45,
    timeout_sec: float = 90.0,
) -> dict:
    deadline = time.time() + timeout_sec
    last: dict = {}
    while time.time() < deadline:
        last = leak_realdata_freshness(db, proj_uuid, max_age_sec)
        if last.get("fresh"):
            return last
        time.sleep(3)
    return last


def print_env(root: Path | None = None, proj_uuid: str | None = None, strict_mock: bool = False) -> dict:
    root = root or resolve_root()
    db = resolve_db(root)
    port = mock_port()
    server_cwd = ism_server_cwd()
    uptime = ism_server_uptime_sec()
    expected_db = (root / "ism_server_user" / "data" / "db" / "ism.db").resolve()
    mock_ok = tcp_port_listening(port=port)
    info: dict = {
        "script_root": str(root),
        "ism_db": str(db),
        "ism_db_exists": db.is_file(),
        "ism_db_bytes": db.stat().st_size if db.is_file() else 0,
        "ism_http_port": ism_http_port(),
        "ism_server_listen_http": ism_server_listen_http(),
        "ism_server_running": ism_server_running(),
        "ism_server_pids": ism_server_pids(),
        "ism_server_cwd": server_cwd,
        "ism_server_cwds": ism_server_cwds(),
        "ism_server_uptime_sec": uptime,
        "mock_host": "127.0.0.1",
        "mock_port": port,
        "mock_port_listening": mock_ok,
        "mock_tcp_established": mock_tcp_established_count(port),
        "db_matches_script_tree": db == expected_db if db.is_file() else None,
        "server_cwd_matches_script_tree": (
            None
            if server_cwd is None
            else Path(server_cwd).resolve() == (root / "ism_server_user").resolve()
        ),
    }
    if proj_uuid and db.is_file():
        info["startup_delay_minutes"] = read_startup_delay_minutes(db, proj_uuid)
        info["startup_suppression_remaining_sec"] = startup_suppression_remaining_sec(db, proj_uuid, uptime)
        info["leak_realdata"] = leak_realdata_freshness(db, proj_uuid)
    print(json.dumps(info, indent=2, ensure_ascii=False))
    if not ism_server_running():
        msg = f"WARN: ism_server 未检测到（HTTP:{ism_http_port()} LISTEN 或 pgrep ism_server）"
        print(msg, file=sys.stderr)
    if not mock_ok:
        msg = f"FAIL: Mock 未监听 127.0.0.1:{port}（ISM 会 dial refused，无边沿）"
        if strict_mock:
            raise SystemExit(msg)
        print(f"WARN: {msg}", file=sys.stderr)
    if (
        ism_server_running()
        and server_cwd
        and Path(server_cwd).resolve() != (root / "ism_server_user").resolve()
    ):
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
    print_env(root, proj_uuid, strict_mock=False)
    pf = db_preflight(db, proj_uuid)
    print(json.dumps(pf, indent=2, ensure_ascii=False), file=sys.stderr)
    if pf["leak_device_rows"] == 0:
        print(
            f"FAIL: 库中无设备 {LEAK_DEVICE}；请在本树执行 python3 scripts/import_xunzhi_jf_project.py",
            file=sys.stderr,
        )
    if pf["leak_trigger_rows"] == 0:
        print(f"FAIL: 库中无触发器「{LEAK_TRIGGER}」；请运行 xunzhi_jf_sync_alarm_db.py", file=sys.stderr)
    fr = leak_realdata_freshness(db, proj_uuid)
    if not fr.get("fresh"):
        print(
            f"HINT: LEAK realdata 未刷新（value={fr.get('value')} age={fr.get('age_sec')}s）；"
            "确认 Mock:1502 LISTEN 且 ism_server 已采集（见 cycle_log_th_leak）",
            file=sys.stderr,
        )
