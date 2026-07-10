"""
ISM Automation — 组态大屏生成器
从层级树和布局模板生成 ISM 组态大屏页面数据。
"""
from __future__ import annotations

import uuid as _uuid
from typing import Any, Dict, List, Optional

from ism_automation.core.logger import ism_logger
from ism_automation.config.loader import ProjectConfig
from ism_automation.extractors.layout_engine import LayoutEngine, CanvasConfig, LayoutConfig
from ism_automation.extractors.hierarchy_builder import HierarchyBuilder


class DashboardGenerator:
    """组态大屏生成器：层级树 + 模板 → ISM 页面配置"""

    def __init__(self, config: ProjectConfig):
        self.config = config
        self.canvas = CanvasConfig(
            width=config.dashboard.canvas_width,
            height=config.dashboard.canvas_height,
            theme=config.dashboard.theme,
        )
        self.layout = LayoutConfig(
            left_sidebar={
                "x": 0, "y": 0,
                "w": config.dashboard.left_sidebar_width,
                "h": config.dashboard.canvas_height,
            },
            header={
                "x": config.dashboard.left_sidebar_width + 10,
                "y": 0,
                "w": config.dashboard.canvas_width - config.dashboard.left_sidebar_width - 10,
                "h": config.dashboard.header_height,
            },
        )
        self.engine = LayoutEngine(canvas=self.canvas, layout=self.layout)

    @staticmethod
    def _generate_uuid(seed: str) -> str:
        """基于 seed 生成确定性 UUID"""
        return _uuid.uuid5(_uuid.NAMESPACE_DNS, seed).hex

    # ── 页面生成 ──────────────────────────────────────

    def generate_pages(
        self,
        tree: Dict[str, Any],
    ) -> List[Dict[str, Any]]:
        """
        生成完整的多页面组态大屏

        返回格式:
        [
            {
                "name": "overview",
                "uuid": "...",
                "width": 1920, "height": 1080,
                "components": [...],
                "is_home": True,
            },
            {
                "name": "building-1A1_U11",
                "uuid": "...",
                "width": 1920, "height": 1080,
                "components": [...],
            },
            ...
        ]
        """
        pages = []
        stats = self._calc_stats(tree)

        # 1. 概览页面 (overview)
        overview = self._generate_overview_page(tree, stats)
        pages.append(overview)

        # 2. 柜级页面 (building)
        for zone in tree.get("zones", []):
            for room in zone.get("rooms", []):
                for cabinet in room.get("cabinets", []):
                    page = self._generate_building_page(cabinet, zone, room)
                    pages.append(page)

        # 3. 设备组页面 (floor)
        for zone in tree.get("zones", []):
            for room in zone.get("rooms", []):
                for cabinet in room.get("cabinets", []):
                    for group in cabinet.get("floor_groups", []):
                        page = self._generate_floor_page(group, cabinet, zone, room)
                        pages.append(page)

        # 4. 设备详情页面 (device-detail)
        # 通常只创建一个通用模板页面
        device_page = self._generate_device_detail_page()
        pages.append(device_page)

        ism_logger.info(f"🎨 生成 {len(pages)} 个页面")
        return pages

    def _calc_stats(self, tree: Dict[str, Any]) -> Dict[str, Any]:
        """计算层级统计信息"""
        zones = tree.get("zones", [])
        return {
            "zone_count": len(zones),
            "room_count": sum(len(z.get("rooms", [])) for z in zones),
            "cabinet_count": sum(
                len(r.get("cabinets", []))
                for z in zones for r in z.get("rooms", [])
            ),
            "device_count": sum(
                c.get("device_count", 0)
                for z in zones for r in z.get("rooms", []) for c in r.get("cabinets", [])
            ),
            "online_count": sum(
                c.get("online", 0)
                for z in zones for r in z.get("rooms", []) for c in r.get("cabinets", [])
            ),
        }

    def _generate_overview_page(self, tree: Dict[str, Any], stats: Dict[str, Any]) -> Dict[str, Any]:
        """生成概览页面"""
        components = self.engine.generate_overview_layout(stats)

        return {
            "name": "overview",
            "uuid": self._generate_uuid(f"{self.config.project_uuid}-overview"),
            "width": self.canvas.width,
            "height": self.canvas.height,
            "components": components,
            "is_home": True,
        }

    def _generate_building_page(
        self,
        cabinet: Dict[str, Any],
        zone: Dict[str, Any],
        room: Dict[str, Any],
    ) -> Dict[str, Any]:
        """生成柜级页面"""
        cabinet_name = cabinet.get("name", "柜")
        components = self.engine.generate_building_layout(cabinet)

        return {
            "name": f"building-{cabinet_name}",
            "uuid": self._generate_uuid(f"{self.config.project_uuid}-building-{cabinet_name}"),
            "width": self.canvas.width,
            "height": self.canvas.height,
            "components": components,
            "is_home": False,
        }

    def _generate_floor_page(
        self,
        group: Dict[str, Any],
        cabinet: Dict[str, Any],
        zone: Dict[str, Any],
        room: Dict[str, Any],
    ) -> Dict[str, Any]:
        """生成设备组页面"""
        group_name = group.get("name", "设备组")
        components = self.engine.generate_floor_layout(group)

        return {
            "name": f"floor-{group_name}",
            "uuid": self._generate_uuid(f"{self.config.project_uuid}-floor-{group_name}"),
            "width": self.canvas.width,
            "height": self.canvas.height,
            "components": components,
            "is_home": False,
        }

    def _generate_device_detail_page(self) -> Dict[str, Any]:
        """生成设备详情页面"""
        # 使用占位设备信息
        placeholder_device = {"name": "设备详情"}
        components = self.engine.generate_device_detail_layout(placeholder_device)

        return {
            "name": "device-detail",
            "uuid": self._generate_uuid(f"{self.config.project_uuid}-device-detail"),
            "width": self.canvas.width,
            "height": self.canvas.height,
            "components": components,
            "is_home": False,
        }

    # ── 模型创建 ──────────────────────────────────────

    def create_display_model(
        self,
        pages: List[Dict[str, Any]],
        name: Optional[str] = None,
    ) -> Dict[str, Any]:
        """
        创建完整的显示模型数据

        返回格式:
        {
            "name": "航信机房大屏",
            "project_uuid": "...",
            "width": 1920, "height": 1080,
            "pages": [...],
        }
        """
        return {
            "name": name or f"{self.config.project_name}大屏",
            "project_uuid": self.config.project_uuid,
            "width": self.canvas.width,
            "height": self.canvas.height,
            "pages": pages,
        }
