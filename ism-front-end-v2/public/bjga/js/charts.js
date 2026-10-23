/* =========================================================
   图表引擎：纯 SVG 字符串生成，零依赖
   ========================================================= */
(function (A) {
  'use strict';

  var uid = 0;
  function gid(p) { uid += 1; return p + uid; }

  var C = {};

  /* ---------- 折线 / 面积图 ---------- */
  /**
   * opts: { w,h, series:[{name,color,data,area,fill}], labels:[], min,max, unit, pad, showLast }
   */
  C.line = function (o) {
    o = o || {};
    var w = o.w || 396, h = o.h || 200;
    var padL = o.padL == null ? 34 : o.padL, padR = o.padR == null ? 12 : o.padR;
    var padT = o.padT == null ? 14 : o.padT, padB = o.padB == null ? 24 : o.padB;
    var ss = (o.series || []).filter(function (s) { return s.data && s.data.length; });
    if (!ss.length) return '<svg viewBox="0 0 ' + w + ' ' + h + '"></svg>';

    var all = [];
    ss.forEach(function (s) { s.data.forEach(function (v) { if (isFinite(v)) all.push(v); }); });
    var mn = o.min != null ? o.min : Math.min.apply(null, all);
    var mx = o.max != null ? o.max : Math.max.apply(null, all);
    if (mx === mn) { mx = mn + 1; }
    var span = mx - mn;
    mn -= span * 0.12; mx += span * 0.12; span = mx - mn;

    /* 每个系列可自带 norm:[min,max]，实现双 Y 轴视觉 */
    ss.forEach(function (s) {
      if (s.norm) { s._mn = s.norm[0]; s._mx = s.norm[1]; }
      else { s._mn = mn; s._mx = mx; }
      var d = s._mx - s._mn || 1;
      s._mn -= d * 0.12; s._mx += d * 0.12; s._span = s._mx - s._mn;
    });
    if (ss[0].norm) { mn = ss[0]._mn; mx = ss[0]._mx; span = ss[0]._span; }

    var iw = w - padL - padR, ih = h - padT - padB;
    var n = ss[0].data.length;
    var X = function (i) { return padL + (n <= 1 ? iw / 2 : iw * i / (n - 1)); };
    var Y = function (v, s) { return padT + ih - (v - s._mn) / s._span * ih; };

    var g = ['<svg viewBox="0 0 ' + w + ' ' + h + '" width="100%" height="100%" fill="none" preserveAspectRatio="none">'];
    g.push('<defs>');
    ss.forEach(function (s, i) {
      if (!s.area) return;
      var id = gid('ga');
      s._g = id;
      g.push('<linearGradient id="' + id + '" x1="0" y1="0" x2="0" y2="1">' +
        '<stop offset="0" stop-color="' + s.color + '" stop-opacity="' + (s.areaOp || 0.26) + '"/>' +
        '<stop offset="1" stop-color="' + s.color + '" stop-opacity="0"/></linearGradient>');
    });
    g.push('</defs>');

    /* 网格 + Y 轴刻度 */
    var lines = o.grid == null ? 4 : o.grid;
    for (var k = 0; k <= lines; k++) {
      var y = padT + ih * k / lines;
      g.push('<path d="M' + padL + ' ' + y.toFixed(1) + 'H' + (w - padR) + '" stroke="' +
        (k === lines ? '#1B4076' : '#12305E') + '" stroke-width="1"/>');
      if (o.yLabels !== false) {
        var v = mx - span * k / lines;
        g.push('<text x="' + (padL - 6) + '" y="' + (y + 3).toFixed(1) +
          '" text-anchor="end" font-family="sans-serif" font-size="9" fill="#5E7FA8">' +
          (Math.round(v * 10) / 10) + '</text>');
      }
    }

    /* X 轴标签 */
    if (o.labels && o.labels.length) {
      var step = Math.max(1, Math.ceil(o.labels.length / (o.labelCount || 7)));
      o.labels.forEach(function (lb, i) {
        if (i % step !== 0 && i !== o.labels.length - 1) return;
        var anchor = i === 0 ? 'start' : (i === o.labels.length - 1 ? 'end' : 'middle');
        g.push('<text x="' + X(i).toFixed(1) + '" y="' + (h - padB + 14) +
          '" text-anchor="' + anchor + '" font-family="sans-serif" font-size="9" fill="#5E7FA8">' +
          A.esc(lb) + '</text>');
      });
    }

    /* 系列 */
    ss.forEach(function (s) {
      var pts = [], d = '';
      s.data.forEach(function (v, i) {
        if (!isFinite(v)) return;
        var x = X(i), y = Y(v, s);
        pts.push([x, y]);
        d += (d ? ' L' : 'M') + x.toFixed(1) + ' ' + y.toFixed(1);
      });
      if (s.area && pts.length) {
        g.push('<path d="' + d + ' L' + pts[pts.length - 1][0].toFixed(1) + ' ' + (padT + ih) +
          ' L' + pts[0][0].toFixed(1) + ' ' + (padT + ih) + ' Z" fill="url(#' + s._g + ')"/>');
      }
      g.push('<path d="' + d + '" stroke="' + s.color + '" stroke-width="' + (s.width || 2) +
        '" stroke-linecap="round" stroke-linejoin="round"/>');
      if (o.dots !== false && pts.length) {
        pts.forEach(function (p, i) {
          if (pts.length > 26 && i % 2) return;
          g.push('<circle cx="' + p[0].toFixed(1) + '" cy="' + p[1].toFixed(1) + '" r="2.5" fill="' + s.color + '" fill-opacity=".9"/>');
        });
      }
      if (o.showLast !== false && pts.length) {
        var lp = pts[pts.length - 1];
        g.push('<circle cx="' + lp[0].toFixed(1) + '" cy="' + lp[1].toFixed(1) + '" r="7" stroke="' + s.color + '" stroke-opacity=".5" fill="none"/>');
        g.push('<circle cx="' + lp[0].toFixed(1) + '" cy="' + lp[1].toFixed(1) + '" r="3.5" fill="#FFFFFF"/>');
      }
    });

    g.push('</svg>');
    return g.join('');
  };

  /* ---------- 环形进度 ---------- */
  C.ring = function (o) {
    o = o || {};
    var size = o.size || 150, r = o.r || (size / 2 - 12), sw = o.sw || 12;
    var c = size / 2, circ = 2 * Math.PI * r;
    var f = A.clamp(o.rate == null ? 0 : o.rate, 0, 1);
    var id = gid('rg');
    var col = o.color || '#00D8FF', col2 = o.color2 || '#2E7CF6';
    var s = '<svg viewBox="0 0 ' + size + ' ' + size + '" fill="none" width="' + size + '" height="' + size + '">';
    s += '<defs><linearGradient id="' + id + '" x1="0" y1="0" x2="1" y2="1">' +
      '<stop offset="0" stop-color="' + col + '"/><stop offset="1" stop-color="' + col2 + '"/></linearGradient></defs>';
    s += '<circle cx="' + c + '" cy="' + c + '" r="' + r + '" stroke="#163A6E" stroke-width="' + sw + '"/>';
    s += '<circle cx="' + c + '" cy="' + c + '" r="' + r + '" stroke="url(#' + id + ')" stroke-width="' + sw +
      '" stroke-linecap="round" stroke-dasharray="' + (circ * f).toFixed(2) + ' ' + (circ * 2).toFixed(2) +
      '" transform="rotate(-90 ' + c + ' ' + c + ')"/>';
    if (o.knob !== false) {
      var ang = 360 * f * Math.PI / 180;
      var kx = c + r * Math.sin(ang), ky = c - r * Math.cos(ang);
      s += '<circle cx="' + kx.toFixed(2) + '" cy="' + ky.toFixed(2) + '" r="' + (sw / 2 + 2) + '" stroke="' + col + '" stroke-opacity=".5" fill="none"/>';
      s += '<circle cx="' + kx.toFixed(2) + '" cy="' + ky.toFixed(2) + '" r="' + (sw / 5) + '" fill="#7FE9FF"/>';
    }
    s += '</svg>';
    return s;
  };

  /* ---------- 仪表盘（270°） ---------- */
  C.gauge = function (o) {
    o = o || {};
    var v = o.value || 0, mn = o.min == null ? 1 : o.min, mx = o.max == null ? 2 : o.max;
    var f = A.clamp((v - mn) / (mx - mn || 1), 0, 1);
    var w = o.w || 200, h = o.h || 176, c = 100, r = 78;
    var id = gid('gg');
    var ang = (225 + 270 * f) * Math.PI / 180;
    var ex = c + r * Math.sin(ang), ey = c - r * Math.cos(ang);
    var sx = c + r * Math.sin(225 * Math.PI / 180), sy = c - r * Math.cos(225 * Math.PI / 180);
    var exx = c + r * Math.sin(135 * Math.PI / 180), eyy = c - r * Math.cos(135 * Math.PI / 180);
    /* 270° 仪表：扫过超过 180° 才用大弧。f>0.5 会在 PUE≈1.5 时把绿弧甩到另一条圆上 */
    var large = 270 * f > 180 ? 1 : 0;
    return '<svg viewBox="0 0 ' + w + ' ' + h + '" fill="none">' +
      '<defs><linearGradient id="' + id + '" x1="0" y1="1" x2="1" y2="0">' +
      '<stop offset="0" stop-color="#00E39A"/><stop offset="1" stop-color="#00D8FF"/></linearGradient></defs>' +
      '<path d="M' + sx.toFixed(2) + ' ' + sy.toFixed(2) + ' A' + r + ' ' + r + ' 0 1 1 ' + exx.toFixed(2) + ' ' + eyy.toFixed(2) +
      '" stroke="#163A6E" stroke-width="12" stroke-linecap="round"/>' +
      '<path d="M' + sx.toFixed(2) + ' ' + sy.toFixed(2) + ' A' + r + ' ' + r + ' 0 ' + large + ' 1 ' +
      ex.toFixed(2) + ' ' + ey.toFixed(2) + '" stroke="url(#' + id + ')" stroke-width="12" stroke-linecap="round"/>' +
      '<circle cx="' + ex.toFixed(2) + '" cy="' + ey.toFixed(2) + '" r="8" stroke="#00D8FF" stroke-opacity=".4"/>' +
      '<circle cx="' + ex.toFixed(2) + '" cy="' + ey.toFixed(2) + '" r="3.5" fill="#7FE9FF"/></svg>';
  };

  /* ---------- 纵向柱状 ---------- */
  C.bars = function (o) {
    o = o || {};
    var w = o.w || 396, h = o.h || 180, data = o.data || [];
    if (!data.length) return '<svg viewBox="0 0 ' + w + ' ' + h + '"></svg>';
    var padT = 14, padB = o.labels === false ? 8 : 30, padL = 8, padR = 8;
    var ih = h - padT - padB, iw = w - padL - padR;
    var mx = Math.max.apply(null, data.map(function (d) { return d.value; })) || 1;
    mx = mx * 1.15;
    var bw = iw / data.length * 0.56, gap = iw / data.length;
    var id = gid('bg');
    var s = '<svg viewBox="0 0 ' + w + ' ' + h + '" width="100%" height="100%" fill="none" preserveAspectRatio="none">';
    s += '<defs><linearGradient id="' + id + '" x1="0" y1="0" x2="0" y2="1">' +
      '<stop offset="0" stop-color="' + (o.color || '#00D8FF') + '"/>' +
      '<stop offset="1" stop-color="' + (o.color2 || '#2E7CF6') + '" stop-opacity=".35"/></linearGradient></defs>';
    for (var k = 0; k <= 3; k++) {
      var y = padT + ih * k / 3;
      s += '<path d="M' + padL + ' ' + y.toFixed(1) + 'H' + (w - padR) + '" stroke="' + (k === 3 ? '#1B4076' : '#12305E') + '" stroke-width="1"/>';
    }
    data.forEach(function (d, i) {
      var bh = Math.max(2, d.value / mx * ih);
      var x = padL + gap * i + (gap - bw) / 2;
      var y = padT + ih - bh;
      s += '<rect x="' + x.toFixed(1) + '" y="' + y.toFixed(1) + '" width="' + bw.toFixed(1) +
        '" height="' + bh.toFixed(1) + '" rx="2" fill="url(#' + id + ')"/>';
      if (o.value !== false) {
        s += '<text x="' + (x + bw / 2).toFixed(1) + '" y="' + (y - 4).toFixed(1) +
          '" text-anchor="middle" font-family="sans-serif" font-size="9" fill="#9FC3E8">' + d.value + '</text>';
      }
      if (o.labels !== false) {
        s += '<text x="' + (x + bw / 2).toFixed(1) + '" y="' + (h - padB + 14) +
          '" text-anchor="middle" font-family="sans-serif" font-size="9" fill="#5E7FA8">' + A.esc(d.label) + '</text>';
      }
    });
    s += '</svg>';
    return s;
  };

  /* ---------- 多段环形（占比） ---------- */
  C.donut = function (o) {
    o = o || {};
    var size = o.size || 160, r = o.r || 58, sw = o.sw || 18;
    var c = size / 2, data = o.data || [];
    var total = data.reduce(function (a, d) { return a + d.value; }, 0) || 1;
    var circ = 2 * Math.PI * r, acc = 0;
    var s = '<svg viewBox="0 0 ' + size + ' ' + size + '" fill="none" width="' + size + '" height="' + size + '">';
    s += '<circle cx="' + c + '" cy="' + c + '" r="' + r + '" stroke="#12305E" stroke-width="' + sw + '"/>';
    data.forEach(function (d) {
      var len = circ * d.value / total;
      s += '<circle cx="' + c + '" cy="' + c + '" r="' + r + '" stroke="' + d.color +
        '" stroke-width="' + sw + '" stroke-dasharray="' + Math.max(0, len - 1.5).toFixed(2) + ' ' + (circ - len + 1.5).toFixed(2) +
        '" stroke-dashoffset="' + (-acc).toFixed(2) + '" transform="rotate(-90 ' + c + ' ' + c + ')"/>';
      acc += len;
    });
    s += '</svg>';
    return s;
  };

  /* ---------- 迷你折线 ---------- */
  C.spark = function (o) {
    o = o || {};
    var w = o.w || 120, h = o.h || 34, d = o.data || [];
    if (d.length < 2) return '<svg viewBox="0 0 ' + w + ' ' + h + '"></svg>';
    var mn = Math.min.apply(null, d), mx = Math.max.apply(null, d);
    if (mx === mn) mx = mn + 1;
    var span = mx - mn;
    var X = function (i) { return i / (d.length - 1) * (w - 2) + 1; };
    var Y = function (v) { return h - 3 - (v - mn) / span * (h - 6); };
    var p = '', a = '';
    d.forEach(function (v, i) {
      var x = X(i), y = Y(v);
      p += (i ? ' L' : 'M') + x.toFixed(1) + ' ' + y.toFixed(1);
    });
    a = p + ' L' + X(d.length - 1).toFixed(1) + ' ' + h + ' L1 ' + h + ' Z';
    var id = gid('sp');
    var col = o.color || '#00D8FF';
    return '<svg class="spark" viewBox="0 0 ' + w + ' ' + h + '" fill="none" width="' + w + '" height="' + h + '">' +
      '<defs><linearGradient id="' + id + '" x1="0" y1="0" x2="0" y2="1">' +
      '<stop offset="0" stop-color="' + col + '" stop-opacity=".3"/><stop offset="1" stop-color="' + col + '" stop-opacity="0"/></linearGradient></defs>' +
      '<path d="' + a + '" fill="url(#' + id + ')"/>' +
      '<path class="line" d="' + p + '" stroke="' + col + '"/>' +
      '<circle cx="' + X(d.length - 1).toFixed(1) + '" cy="' + Y(d[d.length - 1]).toFixed(1) + '" r="2" fill="#fff"/></svg>';
  };

  A.Charts = C;
})(window.App);
