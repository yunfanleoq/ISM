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

  function pos(n) {
    /* 工程里的 z 有到 -1000 的框。加 120 仍是负数，会画到背景后面，红框就看不见。 */
    return 'left:' + n.x + 'px;top:' + n.y + 'px;width:' + n.w + 'px;height:' + n.h + 'px;z-index:' + ((Number(n.z) || 0) + 2000);
  }

  /* 母线、UPS 出线用原图固定色。合闸红、分闸绿只留给开关本身。 */
  function lineColor(page, n) {
    if (!page || page.kind === 'column') return n.c || '#ffffff';
    var k = String(n.k || '');
    /* UPS 进出线、旁路、汇流固定绿。普通母线和馈线固定红，不跟开关状态变。 */
    if (/_U\d+_X|_U\d+_Y|_XBP|_U\d{2}$/.test(k)) return '#00ff00';
    if (!k && n.c === '#8d9191' && Math.abs((Number(n.a) || 0) % 180) < 1 && Number(n.w) > 400) return '#00ff00';
    return '#ce0c18';
  }

  function capNode(n) {
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

  function knifeHtml(n, i, decor) {
    var cls = 'hx-abs hx-sym is-open' + (decor ? ' hx-sym-full' : '');
    return '<div class="' + cls + '" data-hx-knife' + keyAttr(n, i) + ' style="color:#00ff00;' + pos(n) + '">' + knifeInner(decor) + '</div>';
  }

  function html(page) {
    var nodes = page.nodes || [];
    var lines = [];
    var rest = [];
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
        var e = ends(n);
        var sw = Math.max(1.5, Number(n.sw) || 2);
        lines.push('<line data-hx-line' + keyAttr(n, i) + ' x1="' + e[0] + '" y1="' + e[1] + '" x2="' + e[2] + '" y2="' + e[3] + '" stroke="' + esc(lineColor(page, n)) + '" stroke-width="' + sw + '" stroke-linecap="round"/>');
        if (n.k) {
          lines.push('<line' + keyAttr(n, i) + ' x1="' + e[0] + '" y1="' + e[1] + '" x2="' + e[2] + '" y2="' + e[3] + '" stroke="transparent" stroke-width="16"/>');
        }
        return;
      }
      if (n.t === 'S' || n.t === 'I') {
        if (pick[knifePickId(n)] !== i) return;
        rest.push(knifeHtml(n, i, n.t === 'S' && page.kind === 'column' && n.decor));
        return;
      }
      if (n.t === 'T') {
        var pe = n.k ? '' : 'pointer-events:none;';
        var wrap = n.s && String(n.s).indexOf('\n') >= 0 ? 'white-space:normal;text-align:center;line-height:1.15;' : '';
        rest.push('<div class="hx-abs hx-tx"' + keyAttr(n, i) + ' style="' + pe + wrap + 'color:' + esc(n.c || '#fff') + ';font-size:' + (n.fs || 18) + 'px;' + pos(n) + '">' + esc(n.s || '').replace(/\n/g, '<br>') + '</div>');
        return;
      }
      if (n.t === 'B') {
        var frame = page.kind === 'power' ? 'hx-frame' : 'hx-glow';
        rest.push('<div class="hx-abs ' + frame + '" style="' + pos(n) + '"></div>');
      }
    });
    var x0 = page.x0 || 0;
    var bw = (page.w || 3150) - x0;
    var bh = page.h || 1400;
    return '<div class="hx-fit" id="hxFit"><div class="hx-board" style="width:' + bw + 'px;height:' + bh + 'px">' +
      '<div class="hx-shift" style="left:-' + x0 + 'px;width:' + (page.w || 3150) + 'px;height:' + bh + 'px">' +
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
