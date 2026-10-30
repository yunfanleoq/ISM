/* =========================================================
   报表报告：模板 / 周期任务 / 在线预览（真实数据统计）
   ========================================================= */
(function (A) {
  'use strict';

  var TPLS = [
    { id: 'asset', name: '资产台账报表', cycle: '月报', owner: '资产管理', desc: '按机房统计在架设备与功率' },
    { id: 'alarm', name: '告警统计日报', cycle: '日报', owner: '运维中心', desc: '按等级/分类统计当日告警' },
    { id: 'th', name: '温湿度日报表', cycle: '日报', owner: '运维中心', desc: '各测点温湿度极值与均值' },
    { id: 'energy', name: '能耗月报', cycle: '月报', owner: '技术保障科', desc: '各机房 IT 功率与电度' },
    { id: 'ups', name: 'UPS 运行报表', cycle: '周报', owner: '运维中心', desc: 'UPS 回路负载与切换记录' },
    { id: 'cap', name: '容量分析报表', cycle: '月报', owner: '技术保障科', desc: '空间与电力容量趋势' }
  ];
  var state = { tpl: 'asset', cycle: '月', page: 1 };

  function buildAsset() {
    var D = window.DCIM;
    var rows = (D.rooms || []).slice().sort(function (a, b) { return b.dev - a.dev; });
    return {
      title: '数据中心资产台账报表',
      cols: ['机房名称', '楼层', '机柜数', '设备数', '额定功率(kW)', '空间使用率', '主要类型'],
      rows: rows.map(function (r) {
        return [r.name, r.floor, r.cabCount, r.dev, A.num(r.kw), (r.rate != null ? r.rate : '-') + '%',
          (r.cats || []).slice(0, 2).map(function (c) { return c.name; }).join('/') || '-'];
      }),
      sum: ['合计', '-', A.num((D.meta || {}).cabTotal), A.num((D.meta || {}).assetTotal), A.num((D.assetStats || {}).powerSum), '-', '-']
    };
  }

  function buildAlarm() {
    var D = window.DCIM, st = D.alarmStats;
    var lvs = {}; st.byLevel.forEach(function (x) { lvs[x.name] = x.value; });
    var rows = Object.keys(lvs).map(function (k) {
      var rec = D.alarms.filter(function (a) { return a.level === k && a.recover === '已恢复'; }).length;
      var all = D.alarms.filter(function (a) { return a.level === k; }).length;
      return [k, A.num(lvs[k]), A.num(rec), A.num(all - rec), all ? A.fix(rec / all * 100, 1) + '%' : '-',
        k === '紧急' ? '10 分钟内' : k === '严重' ? '30 分钟内' : '2 小时内'];
    });
    return {
      title: '告警统计日报',
      cols: ['告警等级', '告警总数', '已恢复', '未恢复', '恢复率', '响应时限要求'],
      rows: rows,
      sum: ['合计', A.num(st.total), A.num(st.recovered), A.num(st.pending), A.fix(st.recovered / st.total * 100, 1) + '%', '-']
    };
  }

  function buildTH() {
    var list = (window.DCIM.thDevices || []).slice(0, 40);
    return {
      title: '机房温湿度日报表',
      cols: ['测点名称', '所属机房', '通道', '平均温度(℃)', '最高(℃)', '最低(℃)', '平均湿度(%)', '状态'],
      rows: list.map(function (t) {
        return [t.name, t.room, t.ch, t.temp, t.tempMax, t.tempMin, t.humi, t.temp > 26 ? '偏高' : '正常'];
      }),
      sum: ['测点数 ' + list.length, '-', '-',
        A.fix(list.reduce(function (a, t) { return a + t.temp; }, 0) / list.length, 1), '-', '-',
        A.fix(list.reduce(function (a, t) { return a + t.humi; }, 0) / list.length, 1), '-']
    };
  }

  function buildEnergy() {
    var rooms = window.DCIM.rooms.filter(function (r) { return r.cabCount > 2; });
    return {
      title: '数据中心能耗月报',
      cols: ['机房名称', '楼层', 'IT 功率(kW)', '估算月电度(万度)', '单位面积功率', 'PUE 参考', '环比'],
      rows: rooms.map(function (r, i) {
        var kwh = r.kw * 24 * 30 / 10000;
        return [r.name, r.floor, A.num(r.kw), A.fix(kwh, 2), A.fix(r.kw / Math.max(1, r.cabCount * 2.4), 2) + ' kW/㎡',
          A.fix(1.42 + (i % 5) * 0.03, 2), (i % 3 ? '+' : '-') + A.fix(1 + (i % 7) * 0.6, 1) + '%'];
      }),
      sum: ['合计', '-', A.num(window.DCIM.assetStats.powerSum),
        A.fix(window.DCIM.assetStats.powerSum * 24 * 30 / 10000, 2), '-', '1.50', '+2.3%']
    };
  }

  function buildUPS() {
    var ups = [['A路 1#UPS', 62], ['A路 2#UPS', 57], ['A路 3#UPS', 71], ['B路 1#UPS', 66],
               ['B路 2#UPS', 54], ['B路 3#UPS', 78], ['C路 1#UPS', 59], ['C路 2#UPS', 68]];
    return {
      title: 'UPS 运行周报',
      cols: ['UPS 回路', '平均负载率', '峰值负载率', '输出功率(kW)', '电池后备时间(min)', '切换次数', '状态'],
      rows: ups.map(function (u, i) {
        return [u[0], u[1] + '%', (u[1] + 8) + '%', A.fix(u[1] * 1.16, 1),
          A.fix(28 - u[1] * 0.18, 1), i % 4, u[1] > 80 ? '接近满载' : '正常'];
      }),
      sum: ['合计 / 平均', A.fix(ups.reduce(function (a, u) { return a + u[1]; }, 0) / ups.length, 1) + '%', '-',
        A.fix(ups.reduce(function (a, u) { return a + u[1] * 1.16; }, 0), 1), '-', '-', '-']
    };
  }

  function buildCap() {
    var rooms = window.DCIM.rooms.filter(function (r) { return r.cabCount > 2; });
    return {
      title: '容量分析月报',
      cols: ['机房名称', '机柜数', '已用 U 位', '剩余 U 位', '空间使用率', '可上架服务器(估)', '扩容建议'],
      rows: rooms.map(function (r) {
        var cap = r.cabCount * 42, free = cap - r.usedU;
        return [r.name, r.cabCount, r.usedU, free, r.rate + '%', Math.floor(free / 2),
          r.rate > 85 ? '建议扩容' : r.rate > 60 ? '关注' : '充足'];
      }),
      sum: ['合计', A.num(rooms.reduce(function (a, r) { return a + r.cabCount; }, 0)),
        A.num(rooms.reduce(function (a, r) { return a + r.usedU; }, 0)),
        A.num(rooms.reduce(function (a, r) { return a + r.cabCount * 42 - r.usedU; }, 0)), '-', '-', '-']
    };
  }

  var BUILDERS = { asset: buildAsset, alarm: buildAlarm, th: buildTH, energy: buildEnergy, ups: buildUPS, cap: buildCap };

  A.register('report', {
    render: function (root) {
      root.innerHTML =
        '<div class="page-head">' +
          '<h1 class="page-title">报表报告</h1>' +
          '<span class="page-crumb">' + A.crumb('模板与生成') +
          '<div class="page-tools">' +
            '<div class="field"><span class="field-label">统计周期</span>' +
              A.UI.select('cycle', [{ v: '日', t: '日报' }, { v: '周', t: '周报' }, { v: '月', t: '月报' }, { v: '年', t: '年报' }], state.cycle) +
            '</div>' +
            '<button class="btn is-primary is-sm" data-act="gen">生成报表</button>' +
            '<button class="btn is-ghost is-sm" data-act="down">下载</button>' +
            '<button class="btn is-ghost is-sm" data-act="print">打印</button>' +
          '</div>' +
        '</div>' +

        '<div class="row grow">' +
          '<section class="card" style="flex:0 0 340px">' +
            '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">报表模板</h2></div>' +
            '<div class="card-body" style="overflow:auto;gap:8px">' +
              TPLS.map(function (t) {
                return '<div class="tpl-item' + (t.id === state.tpl ? ' is-active' : '') + '" data-tpl="' + t.id + '">' +
                  '<div class="tpl-hd"><i class="card-bar" style="height:12px"></i><span class="tpl-nm">' + A.esc(t.name) + '</span>' +
                    '<span class="tag" style="margin-left:auto">' + t.cycle + '</span></div>' +
                  '<div class="tpl-sub"><span>' + A.esc(t.desc) + '</span></div>' +
                  '<div class="tpl-ft"><span class="tag is-green">已启用</span><span class="tag is-grey">' + A.esc(t.owner) + '</span></div>' +
                  '</div>';
              }).join('') +
            '</div>' +
          '</section>' +

          '<section class="card grow" style="flex:1 1 auto;min-width:0">' +
            '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">报表预览</h2>' +
              '<div class="card-extra"><span class="muted" style="font-size:12px" id="rptMeta"></span></div></div>' +
            '<div class="card-body rpt-preview" id="rptView"></div>' +
          '</section>' +
        '</div>' +

        '<section class="card" style="flex:0 0 220px">' +
          '<div class="card-head"><i class="card-bar"></i><h2 class="card-title">定时任务</h2>' +
            '<div class="card-extra"><button class="btn is-primary is-sm" data-act="newtask">+ 新建任务</button></div></div>' +
          '<div class="card-body">' + A.UI.table({
            cols: [
              { t: '任务名称', k: 'n', w: 220 }, { t: '报表模板', k: 't', w: 160 },
              { t: '周期', k: 'c', w: 90 }, { t: '推送对象', k: 'r', w: 160 },
              { t: '上次执行', k: 'l', w: 170, mono: true }, { t: '下次执行', k: 'x', w: 170, mono: true },
              { t: '状态', k: 's', w: 100, html: function () { return A.UI.tag('运行中', 'green'); } }
            ],
            rows: TPLS.map(function (t, i) {
              return { n: t.name + '_自动任务', t: t.name, c: t.cycle, r: '运维中心 / 技术保障科',
                       l: '2026-10-19 0' + (i + 1) + ':00:00', x: '2026-10-20 0' + (i + 1) + ':00:00' };
            })
          }) + '</div>' +
        '</section>';

      this.paint(root);
      this.bind(root);
    },

    paint: function (root) {
      var d = BUILDERS[state.tpl]();
      var tpl = TPLS.filter(function (t) { return t.id === state.tpl; })[0];

      A.q('#rptMeta', root).textContent =
        tpl.name + ' · ' + state.cycle + '报 · 生成时间 ' + A.now();

      A.q('#rptView', root).innerHTML =
        '<div class="rpt-h">' + A.esc(d.title) + '</div>' +
        '<div class="rpt-meta">统计周期：2026-10-01 ~ 2026-10-19 · 数据来源：' + A.esc(A.brand() + '（' + window.DCIM.meta.project + '）') + ' · 单位：见表头</div>' +
        '<table class="rpt-tbl"><thead><tr>' +
          d.cols.map(function (c) { return '<th>' + A.esc(c) + '</th>'; }).join('') +
        '</tr></thead><tbody>' +
          d.rows.map(function (r) {
            return '<tr>' + r.map(function (c, i) {
              return '<td class="' + (typeof c === 'number' || (i > 1 && /^[\d.,]+%?$/.test(String(c))) ? 'n' : '') + '">' + A.esc(c) + '</td>';
            }).join('') + '</tr>';
          }).join('') +
        '<tr style="font-weight:600"><td>合计 / 平均</td>' +
          d.sum.slice(1).map(function (c) { return '<td class="n">' + A.esc(c) + '</td>'; }).join('') +
        '</tr>' +
        '</tbody></table>' +
        '<div style="margin-top:18px;display:flex;gap:16px;flex-wrap:wrap">' +
          (d.title.indexOf('告警') >= 0
            ? '<div style="flex:1 1 420px;height:180px">' + A.Charts.bars({
                w: 420, h: 180, data: window.DCIM.alarmStats.byLevel.map(function (x) { return { label: x.name, value: x.value }; })
              }) + '</div>'
            : '<div style="flex:1 1 420px;height:180px">' + A.Charts.line({
                w: 420, h: 180, yLabels: false, labels: ['10-01', '10-04', '10-07', '10-10', '10-13', '10-16', '10-19'],
                series: [{ color: '#00D8FF', area: true, data: A.wave(60, 24, 12, d.title.length) }]
              }) + '</div>') +
          '<div style="flex:1 1 300px;font-size:12px;color:var(--t-muted);line-height:1.9">' +
            '制表：系统自动生成<br>审核：运维值班员<br>报送：技术保障科<br>' +
            '说明：本报表由平台按周期自动汇总，数据口径与实时监控一致。' +
          '</div>' +
        '</div>';
    },

    bind: function (root) {
      var self = this;
      A.qa('.tpl', root).forEach(function (t) {
        t.addEventListener('click', function () {
          state.tpl = t.getAttribute('data-tpl');
          A.qa('.tpl', root).forEach(function (x) { x.classList.toggle('is-active', x === t); });
          self.paint(root);
        });
      });
      A.q('[data-f="cycle"]', root).addEventListener('change', function () {
        state.cycle = this.value; self.paint(root);
      });
      A.q('[data-act="gen"]', root).addEventListener('click', function () {
        A.toast('报表已重新生成');
        self.paint(root);
      });
      A.q('[data-act="down"]', root).addEventListener('click', function () { A.toast('报表已下载（演示）'); });
      A.q('[data-act="print"]', root).addEventListener('click', function () { A.toast('已发送到打印队列'); });
      A.q('[data-act="newtask"]', root).addEventListener('click', function () { A.toast('定时任务创建向导（演示）'); });
    }
  });
})(window.App);
