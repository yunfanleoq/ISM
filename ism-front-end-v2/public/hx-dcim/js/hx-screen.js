/* 航信一次图和列头柜画面。布局来自他们画的工程，取数走 HXLive。 */
(function (w) {
  'use strict';

  var L = function () { return w.HXLive; };
  var scenes = [];
  var modelId = '';
  /** 配电室名 → pageId（静态 index + API meta 合并） */
  var powerByName = {};
  var powerIndexP = null;
  var page = null;
  var pageRows = [];
  var hoverRows = [];
  var hotKey = '';
  var pinned = false;
  var hoverAnchor = null;
  var hostEl = null;
  var onResize = null;
  var treeFlat = null;
  var drillPick = null;
  var drillActive = null;
  var treeMode = 'full';

  function esc(s) {
    return String(s == null ? '' : s)
      .replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
  }

  function apiPost(path, body) {
    var h = { 'Content-Type': 'application/json' };
    if (L().headers) {
      var lh = L().headers();
      Object.keys(lh).forEach(function (k) { h[k] = lh[k]; });
    }
    return fetch('/api' + path, {
      method: 'POST',
      headers: h,
      body: JSON.stringify(body || {})
    }).then(function (r) { return r.json(); });
  }

  function mergePowerIndex(name, pageId) {
    if (!name || !pageId || !w.HXSceneImport || !w.HXSceneImport.powerPageRe.test(name)) return;
    powerByName[name] = pageId;
  }

  function loadPowerIndex() {
    if (powerIndexP) return powerIndexP;
    powerIndexP = fetch('/hx-dcim/data/hx-power-page-index.json?v=20261102f5').then(function (r) {
      if (!r.ok) return null;
      return r.json();
    }).then(function (j) {
      if (!j) return;
      if (j.model && !modelId) modelId = j.model;
      (j.pages || []).forEach(function (p) {
        mergePowerIndex(p.name, p.pageId);
      });
    }).catch(function () {});
    return powerIndexP;
  }

  function syncPowerMetaFromApi() {
    if (!modelId) return Promise.resolve();
    return apiPost('/getDisplayModelLayerData', { muid: modelId, metaOnly: true }).then(function (res) {
      if (!res || res.code !== 0) return;
      (res.layer || []).forEach(function (row) {
        mergePowerIndex(row.PageName || row.page_name, row.PageId || row.page_id);
      });
    }).catch(function () {});
  }

  function loadScenes() {
    var jobs = [];
    if (w.HX5A1) {
      jobs.push(w.HX5A1.load('20261107merge').catch(function (e) {
        if (w.console) w.console.warn('[hx-5a1] 布局未加载，5A1 退回通用绘制', e);
      }));
    }
    if (w.HXRooms) {
      jobs.push(w.HXRooms.load('20261107merge').catch(function (e) {
        if (w.console) w.console.warn('[hx-rooms] 配电室布局未加载', e);
      }));
    }
    return Promise.all(jobs).then(function () {
      return fetch('/hx-dcim/data/hx-scenes.json?v=20261102f5');
    }).then(function (r) {
      if (!r.ok) throw new Error('no scenes');
      return r.json();
    }).then(function (j) {
      modelId = (j && j.model) || modelId;
      scenes = (j && j.pages) || [];
      scenes.forEach(function (p) {
        if (p.kind === 'power' && p.name && p.id) mergePowerIndex(p.name, p.id);
      });
    }).catch(function () { scenes = []; });
  }

  function ensureSceneLoaded(id) {
    var cached = pageById(id);
    if (cached) return Promise.resolve(cached);
    if (!id || !w.HXSceneImport) return Promise.resolve(null);
    return apiPost('/getDisplayModelPagerLayerData', { pageid: id }).then(function (res) {
      if (!res || res.code !== 0 || !res.layer) return null;
      var pg = w.HXSceneImport.buildPageFromLayerRow(res.layer);
      if (!pg || !pg.id) return null;
      var i;
      for (i = 0; i < scenes.length; i++) {
        if (scenes[i].id === pg.id) return scenes[i];
      }
      scenes.push(pg);
      if (pg.kind === 'power') mergePowerIndex(pg.name, pg.id);
      return pg;
    }).catch(function () { return null; });
  }

  function devOf(key) {
    var i, n;
    if (!page) return '';
    for (i = 0; i < page.nodes.length; i++) {
      n = page.nodes[i];
      if (n.k === key && n.dev) return n.dev;
    }
    return '';
  }

  function paintTopoDevs() {
    if (!hostEl) return;
    var root = hostEl.querySelector('.hx-topo');
    if (!root) return;
    root.querySelectorAll('.hx-topo-dev[data-dev]').forEach(function (el) {
      var st = topoDevState(el.getAttribute('data-dev'), el.getAttribute('data-uuid'));
      applyTopoDevEl(el, st);
    });
  }

  /** 合闸=红=状态1，分闸=绿=状态0。优先读「开关状态」，不用电压>0 冒充合闸。 */
  function knifeIsClosed(n) {
    if (!n) return false;
    var live = L();
    var base = String(n.k || '').replace(/_+$/, '');
    var statusName = base ? base + '_开关状态' : '';
    var pt = statusName && live.pointOf ? live.pointOf(statusName, n.dev) : null;
    if (!pt && n.p && String(n.p).slice(-4) === '开关状态' && live.pointOf) {
      pt = live.pointOf(n.p, n.dev);
    }
    if (pt) {
      var sv = live.uuidValue ? live.uuidValue(pt.uuid) : (pt.value != null ? String(pt.value) : null);
      if (sv == null || sv === '') sv = pt.value != null ? String(pt.value) : null;
      return sv != null && sv !== '' && live.matchOp(sv, '==', '1');
    }
    var v = live.nodeValue(n);
    if (n.op) return v != null && v !== '' && live.matchOp(v, n.op, n.val);
    return v != null && v !== '' && live.matchOp(v, '==', '1');
  }

  function lineEndsOf(n) {
    var cx = n.x + n.w / 2;
    var cy = n.y + (n.h || 26) / 2;
    var half = n.w / 2;
    if (Math.abs((Number(n.a) || 0) % 180) < 1) {
      return [{ x: cx - half, y: cy }, { x: cx + half, y: cy }];
    }
    return [{ x: cx, y: cy - half }, { x: cx, y: cy + half }];
  }

  function isGreenStroke(c) {
    var s = String(c || '').toLowerCase();
    return s === '#00ff00' || s === '#00e676';
  }

  function hitsPlate(ends, tx, pad) {
    var i;
    for (i = 0; i < ends.length; i++) {
      var e = ends[i];
      if (e.x < tx.x - pad || e.x > tx.x + (Number(tx.w) || 0) + pad) continue;
      if (e.y < tx.y - pad || e.y > tx.y + (Number(tx.h) || 0) + pad) continue;
      return true;
    }
    return false;
  }

  function normPaintColor(c) {
    var s = String(c || '').replace(/\s/g, '').toLowerCase();
    if (s === '#00ff00' || s === '#0f0' || s === 'rgb(0,255,0)') return '#00ff00';
    if (s === '#ff0000' || s === '#f00' || s === 'rgb(255,0,0)') return '#ff0000';
    return '';
  }

  /* 端点先收到 0.01，再当同一点。否则 116.25 和 116.25000000000001 会把同一只铭牌框拆成两截。 */
  function frameCornerKey(x, y) {
    return (Math.round(x * 100) / 100) + ',' + (Math.round(y * 100) / 100);
  }

  function groupFrameRects(segs) {
    var parent = [];
    var k;
    for (k = 0; k < segs.length; k++) parent.push(k);
    function find(a) {
      while (parent[a] !== a) {
        parent[a] = parent[parent[a]];
        a = parent[a];
      }
      return a;
    }
    function unite(a, b) {
      a = find(a);
      b = find(b);
      if (a !== b) parent[b] = a;
    }
    var at = {};
    for (k = 0; k < segs.length; k++) {
      var s = segs[k];
      var c1 = frameCornerKey(s.x1, s.y1);
      var c2 = frameCornerKey(s.x2, s.y2);
      if (at[c1] != null) unite(k, at[c1]);
      else at[c1] = k;
      if (at[c2] != null) unite(k, at[c2]);
      else at[c2] = k;
    }
    var groups = {};
    for (k = 0; k < segs.length; k++) {
      var id = find(k);
      if (!groups[id]) groups[id] = [];
      groups[id].push(segs[k]);
    }
    var rects = [];
    Object.keys(groups).forEach(function (id) {
      var g = groups[id];
      var x0 = Infinity;
      var y0 = Infinity;
      var x1 = -Infinity;
      var y1 = -Infinity;
      var i;
      for (i = 0; i < g.length; i++) {
        x0 = Math.min(x0, g[i].x1, g[i].x2);
        y0 = Math.min(y0, g[i].y1, g[i].y2);
        x1 = Math.max(x1, g[i].x1, g[i].x2);
        y1 = Math.max(y1, g[i].y1, g[i].y2);
      }
      rects.push({ x0: x0, y0: y0, x1: x1, y1: y1, segs: g });
    });
    return rects;
  }

  /* 铭牌只留 CAD 那一圈框，颜色跟字。文字上的 CSS 边框藏掉，避免绿框外面再套一圈红。
     边长超过 160 的框当房间框，不改色。母线和进线不是 net:frame，这里碰不到。 */
  function syncPlateFrames() {
    if (!page || page.kind !== 'power' || !hostEl) return;
    var nodes = page.nodes || [];
    var plates = [];
    hostEl.querySelectorAll('.hx-tx.hx-plate-live').forEach(function (el) {
      if (el.getAttribute('contenteditable') === 'true') return;
      var tx = nodes[Number(el.getAttribute('data-i'))];
      if (!tx || !tx.w || !tx.h) return;
      var col = normPaintColor(el.style.color);
      if (!col) return;
      plates.push({ el: el, tx: tx, col: col });
    });
    if (!plates.length) return;
    var segs = [];
    hostEl.querySelectorAll('[data-hx-line]').forEach(function (el) {
      if ((el.getAttribute('stroke') || '') === 'transparent') return;
      var n = nodes[Number(el.getAttribute('data-i'))];
      if (!n || n.t !== 'L' || n.net !== 'frame') return;
      var ends = lineEndsOf(n);
      if (!ends || ends.length < 2) return;
      segs.push({ el: el, x1: ends[0].x, y1: ends[0].y, x2: ends[1].x, y2: ends[1].y });
    });
    if (!segs.length) return;
    var rects = groupFrameRects(segs);
    var i;
    var j;
    for (i = 0; i < rects.length; i++) {
      var r = rects[i];
      var rw = r.x1 - r.x0;
      var rh = r.y1 - r.y0;
      if (!(rw >= 8 && rh >= 8 && rw <= 160 && rh <= 160)) continue;
      var best = null;
      var bestScore = Infinity;
      var area = rw * rh;
      for (j = 0; j < plates.length; j++) {
        var tx = plates[j].tx;
        var cx = tx.x + tx.w / 2;
        var cy = tx.y + tx.h / 2;
        if (cx < r.x0 - 2 || cx > r.x1 + 2 || cy < r.y0 - 2 || cy > r.y1 + 2) continue;
        var score = Math.abs(area - tx.w * tx.h);
        if (score < bestScore) {
          bestScore = score;
          best = plates[j];
        }
      }
      if (best) r.col = best.col;
    }
    for (i = 0; i < rects.length; i++) {
      var inner = rects[i];
      if (inner.col) continue;
      var iw = inner.x1 - inner.x0;
      var ih = inner.y1 - inner.y0;
      if (!(iw >= 8 && ih >= 8 && iw <= 160 && ih <= 160)) continue;
      var host = null;
      var hostArea = Infinity;
      for (j = 0; j < rects.length; j++) {
        var outer = rects[j];
        if (!outer.col || outer === inner) continue;
        if (inner.x0 < outer.x0 - 0.8 || inner.y0 < outer.y0 - 0.8) continue;
        if (inner.x1 > outer.x1 + 0.8 || inner.y1 > outer.y1 + 0.8) continue;
        var oa = (outer.x1 - outer.x0) * (outer.y1 - outer.y0);
        if ((iw * ih) >= oa - 1) continue;
        if (oa < hostArea) {
          hostArea = oa;
          host = outer;
        }
      }
      if (host) inner.col = host.col;
    }
    var hide = {};
    for (i = 0; i < rects.length; i++) {
      var rr = rects[i];
      if (!rr.col) continue;
      for (j = 0; j < rr.segs.length; j++) rr.segs[j].el.setAttribute('stroke', rr.col);
      for (j = 0; j < plates.length; j++) {
        var p = plates[j];
        var pcx = p.tx.x + p.tx.w / 2;
        var pcy = p.tx.y + p.tx.h / 2;
        if (pcx < rr.x0 - 2 || pcx > rr.x1 + 2 || pcy < rr.y0 - 2 || pcy > rr.y1 + 2) continue;
        hide[j] = 1;
      }
    }
    for (j = 0; j < plates.length; j++) {
      if (hide[j]) plates[j].el.style.borderColor = 'transparent';
    }
  }

  function syncFloor5OpenPlates() {
    var nodes = page.nodes || [];
    var lines = [];
    hostEl.querySelectorAll('[data-hx-line]').forEach(function (el) {
      if ((el.getAttribute('stroke') || '') === 'transparent') return;
      var n = nodes[Number(el.getAttribute('data-i'))];
      if (!n || n.t !== 'L' || n.net === 'frame' || n.net === 'bus') return;
      lines.push({
        el: el,
        n: n,
        ends: lineEndsOf(n),
        green: isGreenStroke(el.getAttribute('stroke'))
      });
    });
    hostEl.querySelectorAll('.hx-tx').forEach(function (el) {
      if (el.getAttribute('contenteditable') === 'true') return;
      var tx = nodes[Number(el.getAttribute('data-i'))];
      if (!tx || !tx.w || !tx.h) return;
      var touch = [];
      lines.forEach(function (ln) {
        if (hitsPlate(ln.ends, tx, 12)) touch.push(ln);
      });
      var openLine = touch.some(function (ln) { return ln.green; });
      if (openLine) {
        el.style.color = '#00ff00';
        el.style.borderColor = '#00ff00';
        el.classList.add('hx-plate-live');
        touch.forEach(function (ln) {
          if (ln.green) return;
          ln.el.setAttribute('stroke', '#00ff00');
          ln.green = true;
          ln.n.followKnife = 1;
        });
        return;
      }
      if (touch.some(function (ln) { return !ln.green; })) {
        el.style.color = '#ff0000';
        el.style.borderColor = '#ff0000';
        el.classList.add('hx-plate-live');
      }
    });
  }

  function paint() {
    paintTopoDevs();
    if (!page || !hostEl) return;
    var closedByKey = {};
    (page.nodes || []).forEach(function (n) {
      if (!n || (n.t !== 'I' && n.t !== 'S') || !n.k) return;
      closedByKey[n.k] = knifeIsClosed(n);
      closedByKey[String(n.k).replace(/_+$/, '')] = closedByKey[n.k];
    });
    hostEl.querySelectorAll('[data-hx-line]').forEach(function (el) {
      var n = page.nodes[Number(el.getAttribute('data-i'))];
      if (!n) return;
      var col = null;
      /* 刀下游出线（followKnife）整段跟闸：合红分绿；其余也只落红/绿，不留 CAD 蓝紫灰 */
      if (n.followKnife && n.k && Object.prototype.hasOwnProperty.call(closedByKey, n.k)) {
        col = closedByKey[n.k] ? '#ff0000' : '#00ff00';
      } else if (n.k && !n.fixc && Object.prototype.hasOwnProperty.call(closedByKey, n.k)) {
        col = closedByKey[n.k] ? '#ff0000' : '#00ff00';
      } else if (w.HXDraw && w.HXDraw.lineColor) {
        col = w.HXDraw.lineColor(page, n);
      }
      if (w.HXDraw && w.HXDraw.powerWireColor) {
        col = w.HXDraw.powerWireColor(col, true);
        if (col === '#ce0c18') col = '#ff0000';
      } else if (col && col !== '#00ff00' && col !== '#ff0000') {
        col = '#ff0000';
      }
      if (col) el.setAttribute('stroke', col);
    });
    hostEl.querySelectorAll('[data-hx-img]').forEach(function (el) {
      var n = page.nodes[Number(el.getAttribute('data-i'))];
      if (!n) return;
      var v = L().nodeValue(n);
      var show = !n.op || L().matchOp(v, n.op, n.val);
      if (v == null || v === '') show = n.st !== '1';
      el.style.display = show ? 'block' : 'none';
    });
    hostEl.querySelectorAll('[data-hx-knife]').forEach(function (el) {
      var n = page.nodes[Number(el.getAttribute('data-i'))];
      if (!n) return;
      var closed = knifeIsClosed(n);
      el.classList.toggle('is-closed', closed);
      el.classList.toggle('is-open', !closed);
      el.style.color = closed ? '#ff0000' : '#00ff00';
    });
    /* 出线侧铭牌跟闸；UPS 铭牌：进线或出线任一侧有电则红 */
    hostEl.querySelectorAll('.hx-tx').forEach(function (el) {
      var n = page.nodes[Number(el.getAttribute('data-i'))];
      if (!n) return;
      var col = null;
      if (n.upsPlate && n.w && n.h) {
        var inLive = false;
        var outLive = false;
        (page.nodes || []).forEach(function (ln) {
          if (!ln || ln.t !== 'L') return;
          var cx = ln.x + ln.w / 2;
          var cy = ln.y + 13;
          var half = ln.w / 2;
          var horiz = Math.abs((Number(ln.a) || 0) % 180) < 1;
          var ends = horiz
            ? [{ x: cx - half, y: cy }, { x: cx + half, y: cy }]
            : [{ x: cx, y: cy - half }, { x: cx, y: cy + half }];
          var live;
          if (ln.followKnife && ln.k && Object.prototype.hasOwnProperty.call(closedByKey, ln.k)) {
            live = !!closedByKey[ln.k];
          } else {
            live = true;
          }
          ends.forEach(function (e) {
            if (e.x < n.x - 8 || e.x > n.x + n.w + 8) return;
            if (e.y >= n.y - 10 && e.y <= n.y + 14) {
              if (live) inLive = true;
            }
            var bot = n.y + n.h;
            if (e.y >= bot - 14 && e.y <= bot + 10) {
              if (live) outLive = true;
            }
          });
        });
        col = (inLive || outLive) ? '#ff0000' : '#00ff00';
      } else if (n.followKnife && n.k && Object.prototype.hasOwnProperty.call(closedByKey, n.k)) {
        col = closedByKey[n.k] ? '#ff0000' : '#00ff00';
      }
      if (!col) return;
      el.style.color = col;
      el.style.borderColor = col;
      el.classList.add('hx-plate-live');
    });
    /* 五层：开闸绿线碰到的铭牌改绿，隔着铭牌的出线也改绿。母线保持红。 */
    if (page && /^配电室5/.test(String(page.name || ''))) syncFloor5OpenPlates();
    /* 铭牌框跟字同色，并藏掉文字上的第二圈 CSS 边框。 */
    syncPlateFrames();
    hostEl.querySelectorAll('[data-hx-uuid]').forEach(function (el) {
      var known = L().uuidValue(el.getAttribute('data-hx-uuid'));
      if (known) el.textContent = L().showValue(known);
    });
  }

  /** 设备名 → 通讯状态测点（能接到就接，接不上返回 null，不编假数据） */
  function topoCommPoint(devName) {
    var name = String(devName || '');
    if (!name || !L().pointOf) return null;
    var cands = [name + '_通讯状态'];
    if (name.indexOf('配电室') !== 0) cands.push('配电室' + name + '_通讯状态');
    else cands.push(name.replace(/^配电室/, '') + '_通讯状态');
    var i, pt;
    for (i = 0; i < cands.length; i++) {
      pt = L().pointOf(cands[i]);
      if (pt && pt.uuid) return pt;
    }
    return null;
  }

  function topoCommOk(value) {
    var s = String(value == null ? '' : value).trim();
    if (s === '') return null;
    if (s === '1' || s === '正常' || s === 'true' || s === 'OK' || s === 'ok') return true;
    if (s === '0' || s === '故障' || s === '中断' || s === '异常' || s === 'false') return false;
    var n = parseFloat(s);
    if (!isNaN(n)) return n !== 0;
    return null;
  }

  function topoDevState(devName, uuidHint) {
    var pt = topoCommPoint(devName);
    if (!pt || !pt.uuid) return { kind: 'unk', label: '未绑' };
    var known = (uuidHint && L().uuidValue(uuidHint)) || L().remembered(pt);
    var raw = known && known.value != null && known.value !== '' ? known.value : pt.value;
    if (raw == null || raw === '') return { kind: 'wait', label: '—' };
    var ok = topoCommOk(raw);
    if (ok === true) return { kind: 'ok', label: '正常' };
    if (ok === false) return { kind: 'bad', label: '故障' };
    return { kind: 'wait', label: '—' };
  }

  function topoShortName(devName) {
    return String(devName || '').replace(/^配电室/, '');
  }

  function applyTopoDevEl(el, st) {
    var short = topoShortName(el.getAttribute('data-dev') || '');
    el.classList.remove('is-ok', 'is-bad', 'is-unk', 'is-wait');
    el.classList.add('is-' + (st.kind || 'unk'));
    var flow = el.querySelector('.hx-topo-flow');
    if (flow) flow.hidden = st.kind !== 'ok';
    var lab = el.querySelector('.hx-topo-dev-lab');
    if (lab) lab.textContent = short + (st.label || '');
    el.title = st.kind === 'unk'
      ? (short + ' · 未找到「' + short + '_通讯状态」测点，占位显示')
      : (short + ' · 通讯' + (st.label || ''));
  }

  function labelOf(key) {
    var i, n;
    for (i = 0; i < page.nodes.length; i++) {
      n = page.nodes[i];
      if (n.t === 'T' && n.k === key && n.s) return n.s;
    }
    return key;
  }

  function powerDrillBtn(key) {
    if (!page || page.kind !== 'power' || !w.HXArch || !w.HXArch.suggestFromPower) return '';
    var code = (key || '').match(/(\d[AB]\d)/);
    if (!code) return '';
    var sug = w.HXArch.suggestFromPower('配电室' + code[1]);
    if (!sug || !sug.item) return '';
    return '<button type="button" class="hx-drill-go" data-hx-drill="1">下钻 · ' + esc(sug.item.name) + '</button>';
  }

  function hoverHtml(key, rows) {
    var body = rows.length ? rows.map(function (p) {
      var known = L().remembered(p) || p;
      return '<div class="hx-pt"><span>' + esc(p.name) + '</span><b data-hx-uuid="' + esc(p.uuid) + '">' + esc(L().showValue(known)) + '</b></div>';
    }).join('') : '<div class="hx-pt"><span>这一回路还没有测点</span></div>';
    return '<i class="hx-hover-caret" aria-hidden="true"></i><div class="hx-hover-body"><div class="hx-hover-title">' +
      esc(labelOf(key)) + '</div>' + body + powerDrillBtn(key) + '</div>';
  }

  function knifeOf(node) {
    return node && node.closest ? node.closest('[data-hx-knife]') : null;
  }

  function clearKnifeClass(root, cls, except) {
    if (!root) return;
    var nodes = root.querySelectorAll('[data-hx-knife].' + cls);
    var i;
    for (i = 0; i < nodes.length; i++) {
      if (nodes[i] !== except) nodes[i].classList.remove(cls);
    }
  }

  /* fixed 提示框的包含块是带 transform 的 .page，不是视口。把刀的视口矩形换到这个包含块里。 */
  function localIn(origin, el) {
    var r = el.getBoundingClientRect();
    if (!origin) {
      return { left: r.left, right: r.right, top: r.top, bottom: r.bottom, cx: (r.left + r.right) / 2, cy: (r.top + r.bottom) / 2, viewW: window.innerWidth, viewH: window.innerHeight };
    }
    var o = origin.getBoundingClientRect();
    var sx = o.width ? origin.offsetWidth / o.width : 1;
    var sy = o.height ? origin.offsetHeight / o.height : 1;
    var left = (r.left - o.left) * sx;
    var right = (r.right - o.left) * sx;
    var top = (r.top - o.top) * sy;
    var bottom = (r.bottom - o.top) * sy;
    return { left: left, right: right, top: top, bottom: bottom, cx: (left + right) / 2, cy: (top + bottom) / 2, viewW: origin.offsetWidth, viewH: origin.offsetHeight };
  }

  /* 提示框贴在刀的一侧，三角指向刀心。靠近屏幕边缘时翻到另一侧，三角仍指向这把刀。 */
  function placeHoverOnKnife(box, el) {
    var k = localIn(box.offsetParent, el);
    var bw = box.offsetWidth || 420;
    var bh = box.offsetHeight || 160;
    var gap = 18;
    var left = k.right + gap;
    var side = 'left';
    if (left + bw > k.viewW - 8) {
      left = k.left - gap - bw;
      side = 'right';
    }
    if (left < 8) {
      left = 8;
      side = 'left';
    }
    var top = k.cy - bh / 2;
    if (top < 8) top = 8;
    if (top + bh > k.viewH - 8) top = Math.max(8, k.viewH - bh - 8);
    box.style.left = left + 'px';
    box.style.top = top + 'px';
    var caret = k.cy - top - 8;
    if (caret < 12) caret = 12;
    if (caret > bh - 28) caret = Math.max(12, bh - 28);
    box.style.setProperty('--hx-caret-y', caret + 'px');
    box.classList.toggle('is-caret-left', side === 'left');
    box.classList.toggle('is-caret-right', side === 'right');
  }

  function placeHover(box, ev) {
    if (!box) return;
    /* 一律贴在目标旁，不盖住目标，否则点不到别的设备 */
    if (hoverAnchor && hoverAnchor.isConnected) {
      placeHoverOnKnife(box, hoverAnchor);
      return;
    }
    box.classList.remove('is-caret-left', 'is-caret-right');
    if (!ev) return;
    var origin = box.offsetParent;
    var r = ev.target && ev.target.getBoundingClientRect ? ev.target.getBoundingClientRect() : null;
    var bw = box.offsetWidth || 420;
    var bh = box.offsetHeight || 160;
    var gap = 16;
    var left;
    var top;
    if (r && r.width) {
      var loc = localIn(origin, ev.target.closest ? (ev.target.closest('[data-hx-key]') || ev.target) : ev.target);
      left = loc.right + gap;
      if (left + bw > loc.viewW - 8) left = loc.left - gap - bw;
      if (left < 8) left = 8;
      top = loc.cy - bh / 2;
      if (top < 8) top = 8;
      if (top + bh > loc.viewH - 8) top = Math.max(8, loc.viewH - bh - 8);
    } else {
      left = ev.clientX + gap;
      top = ev.clientY + gap;
      if (left + bw > window.innerWidth) left = Math.max(8, ev.clientX - bw - gap);
      if (top + bh > window.innerHeight) top = Math.max(8, ev.clientY - bh - gap);
    }
    box.style.left = left + 'px';
    box.style.top = top + 'px';
  }

  function openHover(key, ev, anchor) {
    hoverAnchor = (anchor && anchor.getBoundingClientRect) ? anchor : null;
    hoverRows = L().pointsFor(key, devOf(key));
    L().allow(pageRows.concat(hoverRows));
    var box = document.getElementById('hxHover');
    if (box) {
      box.innerHTML = hoverHtml(key, hoverRows);
      box.hidden = false;
      placeHover(box, ev);
    }
    L().fetch(hoverRows).then(function () {
      if (box && !box.hidden && hotKey === key) {
        box.innerHTML = hoverHtml(key, hoverRows);
        placeHover(box, ev);
      }
    });
  }

  function closeHover() {
    if (pinned) return;
    hotKey = '';
    hoverRows = [];
    hoverAnchor = null;
    var box = document.getElementById('hxHover');
    if (box) {
      box.hidden = true;
      box.classList.remove('is-caret-left', 'is-caret-right');
    }
    L().allow(pageRows);
  }

  function editOn() {
    return !!(w.HXEdit && w.HXEdit.blocksStage && w.HXEdit.blocksStage());
  }

  function bindHover(host) {
    host.addEventListener('mouseover', function (ev) {
      if (editOn()) return;
      var knife = knifeOf(ev.target);
      if (knife) {
        clearKnifeClass(host, 'hx-knife-hot', knife);
        knife.classList.add('hx-knife-hot');
        if (pinned) return;
        var key = knife.getAttribute('data-hx-key');
        if (!key || key === hotKey) {
          if (key && hoverAnchor !== knife) {
            hoverAnchor = knife;
            var stay = document.getElementById('hxHover');
            if (stay && !stay.hidden) placeHover(stay, ev);
          }
          return;
        }
        hotKey = key;
        openHover(key, ev, knife);
        return;
      }
      if (pinned) return;
      var el = ev.target.closest ? ev.target.closest('[data-hx-key]') : null;
      if (!el) return;
      var key2 = el.getAttribute('data-hx-key');
      if (!key2 || key2 === hotKey) return;
      hotKey = key2;
      openHover(key2, ev, null);
    });
    host.addEventListener('mousemove', function (ev) {
      var box = document.getElementById('hxHover');
      if (box && !box.hidden && !pinned) placeHover(box, ev);
    });
    host.addEventListener('mouseout', function (ev) {
      var knife = knifeOf(ev.target);
      if (knife) {
        var inside = ev.relatedTarget;
        if (!inside || !knife.contains(inside)) knife.classList.remove('hx-knife-hot');
      }
      var to = ev.relatedTarget;
      if (to && to.closest && (to.closest('[data-hx-key]') || to.closest('#hxHover'))) return;
      closeHover();
    });
    host.addEventListener('click', function (ev) {
      if (editOn()) return;
      var knife = knifeOf(ev.target);
      if (knife) {
        clearKnifeClass(host, 'hx-knife-pin', knife);
        knife.classList.add('hx-knife-pin');
        pinned = true;
        hotKey = knife.getAttribute('data-hx-key');
        openHover(hotKey, ev, knife);
        return;
      }
      var el = ev.target.closest ? ev.target.closest('[data-hx-key]') : null;
      if (!el) {
        pinned = false;
        hotKey = '';
        hoverRows = [];
        hoverAnchor = null;
        clearKnifeClass(host, 'hx-knife-pin', null);
        var box = document.getElementById('hxHover');
        if (box) {
          box.hidden = true;
          box.classList.remove('is-caret-left', 'is-caret-right');
        }
        L().allow(pageRows);
        return;
      }
      clearKnifeClass(host, 'hx-knife-pin', null);
      pinned = true;
      hotKey = el.getAttribute('data-hx-key');
      openHover(hotKey, ev, null);
    });
  }

  var stage = { user: false, tx: 0, ty: 0, sc: 1, baseS: 1, baseDx: 0, marquee: false };

  function stageBoard() { return hostEl && hostEl.querySelector('.hx-board'); }
  function stageFit() { return hostEl && hostEl.querySelector('.hx-fit'); }

  /* #screen 有 transform:scale。getBoundingClientRect 是缩放后的视觉盒，
     left/top 要用 .hx-fit 自己的 CSS 像素，否则选框会偏离光标、像从中间起笔。 */
  function fitLocal(fit, clientX, clientY) {
    var r = fit.getBoundingClientRect();
    var kx = r.width ? fit.clientWidth / r.width : 1;
    var ky = r.height ? fit.clientHeight / r.height : 1;
    return { x: (clientX - r.left) * kx, y: (clientY - r.top) * ky };
  }

  /* #screen 自己还有一层 scale（core.js）。画布平移要再乘上这一层，才落到物理像素。 */
  function outerScreenScale() {
    var s = w.App && Number(w.App.scale);
    if (s > 0 && isFinite(s)) return s;
    var el = document.getElementById('screen');
    if (el && el.style && el.style.transform) {
      var m = /scale\(\s*([0-9.]+)\s*\)/.exec(el.style.transform);
      if (m) {
        s = parseFloat(m[1]);
        if (s > 0 && isFinite(s)) return s;
      }
    }
    return 1;
  }

  /* 平移落到物理像素上。Windows 上半像素的 translate 会把细线整段裁掉。 */
  function snapPx(v) {
    var dpr = w.devicePixelRatio || 1;
    var unit = dpr * outerScreenScale();
    if (!(unit > 0)) unit = 1;
    return Math.round(v * unit) / unit;
  }

  /* 缩进小窗口后，线在屏幕上至少 1.35px。画布缩放和整屏缩放都要算进去。
     否则华硕这类屏上 Chrome 会把线丢掉，苹果视网膜还能看见。 */
  function thickenStrokes(root, selector, screenScale) {
    if (!root || !(screenScale > 0)) return;
    var outer = outerScreenScale();
    var min = 1.35 / (screenScale * outer);
    root.querySelectorAll(selector).forEach(function (el) {
      var orig = el.getAttribute('data-hx-sw');
      if (orig == null) {
        var raw = (el.style && el.style.strokeWidth) || el.getAttribute('stroke-width') || '1.5';
        orig = String(parseFloat(raw) || 1.5);
        el.setAttribute('data-hx-sw', orig);
      }
      var base = parseFloat(orig) || 1.5;
      var sw = Math.max(base, min);
      if (el.style && el.style.strokeWidth) el.style.strokeWidth = sw.toFixed(3);
      else el.setAttribute('stroke-width', sw.toFixed(3));
    });
  }

  function applyStage(tx, ty, sc) {
    var board = stageBoard();
    if (!board) return;
    tx = snapPx(tx);
    ty = snapPx(ty);
    stage.tx = tx;
    stage.ty = ty;
    stage.sc = sc;
    board.style.transformOrigin = '0 0';
    board.style.transform = 'translate(' + tx + 'px,' + ty + 'px) scale(' + sc + ')';
    thickenStrokes(board, '[data-hx-line]', sc);
  }

  function zoomAt(clientX, clientY, next) {
    var fit = stageFit();
    if (!fit || !page || page.kind !== 'power') return;
    var rect = fit.getBoundingClientRect();
    var px = clientX - rect.left;
    var py = clientY - rect.top;
    var sc = stage.sc || 1;
    next = Math.max(0.05, Math.min(8, next));
    var bx = (px - stage.tx) / sc;
    var by = (py - stage.ty) / sc;
    stage.user = true;
    applyStage(px - bx * next, py - by * next, next);
  }

  function zoomCenter(factor) {
    var fit = stageFit();
    if (!fit) return;
    var rect = fit.getBoundingClientRect();
    zoomAt(rect.left + rect.width / 2, rect.top + rect.height / 2, (stage.sc || 1) * factor);
  }

  function fitWindow() {
    stage.user = false;
    stage.marquee = false;
    var btn = hostEl && hostEl.querySelector('[data-hx-zoom="area"]');
    if (btn) btn.classList.remove('is-on');
    var fitEl = stageFit();
    if (fitEl) fitEl.classList.remove('is-marquee');
    fit();
  }

  function zoom100() {
    var fit = stageFit();
    var board = stageBoard();
    if (!fit || !board) return;
    stage.user = true;
    applyStage((fit.clientWidth - board.offsetWidth) / 2, (fit.clientHeight - board.offsetHeight) / 2, 1);
  }

  /* 刀、线，以及贴着这团刀线的符号和标注。远离刀群的超大坐标不进包围盒。 */
  function powerContentBox(pg) {
    var nodes = (pg && pg.nodes) || [];
    var kminX = Infinity;
    var kminY = Infinity;
    var kmaxX = -Infinity;
    var kmaxY = -Infinity;
    var kn = 0;
    nodes.forEach(function (n) {
      if (n.t !== 'I' && n.t !== 'S') return;
      var x = Number(n.x);
      var y = Number(n.y);
      var w = Number(n.w) || 60;
      var h = Number(n.h) || 98;
      if (!isFinite(x) || !isFinite(y)) return;
      kn += 1;
      if (x < kminX) kminX = x;
      if (y < kminY) kminY = y;
      if (x + w > kmaxX) kmaxX = x + w;
      if (y + h > kmaxY) kmaxY = y + h;
    });
    if (!kn || !(kmaxX > kminX) || !(kmaxY > kminY)) return null;
    var span = Math.max(kmaxX - kminX, kmaxY - kminY);
    var near = Math.max(480, span * 0.12);
    var box = { minX: kminX, minY: kminY, maxX: kmaxX, maxY: kmaxY };
    function take(x0, y0, x1, y1) {
      if (!isFinite(x0) || !isFinite(y0) || !isFinite(x1) || !isFinite(y1)) return;
      var cx = (x0 + x1) / 2;
      var cy = (y0 + y1) / 2;
      if (cx < kminX - near || cx > kmaxX + near || cy < kminY - near || cy > kmaxY + near) return;
      if (x0 < box.minX) box.minX = x0;
      if (y0 < box.minY) box.minY = y0;
      if (x1 > box.maxX) box.maxX = x1;
      if (y1 > box.maxY) box.maxY = y1;
    }
    nodes.forEach(function (n) {
      if (n.t === 'L') {
        var e = lineEnds(n);
        var x0 = Math.min(e.x1, e.x2);
        var x1 = Math.max(e.x1, e.x2);
        var y0 = Math.min(e.y1, e.y2);
        var y1 = Math.max(e.y1, e.y2);
        var thick = Math.max(4, Number(n.sw) || 2) / 2;
        if (isHorizLine(n)) {
          y0 -= thick;
          y1 += thick;
        } else if (isVertLine(n)) {
          x0 -= thick;
          x1 += thick;
        } else {
          x0 -= thick;
          x1 += thick;
          y0 -= thick;
          y1 += thick;
        }
        take(x0, y0, x1, y1);
        return;
      }
      if (n.t === 'B' || n.t === 'T') {
        var x = Number(n.x);
        var y = Number(n.y);
        take(x, y, x + (Number(n.w) || 0), y + (Number(n.h) || 0));
      }
    });
    if (!(box.maxX > box.minX) || !(box.maxY > box.minY)) return null;
    return box;
  }

  /* 5A4：标题和线路的外接矩形。空画布不算进去，否则图会沉在底下。 */
  function contentBoxAll(pg) {
    var nodes = (pg && pg.nodes) || [];
    var box = null;
    function take(x0, y0, x1, y1) {
      if (!isFinite(x0) || !isFinite(y0) || !isFinite(x1) || !isFinite(y1)) return;
      var loX = Math.min(x0, x1);
      var loY = Math.min(y0, y1);
      var hiX = Math.max(x0, x1);
      var hiY = Math.max(y0, y1);
      if (!box) box = { minX: loX, minY: loY, maxX: hiX, maxY: hiY };
      else {
        if (loX < box.minX) box.minX = loX;
        if (loY < box.minY) box.minY = loY;
        if (hiX > box.maxX) box.maxX = hiX;
        if (hiY > box.maxY) box.maxY = hiY;
      }
    }
    nodes.forEach(function (n) {
      if (!n) return;
      if (n.t === 'L') {
        var e = lineEnds(n);
        var x0 = Math.min(e.x1, e.x2);
        var x1 = Math.max(e.x1, e.x2);
        var y0 = Math.min(e.y1, e.y2);
        var y1 = Math.max(e.y1, e.y2);
        var thick = Math.max(4, Number(n.sw) || 2) / 2;
        if (isHorizLine(n)) { y0 -= thick; y1 += thick; }
        else if (isVertLine(n)) { x0 -= thick; x1 += thick; }
        else { x0 -= thick; x1 += thick; y0 -= thick; y1 += thick; }
        take(x0, y0, x1, y1);
        return;
      }
      if (n.t === 'I' || n.t === 'S' || n.t === 'T' || n.t === 'B') {
        var x = Number(n.x);
        var y = Number(n.y);
        take(x, y, x + (Number(n.w) || 0), y + (Number(n.h) || 0));
      }
    });
    if (!box || !(box.maxX > box.minX) || !(box.maxY > box.minY)) return null;
    return box;
  }

  function fit() {
    var main = hostEl && hostEl.querySelector('.hx-main');
    var fitEl = stageFit();
    var board = stageBoard();
    if (!main || !fitEl || !board) return;
    var cw = main.clientWidth;
    var ch = main.clientHeight;
    if (cw < 40 || ch < 40) return;
    fitEl.style.width = cw + 'px';
    fitEl.style.height = ch + 'px';
    fitEl.style.transform = 'none';
    /* 整张画布放进 #hxFit，四周留边。扁长室不再改按高度铺满，否则宽度溢出被裁。 */
    var bw = board.offsetWidth;
    var bh = board.offsetHeight;
    var ox = 0;
    var oy = 0;
    if (page && page.kind === 'power' && roomCodeOf(page) === '5A4') {
      var box = contentBoxAll(page);
      if (box) {
        var x0 = Number(page.x0) || 0;
        bw = box.maxX - box.minX;
        bh = box.maxY - box.minY;
        ox = box.minX - x0;
        oy = box.minY;
      }
    }
    if (!(bw > 1) || !(bh > 1)) return;
    var pad = Math.max(16, Math.round(Math.min(cw, ch) * 0.03));
    var availW = Math.max(1, cw - pad * 2);
    var availH = Math.max(1, ch - pad * 2);
    var s = Math.min(availW / bw, availH / bh);
    if (!isFinite(s) || s <= 0) s = 0.4;
    s = Math.max(0.05, Math.min(8, s));
    var contentW = bw * s;
    var contentH = bh * s;
    var dx = (cw - contentW) / 2 - ox * s;
    var dy = (ch - contentH) / 2 - oy * s;
    stage.baseS = s;
    stage.baseDx = dx;
    if (stage.user && page && page.kind === 'power') return;
    applyStage(dx, dy, s);
  }

  function bindStage(main) {
    if (main._hxStage) return;
    main._hxStage = 1;
    var drag = null;
    function blank(el) {
      if (!el || !el.closest) return false;
      if (el.closest('.hx-zoom-bar, [data-hx-key], .hx-sym')) return false;
      return !!el.closest('.hx-fit');
    }
    /* 编辑态仍要能缩放看清：滚轮/双击缩放、局部框选可用；空白平移仍让给编辑手势。 */
    main.addEventListener('wheel', function (ev) {
      if (!page || page.kind !== 'power' || !blank(ev.target) && !(ev.target.closest && ev.target.closest('.hx-fit'))) return;
      if (ev.target.closest && ev.target.closest('.hx-zoom-bar')) return;
      ev.preventDefault();
      var dy = ev.deltaY;
      if (ev.deltaMode === 1) dy *= 16;
      zoomAt(ev.clientX, ev.clientY, (stage.sc || 1) * (dy > 0 ? 0.9 : 1.12));
    }, { passive: false });
    main.addEventListener('dblclick', function (ev) {
      if (!page || page.kind !== 'power') return;
      if (ev.target.closest && ev.target.closest('.hx-cad-pane, .hx-zoom-bar')) return;
      /* 编辑改字优先；其它位置仍按原手势缩放 */
      if (editOn() && ev.target.closest && ev.target.closest('.hx-tx')) return;
      zoomAt(ev.clientX, ev.clientY, (stage.sc || 1) * (ev.shiftKey ? 1 / 1.35 : 1.35));
    });
    main.addEventListener('pointerdown', function (ev) {
      if (!page || page.kind !== 'power' || ev.button !== 0) return;
      if (ev.target.closest && ev.target.closest('.hx-zoom-bar')) return;
      var fit = stageFit();
      if (!fit) return;
      if (stage.marquee) {
        if (!blank(ev.target) && !(ev.target.closest && ev.target.closest('.hx-fit'))) return;
        var p0 = fitLocal(fit, ev.clientX, ev.clientY);
        drag = { mode: 'area', x: p0.x, y: p0.y, px: ev.clientX, py: ev.clientY };
        var box = document.createElement('div');
        box.className = 'hx-marquee';
        box.style.left = drag.x + 'px';
        box.style.top = drag.y + 'px';
        fit.appendChild(box);
        drag.box = box;
        if (main.setPointerCapture) main.setPointerCapture(ev.pointerId);
        return;
      }
      if (editOn()) return;
      if (!blank(ev.target)) return;
      drag = { mode: 'pan', x: ev.clientX, y: ev.clientY, tx: stage.tx, ty: stage.ty, moved: false };
      if (main.setPointerCapture) main.setPointerCapture(ev.pointerId);
    });
    main.addEventListener('pointermove', function (ev) {
      if (!drag) return;
      if (drag.mode === 'pan') {
        var dx = ev.clientX - drag.x;
        var dy = ev.clientY - drag.y;
        if (!drag.moved && Math.hypot(dx, dy) < 4) return;
        drag.moved = true;
        stage.user = true;
        applyStage(drag.tx + dx, drag.ty + dy, stage.sc);
        return;
      }
      var fit = stageFit();
      if (!fit || !drag.box) return;
      var p = fitLocal(fit, ev.clientX, ev.clientY);
      var x = p.x;
      var y = p.y;
      var x0 = Math.min(drag.x, x);
      var y0 = Math.min(drag.y, y);
      drag.box.style.left = x0 + 'px';
      drag.box.style.top = y0 + 'px';
      drag.box.style.width = Math.abs(x - drag.x) + 'px';
      drag.box.style.height = Math.abs(y - drag.y) + 'px';
    });
    function endDrag(ev) {
      if (!drag) return;
      var d = drag;
      drag = null;
      if (d.mode === 'area' && d.box) {
        var w = parseFloat(d.box.style.width) || 0;
        var h = parseFloat(d.box.style.height) || 0;
        var x0 = parseFloat(d.box.style.left) || 0;
        var y0 = parseFloat(d.box.style.top) || 0;
        if (d.box.parentNode) d.box.parentNode.removeChild(d.box);
        if (w > 12 && h > 12) {
          var fit = stageFit();
          var cw = fit.clientWidth;
          var ch = fit.clientHeight;
          var sc0 = stage.sc || 1;
          var b0x = (x0 - stage.tx) / sc0;
          var b0y = (y0 - stage.ty) / sc0;
          var bw = w / sc0;
          var bh = h / sc0;
          var sc = Math.min(cw / bw, ch / bh) * 0.94;
          stage.user = true;
          applyStage((cw - bw * sc) / 2 - b0x * sc, (ch - bh * sc) / 2 - b0y * sc, sc);
        }
        stage.marquee = false;
        var fitEl = stageFit();
        if (fitEl) fitEl.classList.remove('is-marquee');
        var btn = main.querySelector('[data-hx-zoom="area"]');
        if (btn) btn.classList.remove('is-on');
      }
    }
    main.addEventListener('pointerup', endDrag);
    main.addEventListener('pointercancel', endDrag);
    main.addEventListener('click', function (ev) {
      var btn = ev.target.closest ? ev.target.closest('[data-hx-zoom]') : null;
      if (!btn || !page || page.kind !== 'power') return;
      ev.preventDefault();
      ev.stopPropagation();
      var act = btn.getAttribute('data-hx-zoom');
      if (act === 'in') zoomCenter(1.2);
      else if (act === 'out') zoomCenter(1 / 1.2);
      else if (act === 'fit') fitWindow();
      else if (act === '100') zoom100();
      else if (act === 'area') {
        stage.marquee = !stage.marquee;
        btn.classList.toggle('is-on', stage.marquee);
        var fit = stageFit();
        if (fit) fit.classList.toggle('is-marquee', stage.marquee);
      }
    });
  }

  function guessKey(pg, text) {
    var s = String(text || '').replace(/\s+/g, '');
    if (!s || !L().hasCircuit) return '';
    if (pg.kind === 'column' && pg.cab) {
      var k = '';
      if (s.indexOf('进线') >= 0) k = pg.cab + '_主路';
      var m = s.match(/^(\d+)([ABC])$/);
      if (m) k = pg.cab + '_支路' + m[1] + m[2];
      if (k && L().hasCircuit(k)) return k;
      return '';
    }
    if (pg.kind !== 'power') return '';
    var norm = s.replace(/-/g, '_');
    var cands = [];
    var rm = norm.match(/^(\d[AB]\d)_(.+)$/);
    if (rm) cands.push('配电室' + rm[1] + '_' + rm[2]);
    cands.push('配电室' + norm);
    for (var i = 0; i < cands.length; i++) {
      if (L().hasCircuit(cands[i])) return cands[i];
    }
    return '';
  }

  function lineEnds(n) {
    var rad = (Number(n.a) || 0) * Math.PI / 180;
    var cos = Math.cos(rad);
    var sin = Math.sin(rad);
    var cx = Number(n.x) + Number(n.w || 0) / 2;
    var cy = Number(n.y) + Number(n.h || 0) / 2;
    var hw = Number(n.w || 0) / 2;
    return {
      cx: cx,
      cy: cy,
      y1: cy - hw * sin,
      y2: cy + hw * sin,
      x1: cx - hw * cos,
      x2: cx + hw * cos
    };
  }

  function isVertLine(n) {
    var a = Math.abs((Number(n.a) || 0) % 180);
    return a > 75 && a < 105;
  }

  function isHorizLine(n) {
    var a = Math.abs((Number(n.a) || 0) % 180);
    return a <= 5 || a >= 175;
  }

  /** 并联分组基键：T1/T2 馈线、U1/U2 Yxx/Xxx 支路等。 */
  function feederBaseKey(k) {
    var s = String(k || '');
    var m = s.match(/^(配电室\d[AB]\d_T[12]_\d+)/);
    if (m) return m[1];
    m = s.match(/^(配电室\d[AB]\d_U[12]_Y\d+)/);
    if (m) return m[1];
    m = s.match(/^(配电室\d[AB]\d_U[12]_X[A-Z0-9]+)/);
    if (m) return m[1];
    m = s.match(/^(配电室\d[AB]\d)_445$/);
    if (m) return m[1] + '_445';
    return '';
  }

  function feederSuffix(k) {
    var s = String(k || '');
    var m = s.match(/_(\d+)_?$/);
    if (m) return parseInt(m[1], 10);
    if (/_$/.test(s)) return 1;
    return 0;
  }

  function isMainFeederKey(k) {
    return /(T[12]_4|_401|_402|_411|_412|_413|_414|_415|_416|_421|_422|_423|_445)/.test(String(k || ''));
  }

  function isParallelKnifeKey(k, pageName) {
    var s = String(k || '');
    if (!/^配电室\d[AB]\d_/.test(s)) return false;
    /* 5A1 走 hx-5a1-cad-layout.json 显式竖刀，禁止旧横置并联重排 */
    if (pageName === '配电室5A1') return false;
    if (isMainFeederKey(s)) return true;
    if (/配电室\d[AB]\d_U[12]_Y\d+/.test(s)) return true;
    if (/配电室\d[AB]\d_U[12]_X/.test(s)) return true;
    return false;
  }

  /** scene 合闸帧常导出 88×100，与标准竖刀 60×98 混用会导致 401/411 视觉大小不一。 */
  var KNIFE_W = 60;
  var KNIFE_H = 98;
  var PARALLEL_Y_PITCH_MIN = 56;
  var PARALLEL_COL_GAP = 28;
  var PARALLEL_BUS_OFFSET = 32;

  var BLADE_VX = 18;
  var BLADE_VW = 48;
  var BLADE_VH = 124;
  /* 与 hx.css --hx-blade-long 相同：0.92 * 60。竖刀只在 bladeFit 的五层单室用这个长轴。 */
  var BLADE_LONG = 0.92 * KNIFE_W;

  function symDisplayBox(n) {
    if (knifeHorizontal(n)) return { w: KNIFE_H, h: KNIFE_W };
    return { w: KNIFE_W, h: KNIFE_H };
  }

  /** 与 hx.css 中 .hx-blade 缩放一致，供接线与 SVG 铰接圆 (18,0)/(18,114) 对齐。 */
  function knifeBladeLayout(n) {
    var box = symDisplayBox(n);
    var decor = !!n.decor;
    var horiz = knifeHorizontal(n);
    var s;
    var rw;
    var rh;
    var ox;
    var oy;
    if (horiz && decor) {
      rh = box.h * 0.78;
      s = rh / BLADE_VH;
      rw = BLADE_VW * s;
      ox = (box.w - rw) / 2;
      oy = (box.h - rh) / 2;
    } else if (horiz) {
      rh = BLADE_LONG;
      s = rh / BLADE_VH;
      rw = BLADE_VW * s;
      ox = (box.w - rw) / 2;
      oy = (box.h - rh) / 2;
    } else if (n.bladeFit) {
      s = BLADE_LONG / BLADE_VH;
      rw = BLADE_VW * s;
      rh = BLADE_LONG;
      ox = (box.w - rw) / 2;
      oy = 0;
    } else if (decor) {
      rw = box.w * 0.84;
      s = rw / BLADE_VW;
      rh = BLADE_VH * s;
      ox = (box.w - rw) / 2;
      oy = 0;
    } else {
      s = Math.min(box.w / BLADE_VW, box.h / BLADE_VH);
      rw = BLADE_VW * s;
      rh = BLADE_VH * s;
      ox = (box.w - rw) / 2;
      oy = (box.h - rh) / 2;
    }
    return { box: box, s: s, ox: ox, oy: oy, horiz: horiz };
  }

  function bladePointInBox(n, vx, vy) {
    var lay = knifeBladeLayout(n);
    var px = lay.ox + vx * lay.s;
    var py = lay.oy + vy * lay.s;
    if (lay.horiz) {
      var cx = lay.box.w / 2;
      var cy = lay.box.h / 2;
      var dx = px - cx;
      var dy = py - cy;
      px = cx + dy;
      py = cy - dx;
    }
    return { x: Number(n.x) + px, y: Number(n.y) + py };
  }

  function knifeMidY(n) {
    return bladePointInBox(n, BLADE_VX, 44).y;
  }

  /** 统一 power 页刀闸为 60×98；归一化时锚定 SVG 进线圆心，避免红线与铰接点错位。 */
  function normalizePowerKnifeSize(n) {
    if (n.t !== 'I' && n.t !== 'S') return;
    var oldW = Number(n.w) || KNIFE_W;
    var oldH = Number(n.h) || KNIFE_H;
    if (oldW === KNIFE_W && oldH === KNIFE_H) return;
    var anchor = knifeHorizontal(n)
      ? bladePointInBox(n, BLADE_VX, 44)
      : bladePointInBox(n, BLADE_VX, 0);
    n.w = KNIFE_W;
    n.h = KNIFE_H;
    var after = knifeHorizontal(n)
      ? bladePointInBox(n, BLADE_VX, 44)
      : bladePointInBox(n, BLADE_VX, 0);
    n.x += anchor.x - after.x;
    n.y += anchor.y - after.y;
  }

  function normalizePowerKnifeSizes(nodes) {
    (nodes || []).forEach(normalizePowerKnifeSize);
  }

  function findTopBusY(nodes) {
    var busY = null;
    var busSpan = 0;
    nodes.forEach(function (n) {
      if (n.t !== 'L' || !isHorizLine(n)) return;
      var y = Number(n.y) + Number(n.h || 0) / 2;
      var w = Math.abs(Number(n.w) || 0);
      if (y >= 72 && y <= 132 && w > 400 && w > busSpan) {
        busY = y;
        busSpan = w;
      }
    });
    return busY;
  }

  /** 顶母线 scene 横线段：cy 与 x0~x1（与 pushHorizLine / lineEnds 同一套中心线坐标）。 */
  function mainBusHorizAt(nodes, busY, xProbe) {
    var segs = listTopBusHoriz(nodes, busY);
    if (!segs.length) return null;
    var pick = pickBusSegForX(segs, xProbe);
    return pick || segs[0];
  }

  /** 同一顶母线高度上的全部长横母（T1 左段 / T2 右段等）。 */
  function listTopBusHoriz(nodes, busY) {
    var segs = [];
    nodes.forEach(function (n) {
      if (n.t !== 'L' || !isHorizLine(n)) return;
      var cy = Number(n.y) + Number(n.h || 0) / 2;
      if (Math.abs(cy - busY) > 10) return;
      var w = Math.abs(Number(n.w) || 0);
      if (w < 200) return;
      var x0 = Number(n.x);
      segs.push({ cy: cy, x0: x0, x1: x0 + w, w: w });
    });
    segs.sort(function (a, b) { return a.x0 - b.x0; });
    return segs;
  }

  function pickBusSegForX(segs, x) {
    var i;
    for (i = 0; i < segs.length; i++) {
      var s = segs[i];
      if (x >= s.x0 - 48 && x <= s.x1 + 48) return s;
    }
    var best = null;
    var bestD = 1e9;
    for (i = 0; i < segs.length; i++) {
      var seg = segs[i];
      var d = x < seg.x0 ? seg.x0 - x : (x > seg.x1 ? x - seg.x1 : 0);
      if (d < bestD) {
        bestD = d;
        best = seg;
      }
    }
    return best;
  }

  function pickBackboneJoinX(tap, backboneCols, seg) {
    var best = null;
    var i;
    for (i = 0; i < backboneCols.length; i++) {
      var b = backboneCols[i];
      if (seg && (b.cx < seg.x0 - 80 || b.cx > seg.x1 + 80)) continue;
      if (b.cx + 6 < tap.x) continue;
      if (!best || b.cx < best.cx) best = b;
    }
    if (best) return best.cx;
    if (seg && tap.x < seg.x0 - 4) return seg.x0;
    if (seg && tap.x > seg.x1 + 4) return seg.x1;
    return tap.x;
  }

  function snapshotKnifeNodes(nodes) {
    var snap = {};
    nodes.forEach(function (n, i) {
      if (n.t !== 'I' && n.t !== 'S') return;
      var k = String(n.k || '');
      var id = k || ('@' + i);
      snap[id] = {
        x: n.x, y: n.y, w: n.w, h: n.h, z: n.z,
        orient: n.orient, rot: n.rot, left: n.left, cx: n.cx
      };
    });
    return snap;
  }

  function restoreNonGroupedKnives(nodes, groupedKeys, snap) {
    nodes.forEach(function (n, i) {
      if (n.t !== 'I' && n.t !== 'S') return;
      var k = String(n.k || '');
      if (groupedKeys[k]) return;
      var id = k || ('@' + i);
      var o = snap[id];
      if (!o) {
        if (n.orient) delete n.orient;
        if (n.rot) delete n.rot;
        if (n.left) delete n.left;
        if (n.cx) delete n.cx;
        return;
      }
      n.x = o.x;
      n.y = o.y;
      if (o.w != null) n.w = o.w;
      if (o.h != null) n.h = o.h;
      if (o.z != null) n.z = o.z;
      if (o.orient) n.orient = o.orient;
      else if (n.orient) delete n.orient;
      if (o.rot != null) n.rot = o.rot;
      else if (n.rot) delete n.rot;
      if (o.left != null) n.left = o.left;
      else if (n.left) delete n.left;
      if (o.cx != null) n.cx = o.cx;
      else if (n.cx) delete n.cx;
    });
  }

  /** 刀列上方就近横母线（下半区 UPS/支路用）。 */
  function localFeedBusY(nodes, yMin, mainBusY) {
    if (yMin < mainBusY + 280) return mainBusY;
    var best = null;
    var bestW = 0;
    nodes.forEach(function (n) {
      if (n.t !== 'L' || !isHorizLine(n)) return;
      var cy = Number(n.y) + Number(n.h || 0) / 2;
      var w = Math.abs(Number(n.w) || 0);
      if (w < 80 || cy > yMin + 40 || cy < yMin - 120) return;
      if (w > bestW) {
        bestW = w;
        best = cy;
      }
    });
    return best != null ? best : mainBusY;
  }

  function pickKnifeSlot(slot, n, cx, y, suf) {
    var prev = slot[suf];
    if (!prev || n.st === '1' || (prev.n.st !== '1' && y < prev.y)) {
      slot[suf] = {
        n: n, cx: cx, y: y, h: Number(n.h || 98), suf: suf, left: Number(n.x), base: ''
      };
    }
  }

  function collectParallelKnifeSlots(nodes, busY, pageName) {
    var byBase = {};
    nodes.forEach(function (n) {
      if (n.t !== 'I' && n.t !== 'S') return;
      var k = String(n.k || '');
      if (!isParallelKnifeKey(k, pageName)) return;
      var base = feederBaseKey(k);
      if (!base) return;
      /* 总图 CAD：T1-401/T2-402 为进线单刀，不参与并联组排版 */
      if (/_T[12]_40[12]$/.test(base)) return;
      var cx = Number(n.x) + Number(n.w || 0) / 2;
      var y = Number(n.y);
      if (y < busY - 12 && y < 400) return;
      var suf = feederSuffix(k) || 1;
      var slot = byBase[base] || (byBase[base] = {});
      pickKnifeSlot(slot, n, cx, y, suf);
    });
    return byBase;
  }

  function shiftKnifeLabel(nodes, oldCx, oldY, dx, dy) {
    nodes.forEach(function (t) {
      if (t.t !== 'T') return;
      var s = String(t.s || '');
      if (s !== '备用' && s.indexOf('市电测试箱') < 0 && s.indexOf('测试箱') < 0) return;
      if (Math.abs(Number(t.x) - oldCx) > 52) return;
      if (Math.abs(Number(t.y) - oldY) > 52) return;
      t.x = Number(t.x) + dx;
      t.y = Number(t.y) + (dy || 0);
    });
  }

  function parallelGroupEligible(uniq) {
    if (!uniq || uniq.length < 2) return false;
    var xs = uniq.map(function (u) { return u.cx; });
    var ys = uniq.map(function (u) { return u.y; });
    var xSpread = Math.max.apply(null, xs) - Math.min.apply(null, xs);
    var ySpread = Math.max.apply(null, ys) - Math.min.apply(null, ys);
    if (ySpread >= 32 && xSpread <= 64) return true;
    if (xSpread >= 32 && ySpread <= 64) return true;
    return xSpread >= 20 || ySpread >= 20;
  }

  function pushVertLine(extra, cx, yTop, yBot, z) {
    var top = Math.min(yTop, yBot);
    var bot = Math.max(yTop, yBot);
    var len = bot - top;
    if (len < 6) return;
    var h = 26;
    /* 工程竖线：长度在 w（a=90）；x/y 按 scene 导出惯例对齐到 cx 与 yTop~yBot */
    extra.push({
      t: 'L', x: cx - len / 2, y: top + len / 2 - h / 2, w: len, h: h, a: 90, sw: 2,
      c: '#ce0c18', z: z, _rt: 1
    });
  }

  function pushHorizLine(extra, xLeft, xRight, y, z) {
    var xl = Math.min(xLeft, xRight);
    var xr = Math.max(xLeft, xRight);
    var w = Math.max(12, xr - xl);
    var h = 26;
    extra.push({
      t: 'L', x: xl, y: y - h / 2, w: w, h: h, a: 0, sw: 2,
      c: '#ce0c18', z: z, _rt: 1
    });
  }

  function lineSpanLen(n) {
    var e = lineEnds(n);
    return Math.abs(e.x2 - e.x1) + Math.abs(e.y2 - e.y1);
  }

  /** 同 cy/cx 近重合的 scene/runtime 红线去重，减堆叠交叉。 */
  function linesNearlySame(a, b) {
    if (Math.abs((Number(a.a) || 0) - (Number(b.a) || 0)) > 3) return false;
    var ea = lineEnds(a);
    var eb = lineEnds(b);
    var ha = Math.abs(ea.y1 - ea.y2) < 2;
    var hb = Math.abs(eb.y1 - eb.y2) < 2;
    if (ha !== hb) return false;
    if (ha) {
      if (Math.abs(ea.cy - eb.cy) > 2.5) return false;
      var a0 = Math.min(ea.x1, ea.x2);
      var a1 = Math.max(ea.x1, ea.x2);
      var b0 = Math.min(eb.x1, eb.x2);
      var b1 = Math.max(eb.x1, eb.x2);
      var overlap = Math.min(a1, b1) - Math.max(a0, b0);
      if (overlap < 10) return false;
      var span = Math.max(a1 - a0, b1 - b0);
      return overlap / span > 0.82;
    }
    if (Math.abs(ea.cx - eb.cx) > 2.5) return false;
    var ya0 = Math.min(ea.y1, ea.y2);
    var ya1 = Math.max(ea.y1, ea.y2);
    var yb0 = Math.min(eb.y1, eb.y2);
    var yb1 = Math.max(eb.y1, eb.y2);
    var yOverlap = Math.min(ya1, yb1) - Math.max(ya0, yb0);
    if (yOverlap < 10) return false;
    var ySpan = Math.max(ya1 - ya0, yb1 - yb0);
    return yOverlap / ySpan > 0.82;
  }

  function pruneColinearDupLines(nodes) {
    var lines = [];
    var rest = [];
    nodes.forEach(function (n) {
      if (n.t === 'L') lines.push(n);
      else rest.push(n);
    });
    var drop = {};
    var i;
    var j;
    for (i = 0; i < lines.length; i++) {
      if (drop[i]) continue;
      for (j = i + 1; j < lines.length; j++) {
        if (drop[j]) continue;
        if (!linesNearlySame(lines[i], lines[j])) continue;
        var keepI = true;
        if (lines[i]._rt && !lines[j]._rt) keepI = true;
        else if (!lines[i]._rt && lines[j]._rt) keepI = false;
        else if (lineSpanLen(lines[i]) < lineSpanLen(lines[j])) keepI = false;
        if (keepI) drop[j] = 1;
        else {
          drop[i] = 1;
          break;
        }
      }
    }
    lines.forEach(function (n, idx) {
      if (!drop[idx]) rest.push(n);
    });
    return rest;
  }

  /** 仅左邻列（与 colLeft 横向重叠），避免 T2 远列把并联组整体右移。 */
  function parallelNeighborRight(nodes, yLo, yHi, skipKeys, colLeft) {
    var maxR = -1e9;
    nodes.forEach(function (n) {
      if (n.t !== 'I' && n.t !== 'S') return;
      var k = String(n.k || '');
      if (skipKeys[k]) return;
      var ny = Number(n.y);
      if (ny < yLo - 52 || ny > yHi + KNIFE_H + 48) return;
      var nx = Number(n.x);
      var rw = knifeHorizontal(n) ? KNIFE_H : Number(n.w || KNIFE_W);
      var right = nx + rw;
      if (right + PARALLEL_COL_GAP <= colLeft) return;
      if (nx >= colLeft + KNIFE_W) return;
      maxR = Math.max(maxR, right);
    });
    return maxR;
  }

  /**
   * 并联组（5A1-T1-411 参照）：左竖母线；同列 X、沿 Y 纵向错开；刀横置；每刀高度横引。
   */
  function applyParallelKnifeGroup(uniq, nodes, busY, mainBusY, extra, backboneCols, groupedKeys, z) {
    if (!parallelGroupEligible(uniq)) return false;

    uniq.sort(function (a, b) {
      var dy = (a.y || 0) - (b.y || 0);
      if (Math.abs(dy) > 4) return dy;
      return (a.suf || 1) - (b.suf || 1);
    });
    var ys = uniq.map(function (u) { return u.y; });
    var yMin = Math.min.apply(null, ys);
    var feedY = localFeedBusY(nodes, yMin, mainBusY);
    var maxW = KNIFE_W;
    var maxH = KNIFE_H;
    var lefts = uniq.map(function (u) { return Number(u.n.x); });
    var colLeft = Math.min.apply(null, lefts);
    var yPitch = Math.max(PARALLEL_Y_PITCH_MIN, maxH * 0.55);
    var groupSkip = {};
    uniq.forEach(function (u) { groupSkip[String(u.n.k || '')] = 1; });
    var ySpanTop = yMin;
    var ySpanBot = yMin + yPitch * Math.max(1, uniq.length - 1) + KNIFE_H;
    var neighR = parallelNeighborRight(nodes, ySpanTop, ySpanBot, groupSkip, colLeft);
    if (neighR > -1e8 && colLeft < neighR + PARALLEL_COL_GAP) {
      colLeft = neighR + PARALLEL_COL_GAP;
    }
    var busColX = colLeft - PARALLEL_BUS_OFFSET;
    var busSeg = mainBusHorizAt(nodes, mainBusY, busColX);
    var busTopY = busSeg ? busSeg.cy : mainBusY;

    uniq.forEach(function (u, i) {
      groupedKeys[String(u.n.k || '')] = 1;
      var oldCx = u.cx;
      var oldY = u.y;
      var targetY = yMin + i * yPitch;
      var dx = colLeft - Number(u.n.x);
      u.n.x = colLeft;
      u.n.y = targetY;
      u.n.orient = 'h';
      u.left = colLeft;
      u.y = targetY;
      u.cx = colLeft + maxW / 2;
      shiftKnifeLabel(nodes, oldCx, oldY, dx, targetY - oldY);
    });

    var pi;
    for (pi = 1; pi < uniq.length; pi++) {
      var prevY = Number(uniq[pi - 1].n.y);
      if (Number(uniq[pi].n.y) - prevY < PARALLEL_Y_PITCH_MIN) {
        uniq[pi].n.y = prevY + PARALLEL_Y_PITCH_MIN;
        uniq[pi].y = uniq[pi].n.y;
      }
    }

    uniq.forEach(function (u) { normalizePowerKnifeSize(u.n); });

    var midYs = uniq.map(function (u) { return knifeMidY(u.n); });
    var yBot = Math.max.apply(null, midYs) + 32;
    backboneCols.push({ cx: busColX, y0: busTopY - 8, y1: yBot + 24 });

    pushVertLine(extra, busColX, busTopY, yBot, z);
    if (busSeg && busColX > busSeg.x0 + 6) {
      pushHorizLine(extra, busSeg.x0, busColX, busTopY, z);
    }
    pushHorizLine(extra, busColX, colLeft + 8, busTopY, z);
    if (feedY > mainBusY + 24 && feedY < midYs[0] - 20) {
      pushVertLine(extra, busColX, feedY, yBot, z);
    }

    uniq.forEach(function (u) {
      var tapY = knifeMidY(u.n);
      pushHorizLine(extra, busColX, colLeft, tapY, z);
    });

    return true;
  }

  /** 按回路 key 分组并联（全 power 页）。 */
  function fixParallelKnifeGroups(nodes, busY, maxZ, pageName) {
    var byBase = collectParallelKnifeSlots(nodes, busY, pageName);
    var extra = [];
    var backboneCols = [];
    var groupedKeys = {};
    var z = maxZ + 1;
    var minCount = 2;

    Object.keys(byBase).forEach(function (base) {
      var slot = byBase[base];
      var uniq = Object.keys(slot).map(function (s) { return slot[s]; });
      if (uniq.length < minCount) return;
      if (pageName === '配电室5A1') return;
      applyParallelKnifeGroup(uniq, nodes, busY, busY, extra, backboneCols, groupedKeys, z);
    });

    return { extra: extra, backboneCols: backboneCols, groupedKeys: groupedKeys, z: z };
  }

  function knifeHorizontal(n) {
    return n.orient === 'h' || Number(n.rot) === 90;
  }

  function dist2(x1, y1, x2, y2) {
    var dx = x1 - x2;
    var dy = y1 - y2;
    return dx * dx + dy * dy;
  }

  function nearPoint(x, y, px, py, tol) {
    return dist2(x, y, px, py) <= tol * tol;
  }

  function pickPrimaryKnives(nodes) {
    var pick = {};
    nodes.forEach(function (n, i) {
      if (n.t !== 'I' && n.t !== 'S') return;
      var kid = n.k ? String(n.k) : ('@' + i);
      if (pick[kid] == null) pick[kid] = n;
      else if (n.st === '1') pick[kid] = n;
      else if (pick[kid].st !== '1' && Number(n.w || 0) < Number(pick[kid].w || 999)) pick[kid] = n;
    });
    return pick;
  }

  function knifeInTap(n) {
    return bladePointInBox(n, BLADE_VX, 0);
  }

  function knifeOutTap(n) {
    return bladePointInBox(n, BLADE_VX, 114);
  }

  function lineNearPoint(n, px, py, tol) {
    var e = lineEnds(n);
    return nearPoint(e.x1, e.y1, px, py, tol) || nearPoint(e.x2, e.y2, px, py, tol);
  }

  function isIncomingFeederKey(k) {
    return /_T[12]_401_?$/.test(String(k || '')) || /_T[12]_402_?$/.test(String(k || ''));
  }

  /** scene 里 O / decor / 出线端已有线段 → 有 OUT；主馈线（非 401/402 进线）默认有 OUT。 */
  function buildKnifeOutKeys(rawNodes) {
    var out = {};
    var oRects = [];
    rawNodes.forEach(function (n) {
      if (n.t === 'O') oRects.push(n);
    });
    var prim = pickPrimaryKnives(rawNodes);
    Object.keys(prim).forEach(function (k) {
      var kn = prim[k];
      if (kn.decor) out[k] = 1;
      if (!out[k] && isMainFeederKey(k) && !isIncomingFeederKey(k)) out[k] = 1;
      if (!out[k] && /配电室\d[AB]\d_U[12]_/.test(k)) out[k] = 1;
      if (!isIncomingFeederKey(k)) {
        var tap = knifeOutTap(kn);
        rawNodes.forEach(function (ln) {
          if (ln.t !== 'L') return;
          if (lineNearPoint(ln, tap.x, tap.y, 18)) out[k] = 1;
        });
      }
      var bx = Number(kn.x);
      var by = Number(kn.y);
      var bw = Number(kn.w || 60);
      var bh = Number(kn.h || 98);
      oRects.forEach(function (o) {
        var ox = Number(o.x);
        var oy = Number(o.y);
        var ow = Number(o.w || 20);
        var oh = Number(o.h || 14);
        if (ox < bx + bw + 48 && ox + ow > bx - 48 && oy < by + bh + 48 && oy + oh > by - 48) {
          out[k] = 1;
        }
      });
    });
    return out;
  }

  function lineLeavesFromNonOutKnife(n, prim, outKeys) {
    var e = lineEnds(n);
    var tol = 18;
    var ks = Object.keys(prim);
    for (var i = 0; i < ks.length; i++) {
      var k = ks[i];
      if (outKeys[k]) continue;
      var kn = prim[k];
      var out = knifeOutTap(kn);
      var hitOut = nearPoint(e.x1, e.y1, out.x, out.y, tol) || nearPoint(e.x2, e.y2, out.x, out.y, tol);
      if (!hitOut) continue;
      var inn = knifeInTap(kn);
      var hitIn = nearPoint(e.x1, e.y1, inn.x, inn.y, tol) || nearPoint(e.x2, e.y2, inn.x, inn.y, tol);
      if (hitIn) continue;
      return true;
    }
    return false;
  }

  /** 进线竖线落到 busY 但在顶母线横段左外，与补的横引重复且不通。 */
  function isDanglingIncomingVert(n, nodes, prim, busY, groupedKeys) {
    if (!isVertLine(n)) return false;
    var e = lineEnds(n);
    var yLo = Math.min(e.y1, e.y2);
    if (Math.abs(yLo - busY) > 14) return false;
    var ks = Object.keys(prim);
    for (var i = 0; i < ks.length; i++) {
      var k = ks[i];
      if (groupedKeys[k] || !isIncomingFeederKey(k)) continue;
      var tap = knifeInTap(prim[k]);
      if (Math.abs(e.cx - tap.x) > 24) continue;
      var seg = mainBusHorizAt(nodes, busY, tap.x);
      if (!seg) return false;
      return e.cx < seg.x0 - 6;
    }
    return false;
  }

  /**
   * 总图 CAD：T1-401/T2-402 上端竖落顶母线，再沿顶母横引至同段首组并联竖母（如 T1-411 / T2-421）。
   */
  function connectIncomingFeeders(nodes, extra, busY, backboneCols, groupedKeys, z) {
    if (busY == null) return;
    var busSegs = listTopBusHoriz(nodes, busY);
    var prim = pickPrimaryKnives(nodes);
    Object.keys(prim).forEach(function (k) {
      if (groupedKeys[k] || !isIncomingFeederKey(k)) return;
      var kn = prim[k];
      var tap = knifeInTap(kn);
      if (tap.y < busY + 6) return;
      var seg = pickBusSegForX(busSegs, tap.x);
      var tieY = seg ? seg.cy : busY;
      var joinX = pickBackboneJoinX(tap, backboneCols, seg);
      /* CAD：进线竖线必须在刀位正上方，再沿顶母横引至并联竖母 */
      pushVertLine(extra, tap.x, tieY, tap.y, z);
      if (seg && tap.x > seg.x1 + 8) {
        pushHorizLine(extra, seg.x1, tap.x, tieY, z);
      }
      var hx0 = Math.min(tap.x, joinX);
      var hx1 = Math.max(tap.x, joinX);
      if (hx1 - hx0 > 6) {
        pushHorizLine(extra, hx0, hx1, tieY, z);
      }
    });
  }

  /** 401/402 无 OUT：去掉 scene 在出线端下方多余的竖馈线（与 runtime 重复）。 */
  function isIncomingFeederDownVert(n, primKnives, groupedKeys, outKeys) {
    if (!isVertLine(n) || n._rt) return false;
    var e = lineEnds(n);
    var keys = Object.keys(primKnives);
    var i;
    for (i = 0; i < keys.length; i++) {
      var k = keys[i];
      if (groupedKeys[k] || !isIncomingFeederKey(k)) continue;
      if (outKeys[k]) continue;
      var kn = primKnives[k];
      var inn = knifeInTap(kn);
      var out = knifeOutTap(kn);
      if (Math.abs(e.cx - inn.x) > 14) continue;
      var yLo = Math.min(e.y1, e.y2);
      if (yLo >= out.y - 10) return true;
    }
    return false;
  }

  /**
   * 按总图 DWG「配电室*」块对齐非并联 L 线（并联组已在 fixMainBranchParallel 内完成）。
   * 标准拓扑：进线 401 上端 → 顶母线 → 与 411 并联竖母共母；402 → 同段 T2 顶母 → 421 组竖母。
   */
  function shiftLineCx(n, newCx) {
    var e = lineEnds(n);
    n.x = Number(n.x) + (newCx - e.cx);
  }

  function shiftLineCy(n, newCy) {
    var e = lineEnds(n);
    n.y = Number(n.y) + (newCy - e.cy);
  }

  function spanLineToY(n, y) {
    var e = lineEnds(n);
    var lo = Math.min(e.y1, e.y2);
    var hi = Math.max(e.y1, e.y2);
    if (y < lo) lo = y;
    if (y > hi) hi = y;
    var len = Math.max(4, hi - lo);
    var h = Number(n.h) || 26;
    n.w = len;
    n.x = e.cx - len / 2;
    n.y = lo + len / 2 - h / 2;
  }

  function spanLineToX(n, x) {
    var e = lineEnds(n);
    var lo = Math.min(e.x1, e.x2);
    var hi = Math.max(e.x1, e.x2);
    if (x < lo) lo = x;
    if (x > hi) hi = x;
    var w = Math.max(4, hi - lo);
    var hh = Number(n.h) || 26;
    n.w = w;
    n.x = lo;
    n.y = e.cy - hh / 2;
  }

  function wireNearPoint(nodes, px, py, tol) {
    var t2 = tol * tol;
    var i;
    for (i = 0; i < nodes.length; i++) {
      var n = nodes[i];
      if (n.t !== 'L') continue;
      if (n.c && n.c !== '#ce0c18') continue;
      var e = lineEnds(n);
      if (isHorizLine(n)) {
        var cy = (e.y1 + e.y2) / 2;
        if (Math.abs(cy - py) > tol) continue;
        var x0 = Math.min(e.x1, e.x2);
        var x1 = Math.max(e.x1, e.x2);
        if (px >= x0 - tol && px <= x1 + tol) return true;
      }
      if (isVertLine(n)) {
        if (Math.abs(e.cx - px) > tol) continue;
        var y0 = Math.min(e.y1, e.y2);
        var y1 = Math.max(e.y1, e.y2);
        if (py >= y0 - tol && py <= y1 + tol) return true;
      }
      if (dist2(e.x1, e.y1, px, py) <= t2 || dist2(e.x2, e.y2, px, py) <= t2) return true;
    }
    return false;
  }

  function nearestRedBusY(nodes, x, belowY) {
    var best = null;
    nodes.forEach(function (n) {
      if (n.t !== 'L' || !isHorizLine(n)) return;
      if (n.c && n.c !== '#ce0c18') return;
      var e = lineEnds(n);
      var cy = (e.y1 + e.y2) / 2;
      if (cy > belowY - 4) return;
      var x0 = Math.min(e.x1, e.x2);
      var x1 = Math.max(e.x1, e.x2);
      if (x < x0 - 12 || x > x1 + 12) return;
      if (best == null || cy < best) best = cy;
    });
    return best;
  }

  /** 5A1：scene 线 cx 与 normalize 后铰接点对齐，并补缺失竖落线。 */
  function stitch5A1KnifeWires(nodes, primKnives, outKeys, busY, maxZ) {
    var snapX = 28;
    var klist = Object.keys(primKnives);
    klist.forEach(function (k) {
      var kn = primKnives[k];
      var inn = knifeInTap(kn);
      var out = knifeOutTap(kn);
      var horiz = knifeHorizontal(kn);
      nodes.forEach(function (n) {
        if (n.t !== 'L' || (n.c && n.c !== '#ce0c18')) return;
        var e = lineEnds(n);
        if (isVertLine(n) && Math.abs(e.cx - inn.x) <= snapX) {
          shiftLineCx(n, inn.x);
          spanLineToY(n, inn.y);
          if (outKeys[k]) spanLineToY(n, out.y);
        }
        if (horiz && isHorizLine(n)) {
          var mid = knifeMidY(kn);
          if (Math.abs(e.cy - mid) <= snapX || Math.abs(e.cy - inn.y) <= snapX) {
            shiftLineCy(n, mid);
            spanLineToX(n, inn.x);
            if (outKeys[k]) spanLineToX(n, out.x);
          }
        }
      });
    });

    var z = (maxZ || 900) + 50;
    klist.forEach(function (k) {
      var kn = primKnives[k];
      var inn = knifeInTap(kn);
      if (wireNearPoint(nodes, inn.x, inn.y, 2.5)) return;
      if (knifeHorizontal(kn)) return;
      var tieY = nearestRedBusY(nodes, inn.x, inn.y) || busY;
      if (tieY == null || inn.y <= tieY + 4) return;
      pushVertLine(nodes, inn.x, tieY, inn.y, z);
      var hx0 = Math.min(inn.x, inn.x);
      var segX = inn.x;
      if (Math.abs(tieY - busY) < 12) {
        /* 顶母横引已在 connectIncomingFeeders */
      }
    });

    klist.forEach(function (k) {
      var kn = primKnives[k];
      if (!knifeHorizontal(kn)) return;
      var inn = knifeInTap(kn);
      var out = knifeOutTap(kn);
      var mid = knifeMidY(kn);
      var busX = inn.x - PARALLEL_BUS_OFFSET;
      nodes.forEach(function (n) {
        if (n.t !== 'L' || !isHorizLine(n) || (n.c && n.c !== '#ce0c18')) return;
        var e = lineEnds(n);
        var cy = (e.y1 + e.y2) / 2;
        if (Math.abs(cy - mid) > 4) return;
        spanLineToX(n, Math.min(busX, inn.x));
        spanLineToX(n, out.x);
      });
      if (!wireNearPoint(nodes, inn.x, inn.y, 2.5)) {
        pushHorizLine(nodes, busX, out.x, mid, z);
      }
      if (outKeys[k] && !wireNearPoint(nodes, out.x, out.y, 2.5)) {
        pushHorizLine(nodes, inn.x, out.x, mid, z);
      }
    });

    klist.forEach(function (k) {
      if (!outKeys[k]) return;
      var kn = primKnives[k];
      var out = knifeOutTap(kn);
      if (wireNearPoint(nodes, out.x, out.y, 2.5)) return;
      if (knifeHorizontal(kn)) return;
      var yLo = out.y + 4;
      var yHi = out.y + 48;
      nodes.forEach(function (n) {
        if (n.t !== 'L' || !isVertLine(n) || (n.c && n.c !== '#ce0c18')) return;
        var e = lineEnds(n);
        if (Math.abs(e.cx - out.x) > 3) return;
        var lo = Math.min(e.y1, e.y2);
        if (lo <= out.y + 8) spanLineToY(n, out.y);
      });
      if (!wireNearPoint(nodes, out.x, out.y, 2.5)) {
        pushVertLine(nodes, out.x, out.y, out.y + 36, z);
      }
    });
  }

  function fixPowerLinesFromCadHints(pg) {
    if (!pg || pg.kind !== 'power' || pg._cadLinesFixed) return;
    if (w.HX5A1 && w.HX5A1.applies(pg)) {
      w.HX5A1.apply(pg);
      return;
    }
    if (w.HXRooms && w.HXRooms.applies(pg)) {
      w.HXRooms.apply(pg);
      return;
    }
    if (!pg._mainParallelFixed) fixMainBranchParallel(pg);
    pg._cadLinesFixed = 1;
  }

  /* 顶母线 + 全页并联闸刀组；仅 groupedKeys 内改刀与线，单刀保持 scene 原样。 */
  function fixMainBranchParallel(pg) {
    if (!pg || pg.kind !== 'power' || pg._mainParallelFixed) return;
    var nodes = pg.nodes || [];
    var busY = findTopBusY(nodes);
    if (busY == null) {
      pg._mainParallelFixed = 1;
      return;
    }

    var pageName = pg.name || '';
    var rawNodes = nodes.slice();
    var outKeys = buildKnifeOutKeys(rawNodes);
    var knifeSnap = snapshotKnifeNodes(nodes);
    var maxZ = 0;
    nodes.forEach(function (n) { maxZ = Math.max(maxZ, Number(n.z) || 0); });

    var parallel = fixParallelKnifeGroups(nodes, busY, maxZ, pageName);
    maxZ = parallel.z + parallel.extra.length;
    var groupedKeys = parallel.groupedKeys;
    var backboneCols = parallel.backboneCols;
    var extraLines = parallel.extra.slice();

    restoreNonGroupedKnives(nodes, groupedKeys, knifeSnap);
    normalizePowerKnifeSizes(nodes);
    connectIncomingFeeders(nodes, extraLines, busY, backboneCols, groupedKeys, maxZ + 2);
    var primKnives = pickPrimaryKnives(nodes);

    function lineHitsBackbone(n) {
      if (!isVertLine(n)) return false;
      var e = lineEnds(n);
      var span = Math.abs(e.y2 - e.y1);
      if (span < 68) return false;
      var yLo = Math.min(e.y1, e.y2);
      var yHi = Math.max(e.y1, e.y2);
      for (var i = 0; i < backboneCols.length; i++) {
        var b = backboneCols[i];
        if (Math.abs(e.cx - b.cx) > 36) continue;
        if (yHi > b.y0 && yLo < b.y1) return true;
      }
      return false;
    }

    /** 仅去掉并联竖母线列上的旧 scene 串联横引（不碰单刀原线）。 */
    function isParallelSeriesHorizTap(n) {
      if (!isHorizLine(n)) return false;
      if (!backboneCols.length) return false;
      var cy = Number(n.y) + Number(n.h || 0) / 2;
      var w = Math.abs(Number(n.w) || 0);
      if (w < 40 || w > 130) return false;
      var e = lineEnds(n);
      for (var i = 0; i < backboneCols.length; i++) {
        var b = backboneCols[i];
        if (Math.abs(e.cy - cy) > 10) continue;
        if (Math.abs(e.cx - b.cx) > 90) continue;
        if (cy >= b.y0 - 20 && cy <= b.y1 + 20) return true;
      }
      return false;
    }

    var kept = [];
    nodes.forEach(function (n) {
      if (n.t !== 'L') {
        if ((n.t === 'I' || n.t === 'S') && n.k && outKeys[String(n.k)]) n.decor = 1;
        kept.push(n);
        return;
      }
      if (isParallelSeriesHorizTap(n)) return;
      if (lineHitsBackbone(n)) return;
      if (isDanglingIncomingVert(n, nodes, primKnives, busY, groupedKeys)) return;
      if (isIncomingFeederDownVert(n, primKnives, groupedKeys, outKeys)) return;
      if (lineLeavesFromNonOutKnife(n, primKnives, outKeys)) return;
      kept.push(n);
    });

    extraLines.forEach(function (ln) { kept.push(ln); });

    if (pageName === '配电室5A1') {
      stitch5A1KnifeWires(kept, primKnives, outKeys, busY, maxZ + 20);
    }

    pg.nodes = pruneColinearDupLines(kept);
    pg._mainParallelFixed = 1;
  }

  function prep(pg) {
    if (!pg || pg._prepped) return pg;
    fixPowerLinesFromCadHints(pg);
    (pg.nodes || []).forEach(function (n) {
      if (n.t === 'T' && !n.k && !n.lay) {
        var k = guessKey(pg, n.s);
        if (k) n.k = k;
      }
    });
    pg._prepped = 1;
    return pg;
  }

  function pageById(id) {
    if (id && String(id).indexOf('live:') === 0) return L().cabinetPage(id);
    for (var i = 0; i < scenes.length; i++) if (scenes[i].id === id) return scenes[i];
    return null;
  }

  /* 手绘列头页有的是 A 列头坐标整页复制，测点仍挂在 A列头。这种页不能拿来当 B/C 的画面。 */
  function boundToCab(pg, cab) {
    var own = 0;
    var other = 0;
    (pg.nodes || []).forEach(function (n) {
      [n.k, n.p].forEach(function (mark) {
        if (!mark) return;
        if (String(mark).indexOf(cab) === 0) own++;
        else other++;
      });
    });
    return own > 0 && other === 0;
  }

  function sceneFor(room, cab) {
    var code = (String(room).match(/\d[AB]\d/) || [''])[0];
    if (!code) return null;
    for (var i = 0; i < scenes.length; i++) {
      var p = scenes[i];
      if (p.kind === 'column' && p.cab === cab && String(p.group || '').indexOf(code) >= 0 && boundToCab(p, cab)) return p;
    }
    return null;
  }

  function showTopoComm(main, resolved, item) {
    if (!main || !w.HXArch) return;
    page = null;
    var mgrName = (resolved && resolved.mgr) || (item && item.name) || '';
    var coms = w.HXArch.mgrComPorts(mgrName);
    var title = (resolved && resolved.pageHint) || (mgrName.split('(')[0] + '通讯图');
    var rules = (resolved && resolved.rules) || w.HXArch.netRules || {};
    var rows = [];
    var bound = 0;
    var total = 0;
    var i, j, com, devs, dev, pt, st, short, uuid;

    function comBlockHtml(comNode, role) {
      if (!comNode) return '<div class="hx-topo-com is-empty"></div>';
      var kids = comNode.children || [];
      var body = kids.map(function (d) {
        var nm = d.name || '';
        total++;
        pt = topoCommPoint(nm);
        uuid = pt && pt.uuid ? pt.uuid : '';
        if (uuid) {
          bound++;
          rows.push(pt);
        }
        st = topoDevState(nm, uuid);
        short = topoShortName(nm);
        return '<div class="hx-topo-dev is-' + esc(st.kind) + '" data-dev="' + esc(nm) + '"' +
          (uuid ? ' data-uuid="' + esc(uuid) + '"' : '') +
          ' title="' + esc(st.kind === 'unk' ? (short + ' · 未绑通讯状态测点') : (short + ' · 通讯' + st.label)) + '">' +
          '<span class="hx-topo-flow" ' + (st.kind === 'ok' ? '' : 'hidden') + '></span>' +
          '<span class="hx-topo-dev-lab">' + esc(short + (st.label || '')) + '</span>' +
          '</div>';
      }).join('');
      var roleCls = role ? (' hx-topo-com--' + role) : '';
      /* 列竖干只接到 COM 标题；设备挂在 COM 下独立竖干上 */
      var branch = kids.length
        ? ('<div class="hx-topo-com-branch">' +
            '<div class="hx-topo-com-spine" aria-hidden="true"></div>' +
            '<div class="hx-topo-com-devs">' + body + '</div>' +
          '</div>')
        : '';
      return '<div class="hx-topo-com' + roleCls + '">' +
        '<div class="hx-topo-com-h">' + esc(comNode.name) + '</div>' +
        branch +
        '</div>';
    }

    /* 奇偶 COM 成列：主母线垂下竖干 → 短横接 COM → COM 下再竖干树状挂设备 */
    var cols = [];
    for (i = 0; i < coms.length; i += 2) {
      cols.push({ top: coms[i], bottom: coms[i + 1] || null });
    }

    var colsHtml = cols.map(function (col) {
      var hasBot = !!col.bottom;
      return '<div class="hx-topo-col' + (hasBot ? ' has-bot' : '') + '">' +
        '<div class="hx-topo-col-spine" aria-hidden="true"></div>' +
        comBlockHtml(col.top, 'top') +
        (col.bottom ? comBlockHtml(col.bottom, 'bot') : '<div class="hx-topo-com is-empty"></div>') +
        '</div>';
    }).join('');

    var note = bound
      ? ('已绑通讯状态 ' + bound + '/' + total + '；绿=正常（流动框可见），红=故障（流动框隐藏）；灰·未绑为占位')
      : ('本机 ' + total + ' 个设备暂未接到「*_通讯状态」测点，全部占位显示，不假闪');

    main.innerHTML =
      '<div class="hx-topo">' +
        '<div class="hx-topo-head">' +
          '<div class="hx-topo-title">' + esc(title) + '</div>' +
          '<div class="hx-topo-sub">' + esc(mgrName) +
            (coms.length ? (' · ' + coms.length + ' 个 COM') : '') +
            ' · ' + esc(note) + '</div>' +
        '</div>' +
        '<div class="hx-topo-canvas">' +
          '<div class="hx-topo-mgr" title="' + esc(mgrName) + '">' +
            '<div class="hx-topo-mgr-ico" aria-hidden="true">' +
              '<svg viewBox="0 0 48 40" width="44" height="36">' +
                '<rect x="4" y="6" width="40" height="28" rx="3" fill="#1a4a7a" stroke="#9fd0ff" stroke-width="1.6"/>' +
                '<path d="M12 20c4-8 20-8 24 0" fill="none" stroke="#dff4ff" stroke-width="2"/>' +
                '<circle cx="24" cy="20" r="3" fill="#7fe9ff"/>' +
              '</svg>' +
            '</div>' +
            '<div class="hx-topo-mgr-name">' + esc(mgrName.split('(')[0]) + '</div>' +
          '</div>' +
          '<div class="hx-topo-stem"></div>' +
          '<div class="hx-topo-bus"></div>' +
          '<div class="hx-topo-cols">' + (colsHtml || '<div class="hx-pt">该管理机下没有 COM 口</div>') + '</div>' +
        '</div>' +
        '<div class="hx-topo-legend">' +
          esc(rules.ok || '绿=正常且外框流动可见') + ' · ' +
          esc(rules.bad || '红=异常且流动框不可见') + ' · ' +
          esc(rules.bind || '设备只绑通讯状态') +
        '</div>' +
      '</div>';

    pageRows = rows;
    L().allow(pageRows);
    if (pageRows.length) L().fetch(pageRows).then(paintTopoDevs);
    else paintTopoDevs();
  }

  function showStub(main, resolved, item) {
    if (!main) return;
    var links = (resolved && resolved.links) || [];
    var hint = esc((resolved && resolved.reason) || '暂无画面');
    var btns = links.map(function (lk, i) {
      return '<button type="button" class="hx-drill-go" data-link-i="' + i + '">前往 ' + esc(lk.label || lk.design) + '</button>';
    }).join('');
    var sug = resolved && resolved.kind === 'power' && w.HXArch ? w.HXArch.suggestFromPower(resolved.design) : null;
    if (sug && sug.item) {
      btns = '<button type="button" class="hx-drill-go is-primary" data-suggest="1">下钻 · ' +
        esc(sug.item.name) + (sug.item.resolved && sug.item.resolved.cab ? ' · ' + sug.item.resolved.cab.name : '') +
        '</button>' + btns;
    }
    main.innerHTML =
      '<div class="hx-site-note">' +
        '<div class="hx-title">' + esc((item && item.name) || (resolved && resolved.design) || '现场节点') + '</div>' +
        '<div class="hx-pt">' + hint + '。可通过左侧架构树继续钻取，或点下方关联入口。</div>' +
        (btns ? '<div class="hx-drill-actions">' + btns + '</div>' : '') +
      '</div>';
    main.querySelectorAll('[data-suggest]').forEach(function (btn) {
      btn.addEventListener('click', function () {
        if (drillPick && sug && sug.item) drillPick(sug.item.resolved.room, sug.item.resolved.cab, sug.item);
      });
    });
    main.querySelectorAll('[data-link-i]').forEach(function (btn) {
      btn.addEventListener('click', function () {
        var lk = links[Number(btn.getAttribute('data-link-i'))];
        if (!lk || !w.HXArch || !treeFlat) return;
        var picked = w.HXArch.pickByDesign(treeFlat, lk.design);
        if (picked && drillPick) drillPick(picked.resolved.room, picked.resolved.cab, picked);
        else if (drillPick && lk.design) {
          var path = w.HXArch.findLeafPath(lk.design);
          if (path) {
            var r = w.HXArch.resolveLeaf(path.name, path.ancestors);
            drillPick(r.room, r.cab, { name: path.name, ancestors: path.ancestors, resolved: r });
          }
        }
      });
    });
    L().allow([]);
    pageRows = [];
  }

  function findPageByName(name) {
    if (!name) return null;
    var i, p;
    for (i = 0; i < scenes.length; i++) {
      p = scenes[i];
      if (p && p.name === name) return p;
    }
    if (powerByName[name]) return { name: name, id: powerByName[name], kind: 'topo' };
    return null;
  }

  function setCrumb(labels) {
    var el = hostEl && hostEl.querySelector('#hxCrumb');
    if (!el) return;
    var parts = labels || [];
    if (!parts.length) {
      el.innerHTML = '<span class="hx-crumb-item is-cur">架构钻取</span>';
      return;
    }
    el.innerHTML = parts.map(function (p, i) {
      var cur = i === parts.length - 1 ? ' is-cur' : '';
      return '<span class="hx-crumb-item' + cur + '">' + esc(p) + '</span>' +
        (i < parts.length - 1 ? '<span class="hx-crumb-sep">›</span>' : '');
    }).join('');
  }

  function setDrillLinks(resolved) {
    var bar = hostEl && hostEl.querySelector('#hxDrillLinks');
    if (!bar) return;
    if (!resolved || resolved.kind !== 'power' || !w.HXArch) {
      bar.innerHTML = '';
      return;
    }
    var sug = w.HXArch.suggestFromPower(resolved.design || resolved.room && resolved.room.name);
    if (!sug || !sug.item) {
      bar.innerHTML = '';
      return;
    }
    bar.innerHTML = '<span class="hx-drill-hint">关联</span>' +
      '<button type="button" class="hx-drill-chip" data-suggest="1">' +
      esc(sug.target.label || sug.item.name) + ' ›</button>';
    bar.querySelector('[data-suggest]').addEventListener('click', function () {
      if (drillPick) drillPick(sug.item.resolved.room, sug.item.resolved.cab, sug.item);
    });
  }

  function applyPick(room, cab, item) {
    if (!hostEl) return;
    var main = hostEl.querySelector('.hx-main');
    drillActive = item || null;
    var labels = w.HXArch ? w.HXArch.crumbLabels(item && item.ancestors, item && item.name, cab && cab.name) : [];
    setCrumb(labels);
    var resolved = item && item.resolved;
    if (resolved && resolved.kind === 'topo') {
      if (resolved.ok && resolved.mgr && !resolved.com) {
        showTopoComm(main, resolved, item);
      } else {
        showStub(main, resolved, item);
      }
      setDrillLinks(null);
      return;
    }
    if (resolved && resolved.ok && room) {
      openRoomCab(hostEl, room, cab);
      if (resolved.kind === 'power' || (item && /^配电室/.test(item.name || ''))) {
        setDrillLinks({ kind: 'power', design: (item && item.name) || room.name });
      } else {
        setDrillLinks(null);
      }
      return;
    }
    if (resolved && !resolved.ok) {
      showStub(main, resolved, item);
      setDrillLinks(resolved);
      return;
    }
    if (room) openRoomCab(hostEl, room, cab);
  }

  function roomCodeOf(pg) {
    return String((pg && pg.name) || '').replace(/^配电室/, '');
  }

  /* CAD 原图：100% = 8px / 图上单位。放大、缩小、滚轮、拖动、局部框选跟一次图同一套。 */
  var CAD_PX = 8;

  function cadState(pane) {
    return pane && pane._cad;
  }

  function applyCad(pane) {
    var st = cadState(pane);
    var stage = pane && pane.querySelector('.hx-cad-stage');
    if (!st || !stage) return;
    stage.style.width = st.baseW + 'px';
    stage.style.height = st.baseH + 'px';
    st.tx = snapPx(st.tx);
    st.ty = snapPx(st.ty);
    stage.style.transformOrigin = '0 0';
    stage.style.transform = 'translate(' + st.tx + 'px,' + st.ty + 'px) scale(' + st.sc + ')';
    /* CAD 用户单位先按 8px 铺开，再被适应窗口缩放。线宽仍按原文件比例，只抬高到屏幕上看得见。 */
    thickenStrokes(stage, 'line, path, polyline, polygon, circle, ellipse', CAD_PX * st.sc);
  }

  function cadFit(pane) {
    var st = cadState(pane);
    var view = pane && pane.querySelector('.hx-cad-view');
    if (!st || !view) return;
    var cw = view.clientWidth;
    var ch = view.clientHeight;
    if (cw < 40 || ch < 40 || !(st.baseW > 1) || !(st.baseH > 1)) return;
    var pad = 28;
    var s = Math.min((cw - pad * 2) / st.baseW, (ch - pad * 2) / st.baseH);
    if (!isFinite(s) || s <= 0) s = 0.4;
    st.sc = Math.max(0.05, Math.min(12, s));
    st.tx = (cw - st.baseW * st.sc) / 2;
    st.ty = (ch - st.baseH * st.sc) / 2;
    st.user = false;
    st.marquee = false;
    view.classList.remove('is-marquee');
    var btn = pane.querySelector('[data-hx-cad-zoom="area"]');
    if (btn) btn.classList.remove('is-on');
    applyCad(pane);
  }

  function cadZoomAt(pane, clientX, clientY, next) {
    var st = cadState(pane);
    var view = pane && pane.querySelector('.hx-cad-view');
    if (!st || !view) return;
    var rect = view.getBoundingClientRect();
    var px = clientX - rect.left;
    var py = clientY - rect.top;
    var sc = st.sc || 1;
    next = Math.max(0.05, Math.min(12, next));
    var bx = (px - st.tx) / sc;
    var by = (py - st.ty) / sc;
    st.user = true;
    st.sc = next;
    st.tx = px - bx * next;
    st.ty = py - by * next;
    applyCad(pane);
  }

  function cadZoomCenter(pane, factor) {
    var view = pane && pane.querySelector('.hx-cad-view');
    if (!view) return;
    var rect = view.getBoundingClientRect();
    var st = cadState(pane);
    cadZoomAt(pane, rect.left + rect.width / 2, rect.top + rect.height / 2, (st && st.sc || 1) * factor);
  }

  function cadZoom100(pane) {
    var st = cadState(pane);
    var view = pane && pane.querySelector('.hx-cad-view');
    if (!st || !view) return;
    st.user = true;
    st.sc = 1;
    st.tx = (view.clientWidth - st.baseW) / 2;
    st.ty = (view.clientHeight - st.baseH) / 2;
    applyCad(pane);
  }

  function bindCadControls(main) {
    if (main._hxCad) return;
    main._hxCad = 1;
    var drag = null;
    function paneOf(el) {
      return el && el.closest ? el.closest('.hx-cad-pane') : null;
    }
    function viewOf(pane) {
      return pane && pane.querySelector('.hx-cad-view');
    }
    main.addEventListener('wheel', function (ev) {
      if (!main.classList.contains('is-cad')) return;
      var pane = paneOf(ev.target);
      if (!pane || !cadState(pane)) return;
      if (ev.target.closest && ev.target.closest('.hx-zoom-bar')) return;
      ev.preventDefault();
      ev.stopPropagation();
      var dy = ev.deltaY;
      if (ev.deltaMode === 1) dy *= 16;
      var st = cadState(pane);
      cadZoomAt(pane, ev.clientX, ev.clientY, (st.sc || 1) * (dy > 0 ? 0.9 : 1.12));
    }, { passive: false, capture: true });
    main.addEventListener('dblclick', function (ev) {
      if (!main.classList.contains('is-cad')) return;
      var pane = paneOf(ev.target);
      if (!pane || !cadState(pane)) return;
      if (ev.target.closest && ev.target.closest('.hx-zoom-bar')) return;
      ev.preventDefault();
      ev.stopPropagation();
      var st = cadState(pane);
      cadZoomAt(pane, ev.clientX, ev.clientY, (st.sc || 1) * (ev.shiftKey ? 1 / 1.35 : 1.35));
    }, true);
    main.addEventListener('pointerdown', function (ev) {
      if (!main.classList.contains('is-cad') || ev.button !== 0) return;
      var pane = paneOf(ev.target);
      var st = cadState(pane);
      var view = viewOf(pane);
      if (!st || !view) return;
      if (ev.target.closest && ev.target.closest('.hx-zoom-bar, .hx-view-tabs')) return;
      if (!ev.target.closest('.hx-cad-view')) return;
      if (st.marquee) {
        /* 框画在 .hx-cad-view 上。跟一次图一样用 fitLocal：client 减去这块的
           getBoundingClientRect，再换回它自己的 CSS 像素。不要用 SVG 局部坐标。 */
        var p0 = fitLocal(view, ev.clientX, ev.clientY);
        drag = { mode: 'area', pane: pane, view: view, x: p0.x, y: p0.y };
        var box = document.createElement('div');
        box.className = 'hx-marquee';
        box.style.left = p0.x + 'px';
        box.style.top = p0.y + 'px';
        view.appendChild(box);
        drag.box = box;
        if (view.setPointerCapture) view.setPointerCapture(ev.pointerId);
        return;
      }
      drag = { mode: 'pan', pane: pane, view: view, x: ev.clientX, y: ev.clientY, tx: st.tx, ty: st.ty, moved: false };
      view.classList.add('is-pan');
      if (view.setPointerCapture) view.setPointerCapture(ev.pointerId);
    });
    main.addEventListener('pointermove', function (ev) {
      if (!drag) return;
      if (drag.mode === 'pan') {
        var dx = ev.clientX - drag.x;
        var dy = ev.clientY - drag.y;
        if (!drag.moved && Math.hypot(dx, dy) < 4) return;
        drag.moved = true;
        var st = cadState(drag.pane);
        if (!st) return;
        st.user = true;
        st.tx = drag.tx + dx;
        st.ty = drag.ty + dy;
        applyCad(drag.pane);
        return;
      }
      if (!drag.box) return;
      var p = fitLocal(drag.view, ev.clientX, ev.clientY);
      var x = p.x;
      var y = p.y;
      var x0 = Math.min(drag.x, x);
      var y0 = Math.min(drag.y, y);
      drag.box.style.left = x0 + 'px';
      drag.box.style.top = y0 + 'px';
      drag.box.style.width = Math.abs(x - drag.x) + 'px';
      drag.box.style.height = Math.abs(y - drag.y) + 'px';
    });
    function endCadDrag() {
      if (!drag) return;
      var d = drag;
      drag = null;
      if (d.view) d.view.classList.remove('is-pan');
      if (d.mode !== 'area' || !d.box) return;
      var w = parseFloat(d.box.style.width) || 0;
      var h = parseFloat(d.box.style.height) || 0;
      var x0 = parseFloat(d.box.style.left) || 0;
      var y0 = parseFloat(d.box.style.top) || 0;
      if (d.box.parentNode) d.box.parentNode.removeChild(d.box);
      var st = cadState(d.pane);
      if (st) {
        st.marquee = false;
        if (d.view) d.view.classList.remove('is-marquee');
        var areaBtn = d.pane.querySelector('[data-hx-cad-zoom="area"]');
        if (areaBtn) areaBtn.classList.remove('is-on');
      }
      if (!st || w < 12 || h < 12) return;
      var cw = d.view.clientWidth;
      var ch = d.view.clientHeight;
      var sc0 = st.sc || 1;
      var b0x = (x0 - st.tx) / sc0;
      var b0y = (y0 - st.ty) / sc0;
      var bw = w / sc0;
      var bh = h / sc0;
      var sc = Math.min(cw / bw, ch / bh) * 0.94;
      st.user = true;
      st.sc = Math.max(0.05, Math.min(12, sc));
      st.tx = (cw - bw * st.sc) / 2 - b0x * st.sc;
      st.ty = (ch - bh * st.sc) / 2 - b0y * st.sc;
      applyCad(d.pane);
    }
    main.addEventListener('pointerup', endCadDrag);
    main.addEventListener('pointercancel', endCadDrag);
    main.addEventListener('click', function (ev) {
      var btn = ev.target.closest ? ev.target.closest('[data-hx-cad-zoom]') : null;
      if (!btn) return;
      var pane = paneOf(btn);
      if (!pane || !cadState(pane)) return;
      ev.preventDefault();
      ev.stopPropagation();
      var act = btn.getAttribute('data-hx-cad-zoom');
      var st = cadState(pane);
      var view = viewOf(pane);
      if (act === 'in') cadZoomCenter(pane, 1.2);
      else if (act === 'out') cadZoomCenter(pane, 1 / 1.2);
      else if (act === 'fit') cadFit(pane);
      else if (act === '100') cadZoom100(pane);
      else if (act === 'area' && view) {
        st.marquee = !st.marquee;
        btn.classList.toggle('is-on', st.marquee);
        view.classList.toggle('is-marquee', st.marquee);
      }
    });
    window.addEventListener('resize', function () {
      if (!main.classList.contains('is-cad')) return;
      var pane = main.querySelector('.hx-cad-pane');
      var st = cadState(pane);
      if (!st || st.user) return;
      cadFit(pane);
    });
  }

  function mountCadSvg(pane, code, svg) {
    svg = String(svg).replace(/<script[\s\S]*?<\/script>/gi, '');
    pane.innerHTML =
      '<div class="hx-cad-view"><div class="hx-cad-stage">' + svg + '</div></div>' +
      '<div class="hx-cad-cap">' + esc(code) + ' · ' + ((code === '5A4' || code === '5B4') ? '总图裁出' : (/^[234][AB][12]$/.test(code) ? '2026-11-05 单室 CAD' : '2026-10-30 单室 CAD')) + '</div>' +
      '<div class="hx-zoom-bar">' +
        '<button type="button" data-hx-cad-zoom="in">放大</button>' +
        '<button type="button" data-hx-cad-zoom="out">缩小</button>' +
        '<button type="button" data-hx-cad-zoom="fit">适应窗口</button>' +
        '<button type="button" data-hx-cad-zoom="100">100%</button>' +
        '<button type="button" data-hx-cad-zoom="area">局部</button>' +
      '</div>';
    var svgEl = pane.querySelector('svg');
    var vb = svgEl && svgEl.viewBox && svgEl.viewBox.baseVal;
    var uw = vb && vb.width > 1 ? vb.width : 240;
    var uh = vb && vb.height > 1 ? vb.height : 160;
    pane._cad = { sc: 1, tx: 0, ty: 0, baseW: uw * CAD_PX, baseH: uh * CAD_PX, user: false, marquee: false };
    function fitReady(tries) {
      var view = pane.querySelector('.hx-cad-view');
      if (view && view.clientWidth > 40) {
        cadFit(pane);
        return;
      }
      if (tries < 8) requestAnimationFrame(function () { fitReady(tries + 1); });
    }
    fitReady(0);
  }

  function loadCadPane(pane, code) {
    if (!pane) return;
    if (pane.getAttribute('data-loaded') === code) return;
    pane.setAttribute('data-loaded', code);
    pane.innerHTML = '<div class="hx-pt">正在打开 ' + esc(code) + ' 的 CAD…</div>';
    fetch('/hx-dcim/cad/' + encodeURIComponent(code) + '.svg?v=20261107merge')
      .then(function (r) {
        var ct = (r.headers.get('content-type') || '').toLowerCase();
        if (!r.ok || ct.indexOf('svg') < 0) throw new Error('missing');
        return r.text();
      })
      .then(function (svg) {
        if (String(svg).indexOf('<svg') < 0) throw new Error('missing');
        mountCadSvg(pane, code, svg);
      })
      .catch(function () {
        pane.innerHTML = '<div class="hx-pt">这间配电室还没有单独的 CAD 原图。<br>目前嵌进来的是 2A1、2A2、2B1、2B2、3A1、3A2、3B1、3B2、4A1、4A2、4B1、4B2、5A1、5A2、5A3、5A4、5B1、5B2、5B3、5B4。</div>';
      });
  }

  function showMainPane(main, which) {
    if (!main) return;
    var cad = main.querySelector('.hx-cad-pane');
    main.querySelectorAll('[data-hx-pane]').forEach(function (b) {
      b.classList.toggle('is-on', b.getAttribute('data-hx-pane') === which);
    });
    main.classList.toggle('is-cad', which === 'cad');
    if (which === 'cad') {
      if (cad) loadCadPane(cad, cad.getAttribute('data-room') || '');
      return;
    }
    fit();
  }

  function renderLoadedPage(host) {
    var main = host.querySelector('.hx-main');
    if (!page || !w.HXDraw) {
      if (main) main.innerHTML = '<div class="hx-pt">这张画面没有从工程备份里解开。</div>';
      L().allow([]);
      pageRows = [];
      return;
    }
    if (!main) return;
    main.classList.remove('is-cad');
    stage.user = false;
    stage.marquee = false;
    main.innerHTML = w.HXDraw.html(page);
    if (page.kind === 'power') {
      var roomCode = roomCodeOf(page);
      main.insertAdjacentHTML('afterbegin',
        '<div class="hx-view-tabs">' +
          '<button type="button" class="is-on" data-hx-pane="screen">一次图</button>' +
          '<button type="button" data-hx-pane="cad">CAD原图</button>' +
        '</div>' +
        '<div class="hx-cad-pane" hidden data-room="' + esc(roomCode) + '"></div>');
      main.insertAdjacentHTML('beforeend',
        '<div class="hx-zoom-bar">' +
        '<button type="button" data-hx-zoom="in">放大</button>' +
        '<button type="button" data-hx-zoom="out">缩小</button>' +
        '<button type="button" data-hx-zoom="fit">适应窗口</button>' +
        '<button type="button" data-hx-zoom="100">100%</button>' +
        '<button type="button" data-hx-zoom="area">局部</button>' +
        '</div>');
    }
    bindStage(main);
    bindCadControls(main);
    if (!main._hxPane) {
      main._hxPane = 1;
      main.addEventListener('click', function (ev) {
        var paneBtn = ev.target.closest ? ev.target.closest('[data-hx-pane]') : null;
        if (!paneBtn || !main.contains(paneBtn)) return;
        ev.stopPropagation();
        showMainPane(main, paneBtn.getAttribute('data-hx-pane'));
      });
    }
    if (!main._hxBound) {
      main._hxBound = 1;
      bindHover(main);
      main.addEventListener('click', function (ev) {
        var drillBtn = ev.target.closest ? ev.target.closest('[data-hx-drill]') : null;
        if (!drillBtn || !hotKey || !w.HXArch) return;
        ev.stopPropagation();
        var code = hotKey.match(/(\d[AB]\d)/);
        if (!code) return;
        var sug = w.HXArch.suggestFromPower('配电室' + code[1]);
        if (sug && sug.item && drillPick) drillPick(sug.item.resolved.room, sug.item.resolved.cab, sug.item);
      });
    }
    pageRows = L().drivers(page);
    L().allow(pageRows);
    paint();
    L().fetch(pageRows);
    fit();
    requestAnimationFrame(fit);
    if (w.HXEdit && w.HXEdit.afterRender) w.HXEdit.afterRender(host, page);
  }

  function showPage(id, host) {
    pinned = false;
    hotKey = '';
    hoverRows = [];
    hoverAnchor = null;
    var main = host.querySelector('.hx-main');
    if (!id) {
      page = null;
      if (main) main.innerHTML = '<div class="hx-pt">这张画面没有从工程备份里解开。</div>';
      L().allow([]);
      pageRows = [];
      return;
    }
    var cached = pageById(id);
    if (cached) {
      page = prep(cached);
      renderLoadedPage(host);
      return;
    }
    if (String(id).indexOf('cad-') === 0) {
      var code = String(id).slice(4);
      page = prep({ id: id, name: '配电室' + code, kind: 'power', nodes: [] });
      if (!page.nodes || !page.nodes.length) {
        if (main) main.innerHTML = '<div class="hx-pt">这间配电室还没有从 CAD 摆进画面。</div>';
        return;
      }
      renderLoadedPage(host);
      return;
    }
    if (main) main.innerHTML = '<div class="hx-pt">正在解开一次图…</div>';
    ensureSceneLoaded(id).then(function (pg) {
      page = prep(pg);
      renderLoadedPage(host);
    });
  }

  function catalog() {
    var power = [];
    var seen = {};
    scenes.forEach(function (p) {
      if (p.kind !== 'power') return;
      seen[p.name] = 1;
      power.push({ name: p.name, kind: 'power', pageId: p.id });
    });
    Object.keys(powerByName).forEach(function (name) {
      if (seen[name]) return;
      seen[name] = 1;
      power.push({ name: name, kind: 'power', pageId: powerByName[name] });
    });
    if (w.HXRooms && w.HXRooms.list) {
      w.HXRooms.list().forEach(function (room) {
        var name = room.name || ('配电室' + room.code);
        if (!room.code || seen[name]) return;
        seen[name] = 1;
        power.push({ name: name, kind: 'power', pageId: 'cad-' + room.code });
      });
    }
    power.sort(function (a, b) {
      if (a.name === '配电室5A1') return -1;
      if (b.name === '配电室5A1') return 1;
      return a.name < b.name ? -1 : 1;
    });
    var live = L().tree ? L().tree() : { modules: [], dataRooms: [] };
    /* 列头按各自测点现排。2A1 的手绘页是整页复制，不能拿去当别的列头。 */
    var modules = live.modules || [];
    if (!modules.length) {
      var groups = {};
      scenes.forEach(function (p) {
        if (p.kind !== 'column') return;
        var g = groups[p.group] || (groups[p.group] = { name: p.group, kind: 'column', cabs: [] });
        g.cabs.push({ name: p.cab || p.name, pageId: p.id });
      });
      modules = Object.keys(groups).map(function (k) { return groups[k]; });
    }
    return [
      { name: '配电室', rooms: power },
      { name: '机房模块', rooms: modules },
      { name: '数据机房', rooms: live.dataRooms || [] }
    ];
  }

  function openRoomCab(host, room, cab) {
    if (!room) return;
    var id = cab && cab.pageId ? cab.pageId : room.pageId;
    if (!id && room.cabs && room.cabs[0]) id = room.cabs[0].pageId;
    showPage(id, host);
  }

  function mountFlatTree(host, tree, mode) {
    var folders = catalog().filter(function (f) {
      if (!f.rooms.length) return false;
      return mode === 'room' ? f.name !== '配电室' : true;
    });
    var flat = [];
    tree.innerHTML = '<div class="hx-folder">工程树</div>' + folders.map(function (f) {
      if (!f.rooms.length) return '';
      return '<div class="hx-folder">' + esc(f.name) + '</div>' + f.rooms.map(function (r) {
        var i = flat.length;
        flat.push(r);
        var head = '<button type="button" class="hx-node" data-i="' + i + '">' + esc(r.name) + '</button>';
        var subs = (r.cabs || []).map(function (c, ci) {
          return '<button type="button" class="hx-node hx-sub" data-i="' + i + '" data-c="' + ci + '">' + esc(c.name) + '</button>';
        }).join('');
        return head + subs;
      }).join('');
    }).join('');
    function open(i, ci) {
      var room = flat[i];
      if (!room) return;
      tree.querySelectorAll('.hx-node').forEach(function (b) {
        var on = Number(b.getAttribute('data-i')) === i && (ci == null ? !b.hasAttribute('data-c') : b.getAttribute('data-c') === String(ci));
        b.classList.toggle('is-on', on);
      });
      var cab = ci == null ? null : (room.cabs || [])[ci];
      openRoomCab(host, room, cab);
    }
    tree.querySelectorAll('.hx-node').forEach(function (b) {
      b.addEventListener('click', function () {
        var ci = b.hasAttribute('data-c') ? Number(b.getAttribute('data-c')) : null;
        if (ci == null) {
          var room = flat[Number(b.getAttribute('data-i'))];
          if (room && room.cabs && room.cabs.length) ci = 0;
        }
        open(Number(b.getAttribute('data-i')), ci);
      });
    });
    var prefer = mode === 'room' ? '机房模块2A1' : '配电室5A1';
    var idx = 0;
    var cabIdx = mode === 'room' ? 0 : null;
    flat.forEach(function (r, i) { if (r.name === prefer) idx = i; });
    if (mode === 'room' && flat[idx] && flat[idx].cabs) {
      flat[idx].cabs.forEach(function (c, i) { if (c.name === 'A列头') cabIdx = i; });
    }
    open(idx, cabIdx);
  }

  function mount(host, mode) {
    if (!host) return;
    hostEl = host;
    treeFlat = null;
    drillPick = function (room, cab, item) { applyPick(room, cab, item); };
    host.innerHTML =
      '<div class="hx-wrap">' +
        '<div class="hx-tree-slot" id="hxTreeSlot">' +
          '<div class="hx-tree-bar"><span>架构树</span><span class="hx-tree-acts">' +
            '<button type="button" data-hx-tree="min" title="最小化" aria-label="最小化菜单">−</button>' +
            '<button type="button" data-hx-tree="hide" title="隐藏" aria-label="隐藏菜单">×</button>' +
          '</span></div>' +
          '<aside class="hx-tree" id="hxTree"></aside>' +
          '<button type="button" class="hx-tree-nub" data-hx-tree="show" title="展开架构树" aria-label="展开架构树"></button>' +
        '</div>' +
        '<div class="hx-body">' +
          '<div class="hx-drill-bar">' +
            '<button type="button" class="hx-tree-dock" data-hx-tree="show" hidden>架构树</button>' +
            '<div class="hx-crumb" id="hxCrumb"></div>' +
            '<span class="hx-drill-hint" id="hxDrillHint">已解开一次图正常显示；灰显配电室可下钻关联机房</span>' +
            '<div class="hx-drill-links" id="hxDrillLinks"></div>' +
          '</div>' +
          '<section class="hx-main" id="hxMain"><div class="hx-pt">正在摆画面</div></section>' +
        '</div>' +
      '</div>';
    var tree = host.querySelector('#hxTree');
    var arch = w.HXArch;
    if (arch && arch.renderTree && arch.branches && arch.branches.length) {
      var pickMode = mode === 'room' ? 'room' : (mode === 'power' ? 'power' : (mode === 'all' || mode === 'drill' ? 'drill' : ''));
      var preferRoom = mode === 'room' ? '机房模块2A1' : '配电室5A1';
      var preferCab = mode === 'room' ? 'A列头' : null;
      treeFlat = arch.renderTree(tree, {
        mode: pickMode,
        active: drillActive,
        preferRoom: preferRoom,
        onPick: function (room, cab, item) { applyPick(room, cab, item); }
      });
      var pi = arch.findPrefer(treeFlat.flat, preferRoom, preferCab);
      if (pi >= 0) {
        var item = treeFlat.flat[pi];
        arch.activateFlatItem(treeFlat, pi);
        applyPick(item.resolved.room, item.resolved.cab, item);
      } else {
        setCrumb([]);
      }
    } else if (scenes.length) {
      mountFlatTree(host, tree, mode);
    } else {
      tree.innerHTML = '<div class="hx-pt">还没有解开航信画的一次图；请从架构树选择已绑点的机房列。</div>';
      setCrumb([]);
    }
    if (onResize) window.removeEventListener('resize', onResize);
    onResize = function () { fit(); };
    window.addEventListener('resize', onResize);
    var wrap = host.querySelector('.hx-wrap');
    if (wrap && !wrap._hxTreeChrome) {
      wrap._hxTreeChrome = 1;
      wrap.addEventListener('click', function (ev) {
        var b = ev.target.closest ? ev.target.closest('[data-hx-tree]') : null;
        if (!b || !wrap.contains(b)) return;
        ev.preventDefault();
        ev.stopPropagation();
        var act = b.getAttribute('data-hx-tree');
        if (act === 'min' || act === 'hide') treeMode = 'peek';
        else if (act === 'show') treeMode = 'full';
        applyTreeMode();
        fit();
      });
    }
    applyTreeMode();
  }

  function applyTreeMode() {
    var slot = hostEl && hostEl.querySelector('#hxTreeSlot');
    var dock = hostEl && hostEl.querySelector('.hx-tree-dock');
    if (slot) {
      slot.hidden = false;
      slot.classList.toggle('is-peek', treeMode === 'peek');
    }
    if (dock) dock.hidden = true;
  }

  function bootBrand() {
    var nameEl = document.getElementById('hxBrandName');
    var logo = document.getElementById('hxBrandLogo');
    if (logo) logo.src = '/static/branding/logo-xunan-hexagon.png';
    fetch('/api/GetAuthLicenseInfo', { headers: L().headers() }).then(function (r) { return r.json(); }).then(function (auth) {
      var raw = ((auth && auth.systemName) || '') + ' ' + ((auth && auth.SystemAPPName) || '');
      var legacy = /零界|zerobound|临界|北京公安局|数据中心动环/.test(raw);
      var title = !auth || legacy || !auth.systemName ? '监控大屏' : auth.systemName;
      if (nameEl) nameEl.textContent = title;
      document.title = title;
      if (logo && auth && auth.SystemLogo && String(auth.SystemLogo).indexOf('data:image') !== 0 && !legacy) {
        logo.src = auth.SystemLogo;
      }
    }).catch(function () {
      if (nameEl) nameEl.textContent = '监控大屏';
    });
  }

  function leave() {
    L().clear();
    pageRows = [];
    hoverRows = [];
    hotKey = '';
    pinned = false;
    hoverAnchor = null;
    page = null;
    if (onResize) window.removeEventListener('resize', onResize);
    onResize = null;
  }

  function openIn(host, room, cab) {
    hostEl = host;
    if (!host) return;
    if (!host.querySelector('.hx-main')) {
      host.innerHTML = '<div class="hx-wrap hx-wrap-solo"><section class="hx-main" id="hxMain"></section></div>';
    }
    var id = cab && cab.pageId ? cab.pageId : (room && room.pageId);
    if (!id && room && room.cabs && room.cabs[0]) id = room.cabs[0].pageId;
    showPage(id, host);
    if (onResize) window.removeEventListener('resize', onResize);
    onResize = function () { fit(); };
    window.addEventListener('resize', onResize);
  }

  function withTimeout(p, ms, label) {
    return Promise.race([
      p,
      new Promise(function (resolve) {
        setTimeout(function () {
          console.warn('[HXScreen] wait timeout: ' + label);
          if (w.HXBoot && w.HXBoot.warn) {
            w.HXBoot.warn(label + ' 等待超时，先展示画面（数据可能不完整）。');
          }
          resolve();
        }, ms);
      })
    ]);
  }

  L().onChange = paint;

  if (w.HXEdit && w.HXEdit.bindHost) {
    w.HXEdit.bindHost({
      dismissHover: function () {
        pinned = false;
        closeHover();
      }
    });
  }

  w.HXScreen = {
    reloadCurrent: function () {
      if (!page || !hostEl) return;
      page._prepped = 0;
      page._layRoom = 0;
      page._cadLinesFixed = 0;
      page._mainParallelFixed = 0;
      page = prep(page);
      renderLoadedPage(hostEl);
    },
    ready: Promise.all([
      withTimeout(L().ready, 130000, '测点索引'),
      withTimeout(loadScenes(), 45000, '一次图场景'),
      withTimeout(loadPowerIndex(), 20000, '配电页索引'),
      withTimeout((w.HXArch && w.HXArch.ready) || Promise.resolve(), 45000, '架构树')
    ]).then(function () {
      return syncPowerMetaFromApi();
    }).then(function () {
      if (w.HXBoot && w.HXBoot.hide) w.HXBoot.hide();
      L().connect();
    }).catch(function (err) {
      if (w.HXBoot && w.HXBoot.warn) w.HXBoot.warn((err && err.message) || '初始化异常');
      if (w.HXBoot && w.HXBoot.hide) w.HXBoot.hide();
    }),
    mount: mount,
    leave: leave,
    bootBrand: bootBrand,
    catalog: catalog,
    openIn: openIn,
    applyPick: applyPick,
    setCrumb: setCrumb,
    findPageByName: findPageByName
  };
})(window);
