#!/bin/bash
# 现场热修：大屏右上角齿轮进后台后「项目ID错误」、数据仓库空白。
# 只改前端静态文件，不重启服务，不碰数据库、ism_server、密码。
#
# 用法:
#   bash apply-keep-project-uuid.sh /opt/ISM/ism-release-oceanbase-20260709
# 若目录是 20260708，把路径换成实际主包目录。
# 套完后浏览器 Ctrl+F5。已经停在空白数据仓库的，回大屏再点一次齿轮。
set -euo pipefail

TARGET="${1:-}"
if [[ -z "$TARGET" ]]; then
  echo "用法: bash apply-keep-project-uuid.sh <主包目录>"
  echo "例如: bash apply-keep-project-uuid.sh /opt/ISM/ism-release-oceanbase-20260709"
  exit 1
fi

DIST=""
for cand in "$TARGET/web/dist" "$TARGET/dist" "$TARGET"; do
  if [[ -f "$cand/index.html" ]]; then
    DIST="$cand"
    break
  fi
done
if [[ -z "$DIST" ]]; then
  echo "错误: 在 $TARGET 下找不到 index.html（试过 web/dist、dist、目录本身）"
  exit 1
fi

JS_DIR="$DIST/static/js"
mkdir -p "$JS_DIR"
JS_FILE="$JS_DIR/ism-keep-project-uuid.js"
MARK="ism-keep-project-uuid.js"
STAMP="$(date +%Y%m%d%H%M%S)"
BACKUP="$DIST/backup-keep-project-uuid-$STAMP"
mkdir -p "$BACKUP"
cp -a "$DIST/index.html" "$BACKUP/index.html"
if [[ -f "$JS_FILE" ]]; then
  cp -a "$JS_FILE" "$BACKUP/ism-keep-project-uuid.js"
fi

cat > "$JS_FILE" << 'ENDJS'
/* ISM keep ProjectUuid: ProjectList/Add/Edit/Del 只去掉本次请求头，不清除已保存的项目 ID */
(function () {
  var KEY = 'ProjectUuid'
  var MARK = '__ismKeepProjectUuid'

  function readProjectUuid() {
    try {
      var ss = sessionStorage.getItem(KEY)
      if (ss) return ss
    } catch (e) {}
    var m = document.cookie.match(/(?:^|; )ProjectUuid=([^;]*)/)
    return m ? decodeURIComponent(m[1]) : ''
  }

  function writeProjectUuid(v) {
    if (!v) return
    try { sessionStorage.setItem(KEY, v) } catch (e) {}
    document.cookie = KEY + '=' + encodeURIComponent(v) + '; path=/'
  }

  function shouldKeep(url) {
    url = String(url || '')
    if (url.indexOf('login') !== -1) return false
    return url.indexOf('ProjectList') !== -1 ||
      url.indexOf('ProjectAdd') !== -1 ||
      url.indexOf('ProjectEdit') !== -1 ||
      url.indexOf('ProjectDel') !== -1
  }

  function wrapFulfilled(orig) {
    if (!orig || orig[MARK]) return orig
    function wrapped(config) {
      var url = config && config.url
      var saved = shouldKeep(url) ? readProjectUuid() : ''
      var ret = orig(config)
      if (!saved) return ret
      if (ret && typeof ret.then === 'function') {
        return ret.then(function (c) {
          writeProjectUuid(saved)
          return c
        }, function (err) {
          writeProjectUuid(saved)
          throw err
        })
      }
      writeProjectUuid(saved)
      return ret
    }
    wrapped[MARK] = true
    return wrapped
  }

  function patchAxios() {
    var ax = window.axios
    if (!ax || !ax.interceptors || !ax.interceptors.request) return false
    var mgr = ax.interceptors.request
    var handlers = mgr.handlers
    if (!handlers || !handlers.length) return false
    for (var i = 0; i < handlers.length; i++) {
      var h = handlers[i]
      if (h && h.fulfilled) h.fulfilled = wrapFulfilled(h.fulfilled)
    }
    if (!mgr.use[MARK]) {
      var origUse = mgr.use
      var wrappedUse = function (fulfilled, rejected) {
        return origUse.call(mgr, wrapFulfilled(fulfilled), rejected)
      }
      wrappedUse[MARK] = true
      mgr.use = wrappedUse
    }
    return true
  }

  var n = 0
  var timer = setInterval(function () {
    n += 1
    if (patchAxios() || n > 200) clearInterval(timer)
  }, 50)
})()
ENDJS
chmod 644 "$JS_FILE"

SRC="$MARK"
if grep -q 'src="/static/js/' "$DIST/index.html"; then
  SRC="/static/js/$MARK"
elif grep -q 'src="./static/js/' "$DIST/index.html"; then
  SRC="./static/js/$MARK"
elif grep -q 'src="static/js/' "$DIST/index.html"; then
  SRC="static/js/$MARK"
fi

if grep -q "$MARK" "$DIST/index.html"; then
  echo "[跳过] index.html 已引用 ${MARK}，只更新了 js"
else
  # 插在 </body> 前，保证排在页面原有脚本之后
  python3 - "$DIST/index.html" "$SRC" << 'PY'
import sys
path, src = sys.argv[1], sys.argv[2]
html = open(path, encoding="utf-8", errors="replace").read()
tag = '<script src="%s"></script>' % src
idx = html.lower().rfind("</body>")
if idx < 0:
    sys.stderr.write("错误: index.html 没有 </body>\n")
    sys.exit(1)
html = html[:idx] + tag + "\n" + html[idx:]
open(path, "w", encoding="utf-8").write(html)
PY
  echo "[写入] index.html 已挂上 $SRC"
fi

if [[ -f "$DIST/index.html.gz" ]]; then
  if command -v gzip >/dev/null 2>&1; then
    gzip -c -9 "$DIST/index.html" > "$DIST/index.html.gz"
    echo "[更新] index.html.gz"
  else
    rm -f "$DIST/index.html.gz"
    echo "[删除] 无 gzip 命令，已去掉旧的 index.html.gz，避免继续下发旧页面"
  fi
fi

echo "完成。备份在 $BACKUP"
echo "请让浏览器 Ctrl+F5。若当前已是空白的数据仓库，先回大屏，再点右上角齿轮。"
echo "本脚本不重启服务，不修改数据库。"
