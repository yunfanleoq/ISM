"""
ISM Automation — 项目管理 API
封装项目的创建、查询、导入、导出等操作。
"""
from __future__ import annotations

from typing import Any, Dict, List, Optional

from ism_automation.core.logger import ism_logger
from ism_automation.api.client import ISMClient


class ProjectAPI:
    """项目管理 API 封装"""

    def __init__(self, client: ISMClient):
        self.client = client

    # ── 项目 CRUD ─────────────────────────────────────

    def create(
        self,
        name: str,
        industry: str = "",
        description: str = "",
        **kwargs,
    ) -> Dict[str, Any]:
        """创建新项目"""
        data = {
            "Name": name,
            "Industry": industry,
            "Description": description,
            **kwargs,
        }
        ism_logger.info(f"🏗️ 创建项目: {name}")
        return self.client.post("/ProjectAdd", json_data=data)

    def list(self) -> List[Dict]:
        """查询项目列表"""
        resp = self.client.get("/ProjectList")
        return resp.get("data", {}).get("list", [])

    def edit(self, uuid: str, **kwargs) -> Dict[str, Any]:
        """编辑项目"""
        data = {"uuid": uuid, **kwargs}
        return self.client.post("/ProjectEdit", json_data=data)

    def delete(self, uuid: str) -> Dict[str, Any]:
        """删除项目"""
        ism_logger.warning(f"🗑️ 删除项目: {uuid}")
        return self.client.post("/ProjectDel", json_data={"uuid": uuid})

    # ── 导入/导出 ─────────────────────────────────────

    def export_project(self, uuid: str) -> Dict[str, Any]:
        """导出项目数据包"""
        return self.client.get("/ExportProject", params={"uuid": uuid})

    def import_project(self, file_path: str) -> Dict[str, Any]:
        """导入项目数据包"""
        with open(file_path, "rb") as f:
            resp = self.client.post(
                "/ImportProject",
                files={"file": (file_path, f, "application/json")},
            )
        return resp

    # ── 辅助 ──────────────────────────────────────────

    def get_by_name(self, name: str) -> Optional[Dict]:
        """根据名称查找项目"""
        projects = self.list()
        for p in projects:
            if p.get("Name") == name or p.get("name") == name:
                return p
        return None

    def get_or_create(self, name: str, **kwargs) -> Dict[str, Any]:
        """获取或创建项目"""
        existing = self.get_by_name(name)
        if existing:
            ism_logger.info(f"📁 项目已存在: {name}")
            return existing
        return self.create(name, **kwargs)
