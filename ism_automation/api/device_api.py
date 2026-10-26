"""
ISM Automation — 设备管理 API
封装设备树操作：添加设备/区域、编辑、删除、查询实时数据等。
"""
from __future__ import annotations

from typing import Any, Dict, List, Optional

from ism_automation.core.logger import ism_logger
from ism_automation.api.client import ISMClient


class DeviceAPI:
    """设备管理 API 封装"""

    def __init__(self, client: ISMClient):
        self.client = client

    # ── 设备树操作 ────────────────────────────────────

    def add(
        self,
        name: str,
        project_uuid: str,
        parent_sid: int = 0,
        device_type: int = 1,  # 0=区域/分组, 1=设备
        model_uuid: Optional[str] = None,
        protocol_type: int = 2,  # 2=ModbusTCP
        ip: Optional[str] = None,
        port: Optional[int] = None,
        slave_id: Optional[int] = None,
        status: int = 1,
        **kwargs,
    ) -> Dict[str, Any]:
        """
        添加设备或区域到设备树

        Args:
            name: 设备名称
            project_uuid: 项目UUID
            parent_sid: 父节点 sid (0=根节点)
            device_type: 0=区域/分组, 1=设备
            model_uuid: 设备模型UUID (type=1时必填)
            protocol_type: 通信协议类型
            ip: 设备IP地址
            port: 设备端口
            slave_id: Modbus 从机地址
            status: 1=启用, 0=停用
        """
        data = {
            "Name": name,
            "ProjectUuid": project_uuid,
            "Pid": parent_sid,
            "Type": device_type,
            "Muid": model_uuid or "",
            "ProtocolType": protocol_type,
            "Ip": ip or "",
            "Port": port or 0,
            "SlaveId": slave_id or 1,
            "Status": status,
            **kwargs,
        }
        ism_logger.info(f"📦 添加设备/区域: {name} (pid={parent_sid}, type={device_type})")
        return self.client.post("/monitorAdd", json_data=data)

    def edit(self, sid: int, **kwargs) -> Dict[str, Any]:
        """编辑设备或区域"""
        data = {"Sid": sid, **kwargs}
        return self.client.post("/monitorEdit", json_data=data)

    def delete(self, sid: int) -> Dict[str, Any]:
        """删除设备或区域"""
        ism_logger.warning(f"🗑️ 删除设备: sid={sid}")
        return self.client.post("/monitorDel", json_data={"Sid": sid})

    def copy(self, sid: int, new_name: Optional[str] = None) -> Dict[str, Any]:
        """复制设备"""
        data = {"Sid": sid}
        if new_name:
            data["Name"] = new_name
        return self.client.post("/monitorCopy", json_data=data)

    def get_tree(self, project_uuid: str) -> List[Dict]:
        """获取设备树"""
        resp = self.client.get("/monitortree", params={"ProjectUuid": project_uuid})
        return resp.get("data", {}).get("tree", [])

    # ── 批量操作 ──────────────────────────────────────

    def batch_add(
        self,
        devices: List[Dict[str, Any]],
        project_uuid: str,
        parent_sid: int = 0,
    ) -> List[Dict[str, Any]]:
        """批量添加设备"""
        results = []
        for dev in devices:
            result = self.add(
                name=dev["name"],
                project_uuid=project_uuid,
                parent_sid=dev.get("parent_sid", parent_sid),
                device_type=dev.get("device_type", 1),
                model_uuid=dev.get("model_uuid"),
                protocol_type=dev.get("protocol_type", 2),
                ip=dev.get("ip"),
                port=dev.get("port"),
                slave_id=dev.get("slave_id"),
                status=dev.get("status", 1),
            )
            results.append(result)
        return results

    def batch_delete(self, sids: List[int]) -> List[Dict[str, Any]]:
        """批量删除设备"""
        return [self.delete(sid) for sid in sids]

    def batch_set_status(self, sids: List[int], status: int) -> Dict[str, Any]:
        """批量设置设备状态"""
        return self.client.post("/MonitorBatchSetStatus", json_data={
            "Sids": sids,
            "Status": status,
        })

    # ── 实时数据 ──────────────────────────────────────

    def get_real_data(self, sid: int) -> Dict[str, Any]:
        """获取设备实时数据"""
        return self.client.get("/getRealData", params={"Sid": sid})

    def get_real_data_by_uuid(self, uuid: str) -> Dict[str, Any]:
        """通过 UUID 获取实时数据"""
        return self.client.get("/getRealDataByUuid", params={"Uuid": uuid})

    def set_data(self, sid: int, data_id: str, value: Any) -> Dict[str, Any]:
        """写入设备数据（需支持写入）"""
        return self.client.post("/setData", json_data={
            "Sid": sid,
            "DataId": data_id,
            "Value": value,
        })

    def sync_real_data(self, project_uuid: str) -> Dict[str, Any]:
        """同步设备实时数据表"""
        return self.client.post("/syncDeviceRealData", json_data={
            "ProjectUuid": project_uuid,
        })

    # ── 数据模型绑定 ──────────────────────────────────

    def get_data_model_data(self, model_uuid: str) -> List[Dict]:
        """获取设备模型的数据点列表"""
        resp = self.client.get("/getDataModelData", params={"Muid": model_uuid})
        return resp.get("data", {}).get("list", [])

    def get_device_model_data_list(self, project_uuid: str) -> List[Dict]:
        """获取项目下所有设备模型数据点"""
        resp = self.client.get("/GetDeviceModelDataList", params={"ProjectUuid": project_uuid})
        return resp.get("data", {}).get("list", [])
