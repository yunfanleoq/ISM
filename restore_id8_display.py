#!/usr/bin/env python3
"""
Restore id=8 display_model_layer with working view-svg-text components.
Rebuilds the sci-fi display model with both text and DataV components.
"""
import pymysql
import json
import base64
import uuid

conn = pymysql.connect(
    host='127.0.0.1',
    port=2881,
    user='root@ism_tenant',
    password='ism2024!',
    database='ism',
    charset='utf8mb4'
)
cursor = conn.cursor()

MODEL_ID = '043135ad-44be-e5d8-89be-3e54883c23a8'

def uid():
    return uuid.uuid4().hex

def make_cell(shape, x, y, w, h, z_idx, **detail_kw):
    cid = uid()
    sty = detail_kw.get('style', {})
    sty['position'] = {"x": x, "y": y, "w": w, "h": h}
    return {
        "shape": shape,
        "id": cid,
        "x": x, "y": y, "width": w, "height": h, "zIndex": z_idx,
        "visible": True,
        "position": {"x": x, "y": y},
        "size": {"width": w, "height": h},
        "data": {
            "detail": {
                "type": shape,
                "identifier": cid,
                "name": detail_kw.get("name", shape),
                "style": sty,
                "animate": detail_kw.get("animate", {
                    "selected": [],
                    "animateElement": [],
                    "condition": {"deviceSN":"","selectVideoType":0,"isBandDevice":False,"bandType":1,"dataID":"","dataName":"","operator":"","OperatorValue":"","OperatorMaxValue":""},
                    "isExpression": False,
                    "animateList": [],
                    "move": {"x":{"deviceSN":"","selectVideoType":0,"isBandDevice":False,"bandType":1,"dataID":"","dataName":""},"y":{"deviceSN":"","selectVideoType":0,"isBandDevice":False,"bandType":1,"dataID":"","dataName":""}}
                }),
                "action": [],
                "active": [],
                "dataBind": []
            }
        }
    }

def text_style(text, font_size=14, fore_color="#ffffff"):
    return {
        "text": text, "fontSize": font_size, "foreColor": fore_color,
        "visible": 1, "borderWidth": 0, "BorderEdges": 0, "opacity": 1,
        "diy": []
    }

components = []
z = 0

# Title
components.append(make_cell("view-svg-text", 50, 20, 800, 60, z,
    name="航信机房", style=text_style("航信机房", 28, "#60a5fa")))
z += 1

# Subtitle  
components.append(make_cell("view-svg-text", 50, 75, 600, 30, z,
    name="NCC SCADA Power Monitor", style=text_style("NCC SCADA Power Monitor — 智能电力监控", 14, "#64748b")))
z += 1

# dv-border-box8 as outer border (test that the fix works)
components.append(make_cell("dv-border-box8", 40, 110, 1840, 920, z,
    name="Outer Border Box",
    style={"visible": 1, "diy": [
        {"name":"border89cur","type":1,"value":10,"min":1,"key":"border89cur"},
        {"name":"border89Direction","type":6,"value":1,"min":1,"key":"border89Direction",
         "enumList":[{"value":0,"option":"Forward"},{"value":1,"option":"Negative"}]}
    ]}))
z += 1

# dv-border-box1 left panel
components.append(make_cell("dv-border-box1", 80, 160, 880, 830, z,
    name="Left Panel"))
z += 1

# Panel title
components.append(make_cell("view-svg-text", 100, 170, 400, 30, z,
    name="Left Panel Title", style=text_style("实时电力趋势", 16, "#60a5fa")))
z += 1

# dv-border-box1 right panel
components.append(make_cell("dv-border-box1", 990, 160, 850, 830, z,
    name="Right Panel"))
z += 1

# Right panel title
components.append(make_cell("view-svg-text", 1010, 170, 300, 30, z,
    name="Right Panel Title", style=text_style("设备实时参数", 16, "#60a5fa")))
z += 1

# Stat boxes in left panel
stat_defs = [
    {"x": 100, "y": 220, "value": "12", "label": "在线设备", "color": "#34d399"},
    {"x": 340, "y": 220, "value": "486", "label": "总功率/kW", "color": "#f59e0b"},
    {"x": 580, "y": 220, "value": "2.8k", "label": "今日用电/kWh", "color": "#60a5fa"},
    {"x": 820, "y": 220, "value": "正常", "label": "系统状态", "color": "#10b981"},
]

for stat in stat_defs:
    sx, sy = stat['x'], stat['y']
    components.append(make_cell("dv-border-box1", sx, sy, 210, 100, z,
        name=f"Stat-{stat['label']}"))
    z += 1
    components.append(make_cell("view-svg-text", sx+10, sy+10, 190, 45, z,
        name=f"Value-{stat['label']}", style=text_style(stat['value'], 24, stat['color'])))
    z += 1
    components.append(make_cell("view-svg-text", sx+10, sy+55, 190, 30, z,
        name=f"Label-{stat['label']}", style=text_style(stat['label'], 12, "#94a3b8")))
    z += 1

# dv-decoration3 headers
components.append(make_cell("dv-decoration3", 50, 10, 200, 50, z,
    name="Header Decoration Left"))
z += 1
components.append(make_cell("dv-decoration3", 1670, 10, 200, 50, z,
    name="Header Decoration Right"))
z += 1

print(f"Total cells: {len(components)}")
print("Cell types:")
from collections import Counter
types = Counter(c['shape'] for c in components)
for t, n in sorted(types.items()):
    print(f"  {t}: {n}")

# CRITICAL: wrap in cells!
components_json = json.dumps({"cells": components}, ensure_ascii=False)
layer_json = json.dumps({"backColor": "#0a0e17", "width": 1920, "height": 1080}, ensure_ascii=False)

components_b64 = base64.b64encode(components_json.encode('utf-8')).decode('utf-8')
layer_b64 = base64.b64encode(layer_json.encode('utf-8')).decode('utf-8')

cursor.execute(
    "UPDATE display_model_layer SET components = %s WHERE model_id = %s AND page_name = 'main' AND deleted_at IS NULL",
    (components_b64, MODEL_ID)
)
print(f"Components UPDATE: {cursor.rowcount} rows")

cursor.execute(
    "UPDATE display_model_layer SET layer = %s WHERE model_id = %s AND page_name = 'main' AND deleted_at IS NULL",
    (layer_b64, MODEL_ID)
)
print(f"Layer UPDATE: {cursor.rowcount} rows")

conn.commit()
print("Restore complete! The display model has 3 dv-border-box1 + 1 dv-border-box8 + 2 dv-decoration3 + view-svg-text cells.")
cursor.close()
conn.close()
