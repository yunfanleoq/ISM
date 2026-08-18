#!/bin/bash
# ISM Docker 日志限制：防止容器 json 日志撑满 /var
# - 全局：/etc/docker/daemon.json（需 root，重启 Docker 后对新默认容器生效）
# - 单容器：docker run 使用 ISM_DOCKER_LOG_OPTS（约 500m × 20 ≈ 10GB 上限）
#
# 用法:
#   sudo bash scripts/ensure_docker_log_limits.sh --apply-daemon
#   或在 start-all.sh 中 source 本脚本后使用 "${ISM_DOCKER_LOG_OPTS[@]}"
set -euo pipefail

ISM_DOCKER_LOG_MAX_SIZE="${ISM_DOCKER_LOG_MAX_SIZE:-500m}"
ISM_DOCKER_LOG_MAX_FILE="${ISM_DOCKER_LOG_MAX_FILE:-20}"

if [[ -z "${ISM_DOCKER_LOG_OPTS_LOADED:-}" ]]; then
  ISM_DOCKER_LOG_OPTS=(
    --log-driver json-file
    --log-opt "max-size=${ISM_DOCKER_LOG_MAX_SIZE}"
    --log-opt "max-file=${ISM_DOCKER_LOG_MAX_FILE}"
  )
  ISM_DOCKER_LOG_OPTS_LOADED=1
fi

apply_daemon_json() {
  if [[ "$(id -u)" -ne 0 ]]; then
    echo "  [SKIP] 写入 daemon.json 需要 root；start-all 仍会对 OB/TD 容器加 --log-opt"
    return 0
  fi

  mkdir -p /etc/docker
  local target="/etc/docker/daemon.json"
  local tmp
  tmp="$(mktemp)"
  cat > "$tmp" <<EOF
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "${ISM_DOCKER_LOG_MAX_SIZE}",
    "max-file": "${ISM_DOCKER_LOG_MAX_FILE}"
  }
}
EOF

  if [[ -f "$target" ]] && cmp -s "$tmp" "$target" 2>/dev/null; then
    echo "  Docker daemon.json 日志策略已是最新（${ISM_DOCKER_LOG_MAX_SIZE} × ${ISM_DOCKER_LOG_MAX_FILE}）"
    rm -f "$tmp"
    return 0
  fi

  if [[ -f "$target" ]]; then
    cp "$target" "${target}.bak.$(date +%Y%m%d%H%M%S)"
  fi
  mv "$tmp" "$target"
  echo "  已写入 ${target}（单容器日志上限约 ${ISM_DOCKER_LOG_MAX_SIZE} × ${ISM_DOCKER_LOG_MAX_FILE}）"
  if systemctl is-active docker >/dev/null 2>&1; then
    echo "  提示: 全局默认策略需重启 Docker 后生效: systemctl restart docker"
  fi
}

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  case "${1:---apply-daemon}" in
    --apply-daemon) apply_daemon_json ;;
    *)
      echo "用法: $0 [--apply-daemon]" >&2
      exit 1
      ;;
  esac
fi
