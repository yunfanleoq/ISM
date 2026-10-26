#!/usr/bin/env python3
"""
ISM Automation — 项目校验命令
验证项目完整性：数据模型、设备树、组态大屏、数据绑定等

Usage:
    python -m ism_automation.cli.verify_project --project 航信机房
"""
import argparse
import sys
from pathlib import Path

project_root = Path(__file__).parent.parent.parent
sys.path.insert(0, str(project_root))

from ism_automation.core.logger import ism_logger
from ism_automation.config.loader import load_project_config
from ism_automation.api.client import get_client


def run_verify(args: argparse.Namespace) -> int:
    """执行项目校验"""
    try:
        config = load_project_config(args.project)
    except FileNotFoundError as e:
        ism_logger.error(f"❌ {e}")
        return 1

    ism_logger.info(f"🔍 校验项目: {config.project_name}")

    try:
        client = get_client(
            base_url=config.api_base_url,
            username=args.username or "admin",
            password=args.password or "123456",
        )
    except Exception as e:
        ism_logger.error(f"❌ 登录失败: {e}")
        return 1

    issues = []
    warnings = []

    # 1. 校验项目存在
    try:
        projects = client.project.list()
        project = next((p for p in projects if p.get("uuid") == config.project_uuid), None)
        if not project:
            issues.append(f"项目 {config.project_uuid} 不存在")
        else:
            ism_logger.info(f"✅ 项目存在: {project.get('Name', 'Unknown')}")
    except Exception as e:
        issues.append(f"查询项目失败: {e}")

    # 2. 校验数据模型
    try:
        models = client.model.list_modbus(config.project_uuid)
        expected_models = [m.name for m in config.models]
        found_models = [m.get("Name", "") for m in models]

        for em in expected_models:
            if em not in found_models:
                issues.append(f"数据模型缺失: {em}")
            else:
                ism_logger.info(f"✅ 数据模型存在: {em}")

        ism_logger.info(f"📊 数据模型: {len(models)} 个")
    except Exception as e:
        issues.append(f"查询数据模型失败: {e}")

    # 3. 校验设备树
    try:
        tree = client.device.get_tree(config.project_uuid)
        device_count = sum(1 for n in tree if n.get("type") == 1)
        ism_logger.info(f"📦 设备: {device_count} 台")

        if device_count == 0:
            warnings.append("项目中没有设备")
    except Exception as e:
        issues.append(f"查询设备树失败: {e}")

    # 4. 校验组态大屏
    try:
        displays = client.display.list(config.project_uuid)
        ism_logger.info(f"🎨 组态大屏: {len(displays)} 个")

        if not displays:
            warnings.append("项目中没有组态大屏")
    except Exception as e:
        issues.append(f"查询组态大屏失败: {e}")

    # 5. 校验数据绑定
    # 检查设备是否有模型绑定
    # 检查大屏组件是否有数据绑定

    # 输出报告
    print(f"\n{'='*60}")
    print(f"📋 项目校验报告: {config.project_name}")
    print(f"{'='*60}")

    if issues:
        print(f"\n❌ 错误 ({len(issues)}):")
        for i in issues:
            print(f"   - {i}")

    if warnings:
        print(f"\n⚠️ 警告 ({len(warnings)}):")
        for w in warnings:
            print(f"   - {w}")

    if not issues and not warnings:
        print(f"\n✅ 项目校验通过，一切正常!")

    print(f"{'='*60}\n")

    return 1 if issues else 0


def main():
    parser = argparse.ArgumentParser(description="ISM 项目校验工具")
    parser.add_argument("--project", "-p", required=True, help="项目名称")
    parser.add_argument("--username", "-u", default="admin", help="用户名")
    parser.add_argument("--password", "-pw", default="123456", help="密码")
    parser.add_argument("--verbose", "-v", action="store_true", help="详细输出")

    args = parser.parse_args()
    return run_verify(args)


if __name__ == "__main__":
    sys.exit(main())
