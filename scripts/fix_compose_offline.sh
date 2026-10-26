#!/bin/bash
# 修复 docker-compose 段错误：安装 Compose 插件，删除损坏的 docker-compose 命令
# 用法: sudo bash scripts/fix_compose_offline.sh
# 说明: 正确命令是「docker compose」(空格插件)，不是「docker-compose」(连字符)
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SRC="$ROOT/docker-offline/cli-plugins/docker-compose"

if [[ "$(id -u)" -ne 0 ]]; then
  echo "请使用 root: sudo bash scripts/fix_compose_offline.sh"
  exit 1
fi

if [[ ! -x "$SRC" ]]; then
  echo "错误: 缺少 $SRC"
  echo "  Compose 在 docker-offline/cli-plugins/ 下，不在 bin/ 下"
  exit 1
fi

echo "=== 修复 Docker Compose ==="

# 删除损坏的独立 docker-compose（麒麟现场常见段错误来源）
for bad in /usr/local/bin/docker-compose /usr/bin/docker-compose; do
  if [[ -e "$bad" ]]; then
    echo "  删除损坏命令: $bad"
    rm -f "$bad"
  fi
done

mkdir -p /usr/local/lib/docker/cli-plugins
install -m 755 "$SRC" /usr/local/lib/docker/cli-plugins/docker-compose

echo "=== 完成 ==="
echo "  请使用: docker compose version   （有空格，插件方式）"
echo "  不要用: docker-compose           （连字符，易段错误）"
docker compose version 2>&1 | head -1 || echo "  警告: docker compose 仍不可用，start-all.sh 会用 docker run 启动 OB"
