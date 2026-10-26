#!/bin/bash
# ISM 现场诊断日志采集：环境检查 + 服务状态 + API 探测，输出到 logs/ 便于 FinalShell 下载
# 用法:
#   cd /opt/ISM/ism-release-oceanbase-20260708
#   bash scripts/collect_diagnose_log.sh
# 或启动并等待后采集:
#   bash scripts/collect_diagnose_log.sh --start --wait 120
set -u

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
mkdir -p "$ROOT/logs"

TS="$(date '+%Y%m%d_%H%M%S')"
LOG="$ROOT/logs/ism_diagnose_${TS}.log"
ARCHIVE="$ROOT/logs/ism_diagnose_${TS}.tar.gz"

DO_START=0
WAIT_SEC=0
while [[ $# -gt 0 ]]; do
  case "$1" in
    --start) DO_START=1 ;;
    --wait)  WAIT_SEC="${2:-120}"; shift ;;
    -h|--help)
      echo "用法: bash scripts/collect_diagnose_log.sh [--start] [--wait 秒数]"
      echo "  --start   先执行 stop-all + start-all"
      echo "  --wait N  start 后等待 N 秒再采集（默认 0；与 --start 联用建议 120）"
      exit 0
      ;;
  esac
  shift
done

exec > >(tee -a "$LOG") 2>&1

echo "=============================================="
echo " ISM 诊断日志采集"
echo " 时间: $(date '+%F %T %z')"
echo " 主机: $(hostname -f 2>/dev/null || hostname)"
echo " 目录: $ROOT"
echo " 日志: $LOG"
echo "=============================================="
echo ""

if [[ -f "$ROOT/ports.env" ]]; then
  # shellcheck disable=SC1091
  source "$ROOT/ports.env"
fi
FE_PORT="${ISM_FE_PORT:-7090}"
BE_PORT="${ISM_BE_PORT:-8091}"
OB_PORT="${OB_PORT:-2881}"
OB_TENANT="${OB_TENANT:-ism_tenant}"
OB_PASSWORD="${OB_PASSWORD:-ism2024!}"

section() {
  echo ""
  echo "########## $* ##########"
  echo ""
}

if [[ "$DO_START" == "1" ]]; then
  section "启动服务 (stop-all + start-all)"
  bash "$ROOT/stop-all.sh" 2>/dev/null || true
  sleep 3
  bash "$ROOT/start-all.sh" || echo "[WARN] start-all.sh 退出码非 0"
  if [[ "$WAIT_SEC" -gt 0 ]]; then
    echo "等待 ${WAIT_SEC} 秒..."
    sleep "$WAIT_SEC"
  fi
fi

section "环境检查 check_env_kylin.sh"
if [[ -x "$ROOT/scripts/check_env_kylin.sh" ]]; then
  bash "$ROOT/scripts/check_env_kylin.sh" || true
else
  echo "  无 scripts/check_env_kylin.sh"
fi

section "快速诊断 diagnose_kylin.sh"
if [[ -x "$ROOT/scripts/diagnose_kylin.sh" ]]; then
  bash "$ROOT/scripts/diagnose_kylin.sh" || true
else
  echo "  无 scripts/diagnose_kylin.sh"
fi

section "端口与进程"
echo "--- ss 监听 ---"
ss -lntp 2>/dev/null | grep -E ":${OB_PORT}|:${BE_PORT}|:${FE_PORT} " || echo "  无 ${OB_PORT}/${BE_PORT}/${FE_PORT} 监听"
echo ""
echo "--- pgrep ---"
pgrep -fl 'ism_server|serve_test_frontend|docker-proxy' 2>/dev/null || echo "  无相关进程"
echo ""
echo "--- pid 文件 ---"
[[ -f "$ROOT/.backend.pid" ]] && echo "  .backend.pid=$(cat "$ROOT/.backend.pid")" || echo "  无 .backend.pid"
[[ -f "$ROOT/.frontend.pid" ]] && echo "  .frontend.pid=$(cat "$ROOT/.frontend.pid")" || echo "  无 .frontend.pid"

section "API 探测"
echo "--- 直连后端 login :${BE_PORT} ---"
LOGIN_RESP="$(curl -s -m 15 -X POST "http://127.0.0.1:${BE_PORT}/login" \
  -H 'Content-Type: application/json' \
  -d '{"Username":"admin","password":"e10adc3949ba59abbe56e057f20f883e"}' 2>&1)" || LOGIN_RESP="curl failed"
echo "$LOGIN_RESP"
echo ""
echo "--- 经前端代理 login :${FE_PORT}/api/login ---"
FE_LOGIN="$(curl -s --compressed -m 15 -X POST "http://127.0.0.1:${FE_PORT}/api/login" \
  -H 'Content-Type: application/json' \
  -H 'Accept-Encoding: gzip' \
  -d '{"Username":"admin","password":"e10adc3949ba59abbe56e057f20f883e"}' 2>&1)" || FE_LOGIN="curl failed"
echo "$FE_LOGIN"

section "数据库抽样"
docker exec oceanbase obclient --default-character-set=utf8mb4 \
  -h127.0.0.1 -P"${OB_PORT}" -uroot@"${OB_TENANT}" -p"${OB_PASSWORD}" ism -e "
SELECT COUNT(*) AS user_cnt FROM user;
SELECT id, username, HEX(username) AS username_hex, name, role, deleted_at FROM user LIMIT 10;
SELECT COUNT(*) AS project_cnt FROM project_lists;
SELECT COUNT(*) AS monitor_cnt FROM monitor_list;
SELECT COUNT(*) AS real_data_cnt FROM device_real_data;
SELECT id, name, HEX(name) AS name_hex FROM monitor_list WHERE id IN (771,772) LIMIT 5;
" 2>&1 || echo "  obclient 查询失败"

section "后端业务日志尾部 ism_server_user/logs/ism.log"
if [[ -f "$ROOT/ism_server_user/logs/ism.log" ]]; then
  tail -80 "$ROOT/ism_server_user/logs/ism.log"
else
  echo "  无 ism_server_user/logs/ism.log"
fi

section "打包附件（供 FinalShell 一次下载）"
TAR_LIST=()
TAR_LIST+=("logs/$(basename "$LOG")")
[[ -f "$ROOT/logs/ism_server.log" ]] && TAR_LIST+=("logs/ism_server.log")
[[ -f "$ROOT/logs/frontend.log" ]] && TAR_LIST+=("logs/frontend.log")
[[ -f "$ROOT/ports.env" ]] && TAR_LIST+=("ports.env")
if [[ -f "$ARCHIVE" ]]; then rm -f "$ARCHIVE"; fi
tar -czf "$ARCHIVE" -C "$ROOT" "${TAR_LIST[@]}" 2>/dev/null || true

echo ""
echo "=============================================="
echo " 采集完成"
echo " 主日志: $LOG"
[[ -f "$ARCHIVE" ]] && echo " 压缩包: $ARCHIVE ($(du -h "$ARCHIVE" | cut -f1))"
echo ""
echo " FinalShell 下载路径:"
echo "   $LOG"
[[ -f "$ARCHIVE" ]] && echo "   $ARCHIVE"
echo "=============================================="
