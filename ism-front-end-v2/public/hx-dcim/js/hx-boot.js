(function (w) {
  'use strict';

  var bar = null;

  function ensureBar() {
    if (bar) return bar;
    var root = document.getElementById('pageRoot');
    if (!root) return null;
    bar = document.createElement('div');
    bar.id = 'hxBootBar';
    bar.className = 'hx-boot-bar';
    bar.innerHTML =
      '<div class="hx-boot-inner">' +
        '<div class="hx-boot-title" id="hxBootTitle">新大屏加载中</div>' +
        '<div class="hx-boot-sub" id="hxBootSub"></div>' +
      '</div>';
    root.appendChild(bar);
    return bar;
  }

  function hide() {
    if (bar) {
      bar.remove();
      bar = null;
    }
  }

  w.HXBoot = {
    status: function (msg) {
      var b = ensureBar();
      if (!b) return;
      var t = document.getElementById('hxBootTitle');
      var s = document.getElementById('hxBootSub');
      if (t) t.textContent = msg || '加载中…';
      if (s && (!s.textContent || s.textContent.indexOf('⚠') !== 0)) s.textContent = '';
    },
    warn: function (msg) {
      var b = ensureBar();
      if (!b) return;
      var s = document.getElementById('hxBootSub');
      if (s) s.textContent = '⚠ ' + msg;
    },
    hide: hide
  };
})(window);
