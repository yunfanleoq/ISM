/* 同源 iframe 内读登录态：cookie + 本帧/父页 sessionStorage（父页 axios 常只写 sessionStorage） */
(function (w) {
  'use strict';

  var KEYS = ['Authorization', 'ProjectUuid'];
  /** 航信新大屏对照项目（docs/ISM-麒麟V10部署方案.md；仅在本库存在时启用） */
  var HX_DCIM_PROJECT_UUID = '3ec5821f-b512-2adb-3e1c-473720d0a93e';
  /** 本机预览库常见默认 */
  var XUNZHI_PROJECT_UUID = 'c105271e-9642-c7d2-5c4c-1131a276a81d';

  var parentProjectSnapshot = null;
  var dcimProjectCache = '';
  var projectListPromise = null;

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

  function captureParentProjectOnce() {
    if (parentProjectSnapshot !== null) return parentProjectSnapshot;
    parentProjectSnapshot = get('ProjectUuid') || '';
    return parentProjectSnapshot;
  }

  function writeProjectUuid(pu, opts) {
    if (!pu) return '';
    opts = opts || {};
    var iframeOnly = opts.persist !== true;

    var targets = [w.sessionStorage, w.localStorage];
    if (!iframeOnly) {
      try {
        if (w.top && w.top !== w) {
          targets.push(w.top.sessionStorage, w.top.localStorage);
        }
      } catch (e) {}
    }

    targets.forEach(function (store) {
      if (!store) return;
      try { store.setItem('ProjectUuid', pu); } catch (e2) {}
    });

    if (!iframeOnly) {
      var cookie = 'ProjectUuid=' + encodeURIComponent(pu) + '; path=/; SameSite=Lax';
      try { document.cookie = cookie; } catch (e3) {}
      try {
        if (w.top && w.top !== w) w.top.document.cookie = cookie;
      } catch (e4) {}
    }
    return pu;
  }

  function fetchProjectUuids() {
    if (projectListPromise) return projectListPromise;
    projectListPromise = fetch('/api/ProjectList', {
      method: 'POST',
      headers: (function () {
        var h = { 'Content-Type': 'application/json' };
        var token = get('Authorization');
        if (token) h.Authorization = token;
        return h;
      })(),
      body: '{}'
    }).then(function (r) {
      if (!r.ok) throw new Error('HTTP ' + r.status);
      return r.json();
    }).then(function (body) {
      var list = (body && body.list) || [];
      return list.map(function (row) {
        return (row.ProjectInfo && row.ProjectInfo.uuid) || '';
      }).filter(Boolean);
    }).catch(function () {
      return [];
    });
    return projectListPromise;
  }

  function pickDcimProjectUuid(uuids) {
    captureParentProjectOnce();
    var parent = parentProjectSnapshot || get('ProjectUuid') || '';
    var hasHx = uuids.indexOf(HX_DCIM_PROJECT_UUID) >= 0;

    if (hasHx) {
      return HX_DCIM_PROJECT_UUID;
    }
    if (parent && uuids.indexOf(parent) >= 0) {
      return parent;
    }
    if (parent) return parent;
    if (uuids.length) return uuids[0];
    return XUNZHI_PROJECT_UUID;
  }

  function refreshDcimProjectAsync() {
    return fetchProjectUuids().then(function (uuids) {
      dcimProjectCache = pickDcimProjectUuid(uuids);
      writeProjectUuid(dcimProjectCache, { iframeOnly: true });
      return dcimProjectCache;
    });
  }

  /** 新大屏 API 用项目：仅写 iframe 存储；本库无航信时保留登录时的寻智等项目，不污染父页。 */
  function ensureHangxinProject() {
    captureParentProjectOnce();
    if (dcimProjectCache) {
      writeProjectUuid(dcimProjectCache, { iframeOnly: true });
      return dcimProjectCache;
    }
    var fallback = parentProjectSnapshot || get('ProjectUuid') || XUNZHI_PROJECT_UUID;
    dcimProjectCache = fallback;
    writeProjectUuid(fallback, { iframeOnly: true });
    refreshDcimProjectAsync();
    return fallback;
  }

  /** 齿轮回主系统：恢复进入新大屏前的父页项目，避免航信 UUID 污染后台应用列表。 */
  function persistProject() {
    captureParentProjectOnce();
    var pu = parentProjectSnapshot || dcimProjectCache || get('ProjectUuid');
    if (!pu) return '';
    return writeProjectUuid(pu, { persist: true });
  }

  function headers() {
    ensureHangxinProject();
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
    ensureHangxinProject: ensureHangxinProject,
    persistProject: persistProject,
    authHint: function () {
      if (hasLogin()) return '';
      return '未检测到登录令牌（Authorization）。请先在主系统登录，再通过菜单「新大屏」或 /#/hx-dcim 进入；勿在未登录时直接打开 /hx-dcim/index.html。';
    }
  };
})(window);
