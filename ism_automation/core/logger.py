"""
ISM Automation — 结构化日志系统
基于 loguru，支持文件轮转、彩色输出、结构化记录
"""
import sys
from pathlib import Path
from loguru import logger

LOG_DIR = Path(__file__).parent.parent / "logs"
LOG_DIR.mkdir(exist_ok=True)


def setup_logger(name: str = "ism_automation", level: str = "INFO"):
    """初始化日志系统，返回 logger 实例"""
    # 移除默认 handler
    logger.remove()

    # 控制台输出（带颜色）
    logger.add(
        sys.stdout,
        level=level,
        format="<green>{time:YYYY-MM-DD HH:mm:ss}</green> | "
               "<level>{level: <8}</level> | "
               "<cyan>{name}</cyan> | "
               "{message}",
        colorize=True,
    )

    # 文件输出（按日期轮转）
    logger.add(
        LOG_DIR / f"{name}.log",
        rotation="10 MB",
        retention="30 days",
        level=level,
        format="{time:YYYY-MM-DD HH:mm:ss.SSS} | {level: <8} | {name} | {message}",
        enqueue=True,
    )

    # 错误日志单独文件
    logger.add(
        LOG_DIR / f"{name}.error.log",
        rotation="10 MB",
        retention="30 days",
        level="ERROR",
        format="{time:YYYY-MM-DD HH:mm:ss.SSS} | {level: <8} | {name} | {message}",
        enqueue=True,
    )

    return logger.bind(name=name)


# 全局 logger 实例
ism_logger = setup_logger()
