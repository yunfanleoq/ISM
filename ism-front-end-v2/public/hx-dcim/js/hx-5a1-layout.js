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
    /* 文件里 rot≈90 的是横刀，各自从竖母线引出，不是串成一列的竖干线。 */
    delete n.hfull;
    delete n.rot;
    if (spec.orient === 'h') n.orient = 'h';
    else delete n.orient;
    if (spec.out) n.decor = 1; else delete n.decor;
    if (spec.outWire) n.outWire = 1; else delete n.outWire;
    if (spec.bladeLeft) n.bladeLeft = 1;
    if (spec.role) n.role = spec.role;
    if (spec.use) n.use = spec.use;
    n.bladeFit = 1;
    return n;
  }

  /* 母联：两端都接带电母线，开合都不把两侧线涂成「无电绿」 */
  function isBusTieKnife(kn) {
    if (!kn) return false;
    if (kn.role === 'tie') return true;
    var blob = String(kn.cad || '') + ' ' + String(kn.use || '') + ' ' + String(kn.s || '');
    return /母\s*联|^445$/.test(blob);
  }

  function lineNode(a, b, wire) {
    var horiz = Math.abs(a[1] - b[1]) < 0.01;
    var len = horiz ? Math.abs(b[0] - a[0]) : Math.abs(b[1] - a[1]);
    if (len < 0.5) return null;
    var cx = (a[0] + b[0]) / 2;
    var cy = (a[1] + b[1]) / 2;
    var n = {
      t: 'L', x: cx - len / 2, y: cy - 13, w: len, h: 26, a: horiz ? 0 : 90,
      sw: wire.sw || 2, c: wire.c || RED, z: 2, lay: 1, net: wire.net, fixc: 1
    };
    /* 保留 421-spine / 411-spine 等标记，编辑态禁止误删 */
    if (wire.fix) n.fix = wire.fix;
    return n;
  }

  /* 竖刀：上进下出。横刀（rotate -90°）：左进右出。
     横刀导线轴在刀框顶向下 30px，不在 98px 框中心（中心低 19px，吸不到出线）。
     出线点放到框右沿外侧 18px，盖住刀尖到出线起点约 12px 的空档。 */
  function knifeTaps(n) {
    var midX = n.x + KNIFE_W / 2;
    var midY = n.y + KNIFE_H / 2;
    if (n.orient === 'h') {
      var axisY = n.y + 30;
      var long = n.bladeFit ? 55.2 : (114 * KNIFE_H / 124);
      return { in: { x: midX - long / 2, y: axisY }, out: { x: midX + KNIFE_W / 2 + 18, y: axisY } };
    }
    var d = 114 * KNIFE_H / 124;
    return { in: { x: midX, y: n.y }, out: { x: midX, y: n.y + d } };
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
      if (t.s && String(t.s).indexOf('\n') >= 0) n.stack = 1;
      if (t.capDraw) n.capDraw = 1;
      if (t.k && cadToKey[t.k]) n.k = cadToKey[t.k];
      /* UPS 机组铭牌：5A1-U11 / 500kVA → 配电室5A1_U11；上沿进线、下沿出线 */
      if (!n.k && t.s) {
        var ups = String(t.s).match(/(?:5A1|配电室5A1)[-_]?U(\d{2})/i) || String(t.s).match(/\bU(\d{2})\b/);
        if (ups && /U1[1-4]|U2[1-4]/i.test('U' + ups[1])) {
          n.k = '配电室5A1_U' + ups[1];
          n.devType = 'UPS';
          n.upsPlate = 1;
        }
      }
      if (t.s && /500\s*kVA/i.test(String(t.s))) n.upsPlate = 1;
      nodes.push(n);
    });
    tagOutletFollow(nodes);
    return { nodes: nodes, missing: missing };
  }

  function lineEnds(ln) {
    var cx = ln.x + ln.w / 2;
    var cy = ln.y + 13;
    var half = ln.w / 2;
    if (Math.abs((Number(ln.a) || 0) % 180) < 1) {
      return [{ x: cx - half, y: cy }, { x: cx + half, y: cy }];
    }
    return [{ x: cx, y: cy - half }, { x: cx, y: cy + half }];
  }

  /* 出线侧（竖刀下方 / 横刀右侧）跟闸合红分绿；母联跳过；不爬回进线、不吞母线；铭牌同色 */
  function tagOutletFollow(nodes) {
    var knives = [];
    var lines = [];
    var texts = [];
    nodes.forEach(function (n) {
      if ((n.t === 'I' || n.t === 'S') && n.k) knives.push(n);
      if (n.t === 'L' && n.net !== 'frame') lines.push(n);
      if (n.t === 'T') texts.push(n);
    });
    function mark(ln, k) {
      ln.k = k;
      ln.followKnife = 1;
      delete ln.fixc;
    }
    function touchesOutlet(ln, ox, oy) {
      var ends = lineEnds(ln);
      var i;
      for (i = 0; i < ends.length; i++) {
        if (Math.abs(ends[i].x - ox) <= 12 && Math.abs(ends[i].y - oy) <= 12) return true;
      }
      var horiz = Math.abs((Number(ln.a) || 0) % 180) < 1;
      var cx = ln.x + ln.w / 2;
      var cy = ln.y + 13;
      if (horiz) {
        var left = cx - ln.w / 2;
        return Math.abs(cy - oy) <= 10 && left <= ox + 12 && cx + ln.w / 2 >= ox - 6;
      }
      var top = cy - ln.w / 2;
      return Math.abs(cx - ox) <= 10 && top <= oy + 14 && cy + ln.w / 2 >= oy - 6;
    }
    function isMainBus(ln) {
      if (ln.net === 'bus') return true;
      return Number(ln.w) > 280;
    }
    function isUpstream(e, taps, horizKnife) {
      if (horizKnife) return e.x < taps.in.x - 2;
      return e.y < taps.in.y - 2;
    }
    knives.forEach(function (kn) {
      if (isBusTieKnife(kn)) return;
      var taps = knifeTaps(kn);
      var ox = taps.out.x;
      var oy = taps.out.y;
      var horizKnife = kn.orient === 'h';
      var queue = [];
      var seen = {};
      lines.forEach(function (ln, i) {
        if (isMainBus(ln)) return;
        if (!touchesOutlet(ln, ox, oy)) return;
        seen[i] = 1;
        queue.push(i);
      });
      function isUpsPlate(tx) {
        var s = String(tx.s || '');
        return /500\s*kVA/i.test(s) || /(?:^|[^0-9])U1[1-4](?:\b|$)|(?:^|[^0-9])U2[1-4](?:\b|$)/.test(s);
      }
      function hitsPlateTop(e, tx) {
        return e.x >= tx.x - 6 && e.x <= tx.x + tx.w + 6 &&
          e.y >= tx.y - 10 && e.y <= tx.y + 12;
      }
      function hitsPlateBottom(e, tx) {
        var bot = tx.y + tx.h;
        return e.x >= tx.x - 6 && e.x <= tx.x + tx.w + 6 &&
          e.y >= bot - 12 && e.y <= bot + 10;
      }
      function throughLabel(e, other, j) {
        /* 共线竖缝被铭牌隔开：电气连通，几何不穿框；UPS 铭牌挡住，不往出线侧继续跟上游刀 */
        if (Math.abs((Number(other.a) || 0) % 180 - 90) > 1) return false;
        var oends = lineEnds(other);
        var oe = oends[0];
        if (Math.abs(oe.x - e.x) > 10) return false;
        var oy0 = Math.min(oends[0].y, oends[1].y);
        var oy1 = Math.max(oends[0].y, oends[1].y);
        var gapLo = Math.min(e.y, oy1);
        var gapHi = Math.max(e.y, oy0);
        if (oy0 >= e.y - 2) {
          gapLo = e.y;
          gapHi = oy0;
        } else if (oy1 <= e.y + 2) {
          gapLo = oy1;
          gapHi = e.y;
        } else return false;
        if (gapHi - gapLo < 8 || gapHi - gapLo > 120) return false;
        var i;
        for (i = 0; i < texts.length; i++) {
          var tx = texts[i];
          if (!tx.w || !tx.h) continue;
          if (e.x < tx.x - 8 || e.x > tx.x + tx.w + 8) continue;
          if (tx.y <= gapLo + 4 && tx.y + tx.h >= gapHi - 4) {
            if (isUpsPlate(tx)) return 'ups';
            return 'label';
          }
        }
        return false;
      }
      while (queue.length) {
        var idx = queue.shift();
        var cur = lines[idx];
        mark(cur, kn.k);
        lineEnds(cur).forEach(function (e) {
          if (isUpstream(e, taps, horizKnife)) return;
          /* 碰到 UPS 上沿：进线跟闸到此为止，出线不跟这把上游刀 */
          var stopUps = texts.some(function (tx) {
            return tx.w && tx.h && isUpsPlate(tx) && hitsPlateTop(e, tx);
          });
          if (stopUps) return;
          lines.forEach(function (other, j) {
            if (seen[j] || isMainBus(other)) return;
            var oends = lineEnds(other);
            var conn = oends.some(function (oe) {
              return Math.abs(oe.x - e.x) <= 8 && Math.abs(oe.y - e.y) <= 8;
            });
            if (!conn) {
              var bridge = throughLabel(e, other, j);
              if (bridge === 'ups' || bridge === false) return;
            }
            if (oends.some(function (oe) { return isUpstream(oe, taps, horizKnife); }) &&
                !oends.some(function (oe) { return !isUpstream(oe, taps, horizKnife); })) return;
            /* 下一段若已在 UPS 下沿之外，不跟 */
            if (oends.some(function (oe) {
              return texts.some(function (tx) {
                return tx.w && tx.h && isUpsPlate(tx) && hitsPlateBottom(oe, tx);
              });
            })) return;
            seen[j] = 1;
            queue.push(j);
          });
        });
      }
      texts.forEach(function (tx) {
        if (!tx.w || !tx.h) return;
        if (isUpsPlate(tx)) {
          tx.upsPlate = 1;
          return;
        }
        var already = tx.k && String(tx.k) === String(kn.k);
        var hit = lines.some(function (ln, i) {
          if (!seen[i]) return false;
          return lineEnds(ln).some(function (e) {
            return e.x >= tx.x - 4 && e.x <= tx.x + tx.w + 4 &&
              e.y >= tx.y - 4 && e.y <= tx.y + tx.h + 4;
          });
        });
        if (already || hit) {
          tx.k = kn.k;
          tx.followKnife = 1;
        }
      });
    });
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
