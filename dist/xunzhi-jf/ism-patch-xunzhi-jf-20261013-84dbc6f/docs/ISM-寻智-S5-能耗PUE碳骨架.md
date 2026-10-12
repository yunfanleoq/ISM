# S5 骨架 · dcim-energy（Mock · 非自动全量能耗产品）

> **诚实口径**：仅 **单机房样例功耗 + 手工补录一条 + PUE 重算 + 电网碳因子样例**（JSON/脚本）。**无** 分项计量自动汇总、**无** 全站能耗大屏、**无** `dcim_pue_manual_entry` 表迁移。正式版须读 W2 功率 RealData + 补录 API。

## 契约对齐（规划）

| 表/能力 | Mock | 正式（未做） |
|---------|------|----------------|
| `dcim_pue_manual_entry` | `manual_entries[]` | GORM/SQLite |
| `dcim_carbon_factor` | `carbon_factors[]` | 配置中心 |
| PUE | `facility_kw / it_kw` | 分层级用电 |

**room_id**：`88888888888888888888888888888888`（`F3-模块机房-01`）。

## 路径

| 文件 | 说明 |
|------|------|
| `projects-import/xunzhi-jf/xunzhi-jf_dcim_energy_mock.json` | 种子（IT 50kW / 设施 75kW → PUE 1.5） |
| `scripts/xunzhi_jf_energy_store.py` | init / manual / pue / **demo** |
| `scripts/xunzhi_jf_energy_demo.sh` | 一步自证 |
| `data/xunzhi-jf/dcim_energy_store.json` | 运行态（gitignore） |

## 自证

```bash
bash scripts/xunzhi_jf_energy_demo.sh
```

### 期望输出（关键行）

```
ASSERT OK: baseline PUE=1.5 facility=75.0kW carbon_sample=...kg (...kWh)
=== manual entry +5 kW facility (手工补录) ===
ASSERT OK: after manual PUE=1.6 facility=80.0kW carbon_sample=...kg
DEMO PASS · Store: .../data/xunzhi-jf/dcim_energy_store.json
```

查询：`python3 scripts/xunzhi_jf_energy_store.py pue`
