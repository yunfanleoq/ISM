"""
ISM Automation — 布局引擎
负责计算组态大屏组件的坐标、尺寸、颜色等视觉属性。

支持模板驱动的布局计算，可复用到不同项目。
"""
from __future__ import annotations

from typing import Any, Dict, List, Optional, Tuple
from dataclasses import dataclass, field

from ism_automation.core.logger import ism_logger


# ── 布局常量 ──────────────────────────────────────

@dataclass
class CanvasConfig:
    """画布配置"""
    width: int = 1920
    height: int = 1080
    bg_color: str = "#0a0e17"
    theme: str = "dark"


@dataclass
class LayoutConfig:
    """布局区域配置"""
    # 左侧导航栏
    left_sidebar: Dict[str, Any] = field(default_factory=lambda: {
        "x": 0, "y": 0, "w": 230, "h": 1080
    })
    # 顶部标题栏
    header: Dict[str, Any] = field(default_factory=lambda: {
        "x": 240, "y": 0, "w": 1680, "h": 80
    })
    # 统计卡片区域
    stats_cards: Dict[str, Any] = field(default_factory=lambda: {
        "x": 290, "y": 100, "w": 1560, "h": 110, "count": 4, "gap": 10
    })
    # 左侧主面板
    left_panel: Dict[str, Any] = field(default_factory=lambda: {
        "x": 290, "y": 230, "w": 780, "h": 400
    })
    # 右上面板
    right_upper: Dict[str, Any] = field(default_factory=lambda: {
        "x": 1090, "y": 230, "w": 400, "h": 400
    })
    # 右下面板
    right_lower: Dict[str, Any] = field(default_factory=lambda: {
        "x": 1090, "y": 650, "w": 400, "h": 220
    })
    # 底部面板
    bottom_panel: Dict[str, Any] = field(default_factory=lambda: {
        "x": 290, "y": 650, "w": 780, "h": 220
    })


# ── 主题配色 ──────────────────────────────────────

THEMES = {
    "dark": {
        "bg": "#0a0e17",
        "card_bg": "#1a2332",
        "card_border": "#2a3a4a",
        "text_primary": "#e0e6ed",
        "text_secondary": "#8a94a6",
        "accent": "#00d4aa",
        "accent_warning": "#ff9f43",
        "accent_danger": "#ff6b6b",
        "chart_colors": ["#00d4aa", "#3498db", "#9b59b6", "#e74c3c", "#f1c40f"],
    },
    "light": {
        "bg": "#f5f7fa",
        "card_bg": "#ffffff",
        "card_border": "#e0e6ed",
        "text_primary": "#2c3e50",
        "text_secondary": "#7f8c8d",
        "accent": "#3498db",
        "accent_warning": "#e67e22",
        "accent_danger": "#e74c3c",
        "chart_colors": ["#3498db", "#2ecc71", "#9b59b6", "#e74c3c", "#f1c40f"],
    },
}


class LayoutEngine:
    """
    组态大屏布局引擎

    支持基于模板的布局计算，输出标准 ISM 组件配置。
    """

    def __init__(self, canvas: Optional[CanvasConfig] = None, layout: Optional[LayoutConfig] = None):
        self.canvas = canvas or CanvasConfig()
        self.layout = layout or LayoutConfig()
        self.theme = THEMES.get(self.canvas.theme, THEMES["dark"])

    # ── 基础组件工厂 ──────────────────────────────────

    def create_text(
        self,
        x: int, y: int, w: int, h: int,
        text: str = "",
        font_size: int = 14,
        color: Optional[str] = None,
        align: str = "center",
        bold: bool = False,
        **kwargs,
    ) -> Dict[str, Any]:
        """创建文本组件"""
        return {
            "type": "ism-view-text",
            "style": {
                "position": {"x": x, "y": y, "w": w, "h": h},
                "fontSize": font_size,
                "foreColor": color or self.theme["text_primary"],
                "backColor": "transparent",
                "textAlign": align,
                "fontWeight": "bold" if bold else "normal",
                "borderWidth": 0,
                "opacity": 1,
                **kwargs.get("style", {}),
            },
            "text": text,
            **{k: v for k, v in kwargs.items() if k != "style"},
        }

    def create_card(
        self,
        x: int, y: int, w: int, h: int,
        title: str = "",
        value: str = "",
        unit: str = "",
        icon: str = "",
        color: Optional[str] = None,
        **kwargs,
    ) -> Dict[str, Any]:
        """创建统计卡片组件"""
        card_color = color or self.theme["card_bg"]
        return {
            "type": "dv-border-box-8",
            "style": {
                "position": {"x": x, "y": y, "w": w, "h": h},
                "backColor": card_color,
                "borderColor": self.theme["card_border"],
                "borderWidth": 1,
                "BorderEdges": 4,
                "opacity": 1,
                **kwargs.get("style", {}),
            },
            "title": title,
            "children": [
                # 标题文本
                self.create_text(
                    x=x + 10, y=y + 10, w=w - 20, h=30,
                    text=title, font_size=12, color=self.theme["text_secondary"], align="left"
                ),
                # 数值文本
                self.create_text(
                    x=x + 10, y=y + 45, w=w - 20, h=40,
                    text=f"{value} {unit}", font_size=24, color=self.theme["accent"], bold=True
                ),
            ],
            **{k: v for k, v in kwargs.items() if k not in ("style", "children")},
        }

    def create_chart(
        self,
        x: int, y: int, w: int, h: int,
        chart_type: str = "line",
        title: str = "",
        data_source: str = "",
        **kwargs,
    ) -> Dict[str, Any]:
        """创建图表组件"""
        return {
            "type": f"ism-chart-{chart_type}",
            "style": {
                "position": {"x": x, "y": y, "w": w, "h": h},
                "backColor": self.theme["card_bg"],
                "borderColor": self.theme["card_border"],
                "borderWidth": 1,
                "BorderEdges": 4,
                **kwargs.get("style", {}),
            },
            "title": title,
            "dataBind": {
                "dataSource": data_source,
                **kwargs.get("dataBind", {}),
            },
            "chartConfig": {
                "colors": self.theme["chart_colors"],
                **kwargs.get("chartConfig", {}),
            },
            **{k: v for k, v in kwargs.items() if k not in ("style", "dataBind", "chartConfig")},
        }

    def create_table(
        self,
        x: int, y: int, w: int, h: int,
        title: str = "",
        columns: Optional[List[str]] = None,
        **kwargs,
    ) -> Dict[str, Any]:
        """创建表格组件"""
        return {
            "type": "ism-view-table",
            "style": {
                "position": {"x": x, "y": y, "w": w, "h": h},
                "backColor": self.theme["card_bg"],
                "borderColor": self.theme["card_border"],
                "borderWidth": 1,
                **kwargs.get("style", {}),
            },
            "title": title,
            "columns": columns or [],
            **{k: v for k, v in kwargs.items() if k not in ("style", "columns")},
        }

    def create_border_box(
        self,
        x: int, y: int, w: int, h: int,
        title: str = "",
        **kwargs,
    ) -> Dict[str, Any]:
        """创建 DataV 边框容器组件"""
        return {
            "type": "dv-border-box-8",
            "style": {
                "position": {"x": x, "y": y, "w": w, "h": h},
                "backColor": self.theme["card_bg"],
                "borderColor": self.theme["accent"],
                "borderWidth": 2,
                **kwargs.get("style", {}),
            },
            "title": title,
            **{k: v for k, v in kwargs.items() if k != "style"},
        }

    # ── 布局计算 ──────────────────────────────────────

    def calc_stats_card_positions(self) -> List[Tuple[int, int, int, int]]:
        """计算统计卡片位置"""
        cfg = self.layout.stats_cards
        count = cfg.get("count", 4)
        gap = cfg.get("gap", 10)
        total_width = cfg["w"]
        card_width = (total_width - gap * (count - 1)) // count
        card_height = cfg["h"]

        positions = []
        for i in range(count):
            x = cfg["x"] + i * (card_width + gap)
            y = cfg["y"]
            positions.append((x, y, card_width, card_height))

        return positions

    def calc_grid_positions(
        self,
        start_x: int, start_y: int,
        cols: int, rows: int,
        cell_w: int, cell_h: int,
        gap_x: int = 10, gap_y: int = 10,
    ) -> List[Tuple[int, int, int, int]]:
        """计算网格布局位置"""
        positions = []
        for r in range(rows):
            for c in range(cols):
                x = start_x + c * (cell_w + gap_x)
                y = start_y + r * (cell_h + gap_y)
                positions.append((x, y, cell_w, cell_h))
        return positions

    # ── 页面布局生成 ──────────────────────────────────

    def generate_overview_layout(self, stats: Dict[str, Any]) -> List[Dict[str, Any]]:
        """
        生成概览页面布局

        stats 包含: zone_count, room_count, cabinet_count, device_count, online_count
        """
        components = []
        theme = self.theme

        # 1. 左侧导航栏背景
        sidebar = self.layout.left_sidebar
        components.append(self.create_border_box(
            x=sidebar["x"], y=sidebar["y"], w=sidebar["w"], h=sidebar["h"],
            title="",
            style={"backColor": "#0d1320", "borderWidth": 0}
        ))

        # 2. 顶部标题
        header = self.layout.header
        components.append(self.create_text(
            x=header["x"] + 20, y=header["y"] + 20,
            w=header["w"] - 40, h=40,
            text="电力监控系统", font_size=28, color=theme["text_primary"], bold=True, align="left"
        ))

        # 3. 统计卡片
        stats_data = [
            {"title": "总设备数", "value": stats.get("device_count", 0), "unit": "台", "color": theme["accent"]},
            {"title": "在线设备", "value": stats.get("online_count", 0), "unit": "台", "color": theme["accent"]},
            {"title": "告警设备", "value": stats.get("device_count", 0) - stats.get("online_count", 0), "unit": "台", "color": theme["accent_danger"]},
            {"title": "配电室数", "value": stats.get("room_count", 0), "unit": "个", "color": theme["accent_warning"]},
        ]
        positions = self.calc_stats_card_positions()
        for (x, y, w, h), stat in zip(positions, stats_data):
            components.append(self.create_card(
                x=x, y=y, w=w, h=h,
                title=stat["title"], value=str(stat["value"]), unit=stat["unit"],
                color=stat["color"],
            ))

        # 4. 左侧主图表（功率趋势）
        lp = self.layout.left_panel
        components.append(self.create_chart(
            x=lp["x"], y=lp["y"], w=lp["w"], h=lp["h"],
            chart_type="line", title="功率趋势图",
            data_source="power_trend",
        ))

        # 5. 右上设备状态饼图
        ru = self.layout.right_upper
        components.append(self.create_chart(
            x=ru["x"], y=ru["y"], w=ru["w"], h=ru["h"],
            chart_type="pie", title="设备状态分布",
            data_source="device_status",
        ))

        # 6. 右下告警列表
        rl = self.layout.right_lower
        components.append(self.create_table(
            x=rl["x"], y=rl["y"], w=rl["w"], h=rl["h"],
            title="实时告警", columns=["时间", "设备", "告警内容", "等级"],
        ))

        # 7. 底部设备列表
        bp = self.layout.bottom_panel
        components.append(self.create_table(
            x=bp["x"], y=bp["y"], w=bp["w"], h=bp["h"],
            title="设备运行状态", columns=["设备名称", "状态", "A相电压", "A相电流", "功率"],
        ))

        ism_logger.info(f"🎨 生成概览布局: {len(components)} 个组件")
        return components

    def generate_building_layout(self, cabinet: Dict[str, Any]) -> List[Dict[str, Any]]:
        """生成柜级页面布局（设备组卡片网格）"""
        components = []
        floor_groups = cabinet.get("floor_groups", [])

        # 标题
        components.append(self.create_text(
            x=290, y=20, w=800, h=50,
            text=f"{cabinet.get('name', '柜')} — 设备组概览",
            font_size=24, bold=True, align="left"
        ))

        # 设备组卡片网格
        cols = 4
        rows = (len(floor_groups) + cols - 1) // cols
        positions = self.calc_grid_positions(
            start_x=290, start_y=100,
            cols=cols, rows=rows,
            cell_w=350, cell_h=180,
            gap_x=20, gap_y=20,
        )

        for (x, y, w, h), group in zip(positions, floor_groups):
            components.append(self.create_card(
                x=x, y=y, w=w, h=h,
                title=group["name"],
                value=str(group["count"]),
                unit="台",
                color=self.theme["accent"],
            ))

        return components

    def generate_floor_layout(self, group: Dict[str, Any]) -> List[Dict[str, Any]]:
        """生成设备组级页面布局（设备列表表格）"""
        components = []
        devices = group.get("devices", [])

        # 标题
        components.append(self.create_text(
            x=290, y=20, w=800, h=50,
            text=f"{group.get('name', '设备组')} — 设备列表",
            font_size=24, bold=True, align="left"
        ))

        # 设备表格
        components.append(self.create_table(
            x=290, y=100, w=1400, h=800,
            title="设备实时数据",
            columns=["设备名称", "状态", "A相电压", "B相电压", "C相电压", "A相电流", "功率"],
        ))

        return components

    def generate_device_detail_layout(self, device: Dict[str, Any]) -> List[Dict[str, Any]]:
        """生成设备详情页面布局"""
        components = []

        # 标题
        components.append(self.create_text(
            x=290, y=20, w=800, h=50,
            text=f"{device.get('name', '设备')} — 参数详情",
            font_size=24, bold=True, align="left"
        ))

        # 参数卡片网格（假设 91 个参数）
        positions = self.calc_grid_positions(
            start_x=290, start_y=100,
            cols=7, rows=13,
            cell_w=180, cell_h=60,
            gap_x=10, gap_y=10,
        )

        for i, (x, y, w, h) in enumerate(positions):
            components.append(self.create_card(
                x=x, y=y, w=w, h=h,
                title=f"参数 {i+1}", value="--", unit="",
                color=self.theme["card_bg"],
            ))

        return components
