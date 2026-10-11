# S3 骨架 · dcim-capacity / NSPC（契约 v2 · Mock 未宣称闭环交付）

> **诚实口径**：本迭代仅 **JSON 种子 + 本地脚本** 演示「单机房空间/供电水位 + 一笔预占」；**无** 组态页、**无** 正式 `dcim_capacity_*` 表与 HTTP API。**未实现部分不得记为已交付。**

## 与契约对齐（规划）

| 契约模块 | 本仓 Mock | 正式版（未做） |
|----------|-----------|----------------|
| `dcim-capacity` | `xunzhi-jf_dcim_capacity_mock.json` | `dcim_capacity_snapshot` / `dcim_capacity_reservation` |
| 维度 N/S/P/C | snapshot.dims 四水位 | 读 MonitorList + RealData 功率测点 |
| 预占 | `reserve` → `409` 冲突 | 同步 API |

**room_id** 与 MonitorList **`F3-模块机房-01`** uuid `88888888888888888888888888888888` 一致（与 S2 `rack_id` 同源）。

## 路径

| 文件 | 说明 |
|------|------|
| `projects-import/xunzhi-jf/xunzhi-jf_dcim_capacity_mock.json` | 种子：N/S/P/C 水位 |
| `data/xunzhi-jf/dcim_capacity_store.json` | 运行态（gitignore） |
| `scripts/xunzhi_jf_capacity_store.py` | init / snapshot / reserve / demo |
| `scripts/xunzhi_jf_capacity_demo.sh` | 一步抽样 |

## 一步命令

```bash
bash scripts/xunzhi_jf_capacity_demo.sh
```

期望：`DEMO PASS`，供电维 `P` 的 `reserved` 增加 3.5 kW，且仅 1 条 `active` 预占记录。

## 与 S2 关系

- **S2 dcim-asset**：U 位占用真相（上架/下架）。
- **S3 dcim-capacity**：机房级水位与预占；正式版读 asset + 树，**禁止**采集侧反写容量。
