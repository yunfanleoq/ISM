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
