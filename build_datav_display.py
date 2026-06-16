#!/usr/bin/env python3
"""
Build comprehensive sci-fi DataV display model for 航信机房 project.
Creates components JSON with ISMDisPlay DataV cells and updates display_model_layer.
"""

import pymysql
import json
import uuid
import base64
import sys

# ── Database connection ──────────────────────────────────────────────────
conn = pymysql.connect(
    host='127.0.0.1',
    port=2881,
    user='root@ism_tenant',
    password='ism2024!',
    database='ism',
    charset='utf8mb4'
)
cursor = conn.cursor(pymysql.cursors.DictCursor)

PROJECT_UUID = '31bc90be-ebc4-dd61-ba9d-ce6e075e40e2'
MODEL_ID = '043135ad-44be-e5d8-89be-3e54883c23a8'

def uid():
    """Generate a UUID v4 hex string without dashes."""
    return uuid.uuid4().hex


# ══════════════════════════════════════════════════════════════════════════
# STEP 1: Query device and data model points
# ══════════════════════════════════════════════════════════════════════════

print("=" * 60)
print("STEP 1: Querying device and data model points...")
print("=" * 60)

# 1a. Get the first A20电力仪表 device (type=1) — uses ml.muid as device model UUID
cursor.execute(
    "SELECT ml.uuid as dev_uuid, ml.name, ml.muid as model_uuid "
    "FROM monitor_list ml "
    "WHERE ml.project_uuid=%s AND ml.type=1 AND ml.deleted_at IS NULL LIMIT 1",
    (PROJECT_UUID,)
)
device = cursor.fetchone()
if not device:
    print("ERROR: No device found with type=1 in the project!")
    sys.exit(1)

dev_uuid = device['dev_uuid']
dev_name = device['name']
dev_model_uuid = device.get('model_uuid', '')
print(f"  Device found: {dev_name} (uuid={dev_uuid})")
print(f"  Device model UUID (monitor_list.muid): {dev_model_uuid or '(not set!)'}")

if not dev_model_uuid:
    print("ERROR: Could not determine device model UUID from monitor_list.muid!")
    sys.exit(1)

# 1b. Query ALL data model points for this device model (to find best matches)
cursor.execute(
    "SELECT dm.name, dm.uuid, dm.muid, dm.register_address "
    "FROM modbus_devices_data_model dm "
    "JOIN modbus_devices_register_group rg ON dm.register_group_uuid = rg.uuid "
    "WHERE rg.muid=%s AND dm.deleted_at IS NULL "
    "ORDER BY dm.register_address",
    (dev_model_uuid,)
)
all_points = cursor.fetchall()
print(f"\n  All data model points for this device ({len(all_points)} total):")
for dp in all_points:
    print(f"    - {dp['name']} (muid={dp['muid']})")

# 1c. Find the best matching points for display
# This device has: AB线电压, BC线电压, CA线电压 (no A相电压), A相电流, 总有功功率
phase_voltage = None
phase_current = None
total_power = None

# Priority order for voltage: A相电压 > AB线电压 > any 电压
# Priority order for current: A相电流 > any 电流
# Priority order for power: 总有功功率 > 总功率 > any 功率

for dp in all_points:
    name = dp['name']
    if phase_voltage is None:
        if 'A相电压' in name:
            phase_voltage = dp
        elif 'AB线电压' in name:
            phase_voltage = dp
        elif '电压' in name and phase_voltage is None:
            phase_voltage = dp
    if phase_current is None:
        if 'A相电流' in name:
            phase_current = dp
        elif '电流' in name and '谐波' not in name:
            phase_current = dp
    if total_power is None:
        if '总有功功率' in name:
            total_power = dp
        elif '有功功' in name:
            total_power = dp
        elif '总功率' in name:
            total_power = dp

print(f"\n  Mapped data points for display:")
print(f"    Voltage:  {phase_voltage['name'] if phase_voltage else 'NOT FOUND'} (muid={phase_voltage['muid']})" if phase_voltage else "    Voltage:  NOT FOUND!")
print(f"    Current:  {phase_current['name'] if phase_current else 'NOT FOUND'} (muid={phase_current['muid']})" if phase_current else "    Current:  NOT FOUND!")
print(f"    Power:    {total_power['name'] if total_power else 'NOT FOUND'} (muid={total_power['muid']})" if total_power else "    Power:    NOT FOUND!")


# ══════════════════════════════════════════════════════════════════════════
# STEP 2: Build the components JSON
# ══════════════════════════════════════════════════════════════════════════

print("\n" + "=" * 60)
print("STEP 2: Building components JSON...")
print("=" * 60)

components = []
z = 0  # zIndex counter


# ── Helper: make animate stub ────────────────────────────────────────────
def make_animate(anim_type='selected', val=None):
    anim = {
        "selected": [] if anim_type == 'selected' else [],
        "animateElement": [],
        "move": {
            "x": {
                "deviceSN": "",
                "selectVideoType": 0,
                "isBandDevice": False,
                "bandType": 1,
                "dataID": "",
                "dataName": ""
            },
            "y": {
                "deviceSN": "",
                "selectVideoType": 0,
                "isBandDevice": False,
                "bandType": 1,
                "dataID": "",
                "dataName": ""
            }
        }
    }
    return anim


# ── Helper: make cell base ───────────────────────────────────────────────
def make_cell(shape, x, y, w, h, z_index, data_detail):
    cell_id = uid()
    return {
        "shape": shape,
        "id": cell_id,
        "x": x,
        "y": y,
        "width": w,
        "height": h,
        "zIndex": z_index,
        "visible": True,
        "position": {"x": x, "y": y},
        "size": {"width": w, "height": h},
        "data": {
            "detail": {
                "type": shape,
                "identifier": cell_id,
                "name": data_detail.get("name", shape),
                "style": data_detail.get("style", {}),
                "animate": data_detail.get("animate", make_animate()),
                "active": data_detail.get("active", []),
                "dataBind": data_detail.get("dataBind", [])
            }
        }
    }


# ── Helper: make dv-style (position-based) ───────────────────────────────
def dv_style(x, y, w, h, **extra):
    s = {"position": {"x": x, "y": y, "w": w, "h": h}}
    s.update(extra)
    return s


# ── Helper: make text-style ──────────────────────────────────────────────
def text_style(text, font_size=14, fore_color="#ffffff", **extra):
    s = {
        "text": text,
        "fontSize": font_size,
        "foreColor": fore_color,
        "visible": 1,
        "borderWidth": 0,
        "BorderEdges": 0,
        "opacity": 1,
        "diy": []
    }
    s.update(extra)
    return s


# ── 1️⃣ dv-decoration3 × 2 (header decorations) ──────────────────────────
components.append(make_cell("dv-decoration3", 50, 10, 200, 60, z, {
    "style": dv_style(50, 10, 200, 60),
    "name": "Header Decoration Left"
}))
z += 1

components.append(make_cell("dv-decoration3", 1670, 10, 200, 60, z, {
    "style": dv_style(1670, 10, 200, 60),
    "name": "Header Decoration Right"
}))
z += 1

# ── 2️⃣ Header title ─────────────────────────────────────────────────────
components.append(make_cell("view-svg-text", 80, 20, 800, 50, z, {
    "style": text_style("航信机房电力监控系统", font_size=28, fore_color="#60a5fa"),
    "name": "Header Title"
}))
z += 1

# ── 3️⃣ Subtitle ─────────────────────────────────────────────────────────
components.append(make_cell("view-svg-text", 80, 65, 600, 30, z, {
    "style": text_style("SCADA Power Monitor — 智能电力监控大屏", font_size=14, fore_color="#64748b"),
    "name": "Header Subtitle"
}))
z += 1

# ── 4️⃣ dv-border-box8 outer container ───────────────────────────────────
components.append(make_cell("dv-border-box8", 40, 100, 1840, 940, z, {
    "style": {
        **dv_style(40, 100, 1840, 940),
        "dur": 3,
        "reverse": False
    },
    "animate": {**make_animate(), "dur": 3, "reverse": False},
    "name": "Outer Border Box"
}))
z += 1

# ── 5️⃣ Stats row — 4 stat panels ────────────────────────────────────────
stat_defs = [
    {"x": 80,  "y": 140, "value": "12",  "label": "在线设备", "color": "#34d399"},
    {"x": 310, "y": 140, "value": "486", "label": "总功率/kW",  "color": "#f59e0b"},
    {"x": 540, "y": 140, "value": "2.8k", "label": "今日用电/kWh", "color": "#60a5fa"},
    {"x": 770, "y": 140, "value": "正常", "label": "系统状态",    "color": "#10b981"},
]

for stat in stat_defs:
    sx, sy = stat['x'], stat['y']
    # dv-border-box1 container
    components.append(make_cell("dv-border-box1", sx, sy, 210, 120, z, {
        "style": dv_style(sx, sy, 210, 120),
        "name": f"Stat Panel - {stat['label']}"
    }))
    z += 1

    # dv-decoration5 decorative background (opacity 0.1)
    components.append(make_cell("dv-decoration5", sx + 5, sy + 5, 200, 50, z, {
        "style": dv_style(sx + 5, sy + 5, 200, 50),
        "name": f"Stat Decoration - {stat['label']}"
    }))
    z += 1

    # Value text
    components.append(make_cell("view-svg-text", sx + 15, sy + 15, 180, 45, z, {
        "style": text_style(stat['value'], font_size=32, fore_color=stat['color']),
        "name": f"Stat Value - {stat['label']}"
    }))
    z += 1

    # Label text
    components.append(make_cell("view-svg-text", sx + 15, sy + 65, 180, 40, z, {
        "style": text_style(stat['label'], font_size=13, fore_color="#94a3b8"),
        "name": f"Stat Label - {stat['label']}"
    }))
    z += 1

# ── 6️⃣ Left main panel — dv-border-box1 + smooth chart ─────────────────
components.append(make_cell("dv-border-box1", 80, 300, 900, 650, z, {
    "style": dv_style(80, 300, 900, 650),
    "name": "Left Panel - Power Trend"
}))
z += 1

# Panel title
components.append(make_cell("view-svg-text", 100, 310, 300, 30, z, {
    "style": text_style("实时电力趋势", font_size=16, fore_color="#60a5fa"),
    "name": "Left Panel Title"
}))
z += 1

# Build chart active binding using the phase_voltage data point
chart_active = []
if phase_voltage:
    chart_active = [{
        "condition": {
            "deviceSN": dev_uuid,
            "dataID": phase_voltage['muid'],
            "isBandDevice": True,
            "bandType": 1,
            "selectVideoType": 0
        }
    }]

# Smooth chart for power trend
components.append(make_cell("view-real-data-smooth-chart", 100, 330, 860, 600, z, {
    "style": dv_style(100, 330, 860, 600),
    "animate": make_animate(),
    "active": chart_active,
    "dataBind": [],
    "name": "Power Trend Chart"
}))
z += 1

# ── 7️⃣ Right panel — dv-border-box1 + 2 gauges + labels ────────────────
components.append(make_cell("dv-border-box1", 1010, 300, 830, 650, z, {
    "style": dv_style(1010, 300, 830, 650),
    "name": "Right Panel - Device Gauges"
}))
z += 1

# Panel title
components.append(make_cell("view-svg-text", 1030, 310, 300, 30, z, {
    "style": text_style("设备实时参数", font_size=16, fore_color="#60a5fa"),
    "name": "Right Panel Title"
}))
z += 1

# Gauge 1: Voltage (may be AB线电压 or A相电压)
voltage_name = phase_voltage['name'] if phase_voltage else "电压"
current_name = phase_current['name'] if phase_current else "电流"
power_name = total_power['name'] if total_power else "功率"

gauge1_active = []
if phase_voltage:
    gauge1_active = [{
        "condition": {
            "deviceSN": dev_uuid,
            "dataID": phase_voltage['muid'],
            "isBandDevice": True,
            "bandType": 1,
            "selectVideoType": 0
        }
    }]

# Label above gauge 1
components.append(make_cell("view-svg-text", 1040, 320, 350, 30, z, {
    "style": text_style(f"{voltage_name} (V)", font_size=14, fore_color="#cbd5e1"),
    "name": "Gauge 1 Label"
}))
z += 1

components.append(make_cell("view-chart-gauge", 1040, 350, 350, 280, z, {
    "style": dv_style(1040, 350, 350, 280),
    "animate": make_animate(),
    "active": gauge1_active,
    "dataBind": [],
    "name": f"{voltage_name} Gauge"
}))
z += 1

# Gauge 2: Current
gauge2_active = []
if phase_current:
    gauge2_active = [{
        "condition": {
            "deviceSN": dev_uuid,
            "dataID": phase_current['muid'],
            "isBandDevice": True,
            "bandType": 1,
            "selectVideoType": 0
        }
    }]

# Label above gauge 2
components.append(make_cell("view-svg-text", 1420, 320, 350, 30, z, {
    "style": text_style(f"{current_name} (A)", font_size=14, fore_color="#cbd5e1"),
    "name": "Gauge 2 Label"
}))
z += 1

components.append(make_cell("view-chart-gauge", 1420, 350, 350, 280, z, {
    "style": dv_style(1420, 350, 350, 280),
    "animate": make_animate(),
    "active": gauge2_active,
    "dataBind": [],
    "name": f"{current_name} Gauge"
}))
z += 1

# ── 8️⃣ Device stats text in right panel (static labels) ────────────────
stat_texts = [
    {"y": 660, "text": f"{voltage_name}: -- V", "color": "#34d399"},
    {"y": 710, "text": f"{current_name}: -- A", "color": "#f59e0b"},
    {"y": 760, "text": f"{power_name}: -- kW", "color": "#60a5fa"},
]

for st in stat_texts:
    components.append(make_cell("view-svg-text", 1040, st['y'], 350, 40, z, {
        "style": text_style(st['text'], font_size=18, fore_color=st['color']),
        "name": f"Device Stat - {st['text']}"
    }))
    z += 1

# ── 9️⃣ Device info text at bottom ───────────────────────────────────────
components.append(make_cell("view-svg-text", 1040, 820, 400, 35, z, {
    "style": text_style(f"设备名称: {dev_name}", font_size=14, fore_color="#64748b"),
    "name": "Device Name Label"
}))
z += 1

components.append(make_cell("view-svg-text", 1040, 860, 400, 35, z, {
    "style": text_style(f"设备编号: {dev_uuid[:16]}...", font_size=12, fore_color="#475569"),
    "name": "Device UUID Label"
}))
z += 1


# ── Build final JSON ─────────────────────────────────────────────────────
layer_json = json.dumps({
    "backColor": "#0a0e17",
    "width": 1920,
    "height": 1080
}, ensure_ascii=False)

components_json = json.dumps({"cells": components}, ensure_ascii=False)

print(f"\n  Total cells created: {len(components)}")
print(f"  Max zIndex: {z - 1}")

# Count by type
shape_counts = {}
for c in components:
    s = c['shape']
    shape_counts[s] = shape_counts.get(s, 0) + 1
print("\n  Cell types breakdown:")
for shape, count in sorted(shape_counts.items()):
    print(f"    {shape}: {count}")


# ══════════════════════════════════════════════════════════════════════════
# STEP 3: Base64 encode and update database
# ══════════════════════════════════════════════════════════════════════════

print("\n" + "=" * 60)
print("STEP 3: Base64 encoding and updating database...")
print("=" * 60)

components_b64 = base64.b64encode(components_json.encode('utf-8')).decode('utf-8')
layer_b64 = base64.b64encode(layer_json.encode('utf-8')).decode('utf-8')

print(f"  Components JSON size: {len(components_json)} chars")
print(f"  Components base64 size: {len(components_b64)} chars")
print(f"  Layer base64 size: {len(layer_b64)} chars")

# Update components field
sql_components = (
    "UPDATE display_model_layer "
    "SET components = %s "
    "WHERE model_id = %s AND page_name = 'main' AND deleted_at IS NULL"
)
cursor.execute(sql_components, (components_b64, MODEL_ID))
affected_components = cursor.rowcount
print(f"\n  Components UPDATE affected rows: {affected_components}")

# Update layer field
sql_layer = (
    "UPDATE display_model_layer "
    "SET layer = %s "
    "WHERE model_id = %s AND page_name = 'main' AND deleted_at IS NULL"
)
cursor.execute(sql_layer, (layer_b64, MODEL_ID))
affected_layer = cursor.rowcount
print(f"  Layer UPDATE affected rows: {affected_layer}")

conn.commit()
print("\n  Database commit successful!")


# ══════════════════════════════════════════════════════════════════════════
# STEP 4: Summary
# ══════════════════════════════════════════════════════════════════════════

print("\n" + "=" * 60)
print("BUILD SUMMARY")
print("=" * 60)
print(f"""
  Project:    航信机房 (Power Monitor)
  Model ID:   {MODEL_ID}
  Device:     {dev_name} ({dev_uuid})
  
  Canvas:     1920 × 1080, Background #0a0e17
  Total Cells: {len(components)}
  
  Layer structure:
    - dv-decoration3 × 2       (header decoration)
    - view-svg-text × 2        (title + subtitle)
    - dv-border-box8 × 1       (outer container)
    - dv-border-box1 × 4       (stat panel containers)
    - dv-decoration5 × 4       (stat backgrounds)
    - view-svg-text × 8        (stat values + labels)
    - dv-border-box1 × 1       (left panel)
    - view-svg-text × 2        (panel titles)
    - view-real-data-smooth-chart × 1  (power trend)
    - dv-border-box1 × 1       (right panel)
    - view-chart-gauge × 2     (voltage + current gauges)
    - view-svg-text × 5        (device stat labels + device info)

  Data bindings:
    {phase_voltage and f"- {voltage_name} → smooth chart + gauge 1 (muid: {phase_voltage['muid']})" or '- Voltage → NOT FOUND'}
    {phase_current and f"- {current_name} → gauge 2 (muid: {phase_current['muid']})" or '- Current → NOT FOUND'}
    {total_power and f"- {power_name} → referenced (muid: {total_power['muid']})" or '- Power → NOT FOUND'}

  Database updates:
    - display_model_layer.components = base64(components JSON)
    - display_model_layer.layer      = base64(layer JSON)
    
  Done! The display model is ready for rendering.
""")

# Close connection
cursor.close()
conn.close()
