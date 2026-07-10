#!/bin/bash
# ISM 现场一键排查（推荐）：重启服务 → 等待 → 全量诊断 + 登录检查 → 打压缩包
# 用法:
#   cd /opt/ISM/ism-release-oceanbase-20260708
#   bash scripts/run_full_field_check.sh
#
# 不重启、只采集:
#   bash scripts/run_full_field_check.sh --no-start
set -u

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
mkdir -p "$ROOT/logs"

TS="$(date '+%Y%m%d_%H%M%S')"
LOG="$ROOT/logs/ism_field_check_${TS}.log"
ARCHIVE="$ROOT/logs/ism_field_check_${TS}.tar.gz"

DO_START=1
WAIT_SEC=120
while [[ $# -gt 0 ]]; do
  case "$1" in
    --no-start) DO_START=0 ;;
    --wait)     WAIT_SEC="${2:-120}"; shift ;;
    -h|--help)
      echo "用法: bash scripts/run_full_field_check.sh [--no-start] [--wait 秒数]"
      echo "  默认: stop-all → start-all → 等 120 秒 → 全量诊断"
      echo "  --no-start  跳过重启，只采集当前状态"
      exit 0
      ;;
  esac
  shift
done

exec > >(tee -a "$LOG") 2>&1

echo "=============================================="
echo " ISM 现场一键排查"
echo " 时间: $(date '+%F %T %z')"
echo " 主机: $(hostname 2>/dev/null || echo unknown)"
echo " 目录: $ROOT"
echo " 主日志: $LOG"
echo "=============================================="

if [[ "$DO_START" == "1" ]]; then
  echo ""
  echo ">>> 停止旧服务..."
  bash "$ROOT/stop-all.sh" 2>/dev/null || true
  sleep 3
  echo ">>> 启动服务 (start-all.sh)..."
  bash "$ROOT/start-all.sh" || echo "[WARN] start-all 非零退出"
  echo ">>> 等待 ${WAIT_SEC} 秒（后端启动较慢，请勿中断）..."
  sleep "$WAIT_SEC"
fi

echo ""
echo ">>> 运行 collect_diagnose_log.sh ..."
bash "$ROOT/scripts/collect_diagnose_log.sh" 2>&1 || true

echo ""
echo ">>> 运行 check_login_and_user.sh ..."
bash "$ROOT/scripts/check_login_and_user.sh" 2>&1 || true

echo ""
echo ">>> 打包所有诊断日志 ..."
TAR_FILES=()
TAR_FILES+=("logs/$(basename "$LOG")")
for f in "$ROOT"/logs/ism_diagnose_*.log "$ROOT"/logs/ism_login_check_*.log; do
  [[ -f "$f" ]] && TAR_FILES+=("logs/$(basename "$f")")
done
[[ -f "$ROOT/logs/ism_server.log" ]] && TAR_FILES+=("logs/ism_server.log")
[[ -f "$ROOT/logs/frontend.log" ]] && TAR_FILES+=("logs/frontend.log")
[[ -f "$ROOT/ports.env" ]] && TAR_FILES+=("ports.env")

rm -f "$ARCHIVE"
tar -czf "$ARCHIVE" -C "$ROOT" "${TAR_FILES[@]}" 2>/dev/null || true

echo ""
echo "=============================================="
echo " 一键排查完成"
echo ""
echo " FinalShell 下载（任选其一）:"
echo "   压缩包（推荐）: $ARCHIVE"
echo "   主日志:         $LOG"
ls -lh "$ARCHIVE" "$LOG" 2>/dev/null || ls -lh "$LOG"
echo "=============================================="
