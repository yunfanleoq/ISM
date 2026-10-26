/* 安全审计 */
(function (A) {
  'use strict';
  var SUBS = [
    { id: 'log', name: '审计日志' },
    { id: 'policy', name: '安全策略' },
    { id: 'eq', name: '等保对照' },
    { id: 'lic', name: '授权台账' }
  ];
  var S = { sub: 'log' };

  A.register('audit', {
    render: function (root) {
      if (A.subArg && SUBS.some(function (s) { return s.id === A.subArg; })) S.sub = A.subArg;
      root.innerHTML = '<div class="page-head"><h1 class="page-title">安全审计</h1>' +
        '<span class="page-crumb">' + A.crumb('人 · 时间 · 对象 · 前后值') + '</span></div>' +
        '<div class="subtabs">' + SUBS.map(function (s) {
          return '<button class="subtab' + (s.id === S.sub ? ' is-active' : '') + '" data-sub="' + s.id + '">' + s.name + '</button>';
        }).join('') + '</div><div id="auBody" class="grow" style="overflow:auto"></div>';
      this.paint(root);
      var self = this;
      A.qa('.subtab', root).forEach(function (b) {
        b.addEventListener('click', function () { S.sub = b.getAttribute('data-sub'); self.render(root); });
      });
    },
    paint: function (root) {
      var box = A.q('#auBody', root);
      if (S.sub === 'log') {
        A.api('/bjga/audit').then(function (j) {
          var rows = j.logs || [];
          box.innerHTML = '<section class="card"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">操作留痕</h2>' +
            '<div class="card-extra">记录追加保存，页面不提供删除</div></div><div class="card-body"><table class="tbl"><thead><tr>' +
            '<th>时间</th><th>人</th><th>动作</th><th>对象</th><th>前</th><th>后</th></tr></thead><tbody>' +
            rows.map(function (r) {
              return '<tr><td class="mono">' + A.esc(r.at) + '</td><td>' + A.esc(r.actor) + '</td><td>' + A.esc(r.action) + '</td><td>' + A.esc(r.object) + '</td><td>' + A.esc(r.before) + '</td><td>' + A.esc(r.after) + '</td></tr>';
            }).join('') + '</tbody></table></div></section>';
        });
        return;
      }
      A.api('/bjga/audit/policy').then(function (j) {
        var p = j.policy || {};
        if (S.sub === 'policy') {
          box.innerHTML = '<section class="card"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">口令与锁定</h2></div><div class="card-body">' +
            field('最短长度', 'minLen', p.minLen) + field('字符类别至少', 'minClasses', p.minClasses) +
            field('更换周期（天）', 'maxAgeDays', p.maxAgeDays) + field('历史不重复次数', 'historyN', p.historyN) +
            field('失败锁定次数', 'failLimit', p.failLimit) + field('锁定分钟', 'lockMin', p.lockMin) +
            '<button class="btn is-primary is-sm" id="auSave">保存策略</button></div></section>';
          A.q('#auSave', box).addEventListener('click', function () {
            var body = {};
            ['minLen', 'minClasses', 'maxAgeDays', 'historyN', 'failLimit', 'lockMin'].forEach(function (k) {
              body[k] = Number(A.q('[data-k="' + k + '"]', box).value);
            });
            A.api('/bjga/audit/policy', { method: 'POST', body: body }).then(function (res) {
              A.toast(res.code === 1000 ? '策略已保存并写入审计' : '失败');
            });
          });
        } else if (S.sub === 'eq') {
          box.innerHTML = tbl('等保 2.0 对照（条款 → 本系统位置）', ['条款', '要求', '实现位置', '状态'], (j.compliance || []).map(function (r) {
            return [r.clause, r.item, r.where, r.status];
          }));
        } else {
          box.innerHTML = tbl('软件资产与授权', ['名称', '类型', '版本', '数量', '证书', '到期', '位置'], (j.licenses || []).map(function (r) {
            return [r.name, r.kind, r.version, r.count, r.certNo, r.expire, r.place];
          }));
        }
      });
    }
  });

  function field(label, key, val) {
    return '<div style="margin:8px 0">' + label + ' <input data-k="' + key + '" value="' + A.esc(val) + '" style="width:80px;background:#071228;border:1px solid rgba(28,90,158,.8);color:#E6F3FF"/></div>';
  }
  function tbl(title, heads, rows) {
    return '<section class="card"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">' + A.esc(title) + '</h2></div><div class="card-body"><table class="tbl"><thead><tr>' +
      heads.map(function (h) { return '<th>' + A.esc(h) + '</th>'; }).join('') + '</tr></thead><tbody>' +
      rows.map(function (r) { return '<tr>' + r.map(function (c) { return '<td>' + A.esc(c) + '</td>'; }).join('') + '</tr>'; }).join('') +
      '</tbody></table></div></section>';
  }
})(window.App);
