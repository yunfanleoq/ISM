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
      if (ev.target.closest && ev.target.closest('.hx-zoom-bar, .hx-edit-pal')) return;
      /* 「局部」框选放大时让给舞台缩放，结束后仍回到编辑拖拽 */
      var host0 = A.host();
      var fit0 = host0 && host0.querySelector('.hx-fit');
      if (fit0 && fit0.classList.contains('is-marquee')) return;
      var handle = ev.target.closest && ev.target.closest('.hx-edit-handle');
      if (handle) {
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
      A.setSel({ t: 'base', i: i });
      A.markSel();
      armDrag(ev, {
        kind: 'base', i: i, n: n,
        cx: ev.clientX, cy: ev.clientY,
        dx: mv.dx || 0, dy: mv.dy || 0, s: mv.s,
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

    function shiftEnd(it, end, dx, dy) {
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
    }

    function onMove(ev) {
      if (!drag || ev.pointerId !== drag.pointer) return;
      var sc = A.scale();
      var dx = (ev.clientX - drag.cx) / sc;
      var dy = (ev.clientY - drag.cy) / sc;
      if (!drag.moved && Math.hypot(ev.clientX - drag.cx, ev.clientY - drag.cy) < 3) return;
      drag.moved = true;
      if (drag.kind === 'base') {
        A.elsFor(drag.i).forEach(function (el) {
          A.placeOne(el, drag.n, { dx: drag.dx + dx, dy: drag.dy + dy });
        });
        return;
      }
      var it = JSON.parse(JSON.stringify(drag.snap));
      if (drag.end === '1' || drag.end === '2' || drag.end === 'm') shiftEnd(it, drag.end, dx, dy);
      else shiftItem(it, dx, dy);
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
      disarm();
      if (!d.moved || !ev) return;
      var sc = A.scale();
      var dx = (ev.clientX - d.cx) / sc;
      var dy = (ev.clientY - d.cy) / sc;
      if (d.kind === 'base') {
        A.commit(function (ov) {
          var rec = { dx: A.r1(d.dx + dx), dy: A.r1(d.dy + dy) };
          if (d.s != null) rec.s = d.s;
          if (!rec.dx && !rec.dy && rec.s == null) delete ov.moves[d.i];
          else ov.moves[d.i] = rec;
        });
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

    function onDbl(ev) {
      if (!A.isOn()) return;
      var el = ev.target.closest ? ev.target.closest('.hx-tx') : null;
      if (!el) return;
      ev.preventDefault();
      ev.stopPropagation();
      /* 设备名牌等多行：用换行保留，prompt 里用 \n 表示 */
      var cur = (el.innerText || el.textContent || '').replace(/\r\n/g, '\n').replace(/\r/g, '\n');
      var next = w.prompt('设备名牌文字（换行输入 \\n）', cur.replace(/\n/g, '\\n'));
      if (next == null) return;
      next = String(next).replace(/\\n/g, '\n');
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
      cancel: function () { drag = null; disarm(); }
    };
  }

  w.HXEditGestures = { install: install };
})(window);
