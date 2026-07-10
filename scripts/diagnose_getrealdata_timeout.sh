#!/bin/bash
# 诊断：数据仓库 getRealData 120s 超时
#
# 用法（在麒麟一体包根目录执行）:
#   cd /opt/ISM/ism-release-oceanbase-20260709   # 或你的实际目录
#   bash scripts/diagnose_getrealdata_timeout.sh
#
# 指定控制台里超时的设备 UUID（推荐）:
#   DEVICE_UUID=05a70c46-xxxx-xxxx-xxxx-xxxxxxxxxxxx \
#     bash scripts/diagnose_getrealdata_timeout.sh
#
# 登录若因 OceanBase GORM/`user` 表问题返回 1003，可跳过登录直接测接口:
#   TOKEN='浏览器 Application→Local Storage 里的 Authorization' \
#   PROJECT_UUID='3ec5821f-...' \
#   DEVICE_UUID='343e5e94-...' \
#     bash scripts/diagnose_getrealdata_timeout.sh
#
# 可选环境变量:
#   TOKEN / AUTH_TOKEN / PROJECT_UUID
#   ISM_BE_PORT / ISM_FE_PORT / OB_PORT / OB_TENANT / OB_PASSWORD / OB_USER / OB_DB
set -u

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
mkdir -p "$ROOT/logs"

if [[ -f "$ROOT/ports.env" ]]; then
  # shellcheck disable=SC1091
  source "$ROOT/ports.env"
fi

BE_PORT="${ISM_BE_PORT:-8091}"
FE_PORT="${ISM_FE_PORT:-7090}"
OB_PORT="${OB_PORT:-2881}"
OB_TENANT="${OB_TENANT:-ism_tenant}"
OB_PASSWORD="${OB_PASSWORD:-ism2024!}"
OB_USER="${OB_USER:-root@${OB_TENANT}}"
OB_DB="${OB_DB:-ism}"
DEVICE_UUID="${DEVICE_UUID:-}"
PROJECT_UUID="${PROJECT_UUID:-}"
# 兼容 TOKEN / AUTH_TOKEN；去掉可能的 Bearer 前缀空格
TOKEN="${TOKEN:-${AUTH_TOKEN:-}}"
TOKEN="${TOKEN#Bearer }"
TOKEN="${TOKEN#bearer }"
TOKEN="$(echo -n "$TOKEN" | tr -d '\r\n')"

TS="$(date '+%Y%m%d_%H%M%S')"
LOG="$ROOT/logs/ism_getrealdata_diag_${TS}.log"
exec > >(tee -a "$LOG") 2>&1

section() { echo ""; echo "########## $* ##########"; echo ""; }
ok() { echo "[OK] $*"; }
bad() { echo "[FAIL] $*"; }
info() { echo "[INFO] $*"; }

ms_now() {
  # 优先毫秒；无 %3N 时退化为秒*1000
  local t
  t="$(date +%s%3N 2>/dev/null || true)"
  if [[ "$t" =~ ^[0-9]+$ ]]; then
    echo "$t"
  else
    echo $(($(date +%s) * 1000))
  fi
}

ob_sql() {
  local sql="$1"
  if docker ps --format '{{.Names}}' 2>/dev/null | grep -qx oceanbase; then
    docker exec oceanbase obclient --default-character-set=utf8mb4 \
      -h127.0.0.1 -P"${OB_PORT}" -u"${OB_USER}" -p"${OB_PASSWORD}" "${OB_DB}" -N -e "${sql}" 2>/dev/null
  elif command -v obclient >/dev/null 2>&1; then
    obclient --default-character-set=utf8mb4 \
      -h127.0.0.1 -P"${OB_PORT}" -u"${OB_USER}" -p"${OB_PASSWORD}" "${OB_DB}" -N -e "${sql}" 2>/dev/null
  else
    return 1
  fi
}

echo "=============================================="
echo " ISM getRealData 超时诊断"
echo " 时间: $(date '+%F %T')"
echo " 根目录: $ROOT"
echo " 后端端口: $BE_PORT  前端端口: $FE_PORT"
echo " 日志: $LOG"
echo "=============================================="

# ---------- 1) 进程 / 端口 ----------
section "1) 进程与端口"
if pgrep -fl ism_server >/dev/null 2>&1; then
  pgrep -fl ism_server | head -5
  ok "ism_server 进程存在"
else
  bad "未找到 ism_server 进程"
fi
if lsof -nP -iTCP:"${BE_PORT}" -sTCP:LISTEN >/dev/null 2>&1; then
  lsof -nP -iTCP:"${BE_PORT}" -sTCP:LISTEN | head -3
  ok "后端 ${BE_PORT} 在监听"
else
  bad "后端 ${BE_PORT} 未监听"
fi

# ---------- 2) 二进制是否含分页 ----------
section "2) 后端是否含 GetRealDataPaged（分页修复）"
BIN=""
for cand in \
  "$ROOT/ism_server_user/ism_server" \
  "$ROOT/ism_server" \
  "$(command -v ism_server 2>/dev/null || true)"; do
  if [[ -n "$cand" && -f "$cand" ]]; then
    BIN="$cand"
    break
  fi
done

if [[ -z "$BIN" ]]; then
  bad "未找到 ism_server 二进制"
else
  info "二进制: $BIN"
  info "大小: $(ls -lh "$BIN" | awk '{print $5}')  mtime: $(ls -l "$BIN" | awk '{print $6,$7,$8}')"
  PAGED_CNT="$(strings "$BIN" 2>/dev/null | grep -c 'GetRealDataPaged' || true)"
  if [[ "${PAGED_CNT}" -ge 1 ]]; then
    ok "含 GetRealDataPaged（count=${PAGED_CNT}）— 已打过分页补丁"
  else
    bad "不含 GetRealDataPaged — 仍是旧后端，会一次拉全量测点 → 极易 120s 超时"
    info "处理: 应用 ism-patch-kylin-full-20260709 后重启"
  fi
fi

# ---------- 3) 登录（可被 TOKEN 跳过；1003=OceanBase GORM/`user` 老问题） ----------
section "3) 登录"
LOGIN_OK=0
if [[ -n "$TOKEN" ]]; then
  ok "使用外部 TOKEN 跳过 /login（长度=${#TOKEN}）— 绕过 GORM 1003"
  LOGIN_OK=1
  CODE="skipped"
else
  LOGIN="$(curl -s --compressed -m 20 -X POST "http://127.0.0.1:${BE_PORT}/login" \
    -H 'Content-Type: application/json' \
    -d '{"Username":"admin","password":"e10adc3949ba59abbe56e057f20f883e"}' || true)"
  TOKEN="$(echo "$LOGIN" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('data',{}).get('token',''))" 2>/dev/null || true)"
  CODE="$(echo "$LOGIN" | python3 -c "import sys,json; print(json.load(sys.stdin).get('code',''))" 2>/dev/null || true)"
  if [[ -z "$TOKEN" ]]; then
    bad "登录失败 code=${CODE}（1003=GORM 查不到 user，不一定是没 admin / 1002=密码错误）"
    echo "$LOGIN" | head -c 400
    echo ""
    section "3b) user 表排查"
    echo "--- SELECT FROM \`user\` ---"
    ob_sql "SELECT id, username, HEX(username) AS username_hex, role, deleted_at FROM \`user\` ORDER BY id LIMIT 20;" \
      || bad "无法查 user 表"
    echo "--- GORM 等价条件 ---"
    ob_sql "SELECT id, username, role, deleted_at FROM \`user\` WHERE username='admin' AND deleted_at IS NULL LIMIT 1;" \
      || true
    USER_CNT="$(ob_sql "SELECT COUNT(*) FROM \`user\`;" | tr -d '\r' | head -1 || true)"
    info "user 表行数: ${USER_CNT:-?}"
    if [[ "${CODE}" == "1003" ]]; then
      info "code=1003 已知根因：OceanBase 下 GORM 查保留表 user 匹配失败（不是真没用户）"
      echo "  修复: 后端需含 userTable()/lookupUserByUsername 并重启"
      echo "  绕过: TOKEN=<浏览器 Authorization> PROJECT_UUID=<项目uuid> 再跑本脚本"
    fi
    info "无 TOKEN：跳过 HTTP getRealData，继续测点数/索引/SQL 直查"
  else
    ok "登录成功 code=${CODE}"
    LOGIN_OK=1
  fi
fi

# ---------- 3c) 后端是否含登录修复 ----------
section "3c) 后端是否含 userTable（登录 1003 / GORM 修复）"
if [[ -n "$BIN" ]]; then
  UT_CNT="$(strings "$BIN" 2>/dev/null | grep -c 'userTable\|lookupUserByUsername' || true)"
  if [[ "${UT_CNT}" -ge 1 ]]; then
    ok "含 userTable/lookupUserByUsername（count=${UT_CNT}）— 若 curl 仍 1003，确认跑的是这个二进制并已重启"
  else
    bad "不含 userTable — 登录 1003 就是旧登录逻辑；换含登录修复的 ism_server 后重启"
  fi
fi

# ---------- 4) 项目 / 设备 UUID ----------
section "4) 项目与目标设备"
if [[ -z "$PROJECT_UUID" ]]; then
  PROJECT_UUID="$(ob_sql "SELECT uuid FROM project_lists LIMIT 1;" | tr -d '\r' | head -1 || true)"
fi
if [[ -z "$PROJECT_UUID" && "$LOGIN_OK" -eq 1 ]]; then
  PROJECT_UUID="$(curl -s --compressed -m 20 -X POST "http://127.0.0.1:${BE_PORT}/ProjectList" \
    -H "Authorization: ${TOKEN}" -H 'Content-Type: application/json' -d '{}' \
    | python3 -c "import sys,json; d=json.load(sys.stdin); print((d.get('list') or [{}])[0].get('ProjectInfo',{}).get('uuid',''))" 2>/dev/null || true)"
fi
info "project_uuid=${PROJECT_UUID:-<空>}"
if [[ -z "$PROJECT_UUID" ]]; then
  bad "未取到 ProjectUuid — HTTP getRealData 会 code=-6；DB 直查仍可继续"
fi

if [[ -z "$DEVICE_UUID" ]]; then
  DEVICE_UUID="$(ob_sql "SELECT uuid FROM monitor_list WHERE type=1 LIMIT 1;" | tr -d '\r' | head -1 || true)"
  info "未指定 DEVICE_UUID，使用首个设备: ${DEVICE_UUID:-<空>}"
else
  info "使用指定 DEVICE_UUID=${DEVICE_UUID}"
fi
if [[ -z "$DEVICE_UUID" ]]; then
  bad "无可用设备 UUID"
  echo "日志已写: $LOG"
  exit 1
fi

DEV_NAME="$(ob_sql "SELECT name FROM monitor_list WHERE uuid='${DEVICE_UUID}' LIMIT 1;" | tr -d '\r' | head -1 || true)"
info "设备名: ${DEV_NAME:-<未知>}"

# ---------- 5) 测点数 + 索引 ----------
section "5) device_real_data 测点数与索引"
TOTAL_ALL="$(ob_sql "SELECT COUNT(*) FROM device_real_data;" | tr -d '\r' | head -1 || true)"
TOTAL_DEV="$(ob_sql "SELECT COUNT(*) FROM device_real_data WHERE device_uuid='${DEVICE_UUID}';" | tr -d '\r' | head -1 || true)"
info "全表测点: ${TOTAL_ALL:-?}"
info "本设备测点: ${TOTAL_DEV:-?}"

if [[ -n "${TOTAL_DEV}" && "${TOTAL_DEV}" =~ ^[0-9]+$ ]]; then
  if [[ "$TOTAL_DEV" -gt 50000 ]]; then
    bad "本设备测点 ${TOTAL_DEV} > 5万 — 旧后端全量拉取必超时；新后端分页也应很快，若仍超时查锁/索引"
  elif [[ "$TOTAL_DEV" -gt 5000 ]]; then
    info "本设备测点 ${TOTAL_DEV} 偏多，分页首屏应仍 <2s；若满 120s 优先查后端版本/DB 锁"
  else
    ok "本设备测点 ${TOTAL_DEV} 正常规模"
  fi
fi

echo "--- SHOW INDEX FROM device_real_data (device_uuid 相关) ---"
ob_sql "SHOW INDEX FROM device_real_data;" | grep -i 'device_uuid\|Key_name\|Column' || \
  ob_sql "SHOW INDEX FROM device_real_data;" || \
  bad "无法查询索引（obclient/docker 不可用）"

IDX_HIT="$(ob_sql "SHOW INDEX FROM device_real_data;" 2>/dev/null | grep -ci 'device_uuid' || true)"
if [[ "${IDX_HIT}" -ge 1 ]]; then
  ok "存在 device_uuid 相关索引"
else
  bad "未看到 device_uuid 索引 — 建议: CREATE INDEX idx_device_real_data_device_uuid ON device_real_data(device_uuid);"
fi

# ---------- 6) SQL 直查耗时（绕过 HTTP） ----------
section "6) OceanBase 直查耗时（绕过 HTTP）"
START="$(ms_now)"
ob_sql "SELECT COUNT(*) FROM device_real_data WHERE device_uuid='${DEVICE_UUID}';" >/dev/null
END="$(ms_now)"
info "COUNT(*) 耗时 ms: $((END-START))"

START="$(ms_now)"
ob_sql "SELECT id,name,value,uuid FROM device_real_data WHERE device_uuid='${DEVICE_UUID}' ORDER BY id ASC LIMIT 50;" >/dev/null
END="$(ms_now)"
info "LIMIT 50 耗时 ms: $((END-START))"

START="$(ms_now)"
ROWCNT="$(ob_sql "SELECT COUNT(*) FROM (SELECT id FROM device_real_data WHERE device_uuid='${DEVICE_UUID}' LIMIT 5000) t;" | tr -d '\r' | head -1 || true)"
END="$(ms_now)"
info "抽样至多 5000 行耗时 ms: $((END-START))  rows=${ROWCNT:-?}"

# ---------- 7) HTTP getRealData（分页 vs 无分页） ----------
if [[ "$LOGIN_OK" -ne 1 || -z "$PROJECT_UUID" ]]; then
  section "7) HTTP /getRealData — 跳过"
  bad "登录失败或无 ProjectUuid，无法测 HTTP；请先修复登录后再跑本脚本"
  section "8) HTTP 无分页 — 跳过"
else
  section "7) HTTP /getRealData 分页请求（page=1,pageSize=50）"
  START="$(ms_now)"
  REAL_PAGED="$(curl -s --compressed -m 130 -X POST "http://127.0.0.1:${BE_PORT}/getRealData" \
    -H "Authorization: ${TOKEN}" \
    -H "ProjectUuid: ${PROJECT_UUID}" \
    -H 'Content-Type: application/json' \
    -d "{\"uuid\":\"${DEVICE_UUID}\",\"page\":1,\"pageSize\":50,\"IsRemoveGW\":false}" || true)"
  END="$(ms_now)"
  PAGED_MS=$((END-START))
  info "耗时 ms: ${PAGED_MS}"
  echo "$REAL_PAGED" | python3 -c "
import sys,json
raw=sys.stdin.read()
try:
  d=json.loads(raw)
except Exception as e:
  print('解析失败:', e, '前200字节:', raw[:200])
  sys.exit(0)
rd=d.get('realData') or []
print('code=', d.get('code'), 'rows=', len(rd), 'total=', d.get('total'), 'page=', d.get('page'), 'pageSize=', d.get('pageSize'), 'hasMore=', d.get('hasMore'))
" 2>/dev/null || echo "$REAL_PAGED" | head -c 300

  if [[ "$PAGED_MS" -ge 110000 ]]; then
    bad "分页请求接近/超过 120s — 与前端 timeout of 120000ms exceeded 一致"
  elif [[ "$PAGED_MS" -ge 5000 ]]; then
    bad "分页请求偏慢 (${PAGED_MS}ms)，需查 OB 锁/慢查询"
  else
    ok "分页请求正常 (${PAGED_MS}ms)"
  fi

  section "8) HTTP /getRealData 无分页（模拟旧前端，危险：可能很慢）"
  info "仅测 25s 上限；若旧后端全量拉取，这里会卡住直到超时"
  START="$(ms_now)"
  REAL_FULL="$(curl -s --compressed -m 25 -X POST "http://127.0.0.1:${BE_PORT}/getRealData" \
    -H "Authorization: ${TOKEN}" \
    -H "ProjectUuid: ${PROJECT_UUID}" \
    -H 'Content-Type: application/json' \
    -d "{\"uuid\":\"${DEVICE_UUID}\",\"IsRemoveGW\":false}" || true)"
  END="$(ms_now)"
  FULL_MS=$((END-START))
  info "耗时 ms: ${FULL_MS}（curl -m 25，超时则≈25000）"
  if [[ -z "$REAL_FULL" ]]; then
    bad "无分页请求在 25s 内无响应 — 强烈暗示后端仍在全量拉测点或 DB 卡住"
  else
    echo "$REAL_FULL" | python3 -c "
import sys,json
raw=sys.stdin.read()
try:
  d=json.loads(raw)
  rd=d.get('realData') or []
  print('code=', d.get('code'), 'rows=', len(rd), 'has_page_fields=', 'page' in d)
except Exception as e:
  print('解析失败/非JSON:', e, '前200字节:', raw[:200])
" 2>/dev/null || echo "$REAL_FULL" | head -c 200
  fi
fi

# ---------- 9) 结论 ----------
section "9) 结论速判"
echo "对照表:"
echo "  0. 步骤3 登录 code=1003               → 先修 user/admin（check_login_and_user.sh / 重置密码），否则浏览器也会登不上或 token 异常"
echo "  A. 步骤2 无 GetRealDataPaged          → 打 ism-patch-kylin-full-20260709，替换 ism_server 并重启"
echo "  B. 步骤2 有分页，但步骤7 ≥110s         → 查 OB 锁/连接池；对该设备 COUNT 与 SHOW PROCESSLIST"
echo "  C. 步骤5 无 device_uuid 索引           → CREATE INDEX idx_device_real_data_device_uuid ON device_real_data(device_uuid);"
echo "  D. 步骤6 SQL 很快、步骤7 HTTP 很慢     → 卡在后端 gzip/编码或连接池排队，看 ism_server 日志"
echo "  E. 步骤7 很快(<2s) 但浏览器仍超时      → 浏览器打的不是本机后端，或前端未 Ctrl+F5 / 代理到旧实例"
echo ""
echo "浏览器复现时请在 Network 看 /api/getRealData 的 Waiting 时间；"
echo "控制台 DW.getRealData.catch {message: timeout of 120000ms exceeded} = 本脚本步骤7 同类问题。"
echo ""
echo "完成。请下载/回传日志: $LOG"
