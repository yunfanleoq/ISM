#!/usr/bin/env bash
# ============================================================
# ISM Web SCADA — Cloud Agent start script
# Per-boot runtime reconciliation (idempotent, must return).
# The actual servers run as `terminals` (backend + frontend);
# this only prepares config/state they depend on:
#   1. Use the bundled SQLite DB (dbtype=1) for a self-contained
#      dev backend, overriding the committed OceanBase default (=4).
#   2. Create the runtime directories the backend expects.
#   3. Normalize the demo admin login to  admin / 123456 .
# ============================================================
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT/ism_server_user"

# 1) Dev database = bundled SQLite. The committed app.conf ships dbtype=4
#    (OceanBase, external) for production; a fresh checkout resets it, so we
#    re-apply the dev override on every boot. Idempotent.
sed -i 's/^dbtype=.*/dbtype=1/' conf/app.conf
echo "[start] backend DB set to SQLite (dbtype=1)"

# 2) Runtime directories the backend creates data under.
mkdir -p data/auth data/sessionon static/HistoryData static/reportTemplete static/RecordVideo logs

# 3) Deterministic dev credentials: admin / 123456.
#    Frontend sends MD5(password); backend stores bcrypt(MD5(password)).
#    The literal below is bcrypt(MD5("123456")) — the same value db.go seeds
#    for a fresh admin — so no bcrypt dependency is needed here.
python3 - <<'PY'
import os, sqlite3
db = os.path.join("data", "db", "ism.db")
if os.path.exists(db):
    known = "$2a$10$h9swLjbTTcSVUCqQDt6nAetw.FVRLPE0WPDzqloprYRO7PDtLC5Ii"  # bcrypt(MD5("123456"))
    con = sqlite3.connect(db)
    con.execute("UPDATE user SET password=? WHERE username='admin'", (known,))
    con.commit()
    con.close()
    print("[start] admin login normalized to admin/123456")
else:
    print("[start] no seeded DB found; db.go will create admin (admin/123456) on first run")
PY

echo "[start] done."
