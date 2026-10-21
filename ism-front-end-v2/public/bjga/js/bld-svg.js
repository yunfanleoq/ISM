/* =========================================================
   楼体线框 + 全息光效 SVG
   ---------------------------------------------------------
   形制（按项目实景截图实测）：
     · 四个体量自左向右逐级抬升，形成明显的台阶式错落
     · 每个体量都有真实进深的屋顶面（不是一条薄边）
     · 最右端露出宽大的侧立面
     · 立面为密集竖向长条窗，四层
     · 全息光效：背景光晕 / 扫描弧环 / 地面扩散波纹 /
                楼后光柱 / 掠过楼体的扫光 / 上浮粒子
   几何全部由「实测比例」参数化。
   ========================================================= */
(function (global) {
  'use strict';

  /* ---------- 画布 / 透视 ---------- */
  var VB_X = 26, VB_Y = -62, VB_W = 1078, VB_H = 482;

  var X0 = 58, X1 = 988;         // 楼体横向范围
  var BL = 340, BR = 372;        // 底面线：左端 / 右端 y
  var KX = 0.72, KY = 0.21;      // 进深方向单位向量（向右后上方）
  var N  = 4;                    // 地上层数
  var ROOF = 80;                 // 屋顶面进深
  var DEP  = 112;                // 最右端侧立面进深

  /* ---------- 体量分段 ----------
     h    = 该体量的基准层高总高（含透视前）
     dd   = 深度偏移（负=更靠近观者，形成前后错落）
     cols = 竖向窗列数
     side = 是否绘制最右端大侧立面 */
  var SEGS = [
    { x0:  58, x1: 149, cols: 4,  h: 168, dd:   0 },
    { x0: 149, x1: 383, cols: 9,  h: 201, dd:  -8 },
    { x0: 383, x1: 720, cols: 11, h: 257, dd:  -4 },
    { x0: 720, x1: 988, cols: 9,  h: 338, dd:   0, side: true }
  ];

  /* ---------- 基础函数 ---------- */
  function r2(n) { return Math.round(n * 10) / 10; }
  function uu(x) { return (x - X0) / (X1 - X0); }
  function gy(x) { return BL + (BR - BL) * uu(x); }            // 地面 y
  function pf(x) { return 0.99 + 0.03 * uu(x); }               // 透视放大率
  function topY(s, x) { return gy(x) - s.h * pf(x); }          // 屋面 y
  function ox(x, d) { return x + KX * (d || 0); }
  function oy(y, d) { return y - KY * (d || 0); }
  function rnd(i) { var v = Math.sin(i * 12.9898 + 4.1414) * 43758.5453; return v - Math.floor(v); }

  /* ---------- 节点 / 粒子配色：青为主，金 / 琥珀 / 紫 / 薄荷点缀 ---------- */
  /* 权重：金最多，其次冰蓝 / 琥珀，紫与薄荷少量点缀 */
  var TN = [['tn-gold', 0.30], ['tn-amber', 0.45], ['tn-ice', 0.58], ['tn-violet', 0.82], ['tn-mint', 1.01]];
  /* p = 非主色（彩色）占比；返回带前导空格的 class 片段，主色返回空串 */
  function tintCls(i, p) {
    if (rnd(i * 12.9898 + 7.7) > (p === undefined ? 0.55 : p)) return '';
    var v = rnd(i * 31.7 + 3.3), k;
    for (k = 0; k < TN.length; k++) { if (v <= TN[k][1]) return ' ' + TN[k][0]; }
    return ' ' + TN[0][0];
  }
  /* 随机闪烁节奏：不同延迟 + 不同周期，避免整齐划一 */
  function twk(i, base, span) {
    return ' style="animation-delay:' + r2(-rnd(i * 5.3 + 1.7) * 4) + 's;--td:' +
           r2(base + rnd(i * 9.1 + 4.2) * span) + 's"';
  }

  /* ---------- 路径 ---------- */
  function winPath(s, cx0, cx1, j) {
    var d = s.dd, a = 0.35, b = 0.10;
    function H(x) { return s.h * pf(x) / N; }
    function yt(x) { return oy(gy(x) - (j + 1 - a) * H(x), d); }
    function yb(x) { return oy(gy(x) - (j + b) * H(x), d); }
    return 'M' + r2(ox(cx0, d)) + ',' + r2(yt(cx0)) +
           'L' + r2(ox(cx1, d)) + ',' + r2(yt(cx1)) +
           'L' + r2(ox(cx1, d)) + ',' + r2(yb(cx1)) +
           'L' + r2(ox(cx0, d)) + ',' + r2(yb(cx0)) + 'Z';
  }

  function slabPath(s, j) {
    var d = s.dd;
    function y(x) { return oy(gy(x) - j * s.h * pf(x) / N, d); }
    return 'M' + r2(ox(s.x0, d)) + ',' + r2(y(s.x0)) +
           'L' + r2(ox(s.x1, d)) + ',' + r2(y(s.x1));
  }

  function facePath(s) {
    var d = s.dd;
    return 'M' + r2(ox(s.x0, d)) + ',' + r2(oy(gy(s.x0), d)) +
           'L' + r2(ox(s.x1, d)) + ',' + r2(oy(gy(s.x1), d)) +
           'L' + r2(ox(s.x1, d)) + ',' + r2(oy(topY(s, s.x1), d)) +
           'L' + r2(ox(s.x0, d)) + ',' + r2(oy(topY(s, s.x0), d)) + 'Z';
  }

  /* 屋顶面：沿进深方向拉出的完整平行四边形 */
  function roofPath(s, dep) {
    var d = s.dd, dp = (dep === undefined ? ROOF : dep);
    return 'M' + r2(ox(s.x0, d)) + ',' + r2(oy(topY(s, s.x0), d)) +
           'L' + r2(ox(s.x1, d)) + ',' + r2(oy(topY(s, s.x1), d)) +
           'L' + r2(ox(s.x1, d + dp)) + ',' + r2(oy(topY(s, s.x1), d + dp)) +
           'L' + r2(ox(s.x0, d + dp)) + ',' + r2(oy(topY(s, s.x0), d + dp)) + 'Z';
  }

  function sideFace(s, dep) {
    var d = s.dd, x = s.x1;
    var ty = topY(s, x), by = gy(x);
    return 'M' + r2(ox(x, d)) + ',' + r2(oy(by, d)) +
           'L' + r2(ox(x, d + dep)) + ',' + r2(oy(by, d + dep)) +
           'L' + r2(ox(x, d + dep)) + ',' + r2(oy(ty, d + dep)) +
           'L' + r2(ox(x, d)) + ',' + r2(oy(ty, d)) + 'Z';
  }

  /* 斜面上的竖窗 */
  function sideWins(s, depth, cols, gid) {
    var x = s.x1, d = s.dd, ty = topY(s, x), by = gy(x), out = '', i, j;
    for (j = 0; j < N; j++) {
      for (i = 0; i < cols; i++) {
        var t0 = i / cols + 0.014, t1 = (i + 1) / cols - 0.014;
        var d0 = d + depth * t0, d1 = d + depth * t1;
        var ya0 = oy(by - (j + 1 - 0.22) * (by - ty) / N, d0);
        var yb0 = oy(by - (j + 0.14) * (by - ty) / N, d0);
        var ya1 = oy(by - (j + 1 - 0.22) * (by - ty) / N, d1);
        var yb1 = oy(by - (j + 0.14) * (by - ty) / N, d1);
        out += '<path class="win' + (((i * 5 + j * 9) % 7 < 2) ? ' is-dim' : '') + '" fill="url(#' + gid + ')" d="M' +
               r2(ox(x, d0)) + ',' + r2(ya0) + 'L' + r2(ox(x, d1)) + ',' + r2(ya1) +
               'L' + r2(ox(x, d1)) + ',' + r2(yb1) + 'L' + r2(ox(x, d0)) + ',' + r2(yb0) + 'Z"/>';
      }
    }
    return out;
  }

  /* ---------- 光效层（楼体之后 / 之前） ---------- */
  function fxBack(gid) {
    var cx = (X0 + X1) / 2 + 6;
    var cy = gy(cx);
    var o = ['<g class="fx">'];
    var i, rx, ry;

    /* 1. 后方大气氛光 */
    o.push('<ellipse cx="' + r2(cx + 12) + '" cy="' + r2(cy - 148) + '" rx="560" ry="256" fill="url(#' + gid + 'halo)"/>');

    /* 2. 放射光线（自楼体中心向后上方散开） */
    o.push('<g class="fx-rays">');
    for (i = 0; i < 68; i++) {
      var a = -Math.PI * (0.012 + 0.976 * i / 67);
      var rr = 290 + rnd(i + 7) * 340;
      o.push('<line x1="' + r2(cx) + '" y1="' + r2(cy - 40) + '" x2="' +
             r2(cx + Math.cos(a) * rr) + '" y2="' + r2(cy - 40 + Math.sin(a) * rr * 0.82) + '"/>');
    }
    o.push('</g>');

    /* 3. 全息穹顶：五层同心椭圆弧 */
    for (i = 0; i < 5; i++) {
      rx = 574 - i * 84;
      ry = rx * 0.50;
      o.push('<path class="fx-dome" d="M' + r2(cx - rx) + ',' + r2(cy + 4) +
             'A' + r2(rx) + ' ' + r2(ry) + ' 0 0 1 ' + r2(cx + rx) + ',' + r2(cy + 4) + '"' +
             ' style="animation-delay:' + r2(-i * 1.15) + 's"/>');
    }

    /* 3b. 雷达扫描线（绕楼体旋转） */
    o.push('<g class="fx-radar"><g class="fx-radar-spin">' +
           '<line x1="' + r2(cx) + '" y1="' + r2(cy) + '" x2="' + r2(cx + 580) + '" y2="' + r2(cy) + '"/>' +
           '</g></g>');

    /* 4. 环绕扫描点环（反向旋转） */
    o.push('<circle class="fx-arc" cx="' + r2(cx + 4) + '" cy="' + r2(cy - 40) + '" r="486"/>');
    o.push('<circle class="fx-arc rev" cx="' + r2(cx + 4) + '" cy="' + r2(cy - 40) + '" r="420"/>');
    o.push('<circle class="fx-arc thin" cx="' + r2(cx + 4) + '" cy="' + r2(cy - 40) + '" r="352"/>');

    /* 4b. 外围节点网络（数据云） */
    var pts = [], pn = 46;
    for (i = 0; i < pn; i++) {
      var ag = Math.PI * (0.04 + 0.92 * rnd(i * 3 + 1));
      var rr2 = 330 + rnd(i * 7 + 2) * 300;
      pts.push([cx + Math.cos(ag) * rr2, cy - 130 - Math.sin(ag) * rr2 * 0.62]);
    }
    var link = '';
    for (i = 0; i < pts.length; i++) {
      var best = -1, bd = 1e9, jj;
      for (jj = 0; jj < pts.length; jj++) {
        if (jj === i) continue;
        var dx = pts[jj][0] - pts[i][0], dy = pts[jj][1] - pts[i][1];
        var dd2 = dx * dx + dy * dy;
        if (dd2 < bd) { bd = dd2; best = jj; }
      }
      if (best >= 0 && bd < 30000) {
        link += 'M' + r2(pts[i][0]) + ',' + r2(pts[i][1]) + 'L' + r2(pts[best][0]) + ',' + r2(pts[best][1]);
      }
    }
    o.push('<path class="fx-link" d="' + link + '"/>');
    for (i = 0; i < pts.length; i++) {
      var nr = r2(1.2 + rnd(i + 11) * 1.7);
      var ntc = tintCls(i * 3 + 5, 0.62);
      o.push('<circle class="fx-node' + ntc + '" cx="' + r2(pts[i][0]) + '" cy="' + r2(pts[i][1]) +
             '" r="' + nr + '"' + twk(i + 3, 2.3, 3.0) + '/>');
      /* 部分节点带四芒星光斑，制造远近错落的闪烁感 */
      if (nr > 2 && rnd(i * 17 + 9) < 0.42) {
        var sx = r2(pts[i][0]), sy = r2(pts[i][1]), sa = r2(nr * 3.1);
        o.push('<path class="fx-star' + ntc + '" d="M' + r2(sx - sa) + ',' + sy + 'L' + r2(sx + sa) + ',' + sy +
               'M' + sx + ',' + r2(sy - sa) + 'L' + sx + ',' + r2(sy + sa) + '"' + twk(i + 11, 1.9, 2.2) + '/>');
      }
    }

    /* 5. 楼后升腾的光柱 */
    o.push('<path class="fx-beam" fill="url(#' + gid + 'beam)" d="M' +
           r2(738) + ',' + r2(topY(SEGS[3], 738) + 2) + 'L' + r2(806) + ',' + r2(topY(SEGS[3], 806) + 2) +
           'L' + r2(796) + ',-34' + 'L' + r2(746) + ',-34' + 'Z"/>');
    o.push('<path class="fx-beam b2" fill="url(#' + gid + 'beam)" d="M' +
           r2(452) + ',' + r2(topY(SEGS[2], 452) + 2) + 'L' + r2(492) + ',' + r2(topY(SEGS[2], 492) + 2) +
           'L' + r2(487) + ',' + r2(topY(SEGS[2], 452) - 96) + 'L' + r2(456) + ',' + r2(topY(SEGS[2], 452) - 96) + 'Z"/>');

    /* 6. 地面全息圆盘 + 同心扩散波纹 */
    o.push('<ellipse class="fx-disc" cx="' + r2(cx) + '" cy="' + r2(cy + 10) + '" rx="556" ry="96"/>');
    o.push('<ellipse class="fx-disc2" cx="' + r2(cx) + '" cy="' + r2(cy + 10) + '" rx="470" ry="80"/>');
    for (i = 0; i < 4; i++) {
      o.push('<ellipse class="fx-ring" cx="' + r2(cx) + '" cy="' + r2(cy + 10) + '" rx="474" ry="82"' +
             ' style="animation-delay:' + (i * 1.15) + 's"/>');
    }

    /* 7. 楼体底部接触光带 + 向上打光 */
    o.push('<path class="fx-base" d="M' + r2(X0 - 8) + ',' + r2(gy(X0) + 2) +
           'L' + r2(ox(X1, DEP) + 8) + ',' + r2(oy(gy(X1), DEP) + 2) +
           '" stroke="url(#' + gid + 'base)" stroke-width="2.6" fill="none" stroke-linecap="round"/>');

    o.push('</g>');
    return o.join('');
  }

  function fxFront(gid) {
    var cx = (X0 + X1) / 2 + 6;
    var o = ['<g class="fx">'];
    var i, px, py;

    /* 掠过楼体的竖向扫光 */
    o.push('<g clip-path="url(#' + gid + 'clip)">' +
           '<rect class="fx-sweep" x="0" y="-80" width="148" height="520" fill="url(#' + gid + 'sweep)"/>' +
           '</g>');

    /* 自下而上的水平扫描光带 */
    o.push('<g clip-path="url(#' + gid + 'clip)">' +
           '<rect class="fx-hscan" x="0" y="0" width="1100" height="44" fill="url(#' + gid + 'hscan)"/>' +
           '</g>');

    /* 环绕轨道光环（前景，穿过楼体） */
    o.push('<ellipse class="fx-orbit" cx="' + r2(cx) + '" cy="' + r2(gy(cx) - 92) + '" rx="640" ry="124"/>');
    o.push('<ellipse class="fx-orbit o2" cx="' + r2(cx) + '" cy="' + r2(gy(cx) - 200) + '" rx="588" ry="98"/>');
    o.push('<ellipse class="fx-orbit o3" cx="' + r2(cx) + '" cy="' + r2(gy(cx) - 300) + '" rx="512" ry="74"/>');

    /* 地面向上漫射的光 */
    o.push('<g class="fx-up" fill="url(#' + gid + 'beam)">');
    [232, 466, 700, 900].forEach(function (bx) {
      var by = gy(bx) + 4;
      o.push('<path d="M' + r2(bx - 30) + ',' + r2(by) + 'L' + r2(bx + 30) + ',' + r2(by) +
             'L' + r2(bx + 62) + ',' + r2(by - 130) + 'L' + r2(bx - 62) + ',' + r2(by - 130) + 'Z"/>');
    });
    o.push('</g>');

    /* 上浮粒子（多色 + 独立闪烁相位） */
    for (i = 0; i < 50; i++) {
      px = X0 + rnd(i) * (X1 - X0);
      py = 40 + rnd(i + 60) * 320;
      o.push('<circle class="fx-dot' + tintCls(i * 5 + 2, 0.66) + '" cx="' + r2(px) + '" cy="' + r2(py) +
             '" r="' + r2(1.1 + rnd(i + 30) * 1.5) + '"' +
             ' style="animation-delay:' + r2(-rnd(i + 90) * 6) + 's;--td:' + r2(4 + rnd(i + 120) * 4) + 's"/>');
    }

    o.push('</g>');
    return o.join('');
  }

  /* ---------- 屋面设备层：女儿墙栏杆 / 检修马道 / 空调机组 / 冷却塔 / 管廊 / 天线 ---------- */
  function roofKit(s, gid) {
    var o = ['<g class="bld-rk">'];
    var d = s.dd, i, k;
    var W = s.x1 - s.x0;

    function px(x, dd) { return r2(ox(x, d + dd)); }
    function py(x, dd, h) { return r2(oy(topY(s, x), d + dd) - (h || 0)); }
    function pt(x, dd, h) { return px(x, dd) + ',' + py(x, dd, h); }
    function quad(p1, p2, p3, p4, cls) {
      return '<path class="' + cls + '" d="M' + p1 + 'L' + p2 + 'L' + p3 + 'L' + p4 + 'Z"/>';
    }

    /* 矮体量屋顶：只画剪影（屋顶面在透视里太窄，画立体反而糊） */
    function flat(cx, dd, w, dp, hh) {
      var x0 = cx - w, x1 = cx + w, d1 = dd + dp;
      return quad(pt(x0, dd, 0), pt(x1, dd, 0), pt(x1, d1, 0), pt(x0, d1, 0), 'rk-shadow') +
             quad(pt(x0, dd, 0), pt(x1, dd, 0), pt(x1, dd, hh), pt(x0, dd, hh), 'rk-face') +
             '<path class="rk-edge" d="M' + pt(x0, dd, hh) + 'L' + pt(x1, dd, hh) +
             'L' + pt(x1, d1, hh) + '"/>';
    }

    /* 屋顶上的立方体：底面投影 / 右侧面 / 顶面 / 正面 / 棱线 */
    function cube(cx, dd, w, dp, hh, cls) {
      var x0 = cx - w, x1 = cx + w, d1 = dd + dp, c = cls || '';
      return quad(pt(x0, dd, 0), pt(x1, dd, 0), pt(x1, d1, 0), pt(x0, d1, 0), 'rk-shadow') +
             quad(pt(x1, dd, 0), pt(x1, d1, 0), pt(x1, d1, hh), pt(x1, dd, hh), 'rk-side' + c) +
             quad(pt(x0, dd, hh), pt(x1, dd, hh), pt(x1, d1, hh), pt(x0, d1, hh), 'rk-top' + c) +
             quad(pt(x0, dd, 0), pt(x1, dd, 0), pt(x1, dd, hh), pt(x0, dd, hh), 'rk-face' + c) +
             '<path class="rk-edge" d="M' + pt(x0, dd, hh) + 'L' + pt(x1, dd, hh) + 'L' + pt(x1, d1, hh) + '"/>';
    }

    /* 1. 女儿墙：前沿矮墙 + 压顶线 */
    var PW = 8;
    o.push(quad(pt(s.x0, 0, 0), pt(s.x1, 0, 0), pt(s.x1, 0, PW), pt(s.x0, 0, PW), 'rk-parapet'));
    o.push('<path class="rk-pline" d="M' + pt(s.x0, 0, PW) + 'L' + pt(s.x1, 0, PW) + '"/>');
    o.push('<path class="rk-pline" d="M' + pt(s.x0, ROOF, PW) + 'L' + pt(s.x1, ROOF, PW) + '"/>');
    o.push('<path class="rk-pline" d="M' + pt(s.x0, 0, PW) + 'L' + pt(s.x0, ROOF, PW) + '"/>');
    o.push('<path class="rk-pline" d="M' + pt(s.x1, 0, PW) + 'L' + pt(s.x1, ROOF, PW) + '"/>');

    /* 1b. 前沿压顶上的栏杆竖柱 */
    var rail = '', rn = Math.max(4, Math.round(W / 15));
    for (k = 0; k <= rn; k++) {
      var rx = s.x0 + W * k / rn;
      rail += 'M' + pt(rx, 0, PW) + 'L' + pt(rx, 0, PW + 6);
    }
    rail += 'M' + pt(s.x0, 0, PW + 6) + 'L' + pt(s.x1, 0, PW + 6);
    o.push('<path class="rk-rail" d="' + rail + '"/>');

    /* 2. 屋面马道：一条贯通小径（双线 + 踏步） */
    var walkD = ROOF * 0.56;
    o.push('<path class="rk-walk" d="M' + pt(s.x0 + 6, walkD) + 'L' + pt(s.x1 - 6, walkD) + '"/>');
    o.push('<path class="rk-walk" d="M' + pt(s.x0 + 6, walkD + 7) + 'L' + pt(s.x1 - 6, walkD + 7) + '"/>');
    var stepsN = Math.max(4, Math.round(W / 26)), steps = '';
    for (k = 1; k < stepsN; k++) {
      var sx2 = s.x0 + 6 + (W - 12) * k / stepsN;
      steps += 'M' + pt(sx2, walkD) + 'L' + pt(sx2, walkD + 7);
    }
    o.push('<path class="rk-step" d="' + steps + '"/>');

    /* 2b. 屋面斜向结构线（网格罩感） */
    var diag = '', dn = Math.max(4, Math.round(W / 62));
    for (k = 1; k < dn; k++) {
      var dx2 = s.x0 + W * k / dn;
      diag += 'M' + pt(dx2, 0, PW) + 'L' + pt(s.x1, ROOF, 0);
    }
    for (k = 1; k < 3; k++) {
      var dx3 = s.x0 + W * k / 3;
      diag += 'M' + pt(s.x0, 0, PW) + 'L' + pt(dx3, ROOF, 0);
    }
    o.push('<path class="rk-diag" d="' + diag + '"/>');

    /* 3. 屋面分格（找平层分仓缝） */
    var grid = '', gn = Math.max(3, Math.round(W / 52));
    for (k = 1; k < gn; k++) {
      var gx = s.x0 + W * k / gn;
      grid += 'M' + pt(gx, 3, 0) + 'L' + pt(gx, ROOF - 3, 0);
    }
    for (k = 1; k < 3; k++) {
      var gd = 8 + (ROOF - 16) * k / 3;
      grid += 'M' + pt(s.x0, gd, 0) + 'L' + pt(s.x1, gd, 0);
    }
    o.push('<path class="rk-grid" d="' + grid + '"/>');

    var rich = s.h >= 230;          // 高层体量才有足够屋顶进深画立体机组
    var hasPlant = W > 118;         // 过窄的体量只留女儿墙与檐口

    /* 4. 前排大型空调机组 */
    var nB = !hasPlant ? 0 : (rich ? Math.max(1, Math.min(4, Math.round(W / 105))) : 1);
    var bw = Math.min(23, W / (nB * 3.4));
    var bdd = ROOF * 0.16;
    for (i = 0; i < nB; i++) {
      var bx = s.x0 + W * (i + 0.5) / nB;
      var bh = rich ? 21 + rnd(i * 7 + s.x1) * 8 : 9;
      o.push(rich ? cube(bx, bdd, bw, bw * 0.9, bh) : flat(bx, bdd, bw, bw * 0.9, bh));
      /* 正面检修格栅 */
      var gr = '';
      for (k = 1; k < 5; k++) {
        var gxx = bx - bw + bw * 2 * k / 5;
        gr += 'M' + pt(gxx, bdd, 2) + 'L' + pt(gxx, bdd, bh - 2);
      }
      o.push('<path class="rk-duct" d="' + gr + '"/>');
      /* 顶面双风机 */
      var f2 = '';
      for (k = 0; k < (rich ? 2 : 1); k++) {
        var fx = bx + (k ? 1 : -1) * bw * 0.45;
        f2 += '<ellipse class="rk-fan" cx="' + px(fx, bdd + bw * 0.45) + '" cy="' +
              py(fx, bdd + bw * 0.45, bh) + '" rx="' + r2(bw * 0.27) + '" ry="' + r2(bw * 0.14) + '"/>';
        f2 += '<path class="rk-blade" d="M' + px(fx, bdd + bw * 0.45) + ',' + py(fx, bdd + bw * 0.45, bh) +
              'l' + r2(bw * 0.30) + ',' + r2(-bw * 0.10) + 'M' + px(fx, bdd + bw * 0.45) + ',' +
              py(fx, bdd + bw * 0.45, bh) + 'l' + r2(-bw * 0.28) + ',' + r2(-bw * 0.12) + '"/>';
      }
      o.push(f2);
    }

    /* 5. 后排小型风机组（沿后缘一字排开） */
    var nS = !hasPlant ? 0 : (rich ? Math.max(2, Math.min(7, Math.round(W / 62))) : Math.max(1, Math.min(3, Math.round(W / 120))));
    var sw = Math.min(13, W / (nS * 2.2));
    var sdd = ROOF - 30;
    for (i = 0; i < nS; i++) {
      var cx2 = s.x0 + W * (i + 0.5) / nS;
      var sh = rich ? 10 + rnd(i * 3 + s.x0) * 5 : 7;
      o.push(rich ? cube(cx2, sdd, sw, sw * 1.5, sh, ' is-sm') : flat(cx2, sdd, sw, sw * 1.5, sh));
      o.push('<ellipse class="rk-fan" cx="' + px(cx2, sdd + sw * 0.75) + '" cy="' +
             py(cx2, sdd + sw * 0.75, sh) + '" rx="' + r2(sw * 0.42) + '" ry="' + r2(sw * 0.20) + '"/>');
    }

    /* 6. 管廊：机组到后排的桥架 + 立管 */
    var duct = '';
    for (i = 0; i < nB; i++) {
      var dx = s.x0 + W * (i + 0.5) / nB;
      var dy = 10;
      duct += 'M' + pt(dx, bdd + bw * 0.9, dy) + 'L' + pt(dx, walkD - 2, dy) +
              'M' + pt(dx + 5, walkD + 9, dy) + 'L' + pt(dx + 5, sdd - 3, dy);
    }
    o.push('<path class="rk-duct" d="' + duct + '"/>');

    /* 7. 出屋面楼梯间（只在最宽体量上） */
    if (rich && W > 250) {
      var tx = s.x0 + W * 0.12, tw = Math.min(20, W * 0.06);
      o.push(cube(tx, ROOF * 0.60, tw, tw * 1.2, 22, ' is-tall'));
      o.push('<path class="rk-edge" d="M' + pt(tx - tw * 0.7, ROOF * 0.60 + tw * 1.3, 22) +
             'L' + pt(tx + tw * 0.7, ROOF * 0.60 + tw * 1.3, 22) + '"/>');
      o.push('<path class="rk-door" d="M' + pt(tx - tw * 0.34, ROOF * 0.60, 0) +
             'L' + pt(tx - tw * 0.34, ROOF * 0.60, 11) + 'L' + pt(tx + tw * 0.34, ROOF * 0.60, 11) +
             'L' + pt(tx + tw * 0.34, ROOF * 0.60, 0) + '"/>');
    }

    /* 8. 天线杆 / 避雷针 + 航空信标（最高体量） */
    if (s.side || s.h > 300) {
      var mx = s.x0 + W * 0.90, mh = 38, mdd = ROOF * 0.30;
      o.push('<path class="rk-mast" d="M' + pt(mx, mdd, 0) + 'L' + pt(mx, mdd, mh) + '"/>');
      for (k = 1; k <= 3; k++) {
        var hy = mh * k / 4.2, half = 5.5 * (4 - k) / 4 + 1.5;
        o.push('<path class="rk-mast thin" d="M' + r2(px(mx, mdd) - half) + ',' + py(mx, mdd, hy) +
               'L' + r2(px(mx, mdd) + half) + ',' + py(mx, mdd, hy) + '"/>');
      }
      /* 斜拉索 */
      o.push('<path class="rk-guy" d="M' + pt(mx, mdd, mh * 0.72) + 'L' + pt(mx - 22, mdd, 0) +
             'M' + pt(mx, mdd, mh * 0.72) + 'L' + pt(mx + 20, mdd + 6, 0) + '"/>');
      o.push('<circle class="rk-beacon" cx="' + px(mx, mdd) + '" cy="' + py(mx, mdd, mh + 2) + '" r="2.2"/>');
      /* 第二根短杆 */
      o.push('<path class="rk-mast thin" d="M' + pt(mx - 40, mdd + 24, 0) + 'L' + pt(mx - 40, mdd + 24, 16) + '"/>');
    }

    /* 9. 屋角航空障碍灯 */
    if (rich) {
      [[s.x0 + 3, 3], [s.x1 - 3, 3], [s.x0 + 3, ROOF - 3]].forEach(function (p, q) {
        o.push('<circle class="rk-beacon b2" cx="' + px(p[0], p[1]) + '" cy="' +
               py(p[0], p[1], PW + 3) + '" r="1.7" style="animation-delay:' + (q * 0.7) + 's"/>');
      });
    }

    o.push('</g>');
    return o.join('');
  }

  /* ---------- 内部构造：半透明玻璃体内可见的机房 / 机柜 / 节点 ---------- */
  var RCACHE = null;
  function roomsOf(f) {
    if (!RCACHE) {
      RCACHE = {};
      var D = (global.DCIM && global.DCIM.rooms) || [];
      D.forEach(function (r) { (RCACHE[r.floor] = RCACHE[r.floor] || []).push(r); });
    }
    return RCACHE[f] || [];
  }

  /* 把某层的机房数按体量宽度分配到四个体量（最大余数法） */
  function allocRooms(total) {
    var ws = [], sum = 0, i, raw = [], base = [];
    for (i = 0; i < SEGS.length; i++) { ws.push(SEGS[i].x1 - SEGS[i].x0); sum += ws[i]; }
    for (i = 0; i < ws.length; i++) { raw.push(ws[i] / sum * total); base.push(Math.floor(raw[i])); }
    var got = 0, order = [];
    for (i = 0; i < base.length; i++) got += base[i];
    for (i = 0; i < raw.length; i++) order.push({ i: i, f: raw[i] - base[i] });
    order.sort(function (a, b) { return b.f - a.f; });
    for (i = 0; got < total; i++, got++) base[order[i % order.length].i]++;
    return base;
  }

  function innerSeg(s, gid, hlIdx, alloc, si) {
    var o = ['<g class="bld-inner">'];
    /* 体内大气氛光（让内部结构像被打亮） */
    o.push('<ellipse class="inner-glow" cx="' + r2(ox((s.x0 + s.x1) / 2, dM)) + '" cy="' +
           r2((yF(s.x0, 0, dM) + yC(s.x1, N - 1, dM)) / 2) + '" rx="' +
           r2((s.x1 - s.x0) * 0.62) + '" ry="' + r2(s.h * 0.42) +
           '" fill="url(#' + gid + 'halo)"/>');
    var dF = s.dd + 13, dB = s.dd + ROOF - 5, dM = (dF + dB) / 2;
    var j, k, c, r;
    var allNodes = [];

    function H(x) { return s.h * pf(x) / N; }
    function yF(x, jj, d) { return oy(gy(x) - (jj + 0.10) * H(x), d); }   // 该层地面
    function yC(x, jj, d) { return oy(gy(x) - (jj + 0.94) * H(x), d); }   // 该层天花

    for (j = 0; j < N; j++) {
      var n = alloc[j][si];
      if (!n) continue;
      var isHL = (j === hlIdx);
      o.push('<g class="flr' + (isHL ? ' is-hl' : '') + '" data-f="' + (j + 1) + '">');

      var cw = (s.x1 - s.x0) / n;
      var list = roomsOf((j + 1) + 'F');
      var nodes = [];


      for (k = 0; k < n; k++) {
        var rx0 = s.x0 + cw * (k + 0.05), rx1 = s.x0 + cw * (k + 0.95);
        var rm = list.length ? list[(si * 9 + k * 5 + j) % list.length] : null;
        var hot = rm ? Math.min(0.06, 0.01 + rm.alarmN / 2200) : 0.03;
        var mx = (rx0 + rx1) / 2;

        /* 房间地板（沿进深的透视面） */
        o.push('<path class="rm-floor" fill="url(#' + gid + 'inF)" d="M' +
               r2(ox(rx0, dF)) + ',' + r2(yF(rx0, j, dF)) +
               'L' + r2(ox(rx1, dF)) + ',' + r2(yF(rx1, j, dF)) +
               'L' + r2(ox(rx1, dB)) + ',' + r2(yF(rx1, j, dB)) +
               'L' + r2(ox(rx0, dB)) + ',' + r2(yF(rx0, j, dB)) + 'Z"/>');

        /* 后墙 */
        o.push('<path class="rm-back" fill="url(#' + gid + 'inB)" d="M' +
               r2(ox(rx0, dB)) + ',' + r2(yF(rx0, j, dB)) +
               'L' + r2(ox(rx1, dB)) + ',' + r2(yF(rx1, j, dB)) +
               'L' + r2(ox(rx1, dB)) + ',' + r2(yC(rx1, j, dB)) +
               'L' + r2(ox(rx0, dB)) + ',' + r2(yC(rx0, j, dB)) + 'Z"/>');

        /* 地板透视线（沿进深 3 道 + 横向 2 道） */
        var grid = '', gi;
        for (gi = 1; gi <= 3; gi++) {
          var gx = rx0 + (rx1 - rx0) * gi / 4;
          grid += 'M' + r2(ox(gx, dF)) + ',' + r2(yF(gx, j, dF)) +
                  'L' + r2(ox(gx, dB)) + ',' + r2(yF(gx, j, dB));
        }
        for (gi = 1; gi <= 2; gi++) {
          var gd = dF + (dB - dF) * gi / 3;
          grid += 'M' + r2(ox(rx0, gd)) + ',' + r2(yF(rx0, j, gd)) +
                  'L' + r2(ox(rx1, gd)) + ',' + r2(yF(rx1, j, gd));
        }
        o.push('<path class="rm-grid" d="' + grid + '"/>');

        /* 后墙机柜侧影竖线 */
        var wall = '', wc = Math.max(5, Math.round((rx1 - rx0) / 11));
        for (gi = 1; gi < wc; gi++) {
          var wx = rx0 + (rx1 - rx0) * gi / wc;
          wall += 'M' + r2(ox(wx, dB)) + ',' + r2(yF(wx, j, dB)) +
                  'L' + r2(ox(wx, dB)) + ',' + r2(yC(wx, j, dB));
        }
        o.push('<path class="rm-wall" d="' + wall + '"/>');

        /* 顶部灯带 */
        o.push('<path class="rm-lite" d="M' + r2(ox(rx0, dB)) + ',' + r2(yC(rx0, j, dB) + 3) +
               'L' + r2(ox(rx1, dB)) + ',' + r2(yC(rx1, j, dB) + 3) + '"/>');

        /* 机柜排：整排框架 + 密集 U 位横纹（最接近真实机柜的表达） */
        var nCol = Math.max(3, Math.min(12, Math.round((rx1 - rx0) / 14)));
        var uw = (rx1 - rx0) / nCol;
        var cabH0 = H(mx) * 0.50;
        for (r = 0; r < 2; r++) {
          var dd = dF + (dB - dF) * (0.09 + r * 0.68);
          var cH = cabH0 * (r ? 0.72 : 1);
          var frame = '', urows = '', hotP = '';

          /* 排框（底边 + 顶边 + 柜缝） */
          frame += 'M' + r2(ox(rx0, dd)) + ',' + r2(yF(rx0, j, dd)) +
                   'L' + r2(ox(rx1, dd)) + ',' + r2(yF(rx1, j, dd));
          frame += 'M' + r2(ox(rx0, dd)) + ',' + r2(yF(rx0, j, dd) - cH) +
                   'L' + r2(ox(rx1, dd)) + ',' + r2(yF(rx1, j, dd) - cH);
          for (c = 1; c < nCol; c++) {
            var sx = rx0 + uw * c;
            frame += 'M' + r2(ox(sx, dd)) + ',' + r2(yF(sx, j, dd)) +
                     'L' + r2(ox(sx, dd)) + ',' + r2(yF(sx, j, dd) - cH);
          }

          /* U 位横纹 */
          var nU = 3;
          for (c = 1; c < nU; c++) {
            var yo = cH * c / nU;
            urows += 'M' + r2(ox(rx0, dd)) + ',' + r2(yF(rx0, j, dd) - yo) +
                     'L' + r2(ox(rx1, dd)) + ',' + r2(yF(rx1, j, dd) - yo);
          }

          /* 告警柜高亮 */
          for (c = 0; c < nCol; c++) {
            if (rnd(j * 43 + si * 17 + k * 11 + c * 5 + r * 3) < hot) {
              var ax = rx0 + uw * (c + 0.06), bx = rx0 + uw * (c + 0.94);
              hotP += 'M' + r2(ox(ax, dd)) + ',' + r2(yF(ax, j, dd)) +
                      'L' + r2(ox(bx, dd)) + ',' + r2(yF(bx, j, dd)) +
                      'L' + r2(ox(bx, dd)) + ',' + r2(yF(bx, j, dd) - cH) +
                      'L' + r2(ox(ax, dd)) + ',' + r2(yF(ax, j, dd) - cH) + 'Z';
            }
          }

          var ro = r ? 0.5 : 1;
          if (hotP) o.push('<path class="cab-hot" opacity="' + ro + '" d="' + hotP + '"/>');
          o.push('<path class="cab-row" opacity="' + ro + '" d="' + frame + '"/>');
          o.push('<path class="cab-u" opacity="' + ro + '" d="' + urows + '"/>');
        }

        /* 机柜指示灯：多色、硬切闪烁，像真实机房的硬盘灯 */
        if (s.h >= 230) {
          var ledD = dF + (dB - dF) * 0.09;
          var ln = 3 + Math.round(rnd(j * 7 + k * 5 + si) * 2);
          for (c = 0; c < ln; c++) {
            var lx = rx0 + (rx1 - rx0) * (0.10 + 0.80 * rnd(j * 31 + k * 17 + c * 9 + si));
            var ly = yF(lx, j, ledD) - cabH0 * (0.16 + 0.58 * rnd(k * 13 + c * 7 + j));
            o.push('<circle class="cab-led-dot' + tintCls(j * 41 + k * 19 + c * 7, 0.78) +
                   '" cx="' + r2(ox(lx, ledD)) + '" cy="' + r2(ly) + '" r="0.95"' +
                   twk(j * 23 + k * 11 + c * 3, 1.0, 1.3) + '/>');
          }
        }

        /* 房间玻璃隔断（左边界竖线） */
        o.push('<path class="rm-part" d="M' + r2(ox(rx0, dF)) + ',' + r2(yF(rx0, j, dF)) +
               'L' + r2(ox(rx0, dF)) + ',' + r2(yC(rx0, j, dF)) + '"/>');

        /* 房间节点 */
        var ndX = ox(mx, dF + 4), ndY = yF(mx, j, dF + 4) - H(mx) * 0.52;
        nodes.push([ndX, ndY]);
      }

      /* 本层层板（带厚度的暗色梁） */
      var slabD = 'M' + r2(ox(s.x0, dF)) + ',' + r2(yF(s.x0, j, dF)) +
                  'L' + r2(ox(s.x1, dF)) + ',' + r2(yF(s.x1, j, dF));
      o.push('<path class="slab-band" d="' + slabD + '"/>');
      o.push('<path class="slab-edge" d="' + slabD + '"/>');

      /* 本层水平母线（实线 + 流动光带） */
      var busY = yC(s.x0 + 8, j, dF) + H(s.x0 + 8) * 0.30;
      var busD = 'M' + r2(ox(s.x0 + 3, dF)) + ',' + r2(busY) +
                 'L' + r2(ox(s.x1 - 3, dF)) + ',' + r2(yC(s.x1 - 3, j, dF) + H(s.x1 - 3) * 0.30);
      o.push('<path class="flr-bus" d="' + busD + '"/>');
      o.push('<path class="flr-flow" d="' + busD + '" style="animation-delay:' + r2(-j * 0.5) + 's"/>');

      /* 节点：挂在母线上 */
      var layerNodes = [];
      for (k = 0; k < nodes.length; k++) {
        var nx = nodes[k][0], ny = busY;
        layerNodes.push([nx, ny]);
        var tk = j * 9 + k * 5 + si * 3;
        var tc = tintCls(tk, 0.58);
        o.push('<path class="nd-tick" d="M' + r2(nx) + ',' + r2(ny) + 'L' + r2(nx) + ',' + r2(nodes[k][1] + 3) + '"/>');
        o.push('<circle class="nd-halo' + tc + '" cx="' + r2(nx) + '" cy="' + r2(ny) + '" r="' + r2(5.6 + rnd(j * 3 + k) * 2.4) +
               '" style="animation-delay:' + r2(-rnd(k * 5 + j) * 3.4) + 's;--td:' + r2(2.2 + rnd(k * 7 + j * 3) * 1.6) + 's"/>');
        /* 部分节点带十字星芒 */
        if (rnd(tk * 13 + 7) < 0.38) {
          var sa2 = r2(4.6 + rnd(tk * 3 + 1) * 2.6);
          o.push('<path class="nd-star' + tc + '" d="M' + r2(nx - sa2) + ',' + r2(ny) + 'L' + r2(nx + sa2) + ',' + r2(ny) +
                 'M' + r2(nx) + ',' + r2(ny - sa2) + 'L' + r2(nx) + ',' + r2(ny + sa2) + '"' + twk(tk + 21, 1.7, 2.0) + '/>');
        }
        o.push('<circle class="nd' + tc + '" cx="' + r2(nx) + '" cy="' + r2(ny) + '" r="' + r2(2.4 + rnd(j * 7 + k * 3 + si) * 1.3) +
               '"' + twk(tk + 4, 1.9, 2.2) + '/>');
      }

      allNodes.push(layerNodes);

      /* 层间竖向干线 */
      var tx = s.x0 + (s.x1 - s.x0) * 0.5;
      if (j < N - 1) {
        o.push('<path class="trunk" d="M' + r2(ox(tx, dM)) + ',' + r2(yC(tx, j, dM)) +
               'L' + r2(ox(tx, dM)) + ',' + r2(yC(tx, j + 1, dM)) + '"/>');
      }

      o.push('</g>');
    }

    /* 跨层斜向网络链路（形成拓扑网状感） */
    var mesh = '', mi, mk;
    for (mi = 0; mi + 1 < allNodes.length; mi++) {
      var lo = allNodes[mi], hi = allNodes[mi + 1], cnt = Math.min(lo.length, hi.length);
      for (mk = 0; mk < cnt; mk++) {
        mesh += 'M' + r2(lo[mk][0]) + ',' + r2(lo[mk][1]) + 'L' + r2(hi[mk][0]) + ',' + r2(hi[mk][1]);
        if (mk + 1 < cnt) {
          mesh += 'M' + r2(lo[mk][0]) + ',' + r2(lo[mk][1]) + 'L' + r2(hi[mk + 1][0]) + ',' + r2(hi[mk + 1][1]);
        }
      }
    }
    if (mesh) o.push('<path class="nd-mesh" d="' + mesh + '"/>');

    /* 主干数据总线：贯穿本体量全高 + 上行流动光点 */
    var tX = s.x0 + (s.x1 - s.x0) * 0.5;
    var tp = 'M' + r2(ox(tX, dM)) + ',' + r2(yF(tX, 0, dM) - 4) +
             'L' + r2(ox(tX, dM)) + ',' + r2(yC(tX, N - 1, dM) + 2);
    o.push('<path class="bus" d="' + tp + '"/>');
    o.push('<path class="bus-flow' + (s.side ? ' g2' : '') + '" d="' + tp + '"/>');

    o.push('</g>');
    return o.join('');
  }

  /* ---------- 主渲染 ---------- */
  function svg(opt) {
    opt = opt || {};
    var hl = opt.floor;
    var hlIdx = (hl && /^[1-9]/.test(hl)) ? Math.min(N, parseInt(hl, 10)) - 1 : -1;
    var gid = 'blW' + (opt.uid || 'A');

    var o = [], i, j;
    var alloc = null;
    if (opt.inner !== false) {
      var FL = (global.DCIM && global.DCIM.floors) || [];
      alloc = [];
      for (j = 0; j < N; j++) {
        var tot = 6, k2;
        for (k2 = 0; k2 < FL.length; k2++) if (FL[k2].floor === (j + 1) + 'F') tot = FL[k2].rooms || 6;
        alloc.push(allocRooms(Math.max(2, tot)));
      }
    }
    o.push('<svg class="bld-svg" width="100%" height="100%" viewBox="' + VB_X + ' ' + VB_Y + ' ' + VB_W + ' ' + VB_H + '" ' +
           'preserveAspectRatio="xMidYMid meet" xmlns="http://www.w3.org/2000/svg">');

    /* ---- defs ---- */
    var clipD = '';
    SEGS.forEach(function (s) { clipD += '<path d="' + facePath(s) + '"/><path d="' + roofPath(s) + '"/>'; });

    o.push('<defs>' +
      '<linearGradient id="' + gid + '" x1="0" y1="0" x2="0" y2="1">' +
        '<stop offset="0" stop-color="#EAFDFF" stop-opacity="0.05"/>' +
        '<stop offset="0.42" stop-color="#6FE2FB" stop-opacity="0.04"/>' +
        '<stop offset="1" stop-color="#1199D4" stop-opacity="0.02"/>' +
      '</linearGradient>' +
      '<linearGradient id="' + gid + 's" x1="0" y1="0" x2="0" y2="1">' +
        '<stop offset="0" stop-color="#2FC4F0" stop-opacity="0.12"/>' +
        '<stop offset="1" stop-color="#04182C" stop-opacity="0.48"/>' +
      '</linearGradient>' +
      '<linearGradient id="' + gid + 'f" x1="0" y1="0" x2="0" y2="1">' +
        '<stop offset="0" stop-color="#0B5AA0" stop-opacity="0.13"/>' +
        '<stop offset="1" stop-color="#020C1A" stop-opacity="0.42"/>' +
      '</linearGradient>' +
      '<linearGradient id="' + gid + 'rf" x1="0" y1="0" x2="1" y2="1">' +
        '<stop offset="0" stop-color="#0E5CA4" stop-opacity="0.26"/>' +
        '<stop offset="1" stop-color="#061E38" stop-opacity="0.42"/>' +
      '</linearGradient>' +
      '<radialGradient id="' + gid + 'halo" cx="0.5" cy="0.5" r="0.5">' +
        '<stop offset="0" stop-color="#12AEF8" stop-opacity="0.44"/>' +
        '<stop offset="0.45" stop-color="#0A6CC4" stop-opacity="0.22"/>' +
        '<stop offset="1" stop-color="#04101F" stop-opacity="0"/>' +
      '</radialGradient>' +
      '<linearGradient id="' + gid + 'beam" x1="0" y1="1" x2="0" y2="0">' +
        '<stop offset="0" stop-color="#C9FAFF" stop-opacity="0.85"/>' +
        '<stop offset="0.40" stop-color="#37D0F5" stop-opacity="0.40"/>' +
        '<stop offset="1" stop-color="#0A78C8" stop-opacity="0"/>' +
      '</linearGradient>' +
      '<linearGradient id="' + gid + 'sweep" x1="0" y1="0" x2="1" y2="0">' +
        '<stop offset="0" stop-color="#9DF2FF" stop-opacity="0"/>' +
        '<stop offset="0.5" stop-color="#BFF8FF" stop-opacity="0.15"/>' +
        '<stop offset="1" stop-color="#9DF2FF" stop-opacity="0"/>' +
      '</linearGradient>' +
      '<linearGradient id="' + gid + 'hscan" x1="0" y1="0" x2="0" y2="1">' +
        '<stop offset="0" stop-color="#7DF0FF" stop-opacity="0"/>' +
        '<stop offset="0.62" stop-color="#9DF6FF" stop-opacity="0.09"/>' +
        '<stop offset="1" stop-color="#D8FBFF" stop-opacity="0.24"/>' +
      '</linearGradient>' +
      '<linearGradient id="' + gid + 'base" x1="0" y1="0" x2="1" y2="0">' +
        '<stop offset="0" stop-color="#00D6FF" stop-opacity="0.10"/>' +
        '<stop offset="0.42" stop-color="#7DF0FF" stop-opacity="0.90"/>' +
        '<stop offset="1" stop-color="#00D6FF" stop-opacity="0.22"/>' +
      '</linearGradient>' +
      '<linearGradient id="' + gid + 'inF" x1="0" y1="0" x2="0" y2="1">' +
        '<stop offset="0" stop-color="#0E6FAE" stop-opacity="0.22"/>' +
        '<stop offset="1" stop-color="#031B34" stop-opacity="0.46"/>' +
      '</linearGradient>' +
      '<linearGradient id="' + gid + 'inB" x1="0" y1="0" x2="0" y2="1">' +
        '<stop offset="0" stop-color="#0A3A63" stop-opacity="0.52"/>' +
        '<stop offset="1" stop-color="#041A31" stop-opacity="0.66"/>' +
      '</linearGradient>' +
      '<radialGradient id="' + gid + 'nd" cx="0.5" cy="0.5" r="0.5">' +
        '<stop offset="0" stop-color="#E8FFFF" stop-opacity="1"/>' +
        '<stop offset="0.4" stop-color="#5FE6FF" stop-opacity="0.8"/>' +
        '<stop offset="1" stop-color="#00A8E8" stop-opacity="0"/>' +
      '</radialGradient>' +
      '<clipPath id="' + gid + 'clip">' + clipD + '</clipPath>' +
    '</defs>');

    /* ---- 背景光效 ---- */
    o.push(fxBack(gid));

    /* ---- 地面透视网格 ---- */
    o.push('<g class="bld-ground">');
    for (i = 0; i < 8; i++) {
      var y = gy(X0) + 10 + i * 12;
      var xa = X0 - 12 + i * 4;
      var xb = ox(X1, DEP) + 6 + i * 20;
      o.push('<line x1="' + r2(xa) + '" y1="' + r2(y) + '" x2="' + r2(xb) +
             '" y2="' + r2(y + (xb - xa) * 0.118) + '"/>');
    }
    for (i = 0; i <= 12; i++) {
      var t = i / 12, bx = X0 + (X1 - X0) * t;
      o.push('<line x1="' + r2(bx) + '" y1="' + r2(gy(bx)) + '" x2="' +
             r2(bx - 4 + t * 34) + '" y2="' + r2(gy(bx) + 96) + '"/>');
    }
    o.push('</g>');

    /* ---- 基座台阶 ---- */
    o.push('<g class="bld-base">');
    o.push('<path d="M' + r2(X0 - 5) + ',' + r2(gy(X0)) +
           'L' + r2(ox(X1, DEP) + 5) + ',' + r2(oy(gy(X1), DEP)) +
           'L' + r2(ox(X1, DEP) + 5) + ',' + r2(oy(gy(X1), DEP) + 10) +
           'L' + r2(X0 - 5) + ',' + r2(gy(X0) + 10) + 'Z"/>');
    o.push('</g>');

    /* ---- 楼体：按体量逐段绘制（自左向右，右侧遮住左侧） ---- */
    SEGS.forEach(function (s) {
      o.push('<g class="seg">');

      /* 前后错落形成的窄侧影 */
      if (s.dd) o.push('<path class="side" fill="url(#' + gid + 's)" d="' + sideFace(s, -s.dd) + '"/>');

      /* 最右端大侧立面 */
      if (s.side) {
        o.push('<path class="side" fill="url(#' + gid + 's)" d="' + sideFace(s, DEP) + '"/>');
        o.push(sideWins(s, DEP, 8, gid));
        for (var k = 1; k < N; k++) {
          var ty = topY(s, s.x1), by = gy(s.x1);
          o.push('<path class="slab" d="M' + r2(ox(s.x1, s.dd)) + ',' + r2(oy(by - k * (by - ty) / N, s.dd)) +
                 'L' + r2(ox(s.x1, s.dd + DEP)) + ',' + r2(oy(by - k * (by - ty) / N, s.dd + DEP)) + '"/>');
        }
      }

      /* 内部构造（先画，玻璃立面覆盖其上） */
      if (alloc) o.push(innerSeg(s, gid, hlIdx, alloc, SEGS.indexOf(s)));

      /* 正立面（半透明玻璃幕墙） */
      o.push('<path class="face" fill="url(#' + gid + 'f)" d="' + facePath(s) + '"/>');

      var cw = (s.x1 - s.x0) / s.cols;
      var pad = cw * 0.25;
      var wins = '', mull = '';

      for (j = 0; j < N; j++) {
        for (i = 0; i < s.cols; i++) {
          var cls = 'win' + (j === hlIdx ? ' is-hl' : '') +
                    (((i * 5 + j * 3 + s.x0) % 11 < 3) ? ' is-dim' : '');
          wins += '<path class="' + cls + '" fill="url(#' + gid + ')" data-f="' + (j + 1) + '" d="' +
                  winPath(s, s.x0 + i * cw + pad, s.x0 + (i + 1) * cw - pad, j) + '"/>';
        }
      }
      o.push(wins);

      /* 立面竖向分格（极淡） */
      for (i = 1; i < s.cols; i++) {
        var mx = s.x0 + i * cw;
        mull += '<path class="mull" d="M' + r2(ox(mx, s.dd)) + ',' + r2(oy(gy(mx), s.dd)) +
                'L' + r2(ox(mx, s.dd)) + ',' + r2(oy(topY(s, mx), s.dd)) + '"/>';
      }
      o.push(mull);

      for (var k2 = 1; k2 < N; k2++) o.push('<path class="slab" d="' + slabPath(s, k2) + '"/>');

      /* 层间墙裙线（每层窗台下方一道贯通的横向结构线） */
      var spand = '';
      for (var q = 0; q < N; q++) {
        function wy(x) { return oy(gy(x) - (q + 0.06) * s.h * pf(x) / N, s.dd); }
        spand += 'M' + r2(ox(s.x0, s.dd)) + ',' + r2(wy(s.x0)) +
                 'L' + r2(ox(s.x1, s.dd)) + ',' + r2(wy(s.x1));
      }
      o.push('<path class="spandrel" d="' + spand + '"/>');

      /* 竖向立管（每 3 列一条，自上而下贯通） */
      var riser = '';
      for (i = 3; i < s.cols; i += 3) {
        var rx2 = s.x0 + cw * i;
        riser += 'M' + r2(ox(rx2, s.dd)) + ',' + r2(oy(gy(rx2), s.dd)) +
                 'L' + r2(ox(rx2, s.dd)) + ',' + r2(oy(topY(s, rx2), s.dd));
      }
      if (riser) o.push('<path class="riser" d="' + riser + '"/>');

      /* 挑出的竖向遮阳翼（高层体量：每 3 列一片） */
      if (s.h >= 230) {
        var fin = '', FOUT = 7;
        for (i = 2; i < s.cols; i += 3) {
          var fx = s.x0 + cw * i;
          var fty = topY(s, fx), fby = gy(fx);
          fin += 'M' + r2(ox(fx, s.dd)) + ',' + r2(oy(fby, s.dd)) +
                 'L' + r2(ox(fx, s.dd)) + ',' + r2(oy(fty, s.dd));
          /* 顶部挑出段 + 底部挑出段 */
          fin += 'M' + r2(ox(fx, s.dd)) + ',' + r2(oy(fty, s.dd)) +
                 'L' + r2(ox(fx - FOUT, s.dd)) + ',' + r2(oy(fty, s.dd) + KY * FOUT);
          fin += 'M' + r2(ox(fx, s.dd)) + ',' + r2(oy(fty, s.dd) + 9) +
                 'L' + r2(ox(fx - FOUT, s.dd)) + ',' + r2(oy(fty, s.dd) + 9 + KY * FOUT);
        }
        if (fin) o.push('<path class="fin" d="' + fin + '"/>');
      }

      /* 檐口下方的设备百叶带 */
      var louv = '', lx2;
      for (lx2 = s.x0 + 4; lx2 < s.x1 - 3; lx2 += 6.5) {
        louv += 'M' + r2(ox(lx2, s.dd)) + ',' + r2(oy(topY(s, lx2) + 7, s.dd)) +
                'L' + r2(ox(lx2, s.dd)) + ',' + r2(oy(topY(s, lx2) + 16, s.dd));
      }
      louv += 'M' + r2(ox(s.x0 + 2, s.dd)) + ',' + r2(oy(topY(s, s.x0 + 2) + 6, s.dd)) +
              'L' + r2(ox(s.x1 - 2, s.dd)) + ',' + r2(oy(topY(s, s.x1 - 2) + 6, s.dd));
      louv += 'M' + r2(ox(s.x0 + 2, s.dd)) + ',' + r2(oy(topY(s, s.x0 + 2) + 17, s.dd)) +
              'L' + r2(ox(s.x1 - 2, s.dd)) + ',' + r2(oy(topY(s, s.x1 - 2) + 17, s.dd));
      o.push('<path class="louvre" d="' + louv + '"/>');

      /* 体量转角竖线（暗缝 + 结构立柱亮线） */
      o.push('<path class="seam" d="M' + r2(ox(s.x0, s.dd)) + ',' + r2(oy(gy(s.x0), s.dd)) +
             'L' + r2(ox(s.x0, s.dd)) + ',' + r2(oy(topY(s, s.x0), s.dd)) + '"/>');
      o.push('<path class="post" d="M' + r2(ox(s.x0, s.dd) - 1.6) + ',' + r2(oy(gy(s.x0), s.dd)) +
             'L' + r2(ox(s.x0, s.dd) - 1.6) + ',' + r2(oy(topY(s, s.x0), s.dd)) + '"/>');
      if (!s.side) {
        o.push('<path class="post" d="M' + r2(ox(s.x1, s.dd) + 1.4) + ',' + r2(oy(gy(s.x1), s.dd)) +
               'L' + r2(ox(s.x1, s.dd) + 1.4) + ',' + r2(oy(topY(s, s.x1), s.dd)) + '"/>');
      }

      /* 檐口挑板（顶板向外挑出，形成层叠托盘的阴影厚度） */
      var EAVE = 11, ed0 = s.dd - EAVE;
      o.push('<path class="eave-under" d="M' + r2(ox(s.x0, ed0)) + ',' + r2(oy(topY(s, s.x0), ed0)) +
             'L' + r2(ox(s.x1, ed0)) + ',' + r2(oy(topY(s, s.x1), ed0)) +
             'L' + r2(ox(s.x1, s.dd)) + ',' + r2(oy(topY(s, s.x1), s.dd)) +
             'L' + r2(ox(s.x0, s.dd)) + ',' + r2(oy(topY(s, s.x0), s.dd)) + 'Z"/>');
      o.push('<path class="eave-edge" d="M' + r2(ox(s.x0, ed0)) + ',' + r2(oy(topY(s, s.x0), ed0)) +
             'L' + r2(ox(s.x1, ed0)) + ',' + r2(oy(topY(s, s.x1), ed0)) + '"/>');
      o.push('<path class="eave-line" d="M' + r2(ox(s.x0, ed0)) + ',' + r2(oy(topY(s, s.x0), ed0)) +
             'L' + r2(ox(s.x0, s.dd)) + ',' + r2(oy(topY(s, s.x0), s.dd)) +
             'M' + r2(ox(s.x1, ed0)) + ',' + r2(oy(topY(s, s.x1), ed0)) +
             'L' + r2(ox(s.x1, s.dd)) + ',' + r2(oy(topY(s, s.x1), s.dd)) + '"/>');

      /* 屋面（紧接着本段立面绘制，保证被右侧体量正确遮挡） */
      o.push('<path class="roof" fill="url(#' + gid + 'rf)" d="' + roofPath(s) + '"/>');
      o.push('<path class="ridge" d="M' + r2(ox(s.x0, s.dd + ROOF)) + ',' + r2(oy(topY(s, s.x0), s.dd + ROOF)) +
             'L' + r2(ox(s.x1, s.dd + ROOF)) + ',' + r2(oy(topY(s, s.x1), s.dd + ROOF)) + '"/>');
      o.push('<path class="ridge front" d="M' + r2(ox(s.x0, s.dd)) + ',' + r2(oy(topY(s, s.x0), s.dd)) +
             'L' + r2(ox(s.x1, s.dd)) + ',' + r2(oy(topY(s, s.x1), s.dd)) + '"/>');
      o.push('<path class="ridge side" d="M' + r2(ox(s.x0, s.dd)) + ',' + r2(oy(topY(s, s.x0), s.dd)) +
             'L' + r2(ox(s.x0, s.dd + ROOF)) + ',' + r2(oy(topY(s, s.x0), s.dd + ROOF)) + '"/>');
      o.push('<path class="ridge side" d="M' + r2(ox(s.x1, s.dd)) + ',' + r2(oy(topY(s, s.x1), s.dd)) +
             'L' + r2(ox(s.x1, s.dd + ROOF)) + ',' + r2(oy(topY(s, s.x1), s.dd + ROOF)) + '"/>');

      o.push(roofKit(s, gid));

      o.push('</g>');
    });

    /* ---- 前台光效 ---- */
    o.push(fxFront(gid));

    /* ---- 楼层交互热区 ---- */
    if (opt.interactive) {
      o.push('<g class="bld-hit">');
      for (var j2 = N - 1; j2 >= 0; j2--) {
        SEGS.forEach(function (s) {
          var d = s.dd;
          function yt(x) { return oy(gy(x) - (j2 + 1) * s.h * pf(x) / N, d); }
          function yb(x) { return oy(gy(x) - j2 * s.h * pf(x) / N, d); }
          var p = 'M' + r2(ox(s.x0, d)) + ',' + r2(yt(s.x0)) +
                  'L' + r2(ox(s.x1, d)) + ',' + r2(yt(s.x1)) +
                  'L' + r2(ox(s.x1, d)) + ',' + r2(yb(s.x1)) +
                  'L' + r2(ox(s.x0, d)) + ',' + r2(yb(s.x0)) + 'Z';
          o.push('<path class="hit' + (j2 === hlIdx ? ' is-on' : '') + '" data-floor="' + (j2 + 1) + 'F" d="' + p + '"/>');
        });
      }
      o.push('</g>');
    }

    o.push('</svg>');
    return o.join('');
  }

  global.BLD = { svg: svg, VB_W: VB_W, VB_H: VB_H, floors: N };
})(window);
