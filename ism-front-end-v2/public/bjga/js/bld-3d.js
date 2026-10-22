/* =========================================================
   交互式 3D 仿真查看器（Canvas 2D 软件渲染 · 零依赖）
   楼层 / 机房 / 机柜 三级共用：
     · 轨道相机 —— 左键拖拽旋转 / 滚轮缩放 / 右键或 Shift 平移
     · WASD + QE 漫游 / R 复位 / F 聚焦 / 双击飞向对象
     · 点击房间·机柜直接下钻，镜头连续不跳
     · 空闲 4 秒自动环绕（游戏展示模式），一碰即停
     · 冷/热通道气流粒子、设备运行灯呼吸、满载告警脉冲
   约束：不用 SVG filter / SMIL / CSS animation（放大浏览会出黑块）
   ========================================================= */
(function (A) {
  'use strict';

  /* ==================== 基础工具 ==================== */
  function clamp(v, a, b) { return v < a ? a : v > b ? b : v; }
  function lerp(a, b, t) { return a + (b - a) * t; }

  var BG = [7, 17, 30];                 /* 雾色 = 舞台底色，远处自然融进背景 */
  var FONT = '-apple-system,BlinkMacSystemFont,"PingFang SC","Sarasa Gothic SC","Microsoft YaHei",sans-serif';
  var MONO = '"SF Mono",Menlo,Consolas,"Sarasa Mono SC",monospace';

  /* ---- 颜色：亮度缩放 + 雾化，全部带缓存（每帧上千次调用） ---- */
  var _cCache = {};
  /* 解析颜色：支持 '#RGB'/'#RRGGBB' 和 shade 产出的 'rgb(r,g,b)' 字符串。
     后者若走 hex 分支会被解析成 rgb(0,0,22) 一类的腐烂色（'rg'→0、'16'→0x16），
     面板会被自身的贴花盖成近黑 */
  function parseCol(hex) {
    var s = String(hex);
    if (s.indexOf('rgb') === 0) {
      var m = s.match(/-?\d+(\.\d+)?/g) || [];
      return [+m[0] || 0, +m[1] || 0, +m[2] || 0];
    }
    var h = s.replace('#', '');
    if (h.length === 3) h = h[0] + h[0] + h[1] + h[1] + h[2] + h[2];
    return [parseInt(h.slice(0, 2), 16) || 0, parseInt(h.slice(2, 4), 16) || 0, parseInt(h.slice(4, 6), 16) || 0];
  }
  function shade(hex, k) {
    var key = hex + '|' + k;
    var c = _cCache[key];
    if (c) return c;
    var p = parseCol(hex);
    var r = clamp(Math.round(p[0] * k), 0, 255),
      g = clamp(Math.round(p[1] * k), 0, 255),
      b = clamp(Math.round(p[2] * k), 0, 255);
    c = 'rgb(' + r + ',' + g + ',' + b + ')';
    if (_cCache.n > 6000) { _cCache = { n: 0 }; }
    _cCache[key] = c; _cCache.n = (_cCache.n || 0) + 1;
    return c;
  }
  /* t: 0=原色 1=完全雾化 */
  var _fCache = {};
  function fog(hex, k, t) {
    if (t <= 0.001) return shade(hex, k);
    var q = Math.round(clamp(t, 0, 1) * 24) / 24;
    var key = hex + '|' + k + '|' + q;
    var c = _fCache[key];
    if (c) return c;
    var p = parseCol(hex);
    var r = clamp(Math.round(p[0] * k), 0, 255), g = clamp(Math.round(p[1] * k), 0, 255), b = clamp(Math.round(p[2] * k), 0, 255);
    r = Math.round(lerp(r, BG[0], q)); g = Math.round(lerp(g, BG[1], q)); b = Math.round(lerp(b, BG[2], q));
    c = 'rgb(' + r + ',' + g + ',' + b + ')';
    if (_fCache.n > 6000) { _fCache = { n: 0 }; }
    _fCache[key] = c; _fCache.n = (_fCache.n || 0) + 1;
    return c;
  }

  /* ==================== 配色（与 2D / 2.5D 同一套口径） ==================== */
  var C = {
    floor: '#0B1A2A', floorEdge: '#1B3350', roomMain: '#0F2E44', roomEmpty: '#16202E',
    wall: '#24405C', wallTop: '#2F5478',
    cabOff: '#3A4C60', cabBlue: '#12A6D8', cabYellow: '#E89A2E', cabRed: '#E0524F',
    cold: '#1C68A8', hot: '#A8501F', lhd: '#C8A24D', ac: '#12A88C', ups: '#8A6BD1',
    /* 类型饰条色（设备顶部色带，与图例同色） */
    srv: '#2C7CB5', sto: '#7A5CC6', sw: '#12A08C', sec: '#D2703C', oth: '#5A6B82',
    /* 设备前面板：亮银灰，贴近真实服务器；类型辨识靠顶部饰条 */
    psrv: '#A9B1BB', psto: '#98A1AC', psw: '#B3BAC3', psec: '#A3A7AE', poth: '#93A0AB',
    /* 机柜结构：深灰（真实机柜色），与亮银灰设备形成清晰对比 */
    slot: '#1B2530', rail: '#232B34', top: '#2F3843', blank: '#1B232C',
    led: '#25F0A8', warn: '#FF5A6A'
  };

  function lvColor(c) {
    var lv = c.lv || (!c.dev ? 'off' : c.comp > 90 ? 'red' : c.comp >= 70 ? 'yellow' : 'blue');
    return lv === 'off' ? C.cabOff : lv === 'red' ? C.cabRed : lv === 'yellow' ? C.cabYellow : C.cabBlue;
  }
  /* 面板基色：亮银灰（贴近真实服务器） */
  function catColor(cat) {
    if (cat === '存储设备' || cat === '磁盘阵列') return C.psto;
    if (cat === '交换机') return C.psw;
    if (cat === '网络安全设备' || cat === '网闸') return C.psec;
    if (cat === '其他设备' || cat === '其他') return C.poth;
    return C.psrv;
  }
  /* 类型饰条色：设备顶部色带，与图例同色 */
  function catEdge(cat) {
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
  function capName(r) {
    if (r.type === '配电室') return '配电柜';
    if (r.type === '电池室') return '电池架';
    if (r.type === '空调间') return '精密空调';
    if (r.type === 'UPS间') return 'UPS 主机';
    return '机柜';
  }

  /* ==================== 图元 ==================== */
  var _n = [0, 0, 0];
  function normal(a, b, c) {
    var ux = b[0] - a[0], uy = b[1] - a[1], uz = b[2] - a[2];
    var vx = c[0] - a[0], vy = c[1] - a[1], vz = c[2] - a[2];
    var x = uy * vz - uz * vy, y = uz * vx - ux * vz, z = ux * vy - uy * vx;
    var l = Math.sqrt(x * x + y * y + z * z) || 1;
    _n[0] = x / l; _n[1] = y / l; _n[2] = z / l;
    return [_n[0], _n[1], _n[2]];
  }
  function centroid(p) {
    var x = 0, y = 0, z = 0, i;
    for (i = 0; i < p.length; i++) { x += p[i][0]; y += p[i][1]; z += p[i][2]; }
    return [x / p.length, y / p.length, z / p.length];
  }

  /* 面：p = 3~4 个顶点，n = 法线（朝外） */
  function F(p, fill, opt) {
    opt = opt || {};
    this.p = p;
    this.n = normal(p[0], p[1], p[2]);
    this.c = centroid(p);
    this.fill = fill;
    this.k = opt.k || 1;
    this.stroke = opt.stroke || null;
    this.sw = opt.sw || 1;
    this.op = opt.op == null ? 1 : opt.op;
    this.obj = opt.obj || null;
    this.blend = opt.blend || 0;      /* 半透明：第二趟画 */
    this.noFog = opt.noFog || 0;      /* 灯/粒子：不吃雾 */
    this.dyn = opt.dyn || null;       /* 动态面：fn(t) 返回 [r,g,b] 或 null */
    /* lay：渲染层。0 = 背景（地面/楼板），1 = 常规，2 = 覆盖（半透明罩）。
       地板是一整块几十米的四边形，画家算法用「四个角的平均深度」排序 ——
       近边 5m、远边 30m 平均出 20m，正好落在机柜阵列的深度区间里，于是
       地板被误判成"比机柜更近"而后画，把整个机柜阵列盖掉（用户看到的是
       一片黑色地板）。地面永远在所有物体之下，交给深度排序只会出错，
       所以显式分层：先画背景层，再按深度画常规层。 */
    this.lay = opt.lay || 0;
    /* fade：可淡出面（墙体）。渲染时若该面挡在相机与场景中心之间，
       自动把不透明度压到很低 —— 建筑可视化的标准做法：不拆墙，
       只是在挡住视线的那一瞬间让它透明。近侧墙不淡出的话，
       低角度看机柜会被齐腰高的墙切掉一截。 */
    this.fade = opt.fade ? 1 : 0;
    /* plate：贴在这个面上的文字（名称/编号）。跟面一起画、一起裁切，
       不进 pickLbl，点字等于点这块面板，不会另起一层抢走点击。 */
    this.plate = opt.plate || null;
    this.skin = opt.skin || null;
    this.cab = opt.cab || null;
    this.ax = opt.ax; this.ay = opt.ay; this.az = opt.az;
    this.frame = opt.frame || 0;
    /* quiet：细节面（网孔、色条、底座）仍可点中，但悬停时不各描一圈白边。 */
    this.quiet = opt.quiet ? 1 : 0;
    this.sv = [];                     /* 复用：屏幕点 */
  }
  /* 立方体六面：col 必须传原色，明暗由 k 分配 —— 预先 shade 会二次压暗成一团黑 */
  function box(out, x, y, z, w, d, h, col, opt) {
    opt = opt || {};
    var o = opt.o || {};
    var T = opt.topK == null ? 1.46 : opt.topK,
      R = opt.rightK == null ? 1.06 : opt.rightK,
      Fr = opt.frontK == null ? 0.90 : opt.frontK,
      L = opt.leftK == null ? 0.64 : opt.leftK,
      Bk = opt.backK == null ? 0.56 : opt.backK,
      Bt = opt.bottomK == null ? 0.42 : opt.bottomK;
    function mk(p, k) { out.push(new F(p, col, { k: k, o: o, stroke: o.stroke, sw: o.sw, op: o.op, obj: o.obj, lay: opt.lay == null ? 1 : opt.lay, fade: opt.fade, quiet: opt.quiet })); }
    mk([[x, y, z + h], [x + w, y, z + h], [x + w, y + d, z + h], [x, y + d, z + h]], T);   /* 顶 +z */
    mk([[x, y, z], [x, y + d, z], [x + w, y + d, z], [x + w, y, z]], Bt);               /* 底 -z */
    if (!opt.skipRight) mk([[x + w, y, z], [x + w, y + d, z], [x + w, y + d, z + h], [x + w, y, z + h]], R); /* 右 +x */
    if (!opt.skipLeft) mk([[x, y, z], [x, y, z + h], [x, y + d, z + h], [x, y + d, z]], L);               /* 左 -x */
    mk([[x, y + d, z], [x, y + d, z + h], [x + w, y + d, z + h], [x + w, y + d, z]], Fr); /* 前 +y */
    mk([[x, y, z], [x + w, y, z], [x + w, y, z + h], [x, y, z + h]], Bk);               /* 后 -y */
    return out;
  }
  /* 水平面（地板/楼板/通道色带）。lay 默认 0 = 背景层：
     这些面都贴着 z≈0 的地面，物理上永远在所有柜体之下，不该参与深度竞争 */
  function plate(out, x, y, z, w, d, col, opt) {
    opt = opt || {};
    var o = opt.o || {};
    out.push(new F([[x, y, z], [x + w, y, z], [x + w, y + d, z], [x, y + d, z]], col, {
      k: opt.k == null ? 1 : opt.k, o: o, stroke: o.stroke, sw: o.sw, op: o.op, obj: o.obj,
      blend: o.blend, noFog: o.noFog, lay: opt.lay == null ? 0 : opt.lay, quiet: opt.quiet
    }));
    return out;
  }
  /* 竖直面：贴在 y = y 平面上（设备面板） */
  function face(out, y, x, z, w, h, col, opt) {
    opt = opt || {};
    var o = opt.o || {};
    out.push(new F([[x, y, z], [x, y, z + h], [x + w, y, z + h], [x + w, y, z]], col, {
      k: opt.k == null ? 1 : opt.k, o: o, stroke: o.stroke, sw: o.sw, op: o.op, obj: o.obj,
      blend: o.blend, noFog: o.noFog, dyn: o.dyn, lay: 1
    }));
    return out;
  }
  /* 竖直面：贴在 x = x 平面上（机柜侧板细节） */
  function side(out, x, y, z, d, h, col, opt) {
    opt = opt || {};
    var o = opt.o || {};
    out.push(new F([[x, y, z], [x, y + d, z], [x, y + d, z + h], [x, y, z + h]], col, {
      k: opt.k == null ? 1 : opt.k, o: o, stroke: o.stroke, sw: o.sw, op: o.op, obj: o.obj,
      blend: o.blend, noFog: o.noFog, dyn: o.dyn, lay: 1
    }));
    return out;
  }
  /* 线：带法线，贴背面的线不画（否则柜内线条乱飞） */
  function L(out, a, b, color, opt) {
    opt = opt || {};
    var n = normal(a, b, [a[0], a[1], a[2] + 1]);
    if (opt.flat) n = opt.flat;
    out.push({ t: 1, a: a, b: b, n: n, c: color, w: opt.w || 1, op: opt.op == null ? 1 : opt.op,
      dash: opt.dash || null, obj: opt.obj || null, blend: opt.blend || 0 });
    return out;
  }
  /* 标签：屏幕空间文字。
     obj：可选 —— 带 obj 的标签参与拾取（点文字 = 选它代表的对象）。
           之前标签只进渲染队列、不进拾取表，射线直接穿过字打到背后的地板，
           于是"点 209 空调却选中 207 空调间"。
     lay：渲染层。1 = 与常规几何一起按深度排序（场景内的刻度/地面文字）；
           2 = 覆盖层，永远压在最上面（可交互的房号/图注）。
           标签必须带 lay —— 第 29 轮加分层时漏了这个字段，标签掉进
           layer 0 被机柜整片盖住，用户根本看不见房号。
     dim：淡显。密集场景（110 台机柜各自带编号）里全量常显会糊成一片，
           标了 dim 的标签默认压到低透明度，hover / 选中对应的对象时才
           高亮显形 —— 既能认出台数，又不会淹掉画面。 */
  function T(out, x, y, z, text, opt) {
    opt = opt || {};
    out.push({ t: 2, p: [x, y, z], s: text, size: opt.size || 11, col: opt.col || '#9FB6CE',
      align: opt.align || 'middle', dy: opt.dy || 0, weight: opt.weight || '400',
      family: opt.family || FONT, bg: opt.bg || null, mono: opt.mono || 0, glow: opt.glow || 0,
      obj: opt.obj || null, lay: opt.lay == null ? 1 : opt.lay, dim: opt.dim || 0 });
    return out;
  }

  /* ==================== 场景构建 ==================== */
  var S = {};   /* S.floor / S.room / S.rack —— 返回 {faces,lines,labels,bb,cam,objs,parts} */

  /* ---------- L1 楼层 ---------- */
  S.floor = function (floor) {
    var D = window.DCIM;
    var rs = D.rooms.filter(function (r) { return r.floor === floor; });
    var faces = [], lines = [], labels = [];

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
    var objs = {}, parts = [];

    /* 楼板 + 地面网格 */
    plate(faces, -1.4, -1.4, -0.5, W + 2.8, H + 2.8, C.floor, { o: { stroke: C.floorEdge, sw: 1.4 } });
    for (var gx = 0; gx <= W + 2.8; gx += 3) {
      L(lines, [-1.4 + gx, -1.4, -0.49], [-1.4 + gx, H + 1.4, -0.49], '#132A42', { w: 1, op: 0.5 });
    }
    for (var gy = 0; gy <= H + 2.8; gy += 3) {
      L(lines, [-1.4, -1.4 + gy, -0.49], [W + 1.4, -1.4 + gy, -0.49], '#132A42', { w: 1, op: 0.5 });
    }

    /* 走廊 */
    plate(faces, -0.7, yC, -0.3, W + 1.4, COR, '#0E2436', { o: { op: 0.92, stroke: '#1D4A66', sw: 1 } });

    function drawRoom(it, y0) {
      var r = it.r, sz = it.sz, x = it.x, y = y0;
      var tc = TYPE_COLOR[r.type] || '#8FA8C4';
      var lv = !r.cabCount ? 'off' : r.comp > 90 ? 'red' : r.comp >= 70 ? 'yellow' : 'blue';
      var edge = lv === 'red' ? C.cabRed : lv === 'yellow' ? C.cabYellow : lv === 'blue' ? C.cabBlue : '#3E5470';
      var id = 'room:' + r.name;
      var oo = { obj: id, stroke: shade(edge, 0.72), sw: 1.2 };

      plate(faces, x, y, 0, sz.w, sz.d, r.cabCount ? C.roomMain : C.roomEmpty, { o: oo });
      objs[id] = {
        kind: 'room', name: r.name, drill: r.name,
        focus: [x + sz.w / 2, y + sz.d / 2, 1.2],
        hint: r.name + ' · ' + r.type + (r.cabCount ? ' · ' + r.cabCount + ' 柜 · 综合 ' + Math.round(r.comp) + '%' : ' · ' + r.dev + ' 台独立设备')
      };

      if (!r.cabCount) {
        /* 功能间：矮体块代表配电柜 / 电池架 / 空调 */
        var n = Math.max(2, Math.min(6, Math.round(sz.w / 2)));
        for (var q = 0; q < n; q++) {
          box(faces, x + 0.8 + q * (sz.w - 1.6) / n, y + sz.d * 0.45, 0, (sz.w - 1.6) / n - 0.4, sz.d * 0.34, 1.5,
            tc, { topK: 0.55, rightK: 0.78, frontK: 0.62 });
        }
        T(labels, x + sz.w / 2, y + sz.d * 0.62, 1.85, capName(r) + ' ×' + n,
          { size: 10, col: shade(tc, 1.5), obj: id, lay: 2 });
      } else {
        /* 主机房：按列画机柜长条 */
        var lay = r.layout, cols = sz.cols, per = sz.per;
        var cw = 0.66, gap = 0.05;
        var gx2 = (sz.w - cols * 2.7) / 2;
        for (var k = 0; k < cols; k++) {
          var list = r.cabinets.filter(function (c) { return c.col === k + 1; });
          var perN = list.length || Math.round(per);
          var sum = 0, cnt = 0;
          list.forEach(function (c) { if (c.dev) { sum += c.comp; cnt++; } });
          var comp = cnt ? sum / cnt : 0;
          var cc = !cnt ? C.cabOff : comp > 90 ? C.cabRed : comp >= 70 ? C.cabYellow : C.cabBlue;
          var bx = x + gx2 + k * 2.7 + 0.5;
          var by = y + (sz.d - perN * (cw + gap)) / 2;
          var cid = 'room:' + r.name;
          box(faces, bx, by, 0.02, 1.15, Math.max(0.6, perN * (cw + gap) - gap), 2.0, cc,
            { topK: 1.24, rightK: 0.66, frontK: 0.6, o: { obj: cid } });
          /* 顶面柜位分隔 */
          for (var j = 1; j < perN; j++) {
            L(lines, [bx, by + j * (cw + gap), 2.03], [bx + 1.15, by + j * (cw + gap), 2.03],
              shade(cc, 1.7), { w: 1, op: 0.55, flat: [0, 0, 1] });
          }
          /* 满载列：柜顶红边脉冲 */
          if (comp > 90) {
            faces.push(new F([[bx + 0.06, by + 0.06, 2.03], [bx + 1.09, by + 0.06, 2.03],
              [bx + 1.09, by + Math.max(0.6, perN * (cw + gap) - gap) - 0.06, 2.03],
              [bx + 0.06, by + Math.max(0.6, perN * (cw + gap) - gap) - 0.06, 2.03]],
              C.warn, { k: 1, op: 0.8, obj: cid, noFog: 1, dyn: pulse }));
          }
        }
        /* 冷热通道地面色带 */
        if (lay && lay.aisles && lay.aisles.length) {
          for (var ai = 0; ai < lay.aisles.length && ai < cols - 1; ai++) {
            var isCold = lay.aisles[ai].t === 'cold';
            var ax = x + gx2 + ai * 2.7 + 1.65;
            plate(faces, ax, y + 0.9, 0.03, 1.0, sz.d - 1.8, isCold ? C.cold : C.hot,
              { o: { op: 0.5, obj: id } });
          }
        }
        /* 房间告警 */
        if (r.alarmN > 300) {
          faces.push(new F([[x + 0.1, y + 0.1, 2.04], [x + sz.w - 0.1, y + 0.1, 2.04],
            [x + sz.w - 0.1, y + sz.d - 0.1, 2.04], [x + 0.1, y + sz.d - 0.1, 2.04]],
            C.warn, { k: 1, op: 0.75, obj: id, noFog: 1, dyn: pulseSlow }));
        }
      }

      /* 标签：房号 + 一行摘要。
         obj：标签参与拾取 —— 房号是楼层视图里唯一能看见机房的标识，
         不点它就只能靠"点地面色块"，而地面色块被机柜一压就选不准。
         lay:2：覆盖层，永远压在机柜之上，否则房号被自家机柜整片盖住。 */
      var no = roomNo(r.name) || r.name;
      var tn = r.type.replace('机房', '').replace('间', '');
      var lz = r.cabCount ? 3.4 : 2.7;
      T(labels, x + sz.w / 2, y + sz.d / 2, lz, r.cabCount ? String(no) : (no + ' ' + tn),
        { size: 14, col: '#E6F4FF', weight: '600', bg: 'rgba(3,12,24,.86)', obj: id, lay: 2 });
      if (r.cabCount) {
        T(labels, x + sz.w / 2, y + sz.d / 2, lz, r.type + ' · ' + r.cabCount + ' 柜 · 综合 ' +
          Math.round(r.comp) + '%' + (r.alarmN > 300 ? ' · ▲' + A.num(r.alarmN) : ''),
          { size: 10.5, col: '#8FB0CC', dy: 16, bg: 'rgba(3,12,24,.7)', obj: id, lay: 2 });
      }
    }

    north.items.forEach(function (it) { drawRoom(it, yN); });
    south.items.forEach(function (it) { drawRoom(it, yS); });

    T(labels, W / 2, yC + COR / 2, 0.08, '走 廊 · COLD AISLE', { size: 10, col: '#4E6B86', align: 'start' });
    T(labels, W + 1.6, yC + COR / 2, 0.08, floor + ' · ' + (D.floors.filter(function (f) { return f.floor === floor; })[0] || {}).cabCount + ' 柜',
      { size: 12, col: '#9FB6CE', align: 'start' });

    return {
      faces: faces, lines: lines, labels: labels,
      bb: { x0: -1.4, y0: -1.4, z0: -0.5, x1: W + 1.4, y1: H + 1.4, z1: 3.6 },
      cam: { yaw: 0.62, pitch: 0.62, dist: Math.max(W, H) * 1.2, tx: W / 2, ty: H / 2, tz: 0.6 },
      objs: objs, parts: parts, kind: 'floor'
    };
  };

  /* ---------- L2 机房 ---------- */
  S.room = function (name) {
    var D = window.DCIM;
    var r = D.rooms.filter(function (x) { return x.name === name; })[0];
    if (!r) return null;
    var faces = [], lines = [], labels = [], objs = {}, parts = [];

    if (!r.cabCount) {
      /* 无机柜功能间 */
      var tc = TYPE_COLOR[r.type] || '#8FA8C4';
      var W0 = 9, D0 = 7;
      plate(faces, 0, 0, -0.06, W0, D0, C.roomEmpty, { o: { stroke: shade(tc, 0.8), sw: 1.6 } });
      var n0 = Math.max(3, Math.min(7, Math.round(r.dev / 2) || 4));
      var cid0 = 'room:' + r.name;
      for (var q0 = 0; q0 < n0; q0++) {
        box(faces, 0.9 + q0 * ((W0 - 1.8) / n0), D0 * 0.28, 0, (W0 - 1.8) / n0 - 0.35, D0 * 0.42, 2.0,
          tc, { topK: 0.72, rightK: 0.8, frontK: 0.6, o: { obj: cid0 } });
      }
      T(labels, W0 / 2, D0 * 0.34, 2.4, capName(r) + ' ×' + n0, { size: 11, col: shade(tc, 1.6), obj: cid0, lay: 2 });
      T(labels, W0 / 2, D0 + 0.6, 0.1, r.name + ' · ' + r.type + ' · ' + r.dev + ' 台独立设备',
        { size: 13, col: '#E6F4FF', weight: '600', bg: 'rgba(3,12,24,.86)', obj: cid0, lay: 2 });
      objs[cid0] = { kind: 'room', name: r.name, drill: null, focus: [W0 / 2, D0 / 2, 1], hint: r.name + ' · ' + r.type };
      return {
        faces: faces, lines: lines, labels: labels,
        bb: { x0: 0, y0: 0, z0: -0.06, x1: W0, y1: D0, z1: 2.6 },
        cam: { yaw: 0.55, pitch: 0.5, dist: 15, tx: W0 / 2, ty: D0 / 2, tz: 0.8 },
        objs: objs, parts: parts, kind: 'room'
      };
    }

    var lay = r.layout || {};
    var CD = 1.15, CW = 0.70, GAP = 0.06, AW = 1.75, H = 2.0;
    /* 机房自身的对象：列头柜 / 列标 / 房号标签都指向它。
       没有这一项时标签的 obj 会是 undefined —— 标签既不亮也不响应，
       看起来就是"点了没反应" */
    var roomId = 'room:' + r.name;
    objs[roomId] = {
      kind: 'room', name: r.name, drill: null,
      focus: [0, 0, 1], hint: r.name + ' · ' + r.type + ' · ' + r.cabCount + ' 柜'
    };
    var colsMap = {};
    r.cabinets.forEach(function (c) { (colsMap[c.col] = colsMap[c.col] || []).push(c); });
    var colKeys = Object.keys(colsMap).map(Number).sort(function (a, b) { return a - b; });
    var maxPer = Math.max.apply(null, colKeys.map(function (k) { return colsMap[k].length; }));
    var laneD = maxPer * (CW + GAP);
    var W = colKeys.length * (CD + AW) - AW;
    var Dp = laneD + 3.2;

    /* 地板 + 网格 */
    plate(faces, -1.2, -1.8, -0.06, W + 2.4, Dp + 2.2, C.roomMain, { o: { stroke: '#1D4A66', sw: 1.4 } });
    for (var gx3 = -1.2; gx3 <= W + 1.2; gx3 += 2) {
      L(lines, [gx3, -1.8, -0.05], [gx3, Dp + 0.4, -0.05], '#14304A', { w: 1, op: 0.45, flat: [0, 0, 1] });
    }
    for (var gy3 = -1.8; gy3 <= Dp + 0.4; gy3 += 2) {
      L(lines, [-1.2, gy3, -0.05], [W + 1.2, gy3, -0.05], '#14304A', { w: 1, op: 0.45, flat: [0, 0, 1] });
    }

    /* 四墙：房间才围得起来。四面都标fade —— 哪一面挡在相机和房间
       之间就在渲染时自动淡出，不用拆墙也能看清室内柜列 */
    var WH = 3.3, wo = { op: 0.9 }, wopt = { topK: 1.3, rightK: 0.7, frontK: 0.5, o: wo, fade: 1 };
    box(faces, -1.2, -1.8, -0.06, W + 2.4, 0.22, WH, C.wall, wopt);
    box(faces, -1.2, Dp + 0.2, -0.06, W + 2.4, 0.22, WH, C.wall, wopt);
    box(faces, -1.2, -1.8, -0.06, 0.22, Dp + 2.2, WH, C.wall, wopt);
    box(faces, W + 1.0, -1.8, -0.06, 0.22, Dp + 2.2, WH, C.wall, wopt);

    /* 俯视机柜：深灰柜体。台账没写门型；本屏机柜按高开孔率网孔门（进柜门贴图同结构）。
       冷通道前门、热通道尾门都是黑框黑网，顶上留一条负载色。列头柜仍是实心门。
       doorSign>0 冷通道在 +x。 */
    function makeAisleMesh() {
      var cnv = document.createElement('canvas');
      var W = 160, Hn = 520;
      cnv.width = W; cnv.height = Hn;
      var g = cnv.getContext('2d');
      g.fillStyle = '#0C1016';
      g.fillRect(0, 0, W, Hn);
      g.fillStyle = '#3A424C';
      g.fillRect(16, 18, W - 32, Hn - 36);
      g.fillStyle = '#05070A';
      var gap = 16, hw = 11, hh = 8, row = 0, yy, xx;
      for (yy = 24; yy < Hn - 24; yy += gap) {
        var shift = (row++ % 2) ? gap * 0.5 : 0;
        for (xx = 20 + shift; xx < W - 22; xx += gap) g.fillRect(xx, yy, hw, hh);
      }
      g.fillStyle = '#10141A';
      g.fillRect(0, 0, 16, Hn);
      g.fillRect(W - 16, 0, 16, Hn);
      g.fillRect(0, 0, W, 18);
      g.fillRect(0, Hn - 18, W, 18);
      g.strokeStyle = '#5C6772';
      g.lineWidth = 2;
      g.strokeRect(8, 9, W - 16, Hn - 18);
      return cnv;
    }
    function aisleSkin(kind) {
      if (!aisleSkin.cache) aisleSkin.cache = { mesh: makeAisleMesh() };
      return aisleSkin.cache[kind] || null;
    }
    function aisleCold(i) {
      var right = i < colKeys.length - 1 ? ((lay.aisles || [])[i] || { t: i % 2 === 0 ? 'cold' : 'hot' }) : null;
      var left = i > 0 ? ((lay.aisles || [])[i - 1] || { t: (i - 1) % 2 === 0 ? 'cold' : 'hot' }) : null;
      if (right && right.t === 'cold') return 1;
      if (left && left.t === 'cold') return -1;
      return right ? 1 : -1;
    }
    function drawAisleRack(x, y, depth, width, h, tint, id, doorSign, solid) {
      var o = { obj: id };
      var mesh = !solid;
      box(faces, x - 0.018, y - 0.014, 0, depth + 0.036, width + 0.028, 0.05, '#121820', {
        topK: 1.05, rightK: 0.7, frontK: 0.6, leftK: 0.48, backK: 0.42, o: o, quiet: 1
      });
      box(faces, x, y, 0.045, depth, width, h - 0.045, '#2A3442', {
        topK: 1.7, rightK: 0.86, frontK: 0.78, leftK: 0.58, backK: 0.46, o: o,
        skipRight: mesh || doorSign > 0, skipLeft: mesh || doorSign < 0
      });
      plate(faces, x + 0.07, y + 0.055, h + 0.016, depth - 0.14, width - 0.11, tint, {
        k: 1.2, lay: 1, o: o
      });
      plate(faces, x + depth * 0.34, y + width * 0.34, h + 0.03, depth * 0.28, width * 0.28, '#0B1016', {
        k: 0.8, lay: 1, quiet: 1, o: o
      });
      var doorY = y + 0.035, doorD = width - 0.07, doorZ = 0.11, doorH = h - 0.2;
      function doorFace(px, yy, zz, dd, hh, col, k, sign, skin) {
        /* 下左、上左、上右、下右，和铭牌同一套贴图坐标；法线朝柜外。 */
        var p = sign > 0
          ? [[px, yy + dd, zz], [px, yy + dd, zz + hh], [px, yy, zz + hh], [px, yy, zz]]
          : [[px, yy, zz], [px, yy, zz + hh], [px, yy + dd, zz + hh], [px, yy + dd, zz]];
        faces.push(new F(p, col, { k: k, obj: id, lay: 1, quiet: 1, skin: skin || null }));
      }
      function oneDoor(sign, useMesh) {
        var dx = sign > 0 ? x + depth : x;
        var hx = sign > 0 ? dx + 0.016 : dx - 0.016;
        doorFace(dx, doorY, doorZ, doorD, doorH, useMesh ? '#10141A' : '#3A3428', 1.05, sign, useMesh ? 'mesh' : null);
        doorFace(hx, doorY + 0.012, doorZ + doorH - 0.055, doorD - 0.024, 0.04, tint, 1.2, sign, null);
        if (!useMesh) {
          doorFace(hx, doorY + 0.02, doorZ + doorH * 0.42, doorD - 0.04, 0.05, '#E4D2A2', 1.15, sign, null);
        }
      }
      oneDoor(doorSign, mesh);
      if (mesh) oneDoor(-doorSign, true);
    }

    /* 列 / 机柜 / 通道 */
    colKeys.forEach(function (key, i) {
      var list = colsMap[key].slice().sort(function (a, b) { return a.seq - b.seq; });
      var x = i * (CD + AW);
      var y0 = 1.0;
      var doorSign = aisleCold(i);

      list.forEach(function (c, j) {
        var y = y0 + j * (CW + GAP);
        drawAisleRack(x, y, CD, CW, H, lvColor(c), 'cab:' + c.name, doorSign, false);
        objs['cab:' + c.name] = {
          kind: 'cab', name: c.name, drill: c.name,
          focus: [x + CD / 2, y + CW / 2, 1],
          hint: c.name + ' · ' + c.dev + ' 台 · 综合 ' + Math.round(c.comp) + '%'
        };
        /* 机柜编号（列-序号）。此前机房 3D 里一个编号都没有：110 台柜子
           在画面上只有 30×90px 的窄条，用户既认不出哪台是哪台，也无从
           确认自己点中的对不对 —— "机柜选不中"其实是"看不出选的是哪台"。
           全部常显会糊成一片（22 柜 × 5 列），所以默认淡显，
           hover / 选中时高亮，视线一过去就知道是哪台。 */
        T(labels, x + CD / 2, y + CW / 2, H + 0.02, c.col + '-' + c.seq,
          { size: 9, col: '#7FA6C6', align: 'middle', lay: 2, obj: 'cab:' + c.name, dim: 1 });
      });

      /* 列头柜：同样是深灰柜体，金色只留在柜顶，用来和设备柜分开。 */
      drawAisleRack(x, y0 - 1.0, CD, 0.72, H, C.lhd, roomId, doorSign, true);
      /* 列标是纯说明，不参与拾取：它挂在列头柜上方，投影后可能压到机柜，
         点了会得到"机房"而不是机柜 */
      T(labels, x + CD / 2, y0 - 1.05, H + 0.6, 'L' + key,
        { size: 12, col: '#FFC978', weight: '600', bg: 'rgba(3,12,24,.8)', lay: 2 });

      /* 列间通道 */
      if (i < colKeys.length - 1) {
        var as = (lay.aisles || [])[i] || { t: i % 2 === 0 ? 'cold' : 'hot' };
        var isCold = as.t === 'cold';
        var ax = x + CD;
        var aisleId = 'aisle:' + key;
        plate(faces, ax, y0 - 0.2, 0.02, AW, laneD + 0.4, isCold ? C.cold : C.hot,
          { o: { op: 0.85, obj: aisleId } });
        /* 通道对象必须真的登记：地板面带着 obj:'aisle:N'，若 objs 里没有这一项，
           click() 拿不到 drill/hint，点了就是死对象（高亮变了但什么都不发生） */
        objs[aisleId] = {
          kind: 'aisle', name: (isCold ? '冷通道' : '热通道') + ' ' + key,
          drill: null, focus: [ax + AW / 2, y0 + laneD / 2, 0.1],
          hint: (isCold ? '冷通道' : '热通道') + ' · 列 ' + key + '—' + (key + 1) + ' 之间 · ' +
            (isCold ? '冷风送入机柜前门' : '热风上行回风')
        };
        /* 气流：冷通道向柜、热通道上行 —— 用粒子，3D 里真的在动 */
        parts.push({
          x0: ax + 0.1, x1: ax + AW - 0.1, y0: y0, y1: y0 + laneD, cold: isCold,
          n: 14, seed: key * 7 + (isCold ? 0 : 3)
        });
        /* 通道文字是纯说明，不参与拾取。
           它挂在通道尽头、贴地（z=0.12），投影后正好落在下一列机柜上：
           实测压住 91 个机柜像素，点那些位置得到"通道"而不是机柜 ——
           正是"点柜却选到旁边"的一个来源。通道本身仍可点（地面色带
           带着 obj），文字只负责说明，不再抢点击。 */
        T(labels, ax + AW / 2, y0 + laneD + 0.55, 0.12, isCold ? '冷通道 COLD' : '热通道 HOT',
          { size: 10, col: isCold ? '#7FE9FF' : '#FFB27F', bg: 'rgba(3,12,24,.72)', lay: 2 });
      }
    });

    /* 精密空调 */
    var acN = lay.acIn || 0;
    if (acN) {
      for (var q2 = 0; q2 < acN; q2++) {
        box(faces, -1.05, 1.0 + q2 * (laneD / Math.max(1, acN)), 0, 0.8, 1.2, 2.2, C.ac,
          { topK: 1.4, rightK: 0.8, frontK: 0.6 });
      }
      T(labels, -0.65, 1.0, 2.5, 'AC ×' + acN, { size: 11, col: '#5BE3C8', align: 'start' });
    } else {
      var acr = (lay.acRoom && lay.acRoom.length) ? lay.acRoom.join('、') : '楼层空调间';
      T(labels, W / 2, Dp + 0.9, 0.12, '↗ 冷源：' + acr + ' 送风 · 地板下送风 / 冷通道封闭',
        { size: 10.5, col: '#7E97B0', bg: 'rgba(3,12,24,.7)' });
    }
    /* 机房总标签原来挂在画面右侧（W + 1.3）并带 obj。投影后它正好落在
       最后一列机柜上：实测压住 477 个机柜像素，且这些像素 hit() 全部返回
       "机房"而不是机柜 —— 用户点柜却进了机房。信息量与 HUD 标题重复，
       索性搬进 HUD，画布上不再放这条总标签。 */

    /* 相机：机房是「多列纵深排布」（205/407 是 5 列 × 22 柜）。
       俯角 0.60 rad（约 34°）是扫 12 组参数后挑出来的：再低前后柜互相遮挡、
       整排压成一条线；再高柜列缩到画面中间一条带、四边大片留白。
       yaw 0.85 让通道斜着走，5 列在屏幕上横向铺开而不是叠在一起。 */
    return {
      faces: faces, lines: lines, labels: labels, shellSkin: aisleSkin,
      sub: colKeys.length + ' 列 × ' + maxPer + ' 柜 · ' + r.cabCount + ' 架',
      bb: { x0: -1.3, y0: -1.9, z0: -0.06, x1: W + 1.3, y1: Dp + 0.5, z1: 3.5 },
      /* 取景只算柜列本体（不含 3.3m 高的墙和远处地板）：
         墙一旦进框，视距会被撑到十几米，柜子缩成一条细线 */
      fitBB: { x0: -1.0, y0: -0.6, z0: -0.06, x1: W + 1.0, y1: Dp + 0.2, z1: H + 0.3 },
      cam: { yaw: 0.85, pitch: 0.60, dist: Math.max(W, Dp) * 1.25, tx: W / 2, ty: Dp / 2 - 0.3, tz: 0.9 },
      objs: objs, parts: parts, kind: 'room'
    };
  };

  /* 局部 +y = 前面板法线（朝通道），+x = 柜宽，+z = 上。
     right = (ny, -nx)，是旋转不是镜像。 */
  function posePoint(pose, lx, ly, lz) {
    return [pose.ox + lx * pose.rx + ly * pose.nx, pose.oy + lx * pose.ry + ly * pose.ny, lz];
  }
  function posedFace(out, ly, lx, lz, w, h, col, opt, pose) {
    opt = opt || {};
    var o = opt.o || {};
    if (w <= 0.001 || h <= 0.001) return;
    out.push(new F([
      posePoint(pose, lx, ly, lz),
      posePoint(pose, lx, ly, lz + h),
      posePoint(pose, lx + w, ly, lz + h),
      posePoint(pose, lx + w, ly, lz)
    ], col, {
      k: opt.k == null ? 1 : opt.k, obj: o.obj || null, op: o.op, lay: 1,
      noFog: o.noFog, dyn: o.dyn, plate: o.plate || null, skin: o.skin || null, stroke: o.stroke, sw: o.sw,
      frame: o.frame || 0,
      cab: pose && pose.cab && pose.cab.name, ax: pose && pose.cx, ay: pose && pose.cy, az: 1
    }));
  }
  function posedBox(out, lx, ly, lz, w, d, h, col, opt, pose) {
    opt = opt || {};
    var o = opt.o || {};
    var Tk = opt.topK == null ? 1.35 : opt.topK;
    var Rk = opt.rightK == null ? 0.72 : opt.rightK;
    var Fk = opt.frontK == null ? 0.9 : opt.frontK;
    var Lk = opt.leftK == null ? 0.5 : opt.leftK;
    var Bk = opt.backK == null ? 0.4 : opt.backK;
    function P(x, y, z) { return posePoint(pose, x, y, z); }
    function mk(pts, k) {
      out.push(new F(pts, col, {
        k: k, obj: o.obj || null, op: o.op, lay: 1,
        cab: pose && pose.cab && pose.cab.name, ax: pose && pose.cx, ay: pose && pose.cy, az: 1
      }));
    }
    mk([P(lx, ly, lz + h), P(lx + w, ly, lz + h), P(lx + w, ly + d, lz + h), P(lx, ly + d, lz + h)], Tk);
    mk([P(lx, ly, lz), P(lx, ly + d, lz), P(lx + w, ly + d, lz), P(lx + w, ly, lz)], 0.34);
    mk([P(lx + w, ly, lz), P(lx + w, ly + d, lz), P(lx + w, ly + d, lz + h), P(lx + w, ly, lz + h)], Rk);
    mk([P(lx, ly, lz), P(lx, ly, lz + h), P(lx, ly + d, lz + h), P(lx, ly + d, lz)], Lk);
    if (!opt.skipFront) mk([P(lx, ly + d, lz), P(lx, ly + d, lz + h), P(lx + w, ly + d, lz + h), P(lx + w, ly + d, lz)], Fk);
    mk([P(lx, ly, lz), P(lx + w, ly, lz), P(lx + w, ly, lz + h), P(lx, ly, lz + h)], Bk);
  }
  function devPlate(a) {
    var name = [a.brand, a.model].filter(Boolean).join(' ');
    if (!name) name = a.type || a.cat || '设备';
    var uBit = 'U' + a.uStart + (a.uEnd && a.uEnd !== a.uStart ? '-' + a.uEnd : '');
    var num = a.sn ? String(a.sn) : uBit;
    var uH = a.uH || 1;
    var lines = [
      { t: name, col: '#EAF6FF', weight: '600' },
      { t: a.sn ? (num + '  ' + uBit) : num, col: '#8FDCFF', weight: '500' }
    ];
    if (uH >= 2) {
      var extra = a.cat || a.type || '';
      if (a.status) extra = extra ? (extra + ' · ' + a.status) : a.status;
      if (extra) lines.push({ t: extra, col: '#A9C0D0', weight: '400' });
    }
    /* short：1U 太矮才允许省略。更高的面板先缩小字号，尽量把全名画完。 */
    return { lines: lines, one: name + '  ' + num, keep: 1, short: uH < 2 };
  }

  /* ---------- L3：人在通道里走。巡视只画外壳；双击进入后才画这一台的内部面板 ---------- */
  S.rack = function (cabName, opts) {
    opts = opts || {};
    var inside = !!opts.inside;
    var D = window.DCIM;
    var cab0 = null, room = null;
    D.rooms.forEach(function (x) {
      x.cabinets.forEach(function (c) { if (c.name === cabName) { cab0 = c; room = x; } });
    });
    if (!cab0 || !room) return null;

    var depth = 1.05, width = 0.64, gap = 0.08, aisle = 1.70, H = 2.0, HU = H / 42;
    var colsMap = {};
    room.cabinets.forEach(function (c) {
      var col = c.col;
      if (col == null || col === '') {
        var p = String(c.name).split('-');
        col = Number(p[p.length - 2]);
        if (!col) col = 1;
      }
      (colsMap[col] = colsMap[col] || []).push(c);
    });
    var keys = Object.keys(colsMap).map(Number).sort(function (a, b) { return a - b; });
    if (!keys.length) return null;
    var maxPer = 1;
    keys.forEach(function (k) {
      colsMap[k].sort(function (a, b) {
        var sa = a.seq != null ? Number(a.seq) : Number(String(a.name).split('-').pop());
        var sb = b.seq != null ? Number(b.seq) : Number(String(b.name).split('-').pop());
        if (isNaN(sa)) sa = 0;
        if (isNaN(sb)) sb = 0;
        return sa - sb;
      });
      if (colsMap[k].length > maxPer) maxPer = colsMap[k].length;
    });

    var byCab = {};
    D.assets.forEach(function (a) {
      if (!a.cabinet || a.uStart > 42) return;
      (byCab[a.cabinet] = byCab[a.cabinet] || []).push(a);
    });

    var y0 = 0.45;
    var poses = [];
    keys.forEach(function (key, i) {
      var face = (i % 2 === 0) ? 1 : -1;
      var gx = i * (depth + aisle);
      var nx = face, ny = 0, rx = ny, ry = -nx;
      colsMap[key].forEach(function (c, j) {
        var gy = y0 + j * (width + gap);
        var ox = face > 0 ? gx : gx + depth;
        var oy = face > 0 ? gy + width : gy;
        var seq = c.seq != null ? c.seq : (Number(String(c.name).split('-').pop()) || (j + 1));
        poses.push({
          cab: c, key: key, seq: seq,
          gx: gx, gy: gy,
          nx: nx, ny: ny, rx: rx, ry: ry,
          ox: ox, oy: oy,
          fx: face > 0 ? gx + depth : gx,
          fy: gy + width / 2,
          cx: gx + depth / 2,
          cy: gy + width / 2,
          halfW: width / 2
        });
      });
    });

    var focus = poses[0];
    for (var fi = 0; fi < poses.length; fi++) {
      if (poses[fi].cab.name === cabName) focus = poses[fi];
    }
    var Wroom = keys.length * (depth + aisle) - aisle;
    var Dp = y0 + maxPer * (width + gap) + 0.4;

    var faces = [], lines = [], labels = [], objs = {};
    var blocks = [];
    plate(faces, -1.4, -0.8, -0.04, Wroom + 3.2, Dp + 1.6, '#101820', { o: { stroke: '#1A3348', sw: 1 } });
    function downFace(x, y, z, w, d, col, k) {
      faces.push(new F([[x, y, z], [x, y + d, z], [x + w, y + d, z], [x + w, y, z]], col, { k: k || 0.75, lay: 1 }));
    }
    downFace(-1.2, -0.5, 2.78, Wroom + 2.6, Dp + 1.2, '#121820', 0.55);
    keys.forEach(function (key, i) {
      if (i >= keys.length - 1) return;
      var cold = i % 2 === 0;
      var xA = i * (depth + aisle) + depth;
      var len = maxPer * (width + gap) + 0.2;
      plate(faces, xA + 0.06, y0 - 0.08, 0.008, aisle - 0.12, len, cold ? '#16324C' : '#3A2418', { o: { op: 0.55 } });
      plate(faces, xA + aisle * 0.42, y0, 0.02, aisle * 0.16, len - 0.2, cold ? '#1E6A9A' : '#6A3818', { o: { op: 0.35 } });
      if (!cold) return;
      var stepL = 1.7, yy, nLit = Math.ceil(len / stepL);
      for (yy = 0; yy < nLit; yy++) {
        downFace(xA + aisle * 0.28, y0 + 0.35 + yy * stepL, 2.66, aisle * 0.44, 0.38, '#E4F0FA', 1.65);
      }
    });

    /* 高端 42U 外壳：参照 NetShelter SX / Rittal TS IT 这类机柜的共同结构
       （单扇高开孔率前门、左铰链右旋锁、底座调平脚、顶盖走线口、开门后的 19 英寸立柱）。
       不贴厂商商标。网孔和 U 位刻度是各画一次的贴图，再贴到面上。 */
    function shellSkin(kind) {
      if (!shellSkin.cache) {
        shellSkin.cache = { door: makeDoorSkin(), rail: makeRailSkin() };
      }
      return shellSkin.cache[kind] || null;
    }
    function makeDoorSkin() {
      var cnv = document.createElement('canvas');
      cnv.width = 420;
      cnv.height = 1280;
      var g = cnv.getContext('2d');
      g.fillStyle = '#14181E';
      g.fillRect(0, 0, 420, 1280);
      g.fillStyle = '#2C343E';
      g.fillRect(0, 0, 34, 1280);
      g.fillRect(386, 0, 34, 1280);
      g.fillRect(0, 0, 420, 42);
      g.fillRect(0, 1238, 420, 42);
      var mx = 48, my = 58, mw = 324, mh = 1164;
      var grd = g.createLinearGradient(mx, 0, mx + mw, 0);
      grd.addColorStop(0, '#12181E');
      grd.addColorStop(0.38, '#5A6672');
      grd.addColorStop(0.7, '#2A323C');
      grd.addColorStop(1, '#161C22');
      g.fillStyle = grd;
      g.fillRect(mx, my, mw, mh);
      g.fillStyle = '#030508';
      var gapX = 22, gapY = 16, hw = 13, hh = 8, y, x, row = 0;
      for (y = my + 8; y < my + mh - 8; y += gapY) {
        var shift = (row++ % 2) ? gapX * 0.5 : 0;
        for (x = mx + 6 + shift; x < mx + mw - 14; x += gapX) g.fillRect(x, y, hw, hh);
      }
      g.strokeStyle = 'rgba(180,196,210,.45)';
      g.lineWidth = 2;
      g.strokeRect(mx + 1, my + 1, mw - 2, mh - 2);
      g.fillStyle = '#D5DEE8';
      var hy, hinges = [210, 620, 1040];
      for (hy = 0; hy < hinges.length; hy++) {
        g.fillRect(8, hinges[hy], 16, 52);
        g.fillStyle = '#6E7A88';
        g.fillRect(12, hinges[hy] + 10, 8, 32);
        g.fillStyle = '#D5DEE8';
      }
      g.fillStyle = '#E4EAF0';
      g.fillRect(392, 520, 14, 168);
      g.fillStyle = '#1A2028';
      g.fillRect(396, 568, 6, 78);
      g.beginPath();
      g.arc(399, 500, 6, 0, Math.PI * 2);
      g.fillStyle = '#F2F6FA';
      g.fill();
      g.fillStyle = '#8EA0B0';
      g.fillRect(168, 16, 84, 5);
      return cnv;
    }
    function makeRailSkin() {
      var U = 42, rowH = 36, W = 72, Hn = U * rowH;
      var cnv = document.createElement('canvas');
      cnv.width = W;
      cnv.height = Hn;
      var g = cnv.getContext('2d');
      g.fillStyle = '#C5CED6';
      g.fillRect(0, 0, W, Hn);
      g.fillStyle = '#8E99A4';
      g.fillRect(0, 0, 6, Hn);
      g.font = '600 18px sans-serif';
      g.textBaseline = 'middle';
      g.textAlign = 'right';
      var i, u, y;
      for (i = 0; i < U; i++) {
        u = U - i;
        y = i * rowH;
        g.fillStyle = (u % 2) ? '#B7C2CC' : '#D5DEE6';
        g.fillRect(6, y + 1, W - 8, rowH - 2);
        g.fillStyle = '#1A1E24';
        g.fillRect(W - 28, y + 8, 16, 18);
        g.fillStyle = '#1C2834';
        g.fillText(String(u), W - 34, y + rowH / 2);
      }
      return cnv;
    }
    function addPremiumShell(pose, id) {
      var CW = width, CD = depth, post = 0.032;
      /* 内衬只留侧壁和背板。正面若画成一整块，质心在视高附近，画家算法会把它排到
         低 U 设备前面，服务器面板就被这块近黑盖住。设备面和空位盲板已经封住开口。 */
      posedBox(faces, post, 0.04, 0.02, CW - post * 2, CD - 0.1, H - 0.04, '#12161C', {
        topK: 0.55, frontK: 0.38, rightK: 0.62, leftK: 0.4, skipFront: true, o: { obj: id }
      }, pose);
      posedBox(faces, 0, 0, 0, post, CD, H, '#2C343C', {
        topK: 1.4, rightK: 0.86, frontK: 1.05, o: { obj: id }
      }, pose);
      posedBox(faces, CW - post, 0, 0, post, CD, H, '#262E36', {
        topK: 1.32, leftK: 0.48, frontK: 0.9, o: { obj: id }
      }, pose);
      posedBox(faces, -0.016, -0.014, -0.072, CW + 0.032, CD + 0.022, 0.072, '#0E1218', {
        topK: 1.08, frontK: 0.72, o: { obj: id }
      }, pose);
      var foot = [[0.046, 0.05], [CW - 0.08, 0.05], [0.046, CD - 0.09], [CW - 0.08, CD - 0.09]];
      var fi;
      for (fi = 0; fi < foot.length; fi++) {
        posedBox(faces, foot[fi][0], foot[fi][1], -0.104, 0.034, 0.03, 0.034, '#D7E0E8', {
          topK: 1.65, frontK: 1.15
        }, pose);
      }
      posedBox(faces, -0.008, -0.008, H, CW + 0.016, CD + 0.014, 0.034, '#3A424C', {
        topK: 1.62, frontK: 0.95, o: { obj: id }
      }, pose);
      posedBox(faces, CW * 0.3, CD * 0.28, H + 0.03, CW * 0.4, CD * 0.3, 0.016, '#10161C', {
        topK: 0.75, o: { obj: id }
      }, pose);
      posedFace(faces, CD + 0.006, 0.004, 0.05, 0.006, H - 0.08, '#A8B4C0', {}, pose);
      posedFace(faces, CD + 0.006, CW - 0.01, 0.05, 0.006, H - 0.08, '#46525C', {}, pose);
    }
    function addCabFrame(pose) {
      posedFace(faces, depth + 0.018, -0.016, -0.04, width + 0.032, H + 0.1, '#FFFFFF', {
        o: { op: 0, frame: 1 }
      }, pose);
    }

    function drawDetail(pose) {
      var c = pose.cab;
      var idRack = 'rack:' + c.name;
      var CW = width, CD = depth;
      objs[idRack] = {
        kind: 'rack', name: c.name, drill: null,
        focus: posePoint(pose, CW / 2, CD, H * 0.5),
        hint: c.name + ' · 当前机柜'
      };
      addPremiumShell(pose, idRack);
      posedBox(faces, 0.04, 0.04, 0.02, CW - 0.08, 0.03, H - 0.04, '#0C1016', { frontK: 0.45 }, pose);
      /* 打开后不画门。19 英寸立柱留在门口两侧，U 位从下往上数，1 在最底。 */
      posedFace(faces, CD + 0.012, 0.028, 0.04, 0.03, H - 0.08, '#C5CED6', { o: { skin: 'rail' } }, pose);
      posedFace(faces, CD + 0.012, CW - 0.058, 0.04, 0.03, H - 0.08, '#C5CED6', { o: { skin: 'rail' } }, pose);

      var items = (byCab[c.name] || []).slice().sort(function (a, b) { return a.uStart - b.uStart; });
      var used = {};
      items.forEach(function (a) {
        var uH = a.uH || (a.uEnd ? a.uEnd - a.uStart + 1 : 1);
        for (var u = a.uStart; u < a.uStart + uH && u <= 42; u++) used[u] = 1;
      });
      var x0 = 0.058, fw = CW - 0.116;
      function blank(uStart, n) {
        var z0 = (uStart - 1) * HU + 0.0015, hh = n * HU - 0.003;
        posedFace(faces, CD + 0.004, x0, z0, fw, hh, '#121820', { o: { obj: idRack } }, pose);
        var ss = 0.007;
        posedFace(faces, CD + 0.01, x0 + 0.006, z0 + 0.004, ss, ss, '#8EA0B2', {}, pose);
        posedFace(faces, CD + 0.01, x0 + fw - ss - 0.006, z0 + 0.004, ss, ss, '#8EA0B2', {}, pose);
      }
      var run = 0, from = 1, uu;
      for (uu = 1; uu <= 42; uu++) {
        if (!used[uu]) { if (!run) from = uu; run++; continue; }
        if (run) { blank(from, run); run = 0; }
      }
      if (run) blank(from, run);

      items.forEach(function (a) {
        var uH = a.uH || (a.uEnd ? a.uEnd - a.uStart + 1 : 1);
        if (a.uStart > 42) return;
        var z0 = (a.uStart - 1) * HU + 0.0016;
        var hh = Math.max(HU * 0.82, uH * HU - 0.0032);
        var id = 'dev:' + a.sn;
        var yBack = CD - 0.50;
        posedBox(faces, x0, yBack, z0, fw, 0.50, hh, '#4A6278', {
          topK: 1.55, frontK: 0.78, rightK: 0.62, leftK: 0.48, o: { obj: id }
        }, pose);
        posedFace(faces, CD + 0.008, x0 + 0.006, z0 + 0.003, fw - 0.012, hh - 0.006, '#2C4158',
          { o: { obj: id } }, pose);
        var yDet = CD + 0.014;
        var ss = Math.min(0.009, hh * 0.2);
        posedFace(faces, yDet, x0 + 0.008, z0 + 0.005, ss, ss, '#D5E2EE', {}, pose);
        posedFace(faces, yDet, x0 + fw - ss - 0.008, z0 + 0.005, ss, ss, '#D5E2EE', {}, pose);
        posedFace(faces, yDet, x0 + 0.008, z0 + hh - ss - 0.006, ss, ss, '#D5E2EE', {}, pose);
        posedFace(faces, yDet, x0 + fw - ss - 0.008, z0 + hh - ss - 0.006, ss, ss, '#D5E2EE', {}, pose);
        if (hh > 0.02) {
          posedFace(faces, yDet, x0 + 0.006, z0 + hh - Math.min(0.008, hh * 0.16), fw - 0.012,
            Math.min(0.0055, hh * 0.1), catEdge(a.cat), {}, pose);
        }
        if (uH >= 2) {
          posedFace(faces, yDet, x0 + 0.02, z0 + hh * 0.22, 0.016, hh * 0.5, '#6A8094', {}, pose);
          var vs;
          for (vs = 0; vs < 3; vs++) {
            posedFace(faces, yDet, x0 + fw * 0.76, z0 + hh * 0.22 + vs * hh * 0.16, fw * 0.12, hh * 0.045, '#0A121A', {}, pose);
          }
        } else {
          posedFace(faces, yDet, x0 + 0.02, z0 + hh * 0.28, 0.012, hh * 0.4, '#6A8094', {}, pose);
        }
        posedFace(faces, yDet, x0 + fw - 0.02, z0 + hh * 0.38, 0.008, Math.max(0.006, Math.min(0.012, hh * 0.22)),
          C.led, { o: { noFog: 1, dyn: led } }, pose);
        var plate = devPlate(a);
        posedFace(faces, yDet, x0 + fw * 0.12, z0 + hh * 0.16, fw * 0.56, hh * 0.58, '#0C141E',
          { o: { obj: id, plate: plate } }, pose);
        var fc = posePoint(pose, x0 + fw / 2, CD + 0.02, z0 + hh / 2);
        objs[id] = {
          kind: 'dev', name: (a.brand || '') + ' ' + (a.model || ''), sn: a.sn,
          focus: fc,
          hint: (a.cat || '') + ' · ' + (a.brand || '') + ' ' + (a.model || '') + ' · ' + uH + 'U · U' +
            a.uStart + '-' + a.uEnd + (a.sn ? ' · ' + a.sn : '')
        };
      });
      addCabFrame(pose);
    }

    function drawSimple(pose) {
      var c = pose.cab;
      var cid = 'cab:' + c.name;
      var CW = width, CD = depth;
      var nDev = c.dev || (byCab[c.name] ? byCab[c.name].length : 0);
      objs[cid] = {
        kind: 'cab', name: c.name, drill: null,
        focus: posePoint(pose, CW / 2, CD, H * 0.5),
        hint: c.name + ' · ' + nDev + ' 台 · 综合 ' + Math.round(c.comp || 0) + '%'
      };
      addPremiumShell(pose, cid);
      /* 单扇网孔门：铰链在左、旋锁在右，贴图跟门面一起透视。 */
      var doorZ = 0.04, doorH = H - 0.1, doorX = 0.04, doorW = CW - 0.08;
      posedBox(faces, doorX, CD - 0.018, doorZ, doorW, 0.02, doorH, '#1A2028', {
        frontK: 0.8, topK: 1.1, skipFront: true, o: { obj: cid }
      }, pose);
      posedFace(faces, CD + 0.01, doorX + 0.002, doorZ + 0.012, doorW - 0.004, doorH - 0.02, '#14181E', {
        o: { obj: cid, skin: 'door' }
      }, pose);
      var items = byCab[c.name] || [];
      items.forEach(function (a) {
        if (a.uStart > 42) return;
        var uH = a.uH || (a.uEnd ? a.uEnd - a.uStart + 1 : 1);
        var z0 = (a.uStart - 1) * HU + 0.004;
        var hh = Math.max(0.012, uH * HU - 0.006);
        posedFace(faces, CD + 0.006, 0.008, z0, 0.012, hh, catEdge(a.cat), {}, pose);
      });
      addCabFrame(pose);
    }

    poses.forEach(function (pose) {
      /* 只挡柜体。打开的门不占通道，人可以左右、前后挪到正对设备。 */
      blocks.push({ x0: pose.gx, y0: pose.gy, x1: pose.gx + depth, y1: pose.gy + width });
      if (inside && pose.cab.name === cabName) drawDetail(pose);
      else drawSimple(pose);
    });

    /* 通道宽 1.7 米。站得太近（原先 0.82 米）一进画面就切掉柜顶和柜底。
       默认退到约 1.38 米，整柜落在画面里。 */
    var stand = 1.38;
    var px = focus.fx + focus.nx * stand;
    var py = focus.cy;
    var spawnYaw = Math.atan2(-focus.nx, -focus.ny || 0);
    /* 设备一般先占下半柜，视高从 1.62 降到 1.18，并略向下看。 */
    var eyeZ = 1.18;
    var pit = -0.16;
    var zoom0 = 0.76;
    var devN = (byCab[cabName] || []).length;
    return {
      faces: faces, lines: lines, labels: labels, shellSkin: shellSkin,
      sub: devN + ' 台 · 通道走动',
      bb: { x0: -1.2, y0: -0.6, z0: -0.05, x1: Wroom + 1.8, y1: Dp + 0.4, z1: H + 0.4 },
      fitBB: { x0: px - 0.4, y0: py - 0.4, z0: 0, x1: px + 0.4, y1: py + 0.4, z1: 2 },
      cam: { yaw: spawnYaw, pitch: pit, dist: 6, tx: px, ty: py, tz: eyeZ },
      objs: objs, parts: [], kind: 'rack',
      walk: {
        x: px, y: py, yaw: spawnYaw, pitch: pit, eye: eyeZ, zoom: zoom0,
        stand: stand,
        bounds: { x0: -1.1, y0: -0.5, x1: Wroom + 1.7, y1: Dp + 0.5 },
        blocks: blocks,
        cabs: poses.map(function (p) {
          return {
            name: p.cab.name, title: p.cab.name, fx: p.fx, fy: p.fy, nx: p.nx, ny: p.ny,
            cx: p.cx, cy: p.cy, rx: p.rx, ry: p.ry, halfW: p.halfW
          };
        })
      }
    };
  };

  /* ---- 动态着色：运行灯呼吸 / 告警脉冲 ---- */
  function led(t, seed) {
    var a = 0.45 + 0.55 * (0.5 + 0.5 * Math.sin(t * 2.2 + (seed || 0)));
    return [Math.round(37 * a + 20), Math.round(240 * a + 15), Math.round(168 * a + 40)];
  }
  function pulse(t) {
    var a = 0.45 + 0.55 * (0.5 + 0.5 * Math.sin(t * 3.4));
    return [Math.round(255 * a), Math.round(90 * a), Math.round(106 * a)];
  }
  function pulseSlow(t) {
    var a = 0.35 + 0.5 * (0.5 + 0.5 * Math.sin(t * 1.8));
    return [Math.round(255 * a), Math.round(90 * a), Math.round(106 * a)];
  }

  /* ==================== 查看器 ==================== */
  var active = null;
  /* 全局记录最后一次指针位置：下钻切换视图时光标已停在新画布上方、
     不会再产生 mousemove，新查看器靠它补一次悬停并抑制自动环绕 */
  var lastPointer = null;

  function Viewer(host, kind, key, api, inherit) {
    this.host = host;
    this.kind = kind;
    this.key = key;
    this.api = api || {};
    this.dead = false;
    this.t0 = (typeof performance !== 'undefined' ? performance.now() : Date.now());
    this.lastInput = this.t0;
    this.orbit = true;
    this.hover = null;
    this.sel = null;
    this.drag = null;
    this._mp = null;      /* 指针在画布内的最后位置：帧循环据此持续校正悬停，并暂停自动环绕 */
    this._mpInit = true;  /* 首帧后用 lastPointer 补一次指针位置（切换视图时不会有新 mousemove） */
    this.keys = {};
    this.q = [];
    this.pick = [];
    this.pickLbl = [];   /* 文字标签的屏幕包围盒：标签是屏幕空间的东西，
                           不进多边形拾取表，靠这张表单独判 */
    this.frames = 0;
    this.fps = 0;
    this._fpsT = this.t0;
    this._fpsN = 0;

    var sc = kind === 'rack' ? S.rack(key, { inside: true }) : S[kind](key);
    if (!sc) { this.dead = true; return; }
    this.sc = sc;
    var c = sc.cam;
    this.to = { yaw: c.yaw, pitch: c.pitch, dist: c.dist, tx: c.tx, ty: c.ty, tz: c.tz };
    this.cur = { yaw: this.to.yaw, pitch: this.to.pitch, dist: this.to.dist * 1.9, tx: this.to.tx, ty: this.to.ty, tz: this.to.tz };
    this.home = { yaw: c.yaw, pitch: c.pitch, dist: c.dist, tx: c.tx, ty: c.ty, tz: c.tz };
    this.minD = kind === 'rack' ? 1.1 : 2.2;
    this.maxD = kind === 'rack' ? 12 : (kind === 'room' ? 42 : 150);
    if (inherit && !sc.walk) { this.cur.yaw = this.to.yaw = inherit.yaw; this.cur.pitch = this.to.pitch = inherit.pitch; }
    if (sc.walk) {
      this.walk = true;
      this.player = { x: sc.walk.x, y: sc.walk.y };
      this.eyeZ = sc.walk.eye || 1.18;
      this.zoom = sc.walk.zoom || 0.76;
      this.to.yaw = this.cur.yaw = sc.walk.yaw;
      this.to.pitch = this.cur.pitch = sc.walk.pitch;
      this.orbit = false;
      this.selCab = key;
      /* 进来就打开目标柜的门，人仍站在通道里看全貌。 */
      this.insideView = true;
      this._sceneInside = true;
    }

    this.buildDom();
    this.bind();
    this.resize();
    /* 视距自适应：dist 写死会导致窄画布里内容只占五分之一、
       22 个机柜挤在 20px 间距上（指针差 10px 就选错邻居）。
       首次拿到真实画布尺寸后按包围盒反算视距，让内容填满画面 */
    this._fitPending = 1;
    this.frameFit();
    var self = this;
    this.ro = null;
    if (window.ResizeObserver) {
      this.ro = new ResizeObserver(function () { self.resize(); });
      this.ro.observe(this.cv);
      this.ro.observe(this.host);
    }
    /* 每帧兜底校验：宿主尺寸变了（flex 重排 / 侧栏折叠）要及时跟上 */
    this._lastW = 0;
    this.loop = function (t) { self.frame(t); };
    this.raf = requestAnimationFrame(this.loop);
    if (this.walk) this.syncInsideBtn();
  }

  /* ---------- DOM / HUD ---------- */
  Viewer.prototype.buildDom = function () {
    var self = this;
    var host = this.host;
    host.innerHTML =
      '<canvas class="v3-cv" tabindex="0"></canvas>' +
      '<div class="v3-hud">' +
        '<div class="v3-tl"><b class="v3-lv"></b><span class="v3-sub"></span></div>' +
        '<div class="v3-tr">' +
          '<span class="v3-btn" data-v="iso">等轴</span>' +
          '<span class="v3-btn" data-v="top">顶视</span>' +
          '<span class="v3-btn" data-v="front">正面</span>' +
          '<span class="v3-btn" data-v="side">侧视</span>' +
          '<span class="v3-btn" data-v="reset">复位</span>' +
          '<span class="v3-btn is-on" data-v="orbit">环绕</span>' +
        '</div>' +
        '<div class="v3-bl" style="max-width:calc(100% - 24px);white-space:normal;overflow:visible;bottom:8px;z-index:6">' +
          '<b>左键拖拽</b> 旋转 · <b>滚轮</b> 缩放 · <b>右键/Shift+拖拽</b> 平移 · ' +
          '<b>WASD</b> 漫游 · <b>QE</b> 升降 · <b>双击</b> 飞向 · <b>R</b> 复位 · <b>Esc</b> 返回' +
        '</div>' +
        '<div class="v3-br"><span class="v3-rd"></span></div>' +
        '<div class="v3-tip"></div>' +
      '</div>';
    this.cv = host.querySelector('.v3-cv');
    this.tip = host.querySelector('.v3-tip');
    this.rd = host.querySelector('.v3-rd');
    this.lv = host.querySelector('.v3-lv');
    this.sub = host.querySelector('.v3-sub');
    this.ctx = this.cv.getContext('2d');

    var names = { floor: '楼层立体', room: '机房立体', rack: '机柜 42U' };
    this.lv.textContent = names[this.kind] || '立体';
    /* 副标题：机房名 + 排布概览。原先这条信息是画在画布右下角的标签，
       但它投影后压住最后一列机柜（477px），点击全被它抢走 ——
       概览类信息放 HUD，画布上只留"能指向实体"的标签。 */
    this.sub.textContent = this.key + (this.sc && this.sc.sub ? ' · ' + this.sc.sub : '');

    if (this.walk) {
      var tr = host.querySelector('.v3-tr');
      if (tr) tr.style.display = 'none';
      var bl = host.querySelector('.v3-bl');
      bl.innerHTML = '柜门已打开 · <b>WASD</b> 前后左右、<b>Q/E</b> 蹲起、<b>拖动</b> 转头，滚轮可再拉远或靠近';
      bl.style.maxWidth = 'calc(100% - 24px)';
      bl.style.whiteSpace = 'normal';
      bl.style.overflow = 'visible';
      var hud = host.querySelector('.v3-hud');
      if (hud) {
        var lg = document.createElement('div');
        lg.className = 'v3-lg';
        lg.style.cssText = 'position:absolute;left:12px;top:40px;display:flex;gap:8px;align-items:center;flex-wrap:wrap;' +
          'max-width:calc(100% - 24px);font-size:12px;line-height:1.4;color:#E7F2FF;background:rgba(5,15,27,.78);' +
          'padding:4px 8px;border-radius:4px;z-index:6';
        lg.innerHTML =
          '<span><i style="display:inline-block;width:9px;height:9px;background:#2C7CB5;margin-right:3px;vertical-align:-1px"></i>服务器</span>' +
          '<span><i style="display:inline-block;width:9px;height:9px;background:#12A08C;margin-right:3px;vertical-align:-1px"></i>交换机</span>' +
          '<span><i style="display:inline-block;width:9px;height:9px;background:#7A5CC6;margin-right:3px;vertical-align:-1px"></i>存储</span>' +
          '<span><i style="display:inline-block;width:9px;height:9px;background:#D2703C;margin-right:3px;vertical-align:-1px"></i>网安/网闸</span>' +
          '<span><i style="display:inline-block;width:9px;height:9px;background:#5A6B82;margin-right:3px;vertical-align:-1px"></i>其他</span>';
        var back = document.createElement('button');
        back.type = 'button';
        back.setAttribute('data-act', 'aisle');
        back.textContent = '返回走道';
        back.style.cssText = 'display:none;position:absolute;left:12px;bottom:40px;z-index:7;pointer-events:auto;cursor:pointer;' +
          'border:1px solid rgba(255,209,90,.7);background:rgba(40,28,6,.9);color:#FFD15A;border-radius:4px;padding:4px 10px;font-size:12px;font-weight:600';
        back.addEventListener('click', function (ev) {
          ev.stopPropagation();
          self.exitInside();
        });
        hud.appendChild(lg);
        hud.appendChild(back);
        this._backBtn = back;
      }
    }
    if (this.walk && this.cv) {
      try { this.cv.focus({ preventScroll: true }); } catch (e1) { /* 焦点留给用户点进画布 */ }
    }
    host.querySelectorAll('.v3-btn').forEach(function (b) {
      b.addEventListener('click', function (e) {
        e.stopPropagation();
        var v = b.getAttribute('data-v');
        self.view(v);
        if (v === 'orbit') {
          self.orbit = !self.orbit;
          b.classList.toggle('is-on', self.orbit);
        }
      });
    });
  };

  /* ---------- 视角预设 ---------- */
  Viewer.prototype.view = function (v) {
    if (this.walk) {
      if (v === 'reset' || v === 'front' || v === 'iso') {
        var w = this.sc.walk;
        this.player.x = w.x; this.player.y = w.y;
        this.to.yaw = w.yaw; this.to.pitch = w.pitch;
      }
      this.lastInput = this.now();
      return;
    }
    var c = this.cur, h = this.home;
    if (v === 'iso') { this.to.yaw = h.yaw; this.to.pitch = h.pitch; }
    else if (v === 'top') { this.to.pitch = 1.45; this.to.yaw = 0.001; }
    else if (v === 'front') { this.to.pitch = 0.10; this.to.yaw = 0.0; }
    else if (v === 'side') { this.to.pitch = 0.12; this.to.yaw = 1.45; }
    else if (v === 'reset') { this.to.yaw = h.yaw; this.to.pitch = h.pitch; this.to.dist = h.dist; this.to.tx = h.tx; this.to.ty = h.ty; this.to.tz = h.tz; }
    /* 换姿态后包围盒的投影跨度会变（顶视最扁、侧视最宽），
       重新按当前姿态算一次"刚好填满"的视距，否则切视角就又空出一大片 */
    if (v !== 'orbit') this.frameFit(v === 'reset' ? 0.90 : 0.94);
    this.lastInput = this.now();
  };

  /* ---------- 尺寸 ---------- */
  Viewer.prototype.resize = function () {
    /* 先放开画布上一次写死的像素宽高，否则侧栏收起后宿主变大，画布仍停在旧尺寸。 */
    if (this.cv) { this.cv.style.width = '100%'; this.cv.style.height = '100%'; }
    var w = this.host ? this.host.clientWidth : 0, h = this.host ? this.host.clientHeight : 0;
    if (!w || !h) {
      w = this.cv ? this.cv.clientWidth : 0;
      h = this.cv ? this.cv.clientHeight : 0;
    }
    if (!w || !h) return;
    /* 位图固定超采样：一是边缘更平滑，二是位图必须大于 CSS 盒 ——
       等大时（dpr=1）合成器可能一直挂着空白纹理，画面整个发黑 */
    var bm = (w * h > 900000) ? 1.5 : 2;
    var bw = Math.round(w * bm), bh = Math.round(h * bm);
    if (this.cv.width !== bw || this.cv.height !== bh) {
      this.cv.width = bw;
      this.cv.height = bh;
      this.cv.style.width = w + 'px';
      this.cv.style.height = h + 'px';
      this.needDraw = 1;         /* 位图被清空，下一帧必须重画 */
    }
    this.W = w; this.H = h; this.dpr = bm;
    this.cx = w / 2; this.cy = h / 2;
    this.scale = (h / 2) / Math.tan(0.42);   /* fov ≈ 48° */
    /* 画布尺寸变了就重算视距（窗口缩放、flex 重排都会走到这里） */
    if (this._fitPending || (this._fitW !== w || this._fitH !== h)) {
      this._fitW = w; this._fitH = h;
      this.frameFit(this._fitPending ? 0.90 : 0.92);
      this._fitPending = 0;
      this.needDraw = 1;
      /* 不直接改 cur：保留"由远及近推入"的入场动画，lerp 会自然跟上 */
    }
  };

  /* ---------- 视距自适应 ----------
     把包围盒 8 个角投到当前相机姿态下，量出它们在屏幕上的跨度，
     再按画布宽高比反算"刚好填满"的视距。取宽高两个约束里更紧的那个，
     这样窄画布（机柜栏 836px）和宽画布（机房 1480px）都能填满。 */
  Viewer.prototype.frameFit = function (fill) {
    if (this.walk) return;
    var bb = this.sc.fitBB || this.sc.bb, c = this.cur;
    if (!bb || !this.W) return;
    var k = fill == null ? 0.90 : fill;
    /* 采样包围盒 8 角 + 各面中心，姿态用目标姿态（to），避免动画中途抖动 */
    var src = this.to;
    var cp = Math.cos(src.pitch), sp = Math.sin(src.pitch);
    var cy = Math.cos(src.yaw), sy = Math.sin(src.yaw);
    var ex = src.tx + src.dist * sy * cp, ey = src.ty + src.dist * cy * cp, ez = src.tz + src.dist * sp;
    var fx = src.tx - ex, fy = src.ty - ey, fz = src.tz - ez;
    var fl = Math.sqrt(fx * fx + fy * fy + fz * fz) || 1; fx /= fl; fy /= fl; fz /= fl;
    var rx = fy, ry = -fx, rl = Math.sqrt(rx * rx + ry * ry);
    if (rl < 1e-6) { rx = Math.cos(src.yaw); ry = -Math.sin(src.yaw); rl = 1; }
    rx /= rl; ry /= rl;
    var ux = ry * fz, uy = -rx * fz, uz = rx * fy - ry * fx;

    /* 以 target 为中心，取各角相对 target 的偏移在相机基上的投影极值 */
    var mx = 0, my = 0;
    var pts = [];
    var X = [bb.x0, bb.x1], Y = [bb.y0, bb.y1], Z = [bb.z0, bb.z1];
    for (var i = 0; i < 8; i++) {
      pts.push([X[i & 1], Y[(i >> 1) & 1], Z[(i >> 2) & 1]]);
    }
    pts.push([(bb.x0 + bb.x1) / 2, (bb.y0 + bb.y1) / 2, (bb.z0 + bb.z1) / 2]);
    for (i = 0; i < pts.length; i++) {
      var a = pts[i][0] - src.tx, b = pts[i][1] - src.ty, cc = pts[i][2] - src.tz;
      var vx = Math.abs(a * rx + b * ry);
      var vy = Math.abs(a * ux + b * uy + cc * uz);
      if (vx > mx) mx = vx;
      if (vy > my) my = vy;
    }
    if (mx < 1e-4 || my < 1e-4) return;
    /* 投影公式：screen = center ± scale * v / (dist - v_forward)
       保守取 denom ≈ dist（视距远大于场景尺度），解出两个方向的最小 dist */
    var dW = this.scale * mx / ((this.W / 2) * k);
    var dH = this.scale * my / ((this.H / 2) * k);
    var d = Math.max(dW, dH);
    /* 别把镜头拉得太近（会穿进柜体）或太远 */
    d = clamp(d, this.minD * 1.05, this.maxD * 0.98);
    this.to.dist = d;
    this.home.dist = d;
    this.fitDist = d;
  };

  Viewer.prototype.now = function () {
    return (typeof performance !== 'undefined' ? performance.now() : Date.now());
  };

  /* 悬停拾取：抽成独立方法 —— 帧循环在相机移动时也会复算。
     早期只在 mousemove 里算一次，镜头自转 / 切视角 / 漫游之后，高亮和 tooltip
     还指向旧对象，用户看到的就是"选中的机柜和鼠标位置对不上" */
  Viewer.prototype.hoverAt = function (p) {
    var hit = this.hit(p.x, p.y);
    if (hit !== this.hover) {
      this.hover = hit;
      this.cv.style.cursor = hit ? 'pointer' : 'grab';
      if (hit && this.sc.objs[hit]) {
        this.tip.textContent = this.sc.objs[hit].hint;
        this.tip.classList.add('is-on');
      } else this.tip.classList.remove('is-on');
    }
    if (hit) {
      this.tip.style.left = (p.x + 16) + 'px';
      this.tip.style.top = (p.y + 16) + 'px';
    }
    return hit;
  };

  /* ---------- 输入 ---------- */
  Viewer.prototype.bind = function () {
    var self = this, cv = this.cv;

    /* 视口像素 → 画布自身 CSS 像素。
       整屏 #screen 被 core.js 的 fit() 套了 transform: scale(s)（s = min(窗宽/1920, 窗高/1080)，
       窗口小于 1920×1080 时 s<1）。getBoundingClientRect() 返回的是**缩放后**的视觉盒，
       于是 clientX - rect.left 得到的是"视觉像素"，而引擎内部（cx/cy/scale、拾取多边形、
       标签包围盒）一律用画布自身的 CSS 像素 —— 两者差一个 s。
       后果：指针坐标被压向画布左上角，越往右下偏得越多，点 A 高亮的却是左上方的 B。
       实测 1600×1030（s=0.833）：想点 205机房-3-10，高亮的却是 3-5，偏了 5 台。
       1920×1080 下 s=1，所以这个问题只在"窗口没占满设计尺寸"时暴露。 */
    function local(e) {
      var r = cv.getBoundingClientRect();
      var k = (r.width && cv.clientWidth) ? cv.clientWidth / r.width : 1;
      return { x: (e.clientX - r.left) * k, y: (e.clientY - r.top) * k };
    }
    this._md = function (e) {
      var p = local(e);
      self.cv.focus({ preventScroll: true });
      /* 巡视里鼠标只转头，不平移人所在位置。右键在轨道模式仍是平移。 */
      if (self.walk || !(e.button === 2 || e.shiftKey || e.button === 1)) {
        self.drag = { m: 'rot', x: p.x, y: p.y, sx: p.x, sy: p.y, moved: 0, t: self.now() };
      } else {
        self.drag = { m: 'pan', x: p.x, y: p.y, moved: 0, t: self.now() };
      }
      cv.style.cursor = 'grabbing';
      e.preventDefault();
    };
    this._mm = function (e) {
      lastPointer = { x: e.clientX, y: e.clientY };
      var p = local(e);
      /* 鼠标移动即算交互：自动环绕必须让位 —— 否则用户盯着机柜找的几秒里
         场景自己转走，再点选中的就不是鼠标底下那个柜子 */
      self.lastInput = self.now();
      if (self.drag) {
        var dx, dy;
        if (self.walk && document.pointerLockElement === cv) {
          dx = e.movementX || 0;
          dy = e.movementY || 0;
        } else {
          dx = p.x - self.drag.x;
          dy = p.y - self.drag.y;
          self.drag.x = p.x;
          self.drag.y = p.y;
        }
        self.drag.moved = (self.drag.moved || 0) + Math.abs(dx) + Math.abs(dy);
        if (self.drag.m === 'rot') {
          if (self.walk) {
            /* 第一人称转头：右拖看向右侧，上拖抬头。超过 6px 不算点击。
               拖过阈值后尝试指针锁定；浏览器拒绝时仍靠按住拖动。 */
            self.to.yaw += dx * 0.0052;
            /* 下拖低头、上拖抬头。俯仰为 0 时视线水平，和眼前这一排 U 位正对。 */
            self.to.pitch = clamp(self.to.pitch - dy * 0.007, -1.25, 0.85);
            if (!self._plNo && !self.drag.asked && self.drag.moved > 8 && cv.requestPointerLock) {
              self.drag.asked = 1;
              var rq = cv.requestPointerLock();
              if (rq && rq.catch) rq.catch(function () { self._plNo = 1; });
            }
          } else {
            /* 抓取感：物体跟着手走 —— 右拖场景右移、左拖场景左移；
               下拖应"看到更多柜顶"（pitch 减小 = 相机降低 = 俯角变大）
               这里 pitch 是相机仰角，手往下拖要看到顶面 → pitch 减小 */
            self.to.yaw -= dx * 0.0058;
            self.to.pitch = clamp(self.to.pitch - dy * 0.0058, -1.35, 1.48);
          }
        } else {
          self.panBy(-dx, dy);
        }
        return;
      }
      /* 指针在画布内：记下位置，帧循环持续校正悬停目标；离开则收起 tooltip */
      if (p.x >= 0 && p.y >= 0 && p.x <= self.W && p.y <= self.H) {
        self._mp = p;
        self.hoverAt(p);
      } else {
        self._mp = null;
        if (self.hover) {
          self.hover = null;
          self.tip.classList.remove('is-on');
          cv.style.cursor = 'grab';
        }
      }
    };
    this._ml = function () {
      if (self.drag) return;
      self._mp = null;
      if (self.hover) {
        self.hover = null;
        self.tip.classList.remove('is-on');
        cv.style.cursor = 'grab';
      }
    };
    this._mu = function (e) {
      if (!self.drag) return;
      var d = self.drag; self.drag = null;
      if (document.pointerLockElement === cv && document.exitPointerLock) document.exitPointerLock();
      cv.style.cursor = self.hover ? 'pointer' : 'grab';
      var p = local(e);
      var sx = d.sx != null ? d.sx : p.x, sy = d.sy != null ? d.sy : p.y;
      if (d.m === 'rot' && d.moved < 6 && self.now() - d.t < 420) self.click(sx, sy);
    };
    this._wc = function (e) {
      e.preventDefault();
      if (self.walk) {
        var zk = Math.exp(-clamp(e.deltaY, -120, 120) * 0.0015);
        self.zoom = clamp((self.zoom || 1) * zk, 0.7, 3.8);
        self.needDraw = 1;
        self.lastInput = self.now();
        return;
      }
      var k = Math.exp(clamp(e.deltaY, -120, 120) * 0.0013);
      self.to.dist = clamp(self.to.dist * k, self.minD, self.maxD);
      self.lastInput = self.now();
    };
    this._cc = function (e) { e.preventDefault(); };
    this._db = function (e) {
      e.preventDefault();
      if (self._clk) self._clk.dbl = 1;
      clearTimeout(self._clkTimer);
      var p = local(e);
      var hit = self.hit(p.x, p.y);
      if (!self.walk) { if (hit) self.flyTo(hit); return; }
      if (self.insideView) {
        if (!hit) { self.exitInside(); return; }
        var din = self.sc.objs[hit];
        if (!din || din.kind === 'cab') {
          if (din && din.kind === 'cab') self.enterInside(din.name);
          else self.exitInside();
          return;
        }
        if (din.kind === 'dev') {
          self.sel = hit;
          if (self.api.toast) self.api.toast(din.hint);
          if (self.api.onDev) self.api.onDev(din.sn);
        }
        return;
      }
      if (!hit) return;
      var d = self.sc.objs[hit];
      if (d && (d.kind === 'cab' || d.kind === 'rack')) self.enterInside(d.name);
    };
    /* 触摸：单指旋转、双指缩放+平移 */
    this._ts = function (e) {
      if (e.touches.length === 1) {
        self.drag = { m: 'rot', x: e.touches[0].clientX, y: e.touches[0].clientY, moved: 0, t: self.now() };
      } else if (e.touches.length === 2) {
        var d = Math.hypot(e.touches[0].clientX - e.touches[1].clientX, e.touches[0].clientY - e.touches[1].clientY);
        self.drag = { m: 'pin', d: d, x: (e.touches[0].clientX + e.touches[1].clientX) / 2, y: (e.touches[0].clientY + e.touches[1].clientY) / 2 };
      }
      self.lastInput = self.now();
    };
    this._tm = function (e) {
      if (!self.drag) return;
      e.preventDefault();
      if (self.drag.m === 'rot' && e.touches.length === 1) {
        var dx = e.touches[0].clientX - self.drag.x, dy = e.touches[0].clientY - self.drag.y;
        self.drag.x = e.touches[0].clientX; self.drag.y = e.touches[0].clientY;
        self.to.yaw -= dx * 0.006;
        self.to.pitch = clamp(self.to.pitch - dy * 0.006, -1.35, 1.48);
      } else if (self.drag.m === 'pin' && e.touches.length === 2) {
        var d = Math.hypot(e.touches[0].clientX - e.touches[1].clientX, e.touches[0].clientY - e.touches[1].clientY);
        var cx = (e.touches[0].clientX + e.touches[1].clientX) / 2, cy = (e.touches[0].clientY + e.touches[1].clientY) / 2;
        if (self.drag.d) self.to.dist = clamp(self.to.dist * (self.drag.d / d), self.minD, self.maxD);
        self.panBy(self.drag.x - cx, cy - self.drag.y);
        self.drag.d = d; self.drag.x = cx; self.drag.y = cy;
      }
      self.lastInput = self.now();
    };
    this._te = function () { self.drag = null; };

    /* 键盘：绑 document，只在 3D 激活且焦点不在输入框时生效；Esc 不拦 */
    this._kd = function (e) {
      if (active !== self) return;
      var t = e.target;
      if (t && (t.tagName === 'INPUT' || t.tagName === 'TEXTAREA' || t.isContentEditable)) return;
      var k = e.key.toLowerCase();
      if (self.walk) {
        if (k === 'escape') {
          if (self.insideView || (self.slide && (self.slide.mode === 'in' || self.slide.mode === 'out'))) {
            self.exitInside();
          }
          if (document.pointerLockElement && document.exitPointerLock) document.exitPointerLock();
          return;
        }
        if (k === '[' || k === ']' || e.code === 'BracketLeft' || e.code === 'BracketRight') {
          var zin = (k === ']' || e.code === 'BracketRight') ? 1.12 : (1 / 1.12);
          self.zoom = clamp((self.zoom || 1) * zin, 0.7, 3.8);
          self.needDraw = 1;
          e.preventDefault();
          return;
        }
        if (k === 'w' || k === 'a' || k === 's' || k === 'd' || k === 'q' || k === 'e' ||
            k === 'arrowup' || k === 'arrowdown' || k === 'arrowleft' || k === 'arrowright') {
          self.keys[k] = 1;
          e.preventDefault();
          return;
        }
        if (k === 'r' && self.sc.walk) {
          if (self.insideView || (self.slide && self.slide.mode === 'in')) {
            self.exitInside();
          } else {
            var w0 = self.sc.walk;
            self.player.x = w0.x; self.player.y = w0.y;
            self.to.yaw = w0.yaw; self.to.pitch = w0.pitch;
            self.eyeZ = w0.eye || 1.62;
            self.zoom = 1;
            self.needDraw = 1;
          }
          e.preventDefault();
        }
        return;
      }
      if (k === 'escape') return;
      var used = true;
      if (k === 'w' || k === 'arrowup') self.move(0, 1);
      else if (k === 's' || k === 'arrowdown') self.move(0, -1);
      else if (k === 'a' || k === 'arrowleft') self.move(-1, 0);
      else if (k === 'd' || k === 'arrowright') self.move(1, 0);
      else if (k === 'q') self.to.tz = clamp(self.to.tz + self.to.dist * 0.05, -3, 12);
      else if (k === 'e') self.to.tz = clamp(self.to.tz - self.to.dist * 0.05, -3, 12);
      else if (k === 'r') self.view('reset');
      else if (k === 'f') { if (self.sel) self.flyTo(self.sel); }
      else if (k === 'z') self.view('top');
      else if (k === 'x') self.view('front');
      else if (k === 'c') self.view('iso');
      else used = false;
      if (used) { e.preventDefault(); self.lastInput = self.now(); }
    };

    cv.addEventListener('mousedown', this._md);
    window.addEventListener('mousemove', this._mm);
    window.addEventListener('mouseup', this._mu);
    cv.addEventListener('wheel', this._wc, { passive: false });
    cv.addEventListener('contextmenu', this._cc);
    cv.addEventListener('dblclick', this._db);
    cv.addEventListener('mouseleave', this._ml);
    cv.addEventListener('touchstart', this._ts, { passive: true });
    cv.addEventListener('touchmove', this._tm, { passive: false });
    cv.addEventListener('touchend', this._te);
    this._ku = function (e) {
      if (!e.key) return;
      self.keys[e.key.toLowerCase()] = 0;
    };
    this._blur = function () { self.keys = {}; };
    window.addEventListener('keydown', this._kd);
    window.addEventListener('keyup', this._ku);
    window.addEventListener('blur', this._blur);
    cv.style.cursor = 'grab';
  };

  /* 屏幕位移 → 世界平移（沿相机右向量 / 世界上方向） */
  Viewer.prototype.panBy = function (dx, dy) {
    var c = this.cur;
    var cy = Math.cos(c.yaw), sy = Math.sin(c.yaw);
    var k = c.dist * 0.0016;
    this.to.tx += (-dx * cy) * k;
    this.to.ty += (dx * sy) * k;
    this.to.tz = clamp(this.to.tz + dy * k, -3, 14);
  };
  /* WASD：沿相机水平朝向前后左右（W=深入场景，与游戏一致） */
  Viewer.prototype.move = function (rx, rz) {
    var c = this.cur;
    var cy = Math.cos(c.yaw), sy = Math.sin(c.yaw);
    var k = c.dist * 0.045;
    /* 水平前向 = -（sy, cy)（eye 在 target 的 +（sy,cy) 一侧）；右向 = (-cy, sy) */
    this.to.tx += (-sy * rz - cy * rx) * k;
    this.to.ty += (-cy * rz + sy * rx) * k;
  };

  /* ---------- 拾取 ----------
     不能靠"画家顺序倒扫"决定命中：dep 是面中心平均深度，俯视角下
     近处柜顶（dep 小）与远处柜正面（dep 大）屏幕投影重叠，画家算法
     为正确成像会把顶面画在远处柜正面之上，于是倒扫必然命中近处柜 ——
     正是"鼠标指着远处柜、却选中近处柜"的病灶。
     改为：在所有覆盖该点的多边形里，取该点真实深度最小的那个。
     穿越次数与深度累积必须解耦：边遍历顺序和"进入/离开"无关，
     若把深度累积挂在 inside 翻转上，结果会随多边形起始顶点变化。 */
  Viewer.prototype.hit = function (sx, sy) {
    var list = this.pick, i, a, b, it, p, n, pa, pb, t, xc, cross, pd, za, zb, z;

    /* ---- 第一优先：屏幕空间文字 ----
       标签（T()）是画在屏幕上的字，不属于任何多边形，射线会直接穿过去
       打到字背后的地板或邻房 —— 表现为"点 209 空调却选中 207 空调间"。
       标签在渲染上是覆盖层（画在最上面），用户看到的就是那行字，
       所以点它就该选它代表的对象：命中文字矩形即返回其 obj。
       多个标签重叠时取离相机最近的那个（渲染时也正是它压在最上面）。 */
    var lb = this.pickLbl, lbHit = null, lbDep = 1e30;
    for (i = 0; i < lb.length; i++) {
      it = lb[i];
      if (sx >= it.x0 && sx <= it.x1 && sy >= it.y0 && sy <= it.y1) {
        if (it.dep < lbDep) { lbDep = it.dep; lbHit = it; }
      }
    }
    if (lbHit) return lbHit.obj;

    /* 第一遍：命中判定 + 记录每个命中的真实深度 */
    var hits = this._hits || (this._hits = []);
    hits.length = 0;
    var minD = 1e30;
    for (i = 0; i < list.length; i++) {
      it = list[i]; p = it.poly; n = p.length;
      cross = 0; pd = 1e30;
      for (a = 0, b = n - 1; a < n; b = a++) {
        pa = p[a]; pb = p[b];
        if ((pa[1] > sy) === (pb[1] > sy)) continue;
        t = (sy - pa[1]) / (pb[1] - pa[1]);
        xc = pa[0] + (pb[0] - pa[0]) * t;
        if (sx >= xc) continue;
        cross++;
        za = pa.length > 2 ? pa[2] : it.dep;
        zb = pb.length > 2 ? pb[2] : it.dep;
        z = za + (zb - za) * t;
        if (z < pd) pd = z;
      }
      if (cross & 1) {
        hits.push({ obj: it.obj, dep: pd, pri: it.pri || 0 });
        if (pd < minD) minD = pd;
      }
    }
    if (!hits.length) return null;
    /* 第二遍：深度最浅的命中为基准。设备面板（pri=1）只有 80×14 px，
       常与柜体侧板/后板在屏幕上重叠且深度只差几毫米 —— 这类"面挨着面"
       的贴合关系下让设备赢，用户点了设备就不会选中背后的柜体。
       容差 0.15m 保证真正被柜体挡在后面的设备仍然选不中。 */
    var TOL = 0.15, best = null, bestP = -1, bestD = 1e30;
    for (i = 0; i < hits.length; i++) {
      it = hits[i];
      if (it.dep > minD + TOL) continue;
      if (it.pri > bestP || (it.pri === bestP && it.dep < bestD)) {
        bestP = it.pri; bestD = it.dep; best = it.obj;
      }
    }
    return best || hits[0].obj;
  };
  Viewer.prototype.click = function (sx, sy) {
    var o = this.hit(sx, sy);
    if (!o) { this.sel = null; return; }
    var d = this.sc.objs[o];
    if (!d) { this.sel = o; return; }
    if (this.walk && (d.kind === 'cab' || d.kind === 'rack')) {
      this.armSelect(d.name);
      return;
    }
    this.sel = o;
    if (d.drill && this.api.drill) this.api.drill(d.drill, d.kind);
    else if (d.kind === 'dev') {
      if (this.api.toast) this.api.toast(d.hint);
      if (this.api.onDev) this.api.onDev(d.sn);
    } else if (this.api.toast) this.api.toast(d.hint);
  };
  Viewer.prototype.cabByName = function (name) {
    var list = (this.sc.walk && this.sc.walk.cabs) || [];
    var i;
    for (i = 0; i < list.length; i++) if (list[i].name === name) return list[i];
    return null;
  };
  /* 单击：立刻打开这一台的设备面板，人留在通道里。双击再平滑走到正对。 */
  Viewer.prototype.armSelect = function (name) {
    if (!name || !this.walk) return;
    var now = this.now();
    if (this._clk && this._clk.name === name && now - this._clk.t < 420) {
      this._clk = null;
      this.enterInside(name);
      return;
    }
    this._clk = { name: name, t: now, dbl: 0 };
    if (!this.insideView) {
      this._aislePose = {
        x: this.player.x, y: this.player.y,
        yaw: this.to.yaw, pitch: this.to.pitch, zoom: this.zoom || 1
      };
    }
    this.selCab = name;
    this._pickHold = now + 5200;
    if (this.api && this.api.onCab) this.api.onCab(name);
    this.switchCab(name, true);
    this.needDraw = 1;
  };
  Viewer.prototype.facingCab = function (c) {
    if (!c) return false;
    var lookx = Math.sin(this.to.yaw), looky = Math.cos(this.to.yaw);
    return lookx * (-c.nx) + looky * (-c.ny) > 0.55;
  };
  Viewer.prototype.slideToFront = function (name, from) {
    var c = this.cabByName(name);
    if (!c || !this.player || !this.facingCab(c)) return;
    if (this.slide && this.slide.name === name && this.slide.mode === 'row') return;
    var stand = (this.sc.walk && this.sc.walk.stand) || 0.74;
    this.slide = {
      mode: 'row',
      x0: this.player.x, y0: this.player.y,
      x1: c.fx + c.nx * stand, y1: c.fy,
      yaw0: this.to.yaw, yaw1: this.to.yaw,
      pitch0: this.to.pitch, pitch1: this.to.pitch,
      t0: this.now(), dur: 560, name: name, from: from || '', swapped: true
    };
  };
  /* 双击：平滑推到柜门前，再换成内部面板。人仍在门外，只是这一台改画内部。 */
  Viewer.prototype.enterInside = function (name) {
    if (!name || !this.walk || !this.player) return;
    if (this.slide && this.slide.mode === 'in' && this.slide.name === name) return;
    if (this._clk) this._clk.dbl = 1;
    clearTimeout(this._clkTimer);
    var c = this.cabByName(name);
    if (!c) return;
    if (!this._aislePose || !this.insideView) {
      this._aislePose = {
        x: this.player.x, y: this.player.y,
        yaw: this.to.yaw, pitch: this.to.pitch, zoom: this.zoom || 1,
        sel: this.selCab
      };
    }
    var stand = (this.sc.walk && this.sc.walk.stand) || 1.38;
    var yaw = Math.atan2(-c.nx, -(c.ny || 0));
    this.selCab = name;
    this._pickHold = this.now() + 8000;
    if (this.api && this.api.onCab) this.api.onCab(name);
    this.slide = {
      mode: 'in',
      x0: this.player.x, y0: this.player.y,
      x1: c.fx + c.nx * stand, y1: c.fy,
      yaw0: this.to.yaw, yaw1: yaw,
      pitch0: this.to.pitch, pitch1: (this.sc.walk && this.sc.walk.pitch) || -0.16,
      zoom0: this.zoom || 1, zoom1: (this.sc.walk && this.sc.walk.zoom) || 0.76,
      t0: this.now(), dur: 640, name: name, swapped: false
    };
    this.needDraw = 1;
  };
  Viewer.prototype.exitInside = function () {
    if (this.slide && this.slide.mode === 'out') return;
    if (!this.insideView && !(this.slide && this.slide.mode === 'in')) return;
    var pose = this._aislePose || {
      x: this.player.x, y: this.player.y, yaw: this.to.yaw, pitch: 0, zoom: 1
    };
    this.slide = {
      mode: 'out',
      x0: this.player.x, y0: this.player.y,
      x1: pose.x, y1: pose.y,
      yaw0: this.to.yaw, yaw1: pose.yaw,
      pitch0: this.to.pitch, pitch1: pose.pitch || 0,
      zoom0: this.zoom || 1, zoom1: pose.zoom || 1,
      t0: this.now(), dur: 620,
      name: this.selCab || this.key, swapped: false
    };
    this.needDraw = 1;
  };
  Viewer.prototype.syncInsideBtn = function () {
    if (this._backBtn) this._backBtn.style.display = this.insideView ? '' : 'none';
    if (this.api && this.api.onInside) this.api.onInside(!!this.insideView);
  };
  Viewer.prototype.stepSlide = function (t) {
    var s = this.slide;
    if (!s) return;
    var u = (t - s.t0) / s.dur;
    if (u > 1) u = 1;
    if (u < 0) u = 0;
    var e = u < 0.5 ? 2 * u * u : 1 - Math.pow(-2 * u + 2, 2) / 2;
    this.player.x = s.x0 + (s.x1 - s.x0) * e;
    this.player.y = s.y0 + (s.y1 - s.y0) * e;
    this.keepOutside();
    this.to.tx = this.cur.tx = this.player.x;
    this.to.ty = this.cur.ty = this.player.y;
    var dyaw = s.yaw1 - s.yaw0;
    while (dyaw > Math.PI) dyaw -= Math.PI * 2;
    while (dyaw < -Math.PI) dyaw += Math.PI * 2;
    this.to.yaw = this.cur.yaw = s.yaw0 + dyaw * e;
    this.to.pitch = this.cur.pitch = s.pitch0 + (s.pitch1 - s.pitch0) * e;
    if (s.zoom0 != null) this.zoom = s.zoom0 + (s.zoom1 - s.zoom0) * e;
    this.needDraw = 1;
    if (s.mode === 'in' && !s.swapped && u >= 0.4) {
      s.swapped = true;
      this.switchCab(s.name, true);
    }
    if (s.mode === 'out' && !s.swapped && u >= 0.12) {
      s.swapped = true;
      this.switchCab(s.name, false);
    }
    if (u >= 1) {
      this.slide = null;
      if (s.mode === 'in') {
        if (!s.swapped) this.switchCab(s.name, true);
        this.insideView = true;
        this._sceneInside = true;
        this.syncInsideBtn();
      } else if (s.mode === 'out') {
        if (!s.swapped) this.switchCab(s.name, false);
        this.insideView = false;
        this._sceneInside = false;
        this._aislePose = null;
        this.syncInsideBtn();
      }
    }
  };
  /* 滚轮/[ ] 的 1 倍：整柜顶到底落在当前画布高度里。 */
  Viewer.prototype.fitRack = function () {
    if (!this.walk) return;
    this.zoom = 1;
    this.to.pitch = -0.10;
    this.needDraw = 1;
  };
  /* 镜头飞向对象 */
  Viewer.prototype.flyTo = function (id) {
    var d = this.sc.objs[id];
    if (!d) return;
    if (this.walk) {
        if (d.kind === 'cab' || d.kind === 'rack') {
          this.armSelect(d.name);
          return;
        }
      this.sel = id;
      if (d.kind === 'dev' && this.api.onDev) this.api.onDev(d.sn);
      this.lastInput = this.now();
      return;
    }
    var f = d.focus;
    this.to.tx = f[0]; this.to.ty = f[1]; this.to.tz = f[2];
    var nd = d.kind === 'dev' ? this.minD * 2.4 : d.kind === 'rack' ? this.home.dist : this.to.dist;
    this.to.dist = clamp(nd, this.minD, this.maxD);
    this.sel = id;
    this.lastInput = this.now();
  };

  /* ---------- 帧 ---------- */
  Viewer.prototype.hitsBlock = function (x, y, rad) {
    var blocks = this.sc.walk && this.sc.walk.blocks;
    if (!blocks) return false;
    var i, b;
    for (i = 0; i < blocks.length; i++) {
      b = blocks[i];
      if (x >= b.x0 - rad && x <= b.x1 + rad && y >= b.y0 - rad && y <= b.y1 + rad) return true;
    }
    return false;
  };
  /* 按住 WASD 沿当前朝向走。先试横向再试纵向，贴到机柜就停，不穿进去。 */
  Viewer.prototype.stepWalk = function (t) {
    if (this.slide) return;
    var dt = Math.min(0.05, Math.max(0, (t - (this._lt || t)) / 1000));
    this._lt = t;
    var kk = this.keys || {};
    var f = (kk.w || kk.arrowup ? 1 : 0) - (kk.s || kk.arrowdown ? 1 : 0);
    var s = (kk.d || kk.arrowright ? 1 : 0) - (kk.a || kk.arrowleft ? 1 : 0);
    if (!f && !s) return;
    var yaw = this.cur.yaw;
    var spd = 1.25 * dt;
    var fx = Math.sin(yaw), fy = Math.cos(yaw);
    var dx = (fx * f + fy * s) * spd;
    var dy = (fy * f - fx * s) * spd;
    /* 开门后要能贴近柜门左右挪。半径小一点，柜体本身仍挡着，走不进去。 */
    var nx = this.player.x + dx, ny = this.player.y + dy, rad = this.insideView ? 0.16 : 0.28;
    if (!this.hitsBlock(nx, this.player.y, rad)) this.player.x = nx;
    if (!this.hitsBlock(this.player.x, ny, rad)) this.player.y = ny;
    this.keepOutside();
    var b = this.sc.walk && this.sc.walk.bounds;
    if (b) {
      this.player.x = clamp(this.player.x, b.x0 + rad, b.x1 - rad);
      this.player.y = clamp(this.player.y, b.y0 + rad, b.y1 - rad);
    }
    this.to.tx = this.player.x;
    this.to.ty = this.player.y;
    this.needDraw = 1;
    this.lastInput = t;
  };
  /* Q 降低视高，E 升高。最低约 0.22 米，可平视最底下的 U 位；最高约 1.9 米。 */
  Viewer.prototype.stepEye = function (t) {
    var kk = this.keys || {};
    var dir = (kk.e ? 1 : 0) - (kk.q ? 1 : 0);
    if (!dir) return;
    var dt = Math.min(0.05, Math.max(0.001, (t - (this._et || t)) / 1000));
    this._et = t;
    var next = clamp((this.eyeZ || 1.62) + dir * 1.05 * dt, 0.22, 1.9);
    if (Math.abs(next - this.eyeZ) < 1e-4) return;
    this.eyeZ = next;
    this.needDraw = 1;
    this.lastInput = t;
  };
  /* 人若被推进柜体，沿最近的柜门法线推回通道。转头不改变站位，也就不会转进柜子里面。 */
  Viewer.prototype.keepOutside = function () {
    var blocks = this.sc.walk && this.sc.walk.blocks;
    if (!blocks || !this.player) return;
    var x = this.player.x, y = this.player.y, i, b, guard = 0;
    while (guard++ < 4 && this.hitsBlock(x, y, 0.02)) {
      var best = null, bestD = 1e9;
      for (i = 0; i < blocks.length; i++) {
        b = blocks[i];
        if (x < b.x0 - 0.02 || x > b.x1 + 0.02 || y < b.y0 - 0.02 || y > b.y1 + 0.02) continue;
        var dl = Math.abs(x - b.x0), dr = Math.abs(b.x1 - x), db = Math.abs(y - b.y0), dt = Math.abs(b.y1 - y);
        var m = Math.min(dl, dr, db, dt);
        if (m < bestD) {
          bestD = m;
          if (m === dl) best = { x: b.x0 - 0.08, y: y };
          else if (m === dr) best = { x: b.x1 + 0.08, y: y };
          else if (m === db) best = { x: x, y: b.y0 - 0.08 };
          else best = { x: x, y: b.y1 + 0.08 };
        }
      }
      if (!best) break;
      x = best.x; y = best.y;
    }
    this.player.x = x;
    this.player.y = y;
  };
  /* 走近或转头都不换内部模型。当前柜只由单击选中、双击进入。 */
  Viewer.prototype.considerCab = function () {
    return;
  };
  Viewer.prototype.switchCab = function (name, inside) {
    var want = !!inside;
    if (!name) return;
    if (name === this.key && !!this._sceneInside === want) return;
    var pose = { x: this.player.x, y: this.player.y, yaw: this.to.yaw, pitch: this.to.pitch };
    var sc = S.rack(name, { inside: want });
    if (!sc || !sc.walk) return;
    this.key = name;
    this.sc = sc;
    this.sel = null;
    this.hover = null;
    this.player.x = pose.x;
    this.player.y = pose.y;
    this.to.yaw = pose.yaw;
    this.to.pitch = pose.pitch;
    this._sceneInside = want;
    this.insideView = want;
    this.keepOutside();
    if (this.sub) this.sub.textContent = name + (sc.sub ? ' · ' + sc.sub : '');
    if (this.api && this.api.onCab) this.api.onCab(name);
    this.syncInsideBtn();
    this.needDraw = 1;
  };

  Viewer.prototype.frame = function (t) {
    if (this.dead) return;
    var self = this;
    this.raf = requestAnimationFrame(this.loop);
    var time = (t - this.t0) / 1000;

    /* 走道里转头跟手，轨道相机仍稍缓。 */
    var c = this.cur, to = this.to, sm = this.walk ? 0.55 : 0.16;
    c.yaw = lerp(c.yaw, to.yaw, sm);
    c.pitch = lerp(c.pitch, to.pitch, sm);
    c.dist = lerp(c.dist, to.dist, sm);
    c.tx = lerp(c.tx, to.tx, sm);
    c.ty = lerp(c.ty, to.ty, sm);
    c.tz = lerp(c.tz, to.tz, sm);

    /* 空闲自动环绕：指针停在画布上就不转 —— 用户正在挑机柜，
       镜头自顾自转会让"选中的柜子"和"鼠标底下的柜子"脱节 */
    if (this.orbit && !this._mp && !this.drag && t - this.lastInput > 4000) to.yaw += 0.0022;
    if (this.walk && this.player) {
      if (this.slide) this.stepSlide(t);
      this.stepWalk(t);
      this.stepEye(t);
    }

    /* 静止判定：相机落位后降低重绘频率（呼吸/粒子仍有生命感），
       但保持每帧都提交一次 —— 低频重绘会让合成器读到陈旧纹理（表现为画面发黑） */
    var moving = Math.abs(c.yaw - to.yaw) + Math.abs(c.pitch - to.pitch) + Math.abs(c.dist - to.dist) * 0.2 +
      Math.abs(c.tx - to.tx) + Math.abs(c.ty - to.ty) + Math.abs(c.tz - to.tz);
    var idle = moving <= 0.0015 && !this.drag;
    /* 宿主尺寸兜底校验（RO 之外的补丁） */
    var hw = this.host ? this.host.clientWidth : this.cv.clientWidth;
    var hh = this.host ? this.host.clientHeight : this.cv.clientHeight;
    if (hw !== this._lastW || hh !== this._lastH) { this._lastW = hw; this._lastH = hh; this.resize(); }
    this.frames++;
    if (idle && this.frames % 3 !== 0 && !this.needDraw) return;
    this.needDraw = 0;

    this.draw(time);
    /* 首帧：光标可能已经停在新画布上方（下钻/切视图不产生新的 mousemove）——
       借上一次已知指针位置补一次悬停，否则自动环绕会把场景转走，
       高亮和 tooltip 也要等用户动一下鼠标才出现 */
    if (this._mpInit) {
      this._mpInit = false;
      if (lastPointer) {
        var rr = this.cv.getBoundingClientRect();
        /* 同样要抵掉整屏 transform: scale(s)，否则补算出的悬停点也是压缩过的 */
        var kk = (rr.width && this.cv.clientWidth) ? this.cv.clientWidth / rr.width : 1;
        var lx = (lastPointer.x - rr.left) * kk, ly = (lastPointer.y - rr.top) * kk;
        if (lx >= 0 && ly >= 0 && lx <= this.W && ly <= this.H) {
          this._mp = { x: lx, y: ly };
          this.hoverAt(this._mp);
        }
      }
    }
    /* 相机在动而鼠标没动（自动环绕 / 视角预设 / WASD 漫游 / 双击飞向）时，
       悬停目标要跟着场景重算 —— 保证高亮和 tooltip 始终指向鼠标底下那个对象 */
    if (this._mp && !this.drag && !idle) this.hoverAt(this._mp);
    /* 相机移动时读数实时跟上（否则 HUD 方位/俯仰会滞后整整一秒） */
    if (!idle && this.frames % 6 === 0) this.readout();
    this._fpsN++;
    if (t - this._fpsT > 1000) { this.fps = this._fpsN; this._fpsN = 0; this._fpsT = t; this.readout(); }
  };

  Viewer.prototype.readout = function () {
    if (this.walk) {
      var nm = this.selCab || this.key || '';
      var tag = this.insideView ? '内部' : '走道';
      var pit = this.cur ? Math.round(this.cur.pitch * 180 / Math.PI) : 0;
      this.rd.textContent = tag + ' · ' + nm + ' · 视高 ' + (this.eyeZ || 1.62).toFixed(2) +
        'm · 俯仰 ' + pit + '° · ' + this.fps + 'fps';
      return;
    }
    var c = this.cur;
    var deg = Math.round((-c.yaw * 180 / Math.PI) % 360);
    if (deg < 0) deg += 360;
    this.rd.textContent = '方位 ' + deg + '° · 俯仰 ' + Math.round(c.pitch * 180 / Math.PI) +
      '° · 视距 ' + c.dist.toFixed(1) + 'm · ' + this.fps + 'fps';
  };

  /* ---------- 渲染 ---------- */
  Viewer.prototype.draw = function (time) {
    var ctx = this.ctx, W = this.W, H = this.H, dpr = this.dpr;
    if (!W) return;
    var c = this.cur;

    /* 相机基向量 */
    var cp = Math.cos(c.pitch), sp = Math.sin(c.pitch);
    var cy = Math.cos(c.yaw), sy = Math.sin(c.yaw);
    /* 轨道：eye = target + dist * [sy*cp, cy*cp, sp]（从 +y 一侧看）
       走动：眼睛就在人的位置，朝向 (sin yaw, cos yaw, sin pitch) */
    var ex, ey, ez, fx, fy, fz;
    if (this.walk && this.player) {
      ex = this.player.x; ey = this.player.y; ez = this.eyeZ || 1.15;
      fx = sy * cp; fy = cy * cp; fz = sp;
    } else {
      ex = c.tx + c.dist * sy * cp; ey = c.ty + c.dist * cy * cp; ez = c.tz + c.dist * sp;
      fx = c.tx - ex; fy = c.ty - ey; fz = c.tz - ez;
      var fl0 = Math.sqrt(fx * fx + fy * fy + fz * fz) || 1;
      fx /= fl0; fy /= fl0; fz /= fl0;
    }
    /* right = normalize(forward × worldUp)，worldUp = +z（柜/房都沿 z 向立高）。
       写反参照轴会让整台相机绕视线 roll 30~50°：机柜看起来是"躺"着的 */
    var rx = fy, ry = -fx, rz = 0;
    var rl = Math.sqrt(rx * rx + ry * ry);
    if (rl < 1e-6) {                     /* 近垂直俯视退化：用 yaw 导水平右向量 */
      rx = Math.cos(c.yaw); ry = -Math.sin(c.yaw); rz = 0; rl = 1;
    }
    rx /= rl; ry /= rl;
    /* up = right × forward（右手系：屏幕上方为 +z） */
    var ux = ry * fz - rz * fy, uy = rz * fx - rx * fz, uz = rx * fy - ry * fx;

    var scale = this.scale;
    if (this.walk) {
      /* 约 70° 竖向视野，配合 1.6 米视高，像人站着平视。 */
      var fitTan = 0.70;
      scale = (H / 2) / (fitTan / clamp(this.zoom || 1, 0.7, 3.8));
    }
    this._scrScale = scale;
    /* 雾按层级配系数：机柜/机房近距离要看清面板，楼层俯视全景不该发糊 */
    var fg = this.kind === 'rack' ? [1.4, 5.0] : this.kind === 'room' ? [1.2, 4.0] : [2.2, 6.5];
    var fogN = this.walk ? 8.0 : c.dist * fg[0], fogF = this.walk ? 16 : c.dist * fg[1];

    ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
    /* 背景：上深下浅的舞台渐变（比雾色亮一档，暗面才不会融进背景） */
    var g = ctx.createLinearGradient(0, 0, 0, H);
    g.addColorStop(0, '#08111F');
    g.addColorStop(0.55, '#0C1B2F');
    g.addColorStop(1, '#11283E');
    ctx.fillStyle = g;
    ctx.fillRect(0, 0, W, H);

    var list = this.pick = [];
    var lbls = this.pickLbl = [];
    var emphMap = this._emph || (this._emph = {});
    var ek;
    for (ek in emphMap) delete emphMap[ek];
    if (this.walk && this.player && this.sc.walk && this.sc.walk.cabs) {
      var sel = this.selCab || this.key;
      var slide = this.slide;
      var su = 0;
      var blend = slide && slide.mode === 'row' && slide.from && slide.from !== slide.name;
      if (blend) {
        var uu = (this.now() - slide.t0) / slide.dur;
        if (uu < 0) uu = 0;
        if (uu > 1) uu = 1;
        su = uu < 0.5 ? 2 * uu * uu : 1 - Math.pow(-2 * uu + 2, 2) / 2;
      }
      var pcabs = this.sc.walk.cabs, pi;
      for (pi = 0; pi < pcabs.length; pi++) {
        var pc = pcabs[pi];
        var emph = 0;
        if (blend) {
          if (pc.name === slide.name) emph = su;
          else if (pc.name === slide.from) emph = 1 - su;
        } else if (pc.name === sel || (this.insideView && pc.name === this.key)) emph = 1;
        emphMap[pc.name] = emph;
      }
    }
    var q = this.q; q.length = 0;
    var bxs = this.sc.faces, bxl = this.sc.lines, bxt = this.sc.labels;
    var i, j, f, p, n, cx2, cy2, cz2, vx, vy, vz, dep, dot, k, sx, sy, ar, ok, plateCam;

    /* ---- 面 ---- */
    for (i = 0; i < bxs.length; i++) {
      f = bxs[i];
      cx2 = f.c[0] - ex; cy2 = f.c[1] - ey; cz2 = f.c[2] - ez;
      /* 背面剔除：法线指向相机（eye - 面心）才画 */
      var ex2 = ex - f.c[0], ey2 = ey - f.c[1], ez2 = ez - f.c[2];
      dot = f.n[0] * ex2 + f.n[1] * ey2 + f.n[2] * ez2;
      if (dot <= 0) continue;
      /* 距离剔除 */
      cx2 = f.c[0] - ex; cy2 = f.c[1] - ey; cz2 = f.c[2] - ez;
      vx = cx2 * fx + cy2 * fy + cz2 * fz;             /* 前方距离 */
      if (vx > (this.walk ? 26 : c.dist * 4.2) || vx < 0.02) continue;

      /* 逐顶点变换 + 近平面裁剪（vz = 沿 forward 的正向距离） */
      p = f.p; n = p.length;
      var sxA = [], syA = [], vzA = [], zA = [];
      var clip = false, anyNear = false;
      var emph = (f.cab && emphMap[f.cab]) || 0;
      var cabScl = f.cab ? (0.97 + 0.19 * emph) : 1;
      for (j = 0; j < n; j++) {
        var qv = p[j];
        var qx = qv[0], qy = qv[1], qz = qv[2];
        if (cabScl !== 1) {
          qx = f.ax + (qx - f.ax) * cabScl;
          qy = f.ay + (qy - f.ay) * cabScl;
        }
        vx = (qx - ex) * rx + (qy - ey) * ry + (qz - ez) * rz;
        vy = (qx - ex) * ux + (qy - ey) * uy + (qz - ez) * uz;
        vz = (qx - ex) * fx + (qy - ey) * fy + (qz - ez) * fz;
        if (vz < 0.06) { anyNear = true; }
        sxA.push(vx); syA.push(vy); vzA.push(vz);
      }
      /* 铭牌要贴在这块面上。留下裁剪前的相机坐标，后面按透视把字印上去。 */
      plateCam = null;
      if ((f.plate || f.skin) && n === 4 && !anyNear) {
        plateCam = [
          sxA[0], syA[0], vzA[0], sxA[1], syA[1], vzA[1],
          sxA[2], syA[2], vzA[2], sxA[3], syA[3], vzA[3]
        ];
      }
      if (anyNear) {
        /*  Sutherland–Hodgman 裁 vz >= 0.06 */
        var A1 = [], B1 = [], C1 = [], D1 = [];
        for (j = 0; j < n; j++) {
          var j2 = (j + 1) % n;
          var inA = vzA[j] >= 0.06, inB = vzA[j2] >= 0.06;
          if (inA) { A1.push(sxA[j]); B1.push(syA[j]); C1.push(vzA[j]); }
          if (inA !== inB) {
            var t2 = (0.06 - vzA[j]) / (vzA[j2] - vzA[j]);
            A1.push(sxA[j] + (sxA[j2] - sxA[j]) * t2);
            B1.push(syA[j] + (syA[j2] - syA[j]) * t2);
            C1.push(0.06);
          }
        }
        n = A1.length;
        if (n < 3) continue;
        sxA = A1; syA = B1; vzA = C1;
        clip = true;
      }
      var px = [], py = [], zsum = 0;
      for (j = 0; j < n; j++) {
        var inv = scale / vzA[j];
        px.push(this.cx + sxA[j] * inv);
        py.push(this.cy - syA[j] * inv);
        zsum += vzA[j];
      }
      dep = zsum / n;
      /* 屏幕面积剔除 */
      ar = 0;
      for (j = 0; j < n; j++) {
        var j3 = (j + 1) % n;
        ar += px[j] * py[j3] - px[j3] * py[j];
      }
      ar = Math.abs(ar) * 0.5;
      if (ar < 0.7) continue;

      /* 屏幕外剔除：整块落在画布外的面直接丢，省掉一次无谓的路径构建 */
      var gx0 = px[0], gx1 = px[0], gy0 = py[0], gy1 = py[0];
      for (j = 1; j < n; j++) {
        if (px[j] < gx0) gx0 = px[j]; else if (px[j] > gx1) gx1 = px[j];
        if (py[j] < gy0) gy0 = py[j]; else if (py[j] > gy1) gy1 = py[j];
      }
      if (gx1 < 0 || gx0 > W || gy1 < 0 || gy0 > H) continue;
      /* 巨面（地板 / 大墙）近平面裁剪后会被拉到十万像素级。填充本身是对的
         —— Canvas 光栅化会自己裁到画布边界 —— 但描边会变成横贯全屏的
         长线，而且这种路径极慢。所以巨面只填不描。 */
      var huge = (gx1 - gx0 > W * 6 || gy1 - gy0 > H * 6);

      /* 近侧墙淡出：面心比场景中心更靠近相机时，它就横在视线前面。
         不拆墙（房间感没了），而是把不透明度压到 0.06、描边也去掉，
         低角度看机柜就不会被齐腰的墙切掉一截。判定用面心到相机的
         距离与到场景中心的距离比较，与具体朝向无关，四个方向都成立。 */
      var faceOp = f.op == null ? 1 : f.op;
      var noStroke2 = huge || !!f.quiet;
      if (f.fade) {
        var dcx = f.c[0] - ex, dcy = f.c[1] - ey, dcz = f.c[2] - ez;
        var dFace = dcx * dcx + dcy * dcy + dcz * dcz;
        var tcx = c.tx - ex, tcy = c.ty - ey, tcz = c.tz - ez;
        var dTgt = tcx * tcx + tcy * tcy + tcz * tcz;
        if (dFace < dTgt * 0.80) { faceOp = Math.min(faceOp, 0.06); noStroke2 = true; }
      }

      var tt = (dep - fogN) / (fogF - fogN);
      var kLit = f.k * (f.cab ? (0.7 + 0.46 * emph) : 1);
      var col = f.noFog ? shade(f.fill, kLit) : fog(f.fill, kLit, tt);
      if (f.dyn) {
        var rgb = f.dyn(time, (i % 7) * 0.9);
        if (rgb) col = 'rgba(' + rgb[0] + ',' + rgb[1] + ',' + rgb[2] + ',' + faceOp + ')';
      }
      /* pick 存逐顶点深度（vzA 已在上面算好）：绘制多边形被简化成 3~4 点用于
         填充，但拾取必须用真实形状 + 真实深度，否则俯视角下"近处柜顶"会盖住
         "远处柜正面"，鼠标指着远处柜却命中近处柜 */
      this.q.push({ k: 0, lay: f.lay, op: faceOp, noStroke: noStroke2,
        poly: px.length === 4 ? [px[0], py[0], px[1], py[1], px[2], py[2], px[3], py[3]] : [px[0], py[0], px[1], py[1], px[2], py[2]],
        n: n, dep: dep, col: col, sw: f.sw, st: f.stroke, blend: f.blend, obj: f.obj, plate: f.plate, skin: f.skin, ar: ar,
        emph: emph, frame: f.frame || 0,
        hover: f.obj != null && this.hover === f.obj,
        sel: f.obj != null && this.sel === f.obj,
        /* pickPri：设备面板只有 80×14 px，而柜体侧板/后板这类结构件
           在屏幕上和它重叠。深度几乎相同时（差 < 0.12m，恰好是这种
           "面挨着面"的贴合关系），让设备赢 —— 否则用户明明点在设备
           面板上，命中的却是背后的柜体，观感就是"鼠标和选中的东西对不上" */
        pickPri: f.obj ? (String(f.obj).indexOf('dev:') === 0 ? 1 : 0) : 0,
        pickPoly: this.pickPts(px, py, vzA), pickDep: dep, plateCam: plateCam });
    }

    /* ---- 线 ---- */
    for (i = 0; i < bxl.length; i++) {
      var ln = bxl[i];
      var m0 = [ex - ln.a[0], ey - ln.a[1], ez - ln.a[2]];
      dot = ln.n[0] * m0[0] + ln.n[1] * m0[1] + ln.n[2] * m0[2];
      if (dot <= 0 && !ln.obj) continue;
      var a1 = this.proj(ln.a, ex, ey, ez, rx, ry, rz, ux, uy, uz, fx, fy, fz, scale);
      var b1 = this.proj(ln.b, ex, ey, ez, rx, ry, rz, ux, uy, uz, fx, fy, fz, scale);
      if (!a1 || !b1) continue;
      dep = (a1[2] + b1[2]) / 2;
      var tt2 = (dep - fogN) / (fogF - fogN);
      this.q.push({ k: 1, x1: a1[0], y1: a1[1], x2: b1[0], y2: b1[1], dep: dep,
        col: ln.noFog ? ln.c : fogLine(ln.c, tt2), op: ln.op, w: ln.w, dash: ln.dash, blend: ln.blend });
    }

    /* ---- 标签 ----
       1) dim 标签（机柜编号，110 个）只在对应对象被 hover / 选中时才绘制。
          它们是覆盖层（lay 2），一旦画出来就永远压在机柜之上 —— 实测一个
          "2-20" 编号能盖住 6 台机柜的 660 个像素，点了那些位置只会选中
          这个编号代表的柜子，旁边的柜子反而点不到。索性不 hover 就不画：
          既不糊画面，也一像素都不占。视线扫到哪台，编号就浮现在哪台上。
       2) 其余可点标签把屏幕包围盒压进 pickLbl，供 hit() 优先判定。 */
    for (i = 0; i < bxt.length; i++) {
      var lb = bxt[i];
      var pp = this.proj(lb.p, ex, ey, ez, rx, ry, rz, ux, uy, uz, fx, fy, fz, scale);
      if (!pp) continue;
      var lbHov = lb.obj != null && this.hover === lb.obj;
      var lbSel = lb.obj != null && this.sel === lb.obj;
      /* 淡显标签：没被指着就完全不画 */
      if (lb.dim && !lbHov && !lbSel) continue;
      var tt3 = (pp[2] - fogN) / (fogF - fogN);
      this.q.push({ k: 2, lay: lb.lay, x: pp[0], y: pp[1], dep: pp[2], s: lb.s, size: lb.size, col: lb.col,
        align: lb.align, dy: lb.dy, weight: lb.weight, family: lb.family, bg: lb.bg, mono: lb.mono,
        glow: lb.glow, fogT: tt3, obj: lb.obj, dim: lb.dim,
        hov: lbHov, sel: lbSel });
      /* 标签是屏幕空间文字，不在多边形拾取表里 —— 射线会直接穿过字，
         打到字背后的地板/邻房，于是"点 209 却选中 207"。
         这里把每个字的屏幕包围盒记下来，点击时先判标签。 */
      if (lb.obj) {
        ctx.save();
        ctx.font = (lb.weight || '400') + ' ' + lb.size + 'px ' + (lb.mono ? MONO : (lb.family || FONT));
        var tw = ctx.measureText(lb.s).width;
        ctx.restore();
        /* 包围盒必须与 text() 的实际绘制对齐方式一致（同一套 ALIGN 语义），
           否则"看得见的字"和"判定的区域"会错开半个字宽 */
        var al = alignOf(lb.align);
        var bx = al === 'left' ? pp[0] : al === 'right' ? pp[0] - tw : pp[0] - tw / 2;
        var by = pp[1] + (lb.dy || 0);
        lbls.push({ x0: bx - 1, y0: by - lb.size * 0.85, x1: bx + tw + 1, y1: by + lb.size * 0.28,
          obj: lb.obj, dep: pp[2] });
      }
    }

    /* ---- 气流粒子 ---- */
    var parts = this.sc.parts;
    for (i = 0; i < parts.length; i++) {
      var pt = parts[i];
      for (j = 0; j < pt.n; j++) {
        var seed = pt.seed * 13.7 + j * 5.31;
        var ph = (time * (pt.cold ? 0.10 : 0.17) + (seed % 1)) % 1;
        var pxp = pt.cold ? lerp(pt.x0, pt.x1, ((j / pt.n) + ph) % 1) : lerp(pt.x0, pt.x1, (j % 3) / 2);
        var pyp = pt.cold ? lerp(pt.y0, pt.y1, (j * 7 % pt.n) / pt.n) : lerp(pt.y0, pt.y1, (j % 5) / 4);
        var pzp = pt.cold ? 0.06 : 0.05 + ph * 2.3;
        var pa = pt.cold ? (0.85 - ph * 0.85) : (1 - ph);
        if (pa <= 0.03) continue;
        var pj = this.proj([pxp, pyp, pzp], ex, ey, ez, rx, ry, rz, ux, uy, uz, fx, fy, fz, scale);
        if (!pj) continue;
        var rr = clamp(scale / pj[2] * 0.02, 0.7, 4.2);
        this.q.push({ k: 3, x: pj[0], y: pj[1], dep: pj[2] - 0.001, r: rr,
          col: pt.cold ? 'rgba(0,216,255,' + pa.toFixed(2) + ')' : 'rgba(255,140,80,' + pa.toFixed(2) + ')' });
      }
    }

    /* ---- 排序：先按渲染层（0 地面背景 → 1 常规 → 2 覆盖），层内远先画。
       半透明面（通道/剪影）与不透明面混排才是正确的画家算法；
       提前到尾会洗淡前景，拆两趟会让半透明盖住不透明 ---- */
    q.sort(function (a, b) {
      var la = a.lay || 0, lb = b.lay || 0;
      if (la !== lb) return la - lb;
      return b.dep - a.dep;
    });

    /* ---- 主体：面 / 线 / 标签 一次遍历 ---- */
    ctx.lineJoin = 'round';
    for (i = 0; i < q.length; i++) {
      var d = q[i];
      if (d.k === 0) this.poly(ctx, d);
      else if (d.k === 1) this.seg(ctx, d);
      else if (d.k === 2) this.text(ctx, d);
    }
    /* ---- 气流粒子：加色合成，像真的在发光，压在最终画面之上 ---- */
    ctx.save();
    ctx.globalCompositeOperation = 'lighter';
    for (i = 0; i < q.length; i++) {
      var d3 = q[i];
      if (d3.k !== 3) continue;
      this.dot(ctx, d3);
    }
    ctx.restore();

    this.paintCabNames(ctx, ex, ey, ez, rx, ry, rz, ux, uy, uz, fx, fy, fz, scale);
    this.minimap();
  };

  Viewer.prototype.proj = function (p, ex, ey, ez, rx, ry, rz, ux, uy, uz, fx, fy, fz, scale) {
    var a = p[0] - ex, b = p[1] - ey, c = p[2] - ez;
    var vz = a * fx + b * fy + c * fz;
    if (vz < 0.06) return null;
    var vx = a * rx + b * ry + c * rz;
    var vy = a * ux + b * uy + c * uz;
    var inv = scale / vz;
    return [this.cx + vx * inv, this.cy - vy * inv, vz];
  };

  Viewer.prototype.poly = function (ctx, d) {
    var n = d.n;
    ctx.beginPath();
    ctx.moveTo(d.poly[0], d.poly[1]);
    if (n === 4) { ctx.lineTo(d.poly[2], d.poly[3]); ctx.lineTo(d.poly[4], d.poly[5]); ctx.lineTo(d.poly[6], d.poly[7]); }
    else { ctx.lineTo(d.poly[2], d.poly[3]); ctx.lineTo(d.poly[4], d.poly[5]); }
    ctx.closePath();
    if (d.frame) {
      var em = d.emph || 0;
      ctx.globalAlpha = 0.18 + 0.82 * em;
      ctx.strokeStyle = em > 0.45 ? '#FFF4CC' : '#FFFFFF';
      ctx.lineWidth = 1.1 + 2.1 * em;
      ctx.stroke();
      ctx.globalAlpha = 1;
      return;
    }
    ctx.globalAlpha = d.op == null ? 1 : d.op;
    ctx.fillStyle = d.col;
    ctx.fill();
    if (d.skin && d.plateCam) this.drawSkin(ctx, d);
    /* 贴图会改掉当前路径。不重画的话，悬停描边会沿贴图用的三角形走，门上就多出一个白三角。 */
    ctx.beginPath();
    ctx.moveTo(d.poly[0], d.poly[1]);
    if (n === 4) { ctx.lineTo(d.poly[2], d.poly[3]); ctx.lineTo(d.poly[4], d.poly[5]); ctx.lineTo(d.poly[6], d.poly[7]); }
    else { ctx.lineTo(d.poly[2], d.poly[3]); ctx.lineTo(d.poly[4], d.poly[5]); }
    ctx.closePath();
    /* 悬停 / 选中：先叠一层提亮再描边。
       柜体底色有青、黄、红三种，1~2px 的细描边落在同色系柜顶上很难辨认，
       用户盯着看也确认不了"指着的到底是不是这台"。提亮整面 + 描边双管齐下，
       视线扫过来就能定位，且不改变实际配色语义。 */
    if (d.hover || d.sel) {
      ctx.globalAlpha = 1;
      /* 提亮只给 10%：柜体颜色承载四级色阶（灰/蓝/黄/红），
         洗白了就看不出负载等级 —— 主要靠下面的粗描边来抓眼。 */
      ctx.fillStyle = d.sel ? 'rgba(255,209,102,.10)' : 'rgba(226,246,255,.10)';
      ctx.fill();
    }
    if (d.noStroke) {
      /* 巨面只填不描 */
    } else if (d.hover || d.sel) {
      ctx.globalAlpha = 1;
      /* 悬停描边原来用 #7FE9FF + 1.1px，落在青色柜顶/黄色柜面上几乎看不见 ——
         用户没法确认"我指着的到底是不是这台"。改成近白 + 加粗：
         在青、黄、红三种柜体底色上都能拉开对比，选中用暖黄更能一眼分开。 */
      ctx.strokeStyle = d.sel ? '#FFD166' : '#FFFFFF';
      ctx.lineWidth = d.sel ? 2.6 : 2.3;
      ctx.stroke();
    } else if (d.st) {
      ctx.globalAlpha = Math.min(1, (d.op == null ? 1 : d.op) + 0.15);
      ctx.strokeStyle = d.st;
      ctx.lineWidth = d.sw || 1;
      ctx.stroke();
    }
    ctx.globalAlpha = 1;
    if (d.plate) this.drawPlate(ctx, d);
    if (d.obj) {
      var pp = d.pickPoly || this.polyPts(d);
      this.pick.push({ poly: pp, obj: d.obj, dep: d.pickDep == null ? d.dep : d.pickDep, pri: d.pickPri || 0 });
    }
  };
  function ellipsize(ctx, text, maxW) {
    text = String(text || '');
    if (!text) return '';
    if (ctx.measureText(text).width <= maxW) return text;
    var s = text;
    while (s.length > 1 && ctx.measureText(s + '…').width > maxW) s = s.slice(0, -1);
    return s + '…';
  }
  Viewer.prototype.drawSkin = function (ctx, d) {
    var fn = this.sc && this.sc.shellSkin;
    var img = fn ? fn(d.skin) : null;
    if (!img || !d.plateCam) return;
    this.paintOnPlate(ctx, img, d.plateCam);
  };
  /* 铭牌先画成一张平面图，再按这块面板的四个角做透视，字跟面板是一体的。 */
  Viewer.prototype.drawPlate = function (ctx, d) {
    var lines = d.plate && d.plate.lines;
    if (!lines || !lines.length || !d.plateCam) return;
    var img = this.plateStamp(d.plate);
    if (!img) return;
    this.paintOnPlate(ctx, img, d.plateCam);
  };
  Viewer.prototype.plateStamp = function (plate) {
    var lines = plate.lines || [];
    var key = lines.map(function (ln) { return (ln.t || '') + '|' + (ln.col || ''); }).join('\n');
    if (!this._plateStamp) this._plateStamp = {};
    if (this._plateStamp[key]) return this._plateStamp[key];
    var cnv = document.createElement('canvas');
    var W = 640, n = Math.min(lines.length, 3);
    var H = n >= 3 ? 220 : n === 2 ? 168 : 112;
    cnv.width = W; cnv.height = H;
    var g = cnv.getContext('2d');
    var sizes = n >= 3 ? [54, 38, 32] : n === 2 ? [58, 40] : [64];
    var weights = ['600', '500', '400'];
    var block = 0, i, texts = [];
    var maxW = W - 56;
    for (i = 0; i < n; i++) {
      g.font = (weights[i] || '500') + ' ' + sizes[i] + 'px ' + FONT;
      texts.push(ellipsize(g, lines[i].t || '', maxW));
      block += sizes[i];
    }
    block += (n - 1) * 10;
    var y = (H - block) / 2;
    g.textAlign = 'center';
    g.textBaseline = 'top';
    for (i = 0; i < n; i++) {
      g.font = (weights[i] || '500') + ' ' + sizes[i] + 'px ' + FONT;
      g.fillStyle = lines[i].col || '#EAF6FF';
      g.fillText(texts[i], W / 2, y);
      y += sizes[i] + 10;
    }
    this._plateStamp[key] = cnv;
    return cnv;
  };
  Viewer.prototype.paintOnPlate = function (ctx, img, cam) {
    var self = this;
    function at(u, v) {
      var topx = cam[3] + (cam[6] - cam[3]) * u;
      var topy = cam[4] + (cam[7] - cam[4]) * u;
      var topz = cam[5] + (cam[8] - cam[5]) * u;
      var botx = cam[0] + (cam[9] - cam[0]) * u;
      var boty = cam[1] + (cam[10] - cam[1]) * u;
      var botz = cam[2] + (cam[11] - cam[2]) * u;
      var x = topx + (botx - topx) * v;
      var y = topy + (boty - topy) * v;
      var z = topz + (botz - topz) * v;
      if (z < 0.06) return null;
      var inv = (self._scrScale || self.scale) / z;
      return [self.cx + x * inv, self.cy - y * inv];
    }
    var p00 = at(0, 0), p10 = at(1, 0), p01 = at(0, 1), p11 = at(1, 1);
    if (!p00 || !p10 || !p01 || !p11) return;
    var spanW = Math.hypot(p10[0] - p00[0], p10[1] - p00[1]);
    var spanH = Math.hypot(p01[0] - p00[0], p01[1] - p00[1]);
    if (spanW < 16 || spanH < 8) return;
    function bands(za, zb) {
      var r = Math.max(za, zb) / Math.max(0.06, Math.min(za, zb));
      if (r < 1.05) return 1;
      return Math.min(6, Math.ceil((r - 1) * 16));
    }
    var cols = Math.max(2, bands((cam[2] + cam[5]) * 0.5, (cam[11] + cam[8]) * 0.5));
    var rows = Math.max(3, bands((cam[5] + cam[8]) * 0.5, (cam[2] + cam[11]) * 0.5));
    var iw = img.width, ih = img.height;
    function tri(s0, s1, s2, d0, d1, d2) {
      var x0 = s0[0], y0 = s0[1], x1 = s1[0], y1 = s1[1], x2 = s2[0], y2 = s2[1];
      var u0 = d0[0], v0 = d0[1], u1 = d1[0], v1 = d1[1], u2 = d2[0], v2 = d2[1];
      var det = x0 * (y1 - y2) + x1 * (y2 - y0) + x2 * (y0 - y1);
      if (Math.abs(det) < 1e-4) return;
      var a = (u0 * (y1 - y2) + u1 * (y2 - y0) + u2 * (y0 - y1)) / det;
      var c = (u0 * (x2 - x1) + u1 * (x0 - x2) + u2 * (x1 - x0)) / det;
      var e = (u0 * (x1 * y2 - x2 * y1) + u1 * (x2 * y0 - x0 * y2) + u2 * (x0 * y1 - x1 * y0)) / det;
      var b = (v0 * (y1 - y2) + v1 * (y2 - y0) + v2 * (y0 - y1)) / det;
      var d = (v0 * (x2 - x1) + v1 * (x0 - x2) + v2 * (x1 - x0)) / det;
      var f = (v0 * (x1 * y2 - x2 * y1) + v1 * (x2 * y0 - x0 * y2) + v2 * (x0 * y1 - x1 * y0)) / det;
      ctx.save();
      ctx.beginPath();
      ctx.moveTo(d0[0], d0[1]); ctx.lineTo(d1[0], d1[1]); ctx.lineTo(d2[0], d2[1]);
      ctx.closePath();
      ctx.clip();
      var k = self.dpr || 1;
      ctx.setTransform(a * k, b * k, c * k, d * k, e * k, f * k);
      ctx.drawImage(img, 0, 0, iw, ih);
      ctx.restore();
    }
    var r, c, u0, u1, v0, v1, q00, q10, q01, q11;
    ctx.save();
    ctx.beginPath();
    ctx.moveTo(p00[0], p00[1]); ctx.lineTo(p10[0], p10[1]); ctx.lineTo(p11[0], p11[1]); ctx.lineTo(p01[0], p01[1]);
    ctx.closePath();
    ctx.clip();
    for (r = 0; r < rows; r++) {
      v0 = r / rows; v1 = (r + 1) / rows;
      for (c = 0; c < cols; c++) {
        u0 = c / cols; u1 = (c + 1) / cols;
        q00 = at(u0, v0); q10 = at(u1, v0); q01 = at(u0, v1); q11 = at(u1, v1);
        if (!q00 || !q10 || !q01 || !q11) continue;
        /* 面的 +x 投到屏幕上是反的，图片左右对调后铭牌和 U 位才正读。 */
        var x0 = (1 - u0) * iw, x1 = (1 - u1) * iw;
        tri([x0, v0 * ih], [x1, v0 * ih], [x1, v1 * ih], q00, q10, q11);
        tri([x0, v0 * ih], [x1, v1 * ih], [x0, v1 * ih], q00, q11, q01);
      }
    }
    ctx.restore();
  };
  /* 拾取用多边形：真实顶点 + 每点深度 */
  Viewer.prototype.pickPts = function (px, py, vz) {
    var o = [], i;
    for (i = 0; i < px.length; i++) o.push([px[i], py[i], vz[i]]);
    return o;
  };
  Viewer.prototype.polyPts = function (d) {
    var o = [], i;
    for (i = 0; i < d.poly.length; i += 2) o.push([d.poly[i], d.poly[i + 1]]);
    return o;
  };
  Viewer.prototype.seg = function (ctx, d) {
    ctx.globalAlpha = d.op == null ? 1 : d.op;
    ctx.strokeStyle = d.col;
    ctx.lineWidth = d.w || 1;
    if (d.dash) ctx.setLineDash(d.dash.split(' ')); else ctx.setLineDash([]);
    ctx.beginPath(); ctx.moveTo(d.x1, d.y1); ctx.lineTo(d.x2, d.y2); ctx.stroke();
    ctx.setLineDash([]);
    ctx.globalAlpha = 1;
  };
  Viewer.prototype.dot = function (ctx, d) {
    ctx.globalAlpha = 1;
    ctx.fillStyle = d.col;
    ctx.beginPath(); ctx.arc(d.x, d.y, d.r, 0, 6.2832); ctx.fill();
    ctx.globalAlpha = 1;
  };
  /* 文本对齐：Canvas 的 textAlign 合法枚举只有 left / right / center / start / end。
     写 'middle' 会被浏览器静默忽略（控制台报 "not a valid enum value"），
     文字实际按上一次残留的值绘制 —— 默认就是左对齐。而拾取包围盒按
     "居中"算，于是标签的实际位置和判定区域差了半个字宽：
     鼠标指着字，命中的是字右边那个对象 —— 典型症状就是"点到隔壁"。
     统一把 'middle' 归一成 'center'，绘制与包围盒用同一套语义。 */
  var ALIGN = { middle: 'center', center: 'center', start: 'left', left: 'left', end: 'right', right: 'right' };
  function alignOf(a) { return ALIGN[a] || 'center'; }

  Viewer.prototype.text = function (ctx, d) {
    var size = d.size;
    ctx.font = (d.weight || '400') + ' ' + size + 'px ' + (d.mono ? MONO : (d.family || FONT));
    ctx.textAlign = alignOf(d.align);
    ctx.textBaseline = 'alphabetic';
    var y = d.y + (d.dy || 0);
    /* 淡显标签（dim）：密集场景里的机柜编号这类。默认压暗，
       hover / 选中对应对象时恢复全亮 —— 视线扫到哪台，哪台的编号就亮起来。 */
    var dimA = d.dim && !d.hov && !d.sel ? 0.26 : 1;
    if (d.bg) {
      var w = ctx.measureText(d.s).width;
      var bx = d.align === 'start' || d.align === 'left' ? d.x - 3 :
        d.align === 'end' || d.align === 'right' ? d.x - w - 3 : d.x - w / 2 - 3;
      ctx.globalAlpha = 0.82 * dimA;
      ctx.fillStyle = d.bg.indexOf('rgba') === 0 ? d.bg : '#030C18';
      ctx.fillRect(bx, y - size, w + 6, size + 4);
      ctx.globalAlpha = 1;
    }
    ctx.lineWidth = 2.6;
    ctx.strokeStyle = 'rgba(3,12,24,.9)';
    ctx.lineJoin = 'round';
    ctx.globalAlpha = dimA;
    ctx.strokeText(d.s, d.x, y);
    /* 悬停/选中时文字提亮：标签现在参与拾取，必须给出"这一行字可点"的反馈 */
    ctx.fillStyle = d.sel ? '#FFD166' : d.hov ? '#7FE9FF' : d.col;
    ctx.fillText(d.s, d.x, y);
    ctx.globalAlpha = 1;
  };

  /* 柜名贴在该柜顶前缘上，字的走向跟这条柜顶边在画面里的走向一致。
     只画镜头前方、顶边还落在画面里的近处柜子，左右各 3、最多 6；叠住丢掉更远的。 */
  Viewer.prototype.paintCabNames = function (ctx, ex, ey, ez, rx, ry, rz, ux, uy, uz, fx, fy, fz, scale) {
    if (!this.walk || !this.sc.walk || !this.sc.walk.cabs) return;
    var W = this.W, H = this.H;
    var cabs = this.sc.walk.cabs;
    var emphMap = this._emph || {};
    var sel = this.selCab || this.key;
    var items = [], i, c, dx, dy, ahead, lat, half, pL, pR, edge, ang, mx, my;
    var zTop = 2.0;
    for (i = 0; i < cabs.length; i++) {
      c = cabs[i];
      dx = c.fx - ex; dy = c.fy - ey;
      ahead = dx * fx + dy * fy;
      lat = dx * rx + dy * ry;
      if (ahead < 0.35) continue;
      if (ahead > 5.4 && c.name !== sel) continue;
      if (ahead > 7.2) continue;
      half = c.halfW || 0.32;
      pL = this.proj([c.fx - c.rx * half, c.fy - c.ry * half, zTop], ex, ey, ez, rx, ry, rz, ux, uy, uz, fx, fy, fz, scale);
      pR = this.proj([c.fx + c.rx * half, c.fy + c.ry * half, zTop], ex, ey, ez, rx, ry, rz, ux, uy, uz, fx, fy, fz, scale);
      if (!pL || !pR) continue;
      edge = Math.hypot(pR[0] - pL[0], pR[1] - pL[1]);
      if (edge < 28) continue;
      if (edge < 42 && c.name !== sel) continue;
      mx = (pL[0] + pR[0]) / 2;
      my = (pL[1] + pR[1]) / 2;
      if (mx < 8 || mx > W - 8 || my < 10 || my > H - 10) continue;
      ang = Math.atan2(pR[1] - pL[1], pR[0] - pL[0]);
      if (ang > Math.PI / 2) ang -= Math.PI;
      else if (ang < -Math.PI / 2) ang += Math.PI;
      items.push({
        name: c.name, title: c.title || c.name,
        x: mx, y: my, ang: ang, edge: edge, ahead: ahead, lat: lat,
        emph: emphMap[c.name] || 0, sel: c.name === sel
      });
    }
    var left = [], right = [];
    for (i = 0; i < items.length; i++) (items[i].lat < 0 ? left : right).push(items[i]);
    left.sort(function (a, b) { return a.ahead - b.ahead; });
    right.sort(function (a, b) { return a.ahead - b.ahead; });
    var picked = left.slice(0, 3).concat(right.slice(0, 3));
    var have = {};
    for (i = 0; i < picked.length; i++) have[picked[i].name] = 1;
    for (i = 0; i < items.length; i++) {
      if (items[i].sel && !have[items[i].name]) picked.push(items[i]);
    }
    picked.sort(function (a, b) {
      if (a.sel !== b.sel) return a.sel ? -1 : 1;
      return a.ahead - b.ahead;
    });
    if (picked.length > 6) {
      var keep = [], dropped = 0;
      for (i = 0; i < picked.length; i++) {
        if (picked[i].sel || keep.length < 6) keep.push(picked[i]);
        else dropped++;
      }
      while (keep.length > 6) {
        var cut = -1, far = -1, k;
        for (k = 0; k < keep.length; k++) {
          if (keep[k].sel) continue;
          if (keep[k].ahead > far) { far = keep[k].ahead; cut = k; }
        }
        if (cut < 0) break;
        keep.splice(cut, 1);
      }
      picked = keep;
    }
    function overlaps(box, list) {
      var k, o;
      for (k = 0; k < list.length; k++) {
        o = list[k];
        if (box.x0 < o.x1 && box.x1 > o.x0 && box.y0 < o.y1 && box.y1 > o.y0) return true;
      }
      return false;
    }
    function mixByte(a, b, t) { return Math.round(a + (b - a) * t); }
    var placed = [];
    ctx.save();
    ctx.textAlign = 'center';
    ctx.textBaseline = 'bottom';
    for (i = 0; i < picked.length; i++) {
      var it = picked[i];
      var em = it.sel ? 1 : 0;
      var sz = clamp(it.edge / 7.2, 11, 16);
      if (em) sz = Math.min(18, sz + 2);
      var text = it.title;
      ctx.font = '700 ' + sz.toFixed(1) + 'px ' + FONT;
      var guard = 0, tw;
      while (guard++ < 6 && sz > 11 && ctx.measureText(text).width > it.edge * 0.92) {
        sz -= 1;
        ctx.font = '700 ' + sz.toFixed(1) + 'px ' + FONT;
      }
      tw = ctx.measureText(text).width;
      if (tw > it.edge * 1.08) continue;
      var cth = Math.cos(it.ang), sth = Math.sin(it.ang);
      var hw = tw / 2, hh = sz;
      function corner(lx, ly) {
        return [it.x + lx * cth - ly * sth, it.y + lx * sth + ly * cth];
      }
      var cs = [corner(-hw, -hh), corner(hw, -hh), corner(hw, 0), corner(-hw, 0)];
      var x0 = 1e9, y0 = 1e9, x1 = -1e9, y1 = -1e9, ci, pt;
      for (ci = 0; ci < cs.length; ci++) {
        pt = cs[ci];
        if (pt[0] < x0) x0 = pt[0];
        if (pt[1] < y0) y0 = pt[1];
        if (pt[0] > x1) x1 = pt[0];
        if (pt[1] > y1) y1 = pt[1];
      }
      if (x1 < 4 || x0 > W - 4 || y1 < 4 || y0 > H - 4) continue;
      var box = { x0: x0 - 2, y0: y0 - 2, x1: x1 + 2, y1: y1 + 2 };
      if (!it.sel && overlaps(box, placed)) continue;
      placed.push(box);
      var col = 'rgb(' + mixByte(244, 255, em) + ',' + mixByte(248, 209, em) + ',' + mixByte(255, 90, em) + ')';
      ctx.save();
      ctx.translate(it.x, it.y);
      ctx.rotate(it.ang);
      ctx.font = '700 ' + sz.toFixed(1) + 'px ' + FONT;
      if (em) {
        ctx.shadowColor = 'rgba(255, 196, 48, .9)';
        ctx.shadowBlur = 16;
      }
      ctx.lineWidth = 3;
      ctx.strokeStyle = em ? 'rgba(90, 48, 0, .8)' : 'rgba(4, 10, 18, .88)';
      ctx.strokeText(text, 0, -1);
      ctx.fillStyle = col;
      ctx.fillText(text, 0, -1);
      ctx.restore();
    }
    ctx.restore();
  };
  /* ---------- 小地图（右下角俯视 + 相机扇形） ---------- */
  Viewer.prototype.minimap = function () {
    var ctx = this.ctx, W = this.W, H = this.H;
    var mw = 148, mh = 108, mx = W - mw - 14, my = H - mh - 14;
    var bb = this.sc.bb;
    var bw = bb.x1 - bb.x0, bh = bb.y1 - bb.y0;
    var sc2 = Math.min((mw - 14) / bw, (mh - 14) / bh);
    var ox = mx + mw / 2 - (bb.x0 + bw / 2) * sc2;
    var oy = my + mh / 2 - (bb.y0 + bh / 2) * sc2;

    ctx.save();
    ctx.globalAlpha = 0.9;
    ctx.fillStyle = 'rgba(4,12,22,.78)';
    ctx.strokeStyle = 'rgba(45,90,130,.55)';
    ctx.lineWidth = 1;
    ctx.beginPath();
    if (ctx.roundRect) ctx.roundRect(mx, my, mw, mh, 6); else ctx.rect(mx, my, mw, mh);
    ctx.fill(); ctx.stroke();

    /* 俯视轮廓 */
    ctx.globalAlpha = 0.5;
    ctx.strokeStyle = 'rgba(80,150,200,.7)';
    ctx.beginPath();
    ctx.rect(ox + bb.x0 * sc2, oy + bb.y0 * sc2, bw * sc2, bh * sc2);
    ctx.stroke();
    ctx.globalAlpha = 1;

    /* 相机位置 + 朝向扇形 */
    var c = this.cur;
    var cxp = ox + (c.tx + c.dist * Math.sin(c.yaw) * Math.cos(c.pitch)) * sc2;
    var cyp = oy + (c.ty + c.dist * Math.cos(c.yaw) * Math.cos(c.pitch)) * sc2;
    var txp = ox + c.tx * sc2, typ = oy + c.ty * sc2;
    var ang = Math.atan2(typ - cyp, txp - cxp);
    ctx.fillStyle = 'rgba(0,216,255,.18)';
    ctx.beginPath();
    ctx.moveTo(cxp, cyp);
    ctx.arc(cxp, cyp, 26, ang - 0.42, ang + 0.42);
    ctx.closePath(); ctx.fill();
    ctx.fillStyle = '#00D8FF';
    ctx.beginPath(); ctx.arc(cxp, cyp, 2.6, 0, 6.2832); ctx.fill();
    ctx.strokeStyle = 'rgba(0,216,255,.85)';
    ctx.beginPath(); ctx.arc(txp, typ, 3, 0, 6.2832); ctx.stroke();

    ctx.globalAlpha = 0.75;
    ctx.fillStyle = '#7FA8C8';
    ctx.font = '9px ' + MONO;
    ctx.textAlign = 'left'; ctx.textBaseline = 'alphabetic';
    ctx.fillText('PLAN VIEW', mx + 8, my + 13);
    ctx.restore();
  };

  /* ---------- 销毁 ---------- */
  Viewer.prototype.dispose = function () {
    if (this.dead) return;
    this.dead = true;
    cancelAnimationFrame(this.raf);
    var cv = this.cv;
    cv.removeEventListener('mousedown', this._md);
    window.removeEventListener('mousemove', this._mm);
    window.removeEventListener('mouseup', this._mu);
    cv.removeEventListener('wheel', this._wc);
    cv.removeEventListener('contextmenu', this._cc);
    cv.removeEventListener('dblclick', this._db);
    cv.removeEventListener('mouseleave', this._ml);
    cv.removeEventListener('touchstart', this._ts);
    cv.removeEventListener('touchmove', this._tm);
    cv.removeEventListener('touchend', this._te);
    window.removeEventListener('keydown', this._kd);
    if (this._ku) window.removeEventListener('keyup', this._ku);
    if (this._blur) window.removeEventListener('blur', this._blur);
    if (this.ro) this.ro.disconnect();
    this.host.innerHTML = '';
    if (active === this) active = null;
  };

  function fogLine(hex, t) {
    if (t <= 0.02) return hex;
    var q = Math.round(clamp(t, 0, 1) * 16) / 16;
    var key = 'L' + hex + '|' + q;
    var c = _fCache[key];
    if (c) return c;
    var p = parseCol(hex);
    var r = Math.round(lerp(p[0], BG[0], q)), g = Math.round(lerp(p[1], BG[1], q)), b = Math.round(lerp(p[2], BG[2], q));
    c = 'rgb(' + r + ',' + g + ',' + b + ')';
    _fCache[key] = c;
    return c;
  }

  /* ==================== 对外 API ==================== */
  A.V3 = {
    active: null,
    mount: function (host, kind, key, api) {
      if (active) active.dispose();
      var inherit = null;
      if (active && active.cur) inherit = { yaw: active.cur.yaw, pitch: active.cur.pitch };
      var v = new Viewer(host, kind, key, api, inherit);
      if (v.dead) { v.dispose(); return null; }
      active = v;
      A.V3.active = v;
      return v;
    },
    dispose: function () {
      if (active) { active.dispose(); active = null; A.V3.active = null; }
    }
  };
  window.ISO = window.ISO || {};
})(window.App);
