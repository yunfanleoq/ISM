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
          var fileVer = ver;
          if (/^[34][AB][12]$/.test(room.code)) fileVer = '20261108px';
          if (room.code === '5A4' || room.code === '5B4') fileVer = '20261107merge';
          if (room.code === '5A2' || room.code === '5A3') fileVer = '20261107merge';
          if (/^2[AB][12]$/.test(room.code)) fileVer = '20261108bt';
          return fetch('/hx-dcim/data/hx-power-rooms/' + room.file + '?v=' + fileVer, { cache: 'no-store' })
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

  /* CAD 柜号 → 正式测点回路名。无 scene 时也必须写出 k，否则悬停永远「没有测点」。
   * 例：5A4-413-8 → 配电室5A4_T1_413_8；T1-411-1 → 配电室3A1_T1_411_1；U1-Y11-4 同理。 */
  function circuitFromCad(roomName, cad) {
    var id = String(cad || '').split('#')[0].trim();
    if (!id || !roomName) return '';
    var code = codeOf(roomName);
    var m = id.match(new RegExp('^' + code.replace(/[.*+?^${}()|[\]\\]/g, '\\$&') + '-(.+)$', 'i'));
    if (m) {
      /* 单室图柜号常带室号前缀，测点在 T1/T2 下 */
      return roomName + '_T1_' + m[1].replace(/-/g, '_');
    }
    return roomName + '_' + id.replace(/-/g, '_');
  }

  function knifeNode(spec, src, fit, roomName) {
    var n = {};
    if (src) Object.keys(src).forEach(function (key) { n[key] = src[key]; });
    delete n.rot;
    delete n.orient;
    delete n.hfull;
    if (spec.orient === 'h') n.orient = 'h';
    n.t = 'I';
    n.x = spec.x;
    n.y = spec.y;
    n.w = KNIFE_W;
    n.h = KNIFE_H;
    n.z = 5;
    n.lay = 1;
    n.cad = spec.cad;
    if (!n.k) {
      n.k = circuitFromCad(roomName, spec.cad);
      if (n.k) n.p = n.k + '_开关状态';
    }
    if (spec.out) n.decor = 1; else delete n.decor;
    if (spec.role === 'cap') n.capDraw = 1;
    if (spec.role) n.role = spec.role;
    if (spec.use) n.use = spec.use;
    if (fit) n.bladeFit = 1;
    if (spec.bladeLeft) n.bladeLeft = 1;
    if (spec.bladeFlip) n.bladeFlip = 1;
    return n;
  }

  function lineNode(a, b, wire) {
    var dx = b[0] - a[0];
    var dy = b[1] - a[1];
    var elen = Math.sqrt(dx * dx + dy * dy);
    if (elen < 0.5) return null;
    var horiz = Math.abs(dy) < 0.8;
    var vert = Math.abs(dx) < 0.8;
    var len = horiz ? Math.abs(dx) : (vert ? Math.abs(dy) : elen);
    var cx = (a[0] + b[0]) / 2;
    var cy = (a[1] + b[1]) / 2;
    var net = wire.net;
    if (String(wire.c || '').replace(/\s/g, '').toLowerCase() === '#8aa0b8') net = 'frame';
    var n = {
      t: 'L', x: cx - len / 2, y: cy - 13, w: len, h: 26, a: horiz ? 0 : 90,
      sw: wire.sw || 2, c: wire.c || RED, z: 2, lay: 1, net: net, fixc: 1
    };
    if (!horiz && !vert) {
      n.x1 = a[0];
      n.y1 = a[1];
      n.x2 = b[0];
      n.y2 = b[1];
      n.a = 0;
    }
    if (wire.fix) n.fix = wire.fix;
    if (wire.noHop) n.noHop = 1;
    return n;
  }

  /* 横刀导线轴在刀框顶向下 30px。出线点在框右沿外侧，刀开着则这一侧整段跟绿。 */
  function knifeTaps(n) {
    var midX = n.x + KNIFE_W / 2;
    var midY = n.y + KNIFE_H / 2;
    if (n.orient === 'h') {
      var axisY = n.y + 30;
      var long = n.bladeFit ? 140 : (114 * KNIFE_H / 124);
      return { in: { x: midX - long / 2, y: axisY }, out: { x: midX + KNIFE_W / 2 + 18, y: axisY } };
    }
    var d = 114 * KNIFE_H / 124;
    /* 竖刀长轴与横刀同为 140，出线端跟刀尖。 */
    if (n.bladeFit) d = 140 * 114 / 124;
    return { in: { x: midX, y: n.y }, out: { x: midX, y: n.y + d } };
  }

  function isBusTieKnife(kn) {
    if (!kn) return false;
    if (kn.role === 'tie') return true;
    var blob = String(kn.cad || '') + ' ' + String(kn.use || '') + ' ' + String(kn.s || '');
    return /母\s*联|^445$/.test(blob);
  }

  function build(pg, lay) {
    var prefix = pg.name + '_';
    var code = codeOf(pg.name);
    var fit = code === '5A2' || code === '5A3' || code === '5A4' || code === '5B1' || code === '5B2' || code === '5B3' || code === '5B4' || /^2[AB][12]$/.test(code) || /^[34][AB][12]$/.test(code);
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
      nodes.push(knifeNode(spec, src, fit, pg.name));
    });
    (lay.texts || []).forEach(function (t) {
      var n = { t: 'T', s: t.s, x: t.x, y: t.y, w: t.w, h: t.h, fs: t.fs, c: t.c || '#ffffff', z: 6, lay: 1, al: t.al || 'c' };
      if (t.capDraw) n.capDraw = 1;
      if (t.plate) n.plate = 1;
      if (t.stack) n.stack = 1;
      /* UPS / 机组铭牌：补回路 key，悬停才能取数（如 5A1-U11、UPS1进线） */
      if (!n.k && t.s) {
        var ups = String(t.s).match(/(?:配电室)?(\d[AB]\d)[-_]?U(\d{2})/i);
        if (ups) {
          n.k = '配电室' + ups[1].toUpperCase() + '_U' + ups[2];
          n.devType = 'UPS';
        } else if (/UPS|旁路|XBP/i.test(t.s) && t.k) n.k = t.k;
      }
      nodes.push(n);
    });
    if (/^[34][AB][12]$/.test(code)) {
      addFloor34Frames(nodes);
      landPlateLeads(nodes);
    }
    tagOutletFollow(nodes);
    return { nodes: nodes, missing: missing };
  }

  /* 柜号、kVA、母联是铭牌。房间标题和母线名称不是。 */
  function isNameplateText(n) {
    if (!n || n.t !== 'T' || !(n.w > 28) || !(n.h > 16)) return false;
    if (n.w > 220 || n.h > 160) return false;
    var s = String(n.s || '');
    if (/母线/.test(s)) return false;
    if (/配电室$/.test(s.replace(/\s/g, '')) && s.indexOf('\n') < 0) return false;
    return /T[12]-|U[12]-|445|kVA/.test(s);
  }

  function addFloor34Frames(nodes) {
    var plates = [];
    nodes.forEach(function (n) {
      if (!isNameplateText(n)) return;
      n.plate = 1;
      plates.push(n);
    });
    plates.forEach(function (n) {
      var x0 = n.x;
      var y0 = n.y;
      var x1 = n.x + n.w;
      var y1 = n.y + n.h;
      var edges = [[[x0, y0], [x1, y0]], [[x1, y0], [x1, y1]], [[x1, y1], [x0, y1]], [[x0, y1], [x0, y0]]];
      edges.forEach(function (e) {
        var ln = lineNode(e[0], e[1], { c: '#8aa0b8', sw: 1.25, net: 'frame', noHop: 1 });
        if (!ln) return;
        ln.holdGray = 1;
        ln.noHop = 1;
        nodes.push(ln);
      });
    });
  }

  /* 朝铭牌去的线收到框边。缝不超过 36px 就补上，已经伸进框的收到近边。 */
  function landPlateLeads(nodes) {
    var plates = [];
    nodes.forEach(function (n) {
      if (isNameplateText(n)) plates.push(n);
    });
    nodes.forEach(function (ln) {
      if (ln.t !== 'L' || ln.net === 'frame' || ln.net === 'bus' || ln.holdGray) return;
      var ends = lineEnds(ln);
      if (!ends || ends.length < 2) return;
      var a = { x: ends[0].x, y: ends[0].y };
      var b = { x: ends[1].x, y: ends[1].y };
      var hitA = landPlateEnd(a, b, plates);
      var hitB = landPlateEnd(b, a, plates);
      if (!hitA && !hitB) return;
      if (hitA) a = { x: hitA.x, y: hitA.y };
      if (hitB) b = { x: hitB.x, y: hitB.y };
      var dx = b.x - a.x;
      var dy = b.y - a.y;
      var elen = Math.sqrt(dx * dx + dy * dy);
      if (elen < 0.8) return;
      var horiz = Math.abs(dy) <= 1.5;
      var len = horiz ? Math.abs(dx) : Math.abs(dy);
      var cx = (a.x + b.x) / 2;
      var cy = (a.y + b.y) / 2;
      ln.a = horiz ? 0 : 90;
      ln.w = len;
      ln.x = cx - len / 2;
      ln.y = cy - 13;
      if (ln.x1 != null) {
        ln.x1 = a.x;
        ln.y1 = a.y;
        ln.x2 = b.x;
        ln.y2 = b.y;
      }
    });
  }

  function landPlateEnd(end, other, plates) {
    var horiz = Math.abs(end.y - other.y) <= 1.5 && Math.abs(end.x - other.x) > 4;
    var vert = Math.abs(end.x - other.x) <= 1.5 && Math.abs(end.y - other.y) > 4;
    if (!horiz && !vert) return null;
    var best = null;
    plates.forEach(function (p) {
      var x0 = p.x;
      var y0 = p.y;
      var x1 = p.x + p.w;
      var y1 = p.y + p.h;
      var g;
      var hit = null;
      if (vert && end.x > x0 + 1 && end.x < x1 - 1 && other.y < end.y - 2) {
        if (end.y <= y0 && y0 - end.y <= 36) hit = { g: y0 - end.y, x: end.x, y: y0 };
        else if (end.y > y0 && end.y < y1) hit = { g: end.y - y0, x: end.x, y: y0 };
      } else if (horiz && end.y > y0 + 1 && end.y < y1 - 1) {
        if (other.x < end.x - 2) {
          g = Math.abs(end.x - x0);
          if (end.x <= x0 + 2 && g <= 36) hit = { g: g, x: x0, y: end.y };
        } else if (other.x > end.x + 2) {
          g = Math.abs(end.x - x1);
          if (end.x >= x1 - 2 && g <= 36) hit = { g: g, x: x1, y: end.y };
        }
      }
      if (hit && (!best || hit.g < best.g)) best = hit;
    });
    return best;
  }

  function lineEnds(ln) {
    if (ln.x1 != null && ln.y1 != null && ln.x2 != null && ln.y2 != null) {
      return [{ x: ln.x1, y: ln.y1 }, { x: ln.x2, y: ln.y2 }];
    }
    var cx = ln.x + ln.w / 2;
    var cy = ln.y + 13;
    var half = ln.w / 2;
    if (Math.abs((Number(ln.a) || 0) % 180) < 1) {
      return [{ x: cx - half, y: cy }, { x: cx + half, y: cy }];
    }
    return [{ x: cx, y: cy - half }, { x: cx, y: cy + half }];
  }

  /* 出线侧跟闸合红分绿；母联跳过；铭牌同色 */
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
      var horiz = Math.abs((Number(ln.a) || 0) % 180) < 1;
      if (!horiz) return false;
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
      var floor5 = /^配电室5/.test(String(kn.k || ''));
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
      /* 五层：铭牌把竖线隔开时，出线侧仍跟这把刀。母线不跟着涂绿。 */
      function throughLabel(e, other) {
        if (!floor5) return false;
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
          var stopUps = floor5 && texts.some(function (tx) {
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
              var bridge = throughLabel(e, other);
              if (bridge === 'ups' || bridge === false) return;
            }
            if (oends.some(function (oe) { return isUpstream(oe, taps, horizKnife); }) &&
                !oends.some(function (oe) { return !isUpstream(oe, taps, horizKnife); })) return;
            if (floor5 && oends.some(function (oe) {
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

  function getLayout(code) {
    return layouts[code] || null;
  }

  function setLayout(code, lay) {
    if (code && lay) layouts[code] = lay;
  }

  w.HXRooms = {
    load: load,
    setLayouts: setLayouts,
    getLayout: getLayout,
    setLayout: setLayout,
    applies: applies,
    apply: apply,
    knifeTaps: knifeTaps,
    list: list
  };
})(window);
