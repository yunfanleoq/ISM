#!/usr/bin/env python3
"""ISM 项目级备份/恢复 共享核心库 (DRY)。

backup_project.py 和 restore_project.py 都依赖本模块，集中处理:
  - app.conf 解析与 dbtype 判定 (1=SQLite / 4=OceanBase / 0=MySQL)
  - 统一的数据库句柄抽象 (pymysql ↔ sqlite3，屏蔽占位符/标识符差异)
  - 「项目作用域」解析 (按 projectUUID/displayUUID 圈定所有相关表与行)
  - 行数据 JSON 序列化/反序列化 (datetime/Decimal/bytes)

设计原则: 不假设各表 muid/project_uuid 一定自洽 (该项目曾被误删后部分回填，
数据层可能不一致)，因此作用域用「并集」方式尽量把相关行都纳入备份。
"""
from __future__ import annotations

import base64
import os
import re
from datetime import date, datetime
from decimal import Decimal
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parents[1]
DEFAULT_CONF = REPO_ROOT / "ism_server_user" / "conf" / "app.conf"
DEFAULT_SQLITE = REPO_ROOT / "ism_server_user" / "data" / "db" / "ism.db"

# 本项目的固定标识 (航信机房 SCADA 大屏)
PROJECT_UUID = "31bc90be-ebc4-dd61-ba9d-ce6e075e40e2"
DISPLAY_UUID = "043135ad-44be-e5d8-89be-3e54883c23a8"
PROJECT_LABEL = "航信机房"

# 大表分类: 时序/日志默认不纳入「最小 git 集合」(见 backup_project.py)
HISTORY_PATTERNS = (re.compile(r"history", re.I), re.compile(r"alarm_list$", re.I))
LOG_TABLES = frozenset({"system_journal"})

# ───────────────────────── app.conf ─────────────────────────


def parse_conf(path: Path | str = DEFAULT_CONF) -> dict:
    """解析 beego 风格 key=value 配置 (忽略空行/注释)。"""
    conf: dict[str, str] = {}
    with open(path, "r", encoding="utf-8") as fh:
        for line in fh:
            line = line.strip()
            if not line or line.startswith("#") or "=" not in line:
                continue
            key, _, val = line.partition("=")
            conf[key.strip()] = val.strip()
    return conf


# ───────────────────────── DB 抽象 ─────────────────────────


class Db:
    """统一的最小数据库句柄。SQL 一律用 ``?`` 占位符 + ``{q}`` 标识符占位，
    由本类按方言转换，从而 backup/restore 逻辑无需关心底层是 MySQL 还是 SQLite。
    """

    def __init__(self, kind: str, conn):
        self.kind = kind  # 'mysql' | 'sqlite'
        self.conn = conn

    # -- 方言细节 --
    @property
    def placeholder(self) -> str:
        return "%s" if self.kind == "mysql" else "?"

    def quote(self, ident: str) -> str:
        return f"`{ident}`" if self.kind == "mysql" else f'"{ident}"'

    def _render(self, sql: str) -> str:
        """把中性 SQL 里的 ``?`` 占位符与 ``{q}name{/q}`` 标识符转成方言形式。"""
        sql = re.sub(r"\{q\}(.*?)\{/q\}", lambda m: self.quote(m.group(1)), sql)
        if self.kind == "mysql":
            sql = sql.replace("?", "%s")
        return sql

    # -- 执行 --
    def query(self, sql: str, params=()) -> list[dict]:
        cur = self.conn.cursor()
        cur.execute(self._render(sql), tuple(params))
        rows = cur.fetchall()
        if self.kind == "sqlite":
            rows = [dict(r) for r in rows]
        cur.close()
        return rows

    def execute(self, sql: str, params=()) -> int:
        cur = self.conn.cursor()
        cur.execute(self._render(sql), tuple(params))
        n = cur.rowcount
        cur.close()
        return n

    def executemany(self, sql: str, seq) -> int:
        cur = self.conn.cursor()
        cur.executemany(self._render(sql), list(seq))
        n = cur.rowcount
        cur.close()
        return n

    def list_tables(self) -> list[str]:
        if self.kind == "mysql":
            rows = self.query("SHOW TABLES")
            return sorted(list(r.values())[0] for r in rows)
        rows = self.query(
            "SELECT name FROM sqlite_master WHERE type='table' "
            "AND name NOT LIKE 'sqlite_%' ORDER BY name"
        )
        return [r["name"] for r in rows]

    def columns(self, table: str) -> list[str]:
        if self.kind == "mysql":
            rows = self.query(
                "SELECT COLUMN_NAME FROM information_schema.COLUMNS "
                "WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME=? ORDER BY ORDINAL_POSITION",
                (table,),
            )
            return [r["COLUMN_NAME"] for r in rows]
        rows = self.query(f"PRAGMA table_info({self.quote(table)})")
        return [r["name"] for r in rows]

    def commit(self):
        self.conn.commit()

    def close(self):
        try:
            self.conn.close()
        except Exception:
            pass


def connect(conf: dict) -> Db:
    """按 app.conf 的 dbtype 建立连接。支持 4=OceanBase / 0=MySQL / 1=SQLite。"""
    dbtype = int(conf.get("dbtype", "4"))
    if dbtype == 1:
        import sqlite3

        sqlite_path = os.environ.get("ISM_SQLITE", str(DEFAULT_SQLITE))
        conn = sqlite3.connect(sqlite_path)
        conn.row_factory = sqlite3.Row
        return Db("sqlite", conn)

    import pymysql

    if dbtype == 4:
        cfg = dict(
            host=conf.get("oceanbasehost", "127.0.0.1"),
            port=int(conf.get("oceanbaseport", "2881")),
            user=conf.get("oceanbaseuser", "root@ism_tenant"),
            password=conf.get("oceanbasepwd", ""),
            database=conf.get("oceanbasedbname", "ism"),
        )
    elif dbtype == 0:
        cfg = dict(
            host=conf.get("mysqlhost", "127.0.0.1"),
            port=int(conf.get("mysqlport", "3306")),
            user=conf.get("mysqluser", "root"),
            password=conf.get("mysqlpwd", ""),
            database=conf.get("mysqldbname", "ism"),
        )
    else:
        raise NotImplementedError(
            f"dbtype={dbtype} 暂不支持本脚本 (仅支持 1=SQLite, 0=MySQL, 4=OceanBase)"
        )
    conn = pymysql.connect(charset="utf8mb4", cursorclass=pymysql.cursors.DictCursor, **cfg)
    return Db("mysql", conn)


# ───────────────────────── 作用域解析 ─────────────────────────


class TablePlan:
    """描述某张表如何按项目作用域筛选 (用结构化 scope，跨方言可重建 WHERE)。"""

    def __init__(self, table: str, columns: list[str], kind: str, category: str):
        self.table = table
        self.columns = columns
        self.kind = kind  # uuid_eq_project|project_uuid|model_id|muid|devices_model|device_real
        self.category = category  # config|history|log


def compute_scope(db: Db, project_uuid: str, display_uuid: str) -> dict:
    """收集与项目关联的 muid 集合与设备 uuid 集合 (并集，尽量全)。"""
    muids: set[str] = set()
    rows = db.query(
        "SELECT DISTINCT muid FROM {q}monitor_list{/q} "
        "WHERE project_uuid=? AND muid IS NOT NULL AND muid<>'' "
        "AND deleted_at IS NULL",
        (project_uuid,),
    )
    muids.update(r["muid"] for r in rows)
    # 该项目自有的 devices_model.uuid 也作为 muid 纳入
    if "devices_model" in db.list_tables():
        rows = db.query(
            "SELECT uuid FROM {q}devices_model{/q} WHERE project_uuid=? AND deleted_at IS NULL",
            (project_uuid,),
        )
        muids.update(r["uuid"] for r in rows)

    device_uuids: set[str] = set()
    rows = db.query(
        "SELECT uuid FROM {q}monitor_list{/q} WHERE project_uuid=? AND deleted_at IS NULL",
        (project_uuid,),
    )
    device_uuids.update(r["uuid"] for r in rows)

    return {
        "project_uuid": project_uuid,
        "display_uuid": display_uuid,
        "muids": sorted(muids),
        "device_uuids": sorted(device_uuids),
    }


def _classify(table: str, columns: list[str]) -> tuple[str, str] | None:
    """返回 (kind, category)；与项目无关的表返回 None。"""
    cset = set(columns)
    if table == "project_lists":
        return ("uuid_eq_project", "config")
    if table == "display_model_layer":
        return ("model_id", "config")
    if table == "devices_model":
        return ("devices_model", "config")
    if table == "device_real_data":
        return ("device_real", "config")
    if "project_uuid" in cset:
        if any(p.search(table) for p in HISTORY_PATTERNS):
            return ("project_uuid", "history")
        if table in LOG_TABLES:
            return ("project_uuid", "log")
        return ("project_uuid", "config")
    if "muid" in cset:
        return ("muid", "config")
    return None


def plan_tables(db: Db, scope: dict) -> list[TablePlan]:
    plans: list[TablePlan] = []
    for table in db.list_tables():
        cols = db.columns(table)
        verdict = _classify(table, cols)
        if verdict is None:
            continue
        kind, category = verdict
        plans.append(TablePlan(table, cols, kind, category))
    return plans


def _in_clause(col: str, values: list, db: Db) -> tuple[str, list]:
    if not values:
        return ("1=0", [])
    ph = ", ".join("?" for _ in values)
    return (f"{{q}}{col}{{/q}} IN ({ph})", list(values))


def build_where(plan: TablePlan, scope: dict, db: Db) -> tuple[str, list]:
    """根据 scope 为某表生成 WHERE 子句 (中性占位符) + 参数。"""
    P = scope["project_uuid"]
    M = scope["display_uuid"]
    if plan.kind == "uuid_eq_project":
        return ("{q}uuid{/q}=?", [P])
    if plan.kind == "project_uuid":
        return ("{q}project_uuid{/q}=?", [P])
    if plan.kind == "model_id":
        return ("{q}model_id{/q}=?", [M])
    if plan.kind == "muid":
        return _in_clause("muid", scope["muids"], db)
    if plan.kind == "devices_model":
        clause, params = _in_clause("uuid", scope["muids"], db)
        return (f"{{q}}project_uuid{{/q}}=? OR {clause}", [P, *params])
    if plan.kind == "device_real":
        clause, params = _in_clause("device_uuid", scope["device_uuids"], db)
        return (f"{{q}}project_uuid{{/q}}=? OR {clause}", [P, *params])
    raise ValueError(f"未知 scope kind: {plan.kind}")


# ───────────────────────── 行序列化 ─────────────────────────


def encode_value(value):
    if value is None:
        return None
    if isinstance(value, datetime):
        return {"__dt__": value.isoformat(sep=" ")}
    if isinstance(value, date):
        return {"__date__": value.isoformat()}
    if isinstance(value, Decimal):
        return float(value)
    if isinstance(value, (bytes, bytearray)):
        return {"__b64__": base64.b64encode(bytes(value)).decode()}
    return value


def decode_value(value):
    if isinstance(value, dict):
        if "__dt__" in value:
            return value["__dt__"]
        if "__date__" in value:
            return value["__date__"]
        if "__b64__" in value:
            return base64.b64decode(value["__b64__"])
    return value


def encode_row(row: dict, columns: list[str]) -> dict:
    return {c: encode_value(row.get(c)) for c in columns}
