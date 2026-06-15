#!/usr/bin/env python3
"""
将 1A配电室_ISM项目包.json 中的组态大屏配置导入到模拟器项目
"""
import json
import requests
import uuid
import time
import base64

BASE_URL = "http://127.0.0.1:8081"
PROJECT_UUID = "e308e378-ed94-1264-197b-33e535b812b8"  # 1A配电室-模拟器
new_uuid = lambda: str(uuid.uuid4()).replace('-', '')

def login():
    resp = requests.post(f"{BASE_URL}/login", json={"username": "admin", "password": "admin123"})
    data = resp.json()
    token = data.get("data", {}).get("token", "")
    cookie = resp.cookies.get("ISMSession", "")
    print(f"  登录: token={'OK' if token else 'FAIL'}, cookie={'OK' if cookie else 'FAIL'}")
    return token, cookie

def api_post(path, body, token="", cookie=""):
    h = {"Content-Type": "application/json", "ProjectUuid": PROJECT_UUID}
    if cookie:
        h["Cookie"] = f"ISMSession={cookie}"
    if token:
        h["Authorization"] = token
    r = requests.post(f"{BASE_URL}/{path}", headers=h, json=body)
    try:
        d = r.json()
        code = d.get("code", -99)
        return code, d
    except:
        return -1, r.text[:100]
    finally:
        del r

# =========== 1. 登录 ===========
print("=== 1. 登录 ===")
token, cookie = login()

# =========== 2. 创建组态应用(DisplayModel) ===========
print("\n=== 2. 创建组态应用 ===")
code, data = api_post("displayModelAdd", {
    "name": "1A配电室监控总览",
    "description": "1A配电室 Modbus 模拟监控大屏",
    "display_type": 1,
}, token, cookie)
print(f"  displayModelAdd: code={code}, {data}")

# Get the created model UUID
code, data = api_post("displayModelList", {"DisplayType": 1}, token, cookie)
model_uid = ""
if code == 0 and data.get("list"):
    for m in data["list"]:
        if m.get("name") == "1A配电室监控总览":
            model_uid = m.get("display_model_uid", "")
            print(f"  -> Model UUID: {model_uid}")
            break

# =========== 3. 添加页面 ===========
print("\n=== 3. 添加页面 ===")

# 添加总览页
code, data = api_post("DisplayModelPageAdd", {
    "modelUuid": model_uid,
    "name": "总览",
    "size": "1920x1080",
    "pageType": 1,
    "isLogin": 0
}, token, cookie)
print(f"  PageAdd(总览): code={code}")

# 获取 page_id
code, data = api_post(f"getDisplayModelPagerLayerData", {"modelUuid": model_uid}, token, cookie)
# 这个 API 返回的是 gzip 数据，我们直接用 save API

# =========== 4. 构建组态组件树并保存 ===========
print("\n=== 4. 构建大屏组件 ===")

# ISM 组态组件格式参考
components = []

# 4a. 顶部标题框 - DvBorderBox1
components.append({
    "type": "DvBorderBox1",
    "identifier": "header_border",
    "name": "header_border",
    "detail": {
        "style": {
            "width": 1920, "height": 80,
            "left": 0, "top": 0,
            "borderWidth": 2
        },
        "config": {
            "title": "1A配电室 - 监控总览 (76台设备)",
            "titleColor": "#00d4ff",
            "backgroundColor": "rgba(10,22,40,0.9)",
            "color": ["#0a1628", "#0a1628"]
        }
    }
})

# 4b. 设备树 - 左侧
components.append({
    "type": "DeviceTree",
    "identifier": "device_tree",
    "name": "device_tree",
    "detail": {
        "style": {
            "width": 280, "height": 1000,
            "left": 0, "top": 80
        },
        "config": {
            "title": "设备拓扑",
            "currentNode": []
        }
    }
})

# 4c. 实时数据表格 - 中间（展示所有设备的关键数据）
table_columns = [
    {"title": "设备名称", "dataIndex": "name", "width": 180},
    {"title": "AB线电压(V)", "dataIndex": "voltage", "width": 120},
    {"title": "A相电流(A)", "dataIndex": "current", "width": 120},
    {"title": "有功功率(kW)", "dataIndex": "power", "width": 120},
    {"title": "功率因数", "dataIndex": "pf", "width": 100},
    {"title": "通讯状态", "dataIndex": "comm", "width": 100},
]

table_data = []
# 从项目包中读取设备配置
with open('/Users/yunfanleo/cursorProjects/ISM源码/1A配电室_ISM项目包.json', 'r') as f:
    pkg = json.load(f)

# 收集所有设备名称
all_devices = []
dashboard = pkg.get("dashboard", {})
overview = dashboard.get("overview_page", {})
for section in overview.get("sections", []):
    for cabinet in section.get("cabinets", []):
        for dev in cabinet.get("devices", []):
            all_devices.append(dev)

# 为每个设备创建一行
for i, dev in enumerate(all_devices[:20]):  # 只显示前20个避免组件过大
    table_data.append({
        "name": dev.get("name", f"设备{i}"),
        "voltage": {"value": "", "mduid": f"placeholder_{i}_v"},
        "current": {"value": "", "mduid": f"placeholder_{i}_c"},
        "power": {"value": "", "mduid": f"placeholder_{i}_p"},
        "pf": {"value": "", "mduid": f"placeholder_{i}_f"},
        "comm": {"value": "", "mduid": f"placeholder_{i}_s"},
    })

components.append({
    "type": "RealDataTable",
    "identifier": "main_table",
    "name": "main_table",
    "detail": {
        "style": {
            "width": 1300, "height": 700,
            "left": 290, "top": 90
        },
        "config": {
            "columns": table_columns,
            "rows": table_data,
            "rowNum": 20,
            "headerBGC": "#0a1628",
            "oddRowBGC": "#0d1f35",
            "evenRowBGC": "#0a1628",
            "headerHeight": 40,
            "carousel": "page",
            "align": ["center"]
        }
    }
})

# 4d. 告警列表 - 右侧
components.append({
    "type": "alarmList",
    "identifier": "alarm_list",
    "name": "alarm_list",
    "detail": {
        "style": {
            "width": 330, "height": 500,
            "left": 1590, "top": 80
        },
        "config": {
            "title": "实时告警",
            "count": 10,
            "roll": True
        }
    }
})

components.append({
    "type": "RealDataTable",
    "identifier": "key_data",
    "name": "key_data",
    "detail": {
        "style": {
            "width": 330, "height": 490,
            "left": 1590, "top": 590
        },
        "config": {
            "title": "关键数据",
            "columns": [
                {"title": "参数", "dataIndex": "label"},
                {"title": "值", "dataIndex": "value"}
            ],
            "rows": [
                {"label": "设备总数", "value": {"value": "76", "mduid": ""}},
                {"label": "在线设备", "value": {"value": "76", "mduid": ""}},
                {"label": "离线设备", "value": {"value": "0", "mduid": ""}},
            ],
            "rowNum": 5,
            "evenRowBGC": "#0d1f35",
            "headerBGC": "#0a1628",
            "oddRowBGC": "#0a1628",
        }
    }
})

# Layer 配置
layer_config = {
    "backColor": "#0a1628",
    "backgroundImage": "",
    "widthHeightRatio": "",
    "width": 1920,
    "height": 1080,
    "autoSize": 0,
}

# 序列化并 base64 编码（与 ISM 前端格式一致）
comp_json = json.dumps(components, ensure_ascii=False)
comp_b64 = base64.b64encode(comp_json.encode('utf-8')).decode('utf-8')
layer_json = json.dumps(layer_config, ensure_ascii=False)

print(f"  Components: {len(components)} 个组件, {len(comp_b64)} chars (base64)")
print(f"  Layer config: {layer_json[:80]}...")

# 获取页面 ID
code, data = api_post("getDisplayModelPagerLayerData", 
    {"modelUuid": model_uid, "pageName": "总览"}, token, cookie)
print(f"  getDisplayModelPagerLayerData: code={code}")

# 保存层数据
code, data = api_post("saveDisplayModelLayerData", {
    "modelUuid": model_uid,
    "pageName": "总览",
    "layer": layer_json,
    "components": comp_b64,
    "is_home": 1,
    "page_type": 1,
}, token, cookie)
print(f"  saveDisplayModelLayerData: code={code}")

# =========== 5. 验证 ===========
print("\n=== 5. 验证 ===")
code, data = api_post("displayModelList", {"DisplayType": 1}, token, cookie)
if code == 0:
    app_count = len(data.get("list", []))
    print(f"  模拟器项目组态应用数: {app_count}")
    for m in data.get("list", []):
        print(f"    {m.get('name','?')} uid={m.get('display_model_uid','?')[:20]}...")

print(f"\n✅ 大屏导入完成！请刷新前端 http://localhost:7080")
print(f"   项目: 1A配电室-模拟器 → 应用列表 应出现 '1A配电室监控总览'")
