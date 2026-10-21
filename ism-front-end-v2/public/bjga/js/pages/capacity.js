/* =========================================================
   容量总览 —— 按《资产管理·口径定稿》执行
   · 空间利用率 = (Σ设备U + 台数 × 空U系数) ÷ 42
   · 综合使用率 = max(空间, 电力, 接口)  取极高，不是平均
   · 三数并看：额定 / 平均负荷 / 峰值
   · 四级色阶：灰未启用 / 蓝 <70 / 黄 70~90 / 红 >90
   ========================================================= */
(function (A) {
  'use strict';

  var LV = {
    red:    { c: 'var(--lv1)',    t: '满载',   tag: 'red' },
    yellow: { c: 'var(--orange)', t: '偏高',   tag: 'orange' },
    blue:   { c: 'var(--cyan)',   t: '正常',   tag: 'blue' },
    off:    { c: 'var(--t-muted)', t: '未启用', tag: 'grey' }
  };

  function lvOf(v, off) {
    if (off) return 'off';
    return v > 90 ? 'red' : v >= 70 ? 'yellow' : 'blue';
  }

  function pct(v, off) {
    var l = LV[lvOf(v, off)];
    return '<span style="font-family:var(--num);color:' + l.c + '">' + A.fix(v, 1) + '%</span>';
  }

  /* 可部署预测：连续空间 + 电力余量 + 接口余量 三者取小 */
  function deployable(room, size) {
    var n = 0, cabs = room.cabinets || [];
    for (var i = 0; i < cabs.length; i++) {
      var c = cabs[i];
      if (!c.dev) {                          // 空柜：全新可用
        var byU0 = Math.floor(42 / (size + 1));
        var byP0 = Math.floor((c.rated || 4) / 0.6);
        var byI0 = 20;
        n += Math.min(byU0, byP0, byI0);
        continue;
      }
      var coef = c.coef === undefined ? 1 : c.coef;
      var freeU = 42 - c.usedU - c.dev * coef;
      var byU = Math.floor(freeU / (size + 1));          // 连续空段，不把零散 U 相加
      var byP = Math.floor(Math.max(0, (c.rated || 4) - (c.rtKw || 0)) / 0.6);
      var byI = Math.max(0, 20 - (c.ifUsed || 0));
      n += Math.max(0, Math.min(byU, byP, byI));
    }
    return n;
  }

  function capRoute() {
    var raw = (location.hash || '').replace(/^#\//, '');
    var parts = raw.split('/');
    var name = parts[2] || '';
    try { name = decodeURIComponent(name); } catch (e) {}
    return { kind: parts[0] === 'capacity' ? (parts[1] || '') : '', name: name };
  }

  function roomsOnFloor(floor) {
    var all = (window.DCIM.rooms || []).filter(function (r) { return r.floor === floor; });
    var main = all.filter(function (r) { return r.cabCount > 2; });
    return (main.length ? main : all).slice().sort(function (a, b) { return b.comp - a.comp; });
  }

  /* 与总览楼层卡同一套加权：空间按 U，电力按柜位，接口取平均，综合取极高 */
  function floorRollup(rooms) {
    var uEff = 0, uCap = 0, rt = 0, rated = 0, ifU = 0, ifT = 0, cab = 0;
    rooms.forEach(function (r) {
      var w = (r.cabCount || 0) * 42;
      cab += r.cabCount || 0;
      uEff += r.usedUEff || 0;
      uCap += w;
      rt += (r.powRate || 0) * w / 100;
      rated += w;
      ifU += r.ifRate || 0;
      ifT++;
    });
    var spaceRate = uCap ? uEff / uCap * 100 : 0;
    var powRate = rated ? rt / rated * 100 : 0;
    var ifRate = ifT ? ifU / ifT : 0;
    return {
      cab: cab, n: rooms.length,
      spaceRate: spaceRate, powRate: powRate, ifRate: ifRate,
      comp: Math.max(spaceRate, powRate, ifRate)
    };
  }

  function capHead(title, sub, backLabel, backTo) {
    return '<div class="page-head">' +
      '<h1 class="page-title">' + A.esc(title) + '</h1>' +
      '<span class="page-crumb">' + A.crumb(sub) + '</span>' +
      '<div class="page-tools"><button class="btn is-ghost is-sm" type="button" data-act="cap-back" data-to="' +
        A.esc(backTo) + '">' + A.esc(backLabel) + '</button></div></div>';
  }

  function dimLine(name, v, sub) {
    var l = LV[lvOf(v)];
    return '<div class="dim-row"><span class="dim-nm">' + name + '</span>' +
      '<span class="dim-track"><i style="width:' + A.clamp(v, 1, 100) + '%;background:' + l.c + '"></i></span>' +
      '<span class="dim-v" style="color:' + l.c + '">' + A.fix(v, 1) + '%</span>' +
      '<span class="dim-sub">' + sub + '</span></div>';
  }

  function bindCapBack(root) {
    A.qa('[data-act="cap-back"]', root).forEach(function (b) {
      b.addEventListener('click', function () { A.go(b.getAttribute('data-to') || 'capacity'); });
    });
  }

  function bindCapRooms(root) {
    A.qa('.cap-room', root).forEach(function (el) {
      el.addEventListener('click', function () {
        var room = el.getAttribute('data-room');
        if (room) A.go('capacity/room/' + encodeURIComponent(room));
      });
    });
  }

  function bindCapFloors(root) {
    A.qa('[data-cap-floor]', root).forEach(function (el) {
      el.addEventListener('click', function () {
        var floor = el.getAttribute('data-cap-floor');
        if (floor) A.go('capacity/floor/' + encodeURIComponent(floor));
      });
    });
  }

  function renderFloor(root, floor) {
    var D = window.DCIM || {};
    var rooms = roomsOnFloor(floor);
    var known = rooms.length || (D.floors || []).some(function (f) { return f.floor === floor; });
    if (!known) {
      root.innerHTML = capHead((floor || '楼层') + ' 容量', '未找到该楼层', '返回容量总览', 'capacity') +
        '<div class="empty">没有这一层的容量数据</div>';
      bindCapBack(root);
      return;
    }
    var f = floorRollup(rooms);
    var lv = LV[lvOf(f.comp)];
    root.innerHTML = capHead(floor + ' 容量', '该层机房的空间 · 电力 · 综合占用', '返回容量总览', 'capacity') +
      '<div class="card" style="flex:0 0 auto"><div class="card-body" style="padding:2px 0">' +
        A.UI.stats([
          { v: rooms.length, u: '间', l: '机房' },
          { v: f.cab, u: '个', l: '机柜' },
          { v: A.fix(f.comp, 1) + '%', l: '综合使用率（取极高）', color: lv.c },
          { v: A.fix(f.spaceRate, 1) + '%', l: '空间使用率' },
          { v: A.fix(f.powRate, 1) + '%', l: '电力使用率' },
          { v: A.fix(f.ifRate, 1) + '%', l: '接口使用率' }
        ]) +
      '</div></div>' +
      '<section class="card grow"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">' + A.esc(floor) + ' 各机房容量</h2>' +
        '<div class="card-extra"><span class="muted" style="font-size:12px">点机房查看该机房各机柜占用</span></div></div>' +
        '<div class="card-body" style="overflow:auto">' +
          A.UI.table({
            cols: [
              { t: '机房', k: 'name', w: 168, html: function (r) {
                return '<span class="cap-room" data-room="' + A.esc(r.name) + '">' + A.esc(r.name) + '</span>';
              } },
              { t: '机柜', k: 'cabCount', align: 'r', w: 62 },
              { t: '空间', k: 'spaceRate', align: 'r', w: 78, html: function (r) { return pct(r.spaceRate); } },
              { t: '电力', k: 'powRate', align: 'r', w: 78, html: function (r) { return pct(r.powRate); } },
              { t: '接口', k: 'ifRate', align: 'r', w: 78, html: function (r) { return pct(r.ifRate); } },
              { t: '综合', k: 'comp', align: 'r', w: 96, html: function (r) {
                return '<b style="font-family:var(--num);font-size:14px;color:' + LV[lvOf(r.comp)].c + '">' +
                  A.fix(r.comp, 1) + '%</b>';
              } },
              { t: '可部署 2U/4U', k: 'd', align: 'r', w: 116, html: function (r) {
                return '<span style="font-family:var(--num);color:var(--t-sub)">' +
                  deployable(r, 2) + ' / ' + deployable(r, 4) + ' 台</span>';
              } },
              { t: '状态', k: 'lv', w: 86, html: function (r) {
                var l = LV[lvOf(r.comp)]; return A.UI.tag(l.t, l.tag);
              } }
            ],
            rows: rooms
          }) +
          '<div class="cap-foot">综合使用率 = max(空间, 电力, 接口)。数字来自现有机房容量，不另建数据集。</div>' +
        '</div></section>';
    bindCapBack(root);
    bindCapRooms(root);
  }

  function renderRoom(root, name) {
    var D = window.DCIM || {};
    var r = (D.rooms || []).filter(function (x) { return x.name === name; })[0];
    if (!r) {
      root.innerHTML = capHead((name || '机房') + ' 容量', '未找到该机房', '返回容量总览', 'capacity') +
        '<div class="empty">没有这一间的容量数据</div>';
      bindCapBack(root);
      return;
    }
    var back = 'capacity/floor/' + encodeURIComponent(r.floor);
    var cabs = (r.cabinets || []).slice().sort(function (a, b) { return (b.comp || 0) - (a.comp || 0); });
    var lv = LV[lvOf(r.comp)];
    root.innerHTML = capHead(r.name + ' 容量', r.floor + ' · 各机柜空间 · 电力 · 综合占用', '返回 ' + r.floor, back) +
      '<div class="card" style="flex:0 0 auto"><div class="card-body" style="padding:2px 0">' +
        A.UI.stats([
          { v: r.cabCount || cabs.length, u: '个', l: '机柜' },
          { v: A.num(r.dev || 0), u: '台', l: '设备' },
          { v: A.fix(r.comp, 1) + '%', l: '综合使用率（取极高）', color: lv.c },
          { v: A.fix(r.spaceRate, 1) + '%', l: '空间使用率' },
          { v: A.fix(r.powRate, 1) + '%', l: '电力使用率' },
          { v: A.fix(r.ifRate, 1) + '%', l: '接口使用率' }
        ]) +
      '</div></div>' +
      '<div class="row grow">' +
        '<section class="card" style="flex:0 0 372px"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">三维占用</h2></div>' +
          '<div class="card-body col" style="gap:14px">' +
            dimLine('空间', r.spaceRate, '含空U系数') +
            dimLine('电力', r.powRate, A.fix(r.kw, 1) + ' kW 额定侧') +
            dimLine('接口', r.ifRate, '口占用') +
            '<div class="dp-note">可再部署 2U ' + deployable(r, 2) + ' 台 · 4U ' + deployable(r, 4) +
              ' 台。按连续空段，并与电力、接口余量取小。</div>' +
          '</div></section>' +
        '<section class="card grow"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">' + A.esc(r.name) + ' 机柜占用</h2>' +
          '<div class="card-extra"><span class="muted" style="font-size:12px">' + cabs.length + ' 架</span></div></div>' +
          '<div class="card-body" style="overflow:auto">' +
            (cabs.length ? A.UI.table({
              cols: [
                { t: '机柜', k: 'name', w: 180, html: function (c) { return A.esc(c.name); } },
                { t: '设备', k: 'dev', align: 'r', w: 62, html: function (c) { return A.num(c.dev || 0); } },
                { t: '空间', k: 'spaceRate', align: 'r', w: 78, html: function (c) { return pct(c.spaceRate, !c.dev); } },
                { t: '电力', k: 'powRate', align: 'r', w: 78, html: function (c) { return pct(c.powRate, !c.dev); } },
                { t: '接口', k: 'ifRate', align: 'r', w: 78, html: function (c) { return pct(c.ifRate, !c.dev); } },
                { t: '综合', k: 'comp', align: 'r', w: 96, html: function (c) {
                  var k = c.lv && LV[c.lv] ? c.lv : lvOf(c.comp, !c.dev);
                  return '<b style="font-family:var(--num);font-size:14px;color:' + LV[k].c + '">' +
                    A.fix(c.comp || 0, 1) + '%</b>';
                } },
                { t: '状态', k: 'lv', w: 86, html: function (c) {
                  var k = c.lv && LV[c.lv] ? c.lv : lvOf(c.comp, !c.dev);
                  return A.UI.tag(LV[k].t, LV[k].tag);
                } }
              ],
              rows: cabs
            }) : '<div class="empty">该机房没有机柜占用数据</div>') +
            '<div class="cap-foot">机柜占用沿用该机房 cabinets 上的空间 / 电力 / 接口 / 综合。未启用柜按灰色显示。</div>' +
          '</div></section>' +
      '</div>';
    bindCapBack(root);
  }

  A.register('capacity', {
    render: function (root) {
      var route = capRoute();
      if (route.kind === 'floor') { renderFloor(root, route.name); return; }
      if (route.kind === 'room') { renderRoom(root, route.name); return; }
      var D = window.DCIM;
      var cap = D.capacity || {};
      var rooms = D.rooms.filter(function (r) { return r.cabCount > 2; })
        .slice().sort(function (a, b) { return b.comp - a.comp; });

      var cabTotal = D.meta.cabTotal, cabUsed = D.meta.cabUsed, cabFree = D.meta.cabFree;

      /* 色阶分布 */
      var dist = { red: 0, yellow: 0, blue: 0, off: 0 };
      D.rooms.forEach(function (r) {
        (r.cabinets || []).forEach(function (c) { dist[c.lv || lvOf(c.comp, !c.dev)]++; });
      });

      var d2 = rooms.reduce(function (a, r) { return a + deployable(r, 2); }, 0);
      var d4 = rooms.reduce(function (a, r) { return a + deployable(r, 4); }, 0);

      /* 楼层维度：按机房加权汇总，综合仍取极高 */
      var fl = {};
      rooms.forEach(function (r) {
        var f = fl[r.floor] = fl[r.floor] || { n: 0, cab: 0, uEff: 0, uCap: 0, rt: 0, rated: 0, ifU: 0, ifT: 0 };
        f.n++; f.cab += r.cabCount;
        f.uEff += r.usedUEff || 0; f.uCap += r.cabCount * 42;
      });
      Object.keys(fl).forEach(function (k) {
        var f = fl[k], rt = 0, rated = 0, ifU = 0, ifT = 0;
        rooms.filter(function (r) { return r.floor === k; }).forEach(function (r) {
          rt += r.powRate * (r.cabCount * 42) / 100;   // 以柜数为权重
          rated += r.cabCount * 42;
          ifU += r.ifRate; ifT++;
        });
        f.spaceRate = f.uCap ? f.uEff / f.uCap * 100 : 0;
        f.powRate = rated ? rt / rated * 100 : 0;
        f.ifRate = ifT ? ifU / ifT : 0;
        f.comp = Math.max(f.spaceRate, f.powRate, f.ifRate);
      });
      var FL = Object.keys(fl).sort().map(function (k) {
        return { floor: k, cab: fl[k].cab, comp: fl[k].comp,
          spaceRate: fl[k].spaceRate, powRate: fl[k].powRate, ifRate: fl[k].ifRate };
      });

      root.innerHTML =
        '<div class="page-head">' +
          '<h1 class="page-title">容量总览</h1>' +
          '<span class="page-crumb">' + A.crumb('空间 · 电力 · 接口 三维取极高') +
          '<div class="page-tools"><button class="btn is-primary is-sm" data-act="recalc">按台账重算</button><button class="btn is-ghost is-sm" data-act="export">导出容量明细</button></div>' +
        '</div>' +

        '<div class="card" style="flex:0 0 auto"><div class="card-body" style="padding:2px 0">' +
          A.UI.stats([
            { v: A.num(cabTotal), u: '个', l: '机柜总数' },
            { v: A.num(cabUsed), u: '个', l: '已启用', color: 'var(--cyan)' },
            { v: A.fix(cap.comp, 1) + '%', l: '综合使用率（取极高）', color: 'var(--orange)' },
            { v: A.fix(cap.spaceRate, 1) + '%', l: '空间使用率' },
            { v: A.fix(cap.powRate, 1) + '%', l: '电力使用率' },
            { v: A.fix(cap.ifRate, 1) + '%', l: '接口使用率' }
          ]) +
        '</div></div>' +

        '<div class="row grow">' +

          /* 综合 + 三维 + 色阶 */
          '<section class="card" style="flex:0 0 372px">' +
            '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">综合使用率</h2>' +
              '<div class="card-extra"><span class="muted" style="font-size:12px">取三维极高</span></div></div>' +
            '<div class="card-body col" style="gap:14px">' +
              '<div style="position:relative;width:200px;height:200px;display:grid;place-items:center;align-self:center">' +
                A.Charts.ring({ size: 200, r: 80, sw: 17, rate: cap.comp / 100,
                  color: cap.comp > 90 ? '#FF4D5E' : cap.comp >= 70 ? '#FFA24D' : '#00D8FF',
                  color2: cap.comp > 90 ? '#8B1B2B' : '#2E7CF6' }) +
                '<div style="position:absolute;text-align:center">' +
                  '<div style="font-family:var(--num);font-size:38px;font-weight:700;color:' +
                    (cap.comp > 90 ? 'var(--lv1)' : cap.comp >= 70 ? 'var(--orange)' : 'var(--t-title)') + '">' +
                    A.fix(cap.comp, 1) + '%</div>' +
                  '<div style="font-size:12px;color:var(--t-label);margin-top:2px">综合 · 均 ' + A.fix(cap.avgRate, 1) +
                    '% · 峰 ' + A.fix(cap.peakRate, 1) + '%</div>' +
                '</div>' +
              '</div>' +
              '<div style="width:100%" data-dim>' +
                dimBar('空间', cap.spaceRate, cap.usedUEff + ' / ' + A.num(cap.capU) + ' U（含空U系数）') +
                dimBar('电力', cap.powRate, A.fix(cap.rtKw, 0) + ' / ' + A.fix(cap.ratedKw, 0) + ' kW') +
                dimBar('接口', cap.ifRate, cap.ifUsed + ' / ' + A.num(cap.ifTotal) + ' 口') +
              '</div>' +
              '<div class="lv-dist" style="width:100%">' +
                lvCell('red', dist.red, '满载 >90%') +
                lvCell('yellow', dist.yellow, '偏高 70~90%') +
                lvCell('blue', dist.blue, '正常 <70%') +
                lvCell('off', dist.off, '未启用') +
              '</div>' +
            '</div>' +
          '</section>' +

          /* 各机房容量（点一行下钻） */
          '<section class="card grow">' +
            '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">各机房容量</h2>' +
              '<div class="card-extra"><span class="muted" style="font-size:12px">点机房看该机房容量 · 点楼层看该层容量</span></div></div>' +
            '<div class="card-body" style="overflow:auto">' +
              A.UI.table({
                cols: [
                  { t: '机房', k: 'name', w: 168, html: function (r) {
                    return '<span class="cap-room" data-room="' + A.esc(r.name) + '">' + A.esc(r.name) + '</span>';
                  } },
                  { t: '机柜', k: 'cabCount', align: 'r', w: 62 },
                  { t: '空间', k: 'spaceRate', align: 'r', w: 78, html: function (r) { return pct(r.spaceRate); } },
                  { t: '电力', k: 'powRate', align: 'r', w: 78, html: function (r) { return pct(r.powRate); } },
                  { t: '接口', k: 'ifRate', align: 'r', w: 78, html: function (r) { return pct(r.ifRate); } },
                  { t: '综合', k: 'comp', align: 'r', w: 96, html: function (r) {
                    return '<b style="font-family:var(--num);font-size:14px;color:' + LV[lvOf(r.comp)].c + '">' +
                      A.fix(r.comp, 1) + '%</b>';
                  } },
                  { t: '平均', k: 'avgRate', align: 'r', w: 72, html: function (r) { return A.fix(r.avgRate, 1) + '%'; } },
                  { t: '峰值', k: 'peakRate', align: 'r', w: 72, html: function (r) {
                    return '<span style="color:var(--t-sub)">' + A.fix(r.peakRate, 1) + '%</span>';
                  } },
                  { t: '可部署 2U/4U', k: 'd', align: 'r', w: 116, html: function (r) {
                    return '<span style="font-family:var(--num);color:var(--t-sub)">' +
                      deployable(r, 2) + ' / ' + deployable(r, 4) + ' 台</span>';
                  } },
                  { t: '状态', k: 'lv', w: 86, html: function (r) {
                    var l = LV[lvOf(r.comp)]; return A.UI.tag(l.t, l.tag);
                  } }
                ],
                rows: rooms
              }) +
              '<div class="cap-floor"><div class="side-title" style="margin-bottom:8px">各楼层综合使用率</div>' +
                '<div class="cf-grid">' + FL.map(function (f) {
                  var lv = lvOf(f.comp);
                  return '<div class="cf-cell is-go" data-cap-floor="' + A.esc(f.floor) + '">' +
                    '<div class="cf-h"><b>' + f.floor + '</b><span>' + f.cab + ' 柜</span></div>' +
                    '<div class="cf-v" style="color:' + LV[lv].c + '">' + A.fix(f.comp, 1) + '%</div>' +
                    '<div class="cf-track"><i style="width:' + A.clamp(f.comp, 2, 100) + '%;background:' + LV[lv].c + '"></i></div>' +
                    '<div class="cf-sub">空间 ' + A.fix(f.spaceRate, 0) + '% · 电力 ' + A.fix(f.powRate, 0) +
                      '% · 接口 ' + A.fix(f.ifRate, 0) + '%</div></div>';
                }).join('') + '</div></div>' +
              '<div class="cap-foot">综合使用率 = max(空间, 电力, 接口) —— 取极高，不做平均；' +
                '空间利用率 =（Σ设备U + 台数 × 空U系数）÷ 42；' +
                '平均 / 峰值为近 30 天统计（每日记一个峰值，创新高才更新）</div>' +
            '</div>' +
          '</section>' +

          /* 可部署预测 */
          '<section class="card" style="flex:0 0 344px">' +
            '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">可部署预测</h2></div>' +
            '<div class="card-body col" style="gap:12px">' +
              '<div class="dp-hl"><b>' + d2 + '</b><span>台 2U 设备</span>' +
                '<i>' + d4 + ' 台 4U</i></div>' +
              '<div class="dp-note">按「连续空段」计算，不把零散 U 相加；再与电力余量、接口余量取小值。</div>' +
              '<div class="side-title" style="margin-top:2px">各机房可再部署（2U）</div>' +
              A.UI.rank(rooms.slice().sort(function (a, b) {
                return deployable(b, 2) - deployable(a, 2);
              }).slice(0, 8).map(function (r) {
                return { name: r.name, value: deployable(r, 2) };
              }), ' 台') +
              '<div class="dp-formula">空间利用率 =（Σ设备U + 台数 × 空U系数）÷ 42<br>' +
                '空U系数：服务器/存储 1 · 网络 0.5 · 光纤柜 0 · 电池柜 0</div>' +
            '</div>' +
          '</section>' +

        '</div>';

      var recalc = A.q('[data-act="recalc"]', root);
      if (recalc) recalc.addEventListener('click', function () {
        A.api('/bjga/capacity/recalc', { method: 'POST', body: {} }).then(function (j) {
          if (!j || j.code !== 1000) { A.toast((j && j.message) || '重算失败'); return null; }
          return A.api('/bjga/bootstrap');
        }).then(function (j) {
          if (!j || !j.dcim) return;
          window.DCIM = j.dcim;
          A.toast('容量已按台账重算');
          A.pages.capacity.render(root);
        });
      });

      bindCapRooms(root);
      bindCapFloors(root);

      function dimBar(name, v, sub) {
        var l = LV[lvOf(v)];
        return '<div class="dim-row"><span class="dim-nm">' + name + '</span>' +
          '<span class="dim-track"><i style="width:' + A.clamp(v, 1, 100) + '%;background:' + l.c + '"></i></span>' +
          '<span class="dim-v" style="color:' + l.c + '">' + A.fix(v, 1) + '%</span>' +
          '<span class="dim-sub">' + sub + '</span></div>';
      }

      function lvCell(k, n, label) {
        return '<span class="lv-item"><em class="lv-dot" style="background:' + LV[k].c + '"></em>' +
          '<b>' + n + '</b><i>' + label + '</i></span>';
      }
    }
  });
})(window.App);
