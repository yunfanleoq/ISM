#!/bin/bash
# 构建 ISM 麒麟现场脚本补丁包 v2（修复 docker-compose 段错误）
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
# shellcheck source=lib_build_id.sh
source "$ROOT/scripts/lib_build_id.sh"
BUILD_ID="$(ism_build_id)"
PKG="ism-patch-kylin-compose-fix-${BUILD_ID}"
STAGING="$ROOT/releases/${PKG}"
ZIP="$ROOT/releases/${PKG}.zip"
REL="$ROOT/releases/ism-release-oceanbase-20260707"

rm -rf "$STAGING"
mkdir -p "$STAGING/scripts"

for f in check_env_kylin.sh fix_compose_offline.sh ensure_python.sh diagnose_kylin.sh; do
  cp "$ROOT/scripts/$f" "$STAGING/scripts/" 2>/dev/null || cp "$REL/scripts/$f" "$STAGING/scripts/"
  chmod +x "$STAGING/scripts/$f"
done

for f in start-all.sh deploy-offline.sh; do
  cp "$REL/$f" "$STAGING/"
  chmod +x "$STAGING/$f"
done

cat > "$STAGING/apply-patch.sh" << 'APPLY'
#!/bin/bash
# 用法: bash apply-patch.sh /opt/ISM/ism-release-oceanbase-20260707
set -euo pipefail
PATCH_ROOT="$(cd "$(dirname "$0")" && pwd)"
TARGET="${1:?用法: bash apply-patch.sh <主包目录>}"

[[ -d "$TARGET/scripts" ]] || { echo "错误: 无效目录 $TARGET"; exit 1; }

echo "=== 应用 Compose 修复补丁 ==="
for f in "$PATCH_ROOT/scripts/"*; do
  cp "$f" "$TARGET/scripts/$(basename "$f")"
  chmod +x "$TARGET/scripts/$(basename "$f")"
  echo "  scripts/$(basename "$f")"
done
for f in start-all.sh deploy-offline.sh; do
  cp "$PATCH_ROOT/$f" "$TARGET/$f"
  chmod +x "$TARGET/$f"
  echo "  $f"
done

echo ""
echo "=== 下一步（在目标机执行）==="
echo "  cd $TARGET"
echo "  sudo rm -f /usr/local/bin/docker-compose /usr/bin/docker-compose"
echo "  sudo bash scripts/fix_compose_offline.sh"
echo "  docker compose version          # 有空格，应正常"
echo "  sudo bash start-all.sh          # 或 deploy-offline.sh"
APPLY
chmod +x "$STAGING/apply-patch.sh"

cat > "$STAGING/README-补丁说明.md" << 'README'
# ISM 补丁：修复 docker-compose 段错误

## 原因说明

| 命令 | 状态 | 说明 |
|------|------|------|
| `docker compose version` | ✅ 正常 | **插件方式**（有空格），应使用这个 |
| `docker-compose --version` | ❌ 段错误 | **独立命令**（连字符），系统里的是损坏二进制 |

Compose **不在** `docker-offline/bin/`，而在：
```
docker-offline/cli-plugins/docker-compose
```

## 应用补丁

```bash
bash apply-patch.sh /opt/ISM/ism-release-oceanbase-20260707
cd /opt/ISM/ism-release-oceanbase-20260707

# 删除损坏的 docker-compose
sudo rm -f /usr/local/bin/docker-compose /usr/bin/docker-compose

# 安装 compose 插件
sudo bash scripts/fix_compose_offline.sh

# 验证（必须有空格）
docker compose version

# 启动
sudo bash start-all.sh
```

新版 `start-all.sh` **不再调用** `docker-compose`，直接用 `docker run` 启动 OceanBase。
README

rm -f "$ZIP"
(cd "$ROOT/releases" && COPYFILE_DISABLE=1 zip -r -q "$(basename "$ZIP")" "$(basename "$STAGING")")
echo "=== 完成 ==="
ls -lh "$ZIP"
