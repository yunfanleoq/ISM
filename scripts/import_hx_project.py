#!/usr/bin/env python3
"""
ISM 航信机房项目导入脚本（OceanBase 适配版）
基于 import_1a_project.py v3.0，适配 OceanBase 数据库

用法:
  python3 scripts/import_hx_project.py
"""

import json, subprocess, time, os, sys
from collections import defaultdict
import pymysql

# ========== 配置 ==========
BASE_URL = "http://localhost:8081"
PROJECT_ROOT = os.path.join(os.path.dirname(__file__), "..")
PKG_PATH = os.path.join(PROJECT_ROOT, "liu-chang-1A-dev", "1A_complete_project_package.json")
MODELS_PATH = os.path.join(PROJECT_ROOT, "liu-chang-1A-dev", "ism_data_models.json")
PROJECT_NAME = "航信机房"

# OceanBase 连接
OB_CONFIG = {
    "host": "127.0.0.1", "port": 2881,
    "user": "root@ism_tenant", "password": "ism2024!",
    "database": "ism"
}

SUCCESS_CODES = {0, 200, 2002, 4002, 3001}

def ob_read(sql, params=()):
    """只读查询 OceanBase，返回 fetchall() 结果列表"""
    conn = pymysql.connect(**OB_CONFIG)
    cur = conn.cursor()
    cur.execute(sql, params)
    rows = cur.fetchall()
    cur.close()
    conn.close()
    return rows

def api(url, json_data=None, headers=None):
    """统一 API 调用（curl 子进程）"""
    try:
        h = headers or HEADERS
        args = ["-s", "-X", "POST", f"{BASE_URL}{url}"]
        for k, v in h.items():
            args.extend(["-H", f"{k}:{v}"])
        if json_data is not None:
            body = json.dumps(json_data)
            args.extend(["-d", body])
        result = subprocess.run(["curl"] + args, capture_output=True, text=True, timeout=15)
        j = json.loads(result.stdout)
        return j.get('code', -1), j.get('data', {})
    except Exception as e:
        return -1, str(e)

def api_proj(url, json_data=None):
    h = dict(HEADERS)
    h["ProjectUuid"] = PROJECT_UUID
    return api(url, json_data, headers=h)

def login(username="admin", password="123456"):
    """登录并返回 token 和用户 UUID"""
    payload = json.dumps({"username": username, "password": password})
    result = subprocess.run(
        ["curl", "-s", "-X", "POST", f"{BASE_URL}/login",
         "-H", "Content-Type:application/json", "-d", payload],
        capture_output=True, text=True, timeout=10
    )
    data = json.loads(result.stdout)
    token = data['data']['token']
    # 从 OceanBase user 表查 UUID
    conn = pymysql.connect(**OB_CONFIG)
    cur = conn.cursor()
    cur.execute("SELECT uuid, name FROM user WHERE username=%s AND deleted_at IS NULL", (username,))
    row = cur.fetchone()
    cur.close()
    conn.close()
    if row:
        user_uuid, user_name = row[0], row[1]
    else:
        user_uuid, user_name = "", "unknown"
    print(f"  登录成功: {user_name} (uuid={user_uuid[:24]}...)")
    return token, user_uuid, user_name

def add_zone(pid, name, dt=0):
    api("/monitorAdd", {
        "sid": 0, "pid": pid, "name": name, "type": dt,
        "timeout": 5, "IsEnable": 1, "project_uuid": PROJECT_UUID,
        "interval": 5, "failedTimes": 5, "description": name,
        "offlineClear": 0, "offlineDefaultValue": "0", "deviceType": 0,
        "muid": "", "configUid": "", "PageUUID": "", "extra": "",
        "Status": 0, "longitude": "", "latitude": ""
    })
    return get_sid(name, dt)

def get_sid(name, dt):
    r = ob_read(
        "SELECT sid FROM monitor_list WHERE name=%s AND type=%s AND project_uuid=%s AND deleted_at IS NULL",
        (name, dt, PROJECT_UUID)
    )
    return r[0][0] if r else None

# ========== Step 0: 环境检查 ==========
print("=" * 60)
print("  ISM 航信机房项目导入（OceanBase 适配版）")
print("=" * 60)

print("\n[0] 环境检查")
try:
    result = subprocess.run(["curl", "-s", "-o", "/dev/null", "-w", "%{http_code}", f"{BASE_URL}/"],
        capture_output=True, text=True, timeout=3)
    if result.stdout.strip() in ("200", "404"):
        print(f"  后端就绪: {BASE_URL}")
    else:
        print("  后端异常!"); sys.exit(1)
except:
    print("  后端不可达!"); sys.exit(1)

# 验证 OceanBase 连接
try:
    conn = pymysql.connect(**OB_CONFIG)
    conn.close()
    print("  OceanBase 连接: OK")
except Exception as e:
    print(f"  OceanBase 连接失败: {e}"); sys.exit(1)

# ========== Step 1: 登录 ==========
print("\n[1] 登录获取用户 UUID")
TOKEN, USER_UUID, USER_NAME = login()
HEADERS = {"Content-Type": "application/json", "Authorization": TOKEN}

# ========== Step 2: 创建项目 ==========
print(f"\n[2] 创建项目: {PROJECT_NAME}")

existing = ob_read(
    "SELECT uuid FROM project_lists WHERE name=%s AND deleted_at IS NULL", (PROJECT_NAME,)
)
if existing:
    print(f"  已有同名项目，加时间戳")
    PROJECT_NAME = f"{PROJECT_NAME}-{time.strftime('%m%d%H%M')}"

code, data = api("/ProjectAdd", {
    "name": PROJECT_NAME, "description": "航信机房 Modbus 监控项目",
    "industry": 1
})
PROJECT_UUID = ob_read(
    "SELECT uuid FROM project_lists WHERE name=%s AND deleted_at IS NULL ORDER BY id DESC LIMIT 1",
    (PROJECT_NAME,)
)[0][0]
print(f"  新建项目: {PROJECT_UUID}")

HEADERS["ProjectUuid"] = PROJECT_UUID
code, _ = api_proj("/ProjectFixCreator", {})
print(f"  修正 creator_uuid → {USER_UUID[:24]}... (code={code})")

# ========== Step 3: 清理跨项目同名设备冲突 ==========
print(f"\n[3] 清理跨项目同名设备冲突")
with open(MODELS_PATH) as f:
    dm_data = json.load(f)
all_dev_names = [d['name'] for d in dm_data.get('devices', [])]

if all_dev_names:
    placeholders = ','.join(['%s' for _ in all_dev_names])
    conflicts = ob_read(
        f"SELECT uuid, name, project_uuid FROM monitor_list WHERE name IN ({placeholders}) AND deleted_at IS NULL",
        tuple(all_dev_names)
    )
    if conflicts:
        proj_devs = defaultdict(list)
        for dev_uuid, dev_name, proj_uuid in conflicts:
            proj_devs[proj_uuid].append(dev_uuid)
        for proj_uuid, dev_uuids in proj_devs.items():
            old_puuid = HEADERS.get("ProjectUuid", "")
            HEADERS["ProjectUuid"] = proj_uuid
            code, _ = api("/monitorAllDel", {"uuid": dev_uuids}, headers=HEADERS)
            HEADERS["ProjectUuid"] = old_puuid
            print(f"  清理: {len(dev_uuids)} 个设备 (project={proj_uuid[:12]}...) code={code}")
    else:
        print("  无跨项目冲突")

HEADERS["ProjectUuid"] = PROJECT_UUID

# ========== Step 4: RootZone ==========
print("\n[4] RootZone")
code, _ = api("/monitorAdd", {
    "sid": 1, "pid": 0, "name": "RootZone", "type": 0,
    "timeout": 5, "IsEnable": 1, "project_uuid": PROJECT_UUID,
    "interval": 5, "failedTimes": 5, "description": "根区域",
    "offlineClear": 0, "offlineDefaultValue": "0", "deviceType": 0,
    "muid": "", "configUid": "", "PageUUID": "", "extra": "",
    "Status": 0, "longitude": "", "latitude": ""
})
rz = ob_read(
    "SELECT uuid FROM monitor_list WHERE name='RootZone' AND project_uuid=%s AND deleted_at IS NULL LIMIT 1",
    (PROJECT_UUID,)
)
if rz:
    api_proj("/monitorEdit", {"data": {"Sid": 1, "uuid": rz[0][0]}})
print(f"  RootZone code={code}")

# ========== Step 5: 数据模型 ==========
print("\n[5] 创建数据模型")
MODELS_DEF = {
    "A20电力仪表": {"dec": "A20电力仪表", "type": 2, "gatherNumber": 30, "port": 502},
    "A40电力仪表": {"dec": "A40电力仪表", "type": 2, "gatherNumber": 30, "port": 502},
    "施耐德UPS":    {"dec": "施耐德UPS",    "type": 2, "gatherNumber": 30, "port": 502},
}
model_uuids = {}
for mname, cfg in MODELS_DEF.items():
    code, _ = api("/modbusModelAdd", {
        **cfg, "name": mname, "timeout": 5, "DataFormat": "CDAB",
        "modbusConnectType": "TCPClient", "modbusConnectMode": "TCP/IP",
        "modbusClientIpaddress": "127.0.0.1",
        "configUid": "", "PageUUID": "", "version": 1,
        "modbusCom": "", "serialBaud": 9600, "serialBits": 8,
        "serialParity": "None", "serialStopBits": "1", "serialFlow": "None"
    })
    row = ob_read(
        "SELECT uuid FROM devices_model WHERE name=%s AND project_uuid=%s AND deleted_at IS NULL ORDER BY id DESC LIMIT 1",
        (mname, PROJECT_UUID)
    )
    if row:
        model_uuids[mname] = row[0][0]
        print(f"  {mname}: {row[0][0][:12]}... code={code}")
    else:
        print(f"  {mname}: FAILED code={code}")
        sys.exit(1)

# ========== Step 6: 寄存器组 ==========
print("\n[6] 创建寄存器组")
RG_DEFS = [
    ("AI数据", "A20电力仪表", 3, 0, 30),
    ("DI数据", "A20电力仪表", 2, 0, 3),
    ("AI数据", "A40电力仪表", 3, 0, 41),
    ("DI数据", "A40电力仪表", 2, 0, 3),
    ("AI数据", "施耐德UPS",    3, 0, 44),
    ("DI数据", "施耐德UPS",    2, 0, 1),
]
group_uuids = {}
for gname, mname, fn, start, cnt in RG_DEFS:
    muid = model_uuids.get(mname)
    if not muid: continue
    code, _ = api("/modbusModelRegisterGroupAdd", {
        "name": gname, "muid": muid, "function": fn,
        "registerStart": start, "registerCount": cnt
    })
    row = ob_read(
        "SELECT uuid FROM modbus_devices_register_group WHERE name=%s AND muid=%s AND deleted_at IS NULL ORDER BY id DESC LIMIT 1",
        (gname, muid)
    )
    k = f"{gname}_{muid}"
    if row: group_uuids[k] = row[0][0]
    print(f"  {gname}({mname}): {'OK' if row else 'FAIL'} code={code}")

# ========== Step 7: 数据点 ==========
print("\n[7] 创建 Modbus 数据点")
with open(PKG_PATH) as f:
    pkg = json.load(f)

pkg_model_names = [dm['name'] for dm in pkg['deviceModels']]
pkg_model_uuids = [dm['uuid'] for dm in pkg['deviceModels']]

NAME_MAP = {
    'A20': 'A20电力仪表', 'A40': 'A40电力仪表', '施耐德UPS': '施耐德UPS',
    'A20电力仪表': 'A20电力仪表', 'A40电力仪表': 'A40电力仪表',
}

pkg2db_muid = {}
for pkg_name, pkg_mu in zip(pkg_model_names, pkg_model_uuids):
    db_name = NAME_MAP.get(pkg_name)
    if not db_name:
        for pk, dv in NAME_MAP.items():
            if pk in pkg_name:
                db_name = dv; break
    if db_name and db_name in model_uuids:
        pkg2db_muid[pkg_mu] = model_uuids[db_name]

pkg_grp_map = {}
for rg in pkg.get('registerGroups', []):
    pkg_grp_map[(rg['muid'], rg['name'])] = rg['uuid']

pkg_rg2db = {}
for (pkg_mu, gn), pkg_rg_uuid in pkg_grp_map.items():
    db_mu = pkg2db_muid.get(pkg_mu)
    if not db_mu: continue
    k = f"{gn}_{db_mu}"
    if k in group_uuids:
        pkg_rg2db[pkg_rg_uuid] = group_uuids[k]
    else:
        row = ob_read(
            "SELECT uuid FROM modbus_devices_register_group WHERE name=%s AND muid=%s AND deleted_at IS NULL ORDER BY id DESC LIMIT 1",
            (gn, db_mu)
        )
        if row: pkg_rg2db[pkg_rg_uuid] = row[0][0]

n_pts = 0
for pt in pkg.get('registerPoints', []):
    pkg_mu = pt.get('muid', '')
    pkg_rg = pt.get('registerGroupUuid', '')
    db_mu = pkg2db_muid.get(pkg_mu)
    db_rg = pkg_rg2db.get(pkg_rg)
    if not db_mu or not db_rg: continue
    code, _ = api("/modbusModelRegisterAdd", {
        "name": pt['name'], "muid": db_mu,
        "registerAddress": pt['registerAddress'], "registerGroupUuid": db_rg,
        "auth": pt.get('auth', 'ReadOnly'), "type": pt.get('type', 'Float'),
        "ByteOrder": pt.get('ByteOrder', 'CDAB'), "modeltype": pt.get('modeltype', 2),
        "unit": pt.get('unit', ''),
        "conversionExpression": "",
        "alarm": 0, "alarmLevel": 0,
        "AlarmMessage": "", "AlarmClearMessage": "",
        "record": pt.get('record', 0), "RecordType": pt.get('RecordType', 1),
        "RecordInterval": pt.get('RecordInterval', 5),
        "RecordDataCharge": pt.get('RecordDataCharge', '0'),
        "RecordDataTimely": pt.get('RecordDataTimely', '0'),
        "FloatAccuracy": pt.get('FloatAccuracy', '0.01'),
    })
    if code in SUCCESS_CODES: n_pts += 1
print(f"  添加了 {n_pts}/{len(pkg.get('registerPoints',[]))} 个数据点")

# 清除 registerN 默认命名
for mname in model_uuids:
    mu = model_uuids[mname]
    all_pts = ob_read(
        "SELECT uuid, name FROM modbus_devices_data_model WHERE muid=%s AND deleted_at IS NULL", (mu,)
    )
    register_n = [pt[0] for pt in all_pts if pt[1].lower().startswith('register')]
    if register_n:
        api_proj("/modbusModelRegisterDel", {"uuid": register_n})
        print(f"  清除 {mname} 的 {len(register_n)} 个 registerN 默认命名")

# ========== Step 8: 设备树 + 设备实例 ==========
print("\n[8] 设备树和实例")

root = get_sid("RootZone", 0) or 1
site = add_zone(root, "航信机房")
floor = add_zone(site, "1楼配电室")
cab_a1 = add_zone(floor, "1A1_U11柜")
cab_a3 = add_zone(floor, "1A3_U11柜")
cab_ups = add_zone(floor, "UPS柜")
print(f"  站点:{site} 楼层:{floor} A1柜:{cab_a1} A3柜:{cab_a3} UPS柜:{cab_ups}")

TEMPLATE_MAP = {"A20": "A20电力仪表", "A40": "A40电力仪表", "施耐德UPS": "施耐德UPS"}

def get_cab_for_device(dev_name):
    # UPS设备统一放到UPS柜
    if 'UPS' in dev_name:
        return cab_ups
    if '1A3_' in dev_name or '1A3U' in dev_name:
        return cab_a3
    elif '1A1_' in dev_name or '1A1U' in dev_name:
        return cab_a1
    return cab_a1

slave_count = {"A20": 0, "A40": 0, "施耐德UPS": 0}
n_dev = 0

for dev in dm_data.get('devices', []):
    t = dev['templateType']
    muid = model_uuids[TEMPLATE_MAP[t]]
    pid = get_cab_for_device(dev['name'])
    slave_count[t] += 1

    if t == "A20":
        slave = slave_count[t]
    elif t == "A40":
        slave = 60 + slave_count[t]
    else:
        slave = 69 + slave_count[t]

    extra = json.dumps({
        "modbus": {
            "IPAddress": "127.0.0.1", "Port": "502",
            "address": str(slave), "RegisterPack": 0,
            "packTime": 100
        }
    })

    code, _ = api("/monitorAdd", {
        "sid": 0, "pid": pid, "name": dev['name'], "type": 1,
        "timeout": 5, "IsEnable": 1, "project_uuid": PROJECT_UUID,
        "interval": 5, "failedTimes": 5, "description": dev['name'],
        "offlineClear": 0, "offlineDefaultValue": "0", "deviceType": 2,
        "muid": muid, "configUid": "", "PageUUID": "",
        "extra": extra, "Status": 0, "longitude": "", "latitude": ""
    })
    if code in SUCCESS_CODES: n_dev += 1

print(f"  设备: {n_dev}/{len(dm_data.get('devices',[]))}")
print(f"  A20: 1-{slave_count['A20']}  A40: 61-{60+slave_count['A40']}  UPS: 70-{69+slave_count['施耐德UPS']}")

# ========== Step 9: 补建 device_real_data + 禁用告警 ==========
print("\n[9] 补建实时数据 + 禁用告警")
code, _ = api_proj("/syncDeviceRealData", {})
print(f"  syncDeviceRealData → code={code}")
code, _ = api_proj("/DeviceRealDataDisableAlarm", {})
print(f"  禁用告警 → code={code}")

# ========== Step 10: 验证 ==========
print("\n[10] 最终验证")
models = ob_read("SELECT uuid FROM devices_model WHERE project_uuid=%s AND deleted_at IS NULL", (PROJECT_UUID,))
muids = [m[0] for m in models]

stats = {}
stats['devices_model'] = len(models)
stats['monitor_list'] = ob_read("SELECT count(*) FROM monitor_list WHERE project_uuid=%s AND deleted_at IS NULL", (PROJECT_UUID,))[0][0]
stats['device_real_data'] = ob_read("SELECT count(*) FROM device_real_data WHERE project_uuid=%s", (PROJECT_UUID,))[0][0]
dev_count = ob_read("SELECT count(DISTINCT device_uuid) FROM device_real_data WHERE project_uuid=%s", (PROJECT_UUID,))[0][0]

if muids:
    ph = ','.join(['%s' for _ in muids])
    stats['register_groups'] = ob_read(f"SELECT count(*) FROM modbus_devices_register_group WHERE muid IN ({ph}) AND deleted_at IS NULL", muids)[0][0]
    stats['data_points'] = ob_read(f"SELECT count(*) FROM modbus_devices_data_model WHERE muid IN ({ph}) AND deleted_at IS NULL", muids)[0][0]

print(f"""
{'='*60}
  项目名称: {PROJECT_NAME}
  项目UUID: {PROJECT_UUID}
  ─────────────────────────────────────────
  devices_model:                   {stats.get('devices_model', 0):>4d} 个
  modbus_devices_register_group:   {stats.get('register_groups', 0):>4d} 个
  modbus_devices_data_model:       {stats.get('data_points', 0):>4d} 条
  monitor_list (设备实例):          {stats.get('monitor_list', 0):>4d} 个
  device_real_data (实时数据点):    {stats.get('device_real_data', 0):>4d} 条 ({dev_count} 个设备)
  ─────────────────────────────────────────
  A20 从站: 1 ~ {slave_count['A20']}
  A40 从站: 61 ~ {60+slave_count['A40']}
  UPS 从站: 70 ~ {69+slave_count['施耐德UPS']}
{'='*60}
""")

code, _ = api_proj("/MonitorBatchSetStatus", {"status": 1})
print(f"  已设置所有设备 status=1 (code={code})")
print("\n✅ 航信机房项目导入完成!")
