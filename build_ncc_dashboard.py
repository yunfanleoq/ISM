#!/usr/bin/env python3
"""
构建航信机房炫酷科技感大屏 v3 — 左侧导航 + 多层级钻探
Canvas: 1920×1080, bg #0a0e17

4 层页面:
  Level 0: overview (id=8)    — 230px 左侧导航树 + 主大屏面板
  Level 1: building-{sid}     — 每个柜独立 page（设备组卡片网格）
  Level 2: floor-{sid}-{key}  — 每个设备组独立 page（设备列表表格）
  Level 3: device-detail (id=10)— 91 cells 设备参数 + 趋势图 + 状态监控

布局网格 (overview):
  Left Sidebar:    x:0-230    w:230  (dv-border-box8 + 面包屑 + 设备树)
  Breadcrumb:      x:10, y:5, w:210
  Tree Nav:         inside sidebar (单行: 🏢 柜名 · N台 / 📋 组名 · N台)
  Header:          x:240, y:0-80, w:1680
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
# Per-entity page UUIDs (generated from sid/floor_key below)
def page_id_room(sid):
    return _uuid.uuid5(_uuid.NAMESPACE_DNS, f'ncc-dash-room-{sid}').hex

def page_id_building(sid):
    return _uuid.uuid5(_uuid.NAMESPACE_DNS, f'ncc-dash-bldg-{sid}').hex

def page_id_floor(bldg_sid, floor_key):
    return _uuid.uuid5(_uuid.NAMESPACE_DNS, f'ncc-dash-floor-{bldg_sid}-{floor_key}').hex

def page_id_device(dev_sid):
    return _uuid.uuid5(_uuid.NAMESPACE_DNS, f'ncc-dash-dev-{dev_sid}').hex
DEVICE_UUID = '68db26b1-113d-ad7e-79ff-10dbcc1c18d2'
DEVICE_NAME = '1A1_U11_S18_1'
DEV_MODEL_UUID = '3d734984-56f6-5494-ad4c-dfc67ca28ac8'
PROJECT_UUID = '31bc90be-ebc4-dd61-ba9d-ce6e075e40e2'

# Fetch data points for the sample (default) model
cur.execute(
    'SELECT name, uuid, data_unit FROM modbus_devices_data_model WHERE muid=%s ORDER BY id',
    (DEV_MODEL_UUID,)
)
dp_rows = cur.fetchall()
DP_MAP = {r[0]: {'uuid': r[1], 'unit': r[2] or ''} for r in dp_rows}

print(f"Data points ({len(DP_MAP)}): {list(DP_MAP.keys())}")

# Fetch data points for EVERY modbus model used by this project's devices, so
# per-device real-time binding resolves the correct point uuid for each device
# (avoids 张冠李戴 where every device showed the sample device's values).
cur.execute('SELECT DISTINCT muid FROM modbus_devices_data_model')
_all_muids = [r[0] for r in cur.fetchall()]
MODEL_DP = {}
for _m in _all_muids:
    cur.execute('SELECT name, uuid, data_unit FROM modbus_devices_data_model WHERE muid=%s ORDER BY id', (_m,))
    MODEL_DP[_m] = {r[0]: {'uuid': r[1], 'unit': r[2] or ''} for r in cur.fetchall()}


def dp_map_for(muid):
    """Return the point map for a device's model, falling back to the sample model."""
    return MODEL_DP.get(muid) or DP_MAP

# ── Query real device hierarchy ─────────────────────
cur.execute("""
    SELECT uuid, name, sid, pid, type, muid, status
    FROM monitor_list
    WHERE project_uuid = %s AND deleted_at IS NULL
    ORDER BY pid, type, name
""", (PROJECT_UUID,))
all_devices = cur.fetchall()

# 仅统计本项目 type=1 实际设备（航信机房 = 76，非全库 monitor_list 总数）
TOTAL_DEVICES = sum(1 for row in all_devices if row[4] == 1)
print(f"Project devices (type=1): {TOTAL_DEVICES}")

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
        'uuid': uuid, 'name': name, 'sid': sid, 'pid': pid,
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
                'uuid': uuid, 'name': name, 'sid': sid, 'pid': pid,
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
                'uuid': c0['uuid'], 'name': c0['name'], 'sid': c0['sid'], 'pid': building['sid'],
                'devices': type1_gc, 'device_count': len(type1_gc),
                'floors': [{'key': k, 'name': f'{k}设备组', 'devices': v, 'count': len(v)}
                           for k, v in sorted(floor_groups.items())]
            }
            buildings.append(sub)

# Remove root-level buildings that have no direct devices (they're just containers)
buildings = [b for b in buildings if b['device_count'] > 0]

# Assign unique page UUID per cabinet / device group (multi-page drill-down)
# and compute per-cabinet aggregate counts (scalable: counts, not per-point detail).
for b in buildings:
    b['page_id'] = page_id_building(b['sid'])
    b['online'] = sum(1 for f in b['floors'] for d in f['devices'] if d['status'] == 1)
    b['alarm'] = b['device_count'] - b['online']
    for f in b['floors']:
        f['page_id'] = page_id_floor(b['sid'], f['key'])
        f['online'] = sum(1 for d in f['devices'] if d['status'] == 1)

# ── ROOM / FLOOR level aggregation (真实层级: 机房 → 配电室 → 柜 → 设备) ──
# Cabinets are grouped by their real monitor_list parent (the 配电室/楼层 node),
# giving a top-level "zone" layer for the overview + a dedicated drill page.
ROOT_NODE = next((device_by_sid[s] for s in device_by_sid
                  if device_by_sid[s]['type'] == 0
                  and device_by_sid[s]['pid'] not in device_by_sid), None)
ROOT_NAME = ROOT_NODE['name'] if ROOT_NODE else '航信机房'

_room_map = {}
for b in buildings:
    rsid = b.get('pid')
    rnode = device_by_sid.get(rsid)
    rname = rnode['name'] if rnode else '配电区'
    room = _room_map.setdefault(rsid, {
        'sid': rsid, 'name': rname, 'page_id': page_id_room(rsid), 'cabinets': []
    })
    room['cabinets'].append(b)
    b['room_sid'] = rsid
    b['room_name'] = rname
    b['room_page_id'] = page_id_room(rsid)
rooms = list(_room_map.values())
for r in rooms:
    r['device_count'] = sum(c['device_count'] for c in r['cabinets'])
    r['online'] = sum(c['online'] for c in r['cabinets'])
    r['alarm'] = r['device_count'] - r['online']
    r['cabinet_count'] = len(r['cabinets'])

print(f"\n=== Aggregation ({ROOT_NAME}) ===")
for r in rooms:
    print(f"  Room: {r['name']} — {r['cabinet_count']}柜 / {r['device_count']}台 / 在线{r['online']} / 异常{r['alarm']}")
    for c in r['cabinets']:
        print(f"    Cabinet: {c['name']} — {c['device_count']}台 / 在线{c['online']} / 异常{c['alarm']} / {len(c['floors'])}组")

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

def _make_active(dp_name, device_uuid=None, device_name=None, dp_map=None):
    """Build a ShowData live-binding for a data point on a specific device.
    Defaults to the sample device/model when no device is supplied."""
    dp_map = dp_map if dp_map is not None else DP_MAP
    if dp_name not in dp_map:
        return []
    dp = dp_map[dp_name]
    return [{
        "id": "ShowData",
        "name": "configComponent.variable.ShowData",
        "result": "",
        "isExpression": False,
        "condition": {
            "deviceSN": device_uuid or DEVICE_UUID,
            "DeviceName": device_name or DEVICE_NAME,
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
SIDEBAR_W = 230
SIDEBAR_X = 0
MAIN_X = SIDEBAR_W + 10          # 240 — tighter sidebar/main gap
MAIN_W = 1920 - MAIN_X - 16      # 1664
BODY_Y = HEADER_H                # 56

BLDG_ROW_H = 34
BLDG_ROW_GAP = 2
FLR_ROW_H = 28
FLR_ROW_GAP = 2
FLR_INDENT = 20

# Neon-tech palette (deep-space blue base + cyan/blue/green/orange accents)
C_BG = '#0a0e17'          # deep space background
C_SIDEBAR = '#0b1322'     # sidebar fill (slightly translucent feel)
C_HEADER = '#0e1a2e'      # header bar fill
C_PANEL = '#0d1726'       # inner panel fill (sits inside neon frames)
C_PANEL_CARD = '#101d33'  # raised card fill
C_BORDER = '#1e3a5f'      # subtle separators
C_TEXT = '#e8f1ff'        # primary text
C_TEXT_MUTED = '#9fb6d6'  # secondary text
C_TEXT_DIM = '#5f7799'    # tertiary / labels
C_ACCENT = '#00e5ff'      # neon cyan — titles, links, highlights
C_BLUE = '#3b82f6'        # electric blue
C_GREEN = '#10e0a0'       # data green
C_ORANGE = '#ff6b35'      # alarm orange-red


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
    # Only real, visible glyphs count as "文字". Transparent/empty panel
    # backgrounds (shape view-svg-text with text="") are layout layers, not text.
    texts = [c for c in cells
             if c.get('shape') == 'view-svg-text'
             and (c['data']['detail']['style'].get('text', '') or '').strip()]
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

def make_panel_bg(seed, x, y, w, h, color=C_HEADER, z=0, opacity=1, action=None):
    """Flat panel background — no DataV corner SVG decorations."""
    cell_id = gen_uid(seed)
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
                "name": seed,
                "style": _make_style({"x": x, "y": y, "w": w, "h": h},
                                     text="", fontSize=1, foreColor="transparent",
                                     backColor=color, borderWidth=0, BorderEdges=4,
                                     opacity=opacity, diy=[]),
                "animate": _base_animate(),
                "action": action or [], "active": [], "dataBind": []
            }
        }
    }


def make_svg_time(seed, x, y, w, h, z=10, color=C_ACCENT, font_size=FONT_SUBTITLE,
                  time_format='YYYY/MM/DD HH:mm:ss', show_week=1):
    """Live clock via view-svg-time (updates every 500ms in AppRun)."""
    cell_id = gen_uid(seed)
    return {
        "shape": "view-svg-time",
        "id": cell_id, "x": x, "y": y, "width": w, "height": h,
        "zIndex": z, "visible": True,
        "position": {"x": x, "y": y},
        "size": {"width": w, "height": h},
        "data": {
            "detail": {
                "type": "view-svg-time",
                "identifier": cell_id,
                "name": "系统时间",
                "style": _make_style({"x": x, "y": y, "w": w, "h": h},
                                     text="", fontSize=font_size, foreColor=color,
                                     backColor="transparent", borderWidth=0, BorderEdges=0,
                                     fontFamily="Courier New", diy=[
                                         {"name": "component.public.fillOpacity", "type": 7,
                                          "value": 1, "min": 0, "max": 1, "key": "fillOpacity"},
                                         {"name": "configComponent.time.IsShowWeek", "type": 6,
                                          "value": show_week,
                                          "enumList": [{"value": 1, "option": "Yes"},
                                                       {"value": 0, "option": "No"}],
                                          "key": "IsShowWeek"},
                                         {"name": "configComponent.time.TimeFormat", "type": 4,
                                          "value": time_format, "key": "TimeFormat"},
                                     ]),
                "animate": _base_animate(),
                "action": [], "active": [], "dataBind": []
            }
        }
    }


def make_text(seed, x, y, w, h, text, color='#c8d6e5', font_size=14,
              z=10, data_bound=False, dp_name=None, action=None, device=None):
    cell_id = gen_uid(seed)
    if data_bound and dp_name:
        dev = device or {}
        active = _make_active(dp_name, dev.get('uuid'), dev.get('name'),
                              dp_map_for(dev.get('muid')) if dev else None)
    else:
        active = []
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

# ── DataV decorative frames (controlled "decoration budget") ──
# These wrappers all share one builder. The underlying DataV components only
# read detail.style.position.{w,h} to size themselves and render their own
# neon-cyan border/decoration SVG; extra style keys are ignored harmlessly.

def make_dv_frame(shape, seed, x, y, w, h, z=1, diy=None):
    """Generic DataV border-box / decoration cell. Non-interactive frame layer."""
    cell_id = gen_uid(seed)
    style = {"position": {"x": x, "y": y, "w": w, "h": h}, "visible": 1}
    if diy is not None:
        style["diy"] = diy
    return {
        "shape": shape,
        "id": cell_id, "x": x, "y": y, "width": w, "height": h,
        "zIndex": z, "visible": True,
        "position": {"x": x, "y": y},
        "size": {"width": w, "height": h},
        "data": {
            "detail": {
                "type": shape,
                "identifier": cell_id,
                "name": seed,
                "style": style,
                "animate": _base_animate(),
                "action": [], "active": [], "dataBind": []
            }
        }
    }


# box8 = slow rotating corner light (use sparingly, 1 per page max)
_BOX8_DIY = [
    {"name": "border89cur", "type": 1, "value": 8, "min": 1, "key": "border89cur"},
    {"name": "border89Direction", "type": 6, "value": 0, "min": 1, "key": "border89Direction",
     "enumList": [{"value": 0, "option": "Forward"}, {"value": 1, "option": "Negative"}]},
]

def make_box8(seed, x, y, w, h, z=1):
    return make_dv_frame("dv-border-box8", seed, x, y, w, h, z=z, diy=_BOX8_DIY)

def make_box12(seed, x, y, w, h, z=2):
    """Clean rounded tech frame with corner glow — ideal for KPI cards."""
    return make_dv_frame("dv-border-box12", seed, x, y, w, h, z=z)

def make_box13(seed, x, y, w, h, z=1):
    """Panel frame with top-left title bracket — ideal for main content panels."""
    return make_dv_frame("dv-border-box13", seed, x, y, w, h, z=z)

def make_decoration1(seed, x, y, w, h, z=3):
    """Animated flowing squares — use as a key horizontal divider."""
    return make_dv_frame("dv-decoration1", seed, x, y, w, h, z=z)

def make_decoration8(seed, x, y, w, h, z=3):
    """Angled tech title underline."""
    return make_dv_frame("dv-decoration8", seed, x, y, w, h, z=z)


def make_hud_corners(seed, x, y, w, h, arm=46, thick=2, color=C_ACCENT,
                     opacity=0.55, z=1):
    """Self-drawn L-shaped neon corner brackets (no animation, no frame-in-frame).
    Eight thin bars form 4 right-angle corners — a stable command-center HUD edge."""
    out = []
    x2, y2 = x + w, y + h
    # (bar_x, bar_y, bar_w, bar_h) for each of the 8 segments
    segs = [
        (x, y, arm, thick), (x, y, thick, arm),                       # top-left
        (x2 - arm, y, arm, thick), (x2 - thick, y, thick, arm),       # top-right
        (x, y2 - thick, arm, thick), (x, y2 - arm, thick, arm),       # bottom-left
        (x2 - arm, y2 - thick, arm, thick), (x2 - thick, y2 - arm, thick, arm),  # bottom-right
    ]
    for i, (bx, by, bw, bh) in enumerate(segs):
        out.append(make_panel_bg(f'{seed}-hud-{i}', bx, by, bw, bh,
                                 color=color, z=z, opacity=opacity))
    return out


def build_screen_decor(seed_prefix):
    """Restrained full-screen HUD layer shared by every page:
    low-contrast neon corner brackets + a soft title glow band. z kept low so
    interactive text (z>=4) always wins clicks; nothing blinks."""
    out = []
    # Screen-edge HUD corners (inset 6px, behind everything)
    out.extend(make_hud_corners(f'{seed_prefix}-screen', 6, 6, 1920 - 12, 1080 - 12,
                                arm=52, thick=2, color=C_ACCENT, opacity=0.45, z=1))
    # Title glow band: a wide, dim cyan halo behind the header title (static)
    out.append(make_panel_bg(f'{seed_prefix}-title-glow', 256, 6, 470, 44,
                             color=C_ACCENT, z=1, opacity=0.07))
    return out


def make_panel_title(seed, x, y, text, color=C_ACCENT, font_size=FONT_PANEL, z=6, w=420):
    """Section title with a neon accent bar — flat & cheap (not a DataV component)."""
    out = []
    bar_h = font_size + 4
    out.append(make_panel_bg(f'{seed}-bar', x, y + 1, 4, bar_h, color=color, z=z))
    out.append(make_text(f'{seed}-txt', x + 12, y, w, bar_h + 2, text,
                         color=color, font_size=font_size, z=z))
    return out

def make_smooth_chart(seed, x, y, w, h, title, dp_names, z=5, device=None):
    cell_id = gen_uid(seed)
    active = []
    dev = device or {}
    d_uuid = dev.get('uuid') or DEVICE_UUID
    d_name = dev.get('name') or DEVICE_NAME
    d_map = dp_map_for(dev.get('muid')) if dev else DP_MAP
    var_ids = ['ShowChartVariable1', 'ShowChartVariable2', 'ShowChartVariable3',
               'ShowChartVariable4', 'ShowChartVariable5']
    for i, dpn in enumerate(dp_names[:5]):
        if dpn not in d_map:
            continue
        dp = d_map[dpn]
        active.append({
            "id": var_ids[i],
            "name": "configComponent.variable.ShowData",
            "result": "",
            "isExpression": False,
            "condition": {
                "deviceSN": d_uuid, "DeviceName": d_name,
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
    out.append(make_panel_bg(f'{seed_prefix}-header-bg', 0, 0, 1920, HEADER_H, color=C_HEADER, z=0))
    # Shared restrained HUD decoration (screen-edge corners + title glow), behind text.
    out.extend(build_screen_decor(seed_prefix))
    out.append(make_text(f'{seed_prefix}-header-logo', 220, 10, 36, 36, '⚡',
                         color=C_ACCENT, font_size=26, z=10))
    out.append(make_text(f'{seed_prefix}-header-title', 260, 3, 460, 24, '航信机房电力监控系统',
                         color=C_TEXT, font_size=FONT_TITLE, z=10))
    out.append(make_text(f'{seed_prefix}-header-subtitle', 262, 35, 460, 14, 'NCC ROOM POWER SCADA MONITOR',
                         color=C_TEXT_DIM, font_size=FONT_SUBTITLE, z=10))
    out.extend(make_breadcrumb(f'{seed_prefix}-header-crumb', 760, 20, breadcrumb_segments, z=20))
    out.append(make_svg_time(f'{seed_prefix}-header-clock', 1480, 16, 280, 24,
                             z=10, color=C_ACCENT, font_size=FONT_SUBTITLE,
                             time_format='YYYY/MM/DD HH:mm:ss', show_week=1))
    out.append(make_text(f'{seed_prefix}-header-status', 1770, 18, 130, 22, '🟢 在线',
                         color=C_GREEN, font_size=FONT_SUBTITLE, z=10))
    # STATIC HUD divider line (no marquee animation): a dim full-width rule with a
    # brighter cyan accent segment under the title -> tech feel, zero flicker.
    out.append(make_panel_bg(f'{seed_prefix}-header-rule', MAIN_X, HEADER_H - 1, MAIN_W, 1,
                             color=C_BORDER, z=3, opacity=0.9))
    out.append(make_panel_bg(f'{seed_prefix}-header-accent', MAIN_X, HEADER_H - 2, 220, 2,
                             color=C_ACCENT, z=4, opacity=0.9))
    return out


def build_sidebar_cells(seed_prefix='nav'):
    """230px sidebar: compact single-line tree rows + alarm panel."""
    out = []
    inner_w = SIDEBAR_W - 12
    out.append(make_panel_bg(f'{seed_prefix}-sidebar-frame', SIDEBAR_X + 2, BODY_Y + 2,
                             SIDEBAR_W - 4, 1080 - BODY_Y - 4, color=C_SIDEBAR, z=0))
    # Single STATIC neon frame around the whole sidebar (1 box per zone).
    # box13 has no self-animation (unlike box8's rotating marquee) -> no flicker.
    # z=1 keeps it behind the z=20 nav text, so drill-down clicks still land.
    out.append(make_box13(f'{seed_prefix}-sidebar-box', SIDEBAR_X + 2, BODY_Y + 2,
                          SIDEBAR_W - 4, 1080 - BODY_Y - 6, z=1))
    nav_y = BODY_Y + 12
    out.append(make_text(f'{seed_prefix}-nav-title', SIDEBAR_X + 10, nav_y, inner_w, 20, '📍 层级导航',
                         color=C_TEXT_MUTED, font_size=FONT_PANEL, z=8))
    nav_y += 28
    out.append(make_text(f'{seed_prefix}-nav-root', SIDEBAR_X + 10, nav_y, inner_w, 18,
                         f'🏭 {ROOT_NAME} · {TOTAL_DEVICES}台',
                         color=C_ACCENT, font_size=FONT_NAV_BLDG, z=20,
                         action=_nav_action(PAGE_ID_MAIN)))
    nav_y += 22
    # Hierarchical, bounded nav: 机房 → 配电室(room) → 机柜(cabinet).
    # Device groups are NOT listed here (reachable by drilling) so the sidebar
    # stays bounded under 2万+ points. Rows are capped per level.
    NAV_ROOM_CAP = 6
    NAV_CAB_CAP = 8
    NAV_LIMIT_Y = 820
    for ri, room in enumerate(rooms[:NAV_ROOM_CAP]):
        if nav_y > NAV_LIMIT_Y:
            break
        room_action = _nav_action(room['page_id'])
        ty = vcenter(nav_y, BLDG_ROW_H, 18)
        out.append(make_text(f'{seed_prefix}-nav-room-{ri}', SIDEBAR_X + 10, ty, inner_w - 8, 18,
                             f"🏛 {room['name']} · {room['device_count']}台",
                             color=C_TEXT, font_size=FONT_NAV_BLDG, z=20, action=room_action))
        nav_y += BLDG_ROW_H
        for ci, cab in enumerate(room['cabinets'][:NAV_CAB_CAP]):
            if nav_y > NAV_LIMIT_Y:
                break
            cab_action = _nav_action(cab['page_id'])
            flr_w = SIDEBAR_W - FLR_INDENT - 10
            fty = vcenter(nav_y, FLR_ROW_H, 16)
            mark = '🟢' if cab['alarm'] == 0 else '🔴'
            out.append(make_text(f'{seed_prefix}-nav-cab-{ri}-{ci}', SIDEBAR_X + FLR_INDENT, fty, flr_w - 8, 16,
                                 f"{mark} {cab['name']} · {cab['device_count']}台",
                                 color=C_TEXT_MUTED, font_size=FONT_NAV_FLR, z=20, action=cab_action))
            nav_y += FLR_ROW_H + FLR_ROW_GAP
        if len(room['cabinets']) > NAV_CAB_CAP:
            out.append(make_text(f'{seed_prefix}-nav-cabmore-{ri}', SIDEBAR_X + FLR_INDENT, nav_y, flr_w - 8, 16,
                                 f"+{len(room['cabinets']) - NAV_CAB_CAP} 更多 ›", color=C_TEXT_DIM,
                                 font_size=FONT_NAV_FLR, z=20, action=room_action))
            nav_y += FLR_ROW_H + FLR_ROW_GAP
        nav_y += BLDG_ROW_GAP
    alarm_y = max(nav_y + 16, 858)
    out.append(make_panel_bg(f'{seed_prefix}-alarm-panel', SIDEBAR_X + 10, alarm_y, SIDEBAR_W - 20, 196,
                             color=C_PANEL_CARD, z=2, opacity=0.85))
    out.extend(make_panel_title(f'{seed_prefix}-alarm', SIDEBAR_X + 16, alarm_y + 10,
                                '实时告警', color=C_ORANGE, font_size=FONT_PANEL, z=8, w=110))
    out.append(make_text(f'{seed_prefix}-alarm-count', SIDEBAR_X + 150, alarm_y + 10, 40, 20,
                         '0', color=C_TEXT, font_size=11, z=8))
    out.append(make_text(f'{seed_prefix}-alarm-empty', SIDEBAR_X + 16, alarm_y + 44, SIDEBAR_W - 36, 140,
                         '✅ 当前无告警', color=C_TEXT_DIM, font_size=12, z=8))
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
online_count = sum(1 for row in all_devices if row[4] == 1 and row[6] == 1)
alarm_count = TOTAL_DEVICES - online_count

stats_y = BODY_Y + 16
stats_h = 100
card_w = int((MAIN_W - 48) / 4)
card_gap = 16
card_xs = [MAIN_X + i * (card_w + card_gap) for i in range(4)]

stat_configs = [
    ('stat-power', '⚡', '0', 'kW', '总功率', '总有功功率', C_ACCENT),
    ('stat-energy', '📊', '---', 'kWh', '今日用电量', None, C_BLUE),
    ('stat-online', '🖥', f'{online_count}/{TOTAL_DEVICES}', '', '在线设备', None, C_GREEN),
    ('stat-alarm', '🔔', str(alarm_count), '', '活跃告警', None, C_ORANGE),
]

for i, (seed, icon, val, unit, label, dp_name, accent) in enumerate(stat_configs):
    cx = card_xs[i]
    # depth: faint inner fill (z=1) under the glowing box12 frame (z=2)
    cells.append(make_panel_bg(f'{seed}-fill', cx + 5, stats_y + 5, card_w - 10, stats_h - 10,
                               color=C_PANEL_CARD, z=1, opacity=0.55))
    cells.append(make_box12(f'{seed}-bg', cx, stats_y, card_w, stats_h, z=2))
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

# Topo panel — HIERARCHICAL drill (从大到小): 配电室/楼层 大节点 → 其下机柜（含UPS柜）。
# Scalable: top level only renders ROOM nodes (capped); cabinets are aggregate rows
# (count + status), never per-device. Click a room → room page; click a cabinet →
# cabinet page. Counts come from real monitor_list aggregation, so 2万+ points still
# render as a handful of nodes here.
cells.append(make_box13('panel-topo-frame', MAIN_X, panel_top_y, left_w, panel_h, z=1))
cells.extend(make_panel_title('panel-topo', MAIN_X + 16, panel_top_y + 10,
                              f'拓扑概览 · {ROOT_NAME}', color=C_ACCENT, font_size=FONT_PANEL, z=6))
topo_y = panel_top_y + 44
topo_limit = panel_top_y + panel_h - 12
CAB_LINE_H = 26
TOPO_ROOM_CAP = 6
TOPO_CAB_CAP = 6
for ri, room in enumerate(rooms[:TOPO_ROOM_CAP]):
    avail = topo_limit - topo_y - 34
    if avail < CAB_LINE_H:
        break
    cab_show = min(len(room['cabinets']), max(1, avail // CAB_LINE_H), TOPO_CAB_CAP)
    block_h = 34 + cab_show * CAB_LINE_H
    room_action = _nav_action(room['page_id'])
    cells.append(make_panel_bg(f'topo-room-bg-{ri}', MAIN_X + 16, topo_y, left_w - 32, block_h,
                               color=C_PANEL_CARD, z=3, opacity=0.85, action=room_action))
    cells.append(make_text(f'topo-room-name-{ri}', MAIN_X + 28, topo_y + 8, 260, 18,
                           f'🏛 {room["name"]}', color=C_ACCENT, font_size=13, z=6, action=room_action))
    cells.append(make_text(f'topo-room-agg-{ri}', MAIN_X + 300, topo_y + 9, left_w - 332, 16,
                           f'{room["cabinet_count"]}柜 · {room["device_count"]}台 · 在线{room["online"]} · 异常{room["alarm"]}',
                           color=C_TEXT_MUTED, font_size=11, z=6, action=room_action))
    cab_y = topo_y + 34
    for ci, cab in enumerate(room['cabinets'][:cab_show]):
        cab_action = _nav_action(cab['page_id'])
        cells.append(make_text(f'topo-cab-name-{ri}-{ci}', MAIN_X + 44, cab_y + 3, 220, 16,
                               f'🗄 {cab["name"]}', color=C_TEXT, font_size=11, z=6, action=cab_action))
        dot_color = C_GREEN if cab['alarm'] == 0 else C_ORANGE
        cells.append(make_text(f'topo-cab-stat-{ri}-{ci}', MAIN_X + 300, cab_y + 4, left_w - 340, 14,
                               f'{cab["device_count"]}台 · 在线{cab["online"]}/{cab["device_count"]}',
                               color=dot_color, font_size=10, z=6, action=cab_action))
        cab_y += CAB_LINE_H
    if len(room['cabinets']) > cab_show:
        cells.append(make_text(f'topo-cab-more-{ri}', MAIN_X + 44, cab_y - CAB_LINE_H + 3, 220, 14,
                               f'+{len(room["cabinets"]) - cab_show} 更多机柜 ›', color=C_TEXT_DIM,
                               font_size=10, z=6, action=room_action))
    topo_y += block_h + 12
if len(rooms) > TOPO_ROOM_CAP:
    cells.append(make_text('topo-room-more', MAIN_X + 28, topo_limit - 18, left_w - 56, 16,
                           f'+{len(rooms) - TOPO_ROOM_CAP} 更多区域 · 共 {len(rooms)} 区',
                           color=C_TEXT_DIM, font_size=11, z=6))

# Chart panel — 1 neon frame (box13) + accent title; chart inset on top
cells.append(make_box13('panel-chart-frame', right_x, panel_top_y, right_w, panel_h, z=1))
cells.extend(make_panel_title('panel-chart', right_x + 16, panel_top_y + 10,
                              '功率趋势 (24h)', color=C_ACCENT, font_size=FONT_PANEL, z=6))
cells.append(make_smooth_chart(
    'chart-trend', right_x + 14, panel_top_y + 40, right_w - 28, panel_h - 52,
    title='功率趋势 (24h)',
    dp_names=['总有功功率', '总无功功率', '总视在功率'],
    z=5
))

grid_y = panel_top_y + panel_h + 16
grid_h = 1080 - grid_y - 16
# Device-overview panel — HIERARCHICAL AGGREGATION (按区域/机柜聚合，不平铺设备)。
# Each card = one cabinet's rolled-up stats (台数/在线/异常 + 状态条)。Per-device
# detail lives only at the leaf device page. Cards are capped (TopN) so the panel
# stays bounded under 2万+ points; overflow shows a "+N 更多 / 共M柜" footer.
cells.append(make_box13('panel-grid-frame', MAIN_X, grid_y, MAIN_W, grid_h, z=1))
cells.extend(make_panel_title('panel-grid', MAIN_X + 16, grid_y + 8,
                              '设备运行总览 · 按区域聚合', color=C_ACCENT, font_size=FONT_PANEL, z=6, w=240))
cells.append(make_text('panel-grid-hint', MAIN_X + 280, grid_y + 10, 560, 18,
                       '点击卡片逐级下钻：区域 › 机柜 › 设备组 › 设备', color=C_TEXT_DIM,
                       font_size=11, z=6))

agg_card_w = 392
agg_card_h = 132
agg_gap = 16
agg_cols = max(1, (MAIN_W - 32 + agg_gap) // (agg_card_w + agg_gap))
agg_x0 = MAIN_X + 16
agg_y = grid_y + 40
GRID_CARD_CAP = 9            # bounded node budget for the overview
flat_cabs = [(r, c) for r in rooms for c in r['cabinets']]
shown = flat_cabs[:GRID_CARD_CAP]
for idx, (room, cab) in enumerate(shown):
    rrow = idx // agg_cols
    ccol = idx % agg_cols
    cx = agg_x0 + ccol * (agg_card_w + agg_gap)
    cy = agg_y + rrow * (agg_card_h + agg_gap)
    if cy + agg_card_h > grid_y + grid_h - 28:
        break
    cab_action = _nav_action(cab['page_id'])
    ok = cab['alarm'] == 0
    accent = C_GREEN if ok else C_ORANGE
    cells.append(make_panel_bg(f'agg-bg-{idx}', cx, cy, agg_card_w, agg_card_h,
                               color=C_PANEL_CARD, z=3, opacity=0.9, action=cab_action))
    cells.append(make_panel_bg(f'agg-bar-{idx}', cx, cy, 4, agg_card_h, color=accent, z=4))
    txt_w = agg_card_w - 32 - 96      # leave room for the "N组 ›" badge on the right
    cells.append(make_text(f'agg-room-{idx}', cx + 16, cy + 10, txt_w, 14,
                           f'🏛 {room["name"]}', color=C_TEXT_DIM, font_size=11, z=6, action=cab_action))
    cells.append(make_text(f'agg-name-{idx}', cx + 16, cy + 32, txt_w, 24,
                           f'🗄 {cab["name"]}', color=C_ACCENT, font_size=16, z=6, action=cab_action))
    # three aggregate metrics (count / online / abnormal) — numbers, not per-device
    metrics = [('设备', cab['device_count'], C_TEXT), ('在线', cab['online'], C_GREEN),
               ('异常', cab['alarm'], C_ORANGE if cab['alarm'] else C_TEXT_DIM)]
    mw = (agg_card_w - 32) // 3
    for mi, (mlabel, mval, mcolor) in enumerate(metrics):
        mx = cx + 16 + mi * mw
        cells.append(make_text(f'agg-mv-{idx}-{mi}', mx, cy + 64, mw - 8, 28,
                               str(mval), color=mcolor, font_size=24, z=6, action=cab_action))
        cells.append(make_text(f'agg-ml-{idx}-{mi}', mx, cy + 100, mw - 8, 18,
                               mlabel, color=C_TEXT_DIM, font_size=11, z=6, action=cab_action))
    cells.append(make_text(f'agg-grp-{idx}', cx + agg_card_w - 96, cy + 12, 84, 16,
                           f'{len(cab["floors"])}组 ›', color=C_TEXT_DIM, font_size=11, z=6, action=cab_action))
if len(flat_cabs) > len(shown):
    cells.append(make_text('agg-more', agg_x0, grid_y + grid_h - 24, MAIN_W - 32, 18,
                           f'+{len(flat_cabs) - len(shown)} 更多机柜 · 共 {len(flat_cabs)} 柜 / {TOTAL_DEVICES} 台 · 逐级下钻查看明细',
                           color=C_TEXT_DIM, font_size=12, z=6))

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

# Locate the overview/home row robustly. After a snapshot restore the row id is
# NOT guaranteed to be 8, so resolve by (model_id + home page) instead of id=8.
# autoSize=1 → AppRun fit-scales the 1920×1080 canvas to the viewport (setScale).
# With autoSize=0 the canvas renders at native 1920px and the right edge gets
# clipped on any window narrower than 1920 — the reported "右侧被切" root cause.
_dark_layer = '{"width":1920,"height":1080,"autoSize":1,"Padding":0,"gridSize":10,"background":"#0a0e17"}'
cur.execute(
    """SELECT id FROM display_model_layer
       WHERE model_id=%s AND deleted_at IS NULL AND (page_id=%s OR is_home=1)
       ORDER BY is_home DESC, id LIMIT 1""",
    (MODEL_ID, PAGE_ID_MAIN)
)
_ov = cur.fetchone()
if _ov:
    overview_id = _ov[0]
    cur.execute("UPDATE display_model_layer SET components=%s, layer=%s, updated_at=NOW() WHERE id=%s",
                (comp_b64_main, _dark_layer, overview_id))
    conn.commit()
    print(f"Database UPDATE executed for overview id={overview_id}, rows affected: {cur.rowcount}")
else:
    cur.execute(
        """INSERT INTO display_model_layer
           (model_id, page_name, page_id, is_home, is_login, page_type, layer, components, created_at, updated_at)
           VALUES (%s, 'main', %s, 1, 0, 1, %s, %s, NOW(), NOW())""",
        (MODEL_ID, PAGE_ID_MAIN, _dark_layer, comp_b64_main)
    )
    conn.commit()
    print(f"Inserted overview/home page id={cur.lastrowid}")


# ════════════════════════════════════════════════════════
# LEVEL 1 & 2: PER-CABINET / PER-GROUP PAGES
# ════════════════════════════════════════════════════════

def build_building_detail_cells(bldg, seed_prefix='bldg'):
    """One page per cabinet — shows device-group cards for that cabinet only."""
    out = []
    out.extend(build_header_cells(seed_prefix, [
        ('📊 全局总览', C_TEXT_DIM, PAGE_ID_MAIN),
        (bldg.get('room_name', ROOT_NAME), C_TEXT_DIM, bldg.get('room_page_id')),
        (bldg['name'], C_TEXT, None),
    ]))
    out.extend(build_sidebar_cells(f'{seed_prefix}-nav'))
    back_to_room = _nav_action(bldg.get('room_page_id') or PAGE_ID_MAIN)
    level_y = BODY_Y + 16
    # Single neon frame around the cabinet content panel (1 box per zone)
    out.append(make_box13(f'{seed_prefix}-content-frame', MAIN_X, level_y - 4,
                          MAIN_W, 1080 - (level_y - 4) - 16, z=1))
    out.append(make_text(f'{seed_prefix}-back-btn', MAIN_X + 16, level_y, 180, 32,
                         f'← {bldg.get("room_name", "返回")}',
                         color=C_ACCENT, font_size=14, z=20, action=back_to_room))
    out.append(make_text(f'{seed_prefix}-level-title', MAIN_X + 214, level_y, 760, 28,
                         f'🏢 {bldg["name"]}', color=C_TEXT, font_size=22, z=10))
    alarm_cnt = sum(1 for f in bldg['floors'] for d in f['devices'] if d['status'] != 1)
    out.append(make_text(f'{seed_prefix}-level-sub', MAIN_X + 214, level_y + 40, 760, 18,
                         f'{bldg["device_count"]}台设备 · {alarm_cnt}条异常',
                         color=C_TEXT_DIM, font_size=13, z=10))
    card_start_x = MAIN_X + 16
    card_start_y = level_y + 76
    card_w_item = 300
    card_h_item = 140
    cards_per_row = max(1, (MAIN_W - 32 + 16) // (card_w_item + 16))
    card_gap_x = 16
    card_gap_y = 16
    for fi, floor in enumerate(bldg['floors']):
        row = fi // cards_per_row
        col = fi % cards_per_row
        cx = card_start_x + col * (card_w_item + card_gap_x)
        cy = card_start_y + row * (card_h_item + card_gap_y)
        floor_nav = _nav_action(floor['page_id'])
        running = sum(1 for d in floor['devices'] if d['status'] == 1)
        offline = floor['count'] - running
        out.append(make_panel_bg(f'{seed_prefix}-card-bg-{fi}', cx, cy, card_w_item, card_h_item,
                                 color=C_PANEL_CARD, z=3, opacity=0.9, action=floor_nav))
        out.append(make_text(f'{seed_prefix}-card-name-{fi}', cx + 14, cy + 12, card_w_item - 28, 22,
                             f'📋 {floor["name"]}', color=C_ACCENT, font_size=15, z=5, action=floor_nav))
        out.append(make_text(f'{seed_prefix}-card-count-{fi}', cx + 14, cy + 46, card_w_item - 28, 18,
                             f'{floor["count"]}台设备', color=C_TEXT_DIM, font_size=12, z=5, action=floor_nav))
        out.append(make_text(f'{seed_prefix}-card-run-{fi}', cx + 14, cy + 72, 90, 18,
                             f'🟢 {running}运行', color=C_GREEN, font_size=11, z=5, action=floor_nav))
        out.append(make_text(f'{seed_prefix}-card-alarm-{fi}', cx + 110, cy + 72, 80, 18,
                             f'🔴 0告警', color='#ef4444', font_size=11, z=5, action=floor_nav))
        out.append(make_text(f'{seed_prefix}-card-stop-{fi}', cx + 200, cy + 72, 80, 18,
                             f'⏸ {offline}停止', color=C_TEXT_DIM, font_size=11, z=5, action=floor_nav))
    return out


def build_floor_detail_cells(bldg, floor, seed_prefix='floor'):
    """One page per device group — shows Modbus device list for that group only."""
    out = []
    out.extend(build_header_cells(seed_prefix, [
        ('📊 全局总览', C_TEXT_DIM, PAGE_ID_MAIN),
        (bldg.get('room_name', ROOT_NAME), C_TEXT_DIM, bldg.get('room_page_id')),
        (bldg['name'], C_TEXT_DIM, bldg['page_id']),
        (floor['name'], C_TEXT, None),
    ]))
    out.extend(build_sidebar_cells(f'{seed_prefix}-nav'))
    back_to_bldg = _nav_action(bldg['page_id'])
    level_y = BODY_Y + 16
    # Single neon frame around the device-list panel (1 box per zone)
    out.append(make_box13(f'{seed_prefix}-content-frame', MAIN_X, level_y - 4,
                          MAIN_W, 1080 - (level_y - 4) - 16, z=1))
    out.append(make_text(f'{seed_prefix}-back-btn', MAIN_X + 16, level_y, 160, 32, f'← {bldg["name"]}',
                         color=C_ACCENT, font_size=14, z=20, action=back_to_bldg))
    out.append(make_text(f'{seed_prefix}-level-title', MAIN_X + 196, level_y, 800, 28,
                         f'📋 {floor["name"]}', color=C_TEXT, font_size=22, z=10))
    out.append(make_text(f'{seed_prefix}-level-sub', MAIN_X + 196, level_y + 40, 600, 18,
                         f'{floor["count"]}台设备',
                         color=C_TEXT_DIM, font_size=13, z=10))
    table_x = MAIN_X + 16
    table_y = level_y + 72
    table_w = MAIN_W - 32
    row_h = 42
    # Columns now bind REAL per-device live data (power / current / voltage)
    col_widths = [50, 300, 90, 160, 150, 160, 110, 110]
    col_headers = ['#', '设备名称', '状态', '实时功率(kW)', 'A相电流(A)', 'AB线电压(V)', '协议', '操作']
    # (data point name, unit) bound per row for the data columns (idx 3,4,5)
    bound_cols = {3: '总有功功率', 4: 'A相电流', 5: 'AB线电压'}
    col_starts = [table_x]
    for w in col_widths[:-1]:
        col_starts.append(col_starts[-1] + w)
    out.append(make_panel_bg(f'{seed_prefix}-th-bg', table_x, table_y + 4, table_w, row_h,
                             color=C_PANEL_CARD, z=3, opacity=0.9))
    for hi, (hdr, cs) in enumerate(zip(col_headers, col_starts)):
        if not hdr:
            continue
        out.append(make_text(f'{seed_prefix}-th-{hi}', cs + 8, table_y + 18, col_widths[hi] - 12, 20,
                             hdr, color=C_TEXT_DIM, font_size=11, z=4))
    for di, dev in enumerate(floor['devices'][:15]):
        ry = table_y + 8 + row_h + di * row_h
        dev_page = page_id_device(dev['sid'])
        dev_action = _nav_action(dev_page)
        out.append(make_text(f'{seed_prefix}-row-num-{di}', col_starts[0] + 8, ry + 10, 30, 20,
                             str(di + 1), color=C_TEXT_DIM, font_size=11, z=4))
        out.append(make_text(f'{seed_prefix}-row-name-{di}', col_starts[1] + 8, ry + 10, col_widths[1] - 12, 20,
                             dev['name'], color=C_ACCENT, font_size=13, z=4, action=dev_action))
        status_text = '运行中' if dev['status'] == 1 else '离线'
        status_color = C_GREEN if dev['status'] == 1 else C_TEXT_DIM
        out.append(make_text(f'{seed_prefix}-row-stat-{di}', col_starts[2] + 8, ry + 10, col_widths[2] - 12, 20,
                             status_text, color=status_color, font_size=12, z=4))
        for ci, dpn in bound_cols.items():
            out.append(make_text(f'{seed_prefix}-row-data-{di}-{ci}', col_starts[ci] + 8, ry + 10,
                                 col_widths[ci] - 12, 20, '—', color=C_TEXT, font_size=12, z=4,
                                 data_bound=True, dp_name=dpn, device=dev))
        out.append(make_text(f'{seed_prefix}-row-proto-{di}', col_starts[6] + 8, ry + 10, col_widths[6] - 12, 20,
                             'MODBUS', color=C_TEXT_DIM, font_size=12, z=4))
        out.append(make_text(f'{seed_prefix}-row-btn-{di}', col_starts[7] + 14, ry + 8, 78, 20,
                             '详情 ›', color=C_ACCENT, font_size=12, z=5, action=dev_action))
    table_end_y = table_y + 8 + row_h + min(len(floor['devices']), 15) * row_h + 20
    running = sum(1 for d in floor['devices'] if d['status'] == 1)
    offline = floor['count'] - running
    out.append(make_text(f'{seed_prefix}-summary', table_x, table_end_y, table_w, 22,
                         f'共 {floor["count"]} 台设备 | 运行: {running}台 | 离线: {offline}台',
                         color=C_TEXT_MUTED, font_size=13, z=5))
    return out


def upsert_layer_page(page_name, page_uuid, comp_b64, layer_json=None):
    """Insert or update a display_model_layer row by model_id + page_id."""
    layer_json = layer_json or '{"width":1920,"height":1080,"autoSize":1,"Padding":0,"gridSize":10,"background":"#0a0e17"}'
    cur.execute(
        "SELECT id FROM display_model_layer WHERE model_id=%s AND page_id=%s AND deleted_at IS NULL",
        (MODEL_ID, page_uuid)
    )
    existing = cur.fetchone()
    if existing:
        cur.execute(
            "UPDATE display_model_layer SET page_name=%s, components=%s, layer=%s, updated_at=NOW() WHERE id=%s",
            (page_name, comp_b64, layer_json, existing[0])
        )
        return existing[0], 'updated'
    cur.execute(
        """INSERT INTO display_model_layer
           (model_id, page_name, page_id, is_home, is_login, page_type, layer, components, created_at, updated_at)
           VALUES (%s, %s, %s, 0, 0, 1, %s, %s, NOW(), NOW())""",
        (MODEL_ID, page_name, page_uuid, layer_json, comp_b64)
    )
    return cur.lastrowid, 'inserted'


def build_room_detail_cells(room, seed_prefix='room'):
    """One page per 配电室/楼层 — aggregate cabinet cards (含UPS柜), drill to cabinet.
    Renders cabinet-level rolled-up counts only (scalable), never per-device."""
    out = []
    out.extend(build_header_cells(seed_prefix, [
        ('📊 全局总览', C_TEXT_DIM, PAGE_ID_MAIN),
        (room['name'], C_TEXT, None),
    ]))
    out.extend(build_sidebar_cells(f'{seed_prefix}-nav'))
    level_y = BODY_Y + 16
    out.append(make_box13(f'{seed_prefix}-content-frame', MAIN_X, level_y - 4,
                          MAIN_W, 1080 - (level_y - 4) - 16, z=1))
    out.append(make_text(f'{seed_prefix}-back-btn', MAIN_X + 16, level_y, 140, 32, '← 返回总览',
                         color=C_ACCENT, font_size=14, z=20, action=_nav_action(PAGE_ID_MAIN)))
    out.append(make_text(f'{seed_prefix}-level-title', MAIN_X + 176, level_y, 760, 28,
                         f'🏛 {room["name"]}', color=C_TEXT, font_size=22, z=10))
    out.append(make_text(f'{seed_prefix}-level-sub', MAIN_X + 176, level_y + 40, 760, 18,
                         f'{room["cabinet_count"]}个机柜 · {room["device_count"]}台设备 · 在线{room["online"]} · 异常{room["alarm"]}',
                         color=C_TEXT_DIM, font_size=13, z=10))
    card_w_item = 300
    card_h_item = 150
    gap = 16
    per_row = max(1, (MAIN_W - 32 + gap) // (card_w_item + gap))
    start_x = MAIN_X + 16
    start_y = level_y + 80
    CAB_CARD_CAP = 24
    for ci, cab in enumerate(room['cabinets'][:CAB_CARD_CAP]):
        r, c = ci // per_row, ci % per_row
        cx = start_x + c * (card_w_item + gap)
        cy = start_y + r * (card_h_item + gap)
        if cy + card_h_item > 1080 - 40:
            break
        cab_nav = _nav_action(cab['page_id'])
        ok = cab['alarm'] == 0
        accent = C_GREEN if ok else C_ORANGE
        out.append(make_panel_bg(f'{seed_prefix}-cab-bg-{ci}', cx, cy, card_w_item, card_h_item,
                                 color=C_PANEL_CARD, z=3, opacity=0.9, action=cab_nav))
        out.append(make_panel_bg(f'{seed_prefix}-cab-bar-{ci}', cx, cy, 4, card_h_item, color=accent, z=4))
        out.append(make_text(f'{seed_prefix}-cab-name-{ci}', cx + 16, cy + 14, card_w_item - 32, 24,
                             f'🗄 {cab["name"]}', color=C_ACCENT, font_size=16, z=5, action=cab_nav))
        out.append(make_text(f'{seed_prefix}-cab-sub-{ci}', cx + 16, cy + 48, card_w_item - 32, 18,
                             f'{len(cab["floors"])}个设备组', color=C_TEXT_DIM, font_size=12, z=5, action=cab_nav))
        metrics = [('设备', cab['device_count'], C_TEXT), ('在线', cab['online'], C_GREEN),
                   ('异常', cab['alarm'], C_ORANGE if cab['alarm'] else C_TEXT_DIM)]
        mw = (card_w_item - 32) // 3
        for mi, (mlabel, mval, mcolor) in enumerate(metrics):
            mx = cx + 16 + mi * mw
            out.append(make_text(f'{seed_prefix}-cab-mv-{ci}-{mi}', mx, cy + 80, mw - 8, 28,
                                 str(mval), color=mcolor, font_size=24, z=5, action=cab_nav))
            out.append(make_text(f'{seed_prefix}-cab-ml-{ci}-{mi}', mx, cy + 116, mw - 8, 18,
                                 mlabel, color=C_TEXT_DIM, font_size=11, z=5, action=cab_nav))
    if len(room['cabinets']) > CAB_CARD_CAP:
        out.append(make_text(f'{seed_prefix}-cab-more', start_x, 1080 - 36, MAIN_W - 32, 18,
                             f'+{len(room["cabinets"]) - CAB_CARD_CAP} 更多机柜 · 共 {room["cabinet_count"]} 柜',
                             color=C_TEXT_DIM, font_size=12, z=5))
    return out


print(f"\n=== LEVEL R: ROOM PAGES (per 配电室/楼层) ===")
room_pages = []
for room in rooms:
    seed = f'room-{room["sid"]}'
    room_cells = build_room_detail_cells(room, seed_prefix=seed)
    report_overlaps(room_cells, f'Room {room["name"]}')
    comp_b64 = base64.b64encode(json.dumps({"cells": room_cells}, ensure_ascii=False).encode()).decode()
    row_id, action = upsert_layer_page(f'room-{room["sid"]}', room['page_id'], comp_b64)
    room_pages.append((room['name'], room['page_id'], len(room_cells), action))
    print(f"  {action} room page: {room['name']} id={row_id} cells={len(room_cells)}")
conn.commit()

print(f"\n=== LEVEL 1: BUILDING PAGES (per cabinet) ===")
building_pages = []
for bldg in buildings:
    seed = f'bldg-{bldg["sid"]}'
    bldg_cells = build_building_detail_cells(bldg, seed_prefix=seed)
    report_overlaps(bldg_cells, f'Building {bldg["name"]}')
    comp_json = json.dumps({"cells": bldg_cells}, ensure_ascii=False)
    comp_b64 = base64.b64encode(comp_json.encode()).decode()
    page_name = f'building-{bldg["sid"]}'
    row_id, action = upsert_layer_page(page_name, bldg['page_id'], comp_b64)
    building_pages.append((bldg['name'], bldg['page_id'], len(bldg_cells), action))
    print(f"  {action} building page: {bldg['name']} id={row_id} cells={len(bldg_cells)} page_id={bldg['page_id'][:12]}...")
conn.commit()

print(f"\n=== LEVEL 2: FLOOR PAGES (per device group) ===")
floor_pages = []
for bldg in buildings:
    for floor in bldg['floors']:
        seed = f'floor-{bldg["sid"]}-{floor["key"]}'
        floor_cells = build_floor_detail_cells(bldg, floor, seed_prefix=seed)
        report_overlaps(floor_cells, f'Floor {bldg["name"]}/{floor["name"]}')
        comp_json = json.dumps({"cells": floor_cells}, ensure_ascii=False)
        comp_b64 = base64.b64encode(comp_json.encode()).decode()
        page_name = f'floor-{bldg["sid"]}-{floor["key"]}'
        row_id, action = upsert_layer_page(page_name, floor['page_id'], comp_b64)
        floor_pages.append((f'{bldg["name"]}/{floor["name"]}', floor['page_id'], len(floor_cells), action))
        print(f"  {action} floor page: {bldg['name']}/{floor['name']} id={row_id} cells={len(floor_cells)}")
conn.commit()

# Retire legacy single building-detail / floor-detail pages
cur.execute(
    """UPDATE display_model_layer SET deleted_at=NOW()
       WHERE model_id=%s AND page_name IN ('building-detail', 'floor-detail') AND deleted_at IS NULL""",
    (MODEL_ID,)
)
if cur.rowcount:
    print(f"\nRetired {cur.rowcount} legacy generic detail page(s)")
conn.commit()


# ════════════════════════════════════════════════════════
# LEVEL 3: DEVICE DETAIL PAGE (page_name='device-detail', id=10)
# ════════════════════════════════════════════════════════
print(f"\n=== LEVEL 3: PER-DEVICE DETAIL PAGES ===")


def build_device_detail_cells(dev, bldg, floor, seed_prefix):
    """One page per physical device. All params bind to THIS device's live data
    (deviceSN = dev.uuid) resolved against its own modbus model -> no 张冠李戴."""
    out = []
    dev_name = dev['name']
    out.extend(build_header_cells(seed_prefix, [
        ('📊 全局总览', C_TEXT_DIM, PAGE_ID_MAIN),
        (bldg.get('room_name', ROOT_NAME), C_TEXT_DIM, bldg.get('room_page_id')),
        (bldg['name'], C_TEXT_DIM, bldg['page_id']),
        (floor['name'], C_TEXT_DIM, floor['page_id']),
        (dev_name, C_TEXT, None),
    ]))
    out.extend(build_sidebar_cells(f'{seed_prefix}-nav'))

    level_y = BODY_Y + 16
    out.append(make_text(f'{seed_prefix}-back', MAIN_X, level_y, 200, 32, f'← {floor["name"]}',
                         color=C_ACCENT, font_size=14, z=20, action=_nav_action(floor['page_id'])))
    out.append(make_text(f'{seed_prefix}-title', MAIN_X + 220, level_y, 760, 28,
                         f'🔧 {dev_name}', color=C_TEXT, font_size=22, z=10))
    on = dev['status'] == 1
    out.append(make_text(f'{seed_prefix}-status', MAIN_X + 980, level_y + 4, 200, 22,
                         '● 运行中' if on else '● 离线', color=C_GREEN if on else C_TEXT_DIM,
                         font_size=13, z=10))

    panel_top = level_y + 64
    panel_h = 360
    col_w = (MAIN_W - 32) // 3
    left_x = MAIN_X
    mid_x = MAIN_X + col_w + 16
    right_x = MAIN_X + (col_w + 16) * 2

    # Left: basic info (static, device-specific)
    out.append(make_panel_bg(f'{seed_prefix}-lp', left_x, panel_top, col_w, panel_h, color=C_PANEL, z=2, opacity=0.6))
    out.append(make_box13(f'{seed_prefix}-lf', left_x, panel_top, col_w, panel_h, z=3))
    out.append(make_text(f'{seed_prefix}-lt', left_x + 16, panel_top + 10, 400, 22, '📋 基本参数',
                         color=C_ACCENT, font_size=FONT_PANEL + 2, z=6))
    basic_params = [
        ('设备名称', dev_name), ('设备类型', '多功能电力仪表'),
        ('通信协议', 'Modbus RTU'), ('设备编号', str(dev['sid'])),
        ('所属机房', 'NCC 航信机房'), ('所属区域', f'{bldg["name"]} {floor["name"]}'),
        ('采集周期', '500ms'), ('在线状态', '🟢 运行中' if on else '⚪ 离线'),
    ]
    bp_y = panel_top + 44
    for bi, (bk, bv) in enumerate(basic_params):
        by = bp_y + bi * 36
        out.append(make_text(f'{seed_prefix}-bpk-{bi}', left_x + 16, by, 130, 22,
                             bk, color=C_TEXT_DIM, font_size=FONT_PARAM_LABEL + 1, z=6))
        out.append(make_text(f'{seed_prefix}-bpv-{bi}', left_x + 152, by, col_w - 168, 22,
                             bv, color=C_TEXT, font_size=FONT_PARAM_VAL - 2, z=6))

    # Mid: live electrical params (bound to this device)
    out.append(make_panel_bg(f'{seed_prefix}-mp', mid_x, panel_top, col_w, panel_h, color=C_PANEL, z=2, opacity=0.6))
    out.append(make_box13(f'{seed_prefix}-mf', mid_x, panel_top, col_w, panel_h, z=3))
    out.append(make_text(f'{seed_prefix}-mt', mid_x + 16, panel_top + 10, 400, 22, '📊 实时参数',
                         color=C_ACCENT, font_size=FONT_PANEL + 2, z=6))
    rt_params = [
        ('AB线电压', 'V'), ('BC线电压', 'V'), ('CA线电压', 'V'),
        ('A相电流', 'A'), ('B相电流', 'A'), ('C相电流', 'A'),
        ('中性线电流', 'A'), ('频率', 'Hz'),
    ]
    rtp_y = panel_top + 44
    for ri, (rname, runit) in enumerate(rt_params):
        ry = rtp_y + ri * 36
        out.append(make_text(f'{seed_prefix}-rk-{ri}', mid_x + 16, ry, 150, 22,
                             rname, color=C_TEXT_DIM, font_size=FONT_PARAM_LABEL + 1, z=6))
        out.append(make_text(f'{seed_prefix}-rv-{ri}', mid_x + 172, ry, 110, 22,
                             '—', color=C_ACCENT, font_size=FONT_PARAM_VAL, z=6,
                             data_bound=True, dp_name=rname, device=dev))
        out.append(make_text(f'{seed_prefix}-ru-{ri}', mid_x + 300, ry, 60, 22,
                             runit, color=C_TEXT_DIM, font_size=FONT_PARAM_LABEL, z=6))

    # Right: live power params (bound to this device)
    out.append(make_panel_bg(f'{seed_prefix}-rp', right_x, panel_top, col_w, panel_h, color=C_PANEL, z=2, opacity=0.6))
    out.append(make_box13(f'{seed_prefix}-rf', right_x, panel_top, col_w, panel_h, z=3))
    out.append(make_text(f'{seed_prefix}-rt2', right_x + 16, panel_top + 10, 400, 22, '⚡ 功率参数',
                         color=C_ACCENT, font_size=FONT_PANEL + 2, z=6))
    pw_params = [
        ('总有功功率', 'kW'), ('总无功功率', 'kW'), ('总视在功率', 'kW'),
        ('总功率因数', ''), ('正有功电度', 'kWh'),
    ]
    pw_y = panel_top + 44
    for pi, (pname, punit) in enumerate(pw_params):
        py = pw_y + pi * 44
        out.append(make_text(f'{seed_prefix}-pk-{pi}', right_x + 16, py, 150, 22,
                             pname, color=C_TEXT_DIM, font_size=FONT_PARAM_LABEL + 1, z=6))
        out.append(make_text(f'{seed_prefix}-pv-{pi}', right_x + 172, py, 140, 30,
                             '—', color=C_GREEN, font_size=FONT_KPI_VALUE - 4, z=6,
                             data_bound=True, dp_name=pname, device=dev))
        out.append(make_text(f'{seed_prefix}-pu-{pi}', right_x + 322, py, 60, 22,
                             punit, color=C_TEXT_DIM, font_size=FONT_PARAM_LABEL, z=6))

    chart_y = panel_top + panel_h + 16
    chart_w = int(MAIN_W * 0.65)
    out.append(make_panel_bg(f'{seed_prefix}-cp', MAIN_X, chart_y, chart_w, 280, color=C_PANEL, z=2, opacity=0.6))
    out.append(make_box13(f'{seed_prefix}-cf', MAIN_X, chart_y, chart_w, 280, z=3))
    out.append(make_text(f'{seed_prefix}-ct', MAIN_X + 16, chart_y + 10, 400, 22, '📈 24小时功率曲线',
                         color=C_ACCENT, font_size=FONT_PANEL + 2, z=6))
    out.append(make_smooth_chart(
        f'{seed_prefix}-chart', MAIN_X + 16, chart_y + 40, chart_w - 32, 226,
        title='设备功率趋势', dp_names=['总有功功率', '总无功功率'], z=5, device=dev))

    status_x = MAIN_X + chart_w + 16
    status_w = MAIN_W - chart_w - 16
    out.append(make_panel_bg(f'{seed_prefix}-sp', status_x, chart_y, status_w, 280, color=C_PANEL, z=2, opacity=0.6))
    out.append(make_box13(f'{seed_prefix}-sf', status_x, chart_y, status_w, 280, z=3))
    out.append(make_text(f'{seed_prefix}-stt', status_x + 16, chart_y + 10, 400, 22, '🔔 设备告警',
                         color=C_ORANGE, font_size=FONT_PANEL + 2, z=6))
    out.append(make_text(f'{seed_prefix}-se', status_x + 16, chart_y + 48, status_w - 32, 200,
                         '✅ 该设备无告警记录', color=C_TEXT_DIM, font_size=13, z=6))
    return out


device_pages = 0
for bldg in buildings:
    for floor in bldg['floors']:
        for dev in floor['devices']:
            seed = f'dev-{dev["sid"]}'
            dcells = build_device_detail_cells(dev, bldg, floor, seed_prefix=seed)
            comp_b64 = base64.b64encode(json.dumps({"cells": dcells}, ensure_ascii=False).encode()).decode()
            upsert_layer_page(f'device-{dev["sid"]}', page_id_device(dev['sid']), comp_b64)
            device_pages += 1
conn.commit()
report_overlaps(build_device_detail_cells(buildings[0]['floors'][0]['devices'][0],
                                          buildings[0], buildings[0]['floors'][0], 'dev-sample'),
                'Device detail (sample)')
print(f"  wrote {device_pages} per-device detail pages")

# Keep the legacy shared 'device-detail' page (id=10) pointing at a real device
# so any old reference still resolves to a coherent page.
_sb = buildings[0]; _sf = _sb['floors'][0]; _sd = _sf['devices'][0]
legacy_cells = build_device_detail_cells(_sd, _sb, _sf, 'detail')
comp_b64_detail = base64.b64encode(json.dumps({"cells": legacy_cells}, ensure_ascii=False).encode()).decode()
cur.execute(
    "SELECT id FROM display_model_layer WHERE model_id=%s AND page_name='device-detail' AND deleted_at IS NULL",
    (MODEL_ID,)
)
existing_dev = cur.fetchone()
if existing_dev:
    cur.execute("UPDATE display_model_layer SET page_id=%s, components=%s, updated_at=NOW() WHERE id=%s",
                (PAGE_ID_DEVICE, comp_b64_detail, existing_dev[0]))
else:
    cur.execute(
        """INSERT INTO display_model_layer
           (model_id, page_name, page_id, is_home, is_login, page_type, layer, components, created_at, updated_at)
           VALUES (%s, 'device-detail', %s, 0, 0, 1, '{"height":1080,"width":1920,"autoSize":1}', %s, NOW(), NOW())""",
        (MODEL_ID, PAGE_ID_DEVICE, comp_b64_detail)
    )
conn.commit()
detail_cells = legacy_cells

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
total_bldg_cells = sum(p[2] for p in building_pages)
total_floor_cells = sum(p[2] for p in floor_pages)
print(f"\n{'='*60}")
print(f"✅ Build complete!")
print(f"   Level 0 (overview):        {len(cells)} cells")
print(f"   Level 1 (building pages):  {len(building_pages)} pages, {total_bldg_cells} cells total")
print(f"   Level 2 (floor pages):     {len(floor_pages)} pages, {total_floor_cells} cells total")
print(f"   Level 3 (device-detail):   {len(detail_cells)} cells")
print(f"   Sidebar width: {SIDEBAR_W}px, MAIN_X={MAIN_X}px")
print(f"\n   Page IDs:")
print(f"     MAIN (overview):    {PAGE_ID_MAIN}")
print(f"     DEVICE_DETAIL:      {PAGE_ID_DEVICE}")
if buildings:
    print(f"     Example building:   {buildings[0]['name']} → {buildings[0]['page_id']}")
    if buildings[0]['floors']:
        f0 = buildings[0]['floors'][0]
        print(f"     Example floor:      {f0['name']} → {f0['page_id']}")
print(f"\n   Open in browser: http://localhost:7080/#/AppRun/{MODEL_ID}")
print(f"{'='*60}")

conn.close()
