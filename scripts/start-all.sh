#!/bin/bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")" && pwd)"
[[ -f "$ROOT/ports.env" ]] && source "$ROOT/ports.env"

# Docker 容器日志上限（防 /var/lib/docker 撑满系统盘）
# shellcheck disable=SC1091
source "$ROOT/scripts/ensure_docker_log_limits.sh"
if [[ "$(id -u)" -eq 0 ]]; then
  apply_daemon_json || true
fi

FE_PORT="${ISM_FE_PORT:-7090}"
BE_PORT="${ISM_BE_PORT:-8091}"
OB_PORT="${OB_PORT:-2881}"
TD_PORT="${TD_PORT:-6041}"
TD_NATIVE_PORT="${TD_NATIVE_PORT:-6030}"
TD_USER="${TD_USER:-root}"
TD_PASSWORD="${TD_PASSWORD:-taosdata}"
TD_CONTAINER="${TD_CONTAINER:-tdengine}"
TD_IMAGE="${TD_IMAGE:-tdengine/tdengine:3.3.6.13}"
CONF="$ROOT/ism_server_user/conf/app.conf"
HISTORY_CONF="$ROOT/ism_server_user/conf/historyData.conf"
ISM_PYTHON="${ISM_PYTHON:-$(bash "$ROOT/scripts/ensure_python.sh")}"

mkdir -p "$ROOT/logs" "$ROOT/ism_server_user/data/sessionon"

if [[ "$(uname -s)" == "Darwin" ]]; then
  sed -i '' "s/^httpport=.*/httpport=${BE_PORT}/" "$CONF"
  sed -i '' "s/^oceanbaseport=.*/oceanbaseport=${OB_PORT}/" "$CONF"
  if [[ -f "$HISTORY_CONF" ]]; then
    sed -i '' "s/^tdengineport=.*/tdengineport=${TD_PORT}/" "$HISTORY_CONF"
    sed -i '' "s/^tdenginehost=.*/tdenginehost=127.0.0.1/" "$HISTORY_CONF"
    sed -i '' "s/^historyrecorddbtype=.*/historyrecorddbtype=2/" "$HISTORY_CONF"
  fi
else
  sed -i "s/^httpport=.*/httpport=${BE_PORT}/" "$CONF"
  sed -i "s/^oceanbaseport=.*/oceanbaseport=${OB_PORT}/" "$CONF"
  if [[ -f "$HISTORY_CONF" ]]; then
    sed -i "s/^tdengineport=.*/tdengineport=${TD_PORT}/" "$HISTORY_CONF"
    sed -i "s/^tdenginehost=.*/tdenginehost=127.0.0.1/" "$HISTORY_CONF"
    sed -i "s/^historyrecorddbtype=.*/historyrecorddbtype=2/" "$HISTORY_CONF"
  fi
fi

ob_table_count() {
  docker exec oceanbase obclient -h127.0.0.1 -P2881 \
    -uroot@"${OB_TENANT:-ism_tenant}" -p"${OB_PASSWORD:-ism2024!}" ism \
    -N -e "SELECT COUNT(*) FROM user;" 2>/dev/null || echo "0"
}

if ! command -v docker >/dev/null 2>&1 || ! docker info >/dev/null 2>&1; then
  echo "错误: Docker 未安装或未启动"
  echo "  完全离线环境请先执行: sudo bash deploy-offline.sh"
  echo "  或仅安装 Docker: sudo bash scripts/install_docker_kylin_sp3.sh"
  exit 1
fi

ob_start_container() {
  docker rm -f oceanbase 2>/dev/null || true
  docker run -d --name oceanbase --restart unless-stopped \
    "${ISM_DOCKER_LOG_OPTS[@]}" \
    --ulimit nofile=65536:65536 --ulimit nproc=65536:65536 \
    -p "${OB_PORT}:2881" \
    -e MODE=mini -e OB_MEMORY_LIMIT=8G -e OB_DATAFILE_SIZE=10G -e OB_LOG_DISK_SIZE=5G \
    -e OB_CLUSTER_NAME=ism_cluster \
    -e OB_TENANT_NAME="${OB_TENANT:-ism_tenant}" \
    -e OB_TENANT_PASSWORD="${OB_PASSWORD:-ism2024!}" \
    oceanbase/oceanbase-ce:latest
}

resolve_td_image() {
  # 优先 ports.env 指定名，其次包内 IMAGE_TAG，再尝试镜像加速前缀名
  local candidates=()
  [[ -n "${TD_IMAGE:-}" ]] && candidates+=("$TD_IMAGE")
  if [[ -f "$ROOT/tdengine/IMAGE_TAG" ]]; then
    while IFS= read -r line || [[ -n "$line" ]]; do
      [[ -n "$line" ]] && candidates+=("$line")
    done < "$ROOT/tdengine/IMAGE_TAG"
  fi
  candidates+=("tdengine/tdengine:3.3.6.13" "docker.1ms.run/tdengine/tdengine:3.3.6.13")
  local c
  for c in "${candidates[@]}"; do
    if docker image inspect "$c" >/dev/null 2>&1; then
      echo "$c"
      return 0
    fi
  done
  echo "${TD_IMAGE:-tdengine/tdengine:3.3.6.13}"
}

td_start_container() {
  local img
  img="$(resolve_td_image)"
  docker rm -f "${TD_CONTAINER}" 2>/dev/null || true
  docker run -d --name "${TD_CONTAINER}" --restart unless-stopped \
    "${ISM_DOCKER_LOG_OPTS[@]}" \
    --hostname tdengine \
    -e TAOS_FQDN=localhost \
    -p "${TD_PORT}:6041" \
    -p "${TD_NATIVE_PORT}:6030" \
    -v "${ROOT}/tdengine/data:/var/lib/taos" \
    -v "${ROOT}/tdengine/log:/var/log/taos" \
    "$img"
}

ob_compose_up() {
  # 麒麟现场常见：独立 docker-compose 段错误，仅使用 compose 插件或 docker run
  local compose_file="$1"
  if docker compose version >/dev/null 2>&1; then
    OB_PORT="$OB_PORT" OB_TENANT="${OB_TENANT:-ism_tenant}" OB_PASSWORD="${OB_PASSWORD:-ism2024!}" \
      docker compose -f "$compose_file" up -d
  else
    echo "  compose 插件不可用，使用 docker run 启动 OceanBase（推荐离线环境）"
    ob_start_container
  fi
}

ensure_compose_plugin() {
  docker compose version >/dev/null 2>&1 && return 0
  if [[ -x "$ROOT/scripts/fix_compose_offline.sh" ]] && [[ "$(id -u)" -eq 0 ]]; then
    echo "  尝试安装包内 Compose 插件 ..."
    bash "$ROOT/scripts/fix_compose_offline.sh" >/dev/null 2>&1 || true
  fi
}

ensure_compose_plugin || true

echo "=== [1/6] 启动 OceanBase ==="
if ! docker ps --format '{{.Names}}' 2>/dev/null | grep -q '^oceanbase$'; then
  if [[ -f "$ROOT/.data_preloaded" && -f "$ROOT/oceanbase/oceanbase-ce-preloaded.tar" ]]; then
    echo "  加载预导入 OceanBase 镜像（含业务数据）..."
    docker load -i "$ROOT/oceanbase/oceanbase-ce-preloaded.tar"
    PRELOAD_IMG="$(cat "$ROOT/oceanbase/PRELOADED_IMAGE_TAG" 2>/dev/null || docker images --format '{{.Repository}}:{{.Tag}}' | grep 'oceanbase-ce-ism-preloaded' | head -1)"
    OB_PORT="$OB_PORT" OB_TENANT="${OB_TENANT:-ism_tenant}" OB_PASSWORD="${OB_PASSWORD:-ism2024!}" \
      docker run -d --name oceanbase --restart unless-stopped \
      "${ISM_DOCKER_LOG_OPTS[@]}" \
      --ulimit nofile=65536:65536 --ulimit nproc=65536:65536 \
      -p "${OB_PORT}:2881" \
      -e MODE=mini -e OB_MEMORY_LIMIT=8G -e OB_DATAFILE_SIZE=10G -e OB_LOG_DISK_SIZE=5G \
      -e OB_CLUSTER_NAME=ism_cluster -e OB_TENANT_NAME="${OB_TENANT:-ism_tenant}" \
      -e OB_TENANT_PASSWORD="${OB_PASSWORD:-ism2024!}" \
      "$PRELOAD_IMG"
  else
    if [[ -f "$ROOT/oceanbase/oceanbase-ce.tar" ]]; then
      docker load -i "$ROOT/oceanbase/oceanbase-ce.tar"
    fi
    # 离线环境优先 docker run，避免损坏的 docker-compose 段错误
    ob_start_container
  fi
fi
bash "$ROOT/scripts/init_oceanbase.sh"

echo "=== [2/6] 启动 TDengine（历史库 REST ${TD_PORT}）==="
mkdir -p "$ROOT/tdengine/data" "$ROOT/tdengine/log"
if ! docker ps --format '{{.Names}}' 2>/dev/null | grep -q "^${TD_CONTAINER}$"; then
  if ! docker image inspect "$(resolve_td_image)" >/dev/null 2>&1; then
    if [[ -f "$ROOT/tdengine/tdengine.tar" ]]; then
      echo "  加载离线 TDengine 镜像 ..."
      docker load -i "$ROOT/tdengine/tdengine.tar"
      # crane/加速源 load 后可能只有 docker.1ms.run/... 标签，补打标准名
      if ! docker image inspect "tdengine/tdengine:3.3.6.13" >/dev/null 2>&1; then
        ALT="$(docker images --format '{{.Repository}}:{{.Tag}}' | grep -E 'tdengine/(tdengine|tsdb):' | head -1 || true)"
        [[ -n "$ALT" ]] && docker tag "$ALT" "tdengine/tdengine:3.3.6.13" || true
      fi
    else
      echo "  错误: 缺少 tdengine/tdengine.tar，无法离线启动 TDengine"
      echo "  请将镜像放入 $ROOT/tdengine/tdengine.tar 后重试"
      exit 1
    fi
  fi
  td_start_container
fi
bash "$ROOT/scripts/init_tdengine.sh"

echo "=== [3/6] 检查/导入业务数据 ==="
USER_CNT="$(ob_table_count)"
if [[ "$USER_CNT" == "0" || "$USER_CNT" == "" ]]; then
  if [[ -f "$ROOT/.data_preloaded" ]]; then
    echo "  警告: 预导入镜像中 user 表为空，尝试从 SQL 重新导入..."
  else
    echo "  首次部署，自动导入 MySQL 业务数据（约 10–15 分钟）..."
  fi
  bash "$ROOT/scripts/import_mysql_to_oceanbase.sh"
  touch "$ROOT/.data_imported"
else
  if [[ -f "$ROOT/.data_preloaded" ]]; then
    echo "  预导入业务数据已就绪（user 表 ${USER_CNT} 条，无需再导入 SQL）"
  else
    echo "  业务数据已就绪（user 表 ${USER_CNT} 条）"
  fi
fi

echo "=== 硬删除旧大屏预生成页（启动清理）==="
if [[ -x "$ROOT/scripts/prune_legacy_dashboard_pages_on_start.sh" ]]; then
  bash "$ROOT/scripts/prune_legacy_dashboard_pages_on_start.sh" || true
fi

echo "=== device_real_data 索引自愈（VARCHAR + idx_drd_project_deleted）==="
if [[ -x "$ROOT/scripts/fix_device_real_data_index.sh" ]]; then
  bash "$ROOT/scripts/fix_device_real_data_index.sh" || true
else
  echo "  [SKIP] 无 fix_device_real_data_index.sh"
fi

echo "=== OceanBase 备份/还原调参（max_allowed_packet + ob_query_timeout）==="
if [[ -x "$ROOT/scripts/tune_ob_max_allowed_packet.sh" ]]; then
  OB_HOST=127.0.0.1 OB_PORT="${OB_PORT}" \
    OB_USER="root@${OB_TENANT:-ism_tenant}" OB_PASS="${OB_PASSWORD:-ism2024!}" OB_DB=ism \
    bash "$ROOT/scripts/tune_ob_max_allowed_packet.sh" || true
else
  echo "  [SKIP] 无 tune_ob_max_allowed_packet.sh"
fi

echo "=== [4/6] 启动后端 ism_server (端口 ${BE_PORT}) ==="
if [[ ! -x "$ROOT/ism_server_user/ism_server" ]]; then
  echo "错误: 缺少 ism_server_user/ism_server"
  exit 1
fi
wait_port() {
  local port="$1" label="$2" tries="${3:-20}"
  for _ in $(seq 1 "$tries"); do
    if ss -lnt 2>/dev/null | grep -qE ":${port}[[:space:]]"; then
      echo "  ${label} 端口 ${port} 已监听"
      return 0
    fi
    sleep 1
  done
  echo "  错误: ${label} 端口 ${port} 未监听"
  return 1
}

# 后端首次启动会 AutoMigrate（「检查系统表」），OceanBase 上常需 1–3 分钟才监听 HTTP。
# 旧逻辑只等 25s，进程还在跑就被误判失败；第二次再跑往往已迁移完所以「秒起」。
wait_backend_port() {
  local port="$1" tries="${2:-180}" pid_file="$3"
  local i=0
  echo "  等待后端监听 :${port}（最长 ${tries}s；首次可能在检查系统表）..."
  for i in $(seq 1 "$tries"); do
    if ss -lnt 2>/dev/null | grep -qE ":${port}[[:space:]]"; then
      echo "  后端 端口 ${port} 已监听（约 ${i}s）"
      return 0
    fi
    if [[ -n "$pid_file" && -f "$pid_file" ]] && ! kill -0 "$(cat "$pid_file")" 2>/dev/null; then
      echo "  错误: ism_server 进程已退出（等待第 ${i}s）"
      return 1
    fi
    if (( i % 15 == 0 )); then
      echo "  ...仍在启动（${i}/${tries}s），若日志停在「正在检查系统表」属正常，请继续等待"
    fi
    sleep 1
  done
  echo "  错误: 后端 端口 ${port} 未监听（已等 ${tries}s）"
  return 1
}

cd "$ROOT/ism_server_user"
if [[ -f "$ROOT/.backend.pid" ]] && kill -0 "$(cat "$ROOT/.backend.pid")" 2>/dev/null; then
  echo "后端已在运行"
else
  rm -f "$ROOT/.backend.pid"
  nohup ./ism_server > "$ROOT/logs/ism_server.log" 2>&1 &
  echo $! > "$ROOT/.backend.pid"
  sleep 2
  if ! kill -0 "$(cat "$ROOT/.backend.pid")" 2>/dev/null; then
    echo "  错误: ism_server 启动后立即退出，日志:"
    tail -30 "$ROOT/logs/ism_server.log" 2>/dev/null || true
    exit 1
  fi
fi
cd "$ROOT"
wait_backend_port "$BE_PORT" "${ISM_BE_WAIT_SECS:-180}" "$ROOT/.backend.pid" || {
  echo "  请执行: bash scripts/diagnose_kylin.sh"
  echo "  若日志仍停在「正在检查系统表」：进程可能还在迁移，可再等或调大 ISM_BE_WAIT_SECS 后重试"
  tail -30 "$ROOT/logs/ism_server.log" 2>/dev/null || true
  exit 1
}

echo "=== [5/6] 启动前端静态服务 (端口 ${FE_PORT}) ==="
if [[ ! -d "$ROOT/web/dist" ]]; then
  echo "错误: 缺少 web/dist 目录"
  exit 1
fi
if [[ -f "$ROOT/.frontend.pid" ]] && kill -0 "$(cat "$ROOT/.frontend.pid")" 2>/dev/null; then
  echo "前端已在运行"
else
  rm -f "$ROOT/.frontend.pid"
  nohup "$ISM_PYTHON" "$ROOT/scripts/serve_test_frontend.py" \
    --port "${FE_PORT}" \
    --dist "$ROOT/web/dist" \
    --backend "http://127.0.0.1:${BE_PORT}" \
    > "$ROOT/logs/frontend.log" 2>&1 &
  echo $! > "$ROOT/.frontend.pid"
  sleep 2
  if ! kill -0 "$(cat "$ROOT/.frontend.pid")" 2>/dev/null; then
    echo "  错误: 前端进程退出，日志:"
    tail -20 "$ROOT/logs/frontend.log" 2>/dev/null || true
    exit 1
  fi
fi
wait_port "$FE_PORT" "前端" 15 || {
  tail -20 "$ROOT/logs/frontend.log" 2>/dev/null || true
  exit 1
}

echo ""
echo "=== [6/6] 验证历史库配置 ==="
if [[ -f "$HISTORY_CONF" ]]; then
  echo "  historyData.conf: historyrecorddbtype=$(grep -i '^historyrecorddbtype=' "$HISTORY_CONF" | head -1 | cut -d= -f2)"
  echo "  TDengine: 127.0.0.1:${TD_PORT}  user=${TD_USER}"
fi
if curl -sf -u "${TD_USER}:${TD_PASSWORD}" -d "show databases;" "http://127.0.0.1:${TD_PORT}/rest/sql" >/dev/null 2>&1; then
  echo "  TDengine REST: OK"
else
  echo "  警告: TDengine REST 未响应，历史写入可能失败"
fi

echo ""
echo "=== ISM OceanBase + TDengine 一体环境已启动 ==="
echo "  OceanBase: 127.0.0.1:${OB_PORT}  租户 root@${OB_TENANT:-ism_tenant}"
echo "  TDengine:  127.0.0.1:${TD_PORT}  (REST，历史库)"
echo "  后端:      ${BE_PORT}  前端: ${FE_PORT}"
echo "  访问: http://<本机IP>:${FE_PORT}/#/login  账号 admin / 123456"
echo "  系统参数 → 历史数据库 → TDengine 即可直连本机 ${TD_PORT}"
echo "  停止: bash stop-all.sh"
