#!/usr/bin/env python3
"""
ISM Automation — CLI 入口
项目导入命令：从 Excel 配置表一键导入完整项目（数据模型 + 设备 + 大屏）

Usage:
    python -m ism_automation.cli.import_project --project 航信机房 --dry-run
    python -m ism_automation.cli.import_project --project 航信机房 --steps model,device
"""
import argparse
import sys
from pathlib import Path

# 添加项目根目录到路径
project_root = Path(__file__).parent.parent.parent
sys.path.insert(0, str(project_root))

from ism_automation.core.logger import ism_logger
from ism_automation.core.transaction import Transaction, Operation
from ism_automation.config.loader import load_project_config
from ism_automation.api.client import get_client
from ism_automation.extractors.excel_parser import ExcelParser
from ism_automation.extractors.hierarchy_builder import HierarchyBuilder
from ism_automation.extractors.layout_engine import LayoutEngine, CanvasConfig, LayoutConfig
from ism_automation.generators.model_generator import ModelGenerator
from ism_automation.generators.device_generator import DeviceGenerator
from ism_automation.generators.dashboard_generator import DashboardGenerator


def preview_report(config, parsed_data, model_configs, device_configs, dashboard_model) -> str:
    """生成预览报告"""
    template = parsed_data.get("template", {})
    total_devices = sum(len(m["devices"]) for m in template.values())
    total_ai = sum(len(m["ai_points"]) for m in template.values())
    total_di = sum(len(m["di_points"]) for m in template.values())
    total_pages = len(dashboard_model.get("pages", []))

    lines = [
        "=" * 60,
        "📋 项目导入预览报告",
        "=" * 60,
        f"",
        f"📁 项目: {config.project_name} ({config.project_uuid})",
        f"📖 Excel: {config.excel_path}",
        f"",
        f"📊 数据模型:",
        f"   创建 {len(model_configs)} 个模型",
        f"   AI 数据点: {total_ai}",
        f"   DI 数据点: {total_di}",
        f"",
        f"📦 设备:",
        f"   新增 {len(device_configs)} 台设备",
        f"",
        f"🎨 组态大屏:",
        f"   创建 {total_pages} 个页面",
        f"   模型名称: {dashboard_model.get('name', '')}",
        f"",
        f"💾 预计操作:",
    ]

    for mc in model_configs:
        for rg in mc.get("register_groups", []):
            lines.append(f"   INSERT devices_model + register_group + {len(rg.get('registers', []))} registers")
    lines.append(f"   INSERT monitor_list × {len(device_configs)}")
    lines.append(f"   INSERT display_model + display_model_page × {total_pages}")
    lines.append(f"   INSERT display_model_layer (大量组件)")
    lines.append(f"")
    lines.append(f"=" * 60)

    return "\n".join(lines)


def run_import(args: argparse.Namespace) -> int:
    """执行项目导入"""
    # 1. 加载配置
    try:
        config = load_project_config(args.project)
    except FileNotFoundError as e:
        ism_logger.error(f"❌ {e}")
        return 1

    ism_logger.info(f"🚀 开始导入项目: {config.project_name}")
    if args.dry_run:
        ism_logger.info("🔍 预览模式 (dry-run): 不会修改任何数据")

    # 2. 解析 Excel
    try:
        parser = ExcelParser(config)
        parsed_data = parser.parse_all()
    except Exception as e:
        ism_logger.error(f"❌ Excel 解析失败: {e}")
        return 1

    # 3. 生成配置
    model_gen = ModelGenerator(config)
    model_configs = model_gen.generate(parsed_data)

    device_list = parsed_data.get("device_list", [])
    # 模型 UUID 映射（后续从 API 返回获取）
    model_mapping = {}  # 在 API 调用后填充

    device_gen = DeviceGenerator(config)
    device_configs = device_gen.generate(device_list, model_mapping)

    # 构建层级树
    # 注意：这里需要从 DB 查询现有设备树，或使用设备清单构建
    # 简化：使用设备列表构建层级
    tree = device_gen.build_tree(device_configs)

    # 生成大屏
    dashboard_gen = DashboardGenerator(config)
    # 这里需要 stats，简化使用空 stats
    stats = {"device_count": len(device_configs), "online_count": 0, "room_count": 1}
    pages = dashboard_gen.generate_pages({"zones": [{"name": "默认", "rooms": [{"name": "默认", "cabinets": []}]}]})
    dashboard_model = dashboard_gen.create_display_model(pages)

    # 4. 预览报告
    report = preview_report(config, parsed_data, model_configs, device_configs, dashboard_model)
    print(report)

    if args.dry_run:
        ism_logger.info("🔍 预览完成，未执行任何操作")
        return 0

    # 5. 确认执行
    if not args.yes:
        confirm = input("\n确认执行? [y/N]: ")
        if confirm.lower() != 'y':
            ism_logger.info("❌ 用户取消操作")
            return 0

    # 6. 执行导入（通过 API）
    try:
        client = get_client(
            base_url=config.api_base_url,
            username=args.username or "admin",
            password=args.password or "123456",
        )
    except Exception as e:
        ism_logger.error(f"❌ 登录失败: {e}")
        return 1

    # 使用事务执行导入
    tx = Transaction(dry_run=args.dry_run)
    created_models = []
    created_devices = []

    # 步骤 1: 创建数据模型
    if "model" in args.steps or "all" in args.steps:
        for mc in model_configs:
            tx.add(Operation(
                name=f"创建模型 {mc['name']}",
                execute=lambda mc=mc: client.model.create_with_registers(
                    protocol=mc["protocol"],
                    name=mc["name"],
                    project_uuid=mc["project_uuid"],
                    register_groups=mc["register_groups"],
                ),
                rollback=lambda result: client.model.delete_modbus(
                    result.get("model", {}).get("uuid", "")
                ) if result else None,
            ))

    # 步骤 2: 创建设备
    if "device" in args.steps or "all" in args.steps:
        # 需要获取已创建模型的 UUID 映射
        for dc in device_configs:
            tx.add(Operation(
                name=f"添加设备 {dc['name']}",
                execute=lambda dc=dc: client.device.add(
                    name=dc["name"],
                    project_uuid=dc["project_uuid"],
                    model_uuid=dc["model_uuid"],
                    protocol_type=dc["protocol_type"],
                    ip=dc["ip"],
                    port=dc["port"],
                    slave_id=dc["slave_id"],
                ),
                rollback=lambda result: client.device.delete(
                    result.get("data", {}).get("sid", 0)
                ) if result else None,
            ))

    # 步骤 3: 创建组态大屏
    if "dashboard" in args.steps or "all" in args.steps:
        tx.add(Operation(
            name="创建组态大屏模型",
            execute=lambda: client.display.create(
                name=dashboard_model["name"],
                project_uuid=dashboard_model["project_uuid"],
                width=dashboard_model["width"],
                height=dashboard_model["height"],
            ),
            rollback=lambda result: client.display.delete(
                result.get("data", {}).get("uuid", "")
            ) if result else None,
        ))

    # 执行事务
    result = tx.commit()
    if result.success:
        ism_logger.info("🎉 项目导入成功!")
        return 0
    else:
        ism_logger.error(f"❌ 导入失败: {result.failed_op}")
        if result.rollback_errors:
            for err in result.rollback_errors:
                ism_logger.error(f"   回滚错误: {err}")
        return 1


def main():
    parser = argparse.ArgumentParser(
        description="ISM 项目导入工具",
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )
    parser.add_argument(
        "--project", "-p",
        required=True,
        help="项目名称（对应 config/projects/ 下的 YAML 配置文件）",
    )
    parser.add_argument(
        "--steps", "-s",
        default="all",
        help="执行步骤: all, model, device, dashboard (逗号分隔)",
    )
    parser.add_argument(
        "--dry-run", "-n",
        action="store_true",
        help="预览模式：只输出报告，不执行实际操作",
    )
    parser.add_argument(
        "--yes", "-y",
        action="store_true",
        help="跳过确认提示，直接执行",
    )
    parser.add_argument(
        "--username", "-u",
        default="admin",
        help="ISM 登录用户名",
    )
    parser.add_argument(
        "--password", "-pw",
        default="123456",
        help="ISM 登录密码",
    )
    parser.add_argument(
        "--verbose", "-v",
        action="store_true",
        help="详细输出模式",
    )

    args = parser.parse_args()
    args.steps = [s.strip() for s in args.steps.split(",")]

    return run_import(args)


if __name__ == "__main__":
    sys.exit(main())
