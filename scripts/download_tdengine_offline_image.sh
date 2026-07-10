#!/bin/bash
# 下载 TDengine Docker 镜像并导出为离线 tar（linux/amd64）
# 用法: bash scripts/download_tdengine_offline_image.sh [输出目录]
#
# 说明: 部分 Docker/Colima 环境 docker save 会因 digest 缺失失败，
#       本脚本优先 docker pull+save，失败则用 crane pull。
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
OUT_DIR="${1:-$ROOT/releases/tdengine-offline}"
TD_IMAGE="${TD_IMAGE:-tdengine/tdengine:3.3.6.13}"
TD_TAR="${OUT_DIR}/tdengine.tar"

mkdir -p "$OUT_DIR"

PULL_SOURCES=(
  "docker.1ms.run/${TD_IMAGE}"
  "docker.m.daocloud.io/${TD_IMAGE}"
  "${TD_IMAGE}"
)

ensure_crane() {
  if command -v crane >/dev/null 2>&1; then
    command -v crane
    return 0
  fi
  if [[ -x /tmp/crane ]]; then
    echo /tmp/crane
    return 0
  fi
  local arch carch
  arch="$(uname -m)"
  case "$arch" in
    arm64|aarch64) carch=arm64 ;;
    *) carch=x86_64 ;;
  esac
  local os
  os="$(uname -s)"
  local url
  if [[ "$os" == "Darwin" ]]; then
    url="https://github.com/google/go-containerregistry/releases/download/v0.20.2/go-containerregistry_Darwin_${carch}.tar.gz"
  else
    url="https://github.com/google/go-containerregistry/releases/download/v0.20.2/go-containerregistry_Linux_${carch}.tar.gz"
  fi
  curl -fsSL "$url" -o /tmp/crane.tgz
  tar -xzf /tmp/crane.tgz -C /tmp crane
  echo /tmp/crane
}

export_with_crane() {
  local src="$1"
  local crane_bin
  crane_bin="$(ensure_crane)"
  echo "  使用 crane pull --platform linux/amd64 $src ..."
  "$crane_bin" pull --platform linux/amd64 "$src" "$TD_TAR"
}

if ! command -v docker >/dev/null 2>&1; then
  echo "警告: 无 Docker，仅尝试 crane ..."
  for src in "${PULL_SOURCES[@]}"; do
    if export_with_crane "$src"; then
      printf '%s\n' "$TD_IMAGE" "$src" > "$OUT_DIR/IMAGE_TAG"
      echo "已导出: $(du -sh "$TD_TAR" | cut -f1) -> $TD_TAR"
      exit 0
    fi
  done
  echo "错误: 无法拉取 TDengine 镜像"
  exit 1
fi

pulled=""
for src in "${PULL_SOURCES[@]}"; do
  echo "拉取 --platform linux/amd64 $src ..."
  if docker pull --platform linux/amd64 "$src" >/dev/null 2>&1; then
    docker tag "$src" "$TD_IMAGE" 2>/dev/null || true
    pulled="$src"
    echo "  已拉取: $src"
    break
  fi
done

if [[ -z "$pulled" ]]; then
  echo "docker pull 失败，改用 crane ..."
  for src in "${PULL_SOURCES[@]}"; do
    if export_with_crane "$src"; then
      printf '%s\n' "$TD_IMAGE" "$src" > "$OUT_DIR/IMAGE_TAG"
      echo "已导出: $(du -sh "$TD_TAR" | cut -f1) -> $TD_TAR"
      exit 0
    fi
  done
  echo "错误: 无法拉取 TDengine 镜像"
  exit 1
fi

rm -f "$TD_TAR"
if docker save "$TD_IMAGE" -o "$TD_TAR" 2>/dev/null || docker save "$pulled" -o "$TD_TAR" 2>/dev/null; then
  :
else
  echo "  docker save 失败，改用 crane ..."
  export_with_crane "$pulled"
fi

printf '%s\n' "$TD_IMAGE" "$pulled" > "$OUT_DIR/IMAGE_TAG"
echo "已导出: $(du -sh "$TD_TAR" | cut -f1) -> $TD_TAR"
