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
   重导后抽查 SQLite：`devices_model.data_format`（或导入字段 `DataFormat`）对温湿度物模型为 **`ABCD`**（勿 **CDAB**，否则 TH 易炸成 ~5xxx）。  
   浏览器 **ProjectUuid** Cookie 须为 **当前库** `project_lists` 中「寻智机房环境-xunzhi-jf」的 uuid（重导后会变，**禁止**记死某一版 uuid）；与 `GetSystemHomeDashboard.projectUuid` 一致，勿串航信库。直链 `/#/xunzhi/jf/viz/plane` 依赖 `primeXunzhiJfRouteContext` 自动对齐 Cookie。

6. **平面**  
   菜单 **环境平面大屏**（`/xunzhi/jf/viz/plane`）→ iframe 应为  
   `#/AppRun/<displayUid>?pageId=…&pageName=xunzhi-jf-env-plane&noRunTree=1`  
   无左侧 **组织层级总览** / RunTree 叠层；实时表温湿度约 **20–25℃ / 40–55%RH**。

## 喊测 / 保活注意

- **勿用** `/tmp/keep_xunzhi_be.sh`（或同类脚本）保活双口 `ism_server`：多实例/错 cwd 易**写坏或锁错** `ism.db`。
- **ProjectUuid** 以喊测当日 **现库**为准：`sqlite3 ism_server_user/data/db/ism.db "SELECT uuid,name FROM project_lists WHERE name LIKE '%寻智%';"` 或登录后看 Cookie / `ProjectList` API；**禁止**死记文档里的示例 uuid。

## 常见 FAIL

| 现象 | 处理 |
|------|------|
| 登录 **code 1003**（用户不存在） | 先确认 **`ism_server` 实际连的库**（`ism_server_user/conf/app.conf` 的 **`dbtype`** + 工作目录）。双口本机多为 **`dbtype=1`** → `ism_server_user/data/db/ism.db`（须在 `ism_server_user` 下启动，相对路径才正确）。**勿**只改 `ism.db` 而进程仍连 OceanBase（`dbtype=4`）空库。然后：`python3 scripts/xunzhi_jf_seed_admin.py --dbtype 1`（或按 conf 自动）→ 可选 `--verify-login` + `BASE_URL=http://127.0.0.1:8082`。仍 1003 时 `sqlite3 ism_server_user/data/db/ism.db "SELECT id,username,deleted_at FROM user WHERE username='admin';"`。`ProjectDel`/重导**不删** `user` 表。 |
| TH ~5xxx / 61℃ | 多为 **`DataFormat=CDAB`**（改 **ABCD** 后重导 + 重启 ism）；或 Mock 非 **1503** |
| 菜单无「循智·机房环境」 | **退出重新登录** |
| 组织总览仍叠层 | 查 iframe query 是否带 `pageName` / `noRunTree=1`；硬刷新 |

### admin 种子（1003 复测）

```bash
# 双口 SQLite（与 app.conf dbtype=1 一致时）
python3 scripts/xunzhi_jf_seed_admin.py --dbtype 1

# 种子 + 直连后端登录（8082）
BASE_URL=http://127.0.0.1:8082 python3 scripts/xunzhi_jf_seed_admin.py --dbtype 1 --verify-login
```

账号 **`admin` / `123456`**（前端发 MD5：`e10adc3949ba59abbe56e057f20f883e`）。空库首次启动 `ism_server` 也会写 admin，但历史上 `db.go` 的 `errors.Is` 误判会导致**有表无 admin**；种子脚本与已修启动逻辑均写入 `bcrypt(MD5(123456))` 同款哈希。
