import json
from pathlib import Path

# 读取已生成的数据模型
with open('/Users/yunfanleo/cursorProjects/ISM源码/ism_data_models.json', 'r', encoding='utf-8') as f:
    data = json.load(f)

# 修正A40和UPS的谐波参数（系数和解析方式交换）
for template in data['templates']:
    if template['type'] in ['A40', '施耐德UPS']:
        for ai in template['aiParams']:
            # 检查是否需要交换：如果系数是整数（如73, 177, 179）且解析方式是小数（如0.1, 0.01）
            try:
                factor_val = float(ai['factor']) if ai['factor'] != '' else 0
                parse_val = float(ai['parseMode']) if ai['parseMode'] != '' else 0
                
                # 如果系数大于10且解析方式小于1，则交换
                if factor_val > 10 and parse_val < 1 and parse_val > 0:
                    print(f"[修正] {template['type']} {ai['name']}: 交换系数({ai['factor']})和解析方式({ai['parseMode']})")
                    ai['factor'], ai['parseMode'] = ai['parseMode'], ai['factor']
            except:
                pass

# 重新生成ISM配置
def generate_ism_model(template):
    model = {
        "model_name": template['type'],
        "description": f"{template['type']} 电力设备模型",
        "protocol": "ModbusTCP",
        "register_groups": []
    }
    
    if template['aiParams']:
        ai_group = {
            "group_name": "AI数据",
            "register_type": "HOLDING",
            "start_address": 0,
            "points": []
        }
        for ai in template['aiParams']:
            ai_group["points"].append({
                "name": ai['name'],
                "data_type": "FLOAT32",
                "register_count": 2,
                "factor": ai['factor'],
                "parse_mode": ai['parseMode'],
                "offset": ai['offset']
            })
        model["register_groups"].append(ai_group)
    
    if template['diParams']:
        di_group = {
            "group_name": "DI数据",
            "register_type": "DI",
            "start_address": 0,
            "points": []
        }
        for di in template['diParams']:
            di_group["points"].append({
                "name": di['name'],
                "data_type": "BOOL",
                "register_count": 1,
                "register_offset": di['registerOffset'],
                "bit_offset": di['bitOffset']
            })
        model["register_groups"].append(di_group)
    
    return model

ism_models = {}
for template in data['templates']:
    model = generate_ism_model(template)
    ism_models[template['type']] = model

device_instances = []
for dev in data['devices']:
    device_instances.append({
        "name": dev['name'],
        "template_type": dev['templateType'],
        "ai_start_address": dev['aiStartAddr'],
        "di_start_address": dev['diStartAddr'],
        "modbus_config": {
            "ip": "172.31.4.14",
            "port": 502,
            "slave_id": 1
        }
    })

project_config = {
    "project_name": "1A配电室",
    "description": "1A配电室 ModbusTCP 数据采集项目",
    "protocol": "ModbusTCP",
    "gateway": {
        "ip": "172.31.4.14",
        "port": 502,
        "backup_ip": "172.20.255.14"
    },
    "device_models": ism_models,
    "device_instances": device_instances,
    "statistics": {
        "total_devices": len(device_instances),
        "a20_count": sum(1 for d in device_instances if d['template_type'] == 'A20'),
        "a40_count": sum(1 for d in device_instances if d['template_type'] == 'A40'),
        "ups_count": sum(1 for d in device_instances if d['template_type'] == '施耐德UPS')
    }
}

output_path = Path('/Users/yunfanleo/cursorProjects/ISM源码/ism_project_config.json')
with open(output_path, 'w', encoding='utf-8') as f:
    json.dump(project_config, f, ensure_ascii=False, indent=2)

print(f"✅ ISM项目配置已保存到: {output_path}")
print(f"\n项目统计:")
print(f"  总设备数: {project_config['statistics']['total_devices']}")
print(f"  A20设备: {project_config['statistics']['a20_count']}")
print(f"  A40设备: {project_config['statistics']['a40_count']}")
print(f"  UPS设备: {project_config['statistics']['ups_count']}")
