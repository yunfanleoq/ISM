---
name: ism-默认大屏瘦身
overview: 将默认「组态界面」从 17 页/42 组件瘦身为 4 页/26 组件，合并重复实时表页，精简导航，并产出可现场执行的 SQL 补丁。
todos:
  - id: audit
    content: 审计备份 SQL 中 17 页组件分布
    status: completed
  - id: script
    content: 实现 scripts/slim_default_dashboard.py（17→4）
    status: completed
  - id: sql
    content: 生成 releases/sql/slim_default_dashboard.sql
    status: completed
  - id: docs
    content: 同步 docs/ISM-默认大屏瘦身方案.md 与本 plan
    status: completed
  - id: apply-field
    content: 现场执行 SQL/--apply 并验收 AppRun/设计页
    status: pending
---

# ISM 默认大屏瘦身 Plan

详见人类可读正文：[`docs/ISM-默认大屏瘦身方案.md`](../../docs/ISM-默认大屏瘦身方案.md)

## 目标

| 指标 | 前 | 后 |
|------|----|----|
| 页 | 17 | 4 |
| 组件 | 42 | 26 |

结构：主界面 → 配电室总览 → 统一实时数据 → 历史报警。

## 产物

- `scripts/slim_default_dashboard.py`
- `releases/sql/slim_default_dashboard.sql`
- `docs/ISM-默认大屏瘦身方案.md`

## 现场下一步

```bash
python3 scripts/slim_default_dashboard.py --apply
# 或执行 releases/sql/slim_default_dashboard.sql
```

验收：仅 4 页存活；菜单 3 项可跳转；设计页不再因重复表页长时间 Loading。
