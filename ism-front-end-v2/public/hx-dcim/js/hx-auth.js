/* 同源 iframe 内读登录态：cookie + 本帧/父页 sessionStorage（父页 axios 常只写 sessionStorage） */
(function (w) {
  'use strict';

  var KEYS = ['Authorization', 'ProjectUuid'];
  /** 航信新大屏对照项目（docs/ISM-麒麟V10部署方案.md） */
  var HX_DCIM_PROJECT_UUID = '3ec5821f-b512-2adb-3e1c-473720d0a93e';
  /** 本机预览库常见默认，勿与新大屏混用 */
  var XUNZHI_PROJECT_UUID = 'c105271e-9642-c7d2-5c4c-1131a276a81d';

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

  function writeProjectUuid(pu) {
    if (!pu) return '';
    var targets = [w.sessionStorage, w.localStorage];
    try {
      if (w.top && w.top !== w) {
        targets.push(w.top.sessionStorage, w.top.localStorage);
      }
    } catch (e) {}
    targets.forEach(function (store) {
      if (!store) return;
      try { store.setItem('ProjectUuid', pu); } catch (e2) {}
    });
    var cookie = 'ProjectUuid=' + encodeURIComponent(pu) + '; path=/; SameSite=Lax';
    try { document.cookie = cookie; } catch (e3) {}
    try {
      if (w.top && w.top !== w) w.top.document.cookie = cookie;
    } catch (e4) {}
    return pu;
  }

  /** 新大屏须绑航信项目；登录后若仍带寻智/空 UUID 则纠正（iframe 与父页同步）。 */
  function ensureHangxinProject() {
    var pu = get('ProjectUuid');
    if (pu && pu !== XUNZHI_PROJECT_UUID) return pu;
    return writeProjectUuid(HX_DCIM_PROJECT_UUID);
  }

  /** 齿轮回主系统前：把当前大屏用的 ProjectUuid 写回父页 cookie/session，避免仍落在寻智等项目。 */
  function persistProject() {
    ensureHangxinProject();
    var pu = get('ProjectUuid');
    if (!pu) return '';
    return writeProjectUuid(pu);
  }

  ensureHangxinProject();

  w.HXAuth = {
    get: get,
    headers: headers,
    hasLogin: hasLogin,
    ensureHangxinProject: ensureHangxinProject,
    persistProject: persistProject,
    authHint: function () {
      if (hasLogin()) return '';
      return '未检测到登录令牌（Authorization）。请先在主系统登录，再通过菜单「新大屏」或 /#/hx-dcim 进入；勿在未登录时直接打开 /hx-dcim/index.html。';
    }
  };
})(window);
