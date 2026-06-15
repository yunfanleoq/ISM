#!/usr/bin/env python3
"""修复 1A配电室-模拟器 的组态大屏组件数据，使用正确的 ISM 组件格式"""
import sqlite3, json, base64
from datetime import datetime

DB = '/Users/yunfanleo/cursorProjects/ISM源码/ism_server_user/data/db/ism.db'
MODEL_UID = 'f7cc6ef3-afd2-4e28-aa75-a782683a84f9'

# 通用动画字段（每个组件必备）
def default_animate():
    return {
        "selected": [],
        "condition": {
            "deviceSN": "", "isBandDevice": False, "bandType": 1,
            "dataID": "", "dataName": "",
            "operator": "", "OperatorValue": "", "OperatorMaxValue": ""
        },
        "isExpression": False,
        "animateList": [],
        "animateElement": []
    }

# 通用 style 基础字段
def default_style(x, y, w, h):
    return {
        "position": {"x": x, "y": y, "w": w, "h": h},
        "visible": 1, "zIndex": 10, "opacity": 1, "transform": 0,
        "borderWidth": 0, "BorderEdges": 0,
        "borderStyle": "solid", "borderColor": "#ccc",
        "diy": []
    }

# 构建组件列表
components = []

# 1. 顶部标题框 DvBorderBox1
c1 = {
    "type": "DvBorderBox1",
    "identifier": "hdr",
    "name": "hdr",
    "animate": default_animate(),
    "action": [], "active": [], "dataBind": [],
}
c1["style"] = default_style(0, 0, 1920, 55)
c1["style"]["diy"] = [
    {"name": "text", "type": 4, "value": "1A配电室 - 模拟监控总览 (76台设备)", "key": "text"},
    {"name": "foreColor", "type": 2, "value": "#00d4ff", "key": "foreColor"},
    {"name": "fontSize", "type": 1, "value": 24, "key": "fontSize"},
]
components.append(c1)

# 2. 左侧设备树 DeviceTree
c2 = {
    "type": "DeviceTree",
    "identifier": "tree",
    "name": "tree",
    "animate": default_animate(),
    "action": [], "active": [], "dataBind": [],
}
c2["style"] = default_style(10, 65, 240, 950)
components.append(c2)

# 3. 右侧告警列表 alarmList
c3 = {
    "type": "alarmList",
    "identifier": "alarms",
    "name": "alarms",
    "animate": default_animate(),
    "action": [], "active": [], "dataBind": [],
}
c3["style"] = default_style(1660, 65, 250, 950)
components.append(c3)

# 4. ViewSvgText - 设备实时数据标签（76个设备实例）
conn_src = sqlite3.connect(DB)
cur = conn_src.cursor()
cur.execute("""SELECT name FROM monitor_list 
    WHERE project_uuid='e308e378-ed94-1264-197b-33e535b812b8' 
    AND type != 0 AND sid > 0
    ORDER BY name""")
devices = cur.fetchall()
conn_src.close()

# 布局参数：中间区域从 x=270 开始，共 1380 宽度，每行放 4 个设备 label
col_w, row_h, left_margin = 340, 28, 270
cols_per_row = 4
for i, (name,) in enumerate(devices):
    row = i // cols_per_row
    col = i % cols_per_row
    x = left_margin + col * col_w + 10
    y = 70 + row * row_h
    
    c = {
        "type": "ViewSvgText",
        "identifier": f"d_{name}",
        "name": f"d_{name}",
        "animate": default_animate(),
        "action": [], "active": [], "dataBind": [],
    }
    c["style"] = default_style(x, y, 200, 18)
    c["style"]["diy"] = [
        {"name": "text", "type": 4, "value": f"{name}", "key": "text"},
        {"name": "foreColor", "type": 2, "value": "#aabbcc", "key": "foreColor"},
        {"name": "fontSize", "type": 1, "value": 14, "key": "fontSize"},
    ]
    components.append(c)

print(f"Built {len(components)} components ({len(devices)} devices + 3 containers)")

# 序列化并编码
comp_json = json.dumps(components, ensure_ascii=False)
comp_b64 = base64.b64encode(comp_json.encode('utf-8')).decode('utf-8')

# Layer 配置
layer = json.dumps({
    "backColor": "#0a1628",
    "backgroundImage": "",
    "widthHeightRatio": "",
    "width": 1920,
    "height": 1080,
    "autoSize": 0,
}, ensure_ascii=False)

# 写入数据库
conn = sqlite3.connect(DB)
cur = conn.cursor()

# 先更新 layer 数据
cur.execute("""
    UPDATE display_model_layer 
    SET layer = ?, components = ?, is_home = 1, page_type = 1
    WHERE model_id = ? AND page_name = '总览'
""", (layer, comp_b64, MODEL_UID))

if cur.rowcount == 0:
    # 没有总览页，插入
    import uuid
    cur.execute("""
        INSERT INTO display_model_layer (created_at, model_id, page_name, page_id, is_home, page_type, layer, components, is_login)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
    """, (
        datetime.now().isoformat(), MODEL_UID, '总览',
        str(uuid.uuid4()), 1, 1, layer, comp_b64, 0
    ))
    print("Inserted new layer row")
else:
    print(f"Updated layer row: {comp_b64[:40]}...")

conn.commit()
conn.close()

print(f"""
✅ 大屏组件已修复！数据量: {len(comp_json):,} bytes, base64: {len(comp_b64):,} chars
   Components: 1 DvBorderBox1 + 1 DeviceTree + 1 alarmList + {len(devices)} ViewSvgText
   = {len(components)} 个组件

⚠️  请重启后端（cd ism_server_user && ./ism_server），然后刷新页面
""")
