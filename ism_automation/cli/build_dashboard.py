#!/usr/bin/env python3
"""
ISM Automation — 大屏构建命令
根据模板自动生成组态大屏

Usage:
    python -m ism_automation.cli.build_dashboard --project 航信机房 --template industrial_4level --dry-run
"""
import argparse
import sys
from pathlib import Path

project_root = Path(__file__).parent.parent.parent
sys.path.insert(0, str(project_root))

from ism_automation.core.logger import ism_logger
from ism_automation.core.transaction import Transaction, Operation
from ism_automation.config.loader import load_project_config, load_template_config
from ism_automation.api.client import get_client
from ism_automation.extractors.hierarchy_builder import HierarchyBuilder
from ism_automation.generators.dashboard_generator import DashboardGenerator


def run_build(args: argparse.Namespace) -> int:
    """执行大屏构建"""
    # 1. 加载配置
    try:
        config = load_project_config(args.project)
    except FileNotFoundError as e:
        ism_logger.error(f"❌ {e}")
        return 1

    # 2. 加载模板
    try:
        template = load_template_config(args.template)
    except FileNotFoundError as e:
        ism_logger.error(f"❌ {e}")
        return 1

    ism_logger.info(f"🚀 构建大屏: {config.project_name} / 模板: {args.template}")
    if args.dry_run:
        ism_logger.info("🔍 预览模式")

    # 3. 查询现有设备层级（从 API 或 DB）
    try:
        client = get_client(
            base_url=config.api_base_url,
            username=args.username or "admin",
            password=args.password or "123456",
        )
        device_tree = client.device.get_tree(config.project_uuid)
        ism_logger.info(f"📊 获取设备树: {len(device_tree)} 个节点")
    except Exception as e:
        ism_logger.warning(f"⚠️ 获取设备树失败: {e}，使用空层级")
        device_tree = []

    # 4. 构建层级
    builder = HierarchyBuilder(device_tree)
    tree = builder.build_tree()

    # 5. 生成大屏
    generator = DashboardGenerator(config)
    pages = generator.generate_pages(tree)
    dashboard = generator.create_display_model(pages, name=args.name)

    # 6. 预览报告
    print(f"\n{'='*60}")
    print(f"📋 大屏构建预览报告")
    print(f"{'='*60}")
    print(f"项目: {config.project_name}")
    print(f"模板: {args.template}")
    print(f"页面数: {len(pages)}")
    print(f"首页: {pages[0]['name'] if pages else '无'}")
    print(f"组件数: {sum(len(p.get('components', [])) for p in pages)}")
    print(f"{'='*60}\n")

    if args.dry_run:
        return 0

    # 7. 确认执行
    if not args.yes:
        confirm = input("确认执行? [y/N]: ")
        if confirm.lower() != 'y':
            return 0

    # 8. 执行创建
    tx = Transaction(dry_run=args.dry_run)
    tx.add(Operation(
        name="创建组态大屏模型",
        execute=lambda: client.display.create(
            name=dashboard["name"],
            project_uuid=dashboard["project_uuid"],
            width=dashboard["width"],
            height=dashboard["height"],
        ),
        rollback=lambda result: client.display.delete(
            result.get("data", {}).get("uuid", "")
        ) if result else None,
    ))

    for page in pages:
        tx.add(Operation(
            name=f"添加页面 {page['name']}",
            execute=lambda p=page: client.display.add_page(
                model_uuid="",  # 从上一个操作结果获取
                name=p["name"],
                width=p["width"],
                height=p["height"],
                is_home=p.get("is_home", False),
            ),
            rollback=lambda result: client.display.delete_page(
                result.get("data", {}).get("uuid", "")
            ) if result else None,
        ))

    result = tx.commit()
    if result.success:
        ism_logger.info("🎉 大屏构建成功!")
        return 0
    else:
        ism_logger.error(f"❌ 构建失败: {result.failed_op}")
        return 1


def main():
    parser = argparse.ArgumentParser(description="ISM 大屏构建工具")
    parser.add_argument("--project", "-p", required=True, help="项目名称")
    parser.add_argument("--template", "-t", default="industrial_4level", help="模板名称")
    parser.add_argument("--name", "-n", default=None, help="大屏名称（默认使用项目名称）")
    parser.add_argument("--dry-run", action="store_true", help="预览模式")
    parser.add_argument("--yes", "-y", action="store_true", help="跳过确认")
    parser.add_argument("--username", "-u", default="admin", help="用户名")
    parser.add_argument("--password", "-pw", default="123456", help="密码")

    args = parser.parse_args()
    return run_build(args)


if __name__ == "__main__":
    sys.exit(main())
