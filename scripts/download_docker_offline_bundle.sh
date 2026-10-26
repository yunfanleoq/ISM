#!/bin/bash
# 下载 Docker 静态二进制 + Docker Compose（linux/amd64），供麒麟 V10 SP3 离线安装
# 用法: bash scripts/download_docker_offline_bundle.sh
# 产出: docker-offline/ 目录（打入 ISM 部署包）

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
OUT="$ROOT/docker-offline"
DOCKER_VER="${DOCKER_VER:-24.0.9}"
COMPOSE_VER="${COMPOSE_VER:-2.29.7}"

mkdir -p "$OUT/bin" "$OUT/cli-plugins"

echo "=== 下载 Docker 离线组件 (linux/amd64) ==="
echo "  Docker:  ${DOCKER_VER}"
echo "  Compose: ${COMPOSE_VER}"
echo "  输出:    $OUT"
echo ""

DOCKER_TGZ="$OUT/docker-${DOCKER_VER}.tgz"
COMPOSE_BIN="$OUT/cli-plugins/docker-compose"

if [[ ! -f "$DOCKER_TGZ" ]]; then
  url="https://download.docker.com/linux/static/stable/x86_64/docker-${DOCKER_VER}.tgz"
  echo "[1/2] 下载 Docker 静态包 ..."
  for mirror in "$url" "https://mirrors.aliyun.com/docker-ce/linux/static/stable/x86_64/docker-${DOCKER_VER}.tgz"; do
    if curl -fsSL --connect-timeout 30 -o "$DOCKER_TGZ" "$mirror"; then
      echo "  来源: $mirror"
      break
    fi
  done
  [[ -f "$DOCKER_TGZ" ]] || { echo "错误: Docker 静态包下载失败"; exit 1; }
else
  echo "[1/2] 已存在 Docker 静态包，跳过"
fi

if [[ ! -f "$COMPOSE_BIN" ]]; then
  echo "[2/2] 下载 Docker Compose 二进制 ..."
  compose_url="https://github.com/docker/compose/releases/download/v${COMPOSE_VER}/docker-compose-linux-x86_64"
  for mirror in "$compose_url" "https://ghproxy.net/${compose_url}"; do
    if curl -fsSL --connect-timeout 60 -o "$COMPOSE_BIN" "$mirror"; then
      echo "  来源: $mirror"
      break
    fi
  done
  [[ -f "$COMPOSE_BIN" ]] || { echo "错误: Docker Compose 下载失败"; exit 1; }
  chmod +x "$COMPOSE_BIN"
else
  echo "[2/2] 已存在 Docker Compose，跳过"
fi

# 解压 docker 二进制到 bin/ 便于安装脚本直接使用
rm -rf "$OUT/bin"/*
tar -xzf "$DOCKER_TGZ" -C "$OUT/bin" --strip-components=1 docker/docker dockerd containerd containerd-shim-runc-v2 runc docker-init docker-proxy 2>/dev/null || \
tar -xzf "$DOCKER_TGZ" -C "$OUT" && mv "$OUT/docker/"* "$OUT/bin/" 2>/dev/null || true

chmod +x "$OUT/bin/"* "$OUT/cli-plugins/"* 2>/dev/null || true

cat > "$OUT/VERSION.txt" << EOF
Docker static: ${DOCKER_VER}
Docker Compose: v${COMPOSE_VER}
Platform: linux/amd64
Target OS: 银河麒麟 V10 SP3 x86_64
Downloaded: $(date '+%Y-%m-%d %H:%M:%S')
EOF

echo ""
echo "=== 完成 ==="
du -sh "$OUT" "$DOCKER_TGZ" "$COMPOSE_BIN" 2>/dev/null
ls -lh "$OUT/bin/" | head -10
