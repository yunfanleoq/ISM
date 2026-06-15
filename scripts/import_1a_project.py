#!/usr/bin/env python3
"""
ISM 完整项目导入脚本 (v3.0 — 全 API 调用，零 SQL 直写)
用法:
  python3 scripts/import_1a_project.py              # 默认: 更新已有项目
  python3 scripts/import_1a_project.py --new        # 新建项目
  python3 scripts/import_1a_project.py --name "2B配电室"  # 指定项目名

=== 两种模式 ===
  新建模式 (--new): 创建新 project_lists 记录，命名基于 PROJECT_NAME
  更新模式 (默认): 查找已有项目 → 全量清理旧数据 → 重新导入

=== v3.0 核心理念 ===
  所有写操作均通过后端 API 完成，零直接 SQL UPDATE/DELETE/INSERT。
  仅保留 SELECT 用于验证和 UUID 查询。

=== 新增后端 API ===
  POST /ProjectFixCreator           — 修正项目 creator_uuid
  POST /syncDeviceRealData           — 补建 device_real_data (type>=1)
  POST /DeviceRealDataDisableAlarm   — 批量禁用告警
  POST /MonitorBatchSetStatus        — 批量设置设备状态

=== 已固化的坑点 ===
  [×] project_lists.creator_uuid 必须等于 user 表 UUID（不是 JWT 解码）
  [×] Authorization header 不加 "Bearer " 前缀
  [×] ProjectUuid header 全流程携带
  [×] 设备名全局唯一 (跨项目查重，清理后再导入)
  [×] 禁用所有 is_alarm + 设置 alarm_shield=1
  [×] Modbus 从站 ID 与设备顺序严格对应: A20=1-60, A40=61-69, UPS=70-76
  [×] UUID 映射按模型名称动态匹配
  [×] 导入后清除系统默认 registerN 命名
  [×] user 表与 project_user 表的 admin 密码必须一致 (防登录走错路径)
"""

import json, sqlite3, subprocess, time, os, sys, uuid
from collections import defaultdict

# ========== 配置 ==========
BASE_URL = "http://localhost:8081"
DB_PATH = os.path.join(os.path.dirname(__file__), "..", "ism_server_user", "data", "db", "ism.db")
PROJECT_ROOT = os.path.join(os.path.dirname(__file__), "..")
PKG_PATH = os.path.join(PROJECT_ROOT, "liu-chang-1A-dev", "1A_complete_project_package.json")
MODELS_PATH = os.path.join(PROJECT_ROOT, "liu-chang-1A-dev", "ism_data_models.json")

# 命令行参数: 新建模式 vs 更新已有项目
IS_NEW = "--new" in sys.argv
ARG_NAME = None
for a in sys.argv[1:]:
    if a.startswith("--name="):
        ARG_NAME = a.split("=", 1)[1]
PROJECT_NAME = ARG_NAME or "1A配电室-模拟器"
MODE_LABEL = "新建项目" if IS_NEW else "更新已有项目"

SUCCESS_CODES = {0, 200, 2002, 4002, 3001}  # 3001=已存在也算成功

def api(url, json_data=None, headers=None):
    """统一 API 调用（使用 curl 子进程，避免 requests 库 header 丢失问题）"""
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
    """API 调用，自动携带 ProjectUuid header"""
    h = dict(HEADERS)
    h["ProjectUuid"] = PROJECT_UUID
    return api(url, json_data, headers=h)

def db_read(sql, params=()):
    """只读查询，返回 fetchall() 结果列表"""
    conn = sqlite3.connect(DB_PATH)
    rows = conn.execute(sql, params).fetchall()
    conn.close()
    return rows

def login(username="admin", password="123456"):
    """登录并返回 token 和用户 UUID（从 DB user 表查，避免 JWT 解析错误）"""
    payload = json.dumps({"username": username, "password": password})
    result = subprocess.run(
        ["curl", "-s", "-X", "POST", f"{BASE_URL}/login",
         "-H", "Content-Type:application/json",
         "-d", payload],
        capture_output=True, text=True, timeout=10
    )
    data = json.loads(result.stdout)
    token = data['data']['token']
    # ★ 不解析 JWT（可能因 project_user/user 表差异导致 uuid 不一致），直接查 DB ★
    db_conn = sqlite3.connect(DB_PATH)
    row = db_conn.execute("SELECT uuid, name FROM user WHERE username=? AND deleted_at IS NULL", (username,)).fetchone()
    db_conn.close()
    if row:
        user_uuid, user_name = row[0], row[1]
    else:
        user_uuid, user_name = "", "unknown"
    print(f"  登录成功: {user_name} (uuid={user_uuid[:24]}...)")
    return token, user_uuid, user_name

# ========== Step 0: 环境检查 ==========
print("=" * 60)
print("  ISM 完整项目导入 v3.0 (全API, 零SQL直写)")
print("=" * 60)

# 后端
print("\n[0] 环境检查")
pkill = subprocess.run(["pkill", "-9", "ism_server"], capture_output=True)
time.sleep(1)
subprocess.Popen(
    ["./ism_server"],
    cwd=os.path.join(PROJECT_ROOT, "ism_server_user"),
    stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL,
    preexec_fn=os.setpgrp
)
for i in range(30):
    time.sleep(1)
    try:
        result = subprocess.run(["curl", "-s", "-o", "/dev/null", "-w", "%{http_code}", f"{BASE_URL}/"],
            capture_output=True, text=True, timeout=3)
        if result.stdout.strip() in ("200", "404"):
            print(f"  后端就绪 ({i+1}s)")
            break
    except:
        pass
else:
    print("  后端启动超时!"); sys.exit(1)

# 模拟器
import socket
sock = socket.socket(); sock.settimeout(2)
sim_ok = False
try:
    sock.connect(("127.0.0.1", 502)); sock.close()
    sim_ok = True
except:
    pass
if sim_ok:
    print("  模拟器: OK (127.0.0.1:502)")
else:
    print("  模拟器: 不可用，跳过自动启动 (请手动启动)")

# ========== Step 1: 登录 + UUID ==========
print("\n[1] 登录获取用户 UUID")
TOKEN, USER_UUID, USER_NAME = login()
HEADERS = {"Content-Type": "application/json", "Authorization": TOKEN}
print(f"  Authorization: {TOKEN[:30]}...")

# ========== Step 2: 项目 (支持新建/更新双模式) ==========
print(f"\n[2] 项目 — {MODE_LABEL}: {PROJECT_NAME}")

existing = db_read(
    "SELECT uuid FROM project_lists WHERE name=? AND deleted_at IS NULL",
    (PROJECT_NAME,)
)

if IS_NEW:
    if existing:
        PROJECT_NAME = f"{PROJECT_NAME}-{time.strftime('%m%d%H%M')}"
        print(f"  已有同名项目，新建为: {PROJECT_NAME}")
        existing = None
    code, data = api("/ProjectAdd", {
        "name": PROJECT_NAME, "description": "ISM Modbus 监控项目 (AI生成)",
        "industry": 1
    })
    PROJECT_UUID = db_read(
        "SELECT uuid FROM project_lists WHERE name=? AND deleted_at IS NULL ORDER BY id DESC LIMIT 1",
        (PROJECT_NAME,)
    )[0][0]
    print(f"  新建项目: {PROJECT_UUID}")
else:
    if existing:
        PROJECT_UUID = existing[0][0]
        print(f"  已有项目: {PROJECT_UUID}")
    else:
        code, data = api("/ProjectAdd", {
            "name": PROJECT_NAME, "description": "ISM Modbus 监控项目 (AI生成)",
            "industry": 1
        })
        PROJECT_UUID = db_read(
            "SELECT uuid FROM project_lists WHERE name=? AND deleted_at IS NULL ORDER BY id DESC LIMIT 1",
            (PROJECT_NAME,)
        )[0][0]
        print(f"  自动创建项目: {PROJECT_UUID}")

# ★ v3.0: API 修正 creator_uuid 替代直接 SQL ★
HEADERS["ProjectUuid"] = PROJECT_UUID
code, _ = api_proj("/ProjectFixCreator", {})
print(f"  修正 creator_uuid → {USER_UUID[:24]}... (API: code={code})")

# ========== Step 3: 全量清理 ==========
print(f"\n[3] 清理旧数据（含跨项目冲突）")
if IS_NEW:
    print("  新建模式: 只清理同名冲突设备")
else:
    print("  更新模式: 全量清理当前项目 + 跨项目冲突")

# ★ v3.0: 无论如何都要清理跨项目同名设备冲突 ★
# Step 8 调用 monitorAdd 时，若同名设备在其他项目中已存在，
# 后端返回 3001 但新项目并未创建设备 → n_dev 计数正常但 DB 中无设备
# 必须在创建设备前先把所有旧项目中的同名设备删掉
with open(MODELS_PATH) as f:
    _dm_data = json.load(f)
_all_dev_names = [d['name'] for d in _dm_data.get('devices', [])]

if _all_dev_names:
    placeholders = ','.join(['?' for _ in _all_dev_names])
    conflicts = db_read(
        f"SELECT uuid, name, project_uuid FROM monitor_list WHERE name IN ({placeholders}) AND deleted_at IS NULL",
        _all_dev_names
    )
    if conflicts:
        # 按 project_uuid 分组，然后用 API 删除
        proj_devs = defaultdict(list)
        for dev_uuid, dev_name, proj_uuid in conflicts:
            proj_devs[proj_uuid].append(dev_uuid)
        
        for proj_uuid, dev_uuids in proj_devs.items():
            # 临时切换 ProjectUuid header
            old_puuid = HEADERS.get("ProjectUuid", "")
            HEADERS["ProjectUuid"] = proj_uuid
            code, _ = api("/monitorAllDel", {"uuid": dev_uuids}, headers=HEADERS)
            HEADERS["ProjectUuid"] = old_puuid
            print(f"  清理跨项目冲突: {len(dev_uuids)} 个设备 (project={proj_uuid[:12]}...) code={code}")
    else:
        print("  无跨项目冲突")
# 恢复 ProjectUuid
HEADERS["ProjectUuid"] = PROJECT_UUID

# 更新模式的额外清理
if not IS_NEW:
    # 清理当前项目的区域节点
    zones = db_read(
        "SELECT uuid FROM monitor_list WHERE project_uuid=? AND deleted_at IS NULL",
        (PROJECT_UUID,)
    )
    if zones:
        zone_uuids = [z[0] for z in zones]
        code, _ = api_proj("/monitorAllDel", {"uuid": zone_uuids})
        print(f"  区域节点: 删除 {len(zone_uuids)} 个 (API code={code})")
    
    # 清理模型（级联删除 register_group + data_model）
    models_to_del = db_read(
        "SELECT uuid, name FROM devices_model WHERE project_uuid=? AND deleted_at IS NULL",
        (PROJECT_UUID,)
    )
    for muid, mname in models_to_del:
        code, _ = api_proj("/modbusModelDel", {"uuid": muid})
        print(f"  {mname}: 模型已删除 (API code={code})")

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
# ★ v3.0: 用 API 修正 sid，不再直接 SQL ★
rz = db_read(
    "SELECT uuid FROM monitor_list WHERE name='RootZone' AND project_uuid=? AND deleted_at IS NULL LIMIT 1",
    (PROJECT_UUID,)
)
if rz:
    api_proj("/monitorEdit", {"data": {"Sid": 1, "uuid": rz[0][0]}})
print(f"  RootZone code={code}")

# ========== Step 5: 数据模型 ==========
print("\n[5] 创建数据模型 (devices_model)")
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
    row = db_read(
        "SELECT uuid FROM devices_model WHERE name=? AND project_uuid=? AND deleted_at IS NULL ORDER BY id DESC LIMIT 1",
        (mname, PROJECT_UUID)
    )
    if row:
        model_uuids[mname] = row[0][0]
        print(f"  {mname}: {row[0][0][:12]}... code={code}")
    else:
        print(f"  {mname}: FAILED code={code}")
        sys.exit(1)
print(f"  共 {len(model_uuids)} 个模型")

# ========== Step 6: 寄存器组 ==========
print("\n[6] 创建寄存器组 (modbus_devices_register_group)")
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
    row = db_read(
        "SELECT uuid FROM modbus_devices_register_group WHERE name=? AND muid=? AND deleted_at IS NULL ORDER BY id DESC LIMIT 1",
        (gname, muid)
    )
    k = f"{gname}_{muid}"
    if row:
        group_uuids[k] = row[0][0]
    print(f"  {gname}({mname}): {'OK' if row else 'FAIL'} code={code}")
print(f"  共 {len(group_uuids)} 个寄存器组")

# ========== Step 7: 数据点 (按模型名称映射 UUID) ==========
print("\n[7] 创建 Modbus 数据点 (modbus_devices_data_model)")
with open(PKG_PATH) as f:
    pkg = json.load(f)

# ★ v2.1: 按模型名称动态映射，不依赖硬编码的旧 UUID ★
# 包中的 deviceModels[i].name → DB 中同名模型的 uuid
pkg_model_names = [dm['name'] for dm in pkg['deviceModels']]
pkg_model_uuids = [dm['uuid'] for dm in pkg['deviceModels']]

# 模型名规范化映射: 包中可能是 'A20'/'A40'/'施耐德UPS' ↔ DB中 'A20电力仪表'/'A40电力仪表'/'施耐德UPS'
NAME_MAP = {
    'A20': 'A20电力仪表', 'A40': 'A40电力仪表', '施耐德UPS': '施耐德UPS',
    'A20电力仪表': 'A20电力仪表', 'A40电力仪表': 'A40电力仪表',
}

# pkg_muid → db_muid
pkg2db_muid = {}
for pkg_name, pkg_mu in zip(pkg_model_names, pkg_model_uuids):
    db_name = NAME_MAP.get(pkg_name)
    if not db_name:
        # 尝试模糊匹配
        for pk, dv in NAME_MAP.items():
            if pk in pkg_name:
                db_name = dv
                break
    if db_name and db_name in model_uuids:
        pkg2db_muid[pkg_mu] = model_uuids[db_name]
        print(f"  pkg:{pkg_name} → db:{db_name} ({model_uuids[db_name][:12]}...)")

# pkg_group_uuid → db_group_uuid
# 先用 pkg registerGroups 建立 (pkg_muid, group_name) → pkg_group_uuid
pkg_grp_map = {}
for rg in pkg.get('registerGroups', []):
    pkg_grp_map[(rg['muid'], rg['name'])] = rg['uuid']

# 再用 model_uuids 建立 (db_muid, group_name) → db_group_uuid
pkg_rg2db = {}
for (pkg_mu, gn), pkg_rg_uuid in pkg_grp_map.items():
    db_mu = pkg2db_muid.get(pkg_mu)
    if not db_mu:
        continue
    k = f"{gn}_{db_mu}"
    if k in group_uuids:
        pkg_rg2db[pkg_rg_uuid] = group_uuids[k]
    else:
        # 回退查找: 按 (gn, db_mu) 查 DB
        row = db_read(
            "SELECT uuid FROM modbus_devices_register_group WHERE name=? AND muid=? AND deleted_at IS NULL ORDER BY id DESC LIMIT 1",
            (gn, db_mu)
        )
        if row:
            pkg_rg2db[pkg_rg_uuid] = row[0][0]

print(f"  Mappings: model={len(pkg2db_muid)}, group={len(pkg_rg2db)}")

n_pts = 0
skipped = 0
for pt in pkg.get('registerPoints', []):
    pkg_mu = pt.get('muid', '')
    pkg_rg = pt.get('registerGroupUuid', '')
    db_mu = pkg2db_muid.get(pkg_mu)
    db_rg = pkg_rg2db.get(pkg_rg)
    if not db_mu or not db_rg:
        skipped += 1
        if skipped <= 3:
            print(f"  跳过 {pt.get('name','?'):20s}: pkg_muid={pkg_mu[:12] if pkg_mu else 'NONE'} pkg_rg={pkg_rg[:12] if pkg_rg else 'NONE'}")
        continue
    code, _ = api("/modbusModelRegisterAdd", {
        "name": pt['name'], "muid": db_mu,
        "registerAddress": pt['registerAddress'], "registerGroupUuid": db_rg,
        "auth": pt.get('auth', 'ReadOnly'), "type": pt.get('type', 'Float'),
        "ByteOrder": pt.get('ByteOrder', 'CDAB'), "modeltype": pt.get('modeltype', 2),
        "unit": pt.get('unit', ''),
        "conversionExpression": "",  # ★ 清空，防后端崩溃 ★
        "alarm": 0,                  # ★ 禁用告警 ★
        "alarmLevel": 0,
        "AlarmMessage": "", "AlarmClearMessage": "",
        "record": pt.get('record', 0), "RecordType": pt.get('RecordType', 1),
        "RecordInterval": pt.get('RecordInterval', 5),
        "RecordDataCharge": pt.get('RecordDataCharge', '0'),
        "RecordDataTimely": pt.get('RecordDataTimely', '0'),
        "FloatAccuracy": pt.get('FloatAccuracy', '0.01'),
    })
    if code in SUCCESS_CODES:
        n_pts += 1
print(f"  添加了 {n_pts}/{len(pkg.get('registerPoints',[]))} 个数据点 (跳过 {skipped})")

# ★ v3.0: 清除 registerN 默认命名（通过 API 组合实现，不再直接 SQL） ★
for mname in model_uuids:
    mu = model_uuids[mname]
    all_pts = db_read(
        "SELECT uuid, name FROM modbus_devices_data_model WHERE muid=? AND deleted_at IS NULL",
        (mu,)
    )
    register_n = [pt[0] for pt in all_pts if pt[1].lower().startswith('register')]
    if register_n:
        code, _ = api_proj("/modbusModelRegisterDel", {"uuid": register_n})
        print(f"  清除 {mname} 的 {len(register_n)} 个 registerN 默认命名 (API code={code})")

# ========== Step 8: 设备树 + 设备实例 ==========
print("\n[8] 设备树和实例 (monitor_list)")

# 加载设备模型数据
with open(MODELS_PATH) as f:
    dm = json.load(f)

# 创建设备树节点
def get_sid(name, dt):
    r = db_read(
        "SELECT sid FROM monitor_list WHERE name=? AND type=? AND project_uuid=? AND deleted_at IS NULL",
        (name, dt, PROJECT_UUID)
    )
    return r[0][0] if r else None

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

root = get_sid("RootZone", 0) or 1
site = add_zone(root, "1A配电室")
floor = add_zone(site, "1楼配电室")
cab_a1 = add_zone(floor, "1A1_U11柜")
cab_a3 = add_zone(floor, "1A3_U11柜")
cab_ups = add_zone(floor, "UPS柜")
print(f"  站点:{site} 楼层:{floor} A1柜:{cab_a1} A3柜:{cab_a3} UPS柜:{cab_ups}")

# 设备实例
TEMPLATE_MAP = {"A20": "A20电力仪表", "A40": "A40电力仪表", "施耐德UPS": "施耐德UPS"}
# ★ 坑点: 之前按设备型号硬编码柜号 (A20→1A1柜, A40→1A3柜)，
# 导致1A3楼的A20设备和1A1楼的A40设备全部归错柜。
# 修复: 解析设备名中的楼层号 (1A1_/1A3_) 动态分配 ★
def get_cab_for_device(dev_name):
    """根据设备名中的楼层号返回正确的柜号"""
    if '1A3_' in dev_name or '1A3U' in dev_name:
        return cab_a3  # 1A3_U11柜
    elif '1A1_' in dev_name or '1A1U' in dev_name:
        return cab_a1  # 1A1_U11柜
    # 默认: UPS设备也根据楼层分配, 模糊匹配
    if 'UPS' in dev_name or '施耐德' in dev_name:
        if '1A3' in dev_name:
            return cab_a3
        return cab_a1
    return cab_a1  # fallback

slave_count = {"A20": 0, "A40": 0, "施耐德UPS": 0}
n_dev = 0

for dev in dm.get('devices', []):
    t = dev['templateType']
    muid = model_uuids[TEMPLATE_MAP[t]]
    pid = get_cab_for_device(dev['name'])  # ★ 动态柜号,不再按型号硬编码
    slave_count[t] += 1

    # ★ 从站地址映射: A20=1-60, A40=61-69, UPS=70-76 ★
    if t == "A20":
        slave = slave_count[t]
    elif t == "A40":
        slave = 60 + slave_count[t]
    else:
        slave = 69 + slave_count[t]

    extra = json.dumps({
        "modbus": {
            "IPAddress": "127.0.0.1",
            "Port": "502",
            "address": str(slave),
            "RegisterPack": -1
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
    if code in SUCCESS_CODES:
        n_dev += 1

print(f"  设备: {n_dev}/{len(dm.get('devices',[]))}")
print(f"  A20: 1-{slave_count['A20']}  A40: 61-{60+slave_count['A40']}  UPS: 70-{69+slave_count['施耐德UPS']}")

# ========== Step 9: 补建 device_real_data ==========
print("\n[9] 补建实时数据 + 禁用告警 (API)")
# ★ v3.0: 全部使用 API，不再直接 SQL ★
# 1) 补建 device_real_data
code, _ = api_proj("/syncDeviceRealData", {})
print(f"  syncDeviceRealData → code={code}")

# 2) 批量禁用告警
code, _ = api_proj("/DeviceRealDataDisableAlarm", {})
print(f"  已禁用所有告警 (API code={code})")

# ========== Step 10: 验证 ==========
print("\n[10] 最终验证")
stats = {}
for tbl, name_col in [("devices_model", "name"), ("monitor_list", "name"),
                        ("device_real_data", "device_name")]:
    cnt = db_read(f"SELECT count(*) FROM {tbl} WHERE project_uuid=?", (PROJECT_UUID,))[0][0]
    stats[tbl] = cnt

# modbus_devices_register_group 和 modbus_devices_data_model 用 muid 关联
models = db_read("SELECT uuid FROM devices_model WHERE project_uuid=? AND deleted_at IS NULL", (PROJECT_UUID,))
muids = [m[0] for m in models]
if muids:
    placeholders = ','.join(['?' for _ in muids])
    stats['modbus_devices_register_group'] = db_read(f"SELECT count(*) FROM modbus_devices_register_group WHERE muid IN ({placeholders}) AND deleted_at IS NULL", muids)[0][0]
    stats['modbus_devices_data_model'] = db_read(f"SELECT count(*) FROM modbus_devices_data_model WHERE muid IN ({placeholders}) AND deleted_at IS NULL", muids)[0][0]
else:
    stats['modbus_devices_register_group'] = 0
    stats['modbus_devices_data_model'] = 0

dev_types = db_read(
    "SELECT count(*), count(DISTINCT device_uuid) FROM device_real_data WHERE project_uuid=?",
    (PROJECT_UUID,)
)[0]

print(f"""
{'='*60}
  项目名称: {PROJECT_NAME}
  项目UUID: {PROJECT_UUID}
  creator_uuid: {USER_UUID[:16]}...
  ─────────────────────────────────────────
  devices_model:                   {stats.get('devices_model', 0):>4d} 个
  modbus_devices_register_group:   {stats.get('modbus_devices_register_group', 0):>4d} 个
  modbus_devices_data_model:       {stats.get('modbus_devices_data_model', 0):>4d} 条
  monitor_list (设备实例):          {stats.get('monitor_list', 0):>4d} 个
  device_real_data (实时数据点):    {stats.get('device_real_data', 0):>4d} 条 ({dev_types[1]} 个设备)
  ─────────────────────────────────────────
  A20 从站: 1 ~ {slave_count['A20']}
  A40 从站: 61 ~ {60+slave_count['A40']}
  UPS 从站: 70 ~ {69+slave_count['施耐德UPS']}
  ─────────────────────────────────────────
  模拟器 API: http://localhost:5040/api/slaves
  前端页面:   http://localhost:8080/#/SimulatorMonitor
{'='*60}
""")

# ★ v3.0: 用 API 批量设置设备在线状态 ★
code, _ = api_proj("/MonitorBatchSetStatus", {"status": 1})
print(f"  已设置所有设备 status=1 (API code={code})")
print("\n✅ 导入完成!")
