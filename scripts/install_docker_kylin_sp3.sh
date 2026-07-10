#!/bin/bash
# 麒麟 V10 SP3 离线安装 Docker + Docker Compose（使用部署包内 docker-offline/）
# 用法: sudo bash scripts/install_docker_kylin_sp3.sh
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
OFFLINE="$ROOT/docker-offline"
BIN_DIR="$OFFLINE/bin"
PLUGIN_DIR="$OFFLINE/cli-plugins"

if [[ "$(id -u)" -ne 0 ]]; then
  echo "请使用 root 执行: sudo bash scripts/install_docker_kylin_sp3.sh"
  exit 1
fi

if [[ ! -d "$BIN_DIR" ]] || [[ ! -x "$BIN_DIR/dockerd" ]]; then
  echo "错误: 缺少 $BIN_DIR/dockerd，请确认部署包完整（含 docker-offline/）"
  exit 1
fi

echo "=== 麒麟 V10 SP3 离线安装 Docker ==="
echo "  来源: $OFFLINE"
cat "$OFFLINE/VERSION.txt" 2>/dev/null | sed 's/^/  /' || true

echo "[1/5] 安装二进制到 /usr/local/bin ..."
install -m 755 "$BIN_DIR/docker" /usr/local/bin/docker
install -m 755 "$BIN_DIR/dockerd" /usr/local/bin/dockerd
for f in containerd containerd-shim-runc-v2 runc docker-init docker-proxy; do
  [[ -x "$BIN_DIR/$f" ]] && install -m 755 "$BIN_DIR/$f" "/usr/local/bin/$f" || true
done

echo "[2/5] 安装 Docker Compose 插件 ..."
mkdir -p /usr/local/lib/docker/cli-plugins
if [[ -x "$PLUGIN_DIR/docker-compose" ]]; then
  install -m 755 "$PLUGIN_DIR/docker-compose" /usr/local/lib/docker/cli-plugins/docker-compose
  ln -sf /usr/local/lib/docker/cli-plugins/docker-compose /usr/local/bin/docker-compose 2>/dev/null || true
else
  echo "  警告: 未找到 compose 插件，仅安装 docker 引擎"
fi

echo "[3/5] 创建 docker 用户组与数据目录 ..."
getent group docker >/dev/null || groupadd docker
mkdir -p /var/lib/docker /etc/docker
chmod 711 /var/lib/docker

echo "[4/5] 配置 systemd 服务 ..."
cat > /etc/systemd/system/docker.service << 'UNIT'
[Unit]
Description=Docker Application Container Engine (ISM offline bundle)
Documentation=https://docs.docker.com
After=network-online.target containerd.service
Wants=network-online.target

[Service]
Type=notify
ExecStart=/usr/local/bin/dockerd --host=unix:///var/run/docker.sock
ExecReload=/bin/kill -s HUP $MAINPID
TimeoutSec=0
RestartSec=2
Restart=always
LimitNOFILE=infinity
LimitNPROC=infinity
LimitCORE=infinity
TasksMax=infinity
Delegate=yes
KillMode=process

[Install]
WantedBy=multi-user.target
UNIT

cat > /etc/systemd/system/containerd.service << 'UNIT'
[Unit]
Description=containerd container runtime (ISM offline bundle)
Documentation=https://containerd.io
After=network.target

[Service]
ExecStart=/usr/local/bin/containerd
Type=notify
Delegate=yes
KillMode=process
Restart=always
RestartSec=5
LimitNPROC=infinity
LimitCORE=infinity
LimitNOFILE=1048576
TasksMax=infinity

[Install]
WantedBy=multi-user.target
UNIT

systemctl daemon-reload
systemctl enable containerd docker

echo "[5/5] 启动 Docker ..."
systemctl start containerd || true
sleep 2
systemctl start docker
sleep 2

docker --version
docker compose version 2>/dev/null || docker-compose --version 2>/dev/null || echo "Compose: 请检查 cli-plugins"

echo ""
echo "=== Docker 离线安装完成 ==="
echo "  验证: docker info"
echo "  继续部署: cd $ROOT && bash start-all.sh"
