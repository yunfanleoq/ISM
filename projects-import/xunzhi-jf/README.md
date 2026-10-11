# 寻智机房环境 · xunzhi-jf W1 竖切包

温湿度 + 漏水 Modbus TCP 竖切，走 **协议 → DeviceRealData → WS(:10215 RealData) → 组态**。

## 资产路径

| 文件 | 说明 |
|------|------|
| `xunzhi-jf_complete_project_package.json` | `POST /ImportProject` 正式导入包 |
| `xunzhi-jf_device_points.json` | 点表说明（含 `gather_interval_ms`） |
| `generate_xunzhi_jf_package.py` | 重新生成上述 JSON |

组态画面（导入后）：项目 **寻智机房环境-xunzhi-jf** → 组态模型 **寻智机房环境平面监视-xunzhi-jf** → 页面 `xunzhi-jf-env-plane`。

## 自证步骤（开发机）

```bash
# 1) Mock 采集（1502，周期 5s）
python3 scripts/xunzhi_jf_modbus_mock.py

# 2) 启动后端（需已编译 ism_server）
cd ism_server_user && ./ism_server >> /tmp/ism_be.log 2>&1 &

# 3) 导入项目包
python3 scripts/import_xunzhi_jf_project.py

# 4) 周期日志 P-COL-01
grep cycle_log_th_leak /tmp/ism_be.log | tail -5
# 期望：gather_interval_ms=5000 cycle_ms<10000 points=3

# 5) 实时值（导入后替换 PROJECT_UUID）
curl -s -X POST http://127.0.0.1:8081/login -H 'Content-Type: application/json' \
  -d '{"Username":"admin","password":"e10adc3949ba59abbe56e057f20f883e"}'
```

前端：登录 → 切换项目 **寻智机房环境-xunzhi-jf** → 打开组态 **寻智机房环境平面监视-xunzhi-jf**，表格应刷新温湿度/漏水列。

## 复用说明（相对 1A 配电室）

- 导入链路同 `liu-chang-1A-dev/1A_complete_project_package.json`（`ImportProject` + `monitorTree` + `displayLayer`）
- `extra` 字段沿用 `modbusIP` / `slaveId` / `packTime` 约定（`normalizeExtraData`）
- 组态绑点采用 `ism-view-real-table` 的 `设备名->测点名` 绑定（同 `build_ncc_dashboard.py`）

## 禁止

- 不覆盖麒麟生产 / `patches/` / `releases/` 现网包
