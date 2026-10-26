/* =========================================================
   设施监控：空间总览 / 温湿度 / UPS / 蓄电池 / 列头柜 / 配电柜 /
             空调 / 漏水 / 氢气 / 防雷 / 告警
   ========================================================= */
(function (A) {
  'use strict';

  var SUBS = [
    { id: 'space',  name: '空间总览' },
    { id: 'th',     name: '温湿度' },
    { id: 'ups',    name: 'UPS' },
    { id: 'batt',   name: '蓄电池' },
    { id: 'cab',    name: '列头柜' },
    { id: 'pdu',    name: '配电柜' },
    { id: 'ac',     name: '空调' },
    { id: 'leak',   name: '漏水' },
    { id: 'h2',     name: '氢气' },
    { id: 'lp',     name: '防雷' },
    { id: 'alarm',  name: '告警' },
    { id: 'video',  name: '视频' },
    { id: 'door',   name: '门禁' },
    { id: 'comm',   name: '通信状态' },
    { id: 'fire',   name: '消防' }
  ];

  var state = { sub: 'space', room: null, floor: null, page: 1, level: '', cat: '', kw: '', dev: null };

  function D() { return window.DCIM; }

  function alarmSource() {
    var feed = A.alarmFeed && A.alarmFeed();
    if (feed) return feed.list;
    return (D() && D().alarms) || [];
  }

  function alarmStatSource() {
    var feed = A.alarmFeed && A.alarmFeed();
    if (feed) return feed.stats;
    return D().alarmStats;
  }

  function roomOf(name) {
    var r = D().rooms.filter(function (x) { return x.name === name; });
    return r[0] || D().rooms[0];
  }

  /* ---------------- 左侧树 ---------------- */
  function treeHtml() {
    var byFloor = {};
    D().rooms.forEach(function (r) { (byFloor[r.floor] = byFloor[r.floor] || []).push(r); });
    var floors = Object.keys(byFloor).sort();
    var cur = state.room || (D().rooms[0] && D().rooms[0].name);

    return '<div class="side-title">空间层级</div>' +
      '<div class="side-scroll"><div class="tree">' +
      floors.map(function (f) {
        var open = !state.floor || state.floor === f;
        return '<div class="tree-group">' +
          '<div class="tree-node' + (open ? ' is-open' : '') + '" data-floor="' + A.esc(f) + '">' +
            '<span class="caret">▶</span><span class="nm">' + A.esc(f) + '</span>' +
            '<span class="ct">' + byFloor[f].length + '</span></div>' +
          '<div class="tree-children"' + (open ? '' : ' style="display:none"') + '>' +
          byFloor[f].map(function (r) {
            return '<div class="tree-node' + (r.name === cur ? ' is-active' : '') + '" data-room="' + A.esc(r.name) + '">' +
              '<span class="dot" style="background:' + (r.rate > 85 ? 'var(--orange)' : 'var(--cyan)') + '"></span>' +
              '<span class="nm" title="' + A.esc(r.name) + '">' + A.esc(r.name) + '</span>' +
              '<span class="ct">' + r.dev + '</span></div>';
          }).join('') +
          '</div></div>';
      }).join('') +
      '</div></div>';
  }

  /* ---------------- 各子视图 ---------------- */
  function viewSpace() {
    var rooms = D().rooms;
    var byFloor = {};
    rooms.forEach(function (r) { (byFloor[r.floor] = byFloor[r.floor] || []).push(r); });
    var floors = Object.keys(byFloor).sort();

    var totalCab = rooms.reduce(function (a, r) { return a + r.cabCount; }, 0);
    var totalDev = rooms.reduce(function (a, r) { return a + r.dev; }, 0);
    var totalKw = rooms.reduce(function (a, r) { return a + r.kw; }, 0);
    var avgRate = rooms.reduce(function (a, r) { return a + r.rate; }, 0) / rooms.length;

    var left =
      '<div class="card" style="flex:1 1 0">' +
        '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">楼层平面 · 点击机房下钻</h2>' +
          '<div class="card-extra"><span class="muted" style="font-size:12px">共 ' + rooms.length + ' 个机房</span>' +
            '<button class="btn is-ghost" type="button" id="btnTour2" style="margin-left:10px;font-size:12px;padding:5px 10px">空间巡检 ›</button></div></div>' +
        '<div class="card-body plan"><div class="plan-inner">' +
          floors.map(function (f) {
            return '<div><div style="font-size:12px;color:var(--t-muted);margin-bottom:6px">' + A.esc(f) + ' 层</div>' +
              '<div style="display:grid;grid-template-columns:repeat(auto-fill,minmax(196px,1fr));gap:10px">' +
              byFloor[f].map(function (r) {
                return '<div class="zone' + (r.name === (state.room || rooms[0].name) ? ' is-active' : '') + '" data-room="' + A.esc(r.name) + '">' +
                  '<div class="zone-nm">' + A.esc(r.name) +
                    (r.rate > 85 ? '<span class="tag is-orange">高负载</span>' : '') + '</div>' +
                  '<div class="zone-meta"><span>机柜 <b>' + r.cabCount + '</b></span><span>设备 <b>' + r.dev + '</b></span>' +
                    '<span>功率 <b>' + A.num(r.kw) + '</b> kW</span></div>' +
                  '<div class="zone-bar"><i style="width:' + A.clamp(r.rate, 0, 100).toFixed(0) + '%"></i></div>' +
                  '</div>';
              }).join('') + '</div></div>';
          }).join('') +
        '</div>' +
        '<div class="plan-legend"><span><i style="background:rgba(0,216,255,.7)"></i>容量正常</span>' +
          '<span><i style="background:rgba(255,162,77,.8)"></i>使用率 &gt;85%</span></div>' +
        '</div></div>';

    var rm = roomOf(state.room || rooms[0].name);
    var cabs = rm.cabinets.slice(0, 48);

    var right =
      '<div class="card" style="flex:0 0 700px;min-width:0">' +
        '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">' + A.esc(rm.name) + ' · 机柜布局</h2>' +
          '<div class="card-extra"><span class="muted" style="font-size:12px">' + rm.floor + ' · ' + rm.cabCount + ' 机柜</span></div></div>' +
        '<div class="card-body" style="overflow:auto"><div class="cab-grid">' +
          cabs.map(function (c) {
            var cls = c.topCat === '服务器' ? 'is-srv' : (c.topCat === '交换机' || c.topCat === '网络安全设备' ? 'is-net' : 'is-sto');
            var bars = '';
            for (var u = 0; u < 20; u++) {
              bars += '<i class="cab-u ' + (u < Math.round(c.usedU / 42 * 20) ? 'is-used ' + cls : '') + '"></i>';
            }
            return '<div class="cab" data-cab="' + A.esc(c.name) + '">' +
              '<div class="cab-nm">' + A.esc(c.name.split('-').slice(-2).join('-')) + '</div>' +
              '<div class="cab-stack">' + bars + '</div>' +
              '<div class="cab-foot"><span>' + c.dev + '台</span><span class="cab-rate">' + c.rate + '%</span></div>' +
              '</div>';
          }).join('') +
        '</div></div></div>';

    return {
      stats: A.UI.stats([
        { v: rooms.length, u: '个', l: '机房总数' },
        { v: A.num(totalCab), u: '个', l: '机柜总数' },
        { v: A.num(totalDev), u: '台', l: '在架设备' },
        { v: A.fix(avgRate, 1) + '%', l: '平均空间使用率' },
        { v: A.num(totalKw), u: 'kW', l: '额定功率合计' }
      ]),
      body: '<div class="row grow">' + left + right + '</div>'
    };
  }

  function viewTH() {
    var list = D().thDevices;
    var cur = state.room || (D().rooms[0] && D().rooms[0].name);
    var inRoom = list.filter(function (t) { return t.room === cur; });
    var src = inRoom.length ? inRoom : list.slice(0, 48);
    var avgT = src.reduce(function (a, t) { return a + t.temp; }, 0) / (src.length || 1);
    var avgH = src.reduce(function (a, t) { return a + t.humi; }, 0) / (src.length || 1);
    var warm = src.filter(function (t) { return t.temp > 26; }).length;

    var tSeries = src.length ? src[0].tSeries : A.wave(23, 24, 1, 3);
    var hSeries = src.length ? src[0].hSeries : A.wave(48, 24, 2, 5);

    var body =
      '<div class="row grow">' +
        '<div class="card" style="flex:1 1 0">' +
          '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">' + A.esc(cur) + ' · 温湿度点位</h2>' +
            '<div class="card-extra"><span class="muted" style="font-size:12px">' + src.length + ' 个测点 · 点击查看详情</span></div></div>' +
          '<div class="card-body" style="overflow:auto"><div class="th-grid">' +
            src.map(function (t) {
              return '<div class="th-card' + (t.temp > 26 ? ' is-warm' : '') + '" data-th="' + A.esc(t.name) + '">' +
                '<div class="th-nm" title="' + A.esc(t.name) + '">' + A.esc(t.name) + '</div>' +
                '<div><span class="th-v">' + t.temp + '</span><span class="pt-u">℃</span> ' +
                '<span class="th-h">' + t.humi + '</span><span class="pt-u">%</span></div>' +
                '<div class="pt-sub"><span>最高 ' + t.tempMax + '℃</span><span>最低 ' + t.tempMin + '℃</span></div>' +
                '</div>';
            }).join('') +
          '</div></div>' +
        '</div>' +
        '<div class="card" style="flex:0 0 460px">' +
          '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">24 小时趋势</h2></div>' +
          '<div class="card-body">' +
            '<div class="legend" style="margin-bottom:8px">' +
              '<span class="legend-item"><i class="legend-dot" style="background:#FFA24D"></i><span class="legend-text">温度 ℃</span></span>' +
              '<span class="legend-item"><i class="legend-dot" style="background:#2E9CF6"></i><span class="legend-text">湿度 %</span></span></div>' +
            '<div style="flex:1 1 auto;min-height:0;height:0">' + A.Charts.line({
              w: 428, h: 260, min: 10, max: 80, yLabels: false,
              labels: ['00:00', '04:00', '08:00', '12:00', '16:00', '20:00', '24:00'],
              series: [{ color: '#FFA24D', area: true, data: tSeries }, { color: '#2E9CF6', data: hSeries }]
            }) + '</div>' +
            '<div style="margin-top:12px">' + A.UI.bars([
              { label: '冷通道均温', value: A.fix(avgT - 3.2, 1), cls: 'p' },
              { label: '热通道均温', value: A.fix(avgT + 4.6, 1), cls: 'o' },
              { label: '平均湿度', value: A.fix(avgH, 1), cls: '' }
            ], { max: 60, unit: '' }) + '</div>' +
          '</div>' +
        '</div>' +
      '</div>';

    return {
      stats: A.UI.stats([
        { v: src.length, u: '个', l: '温湿度测点' },
        { v: A.fix(avgT, 1), u: '℃', l: '平均温度' },
        { v: A.fix(avgH, 1), u: '%', l: '平均湿度' },
        { v: warm, u: '个', l: '温度偏高测点', color: warm ? 'var(--orange)' : '' },
        { v: '正常', l: '环境状态', color: 'var(--green)' }
      ]),
      body: body
    };
  }

  /* 通用：设备 → 点位 时序视图 */
  function viewSeries(key, unitHint) {
    var s = D().series[key] || { dates: [], devices: [] };
    var devs = s.devices;
    var dates = s.dates.map(function (d) { return d.slice(5, 10); });

    var totalPt = devs.reduce(function (a, d) { return a + d.points.length; }, 0);
    var warn = devs.reduce(function (a, d) {
      return a + d.points.filter(function (p) { return p.max > p.avg * 1.15; }).length;
    }, 0);

    var curDev = state.dev && devs.filter(function (d) { return d.name === state.dev; })[0];
    if (!curDev) curDev = devs[0];

    var left =
      '<div class="card" style="flex:0 0 420px">' +
        '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">设备列表</h2>' +
          '<div class="card-extra"><span class="muted" style="font-size:12px">' + devs.length + ' 台</span></div></div>' +
        '<div class="card-body" style="overflow:auto"><div class="tree">' +
          devs.map(function (d) {
            return '<div class="tree-node' + (curDev && d.name === curDev.name ? ' is-active' : '') + '" data-dev="' + A.esc(d.name) + '">' +
              '<span class="dot"></span><span class="nm" title="' + A.esc(d.name) + '">' + A.esc(d.name) + '</span>' +
              '<span class="ct">' + d.points.length + '</span></div>';
          }).join('') +
        '</div></div></div>';

    var pts = curDev ? curDev.points : [];
    var series = pts.length ? pts[0].vals : [];

    var right =
      '<div class="card grow">' +
        '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">' +
          A.esc(curDev ? curDev.name : '设备') + ' · 监测点位</h2>' +
          '<div class="card-extra"><span class="muted" style="font-size:12px">' +
            (curDev ? A.esc(curDev.loc) : '') + ' · ' + pts.length + ' 个点位</span></div></div>' +
        '<div class="card-body col" style="gap:12px">' +
          '<div style="flex:0 0 240px">' + A.Charts.line({
            w: 900, h: 230, labels: dates, min: null, max: null,
            series: pts.slice(0, 3).map(function (p, i) {
              return { color: ['#00D8FF', '#00E39A', '#FFA24D'][i], data: p.vals, area: i === 0 };
            })
          }) + '</div>' +
          '<div style="flex:0 0 190px;min-height:0;overflow:auto"><div class="pt-grid">' +
            pts.map(function (p) {
              var isWarn = p.max > p.avg * 1.15;
              return '<div class="pt' + (isWarn ? ' is-warn' : '') + '" data-pt="' + A.esc(p.name) + '">' +
                '<div class="pt-nm" title="' + A.esc(p.name) + '">' + A.esc(p.name) + '</div>' +
                '<div class="pt-v">' + A.fix(p.last, 2) + '<span class="pt-u">' + A.esc(p.unit || unitHint || '') + '</span></div>' +
                '<div class="pt-sub"><span>均 ' + p.avg + '</span><span>峰 ' + p.max + '</span><span>谷 ' + p.min + '</span></div>' +
                '</div>';
            }).join('') +
          '</div></div>' +
          '<div style="flex:1 1 auto;min-height:0;overflow:auto">' + A.UI.table({
            cols: [
              { t: '点位名称', k: 'name' },
              { t: '最新值', k: 'last', align: 'r', fmt: function (r) { return A.fix(r.last, 3); } },
              { t: '最大值', k: 'max', align: 'r' },
              { t: '最小值', k: 'min', align: 'r' },
              { t: '平均值', k: 'avg', align: 'r' },
              { t: '状态', k: 'x', html: function (r) { return r.max > r.avg * 1.15 ? A.UI.tag('波动', 'orange') : A.UI.tag('正常', 'green'); } }
            ],
            rows: pts
          }) + '</div>' +
        '</div>' +
      '</div>';

    return {
      stats: A.UI.stats([
        { v: devs.length, u: '台', l: '监测设备' },
        { v: totalPt, u: '个', l: '点位总数' },
        { v: warn, u: '个', l: '波动点位', color: warn ? 'var(--orange)' : '' },
        { v: dates.length, u: '天', l: '统计周期' },
        { v: '在线', l: '采集状态', color: 'var(--green)' }
      ]),
      body: '<div class="row grow">' + left + right + '</div>'
    };
  }

  /* 状态型视图：漏水 / 防雷 / UPS / 配电柜 */
  function viewStatus(kind) {
    var rooms = D().rooms.filter(function (r) { return r.cabCount > 4; });
    var conf = {
      leak: { name: '漏水监测', unit: '路', n: 6, ok: '未检测到漏水' },
      lp:   { name: '防雷监测', unit: '路', n: 4, ok: '防雷正常' },
      ups:  { name: 'UPS 监测', unit: '台', n: 8, ok: '市电正常' },
      pdu:  { name: '配电柜监测', unit: '面', n: 5, ok: '供电正常' }
    }[kind];

    var items = [];
    rooms.forEach(function (r) {
      for (var i = 1; i <= conf.n; i++) {
        var st = (r.name.length + i) % 17 === 0 ? 'warn' : 'ok';
        items.push({
          room: r.name, floor: r.floor,
          name: r.name + '-' + conf.name.slice(0, 2) + '-' + i,
          load: Math.round(30 + ((r.name.length * 7 + i * 13) % 60)),
          v: A.fix(220 + ((i * 3) % 12) + Math.random(), 1),
          a: A.fix(20 + ((i * 7) % 80), 1),
          st: st
        });
      }
    });

    var warnN = items.filter(function (x) { return x.st === 'warn'; }).length;
    var avgLoad = items.reduce(function (a, x) { return a + x.load; }, 0) / (items.length || 1);

    var body =
      '<div class="card grow">' +
        '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">' + conf.name + '总览</h2>' +
          '<div class="card-extra"><span class="muted" style="font-size:12px">' + items.length + ' ' + conf.unit + ' · 点击查看详情</span></div></div>' +
        '<div class="card-body">' + A.UI.table({
          cols: [
            { t: '监测点', k: 'name', w: 320 },
            { t: '所属机房', k: 'room', w: 180 },
            { t: '楼层', k: 'floor', w: 80 },
            { t: '电压(V)', k: 'v', align: 'r', w: 100 },
            { t: '电流(A)', k: 'a', align: 'r', w: 100 },
            { t: '负载率', k: 'load', align: 'r', w: 100, html: function (r) {
              return '<span style="color:' + (r.load > 80 ? 'var(--orange)' : 'var(--cyan)') + ';font-family:var(--num)">' + r.load + '%</span>';
            } },
            { t: '状态', k: 'st', w: 120, html: function (r) {
              return r.st === 'warn' ? A.UI.tag('告警', 'red') : A.UI.tag('正常', 'green');
            } }
          ],
          rows: items,
          rowAttr: function (r) { return 'data-mp="' + A.esc(r.name) + '"'; }
        }) + '</div>' +
      '</div>';

    return {
      stats: A.UI.stats([
        { v: items.length, u: conf.unit, l: '监测点总数' },
        { v: items.filter(function (x) { return x.st === 'ok'; }).length, u: conf.unit, l: '正常', color: 'var(--green)' },
        { v: warnN, u: conf.unit, l: '告警', color: warnN ? 'var(--lv1)' : '' },
        { v: A.fix(avgLoad, 1) + '%', l: '平均负载率' },
        { v: conf.ok, l: '总体状态', color: 'var(--green)' }
      ]),
      body: body,
      items: items, conf: conf
    };
  }

  /* ---------------- 告警 ---------------- */
  function alarmsFiltered() {
    var list = alarmSource();
    if (state.level) list = list.filter(function (a) { return a.level === state.level; });
    if (state.cat) list = list.filter(function (a) { return a.cat === state.cat; });
    if (state.kw) {
      var k = state.kw.toLowerCase();
      list = list.filter(function (a) {
        return (a.what + a.obj + a.room + a.floor).toLowerCase().indexOf(k) >= 0;
      });
    }
    return list;
  }

  function viewAlarm() {
    var st = alarmStatSource();
    var list = alarmsFiltered();
    var pg = A.paginate(list, state.page, 18);

    var filter =
      '<div class="card" style="flex:0 0 auto">' +
        '<div class="card-body row ac" style="gap:12px;padding:2px 0">' +
          '<div class="field"><span class="field-label">等级</span>' +
            A.UI.select('level', [{ v: '', t: '全部' }].concat(st.byLevel.map(function (x) { return { v: x.name, t: x.name }; })), state.level) + '</div>' +
          '<div class="field"><span class="field-label">分类</span>' +
            A.UI.select('cat', [{ v: '', t: '全部' }].concat(st.byCat.map(function (x) { return { v: x.name, t: x.name }; })), state.cat) + '</div>' +
          '<div class="field"><span class="field-label">关键字</span>' +
            '<input class="input" data-f="kw" placeholder="搜索对象 / 内容 / 机房" value="' + A.esc(state.kw) + '"/></div>' +
          '<button class="btn is-primary" data-act="search">查询</button>' +
          '<button class="btn is-ghost" data-act="reset">重置</button>' +
          '<div style="margin-left:auto;display:flex;gap:8px">' +
            (A.can('edit') ? '<button class="btn is-ghost is-sm" data-act="accept">批量受理</button>' : '') +
            '<button class="btn is-ghost is-sm" data-act="export">导出</button>' +
          '</div>' +
        '</div>' +
      '</div>';

    var body = filter +
      '<div class="card grow">' +
        '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">告警列表</h2>' +
          '<div class="card-extra"><span class="muted" style="font-size:12px">点击行查看处理详情</span></div></div>' +
        '<div class="card-body">' + A.UI.table({
          cols: [
            { t: '等级', k: 'level', w: 84, html: function (r) {
              var c = { 紧急: 'red', 严重: 'orange', 重要: 'orange', 次要: 'blue', 预警: '' }[r.level] || 'grey';
              return '<span class="lv-' + A.esc(r.level) + '" style="font-weight:600">' + A.esc(r.level) + '</span>';
            } },
            { t: '告警内容', k: 'what', w: 300 },
            { t: '监控对象', k: 'obj', w: 240 },
            { t: '位置', k: 'room', w: 180, html: function (r) { return A.esc([r.floor, r.room, r.zone].filter(Boolean).join(' / ')); } },
            { t: '分类', k: 'cat', w: 110 },
            { t: '产生时间', k: 'time', w: 160, mono: true },
            { t: '时长', k: 'dur', w: 100 },
            { t: '处理状态', k: 'state', w: 100, html: function (r) { return A.UI.tag(r.state || '-', r.state === '已处理' ? 'green' : 'grey'); } },
            { t: '恢复', k: 'recover', w: 90, html: function (r) { return r.recover === '已恢复' ? A.UI.tag('已恢复', 'green') : A.UI.tag('未恢复', 'red'); } }
          ],
          rows: pg.rows,
          rowAttr: function (r, i) { return 'data-al="' + ((pg.page - 1) * 18 + i) + '" style="cursor:pointer"'; }
        }) + '</div>' +
      '</div>' +
      '<div class="card" style="flex:0 0 auto;padding:10px 16px"><div class="card-body" style="padding:0">' +
        A.UI.pager(pg) + '</div></div>';

    return {
      stats: A.UI.stats([
        { v: A.num(st.total), u: '条', l: '告警总数' },
        { v: A.num(st.byLevel[0] ? st.byLevel[0].value : 0), u: '条', l: st.byLevel[0] ? st.byLevel[0].name : '最高等级', color: 'var(--lv2)' },
        { v: A.num(st.recovered), u: '条', l: '已恢复', color: 'var(--green)' },
        { v: A.num(st.pending), u: '条', l: '未恢复', color: 'var(--orange)' },
        { v: st.total ? A.fix(st.recovered / st.total * 100, 1) + '%' : '0%', l: '恢复率' }
      ]),
      body: body
    };
  }

  /* ---------------- 页面 ---------------- */
  A.register('facility', {
    render: function (root) {
      if (!state.room) state.room = D().rooms[0] && D().rooms[0].name;
      if (A.subArg && SUBS.some(function (s) { return s.id === A.subArg; })) state.sub = A.subArg;

      root.innerHTML =
        '<div class="page-head">' +
          A.originBack() +
          '<h1 class="page-title">运行监控</h1>' +
          '<span class="page-crumb">' + A.crumb('实时监控') +
          '<div class="page-tools"><span id="clock" class="muted" style="font-family:var(--num)"></span></div>' +
        '</div>' +
        '<div class="subtabs">' +
          SUBS.map(function (s) {
            return '<button class="subtab' + (s.id === state.sub ? ' is-active' : '') + '" data-sub="' + s.id + '">' + s.name + '</button>';
          }).join('') +
        '</div>' +
        '<div class="row grow" style="min-height:0">' +
          '<aside class="side" id="facTree">' + treeHtml() + '</aside>' +
          '<section class="col grow" style="gap:16px">' +
            '<div id="facStats" style="flex:0 0 auto"></div>' +
            '<div id="facBody" class="grow" style="display:flex;flex-direction:column;gap:16px;min-height:0"></div>' +
          '</section>' +
        '</div>';

      this.paint(root);
    },

    paint: function (root) {
      var v;
      var s = state.sub;
      if (s === 'space') v = viewSpace();
      else if (s === 'th') v = viewTH();
      else if (s === 'cab') v = viewSeries('列头柜', 'V');
      else if (s === 'batt') v = viewSeries('蓄电池', 'V');
      else if (s === 'ac') v = viewSeries('空调', '℃');
      else if (s === 'h2') v = viewSeries('氢气', '%');
      else if (s === 'ups') v = (A.liveView && A.liveView('ups', 'UPS')) || viewStatus('ups');
      else if (s === 'pdu') v = viewStatus('pdu');
      else if (s === 'leak') v = (A.liveView && A.liveView('leak', '漏水监测')) || viewStatus('leak');
      else if (s === 'lp') v = viewStatus('lp');
      else if (A.facilityExtra && A.facilityExtra[s]) v = A.facilityExtra[s]();
      else v = viewAlarm();

      A.q('#facStats', root).innerHTML = v.stats;
      A.q('#facBody', root).innerHTML = v.body;
      this._items = v.items;
      this._conf = v.conf;
      this.bind(root);
    },

    bind: function (root) {
      var self = this;

      /* 子页签 */
      A.qa('.subtab', root).forEach(function (b) {
        b.addEventListener('click', function () {
          state.sub = b.getAttribute('data-sub'); state.page = 1; state.dev = null;
          A.qa('.subtab', root).forEach(function (x) { x.classList.toggle('is-active', x === b); });
          self.paint(root);
        });
      });

      /* 树：楼层展开 */
      A.qa('#facTree .tree-node[data-floor]', root).forEach(function (n) {
        n.addEventListener('click', function () {
          var f = n.getAttribute('data-floor');
          state.floor = state.floor === f ? null : f;
          A.q('#facTree', root).innerHTML = treeHtml();
          self.bind(root);
        });
      });

      /* 树 / 平面图：选机房 */
      A.qa('#facTree .tree-node[data-room]', root).forEach(function (n) {
        n.addEventListener('click', function () {
          state.room = n.getAttribute('data-room');
          A.q('#facTree', root).innerHTML = treeHtml();
          self.paint(root);
        });
      });
      A.qa('.zone[data-room]', root).forEach(function (n) {
        n.addEventListener('click', function () {
          state.room = n.getAttribute('data-room');
          self.paint(root);
        });
      });

      /* 设备列表选择 */
      A.qa('.tree-node[data-dev]', root).forEach(function (n) {
        n.addEventListener('click', function () {
          state.dev = n.getAttribute('data-dev');
          self.paint(root);
        });
      });

      /* 跳转空间巡检 */
      var bt2 = A.q('#btnTour2', root);
      if (bt2) bt2.addEventListener('click', function (e) { e.stopPropagation(); A.go('building'); });

      /* 机柜点击 → 抽屉 U 位详图 */
      A.qa('.cab[data-cab]', root).forEach(function (n) {
        n.addEventListener('click', function () { self.cabDrawer(n.getAttribute('data-cab')); });
      });

      /* 温湿度点 */
      A.qa('.th-card[data-th]', root).forEach(function (n) {
        n.addEventListener('click', function () {
          var t = D().thDevices.filter(function (x) { return x.name === n.getAttribute('data-th'); })[0];
          if (!t) return;
          A.drawer(t.name, A.UI.kv([
            ['测点名称', t.name], ['所属机房', t.room], ['通道类型', t.ch],
            ['平均温度', t.temp + ' ℃'], ['温度区间', t.tempMin + ' ~ ' + t.tempMax + ' ℃'],
            ['平均湿度', t.humi + ' %'], ['湿度区间', t.humiMin + ' ~ ' + t.humiMax + ' %'],
            ['状态', t.temp > 26 ? '温度偏高，建议检查空调' : '正常']
          ]) + '<div style="margin-top:16px">' + A.Charts.line({
            w: 476, h: 200, yLabels: false, labels: ['00:00', '04:00', '08:00', '12:00', '16:00', '20:00'],
            series: [{ color: '#FFA24D', area: true, data: t.tSeries }]
          }) + '</div>');
        });
      });

      /* 监测点行 */
      A.qa('tr[data-mp]', root).forEach(function (tr) {
        tr.addEventListener('click', function () {
          var name = tr.getAttribute('data-mp');
          var it = (self._items || []).filter(function (x) { return x.name === name; })[0];
          if (!it) return;
          A.drawer(it.name, A.UI.kv([
            ['监测点', it.name], ['所属机房', it.room], ['楼层', it.floor],
            ['电压', it.v + ' V'], ['电流', it.a + ' A'], ['负载率', it.load + ' %'],
            ['状态', it.st === 'warn' ? '告警' : '正常']
          ]) + '<div style="margin-top:16px">' + A.Charts.line({
            w: 476, h: 190, yLabels: false, labels: ['00:00', '04:00', '08:00', '12:00', '16:00', '20:00'],
            series: [{ color: '#00D8FF', area: true, data: A.wave(it.load, 24, 4, it.name.length) }]
          }) + '</div>');
        });
      });

      /* 告警筛选 */
      var btnSearch = A.q('[data-act="search"]', root);
      if (btnSearch) {
        btnSearch.addEventListener('click', function () {
          state.level = A.q('[data-f="level"]', root).value;
          state.cat = A.q('[data-f="cat"]', root).value;
          state.kw = A.q('[data-f="kw"]', root).value.trim();
          state.page = 1; self.paint(root);
        });
      }
      var btnReset = A.q('[data-act="reset"]', root);
      if (btnReset) {
        btnReset.addEventListener('click', function () {
          state.level = ''; state.cat = ''; state.kw = ''; state.page = 1; self.paint(root);
        });
      }
      var btnAccept = A.q('[data-act="accept"]', root);
      if (btnAccept) btnAccept.addEventListener('click', function () { A.toast('已批量受理当前筛选结果'); });
      var btnExport = A.q('[data-act="export"]', root);
      if (btnExport) btnExport.addEventListener('click', function () { A.toast('告警报表已生成，可到报表报告查看'); });

      /* 告警行 */
      A.qa('tr[data-al]', root).forEach(function (tr) {
        tr.addEventListener('click', function () {
          var a = alarmSource()[Number(tr.getAttribute('data-al'))];
          if (!a) return;
          A.drawer(a.level + '告警 · ' + a.obj, A.UI.kv([
            ['告警等级', '<span class="lv-' + A.esc(a.level) + '">' + A.esc(a.level) + '</span>', true],
            ['告警内容', a.what], ['监控对象', a.obj],
            ['位置', [a.floor, a.room, a.zone].filter(Boolean).join(' / ')],
            ['告警分类', a.cat], ['告警类型', a.kind],
            ['产生时间', a.time], ['告警时长', a.dur], ['恢复时间', a.rtime],
            ['触发值 / 当前值', a.tv + ' / ' + a.cv],
            ['受理人', a.handler], ['受理状态', a.accept], ['确认状态', a.confirm], ['恢复状态', a.recover]
          ].map(function (p) { return p; })) +
          '<div style="margin-top:18px;display:flex;gap:8px">' +
            (A.can('edit')
              ? '<button class="btn is-primary" data-do="accept">受理</button>' +
                '<button class="btn" data-do="confirm">确认</button>'
              : '') +
            '<button class="btn is-ghost" data-do="close">关闭</button>' +
          '</div>');
          A.qa('[data-do]').forEach(function (b) {
            b.addEventListener('click', function () {
              var d = b.getAttribute('data-do');
              if (d === 'close') A.closeDrawer();
              else A.toast('已' + (d === 'accept' ? '受理' : '确认') + '告警：' + a.obj);
            });
          });
        });
      });

      /* 分页 */
      A.qa('.pg', root).forEach(function (b) {
        b.addEventListener('click', function () {
          if (b.disabled) return;
          state.page = Number(b.getAttribute('data-p')); self.paint(root);
        });
      });
    },

    cabDrawer: function (cabName) {
      var devs = D().assets.filter(function (a) { return a.cabinet === cabName; })
        .sort(function (a, b) { return b.uStart - a.uStart; });
      var rm = D().rooms.filter(function (r) {
        return r.cabinets.some(function (c) { return c.name === cabName; });
      })[0];
      var cab = rm && rm.cabinets.filter(function (c) { return c.name === cabName; })[0];
      var kw = devs.reduce(function (a, d) { return a + (d.kw || 0); }, 0);

      var rows = [];
      var byU = {};
      devs.forEach(function (d) { byU[d.uStart] = d; });
      for (var u = 42; u >= 1; u--) {
        var d = byU[u];
        rows.push('<div class="u-row ' + (d ? (d.cat === '服务器' ? 'is-srv' : d.cat === '存储设备' ? 'is-sto' : 'is-net') : 'is-empty') + '">' +
          '<span class="u-no">' + u + '</span>' +
          '<span class="u-nm">' + (d ? A.esc(d.model + ' · ' + d.sn) + (d.source === 'sim' ? ' ' + A.simTag('sim') : '') : '空位') + '</span>' +
          (d ? '<span class="u-kw">' + d.kw + 'kW</span>' : '') + '</div>');
      }

      A.drawer('机柜 ' + cabName, A.UI.kv([
        ['机柜编号', cabName], ['所属机房', rm ? rm.name : '-'], ['楼层', rm ? rm.floor : '-'],
        ['设备数量', (cab ? cab.dev : devs.length) + ' 台'], ['已用 U 位', (cab ? cab.usedU : 0) + ' / 42 U'],
        ['空间使用率', (cab ? cab.rate : 0) + ' %'], ['额定功率', A.fix(kw, 2) + ' kW'],
        ['主要类型', cab ? cab.topCat : '-']
      ]) + '<div style="margin-top:14px;font-size:13px;color:var(--t-label);margin-bottom:8px">U 位排布（自上而下 42U → 1U）</div>' +
      '<div class="rack"><div class="rack-u">' + rows.join('') + '</div></div>');
    }
  });
})(window.App);
