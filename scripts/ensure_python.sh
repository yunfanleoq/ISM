#!/bin/bash
# 检测 Python3；若无可用版本则从包内 python-offline/ 解压/启用便携 Python
# 用法: ISM_PYTHON=$(bash scripts/ensure_python.sh)  或  source <(bash scripts/ensure_python.sh --export)
#  stdout: python3 可执行文件绝对路径
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
OFFLINE="$ROOT/python-offline"
INSTALL="$OFFLINE/install"
MIN_MAJOR=3
MIN_MINOR=6

_export=0
[[ "${1:-}" == "--export" ]] && _export=1

python_ok() {
  local py="$1"
  [[ -x "$py" ]] || return 1
  "$py" -c "import sys; raise SystemExit(0 if sys.version_info >= (${MIN_MAJOR}, ${MIN_MINOR}) else 1)" 2>/dev/null
}

find_bundled() {
  for c in "$INSTALL/bin/python3" "$INSTALL/bin/python3.11" "$INSTALL/bin/python3.12"; do
    if python_ok "$c"; then
      echo "$c"
      return 0
    fi
  done
  return 1
}

extract_bundled() {
  local tar
  tar=$(ls "$OFFLINE"/cpython-*-install_only.tar.gz 2>/dev/null | head -1)
  if [[ -z "$tar" || ! -f "$tar" ]]; then
    return 1
  fi
  echo "  解压包内 Python: $(basename "$tar")" >&2
  rm -rf "$INSTALL"
  mkdir -p "$INSTALL"
  tar -xzf "$tar" -C "$INSTALL"
  if [[ -d "$INSTALL/python/bin" ]]; then
    shopt -s dotglob nullglob
    mv "$INSTALL/python/"* "$INSTALL/"
    rm -rf "$INSTALL/python"
    shopt -u dotglob nullglob
  fi
  find_bundled
}

install_system_link() {
  # 可选：root 时将便携 Python 链到 /usr/local/bin/ism-python3
  local py="$1"
  if [[ "$(id -u)" -eq 0 ]] && [[ -x "$py" ]]; then
    install -m 755 "$py" /usr/local/bin/ism-python3 2>/dev/null || true
    ln -sf /usr/local/bin/ism-python3 /usr/local/bin/python3 2>/dev/null || true
  fi
}

resolve_python() {
  if command -v python3 >/dev/null 2>&1 && python_ok python3; then
    echo "  使用系统 Python: $(python3 --version 2>&1)" >&2
    echo "$(command -v python3)"
    return 0
  fi

  if bundled=$(find_bundled); then
    echo "  使用包内 Python: $($bundled --version 2>&1)" >&2
    install_system_link "$bundled"
    echo "$bundled"
    return 0
  fi

  echo "  系统无 Python3，启用包内离线 Python ..." >&2
  if bundled=$(extract_bundled); then
    install_system_link "$bundled"
    echo "$bundled"
    return 0
  fi

  echo "错误: 未找到可用 Python3，且包内 python-offline/ 不完整" >&2
  echo "  请确认部署包含 python-offline/install/ 或 cpython-*-install_only.tar.gz" >&2
  return 1
}

PY="$(resolve_python)"

if [[ "$_export" == "1" ]]; then
  echo "export ISM_PYTHON='$PY'"
else
  echo "$PY"
fi
