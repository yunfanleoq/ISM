/* 一次图编辑叠加层。不改各室电气 json。按房间写入 localStorage，键 hx-dcim-edit:<室号>。 */
(function (w) {
  'use strict';

  var VER = '20261028cad76';
  var inboxSeen = '';
  var palMode = 'full';
  var editing = false;
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

  function placeOne(el, n, mv) {
    var dx = (mv && mv.dx) || 0;
    var dy = (mv && mv.dy) || 0;
    var tag = (el.tagName || '').toLowerCase();
    if (tag === 'line' || tag === 'path') {
      if (dx || dy) el.setAttribute('transform', 'translate(' + dx + ' ' + dy + ')');
      else el.removeAttribute('transform');
      return;
    }
    el.style.left = (Number(n.x) + dx) + 'px';
    el.style.top = (Number(n.y) + dy) + 'px';
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
    if (it.kind === 'text') return boxHtml(it, 'hx-tx', '#fff', esc(it.s || '文字'), 'font-size:16px;');
    if (it.kind === 'custom') {
      return boxHtml(it, 'hx-tx', '#fff', esc(it.s || '新组件'),
        'font-size:13px;display:flex;align-items:center;justify-content:center;text-align:center;' +
        'border:1px solid #7fd0ff;background:rgba(8,28,52,.88);box-sizing:border-box;');
    }
    return '';
  }

  function handleHtml(it) {
    var h = '<div class="hx-edit-handle" data-hx-edit-id="' + it.id + '" data-end="1" style="left:' + (it.x1 - 6) + 'px;top:' + (it.y1 - 6) + 'px"></div>' +
      '<div class="hx-edit-handle" data-hx-edit-id="' + it.id + '" data-end="2" style="left:' + (it.x2 - 6) + 'px;top:' + (it.y2 - 6) + 'px"></div>';
    if (it.kind === 'ortho' && w.HXEditCmd) {
      var m = w.HXEditCmd.orthoMid(it);
      h += '<div class="hx-edit-handle is-mid" data-hx-edit-id="' + it.id + '" data-end="m" style="left:' + (m.x - 6) + 'px;top:' + (m.y - 6) + 'px"></div>';
    }
    return h;
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
      } else html.push(shapeHtml(it));
    });
    var shift = shiftEl();
    var w = shift ? shift.offsetWidth : 10;
    var h = shift ? shift.offsetHeight : 10;
    layer.innerHTML = '<svg class="hx-edit-wires" width="' + w + '" height="' + h + '">' + lines.join('') + '</svg>' + html.join('') + handles;
  }

  function markSel() {
    var shift = shiftEl();
    if (!shift) return;
    Array.prototype.forEach.call(shift.querySelectorAll('.is-hx-sel'), function (el) {
      el.classList.remove('is-hx-sel');
    });
    if (!selected) return;
    if (selected.t === 'base') {
      elsFor(selected.i).forEach(function (el) {
        if (el.classList.contains('hx-edit-hit')) return;
        if (kindOf(el)) el.classList.add('is-hx-sel');
      });
      return;
    }
    Array.prototype.forEach.call(shift.querySelectorAll('[data-hx-edit-id="' + selected.id + '"]'), function (el) {
      if (el.classList.contains('hx-edit-handle') || el.classList.contains('hx-edit-hit')) return;
      el.classList.add('is-hx-sel');
    });
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
      '<p class="hx-edit-tip">分组标题可折叠。点一下或拖进图里。双击文字、自定义改字。指令例：加UPS、加接地、加母线、加自定义 接地。Delete 删除。改动记在这台浏览器里。</p></aside>' +
      '<button type="button" class="hx-edit-pal-dock" data-hx-pal="show" hidden>组件</button>';
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
  }

  function ensureUi() {
    var body = host && host.querySelector('.hx-body');
    if (body && !body.querySelector('.hx-edit-pal')) {
      body.insertAdjacentHTML('beforeend', palHtml());
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
        '<button type="button" data-hx-edit="del" hidden>删除</button>' +
        '<button type="button" data-hx-edit="export" hidden>导出</button>' +
        '<button type="button" data-hx-edit="restore" hidden>还原</button>');
    }
    var main = host && host.querySelector('.hx-main');
    if (main && !main._hxEdit) {
      main._hxEdit = 1;
      main.addEventListener('click', onBarClick);
      main.addEventListener('pointerdown', gestures.onDown);
      main.addEventListener('dblclick', gestures.onDbl);
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

  function onBarClick(ev) {
    var b = ev.target.closest ? ev.target.closest('[data-hx-edit]') : null;
    if (!b) return;
    ev.preventDefault();
    ev.stopPropagation();
    var act = b.getAttribute('data-hx-edit');
    if (act === 'toggle') toggle();
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
    elsFor: elsFor,
    placeOne: placeOne,
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
