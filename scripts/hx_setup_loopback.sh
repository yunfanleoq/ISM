#!/bin/bash
# 中航信数据中心 —— 为所有真实管理机 IP 创建本机回环别名，
# 使后端可连到 <真实IP>:502（由 hx_simulator.py 真实IP模式监听）。
# 需要 sudo（回环别名 + 绑定 502 端口均为特权操作）。
#
# 用法:
#   sudo bash scripts/hx_setup_loopback.sh validation   # 仅验证集 3 个 IP
#   sudo bash scripts/hx_setup_loopback.sh full         # 全量集所有 IP
#   sudo bash scripts/hx_setup_loopback.sh down         # 移除所有别名
set -e
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SET="${1:-validation}"

ips_for() {
  python3 - "$1" <<'PY'
import json,sys,os
root=os.path.dirname(os.path.dirname(os.path.abspath(sys.argv[0]))) if False else os.getcwd()
base=os.path.join("hx-data", sys.argv[1])
dm=json.load(open(os.path.join(base,"ism_data_models.json")))
ips=sorted({d.get("gatewayIP") for d in dm["devices"] if d.get("gatewayIP")})
print("\n".join(ips))
PY
}

cd "$ROOT"

if [ "$SET" = "down" ]; then
  for S in validation full; do
    [ -d "hx-data/$S" ] || continue
    for ip in $(ips_for "$S"); do
      [ "$ip" = "127.0.0.1" ] && continue
      echo "remove alias $ip"
      ifconfig lo0 -alias "$ip" 2>/dev/null || true
    done
  done
  exit 0
fi

for ip in $(ips_for "$SET"); do
  [ "$ip" = "127.0.0.1" ] && continue
  echo "add alias $ip"
  ifconfig lo0 alias "$ip" 255.255.255.255
done
echo "完成。现在可运行:  sudo python3 scripts/hx_simulator.py --set $SET"
