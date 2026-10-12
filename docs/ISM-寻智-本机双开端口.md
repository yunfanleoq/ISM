# 寻智 xunzhi-jf · 本机双开端口（`ISM源码-xunzhi-jf`）

与旧单口树 **7080 / 8081 / Mock 1502** 区分；投屏复测以本页为准。

| 服务 | 端口 |
|------|------|
| 前端 `vue-cli-service serve` | **7081** |
| `ism_server` | **8082** |
| Modbus Mock | **1503**（`XUNZHI_MOCK_PORT` 可覆盖） |

## 固定复看顺序

1. **Mock**（保持 LISTEN）  
   `python3 scripts/xunzhi_jf_modbus_mock.py --port 1503`  
   或 `export XUNZHI_MOCK_PORT=1503` 后直接 `python3 scripts/xunzhi_jf_modbus_mock.py`

2. **后端**（须含最新 `GetDeviceData` / `finalizeEnvmonEngineering`）  
   `cd ism_server_user && go build -o ism_server . && ./ism_server`  
   确认 **8082** LISTEN。

3. **前端**  
   `cd ism-front-end-v2 && NODE_OPTIONS='--max-old-space-size=20480 --openssl-legacy-provider' npx vue-cli-service serve --port 7081`  
   确认 `vue.config.js` 代理 `/api` → `127.0.0.1:8082`（双口树配置）。

4. **登录**  
   浏览器打开 `http://127.0.0.1:7081`，`admin` / `123456`。  
   **改代码或拉 tip 后须重新登录一次**，以便 `loadRoutes` 挂上 `/xunzhi/jf/*`（含 `viz/plane`）。

5. **项目**  
   **ProjectDel + 重导** `python3 scripts/import_xunzhi_jf_project.py`，确保含页 **`xunzhi-jf-env-plane`** 且设备 `modbusPort` 为 **1503**。

6. **平面**  
   菜单 **环境平面大屏**（`/xunzhi/jf/viz/plane`）→ iframe 应为  
   `#/AppRun/<displayUid>?pageId=…&pageName=xunzhi-jf-env-plane&noRunTree=1`  
   无左侧 **组织层级总览** / RunTree 叠层；实时表温湿度约 **20–25℃ / 40–55%RH**。

## 常见 FAIL

| 现象 | 处理 |
|------|------|
| TH ~5xxx / 61℃ | Mock 未在 **1503** 或导入包仍为 1502；重导 + 重编译 **ism_server** |
| 菜单无「循智·机房环境」 | **退出重新登录** |
| 组织总览仍叠层 | 查 iframe query 是否带 `pageName` / `noRunTree=1`；硬刷新 |
