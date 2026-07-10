"""
ISM Automation — 数据模型生成器
将 Excel 解析的数据模型转换为 ISM API 可用的 JSON 格式。
"""
from __future__ import annotations

from typing import Any, Dict, List

from ism_automation.core.logger import ism_logger
from ism_automation.extractors.excel_parser import ExcelParser
from ism_automation.config.loader import ProjectConfig


class ModelGenerator:
    """数据模型生成器：Excel → ISM API 参数"""

    def __init__(self, config: ProjectConfig):
        self.config = config

    def generate(self, parsed_data: Dict[str, Any]) -> List[Dict[str, Any]]:
        """
        从解析后的 Excel 数据生成 ISM 数据模型配置

        返回格式:
        [
            {
                "protocol": "modbus",
                "name": "A20电力仪表",
                "project_uuid": "...",
                "register_groups": [
                    {
                        "name": "AI寄存器组",
                        "start_address": 0,
                        "end_address": 100,
                        "register_type": 3,
                        "registers": [
                            {"name": "A相电压", "address": 0, "data_type": 1, "unit": "V"},
                            ...
                        ]
                    },
                    ...
                ]
            },
            ...
        ]
        """
        template = parsed_data.get("template", {})
        models = []

        for model_key, model_data in template.items():
            register_groups = []

            # AI 寄存器组
            ai_points = model_data.get("ai_points", [])
            if ai_points:
                ai_addresses = [p["offset"] for p in ai_points if p["offset"] is not None]
                if ai_addresses:
                    register_groups.append({
                        "name": "AI寄存器组",
                        "start_address": min(ai_addresses),
                        "end_address": max(ai_addresses) + 10,
                        "register_type": 3,  # 保持寄存器
                        "registers": [
                            {
                                "name": p["name"],
                                "address": p["offset"],
                                "data_type": 1,  # INT16
                                "unit": ExcelParser.get_unit_for_point(p["name"]),
                                "coefficient": p.get("coeff", 1.0),
                                "parse_type": p.get("parse", 177),
                            }
                            for p in ai_points
                            if p["offset"] is not None
                        ]
                    })

            # DI 寄存器组
            di_points = model_data.get("di_points", [])
            if di_points:
                di_addresses = [p["offset"] for p in di_points if p["offset"] is not None]
                if di_addresses:
                    register_groups.append({
                        "name": "DI寄存器组",
                        "start_address": min(di_addresses),
                        "end_address": max(di_addresses) + 10,
                        "register_type": 1,  # 线圈状态
                        "registers": [
                            {
                                "name": p["name"],
                                "address": p["offset"],
                                "data_type": 5,  # BOOL
                                "unit": "",
                            }
                            for p in di_points
                            if p["offset"] is not None
                        ]
                    })

            models.append({
                "protocol": "modbus",
                "name": model_key,
                "project_uuid": self.config.project_uuid,
                "register_groups": register_groups,
            })

        ism_logger.info(f"🔧 生成 {len(models)} 个数据模型配置")
        return models

    def generate_from_excel(self, excel_path: str) -> List[Dict[str, Any]]:
        """从 Excel 文件直接生成数据模型配置"""
        parser = ExcelParser(self.config)
        parsed_data = parser.parse_all()
        return self.generate(parsed_data)
