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
import os
import socket
import sqlite3

_DB_ROOT = os.path.dirname(os.path.abspath(__file__))
_SQLITE_PATH = os.path.join(_DB_ROOT, "ism_server_user", "data", "db", "ism.db")

def _ob_up():
    s = socket.socket()
    try:
        s.settimeout(0.5)
        s.connect(("127.0.0.1", 2881))
        return True
    except OSError:
        return False
    finally:
        s.close()

USE_SQLITE = os.environ.get("ISM_FORCE_OB") != "1" and not _ob_up()
APPRUN_BASE = os.environ.get("NCC_APPRUN_BASE", "http://localhost:7080")


def apprun_url(model_id, page_id=None, base=APPRUN_BASE):
    """AppRun 预览 URL：子页用 ?pageId=，禁止 /AppRun/{model}/{page} 双路径段。"""
    url = f"{base}/#/AppRun/{model_id}"
    if page_id and page_id != model_id:
        url += f"?pageId={page_id}"
    return url

class _SqlAdapter:
    """sqlite3 兼容层：占位符 %s -> ?"""
    def __init__(self, conn):
        self._conn = conn
    def cursor(self):
        return _CurAdapter(self._conn.cursor())
    def close(self):
        self._conn.close()
    def commit(self):
        self._conn.commit()

class _CurAdapter:
    def __init__(self, cur):
        self._cur = cur
    def execute(self, sql, params=None):
        if USE_SQLITE:
            sql = sql.replace("%s", "?").replace("NOW()", "datetime('now')")
        if params is None:
            return self._cur.execute(sql)
        return self._cur.execute(sql, params)
    def fetchone(self):
        return self._cur.fetchone()
    def fetchall(self):
        return self._cur.fetchall()

    @property
    def rowcount(self):
        return self._cur.rowcount

    @property
    def lastrowid(self):
        return self._cur.lastrowid

import re
import uuid as _uuid
from collections import defaultdict

# ── DB ──────────────────────────────────────────────
import time as _time
def _connect_db():
    if USE_SQLITE:
        return _SqlAdapter(sqlite3.connect(_SQLITE_PATH))
    mysql_port = os.environ.get('NCC_MYSQL_PORT')
    if mysql_port:
        return pymysql.connect(
            host=os.environ.get('NCC_MYSQL_HOST', '127.0.0.1'),
            port=int(mysql_port),
            user=os.environ.get('NCC_MYSQL_USER', 'root'),
            password=os.environ.get('NCC_MYSQL_PWD', 'ism2024!'),
            database=os.environ.get('NCC_MYSQL_DB', 'ism'),
        )
    _last = None
    for _ in range(60):
        try:
            return pymysql.connect(
                host='127.0.0.1', port=2881,
                user='root@ism_tenant', password='ism2024!',
                database='ism'
            )
        except Exception as _e:
            _last = _e
            _time.sleep(0.5)
    raise SystemExit(f'DB 连接失败(端口耗尽?): {_last}')
conn = _connect_db()
cur = conn.cursor()

# ── Constants ───────────────────────────────────────
# 可用环境变量覆盖(支持多项目复用): NCC_MODEL_ID / NCC_PROJECT_UUID
# 默认指向当前活的中航信大屏；重建前仍建议三连查后用环境变量覆盖
MODEL_ID = os.environ.get('NCC_MODEL_ID', 'b8b4c094-faa9-a22a-1d0d-037539b27a6c')
PAGE_ID_MAIN = MODEL_ID
PAGE_ID_DEVICE = _uuid.uuid5(_uuid.NAMESPACE_DNS, 'ncc-dash-device-detail').hex
# Per-entity page UUIDs (generated from sid/floor_key below)
def page_id_room(sid):
    return _uuid.uuid5(_uuid.NAMESPACE_DNS, f'ncc-dash-room-{sid}').hex

def page_id_zone(sid):
    return _uuid.uuid5(_uuid.NAMESPACE_DNS, f'ncc-dash-zone-{sid}').hex

def page_id_building(sid):
    return _uuid.uuid5(_uuid.NAMESPACE_DNS, f'ncc-dash-bldg-{sid}').hex

def page_id_floor(bldg_sid, floor_key):
    return _uuid.uuid5(_uuid.NAMESPACE_DNS, f'ncc-dash-floor-{bldg_sid}-{floor_key}').hex

def page_id_device(dev_sid):
    return _uuid.uuid5(_uuid.NAMESPACE_DNS, f'ncc-dash-dev-{dev_sid}').hex

def page_id_oneline(room_sid):
    # 变电所一次系统总图（单线图）。seed 与技能 advanced-electric.md 约定一致。
    return _uuid.uuid5(_uuid.NAMESPACE_DNS, f'ncc-dash-oneline-{room_sid}').hex
PROJECT_UUID = os.environ.get('NCC_PROJECT_UUID', '3ec5821f-b512-2adb-3e1c-473720d0a93e')

# 采样设备(作为无显式设备绑定时的兜底)动态取自本项目第一台 type=1 设备，
# 避免硬编码 1A 设备导致跨项目张冠李戴。
cur.execute("""
    SELECT uuid, name, muid FROM monitor_list
    WHERE project_uuid=%s AND type=1 AND deleted_at IS NULL AND muid IS NOT NULL AND muid<>''
      AND EXISTS (
          SELECT 1 FROM device_real_data d
          WHERE d.device_uuid=monitor_list.uuid
            AND TRIM(COALESCE(d.value, '')) <> ''
      )
    ORDER BY id LIMIT 1
""", (PROJECT_UUID,))
_samp = cur.fetchone()
if _samp:
    DEVICE_UUID, DEVICE_NAME, DEV_MODEL_UUID = _samp[0], _samp[1], _samp[2]
else:
    DEVICE_UUID = '68db26b1-113d-ad7e-79ff-10dbcc1c18d2'
    DEVICE_NAME = '1A1_U11_S18_1'
    DEV_MODEL_UUID = '3d734984-56f6-5494-ad4c-dfc67ca28ac8'

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

# 设备实时表是运行态唯一可信点位来源：部分虚拟/非 Modbus 模型不在
# modbus_devices_data_model 中，仍可用自身的 model_data_uuid 绑定。
DEVICE_DP = defaultdict(dict)
cur.execute("""
    SELECT device_uuid, name, model_data_uuid, data_unit
    FROM device_real_data
    WHERE device_uuid IS NOT NULL AND device_uuid <> ''
      AND model_data_uuid IS NOT NULL AND model_data_uuid <> ''
""")
for _duuid, _name, _data_uuid, _unit in cur.fetchall():
    if _name and _data_uuid:
        DEVICE_DP[_duuid][_name] = {'uuid': _data_uuid, 'unit': _unit or ''}


def dp_map_for(muid, device_uuid=None):
    """优先返回设备实时点位映射，再回退到模型映射。"""
    return DEVICE_DP.get(device_uuid) or MODEL_DP.get(muid) or DP_MAP

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

# 变电所编码 / 展示名（界面统一「X变电所」，库内原名可能含配电室/模块）
_SUB_KEY_RE = re.compile(r'(\d+[AB]\d+)', re.I)

def _substation_key_from_name(name):
    n = (name or '').strip()
    if not n:
        return None
    if n.upper().startswith('ECC'):
        return 'ECC'
    head = re.split(r'[-_及]', n)[0]
    head = head.replace('配电室', '').replace('模块', '').strip()
    m = _SUB_KEY_RE.search(head)
    if m:
        return m.group(1).upper()
    compact = re.sub(r'[^0-9A-Za-z]', '', n)
    m = _SUB_KEY_RE.search(compact)
    return m.group(1).upper() if m else None

def display_substation_name(name):
    """界面展示名统一为「编码变电所」，不暴露配电室/模块等库内原名。"""
    key = _substation_key_from_name(name)
    if key:
        return f'{key}变电所'
    n = (name or '').strip()
    if not n:
        return '变电所'
    if n.upper().startswith('ECC'):
        return 'ECC变电所'
    m = re.match(r'^(\d+[AB])', n, re.I)
    if m:
        return f'{m.group(1).upper()}变电所'
    cleaned = re.sub(r'(配电室|模块).*', '', n).strip()
    return f'{cleaned}变电所' if cleaned else '变电所'

ROOT_SID = ROOT_NODE['sid'] if ROOT_NODE else 1


def _zone_device_stats(zone_sid):
    """统计区域节点下全部 type=1 设备（含多级子区域）。"""
    total = online = 0
    stack = [zone_sid]
    while stack:
        sid = stack.pop()
        for ch in children_by_pid.get(sid, []):
            if ch['type'] == 1:
                total += 1
                if ch['status'] == 1:
                    online += 1
            elif ch['type'] == 0:
                stack.append(ch['sid'])
    return total, online

_room_map = {}
for b in buildings:
    # RootZone 直下区域（UPS/机房模块/配电室）：自身即 room，不能误挂到 RootZone
    rsid = b['sid'] if b.get('pid') == ROOT_SID else b.get('pid')
    rnode = device_by_sid.get(rsid)
    rname = rnode['name'] if rnode else b['name']
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
    r['display_name'] = r['name']

# ── ZONE level (RootZone 直下，与设备管理树一致) ──

def _zone_sid_for(node_sid):
    """沿 pid 向上，找到 RootZone 直下的区域节点 sid。"""
    node = device_by_sid.get(node_sid)
    while node:
        if node['pid'] == ROOT_SID:
            return node['sid']
        parent = device_by_sid.get(node['pid'])
        if parent and parent['pid'] == ROOT_SID:
            return parent['sid']
        node = parent
    return None

_zone_map = {}
for child in children_by_pid.get(ROOT_SID, []):
    if child['type'] != 0:
        continue
    zsid = child['sid']
    dev_total, dev_online = _zone_device_stats(zsid)
    _zone_map[zsid] = {
        'sid': zsid,
        'name': child['name'],
        'page_id': page_id_zone(zsid),
        'rooms': [],
        'device_count': dev_total,
        'online': dev_online,
        'alarm': dev_total - dev_online,
    }
_orphan_rooms = []
for room in rooms:
    zsid = _zone_sid_for(room['sid'])
    if zsid and zsid in _zone_map:
        _zone_map[zsid]['rooms'].append(room)
    else:
        _orphan_rooms.append(room)
if _orphan_rooms:
    dev_total = sum(r['device_count'] for r in _orphan_rooms)
    dev_online = sum(r['online'] for r in _orphan_rooms)
    _zone_map.setdefault(-1, {
        'sid': -1, 'name': '其他', 'page_id': PAGE_ID_MAIN, 'rooms': _orphan_rooms,
        'device_count': dev_total, 'online': dev_online, 'alarm': dev_total - dev_online,
    })
zones = sorted(
    list(_zone_map.values()),
    key=lambda z: (z['sid'] == -1, z['name']),
)
for z in zones:
    z['room_count'] = len(z['rooms'])
    z['cabinet_count'] = sum(r['cabinet_count'] for r in z['rooms']) if z['rooms'] else 0

# ── 一次系统总图节点（按设备管理顶级区域，不再用正则「变电所」分组）──

def _oneline_node_from_zone(z):
    cabs = []
    for r in z['rooms']:
        cabs.extend(r['cabinets'])
    return {
        'key': str(z['sid']),
        'name': z['name'],
        'zone_name': z['name'],
        'sid': z['sid'],
        'page_id': page_id_oneline(z['sid']),
        'cabinet_count': len(cabs),
        'device_count': z['device_count'],
        'online': z['online'],
        'alarm': z['alarm'],
        'cabinets': cabs,
        'rooms': z['rooms'],
    }

substations = [_oneline_node_from_zone(z) for z in zones]

# 一次系统总图入口 = 全园区区域总览；各区域另有独立单线图 page_id_oneline(sid)
PAGE_ID_ONELINE = page_id_oneline('master') if substations else None


def _zone_child_type0_count(zone_sid):
    return sum(1 for c in children_by_pid.get(zone_sid, []) if c['type'] == 0)


def _topo_stat_line(zone):
    """拓扑概览行统计：与设备管理树一致（子区域 + 台数 + 在线/异常）。"""
    parts = []
    sub_cnt = _zone_child_type0_count(zone['sid'])
    if sub_cnt:
        parts.append(f'{sub_cnt}子区域')
    dev = zone['device_count']
    parts.append(f'{dev}台')
    if dev == 0:
        parts.append('暂无设备')
    else:
        if zone['online']:
            parts.append(f'在线{zone["online"]}')
        if zone['alarm']:
            parts.append(f'异常{zone["alarm"]}')
    return ' · '.join(parts)


def _zone_tree_items(zone):
    """区域下钻卡片：与设备管理树一级子节点一致（含直属设备分组）。"""
    zsid = zone['sid']
    items = []
    has_child_areas = _zone_child_type0_count(zsid) > 0
    self_b = next((b for b in buildings if b['sid'] == zsid), None)
    if self_b and self_b['device_count'] > 0:
        label = f'{zone["name"]}直属设备' if has_child_areas else zone['name']
        items.append({
            'name': label,
            'display_name': label,
            'page_id': self_b['page_id'],
            'device_count': self_b['device_count'],
            'online': self_b['online'],
            'alarm': self_b['alarm'],
            'cabinet_count': 1,
        })
    for ch in sorted(children_by_pid.get(zsid, []), key=lambda c: c['name']):
        if ch['type'] != 0:
            continue
        bldg = next((b for b in buildings if b['sid'] == ch['sid']), None)
        if not bldg:
            continue
        items.append({
            'name': ch['name'],
            'display_name': ch['name'],
            'page_id': bldg['page_id'],
            'device_count': bldg['device_count'],
            'online': bldg['online'],
            'alarm': bldg['alarm'],
            'cabinet_count': 1,
        })
    return items

print(f"\n=== Aggregation ({ROOT_NAME}) ===")
for z in zones:
    print(f"  Zone: {z['name']} — {z['room_count']}变电所 / {z['cabinet_count']}柜 / {z['device_count']}台 / 在线{z['online']} / 异常{z['alarm']}")
    for r in z['rooms']:
        print(f"    Sub: {r['display_name']} — {r['cabinet_count']}柜 / {r['device_count']}台")

print(f"\n=== Substations (一次系统区域) ===")
for s in substations:
    print(f"  {s['name']} — {s['cabinet_count']}柜 / {s['device_count']}台 / 在线{s['online']}")

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
    resolved_name = dp_name
    if resolved_name not in dp_map:
        resolved_name = next(
            (name for name in dp_map
             if name.endswith(f'_{dp_name}') or dp_name in name),
            None,
        )
    if not resolved_name:
        return []
    dp = dp_map[resolved_name]
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
            "dataName": resolved_name,
            "operator": "",
            "OperatorValue": "",
            "OperatorMaxValue": ""
        }
    }]

def _resolve_point_name(device, candidates):
    """从设备模型中解析第一个可用的语义测点，兼容带设备名前缀的点位名称。"""
    if not device:
        return None
    point_map = dp_map_for(device.get('muid'), device.get('uuid'))
    for candidate in candidates:
        if candidate in point_map:
            return candidate
    for name in point_map:
        if any(name.endswith(f'_{candidate}') or candidate in name for candidate in candidates):
            return name
    return None

def _make_status_animate(device, status_point=None):
    """设备合分闸状态为真时启用轻量闪烁，离线/断开时保持静止。"""
    point = status_point or _resolve_point_name(
        device, ['输入状态1（合分闸状态）', '合分闸状态', '合分闸']
    )
    if not point:
        return _base_animate()
    active = _base_animate()
    active['selected'] = ['blink']
    active['isExpression'] = True
    active['condition'] = _make_active(
        point, device.get('uuid'), device.get('name'),
        dp_map_for(device.get('muid'), device.get('uuid'))
    )[0]['condition']
    active['condition'].update({'operator': '>', 'OperatorValue': '0'})
    return active

def _make_flow_active(device, flow_point=None):
    """按功率或电流正负驱动潮流方向；没有可用测点时不生成伪动态效果。"""
    point = flow_point or _resolve_point_name(
        device, ['总有功功率', '输出总有功功率', 'A相电流', '输出A相电流']
    )
    if not point:
        return []
    condition = _make_active(
        point, device.get('uuid'), device.get('name'),
        dp_map_for(device.get('muid'), device.get('uuid'))
    )[0]['condition']
    forward = {**condition, 'operator': '>', 'OperatorValue': '0'}
    reverse = {**condition, 'operator': '<', 'OperatorValue': '0'}
    return [
        {'id': 'Forward', 'name': 'component.ViewCanvasMoveLineArrow.Forward',
         'result': '', 'isExpression': True, 'condition': forward},
        {'id': 'Reverse', 'name': 'component.ViewCanvasMoveLineArrow.Reverse',
         'result': '', 'isExpression': True, 'condition': reverse},
    ]

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
SIDEBAR_W = 0                      # 画布内嵌侧栏已移除，导航由前端 ISMRunTreeNav 浮层承担
SIDEBAR_X = 0
MAIN_X = 16                        # 主内容区左边距
MAIN_W = 1920 - MAIN_X - 16        # 1888 — 全宽主面板
BODY_Y = HEADER_H                # 56
# 内容区标题与 box13 顶边留白，避免 22px 大字贴线被裁切
LEVEL_TITLE_TOP_PAD = 12
LEVEL_TITLE_H = 36
LEVEL_SUB_Y = 44
LEVEL_GRID_Y = 72

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

# 配电室详情页 ViewRealTable：列定义与分页
ROOM_TABLE_COLUMNS = ['AB线电压', 'A相电流', 'B相电流', 'C相电流', '总有功功率', '总功率因数', '频率']
ROOM_TABLE_PAGE_SIZE = 15
ROOM_TABLE_ROW_CAP = 200
ROOM_TABLE_REFRESH_MS = 2000


def device_short_label(name):
    """1A1_U11_S18_1 → S18-1，便于表格「设备名称」列阅读。"""
    parts = (name or '').split('_')
    if len(parts) >= 4:
        return f'{parts[2]}-{parts[3]}'
    return name


def collect_room_devices(room):
    """汇总配电室下全部 type=1 设备（跨机柜/设备组）。"""
    devices = []
    seen = set()
    for cab in room.get('cabinets', []):
        for fl in cab.get('floors', []):
            for d in fl.get('devices', []):
                if d['sid'] not in seen:
                    seen.add(d['sid'])
                    devices.append(d)
        for d in cab.get('devices', []):
            if d['sid'] not in seen:
                seen.add(d['sid'])
                devices.append(d)
    return sorted(devices, key=lambda x: x['name'])


def build_room_table_config(room, columns=None):
    """为 ViewRealTable 生成新四字段配置（columnHeaders / rowDevice* / rowBindings）。"""
    columns = columns or ROOM_TABLE_COLUMNS
    devices = collect_room_devices(room)[:ROOM_TABLE_ROW_CAP]
    if not devices:
        return None
    muid = devices[0].get('muid') or DEV_MODEL_UUID
    dp_map = dp_map_for(muid)
    cols = [c for c in columns if c in dp_map]
    if not cols:
        cols = list(dp_map.keys())[:8]
    row_names = [device_short_label(d['name']) for d in devices]
    row_codes = [d['name'] for d in devices]
    binding_rows = [[f"{d['name']}->{col}" for col in cols] for d in devices]
    return {
        'column_headers': cols,
        'row_names': row_names,
        'row_codes': row_codes,
        'bindings': binding_rows,
        'device_count': len(devices),
    }


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
              z=10, data_bound=False, dp_name=None, action=None, device=None, align=None):
    cell_id = gen_uid(seed)
    if data_bound and dp_name:
        dev = device or {}
        active = _make_active(dp_name, dev.get('uuid'), dev.get('name'),
                              dp_map_for(dev.get('muid'), dev.get('uuid')) if dev else None)
    else:
        active = []
    _extra_style = {}
    if align:
        _extra_style['textAlign'] = align
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
                                     foreColor=color, borderWidth=0, BorderEdges=0,
                                     **_extra_style),
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

def _build_chart_active(dp_names, device=None):
    """构建图表组件的实时绑点 active 列表（smooth / history-trend 共用，DRY）。
    每个数据点对应一个 ShowChartVariableN 变量，dataID 取模型点 uuid（= model_data_uuid，
    实时与历史查询均按它匹配）。"""
    active = []
    dev = device or {}
    d_uuid = dev.get('uuid') or DEVICE_UUID
    d_name = dev.get('name') or DEVICE_NAME
    d_map = dp_map_for(dev.get('muid'), dev.get('uuid')) if dev else DP_MAP
    var_ids = ['ShowChartVariable1', 'ShowChartVariable2', 'ShowChartVariable3',
               'ShowChartVariable4', 'ShowChartVariable5']
    for i, dpn in enumerate(dp_names[:5]):
        binding = _make_active(dpn, d_uuid, d_name, d_map)
        if not binding:
            continue
        active.append({
            **binding[0],
            "id": var_ids[i],
        })
    return active


def _chart_animate_block():
    """图表组件统一的 animate 配置块（每次返回独立 dict，避免共享引用）。"""
    return {
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
    }


def make_smooth_chart(seed, x, y, w, h, title, dp_names, z=5, device=None, show_title=True):
    cell_id = gen_uid(seed)
    chart_title = title if show_title else ''
    active = _build_chart_active(dp_names, device)
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
                        {"name": "configComponent.ChartPublic.ChartTitle", "type": 4, "value": chart_title, "key": "ChartTitle"},
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
                "animate": _chart_animate_block(),
                "action": [], "active": active, "dataBind": []
            }
        }
    }


def _with_energy_overview_role(cell, role):
    """只为首页能源四组件写入后端统计角色，并移除实时绑点。"""
    detail = cell["data"]["detail"]
    detail["energyOverviewRole"] = role
    detail["active"] = []
    return cell


def make_alarm_history(seed, x, y, w, h, z=6):
    """组态告警历史查询组件（AlarmHistoryComponents），深色主题适配大屏。

    首页右侧栏已恢复双趋势图 + 右下角活跃告警；历史查询由 ScadaAlarmPanel 抽屉打开。
    本函数仍保留供组态编辑器/工具箱手工拖入使用。
    """
    cell_id = gen_uid(seed)
    diy = [
        {"name": "configComponent.DeviceTree.ShowCount", "type": 1,
         "value": 8, "min": 1, "max": 100, "key": "ShowCount"},
        {"name": "configComponent.DeviceTree.SearchColor", "type": 2,
         "value": "#9fb6d6", "key": "SearchColor"},
        {"name": "configComponent.DeviceTree.SearchBackColor", "type": 2,
         "value": "#0b1c2b", "key": "SearchBackColor"},
        {"name": "configComponent.DeviceTree.SearchBorderColor", "type": 2,
         "value": "#1e3a5f", "key": "SearchBorderColor"},
        {"name": "configComponent.DataHistoryList.dateSelectColor", "type": 2,
         "value": "#00e5ff", "key": "dateSelectColor"},
        {"name": "configComponent.DataHistoryList.dateSelectBackColor", "type": 2,
         "value": "#0b1c2b", "key": "dateSelectBackColor"},
        {"name": "configComponent.DataHistoryList.dateSelectBorderColor", "type": 2,
         "value": "#1e3a5f", "key": "dateSelectBorderColor"},
        {"name": "configComponent.DataHistoryList.tableHeaderColor", "type": 2,
         "value": "#9fefff", "key": "tableHeaderColor"},
        {"name": "configComponent.DataHistoryList.tableHeaderBackColor", "type": 2,
         "value": "#0d2438", "key": "tableHeaderBackColor"},
        {"name": "configComponent.DataHistoryList.tableSplitColor", "type": 2,
         "value": "#1e3a5f", "key": "tableSplitColor"},
        {"name": "configComponent.DataHistoryList.tableHoverColor", "type": 2,
         "value": "#12304a", "key": "tableHoverColor"},
    ]
    return {
        "shape": "AlarmHistoryComponents",
        "id": cell_id, "x": x, "y": y, "width": w, "height": h,
        "zIndex": z, "visible": True,
        "position": {"x": x, "y": y},
        "size": {"width": w, "height": h},
        "data": {
            "detail": {
                "type": "AlarmHistoryComponents",
                "identifier": cell_id,
                "name": "告警历史查询",
                "style": {
                    "position": {"x": x, "y": y, "w": w, "h": h},
                    "visible": 1,
                    "backColor": "rgba(11, 28, 43, 0.92)",
                    "foreColor": "#e8f1ff",
                    "zIndex": z,
                    "transform": 0,
                    "diy": diy,
                },
                "animate": _base_animate(),
                "action": [],
                "active": [],
                "dataBind": [],
            }
        },
    }


def make_real_table(seed, x, y, w, h, column_headers, row_names, row_codes, binding_rows,
                    z=5, theme='dark', page_size=ROOM_TABLE_PAGE_SIZE,
                    refresh_ms=ROOM_TABLE_REFRESH_MS):
    """ViewRealTable 实时数据表格（新四字段：columnHeaders / rowDeviceNames / rowDeviceCodes / rowBindings）。"""
    cell_id = gen_uid(seed)
    col_text = ', '.join(column_headers)
    names_text = '\n'.join(row_names)
    codes_text = '\n'.join(row_codes)
    bindings_text = ';'.join(','.join(row) for row in binding_rows)
    diy = [
        {"name": "configComponent.viewRealTable.columnHeaders", "type": 9,
         "value": col_text, "key": "columnHeaders"},
        {"name": "configComponent.viewRealTable.rowDeviceNames", "type": 9,
         "value": names_text, "key": "rowDeviceNames"},
        {"name": "configComponent.viewRealTable.rowDeviceCodes", "type": 9,
         "value": codes_text, "key": "rowDeviceCodes"},
        {"name": "configComponent.viewRealTable.rowBindings", "type": 9,
         "value": bindings_text, "key": "rowBindings"},
        {"name": "configComponent.AlarmList.waitTime", "type": 7,
         "value": refresh_ms, "min": 100, "max": 10000, "key": "waitTime"},
        {"name": "configComponent.DeviceTree.ShowCount", "type": 1,
         "value": page_size, "min": 1, "max": 100, "key": "ShowCount"},
        {"name": "configComponent.DeviceTree.SearchColor", "type": 2,
         "value": "#e5eefc", "key": "SearchColor"},
        {"name": "configComponent.DeviceTree.SearchBackColor", "type": 2,
         "value": "#162033", "key": "SearchBackColor"},
        {"name": "configComponent.DeviceTree.SearchBorderColor", "type": 2,
         "value": "#314158", "key": "SearchBorderColor"},
        {"name": "configComponent.DataHistoryList.tableHeaderColor", "type": 2,
         "value": "#f8fbff", "key": "tableHeaderColor"},
        {"name": "configComponent.DataHistoryList.tableHeaderBackColor", "type": 2,
         "value": "#1d3557", "key": "tableHeaderBackColor"},
        {"name": "configComponent.viewRealTable.tableHeaderFont", "type": 3,
         "value": "Arial", "key": "tableHeaderFont"},
        {"name": "configComponent.viewRealTable.tableHeaderFontSize", "type": 1,
         "value": 14, "key": "tableHeaderFontSize"},
        {"name": "configComponent.DataHistoryList.tableSplitColor", "type": 2,
         "value": "#263449", "key": "tableSplitColor"},
        {"name": "configComponent.DataHistoryList.tableHoverColor", "type": 2,
         "value": "#1e3a5f", "key": "tableHoverColor"},
        {"name": "主题风格", "type": 6, "value": theme,
         "enumList": [
             {"value": "light", "option": "极简亮色"},
             {"value": "dark", "option": "深空夜幕"},
             {"value": "ocean", "option": "海岸蓝调"},
             {"value": "amber", "option": "琥珀暖光"},
             {"value": "emerald", "option": "森林翠影"},
         ], "key": "themeName"},
    ]
    return {
        "shape": "ism-view-real-table",
        "id": cell_id, "x": x, "y": y, "width": w, "height": h,
        "zIndex": z, "visible": True,
        "position": {"x": x, "y": y},
        "size": {"width": w, "height": h},
        "data": {
            "detail": {
                "type": "ism-view-real-table",
                "identifier": cell_id,
                "name": "实时数据预览",
                "style": _make_style({"x": x, "y": y, "w": w, "h": h},
                                     text="", fontSize=13, foreColor=C_TEXT,
                                     backColor="transparent", borderWidth=0, BorderEdges=0,
                                     diy=diy),
                "animate": _base_animate(),
                "action": [], "active": [], "dataBind": []
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
            cells_out.append(make_text(sep_id, cx, y, 18, 22, '➜',
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
    out.append(make_text(f'{seed_prefix}-header-logo', 16, 10, 36, 36, '⚡',
                         color=C_ACCENT, font_size=26, z=10))
    out.append(make_text(f'{seed_prefix}-header-title', 56, 3, 520, 24, '中航信数据中心电力监控系统',
                         color=C_TEXT, font_size=FONT_TITLE, z=10, align='left'))
    out.append(make_text(f'{seed_prefix}-header-subtitle', 58, 35, 520, 14, 'AVIC INFO DATA CENTER POWER MONITORING SYSTEM',
                         color=C_TEXT_DIM, font_size=FONT_SUBTITLE, z=10, align='left'))
    # 标题左对齐后，面包屑从中后段起排，中间留空隙
    out.extend(make_breadcrumb(f'{seed_prefix}-header-crumb', 1020, 20, breadcrumb_segments, z=20))
    # 右上角依次紧凑排：日期/时间 → 🟢在线 → ⚙齿轮槽位（最右角，由前端 BackToAdminButton
    # 浮层 right:13px 落入）。这里把时钟右移到 x≈1650、在线右移到 x≈1842（文字左对齐，右端落
    # 在 canvas x≈1881），在 1881~1920 留出空当给齿轮；按 1920×1080 满屏标定，与齿轮约 10px 间距、
    # 与日期约 12px 间距，x+w≤1920 不裁、report_overlaps 无重叠（见 references/build-internals.md）。
    out.append(make_svg_time(f'{seed_prefix}-header-clock', 1650, 16, 190, 24,
                             z=10, color=C_ACCENT, font_size=FONT_SUBTITLE,
                             time_format='YYYY/MM/DD HH:mm:ss', show_week=1))
    out.append(make_text(f'{seed_prefix}-header-status', 1842, 18, 54, 22, '🟢 在线',
                         color=C_GREEN, font_size=FONT_SUBTITLE, z=10))
    # STATIC HUD divider line (no marquee animation): a dim full-width rule with a
    # brighter cyan accent segment under the title -> tech feel, zero flicker.
    out.append(make_panel_bg(f'{seed_prefix}-header-rule', 0, HEADER_H - 1, 1920, 1,
                             color=C_BORDER, z=3, opacity=0.9))
    out.append(make_panel_bg(f'{seed_prefix}-header-accent', MAIN_X, HEADER_H - 2, 220, 2,
                             color=C_ACCENT, z=4, opacity=0.9))
    return out


def build_sidebar_cells(seed_prefix='nav'):
    """已弃用：左侧 dv-border-box-13 内嵌导航由前端 ISMRunTreeNav 替代。"""
    return []


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
# LEVEL 0: OVERVIEW PAGE — 在 append_campus_oneline_diagram 定义后构建（见文件后部）
# ════════════════════════════════════════════════════════


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
    back_to_room = _nav_action(bldg.get('room_page_id') or PAGE_ID_MAIN)
    level_y = BODY_Y + 16
    title_y = level_y + LEVEL_TITLE_TOP_PAD
    # Push title block down from the frame's top border so the title doesn't touch the line
    # Single neon frame around the cabinet content panel (1 box per zone)
    out.append(make_box13(f'{seed_prefix}-content-frame', MAIN_X, level_y - 4,
                          MAIN_W, 1080 - (level_y - 4) - 16, z=1))
    out.append(make_text(f'{seed_prefix}-back-btn', MAIN_X + 16, title_y, 180, 32,
                         f'← {bldg.get("room_name", "返回")}',
                         color=C_ACCENT, font_size=14, z=20, action=back_to_room))
    out.append(make_text(f'{seed_prefix}-level-title', MAIN_X + 214, title_y, 760, LEVEL_TITLE_H,
                         f'🏢 {bldg["name"]}', color=C_TEXT, font_size=22, z=10))
    alarm_cnt = sum(1 for f in bldg['floors'] for d in f['devices'] if d['status'] != 1)
    out.append(make_text(f'{seed_prefix}-level-sub', MAIN_X + 214, title_y + LEVEL_SUB_Y, 760, 18,
                         f'{bldg["device_count"]}台设备 · {alarm_cnt}条异常',
                         color=C_TEXT_DIM, font_size=13, z=10))
    card_start_x = MAIN_X + 16
    card_start_y = title_y + LEVEL_GRID_Y
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
    back_to_bldg = _nav_action(bldg['page_id'])
    level_y = BODY_Y + 16
    title_y = level_y + LEVEL_TITLE_TOP_PAD
    # Single neon frame around the device-list panel (1 box per zone)
    out.append(make_box13(f'{seed_prefix}-content-frame', MAIN_X, level_y - 4,
                          MAIN_W, 1080 - (level_y - 4) - 16, z=1))
    out.append(make_text(f'{seed_prefix}-back-btn', MAIN_X + 16, title_y, 160, 32, f'← {bldg["name"]}',
                         color=C_ACCENT, font_size=14, z=20, action=back_to_bldg))
    out.append(make_text(f'{seed_prefix}-level-title', MAIN_X + 196, title_y, 800, LEVEL_TITLE_H,
                         f'📋 {floor["name"]}', color=C_TEXT, font_size=22, z=10))
    out.append(make_text(f'{seed_prefix}-level-sub', MAIN_X + 196, title_y + LEVEL_SUB_Y, 600, 18,
                         f'{floor["count"]}台设备',
                         color=C_TEXT_DIM, font_size=13, z=10))
    table_x = MAIN_X + 16
    table_y = title_y + LEVEL_GRID_Y
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
            real_dpn = floor_col_dp(dev.get('muid'), dpn)  # UPS 用输出量等价点
            out.append(make_text(f'{seed_prefix}-row-data-{di}-{ci}', col_starts[ci] + 8, ry + 10,
                                 col_widths[ci] - 12, 20, '—', color=C_TEXT, font_size=12, z=4,
                                 data_bound=True, dp_name=real_dpn, device=dev))
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


# ──────────────────────────────────────────────────────
# ELECTRIC ONE-LINE (一次系统总图) BUILDERS
# 复用前端既有 SVG 电气组件：view-svg-electric1~8 / view-svg-line / ViewCanvasMoveLineArrow。
# 铁律同文字单元：animate.selected=[]（无则 includes 崩成 #comment）、style.visible=1、diy=[]。
# ──────────────────────────────────────────────────────

def make_electric(shape, seed, x, y, w, h, color=C_ACCENT, fill=C_PANEL,
                  action=None, z=8, stroke_width=1.5, device=None, status_point=None):
    """电气符号；给定设备时由合分闸状态控制闪烁，仍可挂下钻 action。"""
    cell_id = gen_uid(seed)
    status_name = status_point or _resolve_point_name(
        device, ['输入状态1（合分闸状态）', '合分闸状态', '合分闸']
    )
    active = _make_active(
        status_name, device.get('uuid'), device.get('name'),
        dp_map_for(device.get('muid'), device.get('uuid'))
    ) if device and status_name else []
    style = _make_style(
        {"x": x, "y": y, "w": w, "h": h},
        text="", backColor="transparent", foreColor=color, fontSize=1, borderWidth=0,
        BorderEdges=0,
        diy=[
            {"name": "component.Electric.ElectronicDeviceWidth", "type": 7,
             "value": stroke_width, "min": 0.1, "key": "strokeWidth"},
            {"name": "component.Electric.ElectronicDeviceColor", "type": 2,
             "value": color, "key": "strokeColor"},
            {"name": "component.Electric.ConnectDiameter", "type": 1,
             "value": "1", "min": 1, "key": "ConnectDiameter"},
            {"name": "component.Electric.ConnectColor", "type": 2,
             "value": fill, "key": "strokeFill"},
            {"name": "component.public.fillOpacity", "type": 7,
             "value": 1, "min": 0.1, "max": 1, "key": "fillOpacity"},
        ])
    return {
        "shape": shape,
        "id": cell_id, "x": x, "y": y, "width": w, "height": h,
        "zIndex": z, "visible": True,
        "position": {"x": x, "y": y},
        "size": {"width": w, "height": h},
        "data": {"detail": {
            "type": shape, "identifier": cell_id, "name": seed,
            "style": style, "animate": _make_status_animate(device, status_name) if device else _base_animate(),
            "action": action or [], "active": active, "dataBind": []
        }}
    }


def make_conn_line(seed, x, y, w, h, color=C_ACCENT, z=4, opacity=0.85):
    """静态连线（母线/馈线直线段）。用透明文字单元当细矩形画线，DRY 且天然防崩。"""
    return make_panel_bg(seed, x, y, w, h, color=color, z=z, opacity=opacity)


def make_move_line(seed, x, y, w, h, color=C_ACCENT, z=5, vertical=False,
                   stroke_width=3, direction=0, back_color=C_BORDER,
                   device=None, flow_point=None):
    """ViewCanvasMoveLineArrow：母线/馈线「缓慢流动」潮流箭头（科技感来源，克制使用）。
    必须给 style.points（否则前端 .length 崩）；foreColor=流动色；strokeWidth 细=克制。"""
    cell_id = gen_uid(seed)
    if vertical:
        pts = [{"x": w / 2, "y": 0, "isArrow": False},
               {"x": w / 2, "y": h / 2, "isArrow": True},
               {"x": w / 2, "y": h, "isArrow": False}]
    else:
        pts = [{"x": 0, "y": h / 2, "isArrow": False},
               {"x": w / 2, "y": h / 2, "isArrow": True},
               {"x": w, "y": h / 2, "isArrow": False}]
    _dir_enum = [{"value": 0, "option": "configComponent.bigScreen.border.border89DirectionForward"},
                 {"value": 1, "option": "configComponent.bigScreen.border.border89DirectionNegative"}]
    active = _make_flow_active(device, flow_point) if device else []
    style = {
        "position": {"x": x, "y": y, "w": w, "h": h},
        "points": pts, "visible": 1, "opacity": 1, "transform": 0,
        "backColor": "transparent", "foreColor": color, "zIndex": z,
        "diy": [
            {"name": "displayConfig.ToolBox.Diagram.MoveBrokenLineBackColor", "type": 2,
             "value": back_color, "key": "MoveBrokenLineBackColor"},
            {"name": "component.public.strokeWidth", "type": 1,
             "value": stroke_width, "min": 1, "key": "strokeWidth"},
            {"name": "displayConfig.ToolBox.Diagram.MoveBrokenLineConditionEnable", "type": 6,
             "value": 1 if active else 0, "min": 1, "key": "MoveBrokenLineConditionEnable",
             "enumList": [{"value": 0, "option": "component.public.Forbidden"},
                          {"value": 1, "option": "component.public.Enable"}]},
            {"name": "configComponent.bigScreen.border.border89Direction", "type": 6,
             "value": direction, "min": 1, "key": "spinDirection", "enumList": _dir_enum},
        ]
    }
    return {
        "shape": "ViewCanvasMoveLineArrow",
        "id": cell_id, "x": x, "y": y, "width": w, "height": h,
        "zIndex": z, "visible": True,
        "position": {"x": x, "y": y},
        "size": {"width": w, "height": h},
        "data": {"detail": {
            "type": "ViewCanvasMoveLineArrow", "identifier": cell_id, "name": seed,
            "style": style, "animate": _base_animate(),
            "action": [], "active": active, "dataBind": []
        }}
    }


def _oneline_points_for(muid, device_uuid=None):
    """按设备模型解析该路要显示的电气量（电压/电流/有功/状态），避免张冠李戴。
    返回 [(label, dp_name, unit), ...]，dp_name 必须存在于该模型点表。"""
    m = dp_map_for(muid, device_uuid)
    if '总有功功率' in m and 'AB线电压' in m:           # 标准多功能电力仪表
        return [('母线电压', 'AB线电压', 'V'), ('A相电流', 'A相电流', 'A'),
                ('有功功率', '总有功功率', 'kW'), ('合分闸', '输入状态1（合分闸状态）', '')]
    if '输出AB线电压' in m:                              # 施耐德/伊顿 UPS 模型
        return [('输出电压', '输出AB线电压', 'V'), ('输出电流', '输出A相电流', 'A'),
                ('输出有功', '输出总有功功率', 'kW'), ('运行模式', 'UPS使用模式', '')]
    # 兜底：取该模型前若干个点
    keys = list(m.keys())
    return [(k, k, m[k].get('unit', '')) for k in keys[:4]]


def _is_ups_model(muid):
    """该设备模型是否为 UPS（点名体系与电力仪表完全不同）。"""
    m = dp_map_for(muid)
    return ('输出AB线电压' in m) or ('UPS使用模式' in m)


def detail_params_for(muid):
    """按设备模型返回设备详情页要绑的 (实时参数, 功率参数, 趋势图点, 设备类型名)。
    UPS 模型点名与电表完全不同 → 必须按模型选点，否则 _make_active 全返回 []（张冠李戴/缺绑）。"""
    if _is_ups_model(muid):
        rt = [('输出AB线电压', 'V'), ('输出BC线电压', 'V'), ('输出CA线电压', 'V'),
              ('输出A相电流', 'A'), ('输出B相电流', 'A'), ('输出C相电流', 'A'),
              ('主路输入AB线电压', 'V'), ('输出频率', 'Hz')]
        pw = [('输出总有功功率', 'kW'), ('输出视在功率', 'kW'), ('输出功率因数', ''),
              ('电池电压', 'V'), ('电池剩余运行时间', 'min')]
        chart = ['输出总有功功率', '输出视在功率']
        return rt, pw, chart, 'UPS（施耐德）'
    rt = [('AB线电压', 'V'), ('BC线电压', 'V'), ('CA线电压', 'V'),
          ('A相电流', 'A'), ('B相电流', 'A'), ('C相电流', 'A'),
          ('中性线电流', 'A'), ('频率', 'Hz')]
    pw = [('总有功功率', 'kW'), ('总无功功率', 'kW'), ('总视在功率', 'kW'),
          ('总功率因数', ''), ('正有功电度', 'kWh')]
    chart = ['总有功功率', '总无功功率']
    return rt, pw, chart, '多功能电力仪表'


# 设备组列表表头是固定语义列；UPS 用「输出量」作为等价绑点，避免该列取不到值。
_FLOOR_COL_UPS_ALIAS = {'总有功功率': '输出总有功功率', 'A相电流': '输出A相电流', 'AB线电压': '输出AB线电压'}


def floor_col_dp(muid, canonical):
    """把设备组表格的语义列名解析为该设备模型实际存在的点名（UPS→输出量）。"""
    m = dp_map_for(muid)
    if canonical in m:
        return canonical
    alt = _FLOOR_COL_UPS_ALIAS.get(canonical)
    if alt and alt in m:
        return alt
    return canonical


# 三路馈线使用的电气符号（distinct glyph，仅作视觉区分，非真实型号语义）
_ONELINE_BRANCH_SHAPES = ['view-svg-electric2', 'view-svg-electric3',
                          'view-svg-electric7', 'view-svg-electric4',
                          'view-svg-electric5', 'view-svg-electric6']


def _campus_card_title(sub):
    """与设备管理树区域名一致。"""
    return sub.get('name') or sub.get('zone_name') or '区域'


def _append_campus_substation_card(out, seed_prefix, si, cx, cy, card_w, card_h, sub, accent):
    """组织节点卡片：展示区域与设备统计；上层组织节点不提供旧页面下钻。"""
    sub_action = None
    source_device = next(
        (cab['devices'][0] for cab in sub.get('cabinets', []) if cab.get('devices')),
        None,
    )
    out.append(make_panel_bg(f'{seed_prefix}-card-{si}', cx, cy, card_w, card_h,
                             color=C_PANEL_CARD, z=3, opacity=0.88, action=sub_action))
    out.append(make_panel_bg(f'{seed_prefix}-bar-{si}', cx, cy, 3, card_h, color=accent, z=4))

    bus_y = cy + 10
    bus_l, bus_r = cx + 14, cx + card_w - 14
    out.append(make_conn_line(f'{seed_prefix}-bus-{si}', bus_l, bus_y, bus_r - bus_l, 2,
                              color=C_ACCENT, z=5, opacity=0.9))
    out.append(make_text(f'{seed_prefix}-bus-lab-{si}', cx + 8, bus_y - 14, card_w - 16, 14,
                         '组织区域', color=C_TEXT_DIM, font_size=10, z=10, action=sub_action, align='center'))

    sym_s = 34
    sym_y = bus_y + 16
    left_x = cx + card_w // 2 - sym_s - 6
    right_x = cx + card_w // 2 + 6
    mid_x = cx + card_w // 2
    out.append(make_conn_line(f'{seed_prefix}-drop-l-{si}', left_x + sym_s // 2 - 1, bus_y + 2, 2,
                              sym_y - bus_y - 2, color=C_ACCENT, z=4))
    out.append(make_conn_line(f'{seed_prefix}-drop-r-{si}', right_x + sym_s // 2 - 1, bus_y + 2, 2,
                              sym_y - bus_y - 2, color=C_ACCENT, z=4))
    out.append(make_electric('view-svg-electric1', f'{seed_prefix}-sym-l-{si}',
                             left_x, sym_y, sym_s, sym_s, color=C_ACCENT, action=sub_action, z=8,
                             device=source_device))
    out.append(make_electric('view-svg-electric2', f'{seed_prefix}-sym-r-{si}',
                             right_x, sym_y, sym_s, sym_s, color=C_ACCENT, action=sub_action, z=8,
                             device=source_device))

    name_y = sym_y + sym_s + 10
    name_h = 30
    out.append(make_panel_bg(f'{seed_prefix}-name-bg-{si}', cx + 12, name_y, card_w - 24, name_h,
                             color=C_PANEL, z=5, opacity=0.95, action=sub_action))
    out.append(make_text(f'{seed_prefix}-name-{si}', cx + 12, name_y + 4, card_w - 24, name_h - 8,
                         _campus_card_title(sub), color=C_ACCENT, font_size=13, z=10, action=sub_action, align='center'))

    drop_top = name_y + name_h + 4
    drop_bot = cy + card_h - 28
    if drop_bot > drop_top + 8:
        out.append(make_conn_line(f'{seed_prefix}-out-{si}', mid_x - 1, drop_top, 2,
                                  drop_bot - drop_top, color=C_ACCENT, z=4))
        out.append(make_move_line(f'{seed_prefix}-flow-{si}', mid_x - 8, drop_top, 16,
                                  drop_bot - drop_top, color=C_ACCENT, z=5, vertical=True,
                                  stroke_width=2, direction=0, device=source_device))
    mod_label = f'{sub["device_count"]} 台设备'
    out.append(make_text(f'{seed_prefix}-mod-{si}', cx + 8, cy + card_h - 24, card_w - 16, 18,
                         mod_label, color=C_TEXT, font_size=11, z=10, action=sub_action, align='center'))
    out.append(make_text(f'{seed_prefix}-stat-{si}', cx + card_w - 72, cy + 6, 64, 14,
                         f'在线 {sub["online"]}',
                         color=C_TEXT_MUTED, font_size=9, z=10, action=sub_action))


def _overview_alarm_items(limit=7):
    """汇总离线设备作为活跃告警列表数据源。"""
    items = []
    for b in buildings:
        loc = b.get('room_name') or b.get('name', '')
        for f in b['floors']:
            for d in f['devices']:
                if d.get('status') != 1:
                    items.append({
                        'name': d['name'],
                        'location': loc,
                        'cabinet': b['name'],
                        'device': d,
                        'page_id': page_id_device(d['sid']),
                    })
    return items[:limit]


OVERVIEW_TOPO_W = 400
OVERVIEW_TOPO_H = 220
TOPO_NAME_W = 148


def append_overview_topo_panel(out, seed_prefix, x, y, w, h):
    """左下角拓扑概览浮层：按 RootZone 顶级区域摘要，与设备管理树一致。"""
    out.append(make_panel_bg(f'{seed_prefix}-topo-fill', x + 4, y + 4, w - 8, h - 8,
                             color=C_PANEL, z=15, opacity=0.92))
    out.extend(make_panel_title(f'{seed_prefix}-topo', x + 12, y + 10,
                                f'拓扑概览 · {ROOT_NAME}', color=C_ACCENT, font_size=12, z=16, w=260))
    out.append(make_text(f'{seed_prefix}-topo-hint', x + 12, y + 32, w - 24, 14,
                         '与设备管理树一致 · 点击区域下钻', color=C_TEXT_DIM, font_size=10, z=16))

    line_y = y + 50
    line_h = 24
    line_limit = y + h - 20
    TOPO_ZONE_CAP = 6
    stat_x = x + 12 + TOPO_NAME_W + 8
    stat_w = w - (stat_x - x) - 12
    shown = 0
    for zi, zone in enumerate(zones[:TOPO_ZONE_CAP]):
        if line_y + line_h > line_limit:
            break
        zone_action = _nav_action(zone['page_id'])
        accent = C_GREEN if zone['alarm'] == 0 and zone['device_count'] > 0 else (
            C_TEXT_DIM if zone['device_count'] == 0 else C_ORANGE)
        out.append(make_panel_bg(f'{seed_prefix}-topo-row-{zi}', x + 10, line_y, w - 20, line_h,
                                 color=C_PANEL_CARD, z=16, opacity=0.85, action=zone_action))
        out.append(make_panel_bg(f'{seed_prefix}-topo-bar-{zi}', x + 10, line_y, 3, line_h,
                                 color=accent, z=17))
        out.append(make_text(f'{seed_prefix}-topo-zone-{zi}', x + 16, line_y + 4, TOPO_NAME_W, 16,
                             zone['name'], color=C_ACCENT, font_size=11, z=18, action=zone_action))
        out.append(make_text(f'{seed_prefix}-topo-stat-{zi}', stat_x, line_y + 5, stat_w, 14,
                             _topo_stat_line(zone),
                             color=C_TEXT_MUTED, font_size=10, z=18, action=zone_action))
        line_y += line_h + 4
        shown += 1
    if len(zones) > shown:
        out.append(make_text(f'{seed_prefix}-topo-more', x + 12, line_limit - 2, w - 24, 14,
                             f'+{len(zones) - shown} 更多区域 ›', color=C_TEXT_DIM, font_size=10, z=18))


def append_overview_alarm_panel(out, seed_prefix, x, y, w, h):
    """右侧告警列表面板：展示离线/异常设备，可点进详情。"""
    alarms = _overview_alarm_items(limit=7)
    out.extend(make_panel_title(f'{seed_prefix}-alarm', x + 12, y + 10,
                                '活跃告警', color=C_ORANGE, font_size=FONT_PANEL, z=7, w=160))
    alarm_n = sum(1 for row in all_devices if row[4] == 1 and row[6] != 1)
    out.append(make_text(f'{seed_prefix}-alarm-badge', x + w - 88, y + 12, 76, 16,
                         f'● {alarm_n} 条', color=C_ORANGE, font_size=11, z=8))

    row_y = y + 38
    row_h = 24
    row_limit = y + h - 12
    if not alarms:
        out.append(make_text(f'{seed_prefix}-alarm-ok', x + 16, row_y + 20, w - 32, 40,
                             '✓ 当前无活跃告警\n全园区设备运行正常', color=C_GREEN, font_size=12, z=7))
        return

    for ai, item in enumerate(alarms):
        if row_y + row_h > row_limit:
            break
        dev_action = _nav_action(item['page_id'])
        out.append(make_panel_bg(f'{seed_prefix}-alarm-row-{ai}', x + 10, row_y, w - 20, row_h,
                                 color=C_PANEL_CARD, z=6, opacity=0.88, action=dev_action))
        out.append(make_text(f'{seed_prefix}-alarm-dot-{ai}', x + 16, row_y + 4, 48, 16,
                             '● 离线', color=C_ORANGE, font_size=10, z=7, action=dev_action))
        dev_short = item['name'] if len(item['name']) <= 22 else item['name'][:20] + '…'
        out.append(make_text(f'{seed_prefix}-alarm-dev-{ai}', x + 64, row_y + 3, w - 140, 16,
                             dev_short, color=C_TEXT, font_size=11, z=7, action=dev_action))
        loc_short = item['location'] if len(item['location']) <= 10 else item['location'][:8] + '…'
        out.append(make_text(f'{seed_prefix}-alarm-loc-{ai}', x + w - 72, row_y + 4, 60, 14,
                             loc_short, color=C_TEXT_DIM, font_size=9, z=7, action=dev_action))
        row_y += row_h + 3

    remaining = alarm_n - len(alarms)
    if remaining > 0:
        out.append(make_text(f'{seed_prefix}-alarm-more', x + 12, min(row_y + 2, row_limit - 14), w - 24, 14,
                             f'+{remaining} 更多告警 · 左侧导航树查看', color=C_TEXT_DIM, font_size=10, z=7))


def append_overview_stats_row(out, seed_prefix):
    """经典大屏顶部四卡：总功率 / 今日用电 / 在线设备 / 活跃告警。"""
    stats_y = BODY_Y + 16
    stats_h = 96
    card_gap = 16
    card_w = int((MAIN_W - card_gap * 3) / 4)
    card_xs = [MAIN_X + i * (card_w + card_gap) for i in range(4)]
    # 在线设备 / 活跃告警数值由 ScadaAlarmPanel 实时覆盖，禁止写静态快照。
    stat_configs = [
        ('stat-power', '⚡', '--', '', '总功率', 'activePower', C_ACCENT, False),
        ('stat-energy', '📊', '--', '', '今日用电量', 'todayEnergy', C_BLUE, False),
        ('stat-online', '🖥', '--/--', '', '在线设备', None, C_GREEN, True),
        ('stat-alarm', '🔔', '000', '', '活跃告警', None, C_ORANGE, True),
    ]
    for i, (seed, icon, val, unit, label, energy_role, accent, overlay_value) in enumerate(stat_configs):
        cx = card_xs[i]
        if not overlay_value:
            out.append(make_panel_bg(f'{seed_prefix}-{seed}-glow', cx, stats_y, card_w, stats_h,
                                     color=accent, z=1, opacity=0.06))
            out.append(make_panel_bg(f'{seed_prefix}-{seed}-fill', cx + 5, stats_y + 5, card_w - 10, stats_h - 10,
                                     color=C_PANEL_CARD, z=2, opacity=0.55))
        out.append(make_box12(f'{seed_prefix}-{seed}-bg', cx, stats_y, card_w, stats_h, z=3))
        out.append(make_panel_bg(f'{seed_prefix}-{seed}-accent', cx, stats_y + stats_h - 3, card_w, 2,
                                 color=accent, z=4, opacity=0.85))
        out.append(make_text(f'{seed_prefix}-{seed}-icon', cx, stats_y + 8, card_w, 24, icon,
                             color=accent, font_size=20, z=6, align='center'))
        display_val = f'{val} {unit}'.strip() if unit else val
        if not overlay_value:
            value_cell = make_text(
                f'{seed_prefix}-{seed}-val', cx, stats_y + 32, card_w, 36,
                display_val, color=C_TEXT, font_size=kpi_val_font(card_w), z=6,
                data_bound=False, align='center'
            )
            if energy_role:
                value_cell = _with_energy_overview_role(value_cell, energy_role)
            out.append(value_cell)
        out.append(make_text(
            f'{seed_prefix}-{seed}-lab', cx, stats_y + 70, card_w, 20,
            label, color=C_TEXT_DIM, font_size=FONT_KPI_LABEL, z=6, align='center'
        ))
    out.append(make_panel_bg(f'{seed_prefix}-stats-rule', MAIN_X, stats_y + stats_h + 8, MAIN_W, 1,
                             color=C_BORDER, z=2, opacity=0.7))
    out.append(make_panel_bg(f'{seed_prefix}-stats-accent', MAIN_X, stats_y + stats_h + 7, 280, 2,
                             color=C_ACCENT, z=3, opacity=0.9))
    return stats_y + stats_h + 16


def append_overview_side_panels(out, seed_prefix, panel_top_y, panel_h, right_x, right_w):
    """右侧科技感监测区：功率趋势 + 右下角活跃告警（已去掉用电量趋势）。

    历史查询不再内嵌侧栏（AlarmHistoryComponents 表单过重），由 ScadaAlarmPanel
    标题栏「历史查询」打开全屏抽屉。

    告警区坐标须与 ScadaAlarmPanel.panelStyle 对齐：
      alarm_x = right_x + 8, alarm_w = right_w - 16, alarm_h = 540
      alarm_y = panel_top_y + hdr_h + chart_h + chart_gap
    无顶部 KPI 时标准布局约为 (1312, 516, 584×540)——功率图再矮、告警区再高。
    """
    chart_gap = 12
    alarm_h = 540  # 与 ScadaAlarmPanel.panelStyle 高度对齐（告警区再加高）
    hdr_h = 34  # 顶部留给“● 实时监测”角标的标题带，避免压在边框线上
    # 单功率图：占满告警区以上空间（告警加高后功率图自动变矮）
    chart_h = panel_h - hdr_h - chart_gap - alarm_h - 8

    out.append(make_box13(f'{seed_prefix}-side-frame', right_x, panel_top_y - 4, right_w, panel_h, z=1))
    out.append(make_panel_bg(f'{seed_prefix}-side-glow', right_x + 4, panel_top_y, right_w - 8, panel_h - 8,
                             color=C_ACCENT, z=2, opacity=0.04))
    out.append(make_text(f'{seed_prefix}-side-badge', right_x + right_w - 108, panel_top_y + 12, 92, 18,
                         '● 实时监测', color=C_GREEN, font_size=11, z=8, align='center'))

    power_y = panel_top_y + hdr_h
    out.append(make_box13(f'{seed_prefix}-chart-power-inner', right_x + 8, power_y, right_w - 16, chart_h - 8, z=3))
    out.extend(make_panel_title(f'{seed_prefix}-chart-power', right_x + 20, power_y + 10,
                                '功率趋势 (24h)', color=C_ACCENT, font_size=FONT_PANEL, z=7, w=200))
    out.append(_with_energy_overview_role(make_smooth_chart(
        f'{seed_prefix}-chart-trend', right_x + 18, power_y + 36, right_w - 36, chart_h - 48,
        title='功率趋势 (24h)',
        dp_names=['总有功功率', '总无功功率', '总视在功率'],
        z=6, show_title=False
    ), 'power24h'))

    alarm_y = power_y + chart_h + chart_gap
    append_overview_alarm_panel(out, seed_prefix, right_x + 8, alarm_y, right_w - 16, alarm_h)


def append_overview_main_body(out, seed_prefix, substations):
    """首页主体：左运行时组织总览外框 + 右功率趋势/告警（无顶部 KPI）。"""
    # RC08bate-20260724：删除顶部四卡，左区结构图上移占满
    content_top = BODY_Y + 16
    panel_h = 1080 - content_top - 16
    left_w = int((MAIN_W - 16) * 0.68)
    right_w = MAIN_W - left_w - 16
    right_x = MAIN_X + left_w + 16
    # 组织内容由 ScadaOrgOverview 在运行时根据 monitortree 递归生成。
    # 这里只保留唯一外框，禁止再生成母线、馈线、设备组等旧静态组件。
    out.append(make_box13(
        f'{seed_prefix}-frame', MAIN_X, content_top - 4, left_w, panel_h + 4, z=1,
    ))
    append_overview_side_panels(out, seed_prefix, content_top, panel_h, right_x, right_w)


def append_campus_oneline_diagram(out, seed_prefix, substations, *, show_back=False, back_page_id=None,
                                  frame_x=MAIN_X, frame_y=None, frame_w=MAIN_W, frame_h=None,
                                  bottom_reserve=0, legend_offset_x=0, reserve_rect=None):
    """园区变电所一次系统总图（dv-border-box-13 内）：设计图1 网格示意，点击下钻单所一次图。"""
    cy0 = frame_y if frame_y is not None else BODY_Y + 16
    if frame_h is None:
        frame_h = 1080 - (cy0 - 4) - 16
    out.append(make_box13(f'{seed_prefix}-frame', frame_x, cy0 - 4, frame_w, frame_h, z=1))
    out.append(make_panel_bg(f'{seed_prefix}-frame-glow', frame_x + 6, cy0 + 2, frame_w - 12, frame_h - 12,
                             color=C_ACCENT, z=2, opacity=0.035))

    title_y = cy0 + LEVEL_TITLE_TOP_PAD
    title_x = frame_x + 16
    if show_back and back_page_id:
        out.append(make_text(f'{seed_prefix}-back', frame_x + 16, title_y, 120, LEVEL_TITLE_H,
                             '← 返回总览', color=C_ACCENT, font_size=14, z=20,
                             action=_nav_action(back_page_id)))
        title_x = frame_x + 148

    sub_n = len(substations)
    dev_n = sum(s['device_count'] for s in substations)
    title_w = frame_x + frame_w - title_x - 16
    out.append(make_text(f'{seed_prefix}-title', title_x, title_y, title_w,
                         LEVEL_TITLE_H, '组织层级总览',
                         color=C_TEXT, font_size=FONT_PANEL, z=10))
    out.append(make_text(f'{seed_prefix}-sub', title_x, title_y + LEVEL_TITLE_H + 4,
                         title_w, 20,
                         f'全园区 {sub_n} 个顶级区域 / {dev_n} 台设备 · 组织节点仅展开，设备叶节点可查看测点',
                         color=C_TEXT_DIM, font_size=12, z=10))

    grid_y = title_y + LEVEL_TITLE_H + 28
    inner_w = frame_w - 32
    card_w, card_h = (196, 178) if frame_w >= 1100 else (168, 156)
    gap = 12
    per_row = max(1, (inner_w + gap) // (card_w + gap))
    footer_h = 32
    frame_bottom = cy0 - 4 + frame_h
    # 有左下角拓扑概览(reserve_rect)时按矩形跳过被遮挡的槽位、整框铺满；
    # 否则用 bottom_reserve 预留底部高度（独立总览页无拓扑浮层）。
    if reserve_rect:
        avail_h = frame_bottom - grid_y - footer_h - 8
    else:
        avail_h = frame_bottom - grid_y - footer_h - 12 - bottom_reserve
    max_rows = max(1, (avail_h + gap) // (card_h + gap))
    base_x = frame_x + 16

    # 行优先生成槽位；跳过与左下角拓扑概览重叠的槽位（左下角留给拓扑）
    slots = []
    for r in range(max_rows):
        for c in range(per_row):
            cx = base_x + c * (card_w + gap)
            cy = grid_y + r * (card_h + gap)
            if reserve_rect and rects_overlap(
                    {'x': cx, 'y': cy, 'w': card_w, 'h': card_h},
                    {'x': reserve_rect[0], 'y': reserve_rect[1],
                     'w': reserve_rect[2], 'h': reserve_rect[3]}):
                continue
            slots.append((cx, cy))

    shown = substations[:len(slots)]
    for si, sub in enumerate(shown):
        cx, cy = slots[si]
        ok = sub['alarm'] == 0
        accent = C_GREEN if ok else C_ORANGE
        _append_campus_substation_card(out, seed_prefix, si, cx, cy, card_w, card_h, sub, accent)

    if len(substations) > len(shown):
        # “+N 更多”放到拓扑概览右侧的底部空白带，避免与卡片/拓扑重叠
        more_x = (reserve_rect[0] + reserve_rect[2] + 16) if reserve_rect else (frame_x + 16)
        more_w = frame_x + frame_w - 16 - more_x
        out.append(make_text(f'{seed_prefix}-more', more_x, frame_bottom - 50, more_w, 18,
                             f'+{len(substations) - len(shown)} 更多区域 · 共 {sub_n} 个 · 左侧导航树可逐级查看',
                             color=C_TEXT_DIM, font_size=12, z=6))
    legend_y = frame_bottom - 30
    legend_x = frame_x + 16 + legend_offset_x
    legend_w = inner_w - legend_offset_x
    out.append(make_text(f'{seed_prefix}-legend', legend_x, legend_y, legend_w, 18,
                         '● 与设备管理树实时组织层级一致 · 每个节点展示其包含设备数量',
                         color=C_TEXT_DIM, font_size=11, z=6))


def build_oneline_campus_cells(substations, seed_prefix='oneline-master'):
    """全园区变电所一次系统总览独立页（面包屑返回用，内容与首页总图一致）。"""
    out = []
    out.extend(build_header_cells(seed_prefix, [
        ('📊 全局总览', C_TEXT_DIM, PAGE_ID_MAIN),
        ('🔌 变电所一次系统总图', C_TEXT, None),
    ]))
    append_campus_oneline_diagram(out, seed_prefix, substations,
                                  show_back=True, back_page_id=PAGE_ID_MAIN)
    return out


# ════════════════════════════════════════════════════════
# LEVEL 0: OVERVIEW PAGE (page_name='main') — 设计图1 园区一次系统总图
# ════════════════════════════════════════════════════════

cells = []
cells.extend(build_header_cells('ov', [
    ('📊 全局总览', C_TEXT, PAGE_ID_MAIN),
]))
append_overview_main_body(cells, 'ov', substations)
report_overlaps(cells, 'Overview layout')

components_json_main = json.dumps({"cells": cells}, ensure_ascii=False)
comp_b64_main = base64.b64encode(components_json_main.encode()).decode()

print(f"\n=== LEVEL 0: OVERVIEW PAGE ===")
print(f"Total cells: {len(cells)}")
print(f"JSON size: {len(components_json_main)} chars")
print(f"Base64 size: {len(comp_b64_main)} chars")

decoded = json.loads(base64.b64decode(comp_b64_main).decode())
assert 'cells' in decoded and len(decoded['cells']) == len(cells)
print("Roundtrip verification PASSED")

DISPLAY_NAME = '中航信数据中心电力监控系统'
cur.execute(
    "UPDATE display_models SET name=%s, updated_at=NOW() WHERE display_model_uid=%s",
    (DISPLAY_NAME, MODEL_ID)
)
conn.commit()
print(f"Display name updated to '{DISPLAY_NAME}', rows affected: {cur.rowcount}")

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
    # 保留既有 page_name / template_kind（如「首页模板」+ home），只刷新图层内容
    cur.execute(
        "UPDATE display_model_layer SET components=%s, layer=%s, page_id=%s, is_home=1, updated_at=NOW() WHERE id=%s",
        (comp_b64_main, _dark_layer, PAGE_ID_MAIN, overview_id),
    )
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


def build_oneline_cells(sub, seed_prefix='oneline'):
    """单座变电所一次系统图：进线 → 10kV 母线 → 全部馈线(机柜)网格。"""
    out = []
    out.extend(build_header_cells(seed_prefix, [
        ('📊 全局总览', C_TEXT_DIM, PAGE_ID_MAIN),
        ('🔌 变电所一次系统总图', C_TEXT_DIM, PAGE_ID_MAIN),
        (sub['name'], C_TEXT, None),
    ]))
    cy0 = BODY_Y + 16
    main_right = MAIN_X + MAIN_W
    frame_h = 1080 - (cy0 - 4) - 16
    out.append(make_box13(f'{seed_prefix}-frame', MAIN_X, cy0 - 4, MAIN_W, frame_h, z=1))

    # 标题行：加高文本框，避免长配电室名被压线/裁切
    title_y = cy0 + 6
    out.append(make_text(f'{seed_prefix}-back', MAIN_X + 16, title_y, 120, 30, '← 返回总图',
                         color=C_ACCENT, font_size=14, z=20, action=_nav_action(PAGE_ID_MAIN)))
    out.append(make_text(f'{seed_prefix}-title', MAIN_X + 148, title_y, MAIN_W - 164, 30,
                         f'🔌 {sub["name"]} · 一次系统图',
                         color=C_TEXT, font_size=17, z=10))
    cab_n = sub['cabinet_count']
    zone_hint = f'（{sub["zone_name"]}）' if sub.get('zone_name') and sub.get('zone_name') != sub.get('name') else ''
    out.append(make_text(f'{seed_prefix}-sub', MAIN_X + 148, title_y + 34, MAIN_W - 164, 20,
                         f'{zone_hint}进线 → 10kV 母线 → {cab_n} 路馈线 · 点击卡片下钻机柜',
                         color=C_TEXT_DIM, font_size=12, z=10))

    cabs = sub['cabinets']
    source_device = next((cab['devices'][0] for cab in cabs if cab.get('devices')), None)
    diagram_top = title_y + 64

    # ── 进线 + 母线（紧凑示意）──
    inc_cx = (MAIN_X + main_right) // 2
    inc_y = diagram_top + 6
    sym = 64
    out.append(make_electric('view-svg-electric6', f'{seed_prefix}-inc', inc_cx - sym // 2, inc_y,
                             sym, sym, color=C_ACCENT, z=8, device=source_device))
    out.append(make_text(f'{seed_prefix}-inc-lab', inc_cx + sym // 2 + 6, inc_y + 16, 200, 18,
                         '⫶ 10kV 进线柜', color=C_TEXT, font_size=13, z=10))
    bus_y = inc_y + sym + 32
    bus_l = MAIN_X + 60
    bus_r = main_right - 60
    out.append(make_conn_line(f'{seed_prefix}-inc-drop', inc_cx - 1, inc_y + sym, 3,
                              bus_y - (inc_y + sym), color=C_ACCENT, z=4))
    out.append(make_move_line(f'{seed_prefix}-flow-inc', inc_cx - 10, inc_y + sym, 20,
                              bus_y - (inc_y + sym), color=C_ACCENT, z=5, vertical=True,
                              stroke_width=2, direction=0, device=source_device))
    out.append(make_conn_line(f'{seed_prefix}-bus', bus_l, bus_y - 2, bus_r - bus_l, 4,
                              color=C_ACCENT, z=4, opacity=0.95))
    out.append(make_move_line(f'{seed_prefix}-flow-bus', bus_l, bus_y - 10, bus_r - bus_l, 20,
                              color=C_ACCENT, z=5, vertical=False, stroke_width=2, direction=0,
                              device=source_device))
    out.append(make_text(f'{seed_prefix}-bus-lab', bus_l + 4, bus_y - 24, 240, 16,
                         f'▭ 10kV I 段母线 · {cab_n} 路馈线', color=C_ACCENT, font_size=11, z=10))

    # ── 全部馈线：网格展示（不再限制 3 路）──
    grid_y = bus_y + 24
    card_w, card_h = 252, 200
    gap = 14
    per_row = max(1, (MAIN_W - 32 + gap) // (card_w + gap))
    footer_h = 28
    max_rows = max(1, (1080 - grid_y - footer_h - 16) // (card_h + gap))
    cap = per_row * max_rows
    shown_cabs = cabs[:cap]

    for bi, cab in enumerate(shown_cabs):
        r, c = bi // per_row, bi % per_row
        cardx = MAIN_X + 16 + c * (card_w + gap)
        cardy = grid_y + r * (card_h + gap)
        cab_action = _nav_action(cab['page_id'])
        rep = cab['devices'][0] if cab['devices'] else None
        shape = _ONELINE_BRANCH_SHAPES[bi % len(_ONELINE_BRANCH_SHAPES)]
        bx = cardx + card_w // 2
        sym_s = 44
        sym_y = cardy + 6
        out.append(make_electric(shape, f'{seed_prefix}-sym-{bi}', bx - sym_s // 2, sym_y,
                                 sym_s, sym_s, color=C_ACCENT, action=cab_action, z=8, device=rep))
        out.append(make_text(f'{seed_prefix}-sym-lab-{bi}', cardx + 8, sym_y + sym_s + 2, card_w - 16, 18,
                             f'🗄 {cab["name"]}', color=C_ACCENT, font_size=12, z=10, action=cab_action))
        data_y = sym_y + sym_s + 22
        data_h = cardy + card_h - data_y - 6
        out.append(make_panel_bg(f'{seed_prefix}-card-bg-{bi}', cardx + 4, data_y, card_w - 8, data_h,
                                 color=C_PANEL_CARD, z=2, opacity=0.9, action=cab_action))
        out.append(make_box13(f'{seed_prefix}-card-fr-{bi}', cardx + 4, data_y, card_w - 8, data_h, z=3))
        on = rep and rep.get('status') == 1
        out.append(make_text(f'{seed_prefix}-card-st-{bi}', cardx + card_w - 72, data_y + 8, 60, 16,
                             '● 运行' if on else '● 离线',
                             color=C_GREEN if on else C_TEXT_DIM, font_size=11, z=6))
        if rep:
            rows = _oneline_points_for(rep.get('muid'), rep.get('uuid'))
            ry = data_y + 28
            for li, (lab, dpn, unit) in enumerate(rows[:3]):
                yy = ry + li * 36
                out.append(make_text(f'{seed_prefix}-rk-{bi}-{li}', cardx + 12, yy, 80, 18,
                                     lab, color=C_TEXT_DIM, font_size=11, z=6))
                out.append(make_text(f'{seed_prefix}-rv-{bi}-{li}', cardx + 96, yy, 100, 20,
                                     '—', color=C_ACCENT if li < 2 else C_GREEN,
                                     font_size=FONT_PARAM_VAL, z=6,
                                     data_bound=True, dp_name=dpn, device=rep))
                if unit:
                    out.append(make_text(f'{seed_prefix}-ru-{bi}-{li}', cardx + card_w - 48, yy, 36, 16,
                                         unit, color=C_TEXT_DIM, font_size=10, z=6))
        else:
            out.append(make_text(f'{seed_prefix}-nodata-{bi}', cardx + 12, data_y + 36, card_w - 24, 32,
                                 '⚠ 无数据回路', color=C_TEXT_DIM, font_size=12, z=6))

    if len(cabs) > len(shown_cabs):
        last_row = (len(shown_cabs) - 1) // per_row
        foot_y = grid_y + (last_row + 1) * (card_h + gap) + 4
        out.append(make_text(f'{seed_prefix}-more', MAIN_X + 16, min(foot_y, 1080 - 40), MAIN_W - 32, 18,
                             f'+{len(cabs) - len(shown_cabs)} 更多馈线 · 共 {cab_n} 柜 · 左侧导航树可逐级查看',
                             color=C_TEXT_DIM, font_size=12, z=6))

    legend_y = min(grid_y + ((len(shown_cabs) - 1) // per_row + 1) * (card_h + gap) + 8, 1080 - 32)
    out.append(make_text(f'{seed_prefix}-legend', MAIN_X + 16, legend_y, MAIN_W - 32, 18,
                         '● 实时值取各路代表设备 · 点击符号或卡片下钻机柜 · 流动箭头表潮流方向',
                         color=C_TEXT_DIM, font_size=11, z=6))
    return out


def build_zone_detail_cells(zone, seed_prefix='zone'):
    """One page per 顶级区域 — 子区域卡片与设备管理树一致。"""
    out = []
    tree_items = _zone_tree_items(zone)
    out.extend(build_header_cells(seed_prefix, [
        ('📊 全局总览', C_TEXT_DIM, PAGE_ID_MAIN),
        (zone['name'], C_TEXT, None),
    ]))
    level_y = BODY_Y + 16
    title_y = level_y + LEVEL_TITLE_TOP_PAD
    out.append(make_box13(f'{seed_prefix}-content-frame', MAIN_X, level_y - 4,
                          MAIN_W, 1080 - (level_y - 4) - 16, z=1))
    out.append(make_text(f'{seed_prefix}-back-btn', MAIN_X + 16, title_y, 140, 32, '← 返回总览',
                         color=C_ACCENT, font_size=14, z=20, action=_nav_action(PAGE_ID_MAIN)))
    out.append(make_text(f'{seed_prefix}-level-title', MAIN_X + 176, title_y, 760, LEVEL_TITLE_H,
                         f'{zone["name"]}', color=C_TEXT, font_size=22, z=10))
    sub_cnt = _zone_child_type0_count(zone['sid'])
    sub_hint = f'{sub_cnt}子区域 · ' if sub_cnt else ''
    out.append(make_text(f'{seed_prefix}-level-sub', MAIN_X + 176, title_y + LEVEL_SUB_Y, 760, 18,
                         f'{sub_hint}{zone["device_count"]}台设备 · 在线{zone["online"]} · 异常{zone["alarm"]}',
                         color=C_TEXT_DIM, font_size=13, z=10))
    card_w_item = 300
    card_h_item = 150
    gap = 16
    per_row = max(1, (MAIN_W - 32 + gap) // (card_w_item + gap))
    start_x = MAIN_X + 16
    start_y = title_y + LEVEL_GRID_Y
    ITEM_CAP = 24
    if not tree_items:
        out.append(make_text(f'{seed_prefix}-empty', start_x, start_y + 40, MAIN_W - 32, 40,
                             '暂无设备或未配置子区域', color=C_TEXT_DIM, font_size=14, z=5))
        return out
    for ri, item in enumerate(tree_items[:ITEM_CAP]):
        r, c = ri // per_row, ri % per_row
        cx = start_x + c * (card_w_item + gap)
        cy = start_y + r * (card_h_item + gap)
        if cy + card_h_item > 1080 - 40:
            break
        item_nav = _nav_action(item['page_id'])
        ok = item['alarm'] == 0
        accent = C_GREEN if ok else C_ORANGE
        out.append(make_panel_bg(f'{seed_prefix}-room-bg-{ri}', cx, cy, card_w_item, card_h_item,
                                 color=C_PANEL_CARD, z=3, opacity=0.9, action=item_nav))
        out.append(make_panel_bg(f'{seed_prefix}-room-bar-{ri}', cx, cy, 4, card_h_item, color=accent, z=4))
        out.append(make_text(f'{seed_prefix}-room-name-{ri}', cx + 16, cy + 14, card_w_item - 32, 24,
                             item['display_name'], color=C_ACCENT, font_size=16, z=5, action=item_nav))
        out.append(make_text(f'{seed_prefix}-room-sub-{ri}', cx + 16, cy + 48, card_w_item - 32, 18,
                             f'{item["device_count"]}台设备', color=C_TEXT_DIM, font_size=12, z=5, action=item_nav))
        metrics = [('设备', item['device_count'], C_TEXT), ('在线', item['online'], C_GREEN),
                   ('异常', item['alarm'], C_ORANGE if item['alarm'] else C_TEXT_DIM)]
        mw = (card_w_item - 32) // 3
        for mi, (mlabel, mval, mcolor) in enumerate(metrics):
            mx = cx + 16 + mi * mw
            out.append(make_text(f'{seed_prefix}-room-mv-{ri}-{mi}', mx, cy + 80, mw - 8, 28,
                                 str(mval), color=mcolor, font_size=24, z=5, action=item_nav))
            out.append(make_text(f'{seed_prefix}-room-ml-{ri}-{mi}', mx, cy + 116, mw - 8, 18,
                                 mlabel, color=C_TEXT_DIM, font_size=11, z=5, action=item_nav))
    if len(tree_items) > ITEM_CAP:
        out.append(make_text(f'{seed_prefix}-room-more', start_x, 1080 - 36, MAIN_W - 32, 18,
                             f'+{len(tree_items) - ITEM_CAP} 更多子区域 · 共 {len(tree_items)} 项',
                             color=C_TEXT_DIM, font_size=12, z=5))
    return out


def build_room_detail_cells(room, seed_prefix='room'):
    """One page per 配电室/楼层 — ViewRealTable 实时数据表格（新四字段绑定，深色工业风）。"""
    out = []
    out.extend(build_header_cells(seed_prefix, [
        ('📊 全局总览', C_TEXT_DIM, PAGE_ID_MAIN),
        (room['display_name'], C_TEXT, None),
    ]))
    level_y = BODY_Y + 16
    title_y = level_y + LEVEL_TITLE_TOP_PAD
    content_h = 1080 - (level_y - 4) - 16
    out.append(make_box13(f'{seed_prefix}-content-frame', MAIN_X, level_y - 4,
                          MAIN_W, content_h, z=1))
    out.append(make_text(f'{seed_prefix}-back-btn', MAIN_X + 16, title_y, 140, 32, '← 返回总览',
                         color=C_ACCENT, font_size=14, z=20, action=_nav_action(PAGE_ID_MAIN)))
    out.append(make_text(f'{seed_prefix}-level-title', MAIN_X + 176, title_y, 760, LEVEL_TITLE_H,
                         f'🏛 {room["display_name"]}', color=C_TEXT, font_size=22, z=10))
    tbl_cfg = build_room_table_config(room)
    sub_text = (f'{room["cabinet_count"]}个机柜 · {room["device_count"]}台设备 · '
                f'在线{room["online"]} · 异常{room["alarm"]}')
    if tbl_cfg:
        sub_text += f' · 表格展示 {tbl_cfg["device_count"]} 台'
    out.append(make_text(f'{seed_prefix}-level-sub', MAIN_X + 176, title_y + LEVEL_SUB_Y, 960, 18,
                         sub_text, color=C_TEXT_DIM, font_size=13, z=10))
    table_x = MAIN_X + 16
    table_y = title_y + LEVEL_GRID_Y
    table_w = MAIN_W - 32
    table_h = 1080 - table_y - 20
    if tbl_cfg:
        out.append(make_real_table(
            f'{seed_prefix}-real-table', table_x, table_y, table_w, table_h,
            column_headers=tbl_cfg['column_headers'],
            row_names=tbl_cfg['row_names'],
            row_codes=tbl_cfg['row_codes'],
            binding_rows=tbl_cfg['bindings'],
            z=8, theme='dark', page_size=ROOM_TABLE_PAGE_SIZE,
        ))
    else:
        out.append(make_text(f'{seed_prefix}-no-data', table_x, table_y + 40, table_w, 40,
                             '暂无设备数据', color=C_TEXT_DIM, font_size=16, z=5))
    return out


GENERATE_LEGACY_PAGES = os.environ.get('NCC_GENERATE_LEGACY_PAGES') == '1'
zone_pages = []
room_pages = []
building_pages = []
floor_pages = []
detail_cells = []
device_pages = 0

if not GENERATE_LEGACY_PAGES:
    print("\n=== LEGACY PAGES SKIPPED (default; set NCC_GENERATE_LEGACY_PAGES=1 to regenerate) ===")
    print("  Runtime templates only: 首页模板 / 设备列表模板 / 点位列表模板")
    cur.execute(
        """
        DELETE FROM display_model_layer
        WHERE model_id=%s
          AND COALESCE(is_home,0)<>1
          AND COALESCE(page_name,'') NOT IN ('首页模板','设备列表模板','点位列表模板')
          AND COALESCE(template_kind,'') NOT IN ('home','deviceList','datapointList')
        """,
        (MODEL_ID,),
    )
    conn.commit()
    print(f"  HARD-deleted leftover legacy pages: {cur.rowcount}")
else:
    print(f"\n=== LEVEL Z: ZONE PAGES (per 顶级区域) ===")
    zone_pages = []
    for zone in zones:
        if zone['sid'] == -1:
            continue
        seed = f'zone-{zone["sid"]}'
        zone_cells = build_zone_detail_cells(zone, seed_prefix=seed)
        report_overlaps(zone_cells, f'Zone {zone["name"]}')
        comp_b64 = base64.b64encode(json.dumps({"cells": zone_cells}, ensure_ascii=False).encode()).decode()
        row_id, action = upsert_layer_page(f'zone-{zone["sid"]}', zone['page_id'], comp_b64)
        zone_pages.append((zone['name'], zone['page_id'], len(zone_cells), action))
        print(f"  {action} zone page: {zone['name']} id={row_id} cells={len(zone_cells)}")
    conn.commit()

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

    # ── LEVEL O: 一次系统总图（园区总览 + 各变电所单线图）──
    if substations and PAGE_ID_ONELINE:
        print(f"\n=== LEVEL O: ONE-LINE DIAGRAM (一次系统总图) ===")
        master_cells = build_oneline_campus_cells(substations)
        report_overlaps(master_cells, 'One-line campus master')
        _ol_master_b64 = base64.b64encode(json.dumps({"cells": master_cells}, ensure_ascii=False).encode()).decode()
        _ol_master_id, _ol_master_act = upsert_layer_page('oneline', PAGE_ID_ONELINE, _ol_master_b64)
        conn.commit()
        print(f"  {_ol_master_act} one-line MASTER: {len(substations)} substations id={_ol_master_id} "
              f"cells={len(master_cells)} page_id={PAGE_ID_ONELINE}")
        for sub in substations:
            seed = f'oneline-{sub["sid"]}'
            sub_cells = build_oneline_cells(sub, seed_prefix=seed)
            report_overlaps(sub_cells, f'One-line {sub["name"]}')
            _ol_b64 = base64.b64encode(json.dumps({"cells": sub_cells}, ensure_ascii=False).encode()).decode()
            _ol_id, _ol_act = upsert_layer_page(seed, sub['page_id'], _ol_b64)
            print(f"  {_ol_act} one-line sub: {sub['name']} id={_ol_id} cells={len(sub_cells)}")
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
        # 按本设备的模型选点（UPS 与电力仪表点名体系不同），否则实时/功率绑点全空。
        rt_params, pw_params, chart_dps, dev_type_label = detail_params_for(dev.get('muid'))
        out.extend(build_header_cells(seed_prefix, [
            ('📊 全局总览', C_TEXT_DIM, PAGE_ID_MAIN),
            (bldg.get('room_name', ROOT_NAME), C_TEXT_DIM, bldg.get('room_page_id')),
            (bldg['name'], C_TEXT_DIM, bldg['page_id']),
            (floor['name'], C_TEXT_DIM, floor['page_id']),
            (dev_name, C_TEXT, None),
        ]))
        level_y = BODY_Y + 16
        title_y = level_y + LEVEL_TITLE_TOP_PAD
        out.append(make_text(f'{seed_prefix}-back', MAIN_X, title_y, 200, 32, f'← {floor["name"]}',
                             color=C_ACCENT, font_size=14, z=20, action=_nav_action(floor['page_id'])))
        out.append(make_text(f'{seed_prefix}-title', MAIN_X + 220, title_y, 760, LEVEL_TITLE_H,
                             f'🔧 {dev_name}', color=C_TEXT, font_size=22, z=10))
        on = dev['status'] == 1
        out.append(make_text(f'{seed_prefix}-status', MAIN_X + 980, title_y + 6, 200, 22,
                             '● 运行中' if on else '● 离线', color=C_GREEN if on else C_TEXT_DIM,
                             font_size=13, z=10))
        panel_top = title_y + LEVEL_GRID_Y
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
            ('设备名称', dev_name), ('设备类型', dev_type_label),
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
        device_points = dp_map_for(dev.get('muid'), dev.get('uuid'))
        chart_points = [point for point in chart_dps if _make_active(
            point, dev.get('uuid'), dev.get('name'), device_points
        )]
        if chart_points:
            out.append(make_smooth_chart(
                f'{seed_prefix}-chart', MAIN_X + 16, chart_y + 40, chart_w - 32, 226,
                title='设备功率趋势', dp_names=chart_points, z=5, device=dev))
        else:
            out.append(make_text(
                f'{seed_prefix}-chart-unavailable', MAIN_X + 16, chart_y + 92, chart_w - 32, 32,
                '当前设备未接入可用的功率趋势测点', color=C_TEXT_DIM, font_size=13, z=6, align='center',
            ))
        status_x = MAIN_X + chart_w + 16
        status_w = MAIN_W - chart_w - 16
        out.append(make_panel_bg(f'{seed_prefix}-sp', status_x, chart_y, status_w, 280, color=C_PANEL, z=2, opacity=0.6))
        out.append(make_box13(f'{seed_prefix}-sf', status_x, chart_y, status_w, 280, z=3))
        out.append(make_text(f'{seed_prefix}-stt', status_x + 16, chart_y + 10, 400, 22, '🔔 设备告警',
                             color=C_ORANGE, font_size=FONT_PANEL + 2, z=6))
        out.append(make_text(f'{seed_prefix}-se', status_x + 16, chart_y + 48, status_w - 32, 200,
                             '✅ 该设备无告警记录', color=C_TEXT_DIM, font_size=13, z=6))
        return out

    GENERATE_LEGACY_DEVICE_PAGES = os.environ.get('NCC_GENERATE_LEGACY_DEVICE_PAGES') == '1'
    device_pages = 0
    if GENERATE_LEGACY_DEVICE_PAGES:
        for bldg in buildings:
            for floor in bldg['floors']:
                for dev in floor['devices']:
                    seed = f'dev-{dev["sid"]}'
                    dcells = build_device_detail_cells(dev, bldg, floor, seed_prefix=seed)
                    comp_b64 = base64.b64encode(json.dumps({"cells": dcells}, ensure_ascii=False).encode()).decode()
                    upsert_layer_page(f'device-{dev["sid"]}', page_id_device(dev['sid']), comp_b64)
                    device_pages += 1
    conn.commit()
    if GENERATE_LEGACY_DEVICE_PAGES:
        report_overlaps(build_device_detail_cells(buildings[0]['floors'][0]['devices'][0],
                                                  buildings[0], buildings[0]['floors'][0], 'dev-sample'),
                        'Device detail (sample)')
    print(f"  wrote {device_pages} per-device detail pages"
          f"{'' if GENERATE_LEGACY_DEVICE_PAGES else ' (disabled: runtime templates are used)'}")

    # Keep the legacy shared 'device-detail' page (id=10) pointing at a real device
    # so any old reference still resolves to a coherent page.
    detail_cells = []
    if buildings and buildings[0].get('floors') and buildings[0]['floors'][0].get('devices'):
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
    else:
        print("  skip legacy device-detail page (no devices in hierarchy)")

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
print(f"   Level 0 (overview/home template): {len(cells)} cells")
if GENERATE_LEGACY_PAGES:
    print(f"   Level Z/R/O/1/2 legacy pages: regenerated (NCC_GENERATE_LEGACY_PAGES=1)")
    print(f"   Level 1 (building pages):  {len(building_pages)} pages, {total_bldg_cells} cells total")
    print(f"   Level 2 (floor pages):     {len(floor_pages)} pages, {total_floor_cells} cells total")
    print(f"   Level 3 (device-detail):   {len(detail_cells)} cells")
else:
    print(f"   Legacy zone/room/building/floor/device pages: SKIPPED + HARD-purged")
    print(f"   Runtime: 首页模板 / 设备列表模板 / 点位列表模板")
print(f"   Main content: MAIN_X={MAIN_X}px, MAIN_W={MAIN_W}px (sidebar removed, use ISMRunTreeNav)")
print(f"\n   Preview: {apprun_url(MODEL_ID)}")
print(f"{'='*60}")

conn.close()
