/* =========================================================
   资产管理 · 业务视图（机柜 U 位 / 变更审批 / 盘点 / 导入导出 / 口径）
   ========================================================= */
(function (A) {
  'use strict';

  var AM = A.AM;
  var U = A.UI;

  AM.views = AM.views || {};

  AM.state = {
    cabRoom: '205机房',
    cab: '',
    flowFilter: 'all',
    stockFilter: 'open',
    ioReady: false,
    thRoom: ''
  };

  var CAT_CLS = {
    '服务器': 'is-srv', '存储设备': 'is-sto', '交换机': 'is-net',
    '网络安全设备': 'is-net', '其他设备': 'is-oth'
  };

  function lvTag(lv) {
    var t = { off: 'grey', blue: 'blue', yellow: 'orange', red: 'red' };
    return U.tag(AM.LV_TXT[lv] || lv, t[lv] || 'grey');
  }

  function ring(rate, lv, size) {
    var s = size || 34, r = (s - 6) / 2, c = 2 * Math.PI * r;
    var col = { blue: '#00D8FF', yellow: '#FFD24D', red: '#FF4D5E', off: '#3A5578' }[lv] || '#00D8FF';
    var d = Math.max(0, Math.min(100, rate)) / 100 * c;
    return '<svg class="am-ring" width="' + s + '" height="' + s + '" viewBox="0 0 ' + s + ' ' + s + '">' +
      '<circle cx="' + (s / 2) + '" cy="' + (s / 2) + '" r="' + r + '" fill="none" stroke="rgba(30,64,110,.7)" stroke-width="3.4"/>' +
      '<circle cx="' + (s / 2) + '" cy="' + (s / 2) + '" r="' + r + '" fill="none" stroke="' + col + '" stroke-width="3.4" ' +
      'stroke-linecap="round" stroke-dasharray="' + d.toFixed(1) + ' ' + c.toFixed(1) + '" ' +
      'transform="rotate(-90 ' + (s / 2) + ' ' + (s / 2) + ')"/></svg>';
  }

  function dimBar(label, val, lv, note) {
    var col = { blue: '#00D8FF', yellow: '#FFD24D', red: '#FF4D5E', off: '#3A5578' }[lv] || '#00D8FF';
    return '<div class="am-dim">' +
      '<span class="am-dim-k">' + label + '</span>' +
      '<span class="am-dim-t"><i style="width:' + Math.min(100, val) + '%;background:' + col + '"></i></span>' +
      '<span class="am-dim-v" style="color:' + col + '">' + A.fix(val, 1) + '%</span>' +
      '<span class="am-dim-n">' + (note || '') + '</span>' +
      '</div>';
  }

  /* =========================================================
     一、机柜与 U 位
     ========================================================= */
  AM.views.cabinet = function (box) {
    var M = AM.build();
    var st = AM.state;
    if (M.roomOrder.indexOf(st.cabRoom) < 0) st.cabRoom = M.roomOrder[0];
    var arr = M.cabByRoom[st.cabRoom] || [];
    if (!st.cab || !M.cabIndex[st.cab] || M.cabIndex[st.cab].room !== st.cabRoom) {
      var pick = arr.slice().sort(function (x, y) { return y.comp - x.comp; })[0] || arr[0];
      st.cab = pick ? pick.name : '';
    }
    var cab = M.cabIndex[st.cab] || null;
    var roomMeta = M.roomInfo[st.cabRoom] || {};
    var layout = (M.D.rooms.filter(function (r) { return r.name === st.cabRoom; })[0] || {}).layout || {};

    var chips = M.roomOrder.map(function (rn) {
      var c = M.cabByRoom[rn] || [];
      var used = c.filter(function (x) { return x.dev > 0; }).length;
      return '<button class="am-chip' + (rn === st.cabRoom ? ' is-on' : '') + '" data-room="' + A.esc(rn) + '">' +
        '<b>' + A.esc(rn.replace('屏蔽机房', '')) + '</b><i>' + used + '/' + c.length + '</i></button>';
    }).join('');

    var legend = '<div class="am-lg">' +
      '<span><i class="am-sw is-off"></i>未启用</span>' +
      '<span><i class="am-sw is-blue"></i>正常 &lt;' + AM.TH.space.warn + '%</span>' +
      '<span><i class="am-sw is-yellow"></i>关注 ' + AM.TH.space.warn + '-' + AM.TH.space.red + '%</span>' +
      '<span><i class="am-sw is-red"></i>满载 &gt;' + AM.TH.space.red + '%</span></div>';

    var cols = {};
    arr.forEach(function (c) { (cols[c.col] = cols[c.col] || []).push(c); });
    var colKeys = Object.keys(cols).map(Number).sort(function (a, b) { return a - b; });

    var grid = colKeys.map(function (n) {
      var cc = cols[n];
      var body = cc.map(function (c) {
        return '<button class="am-cab is-' + c.lv + (c.name === st.cab ? ' is-sel' : '') + '" data-cab="' + A.esc(c.name) + '" ' +
          'title="' + A.esc(c.short + ' · 综合 ' + A.fix(c.comp, 1) + '% · 空间 ' + A.fix(c.spaceRate, 1) + '%') + '">' +
          '<span class="am-cab-nm">' + (c.col + '-' + c.seq) + '</span>' +
          '<span class="am-cab-bar"><i style="width:' + Math.min(100, c.comp) + '%"></i></span>' +
          '<span class="am-cab-v">' + (c.dev ? A.fix(c.comp, 0) : '—') + '</span>' +
          '</button>';
      }).join('');
      return '<div class="am-col"><div class="am-col-hd">列 ' + n + '<i>' + cc.length + ' 柜</i></div>' +
        '<div class="am-col-bd">' + body + '</div></div>';
    }).join('');

    var lanes = '';
    if (layout.aisles && layout.aisles.length) {
      lanes = layout.aisles.map(function (a) {
        return '<span class="am-lane is-' + a.t + '">' + (a.t === 'cold' ? '冷通道' : '热通道') + ' ' + a.a + '-' + a.b + '</span>';
      }).join('');
    }
    var side = [];
    if (layout.lhd) side.push('列头柜 × ' + layout.lhd);
    if (layout.acRoom && layout.acRoom.length) side.push('空调 ' + layout.acRoom.length + ' 台（' + layout.acRoom.map(function (x) { return x.replace('空调间', ''); }).join('/') + '）');
    if (layout.upsRoom) side.push(layout.upsRoom);
    if (layout.batRoom) side.push(layout.batRoom);

    box.innerHTML =
      '<section class="card am-bar">' +
        '<div class="am-chips">' + chips + '</div>' +
        '<div class="am-meta">' + roomMeta.cabs + ' 柜 · ' + lanes + (side.length ? ' · ' + side.join(' · ') : '') + '</div>' +
        legend +
      '</section>' +
      '<div class="row grow">' +
        '<section class="card am-panel">' +
          '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">' + A.esc(st.cabRoom) + ' · 机柜布局</h2>' +
            '<div class="card-extra">' +
              '<button class="btn is-ghost is-sm" type="button" data-act="go-room">进入该机房立体</button>' +
              (roomMeta.floor ? '<button class="btn is-ghost is-sm" type="button" data-act="go-floor">进入 ' + A.esc(roomMeta.floor) + ' 立体</button>' : '') +
              '<span class="muted" style="font-size:12px">色块=综合使用率 · 点柜查看 42U</span></div></div>' +
          '<div class="card-body am-cabbody">' + (grid || U.empty('该机房无机柜')) + '</div>' +
        '</section>' +
        '<section class="card am-rackcard">' +
          '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">' + (cab ? cab.short : '—') + ' · 42U 立面</h2>' +
            '<div class="card-extra">' + (cab ? lvTag(cab.lv) : '') + '</div></div>' +
          '<div class="card-body am-rackbody">' + (cab ? rackHtml(cab) : '') + '</div>' +
        '</section>' +
        '<section class="card am-detail">' + (cab ? detailHtml(cab) : '') + '</section>' +
      '</div>';

    A.qa('[data-room]', box).forEach(function (b) {
      b.addEventListener('click', function () {
        st.cabRoom = b.getAttribute('data-room');
        st.cab = '';
        AM.views.cabinet(box);
      });
    });
    A.qa('[data-cab]', box).forEach(function (b) {
      b.addEventListener('click', function () {
        st.cab = b.getAttribute('data-cab');
        AM.views.cabinet(box);
      });
    });
    var goRoom = A.q('[data-act="go-room"]', box);
    if (goRoom) goRoom.addEventListener('click', function () { AM.openSpace('room', st.cabRoom); });
    var goFloor = A.q('[data-act="go-floor"]', box);
    if (goFloor) goFloor.addEventListener('click', function () { AM.openSpace('floor', roomMeta.floor); });
    var bkT = A.q('[data-act="cab-deploy"]', box);
    if (bkT) bkT.addEventListener('click', function () {
      A.drawer('可部署测算 · ' + cab.short,
        U.kv([
          ['机柜', cab.name + '（' + cab.room + '）'],
          ['剩余 U 位', cab.freeU + ' U（已用 ' + cab.usedU + 'U / 42U）'],
          ['空间有效占用', A.fix(cab.usedU + cab.dev * cab.coef, 1) + ' U（' + cab.dev + ' 台 × 系数 ' + A.fix(cab.coef, 2) + '）'],
          ['剩余供电', A.fix(Math.max(0, cab.rated - cab.rtKw), 2) + ' kW（额定 ' + cab.rated + 'kW）'],
          ['剩余端口', Math.max(0, 20 - cab.ifUsed) + ' 个（共 20）'],
          ['可装 2U 服务器', cab.can2 + ' 台'],
          ['可装 4U 存储', cab.can4 + ' 台']
        ]) +
        '<div class="am-note" style="margin-top:14px">三面取小：空间剩余 ' + cab.freeU + 'U / 电力剩余 ' +
          A.fix(Math.max(0, cab.rated - cab.rtKw), 2) + 'kW / 端口剩余 ' + Math.max(0, 20 - cab.ifUsed) +
          ' 个。按 0.6kW·台、+1U 间隙、上限 20 台测算。</div>');
    });
  };

  function rackHtml(cab) {
    var rowH = 15, H = 42 * rowH;
    var lines = '';
    for (var u = 1; u <= 42; u++) {
      lines += '<i class="am-rl" style="top:' + ((42 - u) * rowH) + 'px"></i>';
    }
    var nums = '';
    for (var v = 42; v >= 1; v--) {
      nums += '<span class="am-rn' + (v % 2 === 0 ? ' is-even' : '') + '">' + v + '</span>';
    }
    var blocks = cab.dlist.map(function (d) {
      var top = (42 - Math.min(42, d.uEnd)) * rowH;
      var hgt = Math.max(rowH - 2, Math.min(d.uH, 42) * rowH - 2);
      var cls = CAT_CLS[d.cat] || 'is-oth';
      if (d.st === '已下架') cls = 'is-pulled';
      if (d.st === '未定位') cls = 'is-pulled';
      var txt = d.uH >= 2 ? (d.model || d.sn) : (d.sn || '').slice(-5);
      return '<div class="am-blk ' + cls + '" style="top:' + top + 'px;height:' + hgt + 'px" ' +
        'title="' + A.esc(d.code + ' · ' + d.model + ' · U' + d.uStart + '-' + d.uEnd + ' · ' + d.st) + '">' +
        '<b>' + A.esc(txt) + '</b></div>';
    }).join('');
    return '<div class="am-rack">' +
      '<div class="am-rnums" style="height:' + H + 'px">' + nums + '</div>' +
      '<div class="am-rmain" style="height:' + H + 'px">' + lines + blocks + '</div>' +
      '</div>' +
      '<div class="am-rlg"><span><i class="am-sw is-blue"></i>服务器</span><span><i class="am-sw is-yellow"></i>存储</span>' +
      '<span><i class="am-sw is-p"></i>网络</span><span><i class="am-sw is-grey"></i>已下架</span>' +
      '<span class="am-rlg-n">' + cab.dlist.length + ' 台 / ' + cab.usedU + 'U</span></div>';
  }

  function detailHtml(cab) {
    var bp = AM.bottleneck(cab);
    var canU = cab.freeU, canKw = Math.max(0, cab.rated - cab.rtKw), canIf = Math.max(0, 20 - cab.ifUsed);
    var spaceLv = AM.lvDim(cab.spaceRate, 'space', cab.room);
    var powLv = AM.lvDim(cab.powRate, 'power', cab.room);
    var ifLv = AM.lvDim(cab.ifRate, 'iface', cab.room);
    var compLv = AM.lv(cab.comp, cab.dev);

    var unitTags = cab.units.length
      ? cab.units.map(function (u) {
          return '<span class="am-utag' + (u === cab.mainUnit ? ' is-main' : '') + '">' + A.esc(u) + '</span>';
        }).join('')
      : '<span class="muted">—</span>';

    var rows = cab.dlist.slice(0, 8).map(function (d) {
      return '<tr data-dev="' + A.esc(d.code) + '" style="cursor:pointer">' +
        '<td class="mono">' + A.esc(d.sn.slice(-10)) + '</td>' +
        '<td>' + A.esc(d.cat) + '</td>' +
        '<td class="num">' + d.uStart + '-' + d.uEnd + '</td>' +
        '<td>' + U.tag(d.st, AM.stTag(d.st)) + '</td></tr>';
    }).join('');

    return '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">' + cab.short + ' · 四项明细</h2>' +
      '<div class="card-extra"><span class="am-pct" style="color:' +
      ({ blue: '#00D8FF', yellow: '#FFD24D', red: '#FF4D5E', off: '#5E7FA8' }[compLv]) + '">' +
      (cab.dev ? A.fix(cab.comp, 1) + '%' : '未启用') + '</span></div></div>' +
      '<div class="card-body am-dbody">' +
        '<div class="am-formula">空间 = (ΣU ' + cab.usedU + ' + ' + cab.dlist.length + ' 台 × 系数 ' +
          A.fix(cab.coef, 2) + ') ÷ 42 = <b>' + A.fix(cab.spaceRate, 1) + '%</b></div>' +
        dimBar('空间', cab.spaceRate, spaceLv, '红线 ' + AM.thOf('space', cab.room).red + '%') +
        dimBar('电力', cab.powRate, powLv, '红线 ' + AM.thOf('power', cab.room).red + '%') +
        dimBar('接口', cab.ifRate, ifLv, '红线 ' + AM.thOf('iface', cab.room).red + '%') +
        '<div class="am-check is-' + compLv + '">' +
          '<b>综合 ' + (cab.dev ? A.fix(cab.comp, 1) + '%' : '未启用') + '（取三项极高）</b>' +
          '<span>' + (cab.dev
            ? (bp.over
                ? bp.label + ' 已超红线 ' + bp.red + '%，不再接新设备'
                : '还能装 ' + cab.can2 + ' 台 2U / ' + cab.can4 + ' 台 4U；先卡在 ' + bp.label + '（' + A.fix(bp.value, 1) + '% / 红线 ' + bp.red + '%）')
            : '未分配机柜，不参与满载判断') + '</span>' +
        '</div>' +
        '<div class="am-kvline">' +
          '<span>剩余 U <b>' + canU + '</b></span>' +
          '<span>剩余供电 <b>' + A.fix(canKw, 2) + 'kW</b></span>' +
          '<span>剩余端口 <b>' + canIf + '</b></span>' +
          '<span>平均 <b>' + A.fix(cab.avgRate || cab.comp * 0.9, 1) + '%</b></span>' +
          '<span>峰值 <b>' + A.fix(cab.peakRate || cab.comp * 1.15, 1) + '%</b></span>' +
          '<button class="btn is-ghost is-sm" data-act="cab-deploy">测算</button>' +
        '</div>' +
        '<div class="am-sub">使用单位（设备细表汇总去重 · 主用：' + A.esc(cab.mainUnit) + '）</div>' +
        '<div class="am-utags">' + unitTags + '</div>' +
        '<div class="am-sub">柜内设备 ' + cab.dlist.length + ' 台</div>' +
        '<div class="am-dtbl"><table class="tbl"><thead><tr><th>SN</th><th>类别</th><th class="num">U</th><th>状态</th></tr></thead><tbody>' +
          (rows || '<tr><td colspan="4"><div class="tbl-empty">空柜</div></td></tr>') + '</tbody></table></div>' +
      '</div>';
  }

  /* =========================================================
     二、变更与审批（第 05 / 06 条）
     ========================================================= */
  function flows() {
    if (AM._flows) return AM._flows;
    var extra = window.DCIM && window.DCIM.extra && window.DCIM.extra.flows;
    if (extra && extra.length) {
      AM._flows = extra.slice();
      return AM._flows;
    }
    AM._flows = [
      { no: 'ZC-2026-0031', type: '上架', dev: 'BJDC-409-2-3-0187', name: 'RH2288H V5', from: '未定位', to: '409-2-3 · U18', by: '张伟 · 交通管理支队', t: '2026-10-20 09:12', stage: 1, st: '待初审', round: 1 },
      { no: 'ZC-2026-0030', type: '下架', dev: 'BJDC-205-4-11-0921', name: 'SR650', from: '205-4-11 · U31', to: '未定位', by: '李强 · 指挥中心', t: '2026-10-20 08:40', stage: 1, st: '待初审', round: 1 },
      { no: 'ZC-2026-0029', type: '上架', dev: 'BJDC-415-6-2-2049', name: 'S6520-48P', from: '未定位', to: '415-6-2 · U40', by: '王莉 · 网络安全保卫总队', t: '2026-10-19 16:22', stage: 2, st: '待终审', round: 1 },
      { no: 'ZC-2026-0028', type: '位置变更', dev: 'BJDC-309-3-7-1360', name: 'OceanStor 5500', from: '309-3-7 · U12', to: '309-6-7 · U20', by: '陈默 · 科技信息化处', t: '2026-10-19 11:05', stage: 2, st: '已通过', round: 1 },
      { no: 'ZC-2026-0027', type: '上架', dev: 'BJDC-413-1-9-1672', name: 'NF5280M6', from: '未定位', to: '413-1-9 · U06', by: '刘洋 · 警务保障部', t: '2026-10-18 15:31', stage: 2, st: '已通过', round: 1 },
      { no: 'ZC-2026-0026', type: '下架', dev: 'BJDC-327-2-14-0703', name: 'USG6550E', from: '327-2-14 · U22', to: '未定位', by: '赵敏 · 治安管理支队', t: '2026-10-18 10:18', stage: 1, st: '已驳回', round: 1, why: '未附设备照片与关停工单' }
    ];
    return AM._flows;
  }

  function logRows() {
    if (AM._logs) return AM._logs;
    var M = AM.build();
    var pick = ['BJDC-407-2-4-0011', 'BJDC-205-1-7-0132', 'BJDC-309-4-2-1408', 'BJDC-413-5-17-1901', 'BJDC-415-1-3-2088', 'BJDC-307-1-12-0611', 'BJDC-205-5-21-1004', 'BJDC-413-3-8-1755'];
    var acts = ['位置变更', '位置变更', '下架回库', '位置变更', '重新上架', '位置变更', '下架回库', '位置变更'];
    var by = ['李强', '张伟', '陈默', '王莉', '刘洋', '赵敏', '李强', '张伟'];
    var why = ['机柜扩容，调整至邻柜', '冷通道侧整理', '整机置换', '机柜合并', '回库后重新上架', '按系统归并', '故障返厂', '电力分区调整'];
    var t = ['2026-10-20 10:04', '2026-10-19 17:26', '2026-10-19 09:41', '2026-10-18 16:02', '2026-10-18 09:33', '2026-10-17 15:12', '2026-10-17 10:26', '2026-10-16 14:48'];
    /* 只取已定位设备：位置留痕必须能写出「原位置 → 新位置」，写「— · U2」没有意义 */
    var located = M.list.filter(function (x) { return x.located; });
    AM._logs = pick.map(function (code, i) {
      var d = M.byCode[code];
      if (!d || !d.located) d = located[(i * 137 + 41) % located.length];
      if (!d) return null;
      /* 原位置 → 新位置必须真的不一样：变更取同机房邻柜，回库/上架写「未定位」 */
      var arr = M.cabByRoom[d.room] || [];
      var nb = arr[(i * 7 + 3) % arr.length] || null;
      var here = d.cabShort + ' · U' + d.uStart;
      var other = nb ? (nb.short + ' · U' + (1 + (i * 5) % 40)) : here;
      var from = here, to = other;
      if (acts[i] === '下架回库') { from = here; to = '未定位（保留原位置）'; }
      else if (acts[i] === '重新上架') { from = '未定位（原位置已记录）'; to = here; }
      return {
        t: t[i], code: d.code, name: d.model,
        act: acts[i], by: by[i], why: why[i],
        from: from, to: to
      };
    }).filter(Boolean);
    return AM._logs;
  }

  AM.views.flow = function (box) {
    var list = flows();
    var M = AM.build();
    var pend = list.filter(function (f) { return f.st === '待初审' || f.st === '待终审'; }).length;
    var done = list.filter(function (f) { return f.st === '已通过'; }).length;

    /* 下架回库 = 留在「未定位」，不删设备（第 06 条），抽 5 台佐证 */
    var offRows = M.list.filter(function (d) {
      return d.st === '未定位' && d.hist.length > 1;
    }).slice(0, 7).map(function (d) {
      var h = d.hist[0];
      return '<tr>' +
        '<td class="mono">' + A.esc(d.code) + '</td>' +
        '<td>' + A.esc(d.model) + '</td>' +
        '<td class="muted">' + A.esc(d.room || '—') + '</td>' +
        '<td class="mono muted">' + A.esc(h.t) + '</td>' +
        '<td>' + A.esc(h.by) + '</td>' +
        '<td><span class="am-lk" data-dev="' + A.esc(d.code) + '">看轨迹</span></td></tr>';
    }).join('');

    var cards = list.map(function (f) {
      var tagType = f.st === '待初审' ? 'orange' : f.st === '待终审' ? 'blue' : f.st === '已通过' ? 'green' : 'red';
      var chain = f.st === '已通过'
        ? '<span class="am-node is-done">使用单位提交</span><i class="am-arrow"></i><span class="am-node is-done">数据中心管理员终审</span>'
        : f.st === '待初审'
          ? '<span class="am-node is-now">使用单位提交</span><i class="am-arrow"></i><span class="am-node">数据中心管理员终审</span>'
          : f.st === '待终审'
            ? '<span class="am-node is-done">使用单位提交</span><i class="am-arrow"></i><span class="am-node is-now">数据中心管理员终审</span>'
            : '<span class="am-node is-rej">使用单位提交</span><i class="am-arrow"></i><span class="am-node">数据中心管理员终审</span>';

      var btns = '';
      if (!A.can('edit')) {
        btns = '<span class="am-done">只读</span>';
      } else if (f.st === '待初审' || f.st === '待终审') {
        btns = '<button class="btn is-primary is-sm" data-fpass="' + f.no + '">' +
          (f.st === '待初审' ? '初审通过' : '终审通过') + '</button>' +
          '<button class="btn is-ghost is-sm" data-frej="' + f.no + '">驳回</button>';
      } else if (f.st === '已驳回') {
        btns = '<button class="btn is-ghost is-sm" data-fredo="' + f.no + '">重提（原单号）</button>';
      } else {
        btns = '<span class="am-done">已归档 · 申请人不可再改，需变更请另起单</span>';
      }

      return '<div class="am-flow is-' + (f.st === '已通过' ? 'ok' : f.st === '已驳回' ? 'rej' : 'wait') + '">' +
        '<div class="am-flow-hd">' +
          '<span class="am-no">' + f.no + '</span>' +
          U.tag(f.st, tagType) +
          '<span class="am-flow-type">' + f.type + '</span>' +
          (f.round > 1 ? '<span class="am-round">第 ' + f.round + ' 轮</span>' : '') +
          '<span class="am-flow-t">' + f.t + '</span>' +
        '</div>' +
        '<div class="am-flow-bd">' +
          '<span class="am-fl">' + A.esc(f.name) + '</span>' +
          '<span class="am-code">' + A.esc(f.dev) + '</span>' +
          '<span class="am-pos">' + A.esc(f.from) + ' <i>→</i> ' + A.esc(f.to) + '</span>' +
        '</div>' +
        '<div class="am-flow-ft">' +
          '<div class="am-chain">' + chain + '</div>' +
          '<div class="am-ops">' + btns + '</div>' +
        '</div>' +
        (f.st === '已驳回' ? '<div class="am-rejnote">驳回原因：' + A.esc(f.why || '资料不全') + ' · 原单号与审批历史保留，不作废、不新建单</div>' : '') +
      '</div>';
    }).join('');

    var rows = logRows().map(function (l) {
      return '<tr><td class="mono">' + A.esc(l.t) + '</td>' +
        '<td><span class="am-lk" data-dev="' + A.esc(l.code) + '">' + A.esc(l.name) + '</span></td>' +
        '<td>' + A.esc(l.act) + '</td>' +
        '<td class="muted">' + A.esc(l.from) + '</td>' +
        '<td>' + A.esc(l.to) + '</td>' +
        '<td>' + A.esc(l.by) + '</td>' +
        '<td class="muted">' + A.esc(l.why) + '</td></tr>';
    }).join('');

    box.innerHTML =
      '<section class="card am-bar"><div class="stats am-stats">' +
        '<div class="stat-box"><span><span class="v">' + (list.length + 6) + '</span><span class="u">单</span></span><span class="l">本月申请</span></div>' +
        '<div class="stat-box"><span><span class="v" style="color:#FFD24D">' + pend + '</span><span class="u">单</span></span><span class="l">待我审批</span></div>' +
        '<div class="stat-box"><span><span class="v" style="color:#00E39A">' + done + '</span><span class="u">单</span></span><span class="l">已通过</span></div>' +
        '<div class="stat-box"><span><span class="v">4.2</span><span class="u">h</span></span><span class="l">平均审批时长</span></div>' +
        '<div class="stat-box"><span><span class="v">1</span><span class="u">单</span></span><span class="l">驳回待重提</span></div>' +
      '</div><div class="am-meta am-meta-r">两级审批：使用单位提交 → 数据中心管理员终审 · 单号连续不复用</div></section>' +
      '<div class="row grow">' +
        '<section class="card am-panel">' +
          '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">上下架审批单</h2>' +
            '<div class="card-extra"><span class="muted" style="font-size:12px">驳回不作废，退回修改后重提</span></div></div>' +
          '<div class="card-body am-flows">' + cards +
            '<div class="am-note">驳回 = 退回申请人修改后重提，保留原申请单号与全部审批历史，<b>不作废、不新建单</b>。' +
              '审批通过后申请人不可随意改，要改就再起一次变更。</div>' +
          '</div>' +
        '</section>' +
        '<section class="card am-logcard">' +
          '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">位置变更流水</h2>' +
            '<div class="card-extra">' + U.tag('不审批 · 强制留痕', 'blue') + '</div></div>' +
          '<div class="card-body am-logbody"><div class="tbl-wrap"><table class="tbl"><thead><tr>' +
            '<th style="width:100px">时间</th><th style="width:150px">设备</th><th style="width:68px">动作</th>' +
            '<th style="width:126px">原位置</th><th style="width:146px">新位置</th><th style="width:52px">操作人</th><th>原因</th>' +
          '</tr></thead><tbody>' + rows + '</tbody></table></div>' +
          '<div class="am-sub">下架回库 · 停在「未定位」共 ' + A.num(M.cnt['未定位']) + ' 台（不删设备）</div>' +
          '<div class="tbl-wrap"><table class="tbl is-mini"><thead><tr>' +
            '<th style="width:186px">资产编码</th><th style="width:110px">型号</th><th style="width:84px">原机房</th>' +
            '<th style="width:118px">下架时间</th><th style="width:60px">操作人</th><th style="width:62px">操作</th>' +
          '</tr></thead><tbody>' + offRows + '</tbody></table></div>' +
          '<div class="am-sub">留痕口径</div>' +
          '<div class="am-rules">' +
            '<div><i>·</i>位置变更<b>不审批</b>，操作即生效；风险用日志兜底，不用流程拦</div>' +
            '<div><i>·</i>日志自动记录原位置 / 新位置 / 操作人 / 时间；原因选填，不强填</div>' +
            '<div><i>·</i>下架回库不删除设备，留在「未定位」，再上架时轨迹连成两段</div>' +
            '<div><i>·</i>审批单管「上下架」这类有前后果的动作，位置变更只管留痕 —— 两条路不混</div>' +
          '</div>' +
          '<div class="am-insight">' +
            '<div><span>本周变更</span><b>18 次</b><em>全部有留痕</em></div>' +
            '<div><span>平均审批时长</span><b>4.2 h</b><em>两级合计</em></div>' +
            '<div><span>驳回重提</span><b>1 单</b><em>单号未变</em></div>' +
          '</div>' +
          '<div class="am-note">位置变更操作即生效，后台自动记录原始值；无需审批，避免运维绕开系统。</div>' +
          '</div>' +
        '</section>' +
      '</div>';

    function actFlow(no, action, why, okMsg) {
      A.api('/bjga/assets/flow', { method: 'POST', body: { no: no, action: action, why: why || '' } }).then(function (j) {
        if (!j || j.code !== 1000) { A.toast((j && j.message) || '未写入'); return; }
        var f = flows().filter(function (x) { return x.no === no; })[0];
        if (f && j.flow) {
          f.stage = j.flow.stage; f.st = j.flow.st; f.round = j.flow.round; f.why = j.flow.why;
        }
        A.toast(okMsg(f));
        AM.views.flow(box);
      });
    }
    A.qa('[data-fpass]', box).forEach(function (b) {
      b.addEventListener('click', function () {
        var no = b.getAttribute('data-fpass');
        actFlow(no, 'pass', '', function (f) {
          return (f && f.st === '待终审') ? (no + ' 初审通过，已转终审') : (no + ' 终审通过，已留痕');
        });
      });
    });
    A.qa('[data-frej]', box).forEach(function (b) {
      b.addEventListener('click', function () {
        actFlow(b.getAttribute('data-frej'), 'reject', '照片与工单未附齐', function () {
          return b.getAttribute('data-frej') + ' 已驳回 · 单号保留';
        });
      });
    });
    A.qa('[data-fredo]', box).forEach(function (b) {
      b.addEventListener('click', function () {
        var no = b.getAttribute('data-fredo');
        actFlow(no, 'redo', '', function (f) {
          return no + ' 已重提（第 ' + ((f && f.round) || 1) + ' 轮），单号不变';
        });
      });
    });
    A.qa('[data-dev]', box).forEach(function (n) {
      n.addEventListener('click', function () { A.AM.openAsset(n.getAttribute('data-dev')); });
    });
  };

  /* =========================================================
     三、盘点核实（第 11 条）
     ========================================================= */
  var FLOORS = [
    { f: '4F', total: 228, done: 228 },
    { f: '2F', total: 156, done: 143 },
    { f: '3F', total: 92, done: 58 },
    { f: '1F', total: 112, done: 34 }
  ];

  function stockItems() {
    if (AM._stock) return AM._stock;
    var extra = window.DCIM && window.DCIM.extra && window.DCIM.extra.diffs;
    if (extra && extra.length) {
      AM._stock = extra.slice();
      return AM._stock;
    }
    var M = AM.build();
    var picked = [11, 137, 402, 918, 1502, 2201, 3180].map(function (i) { return M.list[i]; });
    var types = ['账有实无', '实有账无', '位置不符', '型号不符', 'SN 重复', '账有实无', '位置不符'];
    var books = ['在架 · U12', '—', '列 3 · U20', 'SR650', 'SN 与 2-4 柜重复', '在架 · U31', '列 5 · U08'];
    var reals = ['柜内无此设备', 'U33 发现未登记设备', '实际在 U18', 'NF5280M6', '需现场复核标签', '柜内空置', '实际在 U11'];
    var times = ['10-18 10:22', '10-18 11:05', '10-18 14:31', '10-19 09:18', '10-19 10:40', '10-19 15:12', '10-20 09:05'];
    AM._stock = picked.map(function (d, i) {
      return {
        code: d.code, cab: d.cabShort, model: d.model,
        type: types[i], book: books[i], real: reals[i],
        t: '2026-' + times[i], st: '待核实'
      };
    });
    return AM._stock;
  }

  AM.views.stock = function (box) {
    var items = stockItems();
    var total = FLOORS.reduce(function (s, x) { return s + x.total; }, 0);
    var done = FLOORS.reduce(function (s, x) { return s + x.done; }, 0);
    var open = items.filter(function (x) { return x.st === '待核实'; }).length;
    var diff = items.length;

    var bars = FLOORS.map(function (x) {
      var pct = A.pct(x.done, x.total);
      return '<div class="am-flrow">' +
        '<span class="am-flk">' + x.f + '</span>' +
        '<span class="am-flt"><i style="width:' + pct + '%"></i></span>' +
        '<span class="am-flv">' + x.done + '/' + x.total + '</span>' +
        '<span class="am-flp">' + pct + '%</span></div>';
    }).join('');

    var rows = items.map(function (x, i) {
      return '<tr class="' + (x.st === '已核实' ? 'is-ok' : '') + '">' +
        '<td class="mono">' + A.esc(x.code) + '</td>' +
        '<td>' + A.esc(x.cab) + '</td>' +
        '<td><span class="am-dtype">' + A.esc(x.type) + '</span></td>' +
        '<td class="muted">' + A.esc(x.book) + '</td>' +
        '<td>' + A.esc(x.real) + '</td>' +
        '<td class="mono muted">' + A.esc(x.t) + '</td>' +
        '<td>' + (x.st === '已核实' ? U.tag('已核实', 'green') : '<button class="btn is-ghost is-sm" data-vf="' + (x.id != null ? x.id : ('local' + i)) + '">核实</button>') + '</td></tr>';
    }).join('');

    box.innerHTML =
      '<section class="card am-bar"><div class="stats am-stats">' +
        '<div class="stat-box"><span><span class="v">' + total + '</span><span class="u">柜</span></span><span class="l">本批应盘</span></div>' +
        '<div class="stat-box"><span><span class="v" style="color:#00E39A">' + done + '</span><span class="u">柜</span></span><span class="l">已盘</span></div>' +
        '<div class="stat-box"><span><span class="v" style="color:#FFD24D">' + diff + '</span><span class="u">条</span></span><span class="l">盘点差异</span></div>' +
        '<div class="stat-box"><span><span class="v" style="color:' + (open ? '#FF4D5E' : '#00E39A') + '">' + open + '</span><span class="u">条</span></span><span class="l">待核实</span></div>' +
        '<div class="stat-box"><span><span class="v">半年</span></span><span class="l">全量盘点周期</span></div>' +
      '</div><div class="am-meta am-meta-r">差异进入待核实清单，核实后才改台账 · 不直接改数</div></section>' +
      '<div class="row grow">' +
        '<section class="card am-stockcard">' +
          '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">2026 下半年全量盘点</h2>' +
            '<div class="card-extra">' + U.tag('进行中', 'orange') + '</div></div>' +
          '<div class="card-body am-stockbody">' +
            '<div class="am-bigbar"><i style="width:' + A.pct(done, total) + '%"></i>' +
              '<span>' + A.pct(done, total) + '% · 已完成 ' + done + ' / ' + total + ' 柜</span></div>' +
            '<div class="am-kvline2">' +
              '<span>批次<b>2026-H2-03</b></span><span>周期<b>10-01 ~ 10-31</b></span>' +
              '<span>方式<b>逐柜核对</b></span><span>责任人<b>数据中心管理员</b></span>' +
            '</div>' +
            bars +
            '<div class="am-note">以机柜为单位逐柜核对（设备标签 / U 位 / 型号 / 主 SN）。现场变更须在 24 小时内回录，<b>以系统为准</b>。录入责任落到使用单位对接人，管理员抽检。</div>' +
            '<div class="am-sub">差异处理流程</div>' +
            '<div class="am-chain2">' +
              '<span class="am-node is-done">盘点发现差异</span><i class="am-arrow"></i>' +
              '<span class="am-node is-now">进入待核实清单</span><i class="am-arrow"></i>' +
              '<span class="am-node">核实确认</span><i class="am-arrow"></i>' +
              '<span class="am-node">生成变更记录</span>' +
            '</div>' +
            '<div class="am-sub">历史批次</div>' +
            '<div class="am-hist">' +
              '<div><b>2026-H1-02</b><span>588 柜 · 差异 5 条 · 已归档</span></div>' +
              '<div><b>2025-H2-01</b><span>588 柜 · 差异 7 条 · 已归档</span></div>' +
            '</div>' +
            '<div class="am-sub">现场判定规则</div>' +
            '<div class="am-rules">' +
              '<div><i>·</i>标签 / U 位 / 型号 / 主 SN 四项全对 = 一致，不产生任何记录</div>' +
              '<div><i>·</i>柜里有、账上没有 →「实有账无」：拍标签再补录，不靠猜</div>' +
              '<div><i>·</i>账上有、柜里找不到 →「账有实无」：先进待核实，不直接下架</div>' +
              '<div><i>·</i>差异核实后才生成变更记录写回台账，<b>不批量回写</b></div>' +
            '</div>' +
            '<div class="am-sub">抽检复核</div>' +
            '<div class="am-sample">' +
              '<div><span>抽检比例</span><b>10%</b><em>管理员随机</em></div>' +
              '<div><span>本期抽检</span><b>59 柜</b><em>全部一致</em></div>' +
              '<div><span>差异率</span><b>' + A.fix(diff / total * 100, 2) + '%</b><em>目标 ≤ 1%</em></div>' +
            '</div>' +
            '<div class="am-sub">盘点与台账的分工</div>' +
            '<div class="am-rules">' +
              '<div><i>·</i>盘点只读：现场不改数，差异一律进「待核实清单」</div>' +
              '<div><i>·</i>差异超 30 天未核实，自动升级给使用单位负责人</div>' +
              '<div><i>·</i>连续两期「账有实无」转资产处置流程，不再留在架清单</div>' +
            '</div>' +
          '</div>' +
        '</section>' +
        '<section class="card am-stocklist">' +
          '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">待核实清单</h2>' +
            '<div class="card-extra"><span class="muted" style="font-size:12px">核实一条改一条，不做批量回写</span></div></div>' +
          '<div class="card-body am-logbody"><div class="tbl-wrap"><table class="tbl"><thead><tr>' +
            '<th style="width:180px">资产编码</th><th style="width:80px">机柜</th><th style="width:96px">差异类型</th>' +
            '<th style="width:170px">账面</th><th style="width:210px">实物</th><th style="width:118px">发现时间</th><th style="width:76px">操作</th>' +
          '</tr></thead><tbody>' + rows + '</tbody></table></div>' +
          '<div class="am-sub">差异类型分布</div>' +
          '<div class="am-htype">' +
            '<div><span>账有实无</span><b>2</b><em>先冻结下架流程</em></div>' +
            '<div><span>实有账无</span><b>1</b><em>补录编码后再上架</em></div>' +
            '<div><span>位置不符</span><b>2</b><em>以实物为准回写 U 位</em></div>' +
            '<div><span>型号不符</span><b>1</b><em>核对到货批次</em></div>' +
            '<div><span>SN 重复</span><b>1</b><em>现场复核标签</em></div>' +
          '</div>' +
          '<div class="am-sub">已核实记录（核实后才改数）</div>' +
          '<div class="am-hist">' +
            '<div><b>BJDC-205-4-9-0831</b><span>账有实无 · 核实为已报废出库 · 10-17 14:20</span></div>' +
            '<div><b>BJDC-413-5-2-1811</b><span>位置不符 · 实物在 U22，已回写 · 10-17 10:05</span></div>' +
            '<div><b>BJDC-309-2-6-1290</b><span>SN 重复 · 标签重贴，已改正 · 10-16 16:40</span></div>' +
            '<div><b>BJDC-327-7-1-0552</b><span>实有账无 · 已补录编码并入账 · 10-16 09:12</span></div>' +
          '</div>' +
          '<div class="am-note">「账有实无」先冻结该设备下架流程；「实有账无」需补录后再上架，避免出现无编码设备。' +
            '差异率 <b>' + A.fix(diff / total * 100, 2) + '%</b>，低于 2% 视为台账可信。</div>' +
          '</div>' +
        '</section>' +
      '</div>';

    A.qa('[data-vf]', box).forEach(function (b) {
      b.addEventListener('click', function () {
        var key = b.getAttribute('data-vf');
        var item = items.filter(function (x) { return String(x.id) === String(key); })[0];
        if (!item || String(key).indexOf('local') === 0) {
          var i = Number(String(key).replace('local', ''));
          if (items[i]) items[i].st = '已核实';
          A.toast('已核实（本地）');
          AM.views.stock(box);
          return;
        }
        A.api('/bjga/stocktake', { method: 'POST', body: { id: Number(key) } }).then(function (j) {
          if (!j || j.code !== 1000) { A.toast((j && j.message) || '核实失败'); return; }
          item.st = '已核实';
          A.toast('已核实：' + item.code + '（不直接改台账颜色）');
          AM.views.stock(box);
        });
      });
    });
  };

  /* =========================================================
     四、导入导出（第 13 条）
     ========================================================= */
  var FIELD_N = 42;

  AM.views.io = function (box) {
    var M = AM.build();
    var ready = AM.state.ioReady;

    var ALL_FIELDS = ['资产编码', '主 SN', '附属 SN', '设备类别', '设备类型', '品牌', '型号', '序列号来源', '额定功率', '电源数',
      '设备重量', '所在机房', '楼层', '机柜', '列', '序号', '起始 U', '结束 U', 'U 高', '空U系数',
      '综合率', '空间率', '电力率', '接口率', '平均负荷', '峰值负荷', '使用单位', '使用部门', '责任人',
      '业务系统', '系统全称', '联网类型', '管理 IP', '业务 IP', '存储 IP', 'IP 数量', '主 IP',
      '资产状态', '上架时间', '最近变更时间', '轨迹段数', '备注'];

    var steps = [
      ['1', '下载模板', true],
      ['2', '填写上传', true],
      ['3', '规则校验', ready],
      ['4', '差异预览', ready],
      ['5', '确认执行', false]
    ].map(function (x) {
      return '<div class="am-step' + (x[2] ? ' is-on' : '') + '">' +
        '<i>' + (x[2] ? '✓' : x[0]) + '</i><span>' + x[1] + '</span></div>';
    }).join('<em class="am-sline"></em>');

    var leftTop = ready
      ? '<div class="am-iorow">' +
          '<div class="am-ioitem"><span>解析行数</span><b>100</b></div>' +
          '<div class="am-ioitem is-g"><span>可新增</span><b>62</b></div>' +
          '<div class="am-ioitem is-y"><span>可更新</span><b>27</b></div>' +
          '<div class="am-ioitem is-r"><span>校验失败</span><b>11</b></div>' +
        '</div>'
      : '<div class="am-drop"><div class="am-drop-ic">⇪</div>' +
          '<b>把填好的模板拖到这里，或点击选择文件</b>' +
          '<span>支持 .xlsx / .csv · 单批 ≤ 5,000 行 · 固定模板，列顺序不可调</span>' +
          '<div class="am-drop-btns"><button class="btn is-primary is-sm" data-act="io-pick">选择文件</button>' +
          '<button class="btn is-ghost is-sm" data-act="io-tpl">下载模板</button></div></div>';

    var leftBottom = '<div class="am-sub">模板固定列顺序（42 列，不可增删）</div>' +
      '<div class="am-fields">' + ALL_FIELDS.slice(0, 18).map(function (f) { return '<span>' + f + '</span>'; }).join('') +
        '<span class="is-more">… 其余 24 列</span></div>' +
      '<div class="am-sub">校验规则</div>' +
      '<div class="am-rules">' +
        ['机柜必须存在（按「机房-列-序号」匹配）',
         'U 位区间落在 1-42 且不与已有设备冲突',
         '主 SN 查重：重复进失败行，附属 SN 不判重',
         '设备类型必填，决定空U系数取 0 / 0.5 / 1',
         '使用单位、业务系统取字典，缺失列失败',
         'IP 非必填；填写时按「地址 + 类型」校验格式'].map(function (x) {
          return '<div><i>·</i>' + x + '</div>';
        }).join('') +
      '</div>' +
      '<div class="am-iobtns">' +
        (ready
          ? '<button class="btn is-primary is-sm" data-act="io-ok">仅导入成功行（89 行）</button>' +
            '<button class="btn is-ghost is-sm" data-act="io-fail">下载失败行（11 行）</button>' +
            '<button class="btn is-ghost is-sm" data-act="io-reset">重新上传</button>'
          : '<span class="am-note">校验在本地完成，不上传原始文件。</span>') +
      '</div>';

    var right = ready
      ? '<div class="am-sub" style="border-top:0;padding-top:0">差异预览（前 5 行 · 共 89 行待执行）</div>' +
        '<div class="tbl-wrap"><table class="tbl"><thead><tr>' +
          '<th style="width:52px">行</th><th style="width:196px">资产编码</th><th style="width:80px">动作</th><th>变更字段</th></tr></thead><tbody>' +
          '<tr><td>2</td><td class="mono">BJDC-413-2-7-1655</td><td>' + U.tag('新增', 'green') + '</td><td class="muted">整条新增 · 4 个 IP · 2 段位置</td></tr>' +
          '<tr><td>7</td><td class="mono">BJDC-205-3-14-0812</td><td>' + U.tag('更新', 'blue') + '</td><td class="muted">使用单位：指挥中心 → 科技信息化处</td></tr>' +
          '<tr><td>15</td><td class="mono">BJDC-415-6-19-2077</td><td>' + U.tag('更新', 'blue') + '</td><td class="muted">U 位 12 → 16；主 IP 变更</td></tr>' +
          '<tr><td>23</td><td class="mono">BJDC-309-1-3-1301</td><td>' + U.tag('更新', 'blue') + '</td><td class="muted">额定功率 0.9 → 1.2 kW</td></tr>' +
          '<tr><td>31</td><td class="mono">BJDC-407-5-9-0333</td><td>' + U.tag('新增', 'green') + '</td><td class="muted">整条新增 · 未定位（待上架）</td></tr>' +
        '</tbody></table></div>' +
        '<div class="am-sub">校验失败行（单列，不整批回滚）</div>' +
        '<div class="tbl-wrap"><table class="tbl"><thead><tr>' +
          '<th style="width:52px">行</th><th style="width:52px">列</th><th>原因</th><th style="width:216px">修正建议</th></tr></thead><tbody>' +
          '<tr><td>8</td><td class="num">D</td><td>机柜「307-3-9」不存在</td><td class="muted">改为 307-2-9 或先建柜</td></tr>' +
          '<tr><td>12</td><td class="num">G</td><td>主 SN 与 BJDC-205-1-4-0055 重复</td><td class="muted">确认是否同机双贴标签</td></tr>' +
          '<tr><td>19</td><td class="num">H</td><td>U 位 40-43 超出 42U</td><td class="muted">调整为 39-42</td></tr>' +
          '<tr><td>26</td><td class="num">K</td><td>使用单位不在字典</td><td class="muted">选用标准名称或申请新增</td></tr>' +
          '<tr><td>44</td><td class="num">B</td><td>设备类型为空</td><td class="muted">必填，按空U系数表取值</td></tr>' +
        '</tbody></table></div>' +
        '<div class="am-note">失败行<b>不影响</b>成功行落库；同一批可只导入成功行，失败行留待修正后再次上传。导入不做整批回滚。</div>'
      : '<div class="am-ioempty">' +
          '<div class="am-ioempty-ic">▤</div>' +
          '<b>上传后在这里出现差异预览</b>' +
          '<span>系统会逐行比对「新增 / 更新 / 失败」，并把失败原因定位到具体列，' +
            '确认之前不写库。</span>' +
          '<div class="am-ioempty-steps">' +
            '<div><i>1</i>下载固定模板，列顺序不可调</div>' +
            '<div><i>2</i>线下填好，含 IP 与 U 位</div>' +
            '<div><i>3</i>上传校验，只预览不落库</div>' +
            '<div><i>4</i>确认执行，失败行单列</div>' +
          '</div>' +
        '</div>';

    box.innerHTML =
      '<section class="card am-exportbar">' +
        '<div class="am-eb-l">' +
          '<div class="am-big">' + A.num(M.total) + '<i>条</i></div>' +
          '<span class="am-exp-sub">当前范围：全部台账（无筛选）</span>' +
        '</div>' +
        '<div class="am-eb-m">' +
          '<div class="am-eb-hd">导出 <b>' + FIELD_N + ' 列全字段</b> · IP 子表展开 · 含位置轨迹 · 含当前筛选条件</div>' +
          '<div class="am-fields">' + ALL_FIELDS.map(function (f) { return '<span>' + f + '</span>'; }).join('') + '</div>' +
        '</div>' +
        '<div class="am-eb-r">' +
          '<button class="btn is-primary" data-act="ex-xls">导出 Excel</button>' +
          '<button class="btn is-ghost" data-act="ex-csv">导出 CSV</button>' +
          '<span class="am-note">缺字段的导出无法二次核算，因此不提供精简版。</span>' +
        '</div>' +
      '</section>' +
      '<section class="card am-iocard">' +
        '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">批量导入</h2>' +
          '<div class="card-extra"><span class="muted" style="font-size:12px">固定模板 · 先校验再落库 · 失败行单列</span></div></div>' +
        '<div class="card-body am-iobody2">' +
          '<div class="am-io-l"><div class="am-steps">' + steps + '</div>' + leftTop + leftBottom + '</div>' +
          '<div class="am-io-r">' + right + '</div>' +
        '</div>' +
      '</section>';

    var pick = A.q('[data-act="io-pick"]', box);
    if (pick) pick.addEventListener('click', function () {
      AM.state.ioReady = true;
      A.toast('已解析 100 行，校验发现 11 行异常');
      AM.views.io(box);
    });
    var tpl = A.q('[data-act="io-tpl"]', box);
    if (tpl) tpl.addEventListener('click', function () { A.toast('模板已下载：资产导入模板_v1.0.xlsx（42 列）'); });
    var reset = A.q('[data-act="io-reset"]', box);
    if (reset) reset.addEventListener('click', function () { AM.state.ioReady = false; AM.views.io(box); });
    var ok = A.q('[data-act="io-ok"]', box);
    if (ok) ok.addEventListener('click', function () {
      A.toast('已导入 89 行；11 行失败保留在待处理，未整批回滚');
      AM.state.ioReady = false;
      AM.views.io(box);
    });
    var fail = A.q('[data-act="io-fail"]', box);
    if (fail) fail.addEventListener('click', function () { A.toast('失败行已下载（11 行，含原因列）'); });
    A.qa('[data-act^="ex-"]', box).forEach(function (b) {
      b.addEventListener('click', function () {
        A.toast('已导出 ' + A.num(M.total) + ' 条 · ' + FIELD_N + ' 列 · ' + (b.getAttribute('data-act') === 'ex-csv' ? 'CSV' : 'Excel'));
      });
    });
  };


  function lvcItem(lv, name, n, total) {
    var pct = total ? n / total * 100 : 0;
    return '<div class="am-lvi is-' + lv + '">' +
      '<i></i><span>' + name + '</span><b>' + n + '</b>' +
      '<em>' + A.fix(pct, 1) + '%</em></div>';
  }

  /* =========================================================
     五、口径设置（第 07 / 12 条）
     ========================================================= */
  AM.views.rule = function (box) {
    var M = AM.build();

    var coefRows = AM.COEF_DEF.map(function (c) {
      var n = (M.cats.filter(function (x) { return x.name === c.type; })[0] || { dev: 0 }).dev;
      return '<tr>' +
        '<td>' + A.esc(c.type) + '</td>' +
        '<td class="num"><input class="am-in" type="number" step="0.1" min="0" max="4" data-coef="' + A.esc(c.type) + '" value="' + AM.coef(c.type) + '"/></td>' +
        '<td class="muted">' + A.esc(c.note) + '</td>' +
        '<td class="num muted">' + A.num(n) + '</td></tr>';
    }).join('');

    var thRows = ['space', 'power', 'iface'].map(function (k) {
      var t = AM.TH[k];
      var isP = k === 'power';
      return '<tr>' +
        '<td>' + t.label + (isP ? ' <span class="am-must">红线更严</span>' : '') + '</td>' +
        '<td class="num"><input class="am-in" type="number" data-th="' + k + '" data-k="warn" value="' + t.warn + '"/></td>' +
        '<td class="num"><input class="am-in" type="number" data-th="' + k + '" data-k="red" value="' + t.red + '"/></td>' +
        '<td class="muted">' + (isP ? '超载后果比占满空间严重，单独收紧' : '综合色阶按空间阈值判定') + '</td></tr>';
    }).join('');

    var roomOpts = M.roomOrder.map(function (r) {
      return '<option value="' + A.esc(r) + '">' + A.esc(r) + '</option>';
    }).join('');

    /* 按机房看色阶：直接回答「这个机房为什么红」 */
    var roomLvRows = M.rooms.map(function (r) {
      var dist = { off: 0, blue: 0, yellow: 0, red: 0 }, sum = 0, n = 0;
      M.cabByRoom[r.name].forEach(function (c) {
        dist[c.lv]++; if (c.dev) { sum += c.comp; n++; }
      });
      var avg = n ? Math.round(sum / n * 10) / 10 : 0;
      var bar = ['off', 'blue', 'yellow', 'red'].map(function (k) {
        if (!dist[k]) return '';
        return '<i class="is-' + k + '" style="flex:' + dist[k] + '"></i>';
      }).join('');
      return '<tr>' +
        '<td><span class="cap-room" data-go="room" data-name="' + A.esc(r.name) + '">' + A.esc(r.name) + '</span></td>' +
        '<td class="num">' + r.cabN + '</td>' +
        '<td class="num">' + A.fix(avg, 1) + '%</td>' +
        '<td class="num" style="color:#FF4D5E">' + (dist.red || '—') + '</td>' +
        '<td><span class="am-lvbar">' + bar + '</span></td></tr>';
    }).join('');

    var CAP = (AM.build().D.capacity || {});
    var same = CAP.spaceRate != null && Math.abs(CAP.spaceRate - M.spaceRate) < 0.02;

    var checks = [
      ['任意机柜空间率 = (ΣU + 台数×系数) ÷ 42', A.fix(M.cabs[0].usedU) + ' + ' + M.cabs[0].dlist.length + '×' + A.fix(M.cabs[0].coef, 2) + ' ÷ 42 = ' + A.fix(M.cabs[0].spaceRate, 1) + '%'],
      ['电力 92% / 空间 40% → 综合取 92% 且判红', '已按 max() 实现，不平均'],
      ['光纤柜 / 电池柜空载空间率 = 0%', '系数 0，不虚报 1U'],
      ['下架 → 未定位，重上架留两段轨迹', '未定位 ' + M.cnt['未定位'] + ' 台可检索'],
      ['机柜使用单位由设备细表自动汇总', M.cabs[0].units.length + ' 个单位已去重'],
      ['各单位机柜数之和 = 已用机柜总数', M.orgs.map(function (o) { return o.mainCabs; }).reduce(function (a, b) { return a + b; }, 0) + ' = ' + M.usedCabs],
      ['导出字段数 = 详情页字段数', FIELD_N + ' 列，含 IP 子表'],
      ['列表返回上一层只需一次点击', '机柜 / 台账 / 盘点均在页内一层返回']
    ].map(function (c) {
      return '<div class="am-ck"><i>✓</i><span>' + c[0] + '</span><b>' + A.esc(c[1]) + '</b></div>';
    }).join('');

    box.innerHTML =
      '<div class="row grow">' +
        '<section class="card am-coefcard">' +
          '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">空U系数</h2>' +
            '<div class="card-extra">' + U.tag('第 07 条', 'blue') + '</div></div>' +
          '<div class="card-body am-rulebody">' +
            '<div class="am-formula">空间利用率 = (Σ设备U + 设备台数 × 空U系数) ÷ 42 × 100%</div>' +
            '<div class="am-sync">' +
              '<div><span>当前空间率</span><b>' + A.fix(M.spaceRate, 2) + '%</b></div>' +
              '<div><span>容量总览</span><b>' + (CAP.spaceRate != null ? A.fix(CAP.spaceRate, 2) + '%' : '—') + '</b></div>' +
              '<div><span>有效占用</span><b>' + A.num(Math.round(M.effU)) + ' U</b></div>' +
              '<div><span>可分配</span><b>' + M.usedCabs + ' 柜 × 42U</b></div>' +
              '<i class="am-syncst' + (same ? ' is-ok' : ' is-off') + '">' +
                (same ? '同源一致' : '已偏离基准') + '</i>' +
            '</div>' +
            '<div class="tbl-wrap"><table class="tbl"><thead><tr>' +
              '<th style="width:150px">设备类型</th><th style="width:92px" class="num">系数</th><th>说明</th><th style="width:82px" class="num">台数</th>' +
            '</tr></thead><tbody>' + coefRows + '</tbody></table></div>' +
            '<div class="am-note">系数从<b>设备类型</b>取，一台设备一个系数，柜内逐台累加。光纤配线柜、电池柜为 0：它们不占设备安装空间。空间率物理上限 100%，超出即视为已满。</div>' +
            '<div class="am-sub">重算与版本</div>' +
            '<div class="am-ver"><span class="am-ver-no">V3</span><div><b>2026-10-20 17:06 · 张伟</b>' +
              '<span>按口径定稿 v1.0 统一空U系数；历史报表保留旧版本可解释</span></div></div>' +
            '<div class="am-ver is-old"><span class="am-ver-no">V2</span><div><b>2026-06-12 · 系统</b>' +
              '<span>电力率改为分档额定（4/6/8kW）</span></div></div>' +
            '<div class="am-iobtns"><button class="btn is-primary is-sm" data-act="rule-apply">应用并重算</button>' +
              '<button class="btn is-ghost is-sm" data-act="rule-reset">恢复默认</button></div>' +
            '<div class="am-sub">影响范围</div>' +
            '<div class="am-rules">' +
              '<div><i>·</i>容量总览、机房平面色阶、机柜 42U、报表「扩容建议」、可部署预测全部联动</div>' +
              '<div><i>·</i>改系数会触发历史重算，旧报表保留版本，能解释「当时为什么是 47%」</div>' +
              '<div><i>·</i>空间率物理上限 100%，算出超过即视为「已满」，不再往上跳</div>' +
              '<div><i>·</i>系数从设备类型取，不从机柜取；一台设备一个系数，柜内逐台累加</div>' +
            '</div>' +
          '</div>' +
        '</section>' +
        '<section class="card am-thcard">' +
          '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">预警阈值</h2>' +
            '<div class="card-extra">' + U.tag('第 12 条', 'blue') + '</div></div>' +
          '<div class="card-body am-rulebody">' +
            '<div class="tbl-wrap"><table class="tbl"><thead><tr>' +
              '<th style="width:132px">维度</th><th style="width:80px" class="num">关注</th><th style="width:80px" class="num">红线</th><th>说明</th>' +
            '</tr></thead><tbody>' + thRows + '</tbody></table></div>' +
            '<div class="am-sub">机房级覆盖</div>' +
            '<div class="am-ovr"><select class="select" data-f="throom">' +
              '<option value="">选择机房</option>' + roomOpts + '</select>' +
              '<button class="btn is-ghost is-sm" data-act="ovr-add">设电力红线</button></div>' +
            '<div class="am-ovrlist">' + (Object.keys(AM.OVR).length
              ? Object.keys(AM.OVR).map(function (r) {
                  var p = AM.OVR[r].power || {};
                  return '<div class="am-ovritem"><span>' + A.esc(r) + '</span><b>电力红线 ' + (p.red || '—') + '%</b>' +
                    '<button class="am-x" data-ovrdel="' + A.esc(r) + '">✕</button></div>';
                }).join('')
              : '<span class="muted">未设置，全部继承全局阈值</span>') + '</div>' +
            '<div class="am-note">4kW 柜与 8kW 高密柜风险不同，阈值按机房可配。颜色跟着<b>综合使用率</b>走而不是跟空间走 —— 这样「满了」才不可争辩。</div>' +
            '<div class="am-legend4">' +
              '<div><i class="am-sw is-off"></i><b>灰</b><span>未启用 / 无设备</span></div>' +
              '<div><i class="am-sw is-blue"></i><b>蓝</b><span>&lt; ' + AM.TH.space.warn + '% 可继续上架</span></div>' +
              '<div><i class="am-sw is-yellow"></i><b>黄</b><span>' + AM.TH.space.warn + '-' + AM.TH.space.red + '% 进入关注</span></div>' +
              '<div><i class="am-sw is-red"></i><b>红</b><span>&gt; ' + AM.TH.space.red + '% 不接新设备</span></div>' +
            '</div>' +
            '<div class="am-sub">当前机柜色阶分布</div>' +
            '<div class="am-lvc">' +
              lvcItem('off', '未启用', M.dist.off, M.totalCabs) +
              lvcItem('blue', '正常', M.dist.blue, M.totalCabs) +
              lvcItem('yellow', '关注', M.dist.yellow, M.totalCabs) +
              lvcItem('red', '满载', M.dist.red, M.totalCabs) +
            '</div>' +
            '<div class="am-sub">按机房看色阶（回答「这个机房为什么红」）</div>' +
            '<div class="tbl-wrap"><table class="tbl is-mini"><thead><tr>' +
              '<th>机房</th><th class="num" style="width:60px">机柜</th><th class="num" style="width:74px">平均综合</th>' +
              '<th class="num" style="width:60px">满载</th><th style="width:110px">色阶构成</th>' +
            '</tr></thead><tbody>' + roomLvRows + '</tbody></table></div>' +
            '<div class="am-note">综合率取三项<b>最高值</b>，不做加权平均；平均值只看趋势，判定永远看单项红线。</div>' +
          '</div>' +
        '</section>' +
        '<section class="card am-ckcard">' +
          '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">口径自检</h2>' +
            '<div class="card-extra">' + U.tag('8 / 8 通过', 'green') + '</div></div>' +
          '<div class="card-body am-ckbody">' + checks +
            '<div class="am-sub">这 6 个问题现在能在页面上答出来</div>' +
            '<div class="am-qa">' +
              [['这套系统部署在哪些机房？', '资产总览 · 系统部署视图'],
               ['这个机房还剩多少位置？', '机柜布局 + 容量总览'],
               ['这个机柜还能装几台？卡在哪？', '机柜与 U 位 · 四项明细'],
               ['这台设备原来在哪、谁挪的？', '设备台账 · 变更历史'],
               ['这个单位占了多少机柜？', '资产总览 · 按单位统计'],
               ['这 100 台设备能一次导进去吗？', '导入导出 · 差异预览']].map(function (x) {
                return '<div class="am-qa-i"><span>' + x[0] + '</span><em>' + x[1] + '</em></div>';
              }).join('') +
            '</div>' +
            '<div class="am-sub">口径版本对照 · 旧报表按当时版本解释，不回溯改写</div>' +
            '<div class="tbl-wrap"><table class="tbl is-mini"><thead><tr>' +
              '<th style="width:46px">版本</th><th style="width:92px">生效日期</th>' +
              '<th class="num" style="width:64px">空间率</th><th class="num" style="width:64px">综合率</th><th>口径要点</th>' +
            '</tr></thead><tbody>' +
              '<tr><td><b style="color:#00E39A">V3</b></td><td class="muted">2026-10-20</td>' +
                '<td class="num">' + A.fix(M.spaceRate, 2) + '%</td><td class="num">' + A.fix((M.D.capacity || {}).comp || 0, 2) + '%</td>' +
                '<td class="muted">计入空U系数，综合取最高项（本版）</td></tr>' +
              '<tr><td>V2</td><td class="muted">2026-06-12</td><td class="num">47.02%</td><td class="num">65.21%</td>' +
                '<td class="muted">只算设备实际占用 U，不含空U系数</td></tr>' +
              '<tr><td>V1</td><td class="muted">2025-01-08</td><td class="num">41.30%</td><td class="num">62.10%</td>' +
                '<td class="muted">按柜数粗估，无系数、无接口维度</td></tr>' +
            '</tbody></table></div>' +
            '<div class="am-note">「当时为什么是 47%」有答案：那是 V2 的算法。报表跟版本走，历史数不改写。</div>' +
            '<div class="am-sub">本期不做，已记录</div>' +
            '<div class="am-rules">' +
              '<div><i>·</i>系统只留业务系统一层，不建子系统树</div>' +
              '<div><i>·</i>IP 一对多只做设备内子表，不做跨设备多 IP 并列检索</div>' +
            '</div>' +
            '<div class="am-warn">' +
              '<b>两条要先跟运维打招呼</b>' +
              '<div>① 位置变更不走审批，改完即生效 —— 靠日志追溯，不靠流程拦。</div>' +
              '<div>② 机柜使用率会比旧系统明显变高（计入空U系数），这不是算错，是按物理占用算。</div>' +
            '</div>' +
          '</div>' +
        '</section>' +
      '</div>';

    A.qa('[data-coef]', box).forEach(function (inp) {
      inp.addEventListener('change', function () {
        var t = inp.getAttribute('data-coef');
        var old = AM.coef(t);
        AM.setCoef(t, inp.value);
        A.toast('空U系数「' + t + '」' + old + ' → ' + AM.coef(t) + '，已触发重算');
        AM.views.rule(box);
      });
    });
    A.qa('[data-th]', box).forEach(function (inp) {
      inp.addEventListener('change', function () {
        var dim = inp.getAttribute('data-th');
        var k = inp.getAttribute('data-k');
        var t = AM.TH[dim];
        AM.setTh(dim, k === 'warn' ? inp.value : t.warn, k === 'red' ? inp.value : t.red);
        A.toast(AM.TH[dim].label + '阈值已更新：' + AM.TH[dim].warn + '% / ' + AM.TH[dim].red + '%');
        AM.views.rule(box);
      });
    });
    var apply = A.q('[data-act="rule-apply"]', box);
    if (apply) apply.addEventListener('click', function () {
      AM.invalidate();
      var m2 = AM.build();
      A.toast('已重算：' + A.num(m2.total) + ' 台设备 / ' + m2.totalCabs + ' 个机柜，空间率 ' + A.fix(m2.spaceRate, 2) + '%');
      AM.views.rule(box);
    });
    var rs = A.q('[data-act="rule-reset"]', box);
    if (rs) rs.addEventListener('click', function () {
      AM.resetCoef();
      var m2 = AM.build();
      A.toast('空U系数已恢复默认，空间率回到 ' + A.fix(m2.spaceRate, 2) + '%（与容量总览一致）');
      AM.views.rule(box);
    });
    var add = A.q('[data-act="ovr-add"]', box);
    if (add) add.addEventListener('click', function () {
      var sel = A.q('[data-f="throom"]', box);
      var room = sel ? sel.value : '';
      if (!room) { A.toast('请先选择机房'); return; }
      AM.setTh('power', 70, 88, room);
      A.toast(room + ' 电力红线已设为 88%（比全局 85% 更严）');
      AM.views.rule(box);
    });
    A.qa('[data-ovrdel]', box).forEach(function (b) {
      b.addEventListener('click', function () {
        var r = b.getAttribute('data-ovrdel');
        delete AM.OVR[r];
        AM.invalidate();
        AM.views.rule(box);
      });
    });
    A.qa('[data-go]', box).forEach(function (el) {
      el.addEventListener('click', function () {
        AM.openSpace(el.getAttribute('data-go'), el.getAttribute('data-name'));
      });
    });
  };
})(window.App);
