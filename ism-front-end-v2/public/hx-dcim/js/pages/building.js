/* =========================================================
   空间巡检：按航信监控树（配电室 / 机房模块 / 数据机房）
   点进去用已有的配电一次图或列头柜，不读 DCIM.floors。
   ========================================================= */
(function (A) {
  'use strict';

  function D() { return window.DCIM; }

  /* 机柜 → 设备 索引（只建一次） */
  var CABIX = null;
  function cabIndex() {
    if (CABIX) return CABIX;
    CABIX = {};
    D().assets.forEach(function (a) {
      if (!a.cabinet) return;
      (CABIX[a.cabinet] = CABIX[a.cabinet] || []).push(a);
    });
    Object.keys(CABIX).forEach(function (k) {
      CABIX[k].sort(function (a, b) { return b.uStart - a.uStart; });
    });
    return CABIX;
  }

  var TYPE_COLOR = {
    '主机房': '#00D8FF', '屏蔽机房': '#00D8FF', '空调间': '#00E39A', '电池室': '#FFA24D',
    'UPS间': '#FFA24D', '配电室': '#2E9CF6', '弱电间': '#8FA8C4', '监测中心': '#B98CFF', '发电间': '#FF7A45'
  };
  function isMain(t) { return t === '主机房' || t === '屏蔽机房'; }
  function wOf(r) { return Math.max(2.2, r.cabCount > 0 ? r.cabCount : (isMain(r.type) ? 8 : 4)); }
  function roomNo(n) { var m = String(n).match(/(\d+)/); return m ? Number(m[1]) : 0; }
  function hotOf(r) { return r.alarmN >= 200; }

  /* ---------------- 页面 ---------------- */
  A.register('building', {
    render: function (root) {
      root.innerHTML = '<div id="hxBuilding" class="grow" style="min-height:0"></div>';
      var host = document.getElementById('hxBuilding');
      /* 默认进机房模块列头；架构树 + 面包屑 + 主区钻取 */
      var run = function () {
        if (window.HXScreen) window.HXScreen.mount(host, 'room');
      };
      if (window.HXScreen && window.HXScreen.ready) window.HXScreen.ready.then(run);
      else run();
    },

    /* 只读 hash：#/building/floor/4F/assetmgr 或 #/building/room/机房名/assetmgr
       有 assetmgr 段时，楼层这一级返回资产管理。
       有 dashboard 段时，楼层这一级返回总览，不先落到整栋楼。 */
    openFromHash: function () {
      var raw = (location.hash || '').replace(/^#\//, '');
      var parts = raw.split('/');
      var fromDash = parts.length && parts[parts.length - 1] === '~dashboard';
      if (fromDash) parts.pop();
      if (parts[0] !== 'building' || (parts[1] !== 'floor' && parts[1] !== 'room')) return false;
      var name = parts[2] || '';
      try { name = decodeURIComponent(name); } catch (e) {}
      if (!name) return false;
      var back = parts[3] === 'assetmgr' ? 'assetmgr' : (parts[3] === 'dashboard' ? 'dashboard' : '');
      if (back === 'assetmgr' && fromDash) back = 'assetmgr/~dashboard';
      if (parts[1] === 'floor') {
        var Dd = D();
        var ok = (Dd.floors || []).some(function (x) { return x.floor === name; })
          || (Dd.rooms || []).some(function (x) { return x.floor === name; });
        if (!ok) return false;
        this.entryBack = back;
        this.mode = '3d';
        this.viewFloor(name, null);
        return true;
      }
      var room = (D().rooms || []).filter(function (x) { return x.name === name; })[0];
      if (!room) return false;
      this.entryBack = back;
      this.mode = '3d';
      this.viewRoom(name, null);
      return true;
    },

    unmount: function () {
      if (this._esc) document.removeEventListener('keydown', this._esc);
      if (this._fs) document.removeEventListener('fullscreenchange', this._fs);
      if (window.A && window.A.V3) window.A.V3.dispose();
      if (window.HXScreen && window.HXScreen.leave) window.HXScreen.leave();
    },

    /* 配电室里的温感、STS、UPS 叶子和点表目录不是机房，不摆进巡检树。 */
    keepPlace: function (folder, room) {
      if (folder !== '配电室') return true;
      var n = room.name || '';
      if (/温感|STS/.test(n)) return false;
      if (/^配电室_/.test(n)) return false;
      if (/^配电室ECC_U/.test(n)) return false;
      return true;
    },

    siteFolders: function () {
      var cat = (window.HXScreen && window.HXScreen.catalog) ? window.HXScreen.catalog() : [];
      var self = this;
      return cat.map(function (f) {
        return {
          name: f.name,
          rooms: (f.rooms || []).filter(function (r) { return self.keepPlace(f.name, r); })
        };
      });
    },

    siteRoomFlat: function () {
      var flat = [];
      this.siteFolders().forEach(function (f) {
        (f.rooms || []).forEach(function (r) { flat.push(r); });
      });
      return flat;
    },

    renderTree: function () {
      var self = this;
      var arch = window.HXArch;
      if (arch && arch.renderTree && arch.branches && arch.branches.length) {
        self._flat = self.siteRoomFlat();
        var active = null;
        if (self.state.room) {
          var cur = self._flat.filter(function (x) { return x.name === self.state.room; })[0];
          active = {
            room: cur || { name: self.state.room },
            cab: (self.state.cab != null && cur && cur.cabs) ? cur.cabs[self.state.cab] : null
          };
        }
        var rendered = arch.renderTree(this.treeHost, {
          onPick: function (room, cab, item) {
            self.openByResolved(room, cab, item);
          },
          active: active
        });
        self._archFlat = rendered;
        return;
      }
      var folders = this.siteFolders();
      var flat = [];
      var html = folders.map(function (f) {
        if (!f.rooms.length) return '';
        return '<div class="hx-folder">' + A.esc(f.name) + ' · ' + f.rooms.length + '</div>' +
          f.rooms.map(function (r) {
            var i = flat.length;
            flat.push(r);
            var open = self.state.room === r.name;
            var headOn = open && self.state.cab == null;
            var head = '<button type="button" class="hx-node' + (headOn ? ' is-on' : '') + '" data-i="' + i + '">' + A.esc(r.name) + '</button>';
            var subs = '';
            if (open && r.cabs && r.cabs.length) {
              subs = r.cabs.map(function (c, ci) {
                var on = self.state.cab === ci ? ' is-on' : '';
                return '<button type="button" class="hx-node hx-sub' + on + '" data-i="' + i + '" data-c="' + ci + '">' + A.esc(c.name) + '</button>';
              }).join('');
            }
            return head + subs;
          }).join('');
      }).join('');
      this._flat = flat;
      this.treeHost.innerHTML = html || '<div class="hx-pt">现场索引里没有配电室或机房。</div>';
      this.treeHost.querySelectorAll('.hx-node').forEach(function (b) {
        b.addEventListener('click', function () {
          var ci = b.hasAttribute('data-c') ? Number(b.getAttribute('data-c')) : null;
          self.openSite(Number(b.getAttribute('data-i')), ci);
        });
      });
    },

    paintSite: function () {
      this.state.room = null;
      this.state.cab = null;
      if (window.HXScreen && window.HXScreen.leave) window.HXScreen.leave();
      this.renderTree();
      this.crumb.innerHTML = '<span class="crumb-item is-cur">现场空间</span>';
      var n = (this._flat || []).length;
      this.camHost.innerHTML =
        '<div class="hx-site-note">' +
          '<div class="hx-title">航信现场空间</div>' +
          '<div class="hx-pt">左侧按架构层级机构（配电室/机房 → 楼层 → 房间 → 列）。已绑测点的节点可点；灰显项尚未解开画面或未绑点。共 ' + n + ' 个可巡检房间。</div>' +
        '</div>';
    },

    openByResolved: function (room, cab, item) {
      var self = this;
      if (!window.HXScreen || !window.HXScreen.openIn) return;
      var ci = null;
      if (room && cab && room.cabs) {
        room.cabs.forEach(function (c, idx) { if (c.name === cab.name) ci = idx; });
      } else if (room && room.cabs && room.cabs.length) ci = 0;
      if (room) {
        this.state.room = room.name;
        this.state.cab = ci;
      }
      this.renderTree();
      var labelParts = window.HXArch && item
        ? window.HXArch.crumbLabels(item.ancestors, item.name, cab && cab.name)
        : [(item && item.name) || (room && room.name) || '节点'];
      this.crumb.innerHTML =
        '<span class="crumb-item" data-lv="0">现场空间</span><span class="crumb-sep">›</span>' +
        '<span class="crumb-item is-cur">' + A.esc(labelParts.join(' › ')) + '</span>';
      var resolved = item && item.resolved;
      if (resolved && !resolved.ok) {
        if (window.HXScreen.leave) window.HXScreen.leave();
        var links = resolved.links || [];
        var sug = resolved.kind === 'power' && window.HXArch ? window.HXArch.suggestFromPower(resolved.design) : null;
        var btns = links.map(function (lk, i) {
          return '<button type="button" class="hx-drill-go" data-link-i="' + i + '">前往 ' + A.esc(lk.label || lk.design) + '</button>';
        }).join('');
        if (sug && sug.item) {
          btns = '<button type="button" class="hx-drill-go is-primary" data-suggest="1">下钻 · ' +
            A.esc(sug.item.name) + '</button>' + btns;
        }
        this.camHost.innerHTML =
          '<div class="hx-site-note">' +
            '<div class="hx-title">' + A.esc(item.name || resolved.design || '') + '</div>' +
            '<div class="hx-pt">' + A.esc(resolved.reason || '暂无画面') + '。请点关联入口或左侧树继续。</div>' +
            (btns ? '<div class="hx-drill-actions">' + btns + '</div>' : '') +
          '</div>';
        this.camHost.querySelectorAll('[data-suggest]').forEach(function (btn) {
          btn.addEventListener('click', function () {
            if (sug && sug.item) self.openByResolved(sug.item.resolved.room, sug.item.resolved.cab, sug.item);
          });
        });
        this.camHost.querySelectorAll('[data-link-i]').forEach(function (btn) {
          btn.addEventListener('click', function () {
            var lk = links[Number(btn.getAttribute('data-link-i'))];
            if (!lk || !window.HXArch || !self._archFlat) return;
            var picked = window.HXArch.pickByDesign(self._archFlat, lk.design);
            if (picked) self.openByResolved(picked.resolved.room, picked.resolved.cab, picked);
          });
        });
        return;
      }
      if (!room) return;
      window.HXScreen.openIn(this.camHost, room, cab);
    },

    openSite: function (i, ci) {
      var room = (this._flat || [])[i];
      if (!room || !window.HXScreen || !window.HXScreen.openIn) return;
      if (ci == null && room.cabs && room.cabs.length) ci = 0;
      this.state.room = room.name;
      this.state.cab = ci;
      this.renderTree();
      var cab = ci == null ? null : (room.cabs || [])[ci];
      var label = room.name + (cab ? ' · ' + cab.name : '');
      this.crumb.innerHTML =
        '<span class="crumb-item" data-lv="0">现场空间</span><span class="crumb-sep">›</span>' +
        '<span class="crumb-item is-cur">' + A.esc(label) + '</span>';
      window.HXScreen.openIn(this.camHost, room, cab);
    },

    /* ---------------- 2D / 3D 切换 ---------------- */
    modeSw: function () {
      var m = this.mode;
      return '<div class="mode-sw">' +
        '<b class="' + (m === '2d' ? 'is-on' : '') + '" data-m="2d">2D 平面</b>' +
        '<b class="' + (m === '3d' ? 'is-on' : '') + '" data-m="3d">3D 立体</b>' +
      '</div>';
    },
    bindMode: function (nx) {
      var self = this;
      A.qa('.mode-sw b', nx).forEach(function (b) {
        b.addEventListener('click', function () {
          var m = b.getAttribute('data-m');
          if (m === self.mode) return;
          self.mode = m;
          self.redraw();
        });
      });
    },
    redraw: function () {
      var s = this.state;
      if (s.cab) this.viewCab(s.cab, null);
      else if (s.room) this.viewRoom(s.room, null);
      else if (s.floor) this.viewFloor(s.floor, null);
      else this.viewBuilding(null, true);
    },
    isoLegend: function (tip) {
      return '<div class="iso-legend">' +
        '<span><em class="lg is-ok"></em>正常 &lt;70%</span>' +
        '<span><em class="lg is-warn"></em>偏高 70~90%</span>' +
        '<span><em class="lg is-full"></em>满载 &gt;90%</span>' +
        '<span><em class="lg is-lhd"></em>列头柜</span>' +
        '<span><em class="lg is-cold"></em>冷通道</span>' +
        '<span><em class="lg is-warm"></em>热通道</span>' +
        '<span class="lg-tip">' + tip + '</span>' +
      '</div>';
    },

    /* ---------------- 镜头 ---------------- */

    /* 退场：动画结束或超时兜底都要摘掉，绝不让它留在 DOM 里 */
    _kill: function (el) {
      if (!el) return;
      el.removeAttribute('data-out');
      if (el.parentNode) el.remove();
    },

    fly: function (html, focus, dir, instant) {
      var cam = this.cam;
      /* 3D 查看器跟着这一层一起换：先拆掉上一个，rAF / 监听全部收回 */
      if (A.V3) A.V3.dispose();
      /* 上一张画面 = 我们显式记下的当前层。
         不能再用 querySelector 取「第一个」：快速连点时第一个是更早的旧层，
         真正正在显示的那一层就永远轮不到清理，一层层堆起来叠成一堆。 */
      var keep = cam._cur || null;
      var vs = cam.querySelectorAll('.cam-view, .cam-ring');
      for (var i = 0; i < vs.length; i++) {
        /* 除「上一张画面」外，任何残留（上一轮正在退场的、更早的）一律立即清掉 */
        if (vs[i] !== keep) vs[i].remove();
      }

      var nx = A.h('<div class="cam-view"></div>');
      nx.innerHTML = html;
      var ox = focus ? (focus.x * 100).toFixed(1) + '% ' + (focus.y * 100).toFixed(1) + '%' : '50% 50%';
      nx.style.transformOrigin = ox;
      cam._cur = nx;

      if (keep) {
        keep.style.transformOrigin = ox;
        keep.style.pointerEvents = 'none';
        keep.classList.remove('is-deep-in', 'is-back-in', 'is-instant');
        keep.classList.add(dir > 0 ? 'is-deep-out' : 'is-back-out');
        keep.setAttribute('data-out', '1');
        var self = this;
        keep.addEventListener('animationend', function (e) {
          if (e.target !== keep) return;
          self._kill(keep);
        });
        /* 兜底：动画被降级/未触发（隐藏标签页、减少动效）时也要清，不能只等 animationend */
        setTimeout(function () { self._kill(keep); }, 900);
      }
      /* 首屏：不进缩放动画，直接落位 + 淡入。
         进页面时整屏刚被 transform 缩放过，这时再套一层倍率缩放，
         整栋楼这种上百节点的 SVG 会被要求栅格化数倍面积，
         浏览器放大（125% 以上）时栅格预算不足就会留下黑块；
         而且这个黑块会跟着合成层一直留到下次重绘 —— 也就是「再点一次整栋楼才好」。 */
      if (instant) {
        nx.classList.add('is-instant');
        nx.addEventListener('animationend', function (e) {
          if (e.target !== nx) return;
          nx.classList.remove('is-instant');
          nx.style.opacity = '';
        });
        cam.appendChild(nx);
        A.repaint(cam);
        return nx;
      }
      nx.classList.add(dir > 0 ? 'is-deep-in' : 'is-back-in');
      /* 进场动画结束后立即摘掉 transform，避免长期留一个合成层导致重绘发黑 */
      nx.addEventListener('animationend', function (e) {
        if (e.target !== nx) return;
        nx.classList.remove('is-deep-in', 'is-back-in');
        nx.style.transform = '';
        nx.style.opacity = '';
        /* 动画期间这层是被提升过的合成层，落位后再压一次栅格，
           避免它带着动画期的旧栅格继续用（放大浏览时表现为黑块） */
        A.repaint(cam);
      });
      cam.appendChild(nx);
      if (focus) {
        var ring = A.h('<div class="cam-ring"></div>');
        ring.style.left = (focus.x * 100) + '%';
        ring.style.top = (focus.y * 100) + '%';
        cam.appendChild(ring);
        setTimeout(function () { if (ring.parentNode) ring.remove(); }, 720);
      }
      return nx;
    },

    focusOf: function (el) {
      var cam = this.cam;
      if (!el) return { x: .5, y: .5 };
      var a = cam.getBoundingClientRect(), b = el.getBoundingClientRect();
      return {
        x: A.clamp((b.left + b.width / 2 - a.left) / a.width, .06, .94),
        y: A.clamp((b.top + b.height / 2 - a.top) / a.height, .06, .94)
      };
    },

    /* ---------------- 面包屑 ---------------- */
    crumbs: function () {
      var s = this.state, out = [];
      var backTop = String(this.entryBack || '').split('/')[0];
      var rootLabel = backTop === 'assetmgr' ? '资产管理' : (backTop === 'dashboard' ? '总览' : '整栋楼');
      out.push('<span class="crumb-item' + (!s.floor ? ' is-cur' : '') + '" data-lv="0">' + rootLabel + '</span>');
      if (s.floor) out.push('<span class="crumb-sep">›</span><span class="crumb-item' + (!s.room ? ' is-cur' : '') + '" data-lv="1">' + s.floor + '</span>');
      if (s.room) out.push('<span class="crumb-sep">›</span><span class="crumb-item' + (!s.cab ? ' is-cur' : '') + '" data-lv="2">' + A.esc(s.room) + '</span>');
      if (s.cab) out.push('<span class="crumb-sep">›</span><span class="crumb-item is-cur" data-lv="3">' + A.esc(s.cab.split('-').slice(-2).join('-')) + ' 机柜</span>');
      this.crumb.innerHTML = out.join('');
    },

    jump: function (lv) {
      var s = this.state;
      if (lv === 0) {
        if (this.entryBack) { A.go(this.entryBack); return; }
        s.floor = s.room = s.cab = null; this.viewBuilding(this.cam);
      }
      else if (lv === 1) { s.room = s.cab = null; this.viewFloor(s.floor, this.cam); }
      else if (lv === 2) { s.cab = null; this.viewRoom(s.room, this.cam); }
    },
    up: function () {
      var s = this.state;
      if (s.cab) { s.cab = null; this.viewRoom(s.room, null); }
      else if (s.room) { s.room = null; this.viewFloor(s.floor, null); }
      else if (s.floor) {
        if (this.entryBack) { A.go(this.entryBack); return; }
        s.floor = null; this.viewBuilding(null);
      }
    },

    /* ---------------- L0 整栋楼 ---------------- */
    viewBuilding: function (fromEl, instant) {
      this.cleanChrome(false, true);
      var self = this, Dd = D();
      var floors = Dd.floors.slice().sort(function (a, b) { return b.floor < a.floor ? -1 : 1; });
      var maxAlarm = Math.max.apply(null, Dd.rooms.map(function (r) { return r.alarmN; }));

      var html = '<div class="bld-wrap"><div class="bld-stage2">' +
        window.BLD.svg({ interactive: true, uid: 'T' }) +
        '</div>' +
        '<div class="bld-legend">' +
          '<span><em style="background:rgba(0,190,255,.5)"></em>窗格 · 亮层为可选楼层</span>' +
          '<span><em style="background:rgba(0,118,188,.35)"></em>屋面 / 侧立面</span>' +
          '<span>悬停高亮 · 点击楼层向下钻取</span>' +
        '</div></div>';

      var nx = this.fly(html, this.focusOf(fromEl), fromEl ? -1 : 1, instant);

      /* 楼层热区：hover 显示该层概况，点击下钻 */
      A.qa('.bld-stage2 .hit', nx).forEach(function (el) {
        var fl = el.getAttribute('data-floor');
        var f = floors.filter(function (x) { return x.floor === fl; })[0];
        if (f) {
          var t = document.createElementNS('http://www.w3.org/2000/svg', 'title');
          t.textContent = f.floor + ' · ' + f.rooms + ' 间机房 · ' + f.cabCount + ' 架机柜 · 告警 ' + A.num(f.alarmN);
          el.appendChild(t);
        }
        el.addEventListener('click', function () {
          self.state.floor = fl;
          self.state.room = self.state.cab = null;
          self.viewFloor(self.state.floor, el);
        });
      });
      this.crumbs();
      this.sideBuilding();
    },

    sideBuilding: function () {
      var Dd = D();
      var top = Dd.rooms.slice().sort(function (a, b) { return b.alarmN - a.alarmN; }).slice(0, 6);
      this.side.innerHTML =
        A.UI.card({ title: '楼宇总览', body:
          A.UI.stats([
            { v: Dd.floors.length, u: '层', l: '楼层数' },
            { v: Dd.meta.roomTotal, u: '间', l: '机房数' },
            { v: A.num(Dd.meta.cabTotal), u: '架', l: '机柜总数' },
            { v: A.num(Dd.assetStats.total), u: '台', l: '在架设备' }
          ]) +
          '<div class="side-title" style="margin:14px 0 8px">各层机柜与告警</div>' +
          A.UI.bars(Dd.floors.slice().sort(function (a, b) { return a.floor < b.floor ? -1 : 1; }).map(function (f) {
            return { label: f.floor + ' · 告警 ' + A.num(f.alarmN), value: f.cabCount };
          }), { unit: '柜' }) +
          '<div class="side-title" style="margin:14px 0 8px">告警热点机房 TOP6</div>' +
          A.UI.rank(top.map(function (r) { return { name: r.name + '（' + r.floor + '）', value: r.alarmN }; }), ' 条')
        });
    },

    /* ---------------- L1 图纸平面 ---------------- */
    viewFloorCad: function (floor, fromEl, f, rs) {
      var self = this;
      var w = f.w || 48, d = f.d || 46;
      function Y(y) { return d - y; }
      var walls = (f.walls || []).map(function (seg) {
        return '<line x1="' + seg[0] + '" y1="' + Y(seg[1]) + '" x2="' + seg[2] + '" y2="' + Y(seg[3]) + '" class="cad-wall"/>';
      }).join('');
      var rooms = rs.filter(function (r) {
        return r.type === '主机房' || r.type === '配电室' || r.type === '电池室' || r.type === '操作间' || r.type === '空调间';
      }).map(function (r) {
        var p = r.poly;
        if (!p || p.length < 4) return '';
        var pts = p.map(function (q) { return q[0] + ',' + Y(q[1]); }).join(' ');
        return '<polygon class="cad-room" data-room="' + A.esc(r.name) + '" points="' + pts + '"/>' +
          '<text class="cad-label" x="' + ((p[0][0] + p[1][0]) / 2) + '" y="' + Y((p[0][1] + p[2][1]) / 2) + '">' +
          A.esc(r.code || r.name) + '</text>';
      }).join('');
      var cabs = '';
      rs.forEach(function (r) {
        (r.cabinets || []).forEach(function (c) {
          if (c.x == null) return;
          cabs += '<rect class="cad-cab" data-room="' + A.esc(r.name) + '" data-cab="' + A.esc(c.name) +
            '" x="' + c.x + '" y="' + Y(c.y + (c.d || 1)) + '" width="' + (c.w || 0.6) + '" height="' + (c.d || 1.2) + '"/>';
        });
      });
      var html = '<div class="fl-wrap"><div class="fl-head"><span class="fl-title">' + A.esc(floor) + ' 平面</span>' +
        '<span class="fl-sub">' + f.rooms + ' 间 · ' + f.cabCount + ' 柜 · 按图纸布置 · ' + A.simTag('sim') + '</span>' +
        this.modeSw() + '</div>' +
        '<div class="cad-stage"><svg class="cad-plan" viewBox="-1 -1 ' + (w + 2) + ' ' + (d + 2) + '">' +
        walls + rooms + cabs + '</svg></div></div>';
      var nx = this.fly(html, this.focusOf(fromEl), fromEl ? -1 : 1);
      this.bindMode(nx);
      A.qa('.cad-room,.cad-cab', nx).forEach(function (el) {
        el.addEventListener('click', function () {
          var room = el.getAttribute('data-room');
          var cab = el.getAttribute('data-cab');
          if (cab) { self.state.room = room; self.viewCab(cab, el); return; }
          self.viewRoom(room, el);
        });
      });
      this.crumbs();
    },

    /* ---------------- L1 楼层平面 ---------------- */
    viewFloor: function (floor, fromEl) {
      this.cleanChrome(false, true);
      var self = this, Dd = D();
      var f = Dd.floors.filter(function (x) { return x.floor === floor; })[0] || {};
      var rs = Dd.rooms.filter(function (r) { return r.floor === floor; });
      this.state.floor = floor; this.state.room = null; this.state.cab = null;
      if (f.walls && f.walls.length && this.mode !== '3d') return this.viewFloorCad(floor, fromEl, f, rs);

      /* 机柜排布示意：按真实列 / 冷热通道画缩略阵列 */
      function racks(r) {
        if (!r.cabCount) return '';
        var lay = r.layout, cabs = r.cabinets, s = '', ci, k, list;
        if (lay) {
          var ai = 0;
          for (ci = 0; ci < lay.cols; ci++) {
            list = cabs.filter(function (c) { return c.col === ci + 1; });
            s += '<div class="rk-row">';
            for (k = 0; k < list.length; k++) {
              var lv = list[k].lv || (!list[k].dev ? 'off' : list[k].comp > 90 ? 'red' : list[k].comp >= 70 ? 'yellow' : 'blue');
              s += '<i class="is-' + lv + '"></i>';
            }
            s += '</div>';
            if (ci < lay.cols - 1) {
              var as = (lay.aisles || [])[ai++] || { t: ci % 2 === 0 ? 'cold' : 'hot' };
              s += '<div class="rk-aisle is-' + as.t + '"></div>';
            }
          }
        } else {
          var COLS = Math.max(4, Math.min(12, r.cabCount)), rows = Math.min(Math.ceil(r.cabCount / COLS), 8), n = 0;
          for (ci = 0; ci < rows; ci++) {
            s += '<div class="rk-row">';
            for (k = 0; k < COLS; k++) {
              n++;
              s += '<i class="' + ((((n * 7 + r.cabCount) % 13) < (r.rate / 9)) ? 'is-hot' : '') + '"></i>';
            }
            s += '</div>';
          }
        }
        return '<div class="fl-room-racks">' + s + '</div>';
      }

      function zone(z) {        var list = rs.filter(function (r) { return r.zone === z; }).sort(function (a, b) { return roomNo(a.name) - roomNo(b.name); });
        if (!list.length) return '';
        return '<div class="fl-zone is-' + z.toLowerCase() + '">' + list.map(function (r) {
          return '<div class="fl-room' + (isMain(r.type) ? ' is-main' : '') + (hotOf(r) ? ' is-hot' : '') + '" data-room="' + A.esc(r.name) + '" style="--w:' + wOf(r).toFixed(1) + '">' +
            '<div class="fl-room-nm">' + A.esc(r.name) + '</div>' +
            '<div class="fl-room-tp">' + A.esc(r.type) + ' · ' + r.cabCount + ' 机柜</div>' +
            '<div class="fl-room-kv"><span>设备 <b>' + A.num(r.dev) + '</b></span><span>告警 <b>' + A.num(r.alarmN) + '</b></span>' +
              (r.cabCount ? '<span>综合 <b>' + A.fix(r.comp, 0) + '%</b></span>' : '') + '</div>' +
            racks(r) +
            (r.cabCount ? '<div class="fl-room-bar is-lv-' + (r.comp > 90 ? 'red' : r.comp >= 70 ? 'yellow' : 'blue') +
              '"><i style="width:' + A.clamp(r.comp, 3, 100) + '%"></i></div>' : '') +
          '</div>';
        }).join('') + '</div>';
      }

      var is3 = this.mode === '3d' && A.V3;
      var head = '<div class="fl-head"><span class="fl-title">' + floor + (is3 ? ' 立体' : ' 平面') + '</span>' +
          '<span class="fl-sub">' + (f.rooms || 0) + ' 间机房 · ' + (f.cabCount || 0) + ' 架机柜 · ' +
          A.num(f.alarmN || 0) + ' 条告警 · 总负载 ' + A.fix(f.kw || 0, 1) + ' kW</span>' + this.modeSw() + '</div>';

      var body = is3
        ? '<div class="iso-stage" id="v3Host"></div>' +
          this.isoLegend('拖拽旋转 · 滚轮缩放 · WASD 漫游 · 双击飞向房间 · 点击房间向下钻取')
        : '<div class="fl-body">' + zone('N') + '<div class="fl-corridor">走 廊 · COLD AISLE</div>' + zone('S') + '</div>';

      var html = '<div class="fl-wrap">' + head + body + '</div>';

      var nx = this.fly(html, this.focusOf(fromEl), fromEl ? 1 : 1);
      if (is3) {
        A.V3.mount(A.q('#v3Host', nx), 'floor', floor, {
          drill: function (name) {
            self.state.room = name; self.state.cab = null;
            self.viewRoom(name, null);
          },
          toast: A.toast
        });
      }
      A.qa('.fl-room, .iso-room', nx).forEach(function (el) {
        el.addEventListener('click', function () {
          self.state.room = el.getAttribute('data-room');
          self.state.cab = null;
          self.viewRoom(self.state.room, el);
        });
      });
      this.bindMode(nx);
      this.crumbs();
      this.sideFloor(f, rs);
    },

    sideFloor: function (f, rs) {
      var Dd = D();
      var self = this;
      this.side.innerHTML =
        A.UI.card({ title: f.floor + ' 楼层概况', body:
          A.UI.stats([
            { v: f.rooms, u: '间', l: '机房' },
            { v: f.cabCount, u: '架', l: '机柜' },
            { v: A.num(f.dev), u: '台', l: '设备' },
            { v: A.num(f.alarmN), u: '条', l: '告警', color: f.alarmN > 1000 ? 'var(--red)' : '' }
          ]) +
          '<div class="side-title" style="margin:14px 0 8px">机房综合使用率（取极高）</div>' +
          A.UI.bars(rs.filter(function (r) { return r.cabCount; }).map(function (r) {
            return { label: r.name, value: r.comp, cls: r.comp > 90 ? 'o' : '' };
          }), { max: 100, unit: '%' }) +
          '<div class="side-title" style="margin:14px 0 8px">机房清单</div>' +
          A.UI.table({
            cols: [
              { t: '机房', k: 'name' },
              { t: '类型', k: 'type', w: 74 },
              { t: '柜', k: 'cabCount', align: 'r', w: 42 },
              { t: '告警', k: 'alarmN', align: 'r', w: 58, fmt: function (r) { return A.num(r.alarmN); } }
            ],
            rows: rs.slice().sort(function (a, b) { return b.cabCount - a.cabCount; }),
            rowAttr: function (r) { return 'data-room="' + A.esc(r.name) + '" style="cursor:pointer"'; }
          })
        });
      /* 清单行也要能点：标了 cursor:pointer 却绑不上事件，
         用户看着"机房清单"却点不动，只能回头去画布上找色块。 */
      A.qa('[data-room]', this.side).forEach(function (el) {
        el.addEventListener('click', function () {
          var nm = el.getAttribute('data-room');
          if (!nm) return;
          self.state.room = nm; self.state.cab = null;
          self.viewRoom(nm, null);
        });
      });
    },

    /* ---------------- L2 机房 ---------------- */
    viewRoom: function (name, fromEl) {
      this.cleanChrome(false, true);
      var self = this;
      var hint0 = document.getElementById('tourHint');
      if (hint0) hint0.innerHTML = '点击色块沿镜头向下钻取 · <b>Esc</b> 返回上一层';
      var r = D().rooms.filter(function (x) { return x.name === name; })[0];
      if (!r) return;
      this.state.room = name; this.state.floor = r.floor; this.state.cab = null;

      var html;
      if (!r.cabCount) {
        html = '<div class="rm-wrap"><div class="rm-head"><span class="rm-title">' + A.esc(r.name) + '</span>' +
          '<span class="fl-sub">' + A.esc(r.type) + ' · 无机柜 · ' + r.dev + ' 台独立设备 · 告警 ' + A.num(r.alarmN) + '</span>' +
            this.modeSw() + '</div>' +
          '<div class="rm-grid-wrap">' +
          A.UI.table({
            cols: [{ t: '设备', k: 'sn' }, { t: '类别', k: 'cat' }, { t: '品牌', k: 'brand' }, { t: '型号', k: 'model' }],
            rows: D().assets.filter(function (a) { return a.room === name; })
          }) + '</div></div>';
      } else {
        var lay = r.layout || null;
        var ix = cabIndex();
        var cols = [];
        if (lay) {
          for (var ci = 0; ci < lay.cols; ci++) {
            cols.push(r.cabinets.filter(function (c) { return c.col === ci + 1; })
              .sort(function (a, b) { return a.seq - b.seq; }));
          }
        } else {
          var byCol = {};
          r.cabinets.forEach(function (c) {
            var ck = c.col || 0;
            if (!ck && c.name) {
              var p = c.name.split('-');
              ck = Number(p[p.length - 2]);
            }
            (byCol[ck] = byCol[ck] || []).push(c);
          });
          Object.keys(byCol).sort(function (a, b) { return Number(a) - Number(b); }).forEach(function (ck) {
            cols.push(byCol[ck].sort(function (a, b) {
              return (a.seq || 0) - (b.seq || 0);
            }));
          });
        }

        /* 机柜格：色块 = 综合使用率四级色阶（灰未启用 / 蓝<70 / 黄70-90 / 红>90） */
        function cabCell(c) {
          var label = c.id || (c.col + '-' + c.seq);
          var items = ix[c.name] || [];
          if (!c.dev && !c.id) {
            return '<div class="cab2 is-empty" data-cab="' + A.esc(c.name) + '" title="' +
              A.esc(c.name) + ' · 未启用 · 可用 42U · 可部署 2U×' + Math.floor(42 / 3) + '">' +
              '<div class="cab2-nm">' + label + '</div></div>';
          }
          var lv = c.lv || (c.comp > 90 ? 'red' : c.comp >= 70 ? 'yellow' : 'blue');
          var lvName = lv === 'red' ? '满载' : lv === 'yellow' ? '偏高' : '正常';
          var tip = c.id
            ? A.esc(c.id) + ' · ' + (c.kind || '') + ' mm · 模拟'
            : A.esc(c.name) + ' · 综合 ' + c.comp + '%（' + lvName + '）' +
            '\n空间 ' + c.spaceRate + '%（' + c.usedU + 'U + ' + c.dev + '×' + c.coef + '系数 / 42U）' +
            '\n电力 ' + c.powRate + '%（' + c.rtKw + ' / ' + c.rated + 'kW）' +
            '\n接口 ' + c.ifRate + '%（' + c.ifUsed + ' / 20 口）' +
            '\n平均 ' + c.avgRate + '% · 峰值 ' + c.peakRate + '% · ' + c.dev + ' 台设备';
          return '<div class="cab2 lv-' + lv + '" data-cab="' + A.esc(c.name) + '" title="' + tip + '">' +
            '<div class="cab2-fill lv-' + lv + '" style="height:' + A.clamp(c.comp, 4, 100) + '%"></div>' +
            '<div class="cab2-nm">' + label + '</div>' +
            '<div class="cab2-rate">' + (c.id ? '' : Math.round(c.comp) + '%') + '</div></div>';
        }

        /* 列 + 列间通道（冷热交替），列头放列头柜 */
        var plan = [], ai = 0;
        if (lay) {
          for (var k = 0; k < cols.length; k++) {
            var list = cols[k];
            var hotCol = list.filter(function (c) { return c.rate > 85; }).length;
            plan.push(
              '<div class="rm-lane' + (hotCol ? ' is-hot' : '') + '">' +
                '<div class="rm-lhd" title="列头柜 · A/B 双路供电 · 列头配电">' +
                  '<b>L' + (k + 1) + '</b><span>列头柜</span></div>' +
                '<div class="rm-cabs">' + list.map(cabCell).join('') +
                  (list.length < 16 ? '<div class="rm-reserve">预留扩容区</div>' : '') + '</div>' +
                '<div class="rm-lane-nm"><b>列 ' + (k + 1) + '</b><span>' + list.length + ' 柜</span></div>' +
              '</div>'
            );
            if (k < cols.length - 1) {
              var as = (lay.aisles || [])[ai++] || { t: (k % 2 === 0 ? 'cold' : 'hot') };
              var cold = as.t === 'cold';
              plan.push(
                '<div class="rm-aisle is-' + as.t + '">' +
                  '<span class="rm-aisle-tag">' + (cold ? '冷通道 COLD AISLE' : '热通道 HOT AISLE') + '</span>' +
                  '<span class="rm-aisle-v">' + (cold ? '22.5℃ · 39.4%RH' : '33.6℃ · 28.0%RH') + '</span>' +
                  '<span class="rm-aisle-x">列 ' + (k + 1) + ' ⇄ 列 ' + (k + 2) + '</span>' +
                '</div>'
              );
            }
          }
        } else {
          cols.forEach(function (list) {
            var letter = (list[0] && list[0].id && list[0].id.split('-')[0]) || '列';
            plan.push(
              '<div class="rm-lane">' +
                '<div class="rm-lhd"><b>' + A.esc(letter) + '</b><span>列</span></div>' +
                '<div class="rm-cabs">' + list.map(cabCell).join('') + '</div>' +
                '<div class="rm-lane-nm"><b>列 ' + A.esc(letter) + '</b><span>' + list.length + ' 柜</span></div>' +
              '</div>'
            );
          });
        }

        /* 房间周边设备：空调机组 / 送风来源 / 漏水 / 温湿度 / 配电 */
        var acN = lay ? (lay.acIn || 0) : 0;
        var acRooms = lay && lay.acRoom && lay.acRoom.length ? lay.acRoom : null;
        var topEq = '';
        for (var q = 0; q < acN; q++) {
          topEq += '<i class="rm-eq is-ac">AC-' + (q + 1) + ' 精密空调</i>';
        }
        if (!acN) {
          var acTxt = acRooms ? acRooms.map(function (x) { return String(x).replace('空调间', ''); }).join(' / ') + ' 空调间' : '楼层空调间';
          topEq += '<i class="rm-eq is-ac-out">↗ 冷源：' + A.esc(acTxt) + ' 送风</i>';
        }
        if (lay && lay.upsRoom) topEq += '<i class="rm-eq is-ups">UPS：' + A.esc(lay.upsRoom) + '</i>';
        if (lay && lay.batRoom) topEq += '<i class="rm-eq is-ups">电池：' + A.esc(lay.batRoom) + '</i>';

        var botEq = '<i class="rm-eq is-leak">漏水监测 ×' + (lay ? lay.leak : 2) + ' 回路</i>' +
          '<i class="rm-eq is-th">温湿度 ×' + (lay ? lay.th : 4) + ' 点</i>' +
          '<i class="rm-eq is-pdu">列头配电 ' + (lay ? lay.lhd * 2 : 2) + ' 路（A/B）</i>';
        if (!lay) {
          topEq = '<i class="rm-eq is-ac-out">列号来自平面图</i>';
          botEq = '<i class="rm-eq is-th">柜内没有图上的服务器，功率和温度是模拟</i>';
        }

        html = '<div class="rm-wrap">' +
          '<div class="rm-head"><span class="rm-title">' + A.esc(r.name) + '</span>' +
            '<span class="fl-sub">' + A.esc(r.type) + ' · ' +
            (lay ? lay.cols + ' 列 × ' + (cols[0] ? cols[0].length : 0) + ' 柜' : '') + ' · ' +
            r.cabCount + ' 架机柜 · ' + A.num(r.dev) + ' 台设备 · ' +
            A.fix(r.kw, 1) + ' kW · 综合 ' + A.fix(r.comp, 1) + '% · 告警 ' + A.num(r.alarmN) + '</span>' +
            this.modeSw() + '</div>' +
          '<div class="rm-plan">' +
            '<div class="rm-wall">' + topEq + '</div>' +
            '<div class="rm-cols">' + plan.join('') + '</div>' +
            '<div class="rm-wall is-bot">' + botEq + '</div>' +
          '</div>' +
          '<div class="rm-legend">' +
            '<span><em class="lg is-empty"></em>未启用</span>' +
            '<span><em class="lg is-ok"></em>正常 &lt;70%</span>' +
            '<span><em class="lg is-warn"></em>偏高 70~90%</span>' +
            '<span><em class="lg is-full"></em>满载 &gt;90%</span>' +
            '<span><em class="lg is-lhd"></em>列头柜</span>' +
            '<span><em class="lg is-cold"></em>冷通道</span>' +
            '<span><em class="lg is-warm"></em>热通道</span>' +
            '<span class="lg-tip">色块高度 = 综合使用率（空间 / 电力 / 接口 取极高）· 悬停看四项明细 · 点击看 42U 排布</span>' +
          '</div></div>';
      }

      if (this.mode === '3d' && A.V3) {
        html = '<div class="rm-wrap">' +
          '<div class="rm-head"><span class="rm-title">' + A.esc(r.name) + '</span>' +
            '<span class="fl-sub">' + A.esc(r.type) + ' · ' +
            (r.layout ? r.layout.cols + ' 列' : '无机柜') + ' · ' + r.cabCount + ' 架机柜 · ' +
            A.num(r.dev) + ' 台设备 · ' + A.fix(r.kw, 1) + ' kW · 综合 ' + A.fix(r.comp, 1) + '% · 告警 ' +
            A.num(r.alarmN) + '</span>' + this.modeSw() + '</div>' +
          '<div class="iso-stage" id="v3Host"></div>' +
          this.isoLegend(!r.cabCount && r.type === '弱电间'
            ? '拖拽旋转 · 滚轮缩放 · 点击配线柜，看到的是和楼层同一套柜'
            : '拖拽旋转 · 滚轮缩放 · WASD 漫游 · 双击飞向机柜 · 点击机柜查看 42U 排布') +
        '</div>';
      }

      var nx = this.fly(html, this.focusOf(fromEl), fromEl ? 1 : 1);
      if (this.mode === '3d' && A.V3) {
        A.V3.mount(A.q('#v3Host', nx), 'room', r.name, {
          drill: function (name) { self.state.cab = name; self.viewCab(name, null); },
          toast: A.toast
        });
      }
      A.qa('.cab2, .iso-cab', nx).forEach(function (el) {
        el.addEventListener('click', function () {
          self.state.cab = el.getAttribute('data-cab');
          self.viewCab(self.state.cab, el);
        });
      });
      A.qa('.iso-dev', nx).forEach(function (el) {
        el.addEventListener('click', function () {
          var sn = el.getAttribute('data-sn');
          var a = D().assets.filter(function (x) { return x.sn === sn; })[0];
          if (a) A.toast(a.brand + ' ' + a.model + ' · ' + a.uH + 'U · U' + a.uStart + '-' + a.uEnd + ' · ' + a.kw + 'kW');
        });
      });
      this.bindMode(nx);
      this.crumbs();
      this.sideRoom(r);
    },

    sideRoom: function (r) {
      var top = r.cabinets.filter(function (c) { return c.dev; })
        .sort(function (a, b) { return b.comp - a.comp; }).slice(0, 8);
      this.side.innerHTML =
        A.UI.card({ title: A.esc(r.name), body:
          A.UI.stats([
            { v: r.cabCount, u: '架', l: '机柜' },
            { v: A.num(r.dev), u: '台', l: '设备' },
            { v: A.fix(r.comp, 1) + '%', l: '综合使用率',
              color: r.comp > 90 ? 'var(--lv1)' : r.comp >= 70 ? 'var(--orange)' : 'var(--cyan)' },
            { v: A.num(r.alarmN), u: '条', l: '告警', color: r.alarmN > 500 ? 'var(--red)' : '' }
          ]) +
          A.UI.bars([
            { label: '空间', value: r.spaceRate, cls: r.spaceRate > 90 ? 'o' : '' },
            { label: '电力', value: r.powRate, cls: r.powRate > 90 ? 'o' : '' },
            { label: '接口', value: r.ifRate, cls: r.ifRate > 90 ? 'o' : '' },
            { label: '平均', value: r.avgRate },
            { label: '峰值', value: r.peakRate, cls: r.peakRate > 90 ? 'o' : '' }
          ], { max: 100, unit: '%' }) +
          '<div style="margin-top:12px">' + A.UI.kv([
            ['所属楼层', r.floor], ['机房类型', r.type],
            ['机柜排布', r.layout ? (r.layout.cols + ' 列 × ' + r.layout.per.join('/') + ' 柜') : '-'],
            ['冷/热通道', r.layout ? (r.layout.aisles.filter(function (a) { return a.t === 'cold'; }).length + ' 冷 / ' +
              r.layout.aisles.filter(function (a) { return a.t === 'hot'; }).length + ' 热') : '-'],
            ['精密空调', r.layout ? (r.layout.acIn ? (r.layout.acIn + ' 台（机房内）') : (r.layout.acRoom.join('、') || '楼层空调间送风')) : '-'],
            ['列头柜', r.layout ? (r.layout.lhd + ' 台 · A/B 双路') : '-'],
            ['空闲机柜', r.cabinets.filter(function (c) { return !c.dev; }).length + ' 架'],
            ['额定负载', A.fix(r.kw, 1) + ' kW'],
            ['平均温度', r.temp + ' ℃'], ['平均湿度', r.humi + ' %'],
            ['最近告警', r.lastAlarm || '-']
          ]) + '</div>' +
          (r.cats && r.cats.length ? '<div class="side-title" style="margin:14px 0 8px">设备类别构成</div>' +
            '<div style="display:flex;gap:14px;align-items:center">' +
              A.Charts.donut({ size: 104, r: 40, sw: 13, data: r.cats.map(function (c, i) {
                return { value: c.value, color: ['#00D8FF', '#00E39A', '#FFA24D', '#B98CFF', '#2E9CF6', '#8FA8C4'][i % 6] };
              }) }) +
              '<div style="flex:1 1 auto;min-width:0">' + A.UI.bars(r.cats.map(function (c) {
                return { label: c.name, value: c.value };
              }), { unit: '台' }) + '</div></div>' : '') +
          (top.length ? '<div class="side-title" style="margin:14px 0 8px">使用率最高机柜</div>' +
            A.UI.rank(top.map(function (c) {
              return { name: c.name.split('-').slice(-2).join('-') + ' 柜', value: c.comp };
            }), '%') : '')
        });
    },

    /* ---------------- L3 机柜 ---------------- */
    cabParts: function (cabName) {
      var weak = String(cabName || '').match(/^弱电柜:(.+):(\d+)-(\d+)$/);
      if (weak) {
        var wroom = D().rooms.filter(function (x) { return x.name === weak[1]; })[0];
        var wlabel = '配线柜 ' + weak[2] + '-' + weak[3];
        return {
          cab: null, room: wroom, items: [],
          rows: '<div class="rk-u"><span class="rk-unm">' + A.esc(wlabel) + ' · 与楼层同一套配线柜</span></div>',
          cap: '<b>' + A.esc(wlabel) + '</b><span>弱电配线柜 · 无服务器 U 位门</span>',
          right: A.UI.card({ title: wlabel, body: A.UI.kv([
            ['所属机房', wroom ? wroom.name + '（' + wroom.floor + '）' : weak[1]],
            ['外形', '弱电配线柜'],
            ['对应', '楼层和房间里的同一架']
          ]) })
        };
      }
      var cab = null, room = null;
      D().rooms.forEach(function (x) {
        x.cabinets.forEach(function (c) { if (c.name === cabName) { cab = c; room = x; } });
      });
      var items = cabIndex()[cabName] || [];
      var byU = {};
      items.forEach(function (a) { byU[a.uStart] = a; });
      var rows = [], skip = {};
      items.forEach(function (a) {
        var top = Math.min(42, a.uEnd);
        for (var k = a.uStart + 1; k <= top; k++) skip[k] = 1;
      });
      for (var u = 42; u >= 1; u--) {
        if (skip[u]) continue;
        var d = byU[u];
        var grow = d ? Math.max(1, Math.min(42, d.uH)) : 1;
        var cls = !d ? '' : (d.cat === '服务器' ? ' is-srv' : d.cat === '存储设备' ? ' is-sto' : (d.cat === '交换机' || d.cat === '网络安全设备') ? ' is-net' : ' is-oth');
        var un = d && d.uH > 1 ? (A.pad(u) + '-' + A.pad(Math.min(42, d.uEnd))) : A.pad(u);
        rows.push('<div class="rk-u' + cls + '" style="flex-grow:' + grow + '"' + (d ? ' data-sn="' + A.esc(d.sn) + '"' : '') + '>' +
          '<span class="rk-uno">' + un + '</span>' +
          '<span class="rk-unm">' + (d ? A.esc(d.brand + ' ' + d.model) + ' <b style="opacity:.6">' + d.uH + 'U</b>' : '空位') + '</span>' +
          '<span class="rk-ukw">' + (d ? d.kw + 'kW' : '') + '</span></div>');
      }
      var meta = window.ISO ? window.ISO.rackMeta(cabName) : null;
      var cap = meta
        ? '<b>' + A.esc(meta.name) + ' 机柜</b><span>' + meta.n + ' 台' +
          (meta.spec ? ' · ' + meta.spec : '（空柜）') + '</span>' +
          (meta.over ? '<em>⚠ ' + meta.over + ' 台 U 位冲突待调整</em>' : '')
        : '';
      var right = A.UI.card({ title: A.esc(cabName), body:
        A.UI.stats([
          { v: cab ? cab.dev : 0, u: '台', l: '在架设备' },
          { v: cab ? cab.usedU : 0, u: 'U', l: '已用 U 位' },
          { v: cab ? cab.comp : 0, u: '%', l: '综合使用率',
            color: cab ? (cab.comp > 90 ? 'var(--lv1)' : cab.comp >= 70 ? 'var(--orange)' : 'var(--cyan)') : '' },
          { v: cab ? A.fix(cab.rtKw, 2) : 0, u: 'kW', l: '实测功率' }
        ]) +
        (cab && cab.dev ? A.UI.bars([
          { label: '空间', value: cab.spaceRate, cls: cab.spaceRate > 90 ? 'o' : '' },
          { label: '电力', value: cab.powRate, cls: cab.powRate > 90 ? 'o' : '' },
          { label: '接口', value: cab.ifRate, cls: cab.ifRate > 90 ? 'o' : '' },
          { label: '平均', value: cab.avgRate },
          { label: '峰值', value: cab.peakRate, cls: cab.peakRate > 90 ? 'o' : '' }
        ], { max: 100, unit: '%' }) : '') +
        '<div style="margin-top:12px">' + A.UI.kv([
          ['所属机房', room ? room.name + '（' + room.floor + '）' : '-'],
          ['机柜编号', cabName],
          ['主要类别', cab ? cab.topCat : '-'],
          ['空 U 系数', cab ? cab.coef : '-'],
          ['额定供电', cab ? cab.rated + ' kW' : '-'],
          ['空间利用率', cab ? cab.spaceRate + '%（' + cab.usedU + 'U + ' + cab.dev + '×' + cab.coef + ' / 42U）' : '-'],
          ['剩余 U 位', (42 - (cab ? cab.usedU : 0)) + ' U']
        ]) + '</div>'
      }) +
      A.UI.card({ title: '设备清单（自上而下）', body:
        '<div style="flex:1 1 auto;min-height:0;overflow:auto">' +
        A.UI.table({
          cols: [
            { t: 'U', k: 'uStart', w: 40, align: 'r' },
            { t: '类别', k: 'cat', w: 78 },
            { t: '品牌 / 型号', k: 'x', html: function (a) { return A.esc(a.brand + ' ' + a.model); } },
            { t: '功率', k: 'kw', align: 'r', w: 56 }
          ],
          rows: items
        }) + '</div>'
      });
      return { cab: cab, room: room, items: items, rows: rows.join(''), cap: cap, right: right };
    },
    syncCabPanel: function (cabName) {
      var nx = this.cam && this.cam._cur;
      if (!nx) return;
      var p = this.cabParts(cabName);
      var rack = nx.querySelector('.rk-rack');
      if (rack) rack.innerHTML = p.rows;
      var right = nx.querySelector('.rk-right');
      if (right) right.innerHTML = p.right;
      var cap = nx.querySelector('.rk-cap');
      if (cap) cap.innerHTML = p.cap;
      this.state.cab = cabName;
      if (p.room) { this.state.room = p.room.name; this.state.floor = p.room.floor; }
      this.crumbs();
    },
    markDev: function (sn) {
      var nx = this.cam && this.cam._cur;
      if (!nx) return;
      A.qa('.rk-u.is-sel', nx).forEach(function (x) { x.classList.remove('is-sel'); });
      A.qa('.rk-u[data-sn]', nx).forEach(function (el) {
        if (el.getAttribute('data-sn') === sn) el.classList.add('is-sel');
      });
    },
    rackMode: function () {
      try { return sessionStorage.getItem('bjga.rack.mode') === 'list' ? 'list' : 'patrol'; }
      catch (e) { return 'patrol'; }
    },
    rackCleanOn: function () {
      try { return sessionStorage.getItem('bjga.rack.clean') === '1'; }
      catch (e) { return false; }
    },
    setFloats: function (nx, devOn, infoOn) {
      if (!nx) return;
      var dev = nx.querySelector('[data-float="dev"]');
      var info = nx.querySelector('[data-float="info"]');
      if (dev) dev.classList.toggle('is-off', !devOn);
      if (info) info.classList.toggle('is-off', !infoOn);
      var bd = nx.querySelector('[data-act="dev"]');
      var bi = nx.querySelector('[data-act="info"]');
      if (bd) bd.classList.toggle('is-on', !!devOn);
      if (bi) bi.classList.toggle('is-on', !!infoOn);
    },
    /* keep=true：只摘掉铺满样式，不改会话里记下的纯净巡视。 */
    cleanChrome: function (on, keep) {
      var tour = this.cam && this.cam.closest ? this.cam.closest('.tour') : document.querySelector('.tour');
      if (tour) tour.classList.toggle('is-clean', !!on);
      if (!keep) {
        try { sessionStorage.setItem('bjga.rack.clean', on ? '1' : '0'); } catch (e) { /* 会话不可用就只改当前页 */ }
      }
      var nx = this.cam && this.cam._cur;
      var b = nx && nx.querySelector('[data-act="room"]');
      if (b) b.textContent = on ? '显示机房信息' : '隐藏机房信息';
      if (on && nx) this.setFloats(nx, false, false);
      if (!on && nx && nx.querySelector('.rk-wrap.is-3d')) {
        var list = this.rackMode() === 'list';
        this.setFloats(nx, list, list);
      }
      var v = A.V3 && A.V3.active;
      if (v && v.fitRack) v.fitRack();
    },
    viewCab: function (cabName, fromEl) {
      var self = this;
      var cab = null, room = null;
      D().rooms.forEach(function (x) {
        x.cabinets.forEach(function (c) { if (c.name === cabName) { cab = c; room = x; } });
      });
      this.state.cab = cabName;
      if (room) { this.state.room = room.name; this.state.floor = room.floor; }
      var part = this.cabParts(cabName);

      var is3 = this.mode === '3d' && A.V3;
      var weakCab = String(cabName).indexOf('弱电柜:') === 0;
      var lhead = '<div class="rk-lhead"><span class="side-title">' +
        (weakCab ? '配线柜' : (is3 ? '42U 立体视图' : '42U 正面排布')) + '</span>' + this.modeSw() + '</div>';
      var cap = part.cap ? '<div class="rk-cap">' + part.cap + '</div>' : '';
      var mode = this.rackMode();
      var clean = this.rackCleanOn();
      var floats = is3 && mode === 'list' && !clean;
      var html;
      if (is3) {
        html = '<div class="rk-wrap is-3d">' +
          '<div class="rk-stage" id="rkStage">' +
            '<div class="rk-tools">' +
              '<button type="button" class="rk-tool" data-act="fs">全屏</button>' +
              '<button type="button" class="rk-tool' + (floats ? ' is-on' : '') + '" data-act="info">信息</button>' +
              '<button type="button" class="rk-tool' + (floats ? ' is-on' : '') + '" data-act="dev">设备</button>' +
              '<button type="button" class="rk-tool" data-act="room">' + (clean ? '显示机房信息' : '隐藏机房信息') + '</button>' +
              '<button type="button" class="rk-tool" data-act="back" style="display:none">返回走道</button>' +
              '<span class="rk-modes">' +
                '<b data-view="patrol"' + (mode === 'patrol' ? ' class="is-on"' : '') + '>巡视</b>' +
                '<b data-view="list"' + (mode === 'list' ? ' class="is-on"' : '') + '>清单</b>' +
              '</span>' +
              '<span class="rk-tip">单击打开面板 · 双击正对 · WASD 走动 · 拖动转头</span>' +
            '</div>' +
            '<div class="rk-mid">' +
              '<div class="iso-stage is-rack" id="v3Host"></div>' +
              '<aside class="rk-float rk-float-l' + (floats ? '' : ' is-off') + '" data-float="dev">' +
                lhead + '<div class="rk-rack">' + part.rows + '</div></aside>' +
              '<aside class="rk-float rk-float-r' + (floats ? '' : ' is-off') + '" data-float="info">' +
                cap + '<div class="rk-right">' + part.right + '</div></aside>' +
            '</div>' +
          '</div></div>';
      } else {
        html = '<div class="rk-wrap">' +
          '<div class="rk-left is-wide">' + lhead +
            '<div class="rk-rack">' + part.rows + '</div></div>' +
          '<div class="rk-right">' + part.right + '</div>' +
        '</div>';
      }

      var nx = this.fly(html, this.focusOf(fromEl), fromEl ? 1 : 1);
      if (is3) {
        this.cleanChrome(clean, true);
        A.V3.mount(A.q('#v3Host', nx), 'rack', cabName, {
          drill: null,
          toast: A.toast,
          onCab: function (name) { self.syncCabPanel(name); },
          onDev: function (sn) { self.markDev(sn); },
          onInside: function (on) {
            var b = nx.querySelector('[data-act="back"]');
            if (b) b.style.display = on ? '' : 'none';
          }
        });
        var hint = document.getElementById('tourHint');
        if (hint) hint.innerHTML = '单击打开这一台面板 · <b>双击</b> 走到正对 · WASD 在走道里走 · <b>Esc</b> 收起面板';
        var tools = nx.querySelector('.rk-tools');
        if (tools) tools.addEventListener('click', function (ev) {
          var t = ev.target.closest ? ev.target.closest('[data-act],[data-view]') : null;
          if (!t || !tools.contains(t)) return;
          ev.stopPropagation();
          var act = t.getAttribute('data-act');
          var view = t.getAttribute('data-view');
          if (act === 'fs') {
            var stage = nx.querySelector('#rkStage');
            if (document.fullscreenElement) { document.exitFullscreen(); return; }
            var req = stage && stage.requestFullscreen ? stage.requestFullscreen() : null;
            if (req && req.catch) req.catch(function () { self.cleanChrome(true); });
            else self.cleanChrome(true);
            return;
          }
          if (act === 'room') { self.cleanChrome(!self.rackCleanOn()); return; }
          if (act === 'back') {
            var vv = A.V3 && A.V3.active;
            if (vv && vv.exitInside) vv.exitInside();
            return;
          }
          if (act === 'info' || act === 'dev') {
            var el = nx.querySelector('[data-float="' + act + '"]');
            if (!el) return;
            el.classList.toggle('is-off');
            t.classList.toggle('is-on', !el.classList.contains('is-off'));
            return;
          }
          if (view === 'patrol' || view === 'list') {
            try { sessionStorage.setItem('bjga.rack.mode', view); } catch (e2) { /* 仅当前页 */ }
            nx.querySelectorAll('.rk-modes b').forEach(function (x) {
              x.classList.toggle('is-on', x.getAttribute('data-view') === view);
            });
            var open = view === 'list';
            self.setFloats(nx, open, open);
            if (open) self.cleanChrome(false);
          }
        });
      }
      nx.addEventListener('click', function (e) {
        var el = e.target.closest ? e.target.closest('.rk-u[data-sn]') : null;
        if (!el || !nx.contains(el)) return;
        var sn = el.getAttribute('data-sn');
        var a = D().assets.filter(function (x) { return x.sn === sn; })[0];
        if (is3) {
          var v = A.V3.active;
          if (v && v.kind === 'rack') v.flyTo('dev:' + sn);
          self.markDev(sn);
        }
        if (a) A.toast(a.cat + ' · ' + a.brand + ' ' + a.model + ' · ' + a.uH + 'U · U' +
          a.uStart + '-' + a.uEnd + ' · ' + a.kw + 'kW · ' + (a.psu || 1) + ' 电源');
        else A.toast('设备 ' + sn);
      });
      A.qa('.iso-dev', nx).forEach(function (el) {
        el.addEventListener('click', function () {
          var a = D().assets.filter(function (x) { return x.sn === el.getAttribute('data-sn'); })[0];
          if (a) A.toast(a.cat + ' · ' + a.brand + ' ' + a.model + ' · ' + a.uH + 'U · U' +
            a.uStart + '-' + a.uEnd + ' · ' + a.kw + 'kW · ' + (a.psu || 1) + ' 电源');
        });
      });
      this.bindMode(nx);
      this.crumbs();
    }
  });
})(window.App);
