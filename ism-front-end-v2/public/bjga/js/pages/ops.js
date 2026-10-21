/* 运维中心：值班、报平安、知识库、变更、热备、告警通道 */
(function (A) {
  'use strict';
  var SUBS = [
    { id: 'duty', name: '值班排班' },
    { id: 'peace', name: '报平安' },
    { id: 'kb', name: '知识库' },
    { id: 'change', name: '需求变更' },
    { id: 'ha', name: '热备' },
    { id: 'ch', name: '告警通道' }
  ];
  var S = { sub: 'duty', kind: '全部', kb: '' };

  function lead() {
    return '<p class="page-lead">运维中心给当班用：看一周班表和交接，按机房报平安，查动环处置文章，登记质保期内的需求变更，看四台平台服务器的南向、双机热备和备份，核对告警通知通道。短信网关还没接入，页面上不会假装已经发出短信。</p>';
  }

  function tag(text, tone) {
    return '<span class="tag' + (tone ? ' is-' + tone : '') + '">' + A.esc(text) + '</span>';
  }

  function stateTag(s) {
    if (s === '在岗') return tag(s, 'green');
    if (s === '待接班') return tag(s, 'orange');
    if (s === '含异常') return tag(s, 'red');
    if (s === '全部平安' || s === '已完成' || s === '正常') return tag(s, 'green');
    if (s === '未对接' || s === '未启用') return tag(s, 'grey');
    if (s === '审核中' || s === '实施中') return tag(s, 'orange');
    return tag(s || '-');
  }

  function haRole(r) {
    if (r.kind === 'south' || r.role === '南向') return tag('南向采集', 'blue');
    if (r.kind === 'backup' || r.role === '备份') return tag('备份 · 不参与热备', 'grey');
    if (r.role === '主') return tag('热备主 · 承担业务', 'green');
    if (r.role === '备') return tag('热备备 · 故障后顶上', 'orange');
    return tag(r.role || '-', 'blue');
  }

  A.register('ops', {
    render: function (root) {
      if (A.subArg && SUBS.some(function (s) { return s.id === A.subArg; })) S.sub = A.subArg;
      root.innerHTML = '<div class="page-head"><h1 class="page-title">运维中心</h1>' +
        '<span class="page-crumb">' + A.crumb('值班 · 报平安 · 知识库 · 变更') + '</span></div>' +
        lead() +
        '<div class="subtabs">' + SUBS.map(function (s) {
          return '<button class="subtab' + (s.id === S.sub ? ' is-active' : '') + '" data-sub="' + s.id + '">' + s.name + '</button>';
        }).join('') + '</div><div id="opBody" class="grow" style="overflow:auto"></div>';
      this.paint(root);
      var self = this;
      A.qa('.subtab', root).forEach(function (b) {
        b.addEventListener('click', function () { S.sub = b.getAttribute('data-sub'); self.render(root); });
      });
    },
    paint: function (root) {
      var box = A.q('#opBody', root);
      var self = this;
      A.api('/bjga/ops').then(function (j) {
        if (j.code !== 1000) { box.innerHTML = '<div class="empty">加载失败</div>'; return; }
        if (S.sub === 'duty') self.duty(box, j);
        else if (S.sub === 'peace') self.peace(root, box, j);
        else if (S.sub === 'kb') self.kb(box, j.knowledge || []);
        else if (S.sub === 'change') self.changes(root, box, j.changes || []);
        else if (S.sub === 'ha') self.ha(root, box, j.ha || []);
        else self.channels(box, j.channels || []);
      });
    },
    duty: function (box, j) {
      var cur = j.onDuty || {};
      var nxt = j.nextDuty || {};
      var rows = j.duty || [];
      box.innerHTML = '<p class="sub-note">班表按服务器日期排 7 天，白班 08:00–20:00，夜班 20:00–次日 08:00。高亮是今天。人名是班组岗位，不是私人通讯录。</p>' +
        '<section class="card" style="margin-bottom:12px"><div class="card-body">' +
        '<div>服务器日期 <b>' + A.esc(j.today || '') + '</b> · 当前是 <b>' + A.esc(j.shiftNow || '') + '</b></div>' +
        '<div style="margin-top:8px">现在在班：<b>' + A.esc(cur.who || '未排') + '</b> ' + A.esc(cur.post || '') +
        ' · ' + A.esc(cur.phone || '') + ' · ' + A.esc(cur.scope || '') + '</div>' +
        '<div style="margin-top:6px">下一班：<b>' + A.esc(nxt.day || '') + ' ' + A.esc(nxt.shift || '') + ' ' + A.esc(nxt.who || '未排') +
        '</b> · ' + A.esc(nxt.phone || '') + ' · ' + A.esc(nxt.scope || '') + '</div></div></section>' +
        table('一周班表', ['日期', '班次', '岗位', '值班人', '电话', '负责范围', '交接备注', '状态'], rows.map(function (r) {
          return {
            today: r.day === j.today,
            cells: [r.day, r.shift, r.post || '-', r.who, r.phone || '-', r.scope || '-', r.note || '-', raw(stateTag(r.state))]
          };
        })) +
        '<section class="card" style="margin-top:12px"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">今日交接班</h2></div><div class="card-body">' +
        ((j.handover || []).map(function (r) {
          return '<div style="margin:6px 0"><b>' + A.esc(r.shift) + ' ' + A.esc(r.who) + '</b> · ' + A.esc(r.note || '') + '</div>';
        }).join('') || '<div class="muted">今天没有交接备注</div>') +
        '<div class="muted" style="margin-top:8px">巡检结果不写在这里，到「报平安」按机房提交。</div></div></section>';
    },
    peace: function (root, box, j) {
      var list = j.checklist || [];
      var records = j.peace || [];
      var latest = records[0];
      box.innerHTML = '<p class="sub-note">按机房报本班巡检。每一项选平安或异常；异常必须写一句人能看懂的说明。提交后记在当前登录人的报平安记录里，并写入审计。</p>' +
        '<section class="card" style="margin-bottom:12px"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">本班巡检</h2>' +
        '<div class="card-extra"><button class="btn is-primary is-sm" id="opPeace">提交报平安</button></div></div><div class="card-body"><table class="tbl"><thead><tr>' +
        '<th>机房</th><th>看什么</th><th>结论</th><th>异常说明</th></tr></thead><tbody>' +
        list.map(function (c, i) {
          return '<tr><td>' + A.esc(c.room) + '</td><td>' + A.esc(c.focus || '') + '</td><td>' +
            '<label style="margin-right:10px"><input type="radio" name="pr' + i + '" value="平安" checked/> 平安</label>' +
            '<label><input type="radio" name="pr' + i + '" value="异常"/> 异常</label></td><td>' +
            '<input class="input" data-note="' + i + '" placeholder="选异常时必填" style="width:240px;min-width:0"/></td></tr>';
        }).join('') + '</tbody></table></div></section>' +
        peaceLatest(latest) + peaceHistory(records);
      var self = this;
      A.q('#opPeace', box).addEventListener('click', function () {
        var items = list.map(function (c, i) {
          var picked = box.querySelector('input[name="pr' + i + '"]:checked');
          var noteEl = box.querySelector('[data-note="' + i + '"]');
          return { room: c.room, result: picked ? picked.value : '', note: noteEl ? noteEl.value : '' };
        });
        A.api('/bjga/ops/peace', { method: 'POST', body: { to: '值班长 / 技术保障科', items: items } }).then(function (res) {
          A.toast(res.code === 1000 ? '已记录本班报平安' : (res.message || '失败'));
          if (res.code === 1000) self.paint(root);
        });
      });
    },
    kb: function (box, rows) {
      var kinds = ['全部'];
      rows.forEach(function (k) { if (kinds.indexOf(k.kind) < 0) kinds.push(k.kind); });
      var shown = rows.filter(function (k) { return S.kind === '全部' || k.kind === S.kind; });
      if (!shown.some(function (k) { return String(k.id) === String(S.kb); })) S.kb = shown[0] ? String(shown[0].id) : '';
      var cur = shown.filter(function (k) { return String(k.id) === String(S.kb); })[0];
      box.innerHTML = '<p class="sub-note">现场查的处置文章。先选分类，再点开一篇看适用场景和步骤。告警升级、操作票、PUE、容量都在这里，不另做一套文档库。</p>' +
        '<div class="subtabs" style="margin-bottom:10px">' + kinds.map(function (k) {
          return '<button class="subtab' + (k === S.kind ? ' is-active' : '') + '" data-kind="' + A.esc(k) + '">' + A.esc(k) + '</button>';
        }).join('') + '</div><div class="kb-layout"><div class="kb-side"><div class="kb-list">' +
        shown.map(function (k) {
          return '<button class="kb-item' + (String(k.id) === String(S.kb) ? ' is-active' : '') + '" data-id="' + k.id + '"><b>' +
            A.esc(k.title) + '</b><div class="muted">' + A.esc(k.kind) + '</div></button>';
        }).join('') + '</div></div><section class="card grow"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">' +
        A.esc(cur ? cur.title : '知识库') + '</h2><div class="card-extra">' + A.esc(cur ? cur.kind : '') + '</div></div><div class="card-body">' +
        (cur ? '<div style="margin-bottom:8px"><b>适用场景</b> ' + A.esc(cur.scene || '') + '</div><div class="kb-body">' + A.esc(cur.body || '') + '</div>'
          : '<div class="empty">这个分类还没有文章</div>') +
        '</div></section></div>';
      var self = this;
      A.qa('[data-kind]', box).forEach(function (b) {
        b.addEventListener('click', function () { S.kind = b.getAttribute('data-kind'); self.kb(box, rows); });
      });
      A.qa('[data-id]', box).forEach(function (b) {
        b.addEventListener('click', function () { S.kb = b.getAttribute('data-id'); self.kb(box, rows); });
      });
    },
    changes: function (root, box, rows) {
      box.innerHTML = '<p class="sub-note">质保期内的优化和运行方式调整登记在这里，留时间、状态和登记人。这是变更留痕，不是外部工单系统。</p>' +
        '<div style="margin-bottom:8px"><input id="opTitle" class="input" placeholder="变更标题" style="width:280px"/> ' +
        '<button class="btn is-primary is-sm" id="opAdd">登记</button></div>' +
        table('需求变更', ['时间', '标题', '状态', '登记人'], rows.map(function (r) {
          return { cells: [r.at, r.title, raw(stateTag(r.status)), r.owner] };
        }));
      var self = this;
      A.q('#opAdd', box).addEventListener('click', function () {
        A.api('/bjga/ops/change', { method: 'POST', body: { title: A.q('#opTitle', box).value } }).then(function (res) {
          A.toast(res.code === 1000 ? '已登记' : '失败');
          if (res.code === 1000) self.paint(root);
        });
      });
    },
    ha: function (root, box, rows) {
      box.innerHTML = '<p class="sub-note">四台里只有 2 号和 3 号是双机热备：当前角色为「主」的那台承担业务，「备」的那台在主故障后顶上。1 号接南向设备，4 号只备份数据和安装资料，不参与热备。切换是演练，只对调主备角色并写审计。</p>' +
        '<div style="margin-bottom:8px"><button class="btn is-primary is-sm" id="opSw">演练切换（二次确认）</button></div>' +
        table('热备', ['机器', '角色', '规格', '职责', '部署', '切换', '保活', '状态'], rows.map(function (r) {
          return { cells: [r.name, raw(haRole(r)), r.spec || '未列', r.duty || '', r.online || '', r.switchAt || '', r.keep || '', raw(stateTag(r.status))] };
        }));
      var self = this;
      A.q('#opSw', box).addEventListener('click', function () {
        if (!window.confirm('确认执行主备切换演练？这只对调角色，不操作现场电源。')) return;
        A.api('/bjga/ops/ha/switch', { method: 'POST', body: {} }).then(function (res) {
          A.toast(res.code === 1000 ? '角色已对调' : '失败');
          if (res.code === 1000 && window.DCIM && window.DCIM.extra) window.DCIM.extra.ha = res.ha;
          self.paint(root);
        });
      });
    },
    channels: function (box, rows) {
      box.innerHTML = '<p class="sub-note">招标要求电话、短信、声光、邮件四条告警通道。下面是现有通道和真实对接情况。短信网关未接入，本页不发送短信，也没有送达回执。</p>' +
        table('告警通知通道', ['通道', '状态', '对象', '说明'], rows.map(function (r) {
          return { cells: [r.name, raw(stateTag(r.status)), r.target, r.remark || ''] };
        })) +
        '<section class="card" style="margin-top:12px"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">升级规则</h2></div><div class="card-body">' +
        '<div>紧急告警 5 分钟未确认，电话找值班长。重要及以上 15 分钟未确认，升级技术保障科。</div>' +
        '<div class="muted" style="margin-top:8px">这是值班口径。短信通道未对接，升级不能当成已经发出短信。</div></div></section>';
    }
  });

  function peaceLatest(r) {
    if (!r) return '<section class="card" style="margin-bottom:12px"><div class="card-body">还没有报平安记录。</div></section>';
    var items = (r.items || []).map(function (it) {
      var tone = it.result === '异常' ? 'red' : 'green';
      return '<div style="margin:4px 0">' + A.esc(it.room) + ' ' + tag(it.result || '-', tone) +
        (it.note ? ' ' + A.esc(it.note) : '') + '</div>';
    }).join('');
    return '<section class="card" style="margin-bottom:12px"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">最近一次</h2>' +
      '<div class="card-extra">' + stateTag(r.status) + '</div></div><div class="card-body">' +
      '<div>' + A.esc(r.who || '未知') + ' · ' + A.esc(r.at || '') + ' · 报送 ' + A.esc(r.to || '') + '</div>' +
      '<div class="muted" style="margin:6px 0">' + A.esc(r.note || '') + '</div>' +
      (items || '<div class="muted">这条是早期摘要，没有分机房明细。</div>') + '</div></section>';
  }

  function peaceHistory(rows) {
    return table('报平安记录', ['时间', '报送人', '对象', '状态', '摘要'], rows.map(function (r) {
      return { cells: [r.at, r.who || '-', r.to, raw(stateTag(r.status)), r.note] };
    }));
  }

  function table(title, heads, rows) {
    return '<section class="card"><div class="card-head"><i class="card-bar"></i><h2 class="card-title">' + A.esc(title) + '</h2></div><div class="card-body"><table class="tbl"><thead><tr>' +
      heads.map(function (h) { return '<th>' + A.esc(h) + '</th>'; }).join('') + '</tr></thead><tbody>' +
      rows.map(function (r) {
        var cells = r.cells || r;
        var cls = r.today ? ' class="is-today"' : '';
        return '<tr' + cls + '>' + cells.map(function (c) {
          var text = c && c.__html != null ? c.__html : A.esc(c);
          return '<td>' + text + '</td>';
        }).join('') + '</tr>';
      }).join('') + '</tbody></table></div></section>';
  }

  function raw(html) { return { __html: html }; }
})(window.App);
