# S4 骨架 · dcim-workorder（Mock · 非完整工单产品）

> **诚实口径**：仅演示 **告警 → 建单 → 处理 → 归档** 最小链（JSON + 脚本）；**无** SLA 引擎、**无** 知识库、**无** 前端工单中心。Beego 仅占位查询/建单入口。

## 与契约 / 告警

| 字段 | Mock |
|------|------|
| `alarm_id` | `alm-xz-leak-01`（对应 W1 漏水触发器语义） |
| `dcim_workorder` | `wo_id`, `alarm_id`, `state`, `assignee`, `archived_at` |
| `dcim_wo_event` | `events[]` 流转记录 |

状态机（简化）：`active` 告警 → 建单 `open` → `in_progress` → `archived`；告警 `cleared`。

## 路径

| 文件 | 说明 |
|------|------|
| `projects-import/xunzhi-jf/xunzhi-jf_dcim_workorder_mock.json` | 种子告警 |
| `scripts/xunzhi_jf_workorder_store.py` | create-from-alarm / transit / archive / **demo** |
| `scripts/xunzhi_jf_workorder_demo.sh` | 一步自证 |
| `data/xunzhi-jf/dcim_workorder_store.json` | 运行态（gitignore） |

## HTTP 占位（可选，需 init store）

| 方法 | 路径（加 `/api` 前缀经代理） |
|------|------------------------------|
| GET | `/xunzhi-jf/workorder/:woId` |
| POST | `/xunzhi-jf/workorder/from-alarm` → **501** 提示用脚本建单 |

实现：`controllers/xunzhiJfWorkorderCtl.go`。

## 自证（一步）

```bash
bash scripts/xunzhi_jf_workorder_demo.sh
```

### 期望输出（关键行）

```
ASSERT OK: alarm alm-xz-leak-01 active (漏水报警)
ASSERT OK: created wo-xxxxxxxxxx state=open
ASSERT OK: workorder in_progress
ASSERT OK: workorder archived, alarm cleared
DEMO PASS · Store: .../data/xunzhi-jf/dcim_workorder_store.json
```

**禁止**：宣称已交付完整 dcim-workorder 产品或替代 ISM 原生告警处置全流程。
