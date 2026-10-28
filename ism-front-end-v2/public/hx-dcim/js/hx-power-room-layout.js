/* 1～5 层配电室一次图。每间数据在 data/hx-power-rooms/<室>.json。
 * 5A1 不走这里，仍由 hx-5a1-layout.js 使用已验证的 hx-5a1-cad-layout.json。
 * 刀形不在本文件里画。 */
(function (w) {
  'use strict';

  var KNIFE_W = 60;
  var KNIFE_H = 98;
  var RED = '#ce0c18';
  var layouts = {};
  var index = null;

  function load(ver) {
    if (index) return Promise.resolve(index);
    return fetch('/hx-dcim/data/hx-power-rooms/index.json?v=' + ver)
      .then(function (r) { if (!r.ok) throw new Error('power rooms index ' + r.status); return r.json(); })
      .then(function (j) {
        index = j;
        var rooms = j.rooms || [];
        return Promise.all(rooms.map(function (room) {
          return fetch('/hx-dcim/data/hx-power-rooms/' + room.file + '?v=' + ver)
            .then(function (r) { if (!r.ok) throw new Error(room.file); return r.json(); })
            .then(function (lay) { layouts[room.code] = lay; });
        }));
      })
      .then(function () { return index; });
  }

  function setLayouts(map) {
    layouts = map || {};
    index = { rooms: Object.keys(layouts).map(function (code) { return { code: code }; }) };
  }

  function codeOf(name) {
    var m = String(name || '').match(/^配电室(.+)$/);
    return m ? m[1] : '';
  }

  function applies(pg) {
    var code = codeOf(pg && pg.name);
    if (!code || code === '5A1') return false;
    return !!(layouts[code] && pg.kind === 'power');
  }

  function normKey(k, prefix) {
    return String(k || '').replace(prefix, '').replace(/_+$/, '').replace(/_/g, '-');
  }

  function indexSceneKnives(nodes, prefix) {
    var map = {};
    (nodes || []).forEach(function (n) {
      if (n.t !== 'I' && n.t !== 'S') return;
      if (!n.k || String(n.k).indexOf(prefix) !== 0) return;
      var id = normKey(n.k, prefix);
      if (!map[id] || n.st === '1') map[id] = n;
      /* 场景把母联写成 T1-445，CAD 文字是 445 */
      if (id === 'T1-445' || id === 'T2-445') {
        if (!map['445'] || n.st === '1') map['445'] = n;
      }
    });
    return map;
  }

  function knifeNode(spec, src) {
    var n = {};
    if (src) Object.keys(src).forEach(function (key) { n[key] = src[key]; });
    delete n.rot;
    delete n.orient;
    delete n.hfull;
    n.t = 'I';
    n.x = spec.x;
    n.y = spec.y;
    n.w = KNIFE_W;
    n.h = KNIFE_H;
    n.z = 5;
    n.lay = 1;
    n.cad = spec.cad;
    if (spec.out) n.decor = 1; else delete n.decor;
    if (spec.role === 'cap') n.capDraw = 1;
    return n;
  }

  function lineNode(a, b, wire) {
    var horiz = Math.abs(a[1] - b[1]) < 0.8;
    var len = horiz ? Math.abs(b[0] - a[0]) : Math.abs(b[1] - a[1]);
    if (len < 0.5) return null;
    var cx = (a[0] + b[0]) / 2;
    var cy = (a[1] + b[1]) / 2;
    return {
      t: 'L', x: cx - len / 2, y: cy - 13, w: len, h: 26, a: horiz ? 0 : 90,
      sw: wire.sw || 2, c: wire.c || RED, z: 2, lay: 1, net: wire.net, fixc: 1
    };
  }

  function knifeTaps(n) {
    var d = 114 * KNIFE_H / 124;
    return { in: { x: n.x + KNIFE_W / 2, y: n.y }, out: { x: n.x + KNIFE_W / 2, y: n.y + d } };
  }

  function build(pg, lay) {
    var prefix = pg.name + '_';
    var scene = indexSceneKnives(pg.nodes || [], prefix);
    var nodes = [];
    var missing = [];
    (lay.wires || []).forEach(function (wire) {
      var i;
      for (i = 1; i < wire.pts.length; i++) {
        var ln = lineNode(wire.pts[i - 1], wire.pts[i], wire);
        if (ln) nodes.push(ln);
      }
    });
    (lay.boxes || []).forEach(function (b) {
      var n = { t: 'B', x: b.x, y: b.y, w: b.w, h: b.h, z: b.sym ? 4 : 3, lay: 1, id: b.id };
      if (b.sym) n.sym = b.sym;
      nodes.push(n);
    });
    (lay.knives || []).forEach(function (spec) {
      var src = scene[spec.cad] || scene[String(spec.cad).split('#')[0]];
      if (!src) missing.push(spec.cad);
      nodes.push(knifeNode(spec, src));
    });
    (lay.texts || []).forEach(function (t) {
      var n = { t: 'T', s: t.s, x: t.x, y: t.y, w: t.w, h: t.h, fs: t.fs, c: t.c || '#ffffff', z: 6, lay: 1, al: t.al || 'c' };
      if (t.capDraw) n.capDraw = 1;
      nodes.push(n);
    });
    return { nodes: nodes, missing: missing };
  }

  function apply(pg) {
    if (!applies(pg) || pg._layRoom) return pg;
    var lay = layouts[codeOf(pg.name)];
    var r = build(pg, lay);
    pg.nodes = r.nodes;
    pg.w = lay.canvas.w;
    pg.h = lay.canvas.h;
    pg.x0 = lay.canvas.x0 || 0;
    pg._layRoom = 1;
    pg._layRoomMissing = r.missing;
    pg._cadLinesFixed = 1;
    pg._mainParallelFixed = 1;
    return pg;
  }

  function list() {
    return (index && index.rooms) || [];
  }

  w.HXRooms = {
    load: load,
    setLayouts: setLayouts,
    applies: applies,
    apply: apply,
    knifeTaps: knifeTaps,
    list: list
  };
})(window);
