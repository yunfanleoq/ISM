#!/bin/bash
# ISM OceanBase 一体部署包构建脚本
# 用法: bash scripts/build_oceanbase_release.sh
# 产出: releases/ism-release-oceanbase-YYYYMMDD-HHMM-xxxx-offline.zip
#
# 包内含：ism_server + 前端 dist + OceanBase Docker 编排 + MySQL 业务备份 + 导入脚本
# 正式环境：解压 → 配置 ports.env / app.conf → bash start-all.sh

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
# shellcheck source=lib_build_id.sh
source "$ROOT/scripts/lib_build_id.sh"
BUILD_ID="$(ism_build_id)"
DATE_TAG="${BUILD_ID%%-*}"   # YYYYMMDD（兼容预导入镜像 tag）
PKG_NAME="ism-release-oceanbase-${BUILD_ID}"
STAGING="$ROOT/releases/${PKG_NAME}"
ZIP_OUT="$ROOT/releases/${PKG_NAME}-offline.zip"

BACKEND_SRC="$ROOT/ism_server_user"
FRONTEND_DIST="$ROOT/ism-front-end-v2/dist"
# 可通过 MYSQL_BACKUP=... 指定最新 dump（默认仍用历史权威备份文件名）
MYSQL_BACKUP="${MYSQL_BACKUP:-$ROOT/Mysql_Backup_2026-07-08_15-52-44.sql}"
LINUX_BIN="$STAGING/ism_server_user/ism_server"
BUILD_METHOD="unknown"
OB_IMAGE="oceanbase/oceanbase-ce:latest"
OB_TAR="$STAGING/oceanbase/oceanbase-ce.tar"

DEFAULT_FE_PORT=7090
DEFAULT_BE_PORT=8091
DEFAULT_OB_PORT=2881
DEFAULT_TD_PORT=6041
DEFAULT_TD_NATIVE_PORT=6030
TD_IMAGE="${TD_IMAGE:-tdengine/tdengine:3.3.6.13}"
PRELOAD_OB_DATA="${PRELOAD_OB_DATA:-1}"

ob_start_container() {
  local port="${1:-2881}"
  local tenant="${2:-ism_tenant}"
  local password="${3:-ism2024!}"
  docker rm -f oceanbase 2>/dev/null || true
  docker run -d --name oceanbase --restart unless-stopped \
    --ulimit nofile=65536:65536 --ulimit nproc=65536:65536 \
    -p "${port}:2881" \
    -e MODE=mini \
    -e OB_MEMORY_LIMIT=8G \
    -e OB_DATAFILE_SIZE=10G \
    -e OB_LOG_DISK_SIZE=5G \
    -e OB_CLUSTER_NAME=ism_cluster \
    -e OB_TENANT_NAME="${tenant}" \
    -e OB_TENANT_PASSWORD="${password}" \
    "$OB_IMAGE"
}

ob_compose_up() {
  local compose_file="$1"
  if docker compose version >/dev/null 2>&1; then
    docker compose -f "$compose_file" up -d
  elif command -v docker-compose >/dev/null 2>&1; then
    docker-compose -f "$compose_file" up -d
  else
    local port="${OB_PORT:-2881}"
    local tenant="${OB_TENANT:-ism_tenant}"
    local password="${OB_PASSWORD:-ism2024!}"
    ob_start_container "$port" "$tenant" "$password"
  fi
}

echo "=== ISM OceanBase 一体部署包构建 ==="
echo "构建 ID: ${BUILD_ID}"
echo "包名: ${PKG_NAME}"
echo ""

build_backend_kylin_static() {
  local out="$1"
  mkdir -p "$(dirname "$out")"
  echo "[1/6] 编译麒麟 V10 兼容后端 (linux/amd64, CGO_ENABLED=0, 静态链接) ..."
  if command -v go >/dev/null 2>&1; then
    (cd "$BACKEND_SRC" && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
      go build -ldflags "-w -s" -o "$out" .)
    return 0
  fi
  if ! command -v docker >/dev/null 2>&1; then
    return 1
  fi
  local images=(
    "golang:1.22-bullseye"
    "docker.1ms.run/golang:1.22-bullseye"
    "docker.m.daocloud.io/library/golang:1.22-bullseye"
  )
  local img=""
  for candidate in "${images[@]}"; do
    if docker pull "$candidate" >/dev/null 2>&1; then
      img="$candidate"
      break
    fi
  done
  [[ -n "$img" ]] || return 1
  local out_dir src_dir
  out_dir="$(cd "$(dirname "$out")" && pwd)"
  src_dir="$(cd "$BACKEND_SRC" && pwd)"
  docker run --rm --platform linux/amd64 \
    -v "${src_dir}:/src" -v "${out_dir}:/out" -w /src \
    -e GOOS=linux -e GOARCH=amd64 -e CGO_ENABLED=0 \
    "$img" go build -ldflags "-w -s" -o "/out/$(basename "$out")" .
}

verify_kylin_binary() {
  local bin="$1"
  [[ -f "$bin" ]] || return 1
  file "$bin" | grep -q "ELF 64-bit.*x86-64" || return 1
  # 麒麟 V10 glibc 2.28：拒绝依赖更高版本 GLIBC 的动态链接二进制
  if strings "$bin" 2>/dev/null | grep -qE 'GLIBC_2\.(3[2-9]|[4-9][0-9])'; then
    return 1
  fi
  return 0
}

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
    echo "  无法拉取 golang 镜像"
    return 1
  fi
  echo "  使用镜像: $img"
  mkdir -p "$(dirname "$out")"
  local out_dir src_dir
  out_dir="$(cd "$(dirname "$out")" && pwd)"
  src_dir="$(cd "$BACKEND_SRC" && pwd)"
  echo "[1/6] Docker 编译后端 (linux/amd64, CGO_ENABLED=1) ..."
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

rm -rf "$STAGING"
mkdir -p "$STAGING/ism_server_user/data/db" "$STAGING/web/dist" "$STAGING/scripts" "$STAGING/logs" "$STAGING/oceanbase"

PATCH_BIN="$ROOT/patches/ism-server-kylin-glibc228/ism_server"
LEGACY_PATCH="$ROOT/patches/ism-server-cgo-linux-amd64/ism_server"
if [[ -f "$PATCH_BIN" ]] && verify_kylin_binary "$PATCH_BIN"; then
  echo "[1/6] 使用麒麟 glibc2.28 静态二进制 ..."
  cp "$PATCH_BIN" "$LINUX_BIN"
  chmod 755 "$LINUX_BIN"
  BUILD_METHOD="kylin-glibc228-static"
elif build_backend_kylin_static "$LINUX_BIN"; then
  verify_kylin_binary "$LINUX_BIN" || { echo "错误: 编译产物 glibc 不兼容麒麟 V10"; exit 1; }
  cp "$LINUX_BIN" "$PATCH_BIN" 2>/dev/null || true
  BUILD_METHOD="kylin-glibc228-static-build"
elif [[ -f "$LEGACY_PATCH" ]] && verify_kylin_binary "$LEGACY_PATCH"; then
  echo "[1/6] 使用兼容 CGO 二进制 (linux/amd64) ..."
  cp "$LEGACY_PATCH" "$LINUX_BIN"
  chmod 755 "$LINUX_BIN"
  BUILD_METHOD="patch-cgo-linux-amd64"
elif build_backend_docker "$LINUX_BIN" && verify_kylin_binary "$LINUX_BIN"; then
  BUILD_METHOD="docker-linux-amd64-cgo"
else
  echo "[1/6] 尝试本地静态编译 ..."
  BUILD_METHOD="target-build-required"
  rm -f "$LINUX_BIN"
  if build_backend_kylin_static "$LINUX_BIN" && verify_kylin_binary "$LINUX_BIN"; then
    BUILD_METHOD="kylin-glibc228-static-build"
  else
    for fb in \
      "$ROOT/patches/ism-server-kylin-glibc228/ism_server" \
      "$ROOT/patches/ism-server-cgo-linux-amd64/ism_server" \
      "$ROOT/releases/ism-release-oceanbase-20260706/ism_server_user/ism_server" \
      "$ROOT/releases/ism-test-20260703/ism_server_user/ism_server"; do
      if [[ -f "$fb" ]] && verify_kylin_binary "$fb"; then
        cp "$fb" "$LINUX_BIN"
        chmod 755 "$LINUX_BIN"
        BUILD_METHOD="fallback-local-linux-amd64"
        echo "  使用本地 ELF: $fb"
        break
      fi
    done
  fi
fi
if [[ ! -f "$LINUX_BIN" ]] || ! verify_kylin_binary "$LINUX_BIN"; then
  echo "错误: 未得到麒麟 V10 兼容 ism_server（须 glibc<=2.28 或静态链接）"
  echo "  请执行: bash scripts/build_kylin_ism_server.sh"
  exit 1
fi

if [[ "${SKIP_FRONTEND_BUILD:-0}" == "1" ]] && [[ -f "$FRONTEND_DIST/index.html" ]]; then
  echo "[2/6] 跳过前端构建 (SKIP_FRONTEND_BUILD=1) ..."
else
  echo "[2/6] 编译前端 dist ..."
  rm -rf "$FRONTEND_DIST"
  (cd "$ROOT/ism-front-end-v2" && NODE_OPTIONS="--max-old-space-size=20480 --openssl-legacy-provider" npm run build)
fi
[[ -f "$FRONTEND_DIST/index.html" ]] || { echo "错误: dist/index.html 不存在"; exit 1; }
if rg -q 'sockjs-node|webpack-dev-server/client' "$FRONTEND_DIST/static/js/"*.js 2>/dev/null; then
  echo "错误: dist 含 dev-server 热更新代码（多为 vue-cli-service serve 污染）"
  echo "  请确认未在 serve 运行时打包，并重新 npm run build"
  exit 1
fi

echo "[3/6] 组装目录与 OceanBase 配置 ..."
if [[ -f "$LINUX_BIN" ]]; then chmod 755 "$LINUX_BIN"; fi

rsync -a --delete \
  --exclude 'vendor/' --exclude '*.go' --exclude 'logs/' \
  --exclude 'data/dbbackup/' --exclude 'data/tempDir/' --exclude 'data/upload/' \
  "$BACKEND_SRC/conf/" "$STAGING/ism_server_user/conf/"
rsync -a "$BACKEND_SRC/static/" "$STAGING/ism_server_user/static/" 2>/dev/null || mkdir -p "$STAGING/ism_server_user/static"
rsync -a "$BACKEND_SRC/data/auth/" "$STAGING/ism_server_user/data/auth/" 2>/dev/null || mkdir -p "$STAGING/ism_server_user/data/auth"
rsync -a "$FRONTEND_DIST/" "$STAGING/web/dist/"

cp "$ROOT/scripts/serve_test_frontend.py" "$STAGING/scripts/"
cp "$ROOT/scripts/modbus_simulator.py" "$STAGING/scripts/" 2>/dev/null || true
cp "$ROOT/scripts/import_mysql_to_oceanbase.sh" "$STAGING/scripts/"
cp "$ROOT/scripts/fix_ob_charset_reimport.sh" "$STAGING/scripts/"
cp "$ROOT/scripts/verify_ob_charset.sh" "$STAGING/scripts/"
cp "$ROOT/scripts/fix_oceanbase_schema_alarm_on_value.sh" "$STAGING/scripts/"
cp "$ROOT/scripts/fix_device_real_data_index.sh" "$STAGING/scripts/"
cp "$ROOT/scripts/fix_device_real_data_indexes_oceanbase.sql" "$STAGING/scripts/" 2>/dev/null || true
cp "$ROOT/scripts/prepare_mysql_dump_for_oceanbase.sh" "$STAGING/scripts/"
cp "$ROOT/scripts/migrate_sqlite_to_oceanbase.py" "$STAGING/scripts/"
cp "$ROOT/scripts/export_db_to_sqlite.py" "$STAGING/scripts/" 2>/dev/null || true
cp "$ROOT/scripts/clear_all_alarms.py" "$STAGING/scripts/" 2>/dev/null || true
cp "$ROOT/scripts/prune_legacy_dashboard_pages.py" "$STAGING/scripts/"
cp "$ROOT/scripts/prune_legacy_dashboard_pages_on_start.sh" "$STAGING/scripts/"
cp "$ROOT/scripts/install_docker_kylin_sp3.sh" "$STAGING/scripts/"
cp "$ROOT/scripts/ensure_docker_log_limits.sh" "$STAGING/scripts/"
cp "$ROOT/scripts/ensure_python.sh" "$STAGING/scripts/"
cp "$ROOT/scripts/install_python_kylin_sp3.sh" "$STAGING/scripts/"
cp "$ROOT/scripts/fix_compose_offline.sh" "$STAGING/scripts/"
chmod +x "$STAGING/scripts/install_docker_kylin_sp3.sh" \
  "$STAGING/scripts/ensure_docker_log_limits.sh" \
  "$STAGING/scripts/ensure_python.sh" \
  "$STAGING/scripts/install_python_kylin_sp3.sh" \
  "$STAGING/scripts/fix_compose_offline.sh" \
  "$STAGING/scripts/prune_legacy_dashboard_pages_on_start.sh" \
  "$STAGING/scripts/prune_legacy_dashboard_pages.py" \
  "$STAGING/scripts/fix_device_real_data_index.sh"
cp "$ROOT/docs/ISM-OceanBase部署与切换指南.md" "$STAGING/docs-ISM-OceanBase部署与切换指南.md" 2>/dev/null || true
if [[ -f "$ROOT/docs/ISM-麒麟V10-OceanBase部署操作手册.md" ]]; then
  cp "$ROOT/docs/ISM-麒麟V10-OceanBase部署操作手册.md" "$STAGING/"
fi
if [[ -f "$ROOT/docs/ISM-麒麟V10-OceanBase部署操作手册.pdf" ]]; then
  cp "$ROOT/docs/ISM-麒麟V10-OceanBase部署操作手册.pdf" "$STAGING/"
fi

APP_CONF="$STAGING/ism_server_user/conf/app.conf"
if [[ "$(uname -s)" == "Darwin" ]]; then
  sed -i '' "s/^dbtype=.*/dbtype=4/" "$APP_CONF"
  sed -i '' "s/^httpport=.*/httpport=${DEFAULT_BE_PORT}/" "$APP_CONF"
  sed -i '' "s/^oceanbasehost=.*/oceanbasehost=127.0.0.1/" "$APP_CONF"
  sed -i '' "s/^oceanbaseport=.*/oceanbaseport=${DEFAULT_OB_PORT}/" "$APP_CONF"
  sed -i '' "s/^oceanbaseuser=.*/oceanbaseuser=root@ism_tenant/" "$APP_CONF"
  sed -i '' "s/^oceanbasedbname=.*/oceanbasedbname=ism/" "$APP_CONF"
else
  sed -i "s/^dbtype=.*/dbtype=4/" "$APP_CONF"
  sed -i "s/^httpport=.*/httpport=${DEFAULT_BE_PORT}/" "$APP_CONF"
  sed -i "s/^oceanbasehost=.*/oceanbasehost=127.0.0.1/" "$APP_CONF"
  sed -i "s/^oceanbaseport=.*/oceanbaseport=${DEFAULT_OB_PORT}/" "$APP_CONF"
  sed -i "s/^oceanbaseuser=.*/oceanbaseuser=root@ism_tenant/" "$APP_CONF"
  sed -i "s/^oceanbasedbname=.*/oceanbasedbname=ism/" "$APP_CONF"
fi

cat > "$STAGING/ports.env" << PORTEOF
# ISM OceanBase 一体包默认端口（与客户生产错开）
ISM_FE_PORT=${DEFAULT_FE_PORT}
ISM_BE_PORT=${DEFAULT_BE_PORT}
OB_PORT=${DEFAULT_OB_PORT}
OB_TENANT=ism_tenant
OB_PASSWORD=ism2024!
OB_DATABASE=ism

# TDengine 历史库（REST，ISM historyrecorddbtype=2）
TD_PORT=${DEFAULT_TD_PORT}
TD_NATIVE_PORT=${DEFAULT_TD_NATIVE_PORT}
TD_USER=root
TD_PASSWORD=taosdata
TD_CONTAINER=tdengine
TD_IMAGE=${TD_IMAGE}
PORTEOF

cat > "$STAGING/docker-compose.oceanbase.yml" << 'COMPOSEEOF'
services:
  oceanbase:
    image: oceanbase/oceanbase-ce:latest
    container_name: oceanbase
    restart: unless-stopped
    ulimits:
      nofile:
        soft: 65536
        hard: 65536
      nproc:
        soft: 65536
        hard: 65536
    ports:
      - "${OB_PORT:-2881}:2881"
    logging:
      driver: json-file
      options:
        max-size: "500m"
        max-file: "20"
    environment:
      MODE: mini
      OB_MEMORY_LIMIT: 8G
      OB_DATAFILE_SIZE: 10G
      OB_LOG_DISK_SIZE: 5G
      OB_CLUSTER_NAME: ism_cluster
      OB_TENANT_NAME: ${OB_TENANT:-ism_tenant}
      OB_TENANT_PASSWORD: ${OB_PASSWORD:-ism2024!}
COMPOSEEOF

cat > "$STAGING/docker-compose.tdengine.yml" << 'TDCOMPOSEEOF'
services:
  tdengine:
    image: tdengine/tdengine:3.3.6.13
    container_name: tdengine
    restart: unless-stopped
    hostname: tdengine
    environment:
      TAOS_FQDN: localhost
    ports:
      - "${TD_PORT:-6041}:6041"
      - "${TD_NATIVE_PORT:-6030}:6030"
    logging:
      driver: json-file
      options:
        max-size: "500m"
        max-file: "20"
    volumes:
      - tdengine-data:/var/lib/taos
      - tdengine-log:/var/log/taos

volumes:
  tdengine-data:
  tdengine-log:
TDCOMPOSEEOF

# 确保历史库默认指向本机 TDengine
HISTORY_CONF="$STAGING/ism_server_user/conf/historyData.conf"
if [[ -f "$HISTORY_CONF" ]]; then
  if [[ "$(uname -s)" == "Darwin" ]]; then
    sed -i '' "s/^historyrecorddbtype=.*/historyrecorddbtype=2/" "$HISTORY_CONF"
    sed -i '' "s/^tdenginehost=.*/tdenginehost=127.0.0.1/" "$HISTORY_CONF"
    sed -i '' "s/^tdengineport=.*/tdengineport=${DEFAULT_TD_PORT}/" "$HISTORY_CONF"
  else
    sed -i "s/^historyrecorddbtype=.*/historyrecorddbtype=2/" "$HISTORY_CONF"
    sed -i "s/^tdenginehost=.*/tdenginehost=127.0.0.1/" "$HISTORY_CONF"
    sed -i "s/^tdengineport=.*/tdengineport=${DEFAULT_TD_PORT}/" "$HISTORY_CONF"
  fi
fi

cat > "$STAGING/scripts/init_oceanbase.sh" << 'INITOB'
#!/bin/bash
# 等待 OceanBase 就绪并创建 ism 库（首次部署）
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
[[ -f "$ROOT/ports.env" ]] && source "$ROOT/ports.env"

OB_PORT="${OB_PORT:-2881}"
OB_TENANT="${OB_TENANT:-ism_tenant}"
OB_PASSWORD="${OB_PASSWORD:-ism2024!}"
OB_DATABASE="${OB_DATABASE:-ism}"

echo "等待 OceanBase 端口 ${OB_PORT} ..."
for i in $(seq 1 90); do
  if docker exec oceanbase obclient -h 127.0.0.1 -P 2881 -uroot@${OB_TENANT} -p"${OB_PASSWORD}" -e "SELECT 1" >/dev/null 2>&1; then
    echo "OceanBase 已就绪"
    docker exec oceanbase obclient -h 127.0.0.1 -P 2881 -uroot@${OB_TENANT} -p"${OB_PASSWORD}" -e \
      "CREATE DATABASE IF NOT EXISTS ${OB_DATABASE} DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;"
    echo "数据库 ${OB_DATABASE} 已就绪"
    exit 0
  fi
  sleep 2
done
echo "错误: OceanBase 启动超时，请检查 docker logs oceanbase"
exit 1
INITOB

chmod +x "$STAGING/scripts/init_oceanbase.sh" \
  "$STAGING/scripts/import_mysql_to_oceanbase.sh" \
  "$STAGING/scripts/prepare_mysql_dump_for_oceanbase.sh"

echo "[4/6] 复制 MySQL 业务备份 ..."
mkdir -p "$STAGING/data/source"
if [[ ! -f "$MYSQL_BACKUP" ]]; then
  echo "错误: 缺少权威 MySQL 备份: $MYSQL_BACKUP"
  exit 1
fi
cp "$MYSQL_BACKUP" "$STAGING/data/source/"
rm -f "$STAGING/data/source/ism.db" "$STAGING/data/source/ism.db-shm" "$STAGING/data/source/ism.db-wal"
echo "  已复制: $(du -sh "$STAGING/data/source/$(basename "$MYSQL_BACKUP")" | cut -f1)"

echo "[4b/6] 打包 Docker + Docker Compose 离线安装组件 ..."
if [[ "${SKIP_DOCKER_BUNDLE:-0}" != "1" ]]; then
  bash "$ROOT/scripts/download_docker_offline_bundle.sh" || {
    echo "  警告: Docker 离线包下载失败，请检查网络后重试"
    exit 1
  }
fi
if [[ -d "$ROOT/docker-offline/bin" ]] && [[ -x "$ROOT/docker-offline/bin/dockerd" ]]; then
  rsync -a "$ROOT/docker-offline/" "$STAGING/docker-offline/"
  echo "  已打入: docker-offline/ ($(du -sh "$STAGING/docker-offline" | cut -f1))"
else
  echo "错误: 缺少 docker-offline/，完全离线包必须含 Docker + Compose"
  exit 1
fi

echo "[4c/6] 打包 Python 3 离线组件 ..."
if [[ "${SKIP_PYTHON_BUNDLE:-0}" != "1" ]]; then
  bash "$ROOT/scripts/download_python_offline_bundle.sh" || {
    echo "  警告: Python 离线包下载失败，请检查网络后重试"
    exit 1
  }
fi
if [[ -x "$ROOT/python-offline/install/bin/python3" ]] || [[ -x "$ROOT/python-offline/install/bin/python3.11" ]]; then
  rsync -a "$ROOT/python-offline/" "$STAGING/python-offline/"
  echo "  已打入: python-offline/ ($(du -sh "$STAGING/python-offline" | cut -f1))"
else
  echo "错误: 缺少 python-offline/install/bin/python3，完全离线包必须含 Python"
  exit 1
fi

echo "[5a/7] 导出 TDengine Docker 镜像（linux/amd64，离线历史库）..."
mkdir -p "$STAGING/tdengine"
TD_TAR="$STAGING/tdengine/tdengine.tar"
PREV_TD_TAR=$(ls -t "$ROOT"/releases/ism-release-oceanbase-*/tdengine/tdengine.tar "$ROOT"/releases/tdengine-offline/tdengine.tar 2>/dev/null | head -1 || true)
if [[ -n "${PREV_TD_TAR:-}" && -f "$PREV_TD_TAR" && -s "$PREV_TD_TAR" ]]; then
  cp "$PREV_TD_TAR" "$TD_TAR"
  if [[ -f "$(dirname "$PREV_TD_TAR")/IMAGE_TAG" ]]; then
    cp "$(dirname "$PREV_TD_TAR")/IMAGE_TAG" "$STAGING/tdengine/IMAGE_TAG"
  else
    printf '%s\n' "$TD_IMAGE" > "$STAGING/tdengine/IMAGE_TAG"
  fi
  echo "  已复用 TDengine tar: $PREV_TD_TAR ($(du -sh "$TD_TAR" | cut -f1))"
else
  bash "$ROOT/scripts/download_tdengine_offline_image.sh" "$STAGING/tdengine" || {
    echo "  错误: TDengine 离线镜像导出失败（历史库为出厂默认，必须打入包内）"
    exit 1
  }
fi
cp "$ROOT/scripts/init_tdengine.sh" "$STAGING/scripts/init_tdengine.sh"
chmod +x "$STAGING/scripts/init_tdengine.sh"

echo "[5/7] 导出 OceanBase Docker 镜像（linux/amd64，需网络）..."
# 多架构镜像在 Mac arm64 上须按 amd64 digest 拉取，否则 docker save 会因缺 arm64 层失败
OB_AMD64_DIGEST="sha256:46eefcf1275beae76947aac8f96c47efb745f6902d584de3c26465e4a9f113ae"
OB_PULL_SOURCES=(
  "ghcr.io/oceanbase/oceanbase-ce@${OB_AMD64_DIGEST}"
  "docker.1ms.run/oceanbase/oceanbase-ce@${OB_AMD64_DIGEST}"
  "docker.m.daocloud.io/oceanbase/oceanbase-ce@${OB_AMD64_DIGEST}"
)
if command -v docker >/dev/null 2>&1; then
  pulled=""
  for src in "${OB_PULL_SOURCES[@]}"; do
    if docker pull "$src" >/dev/null 2>&1; then
      docker tag "$src" "$OB_IMAGE"
      pulled=1
      echo "  已拉取: $src"
      break
    fi
  done
  if [[ -n "$pulled" ]]; then
    docker save "$OB_IMAGE" -o "$OB_TAR"
    echo "  已导出: $(du -sh "$OB_TAR" | cut -f1) (linux/amd64)"
  else
    echo "  警告: 无法拉取 oceanbase-ce (amd64)，尝试复用已有 release 镜像 tar ..."
    PREV_OB_TAR=$(ls -t "$ROOT"/releases/ism-release-oceanbase-*/oceanbase/oceanbase-ce.tar 2>/dev/null | head -1)
    if [[ -n "$PREV_OB_TAR" && -f "$PREV_OB_TAR" ]]; then
      cp "$PREV_OB_TAR" "$OB_TAR"
      echo "  已复用: $PREV_OB_TAR"
    else
      echo "  包内不含镜像 tar；现场需自行 docker load 或 docker pull"
      rm -f "$OB_TAR"
    fi
  fi
else
  echo "  跳过镜像导出（无 Docker）"
  PREV_OB_TAR=$(ls -t "$ROOT"/releases/ism-release-oceanbase-*/oceanbase/oceanbase-ce.tar 2>/dev/null | head -1)
  if [[ -n "$PREV_OB_TAR" && -f "$PREV_OB_TAR" ]]; then
    cp "$PREV_OB_TAR" "$OB_TAR"
    echo "  已复用已有 tar: $PREV_OB_TAR"
  fi
fi

OB_PRELOADED_TAR="$STAGING/oceanbase/oceanbase-ce-preloaded.tar"
OB_PRELOADED_TAG="oceanbase-ce-ism-preloaded:${DATE_TAG}"
OB_PRELOADED="no"

preload_oceanbase_data() {
  if [[ "$PRELOAD_OB_DATA" != "1" ]]; then
    echo "[5b/7] 跳过 OB 数据预导入 (PRELOAD_OB_DATA=$PRELOAD_OB_DATA)"
    return 0
  fi
  if ! command -v docker >/dev/null 2>&1; then
    echo "[5b/7] 跳过 OB 数据预导入（无 Docker）"
    return 0
  fi
  if [[ ! -f "$OB_TAR" ]] && ! docker image inspect "$OB_IMAGE" >/dev/null 2>&1; then
    echo "[5b/7] 跳过 OB 数据预导入（无 OB 镜像）"
    return 0
  fi

  echo "[5b/7] 预导入业务数据到 OceanBase（约 10–20 分钟，请耐心等待）..."
  docker rm -f oceanbase-build-preload 2>/dev/null || true
  docker rm -f oceanbase 2>/dev/null || true

  if [[ -f "$OB_TAR" ]]; then
    docker load -i "$OB_TAR"
  fi

  export OB_PORT="$DEFAULT_OB_PORT" OB_TENANT=ism_tenant OB_PASSWORD='ism2024!'
  # 构建机与麒麟现场均优先 docker run（避免 compose 段错误/不可用）
  ob_start_container "$DEFAULT_OB_PORT" "ism_tenant" "ism2024!"

  (cd "$STAGING" && bash scripts/init_oceanbase.sh)
  (cd "$STAGING" && bash scripts/import_mysql_to_oceanbase.sh)

  local user_cnt proj_cnt
  user_cnt="$(docker exec oceanbase obclient -h127.0.0.1 -P2881 -uroot@ism_tenant -p'ism2024!' ism \
    -N -e "SELECT COUNT(*) FROM user;" 2>/dev/null || echo 0)"
  proj_cnt="$(docker exec oceanbase obclient -h127.0.0.1 -P2881 -uroot@ism_tenant -p'ism2024!' ism \
    -N -e "SELECT COUNT(*) FROM project_lists;" 2>/dev/null || echo 0)"
  echo "  预导入验证: user=${user_cnt} project_lists=${proj_cnt}" | tee "$STAGING/data/.preload_verify.txt"
  if [[ "${user_cnt:-0}" -lt 1 ]]; then
    echo "  错误: 预导入失败（user 表为空），不生成 preloaded 镜像"
    echo "  提示: Mac arm64 上 OB 模拟常超时；请在麒麟/Linux amd64 上执行:"
    echo "    bash scripts/export_preloaded_ob.sh <部署目录>"
    docker rm -f oceanbase 2>/dev/null || true
    return 1
  fi

  docker commit oceanbase "$OB_PRELOADED_TAG"
  docker save "$OB_PRELOADED_TAG" -o "$OB_PRELOADED_TAR"
  echo "$OB_PRELOADED_TAG" > "$STAGING/oceanbase/PRELOADED_IMAGE_TAG"
  echo "  已导出预导入镜像: $(du -sh "$OB_PRELOADED_TAR" | cut -f1)"

  docker rm -f oceanbase 2>/dev/null || true
  OB_PRELOADED="yes"
  touch "$STAGING/.data_preloaded"
}

preload_oceanbase_data || {
  echo "  警告: OB 数据预导入失败，包内不含有效 preloaded 镜像"
  echo "  现场首次 start-all.sh 将自动导入 SQL（约 10–15 分钟）"
  echo "  或在麒麟测试机执行: bash scripts/export_preloaded_ob.sh <本目录> 后重新 zip"
  rm -f "$OB_PRELOADED_TAR" "$STAGING/.data_preloaded" "$STAGING/oceanbase/PRELOADED_IMAGE_TAG"
  OB_PRELOADED="no"
}

# 从仓库已验证脚本复制（优先最新一体包）
PKG_SCRIPTS_SRC=""
for cand in \
  "$ROOT/releases/ism-release-oceanbase-20260721-1024-dbb9" \
  "$ROOT/releases/ism-release-oceanbase-20260716-1747-fb6d" \
  "$ROOT/releases/ism-release-oceanbase-20260714-1746-3589" \
  "$ROOT/releases/ism-release-oceanbase-20260709"; do
  if [[ -d "$cand" && -f "$cand/start-all.sh" ]]; then
    PKG_SCRIPTS_SRC="$cand"
    break
  fi
done
[[ -n "$PKG_SCRIPTS_SRC" ]] || PKG_SCRIPTS_SRC="$ROOT/releases/ism-release-oceanbase-20260708"
for script in start-all.sh stop-all.sh deploy-offline.sh; do
  if [[ -f "$ROOT/scripts/$script" ]]; then
    cp "$ROOT/scripts/$script" "$STAGING/$script"
  elif [[ -f "$PKG_SCRIPTS_SRC/$script" ]]; then
    cp "$PKG_SCRIPTS_SRC/$script" "$STAGING/$script"
  fi
done
for script in check_env_kylin.sh diagnose_kylin.sh test_on_kylin_host.sh export_preloaded_ob.sh \
              fix_admin_password_oceanbase.sh check_dw_device_loading.sh check_login_deep.sh \
              check_login_and_user.sh collect_diagnose_log.sh run_full_field_check.sh \
              init_tdengine.sh; do
  if [[ -f "$ROOT/scripts/$script" ]]; then
    cp "$ROOT/scripts/$script" "$STAGING/scripts/$script"
  elif [[ -f "$PKG_SCRIPTS_SRC/scripts/$script" ]]; then
    cp "$PKG_SCRIPTS_SRC/scripts/$script" "$STAGING/scripts/$script"
  fi
done

# 启动时自动硬删旧预生成大屏页：注入 start-all（后端启动前）
if [[ -f "$STAGING/start-all.sh" ]] && ! grep -q 'prune_legacy_dashboard_pages_on_start' "$STAGING/start-all.sh"; then
  awk '
    BEGIN { done=0 }
    /启动后端 ism_server/ && !done {
      print "echo \"=== 硬删除旧大屏预生成页（启动清理）===\""
      print "if [[ -x \"$ROOT/scripts/prune_legacy_dashboard_pages_on_start.sh\" ]]; then"
      print "  bash \"$ROOT/scripts/prune_legacy_dashboard_pages_on_start.sh\" || true"
      print "fi"
      print ""
      done=1
    }
    { print }
  ' "$STAGING/start-all.sh" > "$STAGING/start-all.sh.tmp" \
    && mv "$STAGING/start-all.sh.tmp" "$STAGING/start-all.sh"
fi

# 启动时自动修复 device_real_data 索引（Error 4012 / COUNT 超时）：注入 start-all（后端启动前）
if [[ -f "$STAGING/start-all.sh" ]] && ! grep -q 'fix_device_real_data_index' "$STAGING/start-all.sh"; then
  awk '
    BEGIN { done=0 }
    /启动后端 ism_server/ && !done {
      print "echo \"=== device_real_data 索引自愈（VARCHAR + idx_drd_project_deleted）===\""
      print "if [[ -x \"$ROOT/scripts/fix_device_real_data_index.sh\" ]]; then"
      print "  bash \"$ROOT/scripts/fix_device_real_data_index.sh\" || true"
      print "else"
      print "  echo \"  [SKIP] 无 fix_device_real_data_index.sh\""
      print "fi"
      print ""
      done=1
    }
    { print }
  ' "$STAGING/start-all.sh" > "$STAGING/start-all.sh.tmp" \
    && mv "$STAGING/start-all.sh.tmp" "$STAGING/start-all.sh"
fi

chmod +x "$STAGING/start-all.sh" "$STAGING/stop-all.sh" "$STAGING/deploy-offline.sh" \
  "$STAGING/scripts/"*.sh 2>/dev/null || true
chmod +x "$STAGING/scripts/serve_test_frontend.py" \
  "$STAGING/scripts/prune_legacy_dashboard_pages.py" 2>/dev/null || true

cat > "$STAGING/README-部署说明.md" << READMEEOF
# ISM OceanBase + TDengine 一体部署包

- 版本: V3.01.RC07
- 构建日期: ${DATE_TAG}
- 业务库: **OceanBase** (\`dbtype=4\`)
- 历史库: **TDengine** (\`historyrecorddbtype=2\`，REST **${DEFAULT_TD_PORT}**)
- 默认端口: 前端 **${DEFAULT_FE_PORT}** / 后端 **${DEFAULT_BE_PORT}** / OceanBase **${DEFAULT_OB_PORT}** / TDengine **${DEFAULT_TD_PORT}**

## 目录结构

\`\`\`
${PKG_NAME}/
├── start-all.sh                 # 一键启动 OB + TDengine + 后端 + 前端
├── stop-all.sh
├── ports.env                    # 端口 / OB 租户 / TDengine 账号
├── docker-compose.oceanbase.yml
├── docker-compose.tdengine.yml
├── oceanbase/oceanbase-ce.tar   # OceanBase Docker 镜像
├── tdengine/tdengine.tar        # TDengine Docker 镜像（离线）
├── ism_server_user/             # 后端（dbtype=4，历史库默认 TDengine）
├── web/dist/                    # 前端静态资源
├── data/source/Mysql_Backup_*.sql
└── scripts/
    ├── init_oceanbase.sh
    ├── init_tdengine.sh
    └── import_mysql_to_oceanbase.sh
\`\`\`

## 一键部署（推荐）

1. **解压**到独立目录，例如 \`/opt/ISM/${PKG_NAME}/\`
2. **（可选）改端口**：编辑 \`ports.env\`
3. **一键启动**：
   - 首次、无 Docker：\`sudo bash deploy-offline.sh\`
   - 已有 Docker：\`sudo bash start-all.sh\`

\`start-all.sh\` 会自动：
1. 启动 OceanBase 并导入业务数据（首次）
2. \`docker load\` + 启动 TDengine，等待 REST **${DEFAULT_TD_PORT}** 就绪，预建 \`ISMHistoryDb\`
3. 写入 \`historyData.conf\`（\`historyrecorddbtype=2\` → 127.0.0.1:${DEFAULT_TD_PORT}）
4. 启动后端 / 前端

Web **系统参数 → 历史数据库** 选择 TDengine 即可直连，无需再装软件。

## app.conf / historyData.conf 关键项

\`\`\`ini
# app.conf
dbtype=4
oceanbaseuser=root@ism_tenant
oceanbasepwd=ism2024!
oceanbasehost=127.0.0.1
oceanbaseport=2881
oceanbasedbname=ism
httpport=8091

# historyData.conf
historyrecorddbtype=2
[tdengine]
tdenginehost=127.0.0.1
username=root
password=taosdata
tdengineport=6041
\`\`\`

详细文档见包内 \`docs-ISM-OceanBase部署与切换指南.md\`。
READMEEOF

cat > "$STAGING/BUILD_INFO.txt" << EOF
包名: ${PKG_NAME}-offline
构建 ID: ${BUILD_ID}
构建时间: $(date '+%Y-%m-%d %H:%M:%S')
后端编译: ${BUILD_METHOD}
dbtype: $(grep '^dbtype=' "$STAGING/ism_server_user/conf/app.conf")
historyrecorddbtype: $(grep -i '^historyrecorddbtype=' "$STAGING/ism_server_user/conf/historyData.conf" 2>/dev/null || echo N/A)
MySQL备份: $(du -sh "$STAGING/data/source/$(basename "$MYSQL_BACKUP")" 2>/dev/null | cut -f1 || echo N/A)
OB镜像: $([[ -f "$OB_TAR" ]] && du -sh "$OB_TAR" | cut -f1 || echo 未包含)
TDengine镜像: $([[ -f "$TD_TAR" ]] && du -sh "$TD_TAR" | cut -f1 || echo 未包含)
Docker离线包: $([[ -d "$STAGING/docker-offline/bin" ]] && du -sh "$STAGING/docker-offline" | cut -f1 || echo 未包含)
Python离线包: $([[ -d "$STAGING/python-offline/install" ]] && du -sh "$STAGING/python-offline" | cut -f1 || echo 未包含)
OB预导入: ${OB_PRELOADED}
OB预导入镜像: $([[ -f "$OB_PRELOADED_TAR" ]] && du -sh "$OB_PRELOADED_TAR" | cut -f1 || echo 未包含)
目标平台: 麒麟 V10 SP3 x86_64 (linux/amd64, glibc 2.28, 静态链接后端)
默认端口: FE=${DEFAULT_FE_PORT} BE=${DEFAULT_BE_PORT} OB=${DEFAULT_OB_PORT} TD=${DEFAULT_TD_PORT}
EOF

echo "[6/7] 环境检查脚本已随包复制 ..."

echo "[7/7] 压缩 zip ..."
mkdir -p "$ROOT/releases"
rm -f "$ZIP_OUT"
(cd "$ROOT/releases" && COPYFILE_DISABLE=1 zip -r -y -q "$(basename "$ZIP_OUT")" "$(basename "$STAGING")")

echo ""
echo "=== 构建完成 ==="
ls -lh "$ZIP_OUT"
echo "编译方式: $BUILD_METHOD"
