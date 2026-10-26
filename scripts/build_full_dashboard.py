#!/usr/bin/env python3
"""为全量项目创建 display model 并构建大屏。
端口耗尽(本机多网关模拟轮询导致 ephemeral 端口紧张)时自动重试。
"""
import json
import os
import subprocess
import sys
import time

import pymysql, socket, sqlite3

BASE_URL = "http://localhost:8081"
APPRUN_BASE = os.environ.get("NCC_APPRUN_BASE", "http://localhost:7080")


def apprun_url(model_id, page_id=None, base=APPRUN_BASE):
    url = f"{base}/#/AppRun/{model_id}"
    if page_id and page_id != model_id:
        url += f"?pageId={page_id}"
    return url
PROJECT_UUID = "168430f2-ba63-1a80-afcf-10d4061f0072"
MODEL_NAME = "中航信数据中心电力监控系统"
OB = dict(host="127.0.0.1", port=2881, user="root@ism_tenant",
          password="ism2024!", database="ism")


DB_PATH = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), "ism_server_user", "data", "db", "ism.db")

def _ob_up():
    s = socket.socket()
    try:
        s.settimeout(0.5)
        s.connect(("127.0.0.1", OB["port"]))
        return True
    except OSError:
        return False
    finally:
        s.close()

def q(sql):
    return sql.replace("%s", "?") if USE_SQLITE else sql

USE_SQLITE = os.environ.get("ISM_FORCE_OB") != "1" and not _ob_up()

def db():
    if USE_SQLITE:
        return sqlite3.connect(DB_PATH)
    last = None
    for _ in range(60):
        try:
            return pymysql.connect(**OB)
        except Exception as e:
            last = e
            time.sleep(0.5)
    raise SystemExit(f"DB 连接失败: {last}")


def curl(args, timeout=30):
    for _ in range(10):
        r = subprocess.run(["curl"] + args, capture_output=True, text=True, timeout=timeout)
        if r.returncode == 0 and r.stdout.strip():
            return r.stdout
        time.sleep(0.5)
    return r.stdout


def login():
    out = curl(["-s", "-X", "POST", f"{BASE_URL}/login",
                "-H", "Content-Type:application/json",
                "-d", json.dumps({"username": "admin", "password": "123456"})])
    return json.loads(out)["data"]["token"]


def main():
    token = login()
    print(f"[1] 登录 OK token={token[:12]}...")

    # 已存在则复用
    c = db(); cur = c.cursor()
    cur.execute(q("SELECT display_model_uid FROM display_models WHERE project_uuid=%s AND deleted_at IS NULL ORDER BY id DESC LIMIT 1"), (PROJECT_UUID,))
    row = cur.fetchone()
    c.close()

    if row:
        model_uid = row[0]
        print(f"[2] 复用已有 display model: {model_uid}")
    else:
        out = curl(["-s", "-X", "POST", f"{BASE_URL}/displayModelAdd",
                    "-H", "Content-Type:application/json",
                    "-H", f"Authorization:{token}",
                    "-H", f"ProjectUuid:{PROJECT_UUID}",
                    "-d", json.dumps({"name": MODEL_NAME,
                                      "description": "航信机房全量电力监控大屏",
                                      "display_type": 1})])
        print(f"[2] displayModelAdd -> {out.strip()}")
        time.sleep(1)
        c = db(); cur = c.cursor()
        cur.execute(q("SELECT display_model_uid FROM display_models WHERE project_uuid=%s AND deleted_at IS NULL ORDER BY id DESC LIMIT 1"), (PROJECT_UUID,))
        row = cur.fetchone(); c.close()
        if not row:
            raise SystemExit("创建 display model 后查不到记录")
        model_uid = row[0]
        print(f"    新建 display model: {model_uid}")

    print(f"[3] 构建大屏 MODEL_ID={model_uid} PROJECT={PROJECT_UUID}")
    env = dict(os.environ)
    env["NCC_MODEL_ID"] = model_uid
    env["NCC_PROJECT_UUID"] = PROJECT_UUID
    r = subprocess.run([sys.executable, "build_ncc_dashboard.py"],
                       cwd=os.path.dirname(os.path.dirname(os.path.abspath(__file__))),
                       env=env)
    if r.returncode != 0:
        raise SystemExit(f"build_ncc_dashboard.py 失败 rc={r.returncode}")

    print("=" * 50)
    print(f"✅ 大屏构建完成")
    print(f"MODEL_ID={model_uid}")
    print(f"AppRun: {apprun_url(model_uid)}")


if __name__ == "__main__":
    main()
