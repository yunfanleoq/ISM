/* =========================================================
   核心：工具函数 / 缩放舞台 / 路由 / 页面注册
   ========================================================= */
window.App = window.App || {};

(function (A) {
  'use strict';

  /* ---------------- DOM ---------------- */
  A.q  = function (sel, root) { return (root || document).querySelector(sel); };
  A.qa = function (sel, root) { return Array.prototype.slice.call((root || document).querySelectorAll(sel)); };

  A.h = function (html) {
    var t = document.createElement('template');
    t.innerHTML = String(html).trim();
    return t.content.firstElementChild;
  };

  A.esc = function (s) {
    return String(s == null ? '' : s)
      .replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;').replace(/'/g, '&#39;');
  };

  /* ---------------- 数值 ---------------- */
  A.num  = function (v, d) {
    var n = Number(v);
    if (!isFinite(n)) return '-';
    return n.toLocaleString('zh-CN', { minimumFractionDigits: d || 0, maximumFractionDigits: d || 0 });
  };
  A.fix  = function (v, d) { var n = Number(v); return isFinite(n) ? n.toFixed(d == null ? 2 : d) : '-'; };
  A.pct  = function (a, b) { return b ? Math.round(a / b * 1000) / 10 : 0; };
  A.clamp = function (v, a, b) { return Math.min(b, Math.max(a, v)); };

  A.pad = function (n) { return n < 10 ? '0' + n : '' + n; };
  A.now = function () {
    var d = new Date();
    return d.getFullYear() + '-' + A.pad(d.getMonth() + 1) + '-' + A.pad(d.getDate()) + ' ' +
           A.pad(d.getHours()) + ':' + A.pad(d.getMinutes()) + ':' + A.pad(d.getSeconds());
  };

  /* 生成一个围绕基准波动的序列（用于趋势演示） */
  A.wave = function (base, n, amp, seed) {
    var s = seed || 1, out = [], r = function () {
      s = (s * 9301 + 49297) % 233280; return s / 233280;
    };
    for (var i = 0; i < n; i++) out.push(base + (r() - 0.5) * amp * 2 + Math.sin(i / 3) * amp * 0.35);
    return out;
  };

  /* ---------------- 缩放 ----------------
     设计稿为 1920×1080 固定画布。窗口足够大时等比缩放铺满；
     窗口过窄时不再无限缩小（否则整屏只剩顶部一条），
     而是锁定最小比例 0.46 并允许滚动查看。              */
  var DESIGN_W = 1920, DESIGN_H = 1080, MIN_S = 0.46;
  A.scale = 1;

  /* 强制重铺一次栅格。
     整屏是被 transform 缩放的合成层，缩放比一变（窗口拉伸 / 浏览器 Ctrl+± 缩放），
     浏览器会沿用旧比例的栅格结果；面积大的动画图层（如整栋楼 SVG）就会出现黑块。
     换一次图层提示，等于让合成器丢弃缓存重新栅格 —— 和「重开一次页面」同效。 */
  A.repaint = function (el) {
    var n = el || document.getElementById('screen');
    if (!n) return;
    n.style.willChange = 'transform';
    requestAnimationFrame(function () {
      requestAnimationFrame(function () { n.style.willChange = ''; });
    });
  };

  var lastS = 0;
  A.fit = function () {
    var raw = Math.min(window.innerWidth / DESIGN_W, window.innerHeight / DESIGN_H);
    var s = (raw > 0 && isFinite(raw)) ? Math.max(raw, MIN_S) : 1;
    var changed = Math.abs(s - lastS) > 0.004;
    A.scale = s;
    lastS = s;

    var el = document.getElementById('screen');
    var sizer = document.getElementById('sizer');
    if (el) {
      el.style.transformOrigin = 'top left';
      el.style.transform = 'scale(' + s + ')';
    }
    if (sizer) {
      sizer.style.width = Math.round(DESIGN_W * s) + 'px';
      sizer.style.height = Math.round(DESIGN_H * s) + 'px';
    }
    if (changed) A.repaint(el);
  };
  window.addEventListener('resize', A.fit);

  A.fullscreen = function () {
    var el = document.documentElement;
    if (!document.fullscreenElement && !document.webkitFullscreenElement) {
      (el.requestFullscreen || el.webkitRequestFullscreen || function () {}).call(el);
    } else {
      (document.exitFullscreen || document.webkitExitFullscreen || function () {}).call(document);
    }
  };

  /* ---------------- 提示 ---------------- */
  A.toast = function (msg) {
    var box = document.getElementById('toast');
    if (!box) return;
    var t = A.h('<div class="toast-item">' + A.esc(msg) + '</div>');
    box.appendChild(t);
    setTimeout(function () {
      t.style.transition = 'opacity .3s ease'; t.style.opacity = '0';
      setTimeout(function () { t.remove(); }, 320);
    }, 2200);
  };

  /* ---------------- 抽屉 ---------------- */
  A.drawer = function (title, bodyHtml) {
    var root = A.q('.page') || A.q('#pageRoot');
    A.closeDrawer();
    var wrap = A.h('<div class="drawer-mask"></div>');
    var dr = A.h(
      '<aside class="drawer">' +
        '<div class="drawer-head"><span class="drawer-title">' + A.esc(title) + '</span>' +
        '<button class="drawer-close" type="button">✕</button></div>' +
        '<div class="drawer-body">' + bodyHtml + '</div>' +
      '</aside>');
    wrap.appendChild(dr);
    root.appendChild(wrap);
    wrap.addEventListener('click', function (e) { if (e.target === wrap) A.closeDrawer(); });
    A.q('.drawer-close', dr).addEventListener('click', A.closeDrawer);
  };
  A.closeDrawer = function () {
    A.qa('.drawer-mask').forEach(function (n) { n.remove(); });
  };

  /* ---------------- 派生字段（不入数据，减小体积） ---------------- */
  A.derive = function (a) {
    return {
      psu: a.cat === '服务器' ? 2 : 1,
      weight: a.cat === '服务器' ? 30 : (a.cat === '存储设备' ? 20 : 12),
      bizFull: '北京市公安局' + a.biz,
      dept: String(a.org) + ' · 技术保障科',
      ctime: '2026-' + A.pad(1 + (String(a.sn).length * 7 + a.uStart) % 10) + '-' +
             A.pad(1 + (a.uStart * 3 + a.uH) % 28)
    };
  };

  /* ---------------- 面包屑：系统名 · 项目名 / 子页 ---------------- */
  A.brand = function () {
    var m = (window.DCIM && window.DCIM.meta) || {};
    return [m.system, m.project].filter(Boolean).join(' · ');
  };
  A.crumb = function (sub) {
    var head = A.brand();
    return A.esc(head + (sub ? ' / ' + sub : ''));
  };

  /* ---------------- 分页计算 ---------------- */
  A.paginate = function (list, page, size) {
    var total = list.length;
    var pages = Math.max(1, Math.ceil(total / size));
    var p = A.clamp(page || 1, 1, pages);
    return { rows: list.slice((p - 1) * size, p * size), page: p, pages: pages, total: total };
  };

  /* ---------------- 路由 ---------------- */
  A.pages = {};
  A.current = null;

  A.register = function (id, def) { A.pages[id] = def; };

  A.NAV = [
    { id: 'dashboard', name: '驾驶舱' },
    { id: 'building',  name: '空间巡检' },
    { id: 'facility',  name: '运行监控' },
    { id: 'power',     name: '配电专业' },
    { id: 'energy',    name: '能效碳排' },
    { id: 'capacity',  name: '容量总览' },
    { id: 'assetmgr',  name: '资产管理' },
    { id: 'report',    name: '报表报告' },
    { id: 'ops',       name: '运维中心' },
    { id: 'audit',     name: '安全审计' }
  ];
  A.MORE = [
    { id: 'govern',  name: '数据治理' },
    { id: 'config',  name: '工程配置' },
    { id: 'selfmon', name: '自监控' },
    { id: 'admin',   name: '系统管理' },
    { id: 'ontology', name: '本体中心', href: '/#/xunzhi/jf/ontology' },
    { id: 'backend', name: '进入后台', href: '/#/Project' }
  ];

  A.leaveTo = function (href) {
    try {
      if (window.parent && window.parent !== window) window.parent.location.href = href;
      else window.location.href = href;
    } catch (e) {
      window.open(href, '_top');
    }
  };

  A.go = function (id) {
    var top = String(id || '').split('/')[0];
    if (!A.pages[top]) id = 'dashboard';
    location.hash = '#/' + id;
  };

  A.render = function () {
    var raw = (location.hash || '').replace('#/', '');
    var parts = raw.split('/');
    var id = parts[0] || 'dashboard';
    A.subArg = parts[1] || '';
    if (!A.pages[id]) id = 'dashboard';

    A.qa('.tab').forEach(function (t) {
      t.classList.toggle('is-active', t.getAttribute('data-id') === id);
    });

    if (A.current && A.current.unmount) { try { A.current.unmount(); } catch (e) {} }
    A.closeDrawer();

    var root = A.q('#pageRoot');
    root.innerHTML = '';
    var def = A.pages[id];
    var page = A.h('<section class="page"></section>');
    root.appendChild(page);
    A.current = def;
    try {
      def.render(page);
      if (def.mount) def.mount(page);
    } catch (err) {
      page.innerHTML = '<div class="empty">页面渲染异常<br>' + A.esc(err.message) + '</div>';
      console.error(err);
    }
  };

  window.addEventListener('hashchange', A.render);

  A.start = function (fn) {
    if (window.DCIM) { fn(); return; }
    (A._wait = A._wait || []).push(fn);
  };
  A._ready = function () {
    (A._wait || []).forEach(function (fn) { try { fn(); } catch (e) { console.error(e); } });
    A._wait = [];
  };

  A.boot = function () {
    var tabs = A.q('#tabs');
    tabs.style.gap = '12px';
    tabs.innerHTML = A.NAV.map(function (n) {
      return '<button class="tab" type="button" data-id="' + n.id + '">' + n.name + '</button>';
    }).join('') +
      '<div class="tab-more"><button class="tab" type="button">更多</button><div class="tab-more-pop">' +
      A.MORE.map(function (n) {
        return '<button class="tab" type="button" data-id="' + (n.id || '') + '" data-href="' + (n.href || '') + '">' + n.name + '</button>';
      }).join('') + '</div></div>';
    tabs.addEventListener('click', function (e) {
      var b = e.target.closest('.tab');
      if (!b) return;
      var href = b.getAttribute('data-href');
      if (href) { A.leaveTo(href); return; }
      A.go(b.getAttribute('data-id'));
    });
    A.fit();
    A.render();
    /* 顶部时钟 */
    var clock = A.q('#clock');
    if (clock) {
      var tick = function () { clock.textContent = A.now(); };
      tick(); setInterval(tick, 1000);
    }
  };
})(window.App);
