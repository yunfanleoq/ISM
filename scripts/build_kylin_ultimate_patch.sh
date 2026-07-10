#!/bin/bash
# 构建「最全」麒麟补丁包：
#   - 双库离线镜像：OceanBase + TDengine
#   - 双端：ism_server（麒麟静态）+ web/dist
#   - 配置/脚本/一键 apply + 部署手册
# 用法:
#   bash scripts/build_kylin_ultimate_patch.sh
#   SKIP_REBUILD=1 bash scripts/build_kylin_ultimate_patch.sh   # 跳过前后端重编
#   SKIP_DB_IMAGES=1 ...                                       # 不打入 OB/TD 镜像（体积小）
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
# shellcheck source=lib_build_id.sh
source "$ROOT/scripts/lib_build_id.sh"
BUILD_ID="$(ism_build_id)"
PKG="ism-patch-kylin-ultimate-${BUILD_ID}"
STAGING="$ROOT/releases/${PKG}"
ZIP="$ROOT/releases/${PKG}.zip"
BIN_SRC="$ROOT/patches/ism-server-kylin-glibc228/ism_server"
FE_SRC="$ROOT/ism-front-end-v2/dist"
REL09="$ROOT/releases/ism-release-oceanbase-20260709"
TD_PATCH="$ROOT/releases/ism-patch-kylin-tdengine-20260709"

echo "=== ISM 麒麟最全补丁构建 (${BUILD_ID}) ==="

# ── 1. 后端 ──
echo ""
if [[ "${SKIP_REBUILD:-0}" == "1" ]] && [[ -f "$BIN_SRC" ]]; then
  echo "[1/5] 跳过后端编译 (SKIP_REBUILD=1) ..."
else
  echo "[1/5] 编译 ism_server (linux/amd64 静态) ..."
  bash "$ROOT/scripts/build_kylin_ism_server.sh"
fi

# ── 2. 前端 ──
echo ""
if [[ "${SKIP_REBUILD:-0}" == "1" ]] && [[ -f "$FE_SRC/index.html" ]]; then
  echo "[2/5] 跳过前端编译 (SKIP_REBUILD=1) ..."
else
  echo "[2/5] 编译前端 dist ..."
  if [[ -x "$ROOT/scripts/check_mem_before_compile.sh" ]]; then
    if ! bash "$ROOT/scripts/check_mem_before_compile.sh"; then
      echo "错误: 内存检查 FAIL，禁止编译前端"
      exit 1
    fi
  fi
  # 清场旧 vue-cli，避免内存叠加
  pkill -9 -f "vue-cli-service" 2>/dev/null || true
  sleep 2
  rm -rf "$FE_SRC"
  (
    cd "$ROOT/ism-front-end-v2"
    export NODE_OPTIONS="--max-old-space-size=20480 --openssl-legacy-provider"
    npm run build
  )
  [[ -d "$FE_SRC" && -f "$FE_SRC/index.html" ]] || { echo "错误: dist 未生成"; exit 1; }
  echo "  dist: $(du -sh "$FE_SRC" | cut -f1)"
fi

# ── 3. 组装 ──
echo ""
echo "[3/5] 组装最全补丁目录 ..."
rm -rf "$STAGING"
mkdir -p \
  "$STAGING/ism_server_user/conf" \
  "$STAGING/web/dist" \
  "$STAGING/scripts" \
  "$STAGING/oceanbase" \
  "$STAGING/tdengine" \
  "$STAGING/data/source" \
  "$STAGING/docs"

# 后端
cp "$BIN_SRC" "$STAGING/ism_server_user/ism_server"
chmod 755 "$STAGING/ism_server_user/ism_server"

# 配置：先同步主包完整 conf，再覆盖关键项（避免缺 mqtt_broken_config.json / videoConfig.json 导致 panic）
if [[ -d "$REL09/ism_server_user/conf" ]]; then
  rsync -a "$REL09/ism_server_user/conf/" "$STAGING/ism_server_user/conf/"
elif [[ -d "$ROOT/ism_server_user/conf" ]]; then
  rsync -a "$ROOT/ism_server_user/conf/" "$STAGING/ism_server_user/conf/"
fi
# 源码侧关键 conf 覆盖（保证最新）
for f in mqtt_broken_config.json videoConfig.json historyData.conf mqtt.conf \
         opcuaserver.conf sim.conf systimeconfig.conf MenuConfig.json \
         MonibucaServer.yaml NetWorker.json Speeker.conf ISMNodeConfig.conf; do
  [[ -f "$ROOT/ism_server_user/conf/$f" ]] && cp "$ROOT/ism_server_user/conf/$f" "$STAGING/ism_server_user/conf/"
done
if [[ -f "$ROOT/ism_server_user/conf/app.conf" ]]; then
  cp "$ROOT/ism_server_user/conf/app.conf" "$STAGING/ism_server_user/conf/app.conf.sample"
fi
# 证书目录
if [[ -d "$REL09/ism_server_user/conf/x509" ]]; then
  rsync -a "$REL09/ism_server_user/conf/x509/" "$STAGING/ism_server_user/conf/x509/"
fi
for cert in 192.168.199.120.crt 192.168.199.120.key; do
  [[ -f "$REL09/ism_server_user/conf/$cert" ]] && cp "$REL09/ism_server_user/conf/$cert" "$STAGING/ism_server_user/conf/"
done

# 确保正式环境日志/历史缓冲项存在
CONF="$STAGING/ism_server_user/conf/app.conf"
[[ -f "$CONF" ]] || cp "$ROOT/ism_server_user/conf/app.conf" "$CONF"
grep -q '^loglevel=' "$CONF" || cat >> "$CONF" <<'EOF'

# 正式环境写盘收敛
loglevel=3
logFilesSavaDays=2
logmaxsize_mb=20
log_throttle_seconds=60
HistoryDataBufferSize=10000
HistoryDataFlushInterval=2000
EOF
# 一体包业务库必须是 OceanBase
if grep -q '^dbtype=' "$CONF"; then
  if [[ "$(uname -s)" == "Darwin" ]]; then
    sed -i '' 's/^dbtype=.*/dbtype=4/' "$CONF"
  else
    sed -i 's/^dbtype=.*/dbtype=4/' "$CONF"
  fi
fi
# 正式环境默认关闭内置 MQTT Broker / HTTPS（缺证书会导致 ListenAndServeTLS 拖垮进程）
for kv in enablemqttbreoken=false enablehttps=false; do
  key="${kv%%=*}"
  if grep -q "^${key}=" "$CONF"; then
    if [[ "$(uname -s)" == "Darwin" ]]; then
      sed -i '' "s/^${key}=.*/${kv}/" "$CONF"
    else
      sed -i "s/^${key}=.*/${kv}/" "$CONF"
    fi
  else
    echo "$kv" >> "$CONF"
  fi
done
# 系统脚本目录（缺失只打 Error，但一并带上避免噪音）
mkdir -p "$STAGING/ism_server_user/sys_script"
if [[ -d "$ROOT/ism_server_user/sys_script" ]]; then
  rsync -a "$ROOT/ism_server_user/sys_script/" "$STAGING/ism_server_user/sys_script/" || true
fi
touch "$STAGING/ism_server_user/sys_script/.gitkeep"

HC="$STAGING/ism_server_user/conf/historyData.conf"
[[ -f "$HC" ]] || cp "$ROOT/ism_server_user/conf/historyData.conf" "$HC"
if [[ "$(uname -s)" == "Darwin" ]]; then
  sed -i '' 's/^historyrecorddbtype=.*/historyrecorddbtype=2/' "$HC"
  sed -i '' 's/^tdenginehost=.*/tdenginehost=127.0.0.1/' "$HC"
  sed -i '' 's/^tdengineport=.*/tdengineport=6041/' "$HC"
else
  sed -i 's/^historyrecorddbtype=.*/historyrecorddbtype=2/' "$HC"
  sed -i 's/^tdenginehost=.*/tdenginehost=127.0.0.1/' "$HC"
  sed -i 's/^tdengineport=.*/tdengineport=6041/' "$HC"
fi

# 前端
rsync -a --delete "$FE_SRC/" "$STAGING/web/dist/"

# 双库离线镜像
if [[ "${SKIP_DB_IMAGES:-0}" != "1" ]]; then
  if [[ -f "$REL09/oceanbase/oceanbase-ce.tar" ]]; then
    cp "$REL09/oceanbase/oceanbase-ce.tar" "$STAGING/oceanbase/"
    echo "oceanbase/oceanbase-ce:latest" > "$STAGING/oceanbase/IMAGE_TAG"
  else
    echo "警告: 缺少 OceanBase 镜像 $REL09/oceanbase/oceanbase-ce.tar"
  fi
  # 优先用完整 TD 镜像（~469M）
  if [[ -f "$TD_PATCH/tdengine/tdengine.tar" ]]; then
    cp "$TD_PATCH/tdengine/tdengine.tar" "$STAGING/tdengine/"
    cp "$TD_PATCH/tdengine/IMAGE_TAG" "$STAGING/tdengine/" 2>/dev/null || \
      echo "tdengine/tdengine:3.3.6.13" > "$STAGING/tdengine/IMAGE_TAG"
  elif [[ -f "$REL09/tdengine/tdengine.tar" ]]; then
    cp "$REL09/tdengine/tdengine.tar" "$STAGING/tdengine/"
    cp "$REL09/tdengine/IMAGE_TAG" "$STAGING/tdengine/" 2>/dev/null || true
  else
    echo "警告: 缺少 TDengine 镜像"
  fi
else
  echo "  SKIP_DB_IMAGES=1，跳过 OB/TD 镜像拷贝"
fi

# 业务 SQL（若有）
if [[ -d "$REL09/data/source" ]]; then
  rsync -a "$REL09/data/source/" "$STAGING/data/source/" || true
fi

# compose / ports / 启停脚本
for f in docker-compose.oceanbase.yml docker-compose.tdengine.yml ports.env deploy-offline.sh start-all.sh stop-all.sh; do
  if [[ -f "$REL09/$f" ]]; then
    cp "$REL09/$f" "$STAGING/"
  elif [[ -f "$TD_PATCH/$f" ]]; then
    cp "$TD_PATCH/$f" "$STAGING/"
  fi
done
chmod +x "$STAGING/"*.sh 2>/dev/null || true
[[ -f "$STAGING/ports.env" ]] || cp "$TD_PATCH/ports.env.sample" "$STAGING/ports.env" 2>/dev/null || true

# 诊断/修复脚本
for s in \
  fix_admin_password_oceanbase.sh check_dw_device_loading.sh \
  check_login_deep.sh check_login_and_user.sh collect_diagnose_log.sh \
  run_full_field_check.sh check_env_kylin.sh \
  diagnose_getrealdata_timeout.sh fix_device_real_data_index.sh \
  init_tdengine.sh init_oceanbase.sh import_mysql_to_oceanbase.sh \
  ensure_python.sh serve_test_frontend.py diagnose_kylin.sh \
  fix_compose_offline.sh install_docker_kylin_sp3.sh; do
  if [[ -f "$ROOT/scripts/$s" ]]; then
    cp "$ROOT/scripts/$s" "$STAGING/scripts/"
    [[ "$s" == *.sh ]] && chmod +x "$STAGING/scripts/$s"
  elif [[ -f "$REL09/scripts/$s" ]]; then
    cp "$REL09/scripts/$s" "$STAGING/scripts/"
    [[ "$s" == *.sh ]] && chmod +x "$STAGING/scripts/$s"
  elif [[ -f "$TD_PATCH/scripts/$s" ]]; then
    cp "$TD_PATCH/scripts/$s" "$STAGING/scripts/"
    chmod +x "$STAGING/scripts/$s" 2>/dev/null || true
  fi
done

# 文档
if [[ -f "$ROOT/docs/ISM-麒麟V10-OceanBase部署操作手册.md" ]]; then
  cp "$ROOT/docs/ISM-麒麟V10-OceanBase部署操作手册.md" "$STAGING/docs/"
fi
if [[ -f "$REL09/ISM-麒麟V10-OceanBase部署操作手册.pdf" ]]; then
  cp "$REL09/ISM-麒麟V10-OceanBase部署操作手册.pdf" "$STAGING/docs/" 2>/dev/null || true
fi
if [[ -f "$REL09/docs-ISM-OceanBase部署与切换指南.md" ]]; then
  cp "$REL09/docs-ISM-OceanBase部署与切换指南.md" "$STAGING/docs/"
fi

# ── 校验 ──
echo ""
echo "[校验] 产物自检 ..."
if rg -q 'sockjs-node|webpack-dev-server/client' "$FE_SRC/static/js/"*.js 2>/dev/null; then
  echo "错误: dist 含 dev-server 热更新代码，禁止出包"
  exit 1
fi
set +o pipefail
BIN_OK=0
if strings "$BIN_SRC" 2>/dev/null | rg -q 'GetRealDataPaged|ErrorThrottled|ModbusReconnectSleep|userTable'; then
  BIN_OK=1
fi
BIN_MQTT_OK=0
if strings "$BIN_SRC" 2>/dev/null | rg -q 'mqtt_broken_config\.json'; then
  BIN_MQTT_OK=1
fi
BIN_LICENSE_OK=1
if strings "$BIN_SRC" 2>/dev/null | rg -q '目前使用的是个人免费版本'; then
  BIN_LICENSE_OK=0
fi
set -o pipefail
if [[ "$BIN_OK" -ne 1 ]]; then
  echo "错误: ism_server 未包含写盘节流/分页/登录修复符号"
  exit 1
fi
if [[ "$BIN_MQTT_OK" -ne 1 ]]; then
  echo "错误: ism_server 未包含 MQTT Broker 缺配置防护"
  exit 1
fi
if [[ "$BIN_LICENSE_OK" -ne 1 ]]; then
  echo "错误: ism_server 仍含「个人免费版本」提示，请先重编含企业授权默认的二进制"
  exit 1
fi
if [[ ! -f "$STAGING/ism_server_user/conf/mqtt_broken_config.json" ]]; then
  echo "错误: 补丁缺少 conf/mqtt_broken_config.json"
  exit 1
fi
if [[ ! -f "$STAGING/ism_server_user/conf/videoConfig.json" ]]; then
  echo "错误: 补丁缺少 conf/videoConfig.json"
  exit 1
fi
if ! rg -q 'realDataBatch|REAL_DATA_DEFAULT_PAGE_SIZE|fetchRealDataByUuidBatched|getRealData\.loadTimeout' "$FE_SRC/static/js/" 2>/dev/null; then
  echo "错误: dist 未包含分页瘦身/Loading 兜底编译产物"
  exit 1
fi
echo "  ✓ 后端: GetRealDataPaged / ErrorThrottled / ModbusReconnectSleep / userTable / MQTT panic 防护 / 企业授权默认"
echo "  ✓ conf: mqtt_broken_config.json + videoConfig.json + enablemqttbreoken=false + enablehttps=false"
echo "  ✓ 前端: realDataBatch 分页 + Loading 兜底"
[[ -f "$STAGING/oceanbase/oceanbase-ce.tar" ]] && echo "  ✓ OceanBase 镜像 $(du -sh "$STAGING/oceanbase/oceanbase-ce.tar" | cut -f1)" || echo "  ⚠ 无 OceanBase 镜像"
[[ -f "$STAGING/tdengine/tdengine.tar" ]] && echo "  ✓ TDengine 镜像 $(du -sh "$STAGING/tdengine/tdengine.tar" | cut -f1)" || echo "  ⚠ 无 TDengine 镜像"

# ── apply-patch.sh ──
cat > "$STAGING/apply-patch.sh" << 'APPLY'
#!/bin/bash
# 最全补丁一键应用
# 用法: bash apply-patch.sh /opt/ISM/ism-release-oceanbase-20260709
set -euo pipefail
PATCH_ROOT="$(cd "$(dirname "$0")" && pwd)"
TARGET="${1:?用法: bash apply-patch.sh <主包目录>}"
TARGET="$(cd "$TARGET" && pwd)"
[[ -d "$TARGET/ism_server_user" ]] || { echo "错误: 无效目录 $TARGET"; exit 1; }

echo "=== ISM 最全补丁应用 → $TARGET ==="

echo "=== [1/8] 停止服务 ==="
(cd "$TARGET" && bash stop-all.sh 2>/dev/null) || true
sleep 3

echo "=== [2/8] 替换 ism_server ==="
cp "$PATCH_ROOT/ism_server_user/ism_server" "$TARGET/ism_server_user/ism_server"
chmod 755 "$TARGET/ism_server_user/ism_server"
file "$TARGET/ism_server_user/ism_server" || true

echo "=== [3/8] 合并正式环境配置（写盘收敛 + 历史库 + 防 panic 配置文件）==="
mkdir -p "$TARGET/ism_server_user/conf"
# 补齐缺失 conf（不覆盖现场已有自定义文件，除非是空/缺失）
for f in mqtt_broken_config.json videoConfig.json historyData.conf mqtt.conf \
         opcuaserver.conf sim.conf systimeconfig.conf MenuConfig.json \
         MonibucaServer.yaml NetWorker.json Speeker.conf ISMNodeConfig.conf \
         192.168.199.120.crt 192.168.199.120.key; do
  if [[ -f "$PATCH_ROOT/ism_server_user/conf/$f" ]]; then
    if [[ ! -f "$TARGET/ism_server_user/conf/$f" ]]; then
      cp "$PATCH_ROOT/ism_server_user/conf/$f" "$TARGET/ism_server_user/conf/$f"
      echo "  + 补齐 conf/$f"
    fi
  fi
done
# historyData 强制指向本机 TD（覆盖）
if [[ -f "$PATCH_ROOT/ism_server_user/conf/historyData.conf" ]]; then
  cp "$PATCH_ROOT/ism_server_user/conf/historyData.conf" "$TARGET/ism_server_user/conf/historyData.conf"
fi
CONF="$TARGET/ism_server_user/conf/app.conf"
if [[ -f "$CONF" ]]; then
  for kv in \
    "loglevel=3" \
    "logFilesSavaDays=2" \
    "logmaxsize_mb=20" \
    "log_throttle_seconds=60" \
    "HistoryDataBufferSize=10000" \
    "HistoryDataFlushInterval=2000" \
    "enablemqttbreoken=false" \
    "enablehttps=false"; do
    key="${kv%%=*}"
    if grep -q "^${key}=" "$CONF" 2>/dev/null; then
      sed -i "s/^${key}=.*/${kv}/" "$CONF"
    else
      echo "$kv" >> "$CONF"
    fi
  done
else
  cp "$PATCH_ROOT/ism_server_user/conf/app.conf" "$CONF"
fi
# 强制补齐 HTTPS 证书文件（即使关闭 HTTPS，也避免现场再打开时崩）
for f in 192.168.199.120.crt 192.168.199.120.key opcuaserver.conf; do
  if [[ -f "$PATCH_ROOT/ism_server_user/conf/$f" ]]; then
    cp "$PATCH_ROOT/ism_server_user/conf/$f" "$TARGET/ism_server_user/conf/$f"
    echo "  + conf/$f"
  fi
done
mkdir -p "$TARGET/ism_server_user/sys_script"
if [[ -d "$PATCH_ROOT/ism_server_user/sys_script" ]]; then
  rsync -a "$PATCH_ROOT/ism_server_user/sys_script/" "$TARGET/ism_server_user/sys_script/" || true
fi

echo "=== [4/8] 替换前端 web/dist ==="
mkdir -p "$TARGET/web/dist"
rsync -a --delete "$PATCH_ROOT/web/dist/" "$TARGET/web/dist/"

echo "=== [5/8] 同步脚本与启停 ==="
mkdir -p "$TARGET/scripts"
cp "$PATCH_ROOT/scripts/"* "$TARGET/scripts/" 2>/dev/null || true
chmod +x "$TARGET/scripts/"*.sh 2>/dev/null || true
for f in start-all.sh stop-all.sh deploy-offline.sh docker-compose.oceanbase.yml docker-compose.tdengine.yml; do
  [[ -f "$PATCH_ROOT/$f" ]] && cp "$PATCH_ROOT/$f" "$TARGET/$f"
done
chmod +x "$TARGET/start-all.sh" "$TARGET/stop-all.sh" "$TARGET/deploy-offline.sh" 2>/dev/null || true
if [[ -f "$PATCH_ROOT/ports.env" ]] && [[ ! -f "$TARGET/ports.env" ]]; then
  cp "$PATCH_ROOT/ports.env" "$TARGET/ports.env"
fi
# 合并 TD 端口到 ports.env
if [[ -f "$TARGET/ports.env" ]] && ! grep -q '^TD_PORT=' "$TARGET/ports.env" 2>/dev/null; then
  cat >> "$TARGET/ports.env" <<'PORTEOF'

# TDengine 历史库
TD_PORT=6041
TD_NATIVE_PORT=6030
TD_USER=root
TD_PASSWORD=taosdata
TD_CONTAINER=tdengine
TD_IMAGE=tdengine/tdengine:3.3.6.13
PORTEOF
fi

echo "=== [6/8] 同步双库离线镜像（若补丁内有）==="
if [[ -f "$PATCH_ROOT/oceanbase/oceanbase-ce.tar" ]]; then
  mkdir -p "$TARGET/oceanbase"
  rsync -a "$PATCH_ROOT/oceanbase/" "$TARGET/oceanbase/"
  echo "  OceanBase 镜像已同步"
fi
if [[ -f "$PATCH_ROOT/tdengine/tdengine.tar" ]]; then
  mkdir -p "$TARGET/tdengine"
  rsync -a "$PATCH_ROOT/tdengine/" "$TARGET/tdengine/"
  echo "  TDengine 镜像已同步"
fi
if [[ -d "$PATCH_ROOT/data/source" ]]; then
  mkdir -p "$TARGET/data/source"
  rsync -a "$PATCH_ROOT/data/source/" "$TARGET/data/source/" || true
fi
if [[ -d "$PATCH_ROOT/docs" ]]; then
  mkdir -p "$TARGET/docs"
  rsync -a "$PATCH_ROOT/docs/" "$TARGET/docs/" || true
  cp "$PATCH_ROOT/README-部署手册.md" "$TARGET/docs/" 2>/dev/null || true
fi

echo "=== [7/8] 修复 device_real_data 索引 ==="
if [[ -x "$TARGET/scripts/fix_device_real_data_index.sh" ]]; then
  (cd "$TARGET" && bash scripts/fix_device_real_data_index.sh) || \
    echo "  [WARN] 索引脚本失败（可稍后手工执行），继续"
fi

echo "=== [8/8] 启动并验证 ==="
(cd "$TARGET" && bash start-all.sh)
echo "  等待 90 秒 ..."
sleep 90

BE_PORT="${ISM_BE_PORT:-8091}"
FE_PORT="${ISM_FE_PORT:-7090}"
echo "--- 登录 ---"
curl -s --compressed -m 15 -X POST "http://127.0.0.1:${BE_PORT}/login" \
  -H 'Content-Type: application/json' \
  -d '{"Username":"admin","password":"e10adc3949ba59abbe56e057f20f883e"}' || true
echo ""
echo "--- 端口 ---"
ss -lntp 2>/dev/null | grep -E ":${BE_PORT}|:${FE_PORT}|:2881|:6041" || true
echo ""
echo "--- TDengine ---"
curl -s -u root:taosdata -d "show databases;" http://127.0.0.1:6041/rest/sql 2>/dev/null | head -c 400 || echo "(TD 未就绪可稍后重试)"
echo ""
echo "期望: login code=1000；8091/7090/2881/6041 监听"
echo "浏览器 Ctrl+F5 后测数据仓库分页（默认 30 行）与大屏"
APPLY
chmod +x "$STAGING/apply-patch.sh"

# ── 部署手册 ──
cat > "$STAGING/README-部署手册.md" << README
# ISM 麒麟 V10 · 最全补丁包部署手册

构建 ID: **${BUILD_ID}**  
适用主包: \`ism-release-oceanbase-20260709\` / \`20260708\`（同结构 OceanBase 一体包）  
版本: V3.01.RC07

---

## 1. 本包包含什么

| 类别 | 内容 | 说明 |
|------|------|------|
| **业务库** | \`oceanbase/oceanbase-ce.tar\` | OceanBase CE 离线镜像（\`dbtype=4\`） |
| **历史库** | \`tdengine/tdengine.tar\` | TDengine 3.3.x 离线镜像（\`historyrecorddbtype=2\`，REST 6041） |
| **后端** | \`ism_server_user/ism_server\` | 麒麟静态；写盘节流 + 强制分页 + MQTT/HTTPS 防护 + **源码企业版默认已授权** |
| **前端** | \`web/dist/\` | 生产构建；数据仓库/大屏分批 20~50（上限 100） |
| **配置** | 完整 \`conf/\`（含 mqtt/video/history） | 日志收敛、TD 本机、\`enablemqttbreoken=false\` |
| **脚本** | \`apply-patch.sh\` / \`start-all.sh\` / 诊断脚本 | 一键应用与排障 |
| **数据** | \`data/source/*.sql\` | 业务初始 SQL（若主包无数据可导入） |
| **文档** | \`docs/\` + 本手册 | 部署与切换说明 |

默认端口：前端 **7090** / 后端 **8091** / OceanBase **2881** / TDengine **6041**  
登录：**admin / 123456**（前端 MD5 后登录）

---

## 2. 本包相对旧现场的关键修复

| # | 问题 | 修复 |
|---|------|------|
| 1 | **ism_server panic 崩溃** | 缺 \`mqtt_broken_config.json\` 时空指针；现安全返回 + 默认关闭内置 MQTT Broker |
| 1b | **HTTPS 证书缺失拖垮进程** | \`enablehttps=true\` 但缺 crt/key → ListenAndServeTLS 失败退出；现默认 \`enablehttps=false\` 并补齐证书 |
| 1c | **「个人免费版本」授权提示** | 源码企业版默认 \`IsLicense/IsOem=true\`，不依赖 \`license.lic\` / \`active.dat\`，不再弹官网 |
| 2 | 磁盘写入 40~80MB/s | 日志节流、单文件 20MB、Modbus 重连≥5s、默认不写 modbus 明细日志 |
| 3 | 上万点一次加载卡死 | \`getRealData\` **强制分页**默认 30、硬上限 100 |
| 4 | 大屏绑点一次打满内存 | 前端 \`realDataBatch\` 分批请求 |
| 5 | 登录 code:1003 | OceanBase \`user\` 表反引号查询 |
| 6 | 数据仓库一直 Loading | 去掉双重 gzip + 前端 15s 强制关 Loading |
| 7 | getRealData 超时 | 分页 + \`device_uuid\` 前缀索引脚本 |
| 8 | 历史库 6041 refused | 包内 TD 镜像 + start-all 自动拉起 |
| 9 | 缺 videoConfig.json | 补齐 \`conf/videoConfig.json\` |

---

## 3. 推荐：叠加到已有一体包（补丁模式）

\`\`\`bash
# 1) 上传 zip 到麒麟服务器
unzip ism-patch-kylin-ultimate-${BUILD_ID}.zip
cd ism-patch-kylin-ultimate-${BUILD_ID}

# 2) 一键应用（会停服 → 替换双端/配置/镜像 → 启服 → 自检）
bash apply-patch.sh /opt/ISM/ism-release-oceanbase-20260709
\`\`\`

应用后浏览器 **Ctrl+F5**（务必强刷）。

---

## 4. 全新离线部署（无旧包时）

若机器上还没有一体包，可把本补丁目录当作「可启动目录」使用（需已有 Docker）：

\`\`\`bash
cd ism-patch-kylin-ultimate-${BUILD_ID}

# 无 Docker 时先装（包内脚本，需 root）
sudo bash scripts/install_docker_kylin_sp3.sh   # 若存在

# 一键：load OB/TD → 导入数据 → 起后端/前端
sudo bash start-all.sh
# 或完全离线入口：
sudo bash deploy-offline.sh
\`\`\`

首次导入业务 SQL 约 **10~15 分钟**。

---

## 5. 应用后验证清单

\`\`\`bash
cd /opt/ISM/ism-release-oceanbase-20260709

# 登录（期望 code=1000）
bash scripts/check_login_deep.sh

# 端口
ss -lntp | grep -E ':8091|:7090|:2881|:6041'

# TDengine
curl -u root:taosdata -d "show databases;" http://127.0.0.1:6041/rest/sql

# 数据仓库分页（期望 <2s，pageSize≤100）
DEVICE_UUID=<有测点设备UUID> bash scripts/diagnose_getrealdata_timeout.sh

# 日志不应再狂涨（同类错误 60s 一条）
ls -lah ism_server_user/logs/
tail -20 logs/ism_server.log
\`\`\`

浏览器：
1. 打开 \`http://<IP>:7090\`，Ctrl+F5
2. 数据仓库选设备 → 表格约 **30 行/页**，翻页流畅
3. 应用大屏 → 首屏绑点分批加载，不应整页卡死

---

## 6. 手工步骤（不用 apply-patch 时）

\`\`\`bash
cd /opt/ISM/ism-release-oceanbase-20260709
bash stop-all.sh

cp <补丁>/ism_server_user/ism_server ism_server_user/
rsync -a --delete <补丁>/web/dist/ web/dist/
cp <补丁>/ism_server_user/conf/historyData.conf ism_server_user/conf/
# 手工把 app.conf 补上 loglevel / logmaxsize_mb / log_throttle_seconds 等

rsync -a <补丁>/oceanbase/ oceanbase/     # 可选
rsync -a <补丁>/tdengine/ tdengine/       # 可选
cp <补丁>/scripts/*.sh scripts/
cp <补丁>/start-all.sh <补丁>/stop-all.sh .

bash scripts/fix_device_real_data_index.sh
bash start-all.sh
\`\`\`

---

## 7. 配置要点（正式环境）

\`app.conf\` 关键项：

\`\`\`ini
dbtype=4
runmode=prod
isdebug=false
loglevel=3
logFilesSavaDays=2
logmaxsize_mb=20
log_throttle_seconds=60
HistoryDataBufferSize=10000
HistoryDataFlushInterval=2000
history_keep_days=7
\`\`\`

\`historyData.conf\`：

\`\`\`ini
historyrecorddbtype=2
[tdengine]
tdenginehost=127.0.0.1
tdengineport=6041
username=root
password=taosdata
\`\`\`

**不要**在正式环境默认启动 \`modbus_simulator.py\`。

---

## 8. 回滚

\`\`\`bash
# 应用前建议：
cp ism_server_user/ism_server ism_server_user/ism_server.bak.\$(date +%Y%m%d)
cp -a web/dist web/dist.bak.\$(date +%Y%m%d)

# 回滚：
bash stop-all.sh
cp ism_server_user/ism_server.bak.YYYYMMDD ism_server_user/ism_server
rsync -a --delete web/dist.bak.YYYYMMDD/ web/dist/
bash start-all.sh
\`\`\`

---

## 9. 包体积与磁盘建议

| 组件 | 约占用 |
|------|--------|
| 前端 dist | ~1.5 GB |
| OceanBase 镜像 | ~480 MB |
| TDengine 镜像 | ~470 MB |
| 后端二进制 | ~65 MB |
| **解压后合计** | **约 2.5~3 GB** |

现场建议预留磁盘 **≥ 20 GB**（含 OB 数据文件与日志）。

---

## 10. 常见问题

| 现象 | 处理 |
|------|------|
| 登录失败 | 先 \`ss\` 看 8091；再 \`check_login_deep.sh\`；勿先改密码 |
| 6041 refused | \`docker ps \| grep tdengine\`；\`bash scripts/init_tdengine.sh\` |
| 磁盘仍很高 | \`iotop -oP\`：若是 observer 属 OB 正常；若是 ism_server 查日志是否仍刷屏 |
| 数据仓库仍一次很多行 | 确认已 Ctrl+F5；接口响应应带 \`pageSize≤100\` |
| glibc 报错 | 本包后端为静态链接，不应依赖目标机高版本 glibc |

更多细节见 \`docs/ISM-麒麟V10-OceanBase部署操作手册.md\`。
README

# ── 4. BUILD_INFO ──
cat > "$STAGING/BUILD_INFO.txt" << INFO
PKG=${PKG}
BUILD_ID=${BUILD_ID}
DATE=$(date '+%Y-%m-%d %H:%M:%S %z')
BACKEND=patches/ism-server-kylin-glibc228/ism_server
FRONTEND=ism-front-end-v2/dist
FEATURES=disk-throttle,forced-pagination-30-100,mqtt-panic-fix,ob+td-offline,full-conf,apply-patch
INFO

# ── 5. zip ──
echo ""
echo "[4/5] 压缩（可能较久）..."
rm -f "$ZIP"
(
  cd "$ROOT/releases"
  COPYFILE_DISABLE=1 zip -r -q "$(basename "$ZIP")" "$(basename "$STAGING")"
)

echo ""
echo "[5/5] 完成"
ls -lh "$ZIP"
du -sh "$STAGING" "$STAGING/web/dist" "$STAGING/oceanbase" "$STAGING/tdengine" 2>/dev/null || true
echo ""
echo "产出: $ZIP"
echo "应用: unzip ... && bash apply-patch.sh /opt/ISM/ism-release-oceanbase-20260709"
