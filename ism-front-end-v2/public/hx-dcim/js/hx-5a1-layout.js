/* 配电室5A1 显式一次图：拓扑/坐标全部来自 data/hx-5a1-cad-layout.json（按 CAD 网表生成），
 * 不经 hx-screen 的通用并联/补线启发式。刀的点位绑定（k/p/dev/op/val/st）按名称取自 hx-scenes.json 原节点。 */
(function (w) {
  'use strict';

  var PAGE = '配电室5A1';
  var PREFIX = '配电室5A1_';
  var KNIFE_W = 60;
  var KNIFE_H = 98;
  var RED = '#ce0c18';
  var layout = null;

  function load(ver) {
    if (layout) return Promise.resolve(layout);
    return fetch('/hx-dcim/data/hx-5a1-cad-layout.json?v=' + ver)
      .then(function (r) { if (!r.ok) throw new Error('5A1 layout HTTP ' + r.status); return r.json(); })
      .then(function (j) { layout = j; return j; });
  }

  function setLayout(j) { layout = j; }

  function applies(pg) {
    return !!(layout && pg && pg.name === PAGE && pg.kind === 'power');
  }

  function normKey(k) {
    return String(k || '').replace(PREFIX, '').replace(/_+$/, '');
  }

  function cadKey(cad) {
    return String(cad).replace(/-/g, '_');
  }

  /* 同名多帧（分/合两帧）时优先 st='1'，与 hx-draw knifePickId 去重口径一致 */
  function indexSceneKnives(nodes) {
    var map = {};
    nodes.forEach(function (n) {
      if (n.t !== 'I' && n.t !== 'S') return;
      if (!n.k || String(n.k).indexOf(PREFIX) !== 0) return;
      var id = normKey(n.k);
      if (!map[id] || n.st === '1') map[id] = n;
    });
    return map;
  }

  function knifeNode(spec, src) {
    var n = {};
    if (src) Object.keys(src).forEach(function (key) { n[key] = src[key]; });
    delete n.rot;
    n.t = 'I';
    n.x = spec.x;
    n.y = spec.y;
    n.w = KNIFE_W;
    n.h = KNIFE_H;
    n.z = 5;
    n.lay = 1;
    n.cad = spec.cad;
    /* 5A1 全部竖刀。横置 orient 是已废弃的 pargroup 画法。 */
    delete n.orient;
    delete n.hfull;
    delete n.rot;
    if (spec.out) n.decor = 1; else delete n.decor;
    if (spec.outWire) n.outWire = 1; else delete n.outWire;
    return n;
  }

  function lineNode(a, b, wire) {
    var horiz = Math.abs(a[1] - b[1]) < 0.01;
    var len = horiz ? Math.abs(b[0] - a[0]) : Math.abs(b[1] - a[1]);
    if (len < 0.5) return null;
    var cx = (a[0] + b[0]) / 2;
    var cy = (a[1] + b[1]) / 2;
    return {
      t: 'L', x: cx - len / 2, y: cy - 13, w: len, h: 26, a: horiz ? 0 : 90,
      sw: wire.sw || 2, c: wire.c || RED, z: 2, lay: 1, net: wire.net, fixc: 1
    };
  }

  /* 刀在 60×98 框内的真实触点（blade viewBox "-6 0 48 124"，缩放 98/124，竖线落在框中线） */
  function knifeTaps(n) {
    var d = 114 * KNIFE_H / 124;
    return { in: { x: n.x + KNIFE_W / 2, y: n.y }, out: { x: n.x + KNIFE_W / 2, y: n.y + d } };
  }

  function build(pg) {
    var lay = layout;
    var scene = indexSceneKnives(pg.nodes || []);
    var nodes = [];
    var cadToKey = {};
    var missing = [];
    (lay.wires || []).forEach(function (wire) {
      var i;
      for (i = 1; i < wire.pts.length; i++) {
        var ln = lineNode(wire.pts[i - 1], wire.pts[i], wire);
        if (ln) nodes.push(ln);
      }
    });
    lay.boxes.forEach(function (b) {
      var n = { t: 'B', x: b.x, y: b.y, w: b.w, h: b.h, z: b.sym ? 4 : 3, lay: 1, id: b.id };
      if (b.sym) n.sym = b.sym;
      nodes.push(n);
    });
    lay.knives.forEach(function (spec) {
      var src = scene[cadKey(spec.cad)];
      if (!src) missing.push(spec.cad);
      var n = knifeNode(spec, src);
      if (src) cadToKey[spec.cad] = src.k;
      nodes.push(n);
    });
    lay.texts.forEach(function (t) {
      var n = { t: 'T', s: t.s, x: t.x, y: t.y, w: t.w, h: t.h, fs: t.fs, c: t.c || '#ffffff', z: 6, lay: 1, al: t.al || 'c' };
      if (t.capDraw) n.capDraw = 1;
      if (t.k && cadToKey[t.k]) n.k = cadToKey[t.k];
      nodes.push(n);
    });
    return { nodes: nodes, missing: missing };
  }

  function apply(pg) {
    if (!applies(pg) || pg._lay5a1) return pg;
    var r = build(pg);
    if (r.missing.length && w.console) w.console.warn('[hx-5a1] 场景中缺少刀点位:', r.missing.join(','));
    pg.nodes = r.nodes;
    pg.w = layout.canvas.w;
    pg.h = layout.canvas.h;
    pg.x0 = layout.canvas.x0;
    pg._lay5a1 = 1;
    pg._lay5a1Missing = r.missing;
    pg._openEnds = layout.openEnds || [];
    pg._cadLinesFixed = 1;
    pg._mainParallelFixed = 1;
    return pg;
  }

  w.HX5A1 = {
    load: load,
    setLayout: setLayout,
    applies: applies,
    apply: apply,
    knifeTaps: knifeTaps,
    normKey: normKey
  };
})(window);
