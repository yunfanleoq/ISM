/* 按航信「一次图界面」元件位置重画。不走旧组态运行时。 */
(function (w) {
  'use strict';

  function esc(s) {
    return String(s == null ? '' : s)
      .replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
  }

  function ends(n) {
    var rad = (Number(n.a) || 0) * Math.PI / 180;
    var cos = Math.cos(rad);
    var sin = Math.sin(rad);
    var cx = n.x + n.w / 2;
    var cy = n.y + n.h / 2;
    var hw = n.w / 2;
    return [cx - hw * cos, cy - hw * sin, cx + hw * cos, cy + hw * sin];
  }

  function knifeHorizontal(n) {
    return n.orient === 'h' || Number(n.rot) === 90;
  }

  var KNIFE_W = 60;
  var KNIFE_H = 98;

  /** 横置刀闸：固定 60×98 数据框，交换 w/h 作点击区；避免合闸帧 88×100 放大。 */
  function symBox(n) {
    if (knifeHorizontal(n)) return { w: KNIFE_H, h: KNIFE_W };
    return { w: KNIFE_W, h: KNIFE_H };
  }

  var TEXT_ALIGN = {
    l: 'justify-content:flex-start;text-align:left;',
    r: 'justify-content:flex-end;text-align:right;'
  };

  function pos(n) {
    /* 工程里的 z 有到 -1000 的框。加 120 仍是负数，会画到背景后面，红框就看不见。 */
    var b = n.lay && n.t !== 'I' && n.t !== 'S' ? { w: n.w, h: n.h } : symBox(n);
    return 'left:' + n.x + 'px;top:' + n.y + 'px;width:' + b.w + 'px;height:' + b.h + 'px;z-index:' + ((Number(n.z) || 0) + 2000);
  }

  /* 刀下柜号、电流、负荷名。白三角在刀框底边外约 8px，文字再下移 18px 离开三角。不改刀、不改三角。 */
  var KNIFE_TEXT_DROP = 18;

  function underKnifeTextDrop(n, knives) {
    var tw = Number(n.w) || 0;
    var cx = Number(n.x) + tw / 2;
    var ty = Number(n.y);
    var i;
    for (i = 0; i < knives.length; i++) {
      var k = knives[i];
      if (knifeHorizontal(k)) continue;
      var kcx = Number(k.x) + (Number(k.w) || KNIFE_W) / 2;
      if (Math.abs(cx - kcx) > 40) continue;
      var rel = ty - (Number(k.y) + (Number(k.h) || KNIFE_H));
      if (rel >= -12 && rel <= 60) return KNIFE_TEXT_DROP;
    }
    return 0;
  }

  /* 母线、UPS 出线用原图固定色。合闸红、分闸绿只留给开关本身。 */
  function lineColor(page, n) {
    if (!page || page.kind === 'column') return n.c || '#ffffff';
    /* 5A1 显式布局线（fixc）保留蓝/紫/绿/红，不改其它配电室。 */
    if (n.fixc && n.c) return n.c;
    var k = String(n.k || '');
    /* UPS 进出线、旁路、汇流固定绿。普通母线和馈线固定红，不跟开关状态变。 */
    if (/_U\d+_X|_U\d+_Y|_XBP|_U\d{2}$/.test(k)) return '#00ff00';
    if (!k && n.c === '#8d9191' && Math.abs((Number(n.a) || 0) % 180) < 1 && Number(n.w) > 400) return '#00ff00';
    return '#ce0c18';
  }

  function capNode(n) {
    if (n && n.capDraw) return false;
    var blob = String((n && n.k) || '') + String((n && n.p) || '') + String((n && n.s) || '');
    return blob.indexOf('电容补偿') >= 0;
  }

  function keyAttr(n, i) {
    var k = n.k ? ' data-hx-key="' + esc(n.k) + '"' : '';
    return ' data-i="' + i + '"' + k;
  }

  /* 工程里同一开关常导出两把刀（分/合各一帧、坐标略偏）。有 k 用 k；无 k 用 dev，避免叠绘红绿。 */
  function knifePickId(n) {
    if (n.t === 'S') return 'S:' + n.x + ',' + n.y;
    if (n.t !== 'I') return '';
    if (n.k) return 'I:' + n.k;
    if (n.dev) return 'I:dev:' + n.dev;
    return 'I:' + n.x + ',' + n.y;
  }

  function knifeInner(decor) {
    /* 合闸红：半圆弧留着，竖线从母线连通到底，不画斜刀片。分闸绿：竖线在下触点处断开。刀片是同一条直线：上端在上方铰接圆点圆心 (18,44)，下端在半圆弧弧长 2/3 处、下部 1/3 分界 (39.51,70.92)，不贴上沿、不贴正中。 */
    var arc = '<path d="M30 52 A13.2 13.2 0 0 1 30 78"/>';
    var dots = '<circle cx="18" cy="14" r="4.15" fill="currentColor" stroke="none"/>' +
      '<circle cx="18" cy="44" r="4.15" fill="currentColor" stroke="none"/>' +
      '<circle cx="18" cy="86" r="4.15" fill="currentColor" stroke="none"/>';
    var bladeLine = '<line x1="18" y1="44" x2="39.51" y2="70.92" stroke-linecap="butt"/>';
    var blade = '<svg class="hx-blade" viewBox="-6 0 48 124" preserveAspectRatio="xMidYMid meet">' +
      '<g class="hx-k-closed" fill="none" stroke="currentColor" stroke-width="2.15" stroke-linecap="round" stroke-linejoin="round">' +
        '<line x1="18" y1="0" x2="18" y2="114" stroke-linecap="butt"/>' +
        arc +
        dots +
      '</g>' +
      '<g class="hx-k-open" fill="none" stroke="currentColor" stroke-width="2.15" stroke-linecap="round" stroke-linejoin="round">' +
        '<line x1="18" y1="0" x2="18" y2="44" stroke-linecap="butt"/>' +
        '<line x1="18" y1="86" x2="18" y2="114" stroke-linecap="butt"/>' +
        arc +
        bladeLine +
        dots +
      '</g>' +
      '<circle cx="18" cy="114" r="4.2" fill="#ffffff"/>' +
      '</svg>';
    if (!decor) return blade;
    return blade +
      '<div class="hx-k-out">OUT</div>' +
      '<svg class="hx-k-tri" viewBox="0 0 20 14" aria-hidden="true"><polygon points="1.5,1 18.5,1 10,13" fill="#ffffff"/></svg>';
  }

  function knifeHtml(n, i, decor, page) {
    var power = page && page.kind === 'power';
    var cls = 'hx-abs hx-sym is-open';
    if (decor && power) cls += ' hx-sym-has-out';
    else if (decor) cls += ' hx-sym-full';
    if (knifeHorizontal(n)) cls += ' hx-sym-h';
    if (n.hfull) cls += ' hx-sym-hf';
    var inner = knifeInner(decor && !power);
    if (decor && power) {
      inner = knifeInner(false) +
        '<svg class="hx-k-tri hx-k-tri-pwr" viewBox="0 0 20 14" aria-hidden="true"><polygon points="1.5,1 18.5,1 10,13" fill="#ffffff"/></svg>';
    }
    return '<div class="' + cls + '" data-hx-knife' + keyAttr(n, i) + ' style="color:#00ff00;' + pos(n) + '">' + inner + '</div>';
  }

  var WIRE_JUMP_R = 7;
  var WIRE_JUMP_GAP = 5;
  var WIRE_END_TOL = 4;

  function axisAlignedSeg(n) {
    var e = ends(n);
    var x1 = e[0];
    var y1 = e[1];
    var x2 = e[2];
    var y2 = e[3];
    if (Math.abs(y1 - y2) < 1.5) {
      return { kind: 'h', y: (y1 + y2) / 2, x0: Math.min(x1, x2), x1: Math.max(x1, x2), n: n };
    }
    if (Math.abs(x1 - x2) < 1.5) {
      return { kind: 'v', x: (x1 + x2) / 2, y0: Math.min(y1, y2), y1: Math.max(y1, y2), n: n };
    }
    return null;
  }

  function wireCrossPoint(a, b) {
    if (!a || !b || a.kind === b.kind) return null;
    var h = a.kind === 'h' ? a : b;
    var v = a.kind === 'v' ? a : b;
    var ix = v.x;
    var iy = h.y;
    if (ix <= h.x0 + WIRE_END_TOL || ix >= h.x1 - WIRE_END_TOL) return null;
    if (iy <= v.y0 + WIRE_END_TOL || iy >= v.y1 - WIRE_END_TOL) return null;
    if (Math.abs(ix - h.x0) <= WIRE_END_TOL || Math.abs(ix - h.x1) <= WIRE_END_TOL) return null;
    if (Math.abs(iy - v.y0) <= WIRE_END_TOL || Math.abs(iy - v.y1) <= WIRE_END_TOL) return null;
    return { x: ix, y: iy };
  }

  /** power 页红线正交交叉：后绘线段做小拱绕开（T 接端点不算交叉）。 */
  function computePowerWireJumps(lineItems) {
    var segs = [];
    var i;
    for (i = 0; i < lineItems.length; i++) {
      var ax = axisAlignedSeg(lineItems[i].n);
      if (ax) {
        ax.li = i;
        ax.idx = lineItems[i].i;
        segs.push(ax);
      }
    }
    var jumps = {};
    for (i = 0; i < segs.length; i++) {
      var j;
      for (j = i + 1; j < segs.length; j++) {
        var pt = wireCrossPoint(segs[i], segs[j]);
        if (!pt) continue;
        var target = segs[j].li;
        jumps[target] = jumps[target] || [];
        jumps[target].push(pt);
      }
    }
    return jumps;
  }

  function buildJumpPath(x1, y1, x2, y2, jumps) {
    if (!jumps || !jumps.length) return '';
    var horiz = Math.abs(y1 - y2) < 1.5;
    var vert = Math.abs(x1 - x2) < 1.5;
    if (!horiz && !vert) return '';
    var parts = ['M', x1, y1];
    var g = WIRE_JUMP_GAP;
    var r = WIRE_JUMP_R;
    if (horiz) {
      var y = (y1 + y2) / 2;
      var step = x2 >= x1 ? 1 : -1;
      var xs = jumps.map(function (p) { return p.x; });
      xs.sort(function (a, b) { return step > 0 ? a - b : b - a; });
      var cursor = x1;
      xs.forEach(function (ix) {
        var xBefore = ix - g * step;
        var xAfter = ix + g * step;
        if (step > 0 && (xBefore <= cursor + 0.5 || xAfter >= x2 - 0.5)) return;
        if (step < 0 && (xBefore >= cursor - 0.5 || xAfter <= x2 + 0.5)) return;
        parts.push('L', xBefore, y);
        parts.push('Q', ix, (y - r), xAfter, y);
        cursor = xAfter;
      });
      parts.push('L', x2, y2);
    } else {
      var x = (x1 + x2) / 2;
      var step = y2 >= y1 ? 1 : -1;
      var ys = jumps.map(function (p) { return p.y; });
      ys.sort(function (a, b) { return step > 0 ? a - b : b - a; });
      var cursor = y1;
      ys.forEach(function (iy) {
        var yBefore = iy - g * step;
        var yAfter = iy + g * step;
        if (step > 0 && (yBefore <= cursor + 0.5 || yAfter >= y2 - 0.5)) return;
        if (step < 0 && (yBefore >= cursor - 0.5 || yAfter <= y2 + 0.5)) return;
        parts.push('L', x, yBefore);
        parts.push('Q', (x + r), iy, x, yAfter);
        cursor = yAfter;
      });
      parts.push('L', x2, y2);
    }
    return parts.join(' ');
  }

  function renderLineMarkup(page, n, i, jumps) {
    var e = ends(n);
    var sw = Math.max(1.5, Number(n.sw) || 2);
    var col = esc(lineColor(page, n));
    var d = buildJumpPath(e[0], e[1], e[2], e[3], jumps);
    var vis = d
      ? '<path data-hx-line' + keyAttr(n, i) + ' d="' + d + '" fill="none" stroke="' + col + '" stroke-width="' + sw + '" stroke-linecap="round" stroke-linejoin="round"/>'
      : '<line data-hx-line' + keyAttr(n, i) + ' x1="' + e[0] + '" y1="' + e[1] + '" x2="' + e[2] + '" y2="' + e[3] + '" stroke="' + col + '" stroke-width="' + sw + '" stroke-linecap="round"/>';
    var hit = '';
    if (n.k) {
      hit = d
        ? '<path' + keyAttr(n, i) + ' d="' + d + '" fill="none" stroke="transparent" stroke-width="16"/>'
        : '<line' + keyAttr(n, i) + ' x1="' + e[0] + '" y1="' + e[1] + '" x2="' + e[2] + '" y2="' + e[3] + '" stroke="transparent" stroke-width="16"/>';
    }
    return vis + hit;
  }

  /* 补偿柜静态刀：与一次图刀同一 viewBox「-6 0 48 124」缩进 60×98，上口贴在符号顶边，供母线竖线接入。 */
  function capCabinetSvg() {
    return '<svg viewBox="0 0 120 252" width="100%" height="100%" overflow="visible" aria-hidden="true">' +
      '<g fill="none" stroke="#e8eef6" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round">' +
      '<g transform="translate(41.032,0) scale(0.790323) translate(6,0)" stroke="#00ff00" stroke-width="2.15">' +
      '<line x1="18" y1="0" x2="18" y2="44" stroke-linecap="butt"/>' +
      '<line x1="18" y1="86" x2="18" y2="114" stroke-linecap="butt"/>' +
      '<path d="M30 52 A13.2 13.2 0 0 1 30 78"/>' +
      '<line x1="18" y1="44" x2="39.51" y2="70.92" stroke-linecap="butt"/>' +
      '<circle cx="18" cy="14" r="4.15" fill="#00ff00" stroke="none"/>' +
      '<circle cx="18" cy="44" r="4.15" fill="#00ff00" stroke="none"/>' +
      '<circle cx="18" cy="86" r="4.15" fill="#00ff00" stroke="none"/>' +
      '<circle cx="18" cy="114" r="4.2" fill="#ffffff" stroke="none"/>' +
      '</g>' +
      '<path d="M60 90.1 V106"/>' +
      '<circle cx="44" cy="116" r="6.2"/>' +
      '<circle cx="60" cy="116" r="6.2"/>' +
      '<circle cx="76" cy="116" r="6.2"/>' +
      '<path d="M60 106 V132 M16 132 H104"/>' +
      '<path d="M30 132 V140 L24 148 L36 156 L24 164 L36 172 L30 180 V188"/>' +
      '<path d="M90 132 V140 L96 148 L84 156 L96 164 L84 172 L90 180 V188"/>' +
      '<path d="M30 188 V198 M18 216 L30 198 L42 216 Z M16 222 H44"/>' +
      '<path d="M90 188 V198 M78 216 L90 198 L102 216 Z M76 222 H104"/>' +
      '<path d="M60 132 V158 M50 158 H70 M50 165 H70 M60 165 V184"/>' +
      '<path d="M46 184 H74 M50 189 H70 M54 194 H66"/>' +
      '<circle cx="54" cy="236" r="1.7" fill="#e8eef6" stroke="none"/>' +
      '<circle cx="60" cy="236" r="1.7" fill="#e8eef6" stroke="none"/>' +
      '<circle cx="66" cy="236" r="1.7" fill="#e8eef6" stroke="none"/>' +
      '</g></svg>';
  }

  function symbolMarkup(n) {
    var svg = {
      xfmr: '<svg viewBox="0 0 64 88" width="100%" height="100%" aria-hidden="true"><g fill="none" stroke="#e8eef6" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"><circle cx="26" cy="16" r="15"/><polygon points="26,6 15,24 37,24"/><circle cx="26" cy="46" r="15"/><path d="M26 46 V34 M26 46 L16 56 M26 46 L36 56"/><path d="M26 46 H48 V64 M42 64 H54 M44 68 H52 M46 72 H50"/></g></svg>',
      rect: '<svg viewBox="0 0 26 22" width="100%" height="100%" aria-hidden="true"><g fill="none" stroke="#d24dff" stroke-width="1.4" stroke-linecap="round"><rect x="1.2" y="1.2" width="23.6" height="19.6" rx="1.5"/><path d="M5 8 Q8.5 4.2 13 8 T21 8"/><path d="M6 15.5 H20"/></g></svg>',
      inv: '<svg viewBox="0 0 26 22" width="100%" height="100%" aria-hidden="true"><g fill="none" stroke="#d24dff" stroke-width="1.4" stroke-linecap="round"><rect x="1.2" y="1.2" width="23.6" height="19.6" rx="1.5"/><path d="M6 6.5 H20"/><path d="M5 15 Q8.5 11.2 13 15 T21 15"/></g></svg>',
      ssw: '<svg viewBox="0 0 28 26" width="100%" height="100%" aria-hidden="true"><g fill="none" stroke="#00e676" stroke-width="1.4" stroke-linejoin="round" stroke-linecap="round"><rect x="1.2" y="1.2" width="25.6" height="23.6" rx="1.5"/><polygon points="7,7 16,13 7,19"/><polygon points="21,7 12,13 21,19"/><path d="M16 13 H22 M6 13 H12"/></g></svg>',
      cap: capCabinetSvg()
    }[n.sym];
    if (!svg) return '';
    return '<div class="hx-abs" style="pointer-events:none;' + pos(n) + '">' + svg + '</div>';
  }

  function html(page) {
    var nodes = page.nodes || [];
    var knifeNodes = [];
    var ki;
    for (ki = 0; ki < nodes.length; ki++) {
      if (nodes[ki].t === 'I' || nodes[ki].t === 'S') knifeNodes.push(nodes[ki]);
    }
    var lines = [];
    var rest = [];
    var powerRedLines = [];
    var pick = {};
    nodes.forEach(function (n, i) {
      if (capNode(n)) return;
      if (n.t !== 'S' && n.t !== 'I') return;
      var kid = knifePickId(n);
      if (!kid) return;
      if (pick[kid] == null || n.st === '1') pick[kid] = i;
    });
    nodes.forEach(function (n, i) {
      if (capNode(n)) return;
      if (n.t === 'O') return;
      if (n.t === 'L') {
        var col = lineColor(page, n);
        if (page.kind === 'power' && (col === '#ce0c18' || n.fixc)) {
          powerRedLines.push({ n: n, i: i });
        } else {
          lines.push(renderLineMarkup(page, n, i, null));
        }
        return;
      }
      if (n.t === 'S' || n.t === 'I') {
        if (pick[knifePickId(n)] !== i) return;
        var showOut = n.decor && (page.kind === 'column' ? n.t === 'S' : n.t === 'I' || n.t === 'S');
        rest.push(knifeHtml(n, i, showOut, page));
        return;
      }
      if (n.t === 'T') {
        var pe = n.k ? '' : 'pointer-events:none;';
        var wrap = n.s && String(n.s).indexOf('\n') >= 0 ? 'white-space:normal;text-align:center;line-height:1.15;' : '';
        var align = TEXT_ALIGN[n.al] || '';
        var drop = page.kind === 'power' ? underKnifeTextDrop(n, knifeNodes) : 0;
        var txCls = 'hx-abs hx-tx' + (drop ? ' hx-under-knife' : '');
        rest.push('<div class="' + txCls + '"' + keyAttr(n, i) + ' style="' + pe + wrap + align + 'color:' + esc(n.c || '#fff') + ';font-size:' + (n.fs || 18) + 'px;' + pos(n) + '">' + esc(n.s || '').replace(/\n/g, '<br>') + '</div>');
        return;
      }
      if (n.t === 'B') {
        if (n.sym) {
          rest.push(symbolMarkup(n));
          return;
        }
        var frame = page.kind === 'power' ? 'hx-frame' : 'hx-glow';
        rest.push('<div class="hx-abs ' + frame + '" style="' + pos(n) + '"></div>');
      }
    });
    if (powerRedLines.length) {
      var jumpMap = computePowerWireJumps(powerRedLines);
      powerRedLines.forEach(function (item, li) {
        lines.push(renderLineMarkup(page, item.n, item.i, jumpMap[li] || null));
      });
    }
    var x0 = page.x0 || 0;
    var bw = (page.w || 3150) - x0;
    var bh = page.h || 1400;
    return '<div class="hx-fit" id="hxFit"><div class="hx-board" style="width:' + bw + 'px;height:' + bh + 'px">' +
      '<div class="hx-shift" style="left:' + (-Number(x0)) + 'px;width:' + (page.w || 3150) + 'px;height:' + bh + 'px">' +
      '<svg class="hx-wires" width="' + (page.w || 3150) + '" height="' + bh + '">' + lines.join('') + '</svg>' +
      rest.join('') +
      '</div></div></div><div id="hxHover" class="hx-hover" hidden></div>';
  }

  /* 图3：一扇柜，上面进线，下面支路。三相占一格，纯数字单相单独占一格。 */
  function cabinet(spec) {
    var nodes = [];
    var title = spec.incoming || '进线';
    nodes.push({ t: 'L', x: 780, y: 8, w: 36, h: 8, a: 90, sw: 3, c: '#ff3b3b', z: 2 });
    nodes.push({
      t: 'T', x: 700, y: 16, w: 200, h: 28, s: title, k: spec.mainKey, c: '#fff', fs: 18, z: 6
    });
    if (spec.mainPoint) {
      nodes.push({
        t: 'S', x: 772, y: 46, w: 56, h: 78, k: spec.mainKey,
        p: spec.mainPoint, dev: spec.mainDev, z: 5
      });
      nodes.push({ t: 'L', x: 798, y: 124, w: 4, h: 48, a: 90, sw: 2, c: '#ff3b3b', z: 2 });
    }
    var x = 28;
    var y = 146;
    var rowH = 196;
    var maxW = 1560;
    var rowReady = false;
    function openRow() {
      if (rowReady) return;
      nodes.push({ t: 'L', x: 20, y: y + 26, w: 1560, h: 4, a: 0, sw: 1.35, c: '#ffffff', z: 2 });
      rowReady = true;
    }
    (spec.branches || []).forEach(function (br) {
      var phases = br.phases || [];
      var single = !phases.length;
      var w = single ? 108 : 252;
      if (x + w > maxW) {
        x = 28;
        y += rowH;
        rowReady = false;
      }
      openRow();
      var busY = y + 26;
      nodes.push({
        t: 'T', x: x, y: y, w: w, h: 22,
        s: String(br.num), k: single ? br.key : '', c: '#fff', fs: 16, z: 6
      });
      if (!single) nodes.push({ t: 'B', x: x + 6, y: y + 34, w: w - 12, h: 92, z: 3 });
      var kw = 46;
      var kh = 132;
      var ky = y + 28;
      var ly = y + 164;
      if (single) {
        nodes.push({
          t: 'L', x: x + w / 2 - 2, y: busY, w: 4, h: Math.max(4, ky - busY), a: 90, sw: 1.35, c: '#ffffff', z: 2
        });
        nodes.push({
          t: 'S', x: x + (w - kw) / 2, y: ky, w: kw, h: kh, decor: 1,
          k: br.key, p: br.p, dev: br.dev, z: 5
        });
        nodes.push({
          t: 'T', x: x, y: ly, w: w, h: 20,
          s: String(br.num), k: br.key, c: '#fff', fs: 13, z: 6
        });
      } else {
        var gap = (w - 12 - phases.length * kw) / (phases.length + 1);
        var shared = br.p && String(br.dev || '').indexOf('数据机房') === 0;
        phases.forEach(function (ph, pi) {
          var kx = x + 6 + gap + pi * (kw + gap);
          var pk = (br.phaseKey && br.phaseKey[ph]) || '';
          var pp = shared ? br.p : ((br.phaseP && br.phaseP[ph]) || br.p || '');
          var pd = shared ? br.dev : ((br.phaseDev && br.phaseDev[ph]) || br.dev || '');
          nodes.push({
            t: 'L', x: kx + kw / 2 - 2, y: busY, w: 4, h: Math.max(4, ky - busY), a: 90, sw: 1.35, c: '#ffffff', z: 2
          });
          nodes.push({ t: 'S', x: kx, y: ky, w: kw, h: kh, decor: 1, k: pk, p: pp, dev: pd, z: 5 });
          nodes.push({
            t: 'T', x: kx - 12, y: ly, w: kw + 24, h: 20,
            s: br.num + ph, k: pk, c: '#fff', fs: 13, z: 6
          });
        });
      }
      x += w + 16;
    });
    return {
      id: spec.id,
      name: spec.name,
      kind: 'column',
      group: spec.room,
      cab: spec.cab,
      w: 1600,
      h: y + rowH,
      x0: 0,
      nodes: nodes
    };
  }

  w.HXDraw = { html: html, cabinet: cabinet, lineColor: lineColor };
})(window);
