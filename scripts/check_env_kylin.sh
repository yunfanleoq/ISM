#!/bin/bash
echo "=== ISM 麒麟 V10 SP3 环境检查 ==="
echo "[系统] $(uname -a)"
echo "[CPU]  $(uname -m)  (需 x86_64 / amd64)"
echo "[内存] $(free -h 2>/dev/null | awk '/Mem:/{print $2" 可用 "$7}' || echo '请安装 free 命令')"
echo "[磁盘] $(df -h . | tail -1)"
if command -v python3 >/dev/null 2>&1 && python3 -c 'import sys; exit(0 if sys.version_info>=(3,6) else 1)' 2>/dev/null; then
  echo "[Python3] 系统: $(python3 --version)"
else
  echo "[Python3] 系统未安装 → 将使用包内 python-offline/"
fi
bash scripts/ensure_python.sh 2>/dev/null | sed 's|^|[Python3] 实际使用: |' || echo "[Python3] 检测失败"
command -v docker >/dev/null && echo "[Docker] $(docker --version)" || echo "[Docker] 未安装 → sudo bash scripts/install_docker_kylin_sp3.sh"
compose_ok=0
if docker compose version >/dev/null 2>&1; then
  echo "[Compose] $(docker compose version 2>&1 | head -1)"
  compose_ok=1
fi
if [[ "$compose_ok" == "0" ]]; then
  # 禁止执行 docker-compose（连字符）：麒麟现场该二进制常段错误，一运行就崩
  broken=""
  for dc in /usr/local/bin/docker-compose /usr/bin/docker-compose; do
    [[ -e "$dc" ]] && broken="$dc"
  done
  if [[ -n "$broken" ]]; then
    echo "[Compose] 警告: 发现损坏的独立 docker-compose ($broken)，请勿执行 docker-compose"
    echo "          删除: sudo rm -f /usr/local/bin/docker-compose /usr/bin/docker-compose"
    echo "          可选: sudo bash scripts/fix_compose_offline.sh"
  elif [[ -x docker-offline/cli-plugins/docker-compose ]]; then
    echo "[Compose] 插件未就绪；包内 compose 可用 → sudo bash scripts/fix_compose_offline.sh"
  else
    echo "[Compose] 未检测到（不影响部署，start-all.sh 用 docker run 启 OB）"
  fi
fi
file ism_server_user/ism_server 2>/dev/null | sed 's/^/[后端] /'
ls -lh oceanbase/*.tar 2>/dev/null | sed 's/^/[OB镜像] /' || echo "[OB镜像] 未找到"
ls -lh docker-offline/bin/dockerd 2>/dev/null | sed 's/^/[Docker离线] /' || echo "[Docker离线] 未找到 docker-offline/"
echo ""
echo "=== 服务端口（ISM 是否在跑）==="
ss -lntp 2>/dev/null | grep -E ':2881|:8091|:7090 ' || echo "  无 2881/8091/7090 监听 → 若已 start-all，请: bash scripts/diagnose_kylin.sh"
echo "=== 检查完成 ==="
echo "完全离线部署: sudo bash deploy-offline.sh"
echo "仅启动服务:   sudo bash start-all.sh"
