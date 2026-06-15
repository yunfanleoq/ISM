#!/usr/bin/env python3
"""
ISM 1A配电室项目导入脚本
通过后端API批量导入设备模型、寄存器组、数据点、设备树、告警和组态
"""
import json
import requests
import uuid
from pathlib import Path

BASE_URL = "http://localhost:8081"
TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXIiLCJyb2xlIjoiVXNlciIsIm5hbWUiOiJ1c2VyIiwidXVpZCI6ImY4ZmE3MWFjLTBhNjEtNzlhYy01ZjkwLTJlZTNlNzEwYmU2MiIsImV4cCI6MTc4Mzg3MDU5MywiaXNzIjoiZ2luYmxvZyJ9.Hdwsfq8xWwdpVs9iX8YK8DXK1NQB05ZFMBN72y5vU0E"
PROJECT_UUID = "420a696f-2e2e-53cc-e616-1f2ccc11f51c"

HEADERS = {
    "Content-Type": "application/json",
    "Authorization": TOKEN,
    "ProjectUuid": PROJECT_UUID
}

new_uuid = lambda: str(uuid.uuid4()).replace('-', '')

def log_response(name, resp):
    try:
        data = resp.json()
        print(f"  [{name}] code={data.get('code', '?')}, msg={data.get('msg', data.get('message', ''))}")
        return data
    except:
        print(f"  [{name}] status={resp.status_code}, body={resp.text[:100]}")
        return None

# ========== 1. 创建设备模型 ==========
print("\n=== 1. 创建设备模型 ===")
models = [
    {
        "name": "A20电力仪表",
        "dec": "A20 基础电力仪表模型（30AI+3DI）",
        "type": 2,
        "gatherNumber": 30,
        "project_uuid": PROJECT_UUID,
        "modbusConnectType": "TCP",
        "modbusConnectMode": "Client",
        "modbusClientIpaddress": "172.31.4.14",
        "DataFormat": "CDAB",
        "timeout": 5,
        "port": 502
    },
    {
        "name": "A40电力仪表",
        "dec": "A40 增强电力仪表模型（41AI+3DI）",
        "type": 2,
        "gatherNumber": 41,
        "project_uuid": PROJECT_UUID,
        "modbusConnectType": "TCP",
        "modbusConnectMode": "Client",
        "modbusClientIpaddress": "172.31.4.14",
        "DataFormat": "CDAB",
        "timeout": 5,
        "port": 502
    },
    {
        "name": "施耐德UPS",
        "dec": "施耐德UPS监控模型（44AI+1DI）",
        "type": 2,
        "gatherNumber": 44,
        "project_uuid": PROJECT_UUID,
        "modbusConnectType": "TCP",
        "modbusConnectMode": "Client",
        "modbusClientIpaddress": "172.31.4.14",
        "DataFormat": "CDAB",
        "timeout": 5,
        "port": 502
    }
]

model_uuids = {}
for m in models:
    resp = requests.post(f"{BASE_URL}/modbusModelAdd", headers=HEADERS, json=m)
    data = log_response(m['name'], resp)
    if data and data.get('code') == 0:
        # 获取刚创建的模型UUID
        list_resp = requests.post(f"{BASE_URL}/modbusModelList", headers=HEADERS, json={"type": 2})
        list_data = list_resp.json()
        if list_data.get('code') == 0:
            for item in list_data.get('list', []):
                if item['name'] == m['name']:
                    model_uuids[m['name']] = item['uuid']
                    print(f"    -> UUID: {item['uuid']}")
                    break

print(f"  已创建 {len(model_uuids)} 个模型")

# ========== 2. 创建寄存器组 ==========
print("\n=== 2. 创建寄存器组 ===")
register_groups = []
for model_name, muid in model_uuids.items():
    if model_name == "A20电力仪表":
        register_groups.append({"name": "A20_AI数据", "muid": muid, "function": 3, "registerStart": 0, "registerCount": 30})
        register_groups.append({"name": "A20_DI数据", "muid": muid, "function": 2, "registerStart": 0, "registerCount": 3})
    elif model_name == "A40电力仪表":
        register_groups.append({"name": "A40_AI数据", "muid": muid, "function": 3, "registerStart": 0, "registerCount": 41})
        register_groups.append({"name": "A40_DI数据", "muid": muid, "function": 2, "registerStart": 0, "registerCount": 3})
    elif model_name == "施耐德UPS":
        register_groups.append({"name": "UPS_AI数据", "muid": muid, "function": 3, "registerStart": 0, "registerCount": 44})
        register_groups.append({"name": "UPS_DI数据", "muid": muid, "function": 2, "registerStart": 0, "registerCount": 1})

group_uuids = {}
for g in register_groups:
    resp = requests.post(f"{BASE_URL}/modbusModelRegisterGroupAdd", headers=HEADERS, json=g)
    data = log_response(g['name'], resp)
    if data and data.get('code') == 0:
        # 获取刚创建的组UUID
        list_resp = requests.post(f"{BASE_URL}/modbusModelRegisterGroupList", headers=HEADERS, json={"muid": g['muid']})
        list_data = list_resp.json()
        if list_data.get('code') == 0:
            for item in list_data.get('list', []):
                if item['name'] == g['name']:
                    group_uuids[g['name']] = item['uuid']
                    print(f"    -> UUID: {item['uuid']}")
                    break

print(f"  已创建 {len(group_uuids)} 个寄存器组")

# ========== 3. 添加寄存器数据点（简化版） ==========
print("\n=== 3. 添加寄存器数据点 ===")
# 读取预生成的数据点配置
with open('/Users/yunfanleo/cursorProjects/ISM源码/liu-chang-1A-dev/1A_complete_project_package.json', 'r') as f:
    pkg = json.load(f)

points_added = 0
for pt in pkg['registerPoints'][:20]:  # 先导入前20个测试
    pt['uuid'] = new_uuid()
    resp = requests.post(f"{BASE_URL}/modbusModelRegisterAdd", headers=HEADERS, json=pt)
    data = log_response(pt['name'], resp)
    if data and data.get('code') == 0:
        points_added += 1

print(f"  已添加 {points_added} 个数据点")

# ========== 4. 创建设备树 ==========
print("\n=== 4. 创建设备树 ===")
# 创建根区域
zones = [
    {"sid": 1000, "pid": 0, "name": "1A配电室", "type": 0, "description": "1A配电室"},
    {"sid": 1001, "pid": 1000, "name": "1楼配电室", "type": 0, "description": "1楼配电室"},
    {"sid": 1002, "pid": 1001, "name": "1A1_U11柜", "type": 0, "description": "1A1_U11柜"},
    {"sid": 1003, "pid": 1001, "name": "1A3_U11柜", "type": 0, "description": "1A3_U11柜"},
    {"sid": 1004, "pid": 1001, "name": "UPS柜", "type": 0, "description": "UPS柜"},
]

for z in zones:
    payload = {
        "sid": z['sid'],
        "pid": z['pid'],
        "name": z['name'],
        "type": z['type'],
        "timeout": 5,
        "IsEnable": 1,
        "project_uuid": PROJECT_UUID,
        "interval": 5,
        "failedTimes": 5,
        "description": z['description'],
        "offlineClear": 0,
        "offlineDefaultValue": "0",
        "deviceType": 0,
        "muid": "",
        "configUid": "",
        "PageUUID": "",
        "extra": "",
        "Status": 0,
        "longitude": "",
        "latitude": ""
    }
    resp = requests.post(f"{BASE_URL}/monitorAdd", headers=HEADERS, json=payload)
    log_response(z['name'], resp)

# 添加设备实例（先添加5个测试）
print("\n  添加设备实例...")
with open('/Users/yunfanleo/cursorProjects/ISM源码/liu-chang-1A-dev/ism_data_models.json', 'r') as f:
    data_models = json.load(f)

cabinet_map = {
    "A20": 1002,
    "A40": 1002,
    "施耐德UPS": 1004
}

devices_added = 0
for dev in data_models['devices'][:5]:
    model_type = dev['templateType']
    model_name = None
    if model_type == "A20":
        model_name = "A20电力仪表"
    elif model_type == "A40":
        model_name = "A40电力仪表"
    elif model_type == "施耐德UPS":
        model_name = "施耐德UPS"
    
    muid = model_uuids.get(model_name, "")
    pid = cabinet_map.get(model_type, 1002)
    
    payload = {
        "sid": 2000 + devices_added,
        "pid": pid,
        "name": dev['name'],
        "type": 1,
        "timeout": 5,
        "IsEnable": 1,
        "project_uuid": PROJECT_UUID,
        "interval": 5,
        "failedTimes": 5,
        "description": dev['name'],
        "offlineClear": 0,
        "offlineDefaultValue": "0",
        "deviceType": 2,
        "muid": muid,
        "configUid": "",
        "PageUUID": "",
        "extra": json.dumps({"aiStartAddr": dev['aiStartAddr'], "diStartAddr": dev['diStartAddr']})
    }
    resp = requests.post(f"{BASE_URL}/monitorAdd", headers=HEADERS, json=payload)
    log_response(dev['name'], resp)
    devices_added += 1

print(f"  已添加 {devices_added} 个设备")

print("\n✅ 导入完成！请访问 http://localhost:8080 查看效果")
print("   项目UUID:", PROJECT_UUID)
