# 寻智 xunzhi-jf · 演示脚本索引（PR#5 收口）

## 明早演示冻结（仓内收口 · 未宣称正式验收）

| 项 | 值 |
|----|-----|
| **Git tip** | `125e9e91` |
| **试打样包** | `dist/xunzhi-jf/ism-patch-xunzhi-jf-20261012-e8bb7f8/` |
| **分支** | `feature/xunzhi-jf-env-monitor` · PR [#5](https://github.com/yunfanleoq/ISM/pull/5) |
| **口径** | 脚本/Mock/文档在仓内冻结；**ism_server / 7080 / 本机 DB 运行态由演示机自启**；不覆盖 `patches/*kylin`、`releases/` 现网 |

图表挂点计数（导入包 `_meta`）：**shape 52** · **仪表盘实例 28**（见 `docs/ISM-寻智-S1-图表仪表盘清单草稿.md`）。

---

## 唯一演示根目录（必读）

- 全链路只认 **一个** 仓库根 `XUNZHI_ISM_ROOT`（默认=脚本上级目录；业主常见 `~/cursorProjects/ISM源码`）。
- **禁止** `ISM源码` 与 `ISM源码-xunzhi-jf` 各起 `ism_server` 混用 sqlite。
- 起服顺序（告警复测）：**先** `python3 scripts/xunzhi_jf_modbus_mock.py --port 1502`（保持 LISTEN）→ **再** `cd ism_server_user && ./ism_server`（本树）
- 自检：`python3 scripts/xunzhi_jf_sync_alarm_db.py --print-env --strict-mock`（`ism_server` cwd、`ism.db`、**Mock:1502**、LEAK `realdata.updated_at`）

**StartupAlarmDelay**：`sync` / `alarm_demo` 会为寻智项目写入 `DelayMinutes:0`；若 `ism_server` 在写入前已启动，须 **保持 Mock:1502 后重启 ism_server**（或 `export XUNZHI_WAIT_STARTUP_DELAY=1` 等旧 10min 窗结束）。  
**采集判定**：以 LEAK `device_real_data` 刷新 / 日志 `cycle_log_th_leak` 为准；`monitor_list.status=2` 可滞后。  
**告警 sync**：`AlarmTriggerEdit` 热加载；`alarm_demo` 用 `xunzhi_jf_mock_leak_ctl.sh off|on` **热切换** 0→1（禁止 pkill Mock）。  
**ism 检测**：以 **`ism_server_running()`**（HTTP **8081** LISTEN）为准；cwd 用 Linux `/proc/PID/cwd` 或 macOS `lsof -a -p PID -d cwd`（cwd 未知时不判串库）。

### 告警闭环完整复测（旧树 `~/cursorProjects/ISM源码`）

- `python3 scripts/xunzhi_jf_sync_alarm_db.py` 写入 `alarm_notice` · `StartupAlarmDelay` · **`{"DelayMinutes":0}`**，并 `AlarmTriggerEdit` 热加载。
- `bash scripts/xunzhi_jf_alarm_demo.sh`：**GetCurrentAlarmList** 非空 → **AlarmOpt**（`update` 须 **`duid`**=DeviceUuid，勿 `device_uuid`）→ 响应 **`code==0`** 须绿；Mock 漏水 **0→1** 用热切换（勿 pkill Mock）。
- **sync 写 Delay=0 后须保持 Mock:1502，重启 `ism_server`**，抑警窗才关闭。

```bash
export XUNZHI_ISM_ROOT=~/cursorProjects/ISM源码
export ARTIFACTS_DIR=/tmp/xunzhi-artifacts
mkdir -p "$ARTIFACTS_DIR"
cd "$XUNZHI_ISM_ROOT"
git fetch && git checkout feature/xunzhi-jf-env-monitor && git pull
git rev-parse --short=8 HEAD
# Mock 先在线（--flip）
pgrep -f xunzhi_jf_modbus_mock || python3 scripts/xunzhi_jf_modbus_mock.py --flip --port 1502 &
sleep 2
python3 scripts/xunzhi_jf_sync_alarm_db.py --print-env   # 写入 DelayMinutes=0 + 热加载
# 保持 Mock，重启 ism 使 Delay=0 生效
cd "$XUNZHI_ISM_ROOT/ism_server_user" && (pkill -x ism_server || true); sleep 1; ./ism_server &
sleep 5
cd "$XUNZHI_ISM_ROOT" && bash scripts/xunzhi_jf_alarm_demo.sh
# 通过：末尾「告警闭环脚本完成」= 入列 + AlarmOpt code==0 均绿（500 不得算过）
```

> 自 `1101f209` 起 **AlarmOpt `duid` + code 断言**；自 `819fd616` 起 mac **`lsof` cwd / `ps etime`**。pull 后 `git rev-parse --short=8 HEAD` 再跑。

---

## 导入包 regenerate / append 顺序

从干净 W1 重生成时（**已提交仓内的 `xunzhi-jf_complete_project_package.json` 通常已跑完全部 append**，仅改脚本时需重跑）：

```bash
python3 projects-import/xunzhi-jf/generate_xunzhi_jf_package.py
python3 projects-import/xunzhi-jf/append_w2_power_ups.py
python3 projects-import/xunzhi-jf/append_w3_ac_h2.py
python3 projects-import/xunzhi-jf/append_s1_envmon_demo.py
python3 projects-import/xunzhi-jf/append_s1_chart_thicken.py
python3 projects-import/xunzhi-jf/append_s1_gauge_plus.py
python3 projects-import/xunzhi-jf/append_w4_battery_sample.py
python3 projects-import/xunzhi-jf/append_s6_25d_nav_shell.py
python3 projects-import/xunzhi-jf/append_s6_chart_catalog_push.py
```

**导入项目**（新库或删项目后）：

```bash
python3 scripts/import_xunzhi_jf_project.py
```

项目已存在且仅告警字段需对齐：

```bash
python3 scripts/xunzhi_jf_sync_alarm_db.py
```

---

## 后端周期探针（`ism_server` 日志）

| 探针 | 设备范围 | 说明 |
|------|----------|------|
| **P-COL-01** `cycle_log_th_leak` | W1 温湿度/漏水 | `ism_server_user/protocol/common/cycle_obs.go` |
| **P-COL-02** `cycle_log_power_ups` | W2 列头柜/UPS | slave 3/4 |
| **P-COL-03** `cycle_log_ac_h2` | W3 空调/氢气 | slave 5/6 |

Mock 对照：`scripts/xunzhi_jf_modbus_mock.py`（slave **1–7**，7=电池抽样）。  
WebSocket：以 `ism_server_user/conf/app.conf` 的 **wsport** 为准（dev 常见 **10225**）。

---

## 环境冒烟（非 Mock 域 demo）

| 脚本 | 用途 | 成功判据 |
|------|------|----------|
| `scripts/xunzhi_jf_smoke_env.sh` | Mock + 后端 + 导入 | 登录 API、端口监听 |
| `scripts/xunzhi_jf_collect_evidence.sh` | 证据落盘 | 日志含 `cycle_log_*` 等 |

---

## Mock / 竖切 demo.sh 一览

| 域 | 命令 | 自证期望（末尾关键行） |
|----|------|------------------------|
| **告警** | `bash scripts/xunzhi_jf_alarm_demo.sh` | `告警闭环脚本完成（非正式验收结论）`（sync+热加载；Mock 0→1；入列断言） |
| **S2 资产/U 位** | `bash scripts/xunzhi_jf_asset_u_demo.sh` | `DEMO PASS · Store: .../dcim_asset_store.json` |
| **S3 容量/NSPC** | `bash scripts/xunzhi_jf_capacity_demo.sh` | `DEMO PASS · Store: .../dcim_capacity_store.json` |
| **S4 工单** | `bash scripts/xunzhi_jf_workorder_demo.sh` | `DEMO PASS · Store: .../dcim_workorder_store.json` |
| **S5 能耗/PUE** | `bash scripts/xunzhi_jf_energy_demo.sh` | `DEMO PASS · Store: .../dcim_energy_store.json` |
| **S6 安防/消防** | `bash scripts/xunzhi_jf_security_demo.sh` | `DEMO PASS · Store: .../security_event_store.json` |
| **电池抽样** | `bash scripts/xunzhi_jf_battery_demo.sh` | `DEMO PASS · spec: .../xunzhi-jf_battery_device_points.json` |

一键跑通 **JSON Mock 域**（无需后端，告警除外）：

```bash
bash scripts/xunzhi_jf_asset_u_demo.sh && \
bash scripts/xunzhi_jf_capacity_demo.sh && \
bash scripts/xunzhi_jf_workorder_demo.sh && \
bash scripts/xunzhi_jf_energy_demo.sh && \
bash scripts/xunzhi_jf_security_demo.sh && \
bash scripts/xunzhi_jf_battery_demo.sh
```

---

## 组态与文档索引

| 主题 | 文档 |
|------|------|
| 平面页 / S1 | `docs/ISM-寻智-S1-envmon演示.md`、`docs/ISM-寻智-S1-图表仪表盘清单草稿.md` |
| 2.5D 壳（降级平面） | `docs/ISM-寻智-S6-2.5D导航壳.md` |
| S2–S6 骨架 | `docs/ISM-寻智-S2-资产U位骨架.md` … `S6-安防消防事件桩.md` |
| 电池 / 双机 | `docs/ISM-寻智-电池1920抽样.md`、`docs/ISM-寻智-双机软件就绪演练.md` |
| W1 竖切 | `docs/ISM-寻智-W1竖切验收.md` |

**禁止**：将 Mock demo 记为 dcim 全产品或正式商务验收。
