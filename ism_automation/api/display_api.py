"""
ISM Automation — 组态大屏 API
封装显示模型、页面、图层、组件的创建与管理。
"""
from __future__ import annotations

import json
from typing import Any, Dict, List, Optional

from ism_automation.core.logger import ism_logger
from ism_automation.api.client import ISMClient


class DisplayAPI:
    """组态大屏 API 封装"""

    def __init__(self, client: ISMClient):
        self.client = client

    # ── 显示模型 ──────────────────────────────────────

    def create(
        self,
        name: str,
        project_uuid: str,
        description: str = "",
        width: int = 1920,
        height: int = 1080,
        **kwargs,
    ) -> Dict[str, Any]:
        """创建组态大屏模型"""
        data = {
            "Name": name,
            "ProjectUuid": project_uuid,
            "Description": description,
            "Width": width,
            "Height": height,
            **kwargs,
        }
        ism_logger.info(f"🎨 创建组态模型: {name}")
        return self.client.post("/displayModelAdd", json_data=data)

    def list(self, project_uuid: Optional[str] = None) -> List[Dict]:
        """查询组态模型列表"""
        params = {}
        if project_uuid:
            params["ProjectUuid"] = project_uuid
        resp = self.client.get("/displayModelList", params=params)
        return resp.get("data", {}).get("list", [])

    def get(self, uuid: str) -> Dict[str, Any]:
        """获取单个组态模型详情"""
        return self.client.get("/displayModelGet", params={"uuid": uuid})

    def edit(self, uuid: str, **kwargs) -> Dict[str, Any]:
        """编辑组态模型"""
        data = {"uuid": uuid, **kwargs}
        return self.client.post("/displayModelEdit", json_data=data)

    def delete(self, uuid: str) -> Dict[str, Any]:
        """删除组态模型（移到回收站）"""
        ism_logger.warning(f"🗑️ 删除组态模型: {uuid}")
        return self.client.post("/displayModelDel", json_data={"uuid": uuid})

    def force_delete(self, uuid: str) -> Dict[str, Any]:
        """彻底删除组态模型"""
        return self.client.post("/displayModelForceDel", json_data={"uuid": uuid})

    # ── 页面管理 ──────────────────────────────────────

    def add_page(
        self,
        model_uuid: str,
        name: str,
        width: int = 1920,
        height: int = 1080,
        is_home: bool = False,
        **kwargs,
    ) -> Dict[str, Any]:
        """为大屏模型添加子页面"""
        data = {
            "ModelUuid": model_uuid,
            "Name": name,
            "Width": width,
            "Height": height,
            "IsHome": is_home,
            **kwargs,
        }
        ism_logger.info(f"📄 添加页面: {name} → 模型 {model_uuid}")
        return self.client.post("/DisplayModelPageAdd", json_data=data)

    def delete_page(self, page_uuid: str) -> Dict[str, Any]:
        """删除子页面"""
        return self.client.post("/DisplayModelPageDel", json_data={"uuid": page_uuid})

    def edit_page(self, page_uuid: str, **kwargs) -> Dict[str, Any]:
        """编辑子页面"""
        data = {"uuid": page_uuid, **kwargs}
        return self.client.post("/DisplayModelPageEdit", json_data=data)

    def set_home_page(self, model_uuid: str, page_uuid: str) -> Dict[str, Any]:
        """设置首页"""
        return self.client.post("/DisplayModelPageSetHome", json_data={
            "ModelUuid": model_uuid,
            "PageUuid": page_uuid,
        })

    # ── 图层数据 ──────────────────────────────────────

    def get_layer_data(self, model_uuid: str, page_uuid: Optional[str] = None) -> Dict[str, Any]:
        """获取页面图层数据（组件列表）"""
        params = {"uuid": model_uuid}
        if page_uuid:
            params["pageUuid"] = page_uuid
        return self.client.get("/getDisplayModelLayerData", params=params)

    def save_layer_data(
        self,
        model_uuid: str,
        page_uuid: str,
        components: List[Dict[str, Any]],
    ) -> Dict[str, Any]:
        """
        保存页面图层数据（组件列表）

        components 格式为 ISM 组件 JSON 数组，每个组件包含:
        - type: 组件类型 (ism-view-text, ism-chart-line, etc.)
        - style: 样式属性 (position, colors, etc.)
        - dataBind: 数据绑定
        - animate: 动画配置
        - active: 交互动作
        """
        # ISM 后端可能要求特定的数据格式
        data = {
            "uuid": model_uuid,
            "pageUuid": page_uuid,
            "components": components,
        }
        ism_logger.info(f"💾 保存图层数据: {len(components)} 个组件 → 页面 {page_uuid}")
        return self.client.post("/saveDisplayModelLayerData", json_data=data)

    def save_layer_data_raw(
        self,
        model_uuid: str,
        page_uuid: str,
        layer_data: Dict[str, Any],
    ) -> Dict[str, Any]:
        """直接保存完整的图层数据对象（包含 metadata）"""
        data = {
            "uuid": model_uuid,
            "pageUuid": page_uuid,
            **layer_data,
        }
        return self.client.post("/saveDisplayModelLayerData", json_data=data)

    # ── 模板相关 ──────────────────────────────────────

    def list_templates(self) -> List[Dict]:
        """获取系统模板列表"""
        resp = self.client.get("/DisplayTempleteList")
        return resp.get("data", {}).get("list", [])

    def get_template(self, uuid: str) -> Dict[str, Any]:
        """获取模板详情"""
        return self.client.get("/DisplayTempleteGet", params={"uuid": uuid})

    # ── 批量页面操作 ──────────────────────────────────

    def create_multi_page_dashboard(
        self,
        model_uuid: str,
        pages: List[Dict[str, Any]],
        set_home_index: int = 0,
    ) -> List[Dict[str, Any]]:
        """
        批量创建多页面组态大屏

        pages 格式:
        [
            {
                "name": "overview",
                "width": 1920,
                "height": 1080,
                "components": [...],  # 可选，同时保存图层数据
            },
            ...
        ]
        """
        results = []
        home_page_uuid = None

        for i, page in enumerate(pages):
            # 创建页面
            page_resp = self.add_page(
                model_uuid=model_uuid,
                name=page["name"],
                width=page.get("width", 1920),
                height=page.get("height", 1080),
            )
            page_data = page_resp.get("data", {})
            page_uuid = page_data.get("uuid") or page_data.get("Uuid")
            results.append(page_resp)

            if i == set_home_index:
                home_page_uuid = page_uuid

            # 如果有组件数据，保存图层
            components = page.get("components")
            if components and page_uuid:
                self.save_layer_data(model_uuid, page_uuid, components)

        # 设置首页
        if home_page_uuid:
            self.set_home_page(model_uuid, home_page_uuid)

        return results

    # ── 用户显示列表 ──────────────────────────────────

    def get_user_display_list(self, project_uuid: Optional[str] = None) -> List[Dict]:
        """获取用户显示列表（应用管理列表）"""
        params = {}
        if project_uuid:
            params["ProjectUuid"] = project_uuid
        resp = self.client.get("/GetUserDisplayList", params=params)
        return resp.get("data", {}).get("list", [])

    # ── 辅助：构建标准组件 ────────────────────────────

    @staticmethod
    def build_text_component(
        x: int, y: int, w: int, h: int,
        text: str = "",
        font_size: int = 14,
        fore_color: str = "#ffffff",
        back_color: str = "transparent",
        **kwargs,
    ) -> Dict[str, Any]:
        """构建文本组件"""
        return {
            "type": "ism-view-text",
            "style": {
                "position": {"x": x, "y": y, "w": w, "h": h},
                "fontSize": font_size,
                "foreColor": fore_color,
                "backColor": back_color,
                "textAlign": "center",
                "borderWidth": 0,
                "opacity": 1,
                **kwargs.get("style", {}),
            },
            "text": text,
            **{k: v for k, v in kwargs.items() if k != "style"},
        }

    @staticmethod
    def build_chart_component(
        x: int, y: int, w: int, h: int,
        chart_type: str = "line",
        data_source: str = "",
        **kwargs,
    ) -> Dict[str, Any]:
        """构建图表组件"""
        return {
            "type": f"ism-chart-{chart_type}",
            "style": {
                "position": {"x": x, "y": y, "w": w, "h": h},
                **kwargs.get("style", {}),
            },
            "dataBind": {
                "dataSource": data_source,
                **kwargs.get("dataBind", {}),
            },
            **{k: v for k, v in kwargs.items() if k not in ("style", "dataBind")},
        }

    @staticmethod
    def build_border_box_component(
        x: int, y: int, w: int, h: int,
        title: str = "",
        **kwargs,
    ) -> Dict[str, Any]:
        """构建边框容器组件（DataV 风格）"""
        return {
            "type": "dv-border-box-8",
            "style": {
                "position": {"x": x, "y": y, "w": w, "h": h},
                **kwargs.get("style", {}),
            },
            "title": title,
            **{k: v for k, v in kwargs.items() if k != "style"},
        }
