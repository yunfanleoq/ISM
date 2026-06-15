#!/usr/bin/env python3
"""
ISM device_real_data 修复工具 v2.0 — 通用版本

=== UUID 关系说明 ===
  monitor_list.uuid                → device_real_data.device_uuid      ★ 最关键
  modbus_devices_data_model.uuid   → device_real_data.model_data_uuid  ★ 最关键
  devices_model.uuid               → device_real_data.muid (冗余字段)
  devices_model.uuid               → modbus_devices_data_model.muid
  devices_model.uuid               → modbus_devices_register_group.muid

=== 使用场景 ===
  1. 用直接 SQL INSERT 创建 monitor_list 后, device_real_data 缺失 → 补建
  2. 导入脚本的 monitorAdd API 失败后用 SQL 旁路 → 补建
  3. 跨项目迁移后 UUID 断裂 → 重建关联

=== 支持的协议类型 ===
  deviceType=1: SNMP    deviceType=2: Modbus (默认)
  deviceType=3: OPCUA   deviceType=5: RESTFul
  deviceType=15: SimS7  deviceType=20: MQTT
  deviceType=30: DLT645 deviceType=40: IEC104
  deviceType=350: IEC61850  deviceType=470: HJ212
  deviceType=480: VirtualDevice  deviceType=490: CJT188
  deviceType=500: BACnet
"""
import sqlite3, uuid, argparse, sys, os

# 协议类型 → (数据模型表名, 额外字段映射)
PROTO_TABLES = {
    1:  "snmp_devices_data_model",
    2:  "modbus_devices_data_model",
    3:  "opcua_devices_data_model",
    5:  "restful_data_model",
    15: "sim_s7_data_model",
    20: "mqtt_devices_data_model",
    30: "dlt645_devices_data_model",
    40: "iec104_devices_data_model",
    350: "iec61850_devices_data_model",
    470: "hj212_devices_data_model",
    480: "virtual_device_data_model",
    490: "cjt188_devices_data_model",
    500: "bacnet_devices_data_model",
}

def build_real_data(db_path, project_uuid, device_type=2, dry_run=False):
    conn = sqlite3.connect(db_path)
    conn.row_factory = sqlite3.Row

    tbl = PROTO_TABLES.get(device_type)
    if not tbl:
        print(f"不支持的 device_type={device_type}")
        conn.close()
        return

    # 1. 查所有数据模型
    models = conn.execute(
        "SELECT uuid, name FROM devices_model WHERE project_uuid=? AND type=?",
        (project_uuid, device_type)
    ).fetchall()

    # 2. 为每个模型加载数据点定义
    model_points = {}
    for m in models:
        muid = m['uuid']
        groups = conn.execute(
            "SELECT uuid FROM modbus_devices_register_group WHERE muid=?", (muid,)
        ).fetchall()
        pts = []
        for (g_uuid,) in groups:
            rows = conn.execute(f"""
                SELECT uuid, name, register_address, conversion_expression,
                       is_alarm, alarm_level, alarm_message, alarm_clear_message,
                       data_unit, is_record, record_type, record_interval, record_data_charge
                FROM {tbl}
                WHERE register_group_uuid=?
                ORDER BY register_address
            """, (g_uuid,)).fetchall()
            for r in rows:
                pts.append(dict(r))
        model_points[muid] = pts
        print(f"  模型 {m['name']} ({muid[:12]}): {len(pts)} 数据点")

    if not model_points:
        print("没有找到任何数据模型!")
        conn.close()
        return

    # 3. 查所有设备
    devices = conn.execute(
        "SELECT uuid, name, muid FROM monitor_list WHERE project_uuid=? AND type=1 AND deleted_at IS NULL",
        (project_uuid,)
    ).fetchall()
    print(f"\n  共 {len(devices)} 个设备")

    # 4. 删旧 + 建新
    old = conn.execute(
        "SELECT count(*) FROM device_real_data WHERE project_uuid=?",
        (project_uuid,)
    ).fetchone()[0]
    if old:
        if dry_run:
            print(f"  [干跑] 将删除 {old} 条旧数据")
        else:
            conn.execute("DELETE FROM device_real_data WHERE project_uuid=?", (project_uuid,))
            conn.commit()
            print(f"  已删除 {old} 条旧数据")

    if dry_run:
        total = sum(len(model_points.get(d['muid'], [])) for d in devices)
        print(f"  [干跑] 将创建 {total} 条 device_real_data (alarm_shield=1)")
        conn.close()
        return

    # 5. 插入
    sql = f"""
        INSERT INTO device_real_data
        (created_at, updated_at, deleted_at, name, device_name, oid, uuid, project_uuid,
         auth, type, value, muid, model_data_uuid, conversion_expression, device_uuid,
         device_type, data_unit, is_alarm, alarm_level, alarm_message, alarm_clear_message,
         is_record, record_type, record_interval, record_data_charge, alarm_shield)
        VALUES (datetime('now'), datetime('now'), NULL,
         ?, ?, '', ?, ?, 2, 1, '', ?, ?,
         NULL, ?, ?, ?, 0, 0, '', '',
         ?, ?, ?, ?, 1)
    """

    total = 0
    for dev in devices:
        muid = dev['muid']
        if muid not in model_points:
            print(f"  ⚠ 设备 {dev['name']} 的模型 {muid} 无数据点, 跳过")
            continue
        for pt in model_points[muid]:
            conn.execute(sql, (
                pt['name'], dev['name'],
                str(uuid.uuid4()).replace('-', ''), project_uuid,
                muid, pt['uuid'],
                dev['uuid'], device_type, pt['data_unit'] or '',
                pt['is_record'] or 0, pt['record_type'] or 0,
                pt['record_interval'] or 0, pt['record_data_charge'] or '',
            ))
            total += 1

    conn.commit()

    # 验证
    v = conn.execute(
        "SELECT count(*) as cnt, count(DISTINCT device_uuid) as devs FROM device_real_data WHERE project_uuid=?",
        (project_uuid,)
    ).fetchone()
    print(f"\n  创建 {total} 条")
    print(f"  验证: {v['cnt']} 条, {v['devs']} 个设备")
    conn.close()

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='ISM device_real_data 修复工具')
    parser.add_argument('project_uuid', help='目标项目的 UUID')
    parser.add_argument('--db', default=os.path.join(os.path.dirname(__file__), '..',
                        'ism_server_user', 'data', 'db', 'ism.db'),
                        help='SQLite 数据库路径')
    parser.add_argument('--type', type=int, default=2, help='协议类型 (默认 2=Modbus)')
    parser.add_argument('--dry-run', action='store_true', help='仅预览不执行')
    args = parser.parse_args()

    print(f"ISM device_real_data 修复工具 v2.0")
    print(f"  DB:   {args.db}")
    print(f"  项目: {args.project_uuid}")
    print(f"  类型: {args.type} ({PROTO_TABLES.get(args.type, 'Unknown')})")
    print(f"  {'[干跑模式]' if args.dry_run else '[执行模式]'}")
    print()

    build_real_data(args.db, args.project_uuid, args.type, args.dry_run)
