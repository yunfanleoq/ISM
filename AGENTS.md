# AGENTS.md

## Cursor Cloud specific instructions

This repo is the **ISM** web SCADA / configuration ("组态") monitoring platform:
a **Go/beego backend** (`ism_server_user`, port **8081** http + 5622 https) and a
**Vue 2 / vue-cli frontend** (`ism-front-end-v2`, port **7080**, which proxies `/api` → `127.0.0.1:8081`).
Default login: `admin` / `123456`.

> Note: the many `.cursor/rules/*.mdc` and `.cursor/skills/*` files are written for the
> maintainer's **macOS** workstation with lots of RAM. This Cursor Cloud VM is **Linux with ~15 GB
> RAM and no swap**, so several of those rules do NOT apply here (details below). Prefer the
> guidance in this section for the cloud VM.

### Backend (`ism_server_user`, Go)
- Build (needs `gcc`; CGO is required for the SQLite driver): from `ism_server_user/`
  `CGO_ENABLED=1 go build -mod=vendor -o ism_server .`
- Run from **inside `ism_server_user/`** (paths are relative: `conf/`, `data/db/ism.db`, `logs/`):
  `./ism_server` — logs to `logs/ism.log` and stdout.
- Verify: `curl -s http://127.0.0.1:8081/login -X POST -H 'Content-Type: application/json' -d '{"Username":"admin","password":"e10adc3949ba59abbe56e057f20f883e"}'` → expect `"code":1000`
  (that hash is `md5("123456")`; the frontend MD5s the password — see `.cursor/rules/password-chain.mdc`).
- **Database**: `conf/app.conf` `dbtype` selects the DB. This branch sets `dbtype=1` (SQLite,
  using the committed `data/db/ism.db`). The upstream default is `dbtype=4` (OceanBase on :2881),
  which is unavailable locally and makes the backend panic on startup — keep `dbtype=1` for local dev.
- **Vendored-deps gotcha (important):** the root `.gitignore` has a broad `logs/` pattern that
  accidentally excludes the vendored `vendor/github.com/beego/beego/v2/core/logs` and
  `.../adapter/logs` package dirs. They are **force-added** on this branch. If a clean checkout is
  missing them, `go build -mod=vendor` fails with `cannot find module providing package
  .../core/logs`. Restore with:
  `go mod download github.com/beego/beego/v2` then copy `core/logs` and `adapter/logs` from
  `$(go env GOMODCACHE)/github.com/beego/beego/v2@v2.1.0/` into `vendor/github.com/beego/beego/v2/`
  (delete the copied `alils`/`es` subdirs and `*_test.go`).
  Do **not** run `go mod vendor` / `go mod tidy` to "fix" vendoring — the vendor tree contains
  patched forks (modbus / go645 / opcua / iec104) that differ from the go.mod-pinned upstream, so
  module mode (`-mod=mod`) fails to compile. Always build with `-mod=vendor`.
- Optional dependencies the backend logs `connection refused` for but does NOT need to serve:
  TDengine history DB (`:6041`), device protocol endpoints (Modbus `:502`, etc.), external MQTT.

### Frontend (`ism-front-end-v2`, Vue 2 + vue-cli 4 / webpack 4)
- Install deps with **npm**, not yarn: `npm ci --legacy-peer-deps`
  (`yarn install` fails on a git-dependency `Invariant Violation: Commit hash required` bug;
  `--legacy-peer-deps` is required for a `webpack-obfuscator` peer conflict).
- Dev server: `NODE_OPTIONS="--max-old-space-size=8192 --openssl-legacy-provider" npx vue-cli-service serve --port 7080`
  - `--openssl-legacy-provider` is required for webpack 4 on modern Node (Node 22 here; `.nvmrc`=16 is stale).
  - First compile takes ~1 min over ~8,300 modules and **sits at 62–69% for a while — that is normal**,
    not a hang. Wait for `Compiled successfully`.
  - Ignore the macOS-only `scripts/check_mem_before_compile.sh` (it hard-exits FAIL on non-Darwin) and
    the "20 GB Node heap / no setsid / launchctl" rules — those are macOS-specific. 8 GB heap fits this VM.
- Verify proxy end-to-end once compiled:
  `curl http://127.0.0.1:7080/ -o /dev/null -w '%{http_code}\n'` → `200`, and
  `curl http://127.0.0.1:7080/api/login -X POST -H 'Content-Type: application/json' -d '{"Username":"admin","password":"e10adc3949ba59abbe56e057f20f883e"}'` → `"code":1000`.
- Lint: `npx vue-cli-service lint` runs but reports **pre-existing** errors in the checked-in minified
  vendor file `src/assets/js/axios.min.js` (not application source) — unrelated to your changes.

### Known UI caveat
Some heavy configuration pages (e.g. `#/Setting/UserManager`) can crash the Chrome tab
("Aw, Snap!") under memory pressure on this constrained VM. Core flows (login, project CRUD, dashboards)
work; if a specific page crashes, it is usually tab memory pressure, not a backend failure.
