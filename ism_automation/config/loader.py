"""
ISM Automation — 配置管理系统

支持 YAML 配置文件加载、验证、多项目隔离。
"""
from __future__ import annotations

import os
import yaml
from pathlib import Path
from typing import Dict, Any, Optional
from dataclasses import dataclass, field, asdict
from pydantic import BaseModel, Field, validator

from ism_automation.core.logger import ism_logger


# ── 配置路径 ──────────────────────────────────────
CONFIG_ROOT = Path(__file__).parent
PROJECTS_DIR = CONFIG_ROOT / "projects"
TEMPLATES_DIR = Path(__file__).parent.parent / "templates"


# ── Pydantic 模型 ─────────────────────────────────

class ModelConfig(BaseModel):
    """数据模型配置"""
    name: str
    template_sheet: str
    ai_count: Optional[int] = None
    di_count: Optional[int] = None
    protocol: str = "modbus"


class DashboardConfig(BaseModel):
    """组态大屏配置"""
    theme: str = "dark"
    canvas_width: int = 1920
    canvas_height: int = 1080
    left_sidebar_width: int = 230
    header_height: int = 80
    drill_down_levels: int = 4
    layout_template: str = "industrial_4level"
    show_alarm_panel: bool = True


class ExcelMappingConfig(BaseModel):
    """Excel 列映射配置"""
    device_name_col: str = "A"
    ai_start_col: str = "O"
    di_start_col: str = "P"
    model_name_col: str = "A"
    ai_offset_col: str = "B"
    ai_coeff_col: str = "C"
    ai_parse_col: str = "D"
    di_offset_col: str = "I"
    di_bit_col: str = "J"
    template_sheet_name: str = "模板"


class ProjectConfig(BaseModel):
    """项目配置根模型"""
    project_name: str
    project_uuid: str
    excel_path: str
    models: list[ModelConfig] = Field(default_factory=list)
    dashboard: DashboardConfig = Field(default_factory=DashboardConfig)
    excel_mapping: ExcelMappingConfig = Field(default_factory=ExcelMappingConfig)
    api_base_url: str = "http://127.0.0.1:8081"
    db_host: str = "127.0.0.1"
    db_port: int = 2881
    db_user: str = "root@ism_tenant"
    db_password: str = "ism2024!"
    db_name: str = "ism"

    @validator("excel_path")
    def excel_path_exists(cls, v, values):
        """校验 Excel 路径是否存在（如果传入的是绝对路径）"""
        if os.path.isabs(v):
            if not os.path.exists(v):
                ism_logger.warning(f"Excel 路径不存在: {v}")
        return v


# ── 加载函数 ──────────────────────────────────────

def load_project_config(project_name: str) -> ProjectConfig:
    """加载指定项目的 YAML 配置"""
    config_path = PROJECTS_DIR / f"{project_name}.yaml"
    if not config_path.exists():
        raise FileNotFoundError(
            f"项目配置不存在: {config_path}\n"
            f"可用项目: {list_available_projects()}"
        )

    with open(config_path, "r", encoding="utf-8") as f:
        data = yaml.safe_load(f)

    # 解析 Excel 路径（支持相对项目根目录）
    excel_path = data.get("project", {}).get("excel_path", "")
    if not os.path.isabs(excel_path):
        project_root = Path(__file__).parent.parent.parent
        excel_path = str(project_root / excel_path)
        data["project"]["excel_path"] = excel_path

    # 构建 Pydantic 模型
    project_data = data.get("project", {})

    config = ProjectConfig(
        project_name=project_data.get("name", project_name),
        project_uuid=project_data.get("uuid", ""),
        excel_path=project_data.get("excel_path", ""),
        models=[
            ModelConfig(**m) for m in project_data.get("models", [])
        ],
        dashboard=DashboardConfig(**project_data.get("dashboard", {})),
        excel_mapping=ExcelMappingConfig(**project_data.get("excel_mapping", {})),
        api_base_url=project_data.get("api_base_url", "http://127.0.0.1:8081"),
        db_host=project_data.get("db_host", "127.0.0.1"),
        db_port=project_data.get("db_port", 2881),
        db_user=project_data.get("db_user", "root@ism_tenant"),
        db_password=project_data.get("db_password", "ism2024!"),
        db_name=project_data.get("db_name", "ism"),
    )

    ism_logger.info(f"✅ 加载项目配置: {config.project_name} ({config.project_uuid})")
    return config


def list_available_projects() -> list[str]:
    """列出所有可用的项目配置"""
    if not PROJECTS_DIR.exists():
        return []
    return [
        p.stem for p in PROJECTS_DIR.glob("*.yaml")
    ]


def save_project_config(project_name: str, config: ProjectConfig) -> None:
    """保存项目配置到 YAML 文件"""
    config_path = PROJECTS_DIR / f"{project_name}.yaml"
    PROJECTS_DIR.mkdir(parents=True, exist_ok=True)

    # 转换为字典
    data = {
        "project": {
            "name": config.project_name,
            "uuid": config.project_uuid,
            "excel_path": config.excel_path,
            "models": [m.dict() for m in config.models],
            "dashboard": config.dashboard.dict(),
            "excel_mapping": config.excel_mapping.dict(),
            "api_base_url": config.api_base_url,
            "db_host": config.db_host,
            "db_port": config.db_port,
            "db_user": config.db_user,
            "db_password": config.db_password,
            "db_name": config.db_name,
        }
    }

    with open(config_path, "w", encoding="utf-8") as f:
        yaml.dump(data, f, allow_unicode=True, sort_keys=False, default_flow_style=False)

    ism_logger.info(f"💾 保存项目配置: {config_path}")


def load_template_config(template_name: str) -> dict:
    """加载大屏模板配置"""
    template_path = TEMPLATES_DIR / f"{template_name}.yaml"
    if not template_path.exists():
        raise FileNotFoundError(f"模板不存在: {template_path}")

    with open(template_path, "r", encoding="utf-8") as f:
        return yaml.safe_load(f)
