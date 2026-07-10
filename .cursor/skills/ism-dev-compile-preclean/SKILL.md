---
name: ism-dev-compile-preclean
description: >
  ISM 前端 dev server 编译前内存检查、清场与释内存标准流程（macOS）。在启动 vue-cli-service serve 之前
  必须先跑 check_mem_before_compile.sh（≥12GB 可用+可释放）、杀旧进程、给足 20G heap，避免 60%~70% 因 OOM SIGKILL 卡死。
  触发词：ISM 编译、vue-cli-service serve、7080 dev server、编译卡住 60%、68%、69%、内存、
  dev server 挂掉、Webpack 构建卡死、OOM、check_mem_before_compile。
disable-model-invocation: false
---

# ISM 编译前内存检查、清场与释内存

ISM 前端（`ism-front-end-v2`）dev server 在 **60%~70%** 卡住，根因是 **内存不足（OOM SIGKILL）**。
每次启动 `vue-cli-service serve` **之前**必须先执行本技能，**不得跳过**。

完整编译铁律见 `.cursor/rules/ism-compilation.mdc`。

---

## 触发词

ISM 编译、vue-cli-service serve、7080 dev server、编译卡住 60%、68%、69%、内存、check_mem_before_compile

---

## Agent 铁律

1. **serve 前必须跑** `scripts/check_mem_before_compile.sh`；输出 `RESULT: FAIL` 时 **禁止启动** dev server。
2. **禁止** 不杀旧 vue-cli 就启动新编译（HMR 叠加吃内存）。
3. **禁止** 不用 20G heap 直接 serve 大项目（`--max-old-space-size=20480`）。
4. 内存不足时按下文顺序主动释内存；仍 FAIL 则告知用户，不编译。

---

## 内存阈值

| 指标 | 门槛 |
|------|------|
| 可用 + 可释放内存 | **≥ 12 GB** |
| Node V8 heap | **20 GB**（`--max-old-space-size=20480`） |

「可用」= `vm_stat` 的 Pages free；「可释放」= inactive + speculative + purgeable 页。

---

## 第 0 步：内存检查（必须先做）

```bash
# 项目根目录
./scripts/check_mem_before_compile.sh
```

脚本输出 `RESULT: PASS` 或 `RESULT: FAIL` 及具体建议。仅 PASS 时可进入清场与启动。

### macOS 手工检查命令（脚本已内置，排障时可单独跑）

```bash
# 页级内存分布
vm_stat | head -12

# 系统内存压力
memory_pressure

# 物理内存总量
sysctl hw.memsize

# Top 内存进程（RSS）
ps -Aem -o rss,pid,comm | sort -k1 -nr | head -12
```

---

## 内存不足时的释内存顺序（FAIL 时执行）

按顺序执行，每步后 **重新跑** `./scripts/check_mem_before_compile.sh`，直到 PASS 或确认无法继续。

```bash
# 1. 杀旧编译 / ISM 相关进程
launchctl remove com.ism.frontend 2>/dev/null
pkill -9 -f "vue-cli-service" 2>/dev/null
pkill -9 -f "ism_server" 2>/dev/null
lsof -ti :7080 | grep -v Cursor | xargs kill -9 2>/dev/null
lsof -ti :8081 | xargs kill -9 2>/dev/null
sleep 3

# 2. 停 Colima / Docker
colima stop 2>/dev/null || true
docker stop $(docker ps -q) 2>/dev/null || true

# 3. 提示用户手动关闭高内存应用
#    Chrome、Safari 多标签、IDE 多余窗口、模拟器、未用 Electron 应用等
#    参考 check_mem_before_compile.sh 输出的 Top 内存进程列表

# 4. 可选：强制释放 inactive 内存（需 sudo）
sudo purge

# 5. 再次检查
./scripts/check_mem_before_compile.sh
```

**仍 FAIL** → 停止，告知用户：当前机器可用+可释放内存不足 12GB，无法安全启动 20GB heap 的 dev server。

---

## 编译前清场 + 启动标准流程（PASS 后执行）

```bash
# 1. 清场（与释内存第 1 步相同，确保无残留）
launchctl remove com.ism.frontend 2>/dev/null
pkill -9 -f "vue-cli-service" 2>/dev/null
pkill -9 -f "ism_server" 2>/dev/null
lsof -ti :7080 | grep -v Cursor | xargs kill -9 2>/dev/null
lsof -ti :8081 | xargs kill -9 2>/dev/null
sleep 3

# 2. 确认无残留
pgrep -fl vue-cli-service   # 应无输出

# 3. 启动（20G heap）
cd ism-front-end-v2
NODE_OPTIONS="--max-old-space-size=20480 --openssl-legacy-provider" \
  npx vue-cli-service serve --port 7080
```

### 编译前确认清单

- [ ] `check_mem_before_compile.sh` 输出 **RESULT: PASS**
- [ ] 无残留 `vue-cli-service` 进程
- [ ] 7080 / 8081 端口无旧进程占用
- [ ] `launchctl` 守护 `com.ism.frontend` 已 remove
- [ ] 启动命令含 `--max-old-space-size=20480`

---

## 66-69% 卡死排障

1. 查进程是否被 SIGKILL：`log show --predicate 'eventMessage contains "killed"' --last 5m`
2. 确认日志停在 66-69% 且无新输出
3. 重新跑 `./scripts/check_mem_before_compile.sh`；若 FAIL 按释内存顺序处理
4. 按上文「编译前清场」完整重跑
5. 后台启动：**不要**在 Cursor 短 shell 里单独 `nohup ... &`（进程会被带走）。完整前后端启动见 **`ism-service-startup`** 技能：

```bash
./scripts/start_ism_dev.sh              # 推荐：含验证
# 或 Agent 内 block_until_ms=0 + exec npx vue-cli-service serve --port 7080
```

macOS **禁止 `setsid`**（Linux only，日志会出现 `command not found: setsid`）。

---

## 禁止事项

- **禁止** 跳过 `check_mem_before_compile.sh` 直接 serve
- **禁止** `RESULT: FAIL` 时仍启动 dev server
- **禁止** 用 `python3 -m http.server` 或 `http-server` serve `dist/` —— 不支持 `/api` 代理
- **禁止** 用已编译的 `dist/` 充当最新代码
- **禁止** 在不杀 `launchctl` 的情况下反复重启
- **禁止** 用 `--no-progress` 或其他 webpack 魔改参数绕过 —— 问题是内存，不是进度条
- **禁止** macOS 使用 `setsid`；服务启动见 `ism-service-startup` 技能
