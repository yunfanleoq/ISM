#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""中航信数据中心电力监控系统 —— 通用数据驱动导入器。

读取 hx_parse_and_package.py 产出的两份契约:
  <set>/ism_data_models.json               (devices[])
  <set>/中航信_complete_project_package.json (deviceModels/registerGroups/registerPoints)

按 楼栋→室→柜→设备 自动建监控树、按包定义建模型/寄存器组/数据点、
为每台设备写真实管理机 IP + 全局唯一 slave id 的 extra_data。

用法:
  python3 scripts/import_hx_dc.py --set validation [--new]
  python3 scripts/import_hx_dc.py --set full --new
"""
import json, subprocess, time, os, sys, argparse, socket, sqlite3
from collections import defaultdict

try:
    import pymysql
except ImportError:
    pymysql = None

BASE_URL = "http://localhost:8081"
ROOT = os.path.join(os.path.dirname(__file__), "..")
DB_PATH = os.path.join(ROOT, "ism_server_user", "data", "db", "ism.db")
PROJECT_NAME = "中航信数据中心电力监控系统"
OB_CONFIG = {"host": "127.0.0.1", "port": 2881, "user": "root@ism_tenant",
             "password": "ism2024!", "database": "ism"}
SUCCESS_CODES = {0, 200, 2002, 4002, 3001}

PROJECT_UUID = ""
HEADERS = {}


def _ob_available():
    if pymysql is None:
        return False
    s = socket.socket()
    try:
        s.settimeout(0.5)
        s.connect(("127.0.0.1", OB_CONFIG["port"]))
        return True
    except OSError:
        return False
    finally:
        s.close()


USE_SQLITE = os.environ.get("ISM_FORCE_OB") != "1" and not _ob_available()


def ob_read(sql, params=()):
    if USE_SQLITE:
        q = sql.replace("%s", "?")
        conn = sqlite3.connect(DB_PATH)
        cur = conn.cursor()
        cur.execute(q, params)
        rows = cur.fetchall()
        cur.close()
        conn.close()
        return rows
    conn = pymysql.connect(**OB_CONFIG)
    cur = conn.cursor()
    cur.execute(sql, params)
    rows = cur.fetchall()
    cur.close()
    conn.close()
    return rows


def api(url, json_data=None, headers=None):
    try:
        h = headers or HEADERS
        args = ["-s", "-X", "POST", f"{BASE_URL}{url}"]
        for k, v in h.items():
            args.extend(["-H", f"{k}:{v}"])
        if json_data is not None:
            args.extend(["-d", json.dumps(json_data)])
        result = subprocess.run(["curl"] + args, capture_output=True, text=True, timeout=30)
        j = json.loads(result.stdout)
        return j.get("code", -1), j.get("data", {})
    except Exception as e:
        return -1, str(e)


def api_proj(url, json_data=None):
    h = dict(HEADERS); h["ProjectUuid"] = PROJECT_UUID
    return api(url, json_data, headers=h)


def login(username="admin", password="123456"):
    payload = json.dumps({"username": username, "password": password})
    r = subprocess.run(["curl", "-s", "-X", "POST", f"{BASE_URL}/login",
                        "-H", "Content-Type:application/json", "-d", payload],
                       capture_output=True, text=True, timeout=10)
    data = json.loads(r.stdout)
    token = data["data"]["token"]
    row = ob_read("SELECT uuid,name FROM user WHERE username=%s AND deleted_at IS NULL", (username,))
    uid = row[0][0] if row else ""
    print(f"  登录: {row[0][1] if row else '?'} ({uid[:16]}...)")
    return token, uid


def main():
    global PROJECT_UUID, HEADERS, PROJECT_NAME
    ap = argparse.ArgumentParser()
    ap.add_argument("--set", choices=["validation", "full"], default="validation")
    ap.add_argument("--new", action="store_true", help="同名项目存在则加时间戳新建")
    ap.add_argument("--gw-override", default="", help="测试用: 覆盖所有设备网关为 IP:PORT(如 127.0.0.1:1502)")
    ap.add_argument("--local-sim", action="store_true",
                    help="无 sudo 验证: 每网关映射到 127.0.0.1:(1502+序号)，配合 hx_simulator.py --local")
    args = ap.parse_args()
    gw_ip, gw_port = ("", "")
    if args.gw_override:
        gw_ip, gw_port = args.gw_override.split(":")

    base = os.path.join(ROOT, "hx-data", args.set)
    dm_data = json.load(open(os.path.join(base, "ism_data_models.json")))
    pkg = json.load(open(os.path.join(base, "中航信_complete_project_package.json")))
    devices = dm_data["devices"]
    # 网关 → 本机端口 的确定性映射(与 hx_simulator.py --local 完全一致)
    local_port_map = {}
    if args.local_sim:
        for i, ip in enumerate(sorted({d["gatewayIP"] for d in devices})):
            local_port_map[ip] = 1502 + i

    print("=" * 60)
    print(f"  中航信数据中心导入  set={args.set}  设备={len(devices)}")
    print("=" * 60)

    # [0] 环境
    code = subprocess.run(["curl", "-s", "-o", "/dev/null", "-w", "%{http_code}", f"{BASE_URL}/"],
                          capture_output=True, text=True, timeout=5).stdout.strip()
    assert code in ("200", "404"), f"后端不可达 {code}"
    if USE_SQLITE:
        assert os.path.isfile(DB_PATH), f"SQLite 不存在: {DB_PATH}"
        sqlite3.connect(DB_PATH).close()
        print("[0] 后端 + SQLite OK")
    else:
        pymysql.connect(**OB_CONFIG).close()
        print("[0] 后端 + OceanBase OK")

    # [1] 登录
    print("[1] 登录")
    token, user_uuid = login()
    HEADERS = {"Content-Type": "application/json", "Authorization": token}

    # [2] 创建项目
    print(f"[2] 项目: {PROJECT_NAME}")
    existing = ob_read("SELECT uuid FROM project_lists WHERE name=%s AND deleted_at IS NULL", (PROJECT_NAME,))
    if existing and args.new:
        PROJECT_NAME = f"{PROJECT_NAME}-{time.strftime('%m%d%H%M')}"
        existing = []
    if existing:
        PROJECT_UUID = existing[0][0]
        print(f"  复用已有项目 {PROJECT_UUID}")
    else:
        api("/ProjectAdd", {"name": PROJECT_NAME, "description": "中航信数据中心电力监控系统", "industry": 1})
        PROJECT_UUID = ob_read("SELECT uuid FROM project_lists WHERE name=%s AND deleted_at IS NULL ORDER BY id DESC LIMIT 1",
                               (PROJECT_NAME,))[0][0]
        print(f"  新建项目 {PROJECT_UUID}")
    HEADERS["ProjectUuid"] = PROJECT_UUID
    api_proj("/ProjectFixCreator", {})

    # [3] 清理跨项目同名设备冲突
    names = [d["name"] for d in devices]
    if names:
        ph = ",".join(["%s"] * len(names))
        conflicts = ob_read(f"SELECT uuid,project_uuid FROM monitor_list WHERE name IN ({ph}) AND deleted_at IS NULL", tuple(names))
        byproj = defaultdict(list)
        for u, p in conflicts:
            if p != PROJECT_UUID:
                byproj[p].append(u)
        for p, us in byproj.items():
            h = dict(HEADERS); h["ProjectUuid"] = p
            api("/monitorAllDel", {"uuid": us}, headers=h)
            print(f"  清理跨项目冲突 {len(us)} 台 (proj={p[:12]})")

    # [4] RootZone — 项目创建时 ProjectModelAdd 已自动生成 sid=1 的根，勿重复插入
    existing_rz = ob_read(
        "SELECT uuid, sid FROM monitor_list WHERE name='RootZone' AND project_uuid=%s AND pid=0 AND deleted_at IS NULL",
        (PROJECT_UUID,))
    sid1 = [r for r in existing_rz if r[1] == 1]
    if sid1:
        api_proj("/monitorEdit", {"data": {"Sid": 1, "uuid": sid1[0][0]}})
        print(f"[4] RootZone 已存在 sid=1，跳过创建 ({len(existing_rz)} 条根记录)")
    else:
        api("/monitorAdd", {"sid": 1, "pid": 0, "name": "RootZone", "type": 0, "timeout": 5,
                            "IsEnable": 1, "project_uuid": PROJECT_UUID, "interval": 5, "failedTimes": 5,
                            "description": "根区域", "offlineClear": 0, "offlineDefaultValue": "0",
                            "deviceType": 0, "muid": "", "configUid": "", "PageUUID": "", "extra": "",
                            "Status": 0, "longitude": "", "latitude": ""})
        rz = ob_read("SELECT uuid FROM monitor_list WHERE name='RootZone' AND project_uuid=%s AND deleted_at IS NULL LIMIT 1", (PROJECT_UUID,))
        if rz:
            api_proj("/monitorEdit", {"data": {"Sid": 1, "uuid": rz[0][0]}})
        print("[4] RootZone OK")

    # [5] 模型 + [6] 寄存器组 + [7] 数据点（按包定义）
    print("[5] 创建模型")
    pkg2db_muid = {}
    for dm in pkg["deviceModels"]:
        api("/modbusModelAdd", {
            "name": dm["name"], "dec": dm.get("dec", dm["name"]), "type": 2,
            "gatherNumber": dm.get("gatherNumber", 30), "port": 502,
            "timeout": 5, "DataFormat": dm.get("DataFormat", "ABCD"), "modbusConnectType": "TCPClient",
            "modbusConnectMode": "TCP/IP", "modbusClientIpaddress": "127.0.0.1",
            "configUid": "", "PageUUID": "", "version": 1, "modbusCom": "",
            "serialBaud": 9600, "serialBits": 8, "serialParity": "None",
            "serialStopBits": "1", "serialFlow": "None"})
        row = ob_read("SELECT uuid FROM devices_model WHERE name=%s AND project_uuid=%s AND deleted_at IS NULL ORDER BY id DESC LIMIT 1",
                      (dm["name"], PROJECT_UUID))
        if not row:
            print(f"  !! 模型创建失败 {dm['name']}"); sys.exit(1)
        pkg2db_muid[dm["uuid"]] = row[0][0]
    print(f"  {len(pkg2db_muid)} 个模型")

    print("[6] 创建寄存器组")
    pkg_rg2db = {}
    for rg in pkg["registerGroups"]:
        db_mu = pkg2db_muid[rg["muid"]]
        api("/modbusModelRegisterGroupAdd", {"name": rg["name"], "muid": db_mu,
            "function": rg["function"], "registerStart": rg["registerStart"],
            "registerCount": rg["registerCount"]})
        row = ob_read("SELECT uuid FROM modbus_devices_register_group WHERE name=%s AND muid=%s AND deleted_at IS NULL ORDER BY id DESC LIMIT 1",
                      (rg["name"], db_mu))
        if row:
            pkg_rg2db[rg["uuid"]] = row[0][0]
    print(f"  {len(pkg_rg2db)} 个寄存器组")

    print("[7] 创建数据点")
    n_pts = 0; total = len(pkg["registerPoints"])
    for i, pt in enumerate(pkg["registerPoints"], 1):
        db_mu = pkg2db_muid.get(pt["muid"]); db_rg = pkg_rg2db.get(pt["registerGroupUuid"])
        if not db_mu or not db_rg:
            continue
        code, _ = api("/modbusModelRegisterAdd", {
            "name": pt["name"], "muid": db_mu, "registerAddress": pt["registerAddress"],
            "registerGroupUuid": db_rg, "auth": pt.get("auth", "ReadOnly"),
            "type": pt.get("type", "Float"), "ByteOrder": pt.get("ByteOrder", "CDAB"),
            "modeltype": pt.get("modeltype", 2), "unit": pt.get("unit", ""),
            "conversionExpression": pt.get("conversionExpression", ""),
            "alarm": 0, "alarmLevel": 0, "AlarmMessage": "", "AlarmClearMessage": "",
            "record": pt.get("record", 0), "RecordType": pt.get("RecordType", 1),
            "RecordInterval": pt.get("RecordInterval", 5), "RecordDataCharge": "0",
            "RecordDataTimely": "0", "FloatAccuracy": pt.get("FloatAccuracy", "0.01")})
        if code in SUCCESS_CODES:
            n_pts += 1
        if i % 200 == 0:
            print(f"    {i}/{total} ...")
    print(f"  数据点 {n_pts}/{total}")
    # 清 registerN 默认命名
    for pkg_mu, db_mu in pkg2db_muid.items():
        pts = ob_read("SELECT uuid,name FROM modbus_devices_data_model WHERE muid=%s AND deleted_at IS NULL", (db_mu,))
        defaults = [p[0] for p in pts if str(p[1]).lower().startswith("register")]
        if defaults:
            api_proj("/modbusModelRegisterDel", {"uuid": defaults})

    # [8] 监控树（楼栋→室→柜）+ 设备
    print("[8] 监控树 + 设备")
    root_sid = ob_read("SELECT sid FROM monitor_list WHERE name='RootZone' AND project_uuid=%s AND deleted_at IS NULL LIMIT 1", (PROJECT_UUID,))[0][0]
    zone_cache = {}

    def add_zone(pid, name):
        key = (pid, name)
        if key in zone_cache:
            return zone_cache[key]
        api("/monitorAdd", {"sid": 0, "pid": pid, "name": name, "type": 0, "timeout": 5,
                            "IsEnable": 1, "project_uuid": PROJECT_UUID, "interval": 5,
                            "failedTimes": 5, "description": name, "offlineClear": 0,
                            "offlineDefaultValue": "0", "deviceType": 0, "muid": "",
                            "configUid": "", "PageUUID": "", "extra": "", "Status": 0,
                            "longitude": "", "latitude": ""})
        row = ob_read("SELECT sid FROM monitor_list WHERE name=%s AND pid=%s AND type=0 AND project_uuid=%s AND deleted_at IS NULL ORDER BY id DESC LIMIT 1",
                      (name, pid, PROJECT_UUID))
        sid = row[0][0] if row else None
        zone_cache[key] = sid
        return sid

    # 模型名 → db muid
    name2muid = {dm["name"]: pkg2db_muid[dm["uuid"]] for dm in pkg["deviceModels"]}

    n_dev = 0
    for d in devices:
        b = d.get("building") or "未分组"
        r = d.get("room") or b
        c = d.get("cabinet") or r
        b_sid = add_zone(root_sid, b)
        r_sid = add_zone(b_sid, r) if r != b else b_sid
        c_sid = add_zone(r_sid, c) if c != r else r_sid
        muid = name2muid.get(d["modelName"])
        if not muid:
            continue
        if args.local_sim:
            dev_ip, dev_port = "127.0.0.1", str(local_port_map[d["gatewayIP"]])
        else:
            dev_ip = gw_ip or d.get("gatewayIP", "127.0.0.1")
            dev_port = gw_port or str(d.get("port", 502))
        extra = json.dumps({"modbus": {"IPAddress": dev_ip, "Port": dev_port,
                                       "address": str(d["slaveId"]),
                                       "RegisterPack": -1, "packTime": 100}})
        code, _ = api("/monitorAdd", {"sid": 0, "pid": c_sid, "name": d["name"], "type": 1,
                                      "timeout": 3000, "IsEnable": 1, "project_uuid": PROJECT_UUID,
                                      "interval": 500, "failedTimes": 5,
                                      "description": d.get("display", d["name"]),
                                      "offlineClear": 0, "offlineDefaultValue": "0",
                                      "deviceType": 2, "muid": muid, "configUid": "",
                                      "PageUUID": "", "extra": extra, "Status": 0,
                                      "longitude": "", "latitude": ""})
        if code in SUCCESS_CODES:
            n_dev += 1
    print(f"  区域 {len(zone_cache)} 个, 设备 {n_dev}/{len(devices)}")

    # [9] 实时数据 + 禁用告警
    print("[9] syncDeviceRealData + 禁用告警")
    api_proj("/syncDeviceRealData", {})
    api_proj("/DeviceRealDataDisableAlarm", {})
    api_proj("/MonitorBatchSetStatus", {"status": 1})

    # [10] 验证
    models = ob_read("SELECT uuid FROM devices_model WHERE project_uuid=%s AND deleted_at IS NULL", (PROJECT_UUID,))
    muids = [m[0] for m in models]
    ml = ob_read("SELECT count(*) FROM monitor_list WHERE project_uuid=%s AND deleted_at IS NULL", (PROJECT_UUID,))[0][0]
    drd = ob_read("SELECT count(*) FROM device_real_data WHERE project_uuid=%s", (PROJECT_UUID,))[0][0]
    devc = ob_read("SELECT count(DISTINCT device_uuid) FROM device_real_data WHERE project_uuid=%s", (PROJECT_UUID,))[0][0]
    dp = 0
    if muids:
        ph = ",".join(["%s"] * len(muids))
        dp = ob_read(f"SELECT count(*) FROM modbus_devices_data_model WHERE muid IN ({ph}) AND deleted_at IS NULL", muids)[0][0]
    print(f"""
{'='*60}
  项目: {PROJECT_NAME}
  UUID: {PROJECT_UUID}
  模型: {len(models)}   数据点: {dp}
  monitor_list: {ml}   device_real_data: {drd} ({devc} 设备)
{'='*60}""")
    print("\n✅ 导入完成")
    print(f"PROJECT_UUID={PROJECT_UUID}")


if __name__ == "__main__":
    main()
