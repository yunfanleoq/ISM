/* 资产管理附加：运维计划、出入库 */
(function (A) {
  'use strict';
  var AM = A.AM = A.AM || {};
  AM.views = AM.views || {};

  AM.views.maint = function (box) {
    var list = ((window.DCIM.extra || {}).maint) || [];
    box.innerHTML = '<section class="card"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">运行维护计划</h2>' +
      '<div class="card-extra">配电柜 / 列头柜 / 精密空调按周期提示</div></div><div class="card-body"><table class="tbl"><thead><tr>' +
      '<th>设备</th><th>类别</th><th>项目</th><th>周期</th><th>下次</th><th>状态</th><th>依据</th></tr></thead><tbody>' +
      list.map(function (r) {
        var c = r.status === '到期' ? '#FF4D5E' : r.status === '即将到期' ? '#FFD24D' : '#00E39A';
        return '<tr><td>' + A.esc(r.device) + '</td><td>' + A.esc(r.cat) + '</td><td>' + A.esc(r.item) + '</td><td>' + A.esc(r.cycle) + '</td><td class="mono">' + A.esc(r.nextDue) + '</td><td style="color:' + c + '">' + A.esc(r.status) + '</td><td>' + A.esc(r.spec) + '</td></tr>';
      }).join('') + '</tbody></table></div></section>';
  };

  AM.views.inout = function (box) {
    var list = ((window.DCIM.extra || {}).bills) || [];
    box.innerHTML = '<section class="card"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">出入库单</h2>' +
      '<div class="card-extra"><select id="ioKind"><option>入库</option><option>出库</option></select> ' +
      '<input id="ioSn" placeholder="SN" style="width:180px;background:#071228;border:1px solid rgba(28,90,158,.8);color:#E6F3FF"/> ' +
      '<button class="btn is-primary is-sm" id="ioAdd">登记</button></div></div><div class="card-body"><table class="tbl"><thead><tr>' +
      '<th>单号</th><th>类型</th><th>SN</th><th>状态</th><th>说明</th><th>时间</th></tr></thead><tbody>' +
      list.map(function (r) {
        return '<tr><td class="mono">' + A.esc(r.no) + '</td><td>' + A.esc(r.kind) + '</td><td class="mono">' + A.esc(r.sn) + '</td><td>' + A.esc(r.status) + '</td><td>' + A.esc(r.note) + '</td><td class="mono">' + A.esc(r.at) + '</td></tr>';
      }).join('') + '</tbody></table></div></section>';
    A.q('#ioAdd', box).addEventListener('click', function () {
      A.api('/bjga/stock-bill', { method: 'POST', body: { kind: A.q('#ioKind', box).value, sn: A.q('#ioSn', box).value, note: '页面登记' } }).then(function (j) {
        if (j.code !== 1000) { A.toast(j.message || '失败'); return; }
        window.DCIM.extra.bills = [j.bill].concat(list);
        A.toast('已登记 ' + j.bill.no);
        AM.views.inout(box);
      });
    });
  };
})(window.App);
