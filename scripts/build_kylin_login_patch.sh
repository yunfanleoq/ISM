#!/bin/bash
# 构建「登录 1003 修复」麒麟后端补丁包
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
# shellcheck source=lib_build_id.sh
source "$ROOT/scripts/lib_build_id.sh"
BUILD_ID="$(ism_build_id)"
PKG="ism-patch-kylin-login-${BUILD_ID}"
STAGING="$ROOT/releases/${PKG}"
ZIP="$ROOT/releases/${PKG}.zip"
BIN_SRC="$ROOT/patches/ism-server-kylin-glibc228/ism_server"

bash "$ROOT/scripts/build_kylin_ism_server.sh"

rm -rf "$STAGING"
mkdir -p "$STAGING/ism_server_user" "$STAGING/scripts"

cp "$BIN_SRC" "$STAGING/ism_server_user/ism_server"
chmod 755 "$STAGING/ism_server_user/ism_server"
cp "$ROOT/scripts/fix_admin_password_oceanbase.sh" "$STAGING/scripts/"
chmod +x "$STAGING/scripts/fix_admin_password_oceanbase.sh"

cat > "$STAGING/apply-patch.sh" << 'APPLY'
#!/bin/bash
# 用法: bash apply-patch.sh /opt/ISM/ism-release-oceanbase-20260708
set -euo pipefail
PATCH_ROOT="$(cd "$(dirname "$0")" && pwd)"
TARGET="${1:?用法: bash apply-patch.sh <主包目录>}"
DST="$TARGET/ism_server_user/ism_server"
[[ -d "$TARGET/ism_server_user" ]] || { echo "错误: 无效目录 $TARGET"; exit 1; }

echo "=== [1/4] 停止服务 ==="
(cd "$TARGET" && bash stop-all.sh 2>/dev/null) || true
sleep 2

echo "=== [2/4] 替换 ism_server ==="
cp "$PATCH_ROOT/ism_server_user/ism_server" "$DST"
chmod 755 "$DST"
file "$DST"

echo "=== [3/4] 重置 admin 密码（bcrypt(MD5(123456))）==="
cp "$PATCH_ROOT/scripts/fix_admin_password_oceanbase.sh" "$TARGET/scripts/" 2>/dev/null || mkdir -p "$TARGET/scripts"
cp "$PATCH_ROOT/scripts/fix_admin_password_oceanbase.sh" "$TARGET/scripts/"
bash "$TARGET/scripts/fix_admin_password_oceanbase.sh" || echo "  警告: 密码脚本需在 OceanBase 运行后执行"

echo "=== [4/4] 启动并验证 ==="
(cd "$TARGET" && bash start-all.sh)
echo "  等待 120 秒 ..."
sleep 120
BE_PORT="${ISM_BE_PORT:-8091}"
curl -s -m 15 -X POST "http://127.0.0.1:${BE_PORT}/login" \
  -H 'Content-Type: application/json' \
  -d '{"Username":"admin","password":"e10adc3949ba59abbe56e057f20f883e"}'
echo ""
echo "期望 code: 1000"
APPLY
chmod +x "$STAGING/apply-patch.sh"

cat > "$STAGING/README-补丁说明.md" << 'README'
# ISM 麒麟 V10 登录 1003 修复补丁

## 修复内容

1. **GORM 查 `user` 表**：OceanBase 下保留表名 `user` 导致 `CheckLogin` 返回 1003；改为显式 `` `user` `` 表查询。
2. **admin 密码**：SQL 导入的 hash 与前端 MD5 链路不一致；补丁附带 `fix_admin_password_oceanbase.sh` 重置为 `bcrypt(MD5("123456"))`。

## 应用（推荐一键）

```bash
unzip ism-patch-kylin-login-*.zip
cd ism-patch-kylin-login-*
bash apply-patch.sh /opt/ISM/ism-release-oceanbase-20260708
```

## 手工步骤

```bash
bash stop-all.sh
cp ism_server_user/ism_server /opt/ISM/ism-release-oceanbase-20260708/ism_server_user/
bash scripts/fix_admin_password_oceanbase.sh
bash start-all.sh
# 等 2 分钟后
curl -s -X POST http://127.0.0.1:8091/login \
  -H 'Content-Type: application/json' \
  -d '{"Username":"admin","password":"e10adc3949ba59abbe56e057f20f883e"}'
```

期望：`"code":1000`
README

rm -f "$ZIP"
(cd "$ROOT/releases" && COPYFILE_DISABLE=1 zip -r -q "$(basename "$ZIP")" "$(basename "$STAGING")")
echo "=== 补丁包已生成 ==="
ls -lh "$ZIP" "$BIN_SRC"
