#!/usr/bin/env python3
"""
Test: Insert single dv-border-box8 cell into display_model_layer id=8.
Verifies the cells-wrapping fix.
"""
import pymysql
import json
import uuid
import base64

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
cell_id = uuid.uuid4().hex

# Build a single dv-border-box8 cell with complete properties
cell = {
    "shape": "dv-border-box8",
    "id": cell_id,
    "x": 100,
    "y": 100,
    "width": 400,
    "height": 300,
    "zIndex": -1,
    "visible": True,
    "position": {"x": 100, "y": 100},
    "size": {"width": 400, "height": 300},
    "data": {
        "detail": {
            "type": "dv-border-box8",
            "identifier": cell_id,
            "name": "Test Border Box 8",
            "style": {
                "position": {"x": 100, "y": 100, "w": 400, "h": 300},
                "visible": 1,
                "backColor": "transparent",
                "zIndex": -1,
                "transform": 0,
                "diy": [
                    {"name": "configComponent.bigScreen.border.border89cur", "type": 1, "value": 10, "min": 1, "key": "border89cur"},
                    {"name": "configComponent.bigScreen.border.border89Direction", "type": 6, "value": 1, "min": 1, "key": "border89Direction",
                     "enumList": [
                         {"value": 0, "option": "configComponent.bigScreen.border.border89DirectionForward"},
                         {"value": 1, "option": "configComponent.bigScreen.border.border89DirectionNegative"}
                     ]}
                ]
            },
            "animate": {
                "selected": [],
                "animateElement": [],
                "condition": {
                    "deviceSN": "",
                    "selectVideoType": 0,
                    "isBandDevice": False,
                    "bandType": 1,
                    "dataID": "",
                    "dataName": "",
                    "operator": "",
                    "OperatorValue": "",
                    "OperatorMaxValue": ""
                },
                "isExpression": False,
                "animateList": [
                    {"id": "blink", "name": "component.public.animateBlink"}
                ],
                "move": {
                    "x": {"deviceSN": "", "selectVideoType": 0, "isBandDevice": False, "bandType": 1, "dataID": "", "dataName": ""},
                    "y": {"deviceSN": "", "selectVideoType": 0, "isBandDevice": False, "bandType": 1, "dataID": "", "dataName": ""}
                }
            },
            "action": [],
            "active": [],
            "dataBind": []
        }
    }
}

# CRITICAL: wrap in cells
components_json = json.dumps({"cells": [cell]}, ensure_ascii=False)
layer_json = json.dumps({
    "backColor": "#0a0e17",
    "width": 1920,
    "height": 1080
}, ensure_ascii=False)

components_b64 = base64.b64encode(components_json.encode('utf-8')).decode('utf-8')
layer_b64 = base64.b64encode(layer_json.encode('utf-8')).decode('utf-8')

print(f"Test cell id: {cell_id}")
print(f"Components JSON first 200 chars: {components_json[:200]}")

# Update
cursor.execute(
    "UPDATE display_model_layer SET components = %s WHERE model_id = %s AND page_name = 'main' AND deleted_at IS NULL",
    (components_b64, MODEL_ID)
)
print(f"Components UPDATE affected: {cursor.rowcount}")

cursor.execute(
    "UPDATE display_model_layer SET layer = %s WHERE model_id = %s AND page_name = 'main' AND deleted_at IS NULL",
    (layer_b64, MODEL_ID)
)
print(f"Layer UPDATE affected: {cursor.rowcount}")

conn.commit()
print("Done! Refresh the browser at http://localhost:7080/#/AppRun/043135ad-44be-e5d8-89be-3e54883c23a8")

cursor.close()
conn.close()
