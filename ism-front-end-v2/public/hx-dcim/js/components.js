/* =========================================================
   通用 UI 组件：卡片 / 表格 / 分页 / 统计 / 分布条
   ========================================================= */
(function (A) {
  'use strict';

  var U = {};

  /* ---------- 卡片 ---------- */
  U.card = function (o) {
    o = o || {};
    var extra = o.extra ? '<div class="card-extra">' + o.extra + '</div>' : '';
    return '<section class="card ' + (o.cls || '') + '"' + (o.style ? ' style="' + o.style + '"' : '') + '>' +
      (o.title !== false ? '<div class="card-head"><i class="card-bar"></i>' +
        '<h2 class="card-title">' + A.esc(o.title) + '</h2>' + extra + '</div>' : extra) +
      '<div class="card-body ' + (o.bodyCls || '') + '">' + (o.body || '') + '</div>' +
      '</section>';
  };

  /* ---------- 表格 ----------
     cols: [{k, t, w, cls, fmt(row), html(row), align}]
  */
  U.table = function (o) {
    o = o || {};
    var cols = o.cols || [], rows = o.rows || [];
    var head = cols.map(function (c) {
      return '<th class="' + (c.align === 'r' ? 'num' : '') + '"' +
        (c.w ? ' style="width:' + c.w + 'px"' : '') + '>' + A.esc(c.t) + '</th>';
    }).join('');

    var body = rows.map(function (r, i) {
      var tds = cols.map(function (c) {
        var v;
        if (c.html) v = c.html(r, i);
        else if (c.fmt) v = c.fmt(r, i);
        else v = A.esc(r[c.k]);
        return '<td class="' + (c.align === 'r' ? 'num' : '') + (c.mono ? ' mono' : '') + '">' + v + '</td>';
      }).join('');
      return '<tr data-i="' + i + '"' + (o.rowAttr ? ' ' + o.rowAttr(r, i) : '') + '>' + tds + '</tr>';
    }).join('');

    var empty = rows.length ? '' : '<tr><td colspan="' + cols.length + '"><div class="tbl-empty">暂无数据</div></td></tr>';

    return '<div class="tbl-wrap"><table class="tbl"><thead><tr>' + head + '</tr></thead><tbody>' +
      (body || empty) + '</tbody></table></div>';
  };

  /* ---------- 分页 ---------- */
  U.pager = function (o) {
    o = o || {};
    var p = o.page || 1, pages = o.pages || 1, total = o.total || 0, size = o.size || 20;
    var from = total ? (p - 1) * size + 1 : 0, to = Math.min(total, p * size);
    var nums = [];
    var start = Math.max(1, p - 2), end = Math.min(pages, start + 4);
    start = Math.max(1, Math.min(start, end - 4));
    for (var i = start; i <= end; i++) nums.push(i);

    return '<div class="pager">' +
      '<button class="pg" data-p="1"' + (p === 1 ? ' disabled' : '') + '>«</button>' +
      '<button class="pg" data-p="' + (p - 1) + '"' + (p === 1 ? ' disabled' : '') + '>‹</button>' +
      nums.map(function (n) {
        return '<button class="pg' + (n === p ? ' is-active' : '') + '" data-p="' + n + '">' + n + '</button>';
      }).join('') +
      '<button class="pg" data-p="' + (p + 1) + '"' + (p === pages ? ' disabled' : '') + '>›</button>' +
      '<button class="pg" data-p="' + pages + '"' + (p === pages ? ' disabled' : '') + '>»</button>' +
      '<span class="pager-info">共 ' + A.num(total) + ' 条 · 第 ' + from + '-' + to + ' 条 · ' + pages + ' 页</span>' +
      '</div>';
  };

  /* ---------- 统计块 ---------- */
  U.stats = function (list) {
    return '<div class="stats">' + list.map(function (s) {
      return '<div class="stat-box">' +
        '<span><span class="v"' + (s.color ? ' style="color:' + s.color + '"' : '') + '>' + s.v + '</span>' +
        (s.u ? '<span class="u">' + s.u + '</span>' : '') + '</span>' +
        '<span class="l">' + A.esc(s.l) + '</span></div>';
    }).join('') + '</div>';
  };

  /* ---------- 横向分布条 ---------- */
  U.bars = function (list, o) {
    o = o || {};
    var mx = o.max || Math.max.apply(null, list.map(function (d) { return d.value; }).concat([1]));
    return '<div class="bar-list">' + list.map(function (d) {
      return '<div class="bar-row">' +
        '<span class="bl" title="' + A.esc(d.label) + '">' + A.esc(d.label) + '</span>' +
        '<span class="bar-track"><i class="bar-fill ' + (d.cls || '') + '" style="width:' +
        A.clamp(d.value / mx * 100, 0, 100).toFixed(1) + '%"></i></span>' +
        '<span class="bv">' + A.num(d.value) + (o.unit || '') + '</span>' +
        '</div>';
    }).join('') + '</div>';
  };

  /* ---------- 排行榜 ---------- */
  U.rank = function (list, unit) {
    return '<div class="rank">' + list.map(function (d, i) {
      return '<div class="rank-item"><span class="rank-no">' + (i + 1) + '</span>' +
        '<span class="rank-nm" title="' + A.esc(d.name) + '">' + A.esc(d.name) + '</span>' +
        '<span class="rank-v">' + A.num(d.value) + (unit || '') + '</span></div>';
    }).join('') + '</div>';
  };

  /* ---------- 键值详情 ---------- */
  U.kv = function (pairs) {
    return '<dl class="kv">' + pairs.filter(function (p) {
      return p && p[1] !== '' && p[1] != null;
    }).map(function (p) {
      return '<dt>' + A.esc(p[0]) + '</dt><dd>' + (p[2] ? p[1] : A.esc(p[1])) + '</dd>';
    }).join('') + '</dl>';
  };

  /* ---------- 标签 ---------- */
  U.tag = function (text, type) {
    return '<span class="tag' + (type ? ' is-' + type : '') + '">' + A.esc(text) + '</span>';
  };

  /* ---------- 下拉 ---------- */
  U.select = function (name, options, value, cls) {
    return '<select class="select ' + (cls || '') + '" data-f="' + name + '">' +
      options.map(function (o) {
        var v = (typeof o === 'object') ? o.v : o, t = (typeof o === 'object') ? o.t : o;
        return '<option value="' + A.esc(v) + '"' + (String(v) === String(value) ? ' selected' : '') + '>' + A.esc(t) + '</option>';
      }).join('') + '</select>';
  };

  /* ---------- 空态 ---------- */
  U.empty = function (msg) {
    return '<div class="empty">' + A.esc(msg || '暂无数据') + '</div>';
  };

  A.UI = U;
})(window.App);
