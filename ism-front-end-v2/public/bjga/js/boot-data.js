(function (A) {
  'use strict';

  function token() {
    try {
      var s = sessionStorage.getItem('Authorization');
      if (s) return s;
    } catch (e) {}
    var m = document.cookie.match(/(?:^|; )Authorization=([^;]*)/);
    if (!m) return '';
    try { return decodeURIComponent(m[1]); } catch (e2) { return m[1]; }
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
    }).then(function (r) { return r.json(); });
  };

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
    if (by.th && by.th.length) dcim.thDevices = by.th.map(toTh);
    var keys = { pwr: '列头柜', ac: '空调', h2: '氢气', batt: '蓄电池' };
    Object.keys(keys).forEach(function (k) {
      if (by[k] && by[k].length) dcim.series[keys[k]] = toSeries(by[k]);
    });
  }

  var pr = document.getElementById('pageRoot');
  if (pr) {
    pr.innerHTML = '<div class="empty">正在加载动环数据…</div>';
  }
  A.api('/bjga/bootstrap').then(function (j) {
    if (!j || j.code !== 1000 || !j.dcim) {
      throw new Error((j && (j.message || j.msg)) || 'bootstrap 失败');
    }
    window.DCIM = j.dcim;
    applyCollect(window.DCIM);
    if (A._ready) A._ready();
  }).catch(function (err) {
    if (!pr) return;
    pr.innerHTML = '<div class="empty">动环数据未加载<br>' + A.esc(err.message || err) +
      '<br><span style="color:#7FA8D9">请先登录系统，再打开本页。</span></div>';
  });
})(window.App);
