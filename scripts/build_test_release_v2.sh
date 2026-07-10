#!/bin/bash
# ISM 正式测试环境部署包构建脚本 v2
# 用法: bash scripts/build_test_release_v2.sh
# 产出: releases/ism-release-sqlite-YYYYMMDD-HHMM-xxxx.zip
#
# 后端编译策略（macOS arm64 → Linux amd64 + SQLite/CGO）：
#   1. Docker --platform linux/amd64 + CGO_ENABLED=1（首选，最可靠）
#   2. 若 Docker 失败，包内不含二进制，附带 build-on-target.sh 在目标机编译
#   3. 部署包自包含 ism_server，严禁依赖 /opt/ISMCode/ism_web* 客户目录

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
# shellcheck source=lib_build_id.sh
source "$ROOT/scripts/lib_build_id.sh"
BUILD_ID="$(ism_build_id)"
DATE_TAG="$BUILD_ID"
PKG_NAME="ism-release-sqlite-${BUILD_ID}"
STAGING="$ROOT/releases/${PKG_NAME}"
ZIP_OUT="$ROOT/releases/${PKG_NAME}.zip"

BACKEND_SRC="$ROOT/ism_server_user"
FRONTEND_DIST="$ROOT/ism-front-end-v2/dist"
LINUX_BIN="$STAGING/ism_server_user/ism_server"
BUILD_METHOD="unknown"

echo "=== ISM 正式测试环境部署包构建 v2 ==="
echo "包名: ${PKG_NAME}"
echo ""

build_backend_docker() {
  local out="$1"
  if ! command -v docker >/dev/null 2>&1; then
    echo "  Docker 不可用"
    return 1
  fi
  local images=(
    "golang:1.22-bookworm"
    "docker.1ms.run/golang:1.22-bookworm"
    "docker.m.daocloud.io/library/golang:1.22-bookworm"
  )
  local img=""
  for candidate in "${images[@]}"; do
    if docker pull "$candidate" >/dev/null 2>&1; then
      img="$candidate"
      break
    fi
  done
  if [[ -z "$img" ]]; then
    echo "  无法拉取 golang 镜像（网络问题）"
    return 1
  fi
  echo "  使用镜像: $img"
  mkdir -p "$(dirname "$out")"
  local out_dir
  out_dir="$(cd "$(dirname "$out")" && pwd)"
  local src_dir
  src_dir="$(cd "$BACKEND_SRC" && pwd)"
  echo "[1/5] Docker 编译后端 (linux/amd64, CGO_ENABLED=1) ..."
  docker run --rm --platform linux/amd64 \
    -v "${src_dir}:/src" \
    -v "${out_dir}:/out" \
    -w /src \
    -e GOOS=linux \
    -e GOARCH=amd64 \
    -e CGO_ENABLED=1 \
    "$img" \
    go build -ldflags "-w -s" -o /out/ism_server .
}

# 1. 后端
rm -rf "$STAGING"
mkdir -p "$STAGING/ism_server_user/data/db" "$STAGING/web/dist" "$STAGING/scripts" "$STAGING/logs"

if build_backend_docker "$LINUX_BIN"; then
  BUILD_METHOD="docker-linux-amd64-cgo"
  file "$LINUX_BIN" | grep -q "ELF 64-bit" || { echo "错误: Docker 产物不是 linux amd64 ELF"; exit 1; }
  echo "  后端构建成功: $(file -b "$LINUX_BIN")"
else
  echo "[1/5] Docker 构建失败，尝试使用工作区内已编译的 linux/amd64 二进制 ..."
  BUILD_METHOD="target-build-required"
  rm -f "$LINUX_BIN"
  FALLBACK_BINS=(
    "$ROOT/releases/ism-test-20260703/ism_server_user/ism_server"
    "$ROOT/patches/alarm-clear-all-v1/ism_server"
  )
  for fb in "${FALLBACK_BINS[@]}"; do
    if [[ -f "$fb" ]] && file "$fb" | grep -q "ELF 64-bit.*x86-64"; then
      cp "$fb" "$LINUX_BIN"
      chmod 755 "$LINUX_BIN"
      BUILD_METHOD="fallback-local-linux-amd64"
      echo "  使用本地 ELF: $fb"
      break
    fi
  done
  if [[ ! -f "$LINUX_BIN" ]]; then
    echo "  无可用 linux/amd64 二进制，包内将不含 ism_server（目标机执行 build-on-target.sh）"
  fi
fi

# 2. 前端 dist
if [[ "${SKIP_FRONTEND_BUILD:-0}" == "1" ]] && [[ -f "$FRONTEND_DIST/index.html" ]]; then
  echo "[2/5] 跳过前端构建 (SKIP_FRONTEND_BUILD=1，使用已有 dist) ..."
else
  echo "[2/5] 编译前端 dist ..."
  (cd "$ROOT/ism-front-end-v2" && NODE_OPTIONS="--max-old-space-size=20480 --openssl-legacy-provider" npm run build)
fi
[[ -f "$FRONTEND_DIST/index.html" ]] || { echo "错误: dist/index.html 不存在"; exit 1; }

# 3. 组装目录
echo "[3/5] 组装部署目录 ..."
if [[ -f "$LINUX_BIN" ]]; then
  chmod 755 "$LINUX_BIN"
fi
rsync -a --delete \
  --exclude 'vendor/' \
  --exclude '*.go' \
  --exclude 'logs/' \
  --exclude 'data/dbbackup/' \
  --exclude 'data/tempDir/' \
  --exclude 'data/upload/' \
  "$BACKEND_SRC/conf/" "$STAGING/ism_server_user/conf/"
rsync -a "$BACKEND_SRC/static/" "$STAGING/ism_server_user/static/" 2>/dev/null || {
  echo "  警告: rsync static 部分失败，改用 cp -a ..."
  mkdir -p "$STAGING/ism_server_user/static"
  cp -a "$BACKEND_SRC/static/." "$STAGING/ism_server_user/static/" || true
}
rsync -a "$BACKEND_SRC/data/auth/" "$STAGING/ism_server_user/data/auth/" 2>/dev/null || mkdir -p "$STAGING/ism_server_user/data/auth"
rsync -a "$FRONTEND_DIST/" "$STAGING/web/dist/"
cp "$ROOT/scripts/serve_test_frontend.py" "$STAGING/scripts/"
cp "$ROOT/scripts/modbus_simulator.py" "$STAGING/scripts/" 2>/dev/null || true

# 测试包专用端口：前端 7080 对齐 cpolar largescreen；后端 8091 避开客户 8081/8082
DEFAULT_FE_PORT=7080
DEFAULT_BE_PORT=8091
echo "  设置测试包端口: 前端 ${DEFAULT_FE_PORT}, 后端 ${DEFAULT_BE_PORT} ..."
APP_CONF="$STAGING/ism_server_user/conf/app.conf"
if [[ "$(uname -s)" == "Darwin" ]]; then
  sed -i '' "s/^httpport=.*/httpport=${DEFAULT_BE_PORT}/" "$APP_CONF"
else
  sed -i "s/^httpport=.*/httpport=${DEFAULT_BE_PORT}/" "$APP_CONF"
fi
cat > "$STAGING/ports.env" << PORTEOF
# ISM 测试包端口（与客户生产错开，勿占用 ism_web / ism_webchaifa 端口）
# 启动前可 source ports.env 或 export 覆盖
ISM_FE_PORT=${DEFAULT_FE_PORT}
ISM_BE_PORT=${DEFAULT_BE_PORT}
PORTEOF

echo "  备份 ism.db ..."
if command -v sqlite3 >/dev/null 2>&1; then
  sqlite3 "$BACKEND_SRC/data/db/ism.db" ".backup '$STAGING/ism_server_user/data/db/ism.db'"
else
  cp -a "$BACKEND_SRC/data/db/ism.db" "$STAGING/ism_server_user/data/db/ism.db"
fi

# 目标机构建脚本（Docker 失败时的兜底）
cat > "$STAGING/build-on-target.sh" << 'BUILDEOF'
#!/bin/bash
# 在麒麟/Linux 目标机上本地编译后端（需安装 go 1.22+ 与 gcc）
set -euo pipefail
ROOT="$(cd "$(dirname "$0")" && pwd)"
SRC="$ROOT/ism_server_user"
OUT="$SRC/ism_server"
if [[ -f "$OUT" ]]; then
  echo "已存在 $OUT，跳过编译"
  exit 0
fi
if ! command -v go >/dev/null 2>&1; then
  echo "错误: 未安装 go。完整部署包应已包含 ism_server_user/ism_server"
  echo "  请从开发机重新构建 release 包，或安装 Go 1.22+ 与 gcc 后重试"
  exit 1
fi
echo "在目标机编译后端 CGO_ENABLED=1 ..."
cd "$SRC"
# 若包内无 .go 源码，需从开发机重新下发完整包
if [[ ! -f main.go ]] && [[ ! -f ../ism_server_user/main.go ]]; then
  echo "错误: 包内无 Go 源码。请从开发机重新构建并上传完整部署包"
  echo "  禁止从 /opt/ISMCode/ism_web 或 ism_webchaifa 复制客户生产二进制"
  exit 1
fi
CGO_ENABLED=1 go build -ldflags "-w -s" -o "$OUT" .
chmod +x "$OUT"
echo "编译完成: $(file -b "$OUT")"
BUILDEOF
chmod +x "$STAGING/build-on-target.sh"

# 4. 启动脚本（默认 7080/8091；可用 ISM_FE_PORT / ISM_BE_PORT 覆盖）
cat > "$STAGING/start-test.sh" << 'STARTEOF'
#!/bin/bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")" && pwd)"
BACKEND_DIR="$ROOT/ism_server_user"
CONF="$BACKEND_DIR/conf/app.conf"
PID_BACKEND="$ROOT/.backend.pid"
PID_FRONTEND="$ROOT/.frontend.pid"
LOG_BACKEND="$ROOT/logs/ism_server.log"
LOG_FRONTEND="$ROOT/logs/frontend.log"

# 测试包默认端口（前端 7080=cpolar；后端 8091 勿占 8081/8082）
DEFAULT_FE_PORT=7080
DEFAULT_BE_PORT=8091

[[ -f "$ROOT/ports.env" ]] && source "$ROOT/ports.env"

read_httpport() {
  local port
  port="$(grep -E '^httpport=' "$CONF" 2>/dev/null | head -1 | cut -d= -f2 | tr -d '[:space:]')"
  echo "${port:-8091}"
}

FRONTEND_PORT="${ISM_FE_PORT:-${ISM_FRONTEND_PORT:-${DEFAULT_FE_PORT}}}"
HTTPPORT="${ISM_BE_PORT:-$(read_httpport)}"

if [[ -n "${ISM_BE_PORT:-}" ]]; then
  if [[ "$(uname -s)" == "Darwin" ]]; then
    sed -i '' "s/^httpport=.*/httpport=${HTTPPORT}/" "$CONF"
  else
    sed -i "s/^httpport=.*/httpport=${HTTPPORT}/" "$CONF"
  fi
fi

mkdir -p "$BACKEND_DIR/data/sessionon" "$BACKEND_DIR/static/HistoryData" "$BACKEND_DIR/static/reportTemplete" "$BACKEND_DIR/static/RecordVideo" "$ROOT/logs"

if [[ ! -x "$BACKEND_DIR/ism_server" ]]; then
  echo "错误: 缺少 ism_server 二进制（部署包应自包含）"
  echo "  1) 执行 bash build-on-target.sh（需 Go 源码与 gcc）"
  echo "  2) 或从开发机重新构建并上传完整 release 包"
  echo "  禁止从 /opt/ISMCode/ism_web 或 ism_webchaifa 复制客户生产文件"
  exit 1
fi

if [[ -f "$PID_BACKEND" ]] && kill -0 "$(cat "$PID_BACKEND")" 2>/dev/null; then
  echo "后端已在运行 (PID $(cat "$PID_BACKEND"))"
else
  echo "启动后端 (端口 ${HTTPPORT}) ..."
  cd "$BACKEND_DIR"
  chmod +x ./ism_server
  nohup ./ism_server > "$LOG_BACKEND" 2>&1 &
  echo $! > "$PID_BACKEND"
  cd "$ROOT"
  ok=0
  for i in $(seq 1 30); do
    if ss -tlnp 2>/dev/null | grep -q ":${HTTPPORT}\\b"; then
      ok=1
      break
    fi
    sleep 1
  done
  if [[ "$ok" -ne 1 ]]; then
    echo "警告: 后端 ${HTTPPORT} 健康检查未通过，请查看 $LOG_BACKEND"
    tail -20 "$LOG_BACKEND" 2>/dev/null || true
  fi
fi

if [[ -f "$PID_FRONTEND" ]] && kill -0 "$(cat "$PID_FRONTEND")" 2>/dev/null; then
  echo "前端已在运行 (PID $(cat "$PID_FRONTEND"))"
else
  echo "启动前端静态服务 (端口 ${FRONTEND_PORT}, /api -> ${HTTPPORT}) ..."
  nohup python3 "$ROOT/scripts/serve_test_frontend.py" \
    --port "${FRONTEND_PORT}" \
    --dist "$ROOT/web/dist" \
    --backend "http://127.0.0.1:${HTTPPORT}" \
    > "$LOG_FRONTEND" 2>&1 &
  echo $! > "$PID_FRONTEND"
fi

echo ""
echo "=== ISM 测试环境已启动 ==="
echo "  后端端口: ${HTTPPORT}"
echo "  前端端口: ${FRONTEND_PORT}"
echo "  访问: http://<本机IP>:${FRONTEND_PORT}/#/login"
echo "  外网(cpolar): https://largescreen.cpolar.cn"
echo "  账号: admin / 123456"
echo "  后端日志: $LOG_BACKEND"
echo "  前端日志: $LOG_FRONTEND"
echo "  停止: bash stop-test.sh"
STARTEOF

cat > "$STAGING/stop-test.sh" << 'STOPEOF'
#!/bin/bash
ROOT="$(cd "$(dirname "$0")" && pwd)"
CONF="$ROOT/ism_server_user/conf/app.conf"
DEFAULT_FE_PORT=7080
DEFAULT_BE_PORT=8091

[[ -f "$ROOT/ports.env" ]] && source "$ROOT/ports.env"

read_httpport() {
  local port
  port="$(grep -E '^httpport=' "$CONF" 2>/dev/null | head -1 | cut -d= -f2 | tr -d '[:space:]')"
  echo "${port:-8091}"
}

FE_PORT="${ISM_FE_PORT:-${ISM_FRONTEND_PORT:-${DEFAULT_FE_PORT}}}"
BE_PORT="${ISM_BE_PORT:-$(read_httpport)}"

for f in .frontend.pid .backend.pid; do
  if [[ -f "$ROOT/$f" ]]; then
    pid=$(cat "$ROOT/$f")
    kill "$pid" 2>/dev/null && echo "已停止 PID $pid ($f)" || true
    rm -f "$ROOT/$f"
  fi
done
pkill -f "serve_test_frontend.py.*${ROOT}" 2>/dev/null || pkill -f "serve_test_frontend.py" 2>/dev/null || true
# 仅清理本测试包端口，勿动客户 7080/8081/8082
for port in "$FE_PORT" "$BE_PORT"; do
  fuser -k "${port}/tcp" 2>/dev/null || true
done
STOPEOF

chmod +x "$STAGING/start-test.sh" "$STAGING/stop-test.sh" "$STAGING/scripts/serve_test_frontend.py"

cat > "$STAGING/README-部署说明.md" << READMEEOF
# ISM 正式测试环境部署包

> ## ⚠️ 硬性隔离警告（必读）
>
> 本包 **必须** 部署到独立目录，例如 \`/opt/ism/${PKG_NAME}/\`。
>
> **严禁** 触碰、覆盖、依赖以下客户生产目录：
> - \`/opt/ISMCode/\`
> - \`/opt/ISMCode/ism_web/\`（电力生产）
> - \`/opt/ISMCode/ism_webchaifa/\`（柴发生产）
>
> 禁止从上述目录复制 \`ISMServer\` 或任何文件。本包 **自包含** \`ism_server_user/ism_server\`。

- 版本: V3.01.RC07
- 构建日期: ${DATE_TAG}
- 目标平台: **Linux x86_64 (amd64)**
- 部署目录: **/opt/ism/${PKG_NAME}/**
- 数据库: SQLite (\`dbtype=1\`)，已含循安电力监控演示数据
- 后端编译: **${BUILD_METHOD}**（SQLite 需 CGO_ENABLED=1，禁止 CGO_ENABLED=0 交叉编译）

## 目录结构

\`\`\`
${PKG_NAME}/
├── start-test.sh          # 一键启动（读 ports.env / app.conf）
├── stop-test.sh           # 停止服务
├── ports.env              # ISM_FE_PORT=7080, ISM_BE_PORT=8091
├── build-on-target.sh     # Docker 构建失败时，在目标机编译（需 Go 源码）
├── ism_server_user/
│   ├── ism_server         # Linux amd64 后端（CGO 编译）
│   ├── conf/app.conf      # 配置（dbtype=1, httpport=8091）
│   ├── data/db/ism.db
│   └── static/
├── web/dist/
└── scripts/
    ├── serve_test_frontend.py  # 静态 + /api 代理（--backend 可配置）
    └── modbus_simulator.py
\`\`\`

## 快速部署

1. 创建独立目录并上传解压:
   \`\`\`bash
   mkdir -p /opt/ism
   cd /opt/ism
   unzip ${PKG_NAME}.zip
   cd ${PKG_NAME}
   \`\`\`
2. 确认 \`ports.env\` 中端口未与客户 **8081/8082** 冲突（默认 **7080/8091**）；若需调整:
   \`\`\`bash
   export ISM_FE_PORT=7080
   export ISM_BE_PORT=8091
   sed -i 's/^httpport=.*/httpport=8091/' ism_server_user/conf/app.conf
   \`\`\`
   > **禁止**占用客户 API **8081/8082**；前端 7080 与 cpolar 一致，若客户 ism_web 前端也在 7080 则勿同时启两套前端
3. 确认包内 \`ism_server_user/ism_server\` 存在且可执行；若缺失:
   - \`bash build-on-target.sh\`（目标机有 Go+gcc 与源码）
   - 或从开发机重新构建完整 release 包
4. 启动: \`bash start-test.sh\`
5. 访问 \`http://<IP>:7080/#/login\` 或外网 \`https://largescreen.cpolar.cn\`，账号 **admin** / **123456**

## 端口说明（与客户生产隔离）

| 归属 | 目录 | 前端 | 后端 | 说明 |
|------|------|------|------|------|
| **客户电力生产** | \`/opt/ISMCode/ism_web\` | 7080 | 8081 | **勿动、勿占用** |
| **客户柴发生产** | \`/opt/ISMCode/ism_webchaifa\` | — | 8082 | **勿动、勿占用** |
| **本测试包** | \`/opt/ism/${PKG_NAME}/\` | **7080** | **8091** | cpolar 隧道固定 7080；与客户 ism_web 同端口时需二选一 |

环境变量 \`ISM_FE_PORT\` / \`ISM_BE_PORT\`（或 \`ports.env\`）可覆盖；\`start-test.sh\` 从 \`app.conf\` 读 \`httpport=\`，与 API 代理保持一致。

## 注意事项

| 项 | 说明 |
|----|------|
| 隔离 | **禁止**读写 \`/opt/ISMCode/ism_web*\`；仅使用 \`/opt/ism/${PKG_NAME}/\` |
| CGO | **必须** CGO_ENABLED=1 编译，否则 SQLite 启动 panic |
| 密码 | admin/123456；前端 MD5 → 后端 bcrypt(MD5)，**勿改密码链路** |
| 日志 | \`logs/ism_server.log\`、\`logs/frontend.log\` |
READMEEOF

echo "[4/5] 写入 BUILD_INFO ..."
if [[ -f "$LINUX_BIN" ]]; then
  BACKEND_DESC="$(file -b "$LINUX_BIN")"
else
  BACKEND_DESC="未包含（需目标机 build-on-target.sh 或重新构建 release 包）"
fi
cat > "$STAGING/BUILD_INFO.txt" << EOF
包名: ${PKG_NAME}
构建时间: $(date '+%Y-%m-%d %H:%M:%S')
构建主机: $(uname -s)/$(uname -m)
后端编译方式: ${BUILD_METHOD}
后端: ${BACKEND_DESC}
前端 dist: $(du -sh "$STAGING/web/dist" | cut -f1)
数据库: $(du -sh "$STAGING/ism_server_user/data/db/ism.db" | cut -f1)（含循安电力监控演示数据）
dbtype: $(grep '^dbtype=' "$STAGING/ism_server_user/conf/app.conf")
httpport: $(grep '^httpport=' "$STAGING/ism_server_user/conf/app.conf")
部署目录: /opt/ism/${PKG_NAME}
隔离警告: 严禁触碰 /opt/ISMCode/ism_web 与 /opt/ISMCode/ism_webchaifa
EOF

echo "[5/5] 压缩 zip ..."
[[ -f "$STAGING/start-test.sh" ]] || { echo "错误: staging 不完整，缺少 start-test.sh"; exit 1; }
[[ -f "$STAGING/web/dist/index.html" ]] || { echo "错误: staging 不完整，缺少 web/dist/index.html"; exit 1; }
[[ -f "$STAGING/ism_server_user/data/db/ism.db" ]] || { echo "错误: staging 不完整，缺少 ism.db"; exit 1; }
mkdir -p "$ROOT/releases"
rm -f "$ZIP_OUT"
# COPYFILE_DISABLE 避免 macOS 资源分叉导致 zip 报 No such file
(cd "$ROOT/releases" && COPYFILE_DISABLE=1 zip -r -y -q "$(basename "$ZIP_OUT")" "$(basename "$STAGING")")
ZIP_RC=$?
if [[ "$ZIP_RC" -ne 0 ]] && [[ "$ZIP_RC" -ne 18 ]]; then
  echo "错误: zip 失败 (exit $ZIP_RC)"
  exit "$ZIP_RC"
fi

echo ""
echo "=== 构建完成 ==="
ls -lh "$ZIP_OUT"
du -sh "$STAGING"
echo "压缩包:   $ZIP_OUT"
echo "编译方式: $BUILD_METHOD"
