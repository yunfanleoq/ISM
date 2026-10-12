# 寻智机房环境 · W1 竖切验收（温湿度 + 漏水）

> 分支：`feature/xunzhi-jf-env-monitor` · PR #5  
> **测试复测入口**：`docs/ISM-寻智-W1复测证据/README.md`

## WebSocket 端口说明

- 架构契约中的 **10215** 为设计/前端配置缺省示例。
- **联调与冒烟以运行中 `ism_server_user/conf/app.conf` 的 `wsport` 为准**；业主本机实测 **`10225`**（`:10215` 未监听）。
- 前端 `App.vue` 使用 `GetSystemParams` 返回的 `wsport` 建链，勿假设固定 10215。

## 交付物路径

| 资产 | 路径 |
|------|------|
| 导入包（含 W2 配电/UPS 扩展） | `projects-import/xunzhi-jf/xunzhi-jf_complete_project_package.json` |
| W1 点表 | `projects-import/xunzhi-jf/xunzhi-jf_device_points.json` |
| W2 点表附录 | `projects-import/xunzhi-jf/xunzhi-jf_w2_device_points.json` |
| 试打样包 | `dist/xunzhi-jf/ism-patch-xunzhi-jf-YYYYMMDD-<sha7>/` |
| 一键冒烟 | `scripts/xunzhi_jf_smoke_env.sh` |

## 组态

项目 **寻智机房环境-xunzhi-jf** → **寻智机房环境平面监视-xunzhi-jf** → **`xunzhi-jf-env-plane`**

## 调用链

Modbus(Mock:1502) → DeviceRealData → WS(`wsport`) RealData → 组态 `ism-view-real-table`

## 探针

- P-COL-01：`cycle_log_th_leak`（温湿度/漏水）
- P-COL-02：`cycle_log_power_ups`（W2 配电/UPS 设备名含 PWR/UPS）
- P-COL-03：`cycle_log_ac_h2`（W3 设备名含 `-AC-` / `-H2-`）

**W2 探针自证（重导或新库后）**：`ImportProject` 含 PWR/UPS 设备 → 起 Mock（slave 3/4）+ 重编/重启 `ism_server` → `grep cycle_log_power_ups` 后端日志应见 `xunzhi-jf-PWR-01` / `xunzhi-jf-UPS-01`。

## 演示唯一工作树（防双目录串库）

| 正确 | 错误 |
|------|------|
| 固定一个仓库根，例如业主本机 `~/cursorProjects/ISM源码` | 同时存在 `ISM源码` 与 `ISM源码-xunzhi-jf` 各起一套 `ism_server` |
| `cd <根>/ism_server_user && ./ism_server` | 在 A 树改 sqlite、在 B 树起后端（库路径不一致） |
| `export XUNZHI_ISM_ROOT=<根>` 后跑脚本 | 未设环境变量且从 patch 子目录误跑脚本 |

脚本会打印 **`ism_server` cwd** 与 **`ism.db` 绝对路径**；二者须指向同一棵树下的 `ism_server_user/data/db/ism.db`。

空库 / 无 `xunzhi-jf-LEAK-01`：在本树 `python3 scripts/import_xunzhi_jf_project.py` 后再 sync。  
**组态/平面需更新**（含 W2 区块）：`python3 scripts/import_xunzhi_jf_project.py --force`（`ProjectDel` + 重导，uuid 会变）。

## 告警复测（关阻塞项）

1. **已存在项目**：`python3 scripts/xunzhi_jf_sync_alarm_db.py`  
   - 写 `alarm_on_value=1`、补 `alarm_trigger`  
   - 默认对每个触发器 `AlarmTriggerEdit` → 等价 `AlarmTriggerCloseChan` 热加载  
   - 若跳过热加载：加 `--skip-reload`，并**必须从本树 `ism_server_user` 重启** `ism_server`
2. **新库/可删项目**：`python3 scripts/import_xunzhi_jf_project.py` 或 `--force`（包内 `alarmTriggers`≥2`）
3. **边沿造数**：勿单独长期 `--leak-on` 指望入列；须 **先漏水=0 再 =1**（见 `xunzhi_jf_alarm_demo.sh`）  
   - 刚重启 `ism_server` 后存在 **StartupAlarmDelay（默认可达 10min）** 抑警窗，期间须保持漏水=0 再打 0→1
4. **起服顺序**：Mock `127.0.0.1:1503`（双口默认，`XUNZHI_MOCK_PORT`）LISTEN → 本树 `ism_server`；无 Mock 会 `connection refused`、无边沿。
5. **抑警窗**：`sync`/`demo` 写入 `StartupAlarmDelay` `{"DelayMinutes":0}`；`UpdateAlarmNotice` 会 **ReconfigureStartupAlarmWindows**（tip `d628f7f9+` 后端）。仍不入列时再重启 `ism_server`（Mock 保持）或等待旧 10min 窗结束。
6. **闭环**：`ISM_HTTP_PORT=8082 BASE_URL=http://127.0.0.1:8082 bash scripts/xunzhi_jf_alarm_demo.sh` → 证据 `docs/ISM-寻智-W1复测证据/F4_*`  
   - 断言：Mock LISTEN、LEAK realdata 新鲜、再 0→1、`devices_alarm_list` / `GetCurrentAlarmList`
