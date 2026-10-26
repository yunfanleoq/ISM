/* 运行监控：漏水 / UPS 直接展示采集实时值。温湿度与列头柜等在 boot-data 里覆盖序列。 */
(function (A) {
  'use strict';

  function devs(kind) {
    var list = (window.DCIM && window.DCIM.live && window.DCIM.live.devices) || [];
    return list.filter(function (d) { return d.kind === kind; });
  }

  function warnPoint(p) {
    if (!/信号|故障|高报|告警/.test(p.name || '')) return false;
    return Number(p.value) !== 0;
  }

  A.liveView = function (kind, title) {
    var list = devs(kind);
    if (!list.length) return null;
    var rows = [];
    list.forEach(function (d) {
      (d.points || []).forEach(function (p) {
        rows.push({
          dev: d.name,
          point: p.name,
          value: p.value,
          unit: p.unit || '',
          st: warnPoint(p) ? 'warn' : 'ok'
        });
      });
    });
    var warnN = rows.filter(function (r) { return r.st === 'warn'; }).length;
    return {
      stats: A.UI.stats([
        { v: list.length, u: '台', l: '采集设备' },
        { v: rows.length, u: '个', l: '实时测点' },
        { v: warnN, u: '个', l: '告警测点', color: warnN ? 'var(--lv1)' : 'var(--green)' },
        { v: warnN ? '告警' : '正常', l: '总体状态', color: warnN ? 'var(--lv1)' : 'var(--green)' }
      ]),
      body: '<div class="card grow">' +
        '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">' + A.esc(title) + ' · 采集实时值</h2>' +
          '<div class="card-extra"><span class="muted" style="font-size:12px">当前采集</span></div></div>' +
        '<div class="card-body">' + A.UI.table({
          cols: [
            { t: '设备', k: 'dev', w: 260 },
            { t: '测点', k: 'point', w: 180 },
            { t: '实时值', k: 'value', align: 'r', w: 120, html: function (r) {
              return '<b class="live-val" style="font-family:var(--num);color:var(--cyan)">' + A.esc(String(r.value)) + '</b>';
            } },
            { t: '单位', k: 'unit', w: 80 },
            { t: '状态', k: 'st', w: 100, html: function (r) {
              return r.st === 'warn' ? A.UI.tag('告警', 'red') : A.UI.tag('正常', 'green');
            } }
          ],
          rows: rows
        }) + '</div></div>',
      items: []
    };
  };

  function pairs(m) {
    return Object.keys(m).map(function (k) { return { name: k, value: m[k] }; });
  }

  A.alarmFeed = function () {
    var D = window.DCIM;
    var live = D && D.live && D.live.alarms;
    var list = (Array.isArray(live) && live.length) ? live : ((D && D.alarms) || []);
    if (!list.length) return null;
    var levels = {}, cats = {}, rec = 0;
    list.forEach(function (a) {
      if (a.level) levels[a.level] = (levels[a.level] || 0) + 1;
      if (a.cat) cats[a.cat] = (cats[a.cat] || 0) + 1;
      if (a.recover === '已恢复') rec++;
    });
    var byLevel = pairs(levels);
    var byCat = pairs(cats);
    return {
      list: list,
      stats: {
        total: list.length,
        recovered: rec,
        pending: list.length - rec,
        byLevel: byLevel.length ? byLevel : [{ name: '无', value: 0 }],
        byCat: byCat.length ? byCat : [{ name: '采集', value: 0 }]
      }
    };
  };
})(window.App);
