# S3 · dcim-capacity / NSPC（Mock 加厚 · 非闭环交付）

> **诚实口径**：JSON 种子 + 脚本演示单机房 **空间(S)+供电(P)水位** 与 **一笔供电预占**；无组态 NSPC 页、无正式表迁移。**未实现不得记为已交付。**

## 与契约对齐

| 项 | Mock | 正式版（未做） |
|----|------|----------------|
| `dcim-capacity` | `xunzhi-jf_dcim_capacity_mock.json` | `dcim_capacity_snapshot` / `dcim_capacity_reservation` |
| room_id | `88888888888888888888888888888888`（`F3-模块机房-01`） | MonitorList 节点 |

种子默认：**S** 42U 总量 / 8 已用 / 2 预占 → **free 32U**；**P** 120kW / 45.5 已用 / 10 预占 → **free 64.5kW**。

## 路径

| 文件 | 说明 |
|------|------|
| `projects-import/xunzhi-jf/xunzhi-jf_dcim_capacity_mock.json` | 种子 |
| `scripts/xunzhi_jf_capacity_store.py` | init / snapshot / reserve / **demo** |
| `scripts/xunzhi_jf_capacity_demo.sh` | 一步自证 |
| `data/xunzhi-jf/dcim_capacity_store.json` | 运行态（gitignore） |

## 自证步骤（一步）

```bash
bash scripts/xunzhi_jf_capacity_demo.sh
```

### 期望输出（关键行）

```
ASSERT OK: 空间(S) total=42 used=8 reserved=2 free=32.0U
ASSERT OK: 供电(P) total=120.0 used=45.5 reserved=10.0 free=64.5kW
=== reserve 3.5 kW on dim P (一笔预占) ===
ASSERT OK: 空间(S) unchanged; 供电(P) reserved +3.5kW; 1 active reservation
DEMO PASS · Store: .../data/xunzhi-jf/dcim_capacity_store.json
```

分步（可选）：

```bash
python3 scripts/xunzhi_jf_capacity_store.py init
python3 scripts/xunzhi_jf_capacity_store.py snapshot
python3 scripts/xunzhi_jf_capacity_store.py reserve --dim P --amount 3.5 --rack-id 88888888888888888888888888888888
```

## 与 S2

- S2：U 位占用（dcim-asset）。
- S3：机房级水位；预占冲突返回 `409`（`reserve` 超额时）。
