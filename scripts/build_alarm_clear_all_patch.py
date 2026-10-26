#!/usr/bin/env python3
"""
构建「实时告警一键清除」现场补丁包 patches/alarm-clear-all-v1.zip

用法:
  python3 scripts/build_alarm_clear_all_patch.py
  python3 scripts/build_alarm_clear_all_patch.py --check-only
  python3 scripts/build_alarm_clear_all_patch.py --skip-frontend
  python3 scripts/build_alarm_clear_all_patch.py --target linux/amd64
"""

from __future__ import annotations

import argparse
import fcntl
import os
import re
import shutil
import subprocess
import sys
import zipfile
from datetime import datetime
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
PATCH_NAME = "alarm-clear-all-v1"
PATCH_DIR = ROOT / "patches" / PATCH_NAME
PATCH_ZIP = ROOT / "patches" / f"{PATCH_NAME}.zip"
SERVER_DIR = ROOT / "ism_server_user"
FRONTEND_DIR = ROOT / "ism-front-end-v2"
CLEAR_SCRIPT = ROOT / "scripts" / "clear_all_alarms.py"

# 源码 marker 检查：文件 -> 必须包含的子串（任一命中即可时用列表）
MARKERS: dict[str, list[str]] = {
    "ism_server_user/models/alarmModel.go": ["AlarmClearAll", "ResyncOfflineDeviceAlarms"],
    "ism_server_user/controllers/alarmCtl.go": ["AlarmClearAll"],
    "ism_server_user/routers/router.go": ["/AlarmClearAll"],
    "ism-front-end-v2/src/pages/alarm/currentAlarm/currentAlarm.vue": ["ClearAllAlarm", "clearAllAlarmFallback", "clearAllBatch"],
    "ism-front-end-v2/src/services/alarm.js": ["ClearAllCurrentAlarm"],
}

INSTALL_SH = """#!/bin/bash
# ISM 实时告警一键清除补丁 - 安装脚本（兼容标准旧版 v0 升级）
# 用法: sudo bash install.sh
# 环境变量: ISM_BIN, ISM_DIST（可覆盖自动探测结果）

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TS="$(date +%Y%m%d%H%M%S)"

# --- 自动探测 ISM 二进制路径 ---
detect_ism_bin() {
  local candidates=()
  if [[ -n "${ISM_BIN:-}" ]]; then candidates+=("$ISM_BIN"); fi
  candidates+=(
    "/opt/ism/ism_server"
    "/opt/ism/ism_server_user/ism_server"
    "${SCRIPT_DIR}/ism_server"
    "./ism_server"
    "$(pwd)/ism_server"
  )
  # 开发仓库相对路径（从补丁目录向上找）
  candidates+=(
    "${SCRIPT_DIR}/../../ism_server_user/ism_server"
    "${SCRIPT_DIR}/../../../ism_server_user/ism_server"
  )
  for p in "${candidates[@]}"; do
    if [[ -f "$p" ]]; then
      echo "$p"
      return 0
    fi
  done
  return 1
}

# --- 自动探测前端 dist 路径 ---
detect_ism_dist() {
  local candidates=()
  if [[ -n "${ISM_DIST:-}" ]]; then candidates+=("$ISM_DIST"); fi
  candidates+=(
    "/opt/ism/web/dist"
    "/opt/ism/dist"
    "/var/www/ism/dist"
    "./dist"
    "$(pwd)/dist"
  )
  candidates+=(
    "${SCRIPT_DIR}/../../ism-front-end-v2/dist"
    "${SCRIPT_DIR}/../../../ism-front-end-v2/dist"
  )
  for p in "${candidates[@]}"; do
    if [[ -d "$p" ]]; then
      echo "$p"
      return 0
    fi
  done
  return 1
}

ISM_BIN="$(detect_ism_bin || true)"
ISM_DIST="$(detect_ism_dist || true)"

if [[ -z "$ISM_BIN" ]]; then
  echo "错误: 未找到 ism_server 二进制，请设置 ISM_BIN=路径"
  exit 1
fi
if [[ -z "$ISM_DIST" ]]; then
  echo "错误: 未找到前端 dist 目录，请设置 ISM_DIST=路径"
  exit 1
fi

echo "=== ISM 告警一键清除补丁安装 ==="
echo "目标二进制: ${ISM_BIN}"
echo "目标前端:   ${ISM_DIST}"
echo ""
echo "【重要】标准旧版无 AlarmClearAll API，必须同时替换后端 ism_server 与前端 dist。"
echo ""

# 安装前检查旧二进制
if [[ -f "$ISM_BIN" ]]; then
  if command -v file >/dev/null 2>&1; then
  if file "$ISM_BIN" | grep -q ELF; then
    echo "[检查] 旧二进制为 ELF 可执行文件"
  else
    echo "[警告] 旧二进制可能不是 ELF: $ISM_BIN"
  fi
  fi
  if [[ ! -x "$ISM_BIN" ]]; then
    echo "[警告] 旧二进制不可执行，安装后将替换为新文件"
  fi
else
  echo "[提示] 目标路径尚无 ism_server，将新建"
fi

# 1. 停服（常见 systemd 服务名，可按现场调整）
STOPPED_SVC=""
for svc in ism_server ism-server ism; do
  if systemctl is-active --quiet "$svc" 2>/dev/null; then
    echo "[1/4] 停止服务: $svc"
    systemctl stop "$svc"
    STOPPED_SVC="$svc"
    break
  fi
done
if [[ -z "${STOPPED_SVC}" ]]; then
  echo "[1/4] 未检测到 systemd 服务，尝试 pkill ism_server..."
  pkill -f "ism_server" 2>/dev/null || true
  sleep 2
fi

# 2. 备份
echo "[2/4] 备份现有文件..."
if [[ -f "$ISM_BIN" ]]; then
  cp -a "$ISM_BIN" "${ISM_BIN}.bak.${TS}"
  echo "  已备份: ${ISM_BIN}.bak.${TS}"
fi
if [[ -d "$ISM_DIST" ]]; then
  cp -a "$ISM_DIST" "${ISM_DIST}.bak.${TS}"
  echo "  已备份: ${ISM_DIST}.bak.${TS}"
fi

# 3. 替换
echo "[3/4] 部署补丁文件..."
install -m 755 "${SCRIPT_DIR}/ism_server" "$ISM_BIN"
if [[ -d "${SCRIPT_DIR}/dist" ]]; then
  rm -rf "$ISM_DIST"
  cp -a "${SCRIPT_DIR}/dist" "$ISM_DIST"
  echo "  前端 dist 已更新"
else
  echo "  错误: 补丁包内无 dist 目录 — 标准版升级必须包含前端，请使用完整补丁包"
  exit 1
fi
if [[ -f "${SCRIPT_DIR}/clear_all_alarms.py" ]]; then
  install -m 755 "${SCRIPT_DIR}/clear_all_alarms.py" /opt/ism/clear_all_alarms.py 2>/dev/null || \\
    cp -a "${SCRIPT_DIR}/clear_all_alarms.py" "${SCRIPT_DIR}/clear_all_alarms.py.deployed" 2>/dev/null || true
  echo "  应急脚本已部署"
fi
if [[ -f "${SCRIPT_DIR}/verify_old_standard.sh" ]]; then
  install -m 755 "${SCRIPT_DIR}/verify_old_standard.sh" /opt/ism/verify_old_standard.sh 2>/dev/null || true
fi

# 4. 重启
echo "[4/4] 重启服务..."
if [[ -n "${STOPPED_SVC}" ]]; then
  systemctl start "$STOPPED_SVC"
  echo "  已启动: $STOPPED_SVC"
else
  echo "  请手动启动 ism_server，例如:"
  echo "    cd $(dirname "$ISM_BIN") && nohup ./ism_server > /var/log/ism_server.log 2>&1 &"
fi

echo ""
echo "=== 安装完成 ==="
echo "验证: bash ${SCRIPT_DIR}/verify_old_standard.sh"
echo "登录 ISM -> 实时告警页面 -> 点击「一键清除」按钮验证。"
echo "回滚: sudo bash ${SCRIPT_DIR}/rollback.sh"
"""

ROLLBACK_SH = """#!/bin/bash
# ISM 实时告警一键清除补丁 - 回滚脚本
# 用法: sudo bash rollback.sh [备份时间戳 TS，如 20260627143000；省略则使用最新 .bak.*]

set -euo pipefail

ISM_BIN="${ISM_BIN:-/opt/ism/ism_server}"
ISM_DIST="${ISM_DIST:-/opt/ism/web/dist}"
TS="${1:-}"

pick_latest() {
  local base="$1"
  ls -dt "${base}.bak."* 2>/dev/null | head -1
}

if [[ -z "$TS" ]]; then
  BIN_BAK="$(pick_latest "$ISM_BIN")"
  DIST_BAK="$(pick_latest "$ISM_DIST")"
else
  BIN_BAK="${ISM_BIN}.bak.${TS}"
  DIST_BAK="${ISM_DIST}.bak.${TS}"
fi

echo "=== ISM 补丁回滚 ==="

for svc in ism_server ism-server ism; do
  if systemctl is-active --quiet "$svc" 2>/dev/null; then
    systemctl stop "$svc"
    STOPPED_SVC="$svc"
    break
  fi
done
pkill -f "ism_server" 2>/dev/null || true
sleep 2

if [[ -n "$BIN_BAK" && -f "$BIN_BAK" ]]; then
  cp -a "$BIN_BAK" "$ISM_BIN"
  echo "已恢复二进制: $BIN_BAK"
else
  echo "警告: 未找到二进制备份"
fi

if [[ -n "$DIST_BAK" && -d "$DIST_BAK" ]]; then
  rm -rf "$ISM_DIST"
  cp -a "$DIST_BAK" "$ISM_DIST"
  echo "已恢复前端: $DIST_BAK"
else
  echo "警告: 未找到 dist 备份"
fi

if [[ -n "${STOPPED_SVC:-}" ]]; then
  systemctl start "$STOPPED_SVC"
fi

echo "=== 回滚完成 ==="
"""

VERIFY_OLD_STANDARD_SH = """#!/bin/bash
# 检查现场 ISM 是否已具备 AlarmClearAll API（标准旧版无此接口）
# 用法: bash verify_old_standard.sh
# 环境变量: ISM_API_BASE (默认 http://127.0.0.1:8081)

set -euo pipefail

API_BASE="${ISM_API_BASE:-http://127.0.0.1:8081}"
URL="${API_BASE%/}/AlarmClearAll"

echo "=== ISM AlarmClearAll API 探测 ==="
echo "请求: POST $URL"

if ! command -v curl >/dev/null 2>&1; then
  echo "错误: 需要 curl"
  exit 1
fi

HTTP_CODE="$(curl -s -o /tmp/ism_alarm_clear_all_probe.json -w '%{http_code}' \\
  -X POST "$URL" \\
  -H 'Content-Type: application/json' \\
  -d '{}' \\
  --connect-timeout 5 \\
  --max-time 10 || echo '000')"

echo "HTTP 状态码: $HTTP_CODE"

if [[ "$HTTP_CODE" == "404" || "$HTTP_CODE" == "405" || "$HTTP_CODE" == "000" ]]; then
  echo ""
  echo "【结果】当前为标准旧版（无 AlarmClearAll）"
  echo "  - 必须执行完整补丁安装（同时替换 ism_server + dist）"
  echo "  - 或应急使用: python3 clear_all_alarms.py"
  echo "  - 仅换前端时 UI 将降级为逐条清除（≤500 条），大批量请用脚本"
  exit 2
fi

if [[ -f /tmp/ism_alarm_clear_all_probe.json ]]; then
  echo "响应体: $(head -c 200 /tmp/ism_alarm_clear_all_probe.json)"
fi

echo ""
echo "【结果】后端已具备 AlarmClearAll API（补丁后端已安装或新版本）"
echo "  可正常使用前端「一键清除」"
exit 0
"""

INSTALL_README = """ISM 实时告警「一键清除」补丁 - 现场安装说明
================================================

【功能】
  在「实时告警」页面增加「一键清除」按钮，批量清除当前项目未消除告警。

【标准旧版 vs 补丁版】
  标准旧版已有:
    - POST /GetCurrentAlarmList  查询实时告警
    - POST /AlarmOpt             单条清除 (type=1)
  标准旧版没有:
    - POST /AlarmClearAll        批量一键清除（补丁新增）

  ★ 从标准版 v0 升级时，必须同时替换 ism_server 与 dist，不能只换前端！

【前提】
  - Linux x86_64 服务器
  - 已有 ISM 部署
  - root 或 sudo 权限

【安装三步】
  1. 上传 alarm-clear-all-v1.zip 到服务器并解压:
       unzip alarm-clear-all-v1.zip -d /tmp/ism-patch
  2. 安装前探测（可选）:
       bash /tmp/ism-patch/alarm-clear-all-v1/verify_old_standard.sh
  3. 进入目录并执行安装（自动探测路径，也可手动指定）:
       cd /tmp/ism-patch/alarm-clear-all-v1
       sudo bash install.sh
       # 或手动指定:
       sudo ISM_BIN=/opt/ism/ism_server ISM_DIST=/opt/ism/web/dist bash install.sh
  4. 浏览器登录 ISM -> 实时告警 -> 点击「一键清除」验证

【三种使用方式】
  A. 完整补丁（推荐）: install.sh 替换 ism_server + dist -> 一键清除 API
  B. 前端降级模式: 仅当后端未更新时，UI 自动逐条调 AlarmOpt（≤500 条）
  C. 纯脚本清库: python3 clear_all_alarms.py（不依赖任何新 API）

【回滚】
  sudo bash rollback.sh
  # 或指定备份时间戳: sudo bash rollback.sh 20260627143000

【应急脚本】
  python3 clear_all_alarms.py --dry-run   # 仅统计
  python3 clear_all_alarms.py               # 直接清库

【注意】
  - 安装前会自动备份 .bak.时间戳
  - 重启后端后若合闸点位仍配置 IsAlarm=1，采集首轮可能再次产生告警
  - 告警超过 6 万条时请勿用 UI 逐条清除，请用 clear_all_alarms.py
"""


def parse_target(value: str) -> tuple[str, str]:
    m = re.fullmatch(r"([a-z0-9_]+)/([a-z0-9_]+)", value.lower())
    if not m:
        raise argparse.ArgumentTypeError(f"无效 target: {value}，示例 linux/amd64")
    return m.group(1), m.group(2)


def check_markers() -> list[str]:
    errors: list[str] = []
    for rel, needles in MARKERS.items():
        path = ROOT / rel
        if not path.exists():
            errors.append(f"文件不存在: {rel}")
            continue
        text = path.read_text(encoding="utf-8", errors="replace")
        if not any(n in text for n in needles):
            errors.append(f"marker 缺失 {rel}: 需要 {needles}")
    return errors


def run(cmd: list[str], *, cwd: Path | None = None, env: dict | None = None, check: bool = True) -> subprocess.CompletedProcess:
    print(f"  $ {' '.join(cmd)}")
    return subprocess.run(cmd, cwd=cwd or ROOT, env=env, check=check)


def detect_binary_platform(path: Path) -> str:
    if not path.is_file():
        return "missing"
    try:
        cp = subprocess.run(
            ["file", "-b", str(path)],
            capture_output=True,
            text=True,
            check=True,
        )
        return cp.stdout.strip() or "unknown"
    except (OSError, subprocess.CalledProcessError):
        return "unknown"


def build_backend_nocgo(goos: str, goarch: str, out: Path) -> bool:
    """CGO_ENABLED=0 静态交叉编译（main 已拆分 compile_datetime_nocgo）。"""
    env = os.environ.copy()
    env["GOOS"] = goos
    env["GOARCH"] = goarch
    env["CGO_ENABLED"] = "0"
    out.parent.mkdir(parents=True, exist_ok=True)
    try:
        run(["go", "build", "-o", str(out), "."], cwd=SERVER_DIR, env=env)
        return out.is_file()
    except subprocess.CalledProcessError:
        return False


def build_backend_native(goos: str, goarch: str, out: Path, *, cc: str | None = None) -> bool:
    env = os.environ.copy()
    env["GOOS"] = goos
    env["GOARCH"] = goarch
    env["CGO_ENABLED"] = "1"
    if cc:
        env["CC"] = cc
    out.parent.mkdir(parents=True, exist_ok=True)
    try:
        run(["go", "build", "-o", str(out), "."], cwd=SERVER_DIR, env=env)
        return True
    except subprocess.CalledProcessError:
        return False


def build_backend_docker(goos: str, goarch: str, out: Path) -> bool:
    if shutil.which("docker") is None:
        print("  docker 不可用，跳过容器构建")
        return False
    out.parent.mkdir(parents=True, exist_ok=True)
    out_abs = out.resolve()
    server_abs = SERVER_DIR.resolve()
    docker_cmd = [
        "docker", "run", "--rm",
        "-v", f"{server_abs}:/src",
        "-w", "/src",
        "-e", f"GOOS={goos}",
        "-e", f"GOARCH={goarch}",
        "-e", "CGO_ENABLED=1",
        "golang:bookworm",
        "go", "build", "-o", "/out/ism_server", ".",
    ]
    # 挂载输出目录
    docker_cmd = [
        "docker", "run", "--rm",
        "-v", f"{server_abs}:/src",
        "-v", f"{out_abs.parent}:/out",
        "-w", "/src",
        "-e", f"GOOS={goos}",
        "-e", f"GOARCH={goarch}",
        "-e", "CGO_ENABLED=1",
        "golang:bookworm",
        "go", "build", "-o", "/out/ism_server", ".",
    ]
    print(f"\n[backend] Docker golang:bookworm GOOS={goos} GOARCH={goarch}")
    try:
        run(docker_cmd)
        return out.is_file()
    except subprocess.CalledProcessError as e:
        print(f"  Docker 构建失败: 退出码 {e.returncode}")
        return False


def build_backend(goos: str, goarch: str, out: Path) -> tuple[str, str, str]:
    """构建后端，返回 (actual_goos, actual_goarch, platform_note)"""
    print(f"\n[backend] 目标 GOOS={goos} GOARCH={goarch} -> {out}")

    # 1) CGO_ENABLED=0 静态交叉编译（推荐，无需 C 工具链）
    print(f"  尝试 CGO_ENABLED=0 GOOS={goos} GOARCH={goarch}")
    if build_backend_nocgo(goos, goarch, out):
        plat = detect_binary_platform(out)
        if goos == "linux" and "ELF" in plat and "x86-64" in plat:
            print(f"  构建成功: {plat}")
            return goos, goarch, plat
        if goos == "linux" and "ELF" in plat:
            print(f"  构建成功: {plat}")
            return goos, goarch, plat
        print(f"  产物平台不符预期: {plat}，继续尝试其他方式")
        out.unlink(missing_ok=True)

    cc_candidates: list[str | None] = [None]
    if goos == "linux" and goarch == "amd64" and sys.platform == "darwin":
        for cc in ("x86_64-linux-gnu-gcc", "x86_64-unknown-linux-gnu-gcc"):
            if shutil.which(cc):
                cc_candidates.insert(0, cc)

    # 2) 本机 CGO 交叉编译
    for cc in cc_candidates:
        label = f"CGO_ENABLED=1 GOOS={goos} GOARCH={goarch}" + (f" CC={cc}" if cc else "")
        print(f"  尝试 {label}")
        if build_backend_native(goos, goarch, out, cc=cc):
            plat = detect_binary_platform(out)
            if goos == "linux" and "ELF" in plat and "x86-64" in plat:
                print(f"  构建成功: {plat}")
                return goos, goarch, plat
            if goos == "linux" and "ELF" in plat:
                print(f"  构建成功: {plat}")
                return goos, goarch, plat
            print(f"  产物平台不符预期: {plat}，继续尝试其他方式")
            out.unlink(missing_ok=True)

    # 3) Docker 容器内构建 linux/amd64
    if goos == "linux":
        if build_backend_docker(goos, goarch, out):
            plat = detect_binary_platform(out)
            print(f"  Docker 构建成功: {plat}")
            return goos, goarch, plat

    raise SystemExit(
        f"后端构建失败: 无法产出 {goos}/{goarch} Linux ELF。"
        "请确认 compile_datetime_nocgo.go 存在，或安装 x86_64-linux-gnu-gcc / 启动 Docker 后重试。"
    )


def _rm_tree(path: Path) -> None:
    if not path.exists():
        return
    if shutil.which("rm"):
        run(["rm", "-rf", str(path)])
    else:
        shutil.rmtree(path, ignore_errors=True)


def copy_dist_tree(src: Path, dest: Path) -> None:
    """复制 dist 目录内容到 dest，避免 cp -a src 已存在 dest 时产生 dest/dist 嵌套。"""
    _rm_tree(dest)
    dest.parent.mkdir(parents=True, exist_ok=True)
    if shutil.which("cp"):
        dest.mkdir(parents=True, exist_ok=True)
        run(["cp", "-a", f"{src}/.", str(dest)])
    else:
        shutil.copytree(src, dest, symlinks=True, ignore_dangling_symlinks=True)


def build_frontend(dist_dest: Path) -> None:
    print("\n[frontend] npm run build")
    env = os.environ.copy()
    env["NODE_OPTIONS"] = "--max-old-space-size=8192 --openssl-legacy-provider"
    run(["npm", "run", "build"], cwd=FRONTEND_DIR, env=env)
    src = FRONTEND_DIR / "dist"
    if not src.is_dir():
        raise SystemExit(f"前端 build 未产出 dist: {src}")
    copy_dist_tree(src, dist_dest)
    print(f"  dist 已复制到 {dist_dest}")


def verify_frontend_dist(dist_dir: Path) -> bool:
    if not dist_dir.is_dir():
        return False
    keywords = ("AlarmClearAll", "ClearAllAlarm", "ClearAllCurrentAlarm", "ClearAll")
    for js in dist_dir.rglob("*.js"):
        try:
            text = js.read_text(encoding="utf-8", errors="ignore")
        except OSError:
            continue
        if any(k in text for k in keywords):
            print(f"  前端验证 OK: 在 {js.relative_to(dist_dir)} 找到关键字")
            return True
    print("  警告: dist JS 中未找到 AlarmClearAll/ClearAll 关键字")
    return False


def write_patch_info(
    patch_dir: Path,
    *,
    goos: str,
    goarch: str,
    has_frontend: bool,
    binary_platform: str,
    platform_note: str = "",
) -> None:
    lines = [
        f"补丁名称: {PATCH_NAME}",
        f"构建时间: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}",
        f"目标平台: {goos}/{goarch}",
        f"ism_server file 输出: {binary_platform}",
    ]
    if platform_note:
        lines.append(f"平台说明: {platform_note}")
    lines.extend([
        f"包含前端: {'是' if has_frontend else '否（仅后端）'}",
        "",
        "标准旧版兼容:",
        "  A. install.sh 完整替换 ism_server + dist（推荐）",
        "  B. 仅 dist — 前端降级逐条清除（≤500 条）",
        "  C. clear_all_alarms.py 应急清库",
        "",
        "源码 marker:",
    ])
    for rel, needles in MARKERS.items():
        lines.append(f"  - {rel}: {', '.join(needles)}")
    lines.extend([
        "",
        "兼容说明:",
        "  - 标准旧版无 POST /AlarmClearAll，升级须同时替换 ism_server + dist",
        "  - 旧版已有 POST /AlarmOpt（单条清除）与 GetCurrentAlarmList",
        "  - 前端含降级逻辑：API 不可用时逐条 AlarmOpt（≤500 条）",
        "  - 应急脚本 clear_all_alarms.py 不依赖新 API",
        "",
        "包内文件:",
        "  ism_server              - 后端二进制",
        "  dist/                   - 前端静态资源（若已编译）",
        "  install.sh              - 安装脚本（自动探测路径）",
        "  rollback.sh             - 回滚脚本",
        "  verify_old_standard.sh  - 探测是否已有 AlarmClearAll",
        "  clear_all_alarms.py     - 应急清库脚本",
        "  现场安装说明.txt",
        "  PATCH_INFO.txt",
    ])
    (patch_dir / "PATCH_INFO.txt").write_text("\n".join(lines) + "\n", encoding="utf-8")


def verify_patch_artifacts(patch_dir: Path, *, has_frontend: bool) -> None:
    required = [
        "ism_server",
        "install.sh",
        "rollback.sh",
        "verify_old_standard.sh",
        "clear_all_alarms.py",
        "现场安装说明.txt",
    ]
    missing = [name for name in required if not (patch_dir / name).exists()]
    if has_frontend and not (patch_dir / "dist").is_dir():
        missing.append("dist/")
    if missing:
        raise SystemExit(f"补丁目录缺少文件: {', '.join(missing)}")


def acquire_build_lock() -> int:
    lock_path = PATCH_DIR.parent / f".{PATCH_NAME}.build.lock"
    lock_path.parent.mkdir(parents=True, exist_ok=True)
    fd = os.open(str(lock_path), os.O_CREAT | os.O_RDWR)
    try:
        fcntl.flock(fd, fcntl.LOCK_EX | fcntl.LOCK_NB)
    except BlockingIOError:
        os.close(fd)
        raise SystemExit("另一个补丁构建正在进行，请稍后再试")
    return fd


def release_build_lock(fd: int) -> None:
    try:
        fcntl.flock(fd, fcntl.LOCK_UN)
    finally:
        os.close(fd)


def create_zip(patch_dir: Path, zip_path: Path) -> None:
    print(f"\n[zip] {zip_path}")
    zip_path = zip_path.resolve()
    tmp_path = zip_path.with_suffix(".zip.tmp")
    zip_path.parent.mkdir(parents=True, exist_ok=True)
    if tmp_path.exists():
        tmp_path.unlink()
    file_count = 0
    with zipfile.ZipFile(str(tmp_path), "w", zipfile.ZIP_DEFLATED) as zf:
        for f in sorted(patch_dir.rglob("*")):
            if not f.is_file():
                continue
            try:
                arc = f.relative_to(patch_dir.parent)
                zf.write(f, str(arc))
                file_count += 1
            except OSError as e:
                print(f"  警告: 跳过 {f}: {e}")
    if not tmp_path.is_file():
        raise SystemExit(f"ZIP 未生成: {tmp_path}（已打包 {file_count} 个文件）")
    if zip_path.exists():
        zip_path.unlink()
    tmp_path.replace(zip_path)
    size_mb = zip_path.stat().st_size / (1024 * 1024)
    print(f"  文件数: {file_count}")
    print(f"  大小: {zip_path.stat().st_size:,} bytes ({size_mb:.2f} MB)")


def main() -> int:
    parser = argparse.ArgumentParser(description="构建告警一键清除补丁包")
    parser.add_argument("--check-only", action="store_true", help="仅检查源码 marker")
    parser.add_argument("--skip-frontend", action="store_true", help="跳过后端 npm build")
    parser.add_argument("--target", type=parse_target, default=("linux", "amd64"), help="GOOS/GOARCH，默认 linux/amd64")
    args = parser.parse_args()
    goos, goarch = args.target

    print("=== 检查源码 marker ===")
    errors = check_markers()
    if errors:
        for e in errors:
            print(f"  FAIL: {e}")
        return 1
    for rel in MARKERS:
        print(f"  OK: {rel}")
    if args.check_only:
        print("\n--check-only: marker 检查通过")
        return 0

    lock_fd = acquire_build_lock()
    try:
        return _build_patch(args, goos, goarch)
    finally:
        release_build_lock(lock_fd)


def _build_patch(args: argparse.Namespace, goos: str, goarch: str) -> int:
    if PATCH_DIR.exists():
        _rm_tree(PATCH_DIR)
    PATCH_DIR.mkdir(parents=True)

    try:
        actual_goos, actual_goarch, platform_note = build_backend(goos, goarch, PATCH_DIR / "ism_server")
        ism_server_path = PATCH_DIR / "ism_server"
        if not ism_server_path.is_file():
            raise SystemExit(f"后端构建未产出 ism_server: {ism_server_path}")
        binary_platform = detect_binary_platform(ism_server_path)

        has_frontend = False
        if not args.skip_frontend:
            build_frontend(PATCH_DIR / "dist")
            has_frontend = True
            verify_frontend_dist(PATCH_DIR / "dist")
        else:
            print("\n[frontend] 已跳过 (--skip-frontend)")

        (PATCH_DIR / "install.sh").write_text(INSTALL_SH, encoding="utf-8")
        (PATCH_DIR / "rollback.sh").write_text(ROLLBACK_SH, encoding="utf-8")
        (PATCH_DIR / "verify_old_standard.sh").write_text(VERIFY_OLD_STANDARD_SH, encoding="utf-8")
        (PATCH_DIR / "现场安装说明.txt").write_text(INSTALL_README, encoding="utf-8")
        os.chmod(PATCH_DIR / "install.sh", 0o755)
        os.chmod(PATCH_DIR / "rollback.sh", 0o755)
        os.chmod(PATCH_DIR / "verify_old_standard.sh", 0o755)
        os.chmod(PATCH_DIR / "ism_server", 0o755)

        if CLEAR_SCRIPT.exists():
            shutil.copy2(CLEAR_SCRIPT, PATCH_DIR / "clear_all_alarms.py")
            os.chmod(PATCH_DIR / "clear_all_alarms.py", 0o755)
        else:
            print(f"  警告: 未找到 {CLEAR_SCRIPT}")

        write_patch_info(
            PATCH_DIR,
            goos=actual_goos,
            goarch=actual_goarch,
            has_frontend=has_frontend,
            binary_platform=binary_platform,
            platform_note=platform_note if actual_goos != goos or actual_goarch != goarch else "",
        )
        verify_patch_artifacts(PATCH_DIR, has_frontend=has_frontend)
        create_zip(PATCH_DIR, PATCH_ZIP)

        print("\n=== 构建完成 ===")
        print(f"补丁目录: {PATCH_DIR}")
        print(f"ZIP:      {PATCH_ZIP} ({PATCH_ZIP.stat().st_size:,} bytes)")
        print(f"ism_server: {binary_platform}")
        return 0
    except subprocess.CalledProcessError as e:
        print(f"\n构建失败: 命令退出码 {e.returncode}", file=sys.stderr)
        return e.returncode or 1
    except SystemExit as e:
        print(f"\n构建失败: {e}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    raise SystemExit(main())
