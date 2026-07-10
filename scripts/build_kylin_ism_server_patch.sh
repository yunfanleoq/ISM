#!/bin/bash
# 构建 ISM 麒麟 glibc 兼容后端补丁（仅替换 ism_server）
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
# shellcheck source=lib_build_id.sh
source "$ROOT/scripts/lib_build_id.sh"
BUILD_ID="$(ism_build_id)"
PKG="ism-patch-kylin-ism-server-${BUILD_ID}"
STAGING="$ROOT/releases/${PKG}"
ZIP="$ROOT/releases/${PKG}.zip"
BIN_SRC="$ROOT/patches/ism-server-kylin-glibc228/ism_server"

bash "$ROOT/scripts/build_kylin_ism_server.sh"

rm -rf "$STAGING"
mkdir -p "$STAGING/ism_server_user"

cp "$BIN_SRC" "$STAGING/ism_server_user/ism_server"
chmod 755 "$STAGING/ism_server_user/ism_server"

cat > "$STAGING/apply-patch.sh" << 'APPLY'
#!/bin/bash
# 用法: bash apply-patch.sh /opt/ISM/ism-release-oceanbase-20260707
set -euo pipefail
PATCH_ROOT="$(cd "$(dirname "$0")" && pwd)"
TARGET="${1:?用法: bash apply-patch.sh <主包目录>}"
DST="$TARGET/ism_server_user/ism_server"
[[ -d "$TARGET/ism_server_user" ]] || { echo "错误: 无效目录 $TARGET"; exit 1; }

echo "=== 应用麒麟 glibc 兼容后端补丁 ==="
if pgrep -fl ism_server >/dev/null 2>&1; then
  echo "  停止旧后端 ..."
  (cd "$TARGET" && bash stop-all.sh 2>/dev/null) || true
  sleep 2
fi
cp "$PATCH_ROOT/ism_server_user/ism_server" "$DST"
chmod 755 "$DST"
file "$DST"
if strings "$DST" | grep -qE 'GLIBC_2\.(3[2-9]|[4-9][0-9])'; then
  echo "错误: 补丁二进制仍依赖高版本 glibc"
  exit 1
fi
echo ""
echo "=== 完成 ==="
echo "  cd $TARGET"
echo "  sudo bash start-all.sh"
echo "  ss -lntp | grep -E '8091|7090'"
APPLY
chmod +x "$STAGING/apply-patch.sh"

cat > "$STAGING/README-补丁说明.md" << 'README'
# ISM 麒麟 V10 后端 glibc 补丁

## 原因

麒麟 V10 SP3 系统 **glibc 为 2.28，不能升级**。

旧包内 `ism_server` 依赖 **GLIBC_2.32~2.34**，启动即退出，导致 **8091 端口不监听**（`start-all.sh` 误报成功）。

本补丁替换为 **CGO_ENABLED=0 静态链接** 二进制，**不依赖目标机 glibc 版本**（OceanBase `dbtype=4` 部署专用）。

## 应用

```bash
bash apply-patch.sh /opt/ISM/ism-release-oceanbase-20260707
cd /opt/ISM/ism-release-oceanbase-20260707
sudo bash start-all.sh
ss -lntp | grep -E '8091|7090|2881'
curl -s http://127.0.0.1:8091/login -X POST -H 'Content-Type: application/json' \
  -d '{"Username":"admin","password":"e10adc3949ba59abbe56e057f20f883e"}'
```

## 验证

```bash
file ism_server_user/ism_server   # 应含 statically linked
ldd ism_server_user/ism_server    # 静态链接应显示 not a dynamic executable
```
README

rm -f "$ZIP"
(cd "$ROOT/releases" && COPYFILE_DISABLE=1 zip -r -q "$(basename "$ZIP")" "$(basename "$STAGING")")
echo "=== 完成 ==="
ls -lh "$ZIP" "$BIN_SRC"
