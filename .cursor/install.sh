#!/usr/bin/env bash
# ============================================================
# ISM Web SCADA — Cloud Agent install script
# Idempotent, source-derived setup run after checkout:
#   1. Build the Go backend (ism_server_user) with CGO+SQLite.
#   2. Install the Vue frontend dependencies (ism-front-end-v2).
# Long-running servers live in `terminals`; per-boot runtime
# reconciliation lives in start.sh.
# ============================================================
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# ----------------------------------------------------------------
# Backend: Go build (vendored deps, patched → must use -mod=vendor)
# ----------------------------------------------------------------
cd "$ROOT/ism_server_user"

# go.mod pins `toolchain go1.22.7`; let Go fetch it if the base image
# ships an older 1.22.x. go-sqlite3 (dbtype=1) requires cgo.
export GOTOOLCHAIN=auto
export CGO_ENABLED=1

# Self-heal: the repo-root .gitignore has an unanchored `logs/` rule that can
# exclude these vendored beego packages from a checkout. Restore them from the
# module cache so `-mod=vendor` stays consistent even if they are absent.
for pkg in core/logs adapter/logs; do
  dst="vendor/github.com/beego/beego/v2/$pkg"
  if [ ! -f "$dst/log.go" ]; then
    echo "[install] restoring vendored beego/$pkg from module cache"
    go mod download github.com/beego/beego/v2 >/dev/null 2>&1 || true
    src="$(go env GOMODCACHE)/github.com/beego/beego/v2@v2.1.0/$pkg"
    mkdir -p "$dst"
    cp -r "$src/." "$dst/"
    chmod -R u+rwX "$dst"
    find "$dst" -name '*_test.go' -delete
  fi
done

echo "[install] building ism_server (CGO_ENABLED=1, -mod=vendor)…"
go build -mod=vendor -o ism_server .
echo "[install] backend build OK: $(ls -la ism_server | awk '{print $5" bytes"}')"

# ----------------------------------------------------------------
# Frontend: install dependencies (webpack 4 / vue-cli, legacy peers)
# ----------------------------------------------------------------
cd "$ROOT/ism-front-end-v2"
echo "[install] installing frontend dependencies (npm, legacy-peer-deps)…"
npm install --legacy-peer-deps
echo "[install] frontend dependencies OK"

echo "[install] done."
