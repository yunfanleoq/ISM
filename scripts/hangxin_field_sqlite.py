#!/usr/bin/env python3
"""把航信现场 MySQL 备份转成新的 SQLite 文件，并写出大屏用的名称索引。不覆盖 ism.db。"""
import json
import os
import re
import sqlite3
import sys
import zipfile

ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
ZIP_PATH = os.environ.get(
    "HANGXIN_SQL_ZIP",
    "/Users/yunfanleo/Library/Containers/com.tencent.xinWeChat/Data/Documents/xwechat_files/wxid_h1l5gtwyy3eu11_eb51/msg/file/2026-10/制作画面文件.zip",
)
OUT_DB = os.path.join(ROOT, "ism_server_user", "data", "db", "ism-hangxin-20261023.db")
OUT_INDEX = os.path.join(ROOT, "ism-front-end-v2", "public", "hx-dcim", "data", "hx-field-index.json")


def convert_create(sql):
    sql = re.sub(r"\)\s*ENGINE\s*=\s*[^;]*", ")", sql, flags=re.I)
    sql = sql.replace("`", '"')
    sql = re.sub(r"\bAUTO_INCREMENT\b", "", sql, flags=re.I)
    sql = re.sub(r"\bUNSIGNED\b", "", sql, flags=re.I)
    sql = re.sub(r"\b(?:BIGINT|INT|TINYINT|SMALLINT|MEDIUMINT)\(\d+\)", "INTEGER", sql, flags=re.I)
    sql = re.sub(r"\b(?:DOUBLE|FLOAT)(?:\(\d+\s*,\s*\d+\))?", "REAL", sql, flags=re.I)
    sql = re.sub(r"\b(?:DATETIME|TIMESTAMP)(?:\(\d+\))?", "TEXT", sql, flags=re.I)
    sql = re.sub(r"\b(?:LONGTEXT|MEDIUMTEXT|TINYTEXT|VARCHAR\(\d+\)|CHAR\(\d+\))", "TEXT", sql, flags=re.I)
    sql = re.sub(r"\bDEFAULT\s+CURRENT_TIMESTAMP(?:\(\d+\))?", "DEFAULT ''", sql, flags=re.I)
    return sql


def convert_sql(sql):
    head = sql.lstrip()[:12].upper()
    if head.startswith("CREATE"):
        return convert_create(sql)
    if head.startswith("INSERT") or head.startswith("DROP"):
        return sql.replace("`", '"').replace("\\'", "''")
    return None


def iter_statements(stream):
    buf = []
    for raw in stream:
        line = raw.decode("utf-8", "replace") if isinstance(raw, bytes) else raw
        if not line.strip() or line.startswith("/*") or line.startswith("--"):
            continue
        buf.append(line)
        if line.rstrip().endswith(";"):
            yield "".join(buf)
            buf = []
    if buf:
        yield "".join(buf)


def main():
    if not OUT_DB.endswith("ism-hangxin-20261023.db"):
        sys.exit("refuse to write unexpected db path")
    if os.path.basename(OUT_DB) == "ism.db":
        sys.exit("refuse to overwrite ism.db")
    os.makedirs(os.path.dirname(OUT_DB), exist_ok=True)
    os.makedirs(os.path.dirname(OUT_INDEX), exist_ok=True)
    if os.path.exists(OUT_DB):
        os.remove(OUT_DB)
    conn = sqlite3.connect(OUT_DB)
    conn.execute("PRAGMA journal_mode=OFF")
    conn.execute("PRAGMA synchronous=OFF")
    conn.execute("PRAGMA foreign_keys=OFF")
    ok = 0
    bad = 0
    with zipfile.ZipFile(ZIP_PATH) as zf:
        name = [n for n in zf.namelist() if n.endswith(".sql")][0]
        with zf.open(name) as raw:
            for stmt in iter_statements(raw):
                sql = convert_sql(stmt)
                if not sql:
                    continue
                try:
                    conn.execute(sql)
                    ok += 1
                    if ok % 200 == 0:
                        conn.commit()
                        print("ok", ok, "bad", bad, flush=True)
                except Exception as exc:
                    bad += 1
                    if bad <= 20:
                        print("skip", type(exc).__name__, str(exc)[:180], flush=True)
    conn.commit()
    print("import done", ok, bad, flush=True)
    nodes = []
    try:
        for sid, pid, name, uuid in conn.execute(
            "SELECT sid, pid, name, uuid FROM monitor_list WHERE deleted_at IS NULL OR deleted_at = ''"
        ):
            if name:
                nodes.append({"sid": sid, "pid": pid, "name": name, "uuid": uuid or ""})
    except sqlite3.Error as exc:
        print("index nodes", exc, flush=True)
    by_device = {}
    try:
        cur = conn.execute(
            """SELECT name, uuid, device_uuid, device_name, data_unit, value
               FROM device_real_data
               WHERE (deleted_at IS NULL OR deleted_at = '')
                 AND (name LIKE '%配电%' OR device_name LIKE '%配电%'
                      OR name LIKE '%数据机房%' OR device_name LIKE '%数据机房%'
                      OR name LIKE '%列头%' OR device_name LIKE '%列头%')"""
        )
        for name, uuid, duid, dname, unit, value in cur:
            if not name or "电容补偿" in name or (dname and "电容补偿" in dname):
                continue
            key = duid or dname or ""
            by_device.setdefault(key, []).append({
                "name": name,
                "uuid": uuid or "",
                "device": dname or "",
                "unit": unit or "",
                "value": value or "",
            })
    except sqlite3.Error as exc:
        print("index points", exc, flush=True)
    with open(OUT_INDEX, "w", encoding="utf-8") as f:
        json.dump({"nodes": nodes, "byDevice": by_device}, f, ensure_ascii=False)
    print("nodes", len(nodes), "devices", len(by_device), "index", OUT_INDEX, flush=True)
    conn.close()


if __name__ == "__main__":
    main()
