(function (A) {
  'use strict';

  // 与主站 request.js 一致：Authorization 存在 sessionStorage / Cookie
  var DEMO_MD5 = 'e10adc3949ba59abbe56e057f20f883e'; // MD5('123456')，仅无令牌时兜底

  function token() {
    try {
      var s = sessionStorage.getItem('Authorization');
      if (s) return s;
    } catch (e) {}
    var m = document.cookie.match(/(?:^|; )Authorization=([^;]*)/);
    if (!m) return '';
    try { return decodeURIComponent(m[1]); } catch (e2) { return m[1]; }
  }

  function saveToken(tok, expireAt) {
    if (!tok) return;
    try { sessionStorage.setItem('Authorization', tok); } catch (e) {}
    var maxAge = '';
    var n = Number(expireAt);
    if (isFinite(n) && n > 1e6) {
      var sec = Math.floor(n - Date.now() / 1000);
      if (sec > 0) maxAge = '; max-age=' + sec;
    }
    document.cookie = 'Authorization=' + encodeURIComponent(tok) + '; path=/' + maxAge;
  }

  function parseJsonResponse(r) {
    return r.text().then(function (t) {
      var s = (t || '').replace(/^\uFEFF/, '').trim();
      if (!s || s.charAt(0) === '<') {
        throw new Error('接口未返回 JSON（HTTP ' + r.status + '）。通常是未登录被重定向到登录页。');
      }
      try {
        return JSON.parse(s);
      } catch (e) {
        throw new Error('接口 JSON 解析失败（HTTP ' + r.status + '）');
      }
    });
  }

  A.api = function (path, opt) {
    opt = opt || {};
    var headers = { 'Content-Type': 'application/json' };
    var t = token();
    if (t) headers.Authorization = t;
    return fetch('/api' + path, {
      method: opt.method || 'GET',
      headers: headers,
      body: opt.body ? JSON.stringify(opt.body) : undefined
    }).then(parseJsonResponse);
  };

  /** 无令牌时：本机演示用 admin/123456（MD5）自动登录，避免直开大屏白屏 */
  function ensureAuth() {
    if (token()) return Promise.resolve();
    return fetch('/api/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ Username: 'admin', password: DEMO_MD5 })
    }).then(parseJsonResponse).then(function (j) {
      if (!j || j.code !== 1000 || !j.data || !j.data.token) {
        throw new Error((j && (j.message || j.msg)) || '登录失败');
      }
      saveToken(j.data.token, j.data.expireAt);
    });
  }

  function num(v) {
    var n = Number(v);
    return isFinite(n) ? n : 0;
  }

  function flat(v, n) {
    var a = [];
    for (var i = 0; i < n; i++) a.push(v);
    return a;
  }

  function pointOf(d, name) {
    var pts = d.points || [];
    for (var i = 0; i < pts.length; i++) {
      if (pts[i].name === name) return pts[i];
    }
    return null;
  }

  function toTh(d) {
    var t = pointOf(d, '温度');
    var h = pointOf(d, '湿度');
    var tv = t ? num(t.value) : 0;
    var hv = h ? num(h.value) : 0;
    return {
      name: d.name, room: '采集', ch: '实时',
      temp: tv, humi: hv,
      tempMin: tv, tempMax: tv, humiMin: hv, humiMax: hv,
      tSeries: flat(tv, 24), hSeries: flat(hv, 24),
      live: true
    };
  }

  function toSeries(devs) {
    var now = new Date();
    var ds = now.getFullYear() + '-' +
      ('0' + (now.getMonth() + 1)).slice(-2) + '-' +
      ('0' + now.getDate()).slice(-2) + ' 00:00';
    return {
      live: true,
      dates: [ds],
      devices: devs.map(function (d) {
        return {
          name: d.name,
          loc: d.model || '采集',
          points: (d.points || []).map(function (p) {
            var v = num(p.value);
            return { name: p.name, last: v, avg: v, max: v, min: v, vals: [v], unit: p.unit || '' };
          })
        };
      })
    };
  }

  function applyCollect(dcim) {
    var devs = (dcim.live && dcim.live.devices) || [];
    if (!devs.length) return;
    if (!dcim.series) dcim.series = {};
    var by = {};
    devs.forEach(function (d) { (by[d.kind] = by[d.kind] || []).push(d); });
    if (by.th && by.th.length) {
      var known = {};
      (dcim.rooms || []).forEach(function (r) { known[r.name] = 1; });
      var named = by.th.map(toTh).filter(function (t) { return known[t.room]; });
      if (named.length) dcim.thDevices = named;
    }
    var keys = { pwr: '列头柜', ac: '空调', h2: '氢气', batt: '蓄电池' };
    Object.keys(keys).forEach(function (k) {
      if (by[k] && by[k].length) dcim.series[keys[k]] = toSeries(by[k]);
    });
  }

  var pr = document.getElementById('pageRoot');
  if (pr) {
    pr.innerHTML = '<div class="empty">正在加载动环数据…</div>';
  }
  ensureAuth().then(function () {
    return A.api('/bjga/bootstrap');
  }).then(function (j) {
    if (!j || j.code !== 1000 || !j.dcim) {
      throw new Error((j && (j.message || j.msg)) || 'bootstrap 失败');
    }
    window.DCIM = j.dcim;
    applyCollect(window.DCIM);
    if (A.startAlarmPoll) A.startAlarmPoll(2000);
    if (A._ready) A._ready();
  }).catch(function (err) {
    if (!pr) return;
    pr.innerHTML = '<div class="empty">动环数据未加载<br>' + A.esc(err.message || err) +
      '<br><span style="color:#7FA8D9">请确认后端 8082 可用，或先登录主站后再打开本页。</span></div>';
  });
})(window.App);
