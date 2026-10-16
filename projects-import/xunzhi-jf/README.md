# 寻智机房环境 · xunzhi-jf（S0 竖切 + W2/W3 扩展）

Modbus TCP：**温湿度 + 漏水**（W1）+ **列头柜 + UPS**（W2）+ **精密空调 + 氢气**（W3），默认 **`gather_interval_ms=5000`**（≤10s 契约）。

**字节序**：物模型 `DataFormat` / 点表 `ByteOrder` 默认 **ABCD**（与 Mock BE 寄存器原码一致）。**勿用 CDAB** — Unsigned short 会按 LittleEndian 误读，温湿度再 `{val}*0.1` 后表内常见 **~5xxx**。

链路：**协议 → DeviceRealData → WebSocket（`conf/app.conf` 的 `wsport`，本仓 dev 常见 10225）→ 组态**。

## 测试复测入口（推荐）

```bash
./scripts/xunzhi_jf_smoke_env.sh
./scripts/xunzhi_jf_collect_evidence.sh
bash scripts/xunzhi_jf_alarm_demo.sh
```

告警：导入包含 **`alarmTriggers`**（漏水 `=`1、温度 `>`35）及告警点 **`alarmOnValue: 1`**。项目已存在时用 `scripts/xunzhi_jf_sync_alarm_db.py`，勿仅重复 ImportProject（会 4001）。

证据目录：`docs/ISM-寻智-W1复测证据/` · 试打样包：`dist/xunzhi-jf/ism-patch-xunzhi-jf-*` · **演示索引**：`docs/ISM-寻智-演示脚本索引.md`

### WebSocket 端口

| 说明 | 端口 |
|------|------|
| 契约/前端默认示例 | 10215 |
| **本仓 `ism_server_user/conf/app.conf`** | **10225（联调以实测为准）** |

业主本机冒烟：`:10225` 可连 RealData，`:10215` 拒连。文档已对齐，**未改现网配置**。

## 资产路径

| 文件 | 说明 |
|------|------|
| `xunzhi-jf_complete_project_package.json` | `POST /ImportProject`（含 W2 设备时需重新导入或新库） |
| `xunzhi-jf_device_points.json` | W1 点表 |
| `xunzhi-jf_w2_device_points.json` | W2 配电/UPS 附录 |
| `xunzhi-jf_w3_device_points.json` | W3 空调/氢气附录 |
| `generate_xunzhi_jf_package.py` | 生成 W1 基线包 |
| `append_w2_power_ups.py` | 向包追加 PWR/UPS（slave 3/4） |
| `append_w2_power_viz.py` | W2：平面实时表 + 配电/UPS 表盘 |
| `append_w3_ac_h2.py` | 向包追加 AC/H2（slave 5/6）· 探针 **P-COL-03** `cycle_log_ac_h2` |
| `append_s1_envmon_demo.py` | S1：组态告警列表 + 第二机房树节点 |
| `append_s1_chart_thicken.py` | S1：平面页挂 ≥20 类图表/仪表 shape（**演示平面勿依赖**） |
| `reflow_xunzhi_jf_env_plane_layout.py` | **演示平面收口**：实时表 + 机房树 + 告警（去图表/箭头堆叠） |
| `append_s1_gauge_plus.py` | S1：补足仪表盘实例 ≥10 |
| `xunzhi-jf_dcim_capacity_mock.json` | S3：dcim-capacity Mock 种子（骨架，非闭环） |
| `xunzhi-jf_dcim_workorder_mock.json` | S4：dcim-workorder Mock 种子（骨架） |
| `xunzhi-jf_dcim_energy_mock.json` | S5：dcim-energy PUE/碳 Mock 种子 |
| `xunzhi-jf_security_event_mock.json` | S6：安防/消防事件队列 Mock 种子 |
| `xunzhi-jf_dcim_asset_mock.json` | S2：dcim-asset Mock 种子（rack_id = MonitorList 节点） |

### 生成完整导入包（推荐顺序）

```bash
python3 projects-import/xunzhi-jf/generate_xunzhi_jf_package.py
python3 projects-import/xunzhi-jf/append_w2_power_ups.py
python3 projects-import/xunzhi-jf/append_w2_power_viz.py
python3 projects-import/xunzhi-jf/append_w3_ac_h2.py
python3 projects-import/xunzhi-jf/append_s1_envmon_demo.py
python3 projects-import/xunzhi-jf/append_s1_chart_thicken.py
python3 projects-import/xunzhi-jf/append_s1_gauge_plus.py
python3 projects-import/xunzhi-jf/append_w4_battery_sample.py
python3 projects-import/xunzhi-jf/append_s6_25d_nav_shell.py
python3 projects-import/xunzhi-jf/append_s6_chart_catalog_push.py
python3 projects-import/xunzhi-jf/reflow_xunzhi_jf_env_plane_layout.py
```

S2 资产 U 位 Mock：`bash scripts/xunzhi_jf_asset_u_demo.sh`（见 `docs/ISM-寻智-S2-资产U位骨架.md`）。  
S3 NSPC 容量 Mock（骨架）：`bash scripts/xunzhi_jf_capacity_demo.sh`（见 `docs/ISM-寻智-S3-NSPC容量骨架.md`）。  
S4 工单 Mock（骨架）：`bash scripts/xunzhi_jf_workorder_demo.sh`（见 `docs/ISM-寻智-S4-工单骨架.md`）。  
S5 能耗 PUE/碳 Mock：`bash scripts/xunzhi_jf_energy_demo.sh`（见 `docs/ISM-寻智-S5-能耗PUE碳骨架.md`）。  
S6 安防/消防事件桩：`bash scripts/xunzhi_jf_security_demo.sh`（见 `docs/ISM-寻智-S6-安防消防事件桩.md`）。  
2.5D 导航壳：`docs/ISM-寻智-S6-2.5D导航壳.md` · 电池抽样：`bash scripts/xunzhi_jf_battery_demo.sh` · 双机演练：`docs/ISM-寻智-双机软件就绪演练.md`。

Mock：`scripts/xunzhi_jf_modbus_mock.py`（slave 1–7）。周期日志：`P-COL-01` 温湿度/漏水 · `P-COL-02` PWR/UPS · `P-COL-03` AC/H2。

组态：**寻智机房环境平面监视-xunzhi-jf** / 页 **`xunzhi-jf-env-plane`**

## 账号与端口

| 服务 | 默认 |
|------|------|
| HTTP API | `http://127.0.0.1:8081` |
| 前端 dev | `http://127.0.0.1:7080` |
| Mock Modbus | `127.0.0.1:1503`（双口树默认；`XUNZHI_MOCK_PORT` 可覆盖） |
| 登录 | `admin` / `123456`（前端 MD5 后传 `e10adc3949ba59abbe56e057f20f883e`） |

## 禁止

不覆盖麒麟生产、`patches/ism-*-kylin`、`releases/` 现网包。
