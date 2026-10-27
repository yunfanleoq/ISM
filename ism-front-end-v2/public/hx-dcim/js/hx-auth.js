/* 同源 iframe 内读登录态：cookie + 本帧/父页 sessionStorage（父页 axios 常只写 sessionStorage） */
(function (w) {
  'use strict';

  var KEYS = ['Authorization', 'ProjectUuid'];

  function fromCookie(key) {
    var m = document.cookie.match(new RegExp('(?:^|; )' + key + '=([^;]*)'));
    if (!m) return '';
    try { return decodeURIComponent(m[1]); } catch (e) { return m[1]; }
  }

  function fromStore(store, key) {
    if (!store) return '';
    try {
      var s = store.getItem(key);
      if (s) return s;
    } catch (e) {}
    return '';
  }

  function get(key) {
    var v = fromStore(w.sessionStorage, key);
    if (v) return v;
    v = fromCookie(key);
    if (v) return v;
    try {
      if (w.parent && w.parent !== w) {
        v = fromStore(w.parent.sessionStorage, key);
        if (v) return v;
      }
    } catch (e2) {}
    try {
      v = fromStore(w.localStorage, key);
      if (v) return v;
    } catch (e3) {}
    try {
      if (w.parent && w.parent !== w) {
        v = fromStore(w.parent.localStorage, key);
        if (v) return v;
      }
    } catch (e4) {}
    return '';
  }

  function headers() {
    var h = { 'Content-Type': 'application/json' };
    KEYS.forEach(function (k) {
      var v = get(k);
      if (v) h[k] = v;
    });
    return h;
  }

  function hasLogin() {
    return !!get('Authorization');
  }

  w.HXAuth = {
    get: get,
    headers: headers,
    hasLogin: hasLogin,
    authHint: function () {
      if (hasLogin()) return '';
      return '未检测到登录令牌（Authorization）。请先在主系统登录，再通过菜单「新大屏」或 /#/hx-dcim 进入；勿在未登录时直接打开 /hx-dcim/index.html。';
    }
  };
})(window);
