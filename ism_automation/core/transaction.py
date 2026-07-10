"""
ISM Automation — 逻辑事务与回滚系统

即使通过 API 调用，也需要在脚本层面记录所有操作，失败时按逆序回滚。
"""
from __future__ import annotations

import traceback
from dataclasses import dataclass, field
from typing import Callable, List, Any, Optional
from enum import Enum

from ism_automation.core.logger import ism_logger


class OpStatus(Enum):
    PENDING = "pending"
    SUCCESS = "success"
    FAILED = "failed"
    ROLLED_BACK = "rolled_back"


@dataclass
class Operation:
    """单个操作单元，包含执行和回滚函数"""
    name: str
    execute: Callable[[], Any]
    rollback: Callable[[Any], Any]
    result: Any = None
    status: OpStatus = OpStatus.PENDING
    error: Optional[str] = None


@dataclass
class TransactionResult:
    """事务执行结果"""
    success: bool
    completed_ops: List[Operation]
    failed_op: Optional[Operation] = None
    rollback_errors: List[str] = field(default_factory=list)


class Transaction:
    """
    逻辑事务：记录所有操作，失败时按逆序回滚。

    使用示例:
        tx = Transaction()
        tx.add(Operation(
            name="创建A20数据模型",
            execute=lambda: api.model.create("A20电力仪表", ...),
            rollback=lambda result: api.model.delete(result.uuid)
        ))
        tx.add(Operation(
            name="添加1A1设备",
            execute=lambda: api.device.add("1A1_U11_S18_1", ...),
            rollback=lambda result: api.device.delete(result.uuid)
        ))
        result = tx.commit()
    """

    def __init__(self, dry_run: bool = False):
        self.operations: List[Operation] = []
        self.completed: List[Operation] = []
        self.dry_run = dry_run

    def add(self, op: Operation) -> "Transaction":
        """添加操作到事务链，支持链式调用"""
        self.operations.append(op)
        return self

    def commit(self) -> TransactionResult:
        """执行所有操作，失败时回滚"""
        ism_logger.info(f"🚀 事务开始 ({len(self.operations)} 个操作, dry_run={self.dry_run})")

        for i, op in enumerate(self.operations, 1):
            ism_logger.info(f"  [{i}/{len(self.operations)}] {op.name} ...")

            if self.dry_run:
                # 预览模式：只记录，不执行
                op.status = OpStatus.SUCCESS
                op.result = {"dry_run": True, "name": op.name}
                self.completed.append(op)
                ism_logger.info(f"     ✅ [DRY-RUN] {op.name}")
                continue

            try:
                op.result = op.execute()
                op.status = OpStatus.SUCCESS
                self.completed.append(op)
                ism_logger.info(f"     ✅ {op.name}")
            except Exception as e:
                op.status = OpStatus.FAILED
                op.error = str(e)
                ism_logger.error(f"     ❌ {op.name}: {e}")
                ism_logger.error(traceback.format_exc())

                # 回滚已完成的操作
                rollback_errors = self._rollback()
                return TransactionResult(
                    success=False,
                    completed_ops=self.completed,
                    failed_op=op,
                    rollback_errors=rollback_errors,
                )

        ism_logger.info(f"🎉 事务成功完成 ({len(self.completed)} 个操作)")
        return TransactionResult(
            success=True,
            completed_ops=self.completed,
        )

    def _rollback(self) -> List[str]:
        """按逆序回滚已完成的操作，返回回滚错误列表"""
        errors = []
        ism_logger.warning(f"🔄 开始回滚 ({len(self.completed)} 个操作)...")

        for op in reversed(self.completed):
            if self.dry_run:
                # 预览模式下回滚也只是记录
                ism_logger.info(f"     🔄 [DRY-RUN] 回滚 {op.name}")
                op.status = OpStatus.ROLLED_BACK
                continue

            try:
                op.rollback(op.result)
                op.status = OpStatus.ROLLED_BACK
                ism_logger.info(f"     🔄 回滚 {op.name}")
            except Exception as e:
                err_msg = f"⚠️ 回滚失败 {op.name}: {e}"
                errors.append(err_msg)
                ism_logger.error(err_msg)

        return errors

    def preview(self) -> str:
        """生成事务预览报告（dry-run 用）"""
        lines = ["📋 事务预览报告", "=" * 50]
        for i, op in enumerate(self.operations, 1):
            lines.append(f"  {i}. {op.name}")
        lines.append(f"\n共计 {len(self.operations)} 个操作")
        return "\n".join(lines)
