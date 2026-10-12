# W1 复测证据目录（测试复测入口）

> 分支：`feature/xunzhi-jf-env-monitor` · PR [#5](https://github.com/yunfanleoq/ISM/pull/5)  
> **非「已验收」口径**；仅供冒烟/复测留痕。

## 一键环境

```bash
./scripts/xunzhi_jf_smoke_env.sh          # Mock + 后端 + 导入
./scripts/xunzhi_jf_smoke_env.sh --with-fe # 另起前端 :7080（需内存 PASS）
./scripts/xunzhi_jf_collect_evidence.sh   # I1/I2/I3 + cycle 日志
./scripts/xunzhi_jf_alarm_demo.sh         # 漏水告警 → AlarmOpt 闭环
```

## WebSocket 端口（重要）

| 来源 | 端口 |
|------|------|
| 契约/前端默认示例 | `10215`（`setting.config.js` 初值） |
| **本仓 `ism_server_user/conf/app.conf`** | **`wsport=10225`（联调以实测为准）** |

前端登录后从 `GetSystemParams` 读取 `wsport` 建链：`ws://<host>:<wsport>/ws?token=...&project=...`  
**禁止**在未协调的情况下改现网 `app.conf`；文档对齐即可。

## 前端目视（F3.3）

1. `http://127.0.0.1:7080` 登录 `admin` / `123456`
2. 切换项目 **寻智机房环境-xunzhi-jf**
3. 组态 → **寻智机房环境平面监视-xunzhi-jf** → 页 **`xunzhi-jf-env-plane`**
4. 截图/录屏保存为本目录 `F3.3_scada_*.png` 或 `artifacts/`

## 业主本机已采证据（运行态部分绿 · 2026-10-12）

见同目录 `业主本机冒烟摘录.md`（摘自项目经理包，含 P-COL-01、getRealData、cycle 日志）。

## 文件命名约定

| 用例 | 文件名模式 |
|------|------------|
| I1 | `I1_login_*.json` |
| I2 | `I2_realdata_*.json` |
| I3 | `I3_ws_realdata_*.log` |
| P-COL-01 | `P-COL-01_cycle_log_*.log` |
| F3.1 | `F3.1_monitor_tree_*.json` |
| F4 告警 | `F4_alarm_*.json` |
