/* =========================================================
   2.5D 等轴测（仿 3D）视图：楼层 / 机房 / 机柜
   纯 SVG —— 不用 filter、不用 SMIL（合成层会让放大浏览出现黑块）
   ========================================================= */
(function (A) {
  'use strict';

  var KX = 0.8660, KY = 0.5, KXY = KX;

  /* x 横向、y 纵深、z 竖向。
     KXY 单独给出：机柜深度是宽度的近 2 倍，若与 KX 相同，
     侧面会在屏幕上占掉大半，正面（也就是设备面板）反而看不清。 */
  function Pt(x, y, z) { return [x * KX - y * KXY, (x + y) * KY - (z || 0)]; }

  /* 楼层/机房用标准等轴测；机柜把纵深压扁，让 42U 面板主导画面 */
  function proj(kx, ky, kxy) { KX = kx; KY = ky; KXY = (kxy == null ? kx : kxy); }

  /* 亮度缩放：k>1 提亮，k<1 压暗 */
  function clr(hex, k) {
    var h = String(hex).replace('#', '');
    if (h.length === 3) h = h[0] + h[0] + h[1] + h[1] + h[2] + h[2];
    var r = parseInt(h.slice(0, 2), 16), g = parseInt(h.slice(2, 4), 16), b = parseInt(h.slice(4, 6), 16);
    r = Math.max(0, Math.min(255, Math.round(r * k)));
    g = Math.max(0, Math.min(255, Math.round(g * k)));
    b = Math.max(0, Math.min(255, Math.round(b * k)));
    return '#' + ((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1);
  }

  function esc(s) { return String(s == null ? '' : s); }

  /* ---------------- 场景（画家算法） ----------------
     先按几何单位收集图元，输出时统一缩放到目标显示尺寸 ——
     否则 viewBox 只有一两百单位，12px 的字会被放大成几十像素。 */
  function Scene() { this.it = []; this.bb = []; }

  Scene.prototype.poly = function (depth, pts, fill, opt) {
    for (var i = 0; i < pts.length; i++) this.bb.push(pts[i]);
    this.it.push([depth, { k: 'p', pts: pts, fill: fill, o: opt || {} }]);
    return this;
  };

  /* 立方体：顶 / 右 / 前三面 */
  Scene.prototype.box = function (x, y, z, w, d, h, col, opt) {
    opt = opt || {};
    var o = opt.o || {};
    var T = clr(col, opt.topK || 1.34), R = clr(col, opt.rightK || 0.90), F = clr(col, opt.frontK || 0.56);
    this.poly(x + y + z + h,
      [Pt(x, y, z + h), Pt(x + w, y, z + h), Pt(x + w, y + d, z + h), Pt(x, y + d, z + h)], T, o);
    this.poly(x + w + y + z + 0.01,
      [Pt(x + w, y, z + h), Pt(x + w, y + d, z + h), Pt(x + w, y + d, z), Pt(x + w, y, z)], R, o);
    this.poly(x + y + d + z + 0.02,
      [Pt(x, y + d, z + h), Pt(x + w, y + d, z + h), Pt(x + w, y + d, z), Pt(x, y + d, z)], F, o);
    return this;
  };

  /* 平板（只画顶面） */
  Scene.prototype.plate = function (x, y, z, w, d, col, opt) {
    var o = (opt || {}).o || {};
    this.poly(x + y + z, [Pt(x, y, z), Pt(x + w, y, z), Pt(x + w, y + d, z), Pt(x, y + d, z)], col, o);
    return this;
  };

  /* 在某个 y 平面上贴一个矩形（设备面板细节用） */
  Scene.prototype.face = function (y, x, z, w, h, col, opt) {
    var o = opt || {};
    o.o = o.o || {};
    this.poly(x + y + z + (o.lift || 0.5),
      [Pt(x, y, z + h), Pt(x + w, y, z + h), Pt(x + w, y, z), Pt(x, y, z)], col, o);
    return this;
  };

  /* 在某个 x 平面上贴一个矩形（机柜侧板细节用） */
  Scene.prototype.side = function (x, y, z, d, h, col, opt) {
    var o = opt || {};
    o.o = o.o || {};
    this.poly(x + y + z + (o.lift || 0.5),
      [Pt(x, y, z + h), Pt(x, y + d, z + h), Pt(x, y + d, z), Pt(x, y, z)], col, o);
    return this;
  };

  Scene.prototype.line = function (a, b, color, opt) {
    var p1 = Pt(a[0], a[1], a[2]), p2 = Pt(b[0], b[1], b[2]);
    this.bb.push(p1); this.bb.push(p2);
    this.it.push([a[0] + a[1] + a[2] + 900, { k: 'l', a: p1, b: p2, c: color, o: opt || {} }]);
    return this;
  };

  /* 文字不参与等轴测变形，直接落在投影点上，始终置于最上层 */
  Scene.prototype.label = function (x, y, z, txt, cls, anchor, dy) {
    var p = Pt(x, y, z);
    this.bb.push(p);
    this.it.push([1e6, { k: 't', p: p, s: txt, c: cls || 'iso-lb', a: anchor, dy: dy || 0 }]);
    return this;
  };

  Scene.prototype.svg = function (opt) {
    opt = opt || {};
    this.it.sort(function (a, b) { return a[0] - b[0]; });
    var xs = [], ys = [];
    this.bb.forEach(function (p) { xs.push(p[0]); ys.push(p[1]); });
    if (!xs.length) { xs = [0]; ys = [0]; }
    var x0 = Math.min.apply(null, xs), x1 = Math.max.apply(null, xs);
    var y0 = Math.min.apply(null, ys), y1 = Math.max.apply(null, ys);
    var pad = opt.pad || 2;
    var w = Math.max(1, (x1 - x0) + pad * 2), h = Math.max(1, (y1 - y0) + pad * 2);
    /* 缩放到接近真实显示尺寸，字号才是它字面上的大小 */
    var sc = Math.min((opt.w || 1280) / w, (opt.h || 640) / h);
    var ox = x0 - pad, oy = y0 - pad;
    function X(v) { return ((v - ox) * sc).toFixed(1); }
    function Y(v) { return ((v - oy) * sc).toFixed(1); }

    var out = this.it.map(function (e) {
      var d = e[1], i, s;
      if (d.k === 'p') {
        s = '<polygon points="';
        for (i = 0; i < d.pts.length; i++) s += (i ? ' ' : '') + X(d.pts[i][0]) + ',' + Y(d.pts[i][1]);
        s += '" fill="' + d.fill + '"' +
          (d.o.op != null ? ' fill-opacity="' + d.o.op + '"' : '') +
          (d.o.stroke ? ' stroke="' + d.o.stroke + '" stroke-width="' + (d.o.sw || 1) + '"' : '') +
          (d.o.cls ? ' class="' + d.o.cls + '"' : '') +
          (d.o.attr ? ' ' + d.o.attr : '') + '/>';
        return s;
      }
      if (d.k === 'l') {
        return '<line x1="' + X(d.a[0]) + '" y1="' + Y(d.a[1]) + '" x2="' + X(d.b[0]) + '" y2="' + Y(d.b[1]) +
          '" stroke="' + d.c + '" stroke-width="' + (d.o.w || 1) + '"' +
          (d.o.dash ? ' stroke-dasharray="' + d.o.dash + '"' : '') +
          (d.o.op != null ? ' opacity="' + d.o.op + '"' : '') + '/>';
      }
      return '<text x="' + X(d.p[0]) + '" y="' + Y(d.p[1] + 0) + '" dy="' + d.dy + '" class="' + d.c + '"' +
        (d.a ? ' text-anchor="' + d.a + '"' : '') + '>' + esc(d.s) + '</text>';
    }).join('');

    return '<svg class="iso-svg" viewBox="0 0 ' + (w * sc).toFixed(1) + ' ' + (h * sc).toFixed(1) +
      '" preserveAspectRatio="xMidYMid meet">' + out + '</svg>';
  };

  /* ---------------- 配色 ---------------- */
  var C = {
    floor: '#0C1A2B', floorEdge: '#1B3350',
    room: '#12263A', roomMain: '#0F2E44', roomEmpty: '#16202E',
    wall: '#24405C',
    cabOff: '#3A4C60', cabBlue: '#12A6D8', cabYellow: '#E89A2E', cabRed: '#E0524F',
    cold: '#1C68A8', hot: '#A8501F',
    lhd: '#C8A24D', ac: '#12A88C', ups: '#8A6BD1',
    srv: '#2C7CB5', sto: '#7A5CC6', sw: '#12A08C', sec: '#D2703C', oth: '#5A6B82',
    slot: '#1A2A3C', rail: '#22405C', top: '#33506E'
  };
  C.floorEdge = '#1B3350';

  function lvColor(c) {
    var lv = c.lv || (!c.dev ? 'off' : c.comp > 90 ? 'red' : c.comp >= 70 ? 'yellow' : 'blue');
    return lv === 'off' ? C.cabOff : lv === 'red' ? C.cabRed : lv === 'yellow' ? C.cabYellow : C.cabBlue;
  }
  function catColor(cat) {
    if (cat === '存储设备' || cat === '磁盘阵列') return C.sto;
    if (cat === '交换机') return C.sw;
    if (cat === '网络安全设备' || cat === '网闸') return C.sec;
    if (cat === '其他设备' || cat === '其他') return C.oth;
    return C.srv;
  }

  var TYPE_COLOR = {
    '主机房': '#00D8FF', '屏蔽机房': '#00D8FF', '空调间': '#00E39A', '电池室': '#FFA24D',
    'UPS间': '#FFA24D', '配电室': '#2E9CF6', '弱电间': '#8FA8C4', '监测中心': '#B98CFF', '发电间': '#FF7A45'
  };
  function roomNo(n) { var m = String(n).match(/(\d+)/); return m ? Number(m[1]) : 0; }

  var ISO = {};

  /* =========================================================
     L1 楼层 2.5D：楼板 + 房间体块 + 主机房机柜列
     ========================================================= */
  ISO.floor = function (floor) {
    proj(0.866, 0.5);
    var D = window.DCIM;
    var rs = D.rooms.filter(function (r) { return r.floor === floor; });
    var s = new Scene();

    /* 房间物理尺寸由机柜排布反推：列数决定开间，每列柜数决定进深 */
    function size(r) {
      if (!r.cabCount) return { w: 7, d: 6, cols: 0, per: 0 };
      var lay = r.layout;
      var cols = lay ? lay.cols : Math.max(2, Math.round(Math.sqrt(r.cabCount)));
      var per = (lay && lay.per && lay.per.length) ? Math.max.apply(null, lay.per) : Math.ceil(r.cabCount / cols);
      return { w: Math.max(7, cols * 2.7), d: Math.max(6, per * 0.72 + 2.2), cols: cols, per: per };
    }

    function group(list) {
      var x = 0, out = [], maxd = 6;
      list.forEach(function (r) {
        var sz = size(r);
        out.push({ r: r, x: x, sz: sz });
        x += sz.w + 1.4;
        if (sz.d > maxd) maxd = sz.d;
      });
      return { items: out, w: Math.max(7, x - 1.4), d: maxd };
    }

    var north = group(rs.filter(function (r) { return r.zone !== 'S'; }).sort(function (a, b) { return roomNo(a.name) - roomNo(b.name); }));
    var south = group(rs.filter(function (r) { return r.zone === 'S'; }).sort(function (a, b) { return roomNo(a.name) - roomNo(b.name); }));
    var COR = 3.4;
    var yN = 0, yC = north.d, yS = north.d + COR;
    var W = Math.max(north.w, south.w), H = yS + south.d;

    /* 楼板 */
    s.plate(-1.2, -1.2, -0.5, W + 2.4, H + 2.4, C.floor, { o: { stroke: C.floorEdge, sw: 1.4 } });

    /* 走廊 */
    s.plate(-0.6, yC, -0.3, W + 1.2, COR, '#0E2436', { o: { op: 0.9, stroke: '#1D4A66', sw: 1 } });
    s.label(0.2, yC + COR / 2, 0.05, '走 廊 COLD AISLE', 'iso-lb is-dim', 'start', 4);

    function drawRoom(it, y0) {
      var r = it.r, sz = it.sz, x = it.x, y = y0;
      var tc = TYPE_COLOR[r.type] || '#8FA8C4';
      var lv = !r.cabCount ? 'off' : r.comp > 90 ? 'red' : r.comp >= 70 ? 'yellow' : 'blue';
      var edge = lv === 'red' ? C.cabRed : lv === 'yellow' ? C.cabYellow : lv === 'blue' ? C.cabBlue : '#3E5470';

      /* 房间地板 */
      s.plate(x, y, 0, sz.w, sz.d, r.cabCount ? C.roomMain : C.roomEmpty,
        { o: { stroke: clr(edge, 0.72), sw: 1.2, cls: 'iso-room', attr: 'data-room="' + r.name + '"' } });

      if (!r.cabCount) {
        /* 无机柜的功能间：画矮体块代表配电柜 / 电池架 / 空调 */
        var cap = r.type === '配电室' ? '配电柜' : r.type === '电池室' ? '电池架' :
          r.type === '空调间' ? '精密空调' : r.type === 'UPS间' ? 'UPS 主机' : '机柜';
        var n = Math.max(2, Math.min(6, Math.round(sz.w / 2)));
        for (var q = 0; q < n; q++) {
          s.box(x + 0.8 + q * (sz.w - 1.6) / n, y + sz.d * 0.45, 0, (sz.w - 1.6) / n - 0.4, sz.d * 0.34, 1.5,
            clr(tc, 0.30));
        }
        s.label(x + sz.w / 2, y + sz.d * 0.62, 1.85, cap + ' ×' + n, 'iso-lb is-sm', 'middle', 0);
      } else {
        /* 主机房：按列画机柜长条，颜色取该列平均综合率 */
        var lay = r.layout, cols = sz.cols, per = sz.per;
        var cw = 0.66, gap = 0.05;
        var gx = (sz.w - cols * 2.7) / 2;
        for (var k = 0; k < cols; k++) {
          var list = r.cabinets.filter(function (c) { return c.col === k + 1; });
          var perN = list.length || Math.round(per);
          var sum = 0, cnt = 0;
          list.forEach(function (c) { if (c.dev) { sum += c.comp; cnt++; } });
          var comp = cnt ? sum / cnt : 0;
          var cc = !cnt ? C.cabOff : comp > 90 ? C.cabRed : comp >= 70 ? C.cabYellow : C.cabBlue;
          var bx = x + gx + k * 2.7 + 0.5;
          var by = y + (sz.d - perN * (cw + gap)) / 2;
          s.box(bx, by, 0.02, 1.15, Math.max(0.6, perN * (cw + gap) - gap), 2.0, cc, { topK: 1.22, rightK: 0.62 });
          /* 顶面柜位分隔 */
          for (var j = 1; j < perN; j++) {
            s.line([bx, by + j * (cw + gap), 2.02], [bx + 1.15, by + j * (cw + gap), 2.02], clr(cc, 1.6), { w: 0.6, op: 0.5 });
          }
        }
        /* 冷热通道地面提示 */
        if (lay && lay.aisles && lay.aisles.length) {
          for (var ai = 0; ai < lay.aisles.length && ai < cols - 1; ai++) {
            var isCold = lay.aisles[ai].t === 'cold';
            var ax = x + gx + ai * 2.7 + 1.65;
            s.plate(ax, y + 0.9, 0.03, 1.0, sz.d - 1.8, isCold ? C.cold : C.hot, { o: { op: 0.55 } });
          }
        }
      }

      /* 房间标签：浮在房间体块正上方，只留房号 + 一行摘要，告警并进摘要避免三行压字 */
      var no = roomNo(r.name) || r.name;
      var tn = r.type.replace('机房', '').replace('间', '');
      var lblZ = r.cabCount ? 3.35 : 2.62;
      s.label(x + sz.w / 2, y + sz.d / 2, lblZ,
        r.cabCount ? String(no) : (no + ' ' + tn), 'iso-lb is-name', 'middle', 0);
      if (r.cabCount) {
        s.label(x + sz.w / 2, y + sz.d / 2, lblZ, r.type + ' · ' + r.cabCount + ' 柜 · 综合 ' +
          Math.round(r.comp) + '%' + (r.alarmN > 300 ? ' · ▲' + A.num(r.alarmN) : ''),
          'iso-lb is-sm', 'middle', 16);
      }
    }

    north.items.forEach(function (it) { drawRoom(it, yN); });
    south.items.forEach(function (it) { drawRoom(it, yS); });

    return s.svg({ pad: 2.5, w: 1300, h: 590 });
  };

  /* =========================================================
     L2 机房 2.5D：地板 + 机柜列 + 冷热通道 + 列头柜 + 空调
     ========================================================= */
  ISO.room = function (name) {
    proj(0.866, 0.5);
    var D = window.DCIM;
    var r = D.rooms.filter(function (x) { return x.name === name; })[0];
    if (!r) return '';
    var s = new Scene();

    /* 无机柜的功能间（配电室 / 电池室 / 空调间…）：画房间体块与主要设备 */
    if (!r.cabCount) {
      var tc = TYPE_COLOR[r.type] || '#8FA8C4';
      var W0 = 9, D0 = 7;
      s.plate(0, 0, -0.06, W0, D0, C.roomEmpty, { o: { stroke: clr(tc, 0.8), sw: 1.6 } });
      var cap = r.type === '配电室' ? '配电柜' : r.type === '电池室' ? '电池架' :
        r.type === '空调间' ? '精密空调' : r.type === 'UPS间' ? 'UPS 主机' : '机柜';
      var n = Math.max(3, Math.min(7, Math.round(r.dev / 2) || 4));
      for (var q0 = 0; q0 < n; q0++) {
        s.box(0.9 + q0 * ((W0 - 1.8) / n), D0 * 0.28, 0, (W0 - 1.8) / n - 0.35, D0 * 0.42, 2.0, clr(tc, 0.45), { topK: 1.4 });
      }
      s.label(W0 / 2, D0 * 0.34, 2.35, cap + ' ×' + n, 'iso-lb is-sm', 'middle', 0);
      s.label(W0 / 2, D0 + 0.5, 0.1, r.name + ' · ' + r.type + ' · ' + r.dev + ' 台独立设备', 'iso-lb is-name', 'middle', 0);
      return s.svg({ pad: 1.6, w: 960, h: 560 });
    }

    var lay = r.layout || {};
    var CD = 1.15, CW = 0.70, GAP = 0.06, AW = 1.75, H = 2.0;
    var colsMap = {};
    r.cabinets.forEach(function (c) { (colsMap[c.col] = colsMap[c.col] || []).push(c); });
    var colKeys = Object.keys(colsMap).map(Number).sort(function (a, b) { return a - b; });
    var maxPer = Math.max.apply(null, colKeys.map(function (k) { return colsMap[k].length; }));
    var laneD = maxPer * (CW + GAP);

    var W = colKeys.length * (CD + AW) - AW;
    var Dp = laneD + 3.2;

    /* 机房地板 */
    s.plate(-1.0, -1.6, -0.06, W + 2.0, Dp + 2.0, C.roomMain, { o: { stroke: '#1D4A66', sw: 1.4 } });

    /* 列 / 通道 */
    colKeys.forEach(function (key, i) {
      var list = colsMap[key].slice().sort(function (a, b) { return a.seq - b.seq; });
      var x = i * (CD + AW);
      var y0 = 1.0;

      list.forEach(function (c, j) {
        var y = y0 + j * (CW + GAP);
        var col = lvColor(c);
        var o = { cls: 'iso-cab', attr: 'data-cab="' + c.name + '"' };
        s.box(x, y, 0, CD, CW, H, col, { o: o, topK: 1.5 });
        /* 正面：使用率越高，色块越高（与 2D 平面同一套读法） */
        var fh = Math.max(0.12, H * Math.min(1, (c.comp || 0) / 100));
        if (c.dev) s.face(y + CW, x + 0.12, 0, CD - 0.24, fh, clr(col, 1.5), { o: { op: 0.55 } });
      });

      /* 列头柜 */
      s.box(x, y0 - 1.0, 0, CD, 0.72, H, C.lhd, { topK: 1.4 });
      s.label(x + CD / 2, y0 - 1.05, H + 0.55, 'L' + key, 'iso-lb is-lhd', 'middle', 0);

      /* 列间通道 */
      if (i < colKeys.length - 1) {
        var as = (lay.aisles || [])[i] || { t: i % 2 === 0 ? 'cold' : 'hot' };
        var isCold = as.t === 'cold';
        var ax = x + CD;
        s.plate(ax, y0 - 0.2, 0.02, AW, laneD + 0.4, isCold ? C.cold : C.hot, { o: { op: 0.88 } });
        /* 气流箭头：冷通道向机柜、热通道向上 */
        var mid = ax + AW / 2;
        for (var q = 0; q < 3; q++) {
          var yy = y0 + laneD * (0.2 + q * 0.3);
          if (isCold) {
            s.line([mid, yy, 0.06], [ax + 0.18, yy, 0.06], '#7FE9FF', { w: 1, op: 0.8 });
            s.line([mid, yy, 0.06], [ax + AW - 0.18, yy, 0.06], '#7FE9FF', { w: 1, op: 0.8 });
          } else {
            s.line([mid, yy, 0.02], [mid, yy, 1.5], '#FFB27F', { w: 1, op: 0.75, dash: '3 3' });
          }
        }
        s.label(mid, y0 + laneD + 0.5, 0.1, isCold ? '冷通道 COLD' : '热通道 HOT', 'iso-lb is-sm', 'middle', 0);
      }
    });

    /* 精密空调：机房内的靠墙布置，外置的画送风方向 */
    var acN = lay.acIn || 0;
    if (acN) {
      for (var q2 = 0; q2 < acN; q2++) {
        var ax2 = -0.85 + q2 * 0.1;
        s.box(ax2, 1.0 + q2 * (laneD / Math.max(1, acN)), 0, 0.8, 1.2, 2.2, C.ac, { topK: 1.4 });
      }
      s.label(-0.45, 1.0, 2.5, 'AC ×' + acN, 'iso-lb is-sm', 'middle', -2);
    } else {
      var rooms = (lay.acRoom && lay.acRoom.length) ? lay.acRoom.join('、') : '楼层空调间';
      /* 冷源说明放到地台前沿之外，避免压住通道标签 */
      s.label(W / 2, Dp + 0.95, 0.1, '↗ 冷源：' + rooms + ' 送风（地板下送风 / 冷通道封闭）', 'iso-lb is-sm', 'middle', 0);
    }

    s.label(W + 0.6, Dp / 2, 0.1, r.name + ' · ' + colKeys.length + ' 列 × ' + maxPer + ' 柜', 'iso-lb is-name', 'start', 0);

    return s.svg({ pad: 2.0, w: 1180, h: 600 });
  };

  /* =========================================================
     L3 机柜 2.5D：42U 柜体 + 按真实 uH 的 1U/2U/4U 设备
     ========================================================= */
  ISO.rack = function (cabName) {
    proj(1.0, 0.30, 0.55);
    var D = window.DCIM;
    var cab = null, room = null;
    D.rooms.forEach(function (x) {
      x.cabinets.forEach(function (c) { if (c.name === cabName) { cab = c; room = x; } });
    });
    if (!cab) return '';

    var items = D.assets.filter(function (a) { return a.cabinet === cabName; })
      .sort(function (a, b) { return a.uStart - b.uStart; });

    var s = new Scene();
    var CW = 0.62, CD = 1.02, H = 2.0, HU = H / 42;   /* CW=柜宽(x) CD=柜深(y) 1U=4.76cm */
    var DD = 0.09;                                    /* 设备面板厚度：薄，正面才看得清 */
    var DY = CD - 0.05 - DD;                          /* 设备后缘：面板朝 +y（观察者一侧） */

    /* 柜体：底座 / 顶盖 / 两侧立柱 / 后板 */
    s.box(0, 0, -0.06, CW, CD, 0.06, C.rail, { topK: 1.2 });
    s.box(0, 0, H, CW, CD, 0.07, C.top, { topK: 0.92, rightK: 0.74 });
    s.box(0, 0, 0, 0.045, CD, H, C.rail, { topK: 1.15 });
    s.box(CW - 0.045, 0, 0, 0.045, CD, H, C.rail, { topK: 1.15 });
    s.box(0, 0.02, 0, CW, 0.04, H, clr(C.rail, 0.7), { topK: 1.1 });   /* 后板 */

    /* 侧板工艺细节：折边缝 + 上部散热格栅 + 铭牌 —— 否则整块侧面是一块死板 */
    var SX = CW + 0.001, SL = 0.03;
    s.side(SX, 0.05, 0.06, CD - 0.10, H - 0.12, clr(C.rail, 0.80), { o: { op: 1, lift: SL } });
    s.line([SX, 0.05, 0.10], [SX, CD - 0.05, 0.10], '#31536F', { w: 0.8, op: 0.85 });
    s.line([SX, 0.05, H - 0.10], [SX, CD - 0.05, H - 0.10], '#31536F', { w: 0.8, op: 0.85 });
    s.line([SX, 0.05, 0.10], [SX, 0.05, H - 0.10], '#31536F', { w: 0.8, op: 0.85 });
    s.line([SX, CD - 0.05, 0.10], [SX, CD - 0.05, H - 0.10], '#31536F', { w: 0.8, op: 0.85 });
    for (var vg = 0; vg < 7; vg++) {
      s.line([SX, 0.16, H - 0.60 + vg * 0.045], [SX, CD - 0.16, H - 0.60 + vg * 0.045], '#1A3149', { w: 1.1, op: 0.9 });
    }
    s.side(SX, CD * 0.42, H * 0.46, CD * 0.18, 0.10, '#16293C', { o: { op: 0.95, lift: SL + 0.01 } });

    /* 42U 刻度：每 5U 一根，标在柜门左侧 */
    for (var u = 1; u <= 42; u++) {
      if (u % 5 && u !== 1) continue;
      var z = (u - 1) * HU;
      s.line([-0.02, CD + 0.02, z], [-0.16, CD + 0.02, z], '#3E5F80', { w: 0.8, op: 0.85 });
      s.label(-0.20, CD + 0.02, z, u, 'iso-lb is-u', 'end', 3);
    }

    /* 空位：按连续空档铺「空盲板」—— 面是实的，柜内不会透出后板变成玻璃箱 */
    var used = {};
    items.forEach(function (a) { for (var u = a.uStart; u <= a.uEnd; u++) used[u] = 1; });
    var FYD = DY + DD;
    function blankFill(uStart, n) {
      var z0 = (uStart - 1) * HU, hh = n * HU - 0.003;
      s.face(FYD, 0.06, z0, CW - 0.12, hh, '#101E2C', { o: { op: 1 } });
      for (var i = 1; i < n; i++) {
        s.line([0.07, FYD + 0.002, z0 + i * HU], [CW - 0.07, FYD + 0.002, z0 + i * HU], '#2A4661', { w: 0.6, op: 0.5 });
      }
      /* 盲板两端固定螺钉，空柜也有工艺感 */
      s.line([0.10, FYD + 0.003, z0 + 0.02], [0.10, FYD + 0.003, z0 + hh - 0.02], '#3E6183', { w: 1.0, op: 0.85 });
      s.line([CW - 0.10, FYD + 0.003, z0 + 0.02], [CW - 0.10, FYD + 0.003, z0 + hh - 0.02], '#3E6183', { w: 1.0, op: 0.85 });
    }
    var blankRun = 0, blankFrom = 1;
    for (var uu = 1; uu <= 42; uu++) {
      if (!used[uu]) { if (!blankRun) blankFrom = uu; blankRun++; continue; }
      if (blankRun) { blankFill(blankFrom, blankRun); blankRun = 0; }
    }
    if (blankRun) blankFill(blankFrom, blankRun);
    for (var ug = 1; ug <= 42; ug += 5) {
      if (used[ug]) continue;
      s.line([0.06, FYD, (ug - 1) * HU], [CW - 0.06, FYD, (ug - 1) * HU], '#3A5F82', { w: 0.9, op: 0.85 });
    }

    /* 设备：按真实 U 高堆叠（1U / 2U / 4U / 8U 各自高度不同） */
    items.forEach(function (a) {
      if (a.uStart > 42) return;                 /* 超容待调整，不画进柜体 */
      var z = (a.uStart - 1) * HU;
      var h = Math.max(HU * 0.93, a.uH * HU - 0.0035);
      var col = catColor(a.cat);
      var FY = CD - 0.05;                        /* 柜门平面（朝观察者） */
      var fx = 0.06, fw = CW - 0.12;

      s.box(0.06, DY, z, fw, DD, h, col, {
        o: { cls: 'iso-dev', attr: 'data-sn="' + a.sn + '"', frontK: 1.0, rightK: 0.62, topK: 0.72 }
      });

      /* 面板底衬：深一档，细节才有对比 */
      s.face(FY, fx, z, fw, h, clr(col, 0.92), { o: { op: 1, lift: 0.01 } });

      var i, iz = z + h * 0.10, ih = h * 0.80;
      var psu = a.psu === 1 ? 1 : 2;
      var pw = fw * 0.055;
      /* 电源模块：1U 竖排、更高 U 更宽 */
      for (i = 0; i < psu; i++) {
        s.face(FY, fx + 0.010 + i * (pw + 0.007), iz, pw, ih, clr(col, 0.34), { o: { op: 1, lift: 0.02, stroke: '#0B1A28', sw: 0.4 } });
      }

      /* 硬盘位：U 越高盘位越多；1U 面板太矮，只留电源与指示灯 */
      var slots = a.cat === '存储设备' || a.cat === '磁盘阵列'
        ? (a.uH >= 4 ? 8 : 4)
        : (a.uH >= 4 ? 6 : a.uH >= 2 ? 4 : 0);
      if (a.uH < 2) slots = 0;
      if (a.cat === '交换机' || a.cat === '网络安全设备') slots = 0;
      var sx0 = fx + 0.010 + psu * (pw + 0.007) + 0.016;
      var room = fw * 0.50;
      var sw = Math.min(fw * 0.045, (room - slots * 0.005) / Math.max(1, slots));
      for (i = 0; i < slots; i++) {
        s.face(FY, sx0 + i * (sw + 0.005), iz + ih * 0.10, sw, ih * 0.80, clr(col, 0.50), { o: { op: 1, lift: 0.02 } });
      }

      /* 网口：交换机满排，服务器几个 */
      var ports = a.cat === '交换机' ? (a.uH >= 2 ? 24 : 12)
        : a.cat === '网络安全设备' ? 8 : (a.cat === '服务器' ? (a.uH >= 2 ? 4 : 2) : 2);
      var pw2 = Math.min(fw * 0.017, (fw * 0.34) / Math.max(1, ports));
      var px0 = fx + fw - 0.012 - ports * (pw2 + 0.004);
      for (i = 0; i < ports; i++) {
        s.face(FY, px0 + i * (pw2 + 0.004), iz + ih * 0.30, pw2, ih * 0.40, '#08131F', { o: { op: 0.95, lift: 0.02 } });
      }

      /* 运行指示灯（右端一颗绿点） */
      s.face(FY, fx + fw - 0.022, iz + ih * 0.62, 0.012, ih * 0.26, '#25F0A8', { o: { op: 1, lift: 0.03 } });

      /* 设备之间的分隔线，U 位层次更清楚 */
      s.line([fx, FY + 0.001, z + h], [fx + fw, FY + 0.001, z + h], '#0A1622', { w: 0.6, op: 0.55 });
    });

    /* 柜顶：后部走线槽 + 前部双路 PDU（贴在顶盖上，不悬空） */
    s.box(0.10, 0.05, H + 0.07, CW - 0.20, 0.19, 0.032, '#24384C', { topK: 1.12, rightK: 0.8 });
    s.box(0.10, CD - 0.30, H + 0.07, CW - 0.20, 0.11, 0.028, '#2E2750', { topK: 1.15, rightK: 0.8 });
    s.box(0.10, CD - 0.16, H + 0.07, CW - 0.20, 0.09, 0.024, '#20283F', { topK: 1.1, rightK: 0.8 });
    for (var pz = 0; pz < 6; pz++) {
      s.line([0.16 + pz * (CW - 0.32) / 5, CD - 0.30, H + 0.099], [0.16 + pz * (CW - 0.32) / 5, CD - 0.19, H + 0.099],
        '#4B3E7A', { w: 0.9, op: 0.9 });
    }

    /* 柜顶进风口：细格栅，给顶面一点工艺细节 */
    for (var gz = 0; gz < 6; gz++) {
      s.line([0.13, 0.34 + gz * 0.05, H + 0.076], [CW - 0.13, 0.34 + gz * 0.05, H + 0.076], '#22384F', { w: 0.7, op: 0.8 });
    }

    return s.svg({ pad: 0.12, w: 400, h: 640 });
  };

  /* 机柜元信息：柜号 / 台数 / 规格构成 / 超容 —— 交给 HTML 图注排版，SVG 里只留刻度 */
  ISO.rackMeta = function (cabName) {
    var items = window.DCIM.assets.filter(function (a) { return a.cabinet === cabName; });
    var byU = {};
    items.forEach(function (a) { byU[a.uH] = (byU[a.uH] || 0) + 1; });
    var spec = Object.keys(byU).sort(function (a, b) { return a - b; })
      .map(function (k) { return k + 'U×' + byU[k]; }).join(' · ');
    return {
      name: cabName.split('-').slice(-2).join('-'),
      n: items.length,
      spec: spec,
      over: items.filter(function (a) { return a.uStart > 42; }).length
    };
  };

  ISO.Scene = Scene;
  ISO.clr = clr;
  window.ISO = ISO;
})(window.App);
