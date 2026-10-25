/* 数据治理与接口目录 */
(function (A) {
  'use strict';
  A.register('govern', {
    render: function (root) {
      root.innerHTML = '<div class="page-head"><h1 class="page-title">数据治理</h1>' +
        '<span class="page-crumb">' + A.crumb('校验规则 · 南向北向') + '</span></div><div id="gvBody" class="grow" style="overflow:auto"></div>';
      A.api('/bjga/govern').then(function (j) {
        var box = A.q('#gvBody', root);
        box.innerHTML =
          '<section class="card" style="margin-bottom:12px"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">校验规则</h2></div><div class="card-body"><table class="tbl"><thead><tr><th>规则</th><th>类型</th><th>范围</th><th>状态</th><th>命中</th></tr></thead><tbody>' +
          (j.rules || []).map(function (r) {
            return '<tr><td>' + A.esc(r.name) + '</td><td>' + A.esc(r.kind) + '</td><td>' + A.esc(r.scope) + '</td><td>' + A.esc(r.status) + '</td><td class="mono">' + r.hits + '</td></tr>';
          }).join('') + '</tbody></table></div></section>' +
          '<section class="card"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">接口目录</h2></div><div class="card-body"><table class="tbl"><thead><tr><th>名称</th><th>方向</th><th>协议</th><th>状态</th><th>示例</th></tr></thead><tbody>' +
          (j.interfaces || []).map(function (r) {
            return '<tr><td>' + A.esc(r.name) + '</td><td>' + A.esc(r.direct) + '</td><td>' + A.esc(r.proto) + '</td><td>' + A.esc(r.status) + '</td><td>' + A.esc(r.sample) + '</td></tr>';
          }).join('') + '</tbody></table><div class="muted" style="margin-top:8px">北向示例即本系统 GET /bjga/bootstrap。门禁码流在厂商 SDK 到位前保持待开放。</div></div></section>';
      });
    }
  });
})(window.App);
