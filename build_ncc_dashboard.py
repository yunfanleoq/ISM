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
import json
import base64
import sqlite3
import uuid as _uuid
from collections import defaultdict
from datetime import datetime
from pathlib import Path

# ── DB (SQLite for cloud / OceanBase via export script locally) ──
DB_PATH = Path(__file__).resolve().parent / 'ism_server_user' / 'data' / 'db' / 'ism.db'
conn = sqlite3.connect(str(DB_PATH))
cur = conn.cursor()

# ── Constants ───────────────────────────────────────
MODEL_ID = '043135ad-44be-e5d8-89be-3e54883c23a8'
PAGE_ID_MAIN = MODEL_ID
PAGE_ID_BUILDING = _uuid.uuid5(_uuid.NAMESPACE_DNS, 'ncc-dash-building-detail').hex
PAGE_ID_FLOOR = _uuid.uuid5(_uuid.NAMESPACE_DNS, 'ncc-dash-floor-detail').hex
UPS_MODEL_UUID = '13b6fe72-1ad2-969e-499c-a85d7cefdb6f'
METER_MODEL_UUID = '3d734984-56f6-5494-ad4c-dfc67ca28ac8'
PROJECT_UUID = '31bc90be-ebc4-dd61-ba9d-ce6e075e40e2'
LAYER_JSON = json.dumps({
    "width": 1920, "height": 1080, "autoSize": 0, "Padding": 0,
    "gridSize": 10, "background": "#070b14"
})

# Runtime binding context (switched per device-detail page)
class DeviceCtx:
    def __init__(self, dev=None):
        dev = dev or {}
        self.uuid = dev.get('uuid', '')
        self.name = dev.get('name', '')
        self.muid = dev.get('muid') or ''
        self.is_ups = self.name.startswith('UPS') or self.muid == UPS_MODEL_UUID
        self.dp_map = load_dp_map(self.muid) if self.muid else {}

CTX = DeviceCtx()

def set_ctx(dev):
    global CTX
    CTX = DeviceCtx(dev)

def load_dp_map(muid):
    cur.execute(
        'SELECT name, uuid, data_unit FROM modbus_devices_data_model WHERE muid=? ORDER BY id',
        (muid,),
    )
    return {r[0]: {'uuid': r[1], 'unit': r[2] or ''} for r in cur.fetchall()}

def device_page_id(dev_uuid):
    return _uuid.uuid5(_uuid.NAMESPACE_DNS, f'ncc-device-{dev_uuid}').hex

def floor_page_id(bldg_sid, floor_key):
    return _uuid.uuid5(_uuid.NAMESPACE_DNS, f'ncc-floor-{bldg_sid}-{floor_key}').hex

def building_page_id(bldg_sid):
    return _uuid.uuid5(_uuid.NAMESPACE_DNS, f'ncc-building-{bldg_sid}').hex

def dev_icon(dev):
    if dev.get('name', '').startswith('UPS') or dev.get('muid') == UPS_MODEL_UUID:
        return '🔋'
    return '⚡'

def floor_key_for_device(d):
    name = d.get('name', '')
    if name.startswith('UPS') or d.get('muid') == UPS_MODEL_UUID:
        return 'UPS设备组'
    parts = name.split('_')
    if len(parts) >= 3:
        return f'{parts[2]}设备组'
    return '其他设备'

cur.execute('SELECT COUNT(*) FROM monitor_list WHERE deleted_at IS NULL')
TOTAL_DEVICES = cur.fetchone()[0]

print(f"Total devices: {TOTAL_DEVICES}")

# ── Query real device hierarchy ─────────────────────
cur.execute("""
    SELECT uuid, name, sid, pid, type, muid, status
    FROM monitor_list
    WHERE project_uuid = ? AND deleted_at IS NULL
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
                floors[floor_key_for_device(d)].append(d)
            building_entry = {
                'uuid': uuid, 'name': name, 'sid': sid,
                'devices': type1_children,
                'device_count': len(type1_children),
                'floors': [{'key': k, 'name': k, 'devices': v, 'count': len(v)}
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
                floor_groups[floor_key_for_device(d)].append(d)
            sub = {
                'uuid': c0['uuid'], 'name': c0['name'], 'sid': c0['sid'],
                'devices': type1_gc, 'device_count': len(type1_gc),
                'floors': [{'key': k, 'name': k, 'devices': v, 'count': len(v)}
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
    if dp_name not in CTX.dp_map:
        return []
    dp = CTX.dp_map[dp_name]
    return [{
        "id": "ShowData",
        "name": "configComponent.variable.ShowData",
        "result": "",
        "isExpression": False,
        "condition": {
            "deviceSN": CTX.uuid,
            "DeviceName": CTX.name,
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

BLDG_ROW_H = 32
BLDG_ROW_GAP = 2
FLR_ROW_H = 26
FLR_ROW_GAP = 1
FLR_INDENT = 20
ALARM_PANEL_Y = 868
ALARM_PANEL_H = 1080 - ALARM_PANEL_Y - 8

# SCADAMonitor palette (v4 — deeper contrast + accent glow)
C_BG = '#070b14'
C_SIDEBAR = '#0a1020'
C_HEADER = '#111827'
C_BORDER = '#1e3a5f'
C_TEXT = '#f1f5f9'
C_TEXT_MUTED = '#94a3b8'
C_TEXT_DIM = '#64748b'
C_ACCENT = '#38bdf8'
C_ACCENT2 = '#818cf8'
C_GREEN = '#34d399'
C_AMBER = '#fbbf24'
C_RED = '#f87171'


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


def truncate_label(text, max_len=10):
    return text if len(text) <= max_len else text[: max_len - 1] + '…'


def all_ups_devices():
    ups = []
    for bldg in buildings:
        for floor in bldg['floors']:
            for dev in floor['devices']:
                if dev_icon(dev) == '🔋':
                    ups.append(dev)
    return ups


def append_device_cards(cells, devices, seed_prefix, start_x, start_y, area_w,
                        card_w=148, card_h=92, gap=12, max_count=None):
    """Render clickable device cards; returns (next_y, count)."""
    cols = max(1, (area_w - 12) // (card_w + gap))
    col = 0
    row_y = start_y
    count = 0
    for dev in devices:
        if max_count is not None and count >= max_count:
            break
        dx = start_x + 12 + col * (card_w + gap)
        dy = row_y
        dev_action = _nav_device_action(dev)
        status_color = C_GREEN if dev['status'] == 1 else C_TEXT_DIM
        status_text = '运行中' if dev['status'] == 1 else '离线'
        icon = dev_icon(dev)
        accent = '#22d3ee' if icon == '🔋' else C_ACCENT
        cells.append(make_border_box1(f'{seed_prefix}-bg-{count}', dx, dy, card_w, card_h, z=3,
                                      action=dev_action))
        cells.append(make_decoration5(f'{seed_prefix}-deco-{count}', dx + card_w - 72, dy - 6, 68, 68, z=1))
        cells.append(make_text(f'{seed_prefix}-icon-{count}', dx + card_w // 2 - 18, dy + 14, 36, 28, icon,
                               color=status_color, font_size=24, z=6, action=dev_action))
        short_name = dev['name'][-14:] if len(dev['name']) > 14 else dev['name']
        cells.append(make_text(f'{seed_prefix}-name-{count}', dx + 8, dy + 48, card_w - 16, 18,
                               short_name, color=C_TEXT, font_size=11, z=6, action=dev_action))
        cells.append(make_text(f'{seed_prefix}-st-{count}', dx + 24, dy + 70, card_w - 48, 16,
                               status_text, color=status_color, font_size=10, z=6, action=dev_action))
        count += 1
        col += 1
        if col >= cols:
            col = 0
            row_y += card_h + gap
    if col > 0:
        row_y += card_h + gap
    return row_y, count

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

def make_decoration6(seed, x, y, w, h, z=0):
    cell_id = gen_uid(seed)
    return {
        "shape": "dv-decoration6",
        "id": cell_id, "x": x, "y": y, "width": w, "height": h,
        "zIndex": z, "visible": True,
        "position": {"x": x, "y": y},
        "size": {"width": w, "height": h},
        "data": {
            "detail": {
                "type": "dv-decoration6",
                "identifier": cell_id,
                "name": seed,
                "style": {"position": {"x": x, "y": y, "w": w, "h": h}, "visible": 1, "diy": []},
                "animate": _base_animate(),
                "action": [], "active": [], "dataBind": []
            }
        }
    }


def make_border_box13(seed, x, y, w, h, z=1, action=None):
    cell_id = gen_uid(seed)
    return {
        "shape": "dv-border-box13",
        "id": cell_id, "x": x, "y": y, "width": w, "height": h,
        "zIndex": z, "visible": True,
        "position": {"x": x, "y": y},
        "size": {"width": w, "height": h},
        "data": {
            "detail": {
                "type": "dv-border-box13",
                "identifier": cell_id,
                "name": seed,
                "style": {"position": {"x": x, "y": y, "w": w, "h": h}, "visible": 1, "diy": []},
                "animate": _base_animate(),
                "action": action or [], "active": [], "dataBind": []
            }
        }
    }


def make_smooth_chart(seed, x, y, w, h, title, dp_names, z=5):
    cell_id = gen_uid(seed)
    active = []
    var_ids = ['ShowChartVariable1', 'ShowChartVariable2', 'ShowChartVariable3',
               'ShowChartVariable4', 'ShowChartVariable5']
    for i, dpn in enumerate(dp_names[:5]):
        dp = CTX.dp_map.get(dpn)
        if not dp:
            continue
        active.append({
            "id": var_ids[i],
            "name": "configComponent.variable.ShowData",
            "result": "",
            "isExpression": False,
            "condition": {
                "deviceSN": CTX.uuid, "DeviceName": CTX.name,
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
    dp = CTX.dp_map.get(dp_name, {'uuid': '', 'unit': unit})
    range_span = max_val - min_val
    a1_end = min_val + range_span * 0.3
    a2_end = min_val + range_span * 0.65
    active = [{
        "id": "ShowData",
        "name": "configComponent.ChartPublic.ShowData",
        "result": "",
        "isExpression": False,
        "condition": {
            "deviceSN": CTX.uuid, "DeviceName": CTX.name,
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


def _nav_device_action(dev):
    return _nav_action(device_page_id(dev['uuid']))


def upsert_page(page_name, page_id, comp_b64, is_home=0):
    cur.execute(
        "SELECT id FROM display_model_layer WHERE model_id=? AND page_id=? AND deleted_at IS NULL",
        (MODEL_ID, page_id),
    )
    row = cur.fetchone()
    if row:
        cur.execute(
            "UPDATE display_model_layer SET page_name=?, components=?, layer=?, updated_at=datetime('now') WHERE id=?",
            (page_name, comp_b64, LAYER_JSON, row[0]),
        )
        return row[0], 'updated'
    cur.execute(
        """INSERT INTO display_model_layer
           (model_id, page_name, page_id, is_home, is_login, page_type, layer, components, created_at, updated_at)
           VALUES (?, ?, ?, ?, 0, 1, ?, ?, datetime('now'), datetime('now'))""",
        (MODEL_ID, page_name, page_id, is_home, LAYER_JSON, comp_b64),
    )
    return cur.lastrowid, 'inserted'


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
    out.append(make_decoration6(f'{seed_prefix}-header-glow', 480, 8, 960, 40, z=1))
    out.append(make_text(f'{seed_prefix}-header-logo', 220, 12, 36, 36, '⚡',
                         color=C_ACCENT, font_size=26, z=10))
    out.append(make_text(f'{seed_prefix}-header-title', 262, 4, 380, 22, '航信机房电力监控',
                         color=C_TEXT, font_size=FONT_TITLE, z=10))
    out.append(make_text(f'{seed_prefix}-header-subtitle', 262, 34, 380, 12, 'NCC POWER SCADA',
                         color=C_TEXT_DIM, font_size=FONT_SUBTITLE - 1, z=10))
    out.append(make_border_box1(f'{seed_prefix}-header-accent', 262, 48, 280, 2, z=9))
    out.extend(make_breadcrumb(f'{seed_prefix}-header-crumb', 720, 18, breadcrumb_segments, z=20))
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
        bldg_action = _nav_action(building_page_id(bldg['sid']))
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
            floor_action = _nav_action(floor_page_id(bldg['sid'], floor['key']))
            fy = nav_y
            if fy + FLR_ROW_H > ALARM_PANEL_Y - 12:
                break
            flr_icon = '🔋' if floor['key'] == 'UPS设备组' else '📋'
            flr_label = truncate_label(floor['name'], 9)
            out.append(make_border_box1(f'{seed_prefix}-nav-flr-row-{bi}-{fi}', SIDEBAR_X + FLR_INDENT - 4, fy,
                                        SIDEBAR_W - FLR_INDENT - 8, FLR_ROW_H, z=3, action=floor_action))
            fty = vcenter(fy, FLR_ROW_H, 14)
            out.append(make_text(f'{seed_prefix}-nav-flr-icon-{bi}-{fi}', SIDEBAR_X + FLR_INDENT, fty, 14, 14, flr_icon,
                                 color=C_TEXT_MUTED, font_size=10, z=20, action=floor_action))
            out.append(make_text(f'{seed_prefix}-nav-flr-name-{bi}-{fi}', SIDEBAR_X + FLR_INDENT + 16, fty, 118, 14,
                                 flr_label, color=C_TEXT_MUTED, font_size=FONT_NAV_FLR - 1, z=20, action=floor_action))
            alarm_cnt = sum(1 for d in floor['devices'] if d['status'] != 1)
            badge_color = '#ef4444' if alarm_cnt else C_TEXT_DIM
            out.append(make_text(f'{seed_prefix}-nav-flr-badge-{bi}-{fi}', SIDEBAR_X + 210, fty, 58, 14,
                                 str(alarm_cnt), color=badge_color, font_size=10, z=20, action=floor_action))
            nav_y += FLR_ROW_H + FLR_ROW_GAP
        nav_y += BLDG_ROW_GAP
    alarm_y = ALARM_PANEL_Y
    out.append(make_border_box13(f'{seed_prefix}-alarm-panel', SIDEBAR_X + 8, alarm_y, SIDEBAR_W - 16, ALARM_PANEL_H, z=2))
    out.append(make_text(f'{seed_prefix}-alarm-title', SIDEBAR_X + 16, alarm_y + 10, 110, 20,
                         '🚨 实时告警', color=C_TEXT_MUTED, font_size=FONT_PANEL, z=8))
    out.append(make_text(f'{seed_prefix}-alarm-count', SIDEBAR_X + 150, alarm_y + 10, 40, 20,
                         '0', color='#ffffff', font_size=11, z=8))
    out.append(make_text(f'{seed_prefix}-alarm-empty', SIDEBAR_X + 16, alarm_y + 40, SIDEBAR_W - 32, ALARM_PANEL_H - 52,
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

# Subtle animated background accents
cells.append(make_decoration6('ov-bg-glow-l', MAIN_X, BODY_Y + 80, 420, 120, z=0))
cells.append(make_decoration6('ov-bg-glow-r', MAIN_X + MAIN_W - 420, BODY_Y + 80, 420, 120, z=0))

# Default KPI/chart context: first power meter (not UPS)
default_meter = None
for _b in buildings:
    for _f in _b['floors']:
        for _d in _f['devices']:
            if not _d['name'].startswith('UPS'):
                default_meter = _d
                break
        if default_meter:
            break
    if default_meter:
        break
if default_meter:
    set_ctx(default_meter)

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
    ('stat-power', '⚡', '0', 'kW', '总功率', '总有功功率' if '总有功功率' in CTX.dp_map else '输出总有功功率', '#667eea'),
    ('stat-energy', '📊', '---', 'kWh', '今日用电量', None, '#f093fb'),
    ('stat-online', '🖥', f'{online_count}/{TOTAL_DEVICES}', '', '在线设备', None, '#4facfe'),
    ('stat-alarm', '🔔', '0', '', '活跃告警', None, '#fa709a'),
]

for i, (seed, icon, val, unit, label, dp_name, accent) in enumerate(stat_configs):
    cx = card_xs[i]
    cells.append(make_border_box13(f'{seed}-bg', cx, stats_y, card_w, stats_h, z=2))
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
panel_h = 360
left_w = int((MAIN_W - 16) * 0.55)
right_w = MAIN_W - left_w - 16
right_x = MAIN_X + left_w + 16

cells.append(make_border_box13('panel-topology', MAIN_X, panel_top_y, left_w, panel_h, z=2))
cells.append(make_text('panel-topo-title', MAIN_X + 12, panel_top_y + 10, 400, 22,
                       '⚡ 系统拓扑概览', color=C_TEXT_MUTED, font_size=FONT_PANEL, z=6))
topo_y = panel_top_y + 44
topo_bottom = panel_top_y + panel_h - 12
for bi, bldg in enumerate(buildings):
    if topo_y + 40 > topo_bottom:
        break
    bldg_action = _nav_action(building_page_id(bldg['sid']))
    remaining = topo_bottom - topo_y
    max_floors = max(1, min(len(bldg['floors']), (remaining - 34) // 34))
    block_h = 34 + max_floors * 34
    cells.append(make_border_box1(f'topo-bldg-bg-{bi}', MAIN_X + 12, topo_y, left_w - 24, block_h, z=3,
                                  action=bldg_action))
    cells.append(make_text(f'topo-bldg-name-{bi}', MAIN_X + 24, topo_y + 8, left_w - 48, 18,
                           f'🏢 {bldg["name"]}', color=C_TEXT, font_size=13, z=6, action=bldg_action))
    floor_y = topo_y + 34
    for fi, floor in enumerate(bldg['floors'][:max_floors]):
        floor_action = _nav_action(floor_page_id(bldg['sid'], floor['key']))
        dot_count = min(len(floor['devices']), 6)
        dots = '🟢' * dot_count
        flr_icon = '🔋' if floor['key'] == 'UPS设备组' else '📋'
        line = f'{flr_icon} {truncate_label(floor["name"], 12)}  {dots}'
        cells.append(make_text(f'topo-flr-line-{bi}-{fi}', MAIN_X + 36, floor_y, left_w - 72, 16,
                               line, color=C_TEXT_MUTED, font_size=10, z=6, action=floor_action))
        floor_y += 34
    topo_y += block_h + 12

cells.append(make_border_box13('panel-chart', right_x, panel_top_y, right_w, panel_h, z=2))
cells.append(make_text('panel-chart-title', right_x + 12, panel_top_y + 10, 300, 22,
                       '📈 功率趋势 (24h)', color=C_TEXT_MUTED, font_size=FONT_PANEL, z=6))
cells.append(make_smooth_chart(
    'chart-trend', right_x + 12, panel_top_y + 38, right_w - 24, panel_h - 50,
    title='功率趋势 (24h)',
    dp_names=['总有功功率', '总无功功率', '总视在功率'] if not CTX.is_ups else ['输出总有功功率', '输出视在功率', '输出功率因数'],
    z=5
))

grid_y = panel_top_y + panel_h + 16
grid_h = 1080 - grid_y - 16
ups_devices = all_ups_devices()
ups_strip_h = 128 if ups_devices else 0

if ups_devices:
    cells.append(make_border_box13('panel-ups-strip', MAIN_X, grid_y, MAIN_W, ups_strip_h, z=2))
    cells.append(make_text('panel-ups-title', MAIN_X + 12, grid_y + 6, 520, 20,
                           f'🔋 UPS 不间断电源 ({len(ups_devices)}台)', color='#22d3ee', font_size=FONT_PANEL + 1, z=6))
    ups_floor_action = None
    for bldg in buildings:
        for floor in bldg['floors']:
            if floor['key'] == 'UPS设备组':
                ups_floor_action = _nav_action(floor_page_id(bldg['sid'], floor['key']))
                break
        if ups_floor_action:
            break
    cells.append(make_text('panel-ups-more', MAIN_X + MAIN_W - 130, grid_y + 8, 118, 22,
                           '查看全部 ›', color=C_ACCENT, font_size=12, z=20,
                           action=ups_floor_action))
    append_device_cards(cells, ups_devices, 'ups-card', MAIN_X, grid_y + 38, MAIN_W,
                        card_w=152, card_h=78, gap=10)

meter_grid_y = grid_y + ups_strip_h + (12 if ups_devices else 0)
meter_grid_h = grid_h - ups_strip_h - (12 if ups_devices else 0)
cells.append(make_border_box13('panel-device-grid', MAIN_X, meter_grid_y, MAIN_W, meter_grid_h, z=2))
cells.append(make_text('panel-grid-title', MAIN_X + 12, meter_grid_y + 10, 400, 22,
                       '⚡ 电力仪表运行状态', color=C_TEXT_MUTED, font_size=FONT_PANEL, z=6))

meter_devices = []
for bldg in buildings:
    for floor in bldg['floors']:
        for dev in floor['devices']:
            if dev_icon(dev) != '🔋':
                meter_devices.append(dev)

append_device_cards(cells, meter_devices, 'meter-card', MAIN_X, meter_grid_y + 36, MAIN_W,
                    card_w=148, card_h=88, gap=10, max_count=20)

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

cur.execute(
    "UPDATE display_model_layer SET components=?, updated_at=datetime('now') WHERE id=8",
    (comp_b64_main,),
)
conn.commit()
print(f"Database UPDATE executed for id=8 (overview), rows affected: {cur.rowcount}")

cur.execute(
    f"UPDATE display_model_layer SET layer='{LAYER_JSON}', updated_at=datetime('now') WHERE model_id=?",
    (MODEL_ID,),
)
conn.commit()
print(f"Layer background updated for id=8, rows affected: {cur.rowcount}")


# ════════════════════════════════════════════════════════
# LEVEL 1 & 2: PER-BUILDING / PER-FLOOR PAGES
# ════════════════════════════════════════════════════════

def build_building_detail_cells(bldg):
    out = []
    out.extend(build_header_cells(f'bldg-{bldg["sid"]}', [
        ('📊 全局总览', C_TEXT_DIM, PAGE_ID_MAIN),
        (bldg['name'], C_TEXT, None),
    ]))
    out.extend(build_sidebar_cells(f'bldg-nav-{bldg["sid"]}'))
    level_y = BODY_Y + 16
    out.append(make_text(f'bldg-{bldg["sid"]}-back', MAIN_X, level_y, 140, 32, '← 返回总览',
                         color=C_ACCENT, font_size=14, z=20, action=_nav_action(PAGE_ID_MAIN)))
    out.append(make_text(f'bldg-{bldg["sid"]}-title', MAIN_X + 160, level_y, 800, 28,
                         f'🏢 {bldg["name"]}', color=C_TEXT, font_size=22, z=10))
    alarm_cnt = sum(1 for f in bldg['floors'] for d in f['devices'] if d['status'] != 1)
    out.append(make_text(f'bldg-{bldg["sid"]}-sub', MAIN_X + 160, level_y + 34, 800, 18,
                         f'{bldg["device_count"]}台设备 · {alarm_cnt}条告警',
                         color=C_TEXT_DIM, font_size=13, z=10))
    card_w_item, card_h_item = 300, 140
    cards_per_row = max(1, (MAIN_W + 16) // (card_w_item + 16))
    card_start_x, card_start_y = MAIN_X, level_y + 72
    for fi, floor in enumerate(bldg['floors']):
        cx = card_start_x + (fi % cards_per_row) * (card_w_item + 16)
        cy = card_start_y + (fi // cards_per_row) * (card_h_item + 16)
        floor_nav = _nav_action(floor_page_id(bldg['sid'], floor['key']))
        running = sum(1 for d in floor['devices'] if d['status'] == 1)
        offline = floor['count'] - running
        icon = '🔋' if floor['key'] == 'UPS设备组' else '📋'
        out.append(make_border_box1(f'bldg-{bldg["sid"]}-card-{fi}', cx, cy, card_w_item, card_h_item, z=2, action=floor_nav))
        out.append(make_text(f'bldg-{bldg["sid"]}-card-name-{fi}', cx + 14, cy + 14, card_w_item - 28, 22,
                             f'{icon} {floor["name"]}', color=C_TEXT, font_size=15, z=5, action=floor_nav))
        out.append(make_text(f'bldg-{bldg["sid"]}-card-count-{fi}', cx + 14, cy + 40, card_w_item - 28, 18,
                             f'{floor["count"]}台设备', color=C_TEXT_DIM, font_size=12, z=5, action=floor_nav))
        out.append(make_text(f'bldg-{bldg["sid"]}-card-run-{fi}', cx + 14, cy + 72, 90, 18,
                             f'🟢 {running}运行', color=C_GREEN, font_size=11, z=5, action=floor_nav))
        out.append(make_text(f'bldg-{bldg["sid"]}-card-off-{fi}', cx + 110, cy + 72, 120, 18,
                             f'⏸ {offline}离线', color=C_TEXT_DIM, font_size=11, z=5, action=floor_nav))
    return out


def build_floor_detail_cells(bldg, floor):
    out = []
    out.extend(build_header_cells(f'flr-{bldg["sid"]}-{floor["key"]}', [
        ('📊 全局总览', C_TEXT_DIM, PAGE_ID_MAIN),
        (bldg['name'], C_TEXT_DIM, building_page_id(bldg['sid'])),
        (floor['name'], C_TEXT, None),
    ]))
    out.extend(build_sidebar_cells(f'flr-nav-{bldg["sid"]}-{floor["key"]}'))
    level_y = BODY_Y + 16
    out.append(make_text(f'flr-{bldg["sid"]}-back', MAIN_X, level_y, 160, 32, f'← {bldg["name"]}',
                         color=C_ACCENT, font_size=14, z=20, action=_nav_action(building_page_id(bldg['sid']))))
    out.append(make_text(f'flr-{bldg["sid"]}-title', MAIN_X + 180, level_y, 800, 28,
                         f'📋 {floor["name"]}', color=C_TEXT, font_size=22, z=10))
    out.append(make_text(f'flr-{bldg["sid"]}-sub', MAIN_X + 180, level_y + 40, 600, 18,
                         f'{floor["count"]}台设备', color=C_TEXT_DIM, font_size=13, z=10))
    table_y = level_y + 72
    row_h = 42
    col_widths = [50, 320, 100, 120, 140, 140, 120, 110, 120]
    col_headers = ['#', '设备名称', '状态', '实时功率', '电流A/B/C', '温度', '协议', '操作', '']
    col_starts = [MAIN_X]
    for w in col_widths[:-1]:
        col_starts.append(col_starts[-1] + w)
    out.append(make_border_box1(f'flr-{bldg["sid"]}-table', MAIN_X, table_y, MAIN_W, 900, z=1))
    out.append(make_border_box1(f'flr-{bldg["sid"]}-th', MAIN_X, table_y + 8, MAIN_W, row_h, z=3))
    for hi, (hdr, cs) in enumerate(zip(col_headers[:-1], col_starts)):
        if hdr:
            out.append(make_text(f'flr-{bldg["sid"]}-th-{hi}', cs + 8, table_y + 18, col_widths[hi] - 12, 20,
                                 hdr, color=C_TEXT_DIM, font_size=11, z=4))
    for di, dev in enumerate(floor['devices'][:20]):
        ry = table_y + 8 + row_h + di * row_h
        dev_action = _nav_device_action(dev)
        icon = dev_icon(dev)
        out.append(make_text(f'flr-{bldg["sid"]}-num-{di}', col_starts[0] + 8, ry + 10, 30, 20,
                             str(di + 1), color=C_TEXT_DIM, font_size=11, z=4))
        out.append(make_text(f'flr-{bldg["sid"]}-name-{di}', col_starts[1] + 8, ry + 10, col_widths[1] - 12, 20,
                             f'{icon} {dev["name"]}', color=C_ACCENT, font_size=13, z=4, action=dev_action))
        status_text = '运行中' if dev['status'] == 1 else '离线'
        status_color = C_GREEN if dev['status'] == 1 else C_TEXT_DIM
        out.append(make_text(f'flr-{bldg["sid"]}-stat-{di}', col_starts[2] + 8, ry + 10, col_widths[2] - 12, 20,
                             status_text, color=status_color, font_size=12, z=4))
        proto = 'UPS/Modbus' if icon == '🔋' else 'MODBUS'
        out.append(make_text(f'flr-{bldg["sid"]}-proto-{di}', col_starts[6] + 8, ry + 10, col_widths[6] - 12, 20,
                             proto, color=C_TEXT_DIM, font_size=12, z=4))
        out.append(make_border_box1(f'flr-{bldg["sid"]}-btn-bg-{di}', col_starts[7] + 8, ry + 4, 90, 28, z=4))
        out.append(make_text(f'flr-{bldg["sid"]}-btn-{di}', col_starts[7] + 14, ry + 8, 78, 20,
                             '详情', color=C_ACCENT, font_size=12, z=5, action=dev_action))
    return out


print("\n=== LEVEL 1 & 2: BUILDING / FLOOR PAGES ===")
building_count = floor_count = 0
for bldg in buildings:
    b64 = base64.b64encode(json.dumps({"cells": build_building_detail_cells(bldg)}, ensure_ascii=False).encode()).decode()
    upsert_page(f'building-{bldg["name"][:30]}', building_page_id(bldg['sid']), b64)
    building_count += 1
    for floor in bldg['floors']:
        f64 = base64.b64encode(json.dumps({"cells": build_floor_detail_cells(bldg, floor)}, ensure_ascii=False).encode()).decode()
        upsert_page(f'floor-{bldg["name"][:20]}-{floor["key"][:20]}', floor_page_id(bldg['sid'], floor['key']), f64)
        floor_count += 1
print(f"Created/updated {building_count} building pages, {floor_count} floor pages")
for old_pid in (PAGE_ID_BUILDING, PAGE_ID_FLOOR, '5100148b34ec5a609ce723e485a41a20'):
    cur.execute("UPDATE display_model_layer SET deleted_at=datetime('now') WHERE model_id=? AND page_id=?", (MODEL_ID, old_pid))
conn.commit()


# ════════════════════════════════════════════════════════
# LEVEL 3: PER-DEVICE DETAIL PAGES (one page_id per device)
# ════════════════════════════════════════════════════════

def build_device_detail_cells(dev, bldg, floor):
    set_ctx(dev)
    seed = dev['uuid'][:8]
    name = dev['name']
    bldg_name = bldg['name']
    floor_name = floor['name']
    is_ups = CTX.is_ups
    icon = dev_icon(dev)
    status_text = '运行中' if dev.get('status') == 1 else '离线'
    status_color = C_GREEN if dev.get('status') == 1 else C_TEXT_DIM
    dev_type = '施耐德 UPS' if is_ups else '多功能电力仪表'

    out = []
    out.extend(build_header_cells(f'dtl-{seed}', [
        ('📊 全局总览', C_TEXT_DIM, PAGE_ID_MAIN),
        (bldg_name, C_TEXT_DIM, building_page_id(bldg['sid'])),
        (floor_name, C_TEXT_DIM, floor_page_id(bldg['sid'], floor['key'])),
        (name, C_TEXT, None),
    ]))
    out.extend(build_sidebar_cells(f'dtl-nav-{seed}'))

    level_y = BODY_Y + 16
    out.append(make_text(f'dtl-{seed}-back', MAIN_X, level_y, 180, 32, f'← {floor_name}',
                         color=C_ACCENT, font_size=14, z=20,
                         action=_nav_action(floor_page_id(bldg['sid'], floor['key']))))
    out.append(make_decoration5(f'dtl-{seed}-title-deco', MAIN_X + 190, level_y - 4, 120, 40, z=1))
    out.append(make_text(f'dtl-{seed}-title', MAIN_X + 200, level_y, 720, 28,
                         f'{icon} {name}', color=C_TEXT, font_size=24, z=10))
    out.append(make_text(f'dtl-{seed}-status', MAIN_X + 940, level_y + 4, 220, 22,
                         status_text, color=status_color, font_size=14, z=10))

    panel_top = level_y + 64
    panel_h = 360
    col_w = (MAIN_W - 32) // 3
    left_x, mid_x, right_x = MAIN_X, MAIN_X + col_w + 16, MAIN_X + (col_w + 16) * 2

    out.append(make_border_box13(f'dtl-{seed}-left', left_x, panel_top, col_w, panel_h, z=2))
    out.append(make_text(f'dtl-{seed}-left-title', left_x + 15, panel_top + 8, 400, 22, '📋 基本参数',
                         color=C_ACCENT2, font_size=FONT_PANEL + 2, z=6))
    basic_params = [
        ('设备名称', name), ('设备类型', dev_type),
        ('通信协议', 'Modbus RTU'), ('所属机房', 'NCC 航信机房'),
        ('所属区域', floor_name), ('在线状态', status_text),
    ]
    for bi, (bk, bv) in enumerate(basic_params):
        by = panel_top + 44 + bi * 48
        out.append(make_text(f'dtl-{seed}-bk-{bi}', left_x + 15, by, 140, 22, bk,
                             color=C_TEXT_DIM, font_size=FONT_PARAM_LABEL + 1, z=6))
        out.append(make_text(f'dtl-{seed}-bv-{bi}', left_x + 160, by, col_w - 180, 22, bv,
                             color=C_TEXT, font_size=FONT_PARAM_VAL - 2, z=6))

    out.append(make_border_box13(f'dtl-{seed}-mid', mid_x, panel_top, col_w, panel_h, z=2))
    out.append(make_text(f'dtl-{seed}-mid-title', mid_x + 15, panel_top + 8, 400, 22, '📊 实时参数',
                         color=C_ACCENT2, font_size=FONT_PANEL + 2, z=6))
    if is_ups:
        rt_params = [
            ('输出A相电压', 'V'), ('输出B相电压', 'V'), ('输出C相电压', 'V'),
            ('输出A相电流', 'A'), ('输出B相电流', 'A'), ('输出C相电流', 'A'),
            ('输出频率', 'Hz'), ('电池电压', 'V'),
        ]
    else:
        rt_params = [
            ('AB线电压', 'V'), ('BC线电压', 'V'), ('CA线电压', 'V'),
            ('A相电流', 'A'), ('B相电流', 'A'), ('C相电流', 'A'),
            ('频率', 'Hz'), ('中性线电流', 'A'),
        ]
    for ri, (rname, runit) in enumerate(rt_params):
        if rname not in CTX.dp_map:
            continue
        ry = panel_top + 44 + ri * 36
        out.append(make_text(f'dtl-{seed}-rk-{ri}', mid_x + 15, ry, 150, 22, rname,
                             color=C_TEXT_DIM, font_size=FONT_PARAM_LABEL + 1, z=6))
        out.append(make_text(f'dtl-{seed}-rv-{ri}', mid_x + 170, ry, 120, 22, '---',
                             color=C_ACCENT, font_size=FONT_PARAM_VAL, z=6, data_bound=True, dp_name=rname))
        out.append(make_text(f'dtl-{seed}-ru-{ri}', mid_x + 300, ry, 60, 22, runit,
                             color=C_TEXT_DIM, font_size=FONT_PARAM_LABEL, z=6))

    out.append(make_border_box13(f'dtl-{seed}-right', right_x, panel_top, col_w, panel_h, z=2))
    title_right = '🔋 UPS 状态' if is_ups else '⚡ 功率参数'
    out.append(make_text(f'dtl-{seed}-right-title', right_x + 15, panel_top + 8, 400, 22, title_right,
                         color=C_ACCENT2, font_size=FONT_PANEL + 2, z=6))
    if is_ups:
        pw_params = [
            ('输出总有功功率', 'kW'), ('输出视在功率', 'kVA'), ('输出功率因数', ''),
            ('电池剩余运行时间', 'min'), ('电池温度', '°C'), ('UPS旁路状态', ''),
        ]
    else:
        pw_params = [
            ('总有功功率', 'kW'), ('总无功功率', 'kVar'), ('总视在功率', 'kVA'),
            ('总功率因数', ''), ('正有功电度', 'kWh'),
        ]
    pi_idx = 0
    for pname, punit in pw_params:
        if pname not in CTX.dp_map:
            continue
        py = panel_top + 44 + pi_idx * 44
        out.append(make_text(f'dtl-{seed}-pk-{pi_idx}', right_x + 15, py, 150, 22, pname,
                             color=C_TEXT_DIM, font_size=FONT_PARAM_LABEL + 1, z=6))
        out.append(make_text(f'dtl-{seed}-pv-{pi_idx}', right_x + 170, py, 140, 30, '---',
                             color=C_GREEN, font_size=FONT_KPI_VALUE - 4, z=6, data_bound=True, dp_name=pname))
        out.append(make_text(f'dtl-{seed}-pu-{pi_idx}', right_x + 320, py, 60, 22, punit,
                             color=C_TEXT_DIM, font_size=FONT_PARAM_LABEL, z=6))
        pi_idx += 1

    chart_y = panel_top + panel_h + 16
    chart_w = int(MAIN_W * 0.65)
    out.append(make_border_box13(f'dtl-{seed}-chart-bg', MAIN_X, chart_y, chart_w, 280, z=2))
    out.append(make_text(f'dtl-{seed}-chart-title', MAIN_X + 15, chart_y + 8, 400, 22, '📈 24小时趋势',
                         color=C_ACCENT2, font_size=FONT_PANEL + 2, z=6))
    chart_dps = (['输出总有功功率', '输出视在功率'] if is_ups else ['总有功功率', '总无功功率'])
    chart_dps = [d for d in chart_dps if d in CTX.dp_map]
    if chart_dps:
        out.append(make_smooth_chart(
            f'dtl-{seed}-chart', MAIN_X + 15, chart_y + 36, chart_w - 30, 230,
            title=f'{name} 趋势', dp_names=chart_dps, z=5,
        ))

    status_x = MAIN_X + chart_w + 16
    status_w = MAIN_W - chart_w - 16
    out.append(make_border_box13(f'dtl-{seed}-alarm-bg', status_x, chart_y, status_w, 280, z=2))
    out.append(make_text(f'dtl-{seed}-alarm-title', status_x + 15, chart_y + 8, 400, 22, '🔔 设备告警',
                         color=C_ACCENT2, font_size=FONT_PANEL + 2, z=6))
    out.append(make_text(f'dtl-{seed}-alarm-empty', status_x + 15, chart_y + 48, status_w - 30, 200,
                         '✅ 该设备无告警记录', color=C_TEXT_DIM, font_size=13, z=6))
    return out


print(f"\n=== LEVEL 3: PER-DEVICE DETAIL PAGES ===")
device_page_count = 0
device_meta = {}
for bldg in buildings:
    for floor in bldg['floors']:
        for dev in floor['devices']:
            pid = device_page_id(dev['uuid'])
            detail_cells = build_device_detail_cells(dev, bldg, floor)
            comp_json = json.dumps({"cells": detail_cells}, ensure_ascii=False)
            comp_b64 = base64.b64encode(comp_json.encode()).decode()
            page_name = f"device-{dev['name'][:40]}"
            layer_id, action = upsert_page(page_name, pid, comp_b64)
            device_page_count += 1
            device_meta[dev['name']] = pid
            if device_page_count <= 5 or dev['name'].startswith('UPS'):
                print(f"  {action} device page: {dev['name']} -> {pid[:16]}... (layer id={layer_id})")

# Remove legacy single shared device-detail page if present
cur.execute(
    "UPDATE display_model_layer SET deleted_at=datetime('now') WHERE model_id=? AND page_name='device-detail' AND page_id NOT IN (SELECT page_id FROM display_model_layer WHERE model_id=? AND page_name LIKE 'device-%')",
    (MODEL_ID, MODEL_ID),
)
# Simpler: soft-delete old generic page by known old page_id
OLD_DEVICE_PAGE = '5100148b34ec5a609ce723e485a41a20'
cur.execute(
    "UPDATE display_model_layer SET deleted_at=datetime('now') WHERE model_id=? AND page_id=?",
    (MODEL_ID, OLD_DEVICE_PAGE),
)

print(f"Created/updated {device_page_count} per-device detail pages")
conn.commit()


cur.execute("""
    SELECT id, page_name, is_home, page_id, LENGTH(components)
    FROM display_model_layer
    WHERE model_id=? AND deleted_at IS NULL
    ORDER BY is_home DESC, id
""", (MODEL_ID,))
pages = cur.fetchall()
print(f"\n=== All pages for model {MODEL_ID} ===")
for p in pages:
    print(f"  id={p[0]}, name={p[1]}, is_home={p[2]}, page_id={p[3][:20]}..., comp_len={p[4]}")

print(f"\n{'='*60}")
print("✅ Build complete!")
print(f"   Level 0 (overview):        {len(cells)} cells")
print(f"   Level 1 (building pages):  {building_count}")
print(f"   Level 2 (floor pages):     {floor_count}")
print(f"   Level 3 (device pages):    {device_page_count} pages")
print(f"\n   UPS device page examples:")
for ups_name in sorted(k for k in device_meta if k.startswith('UPS'))[:3]:
    print(f"     {ups_name} -> {device_meta[ups_name]}")
print(f"\n   Open: http://localhost:7080/#/AppRun/{MODEL_ID}")
print(f"{'='*60}")

conn.close()
