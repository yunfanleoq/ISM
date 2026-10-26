#!/bin/bash
# 适用：银河麒麟 V10 SP3 x86_64 + OceanBase 一体离线包
# 前置：整包已解压；sudo/root；内存≥8G；完全离线；勿改 /opt/ISMCode
# 用法：sudo bash deploy-min.sh   （整包根目录）
#       或：sudo bash scripts/deploy_kylin_oceanbase_min.sh  （仓库内，自动选最新 releases/ism-release-oceanbase-*）
set -euo pipefail

HERE="$(cd "$(dirname "$0")" && pwd)"
if [[ -f "$HERE/deploy-offline.sh" ]]; then
  ROOT="$HERE"
else
  REPO="$(cd "$HERE/.." && pwd)"
  ROOT="$(ls -d "$REPO"/releases/ism-release-oceanbase-* 2>/dev/null | sort | tail -1 || true)"
  [[ -n "${ROOT:-}" && -f "$ROOT/deploy-offline.sh" ]] || {
    echo "错误: 未找到整包（需含 deploy-offline.sh）。请在解压后的整包根目录执行本脚本。"
    exit 1
  }
fi

cd "$ROOT"
chmod +x deploy-offline.sh start-all.sh stop-all.sh ism_server_user/ism_server scripts/*.sh 2>/dev/null || true

# 复用包内离线部署：装 Docker（若无）→ start-all（OB+TDengine+前后端）
bash ./deploy-offline.sh
