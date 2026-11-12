/* 一次图编辑。拖动先记在 localStorage（hx-dcim-edit:<室号>）。点保存才写回这间一次图 JSON。 */
(function (w) {
  'use strict';

  var VER = '20261102f5';
  var inboxSeen = '';
  var palMode = 'full';
  var editing = false;
  var saving = false;
  var host = null;
  var page = null;
  var room = '';
  var selected = null;
  var dismissHover = null;
  var gestures = null;

  var BLADE = '<svg class="hx-blade" viewBox="-6 0 48 124" preserveAspectRatio="xMidYMid meet"><g class="hx-k-closed" fill="none" stroke="currentColor" stroke-width="2.15" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="0" x2="18" y2="114"/><path d="M30 52 A13.2 13.2 0 0 1 30 78"/><circle cx="18" cy="14" r="4.15" fill="currentColor" stroke="none"/><circle cx="18" cy="44" r="4.15" fill="currentColor" stroke="none"/><circle cx="18" cy="86" r="4.15" fill="currentColor" stroke="none"/></g><g class="hx-k-open" fill="none" stroke="currentColor" stroke-width="2.15" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="0" x2="18" y2="44"/><line x1="18" y1="86" x2="18" y2="114"/><path d="M30 52 A13.2 13.2 0 0 1 30 78"/><line x1="18" y1="44" x2="39.51" y2="70.92"/><circle cx="18" cy="14" r="4.15" fill="currentColor" stroke="none"/><circle cx="18" cy="44" r="4.15" fill="currentColor" stroke="none"/><circle cx="18" cy="86" r="4.15" fill="currentColor" stroke="none"/></g><circle cx="18" cy="114" r="4.2" fill="#ffffff"/></svg>';
  var XFMR = '<svg viewBox="0 0 64 88" width="100%" height="100%" aria-hidden="true"><g fill="none" stroke="#e8eef6" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"><circle cx="26" cy="16" r="15"/><polygon points="26,6 15,24 37,24"/><circle cx="26" cy="46" r="15"/><path d="M26 46 V34 M26 46 L16 56 M26 46 L36 56"/><path d="M26 46 H48 V64 M42 64 H54 M44 68 H52 M46 72 H50"/></g></svg>';
  var CAP = '<svg viewBox="0 0 120 252" width="100%" height="100%" overflow="visible" aria-hidden="true"><g fill="none" stroke="#e8eef6" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"><g transform="translate(41.032,0) scale(0.790323) translate(6,0)" stroke="#00ff00" stroke-width="2.15"><line x1="18" y1="0" x2="18" y2="44" stroke-linecap="butt"/><line x1="18" y1="86" x2="18" y2="114" stroke-linecap="butt"/><path d="M30 52 A13.2 13.2 0 0 1 30 78"/><line x1="18" y1="44" x2="39.51" y2="70.92" stroke-linecap="butt"/><circle cx="18" cy="14" r="4.15" fill="#00ff00" stroke="none"/><circle cx="18" cy="44" r="4.15" fill="#00ff00" stroke="none"/><circle cx="18" cy="86" r="4.15" fill="#00ff00" stroke="none"/><circle cx="18" cy="114" r="4.2" fill="#ffffff" stroke="none"/></g><path d="M60 90.1 V106"/><circle cx="44" cy="116" r="6.2"/><circle cx="60" cy="116" r="6.2"/><circle cx="76" cy="116" r="6.2"/><path d="M60 106 V132 M16 132 H104"/><path d="M30 132 V140 L24 148 L36 156 L24 164 L36 172 L30 180 V188"/><path d="M90 132 V140 L96 148 L84 156 L96 164 L84 172 L90 180 V188"/><path d="M30 188 V198 M18 216 L30 198 L42 216 Z M16 222 H44"/><path d="M90 188 V198 M78 216 L90 198 L102 216 Z M76 222 H104"/><path d="M60 132 V158 M50 158 H70 M50 165 H70 M60 165 V184"/><path d="M46 184 H74 M50 189 H70 M54 194 H66"/><circle cx="54" cy="236" r="1.7" fill="#e8eef6" stroke="none"/><circle cx="60" cy="236" r="1.7" fill="#e8eef6" stroke="none"/><circle cx="66" cy="236" r="1.7" fill="#e8eef6" stroke="none"/></g></svg>';

  var UPS = '<svg viewBox="0 0 72 48" width="100%" height="100%" aria-hidden="true"><g fill="none" stroke="#d24dff" stroke-width="1.6" stroke-linecap="round"><rect x="2" y="2" width="68" height="44" rx="3"/><path d="M12 16 Q22 8 36 16 T60 16"/><path d="M14 32 H58"/></g></svg>';
  var INV = '<svg viewBox="0 0 72 48" width="100%" height="100%" aria-hidden="true"><g fill="none" stroke="#d24dff" stroke-width="1.6" stroke-linecap="round"><rect x="2" y="2" width="68" height="44" rx="3"/><path d="M12 14 H60"/><path d="M12 34 Q22 26 36 34 T60 34"/></g></svg>';
  var SSW = '<svg viewBox="0 0 64 52" width="100%" height="100%" aria-hidden="true"><g fill="none" stroke="#00e676" stroke-width="1.5" stroke-linejoin="round" stroke-linecap="round"><rect x="2" y="2" width="60" height="48" rx="3"/><polygon points="14,14 32,26 14,38"/><polygon points="50,14 32,26 50,38"/><path d="M32 26 H46 M18 26 H32"/></g></svg>';

  var GND = '<svg viewBox="0 0 48 36" width="100%" height="100%" aria-hidden="true"><g fill="none" stroke="#e8eef6" stroke-width="1.8" stroke-linecap="round"><path d="M24 2 V14"/><path d="M8 14 H40"/><path d="M13 21 H35"/><path d="M18 28 H30"/></g></svg>';
  var GEN = '<svg viewBox="0 0 64 64" width="100%" height="100%" aria-hidden="true"><g fill="none" stroke="#ffb020" stroke-width="1.8" stroke-linecap="round"><circle cx="32" cy="32" r="22"/><path d="M16 32 H22 Q28 18 32 32 T42 32 H48"/></g></svg>';

  var GROUPS = [
    { name: '开关', tools: [
      ['knife', '竖刀', '<path d="M12 1 V16 M12 24 V35" stroke="#00ff00" stroke-width="1.7" fill="none"/><path d="M16 18 A7 7 0 0 1 16 28" stroke="#00ff00" stroke-width="1.7" fill="none"/>']
    ]},
    { name: '电源', tools: [
      ['xfmr', '变压器', '<circle cx="10" cy="9" r="5" fill="none" stroke="#e8eef6" stroke-width="1.4"/><circle cx="10" cy="20" r="5" fill="none" stroke="#e8eef6" stroke-width="1.4"/>'],
      ['gen', '柴油发电机', '<circle cx="12" cy="16" r="7" fill="none" stroke="#ffb020" stroke-width="1.4"/><path d="M8 16 H16" stroke="#ffb020" stroke-width="1.3"/>'],
      ['ups', 'UPS', '<rect x="3" y="8" width="18" height="14" rx="1.5" fill="none" stroke="#d24dff" stroke-width="1.3"/><path d="M6 13 H18" stroke="#d24dff" stroke-width="1.2"/>'],
      ['inv', '逆变器', '<rect x="3" y="8" width="18" height="14" rx="1.5" fill="none" stroke="#d24dff" stroke-width="1.3"/><path d="M6 12 Q12 16 18 12" fill="none" stroke="#d24dff" stroke-width="1.2"/>']
    ]},
    { name: '补偿与接地', tools: [
      ['cap', '电容补偿柜', '<path d="M12 1 V7" stroke="#00ff00" stroke-width="1.4"/><circle cx="8" cy="11" r="2" fill="none" stroke="#e8eef6"/><circle cx="12" cy="11" r="2" fill="none" stroke="#e8eef6"/><circle cx="16" cy="11" r="2" fill="none" stroke="#e8eef6"/><path d="M5 18 H19 M6 22 H18 M7 26 H17" stroke="#e8eef6" stroke-width="1.3" fill="none"/>'],
      ['gnd', '接地', '<path d="M12 4 V12 M4 12 H20 M7 17 H17 M10 22 H14" stroke="#e8eef6" stroke-width="1.4" fill="none"/>']
    ]},
    { name: '切换', tools: [
      ['ssw', '静态开关', '<rect x="3" y="8" width="18" height="16" rx="1.5" fill="none" stroke="#00e676" stroke-width="1.3"/><path d="M8 12 L12 16 L8 20 M16 12 L12 16 L16 20" fill="none" stroke="#00e676" stroke-width="1.2"/>']
    ]},
    { name: '标注', tools: [
      ['text', '文字', '<path d="M4 6 H20 M12 6 V20" stroke="#fff" stroke-width="1.6" fill="none"/>'],
      ['custom', '自定义', '<rect x="3" y="8" width="18" height="16" rx="1.5" fill="none" stroke="#7fd0ff" stroke-width="1.3"/><path d="M8 16 H16" stroke="#7fd0ff" stroke-width="1.3"/>']
    ]},
    { name: '连线', tools: [
      ['bus', '母线', '<path d="M2 16 H22" stroke="#ce0c18" stroke-width="3.2"/>'],
      ['hline', '横线', '<path d="M2 16 H22" stroke="#ce0c18" stroke-width="2"/>'],
      ['vline', '竖线', '<path d="M12 2 V34" stroke="#ce0c18" stroke-width="2"/>'],
      ['ortho', '智能折线', '<path d="M2 28 H12 V8 H22" fill="none" stroke="#ce0c18" stroke-width="2"/>']
    ]}
  ];

  function esc(s) {
    return String(s == null ? '' : s)
      .replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
  }

  function r1(n) { return Math.round(Number(n) * 10) / 10; }

  function roomOf(pg) {
    var name = (pg && pg.name) || '';
    var m = String(name).match(/配电室\s*([0-9A-Za-z]+)/);
    if (m) return m[1];
    if (pg && pg.id) return String(pg.id);
    return 'room';
  }

  function key() { return 'hx-dcim-edit:' + roomOf(page); }

  function emptyOv() { return { v: 1, moves: {}, deleted: [], items: [] }; }

  function load() {
    var empty = emptyOv();
    try {
      var raw = w.localStorage.getItem(key());
      if (!raw) return empty;
      var j = JSON.parse(raw);
      if (!j || typeof j !== 'object') return empty;
      j.moves = j.moves || {};
      j.deleted = (j.deleted || []).map(String);
      j.items = j.items || [];
      j.v = 1;
      return j;
    } catch (e) { return empty; }
  }

  function commit(fn) {
    var ov = load();
    fn(ov);
    ov.v = 1;
    ov.roomId = roomOf(page);
    try { w.localStorage.setItem(key(), JSON.stringify(ov)); } catch (e) {}
    return ov;
  }

  function shiftEl() { return host && host.querySelector('.hx-shift'); }

  function scale() {
    var shift = shiftEl();
    if (!shift) return 1;
    var r = shift.getBoundingClientRect();
    var s = r.width / (shift.offsetWidth || 1);
    return s > 0.0001 ? s : 1;
  }

  function pageXY(clientX, clientY) {
    var shift = shiftEl();
    if (!shift) return null;
    var r = shift.getBoundingClientRect();
    if (!(r.width > 1) || !(r.height > 1)) return null;
    return {
      x: (clientX - r.left) / r.width * shift.offsetWidth,
      y: (clientY - r.top) / r.height * shift.offsetHeight
    };
  }

  function nid() {
    return 'e' + Date.now().toString(36) + Math.floor(Math.random() * 1296).toString(36);
  }

  function editableNode(n) {
    if (!n) return false;
    if (n.t === 'I' || n.t === 'S' || n.t === 'T' || n.t === 'L') return true;
    return n.t === 'B' && (n.sym === 'xfmr' || n.sym === 'cap');
  }

  /* 411/421 竖母线等布局补线：禁止编辑删除，否则铭牌横抽头悬空 */
  function isProtectedNode(n) {
    if (!n) return false;
    return /spine/i.test(String(n.fix || ''));
  }

  function nodeAt(el) {
    if (!page || !el || !el.getAttribute) return null;
    var i = el.getAttribute('data-i');
    if (i == null || i === '') return null;
    return page.nodes[Number(i)] || null;
  }

  function kindOf(el) {
    if (!el || (el.closest && el.closest('.hx-edit-layer'))) return '';
    if (el.hasAttribute && el.hasAttribute('data-hx-knife')) return 'knife';
    var sym = el.getAttribute && el.getAttribute('data-hx-sym');
    if (sym === 'xfmr' || sym === 'cap') return sym;
    if (el.classList && el.classList.contains('hx-tx')) return 'text';
    var tag = (el.tagName || '').toLowerCase();
    if (tag !== 'line' && tag !== 'path') return '';
    var n = nodeAt(el);
    return n && n.t === 'L' ? 'line' : '';
  }

  function elsFor(i) {
    var shift = shiftEl();
    if (!shift) return [];
    return Array.prototype.filter.call(shift.querySelectorAll('[data-i="' + i + '"]'), function (el) {
      return !(el.closest && el.closest('.hx-edit-layer'));
    });
  }

  function tagSymbols(shift, nodes) {
    var free = [];
    Array.prototype.forEach.call(shift.children, function (el) {
      if (!el.classList || !el.classList.contains('hx-abs')) return;
      if (el.hasAttribute('data-i') || el.classList.contains('hx-tx')) return;
      if (el.classList.contains('hx-edit-layer')) return;
      free.push(el);
    });
    nodes.forEach(function (n, i) {
      if (!n || n.t !== 'B' || (n.sym !== 'xfmr' && n.sym !== 'cap')) return;
      var j, best = -1, bestD = 2;
      for (j = 0; j < free.length; j++) {
        var d = Math.abs(parseFloat(free[j].style.left) - Number(n.x)) + Math.abs(parseFloat(free[j].style.top) - Number(n.y));
        if (d < bestD) { bestD = d; best = j; }
      }
      if (best < 0) return;
      free[best].setAttribute('data-i', String(i));
      free[best].setAttribute('data-hx-sym', n.sym);
      free.splice(best, 1);
    });
  }

  function endsOf(n) {
    if (!n) return [0, 0, 0, 0];
    if (n.x1 != null && n.y1 != null && n.x2 != null && n.y2 != null) {
      return [Number(n.x1), Number(n.y1), Number(n.x2), Number(n.y2)];
    }
    var rad = (Number(n.a) || 0) * Math.PI / 180;
    var cos = Math.cos(rad);
    var sin = Math.sin(rad);
    var cx = Number(n.x) + Number(n.w) / 2;
    var cy = Number(n.y) + Number(n.h) / 2;
    var hw = Number(n.w) / 2;
    return [cx - hw * cos, cy - hw * sin, cx + hw * cos, cy + hw * sin];
  }

  function lineAxis(n) {
    if (!n) return 'h';
    if (n.x1 != null && n.y1 != null && n.x2 != null && n.y2 != null) {
      if (Math.abs(Number(n.y1) - Number(n.y2)) < 1.5) return 'h';
      if (Math.abs(Number(n.x1) - Number(n.x2)) < 1.5) return 'v';
      return 'd';
    }
    if (Math.abs((Number(n.a) || 0) % 180) < 1) return 'h';
    return 'v';
  }

  function lineGeom(n, mv) {
    if (mv && mv.x1 != null) return { x1: Number(mv.x1), y1: Number(mv.y1), x2: Number(mv.x2), y2: Number(mv.y2) };
    var e = endsOf(n);
    var dx = (mv && mv.dx) || 0;
    var dy = (mv && mv.dy) || 0;
    return { x1: e[0] + dx, y1: e[1] + dy, x2: e[2] + dx, y2: e[3] + dy };
  }

  function boxGeom(n, mv) {
    if (mv && mv.w != null && mv.x != null) {
      return { x: Number(mv.x), y: Number(mv.y), w: Number(mv.w), h: Number(mv.h) };
    }
    return {
      x: Number(n.x) + ((mv && mv.dx) || 0),
      y: Number(n.y) + ((mv && mv.dy) || 0),
      w: Number(n.w) || 0,
      h: Number(n.h) || 0
    };
  }

  function resizeLine(g, axis, end, dx, dy) {
    var x1 = g.x1;
    var y1 = g.y1;
    var x2 = g.x2;
    var y2 = g.y2;
    if (end === '1') {
      if (axis !== 'v') x1 += dx;
      if (axis !== 'h') y1 += dy;
    } else {
      if (axis !== 'v') x2 += dx;
      if (axis !== 'h') y2 += dy;
    }
    if (axis === 'h') {
      y1 = g.y1;
      y2 = g.y2;
      if (Math.abs(x2 - x1) < 8) {
        if (end === '1') x1 = x2 + (g.x1 <= g.x2 ? -8 : 8);
        else x2 = x1 + (g.x2 >= g.x1 ? 8 : -8);
      }
    } else if (axis === 'v') {
      x1 = g.x1;
      x2 = g.x2;
      if (Math.abs(y2 - y1) < 8) {
        if (end === '1') y1 = y2 + (g.y1 <= g.y2 ? -8 : 8);
        else y2 = y1 + (g.y2 >= g.y1 ? 8 : -8);
      }
    } else if (Math.hypot(x2 - x1, y2 - y1) < 8) {
      var vx = g.x2 - g.x1;
      var vy = g.y2 - g.y1;
      var len = Math.hypot(vx, vy) || 1;
      if (end === '1') { x1 = x2 - vx / len * 8; y1 = y2 - vy / len * 8; }
      else { x2 = x1 + vx / len * 8; y2 = y1 + vy / len * 8; }
    }
    return { x1: x1, y1: y1, x2: x2, y2: y2 };
  }

  function resizeBox(b, end, dx, dy) {
    var x = b.x;
    var y = b.y;
    var r = b.x + b.w;
    var bt = b.y + b.h;
    if (end.indexOf('w') >= 0) x += dx;
    if (end.indexOf('e') >= 0) r += dx;
    if (end.indexOf('n') >= 0) y += dy;
    if (end.indexOf('s') >= 0) bt += dy;
    if (r - x < 12) {
      if (end.indexOf('w') >= 0) x = r - 12;
      else r = x + 12;
    }
    if (bt - y < 12) {
      if (end.indexOf('n') >= 0) y = bt - 12;
      else bt = y + 12;
    }
    return { x: x, y: y, w: r - x, h: bt - y };
  }

  function paintLineEl(el, x1, y1, x2, y2) {
    el.removeAttribute('transform');
    var tag = (el.tagName || '').toLowerCase();
    if (tag === 'line') {
      el.setAttribute('x1', x1);
      el.setAttribute('y1', y1);
      el.setAttribute('x2', x2);
      el.setAttribute('y2', y2);
      return;
    }
    if (tag === 'path') el.setAttribute('d', 'M ' + x1 + ' ' + y1 + ' L ' + x2 + ' ' + y2);
  }

  function paintTextLook(el, mv, n) {
    if (!el || !mv) return;
    if (mv.fs != null && mv.fs !== '') {
      el.style.fontSize = mv.fs + 'px';
      el.setAttribute('data-hx-fs', String(mv.fs));
    }
    /* 跟闸的字由开合上色，这里不盖。其余文字改色时，字和框用同一个颜色。 */
    if (mv.c && !(n && (n.followKnife || n.upsPlate)) && !(el.classList && el.classList.contains('hx-plate-live'))) {
      el.style.color = mv.c;
      el.style.borderColor = mv.c;
    }
  }

  function placeOne(el, n, mv) {
    var tag = (el.tagName || '').toLowerCase();
    if (tag === 'line' || tag === 'path') {
      if (mv && mv.x1 != null) {
        paintLineEl(el, mv.x1, mv.y1, mv.x2, mv.y2);
        return;
      }
      var dx = (mv && mv.dx) || 0;
      var dy = (mv && mv.dy) || 0;
      if (dx || dy) el.setAttribute('transform', 'translate(' + dx + ' ' + dy + ')');
      else el.removeAttribute('transform');
      return;
    }
    if (mv && mv.w != null && mv.x != null && n && (n.t === 'T' || n.t === 'B')) {
      el.style.left = mv.x + 'px';
      el.style.top = mv.y + 'px';
      el.style.width = mv.w + 'px';
      el.style.height = mv.h + 'px';
    } else {
      var dx2 = (mv && mv.dx) || 0;
      var dy2 = (mv && mv.dy) || 0;
      el.style.left = (Number(n.x) + dx2) + 'px';
      el.style.top = (Number(n.y) + dy2) + 'px';
    }
    if (n && n.t === 'T' && el.classList && el.classList.contains('hx-tx')) paintTextLook(el, mv, n);
  }

  function textOf(el, s) {
    el.innerHTML = esc(s).replace(/\n/g, '<br>');
  }

  function ensureLayer() {
    var shift = shiftEl();
    if (!shift) return null;
    var layer = shift.querySelector('.hx-edit-layer');
    if (!layer) {
      layer = document.createElement('div');
      layer.className = 'hx-edit-layer';
      shift.appendChild(layer);
    }
    layer.style.width = shift.offsetWidth + 'px';
    layer.style.height = shift.offsetHeight + 'px';
    return layer;
  }

  function isWire(kind) {
    return kind === 'hline' || kind === 'vline' || kind === 'ortho' || kind === 'bus';
  }

  function lineSvg(it) {
    var sw = it.kind === 'bus' ? 6 : 3;
    var a = ' data-hx-edit-id="' + it.id + '" x1="' + it.x1 + '" y1="' + it.y1 + '" x2="' + it.x2 + '" y2="' + it.y2 + '"';
    return '<line class="hx-edit-item"' + a + ' stroke="#ce0c18" stroke-width="' + sw + '" stroke-linecap="round"/>' +
      '<line class="hx-edit-hit"' + a + ' stroke="transparent" stroke-width="18" stroke-linecap="round"/>';
  }

  function orthoSvg(it) {
    var C = w.HXEditCmd;
    var d = C.orthoD(C.orthoPts(it));
    var a = ' data-hx-edit-id="' + it.id + '" d="' + d + '"';
    return '<path class="hx-edit-item"' + a + ' fill="none" stroke="#ce0c18" stroke-width="3" stroke-linejoin="round" stroke-linecap="round"/>' +
      '<path class="hx-edit-hit"' + a + ' fill="none" stroke="transparent" stroke-width="16"/>';
  }

  function boxHtml(it, cls, color, inner, more) {
    return '<div class="hx-abs hx-edit-item ' + cls + '" data-hx-edit-id="' + it.id + '" style="' +
      (color ? 'color:' + color + ';' : '') + (more || '') +
      'left:' + it.x + 'px;top:' + it.y + 'px;width:' + it.w + 'px;height:' + it.h + 'px;z-index:6000;">' +
      inner + '</div>';
  }

  function shapeHtml(it) {
    if (it.kind === 'knife') return boxHtml(it, 'hx-sym is-open', '#00ff00', BLADE, '');
    if (it.kind === 'xfmr') return boxHtml(it, '', '', XFMR, '');
    if (it.kind === 'cap') return boxHtml(it, '', '', CAP, '');
    if (it.kind === 'ups') return boxHtml(it, '', '', UPS, '');
    if (it.kind === 'inv') return boxHtml(it, '', '', INV, '');
    if (it.kind === 'ssw') return boxHtml(it, '', '', SSW, '');
    if (it.kind === 'gnd') return boxHtml(it, '', '', GND, '');
    if (it.kind === 'gen') return boxHtml(it, '', '', GEN, '');
    if (it.kind === 'text') {
      return boxHtml(it, 'hx-tx', it.c || '#fff', esc(it.s || '文字'), 'font-size:' + (it.fs || 16) + 'px;');
    }
    if (it.kind === 'custom') {
      return boxHtml(it, 'hx-tx', it.c || '#fff', esc(it.s || '新组件'),
        'font-size:' + (it.fs || 13) + 'px;display:flex;align-items:center;justify-content:center;text-align:center;' +
        'border:1px solid #7fd0ff;background:rgba(8,28,52,.88);box-sizing:border-box;');
    }
    return '';
  }

  var BOX_DOTS = [
    ['nw', 'nwse', function (b) { return [b.x, b.y]; }],
    ['n', 'ns', function (b) { return [b.x + b.w / 2, b.y]; }],
    ['ne', 'nesw', function (b) { return [b.x + b.w, b.y]; }],
    ['e', 'ew', function (b) { return [b.x + b.w, b.y + b.h / 2]; }],
    ['se', 'nwse', function (b) { return [b.x + b.w, b.y + b.h]; }],
    ['s', 'ns', function (b) { return [b.x + b.w / 2, b.y + b.h]; }],
    ['sw', 'nesw', function (b) { return [b.x, b.y + b.h]; }],
    ['w', 'ew', function (b) { return [b.x, b.y + b.h / 2]; }]
  ];

  function dotHtml(extra, x, y) {
    return '<div class="hx-edit-handle" ' + extra + ' style="left:' + (x - 6) + 'px;top:' + (y - 6) + 'px"></div>';
  }

  function handleHtml(it) {
    var cur = (it.kind === 'vline') ? 'ns' : ((it.kind === 'hline' || it.kind === 'bus') ? 'ew' : '');
    var c = cur ? ' data-cursor="' + cur + '"' : '';
    var h = dotHtml('data-hx-edit-id="' + it.id + '" data-end="1"' + c, it.x1, it.y1) +
      dotHtml('data-hx-edit-id="' + it.id + '" data-end="2"' + c, it.x2, it.y2);
    if (it.kind === 'ortho' && w.HXEditCmd) {
      var m = w.HXEditCmd.orthoMid(it);
      h += dotHtml('data-hx-edit-id="' + it.id + '" data-end="m"', m.x, m.y);
    }
    return h;
  }

  function itemBoxHandles(it) {
    var b = { x: Number(it.x), y: Number(it.y), w: Number(it.w) || 0, h: Number(it.h) || 0 };
    if (!(b.w > 2 && b.h > 2)) return '';
    return BOX_DOTS.map(function (d) {
      var p = d[2](b);
      return dotHtml('data-hx-edit-id="' + it.id + '" data-end="' + d[0] + '" data-cursor="' + d[1] + '"', p[0], p[1]);
    }).join('');
  }

  function baseHandleHtml(n, mv, i) {
    if (!n) return '';
    if (n.t === 'L') {
      var g = lineGeom(n, mv);
      var axis = lineAxis(n);
      var cur = axis === 'v' ? 'ns' : (axis === 'h' ? 'ew' : '');
      var c = cur ? ' data-cursor="' + cur + '"' : '';
      var attr = 'data-hx-base="1" data-i="' + i + '"';
      return dotHtml(attr + ' data-end="1"' + c, g.x1, g.y1) + dotHtml(attr + ' data-end="2"' + c, g.x2, g.y2);
    }
    if (n.t === 'T' || n.t === 'B') {
      var b = boxGeom(n, mv);
      if (!(b.w > 2 && b.h > 2)) return '';
      return BOX_DOTS.map(function (d) {
        var p = d[2](b);
        return dotHtml('data-hx-base="1" data-i="' + i + '" data-end="' + d[0] + '" data-cursor="' + d[1] + '"', p[0], p[1]);
      }).join('');
    }
    return '';
  }

  function paintHandles(n, mv, i) {
    var layer = ensureLayer();
    if (!layer) return;
    var html = baseHandleHtml(n, mv, i);
    var old = layer.querySelector('.hx-sel-handles');
    if (!html) {
      if (old) old.remove();
      return;
    }
    if (!old) {
      old = document.createElement('div');
      old.className = 'hx-sel-handles';
      layer.appendChild(old);
    }
    old.innerHTML = html;
  }

  function preview(i, n, mv) {
    elsFor(i).forEach(function (el) { placeOne(el, n, mv); });
    paintHandles(n, mv, i);
  }

  function renderItems() {
    var layer = ensureLayer();
    if (!layer) return;
    var ov = load();
    var lines = [];
    var html = [];
    var handles = '';
    ov.items.forEach(function (it) {
      if (isWire(it.kind)) {
        lines.push(it.kind === 'ortho' ? orthoSvg(it) : lineSvg(it));
        if (selected && selected.t === 'item' && selected.id === it.id) handles = handleHtml(it);
      } else {
        html.push(shapeHtml(it));
        if (selected && selected.t === 'item' && selected.id === it.id) handles = itemBoxHandles(it);
      }
    });
    var shift = shiftEl();
    var w = shift ? shift.offsetWidth : 10;
    var h = shift ? shift.offsetHeight : 10;
    layer.innerHTML = '<svg class="hx-edit-wires" width="' + w + '" height="' + h + '">' + lines.join('') + '</svg>' + html.join('') + handles;
    if (selected && selected.t === 'base' && page) {
      var bn = page.nodes[Number(selected.i)];
      var bmv = (load().moves || {})[String(selected.i)] || {};
      paintHandles(bn, bmv, selected.i);
    }
  }

  function markSel() {
    var shift = shiftEl();
    if (!shift) return;
    Array.prototype.forEach.call(shift.querySelectorAll('.is-hx-sel'), function (el) {
      el.classList.remove('is-hx-sel');
    });
    if (!selected) { syncStyleBar(); return; }
    if (selected.t === 'base') {
      elsFor(selected.i).forEach(function (el) {
        if (el.classList.contains('hx-edit-hit')) return;
        if (kindOf(el)) el.classList.add('is-hx-sel');
      });
      var nSel = page && page.nodes[Number(selected.i)];
      var mvSel = (load().moves || {})[String(selected.i)] || {};
      paintHandles(nSel, mvSel, selected.i);
      syncStyleBar();
      return;
    }
    var oldHandles = shift.querySelector('.hx-sel-handles');
    if (oldHandles) oldHandles.remove();
    Array.prototype.forEach.call(shift.querySelectorAll('[data-hx-edit-id="' + selected.id + '"]'), function (el) {
      if (el.classList.contains('hx-edit-handle') || el.classList.contains('hx-edit-hit')) return;
      el.classList.add('is-hx-sel');
    });
    syncStyleBar();
  }

  function apply() {
    if (!page || page.kind !== 'power') return;
    var shift = shiftEl();
    if (!shift) return;
    var nodes = page.nodes || [];
    tagSymbols(shift, nodes);
    var ov = load();
    var del = {};
    ov.deleted.forEach(function (id) { del[id] = 1; });
    nodes.forEach(function (n, idx) {
      if (!editableNode(n)) return;
      var i = String(idx);
      var mv = ov.moves[i];
      var gone = !!del[i];
      elsFor(i).forEach(function (el) {
        el.style.display = gone ? 'none' : '';
        if (gone) return;
        placeOne(el, n, mv);
        if (n.t === 'T' && el.classList.contains('hx-tx')) {
          textOf(el, mv && mv.s != null ? mv.s : (n.s || ''));
        }
      });
    });
    renderItems();
    markSel();
  }

  function toolBtn(t) {
    return '<button type="button" data-hx-add="' + t[0] + '"><svg viewBox="0 0 24 36" width="16" height="24" aria-hidden="true">' +
      t[2] + '</svg><span>' + t[1] + '</span></button>';
  }

  function palHtml() {
    var groups = GROUPS.map(function (g) {
      return '<details class="hx-edit-grp" open><summary>' + g.name + '</summary><div class="hx-edit-grp-body">' +
        g.tools.map(toolBtn).join('') + '</div></details>';
    }).join('');
    return '<aside class="hx-edit-pal" hidden>' +
      '<div class="hx-edit-pal-h"><span>组件</span><span class="hx-edit-pal-acts">' +
      '<button type="button" data-hx-pal="min" title="最小化" aria-label="最小化组件栏">−</button>' +
      '<button type="button" data-hx-pal="hide" title="隐藏" aria-label="隐藏组件栏">×</button>' +
      '</span></div>' + groups +
      '<form class="hx-edit-cmd"><input name="cmd" placeholder="指令：加3把竖刀" autocomplete="off"/><button type="submit">执行</button>' +
      '<p class="hx-edit-cmd-msg"></p></form>' +
      '<p class="hx-edit-tip">分组标题可折叠。点一下或拖进图里。选中文字后右下角改字号；再点一次改里面的字，Enter 记下，Shift+Enter 换行。跟闸的字不单独改色，字、框、出线仍是同一个颜色，合红分绿。线和框选中后拖黄点改长短、大小。线头靠近另一根线的端点时会吸过去，横平竖直的仍保持平直并对齐接上。点选被盖住的线：按住 Alt 再点，或在同一位置再点一次切换。指令例：加UPS、加接地、加母线、加自定义 接地。Delete 删除。改完要点保存，才写进这间一次图。</p></aside>' +
      '<button type="button" class="hx-edit-pal-dock" data-hx-pal="show" hidden>组件</button>';
  }

  function styleHtml() {
    return '<div class="hx-edit-style" hidden>' +
      '<label>字号 <input type="number" min="8" max="96" step="1" data-hx-fs /></label>' +
      '<span data-hx-color-free>' +
      '<label>颜色 <input type="color" data-hx-color /></label>' +
      '<button type="button" data-hx-swatch="#ffffff" title="白" style="background:#ffffff"></button>' +
      '<button type="button" data-hx-swatch="#ff0000" title="红" style="background:#ff0000"></button>' +
      '<button type="button" data-hx-swatch="#00ff00" title="绿" style="background:#00ff00"></button>' +
      '<button type="button" data-hx-swatch="#ffd24d" title="黄" style="background:#ffd24d"></button>' +
      '</span>' +
      '<span data-hx-knife-note hidden>颜色跟闸，合红分绿</span>' +
      '</div>';
  }

  function hexColor(s, fallback) {
    s = String(s || '').trim();
    if (/^#[0-9a-fA-F]{6}$/.test(s)) return s.toLowerCase();
    if (/^#[0-9a-fA-F]{3}$/.test(s)) return ('#' + s[1] + s[1] + s[2] + s[2] + s[3] + s[3]).toLowerCase();
    var m = s.match(/rgba?\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)/i);
    if (!m) return fallback || '#ffffff';
    function h(n) {
      var v = Number(n).toString(16);
      return v.length < 2 ? '0' + v : v;
    }
    return '#' + h(m[1]) + h(m[2]) + h(m[3]);
  }

  function selectedTextInfo() {
    if (!selected || !page) return null;
    if (selected.t === 'base') {
      var n = page.nodes[Number(selected.i)];
      if (!n || n.t !== 'T') return null;
      var mv = (load().moves || {})[String(selected.i)] || {};
      var el = null;
      elsFor(selected.i).forEach(function (e) {
        if (e.classList && e.classList.contains('hx-tx')) el = e;
      });
      var fs = mv.fs != null ? Number(mv.fs) : ((el && parseFloat(el.style.fontSize)) || Number(n.fs) || 14);
      var c = (el && el.style.color) || n.c || '#ffffff';
      var knife = !!(n.followKnife || n.upsPlate || (el && el.classList.contains('hx-plate-live')));
      return { fs: fs, c: c, knife: knife };
    }
    if (selected.t === 'item') {
      var it = null;
      (load().items || []).forEach(function (x) { if (x.id === selected.id) it = x; });
      if (!it || (it.kind !== 'text' && it.kind !== 'custom')) return null;
      return { fs: it.fs || (it.kind === 'text' ? 16 : 13), c: it.c || '#ffffff', knife: false };
    }
    return null;
  }

  function syncStyleBar() {
    var bar = host && host.querySelector('.hx-edit-style');
    if (!bar) return;
    var on = !!(editing && page && page.kind === 'power');
    var info = on ? selectedTextInfo() : null;
    if (!info) { bar.hidden = true; return; }
    bar.hidden = false;
    var fs = bar.querySelector('[data-hx-fs]');
    var free = bar.querySelector('[data-hx-color-free]');
    var note = bar.querySelector('[data-hx-knife-note]');
    var color = bar.querySelector('[data-hx-color]');
    if (fs && document.activeElement !== fs) fs.value = String(Math.round(info.fs));
    if (free) free.hidden = !!info.knife;
    if (note) note.hidden = !info.knife;
    if (!info.knife && color && document.activeElement !== color) color.value = hexColor(info.c, '#ffffff');
  }

  function setTextStyle(part) {
    if (!selected || !page) return;
    if (selected.t === 'base') {
      var n = page.nodes[Number(selected.i)];
      if (!n || n.t !== 'T') return;
      if (part.c && (n.followKnife || n.upsPlate)) return;
      commit(function (ov) {
        var rec = ov.moves[String(selected.i)] || {};
        if (part.fs != null) rec.fs = part.fs;
        if (part.c && !n.followKnife && !n.upsPlate) rec.c = part.c;
        ov.moves[String(selected.i)] = rec;
      });
      var saved = (load().moves || {})[String(selected.i)] || {};
      elsFor(selected.i).forEach(function (el) {
        if (el.classList && el.classList.contains('hx-tx')) paintTextLook(el, saved, n);
      });
      return;
    }
    if (selected.t !== 'item') return;
    commit(function (ov) {
      var i;
      for (i = 0; i < ov.items.length; i++) {
        if (ov.items[i].id !== selected.id) continue;
        if (ov.items[i].kind !== 'text' && ov.items[i].kind !== 'custom') return;
        if (part.fs != null) ov.items[i].fs = part.fs;
        if (part.c) ov.items[i].c = part.c;
      }
    });
    var layer = ensureLayer();
    var el = layer && layer.querySelector('.hx-tx[data-hx-edit-id="' + selected.id + '"]');
    if (!el) return;
    if (part.fs != null) el.style.fontSize = part.fs + 'px';
    if (part.c) {
      el.style.color = part.c;
      el.style.borderColor = part.c;
    }
  }

  function onStyleInput(ev) {
    var t = ev.target;
    if (!t || !t.matches) return;
    if (t.matches('[data-hx-fs]')) {
      var fs = Number(t.value);
      if (!(fs >= 8 && fs <= 96)) return;
      setTextStyle({ fs: fs });
      return;
    }
    if (t.matches('[data-hx-color]')) setTextStyle({ c: t.value });
  }

  function onStyleClick(ev) {
    var b = ev.target.closest ? ev.target.closest('[data-hx-swatch]') : null;
    if (!b) return;
    ev.preventDefault();
    ev.stopPropagation();
    var c = b.getAttribute('data-hx-swatch');
    var bar = host && host.querySelector('.hx-edit-style');
    var inp = bar && bar.querySelector('[data-hx-color]');
    if (inp) inp.value = c;
    setTextStyle({ c: c });
  }

  function applyPalMode() {
    var on = !!(editing && page && page.kind === 'power');
    var pal = host && host.querySelector('.hx-edit-pal');
    var dock = host && host.querySelector('.hx-edit-pal-dock');
    if (pal) {
      pal.hidden = !on || palMode === 'hide';
      pal.classList.toggle('is-min', palMode === 'min');
      var minBtn = pal.querySelector('[data-hx-pal="min"]');
      if (minBtn) {
        var folded = palMode === 'min';
        minBtn.textContent = folded ? '▢' : '−';
        minBtn.title = folded ? '展开' : '最小化';
        minBtn.setAttribute('aria-label', folded ? '展开组件栏' : '最小化组件栏');
      }
    }
    if (dock) dock.hidden = !on || palMode !== 'hide';
  }

  function syncChrome() {
    var on = !!(editing && page && page.kind === 'power');
    var wrap = host && host.querySelector('.hx-wrap');
    if (wrap) wrap.classList.toggle('is-hx-editing', on);
    applyPalMode();
    if (!host) return;
    Array.prototype.forEach.call(host.querySelectorAll('[data-hx-edit]'), function (b) {
      var act = b.getAttribute('data-hx-edit');
      if (act === 'toggle') {
        b.classList.toggle('is-on', on);
        b.textContent = '编辑';
      } else b.hidden = !on;
    });
    syncStyleBar();
  }

  function ensureUi() {
    var body = host && host.querySelector('.hx-body');
    if (body && !body.querySelector('.hx-edit-pal')) {
      body.insertAdjacentHTML('beforeend', palHtml());
      var dock = body.querySelector('.hx-edit-pal-dock');
      var drillBar = body.querySelector('.hx-drill-bar');
      if (dock && drillBar) drillBar.insertBefore(dock, drillBar.firstChild);
      var pal = body.querySelector('.hx-edit-pal');
      pal.addEventListener('pointerdown', gestures.onPalDown);
      var form = pal.querySelector('.hx-edit-cmd');
      if (form) form.addEventListener('submit', onCmdSubmit);
      body.addEventListener('click', function (ev) {
        var b = ev.target.closest ? ev.target.closest('[data-hx-pal]') : null;
        if (!b || !body.contains(b)) return;
        ev.preventDefault();
        ev.stopPropagation();
        var act = b.getAttribute('data-hx-pal');
        if (act === 'min') palMode = palMode === 'min' ? 'full' : 'min';
        else if (act === 'hide') palMode = 'hide';
        else if (act === 'show') palMode = 'full';
        applyPalMode();
      });
    }
    var bar = host && host.querySelector('.hx-zoom-bar');
    if (bar && !bar.querySelector('[data-hx-edit]')) {
      bar.insertAdjacentHTML('afterbegin',
        '<button type="button" data-hx-edit="toggle">编辑</button>' +
        '<button type="button" data-hx-edit="save" hidden>保存</button>' +
        '<button type="button" data-hx-edit="del" hidden>删除</button>' +
        '<button type="button" data-hx-edit="export" hidden>导出</button>' +
        '<button type="button" data-hx-edit="restore" hidden>还原</button>');
    }
    var main = host && host.querySelector('.hx-main');
    if (main && !main.querySelector('.hx-edit-style')) {
      main.insertAdjacentHTML('beforeend', styleHtml());
      var styleBar = main.querySelector('.hx-edit-style');
      styleBar.addEventListener('input', onStyleInput);
      styleBar.addEventListener('change', onStyleInput);
      styleBar.addEventListener('click', onStyleClick);
    }
    if (main && !main._hxEdit) {
      main._hxEdit = 1;
      main.addEventListener('click', onBarClick);
      main.addEventListener('pointerdown', gestures.onDown);
      main.addEventListener('dblclick', gestures.onDbl, true);
    }
  }

  function makeItem(kind, x, y) {
    var it = { id: nid(), kind: kind };
    if (kind === 'knife') { it.w = 60; it.h = 98; it.x = r1(x - 30); it.y = r1(y - 49); }
    else if (kind === 'xfmr') { it.w = 64; it.h = 88; it.x = r1(x - 32); it.y = r1(y - 44); }
    else if (kind === 'cap') { it.w = 100; it.h = 230; it.x = r1(x - 50); it.y = r1(y - 20); }
    else if (kind === 'ups') { it.w = 72; it.h = 48; it.x = r1(x - 36); it.y = r1(y - 24); }
    else if (kind === 'inv') { it.w = 72; it.h = 48; it.x = r1(x - 36); it.y = r1(y - 24); }
    else if (kind === 'ssw') { it.w = 64; it.h = 52; it.x = r1(x - 32); it.y = r1(y - 26); }
    else if (kind === 'text') { it.w = 96; it.h = 28; it.x = r1(x - 48); it.y = r1(y - 14); it.s = '文字'; }
    else if (kind === 'custom') { it.w = 96; it.h = 40; it.x = r1(x - 48); it.y = r1(y - 20); it.s = '新组件'; }
    else if (kind === 'gnd') { it.w = 48; it.h = 36; it.x = r1(x - 24); it.y = r1(y - 8); }
    else if (kind === 'gen') { it.w = 64; it.h = 64; it.x = r1(x - 32); it.y = r1(y - 32); }
    else if (kind === 'bus') { it.x1 = r1(x - 140); it.y1 = r1(y); it.x2 = r1(x + 140); it.y2 = r1(y); }
    else if (kind === 'hline') { it.x1 = r1(x - 80); it.y1 = r1(y); it.x2 = r1(x + 80); it.y2 = r1(y); }
    else if (kind === 'vline') { it.x1 = r1(x); it.y1 = r1(y - 80); it.x2 = r1(x); it.y2 = r1(y + 80); }
    else if (kind === 'ortho') {
      it.x1 = r1(x - 70); it.y1 = r1(y - 36);
      it.x2 = r1(x + 70); it.y2 = r1(y + 36);
      it.bend = 'auto';
    }
    else return null;
    return it;
  }

  function addAt(kind, clientX, clientY) {
    var p = pageXY(clientX, clientY);
    if (!p) return null;
    var it = makeItem(kind, p.x, p.y);
    if (!it) return null;
    selected = { t: 'item', id: it.id };
    commit(function (ov) { ov.items.push(it); });
    renderItems();
    markSel();
    return it;
  }

  function addAtCenter(kind) {
    var fit = host && host.querySelector('.hx-fit');
    if (!fit) return null;
    var r = fit.getBoundingClientRect();
    return addAt(kind, r.left + r.width / 2, r.top + r.height / 2);
  }

  function showCmd(msg, bad) {
    var el = host && host.querySelector('.hx-edit-cmd-msg');
    if (!el) return;
    el.textContent = msg || '';
    el.style.color = bad ? '#ffb4b4' : '#9eebff';
  }

  var KIND_NAME = {
    knife: '竖刀', xfmr: '变压器', gen: '柴油发电机', ups: 'UPS', inv: '逆变器',
    cap: '电容补偿柜', gnd: '接地', ssw: '静态开关',
    text: '文字', custom: '自定义',
    bus: '母线', hline: '横线', vline: '竖线', ortho: '智能折线'
  };

  function execText(raw) {
    if (!w.HXEditCmd) return '指令模块没加载';
    if (!editing) toggle();
    var cmd = w.HXEditCmd.parse(raw);
    if (cmd.err) { showCmd(cmd.err, true); return cmd.err; }
    if (cmd.act === 'del') {
      if (!selected) { showCmd('先选中要删的', true); return '先选中要删的'; }
      removeSelected();
      showCmd('已删除');
      return '已删除';
    }
    var fit = host && host.querySelector('.hx-fit');
    if (!fit) { showCmd('画面还没出来', true); return '画面还没出来'; }
    var r = fit.getBoundingClientRect();
    var i, lastId = '';
    for (i = 0; i < cmd.count; i++) {
      var it = addAt(cmd.kind, r.left + r.width / 2 + i * 36, r.top + r.height / 2 + i * 28);
      if (it) lastId = it.id;
    }
    if ((cmd.kind === 'text' || cmd.kind === 'custom') && cmd.text && lastId) {
      commit(function (ov) {
        var k;
        for (k = 0; k < ov.items.length; k++) if (ov.items[k].id === lastId) ov.items[k].s = cmd.text;
      });
      renderItems();
      markSel();
    }
    var msg = '已加 ' + cmd.count + ' 个' + (KIND_NAME[cmd.kind] || cmd.kind);
    showCmd(msg);
    return msg;
  }

  function onCmdSubmit(ev) {
    ev.preventDefault();
    var input = ev.target.querySelector('input[name="cmd"]');
    if (!input) return;
    execText(input.value);
    input.select();
  }

  function pollInbox() {
    if (!editing || !page) return;
    fetch('/hx-dcim/data/hx-edit-inbox.json?t=' + Date.now(), { cache: 'no-store' }).then(function (res) {
      return res.ok ? res.json() : null;
    }).then(function (j) {
      if (!j || !j.seq || j.seq === inboxSeen || !j.text) return;
      if (j.room && String(j.room) !== roomOf(page)) return;
      inboxSeen = j.seq;
      execText(j.text);
    }).catch(function () {});
  }

  function toggle() {
    if (!page || page.kind !== 'power') return;
    editing = !editing;
    if (editing && dismissHover) dismissHover();
    if (!editing) selected = null;
    syncChrome();
    renderItems();
    markSel();
  }

  function removeSelected() {
    if (!selected) return;
    var sel = selected;
    if (sel.t === 'base') {
      var n = page && page.nodes[Number(sel.i)];
      if (isProtectedNode(n)) {
        w.alert('这是 411/421 竖母线，不能删。删了铭牌横抽头会悬空。');
        return;
      }
    }
    selected = null;
    if (sel.t === 'base') {
      commit(function (ov) {
        if (ov.deleted.indexOf(sel.i) < 0) ov.deleted.push(sel.i);
        delete ov.moves[sel.i];
      });
    } else {
      commit(function (ov) {
        ov.items = ov.items.filter(function (it) { return it.id !== sel.id; });
      });
    }
    apply();
  }

  function exportJson() {
    if (!page) return;
    var ov = load();
    var payload = {
      version: VER,
      roomId: roomOf(page),
      roomName: page.name || '',
      pageId: page.id || '',
      moves: ov.moves,
      deleted: ov.deleted,
      items: ov.items
    };
    var blob = new Blob([JSON.stringify(payload, null, 2)], { type: 'application/json' });
    var a = document.createElement('a');
    a.href = URL.createObjectURL(blob);
    a.download = 'hx-edit-' + payload.roomId + '.json';
    document.body.appendChild(a);
    a.click();
    a.remove();
    w.setTimeout(function () { URL.revokeObjectURL(a.href); }, 1200);
  }

  function restore() {
    if (!page || page.kind !== 'power') return;
    var name = page.name || roomOf(page);
    if (!w.confirm('清掉「' + name + '」在这台浏览器里的编辑，回到原来的图？')) return;
    try { w.localStorage.removeItem(key()); } catch (e) {}
    selected = null;
    apply();
  }

  function r2(n) { return Math.round(Number(n) * 100) / 100; }

  function segKept(a, b) {
    var horiz = Math.abs(a[1] - b[1]) < 0.8;
    var len = horiz ? Math.abs(b[0] - a[0]) : Math.abs(b[1] - a[1]);
    return len >= 0.5;
  }

  function buildPatch(lay, ov) {
    var baked = JSON.parse(JSON.stringify(lay));
    var patch = {
      room: roomOf(page),
      texts: [], knives: [], boxes: [], wires: [],
      remove: { texts: [], knives: [], boxes: [], wires: [] },
      add: { texts: [], knives: [], boxes: [], wires: [] }
    };
    var del = {};
    (ov.deleted || []).forEach(function (id) { del[String(id)] = 1; });
    var moves = ov.moves || {};
    var i = 0;
    var wireNode = [];
    (baked.wires || []).forEach(function (w, wi) {
      var s, n = 0;
      for (s = 1; s < (w.pts || []).length; s++) {
        if (!segKept(w.pts[s - 1], w.pts[s])) continue;
        wireNode.push({ wi: wi, seg: s, i: i });
        n++;
        i++;
      }
      w._n = n;
    });
    var boxAt = (baked.boxes || []).map(function () { return i++; });
    var knifeAt = (baked.knives || []).map(function () { return i++; });
    var textAt = (baked.texts || []).map(function () { return i++; });
    var byWire = {};
    wireNode.forEach(function (seg) {
      var id = String(seg.i);
      if (!byWire[seg.wi]) byWire[seg.wi] = { del: 0, segs: [] };
      var info = byWire[seg.wi];
      if (del[id]) info.del++;
      var mv = moves[id];
      if (!mv) return;
      if (mv.x1 != null) info.segs.push({ seg: seg.seg, i: seg.i, abs: 1, x1: mv.x1, y1: mv.y1, x2: mv.x2, y2: mv.y2 });
      else if ((mv.dx || 0) || (mv.dy || 0)) info.segs.push({ seg: seg.seg, dx: mv.dx || 0, dy: mv.dy || 0 });
    });
    Object.keys(byWire).forEach(function (wis) {
      var wi = Number(wis);
      var info = byWire[wi];
      var ww = baked.wires[wi];
      if (info.del && info.del === ww._n) {
        patch.remove.wires.push(wi);
        return;
      }
      if (!info.segs.length) return;
      var pts = ww.pts.map(function (p) { return [p[0], p[1]]; });
      var target = {};
      var conflict = false;
      function put(pi, x, y) {
        var prev = target[pi];
        if (!prev) target[pi] = [x, y];
        else if (Math.abs(prev[0] - x) > 0.05 || Math.abs(prev[1] - y) > 0.05) conflict = true;
      }
      info.segs.forEach(function (s) {
        if (s.abs) {
          var node = page.nodes[s.i];
          var e = endsOf(node);
          var a = ww.pts[s.seg - 1];
          var d1 = Math.hypot(a[0] - e[0], a[1] - e[1]);
          var d2 = Math.hypot(a[0] - e[2], a[1] - e[3]);
          if (d1 <= d2) {
            put(s.seg - 1, s.x1, s.y1);
            put(s.seg, s.x2, s.y2);
          } else {
            put(s.seg - 1, s.x2, s.y2);
            put(s.seg, s.x1, s.y1);
          }
          return;
        }
        put(s.seg - 1, ww.pts[s.seg - 1][0] + s.dx, ww.pts[s.seg - 1][1] + s.dy);
        put(s.seg, ww.pts[s.seg][0] + s.dx, ww.pts[s.seg][1] + s.dy);
      });
      if (conflict) return;
      Object.keys(target).forEach(function (pi) {
        pts[pi][0] = r2(target[pi][0]);
        pts[pi][1] = r2(target[pi][1]);
      });
      ww.pts = pts;
      patch.wires.push({ i: wi, pts: pts });
    });
    function takeBox(arr, at, bucket, removed, withText) {
      arr.forEach(function (n, bi) {
        var idx = at[bi];
        if (del[String(idx)]) { removed.push(bi); return; }
        var mv = moves[String(idx)];
        if (!mv) return;
        var dx = mv.dx || 0;
        var dy = mv.dy || 0;
        var rec = null;
        if (mv.w != null && mv.x != null) {
          n.x = r2(mv.x);
          n.y = r2(mv.y);
          n.w = r2(mv.w);
          n.h = r2(mv.h);
          rec = { i: bi, x: n.x, y: n.y, w: n.w, h: n.h };
        } else if (dx || dy) {
          n.x = r2(Number(n.x) + dx);
          n.y = r2(Number(n.y) + dy);
          rec = { i: bi, x: n.x, y: n.y };
        }
        if (withText && mv.s != null) {
          n.s = mv.s;
          rec = rec || { i: bi };
          rec.s = mv.s;
        }
        if (rec) bucket.push(rec);
      });
    }
    baked.boxes = baked.boxes || [];
    baked.knives = baked.knives || [];
    baked.texts = baked.texts || [];
    baked.wires = baked.wires || [];
    takeBox(baked.boxes, boxAt, patch.boxes, patch.remove.boxes, false);
    takeBox(baked.knives, knifeAt, patch.knives, patch.remove.knives, false);
    /* 铭牌在画面节点里的顺序就是 JSON texts 的顺序。5A1 的线段过滤和这里不一致，不能拿节点总序号去对文字下标。 */
    var textOrd = 0;
    (page.nodes || []).forEach(function (n, idx) {
      if (!n || n.t !== 'T') return;
      var bi = textOrd++;
      if (bi >= baked.texts.length) return;
      if (del[String(idx)]) { patch.remove.texts.push(bi); return; }
      var mv = moves[String(idx)];
      if (!mv) return;
      var dx = mv.dx || 0;
      var dy = mv.dy || 0;
      var rec = null;
      var tx = baked.texts[bi];
      if (mv.w != null && mv.x != null) {
        tx.x = r2(mv.x);
        tx.y = r2(mv.y);
        tx.w = r2(mv.w);
        tx.h = r2(mv.h);
        rec = { i: bi, x: tx.x, y: tx.y, w: tx.w, h: tx.h };
      } else if (dx || dy) {
        tx.x = r2(Number(tx.x) + dx);
        tx.y = r2(Number(tx.y) + dy);
        rec = { i: bi, x: tx.x, y: tx.y };
      }
      if (mv.s != null) {
        tx.s = mv.s;
        rec = rec || { i: bi };
        rec.s = mv.s;
      }
      if (mv.fs != null) {
        tx.fs = r2(mv.fs);
        tx.fsLock = 1;
        rec = rec || { i: bi };
        rec.fs = tx.fs;
        rec.fsLock = true;
      }
      if (mv.c && !n.followKnife && !n.upsPlate) {
        tx.c = mv.c;
        rec = rec || { i: bi };
        rec.c = tx.c;
      }
      if (rec) patch.texts.push(rec);
    });
    (baked.wires || []).forEach(function (w) { delete w._n; });
    function drop(arr, idx) {
      idx.slice().sort(function (a, b) { return b - a; }).forEach(function (k) { arr.splice(k, 1); });
    }
    drop(baked.wires, patch.remove.wires);
    drop(baked.boxes, patch.remove.boxes);
    drop(baked.knives, patch.remove.knives);
    drop(baked.texts, patch.remove.texts);
    (ov.items || []).forEach(function (it) {
      if (it.kind === 'text' || it.kind === 'custom') {
        var tx = { s: it.s || '', x: r2(it.x), y: r2(it.y), w: it.w || 96, h: it.h || 28, fs: it.fs || 14, c: it.c || '#ffffff', al: 'c' };
        if (it.fs != null) tx.fsLock = 1;
        patch.add.texts.push(tx);
        baked.texts.push(tx);
      } else if (it.kind === 'knife') {
        var kn = { cad: it.id || 'new', x: r2(it.x), y: r2(it.y), out: false, role: '', amp: '', use: '' };
        patch.add.knives.push(kn);
        baked.knives.push(kn);
      } else if (it.kind === 'xfmr' || it.kind === 'cap') {
        var bx = { x: r2(it.x), y: r2(it.y), w: it.w, h: it.h, sym: it.kind, id: it.id || it.kind };
        patch.add.boxes.push(bx);
        baked.boxes.push(bx);
      } else if (it.kind === 'hline' || it.kind === 'vline' || it.kind === 'bus' || it.kind === 'ortho') {
        var wr = {
          net: 'w',
          pts: [[r2(it.x1), r2(it.y1)], [r2(it.x2), r2(it.y2)]],
          c: '#ce0c18',
          sw: it.kind === 'bus' ? 3 : 2
        };
        patch.add.wires.push(wr);
        baked.wires.push(wr);
      }
    });
    return { patch: patch, baked: baked };
  }

  function flashSave(text) {
    saving = false;
    var btn = host && host.querySelector('[data-hx-edit="save"]');
    if (!btn) return;
    btn.disabled = false;
    btn.textContent = text;
    if (text !== '保存') {
      w.setTimeout(function () {
        if (btn.textContent === text) btn.textContent = '保存';
      }, 1600);
    }
  }

  function saveRoom() {
    if (!page || page.kind !== 'power') return;
    if (saving) return;
    var code = roomOf(page);
    var lay = w.HXRooms && w.HXRooms.getLayout && w.HXRooms.getLayout(code);
    if (!lay && code === '5A1' && w.HX5A1 && w.HX5A1.getLayout) lay = w.HX5A1.getLayout();
    var btn = host && host.querySelector('[data-hx-edit="save"]');
    if (!lay) {
      if (btn) btn.textContent = '保存不了';
      flashSave('保存不了');
      return;
    }
    saving = true;
    if (btn) { btn.disabled = true; btn.textContent = '保存中'; }
    var built = buildPatch(lay, load());
    var headers = (w.HXAuth && w.HXAuth.headers) ? w.HXAuth.headers() : { 'Content-Type': 'application/json' };
    if (!headers['Content-Type']) headers['Content-Type'] = 'application/json';
    fetch('/api/hxPowerRoomSave', {
      method: 'POST',
      credentials: 'same-origin',
      headers: headers,
      body: JSON.stringify(built.patch)
    }).then(function (res) {
      return res.json().then(function (j) { return j; }, function () { throw new Error('保存失败'); });
    }).then(function (j) {
      if (!j || j.code !== 200) throw new Error('保存失败');
      try {
        if (code === '5A1' && w.HX5A1 && w.HX5A1.setLayout) w.HX5A1.setLayout(built.baked);
        else if (w.HXRooms && w.HXRooms.setLayout) w.HXRooms.setLayout(code, built.baked);
        try { w.localStorage.removeItem(key()); } catch (e) {}
        if (w.HXScreen && w.HXScreen.reloadCurrent) w.HXScreen.reloadCurrent();
      } catch (e) {}
      flashSave('保存成功');
    }).catch(function () {
      flashSave('保存失败');
    });
  }

  function onBarClick(ev) {
    var b = ev.target.closest ? ev.target.closest('[data-hx-edit]') : null;
    if (!b) return;
    ev.preventDefault();
    ev.stopPropagation();
    var act = b.getAttribute('data-hx-edit');
    if (act === 'toggle') toggle();
    else if (act === 'save') saveRoom();
    else if (act === 'del') removeSelected();
    else if (act === 'export') exportJson();
    else if (act === 'restore') restore();
  }

  gestures = w.HXEditGestures.install({
    isOn: function () { return !!(editing && page && page.kind === 'power'); },
    host: function () { return host; },
    page: function () { return page; },
    sel: function () { return selected; },
    setSel: function (s) { selected = s; },
    load: load,
    commit: commit,
    renderItems: renderItems,
    markSel: markSel,
    scale: scale,
    pageXY: pageXY,
    textOf: textOf,
    elsFor: elsFor,
    placeOne: placeOne,
    preview: preview,
    lineGeom: lineGeom,
    boxGeom: boxGeom,
    lineAxis: lineAxis,
    resizeLine: resizeLine,
    resizeBox: resizeBox,
    kindOf: kindOf,
    ensureLayer: ensureLayer,
    addAt: addAt,
    addAtCenter: addAtCenter,
    removeSelected: removeSelected,
    r1: r1
  });

  function afterRender(h, pg) {
    host = h;
    if (gestures) gestures.cancel();
    var prev = room;
    page = pg;
    room = roomOf(pg);
    if (room !== prev) selected = null;
    if (!pg || pg.kind !== 'power') {
      editing = false;
      syncChrome();
      return;
    }
    ensureUi();
    apply();
    syncChrome();
  }

  w.setInterval(pollInbox, 2000);

  w.HXEdit = {
    afterRender: afterRender,
    blocksStage: function () { return !!(editing && page && page.kind === 'power'); },
    bindHost: function (api) { dismissHover = api && api.dismissHover; },
    run: execText
  };
})(window);
