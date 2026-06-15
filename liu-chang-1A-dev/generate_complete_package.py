import json
import uuid
from pathlib import Path
from datetime import datetime

# 读取已生成的数据模型
with open('/Users/yunfanleo/cursorProjects/ISM源码/liu-chang-1A-dev/ism_data_models.json', 'r', encoding='utf-8') as f:
    data = json.load(f)

# 生成UUID
new_uuid = lambda: str(uuid.uuid4()).replace('-', '')

# ─── parse_mode → (Go data_type, register_count, unit) 映射 ───
# 对应 Go 后端 modbusPthread.go 中 data.Type 的分支逻辑
#   177 → uint16, 1寄存器（Unsigned short）
#   179 → int32,  2寄存器（Long）
#    73 → int16,  1寄存器（Short, 谐波畸变率）
#    71 → uint16, 1寄存器（Unsigned short, UPS参数）
#     1 → uint16, 1寄存器（code_1原值, UPS使用）
# 默认 → Float,  2寄存器
PARSE_MODE_TYPE_MAP = {
    177: ('Unsigned short', 1),
    179: ('Long', 2),
    73:  ('Short', 1),
    71:  ('Unsigned short', 1),
    1:   ('Unsigned short', 1),
}

def infer_register_type(parse_mode, data_name=''):
    """根据 parse_mode 推断 Go 后端所需的 data_type 和 register_count"""
    try:
        pm = int(parse_mode) if parse_mode != '' else -1
    except (ValueError, TypeError):
        pm = -1
    if pm in PARSE_MODE_TYPE_MAP:
        return PARSE_MODE_TYPE_MAP[pm]
    # 空 parse_mode → Short (状态类)
    if pm == -1:
        return ('Short', 1)
    return ('Float', 2)

def infer_unit(data_name):
    """根据数据点名推断单位"""
    if '电流' in data_name: return 'A'
    if '电压' in data_name: return 'V'
    if '功率' in data_name and '因数' not in data_name: return 'kW'
    if '频率' in data_name: return 'Hz'
    if '电度' in data_name: return 'kWh'
    if '谐波' in data_name or '畸变率' in data_name: return '%'
    if '温度' in data_name: return '°C'
    if '时间' in data_name: return 'min'
    return ''

# 项目基础信息
project_uuid = new_uuid()
project_name = "1A配电室"

# 先为组态大屏生成 UUID（后续构造设备模型/监控树节点时直接填入）
display_model_uuid = new_uuid()
display_uid = new_uuid()
page_id = new_uuid()

# ==================== 1. 设备模型 (DevicesModel) ====================
def generate_device_model(template):
    """生成ISM后端DevicesModel格式"""
    model_uuid = new_uuid()
    return {
        "uuid": model_uuid,
        "name": template['type'],
        "dec": f"{template['type']} 电力设备模型",
        "type": 2,  # Modbus = 2
        "gatherNumber": 30,
        "project_uuid": project_uuid,
        "configUid": display_uid,
        "PageUUID": page_id,
        "version": 1,
        "port": 502,
        "timeout": 5,
        "DataFormat": "CDAB",
        "modbusConnectType": "TCPClient",  # TCP/RTU
        "modbusConnectMode": "TCP/IP",
        "modbusClientIpaddress": "127.0.0.1",
        "modbusCom": "",
        "serialBaud": 9600,
        "serialBits": 8,
        "serialParity": "None",
        "serialStopBits": "1",
        "serialFlow": "None"
    }, model_uuid

# 生成3种设备模型
device_models = []
model_uuid_map = {}  # template_type -> uuid
for template in data['templates']:
    model, muid = generate_device_model(template)
    device_models.append(model)
    model_uuid_map[template['type']] = muid

# ==================== 2. 寄存器组 (ModbusDevicesRegisterGroup) ====================
def generate_register_group(model_name, muid, group_name, function_code, start_addr, count, explicit_uuid=None):
    """生成寄存器组"""
    return {
        "uuid": explicit_uuid if explicit_uuid else new_uuid(),
        "muid": muid,
        "name": group_name,
        "function": function_code,  # 3=Holding, 2=DI
        "registerStart": start_addr,
        "registerCount": count
    }

register_groups = []
# ★ 关键修复: 按模板类型追踪各模板的 AI/DI 寄存器组 UUID ★
# 坑点: 之前使用全局 list，所有数据点都引用 register_groups[0]（A20的组），
# 导致 A40 和 UPS 的数据点被错误挂到 A20 的寄存器组下，offset 超出 range 的寄存器读不到数据
template_groups = {}  # {template_type: (ai_group_uuid, di_group_uuid)} 或 (None, None)
for template in data['templates']:
    muid = model_uuid_map[template['type']]
    ai_count = template['aiCount']
    di_count = template['diCount']
    
    ai_uuid = None
    di_uuid = None
    
    if ai_count > 0:
        ai_uuid = new_uuid()
        register_groups.append(generate_register_group(
            template['type'], muid, "AI数据", 3, 0, ai_count, ai_uuid
        ))
    if di_count > 0:
        di_uuid = new_uuid()
        register_groups.append(generate_register_group(
            template['type'], muid, "DI数据", 2, 0, di_count, di_uuid
        ))
    
    template_groups[template['type']] = (ai_uuid, di_uuid)

# ==================== 3. 寄存器数据点 (ModbusDevicesDataModel) ====================
def generate_register_point(name, register_addr, group_uuid, muid, data_type, unit, factor, parse_mode, is_alarm=0, alarm_level=0):
    """生成寄存器数据点"""
    return {
        "uuid": new_uuid(),
        "muid": muid,
        "name": name,
        "registerAddress": register_addr,
        "registerGroupUuid": group_uuid,
        "auth": "ReadOnly",
        "type": data_type,  # Short, Float, Bool
        "ByteOrder": "CDAB",
        "modeltype": 2,
        "unit": unit,
        "conversionExpression": f"{{val}}*{factor}" if factor else "",
        "alarm": is_alarm,
        "alarmLevel": alarm_level,
        "AlarmMessage": f"{name}告警" if is_alarm else "",
        "AlarmClearMessage": f"{name}告警消除" if is_alarm else "",
        "record": 1 if is_alarm else 0,
        "RecordType": 1,
        "RecordInterval": 5,
        "RecordDataCharge": "",
        "RecordDataTimely": "",
        "FloatAccuracy": "0.01"
    }

register_points = []
# 为每个模板生成数据点
for template in data['templates']:
    muid = model_uuid_map[template['type']]
    ttype = template['type']
    ai_guid, di_guid = template_groups.get(ttype, (None, None))
    
    # AI数据点
    for i, ai in enumerate(template['aiParams']):
        go_type, reg_count = infer_register_type(ai.get('parseMode', ''), ai.get('name', ''))
        unit = infer_unit(ai.get('name', ''))
        
        # 告警判断
        is_alarm = 1 if any(k in ai.get('name', '') for k in ['电流', '电压', '功率', '温度', '故障']) else 0
        alarm_level = 2 if '故障' in ai.get('name', '') else 1
        
        register_points.append(generate_register_point(
            ai['name'], ai['offset'], ai_guid,
            muid, go_type, unit, ai.get('factor', ''), ai.get('parseMode', ''),
            is_alarm, alarm_level
        ))
    
    # DI数据点
    for di in template['diParams']:
        register_points.append(generate_register_point(
            di['name'], di['registerOffset'], di_guid,
            muid, "Bool", "", "", "",
            1, 3  # DI默认告警，严重级别
        ))

# ==================== 4. 设备树/监控列表 (MonitorList) ====================
# 创建楼层/配电柜结构
def generate_zone(pid, sid, name, zone_type=0):
    """生成区域/楼层节点"""
    return {
        "uuid": new_uuid(),
        "sid": sid,
        "pid": pid,
        "name": name,
        "type": zone_type,  # 0=区域, 1=设备
        "timeout": 5,
        "IsEnable": 1,
        "project_uuid": project_uuid,
        "interval": 5000,
        "failedTimes": 5,
        "description": name,
        "offlineClear": 0,
        "offlineDefaultValue": "0",
        "deviceType": 0,
        "muid": "",
        "configUid": display_uid,
        "PageUUID": page_id,
        "extra": "",
        "Status": 0,
        "longitude": "",
        "latitude": ""
    }

def generate_device(pid, sid, name, device_type, muid, extra_data):
    """生成设备节点"""
    return {
        "uuid": new_uuid(),
        "sid": sid,
        "pid": pid,
        "name": name,
        "type": 1,  # 1=设备
        "timeout": 5,
        "IsEnable": 1,
        "project_uuid": project_uuid,
        "interval": 5000,
        "failedTimes": 5,
        "description": name,
        "offlineClear": 0,
        "offlineDefaultValue": "0",
        "deviceType": device_type,  # 2=Modbus
        "muid": muid,
        "configUid": display_uid,
        "PageUUID": page_id,
        "extra": extra_data,  # Modbus地址等配置
        "Status": 2,
        "longitude": "",
        "latitude": ""
    }

# 构建楼层/配电柜/设备树
monitor_list = []
sid_counter = 1000

# 根节点：1A配电室
root_sid = sid_counter
monitor_list.append(generate_zone(0, root_sid, "1A配电室"))
sid_counter += 1

# 楼层：1楼
floor1_sid = sid_counter
monitor_list.append(generate_zone(root_sid, floor1_sid, "1楼配电室"))
sid_counter += 1

# 按源节点分组创建设备柜
source_nodes = {
    "P1_仪表_": {"count": 306, "template": "A20", "cabinet": "1A1_U11柜"},
    "P2_仪表_": {"count": 306, "template": "A20", "cabinet": "1A1_U11柜"},
    "P4_仪表_": {"count": 306, "template": "A20", "cabinet": "1A1_U11柜"},
    "P5_仪表_": {"count": 306, "template": "A20", "cabinet": "1A1_U11柜"},
    "P7_仪表_": {"count": 306, "template": "A20", "cabinet": "1A3_U11柜"},
    "P3_仪表_": {"count": 288, "template": "A40", "cabinet": "1A1_U11柜"},
    "P6_仪表_": {"count": 238, "template": "A20", "cabinet": "1A3_U11柜"},
    "P9_仪表_": {"count": 204, "template": "A20", "cabinet": "1A3_U11柜"},
    "P10_施耐德UPS_": {"count": 180, "template": "施耐德UPS", "cabinet": "UPS柜"},
    "P8_仪表_": {"count": 176, "template": "A40", "cabinet": "1A3_U11柜"},
    "P11_施耐德UPS_": {"count": 135, "template": "施耐德UPS", "cabinet": "UPS柜"},
}

# 创建配电柜节点
cabinets = {}
for src, info in source_nodes.items():
    cabinet_name = info['cabinet']
    if cabinet_name not in cabinets:
        cabinet_sid = sid_counter
        cabinets[cabinet_name] = cabinet_sid
        monitor_list.append(generate_zone(floor1_sid, cabinet_sid, cabinet_name))
        sid_counter += 1

# 创建设备实例
devices = []
# ★ 按模型类型排序分配 slave ID（匹配模拟器布局：A20→1-60, A40→61-69, UPS→70-76）
template_order = {'A20': 0, 'A40': 1, '施耐德UPS': 2}
sorted_devices = sorted(data['devices'], key=lambda d: template_order.get(d['templateType'], 99))
slave_by_model = {'A20': 1, 'A40': 61, '施耐德UPS': 70}
next_slave = {'A20': 1, 'A40': 61, '施耐德UPS': 70}

for dev in sorted_devices:
    template_type = dev['templateType']
    muid = model_uuid_map.get(template_type, "")
    
    # ★ 按设备名前缀确定所属配电柜（而非模板类型）
    name = dev['name']
    if name.startswith('1A1_U11_') or name.startswith('1A1_U12_') or name.startswith('1A1_U13_') or name.startswith('1A1_U14_'):
        cabinet_name = "1A1_U11柜"
    elif name.startswith('1A3_U11_') or name.startswith('1A3_U12_') or name.startswith('1A3_U13_'):
        cabinet_name = "1A3_U11柜"
    elif name.startswith('UPS_'):
        cabinet_name = "UPS柜"
    else:
        cabinet_name = "其他柜"
    
    cabinet_sid = cabinets.get(cabinet_name, floor1_sid)
    
    # ★ 使用模拟器地址 127.0.0.1，按模型分配 slaveId
    sim_slave_id = next_slave.get(template_type, 1)
    extra_data = json.dumps({
        "aiStartAddr": dev['aiStartAddr'],
        "diStartAddr": dev['diStartAddr'],
        "modbusIP": "127.0.0.1",
        "modbusPort": 502,
        "slaveId": sim_slave_id,
        "packTime": 5000
    })
    
    device = generate_device(cabinet_sid, sid_counter, name, 2, muid, extra_data)
    devices.append(device)
    monitor_list.append(device)
    sid_counter += 1
    next_slave[template_type] += 1

# ==================== 5. 告警触发器 (AlarmTrigger) ====================
alarm_triggers = []

# 为每个设备模型的关键数据点创建告警触发器
critical_params = {
    "A20": ["A相电流", "B相电流", "C相电流", "AB线电压", "BC线电压", "CA线电压"],
    "A40": ["A相电流", "B相电流", "C相电流", "AB线电压", "BC线电压", "CA线电压"],
    "施耐德UPS": ["主路A相电流", "输出A相电流", "电池电压", "电池温度"]
}

for template in data['templates']:
    template_type = template['type']
    muid = model_uuid_map[template_type]
    
    for ai in template['aiParams']:
        if ai['name'] in critical_params.get(template_type, []):
            # 判断告警条件
            if '电流' in ai['name']:
                condition = ">"
                x_value = "100"  # 超过100A告警
                alarm_text = f"{ai['name']}超过100A，设备过载"
                clear_text = f"{ai['name']}恢复正常"
                level = 2
            elif '电压' in ai['name']:
                condition = "!="
                x_value = "380"
                alarm_text = f"{ai['name']}异常，偏离380V"
                clear_text = f"{ai['name']}恢复正常"
                level = 1
            elif '温度' in ai['name']:
                condition = ">"
                x_value = "40"
                alarm_text = f"{ai['name']}超过40°C，设备过热"
                clear_text = f"{ai['name']}温度恢复正常"
                level = 3
            else:
                continue
            
            alarm_triggers.append({
                "Uuid": new_uuid(),
                "TriggerName": f"{ai['name']}告警",
                "project_uuid": project_uuid,
                "TriggerDeviceUuid": "",  # 绑定所有设备
                "TriggerDeviceName": "",
                "TriggerDataUuid": "",
                "TriggerDeviceType": 2,  # Modbus
                "TriggerDeviceModelUuid": muid,
                "TriggerModelDataUuid": "",
                "TriggerAlarmHideText": clear_text,
                "TriggerAlarmShowText": alarm_text,
                "TriggerCondition": condition,
                "TriggerXValue": x_value,
                "TriggerYValue": "",
                "TriggerAlarmLevel": level,
                "TriggerKeepTime": 5,
                "TriggerLinkDeviceType": 0,
                "TriggerLinkdeviceModelUuid": "",
                "TriggerLinkModelDataUuid": "",
                "TriggerLinkageAlarmValue": "",
                "TriggerLinkageAlarmClearValue": "",
                "TriggerType": 0
            })

# 添加通讯状态告警
for template_type, muid in model_uuid_map.items():
    alarm_triggers.append({
        "Uuid": new_uuid(),
        "TriggerName": f"{template_type}通讯中断",
        "project_uuid": project_uuid,
        "TriggerDeviceUuid": "",
        "TriggerDeviceName": "",
        "TriggerDataUuid": "",
        "TriggerDeviceType": 2,
        "TriggerDeviceModelUuid": muid,
        "TriggerModelDataUuid": "",
        "TriggerAlarmHideText": "设备通讯恢复",
        "TriggerAlarmShowText": "设备通讯中断，请检查连接",
        "TriggerCondition": "==",
        "TriggerXValue": "0",
        "TriggerYValue": "",
        "TriggerAlarmLevel": 4,  # 致命
        "TriggerKeepTime": 10,
        "TriggerLinkDeviceType": 0,
        "TriggerLinkdeviceModelUuid": "",
        "TriggerLinkModelDataUuid": "",
        "TriggerLinkageAlarmValue": "",
        "TriggerLinkageAlarmClearValue": "",
        "TriggerType": 0
    })

# ==================== 6. 组态大屏布局 (DisplayModelLayer) ====================
# ============================================
# 生成大屏组态布局 (ISM 标准 cell 格式)
# cell 必须包含: shape, id, x, y, width, height, zIndex, visible,
#               position, size, data.detail 完整结构
# 缺少 shape → normalizeISMCell 返回 null → 画布空
# ============================================

def make_text_cell(cid, text, x, y, w=300, h=36, font_size=16, font_weight="normal",
                   color="#e0e8f0", text_align="left", name="", z_index=-1):
    """生成 ISM view-svg-text 组件的标准 cell"""
    if not name:
        name = text[:20]
    return {
        "shape": "view-svg-text",
        "id": cid,
        "x": x, "y": y,
        "width": w, "height": h,
        "zIndex": z_index,
        "visible": True,
        "position": {"x": x, "y": y},
        "size": {"width": w, "height": h},
        "data": {
            "detail": {
                "type": "view-svg-text",
                "name": name,
                "identifier": cid,
                "style": {
                    "position": {"x": x, "y": y, "w": w, "h": h},
                    "text": text,
                    "visible": 1,
                    "fontSize": font_size,
                    "fontWeight": font_weight,
                    "foreColor": color,
                    "textAlign": text_align,
                    "zIndex": z_index
                },
                "animate": {
                    "condition": {
                        "deviceSN": "",
                        "selectVideoType": 0,
                        "isBandDevice": False,
                        "bandType": 1,
                        "dataID": "",
                        "dataName": "",
                        "operator": "",
                        "OperatorValue": "",
                        "OperatorMaxValue": ""
                    },
                    "isExpression": False,
                    "animateList": [],
                    "animateElement": []
                },
                "action": [],
                "active": [],
                "dataBind": []
            }
        }
    }

def make_rect_cell(cid, name, x, y, w, h, fill="#1a2a40", stroke="#334466",
                   stroke_width=1, z_index=-2, label=""):
    """生成 ISM view-svg-rect 组件的标准 cell (带标签)"""
    cell = {
        "shape": "view-svg-rect",
        "id": cid,
        "x": x, "y": y,
        "width": w, "height": h,
        "zIndex": z_index,
        "visible": True,
        "position": {"x": x, "y": y},
        "size": {"width": w, "height": h},
        "data": {
            "detail": {
                "type": "view-svg-rect",
                "name": name,
                "identifier": cid,
                "style": {
                    "position": {"x": x, "y": y, "w": w, "h": h},
                    "fill": fill,
                    "stroke": stroke,
                    "strokeWidth": stroke_width,
                    "zIndex": z_index
                },
                "animate": {
                    "condition": {
                        "deviceSN": "",
                        "selectVideoType": 0,
                        "isBandDevice": False,
                        "bandType": 1,
                        "dataID": "",
                        "dataName": "",
                        "operator": "",
                        "OperatorValue": "",
                        "OperatorMaxValue": ""
                    },
                    "isExpression": False,
                    "animateList": [],
                    "animateElement": []
                },
                "action": [],
                "active": [],
                "dataBind": []
            }
        }
    }
    return cell

components = []

# === 大屏标题区 ===
components.append(make_text_cell(
    new_uuid(), "1A配电室实时监控系统",
    50, 20, 600, 56, 28, "bold", "#ffffff", "center", "主标题"
))
components.append(make_text_cell(
    new_uuid(), "电力监控大屏 · 设备运行状态总览",
    50, 82, 600, 30, 14, "normal", "#90a4be", "center", "副标题"
))

# === 配电柜卡片区 ===
card_y = 140
cabinet_config = [
    ("1A1_U11柜", "进线柜 · 60台设备", "#4fc3f7"),
    ("1A3_U11柜", "馈线柜 · 12台设备", "#81c784"),
    ("UPS柜",    "电源柜 · 4台设备",  "#ffb74d"),
]

for i, (cab_name, cab_desc, cab_color) in enumerate(cabinet_config):
    cx = 50 + i * 420

    # 柜名标签
    components.append(make_text_cell(
        new_uuid(), cab_name,
        cx, card_y, 380, 32, 20, "bold", cab_color, "left", cab_name
    ))
    # 柜描述
    components.append(make_text_cell(
        new_uuid(), cab_desc,
        cx, card_y + 36, 380, 24, 13, "normal", "#78909c", "left", f"{cab_name}描述"
    ))
    # 柜卡片背景
    components.append(make_rect_cell(
        new_uuid(), f"{cab_name}卡片",
        cx, card_y + 68, 380, 120, "#0d2137", "#1a3a5c", 1, -3
    ))
    # 柜卡片内状态文字
    components.append(make_text_cell(
        new_uuid(), "在线 · 运行中",
        cx + 15, card_y + 82, 350, 28, 14, "normal", "#4caf50", "left", f"{cab_name}状态"
    ))
    # 柜卡片内数据提示
    components.append(make_text_cell(
        new_uuid(), "点击查看详细数据 →",
        cx + 15, card_y + 118, 350, 24, 12, "normal", "#607d8b", "left", f"{cab_name}入口"
    ))

# === 底部信息区 ===
bottom_y = 480
components.append(make_text_cell(
    new_uuid(), "系统状态: 正常监控中",
    50, bottom_y, 400, 28, 13, "normal", "#66bb6a", "left", "运行状态"
))
components.append(make_text_cell(
    new_uuid(), "数据刷新周期: 5秒 · 最新采集: --",
    50, bottom_y + 28, 400, 24, 12, "normal", "#546e7a", "left", "刷新周期"
))

display_model = {
    "uuid": display_model_uuid,
    "name": "1A配电室监控大屏",
    "project_uuid": project_uuid,
    "description": "1A配电室电力监控主屏",
    "displayUid": display_uid,
    "DisplayImage": "",
    "DisplayUserList": "",
    "DisplayType": 1
}

display_layer = {
    "modelId": display_model['displayUid'],
    "pageName": "main",
    "pageId": page_id,
    "isHome": 1,
    "isLogin": 0,
    "pageType": 1,
    "layer": {
        "backColor": "#f0f2f5",
        "backgroundImage": "",
        "widthHeightRatio": "16:9",
        "width": 1920,
        "height": 1080
    },
    "components": {"cells": components}
}

# ==================== 7. 组装完整项目数据包 ====================
project_package = {
    "project": {
        "uuid": project_uuid,
        "name": project_name,
        "description": "1A配电室 ModbusTCP 电力监控项目",
        "createTime": datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
        "protocol": "ModbusTCP",
        "gateway": {
            "ip": "172.31.4.14",
            "port": 502,
            "backupIp": "172.20.255.14",
            "connectType": "TCP",
            "dataFormat": "CDAB"
        }
    },
    "deviceModels": device_models,
    "registerGroups": register_groups,
    "registerPoints": register_points,
    "monitorTree": monitor_list,
    "alarmTriggers": alarm_triggers,
    "displayModel": display_model,
    "displayLayer": display_layer,
    "statistics": {
        "totalDevices": len(devices),
        "deviceModelsCount": len(device_models),
        "registerGroupsCount": len(register_groups),
        "registerPointsCount": len(register_points),
        "alarmTriggersCount": len(alarm_triggers),
        "cabinetsCount": len(cabinets),
        "zonesCount": len([n for n in monitor_list if n['type'] == 0]),
        "floorCount": 1
    }
}

# 保存完整项目数据包
output_path = Path('/Users/yunfanleo/cursorProjects/ISM源码/liu-chang-1A-dev/1A_complete_project_package.json')
with open(output_path, 'w', encoding='utf-8') as f:
    json.dump(project_package, f, ensure_ascii=False, indent=2)

print(f"✅ 完整项目数据包已生成: {output_path}")
print(f"   文件大小: {output_path.stat().st_size / 1024:.1f} KB")
print(f"\n📊 项目统计:")
print(f"   设备模型: {len(device_models)} 种")
print(f"   寄存器组: {len(register_groups)} 个")
print(f"   寄存器数据点: {len(register_points)} 个")
print(f"   设备实例: {len(devices)} 个")
print(f"   告警触发器: {len(alarm_triggers)} 个")
print(f"   配电柜: {len(cabinets)} 个")
print(f"   组态组件: {len(components)} 个")
