#!/usr/bin/env python3
"""
构建航信机房炫酷科技感大屏 v3 — 左侧导航 + 多层级钻探
Canvas: 1920×1080, bg #0a0e17

4 层页面:
  Level 0: overview (id=8)    — 280px 左侧导航树 + 主大屏面板
  Level 1: building-detail     — 建筑名标题 + 楼层设备组卡片网格
  Level 2: floor-detail         — 楼层名标题 + 设备列表表格
  Level 3: device-detail (id=10)— 91 cells 设备参数 + 趋势图 + 状态监控

布局网格 (overview):
  Left Sidebar:    x:0-280    w:280  (dv-border-box8 + 面包屑 + 设备树)
  Breadcrumb:      x:10, y:5, w:260
  Tree Nav:         inside sidebar
  Header:          x:290, y:0-80, w:1630
  Stats Cards:     x:290, 690, 1090, 1490, w:390 each, y:100-210
  Left Panel:      x:290, y:230, w:780, h:400
  Right Upper:     x:1090, y:230, w:400, h:400
  Right Lower:     x:1090, y:650, w:400, h:220
  Bottom Panel:    x:290, y:650, w:780, h:220
"""
import pymysql
import json
import base64
import uuid as _uuid
from collections import defaultdict

# ── DB ──────────────────────────────────────────────
conn = pymysql.connect(
    host='127.0.0.1', port=2881,
    user='root@ism_tenant', password='ism2024!',
    database='ism'
)
cur = conn.cursor()

# ── Constants ───────────────────────────────────────
MODEL_ID = '043135ad-44be-e5d8-89be-3e54883c23a8'
PAGE_ID_MAIN = MODEL_ID
PAGE_ID_DEVICE = _uuid.uuid5(_uuid.NAMESPACE_DNS, 'ncc-dash-device-detail').hex
PAGE_ID_BUILDING = _uuid.uuid5(_uuid.NAMESPACE_DNS, 'ncc-dash-building-detail').hex
PAGE_ID_FLOOR = _uuid.uuid5(_uuid.NAMESPACE_DNS, 'ncc-dash-floor-detail').hex
DEVICE_UUID = '68db26b1-113d-ad7e-79ff-10dbcc1c18d2'
DEVICE_NAME = '1A1_U11_S18_1'
DEV_MODEL_UUID = '3d734984-56f6-5494-ad4c-dfc67ca28ac8'
PROJECT_UUID = '31bc90be-ebc4-dd61-ba9d-ce6e075e40e2'

# Fetch data points
cur.execute(
    'SELECT name, uuid, data_unit FROM modbus_devices_data_model WHERE muid=%s ORDER BY id',
    (DEV_MODEL_UUID,)
)
dp_rows = cur.fetchall()
DP_MAP = {r[0]: {'uuid': r[1], 'unit': r[2] or ''} for r in dp_rows}

cur.execute('SELECT COUNT(*) FROM monitor_list')
TOTAL_DEVICES = cur.fetchone()[0]

print(f"Total devices: {TOTAL_DEVICES}")
print(f"Data points ({len(DP_MAP)}): {list(DP_MAP.keys())}")

# ── Query real device hierarchy ─────────────────────
cur.execute("""
    SELECT uuid, name, sid, pid, type, muid, status
    FROM monitor_list
    WHERE project_uuid = %s AND deleted_at IS NULL
    ORDER BY pid, type, name
""", (PROJECT_UUID,))
all_devices = cur.fetchall()

# Build hierarchy: pid → list of children
children_by_pid = defaultdict(list)
device_by_sid = {}
for row in all_devices:
    uuid, name, sid, pid, dtype, muid, status = row
    children_by_pid[pid].append({
        'uuid': uuid, 'name': name, 'sid': sid,
        'type': dtype, 'muid': muid, 'status': status
    })
    device_by_sid[sid] = {
        'uuid': uuid, 'name': name, 'sid': sid,
        'type': dtype, 'muid': muid, 'status': status
    }

# Extract buildings: type=0 nodes that have type=1 children
buildings = []
for node in all_devices:
    _, name, sid, pid, dtype, muid, status = node
    if dtype == 0:
        children = children_by_pid.get(sid, [])
        type1_children = [c for c in children if c['type'] == 1]
        if type1_children:
            # Group devices by prefix (for "floor" grouping)
            floors = defaultdict(list)
            for d in type1_children:
                # Extract group key from name (e.g., "1A1_U11_S18_1" → "S18")
                parts = d['name'].split('_')
                if len(parts) >= 3:
                    floor_key = parts[2]  # e.g., "S18"
                else:
                    floor_key = 'default'
                floors[floor_key].append(d)
            building_entry = {
                'uuid': uuid, 'name': name, 'sid': sid,
                'devices': type1_children,
                'device_count': len(type1_children),
                'floors': [{'key': k, 'name': f'{k}设备组', 'devices': v, 'count': len(v)}
                           for k, v in sorted(floors.items())]
            }
            buildings.append(building_entry)

# Also look for second-level type=0 nodes (cabinet groups)
# e.g., 1A1_U11柜 → group by parent sid
for building in buildings[:]:
    children = children_by_pid.get(building['sid'], [])
    type0_children = [c for c in children if c['type'] == 0]
    for c0 in type0_children:
        grand_children = children_by_pid.get(c0['sid'], [])
        type1_gc = [g for g in grand_children if g['type'] == 1]
        if type1_gc:
            floor_groups = defaultdict(list)
            for d in type1_gc:
                parts = d['name'].split('_')
                fk = parts[2] if len(parts) >= 3 else 'default'
                floor_groups[fk].append(d)
            sub = {
                'uuid': c0['uuid'], 'name': c0['name'], 'sid': c0['sid'],
                'devices': type1_gc, 'device_count': len(type1_gc),
                'floors': [{'key': k, 'name': f'{k}设备组', 'devices': v, 'count': len(v)}
                           for k, v in sorted(floor_groups.items())]
            }
            buildings.append(sub)

# Remove root-level buildings that have no direct devices (they're just containers)
buildings = [b for b in buildings if b['device_count'] > 0]

print(f"\n=== Device Hierarchy ===")
for b in buildings:
    print(f"  Building: {b['name']} ({b['device_count']} devices)")
    for f in b['floors']:
        print(f"    Floor: {f['name']} ({f['count']} devices)")
        for d in f['devices'][:3]:
            status_str = '运行' if d['status']==1 else '离线'
            print(f"      - {d['name']} [{status_str}]")
        if f['count'] > 3:
            print(f"      ... +{f['count']-3} more")

# ── Helpers ─────────────────────────────────────────

def gen_uid(seed):
    return _uuid.uuid5(_uuid.NAMESPACE_DNS, f'ncc-v3-{seed}').hex

def _base_animate():
    return {
        "selected": [],
        "animateElement": [],
        "condition": {
            "deviceSN": "", "selectVideoType": 0, "isBandDevice": False,
            "bandType": 1, "dataID": "", "dataName": "",
            "operator": "", "OperatorValue": "", "OperatorMaxValue": ""
        },
        "isExpression": False, "animateList": [],
        "move": {
            "x": {"deviceSN": "", "selectVideoType": 0, "isBandDevice": False,
                  "bandType": 1, "dataID": "", "dataName": ""},
            "y": {"deviceSN": "", "selectVideoType": 0, "isBandDevice": False,
                  "bandType": 1, "dataID": "", "dataName": ""}
        }
    }

def _make_active(dp_name):
    if dp_name not in DP_MAP:
        return []
    dp = DP_MAP[dp_name]
    return [{
        "id": "ShowData",
        "name": "configComponent.variable.ShowData",
        "result": "",
        "isExpression": False,
        "condition": {
            "deviceSN": DEVICE_UUID,
            "DeviceName": DEVICE_NAME,
            "selectVideoType": 0,
            "isBandDevice": False,
            "bandType": 1,
            "dataID": dp['uuid'],
            "dataName": dp_name,
            "operator": "",
            "OperatorValue": "",
            "OperatorMaxValue": ""
        }
    }]

def _make_style(pos, **extras):
    s = {"position": pos, "visible": 1, "opacity": 1, "diy": []}
    s.update(extras)
    return s

# ── Typography & grid layout (SCADAMonitor-aligned) ──
FONT_TITLE = 20
FONT_SUBTITLE = 12
FONT_KPI_VALUE = 28
FONT_KPI_LABEL = 13
FONT_PANEL = 13
FONT_PARAM_VAL = 16
FONT_PARAM_LABEL = 11
FONT_NAV_BLDG = 13
FONT_NAV_FLR = 12
FONT_BREAD = 14

HEADER_H = 56
SIDEBAR_W = 280
SIDEBAR_X = 0
MAIN_X = SIDEBAR_W + 16          # 296, matches scada-main padding
MAIN_W = 1920 - MAIN_X - 16      # 1608
BODY_Y = HEADER_H                # 56

BLDG_ROW_H = 34
BLDG_ROW_GAP = 2
FLR_ROW_H = 28
FLR_ROW_GAP = 2
FLR_INDENT = 20

# SCADAMonitor palette
C_BG = '#0a0e17'
C_SIDEBAR = '#0d1220'
C_HEADER = '#141c2b'
C_BORDER = '#1e293b'
C_TEXT = '#e2e8f0'
C_TEXT_MUTED = '#94a3b8'
C_TEXT_DIM = '#64748b'
C_ACCENT = '#60a5fa'
C_GREEN = '#22c55e'


def kpi_val_font(card_w):
    """KPI 数值 28-36px，按卡片宽度缩放"""
    if card_w >= 380:
        return 32
    if card_w >= 320:
        return 28
    return 24


def rects_overlap(a, b):
    return not (a['x'] + a['w'] <= b['x'] or b['x'] + b['w'] <= a['x']
                or a['y'] + a['h'] <= b['y'] or b['y'] + b['h'] <= a['y'])


def cell_rect(c, pad=0):
    fs = c['data']['detail']['style'].get('fontSize', 14)
    vpad = max(pad, int(fs * 0.2))
    return {
        'text': c['data']['detail']['style'].get('text', '')[:30],
        'x': c['x'], 'y': c['y'] - vpad,
        'w': c['width'], 'h': c['height'] + vpad * 2,
    }


def find_text_overlaps(cells, pad=2):
    texts = [c for c in cells if c.get('shape') == 'view-svg-text']
    rects = [cell_rect(c, pad) for c in texts]
    pairs = []
    for i in range(len(rects)):
        for j in range(i + 1, len(rects)):
            if rects_overlap(rects[i], rects[j]):
                pairs.append((rects[i], rects[j]))
    return pairs


def vcenter(row_y, row_h, text_h):
    """行内垂直居中"""
    return row_y + max(0, (row_h - text_h) // 2)

# ──────────────────────────────────────────────────────
# CELL BUILDERS
# ──────────────────────────────────────────────────────

def make_text(seed, x, y, w, h, text, color='#c8d6e5', font_size=14,
              z=10, data_bound=False, dp_name=None, action=None):
    cell_id = gen_uid(seed)
    active = _make_active(dp_name) if (data_bound and dp_name) else []
    return {
        "shape": "view-svg-text",
        "id": cell_id, "x": x, "y": y, "width": w, "height": h,
        "zIndex": z, "visible": True,
        "position": {"x": x, "y": y},
        "size": {"width": w, "height": h},
        "data": {
            "detail": {
                "type": "view-svg-text",
                "identifier": cell_id,
                "name": text,
                "style": _make_style({"x": x, "y": y, "w": w, "h": h},
                                     text=text, fontSize=font_size,
                                     foreColor=color, borderWidth=0, BorderEdges=0),
                "animate": _base_animate(),
                "action": action or [],
                "active": active,
                "dataBind": []
            }
        }
    }

def make_border_box1(seed, x, y, w, h, z=1, action=None):
    cell_id = gen_uid(seed)
    return {
        "shape": "dv-border-box1",
        "id": cell_id, "x": x, "y": y, "width": w, "height": h,
        "zIndex": z, "visible": True,
        "position": {"x": x, "y": y},
        "size": {"width": w, "height": h},
        "data": {
            "detail": {
                "type": "dv-border-box1",
                "identifier": cell_id,
                "name": seed,
                "style": {"position": {"x": x, "y": y, "w": w, "h": h}},
                "animate": _base_animate(),
                "action": action or [], "active": [], "dataBind": []
            }
        }
    }

def make_decoration3(seed, x, y, w, h, z=1):
    cell_id = gen_uid(seed)
    return {
        "shape": "dv-decoration3",
        "id": cell_id, "x": x, "y": y, "width": w, "height": h,
        "zIndex": z, "visible": True,
        "position": {"x": x, "y": y},
        "size": {"width": w, "height": h},
        "data": {
            "detail": {
                "type": "dv-decoration3",
                "identifier": cell_id,
                "name": seed,
                "style": {"position": {"x": x, "y": y, "w": w, "h": h}},
                "animate": _base_animate(),
                "action": [], "active": [], "dataBind": []
            }
        }
    }

def make_decoration5(seed, x, y, w, h, z=1):
    cell_id = gen_uid(seed)
    return {
        "shape": "dv-decoration5",
        "id": cell_id, "x": x, "y": y, "width": w, "height": h,
        "zIndex": z, "visible": True,
        "position": {"x": x, "y": y},
        "size": {"width": w, "height": h},
        "data": {
            "detail": {
                "type": "dv-decoration5",
                "identifier": cell_id,
                "name": seed,
                "style": {"position": {"x": x, "y": y, "w": w, "h": h}},
                "animate": _base_animate(),
                "action": [], "active": [], "dataBind": []
            }
        }
    }

def make_border_box8(seed, x, y, w, h, z=0):
    cell_id = gen_uid(seed)
    return {
        "shape": "dv-border-box8",
        "id": cell_id, "x": x, "y": y, "width": w, "height": h,
        "zIndex": z, "visible": True,
        "position": {"x": x, "y": y},
        "size": {"width": w, "height": h},
        "data": {
            "detail": {
                "type": "dv-border-box8",
                "identifier": cell_id,
                "name": seed,
                "style": {
                    "visible": 1,
                    "diy": [
                        {"name": "border89cur", "type": 1, "value": 10, "min": 1, "key": "border89cur"},
                        {"name": "border89Direction", "type": 6, "value": 1, "min": 1, "key": "border89Direction",
                         "enumList": [{"value": 0, "option": "Forward"}, {"value": 1, "option": "Negative"}]}
                    ],
                    "position": {"x": x, "y": y, "w": w, "h": h}
                },
                "animate": _base_animate(),
                "action": [], "active": [], "dataBind": []
            }
        }
    }

def make_smooth_chart(seed, x, y, w, h, title, dp_names, z=5):
    cell_id = gen_uid(seed)
    active = []
    var_ids = ['ShowChartVariable1', 'ShowChartVariable2', 'ShowChartVariable3',
               'ShowChartVariable4', 'ShowChartVariable5']
    for i, dpn in enumerate(dp_names[:5]):
        dp = DP_MAP[dpn]
        active.append({
            "id": var_ids[i],
            "name": "configComponent.variable.ShowData",
            "result": "",
            "isExpression": False,
            "condition": {
                "deviceSN": DEVICE_UUID, "DeviceName": DEVICE_NAME,
                "selectVideoType": 0, "isBandDevice": False, "bandType": 1,
                "dataID": dp['uuid'], "dataName": dpn,
                "operator": "", "OperatorValue": "", "OperatorMaxValue": ""
            }
        })
    return {
        "shape": "ism-view-real-data-smooth-chart",
        "id": cell_id, "x": x, "y": y, "width": w, "height": h,
        "zIndex": z, "visible": True,
        "position": {"x": x, "y": y},
        "size": {"width": w, "height": h},
        "data": {
            "detail": {
                "type": "ism-view-real-data-smooth-chart",
                "identifier": cell_id, "name": title,
                "style": {
                    "position": {"x": x, "y": y, "w": w, "h": h},
                    "backColor": "transparent", "foreColor": "#ffffff",
                    "fontSize": 14, "fontFamily": "Arial", "zIndex": 1, "transform": 0,
                    "diy": [
                        {"name": "configComponent.ChartPublic.ChartTitle", "type": 4, "value": title, "key": "ChartTitle"},
                        {"name": "configComponent.ChartPublic.TimelyInitEcharts", "type": 1, "value": 60, "key": "TimelyInitEcharts"},
                        {"name": "configComponent.ChartPublic.ChartTimelyRefresh", "type": 1, "value": 60, "key": "ChartTimelyRefresh"},
                        {"name": "configComponent.ChartPublic.YMax", "type": 1, "value": 0, "key": "YMax"},
                        {"name": "configComponent.ChartPublic.YMin", "type": 1, "value": 0, "key": "YMin"},
                        {"name": "configComponent.ChartPublic.EchartsWidth", "type": 1, "value": 2, "key": "EchartsWidth"},
                        {"name": "configComponent.ChartPublic.EchartsXRotate", "type": 1, "value": 30, "key": "EchartsXRotate"},
                        {"name": "configComponent.ChartPublic.EchartsXFormat", "type": 4, "value": "HH:mm:ss", "key": "EchartsXFormat"},
                        {"name": "configComponent.ChartPublic.EchartsXTheme", "type": 6, "value": "dark", "key": "EchartsXTheme",
                         "enumList": [{"value": v, "option": v} for v in
                                      ["chalk","essos","dark","infographic","macarons","roma","shine","vintage",
                                       "purplePassion","walden","westeros","wonderland"]]}
                    ]
                },
                "animate": {
                    "selected": [],
                    "condition": {
                        "deviceSN": "", "selectVideoType": 0, "isBandDevice": False,
                        "bandType": 1, "dataID": "", "dataName": "",
                        "operator": "", "OperatorValue": "", "OperatorMaxValue": ""
                    },
                    "isExpression": False,
                    "animateList": [
                        {"id": "Forbidden", "name": "component.public.Forbidden"},
                        {"id": "blink", "name": "component.public.animateBlink"},
                        {"id": "Zoom", "name": "component.public.Zoom"},
                        {"id": "animateSpin", "name": "component.public.animateSpin"}
                    ],
                    "animateElement": [
                        {"id": "blink", "elementList": [
                            {"name": "component.public.animateSpeed", "type": 7, "value": 1, "min": 0.1, "key": "blinkSpeed"}
                        ]},
                        {"id": "millcolorGrad", "elementList": [
                            {"name": "component.public.startColor", "type": 2, "value": "#74f808", "key": "startColor"},
                            {"name": "component.public.stopColor", "type": 2, "value": "#f30b0b", "key": "stopColor"},
                            {"name": "component.public.animateSpeed", "type": 7, "value": 1, "min": 0.1, "key": "animateSpeed"}
                        ]},
                        {"id": "animateSpin", "elementList": [
                            {"name": "component.public.animateSpinSpeed", "type": 7, "value": 1, "min": 0.1, "key": "spinSpeed"},
                            {"name": "configComponent.bigScreen.border.border89Direction", "type": 6, "value": 0, "min": 1,
                             "key": "spinDirection",
                             "enumList": [
                                 {"value": 0, "option": "configComponent.bigScreen.border.border89DirectionForward"},
                                 {"value": 1, "option": "configComponent.bigScreen.border.border89DirectionNegative"}
                             ]}
                        ]}
                    ]
                },
                "action": [], "active": active, "dataBind": []
            }
        }
    }

def make_gauge(seed, x, y, w, h, title, dp_name, unit, min_val, max_val, z=5):
    cell_id = gen_uid(seed)
    dp = DP_MAP.get(dp_name, {'uuid': '', 'unit': unit})
    range_span = max_val - min_val
    a1_end = min_val + range_span * 0.3
    a2_end = min_val + range_span * 0.65
    active = [{
        "id": "ShowData",
        "name": "configComponent.ChartPublic.ShowData",
        "result": "",
        "isExpression": False,
        "condition": {
            "deviceSN": DEVICE_UUID, "DeviceName": DEVICE_NAME,
            "selectVideoType": 0, "isBandDevice": False, "bandType": 1,
            "dataID": dp['uuid'], "dataName": dp_name,
            "operator": "", "OperatorValue": "", "OperatorMaxValue": ""
        }
    }]
    return {
        "shape": "ism-view-chart-gauge-0",
        "id": cell_id, "x": x, "y": y, "width": w, "height": h,
        "zIndex": z, "visible": True,
        "position": {"x": x, "y": y},
        "size": {"width": w, "height": h},
        "data": {
            "detail": {
                "type": "ism-view-chart-gauge-0",
                "identifier": cell_id, "name": title,
                "style": {
                    "position": {"x": x, "y": y, "w": w, "h": h},
                    "backColor": "transparent", "zIndex": 1, "transform": 0,
                    "diy": [
                        {"name": "configComponent.ChartPublic.splitNumber", "type": 1, "value": 10, "key": "splitNumber"},
                        {"name": "configComponent.ChartPublic.ChartTitle", "type": 4, "value": title, "key": "ChartTitle"},
                        {"name": "configComponent.ChartPublic.ChartUnit", "type": 4, "value": unit, "key": "ChartUnit"},
                        {"name": "configComponent.ChartPublic.ChartTitleFontSize", "type": 1, "value": 13, "key": "ChartTitleFontSize"},
                        {"name": "configComponent.ChartPublic.ChartTitleFontColor", "type": 2, "value": "#00e5ff", "key": "ChartTitleFontColor"},
                        {"name": "configComponent.ChartPublic.ChartMinValue", "type": 7, "value": min_val, "key": "ChartMinValue"},
                        {"name": "configComponent.ChartPublic.ChartMaxValue", "type": 7, "value": max_val, "key": "ChartMaxValue"},
                        {"name": "configComponent.ChartPublic.ChartAxisTickColor", "type": 2, "value": "#5c8db8", "key": "ChartAxisTickColor"},
                        {"name": "configComponent.ChartPublic.ChartWidth", "type": 7, "value": 12, "key": "ChartWidth"},
                        {"name": "configComponent.ChartPublic.LabelDis", "type": 7, "value": 18, "key": "LabelDis"},
                        {"name": "configComponent.ChartPublic.ChartSplitLineWidth", "type": 7, "value": 2, "key": "ChartSplitLineWidth"},
                        {"name": "configComponent.ChartPublic.ChartSplitLineHeight", "type": 7, "value": 2, "key": "ChartSplitLineHeight"},
                        {"name": "configComponent.ChartPublic.Area1Range", "type": 4, "value": f"{int(min_val)}~{int(a1_end)}", "key": "Area1Range"},
                        {"name": "configComponent.ChartPublic.Area1Color", "type": 2, "value": "#4dabf7", "key": "Area1Color"},
                        {"name": "configComponent.ChartPublic.Area2Range", "type": 4, "value": f"{int(a1_end)}~{int(a2_end)}", "key": "Area2Range"},
                        {"name": "configComponent.ChartPublic.Area2Color", "type": 2, "value": "#69db7c", "key": "Area2Color"},
                        {"name": "configComponent.ChartPublic.Area3Color", "type": 2, "value": "#ff6b6b", "key": "Area3Color"},
                    ]
                },
                "animate": {
                    "selected": [],
                    "condition": {"deviceSN": "", "selectVideoType": 0, "isBandDevice": False,
                                  "bandType": 1, "dataID": "", "dataName": "",
                                  "operator": "", "OperatorValue": "", "OperatorMaxValue": ""},
                    "isExpression": False,
                    "animateList": [
                        {"id": "blink", "name": "component.public.animateBlink"},
                        {"id": "Zoom", "name": "component.public.Zoom"},
                        {"id": "animateSpin", "name": "component.public.animateSpin"}
                    ],
                    "animateElement": [
                        {"id": "blink", "elementList": [
                            {"name": "component.public.animateSpeed", "type": 7, "value": 1, "min": 0.1, "key": "blinkSpeed"}
                        ]},
                        {"id": "millcolorGrad", "elementList": [
                            {"name": "component.public.startColor", "type": 2, "value": "#74f808", "key": "startColor"},
                            {"name": "component.public.stopColor", "type": 2, "value": "#f30b0b", "key": "stopColor"},
                            {"name": "component.public.animateSpeed", "type": 7, "value": 1, "min": 0.1, "key": "animateSpeed"}
                        ]},
                        {"id": "animateSpin", "elementList": [
                            {"name": "component.public.animateSpinSpeed", "type": 7, "value": 1, "min": 0.1, "key": "spinSpeed"},
                            {"name": "configComponent.bigScreen.border.border89Direction", "type": 6, "value": 0, "min": 1,
                             "key": "spinDirection",
                             "enumList": [
                                 {"value": 0, "option": "configComponent.bigScreen.border.border89DirectionForward"},
                                 {"value": 1, "option": "configComponent.bigScreen.border.border89DirectionNegative"}
                             ]}
                        ]}
                    ]
                },
                "action": [], "active": active, "dataBind": []
            }
        }
    }

def _nav_action(target_page_id):
    """Generate navigation action to a target page (ISMRender expects type=click + action=link)."""
    return [{
        "type": "click",
        "action": "link",
        "link": {
            "linkType": "Inside",
            "isPopUp": False,
            "Inside": {
                "displayUUID": MODEL_ID,
                "pageUUID": target_page_id,
                "displayType": 1
            }
        }
    }]


def make_breadcrumb(seed, x, y, segments, z=20):
    """Build breadcrumb text cells. segments: [(text, color, target_page_id|None), ...]"""
    cells_out = []
    cx = x
    for si, (text, color, target) in enumerate(segments):
        if si > 0:
            sep_id = f'{seed}-sep-{si}'
            cells_out.append(make_text(sep_id, cx, y, 18, 22, '›',
                                       color='#475569', font_size=FONT_BREAD - 1, z=z))
            cx += 18
        action = _nav_action(target) if target else None
        seg_w = max(len(text) * 14 + 10, 80)
        cells_out.append(make_text(f'{seed}-{si}', cx, y, seg_w, 22, text,
                                   color=color, font_size=FONT_BREAD - 1, z=z, action=action))
        cx += seg_w + 4
    return cells_out


def build_header_cells(seed_prefix, breadcrumb_segments):
    """SCADAMonitor-style top header bar (56px)."""
    out = []
    out.append(make_border_box1(f'{seed_prefix}-header-bg', 0, 0, 1920, HEADER_H, z=0))
    out.append(make_decoration3(f'{seed_prefix}-header-deco-l', 0, 0, 200, HEADER_H, z=1))
    out.append(make_decoration3(f'{seed_prefix}-header-deco-r', 1720, 0, 200, HEADER_H, z=1))
    out.append(make_text(f'{seed_prefix}-header-logo', 220, 10, 36, 36, '⚡',
                         color=C_ACCENT, font_size=26, z=10))
    out.append(make_text(f'{seed_prefix}-header-title', 260, 4, 420, 22, '航信机房电力监控系统',
                         color=C_TEXT, font_size=FONT_TITLE, z=10))
    out.append(make_text(f'{seed_prefix}-header-subtitle', 268, 28, 420, 14, 'NCC ROOM POWER SCADA MONITOR',
                         color=C_TEXT_DIM, font_size=FONT_SUBTITLE, z=10))
    out.extend(make_breadcrumb(f'{seed_prefix}-header-crumb', 760, 20, breadcrumb_segments, z=20))
    out.append(make_text(f'{seed_prefix}-header-clock', 1540, 18, 180, 22, '系统运行中',
                         color=C_TEXT_MUTED, font_size=FONT_SUBTITLE, z=10))
    out.append(make_text(f'{seed_prefix}-header-status', 1730, 18, 160, 22, '🟢 在线',
                         color=C_GREEN, font_size=FONT_SUBTITLE, z=10))
    return out


def build_sidebar_cells(seed_prefix='nav'):
    """280px sidebar: building tree + alarm panel (mirrors SCADAMonitor)."""
    out = []
    out.append(make_border_box8(f'{seed_prefix}-sidebar-frame', SIDEBAR_X + 2, BODY_Y + 2,
                                SIDEBAR_W - 4, 1080 - BODY_Y - 4, z=0))
    nav_y = BODY_Y + 12
    out.append(make_text(f'{seed_prefix}-nav-title', SIDEBAR_X + 12, nav_y, 240, 20, '📍 建筑导航',
                         color=C_TEXT_MUTED, font_size=FONT_PANEL, z=8))
    nav_y += 28
    out.append(make_text(f'{seed_prefix}-nav-root', SIDEBAR_X + 10, nav_y, 240, 18, '🏭 航信机房',
                         color=C_ACCENT, font_size=FONT_NAV_BLDG, z=20))
    nav_y += 22
    for bi, bldg in enumerate(buildings):
        bldg_action = _nav_action(PAGE_ID_BUILDING)
        row_y = nav_y
        out.append(make_border_box1(f'{seed_prefix}-nav-bldg-row-{bi}', SIDEBAR_X + 6, row_y,
                                    SIDEBAR_W - 12, BLDG_ROW_H, z=3, action=bldg_action))
        ty = vcenter(row_y, BLDG_ROW_H, 18)
        out.append(make_text(f'{seed_prefix}-nav-bldg-icon-{bi}', SIDEBAR_X + 10, ty, 18, 18, '🏢',
                             color=C_ACCENT, font_size=14, z=20, action=bldg_action))
        out.append(make_text(f'{seed_prefix}-nav-bldg-name-{bi}', SIDEBAR_X + 32, ty, 150, 18, bldg['name'],
                             color=C_TEXT, font_size=FONT_NAV_BLDG, z=20, action=bldg_action))
        out.append(make_text(f'{seed_prefix}-nav-bldg-count-{bi}', SIDEBAR_X + 200, ty, 68, 18,
                             f'{bldg["device_count"]}台', color=C_TEXT_DIM, font_size=11, z=20,
                             action=bldg_action))
        nav_y += BLDG_ROW_H
        for fi, floor in enumerate(bldg['floors']):
            floor_action = _nav_action(PAGE_ID_FLOOR)
            fy = nav_y
            out.append(make_border_box1(f'{seed_prefix}-nav-flr-row-{bi}-{fi}', SIDEBAR_X + FLR_INDENT - 4, fy,
                                        SIDEBAR_W - FLR_INDENT - 8, FLR_ROW_H, z=3, action=floor_action))
            fty = vcenter(fy, FLR_ROW_H, 16)
            out.append(make_text(f'{seed_prefix}-nav-flr-icon-{bi}-{fi}', SIDEBAR_X + FLR_INDENT, fty, 14, 16, '📋',
                                 color=C_TEXT_MUTED, font_size=11, z=20, action=floor_action))
            out.append(make_text(f'{seed_prefix}-nav-flr-name-{bi}-{fi}', SIDEBAR_X + FLR_INDENT + 18, fty, 130, 16,
                                 floor['name'], color=C_TEXT_MUTED, font_size=FONT_NAV_FLR, z=20, action=floor_action))
            alarm_cnt = sum(1 for d in floor['devices'] if d['status'] != 1)
            badge_color = '#ef4444' if alarm_cnt else C_TEXT_DIM
            out.append(make_text(f'{seed_prefix}-nav-flr-badge-{bi}-{fi}', SIDEBAR_X + 210, fty, 58, 16,
                                 str(alarm_cnt), color=badge_color, font_size=10, z=20, action=floor_action))
            nav_y += FLR_ROW_H + FLR_ROW_GAP
        nav_y += BLDG_ROW_GAP
    alarm_y = max(nav_y + 16, 860)
    out.append(make_border_box1(f'{seed_prefix}-alarm-panel', SIDEBAR_X + 8, alarm_y, SIDEBAR_W - 16, 200, z=2))
    out.append(make_text(f'{seed_prefix}-alarm-title', SIDEBAR_X + 16, alarm_y + 10, 110, 20,
                         '🚨 实时告警', color=C_TEXT_MUTED, font_size=FONT_PANEL, z=8))
    out.append(make_text(f'{seed_prefix}-alarm-count', SIDEBAR_X + 150, alarm_y + 10, 40, 20,
                         '0', color='#ffffff', font_size=11, z=8))
    out.append(make_text(f'{seed_prefix}-alarm-empty', SIDEBAR_X + 16, alarm_y + 40, SIDEBAR_W - 32, 140,
                         '✅ 当前无告警', color=C_TEXT_DIM, font_size=12, z=8))
    out.append(make_border_box1(f'{seed_prefix}-sidebar-divider', SIDEBAR_W + 4, BODY_Y + 8, 2, 1016, z=1))
    return out


def report_overlaps(cells, label):
    overlaps = find_text_overlaps(cells, pad=4)
    if overlaps:
        print(f"⚠️  {label}: {len(overlaps)} visual text overlaps detected")
        for a, b in overlaps[:8]:
            print(f"   '{a['text']}' ↔ '{b['text']}'")
    else:
        print(f"✓ {label}: no text overlaps")
    return overlaps


# ════════════════════════════════════════════════════════
# LEVEL 0: OVERVIEW PAGE (page_name='main', id=8)
# Layout mirrors SCADAMonitor: header + 280px sidebar + main panels
# ════════════════════════════════════════════════════════

cells = []

# ── Shared SCADAMonitor shell: header + sidebar ──
cells.extend(build_header_cells('ov', [
    ('📊 全局总览', C_TEXT, PAGE_ID_MAIN),
]))
cells.extend(build_sidebar_cells('ov'))

# ── Main content (SCADAMonitor drillLevel=0 layout) ──
online_count = sum(1 for b in buildings for f in b['floors'] for d in f['devices'] if d['status'] == 1)

stats_y = BODY_Y + 16
stats_h = 100
card_w = int((MAIN_W - 48) / 4)
card_gap = 16
card_xs = [MAIN_X + i * (card_w + card_gap) for i in range(4)]

stat_configs = [
    ('stat-power', '⚡', '0', 'kW', '总功率', '总有功功率', '#667eea'),
    ('stat-energy', '📊', '---', 'kWh', '今日用电量', None, '#f093fb'),
    ('stat-online', '🖥', f'{online_count}/{TOTAL_DEVICES}', '', '在线设备', None, '#4facfe'),
    ('stat-alarm', '🔔', '0', '', '活跃告警', None, '#fa709a'),
]

for i, (seed, icon, val, unit, label, dp_name, accent) in enumerate(stat_configs):
    cx = card_xs[i]
    cells.append(make_border_box1(f'{seed}-bg', cx, stats_y, card_w, stats_h, z=2))
    cells.append(make_decoration5(f'{seed}-deco', cx + card_w - 90, stats_y - 8, 80, 80, z=1))
    cells.append(make_text(f'{seed}-icon', cx + 16, stats_y + 24, 40, 40, icon,
                           color=accent, font_size=24, z=5))
    display_val = f'{val} {unit}'.strip() if unit else val
    cells.append(make_text(
        f'{seed}-val', cx + 68, stats_y + 16, card_w - 80, 36,
        display_val, color=C_TEXT, font_size=FONT_KPI_VALUE, z=5,
        data_bound=(dp_name is not None), dp_name=dp_name
    ))
    cells.append(make_text(
        f'{seed}-lab', cx + 68, stats_y + 62, card_w - 80, 22,
        label, color=C_TEXT_DIM, font_size=FONT_KPI_LABEL, z=5
    ))

panel_top_y = stats_y + stats_h + 16
panel_h = 380
left_w = int((MAIN_W - 16) * 0.55)
right_w = MAIN_W - left_w - 16
right_x = MAIN_X + left_w + 16

cells.append(make_border_box1('panel-topology', MAIN_X, panel_top_y, left_w, panel_h, z=2))
cells.append(make_text('panel-topo-title', MAIN_X + 12, panel_top_y + 10, 400, 22,
                       '⚡ 系统拓扑概览', color=C_TEXT_MUTED, font_size=FONT_PANEL, z=6))
topo_y = panel_top_y + 44
for bi, bldg in enumerate(buildings[:3]):
    bldg_action = _nav_action(PAGE_ID_BUILDING)
    floor_count = min(len(bldg['floors']), 3)
    block_h = 34 + floor_count * 30
    cells.append(make_border_box1(f'topo-bldg-bg-{bi}', MAIN_X + 12, topo_y, left_w - 24, block_h, z=3,
                                  action=bldg_action))
    cells.append(make_text(f'topo-bldg-name-{bi}', MAIN_X + 24, topo_y + 8, left_w - 48, 18,
                           f'🏢 {bldg["name"]}', color=C_TEXT, font_size=13, z=6, action=bldg_action))
    floor_y = topo_y + 34
    for fi, floor in enumerate(bldg['floors'][:floor_count]):
        floor_action = _nav_action(PAGE_ID_FLOOR)
        dot_count = min(len(floor['devices']), 8)
        dots = '🟢' * dot_count
        cells.append(make_text(f'topo-flr-name-{bi}-{fi}', MAIN_X + 36, floor_y, left_w - 72, 14,
                               floor['name'], color=C_TEXT_MUTED, font_size=10, z=6, action=floor_action))
        cells.append(make_text(f'topo-flr-dots-{bi}-{fi}', MAIN_X + 36, floor_y + 16, left_w - 72, 12,
                               dots, color=C_GREEN, font_size=8, z=6, action=floor_action))
        floor_y += 30
    topo_y += block_h + 12

cells.append(make_border_box1('panel-chart', right_x, panel_top_y, right_w, panel_h, z=2))
cells.append(make_text('panel-chart-title', right_x + 12, panel_top_y + 10, 300, 22,
                       '📈 功率趋势 (24h)', color=C_TEXT_MUTED, font_size=FONT_PANEL, z=6))
cells.append(make_smooth_chart(
    'chart-trend', right_x + 12, panel_top_y + 38, right_w - 24, panel_h - 50,
    title='功率趋势 (24h)',
    dp_names=['总有功功率', '总无功功率', '总视在功率'],
    z=5
))

grid_y = panel_top_y + panel_h + 16
grid_h = 1080 - grid_y - 16
cells.append(make_border_box1('panel-device-grid', MAIN_X, grid_y, MAIN_W, grid_h, z=2))
cells.append(make_text('panel-grid-title', MAIN_X + 12, grid_y + 10, 400, 22,
                       '🏭 设备运行状态总览', color=C_TEXT_MUTED, font_size=FONT_PANEL, z=6))

dev_card_w = 140
dev_card_h = 88
dev_gap = 10
dev_cols = max(1, (MAIN_W - 24) // (dev_card_w + dev_gap))
dev_idx = 0
dev_start_y = grid_y + 40
row_y = dev_start_y
col = 0
max_rows = 4
placed_rows = 0
for bldg in buildings:
    if placed_rows >= max_rows:
        break
    cells.append(make_text(f'grid-bldg-title-{bldg["sid"]}', MAIN_X + 12, row_y, MAIN_W - 24, 20,
                           f'🏢 {bldg["name"]}', color=C_TEXT_MUTED, font_size=FONT_PANEL, z=6))
    row_y += 28
    col = 0
    for floor in bldg['floors']:
        for dev in floor['devices']:
            if placed_rows >= max_rows:
                break
            dx = MAIN_X + 12 + col * (dev_card_w + dev_gap)
            dy = row_y
            dev_action = _nav_action(PAGE_ID_DEVICE)
            status_color = C_GREEN if dev['status'] == 1 else C_TEXT_DIM
            status_text = '运行中' if dev['status'] == 1 else '离线'
            cells.append(make_border_box1(f'dev-card-bg-{dev_idx}', dx, dy, dev_card_w, dev_card_h, z=3,
                                          action=dev_action))
            cells.append(make_text(f'dev-card-icon-{dev_idx}', dx + 52, dy + 8, 36, 28, '⚡',
                                   color=status_color, font_size=22, z=6, action=dev_action))
            short_name = dev['name'][-12:] if len(dev['name']) > 12 else dev['name']
            cells.append(make_text(f'dev-card-name-{dev_idx}', dx + 6, dy + 40, dev_card_w - 12, 18,
                                   short_name, color=C_TEXT_MUTED, font_size=10, z=6, action=dev_action))
            cells.append(make_text(f'dev-card-st-{dev_idx}', dx + 20, dy + 62, dev_card_w - 40, 16,
                                   status_text, color=status_color, font_size=10, z=6, action=dev_action))
            dev_idx += 1
            col += 1
            if col >= dev_cols:
                col = 0
                row_y += dev_card_h + dev_gap
                placed_rows += 1
        if col > 0:
            col = 0
            row_y += dev_card_h + dev_gap
            placed_rows += 1
    row_y += 8

report_overlaps(cells, 'Overview layout')


# ════════════════════════════════════════════════════════
# ENCODE & WRITE OVERVIEW PAGE (id=8)
# ════════════════════════════════════════════════════════
components_json_main = json.dumps({"cells": cells}, ensure_ascii=False)
comp_b64_main = base64.b64encode(components_json_main.encode()).decode()

print(f"\n=== LEVEL 0: OVERVIEW PAGE ===")
print(f"Total cells: {len(cells)}")
print(f"JSON size: {len(components_json_main)} chars")
print(f"Base64 size: {len(comp_b64_main)} chars")

decoded = json.loads(base64.b64decode(comp_b64_main).decode())
assert 'cells' in decoded and len(decoded['cells']) == len(cells)
print("Roundtrip verification PASSED")

cur.execute("UPDATE display_model_layer SET components=%s, updated_at=NOW() WHERE id=8", (comp_b64_main,))
conn.commit()
print(f"Database UPDATE executed for id=8 (overview), rows affected: {cur.rowcount}")

# Also update layer background to dark theme
cur.execute(
    "UPDATE display_model_layer SET layer='{\"width\":1920,\"height\":1080,\"autoSize\":0,\"Padding\":0,\"gridSize\":10,\"background\":\"#0a0e17\"}', updated_at=NOW() WHERE id=8"
)
conn.commit()
print(f"Layer background updated for id=8, rows affected: {cur.rowcount}")


# ════════════════════════════════════════════════════════
# LEVEL 1: BUILDING DETAIL PAGE (page_name='building-detail')
# ════════════════════════════════════════════════════════
print(f"\n=== LEVEL 1: BUILDING DETAIL PAGE ===")

first_bldg_name = buildings[0]['name'] if buildings else '航信机房'
sample_building = buildings[0] if buildings else {'floors': [], 'device_count': 0}
first_floor_name = sample_building['floors'][0]['name'] if sample_building['floors'] else 'S18设备组'
sample_floor = sample_building['floors'][0] if sample_building['floors'] else {'devices': [], 'count': 0}

bldg_cells = []
bldg_cells.extend(build_header_cells('bldg', [
    ('📊 全局总览', C_TEXT_DIM, PAGE_ID_MAIN),
    (first_bldg_name, C_TEXT, None),
]))
bldg_cells.extend(build_sidebar_cells('bldg-nav'))

back_to_main = _nav_action(PAGE_ID_MAIN)
level_y = BODY_Y + 16
bldg_cells.append(make_text('bldg-back-btn', MAIN_X, level_y, 140, 32, '← 返回总览',
                            color=C_ACCENT, font_size=14, z=20, action=back_to_main))
bldg_cells.append(make_text('bldg-level-title', MAIN_X + 160, level_y, 800, 28,
                            f'🏢 {first_bldg_name}', color=C_TEXT, font_size=22, z=10))
alarm_cnt = sum(1 for f in sample_building['floors'] for d in f['devices'] if d['status'] != 1)
bldg_cells.append(make_text('bldg-level-sub', MAIN_X + 160, level_y + 34, 800, 18,
                            f'{sample_building["device_count"]}台设备 · {alarm_cnt}条告警',
                            color=C_TEXT_DIM, font_size=13, z=10))

card_start_x = MAIN_X
card_start_y = level_y + 72
card_w_item = 300
card_h_item = 140
cards_per_row = max(1, (MAIN_W + 16) // (card_w_item + 16))
card_gap_x = 16
card_gap_y = 16

for fi, floor in enumerate(sample_building['floors']):
    row = fi // cards_per_row
    col = fi % cards_per_row
    cx = card_start_x + col * (card_w_item + card_gap_x)
    cy = card_start_y + row * (card_h_item + card_gap_y)
    floor_nav = _nav_action(PAGE_ID_FLOOR)
    running = sum(1 for d in floor['devices'] if d['status'] == 1)
    offline = floor['count'] - running
    bldg_cells.append(make_border_box1(f'bldg-card-bg-{fi}', cx, cy, card_w_item, card_h_item, z=2, action=floor_nav))
    bldg_cells.append(make_text(f'bldg-card-name-{fi}', cx + 14, cy + 14, card_w_item - 28, 22,
                                f'📋 {floor["name"]}', color=C_TEXT, font_size=15, z=5, action=floor_nav))
    bldg_cells.append(make_text(f'bldg-card-count-{fi}', cx + 14, cy + 40, card_w_item - 28, 18,
                                f'{floor["count"]}台设备', color=C_TEXT_DIM, font_size=12, z=5, action=floor_nav))
    bldg_cells.append(make_text(f'bldg-card-run-{fi}', cx + 14, cy + 72, 90, 18,
                                f'🟢 {running}运行', color=C_GREEN, font_size=11, z=5, action=floor_nav))
    bldg_cells.append(make_text(f'bldg-card-alarm-{fi}', cx + 110, cy + 72, 80, 18,
                                f'🔴 0告警', color='#ef4444', font_size=11, z=5, action=floor_nav))
    bldg_cells.append(make_text(f'bldg-card-stop-{fi}', cx + 200, cy + 72, 80, 18,
                                f'⏸ {offline}停止', color=C_TEXT_DIM, font_size=11, z=5, action=floor_nav))

report_overlaps(bldg_cells, 'Building detail layout')

# Encode building-detail page
components_json_bldg = json.dumps({"cells": bldg_cells}, ensure_ascii=False)
comp_b64_bldg = base64.b64encode(components_json_bldg.encode()).decode()
print(f"Total building-detail cells: {len(bldg_cells)}")
print(f"JSON size: {len(components_json_bldg)} chars")


# ════════════════════════════════════════════════════════
# LEVEL 2: FLOOR DETAIL PAGE (page_name='floor-detail')
# ════════════════════════════════════════════════════════
print(f"\n=== LEVEL 2: FLOOR DETAIL PAGE ===")

floor_cells = []
floor_cells.extend(build_header_cells('floor', [
    ('📊 全局总览', C_TEXT_DIM, PAGE_ID_MAIN),
    (first_bldg_name, C_TEXT_DIM, PAGE_ID_BUILDING),
    (first_floor_name, C_TEXT, None),
]))
floor_cells.extend(build_sidebar_cells('floor-nav'))

back_to_bldg = _nav_action(PAGE_ID_BUILDING)
level_y = BODY_Y + 16
floor_cells.append(make_text('floor-back-btn', MAIN_X, level_y, 160, 32, f'← {first_bldg_name}',
                             color=C_ACCENT, font_size=14, z=20, action=back_to_bldg))
floor_cells.append(make_text('floor-level-title', MAIN_X + 180, level_y, 800, 28,
                             f'📋 {first_floor_name}', color=C_TEXT, font_size=22, z=10))
floor_cells.append(make_text('floor-level-sub', MAIN_X + 180, level_y + 40, 600, 18,
                             f'{sample_floor.get("count", len(sample_floor.get("devices", [])))}台设备',
                             color=C_TEXT_DIM, font_size=13, z=10))

table_x = MAIN_X
table_y = level_y + 72
table_w = MAIN_W
row_h = 42
col_widths = [50, 320, 100, 120, 140, 140, 120, 110, 120]
col_headers = ['#', '设备名称', '状态', '实时功率', '电流A/B/C', '温度', '协议', '操作', '']
col_starts = [table_x]
for w in col_widths[:-1]:
    col_starts.append(col_starts[-1] + w)

floor_cells.append(make_border_box1('floor-table-frame', MAIN_X, table_y, MAIN_W, 900, z=1))
floor_cells.append(make_border_box1('floor-th-bg', table_x, table_y + 8, table_w, row_h, z=3))
for hi, (hdr, cs) in enumerate(zip(col_headers[:-1], col_starts)):
    if not hdr:
        continue
    floor_cells.append(make_text(f'floor-th-{hi}', cs + 8, table_y + 18, col_widths[hi] - 12, 20,
                                  hdr, color=C_TEXT_DIM, font_size=11, z=4))

for di, dev in enumerate(sample_floor['devices'][:15]):
    ry = table_y + 8 + row_h + di * row_h
    floor_cells.append(make_text(f'floor-row-num-{di}', col_starts[0] + 8, ry + 10, 30, 20,
                                  str(di + 1), color=C_TEXT_DIM, font_size=11, z=4))
    dev_action = _nav_action(PAGE_ID_DEVICE)
    floor_cells.append(make_text(f'floor-row-name-{di}', col_starts[1] + 8, ry + 10, col_widths[1] - 12, 20,
                                  dev['name'], color=C_ACCENT, font_size=13, z=4, action=dev_action))
    status_text = '运行中' if dev['status'] == 1 else '离线'
    status_color = C_GREEN if dev['status'] == 1 else C_TEXT_DIM
    floor_cells.append(make_text(f'floor-row-stat-{di}', col_starts[2] + 8, ry + 10, col_widths[2] - 12, 20,
                                  status_text, color=status_color, font_size=12, z=4))
    for ci, placeholder in enumerate(['-- kW', '--/--/-- A', '--°C', 'MODBUS'], start=3):
        floor_cells.append(make_text(f'floor-row-data-{di}-{ci}', col_starts[ci] + 8, ry + 10,
                                      col_widths[ci] - 12, 20, placeholder, color=C_TEXT_DIM, font_size=12, z=4))
    detail_btn = _nav_action(PAGE_ID_DEVICE)
    floor_cells.append(make_border_box1(f'floor-row-btn-bg-{di}', col_starts[7] + 8, ry + 4, 90, 28, z=4))
    floor_cells.append(make_text(f'floor-row-btn-{di}', col_starts[7] + 14, ry + 8, 78, 20,
                                  '详情', color=C_ACCENT, font_size=12, z=5, action=detail_btn))

table_end_y = table_y + 8 + row_h + len(sample_floor['devices'][:15]) * row_h + 20
floor_cells.append(make_text('floor-summary', table_x, table_end_y, table_w, 22,
                              f'共 {sample_floor["count"]} 台设备 | 运行: {sum(1 for d in sample_floor["devices"] if d["status"]==1)}台 | 离线: {sum(1 for d in sample_floor["devices"] if d["status"]!=1)}台',
                              color=C_TEXT_MUTED, font_size=13, z=5))

report_overlaps(floor_cells, 'Floor detail layout')

# Encode floor-detail page
components_json_floor = json.dumps({"cells": floor_cells}, ensure_ascii=False)
comp_b64_floor = base64.b64encode(components_json_floor.encode()).decode()
print(f"Total floor-detail cells: {len(floor_cells)}")
print(f"JSON size: {len(components_json_floor)} chars")


# ════════════════════════════════════════════════════════
# LEVEL 3: DEVICE DETAIL PAGE (page_name='device-detail', id=10)
# ════════════════════════════════════════════════════════
print(f"\n=== LEVEL 3: DEVICE DETAIL PAGE ===")

detail_cells = []
detail_cells.extend(build_header_cells('detail', [
    ('📊 全局总览', C_TEXT_DIM, PAGE_ID_MAIN),
    (first_bldg_name, C_TEXT_DIM, PAGE_ID_BUILDING),
    (first_floor_name, C_TEXT_DIM, PAGE_ID_FLOOR),
    (DEVICE_NAME, C_TEXT, None),
]))
detail_cells.extend(build_sidebar_cells('detail-nav'))

back_to_floor = _nav_action(PAGE_ID_FLOOR)
level_y = BODY_Y + 16
detail_cells.append(make_text('detail-back-floor-btn', MAIN_X, level_y, 180, 32, f'← {first_floor_name}',
                              color=C_ACCENT, font_size=14, z=20, action=back_to_floor))
detail_cells.append(make_text('detail-level-title', MAIN_X + 200, level_y, 700, 28,
                              f'🔧 {DEVICE_NAME}', color=C_TEXT, font_size=22, z=10))
detail_cells.append(make_text('detail-level-status', MAIN_X + 920, level_y + 4, 200, 22,
                              '运行中', color=C_GREEN, font_size=13, z=10))

panel_top = level_y + 64
panel_h = 360
col_w = (MAIN_W - 32) // 3
left_x = MAIN_X
mid_x = MAIN_X + col_w + 16
right_x = MAIN_X + (col_w + 16) * 2

# Left panel: Basic parameters
detail_cells.append(make_border_box1('detail-left-panel', left_x, panel_top, col_w, panel_h, z=2))
detail_cells.append(make_text('detail-left-title', left_x + 15, panel_top + 8, 400, 22, '📋 基本参数',
                              color=C_ACCENT, font_size=FONT_PANEL + 2, z=6))
basic_params = [
    ('设备名称', DEVICE_NAME), ('设备类型', '多功能电力仪表'),
    ('通信协议', 'Modbus RTU'), ('通信地址', '1'),
    ('所属机房', 'NCC 航信机房'), ('所属区域', '1A1_U11柜 S18设备组'),
    ('打包时间', '500ms'), ('在线状态', '🟢 运行中'),
]
bp_y = panel_top + 44
for bi, (bk, bv) in enumerate(basic_params):
    by = bp_y + bi * 36
    detail_cells.append(make_text(f'detail-bp-key-{bi}', left_x + 15, by, 140, 22,
                                  bk, color=C_TEXT_DIM, font_size=FONT_PARAM_LABEL + 1, z=6))
    detail_cells.append(make_text(f'detail-bp-val-{bi}', left_x + 160, by, col_w - 180, 22,
                                  bv, color=C_TEXT, font_size=FONT_PARAM_VAL - 2, z=6))

detail_cells.append(make_border_box1('detail-mid-panel', mid_x, panel_top, col_w, panel_h, z=2))
detail_cells.append(make_text('detail-mid-title', mid_x + 15, panel_top + 8, 400, 22, '📊 实时参数',
                              color=C_ACCENT, font_size=FONT_PANEL + 2, z=6))
rt_params = [
    ('AB线电压', 'V'), ('BC线电压', 'V'), ('CA线电压', 'V'),
    ('A相电流', 'A'), ('B相电流', 'A'), ('C相电流', 'A'),
    ('中性线电流', 'A'), ('频率', 'Hz'),
]
rtp_y = panel_top + 44
for ri, (rname, runit) in enumerate(rt_params):
    ry = rtp_y + ri * 36
    detail_cells.append(make_text(f'detail-rtp-key-{ri}', mid_x + 15, ry, 150, 22,
                                  rname, color=C_TEXT_DIM, font_size=FONT_PARAM_LABEL + 1, z=6))
    detail_cells.append(make_text(f'detail-rtp-val-{ri}', mid_x + 170, ry, 120, 22,
                                  '---', color=C_ACCENT, font_size=FONT_PARAM_VAL, z=6,
                                  data_bound=True, dp_name=rname))
    detail_cells.append(make_text(f'detail-rtp-unit-{ri}', mid_x + 300, ry, 60, 22,
                                  runit, color=C_TEXT_DIM, font_size=FONT_PARAM_LABEL, z=6))

detail_cells.append(make_border_box1('detail-right-panel', right_x, panel_top, col_w, panel_h, z=2))
detail_cells.append(make_text('detail-right-title', right_x + 15, panel_top + 8, 400, 22, '⚡ 功率参数',
                              color=C_ACCENT, font_size=FONT_PANEL + 2, z=6))
pw_params = [
    ('总有功功率', 'kW'), ('总无功功率', 'kW'), ('总视在功率', 'kW'),
    ('总功率因数', ''), ('正有功电度', 'kWh'),
]
pw_y = panel_top + 44
for pi, (pname, punit) in enumerate(pw_params):
    py = pw_y + pi * 44
    detail_cells.append(make_text(f'detail-pw-key-{pi}', right_x + 15, py, 150, 22,
                                  pname, color=C_TEXT_DIM, font_size=FONT_PARAM_LABEL + 1, z=6))
    detail_cells.append(make_text(f'detail-pw-val-{pi}', right_x + 170, py, 140, 30,
                                  '---', color=C_GREEN, font_size=FONT_KPI_VALUE - 4, z=6,
                                  data_bound=True, dp_name=pname))
    detail_cells.append(make_text(f'detail-pw-unit-{pi}', right_x + 320, py, 60, 22,
                                  punit, color=C_TEXT_DIM, font_size=FONT_PARAM_LABEL, z=6))

chart_y = panel_top + panel_h + 16
chart_w = int(MAIN_W * 0.65)
detail_cells.append(make_border_box1('detail-chart-panel', MAIN_X, chart_y, chart_w, 280, z=2))
detail_cells.append(make_text('detail-chart-title', MAIN_X + 15, chart_y + 8, 400, 22, '📈 24小时功率曲线',
                              color=C_ACCENT, font_size=FONT_PANEL + 2, z=6))
detail_cells.append(make_smooth_chart(
    'detail-chart', MAIN_X + 15, chart_y + 36, chart_w - 30, 230,
    title='设备功率趋势',
    dp_names=['总有功功率', '总无功功率'],
    z=5
))

status_x = MAIN_X + chart_w + 16
status_w = MAIN_W - chart_w - 16
detail_cells.append(make_border_box1('detail-status-panel', status_x, chart_y, status_w, 280, z=2))
detail_cells.append(make_text('detail-st-title', status_x + 15, chart_y + 8, 400, 22, '🔔 设备告警',
                              color=C_ACCENT, font_size=FONT_PANEL + 2, z=6))
detail_cells.append(make_text('detail-alarm-empty', status_x + 15, chart_y + 48, status_w - 30, 200,
                              '✅ 该设备无告警记录', color=C_TEXT_DIM, font_size=13, z=6))

report_overlaps(detail_cells, 'Device detail layout')

# Encode device detail page
components_json_detail = json.dumps({"cells": detail_cells}, ensure_ascii=False)
comp_b64_detail = base64.b64encode(components_json_detail.encode()).decode()
print(f"Total device-detail cells: {len(detail_cells)}")
print(f"JSON size: {len(components_json_detail)} chars")


# ════════════════════════════════════════════════════════
# WRITE ALL PAGES TO DATABASE
# ════════════════════════════════════════════════════════

# Upsert building-detail page
cur.execute(
    "SELECT id FROM display_model_layer WHERE model_id=%s AND page_name='building-detail' AND deleted_at IS NULL",
    (MODEL_ID,)
)
existing_bldg = cur.fetchone()
if existing_bldg:
    cur.execute("UPDATE display_model_layer SET components=%s, updated_at=NOW() WHERE id=%s",
                (comp_b64_bldg, existing_bldg[0]))
    print(f"Updated existing building-detail page id={existing_bldg[0]}")
else:
    cur.execute(
        """INSERT INTO display_model_layer
           (model_id, page_name, page_id, is_home, is_login, page_type, layer, components, created_at, updated_at)
           VALUES (%s, 'building-detail', %s, 0, 0, 1, '{"height":1080,"width":1920,"autoSize":1}', %s, NOW(), NOW())""",
        (MODEL_ID, PAGE_ID_BUILDING, comp_b64_bldg)
    )
    print(f"Inserted new building-detail page, id={cur.lastrowid}")
conn.commit()

# Upsert floor-detail page
cur.execute(
    "SELECT id FROM display_model_layer WHERE model_id=%s AND page_name='floor-detail' AND deleted_at IS NULL",
    (MODEL_ID,)
)
existing_floor = cur.fetchone()
if existing_floor:
    cur.execute("UPDATE display_model_layer SET components=%s, updated_at=NOW() WHERE id=%s",
                (comp_b64_floor, existing_floor[0]))
    print(f"Updated existing floor-detail page id={existing_floor[0]}")
else:
    cur.execute(
        """INSERT INTO display_model_layer
           (model_id, page_name, page_id, is_home, is_login, page_type, layer, components, created_at, updated_at)
           VALUES (%s, 'floor-detail', %s, 0, 0, 1, '{"height":1080,"width":1920,"autoSize":1}', %s, NOW(), NOW())""",
        (MODEL_ID, PAGE_ID_FLOOR, comp_b64_floor)
    )
    print(f"Inserted new floor-detail page, id={cur.lastrowid}")
conn.commit()

# Upsert device-detail page
cur.execute(
    "SELECT id FROM display_model_layer WHERE model_id=%s AND page_name='device-detail' AND deleted_at IS NULL",
    (MODEL_ID,)
)
existing_dev = cur.fetchone()
if existing_dev:
    cur.execute("UPDATE display_model_layer SET components=%s, updated_at=NOW() WHERE id=%s",
                (comp_b64_detail, existing_dev[0]))
    print(f"Updated existing device-detail page id={existing_dev[0]}")
else:
    cur.execute(
        """INSERT INTO display_model_layer
           (model_id, page_name, page_id, is_home, is_login, page_type, layer, components, created_at, updated_at)
           VALUES (%s, 'device-detail', %s, 0, 0, 1, '{"height":1080,"width":1920,"autoSize":1}', %s, NOW(), NOW())""",
        (MODEL_ID, PAGE_ID_DEVICE, comp_b64_detail)
    )
    print(f"Inserted new device-detail page, id={cur.lastrowid}")
conn.commit()

# ── Verify all pages ──
cur.execute("""
    SELECT id, page_name, is_home, page_id, LENGTH(components)
    FROM display_model_layer
    WHERE model_id=%s AND deleted_at IS NULL
    ORDER BY is_home DESC, id
""", (MODEL_ID,))
pages = cur.fetchall()
print(f"\n=== All pages for model {MODEL_ID} ===")
for p in pages:
    print(f"  id={p[0]}, name={p[1]}, is_home={p[2]}, page_id={p[3][:20]}..., comp_len={p[4]}")

# ── Summary ──
print(f"\n{'='*60}")
print(f"✅ Build complete!")
print(f"   Level 0 (overview):     {len(cells)} cells, {len(components_json_main)} chars")
print(f"   Level 1 (building-detail): {len(bldg_cells)} cells, {len(components_json_bldg)} chars")
print(f"   Level 2 (floor-detail):    {len(floor_cells)} cells, {len(components_json_floor)} chars")
print(f"   Level 3 (device-detail):   {len(detail_cells)} cells, {len(components_json_detail)} chars")
print(f"   Total cells across all pages: {len(cells) + len(bldg_cells) + len(floor_cells) + len(detail_cells)}")
print(f"\n   Page IDs:")
print(f"     MAIN (overview):    {PAGE_ID_MAIN}")
print(f"     BUILDING_DETAIL:    {PAGE_ID_BUILDING}")
print(f"     FLOOR_DETAIL:       {PAGE_ID_FLOOR}")
print(f"     DEVICE_DETAIL:      {PAGE_ID_DEVICE}")
print(f"\n   Open in browser: http://localhost:7080/#/ISMDisPlay/DisPlayRunApp?displayUUID={MODEL_ID}")
print(f"{'='*60}")

conn.close()
