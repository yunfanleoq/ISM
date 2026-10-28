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
    powerIndexP = fetch('/hx-dcim/data/hx-power-page-index.json?v=20261028cad76').then(function (r) {
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
      jobs.push(w.HX5A1.load('20261028cad76').catch(function (e) {
        if (w.console) w.console.warn('[hx-5a1] 布局未加载，5A1 退回通用绘制', e);
      }));
    }
    if (w.HXRooms) {
      jobs.push(w.HXRooms.load('20261028cad76').catch(function (e) {
        if (w.console) w.console.warn('[hx-rooms] 配电室布局未加载', e);
      }));
    }
    return Promise.all(jobs).then(function () {
      return fetch('/hx-dcim/data/hx-scenes.json?v=20261028cad76');
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

  function paint() {
    if (!page || !hostEl) return;
    hostEl.querySelectorAll('[data-hx-line]').forEach(function (el) {
      var n = page.nodes[Number(el.getAttribute('data-i'))];
      if (n && w.HXDraw && w.HXDraw.lineColor) el.setAttribute('stroke', w.HXDraw.lineColor(page, n));
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
      var v = L().nodeValue(n);
      var closed = false;
      if (n.op) closed = v != null && v !== '' && L().matchOp(v, n.op, n.val);
      else closed = v != null && v !== '' && L().matchOp(v, '==', '1');
      el.classList.toggle('is-closed', closed);
      el.classList.toggle('is-open', !closed);
      el.style.color = closed ? '#ff0000' : '#00ff00';
    });
    hostEl.querySelectorAll('[data-hx-uuid]').forEach(function (el) {
      var known = L().uuidValue(el.getAttribute('data-hx-uuid'));
      if (known) el.textContent = L().showValue(known);
    });
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
    return '<div class="hx-hover-title">' + esc(labelOf(key)) + '</div>' + body + powerDrillBtn(key);
  }

  function placeHover(box, ev) {
    if (!box || !ev) return;
    var x = ev.clientX + 14;
    var y = ev.clientY + 14;
    if (x + 440 > window.innerWidth) x = Math.max(8, ev.clientX - 440);
    if (y + 340 > window.innerHeight) y = Math.max(8, ev.clientY - 280);
    box.style.left = x + 'px';
    box.style.top = y + 'px';
  }

  function openHover(key, ev) {
    hoverRows = L().pointsFor(key, devOf(key));
    L().allow(pageRows.concat(hoverRows));
    var box = document.getElementById('hxHover');
    if (box) {
      box.innerHTML = hoverHtml(key, hoverRows);
      box.hidden = false;
      placeHover(box, ev);
    }
    L().fetch(hoverRows).then(function () {
      if (box && !box.hidden && hotKey === key) box.innerHTML = hoverHtml(key, hoverRows);
    });
  }

  function closeHover() {
    if (pinned) return;
    hotKey = '';
    hoverRows = [];
    var box = document.getElementById('hxHover');
    if (box) box.hidden = true;
    L().allow(pageRows);
  }

  function editOn() {
    return !!(w.HXEdit && w.HXEdit.blocksStage && w.HXEdit.blocksStage());
  }

  function bindHover(host) {
    host.addEventListener('mouseover', function (ev) {
      if (editOn()) return;
      if (pinned) return;
      var el = ev.target.closest ? ev.target.closest('[data-hx-key]') : null;
      if (!el) return;
      var key = el.getAttribute('data-hx-key');
      if (!key || key === hotKey) return;
      hotKey = key;
      openHover(key, ev);
    });
    host.addEventListener('mousemove', function (ev) {
      var box = document.getElementById('hxHover');
      if (box && !box.hidden && !pinned) placeHover(box, ev);
    });
    host.addEventListener('mouseout', function (ev) {
      var to = ev.relatedTarget;
      if (to && to.closest && (to.closest('[data-hx-key]') || to.closest('#hxHover'))) return;
      closeHover();
    });
    host.addEventListener('click', function (ev) {
      if (editOn()) return;
      var el = ev.target.closest ? ev.target.closest('[data-hx-key]') : null;
      if (!el) {
        pinned = false;
        hotKey = '';
        hoverRows = [];
        var box = document.getElementById('hxHover');
        if (box) box.hidden = true;
        L().allow(pageRows);
        return;
      }
      pinned = true;
      hotKey = el.getAttribute('data-hx-key');
      openHover(hotKey, ev);
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

  function applyStage(tx, ty, sc) {
    var board = stageBoard();
    if (!board) return;
    stage.tx = tx;
    stage.ty = ty;
    stage.sc = sc;
    board.style.transformOrigin = '0 0';
    board.style.transform = 'translate(' + tx + 'px,' + ty + 'px) scale(' + sc + ')';
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
    var box = page && page.kind === 'power' ? powerContentBox(page) : null;
    var x0 = (page && page.x0) || 0;
    var ox = 0;
    var oy = 0;
    var bw = board.offsetWidth;
    var bh = board.offsetHeight;
    if (box) {
      ox = box.minX - x0;
      oy = box.minY;
      bw = box.maxX - box.minX;
      bh = box.maxY - box.minY;
    }
    if (!(bw > 1) || !(bh > 1)) return;
    var pad = Math.max(16, Math.round(Math.min(cw, ch) * 0.03));
    var availW = Math.max(1, cw - pad * 2);
    var availH = Math.max(1, ch - pad * 2);
    var sW = availW / bw;
    var sH = availH / bh;
    var s = Math.min(sW, sH);
    if (!isFinite(s) || s <= 0) s = 0.4;
    /* 扁长图：整图塞进窗口后画高不到视口 55%，改按高度铺满，多出来的宽度平移。
       1A2 约 12000×1300，按宽度适配只剩一条带子。接近视口比例的室（如 5A1）仍整图放进。 */
    var viewAspect = cw / ch;
    var contentAspect = bw / bh;
    if (contentAspect > viewAspect * 2.5 && bh * s < ch * 0.55) s = sH;
    s = Math.max(0.05, Math.min(8, s));
    var contentW = bw * s;
    var contentH = bh * s;
    var dx = contentW <= cw - pad ? (cw - contentW) / 2 - ox * s : pad - ox * s;
    var dy = contentH <= ch - pad ? (ch - contentH) / 2 - oy * s : pad - oy * s;
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
    main.addEventListener('wheel', function (ev) {
      if (editOn()) return;
      if (!page || page.kind !== 'power' || !blank(ev.target) && !(ev.target.closest && ev.target.closest('.hx-fit'))) return;
      if (ev.target.closest && ev.target.closest('.hx-zoom-bar')) return;
      ev.preventDefault();
      var dy = ev.deltaY;
      if (ev.deltaMode === 1) dy *= 16;
      zoomAt(ev.clientX, ev.clientY, (stage.sc || 1) * (dy > 0 ? 0.9 : 1.12));
    }, { passive: false });
    main.addEventListener('dblclick', function (ev) {
      if (editOn()) return;
      if (!page || page.kind !== 'power') return;
      if (ev.target.closest && ev.target.closest('.hx-zoom-bar')) return;
      zoomAt(ev.clientX, ev.clientY, (stage.sc || 1) * (ev.shiftKey ? 1 / 1.35 : 1.35));
    });
    main.addEventListener('pointerdown', function (ev) {
      if (editOn()) return;
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
      rh = box.h * 0.92;
      s = rh / BLADE_VH;
      rw = BLADE_VW * s;
      ox = (box.w - rw) / 2;
      oy = (box.h - rh) / 2;
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

  function renderLoadedPage(host) {
    var main = host.querySelector('.hx-main');
    if (!page || !w.HXDraw) {
      if (main) main.innerHTML = '<div class="hx-pt">这张画面没有从工程备份里解开。</div>';
      L().allow([]);
      pageRows = [];
      return;
    }
    if (!main) return;
    stage.user = false;
    stage.marquee = false;
    main.innerHTML = w.HXDraw.html(page);
    if (page.kind === 'power') {
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
          '<div class="hx-tree-bar"><span>配电室</span><span class="hx-tree-acts">' +
            '<button type="button" data-hx-tree="min" title="最小化" aria-label="最小化菜单">−</button>' +
            '<button type="button" data-hx-tree="hide" title="隐藏" aria-label="隐藏菜单">×</button>' +
          '</span></div>' +
          '<aside class="hx-tree" id="hxTree"></aside>' +
          '<button type="button" class="hx-tree-nub" data-hx-tree="show" title="展开配电室" aria-label="展开配电室"></button>' +
        '</div>' +
        '<div class="hx-body">' +
          '<div class="hx-drill-bar">' +
            '<button type="button" class="hx-tree-dock" data-hx-tree="show" hidden>配电室</button>' +
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
      treeFlat = arch.renderTree(tree, {
        mode: pickMode,
        active: drillActive,
        onPick: function (room, cab, item) { applyPick(room, cab, item); }
      });
      var preferRoom = mode === 'room' ? '机房模块2A1' : '配电室5A1';
      var preferCab = mode === 'room' ? 'A列头' : null;
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
    setCrumb: setCrumb
  };
})(window);
