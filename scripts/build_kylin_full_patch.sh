#!/bin/bash
# 构建「登录 + 数据仓库/设备树/大屏懒加载」全量补丁（最新后端 + 最新前端 dist）
# 用法: bash scripts/build_kylin_full_patch.sh
# 产出: releases/ism-patch-kylin-full-YYYYMMDD-HHMM-xxxx.zip
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
# shellcheck source=lib_build_id.sh
source "$ROOT/scripts/lib_build_id.sh"
BUILD_ID="$(ism_build_id)"
DATE_TAG="$BUILD_ID"
PKG="ism-patch-kylin-full-${BUILD_ID}"
STAGING="$ROOT/releases/${PKG}"
ZIP="$ROOT/releases/${PKG}.zip"
BIN_SRC="$ROOT/patches/ism-server-kylin-glibc228/ism_server"
FE_SRC="$ROOT/ism-front-end-v2/dist"

echo "=== ISM 麒麟全量补丁构建 (${BUILD_ID}) ==="

# ── 1. 后端 ──
echo ""
if [[ "${SKIP_REBUILD:-0}" == "1" ]] && [[ -f "$BIN_SRC" ]]; then
  echo "[1/4] 跳过后端编译 (SKIP_REBUILD=1) ..."
else
  echo "[1/4] 编译 ism_server ..."
  bash "$ROOT/scripts/build_kylin_ism_server.sh"
fi

# ── 2. 前端 ──
echo ""
if [[ "${SKIP_REBUILD:-0}" == "1" ]] && [[ -f "$FE_SRC/index.html" ]]; then
  echo "[2/4] 跳过前端编译 (SKIP_REBUILD=1) ..."
else
  echo "[2/4] 编译前端 dist ..."
if [[ -x "$ROOT/scripts/check_mem_before_compile.sh" ]]; then
  if ! bash "$ROOT/scripts/check_mem_before_compile.sh"; then
    echo "错误: 内存检查 FAIL，禁止编译前端"
    exit 1
  fi
fi
rm -rf "$FE_SRC"
(
  cd "$ROOT/ism-front-end-v2"
  export NODE_OPTIONS="--max-old-space-size=20480 --openssl-legacy-provider"
  npm run build
)
[[ -d "$FE_SRC" && -f "$FE_SRC/index.html" ]] || { echo "错误: dist 未生成"; exit 1; }
echo "  dist: $(du -sh "$FE_SRC" | cut -f1)"
fi

# ── 3. 组装补丁目录 ──
echo ""
echo "[3/4] 组装补丁包 ..."
rm -rf "$STAGING"
mkdir -p "$STAGING/ism_server_user" "$STAGING/web/dist" "$STAGING/scripts"

cp "$BIN_SRC" "$STAGING/ism_server_user/ism_server"
chmod 755 "$STAGING/ism_server_user/ism_server"

# 同步 dist（排除 .map 减小体积可选，此处保留完整 dist）
rsync -a --delete "$FE_SRC/" "$STAGING/web/dist/"

for s in fix_admin_password_oceanbase.sh check_dw_device_loading.sh \
         check_login_deep.sh check_login_and_user.sh collect_diagnose_log.sh \
         run_full_field_check.sh check_env_kylin.sh \
         diagnose_getrealdata_timeout.sh fix_device_real_data_index.sh; do
  [[ -f "$ROOT/scripts/$s" ]] && cp "$ROOT/scripts/$s" "$STAGING/scripts/" && chmod +x "$STAGING/scripts/$s"
done

# ── 构建后校验（未通过则禁止出包）──
echo ""
echo "[校验] 产物自检 ..."
if rg -q 'sockjs-node|webpack-dev-server/client' "$FE_SRC/static/js/"*.js 2>/dev/null; then
  echo "错误: dist 含 dev-server 热更新代码，禁止出包"
  exit 1
fi
set +o pipefail
BIN_HAS_FIX=0
if strings "$BIN_SRC" 2>/dev/null | rg -q 'GetRealDataPaged|GetMonitorTreeLazy|userTable'; then
  BIN_HAS_FIX=1
fi
set -o pipefail
if [[ "$BIN_HAS_FIX" -ne 1 ]]; then
  echo "错误: ism_server 未包含关键修复符号"
  exit 1
fi
if ! rg -q 'pageSize|metaOnly|requestIdleCallback' "$FE_SRC/static/js/" 2>/dev/null; then
  echo "错误: dist 未包含分页/懒加载编译产物"
  exit 1
fi
# 前端 Loading 兜底：ismDebug + finishLoad / forcedTimeout 必须进包
if ! rg -q 'getRealData\.loadTimeout|forcedTimeout|__ISM_DEBUG' "$FE_SRC/static/js/" 2>/dev/null; then
  echo "错误: dist 未包含数据仓库 Loading 兜底/调试日志"
  exit 1
fi
echo "  ✓ 后端符号: userTable / GetRealDataPaged / GetMonitorTreeLazy"
echo "  ✓ 前端 dist: 无 dev-server，含分页/懒加载/Loading 兜底"

cat > "$STAGING/apply-patch.sh" << 'APPLY'
#!/bin/bash
# 用法: bash apply-patch.sh /opt/ISM/ism-release-oceanbase-20260708
set -euo pipefail
PATCH_ROOT="$(cd "$(dirname "$0")" && pwd)"
TARGET="${1:?用法: bash apply-patch.sh <主包目录>}"
[[ -d "$TARGET/ism_server_user" ]] || { echo "错误: 无效目录 $TARGET"; exit 1; }

echo "=== [1/6] 停止服务 ==="
(cd "$TARGET" && bash stop-all.sh 2>/dev/null) || true
sleep 3

echo "=== [2/6] 替换 ism_server ==="
cp "$PATCH_ROOT/ism_server_user/ism_server" "$TARGET/ism_server_user/ism_server"
chmod 755 "$TARGET/ism_server_user/ism_server"
file "$TARGET/ism_server_user/ism_server"

echo "=== [3/6] 替换前端 web/dist ==="
mkdir -p "$TARGET/web/dist"
rsync -a --delete "$PATCH_ROOT/web/dist/" "$TARGET/web/dist/"

echo "=== [4/6] 同步诊断脚本 ==="
mkdir -p "$TARGET/scripts"
cp "$PATCH_ROOT/scripts/"*.sh "$TARGET/scripts/" 2>/dev/null || true
chmod +x "$TARGET/scripts/"*.sh 2>/dev/null || true

echo "=== [5/6] 修复 device_real_data.device_uuid 索引（longtext 前缀索引）==="
if [[ -x "$TARGET/scripts/fix_device_real_data_index.sh" ]]; then
  (cd "$TARGET" && bash scripts/fix_device_real_data_index.sh) || \
    echo "  [WARN] 索引脚本失败（可稍后手工执行），继续启动"
else
  echo "  [SKIP] 无 fix_device_real_data_index.sh"
fi

echo "=== [6/6] 启动并验证 ==="
(cd "$TARGET" && bash start-all.sh)
echo "  等待 120 秒 ..."
sleep 120

BE_PORT="${ISM_BE_PORT:-8091}"
FE_PORT="${ISM_FE_PORT:-7090}"
echo "--- 登录 ---"
curl -s --compressed -m 15 -X POST "http://127.0.0.1:${BE_PORT}/login" \
  -H 'Content-Type: application/json' \
  -d '{"Username":"admin","password":"e10adc3949ba59abbe56e057f20f883e"}'
echo ""
echo "--- 端口 ---"
ss -lntp 2>/dev/null | grep -E ":${BE_PORT}|:${FE_PORT} " || true
echo ""
echo "期望: login code=1000, 8091/7090 均 LISTEN"
echo "浏览器务必 Ctrl+F5（或清缓存）后测: 数据仓库选设备 → Loading 应秒关并出表"
echo "Console 应出现: DW.getRealData.response / rendered / finally"
APPLY
chmod +x "$STAGING/apply-patch.sh"

cat > "$STAGING/README-补丁说明.md" << README
# ISM 麒麟 V10 全量补丁（最终版）

构建日期: ${DATE_TAG}
适用主包: \`ism-release-oceanbase-20260709\` / \`ism-release-oceanbase-20260708\`（同结构 OceanBase 一体包）

---

## 本包解决的问题（现场已出现过的全部）

| # | 现象 | 根因 | 本包修复 |
|---|------|------|----------|
| 1 | 登录 \`code:1003\` 用户不存在 | GORM 查 OceanBase 保留表 \`user\` 匹配失败 | 后端 \`userTable()\` / \`lookupUserByUsername\` |
| 2 | \`getRealData\` 120s 超时 | \`device_uuid\` 为 longtext 无索引，全表扫 20 万+ 行 | 分页 API + \`fix_device_real_data_index.sh\` 前缀索引 |
| 3 | 数据仓库接口已返回但页面一直 Loading | \`enablegzip\` + 手写 gzip 双重压缩，axios Promise 不 settle | 后端 \`GetRealData\` 改 \`ServeJSON\`；前端 \`finishLoad\` + 15s 强制关 Loading |
| 4 | 设备管理 / 设备树加载慢或卡死 | 整树一次性加载 | \`monitortree\` 懒加载 + 前端 \`load-data\` |
| 5 | 大屏 AppRun 长期 Loading | 全量 \`getLayerDataStruct\` 过重 | \`metaOnly\` 首屏 + 空闲 LRU 预取 |
| 6 | \`getRealData\` 参数异常崩溃 | \`IsRemoveGW\` 类型断言 panic | 安全 bool 解析 |
| 7 | 现场难排查 | 缺一键脚本 | \`diagnose_getrealdata_timeout.sh\` 等 |

**不在本包范围（环境/网络类，非代码缺陷）：**
- Modbus 设备 IP/端口不通 → 设备离线，需现场网络
- \`system_journal\` ALTER 告警 → 非致命
- 前端 7090 未监听 → 执行 \`start-all.sh\` 并等待 2 分钟

---

## 包内容

\`\`\`
ism_server_user/ism_server              # 麒麟 V10 静态链接后端
web/dist/                               # 最新生产前端（含 Loading 兜底）
scripts/fix_device_real_data_index.sh  # longtext 前缀索引
scripts/diagnose_getrealdata_timeout.sh # getRealData 超时诊断
scripts/check_*.sh / collect_*.sh       # 登录/现场检查
apply-patch.sh                          # 一键应用（含索引修复）
\`\`\`

## 应用（推荐）

\`\`\`bash
# 上传到麒麟服务器后
unzip ism-patch-kylin-full-${DATE_TAG}.zip
cd ism-patch-kylin-full-${DATE_TAG}
bash apply-patch.sh /opt/ISM/ism-release-oceanbase-20260709
\`\`\`

## 应用后验证

\`\`\`bash
cd /opt/ISM/ism-release-oceanbase-20260709

# 1. 登录（期望 code=1000）
bash scripts/check_login_deep.sh

# 2. 索引（期望有 idx_drd_device_uuid）
bash scripts/fix_device_real_data_index.sh

# 3. 数据仓库 API（期望 getRealData <2s）
DEVICE_UUID=<有测点设备> bash scripts/diagnose_getrealdata_timeout.sh

# 4. 浏览器 Ctrl+F5 → 数据仓库选设备
#    Console 期望: DW.getRealData.response / rendered / finally
#    Loading 应立刻消失并出表
\`\`\`

登录: **admin / 123456**

## 手工步骤

\`\`\`bash
cd /opt/ISM/ism-release-oceanbase-20260709
bash stop-all.sh
cp <补丁目录>/ism_server_user/ism_server ism_server_user/
rsync -a --delete <补丁目录>/web/dist/ web/dist/
cp <补丁目录>/scripts/*.sh scripts/
bash scripts/fix_device_real_data_index.sh
bash start-all.sh
# 等 2~3 分钟后浏览器 Ctrl+F5
\`\`\`
README

# ── 4. 打 zip ──
echo ""
echo "[4/4] 压缩 ..."
rm -f "$ZIP"
(
  cd "$ROOT/releases"
  COPYFILE_DISABLE=1 zip -r -q "$(basename "$ZIP")" "$(basename "$STAGING")"
)

# 同步到当前 release 目录（便于本机对照）
RELEASE_DIR="$ROOT/releases/ism-release-oceanbase-20260708"
if [[ -d "$RELEASE_DIR" ]]; then
  cp "$BIN_SRC" "$RELEASE_DIR/ism_server_user/ism_server"
  rsync -a --delete "$FE_SRC/" "$RELEASE_DIR/web/dist/"
  for s in check_dw_device_loading.sh check_login_deep.sh check_login_and_user.sh \
           collect_diagnose_log.sh run_full_field_check.sh fix_admin_password_oceanbase.sh \
           diagnose_getrealdata_timeout.sh fix_device_real_data_index.sh; do
    [[ -f "$ROOT/scripts/$s" ]] && cp "$ROOT/scripts/$s" "$RELEASE_DIR/scripts/" && chmod +x "$RELEASE_DIR/scripts/$s"
  done
  echo "  已同步到 $RELEASE_DIR"
fi

echo ""
echo "=== 全量补丁包已生成 ==="
ls -lh "$ZIP" "$BIN_SRC"
du -sh "$STAGING/web/dist"
