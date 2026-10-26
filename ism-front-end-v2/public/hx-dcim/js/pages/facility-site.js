/* 运行监控附加页签：视频 / 门禁 / 通信 / 消防 */
(function (A) {
  'use strict';

  function rowsOf(list, cols) {
    return '<table class="tbl"><thead><tr>' + cols.map(function (c) { return '<th>' + c.t + '</th>'; }).join('') +
      '</tr></thead><tbody>' + (list || []).map(function (r) {
        return '<tr>' + cols.map(function (c) { return '<td>' + A.esc(r[c.k]) + '</td>'; }).join('') + '</tr>';
      }).join('') + '</tbody></table>';
  }

  function card(title, note, html) {
    return {
      stats: '<div class="muted" style="padding:8px 0">' + A.esc(note) + '</div>',
      body: '<section class="card grow"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">' + A.esc(title) + '</h2></div><div class="card-body">' + html + '</div></section>'
    };
  }

  A.facilityExtra = {
    video: function () {
      var list = ((window.DCIM.extra || {}).video) || [];
      var on = list.filter(function (x) { return x.status === '在线'; }).length;
      return card('视频点位', '在线 ' + on + ' / ' + list.length + '。没有厂商码流时只展示在线状态。',
        rowsOf(list, [{ t: '点位', k: 'name' }, { t: '位置', k: 'room' }, { t: '状态', k: 'status' }]));
    },
    door: function () {
      var list = ((window.DCIM.extra || {}).doors) || [];
      return card('门禁', '异常开门单独标出，事件来自台账。',
        rowsOf(list, [{ t: '门', k: 'name' }, { t: '位置', k: 'room' }, { t: '状态', k: 'status' }, { t: '最近', k: 'last' }]));
    },
    comm: function () {
      var list = ((window.DCIM.extra || {}).comms) || [];
      return card('通信状态', '通道 → 协议 → 丢包与时延。',
        rowsOf(list, [{ t: '通道', k: 'name' }, { t: '协议', k: 'proto' }, { t: '设备数', k: 'deviceN' }, { t: '丢包%', k: 'lossPct' }, { t: '时延ms', k: 'delayMs' }, { t: '状态', k: 'status' }]));
    },
    fire: function () {
      return card('消防联动', '与漏水、烟感告警同一告警流。此处给出联动原则。',
        '<div>烟感或漏水达到紧急等级时，通知值班组并在告警通道里留发送记录。</div>' +
        '<div style="margin-top:8px">现场喷淋、排烟由消防主机执行，本页不直接下发控制。</div>');
    }
  };
})(window.App);
