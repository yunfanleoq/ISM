/* 配电专业：单线图 / 定值 / 电能质量 / 录波 / 操作票 / 仿真 */
(function (A) {
  'use strict';
  var SUBS = [
    { id: 'line', name: '单线图' },
    { id: 'set', name: '定值' },
    { id: 'pq', name: '电能质量' },
    { id: 'wave', name: '录波分析' },
    { id: 'ticket', name: '操作票' },
    { id: 'sim', name: '仿真' }
  ];
  var NOTES = {
    line: '本站简化一次系统：10kV 进线到变压器、0.4kV 母线、UPS 和列头柜。这不是 IEEE33，也不是竣工图扫描件。电压电流功率只显示已接入的实时测点。',
    set: '按保护装置核对当前值、定值单和上次整定。相对定值单偏差超过 5% 标成超差。这里只校核和留痕，不向保护装置下装。',
    pq: '电压偏差、频率、不平衡、谐波、闪变和功率因数按国标判限。这组数是评估示例，不是电能质量装置的现场采集。',
    wave: '示例录波，用来看三相波形、相量、序分量和故障类型建议。不是现场保护装置导出的 COMTRADE 文件。',
    ticket: '操作票从待审核到批准、逐步执行、归档。执行前核对设备状态，防止误分误合。',
    sim: '本站短路和灵敏度是简化演算，改变压器容量会重算。IEEE33 是标准算例，非本站拓扑。结果不能当成现场故障电流。'
  };
  var S = { sub: 'line', simTab: 'station' };
  var TICKET = { review: '待审核', approved: '已批准', executing: '执行中', archived: '已归档' };

  function head() {
    return '<div class="page-head"><h1 class="page-title">配电专业</h1>' +
      '<span class="page-crumb">' + A.crumb('本站拓扑与 IEEE33 分开') + '</span></div>' +
      '<p class="page-lead">这一页给配电运维用：看供配电一次系统，核对保护定值，看电能质量和录波，办操作票，做短路与灵敏度校核。没有采集源的计算和录波都标明演示。IEEE33 只用来校验算法，不是本站接线。</p>' +
      '<div class="subtabs">' + SUBS.map(function (s) {
        return '<button class="subtab' + (s.id === S.sub ? ' is-active' : '') + '" data-sub="' + s.id + '">' + s.name + '</button>';
      }).join('') + '</div><div id="pwBody" class="grow" style="min-height:0;overflow:auto"></div>';
  }

  function note() {
    return '<p class="sub-note">' + A.esc(NOTES[S.sub] || '') + '</p>';
  }

  function lineSvg() {
    return '<svg viewBox="0 0 920 250" width="100%" height="230">' +
      '<text x="16" y="22" fill="#7FE9FF" font-size="13">市电 I</text>' +
      '<line x1="70" y1="18" x2="180" y2="18" stroke="#00D8FF" stroke-width="2"/>' +
      '<rect x="180" y="4" width="64" height="28" fill="none" stroke="#00D8FF"/><text x="188" y="23" fill="#E6F3FF" font-size="12">进线</text>' +
      '<line x1="244" y1="18" x2="320" y2="18" stroke="#00D8FF"/>' +
      '<rect x="320" y="4" width="72" height="28" fill="none" stroke="#00D8FF"/><text x="328" y="23" fill="#E6F3FF" font-size="12">T1</text>' +
      '<line x1="392" y1="18" x2="470" y2="70" stroke="#00D8FF" stroke-width="2"/>' +
      '<text x="16" y="128" fill="#7FE9FF" font-size="13">市电 II</text>' +
      '<line x1="70" y1="124" x2="180" y2="124" stroke="#7FE9FF" stroke-width="2"/>' +
      '<rect x="180" y="110" width="64" height="28" fill="none" stroke="#7FE9FF"/><text x="188" y="129" fill="#E6F3FF" font-size="12">进线</text>' +
      '<line x1="244" y1="124" x2="320" y2="124" stroke="#7FE9FF"/>' +
      '<rect x="320" y="110" width="72" height="28" fill="none" stroke="#7FE9FF"/><text x="328" y="129" fill="#E6F3FF" font-size="12">T2</text>' +
      '<line x1="392" y1="124" x2="470" y2="90" stroke="#7FE9FF" stroke-width="2"/>' +
      '<line x1="470" y1="70" x2="760" y2="70" stroke="#00D8FF" stroke-width="3"/>' +
      '<line x1="470" y1="90" x2="760" y2="90" stroke="#7FE9FF" stroke-width="3"/>' +
      '<line x1="620" y1="70" x2="620" y2="90" stroke="#FFD24D" stroke-width="2"/>' +
      '<text x="590" y="64" fill="#FFD24D" font-size="11">母联</text>' +
      '<text x="480" y="60" fill="#7FA8D9" font-size="11">0.4kV I 段</text>' +
      '<text x="480" y="108" fill="#7FA8D9" font-size="11">0.4kV II 段</text>' +
      '<rect x="760" y="54" width="70" height="28" fill="none" stroke="#C8A24D"/><text x="772" y="73" fill="#E6F3FF" font-size="12">UPS-A</text>' +
      '<rect x="760" y="100" width="70" height="28" fill="none" stroke="#C8A24D"/><text x="772" y="119" fill="#E6F3FF" font-size="12">UPS-B</text>' +
      '<text x="848" y="73" fill="#E6F3FF" font-size="12">列头柜 A</text>' +
      '<text x="848" y="119" fill="#E6F3FF" font-size="12">列头柜 B</text>' +
      '<text x="16" y="190" fill="#5E7FA8" font-size="12">本图是本站简化单线，不是 IEEE33，也不是竣工图。短路计算目前按单台变压器简化，见仿真。</text>' +
      '<text x="16" y="214" fill="#5E7FA8" font-size="12">下面的电压、电流、功率来自已接入测点；没有测点就写未采集，不用算例数填上。</text>' +
      '</svg>';
  }

  function pointOf(d, keys) {
    var pts = d.points || [];
    for (var i = 0; i < pts.length; i++) {
      var name = pts[i].name || '';
      for (var k = 0; k < keys.length; k++) {
        if (name.indexOf(keys[k]) >= 0) return pts[i];
      }
    }
    return null;
  }

  function fmtPoint(p) {
    if (!p || p.value === '' || p.value == null) return '未采集';
    return String(p.value) + (p.unit ? ' ' + p.unit : '');
  }

  function liveRows(kind) {
    var devs = (window.DCIM && DCIM.live && DCIM.live.devices) || [];
    return devs.filter(function (d) { return d.kind === kind; }).slice(0, 6);
  }

  function liveTable(title, kind) {
    var rows = liveRows(kind);
    var all = ((window.DCIM && DCIM.live && DCIM.live.devices) || []).filter(function (d) { return d.kind === kind; });
    if (!rows.length) {
      return '<section class="card" style="margin-top:12px"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">' + A.esc(title) + '</h2></div>' +
        '<div class="card-body"><div class="muted">当前没有这类实时设备。单线图仍按本站简化拓扑显示，不拿演示数冒充采集。</div></div></section>';
    }
    return '<section class="card" style="margin-top:12px"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">' + A.esc(title) + '</h2>' +
      '<div class="card-extra"><span class="muted">实时测点' + (all.length > rows.length ? '，列出前 ' + rows.length + ' 台 / 共 ' + all.length + ' 台' : '') + '</span></div></div>' +
      '<div class="card-body"><table class="tbl"><thead><tr><th>设备</th><th>电压</th><th>电流</th><th>功率</th></tr></thead><tbody>' +
      rows.map(function (d) {
        return '<tr><td>' + A.esc(d.name) + '</td><td class="mono">' + A.esc(fmtPoint(pointOf(d, ['电压']))) + '</td>' +
          '<td class="mono">' + A.esc(fmtPoint(pointOf(d, ['电流']))) + '</td>' +
          '<td class="mono">' + A.esc(fmtPoint(pointOf(d, ['有功', '功率']))) + '</td></tr>';
      }).join('') + '</tbody></table></div></section>';
  }

  function poly(data, color) {
    if (!data || !data.length) return '';
    var min = Math.min.apply(null, data), max = Math.max.apply(null, data), span = (max - min) || 1;
    var w = 640, h = 140;
    var d = data.map(function (v, i) {
      var x = (data.length === 1 ? 0 : i / (data.length - 1)) * w;
      var y = h - ((v - min) / span) * (h - 8) - 4;
      return (i ? 'L' : 'M') + x.toFixed(1) + ' ' + y.toFixed(1);
    }).join(' ');
    return '<path d="' + d + '" fill="none" stroke="' + color + '" stroke-width="1.6"/>';
  }

  A.register('power', {
    render: function (root) {
      if (A.subArg && SUBS.some(function (s) { return s.id === A.subArg; })) S.sub = A.subArg;
      root.innerHTML = head();
      this.paint(root);
      var self = this;
      A.qa('.subtab', root).forEach(function (b) {
        b.addEventListener('click', function () {
          S.sub = b.getAttribute('data-sub');
          self.render(root);
        });
      });
    },
    paint: function (root) {
      var box = A.q('#pwBody', root);
      var self = this;
      if (S.sub === 'line') {
        box.innerHTML = note() + '<section class="card"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">本站供电简图</h2>' +
          '<div class="card-extra"><span class="muted">标准算例不画在这里</span></div></div><div class="card-body">' + lineSvg() + '</div></section>' +
          liveTable('UPS 实时', 'ups') + liveTable('列头柜实时', 'pwr');
        return;
      }
      box.innerHTML = note() + '<div class="empty">加载中</div>';
      if (S.sub === 'set') A.api('/bjga/power/settings').then(function (j) { self.settings(box, j.settings || []); });
      else if (S.sub === 'pq') A.api('/bjga/power/quality').then(function (j) { self.quality(box, j.points || []); });
      else if (S.sub === 'wave') A.api('/bjga/power/comtrade').then(function (j) { self.wave(box, j); });
      else if (S.sub === 'ticket') A.api('/bjga/power/tickets').then(function (j) { self.tickets(box, j.tickets || []); });
      else A.api('/bjga/power/sim').then(function (j) { self.sim(box, j); });
    },
    settings: function (box, rows) {
      var bad = rows.filter(function (r) { return r.mismatch; }).length;
      var verdict = bad ? bad + ' 项超差，需要核对定值单，不要直接当现场已整定' : '当前列出的定值与定值单一致';
      box.innerHTML = note() + '<section class="card"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">定值校核</h2>' +
        '<div class="card-extra"><span class="muted">进线 / 馈线 / 母联</span></div></div><div class="card-body">' +
        '<div style="margin-bottom:8px;color:' + (bad ? '#FF4D5E' : '#00E39A') + '">' + A.esc(verdict) + '</div>' +
        '<table class="tbl"><thead><tr><th>装置</th><th>项目</th><th>当前</th><th>定值单</th><th>上次</th><th>结论</th><th></th></tr></thead><tbody>' +
        rows.map(function (r) {
          return '<tr><td>' + A.esc(r.device) + '</td><td>' + A.esc(r.item) + '</td><td class="mono">' + r.current + ' ' + A.esc(r.unit) + '</td>' +
            '<td class="mono">' + r.sheet + '</td><td class="mono">' + r.last + '</td><td>' +
            (r.mismatch ? '<b style="color:#FF4D5E">超差</b>' : '<span style="color:#00E39A">一致</span>') + '</td><td>' +
            (r.mismatch ? '<button class="btn is-ghost is-sm" data-id="' + r.id + '" data-sheet="' + r.sheet + '">回写为定值单</button>' : '') +
            '</td></tr>';
        }).join('') + '</tbody></table><div class="muted" style="margin-top:8px">回写只改本系统里的当前值并记审计，保护装置本身不会被下装。</div></div></section>';
      var self = this;
      A.qa('[data-id]', box).forEach(function (b) {
        b.addEventListener('click', function () {
          A.api('/bjga/power/settings', { method: 'POST', body: { id: Number(b.getAttribute('data-id')), current: Number(b.getAttribute('data-sheet')) } }).then(function (res) {
            if (res.code !== 1000) { A.toast(res.message || '未写入'); return; }
            A.toast('已按定值单回写，并写入审计');
            self.settings(box, res.settings || []);
          });
        });
      });
    },
    quality: function (box, rows) {
      var bad = rows.filter(function (r) { return r.over; }).length;
      box.innerHTML = note() + '<section class="card"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">电能质量评估表</h2>' +
        '<div class="card-extra"><span class="tag is-orange">评估示例</span></div></div><div class="card-body">' +
        '<div style="margin-bottom:8px;color:' + (bad ? '#FF4D5E' : '#00E39A') + '">' +
        (bad ? bad + ' 项越限。示例结论不能当作本站电能质量合格证。' : '列出的示例指标均在限值内。') + '</div>' +
        '<table class="tbl"><thead><tr><th>指标</th><th>示例值</th><th>限值</th><th>依据</th><th>结论</th></tr></thead><tbody>' +
        rows.map(function (r) {
          var limit = r.name === '功率因数' ? '不低于 ' + r.limit : r.limit;
          return '<tr><td>' + A.esc(r.name) + '</td><td class="mono">' + r.value + ' ' + A.esc(r.unit || '') + '</td><td class="mono">' + limit + '</td><td>' + A.esc(r.std) + '</td><td>' +
            (r.over ? '<b style="color:#FF4D5E">越限</b>' : '<span style="color:#00E39A">合格</span>') + '</td></tr>';
        }).join('') + '</tbody></table></div></section>';
    },
    wave: function (box, j) {
      var chs = j.channels || [];
      var colors = ['#00D8FF', '#FFA24D', '#00E39A'];
      var svg = chs.map(function (c, i) { return poly(c.data, colors[i % 3]); }).join('');
      var ph = (j.phasors || []).map(function (p, i) {
        var rad = p.ang * Math.PI / 180;
        var len = Math.max(20, Math.min(70, p.mag));
        var x = 80 + Math.cos(rad) * len, y = 80 - Math.sin(rad) * len;
        return '<line x1="80" y1="80" x2="' + x.toFixed(1) + '" y2="' + y.toFixed(1) + '" stroke="' + colors[i % 3] + '" stroke-width="2"/>' +
          '<text x="' + (x + 4).toFixed(0) + '" y="' + y.toFixed(0) + '" fill="#E6F3FF" font-size="11">' + A.esc(p.name) + ' ' + p.mag + '∠' + p.ang + '</text>';
      }).join('');
      var seq = j.sequence || {};
      box.innerHTML = note() + '<section class="card"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">' + A.esc(j.name || '录波') + '</h2>' +
        '<div class="card-extra"><span class="tag is-orange">演示</span> <button class="btn is-primary is-sm" id="pwLoad">载入示例</button></div></div>' +
        '<div class="card-body"><div style="color:#FFD24D;margin-bottom:8px">' + A.esc(j.demoNote || '') + ' 定性：' + A.esc(j.conclusion || '') + '</div>' +
        '<svg viewBox="0 0 640 140" width="100%" height="150">' + svg + '</svg>' +
        '<div style="display:flex;gap:24px;align-items:flex-start;flex-wrap:wrap"><svg viewBox="0 0 220 160" width="220" height="160"><circle cx="80" cy="80" r="70" fill="none" stroke="rgba(0,216,255,.3)"/>' + ph + '</svg>' +
        '<div><div>序分量（由上面的相量算出，演示判据）</div><div style="margin-top:6px">零序 ' + (seq.zero != null ? seq.zero : '-') +
        ' · 正序 ' + (seq.positive != null ? seq.positive : '-') + ' · 负序 ' + (seq.negative != null ? seq.negative : '-') + '</div>' +
        '<div style="margin-top:6px;color:#FFD24D">' + A.esc(seq.hint || '') + '</div>' +
        '<div class="muted" style="margin-top:10px">粘贴三相采样，每行 va,vb,vc。解析的仍是演示数据，不会变成现场录波。</div>' +
        '<textarea id="pwTxt" style="width:360px;height:80px;background:#071228;border:1px solid rgba(28,90,158,.8);color:#E6F3FF"></textarea>' +
        '<div style="margin-top:8px"><button class="btn is-ghost is-sm" id="pwUp">解析上传</button></div></div></div></div></section>';
      var self = this;
      A.q('#pwLoad', box).addEventListener('click', function () { self.paint(box.parentNode); });
      A.q('#pwUp', box).addEventListener('click', function () {
        A.api('/bjga/power/comtrade', { method: 'POST', body: { name: '手工粘贴采样', text: A.q('#pwTxt', box).value } }).then(function (res) {
          if (res.code && res.code !== 1000) { A.toast(res.message || '解析失败'); return; }
          self.wave(box, res);
        });
      });
    },
    tickets: function (box, list) {
      box.innerHTML = note() + (list.map(function (t) {
        var steps = (t.steps || []).map(function (s) {
          return '<div style="margin:6px 0"><button class="btn is-ghost is-sm" data-step="' + s.seq + '" data-no="' + A.esc(t.no) + '"' +
            (s.done ? ' disabled' : '') + '>' + (s.done ? '已执行' : '执行') + '</button> ' + s.seq + '. ' + A.esc(s.text) + '</div>';
        }).join('');
        return '<section class="card" style="margin-bottom:12px"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">' + A.esc(t.no) + ' ' + A.esc(t.title) + '</h2>' +
          '<div class="card-extra">' + A.esc(TICKET[t.status] || t.status) + ' · ' + A.esc(t.device) + ' 当前' + A.esc(t.deviceState) + ' · 拟票 ' + A.esc(t.author || '') + '</div></div><div class="card-body">' +
          '<div class="muted" style="margin-bottom:6px">批准后才能逐步执行。设备状态和票面不符时，后端会拒绝。</div>' +
          (t.status === 'review' ? '<button class="btn is-primary is-sm" data-ok="' + A.esc(t.no) + '">批准</button>' : '') +
          steps + '</div></section>';
      }).join('') || '<div class="empty">没有操作票</div>');
      var self = this;
      function act(no, action, seq) {
        A.api('/bjga/power/tickets', { method: 'POST', body: { no: no, action: action, seq: seq || 0 } }).then(function (j) {
          if (j.code !== 1000) { A.toast(j.message || '未通过'); return; }
          A.toast(action === 'approve' ? '已批准，可逐步执行' : '已记录一步');
          self.tickets(box, j.tickets || []);
        });
      }
      A.qa('[data-ok]', box).forEach(function (b) {
        b.addEventListener('click', function () { act(b.getAttribute('data-ok'), 'approve'); });
      });
      A.qa('[data-step]', box).forEach(function (b) {
        b.addEventListener('click', function () { act(b.getAttribute('data-no'), 'execute', Number(b.getAttribute('data-step'))); });
      });
    },
    sim: function (box, j) {
      var st = j.station || {}, calc = st.calc || {};
      var ieee = j.ieee33 || {};
      var nodes = (ieee.nodes || []).map(function (n) {
        return '<span style="display:inline-block;width:92px">N' + n.n + ' ' + n.vPu + '</span>';
      }).join('');
      var tabs = '<div class="subtabs" style="margin-bottom:10px">' +
        '<button class="subtab' + (S.simTab === 'station' ? ' is-active' : '') + '" data-sim="station">短路电流</button>' +
        '<button class="subtab' + (S.simTab === 'safety' ? ' is-active' : '') + '" data-sim="safety">安全系数</button>' +
        '<button class="subtab' + (S.simTab === 'ieee' ? ' is-active' : '') + '" data-sim="ieee">IEEE33 标准算例</button></div>';
      var body = '';
      if (S.simTab === 'ieee') {
        body = '<section class="card"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">IEEE33</h2>' +
          '<div class="card-extra"><span class="tag is-orange">标准算例，非本站拓扑</span></div></div><div class="card-body">' +
          '<div style="color:#FFD24D">' + A.esc(ieee.note || '标准算例，非本站拓扑') + '</div>' +
          '<div class="muted" style="margin:8px 0">33 节点辐射状配网，只用于看算法是否沿馈线掉电压。不能拿来当本站单线图或定值依据。</div>' +
          '<div>' + nodes + '</div></div></section>';
      } else if (S.simTab === 'safety') {
        var sensBad = String(calc.conclusion || '').indexOf('不足') >= 0;
        body = '<section class="card"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">灵敏度、载流量与供电连续性</h2>' +
          '<div class="card-extra"><span class="tag is-orange">演示演算</span></div></div><div class="card-body">' +
          '<div>保护灵敏度 ' + (calc.sensitivity != null ? calc.sensitivity : '-') + '（Isc / 整定 ' + (st.isetA || '-') + ' A，低于 1.5 为不足）</div>' +
          '<div style="margin-top:6px;color:' + (sensBad ? '#FF4D5E' : '#00E39A') + '">' + A.esc(calc.conclusion || '') + '</div>' +
          '<div style="margin-top:10px">演示母线载流量 ' + (calc.ampacityA || '-') + ' A，额定电流 ' + (calc.inA || '-') + ' A，裕度 ' + (calc.marginPct != null ? calc.marginPct + '%' : '-') + '</div>' +
          '<div>' + A.esc(calc.ampacityConclusion || '') + '。' + A.esc(calc.ampacityNote || '') + '</div>' +
          '<div style="margin-top:10px;color:#FFD24D">' + A.esc(calc.continuity || '') + '</div>' +
          '<div class="muted" style="margin-top:8px">' + A.esc(calc.model || '') + '</div></div></section>';
      } else {
        body = '<section class="card"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">' + A.esc(st.name || '本站') + '</h2>' +
          '<div class="card-extra"><span class="tag is-orange">演示演算</span></div></div><div class="card-body">' +
          '<div>' + A.esc(st.note || '') + '</div>' +
          '<div style="margin:10px 0">变压器容量 <input id="pwKva" class="input" value="' + (st.kva || 1600) + '" style="width:100px;min-width:0"/> kVA ' +
          '<button class="btn is-primary is-sm" id="pwCalc">重算短路</button></div>' +
          '<div>额定电流 ' + (calc.inA || '-') + ' A · 短路电流 ' + (calc.iscA || '-') + ' A · 灵敏度 ' + (calc.sensitivity || '-') + '</div>' +
          '<div style="margin-top:8px;color:' + (String(calc.conclusion || '').indexOf('不足') >= 0 ? '#FF4D5E' : '#00E39A') + '">' + A.esc(calc.conclusion || '') + '</div>' +
          '<div class="muted" style="margin-top:8px">' + A.esc(calc.model || '') + '</div></div></section>';
      }
      box.innerHTML = note() + tabs + body;
      var self = this;
      A.qa('[data-sim]', box).forEach(function (b) {
        b.addEventListener('click', function () { S.simTab = b.getAttribute('data-sim'); self.sim(box, j); });
      });
      var btn = A.q('#pwCalc', box);
      if (btn) btn.addEventListener('click', function () {
        A.api('/bjga/power/sim', { method: 'POST', body: { kva: Number(A.q('#pwKva', box).value) } }).then(function (res) {
          if (res.code !== 1000) { A.toast(res.message || '计算失败'); return; }
          self.sim(box, res);
        });
      });
    }
  });
})(window.App);
