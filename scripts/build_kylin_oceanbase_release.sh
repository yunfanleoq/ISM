#!/bin/bash
# 构建麒麟 V10 + OceanBase 完整发布包（新装 + 覆盖升级同一份）
# 用法: bash scripts/build_kylin_oceanbase_release.sh
# 产出: releases/ism-release-oceanbase-YYYYMMDD-HHMM-<shortsha>/
#       releases/ism-release-oceanbase-YYYYMMDD-HHMM-<shortsha>.zip
#
# 与 0817 主包同族：复用其离线组件/启动脚本/conf 模板，叠入今日 ism_server + web/dist。
# 禁止 go mod tidy / go mod vendor；禁止把本地 sqlite dbtype=1 打进包。
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

BASE="$ROOT/releases/ism-release-oceanbase-20260817-0001-c851"
BIN_SRC="$ROOT/patches/ism-server-kylin-glibc228/ism_server"
FE_SRC="$ROOT/ism-front-end-v2/dist"
PATCH_BIN="$ROOT/releases/ism-patch-kylin-20260820-20260821-2319-2056/ism_server_user/ism_server"
PATCH_FE="$ROOT/releases/ism-patch-kylin-20260820-20260821-2319-2056/web/dist"

DATE_TAG="$(date +%Y%m%d)"
TIME_TAG="$(date +%H%M)"
SHORT_SHA="$(git rev-parse --short HEAD 2>/dev/null || echo unknown)"
PKG_NAME="ism-release-oceanbase-${DATE_TAG}-${TIME_TAG}-${SHORT_SHA}"
STAGING="$ROOT/releases/${PKG_NAME}"
ZIP_OUT="$ROOT/releases/${PKG_NAME}.zip"

echo "=== ISM OceanBase 完整发布包 ==="
echo "包名: ${PKG_NAME}"
echo "基线: $(basename "$BASE")"
echo ""

[[ -d "$BASE" ]] || { echo "错误: 找不到 0817 主包目录 $BASE"; exit 1; }
[[ -f "$BASE/start-all.sh" ]] || { echo "错误: 0817 主包缺少 start-all.sh"; exit 1; }
[[ -f "$BASE/ism_server_user/conf/app.conf" ]] || { echo "错误: 0817 缺少 conf"; exit 1; }
grep -q '^dbtype=4' "$BASE/ism_server_user/conf/app.conf" || {
  echo "错误: 0817 app.conf 不是 dbtype=4"
  exit 1
}

verify_kylin_binary() {
  local bin="$1"
  [[ -f "$bin" ]] || return 1
  file "$bin" | grep -q "ELF 64-bit.*x86-64" || return 1
  if strings "$bin" 2>/dev/null | grep -qE 'GLIBC_2\.(3[2-9]|[4-9][0-9])'; then
    return 1
  fi
  return 0
}

build_backend_if_needed() {
  echo "[backend] 今日二进制不可用，按麒麟静态链接重编 (-mod=vendor) ..."
  mkdir -p "$(dirname "$BIN_SRC")"
  (cd "$ROOT/ism_server_user" && \
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -mod=vendor -ldflags "-w -s" -o "$BIN_SRC" .)
  chmod 755 "$BIN_SRC"
  verify_kylin_binary "$BIN_SRC" || { echo "错误: 新编二进制不兼容麒麟 V10"; exit 1; }
}

echo "[1/6] 核对后端二进制 ..."
if [[ -f "$BIN_SRC" ]] && verify_kylin_binary "$BIN_SRC"; then
  echo "  复用: $BIN_SRC ($(du -sh "$BIN_SRC" | cut -f1), $(stat -f '%Sm' -t '%Y-%m-%d %H:%M' "$BIN_SRC" 2>/dev/null || stat -c '%y' "$BIN_SRC"))"
elif [[ -f "$PATCH_BIN" ]] && verify_kylin_binary "$PATCH_BIN"; then
  echo "  复用 0820 补丁二进制"
  mkdir -p "$(dirname "$BIN_SRC")"
  cp "$PATCH_BIN" "$BIN_SRC"
  chmod 755 "$BIN_SRC"
else
  build_backend_if_needed
fi
file "$BIN_SRC" | grep -q "statically linked" || echo "  警告: 非静态链接，请现场确认 glibc"

echo "[2/6] 核对前端 dist ..."
if [[ ! -f "$FE_SRC/index.html" ]] && [[ -f "$PATCH_FE/index.html" ]]; then
  echo "  使用 0820 补丁 dist"
  FE_SRC="$PATCH_FE"
fi
[[ -f "$FE_SRC/index.html" ]] || { echo "错误: 无 web/dist（禁止用旧 dist 冒充）"; exit 1; }
if command -v rg >/dev/null 2>&1; then
  if rg -q 'sockjs-node|webpack-dev-server/client' "$FE_SRC/static/js/"*.js 2>/dev/null; then
    echo "错误: dist 含 webpack-dev-server，禁止打包"
    exit 1
  fi
  rg -q 'resolveSavePageId' "$FE_SRC/static/js/"*.js || { echo "错误: dist 缺少 resolveSavePageId"; exit 1; }
  rg -q 'includeEmpty: true' "$FE_SRC/static/js/"*.js || { echo "错误: dist 缺少 includeEmpty: true"; exit 1; }
  rg -q 'dw-category-tags' "$FE_SRC/static/js/"*.js "$FE_SRC/static/css/"*.css 2>/dev/null \
    || rg -q 'dw-category-tags' "$FE_SRC/static/js/"*.js \
    || { echo "错误: dist 缺少 dw-category-tags（0825 设备分类 Tag）"; exit 1; }
  rg -q 'UnboundPage' "$FE_SRC/static/js/"*.js || { echo "错误: dist 缺少 UnboundPage"; exit 1; }
  rg -q 'homeFromModelUuid' "$FE_SRC/static/js/"*.js || { echo "错误: dist 缺少 homeFromModelUuid（0826 模型UUID绑首页）"; exit 1; }
  rg -q '变化百分比' "$FE_SRC/static/js/"*.js || { echo "错误: dist 缺少五种存储类型"; exit 1; }
fi
echo "  dist: $(du -sh "$FE_SRC" | cut -f1), index.html $(stat -f '%Sm' -t '%Y-%m-%d %H:%M' "$FE_SRC/index.html" 2>/dev/null || stat -c '%y' "$FE_SRC/index.html")"

echo "[3/6] 以 0817 主包为骨架组装（硬链，不复制离线镜像）..."
rm -rf "$STAGING"
mkdir -p "$ROOT/releases"
cp -al "$BASE" "$STAGING"

# 替换程序与前端（先断开硬链，避免改到 0817）
rm -f "$STAGING/ism_server_user/ism_server"
cp "$BIN_SRC" "$STAGING/ism_server_user/ism_server"
chmod 755 "$STAGING/ism_server_user/ism_server"

rm -rf "$STAGING/web/dist"
mkdir -p "$STAGING/web"
cp -al "$FE_SRC" "$STAGING/web/dist"

# 运行时目录保持空，不打现场数据
mkdir -p "$STAGING/logs" "$STAGING/ism_server_user/data/db"
find "$STAGING/logs" -mindepth 1 -delete 2>/dev/null || true
rm -f "$STAGING/ism_server_user/data/db/"*.db "$STAGING/ism_server_user/data/db/"*.db-* 2>/dev/null || true
rm -rf "$STAGING/tdengine/data" "$STAGING/tdengine/log" 2>/dev/null || true

# 断开将改写的文件硬链
break_link() {
  local f="$1"
  [[ -f "$f" ]] || return 0
  local tmp="${f}.new.$$"
  cp "$f" "$tmp"
  rm -f "$f"
  mv "$tmp" "$f"
}

break_link "$STAGING/deploy-offline.sh"
# 0817 deploy-offline.sh 的 ROOT=dirname/.. 会指到包外；新包装正为包根
if grep -q 'dirname "$0")/..' "$STAGING/deploy-offline.sh" 2>/dev/null; then
  if [[ "$(uname -s)" == "Darwin" ]]; then
    sed -i '' 's|dirname "$0")/..|dirname "$0")|' "$STAGING/deploy-offline.sh"
  else
    sed -i 's|dirname "$0")/..|dirname "$0")|' "$STAGING/deploy-offline.sh"
  fi
fi
chmod +x "$STAGING/deploy-offline.sh" "$STAGING/start-all.sh" "$STAGING/stop-all.sh"

# 确认 conf 仍是 0817 OceanBase 模板
grep -q '^dbtype=4' "$STAGING/ism_server_user/conf/app.conf" || {
  echo "错误: 组装后 app.conf 不是 dbtype=4（禁止打入本地 sqlite）"
  exit 1
}

echo "[4/6] 写入新装 / 覆盖升级脚本与 README ..."

cat > "$STAGING/upgrade-existing.sh" << 'UPGRADE'
#!/bin/bash
# 把本完整包中的程序/前端覆盖到已有 ISM 安装。
# 用法: bash upgrade-existing.sh /opt/ISM/ism-release-oceanbase-20260817-0001-c851
# 覆盖: ism_server、web/dist
# 不覆盖: data/、logs/、conf/、证书、许可证、OceanBase/TDengine 数据
set -euo pipefail

PKG="$(cd "$(dirname "$0")" && pwd)"
TARGET="${1:-}"
if [[ -z "$TARGET" ]]; then
  echo "用法: bash upgrade-existing.sh <已有ISM安装目录>"
  echo "示例: bash upgrade-existing.sh /opt/ISM/ism-release-oceanbase-20260817-0001-c851"
  exit 1
fi
[[ -d "$TARGET" ]] || { echo "错误: 目录不存在 $TARGET"; exit 1; }
TARGET="$(cd "$TARGET" && pwd)"
[[ -d "$TARGET/ism_server_user" ]] || { echo "错误: 不是 ISM 安装目录: $TARGET"; exit 1; }
[[ -x "$PKG/ism_server_user/ism_server" ]] || { echo "错误: 本包缺少 ism_server"; exit 1; }
[[ -f "$PKG/web/dist/index.html" ]] || { echo "错误: 本包缺少 web/dist"; exit 1; }

TS="$(date +%Y%m%d-%H%M)"
BIN="$TARGET/ism_server_user/ism_server"
DIST="$TARGET/web/dist"
BIN_BAK="${BIN}.bak-${TS}"
DIST_BAK="${TARGET}/web/dist.bak-${TS}"
ROLLED=0

rollback() {
  [[ "$ROLLED" -eq 1 ]] && return 0
  ROLLED=1
  echo "!!! 覆盖失败，尝试回滚 ..."
  if [[ -f "$BIN_BAK" ]]; then
    cp -a "$BIN_BAK" "$BIN" 2>/dev/null || true
    chmod 755 "$BIN" 2>/dev/null || true
    echo "  已回滚 ism_server"
  fi
  if [[ -d "$DIST_BAK" ]] && [[ ! -f "$DIST/index.html" ]]; then
    mv "$DIST_BAK" "$DIST" 2>/dev/null || true
    echo "  已回滚 web/dist"
  fi
}

trap rollback ERR

echo "=== ISM 覆盖升级（不触碰业务数据）==="
echo "  来源: $PKG"
echo "  目标: $TARGET"
echo "  覆盖: ism_server_user/ism_server 、 web/dist"
echo "  保留: data/ logs/ conf/ 证书 许可证 OceanBase TDengine"
echo ""

echo "[1/4] 停止应用进程（不停止 OceanBase / TDengine 容器）"
if [[ -x "$TARGET/stop-all.sh" ]]; then
  (cd "$TARGET" && bash stop-all.sh) || true
else
  echo "  无 stop-all.sh，跳过"
fi
sleep 2

echo "[2/4] 备份将覆盖的二进制"
if [[ -f "$BIN" ]]; then
  cp -a "$BIN" "$BIN_BAK"
  echo "  备份: $BIN_BAK"
else
  echo "  目标尚无 ism_server，跳过备份"
fi
if [[ -d "$DIST" ]] && [[ -f "$DIST/index.html" ]]; then
  mv "$DIST" "$DIST_BAK"
  echo "  备份: $DIST_BAK"
fi

echo "[3/4] 覆盖 ism_server + web/dist"
mkdir -p "$TARGET/ism_server_user" "$TARGET/web"
cp -a "$PKG/ism_server_user/ism_server" "$BIN"
chmod 755 "$BIN"
mkdir -p "$DIST"
if command -v rsync >/dev/null 2>&1; then
  rsync -a "$PKG/web/dist/" "$DIST/"
else
  cp -a "$PKG/web/dist/." "$DIST/"
fi
[[ -f "$DIST/index.html" ]] || { echo "错误: dist 覆盖后缺少 index.html"; exit 1; }
[[ -x "$BIN" ]] || { echo "错误: ism_server 不可执行"; exit 1; }

trap - ERR

echo "[4/4] 覆盖完成，请重启"
echo "  cd \"$TARGET\" && bash start-all.sh"
echo "  浏览器 Ctrl+F5"
echo ""
echo "回滚:"
echo "  cp -a \"$BIN_BAK\" \"$BIN\""
echo "  rm -rf \"$DIST\" && mv \"$DIST_BAK\" \"$DIST\""
echo "验证通过后可删备份以省空间: $DIST_BAK"
UPGRADE
chmod +x "$STAGING/upgrade-existing.sh"

cat > "$STAGING/install.sh" << 'INSTALL'
#!/bin/bash
# 同一份完整包：新装 或 覆盖已有安装。
# 新装:   sudo bash install.sh
# 覆盖:   bash install.sh --upgrade /opt/ISM/ism-release-oceanbase-20260817-0001-c851
set -euo pipefail

PKG="$(cd "$(dirname "$0")" && pwd)"

usage() {
  echo "用法:"
  echo "  sudo bash install.sh"
  echo "      在本目录新装（无 Docker 则离线装 Docker，再 start-all）"
  echo "  bash install.sh --upgrade <已有ISM目录>"
  echo "      只覆盖 ism_server + web/dist，不碰现场数据/配置"
  echo "示例:"
  echo "  bash install.sh --upgrade /opt/ISM/ism-release-oceanbase-20260817-0001-c851"
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage
  exit 0
fi

if [[ "${1:-}" == "--upgrade" ]]; then
  TARGET="${2:-}"
  [[ -n "$TARGET" ]] || { usage; exit 1; }
  bash "$PKG/upgrade-existing.sh" "$TARGET"
  exit 0
fi

if [[ -n "${1:-}" ]]; then
  if [[ -d "$1/ism_server_user" && -f "$1/start-all.sh" ]]; then
    echo "检测到已有 ISM 安装: $1"
    echo "覆盖升级请用:"
    echo "  bash install.sh --upgrade $1"
    echo "禁止把本包直接解压覆盖该目录（会碰到 conf/data）。"
    exit 1
  fi
  usage
  exit 1
fi

echo "=== ISM 新装（本目录）==="
echo "  目录: $PKG"
echo "  将启动 OceanBase + TDengine + 后端 + 前端"
echo "  首次会导入包内出厂 SQL（data/source/），不是现场库"
echo ""
if [[ -x "$PKG/deploy-offline.sh" ]]; then
  bash "$PKG/deploy-offline.sh"
else
  bash "$PKG/start-all.sh"
fi
INSTALL
chmod +x "$STAGING/install.sh"

break_link "$STAGING/README-部署说明.md"
cat > "$STAGING/README-部署说明.md" << READMEEOF
# ISM OceanBase 完整发布包（${PKG_NAME}）

这是**完整发布包**，不是 \`ism-patch-kylin-*\` 补丁。

- 与主包 \`ism-release-oceanbase-20260817-0001-c851\` **同族**（麒麟 V10 + OceanBase + TDengine 离线一体）
- **一份 zip，两种用法**：全新机器整包安装；已有系统只覆盖程序/前端
- 版本基线: V3.01.RC07（0817）+ 0819/0820 五项 + 20260825 + **20260826 问题项**
- 业务库: OceanBase（\`dbtype=4\`）
- 历史库: TDengine（REST 6041 / 原生 6030）
- 默认端口: 前端 **7090** / 后端 **8091** / OceanBase **2881** / TDengine **6041**

**不要和补丁包混用**：不必先装 0817 再打 0819/0820 补丁。本包已含今日程序。若现场已是 0817 且打过补丁，用「覆盖升级」即可，效果与再打补丁相同，但带备份与回滚。

## 禁止

- **禁止**把本 zip \`unzip -o\` 直接覆盖已有安装目录（会碰到 \`conf/\`、\`data/\`）
- **禁止**覆盖/清空现场 OceanBase、TDengine 数据、\`data/\`、许可证、证书
- **禁止**把「备份上传」当成还原
- 新装与覆盖都**不要**改密码链、不要停 OceanBase/TDengine 容器（覆盖脚本默认不停容器）

## 覆盖哪些 / 绝对不覆盖

| 操作 | 路径 |
|------|------|
| **覆盖**（升级时） | \`ism_server_user/ism_server\`、\`web/dist\` |
| **新装才写入** | 默认 \`conf/\`（仅目标没有现场 conf 时）、启动脚本、离线 Docker/Python/镜像、出厂 SQL |
| **绝对不覆盖** | \`ism_server_user/data/\`、\`logs/\`、已有 \`conf/app.conf\` 及现场改过的 conf、证书 \`conf/*.crt *.key\`、许可证、\`data/db/\`、上传组态/资源、OceanBase 容器数据、TDengine \`tdengine/data\` |

升级脚本会先把 \`ism_server\` 拷成 \`ism_server.bak-YYYYMMDD-HHMM\`，并把旧 \`web/dist\` 改名为 \`web/dist.bak-YYYYMMDD-HHMM\`。失败自动回滚。

## 用法一：全新机器整包安装

依赖：麒麟 V10 SP3 x86_64。OceanBase / TDengine 由本包 Docker 镜像拉起；**不要**把现场库打进包，现场库在本机 Docker 卷里。

\`\`\`bash
# 1. 解压到独立目录（不要解压进已有 ISM 目录）
mkdir -p /opt/ISM
cd /opt/ISM
unzip -o /path/to/${PKG_NAME}.zip
cd ${PKG_NAME}

# 2. （可选）改端口: 编辑 ports.env

# 3. 一键安装并启动
sudo bash install.sh
# 等价：
#   无 Docker: sudo bash deploy-offline.sh
#   已有 Docker: sudo bash start-all.sh
\`\`\`

访问: \`http://<本机IP>:7090/#/login\`  默认账号 \`admin / 123456\`（仅出厂库；现场库沿用现场账号）。

\`start-all.sh\` 会：启动 OceanBase（首次导入包内出厂 SQL）→ 启动 TDengine 并预建 \`ISMHistoryDb\` → 启动后端/前端。

## 用法二：覆盖已有系统（现场 0817 主包）

目标示例：\`/opt/ISM/ism-release-oceanbase-20260817-0001-c851\`

\`\`\`bash
# 1. 解压到临时目录（不要解压到现网目录上）
mkdir -p /tmp/ism-full-20260821
cd /tmp/ism-full-20260821
unzip -o /path/to/${PKG_NAME}.zip
cd ${PKG_NAME}

# 2. 覆盖程序（自动停应用、备份二进制、不碰库）
bash upgrade-existing.sh /opt/ISM/ism-release-oceanbase-20260817-0001-c851
# 或: bash install.sh --upgrade /opt/ISM/ism-release-oceanbase-20260817-0001-c851

# 3. 重启（在现网目录）
cd /opt/ISM/ism-release-oceanbase-20260817-0001-c851
bash start-all.sh
\`\`\`

浏览器 **Ctrl+F5**。验证通过后可删除 \`web/dist.bak-*\` 省空间。

## 本次相对 0817 含哪些修复（含 20260826）

**20260826（本包重点）**

1. **按位脚本 -1**：源键按最后一个 \`->\` 切分虚拟柜名；skip 日志带 device/point/aliases。
2. **历史入库 TAGS**：固定 \`TAGS(1)\` + 失败按行重试；不再随机 TAG 整批失败。
3. **入库可观测**：快照 tick 改为 Error，\`loglevel=3\` 也能看到 wrote / noRealtimeValue。
4. **导航绑模型 UUID**：\`pageid===displayUUID\` 时跳该模型首页图层。
5. **Excel 五种存储**：变化 / 定时 / 即时 / 变化百分比 / 整点；触发值列保留。
6. **历史备份超时**：15 分钟超时文案；与业务库备份页区分。

**20260825**

1. **数据仓库设备分类**：撤掉列头名称筛选；测点显示库内全名；选中真设备后出 \`全部 / A列头 / 列尾\` Tag 墙。
2. **组态切页保存**：保存冻结 pageId+快照；切页等待队列，避免改 A 切 B 再回 A 丢失。
3. **运行态导航**：未绑定 →「菜单未绑定页面」；找不到页 →「找不到页面」+ pageUuid。
4. **历史备份报错**：本机 taosdump 优先；docker 失败带容器名 / inspect，不只 \`exit status 125\`。
5. **虚拟设备 Excel**：\`GetCellValue\` 读报警触发值，\`0\` 不再被当成空写成 \`1\`；按数据 ID 更新。
6. **定时入库日志**：\`points loaded: total= N type1= x type4= y\`。

**此前已含（0819/0820）**

- 组态 \`resolveSavePageId\`、animateType fallback
- 历史库备份 taosdump 参数、TDengine TAGS/快照、脚本 Settle、Excel \`includeEmpty\`

另：备份上传 ≠ 还原。**不要**用列头放大镜当分类验收。

## 包内目录（与 0817 对齐）

\`\`\`
${PKG_NAME}/
├── install.sh                   # 新装（无参）/ 覆盖（--upgrade）
├── upgrade-existing.sh          # 覆盖已有安装
├── start-all.sh / stop-all.sh / deploy-offline.sh
├── ports.env
├── docker-compose.oceanbase.yml
├── docker-compose.tdengine.yml
├── oceanbase/oceanbase-ce.tar   # 出厂镜像，不是现场数据
├── tdengine/tdengine.tar
├── docker-offline/ python-offline/
├── ism_server_user/             # 后端（conf 为 dbtype=4 模板）
├── web/dist/                    # 前端（今日 production build）
├── data/source/Mysql_Backup_*.sql   # 出厂 SQL，仅新装空库导入
└── scripts/
\`\`\`

## 未包含 / 现场需已有

- **覆盖升级**时，现场必须已有 OceanBase、TDengine（0817 主包已装即可）
- **新装**时本包自带 Docker 离线组件与 OB/TD 镜像；仍需现场磁盘/内存满足 0817 手册（OB mini 约 8G 内存）
- 不包含现场业务库、历史库、许可证、已改 conf
- zip 的 SHA256 见同目录 \`${PKG_NAME}.zip.sha256\`

详细文档见 \`docs-ISM-OceanBase部署与切换指南.md\`、\`ISM-麒麟V10-OceanBase部署操作手册.md\`。
READMEEOF
cp "$STAGING/README-部署说明.md" "$STAGING/README.md"

cat > "$STAGING/BUILD_INFO.txt" << EOF
包名: ${PKG_NAME}
包类型: 完整发布（非 patch）
基线主包: ism-release-oceanbase-20260817-0001-c851
构建时间: $(date '+%Y-%m-%d %H:%M:%S')
git: ${SHORT_SHA}
后端: 复用今日麒麟静态 ism_server (CGO_ENABLED=0 linux/amd64 -mod=vendor 同源)
前端: 今日 production dist（含 resolveSavePageId / includeEmpty）
dbtype: $(grep '^dbtype=' "$STAGING/ism_server_user/conf/app.conf")
historyrecorddbtype: $(grep -i '^historyrecorddbtype=' "$STAGING/ism_server_user/conf/historyData.conf" 2>/dev/null || echo N/A)
目标平台: 麒麟 V10 SP3 x86_64
默认端口: FE=7090 BE=8091 OB=2881 TD=6041
新装入口: bash install.sh  或  sudo bash deploy-offline.sh
覆盖入口: bash upgrade-existing.sh <已有ISM目录>
EOF

chmod +x "$STAGING/scripts/"*.sh 2>/dev/null || true

echo "[5/6] 完整性检查 ..."
for must in \
  "$STAGING/ism_server_user/ism_server" \
  "$STAGING/web/dist/index.html" \
  "$STAGING/start-all.sh" \
  "$STAGING/stop-all.sh" \
  "$STAGING/deploy-offline.sh" \
  "$STAGING/install.sh" \
  "$STAGING/upgrade-existing.sh" \
  "$STAGING/ports.env" \
  "$STAGING/oceanbase/oceanbase-ce.tar" \
  "$STAGING/tdengine/tdengine.tar" \
  "$STAGING/docker-offline/bin/dockerd" \
  "$STAGING/ism_server_user/conf/app.conf"
do
  [[ -e "$must" ]] || { echo "错误: 缺少 $must"; exit 1; }
done
verify_kylin_binary "$STAGING/ism_server_user/ism_server" || { echo "错误: 包内二进制不兼容麒麟"; exit 1; }
grep -q '^dbtype=4' "$STAGING/ism_server_user/conf/app.conf"
echo "  OK"

echo "[6/6] 压缩 zip（约数分钟到十几分钟）..."
rm -f "$ZIP_OUT"
(cd "$ROOT/releases" && COPYFILE_DISABLE=1 zip -r -y -q "$(basename "$ZIP_OUT")" "$(basename "$STAGING")")

HASH="$(shasum -a 256 "$ZIP_OUT" | awk '{print $1}')"
echo "${HASH}  $(basename "$ZIP_OUT")" > "${ZIP_OUT}.sha256"
printf '%s\n' "$HASH" > "${ZIP_OUT}.sha256.txt"

echo ""
echo "=== 构建完成 ==="
ls -lh "$ZIP_OUT" "${ZIP_OUT}.sha256"
echo "SHA256: $HASH"
echo "目录: $STAGING"
echo "zip:  $ZIP_OUT"
