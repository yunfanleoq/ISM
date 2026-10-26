---
name: ism-service-startup
description: >
  ISM 前后端服务启动与存活验证（macOS）。禁止 setsid、禁止短 shell nohup 假启动；
  登录不了先查进程/端口/API。触发词：ISM 启动、重启 ISM、登录不了、7080、8081、
  ism_server、vue-cli-service serve、服务没起来、setsid、nohup 进程消失。
disable-model-invocation: false
---

# ISM 服务启动（macOS）

启动 ISM 开发环境（后端 8081 + 前端 7080）并**验证真正存活**。与 `ism-dev-compile-preclean` 配合：前端 serve 前必须先内存 PASS。

规则铁律见 `.cursor/rules/ism-service-startup-macos.mdc`。

---

## 触发词

ISM 启动、重启、登录不了、7080、8081、服务没起来、setsid、nohup 进程消失、command not found setsid

---

## 历史教训（必读，禁止重复）

### 教训 1：`setsid` 在 macOS 不存在

- **错误**：`setsid ./ism_server &` 或 `setsid npx vue-cli-service ...`
- **日志**：`(eval):1: command not found: setsid`
- **结果**：命令整段失败，**没有任何服务启动**
- **规则**：macOS **永远不要用 setsid**；这是 Linux 专用

### 教训 2：Cursor 短 shell + `nohup ... &` 进程被带走

- **错误**：在 `block_until_ms=5000` 的 shell 里 `nohup ./ism_server & echo PID`
- **现象**：echo 出 PID，几秒后 `pgrep` 无进程、`lsof` 无 7080/8081
- **根因**：父 shell 结束后后台作业被清理（即使有 nohup）
- **规则**：Agent 内用 **`block_until_ms=0` + `exec`**；脚本用 **`disown`** 或 `start_ism_dev.sh`

### 教训 2b：Agent 内跑 `start_ism_dev.sh` 也会假成功（已多次复现）

- **错误**：`block_until_ms=0 ./scripts/start_ism_dev.sh`（脚本内虽有 nohup+disown）
- **现象**：脚本打印「启动流程结束」、PID、端口均正常；**脚本退出后 3~10 秒** 8081/7080 全无
- **根因**：Cursor Agent 后台 shell 结束时，会清理该 shell 拉起的整棵进程树；`disown` 挡不住
- **规则**：**Agent 禁止调用 `start_ism_dev.sh`**；只用方式 A（两个 `exec` 持久后台终端）。`start_ism_dev.sh` **仅**给用户在本机终端执行

### 教训 3：未验证就报「重启成功」

- **错误**：看到 echo PID 或日志前几行就告诉用户可以去登录
- **现象**：用户仍登录不了
- **规则**：必须完成下方「启动后验证清单」全部通过

---

## Instructions

### 第 0 步：确认当前状态（任何「重启」前先做）

```bash
pgrep -fl "ism_server|vue-cli-service"
lsof -nP -iTCP:7080,8081 -sTCP:LISTEN
curl -s -o /dev/null -w "8081:%{http_code}\n" http://127.0.0.1:8081/
curl -s -o /dev/null -w "7080:%{http_code}\n" http://127.0.0.1:7080/
```

若端口已在监听且登录 API 正常，**不要无谓重启**。

### 第 1 步：内存检查（仅启动前端前）

```bash
./scripts/check_mem_before_compile.sh   # 必须 RESULT: PASS
```

FAIL → 按 `ism-dev-compile-preclean` 释内存，禁止 serve。

### 第 2 步：清场

```bash
launchctl remove com.ism.frontend 2>/dev/null
pkill -9 -f "vue-cli-service" 2>/dev/null
pkill -9 -f "ism_server" 2>/dev/null
lsof -ti :7080 | grep -v Cursor | xargs kill -9 2>/dev/null
lsof -ti :8081 | xargs kill -9 2>/dev/null
sleep 3
pgrep -fl "ism_server|vue-cli-service"   # 应无输出
```

### 第 3 步：启动（二选一）

#### 方式 A — Cursor Agent（推荐）

两个独立 Shell 调用，均 **`block_until_ms=0`**：

```bash
cd ism_server_user && exec ./ism_server >> /tmp/ism_be.log 2>&1
```

```bash
cd ism-front-end-v2 && NODE_OPTIONS="--max-old-space-size=20480 --openssl-legacy-provider" \
  exec npx vue-cli-service serve --port 7080 >> /tmp/ism_fe.log 2>&1
```

#### 方式 B — 项目脚本（终端 / 脱离 Cursor）

```bash
./scripts/start_ism_dev.sh
```

### 第 4 步：启动后验证清单（全部通过才算成功）

- [ ] `pgrep -fl ism_server` 有进程
- [ ] `pgrep -fl vue-cli-service` 有进程
- [ ] `lsof` 显示 `*:8081` 和 `*:7080` LISTEN
- [ ] 后端登录：`curl ... 8081/login` → `code: 1000`
- [ ] 前端代理：`curl ... 7080/api/login` → `code: 1000`（需等编译完成）
- [ ] `/tmp/ism_fe.log` 含 `Compiled successfully`
- [ ] `curl -o /dev/null -w "%{http_code}" http://127.0.0.1:7080/` → `200`

**任一项失败 → 不能告诉用户「已启动」，继续查日志排障。**

---

## Examples

### 登录不了 — 正确排障

```bash
# 1. 进程没了 → 启动方式错了（setsid / 短 shell nohup），按本技能重起
# 2. 8081 无监听 → 查 /tmp/ism_be.log 开头是否 "http server Running on http://:8081"
# 3. 8081 正常、7080 无 → 前端未起或仍在 69% 编译，tail -f /tmp/ism_fe.log
# 4. 都监听、login 1000 → 密码/浏览器缓存问题，非服务问题
```

### 日志里出现 setsid 错误

```text
(eval):1: command not found: setsid
```

**判定**：本次启动完全失败。改用方式 A 或 `start_ism_dev.sh`，勿再试 setsid。

---

## Performance Notes

- 前端首次编译约 15~30s，69% 停很久属正常（大项目 webpack），以 `Compiled successfully` 为准
- 后端启动约 2s 即可监听 8081；登录 API 可比前端编译先测
- 内存门槛：可用+可释放 ≥ 12GB，Node heap 20GB（见 `ism-dev-compile-preclean`）

---

## Troubleshooting

| 现象 | 最可能原因 | 处理 |
|---|---|---|
| `command not found: setsid` | macOS 无此命令 | 禁止 setsid；用 exec 后台或 start_ism_dev.sh |
| PID 打印后进程消失 | 短 shell 清理子进程 | block_until_ms=0 + exec |
| 7080:000 / 连接拒绝 | 前端未起或编译中/崩了 | tail /tmp/ism_fe.log；查 OOM |
| 8081 无响应 | 后端未起或崩了 | tail /tmp/ism_be.log |
| 7080 有、/api/login 失败 | 前端未编译完或代理未就绪 | 等 Compiled successfully |
| 直连 8081 login 1000、浏览器登不了 | 前端代理或缓存 | 硬刷新 Cmd+Shift+R |
| 反复「重启成功」仍登不了 | 从未真正验证端口 | 严格执行验证清单 |

---

## 禁止事项

- **禁止** macOS 使用 `setsid`
- **禁止** 短生命周期 shell 里 `nohup ... &` 后未验证就报成功
- **禁止** 登录失败时跳过进程/端口检查直接改密码
- **禁止** 用 `http-server` / `dist/` 替代 dev server（无 `/api` 代理）
