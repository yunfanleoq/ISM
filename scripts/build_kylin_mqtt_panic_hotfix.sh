#!/bin/bash
# 热修复：MQTT panic + HTTPS 缺证书拖垮进程
# 产出小包（后端 + 完整 conf），可快速下发现场
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
# shellcheck source=lib_build_id.sh
source "$ROOT/scripts/lib_build_id.sh"
BUILD_ID="$(ism_build_id)"
PKG="ism-patch-kylin-runtime-fix-${BUILD_ID}"
STAGING="$ROOT/releases/${PKG}"
ZIP="$ROOT/releases/${PKG}.zip"
BIN_SRC="$ROOT/patches/ism-server-kylin-glibc228/ism_server"
REL09="$ROOT/releases/ism-release-oceanbase-20260709"

echo "=== 运行时崩溃热修复包 (${BUILD_ID}) ==="
if [[ ! -f "$BIN_SRC" ]]; then
  bash "$ROOT/scripts/build_kylin_ism_server.sh"
fi

rm -rf "$STAGING"
mkdir -p "$STAGING/ism_server_user/conf" "$STAGING/ism_server_user/sys_script" "$STAGING/scripts"
cp "$BIN_SRC" "$STAGING/ism_server_user/ism_server"
chmod 755 "$STAGING/ism_server_user/ism_server"

if [[ -d "$REL09/ism_server_user/conf" ]]; then
  rsync -a "$REL09/ism_server_user/conf/" "$STAGING/ism_server_user/conf/"
fi
for f in mqtt_broken_config.json videoConfig.json historyData.conf mqtt.conf \
         app.conf opcuaserver.conf 192.168.199.120.crt 192.168.199.120.key; do
  [[ -f "$ROOT/ism_server_user/conf/$f" ]] && cp "$ROOT/ism_server_user/conf/$f" "$STAGING/ism_server_user/conf/"
done
if [[ -d "$ROOT/ism_server_user/sys_script" ]]; then
  rsync -a "$ROOT/ism_server_user/sys_script/" "$STAGING/ism_server_user/sys_script/" || true
fi
touch "$STAGING/ism_server_user/sys_script/.gitkeep"

CONF="$STAGING/ism_server_user/conf/app.conf"
for kv in enablemqttbreoken=false enablehttps=false dbtype=4; do
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
HC="$STAGING/ism_server_user/conf/historyData.conf"
if [[ "$(uname -s)" == "Darwin" ]]; then
  sed -i '' 's/^historyrecorddbtype=.*/historyrecorddbtype=2/' "$HC"
  sed -i '' 's/^tdenginehost=.*/tdenginehost=127.0.0.1/' "$HC"
  sed -i '' 's/^tdengineport=.*/tdengineport=6041/' "$HC"
else
  sed -i 's/^historyrecorddbtype=.*/historyrecorddbtype=2/' "$HC"
  sed -i 's/^tdenginehost=.*/tdenginehost=127.0.0.1/' "$HC"
  sed -i 's/^tdengineport=.*/tdengineport=6041/' "$HC"
fi

cat > "$STAGING/apply-hotfix.sh" << 'APPLY'
#!/bin/bash
# 用法: bash apply-hotfix.sh /opt/ISM/ism-release-oceanbase-20260709
set -euo pipefail
PATCH_ROOT="$(cd "$(dirname "$0")" && pwd)"
TARGET="${1:?用法: bash apply-hotfix.sh <主包目录>}"
TARGET="$(cd "$TARGET" && pwd)"
[[ -d "$TARGET/ism_server_user" ]] || { echo "错误: 无效目录"; exit 1; }

echo "=== [1/4] 停止后端 ==="
(cd "$TARGET" && bash stop-all.sh 2>/dev/null) || true
pkill -9 -f '[.]/ism_server' 2>/dev/null || true
sleep 2

echo "=== [2/4] 替换 ism_server ==="
cp "$PATCH_ROOT/ism_server_user/ism_server" "$TARGET/ism_server_user/ism_server"
chmod 755 "$TARGET/ism_server_user/ism_server"

echo "=== [3/4] 补齐 conf/证书 + 关闭 MQTT Broker/HTTPS ==="
mkdir -p "$TARGET/ism_server_user/conf" "$TARGET/ism_server_user/sys_script"
for f in mqtt_broken_config.json videoConfig.json opcuaserver.conf \
         192.168.199.120.crt 192.168.199.120.key historyData.conf; do
  if [[ -f "$PATCH_ROOT/ism_server_user/conf/$f" ]]; then
    cp "$PATCH_ROOT/ism_server_user/conf/$f" "$TARGET/ism_server_user/conf/$f"
    echo "  + conf/$f"
  fi
done
if [[ -d "$PATCH_ROOT/ism_server_user/sys_script" ]]; then
  rsync -a "$PATCH_ROOT/ism_server_user/sys_script/" "$TARGET/ism_server_user/sys_script/" || true
fi

CONF="$TARGET/ism_server_user/conf/app.conf"
if [[ -f "$CONF" ]]; then
  for kv in enablemqttbreoken=false enablehttps=false; do
    key="${kv%%=*}"
    if grep -q "^${key}=" "$CONF"; then
      sed -i "s/^${key}=.*/${kv}/" "$CONF"
    else
      echo "$kv" >> "$CONF"
    fi
  done
else
  cp "$PATCH_ROOT/ism_server_user/conf/app.conf" "$CONF"
fi

echo "=== [4/4] 启动 ==="
(cd "$TARGET" && bash start-all.sh)
sleep 8
echo "--- 进程 ---"
pgrep -fl ism_server || echo "警告: ism_server 未运行"
echo "--- 日志尾 ---"
tail -40 "$TARGET/logs/ism_server.log" 2>/dev/null || tail -40 "$TARGET/ism_server_user/logs/ism.log" 2>/dev/null || true
echo ""
echo "期望:"
echo "  - 无 panic"
echo "  - 有「禁止启动 MQTT Broken」"
echo "  - 有 http server Running on http://:8091"
echo "  - 无 ListenAndServeTLS / open conf/*.crt 错误"
echo "  - TDengine 连接成功"
APPLY
chmod +x "$STAGING/apply-hotfix.sh"

cat > "$STAGING/README.md" << README
# 运行时崩溃热修复（MQTT + HTTPS）

构建: ${BUILD_ID}

## 现场现象
1. MQTT Broker 缺配置 → panic
2. \`enablehttps=true\` 但缺 \`conf/*.crt\` → \`ListenAndServeTLS\` 失败，**整个进程退出**

## 本包修复
- 后端 MQTT 缺配置安全返回
- 默认 \`enablemqttbreoken=false\`、\`enablehttps=false\`
- 强制补齐证书 / mqtt / video / opcua / historyData / sys_script

## 应用
\`\`\`bash
unzip ${PKG}.zip
cd ${PKG}
bash apply-hotfix.sh /opt/ISM/ism-release-oceanbase-20260709
\`\`\`

## 应急（不换二进制）
\`\`\`bash
cd /opt/ISM/ism-release-oceanbase-20260709/ism_server_user
sed -i 's/^enablemqttbreoken=.*/enablemqttbreoken=false/' conf/app.conf
sed -i 's/^enablehttps=.*/enablehttps=false/' conf/app.conf
# 或拷入 192.168.199.120.crt / .key 后再开 HTTPS
cd .. && bash stop-all.sh && bash start-all.sh
\`\`\`
README

rm -f "$ZIP"
(
  cd "$ROOT/releases"
  COPYFILE_DISABLE=1 zip -r -q "$(basename "$ZIP")" "$(basename "$STAGING")"
)
ls -lh "$ZIP"
echo "产出: $ZIP"
