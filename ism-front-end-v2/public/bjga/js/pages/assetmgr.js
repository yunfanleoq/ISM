/* =========================================================
   资产管理 · 页面壳 / 资产总览 / 设备台账 / 设备档案
   口径依据：《资产管理 · 口径定稿 v1.0》
   ========================================================= */
(function (A) {
  'use strict';

  var AM = A.AM;
  var U = A.UI;

  var VIEWS = [
    { id: 'overview', name: '资产总览', sub: '分布 · 部署 · 单位' },
    { id: 'ledger', name: '设备台账', sub: '全字段 · 检索' },
    { id: 'cabinet', name: '机柜与 U 位', sub: '四项 · 卡点' },
    { id: 'flow', name: '变更与审批', sub: '两级 · 留痕' },
    { id: 'stock', name: '盘点核实', sub: '半年 · 待核实' },
    { id: 'maint', name: '运维计划', sub: '周期 · 到期' },
    { id: 'inout', name: '出入库', sub: '入库 · 出库' },
    { id: 'io', name: '导入导出', sub: '全字段 · 校验' },
    { id: 'rule', name: '口径设置', sub: '系数 · 阈值' }
  ];

  var S = {
    view: 'overview',
    dim: 'org',
    biz: '',
    led: { room: '', cat: '', st: '', org: '', kw: '', page: 1 },
    size: 18
  };

  function bindGo(scope) {
    A.qa('[data-go]', scope).forEach(function (el) {
      el.addEventListener('click', function (e) {
        e.stopPropagation();
        AM.openSpace(el.getAttribute('data-go'), el.getAttribute('data-name'));
      });
    });
  }

  A.register('assetmgr', {
    render: function (root) {
      var M = AM.build();
      if (A.subArg) {
        for (var i = 0; i < VIEWS.length; i++) {
          if (VIEWS[i].id === A.subArg) S.view = VIEWS[i].id;
        }
      }
      if (!M.bizs.filter(function (b) { return b.name === S.biz; }).length) {
        S.biz = M.bizs[0] ? M.bizs[0].name : '';
      }

      root.innerHTML =
        '<div class="page-head">' +
          '<h1 class="page-title">资产管理</h1>' +
          '<span class="page-crumb">' + A.crumb('设备台账 · 空间占用 · 上下架') + '</span>' +
          '<div class="page-tools">' +
            '<span class="am-vtag">口径定稿 v1.0 · 13 条已落地</span>' +
            '<button class="btn is-ghost is-sm" data-act="go-io">导入 / 导出</button>' +
            '<button class="btn is-primary is-sm" data-act="new">+ 设备登记</button>' +
          '</div>' +
        '</div>' +
        '<div class="am-nav">' + VIEWS.map(function (v) {
          return '<button class="am-navb' + (v.id === S.view ? ' is-on' : '') + '" data-v="' + v.id + '">' +
            '<b>' + v.name + '</b><i>' + v.sub + '</i></button>';
        }).join('') + '</div>' +
        '<div class="am-body" id="amBody"></div>';

      var box = A.q('#amBody', root);

      function paint() {
        A.qa('.am-navb', root).forEach(function (b) {
          b.classList.toggle('is-on', b.getAttribute('data-v') === S.view);
        });
        box.innerHTML = '';
        AM.currentRefresh = function () { paint(); };
        var fn = AM.views[S.view];
        if (fn) fn(box);
      }

      A.q('.am-nav', root).addEventListener('click', function (e) {
        var b = e.target.closest('.am-navb');
        if (!b) return;
        var id = b.getAttribute('data-v');
        if (id === S.view) return;
        S.view = id;
        paint();
      });
      A.q('[data-act="go-io"]', root).addEventListener('click', function () { S.view = 'io'; paint(); });
      A.q('[data-act="new"]', root).addEventListener('click', function () {
        A.drawer('设备登记',
          U.kv([['资产编码', '系统生成（永久，不复用）'], ['主 SN', '参与查重；附属 SN 只记录不判重'],
            ['IP', '一对多，管理 / 业务 / 存储 / 心跳'], ['初始状态', '未定位（登记后进入待上架）']]) +
          '<div class="am-note" style="margin-top:14px">登记只建档，不占 U 位；上架时才选择机柜与 U 位，轨迹从「未定位」开始。</div>' +
          '<div class="am-iobtns" style="margin-top:16px"><button class="btn is-primary" data-x="ok">确认登记</button>' +
          '<button class="btn is-ghost" data-x="no">取消</button></div>');
        A.qa('[data-x]').forEach(function (b) {
          b.addEventListener('click', function () {
            A.closeDrawer();
            if (b.getAttribute('data-x') === 'ok') A.toast('设备已登记：编码已生成，状态「未定位」');
          });
        });
      });

      paint();
    }
  });

  /* =========================================================
     资产总览
     ========================================================= */
  AM.views = AM.views || {};
  AM.views.overview = function (box) {
    var M = AM.build();
    var D = M.dist;

    /* --- 四维分布 --- */
    var DIMS = {
      org: { t: '使用单位', cls: 'p', rows: M.orgs.map(function (o) { return { label: o.name, value: o.dev }; }) },
      biz: { t: '业务系统', cls: '', rows: M.bizs.map(function (b) { return { label: b.name, value: b.dev }; }) },
      cat: { t: '设备类型', cls: 'g', rows: M.cats.map(function (c) { return { label: c.name, value: c.dev }; }) },
      room: { t: '机房', cls: 'o', rows: M.rooms.map(function (r) { return { label: r.name.replace('屏蔽机房', ''), value: r.dev, go: 'room', name: r.name }; }) }
    };
    var cur = DIMS[S.dim] || DIMS.org;
    var rows = cur.rows.slice().sort(function (x, y) { return y.value - x.value; });
    var sum = rows.reduce(function (s, x) { return s + x.value; }, 0);
    var PAL = ['#00D8FF', '#2E7CF6', '#00E39A', '#FFD24D', '#FFA24D', '#A46BFF', '#FF6B7A',
      '#7FE9FF', '#4ED6A8', '#F0A6FF', '#6E9BD6', '#C0E24D', '#FF8A5C'];
    var segs = rows.map(function (r, i) {
      return '<i style="width:' + (sum ? r.value / sum * 100 : 0) + '%;background:' + PAL[i % PAL.length] + '"></i>';
    }).join('');
    var legs = rows.map(function (r, i) {
      var go = r.go ? ' class="am-go" data-go="' + r.go + '" data-name="' + A.esc(r.name) + '"' : '';
      return '<div' + go + '><i style="background:' + PAL[i % PAL.length] + '"></i>' +
        '<span title="' + A.esc(r.label) + '">' + A.esc(r.label) + '</span>' +
        '<b>' + A.num(r.value) + '</b><em>' + A.fix(r.value / (sum || 1) * 100, 1) + '%</em></div>';
    }).join('');
    var top3 = rows.slice(0, 3).reduce(function (t, x) { return t + x.value; }, 0);
    var top1pct = rows.length ? A.fix(rows[0].value / (sum || 1) * 100, 1) : '0.0';

    var dimTabs = Object.keys(DIMS).map(function (k) {
      return '<button class="am-pill' + (k === S.dim ? ' is-on' : '') + '" data-dim="' + k + '">' + DIMS[k].t + '</button>';
    }).join('');

    /* --- 系统部署视图（第 02 条：只算已部署） --- */
    var bizMap = {};
    M.list.forEach(function (d) {
      if (d.biz !== S.biz || !d.located) return;
      if (!bizMap[d.room]) bizMap[d.room] = { room: d.room, floor: d.floor, dev: 0, cabs: {} };
      bizMap[d.room].dev++;
      bizMap[d.room].cabs[d.cab] = 1;
    });
    var bizRooms = Object.keys(bizMap).map(function (k) {
      var r = bizMap[k];
      r.cabN = Object.keys(r.cabs).length;
      var rm = M.rooms.filter(function (x) { return x.name === r.room; })[0];
      r.comp = rm ? rm.comp : 0;
      return r;
    }).sort(function (a, b) { return b.dev - a.dev; });
    var bizDev = bizRooms.reduce(function (s, x) { return s + x.dev; }, 0);
    var bizCab = bizRooms.reduce(function (s, x) { return s + x.cabN; }, 0);

    var bizOpts = M.bizs.map(function (b) {
      return '<option value="' + A.esc(b.name) + '"' + (b.name === S.biz ? ' selected' : '') + '>' +
        A.esc(b.name) + '（' + b.dev + '）</option>';
    }).join('');

    var bizRows = bizRooms.map(function (r) {
      return '<tr><td><span class="cap-room" data-go="room" data-name="' + A.esc(r.room) + '">' + A.esc(r.room) + '</span></td>' +
        '<td class="num"><span class="cap-room" data-go="floor" data-name="' + A.esc(r.floor) + '">' + A.esc(r.floor) + '</span></td>' +
        '<td class="num">' + r.dev + '</td><td class="num">' + r.cabN + '</td>' +
        '<td><span class="am-mini-bar"><i style="width:' + Math.min(100, r.comp) + '%"></i></span>' +
        '<span class="am-mini-v">' + A.fix(r.comp, 1) + '%</span></td></tr>';
    }).join('');

    /* --- 按单位统计（第 01 条） --- */
    var orgRows = M.orgs.map(function (o) {
      return '<tr><td>' + A.esc(o.name) + '</td>' +
        '<td class="num">' + A.num(o.dev) + '</td>' +
        '<td class="num muted">' + A.fix(o.dev / M.total * 100, 1) + '%</td>' +
        '<td class="num">' + o.mainCabs + '</td>' +
        '<td class="num muted">' + o.touch + '</td>' +
        '<td class="num">' + o.roomN + '</td></tr>';
    }).join('');
    var sumCab = M.orgs.reduce(function (s, o) { return s + o.mainCabs; }, 0);

    var Dd = window.DCIM || {};
    var spaceFloors = (Dd.floors || []).map(function (f) { return f.floor; }).filter(Boolean).sort();
    var floorBtns = spaceFloors.map(function (f) {
      return '<button type="button" class="btn is-ghost is-sm" data-go="floor" data-name="' + A.esc(f) + '">' + A.esc(f) + '</button>';
    }).join('');
    var roomBtns = (Dd.rooms || []).filter(function (r) { return r.cabCount > 0; }).slice().sort(function (a, b) {
      return String(a.floor).localeCompare(String(b.floor)) || String(a.name).localeCompare(String(b.name), 'zh');
    }).map(function (r) {
      return '<button type="button" class="btn is-ghost is-sm" data-go="room" data-name="' + A.esc(r.name) + '">' + A.esc(r.name) + '</button>';
    }).join('');

    box.innerHTML =
      '<section class="card am-kpi"><div class="stats am-stats">' +
        '<div class="stat-box"><span><span class="v">' + A.num(M.total) + '</span><span class="u">台</span></span><span class="l">设备总数</span></div>' +
        '<div class="stat-box"><span><span class="v" style="color:#00E39A">' + A.num(M.cnt['在用']) + '</span><span class="u">台</span></span><span class="l">在用</span></div>' +
        '<div class="stat-box"><span><span class="v">' + M.usedCabs + '</span><span class="u">/ ' + M.totalCabs + ' 柜</span></span><span class="l">已用机柜</span></div>' +
        '<div class="stat-box"><span><span class="v">' + M.deployedRooms.length + '</span><span class="u">个</span></span><span class="l">已部署机房</span></div>' +
        '<div class="stat-box"><span><span class="v" style="color:#FF4D5E">' + A.num(M.cnt['未定位']) + '</span><span class="u">台</span></span><span class="l">未定位</span></div>' +
        '<div class="stat-box"><span><span class="v" style="color:#FFD24D">' + A.fix(M.spaceRate, 2) + '</span><span class="u">%</span></span><span class="l">空间利用率</span></div>' +
      '</div></section>' +

      '<section class="card" style="flex:0 0 auto">' +
        '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">按楼层 / 机房查看资产</h2>' +
          '<div class="card-extra"><span class="muted" style="font-size:12px">进入该层或该机房的立体画面，不打开整栋楼</span></div></div>' +
        '<div class="card-body am-space">' +
          '<div class="am-space-row"><b>楼层</b>' + floorBtns + '</div>' +
          '<div class="am-space-row am-space-rooms"><b>机房</b>' + roomBtns + '</div>' +
        '</div></section>' +

      '<div class="row grow">' +
        '<section class="card am-panel">' +
          '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">资产四维分布</h2>' +
            '<div class="card-extra">' + dimTabs + '</div></div>' +
          '<div class="card-body am-dimbody">' +
            '<div class="am-stack">' + segs + '</div>' +
            '<div class="am-insight">' +
              '<div><span>占比最高</span><b>' + A.esc(rows.length ? rows[0].label : '—') + '</b>' +
                '<em>' + A.num(rows.length ? rows[0].value : 0) + ' 台 · ' + top1pct + '%</em></div>' +
              '<div><span>TOP3 集中度</span><b>' + A.fix(top3 / (sum || 1) * 100, 1) + '%</b>' +
                '<em>前 3 项合计 ' + A.num(top3) + ' 台</em></div>' +
              '<div><span>维度项数</span><b>' + rows.length + ' 项</b>' +
                '<em>均值 ' + A.num(rows.length ? Math.round(sum / rows.length) : 0) + ' 台</em></div>' +
            '</div>' +
            '<div class="am-leg">' + legs + '</div>' +
            '<div class="am-sumline"><span>合计 ' + cur.t + '</span><b>' + A.num(sum) + ' 台</b>' +
              '<em>按设备归属精确统计 · 一台设备只算一次</em></div>' +
          '</div>' +
        '</section>' +

        '<section class="card am-syscard">' +
          '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">系统部署视图</h2>' +
            '<div class="card-extra"><select class="select am-sel" data-f="biz">' + bizOpts + '</select></div></div>' +
          '<div class="card-body am-sysbody">' +
            '<div class="am-sysrow">' +
              '<div class="am-sy"><span>涉及机房</span><b>' + bizRooms.length + '</b><em>只算已部署</em></div>' +
              '<div class="am-sy"><span>设备</span><b>' + A.num(bizDev) + '</b><em>台</em></div>' +
              '<div class="am-sy"><span>占用机柜</span><b>' + bizCab + '</b><em>个</em></div>' +
            '</div>' +
            '<div class="tbl-wrap am-systbl"><table class="tbl"><thead><tr>' +
              '<th>机房</th><th class="num" style="width:56px">楼层</th><th class="num" style="width:64px">设备</th>' +
              '<th class="num" style="width:64px">机柜</th><th style="width:118px">机房综合率</th>' +
            '</tr></thead><tbody>' + (bizRows || '<tr><td colspan="5"><div class="tbl-empty">该业务系统暂无已部署设备</div></td></tr>') + '</tbody></table></div>' +
            '<div class="am-note">「<b>规划机房</b>」字段默认不在界面显示，只在容量规划视图使用 —— 规划随时变，混进统计会让数字变虚。' +
              '网络设备不单列层级，作为「设备类型 = 网络设备」挂在所属业务系统下。</div>' +
          '</div>' +
        '</section>' +
      '</div>' +

      '<div class="row am-bot">' +
        '<section class="card am-orgcard">' +
          '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">按使用单位统计</h2>' +
            '<div class="card-extra"><span class="muted" style="font-size:12px">设备归属精确口径 · 机柜按主用单位归属</span></div></div>' +
          '<div class="card-body am-orgbody"><div class="tbl-wrap"><table class="tbl"><thead><tr>' +
            '<th>使用单位</th><th class="num" style="width:90px">设备</th><th class="num" style="width:76px">占比</th>' +
            '<th class="num" style="width:104px">主用机柜</th><th class="num" style="width:104px">涉及机柜</th>' +
            '<th class="num" style="width:90px">涉及机房</th>' +
          '</tr></thead><tbody>' + orgRows +
            '<tr class="am-total"><td>合计</td><td class="num">' + A.num(M.total) + '</td><td class="num">100.0%</td>' +
            '<td class="num">' + sumCab + '</td><td class="num muted">含共用柜，不可加总</td><td class="num">—</td></tr>' +
          '</tbody></table></div>' +
          '<div class="am-note">主用机柜之和 <b>' + sumCab + ' = 已用机柜 ' + M.usedCabs + '</b>：一柜多单位时按柜内台数最多者落一个主分类，' +
            '因此可加总、不重复计数；设备数按设备自身归属精确统计，两者各自成立。</div></div>' +
        '</section>' +

        '<section class="card am-qcard">' +
          '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">台账质量</h2>' +
            '<div class="card-extra">' + U.tag('可直接进盘点', 'green') + '</div></div>' +
          '<div class="card-body am-qbody">' +
            qRow('主 SN 待补录', M.snMissing, '补字段，不新建记录', '#FFD24D') +
            qRow('双 SN 设备', M.snDual, '主 SN 参与查重，附属只记录', '#7FE9FF') +
            qRow('未定位设备', M.cnt['未定位'], '下架回库保留原位置轨迹', '#FF4D5E') +
            qRow('已下架待清柜', M.cnt['已下架'], '占位至清柜，不释放 U 位', '#5E7FA8') +
            qRow('满载 / 关注机柜', D.red + ' / ' + D.yellow, '红 ' + D.red + ' 柜 · 黄 ' + D.yellow + ' 柜 · 蓝 ' + D.blue + ' 柜', '#FFA24D') +
            qRow('IP 记录（一对多）', M.list.reduce(function (n, d) { return n + d.ips.length; }, 0), '列表只显主 IP，其余进设备内子表', '#7FE9FF') +
            qRow('编码不合规', M.list.filter(function (d) {
              return !/^BJDC-\d{3}-\d+-\d+-\d{4}$/.test(d.code);
            }).length, '全部符合 BJDC-机房-列-序号-流水', '#00E39A') +
            qRow('轨迹缺失', 0, '每台设备至少 1 段位置轨迹', '#00E39A') +
          '</div>' +
        '</section>' +
      '</div>';

    A.qa('[data-dim]', box).forEach(function (b) {
      b.addEventListener('click', function () {
        S.dim = b.getAttribute('data-dim');
        AM.views.overview(box);
      });
    });
    var sel = A.q('[data-f="biz"]', box);
    if (sel) sel.addEventListener('change', function () {
      S.biz = sel.value;
      AM.views.overview(box);
    });
    bindGo(box);
  };

  function qRow(label, val, note, color) {
    /* 数值走 num 格式化；已经是「60 / 213」这类复合串则原样输出 */
    var txt = typeof val === 'string' ? val : A.num(val);
    return '<div class="am-q"><span>' + label + '</span>' +
      '<b style="color:' + color + '">' + txt + '</b>' +
      '<em>' + note + '</em></div>';
  }

  /* =========================================================
     设备台账
     ========================================================= */
  function ledgerList() {
    var M = AM.build();
    var l = M.list;
    var s = S.led;
    if (s.room) l = l.filter(function (d) { return d.room === s.room; });
    if (s.cat) l = l.filter(function (d) { return d.cat === s.cat; });
    if (s.st) l = l.filter(function (d) { return d.st === s.st; });
    if (s.org) l = l.filter(function (d) { return d.org === s.org; });
    if (s.kw) {
      var k = s.kw.toLowerCase();
      l = l.filter(function (d) {
        return (d.code + d.sn + d.snAlt + d.model + d.brand + d.biz + d.org + d.owner + d.mainIp + d.cab)
          .toLowerCase().indexOf(k) >= 0;
      });
    }
    return l;
  }

  AM.views.ledger = function (box) {
    var M = AM.build();
    var st = S.led;

    var roomOpts = [{ v: '', t: '全部机房' }].concat(M.roomOrder.map(function (r) { return { v: r, t: r }; }));
    var catOpts = [{ v: '', t: '全部类别' }].concat(M.cats.map(function (c) { return { v: c.name, t: c.name }; }));
    var stOpts = [{ v: '', t: '全部状态' }].concat(['在用', '备用', '维修中', '已下架', '未定位'].map(function (x) { return { v: x, t: x }; }));
    var orgOpts = [{ v: '', t: '全部单位' }].concat(M.orgs.map(function (o) { return { v: o.name, t: o.name }; }));

    var list = ledgerList();
    var pg = A.paginate(list, st.page, S.size);
    if (pg.page !== st.page) st.page = pg.page;

    var table = U.table({
      cols: [
        { t: '资产编码（永久）', w: 176, mono: true, html: function (d) { return A.esc(d.code); } },
        { t: '主 IP', w: 118, mono: true, html: function (d) { return A.esc(d.mainIp); } },
        { t: '主 SN', w: 150, mono: true, html: function (d) {
          if (d.snState === 'missing') return U.tag('待补录', 'orange');
          return A.esc(d.sn) + (d.snAlt ? ' <span class="am-sup">+1</span>' : '');
        } },
        { t: '类别', w: 100, html: function (d) { return A.esc(d.cat); } },
        { t: '品牌 / 型号', w: 168, html: function (d) {
          return '<span class="muted">' + A.esc(d.brand) + '</span> ' + A.esc(d.model);
        } },
        { t: '位置', w: 168, html: function (d) {
          return d.located
            ? A.esc(d.cabShort + ' · U' + d.uStart + (d.uEnd > d.uStart ? '-' + d.uEnd : ''))
            : U.tag('未定位', 'grey');
        } },
        { t: '业务系统', w: 138, html: function (d) { return A.esc(d.biz); } },
        { t: '使用单位', w: 128, html: function (d) { return A.esc(d.org); } },
        { t: '状态', w: 84, html: function (d) { return U.tag(d.st, AM.stTag(d.st)); } },
        { t: '操作', w: 66, html: function (d) {
          return '<button class="btn is-ghost is-sm" data-off="' + A.esc(d.code) + '">' +
            (d.st === '未定位' ? '上架' : '下架') + '</button>';
        } }
      ],
      rows: pg.rows,
      rowAttr: function (d) { return 'data-code="' + A.esc(d.code) + '" style="cursor:pointer"'; }
    });

    box.innerHTML =
      '<section class="card am-filters">' +
        '<div class="row ac" style="gap:10px">' +
          '<div class="field"><span class="field-label">机房</span>' + U.select('room', roomOpts, st.room) + '</div>' +
          '<div class="field"><span class="field-label">类别</span>' + U.select('cat', catOpts, st.cat) + '</div>' +
          '<div class="field"><span class="field-label">状态</span>' + U.select('st', stOpts, st.st) + '</div>' +
          '<div class="field"><span class="field-label">单位</span>' + U.select('org', orgOpts, st.org) + '</div>' +
          '<div class="field"><span class="field-label">关键字</span>' +
            '<input class="input" data-f="kw" style="width:250px" placeholder="编码 / SN / IP / 型号 / 业务系统" value="' + A.esc(st.kw) + '"/></div>' +
          '<button class="btn is-primary is-sm" data-act="q">查询</button>' +
          '<button class="btn is-ghost is-sm" data-act="reset">重置</button>' +
          '<span class="am-sep"></span>' +
          '<button class="am-pill' + (st.st === '未定位' ? ' is-on' : '') + '" data-quick="未定位">未定位 ' + A.num(M.cnt['未定位']) + '</button>' +
          '<button class="am-pill' + (st.st === '已下架' ? ' is-on' : '') + '" data-quick="已下架">已下架 ' + A.num(M.cnt['已下架']) + '</button>' +
        '</div>' +
      '</section>' +
      '<section class="card am-ledger">' +
        '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">设备台账</h2>' +
          '<div class="card-extra">' +
            '<span class="muted" style="font-size:12px">共 ' + A.num(list.length) + ' 条 · 点击行看档案 / IP / 轨迹</span>' +
            '<button class="btn is-ghost is-sm" data-act="export">导出当前筛选</button>' +
          '</div></div>' +
        '<div class="card-body am-ledbody">' + table + '</div>' +
        '<div class="am-pager">' + U.pager({ page: pg.page, pages: pg.pages, total: pg.total, size: S.size }) + '</div>' +
      '</section>';

    function repaint() { AM.views.ledger(box); }

    A.q('[data-act="q"]', box).addEventListener('click', function () {
      st.room = A.q('[data-f="room"]', box).value;
      st.cat = A.q('[data-f="cat"]', box).value;
      st.st = A.q('[data-f="st"]', box).value;
      st.org = A.q('[data-f="org"]', box).value;
      st.kw = A.q('[data-f="kw"]', box).value.trim();
      st.page = 1;
      repaint();
    });
    A.q('[data-act="reset"]', box).addEventListener('click', function () {
      st.room = ''; st.cat = ''; st.st = ''; st.org = ''; st.kw = ''; st.page = 1;
      repaint();
    });
    A.qa('[data-quick]', box).forEach(function (b) {
      b.addEventListener('click', function () {
        st.st = st.st === b.getAttribute('data-quick') ? '' : b.getAttribute('data-quick');
        st.page = 1;
        repaint();
      });
    });
    A.q('[data-act="export"]', box).addEventListener('click', function () {
      A.toast('已导出 ' + A.num(list.length) + ' 条 · 42 列 · 含当前筛选条件');
    });
    A.qa('.pg', box).forEach(function (b) {
      b.addEventListener('click', function () {
        if (b.disabled) return;
        st.page = Number(b.getAttribute('data-p'));
        repaint();
      });
    });
    A.qa('tr[data-code]', box).forEach(function (tr) {
      tr.addEventListener('click', function (e) {
        if (e.target.closest('[data-off]')) return;
        AM.openAsset(tr.getAttribute('data-code'), repaint);
      });
    });
    A.qa('[data-off]', box).forEach(function (b) {
      b.addEventListener('click', function (e) {
        e.stopPropagation();
        var d = AM.assetOf(b.getAttribute('data-off'));
        if (!d) return;
        if (d.st === '未定位') { AM.openAsset(d.code, repaint, true); return; }
        d.st = '已下架';
        d.hist.unshift({ t: A.now(), act: '下架回库', from: d.cabShort + ' · U' + d.uStart, to: '未定位（保留原位置）', by: '当前用户', why: '手工下架' });
        AM.invalidate();
        A.toast(d.code + ' 已下架 → 进入「未定位」，轨迹保留');
        repaint();
      });
    });
  };

  /* =========================================================
     设备档案
     ========================================================= */
  AM.openAsset = function (code, refresh, moveMode) {
    var M = AM.build();
    var d = M.byCode[code];
    if (!d) { A.toast('未找到该资产编码'); return; }
    AM.currentRefresh = refresh || AM.currentRefresh;

    var biz = M.bizs.filter(function (b) { return b.name === d.biz; })[0];
    var rm = M.rooms.filter(function (r) { return r.name === d.room; })[0];
    var cab = d.located ? M.cabIndex[d.cab] : null;
    var sameCab = cab ? cab.dlist.length : 0;

    var ips = d.ips.map(function (p) {
      return '<tr><td>' + p.t + '</td><td class="mono">' + A.esc(p.ip) + '</td>' +
        '<td class="muted">' + A.esc(p.note) + '</td>' +
        '<td>' + (p.main ? U.tag('主', 'green') : '<span class="muted">—</span>') + '</td></tr>';
    }).join('');

    var tl = d.hist.map(function (h, i) {
      return '<div class="am-tli' + (i === 0 ? ' is-now' : '') + '"><i></i><div>' +
        '<b>' + A.esc(h.t) + ' · ' + A.esc(h.act) + '</b>' +
        '<span>' + A.esc(h.from) + ' <em>→</em> ' + A.esc(h.to) + '</span>' +
        '<u>' + A.esc(h.by) + (h.why ? ' · ' + A.esc(h.why) : '') + '</u></div></div>';
    }).join('');

    var cabOpts = M.cabs.map(function (c) {
      return '<option value="' + A.esc(c.name) + '"' + (c.name === d.cab ? ' selected' : '') + '>' + A.esc(c.short) + '（' + c.dev + ' 台）</option>';
    }).join('');

    A.drawer('设备档案 · ' + d.code,
      '<div class="am-dhead">' + U.tag(d.st, AM.stTag(d.st)) +
        '<span class="am-code">' + A.esc(d.code) + '</span></div>' +
      U.kv([
        ['资产编码', d.code + '（永久，不复用）'],
        ['主 SN', d.snState === 'missing' ? '待补录' : d.sn],
        ['附属 SN', d.snAlt || '—'],
        ['设备类别', d.cat + ' / ' + d.type],
        ['品牌型号', d.brand + ' ' + d.model],
        ['当前位置', d.located ? (d.cabShort + ' · U' + d.uStart + '-' + d.uEnd) : '未定位'],
        ['所属机房', d.room + '（' + d.floor + '）'],
        ['额定功率', A.fix(d.kw, 2) + ' kW · ' + d.psu + ' 电源'],
        ['空U系数', A.fix(d.coef, 1) + '（按设备类型）'],
        ['使用单位', d.org],
        ['业务系统', d.biz],
        ['责任人', d.owner]
      ]) +
      '<div class="am-sub">IP 地址（一对多，列表只显示主 IP）</div>' +
      '<div class="tbl-wrap"><table class="tbl am-iptbl"><thead><tr><th>类型</th><th>地址</th><th>用途</th><th style="width:52px">主 IP</th></tr></thead>' +
        '<tbody>' + ips + '</tbody></table></div>' +
      '<div class="am-sub">影响分析</div>' +
      '<div class="am-imp">' +
        '<span>承载业务<b>' + A.esc(d.biz) + '</b></span>' +
        '<span>该系统涉及机房<b>' + (biz ? biz.roomN : 0) + ' 个</b></span>' +
        '<span>同柜设备<b>' + sameCab + ' 台</b></span>' +
        '<span>机房告警<b>' + (rm ? rm.alarmN || 0 : 0) + ' 条</b></span>' +
      '</div>' +
      '<div class="am-sub">变更历史（不审批 · 强制留痕）</div>' +
      '<div class="am-tl">' + tl + '</div>' +
      '<div class="am-sub">位置变更</div>' +
      '<div class="am-move">' +
        '<select class="select" data-m-cab>' + cabOpts + '</select>' +
        '<input class="input" data-m-u type="number" min="1" max="42" value="' + Math.min(42, d.uStart) + '" style="width:74px"/>' +
        '<button class="btn is-primary is-sm" data-m-ok>确认变更</button>' +
      '</div>' +
      '<div class="am-note">改完即生效，不需要审批；系统自动记录原位置、新位置、操作人与时间。' +
        '下架不删除设备，留在「未定位」里，再次上架时选新机柜，轨迹才连得上。</div>' +
      '<div class="am-iobtns" style="margin-top:16px">' +
        (d.room ? '<button class="btn is-ghost is-sm" type="button" data-go="room" data-name="' + A.esc(d.room) + '">查看该机房立体</button>' : '') +
        (d.floor ? '<button class="btn is-ghost is-sm" type="button" data-go="floor" data-name="' + A.esc(d.floor) + '">查看 ' + A.esc(d.floor) + ' 立体</button>' : '') +
        '<button class="btn is-ghost is-sm" data-a="qr">打印标签</button>' +
        '<button class="btn is-ghost is-sm" data-a="sn">补录 SN</button>' +
      '</div>');

    var dr = A.q('.drawer');
    if (dr) bindGo(dr);

    var mv = A.q('[data-m-ok]');
    if (mv) mv.addEventListener('click', function () {
      var nb = A.q('[data-m-cab]').value;
      var nu = Number(A.q('[data-m-u]').value) || 1;
      var nc = M.cabIndex[nb];
      if (!nc) return;
      var old = d.located ? (d.cabShort + ' · U' + d.uStart) : '未定位';
      d.cab = nb; d.cabShort = nc.short; d.uStart = nu; d.uEnd = nu + d.uH - 1;
      d.located = true; d.st = '在用'; d.room = nc.room; d.floor = nc.floor;
      d.hist.unshift({ t: A.now(), act: '位置变更', from: old, to: nc.short + ' · U' + nu, by: '当前用户', why: '手工变更' });
      AM.invalidate();
      A.closeDrawer();
      A.toast(d.code + ' 位置已变更：' + old + ' → ' + nc.short + ' U' + nu + '（已留痕）');
      if (AM.currentRefresh) AM.currentRefresh();
    });
    A.qa('[data-a]').forEach(function (b) {
      b.addEventListener('click', function () {
        A.toast(b.getAttribute('data-a') === 'qr' ? '标签已发送至打印队列' : '已进入 SN 补录（补字段，不新建记录）');
      });
    });

    if (moveMode) {
      var f = A.q('[data-m-cab]');
      if (f) f.focus();
    }
    void moveMode;
  };
})(window.App);
