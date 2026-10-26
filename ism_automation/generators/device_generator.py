"""
ISM Automation — 设备生成器
从 Excel 设备清单生成 ISM 设备树配置。
"""
from __future__ import annotations

import re
from typing import Any, Dict, List, Optional

from ism_automation.core.logger import ism_logger
from ism_automation.config.loader import ProjectConfig


class DeviceGenerator:
    """设备生成器：Excel 设备清单 → ISM 设备树配置"""

    def __init__(self, config: ProjectConfig):
        self.config = config

    def generate(
        self,
        device_list: List[Dict[str, Any]],
        model_mapping: Dict[str, str],  # 设备类型 → 模型UUID
    ) -> List[Dict[str, Any]]:
        """
        生成设备树配置

        device_list 格式:
        [
            {"full_name": "1A1_U11_S18_1", "short_name": "1A1_S18_1", "ai_start": 2000, "di_start": 3000},
            ...
        ]

        model_mapping 格式:
        {"A20电力仪表": "uuid-1", "A40电力仪表": "uuid-2", "施耐德UPS": "uuid-3"}

        返回格式:
        [
            {
                "name": "1A1_U11_S18_1",
                "project_uuid": "...",
                "parent_sid": 0,
                "device_type": 1,
                "model_uuid": "uuid-1",
                "protocol_type": 2,
                "ip": "172.31.4.14",
                "port": 502,
                "slave_id": 1,
                "status": 1,
            },
            ...
        ]
        """
        devices = []

        for dev in device_list:
            full_name = dev.get("full_name", "")
            device_type = self._determine_device_type(full_name)
            model_uuid = model_mapping.get(device_type)

            if not model_uuid:
                ism_logger.warning(f"⚠️ 设备 {full_name} 未找到对应模型: {device_type}")
                continue

            # 解析 IP 和从机地址（从配置或名称推断）
            ip = self._extract_ip(dev)
            slave_id = self._extract_slave_id(dev)

            devices.append({
                "name": full_name,
                "project_uuid": self.config.project_uuid,
                "parent_sid": 0,  # 根节点，后续由层级构建器调整
                "device_type": 1,  # 设备
                "model_uuid": model_uuid,
                "protocol_type": 2,  # Modbus TCP
                "ip": ip,
                "port": 502,
                "slave_id": slave_id,
                "status": 1,
                "ai_start": dev.get("ai_start"),
                "di_start": dev.get("di_start"),
            })

        ism_logger.info(f"🔧 生成 {len(devices)} 台设备配置")
        return devices

    def _determine_device_type(self, full_name: str) -> str:
        """判断设备类型"""
        if 'A40' in full_name.upper():
            return 'A40电力仪表'
        if 'UPS' in full_name.upper() or '施耐德' in full_name:
            return '施耐德UPS'
        return 'A20电力仪表'

    def _extract_ip(self, dev: Dict[str, Any]) -> str:
        """从设备信息提取 IP（从 Excel 文件名或配置）"""
        # 从 Excel 文件名中提取 IP
        excel_path = self.config.excel_path
        ip_match = re.search(r'(\d+\.\d+\.\d+\.\d+)', excel_path)
        if ip_match:
            return ip_match.group(1)
        return "127.0.0.1"

    def _extract_slave_id(self, dev: Dict[str, Any]) -> int:
        """从设备名称提取从机地址"""
        # 从名称末尾的数字提取
        name = dev.get("full_name", "")
        parts = name.split('_')
        if parts:
            last_part = parts[-1]
            try:
                return int(last_part)
            except ValueError:
                pass
        return 1

    def build_tree(
        self,
        devices: List[Dict[str, Any]],
    ) -> Dict[str, Any]:
        """
        构建设备树层级结构（用于 API 批量创建）

        返回格式:
        {
            "zones": [
                {
                    "name": "1A",
                    "children": [
                        {
                            "name": "1A1配电室",
                            "children": [
                                {
                                    "name": "1A1_U11柜",
                                    "children": [
                                        {"name": "1A1_U11_S18_1", ...},
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
        # 按名称前缀分组
        zones = {}
        for dev in devices:
            name = dev["name"]
            parts = name.split('_')
            if len(parts) < 3:
                continue

            zone_key = parts[0][0] + parts[0][1:]  # e.g., "1A1" → "1A"
            room_key = parts[0]  # e.g., "1A1"
            cabinet_key = parts[1]  # e.g., "U11"

            if zone_key not in zones:
                zones[zone_key] = {"name": zone_key, "children": {}}

            rooms = zones[zone_key]["children"]
            if room_key not in rooms:
                rooms[room_key] = {"name": f"{room_key}配电室", "children": {}}

            cabinets = rooms[room_key]["children"]
            if cabinet_key not in cabinets:
                cabinets[cabinet_key] = {"name": f"{room_key}_{cabinet_key}柜", "children": []}

            cabinets[cabinet_key]["children"].append(dev)

        # 转换为列表格式
        return {
            "zones": [
                {
                    "name": z["name"],
                    "children": [
                        {
                            "name": r["name"],
                            "children": [
                                {
                                    "name": c["name"],
                                    "children": c["children"],
                                }
                                for c in r["children"].values()
                            ]
                        }
                        for r in z["children"].values()
                    ]
                }
                for z in zones.values()
            ]
        }
