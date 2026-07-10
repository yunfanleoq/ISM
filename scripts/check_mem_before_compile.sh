#!/usr/bin/env bash
# ISM 前端 dev server 编译前内存检查（macOS）
# 阈值：可用+可释放 ≥ 12GB；Node heap 需预留 20GB
# 用法：./scripts/check_mem_before_compile.sh
# 退出码：0=PASS，1=FAIL

set -euo pipefail

REQUIRED_GB=12
NODE_HEAP_GB=20
PASS=0
FAIL=1

if [[ "$(uname -s)" != "Darwin" ]]; then
  echo "FAIL: 本脚本仅支持 macOS（Darwin）。"
  exit "$FAIL"
fi

# --- 硬件总内存 ---
PAGE_SIZE=$(sysctl -n hw.pagesize 2>/dev/null || echo 4096)
TOTAL_BYTES=$(sysctl -n hw.memsize 2>/dev/null || echo 0)
TOTAL_GB=$(awk -v b="$TOTAL_BYTES" 'BEGIN { printf "%.2f", b / 1073741824 }')

# --- vm_stat 解析（页数 → 字节）---
vm_stat_out=$(vm_stat)

page_count() {
  local key="$1"
  echo "$vm_stat_out" | awk -v k="$key" '
    $0 ~ k {
      gsub(/\./, "", $NF)
      print $NF
      exit
    }
  '
}

pages_free=$(page_count "Pages free")
pages_inactive=$(page_count "Pages inactive")
pages_speculative=$(page_count "Pages speculative")
pages_purgeable=$(page_count "Pages purgeable")

pages_free=${pages_free:-0}
pages_inactive=${pages_inactive:-0}
pages_speculative=${pages_speculative:-0}
pages_purgeable=${pages_purgeable:-0}

bytes_free=$((pages_free * PAGE_SIZE))
bytes_releasable=$(((pages_inactive + pages_speculative + pages_purgeable) * PAGE_SIZE))
bytes_available=$((bytes_free + bytes_releasable))

gb_free=$(awk -v b="$bytes_free" 'BEGIN { printf "%.2f", b / 1073741824 }')
gb_releasable=$(awk -v b="$bytes_releasable" 'BEGIN { printf "%.2f", b / 1073741824 }')
gb_available=$(awk -v b="$bytes_available" 'BEGIN { printf "%.2f", b / 1073741824 }')

# --- memory_pressure（原样输出供人工判断）---
pressure_out=""
if command -v memory_pressure >/dev/null 2>&1; then
  pressure_out=$(memory_pressure 2>&1 || true)
fi

# --- Top 内存进程（RSS，单位 KB）---
# head 提前关闭管道会触发 SIGPIPE；临时关闭 pipefail
set +o pipefail
top_procs=$(
  ps -Aem -o rss=,pid=,comm= 2>/dev/null \
    | awk '{rss=$1; pid=$2; $1=""; $2=""; sub(/^  */, ""); printf "%10d KB  PID %-7s  %s\n", rss, pid, $0}' \
    | sort -rn \
    | head -12 \
    || true
)
set -o pipefail

echo "========================================"
echo " ISM 编译前内存检查"
echo "========================================"
echo ""
echo "[硬件]"
echo "  sysctl hw.memsize     : ${TOTAL_GB} GB（物理内存总量）"
echo "  hw.pagesize           : ${PAGE_SIZE} bytes"
echo ""
echo "[vm_stat 估算]"
echo "  可用（Pages free）    : ${gb_free} GB"
echo "  可释放（inactive+spec+purgeable）: ${gb_releasable} GB"
echo "  合计（可用+可释放）   : ${gb_available} GB"
echo "  编译门槛              : ≥ ${REQUIRED_GB} GB"
echo "  Node heap 预留        : ${NODE_HEAP_GB} GB（--max-old-space-size=20480）"
echo ""

if [[ -n "$pressure_out" ]]; then
  echo "[memory_pressure]"
  echo "$pressure_out" | sed 's/^/  /'
  echo ""
fi

echo "[Top 内存进程（RSS）]"
if [[ -n "$top_procs" ]]; then
  echo "$top_procs" | sed 's/^/  /'
else
  echo "  （无法获取进程列表）"
fi
echo ""

# --- 判定 ---
meets_threshold=$(awk -v avail="$gb_available" -v req="$REQUIRED_GB" 'BEGIN { print (avail >= req) ? "yes" : "no" }')

if [[ "$meets_threshold" == "yes" ]]; then
  echo "RESULT: PASS"
  echo ""
  echo "建议：内存充足，可继续执行编译前清场后启动 dev server。"
  echo "  NODE_OPTIONS=\"--max-old-space-size=20480 --openssl-legacy-provider\" \\"
  echo "    npx vue-cli-service serve --port 7080"
  exit "$PASS"
fi

echo "RESULT: FAIL"
echo ""
echo "可用+可释放仅 ${gb_available} GB，低于 ${REQUIRED_GB} GB 门槛。"
echo "在释放内存之前，禁止启动 vue-cli-service serve。"
echo ""
echo "建议按顺序执行："
echo "  1. 杀旧编译进程"
echo "     launchctl remove com.ism.frontend 2>/dev/null"
echo "     pkill -9 -f \"vue-cli-service\" 2>/dev/null"
echo "     lsof -ti :7080 | grep -v Cursor | xargs kill -9 2>/dev/null"
echo "     sleep 3"
echo "  2. 停 Docker / Colima"
echo "     colima stop 2>/dev/null || true"
echo "     docker stop \$(docker ps -q) 2>/dev/null || true"
echo "  3. 手动关闭高内存应用（Chrome、IDE 多余窗口、模拟器等）"
echo "     参考上方 Top 内存进程列表"
echo "  4. 可选：sudo purge（强制释放 inactive 内存，需管理员密码）"
echo "     sudo purge"
echo "  5. 重新运行本脚本："
echo "     ./scripts/check_mem_before_compile.sh"
echo ""
echo "若仍 FAIL，请勿编译，告知用户当前机器内存不足（需 ≥${REQUIRED_GB}GB 可用+可释放 + ${NODE_HEAP_GB}GB Node heap）。"
exit "$FAIL"
