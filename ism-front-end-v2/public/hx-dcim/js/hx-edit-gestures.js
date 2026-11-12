/* 一次图编辑的点选、拖拽、改字、删除。画面数据在 hx-edit.js。 */
(function (w) {
  'use strict';

  function install(A) {
    var drag = null;

    function onPalDown(ev) {
      var item = ev.target.closest ? ev.target.closest('[data-hx-add]') : null;
      if (!item || !A.isOn()) return;
      ev.preventDefault();
      var start = { x: ev.clientX, y: ev.clientY, kind: item.getAttribute('data-hx-add'), moved: false };
      function move(e) {
        if (Math.hypot(e.clientX - start.x, e.clientY - start.y) > 5) start.moved = true;
      }
      function up(e) {
        document.removeEventListener('pointermove', move);
        document.removeEventListener('pointerup', up);
        if (!A.isOn()) return;
        if (!start.moved) { A.addAtCenter(start.kind); return; }
        if (e.target.closest && e.target.closest('.hx-edit-pal')) return;
        var host = A.host();
        var fit = host && host.querySelector('.hx-fit');
        if (!fit) return;
        var fr = fit.getBoundingClientRect();
        if (e.clientX < fr.left || e.clientX > fr.right || e.clientY < fr.top || e.clientY > fr.bottom) return;
        A.addAt(start.kind, e.clientX, e.clientY);
      }
      document.addEventListener('pointermove', move);
      document.addEventListener('pointerup', up);
    }

    function closestBase(node) {
      if (!node || !node.closest) return null;
      var el = node.closest('[data-hx-knife], [data-hx-sym], .hx-tx, [data-i]');
      if (!el || (el.closest && el.closest('.hx-edit-layer'))) return null;
      return A.kindOf(el) ? el : null;
    }

    /* 点到线段距离（板坐标） */
    function distToSeg(px, py, x1, y1, x2, y2) {
      var dx = x2 - x1;
      var dy = y2 - y1;
      var len2 = dx * dx + dy * dy;
      var t = len2 < 1e-6 ? 0 : Math.max(0, Math.min(1, ((px - x1) * dx + (py - y1) * dy) / len2));
      var qx = x1 + t * dx;
      var qy = y1 + t * dy;
      return Math.hypot(px - qx, py - qy);
    }

    function lineGeom(el) {
      var tag = (el.tagName || '').toLowerCase();
      if (tag === 'line') {
        return {
          x1: Number(el.getAttribute('x1')), y1: Number(el.getAttribute('y1')),
          x2: Number(el.getAttribute('x2')), y2: Number(el.getAttribute('y2'))
        };
      }
      if (tag === 'path') {
        var d = String(el.getAttribute('d') || '');
        var m = d.match(/M\s*([-\d.]+)\s*([-\d.]+)/i);
        var nums = d.match(/[-\d.]+/g) || [];
        if (!m || nums.length < 4) return null;
        return {
          x1: Number(nums[0]), y1: Number(nums[1]),
          x2: Number(nums[nums.length - 2]), y2: Number(nums[nums.length - 1])
        };
      }
      return null;
    }

    /* 在点击处找出被铭牌/刀盖住的电源线（含无 data-hx-key 的 CAD 线） */
    function wiresUnder(clientX, clientY, tol) {
      var host = A.host();
      var shift = host && host.querySelector('.hx-shift');
      if (!shift || !A.pageXY) return [];
      var pt = A.pageXY(clientX, clientY);
      if (!pt) return [];
      var hitTol = tol == null ? 10 : tol;
      var out = [];
      var seen = {};
      Array.prototype.forEach.call(shift.querySelectorAll('.hx-wires [data-hx-line]'), function (el) {
        if (el.getAttribute('stroke') === 'transparent') return;
        var i = el.getAttribute('data-i');
        if (i == null || seen[i]) return;
        var g = lineGeom(el);
        if (!g) return;
        if (distToSeg(pt.x, pt.y, g.x1, g.y1, g.x2, g.y2) > hitTol) return;
        seen[i] = 1;
        out.push(el);
      });
      return out;
    }

    var lastPick = null;

    function pickBase(ev) {
      var top = closestBase(ev.target);
      var wires = wiresUnder(ev.clientX, ev.clientY, 12);
      var pierce = !!(ev.altKey || ev.ctrlKey || ev.metaKey);
      if (pierce && wires.length) return wires[0];
      /* 同一位置再点：在顶层件与下方线之间切换 */
      if (top && wires.length) {
        var key = Math.round(ev.clientX) + ',' + Math.round(ev.clientY);
        var topKind = A.kindOf(top);
        if (topKind && topKind !== 'line') {
          if (lastPick && lastPick.key === key && lastPick.i === top.getAttribute('data-i')) {
            lastPick = { key: key, i: wires[0].getAttribute('data-i') };
            return wires[0];
          }
          lastPick = { key: key, i: top.getAttribute('data-i') };
          return top;
        }
      }
      if (top) return top;
      return wires[0] || null;
    }

    function disarm() {
      document.removeEventListener('pointermove', onMove);
      document.removeEventListener('pointerup', onUp);
      document.removeEventListener('pointercancel', onUp);
    }

    function armDrag(ev, state) {
      drag = state;
      drag.pointer = ev.pointerId;
      document.addEventListener('pointermove', onMove);
      document.addEventListener('pointerup', onUp);
      document.addEventListener('pointercancel', onUp);
      var host = A.host();
      var main = host && host.querySelector('.hx-main');
      if (main && main.setPointerCapture) {
        try { main.setPointerCapture(ev.pointerId); } catch (e) {}
      }
      ev.preventDefault();
      ev.stopPropagation();
    }

    function onDown(ev) {
      if (!A.isOn() || ev.button !== 0) return;
      if (ev.target.closest && ev.target.closest('.hx-zoom-bar, .hx-edit-pal, .hx-edit-style')) return;
      if (ev.target.closest && ev.target.closest('.is-hx-texting, [contenteditable="true"]')) return;
      /* 「局部」框选放大时让给舞台缩放，结束后仍回到编辑拖拽 */
      var host0 = A.host();
      var fit0 = host0 && host0.querySelector('.hx-fit');
      if (fit0 && fit0.classList.contains('is-marquee')) return;
      var handle = ev.target.closest && ev.target.closest('.hx-edit-handle');
      if (handle) {
        if (handle.getAttribute('data-hx-base')) {
          beginBaseResize(ev, handle.getAttribute('data-i'), handle.getAttribute('data-end'));
          return;
        }
        beginItem(ev, handle.getAttribute('data-hx-edit-id'), handle.getAttribute('data-end'));
        return;
      }
      var itemEl = ev.target.closest && ev.target.closest('[data-hx-edit-id]');
      /* Alt/Ctrl：优先点穿到下层线，不抢编辑层新加的线 */
      if (itemEl && !itemEl.classList.contains('hx-edit-handle') && !(ev.altKey || ev.ctrlKey || ev.metaKey)) {
        beginItem(ev, itemEl.getAttribute('data-hx-edit-id'), '');
        return;
      }
      if ((ev.altKey || ev.ctrlKey || ev.metaKey) && itemEl) {
        var under = wiresUnder(ev.clientX, ev.clientY, 12);
        if (under.length) { beginBase(ev, under[0]); return; }
      }
      var base = pickBase(ev);
      if (base) { beginBase(ev, base); return; }
      A.setSel(null);
      A.renderItems();
      A.markSel();
    }

    function beginBase(ev, el) {
      var page = A.page();
      var i = String(el.getAttribute('data-i'));
      var n = page && page.nodes[Number(i)];
      if (!n) return;
      var ov = A.load();
      if (ov.deleted.indexOf(i) >= 0) return;
      var mv = ov.moves[i] || {};
      var prev = A.sel();
      var was = !!(prev && prev.t === 'base' && String(prev.i) === i);
      A.setSel({ t: 'base', i: i });
      A.markSel();
      var geom = null;
      if (n.t === 'L' && A.lineGeom) geom = A.lineGeom(n, mv);
      else if ((n.t === 'T' || n.t === 'B') && A.boxGeom) geom = A.boxGeom(n, mv);
      armDrag(ev, {
        kind: 'base', i: i, n: n, el: el, was: was, end: '',
        geom: geom,
        absLine: !!(n.t === 'L' && mv.x1 != null),
        absBox: !!((n.t === 'T' || n.t === 'B') && mv.w != null && mv.x != null),
        axis: n.t === 'L' && A.lineAxis ? A.lineAxis(n) : '',
        cx: ev.clientX, cy: ev.clientY,
        dx: mv.dx || 0, dy: mv.dy || 0, s: mv.s, fs: mv.fs, c: mv.c,
        moved: false
      });
    }

    function beginBaseResize(ev, i, end) {
      var page = A.page();
      var n = page && page.nodes[Number(i)];
      if (!n) return;
      var ov = A.load();
      var mv = ov.moves[i] || {};
      A.setSel({ t: 'base', i: String(i) });
      A.markSel();
      var geom = n.t === 'L' ? A.lineGeom(n, mv) : A.boxGeom(n, mv);
      armDrag(ev, {
        kind: 'base', i: String(i), n: n, end: end || '',
        geom: geom,
        axis: n.t === 'L' ? A.lineAxis(n) : '',
        cx: ev.clientX, cy: ev.clientY,
        s: mv.s, fs: mv.fs, c: mv.c,
        moved: false
      });
    }

    function beginItem(ev, id, end) {
      var ov = A.load();
      var it = null;
      ov.items.forEach(function (x) { if (x.id === id) it = x; });
      if (!it) return;
      A.setSel({ t: 'item', id: id });
      A.renderItems();
      A.markSel();
      armDrag(ev, {
        kind: 'item', end: end || '',
        snap: JSON.parse(JSON.stringify(it)),
        cx: ev.clientX, cy: ev.clientY,
        moved: false, live: null
      });
    }

    function shiftItem(it, dx, dy) {
      if (it.kind === 'hline' || it.kind === 'vline' || it.kind === 'ortho' || it.kind === 'bus') {
        it.x1 = Number(it.x1) + dx; it.x2 = Number(it.x2) + dx;
        it.y1 = Number(it.y1) + dy; it.y2 = Number(it.y2) + dy;
        if (it.xm != null) it.xm = Number(it.xm) + dx;
        if (it.ym != null) it.ym = Number(it.ym) + dy;
        return;
      }
      it.x = Number(it.x) + dx;
      it.y = Number(it.y) + dy;
    }

    /* 屏幕上大约 14 像素。放大缩小后，按画面坐标换算，手感不变。 */
    var SNAP_PX = 14;

    function axisFrom(g) {
      if (!g || g.x1 == null) return 'd';
      if (Math.abs(g.y1 - g.y2) < 1.5) return 'h';
      if (Math.abs(g.x1 - g.x2) < 1.5) return 'v';
      return 'd';
    }

    function collectEnds(skip) {
      var pg = A.page && A.page();
      var ov = A.load() || {};
      var moves = ov.moves || {};
      var pts = [];
      var nodes = (pg && pg.nodes) || [];
      nodes.forEach(function (n, i) {
        if (!n || n.t !== 'L') return;
        if (skip && skip.base != null && String(skip.base) === String(i)) return;
        var g = A.lineGeom(n, moves[i] || moves[String(i)] || null);
        pts.push({ x: g.x1, y: g.y1 }, { x: g.x2, y: g.y2 });
      });
      (ov.items || []).forEach(function (it) {
        if (!it || it.x1 == null) return;
        if (skip && skip.item && skip.item === it.id) return;
        if (it.kind === 'ortho' && w.HXEditCmd && w.HXEditCmd.orthoPts) {
          var flat = w.HXEditCmd.orthoPts(it);
          var k;
          for (k = 0; k < flat.length; k += 2) pts.push({ x: Number(flat[k]), y: Number(flat[k + 1]) });
          return;
        }
        pts.push({ x: Number(it.x1), y: Number(it.y1) }, { x: Number(it.x2), y: Number(it.y2) });
      });
      return pts;
    }

    function nearestPoint(x, y, pts, radius, hold) {
      if (hold && Math.hypot(hold.x - x, hold.y - y) <= radius * 1.45) return hold;
      var best = null;
      var bestD = radius;
      var i;
      for (i = 0; i < pts.length; i++) {
        var d = Math.hypot(pts[i].x - x, pts[i].y - y);
        if (d > 0.4 && d < bestD) { bestD = d; best = pts[i]; }
      }
      return best;
    }

    /* 横线、竖线吸上之后仍保持平直：整段平移到目标那一轴，被拖的那头落到目标点。 */
    function lockOnto(g, axis, end, hit) {
      var o = { x1: g.x1, y1: g.y1, x2: g.x2, y2: g.y2 };
      if (axis === 'h') {
        o.y1 = hit.y;
        o.y2 = hit.y;
        if (end === '1') o.x1 = hit.x;
        else o.x2 = hit.x;
      } else if (axis === 'v') {
        o.x1 = hit.x;
        o.x2 = hit.x;
        if (end === '1') o.y1 = hit.y;
        else o.y2 = hit.y;
      } else if (end === '1') {
        o.x1 = hit.x;
        o.y1 = hit.y;
      } else {
        o.x2 = hit.x;
        o.y2 = hit.y;
      }
      return o;
    }

    function magnetLine(g, axis, which) {
      var radius = SNAP_PX / (A.scale() || 1);
      var pts = collectEnds(drag && drag.magnetSkip);
      var hit = null;
      var next = g;
      var hold = drag && drag.magnetHold;
      if (which === '1' || which === '2') {
        var px = which === '1' ? g.x1 : g.x2;
        var py = which === '1' ? g.y1 : g.y2;
        hit = nearestPoint(px, py, pts, radius, hold);
        if (hit) next = lockOnto(g, axis, which, hit);
      } else {
        var h1 = nearestPoint(g.x1, g.y1, pts, radius, null);
        var h2 = nearestPoint(g.x2, g.y2, pts, radius, null);
        var d1 = h1 ? Math.hypot(h1.x - g.x1, h1.y - g.y1) : 1e9;
        var d2 = h2 ? Math.hypot(h2.x - g.x2, h2.y - g.y2) : 1e9;
        if (hold) {
          var dh1 = Math.hypot(hold.x - g.x1, hold.y - g.y1);
          var dh2 = Math.hypot(hold.x - g.x2, hold.y - g.y2);
          if (Math.min(dh1, dh2) <= radius * 1.45) {
            hit = hold;
            d1 = dh1;
            d2 = dh2;
            h1 = hold;
            h2 = hold;
          }
        }
        if (d1 < 1e9 || d2 < 1e9) {
          hit = d1 <= d2 ? h1 : h2;
          var sx = d1 <= d2 ? g.x1 : g.x2;
          var sy = d1 <= d2 ? g.y1 : g.y2;
          var dx = hit.x - sx;
          var dy = hit.y - sy;
          next = { x1: g.x1 + dx, y1: g.y1 + dy, x2: g.x2 + dx, y2: g.y2 + dy };
        }
      }
      if (hit && Math.hypot(next.x2 - next.x1, next.y2 - next.y1) < 8) hit = null;
      if (drag) drag.magnetHold = hit || null;
      if (hit) showMagnet(hit.x, hit.y);
      else hideMagnet();
      return hit ? next : g;
    }

    function showMagnet(x, y) {
      var layer = A.ensureLayer();
      if (!layer) return;
      var el = layer.querySelector('.hx-snap-magnet');
      if (!el) {
        el = document.createElement('div');
        el.className = 'hx-snap-magnet';
        layer.appendChild(el);
      }
      el.style.display = 'block';
      el.style.left = (x - 11) + 'px';
      el.style.top = (y - 11) + 'px';
    }

    function hideMagnet() {
      var layer = A.ensureLayer && A.ensureLayer();
      if (!layer) return;
      var el = layer.querySelector('.hx-snap-magnet');
      if (el) el.style.display = 'none';
    }

    function shiftEnd(it, end, dx, dy) {
      if (end === 'nw' || end === 'n' || end === 'ne' || end === 'e' || end === 'se' || end === 's' || end === 'sw' || end === 'w') {
        var box = A.resizeBox({ x: Number(it.x), y: Number(it.y), w: Number(it.w), h: Number(it.h) }, end, dx, dy);
        it.x = box.x;
        it.y = box.y;
        it.w = box.w;
        it.h = box.h;
        return;
      }
      if (it.kind === 'ortho' && w.HXEditCmd) {
        if (end === 'm') {
          var b = w.HXEditCmd.bendOf(it);
          it.bend = b;
          if (b === 'hv') it.xm = (it.xm != null ? Number(it.xm) : (Number(it.x1) + Number(it.x2)) / 2) + dx;
          else it.ym = (it.ym != null ? Number(it.ym) : (Number(it.y1) + Number(it.y2)) / 2) + dy;
          return;
        }
        if (end === '1') { it.x1 = Number(it.x1) + dx; it.y1 = Number(it.y1) + dy; }
        else { it.x2 = Number(it.x2) + dx; it.y2 = Number(it.y2) + dy; }
        delete it.xm;
        delete it.ym;
        return;
      }
      if (it.kind === 'hline' || it.kind === 'bus') {
        if (end === '1') it.x1 = Number(it.x1) + dx;
        else it.x2 = Number(it.x2) + dx;
        it.y2 = it.y1;
        if (Math.abs(it.x2 - it.x1) < 8) it.x2 = it.x1 + (it.x2 >= it.x1 ? 8 : -8);
      } else if (it.kind === 'vline') {
        if (end === '1') it.y1 = Number(it.y1) + dy;
        else it.y2 = Number(it.y2) + dy;
        it.x2 = it.x1;
        if (Math.abs(it.y2 - it.y1) < 8) it.y2 = it.y1 + (it.y2 >= it.y1 ? 8 : -8);
      }
    }

    function paintItem(it) {
      var layer = A.ensureLayer();
      if (!layer) return;
      if (it.kind === 'hline' || it.kind === 'vline' || it.kind === 'ortho' || it.kind === 'bus') {
        if (it.kind === 'ortho' && w.HXEditCmd) {
          var d = w.HXEditCmd.orthoD(w.HXEditCmd.orthoPts(it));
          Array.prototype.forEach.call(layer.querySelectorAll('path[data-hx-edit-id="' + it.id + '"]'), function (p) {
            p.setAttribute('d', d);
          });
        } else {
          Array.prototype.forEach.call(layer.querySelectorAll('line[data-hx-edit-id="' + it.id + '"]'), function (line) {
            line.setAttribute('x1', it.x1); line.setAttribute('y1', it.y1);
            line.setAttribute('x2', it.x2); line.setAttribute('y2', it.y2);
          });
        }
        var h1 = layer.querySelector('.hx-edit-handle[data-end="1"]');
        var h2 = layer.querySelector('.hx-edit-handle[data-end="2"]');
        if (h1) { h1.style.left = (it.x1 - 6) + 'px'; h1.style.top = (it.y1 - 6) + 'px'; }
        if (h2) { h2.style.left = (it.x2 - 6) + 'px'; h2.style.top = (it.y2 - 6) + 'px'; }
        if (it.kind === 'ortho' && w.HXEditCmd) {
          var mid = w.HXEditCmd.orthoMid(it);
          var hm = layer.querySelector('.hx-edit-handle[data-end="m"]');
          if (hm) { hm.style.left = (mid.x - 6) + 'px'; hm.style.top = (mid.y - 6) + 'px'; }
        }
        return;
      }
      var el = layer.querySelector('.hx-edit-item[data-hx-edit-id="' + it.id + '"]');
      if (!el) return;
      el.style.left = it.x + 'px';
      el.style.top = it.y + 'px';
      if (it.w != null) el.style.width = it.w + 'px';
      if (it.h != null) el.style.height = it.h + 'px';
      var dots = {
        nw: [it.x, it.y],
        n: [it.x + it.w / 2, it.y],
        ne: [it.x + it.w, it.y],
        e: [it.x + it.w, it.y + it.h / 2],
        se: [it.x + it.w, it.y + it.h],
        s: [it.x + it.w / 2, it.y + it.h],
        sw: [it.x, it.y + it.h],
        w: [it.x, it.y + it.h / 2]
      };
      Object.keys(dots).forEach(function (end) {
        var h = layer.querySelector('.hx-edit-handle[data-hx-edit-id="' + it.id + '"][data-end="' + end + '"]');
        if (!h) return;
        h.style.left = (dots[end][0] - 6) + 'px';
        h.style.top = (dots[end][1] - 6) + 'px';
      });
    }

    function onMove(ev) {
      if (!drag || ev.pointerId !== drag.pointer) return;
      var sc = A.scale();
      var dx = (ev.clientX - drag.cx) / sc;
      var dy = (ev.clientY - drag.cy) / sc;
      if (!drag.moved && Math.hypot(ev.clientX - drag.cx, ev.clientY - drag.cy) < 3) return;
      drag.moved = true;
      if (drag.kind === 'base') {
        var scb = A.scale();
        var dbx = (ev.clientX - drag.cx) / scb;
        var dby = (ev.clientY - drag.cy) / scb;
        var live = null;
        if (drag.end && drag.n.t === 'L') live = A.resizeLine(drag.geom, drag.axis, drag.end, dbx, dby);
        else if (drag.end) live = A.resizeBox(drag.geom, drag.end, dbx, dby);
        else if (drag.absLine) {
          live = {
            x1: drag.geom.x1 + dbx, y1: drag.geom.y1 + dby,
            x2: drag.geom.x2 + dbx, y2: drag.geom.y2 + dby
          };
        } else if (drag.absBox) {
          live = { x: drag.geom.x + dbx, y: drag.geom.y + dby, w: drag.geom.w, h: drag.geom.h };
        }
        if (live) {
          if (drag.n.t === 'L') {
            drag.magnetSkip = { base: drag.i };
            var axisL = (drag.axis === 'h' || drag.axis === 'v') ? drag.axis : axisFrom(live);
            live = magnetLine(live, axisL, drag.end || '');
            drag.live = { x1: live.x1, y1: live.y1, x2: live.x2, y2: live.y2 };
          } else {
            hideMagnet();
            drag.live = { x: live.x, y: live.y, w: live.w, h: live.h };
          }
          if (drag.s != null) drag.live.s = drag.s;
          if (drag.fs != null) drag.live.fs = drag.fs;
          if (drag.c) drag.live.c = drag.c;
          A.preview(drag.i, drag.n, drag.live);
          return;
        }
        if (drag.n.t === 'L' && drag.geom) {
          drag.magnetSkip = { base: drag.i };
          var body = {
            x1: drag.geom.x1 + dbx, y1: drag.geom.y1 + dby,
            x2: drag.geom.x2 + dbx, y2: drag.geom.y2 + dby
          };
          var snapped = magnetLine(body, (drag.axis === 'h' || drag.axis === 'v') ? drag.axis : axisFrom(body), '');
          if (drag.magnetHold) {
            drag.live = { x1: snapped.x1, y1: snapped.y1, x2: snapped.x2, y2: snapped.y2, s: drag.s, fs: drag.fs, c: drag.c };
            A.preview(drag.i, drag.n, drag.live);
            return;
          }
          drag.live = null;
        } else hideMagnet();
        A.preview(drag.i, drag.n, { dx: drag.dx + dbx, dy: drag.dy + dby, s: drag.s, fs: drag.fs, c: drag.c });
        return;
      }
      var it = JSON.parse(JSON.stringify(drag.snap));
      if (drag.end) shiftEnd(it, drag.end, dx, dy);
      else shiftItem(it, dx, dy);
      if (it.x1 != null && drag.end !== 'm' && (it.kind === 'hline' || it.kind === 'vline' || it.kind === 'bus' || it.kind === 'ortho')) {
        drag.magnetSkip = { item: it.id };
        var axisI = it.kind === 'vline' ? 'v' : (it.kind === 'ortho' ? 'd' : 'h');
        var which = (drag.end === '1' || drag.end === '2') ? drag.end : '';
        var before = { x1: Number(it.x1), y1: Number(it.y1), x2: Number(it.x2), y2: Number(it.y2) };
        var after = magnetLine(before, axisI, which);
        if (drag.magnetHold) {
          if (!which) {
            if (it.xm != null) it.xm = Number(it.xm) + (after.x1 - before.x1);
            if (it.ym != null) it.ym = Number(it.ym) + (after.y1 - before.y1);
          }
          it.x1 = after.x1; it.y1 = after.y1; it.x2 = after.x2; it.y2 = after.y2;
        }
      } else hideMagnet();
      drag.live = it;
      paintItem(it);
    }

    function roundItem(it) {
      var o = JSON.parse(JSON.stringify(it));
      ['x', 'y', 'w', 'h', 'x1', 'y1', 'x2', 'y2', 'xm', 'ym'].forEach(function (k) {
        if (o[k] != null) o[k] = A.r1(o[k]);
      });
      return o;
    }

    function onUp(ev) {
      if (!drag) return;
      if (ev && ev.pointerId != null && ev.pointerId !== drag.pointer) return;
      var d = drag;
      drag = null;
      hideMagnet();
      disarm();
      if (!d.moved) {
        if (d.kind === 'base' && !d.end && d.was && d.n && d.n.t === 'T' && d.el) startTextEdit(d.el);
        return;
      }
      if (!ev) return;
      var sc = A.scale();
      var dx = (ev.clientX - d.cx) / sc;
      var dy = (ev.clientY - d.cy) / sc;
      if (d.kind === 'base') {
        if (d.live) {
          A.commit(function (ov) {
            var rec = {};
            ['x', 'y', 'w', 'h', 'x1', 'y1', 'x2', 'y2'].forEach(function (k) {
              if (d.live[k] != null) rec[k] = A.r1(d.live[k]);
            });
            if (d.s != null) rec.s = d.s;
            if (d.fs != null) rec.fs = d.fs;
            if (d.c) rec.c = d.c;
            ov.moves[d.i] = rec;
          });
          A.markSel();
          return;
        }
        A.commit(function (ov) {
          var rec = { dx: A.r1(d.dx + dx), dy: A.r1(d.dy + dy) };
          if (d.s != null) rec.s = d.s;
          if (d.fs != null) rec.fs = d.fs;
          if (d.c) rec.c = d.c;
          if (!rec.dx && !rec.dy && rec.s == null && rec.fs == null && !rec.c) delete ov.moves[d.i];
          else ov.moves[d.i] = rec;
        });
        A.markSel();
        return;
      }
      var it = roundItem(d.live || d.snap);
      A.commit(function (ov) {
        var i;
        for (i = 0; i < ov.items.length; i++) if (ov.items[i].id === it.id) ov.items[i] = it;
      });
      A.renderItems();
      A.markSel();
    }

    function applyText(el, next) {
      if (A.textOf) A.textOf(el, next);
      else el.innerHTML = String(next == null ? '' : next)
        .replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
        .replace(/\n/g, '<br>');
    }

    function commitText(el, next) {
      if (el.closest('.hx-edit-layer')) {
        var id = el.getAttribute('data-hx-edit-id');
        A.commit(function (ov) {
          var i;
          for (i = 0; i < ov.items.length; i++) if (ov.items[i].id === id) ov.items[i].s = next;
        });
        applyText(el, next);
        return;
      }
      var i = el.getAttribute('data-i');
      if (i == null) return;
      A.commit(function (ov) {
        var mv = ov.moves[i] || { dx: 0, dy: 0 };
        mv.s = next;
        ov.moves[i] = mv;
      });
      applyText(el, next);
    }

    function lineEnds(n) {
      if (!n) return null;
      if (n.x1 != null && n.y1 != null && n.x2 != null && n.y2 != null) {
        return [{ x: Number(n.x1), y: Number(n.y1) }, { x: Number(n.x2), y: Number(n.y2) }];
      }
      if (n.x == null || n.w == null) return null;
      var cx = Number(n.x) + Number(n.w) / 2;
      var cy = Number(n.y) + (Number(n.h) || 26) / 2;
      var half = Number(n.w) / 2;
      if (Math.abs((Number(n.a) || 0) % 180) < 1) {
        return [{ x: cx - half, y: cy }, { x: cx + half, y: cy }];
      }
      return [{ x: cx, y: cy - half }, { x: cx, y: cy + half }];
    }

    function isFrameLine(n) {
      if (!n || n.t !== 'L') return false;
      if (n.net === 'frame') return true;
      return String(n.c || '').replace(/\s/g, '').toLowerCase() === '#8aa0b8';
    }

    /* 点在铭牌框内（字上、框内空白或框线）都算这只铭牌。 */
    function plateTextAt(ev) {
      var host = A.host();
      var page = A.page && A.page();
      var direct = ev.target.closest ? ev.target.closest('.hx-tx') : null;
      if (direct && !(direct.closest && direct.closest('.hx-edit-layer'))) return direct;
      if (!host || !page || !A.pageXY) return null;
      var pt = A.pageXY(ev.clientX, ev.clientY);
      if (!pt) return null;
      var nodes = page.nodes || [];
      var segs = [];
      var i;
      for (i = 0; i < nodes.length; i++) {
        if (!isFrameLine(nodes[i])) continue;
        var ends = lineEnds(nodes[i]);
        if (!ends) continue;
        segs.push({ x1: ends[0].x, y1: ends[0].y, x2: ends[1].x, y2: ends[1].y });
      }
      if (!segs.length) return null;
      var parent = [];
      for (i = 0; i < segs.length; i++) parent.push(i);
      function find(a) {
        while (parent[a] !== a) {
          parent[a] = parent[parent[a]];
          a = parent[a];
        }
        return a;
      }
      function unite(a, b) {
        a = find(a);
        b = find(b);
        if (a !== b) parent[b] = a;
      }
      var at = {};
      for (i = 0; i < segs.length; i++) {
        var s = segs[i];
        var c1 = (Math.round(s.x1 * 100) / 100) + ',' + (Math.round(s.y1 * 100) / 100);
        var c2 = (Math.round(s.x2 * 100) / 100) + ',' + (Math.round(s.y2 * 100) / 100);
        if (at[c1] != null) unite(i, at[c1]);
        else at[c1] = i;
        if (at[c2] != null) unite(i, at[c2]);
        else at[c2] = i;
      }
      var groups = {};
      for (i = 0; i < segs.length; i++) {
        var id = find(i);
        if (!groups[id]) groups[id] = [];
        groups[id].push(segs[i]);
      }
      var hit = null;
      var hitArea = Infinity;
      Object.keys(groups).forEach(function (id) {
        var g = groups[id];
        var x0 = Infinity;
        var y0 = Infinity;
        var x1 = -Infinity;
        var y1 = -Infinity;
        var k;
        for (k = 0; k < g.length; k++) {
          x0 = Math.min(x0, g[k].x1, g[k].x2);
          y0 = Math.min(y0, g[k].y1, g[k].y2);
          x1 = Math.max(x1, g[k].x1, g[k].x2);
          y1 = Math.max(y1, g[k].y1, g[k].y2);
        }
        var rw = x1 - x0;
        var rh = y1 - y0;
        if (!(rw >= 8 && rh >= 8 && rw <= 560 && rh <= 560)) return;
        if (pt.x < x0 - 3 || pt.x > x1 + 3 || pt.y < y0 - 3 || pt.y > y1 + 3) return;
        var area = rw * rh;
        if (area < hitArea) {
          hitArea = area;
          hit = { x0: x0, y0: y0, x1: x1, y1: y1 };
        }
      });
      if (!hit) return null;
      var best = null;
      var bestD = Infinity;
      host.querySelectorAll('.hx-tx').forEach(function (el) {
        if (el.closest && el.closest('.hx-edit-layer')) return;
        var tx = nodes[Number(el.getAttribute('data-i'))];
        if (!tx || tx.t !== 'T') return;
        var cx = Number(tx.x) + (Number(tx.w) || 0) / 2;
        var cy = Number(tx.y) + (Number(tx.h) || 0) / 2;
        if (cx < hit.x0 - 4 || cx > hit.x1 + 4 || cy < hit.y0 - 4 || cy > hit.y1 + 4) return;
        var d = Math.hypot(cx - (hit.x0 + hit.x1) / 2, cy - (hit.y0 + hit.y1) / 2);
        if (d < bestD) {
          bestD = d;
          best = el;
        }
      });
      return best;
    }

    /* 点在字上，或点在盖住铭牌的线上，都改这块铭牌的字。 */
    function textAt(ev) {
      var direct = plateTextAt(ev);
      if (direct) return direct;
      var host = A.host();
      if (!host) return null;
      var best = null;
      var bestArea = Infinity;
      host.querySelectorAll('.hx-tx').forEach(function (el) {
        if (el.closest && el.closest('.hx-edit-layer')) return;
        var r = el.getBoundingClientRect();
        if (ev.clientX < r.left || ev.clientX > r.right || ev.clientY < r.top || ev.clientY > r.bottom) return;
        var area = r.width * r.height;
        if (area < bestArea) {
          bestArea = area;
          best = el;
        }
      });
      return best;
    }

    /* 就地改铭牌上的字。Enter 记下，Shift+Enter 换行，Esc 取消。写进一次图仍要点保存。 */
    function startTextEdit(el) {
      if (!el || el.getAttribute('contenteditable') === 'true') return;
      var cur = (el.innerText || el.textContent || '').replace(/\r\n/g, '\n').replace(/\r/g, '\n');
      var done = false;
      el.classList.add('is-hx-texting');
      el.setAttribute('contenteditable', 'true');
      el.style.userSelect = 'text';
      el.style.webkitUserSelect = 'text';
      w.setTimeout(function () {
        if (done) return;
        el.focus();
        try {
          var sel = w.getSelection();
          var range = w.document.createRange();
          range.selectNodeContents(el);
          sel.removeAllRanges();
          sel.addRange(range);
        } catch (e) {}
      }, 0);
      function finish(save) {
        if (done) return;
        done = true;
        el.removeEventListener('blur', onBlur);
        el.removeEventListener('keydown', onKey);
        var next = save ? (el.innerText || '').replace(/\r\n/g, '\n').replace(/\r/g, '\n').replace(/\n$/, '') : cur;
        el.removeAttribute('contenteditable');
        el.classList.remove('is-hx-texting');
        if (!save || next === cur) {
          applyText(el, cur);
          return;
        }
        commitText(el, next);
      }
      function onBlur() { finish(true); }
      function onKey(e) {
        if (e.key === 'Escape') {
          e.preventDefault();
          e.stopPropagation();
          finish(false);
          return;
        }
        if (e.key === 'Enter' && !e.shiftKey) {
          e.preventDefault();
          e.stopPropagation();
          el.blur();
        }
      }
      el.addEventListener('blur', onBlur);
      el.addEventListener('keydown', onKey);
    }

    function onDbl(ev) {
      if (!A.isOn()) return;
      var el = textAt(ev);
      if (!el) return;
      ev.preventDefault();
      ev.stopPropagation();
      if (ev.stopImmediatePropagation) ev.stopImmediatePropagation();
      startTextEdit(el);
    }

    document.addEventListener('keydown', function (ev) {
      if (!A.isOn()) return;
      if (ev.key !== 'Delete' && ev.key !== 'Backspace') return;
      var t = ev.target;
      if (t && (t.tagName === 'INPUT' || t.tagName === 'TEXTAREA' || t.isContentEditable)) return;
      if (!A.sel()) return;
      ev.preventDefault();
      A.removeSelected();
    });

    return {
      onPalDown: onPalDown,
      onDown: onDown,
      onDbl: onDbl,
      cancel: function () { drag = null; hideMagnet(); disarm(); }
    };
  }

  w.HXEditGestures = { install: install };
})(window);
