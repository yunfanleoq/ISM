"""
ISM Automation — 设备层级构建器
从设备列表中构建机房 → 区域 → 配电室 → 柜 → 设备组的层级树。

解析逻辑继承自 build_ncc_dashboard.py 中的层级计算部分。
"""
from __future__ import annotations

import re
from collections import defaultdict
from typing import Any, Dict, List, Optional, Tuple

from ism_automation.core.logger import ism_logger


class HierarchyBuilder:
    """
    设备层级构建器

    输入: 设备列表 (from DB or from Excel)
    输出: 分层结构树

    层级结构:
    Root → Zone(区域) → Room(配电室) → Cabinet(柜) → FloorGroup(设备组) → Device(设备)
    """

    def __init__(self, devices: List[Dict[str, Any]]):
        """
        Args:
            devices: 设备列表，每个设备包含:
                - uuid, name, sid, pid, type, muid, status
        """
        self.devices = devices
        self.device_by_sid: Dict[int, Dict[str, Any]] = {}
        self.children_by_pid: Dict[int, List[Dict[str, Any]]] = defaultdict(list)

        self._build_index()

    def _build_index(self):
        """构建父子索引"""
        for dev in self.devices:
            sid = dev.get('sid')
            pid = dev.get('pid', 0)
            if sid is not None:
                self.device_by_sid[sid] = dev
                self.children_by_pid[pid].append(dev)

        ism_logger.info(f"📊 层级索引: {len(self.device_by_sid)} 个节点, "
                       f"{len(self.children_by_pid)} 个父节点")

    # ── 层级解析 ──────────────────────────────────────

    def get_root(self) -> Optional[Dict[str, Any]]:
        """获取根节点（pid 不在设备列表中的 type=0 节点）"""
        for sid, dev in self.device_by_sid.items():
            if dev.get('type') == 0 and dev.get('pid') not in self.device_by_sid:
                return dev
        return None

    def get_zones(self, root_sid: Optional[int] = None) -> List[Dict[str, Any]]:
        """获取 Root 直下的区域节点"""
        if root_sid is None:
            root = self.get_root()
            root_sid = root.get('sid') if root else 0

        zones = []
        for child in self.children_by_pid.get(root_sid, []):
            if child.get('type') == 0:
                zones.append(child)
        return zones

    def get_rooms(self, zone_sid: int) -> List[Dict[str, Any]]:
        """获取区域下的配电室节点"""
        rooms = []
        for child in self.children_by_pid.get(zone_sid, []):
            if child.get('type') == 0:
                rooms.append(child)
        return rooms

    def get_cabinets(self, room_sid: int) -> List[Dict[str, Any]]:
        """获取配电室下的柜节点（type=0 且有 type=1 子节点）"""
        cabinets = []
        for child in self.children_by_pid.get(room_sid, []):
            if child.get('type') == 0:
                # 检查是否有 type=1 子节点（实际设备）
                grand_children = self.children_by_pid.get(child.get('sid'), [])
                type1_children = [g for g in grand_children if g.get('type') == 1]
                if type1_children:
                    child['devices'] = type1_children
                    child['device_count'] = len(type1_children)
                    cabinets.append(child)
        return cabinets

    def get_floor_groups(self, cabinet_sid: int) -> List[Dict[str, Any]]:
        """获取柜下的设备组（按设备名称前缀分组，如 S18）"""
        devices = [
            d for d in self.children_by_pid.get(cabinet_sid, [])
            if d.get('type') == 1
        ]

        groups = defaultdict(list)
        for d in devices:
            name = d.get('name', '')
            parts = name.split('_')
            if len(parts) >= 3:
                floor_key = parts[2]  # e.g., "S18"
            else:
                floor_key = 'default'
            groups[floor_key].append(d)

        return [
            {
                'key': k,
                'name': f'{k}设备组',
                'devices': v,
                'count': len(v),
            }
            for k, v in sorted(groups.items())
        ]

    # ── 构建完整层级树 ────────────────────────────────

    def build_tree(self) -> Dict[str, Any]:
        """
        构建完整层级树

        返回格式:
        {
            "root": {...},
            "zones": [
                {
                    "sid": 1, "name": "1A", "uuid": "...",
                    "rooms": [
                        {
                            "sid": 10, "name": "1A1配电室",
                            "cabinets": [
                                {
                                    "sid": 100, "name": "1A1_U11柜",
                                    "device_count": 18,
                                    "floor_groups": [
                                        {"key": "S18", "name": "S18设备组", "count": 9, "devices": [...]},
                                        ...
                                    ]
                                },
                                ...
                            ]
                        },
                        ...
                    ]
                },
                ...
            ]
        }
        """
        root = self.get_root()
        root_sid = root.get('sid') if root else 0

        zones = []
        for zone in self.get_zones(root_sid):
            zone_sid = zone.get('sid')
            rooms = []

            for room in self.get_rooms(zone_sid):
                room_sid = room.get('sid')
                cabinets = self.get_cabinets(room_sid)

                for cabinet in cabinets:
                    cabinet_sid = cabinet.get('sid')
                    cabinet['floor_groups'] = self.get_floor_groups(cabinet_sid)
                    cabinet['online'] = sum(1 for d in cabinet['devices'] if d.get('status') == 1)
                    cabinet['alarm'] = cabinet['device_count'] - cabinet['online']

                rooms.append({
                    'sid': room_sid,
                    'name': room.get('name'),
                    'uuid': room.get('uuid'),
                    'cabinets': cabinets,
                    'device_count': sum(c['device_count'] for c in cabinets),
                })

            zones.append({
                'sid': zone_sid,
                'name': zone.get('name'),
                'uuid': zone.get('uuid'),
                'rooms': rooms,
                'device_count': sum(r['device_count'] for r in rooms),
            })

        tree = {
            'root': root,
            'zones': zones,
            'total_devices': sum(z['device_count'] for z in zones),
        }

        ism_logger.info(f"🌳 层级树构建完成: {len(zones)} 个区域, "
                       f"{sum(len(z['rooms']) for z in zones)} 个配电室, "
                       f"{tree['total_devices']} 台设备")
        return tree

    # ── 名称处理工具 ──────────────────────────────────

    _SUB_KEY_RE = re.compile(r'(\d+[AB]\d+)', re.I)

    @classmethod
    def extract_substation_key(cls, name: str) -> Optional[str]:
        """从配电室名称中提取变电所编码"""
        n = (name or '').strip()
        if not n:
            return None
        if n.upper().startswith('ECC'):
            return 'ECC'
        head = re.split(r'[-_及]', n)[0]
        head = head.replace('配电室', '').replace('模块', '').strip()
        m = cls._SUB_KEY_RE.search(head)
        if m:
            return m.group(1).upper()
        compact = re.sub(r'[^0-9A-Za-z]', '', n)
        m = cls._SUB_KEY_RE.search(compact)
        return m.group(1).upper() if m else None

    @classmethod
    def display_substation_name(cls, name: str) -> str:
        """界面展示名统一为「编码变电所」"""
        key = cls.extract_substation_key(name)
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

    # ── 统计信息 ──────────────────────────────────────

    def get_stats(self) -> Dict[str, Any]:
        """获取层级统计信息"""
        tree = self.build_tree()
        return {
            'zone_count': len(tree['zones']),
            'room_count': sum(len(z['rooms']) for z in tree['zones']),
            'cabinet_count': sum(
                len(r['cabinets']) for z in tree['zones'] for r in z['rooms']
            ),
            'device_count': tree['total_devices'],
            'online_count': sum(
                sum(c['online'] for c in r['cabinets'])
                for z in tree['zones'] for r in z['rooms']
            ),
        }
