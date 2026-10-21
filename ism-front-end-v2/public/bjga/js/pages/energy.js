/* 能效与碳排 */
(function (A) {
  'use strict';
  var SUBS = [
    { id: 'pue', name: 'PUE' },
    { id: 'air', name: '气流与热点' },
    { id: 'cool', name: '冷源能效' },
    { id: 'carbon', name: '碳排' },
    { id: 'diag', name: '诊断' }
  ];
  var S = { sub: 'pue' };

  function head() {
    return '<div class="page-head"><h1 class="page-title">能效碳排</h1>' +
      '<span class="page-crumb">' + A.crumb('人工数据带修订留痕') + '</span></div>' +
      '<div class="subtabs">' + SUBS.map(function (s) {
        return '<button class="subtab' + (s.id === S.sub ? ' is-active' : '') + '" data-sub="' + s.id + '">' + s.name + '</button>';
      }).join('') + '</div><div id="enBody" class="grow" style="overflow:auto"></div>';
  }

  A.register('energy', {
    render: function (root) {
      if (A.subArg && SUBS.some(function (s) { return s.id === A.subArg; })) S.sub = A.subArg;
      root.innerHTML = head();
      this.paint(root);
      var self = this;
      A.qa('.subtab', root).forEach(function (b) {
        b.addEventListener('click', function () { S.sub = b.getAttribute('data-sub'); self.render(root); });
      });
    },
    paint: function (root) {
      var box = A.q('#enBody', root);
      var self = this;
      box.innerHTML = '<div class="empty">加载中</div>';
      var path = S.sub === 'pue' || S.sub === 'diag' ? '/bjga/energy/pue'
        : S.sub === 'air' ? '/bjga/energy/airflow'
        : S.sub === 'cool' ? '/bjga/energy/cooling' : '/bjga/energy/carbon';
      A.api(path).then(function (j) {
        if (S.sub === 'pue' || S.sub === 'diag') self.pue(box, j, S.sub === 'diag');
        else if (S.sub === 'air') self.air(box, j);
        else if (S.sub === 'cool') self.cool(box, j);
        else self.carbon(box, j);
      });
    },
    pue: function (box, j, diag) {
      var live = j.live || {};
      var meters = j.meters || [];
      var revs = j.revisions || [];
      box.innerHTML = '<section class="card"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">PUE ' + A.fix(live.pue, 2) + '</h2>' +
        '<div class="card-extra">本次计算含人工数据 ' + A.fix(live.manualRatio, 1) + '%</div></div><div class="card-body">' +
        '<table class="tbl"><thead><tr><th>计量点</th><th>分类</th><th>来源</th><th>数值</th><th>修订</th></tr></thead><tbody>' +
        meters.map(function (m) {
          var src = m.source === 'manual' ? '人工' : m.source === 'estimate' ? '估算' : '自动';
          return '<tr><td>' + A.esc(m.name) + '</td><td>' + A.esc(m.category) + '</td><td>' + src + '</td>' +
            '<td class="mono">' + m.value + ' ' + A.esc(m.unit) + '</td><td><button class="btn is-ghost is-sm" data-m="' + m.id + '" data-v="' + m.value + '">修改</button></td></tr>';
        }).join('') + '</tbody></table>' +
        (diag ? '<div style="margin-top:12px;color:#FFD24D">高碳时段与高负载机房见碳排页。提高供水温度、延长自然冷却可拉低 PUE。</div>' : '') +
        '<div class="muted" style="margin-top:12px">修订记录</div><table class="tbl"><tbody>' +
        revs.map(function (r) {
          return '<tr><td>' + A.esc(r.at) + '</td><td>' + A.esc(r.actor) + '</td><td class="mono">' + r.before + ' → ' + r.after + '</td><td>' + A.esc(r.reason) + '</td></tr>';
        }).join('') + '</tbody></table></div></section>';
      var self = this;
      A.qa('[data-m]', box).forEach(function (b) {
        b.addEventListener('click', function () {
          var v = window.prompt('新数值', b.getAttribute('data-v'));
          if (v == null) return;
          var reason = window.prompt('修订理由', '现场补录') || '现场补录';
          A.api('/bjga/energy/meters', { method: 'POST', body: { id: Number(b.getAttribute('data-m')), value: Number(v), reason: reason } }).then(function (res) {
            if (res.code !== 1000) { A.toast(res.message || '失败'); return; }
            if (window.DCIM) window.DCIM.live = (res.pue || {}).live || window.DCIM.live;
            A.toast('已修订，PUE ' + A.fix(((res.pue || {}).live || {}).pue, 2));
            self.pue(box, res.pue || res, diag);
          });
        });
      });
    },
    air: function (box, j) {
      var rows = j.all || [];
      box.innerHTML = '<section class="card"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">进风超过 ' + (j.limitC || 27) + '℃ 的机柜</h2></div><div class="card-body"><table class="tbl"><thead><tr>' +
        '<th>机房</th><th>机柜</th><th>进风 ℃</th><th>负载 kW</th><th>建议</th></tr></thead><tbody>' +
        rows.map(function (r) {
          var hot = r.tin > (j.limitC || 27);
          return '<tr><td>' + A.esc(r.room) + '</td><td>' + A.esc(r.cab) + '</td><td style="color:' + (hot ? '#FF4D5E' : '#00E39A') + '">' + r.tin + '</td><td>' + r.loadKw + '</td><td>' + A.esc(r.advice) + '</td></tr>';
        }).join('') + '</tbody></table></div></section>';
    },
    cool: function (box, j) {
      var u = j.unit || {};
      box.innerHTML = '<section class="card"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">' + A.esc(u.name || '冷源') + '</h2></div><div class="card-body">' +
        '<div>回风 ' + u.returnC + '℃ → 送风 ' + u.supplyC + '℃ · 风机功率 ' + u.powerKw + ' kW · 风量 ' + u.airflowM3h + ' m³/h</div>' +
        '<div style="margin:8px 0">估算冷量 ' + (j.qKw || '-') + ' kW · COP ' + (j.cop || '-') + ' · 设计 COP ' + (u.designCop || '-') + '</div>' +
        '<div>' + (j.process || []).map(function (p) { return A.esc(p.name) + ' ' + p.t + '℃'; }).join(' → ') + '</div>' +
        '<div class="muted" style="margin-top:8px">' + A.esc(j.note || '') + '</div></div></section>';
    },
    carbon: function (box, j) {
      box.innerHTML = '<section class="card"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">年碳排约 ' + (j.totalT || '-') + ' t</h2></div><div class="card-body">' +
        '<div>年电量约 ' + (j.annualKwh || '-') + ' kWh · 范围一 ' + (j.scope1kg || 0) + ' kg · 范围二 ' + (j.scope2kg || 0) + ' kg</div>' +
        '<table class="tbl" style="margin-top:10px"><thead><tr><th>来源</th><th>范围</th><th>因子</th><th>依据</th><th>排放 kg</th></tr></thead><tbody>' +
        (j.items || []).map(function (it) {
          return '<tr><td>' + A.esc(it.name) + '</td><td>' + A.esc(it.scope) + '</td><td class="mono">' + it.factor + ' ' + A.esc(it.unit) + '</td><td>' + A.esc(it.source) + '</td><td class="mono">' + it.kg + '</td></tr>';
        }).join('') + '</tbody></table><div style="margin-top:10px">' + A.esc(j.advice || '') + '</div></div></section>';
    }
  });
})(window.App);
