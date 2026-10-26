/* =========================================================
   资产总览：分布 / 构成 / 归属 / 业务视角
   ========================================================= */
(function (A) {
  'use strict';

  var CAT_COLOR = {
    '服务器': '#00D8FF', '交换机': '#00E39A', '网络安全设备': '#A46BFF',
    '存储设备': '#FFA24D', '其他设备': '#2E9CF6'
  };

  A.register('assets', {
    render: function (root) {
      var D = window.DCIM, st = D.assetStats;
      var rooms = D.rooms.slice().sort(function (a, b) { return b.dev - a.dev; });

      var donutData = st.byCat.map(function (c) { return { value: c.value, color: CAT_COLOR[c.name] || '#2E9CF6' }; });

      root.innerHTML =
        '<div class="page-head">' +
          '<h1 class="page-title">资产总览</h1>' +
          '<span class="page-crumb">' + A.crumb('资产分布与构成') +
          '<div class="page-tools"><button class="btn is-ghost is-sm" data-act="export">导出台账</button></div>' +
        '</div>' +

        '<div class="card" style="flex:0 0 auto"><div class="card-body" style="padding:2px 0">' +
          A.UI.stats([
            { v: A.num(D.meta.assetTotal), u: '台', l: '设备总数' },
            { v: D.meta.roomTotal, u: '个', l: '机房数量' },
            { v: A.num(D.meta.cabTotal), u: '个', l: '机柜数量' },
            { v: st.bizCount, u: '个', l: '业务系统' },
            { v: A.num(st.powerSum), u: 'kW', l: '额定功率合计' },
            { v: A.num(st.weightSum / 1000), u: 't', l: '设备总重量' }
          ]) +
        '</div></div>' +

        '<div class="row grow">' +

          /* 机房分布 */
          '<section class="card" style="flex:1 1 0">' +
            '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">各机房设备分布</h2>' +
              '<div class="card-extra"><span class="muted" style="font-size:12px">按在架数量</span></div></div>' +
            '<div class="card-body" style="overflow:auto;gap:18px">' +
              '<div style="flex:0 0 auto">' +
                A.UI.bars(rooms.slice(0, 15).map(function (r) {
                  return { label: r.name, value: r.dev };
                }), {}) +
              '</div>' +
              '<div style="flex:0 0 auto">' +
                '<div class="side-title" style="margin-bottom:10px">各机房额定功率 (kW)</div>' +
                A.UI.bars(rooms.slice(0, 10).map(function (r) {
                  return { label: r.name, value: Math.round(r.kw), cls: 'o' };
                }), {}) +
              '</div>' +
            '</div>' +
          '</section>' +

          /* 类别构成 */
          '<section class="card" style="flex:0 0 420px">' +
            '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">设备类别构成</h2></div>' +
            '<div class="card-body" style="gap:14px">' +
              '<div class="row ac" style="gap:16px;flex:0 0 auto">' +
                '<div style="position:relative;flex:0 0 168px;width:168px;height:168px;display:grid;place-items:center">' +
                  A.Charts.donut({ size: 168, r: 62, sw: 20, data: donutData }) +
                  '<div style="position:absolute;text-align:center">' +
                    '<div style="font-family:var(--num);font-size:24px;font-weight:700;color:var(--cyan)">' + A.num(st.total) + '</div>' +
                    '<div style="font-size:11px;color:var(--t-label)">设备总数</div>' +
                  '</div>' +
                '</div>' +
                '<div style="flex:1 1 auto;min-width:0">' +
                  st.byCat.map(function (c) {
                    return '<div class="rank-item" style="margin-bottom:9px">' +
                      '<i class="dot" style="background:' + (CAT_COLOR[c.name] || '#2E9CF6') + '"></i>' +
                      '<span class="rank-nm">' + A.esc(c.name) + '</span>' +
                      '<span class="rank-v">' + A.num(c.value) + '</span>' +
                      '<span style="width:46px;text-align:right;font-family:var(--num);color:var(--t-muted);font-size:12px">' +
                        A.pct(c.value, st.total) + '%</span></div>';
                  }).join('') +
                '</div>' +
              '</div>' +
              '<div style="flex:1 1 auto;overflow:auto;min-height:0">' +
                '<div class="side-title" style="margin-bottom:10px">品牌分布 TOP 8</div>' +
                A.UI.bars(st.byBrand.slice(0, 8).map(function (b) { return { label: b.name, value: b.value, cls: 'g' }; }), {}) +
              '</div>' +
            '</div>' +
          '</section>' +

          /* 归属与业务 */
          '<section class="card" style="flex:0 0 380px">' +
            '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">使用单位 / 业务系统</h2></div>' +
            '<div class="card-body" style="gap:16px;overflow:auto">' +
              '<div style="flex:0 0 auto"><div class="side-title" style="margin-bottom:8px">按使用单位</div>' +
                A.UI.rank(st.byOrg.slice(0, 6), ' 台') + '</div>' +
              '<div style="flex:0 0 auto"><div class="side-title" style="margin-bottom:8px">按业务系统 TOP 10</div>' +
                A.UI.rank(st.byBiz.slice(0, 10), ' 台') + '</div>' +
              '<div style="flex:0 0 auto"><div class="side-title" style="margin-bottom:8px">按联网类型</div>' +
                A.UI.bars(st.byNet.map(function (n) { return { label: n.name, value: n.value, cls: 'p' }; }), {}) + '</div>' +
              '<div style="flex:0 0 auto"><div class="side-title" style="margin-bottom:8px">按使用状态</div>' +
                A.UI.bars(st.byStatus.map(function (s) {
                  return { label: s.name, value: s.value, cls: s.name === '使用中' ? 'g' : 'o' };
                }), {}) + '</div>' +
            '</div>' +
          '</section>' +

        '</div>' +

        /* 机房明细 */
        '<section class="card" style="flex:0 0 300px">' +
          '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">机房资产明细</h2>' +
            '<div class="card-extra"><span class="muted" style="font-size:12px">点击行查看该机房设备构成</span></div></div>' +
          '<div class="card-body">' + A.UI.table({
            cols: [
              { t: '机房名称', k: 'name', w: 200 },
              { t: '楼层', k: 'floor', w: 80 },
              { t: '机柜数', k: 'cabCount', align: 'r', w: 100 },
              { t: '设备数', k: 'dev', align: 'r', w: 100 },
              { t: '额定功率(kW)', k: 'kw', align: 'r', w: 130, fmt: function (r) { return A.num(r.kw); } },
              { t: '空间使用率', k: 'rate', align: 'r', w: 120, html: function (r) {
                return '<span style="color:' + (r.rate > 85 ? 'var(--orange)' : 'var(--cyan)') + ';font-family:var(--num)">' + r.rate + '%</span>';
              } },
              { t: '主要设备类型', k: 'cats', w: 200, html: function (r) {
                return A.esc(r.cats.slice(0, 2).map(function (c) { return c.name + ' ' + c.value; }).join(' · '));
              } },
              { t: '承载业务系统', k: 'biz', w: 260, html: function (r) {
                return A.esc(r.biz.slice(0, 3).map(function (b) { return b.name; }).join('、'));
              } }
            ],
            rows: rooms,
            rowAttr: function (r) { return 'data-rm="' + A.esc(r.name) + '" style="cursor:pointer"'; }
          }) + '</div>' +
        '</section>';

      /* 交互 */
      A.q('[data-act="export"]', root).addEventListener('click', function () {
        A.toast('资产台账已导出（演示）');
      });

      A.qa('tr[data-rm]', root).forEach(function (tr) {
        tr.addEventListener('click', function () {
          var rm = window.DCIM.rooms.filter(function (x) { return x.name === tr.getAttribute('data-rm'); })[0];
          if (!rm) return;
          A.drawer(rm.name + ' · 资产构成',
            A.UI.kv([
              ['机房', rm.name], ['楼层', rm.floor], ['机柜数', rm.cabCount + ' 个'],
              ['设备数', rm.dev + ' 台'], ['额定功率', A.num(rm.kw) + ' kW'],
              ['空间使用率', rm.rate + ' %'],
              ['承载业务', rm.biz.map(function (b) { return b.name + '(' + b.value + ')'; }).join('、')]
            ]) +
            '<div style="margin-top:16px"><div class="side-title" style="margin-bottom:8px">设备类别构成</div>' +
              A.UI.bars(rm.cats.map(function (c) {
                return { label: c.name, value: c.value, cls: c.name === '服务器' ? '' : c.name === '存储设备' ? 'o' : 'g' };
              }), {}) + '</div>' +
            '<div style="margin-top:16px"><div class="side-title" style="margin-bottom:8px">使用单位构成</div>' +
              A.UI.bars(rm.orgs.map(function (o) { return { label: o.name, value: o.value, cls: 'p' }; }), {}) + '</div>' +
            '<div class="am-iobtns" style="margin-top:16px">' +
              '<button class="btn is-primary is-sm" type="button" data-go-space="room" data-name="' + A.esc(rm.name) + '">查看该机房立体</button>' +
              '<button class="btn is-ghost is-sm" type="button" data-go-space="floor" data-name="' + A.esc(rm.floor) + '">查看 ' + A.esc(rm.floor) + ' 立体</button>' +
            '</div>');
          var dr = A.q('.drawer');
          if (dr && A.AM && A.AM.openSpace) {
            A.qa('[data-go-space]', dr).forEach(function (b) {
              b.addEventListener('click', function () {
                A.AM.openSpace(b.getAttribute('data-go-space'), b.getAttribute('data-name'));
              });
            });
          }
        });
      });
    }
  });
})(window.App);
