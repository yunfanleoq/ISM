#!/usr/bin/env python3
"""
寻智双口 / 本机 ISM：确保业务库存在可登录的 admin（密码链：前端 MD5(123456) → bcrypt 比对）。

用法（在仓库根或 ISM源码-xunzhi-jf 根）:
  python3 scripts/xunzhi_jf_seed_admin.py
  python3 scripts/xunzhi_jf_seed_admin.py --verify-login
  ISM_DB=/path/ism.db python3 scripts/xunzhi_jf_seed_admin.py
  BASE_URL=http://127.0.0.1:8082 python3 scripts/xunzhi_jf_seed_admin.py --verify-login

读取 ism_server_user/conf/app.conf 的 dbtype：
  1 = SQLite（默认 data/db/ism.db，可用 ISM_DB / XUNZHI_ISM_ROOT 覆盖）
  4 = OceanBase（MySQL 协议，需 pymysql）

与 ism_server 启动种子一致：password = bcrypt(MD5("123456"))。
"""
from __future__ import annotations

import argparse
import configparser
import json
import sqlite3
import subprocess
import sys
import uuid
from datetime import datetime, timezone
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
DEFAULT_CONF = ROOT / "ism_server_user" / "conf" / "app.conf"

# models/db.go 单机初始化同款 hash（bcrypt(MD5("123456"))）
ADMIN_BCRYPT = (
    "$2a$10$h9swLjbTTcSVUCqQDt6nAetw.FVRLPE0WPDzqloprYRO7PDtLC5Ii"
)
LOGIN_MD5_123456 = "e10adc3949ba59abbe56e057f20f883e"

DEFAULT_ROLES = (
    ("Admin", "管理员", "add,edit,delete,view"),
    ("Operator", "操作员", "add,edit,delete,view"),
    ("User", "普通用户", "view"),
)


def load_app_conf(conf_path: Path) -> dict:
    if not conf_path.exists():
        return {}
    text = conf_path.read_text(encoding="utf-8", errors="ignore")
    parser = configparser.ConfigParser()
    parser.optionxform = str
    parser.read_string("[app]\n" + text)
    return {k.lower(): v for k, v in parser.items("app")}


def resolve_sqlite_path(conf: dict) -> Path:
    import os

    env_db = os.environ.get("ISM_DB")
    if env_db:
        return Path(env_db).resolve()
    env_root = os.environ.get("XUNZHI_ISM_ROOT") or os.environ.get("ISM_ROOT")
    if env_root:
        return (Path(env_root) / "ism_server_user" / "data" / "db" / "ism.db").resolve()
    return (ROOT / "ism_server_user" / "data" / "db" / "ism.db").resolve()


def connect_db(conf: dict, dbtype: int):
    if dbtype == 1:
        db_path = resolve_sqlite_path(conf)
        if not db_path.parent.exists():
            db_path.parent.mkdir(parents=True, exist_ok=True)
        conn = sqlite3.connect(str(db_path))
        conn.row_factory = sqlite3.Row
        return conn, "sqlite", str(db_path)

    try:
        import pymysql
    except ImportError as exc:
        raise SystemExit("OceanBase/MySQL 种子需要 pymysql：pip3 install pymysql") from exc

    if dbtype == 4:
        host = conf.get("oceanbasehost", "127.0.0.1")
        port = int(conf.get("oceanbaseport", "2881"))
        user = conf.get("oceanbaseuser", "root")
        password = conf.get("oceanbasepwd", "")
        database = conf.get("oceanbasedbname", "ism")
        label = f"mysql://{user}@{host}:{port}/{database}"
    else:
        host = conf.get("mysqlhost", "127.0.0.1")
        port = int(conf.get("mysqlport", "3306"))
        user = conf.get("mysqluser", "root")
        password = conf.get("mysqlpwd", "")
        database = conf.get("mysqldbname", "ism")
        label = f"mysql://{user}@{host}:{port}/{database}"

    conn = pymysql.connect(
        host=host,
        port=port,
        user=user,
        password=password,
        database=database,
        charset="utf8mb4",
        autocommit=False,
        connect_timeout=15,
    )
    return conn, "mysql", label


def utc_now_sql() -> str:
    return datetime.now(timezone.utc).strftime("%Y-%m-%d %H:%M:%S.%f")[:-3] + "+00:00"


def seed_roles(conn, dialect: str) -> int:
    created = 0
    cur = conn.cursor()
    for role_id, role_name, permission in DEFAULT_ROLES:
        if dialect == "sqlite":
            cur.execute(
                "SELECT id FROM roles_list WHERE role_id = ? AND deleted_at IS NULL LIMIT 1",
                (role_id,),
            )
        else:
            cur.execute(
                "SELECT id FROM roles_list WHERE role_id = %s AND deleted_at IS NULL LIMIT 1",
                (role_id,),
            )
        if cur.fetchone():
            continue
        now = utc_now_sql()
        if dialect == "sqlite":
            cur.execute(
                """
                INSERT INTO roles_list (created_at, updated_at, role_name, role_id, permission)
                VALUES (?, ?, ?, ?, ?)
                """,
                (now, now, role_name, role_id, permission),
            )
        else:
            cur.execute(
                """
                INSERT INTO roles_list (created_at, updated_at, role_name, role_id, permission)
                VALUES (%s, %s, %s, %s, %s)
                """,
                (now, now, role_name, role_id, permission),
            )
        created += 1
    return created


def seed_admin(conn, dialect: str, reset_password: bool) -> str:
    cur = conn.cursor()
    if dialect == "sqlite":
        cur.execute(
            "SELECT id, password FROM user WHERE username = ? AND deleted_at IS NULL LIMIT 1",
            ("admin",),
        )
    else:
        cur.execute(
            "SELECT id, password FROM `user` WHERE username = %s AND deleted_at IS NULL LIMIT 1",
            ("admin",),
        )
    row = cur.fetchone()
    now = utc_now_sql()
    if row:
        if reset_password:
            user_id = row[0]
            if dialect == "sqlite":
                cur.execute(
                    "UPDATE user SET password = ?, updated_at = ? WHERE id = ?",
                    (ADMIN_BCRYPT, now, user_id),
                )
            else:
                cur.execute(
                    "UPDATE `user` SET password = %s, updated_at = %s WHERE id = %s",
                    (ADMIN_BCRYPT, now, user_id),
                )
            return "updated_password"
        return "exists"
    admin_uuid = str(uuid.uuid4())
    if dialect == "sqlite":
        cur.execute(
            """
            INSERT INTO user (created_at, updated_at, username, password, name, role, uuid)
            VALUES (?, ?, 'admin', ?, '超级管理员', 'Admin', ?)
            """,
            (now, now, ADMIN_BCRYPT, admin_uuid),
        )
    else:
        cur.execute(
            """
            INSERT INTO `user` (created_at, updated_at, username, password, name, role, uuid)
            VALUES (%s, %s, 'admin', %s, '超级管理员', 'Admin', %s)
            """,
            (now, now, ADMIN_BCRYPT, admin_uuid),
        )
    return "created"


def verify_login(base_url: str) -> dict:
    url = base_url.rstrip("/") + "/login"
    payload = json.dumps({"Username": "admin", "password": LOGIN_MD5_123456})
    proc = subprocess.run(
        [
            "curl",
            "-s",
            "-m",
            "15",
            "-X",
            "POST",
            url,
            "-H",
            "Content-Type: application/json",
            "-d",
            payload,
        ],
        capture_output=True,
        text=True,
    )
    try:
        return json.loads(proc.stdout or "{}")
    except json.JSONDecodeError:
        return {"raw": proc.stdout, "stderr": proc.stderr, "code": -1}


def main() -> int:
    parser = argparse.ArgumentParser(description="确保 ISM 库存在 admin 用户（寻智双口复测）")
    parser.add_argument("--conf", type=Path, default=DEFAULT_CONF, help="app.conf 路径")
    parser.add_argument("--dbtype", type=int, help="覆盖 app.conf 的 dbtype")
    parser.add_argument(
        "--reset-password",
        action="store_true",
        help="admin 已存在时仍重置为 bcrypt(MD5(123456))",
    )
    parser.add_argument(
        "--verify-login",
        action="store_true",
        help="种子后 curl POST /login（需 ism_server 已监听）",
    )
    parser.add_argument(
        "--base-url",
        default=None,
        help="登录验证根 URL，默认从 app.conf httpport 或 BASE_URL 环境变量",
    )
    args = parser.parse_args()

    conf = load_app_conf(args.conf)
    dbtype = args.dbtype if args.dbtype is not None else int(str(conf.get("dbtype", "1")).strip())

    conn, dialect, db_label = connect_db(conf, dbtype)
    try:
        roles_n = seed_roles(conn, dialect)
        admin_action = seed_admin(conn, dialect, args.reset_password)
        conn.commit()
    finally:
        conn.close()

    http_port = conf.get("httpport", "8081")
    base_url = (
        args.base_url
        or __import__("os").environ.get("BASE_URL")
        or f"http://127.0.0.1:{http_port}"
    )

    out = {
        "dbtype": dbtype,
        "database": db_label,
        "conf": str(args.conf.resolve()),
        "admin": admin_action,
        "roles_inserted": roles_n,
        "login_md5_hint": LOGIN_MD5_123456,
    }

    if args.verify_login:
        out["login"] = verify_login(base_url)

    print(json.dumps(out, ensure_ascii=False, indent=2))

    if args.verify_login:
        code = out.get("login", {}).get("code")
        if code != 1000:
            print(
                f"WARN: 登录未成功 (code={code})；确认 ism_server 使用的 dbtype/库路径与本次种子一致。",
                file=sys.stderr,
            )
            return 2
    return 0


if __name__ == "__main__":
    sys.exit(main())
