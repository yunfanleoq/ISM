/* 航信一次图和列头柜画面。布局来自他们画的工程，取数走 HXLive。 */
(function (w) {
  'use strict';

  var L = function () { return w.HXLive; };
  var scenes = [];
  var page = null;
  var pageRows = [];
  var hoverRows = [];
  var hotKey = '';
  var pinned = false;
  var hostEl = null;
  var onResize = null;

  function esc(s) {
    return String(s == null ? '' : s)
      .replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
  }

  function loadScenes() {
    return fetch('/hx-dcim/data/hx-scenes.json?v=20261024draw').then(function (r) {
      if (!r.ok) throw new Error('no scenes');
      return r.json();
    }).then(function (j) {
      scenes = (j && j.pages) || [];
    }).catch(function () { scenes = []; });
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

  function hoverHtml(key, rows) {
    var body = rows.length ? rows.map(function (p) {
      var known = L().remembered(p) || p;
      return '<div class="hx-pt"><span>' + esc(p.name) + '</span><b data-hx-uuid="' + esc(p.uuid) + '">' + esc(L().showValue(known)) + '</b></div>';
    }).join('') : '<div class="hx-pt"><span>这一回路还没有测点</span></div>';
    return '<div class="hx-hover-title">' + esc(labelOf(key)) + '</div>' + body;
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

  function bindHover(host) {
    host.addEventListener('mouseover', function (ev) {
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

  function fit() {
    var main = hostEl && hostEl.querySelector('.hx-main');
    var fitEl = hostEl && hostEl.querySelector('.hx-fit');
    var board = hostEl && hostEl.querySelector('.hx-board');
    if (!main || !fitEl || !board) return;
    var cw = main.clientWidth;
    var ch = main.clientHeight;
    var bw = board.offsetWidth;
    var bh = board.offsetHeight;
    /* 外层铺满可视区域，元件仍按原图宽高比缩放，避免只按宽度缩完后底部空出一条别的颜色。 */
    fitEl.style.width = cw + 'px';
    fitEl.style.height = ch + 'px';
    fitEl.style.transform = 'none';
    var s = Math.min(cw / bw, ch / bh);
    if (!isFinite(s) || s <= 0) s = 0.4;
    var spare = ch - bh * s;
    /* 底部空档大于一截时，只要左右裁掉的还在画布空白里，就按高度铺满。 */
    if (spare > 12) {
      var cover = Math.max(cw / bw, ch / bh);
      var crop = (bw * cover - cw) / 2 / cover;
      if (crop < bw * 0.03) s = cover;
    }
    var dx = (cw - bw * s) / 2;
    board.style.transformOrigin = '0 0';
    board.style.transform = 'translate(' + dx + 'px,0) scale(' + s + ')';
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

  function prep(pg) {
    if (!pg || pg._prepped) return pg;
    (pg.nodes || []).forEach(function (n) {
      if (n.t === 'T' && !n.k) {
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

  function showPage(id, host) {
    page = prep(pageById(id));
    pinned = false;
    hotKey = '';
    hoverRows = [];
    var main = host.querySelector('.hx-main');
    if (!page || !w.HXDraw) {
      if (main) main.innerHTML = '<div class="hx-pt">这张画面没有从工程备份里解开。</div>';
      L().allow([]);
      pageRows = [];
      return;
    }
    if (!main) return;
    main.innerHTML = w.HXDraw.html(page);
    if (!main._hxBound) {
      main._hxBound = 1;
      bindHover(main);
    }
    pageRows = L().drivers(page);
    L().allow(pageRows);
    paint();
    L().fetch(pageRows);
    fit();
    requestAnimationFrame(fit);
  }

  function catalog() {
    var power = [];
    scenes.forEach(function (p) {
      if (p.kind === 'power') power.push({ name: p.name, kind: 'power', pageId: p.id });
    });
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

  function mount(host, mode) {
    if (!host) return;
    hostEl = host;
    var folders = catalog().filter(function (f) {
      if (!f.rooms.length) return false;
      return mode === 'room' ? f.name !== '配电室' : true;
    });
    host.innerHTML = '<div class="hx-wrap"><aside class="hx-tree" id="hxTree"></aside><section class="hx-main" id="hxMain"><div class="hx-pt">正在摆画面</div></section></div>';
    var tree = host.querySelector('#hxTree');
    var flat = [];
    if (!scenes.length) {
      tree.innerHTML = '<div class="hx-pt">还没有解开航信画的一次图。</div>';
      return;
    }
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
      var id = cab && cab.pageId ? cab.pageId : room.pageId;
      if (!id && room.cabs && room.cabs[0]) id = room.cabs[0].pageId;
      showPage(id, host);
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
    if (onResize) window.removeEventListener('resize', onResize);
    onResize = function () { fit(); };
    window.addEventListener('resize', onResize);
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

  L().onChange = paint;

  w.HXScreen = {
    ready: Promise.all([L().ready, loadScenes()]).then(function () { L().connect(); }),
    mount: mount,
    leave: leave,
    bootBrand: bootBrand,
    catalog: catalog,
    openIn: openIn
  };
})(window);
