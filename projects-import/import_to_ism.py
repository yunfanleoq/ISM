#!/usr/bin/env python3
"""
ISM 项目批量导入脚本
读取 _ISM项目包.json → 通过 HTTP API 将项目导入 ISM 系统

用法:
    python3 import_to_ism.py <项目包.json> [--base-url URL] [--token TOKEN] [--dry-run]
    
示例:
    python3 import_to_ism.py projects-import/1A配电室/1A配电室_ISM项目包.json
    
需要:
    pip install requests
"""

import json, sys, os, time, argparse, re
import requests

# ============================================
# 配置
# ============================================

DEFAULT_BASE_URL = "http://8.138.197.243:9081"
DEFAULT_TOKEN = ""
API_PREFIX = "/api"


class ISMImporter:
    def __init__(self, base_url, token=None, project_uuid=None):
        self.base_url = base_url.rstrip('/')
        self.token = token
        self.project_uuid = project_uuid
        self.session = requests.Session()
        if token:
            self.session.headers.update({"Authorization": token})
        # If no token, we'll need to login
        self.logged_in = token is not None and len(token) > 0
    
    def login(self, username="admin", password="admin"):
        """登录获取JWT token"""
        url = f"{self.base_url}/login"
        resp = self.session.post(url, json={
            "username": username,
            "password": password
        })
        if resp.status_code == 200:
            data = resp.json()
            if data.get("code") == 0:
                self.token = data.get("token", "")
                self.session.headers.update({"Authorization": self.token})
                self.logged_in = True
                print(f"  ✓ 登录成功 (用户: {username})")
                return True
        print(f"  ✗ 登录失败: {resp.text[:200]}")
        return False
    
    def _post(self, path, data, timeout=30, project_uuid=None):
        """发送POST请求"""
        url = f"{self.base_url}{path}"
        headers = {}
        if project_uuid or self.project_uuid:
            headers["ProjectUuid"] = project_uuid or self.project_uuid
        resp = self.session.post(url, json=data, headers=headers, timeout=timeout)
        try:
            return resp.json()
        except:
            return {"code": -1, "error": resp.text[:200]}
    
    def _get(self, path, params=None, timeout=30):
        """发送GET请求"""
        url = f"{self.base_url}{path}"
        resp = self.session.get(url, params=params, timeout=timeout)
        try:
            return resp.json()
        except:
            return {"code": -1, "error": resp.text[:200]}
    
    # ---- Phase 0: 创建项目 ----
    def create_project(self, name, description="", creator="admin"):
        """创建新项目"""
        print(f"\n[Phase 0] 创建项目: {name}")
        
        # 需要先获取用户UUID
        user_info = self._get("/GetUserInfo")
        creator_uuid = ""
        if user_info.get("code") == 0:
            creator_uuid = user_info.get("data", {}).get("uuid", "")
        else:
            # Fallback: try routes to get user info
            routes = self._get("/routes")
            creator_uuid = "662929b7-1d42-41fa-baa2-08b1988e4f54"  # admin uuid
        
        data = {
            "name": name,
            "description": description or f"AI批量导入项目",
            "industry": 1,
            "creator": creator,
            "creator_uuid": creator_uuid,
        }
        
        result = self._post("/ProjectAdd", data, project_uuid="")
        if result.get("code") == 0:
            # 获取刚创建的项目 uuid
            proj_list = self._post("/ProjectList", {}, project_uuid="")
            if proj_list.get("code") == 0:
                projects = proj_list.get("list", [])
                for p in projects:
                    pi = p.get("ProjectInfo", p)
                    if pi.get("name") == name:
                        self.project_uuid = pi.get("uuid", "")
                        print(f"  ✓ 项目创建成功: {self.project_uuid}")
                        return self.project_uuid
        else:
            # 可能已存在，尝试获取
            proj_list = self._post("/ProjectList", {}, project_uuid="")
            if proj_list.get("code") == 0:
                projects = proj_list.get("list", [])
                for p in projects:
                    pi = p.get("ProjectInfo", p)
                    if pi.get("name") == name:
                        self.project_uuid = pi.get("uuid", "")
                        print(f"  ✓ 项目已存在，复用: {self.project_uuid}")
                        return self.project_uuid
            print(f"  ✗ 项目创建失败: {result}")
        return None
    
    # ---- Phase 1: 创建数据模型 ----
    def create_data_models(self, models):
        """批量创建 Modbus 数据模型"""
        print(f"\n[Phase 1] 创建 {len(models)} 个数据模型")
        model_uuids = {}
        
        for dm in models:
            name = dm["name"]
            print(f"  创建模型: {name}...", end=" ")
            
            # 1. 创建 DevicesModel
            data = {
                "name": name,
                "dec": f"AI批量导入 - {dm['protocol']}",
                "type": 2,  # Modbus
                "gatherNumber": 30,
                "project_uuid": self.project_uuid,
            }
            
            # Modbus connection params
            conn = dm.get("connection", {})
            data["modbusConnectType"] = conn.get("connect_type", "TCP")
            data["modbusClientIpaddress"] = conn.get("ip", "172.31.4.14")
            data["port"] = conn.get("port", 502)
            data["timeout"] = 5
            data["modbusConnectMode"] = "TCP"
            data["serialBaud"] = 9600
            data["serialBits"] = 8
            data["serialParity"] = "None"
            data["serialStopBits"] = "1"
            
            result = self._post("/modbusModelAdd", data)
            code = result.get("code", -1)
            
            if code in [0, 200, -4]:  # -4 may mean "already exists"
                # Get model uuid
                muid = ""
                model_list = self._post("/modbusModelList", {}, project_uuid="")
                if model_list.get("code") == 0:
                    for m in model_list.get("list", []):
                        if m.get("name") == name:
                            muid = m.get("uuid", "")
                            break
                
                if not muid:
                    muid = f"auto_{name}"
                
                model_uuids[name] = muid
                
                # 2. 创建寄存器组
                for grp in dm.get("register_groups", []):
                    grp_data = {
                        "muid": muid,
                        "name": grp.get("group_name", "数据组"),
                        "function": 3 if grp.get("register_type") == "holding" else 1,
                        "registerStart": 0,
                        "registerCount": len(grp.get("addresses", [])),
                    }
                    self._post("/modbusModelRegisterGroupAdd", grp_data)
                    
                    # 3. 创建寄存器地址
                    for addr in grp.get("addresses", []):
                        addr_data = {
                            "muid": muid,
                            "name": addr.get("name", "数据点"),
                            "offset": addr.get("offset", 0),
                            "dataType": addr.get("data_type", "uint16"),
                            "dataUnit": addr.get("unit", ""),
                            "conversionExpression": addr.get("conversion_expression", "x"),
                            "isAlarm": 1 if addr.get("is_alarm") else 0,
                            "alarmLevel": addr.get("alarm_level", 0),
                            "alarmMessage": addr.get("alarm_message", ""),
                            "isRecord": 1 if addr.get("is_record") else 0,
                        }
                        self._post("/modbusModelRegisterAdd", addr_data, timeout=5)
                
                print(f"✓ ({muid[:8]}...)")
            else:
                print(f"✗ code={code}")
        
        return model_uuids
    
    # ---- Phase 2-3: 创建设备 + 数据点 ----
    def create_devices(self, devices, model_uuids):
        """批量创建设备实例"""
        print(f"\n[Phase 2-3] 创建 {len(devices)} 台设备实例...")
        device_count = 0
        point_count = 0
        
        for i, dev in enumerate(devices):
            muid = model_uuids.get(dev.get("model_name", ""), "")
            if not muid:
                continue
            
            name = dev["name"]
            extra = json.dumps({
                "ai_start": dev.get("register_offset"),
                "di_start": dev.get("di_start"),
                "model_name": dev.get("model_name"),
                "location": dev.get("location", {}),
            })
            
            data = {
                "name": name,
                "sid": i + 2,  # 自增，1已被RootZone占用
                "pid": 1,      # 父节点是RootZone
                "type": 1,     # 设备
                "deviceType": 2,  # Modbus
                "muid": muid,
                "uuid": "",
                "extra": extra,
                "interval": 5,
                "timeout": 5,
                "failedTimes": 5,
                "IsEnable": 1,
                "description": dev.get("display_name", name),
                "project_uuid": self.project_uuid,
                "offlineClear": 0,
                "offlineDefaultValue": "",
                "Status": 0,
                "longitude": "",
                "latitude": "",
            }
            
            result = self._post("/monitorAdd", data)
            if result.get("code") in [0, 200]:
                device_count += 1
                point_count += dev.get("data_point_count", 0)
            
            if (i + 1) % 10 == 0:
                print(f"  进度: {i+1}/{len(devices)} (已成功{device_count}台, ~{point_count}数据点)")
        
        print(f"  完成: {device_count}台设备, ~{point_count}数据点")
        return device_count, point_count
    
    # ---- Phase 4: 创建告警触发器 ----
    def create_alarm_triggers(self, triggers):
        """批量创建告警触发器"""
        print(f"\n[Phase 4] 创建 {len(triggers)} 个告警触发器...")
        alarm_count = 0
        
        for trig in triggers:
            data = {
                "triggerName": trig.get("name", "告警触发器"),
                "triggerCondition": trig.get("condition", "x == 0"),
                "triggerAlarmLevel": trig.get("level", 3),
                "triggerKeepTime": trig.get("keep_time", 10),
                "triggerAlarmShowText": trig.get("alarm_message", "设备告警"),
                "triggerAlarmHideText": "告警已消除",
                "triggerType": 1,
                "triggerXValue": "x",
                "project_uuid": self.project_uuid,
                "triggerDeviceUuid": "",
                "triggerDataUuid": "",
                "triggerDeviceType": 2,
                "triggerDeviceModelUuid": "",
                "triggerModelDataUuid": "",
            }
            
            result = self._post("/AlarmTriggerAdd", data)
            if result.get("code") in [0, 200]:
                alarm_count += 1
        
        print(f"  完成: {alarm_count}个告警触发器")
        return alarm_count
    
    # ---- Phase 5: 创建组态大屏 ----
    def create_dashboard(self, dashboard_config, devices):
        """创建组态大屏"""
        print(f"\n[Phase 5] 创建组态大屏...")
        
        dash_name = dashboard_config.get("name", "监控大屏")
        
        # 创建组态模型
        data = {
            "name": dash_name,
            "description": "AI智能生成",
            "displayType": 1,
            "DisplayImage": "",
        }
        result = self._post("/displayModelAdd", data)
        
        if result.get("code") not in [0, 200]:
            print(f"  ✗ 创建失败: {result}")
            return None
        
        # 获取 dashboard UID
        dash_uid = ""
        dash_list = self._post("/displayModelList", {"DisplayType": 1})
        if dash_list.get("code") == 0:
            for d in dash_list.get("list", []):
                if d.get("name") == dash_name:
                    dash_uid = d.get("displayUid", "")
                    break
        
        if not dash_uid:
            print(f"  ✗ 未找到刚创建的组态")
            return None
        
        print(f"  ✓ 组态创建成功: {dash_uid}")
        
        # 保存总览页面布局
        overview = dashboard_config.get("overview_page", {})
        layer_config = {
            "backColor": "#0a1628",
            "backgroundImage": "",
            "width": 1920,
            "height": 1080,
            "autoSize": 0,
        }
        
        components = generate_dashboard_components(dashboard_config, devices)
        
        save_data = {
            "muid": dash_uid,
            "pageid": "demo_page",
            "saveData": {
                "name": overview.get("name", "总览"),
                "layer": layer_config,
                "components": components,
            }
        }
        
        result = self._post("/saveDisplayModelLayerData", save_data, timeout=30)
        if result.get("code") in [0, 200]:
            print(f"  ✓ 页面布局保存成功")
        
        dash_url = f"{self.base_url}/#/AppRun/{dash_uid}"
        print(f"\n  🎯 大屏地址: {dash_url}")
        return dash_uid


def generate_dashboard_components(dashboard, devices):
    """根据设备列表生成组态组件JSON"""
    components = []
    
    # 标题
    components.append({
        "type": "DvBorderBox1",
        "identifier": "header",
        "name": "header",
        "style": {
            "position": {"x": 0, "y": 0, "w": 1920, "h": 70},
            "diy": [
                {"name": "标题文本", "type": 4, "value": dashboard.get("name", "监控总览"), "key": "text"},
                {"name": "前景色", "type": 2, "value": "#00d4ff", "key": "foreColor"},
                {"name": "字体大小", "type": 1, "value": 32, "key": "fontSize"},
            ],
            "zIndex": 100,
            "visible": 1,
            "opacity": 1,
            "transform": 0,
        },
        "action": [],
    })
    
    # 设备标签
    x, y = 20, 90
    col_count = 0
    for dev in devices:
        components.append({
            "type": "ViewSvgText",
            "identifier": f"label_{dev['name']}",
            "name": dev["name"],
            "style": {
                "position": {"x": x, "y": y, "w": 180, "h": 30},
                "diy": [
                    {"name": "文本", "type": 4, "value": dev.get("display_name", dev["name"]), "key": "text"},
                    {"name": "前景色", "type": 2, "value": "#00d4ff", "key": "foreColor"},
                    {"name": "字体大小", "type": 1, "value": 13, "key": "fontSize"},
                ],
                "zIndex": 10,
                "visible": 1,
                "opacity": 1,
                "transform": 0,
            },
            "action": [],
            "dataBind": [],
        })
        
        y += 35
        col_count += 1
        if col_count >= 25:
            col_count = 0
            y = 90
            x += 200
    
    return components


# ============================================
# Main
# ============================================

def main():
    parser = argparse.ArgumentParser(description="ISM 项目批量导入工具")
    parser.add_argument("package_file", help="ISM项目包JSON文件路径")
    parser.add_argument("--base-url", default=DEFAULT_BASE_URL, help=f"ISM服务地址 (默认: {DEFAULT_BASE_URL})")
    parser.add_argument("--token", default="", help="API Token（可选，不提供则尝试登录）")
    parser.add_argument("--username", default="admin", help="登录用户名")
    parser.add_argument("--password", default="admin", help="登录密码")
    parser.add_argument("--dry-run", action="store_true", help="干跑模式，不实际调用API")
    parser.add_argument("--skip-alarms", action="store_true", help="跳过告警创建")
    parser.add_argument("--skip-dashboard", action="store_true", help="跳过大屏创建")
    
    args = parser.parse_args()
    
    print("=" * 60)
    print("ISM 项目批量导入工具")
    print(f"  包文件: {args.package_file}")
    print(f"  目标: {args.base_url}")
    print("=" * 60)
    
    # 读取项目包
    if not os.path.exists(args.package_file):
        print(f"✗ 文件不存在: {args.package_file}")
        sys.exit(1)
    
    with open(args.package_file, 'r', encoding='utf-8') as f:
        pkg = json.load(f)
    
    meta = pkg.get("_meta", {})
    project = pkg.get("project", {})
    models = pkg.get("data_models", [])
    devices = pkg.get("devices", [])
    alarms = pkg.get("alarm_triggers", [])
    dashboard = pkg.get("dashboard", {})
    
    print(f"\n项目: {project.get('name', 'Unknown')}")
    print(f"数据模型: {len(models)}")
    print(f"设备: {len(devices)}")
    print(f"告警: {len(alarms)}")
    print(f"大屏: {dashboard.get('name', '无')}")
    
    if args.dry_run:
        print("\n[干跑模式] 仅验证，不执行导入")
        return
    
    # 初始化导入器
    importer = ISMImporter(args.base_url, args.token)
    
    # 尝试连接
    print(f"\n[连接测试] 尝试连接 {args.base_url}...")
    try:
        test = importer._get("/login", params={"test": "1"})
        print(f"  响应: {test.get('code', 'no response')}")
    except Exception as e:
        print(f"  ✗ 连接失败: {e}")
        print(f"  提示: 请确保 ISM 服务已启动且地址正确")
        print(f"  可手动指定: --base-url http://your-server:port")
        sys.exit(1)
    
    # 登录
    if not importer.logged_in:
        print(f"\n[登录] 尝试登录 (用户: {args.username})")
        if not importer.login(args.username, args.password):
            print("✗ 无法登录，请检查用户名/密码")
            sys.exit(1)
    
    # ---- 执行导入 ----
    
    # Phase 0: 创建项目
    project_name = project.get("name", "导入项目")
    proj_uuid = importer.create_project(project_name, project.get("description", ""))
    if not proj_uuid:
        print("✗ 项目创建失败，终止")
        sys.exit(1)
    
    # Phase 1: 创建数据模型
    model_uuids = importer.create_data_models(models)
    
    # Phase 2-3: 创建设备
    dev_count, pt_count = importer.create_devices(devices, model_uuids)
    
    # Phase 4: 创建告警
    if not args.skip_alarms:
        alarm_count = importer.create_alarm_triggers(alarms)
    
    # Phase 5: 创建大屏
    if not args.skip_dashboard:
        dash_uid = importer.create_dashboard(dashboard, devices)
    
    print("\n" + "=" * 60)
    print("导入完成！")
    print(f"  项目UUID: {proj_uuid}")
    print(f"  数据模型: {len(models)}个")
    print(f"  设备: {dev_count}台")
    print(f"  数据点: ~{pt_count}个")
    print("=" * 60)


if __name__ == "__main__":
    main()
