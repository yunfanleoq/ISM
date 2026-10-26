#!/usr/bin/env bash
# 柴发独立完整部署包：前端 + 后端 + 柴发 SQLite 库（含大屏）
# 目标：单独物理服务器部署，与中航信/循安电力包互不影响。
#
# 用法:
#   bash scripts/build_chaifa_standalone_release.sh
#   bash scripts/build_chaifa_standalone_release.sh /path/to/Sqlite3_Backup.zip
#
# 产出:
#   releases/ism-release-sqlite-chaifa-YYYYMMDD-HHMM-xxxx/
#   releases/ism-release-sqlite-chaifa-YYYYMMDD-HHMM-xxxx.zip

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
# shellcheck source=lib_build_id.sh
source "$ROOT/scripts/lib_build_id.sh"
BUILD_ID="$(ism_build_id)"
PKG_NAME="ism-release-sqlite-chaifa-${BUILD_ID}"
STAGING="$ROOT/releases/${PKG_NAME}"
ZIP_OUT="$ROOT/releases/${PKG_NAME}.zip"

BASE_REL="$ROOT/releases/ism-release-sqlite-20260714-1549-fe2e"
# 电力最新完整包：同步其前端 dist（后端电力为 CGO=0/OceanBase，柴发 SQLite 仍用下方 CGO 二进制）
POWER_REL="${ISM_POWER_REL:-$ROOT/releases/ism-release-oceanbase-20260726-2256-9a29}"
SQL_SRC="${1:-$ROOT/Sqlite3_Backup_2026-07-13_10-37-19.zip}"

if [[ ! -d "$BASE_REL" ]]; then
  echo "缺少骨架包: $BASE_REL"
  echo "请先有 ism-release-sqlite-20260714-1549-fe2e（含 deploy.sh + 目录骨架）"
  exit 1
fi
if [[ ! -d "$POWER_REL/web/dist" ]] || [[ ! -f "$POWER_REL/web/dist/index.html" ]]; then
  echo "缺少电力最新前端: $POWER_REL/web/dist"
  echo "可设置 ISM_POWER_REL=... 指向最新 ism-release-oceanbase-* 目录"
  exit 1
fi
if [[ ! -f "$SQL_SRC" ]]; then
  echo "缺少柴发备份: $SQL_SRC"
  exit 1
fi
if [[ ! -x "$(command -v python3)" ]]; then
  echo "需要 python3"
  exit 1
fi

echo "=== 柴发独立完整包构建 ==="
echo "包名: $PKG_NAME"
echo "骨架: $BASE_REL"
echo "电力前端同步自: $POWER_REL"
echo "数据: $SQL_SRC"
echo ""

rm -rf "$STAGING"
mkdir -p "$STAGING"

echo "[1/5] 复制骨架（后端二进制 + 前端 dist + 脚本）..."
# APFS clonefile 优先，失败则普通复制
if cp -c -R "$BASE_REL/." "$STAGING/" 2>/dev/null; then
  echo "  使用 APFS clone"
else
  rsync -a --exclude 'logs/*' "$BASE_REL/" "$STAGING/"
  echo "  使用 rsync"
fi
rm -rf "$STAGING/logs"
mkdir -p "$STAGING/logs"
# 去掉骨架里的旧库
rm -f "$STAGING/ism_server_user/data/db/ism.db" \
      "$STAGING/ism_server_user/data/db/ism.db-shm" \
      "$STAGING/ism_server_user/data/db/ism.db-wal"

# ★ SQLite + 麒麟铁律：
# 1) CGO_ENABLED=1（否则 go-sqlite3 stub panic）
# 2) 动态链接 GLIBC ≤ 2.28（否则麒麟报 GLIBC_2.32+ not found）
verify_kylin_cgo_bin() {
  local f="$1"
  [[ -f "$f" ]] || return 1
  file "$f" | grep -q "ELF 64-bit.*x86-64" || return 1
  if strings "$f" 2>/dev/null | grep -qF "go-sqlite3 requires cgo to work. This is a stub"; then
    return 1
  fi
  if strings "$f" 2>/dev/null | grep -qE 'GLIBC_2\.(3[2-9]|[4-9][0-9])'; then
    return 1
  fi
  return 0
}

pick_cgo_linux_bin() {
  local candidates=(
    "$ROOT/patches/ism-server-kylin-glibc228-cgo/ism_server"
    "$ROOT/patches/ism-server-cgo-linux-amd64/ism_server"
    "$ROOT/releases/ism-release-sqlite-20260706/ism_server_user/ism_server"
    "$ROOT/releases/ism-release-sqlite-20260706-offline/ism_server_user/ism_server"
    "$ROOT/releases/ism-release-20260703/ism_server_user/ism_server"
  )
  local f
  for f in "${candidates[@]}"; do
    verify_kylin_cgo_bin "$f" || continue
    echo "$f"
    return 0
  done
  return 1
}

echo "[1b/5] 校验/替换 Linux amd64 CGO 后端（SQLite + 麒麟 glibc2.28）..."
CGO_BIN="$(pick_cgo_linux_bin || true)"
if [[ -z "${CGO_BIN}" ]]; then
  echo "  未找到兼容二进制，本机交叉编译（messense sysroot glibc 2.28）..."
  bash "$ROOT/scripts/build_cgo_linux_amd64.sh" "$ROOT/patches/ism-server-kylin-glibc228-cgo/ism_server"
  CGO_BIN="$ROOT/patches/ism-server-kylin-glibc228-cgo/ism_server"
fi
if ! verify_kylin_cgo_bin "$CGO_BIN"; then
  echo "错误: $CGO_BIN 不满足 CGO + GLIBC≤2.28，禁止打包"
  exit 1
fi
cp "$CGO_BIN" "$STAGING/ism_server_user/ism_server"
chmod 755 "$STAGING/ism_server_user/ism_server"
echo "  使用 CGO/glibc228 二进制: $CGO_BIN ($(du -h "$STAGING/ism_server_user/ism_server" | awk '{print $1}'))"

# 前端静态代理：覆盖为 Python 3.7 兼容版（麒麟常见 3.7.9）
mkdir -p "$STAGING/scripts"
cp "$ROOT/scripts/serve_test_frontend.py" "$STAGING/scripts/serve_test_frontend.py"
if grep -qE 'if [A-Za-z_][A-Za-z0-9_]* :=' "$STAGING/scripts/serve_test_frontend.py"; then
  echo "错误: serve_test_frontend.py 仍含 walrus 赋值（Python 3.7 不支持）"
  exit 1
fi
echo "  已写入 Python3.7 兼容 serve_test_frontend.py"

echo "[1c/5] 用电力最新前端覆盖 web/dist ..."
rm -rf "$STAGING/web/dist"
mkdir -p "$STAGING/web"
if cp -c -R "$POWER_REL/web/dist" "$STAGING/web/dist" 2>/dev/null; then
  echo "  使用 APFS clone 同步前端"
else
  rsync -a "$POWER_REL/web/dist/" "$STAGING/web/dist/"
  echo "  使用 rsync 同步前端"
fi
[[ -f "$STAGING/web/dist/index.html" ]] || { echo "错误: 同步后无 index.html"; exit 1; }
if rg -q 'sockjs-node|webpack-dev-server/client' "$STAGING/web/dist/static/js/"*.js 2>/dev/null; then
  echo "错误: 电力前端 dist 含 dev-server 热更新代码，禁止打包"
  exit 1
fi
echo "  前端已同步: $(du -sh "$STAGING/web/dist" | awk '{print $1}') ← $(basename "$POWER_REL")"

echo "[2/5] 制备柴发 ism.db（还原备份 + 楼层树 + 大屏 + admin 密码）..."
python3 "$ROOT/scripts/prepare_chaifa_release_db.py" \
  --sql "$SQL_SRC" \
  --out "$STAGING/ism_server_user/data/db/ism.db" \
  --project-name "后沙峪改造-柴发部分"

echo "[3/5] 调整端口与配置（柴发独立机默认 7080/8081）..."
# 独立物理机可用标准端口；与电力机房另一台服务器无关
cat > "$STAGING/ports.env" <<'EOF'
# 柴发独立服务器默认端口（可按现场修改）
ISM_FE_PORT=7080
ISM_BE_PORT=8081
EOF

CONF="$STAGING/ism_server_user/conf/app.conf"
if [[ -f "$CONF" ]]; then
  if [[ "$(uname -s)" == "Darwin" ]]; then
    sed -i '' 's/^httpport=.*/httpport=8081/' "$CONF"
    sed -i '' 's/^dbtype=.*/dbtype=1/' "$CONF"
  else
    sed -i 's/^httpport=.*/httpport=8081/' "$CONF"
    sed -i 's/^dbtype=.*/dbtype=1/' "$CONF"
  fi
fi

# deploy.sh 里默认 BE 改为 8081（若脚本硬编码 8091）
if [[ -f "$STAGING/deploy.sh" ]]; then
  if [[ "$(uname -s)" == "Darwin" ]]; then
    sed -i '' 's/BE_PORT="${ISM_BE_PORT:-8091}"/BE_PORT="${ISM_BE_PORT:-8081}"/' "$STAGING/deploy.sh" || true
  else
    sed -i 's/BE_PORT="${ISM_BE_PORT:-8091}"/BE_PORT="${ISM_BE_PORT:-8081}"/' "$STAGING/deploy.sh" || true
  fi
fi

echo "[4/5] 写入部署说明与 BUILD_INFO..."
DB_SIZE="$(du -h "$STAGING/ism_server_user/data/db/ism.db" | awk '{print $1}')"
FE_SIZE="$(du -sh "$STAGING/web/dist" 2>/dev/null | awk '{print $1}')"
BE_INFO="$(file -b "$STAGING/ism_server_user/ism_server" 2>/dev/null || echo missing)"

cat > "$STAGING/BUILD_INFO.txt" <<EOF
包名: ${PKG_NAME}
构建时间: $(date '+%Y-%m-%d %H:%M:%S')
构建主机: $(uname -s)/$(uname -m)
类型: 柴发独立完整包（前后端 + SQLite 业务库 + 监控大屏）
数据来源: $(basename "$SQL_SRC")
后端: ${BE_INFO}
后端来源(CGO): ${CGO_BIN}
后端校验: 无 go-sqlite3 CGO stub
前端 dist: ${FE_SIZE}
前端同步自电力包: $(basename "$POWER_REL")
前端代理: Python 3.7 兼容（无 :=）
数据库: ${DB_SIZE}
dbtype: 1 (SQLite)
httpport: 8081
项目: 后沙峪改造-柴发部分
登录: admin / 123456
说明: 与中航信/循安电力包分属不同物理服务器，本包自包含，无需导入到另一套系统
修复: 同步电力最新前端 + CGO=1 SQLite + serve_test_frontend.py 兼容麒麟 Python 3.7.9
EOF

cat > "$STAGING/README-部署说明.md" <<EOF
# ISM 柴发监控 — 独立完整部署包

> 本包用于 **柴发专用物理服务器** 独立部署，自包含前端、后端、SQLite 数据库与监控大屏。
> **不要** 往中航信/循安那套系统里「再导一个项目」——两套系统分别装在两台机器上。

- 包名: \`${PKG_NAME}\`
- 平台: **Linux x86_64**
- 数据库: SQLite（\`dbtype=1\`），已含「后沙峪改造-柴发部分」全部设备点位 + 楼层树 + 大屏
- 默认端口: 前端 **7080** / 后端 **8081**
- 登录: **admin** / **123456**

## 快速部署

\`\`\`bash
# 上传 zip 后
mkdir -p /opt/ism && cd /opt/ism
unzip -o ${PKG_NAME}.zip
cd ${PKG_NAME}
bash deploy.sh
\`\`\`

浏览器访问: \`http://<柴发服务器IP>:7080/#/login\`

登录后进入项目「后沙峪改造-柴发部分」，默认首页为「柴发楼监控大屏」
（与中航信最新一致的三页运行模板）：

- 首页模板（home）
- 设备列表模板（deviceList）— 树点选楼层/区域
- 点位列表模板（datapointList）— 树点选设备

## 目录结构

\`\`\`
${PKG_NAME}/
├── deploy.sh
├── start-test.sh / stop-test.sh
├── ports.env                 # FE=7080 BE=8081
├── ism_server_user/
│   ├── ism_server            # Linux amd64
│   ├── conf/app.conf         # dbtype=1, httpport=8081
│   └── data/db/ism.db        # 柴发业务库
├── web/dist/                 # 前端静态资源
└── scripts/
\`\`\`

## 与另一套系统的关系

| 系统 | 服务器 | 典型目录 | 本包 |
|------|--------|----------|------|
| 电力 / 中航信 / 循安 | 物理机 A | \`/opt/ISMCode/ism_web\` 等 | **无关** |
| **柴发** | 物理机 B | 建议 \`/opt/ism/${PKG_NAME}\` | **本包** |

两套互不共享数据库，也无需互相导入项目。

## 注意

1. 包内 \`ism_server\` 为 **Linux amd64 + CGO_ENABLED=1 + GLIBC≤2.28**（SQLite + 麒麟 V10）；CGO=0 会 panic，GLIBC>2.28 无法加载。
2. 前端静态服务兼容 **Python 3.7+**（麒麟自带 3.7.9 可用）。
3. 设备 IP/端口来自原现场备份（如 \`172.31.97.x\`）；现场网段变化时在设备扩展参数中修改。
4. 勿覆盖客户已有生产目录；解压到独立路径后执行 \`deploy.sh\`。
5. 改端口: \`ISM_FE_PORT=7090 ISM_BE_PORT=8082 bash deploy.sh\`
EOF

# 附带柴发项目包 JSON（可选二次导入/备份）
mkdir -p "$STAGING/data-source"
if [[ -f "$ROOT/projects-import/柴发监控/后沙峪改造-柴发部分_ISM项目包.json" ]]; then
  cp "$ROOT/projects-import/柴发监控/后沙峪改造-柴发部分_ISM项目包.json" \
     "$STAGING/data-source/houshayu-chaifa-ISM-project-package.json" || true
fi
cp "$SQL_SRC" "$STAGING/data-source/$(basename "$SQL_SRC")" 2>/dev/null || true

echo "[5/5] 打包 zip（较大，约数分钟）..."
rm -f "$ZIP_OUT"
(
  cd "$ROOT/releases"
  zip -r -q "$(basename "$ZIP_OUT")" "$(basename "$STAGING")" \
    -x "*/logs/*" -x "*.DS_Store"
)

echo ""
echo "完成:"
echo "  目录: $STAGING"
echo "  ZIP:  $ZIP_OUT ($(du -h "$ZIP_OUT" | awk '{print $1}'))"
cat "$STAGING/BUILD_INFO.txt"
