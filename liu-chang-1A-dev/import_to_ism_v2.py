#!/usr/bin/env python3
"""
ISM 1A配电室项目导入脚本 - 修正版
通过后端API批量导入，查询数据库获取后端生成的sid
"""
import json
import requests
import sqlite3
from pathlib import Path

BASE_URL = "http://localhost:8081"
TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXIiLCJyb2xlIjoiVXNlciIsIm5hbWUiOiJ1c2VyIiwidXVpZCI6ImY4ZmE3MWFjLTBhNjEtNzlhYy01ZjkwLTJlZTNlNzEwYmU2MiIsImV4cCI6MTc4Mzg3MDU5MywiaXNzIjoiZ2luYmxvZyJ9.Hdwsfq8xWwdpVs9iX8YK8DXK1NQB05ZFMBN72y5vU0E"
PROJECT_UUID = "420a696f-2e2e-53cc-e616-1f2ccc11f51c"
DB_PATH = "/Users/yunfanleo/cursorProjects/ISM源码/ism_server_user/data/db/ism.db"

HEADERS = {
    "Content-Type": "application/json",
    "Authorization": TOKEN,
    "ProjectUuid": PROJECT_UUID
}

def log_response(name, resp):
    try:
        data = resp.json()
        print(f"  [{name}] code={data.get('code', '?')}, msg={data.get('msg', data.get('message', ''))}")
        return data
    except:
        print(f"  [{name}] status={resp.status_code}, body={resp.text[:100]}")
        return None

def get_sid_from_db(name, type_val):
    """从数据库查询刚创建区域的sid"""
    conn = sqlite3.connect(DB_PATH)
    cursor = conn.cursor()
    cursor.execute("SELECT sid FROM monitor_list WHERE name=? AND type=? ORDER BY id DESC LIMIT 1", (name, type_val))
    row = cursor.fetchone()
    conn.close()
    return row[0] if row else None

# ========== 1. 获取现有设备模型UUID ==========
print("\n=== 1. 获取现有设备模型 ===")
resp = requests.post(f"{BASE_URL}/modbusModelList", headers=HEADERS, json={"type": 2})
data = resp.json()
model_uuids = {}
if data.get('code') == 0:
    for item in data.get('list', []):
        model_uuids[item['name']] = item['uuid']
        print(f"  {item['name']}: {item['uuid']}")

# ========== 2. 创建寄存器组 ==========
print("\n=== 2. 创建寄存器组 ===")
register_groups = [
    {"name": "A20_AI数据", "muid": model_uuids.get("A20电力仪表"), "function": 3, "registerStart": 0, "registerCount": 30},
    {"name": "A20_DI数据", "muid": model_uuids.get("A20电力仪表"), "function": 2, "registerStart": 0, "registerCount": 3},
    {"name": "A40_AI数据", "muid": model_uuids.get("A40电力仪表"), "function": 3, "registerStart": 0, "registerCount": 41},
    {"name": "A40_DI数据", "muid": model_uuids.get("A40电力仪表"), "function": 2, "registerStart": 0, "registerCount": 3},
    {"name": "UPS_AI数据", "muid": model_uuids.get("施耐德UPS"), "function": 3, "registerStart": 0, "registerCount": 44},
    {"name": "UPS_DI数据", "muid": model_uuids.get("施耐德UPS"), "function": 2, "registerStart": 0, "registerCount": 1},
]

group_uuids = {}
for g in register_groups:
    if not g['muid']:
        continue
    resp = requests.post(f"{BASE_URL}/modbusModelRegisterGroupAdd", headers=HEADERS, json=g)
    data = log_response(g['name'], resp)
    if data and data.get('code') == 0:
        list_resp = requests.post(f"{BASE_URL}/modbusModelRegisterGroupList", headers=HEADERS, json={"muid": g['muid']})
        list_data = list_resp.json()
        if list_data.get('code') == 0:
            for item in list_data.get('list', []):
                if item['name'] == g['name']:
                    group_uuids[g['name']] = item['uuid']
                    print(f"    -> UUID: {item['uuid']}")
                    break

print(f"  已创建/获取 {len(group_uuids)} 个寄存器组")

# ========== 3. 添加寄存器数据点 ==========
print("\n=== 3. 添加寄存器数据点 ===")
with open('/Users/yunfanleo/cursorProjects/ISM源码/liu-chang-1A-dev/1A_complete_project_package.json', 'r') as f:
    pkg = json.load(f)

points_added = 0
for pt in pkg['registerPoints']:
    # 查找对应的寄存器组UUID
    group_name = None
    if "A20" in pt['name'] or pt['muid'] == model_uuids.get("A20电力仪表"):
        if pt['type'] == 'Bool':
            group_name = "A20_DI数据"
        else:
            group_name = "A20_AI数据"
    elif "A40" in pt['name'] or pt['muid'] == model_uuids.get("A40电力仪表"):
        if pt['type'] == 'Bool':
            group_name = "A40_DI数据"
        else:
            group_name = "A40_AI数据"
    elif pt['muid'] == model_uuids.get("施耐德UPS"):
        if pt['type'] == 'Bool':
            group_name = "UPS_DI数据"
        else:
            group_name = "UPS_AI数据"
    
    group_uuid = group_uuids.get(group_name, "")
    if not group_uuid:
        continue
    
    pt['registerGroupUuid'] = group_uuid
    pt.pop('uuid', None)  # 让后端生成新UUID
    
    resp = requests.post(f"{BASE_URL}/modbusModelRegisterAdd", headers=HEADERS, json=pt)
    data = log_response(pt['name'], resp)
    if data and data.get('code') == 0:
        points_added += 1

print(f"  已添加 {points_added} 个数据点")

# ========== 4. 创建设备树 ==========
print("\n=== 4. 创建设备树 ===")

# 创建根区域
zones = [
    {"name": "1A配电室", "type": 0, "pid": 1},  # pid=1 = RootZone
]

zone_sids = {}
for z in zones:
    payload = {
        "sid": 0, "pid": z['pid'], "name": z['name'], "type": z['type'],
        "timeout": 5, "IsEnable": 1, "project_uuid": PROJECT_UUID,
        "interval": 5, "failedTimes": 5, "description": z['name'],
        "offlineClear": 0, "offlineDefaultValue": "0", "deviceType": 0,
        "muid": "", "configUid": "", "PageUUID": "", "extra": "",
        "Status": 0, "longitude": "", "latitude": ""
    }
    resp = requests.post(f"{BASE_URL}/monitorAdd", headers=HEADERS, json=payload)
    log_response(z['name'], resp)
    sid = get_sid_from_db(z['name'], z['type'])
    if sid:
        zone_sids[z['name']] = sid
        print(f"    -> sid: {sid}")

# 创建配电柜子区域
cabinets = ["1A1_U11柜", "1A3_U11柜", "UPS柜"]
cabinet_sids = {}
for cab in cabinets:
    pid = zone_sids.get("1A配电室", 1)
    payload = {
        "sid": 0, "pid": pid, "name": cab, "type": 0,
        "timeout": 5, "IsEnable": 1, "project_uuid": PROJECT_UUID,
        "interval": 5, "failedTimes": 5, "description": cab,
        "offlineClear": 0, "offlineDefaultValue": "0", "deviceType": 0,
        "muid": "", "configUid": "", "PageUUID": "", "extra": "",
        "Status": 0, "longitude": "", "latitude": ""
    }
    resp = requests.post(f"{BASE_URL}/monitorAdd", headers=HEADERS, json=payload)
    log_response(cab, resp)
    sid = get_sid_from_db(cab, 0)
    if sid:
        cabinet_sids[cab] = sid
        print(f"    -> sid: {sid}")

# 添加设备实例
print("\n  添加设备实例...")
with open('/Users/yunfanleo/cursorProjects/ISM源码/liu-chang-1A-dev/ism_data_models.json', 'r') as f:
    data_models = json.load(f)

cabinet_map = {
    "A20": "1A1_U11柜",
    "A40": "1A3_U11柜",
    "施耐德UPS": "UPS柜"
}

devices_added = 0
for dev in data_models['devices']:
    model_type = dev['templateType']
    model_name_map = {"A20": "A20电力仪表", "A40": "A40电力仪表", "施耐德UPS": "施耐德UPS"}
    model_name = model_name_map.get(model_type)
    muid = model_uuids.get(model_name, "")
    
    cabinet_name = cabinet_map.get(model_type, "1A1_U11柜")
    pid = cabinet_sids.get(cabinet_name, 1)
    
    payload = {
        "sid": 0, "pid": pid, "name": dev['name'], "type": 1,
        "timeout": 5, "IsEnable": 1, "project_uuid": PROJECT_UUID,
        "interval": 5, "failedTimes": 5, "description": dev['name'],
        "offlineClear": 0, "offlineDefaultValue": "0", "deviceType": 2,
        "muid": muid, "configUid": "", "PageUUID": "",
        "extra": json.dumps({"aiStartAddr": dev['aiStartAddr'], "diStartAddr": dev['diStartAddr']})
    }
    resp = requests.post(f"{BASE_URL}/monitorAdd", headers=HEADERS, json=payload)
    log_response(dev['name'], resp)
    devices_added += 1

print(f"  已添加 {devices_added} 个设备")

print("\n✅ 导入完成！")
print("   请访问 http://localhost:8080 登录查看")
print("   用户名: user, 密码: 123456")
