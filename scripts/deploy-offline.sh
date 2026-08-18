#!/bin/bash
# 麒麟 V10 SP3 完全离线一键部署：Docker（若无）→ OceanBase → ISM
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

echo "=== ISM 完全离线部署 ==="
echo "  平台: 银河麒麟 V10 SP3 x86_64"
echo ""

echo "[步骤 1/3] 检测 Python3 ..."
export ISM_PYTHON
ISM_PYTHON="$(bash "$ROOT/scripts/ensure_python.sh")"
echo "  Python: $ISM_PYTHON ($("$ISM_PYTHON" --version 2>&1))"

need_docker=0
if ! command -v docker >/dev/null 2>&1; then
  need_docker=1
elif ! docker info >/dev/null 2>&1; then
  need_docker=1
fi

if [[ "$need_docker" == "1" ]]; then
  echo "[步骤 2/3] 安装 Docker + Docker Compose（包内离线组件）..."
  if [[ ! -x "$ROOT/scripts/install_docker_kylin_sp3.sh" ]]; then
    echo "错误: 缺少 scripts/install_docker_kylin_sp3.sh"
    exit 1
  fi
  if [[ "$(id -u)" -ne 0 ]]; then
    echo "需要 root 权限安装 Docker，请执行:"
    echo "  sudo bash deploy-offline.sh"
    exit 1
  fi
  bash "$ROOT/scripts/install_docker_kylin_sp3.sh"
else
  echo "[步骤 2/3] Docker 已安装: $(docker --version)"
  if docker compose version >/dev/null 2>&1; then
    echo "  Compose: $(docker compose version 2>&1 | head -1)"
  else
    echo "  Compose: 插件未就绪（不影响启动，start-all.sh 使用 docker run）"
    if [[ -x "$ROOT/scripts/fix_compose_offline.sh" ]] && [[ "$(id -u)" -eq 0 ]]; then
      bash "$ROOT/scripts/fix_compose_offline.sh" || true
    fi
  fi
  if [[ -x "$ROOT/scripts/ensure_docker_log_limits.sh" ]]; then
    echo "  检查 Docker 日志轮转策略 ..."
    if [[ "$(id -u)" -eq 0 ]]; then
      bash "$ROOT/scripts/ensure_docker_log_limits.sh" --apply-daemon || true
    else
      echo "  提示: sudo bash scripts/ensure_docker_log_limits.sh 可写入全局 daemon.json"
    fi
  fi
fi

echo ""
echo "[步骤 3/3] 启动 ISM（OceanBase + TDengine + 后端 + 前端）..."
bash "$ROOT/start-all.sh"

echo ""
echo "=== 离线部署完成 ==="
echo "  访问: http://<本机IP>:7090/#/login  账号 admin / 123456"
echo "  历史库: TDengine REST 127.0.0.1:6041（系统参数中选 TDengine 即可）"
