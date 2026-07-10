"""
ISM Automation — 数据模型 API (Modbus / OPCUA / SNMP / DLT645 / MQTT / S7 / HJ212 / IEC104 / CJT188)
封装设备模型的创建、查询、修改、删除，以及寄存器组和数据点的管理。
"""
from __future__ import annotations

from typing import Any, Dict, List, Optional

from ism_automation.core.logger import ism_logger
from ism_automation.api.client import ISMClient


class ModelAPI:
    """数据模型 API 封装"""

    def __init__(self, client: ISMClient):
        self.client = client

    # ── Modbus 模型 ──────────────────────────────────

    def create_modbus(self, name: str, project_uuid: str, **kwargs) -> Dict[str, Any]:
        """创建 Modbus 设备模型"""
        data = {
            "Name": name,
            "ProjectUuid": project_uuid,
            **kwargs,
        }
        ism_logger.info(f"📦 创建 Modbus 模型: {name}")
        return self.client.post("/modbusModelAdd", json_data=data)

    def list_modbus(self, project_uuid: Optional[str] = None) -> List[Dict]:
        """查询 Modbus 模型列表"""
        params = {}
        if project_uuid:
            params["ProjectUuid"] = project_uuid
        resp = self.client.get("/modbusModelList", params=params)
        return resp.get("data", {}).get("list", [])

    def get_modbus(self, uuid: str) -> Dict[str, Any]:
        """获取单个 Modbus 模型详情"""
        return self.client.get("/modbusModelEdit", params={"uuid": uuid})

    def delete_modbus(self, uuid: str) -> Dict[str, Any]:
        """删除 Modbus 模型"""
        ism_logger.warning(f"🗑️ 删除 Modbus 模型: {uuid}")
        return self.client.post("/modbusModelDel", json_data={"uuid": uuid})

    # ── Modbus 寄存器组 ──────────────────────────────

    def add_register_group(
        self,
        model_uuid: str,
        name: str,
        start_address: int,
        end_address: int,
        register_type: int = 3,  # 3=保持寄存器, 4=输入寄存器
    ) -> Dict[str, Any]:
        """为 Modbus 模型添加寄存器组"""
        data = {
            "modelUuid": model_uuid,
            "Name": name,
            "StartAddress": start_address,
            "EndAddress": end_address,
            "RegisterType": register_type,
        }
        ism_logger.info(f"📦 添加寄存器组: {name} ({start_address}-{end_address})")
        return self.client.post("/modbusModelRegisterGroupAdd", json_data=data)

    def list_register_groups(self, model_uuid: str) -> List[Dict]:
        """查询模型的寄存器组列表"""
        resp = self.client.get("/modbusModelRegisterGroupList", params={"modelUuid": model_uuid})
        return resp.get("data", {}).get("list", [])

    def delete_register_group(self, group_uuid: str) -> Dict[str, Any]:
        """删除寄存器组"""
        return self.client.post("/modbusModelRegisterGroupDel", json_data={"uuid": group_uuid})

    # ── Modbus 数据点（寄存器地址）───────────────────

    def add_register(
        self,
        group_uuid: str,
        name: str,
        address: int,
        data_type: int = 1,       # 1=INT16, 2=UINT16, 3=INT32, 4=UINT32, 5=FLOAT32
        unit: str = "",
        coefficient: float = 1.0,
        parse_type: int = 177,    # 177=标准解析
        **kwargs,
    ) -> Dict[str, Any]:
        """在寄存器组中添加数据点"""
        data = {
            "groupUuid": group_uuid,
            "Name": name,
            "Address": address,
            "DataType": data_type,
            "Unit": unit,
            "Coefficient": coefficient,
            "ParseType": parse_type,
            **kwargs,
        }
        ism_logger.info(f"📍 添加数据点: {name} @ {address}")
        return self.client.post("/modbusModelRegisterAdd", json_data=data)

    def list_registers(self, group_uuid: str) -> List[Dict]:
        """查询寄存器组下的数据点列表"""
        resp = self.client.get("/modbusModelRegisterList", params={"groupUuid": group_uuid})
        return resp.get("data", {}).get("list", [])

    def delete_register(self, register_uuid: str) -> Dict[str, Any]:
        """删除数据点"""
        return self.client.post("/modbusModelRegisterDel", json_data={"uuid": register_uuid})

    # ── 其他协议模型（通用接口）────────────────────────

    def _create_by_protocol(self, protocol: str, name: str, project_uuid: str, **kwargs) -> Dict[str, Any]:
        """根据协议创建模型"""
        routes = {
            "modbus": "/modbusModelAdd",
            "opcua": "/opcuaModelAdd",
            "snmp": "/snmpmodeladd",
            "dlt645": "/dlt645ModelAdd",
            "mqtt": "/mqttModelAdd",
            "s7": "/s7ModelAdd",
            "hj212": "/hj212ModelAdd",
            "iec104": "/iec104ModelAdd",
            "cjt188": "/cjt188ModelAdd",
        }
        route = routes.get(protocol.lower())
        if not route:
            raise ValueError(f"不支持的协议: {protocol}")

        data = {"Name": name, "ProjectUuid": project_uuid, **kwargs}
        ism_logger.info(f"📦 创建 {protocol.upper()} 模型: {name}")
        return self.client.post(route, json_data=data)

    def create(self, protocol: str, name: str, project_uuid: str, **kwargs) -> Dict[str, Any]:
        """通用创建模型入口"""
        return self._create_by_protocol(protocol, name, project_uuid, **kwargs)

    def list_by_protocol(self, protocol: str, project_uuid: Optional[str] = None) -> List[Dict]:
        """根据协议查询模型列表"""
        routes = {
            "modbus": "/modbusModelList",
            "opcua": "/opcuaModelList",
            "snmp": "/snmpmodellist",
            "dlt645": "/dlt64ModelList",
            "mqtt": "/mqttModelList",
            "s7": "/s7ModelList",
            "hj212": "/hj212ModelList",
            "iec104": "/iec104ModelList",
            "cjt188": "/cjt188ModelList",
        }
        route = routes.get(protocol.lower())
        if not route:
            raise ValueError(f"不支持的协议: {protocol}")

        params = {}
        if project_uuid:
            params["ProjectUuid"] = project_uuid
        resp = self.client.get(route, params=params)
        return resp.get("data", {}).get("list", [])

    # ── 批量操作 ──────────────────────────────────────

    def create_with_registers(
        self,
        protocol: str,
        name: str,
        project_uuid: str,
        register_groups: List[Dict[str, Any]],
        **kwargs,
    ) -> Dict[str, Any]:
        """
        创建模型并批量添加寄存器组和数据点

        register_groups 格式:
        [
            {
                "name": "保持寄存器组",
                "start_address": 0,
                "end_address": 100,
                "register_type": 3,
                "registers": [
                    {"name": "温度", "address": 0, "data_type": 1, "unit": "°C"},
                    ...
                ]
            }
        ]
        """
        # 1. 创建模型
        model_resp = self.create(protocol, name, project_uuid, **kwargs)
        model_data = model_resp.get("data", {})
        model_uuid = model_data.get("uuid") or model_data.get("Uuid")

        if not model_uuid:
            raise ValueError(f"创建模型后未返回 uuid: {model_resp}")

        # 2. 创建寄存器组和数据点
        created_groups = []
        for group in register_groups:
            group_resp = self.add_register_group(
                model_uuid=model_uuid,
                name=group["name"],
                start_address=group["start_address"],
                end_address=group["end_address"],
                register_type=group.get("register_type", 3),
            )
            group_data = group_resp.get("data", {})
            group_uuid = group_data.get("uuid") or group_data.get("Uuid")

            # 添加数据点
            registers = []
            for reg in group.get("registers", []):
                reg_resp = self.add_register(
                    group_uuid=group_uuid,
                    name=reg["name"],
                    address=reg["address"],
                    data_type=reg.get("data_type", 1),
                    unit=reg.get("unit", ""),
                    coefficient=reg.get("coefficient", 1.0),
                    parse_type=reg.get("parse_type", 177),
                )
                registers.append(reg_resp.get("data", {}))

            created_groups.append({
                "group": group_resp.get("data", {}),
                "registers": registers,
            })

        return {
            "model": model_data,
            "groups": created_groups,
        }
